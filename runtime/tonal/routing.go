package tonal

import "strings"

// SelectionPolicy decides which already-qualified capability candidate Tonal
// should use for one bounded Step. It does not construct executors, mutate the
// Registry, or receive provider/model-specific request parameters.
type SelectionPolicy interface {
	Name() string
	SelectWorker(step Step, candidates []CapabilityCandidate) (workerID string, reason string)
}

// RoutingPolicy is the T1-era name retained as a source-compatible alias.
type RoutingPolicy = SelectionPolicy

// MachineryFirstPolicy is the initial post-T1 R2 policy. Given an already
// eligible candidate set it prefers deterministic machinery, then other
// non-external machinery, and only then external cognition.
//
// Eligibility is a separate concern. In particular, Control Loop R0 filters
// EXTERNAL_MODEL candidates unless the controller explicitly permits them for
// that transition. This policy therefore cannot silently wake Parrot.
type MachineryFirstPolicy struct{}

func (MachineryFirstPolicy) Name() string { return "R2_MACHINERY_FIRST" }

func (MachineryFirstPolicy) SelectWorker(_ Step, candidates []CapabilityCandidate) (string, string) {
	for _, candidate := range candidates {
		if candidate.Kind != CapabilityExternalModel && candidate.Deterministic {
			return candidate.WorkerID, "machinery-first: deterministic non-external capability"
		}
	}
	for _, candidate := range candidates {
		if candidate.Kind != CapabilityExternalModel {
			return candidate.WorkerID, "machinery-first: qualified non-external capability"
		}
	}
	for _, candidate := range candidates {
		if candidate.Kind == CapabilityExternalModel {
			return candidate.WorkerID, "machinery-first: external cognition is the remaining eligible capability"
		}
	}
	return "", "no eligible candidate"
}

// HeterogeneousPolicy is frozen T1 Arm C: Tonal never overrides the frozen
// Tlaloc R1 Registry's ranking.
type HeterogeneousPolicy struct{}

func (HeterogeneousPolicy) Name() string { return "C_HETEROGENEOUS_TONAL" }

func (HeterogeneousPolicy) SelectWorker(Step, []CapabilityCandidate) (string, string) {
	return "", "registry deterministic-first selection (no arm override)"
}

// ParrotCentricPolicy is the frozen Arm B compatibility policy. For the
// cognitive capabilities declared by T1 it forces the EXTERNAL_MODEL
// candidate instead of reusable machinery. Infrastructure capabilities are
// left to the frozen Registry.
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
