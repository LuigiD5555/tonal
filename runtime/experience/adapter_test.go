package experience

import (
	"os"
	"testing"
	"time"

	"tlaloc.local/behaviorlab/prototypelab"
	"tonal.local/runtime/tonal"
)

func TestFromRunRecordDoesNotConfuseExecutionWithCorrectness(t *testing.T) {
	record := tonal.RunRecord{
		WorkflowID: "wf-001",
		Family:     "FAMILY_X",
		Goal:       "produce the right answer",
		Arm:        "B",
		FinalStatus: "OK", // execution completed, but evaluation below says answer was wrong.
		Steps: []tonal.StepTrace{{
			NodeID:         "wf-001::extract",
			Capability:     "EXTRACT_NUMBER",
			SelectedWorker: "parrot-r1",
			OutputJSON:     `{"value":41}`,
			ModelCalls:     1,
			LatencyMS:      18,
		}},
		Accounting: tonal.Accounting{GenerativeCalls: 1, TotalLatencyMS: 18},
	}

	ep, err := FromRunRecord("run-1", "TONAL_PROTO", record, Evaluation{
		Success: false, SemanticCorrect: false, ExactCorrect: false,
	})
	if err != nil {
		t.Fatalf("FromRunRecord: %v", err)
	}
	if ep.Success {
		t.Fatal("episode marked successful solely because FinalStatus was OK")
	}
	if ep.FailureRootCause != "EVALUATION_FAILED" {
		t.Errorf("FailureRootCause = %q, want EVALUATION_FAILED", ep.FailureRootCause)
	}
	if ep.RunID != "run-1" || ep.TaskID != "wf-001" || ep.Arm != "B" {
		t.Errorf("episode provenance = run:%q task:%q arm:%q", ep.RunID, ep.TaskID, ep.Arm)
	}
	if len(ep.Steps) != 1 || ep.Steps[0].RawOutput != `{"value":41}` || ep.Steps[0].ExecutorID != "parrot-r1" {
		t.Errorf("step projection = %+v", ep.Steps)
	}
	if ep.Cost.ModelCalls != 1 || ep.Cost.LatencyMS != 18 {
		t.Errorf("cost = %+v", ep.Cost)
	}
}

func TestWriteBundleUsesPublicPrototypeLabSurface(t *testing.T) {
	manifest := prototypelab.RunManifest{
		Schema:           prototypelab.ManifestSchema,
		RunID:            "run-bundle",
		SourceExperiment: "TONAL_PROTO",
		Prototype: prototypelab.Prototype{
			ID:      "IMPULSE_CONTROLLER",
			Version: "0.1",
		},
	}
	runs := []EvaluatedRun{{
		Record: tonal.RunRecord{
			WorkflowID: "wf-001", Family: "FAMILY_X", Goal: "answer", Arm: "C", FinalStatus: "OK",
			Accounting: tonal.Accounting{},
		},
		Evaluation: Evaluation{Success: true, SemanticCorrect: true, ExactCorrect: true},
	}}

	paths, err := WriteBundle(
		t.TempDir(), manifest, runs,
		time.Date(2026, 9, 4, 20, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("WriteBundle: %v", err)
	}
	if _, err := os.Stat(paths.Summary); err != nil {
		t.Fatalf("summary not written: %v", err)
	}
}
