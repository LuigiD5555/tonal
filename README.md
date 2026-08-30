# Tonal

Tonal is an optional **composition, compatibility, provenance and distribution layer for development toolchains**.

It exists so independently versioned tools can be combined reproducibly without pretending they belong to one codebase or one authority domain.

Examples of tools/components Tonal may compose include:

```text
Tlaloc
Blueprint Framework
Origami revisions used as development targets/artifacts
future alternative or complementary development kits
verification/evaluation tooling
```

Tonal is **not** the authority for Origami semantics, visual grammar or profile promotion. Origami owns its own canonical versions. Tlaloc is one development kit that can experiment on and improve Origami, but Tonal is deliberately broader than the pair `Tlaloc + Origami`.

## Ecosystem roles

```text
                           TONAL
          optional reproducible toolchain composition
        /                 |                    \
   TLALOC          Blueprint Framework      other tools
 development kit      development tool      future tools
        |
        | can target / improve / validate
        v
                           ORIGAMI
          representation language + machine + memory
          owns its own canonical versions/profiles
```

The same pattern works for targets other than Origami: Tlaloc or Blueprint Framework can participate in a Tonal composition without Origami being present.

## Authority rule

```text
Development tool
  -> proposes / builds / tests / produces evidence

Target project
  -> owns its own contracts and canonical releases

Tonal
  -> pins exact revisions
  -> verifies integration/compatibility
  -> records provenance
  -> builds reproducible snapshots/distributions
```

For Origami specifically:

```text
Tlaloc or another development tool
        ↓
candidate improvement + evidence
        ↓
Origami validation
        ↓
Origami canonical version/profile
        ↓
optional Tonal composition with the chosen toolchain
```

Therefore:

```text
TONAL COMPOSITION != ORIGAMI PROFILE PROMOTION
TLALOC != REQUIRED ORIGAMI RUNTIME
ORIGAMI != TLALOC SUBCOMPONENT
```

## Current historical composition

The repository currently contains lock/proposal machinery created while Tonal was initially used to compose exact Tlaloc + Origami revisions. Those records remain valid historical/reproducibility artifacts; the project role is now explicitly generalized so additional development tools can be added without redefining Tonal.

Existing commands remain useful:

```bash
./scripts/component-status.sh
./scripts/update-component.sh <component> <revision>
./scripts/verify-stack.sh
./scripts/build-snapshot.sh
```

As Tonal becomes multi-tool rather than pair-specific, component registration should evolve from hard-coded assumptions toward a registry/manifest of independently versioned tools.

## Shared workflow

Tonal may host project-agnostic workflow skills, integration policies and snapshot/distribution logic that are useful across multiple development systems.

This does not transfer semantic ownership from component/target repositories to Tonal. A stack gate can prove that a selected composition works together; it does not rewrite the meaning or release policy of the composed projects.

## Provenance and Gatekeeper

Gatekeeper remains useful as a cross-repository provenance/integration policy. Technical gates must continue to run regardless of provenance classification.

Its authority in Tonal is about **composition and integration**, not about deciding semantic truth inside Origami or any other target project.

## Direction

Tonal should be able to express toolchains such as:

```text
Tonal composition A
  Tlaloc
  Origami

Tonal composition B
  Blueprint Framework
  Origami

Tonal composition C
  Tlaloc
  Blueprint Framework
  Origami

Tonal composition D
  Tlaloc
  some unrelated target/tool
```

Each component remains independently versioned and authoritative over its own contracts.
