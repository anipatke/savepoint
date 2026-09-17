---
id: E45-safe-migration/T002-take-stock-of-what-exists
title: Take stock of what exists
status: planned
objective: Inventory every project-owned source file by exact bytes and classify each one against the frozen fixture role vocabulary.
depends_on: []
complexity_tier: medium
complexity_reason: New package doing read-only filesystem walking and classification, with no writes and no cross-module ripple.
---

# T002: Take stock of what exists

## Problem

Migration cannot promise that every user file has an accountable destination until it can say, from bytes rather than from a parse, exactly which files exist and what each one is.

V1's loader heals as it reads: unknown status becomes `planned`, `stage: implementation` becomes `build`, a missing title falls back to the objective. A hash taken from a re-marshalled parse would therefore record what the loader wished the user had written. The inventory has to record what the user actually wrote, because that same hash is later the freshness basis that decides whether a reviewed preview still describes the project.

Classification has a ready-made contract. The two frozen E41 fixtures carry hand-authored `manifest.yml` files naming every source file's role, scoped ID, raw status, raw phase, and expected classification. Checking against those manifests is what keeps classification from being checked against itself.

## Context Files

- `internal/migrate/inventory.go`
- `internal/migrate/inventory_test.go`
- `internal/migrate/classify.go`
- `internal/migrate/classify_test.go`
- `internal/data/testdata/migration/README.md`
- `internal/data/testdata/migration/v1-basic/manifest.yml`
- `internal/data/testdata/migration/v1-history/manifest.yml`
- `internal/data/migration_source_test.go`
- `internal/data/discover.go`
- `internal/init/manifest.go`
- `.savepoint/Guardrails.md`

## Acceptance Criteria

- [ ] The inventory records, for every source file, a project-relative normalized path, an exact-byte SHA-256 taken from the file contents, the byte size, and the file mode.
- [ ] Hashes are never computed from a parsed and re-marshalled record, and a test proves a fixture file whose raw status the V1 loader would heal still hashes as its authored bytes.
- [ ] The walk covers project-owned `.savepoint/` content plus the managed `AGENTS.md` and `agent-skills/`, and excludes `.savepoint/.migration/`, which is the operation's own state rather than project source.
- [ ] Path confinement runs before anything else: traversal outside the project root, symlink escape, and case collisions are each rejected with a distinct named diagnostic, reusing the confinement rules `internal/data/discover.go` already established rather than inventing a second set.
- [ ] Paths are built with `filepath` joins and stored in normalized project-relative form, so a Windows target produces the same inventory keys as a Unix one.
- [ ] Every file named in both fixture manifests is classified with the role that manifest records, and the test reads the manifests as the expectation rather than regenerating them.
- [ ] The role vocabulary covers `config`, `router`, `product-prd`, `architecture`, `health-check`, `release-prd`, `epic-detail`, `epic-audit`, `task`, `defect`, `audit-prompt`, `audit-register`, `finding`, `audit-run`, the managed guide, and shipped skills.
- [ ] A file matching no known role is inventoried and reported as unclassified; it is never dropped, never guessed at, and never silently archived without appearing in the report.
- [ ] An absent optional artifact — no `Health-Check.md`, no `audit/`, no `AGENTS.md`, as `v1-basic` and `v1-history` respectively freeze — produces no finding, because absence is the expected V1 shape.
- [ ] The inventory is returned in a deterministic order, so every consumer downstream of it is deterministic for free.
- [ ] Inventory and classification perform no write: a test snapshots bytes and modification times across a full run of both fixtures and asserts nothing changed, on the success path and on each named failure path.

## Implementation Plan

- [ ] Create the `internal/migrate` package with `inventory.go`: the source-file record type, the confined walk, and content hashing.
- [ ] Reuse the existing confinement approach from `internal/data/discover.go` and the content-hash helper shape from `internal/init/manifest.go` rather than writing a third variant of either.
- [ ] Add `classify.go` with the role type and the path-and-shape rules that assign it, keeping the vocabulary in data rather than spread through conditionals.
- [ ] Add fixture-driven tests that read both `manifest.yml` files and assert role agreement file by file, naming any file the manifest covers that classification misses.
- [ ] Test confinement rejections, the unclassified path, absent optional artifacts, deterministic ordering, and hash fidelity against a loader-healed record.
- [ ] Add the byte-and-mtime snapshot assertion around a full inventory run of both fixtures.
- [ ] Run the focused `internal/migrate` and `internal/data` suites, then `make build && make test`.

## Context Log

Pending.
