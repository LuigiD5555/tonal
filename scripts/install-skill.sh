#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

usage() {
  cat <<'USAGE'
Usage:
  ./scripts/install-skill.sh <name> [--project PATH] [--force]

Installs one canonical Tonal skill into both:
  .claude/skills/<name>/
  .agents/skills/<name>/

Existing differing content is not overwritten unless --force is explicit.
USAGE
}

name="${1:-}"
[[ -n "$name" ]] || { usage >&2; exit 2; }
shift
[[ "$name" =~ ^[A-Za-z0-9._-]+$ ]] || { echo "invalid skill name: $name" >&2; exit 2; }

project="$PWD"
force=0
while (($#)); do
  case "$1" in
    --project)
      shift
      (($#)) || { echo "--project requires a path" >&2; exit 2; }
      project="$1"
      ;;
    --force) force=1 ;;
    -h|--help) usage; exit 0 ;;
    *) echo "unknown option: $1" >&2; usage >&2; exit 2 ;;
  esac
  shift
done

src="$ROOT/skills/$name"
[[ -f "$src/SKILL.md" ]] || { echo "unknown Tonal skill: $name" >&2; exit 1; }
command -v git >/dev/null 2>&1 || { echo "git is required" >&2; exit 1; }
repo_root="$(git -C "$project" rev-parse --show-toplevel 2>/dev/null || true)"
[[ -n "$repo_root" ]] || { echo "not inside a Git repository: $project" >&2; exit 1; }

targets=(
  "$repo_root/.claude/skills/$name"
  "$repo_root/.agents/skills/$name"
)

for dst in "${targets[@]}"; do
  if [[ -e "$dst" ]] && ! diff -qr "$src" "$dst" >/dev/null 2>&1 && [[ "$force" -ne 1 ]]; then
    echo "refusing to overwrite existing differing skill: $dst" >&2
    echo "re-run with --force only after reviewing the local copy" >&2
    exit 1
  fi
done

for dst in "${targets[@]}"; do
  if [[ -e "$dst" ]] && diff -qr "$src" "$dst" >/dev/null 2>&1; then
    echo "already current: $dst"
    continue
  fi
  rm -rf -- "$dst"
  mkdir -p "$(dirname "$dst")"
  cp -a "$src" "$dst"
  echo "installed $name -> $dst"
done
