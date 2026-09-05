# Tonal Control Loop R0

**Status:** R2 foundation prototype specification  
**Scope:** post-T1 architecture; does not modify frozen T1

## Purpose

Control Loop R0 is the smallest mechanism that lets Tonal advance a goal through repeated bounded, verified transitions instead of requiring one model call to own the entire task.

```text
GOAL
  ↓
STATE / BLACKBOARD
  ↓
CONTROLLER: what operation is needed next?
  ↓
CAPABILITY REGISTRY
  ↓
SELECTION POLICY: which component should perform it?
  ↓
EXECUTE
  ↓
VERIFY
  ↓
COMMIT / REJECT
  ↓
TRANSITION RECORD
  ↓
next controller decision
```

## Controller is not the router

The Controller decides **what behavior is needed next**.

The SelectionPolicy decides **which eligible component supplies that behavior**.

The Executor performs the behavior.

The Verifier decides whether the produced observations may enter committed state.

No one of these roles implicitly owns the others.

## External cognition rule

Parrot is Tonal's singular `EXTERNAL_MODEL` capability source. Control Loop R0 must never use it merely because machinery failed to appear.

A control decision explicitly declares whether external cognition is allowed for that transition.

```text
request capability X
allow_external_model = false
        ↓
qualified machinery exists?
  yes -> execute machinery
  no  -> record UNAVAILABLE
        ↓
controller may explicitly decide to escalate
        ↓
request capability X
allow_external_model = true
        ↓
Parrot may become eligible
```

This produces a measurable boundary between reusable machinery and purchased/probabilistic cognition.

## Machinery-first policy

When external cognition is explicitly allowed, the initial R0 selection policy still prefers:

1. deterministic non-external machinery;
2. other non-external machinery;
3. `EXTERNAL_MODEL` only when no eligible machinery remains.

This is intentionally simple and deterministic. Later reliability/cost-aware policies are experiments, not assumptions baked into R0.

## Commit rule

Execution output is not automatically system truth.

R0's contract verifier enforces at least:

- execution must produce an observation to commit;
- a non-`VERIFY` capability cannot emit/promote `FACT`;
- rejected transitions do not mutate the Blackboard;
- failures and rejections remain in the transition trace.

Later verifiers may add domain-specific checks without granting Parrot or a Tlaloque global verification authority.

## Transition record

Every iteration records:

```text
iteration
controller decision
external-cognition permission
eligible candidates
selected component kind/id
execution output/usage
verification verdict
commit/reject status
error or reason
```

This is the raw material for Episodes and later Tlaloc analysis.

## First self-hosting target

After the generic loop is stable, the first useful engineering workflow should be a bounded repository-repair cycle:

```text
RUN_TESTS                  Tool/Machine
  ↓
PARSE_FAILURE              Tlaloque/Machine
  ↓
LOCATE_RELEVANT_CODE       Tool/Tlaloque
  ↓
PROPOSE_HYPOTHESIS         machinery if known; otherwise Parrot
  ↓
DESIGN_CHECK / ABLATION    Tlaloque or explicit Parrot escalation
  ↓
PROPOSE_PATCH              machinery if known; otherwise Parrot
  ↓
APPLY_PATCH                constrained Tool
  ↓
RUN_TARGETED_TESTS         Tool
  ↓
RUN_REGRESSION_SUITE       Tool
  ↓
VERIFY_OUTCOME             verifier machinery
  ↓
RECORD_EPISODE
```

A proposed hypothesis, patch or change of direction is a **candidate**, not truth. It gains authority only through observed tests/evidence.

## Tlaloc connection

Tlaloc does not run this loop. It studies its verified Episodes and asks whether recurring uncertain work can become reusable machinery.

```text
repeated external-cognition transition
        ↓
Episode corpus
        ↓
Tlaloc detects bounded recurring structure
        ↓
candidate Tlaloque / Machine
        ↓
holdout + ablation + verification
        ↓
promotion
        ↓
future loop uses machinery instead
```

## R0 non-goals

- no learned controller;
- no autonomous self-promotion;
- no free-form Shponglese generation;
- no hidden chain-of-thought capture;
- no automatic architectural rewrite;
- no mandatory Origami dependency;
- no modification of frozen T1.

## Success criteria

R0 is successful when tests demonstrate all of the following:

1. a multi-transition goal can advance through committed Blackboard state;
2. machinery is selected before an external model when both are eligible;
3. Parrot is never invoked without explicit per-transition permission;
4. a rejected transition cannot mutate committed state;
5. unavailable machinery can be recorded and followed by explicit external escalation;
6. every transition leaves enough structured evidence to form an Episode later.
