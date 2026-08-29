# Tonal Hybrid Receiver Composition R0

Status: **COMPONENTS MERGED / MODEL EVIDENCE PENDING / NOT SUPPORTED**

## Current immutable component revisions

The implementation branches have been integrated independently into their component repositories:

- Tlaloc: `625a254748fab43af29e4f3d44b5a7d2e3979ae2`
- Origami: `f18f0449320d676588d25c135c5c348fb32370d2`

These are immutable merged commits and replace feature-branch references in this composition record.

## What Tonal coordinates

Tonal does not own receiver semantics or swarm behavior. It coordinates an exact tested composition of independently versioned Tlaloc and Origami revisions.

```text
TLALOC
swarm/Tlaloque search
  -> receiver candidate
  -> prompt + BOOT/Rosetta strategy + MicroAgent IR
  -> tournament/evidence
             |
             v
ORIGAMI
semantic validation
  -> self-boot carrier contract
  -> deterministic execution
  -> receiver artifact registry
             |
             v
TONAL
immutable component identities
  -> cross-component gates
  -> compatibility state
```

## Primary runtime mode

Hybrid is the primary target. Native and Computational remain diagnostic baselines.

## Completed integration gates

- Origami deterministic CI and Gatekeeper: PASS before merge.
- Tlaloc verify and Gatekeeper: PASS before merge.
- Both component PRs merged independently.
- Component identities are now immutable merged commits.

## Remaining evidence gates

The composition is intentionally **not** marked SUPPORTED yet. Remaining empirical gates require a real model-facing campaign:

1. `PEAK_ACTIVE_TOKEN_EQ <= 4000` during the actual model-facing run.
2. Synthetic Hybrid end-to-end with a real target model.
3. Carrier-local symbol permutation with the same external Master Prompt.
4. Cross-model symbol-permutation replication.
5. No hidden private source/envelope on the answering side.
6. `FALSE_EXACT = 0`.
7. Origami-owned receiver promotion after evidence.

## Hard anti-shortcut rules

Tonal rejects a SUPPORTED claim when hidden source material is exposed, Rosetta mappings are smuggled into a supposedly universal prompt, false exactness is nonzero, active context exceeds W, global scans are hidden, or promotion authority leaves Origami.

## Current lock state

`tonal.lock` remains on the previously promoted stack. The new merged Hybrid Receiver revisions are recorded here as the next candidate composition, but are not promoted into the release lock until the remaining real-model evidence gates pass. This prevents a merge from being confused with empirical support.
