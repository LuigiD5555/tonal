package tonal

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// ControlAction is the controller's public decision for one iteration.
type ControlAction string

const (
	ControlExecute ControlAction = "EXECUTE"
	ControlSucceed ControlAction = "SUCCEED"
	ControlStop    ControlAction = "STOP"
)

// ControlDecision says what bounded behavior should happen next. It does not
// name an executor or provider. AllowExternalModel is deliberately explicit:
// false means Parrot/external cognition is ineligible for this transition.
type ControlDecision struct {
	Action              ControlAction   `json:"action"`
	Capability          string          `json:"capability,omitempty"`
	Input               json.RawMessage `json:"input,omitempty"`
	PreferDeterministic bool            `json:"prefer_deterministic,omitempty"`
	AllowExternalModel  bool            `json:"allow_external_model,omitempty"`
	ExpectedEffect      string          `json:"expected_effect,omitempty"`
	Reason              string          `json:"reason,omitempty"`
}

// ControlContext is the immutable information a Controller sees when deciding
// the next bounded operation. It exposes committed state and observable
// transition results, never a mutable Blackboard pointer.
type ControlContext struct {
	RunID                string             `json:"run_id"`
	Goal                 string             `json:"goal"`
	Iteration            int                `json:"iteration"`
	RemainingTransitions int                `json:"remaining_transitions"`
	Observations         []Observation      `json:"observations,omitempty"`
	LastTransition       *ControlTransition `json:"last_transition,omitempty"`
}

// Controller decides WHAT capability is needed next. It is not the router;
// SelectionPolicy independently decides WHICH eligible component executes it.
type Controller interface {
	Next(ctx context.Context, state ControlContext) (ControlDecision, error)
}

// FixedProgramController is the deterministic R0/ceiling controller. It is
// intentionally boring: it emits a predeclared sequence, then succeeds.
type FixedProgramController struct {
	Decisions []ControlDecision
}

func (p FixedProgramController) Next(_ context.Context, state ControlContext) (ControlDecision, error) {
	if state.Iteration >= len(p.Decisions) {
		return ControlDecision{Action: ControlSucceed, Reason: "fixed program exhausted"}, nil
	}
	return p.Decisions[state.Iteration], nil
}

// VerificationVerdict controls whether an execution result mutates committed
// Blackboard state.
type VerificationVerdict string

const (
	VerificationPass    VerificationVerdict = "PASS"
	VerificationReject  VerificationVerdict = "REJECT"
	VerificationUnknown VerificationVerdict = "UNKNOWN"
)

type TransitionVerificationRequest struct {
	Decision ControlDecision           `json:"decision"`
	Selected CapabilityCandidate       `json:"selected"`
	Result   CapabilityExecutionResult `json:"result"`
}

type TransitionVerification struct {
	Verdict VerificationVerdict `json:"verdict"`
	Reason  string              `json:"reason,omitempty"`
}

// TransitionVerifier decides whether observations may enter committed state.
// It does not execute the capability and does not select the component.
type TransitionVerifier interface {
	Verify(ctx context.Context, req TransitionVerificationRequest) (TransitionVerification, error)
}

// ContractTransitionVerifier is R0's deterministic safety gate. Domain-level
// verification remains a separate capability/experiment, but these invariants
// hold for every control transition.
type ContractTransitionVerifier struct{}

func (ContractTransitionVerifier) Verify(_ context.Context, req TransitionVerificationRequest) (TransitionVerification, error) {
	if len(req.Result.Observations) == 0 {
		return TransitionVerification{Verdict: VerificationReject, Reason: "execution produced no observation"}, nil
	}
	for _, obs := range req.Result.Observations {
		if strings.EqualFold(obs.Kind, "FACT") && !strings.EqualFold(req.Decision.Capability, "VERIFY") {
			return TransitionVerification{
				Verdict: VerificationReject,
				Reason:  "FACT_PROMOTION_SCOPE_VIOLATION: only VERIFY may promote FACT",
			}, nil
		}
	}
	return TransitionVerification{Verdict: VerificationPass, Reason: "contract invariants satisfied"}, nil
}

type ControlTransitionOutcome string

const (
	TransitionCommitted      ControlTransitionOutcome = "COMMITTED"
	TransitionRejected       ControlTransitionOutcome = "REJECTED"
	TransitionUnavailable    ControlTransitionOutcome = "UNAVAILABLE"
	TransitionExecutionError ControlTransitionOutcome = "EXECUTION_ERROR"
	TransitionVerifierError  ControlTransitionOutcome = "VERIFIER_ERROR"
)

