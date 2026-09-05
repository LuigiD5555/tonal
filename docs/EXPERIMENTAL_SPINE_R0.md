# Experimental Spine R0 — TONAL integration

Status: experimental branch infrastructure. This does not change TONAL's cognitive/runtime authority or Origami semantics.

## Purpose

TONAL already has a rich native `RunRecord` for workflow execution. The goal here is not to replace it. The goal is to make every fast prototype run produce a small, reusable experience bundle that the next prototype iteration can compare.

```text
TONAL RunRecord[] + explicit experiment evaluation
                 |
                 v
        runtime/experience adapter
                 |
                 v
     Tlaloc public prototypelab API
                 |
                 v
          experience/
            manifest.json
            episodes/...
            summary.json
```

## Authority boundary

`runtime/tonal.RunRecord` remains TONAL's execution source of truth. It records routing, selected workers, Blackboard reads/writes, model calls, latency, final observation and accounting.

`tlaloc.local/behaviorlab/prototypelab` owns the generic development/learning projection: Episode, RunManifest, Summary and immutable bundle writing.

Origami is not involved in this integration and is not modified by it.

## Why evaluation is explicit

A workflow may finish with:

```text
FinalStatus = OK
```

while still giving the wrong task answer. Therefore `runtime/experience` never infers semantic correctness from execution completion.

The experiment harness supplies:

```go
experience.Evaluation{
    Success:          ...,
    SemanticCorrect:  ...,
    ExactCorrect:     ...,
    FailureRootCause: ...,
}
```

and the adapter combines that evaluation with the native RunRecord.

```text
EXECUTION SUCCESS != SEMANTIC CORRECTNESS
```

## One-call prototype persistence

After a batch of workflows has been evaluated:

```go
paths, err := experience.WriteBundle(
    outDir,
    manifest,
    evaluatedRuns,
    observedAt,
)
```

This writes the common immutable `experience/` view while leaving native TONAL traces authoritative.

The manifest can link iterations through:

```text
run_id
parent_run_id
prototype.version
prototype.parent_version
prototype.hypothesis
prototype.change_summary
```

That gives later prototypes a direct trail of:

```text
previous run
 -> observed failure frontier
 -> explicit change/hypothesis
 -> new run
 -> outcome
```

without requiring a database or a new orchestration framework.

## Accounting discipline

TONAL RunRecord currently observes generative call counts and runtime latency, so the generic Episode projection records those values.

It does **not** invent lower-level HTTP attempt/failure counts that TONAL's native trace does not observe. T1-specific raw records in Tlaloc have richer transport accounting and are adapted there instead.

## Exact experimental component pin

This branch pins the Tlaloc experimental spine implementation by immutable commit in both `TONAL.json` and `tonal.lock`. The pin is an experimental composition dependency, not a Tlaloc release promotion.

## Development rule

Use the spine to reduce ceremony, not increase it:

```text
DEV            -> touched package tests
SMOKE REAL     -> 1-3 real tasks
BUILD-TO-LEARN -> small expanding batches that expose failure structure
FREEZE         -> full gates only when evidence is worth freezing/promoting
```

## Hard boundaries

```text
TONAL RUNRECORD REMAINS EXECUTION AUTHORITY
TLALOC PROTOTYPELAB OWNS COMMON EXPERIENCE PROJECTION
EXPLICIT EVALUATION OWNS TASK CORRECTNESS
EXECUTION SUCCESS != SEMANTIC CORRECTNESS
UNKNOWN ACCOUNTING IS NOT INVENTED
EXPERIENCE SUMMARY != PROMOTION DECISION
ORIGAMI SEMANTICS REMAIN UNTOUCHED
```
