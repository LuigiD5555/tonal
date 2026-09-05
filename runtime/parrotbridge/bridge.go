// Package parrotbridge adapts an external probabilistic model client into a
// Tonal CapabilityRegistry source. It deliberately contains no Tlaloc
// dependency: Parrot is external cognition, not Tlaloc machinery.
package parrotbridge

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"tonal.local/runtime/tonal"
)

// Client is the provider-neutral Parrot invocation seam. Claude, an
// OpenAI-compatible local endpoint, or another backend may implement this
// interface without changing Tonal's component ontology.
type Client interface {
	Invoke(ctx context.Context, req Request) (Result, error)
}

// Request is one bounded external-cognition operation. Input is the contract
// payload constructed by Tonal; PriorObservations are read-only context.
type Request struct {
	Capability        string              `json:"capability"`
	Input             json.RawMessage     `json:"input"`
	PriorObservations []tonal.Observation `json:"prior_observations,omitempty"`
}

// Result is provider-neutral. Internal Parrot capabilities may return
// OBSERVATION records that still require Tonal verification. Presentation-only
// capabilities such as RENDER_RESPONSE should return no observations.
type Result struct {
	Output       json.RawMessage       `json:"output"`
	Confidence   float64               `json:"confidence,omitempty"`
	Notes        string                `json:"notes,omitempty"`
	Observations []tonal.Observation   `json:"observations,omitempty"`
	Usage        *tonal.CapabilityUsage `json:"usage,omitempty"`
}

// CapabilitySpec declares one bounded behavior exposed through Parrot. It is
// a runtime contract, not a claim that the external model reliably possesses
// the capability; reliability belongs in evidence/profile data.
type CapabilitySpec struct {
	Capability  string `json:"capability"`
	WorkerID    string `json:"worker_id,omitempty"`
	ProfileRef  string `json:"profile_ref,omitempty"`
	EvidenceRef string `json:"evidence_ref,omitempty"`
}

// Registry exposes a configured external model client as EXTERNAL_MODEL
// candidates only. CompositeRegistry is responsible for stamping SourceID
// (normally "parrot").
type Registry struct {
	client Client
	specs  map[string]CapabilitySpec
}

var _ tonal.CapabilityRegistry = (*Registry)(nil)

func New(client Client, specs ...CapabilitySpec) (*Registry, error) {
	if client == nil {
		return nil, fmt.Errorf("parrotbridge: Client is required")
	}
	if len(specs) == 0 {
		return nil, fmt.Errorf("parrotbridge: at least one capability spec is required")
	}
	registry := &Registry{client: client, specs: make(map[string]CapabilitySpec, len(specs))}
	for _, raw := range specs {
		capability := strings.ToUpper(strings.TrimSpace(raw.Capability))
		if capability == "" {
			return nil, fmt.Errorf("parrotbridge: capability cannot be empty")
		}
		if _, exists := registry.specs[capability]; exists {
			return nil, fmt.Errorf("parrotbridge: duplicate capability %q", capability)
		}
		raw.Capability = capability
		raw.WorkerID = strings.TrimSpace(raw.WorkerID)
		if raw.WorkerID == "" {
			raw.WorkerID = "default"
		}
		registry.specs[capability] = raw
	}
	return registry, nil
}

// ResponseRendererSpec is the canonical final-response Parrot surface. The
// semantic restrictions themselves are enforced by tonal.ResponseComposer and
// ResponseVerifier; this spec merely publishes the bounded rendering ability.
func ResponseRendererSpec(profileRef string) CapabilitySpec {
	return CapabilitySpec{
		Capability: tonal.ResponseRenderCapability,
		WorkerID:   "default",
		ProfileRef: strings.TrimSpace(profileRef),
	}
}

func (r *Registry) Capabilities() []CapabilitySpec {
	if r == nil {
		return nil
	}
	keys := make([]string, 0, len(r.specs))
	for capability := range r.specs {
		keys = append(keys, capability)
	}
	sort.Strings(keys)
	out := make([]CapabilitySpec, 0, len(keys))
	for _, capability := range keys {
		out = append(out, r.specs[capability])
	}
	return out
}

func (r *Registry) Candidates(capability string, _ tonal.CapabilityGoal) []tonal.CapabilityCandidate {
	if r == nil {
		return nil
	}
	capability = strings.ToUpper(strings.TrimSpace(capability))
	spec, ok := r.specs[capability]
	if !ok {
		return nil
	}
	return []tonal.CapabilityCandidate{{
		WorkerID:      spec.WorkerID,
		Capability:    capability,
		Kind:          tonal.CapabilityExternalModel,
		EngineKind:    "external_model",
		ProfileRef:    spec.ProfileRef,
		Deterministic: false,
		Generative:    true,
		Selected:      true,
		Reason:        "configured Parrot external-cognition capability",
	}}
}

func (r *Registry) Execute(ctx context.Context, req tonal.CapabilityExecutionRequest) (tonal.CapabilityExecutionResult, error) {
	if r == nil || r.client == nil {
		return tonal.CapabilityExecutionResult{}, fmt.Errorf("parrotbridge: registry is not configured")
	}
	capability := strings.ToUpper(strings.TrimSpace(req.Capability))
	spec, ok := r.specs[capability]
	if !ok {
		return tonal.CapabilityExecutionResult{}, fmt.Errorf("parrotbridge: capability %q is not configured", capability)
	}
	if req.WorkerID != "" && req.WorkerID != spec.WorkerID {
		return tonal.CapabilityExecutionResult{}, fmt.Errorf("parrotbridge: worker %q does not own capability %q", req.WorkerID, capability)
	}

	result, err := r.client.Invoke(ctx, Request{
		Capability:        capability,
		Input:             append(json.RawMessage(nil), req.Input...),
		PriorObservations: append([]tonal.Observation(nil), req.PriorObservations...),
	})
	if err != nil {
		return tonal.CapabilityExecutionResult{}, err
	}
	if len(result.Output) == 0 || !json.Valid(result.Output) {
		return tonal.CapabilityExecutionResult{}, fmt.Errorf("parrotbridge: client returned invalid JSON output for %s", capability)
	}
	return tonal.CapabilityExecutionResult{
		WorkerID:     spec.WorkerID,
		Output:       append(json.RawMessage(nil), result.Output...),
		Confidence:   result.Confidence,
		Notes:        result.Notes,
		Observations: append([]tonal.Observation(nil), result.Observations...),
		Usage:        result.Usage,
	}, nil
}
