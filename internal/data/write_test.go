package data

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestWriteTaskStatus_updatesStatusAndStage(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "task.md")
	content := `---
id: E01/T-001
status: planned
stage: build
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
		ID:     "E01/T-001",
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

func TestWriteTaskStatus_removesStageWhenStageEmpty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "task.md")
	content := `---
id: E01/T-002
status: in_progress
stage: audit
objective: "Test"
---

# Body`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	fi, _ := os.Stat(path)

	task := &Task{
		ID:     "E01/T-002",
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
	if strings.Contains(string(result), "stage:") {
		t.Error("stage field should be removed when stage is empty")
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

func TestWriteTaskStatus_rewritesAgentCompleteStatusAsCanonicalDone(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "task.md")
	content := `---
id: E01/T-002
status: complete
objective: "Agent completed task"
---

# Body`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	parsed, err := NewParser().ParseTaskFile(path, content)
	if err != nil {
		t.Fatalf("ParseTaskFile() error = %v", err)
	}
	if parsed.Column != ColumnDone {
		t.Fatalf("parsed Column = %q, want done", parsed.Column)
	}

	fi, _ := os.Stat(path)
	if err := WriteTaskStatus(path, parsed, fi.ModTime()); err != nil {
		t.Fatalf("WriteTaskStatus() error = %v", err)
	}

	result, _ := os.ReadFile(path)
	if strings.Contains(string(result), "status: complete") {
		t.Error("agent status alias should not be preserved")
	}
	if !strings.Contains(string(result), "status: done") {
		t.Error("status should be written as canonical done")
	}
}

func TestWriteTaskStatus_removesProgressFieldsWhenStatusPlanned(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "task.md")
	content := `---
id: E01/T-003
status: in_progress
stage: build
---

# Body`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	fi, _ := os.Stat(path)

	task := &Task{
		ID:     "E01/T-003",
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
	if strings.Contains(string(result), "stage:") {
		t.Error("stage field should be removed when status is planned")
	}
}

func TestWriteTaskStatus_mtimeConflict(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "task.md")
	content := `---
id: E01/T-004
status: planned
---`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	oldMtime := time.Now().Add(-time.Hour)

	task := &Task{
		ID:     "E01/T-004",
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

func TestWriteTaskStatus_addsStageWhenStagePresent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "task.md")
	content := `---
id: E01/T-005
status: in_progress
objective: "No phase yet"
---`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	fi, _ := os.Stat(path)

	task := &Task{
		ID:     "E01/T-005",
		Column: ColumnInProgress,
		Stage:  StageAudit,
	}

	if err := WriteTaskStatus(path, task, fi.ModTime()); err != nil {
		t.Fatalf("WriteTaskStatus() error = %v", err)
	}

	result, _ := os.ReadFile(path)

	if !strings.Contains(string(result), "stage: audit") {
		t.Error("stage field should be added when stage is set")
	}
	if strings.Contains(string(result), "phase:") {
		t.Error("legacy phase field should not be written")
	}
}

func TestWriteTaskStatus_rejectsInProgressWhenStageMissing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "task.md")
	content := `---
id: E01/T-010
status: planned
objective: "No phase yet"
---`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	fi, _ := os.Stat(path)

	task := &Task{
		ID:     "E01/T-010",
		Column: ColumnInProgress,
	}

	err := WriteTaskStatus(path, task, fi.ModTime())
	if err == nil {
		t.Fatal("WriteTaskStatus() expected missing stage error")
	}
	if !strings.Contains(err.Error(), "stage is required") {
		t.Fatalf("WriteTaskStatus() error = %v, want missing stage message", err)
	}
}

func TestWriteTaskStatus_removesLegacyPhaseField(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "task.md")
	content := `---
id: E01/T-009
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
		ID:     "E01/T-009",
		Column: ColumnInProgress,
		Stage:  StageTest,
	}

	if err := WriteTaskStatus(path, task, fi.ModTime()); err != nil {
		t.Fatalf("WriteTaskStatus() error = %v", err)
	}

	result, _ := os.ReadFile(path)
	if strings.Contains(string(result), "phase:") {
		t.Error("legacy phase field should be removed")
	}
	if !strings.Contains(string(result), "stage: test") {
		t.Error("stage field should be updated to test")
	}
	parsed, err := NewParser().ParseTaskFile(path, string(result))
	if err != nil {
		t.Fatalf("ParseTaskFile() error = %v", err)
	}
	if parsed.Stage != StageTest {
		t.Errorf("Stage = %q, want test", parsed.Stage)
	}
}

func TestWriteTaskStatus_removesLegacyImplementationFieldsOutsideInProgress(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "task.md")
	content := `---
id: E01/T-011
status: done
stage: implementation
phase: implementation
objective: "Legacy completed task"
---`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	fi, _ := os.Stat(path)
	task := &Task{
		ID:     "E01/T-011",
		Column: ColumnDone,
	}

	if err := WriteTaskStatus(path, task, fi.ModTime()); err != nil {
		t.Fatalf("WriteTaskStatus() error = %v", err)
	}

	result, _ := os.ReadFile(path)
	if strings.Contains(string(result), "phase:") {
		t.Error("legacy phase field should be removed")
	}
	if strings.Contains(string(result), "stage:") {
		t.Error("stale stage field should be removed")
	}

	parsed, err := NewParser().ParseTaskFile(path, string(result))
	if err != nil {
		t.Fatalf("ParseTaskFile() error = %v", err)
	}
	if parsed.Column != ColumnDone {
		t.Errorf("Column = %q, want done", parsed.Column)
	}
	if parsed.Stage != "" {
		t.Errorf("Stage = %q, want empty", parsed.Stage)
	}
}

func TestWriteTaskStatus_preservesBodyWithMultipleLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "task.md")
	content := `---
