# Project Gatekeeper R0

Tonal owns the project-wide provenance policy used by Tonal, Tlaloc and Origami.

## Decision

Every pull request is classified before promotion:

- `OWNER`: PR author is `LuigiD5555` and the head repository is the same canonical repository. Normal CI, integrity and compatibility checks still run. The owner may intentionally override a promotion gate when necessary.
- `EXTERNAL`: every other provenance. CI/integrity/compatibility still run, owner approval is mandatory, and the contributor cannot use an override or auto-promote the change.

The distinction is about **promotion authority**, not code correctness. Owner code is not assumed correct and external code is not assumed wrong.

## Scope

The same policy applies to:

- `LuigiD5555/tonal`
- `LuigiD5555/tlaloc`
- `LuigiD5555/origami`

`gatekeeper.json` is canonical here. Component repositories carry a small machine-readable mirror with `authority: LuigiD5555/tonal` so their CI can enforce the same rule without requiring Tonal to be checked out.

## GitHub behavior

`.github/workflows/gatekeeper.yml` classifies PR provenance. External PRs fail the gate until `LuigiD5555` submits an `APPROVED` review. Review submission/dismissal reruns the gate automatically.

The existing component/stack CI remains separate and must still report its own result. Gatekeeper does not turn a failing test into a passing test.

## Owner override

An override is an explicit promotion decision by the owner, not an automatic bypass. Record the reason in the PR/merge context when a required technical gate is knowingly bypassed. Never rewrite evidence to make the bypass look like a PASS.

## Enforcement boundary

The workflow computes and reports the policy. To make it impossible to merge an external PR while the gate is red, configure the repository ruleset/branch protection to require both the normal CI check and `gatekeeper / provenance`. Until that GitHub repository setting is enabled, a repository administrator can still manually bypass workflow results.

## Simple usage

For owner work: open/update the PR normally; gatekeeper recognizes the canonical owner path automatically.

For external work: contributor opens a PR -> normal CI runs -> gatekeeper waits for owner approval -> owner reviews -> gatekeeper reruns -> promotion can proceed only after the technical gates and approval are satisfied.
