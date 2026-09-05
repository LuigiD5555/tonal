// Package experience adapts TONAL's native RunRecord into Tlaloc's public
// Prototype Lab projection. TONAL RunRecord remains the execution source of
// truth; this package only produces reusable cross-prototype Episodes.
package experience

import (
	"fmt"

	"tlaloc.local/behaviorlab/prototypelab"
	"tonal.local/runtime/tonal"
)

// Evaluation is deliberately supplied by the experiment harness instead of
// inferred from RunRecord.FinalStatus. A workflow can execute successfully and
// still be semantically wrong, so TONAL must not confuse completion with task
// correctness.
type Evaluation struct {
	Success          bool
	SemanticCorrect  bool
	ExactCorrect     bool
	FailureRootCause string
}

// FromRunRecord projects one TONAL workflow execution into one Episode.
// runID and sourceExperiment are required because TONAL's native RunRecord is
// workflow-scoped and intentionally does not invent campaign provenance.
func FromRunRecord(runID, sourceExperiment string, record tonal.RunRecord, evaluation Evaluation) (prototypelab.Episode, error) {
	if runID == "" {
		return prototypelab.Episode{}, fmt.Errorf("tonal experience: run_id is empty")
	}
	if sourceExperiment == "" {
		return prototypelab.Episode{}, fmt.Errorf("tonal experience: source_experiment is empty")
	}
	if record.WorkflowID == "" {
		return prototypelab.Episode{}, fmt.Errorf("tonal experience: workflow_id is empty")
	}

	steps := make([]prototypelab.Step, 0, len(record.Steps))
	for index, trace := range record.Steps {
		status := "OK"
		if trace.Error != "" {
			status = "ERROR"
		}
		steps = append(steps, prototypelab.Step{
			RequestIndex:       index,
			NodeID:             trace.NodeID,
			SelectedCapability: trace.Capability,
			ExecutorID:         trace.SelectedWorker,
			RawOutput:          trace.OutputJSON,
			Status:             status,
			ModelCalls:         trace.ModelCalls,
			LatencyMS:          trace.LatencyMS,
			Error:              trace.Error,
		})
	}

	failureRootCause := evaluation.FailureRootCause
	if failureRootCause == "" && !evaluation.Success {
		if record.FinalStatus != "" && record.FinalStatus != "OK" {
			failureRootCause = "TONAL_" + record.FinalStatus
		} else {
			failureRootCause = "EVALUATION_FAILED"
		}
	}

	return prototypelab.Episode{
		Schema:            prototypelab.EpisodeSchema,
		EpisodeID:         fmt.Sprintf("%s-%s-%s", runID, record.WorkflowID, record.Arm),
		SourceExperiment:  sourceExperiment,
		RunID:             runID,
		TaskID:            record.WorkflowID,
		Goal:              record.Goal,
		Arm:               record.Arm,
		Family:            record.Family,
		Steps:             steps,
		Success:           evaluation.Success,
		SemanticCorrect:   evaluation.SemanticCorrect,
		ExactCorrect:      evaluation.ExactCorrect,
		TerminalStatus:    record.FinalStatus,
		FailureRootCause:  failureRootCause,
		Cost: prototypelab.Cost{
			ModelCalls: record.Accounting.GenerativeCalls,
			LatencyMS:  record.Accounting.TotalLatencyMS,
		},
	}, nil
}