id: E01/T-006
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
		ID:     "E01/T-006",
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
task: E03-board-tui-core/T-004-render
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
		Task:       "E05-phase-transitions/T-004-write-router",
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
	if parsed.Task != "E05-phase-transitions/T-004-write-router" {
		t.Errorf("Task = %q, want E05-phase-transitions/T-004-write-router", parsed.Task)
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

func TestWriteRouterState_quotesNextActionWithColon(t *testing.T) {
	dir := t.TempDir()
	content := `## Current state

` + "```" + `yaml
state: task-building
release: v1
epic: E03-board-tui-core
task: ""
next_action: "Do the thing"
` + "```" + `
`

	path := filepath.Join(dir, "router.md")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	want := "Create tasks: status: planned"
	state := &RouterState{State: "task-building", NextAction: want}
	if err := WriteRouterState(dir, state, fi.ModTime()); err != nil {
		t.Fatalf("WriteRouterState() error = %v", err)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := NewRouterReader().ReadState(string(result))
	if err != nil {
		t.Fatalf("ReadState() error = %v", err)
	}
	if parsed.NextAction != want {
		t.Errorf("NextAction = %q, want %q", parsed.NextAction, want)
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
		ID:     "E01/T-007",
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
id: E01/T-008
status: planned
---`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	fi, _ := os.Stat(path)
	task := &Task{
		ID:     "E01/T-008",
		Column: ColumnDone,
		Stage:  StageAudit,
	}

	err := WriteTaskStatus(path, task, fi.ModTime())
	if err == nil {
		t.Fatal("WriteTaskStatus() expected invalid lifecycle error")
	}
}

func TestWriteTaskStatus_rejectsImplementationStageForInProgress(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "task.md")
	content := `---
id: E01/T-012
status: planned
---`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	fi, _ := os.Stat(path)
	task := &Task{
		ID:     "E01/T-012",
		Column: ColumnInProgress,
		Stage:  LegacyTaskStageImplementation,
	}

	err := WriteTaskStatus(path, task, fi.ModTime())
	if err == nil {
		t.Fatal("WriteTaskStatus() expected invalid stage error")
	}
	if !strings.Contains(err.Error(), `invalid stage "implementation"`) {
		t.Fatalf("WriteTaskStatus() error = %v, want invalid implementation stage message", err)
	}
}

func TestWriteTaskStatus_rejectsInvalidComplexityOnInProgress(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "task.md")
	content := `---
id: E19/T-011
status: planned
objective: "Invalid complexity"
---`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	fi, _ := os.Stat(path)
	task := &Task{
		ID:               "E19/T-011",
		Column:           ColumnInProgress,
		Stage:            StageBuild,
		ComplexityTier:   ComplexityTier("extreme"),
		ComplexityReason: "Invalid tier should be rejected before writing.",
	}

	err := WriteTaskStatus(path, task, fi.ModTime())
	if err == nil {
		t.Fatal("WriteTaskStatus() expected invalid complexity error")
	}
	if !strings.Contains(err.Error(), "invalid complexity_tier") {
		t.Fatalf("WriteTaskStatus() error = %v, want invalid complexity_tier", err)
	}
}

func TestWriteTaskStatus_trimsOverlongComplexityReasonByWordCount(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "task.md")
	content := `---
id: E19/T-012
status: planned
objective: "Invalid complexity"
---`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	fi, _ := os.Stat(path)
	task := &Task{
		ID:               "E19/T-012",
		Column:           ColumnInProgress,
		Stage:            StageBuild,
		ComplexityTier:   ComplexityHigh,
		ComplexityReason: strings.TrimSpace(strings.Repeat("word ", MaxComplexityReasonWords+1)),
	}

	if err := WriteTaskStatus(path, task, fi.ModTime()); err != nil {
		t.Fatalf("WriteTaskStatus() error = %v", err)
	}

	result, _ := os.ReadFile(path)
	parsed, err := NewParser().ParseTaskFile(path, string(result))
	if err != nil {
		t.Fatalf("ParseTaskFile() error = %v", err)
	}
	if got := ComplexityReasonWordCount(parsed.ComplexityReason); got != MaxComplexityReasonWords {
		t.Fatalf("ComplexityReasonWordCount() = %d, want %d", got, MaxComplexityReasonWords)
	}
}

func TestWriteTaskStatus_preservesComplexityFields(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "task.md")
	content := `---
id: E19/T-001
status: planned
complexity_tier: high
complexity_reason: "Requires coordinated changes across multiple packages."
objective: "Complexity test"
---

# Body`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	fi, _ := os.Stat(path)
	task := &Task{
		ID:     "E19/T-001",
		Column: ColumnInProgress,
		Stage:  StageBuild,
	}

	if err := WriteTaskStatus(path, task, fi.ModTime()); err != nil {
		t.Fatalf("WriteTaskStatus() error = %v", err)
	}

	result, _ := os.ReadFile(path)
	if !strings.Contains(string(result), "complexity_tier: high") {
		t.Error("complexity_tier not preserved after WriteTaskStatus")
	}
	if !strings.Contains(string(result), "complexity_reason:") {
		t.Error("complexity_reason not preserved after WriteTaskStatus")
	}

	p := NewParser()
	parsed, err := p.ParseTaskFile(path, string(result))
	if err != nil {
		t.Fatalf("ParseTaskFile() error = %v", err)
	}
	if parsed.ComplexityTier != ComplexityHigh {
		t.Errorf("ComplexityTier = %q, want high", parsed.ComplexityTier)
	}
	if parsed.ComplexityReason != "Requires coordinated changes across multiple packages." {
		t.Errorf("ComplexityReason = %q, want reason text", parsed.ComplexityReason)
	}
}

func TestWriteTaskStatus_selfHealsKnownComplexityBlockers(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "task.md")
	longReason := strings.TrimSpace(strings.Repeat("word ", MaxComplexityReasonWords+1))
	content := "---\nid: E19/T-013\nstatus: planned\ncomplexity_tier: small\ncomplexity_reason: \"" + longReason + "\"\nobjective: \"Complexity repair\"\n---\n\n# Body"

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	parsed, err := NewParser().ParseTaskFile(path, content)
	if err != nil {
		t.Fatalf("ParseTaskFile() error = %v", err)
	}
	parsed.Column = ColumnInProgress
	parsed.Stage = StageBuild

	fi, _ := os.Stat(path)
	if err := WriteTaskStatus(path, parsed, fi.ModTime()); err != nil {
		t.Fatalf("WriteTaskStatus() error = %v", err)
	}

	result, _ := os.ReadFile(path)
	if strings.Contains(string(result), "complexity_tier: small") {
		t.Fatal("complexity_tier alias should be repaired")
	}
	if !strings.Contains(string(result), "complexity_tier: low") {
		t.Fatal("complexity_tier should be written as low")
	}

	reparsed, err := NewParser().ParseTaskFile(path, string(result))
	if err != nil {
		t.Fatalf("ParseTaskFile() after repair error = %v", err)
	}
	if reparsed.ComplexityTier != ComplexityLow {
		t.Fatalf("ComplexityTier = %q, want low", reparsed.ComplexityTier)
	}
	if got := ComplexityReasonWordCount(reparsed.ComplexityReason); got != MaxComplexityReasonWords {
		t.Fatalf("ComplexityReasonWordCount() = %d, want %d", got, MaxComplexityReasonWords)
	}
}

func TestWriteDefectStatus_updatesStatusAndPreservesBody(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "D001.md")
	content := `---
id: v1/D001
release: v1
status: open
severity: high
title: "Crash"
reference: E01/T-001
---

# Body

Keep this text.`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	defect := &Defect{
		ID:       "v1/D001",
		Release:  "v1",
		Status:   DefectResolved,
		Severity: SeverityHigh,
		Title:    "Crash",
	}

	if err := WriteDefectStatus(path, defect, fi.ModTime()); err != nil {
		t.Fatalf("WriteDefectStatus() error = %v", err)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	parsed, err := NewParser().ParseDefectFile(path, string(result))
	if err != nil {
		t.Fatalf("ParseDefectFile() error = %v", err)
	}
	if parsed.Status != DefectResolved {
		t.Errorf("Status = %v, want resolved", parsed.Status)
	}
	if !strings.Contains(string(result), "reference: E01/T-001") {
		t.Error("unrelated frontmatter field not preserved")
	}
	if !strings.Contains(string(result), "Keep this text.") {
		t.Error("body content not preserved")
	}
}

func TestWriteObjectiveV2_updatesStatusPreservesUnknownFieldsAndBody(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "Objective.md")
	content := `---
id: O-002
title: "Load V2 work with stable identity"
status: planned
depends_on: [O-001]
release: R-001
owner:
  team: platform
  contact: "team@example.com"
---

# Objective

## Notes

Authored planning notes that must survive the rewrite.`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	objective, err := DecodeObjectiveV2(path, content)
	if err != nil {
		t.Fatalf("DecodeObjectiveV2() error = %v", err)
	}
	objective.Status = ColumnInProgress

	if err := WriteObjectiveV2(objective); err != nil {
		t.Fatalf("WriteObjectiveV2() error = %v", err)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	reparsed, err := DecodeObjectiveV2(path, string(result))
	if err != nil {
		t.Fatalf("DecodeObjectiveV2() after write error = %v", err)
	}
	if reparsed.Status != ColumnInProgress {
		t.Errorf("Status = %q, want in_progress", reparsed.Status)
	}
	if len(reparsed.DependsOn) != 1 || reparsed.DependsOn[0] != "O-001" {
		t.Errorf("DependsOn = %v, want [O-001] preserved", reparsed.DependsOn)
	}
	if reparsed.Release != "R-001" {
		t.Errorf("Release = %q, want R-001 preserved", reparsed.Release)
	}
	if !strings.Contains(string(result), "team: platform") {
		t.Error("unknown nested field not preserved")
	}
	if !strings.Contains(string(result), "team@example.com") {
		t.Error("unknown nested field value not preserved")
	}
	if !strings.Contains(string(result), "Authored planning notes that must survive the rewrite.") {
		t.Error("authored body content not preserved")
	}
	if !strings.Contains(string(result), "## Notes") {
		t.Error("authored heading not preserved")
	}
}

func TestWriteObjectiveV2_noOpLeavesBytesAndMtimeUnchanged(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "Objective.md")
	content := `---
id: O-003
title: "No-op objective"
status: in_progress
---

# Objective`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	before, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	beforeBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(10 * time.Millisecond)

	objective, err := DecodeObjectiveV2(path, content)
	if err != nil {
		t.Fatalf("DecodeObjectiveV2() error = %v", err)
	}

	if err := WriteObjectiveV2(objective); err != nil {
		t.Fatalf("WriteObjectiveV2() error = %v", err)
	}

	after, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	afterBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if !before.ModTime().Equal(after.ModTime()) {
		t.Errorf("ModTime changed on no-op write: before %v, after %v", before.ModTime(), after.ModTime())
	}
	if string(beforeBytes) != string(afterBytes) {
		t.Error("file bytes changed on no-op write")
	}
}

func TestWriteObjectiveV2_refusesUnsupportedStatusAndLeavesFileUntouched(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "Objective.md")
	content := `---
id: O-004
title: "Guarded objective"
status: planned
---

# Objective`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	objective, err := DecodeObjectiveV2(path, content)
	if err != nil {
		t.Fatalf("DecodeObjectiveV2() error = %v", err)
	}
	objective.Status = ColumnType("archived")

	err = WriteObjectiveV2(objective)
	if err == nil {
		t.Fatal("WriteObjectiveV2() expected error for unsupported status")
	}
	if !errors.Is(err, ErrV2InvalidLifecycle) {
		t.Fatalf("WriteObjectiveV2() error = %v, want ErrV2InvalidLifecycle", err)
	}
	if !strings.Contains(err.Error(), path) {
		t.Errorf("error %v is not path-qualified with %q", err, path)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(result) != content {
		t.Error("file content changed despite validation refusal")
	}
}

func TestWriteTaskV2_updatesStatusAndStagePreservesDependencies(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "T-005.md")
	content := `---
id: T-005
title: "Show clear project errors"
objective: O-002
planned_by: {role: planner, session: planning-fixture}
status: planned
depends_on:
  - task: T-003
  - task: T-004
    requires: accepted
release: R-001
---

# Task

Authored task notes.`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	task, err := DecodeTaskV2(path, content)
	if err != nil {
		t.Fatalf("DecodeTaskV2() error = %v", err)
	}
	task.Status = ColumnInProgress
	task.Stage = StageBuild

	if err := WriteTaskV2(task); err != nil {
		t.Fatalf("WriteTaskV2() error = %v", err)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	reparsed, err := DecodeTaskV2(path, string(result))
	if err != nil {
		t.Fatalf("DecodeTaskV2() after write error = %v", err)
	}
	if reparsed.Status != ColumnInProgress || reparsed.Stage != StageBuild {
		t.Errorf("Status/Stage = %q/%q, want in_progress/build", reparsed.Status, reparsed.Stage)
	}
	if len(reparsed.DependsOn) != 2 {
		t.Fatalf("DependsOn len = %d, want 2 preserved", len(reparsed.DependsOn))
	}
	if reparsed.DependsOn[1].Task != "T-004" || reparsed.DependsOn[1].Requires != TaskDependencyAccepted {
		t.Errorf("DependsOn[1] = %+v, want T-004/accepted preserved", reparsed.DependsOn[1])
	}
	if !strings.Contains(string(result), "Authored task notes.") {
		t.Error("authored body content not preserved")
	}
}

func TestWriteTaskV2_removesStageWhenLeavingInProgress(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "T-006.md")
	content := `---
id: T-006
title: "Finish up"
objective: O-002
planned_by: {role: planner, session: planning-fixture}
status: in_progress
stage: audit
---

# Task`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	task, err := DecodeTaskV2(path, content)
	if err != nil {
		t.Fatalf("DecodeTaskV2() error = %v", err)
	}
	task.Status = ColumnDone
	task.Stage = ""

	if err := WriteTaskV2(task); err != nil {
		t.Fatalf("WriteTaskV2() error = %v", err)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(result), "stage:") {
		t.Error("stage field should be removed when leaving in_progress")
	}

	reparsed, err := DecodeTaskV2(path, string(result))
	if err != nil {
		t.Fatalf("DecodeTaskV2() after write error = %v", err)
	}
	if reparsed.Status != ColumnDone || reparsed.Stage != "" {
		t.Errorf("Status/Stage = %q/%q, want done/empty", reparsed.Status, reparsed.Stage)
	}
}

func TestWriteTaskV2_refusesInProgressWithoutStageAndLeavesFileUntouched(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "T-007.md")
	content := `---
id: T-007
title: "Guarded task"
objective: O-002
planned_by: {role: planner, session: planning-fixture}
status: planned
---

# Task`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	task, err := DecodeTaskV2(path, content)
	if err != nil {
		t.Fatalf("DecodeTaskV2() error = %v", err)
	}
	task.Status = ColumnInProgress
	task.Stage = ""

	err = WriteTaskV2(task)
	if err == nil {
		t.Fatal("WriteTaskV2() expected error for in_progress without stage")
	}
	if !errors.Is(err, ErrV2InvalidLifecycle) {
		t.Fatalf("WriteTaskV2() error = %v, want ErrV2InvalidLifecycle", err)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(result) != content {
		t.Error("file content changed despite validation refusal")
	}
}

func TestWriteTaskV2_noOpLeavesBytesAndMtimeUnchanged(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "T-008.md")
	content := `---
id: T-008
title: "No-op task"
objective: O-002
planned_by: {role: planner, session: planning-fixture}
status: in_progress
stage: test
---

# Task`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	before, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	beforeBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(10 * time.Millisecond)

	task, err := DecodeTaskV2(path, content)
	if err != nil {
		t.Fatalf("DecodeTaskV2() error = %v", err)
	}

	if err := WriteTaskV2(task); err != nil {
		t.Fatalf("WriteTaskV2() error = %v", err)
	}

	after, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	afterBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if !before.ModTime().Equal(after.ModTime()) {
		t.Errorf("ModTime changed on no-op write: before %v, after %v", before.ModTime(), after.ModTime())
	}
	if string(beforeBytes) != string(afterBytes) {
		t.Error("file bytes changed on no-op write")
	}
}

func TestWriteTaskV2_preservesCRLFLineEndings(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "T-009.md")
	content := "---\r\nid: T-009\r\ntitle: \"CRLF task\"\r\nobjective: O-002\r\nplanned_by: {role: planner, session: planning-fixture}\r\nstatus: planned\r\n---\r\n\r\n# Task\r\n\r\nAuthored notes.\r\n"

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	task, err := DecodeTaskV2(path, content)
	if err != nil {
		t.Fatalf("DecodeTaskV2() error = %v", err)
	}
	task.Status = ColumnInProgress
	task.Stage = StageBuild

	if err := WriteTaskV2(task); err != nil {
		t.Fatalf("WriteTaskV2() error = %v", err)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Count(string(result), "\n") != strings.Count(string(result), "\r\n") {
		t.Errorf("result did not preserve CRLF line endings throughout: %q", string(result))
	}
	if !strings.Contains(string(result), "Authored notes.") {
		t.Error("authored body content not preserved across CRLF rewrite")
	}

	reparsed, err := DecodeTaskV2(path, string(result))
	if err != nil {
		t.Fatalf("DecodeTaskV2() after write error = %v", err)
	}
	if reparsed.Status != ColumnInProgress || reparsed.Stage != StageBuild {
		t.Errorf("Status/Stage = %q/%q, want in_progress/build", reparsed.Status, reparsed.Stage)
	}
}

func TestWriteTaskV2_preservesLFLineEndingsByDefault(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "T-010.md")
	content := "---\nid: T-010\ntitle: \"LF task\"\nobjective: O-002\nplanned_by: {role: planner, session: planning-fixture}\nstatus: planned\n---\n\n# Task\n"

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	task, err := DecodeTaskV2(path, content)
	if err != nil {
		t.Fatalf("DecodeTaskV2() error = %v", err)
	}
	task.Status = ColumnInProgress
	task.Stage = StageBuild

	if err := WriteTaskV2(task); err != nil {
		t.Fatalf("WriteTaskV2() error = %v", err)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(result), "\r\n") {
		t.Errorf("LF source should not gain CRLF endings: %q", string(result))
	}
}

