# R2 Foundation implementation slice

## Objective

Make Parrot one interchangeable capability by separating capability identity, execution, verification and selection policy without changing frozen T1 experiment semantics.

## First code targets

Introduce or stabilize the following concepts around the existing T1 runtime:

```text
CapabilityRef / CapabilityContract
CapabilityProfile
Executor
Verifier
SelectionPolicy
Episode adapter
```

Do not create an adaptive learned router yet.

## Compatibility requirement

The first policy implementation should reproduce existing T1 routing decisions. A regression test should compare current T1 workflow resolution with the generic policy path.

Conceptual target:

```text
SelectionPolicy
  ├── FrozenT1Policy
  ├── future StaticPolicy
  ├── future CostAwarePolicy
  ├── future ReliabilityAwarePolicy
  └── future LearnedPolicy
```

## Architectural tests

At minimum prove:

1. a deterministic capability can be selected instead of Parrot without changing the engine;
2. Parrot can be absent from the Registry for workflows that do not require its capabilities;
3. SelectionPolicy can be replaced independently of executors;
4. verifier authority is separate from executor authority;
5. execution can emit Episode-compatible evidence.

## Non-goals

- no changes to frozen T1 corpus/workflows/accounting;
- no learned router;
- no autonomous capability synthesis;
- no Shponglese generation;
- no Origami runtime dependency.
