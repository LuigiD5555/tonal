// Package tonal is the TONAL runtime: goal intake, workflow/DAG expansion,
// capability resolution and routing, the scheduler/execution loop, the
// Blackboard, resource accounting and the deterministic routing trace.
//
// It depends only on the public tlaloquekit.QualifiedRegistry contract. It
// holds no executor-specific knowledge: nothing in this package knows a
// particular model's pixel scale or cropping policy. A Tlaloque enforces
// its own competence envelope inside Execute.
package tonal

import (
	"sort"
	"strings"

	"tlaloc.local/behaviorlab/tlaloquekit"
)

// Blackboard is TONAL's append-only observation ledger. Every Tlaloque
// result lands here as a typed Observation with full provenance; a
// deterministic derived value keeps provenance back to its source
// observations. TONAL owns this — no Tlaloque writes to it directly.
type Blackboard struct {
	taskID  string
	entries []tlaloquekit.Observation
}

// NewBlackboard creates an empty ledger for one workflow run.
func NewBlackboard(taskID string) *Blackboard {
	return &Blackboard{taskID: strings.TrimSpace(taskID)}
}

// TaskID reports the run this Blackboard belongs to.
func (b *Blackboard) TaskID() string { return b.taskID }

// Append records observations in order. Model output arrives here as an
// OBSERVATION; it becomes a FACT only when a VERIFY Tlaloque emits one.
func (b *Blackboard) Append(observations ...tlaloquekit.Observation) {
	b.entries = append(b.entries, observations...)
}

// Snapshot returns a copy of every entry, oldest first.
func (b *Blackboard) Snapshot() []tlaloquekit.Observation {
	out := make([]tlaloquekit.Observation, len(b.entries))
	copy(out, b.entries)
	return out
}

// Latest returns the most recently appended entry for a key, preferring a
// FACT over a raw OBSERVATION for the same key.
func (b *Blackboard) Latest(key string) (tlaloquekit.Observation, bool) {
	key = strings.TrimSpace(key)
	var obs tlaloquekit.Observation
	var fact tlaloquekit.Observation
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
	return tlaloquekit.Observation{}, false
}

// Keys lists every distinct observation key, sorted.
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
