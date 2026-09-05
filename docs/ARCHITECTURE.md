# Tonal Architecture R2

## 1. System definition

Tonal is the complete heterogeneous cognitive/runtime system and research authority for composing bounded capabilities into verified workflows.

Architecture R2 promotes the runtime boundary introduced during T1 from an experiment-local addendum to the canonical system architecture. The historical composition/pinning/provenance layer remains useful infrastructure, but it no longer defines Tonal as a whole.

## 2. Core runtime responsibilities

Tonal owns goal intake, workflow/DAG representation, capability registry consumption, capability selection/routing, scheduling/dependency coordination, Blackboard state, execution state, resource accounting, verification coordination, trace generation, final workflow results and experiment-level system metrics.

## 3. Capability is the common runtime abstraction

Tonal asks:

```text
Which eligible capability should execute this bounded operation?
```

A capability describes behavior, not component identity. Different component kinds may satisfy the same capability.

```text
CapabilityKind
├── TLALOQUE        reusable machinery produced/qualified through Tlaloc
├── MACHINE         compiled deterministic/composite routine
├── TOOL            external deterministic/operational tool
└── EXTERNAL_MODEL  external probabilistic cognition; Parrot is the canonical interface
```

`Verifier` is a role, not necessarily a component kind: verification may itself be implemented by deterministic machinery, a Tlaloque, a tool or another qualified mechanism.

Canonical separation:

```text
Capability          what behavior is requested/offered
CapabilityKind      what class of component supplies it
CapabilityProfile   evidence, competence envelope, cost/reliability history
Executor            performs the work
Verifier            judges result/evidence
SelectionPolicy     chooses among eligible candidates
```

## 4. Tlaloques are machinery

A Tlaloque is a bounded, typed, measurable reusable mechanism created or qualified through Tlaloc. It may be deterministic, algorithmic, symbolic, tool-backed, specialized-model-backed or hybrid, but it participates as reusable machinery with explicit contracts and evidence.

A model-backed Tlaloque is still distinct from Parrot: its bounded behavior is the artifact being qualified. The external general model itself is not a Tlaloque merely because Tonal can call it.

## 5. Parrot is external cognition

Parrot is Tonal's singular interface to external probabilistic cognition. It is not produced by Tlaloc and is not a Tlaloque.

A concrete provider/model/version is configuration beneath the Parrot interface. Tonal should not create architectural categories for Claude, GPT, Gemini, LFM or other providers when they play the same external-cognition role.

Parrot may be valuable for perception, ambiguity, interpretation, hypothesis generation or other operations where machinery is insufficient. Its output is a probabilistic candidate/observation, not automatic truth.

```text
machinery
   ↓
unresolved ambiguity / novelty
   ↓
PARROT
   ↓
probabilistic observation/candidate
   ↓
VERIFY
   ↓
COMMIT or reject
   ↓
machinery continues
```

Parrot is not the router, global planner, verifier, memory, semantic authority or mandatory executor. There is no implicit fallback from an unknown capability to Parrot.

## 6. Final response boundary: semantic authority vs expressive freedom

Parrot may have expressive freedom without having answer authority.

For a final user-facing response, the semantic content must already have been established by the internal Tonal process and committed as verified facts in the Blackboard. Parrot receives a read-only projection of those facts and renders them into natural language.

```text
Tlaloques / Machines / Tools / bounded Parrot operations
                  ↓
              VERIFY
                  ↓
         Blackboard FACTs
                  ↓
          ResponseComposer
                  ↓
 Parrot: wording / ordering / tone
                  ↓
        ResponseVerifier
                  ↓
        release to user
```

The invariant is:

```text
Blackboard decides WHAT may be asserted.
Parrot may decide HOW to express it.
```

Tone policy is separate from semantic policy:

- `AUTO`: Parrot may choose a natural tone/register using the interaction context and its pretrained conversational priors.
- `EXPLICIT`: Tonal or the user may request a particular tone.

Neither mode allows Parrot to add unsupported claims.

The final renderer:

- receives only explicitly selected verified Blackboard facts;
- has no write access to the Blackboard;
- may not emit new semantic observations during rendering;
- must identify which Blackboard keys it used;
- cannot release text unless a `ResponseVerifier` accepts semantic closure against the supplied facts.

This boundary does **not** forbid Parrot from participating earlier in the control loop for bounded perception, interpretation or hypothesis capabilities. It only prevents the final natural-language renderer from becoming a second hidden reasoning authority after Tonal has already established its answer state.

## 7. Tlaloc

Tlaloc is the capability foundry and Behavior Lab. It constructs, qualifies, packages, promotes and deprecates Tlaloques and other reusable machinery; produces CapabilityProfiles; runs behavior experiments; and studies verified Episodes for reusable structure.

Tlaloc may characterize Parrot empirically and may convert recurring successful Parrot behavior into a candidate Tlaloque/Machine. It does not own or promote the external model itself.

## 8. Shponglese

Shponglese begins as semantic operational IR for primitive and composed behavior:

```text
ATOM -> MOTIF -> GRAPH -> PACK
```

Its meaning must remain independent of codec. Claims of language require later empirical evidence for compositionality, systematicity, slot generalization, held-out composition and codec invariance.

## 9. Origami

Origami is an independently testable representation, transport, addressing and virtual-memory substrate. Tonal may use Origami where measured evidence shows value. Origami must not become an undeclared requirement for Tonal correctness and does not own Shponglese semantics.

## 10. Experience loop

Target long-term loop:

```text
novel/unresolved task
   ↓
existing machinery sufficient?
   ├── yes -> execute + verify
   └── no  -> selectively invoke Parrot
                  ↓
             verify outcome
                  ↓
                Episode
                  ↓
            Tlaloc analysis
                  ↓
      recurring reliable structure?
                  ↓
 candidate Tlaloque / Machine / motif / specialist
                  ↓
 qualification + holdout + ablation
                  ↓
               Registry
                  ↓
       less external cognition next time
```

This is a research target, not yet a demonstrated system property.

## 11. T1 compatibility

T1 remains scientifically frozen. Its historical Tlaloc R1 adapter described Parrot as a generative Tlaloque. R2 may adapt that frozen publication into `CapabilityKind=EXTERNAL_MODEL` at the Tonal boundary while preserving exact T1 routing, calls and accounting. This reclassification must not rewrite T1 artifacts or results.

## 12. Architectural invariant

> Use reliable reusable machinery where it is sufficient; spend external probabilistic cognition only where uncertainty or novelty requires it.

For final responses, the companion invariant is:

> Semantic authority lives in verified Blackboard state; Parrot contributes expression, not post-hoc facts.
