# Public API Naming Reset Design

## Goal

Reset the public GoSprite64 API to one coherent idiomatic-Go naming convention, even at the cost of breaking backward compatibility.

The intended outcome is:

- one naming philosophy across the public API
- clearer exported types, functions, and entrypoints
- removal of mixed naming dialects from the public contract
- examples, docs, generated code, and tests all teaching the same language

## Context

GoSprite64 currently exposes a public API with several naming styles mixed together.

Today the exported surface combines:

- idiomatic-Go style names such as `Print`, `GetStick`, `PlayEffect`, and `SetTrackVolume`
- retro shorthand names such as `Btn`, `Btnp`, `Rnd`, and `Flr`
- hybrid or framework-flavored names such as `Rectfill` and `Gamelooper`
- evolving startup naming questions such as `App` versus more domain-specific alternatives

This makes the library feel less coherent than it should.

The problem is not that any single name is unbearable on its own. The problem is that the public surface speaks more than one naming language at once.

## Decision

GoSprite64 will adopt one explicit naming rule for the public API:

**The public API should read like a normal Go library first, and only secondarily like a retro game framework.**

The design direction is:

- exported types should use clear domain nouns or standard Go role names
- exported functions should use clear verbs or descriptive nouns
- shorthand names that are not standard Go or widely accepted technical abbreviations should not remain part of the public API
- duplicate names for the same concept should be eliminated in favor of the idiomatic-Go choice
- the naming reset applies to the public contract as a whole, not only to one new API surface

This is an intentional public API reset, not a compatibility-preserving cleanup.

Implementation should follow a simplicity-first Go style in the spirit of Ken Thompson, Rob Pike, and Robert Griesemer:

- prefer the clearest exported name over the shortest nostalgic one
- avoid ornamental renaming that does not improve public comprehension
- use one consistent naming model across types, functions, examples, and generators
- make the public API explain itself without requiring users to learn a separate dialect

## Naming Policy

### Primary rule

A first-time Go developer reading a GoSprite64 program should experience it as idiomatic Go code using a game library.

That means:

- exported names should be readable in normal Go source
- names should reveal purpose without relying on prior exposure to fantasy-console APIs
- unusual suffixes, abbreviations, or compression should require strong justification to remain public

### What this rejects

The reset explicitly rejects:

- retro shorthand as the primary naming style for exports
- mixed pairs where one name is descriptive and another is legacy shorthand
- partial cleanup that leaves the public API speaking multiple naming dialects

## Public API Categories

The reset should classify the public surface and apply one rule per category.

### 1. Core gameplay types

These should use clear domain nouns.

Examples of desired direction:

- user gameplay type should read as `Game` in examples and docs
- interfaces should avoid awkward suffixes such as `-looper` unless that is truly the clearest role name

The goal is to stop using names that sound like framework internals when they are actually part of the user-facing API.

### 2. Startup and entrypoints

These names should describe startup clearly and explicitly.

Entry functions should read cleanly in `main.go` and should use ordinary Go library naming, not engine jargon or generic framework naming when a clearer domain name exists.

If a startup config type exists, its name should describe setup intent clearly and should not collide with the user’s gameplay type.

### 3. Gameplay helpers

Gameplay helpers should use full descriptive verbs or clear nouns.

This category includes drawing, input, timing, and math helpers exposed as part of the public package API.

Names in this category should:

- prefer full words over shorthand
- line up stylistically with other exported helpers
- avoid preserving retro spellings that conflict with Go naming expectations

### 4. Asset and generated-code setup surfaces

Public names emitted or required by generated code should sound like ordinary Go data and functions.

They should:

- describe the exported value clearly
- avoid encoding internal runtime details into public names unless necessary
- align with the same naming rules as handwritten API types

## Naming Rules

The reset should apply these concrete rules:

1. **Types are nouns**
   - examples: `Game`, `AudioBundle`, `Config`

2. **Functions are verbs or descriptive actions**
   - examples: `Run`, `RunGame`, `IsButtonPressed`

3. **No non-standard shorthand exports**
   - names like `Btn`, `Btnp`, `Rnd`, and `Flr` should not remain public exports

4. **No mixed dialect pairs**
   - if `Rectfill` and `DrawRectFill` compete, one idiomatic name should win

5. **Prefer readability in call sites**
   - exported names should make gameplay code and `main.go` self-explanatory

6. **Keep internal names out of the public contract**
   - public names should not reflect internal runtime machinery unless users truly need that concept

## Naming Table Requirement

The implementation plan must include an explicit old-to-new naming table.

Each renamed symbol should have:

- old public name
- new public name
- category
- short rationale

This is required so the reset happens systematically instead of through ad hoc renaming.

## Candidate Rename Direction

The spec does not need to finalize every symbol name yet, but it does establish the expected direction.

Likely rename targets include:

- `Gamelooper` -> an idiomatic replacement based on the final gameplay/startup model
- `Btn` -> a full-word input helper name
- `Btnp` -> a full-word input-edge helper name
- `Rnd` -> a full-word randomness helper name, if it remains public at all
- `Flr` -> a full-word math helper name, if it remains public at all
- `Rectfill` -> one canonical descriptive draw-helper name

The naming table in the plan should decide the exact replacements.

## Migration Strategy

This reset should be treated as one deliberate API break, not a long compatibility migration.

The expected strategy is:

- rename the public API to the new idiomatic-Go names
- update examples, docs, generated code, and tests in the same phase
- remove the old names from the public contract
- avoid maintaining a broad alias layer for the previous naming dialect

Temporary glue is allowed only when required to land the reset in one coherent implementation cycle, and it should be removed before the phase is declared complete.

## Documentation And Examples

Documentation and examples must teach only the new names after the reset lands.

This includes:

- `README.md`
- getting-started and guide docs
- example programs
- generated example code where relevant

The public story must not say one thing while the examples teach another.

## Generator Impact

Generated code is part of the public experience and must follow the same naming reset.

`cmd/audiogen` and any generated outputs should:

- emit the new public names
- stop teaching old names in generated code comments or emitted code shape
- remain stable and understandable to users reading generated files

## Testing Strategy

Verification should prove public naming consistency, not just build correctness.

### Public naming checks

Add focused tests or checks that ensure:

- old shorthand exports are gone from the intended public API
- generated code uses the new names
- examples compile against the new names only

### Example and docs consistency

At least one updated example should serve as the canonical proof that the new naming style reads cleanly in practice.

Docs and examples should be reviewed as part of the reset, not treated as follow-up cleanup.

### Build and generator verification

Existing build and generation workflows should continue to pass after the reset.

The naming cleanup must not leave the repository in a state where generators, examples, or docs lag behind the renamed API.

## Rollout Order

Implement in this order:

1. define the naming table and final replacements
2. rename core public types and entrypoints
3. rename public gameplay helpers
4. update generated-code outputs and expectations
5. update examples and docs to teach only the new names
6. add or update checks ensuring the old public names do not remain part of the supported contract

This order keeps the reset coherent and reduces the chance of teaching mixed names during transition.

## Non-Goals For This Phase

Do not expand this work into:

- unrelated behavioral refactors justified only by naming work
- broad internal renaming where the names are not part of the public contract
- preserving both the old and new naming dialects long-term
- a retro alias layer presented as part of the official API
- speculative redesign of package structure unless required by the naming model

The reset should improve the public contract, not become an excuse for unrelated churn.

## Final Position

The right naming improvement for GoSprite64 is a deliberate public API reset toward idiomatic Go naming, applied consistently across exported types, entrypoints, helpers, generated code, examples, and documentation, with one explicit naming table and without preserving the old shorthand dialect as an equal public API story.
