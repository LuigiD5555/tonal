package tonal

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func TestControlLoopPrefersMachineryEvenWhenExternalModelIsAllowed(t *testing.T) {
	externalExecutions := 0
	registry := &fakeRegistry{
		descriptors: map[string][]CapabilityDescriptor{
			"ANALYZE": {
				{ID: "analysis-tlaloque", Capability: "ANALYZE", Kind: CapabilityTlaloque, EngineKind: "DETERMINISTIC", Deterministic: true},
				{ID: "parrot", Capability: "ANALYZE", Kind: CapabilityExternalModel, EngineKind: "GENERATIVE"},
			},
		},
		exec: func(req CapabilityExecutionRequest) (CapabilityExecutionResult, error) {
			if req.WorkerID == "parrot" {
				externalExecutions++
			}
			return observationOf(req, "DETERMINISTIC", map[string]any{"result": "ok"}), nil
		},
	}
	controller := FixedProgramController{Decisions: []ControlDecision{{
		Action: ControlExecute, Capability: "ANALYZE", AllowExternalModel: true,
	}}}

	run, err := (&ControlLoop{Registry: registry, MaxTransitions: 4}).Run(
		context.Background(), ControlRunRequest{RunID: "control-1", Goal: "analyze"}, controller,
	)
	if err != nil {
		t.Fatal(err)
	}
	if run.Status != ControlRunSucceeded {
		t.Fatalf("status=%q", run.Status)
	}
	if len(run.Transitions) != 1 || run.Transitions[0].Selected == nil {
		t.Fatalf("transitions=%+v", run.Transitions)
	}
	if got := run.Transitions[0].Selected.Kind; got != CapabilityTlaloque {
		t.Fatalf("selected kind=%q, want TLALOQUE", got)
	}
	if externalExecutions != 0 || run.Accounting.ExternalModelCalls != 0 {
		t.Fatalf("external cognition was used unnecessarily: executions=%d accounting=%d", externalExecutions, run.Accounting.ExternalModelCalls)
	}
}

func TestControlLoopNeverImplicitlyFallsBackToExternalModel(t *testing.T) {
	executions := 0
	registry := &fakeRegistry{
		descriptors: map[string][]CapabilityDescriptor{
			"INTERPRET": {{ID: "parrot", Capability: "INTERPRET", Kind: CapabilityExternalModel, EngineKind: "GENERATIVE"}},
		},
		exec: func(req CapabilityExecutionRequest) (CapabilityExecutionResult, error) {
			executions++
			return observationOf(req, "GENERATIVE", map[string]any{"candidate": "x"}), nil
		},
	}
	controller := FixedProgramController{Decisions: []ControlDecision{
		{Action: ControlExecute, Capability: "INTERPRET", AllowExternalModel: false},
		{Action: ControlStop, Reason: "do not escalate"},
	}}

	run, err := (&ControlLoop{Registry: registry, MaxTransitions: 4}).Run(
		context.Background(), ControlRunRequest{RunID: "control-2", Goal: "interpret"}, controller,
	)
	if err != nil {
		t.Fatal(err)
	}
	if executions != 0 {
		t.Fatalf("external model executed without permission: %d", executions)
	}
	if len(run.Transitions) != 1 || run.Transitions[0].Outcome != TransitionUnavailable {
		t.Fatalf("transitions=%+v", run.Transitions)
	}
	if !strings.Contains(run.Transitions[0].Error, "not explicitly permitted") {
		t.Fatalf("missing explicit external-cognition reason: %q", run.Transitions[0].Error)
	}
	if run.Status != ControlRunStopped {
		t.Fatalf("status=%q, want STOPPED", run.Status)
	}
}

func TestControlLoopAllowsExplicitExternalEscalationAfterMachineryMiss(t *testing.T) {
	executions := 0
	registry := &fakeRegistry{
		descriptors: map[string][]CapabilityDescriptor{
			"INTERPRET": {{ID: "parrot", Capability: "INTERPRET", Kind: CapabilityExternalModel, EngineKind: "GENERATIVE"}},
		},
		exec: func(req CapabilityExecutionRequest) (CapabilityExecutionResult, error) {
			executions++
			result := observationOf(req, "GENERATIVE", map[string]any{"candidate": "likely cause"})
			result.Usage = &CapabilityUsage{ModelCalls: 1, PromptTokens: 20, CompletionTokens: 4}
			return result, nil
		},
	}
	controller := FixedProgramController{Decisions: []ControlDecision{
		{Action: ControlExecute, Capability: "INTERPRET", AllowExternalModel: false, Reason: "try machinery first"},
		{Action: ControlExecute, Capability: "INTERPRET", AllowExternalModel: true, Reason: "explicit escalation after machinery miss"},
	}}

	run, err := (&ControlLoop{Registry: registry, MaxTransitions: 5}).Run(
		context.Background(), ControlRunRequest{RunID: "control-3", Goal: "interpret"}, controller,
	)
	if err != nil {
		t.Fatal(err)
	}
	if run.Status != ControlRunSucceeded {
		t.Fatalf("status=%q error=%q", run.Status, run.Error)
	}
	if len(run.Transitions) != 2 {
		t.Fatalf("transition count=%d", len(run.Transitions))
	}
	if run.Transitions[0].Outcome != TransitionUnavailable || run.Transitions[1].Outcome != TransitionCommitted {
		t.Fatalf("outcomes=%q,%q", run.Transitions[0].Outcome, run.Transitions[1].Outcome)
	}
	if run.Transitions[1].Selected == nil || run.Transitions[1].Selected.Kind != CapabilityExternalModel {
		t.Fatalf("explicit escalation did not select external cognition: %+v", run.Transitions[1].Selected)
	}
	if executions != 1 || run.Accounting.ExternalModelCalls != 1 {
		t.Fatalf("external calls executions=%d accounting=%d", executions, run.Accounting.ExternalModelCalls)
	}
	if len(run.Observations) != 1 {
		t.Fatalf("committed observations=%d, want 1", len(run.Observations))
	}
}

