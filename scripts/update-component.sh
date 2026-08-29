#!/usr/bin/env bash
set -euo pipefail
component=${1:-}
revision=${2:-main}
case "$component" in tlaloc|origami) ;; *) echo "usage: $0 {tlaloc|origami} [revision]" >&2; exit 2;; esac
python3 - "$component" "$revision" <<'PY'
import json, subprocess, sys
name, rev=sys.argv[1:]
p='TONAL.json'; m=json.load(open(p)); target=m['components'][name]
repo=target['repository']; remote=repo+'.git'
sha=subprocess.check_output(['git','ls-remote',remote,rev],text=True).split()[0]
version=subprocess.check_output(['git','show',f'{sha}:VERSION'],text=True,stderr=subprocess.DEVNULL).strip() if False else None
# Resolve VERSION without mutating the working tree.
import urllib.request
raw=repo.replace('https://github.com/','https://raw.githubusercontent.com/')+f'/{sha}/VERSION'
with urllib.request.urlopen(raw) as r: version=r.read().decode().strip()
target['commit']=sha; target['version']=version
open(p,'w').write(json.dumps(m,indent=2)+'\n')
lock=json.load(open('tonal.lock')); lt=lock['components'][name]; lt['commit']=sha; lt['version']=version
open('tonal.lock','w').write(json.dumps(lock,indent=2)+'\n')
print(f'updated {name}: version={version} commit={sha}')
PY
scripts/verify-lock.sh
