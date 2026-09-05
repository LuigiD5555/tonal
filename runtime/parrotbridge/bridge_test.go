package parrotbridge

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"tonal.local/runtime/tonal"
)

type fakeClient struct {
	calls int
	last  Request
	fn    func(Request) (Result, error)
}

func (f *fakeClient) Invoke(_ context.Context, req Request) (Result, error) {
	f.calls++
	f.last = req
	return f.fn(req)
}

func TestRegistryPublishesExternalModelWithoutTlalocOntology(t *testing.T) {
	client := &fakeClient{fn: func(Request) (Result, error) {
		return Result{Output: json.RawMessage(`{"text":"ok","used_keys":["fact"]}`)}, nil
	}}
	registry, err := New(client, ResponseRendererSpec("conversation-model@r2"))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	candidates := registry.Candidates(tonal.ResponseRenderCapability, tonal.CapabilityGoal{})
	if len(candidates) != 1 {
		t.Fatalf("candidates=%d, want 1", len(candidates))
	}
	candidate := candidates[0]
	if candidate.Kind != tonal.CapabilityExternalModel || !candidate.Generative || candidate.Deterministic {
		t.Fatalf("candidate ontology = %+v", candidate)
	}
	if candidate.WorkerID != "default" {
		t.Fatalf("worker=%q", candidate.WorkerID)
	}
}

func TestRegistryDelegatesBoundedCapabilityToProviderNeutralClient(t *testing.T) {
	client := &fakeClient{fn: func(req Request) (Result, error) {
		if req.Capability != "PROPOSE_HYPOTHESIS" {
			t.Fatalf("capability=%q", req.Capability)
		}
		return Result{
			Output: json.RawMessage(`{"hypothesis":"x"}`),
			Observations: []tonal.Observation{{
				Producer: "parrot", Capability: req.Capability, Key: "hypothesis", Kind: "OBSERVATION", Value: json.RawMessage(`"x"`),
			}},
			Usage: &tonal.CapabilityUsage{ModelCalls: 1},
		}, nil
	}}
	registry, _ := New(client, CapabilitySpec{Capability: "PROPOSE_HYPOTHESIS", WorkerID: "reasoner"})
	result, err := registry.Execute(context.Background(), tonal.CapabilityExecutionRequest{
		Capability: "PROPOSE_HYPOTHESIS", WorkerID: "reasoner", Input: json.RawMessage(`{"failure":"boom"}`),
		PriorObservations: []tonal.Observation{{Key: "test.failure", Kind: "FACT", Value: json.RawMessage(`"boom"`)}},
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.calls != 1 || result.Usage == nil || result.Usage.ModelCalls != 1 {
		t.Fatalf("calls=%d result=%+v", client.calls, result)
	}
	if len(client.last.PriorObservations) != 1 || len(result.Observations) != 1 {
		t.Fatalf("bounded context/result lost: req=%+v result=%+v", client.last, result)
	}
}

func TestResponseComposerWorksAgainstParrotR2Source(t *testing.T) {
	client := &fakeClient{fn: func(req Request) (Result, error) {
		var envelope tonal.ResponseGroundingEnvelope
		if err := json.Unmarshal(req.Input, &envelope); err != nil {
			t.Fatalf("decode envelope: %v", err)
		}
		if envelope.SemanticAuthority != "BLACKBOARD_FACTS_ONLY" || envelope.MayAddClaims {
			t.Fatalf("bad semantic envelope: %+v", envelope)
		}
		return Result{Output: json.RawMessage(`{"text":"El valor es 512 MPa.","used_keys":["value","unit"]}`)}, nil
	}}
	parrot, _ := New(client, ResponseRendererSpec("model@r2"))
	composite, err := tonal.NewCompositeRegistry(tonal.RegistrySource{ID: "parrot", Registry: parrot})
	if err != nil {
		t.Fatalf("NewCompositeRegistry: %v", err)
	}
	bb := tonal.NewBlackboard("r")
	bb.Append(
		tonal.Observation{Key: "value", Kind: "FACT", Value: json.RawMessage(`512`)},
		tonal.Observation{Key: "unit", Kind: "FACT", Value: json.RawMessage(`"MPa"`)},
	)
	verifier := responseVerifierFunc(func(req tonal.ResponseVerificationRequest) (tonal.ResponseVerification, error) {
		return tonal.ResponseVerification{Verdict: tonal.ResponseVerificationPass}, nil
	})
	rendered, _, err := (tonal.ResponseComposer{Registry: composite, Verifier: verifier}).Compose(context.Background(), bb, tonal.ResponsePlan{
		GroundingKeys: []string{"value", "unit"},
		Style:         tonal.ResponseStyle{ToneMode: tonal.ResponseToneAuto, Language: "es"},
	})
	if err != nil {
		t.Fatalf("Compose: %v", err)
	}
	if rendered.Text != "El valor es 512 MPa." || client.calls != 1 {
		t.Fatalf("rendered=%+v calls=%d", rendered, client.calls)
	}
}

type responseVerifierFunc func(tonal.ResponseVerificationRequest) (tonal.ResponseVerification, error)

func (f responseVerifierFunc) VerifyResponse(_ context.Context, req tonal.ResponseVerificationRequest) (tonal.ResponseVerification, error) {
	return f(req)
}

func TestRegistryRejectsInvalidJSONAndUnknownCapabilities(t *testing.T) {
	client := &fakeClient{fn: func(Request) (Result, error) {
		return Result{Output: json.RawMessage(`not-json`)}, nil
	}}
	registry, _ := New(client, CapabilitySpec{Capability: "INTERPRET"})
	_, err := registry.Execute(context.Background(), tonal.CapabilityExecutionRequest{Capability: "INTERPRET", WorkerID: "default"})
	if err == nil || !strings.Contains(err.Error(), "invalid JSON") {
		t.Fatalf("expected JSON contract failure, got %v", err)
	}
	_, err = registry.Execute(context.Background(), tonal.CapabilityExecutionRequest{Capability: "OTHER", WorkerID: "default"})
	if err == nil || !strings.Contains(err.Error(), "not configured") {
		t.Fatalf("expected capability failure, got %v", err)
	}
}

func TestNewRejectsEmptyAndDuplicateCapabilitySpecs(t *testing.T) {
	client := &fakeClient{fn: func(Request) (Result, error) { return Result{Output: json.RawMessage(`{}`)}, nil }}
	if _, err := New(client); err == nil {
		t.Fatal("expected no-spec error")
	}
	if _, err := New(client, CapabilitySpec{Capability: "X"}, CapabilitySpec{Capability: "x"}); err == nil {
		t.Fatal("expected duplicate capability error")
	}
}
