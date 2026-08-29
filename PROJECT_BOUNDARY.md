# Tonal project boundary

The complete project is organized as three independently versioned repositories. **Tonal is the composition root**, while Tlaloc and Origami remain authoritative for their own semantics.

```text
TONAL — composition root
composition / lock / compatibility / integrity / distribution / repo-flow
    |
    +-- TLALOC — behavior and execution authority
    |   behavior compilation / swarm training / orchestration / Tlaloque / verification
    |
    +-- ORIGAMI — representation and machine authority
        visual/computational language / state machines / dynamics / evidence
        |
        +-- OHF
            nested carrier/laboratory research track
```

## Repository contract

| Repository | Owns | Must not own |
| --- | --- | --- |
| Tonal | exact component pins, compatibility, stack integrity, reproducible distribution, project-agnostic `repo-flow`, snapshots/releases | Tlaloc behavior semantics or Origami machine semantics |
| Tlaloc | BehaviorSpec/PromptIR behavior compilation, swarm orchestration/training, verification, Tlaloque execution, Tlaloc-specific skills | canonical `repo-flow`, Origami representation/state semantics |
| Origami | representation, state-machine laws, dynamics, perceptual contracts, deterministic experiment/reference evidence, OHF research track | Tlaloc orchestration/training semantics or stack composition policy |

## Promotion rule

A component repository develops and tests independently. Tonal only promotes a component after the component's required gates are green, then records its **exact immutable commit** in `TONAL.json` and `tonal.lock`. Floating branch names are never a released composition.

Changing one component does not automatically change the version of another component. Tonal receives its own version bump when the resolved stack, compatibility contract, distribution behavior, or release metadata changes.

## Workflow authority

`skills/repo-flow/SKILL.md` in Tonal is the canonical project-agnostic repository workflow. Tonal synchronizes its `.claude/skills/` and `.agents/skills/` mirrors. Tlaloc and Origami may own component-specific workflows, but must not fork the canonical project-agnostic `repo-flow` semantics.

## Snapshot boundary

A future Tonal snapshot is an immutable distribution artifact assembled from the exact locked component commits. It is not a monorepo and never becomes a new source of truth for Tlaloc or Origami.
