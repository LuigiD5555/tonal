package tonal

import (
	"context"
	"encoding/json"
)

// CapabilityDescriptor is Tonal's registry-level description of one
// executable component. Capability says what behavior it offers; Kind says
// what class of component supplies it.
type CapabilityDescriptor struct {
	ID             string         `json:"id"`
	Capability     string         `json:"capability"`
	Kind           CapabilityKind `json:"kind"`
	EngineKind     string         `json:"engine_kind,omitempty"`
	Deterministic  bool           `json:"deterministic"`
	ParameterCount int64          `json:"parameter_count,omitempty"`
	InputSchema    string         `json:"input_schema,omitempty"`
	OutputSchema   string         `json:"output_schema,omitempty"`
	Dependencies   []string       `json:"dependencies,omitempty"`
	ProfileRef     string         `json:"profile_ref,omitempty"`
	EvidenceRef    string         `json:"evidence_ref,omitempty"`
}

// CapabilityGoal asks for behavior, never for a specific implementation.
type CapabilityGoal struct {
	Capability          string   `json:"capability"`
	PreferDeterministic bool     `json:"prefer_deterministic,omitempty"`
	MaxParameters       int64    `json:"max_parameters,omitempty"`
	AvailableProducts   []string `json:"available_products,omitempty"`
}

// Observation is Tonal's typed state/evidence record. Executors return
// observations; only the Tonal runtime decides what is committed to its
// Blackboard and what may be promoted to FACT under verification rules.
type Observation struct {
	Producer       string            `json:"producer"`
	Capability     string            `json:"capability"`
	Key            string            `json:"key"`
	Value          json.RawMessage   `json:"value"`
	Kind           string            `json:"kind"` // OBSERVATION | FACT
	Status         string            `json:"status,omitempty"`
	Confidence     float64           `json:"confidence,omitempty"`
	References     []string          `json:"references,omitempty"`
	Provenance     map[string]string `json:"provenance,omitempty"`
	ProfileVersion string            `json:"profile_version,omitempty"`
	RecordedAt     string            `json:"recorded_at,omitempty"`
}

// CapabilityUsage is optional execution accounting. External cognition and
// model-backed machinery may report model calls/tokens; deterministic
// machinery normally leaves Usage nil.
type CapabilityUsage struct {
	PromptTokens     int `json:"prompt_tokens,omitempty"`
	CompletionTokens int `json:"completion_tokens,omitempty"`
	ModelCalls       int `json:"model_calls,omitempty"`
}

// CapabilityExecutionRequest asks the Registry to execute one already
// selected component on one bounded workflow node.
type CapabilityExecutionRequest struct {
	TaskID            string          `json:"task_id"`
	NodeID            string          `json:"node_id"`
	Capability        string          `json:"capability"`
	WorkerID          string          `json:"worker_id"`
	Input             json.RawMessage `json:"input"`
	PriorObservations []Observation   `json:"prior_observations,omitempty"`
}

// CapabilityExecutionResult is the component-neutral result returned to the
// Tonal Engine.
type CapabilityExecutionResult struct {
	WorkerID     string          `json:"worker_id"`
	Output       json.RawMessage `json:"output"`
	Confidence   float64         `json:"confidence,omitempty"`
	Notes        string          `json:"notes,omitempty"`
	Observations []Observation   `json:"observations,omitempty"`
	Usage        *CapabilityUsage `json:"usage,omitempty"`
}

// CapabilityRegistry is the runtime seam Tonal owns. Tlaloc, a local Machine
// registry, tool adapters and Parrot/external-model adapters may all feed this
// surface without becoming Tonal's scheduler or Blackboard authority.
type CapabilityRegistry interface {
	Candidates(capability string, goal CapabilityGoal) []CapabilityCandidate
	Execute(ctx context.Context, req CapabilityExecutionRequest) (CapabilityExecutionResult, error)
}
