package tonal

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

type responseRegistry struct {
	candidates []CapabilityCandidate
	execute    func(CapabilityExecutionRequest) (CapabilityExecutionResult, error)
	calls      int
	last       CapabilityExecutionRequest
}

func (r *responseRegistry) Candidates(capability string, _ CapabilityGoal) []CapabilityCandidate {
	if capability != ResponseRenderCapability {
		return nil
	}
	out := make([]CapabilityCandidate, len(r.candidates))
	copy(out, r.candidates)
	return out
}

func (r *responseRegistry) Execute(_ context.Context, req CapabilityExecutionRequest) (CapabilityExecutionResult, error) {
	r.calls++
	r.last = req
	return r.execute(req)
}

type responseVerifierFunc func(ResponseVerificationRequest) (ResponseVerification, error)

func (f responseVerifierFunc) VerifyResponse(_ context.Context, req ResponseVerificationRequest) (ResponseVerification, error) {
	return f(req)
}

func verifiedBlackboard() *Blackboard {
	bb := NewBlackboard("run-response")
	bb.Append(
		Observation{Producer: "verify", Capability: "VERIFY", Key: "answer.value", Kind: "FACT", Value: json.RawMessage(`512`)},
		Observation{Producer: "verify", Capability: "VERIFY", Key: "answer.unit", Kind: "FACT", Value: json.RawMessage(`"MPa"`)},
		Observation{Producer: "scratch", Capability: "INTERPRET", Key: "unverified.note", Kind: "OBSERVATION", Value: json.RawMessage(`"maybe 600"`)},
	)
	return bb
}

func rendererCandidate() CapabilityCandidate {
	return CapabilityCandidate{
		SourceID:   "parrot",
		WorkerID:   "default",
		Capability: ResponseRenderCapability,
		Kind:       CapabilityExternalModel,
		Generative: true,
		Selected:   true,
	}
}

func TestResponseComposer_AutoToneLeavesExpressionToParrotButGroundsContent(t *testing.T) {
	registry := &responseRegistry{candidates: []CapabilityCandidate{rendererCandidate()}}
	registry.execute = func(req CapabilityExecutionRequest) (CapabilityExecutionResult, error) {
		var envelope ResponseGroundingEnvelope
		if err := json.Unmarshal(req.Input, &envelope); err != nil {
			t.Fatalf("decode envelope: %v", err)
		}
		if envelope.SemanticAuthority != "BLACKBOARD_FACTS_ONLY" || envelope.MayAddClaims {
			t.Fatalf("semantic boundary = authority:%q may_add:%v", envelope.SemanticAuthority, envelope.MayAddClaims)
		}
		if len(envelope.Grounding) != 2 {
			t.Fatalf("grounding count = %d, want 2", len(envelope.Grounding))
		}
		if envelope.Style.ToneMode != ResponseToneAuto {
			t.Fatalf("tone mode = %q, want AUTO", envelope.Style.ToneMode)
		}
		if envelope.Style.Tone != "" {
			t.Fatalf("AUTO should not prescribe a tone, got %q", envelope.Style.Tone)
		}
		output, _ := json.Marshal(RenderedResponse{
			Text:     "El valor verificado es 512 MPa.",
			UsedKeys: []string{"answer.value", "answer.unit"},
		})
		return CapabilityExecutionResult{WorkerID: req.WorkerID, Output: output, Usage: &CapabilityUsage{ModelCalls: 1}}, nil
	}
	verifier := responseVerifierFunc(func(req ResponseVerificationRequest) (ResponseVerification, error) {
		if req.Rendered.Text == "" || len(req.Grounding) != 2 {
			t.Fatalf("verifier received incomplete response request")
		}
		return ResponseVerification{Verdict: ResponseVerificationPass, Reason: "claims supported by supplied facts"}, nil
	})

	rendered, verification, err := (ResponseComposer{Registry: registry, Verifier: verifier}).Compose(
		context.Background(),
		verifiedBlackboard(),
		ResponsePlan{GroundingKeys: []string{"answer.value", "answer.unit"}, Purpose: "answer the user clearly", Style: ResponseStyle{ToneMode: ResponseToneAuto, Language: "es"}},
	)
	if err != nil {
		t.Fatalf("Compose: %v", err)
	}
	if verification.Verdict != ResponseVerificationPass || rendered.Text == "" {
		t.Fatalf("unexpected result: rendered=%+v verification=%+v", rendered, verification)
	}
	if registry.last.SourceID != "parrot" || registry.last.WorkerID != "default" {
		t.Fatalf("dispatch = %s/%s", registry.last.SourceID, registry.last.WorkerID)
	}
	if len(registry.last.PriorObservations) != 2 {
		t.Fatalf("renderer should see only selected grounding facts, got %d observations", len(registry.last.PriorObservations))
	}
}

