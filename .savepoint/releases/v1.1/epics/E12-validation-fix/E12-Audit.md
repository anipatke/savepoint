---
type: audit-findings
audited: 2026-05-03
---
# Audit Findings: E12 Task File Validation & Auto-Fix

## Main Findings

E12 is closed. Parser-side defaults are in place: missing task status parses as `planned`, missing phase/stage on parsed `in_progress` tasks becomes `build`, and invalid status/phase messages include actionable hints.

The write-time lifecycle gap found during audit was fixed. `ValidateTaskLifecycle` now receives a task pointer, so defaulting an `in_progress` task to `StageBuild` persists for callers such as `WriteTaskStatus`. Regression coverage now verifies both lifecycle defaulting and persisted `phase: build` writes.

The E12 task files now include the required `## Context Files` sections. Architecture and agent documentation were reconciled to describe `internal/data` lifecycle validation/defaulting.

Verification after applying proposals: `go build ./...` passed and `go test ./...` passed.

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

## Proposed Changes

### Target File
internal/data/lifecycle.go

### Replace
```go
func ValidateTaskLifecycle(task Task) error {
	if !IsCanonicalColumn(task.Column) {
		return fmt.Errorf("invalid status %q: use planned, in_progress, or done. Add 'status: planned' or 'status: in_progress' to task frontmatter", task.Column)
	}

	if task.Column == ColumnInProgress {
		if task.Stage == "" {
			task.Stage = StageBuild
			return nil
		}
		if !IsCanonicalStage(task.Stage) {
			return fmt.Errorf("invalid phase %q: use build, test, or audit. Add 'phase: build' to task frontmatter", task.Stage)
		}
		return nil
	}

	if task.Stage != "" {
		return fmt.Errorf("phase field %q is only valid when status is in_progress. Remove 'phase' or change status to in_progress", task.Stage)
	}

	return nil
}
```

### With
```go
func ValidateTaskLifecycle(task *Task) error {
	if !IsCanonicalColumn(task.Column) {
		return fmt.Errorf("invalid status %q: use planned, in_progress, or done. Add 'status: planned' or 'status: in_progress' to task frontmatter", task.Column)
	}

	if task.Column == ColumnInProgress {
		if task.Stage == "" {
			task.Stage = StageBuild
			return nil
		}
		if !IsCanonicalStage(task.Stage) {
			return fmt.Errorf("invalid phase %q: use build, test, or audit. Add 'phase: build' to task frontmatter", task.Stage)
		}
		return nil
	}

	if task.Stage != "" {
		return fmt.Errorf("phase field %q is only valid when status is in_progress. Remove 'phase' or change status to in_progress", task.Stage)
	}

	return nil
}
```

### Target File
internal/data/write.go

### Replace
```go
	if err := ValidateTaskLifecycle(*task); err != nil {
		return err
	}
```

### With
```go
	if err := ValidateTaskLifecycle(task); err != nil {
		return err
	}
```

### Target File
internal/data/lifecycle_test.go

### Replace
```go
func TestValidateTaskLifecycle_allowsPlannedWithoutPhase(t *testing.T) {
	task := Task{Column: ColumnPlanned}
	if err := ValidateTaskLifecycle(task); err != nil {
		t.Fatalf("ValidateTaskLifecycle() error = %v", err)
	}
}

func TestValidateTaskLifecycle_allowsInProgressWithPhase(t *testing.T) {
	task := Task{Column: ColumnInProgress, Stage: StageAudit}
	if err := ValidateTaskLifecycle(task); err != nil {
		t.Fatalf("ValidateTaskLifecycle() error = %v", err)
	}
}

func TestValidateTaskLifecycle_rejectsUnknownStatus(t *testing.T) {
	task := Task{Column: "review"}
	if err := ValidateTaskLifecycle(task); err == nil {
		t.Fatal("ValidateTaskLifecycle() expected unknown status error")
	}
}

func TestValidateTaskLifecycle_rejectsPhaseOutsideInProgress(t *testing.T) {
	task := Task{Column: ColumnPlanned, Stage: StageBuild}
	if err := ValidateTaskLifecycle(task); err == nil {
		t.Fatal("ValidateTaskLifecycle() expected phase/status error")
	}
}
```

