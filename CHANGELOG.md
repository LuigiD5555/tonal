# Tonal changelog

## 0.1.0-alpha.5 — Generic Development-Tool Composition R1

- replaces the pair-specific composition model with `tonal.composition.v2` and `tonal.lock.v2`;
- introduces typed components: `development_tool`, `target`, `support`;
- records Tlaloc as a development tool and Origami as an independent target;
- treats Blueprint Framework as an unpinned example only, never fabricating a repository/version/SHA;
- makes component update, lock verification, declared verification and snapshot assembly generic over all locked components;
- moves component-specific deterministic checks into declarative `verification` entries in `TONAL.json`;
- distinguishes `COMPOSITION_VERIFIED` from target/model capability promotion;
- formalizes `TONAL_COMPOSITION != TARGET_PROMOTION` and `TONAL_IS_OPTIONAL_FOR_COMPONENT_RUNTIME`;
- preserves historical Fixed Carrier R2 candidate records/regression gates;
- updates the current Tlaloc development-tool pin to `6.0.0-alpha.14`; the Origami target pin remains on the last merged green release until the pending alpha.11 PR is merged.

## 0.1.0-alpha.4 — Candidate composition records

- records the merged Hybrid Receiver candidate composition without silently promoting unresolved real-model gates;
- adds the Fixed Carrier R2 candidate pairing Origami `6.0.0-alpha.5` at `cf6094cf9d3d9f636afae9bc62c15d063ad4fb3a` with Tlaloc `6.0.0-alpha.11` at `44dca10d1deb78446131a2de84ac37120081e5e0`;
- adds automated candidate identity/contract verification independent of the released `tonal.lock`;
- keeps the lock on the previously promoted stack until cross-model visual BOOT, real external tool-bridge, held-out multi-document and Native visual-machine gates are complete;
- preserves `SUPPORTED` as an explicit historical later promotion decision rather than equating component merge with stack support. Alpha.5 clarifies that this refers to composition/support records, not authority over Origami releases.

## 0.1.0-alpha.3 — Stack Definition R0

- formalized the then-current three-repository composition with Tonal as composition root;
- defined Tlaloc and Origami as independently authoritative for their own semantics;
- established Tonal as the single authority for project-agnostic `repo-flow` and composition/distribution policy;
- locked exact immutable component commits and required component CI before updating the composition;
- later alpha.5 generalizes the architecture beyond a fixed Tlaloc+Origami pair and removes any implication that Tonal owns target releases.

## 0.1.0-alpha.2 — Project Skill Distribution R0

- adds canonical `skills/repo-flow/SKILL.md` for project-agnostic Git/GitHub workflow discipline;
- adds byte-identical `.claude/skills/repo-flow/` and `.agents/skills/repo-flow/` mirrors;
- adds mirror synchronization and regression gates so agent-specific copies cannot silently drift;
- adds safe project installation into both supported layouts with local-edit protection and explicit `--force`;
- updates the tested Tonal composition to Tlaloc `6.0.0-alpha.9` at `6700b3bd9371a69ecc92a3f7bf95643d91f0f4ef`;
- keeps Origami `6.0.0-alpha.3` pinned at `978feef7f286cfe18b312ab8c833569094f12ef7`;
- does not change Tlaloc or Origami semantic ownership.

## 0.1.0-alpha.1 — Composition Contract R0

- establishes Tonal as a composition/distribution layer for independently versioned repositories;
- locks exact Tlaloc and Origami commits;
- introduces `TONAL.json` and immutable-release `tonal.lock` concepts;
- adds exact-commit fetch and identity verification scripts;
- adds CI gates for Tonal metadata plus locked component deterministic test closures;
- intentionally separates composition metadata from component semantic authority.
