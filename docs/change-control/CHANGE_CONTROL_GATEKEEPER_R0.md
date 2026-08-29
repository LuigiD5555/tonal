# Change Control — Project Gatekeeper R0

## Scope
Project-wide provenance and promotion authority for Tonal, Tlaloc and Origami.

## Before
Technical CI existed, but owner-authored and external contributions had no shared machine-readable promotion policy.

## After
- Tonal owns `gatekeeper.json` and `GATEKEEPER.md`.
- All three repositories execute a `gatekeeper / provenance` workflow.
- OWNER is the canonical `LuigiD5555` PR path; technical CI still runs and explicit owner override remains possible.
- EXTERNAL requires technical CI plus an APPROVED owner review and has no override/auto-promotion authority.
- Tonal distributes a canonical `gatekeeper` skill with Claude/Agents mirrors and installer coherence tests.

## Boundary
This change does not alter Tlaloc behavior semantics, Origami representation semantics, or existing technical PASS/FAIL evidence. Gatekeeper governs promotion authority only.

## Enforcement note
Hard prevention of administrator bypass requires GitHub repository rules/branch protection to mark normal CI and `gatekeeper / provenance` as required checks.