### With
```go
func TestValidateTaskLifecycle_allowsPlannedWithoutPhase(t *testing.T) {
	task := Task{Column: ColumnPlanned}
	if err := ValidateTaskLifecycle(&task); err != nil {
		t.Fatalf("ValidateTaskLifecycle() error = %v", err)
	}
}

func TestValidateTaskLifecycle_defaultsInProgressWithoutPhase(t *testing.T) {
	task := Task{Column: ColumnInProgress}
	if err := ValidateTaskLifecycle(&task); err != nil {
		t.Fatalf("ValidateTaskLifecycle() error = %v", err)
	}
	if task.Stage != StageBuild {
		t.Fatalf("Task.Stage = %q, want %q", task.Stage, StageBuild)
	}
}

func TestValidateTaskLifecycle_allowsInProgressWithPhase(t *testing.T) {
	task := Task{Column: ColumnInProgress, Stage: StageAudit}
	if err := ValidateTaskLifecycle(&task); err != nil {
		t.Fatalf("ValidateTaskLifecycle() error = %v", err)
	}
}

func TestValidateTaskLifecycle_rejectsUnknownStatus(t *testing.T) {
	task := Task{Column: "review"}
	if err := ValidateTaskLifecycle(&task); err == nil {
		t.Fatal("ValidateTaskLifecycle() expected unknown status error")
	}
}

func TestValidateTaskLifecycle_rejectsPhaseOutsideInProgress(t *testing.T) {
	task := Task{Column: ColumnPlanned, Stage: StageBuild}
	if err := ValidateTaskLifecycle(&task); err == nil {
		t.Fatal("ValidateTaskLifecycle() expected phase/status error")
	}
}
```

### Target File
internal/data/write_test.go

### Replace
```go
func TestWriteTaskStatus_addsPhaseWhenStagePresent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "task.md")
	content := `---
id: E01/T005
status: in_progress
objective: "No phase yet"
---`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	fi, _ := os.Stat(path)

	task := &Task{
		ID:     "E01/T005",
		Column: ColumnInProgress,
		Stage:  StageAudit,
	}

	if err := WriteTaskStatus(path, task, fi.ModTime()); err != nil {
		t.Fatalf("WriteTaskStatus() error = %v", err)
	}

	result, _ := os.ReadFile(path)

	if !strings.Contains(string(result), "phase: audit") {
		t.Error("phase field should be added when stage is set")
	}
}
```

### With
```go
func TestWriteTaskStatus_addsPhaseWhenStagePresent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "task.md")
	content := `---
id: E01/T005
status: in_progress
objective: "No phase yet"
---`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	fi, _ := os.Stat(path)

	task := &Task{
		ID:     "E01/T005",
		Column: ColumnInProgress,
		Stage:  StageAudit,
	}

	if err := WriteTaskStatus(path, task, fi.ModTime()); err != nil {
		t.Fatalf("WriteTaskStatus() error = %v", err)
	}

	result, _ := os.ReadFile(path)

	if !strings.Contains(string(result), "phase: audit") {
		t.Error("phase field should be added when stage is set")
	}
}

func TestWriteTaskStatus_defaultsInProgressPhaseWhenStageMissing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "task.md")
	content := `---
id: E01/T010
status: planned
objective: "No phase yet"
---`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	fi, _ := os.Stat(path)

	task := &Task{
		ID:     "E01/T010",
		Column: ColumnInProgress,
	}

	if err := WriteTaskStatus(path, task, fi.ModTime()); err != nil {
		t.Fatalf("WriteTaskStatus() error = %v", err)
	}

	result, _ := os.ReadFile(path)

	if !strings.Contains(string(result), "phase: build") {
		t.Error("phase field should default to build for in_progress writes")
	}
}
```

