#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CANON_ROOT="$ROOT/skills"

[[ -d "$CANON_ROOT" ]] || { echo "missing canonical skills directory: $CANON_ROOT" >&2; exit 1; }

for skill_dir in "$CANON_ROOT"/*; do
  [[ -d "$skill_dir" && -f "$skill_dir/SKILL.md" ]] || continue
  name="$(basename "$skill_dir")"
  for family in .claude .agents; do
    dst="$ROOT/$family/skills/$name"
    rm -rf -- "$dst"
    mkdir -p "$(dirname "$dst")"
    cp -a "$skill_dir" "$dst"
  done
  echo "synced $name"
done
