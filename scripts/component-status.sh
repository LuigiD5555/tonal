#!/usr/bin/env bash
set -euo pipefail
python3 - <<'PY'
import json, urllib.request
m=json.load(open('TONAL.json'))
for name, c in sorted(m['components'].items()):
    locked=c['commit']; repo=c['repository'].removesuffix('.git')
    if not repo.startswith('https://github.com/'):
        print(f'{name}: STATUS_UNSUPPORTED repo={c["repository"]}')
        continue
    api=repo.replace('https://github.com/','https://api.github.com/repos/')+'/commits/main'
    try:
        with urllib.request.urlopen(api) as r: latest=json.load(r)['sha']
    except Exception as exc:
        print(f'{name}: STATUS_ERROR locked={locked} error={exc}')
        continue
    status='CURRENT' if latest==locked else 'DRIFT'
    print(f'{name} [{c["kind"]}]: {status} locked={locked} main={latest}')
PY
