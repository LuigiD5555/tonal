# Enforcement boundary

Gatekeeper R0 is executable policy, but GitHub repository settings determine whether a red check is mechanically unmergeable.

When multi-user collaboration begins, require these checks on `main` through GitHub rulesets/branch protection:

- normal repository CI (`verify`/`CI` as appropriate);
- `gatekeeper / provenance`.

Keep administrator bypass only if deliberate owner emergency override is desired. If bypass is used, preserve the failing/pending evidence and record the reason; do not mutate history to manufacture a green result.
