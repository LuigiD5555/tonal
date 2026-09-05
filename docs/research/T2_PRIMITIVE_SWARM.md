# T2 — Primitive Swarm / MICRO-ISA

**Status:** planned post-T1 experiment. Do not modify frozen T1 to implement this.

## Question

Can a small library of bounded, verifiable primitives cover a much larger family of tasks through composition, with lower neural dependence and/or slower degradation than a general-purpose model performing every operation?

## Initial principle

Global task complexity does not imply that every local operation must be complex.

T2 should test whether Tonal can express difficult workflows as a DAG over a small operational vocabulary and route each node to the cheapest reliable eligible capability.

## Candidate initial vocabulary

```text
LOCATE
READ
EXTRACT
NORMALIZE
COMPARE
ADD
SUBTRACT
FILTER
MAP
REDUCE
BRANCH
VERIFY
RETURN
```

The exact vocabulary is experimental. Do not promote primitives because they are aesthetically appealing; retain only operations with measurable reuse and bounded semantics.

## Initial baselines

At minimum compare:

1. general probabilistic executor for the complete task;
2. decomposed workflow where the general model still performs most nodes;
3. primitive composition using deterministic/specialized capabilities where available.

T1 may provide methodological precedent but T2 must have its own frozen protocol before execution.

## Metrics

Candidate metrics:

- task success;
- degradation vs composition depth;
- neural/model calls;
- primitive activations;
- routing mistakes;
- verification failures;
- latency/cost;
- capability reuse;
- coverage of task families by the primitive set.

## Motif direction

Repeated verified primitive trajectories may later become candidate Shponglese motifs or Machines, for example:

```text
EXTRACT -> NORMALIZE -> SUBTRACT -> VERIFY
```

becoming a candidate macro such as `RECONCILE_NUMBERS` only after explicit holdout, verification, reuse and complexity/cost tests.

## Relevant mechanisms

PAL, Least-to-Most, Decomposed Prompting, ViperGPT, VisProg and LLM+P motivate the central experimental mechanism: move well-defined subproblems out of a general model and compose bounded operations instead.