func TestWriteObjectiveV2_writeFailureIsPathQualifiedAndLeavesRecordIntact(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("running as root: permission checks do not apply")
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "Objective.md")
	content := `---
id: O-005
title: "Read-only objective"
status: planned
---

# Objective`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	objective, err := DecodeObjectiveV2(path, content)
	if err != nil {
		t.Fatalf("DecodeObjectiveV2() error = %v", err)
	}
	objective.Status = ColumnInProgress

	if err := os.Chmod(path, 0444); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(path, 0644)

	err = WriteObjectiveV2(objective)
	if err == nil {
		t.Fatal("WriteObjectiveV2() expected error writing to a read-only file")
	}
	if !strings.Contains(err.Error(), path) {
		t.Errorf("error %v is not path-qualified with %q", err, path)
	}

	if err := os.Chmod(path, 0644); err != nil {
		t.Fatal(err)
	}
	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(result) != content {
		t.Error("existing record was truncated or altered despite the write failure")
	}
}

func TestWriteV2Record_resolvesDiscoveredRelativePathFromProjectRoot(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First objective")
	writeV2TaskFixture(t, root, "O-001-first", "T-001-first.md", "T-001", "First task", "O-001")

	objectives, tasks, err := DiscoverV2Records(root)
	if err != nil {
		t.Fatalf("DiscoverV2Records() error = %v", err)
	}
	objective := objectives["O-001"]
	task := tasks["T-001"]
	if objective == nil || task == nil {
		t.Fatal("DiscoverV2Records() did not return both V2 records")
	}
	if filepath.IsAbs(objective.Source.Path) || filepath.IsAbs(task.Source.Path) {
		t.Fatalf("source paths must remain project-relative: objective=%q task=%q", objective.Source.Path, task.Source.Path)
	}
	if objective.Source.ProjectRoot != root || task.Source.ProjectRoot != root {
		t.Fatalf("source project roots = %q/%q, want %q", objective.Source.ProjectRoot, task.Source.ProjectRoot, root)
	}

	otherWorkingDir := t.TempDir()
	t.Chdir(otherWorkingDir)

	objective.Status = ColumnInProgress
	if err := WriteObjectiveV2(objective); err != nil {
		t.Fatalf("WriteObjectiveV2() from unrelated cwd error = %v", err)
	}
	task.Status = ColumnInProgress
	task.Stage = StageBuild
	if err := WriteTaskV2(task); err != nil {
		t.Fatalf("WriteTaskV2() from unrelated cwd error = %v", err)
	}

	objectiveBytes, err := os.ReadFile(filepath.Join(root, objective.Source.Path))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(objectiveBytes), "status: in_progress") {
		t.Errorf("objective file at project root was not updated: %s", objectiveBytes)
	}
	taskBytes, err := os.ReadFile(filepath.Join(root, task.Source.Path))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(taskBytes), "status: in_progress") || !strings.Contains(string(taskBytes), "stage: build") {
		t.Errorf("task file at project root was not updated: %s", taskBytes)
	}
}

