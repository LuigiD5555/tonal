# T6 — Origami carrier / memory

**Status:** planned.

## Question

Can Origami represent, transport, address or selectively unfold the same semantic program more efficiently or usefully than conventional codecs without relying on language-model priors?

## Fair comparison

Hold semantics constant:

```text
same Shponglese program
  ├── JSON
  ├── compact text
  ├── binary
  └── Origami
```

## Candidate metrics

- exact semantic recovery;
- downstream deterministic execution success;
- bytes/tokens transported;
- latency/cost;
- selective reads/unfolds;
- false-known rate;
- robustness under corruption;
- performance under randomized symbol mappings.

DeepSeek-OCR motivates compact visual transport experiments; critical follow-up findings motivate explicit anti-prior controls.
