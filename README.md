# Tonal 0.1.0-alpha.5

Tonal is an optional **composition, compatibility, provenance and distribution layer for development toolchains and targets**.

It exists so independently versioned development systems and target projects can be combined reproducibly without pretending they share one codebase or one semantic authority.

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

Blueprint Framework is an example of another development tool that Tonal may compose later. It is **not currently locked** because no repository/version/SHA has been supplied in this composition. Tonal never invents a component pin from a conceptual example.

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
                 owns its own releases
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
TONAL VERIFICATION != UNIVERSAL MODEL CAPABILITY
TLALOC != REQUIRED ORIGAMI RUNTIME
```

## Generic component manifest

`TONAL.json` now declares each component's:

```text
kind
repository
version
immutable commit
authority/role
deterministic verification commands
```

`tonal.lock` mirrors the immutable identity fields.

Adding a future development tool no longer requires teaching Tonal that the only valid names are `tlaloc` and `origami`. A tool first has to be explicitly registered and pinned; only then is it part of the composition.

## Verification

```bash
./scripts/verify-stack.sh
```

Verification does four distinct things:

1. verifies manifest/lock coherence and component kinds;
2. fetches the exact immutable component commits;
3. verifies each checkout's commit/version identity;
4. executes the component-specific deterministic checks declared in `TONAL.json`.

Historical Fixed Carrier R2 regression gates are retained as additional evidence rather than deleted.

A green result means:

```text
COMPOSITION_VERIFIED
```

It does **not** mean that Origami Native visual support, Hybrid support or another empirical capability has been promoted.

## Generic snapshots

```bash
./scripts/build-snapshot.sh
```

The snapshot builder now includes **every component in `tonal.lock`**, rather than hard-coding directories for only Tlaloc and Origami.

The result remains a reproducible distribution artifact, not a new source repository.

## Component operations

```bash
./scripts/component-status.sh
./scripts/update-component.sh <declared-component> [revision]
```

`update-component.sh` accepts any component already declared in the manifest. It refuses unknown names instead of silently creating an ungoverned component.

## Current exact pins

At this branch stage:

```text
Tlaloc  6.0.0-alpha.14
9b83e76f33c888b8701d3ab6a48049425ba7b8e8

Origami 6.0.0-alpha.10
fe4ba64ef35dc900c84fa5c20c9afad4fceee173
```

The Origami pin will only advance when the pending portable-baseline release is actually merged and green; Tonal will not pin an unmerged branch as a released target revision.

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
```

See `PROJECT_BOUNDARY.md`, `TONAL.json`, `tonal.lock` and `compatibility.json`.
