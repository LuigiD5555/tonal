---
name: repo-flow
description: Use this skill for Git/GitHub repository work: inspect state, create branches, make coherent commits, open or repair pull requests, resolve conflicts, run CI gates, merge safely, verify releases, and compose multiple repositories into tested snapshots.
version: 0.2.0
---

# Repo Flow

Use a repository workflow that preserves user work, keeps history understandable, and makes promotion evidence explicit. The goal is not Git ceremony; it is a change that can be inspected, tested, reviewed, reproduced, merged, and recovered safely.

## Core rules

1. Inspect before changing: repository root, current branch, HEAD, base/upstream, dirty and untracked files, recent relevant history, and project-local instructions.
2. Never discard unknown work. Do not reset, clean, checkout over, stash, amend, rebase, force-push, or delete user changes unless that exact action is understood and explicitly intended.
3. Respect source-of-truth boundaries. Documentation, packaging, composition, generated artifacts, and runtime semantics are different authorities.
4. Make the smallest coherent change, including its consistency closure but excluding unrelated edits.
5. Test the impact closure. Run cheap deterministic checks first, affected tests next, and required release/CI gates last.
6. No false green. Pending, failing, skipped, or unknown required gates do not count as PASS.
7. Resolve conflicts semantically. `ours`, `theirs`, `current`, `incoming`, and `accept both` are implementation choices, not correctness criteria.
8. Verify after merge directly on the target branch.

## 1. Preflight

Establish a baseline equivalent to:

```bash
git rev-parse --show-toplevel
git status --short --branch
git rev-parse HEAD
git remote -v
git log --oneline --decorate -n 12
```

Read applicable repository guidance when present: `CLAUDE.md`, `AGENTS.md`, `CONTRIBUTING.md`, `README.md`, `VERSION`, `CHANGELOG.md`, change-control records, CI workflows, manifests, locks, and machine-readable state.

If unrelated work is present, preserve it and limit the edit scope rather than hiding or deleting it.

## 2. Define the change

Before editing, identify:

```text
component
before
intended after
reason
expected files
impact closure
tests/gates
version impact
promotion condition
```

Do not bump a neighboring repository merely because another component changed. Independently versioned projects remain independent.

## 3. Branch and implement

Normally use a dedicated branch based on the intended target branch:

```text
feature/<topic>
fix/<topic>
docs/<topic>
chore/<topic>
release/<version-or-topic>
```

Change only the declared scope and consistency closure. Examples:

- version change -> canonical version source, README, changelog, installer/package metadata, coherence tests;
- rename -> active docs, CLI/help, tests, machine-readable state;
- historical relocation -> preserve history, update indexes, do not rewrite old evidence as current;
- public contract change -> update downstream compatibility claims without silently implementing unsupported runtime behavior.

Inspect the actual diff:

```bash
git diff --stat
git diff --check
git diff
```

Look for conflict markers, duplicated headings, stale version strings, accidental generated files, lost history, and unrelated formatting churn.

## 4. Verification

Typical order:

```text
syntax / parse
static consistency
unit / targeted tests
integration / affected regressions
required full gates
packaging or install roundtrip when relevant
```

Parse JSON/YAML rather than only reading it. Compare regenerated artifacts or hashes when applicable. A failed required gate blocks promotion until fixed or explicitly rejected/waived by project policy.

## 5. Commits

Before committing, verify that only intended files are staged, no conflict markers remain, the diff is coherent, and required local gates have results.

Prefer effect-oriented messages such as:

```text
fix: keep installer version aligned with VERSION
feat: add reusable repository workflow skill
docs: reconcile component scope
release: <project> <version> <summary>
```

Do not rewrite already-shared history unless explicitly desired and safe.

## 6. Pull requests

Compare the branch against its real base and inspect the changed-file set before opening or updating a PR.

A useful PR states:

```text
why the change exists
what changed
what deliberately did not change
tests/evidence
risk or migration impact
version/release impact
```

Update an existing PR instead of creating a duplicate unless replacement is deliberate.

## Conflict resolution

For each conflicted file determine:

1. What valid information exists only on each side?
2. Which side is authoritative for current state?
3. Is the file cumulative history, generated output, or current source of truth?
4. Can the correct result be regenerated instead of hand-merged?

Useful defaults:

- README/current architecture: keep current authority, then reintroduce still-valid detail.
- CHANGELOG/history: preserve both valid histories, reorder and deduplicate manually.
- Generated files: regenerate from source of truth.
- JSON/YAML/manifests/locks: reconcile facts and schema deliberately; never concatenate blindly.

`Accept both` is correct only when the resulting combined structure is itself valid. Re-read the whole file and rerun its gates after resolving.

## 7. CI and merge

A merge candidate requires:

```text
PR diff reviewed
mergeable state known
required CI complete and successful
no unresolved blocker
expected head commit unchanged
```

If CI fails, inspect and fix the cause before merging. When supported, guard the merge with the exact expected head SHA. Use the repository's configured merge strategy rather than imposing squash/rebase/merge globally.

## 8. Post-merge verification

After merge verify the authoritative branch directly:

```text
merge result
VERSION or equivalent
critical changed files
machine-readable state/lock
changelog/change-control entry
package or installer identity if affected
```

A successful GitHub merge response alone is not proof that all release invariants are correct.

## Multi-repository composition

Keep component repositories independently versioned and authoritative.

```text
repo A exact commit
repo B exact commit
integration gates
      ↓
immutable snapshot / lock
```

Exact Git pins or submodules are useful inside an integration/snapshot builder because they record exact component commits. Do not force application repositories to use submodules merely to consume tools when installed-tool + lockfile is simpler.

A snapshot is a distribution artifact, not a new source of truth. Fix bugs in the owning repository, verify a new composition, then publish a new snapshot. Record both compatible ranges and the exact tested set. Never silently mutate a published snapshot under the same identity.

## Release invariants

When a canonical `VERSION` or equivalent exists:

- prefer it as the single release-version source;
- derive installers/build scripts from it where practical;
- guard against duplicated hard-coded versions;
- update changelog/change-control in the same coherent release change;
- verify installed/package identity matches the declared release.

## Stop conditions

Stop and investigate or ask rather than guessing when unrelated work may be overwritten, the correct base is unclear, two incompatible authorities appear valid, required CI is pending/failing, release ownership/versioning is ambiguous, a destructive Git action would be required, or a snapshot points at an unmerged/mutable state when immutability is required.

## Fast path

```text
inspect -> branch -> edit -> diff -> targeted tests -> consistency gates
-> atomic commit -> PR -> inspect PR -> CI -> merge -> verify target branch
```

Use more ceremony only when the repository risk justifies it.