// ControlTransition is the observable flight-recorder entry for one loop
// iteration. It is designed to be convertible into an Episode later.
type ControlTransition struct {
	Iteration    int                        `json:"iteration"`
	Decision     ControlDecision            `json:"decision"`
	Candidates   []CapabilityCandidate      `json:"candidates,omitempty"`
	Selected     *CapabilityCandidate       `json:"selected,omitempty"`
	Result       *CapabilityExecutionResult `json:"result,omitempty"`
	Verification *TransitionVerification    `json:"verification,omitempty"`
	Outcome      ControlTransitionOutcome   `json:"outcome"`
	Error        string                     `json:"error,omitempty"`
}

type ControlRunStatus string

const (
	ControlRunSucceeded       ControlRunStatus = "SUCCEEDED"
	ControlRunStopped         ControlRunStatus = "STOPPED"
	ControlRunBudgetExhausted ControlRunStatus = "BUDGET_EXHAUSTED"
	ControlRunControllerError ControlRunStatus = "CONTROLLER_ERROR"
)

type ControlAccounting struct {
	Transitions        int `json:"transitions"`
	Committed          int `json:"committed"`
	Rejected           int `json:"rejected"`
	Unavailable        int `json:"unavailable"`
	ExecutionErrors    int `json:"execution_errors"`
	ModelCalls         int `json:"model_calls"`
	ExternalModelCalls int `json:"external_model_calls"`
}

type ControlRun struct {
	RunID        string              `json:"run_id"`
	Goal         string              `json:"goal"`
	Status       ControlRunStatus    `json:"status"`
	Transitions  []ControlTransition `json:"transitions"`
	Observations []Observation       `json:"observations"`
	Accounting   ControlAccounting   `json:"accounting"`
	Error        string              `json:"error,omitempty"`
}

type ControlRunRequest struct {
	RunID string `json:"run_id"`
	Goal  string `json:"goal"`
}

// ControlLoop repeatedly asks a Controller what to do next, then resolves,
// executes, verifies and conditionally commits one bounded transition.
type ControlLoop struct {
	Registry       CapabilityRegistry
	Selection      SelectionPolicy
	Verifier       TransitionVerifier
	MaxTransitions int
}

