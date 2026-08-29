#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
python3 - <<'PY' "$ROOT/gatekeeper.json"
import json,sys
p=json.load(open(sys.argv[1]))
assert p['schema']=='tonal.gatekeeper.r0'
assert p['owner']=='LuigiD5555'
assert set(p['applies_to'])=={'LuigiD5555/tonal','LuigiD5555/tlaloc','LuigiD5555/origami'}
assert p['policy']['OWNER']['may_explicitly_override_promotion_gate'] is True
assert p['policy']['EXTERNAL']['may_explicitly_override_promotion_gate'] is False
assert p['policy']['EXTERNAL']['auto_promotion'] is False
assert 'owner_approval' in p['policy']['EXTERNAL']['requirements']
PY
for mirror in .claude/skills/gatekeeper .agents/skills/gatekeeper; do
  diff -qr "$ROOT/skills/gatekeeper" "$ROOT/$mirror"
done
grep -q 'pull_request_review:' "$ROOT/.github/workflows/gatekeeper.yml"
grep -q 'LuigiD5555' "$ROOT/.github/workflows/gatekeeper.yml"
echo 'PASS gatekeeper policy/workflow/skill coherence'
