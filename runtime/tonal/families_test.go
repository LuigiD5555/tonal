package tonal

import "testing"

// The frozen T1 families must normalize cleanly and their mechanical
// critical-path depth is fixed at freeze time.
func TestFrozenFamilies_NormalizeAndDepth(t *testing.T) {
	wantCriticalPath := map[string]int{
		FamilyReadAndCheck:         4,
		FamilyCompareTwoValues:     4,
		FamilyDifferenceThenVerify: 6,
		FamilyRatioOfDifference:    8,
		FamilyReconciliationChain:  12,
	}
	wantNodeCount := map[string]int{
		FamilyReadAndCheck:         4,
		FamilyCompareTwoValues:     7,
		FamilyDifferenceThenVerify: 9,
		FamilyRatioOfDifference:    14,
		FamilyReconciliationChain:  22,
	}
	for _, family := range Families() {
		normalized, err := family.Normalize()
		if err != nil {
			t.Fatalf("family %s: normalize: %v", family.ID, err)
		}
		if got := normalized.CriticalPathDepth(); got != wantCriticalPath[family.ID] {
			t.Errorf("family %s: critical path depth = %d, want %d", family.ID, got, wantCriticalPath[family.ID])
		}
		if got := len(normalized.Steps); got != wantNodeCount[family.ID] {
			t.Errorf("family %s: node count = %d, want %d", family.ID, got, wantNodeCount[family.ID])
		}
		if _, ok := NaturalDepth[family.ID]; !ok {
			t.Errorf("family %s: no NaturalDepth entry", family.ID)
		}
	}
}

// Only Shapes 3-5 carry a Fact-promoting VERIFY (protocol section 4).
func TestFrozenFamilies_VerifyOnlyInShapes345(t *testing.T) {
	wantVerify := map[string]bool{
		FamilyReadAndCheck:         false,
		FamilyCompareTwoValues:     false,
		FamilyDifferenceThenVerify: true,
		FamilyRatioOfDifference:    true,
		FamilyReconciliationChain:  true,
	}
	for _, family := range Families() {
		if got := family.HasVerify(); got != wantVerify[family.ID] {
			t.Errorf("family %s: HasVerify = %v, want %v", family.ID, got, wantVerify[family.ID])
		}
		verifyCount := 0
		for _, step := range family.Steps {
			if step.Capability == "VERIFY" {
				verifyCount++
			}
		}
		if verifyCount > 1 {
			t.Errorf("family %s: %d VERIFY nodes, want at most 1", family.ID, verifyCount)
		}
	}
}

// The corrected Shape 5 must not reference the withdrawn incoherent spine.
func TestShape5_CorrectedSpine(t *testing.T) {
	family, ok := FamilyByID(FamilyReconciliationChain)
	if !ok {
		t.Fatal("RECONCILIATION_CHAIN family missing")
	}
	byID := map[string]Step{}
	for _, step := range family.Steps {
		byID[step.LocalID] = step
	}
	for _, banned := range []string{"norm_ratio", "norm_resid"} {
		if _, present := byID[banned]; present {
			t.Errorf("Shape 5 still references withdrawn node %q", banned)
		}
	}
	for _, want := range []string{"sub_A", "sub_B", "disagreement_pct", "fraction", "tolerance_margin", "cmp_zero", "verify"} {
		if _, present := byID[want]; !present {
			t.Errorf("Shape 5 missing corrected node %q", want)
		}
	}
	if op, _ := byID["disagreement_pct"].Input.Template["operation"]; op != "PERCENT_DIFFERENCE" {
		t.Errorf("disagreement_pct operation = %v, want PERCENT_DIFFERENCE", op)
	}
	// All four raw observations feed the terminal verify transitively.
	if d := family.CriticalPathDepth(); d != 12 {
		t.Fatalf("Shape 5 critical path = %d, want 12", d)
	}
}
