#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
python3 - "$ROOT" <<'PY'
import json, pathlib, re, sys
root = pathlib.Path(sys.argv[1])
version = (root / "VERSION").read_text().strip()
manifest = json.loads((root / "TONAL.json").read_text())
lock = json.loads((root / "tonal.lock").read_text())
assert manifest["schema"] == "tonal.composition.v1"
assert lock["schema"] == "tonal.lock.v1"
assert version == manifest["tonal"]["version"] == lock["tonal_version"]
expected_repos = {
    "tlaloc": "https://github.com/LuigiD5555/tlaloc.git",
    "origami": "https://github.com/LuigiD5555/origami.git",
}
for name in ("tlaloc", "origami"):
    m = manifest["components"][name]
    l = lock["components"][name]
    assert m["version"] == l["version"]
    assert m["repository"] == l["repository"] == expected_repos[name]
    assert m["commit"] == l["commit"]
    assert re.fullmatch(r"[0-9a-f]{40}", l["commit"]), (name, l["commit"])
print(f"PASS Tonal {version}: manifest and lock are coherent")
PY
