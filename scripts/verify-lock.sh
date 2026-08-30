#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
python3 - "$ROOT" <<'PY'
import json, pathlib, re, sys
root = pathlib.Path(sys.argv[1])
version = (root / "VERSION").read_text().strip()
manifest = json.loads((root / "TONAL.json").read_text())
lock = json.loads((root / "tonal.lock").read_text())
assert manifest["schema"] == "tonal.composition.v2"
assert lock["schema"] == "tonal.lock.v2"
assert version == manifest["tonal"]["version"] == lock["tonal_version"]
assert set(manifest["components"]) == set(lock["components"])
allowed = set(manifest["composition_model"]["component_kinds"])
for name in sorted(manifest["components"]):
    m = manifest["components"][name]
    l = lock["components"][name]
    assert m["kind"] == l["kind"] in allowed, (name, m.get("kind"), l.get("kind"))
    assert m["version"] == l["version"]
    assert m["repository"] == l["repository"]
    assert m["commit"] == l["commit"]
    assert re.fullmatch(r"[0-9a-f]{40}", l["commit"]), (name, l["commit"])
    assert m["repository"].startswith("https://github.com/") and m["repository"].endswith(".git")
    assert isinstance(m.get("verification", []), list)
assert "origami" in manifest["composition_model"]["current_targets"]
assert manifest["components"]["origami"]["kind"] == "target"
assert manifest["components"]["tlaloc"]["kind"] == "development_tool"
assert "Blueprint Framework" in manifest["composition_model"]["unlocked_examples"]
assert "blueprint" not in {name.lower() for name in lock["components"]}
print(f"PASS Tonal {version}: generic manifest and lock are coherent")
PY
