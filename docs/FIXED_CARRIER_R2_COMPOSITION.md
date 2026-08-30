# Fixed Carrier R2 composition candidate

This Tonal composition pairs **Origami Fixed Carrier R2** with **Tlaloc Canonical Memory R2**. Origami owns the frozen visual control plane; Tlaloc owns canonicalization, proposal-only Tlaloque compilation, deterministic reduction, exact/CID/Merkle memory and recursive query orchestration.

## Immutable component pair

- Origami `6.0.0-alpha.5`: `cf6094cf9d3d9f636afae9bc62c15d063ad4fb3a`
- Tlaloc `6.0.0-alpha.11`: `44dca10d1deb78446131a2de84ac37120081e5e0`

Both component pull requests passed their normal CI and Gatekeeper workflows before merge. Tonal verifies these exact commits independently rather than following floating `main` branches.

## What this candidate proves

The deterministic composition has evidence for a fixed 8192-byte carrier profile, corpus-independent physical PNG size, exact address/CID/Merkle access, CanonicalState replay with CDR=1.0 on the reference corpus, and `FALSE_EXACT=0` on that gate. A DeepSeek session also demonstrated that the carrier's visual probe can be read from a real image input.

## What it does not yet prove

This candidate does **not** promote universal Native visual decoding or general cross-model support. Remaining empirical gates include cross-model T0/visual-probe readability, an external real-model tool-bridge end-to-end run, held-out multi-document routing, and Origami promotion of the Native visual-machine profile.

For that reason `tonal.lock` remains the previously promoted stack. `proposals/FIXED_CARRIER_R2.json` is the exact immutable R2 candidate record; a later Tonal release may update the lock only after those empirical gates pass.
