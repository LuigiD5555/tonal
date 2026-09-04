// Command tonal-t1 drives the TONAL T1 heterogeneous-composition
// experiment: freeze the protocol/registry/dataset, run the readiness
// doctor, execute arms A/B/C over the frozen workflow set, and aggregate.
//
// This is the experiment harness skeleton. Dataset construction, the
// pre-inference freeze, the integrity doctor and the A/B/C run loop are
// filled in the later T1 steps; no model inference happens yet.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"tonal.local/runtime/tlalocbridge"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "capabilities":
		capabilities()
	case "doctor", "freeze", "run", "aggregate":
		fmt.Fprintf(os.Stderr, "tonal-t1 %s: not implemented yet (T1 steps 8-11)\n", os.Args[1])
		os.Exit(3)
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `tonal-t1 <capabilities|doctor|freeze|run|aggregate>

  capabilities  print the qualified Tlaloque set published by the pinned Tlaloc
  doctor        pre-inference integrity checks            (not implemented yet)
  freeze        write TONAL_T1_*.json + READY_TONAL_T1    (not implemented yet)
  run           execute arms A/B/C over the frozen set    (not implemented yet)
  aggregate     depth curve, paired tables, report        (not implemented yet)`)
}

func capabilities() {
	registry, err := tlalocbridge.Build(tlalocbridge.Config{})
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	body, _ := json.MarshalIndent(map[string]any{
		"parrot_profile_id":   registry.ParrotProfileID(),
		"parrot_profile_hash": registry.ParrotProfileHash(),
		"capabilities":        registry.Capabilities(),
	}, "", "  ")
	fmt.Println(string(body))
}