### Target File
.savepoint/Design.md

### Replace
```md
- **Go data-reader boundary:** established in epic `E02-data-readers` (2026-05-01). `internal/data` owns Savepoint file parsing and discovery for the Go implementation: task frontmatter models, markdown YAML extraction, router state parsing, config theme defaults, release/epic/task directory listing, and boundary error sentinels.
```

### With
```md
- **Go data-reader boundary:** established in epic `E02-data-readers` (2026-05-01). `internal/data` owns Savepoint file parsing and discovery for the Go implementation: task frontmatter models, markdown YAML extraction, router state parsing, config theme defaults, release/epic/task directory listing, task lifecycle validation/defaulting, write-time status validation, and boundary error sentinels.
```

### Target File
AGENTS.md

### Replace
```md
| `internal/data/` | Task/router models, frontmatter parsing, discovery |
```

### With
```md
| `internal/data/` | Task/router models, frontmatter parsing, lifecycle validation/defaulting, discovery |
```

### Target File
.savepoint/releases/v1.1/epics/E12-validation-fix/E12-Detail.md

### Replace
```md
**Out of scope:**
- New UI features
- New commands
```

### With
```md
**Out of scope:**
- New UI features
- New commands

## Implemented as

- Parser-side defaults live in `internal/data/parser.go`: empty status/column normalizes to `planned`, and parsed `in_progress` tasks without phase/stage default to `build`.
- Lifecycle validation and user-facing hints live in `internal/data/lifecycle.go`.
- Status writes validate lifecycle rules through `internal/data/write.go`; audit follow-up must persist the default `phase: build` when callers write `in_progress` with no stage.
- Regression coverage lives in `internal/data/parser_test.go`, `internal/data/lifecycle_test.go`, and `internal/data/write_test.go`.
```

### Target File
.savepoint/releases/v1.1/epics/E12-validation-fix/tasks/T001-default-phase.md

### Replace
```md
# T001: Default Phase to Build

## Acceptance Criteria
```

### With
```md
# T001: Default Phase to Build

## Context Files

- `internal/data/parser.go`
- `internal/data/lifecycle.go`
- `internal/data/parser_test.go`

## Acceptance Criteria
```

### Target File
.savepoint/releases/v1.1/epics/E12-validation-fix/tasks/T002-default-status.md

### Replace
```md
# T002: Default Status to Planned

## Acceptance Criteria
```

### With
```md
# T002: Default Status to Planned

## Context Files

- `internal/data/parser.go`
- `internal/data/parser_test.go`

## Acceptance Criteria
```

### Target File
.savepoint/releases/v1.1/epics/E12-validation-fix/tasks/T003-better-errors.md

### Replace
```md
# T003: Better Error Messages

## Acceptance Criteria
```

### With
```md
# T003: Better Error Messages

## Context Files

- `internal/data/lifecycle.go`
- `internal/data/lifecycle_test.go`

## Acceptance Criteria
```

### Target File
.savepoint/releases/v1.1/epics/E12-validation-fix/tasks/T004-validate-on-write.md

### Replace
```md
# T004: Validate On Write

## Acceptance Criteria
```

### With
```md
# T004: Validate On Write

## Context Files

- `internal/data/write.go`
- `internal/data/write_test.go`
- `internal/data/lifecycle.go`

## Acceptance Criteria
```

### Target File
.savepoint/releases/v1.1/epics/E12-validation-fix/tasks/T005-tests.md

### Replace
```md
# T005: Tests

## Acceptance Criteria
```

### With
```md
# T005: Tests

## Context Files

- `internal/data/parser_test.go`
- `internal/data/lifecycle_test.go`
- `internal/data/write_test.go`

## Acceptance Criteria
```
