---
type: audit-findings
audited: 2026-05-03
---

# Audit Findings: E13 Codebase Audit Remediation

## Main Findings

Applied outcome: E13 is closed as audited. The audit proposals were applied: the unused `readEpicAuditFile()` helper was removed, `Design.md` now records the E13 architecture delta, `AGENTS.md` describes the updated module responsibilities, and `E13-Detail.md` includes an "Implemented As" reconciliation note.

The critical and high-priority remediation work is present in code: cycle detection now reconstructs cycles from a DFS stack, AtomicWrite has a copy-based rename fallback, theme accent defaults fill missing keys individually, quality gates use `exec.CommandContext` with `gate_timeout`, frontmatter/line-ending helpers were centralized, board I/O was moved behind Bubble Tea commands, buildtool/styles tests were added, and committed binary paths are no longer tracked.

Verification performed during audit: `go build ./...` passed, and `go test ./...` passed across all packages. `golangci-lint run` could not be verified because `golangci-lint` is not installed in this environment.

Residual notes: T003, T004, and T007 task files still contain unchecked acceptance criteria despite implementation being present, and T005 did not meet the original `update.go` 30% line-count reduction criterion. Those are documented process/criterion drift, not remaining code blockers.

## Code Style Review

- [x] One job per file — New helper files stay within existing package responsibilities; `internal/board/io.go` owns board command I/O wrappers.
- [x] One job per function — The unused audit-read helper was removed during audit apply.
- [x] Test branches — New tests cover timeout handling, config accent default fill, accurate cycle paths, buildtool behavior, styles completeness, and async board update behavior.
- [x] Types document intent — Repair matching now uses sentinel errors and `errors.Is()`.
- [x] Build only what is needed — No new user-facing features or speculative command surface were introduced.
- [x] Handle errors at boundaries — Quality gate execution, board write commands, reload failures, and AtomicWrite fallback now surface boundary errors.
- [x] One source of truth — Line-ending normalization, frontmatter splitting, ID shortening, slice lookup, and board text helpers were consolidated.
- [x] Comments explain WHY — Comments are sparse and mostly orient around non-obvious behavior.
- [x] Content in data files — No new hardcoded product copy beyond command/status messages.
- [x] Small diffs — The epic is broad, but each task stayed scoped to the audit remediation area.

## Proposed Changes

### Target File
internal/board/update.go

### Replace
```go
func readEpicAuditFile(epicDir, shortID string) string {
	raw, err := os.ReadFile(filepath.Join(epicDir, shortID+"-Audit.md"))
	if err != nil {
		return "(no audit available)"
	}
	return string(raw)
}

```

### With
```go
```

### Target File
.savepoint/Design.md

