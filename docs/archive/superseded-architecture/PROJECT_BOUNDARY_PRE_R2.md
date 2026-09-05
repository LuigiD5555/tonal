# Tonal project boundary

Tonal is an **optional composition and reproducibility layer** for independently versioned development tools, targets and support components.

It is not the semantic root of the projects it composes.

```text
                    TONAL (optional)
       exact pins / compatibility / provenance / snapshot
            /                 |                 \
       TLALOC          Blueprint Framework     other tools
 development_tool       development_tool        ...
            \                 |                 /
             \------ development ecosystem ---/
                            |
                            v
                         ORIGAMI
                           target
                 owns its own releases
```

Blueprint Framework is an architectural example here, **not a currently locked Tonal component** unless an exact repository/version/commit is explicitly added.

## Component kinds

Tonal composition v2 supports:

```text
development_tool
target
support
```

Every actual component in a released lock requires:

```text
name
kind
repository
version
immutable commit SHA
```

Unpinned examples are not components.

## Authority

| Kind | Owns | Tonal may do | Tonal must not do |
| --- | --- | --- | --- |
| development tool | its own development semantics and artifacts | pin, verify, compose, snapshot | redefine the tool |
| target | its own semantics, releases and accepted artifacts | pin, run declared integration checks, record provenance | promote/redefine target semantics |
| support | its own contract | pin and verify declared integration | silently become target authority |

For the current composition:

- **Tlaloc** is a `development_tool`: behavioral discovery, bounded Tlaloque reference swarms, prompt-first distillation and target-specific laboratories.
- **Origami** is a `target`: representation/state-machine/visual language and owner of Origami semantics/profiles/releases.

## Promotion terminology

Tonal can promote a **composition release** in the sense that a new Tonal lock/snapshot is verified and published.

That is not the same as promoting a target project's semantics or capabilities.

```text
TONAL COMPOSITION RELEASE
        !=
ORIGAMI PROFILE/SEMANTIC PROMOTION
```

Origami decides what is canonical Origami. Tlaloc decides what is a Tlaloc release. Tonal records exact compatible/reproducible combinations.

## Verification

Each component may declare deterministic verification commands in `TONAL.json`. Tonal executes those commands against the exact locked checkout.

This allows future tools to be added without hard-coding their test layout into Tonal itself.

A green Tonal verification means:

> the declared composition, pins and configured checks are reproducible and pass.

It does not mean that every empirical/model-facing capability of every target is universally supported.

## Workflow authority

Tonal remains the canonical distribution owner for project-agnostic `repo-flow` / shared Gatekeeper workflow assets where those are used.

That workflow ownership is independent from target semantics.

## Snapshot boundary

A Tonal snapshot contains the exact locked components plus composition metadata/shared workflow assets. It is a reproducible distribution, not a monorepo and not a new source of truth for the component projects.

## Hard rules

```text
COMPONENT_REPOSITORIES_REMAIN_AUTHORITATIVE
TARGET_REPOSITORY_OWNS_TARGET_RELEASES
TONAL_COMPOSITION != TARGET_PROMOTION
TONAL_IS_OPTIONAL_FOR_COMPONENT_RUNTIME
UNPINNED_TOOL != COMPOSITION_COMPONENT
SNAPSHOT_REFERENCES_EXACT_COMMITS
DECLARED_COMPONENT_VERIFICATION_MUST_PASS
```

---

> Archived under Architecture R2. This exact pre-R2 definition is preserved for provenance and is not current authority.
