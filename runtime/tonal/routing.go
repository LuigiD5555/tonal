package tonal

import (
	"strings"

	"tlaloc.local/behaviorlab/tlaloquekit"
)

// RoutingPolicy is how an experiment arm influences executor selection
// WITHOUT TONAL learning any executor specifics. A policy may only pick a
// worker id that already appears in the Registry's candidate list for that
// capability; it never constructs executor-specific parameters.
type RoutingPolicy interface {
	Name() string
	// SelectWorker returns a worker id to pin for this step, or "" to
	// accept the Registry's own deterministic-first selection. The reason
	// is recorded in the routing trace.
	SelectWorker(step Step, candidates []tlaloquekit.Candidate) (workerID string, reason string)
}

// HeterogeneousPolicy is Arm C: TONAL never overrides the Registry. The
// Registry's deterministic-first, smallest-first, CapabilityProfile-vetoed
// ranking decides every executor.
type HeterogeneousPolicy struct{}

func (HeterogeneousPolicy) Name() string { return "C_HETEROGENEOUS_TONAL" }

func (HeterogeneousPolicy) SelectWorker(Step, []tlaloquekit.Candidate) (string, string) {
	return "", "registry deterministic-first selection (no arm override)"
}

// ParrotCentricPolicy is Arm B: for capabilities that could plausibly be
// asked of a generative executor, force the generative candidate instead
// of the deterministic one. Infrastructure capabilities (locate, crop,
// verify, store) are left to the Registry. This measures external sequence
// + memory without heterogeneous specialisation.
//
// The policy identifies the generative candidate purely by its published
// EngineKind — it holds no model-specific knowledge.
type ParrotCentricPolicy struct {
	// CognitiveCapabilities is the set forced toward the generative
	// executor, e.g. NORMALIZE, COMPARE_NUMBERS, ARITHMETIC,
	// SAME_DIFFERENT, EXTRACT_NUMBER.
	CognitiveCapabilities map[string]bool
}

func (ParrotCentricPolicy) Name() string { return "B_PARROT_CENTRIC_COMPOSITION" }

func (p ParrotCentricPolicy) SelectWorker(step Step, candidates []tlaloquekit.Candidate) (string, string) {
	if !p.CognitiveCapabilities[strings.ToUpper(step.Capability)] {
		return "", "infrastructure capability left to the registry"
	}
	for _, candidate := range candidates {
		if candidate.Descriptor.Engine == tlaloquekit.EngineGenerative {
			return candidate.Descriptor.ID, "arm B forces the generative executor for a cognitive capability"
		}
	}
	return "", "no generative candidate available; falling back to the registry"
}
