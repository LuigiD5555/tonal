package tonal

import "tlaloc.local/behaviorlab/tlaloquekit"

// CapabilityKind identifies what class of component supplies a capability.
// It is deliberately separate from behavioral properties such as
// Deterministic or Generative.
type CapabilityKind string

const (
	CapabilityTlaloque      CapabilityKind = "TLALOQUE"
	CapabilityMachine       CapabilityKind = "MACHINE"
	CapabilityTool          CapabilityKind = "TOOL"
	CapabilityExternalModel CapabilityKind = "EXTERNAL_MODEL"
)

// CapabilityCandidate is TONAL's runtime view of one qualified executor
// candidate for a bounded capability. Selection policies depend on this
// Tonal-owned type rather than Tlaloc registry implementation details.
type CapabilityCandidate struct {
	WorkerID      string         `json:"worker_id"`
	Capability    string         `json:"capability"`
	Kind          CapabilityKind `json:"kind"`
	EngineKind    string         `json:"engine_kind"`
	ProfileRef    string         `json:"profile_ref,omitempty"`
	Deterministic bool           `json:"deterministic"`
	Generative    bool           `json:"generative"`
	Selected      bool           `json:"selected"`
	Reason        string         `json:"reason,omitempty"`
}

// capabilityCandidates is the compatibility adapter for the frozen Tlaloc R1
// publication contract used by T1. That contract historically described
// Parrot as EngineGenerative / a generative Tlaloque. Architecture R2 does
// not rewrite the frozen contract; it reclassifies its generative Parrot at
// the Tonal boundary as EXTERNAL_MODEL.
//
// The frozen R1 registry has only one generative family: Parrot. Future
// adapters MUST publish component Kind explicitly rather than assuming every
// generative executor is external cognition.
func capabilityCandidates(candidates []tlaloquekit.Candidate) []CapabilityCandidate {
	out := make([]CapabilityCandidate, 0, len(candidates))
	for _, candidate := range candidates {
		kind := CapabilityTlaloque
		if candidate.Descriptor.Engine == tlaloquekit.EngineGenerative {
			kind = CapabilityExternalModel
		}
		out = append(out, CapabilityCandidate{
			WorkerID:      candidate.Descriptor.ID,
			Capability:    candidate.Descriptor.Capability,
			Kind:          kind,
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
