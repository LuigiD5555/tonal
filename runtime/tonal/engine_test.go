package tonal

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

// fakeRegistry proves the Engine can run with no Tlaloc package and no model.
type fakeRegistry struct {
	descriptors map[string][]CapabilityDescriptor // first = selected
	exec        func(req CapabilityExecutionRequest) (CapabilityExecutionResult, error)
}

func (f *fakeRegistry) Candidates(capability string, _ CapabilityGoal) []CapabilityCandidate {
	list := f.descriptors[strings.ToUpper(capability)]
	out := make([]CapabilityCandidate, 0, len(list))
	for index, d := range list {
		out = append(out, CapabilityCandidate{
			WorkerID:      d.ID,
			Capability:    d.Capability,
			Kind:          d.Kind,
			EngineKind:    d.EngineKind,
			ProfileRef:    d.ProfileRef,
			Deterministic: d.Deterministic,
			Generative:    d.EngineKind == "GENERATIVE",
			Selected:      index == 0,
			Reason:        "fake",
		})
	}
	return out
}

func (f *fakeRegistry) Execute(_ context.Context, req CapabilityExecutionRequest) (CapabilityExecutionResult, error) {
	return f.exec(req)
}

func observationOf(req CapabilityExecutionRequest, engine string, value any) CapabilityExecutionResult {
	raw, _ := json.Marshal(value)
	kind := "OBSERVATION"
	if req.Capability == "VERIFY" {
		kind = "FACT"
	}
	return CapabilityExecutionResult{
		WorkerID: req.WorkerID,
		Output:   raw,
		Observations: []Observation{{
			Producer: req.WorkerID, Capability: req.Capability, Key: req.NodeID,
			Value: raw, Kind: kind, Confidence: 1,
		}},
		Usage: usageFor(engine),
	}
}

func usageFor(engine string) *CapabilityUsage {
	if engine == "GENERATIVE" {
		return &CapabilityUsage{ModelCalls: 1, PromptTokens: 10, CompletionTokens: 2}
	}
	return nil
}

func depth4Family() TaskFamily {
	return TaskFamily{
		ID:   "LOCATE_EXTRACT_NORMALIZE_COMPARE",
		Goal: "read a value and compare it with a threshold",
		Steps: []Step{
			{LocalID: "locate", Capability: "LOCATE_REGION", Input: InputSpec{Template: map[string]any{
				"mode": "REAL", "question": "${param:question}", "store_dir": "${param:store_dir}",
			}}},
			{LocalID: "extract", Capability: "EXTRACT_NUMBER", DependsOn: []string{"locate"}, Input: InputSpec{Template: map[string]any{
				"image_path": "${param:page_image}", "region": "${obs:locate}",
			}}},
			{LocalID: "normalize", Capability: "NORMALIZE", DependsOn: []string{"extract"}, Input: InputSpec{Template: map[string]any{
				"raw": "${obs:extract:text}", "target_type": "number",
			}}},
			{LocalID: "compare", Capability: "COMPARE_NUMBERS", DependsOn: []string{"normalize"}, Input: InputSpec{Template: map[string]any{
				"a": "${obs:normalize:trimmed}", "b": "${param:threshold}",
			}}},
		},
	}
}

func TestCriticalPathDepthIsMechanical(t *testing.T) {
	if got := depth4Family().CriticalPathDepth(); got != 4 {
		t.Fatalf("critical path depth = %d, want 4", got)
	}
}

