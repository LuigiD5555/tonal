// Command callbudget independently derives the T1 model-call budget per arm
// directly from the frozen TaskFamily DAGs (runtime/tonal/families.go),
// without trusting any previously-reported number.
package main

import (
	"fmt"

	tonal "tonal.local/runtime/tonal"
)

// Capabilities that count as a Parrot model call under each arm's frozen
// routing policy (T1_D5_ARM_B_POLICY.json / T1_D5_ARM_C_POLICY.json).
var armBParrotCapabilities = map[string]bool{
	"EXTRACT_NUMBER":  true,
	"NORMALIZE":       true,
	"COMPARE_NUMBERS": true,
	"ARITHMETIC":      true,
}

var armCParrotCapabilities = map[string]bool{
	"EXTRACT_NUMBER": true, // only this one under Arm C's heterogeneous policy
}

func main() {
	families := tonal.Families()

	totalA := 0
	totalB := 0
	totalC := 0

	fmt.Println("shape,workflows,arm_a_calls_per_wf,arm_a_total,arm_b_calls_per_wf,arm_b_total,arm_c_calls_per_wf,arm_c_total")

	const workflowsPerShape = 12

	for _, family := range families {
		normalized, err := family.Normalize()
		if err != nil {
			fmt.Printf("ERROR normalizing %s: %v\n", family.ID, err)
			continue
		}

		armBPerWorkflow := 0
		armCPerWorkflow := 0
		for _, step := range normalized.Steps {
			if armBParrotCapabilities[step.Capability] {
				armBPerWorkflow++
			}
			if armCParrotCapabilities[step.Capability] {
				armCPerWorkflow++
			}
		}

		armAPerWorkflow := 1 // frozen policy: exactly one composite call per workflow

		armATotal := armAPerWorkflow * workflowsPerShape
		armBTotal := armBPerWorkflow * workflowsPerShape
		armCTotal := armCPerWorkflow * workflowsPerShape

		totalA += armATotal
		totalB += armBTotal
		totalC += armCTotal

		fmt.Printf("%s,%d,%d,%d,%d,%d,%d,%d\n",
			family.ID, workflowsPerShape,
			armAPerWorkflow, armATotal,
			armBPerWorkflow, armBTotal,
			armCPerWorkflow, armCTotal)
	}

	fmt.Println()
	fmt.Printf("EXPECTED_MODEL_CALLS_ARM_A = %d\n", totalA)
	fmt.Printf("EXPECTED_MODEL_CALLS_ARM_B = %d\n", totalB)
	fmt.Printf("EXPECTED_MODEL_CALLS_ARM_C = %d\n", totalC)
	fmt.Printf("EXPECTED_PRIMARY_TOTAL = %d\n", totalA+totalB+totalC)
}
