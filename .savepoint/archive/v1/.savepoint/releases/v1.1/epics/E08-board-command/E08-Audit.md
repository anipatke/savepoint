---
type: audit-findings
audited: 2026-05-03
---
# Audit Findings: E08 Board Command

## Main Findings

E08 is closed as audited. The board command surface, non-TTY fallback, Bubble Tea shell, model grouping, detail overlay, keyboard navigation, status writes, filters, theme profile detection, and integration coverage are in place.

The two audit blockers were applied: board task writes now surface mtime conflicts instead of overwriting externally modified task files, and dependency gating now blocks planned tasks from entering `in_progress` until dependencies are done.

Documentation reconciliation was applied to Design.md, AGENTS.md, and the E08 epic detail. Remaining process note: all eight E08 task files were missing the required `## Context Files` section; that was not backfilled because it is a task-authoring process gap rather than an implementation blocker.

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
- [x] Small diffs

Notes: The applied changes remove the speculative mtime overwrite and add branch coverage for planned-task dependency gating. The broader reducer size remains a known shape of `internal/board/update.go`, but the E08 audit blockers are resolved.

## Proposed Changes

### Target File
internal/board/update.go

### Replace
```go
func writeTaskStatus(task *data.Task, expectedMtime time.Time) error {
	if err := data.WriteTaskStatus(task.Path, task, expectedMtime); err != nil {
		if err != data.ErrMtimeConflict {
			return err
		}
		fi, statErr := os.Stat(task.Path)
		if statErr != nil {
			return statErr
		}
		if retryErr := data.WriteTaskStatus(task.Path, task, fi.ModTime()); retryErr != nil {
			return retryErr
		}
	}
	fi, err := os.Stat(task.Path)
	if err != nil {
		return err
	}
	task.Mtime = fi.ModTime()
	return nil
}
```

### With
```go
func writeTaskStatus(task *data.Task, expectedMtime time.Time) error {
	if err := data.WriteTaskStatus(task.Path, task, expectedMtime); err != nil {
		return err
	}
	fi, err := os.Stat(task.Path)
	if err != nil {
		return err
	}
	task.Mtime = fi.ModTime()
	return nil
}
```

### Target File
internal/board/integration_test.go

### Replace
```go
// TestMtimeConflict_boardRetries verifies the board retries after an mtime conflict and persists the transition.
func TestMtimeConflict_boardRetries(t *testing.T) {
	path := filepath.Join(t.TempDir(), "T001.md")
	content := "---\nid: E01/T001\nstatus: in_progress\nphase: build\n---\n\n# Task\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	task := data.Task{
		ID:     "E01/T001",
		Column: data.ColumnInProgress,
		Stage:  data.StageBuild,
		Path:   path,
		Mtime:  fi.ModTime().Add(-time.Minute), // intentionally stale
	}
	m := NewModel([]data.Task{task}, "v1", "E01")
	m.FocusedColumn = data.ColumnInProgress

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeySpace})
	updated := requireModel(t, got)

	if updated.AllTasks[0].Stage != data.StageTest {
		t.Errorf("Stage = %q after retry, want test", updated.AllTasks[0].Stage)
	}
	if strings.Contains(updated.StatusMessage, "conflict") {
		t.Errorf("conflict error surfaced to user unexpectedly: %q", updated.StatusMessage)
	}
}
```

### With
```go
// TestMtimeConflict_boardWarns verifies the board surfaces an mtime conflict instead of overwriting external edits.
func TestMtimeConflict_boardWarns(t *testing.T) {
	path := filepath.Join(t.TempDir(), "T001.md")
	content := "---\nid: E01/T001\nstatus: in_progress\nphase: build\n---\n\n# Task\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	task := data.Task{
		ID:     "E01/T001",
		Column: data.ColumnInProgress,
		Stage:  data.StageBuild,
		Path:   path,
		Mtime:  fi.ModTime().Add(-time.Minute), // intentionally stale
	}
	m := NewModel([]data.Task{task}, "v1", "E01")
	m.FocusedColumn = data.ColumnInProgress

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeySpace})
	updated := requireModel(t, got)

	if !strings.Contains(updated.StatusMessage, "mtime conflict") {
		t.Errorf("StatusMessage = %q, want mtime conflict warning", updated.StatusMessage)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), "phase: build") {
		t.Errorf("task file was overwritten despite mtime conflict:\n%s", raw)
	}
}
```

### Target File
internal/board/transitions.go

### Replace
```go
	case data.ColumnPlanned:
		return true, ""
	case data.ColumnInProgress:
```

### With
```go
	case data.ColumnPlanned:
		return dependenciesDone(t, allTasks)
	case data.ColumnInProgress:
```

### Target File
internal/board/transitions.go

### Replace
```go
		switch stage {
		case data.StageBuild:
			return true, ""
		case data.StageTest:
			return true, ""
		case data.StageAudit:
			for _, depID := range t.DependsOn {
				dep := findTask(depID, allTasks)
				if dep == nil {
					return false, fmt.Sprintf("dependency %q not found", depID)
				}
				if dep.Column != data.ColumnDone {
					return false, fmt.Sprintf("dependency %q is not done", depID)
				}
			}
			return true, ""
		default:
```

