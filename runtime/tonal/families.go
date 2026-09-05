package tonal

// TONAL T1 FROZEN TASK FAMILIES (protocol sections 2-4).
//
// Five generic workflow templates whose primary difficulty axis is
// WORKFLOW_COMPOSITION. Every family here is executor-agnostic: it names
// capabilities, never workers; it carries no per-instance operand address,
// no expected answer, and no benchmark id. A concrete run binds operands
// only through Instance.Params (the natural-language locate questions, the
// store directory, the per-operand page image path, and the scalar
// threshold/tolerance).
//
// natural_depth is the family's declared logical-composition label from the
// approved T1 design (2/4/6/8/12). critical_path_depth is the mechanical
// CriticalPathDepth() over the resolved DAG counting EVERY node including
// LOCATE_REGION. They coincide for Shape 5 by the approved enumeration and
// differ by a small constant for Shapes 1-4 (LOCATE/READ acquisition nodes
// are counted mechanically but not in the logical label). Both numbers are
// recorded per workflow; section 7 analysis uses critical_path_depth.
//
// Fact semantics (section 4): VERIFY is the only capability that promotes a
// Fact. Shapes 1-2 contain no VERIFY, so their terminal observation is an
// evaluable workflow result that the scorer reads directly — it is NOT a
// promoted Fact. Shapes 3-5 terminate on VERIFY.

// Frozen family ids.
const (
	FamilyReadAndCheck         = "READ_AND_CHECK"
	FamilyCompareTwoValues     = "COMPARE_TWO_VALUES"
	FamilyDifferenceThenVerify = "DIFFERENCE_THEN_VERIFY"
	FamilyRatioOfDifference    = "RATIO_OF_DIFFERENCE"
	FamilyReconciliationChain  = "RECONCILIATION_CHAIN"
)

// NaturalDepth is the approved logical-composition label per family.
var NaturalDepth = map[string]int{
	FamilyReadAndCheck:         2,
	FamilyCompareTwoValues:     4,
	FamilyDifferenceThenVerify: 6,
	FamilyRatioOfDifference:    8,
	FamilyReconciliationChain:  12,
}

// tpl is a terse InputSpec constructor.
func tpl(pairs map[string]any) InputSpec { return InputSpec{Template: pairs} }

// locateStep builds the REAL-mode LOCATE_REGION node for operand `role`.
// Arm C resolves the deterministic region locator; Arm B leaves
// infrastructure capabilities to the registry too.
func locateStep(role string) Step {
	return Step{
		LocalID: "locate_" + role, Capability: "LOCATE_REGION", Role: role,
		Input: tpl(map[string]any{
			"mode":      "REAL",
			"question":  "${param:question_" + role + "}",
			"store_dir": "${param:store_dir}",
			"limit":     5,
		}),
	}
}

// readStep builds the EXTRACT_NUMBER node for operand `role`. It consumes
// the whole located-region observation (page/bbox/page dimensions pass
// through intact) plus the per-operand page image path.
func readStep(role string) Step {
	return Step{
		LocalID: "read_" + role, Capability: "EXTRACT_NUMBER", Role: role,
		DependsOn: []string{"locate_" + role},
		Input: tpl(map[string]any{
			"opcode":          "EXTRACT_NUMBER",
			"page_image_path": "${param:page_image_" + role + "}",
			"region":          "${obs:locate_" + role + "}",
		}),
	}
}

// normReadStep normalizes a raw EXTRACT_NUMBER text output to a number.
func normReadStep(role string) Step {
	return Step{
		LocalID: "norm_" + role, Capability: "NORMALIZE", Role: role,
		DependsOn: []string{"read_" + role},
		Input: tpl(map[string]any{
			"raw":         "${obs:read_" + role + ":text}",
			"target_type": "number",
		}),
		PreferDeterministic: true,
	}
}

// normValueStep normalizes an upstream deterministic numeric field
// (an arithmetic result) back to a canonical number.
func normValueStep(localID, fromLocal, fromField string, deps ...string) Step {
	return Step{
		LocalID: localID, Capability: "NORMALIZE",
		DependsOn: append([]string{fromLocal}, deps...),
		Input: tpl(map[string]any{
			"raw":         "${obs:" + fromLocal + ":" + fromField + "}",
			"target_type": "number",
		}),
		PreferDeterministic: true,
	}
}