### Replace
```md
- **Doctor command** (`savepoint doctor`, `savepoint doctor --epic E##`) runs read-only integrity diagnostics for config, router state, release/epic/task structure, frontmatter validity, acceptance criteria presence, dependencies, duplicate task IDs, stale audit files, orphaned task IDs, and configured quality gates. It prints a human-readable report with repair suggestions and exits 0 when clean, 1 when problems are diagnosed, and 2 for internal or invocation failures.
```

### With
```md
- **Doctor command** (`savepoint doctor`, `savepoint doctor --epic E##`) runs read-only integrity diagnostics for config, router state, release/epic/task structure, frontmatter validity, acceptance criteria presence, dependencies, duplicate task IDs, stale audit files, orphaned task IDs, and configured quality gates. It prints a human-readable report with repair suggestions and exits 0 when clean, 1 when problems are diagnosed, and 2 for internal or invocation failures.
- **Audit remediation baseline** (v1.1 E13) centralizes frontmatter/body splitting and line-ending normalization in `internal/data`, uses typed sentinel errors for doctor repair suggestions, applies a configurable `quality_gates.gate_timeout`, removes tracked build artifacts from source control, adds `.golangci.yml`, and moves board filesystem writes/reads behind Bubble Tea command messages while preserving direct file I/O inside command helpers.
```

### Target File
.savepoint/Design.md

### Replace
```md
- **Board persistence and refresh:** task status transitions write canonical task frontmatter through `internal/data.WriteTaskStatus` with mtime conflict checks. The board treats `Model.Root` as the `.savepoint` directory, watches `.savepoint/releases/` recursively with fsnotify, adds watches for newly-created release/epic/task directories, and reloads task plus release/epic index data plus epic status metadata after debounced file changes. Router priority markers match release + epic + task, not only the short `T###` value; completed cards render with the orange build glyph even if they previously matched router priority. The `p` key explicitly writes the focused non-done task to router state as `task-building`; it does not infer `audit-pending` from task position. Epic status glyphs are cached from each epic's `E##-Detail.md` frontmatter and shown in the wide epic sidebar only.
```

### With
```md
- **Board persistence and refresh:** task status transitions write canonical task frontmatter through `internal/data.WriteTaskStatus` with mtime conflict checks. Board update handlers dispatch filesystem reads and writes through Bubble Tea command helpers (`routerWriteMsg`, `taskWriteMsg`, `epicDetailMsg`, `auditContentMsg`, and `errorMsg`) so `Update()` remains an event/message reducer. The board treats `Model.Root` as the `.savepoint` directory, watches `.savepoint/releases/` recursively with fsnotify, adds watches for newly-created release/epic/task directories, and reloads task plus release/epic index data plus epic status metadata after debounced file changes. Router priority markers match release + epic + task, not only the short `T###` value; completed cards render with the orange build glyph even if they previously matched router priority. The `p` key explicitly writes the focused non-done task to router state as `task-building`; it does not infer `audit-pending` from task position. Epic status glyphs are cached from each epic's `E##-Detail.md` frontmatter and shown in the wide epic sidebar only.
```

### Target File
AGENTS.md

### Replace
```md
| `internal/board/` | TUI board, overlays, epic sidebar, Next Activity line, router priority key, detail checklist rendering, status glyphs, forced color profile |
| `internal/buildtool/` | Makefile helper, cross-compile, archives |
| `internal/doctor/` | Read-only project diagnostics, integrity checks, quality gate execution, report formatting, repair suggestions |
| `internal/data/` | Task/router models, frontmatter parsing, lifecycle validation/defaulting, discovery |
```

### With
```md
| `internal/board/` | TUI board, overlays, epic sidebar, Next Activity line, router priority key, detail checklist rendering, status glyphs, forced color profile, async update I/O commands, shared board utilities |
| `internal/buildtool/` | Makefile helper, cross-compile, archives |
| `internal/doctor/` | Read-only project diagnostics, integrity checks, timed quality gate execution, report formatting, typed repair suggestions |
| `internal/data/` | Task/router models, frontmatter parsing/splitting, lifecycle validation/defaulting, discovery, canonical write helpers |
```

### Target File
.savepoint/releases/v1.1/epics/E13-audit-remediation/E13-Detail.md

### Replace
```md
## Boundaries
```

### With
```md
## Implemented As

- `internal/data` now owns shared line-ending normalization and frontmatter/body splitting for parser/write paths.
- `internal/doctor` now uses stack-based dependency cycle reconstruction, typed repair sentinels, and timed quality-gate command execution.
- `internal/init` uses a copy-based AtomicWrite fallback after failed rename attempts.
- `internal/board` now has shared utility helpers and dispatches board filesystem reads/writes through Bubble Tea command messages; `update.go` was decomposed but did not meet the original 30% line-count reduction target.
- `internal/buildtool` and `internal/styles` gained focused package tests.
- Repository hygiene now ignores build outputs and archives, removes tracked binaries, and includes `.golangci.yml`; lint execution still depends on `golangci-lint` being installed locally.

## Boundaries
```
