# Tonal ecosystem boundaries — Architecture R2

## Tonal owns

Tonal owns runtime authority for goal intake, workflow/DAG representation, capability discovery/selection, routing/scheduling, Blackboard state, dependency coordination, execution state, resource accounting, verification coordination, tracing, final workflow results and system-level experiment metrics.

## Tlaloc owns

Tlaloc owns the lifecycle of reusable machinery:

- bounded Tlaloque construction;
- behavior decomposition and Behavior Lab campaigns;
- capability contracts;
- executor qualification;
- CapabilityProfile evidence;
- candidate discovery/comparison;
- promotion/deprecation of Tlaloc-managed machinery;
- Episode analysis and future experience-to-structure compilation.

Tlaloc does not own Tonal's scheduler, Blackboard or external general-purpose model.

## Tlaloque boundary

A Tlaloque is a bounded, typed, measurable reusable mechanism produced or qualified through Tlaloc. A Tlaloque may be deterministic, algorithmic, symbolic, tool-backed, specialized-model-backed or hybrid.

The defining property is not whether it contains a model. The defining property is that the bounded mechanism itself has a contract, evidence and lifecycle under Tlaloc.

## Parrot boundary

Parrot is **not a Tlaloque**. Parrot is Tonal's singular external probabilistic cognition interface.

A concrete provider/model/version lives beneath that interface as runtime configuration. Parrot may satisfy capabilities involving ambiguity, perception, interpretation, generation or novelty, but its output is an observation/candidate requiring verification.

Parrot does not own system planning, routing, verification, memory, Blackboard state, semantic truth or capability promotion. There is no implicit fallback that sends every unsupported capability to Parrot.

Tlaloc may characterize Parrot behavior and may turn recurring verified Parrot work into a candidate Tlaloque or Machine. It does not promote the external model itself.

### Final-response sub-boundary

Parrot's user-facing response role is narrower than its possible internal cognition roles.

Once Tonal has resolved a task, verified results and committed final facts to the Blackboard, Parrot may render those facts into natural language. At this boundary:

- Blackboard facts define what may be asserted;
- Parrot may choose wording, ordering and tone;
- `AUTO` tone may rely on the interaction and Parrot's pretrained conversational priors;
- an explicit tone may be requested by the user or Tonal policy;
- Parrot receives only an explicit read-only projection of verified facts;
- Parrot cannot append observations or facts during response rendering;
- rendered text must identify its grounding keys;
- rendered text is not released until a `ResponseVerifier` accepts semantic closure against those facts.

In short:

```text
semantic freedom: NO
expressive freedom: YES
```

## Capability boundary

`Capability` is the common Tonal abstraction. It must not collapse component identity.

```text
Capability
  supplied by one of:
    TLALOQUE
    MACHINE
    TOOL
    EXTERNAL_MODEL   <- Parrot
```

`Verifier` is a role that can be implemented by one of those component kinds.

## Shponglese boundary

Shponglese owns semantic operational representation for primitive and composed behavior. Its first role is an IR, not a claim of emergent language. Semantics should remain stable across codecs.

## Origami boundary

Origami owns representation, carrier, addressing, selective unfolding and virtual-memory mechanisms that it independently validates. Origami may carry Shponglese or other semantic structures, but it does not own Shponglese semantics.

## Dependency direction

```text
Tlaloc
  produces qualified Tlaloques/Machines
          │
          ▼
      Tonal Registry  <──── Parrot adapter / external model
          │
          ▼
SelectionPolicy -> Executor -> Verifier -> Blackboard/Episode
                                      │
                                      ▼
                              ResponseComposer
                                      │
                                      ▼
                          Parrot expressive rendering
                                      │
                                      ▼
                              ResponseVerifier
                                      │
                                      ▼
                                     User
          │
          ▼
        Tlaloc
  learns from verified Episodes
```

## Core invariant

Use reliable reusable machinery when sufficient. Invoke external probabilistic cognition only when unresolved uncertainty or novelty justifies it.

For final responses, semantic authority stays in verified Blackboard state; Parrot contributes expression, not post-hoc facts.