func TestResponseComposer_DefaultToneNormalizesToAuto(t *testing.T) {
	registry := &responseRegistry{candidates: []CapabilityCandidate{rendererCandidate()}}
	registry.execute = func(req CapabilityExecutionRequest) (CapabilityExecutionResult, error) {
		var envelope ResponseGroundingEnvelope
		_ = json.Unmarshal(req.Input, &envelope)
		if envelope.Style.ToneMode != ResponseToneAuto {
			t.Fatalf("default tone mode = %q, want AUTO", envelope.Style.ToneMode)
		}
		output, _ := json.Marshal(RenderedResponse{Text: "512", UsedKeys: []string{"answer.value"}})
		return CapabilityExecutionResult{Output: output}, nil
	}
	verifier := responseVerifierFunc(func(ResponseVerificationRequest) (ResponseVerification, error) {
		return ResponseVerification{Verdict: ResponseVerificationPass}, nil
	})
	_, _, err := (ResponseComposer{Registry: registry, Verifier: verifier}).Compose(context.Background(), verifiedBlackboard(), ResponsePlan{
		GroundingKeys: []string{"answer.value"},
	})
	if err != nil {
		t.Fatalf("Compose: %v", err)
	}
}

func TestResponseComposer_ExplicitToneControlsStyleNotFacts(t *testing.T) {
	registry := &responseRegistry{candidates: []CapabilityCandidate{rendererCandidate()}}
	registry.execute = func(req CapabilityExecutionRequest) (CapabilityExecutionResult, error) {
		var envelope ResponseGroundingEnvelope
		_ = json.Unmarshal(req.Input, &envelope)
		if envelope.Style.ToneMode != ResponseToneExplicit || envelope.Style.Tone != "technical" {
			t.Fatalf("style = %+v", envelope.Style)
		}
		output, _ := json.Marshal(RenderedResponse{Text: "Valor verificado: 512 MPa.", UsedKeys: []string{"answer.value", "answer.unit"}})
		return CapabilityExecutionResult{Output: output}, nil
	}
	verifier := responseVerifierFunc(func(ResponseVerificationRequest) (ResponseVerification, error) {
		return ResponseVerification{Verdict: ResponseVerificationPass}, nil
	})
	_, _, err := (ResponseComposer{Registry: registry, Verifier: verifier}).Compose(context.Background(), verifiedBlackboard(), ResponsePlan{
		GroundingKeys: []string{"answer.value", "answer.unit"},
		Style:         ResponseStyle{ToneMode: ResponseToneExplicit, Tone: "technical"},
	})
	if err != nil {
		t.Fatalf("Compose: %v", err)
	}
}

func TestResponseComposer_UnverifiedStateNeverReachesParrot(t *testing.T) {
	registry := &responseRegistry{candidates: []CapabilityCandidate{rendererCandidate()}, execute: func(CapabilityExecutionRequest) (CapabilityExecutionResult, error) {
		t.Fatal("renderer must not be called with unverified grounding")
		return CapabilityExecutionResult{}, nil
	}}
	verifier := responseVerifierFunc(func(ResponseVerificationRequest) (ResponseVerification, error) {
		return ResponseVerification{Verdict: ResponseVerificationPass}, nil
	})
	_, _, err := (ResponseComposer{Registry: registry, Verifier: verifier}).Compose(context.Background(), verifiedBlackboard(), ResponsePlan{
		GroundingKeys: []string{"unverified.note"},
		Style:         ResponseStyle{ToneMode: ResponseToneAuto},
	})
	if err == nil || !strings.Contains(err.Error(), "not a verified FACT") {
		t.Fatalf("expected verified-FACT gate, got %v", err)
	}
	if registry.calls != 0 {
		t.Fatalf("renderer calls = %d, want 0", registry.calls)
	}
}

