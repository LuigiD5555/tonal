// Package r2registry wires the post-T1 Tonal capability sources without
// collapsing their ownership boundaries.
package r2registry

import (
	"fmt"

	"tonal.local/runtime/parrotbridge"
	"tonal.local/runtime/tlalocbridge"
	"tonal.local/runtime/tonal"
)

const (
	SourceTlaloc = "tlaloc"
	SourceParrot = "parrot"
)

// Config builds Tonal's R2 registry. Tlaloc machinery is always present.
// Parrot is optional and requires an explicit provider-neutral client plus
// bounded capability specs. Additional sources may publish Machines or Tools.
type Config struct {
	ParrotClient       parrotbridge.Client
	ParrotCapabilities []parrotbridge.CapabilitySpec
	AdditionalSources  []tonal.RegistrySource
}

func Build(cfg Config) (*tonal.CompositeRegistry, error) {
	machinery, err := tlalocbridge.Build(tlalocbridge.Config{})
	if err != nil {
		return nil, fmt.Errorf("r2registry: build Tlaloc machinery source: %w", err)
	}

	sources := []tonal.RegistrySource{{ID: SourceTlaloc, Registry: machinery}}

	if cfg.ParrotClient != nil || len(cfg.ParrotCapabilities) != 0 {
		if cfg.ParrotClient == nil {
			return nil, fmt.Errorf("r2registry: Parrot capabilities require ParrotClient")
		}
		parrot, err := parrotbridge.New(cfg.ParrotClient, cfg.ParrotCapabilities...)
		if err != nil {
			return nil, fmt.Errorf("r2registry: build Parrot source: %w", err)
		}
		sources = append(sources, tonal.RegistrySource{ID: SourceParrot, Registry: parrot})
	}

	for _, source := range cfg.AdditionalSources {
		sources = append(sources, source)
	}
	registry, err := tonal.NewCompositeRegistry(sources...)
	if err != nil {
		return nil, fmt.Errorf("r2registry: compose sources: %w", err)
	}
	return registry, nil
}
