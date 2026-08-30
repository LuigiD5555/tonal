#!/usr/bin/env bash
set -euo pipefail
component=${1:-}
revision=${2:-main}
if [[ -z "$component" ]]; then
  echo "usage: $0 <declared-component> [revision]" >&2
  exit 2
fi
python3 - "$component" "$revision" <<'PY'
import json, subprocess, sys, urllib.request
name, rev = sys.argv[1:]
p = 'TONAL.json'
m = json.load(open(p))
if name not in m['components']:
    raise SystemExit(f"unknown component {name!r}; declare it in TONAL.json before pinning")
target = m['components'][name]
repo = target['repository']
rows = subprocess.check_output(['git','ls-remote',repo,rev], text=True).splitlines()
if not rows:
    raise SystemExit(f"cannot resolve {rev!r} in {repo}")
sha = rows[0].split()[0]
base = repo.removesuffix('.git').replace('https://github.com/','https://raw.githubusercontent.com/')
raw = f'{base}/{sha}/VERSION'
with urllib.request.urlopen(raw) as r:
    version = r.read().decode().strip()
target['commit'] = sha
target['version'] = version
open(p,'w').write(json.dumps(m,indent=2)+'\n')
lock = json.load(open('tonal.lock'))
lt = lock['components'][name]
lt['commit'] = sha
lt['version'] = version
lt['repository'] = repo
lt['kind'] = target['kind']
open('tonal.lock','w').write(json.dumps(lock,indent=2)+'\n')
print(f'updated {name}: kind={target["kind"]} version={version} commit={sha}')
PY
scripts/verify-lock.sh
