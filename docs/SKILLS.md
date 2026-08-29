# Shared skills

Tonal owns project-agnostic operational skills.

## repo-flow
Use for repository inspection, branches, coherent commits, PRs, CI, merges, releases and multi-repository composition.

## gatekeeper
Use when provenance or promotion authority matters: OWNER vs EXTERNAL, approval requirements, explicit owner override and propagation into Tonal.

Canonical sources are under `skills/`. `.claude/skills/` and `.agents/skills/` are mirrors verified by tests. Use `scripts/sync-skills.sh` after canonical edits and `scripts/install-skill.sh <name> --project <path>` for installation.

Component-specific semantic/behavior skills stay in their owning repositories; do not move them into Tonal merely for convenience.
