# EVIDENCE_PLAN — Plan de cierre de discrepancia

**Versión del plan:** 0.8.0
**Fecha de corte:** 2026-09-01
**Autoridad:** este documento es la fuente de verdad del *plan*. No es fuente de verdad del *estado*.
**Ámbito:** stack `tlaloc` + `origami` + `tonal`.
**Ubicación canónica:** `tonal/docs/EVIDENCE_PLAN.md` (es un asunto de nivel stack, no de un componente).

---

## Instrucciones para el agente

Si estás leyendo esto como agente (Claude Code, Codex, Cursor, Copilot, Gemini CLI), esta sección manda sobre el resto del documento.

### Orden de lectura

1. Esta sección.
2. `TAREA ACTUAL`, más abajo. Es lo único que hay que hacer ahora.
3. La sección del plan que esa tarea referencia.
4. Los archivos reales del repositorio.

No leas el documento completo antes de empezar. Cada tarea es autocontenida a propósito.

### Precedencia

Para el **plan**:

```
prompt explícito del usuario > esta sección > resto del plan > documentos históricos
```

Para el **estado**, la precedencia se invierte y es absoluta:

```
repositorio > ledger > cualquier documento
```

Si el plan y el repo se contradicen: **reporta la contradicción y detente.** No la reconcilies en silencio, no edites el repo para que coincida con el plan, ni edites el plan para que coincida con el repo sin decirlo. Esa reconciliación callada es exactamente el fallo que originó este documento.

### Prohibiciones duras

- **P1.** No escribas estado a mano en ningún `.md`. El estado sale del ledger (§4).
- **P2.** No promuevas el status de una afirmación sin cumplir su requisito de §4.4. Que el código "se vea correcto" no es evidencia.
- **P3.** No toques nada de §9. Si crees que hace falta, anótalo en §11 y sigue.
- **P4.** No añadas arquitectura, capas, nombres ni abstracciones nuevas. Este ciclo es de instrumentación, no de diseño.
- **P5.** No inventes `run_id`, hashes ni resultados. Si no corriste algo, el campo va vacío y lo dices.
- **P6.** No borres ni "limpies" documentos históricos. Son evidencia.

### Definición de terminado

- `go test ./...` y `go vet ./...` pasan en el repo tocado;
- lo producido existe como archivo commiteado, no como propuesta en el chat;
- si la tarea tocó el plan, hay una línea nueva en §12;
- dijiste explícitamente qué **no** hiciste y por qué.

Terminar a medias y decirlo es un resultado válido. Terminar a medias y no decirlo, no.

### TAREA ACTUAL

> **Semana 4 — Markov mínimo.**
>
> Contar transiciones únicamente desde `trace[]` de Run Records comparables, agrupadas por `(modelo, perfil)` y restringidas a corridas que compartan `env_hash` según §7.
> Estimar la matriz `P(s'|s)` por conteo directo, reportar siempre `n` junto a cada probabilidad y mantener el piso duro de `n >= 30` para cualquier transición que pueda alimentar una decisión.
> Emitir el primer reporte que responda las cuatro preguntas de §7.3 sin introducir un MDP, recompensas, acciones aprendidas ni Behavior Lab de tres niveles.
> Si el volumen real no alcanza el piso de conteo, declarar la capa insuficientemente evidenciada en lugar de bajar el umbral.

Cuando la tarea cambie, se cambia **aquí** y se anota en §12. Nunca hay dos tareas actuales.

### Al terminar, reporta en este orden

Qué cambió con rutas; qué comandos corriste con su salida; qué quedó como `designed` y por qué; qué preguntas nuevas abrir en §11.

---

## 0. Qué es este documento y qué no

**Es:**

- El plan completo para eliminar la discrepancia entre lo diseñado y lo implementado.
- El registro histórico de cómo ese plan cambia y por qué.
- El lugar donde se modifica cualquier decisión de este ciclo. Si algo cambia, cambia aquí primero.

**No es:**

- Un reporte de estado. El estado se **genera** desde el ledger (§4), nunca se escribe a mano aquí.
- Un sustituto de `HISTORY_AND_LINEAGE.md`. Aquel documento es genealogía conceptual y hace bien ese trabajo. Este es operativo y acotado a un ciclo.
- Un lugar para ideas nuevas de arquitectura. Ver §1.

---

## 1. Reglas de este documento

El documento maestro anterior creció por absorción: toda idea nueva recibía nombre, sección y etiqueta, y ninguna regla la eliminaba nunca. Este documento arranca con el mecanismo contrario.

