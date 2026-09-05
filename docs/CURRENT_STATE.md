# Tonal current state — R2 Foundation

**Status:** migration foundation in progress  
**Authority:** Architecture R2  
**T1 status:** frozen experiment line; do not retrofit new ideas into T1

## What Tonal is now

Tonal is the complete heterogeneous cognitive/runtime system and research authority for composing bounded capabilities into verified workflows.

The repository contains historical composition/pinning/provenance infrastructure plus the runtime subsystem introduced during T1. Architecture R2 keeps both but makes the runtime/system boundary authoritative.

## Canonical ecosystem roles

```text
TONAL
  complete runtime/research system

TLALOC
  capability foundry + Behavior Lab
  manufactures/qualifies reusable machinery

TLALOQUE
  bounded typed measurable machinery produced/qualified through Tlaloc

PARROT
  singular external probabilistic cognition interface
  NOT a Tlaloque and NOT produced by Tlaloc

SHPONGLESE
  semantic operational IR

ORIGAMI
  representation / transport / virtual-memory substrate
```

`Capability` is the common Tonal abstraction above component kinds. A capability may be supplied by a Tlaloque, Machine, Tool or external model without collapsing those categories.

## Parrot role

Parrot exists to inject useful indeterminism when reusable machinery is insufficient. A concrete model/provider is configuration beneath Parrot. Its output is evidence/candidate material that re-enters Tonal verification and Blackboard mechanics; it is not automatic system truth.

Tlaloc may study recurring verified Parrot Episodes and propose machinery that replaces that recurring uncertainty in future runs.

## Current protected work

T1 asks whether heterogeneous composition degrades less with workflow depth than monolithic Parrot or decomposed-but-Parrot-centric execution.

T1 remains frozen. Its historical R1 adapter modeled Parrot as a generative Tlaloque for experiment execution. R2 reclassifies that adapter at the Tonal boundary as `external_model` without changing T1 experiment semantics, calls, corpus, workflows or accounting.

## R2 Foundation goals

1. reorganize repositories and archive superseded authority documents;
2. align Tonal, Tlaloc and Origami around one shared architecture;
3. separate component kind from capability role;
4. make Tonal's Registry and SelectionPolicy independent from Tlaloc implementation details;
5. stabilize verification and Episode contracts;
6. build the first Tonal Control Loop using reusable machinery plus selective Parrot escalation.

## Immediate implementation sequence

```text
CapabilityKind
  ↓
Tonal-owned CapabilityRegistry
  ↓
Tlaloc adapter + Parrot adapter
  ↓
SelectionPolicy
  ↓
Executor
  ↓
Verifier
  ↓
Blackboard commit/reject
  ↓
Episode
  ↓
next transition
```

## Non-goals of R2 Foundation

- no learned router yet;
- no autonomous self-promotion;
- no free-form Shponglese generation;
- no mandatory Origami dependency;
- no modification of frozen T1 evidence.
