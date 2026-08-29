# Change control — Tonal 0.1.0-alpha.2

Date: 2026-08-28  
Status: `CANDIDATE_PROJECT_SKILL_DISTRIBUTION_R0`

## Component changed

Tonal stack-level project-skill distribution and the exact tested Tlaloc composition pin.

## Before

Tonal `0.1.0-alpha.1` defined exact Tlaloc/Origami composition and CI, but it did not carry a reusable project-agnostic workflow skill. `repo-flow` existed in Tlaloc `6.0.0-alpha.9` only under `.claude/skills/`, so Tonal had no neutral canonical location or `.agents/skills/` mirror.

## After

- Tonal advances to `0.1.0-alpha.2`.
- `skills/repo-flow/SKILL.md` is the canonical Tonal distribution copy.
- `.claude/skills/repo-flow/` and `.agents/skills/repo-flow/` are byte-identical mirrors.
- `scripts/sync-skills.sh` regenerates both mirrors from canonical skill sources.
- `scripts/install-skill.sh` installs a skill into both supported layouts of another Git repository and refuses to overwrite differing content without explicit `--force`.
- `tests/test-skills.sh` verifies source/mirror equality, idempotent installation, local-edit protection and forced replacement.
- Tonal now locks Tlaloc `6.0.0-alpha.9` at `6700b3bd9371a69ecc92a3f7bf95643d91f0f4ef`.
- Origami remains `6.0.0-alpha.3` at `978feef7f286cfe18b312ab8c833569094f12ef7`.

## Ownership

Tonal owns distribution/synchronization of this project-agnostic stack workflow asset. It does not thereby own Tlaloc-specific architecture skills, Origami semantics, BehaviorSpec/PromptIR/Tlaloque semantics, or component source history.

The existing Tlaloc copy of `repo-flow` is not deleted by this release; cross-repository ownership cleanup can be performed separately after Tonal distribution is established.

## Verification gates

- `VERSION`, `TONAL.json` and `tonal.lock` coherence;
- JSON/schema/commit-format validation;
- canonical `repo-flow` vs `.claude` mirror equality;
- canonical `repo-flow` vs `.agents` mirror equality;
- project install/idempotence/local-edit/`--force` regression test;
- exact Tlaloc/Origami checkout and component `VERSION` verification;
- `go test ./...` and `go vet ./...` for Origami;
- `go test ./...` and `go vet ./...` for Tlaloc Behavior Lab;
- no external model campaign.

## Regression risk

Low. The skill and its installer are additive to Tonal. The composition pin moves only Tlaloc from alpha.8 to the already-promoted alpha.9 release; Origami remains unchanged. No component semantics are modified.

## Downstream impact

Future Tonal snapshots can include one project-agnostic repository workflow source and expose it to both Claude-style and generic agent-style project layouts without requiring submodules or duplicate hand-maintenance.

## Promotion decision

Promote only after the pull-request CI verifies the exact locked components and all Tonal skill/coherence gates.