**R1 — Tope duro.** Máximo 500 líneas. Al llegar al tope no se añade sin borrar.

**R2 — El estado no se escribe.** Ninguna sección afirma "X está implementado". Las afirmaciones viven en el ledger y traen `run_id`. Este documento describe qué se va a hacer y qué se decidió.

**R3 — Una idea nueva no entra por ser buena.** Entra si cambia lo que se hace en este ciclo. Si no, va a `EXPERIMENTAL_IDEAS.md` y no vuelve hasta que el ciclo cierre.

**R4 — Superseder es borrar.** Cuando una decisión reemplaza a otra, la anterior se elimina del cuerpo y queda una línea en §12. No se conservan dos versiones vivas.

**R5 — Sin diagramas decorativos.** Un diagrama entra si sustituye texto, no si lo acompaña.

**R6 — El plan también se verifica.** Si una sección de este plan lleva dos semanas sin producir ni un commit ni un `run_id`, se marca en §12 y se decide explícitamente entre continuarla o matarla.

---

## 2. El problema, medido

Estado observado tras `git fetch --prune origin` sobre los `main` publicados al 2026-09-01:

| Repo | VERSION | Commits | Commit observado |
|---|---|---:|---|
| `LuigiD5555/origami` | `6.0.0-alpha.15` | 404 | `88ddb8a` |
| `LuigiD5555/tlaloc` | `6.0.0-alpha.21` | 450 | `1afb315` |
| `LuigiD5555/tonal` | `0.1.0-alpha.5` | 47 | `3ce0871` |

El `tonal.lock` v2 de Tonal `0.1.0-alpha.5` fija Tlaloc `6.0.0-alpha.15@349deaef` y Origami `6.0.0-alpha.12@12176e68`. Que el lock vaya detrás de los `main` es válido: representa una composición fijada, no el estado de desarrollo de cada componente.

La medición de v0.1.0 quedó supersedida porque se hizo contra referencias locales atrasadas y los repos públicos avanzaron. Los números del documento histórico ahora coinciden con los archivos `VERSION` de Tlaloc y Origami, pero esa coincidencia no convierte sus afirmaciones de capacidad en evidencia.

**El problema real no es el número equivocado.** Es que existía un lugar donde se podía escribir un estado sin evidencia y nada lo rechazaba. Corregir el número a mano no arregla nada; en tres semanas vuelve a divergir. Lo que hay que construir es el mecanismo que hace imposible escribirlo.

Esto es especialmente incómodo porque el propio proyecto ya tiene la ley que se violó: `FALSE EXACT = 0` y `predicción ≠ verificación`. La ley existía y no estaba aplicada al proyecto mismo.

---

## 3. Tesis y orden forzado

Una cadena de Markov es, operativamente, conteo de transiciones sobre un espacio de estados definido. No se puede estimar `P(s'|s)` sin trayectorias registradas, y no se pueden registrar trayectorias comparables sin haber fijado el entorno.

El orden no es una preferencia:

```
1. LEDGER      estado declarado con evidencia   → mata el drift
2. RUN RECORD  entorno capturado                → hace replicable
3. TRACES      trayectorias por corrida         → hace contable
4. MARKOV      matriz de transición estimada    → hace automatizable
```

Saltarse el 2 para llegar al 4 es automatizar sobre ruido.

---

## 4. Capa 1 — Claims Ledger

### 4.1 Propósito

Un archivo por componente que enumera cada afirmación de capacidad y el nivel de evidencia que la sostiene. Sustituye a cualquier tabla de estado escrita a mano.

### 4.2 Ubicación

```
origami/state/CLAIMS.json
tlaloc/state/CLAIMS.json
tonal/state/CLAIMS.json
```

Origami ya tiene `state/ORIGAMI_STATE.json` y `changes/`. La intención es extender esa vía, no crear una paralela. Tlaloc ya tiene `docs/CAPABILITY_STATUS.md`, que hoy es prosa; se convierte en salida generada.

### 4.3 Esquema

```json
{
  "id": "ORIGAMI.ADDRESSING.PAGE_RESOLVE",
  "statement": "resolver página lógica a offset físico correcto",
  "status": "evidenced_failing",
  "evidence": ["run:2026-08-14T22:10:03Z:a91f", "test:TestPageResolve"],
  "version_introduced": "6.0.0-alpha.2",
  "last_checked": "2026-09-01",
  "notes": "falla asumiendo offset fijo página lógica = página impresa"
}
```

