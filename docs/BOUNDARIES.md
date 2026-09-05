# Tonal ecosystem boundaries — Architecture R2

## Tonal owns

Tonal owns runtime authority for:

- goal intake;
- workflow/DAG representation;
- capability discovery and selection;
- routing and scheduling;
- Blackboard state;
- dependency coordination;
- execution state;
- resource accounting;
- verification coordination;
- trace generation;
- final workflow result;
- experiment-level system metrics.

Tonal may consume capabilities produced by Tlaloc and representations provided by Origami, but neither subsystem owns Tonal runtime decisions.

## Tlaloc owns

Tlaloc owns the lifecycle of reusable capabilities:

- behavior decomposition for experiments;
- bounded Tlaloque construction;
- capability contracts;
- executor qualification;
- CapabilityProfile evidence;
- Behavior Lab campaigns;
- candidate discovery and comparison;
- promotion and deprecation decisions for Tlaloc-managed capabilities;
- Episode analysis and future experience-to-structure compilation.

Tlaloc does not own Tonal's scheduler or Blackboard and does not self-promote individual workers.

## Tlaloque boundary

A Tlaloque is one bounded, typed, measurable capability.

A Tlaloque may be:

- deterministic;
- neural;
- model-backed;
- tool-backed;
- symbolic;
- hybrid.

A Tlaloque must not gain authority merely because it uses a language model.

## Parrot boundary

Parrot is one probabilistic Tlaloque.

It may provide capabilities such as ambiguous perception, extraction, generation or interpretation when those capabilities are justified by evidence.

Parrot does not own:

- system planning;
- routing;
- verification;
- memory;
- Blackboard state;
- semantic truth;
- capability promotion.

There must be no architectural requirement that Parrot participate in every workflow.

## Shponglese boundary

Shponglese owns semantic operational representation for primitive and composed behavior.

Its first role is an IR, not a claim of emergent language. Semantics should remain stable across codecs so the same program can be represented as text, JSON, binary, Unicode or Origami without changing meaning.

Primitive, motif, graph and pack abstractions may be explored only with explicit compositional and held-out tests.

## Origami boundary

Origami owns representation, carrier, addressing, selective unfolding and virtual-memory mechanisms that it independently validates.

Origami may carry Shponglese or other semantic structures, but it does not own Shponglese semantics.

Tonal correctness must not depend on Origami unless a specific experiment explicitly tests an Origami-dependent configuration.

## Dependency direction

```text
Tlaloc
  produces qualified capabilities/profiles
          │
          ▼
       Tonal Registry
          │
          ▼
Tonal SelectionPolicy -> Executor -> Verifier -> Blackboard/Episode

Shponglese
  semantic program / IR
          │
          ├── conventional codecs
          └── Origami carrier/memory
```

## Core invariant

Prefer the smallest reliable capability that satisfies the required behavior, evidence threshold and resource budget.
