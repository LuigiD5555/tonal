// Package arch holds architecture-boundary tests for the TONAL runtime.
package arch

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// runtimeRoot is the runtime module directory (parent of this package).
func runtimeRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Dir(wd)
}

// walkGo yields every non-test .go file under dir.
func walkGo(t *testing.T, dir string, fn func(path string, src []byte)) {
	t.Helper()
	err := filepath.WalkDir(dir, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		fn(path, src)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

// TONAL's engine and CLI must never DIRECTLY import a Tlaloc internal
// package. They see Tlaloc only through the public tlaloquekit contract.
// (tlalocbridge is the one wiring seam and is allowed to import the
// tlaloquekit builder, which is still public.)
func TestEngineDoesNotImportTlalocInternals(t *testing.T) {
	root := runtimeRoot(t)
	for _, sub := range []string{"tonal", "cmd", "arch"} {
		walkGo(t, filepath.Join(root, sub), func(path string, src []byte) {
			fileSet := token.NewFileSet()
			file, err := parser.ParseFile(fileSet, path, src, parser.ImportsOnly)
			if err != nil {
				t.Fatalf("parse %s: %v", path, err)
			}
			for _, imp := range file.Imports {
				value := strings.Trim(imp.Path.Value, `"`)
				if strings.Contains(value, "/internal/") {
					t.Errorf("%s directly imports internal package %q", path, value)
				}
			}
		})
	}
}

// The public contract package tlaloquekit (the only tlaloc package the
// engine imports) must itself be free of internal/* imports, so importing
// it drags in no engine internals.
func TestPublicContractPackageIsInternalFree(t *testing.T) {
	root := filepath.Join(runtimeRoot(t), "..", ".work", "components", "tlaloc", "behavior-lab", "tlaloquekit")
	if _, err := os.Stat(root); err != nil {
		t.Skipf("pinned Tlaloc checkout not materialised: %v", err)
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		src, err := os.ReadFile(filepath.Join(root, entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		fileSet := token.NewFileSet()
		file, err := parser.ParseFile(fileSet, entry.Name(), src, parser.ImportsOnly)
		if err != nil {
			t.Fatalf("parse %s: %v", entry.Name(), err)
		}
		for _, imp := range file.Imports {
			value := strings.Trim(imp.Path.Value, `"`)
			if strings.Contains(value, "tlaloc.local/behaviorlab/internal/") {
				t.Errorf("tlaloquekit/%s imports internal %q", entry.Name(), value)
			}
		}
	}
}

// TONAL must contain no executor-specific competence constants. A Tlaloque
// enforces its own envelope; the runtime stays model-agnostic.
func TestEngineHasNoParrotSpecificConstants(t *testing.T) {
	forbidden := []string{
		"UPSCALE", "LOW_SCALE", "crop_to_line", "CROP_TO_OPERAND_LINE",
		"line_height_px", "32px", "target_line_height", "padded_value_cue",
		"AdapterR1", "CapabilityProfileR1",
	}
	walkGo(t, filepath.Join(runtimeRoot(t), "tonal"), func(path string, src []byte) {
		text := string(src)
		for _, token := range forbidden {
			if strings.Contains(text, token) {
				t.Errorf("%s contains executor-specific token %q — competence knowledge leaked into the runtime", path, token)
			}
		}
	})
}
