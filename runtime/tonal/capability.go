package tonal

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

// CapabilityCandidate is Tonal's runtime view of one eligible executor
// candidate for a bounded capability. SourceID identifies the registry source
// that owns execution dispatch; WorkerID only needs to be unique within that
// source.
type CapabilityCandidate struct {
	SourceID      string         `json:"source_id,omitempty"`
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
