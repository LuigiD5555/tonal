# Paper map — mechanisms relevant to Tonal R2

This map is not a bibliography of things Tonal should imitate. Each paper belongs here only because it demonstrates a mechanism that may be tested against a Tonal/Tlaloc/Origami hypothesis.

For each paper, future updates should record:

```text
paper_result
what_changed
measured_gain
important_limitations
tonal_mapping
smallest_prototype
falsification_condition
status
```

## Complexity through primitives and decomposition

- PAL
- Least-to-Most Prompting
- Decomposed Prompting
- ViperGPT
- VisProg
- LLM+P

Primary question: can complex global tasks be solved by composing simpler bounded operations, delegating deterministic work outside a general model?

## Execution and composition

- LLMCompiler
- Options framework
- MAXQ

Primary question: can dependency-aware execution, macro-operations and local state abstraction reduce orchestration overhead while preserving correctness?

## Selective / frugal intelligence

- FrugalGPT
- RouteLLM
- Self-Consistency
- ACT / PonderNet
- BLT

Primary question: can Tonal spend expensive cognition only where cheaper capabilities or previous structure are insufficient?

## Verification

- Training Verifiers to Solve Math Word Problems
- Let's Verify Step by Step

Primary question: does explicit process/step verification improve failure localization, routing decisions and safe capability promotion compared with outcome-only success labels?

## Experience to reusable structure

- DreamCoder
- Stitch
- Voyager
- Distilling Step-by-Step
- MDL / library-learning work

Primary question: can repeated verified experience be converted into reusable capabilities that reduce future cost or neural dependence while surviving holdout tests?

## Representation and carrier

- DeepSeek-OCR
- DeepSeek-OCR 2
- critical work on linguistic-prior dependence in optical compression
- BLT

Primary question: can compact representations transport or expose useful semantic structure without relying on model priors, and can granularity adapt to what is actually uncertain?

## Architectural inspirations kept secondary

- HRM
- SpikingBrain
- MoE / sparse routing
- Coconut / latent reasoning work
- Titans / learned memory

These may inspire later mechanisms, but R2 should prefer smaller system-level prototypes before reproducing large neural architectures.

## Current priority order

1. frozen T1 heterogeneous-depth evidence;
2. Primitive Swarm / MICRO-ISA;
3. selective/frugal routing;
4. Episode -> reusable structure / Cognitive JIT;
5. Shponglese compositional IR;
6. Origami carrier/memory comparisons.
