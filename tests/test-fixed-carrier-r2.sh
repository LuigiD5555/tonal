#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
python3 - "$ROOT/proposals/FIXED_CARRIER_R2.json" <<'PY2'
import json,re,sys
p=json.load(open(sys.argv[1]))
assert p['schema']=='tonal.proposal.r0'
assert p['id']=='fixed-carrier-r2'
assert p['status']=='COMPONENTS_MERGED_CI_GREEN_PENDING_EMPIRICAL_PROMOTION'
assert p['origami']['required_version']=='6.0.0-alpha.5'
assert p['tlaloc']['required_version']=='6.0.0-alpha.11'
for name in ('origami','tlaloc'):
    sha=p[name]['merged_commit']
    assert re.fullmatch(r'[0-9a-f]{40}',sha),(name,sha)
    assert p[name]['ci']=='GREEN'
assert p['supported_claim'] is False
assert 'FALSE_EXACT=0' in p['invariants']
assert p['remaining_empirical_gates']
print('PASS Fixed Carrier R2 candidate metadata')
PY2
