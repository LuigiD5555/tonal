package tonal

import (
	"context"
	"fmt"
	"strings"
)

// RegistrySource names one independent provider of Tonal capabilities.
// Examples are "tlaloc", "machines", "tools" and "parrot". WorkerID only
// needs to be unique within this source.
type RegistrySource struct {
	ID       string
	Registry CapabilityRegistry
}

// CompositeRegistry merges independent capability sources without collapsing
// their identities. Candidate SourceID is stamped at this boundary and is
// carried unchanged through selection, trace and execution dispatch.
type CompositeRegistry struct {
	sources []RegistrySource
	byID    map[string]CapabilityRegistry
}

func NewCompositeRegistry(sources ...RegistrySource) (*CompositeRegistry, error) {
	if len(sources) == 0 {
		return nil, fmt.Errorf("composite registry: at least one source is required")
	}
	composite := &CompositeRegistry{
		sources: make([]RegistrySource, 0, len(sources)),
		byID:    make(map[string]CapabilityRegistry, len(sources)),
	}
	for _, raw := range sources {
		id := strings.TrimSpace(raw.ID)
		if id == "" {
			return nil, fmt.Errorf("composite registry: source ID cannot be empty")
		}
		if raw.Registry == nil {
			return nil, fmt.Errorf("composite registry: source %q has nil registry", id)
		}
		if _, exists := composite.byID[id]; exists {
			return nil, fmt.Errorf("composite registry: duplicate source ID %q", id)
		}
		source := RegistrySource{ID: id, Registry: raw.Registry}
		composite.sources = append(composite.sources, source)
		composite.byID[id] = raw.Registry
	}
	return composite, nil
}

func (c *CompositeRegistry) Sources() []string {
	if c == nil {
		return nil
	}
	out := make([]string, 0, len(c.sources))
	for _, source := range c.sources {
		out = append(out, source.ID)
	}
	return out
}

func (c *CompositeRegistry) Candidates(capability string, goal CapabilityGoal) []CapabilityCandidate {
	if c == nil {
		return nil
	}
	var out []CapabilityCandidate
	for _, source := range c.sources {
		candidates := source.Registry.Candidates(capability, goal)
		for _, candidate := range candidates {
			candidate.SourceID = source.ID
			if candidate.Reason == "" {
				candidate.Reason = "candidate from source " + source.ID
			}
			out = append(out, candidate)
		}
	}
	return out
}

func (c *CompositeRegistry) Execute(ctx context.Context, req CapabilityExecutionRequest) (CapabilityExecutionResult, error) {
	if c == nil || len(c.sources) == 0 {
		return CapabilityExecutionResult{}, fmt.Errorf("composite registry: no sources configured")
	}

	sourceID := strings.TrimSpace(req.SourceID)
	if sourceID == "" {
		if len(c.sources) != 1 {
			return CapabilityExecutionResult{}, fmt.Errorf("composite registry: SourceID is required when %d sources are configured", len(c.sources))
		}
		sourceID = c.sources[0].ID
	}
	registry, ok := c.byID[sourceID]
	if !ok {
		return CapabilityExecutionResult{}, fmt.Errorf("composite registry: unknown source %q", sourceID)
	}

	// SourceID belongs to the composite dispatch boundary. Leaf registries keep
	// their existing execution contract and see their source-local WorkerID.
	leafRequest := req
	leafRequest.SourceID = ""
	return registry.Execute(ctx, leafRequest)
}