func TestResolveV2SourcePath_confinesPathsToProjectRoot(t *testing.T) {
	root := t.TempDir()
	inside := filepath.Join(root, "objectives", "O-001", "Objective.md")

	tests := []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		{
			name: "relative inside",
			path: "objectives/O-001/Objective.md",
			want: inside,
		},
		{
			name:    "relative escape",
			path:    filepath.Join("..", "outside.md"),
			wantErr: true,
		},
		{
			name: "absolute inside",
			path: inside,
			want: inside,
		},
		{
			name:    "absolute outside",
			path:    filepath.Join(root, "..", "outside.md"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveV2SourcePath(V2SourceDocument{
				Path:        tt.path,
				ProjectRoot: root,
			})
			if tt.wantErr {
				if !errors.Is(err, ErrV2UnsafePath) {
					t.Fatalf("resolveV2SourcePath() error = %v, want ErrV2UnsafePath", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("resolveV2SourcePath() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("resolveV2SourcePath() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestResolveV2SourcePath_allowsAbsolutePathWithoutProjectRoot(t *testing.T) {
	path := filepath.Join(t.TempDir(), "record.md")

	got, err := resolveV2SourcePath(V2SourceDocument{Path: path})
	if err != nil {
		t.Fatalf("resolveV2SourcePath() error = %v", err)
	}
	if got != path {
		t.Errorf("resolveV2SourcePath() = %q, want %q", got, path)
	}
}

func TestWriteV2Record_refusesStaleLoadedSourceWithoutOverwritingUserEdit(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First objective")
	writeV2TaskFixture(t, root, "O-001-first", "T-001-first.md", "T-001", "First task", "O-001")

	_, tasks, err := DiscoverV2Records(root)
	if err != nil {
		t.Fatalf("DiscoverV2Records() error = %v", err)
	}
	task := tasks["T-001"]
	if task == nil {
		t.Fatal("DiscoverV2Records() did not return T-001")
	}

	path := filepath.Join(root, task.Source.Path)
	userEdit := "---\nid: T-001\ntitle: \"First task\"\nobjective: O-001\nplanned_by: {role: planner, session: planning-fixture}\nstatus: planned\neditor_note: \"owner edit\"\n---\n\n# First task\n\nOwner edit must survive.\n"
	if err := os.WriteFile(path, []byte(userEdit), 0644); err != nil {
		t.Fatal(err)
	}

	task.Status = ColumnInProgress
	task.Stage = StageBuild
	err = WriteTaskV2(task)
	if !errors.Is(err, ErrV2SourceConflict) {
		t.Fatalf("WriteTaskV2() error = %v, want ErrV2SourceConflict", err)
	}
	if !errors.Is(err, ErrMtimeConflict) {
		t.Fatalf("WriteTaskV2() error = %v, want ErrMtimeConflict compatibility marker", err)
	}
	if !strings.Contains(err.Error(), path) {
		t.Errorf("WriteTaskV2() error = %v, want path %q", err, path)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != userEdit {
		t.Fatalf("stale write changed user content:\n got: %s\nwant: %s", got, userEdit)
	}
}

func TestWriteTaskEvidenceV2_setsAllSubBlocksPreservesUnknownFieldsAndBody(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "T-020.md")
	content := `---
id: T-020
title: "No evidence yet"
objective: O-002
planned_by: {role: planner, session: planning-fixture}
status: in_progress
stage: build
depends_on:
  - task: T-010
release: R-001
---

# Task

Authored task notes.`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	task, err := DecodeTaskV2(path, content)
	if err != nil {
		t.Fatalf("DecodeTaskV2() error = %v", err)
	}
	if task.Evidence != nil {
		t.Fatalf("Evidence = %+v, want nil before write", task.Evidence)
	}

	assessedAt := time.Date(2026, 9, 14, 0, 0, 0, 0, time.UTC)
	recordedAt := time.Date(2026, 9, 14, 1, 0, 0, 0, time.UTC)
	replanAt := time.Date(2026, 9, 14, 2, 0, 0, 0, time.UTC)
	waiverAt := time.Date(2026, 9, 14, 3, 0, 0, 0, time.UTC)
	task.Evidence = &Evidence{
		CheckWaiver: &CheckWaiver{
			Task:       "T-020",
			Reason:     "Owner completed via board without requesting a Check.",
			Actor:      Actor{Role: ActorRoleOwner, Session: "board-owner"},
			RecordedAt: waiverAt,
		},
		LastCheck: "C-001",
		Freshness: &Freshness{
			State:      FreshnessCurrent,
			Check:      "C-001",
			AssessedBy: Actor{Role: ActorRoleChecker, Session: "sess-1"},
			AssessedAt: assessedAt,
			Basis:      "Reviewed diff against AC.",
		},
		OwnerValidation: &OwnerValidation{Required: true, AcceptedCheck: "C-001", AcceptedBy: Actor{Role: ActorRoleOwner, Session: "owner-1"}},
		Exception: &Exception{
			Requirements: []string{"TEST-08"},
			Reason:       "Owner accepted known risk.",
			Owner:        "owner-1",
			RecordedAt:   recordedAt,
			Check:        "C-001",
		},
		Replan: &Replan{
			Reason:     "Plan needs revisiting.",
			RecordedBy: Actor{Role: ActorRolePlanner, Session: "sess-2"},
			RecordedAt: replanAt,
		},
	}

	if err := WriteTaskEvidenceV2(task); err != nil {
		t.Fatalf("WriteTaskEvidenceV2() error = %v", err)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	reparsed, err := DecodeTaskV2(path, string(result))
	if err != nil {
		t.Fatalf("DecodeTaskV2() after write error = %v", err)
	}
	if reparsed.Evidence == nil {
		t.Fatal("Evidence = nil, want populated evidence")
	}
	if reparsed.Evidence.LastCheck != "C-001" {
		t.Errorf("LastCheck = %q, want C-001", reparsed.Evidence.LastCheck)
	}
	if reparsed.Evidence.Freshness == nil || reparsed.Evidence.Freshness.State != FreshnessCurrent {
		t.Errorf("Freshness = %+v, want state current", reparsed.Evidence.Freshness)
	}
	if reparsed.Evidence.OwnerValidation == nil || !reparsed.Evidence.OwnerValidation.Required {
		t.Errorf("OwnerValidation = %+v, want required true", reparsed.Evidence.OwnerValidation)
	}
	if reparsed.Evidence.OwnerValidation.AcceptedBy != (Actor{Role: ActorRoleOwner, Session: "owner-1"}) {
		t.Errorf("OwnerValidation.AcceptedBy = %+v, want owner/owner-1", reparsed.Evidence.OwnerValidation.AcceptedBy)
	}
	if reparsed.Evidence.Exception == nil || reparsed.Evidence.Exception.Owner != "owner-1" {
		t.Errorf("Exception = %+v, want owner owner-1", reparsed.Evidence.Exception)
	}
	if reparsed.Evidence.Replan == nil || reparsed.Evidence.Replan.Reason != "Plan needs revisiting." {
		t.Errorf("Replan = %+v, want reason set", reparsed.Evidence.Replan)
	}
	if reparsed.Evidence.CheckWaiver == nil || reparsed.Evidence.CheckWaiver.Task != "T-020" || reparsed.Evidence.CheckWaiver.Actor.Role != ActorRoleOwner {
		t.Errorf("CheckWaiver = %+v, want task T-020 recorded by owner", reparsed.Evidence.CheckWaiver)
	}

	if reparsed.Status != ColumnInProgress || reparsed.Stage != StageBuild {
		t.Errorf("Status/Stage = %q/%q, want in_progress/build preserved", reparsed.Status, reparsed.Stage)
	}
	if len(reparsed.DependsOn) != 1 || reparsed.DependsOn[0].Task != "T-010" {
		t.Errorf("DependsOn = %v, want [T-010] preserved", reparsed.DependsOn)
	}
	if !strings.Contains(string(result), "Authored task notes.") {
		t.Error("authored body content not preserved")
	}
}

func TestWriteTaskEvidenceV2_removesReplanKeyRatherThanEmptyValue(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "T-021.md")
	content := `---
id: T-021
title: "Replan flagged"
objective: O-002
planned_by: {role: planner, session: planning-fixture}
status: in_progress
stage: build
last_check: C-001
replan:
  reason: "Plan needs revisiting."
  recorded_by:
    role: planner
    session: sess-2
  recorded_at: '2026-09-14T02:00:00Z'
---

# Task`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	task, err := DecodeTaskV2(path, content)
	if err != nil {
		t.Fatalf("DecodeTaskV2() error = %v", err)
	}
	if task.Evidence == nil || task.Evidence.Replan == nil {
		t.Fatalf("Evidence = %+v, want replan present before write", task.Evidence)
	}

	task.Evidence.Replan = nil

	if err := WriteTaskEvidenceV2(task); err != nil {
		t.Fatalf("WriteTaskEvidenceV2() error = %v", err)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(result), "replan:") {
		t.Error("replan key should be removed, not written empty")
	}

	reparsed, err := DecodeTaskV2(path, string(result))
	if err != nil {
		t.Fatalf("DecodeTaskV2() after write error = %v", err)
	}
	if reparsed.Evidence == nil {
		t.Fatal("Evidence = nil, want last_check preserved")
	}
	if reparsed.Evidence.Replan != nil {
		t.Errorf("Replan = %+v, want nil after clearing", reparsed.Evidence.Replan)
	}
	if reparsed.Evidence.LastCheck != "C-001" {
		t.Errorf("LastCheck = %q, want C-001 preserved", reparsed.Evidence.LastCheck)
	}
}

func TestWriteTaskEvidenceV2_noOpLeavesBytesAndMtimeUnchanged(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "T-022.md")
	content := `---
id: T-022
title: "Fully evidenced task"
objective: O-002
planned_by: {role: planner, session: planning-fixture}
status: in_progress
stage: build
last_check: C-001
freshness:
  state: current
  check: C-001
  assessed_by:
    role: checker
    session: sess-1
  assessed_at: '2026-09-14T00:00:00Z'
  basis: "Reviewed diff against AC."
owner_validation:
  required: true
  accepted_check: C-001
  accepted_by:
    role: owner
    session: owner-1
exception:
  requirements:
    - TEST-08
  reason: "Owner accepted known risk."
  owner: owner-1
  recorded_at: '2026-09-14T01:00:00Z'
  check: C-001
replan:
  reason: "Plan needs revisiting."
  recorded_by:
    role: planner
    session: sess-2
  recorded_at: '2026-09-14T02:00:00Z'
---

# Task`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	before, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	beforeBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(10 * time.Millisecond)

	task, err := DecodeTaskV2(path, content)
	if err != nil {
		t.Fatalf("DecodeTaskV2() error = %v", err)
	}

	if err := WriteTaskEvidenceV2(task); err != nil {
		t.Fatalf("WriteTaskEvidenceV2() error = %v", err)
	}

	after, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	afterBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if !before.ModTime().Equal(after.ModTime()) {
		t.Errorf("ModTime changed on no-op write: before %v, after %v", before.ModTime(), after.ModTime())
	}
	if string(beforeBytes) != string(afterBytes) {
		t.Error("file bytes changed on no-op write")
	}
}

func TestWriteTaskEvidenceV2_refusesStaleSourceWithoutOverwritingUserEdit(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "T-023.md")
	content := `---
id: T-023
title: "Guarded evidence write"
objective: O-002
planned_by: {role: planner, session: planning-fixture}
status: planned
---

# Task`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	task, err := DecodeTaskV2(path, content)
	if err != nil {
		t.Fatalf("DecodeTaskV2() error = %v", err)
	}

	userEdit := "---\nid: T-023\ntitle: \"Guarded evidence write\"\nobjective: O-002\nplanned_by: {role: planner, session: planning-fixture}\nstatus: planned\neditor_note: \"owner edit\"\n---\n\n# Task\n\nOwner edit must survive.\n"
	if err := os.WriteFile(path, []byte(userEdit), 0644); err != nil {
		t.Fatal(err)
	}

	task.Evidence = &Evidence{LastCheck: "C-001"}

	err = WriteTaskEvidenceV2(task)
	if !errors.Is(err, ErrV2SourceConflict) {
		t.Fatalf("WriteTaskEvidenceV2() error = %v, want ErrV2SourceConflict", err)
	}
	if !errors.Is(err, ErrMtimeConflict) {
		t.Fatalf("WriteTaskEvidenceV2() error = %v, want ErrMtimeConflict compatibility marker", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != userEdit {
		t.Fatalf("stale write changed user content:\n got: %s\nwant: %s", got, userEdit)
	}
}

func TestWriteTaskEvidenceV2_preservesCRLFLineEndings(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "T-024.md")
	content := "---\r\nid: T-024\r\ntitle: \"CRLF evidence task\"\r\nobjective: O-002\r\nplanned_by: {role: planner, session: planning-fixture}\r\nstatus: planned\r\n---\r\n\r\n# Task\r\n\r\nAuthored notes.\r\n"

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	task, err := DecodeTaskV2(path, content)
	if err != nil {
		t.Fatalf("DecodeTaskV2() error = %v", err)
	}
	task.Evidence = &Evidence{LastCheck: "C-001"}

	if err := WriteTaskEvidenceV2(task); err != nil {
		t.Fatalf("WriteTaskEvidenceV2() error = %v", err)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(result), "\n") != strings.Count(string(result), "\r\n") {
		t.Errorf("result did not preserve CRLF line endings throughout: %q", string(result))
	}
	if !strings.Contains(string(result), "Authored notes.") {
		t.Error("authored body content not preserved across CRLF rewrite")
	}
}

func TestWriteTaskEvidenceV2_rejectsMalformedEvidenceLeavesFileUntouched(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "T-025.md")
	content := `---
id: T-025
title: "Guarded malformed evidence"
objective: O-002
planned_by: {role: planner, session: planning-fixture}
status: planned
---

# Task`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	task, err := DecodeTaskV2(path, content)
	if err != nil {
		t.Fatalf("DecodeTaskV2() error = %v", err)
	}
	task.Evidence = &Evidence{
		Freshness: &Freshness{
			State:      FreshnessState("bogus"),
			Check:      "C-001",
			AssessedBy: Actor{Role: ActorRoleChecker, Session: "sess-1"},
			AssessedAt: time.Date(2026, 9, 14, 0, 0, 0, 0, time.UTC),
			Basis:      "Reviewed diff against AC.",
		},
	}

	err = WriteTaskEvidenceV2(task)
	if err == nil {
		t.Fatal("WriteTaskEvidenceV2() expected error for malformed freshness state")
	}
	if !errors.Is(err, ErrV2EvidenceMalformed) {
		t.Fatalf("WriteTaskEvidenceV2() error = %v, want ErrV2EvidenceMalformed", err)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(result) != content {
		t.Error("file content changed despite validation refusal")
	}
}

func TestCreateCheckV2_writesNewFileAllocatesFirstID(t *testing.T) {
	root := t.TempDir()
	index := &V2Index{Checks: map[string]*CheckV2{}}

	fields := NewCheckV2{
		Scope:           CheckScope{Kind: CheckScopeTask, ID: "T-001"},
		Result:          CheckResultClear,
		CheckedBy:       Actor{Role: ActorRoleChecker, Session: "sess-1"},
		ExecutedSession: "build-001",
		CheckedAt:       time.Date(2026, 9, 14, 0, 0, 0, 0, time.UTC),
		Issues:          []string{"I-001"},
		Body:            "\n\n# Check\n\nOutcome notes.\n",
	}

	check, err := CreateCheckV2(root, index, fields)
	if err != nil {
		t.Fatalf("CreateCheckV2() error = %v", err)
	}
	if check.ID != "C-001" {
		t.Errorf("ID = %q, want C-001", check.ID)
	}

	path := filepath.Join(root, "checks", "C-001.md")
	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("checks/C-001.md not written: %v", err)
	}
	if !strings.Contains(string(result), "Outcome notes.") {
		t.Error("authored body content not written")
	}

	reparsed, err := DecodeCheckV2(filepath.Join("checks", "C-001.md"), string(result))
	if err != nil {
		t.Fatalf("DecodeCheckV2() after create error = %v", err)
	}
	if reparsed.Scope.ID != "T-001" || reparsed.Result != CheckResultClear {
		t.Errorf("reparsed = %+v, want scope T-001/CLEAR", reparsed)
	}
}

func TestCreateCheckV2_allocatesNextIDOverPopulatedIndex(t *testing.T) {
	root := t.TempDir()
	index := &V2Index{Checks: map[string]*CheckV2{
		"C-001": {ID: "C-001"},
		"C-002": {ID: "C-002"},
	}}

	fields := NewCheckV2{
		Scope:           CheckScope{Kind: CheckScopeTask, ID: "T-001"},
		Result:          CheckResultNeedsWork,
		CheckedBy:       Actor{Role: ActorRoleChecker, Session: "sess-1"},
		ExecutedSession: "build-001",
		CheckedAt:       time.Date(2026, 9, 14, 0, 0, 0, 0, time.UTC),
		Body:            "\n\n# Check\n",
	}

	check, err := CreateCheckV2(root, index, fields)
	if err != nil {
		t.Fatalf("CreateCheckV2() error = %v", err)
	}
	if check.ID != "C-003" {
		t.Errorf("ID = %q, want C-003 as next unused id", check.ID)
	}
}

func TestCreateCheckV2_refusesExistingPathAndLeavesItUntouched(t *testing.T) {
	root := t.TempDir()
	checksDir := filepath.Join(root, "checks")
	if err := os.MkdirAll(checksDir, 0755); err != nil {
		t.Fatal(err)
	}
	leftover := "not a check record"
	if err := os.WriteFile(filepath.Join(checksDir, "C-001.md"), []byte(leftover), 0644); err != nil {
		t.Fatal(err)
	}

	index := &V2Index{Checks: map[string]*CheckV2{}}
	fields := NewCheckV2{
		Scope:           CheckScope{Kind: CheckScopeTask, ID: "T-001"},
		Result:          CheckResultClear,
		CheckedBy:       Actor{Role: ActorRoleChecker, Session: "sess-1"},
		ExecutedSession: "build-001",
		CheckedAt:       time.Date(2026, 9, 14, 0, 0, 0, 0, time.UTC),
		Body:            "\n\n# Check\n",
	}

	check, err := CreateCheckV2(root, index, fields)
	if check != nil {
		t.Errorf("CreateCheckV2() returned %+v, want nil on collision", check)
	}
	if !errors.Is(err, ErrV2CheckImmutable) {
		t.Fatalf("CreateCheckV2() error = %v, want ErrV2CheckImmutable", err)
	}

	result, err := os.ReadFile(filepath.Join(checksDir, "C-001.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(result) != leftover {
		t.Error("existing file at colliding path was modified")
	}
}

type fakeV2CheckTempFile struct {
	path     string
	contents []byte
	chmodErr error
	writeErr error
	writeN   int
	syncErr  error
	closeErr error
}

func (f *fakeV2CheckTempFile) Name() string { return f.path }

func (f *fakeV2CheckTempFile) Chmod(os.FileMode) error { return f.chmodErr }

func (f *fakeV2CheckTempFile) Write(content []byte) (int, error) {
	if f.writeErr != nil {
		return 0, f.writeErr
	}
	f.contents = append(f.contents[:0], content...)
	if f.writeN > 0 && f.writeN < len(content) {
		return f.writeN, nil
	}
	return len(content), nil
}

func (f *fakeV2CheckTempFile) Sync() error { return f.syncErr }

func (f *fakeV2CheckTempFile) Close() error { return f.closeErr }

type fakeV2CheckFileSystem struct {
	temp       *fakeV2CheckTempFile
	mkdirErr   error
	createErr  error
	removeErr  error
	linkErr    error
	published  map[string][]byte
	removeCall int
}

func (f *fakeV2CheckFileSystem) operations() v2CheckFileOperations {
	return v2CheckFileOperations{
		mkdirAll: func(string, os.FileMode) error {
			return f.mkdirErr
		},
		createTemp: func(dir, _ string) (v2CheckTempFile, error) {
			if f.createErr != nil {
				return nil, f.createErr
			}
			if f.temp == nil {
				f.temp = &fakeV2CheckTempFile{path: filepath.Join(dir, "temp-check")}
			}
			return f.temp, nil
		},
		remove: func(string) error {
			f.removeCall++
			return f.removeErr
		},
		link: func(_, path string) error {
			if f.linkErr != nil {
				return f.linkErr
			}
			if f.published == nil {
				f.published = map[string][]byte{}
			}
			f.published[path] = append([]byte(nil), f.temp.contents...)
			return nil
		},
	}
}

func TestCreateV2CheckFile_operationFailuresPreserveFinalState(t *testing.T) {
	path := filepath.Join(t.TempDir(), "checks", "C-001.md")
	wantContent := []byte("check content")
	failure := errors.New("injected failure")

	cases := []struct {
		name       string
		configure  func(*fakeV2CheckFileSystem)
		wantAbsent bool
		wantError  bool
		wantErrMsg string
	}{
		{name: "directory creation", configure: func(fs *fakeV2CheckFileSystem) { fs.mkdirErr = failure }, wantAbsent: true, wantError: true},
		{name: "temporary creation", configure: func(fs *fakeV2CheckFileSystem) { fs.createErr = failure }, wantAbsent: true, wantError: true},
		{name: "permission change", configure: func(fs *fakeV2CheckFileSystem) {
			fs.temp = &fakeV2CheckTempFile{path: filepath.Join(filepath.Dir(path), "temp-check"), chmodErr: failure}
		}, wantAbsent: true, wantError: true},
		{name: "write", configure: func(fs *fakeV2CheckFileSystem) {
			fs.temp = &fakeV2CheckTempFile{path: filepath.Join(filepath.Dir(path), "temp-check"), writeErr: failure}
		}, wantAbsent: true, wantError: true},
		{name: "short write", configure: func(fs *fakeV2CheckFileSystem) {
			fs.temp = &fakeV2CheckTempFile{path: filepath.Join(filepath.Dir(path), "temp-check"), writeN: 1}
		}, wantAbsent: true, wantError: true},
		{name: "sync", configure: func(fs *fakeV2CheckFileSystem) {
			fs.temp = &fakeV2CheckTempFile{path: filepath.Join(filepath.Dir(path), "temp-check"), syncErr: failure}
		}, wantAbsent: true, wantError: true},
		{name: "close", configure: func(fs *fakeV2CheckFileSystem) {
			fs.temp = &fakeV2CheckTempFile{path: filepath.Join(filepath.Dir(path), "temp-check"), closeErr: failure}
		}, wantAbsent: true, wantError: true},
		{name: "hard-link publication", configure: func(fs *fakeV2CheckFileSystem) { fs.linkErr = failure }, wantAbsent: true, wantError: true},
		{name: "cleanup after successful publication", configure: func(fs *fakeV2CheckFileSystem) { fs.removeErr = failure }, wantAbsent: false, wantError: true},
		{name: "primary and cleanup failure", configure: func(fs *fakeV2CheckFileSystem) {
			fs.temp = &fakeV2CheckTempFile{path: filepath.Join(filepath.Dir(path), "temp-check"), writeErr: failure}
			fs.removeErr = errors.New("cleanup failure")
		}, wantAbsent: true, wantError: true, wantErrMsg: "cleanup failure"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fs := &fakeV2CheckFileSystem{}
			tc.configure(fs)
			previous := v2CheckFileOps
			v2CheckFileOps = fs.operations()
			t.Cleanup(func() { v2CheckFileOps = previous })

			err := createV2CheckFile(path, wantContent)
			if (err != nil) != tc.wantError {
				t.Fatalf("createV2CheckFile() error = %v, wantError = %v", err, tc.wantError)
			}
			if tc.wantErrMsg != "" && !strings.Contains(err.Error(), tc.wantErrMsg) {
				t.Fatalf("createV2CheckFile() error = %v, want it to include %q", err, tc.wantErrMsg)
			}
			published, ok := fs.published[path]
			if tc.wantAbsent {
				if ok {
					t.Fatalf("published final file = %q, want final path absent", published)
				}
				return
			}
			if !ok || string(published) != string(wantContent) {
				t.Fatalf("published final file = %q, want byte-identical %q", published, wantContent)
			}
			if fs.removeCall != 1 {
				t.Errorf("cleanup calls = %d, want 1", fs.removeCall)
			}
		})
	}
}

func TestCreateCheckV2_rejectsMalformedRecordLeavesNoFileBehind(t *testing.T) {
	root := t.TempDir()
	index := &V2Index{Checks: map[string]*CheckV2{}}

	fields := NewCheckV2{
		Scope:           CheckScope{Kind: CheckScopeTask, ID: "T-001"},
		Result:          CheckResultClear,
		CheckedBy:       Actor{Role: ActorRoleChecker, Session: ""}, // missing required session
		ExecutedSession: "build-001",
		CheckedAt:       time.Date(2026, 9, 14, 0, 0, 0, 0, time.UTC),
		Body:            "\n\n# Check\n",
	}

	check, err := CreateCheckV2(root, index, fields)
	if err == nil {
		t.Fatal("CreateCheckV2() expected error for missing checked_by.session")
	}
	if check != nil {
		t.Errorf("CreateCheckV2() returned %+v, want nil on validation failure", check)
	}

	if _, statErr := os.Stat(filepath.Join(root, "checks", "C-001.md")); !os.IsNotExist(statErr) {
		t.Errorf("checks/C-001.md exists after rejected record, statErr = %v", statErr)
	}
}

func TestCreateIssueV2_writesNewFileAllocatesFirstID(t *testing.T) {
	root := t.TempDir()
	index := &V2Index{Issues: map[string]*IssueV2{}}

	fields := NewIssueV2{
		Title: "Broken retry loop",
		Type:  IssueTypeDefect,
		Origin: IssueOrigin{
			Kind:  IssueOriginCheck,
			Check: "C-001",
			Actor: Actor{Role: ActorRoleChecker, Session: "sess-1"},
			At:    time.Date(2026, 9, 14, 0, 0, 0, 0, time.UTC),
		},
		Checks: []string{"C-001"},
		Body:   "\n\n# Issue\n\nSummary notes.\n",
	}

	issue, err := CreateIssueV2(root, index, fields)
	if err != nil {
		t.Fatalf("CreateIssueV2() error = %v", err)
	}
	if issue.ID != "I-001" {
		t.Errorf("ID = %q, want I-001", issue.ID)
	}
	if issue.Status != IssueStatusOpen {
		t.Errorf("Status = %q, want open", issue.Status)
	}

	path := filepath.Join(root, "issues", "I-001-broken-retry-loop.md")
	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("issues/I-001-broken-retry-loop.md not written: %v", err)
	}
	if !strings.Contains(string(result), "Summary notes.") {
		t.Error("authored body content not written")
	}

	reparsed, err := DecodeIssueV2(filepath.Join("issues", "I-001-broken-retry-loop.md"), string(result))
	if err != nil {
		t.Fatalf("DecodeIssueV2() after create error = %v", err)
	}
	if reparsed.Type != IssueTypeDefect || reparsed.Status != IssueStatusOpen {
		t.Errorf("reparsed = %+v, want type defect/status open", reparsed)
	}
	if len(reparsed.Checks) != 1 || reparsed.Checks[0] != "C-001" {
		t.Errorf("Checks = %v, want [C-001]", reparsed.Checks)
	}
}

// TestCreateIssueV2_reportOriginOmitsEmptyCheckKey proves a report-sourced
// Issue, which carries no source Check, is written without a blank
// `check: ""` key under source — an absent optional field must stay absent
// in the written bytes, not just decode back to the same value.
func TestCreateIssueV2_reportOriginOmitsEmptyCheckKey(t *testing.T) {
	root := t.TempDir()
	index := &V2Index{Issues: map[string]*IssueV2{}}

	fields := NewIssueV2{
		Title: "Reported directly",
		Type:  IssueTypeDrift,
		Origin: IssueOrigin{
			Kind:  IssueOriginReport,
			Actor: Actor{Role: ActorRolePlanner, Session: "sess-1"},
			At:    time.Date(2026, 9, 14, 0, 0, 0, 0, time.UTC),
		},
		Body: "\n\n# Issue\n",
	}

	issue, err := CreateIssueV2(root, index, fields)
	if err != nil {
		t.Fatalf("CreateIssueV2() error = %v", err)
	}

	raw, err := os.ReadFile(filepath.Join(root, "issues", issue.ID+"-reported-directly.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), "check:") {
		t.Errorf("written record names an empty source.check key, want it omitted entirely:\n%s", raw)
	}
}

func TestCreateIssueV2_allocatesNextIDOverPopulatedIndex(t *testing.T) {
	root := t.TempDir()
	index := &V2Index{Issues: map[string]*IssueV2{
		"I-001": {ID: "I-001"},
		"I-002": {ID: "I-002"},
	}}

	fields := NewIssueV2{
		Title: "Second issue",
		Type:  IssueTypeDrift,
		Origin: IssueOrigin{
			Kind:  IssueOriginReport,
			Actor: Actor{Role: ActorRolePlanner, Session: "sess-1"},
			At:    time.Date(2026, 9, 14, 0, 0, 0, 0, time.UTC),
		},
		Body: "\n\n# Issue\n",
	}

	issue, err := CreateIssueV2(root, index, fields)
	if err != nil {
		t.Fatalf("CreateIssueV2() error = %v", err)
	}
	if issue.ID != "I-003" {
		t.Errorf("ID = %q, want I-003 as next unused id", issue.ID)
	}
}

func TestCreateIssueV2_refusesExistingPathAndLeavesItUntouched(t *testing.T) {
	root := t.TempDir()
	issuesDir := filepath.Join(root, "issues")
	if err := os.MkdirAll(issuesDir, 0755); err != nil {
		t.Fatal(err)
	}
	leftover := "not an issue record"
	if err := os.WriteFile(filepath.Join(issuesDir, "I-001-collides-with-existing-file.md"), []byte(leftover), 0644); err != nil {
		t.Fatal(err)
	}

	index := &V2Index{Issues: map[string]*IssueV2{}}
	fields := NewIssueV2{
		Title: "Collides with existing file",
		Type:  IssueTypeDefect,
		Origin: IssueOrigin{
			Kind:  IssueOriginReport,
			Actor: Actor{Role: ActorRolePlanner, Session: "sess-1"},
			At:    time.Date(2026, 9, 14, 0, 0, 0, 0, time.UTC),
		},
		Body: "\n\n# Issue\n",
	}

	issue, err := CreateIssueV2(root, index, fields)
	if issue != nil {
		t.Errorf("CreateIssueV2() returned %+v, want nil on collision", issue)
	}
	if !errors.Is(err, ErrV2IssueAlreadyExists) {
		t.Fatalf("CreateIssueV2() error = %v, want ErrV2IssueAlreadyExists", err)
	}

	result, err := os.ReadFile(filepath.Join(issuesDir, "I-001-collides-with-existing-file.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(result) != leftover {
		t.Error("existing file at colliding path was modified")
	}
}

func TestCreateIssueV2_rejectsMalformedRecordLeavesNoFileBehind(t *testing.T) {
	root := t.TempDir()
	index := &V2Index{Issues: map[string]*IssueV2{}}

	fields := NewIssueV2{
		Title: "Missing actor session",
		Type:  IssueTypeDefect,
		Origin: IssueOrigin{
			Kind:  IssueOriginReport,
			Actor: Actor{Role: ActorRolePlanner, Session: ""}, // missing required session
			At:    time.Date(2026, 9, 14, 0, 0, 0, 0, time.UTC),
		},
		Body: "\n\n# Issue\n",
	}

	issue, err := CreateIssueV2(root, index, fields)
	if err == nil {
		t.Fatal("CreateIssueV2() expected error for missing source.actor.session")
	}
	if issue != nil {
		t.Errorf("CreateIssueV2() returned %+v, want nil on validation failure", issue)
	}

	if _, statErr := os.Stat(filepath.Join(root, "issues")); !os.IsNotExist(statErr) {
		t.Errorf("issues/ directory exists after rejected record, statErr = %v", statErr)
	}
}

func issueV2FixtureContent() string {
	return `---
id: I-010
title: "Broken retry loop"
type: defect
status: open
source:
  kind: check
  check: C-001
  actor:
    role: checker
    session: sess-1
  at: "2026-09-14T00:00:00Z"
owner:
  team: platform
---

# Issue

## Summary

Authored issue notes that must survive the rewrite.`
}

func TestWriteIssueV2_updatesManagedFieldsPreservesUnknownFieldsAndBody(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "I-010.md")
	content := issueV2FixtureContent()

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	issue, err := DecodeIssueV2(path, content)
	if err != nil {
		t.Fatalf("DecodeIssueV2() error = %v", err)
	}

	issue.Status = IssueStatusResolved
	issue.Severity = "high"
	issue.Tasks = []string{"T-010"}
	issue.Checks = []string{"C-001"}
	issue.DuplicateOf = ""
	resolvedAt := time.Date(2026, 9, 15, 0, 0, 0, 0, time.UTC)
	issue.Resolution = &IssueResolution{
		Disposition: IssueDispositionVerified,
		Check:       "C-001",
		Actor:       Actor{Role: ActorRoleChecker, Session: "sess-1"},
		At:          resolvedAt,
		Reason:      "Recheck confirmed the fix.",
	}

	if err := WriteIssueV2(issue); err != nil {
		t.Fatalf("WriteIssueV2() error = %v", err)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	reparsed, err := DecodeIssueV2(path, string(result))
	if err != nil {
		t.Fatalf("DecodeIssueV2() after write error = %v", err)
	}
	if reparsed.Status != IssueStatusResolved {
		t.Errorf("Status = %q, want resolved", reparsed.Status)
	}
	if reparsed.Severity != "high" {
		t.Errorf("Severity = %q, want high", reparsed.Severity)
	}
	if len(reparsed.Tasks) != 1 || reparsed.Tasks[0] != "T-010" {
		t.Errorf("Tasks = %v, want [T-010]", reparsed.Tasks)
	}
	if reparsed.Resolution == nil || reparsed.Resolution.Disposition != IssueDispositionVerified {
		t.Errorf("Resolution = %+v, want disposition verified", reparsed.Resolution)
	}
	if reparsed.Title != "Broken retry loop" || reparsed.Origin.Check != "C-001" {
		t.Errorf("unrelated fields not preserved: %+v", reparsed)
	}
	if !strings.Contains(string(result), "team: platform") {
		t.Error("unknown nested field not preserved")
	}
	if !strings.Contains(string(result), "Authored issue notes that must survive the rewrite.") {
		t.Error("authored body content not preserved")
	}
}

func TestWriteIssueV2_clearingResolutionOnReopenRemovesKey(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "I-011.md")
	content := `---
id: I-011
title: "Reopened issue"
type: defect
status: resolved
source:
  kind: report
  actor:
    role: planner
    session: sess-1
  at: "2026-09-14T00:00:00Z"
resolution:
  disposition: accepted
  actor:
    role: owner
    session: owner-1
  at: "2026-09-14T01:00:00Z"
  reason: "Accepted known risk."
---

# Issue`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	issue, err := DecodeIssueV2(path, content)
	if err != nil {
		t.Fatalf("DecodeIssueV2() error = %v", err)
	}
	issue.Status = IssueStatusOpen
	issue.Resolution = nil

	if err := WriteIssueV2(issue); err != nil {
		t.Fatalf("WriteIssueV2() error = %v", err)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(result), "resolution:") {
		t.Error("resolution key should be removed when reopening clears it")
	}

	reparsed, err := DecodeIssueV2(path, string(result))
	if err != nil {
		t.Fatalf("DecodeIssueV2() after write error = %v", err)
	}
	if reparsed.Status != IssueStatusOpen || reparsed.Resolution != nil {
		t.Errorf("reparsed = %+v, want open with no resolution", reparsed)
	}
}

// TestWriteIssueV2_acceptedResolutionOmitsProofCheckKey proves an accepted
// resolution, whose entire contract is that it must not name a proof Check,
// is written without a blank `check: ""` key under resolution — naming an
// empty check there would misleadingly look like an unset proof reference
// rather than a disposition that forbids one.
func TestWriteIssueV2_acceptedResolutionOmitsProofCheckKey(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "I-013.md")
	content := `---
id: I-013
title: "Accepted risk"
type: guardrail
status: open
source:
  kind: report
  actor:
    role: planner
    session: sess-1
  at: "2026-09-14T00:00:00Z"
---

# Issue`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	issue, err := DecodeIssueV2(path, content)
	if err != nil {
		t.Fatalf("DecodeIssueV2() error = %v", err)
	}
	issue.Status = IssueStatusResolved
	issue.Resolution = &IssueResolution{
		Disposition: IssueDispositionAccepted,
		Actor:       Actor{Role: ActorRoleOwner, Session: "owner-1"},
		At:          time.Date(2026, 9, 14, 1, 0, 0, 0, time.UTC),
		Reason:      "Accepted known risk.",
	}

	if err := WriteIssueV2(issue); err != nil {
		t.Fatalf("WriteIssueV2() error = %v", err)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(result), "check:") {
		t.Errorf("accepted resolution names a check key, want it omitted entirely:\n%s", result)
	}
}

func TestWriteIssueV2_noOpLeavesBytesAndMtimeUnchanged(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "I-012.md")
	content := issueV2FixtureContent()

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	before, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	beforeBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(10 * time.Millisecond)

	issue, err := DecodeIssueV2(path, content)
	if err != nil {
		t.Fatalf("DecodeIssueV2() error = %v", err)
	}

	if err := WriteIssueV2(issue); err != nil {
		t.Fatalf("WriteIssueV2() error = %v", err)
	}

	after, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	afterBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if !before.ModTime().Equal(after.ModTime()) {
		t.Errorf("ModTime changed on no-op write: before %v, after %v", before.ModTime(), after.ModTime())
	}
	if string(beforeBytes) != string(afterBytes) {
		t.Error("file bytes changed on no-op write")
	}
}

func TestCreateIssueV2ThenWriteIssueV2_roundTripsAsNoOp(t *testing.T) {
	root := t.TempDir()
	index := &V2Index{Issues: map[string]*IssueV2{}}

	fields := NewIssueV2{
		Title: "Freshly created, nothing to patch",
		Type:  IssueTypeOther,
		Origin: IssueOrigin{
			Kind:  IssueOriginReport,
			Actor: Actor{Role: ActorRolePlanner, Session: "sess-1"},
			At:    time.Date(2026, 9, 14, 0, 0, 0, 0, time.UTC),
		},
		Body: "\n\n# Issue\n",
	}

	issue, err := CreateIssueV2(root, index, fields)
	if err != nil {
		t.Fatalf("CreateIssueV2() error = %v", err)
	}

	path := filepath.Join(root, issue.Source.Path)
	before, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	beforeBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(10 * time.Millisecond)

	if err := WriteIssueV2(issue); err != nil {
		t.Fatalf("WriteIssueV2() error = %v", err)
	}

	after, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	afterBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if !before.ModTime().Equal(after.ModTime()) {
		t.Error("ModTime changed writing back a freshly created issue's own empty fields")
	}
	if string(beforeBytes) != string(afterBytes) {
		t.Error("bytes changed writing back a freshly created issue's own empty fields")
	}
}

func TestWriteIssueV2_refusesUnsupportedStatusAndLeavesFileUntouched(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "I-013.md")
	content := issueV2FixtureContent()

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	issue, err := DecodeIssueV2(path, content)
	if err != nil {
		t.Fatalf("DecodeIssueV2() error = %v", err)
	}
	issue.Status = "planned" // Task lifecycle vocabulary, invalid for an Issue

	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	err = WriteIssueV2(issue)
	if err == nil {
		t.Fatal("WriteIssueV2() expected error for planned status")
	}

	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Error("file changed despite refused write")
	}
}

func TestWriteIssueV2_preservesCRLFLineEndings(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "I-014.md")
	content := "---\r\nid: I-014\r\ntitle: \"CRLF issue\"\r\ntype: defect\r\nstatus: open\r\nsource:\r\n  kind: report\r\n  actor:\r\n    role: planner\r\n    session: sess-1\r\n  at: \"2026-09-14T00:00:00Z\"\r\n---\r\n\r\n# Issue\r\n\r\nAuthored notes.\r\n"

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	issue, err := DecodeIssueV2(path, content)
	if err != nil {
		t.Fatalf("DecodeIssueV2() error = %v", err)
	}
	issue.Status = IssueStatusInProgress

	if err := WriteIssueV2(issue); err != nil {
		t.Fatalf("WriteIssueV2() error = %v", err)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Count(string(result), "\n") != strings.Count(string(result), "\r\n") {
		t.Errorf("result did not preserve CRLF line endings throughout: %q", string(result))
	}
	if !strings.Contains(string(result), "Authored notes.") {
		t.Error("authored body content not preserved across CRLF rewrite")
	}

	reparsed, err := DecodeIssueV2(path, string(result))
	if err != nil {
		t.Fatalf("DecodeIssueV2() after write error = %v", err)
	}
	if reparsed.Status != IssueStatusInProgress {
		t.Errorf("Status = %q, want in_progress", reparsed.Status)
	}
}

func issueHistoryFixtureEntry() IssueHistoryEntry {
	return IssueHistoryEntry{
		At:    time.Date(2026, 9, 14, 0, 0, 0, 0, time.UTC),
		Actor: Actor{Role: ActorRoleChecker, Session: "sess-1"},
		Kind:  IssueHistoryObserved,
		Note:  "Found during review.",
		Check: "C-001",
	}
}

func TestWriteIssueHistoryV2_appendsEntryPreservesEarlierEntries(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "I-020.md")
	content := `---
id: I-020
title: "History issue"
type: defect
status: open
source:
  kind: check
  check: C-001
  actor:
    role: checker
    session: sess-1
  at: "2026-09-14T00:00:00Z"
history:
  - at: "2026-09-14T00:00:00Z"
    actor:
      role: checker
      session: sess-1
    kind: observed
    note: "Found during review."
    check: C-001
---

# Issue`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	issue, err := DecodeIssueV2(path, content)
	if err != nil {
		t.Fatalf("DecodeIssueV2() error = %v", err)
	}
	if len(issue.History) != 1 {
		t.Fatalf("History = %v, want 1 recorded entry", issue.History)
	}

	appended := IssueHistoryEntry{
		At:    time.Date(2026, 9, 15, 0, 0, 0, 0, time.UTC),
		Actor: Actor{Role: ActorRoleExecutor, Session: "sess-2"},
		Kind:  IssueHistoryRepairAttempted,
		Note:  "Patched the retry loop.",
	}
	entries := append(append([]IssueHistoryEntry{}, issue.History...), appended)

	if err := WriteIssueHistoryV2(issue, entries); err != nil {
		t.Fatalf("WriteIssueHistoryV2() error = %v", err)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	reparsed, err := DecodeIssueV2(path, string(result))
	if err != nil {
		t.Fatalf("DecodeIssueV2() after write error = %v", err)
	}
	if len(reparsed.History) != 2 {
		t.Fatalf("History = %v, want 2 entries", reparsed.History)
	}
	if reparsed.History[0].Note != "Found during review." || reparsed.History[0].Check != "C-001" {
		t.Errorf("History[0] = %+v, want earlier entry preserved byte-for-byte", reparsed.History[0])
	}
	if reparsed.History[1].Kind != IssueHistoryRepairAttempted || reparsed.History[1].Note != "Patched the retry loop." {
		t.Errorf("History[1] = %+v, want the appended entry", reparsed.History[1])
	}
}

// TestWriteIssueHistoryV2_appendPreservesPriorEntryBytes proves appending a
// new entry does not rewrite an earlier entry's own rendered YAML: a prior
// entry recorded without a note or check must not gain blank `note: ""` or
// `check: ""` keys just because a later write touched the file.
func TestWriteIssueHistoryV2_appendPreservesPriorEntryBytes(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "I-025.md")
	content := `---
id: I-025
title: "History issue"
type: defect
status: open
source:
  kind: report
  actor:
    role: planner
    session: sess-1
  at: "2026-09-14T00:00:00Z"
history:
  - at: "2026-09-14T00:00:00Z"
    actor:
      role: checker
      session: sess-1
    kind: observed
---

# Issue`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	issue, err := DecodeIssueV2(path, content)
	if err != nil {
		t.Fatalf("DecodeIssueV2() error = %v", err)
	}

	appended := IssueHistoryEntry{
		At:    time.Date(2026, 9, 15, 0, 0, 0, 0, time.UTC),
		Actor: Actor{Role: ActorRoleChecker, Session: "sess-1"},
		Kind:  IssueHistoryRechecked,
	}
	entries := append(append([]IssueHistoryEntry{}, issue.History...), appended)

	if err := WriteIssueHistoryV2(issue, entries); err != nil {
		t.Fatalf("WriteIssueHistoryV2() error = %v", err)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(result), `note: ""`) || strings.Contains(string(result), `check: ""`) {
		t.Errorf("write introduced a blank placeholder key on a field the record never declared:\n%s", result)
	}
}

