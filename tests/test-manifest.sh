#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
"$ROOT/scripts/verify-lock.sh"
python3 - "$ROOT" <<'PY'
import json, pathlib, sys
root = pathlib.Path(sys.argv[1])
m = json.loads((root / "TONAL.json").read_text())
assert "COMPONENT_REPOSITORIES_REMAIN_AUTHORITATIVE" in m["invariants"]
assert "SNAPSHOT_REFERENCES_EXACT_COMMITS" in m["invariants"]
assert m["snapshot"]["physical_artifact"] == "NOT_YET_RELEASED"
print("PASS composition invariants")
PY
