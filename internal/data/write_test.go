package data

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestWriteTaskStatus_updatesStatusAndPhase(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "task.md")
	content := `---
id: E01/T001
status: planned
phase: build
objective: "Test"
depends_on: []
---

# Body text`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	task := &Task{
		ID:     "E01/T001",
		Column: ColumnInProgress,
		Stage:  StageTest,
	}

	if err := WriteTaskStatus(path, task, fi.ModTime()); err != nil {
		t.Fatalf("WriteTaskStatus() error = %v", err)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	p := NewParser()
	parsed, err := p.ParseTaskFile(path, string(result))
	if err != nil {
		t.Fatalf("ParseTaskFile() error = %v", err)
	}

	if parsed.Column != ColumnInProgress {
		t.Errorf("Column = %v, want in_progress", parsed.Column)
	}
	if parsed.Stage != StageTest {
		t.Errorf("Stage = %v, want test", parsed.Stage)
	}

	if !strings.Contains(string(result), "# Body text") {
		t.Error("body content not preserved")
	}
}

func TestWriteTaskStatus_removesPhaseWhenStageEmpty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "task.md")
	content := `---
id: E01/T002
status: in_progress
phase: audit
objective: "Test"
---

# Body`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	fi, _ := os.Stat(path)

	task := &Task{
		ID:     "E01/T002",
		Column: ColumnDone,
		Stage:  "",
	}

	if err := WriteTaskStatus(path, task, fi.ModTime()); err != nil {
		t.Fatalf("WriteTaskStatus() error = %v", err)
	}

	result, _ := os.ReadFile(path)

	if strings.Contains(string(result), "phase:") {
		t.Error("phase field should be removed when stage is empty")
	}

	p := NewParser()
	parsed, err := p.ParseTaskFile(path, string(result))
	if err != nil {
		t.Fatalf("ParseTaskFile() error = %v", err)
	}

	if parsed.Column != ColumnDone {
		t.Errorf("Column = %v, want done", parsed.Column)
	}
	if parsed.Stage != "" {
		t.Errorf("Stage = %v, want empty", parsed.Stage)
	}
}

func TestWriteTaskStatus_removesPhaseWhenStatusPlanned(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "task.md")
	content := `---
id: E01/T003
status: in_progress
phase: build
---

# Body`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	fi, _ := os.Stat(path)

	task := &Task{
		ID:     "E01/T003",
		Column: ColumnPlanned,
		Stage:  "",
	}

	if err := WriteTaskStatus(path, task, fi.ModTime()); err != nil {
		t.Fatalf("WriteTaskStatus() error = %v", err)
	}

	result, _ := os.ReadFile(path)

	if strings.Contains(string(result), "phase:") {
		t.Error("phase field should be removed when status is planned")
	}
}

func TestWriteTaskStatus_mtimeConflict(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "task.md")
	content := `---
id: E01/T004
status: planned
---`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	oldMtime := time.Now().Add(-time.Hour)

	task := &Task{
		ID:     "E01/T004",
		Column: ColumnInProgress,
		Stage:  StageBuild,
	}

	err := WriteTaskStatus(path, task, oldMtime)
	if err == nil {
		t.Fatal("WriteTaskStatus() expected mtime conflict error")
	}
	if err != ErrMtimeConflict {
		t.Fatalf("WriteTaskStatus() error = %v, want ErrMtimeConflict", err)
	}
}

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

func TestWriteTaskStatus_removesLegacyStageField(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "task.md")
	content := `---
id: E01/T009
status: in_progress
stage: build
phase: build
objective: "Legacy mixed fields"
---`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	fi, _ := os.Stat(path)
	task := &Task{
		ID:     "E01/T009",
		Column: ColumnInProgress,
		Stage:  StageTest,
	}

	if err := WriteTaskStatus(path, task, fi.ModTime()); err != nil {
		t.Fatalf("WriteTaskStatus() error = %v", err)
	}

	result, _ := os.ReadFile(path)
	if strings.Contains(string(result), "stage:") {
		t.Error("legacy stage field should be removed")
	}
	if !strings.Contains(string(result), "phase: test") {
		t.Error("phase field should be updated to test")
	}
	parsed, err := NewParser().ParseTaskFile(path, string(result))
	if err != nil {
		t.Fatalf("ParseTaskFile() error = %v", err)
	}
	if parsed.Stage != StageTest {
		t.Errorf("Stage = %q, want test", parsed.Stage)
	}
}