func TestControlLoopRejectedTransitionCannotMutateBlackboard(t *testing.T) {
	registry := &fakeRegistry{
		descriptors: map[string][]CapabilityDescriptor{
			"PROPOSE_HYPOTHESIS": {{ID: "hypothesis-tlaloque", Capability: "PROPOSE_HYPOTHESIS", Kind: CapabilityTlaloque, EngineKind: "SPECIALIST"}},
		},
		exec: func(req CapabilityExecutionRequest) (CapabilityExecutionResult, error) {
			raw, _ := json.Marshal(map[string]any{"hypothesis": "x"})
			return CapabilityExecutionResult{
				WorkerID: req.WorkerID,
				Output:   raw,
				Observations: []Observation{{
					Producer: req.WorkerID, Capability: req.Capability, Key: req.NodeID,
					Value: raw, Kind: "FACT", Confidence: 1,
				}},
			}, nil
		},
	}
	controller := FixedProgramController{Decisions: []ControlDecision{{
		Action: ControlExecute, Capability: "PROPOSE_HYPOTHESIS",
	}}}

	run, err := (&ControlLoop{Registry: registry}).Run(
		context.Background(), ControlRunRequest{RunID: "control-4", Goal: "form hypothesis"}, controller,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(run.Transitions) != 1 || run.Transitions[0].Outcome != TransitionRejected {
		t.Fatalf("transition=%+v", run.Transitions)
	}
	if len(run.Observations) != 0 {
		t.Fatalf("rejected transition mutated committed state: %+v", run.Observations)
	}
	if run.Accounting.Rejected != 1 || run.Accounting.Committed != 0 {
		t.Fatalf("accounting=%+v", run.Accounting)
	}
}

func TestControlLoopFeedsCommittedStateIntoNextTransition(t *testing.T) {
	registry := &fakeRegistry{
		descriptors: map[string][]CapabilityDescriptor{
			"NORMALIZE": {{ID: "normalize-tlaloque", Capability: "NORMALIZE", Kind: CapabilityTlaloque, EngineKind: "DETERMINISTIC", Deterministic: true}},
			"VERIFY":    {{ID: "verify-tlaloque", Capability: "VERIFY", Kind: CapabilityTlaloque, EngineKind: "DETERMINISTIC", Deterministic: true}},
		},
		exec: func(req CapabilityExecutionRequest) (CapabilityExecutionResult, error) {
			if req.Capability == "VERIFY" && len(req.PriorObservations) != 1 {
				t.Fatalf("VERIFY saw %d prior observations, want 1", len(req.PriorObservations))
			}
			if req.Capability == "VERIFY" {
				raw, _ := json.Marshal(map[string]any{"verified": true})
				return CapabilityExecutionResult{
					WorkerID: req.WorkerID,
					Output:   raw,
					Observations: []Observation{{
						Producer: req.WorkerID, Capability: req.Capability, Key: req.NodeID,
						Value: raw, Kind: "FACT", Confidence: 1,
					}},
				}, nil
			}
			return observationOf(req, "DETERMINISTIC", map[string]any{"normalized": 7}), nil
		},
	}
	controller := FixedProgramController{Decisions: []ControlDecision{
		{Action: ControlExecute, Capability: "NORMALIZE", PreferDeterministic: true},
		{Action: ControlExecute, Capability: "VERIFY", PreferDeterministic: true},
	}}

	run, err := (&ControlLoop{Registry: registry}).Run(
		context.Background(), ControlRunRequest{RunID: "control-5", Goal: "normalize and verify"}, controller,
	)
	if err != nil {
		t.Fatal(err)
	}
	if run.Status != ControlRunSucceeded || len(run.Observations) != 2 || run.Accounting.Committed != 2 {
		t.Fatalf("run=%+v", run)
	}
}