func TestWriteIssueHistoryV2_refusesShorterListLeavesFileUntouched(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "I-021.md")
	content := `---
id: I-021
title: "History issue"
type: defect
status: open
source:
  kind: report
  actor:
    role: planner
    session: sess-1
  at: "2026-09-14T00:00:00Z"
history:
  - at: "2026-09-14T00:00:00Z"
    actor:
      role: checker
      session: sess-1
    kind: observed
    note: "Found during review."
---

# Issue`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	issue, err := DecodeIssueV2(path, content)
	if err != nil {
		t.Fatalf("DecodeIssueV2() error = %v", err)
	}

	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	err = WriteIssueHistoryV2(issue, []IssueHistoryEntry{})
	if !errors.Is(err, ErrV2IssueHistoryNotAppendOnly) {
		t.Fatalf("WriteIssueHistoryV2() error = %v, want ErrV2IssueHistoryNotAppendOnly", err)
	}

	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Error("file changed despite refused shorter history write")
	}
}

func TestWriteIssueHistoryV2_refusesEditedEntryLeavesFileUntouched(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "I-022.md")
	content := `---
id: I-022
title: "History issue"
type: defect
status: open
source:
  kind: report
  actor:
    role: planner
    session: sess-1
  at: "2026-09-14T00:00:00Z"
history:
  - at: "2026-09-14T00:00:00Z"
    actor:
      role: checker
      session: sess-1
    kind: observed
    note: "Found during review."
---

# Issue`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	issue, err := DecodeIssueV2(path, content)
	if err != nil {
		t.Fatalf("DecodeIssueV2() error = %v", err)
	}

	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	edited := issue.History[0]
	edited.Note = "Rewritten note."
	err = WriteIssueHistoryV2(issue, []IssueHistoryEntry{edited})
	if !errors.Is(err, ErrV2IssueHistoryNotAppendOnly) {
		t.Fatalf("WriteIssueHistoryV2() error = %v, want ErrV2IssueHistoryNotAppendOnly", err)
	}

	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Error("file changed despite refused edited-entry history write")
	}
}

