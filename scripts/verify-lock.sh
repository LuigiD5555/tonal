#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
COMPONENTS_ROOT="${1:-}"
python3 - "$ROOT" "$COMPONENTS_ROOT" <<'PY'
import json, pathlib, re, sys
root = pathlib.Path(sys.argv[1])
components_root = pathlib.Path(sys.argv[2]) if sys.argv[2] else None
version = (root / "VERSION").read_text().strip()
readme = (root / "README.md").read_text()
manifest = json.loads((root / "TONAL.json").read_text())
lock = json.loads((root / "tonal.lock").read_text())
assert manifest["schema"] == "tonal.composition.v2"
assert lock["schema"] == "tonal.lock.v2"
assert version == manifest["tonal"]["version"] == lock["tonal_version"]
assert readme.startswith(f"# Tonal {version}\n")
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

    if components_root is not None:
        component_root = components_root / name
        assert component_root.is_dir(), f"missing component checkout: {component_root}"
        component_version = (component_root / "VERSION").read_text().strip()
        assert component_version == l["version"], (
            f"{name}: VERSION {component_version} != tonal.lock {l['version']}"
        )
        component_readme = (component_root / "README.md").read_text()
        first_line = component_readme.splitlines()[0] if component_readme else ""
        header_match = re.fullmatch(r"#\s+\S+\s+(\S+)", first_line)
        assert header_match is not None, f"{name}: README version header is missing: {first_line!r}"
        assert header_match.group(1) == l["version"], (
            f"{name}: README version {header_match.group(1)} != tonal.lock {l['version']}"
        )
assert "origami" in manifest["composition_model"]["current_targets"]
assert manifest["components"]["origami"]["kind"] == "target"
assert manifest["components"]["tlaloc"]["kind"] == "development_tool"
assert "Blueprint Framework" in manifest["composition_model"]["unlocked_examples"]
assert "blueprint" not in {name.lower() for name in lock["components"]}
print(f"PASS Tonal {version}: generic manifest and lock are coherent")
PY
python3 "$ROOT/tools/claims.py" validate
python3 "$ROOT/tools/claims.py" check

if [[ -n "$COMPONENTS_ROOT" ]]; then
python3 - "$ROOT" "$COMPONENTS_ROOT" <<'PY'
import json, pathlib, subprocess, sys
root = pathlib.Path(sys.argv[1])
components_root = pathlib.Path(sys.argv[2])
lock = json.loads((root / "tonal.lock").read_text())
for name in sorted(lock["components"]):
    component_root = components_root / name
    claims_path = component_root / "state" / "CLAIMS.json"
    if not claims_path.is_file():
        continue
    command = [
        sys.executable,
        str(root / "tools" / "claims.py"),
        "validate",
        "--claims", str(claims_path),
        "--repository-root", str(component_root),
    ]
    subprocess.run(command, check=True)
    document_path = component_root / "docs" / "CAPABILITY_STATUS.md"
    if document_path.is_file():
        subprocess.run(
            command[:2]
            + [
                "check",
                "--claims", str(claims_path),
                "--document", str(document_path),
                "--repository-root", str(component_root),
            ],
            check=True,
        )
PY
fi