// arithStep builds an ARITHMETIC node (SUBTRACT | RATIO | PERCENT_DIFFERENCE)
// over two upstream :trimmed operands.
func arithStep(localID, op, aLocal, aField, bLocal, bField string) Step {
	deps := []string{}
	if aLocal != "" {
		deps = append(deps, aLocal)
	}
	if bLocal != "" && bLocal != aLocal {
		deps = append(deps, bLocal)
	}
	aRef := "${obs:" + aLocal + ":" + aField + "}"
	bRef := "${obs:" + bLocal + ":" + bField + "}"
	if bLocal == "" {
		bRef = bField // literal (a param placeholder or a constant string)
	}
	return Step{
		LocalID: localID, Capability: "ARITHMETIC", DependsOn: deps,
		Input:               tpl(map[string]any{"operation": op, "a": aRef, "b": bRef}),
		PreferDeterministic: true,
	}
}

// acquire returns the three-node LOCATE -> READ -> NORMALIZE chain for one
// operand role.
func acquire(role string) []Step {
	return []Step{locateStep(role), readStep(role), normReadStep(role)}
}

// --- Shape 1: READ_AND_CHECK (natural 2) ---------------------------------
//
// One operand: read it, normalize it, compare it against a constant
// threshold param. Terminal COMPARE_NUMBERS observation, no VERIFY.
func shapeReadAndCheck() TaskFamily {
	steps := acquire("A")
	steps = append(steps, Step{
		LocalID: "check", Capability: "COMPARE_NUMBERS", DependsOn: []string{"norm_A"},
		Input: tpl(map[string]any{
			"a": "${obs:norm_A:trimmed}", "b": "${param:threshold}",
		}),
		PreferDeterministic: true,
	})
	return TaskFamily{ID: FamilyReadAndCheck, Goal: "read one value and check it against a threshold", Steps: steps}
}

// --- Shape 2: COMPARE_TWO_VALUES (natural 4) ----------------------------
func shapeCompareTwoValues() TaskFamily {
	steps := append(acquire("A"), acquire("B")...)
	steps = append(steps, Step{
		LocalID: "cmp", Capability: "COMPARE_NUMBERS", DependsOn: []string{"norm_A", "norm_B"},
		Input: tpl(map[string]any{
			"a": "${obs:norm_A:trimmed}", "b": "${obs:norm_B:trimmed}",
		}),
		PreferDeterministic: true,
	})
	return TaskFamily{ID: FamilyCompareTwoValues, Goal: "read two values and compare them", Steps: steps}
}

// --- Shape 3: DIFFERENCE_THEN_VERIFY (natural 6, critical path 6) -------
//
// diff = A - B ; norm_diff = NORMALIZE(diff) ; VERIFY(norm_diff : number).
func shapeDifferenceThenVerify() TaskFamily {
	steps := append(acquire("A"), acquire("B")...)
	steps = append(steps,
		arithStep("diff", "SUBTRACT", "norm_A", "trimmed", "norm_B", "trimmed"),
		normValueStep("norm_diff", "diff", "result"),
		verifyNumberStep("verify", "norm_diff"),
	)
	return TaskFamily{ID: FamilyDifferenceThenVerify, Goal: "read two values, take their difference, verify it", Steps: steps}
}

// --- Shape 4: RATIO_OF_DIFFERENCE (natural 8, critical path 8) ----------
//
// Approved 3-operand form (plan stateful-seeking-puddle.md family table;
// original call-budget "Shape 4: 12 x 3"): ratio OF a difference.
//
//	diff       = A - B
//	norm_diff  = NORMALIZE(diff)
//	ratio      = RATIO(norm_diff, C)
//	norm_ratio = NORMALIZE(ratio)
//	verify     = VERIFY(norm_ratio : number)
//
// Critical path (8): locate_A read_A norm_A diff norm_diff ratio norm_ratio
// verify. norm_B feeds diff and norm_C feeds ratio, so all three
// observations are transitively necessary.
func shapeRatioOfDifference() TaskFamily {
	steps := acquire("A")
	steps = append(steps, acquire("B")...)
	steps = append(steps, acquire("C")...)
	steps = append(steps,
		arithStep("diff", "SUBTRACT", "norm_A", "trimmed", "norm_B", "trimmed"),
		normValueStep("norm_diff", "diff", "result"),
		arithStep("ratio", "RATIO", "norm_diff", "trimmed", "norm_C", "trimmed"),
		normValueStep("norm_ratio", "ratio", "result"),
		verifyNumberStep("verify", "norm_ratio"),
	)
	return TaskFamily{ID: FamilyRatioOfDifference, Goal: "read three values, ratio of a difference, verify it", Steps: steps}
}

