# TONAL system boundary — R1

Addendum, not a rewrite. `PROJECT_BOUNDARY.md`, `TONAL.json` and the
existing composition/pinning/snapshot/evidence layer describe the
historical Tonal and remain accurate for what they cover. This document
adds the boundary for the **runtime subsystem** introduced by experiment
**T1**.

## Two concerns in one repository

| | Historical composition layer (repo root) | Runtime subsystem (`runtime/`) — new in T1 |
|---|---|---|
| Purpose | reproducible exact-commit composition of independent tools/targets; provenance; snapshot; evidence ledger | the cognitive/computational system: it takes a goal and produces a verified answer |
| Artifacts | `TONAL.json`, `tonal.lock`, `scripts/`, `.work/components/`, gatekeeper, `docs/EVIDENCE_PLAN.md` | `runtime/` Go module `tonal.local/runtime` |
| Authority | "is this combination of pinned components reproducible?" | "did Tonal route this workflow across heterogeneous Tlaloques and get it right?" |
| Changes here? | unchanged by T1 | entirely new |

Both stay. The composition layer is infrastructure; the runtime subsystem
is the TONAL cognitive authority from T1 onward.

## Canonical responsibilities (T1+)

**TONAL** (this repo, `runtime/`) owns: goal intake · workflow/DAG ·
capability resolution · Registry consumption · routing · scheduler /
execution · Blackboard · execution state · observations · dependency
coordination · resource accounting · verification coordination · final
workflow result.

**TLALOC** (pinned component) owns: Tlaloque construction, packaging,
capability-contract definition, executor qualification, CapabilityProfile
production, competence-envelope evidence, tests, versioning, promotion, and
publication into Tonal's Registry — exposed through the public package
`tlaloc.local/behaviorlab/tlaloquekit`.

**TLALOQUE**: one bounded executable capability with input/output
contracts, declared capabilities, dependencies, and a profile/evidence
reference.

**PARROT**: one Tlaloque — `LFM2-VL 1.6B + CapabilityProfile R1 +
AdapterR1`. No system-level role.

**ORIGAMI**: out of scope for T1.

## Dependency direction

```
tlaloc  (pinned by tonal.lock at b22b856 / alpha.21)
  behavior-lab/internal/*              frozen experiment engine
        │ exposed by
        ▼
  behavior-lab/tlaloquekit             public, no internal/* leak
        │ consumed by  (module replace on the pinned .work checkout)
        ▼
runtime/tlalocbridge                   adapter — owns nothing cognitive
        ▼
runtime/tonal                          TONAL: goal, DAG, Blackboard,
                                       routing, scheduler, accounting
        ▼
runtime/cmd/tonal-t1                   A/B/C experiment harness
```

- Tonal depends on Tlaloc's published kit. Never the reverse.
- Tonal is not a submodule of Tlaloc; `behavior-lab` never imports
  `runtime`.
- `runtime/` consumes the exact Tlaloc commit pinned by `tonal.lock`
  (`scripts/fetch-components.sh` → `.work/components/tlaloc/behavior-lab`),
  not a machine-specific local checkout.

## What T1 does NOT do

No Origami integration · no new model characterization or training · no
Grounding R1 · no autonomous Tlaloque synthesis (that is Tlaloc T2) · no
free-form LLM planner · no bulk Tlaloque creation. Tlaloc's T1 role is
limited to packaging and qualifying a small fixed Tlaloque set.
