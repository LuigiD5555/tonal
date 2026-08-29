# Tonal changelog

## 0.1.0-alpha.1 — Composition Contract R0

- establishes Tonal as the composition/distribution layer above independently versioned Tlaloc and Origami repositories;
- locks Tlaloc `6.0.0-alpha.8` at commit `b222d8bef92a3da2a1753731d351cbe545ff7805`;
- locks Origami `6.0.0-alpha.3` at commit `978feef7f286cfe18b312ab8c833569094f12ef7`;
- introduces `TONAL.json` and immutable-release `tonal.lock` concepts;
- adds exact-commit fetch and identity verification scripts;
- adds CI gates for Tonal metadata plus the locked components' deterministic Go test/vet closures;
- intentionally defers physical snapshot packaging/installer publication to the next layer rather than claiming an unimplemented distribution artifact.
