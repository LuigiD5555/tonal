# Tonal Architecture R2

## 1. System definition

Tonal is the complete heterogeneous cognitive/runtime system and research authority for composing bounded capabilities into verified workflows.

Architecture R2 promotes the runtime boundary introduced during T1 from an experiment-local addendum to the canonical system architecture.

The historical composition/pinning/provenance layer remains useful infrastructure, but it no longer defines Tonal as a whole.

## 2. Core runtime responsibilities

Tonal owns:

- goal intake;
- workflow/DAG representation;
- capability registry consumption;
- capability selection and routing;
- scheduling and dependency coordination;
- Blackboard state;
- execution state;
- resource accounting;
- verification coordination;
- trace generation;
- final workflow results;
- experiment-level system metrics.

## 3. Capability-first architecture

The runtime should answer:

```text
Which capability should execute this bounded operation?
```

not:

```text
Which Parrot prompt should solve this?
```

A capability may be implemented by a deterministic function, symbolic solver, tool, state machine, specialized model, Parrot or another bounded executor.

Canonical separation:

```text
Capability          what behavior is offered
CapabilityProfile   evidence, competence envelope, cost/reliability history
Executor            performs the work
Verifier            judges result/evidence
SelectionPolicy     chooses among eligible capabilities
```

These roles should remain replaceable and independently testable.

## 4. Parrot

Parrot is one probabilistic Tlaloque/capability.

Parrot may be valuable for perception, ambiguity, extraction, interpretation or generation. It must participate through the same registry/selection/execution boundaries as other capabilities.

Parrot is not the router, global planner, verifier, memory, semantic authority or mandatory executor.

## 5. Tlaloc

Tlaloc is the capability foundry and Behavior Lab.

Tlaloc may construct, qualify, package, promote and deprecate capabilities; produce CapabilityProfiles; run behavior experiments; and study verified Episodes for reusable structure.

Tlaloc does not own Tonal workflow execution or Blackboard authority.

## 6. Shponglese

Shponglese begins as semantic operational IR for primitive and composed behavior.

Expected conceptual hierarchy:

```text
ATOM -> MOTIF -> GRAPH -> PACK
```

Its meaning must remain independent of codec. Claims of language require later empirical evidence for compositionality, systematicity, slot generalization, held-out composition and codec invariance.

## 7. Origami

Origami is an independently testable representation, transport, addressing and virtual-memory substrate.

Tonal may use Origami where measured evidence shows value. Origami must not become an undeclared requirement for Tonal correctness and does not own Shponglese semantics.

## 8. Experience loop

Target long-term loop:

```text
novel task
   ↓
selected capabilities / probabilistic reasoning when needed
   ↓
verification
   ↓
Episode
   ↓
Tlaloc analysis
   ↓
recurring reliable structure?
   ↓
candidate Machine / primitive / motif / specialist
   ↓
qualification + holdout + ablation
   ↓
Registry
   ↓
cheaper future execution
```

This is a research target, not yet a demonstrated system property.

## 9. T1 compatibility

T1 remains scientifically frozen.

R2 may later expose a `FrozenT1SelectionPolicy` or equivalent compatibility path that reproduces T1 routing through generic interfaces. The goal is regression compatibility, not alteration of T1.

## 10. Architectural invariant

> Prefer the smallest reliable capability that satisfies the required behavior, evidence threshold and resource budget.
