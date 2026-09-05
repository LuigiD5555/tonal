# Tonal 0.1.0-alpha.5

**Architecture R2 foundation**

Tonal is an experimental **heterogeneous cognitive/runtime system** for composing bounded, measurable machinery with selectively invoked external cognition into verified workflows.

The core research question is whether complex, reliable behavior can emerge from external structure, specialization, verification, memory and reusable experience while reducing unnecessary dependence on general-purpose model inference.

## Canonical roles

```text
TONAL
  complete runtime/research system
  goal -> DAG -> select -> execute -> verify -> Episode

TLALOC
  capability foundry + Behavior Lab
  builds/qualifies reusable machinery

TLALOQUE
  bounded typed measurable machinery produced/qualified through Tlaloc

PARROT
  singular external probabilistic cognition interface
  NOT a Tlaloque; no system-level authority

SHPONGLESE
  semantic operational IR

ORIGAMI
  optional representation / transport / virtual-memory substrate
```

`Capability` is the common runtime abstraction above these component kinds. Parrot is useful where ambiguity, novelty, perception, interpretation or generation actually needs probabilistic cognition. It is **not** Tonal's brain, router, verifier, memory or default executor.

A provider/model such as a local VLM or remote frontier model is configuration beneath Parrot rather than a new architectural species.

## Runtime

```text
goal / task family
       ↓
workflow DAG
       ↓
CapabilityRegistry
       ↓
SelectionPolicy
       ↓
Tlaloque / Machine / Tool / Parrot
       ↓
verification
       ↓
Blackboard commit/reject + accounting
       ↓
RunRecord / Episode
```

When machinery cannot resolve a state, Tonal may invoke Parrot to produce a probabilistic candidate and then return control to verification and machinery.

## Tlaloc learning direction

Tlaloc may study verified Episodes and ask whether recurring Parrot-assisted behavior can be replaced by cheaper reusable structure:

```text
Parrot-assisted success
       ↓
verified Episodes
       ↓
Tlaloc pattern discovery
       ↓
candidate Tlaloque / Machine
       ↓
holdout + ablation + verification
       ↓
promotion or rejection
```

This is a research target, not yet a demonstrated autonomous learning loop.

## T1 remains frozen

T1 asks whether heterogeneous composition degrades less with workflow depth than monolithic Parrot or a decomposed-but-Parrot-centric workflow.

The frozen R1 adapter historically published Parrot through Tlaloc as a generative Tlaloque. Architecture R2 may reclassify that adapter at the Tonal boundary as `EXTERNAL_MODEL`, but does not rewrite T1 artifacts, calls, accounting or experiment semantics.

## Research program

1. **T2 Primitive Swarm / MICRO-ISA** — can a small set of bounded primitives compose into a much larger behavior space?
2. **T3 Selective / Frugal Compute** — can Tonal activate only the minimum sufficient machinery and invoke Parrot only when needed?
3. **T4 Cognitive JIT** — can verified Episodes become cheaper reusable Tlaloques/Machines through Tlaloc?
4. **T5 Shponglese** — can a primitive/motif IR generalize compositionally and remain codec-invariant?
5. **T6 Origami carrier/memory** — can Origami transport/address the same semantics more effectively than conventional codecs under fair anti-prior controls?

See `docs/RESEARCH_PROGRAM.md` and `docs/research/PAPER_MAP.md`.

## Historical composition infrastructure

Tonal originally existed primarily as an optional composition, compatibility, provenance and snapshot layer. That machinery remains useful infrastructure and is retained (`TONAL.json`, `tonal.lock`, compatibility/gatekeeper state and component-fetch/verification scripts). Composition/pinning/provenance are now **part of Tonal**, not the complete definition of Tonal.

## Exact component pins

`TONAL.json` and `tonal.lock` remain the machine-readable authority for exact external revisions used by a frozen composition or experiment.

## Verification

```bash
./scripts/verify-stack.sh
```

The R2 gate also tests/vets the runtime against the exact Tlaloc revision materialized from `tonal.lock`.

## Documentation authority

Start with `CLAUDE.md`, this README, `docs/CURRENT_STATE.md`, `docs/ARCHITECTURE.md`, `docs/BOUNDARIES.md`, then the active experiment specification. Anything under `docs/archive/` is historical.

## Core invariant

> Use reliable reusable machinery where it is sufficient; spend external probabilistic cognition only where unresolved uncertainty or novelty requires it.
