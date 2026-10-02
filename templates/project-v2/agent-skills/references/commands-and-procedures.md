---
type: commands-reference
triggerable: false
---

# Shared Savepoint Commands And Procedures

This reference is not a skill and never triggers on its own. `savepoint-design`
reads it when reconciling a migrated project's commands and procedures,
without restating the mapping here.

## Config Contract

`config.yml`'s `quality_gates` owns project commands: `lint`, `typecheck`, `build`, `test`, `block_on_failure`, and `gate_timeout`. No parallel `checks.technical` surface is introduced.

`build` is an explicit, optional gate command. `internal/doctor`'s `RunQualityGates` runs it in the project root after `typecheck` and before `test`, using the same timeout and failure semantics.

## Three-Part Mapping From A Legacy Health-Check.md

Reconcile a migrated V1 `Health-Check.md` into exactly three homes:

1. Lint, typecheck, build, and test commands move into `config.yml`'s
   `quality_gates`.
2. Quick/Full mechanics stay in `agent-skills/references/check-method.md`; do not copy its prose into project files.
3. Remaining project-specific manual verification prose becomes one preserved optional procedure file.

### What migrate actually does

`internal/migrate`'s `Plan` archives `Health-Check.md` byte-for-byte under `.savepoint/archive/v1/...` and lists fenced-code lines as `ArchiveEntry.CandidateCommands` in preview. It never writes `config.yml` or creates a procedure file. Manually reconcile commands and remaining prose from that archived evidence; this is not something migration performs automatically.

## The Preserved Procedure Is Ordinary, Never A Default

The procedure is an ordinary project file, named by Tasks in Context Files or Technical Verification when needed. It is never installed as a mandatory scaffold default; no project is expected to have one.

## Absence Is Not A Gap

A project without `Health-Check.md` needs no generated substitute; absence is not a finding or an Issue.

## Extra Verification Belongs To The Task

Put one-off verification in the Task's `## Technical Verification`, with implementation evidence even when its Task Check is waived. Do not add a one-off requirement to `agent-skills/references/check-method.md`, `config.yml`, or another shared policy file.
