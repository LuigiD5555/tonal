#!/usr/bin/env bash
set -euo pipefail
chmod +x scripts/*.sh tests/*.sh
tests/test-manifest.sh
tests/test-skills.sh
tests/test-gatekeeper.sh
tests/test-fixed-carrier-r2.sh
scripts/fetch-components.sh
scripts/verify-components.sh
scripts/verify-fixed-carrier-r2.sh
(
 cd .work/components/origami
 go test ./...
 go vet ./...
 mkdir -p evidence
 go run ./cmd/origami-reference --experiment experiments/EXP-001-relational-state/experiment.json --out evidence/EXP-001-reference.json
 go run ./cmd/origami-evidence-gate --experiment experiments/EXP-001-relational-state/experiment.json --evidence evidence/EXP-001-reference.json
)
(
 cd .work/components/tlaloc/behavior-lab
 go test ./...
 go vet ./...
)
echo 'TONAL STACK VERIFY: PASS'
