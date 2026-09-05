package experience

import (
	"fmt"
	"time"

	"tlaloc.local/behaviorlab/prototypelab"
	"tonal.local/runtime/tonal"
)

// EvaluatedRun keeps TONAL's native trace and the experiment's explicit
// correctness judgment side by side. This prevents the recording layer from
// inventing semantic correctness from execution status.
type EvaluatedRun struct {
	Record     tonal.RunRecord
	Evaluation Evaluation
}

// Episodes projects a set of evaluated TONAL workflow traces into the common
// Prototype Lab representation.
func Episodes(runID, sourceExperiment string, runs []EvaluatedRun) ([]prototypelab.Episode, error) {
	out := make([]prototypelab.Episode, 0, len(runs))
	for index, run := range runs {
		ep, err := FromRunRecord(runID, sourceExperiment, run.Record, run.Evaluation)
		if err != nil {
			return nil, fmt.Errorf("tonal experience: run[%d]: %w", index, err)
		}
		out = append(out, ep)
	}
	return out, nil
}

// WriteBundle is the one-call persistence boundary intended for TONAL
// prototype harnesses after they finish a batch of evaluated workflows.
// Native RunRecords remain authoritative; this writes only the reusable
// experience projection.
func WriteBundle(outDir string, manifest prototypelab.RunManifest, runs []EvaluatedRun, observedAt time.Time) (prototypelab.BundlePaths, error) {
	if manifest.RunID == "" {
		return prototypelab.BundlePaths{}, fmt.Errorf("tonal experience: manifest run_id is empty")
	}
	if manifest.SourceExperiment == "" {
		return prototypelab.BundlePaths{}, fmt.Errorf("tonal experience: manifest source_experiment is empty")
	}
	episodes, err := Episodes(manifest.RunID, manifest.SourceExperiment, runs)
	if err != nil {
		return prototypelab.BundlePaths{}, err
	}
	return prototypelab.WriteBundle(outDir, manifest, episodes, observedAt)
}
