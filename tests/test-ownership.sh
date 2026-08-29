#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
COMPONENTS="${1:-$ROOT/.work/components}"

CANON="$ROOT/skills/repo-flow/SKILL.md"
CLAUDE="$ROOT/.claude/skills/repo-flow/SKILL.md"
AGENTS="$ROOT/.agents/skills/repo-flow/SKILL.md"

cmp "$CANON" "$CLAUDE"
cmp "$CANON" "$AGENTS"

TLALOC="$COMPONENTS/tlaloc"
[[ -d "$TLALOC" ]] || { echo "missing fetched Tlaloc checkout: $TLALOC" >&2; exit 1; }
[[ ! -e "$TLALOC/.claude/skills/repo-flow" ]] || { echo "duplicate repo-flow authority found in locked Tlaloc" >&2; exit 1; }

expected=(tlaloc-project tlaloc-behavior tlaloc-tlaloque origami-semantics tlaloc-release)
for name in "${expected[@]}"; do
  [[ -f "$TLALOC/.claude/skills/$name/SKILL.md" ]] || { echo "missing Tlaloc-owned skill: $name" >&2; exit 1; }
done
count="$(find "$TLALOC/.claude/skills" -mindepth 1 -maxdepth 1 -type d | wc -l | tr -d ' ')"
[[ "$count" == "5" ]] || { echo "expected exactly 5 Tlaloc-owned skill directories, got $count" >&2; exit 1; }

grep -q 'REPO_FLOW_SINGLE_CANONICAL_AUTHORITY' "$ROOT/TONAL.json"

echo "PASS repo-flow single authority: Tonal canonical, Tlaloc component-specific only"
