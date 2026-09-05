// Package tlalocbridge is the ONLY place in the Tonal runtime that knows the
// concrete Tlaloc publication contract. It translates Tlaloc's frozen R1
// types into Tonal-owned capability contracts.
package tlalocbridge

import (
	"context"
	"fmt"

	"tlaloc.local/behaviorlab/tlaloquekit"
	"tlaloc.local/behaviorlab/tlaloquekit/tlalocregistry"
	"tonal.local/runtime/tonal"
)

type Config struct {
	OmitDeterministic bool
	Parrot            *ParrotConfig
}

type ParrotConfig struct {
	ProfilePath         string
	ExpectedProfileHash string
	Endpoint            tlaloquekit.ParrotEndpoint
	WorkDir             string
}

// Registry adapts a Tlaloc-qualified registry into Tonal's generic runtime
// seam. Tlaloc remains one machinery source; this adapter does not make it the
// root type of Tonal's Registry.
type Registry struct {
	inner tlaloquekit.QualifiedRegistry
}

var _ tonal.CapabilityRegistry = (*Registry)(nil)

func Build(cfg Config) (*Registry, error) {
	innerCfg := tlalocregistry.Config{OmitDeterministic: cfg.OmitDeterministic}
	if cfg.Parrot != nil {
		innerCfg.Parrot = &tlalocregistry.ParrotConfig{
			ProfilePath:         cfg.Parrot.ProfilePath,
			ExpectedProfileHash: cfg.Parrot.ExpectedProfileHash,
			Endpoint:            cfg.Parrot.Endpoint,
			WorkDir:             cfg.Parrot.WorkDir,
		}
	}
	inner, err := tlalocregistry.BuildQualifiedRegistry(innerCfg)
	if err != nil {
		return nil, fmt.Errorf("tlalocbridge: build qualified registry: %w", err)
	}
	return &Registry{inner: inner}, nil
}

// frozenR1Kind translates the historical T1 publication contract. R1 called
// Parrot a GENERATIVE Tlaloque; R2 keeps those bytes/history intact but
// exposes it to Tonal as EXTERNAL_MODEL. Future publication contracts should
// carry component Kind explicitly rather than infer it from EngineKind.
func frozenR1Kind(d tlaloquekit.Descriptor) tonal.CapabilityKind {
	if d.Engine == tlaloquekit.EngineGenerative {
		return tonal.CapabilityExternalModel
	}
	return tonal.CapabilityTlaloque
}

func descriptor(d tlaloquekit.Descriptor) tonal.CapabilityDescriptor {
	return tonal.CapabilityDescriptor{
		ID:             d.ID,
		Capability:     d.Capability,
		Kind:           frozenR1Kind(d),
		EngineKind:     string(d.Engine),
		Deterministic:  d.Deterministic,
		ParameterCount: d.ParameterCount,
		InputSchema:    d.InputSchema,
		OutputSchema:   d.OutputSchema,
		Dependencies:   append([]string(nil), d.Dependencies...),
		ProfileRef:     d.ProfileRef,
		EvidenceRef:    d.EvidenceRef,
	}
}

func (r *Registry) Capabilities() []tonal.CapabilityDescriptor {
	published := r.inner.Capabilities()
	out := make([]tonal.CapabilityDescriptor, 0, len(published))
	for _, d := range published {
		out = append(out, descriptor(d))
	}
	return out
}

func (r *Registry) Candidates(capability string, goal tonal.CapabilityGoal) []tonal.CapabilityCandidate {
	published := r.inner.Candidates(capability, tlaloquekit.Goal{
		Capability:          goal.Capability,
		PreferDeterministic: goal.PreferDeterministic,
		MaxParameters:       goal.MaxParameters,
		AvailableProducts:   append([]string(nil), goal.AvailableProducts...),
	})
	out := make([]tonal.CapabilityCandidate, 0, len(published))
	for _, candidate := range published {
		d := candidate.Descriptor
		out = append(out, tonal.CapabilityCandidate{
			WorkerID:      d.ID,
			Capability:    d.Capability,
			Kind:          frozenR1Kind(d),
			EngineKind:    string(d.Engine),
			ProfileRef:    d.ProfileRef,
			Deterministic: d.Deterministic,
			Generative:    d.Engine == tlaloquekit.EngineGenerative,
			Selected:      candidate.Selected,
			Reason:        candidate.Reason,
		})
	}
	return out
}

func toTlalocObservation(obs tonal.Observation) tlaloquekit.Observation {
	return tlaloquekit.Observation{
		Producer:       obs.Producer,
		Capability:     obs.Capability,
		Key:            obs.Key,
		Value:          obs.Value,
		Kind:           obs.Kind,
		Status:         obs.Status,
		Confidence:     obs.Confidence,
		References:     append([]string(nil), obs.References...),
		Provenance:     obs.Provenance,
		ProfileVersion: obs.ProfileVersion,
		RecordedAt:     obs.RecordedAt,
	}
}

func toTonalObservation(obs tlaloquekit.Observation) tonal.Observation {
	return tonal.Observation{
		Producer:       obs.Producer,
		Capability:     obs.Capability,
		Key:            obs.Key,
		Value:          obs.Value,
		Kind:           obs.Kind,
		Status:         obs.Status,
		Confidence:     obs.Confidence,
		References:     append([]string(nil), obs.References...),
		Provenance:     obs.Provenance,
		ProfileVersion: obs.ProfileVersion,
		RecordedAt:     obs.RecordedAt,
	}
}

func (r *Registry) Execute(ctx context.Context, req tonal.CapabilityExecutionRequest) (tonal.CapabilityExecutionResult, error) {
	prior := make([]tlaloquekit.Observation, 0, len(req.PriorObservations))
	for _, obs := range req.PriorObservations {
		prior = append(prior, toTlalocObservation(obs))
	}
	result, err := r.inner.Execute(ctx, tlaloquekit.ExecutionRequest{
		TaskID:            req.TaskID,
		NodeID:            req.NodeID,
		Capability:        req.Capability,
		WorkerID:          req.WorkerID,
		Input:             req.Input,
		PriorObservations: prior,
	})
	if err != nil {
		return tonal.CapabilityExecutionResult{}, err
	}
	observations := make([]tonal.Observation, 0, len(result.Observations))
	for _, obs := range result.Observations {
		observations = append(observations, toTonalObservation(obs))
	}
	out := tonal.CapabilityExecutionResult{
		WorkerID:     result.WorkerID,
		Output:       result.Output,
		Confidence:   result.Confidence,
		Notes:        result.Notes,
		Observations: observations,
	}
	if result.Usage != nil {
		out.Usage = &tonal.CapabilityUsage{
			PromptTokens:     result.Usage.PromptTokens,
			CompletionTokens: result.Usage.CompletionTokens,
			ModelCalls:       result.Usage.ModelCalls,
		}
	}
	return out, nil
}

func (r *Registry) ParrotProfileID() string   { return r.inner.ParrotProfileID() }
func (r *Registry) ParrotProfileHash() string { return r.inner.ParrotProfileHash() }