func TestWriteTaskStatus_preservesBodyWithMultipleLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "task.md")
	content := `---
id: E01/T006
status: planned
---

# Title

Some description here.

More content.`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	fi, _ := os.Stat(path)

	task := &Task{
		ID:     "E01/T006",
		Column: ColumnInProgress,
		Stage:  StageBuild,
	}

	if err := WriteTaskStatus(path, task, fi.ModTime()); err != nil {
		t.Fatalf("WriteTaskStatus() error = %v", err)
	}

	result, _ := os.ReadFile(path)

	if !strings.Contains(string(result), "# Title") {
		t.Error("# Title not preserved")
	}
	if !strings.Contains(string(result), "Some description here.") {
		t.Error("description not preserved")
	}
	if !strings.Contains(string(result), "More content.") {
		t.Error("More content not preserved")
	}
}

func TestWriteRouterState_updatesRouterFields(t *testing.T) {
	dir := t.TempDir()
	root := dir
	content := `# Agent State Machine

## Current state

` + "```" + `yaml
state: task-building
release: v1
epic: E03-board-tui-core
task: E03-board-tui-core/T004-render
next_action: "Render the board"
` + "```" + `

## State definitions`

	path := filepath.Join(root, "router.md")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	state := &RouterState{
		State:      "task-building",
		Release:    "v1",
		Epic:       "E05-phase-transitions",
		Task:       "E05-phase-transitions/T004-write-router",
		NextAction: "Write router state",
	}

	if err := WriteRouterState(root, state, fi.ModTime()); err != nil {
		t.Fatalf("WriteRouterState() error = %v", err)
	}

	r := NewRouterReader()
	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	parsed, err := r.ReadState(string(result))
	if err != nil {
		t.Fatalf("ReadState() error = %v", err)
	}

	if parsed.State != "task-building" {
		t.Errorf("State = %q, want task-building", parsed.State)
	}
	if parsed.Epic != "E05-phase-transitions" {
		t.Errorf("Epic = %q, want E05-phase-transitions", parsed.Epic)
	}
	if parsed.Release != "v1" {
		t.Errorf("Release = %q, want v1", parsed.Release)
	}
	if parsed.Task != "E05-phase-transitions/T004-write-router" {
		t.Errorf("Task = %q, want E05-phase-transitions/T004-write-router", parsed.Task)
	}
	if parsed.NextAction != "Write router state" {
		t.Errorf("NextAction = %q, want Write router state", parsed.NextAction)
	}

	if !strings.Contains(string(result), "## State definitions") {
		t.Error("body content after state block not preserved")
	}
}

func TestWriteRouterState_mtimeConflict(t *testing.T) {
	dir := t.TempDir()
	root := dir
	content := `## Current state

` + "```" + `yaml
state: task-building
` + "```" + `
`

	path := filepath.Join(root, "router.md")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	oldMtime := time.Now().Add(-time.Hour)
	state := &RouterState{State: "audit-pending"}

	err := WriteRouterState(root, state, oldMtime)
	if err == nil {
		t.Fatal("WriteRouterState() expected mtime conflict error")
	}
	if err != ErrMtimeConflict {
		t.Fatalf("WriteRouterState() error = %v, want ErrMtimeConflict", err)
	}
}

func TestWriteRouterState_missingStateBlock(t *testing.T) {
	dir := t.TempDir()
	root := dir
	content := `# No state block`

	path := filepath.Join(root, "router.md")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	fi, _ := os.Stat(path)
	state := &RouterState{State: "task-building"}

	err := WriteRouterState(root, state, fi.ModTime())
	if err == nil {
		t.Fatal("WriteRouterState() expected error for missing state block")
	}
}

