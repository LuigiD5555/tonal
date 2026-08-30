# Tonal changelog

## 0.1.0-alpha.6 — Origami Protocol R0 composition

- pins Tlaloc `6.0.0-alpha.16` at `7f3393079d163131ad690f97de59ab2e2a249179`;
- pins Origami `6.0.0-alpha.13` at `7dbd7ba073b227377c6cc3ae592f4c4f2573dabf`;
- records deterministic Origami Protocol R0 support artifacts: S*/E* codec registry, capability negotiation, Master Prompt R4 and rendered profile-3;
- records deterministic Tlaloc protocol evaluation: READ/WRITE/ROUNDTRIP/MULTIHOP, codec discovery, semantic preservation/drift, invention, external-codec dependency and semantic-to-exact escalation;
- records component-green profile-3 640x640 / 8192-byte renderer, profile-1/profile-2 decode non-regression and `S2(E2(INDEX))` roundtrip evidence;
- adds explicit separation `COMPOSITION_VERIFIED != PROTOCOL_INTEROPERABILITY_PROMOTED != NATIVE_SEMANTIC_PROMOTED`;
- keeps held-out Native S2, Native E2 and A -> B -> C real-model evidence pending and component-owned;
- does not make Tlaloc or Tonal runtime requirements for Origami;
- preserves historical Fixed Carrier R2 proposal/regression metadata independently of the new protocol composition.

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
- pins Tlaloc `6.0.0-alpha.15` at `349deaef6f6c94966c814f9f99a1ec1fb78b875f`, including Prompt-First Distillation and failure-driven Native semantic regression tooling;
- pins Origami `6.0.0-alpha.12` at `12176e6829e8cd0aaaa1db03f0b2cc4f4d5ea838`, including semantic-first T2 navigation and the preserved failed index regression;
- explicitly keeps Native semantic empirical promotion false: a reproducibly green composition is not evidence that held-out VLMs now pass the corrected index test.

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
