package tonal

// Accounting is the per-workflow resource picture. The ParrotCalls field is
// retained for T1 schema compatibility; in R2 it counts model calls performed
// by the selected EXTERNAL_MODEL component rather than relying on a worker-id
// naming convention.
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

	switch step.EngineKind {
	case "GENERATIVE":
		a.GenerativeCalls += step.ModelCalls
	case "SPECIALIST":
		a.SpecialistOps++
	case "DETERMINISTIC", "ALGORITHMIC":
		a.DeterministicOps++
	}

	if step.SelectedKind == CapabilityExternalModel {
		a.ParrotCalls += step.ModelCalls
	}
}