func TestWriteIssueHistoryV2_refusesReorderedEntriesLeavesFileUntouched(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "I-023.md")
	content := `---
id: I-023
title: "History issue"
type: defect
status: open
source:
  kind: report
  actor:
    role: planner
    session: sess-1
  at: "2026-09-14T00:00:00Z"
history:
  - at: "2026-09-14T00:00:00Z"
    actor:
      role: checker
      session: sess-1
    kind: observed
    note: "First."
  - at: "2026-09-15T00:00:00Z"
    actor:
      role: checker
      session: sess-1
    kind: rechecked
    note: "Second."
---

# Issue`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	issue, err := DecodeIssueV2(path, content)
	if err != nil {
		t.Fatalf("DecodeIssueV2() error = %v", err)
	}
	if len(issue.History) != 2 {
		t.Fatalf("History = %v, want 2 recorded entries", issue.History)
	}

	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	reordered := []IssueHistoryEntry{issue.History[1], issue.History[0]}
	err = WriteIssueHistoryV2(issue, reordered)
	if !errors.Is(err, ErrV2IssueHistoryNotAppendOnly) {
		t.Fatalf("WriteIssueHistoryV2() error = %v, want ErrV2IssueHistoryNotAppendOnly", err)
	}

	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Error("file changed despite refused reordered history write")
	}
}

