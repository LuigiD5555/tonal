# T4 — Cognitive JIT

**Status:** planned.

## Question

Can repeated verified Episodes be converted by Tlaloc into cheaper reusable capabilities that preserve held-out reliability and reduce future neural/model dependence?

## System loop

```text
Tonal execution
  ↓
Episode
  ↓
Tlaloc pattern discovery
  ↓
candidate reusable structure
  ↓
qualification / holdout / ablation
  ↓
Registry
  ↓
future Tonal execution
```

## Key measurements

- reuse rate;
- held-out reliability;
- model calls avoided;
- cost/latency reduction;
- abstraction complexity;
- failure-rate change;
- fraction of work handled without general-model inference.

DreamCoder, Stitch, Voyager, Distilling Step-by-Step and MDL motivate candidate mechanisms. They do not establish that this loop will work in Tonal.
