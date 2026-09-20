---
type: audit-findings
audited: 2026-05-03
---
# Audit Findings: E09 Doctor Command

## Main Findings

Applied audit proposals for E09. The `doctor` structure check now validates the active `{release}-PRD.md` layout instead of the obsolete `PRD.md` filename, and the related tests and repair suggestion were updated. During verification, the repair switch was also tightened so the specific release PRD hint is not shadowed by the generic release-directory hint.

Design.md, AGENTS.md, and E09-Detail.md now document the implemented doctor command and `internal/doctor/` module. E09 is marked audited and Design.md `last_audited` points to `v1.1/E09-doctor-command`.

Residual process note: the E09 task files still lack the required `## Context Files` sections. That does not block the applied code fix, but it should be corrected in future task authoring.

Verification after apply: `go build ./...` and `go test ./...` passed. `make build && make test` remains unavailable in this environment because `make` is not installed.

## Code Style Review

- [x] One job per file
- [x] One-sentence functions
- [x] Test branches
- [x] Types are documentation
- [x] Build, don't speculate
- [x] Errors at boundaries
- [x] One source of truth
- [x] Comments explain WHY
- [x] Content in data files
- [x] Small diffs - the implementation fix was narrowly scoped; the remaining context-section note is process cleanup, not product code risk.

## Proposed Changes

### Target File
internal/doctor/checks.go

### Replace
```go
		for _, release := range releases {
		checkReleasePRD(release.Path, &problems)
```

### With
```go
	for _, release := range releases {
		checkReleasePRD(release.Path, release.ID, &problems)
```

### Target File
internal/doctor/checks.go

### Replace
```go
func checkReleasePRD(releasePath string, problems *[]Problem) {
	prdPath := filepath.Join(releasePath, "PRD.md")
```

### With
```go
func checkReleasePRD(releasePath string, releaseID string, problems *[]Problem) {
	prdPath := filepath.Join(releasePath, releaseID+"-PRD.md")
```

### Target File
internal/doctor/checks.go

### Replace
```go
		*problems = append(*problems, Problem{File: prdPath, Message: "release PRD.md not found"})
```

### With
```go
		*problems = append(*problems, Problem{File: prdPath, Message: "release PRD file not found"})
```

### Target File
internal/doctor/repairs.go

### Replace
```go
	case contains(p.Message, "PRD.md not found"):
		return "Create a PRD.md with frontmatter for the release"
```

### With
```go
	case contains(p.Message, "release PRD file not found"):
		return "Create a {release}-PRD.md file with frontmatter for the release"
```

### Target File
internal/doctor/checks_test.go

### Replace
```go
	writeFile(t, filepath.Join(releasePath, "PRD.md"), "---\ntype: project-prd\nstatus: active\n---\n\n# Release\n")
```

### With
```go
	writeFile(t, filepath.Join(releasePath, filepath.Base(releasePath)+"-PRD.md"), "---\ntype: project-prd\nstatus: active\n---\n\n# Release\n")
```

### Target File
internal/doctor/report_test.go

### Replace
```go
	writeFile(t, filepath.Join(releasePath, "PRD.md"), "---\ntype: project-prd\nstatus: active\n---\n\n# Release\n")
```

### With
```go
	writeFile(t, filepath.Join(releasePath, "v1-PRD.md"), "---\ntype: project-prd\nstatus: active\n---\n\n# Release\n")
```

### Target File
.savepoint/Design.md

