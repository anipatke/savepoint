---
type: commands-reference
triggerable: false
---

# Shared Savepoint Commands And Procedures

This reference is not a skill and never triggers on its own. `savepoint-design`
reads it when reconciling a migrated project's commands and procedures,
without restating the mapping here.

## Config Contract

`config.yml`'s `quality_gates` key is the single place project commands
live: `lint`, `typecheck`, `build`, `test`, `block_on_failure`, and
`gate_timeout` — the field names `internal/data/config.go`'s `QualityGates`
struct already decodes. No parallel `checks.technical` surface is introduced
alongside it; a project's lint, typecheck, build, and test commands live in
exactly one place.

`build` is an explicit, optional gate command alongside the existing `lint`,
`typecheck`, and `test` gates. When configured, `internal/doctor`'s
`RunQualityGates` runs it in the project root with the same `gate_timeout`
and `block_on_failure` semantics as the other gates, ordered after
`typecheck` and before `test`.

## Three-Part Mapping From A Legacy Health-Check.md

A V1 project's `Health-Check.md` mixes three different things. Reconciling
one after migration splits it exactly three ways — never left as one
undifferentiated file, and never given a fourth home:

1. Lint, typecheck, build, and test commands move into `config.yml`'s
   `quality_gates`.
2. The shared Quick/Full evidence mechanics move nowhere new — they are
   already written once, for every project, in
   `agent-skills/references/check-method.md`. Quick is optional for a requested
   Task Check; Full is mandatory for Objective integration checks.
   Do not copy that method's prose into a project file.
3. Whatever remains — project-specific prose describing a manual
   verification workflow that has no home in either of the above — becomes
   one preserved optional procedure file the project references.

### What migrate actually does

`internal/migrate`'s `Plan` archives a V1 project's `Health-Check.md`
byte-for-byte under `.savepoint/archive/v1/...` and separately records the
fenced-code-block lines it contains as `ArchiveEntry.CandidateCommands`, for
the preview to list as candidates. It never writes to `config.yml` and never
creates a procedure file itself: turning a candidate command into a live
`quality_gates` entry, and turning the remaining prose into a preserved
procedure file, is reconciliation work done by hand after migration — using
the archived original and its candidate commands as the starting material —
not something migration performs automatically.

## The Preserved Procedure Is Ordinary, Never A Default

The procedure file that survives step 3 is an ordinary project file like any
other: a Task names it in its own Context Files or Technical Verification
section when it needs it. It is never installed as a mandatory scaffold
default, and no project is expected to have one.

## Absence Is Not A Gap

A project with no `Health-Check.md` — this repository is one — needs no
generated substitute. Its absence is not a finding, not an Issue, and not
something a Check or an audit should flag as missing.

## Extra Verification Belongs To The Task

When a Task needs verification beyond the shared check method and the
configured gates, it names that in its own `## Technical Verification`
section. The Task still records implementation evidence even when its optional
Task Check is waived. It does not edit `agent-skills/references/check-method.md`,
`config.yml`'s `quality_gates`, or any other shared policy file to add a
one-off requirement.
