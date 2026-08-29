# Project workflow

## Change a component

1. Work in Tlaloc or Origami on a dedicated branch.
2. Open a PR. Gatekeeper classifies OWNER vs EXTERNAL.
3. Run the component's technical CI/evidence.
4. For EXTERNAL, owner approval is additionally required.
5. Promote the component according to the resulting policy decision.

## Spread the promoted change

1. In Tonal run `scripts/component-status.sh`.
2. Run `scripts/update-component.sh <tlaloc|origami> <revision>`.
3. Run `scripts/verify-stack.sh`.
4. Run `scripts/build-snapshot.sh`.
5. Promote the Tonal composition only if the exact pinned combination has the required evidence.

## Shared workflow changes

Project-agnostic workflow policy belongs in Tonal. Edit the canonical skill under `skills/`, run `scripts/sync-skills.sh`, run stack verification, then distribute/install it. Do not independently fork shared policy inside Tlaloc or Origami.

## Component-specific guidance

Tlaloc-specific behavior/orchestration skills remain in Tlaloc. Origami semantic documentation remains in Origami. Tonal should reference those authorities rather than duplicating their semantics.
