package tonal

import "tlaloc.local/behaviorlab/tlaloquekit"

// RunRecord is the mandatory deterministic trace of one workflow
// execution: goal, resolved DAG, every capability requested, candidate
// Tlaloques and rejection reasons, the selected Tlaloque, Blackboard reads
// and writes, deterministic transforms, model calls, and the final
// verification/answer. It is what proves TONAL actually routed the work.
type RunRecord struct {
	WorkflowID        string `json:"workflow_id"`
	Family            string `json:"family"`
	Goal              string `json:"goal"`
	Arm               string `json:"arm"`
	DeclaredDepth     int    `json:"declared_depth"`
	CriticalPathDepth int    `json:"critical_path_depth"`

	Steps []StepTrace `json:"steps"`

	FinalKey    string                  `json:"final_key"`
	FinalValue  tlaloquekit.Observation `json:"final_value"`
	FinalStatus string                  `json:"final_status"` // OK | UNKNOWN | UNSUPPORTED | CONTRACT_FAILURE | ERROR
	Error       string                  `json:"error,omitempty"`

	Accounting Accounting `json:"accounting"`
}

// StepTrace records one node's resolution and execution.
type StepTrace struct {
	LocalID    string `json:"local_id"`
	Role       string `json:"role,omitempty"`
	Capability string `json:"capability"`
	NodeID     string `json:"node_id"`

	Candidates      []tlaloquekit.Candidate `json:"candidates"`
	SelectedWorker  string                  `json:"selected_worker"`
	SelectionReason string                  `json:"selection_reason"`
	EngineKind      string                  `json:"engine_kind"`

	BlackboardReads  []string `json:"blackboard_reads"`
	BlackboardWrites []string `json:"blackboard_writes"`

	InputJSON  string `json:"input_json"`
	OutputJSON string `json:"output_json,omitempty"`

	DeterministicTransforms int   `json:"deterministic_transforms"`
	ModelCalls              int   `json:"model_calls"`
	LatencyMS               int64 `json:"latency_ms"`

	ProfileVersion string `json:"profile_version,omitempty"`
	Notes          string `json:"notes,omitempty"`
	Error          string `json:"error,omitempty"`
}
