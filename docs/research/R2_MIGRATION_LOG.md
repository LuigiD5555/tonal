# R2 migration log

## Started

Architecture R2 migration started from:

- Tonal: `claude/tonal-t1`
- Tlaloc: `main`
- Origami: `main`

Migration branch in all three repositories:

`refactor/r2-foundation`

## Protected invariant

Frozen T1 experiment semantics and artifacts are not to be modified by R2 research ideas.

## Completed foundation steps

- established Tonal R2 `CLAUDE.md`;
- established Tonal current state, architecture, boundaries and research program;
- established paper map focused on experimentally demonstrated mechanisms;
- established T2 Primitive Swarm research scaffold;
- aligned Tlaloc `CLAUDE.md` and nomenclature with capability-foundry role;
- documented Tlaloc role in Tonal and Cognitive JIT direction;
- established Origami `CLAUDE.md`, role in Tonal, Shponglese carrier separation and anti-prior direction;
- began archive areas for superseded documentation;
- removed superseded Origami `PROJECT_BOUNDARY.md` from repository root after recording its pre-R2 status.

## Next implementation step

Inspect the T1 runtime types and introduce the smallest compatible capability/selection abstraction. Preserve T1 routing through a compatibility policy before adding any adaptive policy.
