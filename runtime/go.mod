module tonal.local/runtime

go 1.22

// Tonal's runtime consumes the exact Tlaloc commit pinned by tonal.lock.
// scripts/fetch-components.sh materialises it at
// .work/components/tlaloc/behavior-lab (git-ignored, rebuilt from the pin).
// The path is repo-relative, never machine-specific.
require tlaloc.local/behaviorlab v0.0.0

replace tlaloc.local/behaviorlab => ../.work/components/tlaloc/behavior-lab
