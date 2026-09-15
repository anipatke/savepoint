---
type: epic-design
status: audited
---

# E42: Load V2 work with stable identity

## Purpose

Establish the read and write boundary for V2 project records so every later V2 capability consumes one strict, indexed interpretation of Objectives and Tasks without losing authored content.

This is the V1 delivery epic for V2 product Objective O002. The current V1 lifecycle remains authoritative while it is built.

## What this epic adds

- Project schema detection from `.savepoint/config.yml`: explicit `schema_version: 2` selects V2, absence selects the transitional V1 reader, and unsupported or malformed explicit versions produce named diagnostics.
- V2 Objective and Task records with explicit global IDs, Task-to-Objective ownership, optional release metadata, and typed Task dependencies.
- A project-level index that resolves Objectives, Tasks, ownership, and references by identity rather than by directory-derived numbering.
- Strict structural diagnostics for duplicate or malformed IDs, path/record mismatches, missing ownership and dependency targets, self-dependencies, and reference cycles.
- A raw-document/typed-record boundary that lets managed writes change owned fields while retaining unknown frontmatter fields and authored Markdown bodies.
- Frozen V1 dispatch behind a schema-specific adapter so current projects remain readable during the transition without leaking V1 inference rules into V2 records.

## Components and files

| Module | Purpose |
|--------|---------|
| `internal/data/config.go` | Parse and classify project schema version independently from package, release, and upgrade-manifest versions. |
| `internal/data/parser.go` | Provide the shared frontmatter document boundary used by strict V2 decoding and preserving writes. |
| `internal/data/project.go` (new) | Define the loaded project/index API and dispatch to schema-specific discovery. |
| `internal/data/objective_v2.go` (new) | Define and validate V2 Objective records and source metadata. |
| `internal/data/task_v2.go` (new) | Define and validate V2 Task ownership and typed dependency references without changing the active V1 Task contract. |
| `internal/data/discover.go` | Retain transitional V1 discovery behind the project loader and add confined V2 record discovery. |
| `internal/data/dependency.go` | Build and validate the V2 Task and Objective reference graphs from global IDs. |
| `internal/data/write.go` | Patch managed fields through the parsed document while preserving unknown fields and authored body content. |
| `internal/data/lifecycle.go` | Keep V1 lifecycle defaulting isolated; reject unsafe V2 lifecycle values rather than healing them into completion-capable state. |
| `internal/doctor/checks.go` | Report loader/index diagnostics by stable name without duplicating parsing or reference rules. |
| `internal/data/*_test.go`, `internal/doctor/*_test.go` | Cover schema dispatch, records, identity graphs, preservation, path safety, and diagnostic integration using E41 fixtures and temporary projects. |

## Architectural delta

`internal/data` gains one project-loading boundary for all consumers. `LoadProject` (final signature settled during task breakdown) accepts the `.savepoint` root, detects the schema once, dispatches to an isolated V1 or V2 loader, and returns a project model whose Objective and Task indexes are keyed by globally unique IDs. V2 consumers must not rediscover records or infer ownership from path segments.

V2 identity is record-owned. Objective IDs match `O` plus at least three digits; Task IDs match `T` plus at least three digits. A Task carries exactly one `objective: O###` reference. Titles, paths, and optional release labels may change without changing identity. Stored paths are normalized project-relative paths; discovery rejects traversal, duplicate identity, case collisions that alias on the target filesystem, and symlink escape from the project root.

Dependencies are decoded as records of the form `{task: T###, requires: clear|accepted}`, with `clear` supplied only when the field is omitted. Objective dependencies are `O###` references. The index validates targets, ownership, self-reference, and cycles, but E42 does not interpret Check evidence or decide whether a dependency is satisfied; those gates belong to E43.

Parsing separates the source document from its typed projection. Each loaded record retains source path plus the parsed YAML node/body needed for a preserving rewrite. Managed writes update only fields owned by the operation, retain unknown YAML keys and their values, and leave authored Markdown body content unchanged. Exact archival byte preservation is an E45 migration responsibility; E42 must nevertheless prove that a load followed by a no-op write does not alter the file and that a managed-field write does not discard unknown content.

Schema errors and graph errors are typed/named data diagnostics with record path and identity context. Doctor renders those diagnostics; it does not recreate schema rules. Unsupported explicit schema versions fail closed. Missing V2 fields, unknown lifecycle values, or malformed dependency requirements are never converted to defaults that could grant progress or completion.

The V1 adapter remains available only for transitional live reads and E45 migration input. It retains current V1 semantics and types at its boundary. New V2 records and indexes do not broaden or silently redefine the active V1 Task model.

## Boundaries

**In scope:**

- Schema classification and version dispatch for V1 transition input and V2 projects.
- Strict V2 Objective/Task parsing, global identity, explicit ownership, reference indexing, and structural validation.
- Safe project-relative discovery and named data/doctor diagnostics.
- Lossless preservation of unknown frontmatter and authored bodies during E42-managed record writes.
- Characterization against the frozen V1 source fixtures produced by E41.

**Out of scope:**

- Check records, clearance freshness, owner acceptance, dependency satisfaction, completion authority, or lifecycle transition changes; E43 owns them.
- Issue records, Objective integration Checks, and issue reconciliation; E44 owns them.
- Migration preview/apply, archival mapping, backups, recovery, or schema activation; E45 owns them.
- V2 skills, scaffolding, asset upgrades, resume/Next, board redesign, or removal of live V1 readers.
- Exact-byte preservation of an entire pre-migration V1 tree; E42 preserves the content it manages, while E45 owns archival byte fidelity.

## Dependencies

- E41-migration-source-fixtures must provide the frozen V1 source fixtures used to characterize transitional dispatch and preservation behavior.

## Quality gates

- Schema tests distinguish absent version, valid V2, malformed version, and unsupported explicit version with stable diagnostic names.
- Fixture-backed and temporary-project tests cover duplicate/malformed IDs, moved Tasks retaining identity, explicit ownership, missing targets, path mismatch, traversal/symlink escape, dependency defaults, self-reference, and Task/Objective cycles.
- Parser/writer tests prove no-op stability and preservation of unknown YAML fields, YAML values, Markdown bodies, and supported line-ending form across managed edits.
- Lifecycle fixtures prove malformed or unknown V2 values cannot be healed into a completion-capable state and that V1 compatibility behavior remains isolated.
- Doctor tests prove it consumes data diagnostics and reports actionable path/identity context without writing project files.
- Implementation handoff requires focused package tests plus `make build && make test`, with the named cases and outcomes recorded in the Tasks.
- Epic closeout requires a fresh independent V1 epic audit; this planning session is not that audit. Current Guardrails apply, with STYLE findings advisory.

## Open decisions

None. Exact Go function/type names may be refined during task breakdown, but the ownership, dispatch, preservation, and diagnostic contracts above are fixed.
