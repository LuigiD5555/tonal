# Tonal changelog

## 0.1.0-alpha.2 — Project Skill Distribution R0

- adds canonical `skills/repo-flow/SKILL.md` for project-agnostic Git/GitHub workflow discipline;
- adds byte-identical `.claude/skills/repo-flow/` and `.agents/skills/repo-flow/` mirrors;
- adds mirror synchronization and regression gates so agent-specific copies cannot silently drift;
- adds safe project installation into both supported layouts with local-edit protection and explicit `--force`;
- updates the tested Tonal composition to Tlaloc `6.0.0-alpha.9` at `6700b3bd9371a69ecc92a3f7bf95643d91f0f4ef`;
- keeps Origami `6.0.0-alpha.3` pinned at `978feef7f286cfe18b312ab8c833569094f12ef7`;
- does not change Tlaloc or Origami semantic ownership and still does not claim a physical stack snapshot artifact.

## 0.1.0-alpha.1 — Composition Contract R0

- establishes Tonal as the composition/distribution layer above independently versioned Tlaloc and Origami repositories;
- locks Tlaloc `6.0.0-alpha.8` at commit `b222d8bef92a3da2a1753731d351cbe545ff7805`;
- locks Origami `6.0.0-alpha.3` at commit `978feef7f286cfe18b312ab8c833569094f12ef7`;
- introduces `TONAL.json` and immutable-release `tonal.lock` concepts;
- adds exact-commit fetch and identity verification scripts;
- adds CI gates for Tonal metadata plus the locked components' deterministic Go test/vet closures;
- intentionally defers physical snapshot packaging/installer publication to the next layer rather than claiming an unimplemented distribution artifact.
