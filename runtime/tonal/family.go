package tonal

import (
	"fmt"
	"strings"
)

// TaskFamily is a GENERIC workflow template. It is allowed to know the
// shape of a class of tasks ("locate two values, compare them"); it must
// never carry a per-instance solution plan, a benchmark id, an expected
// answer, or a hidden evidence address. Every Step names a capability, not
// an executor — TONAL resolves the executor through the Registry.
type TaskFamily struct {
	ID    string `json:"id"`
	Goal  string `json:"goal"` // human description, not machine-consumed
	Steps []Step `json:"steps"`
}

// Step is one logical unit of the template.
type Step struct {
	// LocalID is unique within the family and is how later Steps and the
	// InputSpec refer to this Step's output.
	LocalID string `json:"local_id"`
	// Capability is the CapabilityGoal for this Step.
	Capability string `json:"capability"`
	// Role disambiguates repeated capabilities in one workflow (e.g.
	// "A"/"B"). It is descriptive only.
	Role string `json:"role,omitempty"`
	// DependsOn lists LocalIDs whose output this Step consumes.
	DependsOn []string `json:"depends_on,omitempty"`
	// Input declares how to build this Step's CapabilityRequest input from
	// instance parameters and upstream observations. No executor specifics.
	Input InputSpec `json:"input"`
	// PreferDeterministic is passed through to resolution; it never names
	// a worker.
	PreferDeterministic bool `json:"prefer_deterministic,omitempty"`
}

// InputSpec is a flat JSON object template. String leaves may contain
// placeholders:
//
//	${param:NAME}          -> instance Params[NAME]
//	${obs:LOCALID}         -> the upstream Step's observation value (raw)
//	${obs:LOCALID:FIELD}   -> a top-level string/number field of that value
//
// Non-string leaves (numbers, bools) pass through unchanged. This is
// deliberately not a general expression language.
type InputSpec struct {
	Template map[string]any `json:"template"`
}

// Normalize validates the family: unique LocalIDs, resolvable DependsOn, a
// DAG (no cycles), non-empty capabilities.
func (f TaskFamily) Normalize() (TaskFamily, error) {
	f.ID = strings.TrimSpace(f.ID)
	if f.ID == "" {
		return TaskFamily{}, fmt.Errorf("task family: id is required")
	}
	if len(f.Steps) == 0 {
		return TaskFamily{}, fmt.Errorf("task family %s: at least one step is required", f.ID)
	}
	seen := map[string]int{}
	for index := range f.Steps {
		step := &f.Steps[index]
		step.LocalID = strings.TrimSpace(step.LocalID)
		step.Capability = strings.ToUpper(strings.TrimSpace(step.Capability))
		step.Role = strings.TrimSpace(step.Role)
		if step.LocalID == "" || step.Capability == "" {
			return TaskFamily{}, fmt.Errorf("task family %s: step %d needs local_id and capability", f.ID, index)
		}
		if _, dup := seen[step.LocalID]; dup {
			return TaskFamily{}, fmt.Errorf("task family %s: duplicate local_id %q", f.ID, step.LocalID)
		}
		seen[step.LocalID] = index
	}
	for _, step := range f.Steps {
		for _, dep := range step.DependsOn {
			depIndex, ok := seen[dep]
			if !ok {
				return TaskFamily{}, fmt.Errorf("task family %s: step %q depends on unknown %q", f.ID, step.LocalID, dep)
			}
			if depIndex >= seen[step.LocalID] {
				return TaskFamily{}, fmt.Errorf("task family %s: step %q depends on later step %q (not a forward DAG)", f.ID, step.LocalID, dep)
			}
		}
	}
	return f, nil
}

// CriticalPathDepth is the mechanical workflow-depth metric: the number of
// sequential capability executions on the longest dependency chain. It is
// computed from the template alone, before any inference.
func (f TaskFamily) CriticalPathDepth() int {
	index := map[string]Step{}
	for _, step := range f.Steps {
		index[step.LocalID] = step
	}
	memo := map[string]int{}
	var depth func(id string) int
	depth = func(id string) int {
		if value, ok := memo[id]; ok {
			return value
		}
		best := 0
		for _, dep := range index[id].DependsOn {
			if d := depth(dep); d > best {
				best = d
			}
		}
		memo[id] = best + 1
		return memo[id]
	}
	max := 0
	for _, step := range f.Steps {
		if d := depth(step.LocalID); d > max {
			max = d
		}
	}
	return max
}

// topoOrder returns the Steps in a valid topological order. Normalize
// enforces that every dependency is declared earlier, so declaration order
// is already topological.
func (f TaskFamily) topoOrder() []Step { return f.Steps }
