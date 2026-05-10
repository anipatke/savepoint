---
type: audit-findings
audited: 2026-05-10
---

# Audit Findings: E17 Defect Workflow TUI

## Main Findings

Audit proposals were applied on 2026-05-10. The implementation now enforces the defect rule that `stage` is required when `status: in_progress`, and the data/doctor tests were updated to cover the required-stage behavior. The live and scaffolded agent guidance now recognizes `defect-building` as repair flow, and `.savepoint/Design.md` plus `E17-Detail.md` were reconciled with the defect model, board overlay/detail paths, doctor diagnostics, and template responsibilities.

Remaining blocker before close:

- `T006-doctor-docs-and-templates.md` is still `status: planned`, while `.savepoint/router.md` says all E17 tasks are done and the epic is ready for audit. Only the user may mark a task `done`, so E17 has not been closed as audited and the router was not advanced.

Residual process gaps:

- T003, T004, T005, and T006 still have stale task implementation-plan checkboxes/context logs. The code and tests exist, but the task ledger does not fully document the build evidence.

Verified after apply:

- `make build` passed after rerunning outside the sandbox because the first attempt could not access the Go build cache under AppData.
- `make test` passed.

## Code Style Review

- [x] One job per file
- [x] One job per function
- [x] Test branches
- [x] Types document intent
- [x] Build only what is needed
- [x] Handle errors at boundaries
- [x] One source of truth
- [x] Comments explain WHY
- [x] Content in data files
- [ ] Small diffs - source changes are scoped, but task/router bookkeeping drift still prevents clean epic closeout

## Proposed Changes

### Target File
internal/data/lifecycle.go

### Replace
```go
func validateDefectLifecycle(d *Defect) error {
	if d.Status == "" {
		d.Status = ColumnPlanned
	}
	if !IsCanonicalColumn(d.Status) {
		return fmt.Errorf("invalid defect status %q: use planned, in_progress, or done", d.Status)
	}
	if d.Status == ColumnInProgress {
		if d.Stage == "" {
			d.Stage = StageBuild
			return nil
		}
		if !IsCanonicalStage(d.Stage) {
			return fmt.Errorf("invalid stage %q: use build, test, or audit", d.Stage)
		}
	}
	return nil
}
```

### With
```go
func validateDefectLifecycle(d *Defect) error {
	if d.Status == "" {
		d.Status = ColumnPlanned
	}
	if !IsCanonicalColumn(d.Status) {
		return fmt.Errorf("invalid defect status %q: use planned, in_progress, or done", d.Status)
	}
	if d.Status == ColumnInProgress {
		if d.Stage == "" {
			return fmt.Errorf("stage is required when defect status is in_progress")
		}
		if !IsCanonicalStage(d.Stage) {
			return fmt.Errorf("invalid stage %q: use build, test, or audit", d.Stage)
		}
		return nil
	}
	if d.Stage != "" {
		return fmt.Errorf("stage field %q is only valid when defect status is in_progress", d.Stage)
	}
	return nil
}
```

### Target File
internal/data/defect_test.go

### Replace
```go
func TestParseDefectFile_InProgressDefaultsToStageBuild(t *testing.T) {
	p := NewParser()
	content := `---
id: v1.1/D003-in-progress
status: in_progress
severity: medium
title: In progress defect
---
`
	defect, err := p.ParseDefectFile("D003.md", content)
	if err != nil {
		t.Fatalf("ParseDefectFile() error = %v", err)
	}
	if defect.Stage != StageBuild {
		t.Errorf("Stage = %q, want build (default)", defect.Stage)
	}
}
```

### With
```go
func TestParseDefectFile_InProgressRequiresStage(t *testing.T) {
	p := NewParser()
	content := `---
id: v1.1/D003-in-progress
status: in_progress
severity: medium
title: In progress defect
---
`
	_, err := p.ParseDefectFile("D003.md", content)
	if err == nil {
		t.Fatal("ParseDefectFile() error = nil, want error for missing in_progress stage")
	}
}
```

### Target File
internal/data/defect_test.go

### Replace
```go
func TestParseDefectFile_DoneStatus(t *testing.T) {
	p := NewParser()
	content := `---