func (l *ControlLoop) Run(ctx context.Context, req ControlRunRequest, controller Controller) (ControlRun, error) {
	if l.Registry == nil {
		return ControlRun{}, fmt.Errorf("control loop: Registry is required")
	}
	if controller == nil {
		return ControlRun{}, fmt.Errorf("control loop: Controller is required")
	}
	selection := l.Selection
	if selection == nil {
		selection = MachineryFirstPolicy{}
	}
	verifier := l.Verifier
	if verifier == nil {
		verifier = ContractTransitionVerifier{}
	}
	maxTransitions := l.MaxTransitions
	if maxTransitions <= 0 {
		maxTransitions = 32
	}

	blackboard := NewBlackboard(req.RunID)
	run := ControlRun{RunID: req.RunID, Goal: req.Goal}

	for iteration := 0; iteration < maxTransitions; iteration++ {
		var last *ControlTransition
		if len(run.Transitions) > 0 {
			copyOfLast := run.Transitions[len(run.Transitions)-1]
			last = &copyOfLast
		}
		decision, err := controller.Next(ctx, ControlContext{
			RunID:                req.RunID,
			Goal:                 req.Goal,
			Iteration:            iteration,
			RemainingTransitions: maxTransitions - iteration,
			Observations:         blackboard.Snapshot(),
			LastTransition:       last,
		})
		if err != nil {
			run.Status = ControlRunControllerError
			run.Error = err.Error()
			run.Observations = blackboard.Snapshot()
			return run, nil
		}

		switch decision.Action {
		case ControlSucceed:
			run.Status = ControlRunSucceeded
			run.Observations = blackboard.Snapshot()
			return run, nil
		case ControlStop:
			run.Status = ControlRunStopped
			run.Observations = blackboard.Snapshot()
			return run, nil
		case ControlExecute:
			// continue below
		default:
			run.Status = ControlRunControllerError
			run.Error = fmt.Sprintf("unknown controller action %q", decision.Action)
			run.Observations = blackboard.Snapshot()
			return run, nil
		}

		capability := strings.TrimSpace(decision.Capability)
		if capability == "" {
			run.Status = ControlRunControllerError
			run.Error = "EXECUTE decision requires capability"
			run.Observations = blackboard.Snapshot()
			return run, nil
		}
		decision.Capability = capability

		allCandidates := l.Registry.Candidates(capability, CapabilityGoal{
			Capability:          capability,
			PreferDeterministic: decision.PreferDeterministic,
		})
		eligible := filterControlCandidates(allCandidates, decision.AllowExternalModel)
		transition := ControlTransition{
			Iteration:  iteration,
			Decision:   decision,
			Candidates: eligible,
		}

		if len(eligible) == 0 {
			transition.Outcome = TransitionUnavailable
			if !decision.AllowExternalModel && hasExternalCandidate(allCandidates) {
				transition.Error = "machinery unavailable; external cognition exists but was not explicitly permitted"
			} else {
				transition.Error = "no eligible capability candidate"
			}
			run.Transitions = append(run.Transitions, transition)
			run.Accounting.Transitions++
			run.Accounting.Unavailable++
			continue
		}

		selected, reason, ok := chooseControlCandidate(selection, decision, eligible)
		if !ok {
			transition.Outcome = TransitionUnavailable
			transition.Error = "selection policy returned no eligible component"
			run.Transitions = append(run.Transitions, transition)
			run.Accounting.Transitions++
			run.Accounting.Unavailable++
			continue
		}
		selected.Reason = reason
		transition.Selected = &selected

		result, execErr := l.Registry.Execute(ctx, CapabilityExecutionRequest{
			TaskID:            req.RunID,
			NodeID:            fmt.Sprintf("%s::control-%03d", req.RunID, iteration),
			Capability:        capability,
			SourceID:          selected.SourceID,
			WorkerID:          selected.WorkerID,
			Input:             decision.Input,
			PriorObservations: blackboard.Snapshot(),
		})
		if execErr != nil {
			transition.Outcome = TransitionExecutionError
			transition.Error = execErr.Error()
			run.Transitions = append(run.Transitions, transition)
			run.Accounting.Transitions++
			run.Accounting.ExecutionErrors++
			continue
		}
		transition.Result = &result
		if result.Usage != nil {
			run.Accounting.ModelCalls += result.Usage.ModelCalls
			if selected.Kind == CapabilityExternalModel {
				run.Accounting.ExternalModelCalls += result.Usage.ModelCalls
			}
		}

		verification, verifyErr := verifier.Verify(ctx, TransitionVerificationRequest{
			Decision: decision,
			Selected: selected,
			Result:   result,
		})
		if verifyErr != nil {
			transition.Outcome = TransitionVerifierError
			transition.Error = verifyErr.Error()
			run.Transitions = append(run.Transitions, transition)
			run.Accounting.Transitions++
			run.Accounting.Rejected++
			continue
		}
		transition.Verification = &verification
		if verification.Verdict != VerificationPass {
			transition.Outcome = TransitionRejected
			run.Transitions = append(run.Transitions, transition)
			run.Accounting.Transitions++
			run.Accounting.Rejected++
			continue
		}

		blackboard.Append(result.Observations...)
		transition.Outcome = TransitionCommitted
		run.Transitions = append(run.Transitions, transition)
		run.Accounting.Transitions++
		run.Accounting.Committed++
	}

	run.Status = ControlRunBudgetExhausted
	run.Observations = blackboard.Snapshot()
	return run, nil
}

func filterControlCandidates(candidates []CapabilityCandidate, allowExternal bool) []CapabilityCandidate {
	out := make([]CapabilityCandidate, 0, len(candidates))
	for _, candidate := range candidates {
		if candidate.Kind == CapabilityExternalModel && !allowExternal {
			continue
		}
		out = append(out, candidate)
	}
	return out
}

func hasExternalCandidate(candidates []CapabilityCandidate) bool {
	for _, candidate := range candidates {
		if candidate.Kind == CapabilityExternalModel {
			return true
		}
	}
	return false
}

func chooseControlCandidate(policy SelectionPolicy, decision ControlDecision, candidates []CapabilityCandidate) (CapabilityCandidate, string, bool) {
	return selectCapabilityCandidate(policy, Step{
		Capability:          decision.Capability,
		PreferDeterministic: decision.PreferDeterministic,
	}, candidates)
}
