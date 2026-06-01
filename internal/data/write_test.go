package data

import (
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
id: E01/T001
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

func TestWriteTaskStatus_removesStageWhenStageEmpty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "task.md")
	content := `---
id: E01/T002
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
id: E01/T002
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
id: E01/T003
status: in_progress
stage: build
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
	if strings.Contains(string(result), "stage:") {
		t.Error("stage field should be removed when status is planned")
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

func TestWriteTaskStatus_addsStageWhenStagePresent(t *testing.T) {
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
id: E01/T011
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
		ID:     "E01/T011",
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

func TestWriteTaskStatus_rejectsImplementationStageForInProgress(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "task.md")
	content := `---
id: E01/T012
status: planned
---`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	fi, _ := os.Stat(path)
	task := &Task{
		ID:     "E01/T012",
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
id: E19/T011
status: planned
objective: "Invalid complexity"
---`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	fi, _ := os.Stat(path)
	task := &Task{
		ID:               "E19/T011",
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
id: E19/T012
status: planned
objective: "Invalid complexity"
---`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	fi, _ := os.Stat(path)
	task := &Task{
		ID:               "E19/T012",
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
id: E19/T001
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
		ID:     "E19/T001",
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
	content := "---\nid: E19/T013\nstatus: planned\ncomplexity_tier: small\ncomplexity_reason: \"" + longReason + "\"\nobjective: \"Complexity repair\"\n---\n\n# Body"

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
reference: E01/T001
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
	if !strings.Contains(string(result), "reference: E01/T001") {
		t.Error("unrelated frontmatter field not preserved")
	}
	if !strings.Contains(string(result), "Keep this text.") {
		t.Error("body content not preserved")
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
