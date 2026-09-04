package tlalocbridge

import (
	"context"
	"encoding/json"
	"testing"

	"tlaloc.local/behaviorlab/tlaloquekit"
	"tonal.local/runtime/tonal"
)

// The bridge is the seam where a real Tlaloc registry meets the engine.
// This proves the engine runs end to end against the actual published
// deterministic Tlaloque set (no model involved).
func TestEngineRunsAgainstRealDeterministicTlaloques(t *testing.T) {
	registry, err := Build(Config{})
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	family := tonal.TaskFamily{
		ID:   "ARITHMETIC_THEN_COMPARE",
		Goal: "compute A-B and compare it with a threshold",
		Steps: []tonal.Step{
			{LocalID: "diff", Capability: "ARITHMETIC", Input: tonal.InputSpec{Template: map[string]any{
				"operation": "SUBTRACT", "a": "${param:a}", "b": "${param:b}",
			}}},
			{LocalID: "cmp", Capability: "COMPARE_NUMBERS", DependsOn: []string{"diff"}, Input: tonal.InputSpec{Template: map[string]any{
				"a": "${obs:diff:result}", "b": "${param:threshold}",
			}}},
		},
	}
	instance := tonal.Instance{ID: "wf-real-1", Family: family.ID, DeclaredDepth: 2, Params: map[string]string{
		"a": "150", "b": "40", "threshold": "100",
	}}

	record, _, err := (&tonal.Engine{Registry: registry}).RunWorkflow(context.Background(), family, instance, tonal.HeterogeneousPolicy{})
	if err != nil {
		t.Fatalf("RunWorkflow: %v", err)
	}
	if record.FinalStatus != "OK" {
		t.Fatalf("final status = %q (%s)", record.FinalStatus, record.Error)
	}
	var compare struct {
		Comparison string `json:"comparison"`
	}
	if err := json.Unmarshal(record.FinalValue.Value, &compare); err != nil {
		t.Fatalf("decode final value: %v", err)
	}
	if compare.Comparison != "GREATER" { // 150-40 = 110 > 100
		t.Fatalf("A-B vs threshold: got %q, want GREATER", compare.Comparison)
	}
	if record.Accounting.DeterministicOps != 2 || record.Accounting.GenerativeCalls != 0 {
		t.Fatalf("expected 2 deterministic ops and no generative calls, got %+v", record.Accounting)
	}
	// every candidate the registry offered for these capabilities is
	// deterministic
	for _, step := range record.Steps {
		for _, candidate := range step.Candidates {
			if candidate.Descriptor.Engine == tlaloquekit.EngineGenerative {
				t.Fatalf("step %s unexpectedly had a generative candidate", step.LocalID)
			}
		}
	}
}
