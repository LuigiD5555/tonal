package tonal

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

const ResponseRenderCapability = "RENDER_RESPONSE"

// ResponseToneMode separates semantic authority from expressive freedom.
// AUTO means the renderer may choose a natural tone from the interaction and
// its pretrained conversational priors. EXPLICIT means Tonal/user policy has
// requested a particular tone. Neither mode grants freedom to add claims.
type ResponseToneMode string

const (
	ResponseToneAuto     ResponseToneMode = "AUTO"
	ResponseToneExplicit ResponseToneMode = "EXPLICIT"
)

// ResponseStyle controls presentation only. It must never carry factual
// content that bypasses the Blackboard.
type ResponseStyle struct {
	ToneMode ResponseToneMode `json:"tone_mode"`
	Tone     string           `json:"tone,omitempty"`
	Language string           `json:"language,omitempty"`
	Audience string           `json:"audience,omitempty"`
}

// ResponsePlan describes how Tonal wants already-committed semantic state to
// be rendered. GroundingKeys are the only Blackboard entries from which the
// renderer may derive factual claims.
type ResponsePlan struct {
	GroundingKeys []string      `json:"grounding_keys"`
	Purpose       string        `json:"purpose,omitempty"`
	Style         ResponseStyle `json:"style"`
}

// ResponseGroundingEnvelope is the read-only semantic boundary sent to the
// renderer. The renderer may choose wording, ordering and tone, but the
// Blackboard facts remain the sole semantic authority.
type ResponseGroundingEnvelope struct {
	SemanticAuthority string        `json:"semantic_authority"`
	MayAddClaims      bool          `json:"may_add_claims"`
	Purpose           string        `json:"purpose,omitempty"`
	Style             ResponseStyle `json:"style"`
	Grounding         []Observation `json:"grounding"`
}

// RenderedResponse is the renderer's presentation-only result. UsedKeys must
// identify which allowed Blackboard facts support the response text.
type RenderedResponse struct {
	Text     string   `json:"text"`
	UsedKeys []string `json:"used_keys"`
}

type ResponseVerificationVerdict string

const (
	ResponseVerificationPass   ResponseVerificationVerdict = "PASS"
	ResponseVerificationReject ResponseVerificationVerdict = "REJECT"
)

type ResponseVerificationRequest struct {
	Plan      ResponsePlan              `json:"plan"`
	Grounding []Observation             `json:"grounding"`
	Candidate CapabilityCandidate       `json:"candidate"`
	Rendered  RenderedResponse          `json:"rendered"`
	RawResult CapabilityExecutionResult `json:"raw_result"`
}

type ResponseVerification struct {
	Verdict ResponseVerificationVerdict `json:"verdict"`
	Reason  string                      `json:"reason,omitempty"`
}

// ResponseVerifier is the semantic closure gate between Parrot and the user.
// A renderer result is never exposed unless a verifier says it is supported
// by the supplied Blackboard facts. Implementations may be deterministic,
// Tlaloque-backed, or composed, but they must not mutate the Blackboard.
type ResponseVerifier interface {
	VerifyResponse(ctx context.Context, req ResponseVerificationRequest) (ResponseVerification, error)
}

// ResponseComposer owns the final presentation seam. Parrot/external
// cognition is allowed expressive freedom here, but not semantic authority:
// it receives only a projection of verified Blackboard facts, cannot append
// observations, and its text must pass ResponseVerifier before release.
type ResponseComposer struct {
	Registry CapabilityRegistry
	Verifier ResponseVerifier
}

