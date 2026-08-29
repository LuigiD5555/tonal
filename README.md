# Tonal 0.1.0-alpha.3

Tonal is the **reproducible composition and distribution layer** for independently versioned Tlaloc and Origami releases.

Tonal does not own Tlaloc behavior semantics or Origami representation semantics. It resolves exact component revisions, verifies compatibility and integrity, carries stack-level project tooling, and is the place from which physical stack snapshots/releases are built.

## Current composition

| Component | Version | Exact revision |
|---|---|---|
| Tlaloc | `6.0.0-alpha.10` | `a20ab61d0043bc3e0a166f67207b01bdb2678a78` |
| Origami | `6.0.0-alpha.3` | `978feef7f286cfe18b312ab8c833569094f12ef7` |

The component repositories remain authoritative. Tonal pins immutable commits; it does not vendor or fork their source trees.

## Repository roles

```text
LuigiD5555/tlaloc   -> work/orchestration system
LuigiD5555/origami  -> representation/state-machine language
LuigiD5555/tonal    -> composition, compatibility, reproducibility and distribution
```

Within Tlaloc, Tlaloque are bounded specialist agents. Within Origami, OHF/R3.10-LAB is a nested research track rather than the identity of the entire language.

## Cross-project skills

Tonal is the **single canonical owner and distributor** of the project-agnostic `repo-flow` workflow skill.

The source of truth is:

```text
skills/repo-flow/SKILL.md
```

The repository keeps byte-identical mirrors for common agent layouts:

```text
.claude/skills/repo-flow/SKILL.md
.agents/skills/repo-flow/SKILL.md
```

Tlaloc `6.0.0-alpha.10` no longer contains a second `repo-flow` copy. Its own `.claude/skills/` directory now contains only the five Tlaloc/Origami-specific development skills. This prevents two active authorities for the same shared workflow asset.

Update the Tonal canonical copy and run:

```bash
./scripts/sync-skills.sh
./tests/test-skills.sh
```

To install the skill into both layouts of another Git project:

```bash
./scripts/install-skill.sh repo-flow --project /path/to/project
```

Existing differing copies are protected; use `--force` only after reviewing the local changes.

## Files

- `TONAL.json` — human/machine-readable composition policy and ownership contract.
- `tonal.lock` — exact resolved component revisions for this Tonal release.
- `VERSION` — Tonal's independent version.
- `skills/` — canonical stack-level project skills.
- `.claude/skills/` — verified Claude-compatible mirrors.
- `.agents/skills/` — verified agent-compatible mirrors.
- `scripts/sync-skills.sh` — regenerates skill mirrors from canonical sources.
- `scripts/install-skill.sh` — safely copies a Tonal skill into both supported project layouts.
- `scripts/verify-lock.sh` — verifies manifest/lock/version coherence.
- `scripts/fetch-components.sh` — fetches the exact locked commits into a temporary workspace.
- `scripts/verify-components.sh` — verifies fetched commit and component `VERSION` identities.
- `tests/test-ownership.sh` — verifies that shared `repo-flow` ownership is canonical in Tonal and absent from the locked Tlaloc payload.
- `.github/workflows/verify.yml` — composition, ownership and skill-coherence CI without external model campaigns.

## Development flow

Tonal releases are assembled from tested commits that already exist in the independent component repositories:

```text
Tlaloc main ---- exact commit ---\
                                 +--> Tonal lock --> verification --> snapshot/release
Origami main --- exact commit ---/
```

A released `tonal.lock` is immutable. To change one component, create a new Tonal version and a new lock resolution.

## Snapshot status

`0.1.0-alpha.3` closes **Repo-flow Ownership R0**: Tonal is the single canonical authority for the shared repository workflow while Tlaloc returns to component-specific skills only. Physical binary/source snapshot packaging and GitHub Release publication remain the next layer; this release does not claim a final portable stack installer yet.
