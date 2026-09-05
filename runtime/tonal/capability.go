package tonal

import "tlaloc.local/behaviorlab/tlaloquekit"

// CapabilityCandidate is TONAL's runtime view of one qualified executor
// candidate for a bounded capability. It deliberately contains only the
// fields selection policies need; policies must not depend on Tlaloc's
// registry implementation details.
//
// Parrot is not represented by a special type. A Parrot-backed executor is
// simply a candidate whose published role is generative.
type CapabilityCandidate struct {
	WorkerID      string `json:"worker_id"`
	Capability    string `json:"capability"`
	EngineKind    string `json:"engine_kind"`
	ProfileRef    string `json:"profile_ref,omitempty"`
	Deterministic bool   `json:"deterministic"`
	Generative    bool   `json:"generative"`
	Selected      bool   `json:"selected"`
	Reason        string `json:"reason,omitempty"`
}

// capabilityCandidates adapts Tlaloc's current QualifiedRegistry response to
// TONAL-owned selection data. Keeping this translation at the runtime edge
// lets future registries provide the same selection surface without making
// SelectionPolicy depend on tlaloquekit.Candidate.
func capabilityCandidates(candidates []tlaloquekit.Candidate) []CapabilityCandidate {
	out := make([]CapabilityCandidate, 0, len(candidates))
	for _, candidate := range candidates {
		out = append(out, CapabilityCandidate{
			WorkerID:      candidate.Descriptor.ID,
			Capability:    candidate.Descriptor.Capability,
			EngineKind:    string(candidate.Descriptor.Engine),
			ProfileRef:    candidate.Descriptor.ProfileRef,
			Deterministic: candidate.Descriptor.Deterministic,
			Generative:    candidate.Descriptor.Engine == tlaloquekit.EngineGenerative,
			Selected:      candidate.Selected,
			Reason:        candidate.Reason,
		})
	}
	return out
}