func TestResponseComposer_RendererCannotWriteSemanticState(t *testing.T) {
	bb := verifiedBlackboard()
	before := len(bb.Snapshot())
	registry := &responseRegistry{candidates: []CapabilityCandidate{rendererCandidate()}}
	registry.execute = func(req CapabilityExecutionRequest) (CapabilityExecutionResult, error) {
		output, _ := json.Marshal(RenderedResponse{Text: "512", UsedKeys: []string{"answer.value"}})
		return CapabilityExecutionResult{
			Output: output,
			Observations: []Observation{{Producer: "parrot", Capability: ResponseRenderCapability, Key: "invented", Kind: "OBSERVATION", Value: json.RawMessage(`999`)}},
		}, nil
	}
	verifier := responseVerifierFunc(func(ResponseVerificationRequest) (ResponseVerification, error) {
		t.Fatal("semantic verifier must not run after renderer attempts a Blackboard write")
		return ResponseVerification{}, nil
	})
	_, _, err := (ResponseComposer{Registry: registry, Verifier: verifier}).Compose(context.Background(), bb, ResponsePlan{
		GroundingKeys: []string{"answer.value"},
		Style:         ResponseStyle{ToneMode: ResponseToneAuto},
	})
	if err == nil || !strings.Contains(err.Error(), "cannot mutate semantic state") {
		t.Fatalf("expected semantic-state rejection, got %v", err)
	}
	if after := len(bb.Snapshot()); after != before {
		t.Fatalf("Blackboard mutated: before=%d after=%d", before, after)
	}
}

func TestResponseComposer_RejectsUnsupportedKeyAndVerifierRejection(t *testing.T) {
	registry := &responseRegistry{candidates: []CapabilityCandidate{rendererCandidate()}}
	registry.execute = func(req CapabilityExecutionRequest) (CapabilityExecutionResult, error) {
		output, _ := json.Marshal(RenderedResponse{Text: "El valor es 999.", UsedKeys: []string{"not.allowed"}})
		return CapabilityExecutionResult{Output: output}, nil
	}
	verifier := responseVerifierFunc(func(ResponseVerificationRequest) (ResponseVerification, error) {
		return ResponseVerification{Verdict: ResponseVerificationPass}, nil
	})
	_, _, err := (ResponseComposer{Registry: registry, Verifier: verifier}).Compose(context.Background(), verifiedBlackboard(), ResponsePlan{
		GroundingKeys: []string{"answer.value"}, Style: ResponseStyle{ToneMode: ResponseToneAuto},
	})
	if err == nil || !strings.Contains(err.Error(), "non-grounding key") {
		t.Fatalf("expected key-closure rejection, got %v", err)
	}

	registry.execute = func(req CapabilityExecutionRequest) (CapabilityExecutionResult, error) {
		output, _ := json.Marshal(RenderedResponse{Text: "El valor es 999.", UsedKeys: []string{"answer.value"}})
		return CapabilityExecutionResult{Output: output}, nil
	}
	verifier = responseVerifierFunc(func(ResponseVerificationRequest) (ResponseVerification, error) {
		return ResponseVerification{Verdict: ResponseVerificationReject, Reason: "text adds a claim not entailed by answer.value"}, nil
	})
	rendered, verification, err := (ResponseComposer{Registry: registry, Verifier: verifier}).Compose(context.Background(), verifiedBlackboard(), ResponsePlan{
		GroundingKeys: []string{"answer.value"}, Style: ResponseStyle{ToneMode: ResponseToneAuto},
	})
	if err == nil || verification.Verdict != ResponseVerificationReject {
		t.Fatalf("expected semantic verifier rejection, rendered=%+v verification=%+v err=%v", rendered, verification, err)
	}
	if rendered.Text != "" {
		t.Fatalf("rejected response must not be released, got %q", rendered.Text)
	}
}

func TestResponseComposer_RequiresSemanticVerifier(t *testing.T) {
	registry := &responseRegistry{candidates: []CapabilityCandidate{rendererCandidate()}}
	_, _, err := (ResponseComposer{Registry: registry}).Compose(context.Background(), verifiedBlackboard(), ResponsePlan{
		GroundingKeys: []string{"answer.value"}, Style: ResponseStyle{ToneMode: ResponseToneAuto},
	})
	if err == nil || !strings.Contains(err.Error(), "ResponseVerifier is required") {
		t.Fatalf("expected fail-closed verifier requirement, got %v", err)
	}
}