id: v1.1/D005-done
status: done
severity: low
title: Fixed defect
---
`
	defect, err := p.ParseDefectFile("D005.md", content)
	if err != nil {
		t.Fatalf("ParseDefectFile() error = %v", err)
	}
	if defect.Status != ColumnDone {
		t.Errorf("Status = %q, want done", defect.Status)
	}
	if defect.Stage != "" {
		t.Errorf("Stage = %q, want empty for done defect", defect.Stage)
	}
}
```

### With
```go
func TestParseDefectFile_DoneStatus(t *testing.T) {
	p := NewParser()
	content := `---
id: v1.1/D005-done
status: done
severity: low
title: Fixed defect
---
`
	defect, err := p.ParseDefectFile("D005.md", content)
	if err != nil {
		t.Fatalf("ParseDefectFile() error = %v", err)
	}
	if defect.Status != ColumnDone {
		t.Errorf("Status = %q, want done", defect.Status)
	}
	if defect.Stage != "" {
		t.Errorf("Stage = %q, want empty for done defect", defect.Stage)
	}
}

func TestParseDefectFile_StageOnlyValidInProgress(t *testing.T) {
	p := NewParser()
	content := `---
id: v1.1/D005-done
status: done
stage: audit
severity: low
title: Fixed defect
---
`
	_, err := p.ParseDefectFile("D005.md", content)
	if err == nil {
		t.Fatal("ParseDefectFile() error = nil, want error for stage outside in_progress")
	}
}
```

### Target File
internal/doctor/checks_test.go

### Replace
```go
func TestCheckDefects_InProgressMissingStage(t *testing.T) {
	// ParseDefectFile defaults stage to build when in_progress — no error
	root := t.TempDir()
	defectsDir := filepath.Join(root, "releases", "v1", "defects")
	writeDefect(t, defectsDir, "D001-in-progress.md",
		"---\nid: v1/D001-in-progress\nstatus: in_progress\nseverity: medium\ntitle: In progress\n---\n")
	problems := CheckDefects(root)
	if len(problems) > 0 {
		t.Fatalf("CheckDefects() = %v, want no problems (in_progress defaults stage to build)", problems)
	}
}
```

### With
```go
func TestCheckDefects_InProgressMissingStage(t *testing.T) {
	root := t.TempDir()
	defectsDir := filepath.Join(root, "releases", "v1", "defects")
	writeDefect(t, defectsDir, "D001-in-progress.md",
		"---\nid: v1/D001-in-progress\nstatus: in_progress\nseverity: medium\ntitle: In progress\n---\n")
	problems := CheckDefects(root)
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, "stage is required") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckDefects() = %v, want missing stage problem", problems)
	}
}
```

### Target File
AGENTS.md

### Replace
```md
| audit-pending | savepoint-audit |
```

### With
```md
| audit-pending | savepoint-audit |
| defect-building | savepoint-build-task |
```

### Target File
AGENTS.md

### Replace
```md
## Implementation
```

### With
```md
## Defect Workflow

Use a defect conversation when the user reports a concrete bug, regression, broken behavior, or failed expectation that should be repaired without reshaping the planned epic/task backlog.

- Defects live at `.savepoint/releases/{release}/defects/D###-slug.md`.
- Router state may be `defect-building` with a `defect` field naming the active defect id.
- A defect in progress follows the same `status: in_progress` plus `stage: build|test|audit` rule as tasks.
- Use the board `d` overlay to inspect defects; do not turn defects into a fourth task column.