### Replace
```md
- **Board command** (`savepoint board`, and bare `savepoint`) reads project state, renders the Atari-Noir TUI board when stdout is a TTY, falls back to a deterministic plain table in non-TTY mode, supports `--release`/`--epic` filtering, detail overlays, task status transitions with mtime-guarded writes, release/epic-scoped router priority markers, fsnotify-based task auto-refresh (epic E06), header Next Activity display, height-aware column/detail viewport scrolling, stable focused/unfocused column border geometry (v1.1 E01), dedicated phase-colored Next Activity line below the header, sentence-boundary checklist rendering in task details, shared status glyph mapping for task cards and the epic sidebar, a forced ANSI256 Lipgloss color profile for board startup (v1.1 E03), a focusable wide-screen epic sidebar with purple epic focus, epic detail overlays, and status glyphs loaded from epic detail frontmatter (v1.1 E04), and an epic Detail/Audit tab switch that renders user-facing audit findings from `{epic}/E##-Audit.md` (v1.1 E06).
```

### With
```md
- **Board command** (`savepoint board`, and bare `savepoint`) reads project state, renders the Atari-Noir TUI board when stdout is a TTY, falls back to a deterministic plain table in non-TTY mode, supports `--release`/`--epic` filtering, detail overlays, task status transitions with mtime-guarded writes, release/epic-scoped router priority markers, fsnotify-based task auto-refresh (epic E06), header Next Activity display, height-aware column/detail viewport scrolling, stable focused/unfocused column border geometry (v1.1 E01), dedicated phase-colored Next Activity line below the header, sentence-boundary checklist rendering in task details, shared status glyph mapping for task cards and the epic sidebar, a forced ANSI256 Lipgloss color profile for board startup (v1.1 E03), a focusable wide-screen epic sidebar with purple epic focus, epic detail overlays, and status glyphs loaded from epic detail frontmatter (v1.1 E04), and an epic Detail/Audit tab switch that renders user-facing audit findings from `{epic}/E##-Audit.md` (v1.1 E06).
- **Doctor command** (`savepoint doctor`, `savepoint doctor --epic E##`) runs read-only integrity diagnostics for config, router state, release/epic/task structure, frontmatter validity, acceptance criteria presence, dependencies, duplicate task IDs, stale audit files, orphaned task IDs, and configured quality gates. It prints a human-readable report with repair suggestions and exits 0 when clean, 1 when problems are diagnosed, and 2 for internal or invocation failures.
```

### Target File
.savepoint/Design.md

### Replace
```md
| `cmd/` | CLI command arg parsing and dispatch for init and board |
```

### With
```md
| `cmd/` | CLI command arg parsing and dispatch for init, board, and doctor |
```

### Target File
.savepoint/Design.md

### Replace
```md
| `internal/board/` | TUI board, overlays, epic sidebar, Next Activity line, router priority key, detail checklist rendering, status glyphs, forced color profile |
| `internal/buildtool/` | Makefile helper, cross-compile, archives |
```

### With
```md
| `internal/board/` | TUI board, overlays, epic sidebar, Next Activity line, router priority key, detail checklist rendering, status glyphs, forced color profile |
| `internal/buildtool/` | Makefile helper, cross-compile, archives |
| `internal/doctor/` | Read-only project diagnostics, integrity checks, quality gate execution, report formatting, repair suggestions |
```

### Target File
AGENTS.md

### Replace
```md
| `cmd/` | CLI command arg parsing and dispatch for init and board |
```

### With
```md
| `cmd/` | CLI command arg parsing and dispatch for init, board, and doctor |
```

### Target File
AGENTS.md

### Replace
```md
| `internal/board/` | TUI board, overlays, epic sidebar, Next Activity line, router priority key, detail checklist rendering, status glyphs, forced color profile |
| `internal/buildtool/` | Makefile helper, cross-compile, archives |
```

### With
```md
| `internal/board/` | TUI board, overlays, epic sidebar, Next Activity line, router priority key, detail checklist rendering, status glyphs, forced color profile |
| `internal/buildtool/` | Makefile helper, cross-compile, archives |
| `internal/doctor/` | Read-only project diagnostics, integrity checks, quality gate execution, report formatting, repair suggestions |
```

### Target File
.savepoint/releases/v1.1/epics/E09-doctor-command/E09-Detail.md

### Replace
```md
## Boundaries
```

### With
```md
## Implemented As

- `cmd/doctor.go` parses `doctor [--epic <epic>]`, reports help, rejects unsupported arguments, and delegates execution through an injected runner.
- `main.go` wires `savepoint doctor` to `internal/doctor.RunAllChecks` and preserves the required 0/1/2 exit-code contract.
- `internal/doctor/checks.go` implements config, router, structure, dependency, duplicate ID, audit-state, and orphan diagnostics.
- `internal/doctor/gates.go`, `report.go`, and `repairs.go` run configured quality gates, format human-readable reports, and attach repair suggestions.
- Tests live in `cmd/doctor_test.go` and `internal/doctor/*_test.go`.

## Boundaries
```
