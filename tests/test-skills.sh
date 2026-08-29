#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

for name in repo-flow gatekeeper; do
  CANON="$ROOT/skills/$name"
  CLAUDE="$ROOT/.claude/skills/$name"
  AGENTS="$ROOT/.agents/skills/$name"
  for path in "$CANON" "$CLAUDE" "$AGENTS"; do
    [[ -f "$path/SKILL.md" ]] || { echo "missing $name skill: $path" >&2; exit 1; }
  done
  diff -qr "$CANON" "$CLAUDE"
  diff -qr "$CANON" "$AGENTS"
  grep -q "^name: $name$" "$CANON/SKILL.md"

done

grep -q '^version: 0\.2\.0$' "$ROOT/skills/repo-flow/SKILL.md"
grep -q '^version: 0\.1\.0$' "$ROOT/skills/gatekeeper/SKILL.md"

TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT
git init -q "$TMP/project"

for name in repo-flow gatekeeper; do
  CANON="$ROOT/skills/$name"
  "$ROOT/scripts/install-skill.sh" "$name" --project "$TMP/project"
  diff -qr "$CANON" "$TMP/project/.claude/skills/$name"
  diff -qr "$CANON" "$TMP/project/.agents/skills/$name"
done

printf '\nlocal-edit\n' >> "$TMP/project/.agents/skills/gatekeeper/SKILL.md"
if "$ROOT/scripts/install-skill.sh" gatekeeper --project "$TMP/project" >/dev/null 2>&1; then
  echo "installer overwrote or accepted differing local content without --force" >&2
  exit 1
fi

"$ROOT/scripts/install-skill.sh" gatekeeper --project "$TMP/project" --force
diff -qr "$ROOT/skills/gatekeeper" "$TMP/project/.claude/skills/gatekeeper"
diff -qr "$ROOT/skills/gatekeeper" "$TMP/project/.agents/skills/gatekeeper"

echo "PASS canonical/mirror/install coherence for repo-flow and gatekeeper"