### With
```go
		switch stage {
		case data.StageBuild:
			return true, ""
		case data.StageTest:
			return true, ""
		case data.StageAudit:
			return dependenciesDone(t, allTasks)
		default:
```

### Target File
internal/board/transitions.go

### Replace
```go
func findTask(id string, tasks []data.Task) *data.Task {
	for i := range tasks {
		if tasks[i].ID == id {
			return &tasks[i]
		}
	}
	return nil
}
```

### With
```go
func dependenciesDone(t *data.Task, allTasks []data.Task) (bool, string) {
	for _, depID := range t.DependsOn {
		dep := findTask(depID, allTasks)
		if dep == nil {
			return false, fmt.Sprintf("dependency %q not found", depID)
		}
		if dep.Column != data.ColumnDone {
			return false, fmt.Sprintf("dependency %q is not done", depID)
		}
	}
	return true, ""
}

func findTask(id string, tasks []data.Task) *data.Task {
	for i := range tasks {
		if tasks[i].ID == id {
			return &tasks[i]
		}
	}
	return nil
}
```

### Target File
internal/board/transitions_test.go

### Replace
```go
func TestCanAdvance_plannedAlwaysAllowed(t *testing.T) {
	task := data.Task{Column: data.ColumnPlanned}
	ok, reason := CanAdvance(&task, nil)
	if !ok {
		t.Errorf("CanAdvance(planned) = false %q, want true", reason)
	}
}
```

### With
```go
func TestCanAdvance_plannedAllowedWhenDependenciesDone(t *testing.T) {
	allTasks := []data.Task{
		{ID: "T1", Column: data.ColumnPlanned, DependsOn: []string{"T2"}},
		{ID: "T2", Column: data.ColumnDone},
	}
	ok, reason := CanAdvance(&allTasks[0], allTasks)
	if !ok {
		t.Errorf("CanAdvance(planned with done dep) = false %q, want true", reason)
	}
}

func TestCanAdvance_plannedBlockedByDependency(t *testing.T) {
	allTasks := []data.Task{
		{ID: "T1", Column: data.ColumnPlanned, DependsOn: []string{"T2"}},
		{ID: "T2", Column: data.ColumnInProgress},
	}
	ok, reason := CanAdvance(&allTasks[0], allTasks)
	if ok {
		t.Fatal("CanAdvance(planned with unfinished dep) = true, want false")
	}
	if reason != "dependency \"T2\" is not done" {
		t.Errorf("reason = %q, want dependency warning", reason)
	}
}
```

### Target File
AGENTS.md

### Replace
```md
| `cmd/` | Init command arg parsing, dispatch |
```

### With
```md
| `cmd/` | CLI command arg parsing and dispatch for init and board |
```

### Target File
.savepoint/releases/v1.1/epics/E08-board-command/E08-Detail.md

### Replace
```md
## Boundaries
```

### With
```md
## Implemented As

- `main.go` dispatches `board` and uses the no-arg command path as the board default.
- `cmd/board.go` owns board-specific flag parsing for `--release`, `--epic`, and `--help`.
- `internal/board/board.go`, `tui.go`, `model.go`, `plain.go`, `transitions.go`, and related render/update files own project discovery, non-TTY fallback, TUI startup, filtering, detail overlays, keyboard flow, and task status writes.
- Audit entry-point work is represented as the non-TTY audit proposal signal and the epic Detail/Audit tab flow added around E06.

## Boundaries
```

### Target File
.savepoint/Design.md

### Replace
```md
- **Board command** (`savepoint board`) reads project state, renders the Atari-Noir TUI board, supports release/epic filtering, detail overlays, task status transitions with mtime-guarded writes, release/epic-scoped router priority markers, fsnotify-based task auto-refresh (epic E06), header Next Activity display, height-aware column/detail viewport scrolling, stable focused/unfocused column border geometry (v1.1 E01), dedicated phase-colored Next Activity line below the header, sentence-boundary checklist rendering in task details, shared status glyph mapping for task cards and the epic sidebar, a forced ANSI256 Lipgloss color profile for board startup (v1.1 E03), a focusable wide-screen epic sidebar with purple epic focus, epic detail overlays, and status glyphs loaded from epic detail frontmatter (v1.1 E04), and an epic Detail/Audit tab switch that renders user-facing audit findings from `{epic}/E##-Audit.md` (v1.1 E06).
```

### With
```md
- **Board command** (`savepoint board`, and bare `savepoint`) reads project state, renders the Atari-Noir TUI board when stdout is a TTY, falls back to a deterministic plain table in non-TTY mode, supports `--release`/`--epic` filtering, detail overlays, task status transitions with mtime-guarded writes, release/epic-scoped router priority markers, fsnotify-based task auto-refresh (epic E06), header Next Activity display, height-aware column/detail viewport scrolling, stable focused/unfocused column border geometry (v1.1 E01), dedicated phase-colored Next Activity line below the header, sentence-boundary checklist rendering in task details, shared status glyph mapping for task cards and the epic sidebar, a forced ANSI256 Lipgloss color profile for board startup (v1.1 E03), a focusable wide-screen epic sidebar with purple epic focus, epic detail overlays, and status glyphs loaded from epic detail frontmatter (v1.1 E04), and an epic Detail/Audit tab switch that renders user-facing audit findings from `{epic}/E##-Audit.md` (v1.1 E06).
```
