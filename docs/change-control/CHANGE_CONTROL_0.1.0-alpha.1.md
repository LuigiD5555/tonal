# Change control — Tonal 0.1.0-alpha.1

Date: 2026-08-28  
Status: `PROPOSED_COMPOSITION_CONTRACT_R0`

## Component changed

New Tonal composition repository.

## Before

Tlaloc and Origami were independently versioned and installable, but there was no Git-level object defining which exact pair of revisions constituted a reproducible combined stack.

## After

Tonal provides an independent composition version, policy manifest and exact lockfile. It can fetch and verify the locked revisions without submodules or vendoring.

## Locked inputs

- Tlaloc `6.0.0-alpha.8` — `b222d8bef92a3da2a1753731d351cbe545ff7805`
- Origami `6.0.0-alpha.3` — `978feef7f286cfe18b312ab8c833569094f12ef7`

## Invariants

- source authority remains in each component repository;
- Tonal pins exact commits, never floating `main` references;
- released lockfiles are immutable;
- Tonal cannot promote unsupported Tlaloc or Origami capabilities merely by packaging them;
- Tonal versioning is independent from both component version lines.

## Verification gates

- JSON parse and schema/identity sanity checks;
- `VERSION`, `TONAL.json` and `tonal.lock` coherence;
- exact commit format and repository identity checks;
- exact locked component checkout verification;
- component `VERSION` verification;
- `go test ./...` and `go vet ./...` for Origami;
- `go test ./...` and `go vet ./...` for Tlaloc Behavior Lab;
- no external model campaign in composition CI.

## Out of scope

Binary packaging, final Tonal installer/uninstaller, GitHub Release assets and project-local `tonal.lock` generation are not claimed by this revision.
