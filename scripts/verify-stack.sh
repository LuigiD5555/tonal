#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
chmod +x scripts/*.sh tests/*.sh

tests/test-manifest.sh
tests/test-skills.sh
tests/test-gatekeeper.sh
tests/test-fixed-carrier-r2.sh
scripts/fetch-components.sh
scripts/verify-components.sh
scripts/verify-fixed-carrier-r2.sh

python3 - "$ROOT/TONAL.json" "$ROOT/.work/components" <<'PY'
import json, pathlib, subprocess, sys
manifest = json.load(open(sys.argv[1]))
base = pathlib.Path(sys.argv[2]).resolve()
for name, component in manifest['components'].items():
    repo = (base / name).resolve()
    if not repo.is_dir() or base not in repo.parents:
        raise SystemExit(f'invalid component checkout for {name}: {repo}')
    checks = component.get('verification', [])
    if not checks:
        raise SystemExit(f'{name}: no verification declared in TONAL.json')
    print(f'VERIFY {name} [{component["kind"]}]')
    for index, check in enumerate(checks, 1):
        argv = check.get('argv')
        cwd_rel = check.get('cwd', '.')
        if not isinstance(argv, list) or not argv or not all(isinstance(x, str) and x for x in argv):
            raise SystemExit(f'{name}: verification #{index} has invalid argv')
        cwd = (repo / cwd_rel).resolve()
        if cwd != repo and repo not in cwd.parents:
            raise SystemExit(f'{name}: verification #{index} escapes component root')
        if not cwd.is_dir():
            raise SystemExit(f'{name}: verification cwd missing: {cwd}')
        print('  +', ' '.join(argv), f'(cwd={cwd_rel})')
        subprocess.run(argv, cwd=cwd, check=True)
print('DECLARED COMPONENT VERIFICATION: PASS')
PY

echo 'TONAL COMPOSITION VERIFY: PASS'
