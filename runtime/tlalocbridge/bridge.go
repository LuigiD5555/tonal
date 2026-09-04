// Package tlalocbridge is the ONLY place in the TONAL runtime that wires a
// concrete Tlaloc-published registry. It converts a runtime configuration
// into a tlaloquekit.QualifiedRegistry and hands it to the executor-agnostic
// tonal engine. It owns nothing cognitive: no DAG, no Blackboard, no
// scheduler, and no executor-specific competence knowledge.
package tlalocbridge

import (
	"fmt"

	"tlaloc.local/behaviorlab/tlaloquekit"
	"tlaloc.local/behaviorlab/tlaloquekit/tlalocregistry"
)

// Config is the runtime-supplied wiring configuration.
type Config struct {
	// OmitDeterministic is a test hook; the default publishes the full
	// deterministic Tlaloque set.
	OmitDeterministic bool

	// Parrot, when set, makes the R1-aware generative Tlaloque available.
	// It is the operator's explicit decision; there is no implicit
	// capability -> Parrot fallback.
	Parrot *ParrotConfig
}

// ParrotConfig mirrors the frozen-profile-gated Parrot wiring.
type ParrotConfig struct {
	ProfilePath         string
	ExpectedProfileHash string
	Endpoint            tlaloquekit.ParrotEndpoint
	WorkDir             string
}

// Build returns the qualified registry Tlaloc publishes, behind the public
// contract the engine consumes.
func Build(cfg Config) (tlaloquekit.QualifiedRegistry, error) {
	inner := tlalocregistry.Config{OmitDeterministic: cfg.OmitDeterministic}
	if cfg.Parrot != nil {
		inner.Parrot = &tlalocregistry.ParrotConfig{
			ProfilePath:         cfg.Parrot.ProfilePath,
			ExpectedProfileHash: cfg.Parrot.ExpectedProfileHash,
			Endpoint:            cfg.Parrot.Endpoint,
			WorkDir:             cfg.Parrot.WorkDir,
		}
	}
	registry, err := tlalocregistry.BuildQualifiedRegistry(inner)
	if err != nil {
		return nil, fmt.Errorf("tlalocbridge: build qualified registry: %w", err)
	}
	return registry, nil
}
