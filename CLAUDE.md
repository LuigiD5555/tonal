# Tonal project instructions — Architecture R2

## North Star

Tonal investigates whether complex, reliable behavior can emerge from the composition of small, bounded, verifiable and reusable capabilities instead of requiring one increasingly complex general-purpose model.

Tonal is the complete heterogeneous runtime and research system. It owns goal intake, workflow/DAG execution, capability resolution, routing, scheduling, Blackboard state, resource accounting, verification coordination, tracing and final workflow results.

## Ecosystem roles

- **Tonal**: complete system/runtime and research authority for adaptive heterogeneous composition.
- **Tlaloc**: Tonal's capability-development and Behavior Lab subsystem. It discovers, tests, qualifies, packages, promotes, deprecates and studies reusable capabilities.
- **Tlaloque**: one bounded, typed, measurable capability. It may be deterministic or probabilistic.
- **Parrot**: one probabilistic Tlaloque. It has no system-level authority and is not Tonal's brain, router, verifier, memory or source of truth.
- **Shponglese**: semantic operational IR for primitive and composed behavior. Its semantics must remain independent of physical codec.
- **Origami**: independently testable representation, transport and virtual-memory substrate. Tonal may use it, but Origami remains independently measurable.

Prefer the smallest reliable capability that satisfies the required behavior, evidence threshold and resource budget.

## Research method

Do not add architectural complexity because a paper proposes it. Extract the experimentally demonstrated mechanism, map it to a measured Tonal failure or hypothesis, prototype the smallest falsifiable version, and promote only with evidence.

Current research themes include:

- complexity through primitive composition;
- heterogeneous delegation;
- selective/frugal activation;
- process verification;
- experience -> reusable structure;
- semantic IR and codec invariance;
- bounded memory, transport and selective unfolding.

## Parrot invariant

Parrot is useful for ambiguous, generative or perceptual work when its marginal capability is needed. It is replaceable and must be routed through the same capability machinery as other executors.

Never introduce a privileged Parrot execution path.

## T1 freeze rule

Experiment T1 is frozen with respect to its preregistered scientific question and experiment semantics.

New research ideas from PAL, Least-to-Most, Decomposed Prompting, ViperGPT, LLMCompiler, RouteLLM, FrugalGPT, Stitch, DeepSeek-OCR or other papers may motivate post-T1 experiments. They must not be retrofitted into T1.

A future R2 runtime may provide a compatibility policy that reproduces T1 behavior, but frozen T1 artifacts, workflows, corpora, accounting definitions and experiment rules remain authoritative for T1.

## Document authority

Read in this order:

1. `CLAUDE.md`
2. `README.md`
3. `docs/CURRENT_STATE.md`
4. `docs/ARCHITECTURE.md`
5. `docs/BOUNDARIES.md`
6. `docs/RESEARCH_PROGRAM.md` when research direction matters
7. the active experiment specification

Anything under `docs/archive/` is historical and is not current architectural authority.

If an archived document conflicts with a current document, the current document wins. If a frozen experiment specification conflicts with a later architectural idea, the frozen experiment wins for that experiment.

## Change discipline

- Keep routing, execution, verification and capability qualification separable.
- Do not make Tlaloc the Tonal runtime.
- Do not make Origami mandatory for Tonal correctness.
- Do not make Parrot mandatory for every workflow.
- Preserve explicit evidence and accounting boundaries.
- Prefer interchangeable policies and capability contracts over hard-coded model-specific logic.
- Keep repository history; superseded architecture belongs under `docs/archive/`, not deleted.
