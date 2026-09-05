// Package tonal is the TONAL runtime: goal intake, workflow/DAG expansion,
// capability resolution and routing, the scheduler/execution loop, the
// Blackboard, resource accounting and deterministic routing trace.
package tonal

import (
	"sort"
	"strings"
)

// Blackboard is Tonal's append-only observation ledger. Results from
// Tlaloques, Machines, Tools or Parrot land here through the Engine as typed
// Observations with provenance. Executors never mutate the Blackboard
// directly.
type Blackboard struct {
	taskID  string
	entries []Observation
}

func NewBlackboard(taskID string) *Blackboard {
	return &Blackboard{taskID: strings.TrimSpace(taskID)}
}

func (b *Blackboard) TaskID() string { return b.taskID }

func (b *Blackboard) Append(observations ...Observation) {
	b.entries = append(b.entries, observations...)
}

func (b *Blackboard) Snapshot() []Observation {
	out := make([]Observation, len(b.entries))
	copy(out, b.entries)
	return out
}

// Latest returns the most recently appended entry for a key, preferring a
// FACT over a raw OBSERVATION for the same key.
func (b *Blackboard) Latest(key string) (Observation, bool) {
	key = strings.TrimSpace(key)
	var obs Observation
	var fact Observation
	haveObs, haveFact := false, false
	for _, entry := range b.entries {
		if entry.Key != key {
			continue
		}
		if strings.EqualFold(entry.Kind, "FACT") {
			fact, haveFact = entry, true
		} else {
			obs, haveObs = entry, true
		}
	}
	if haveFact {
		return fact, true
	}
	if haveObs {
		return obs, true
	}
	return Observation{}, false
}

func (b *Blackboard) Keys() []string {
	seen := map[string]struct{}{}
	for _, entry := range b.entries {
		seen[entry.Key] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for key := range seen {
		out = append(out, key)
	}
	sort.Strings(out)
	return out
}
