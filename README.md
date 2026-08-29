# Tonal 0.1.0-alpha.1

Tonal is the **reproducible composition and distribution layer** for independently versioned Tlaloc and Origami releases.

Tonal does not own Tlaloc behavior semantics or Origami representation semantics. It resolves exact component revisions, verifies compatibility and integrity, and is the place from which physical stack snapshots/releases are built.

## Current composition

| Component | Version | Exact revision |
|---|---|---|
| Tlaloc | `6.0.0-alpha.8` | `b222d8bef92a3da2a1753731d351cbe545ff7805` |
| Origami | `6.0.0-alpha.3` | `978feef7f286cfe18b312ab8c833569094f12ef7` |

The component repositories remain authoritative. Tonal pins immutable commits; it does not vendor or fork their source trees.

## Repository roles

```text
LuigiD5555/tlaloc   -> work/orchestration system
LuigiD5555/origami  -> representation/state-machine language
LuigiD5555/tonal    -> composition, compatibility, reproducibility and distribution
```

Within Tlaloc, Tlaloque are bounded specialist agents. Within Origami, OHF/R3.10-LAB is a nested research track rather than the identity of the entire language.

## Files

- `TONAL.json` — human/machine-readable composition policy and ownership contract.
- `tonal.lock` — exact resolved component revisions for this Tonal release.
- `VERSION` — Tonal's independent version.
- `scripts/verify-lock.sh` — verifies manifest/lock/version coherence.
- `scripts/fetch-components.sh` — fetches the exact locked commits into a temporary workspace.
- `scripts/verify-components.sh` — verifies fetched commit and component `VERSION` identities.
- `.github/workflows/verify.yml` — composition CI without external model campaigns.

## Development flow

Tonal releases are assembled from tested commits that already exist in the independent component repositories:

```text
Tlaloc main ---- exact commit ---\
                                 +--> Tonal lock --> verification --> snapshot/release
Origami main --- exact commit ---/
```

A released `tonal.lock` is immutable. To change one component, create a new Tonal version and a new lock resolution.

## Snapshot status

`0.1.0-alpha.1` establishes the **Composition Contract R0**. Physical binary/source snapshot packaging and GitHub Release publication are intentionally the next layer; this release does not pretend that a final portable installer already exists.
