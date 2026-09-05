package r2registry

import (
	"context"
	"encoding/json"
	"testing"

	"tonal.local/runtime/parrotbridge"
	"tonal.local/runtime/tonal"
)

type fakeParrotClient struct {
	calls int
}

func (f *fakeParrotClient) Invoke(_ context.Context, req parrotbridge.Request) (parrotbridge.Result, error) {
	f.calls++
	if req.Capability == tonal.ResponseRenderCapability {
		return parrotbridge.Result{Output: json.RawMessage(`{"text":"ok","used_keys":["answer"]}`), Usage: &tonal.CapabilityUsage{ModelCalls: 1}}, nil
	}
	return parrotbridge.Result{Output: json.RawMessage(`{"trimmed":"7"}`), Usage: &tonal.CapabilityUsage{ModelCalls: 1}}, nil
}

func TestBuildSeparatesTlalocMachineryAndParrotSources(t *testing.T) {
	client := &fakeParrotClient{}
	registry, err := Build(Config{
		ParrotClient: client,
		ParrotCapabilities: []parrotbridge.CapabilitySpec{
			{Capability: "NORMALIZE", WorkerID: "default"},
			parrotbridge.ResponseRendererSpec("model@r2"),
		},
	})
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	sources := registry.Sources()
	if len(sources) != 2 || sources[0] != SourceTlaloc || sources[1] != SourceParrot {
		t.Fatalf("sources=%v", sources)
	}

	candidates := registry.Candidates("NORMALIZE", tonal.CapabilityGoal{Capability: "NORMALIZE"})
	var haveTlaloc, haveParrot bool
	for _, candidate := range candidates {
		switch candidate.SourceID {
		case SourceTlaloc:
			haveTlaloc = haveTlaloc || candidate.Kind == tonal.CapabilityTlaloque
		case SourceParrot:
			haveParrot = haveParrot || candidate.Kind == tonal.CapabilityExternalModel
		}
	}
	if !haveTlaloc || !haveParrot {
		t.Fatalf("NORMALIZE candidates did not preserve source ontology: %+v", candidates)
	}

	selected, _, ok := selectForTest(tonal.MachineryFirstPolicy{}, tonal.Step{Capability: "NORMALIZE"}, candidates)
	if !ok || selected.SourceID != SourceTlaloc || selected.Kind != tonal.CapabilityTlaloque {
		t.Fatalf("machinery-first selected %+v", selected)
	}
}

// selectForTest exercises the public policy surface without depending on
// Tonal's unexported resolver. Worker IDs in this assembled registry are
// source-distinct in practice for this test, so the policy choice is enough.
func selectForTest(policy tonal.SelectionPolicy, step tonal.Step, candidates []tonal.CapabilityCandidate) (tonal.CapabilityCandidate, string, bool) {
	workerID, reason := policy.SelectWorker(step, candidates)
	for _, candidate := range candidates {
		if candidate.WorkerID == workerID && candidate.SourceID == SourceTlaloc && candidate.Kind != tonal.CapabilityExternalModel {
			return candidate, reason, true
		}
	}
	for _, candidate := range candidates {
		if candidate.WorkerID == workerID {
			return candidate, reason, true
		}
	}
	return tonal.CapabilityCandidate{}, reason, false
}

func TestBuildPublishesResponseRenderingOnlyThroughParrot(t *testing.T) {
	client := &fakeParrotClient{}
	registry, err := Build(Config{
		ParrotClient:       client,
		ParrotCapabilities: []parrotbridge.CapabilitySpec{parrotbridge.ResponseRendererSpec("model@r2")},
	})
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	candidates := registry.Candidates(tonal.ResponseRenderCapability, tonal.CapabilityGoal{Capability: tonal.ResponseRenderCapability})
	if len(candidates) != 1 {
		t.Fatalf("render candidates=%d, want 1: %+v", len(candidates), candidates)
	}
	if candidates[0].SourceID != SourceParrot || candidates[0].Kind != tonal.CapabilityExternalModel {
		t.Fatalf("render candidate=%+v", candidates[0])
	}
}

func TestBuildWithoutParrotKeepsTonalValid(t *testing.T) {
	registry, err := Build(Config{})
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if got := registry.Sources(); len(got) != 1 || got[0] != SourceTlaloc {
		t.Fatalf("sources=%v", got)
	}
	if candidates := registry.Candidates(tonal.ResponseRenderCapability, tonal.CapabilityGoal{}); len(candidates) != 0 {
		t.Fatalf("unexpected response renderer without Parrot: %+v", candidates)
	}
}

func TestTlalocExecutionDoesNotCallParrot(t *testing.T) {
	client := &fakeParrotClient{}
	registry, err := Build(Config{
		ParrotClient:       client,
		ParrotCapabilities: []parrotbridge.CapabilitySpec{{Capability: "NORMALIZE"}},
	})
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	candidates := registry.Candidates("ARITHMETIC", tonal.CapabilityGoal{Capability: "ARITHMETIC", PreferDeterministic: true})
	if len(candidates) == 0 {
		t.Fatal("expected Tlaloc ARITHMETIC machinery")
	}
	var selected *tonal.CapabilityCandidate
	for _, candidate := range candidates {
		if candidate.SourceID == SourceTlaloc {
			copyCandidate := candidate
			selected = &copyCandidate
			break
		}
	}
	if selected == nil {
		t.Fatalf("no Tlaloc candidate: %+v", candidates)
	}

	_, err = registry.Execute(context.Background(), tonal.CapabilityExecutionRequest{
		TaskID: "r", NodeID: "n", Capability: "ARITHMETIC", SourceID: selected.SourceID, WorkerID: selected.WorkerID,
		Input: json.RawMessage(`{"operation":"SUBTRACT","a":"10","b":"3"}`),
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.calls != 0 {
		t.Fatalf("Parrot calls=%d, want 0", client.calls)
	}
}