Campos obligatorios: `id`, `statement`, `status`, `version_introduced`, `last_checked`.
`evidence` es obligatorio para `implemented` en adelante.

### 4.4 Estados y regla de promoción

| status | requisito para declararlo |
|---|---|
| `designed` | existe la sección en el documento de diseño. Nada más. |
| `implemented` | existe un test que corre en `go test ./...` |
| `evidenced` | existe ≥1 `run_id` presente en `runs/` |
| `evidenced_failing` | existe ≥1 `run_id` y el veredicto es negativo (dato válido) |
| `verified` | comprobación determinista: hash, invariante o igualdad de bytes |
| `rejected` | está en `REJECTED_IDEAS.md` con causa |

**Regla dura:** ningún estado se promueve sin cumplir su requisito, y ningún estado se promueve por consenso de modelos. Consenso perceptual no es verificación.

`unknown` es un veredicto legítimo y nunca se convierte en `verified` por falta de contraejemplo.

### 4.5 Trabajo inicial concreto

Poblar el ledger a partir de lo que hoy existe, sin inventar:

- **Tlaloc**: `BehaviorSpec → PromptIR → prompt compilado`; propuestas de reparación acotadas con autoridad de promoción centralizada; comparación determinista contra semántica de referencia; transporte compatible con OpenAI; perfil `origami.quantum-inspired.r0`.
- **Tlaloc, como `designed`**: ejecución/evaluación de canales perceptuales, soporte nativo Anthropic, SkillIR, backends por familia de modelo. Su README ya los declara no implementados; el ledger solo lo formaliza.
- **Origami**: workbench Go R3.10-LAB, identidad determinista de experimentos, hashing de artefactos, fixtures, plumbing de fallo a regresión, DVR/PSS, Glyph Calculus, SAFE_MICRO_ISA, Context SIMD. Specs de semántica de estado y canales perceptuales entran como `designed` salvo que exista ejecutor.
- **Origami, evidencia real ya disponible**: prueba del libro, 7.55 MB → ~474 KB de carrier, recuperación y hash correctos, direccionamiento fallido. Entra como dos afirmaciones separadas, una `evidenced` y otra `evidenced_failing`. Esta separación es el hallazgo, no un accidente.
- **Tonal**: resolución de revisiones, verificación de coherencia lock/version, distribución de `repo-flow`.

Todo lo demás del documento histórico (Causal Dataflow, Behavior Lab de tres niveles, MDP de diseño, auto-Tlaloque) entra como `designed`. No es un juicio de valor sobre esas ideas.

---

## 5. Capa 2 — Run Record

### 5.1 Propósito

Capturar el entorno completo de cada corrida contra un modelo, de forma que sea replicable y comparable meses después.

### 5.2 Esquema

```json
{
  "run_id": "2026-09-03T14:22:11Z:7f3c9a",
  "env_hash": "sha256:...",
  "variable_axis": "model.id",

  "component": {
    "tlaloc": "6.0.0-alpha.10@6700b3bd",
    "origami": "6.0.0-alpha.3@978feef7",
    "tonal_lock": "0.1.0-alpha.2"
  },

  "model": {
    "provider": "local|openai-compatible|anthropic",
    "id_requested": "qwen2.5-7b-instruct",
    "id_reported": "<lo que devolvió el endpoint>",
    "quantization": "Q4_K_M",
    "context_window": 32768,
    "tokenizer": "<id o hash>",
    "endpoint": "http://localhost:8080",
    "observed_at": "2026-09-03T14:22:11Z"
  },

  "sampling": {
    "temperature": 0.0,
    "top_p": 1.0,
    "seed": 42,
    "max_tokens": 2048,
    "stop": []
  },

  "prompt": {
    "behaviorspec_id": "...",
    "promptir_hash": "sha256:...",
    "compiled_prompt_hash": "sha256:..."
  },

  "fixture": { "id": "...", "sha256": "sha256:..." },

  "host": {
    "os": "Manjaro 6.x",
    "cpu": "...", "gpu": "...", "ram_gb": 32,
    "go": "1.2x", "python": "3.1x"
  },

  "outcome": {
    "output_hash": "sha256:...",
    "parsed": true,
    "verdict": "verify_pass|verify_fail|parse_fail|unknown",
    "latency_ms": 0,
    "tokens_in": 0,
    "tokens_out": 0
  },

  "repetitions": { "n": 1, "verdict_distribution": {} },
  "replay": "tlaloc lab run --spec ... --fixture ... --model ...",
  "trace": []
}
```

