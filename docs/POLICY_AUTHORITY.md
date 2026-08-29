# Policy authority

Tonal is the single authority for project-agnostic repository workflow and provenance/promotion policy. Component repositories may carry executable mirrors needed for independent CI, but those mirrors must identify Tonal as authority and must not evolve independently.

Tlaloc owns behavior/orchestration semantics. Origami owns representation/state/observation semantics. Gatekeeper does not transfer either semantic authority to Tonal.