func TestRunWorkflow_HeterogeneousArm_ChainsObservationsAndTraces(t *testing.T) {
	registry := &fakeRegistry{
		descriptors: map[string][]CapabilityDescriptor{
			"LOCATE_REGION":   {{ID: "region-locate-tlaloque", Capability: "LOCATE_REGION", Kind: CapabilityTlaloque, EngineKind: "DETERMINISTIC", Deterministic: true}},
			"EXTRACT_NUMBER":  {{ID: "external-model:EXTRACT_NUMBER", Capability: "EXTRACT_NUMBER", Kind: CapabilityExternalModel, EngineKind: "GENERATIVE", ProfileRef: "parrot-lfm2-vl-1.6b@r1.0.0"}},
			"NORMALIZE":       {{ID: "normalize-tlaloque", Capability: "NORMALIZE", Kind: CapabilityTlaloque, EngineKind: "DETERMINISTIC", Deterministic: true}},
			"COMPARE_NUMBERS": {{ID: "numeric-tlaloque", Capability: "COMPARE_NUMBERS", Kind: CapabilityTlaloque, EngineKind: "DETERMINISTIC", Deterministic: true}},
		},
		exec: func(req CapabilityExecutionRequest) (CapabilityExecutionResult, error) {
			switch req.Capability {
			case "LOCATE_REGION":
				return observationOf(req, "DETERMINISTIC", map[string]any{"selected_address": "ohf://doc/pages/000100", "page": 100}), nil
			case "EXTRACT_NUMBER":
				return observationOf(req, "GENERATIVE", map[string]any{"text": "512"}), nil
			case "NORMALIZE":
				return observationOf(req, "DETERMINISTIC", map[string]any{"trimmed": "512", "is_number": true}), nil
			case "COMPARE_NUMBERS":
				return observationOf(req, "DETERMINISTIC", map[string]any{"comparison": "GREATER"}), nil
			}
			return CapabilityExecutionResult{}, nil
		},
	}

	engine := &Engine{Registry: registry}
	instance := Instance{
		ID: "wf-001", Family: "LOCATE_EXTRACT_NORMALIZE_COMPARE", DeclaredDepth: 4,
		Params: map[string]string{
			"question": "what is the FashionMNIST training set size", "store_dir": "/store",
			"page_image": "/img/p100.png", "threshold": "100",
		},
	}

	record, blackboard, err := engine.RunWorkflow(context.Background(), depth4Family(), instance, HeterogeneousPolicy{})
	if err != nil {
		t.Fatalf("RunWorkflow: %v", err)
	}
	if record.FinalStatus != "OK" {
		t.Fatalf("final status = %q (%s)", record.FinalStatus, record.Error)
	}
	if len(record.Steps) != 4 {
		t.Fatalf("expected 4 step traces, got %d", len(record.Steps))
	}
	if got := record.Steps[1].SelectedKind; got != CapabilityExternalModel {
		t.Fatalf("extract selected kind = %q, want EXTERNAL_MODEL", got)
	}

	normalize := record.Steps[2]
	if normalize.InputJSON == "" || !strings.Contains(normalize.InputJSON, `"512"`) {
		t.Fatalf("normalize did not receive the extracted value: %s", normalize.InputJSON)
	}
	if len(normalize.BlackboardReads) == 0 {
		t.Fatalf("normalize step recorded no blackboard reads")
	}

	if record.Accounting.ParrotCalls != 1 || record.Accounting.GenerativeCalls != 1 {
		t.Fatalf("accounting external/parrot=%d generative=%d, want 1/1", record.Accounting.ParrotCalls, record.Accounting.GenerativeCalls)
	}
	if record.Accounting.DeterministicOps != 3 {
		t.Fatalf("expected 3 deterministic ops, got %d", record.Accounting.DeterministicOps)
	}
	if len(blackboard.Keys()) != 4 {
		t.Fatalf("expected 4 blackboard keys, got %v", blackboard.Keys())
	}
}

func TestRunWorkflow_ParrotCentricArm_ForcesExternalModelForCognitiveCaps(t *testing.T) {
	registry := &fakeRegistry{
		descriptors: map[string][]CapabilityDescriptor{
			"LOCATE_REGION": {{ID: "region-locate-tlaloque", Capability: "LOCATE_REGION", Kind: CapabilityTlaloque, EngineKind: "DETERMINISTIC", Deterministic: true}},
			"NORMALIZE": {
				{ID: "normalize-tlaloque", Capability: "NORMALIZE", Kind: CapabilityTlaloque, EngineKind: "DETERMINISTIC", Deterministic: true},
				{ID: "bounded-generative-tlaloque", Capability: "NORMALIZE", Kind: CapabilityTlaloque, EngineKind: "GENERATIVE"},
				{ID: "external-model:NORMALIZE", Capability: "NORMALIZE", Kind: CapabilityExternalModel, EngineKind: "GENERATIVE"},
			},
		},
		exec: func(req CapabilityExecutionRequest) (CapabilityExecutionResult, error) {
			engine := "DETERMINISTIC"
			if strings.HasPrefix(req.WorkerID, "external-model") {
				engine = "GENERATIVE"
			}
			return observationOf(req, engine, map[string]any{"trimmed": "7"}), nil
		},
	}
	family := TaskFamily{ID: "F", Goal: "g", Steps: []Step{
		{LocalID: "locate", Capability: "LOCATE_REGION", Input: InputSpec{Template: map[string]any{"mode": "REAL", "question": "${param:q}", "store_dir": "/s"}}},
		{LocalID: "normalize", Capability: "NORMALIZE", DependsOn: []string{"locate"}, Input: InputSpec{Template: map[string]any{"raw": "7", "target_type": "number"}}},
	}}
	policy := ParrotCentricPolicy{CognitiveCapabilities: map[string]bool{"NORMALIZE": true}}

	record, _, err := (&Engine{Registry: registry}).RunWorkflow(context.Background(), family, Instance{ID: "wf", Params: map[string]string{"q": "x"}}, policy)
	if err != nil {
		t.Fatalf("RunWorkflow: %v", err)
	}
	if got := record.Steps[0].SelectedWorker; got != "region-locate-tlaloque" {
		t.Fatalf("locate should stay machinery, got %q", got)
	}
	if got := record.Steps[1].SelectedWorker; got != "external-model:NORMALIZE" {
		t.Fatalf("arm B must force NORMALIZE to EXTERNAL_MODEL, got %q", got)
	}
	if record.Accounting.ParrotCalls != 1 {
		t.Fatalf("arm B external/parrot calls = %d, want 1", record.Accounting.ParrotCalls)
	}
}

