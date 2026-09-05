# TONAL T1 — Terminal-output / Fact semantics (FROZEN)

Status: frozen for the T1 pre-inference package. Do not revise after the
first T1 inference call.

## Rule

1. **VERIFY is the only Tlaloque that may promote a Fact.** The TONAL engine
   enforces this as a hard invariant: if any non-VERIFY node emits an
   Observation with `Kind == "FACT"`, the run terminates with
   `FinalStatus = CONTRACT_FAILURE` and
   `Error = FACT_PROMOTION_SCOPE_VIOLATION: <capability> node <local_id> …`.
   (`runtime/tonal/engine.go`, in the per-step observation loop.)

2. **Shapes 1–2 contain no VERIFY.** `READ_AND_CHECK` and
   `COMPARE_TWO_VALUES` terminate on a `COMPARE_NUMBERS` observation.
   Therefore their terminal output is an **evaluable workflow result**, not
   a promoted Fact.

3. **Shapes 3–5 terminate on VERIFY** and their terminal output is a
   **promoted Fact** (`blackboard.Fact`, `Status ∈ {VERIFIED, UNSUPPORTED}`).

4. Every `RunRecord` carries `terminal_output_kind ∈
   {evaluable_terminal_output, promoted_fact}`, set from
   `TaskFamily.HasVerify()` before execution.

## Scorer contract

- For `terminal_output_kind == evaluable_terminal_output` (Shapes 1–2) the
  scorer reads `RunRecord.FinalValue` directly. It is the last node's
  Observation; no Fact promotion is expected or required, and its absence
  is **not** a failure mode for these shapes.
- For `terminal_output_kind == promoted_fact` (Shapes 3–5) the scorer reads
  the VERIFY node's Fact. `Status == UNSUPPORTED` is a real workflow outcome
  (recorded, counts against `contract_success`), not a crash.
- No other node is permitted to write `Fact`; a run where one does is scored
  as `CONTRACT_FAILURE` with failure class `verification` (the scope
  invariant fired) — see the protocol document's failure taxonomy.

## Protocol limitation (recorded)

The frozen `exocortex.VerifyTlaloque` contract verifies a **single**
`target_key` against `expected_type ∈ {number,text,choice}` plus optional
range/allowed-choice constraints. The multi-argument
`VERIFY(disagreement_fraction, tolerance_margin, cmp_zero)` form sketched in
the driving prompt is **not expressible** without changing the frozen
`tlaloquekit` surface, which T1 does not do.

Resolution used by Shapes 3–5: the terminal VERIFY promotes the single
numeric observation that is transitively downstream of every required
observation (`norm_diff` for Shape 3, `norm_ratio` for Shape 4,
`norm_margin` for Shape 5). Additional required predecessors (e.g. Shape 5's
`cmp_zero` within-tolerance sign) are declared as `DependsOn` of the VERIFY
node so they remain transitively necessary and are computed and left on the
Blackboard for section-7 analysis, even though VERIFY reads only its target.