func TestWriteIssueHistoryV2_noOpLeavesBytesAndMtimeUnchanged(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "I-024.md")
	content := `---
id: I-024
title: "History issue"
type: defect
status: open
source:
  kind: report
  actor:
    role: planner
    session: sess-1
  at: "2026-09-14T00:00:00Z"
history:
  - at: "2026-09-14T00:00:00Z"
    actor:
      role: checker
      session: sess-1
    kind: observed
    note: "Found during review."
    check: C-001
---

# Issue`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	before, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	beforeBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(10 * time.Millisecond)

	issue, err := DecodeIssueV2(path, content)
	if err != nil {
		t.Fatalf("DecodeIssueV2() error = %v", err)
	}

	if err := WriteIssueHistoryV2(issue, issue.History); err != nil {
		t.Fatalf("WriteIssueHistoryV2() error = %v", err)
	}

	after, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	afterBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if !before.ModTime().Equal(after.ModTime()) {
		t.Errorf("ModTime changed on no-op history write: before %v, after %v", before.ModTime(), after.ModTime())
	}
	if string(beforeBytes) != string(afterBytes) {
		t.Error("file bytes changed on no-op history write")
	}
}

// TestWriteIssueHistoryV2_emptyOnHistoryFreeIssueIsNoOp proves writing zero
// entries to an Issue that already has no history key stays a true no-op:
// bytes and modification time unchanged, and no `history: []` key appears.
func TestWriteIssueHistoryV2_emptyOnHistoryFreeIssueIsNoOp(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "I-026.md")
	content := `---
id: I-026
title: "No history yet"
type: defect
status: open
source:
  kind: report
  actor:
    role: planner
    session: sess-1
  at: "2026-09-14T00:00:00Z"
---

# Issue`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	before, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	beforeBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(10 * time.Millisecond)

	issue, err := DecodeIssueV2(path, content)
	if err != nil {
		t.Fatalf("DecodeIssueV2() error = %v", err)
	}
	if len(issue.History) != 0 {
		t.Fatalf("History = %v, want none recorded", issue.History)
	}

	if err := WriteIssueHistoryV2(issue, nil); err != nil {
		t.Fatalf("WriteIssueHistoryV2() error = %v", err)
	}

	after, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	afterBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if !before.ModTime().Equal(after.ModTime()) {
		t.Errorf("ModTime changed on empty-history no-op write: before %v, after %v", before.ModTime(), after.ModTime())
	}
	if string(beforeBytes) != string(afterBytes) {
		t.Errorf("file bytes changed on empty-history no-op write:\nbefore: %s\nafter: %s", beforeBytes, afterBytes)
	}
}

func TestWriteObjectiveEvidenceV2_setsSubBlocksPreservesUnknownFieldsAndBody(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "O-020.md")
	content := `---
id: O-020
title: "No evidence yet"
status: in_progress
depends_on: [O-001]
release: R-001
---

# Objective

Authored objective notes.`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	objective, err := DecodeObjectiveV2(path, content)
	if err != nil {
		t.Fatalf("DecodeObjectiveV2() error = %v", err)
	}
	if objective.Evidence != nil {
		t.Fatalf("Evidence = %+v, want nil before write", objective.Evidence)
	}

	assessedAt := time.Date(2026, 9, 14, 0, 0, 0, 0, time.UTC)
	objective.Evidence = &Evidence{
		LastCheck: "C-001",
		Freshness: &Freshness{
			State:      FreshnessCurrent,
			Check:      "C-001",
			AssessedBy: Actor{Role: ActorRoleChecker, Session: "sess-1"},
			AssessedAt: assessedAt,
			Basis:      "Objective integration Check is current.",
		},
	}

	if err := WriteObjectiveEvidenceV2(objective); err != nil {
		t.Fatalf("WriteObjectiveEvidenceV2() error = %v", err)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	reparsed, err := DecodeObjectiveV2(path, string(result))
	if err != nil {
		t.Fatalf("DecodeObjectiveV2() after write error = %v", err)
	}
	if reparsed.Evidence == nil || reparsed.Evidence.LastCheck != "C-001" {
		t.Errorf("Evidence = %+v, want LastCheck C-001", reparsed.Evidence)
	}
	if reparsed.Evidence.Freshness == nil || reparsed.Evidence.Freshness.State != FreshnessCurrent {
		t.Errorf("Freshness = %+v, want state current", reparsed.Evidence.Freshness)
	}
	if len(reparsed.DependsOn) != 1 || reparsed.DependsOn[0] != "O-001" {
		t.Errorf("DependsOn = %v, want [O-001] preserved", reparsed.DependsOn)
	}
	if reparsed.Release != "R-001" {
		t.Errorf("Release = %q, want R-001 preserved", reparsed.Release)
	}
	if !strings.Contains(string(result), "Authored objective notes.") {
		t.Error("authored body content not preserved")
	}
}

func TestWriteObjectiveEvidenceV2_noOpLeavesBytesAndMtimeUnchanged(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "O-021.md")
	content := `---
id: O-021
title: "No-op objective evidence"
status: in_progress
last_check: C-001
freshness:
  state: current
  check: C-001
  assessed_by:
    role: checker
    session: sess-1
  assessed_at: "2026-09-14T00:00:00Z"
  basis: "Objective integration Check is current."
---

# Objective`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	before, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	beforeBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(10 * time.Millisecond)

	objective, err := DecodeObjectiveV2(path, content)
	if err != nil {
		t.Fatalf("DecodeObjectiveV2() error = %v", err)
	}

	if err := WriteObjectiveEvidenceV2(objective); err != nil {
		t.Fatalf("WriteObjectiveEvidenceV2() error = %v", err)
	}

	after, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	afterBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if !before.ModTime().Equal(after.ModTime()) {
		t.Errorf("ModTime changed on no-op write: before %v, after %v", before.ModTime(), after.ModTime())
	}
	if string(beforeBytes) != string(afterBytes) {
		t.Error("file bytes changed on no-op write")
	}
}

func TestWriteDefectStatus_removesStageWhenDone(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "D002.md")
	content := `---
id: v1/D002
release: v1
status: in_progress
stage: build
severity: medium
title: "Bug"
---

# Body`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	fi, _ := os.Stat(path)
	defect := &Defect{ID: "v1/D002", Release: "v1", Status: DefectResolved, Severity: SeverityMedium, Title: "Bug"}

	if err := WriteDefectStatus(path, defect, fi.ModTime()); err != nil {
		t.Fatalf("WriteDefectStatus() error = %v", err)
	}

	result, _ := os.ReadFile(path)
	if strings.Contains(string(result), "stage:") {
		t.Error("stage field should be removed when defect status is resolved")
	}
}

// routerV2FixtureContent reproduces the shipped V2 router's document shape:
// prose, headings and a second fenced block around the anchor, and the four
// anchor keys with a long quoted next_action. The anchor's key set is closed
// — ReadStateV2 decodes it with KnownFields(true) — so what a selection
// write must preserve is state, next_action, key order, and every byte of
// the document outside the anchor.
func routerV2FixtureContent() string {
	return `# Agent State Machine

This file routes the agent. The active skill is the canonical workflow source.

## Read order

1. This file
2. The matching skill for ` + "`state`" + `

## Current state

` + "```" + `yaml
state: design
objective: none
task: none
next_action: "Turn a rough idea — a single sentence is a valid starting point, no prepared requirements document needed — into .savepoint/Idea.md through a short back-and-forth with the owner."
` + "```" + `

## State Meanings

- ` + "`idea`" + `: intent and boundary are being defined.

` + "```" + `yaml
project_note: a second fenced block the writer must not touch
` + "```" + `
`
}

func writeRouterV2Fixture(t *testing.T, content string) (root, path string, mtime time.Time) {
	t.Helper()
	root = t.TempDir()
	path = filepath.Join(root, "router.md")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	return root, path, info.ModTime()
}

func readFileString(t *testing.T, path string) string {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(raw)
}

func TestWriteRouterStateV2_setsSelectionAndPreservesEveryOtherByte(t *testing.T) {
	content := routerV2FixtureContent()
	root, path, mtime := writeRouterV2Fixture(t, content)

	if err := WriteRouterStateV2(root, RouterSelectionV2{Objective: "O-001", Task: "T-005"}, mtime); err != nil {
		t.Fatalf("WriteRouterStateV2() error = %v", err)
	}

	want := strings.Replace(content, "objective: none", "objective: O-001", 1)
	want = strings.Replace(want, "task: none", "task: T-005", 1)

	if got := readFileString(t, path); got != want {
		t.Errorf("file content changed beyond the selection keys.\n got:\n%s\nwant:\n%s", got, want)
	}
}

