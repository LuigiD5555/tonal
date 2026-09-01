#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

new_fixture() {
  local fixture
  fixture="$(mktemp -d)"
  mkdir -p "$fixture/scripts" "$fixture/tools" "$fixture/state" "$fixture/docs" "$fixture/runs"
  cp "$ROOT/VERSION" "$ROOT/README.md" "$ROOT/TONAL.json" "$ROOT/tonal.lock" "$fixture/"
  cp "$ROOT/scripts/verify-lock.sh" "$fixture/scripts/"
  cp "$ROOT/tools/claims.py" "$fixture/tools/"
  cp "$ROOT/state/CLAIMS.json" "$fixture/state/"
  cp "$ROOT/docs/CAPABILITY_STATUS.md" "$fixture/docs/"
  chmod +x "$fixture/scripts/verify-lock.sh" "$fixture/tools/claims.py"
  printf '%s\n' "$fixture"
}

expect_gate_failure() {
  local label="$1"
  local fixture="$2"
  local expected="$3"
  local output
  if output="$("$fixture/scripts/verify-lock.sh" 2>&1)"; then
    echo "FAIL $label: verify-lock unexpectedly passed" >&2
    rm -rf "$fixture"
    exit 1
  fi
  if [[ -n "$expected" ]] && ! grep -Fq "$expected" <<<"$output"; then
    echo "FAIL $label: rejection did not contain expected evidence" >&2
    echo "$output" >&2
    rm -rf "$fixture"
    exit 1
  fi
  echo "PASS negative gate: $label"
  rm -rf "$fixture"
}

fixture="$(new_fixture)"
printf '%s\n' '0.1.0-alpha.bad' > "$fixture/VERSION"
expect_gate_failure "VERSION drift" "$fixture" ""

fixture="$(new_fixture)"
python3 - "$fixture/README.md" <<'PY'
from pathlib import Path
import sys
path = Path(sys.argv[1])
lines = path.read_text().splitlines()
lines[0] = "# Tonal 0.1.0-alpha.bad"
path.write_text("\n".join(lines) + "\n")
PY
expect_gate_failure "README drift" "$fixture" ""

fixture="$(new_fixture)"
python3 - "$fixture/tonal.lock" <<'PY'
import json, pathlib, sys
path = pathlib.Path(sys.argv[1])
value = json.loads(path.read_text())
value["tonal_version"] = "0.1.0-alpha.bad"
path.write_text(json.dumps(value, indent=2) + "\n")
PY
expect_gate_failure "tonal.lock drift" "$fixture" ""

fixture="$(new_fixture)"
python3 - "$fixture/state/CLAIMS.json" <<'PY'
import json, pathlib, sys
path = pathlib.Path(sys.argv[1])
claims = json.loads(path.read_text())
claims[0]["statement"] += " [drift]"
path.write_text(json.dumps(claims, indent=2) + "\n")
PY
expect_gate_failure "ledger/generated-state drift" "$fixture" "generated claims table is stale"

fixture="$(new_fixture)"
python3 - "$fixture/state/CLAIMS.json" <<'PY'
import json, pathlib, sys
path = pathlib.Path(sys.argv[1])
claims = json.loads(path.read_text())
claims[0]["status"] = "evidenced"
claims[0]["evidence"] = ["run:missing-run-id"]
path.write_text(json.dumps(claims, indent=2) + "\n")
PY
expect_gate_failure "missing evidenced run_id" "$fixture" "run evidence does not exist in runs/: missing-run-id"

fixture="$(new_fixture)"
python3 - "$fixture/state/CLAIMS.json" "$fixture/runs" <<'PY'
import json, pathlib, sys
claims_path = pathlib.Path(sys.argv[1])
runs_root = pathlib.Path(sys.argv[2])
claims = json.loads(claims_path.read_text())
claims[0]["status"] = "evidenced"
claims[0]["evidence"] = ["run:fixture-run-id"]
claims_path.write_text(json.dumps(claims, indent=2) + "\n")
run_dir = runs_root / "2026-09"
run_dir.mkdir(parents=True, exist_ok=True)
(run_dir / "fixture.json").write_text(json.dumps({"run_id": "fixture-run-id"}, indent=2) + "\n")
PY
python3 "$fixture/tools/claims.py" generate >/dev/null
"$fixture/scripts/verify-lock.sh" >/dev/null
echo "PASS positive gate: evidenced run_id exists"
rm -rf "$fixture"

echo "EVIDENCE GATES C1/C2: PASS"