### 5.3 `env_hash` y `variable_axis`

Son los dos campos que hacen todo el trabajo.

`env_hash` es el hash de todo el registro **excepto** el eje declarado en `variable_axis` y el bloque `outcome`.

**Regla de comparabilidad:** dos corridas solo se comparan si comparten `env_hash`. Lo que difiere entre ellas es exactamente lo declarado en `variable_axis`. Sin esto, en dos meses hay 400 corridas incomparables y ninguna conclusión defendible.

### 5.4 Riesgos conocidos

**Deriva de modelos remotos.** Un proveedor puede cambiar pesos detrás del mismo identificador. Por eso se guardan `id_reported` y `observed_at`, y toda comparación entre fechas distintas se trata como sospechosa por defecto. Los modelos locales cuantizados son reproducibles de verdad y funcionan como control: **al menos un modelo local en cada campaña.**

**Temperatura mayor que cero.** Con muestreo estocástico, una corrida no es un dato. Se exige `repetitions.n ≥ N` y se guarda la distribución de veredictos, no el veredicto. `N` se fija empíricamente en el ciclo; punto de partida 5 para exploración, 30 para cualquier número que alimente §7.

**Corridas deterministas.** Con `temperature = 0` y semilla fija en modelo local, `output_hash` debe reproducirse. Es la prueba de que el arnés funciona; si no reproduce, el arnés está roto y nada de lo demás vale.

### 5.5 Layout

```
runs/
  2026-09/
    2026-09-03T14-22-11Z-7f3c9a.json
  index.jsonl          # una línea por corrida, para consulta barata
```

Inmutable. Una corrida nunca se edita. Una corrida errónea se marca con otra que la refiere.

---

## 6. Capa 3 — Trayectorias

### 6.1 Espacio de estados

No es el estado del proyecto. Es el estado de **una corrida**, y sale casi directo del pipeline actual de Tlaloc:

```
PENDING
  → PROMPT_COMPILED
  → MODEL_ANSWERED
      ├─ PARSE_FAIL → REPAIR_PROPOSED → REPAIR_APPLIED → MODEL_ANSWERED
      └─ PARSE_OK → VERIFY_PASS
                  | VERIFY_FAIL
                  | UNKNOWN
                       └─ PROMOTED | ABANDONED
```

Estados absorbentes: `PROMOTED`, `ABANDONED`.

### 6.2 Evento de transición

```json
{ "t": 3, "from": "PARSE_FAIL", "to": "REPAIR_PROPOSED",
  "at": "...", "latency_ms": 0, "actor": "tlaloque:repair" }
```

Se acumulan en `trace[]` dentro del run record. No hay archivo separado: la trayectoria pertenece a la corrida que la produjo.

---

## 7. Capa 4 — Markov

### 7.1 Qué se estima

Para cada par `(modelo, perfil)`, la matriz `P(s'|s)` por conteo directo de transiciones sobre las trayectorias de §6, restringido a corridas que comparten `env_hash`.

### 7.2 Piso de conteo

**Ninguna transición con `n < 30` observaciones se usa para decidir.** Por debajo de ese umbral, el controlador cae a la política por defecto. Sin este piso, la automatización estaría reaccionando a ruido con apariencia de estadística.

Las probabilidades se reportan siempre con su `n` al lado. Un `0.87 (n=4)` no es un `0.87`.

### 7.3 Qué preguntas responde

Números que hoy no existen y que valen más que arquitectura nueva:

- `P(VERIFY_PASS | REPAIR_APPLIED)` — ¿la reparación acotada de los Tlaloque sirve, o solo gasta tokens?
- Número esperado de rondas hasta absorción — da `max_epochs` **medido** en lugar del 3 supuesto por diseño.
- Probabilidad de absorción en `ABANDONED` — qué fracción de tareas nunca cierra.
- Comparación de matrices entre modelos — dónde falla cada uno, no cuál es "mejor".

### 7.4 Qué NO es

**No es un MDP.** Un MDP requiere acciones, recompensa y aprendizaje de política. Aquí solo hay conteo de transiciones y una regla de decisión sobre ellas. Es suficiente para automatizar y no requiere nada que no se pueda auditar a mano.

**No es Behavior Lab.** El Behavior Lab de tres niveles del documento histórico queda `designed`. Esta capa es su precondición, no su implementación.

