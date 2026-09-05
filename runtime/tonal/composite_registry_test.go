package tonal

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

type sourceRegistryStub struct {
	candidates map[string][]CapabilityCandidate
	calls      int
	last       CapabilityExecutionRequest
	label      string
}

func (s *sourceRegistryStub) Candidates(capability string, _ CapabilityGoal) []CapabilityCandidate {
	list := s.candidates[capability]
	out := make([]CapabilityCandidate, len(list))
	copy(out, list)
	return out
}

func (s *sourceRegistryStub) Execute(_ context.Context, req CapabilityExecutionRequest) (CapabilityExecutionResult, error) {
	s.calls++
	s.last = req
	raw, _ := json.Marshal(map[string]any{"source": s.label})
	return CapabilityExecutionResult{
		WorkerID: req.WorkerID,
		Output:   raw,
		Observations: []Observation{{
			Producer:   req.WorkerID,
			Capability: req.Capability,
			Key:        req.NodeID,
			Value:      raw,
			Kind:       "OBSERVATION",
		}},
	}, nil
}

func TestCompositeRegistry_StampsSourceWithoutRewritingWorkerID(t *testing.T) {
	tlaloc := &sourceRegistryStub{label: "tlaloc", candidates: map[string][]CapabilityCandidate{
		"NORMALIZE": {{WorkerID: "shared", Capability: "NORMALIZE", Kind: CapabilityTlaloque, Deterministic: true, Selected: true}},
	}}
	parrot := &sourceRegistryStub{label: "parrot", candidates: map[string][]CapabilityCandidate{
		"NORMALIZE": {{WorkerID: "shared", Capability: "NORMALIZE", Kind: CapabilityExternalModel, Generative: true, Selected: true}},
	}}
	registry, err := NewCompositeRegistry(
		RegistrySource{ID: "tlaloc", Registry: tlaloc},
		RegistrySource{ID: "parrot", Registry: parrot},
	)
	if err != nil {
		t.Fatalf("NewCompositeRegistry: %v", err)
	}

	got := registry.Candidates("NORMALIZE", CapabilityGoal{Capability: "NORMALIZE"})
	if len(got) != 2 {
		t.Fatalf("candidates = %d, want 2", len(got))
	}
	if got[0].SourceID != "tlaloc" || got[1].SourceID != "parrot" {
		t.Fatalf("source IDs = %q, %q", got[0].SourceID, got[1].SourceID)
	}
	if got[0].WorkerID != "shared" || got[1].WorkerID != "shared" {
		t.Fatalf("composite registry rewrote source-local WorkerID: %+v", got)
	}
}

func TestMachineryFirstSelection_IsSourceAwareWithDuplicateWorkerIDs(t *testing.T) {
	candidates := []CapabilityCandidate{
		{SourceID: "parrot", WorkerID: "shared", Capability: "NORMALIZE", Kind: CapabilityExternalModel, Generative: true},
		{SourceID: "tlaloc", WorkerID: "shared", Capability: "NORMALIZE", Kind: CapabilityTlaloque, Deterministic: true},
	}
	selected, _, ok := selectCapabilityCandidate(MachineryFirstPolicy{}, Step{Capability: "NORMALIZE"}, candidates)
	if !ok {
		t.Fatal("selection failed")
	}
	if selected.SourceID != "tlaloc" || selected.WorkerID != "shared" {
		t.Fatalf("selected %+v, want tlaloc/shared", selected)
	}
}

