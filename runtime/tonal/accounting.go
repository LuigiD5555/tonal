package tonal

import "tlaloc.local/behaviorlab/tlaloquekit"

// Accounting is the per-workflow resource picture. TONAL records usage
// transparently rather than equalising it across arms — the resource
// difference between arms is part of what T1 measures.
type Accounting struct {
	TotalSteps              int   `json:"total_steps"`
	SuccessfulSteps         int   `json:"successful_steps"`
	ParrotCalls             int   `json:"parrot_calls"`
	GenerativeCalls         int   `json:"generative_calls"`
	DeterministicOps        int   `json:"deterministic_ops"`
	SpecialistOps           int   `json:"specialist_ops"`
	DeterministicTransforms int   `json:"deterministic_transforms"`
	TotalLatencyMS          int64 `json:"total_latency_ms"`
	MaxStepLatencyMS        int64 `json:"max_step_latency_ms"`
}

func (a *Accounting) recordStep(step StepTrace, ok bool) {
	a.TotalSteps++
	if ok {
		a.SuccessfulSteps++
	}
	a.DeterministicTransforms += step.DeterministicTransforms
	a.TotalLatencyMS += step.LatencyMS
	if step.LatencyMS > a.MaxStepLatencyMS {
		a.MaxStepLatencyMS = step.LatencyMS
	}
	switch tlaloquekit.EngineKind(step.EngineKind) {
	case tlaloquekit.EngineGenerative:
		a.GenerativeCalls += step.ModelCalls
		if isParrot(step.SelectedWorker) {
			a.ParrotCalls += step.ModelCalls
		}
	case tlaloquekit.EngineSpecialist:
		a.SpecialistOps++
	case tlaloquekit.EngineDeterministic, tlaloquekit.EngineAlgorithmic:
		a.DeterministicOps++
	}
}

func isParrot(workerID string) bool {
	return len(workerID) >= 6 && workerID[:6] == "parrot"
}
