package tonal

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"tlaloc.local/behaviorlab/tlaloquekit"
)

// scalarString renders a JSON scalar the way the deterministic Tlaloque
// operand contracts expect it: an integral float without a trailing ".0",
// a fractional float at full precision, bools and strings verbatim.
func scalarString(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case float64:
		return strconv.FormatFloat(typed, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(typed)
	case nil:
		return ""
	default:
		return fmt.Sprint(typed)
	}
}

// Instance is one concrete workflow to run. Params carries only
// runtime-visible values (a question string, a store directory, a page
// image path, a threshold). It must never contain an expected answer or a
// hidden evidence address.
type Instance struct {
	ID            string            `json:"id"`
	Family        string            `json:"family"`
	DeclaredDepth int               `json:"declared_depth"`
	Params        map[string]string `json:"params"`
}

// Engine is the TONAL scheduler/execution loop. It owns the Blackboard for
// each run, resolves each step's executor through the Registry (optionally
// nudged by the arm's RoutingPolicy), executes in dependency order, and
// emits the deterministic RunRecord.
type Engine struct {
	Registry tlaloquekit.QualifiedRegistry
	// Now is injectable for deterministic tests; defaults to time.Now.
	Now func() time.Time
}

func (e *Engine) now() time.Time {
	if e.Now != nil {
		return e.Now()
	}
	return time.Now()
}

var placeholderPattern = regexp.MustCompile(`\$\{(param|obs|node):([a-zA-Z0-9_.-]+)(?::([a-zA-Z0-9_.-]+))?\}`)

