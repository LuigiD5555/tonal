package tonal

import (
	"testing"

	"tlaloc.local/behaviorlab/tlaloquekit"
)

var _ SelectionPolicy = HeterogeneousPolicy{}
var _ RoutingPolicy = HeterogeneousPolicy{}

func TestCapabilityCandidatesAdaptRegistryWithoutParrotSpecialCase(t *testing.T) {
	in := []tlaloquekit.Candidate{
		{
			Descriptor: tlaloquekit.Descriptor{
				ID:            "det-normalize",
				Capability:    "NORMALIZE",
				Engine:        tlaloquekit.EngineDeterministic,
				Deterministic: true,
				ProfileRef:    "normalize@r1",
			},
			Selected: true,
			Reason:   "qualified deterministic candidate",
		},
		{
			Descriptor: tlaloquekit.Descriptor{
				ID:         "generic-model-worker",
				Capability: "NORMALIZE",
				Engine:     tlaloquekit.EngineGenerative,
				ProfileRef: "model@r1",
			},
			Reason: "qualified generative candidate",
		},
	}

	got := capabilityCandidates(in)
	if len(got) != 2 {
		t.Fatalf("got %d candidates, want 2", len(got))
	}
	if !got[0].Deterministic || got[0].Generative {
		t.Fatalf("deterministic candidate flags = deterministic:%v generative:%v", got[0].Deterministic, got[0].Generative)
	}
	if !got[1].Generative || got[1].Deterministic {
		t.Fatalf("generative candidate flags = deterministic:%v generative:%v", got[1].Deterministic, got[1].Generative)
	}
	if got[1].WorkerID != "generic-model-worker" {
		t.Fatalf("worker id = %q", got[1].WorkerID)
	}
}

func TestParrotCentricPolicySelectsGenerativeRoleNotParrotName(t *testing.T) {
	policy := ParrotCentricPolicy{CognitiveCapabilities: map[string]bool{"NORMALIZE": true}}
	candidates := []CapabilityCandidate{
		{WorkerID: "det-normalize", Capability: "NORMALIZE", Deterministic: true, Selected: true},
		{WorkerID: "replaceable-model-x", Capability: "NORMALIZE", Generative: true},
	}

	workerID, _ := policy.SelectWorker(Step{Capability: "NORMALIZE"}, candidates)
	if workerID != "replaceable-model-x" {
		t.Fatalf("selected %q, want generic generative candidate", workerID)
	}
}

func TestHeterogeneousPolicyDefersToRegistrySelection(t *testing.T) {
	workerID, reason := (HeterogeneousPolicy{}).SelectWorker(
		Step{Capability: "ARITHMETIC"},
		[]CapabilityCandidate{{WorkerID: "det-arithmetic", Selected: true, Deterministic: true}},
	)
	if workerID != "" {
		t.Fatalf("heterogeneous policy overrode registry with %q", workerID)
	}
	if reason == "" {
		t.Fatal("heterogeneous policy should record why it deferred")
	}
}