---

## 8. Criterios de cierre del ciclo

El ciclo cierra cuando las cinco condiciones se cumplen simultáneamente:

**C1.** Un comando de Tonal sale con código distinto de cero ante cualquier desacuerdo entre sus superficies `VERSION`, `README`, `TONAL.json`, `tonal.lock` y ledger. Para cada componente fijado por `tonal.lock`, el checkout exacto debe hacer coincidir `VERSION` y el encabezado del README con el lock. Si ese pin histórico contiene `state/CLAIMS.json`, el ledger y su tabla generada también se validan; la ausencia del ledger en una revisión anterior a su introducción no invalida retrospectivamente una composición fijada.

**C2.** Toda afirmación con status `evidenced` o `evidenced_failing` apunta a al menos un `run_id` que existe en `runs/`. Verificable por script.

**C3.** Una corrida cualquiera se re-ejecuta desde su propio `replay` y reproduce el mismo `env_hash`; en rutas deterministas, el mismo `output_hash`.

**C4.** La tabla de estado de cualquier documento se genera desde el ledger. Escribir estado a mano rompe CI.

**C5.** Ninguna transición se usa para decidir sin superar el piso de conteo.

Mientras C1–C3 no se cumplan, C5 no tiene sentido: no hay datos que contar.

---

## 9. Fuera de alcance en este ciclo

Explícitamente congelado. No es rechazo, es secuencia.

- Frente perceptual completo: moiré, estéreo, parallax, `KINETIC_REVEAL`, Native, las 42 familias generadoras.
- Auto-creación de Tlaloque y destilación.
- Behavior Lab de tres niveles y Design Lab.
- Causal Dataflow como runtime.
- Backends por familia de modelo, SkillIR, soporte nativo Anthropic.
- Empaquetado de snapshot portable del stack en Tonal.

Todos dependen del registro que aún no existe. Ninguno se puede evaluar sin él.

---

## 10. Secuencia por semanas

**Semana 1 — Ledger.** Esquema, poblado inicial desde los tres repos según §4.5, script de validación, generador de la tabla de estado. Salida observable: la tabla de `CAPABILITY_STATUS` deja de estar escrita a mano.

**Semana 2 — Run record.** Esquema, emisión desde `behavior-lab`, `env_hash`, layout de `runs/`, prueba de reproducibilidad determinista (C3). Salida observable: una corrida replicada desde su propio JSON.

**Semana 3 — Trayectorias y gate.** Emisión de `trace[]`, C1 y C2 en CI de Tonal. Salida observable: un commit que rompe CI a propósito por declarar estado sin evidencia.

**Semana 4 — Markov mínimo.** Conteo, matriz por `(modelo, perfil)`, piso de conteo, primer reporte con las cuatro preguntas de §7.3. Salida observable: `max_epochs` medido.

La secuencia es la apuesta, no una promesa. Si la semana 2 se estira, la 4 se recorta antes que saltarse la 3.

---

## 11. Preguntas abiertas

Se resuelven dentro del ciclo y se registran en §12 cuando se cierren.

- **P2.** ¿Qué `N` de repeticiones para muestreo estocástico? Punto de partida 5 / 30. Sin medir.
- **P3.** ¿El piso de 30 observaciones aguanta con el volumen real de corridas que se pueden generar en semanas? Si no, la capa 4 se pospone y se declara así, en vez de bajar el piso.

---

## 12. Historial de cambios

Toda modificación al plan se registra aquí. Una línea, con fecha y causa.