func TestWriteRouterStateV2_setsReleaseContextWithoutChangingRouterProse(t *testing.T) {
	content := routerV2FixtureContent()
	root, path, mtime := writeRouterV2Fixture(t, content)
	before, err := NewRouterReader().ReadStateV2(content)
	if err != nil {
		t.Fatalf("ReadStateV2() before release write error = %v", err)
	}

	if err := WriteRouterStateV2(root, RouterSelectionV2{Release: "R-001", Objective: "O-001", Task: "T-005"}, mtime); err != nil {
		t.Fatalf("WriteRouterStateV2() error = %v", err)
	}

	state, err := NewRouterReader().ReadStateV2(readFileString(t, path))
	if err != nil {
		t.Fatalf("ReadStateV2() after release write error = %v", err)
	}
	if state.Release != "R-001" || state.Objective != "O-001" || state.Task != "T-005" {
		t.Fatalf("selections = release %q objective %q task %q, want R-001/O-001/T-005", state.Release, state.Objective, state.Task)
	}
	if state.State != before.State || state.NextAction != before.NextAction {
		t.Fatalf("router lifecycle/prose changed: before state %q next_action %q, after state %q next_action %q", before.State, before.NextAction, state.State, state.NextAction)
	}
	got := readFileString(t, path)
	if !strings.Contains(got, "project_note: a second fenced block the writer must not touch") || !strings.Contains(got, "This file routes the agent.") {
		t.Error("release selection write did not preserve surrounding document bytes")
	}
}

func TestWriteRouterStateV2_repeatingReleaseSelectionLeavesBytesAndMtimeUnchanged(t *testing.T) {
	content := strings.NewReplacer(
		"objective: none", "objective: O-001",
		"task: none", "task: T-005",
	).Replace(routerV2FixtureContent())
	content = strings.Replace(content, "state: design", "state: task\nrelease: R-001", 1)
	root, path, mtime := writeRouterV2Fixture(t, content)

	time.Sleep(10 * time.Millisecond)
	if err := WriteRouterStateV2(root, RouterSelectionV2{Release: "R-001", Objective: "O-001", Task: "T-005"}, mtime); err != nil {
		t.Fatalf("WriteRouterStateV2() error = %v", err)
	}
	after, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if !after.ModTime().Equal(mtime) || readFileString(t, path) != content {
		t.Fatal("repeating unchanged Release selection changed bytes or mtime")
	}
}

func TestWriteRouterStateV2_clearingSelectionWritesTheNoneSentinel(t *testing.T) {
	content := strings.NewReplacer(
		"objective: none", "objective: O-001",
		"task: none", "task: T-005",
	).Replace(routerV2FixtureContent())
	root, path, mtime := writeRouterV2Fixture(t, content)

	if err := WriteRouterStateV2(root, RouterSelectionV2{}, mtime); err != nil {
		t.Fatalf("WriteRouterStateV2() error = %v", err)
	}

	got := readFileString(t, path)
	if !strings.Contains(got, "objective: none") || !strings.Contains(got, "task: none") {
		t.Errorf("cleared selection not written as the none sentinel:\n%s", got)
	}

	state, err := NewRouterReader().ReadStateV2(got)
	if err != nil {
		t.Fatalf("ReadStateV2() error = %v", err)
	}
	if state.Objective != "" || state.Task != "" {
		t.Errorf("cleared selection read back as objective %q task %q", state.Objective, state.Task)
	}
}

func TestWriteRouterStateV2_doesNotAddSentinelKeysTheDocumentLacks(t *testing.T) {
	content := "## Current state\n\n```yaml\nstate: idea\n```\n"
	root, path, mtime := writeRouterV2Fixture(t, content)

	if err := WriteRouterStateV2(root, RouterSelectionV2{}, mtime); err != nil {
		t.Fatalf("WriteRouterStateV2() error = %v", err)
	}

	if got := readFileString(t, path); got != content {
		t.Errorf("clearing an absent selection changed the file:\n%s", got)
	}
}

func TestWriteRouterStateV2_addsSelectionKeysToRecordARealSelection(t *testing.T) {
	content := "## Current state\n\n```yaml\nstate: task\n```\n"
	root, path, mtime := writeRouterV2Fixture(t, content)

	if err := WriteRouterStateV2(root, RouterSelectionV2{Objective: "O-007", Task: "T-042"}, mtime); err != nil {
		t.Fatalf("WriteRouterStateV2() error = %v", err)
	}

	state, err := NewRouterReader().ReadStateV2(readFileString(t, path))
	if err != nil {
		t.Fatalf("ReadStateV2() error = %v", err)
	}
	if state.Objective != "O-007" || state.Task != "T-042" {
		t.Errorf("selection read back as objective %q task %q, want O-007/T-042", state.Objective, state.Task)
	}
	if state.State != RouterPhaseTask {
		t.Errorf("state changed to %q", state.State)
	}
}

func TestWriteRouterStateV2_refusesMalformedSelectionAndLeavesFileUntouched(t *testing.T) {
	cases := []struct {
		name      string
		selection RouterSelectionV2
		wantErr   error
	}{
		{"objective wrong family", RouterSelectionV2{Objective: "T-001"}, ErrV2InvalidID},
		{"release wrong family", RouterSelectionV2{Release: "O-001"}, ErrV2InvalidID},
		{"release too few digits", RouterSelectionV2{Release: "R01"}, ErrV2InvalidID},
		{"objective too few digits", RouterSelectionV2{Objective: "O1"}, ErrV2InvalidID},
		{"objective with trailing slug", RouterSelectionV2{Objective: "O-001-recovery"}, ErrV2InvalidID},
		{"two objectives", RouterSelectionV2{Objective: "O-001 O-002"}, ErrV2InvalidID},
		{"task wrong family", RouterSelectionV2{Objective: "O-001", Task: "C-001"}, ErrV2InvalidID},
		{"task too few digits", RouterSelectionV2{Objective: "O-001", Task: "T4"}, ErrV2InvalidID},
		{"task without objective", RouterSelectionV2{Task: "T-005"}, ErrV2InvalidOwnership},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			content := routerV2FixtureContent()
			root, path, mtime := writeRouterV2Fixture(t, content)

			err := WriteRouterStateV2(root, tc.selection, mtime)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("WriteRouterStateV2() error = %v, want %v", err, tc.wantErr)
			}
			if got := readFileString(t, path); got != content {
				t.Error("file changed despite a refused selection")
			}
		})
	}
}

func TestWriteRouterStateV2_noOpLeavesBytesAndMtimeUnchanged(t *testing.T) {
	content := strings.NewReplacer(
		"objective: none", "objective: O-001",
		"task: none", "task: T-005",
	).Replace(routerV2FixtureContent())
	root, path, mtime := writeRouterV2Fixture(t, content)

	time.Sleep(10 * time.Millisecond)

	if err := WriteRouterStateV2(root, RouterSelectionV2{Objective: "O-001", Task: "T-005"}, mtime); err != nil {
		t.Fatalf("WriteRouterStateV2() error = %v", err)
	}

	after, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if !after.ModTime().Equal(mtime) {
		t.Errorf("ModTime changed on no-op write: before %v, after %v", mtime, after.ModTime())
	}
	if got := readFileString(t, path); got != content {
		t.Error("file bytes changed on no-op write")
	}
}

func TestWriteRouterStateV2_refusesStaleMtimeAndLeavesFileUntouched(t *testing.T) {
	content := routerV2FixtureContent()
	root, path, _ := writeRouterV2Fixture(t, content)

	err := WriteRouterStateV2(root, RouterSelectionV2{Objective: "O-001"}, time.Now().Add(-time.Hour))
	if !errors.Is(err, ErrMtimeConflict) {
		t.Fatalf("WriteRouterStateV2() error = %v, want ErrMtimeConflict", err)
	}
	if got := readFileString(t, path); got != content {
		t.Error("file changed despite an mtime conflict")
	}
}

func TestWriteRouterStateV2_refusesUnreadableDocuments(t *testing.T) {
	cases := []struct {
		name    string
		content string
		wantErr string
	}{
		{"no anchor", "# Agent State Machine\n\nNo current state here.\n", "no Current state block found"},
		{"no yaml fence", "## Current state\n\nstate: idea\n", "no yaml code block found"},
		{"unterminated fence", "## Current state\n\n```yaml\nstate: idea\n", "no closing code block found"},
		{"anchor is not a mapping", "## Current state\n\n```yaml\n- state: idea\n```\n", "not a mapping"},
		{"anchor is not yaml", "## Current state\n\n```yaml\nstate: [idea\n```\n", "router state"},
		// The anchor's key set is closed on the way in, so a key the reader
		// would reject is refused on the way out too rather than being
		// written into a document that would then fail to load.
		{"anchor carries an unknown key", "## Current state\n\n```yaml\nstate: idea\nproject_note: extra\n```\n", "project_note"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root, path, mtime := writeRouterV2Fixture(t, tc.content)

			err := WriteRouterStateV2(root, RouterSelectionV2{Objective: "O-001"}, mtime)
			if err == nil {
				t.Fatal("WriteRouterStateV2() expected an error")
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("WriteRouterStateV2() error = %v, want it to mention %q", err, tc.wantErr)
			}
			if got := readFileString(t, path); got != tc.content {
				t.Error("file changed despite a refused write")
			}
		})
	}
}

func TestWriteRouterStateV2_refusesMissingFile(t *testing.T) {
	root := t.TempDir()

	err := WriteRouterStateV2(root, RouterSelectionV2{Objective: "O-001"}, time.Now())
	if err == nil {
		t.Fatal("WriteRouterStateV2() expected an error for a missing router.md")
	}
	if !strings.Contains(err.Error(), "router.md") {
		t.Errorf("WriteRouterStateV2() error = %v, want it to name router.md", err)
	}
}

func TestWriteRouterStateV2_refusesARouterItCouldNotReadBack(t *testing.T) {
	content := strings.Replace(routerV2FixtureContent(), "state: design", "state: audit-pending", 1)
	root, path, mtime := writeRouterV2Fixture(t, content)

	err := WriteRouterStateV2(root, RouterSelectionV2{Objective: "O-001"}, mtime)
	if !errors.Is(err, ErrV2InvalidLifecycle) {
		t.Fatalf("WriteRouterStateV2() error = %v, want ErrV2InvalidLifecycle", err)
	}
	if got := readFileString(t, path); got != content {
		t.Error("file changed despite content that could not be read back")
	}
}

func TestWriteRouterStateV2_roundTripsThroughReadStateV2(t *testing.T) {
	content := routerV2FixtureContent()
	root, path, mtime := writeRouterV2Fixture(t, content)

	if err := WriteRouterStateV2(root, RouterSelectionV2{Objective: "O-012", Task: "T-034"}, mtime); err != nil {
		t.Fatalf("WriteRouterStateV2() error = %v", err)
	}

	state, err := NewRouterReader().ReadStateV2(readFileString(t, path))
	if err != nil {
		t.Fatalf("ReadStateV2() error = %v", err)
	}
	if state.Objective != "O-012" || state.Task != "T-034" {
		t.Errorf("round trip gave objective %q task %q, want O-012/T-034", state.Objective, state.Task)
	}
	if state.State != RouterPhaseDesign {
		t.Errorf("state changed to %q, want design", state.State)
	}
	if !strings.Contains(state.NextAction, "Turn a rough idea") {
		t.Errorf("next_action changed to %q", state.NextAction)
	}
}

func TestWriteRouterStateV2_preservesCRLFLineEndings(t *testing.T) {
	content := strings.ReplaceAll(routerV2FixtureContent(), "\n", "\r\n")
	root, path, mtime := writeRouterV2Fixture(t, content)

	if err := WriteRouterStateV2(root, RouterSelectionV2{Objective: "O-001"}, mtime); err != nil {
		t.Fatalf("WriteRouterStateV2() error = %v", err)
	}

	got := readFileString(t, path)
	if strings.Contains(strings.ReplaceAll(got, "\r\n", ""), "\n") {
		t.Error("CRLF router gained bare LF line endings")
	}
	want := strings.Replace(content, "objective: none", "objective: O-001", 1)
	if got != want {
		t.Errorf("CRLF router changed beyond the selection key:\n%q", got)
	}
}