// RunWorkflow executes one workflow instance under one routing policy and
// returns the trace plus the Blackboard it produced.
func (e *Engine) RunWorkflow(ctx context.Context, family TaskFamily, instance Instance, policy RoutingPolicy) (RunRecord, *Blackboard, error) {
	family, err := family.Normalize()
	if err != nil {
		return RunRecord{}, nil, err
	}
	blackboard := NewBlackboard(instance.ID)
	record := RunRecord{
		WorkflowID:         instance.ID,
		Family:             family.ID,
		Goal:               family.Goal,
		Arm:                policy.Name(),
		DeclaredDepth:      instance.DeclaredDepth,
		CriticalPathDepth:  family.CriticalPathDepth(),
		TerminalOutputKind: "evaluable_terminal_output",
	}
	if family.HasVerify() {
		record.TerminalOutputKind = "promoted_fact"
	}

	nodeIDByLocal := map[string]string{}
	obsByLocal := map[string]tlaloquekit.Observation{}

	for _, step := range family.topoOrder() {
		nodeID := instance.ID + "::" + step.LocalID
		nodeIDByLocal[step.LocalID] = nodeID

		stepTrace := StepTrace{
			LocalID:    step.LocalID,
			Role:       step.Role,
			Capability: step.Capability,
			NodeID:     nodeID,
		}

		// --- resolution / routing ---
		goal := tlaloquekit.Goal{Capability: step.Capability, PreferDeterministic: step.PreferDeterministic}
		candidates := e.Registry.Candidates(step.Capability, goal)
		stepTrace.Candidates = candidates
		if len(candidates) == 0 {
			stepTrace.Error = "CAPABILITY_UNAVAILABLE: no qualified executor for " + step.Capability
			record.Steps = append(record.Steps, stepTrace)
			record.Accounting.recordStep(stepTrace, false)
			return finish(record, blackboard, "CONTRACT_FAILURE", stepTrace.Error), blackboard, nil
		}

		workerID, reason := policy.SelectWorker(step, candidates)
		if workerID == "" {
			for _, candidate := range candidates {
				if candidate.Selected {
					workerID = candidate.Descriptor.ID
					if reason == "" {
						reason = candidate.Reason
					}
				}
			}
		}
		if workerID == "" {
			workerID = candidates[0].Descriptor.ID
			reason = "fell back to first candidate: " + candidates[0].Reason
		}
		stepTrace.SelectedWorker = workerID
		stepTrace.SelectionReason = reason
		for _, candidate := range candidates {
			if candidate.Descriptor.ID == workerID {
				stepTrace.EngineKind = string(candidate.Descriptor.Engine)
				stepTrace.ProfileVersion = candidate.Descriptor.ProfileRef
			}
		}

		// --- input construction (generic template, no per-instance plan) ---
		input, reads, err := e.buildInput(step, instance, nodeIDByLocal, obsByLocal)
		if err != nil {
			stepTrace.Error = "INPUT_BUILD_FAILURE: " + err.Error()
			record.Steps = append(record.Steps, stepTrace)
			record.Accounting.recordStep(stepTrace, false)
			return finish(record, blackboard, "CONTRACT_FAILURE", stepTrace.Error), blackboard, nil
		}
		stepTrace.BlackboardReads = reads
		stepTrace.InputJSON = string(input)

		// --- execution ---
		start := e.now()
		result, execErr := e.Registry.Execute(ctx, tlaloquekit.ExecutionRequest{
			TaskID:            instance.ID,
			NodeID:            nodeID,
			Capability:        step.Capability,
			WorkerID:          workerID,
			Input:             input,
			PriorObservations: blackboard.Snapshot(),
		})
		stepTrace.LatencyMS = e.now().Sub(start).Milliseconds()
		if execErr != nil {
			stepTrace.Error = "EXECUTION_ERROR: " + execErr.Error()
			record.Steps = append(record.Steps, stepTrace)
			record.Accounting.recordStep(stepTrace, false)
			return finish(record, blackboard, "ERROR", stepTrace.Error), blackboard, nil
		}
		stepTrace.OutputJSON = string(result.Output)
		stepTrace.Notes = result.Notes
		if result.Usage != nil {
			stepTrace.ModelCalls = result.Usage.ModelCalls
		}
		for _, obs := range result.Observations {
			// Fact-promotion scope invariant (protocol section 4): only a
			// VERIFY Tlaloque may write a FACT. Any other node emitting one
			// is a hard protocol error.
			if strings.EqualFold(obs.Kind, "FACT") && !strings.EqualFold(step.Capability, "VERIFY") {
				stepTrace.Error = "FACT_PROMOTION_SCOPE_VIOLATION: " + step.Capability + " node " + step.LocalID + " emitted a FACT"
				record.Steps = append(record.Steps, stepTrace)
				record.Accounting.recordStep(stepTrace, false)
				return finish(record, blackboard, "CONTRACT_FAILURE", stepTrace.Error), blackboard, nil
			}
			blackboard.Append(obs)
			stepTrace.BlackboardWrites = append(stepTrace.BlackboardWrites, obs.Key)
			if obs.Key == nodeID || strings.HasSuffix(obs.Key, step.LocalID) {
				obsByLocal[step.LocalID] = obs
			}
		}
		if _, ok := obsByLocal[step.LocalID]; !ok && len(result.Observations) > 0 {
			obsByLocal[step.LocalID] = result.Observations[len(result.Observations)-1]
		}

		record.Steps = append(record.Steps, stepTrace)
		record.Accounting.recordStep(stepTrace, true)
	}

	// The final step's observation is the workflow answer.
	last := family.Steps[len(family.Steps)-1]
	final, ok := obsByLocal[last.LocalID]
	if !ok {
		return finish(record, blackboard, "CONTRACT_FAILURE", "final step produced no observation"), blackboard, nil
	}
	record.FinalKey = final.Key
	record.FinalValue = final
	// A Tlaloque signals abstention either through the observation Status
	// (a verified/unsupported fact) or through provenance["status"] (a
	// generative Tlaloque that rejected the call inside its own envelope).
	declaredStatus := final.Status
	if declaredStatus == "" {
		declaredStatus = final.Provenance["status"]
	}
	status := "OK"
	switch strings.ToUpper(declaredStatus) {
	case "UNKNOWN":
		status = "UNKNOWN"
	case "UNSUPPORTED":
		status = "UNSUPPORTED"
	}
	return finish(record, blackboard, status, ""), blackboard, nil
}

func finish(record RunRecord, _ *Blackboard, status, errMessage string) RunRecord {
	record.FinalStatus = status
	record.Error = errMessage
	return record
}

