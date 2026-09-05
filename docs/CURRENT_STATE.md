# Tonal current state — R2 Foundation

**Status:** migration foundation in progress  
**Authority:** Architecture R2  
**T1 status:** frozen experiment line; do not retrofit new ideas into T1

## What Tonal is now

Tonal is the complete heterogeneous cognitive/runtime system and research authority for composing bounded capabilities into verified workflows.

The repository currently contains two historical layers:

1. a composition/pinning/provenance layer created before the runtime existed;
2. the T1 runtime subsystem that introduced workflow/DAG execution, routing, Blackboard state, accounting and heterogeneous Tlaloque execution.

Architecture R2 keeps both, but changes their hierarchy:

```text
Tonal
├── runtime / cognitive authority
├── capability registry and selection
├── verification and evidence
├── experiments
└── composition/pinning/provenance infrastructure
```

The composition layer is infrastructure inside Tonal. It no longer defines Tonal as a whole.

## Canonical ecosystem roles

```text
TONAL
  complete runtime/research system

TLALOC
  capability foundry + Behavior Lab

TLALOQUE
  bounded typed measurable capability

PARROT
  one probabilistic Tlaloque; no privileged system role

SHPONGLESE
  semantic operational IR

ORIGAMI
  representation / transport / virtual-memory substrate
```

## Current protected work

T1 asks whether heterogeneous composition degrades less with workflow depth than monolithic Parrot or decomposed-but-Parrot-centric execution.

T1 remains frozen. R2 may build compatibility adapters and future policies around it, but may not change its experiment semantics, corpus, workflows, accounting or preregistered comparison.

## R2 Foundation goals

The current migration has three immediate goals:

1. reorganize repositories and archive superseded authority documents;
2. align Tonal, Tlaloc and Origami documentation around one shared architecture;
3. introduce generic capability and selection abstractions so Parrot becomes one interchangeable executor rather than a privileged path.

## First implementation target after documentation migration

```text
Capability
CapabilityProfile
Executor
Verifier
SelectionPolicy
Episode
```

The first compatibility policy should reproduce frozen T1 routing before any adaptive policy is introduced.

## Non-goals of R2 Foundation

- no DreamCoder/Stitch implementation yet;
- no learned router yet;
- no HRM-like controller yet;
- no Origami dependency for Tonal correctness;
- no free-form Shponglese generation;
- no modification of frozen T1 evidence.
