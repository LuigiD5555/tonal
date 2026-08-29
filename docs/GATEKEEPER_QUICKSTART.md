# Gatekeeper quickstart

## If the change is yours

Open the PR from a branch in the canonical repository under `LuigiD5555`. Gatekeeper classifies it `OWNER`. Run normal CI. Merge normally when green; if you intentionally override a promotion gate, record the reason and do not call the failed check a PASS.

## If the change is from someone else

Let normal CI run. Gatekeeper classifies it `EXTERNAL` and stays red until `LuigiD5555` submits an `APPROVED` review. The external author cannot override or auto-promote.

## If Tlaloc or Origami changed

After component promotion, update Tonal's exact pin and verify the whole composition:

```bash
./scripts/update-component.sh origami main   # or tlaloc
./scripts/verify-stack.sh
./scripts/build-snapshot.sh
```

## Skills

Install the shared operational guidance into any project:

```bash
./scripts/install-skill.sh repo-flow --project /path/to/project
./scripts/install-skill.sh gatekeeper --project /path/to/project
```

Tonal is the authority for both shared skills. Component-specific skills remain owned by their component repositories.