## Implementation
```

### Target File
AGENTS.md

### Replace
```md
| `internal/board/` | TUI board, overlays, epic sidebar, Next Activity line, router priority key, detail checklist rendering, status glyphs, forced color profile, debug logging hooks, async update I/O commands, shared board utilities |
```

### With
```md
| `internal/board/` | TUI board, overlays, epic sidebar, Next Activity line, router priority key, detail checklist rendering, status glyphs, forced color profile, debug logging hooks, async update I/O commands, defect summary/overlay/detail rendering, related-defect card markers, shared board utilities |
```

### Target File
AGENTS.md

### Replace
```md
| `internal/doctor/` | Read-only project diagnostics, integrity checks, timed quality gate execution, report formatting, typed repair suggestions |
```

### With
```md
| `internal/doctor/` | Read-only project diagnostics, integrity checks, defect validation, timed quality gate execution, report formatting, typed repair suggestions |
```

### Target File
AGENTS.md

### Replace
```md
| `internal/data/` | Task/router models, frontmatter parsing/splitting, lifecycle validation/defaulting, discovery including root-dir traversal, unified task status constants, canonical write helpers |
```

### With
```md
| `internal/data/` | Task/router/defect models, frontmatter parsing/splitting, lifecycle validation/defaulting, discovery including root-dir and release defect traversal, unified task status constants, canonical write helpers |
```

### Target File
AGENTS.md

### Replace
```md
| `templates/` | Scaffold markdown, YAML, prompts |
```

### With
```md
| `templates/` | Scaffold markdown, YAML, prompts, and defect workflow guidance |
```

### Target File
templates/project/AGENTS.md

### Replace
```md
| audit-pending | savepoint-audit |
```

### With
```md
| audit-pending | savepoint-audit |
| defect-building | savepoint-build-task |
```

### Target File
templates/project/AGENTS.md

### Replace
```md
- `savepoint doctor` validates defect files and reports malformed frontmatter, invalid status, and broken references.
```

### With
```md
- `savepoint doctor` validates defect files and reports malformed frontmatter, invalid status, missing in-progress stage, and broken references.
- If the router is in `defect-building`, treat the session as repair work for the named defect rather than normal epic planning or task-building.
```

### Target File
templates/project/.savepoint/router.md

### Replace
```md
### `state: task-building`

Task is `in_progress`. All `depends_on` are `done`.

**Next action:** Execute the plan. Tick checkboxes as you complete them. The implementation checklist exists to satisfy the task's acceptance criteria; every checked box should map to an observable outcome in `## Acceptance Criteria`. Edit code per the **Code Style** rules in `AGENTS.md`. Before setting `status: done`, update the task's `## Context Log`. When all checkboxes tick, run the full quality-gate suite, set `status: done`, update the router, and **stop for user review**.

**Do not start the next task without user acknowledgment.**
```

### With
```md
### `state: task-building`

Task is `in_progress`. All `depends_on` are `done`.

**Next action:** Execute the plan. Tick checkboxes as you complete them. The implementation checklist exists to satisfy the task's acceptance criteria; every checked box should map to an observable outcome in `## Acceptance Criteria`. Edit code per the **Code Style** rules in `AGENTS.md`. Before setting `status: done`, update the task's `## Context Log`. When all checkboxes tick, run the full quality-gate suite, set `status: done`, update the router, and **stop for user review**.

**Do not start the next task without user acknowledgment.**

### `state: defect-building`

A release-level defect is the active repair item. This is not epic planning and not normal task-building.

**Next action:** Read the active defect file from `.savepoint/releases/{release}/defects/{defect}.md`, then read only the source/test files needed to reproduce and fix it. Keep the defect file as the repair ledger, run the full quality-gate suite, update resolution notes, and stop for user review.