func (e *Engine) buildInput(step Step, instance Instance, nodeIDByLocal map[string]string, obsByLocal map[string]tlaloquekit.Observation) (json.RawMessage, []string, error) {
	reads := map[string]struct{}{}
	resolved := map[string]any{}
	for key, raw := range step.Input.Template {
		value, err := e.resolveValue(raw, instance, nodeIDByLocal, obsByLocal, reads)
		if err != nil {
			return nil, nil, fmt.Errorf("field %q: %w", key, err)
		}
		resolved[key] = value
	}
	body, err := json.Marshal(resolved)
	if err != nil {
		return nil, nil, err
	}
	readList := make([]string, 0, len(reads))
	for key := range reads {
		readList = append(readList, key)
	}
	return body, readList, nil
}

func (e *Engine) resolveValue(raw any, instance Instance, nodeIDByLocal map[string]string, obsByLocal map[string]tlaloquekit.Observation, reads map[string]struct{}) (any, error) {
	switch typed := raw.(type) {
	case map[string]any:
		out := make(map[string]any, len(typed))
		for key, value := range typed {
			resolved, err := e.resolveValue(value, instance, nodeIDByLocal, obsByLocal, reads)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", key, err)
			}
			out[key] = resolved
		}
		return out, nil
	case []any:
		out := make([]any, len(typed))
		for index, value := range typed {
			resolved, err := e.resolveValue(value, instance, nodeIDByLocal, obsByLocal, reads)
			if err != nil {
				return nil, err
			}
			out[index] = resolved
		}
		return out, nil
	}
	text, ok := raw.(string)
	if !ok {
		return raw, nil
	}
	match := placeholderPattern.FindStringSubmatch(text)
	// A leaf that is exactly one placeholder: an object/array value passes
	// through intact (a bbox, a geometry), while a scalar is rendered as a
	// string so it fits the string-typed operand contracts the
	// deterministic Tlaloques declare.
	if match != nil && match[0] == text {
		value, err := e.resolvePlaceholder(match[1], match[2], match[3], instance, nodeIDByLocal, obsByLocal, reads)
		if err != nil {
			return nil, err
		}
		switch value.(type) {
		case map[string]any, []any, nil:
			return value, nil
		default:
			return scalarString(value), nil
		}
	}
	// Otherwise substitute every placeholder as text.
	var subErr error
	out := placeholderPattern.ReplaceAllStringFunc(text, func(token string) string {
		parts := placeholderPattern.FindStringSubmatch(token)
		value, err := e.resolvePlaceholder(parts[1], parts[2], parts[3], instance, nodeIDByLocal, obsByLocal, reads)
		if err != nil {
			subErr = err
			return ""
		}
		return scalarString(value)
	})
	if subErr != nil {
		return nil, subErr
	}
	return out, nil
}

func (e *Engine) resolvePlaceholder(kind, name, field string, instance Instance, nodeIDByLocal map[string]string, obsByLocal map[string]tlaloquekit.Observation, reads map[string]struct{}) (any, error) {
	switch kind {
	case "param":
		value, ok := instance.Params[name]
		if !ok {
			return nil, fmt.Errorf("instance param %q is not set", name)
		}
		return value, nil
	case "node":
		nodeID, ok := nodeIDByLocal[name]
		if !ok {
			return nil, fmt.Errorf("step %q has not run yet", name)
		}
		return nodeID, nil
	case "obs":
		obs, ok := obsByLocal[name]
		if !ok {
			return nil, fmt.Errorf("step %q produced no observation", name)
		}
		reads[obs.Key] = struct{}{}
		var decoded any
		if err := json.Unmarshal(obs.Value, &decoded); err != nil {
			return nil, fmt.Errorf("step %q observation is not JSON: %w", name, err)
		}
		if field == "" {
			return decoded, nil
		}
		object, ok := decoded.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("step %q observation is not an object; cannot read field %q", name, field)
		}
		fieldValue, ok := object[field]
		if !ok {
			return nil, fmt.Errorf("step %q observation has no field %q", name, field)
		}
		return fieldValue, nil
	default:
		return nil, fmt.Errorf("unknown placeholder kind %q", kind)
	}
}
