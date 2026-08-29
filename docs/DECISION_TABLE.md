# Gatekeeper decision table

| Provenance | Technical gates | Owner approval | Override | Promotion |
|---|---|---|---|---|
| OWNER | PASS | implicit authority | not needed | ALLOW |
| OWNER | FAIL/UNKNOWN | n/a | explicit owner decision | OWNER_OVERRIDE or stop |
| EXTERNAL | PASS | APPROVED required | forbidden | ALLOW after approval |
| EXTERNAL | PASS | missing | forbidden | DENY |
| EXTERNAL | FAIL/UNKNOWN | any | forbidden | DENY |