func (c ResponseComposer) Compose(ctx context.Context, blackboard *Blackboard, plan ResponsePlan) (RenderedResponse, ResponseVerification, error) {
	if c.Registry == nil {
		return RenderedResponse{}, ResponseVerification{}, fmt.Errorf("response composer: Registry is required")
	}
	if c.Verifier == nil {
		return RenderedResponse{}, ResponseVerification{}, fmt.Errorf("response composer: ResponseVerifier is required")
	}
	if blackboard == nil {
		return RenderedResponse{}, ResponseVerification{}, fmt.Errorf("response composer: Blackboard is required")
	}

	style, err := normalizeResponseStyle(plan.Style)
	if err != nil {
		return RenderedResponse{}, ResponseVerification{}, err
	}
	plan.Style = style

	grounding, err := responseGrounding(blackboard, plan.GroundingKeys)
	if err != nil {
		return RenderedResponse{}, ResponseVerification{}, err
	}

	envelope := ResponseGroundingEnvelope{
		SemanticAuthority: "BLACKBOARD_FACTS_ONLY",
		MayAddClaims:      false,
		Purpose:           strings.TrimSpace(plan.Purpose),
		Style:             style,
		Grounding:         grounding,
	}
	input, err := json.Marshal(envelope)
	if err != nil {
		return RenderedResponse{}, ResponseVerification{}, fmt.Errorf("response composer: encode grounding envelope: %w", err)
	}

	candidates := c.Registry.Candidates(ResponseRenderCapability, CapabilityGoal{Capability: ResponseRenderCapability})
	candidate, ok := selectExternalRenderer(candidates)
	if !ok {
		return RenderedResponse{}, ResponseVerification{}, fmt.Errorf("response composer: no EXTERNAL_MODEL renderer is available")
	}

	result, err := c.Registry.Execute(ctx, CapabilityExecutionRequest{
		TaskID:            blackboard.TaskID(),
		NodeID:            blackboard.TaskID() + "::response-render",
		Capability:        ResponseRenderCapability,
		SourceID:          candidate.SourceID,
		WorkerID:          candidate.WorkerID,
		Input:             input,
		PriorObservations: append([]Observation(nil), grounding...),
	})
	if err != nil {
		return RenderedResponse{}, ResponseVerification{}, fmt.Errorf("response composer: render: %w", err)
	}

	// Rendering is presentation-only. Any attempt to emit semantic
	// observations would create a second path into committed state and is
	// therefore rejected before verification.
	if len(result.Observations) != 0 {
		return RenderedResponse{}, ResponseVerification{}, fmt.Errorf("response composer: renderer attempted to emit %d observation(s); rendering cannot mutate semantic state", len(result.Observations))
	}

	var rendered RenderedResponse
	if err := json.Unmarshal(result.Output, &rendered); err != nil {
		return RenderedResponse{}, ResponseVerification{}, fmt.Errorf("response composer: renderer output must be a RenderedResponse JSON envelope: %w", err)
	}
	if strings.TrimSpace(rendered.Text) == "" {
		return RenderedResponse{}, ResponseVerification{}, fmt.Errorf("response composer: renderer returned empty text")
	}
	if len(rendered.UsedKeys) == 0 {
		return RenderedResponse{}, ResponseVerification{}, fmt.Errorf("response composer: renderer must cite at least one grounding key")
	}
	if err := validateUsedKeys(rendered.UsedKeys, grounding); err != nil {
		return RenderedResponse{}, ResponseVerification{}, err
	}

	verification, err := c.Verifier.VerifyResponse(ctx, ResponseVerificationRequest{
		Plan:      plan,
		Grounding: grounding,
		Candidate: candidate,
		Rendered:  rendered,
		RawResult: result,
	})
	if err != nil {
		return RenderedResponse{}, ResponseVerification{}, fmt.Errorf("response composer: verify response: %w", err)
	}
	if verification.Verdict != ResponseVerificationPass {
		return RenderedResponse{}, verification, fmt.Errorf("response composer: response rejected: %s", verification.Reason)
	}
	return rendered, verification, nil
}

func normalizeResponseStyle(style ResponseStyle) (ResponseStyle, error) {
	if style.ToneMode == "" {
		style.ToneMode = ResponseToneAuto
	}
	switch style.ToneMode {
	case ResponseToneAuto:
		// AUTO deliberately leaves tone selection to the pretrained renderer.
		// Ignore any stale explicit tone so presentation policy is unambiguous.
		style.Tone = ""
		return style, nil
	case ResponseToneExplicit:
		style.Tone = strings.TrimSpace(style.Tone)
		if style.Tone == "" {
			return ResponseStyle{}, fmt.Errorf("response composer: EXPLICIT tone mode requires a tone")
		}
		return style, nil
	default:
		return ResponseStyle{}, fmt.Errorf("response composer: unknown tone mode %q", style.ToneMode)
	}
}

func responseGrounding(blackboard *Blackboard, keys []string) ([]Observation, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("response composer: at least one grounding key is required")
	}
	seen := map[string]struct{}{}
	grounding := make([]Observation, 0, len(keys))
	for _, rawKey := range keys {
		key := strings.TrimSpace(rawKey)
		if key == "" {
			return nil, fmt.Errorf("response composer: grounding key cannot be empty")
		}
		if _, duplicate := seen[key]; duplicate {
			continue
		}
		seen[key] = struct{}{}
		entry, ok := blackboard.Latest(key)
		if !ok {
			return nil, fmt.Errorf("response composer: grounding key %q is not present in Blackboard", key)
		}
		if !strings.EqualFold(entry.Kind, "FACT") {
			return nil, fmt.Errorf("response composer: grounding key %q is not a verified FACT", key)
		}
		grounding = append(grounding, entry)
	}
	return grounding, nil
}

func selectExternalRenderer(candidates []CapabilityCandidate) (CapabilityCandidate, bool) {
	var first *CapabilityCandidate
	for _, candidate := range candidates {
		if candidate.Kind != CapabilityExternalModel {
			continue
		}
		copyCandidate := candidate
		if first == nil {
			first = &copyCandidate
		}
		if candidate.Selected {
			return candidate, true
		}
	}
	if first != nil {
		return *first, true
	}
	return CapabilityCandidate{}, false
}

func validateUsedKeys(used []string, grounding []Observation) error {
	allowed := make(map[string]struct{}, len(grounding))
	for _, entry := range grounding {
		allowed[entry.Key] = struct{}{}
	}
	for _, rawKey := range used {
		key := strings.TrimSpace(rawKey)
		if _, ok := allowed[key]; !ok {
			return fmt.Errorf("response composer: renderer cited non-grounding key %q", key)
		}
	}
	return nil
}
