#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
"$ROOT/scripts/verify-lock.sh"
python3 - "$ROOT" <<'PY'
import json, pathlib, sys
root = pathlib.Path(sys.argv[1])
m = json.loads((root / "TONAL.json").read_text())
c = json.loads((root / "compatibility.json").read_text())
assert m['schema'] == 'tonal.composition.v2'
assert 'COMPONENT_REPOSITORIES_REMAIN_AUTHORITATIVE' in m['invariants']
assert 'TARGET_REPOSITORY_OWNS_TARGET_RELEASES' in m['invariants']
assert 'TONAL_COMPOSITION_NE_TARGET_PROMOTION' in m['invariants']
assert 'TONAL_IS_OPTIONAL_FOR_COMPONENT_RUNTIME' in m['invariants']
assert 'UNPINNED_TOOL_NE_COMPOSITION_COMPONENT' in m['invariants']
assert m['components']['tlaloc']['kind'] == 'development_tool'
assert m['components']['origami']['kind'] == 'target'
assert 'Blueprint Framework' in m['composition_model']['unlocked_examples']
assert all(x.lower() != 'blueprint framework' for x in m['components'])
for name, component in m['components'].items():
    assert component['verification'], f'{name}: verification must be declared'
    for check in component['verification']:
        assert isinstance(check['argv'], list) and check['argv']
        assert isinstance(check.get('cwd','.'), str)
assert m['snapshot']['physical_artifact'] == 'BUILT_BY_VERIFY_WORKFLOW'
assert m['snapshot']['builder_status'] == 'IMPLEMENTED'
assert c['schema'] == 'tonal.compatibility.r1'
assert c['current']['state'] == 'COMPOSITION_VERIFIED'
assert c['current']['target_capability_promotion'] is False
print('PASS generic composition invariants')
PY
