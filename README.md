# Tonal — Architecture R2 foundation

Tonal is an experimental **heterogeneous cognitive/runtime system** for composing bounded, measurable capabilities into verified workflows.

The core research question is whether complex, reliable behavior can emerge from external structure, specialization, verification, memory and reusable experience instead of requiring one increasingly complex general-purpose model to perform every operation.

## Canonical roles

```text
TONAL
  complete runtime/research system
  goal -> DAG -> select -> execute -> verify -> Episode

TLALOC
  capability foundry + Behavior Lab
  build -> test -> qualify -> promote/deprecate

TLALOQUE
  one bounded typed measurable capability

PARROT
  one probabilistic Tlaloque
  no system-level authority

SHPONGLESE
  semantic operational IR

ORIGAMI
  optional representation / transport / virtual-memory substrate
```

Parrot is useful where probabilistic perception, extraction, interpretation or generation is actually needed. It is **not** Tonal's brain, router, verifier, memory or default executor for every task.

## Runtime

The T1 runtime introduced the system boundary now promoted by Architecture R2:

```text
goal / task family
       ↓
workflow DAG
       ↓
qualified capability candidates
       ↓
SelectionPolicy
       ↓
executor
       ↓
Blackboard + verification + accounting
       ↓
RunRecord / Episode-compatible evidence
```

See `runtime/` and `docs/ARCHITECTURE.md`.

## T1 remains frozen

T1 asks whether heterogeneous composition degrades less with workflow depth than monolithic Parrot or a decomposed-but-Parrot-centric workflow.

Architecture R2 does not retrofit new research ideas into T1. Frozen T1 artifacts and experiment semantics remain authoritative for T1.

## Research program

Post-T1 research is organized around falsifiable mechanisms rather than copying large architectures:

1. **T2 Primitive Swarm / MICRO-ISA** — can a small set of bounded primitives compose into a much larger behavior space?
2. **T3 Selective / Frugal Compute** — can Tonal activate only the minimum sufficient capabilities?
3. **T4 Cognitive JIT** — can verified Episodes become cheaper reusable capabilities through Tlaloc?
4. **T5 Shponglese** — can a primitive/motif IR generalize compositionally and remain codec-invariant?
5. **T6 Origami carrier/memory** — can Origami transport/address the same semantics more effectively than conventional codecs under fair anti-prior controls?

See `docs/RESEARCH_PROGRAM.md` and `docs/research/PAPER_MAP.md`.

## Historical composition infrastructure

Tonal originally existed primarily as an optional composition, compatibility, provenance and snapshot layer. That machinery remains useful infrastructure and is retained:

```text
TONAL.json
tonal.lock
compatibility.json
gatekeeper.json
scripts/fetch-components.sh
scripts/verify-components.sh
scripts/build-snapshot.sh
```

The change in Architecture R2 is one of hierarchy: composition/pinning/provenance are **part of Tonal**, not the complete definition of Tonal.

The exact pre-R2 project boundary is preserved under `docs/archive/superseded-architecture/`.

## Exact component pins

`TONAL.json` and `tonal.lock` remain the machine-readable authority for the exact external component revisions used by a composition or frozen experiment. Do not infer current pins from prose in this README.

## Verification

Historical composition verification remains available:

```bash
./scripts/verify-stack.sh
```

The runtime Go module is under `runtime/`. It consumes the exact Tlaloc revision materialized by the repository's component-fetch workflow.

## Documentation authority

Start with:

1. `CLAUDE.md`
2. this `README.md`
3. `docs/CURRENT_STATE.md`
4. `docs/ARCHITECTURE.md`
5. `docs/BOUNDARIES.md`
6. the active experiment specification

Anything under `docs/archive/` is historical and not current architectural authority.

## Core invariant

> Prefer the smallest reliable capability that satisfies the required behavior, evidence threshold and resource budget.