| Fecha | Cambio | Causa |
|---|---|---|
| 2026-09-01 | Creación del plan v0.1.0 | Discrepancia detectada entre documento histórico (Tlaloc alpha.21 / Origami alpha.15) y repos (alpha.10 / alpha.3) |
| 2026-09-01 | v0.2.0: bloque de instrucciones para agentes al inicio, con `TAREA ACTUAL` como puntero único | El plan se va a entregar a agentes (Claude, Codex) para ejecutarlo |
| 2026-09-01 | v0.2.0: §13 reescrita con la especificación real de `AGENTS.md` y la distribución `.claude/` + `.agents/` | La convención de raíz es `AGENTS.md`, no un archivo por herramienta |
| 2026-09-01 | v0.3.0: plan materializado en Tonal, §2 actualizado contra `origin/main` y `TAREA ACTUAL` avanzada al validador/generador de Tlaloc | Las referencias locales usadas por v0.1.0 estaban atrasadas; el ledger inicial de Tlaloc ya permite probar el siguiente gate |
| 2026-09-01 | v0.4.0: `TAREA ACTUAL` avanzada al ledger de Origami | Tlaloc ya valida su ledger, genera `CAPABILITY_STATUS` y rechaza divergencia en CI mediante el commit local `9ba1a22` |
| 2026-09-01 | v0.5.0: `TAREA ACTUAL` avanzada al ledger de Tonal y las contradicciones pasan a listas estructuradas pendientes | Origami ya separa 21 capacidades implementadas, 9 diseñadas y 5 discrepancias pendientes mediante el commit local `15820f3` |
| 2026-09-01 | v0.6.0: cerrada Semana 1, ledger por componente con validación local y Tonal como verificador de composición; `TAREA ACTUAL` avanzada al Run Record base | Tlaloc, Origami y Tonal ya validan sus ledgers y rechazan deriva de las tablas generadas; se evita depender de un checkout vecino para validar estado local |
| 2026-09-01 | v0.7.0: `TAREA ACTUAL` avanzada a trayectorias y gates C1/C2 | Tlaloc ya emite Run Records inmutables desde `realcampaign.Prepare`, calcula `env_hash`, mantiene `index.jsonl` y verifica replay determinista con el mismo `env_hash` y `output_hash` |
| 2026-09-01 | v0.8.0: cerrada Semana 3 y `TAREA ACTUAL` avanzada a Markov mínimo; C1 aclara la semántica de pins históricos | `realcampaign.Prepare` ya registra su transición observada en `trace[]`; Tonal rechaza deriva de superficies y `run_id` inexistentes, mientras composiciones históricas anteriores al ledger siguen siendo válidas si sus superficies versionadas coinciden |

---

## 13. Dónde viven las instrucciones para agentes

### `AGENTS.md` es la raíz

`AGENTS.md` es un formato abierto: markdown plano en la raíz del repositorio que los agentes leen a contexto antes de trabajar. No tiene esquema ni campos obligatorios; se usa cualquier encabezado y el agente parsea el texto. Lo lee un ecosistema amplio (Codex de OpenAI, Cursor, Copilot, Jules, Aider, Zed, Gemini CLI, Windsurf, entre otros) y hoy lo administra la Agentic AI Foundation dentro de la Linux Foundation.

Cuatro propiedades del formato que importan para este plan:

- **Anidamiento.** El `AGENTS.md` más cercano al archivo editado gana, así que un repo puede tener varios por subproyecto.
- **Precedencia.** Un prompt explícito del usuario sobrepasa cualquier `AGENTS.md`.
- **Ejecución.** Si listas comandos de prueba, el agente intentará ejecutarlos y arreglar los fallos antes de terminar. Aquí eso es deseable: `go test ./...`, `go vet ./...` y el verificador del ledger deben estar listados.
- **Es documentación viva.** Se actualiza cuando el proyecto cambia.

### Distribución en este stack

```
<repo>/AGENTS.md          fuente de verdad de reglas de agente
<repo>/CLAUDE.md          apunta a AGENTS.md, no lo duplica
<repo>/.claude/skills/    skills de tarea
<repo>/.agents/skills/    espejo byte-idéntico
```

La maquinaria ya existe y no hay que inventarla: Tlaloc tiene `CLAUDE.md` y cinco skills en `.claude/skills/`; Tonal mantiene `skills/repo-flow/` como canónico con espejos verificados en `.claude/skills/` y `.agents/skills/`, más `scripts/sync-skills.sh`, `scripts/install-skill.sh` y `tests/test-skills.sh`. Lo que falta es el `AGENTS.md` de raíz, y que `CLAUDE.md` deje de poder ser una fuente paralela.

### Regla de no duplicación

Tlaloc ya declara que sus skills commiteadas no son una segunda fuente semántica de verdad. Lo mismo rige aquí: `AGENTS.md` **referencia** este plan y el ledger, nunca los copia. Un esquema replicado en `AGENTS.md`, `CLAUDE.md` y dos árboles de skills son cuatro copias que van a divergir, que es literalmente el problema que este ciclo existe para resolver.

### Cuándo escribirlos

Las skills, después de la semana 2, cuando el esquema del run record deje de moverse.

El `AGENTS.md` mínimo, desde ya: solo los comandos de verificación y un puntero a `docs/EVIDENCE_PLAN.md`. No puede quedar obsoleto porque no contiene esquema.
