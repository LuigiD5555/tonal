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

// ParrotCentricPolicy is the frozen Arm B compatibility policy. For the
// cognitive capabilities declared by T1 it forces the EXTERNAL_MODEL
// candidate instead of reusable machinery. Infrastructure capabilities are
// left to the Registry.
//
// The policy knows no provider, model name, worker id or prompt. It selects
// only by Tonal's component Kind. This preserves T1's Parrot-centric arm while
// keeping Parrot external to the Tlaloque machinery ontology.
type ParrotCentricPolicy struct {
	CognitiveCapabilities map[string]bool
}

func (ParrotCentricPolicy) Name() string { return "B_PARROT_CENTRIC_COMPOSITION" }

func (p ParrotCentricPolicy) SelectWorker(step Step, candidates []CapabilityCandidate) (string, string) {
	if !p.CognitiveCapabilities[strings.ToUpper(step.Capability)] {
		return "", "infrastructure capability left to the registry"
	}
	for _, candidate := range candidates {
		if candidate.Kind == CapabilityExternalModel {
			return candidate.WorkerID, "arm B forces external probabilistic cognition for a cognitive capability"
		}
	}
	return "", "no external-model candidate available; falling back to the registry"
}
