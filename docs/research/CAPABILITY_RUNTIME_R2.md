# Capability runtime target — R2

This document defines the first code-facing architectural target.

## Problem

The runtime must select among capabilities, not treat Parrot as the implicit default implementation of cognition.

## Target separation

```text
CapabilityRef
  identity + declared behavior

CapabilityProfile
  evidence + competence + cost/reliability observations

Executor
  performs bounded work

Verifier
  independently evaluates result/evidence

SelectionPolicy
  selects among eligible capabilities
```

## Compatibility-first rule

The first implementation must preserve existing T1 behavior through a compatibility policy before adaptive routing is introduced.

## Required properties

- deterministic executors and Parrot use the same selection surface;
- Parrot may be absent from workflows that do not need it;
- selection policy is independently replaceable;
- verification authority remains independent from executor authority;
- execution traces remain compatible with Episode generation.
