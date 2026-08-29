#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DEST="${1:-$ROOT/.work/components}"
python3 - "$ROOT/tonal.lock" "$DEST" <<'PY'
import json, pathlib, subprocess, sys
lock = json.load(open(sys.argv[1]))
base = pathlib.Path(sys.argv[2])
for name, c in lock["components"].items():
    repo = base / name
    if not repo.is_dir():
        raise SystemExit(f"missing component checkout: {repo}")
    head = subprocess.check_output(["git", "-C", str(repo), "rev-parse", "HEAD"], text=True).strip()
    version = (repo / "VERSION").read_text().strip()
    if head != c["commit"]:
        raise SystemExit(f"{name}: commit mismatch {head} != {c['commit']}")
    if version != c["version"]:
        raise SystemExit(f"{name}: version mismatch {version} != {c['version']}")
    print(f"PASS {name}: {version} @ {head}")
PY
