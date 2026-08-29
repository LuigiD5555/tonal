---
name: gatekeeper
description: Use this skill when deciding whether a Tonal, Tlaloc, or Origami change is owner-authored or external, what promotion gates apply, whether owner approval is required, and whether an explicit owner override is allowed.
version: 0.1.0
---

# Project Gatekeeper

Use the same provenance decision across all three repositories.

## Fast path

1. Identify the canonical repository, PR author and PR head repository.
2. Classify `OWNER` only when the author is `LuigiD5555` and the head repository is that same canonical repository.
3. Classify everything else `EXTERNAL`.
4. Always run the repository's normal technical CI/evidence gates.
5. For `OWNER`, promotion may proceed normally; an intentional owner override is allowed but must never be represented as a technical PASS.
6. For `EXTERNAL`, require an `APPROVED` review from `LuigiD5555`; never grant override or auto-promotion authority.
7. When a component is promoted, use Tonal to update the exact pin and verify the complete stack before distributing a new composition.

## Authority

`Tonal/gatekeeper.json` and `Tonal/GATEKEEPER.md` are the project-wide authority. Tlaloc and Origami contain local mirrors so CI can classify provenance without fetching Tonal.

## Important distinction

Provenance answers **who may promote**, not **whether code is correct**. OWNER does not mean tests may be silently ignored. EXTERNAL does not mean the contribution is untrusted code; it means owner approval is required before project promotion.

## Expected outcomes

- `OWNER + PASS` -> normal promotion.
- `OWNER + failing/unknown technical gate` -> stop by default; owner may explicitly override and record why.
- `EXTERNAL + PASS + owner approval` -> eligible for promotion.
- `EXTERNAL + no owner approval` -> deny promotion.
- `EXTERNAL + failing technical gate` -> deny promotion.

## Multi-repository closure

A component merge is not automatically a Tonal release. After Tlaloc or Origami changes, Tonal must resolve the exact component commit, update its lock/manifest, execute full-stack verification, retain evidence, and build the reproducible snapshot.

## Do not

Do not infer OWNER from commit email/name alone. Do not allow a fork PR to become OWNER merely because its code looks familiar. Do not rewrite failed evidence as green. Do not copy project-wide policy into a second independent authority; component mirrors must continue pointing to Tonal.