// --- Shape 5: RECONCILIATION_CHAIN (natural 12, critical path 12) -------
//
// Dimensionally coherent reconciliation (protocol section 3, corrected):
//
//	sub_A            = A - a
//	sub_B            = B - b
//	disagreement_pct = PERCENT_DIFFERENCE(sub_A, sub_B)   (both are differences)
//	norm_pct         = NORMALIZE(disagreement_pct)
//	fraction         = RATIO(norm_pct, 100)              (percent -> fraction)
//	norm_fraction    = NORMALIZE(fraction)
//	tolerance_margin = fraction - tolerance              (fraction - fraction)
//	norm_margin      = NORMALIZE(tolerance_margin)
//	cmp_zero         = COMPARE_NUMBERS(norm_margin, 0)
//	verify           = VERIFY(norm_margin : number), sequenced after cmp_zero
//
// Critical path (12): locate_A read_A norm_A sub_A disagreement_pct norm_pct
// fraction norm_fraction tolerance_margin norm_margin cmp_zero verify.
// All four observations (norm_A, norm_a, norm_B, norm_b) are transitively
// necessary via sub_A / sub_B; cmp_zero is an explicit predecessor of
// verify. The frozen VERIFY contract is single-target, so verify promotes
// the numeric margin and the within-tolerance sign lives on the Blackboard
// as cmp_zero (protocol limitation, documented).
func shapeReconciliationChain() TaskFamily {
	steps := acquire("A")
	steps = append(steps, acquire("a")...)
	steps = append(steps, acquire("B")...)
	steps = append(steps, acquire("b")...)
	steps = append(steps,
		arithStep("sub_A", "SUBTRACT", "norm_A", "trimmed", "norm_a", "trimmed"),
		arithStep("sub_B", "SUBTRACT", "norm_B", "trimmed", "norm_b", "trimmed"),
		arithStep("disagreement_pct", "PERCENT_DIFFERENCE", "sub_A", "result", "sub_B", "result"),
		normValueStep("norm_pct", "disagreement_pct", "result"),
		arithStep("fraction", "RATIO", "norm_pct", "trimmed", "", "100"),
		normValueStep("norm_fraction", "fraction", "result"),
		arithStep("tolerance_margin", "SUBTRACT", "norm_fraction", "trimmed", "", "${param:tolerance}"),
		normValueStep("norm_margin", "tolerance_margin", "result"),
		Step{
			LocalID: "cmp_zero", Capability: "COMPARE_NUMBERS", DependsOn: []string{"norm_margin"},
			Input:               tpl(map[string]any{"a": "${obs:norm_margin:trimmed}", "b": "0"}),
			PreferDeterministic: true,
		},
		verifyNumberStep("verify", "norm_margin", "cmp_zero"),
	)
	return TaskFamily{ID: FamilyReconciliationChain, Goal: "reconcile two independent differences against a tolerance", Steps: steps}
}

// verifyNumberStep builds the terminal VERIFY node promoting `targetLocal`
// as a number. Extra deps are sequenced predecessors (their observations
// stay on the Blackboard for analysis but VERIFY reads only the target).
func verifyNumberStep(localID, targetLocal string, extraDeps ...string) Step {
	return Step{
		LocalID: localID, Capability: "VERIFY",
		DependsOn: append([]string{targetLocal}, extraDeps...),
		Input: tpl(map[string]any{
			"target_key":    "${node:" + targetLocal + "}",
			"fact_id":       "${node:" + targetLocal + "}",
			"expected_type": "number",
		}),
		PreferDeterministic: true,
	}
}

// Families returns the five frozen T1 task families in shape order.
func Families() []TaskFamily {
	return []TaskFamily{
		shapeReadAndCheck(),
		shapeCompareTwoValues(),
		shapeDifferenceThenVerify(),
		shapeRatioOfDifference(),
		shapeReconciliationChain(),
	}
}

// FamilyByID returns the frozen family with id, normalized.
func FamilyByID(id string) (TaskFamily, bool) {
	for _, family := range Families() {
		if family.ID == id {
			return family, true
		}
	}
	return TaskFamily{}, false
}

// HasVerify reports whether the family terminates on a Fact-promoting
// VERIFY node (Shapes 3-5) rather than an evaluable terminal observation
// (Shapes 1-2).
func (f TaskFamily) HasVerify() bool {
	for _, step := range f.Steps {
		if step.Capability == "VERIFY" {
			return true
		}
	}
	return false
}
