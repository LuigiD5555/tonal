# Tonal 0.1.0-alpha.6

Tonal is an optional **composition, compatibility, provenance and distribution layer for development toolchains and targets**.

It exists so independently versioned development systems and target projects can be combined reproducibly without pretending they share one codebase or one semantic authority.

## Current composition

```text
Tlaloc 6.0.0-alpha.16
  development_tool
  7f3393079d163131ad690f97de59ab2e2a249179

Origami 6.0.0-alpha.13
  target
  7dbd7ba073b227377c6cc3ae592f4c4f2573dabf
```

This pair contains the second-half Origami protocol work:

```text
Origami
  Protocol R0
  S*/E* codec registry
  capability negotiation
  Master Prompt R4
  rendered profile-3 at 640x640 / 8192 bytes
  deterministic S2(E2(INDEX)) roundtrip

Tlaloc
  codec-aware Native regression
  tlaloc-protocol-eval
  deterministic READ/WRITE/ROUNDTRIP/MULTIHOP evaluation
  semantic drift/invention/external-codec/exact-escalation metrics
```

## The most important boundary

A green Tonal composition means:

```text
COMPOSITION_VERIFIED
```

It does **not** mean:

```text
PROTOCOL_INTEROPERABILITY_PROMOTED
NATIVE_SEMANTIC_PROMOTED
UNIVERSAL_MODEL_CAPABILITY
```

The current exact commits have deterministic component evidence for renderer/roundtrip/evaluators. Held-out real-model evidence remains pending for:

```text
Native S2 index recovery
Native E2 write/construction behavior
A -> B -> C cross-model Origami interoperability
```

Those claims remain owned by the relevant component/evidence process, not by Tonal.

## Composition model

Tonal v2 recognizes generic component kinds:

```text
development_tool
target
support
```

The current locked composition contains:

```text
Tlaloc  -> development_tool
Origami -> target
```

Blueprint Framework is an example of another development tool that Tonal may compose later. It is **not currently locked** because no exact repository/version/SHA has been supplied. Tonal never fabricates a component pin from a conceptual example.

## Ecosystem role

```text
                    TONAL (optional)
          exact pins / checks / provenance / snapshot
            /                 |                 \
       TLALOC          Blueprint Framework     other tools
 development_tool       possible future         ...
            \                 |                 /
             \------ development ecosystem ---/
                            |
                            v
                         ORIGAMI
                           target
          owns protocol/profile semantics/releases
```

Tonal is not required to run Tlaloc or Origami.

## Authority rule

```text
Development tool
  -> develops / experiments / proposes / produces evidence

Target project
  -> owns its semantics, contracts and releases

Tonal
  -> pins exact revisions
  -> verifies the declared composition
  -> records provenance
  -> builds reproducible snapshots
```

For Origami:

```text
Tlaloc / another development tool
        ↓
candidate + evidence
        ↓
Origami validation/adoption
        ↓
Origami release/profile
        ↓
optional Tonal composition/pin
```

Therefore:

```text
TONAL COMPOSITION != ORIGAMI PROFILE PROMOTION
TONAL VERIFICATION != PROTOCOL INTEROPERABILITY PROMOTION
TONAL VERIFICATION != UNIVERSAL MODEL CAPABILITY
TLALOC != REQUIRED ORIGAMI RUNTIME
```

## Generic component manifest

`TONAL.json` declares each component's:

```text
kind
repository
version
immutable commit
authority/role
deterministic verification commands
```

`tonal.lock` mirrors the immutable identity fields.

A future tool first has to be explicitly registered and pinned; only then is it part of the composition.

## Verification

```bash
./scripts/verify-stack.sh
```

Verification:

1. verifies manifest/lock coherence and component kinds;
2. fetches the exact immutable component commits;
3. verifies each checkout's commit/version identity;
4. executes the component-specific deterministic checks declared in `TONAL.json`.

Historical Fixed Carrier R2 regression metadata is retained as historical evidence rather than rewritten as current protocol evidence.

## Generic snapshots

```bash
./scripts/build-snapshot.sh
```

The snapshot builder includes every component in `tonal.lock`. The result is a reproducible distribution artifact, not a new source repository or a transfer of semantic authority.

## Component operations

```bash
./scripts/component-status.sh
./scripts/update-component.sh <declared-component> [revision]
```

`update-component.sh` accepts any component already declared in the manifest. It refuses unknown names instead of silently creating an ungoverned component.

## Shared workflow

Tonal remains the canonical distribution owner for project-agnostic `repo-flow` and related shared integration/Gatekeeper workflow assets where those are used.

That ownership concerns repository workflow and composition. It does not transfer semantic authority from component repositories.

## Hard boundaries

```text
COMPONENT REPOSITORIES REMAIN AUTHORITATIVE
TARGET REPOSITORY OWNS TARGET RELEASES
TONAL IS OPTIONAL
UNPINNED TOOL != COMPOSITION COMPONENT
SNAPSHOT REFERENCES EXACT COMMITS
DECLARED COMPONENT VERIFICATION MUST PASS
TONAL COMPOSITION != TARGET PROMOTION
COMPOSITION VERIFIED != PROTOCOL INTEROPERABILITY PROMOTED
COMPOSITION VERIFIED != NATIVE SEMANTIC PROMOTED
REAL MODEL EVIDENCE REMAINS COMPONENT OWNED
```

See `PROJECT_BOUNDARY.md`, `TONAL.json`, `tonal.lock` and `compatibility.json`.
