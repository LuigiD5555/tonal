package tonal

import "strings"

// SelectionPolicy decides which already-qualified capability candidate Tonal
// should use for one bounded Step. It does not construct executors, mutate the
// Registry, or receive provider/model-specific request parameters.
type SelectionPolicy interface {
	Name() string
	SelectWorker(step Step, candidates []CapabilityCandidate) (workerID string, reason string)
}

// SourceAwareSelectionPolicy is the R2 extension for multi-source registries.
// WorkerID is only required to be unique within SourceID, so policies that
// actively choose a candidate should return the complete candidate identity.
// Legacy/frozen policies remain valid through SelectionPolicy.
type SourceAwareSelectionPolicy interface {
	SelectCandidate(step Step, candidates []CapabilityCandidate) (CapabilityCandidate, string, bool)
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

func (p MachineryFirstPolicy) SelectWorker(step Step, candidates []CapabilityCandidate) (string, string) {
	candidate, reason, ok := p.SelectCandidate(step, candidates)
	if !ok {
		return "", reason
	}
	return candidate.WorkerID, reason
}

func (MachineryFirstPolicy) SelectCandidate(_ Step, candidates []CapabilityCandidate) (CapabilityCandidate, string, bool) {
	for _, candidate := range candidates {
		if candidate.Kind != CapabilityExternalModel && candidate.Deterministic {
			return candidate, "machinery-first: deterministic non-external capability", true
		}
	}
	for _, candidate := range candidates {
		if candidate.Kind != CapabilityExternalModel {
			return candidate, "machinery-first: qualified non-external capability", true
		}
	}
	for _, candidate := range candidates {
		if candidate.Kind == CapabilityExternalModel {
			return candidate, "machinery-first: external cognition is the remaining eligible capability", true
		}
	}
	return CapabilityCandidate{}, "no eligible candidate", false
}

// HeterogeneousPolicy is frozen T1 Arm C: Tonal never overrides the frozen
// Tlaloc R1 Registry's ranking. It deliberately remains a legacy
// SelectionPolicy because T1 has one source and must preserve frozen behavior.
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
	candidate, reason, ok := p.SelectCandidate(step, candidates)
	if !ok {
		return "", reason
	}
	return candidate.WorkerID, reason
}

func (p ParrotCentricPolicy) SelectCandidate(step Step, candidates []CapabilityCandidate) (CapabilityCandidate, string, bool) {
	if !p.CognitiveCapabilities[strings.ToUpper(step.Capability)] {
		return CapabilityCandidate{}, "infrastructure capability left to the registry", false
	}
	for _, candidate := range candidates {
		if candidate.Kind == CapabilityExternalModel {
			return candidate, "arm B forces external probabilistic cognition for a cognitive capability", true
		}
	}
	return CapabilityCandidate{}, "no external-model candidate available; falling back to the registry", false
}

// selectCapabilityCandidate is Tonal's common selection resolver. R2-aware
// policies return SourceID+WorkerID directly. Legacy policies are accepted
// only when their worker identity is unambiguous across the candidate set.
func selectCapabilityCandidate(policy SelectionPolicy, step Step, candidates []CapabilityCandidate) (CapabilityCandidate, string, bool) {
	if aware, ok := policy.(SourceAwareSelectionPolicy); ok {
		candidate, reason, selected := aware.SelectCandidate(step, candidates)
		if selected {
			for _, eligible := range candidates {
				if eligible.SourceID == candidate.SourceID && eligible.WorkerID == candidate.WorkerID {
					return eligible, reason, true
				}
			}
			return CapabilityCandidate{}, reason, false
		}
		// A source-aware policy may deliberately defer (for example an
		// infrastructure capability in the frozen Arm B compatibility path).
		if reason != "" && strings.Contains(reason, "left to the registry") {
			return selectedRegistryCandidate(candidates, reason)
		}
	}

	workerID, reason := policy.SelectWorker(step, candidates)
	if workerID != "" {
		matches := make([]CapabilityCandidate, 0, 1)
		for _, candidate := range candidates {
			if candidate.WorkerID == workerID {
				matches = append(matches, candidate)
			}
		}
		if len(matches) == 1 {
			return matches[0], reason, true
		}
		if len(matches) > 1 {
			var selected *CapabilityCandidate
			for _, candidate := range matches {
				if !candidate.Selected {
					continue
				}
				if selected != nil {
					return CapabilityCandidate{}, "ambiguous legacy selection: WorkerID exists in multiple sources", false
				}
				copyCandidate := candidate
				selected = &copyCandidate
			}
			if selected != nil {
				return *selected, reason, true
			}
			return CapabilityCandidate{}, "ambiguous legacy selection: WorkerID exists in multiple sources", false
		}
		return CapabilityCandidate{}, reason, false
	}

	return selectedRegistryCandidate(candidates, reason)
}

func selectedRegistryCandidate(candidates []CapabilityCandidate, reason string) (CapabilityCandidate, string, bool) {
	var selected *CapabilityCandidate
	for _, candidate := range candidates {
		if !candidate.Selected {
			continue
		}
		if selected != nil {
			return CapabilityCandidate{}, "registry selection is ambiguous across multiple sources", false
		}
		copyCandidate := candidate
		selected = &copyCandidate
	}
	if selected != nil {
		if reason == "" {
			reason = selected.Reason
		}
		return *selected, reason, true
	}
	if len(candidates) == 0 {
		return CapabilityCandidate{}, reason, false
	}
	if reason == "" {
		reason = "fallback to first eligible candidate"
	}
	return candidates[0], reason, true
}
