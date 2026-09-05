# T3 — Selective / Frugal Compute

**Status:** planned after Primitive Swarm foundation.

## Question

Can Tonal maintain useful quality while activating only the minimum sufficient capabilities and escalating only when cheaper paths are insufficient?

## Candidate policy ladder

```text
known deterministic Machine
        ↓ insufficient
single specialized Tlaloque
        ↓ uncertain/disagreement
small quorum + verifier
        ↓ unresolved
stronger/global probabilistic capability
        ↓ unresolved
human / explicit failure
```

## Mechanisms to test

- cost-aware cascades;
- reliability-aware routing;
- disagreement-triggered escalation;
- selective self-consistency/quorum;
- DAG parallel execution;
- explicit resource budgets.

RouteLLM, FrugalGPT, Self-Consistency and LLMCompiler motivate these system-level mechanisms. Tonal should first test deterministic/simple policies before learned routing.
