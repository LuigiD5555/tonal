// Package integration proves the Tonal runtime can use the frozen Tlaloc R1
// Parrot adapter while exposing it through Architecture R2 as external
// probabilistic cognition, with a fake OpenAI-compatible endpoint.
package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"tlaloc.local/behaviorlab/tlaloquekit"
	"tonal.local/runtime/tlalocbridge"
	"tonal.local/runtime/tonal"
)

func pinnedTlaloc(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", "..", ".work", "components", "tlaloc", "behavior-lab"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "profiles", "parrot-lfm2-vl-1.6b-r1.json")); err != nil {
		t.Skipf("pinned Tlaloc checkout not materialised (run scripts/fetch-components.sh): %v", err)
	}
	return root
}

type fakeEndpoint struct {
	server *httptest.Server
	calls  int64
	reply  string
}

func newFakeEndpoint(reply string) *fakeEndpoint {
	fake := &fakeEndpoint{reply: reply}
	fake.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&fake.calls, 1)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []any{map[string]any{"message": map[string]any{"content": fake.reply}}},
			"usage":   map[string]any{"prompt_tokens": 20, "completion_tokens": 2},
		})
	}))
	return fake
}

func (f *fakeEndpoint) close()       { f.server.Close() }
func (f *fakeEndpoint) count() int64 { return atomic.LoadInt64(&f.calls) }

func parrotRegistry(t *testing.T, endpointURL string) *tlalocbridge.Registry {
	t.Helper()
	registry, err := tlalocbridge.Build(tlalocbridge.Config{
		Parrot: &tlalocbridge.ParrotConfig{
			ProfilePath:         filepath.Join(pinnedTlaloc(t), "profiles", "parrot-lfm2-vl-1.6b-r1.json"),
			ExpectedProfileHash: "8acc959b",
			Endpoint:            tlaloquekit.ParrotEndpoint{BaseURL: endpointURL, Model: "lfm2-vl-1.6b", MaxTokens: 16},
			WorkDir:             t.TempDir(),
		},
	})
	if err != nil {
		t.Fatalf("build registry: %v", err)
	}
	return registry
}

func extractFamily(regionTemplate map[string]any) tonal.TaskFamily {
	return tonal.TaskFamily{
		ID:   "EXTRACT_ONE_NUMBER",
		Goal: "read one number from a located region",
		Steps: []tonal.Step{{
			LocalID: "extract", Capability: "EXTRACT_NUMBER",
			Input: tonal.InputSpec{Template: map[string]any{
				"opcode":          "EXTRACT_NUMBER",
				"page_image_path": "${param:page_image}",
				"region":          regionTemplate,
			}},
		}},
	}
}

func run(t *testing.T, registry tonal.CapabilityRegistry, family tonal.TaskFamily, params map[string]string) (tonal.RunRecord, tonal.StepTrace) {
	t.Helper()
	record, _, err := (&tonal.Engine{Registry: registry}).RunWorkflow(context.Background(), family,
		tonal.Instance{ID: "wf-int", Family: family.ID, DeclaredDepth: 1, Params: params}, tonal.HeterogeneousPolicy{})
	if err != nil {
		t.Fatalf("RunWorkflow: %v", err)
	}
	return record, record.Steps[len(record.Steps)-1]
}

func decodeParrotOutput(t *testing.T, step tonal.StepTrace) map[string]any {
	t.Helper()
	var out map[string]any
	if err := json.Unmarshal([]byte(step.OutputJSON), &out); err != nil {
		t.Fatalf("decode parrot output: %v", err)
	}
	return out
}

func TestParrot_MissingOperand_ZeroModelCalls_Behind_The_Adapter(t *testing.T) {
	fake := newFakeEndpoint("999")
	defer fake.close()
	registry := parrotRegistry(t, fake.server.URL+"/v1")

	family := extractFamily(map[string]any{"page": 1})
	record, step := run(t, registry, family, map[string]string{"page_image": ""})

	if fake.count() != 0 {
		t.Fatalf("expected zero model calls, got %d", fake.count())
	}
	if record.FinalStatus != "UNSUPPORTED" {
		t.Fatalf("final status = %q, want UNSUPPORTED", record.FinalStatus)
	}
	if step.SelectedKind != tonal.CapabilityExternalModel {
		t.Fatalf("EXTRACT_NUMBER kind=%q, want EXTERNAL_MODEL", step.SelectedKind)
	}
	if step.EngineKind != string(tlaloquekit.EngineGenerative) {
		t.Fatalf("frozen R1 adapter should still report GENERATIVE engine, got %q", step.EngineKind)
	}
	if step.ModelCalls != 0 || record.Accounting.ParrotCalls != 0 {
		t.Fatalf("accounting must show zero external/Parrot calls, got step=%d acct=%d", step.ModelCalls, record.Accounting.ParrotCalls)
	}
}

