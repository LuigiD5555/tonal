# Tonal 0.1.0-alpha.3

Tonal is the **composition, compatibility, evidence and distribution layer** for independently versioned Tlaloc and Origami.

Tonal does not own Tlaloc behavior semantics or Origami representation semantics. Component repositories remain authoritative; Tonal records immutable compatible combinations and builds reproducible stack snapshots from them.

## Current composition

| Component | Version | Exact revision |
|---|---|---|
| Tlaloc | `6.0.0-alpha.10` | `a20ab61d0043bc3e0a166f67207b01bdb2678a78` |
| Origami | `6.0.0-alpha.3` | `bd6e47979fcc8918cefe1302bd34e183d784a14a` |

## Repository roles

```text
LuigiD5555/tlaloc   -> behavior compilation, orchestration, training, verification, Tlaloque
LuigiD5555/origami  -> representation, state machines, dynamics, observation contracts, OHF research
LuigiD5555/tonal    -> exact composition, compatibility evidence, shared workflow and distribution
```

## Normal propagation flow

```text
component change
      |
      v
Tlaloc or Origami CI
      |
      v
promoted component commit
      |
      v
Tonal pin update
      |
      v
full-stack verification + evidence
      |
      v
reproducible snapshot
```

Useful commands:

```bash
./scripts/component-status.sh
./scripts/update-component.sh origami main
./scripts/update-component.sh tlaloc main
./scripts/verify-stack.sh
./scripts/build-snapshot.sh
```

`component-status.sh` detects drift between the Tonal lock and component `main` branches. `update-component.sh` resolves a requested revision to an immutable commit and updates both `TONAL.json` and `tonal.lock`. `verify-stack.sh` fetches the exact commits and runs component tests, vet and Origami deterministic evidence gates. `build-snapshot.sh` creates a normalized source archive plus SHA-256 checksum.

## Shared project workflow

Tonal is the sole authority for project-agnostic workflow skills. Canonical skills live under `skills/`; `.claude/skills/` and `.agents/skills/` are verified mirrors. Run `scripts/sync-skills.sh` after editing a canonical skill and `scripts/install-skill.sh <name> --project /path/to/project` to distribute it safely.

- `repo-flow` handles Git/GitHub change, CI, merge and multi-repository composition.
- `gatekeeper` handles provenance and promotion authority across Tonal, Tlaloc and Origami.

## Project Gatekeeper

The same provenance rule applies to all three repositories. `gatekeeper.json` and `GATEKEEPER.md` are canonical here; component repositories contain local mirrors for CI.

```text
OWNER = LuigiD5555 + canonical repository PR
  -> technical gates still run
  -> explicit owner override is permitted

EXTERNAL = every other PR provenance
  -> technical gates still run
  -> APPROVED review from LuigiD5555 is mandatory
  -> no override / no auto-promotion
```

The distinction controls **promotion authority**, not whether code is presumed correct. See `GATEKEEPER.md` and the `gatekeeper` skill for the operational procedure.

## Promotion rule

A component change is not part of a Tonal stack merely because it exists on component `main`. It becomes part of the stack only when a new exact pin passes the complete Tonal verification closure. Released lock files are immutable; later component changes require a new Tonal version/lock.

GitHub Actions runs stack verification and provenance classification automatically. For hard multi-user enforcement, configure repository rules to require both normal CI and `gatekeeper / provenance`; administrators can otherwise intentionally bypass workflow results.
