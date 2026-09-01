# Tonal capability status

This table is generated from `state/CLAIMS.json`. Edit the ledger and regenerate it with `python3 tools/claims.py generate`; do not edit the table directly.

<!-- BEGIN GENERATED CLAIMS TABLE: do not edit; run python3 tools/claims.py generate -->
| Claim | Statement | Status | Evidence | Version introduced | Last checked | Notes |
|---|---|---|---|---|---|---|
| `TONAL.COMPONENT.EXACT_IDENTITY_VERIFICATION` | Verify that each fetched component checkout matches the exact commit and component version declared by tonal.lock. | `designed` | — | `0.1.0-alpha.1` | 2026-09-01 | scripts/fetch-components.sh and scripts/verify-components.sh implement the mechanism, but the repository has no qualifying go test ./... evidence. |
| `TONAL.COMPONENT.REVISION_RESOLUTION` | Resolve a declared component revision to an immutable commit and update its version and pin coherently in the composition manifest and lock. | `designed` | — | `0.1.0-alpha.5` | 2026-09-01 | scripts/update-component.sh implements this path, but no test executed by go test ./... establishes the contract required for implemented status. |
| `TONAL.COMPOSITION.DECLARED_VERIFICATION` | Execute each component verification command declared in TONAL.json against the corresponding exact locked checkout. | `designed` | — | `0.1.0-alpha.5` | 2026-09-01 | scripts/verify-stack.sh implements this orchestration, but no qualifying go test ./... establishes the complete contract. |
| `TONAL.COMPOSITION.LOCK_VERSION_COHERENCE` | Reject disagreement between Tonal VERSION, TONAL.json, tonal.lock, component identity fields, and the allowed component-kind contract. | `designed` | — | `0.1.0-alpha.1` | 2026-09-01 | scripts/verify-lock.sh and tests/test-manifest.sh exercise part of this contract as shell tests, but Tonal has no go test ./... evidence and the current verifier does not yet cover README or the claims ledger. |
| `TONAL.DISTRIBUTION.REPO_FLOW` | Distribute repo-flow from one canonical Tonal source to byte-identical Claude and Agents mirrors with controlled project installation. | `designed` | — | `0.1.0-alpha.2` | 2026-09-01 | scripts/sync-skills.sh, scripts/install-skill.sh, and tests/test-skills.sh implement and shell-test the path, but no test executed by go test ./... satisfies the ledger promotion rule. |
| `TONAL.POLICY.GATEKEEPER_DISTRIBUTION` | Keep Gatekeeper policy authoritative in Tonal and distribute byte-identical workflow mirrors without transferring component semantic authority. | `designed` | — | `0.1.0-alpha.4` | 2026-09-01 | The policy, canonical skill, mirrors, and shell tests exist, but no qualifying go test ./... establishes this contract. |
| `TONAL.SNAPSHOT.REPRODUCIBLE_DISTRIBUTION` | Build a source snapshot containing all exact locked components plus Tonal composition metadata and shared workflow assets. | `designed` | — | `0.1.0-alpha.5` | 2026-09-01 | scripts/build-snapshot.sh and TONAL.json describe an implemented builder, but no qualifying go test ./... proves reproducibility of the produced archive. |
<!-- END GENERATED CLAIMS TABLE -->

## Interpretation

`designed` means the repository contains a documented design or implementation surface without the test evidence required by the current ledger policy. Shell or Python gates are recorded in claim notes but do not promote a claim while the policy requires a test executed by `go test ./...`.
