#!/usr/bin/env bash
set -euo pipefail
root=$(cd "$(dirname "$0")/.." && pwd)
cd "$root"
rm -rf .work/snapshot dist
mkdir -p .work/snapshot/components dist
scripts/verify-lock.sh
scripts/fetch-components.sh
scripts/verify-components.sh
python3 - "$root/tonal.lock" "$root/.work/components" "$root/.work/snapshot/components" <<'PY'
import json, pathlib, shutil, sys
lock = json.load(open(sys.argv[1]))
src = pathlib.Path(sys.argv[2])
dst = pathlib.Path(sys.argv[3])
for name in sorted(lock['components']):
    source = src / name
    target = dst / name
    if not source.is_dir():
        raise SystemExit(f'missing fetched component {source}')
    shutil.copytree(source, target, symlinks=True)
    gitdir = target / '.git'
    if gitdir.exists():
        shutil.rmtree(gitdir)
PY

# Architecture R2 keeps the historical composition artifacts, but the source
# snapshot must now include the Tonal runtime and the current authority docs.
cp VERSION TONAL.json tonal.lock compatibility.json CLAUDE.md README.md .work/snapshot/
mkdir -p .work/snapshot/docs
cp \
  docs/CURRENT_STATE.md \
  docs/ARCHITECTURE.md \
  docs/BOUNDARIES.md \
  docs/RESEARCH_PROGRAM.md \
  .work/snapshot/docs/
cp -a runtime .work/snapshot/runtime
mkdir -p .work/snapshot/skills
cp -a skills/. .work/snapshot/skills/

version=$(cat VERSION)
archive="dist/tonal-${version}-source.tar.gz"
TZ=UTC tar --sort=name --mtime='UTC 1970-01-01' --owner=0 --group=0 --numeric-owner -C .work/snapshot -cf - . | gzip -n > "$archive"
sha256sum "$archive" > "$archive.sha256"
echo "$archive"
cat "$archive.sha256"