func TestWriteRouterState_preservesNextAction(t *testing.T) {
	dir := t.TempDir()
	root := dir
	content := `## Current state

` + "```" + `yaml
state: task-building
release: v1
epic: E03-board-tui-core
task: ""
next_action: "Do the thing"
` + "```" + `
`

	path := filepath.Join(root, "router.md")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	fi, _ := os.Stat(path)

	state := &RouterState{
		State:      "task-building",
		Release:    "v1",
		Epic:       "E05-phase-transitions",
		Task:       "",
		NextAction: "Do the thing",
	}

	if err := WriteRouterState(root, state, fi.ModTime()); err != nil {
		t.Fatalf("WriteRouterState() error = %v", err)
	}

	r := NewRouterReader()
	result, _ := os.ReadFile(path)
	parsed, _ := r.ReadState(string(result))

	if parsed.NextAction != "Do the thing" {
		t.Errorf("NextAction = %q, want %q", parsed.NextAction, "Do the thing")
	}
}

func TestApplyProposal_replacesText(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "Design.md")
	content := "# Architecture\n\nOld section text.\n\nMore content."
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	if err := ApplyProposal(path, "Old section text.", "New section text."); err != nil {
		t.Fatalf("ApplyProposal() error = %v", err)
	}

	result, _ := os.ReadFile(path)
	if !strings.Contains(string(result), "New section text.") {
		t.Error("replacement not applied")
	}
	if strings.Contains(string(result), "Old section text.") {
		t.Error("old text still present")
	}
	if !strings.Contains(string(result), "More content.") {
		t.Error("surrounding content not preserved")
	}
}

func TestApplyProposal_missingTarget(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "Design.md")
	if err := os.WriteFile(path, []byte("some content"), 0644); err != nil {
		t.Fatal(err)
	}

	err := ApplyProposal(path, "not present", "replacement")
	if err == nil {
		t.Fatal("ApplyProposal() expected error for missing target")
	}
}

func TestUpdateEpicStatus_setsStatusField(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "E06-Detail.md")
	content := "---\ntype: epic-design\nstatus: planned\n---\n\n# E06 Body"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	if err := UpdateEpicStatus(path, "audited"); err != nil {
		t.Fatalf("UpdateEpicStatus() error = %v", err)
	}

	result, _ := os.ReadFile(path)
	if !strings.Contains(string(result), "status: audited") {
		t.Error("status not updated to audited")
	}
	if !strings.Contains(string(result), "# E06 Body") {
		t.Error("body not preserved")
	}
}

func TestUpdateLastAudited_setsField(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "Design.md")
	content := "---\ntype: project-design\nstatus: active\nlast_audited: v1.1/E05-tasking-permissions\n---\n\n# Body"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	if err := UpdateLastAudited(path, "v1.1/E06-audit-command"); err != nil {
		t.Fatalf("UpdateLastAudited() error = %v", err)
	}

	result, _ := os.ReadFile(path)
	if !strings.Contains(string(result), "last_audited: v1.1/E06-audit-command") {
		t.Error("last_audited not updated")
	}
}

func TestUpdateLastAudited_addsFieldIfMissing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "Design.md")
	content := "---\ntype: project-design\nstatus: active\n---\n\n# Body"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	if err := UpdateLastAudited(path, "v1.1/E06-audit-command"); err != nil {
		t.Fatalf("UpdateLastAudited() error = %v", err)
	}

	result, _ := os.ReadFile(path)
	if !strings.Contains(string(result), "last_audited: v1.1/E06-audit-command") {
		t.Error("last_audited not added")
	}
}

func TestWriteTaskStatus_noFrontmatter(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "task.md")
	content := `# No frontmatter here`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	fi, _ := os.Stat(path)

	task := &Task{
		ID:     "E01/T007",
		Column: ColumnPlanned,
	}

	err := WriteTaskStatus(path, task, fi.ModTime())
	if err == nil {
		t.Fatal("WriteTaskStatus() expected error for missing frontmatter")
	}
}

func TestWriteTaskStatus_rejectsInvalidLifecycle(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "task.md")
	content := `---
id: E01/T008
status: planned
---`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	fi, _ := os.Stat(path)
	task := &Task{
		ID:     "E01/T008",
		Column: ColumnDone,
		Stage:  StageAudit,
	}

	err := WriteTaskStatus(path, task, fi.ModTime())
	if err == nil {
		t.Fatal("WriteTaskStatus() expected invalid lifecycle error")
	}
}
