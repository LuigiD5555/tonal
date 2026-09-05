# Tonal research program — R2

## Research thesis

Tonal investigates whether a system can obtain complex and reliable behavior by composing small bounded capabilities, external structure, verification, memory and compiled experience rather than pushing all complexity into one general-purpose model.

A positive result must be measured at system level. Architectural resemblance to a paper is not evidence.

## Method

For every external research result:

1. record what changed in the paper;
2. record the measured gain;
3. record important limitations;
4. map the mechanism to a Tonal failure or hypothesis;
5. build the smallest falsifiable prototype;
6. compare against a frozen baseline;
7. promote only if the evidence survives holdout/ablation where applicable.

## T1 — Heterogeneous depth

**Status:** frozen.

Question:

> Does heterogeneous composition degrade less with workflow depth than monolithic Parrot or a decomposed-but-Parrot-centric workflow?

T1 must not be modified to absorb later research ideas.

## T2 — Primitive Swarm / MICRO-ISA

Question:

> Can a small library of bounded, verifiable primitives cover a much larger family of tasks through composition?

Primary inspirations by demonstrated mechanism:

- PAL — separate semantic decomposition from deterministic execution;
- Least-to-Most — solve complex tasks through simpler sequential subproblems;
- Decomposed Prompting — route subtasks to replaceable specialists/functions;
- ViperGPT / VisProg — compose existing primitives into programs instead of training one monolith;
- LLM+P — delegate search/planning to a purpose-built symbolic solver.

Initial target: a small operational vocabulary such as LOCATE, READ, EXTRACT, NORMALIZE, COMPARE, ADD, SUBTRACT, FILTER, MAP, REDUCE, BRANCH, VERIFY and RETURN.

## T3 — Frugal Selective Compute

Question:

> Can Tonal maintain quality while activating only the minimum sufficient set of capabilities?

Mechanisms to test:

- cheap-first cascades;
- cost/reliability-aware routing;
- disagreement-triggered escalation;
- selective quorum/self-consistency only on uncertain nodes;
- DAG parallelism where dependencies permit.

Relevant results include RouteLLM, FrugalGPT, Self-Consistency and LLMCompiler.

## T4 — Cognitive JIT / Tlaloc experience compiler

Question:

> Can repeated verified Episodes be converted into cheaper reusable capabilities without losing held-out reliability?

Candidate mechanisms:

- recurring trajectory discovery;
- abstraction/library learning;
- MDL/compression-based candidate scoring;
- skill reuse;
- distillation to smaller specialists;
- deterministic compilation when possible.

Relevant research includes DreamCoder, Stitch, Voyager, Distilling Step-by-Step and MDL.

## T5 — Shponglese

Question:

> Can the learned primitive/motif vocabulary generalize compositionally across unseen combinations and remain semantically invariant across codecs?

Shponglese begins as a semantic IR. Claims of language require evidence for systematicity, compositionality, slot generalization, held-out composition and codec invariance.

## T6 — Origami carrier / memory

Question:

> Can Origami carry, address or selectively unfold the same semantic program more efficiently than conventional representations without relying on model priors?

DeepSeek-OCR motivates testing compact visual transport; critical follow-up work motivates anti-prior controls using randomized or semantically corrupted symbol assignments.

Origami must be compared against conventional codecs on the same Shponglese semantics.

## Cross-cutting metrics

Depending on experiment, measure:

- task success/accuracy;
- degradation versus depth/complexity;
- model calls;
- capability activations;
- routing errors;
- verification failures;
- latency;
- cost;
- bytes/tokens transported;
- reuse rate;
- held-out reliability;
- fraction of work requiring neural inference.

## Long-term target

The strongest long-term hypothesis is:

> As verified experience accumulates, Tonal should be able to convert recurring uncertainty into reusable structure, reducing future neural dependence without reducing useful capability.

This is a project hypothesis, not a result established by the cited papers.
