package tonal

import "testing"

var _ SelectionPolicy = HeterogeneousPolicy{}
var _ RoutingPolicy = HeterogeneousPolicy{}

func TestComponentKindIsIndependentFromGenerativeBehavior(t *testing.T) {
	boundedGenerative := CapabilityCandidate{
		WorkerID:   "bounded-generative-specialist",
		Capability: "NORMALIZE",
		Kind:       CapabilityTlaloque,
		Generative: true,
	}
	external := CapabilityCandidate{
		WorkerID:   "external-model",
		Capability: "NORMALIZE",
		Kind:       CapabilityExternalModel,
		Generative: true,
	}
	if boundedGenerative.Kind == external.Kind {
		t.Fatal("generative behavior must not collapse TLALOQUE and EXTERNAL_MODEL kinds")
	}
}

func TestParrotCentricPolicySelectsExternalModelKindNotGenerativeFlag(t *testing.T) {
	policy := ParrotCentricPolicy{CognitiveCapabilities: map[string]bool{"NORMALIZE": true}}
	candidates := []CapabilityCandidate{
		{WorkerID: "det-normalize", Capability: "NORMALIZE", Kind: CapabilityTlaloque, Deterministic: true, Selected: true},
		{WorkerID: "bounded-generative-specialist", Capability: "NORMALIZE", Kind: CapabilityTlaloque, Generative: true},
		{WorkerID: "replaceable-external-model", Capability: "NORMALIZE", Kind: CapabilityExternalModel, Generative: true},
	}

	workerID, _ := policy.SelectWorker(Step{Capability: "NORMALIZE"}, candidates)
	if workerID != "replaceable-external-model" {
		t.Fatalf("selected %q, want EXTERNAL_MODEL candidate", workerID)
	}
}

func TestHeterogeneousPolicyDefersToRegistrySelection(t *testing.T) {
	workerID, reason := (HeterogeneousPolicy{}).SelectWorker(
		Step{Capability: "ARITHMETIC"},
		[]CapabilityCandidate{{WorkerID: "det-arithmetic", Kind: CapabilityTlaloque, Selected: true, Deterministic: true}},
	)
	if workerID != "" {
		t.Fatalf("heterogeneous policy overrode registry with %q", workerID)
	}
	if reason == "" {
		t.Fatal("heterogeneous policy should record why it deferred")
	}
}
