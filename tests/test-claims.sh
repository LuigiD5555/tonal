#!/usr/bin/env bash
set -euo pipefail

REPOSITORY_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CLAIMS_TOOL="$REPOSITORY_ROOT/tools/claims.py"
CAPABILITY_DOCUMENT="$REPOSITORY_ROOT/docs/CAPABILITY_STATUS.md"

python3 "$CLAIMS_TOOL" validate
python3 "$CLAIMS_TOOL" check

python3 - "$REPOSITORY_ROOT" "$CLAIMS_TOOL" "$CAPABILITY_DOCUMENT" <<'PY'
import copy
import subprocess
import sys
import tempfile
from pathlib import Path

repository_root = Path(sys.argv[1])
claims_tool = Path(sys.argv[2])
capability_document = Path(sys.argv[3])
sys.path.insert(0, str(repository_root))
from tools import claims as claims_module

baseline_claims = claims_module.load_claims(repository_root / "state" / "CLAIMS.json")

duplicate_claims = copy.deepcopy(baseline_claims)
duplicate_claims.append(copy.deepcopy(duplicate_claims[0]))
assert any("duplicate claim id" in error for error in claims_module.validate_claims(duplicate_claims))

invalid_status_claims = copy.deepcopy(baseline_claims)
invalid_status_claims[0]["status"] = "claimed"
assert any("status is not allowed" in error for error in claims_module.validate_claims(invalid_status_claims))

missing_evidence_claims = copy.deepcopy(baseline_claims)
missing_evidence_claims[0]["status"] = "implemented"
missing_evidence_claims[0]["evidence"] = []
assert any("requires evidence" in error for error in claims_module.validate_claims(missing_evidence_claims))

missing_test_claims = copy.deepcopy(baseline_claims)
missing_test_claims[0]["status"] = "implemented"
missing_test_claims[0]["evidence"] = ["test:missing-package:TestThatDoesNotExist"]
assert any("package does not exist" in error for error in claims_module.validate_claims(missing_test_claims))

with tempfile.TemporaryDirectory() as test_directory:
    drifted_document = Path(test_directory) / capability_document.name
    document_text = capability_document.read_text(encoding="utf-8")
    drifted_document.write_text(
        document_text.replace(baseline_claims[0]["id"], "DRIFTED.CLAIM", 1),
        encoding="utf-8",
    )
    check_command = [sys.executable, str(claims_tool), "check", "--document", str(drifted_document)]
    assert subprocess.run(check_command, capture_output=True, check=False).returncode != 0
    subprocess.run(
        [sys.executable, str(claims_tool), "generate", "--document", str(drifted_document)],
        check=True,
        capture_output=True,
    )
    subprocess.run(check_command, check=True, capture_output=True)
PY

echo PASS
