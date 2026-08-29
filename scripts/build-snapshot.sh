#!/usr/bin/env bash
set -euo pipefail
root=$(cd "$(dirname "$0")/.." && pwd); cd "$root"
rm -rf .work/snapshot dist; mkdir -p .work/snapshot/components dist
scripts/verify-lock.sh
scripts/fetch-components.sh
scripts/verify-components.sh
cp -a .work/components/tlaloc .work/snapshot/components/tlaloc
cp -a .work/components/origami .work/snapshot/components/origami
rm -rf .work/snapshot/components/*/.git
cp VERSION TONAL.json tonal.lock PROJECT_BOUNDARY.md compatibility.json .work/snapshot/
mkdir -p .work/snapshot/skills; cp -a skills/. .work/snapshot/skills/
version=$(cat VERSION)
archive="dist/tonal-${version}-source.tar.gz"
TZ=UTC tar --sort=name --mtime='UTC 1970-01-01' --owner=0 --group=0 --numeric-owner -C .work/snapshot -cf - . | gzip -n > "$archive"
sha256sum "$archive" > "$archive.sha256"
echo "$archive"
cat "$archive.sha256"