func TestRunWorkflow_NonVerifyNodeEmittingFactIsAScopeViolation(t *testing.T) {
	registry := &fakeRegistry{
		descriptors: map[string][]CapabilityDescriptor{
			"NORMALIZE": {{ID: "normalize-tlaloque", Capability: "NORMALIZE", Kind: CapabilityTlaloque, EngineKind: "DETERMINISTIC", Deterministic: true}},
		},
		exec: func(req CapabilityExecutionRequest) (CapabilityExecutionResult, error) {
			raw, _ := json.Marshal(map[string]any{"trimmed": "7"})
			return CapabilityExecutionResult{
				WorkerID: req.WorkerID, Output: raw,
				Observations: []Observation{{
					Producer: req.WorkerID, Capability: req.Capability, Key: req.NodeID,
					Value: raw, Kind: "FACT", Confidence: 1,
				}},
			}, nil
		},
	}
	family := TaskFamily{ID: "F", Goal: "g", Steps: []Step{
		{LocalID: "normalize", Capability: "NORMALIZE", Input: InputSpec{Template: map[string]any{"raw": "7", "target_type": "number"}}},
	}}
	record, _, err := (&Engine{Registry: registry}).RunWorkflow(context.Background(), family, Instance{ID: "wf"}, HeterogeneousPolicy{})
	if err != nil {
		t.Fatalf("RunWorkflow returned a hard error: %v", err)
	}
	if record.FinalStatus != "CONTRACT_FAILURE" || !strings.Contains(record.Error, "FACT_PROMOTION_SCOPE_VIOLATION") {
		t.Fatalf("status=%q error=%q, want CONTRACT_FAILURE / FACT_PROMOTION_SCOPE_VIOLATION", record.FinalStatus, record.Error)
	}
}

func TestRunWorkflow_TerminalOutputKindFollowsVerifyPresence(t *testing.T) {
	withVerify := TaskFamily{
		ID:   "V",
		Goal: "g",
		Steps: []Step{
			{LocalID: "s", Capability: "VERIFY", Input: InputSpec{Template: map[string]any{}}},
		},
	}
	if !withVerify.HasVerify() {
		t.Fatal("HasVerify should be true")
	}
	noVerify := TaskFamily{
		ID:   "N",
		Goal: "g",
		Steps: []Step{
			{LocalID: "s", Capability: "COMPARE_NUMBERS", Input: InputSpec{Template: map[string]any{}}},
		},
	}
	if noVerify.HasVerify() {
		t.Fatal("HasVerify should be false")
	}
}

func TestRunWorkflow_UnavailableCapabilityIsAContractFailureNotAPanic(t *testing.T) {
	registry := &fakeRegistry{descriptors: map[string][]CapabilityDescriptor{}, exec: func(CapabilityExecutionRequest) (CapabilityExecutionResult, error) {
		return CapabilityExecutionResult{}, nil
	}}
	family := TaskFamily{
		ID:   "F",
		Goal: "g",
		Steps: []Step{
			{LocalID: "s", Capability: "SUMMARIZE", Input: InputSpec{Template: map[string]any{}}},
		},
	}
	record, _, err := (&Engine{Registry: registry}).RunWorkflow(context.Background(), family, Instance{ID: "wf"}, HeterogeneousPolicy{})
	if err != nil {
		t.Fatalf("RunWorkflow returned a hard error: %v", err)
	}
	if record.FinalStatus != "CONTRACT_FAILURE" {
		t.Fatalf("final status = %q, want CONTRACT_FAILURE", record.FinalStatus)
	}
}
