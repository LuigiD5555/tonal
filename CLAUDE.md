# Tonal project instructions — Architecture R2

## North Star

Tonal investigates whether complex, reliable behavior can emerge from the composition of small, bounded, verifiable and reusable machinery plus selectively invoked external cognition, instead of requiring one increasingly complex general-purpose model to perform every operation.

Tonal is the complete heterogeneous runtime and research system. It owns goal intake, workflow/DAG execution, capability resolution, routing, scheduling, Blackboard state, resource accounting, verification coordination, tracing and final workflow results.

## Ecosystem roles

- **Tonal**: complete system/runtime and research authority for adaptive heterogeneous composition.
- **Tlaloc**: Tonal's capability foundry and Behavior Lab. It discovers, constructs, tests, qualifies, packages, promotes, deprecates and studies reusable machinery.
- **Tlaloque**: a bounded, typed, measurable reusable mechanism produced or qualified through Tlaloc. Tlaloques are machinery; they may be deterministic, algorithmic, symbolic, tool-backed, specialized-model-backed or hybrid.
- **Parrot**: Tonal's singular external probabilistic cognition interface. It is **not a Tlaloque**, is not produced by Tlaloc, and has no system-level authority. A concrete provider/model is runtime configuration beneath Parrot.
- **Shponglese**: semantic operational IR for primitive and composed behavior. Its semantics must remain independent of physical codec.
- **Origami**: independently testable representation, transport and virtual-memory substrate. Tonal may use it, but Origami remains independently measurable.

`Capability` is Tonal's common runtime abstraction. A Tlaloque and Parrot may both satisfy a capability without being the same kind of component.

Prefer the smallest reliable machinery that satisfies the required behavior. Invoke Parrot only when unresolved ambiguity, novelty, perception, interpretation or generation justifies external probabilistic cognition.

## Parrot invariant

Parrot exists to inject useful indeterminism into an otherwise structured mechanism.

```text
machinery
  -> unresolved/novel/ambiguous state
  -> Parrot
  -> probabilistic candidate/observation
  -> verification
  -> machinery continues
```

Parrot is not Tonal's brain, planner, router, verifier, memory or source of truth. Never introduce an implicit capability-to-Parrot fallback.

Tlaloc may characterize Parrot behavior and may use verified Episodes to replace recurring Parrot work with a Tlaloque or Machine. Tlaloc does not own or promote the external model itself.

## Research method

Do not add architectural complexity because a paper proposes it. Extract the experimentally demonstrated mechanism, map it to a measured Tonal failure or hypothesis, prototype the smallest falsifiable version, and promote only with evidence.

Current research themes include primitive composition, heterogeneous delegation, selective/frugal activation, process verification, experience-to-structure compilation, semantic IR/codec invariance and bounded memory/transport.

## T1 freeze rule

Experiment T1 is frozen with respect to its preregistered scientific question and experiment semantics.

T1 historically wrapped the external model through Tlaloc's R1 publication contract and called that wrapper a generative Tlaloque. Architecture R2 reclassifies Parrot conceptually as `external_model`, but must preserve frozen T1 behavior and accounting through compatibility adapters. Do not rewrite T1 history or results.

New research ideas may motivate post-T1 experiments. They must not be retrofitted into T1.

## Document authority

Read in this order:

1. `CLAUDE.md`
2. `README.md`
3. `docs/CURRENT_STATE.md`
4. `docs/ARCHITECTURE.md`
5. `docs/BOUNDARIES.md`
6. `docs/RESEARCH_PROGRAM.md` when research direction matters
7. the active experiment specification

Anything under `docs/archive/` is historical and is not current architectural authority. Frozen experiment specifications remain authoritative for their experiment.

## Change discipline

- Keep component kind, capability role, execution, routing and verification separate.
- Keep Tlaloc machinery distinct from external cognition.
- Do not make Tlaloc the Tonal runtime.
- Do not make Origami mandatory for Tonal correctness.
- Do not make Parrot mandatory for every workflow.
- Preserve explicit evidence and accounting boundaries.
- Prefer interchangeable policies and capability contracts over model-specific logic.
- Keep repository history; superseded architecture belongs under `docs/archive/`, not deleted.
