# Owner override

An owner override is a deliberate promotion decision available only to `LuigiD5555` under Gatekeeper R0.

It does not change a technical gate result. If CI is failing, the evidence remains failing. The override means the owner accepts the known risk and chooses promotion anyway.

Recommended record in the PR/merge context:

```text
OWNER_OVERRIDE
Gate: <failed/pending gate>
Reason: <why promotion is intentional>
Risk accepted: <known consequence>
Follow-up: <repair or none>
```

External contributors cannot issue this decision.
