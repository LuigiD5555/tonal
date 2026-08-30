#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DEST="${1:-$ROOT/.work/candidates/fixed-carrier-r2}"
rm -rf "$DEST"
mkdir -p "$DEST"
python3 - "$ROOT/proposals/FIXED_CARRIER_R2.json" <<'PY2' | while IFS=$'\t' read -r name repo commit version; do
import json,sys
p=json.load(open(sys.argv[1]))
rows=[
 ('origami','https://github.com/LuigiD5555/origami.git',p['origami']['merged_commit'],p['origami']['required_version']),
 ('tlaloc','https://github.com/LuigiD5555/tlaloc.git',p['tlaloc']['merged_commit'],p['tlaloc']['required_version']),
]
for row in rows: print(*row,sep='\t')
PY2
  target="$DEST/$name"
  git init -q "$target"
  git -C "$target" remote add origin "$repo"
  git -C "$target" fetch -q --depth 1 origin "$commit"
  git -C "$target" checkout -q --detach FETCH_HEAD
  actual="$(git -C "$target" rev-parse HEAD)"
  [[ "$actual" == "$commit" ]] || { echo "$name candidate commit mismatch" >&2; exit 1; }
  actual_version="$(cat "$target/VERSION")"
  [[ "$actual_version" == "$version" ]] || { echo "$name candidate version mismatch: $actual_version != $version" >&2; exit 1; }
  case "$name" in
    origami)
      [[ -f "$target/spec/FIXED_CARRIER_R2.json" ]]
      grep -q 'origami.fixed-carrier.r2' "$target/spec/FIXED_CARRIER_R2.json"
      (cd "$target" && go test ./... && go vet ./...)
      ;;
    tlaloc)
      [[ -f "$target/docs/FIXED_CARRIER_MEMORY_PLANE_R2.md" ]]
      grep -q 'tlaloc.origami-tools.r2' "$target/behavior-lab/profiles/origami/fixed-carrier-r2.json"
      (cd "$target/behavior-lab" && go test ./... && go vet ./...)
      ;;
  esac
  echo "PASS candidate $name: $version @ $commit"
done
