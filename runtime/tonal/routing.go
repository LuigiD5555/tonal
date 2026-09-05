package tonal

import "strings"

// SelectionPolicy decides which already-qualified capability candidate TONAL
// should use for one bounded Step. It does not construct executors, mutate the
// Registry, or receive model-specific request parameters.
//
// A policy may return "" to accept the Registry's preselected candidate. This
// keeps the current T1 behavior while making selection independently
// replaceable for later R2 experiments.
type SelectionPolicy interface {
	Name() string
	SelectWorker(step Step, candidates []CapabilityCandidate) (workerID string, reason string)
}

// RoutingPolicy is the T1-era name retained as a source-compatible alias.
// New R2 code should prefer SelectionPolicy. Frozen T1 behavior is unchanged.
type RoutingPolicy = SelectionPolicy

// HeterogeneousPolicy is Arm C: TONAL never overrides the Registry. The
// Registry's deterministic-first, smallest-first, CapabilityProfile-vetoed
// ranking decides every executor.
type HeterogeneousPolicy struct{}

func (HeterogeneousPolicy) Name() string { return "C_HETEROGENEOUS_TONAL" }

func (HeterogeneousPolicy) SelectWorker(Step, []CapabilityCandidate) (string, string) {
	return "", "registry deterministic-first selection (no arm override)"
}

// ParrotCentricPolicy is Arm B: for capabilities that could plausibly be
// asked of a generative executor, force a generative candidate instead of the
// deterministic one. Infrastructure capabilities (locate, crop, verify,
// store) are left to the Registry. This measures external sequence + memory
// without heterogeneous specialisation.
//
// The policy does not know a Parrot worker id or model name. It selects by the
// generic generative role published through CapabilityCandidate, so Parrot
// remains one replaceable capability implementation rather than a privileged
// runtime path.
type ParrotCentricPolicy struct {
	// CognitiveCapabilities is the set forced toward a generative executor,
	// e.g. NORMALIZE, COMPARE_NUMBERS, ARITHMETIC, SAME_DIFFERENT,
	// EXTRACT_NUMBER.
	CognitiveCapabilities map[string]bool
}

func (ParrotCentricPolicy) Name() string { return "B_PARROT_CENTRIC_COMPOSITION" }

func (p ParrotCentricPolicy) SelectWorker(step Step, candidates []CapabilityCandidate) (string, string) {
	if !p.CognitiveCapabilities[strings.ToUpper(step.Capability)] {
		return "", "infrastructure capability left to the registry"
	}
	for _, candidate := range candidates {
		if candidate.Generative {
			return candidate.WorkerID, "arm B forces a generative executor for a cognitive capability"
		}
	}
	return "", "no generative candidate available; falling back to the registry"
}
