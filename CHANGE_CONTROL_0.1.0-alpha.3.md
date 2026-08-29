# Change Control — Tonal 0.1.0-alpha.3

## Component
Tonal composition contract and complete-project repository boundary.

## Before
Tonal 0.1.0-alpha.2 pinned Tlaloc 6.0.0-alpha.9 and an older Origami 6.0.0-alpha.3 commit. Project-agnostic repo-flow was already canonical in Tonal, but the promoted composition had not yet incorporated Tlaloc's ownership cleanup or Origami's automated reference-evidence gate.

## After
Tonal 0.1.0-alpha.3 defines the complete project as Tonal + Tlaloc + Origami with independent versions and explicit authority boundaries. It pins Tlaloc 6.0.0-alpha.10 and the promoted Origami Reference Engine R0 commit by exact SHA.

## Evidence
- Tlaloc main: `a20ab61d0043bc3e0a166f67207b01bdb2678a78`, version `6.0.0-alpha.10`.
- Origami PR #7 CI run `33230458227`: SUCCESS before guarded squash merge.
- Origami promoted main commit: `bd6e47979fcc8918cefe1302bd34e183d784a14a`, version remains `6.0.0-alpha.3` because the experimental implementation is an alpha.3 development increment rather than a separately promoted component release.
- Origami CI now executes EXP-001, validates evidence integrity and retains the evidence artifact.

## Required Tonal gates
- manifest/lock coherence;
- canonical repo-flow mirror equality;
- exact component fetch and identity verification;
- locked Tlaloc deterministic Go test/vet closure;
- locked Origami deterministic Go test/vet closure.

## Promotion decision
PENDING until Tonal CI for this branch is green. Merge must use the expected PR head SHA. After merge, verify VERSION, TONAL.json and tonal.lock on main.

## Downstream impact
Future stack automation must read `tonal.lock` as the exact composition source. Component repositories remain authoritative and are never rewritten from Tonal.