func TestParrot_LowScaleOperand_UpscalesThenOneCall(t *testing.T) {
	fake := newFakeEndpoint("512")
	defer fake.close()
	registry := parrotRegistry(t, fake.server.URL+"/v1")

	page := filepath.Join(pinnedTlaloc(t), "internal", "parrotpresent", "testdata", "page.png")
	family := extractFamily(map[string]any{
		"page":        1,
		"bbox":        map[string]any{"x1": 20.0, "y1": 100.0, "x2": 180.0, "y2": 108.0},
		"page_width":  200.0,
		"page_height": 300.0,
	})
	record, step := run(t, registry, family, map[string]string{"page_image": page})

	if fake.count() != 1 {
		t.Fatalf("expected exactly one model call, got %d", fake.count())
	}
	if record.FinalStatus != "OK" || record.Accounting.ParrotCalls != 1 {
		t.Fatalf("status=%q parrotCalls=%d", record.FinalStatus, record.Accounting.ParrotCalls)
	}
	if step.SelectedKind != tonal.CapabilityExternalModel {
		t.Fatalf("selected kind=%q, want EXTERNAL_MODEL", step.SelectedKind)
	}
	out := decodeParrotOutput(t, step)
	decision, _ := out["adapter_decision"].(map[string]any)
	transforms, _ := decision["transformations"].([]any)
	if !hasTransform(transforms, "UPSCALE_TO_PREFERRED") {
		t.Fatalf("low-scale operand did not trigger UPSCALE_TO_PREFERRED: %v", transforms)
	}
	if step.ProfileVersion == "" {
		t.Fatalf("step trace missing the profile version for the external-model call")
	}
}

func TestParrot_AdequateScaleOperand_NoUnnecessaryUpscale(t *testing.T) {
	fake := newFakeEndpoint("77")
	defer fake.close()
	registry := parrotRegistry(t, fake.server.URL+"/v1")

	page := filepath.Join(pinnedTlaloc(t), "internal", "parrotpresent", "testdata", "page.png")
	family := extractFamily(map[string]any{
		"page":        1,
		"bbox":        map[string]any{"x1": 20.0, "y1": 100.0, "x2": 180.0, "y2": 140.0},
		"page_width":  200.0,
		"page_height": 300.0,
	})
	_, step := run(t, registry, family, map[string]string{"page_image": page})

	if fake.count() != 1 {
		t.Fatalf("expected exactly one model call, got %d", fake.count())
	}
	out := decodeParrotOutput(t, step)
	decision, _ := out["adapter_decision"].(map[string]any)
	transforms, _ := decision["transformations"].([]any)
	if hasTransform(transforms, "UPSCALE_TO_PREFERRED") {
		t.Fatalf("adequate-scale operand should not be upscaled: %v", transforms)
	}
}

func TestParrot_ProfileHashIsHardValidated(t *testing.T) {
	_, err := tlalocbridge.Build(tlalocbridge.Config{
		Parrot: &tlalocbridge.ParrotConfig{
			ProfilePath:         filepath.Join(pinnedTlaloc(t), "profiles", "parrot-lfm2-vl-1.6b-r1.json"),
			ExpectedProfileHash: "0000000000",
			Endpoint:            tlaloquekit.ParrotEndpoint{BaseURL: "http://127.0.0.1:0", Model: "m"},
		},
	})
	if err == nil || !strings.Contains(err.Error(), "hash") {
		t.Fatalf("expected a profile hash validation error, got %v", err)
	}
}

func hasTransform(transforms []any, kind string) bool {
	for _, entry := range transforms {
		object, ok := entry.(map[string]any)
		if ok && object["kind"] == kind {
			return true
		}
	}
	return false
}