**Do not** convert a defect into an epic or task unless the user explicitly decides it is scope rather than repair.
```

### Target File
.savepoint/Design.md

### Replace
```md
- **Board command** (`savepoint board`, and bare `savepoint`) reads project state, renders the Atari-Noir TUI board when stdout is a TTY, falls back to a deterministic plain table in non-TTY mode, supports `--release`/`--epic` filtering, detail overlays, task status transitions with mtime-guarded writes, release/epic-scoped router priority markers, fsnotify-based task auto-refresh (epic E06), header Next Activity display, height-aware column/detail viewport scrolling, stable focused/unfocused column border geometry (v1.1 E01), dedicated phase-colored Next Activity line below the header, sentence-boundary checklist rendering in task details, shared status glyph mapping for task cards and the epic sidebar, a forced ANSI256 Lipgloss color profile for board startup (v1.1 E03), a focusable wide-screen epic sidebar with purple epic focus, epic detail overlays, and status glyphs loaded from epic detail frontmatter (v1.1 E04), and an epic Detail/Audit tab switch that renders user-facing audit findings from `{epic}/E##-Audit.md` (v1.1 E06).
```

### With
```md
- **Board command** (`savepoint board`, and bare `savepoint`) reads project state, renders the Atari-Noir TUI board when stdout is a TTY, falls back to a deterministic plain table in non-TTY mode, supports `--release`/`--epic` filtering, detail overlays, task status transitions with mtime-guarded writes, release/epic-scoped router priority markers, fsnotify-based task and defect auto-refresh, header Next Activity display, height-aware column/detail viewport scrolling, stable focused/unfocused column border geometry (v1.1 E01), dedicated phase-colored Next Activity line below the header including `DEFECT` router state, sentence-boundary checklist rendering in task details, shared status glyph mapping for task cards and the epic sidebar, a forced ANSI256 Lipgloss color profile for board startup (v1.1 E03), a focusable wide-screen epic sidebar with purple epic focus, epic detail overlays, status glyphs loaded from epic detail frontmatter (v1.1 E04), an epic Detail/Audit tab switch that renders user-facing audit findings from `{epic}/E##-Audit.md` (v1.1 E06), release-scoped open-defect counts, a keyboard-driven `d` defect overlay, defect detail overlays, and related-defect task card markers (v1.2 E17).
```

### Target File
.savepoint/Design.md

### Replace
```md
- **Doctor command** (`savepoint doctor`, `savepoint doctor --epic E##`) runs read-only integrity diagnostics for config, router state, release/epic/task structure, frontmatter validity, acceptance criteria presence, dependencies, duplicate task IDs, stale audit files, orphaned task IDs, and configured quality gates. It prints a human-readable report with repair suggestions and exits 0 when clean, 1 when problems are diagnosed, and 2 for internal or invocation failures.
```

### With
```md
- **Doctor command** (`savepoint doctor`, `savepoint doctor --epic E##`) runs read-only integrity diagnostics for config, router state, release/epic/task/defect structure, frontmatter validity, acceptance criteria presence, dependencies, duplicate task IDs, stale audit files, orphaned task IDs, broken defect references, and configured quality gates. It prints a human-readable report with repair suggestions and exits 0 when clean, 1 when problems are diagnosed, and 2 for internal or invocation failures.
```

### Target File
.savepoint/Design.md

### Replace
```md
            └── epics/
                └── E##-{epic-name}/
                    ├── E##-Detail.md   ← epic delta
                    ├── E##-Audit.md    ← audit findings + admin apply proposals
                    └── tasks/
                        └── T001-slug.md
```

### With
```md
            ├── defects/
            │   └── D001-slug.md    ← release-level repair record
            └── epics/
                └── E##-{epic-name}/
                    ├── E##-Detail.md   ← epic delta
                    ├── E##-Audit.md    ← audit findings + admin apply proposals
                    └── tasks/
                        └── T001-slug.md
```

### Target File
.savepoint/Design.md

### Replace
```md
| **Task**     | Independently buildable. Objective-led. **Requires implementation plan before build.** |
| **Sub-task** | Inline checklist item — _evidence of the implementation plan_, not standalone work.    |
```

### With
```md
| **Task**     | Independently buildable. Objective-led. **Requires implementation plan before build.** |
| **Defect**   | Release-level repair artifact for observed bugs or regressions; separate from planned epic/task scope. |
| **Sub-task** | Inline checklist item — _evidence of the implementation plan_, not standalone work.    |
```

### Target File
.savepoint/Design.md

### Replace
```md
## 5. Dependencies
```

### With
```md
Defects use the same `planned`, `in_progress`, and `done` status vocabulary. `stage` is required while a defect is `in_progress`. Router state may enter `defect-building` with a `defect` field naming the active repair item, which the board renders as a `DEFECT` Next Activity line.

## 5. Dependencies
```

### Target File
.savepoint/releases/v1.2/epics/E17-defect-workflow-tui/E17-Detail.md

### Replace
```md
## Boundaries
```

### With
```md
## Implemented As

- Defect model, parser, release-scoped discovery, and lifecycle validation live in `internal/data`.
- Board loading reads release defects alongside tasks and renders an open-defect header signal without adding a fourth column.
- The `d` key opens a defect overlay; enter opens a defect detail overlay that renders the evidence sections from the defect markdown body.
- Related task cards can show compact defect markers when a defect `reference` matches a visible task.
- Doctor validates defect files for parse errors, required fields, invalid status/stage, and broken task-like references.
- Scaffolded AGENTS/router templates and README now document the defect lane, but `defect-building` state guidance needs the audit proposal updates before close.

## Boundaries
```
