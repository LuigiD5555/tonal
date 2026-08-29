#!/usr/bin/env bash
set -euo pipefail
python3 - <<'PY'
import json, subprocess, urllib.request
m=json.load(open('TONAL.json'))
for name, c in m['components'].items():
    locked=c['commit']; repo=c['repository']
    url=repo.replace('https://github.com/','https://api.github.com/repos/')+'/commits/main'
    with urllib.request.urlopen(url) as r: latest=json.load(r)['sha']
    status='CURRENT' if latest==locked else 'DRIFT'
    print(f'{name}: {status} locked={locked} main={latest}')
PY