func TestCompositeRegistry_ExecuteDispatchesBySourceID(t *testing.T) {
	tlaloc := &sourceRegistryStub{label: "tlaloc", candidates: map[string][]CapabilityCandidate{}}
	parrot := &sourceRegistryStub{label: "parrot", candidates: map[string][]CapabilityCandidate{}}
	registry, _ := NewCompositeRegistry(
		RegistrySource{ID: "tlaloc", Registry: tlaloc},
		RegistrySource{ID: "parrot", Registry: parrot},
	)

	_, err := registry.Execute(context.Background(), CapabilityExecutionRequest{
		TaskID: "run", NodeID: "node", Capability: "X", SourceID: "parrot", WorkerID: "shared",
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if parrot.calls != 1 || tlaloc.calls != 0 {
		t.Fatalf("dispatch calls tlaloc=%d parrot=%d", tlaloc.calls, parrot.calls)
	}
	if parrot.last.SourceID != "" {
		t.Fatalf("leaf registry should receive source-local request, SourceID=%q", parrot.last.SourceID)
	}
	if parrot.last.WorkerID != "shared" {
		t.Fatalf("leaf WorkerID = %q", parrot.last.WorkerID)
	}
}

func TestCompositeRegistry_MultipleSourcesRequireSourceID(t *testing.T) {
	a := &sourceRegistryStub{candidates: map[string][]CapabilityCandidate{}}
	b := &sourceRegistryStub{candidates: map[string][]CapabilityCandidate{}}
	registry, _ := NewCompositeRegistry(RegistrySource{ID: "a", Registry: a}, RegistrySource{ID: "b", Registry: b})
	_, err := registry.Execute(context.Background(), CapabilityExecutionRequest{WorkerID: "same"})
	if err == nil || !strings.Contains(err.Error(), "SourceID is required") {
		t.Fatalf("expected SourceID requirement, got %v", err)
	}
}

func TestCompositeRegistry_SingleSourceKeepsLegacyExecutionCompatible(t *testing.T) {
	leaf := &sourceRegistryStub{label: "one", candidates: map[string][]CapabilityCandidate{}}
	registry, _ := NewCompositeRegistry(RegistrySource{ID: "only", Registry: leaf})
	_, err := registry.Execute(context.Background(), CapabilityExecutionRequest{WorkerID: "worker"})
	if err != nil {
		t.Fatalf("single-source legacy execution: %v", err)
	}
	if leaf.calls != 1 {
		t.Fatalf("leaf calls = %d", leaf.calls)
	}
}

func TestControlLoop_WithDuplicateWorkerIDsExecutesMachinerySourceNotParrot(t *testing.T) {
	tlaloc := &sourceRegistryStub{label: "tlaloc", candidates: map[string][]CapabilityCandidate{
		"NORMALIZE": {{WorkerID: "shared", Capability: "NORMALIZE", Kind: CapabilityTlaloque, Deterministic: true}},
	}}
	parrot := &sourceRegistryStub{label: "parrot", candidates: map[string][]CapabilityCandidate{
		"NORMALIZE": {{WorkerID: "shared", Capability: "NORMALIZE", Kind: CapabilityExternalModel, Generative: true}},
	}}
	registry, _ := NewCompositeRegistry(
		RegistrySource{ID: "parrot", Registry: parrot},
		RegistrySource{ID: "tlaloc", Registry: tlaloc},
	)
	controller := FixedProgramController{Decisions: []ControlDecision{{
		Action: ControlExecute, Capability: "NORMALIZE", Input: json.RawMessage(`{"raw":"7"}`), AllowExternalModel: true,
	}}}
	loop := &ControlLoop{Registry: registry, Selection: MachineryFirstPolicy{}, MaxTransitions: 2}
	run, err := loop.Run(context.Background(), ControlRunRequest{RunID: "r", Goal: "normalize"}, controller)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if len(run.Transitions) == 0 || run.Transitions[0].Selected == nil {
		t.Fatalf("missing selected transition: %+v", run)
	}
	selected := run.Transitions[0].Selected
	if selected.SourceID != "tlaloc" || selected.Kind != CapabilityTlaloque {
		t.Fatalf("selected %+v, want tlaloc machinery", selected)
	}
	if tlaloc.calls != 1 || parrot.calls != 0 {
		t.Fatalf("execution calls tlaloc=%d parrot=%d", tlaloc.calls, parrot.calls)
	}
}

func TestCompositeRegistryRejectsDuplicateSourceIDs(t *testing.T) {
	leaf := &sourceRegistryStub{candidates: map[string][]CapabilityCandidate{}}
	_, err := NewCompositeRegistry(RegistrySource{ID: "same", Registry: leaf}, RegistrySource{ID: "same", Registry: leaf})
	if err == nil || !strings.Contains(err.Error(), "duplicate source ID") {
		t.Fatalf("expected duplicate source error, got %v", err)
	}
}
