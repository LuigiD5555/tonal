#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DEST="${1:-$ROOT/.work/components}"
mkdir -p "$DEST"
python3 - "$ROOT/tonal.lock" <<'PY' | while IFS=$'\t' read -r name repo commit; do
import json, sys
lock = json.load(open(sys.argv[1]))
for name, c in lock["components"].items():
    print(name, c["repository"], c["commit"], sep="\t")
PY
  target="$DEST/$name"
  rm -rf "$target"
  git init -q "$target"
  git -C "$target" remote add origin "$repo"
  git -C "$target" fetch -q --depth 1 origin "$commit"
  git -C "$target" checkout -q --detach FETCH_HEAD
  actual="$(git -C "$target" rev-parse HEAD)"
  [[ "$actual" == "$commit" ]] || { echo "commit mismatch for $name" >&2; exit 1; }
  echo "fetched $name@$actual"
done
