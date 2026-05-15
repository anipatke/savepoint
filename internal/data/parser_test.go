package data

import (
	"strings"
	"testing"
)

func TestParseFrontmatter(t *testing.T) {
	p := NewParser()
	content := `---
id: E01/T001
status: done
objective: "Test objective"
depends_on: []
---

body content here`

	result, err := p.ParseFrontmatter(content)
	if err != nil {
		t.Fatalf("ParseFrontmatter() error = %v", err)
	}

	if result["id"] != "E01/T001" {
		t.Errorf("ParseFrontmatter() id = %v, want E01/T001", result["id"])
	}
	if result["objective"] != "Test objective" {
		t.Errorf("ParseFrontmatter() objective = %v, want Test objective", result["objective"])
	}
}

func TestParseFrontmatterMissing(t *testing.T) {
	p := NewParser()
	content := `# No frontmatter here`

	_, err := p.ParseFrontmatter(content)
	if err == nil {
		t.Error("ParseFrontmatter() expected error for missing frontmatter")
	}
}

func TestParseFrontmatterMalformedYAML(t *testing.T) {
	p := NewParser()
	content := `---
id: [broken
---

# Task description`

	_, err := p.ParseFrontmatter(content)
	if err == nil {
		t.Fatal("ParseFrontmatter() expected malformed YAML error")
	}
}

func TestParseTaskFile(t *testing.T) {
	p := NewParser()
	content := `---
id: E02/T001
status: in_progress
stage: test
objective: "Define Task struct"
description: "Build the task model"
priority: high
points: 3
tags: [data, parser]
acceptance:
  - parses metadata
notes: "Keep it small"
depends_on: [E01/T003]
---

# Task description`

	task, err := p.ParseTaskFile("test.md", content)
	if err != nil {
		t.Fatalf("ParseTaskFile() error = %v", err)
	}

	if task.ID != "E02/T001" {
		t.Errorf("Task.ID = %v, want E02/T001", task.ID)
	}
	if task.Title != "Define Task struct" {
		t.Errorf("Task.Title = %v, want Define Task struct", task.Title)
	}
	if task.Epic != "E02" {
		t.Errorf("Task.Epic = %v, want E02", task.Epic)
	}
	if task.Release != "v1" {
		t.Errorf("Task.Release = %v, want v1", task.Release)
	}
	if task.Column != ColumnInProgress {
		t.Errorf("Task.Column = %v, want %v", task.Column, ColumnInProgress)
	}
	if task.Stage != StageTest {
		t.Errorf("Task.Stage = %v, want %v", task.Stage, StageTest)
	}
	if len(task.DependsOn) != 1 || task.DependsOn[0] != "E01/T003" {
		t.Errorf("Task.DependsOn = %v, want [E01/T003]", task.DependsOn)
	}
	if task.Priority != "high" || task.Points != 3 {
		t.Errorf("Task priority/points = %v/%v, want high/3", task.Priority, task.Points)
	}
	if len(task.Tags) != 2 || task.Tags[0] != "data" || task.Tags[1] != "parser" {
		t.Errorf("Task.Tags = %v, want [data parser]", task.Tags)
	}
	if len(task.Acceptance) != 1 || task.Acceptance[0] != "parses metadata" {
		t.Errorf("Task.Acceptance = %v, want [parses metadata]", task.Acceptance)
	}
	if task.Notes != "Keep it small" {
		t.Errorf("Task.Notes = %v, want Keep it small", task.Notes)
	}
}

func TestParseTaskFile_normalizesLegacyTodoStatusToPlanned(t *testing.T) {
	p := NewParser()
	content := `---
id: E06/T001
status: todo
objective: "Style the board"
---

# Task`

	task, err := p.ParseTaskFile("test.md", content)
	if err != nil {
		t.Fatalf("ParseTaskFile() error = %v", err)
	}
	if task.Column != ColumnPlanned {
		t.Errorf("Task.Column = %v, want %v", task.Column, ColumnPlanned)
	}
}

func TestParseTaskFile_rejectsUnknownStatus(t *testing.T) {
	p := NewParser()
	content := `---
id: E06/T001
status: review
objective: "Style the board"
---

# Task`

	_, err := p.ParseTaskFile("test.md", content)
	if err == nil {
		t.Fatal("ParseTaskFile() expected unknown status error")
	}
}

func TestParseTaskFile_allowsLegacyPhaseOutsideInProgress(t *testing.T) {
	p := NewParser()
	content := `---
id: E06/T001
status: planned
phase: build
objective: "Style the board"
---

# Task`

	task, err := p.ParseTaskFile("test.md", content)
	if err != nil {
		t.Fatalf("ParseTaskFile() error = %v, want no error for legacy phase field", err)
	}
	if task.Stage != "" {
		t.Fatalf("Task.Stage = %q, want empty for non-in-progress task", task.Stage)
	}
}

func TestParseTaskFile_clearsLegacyImplementationPhaseOutsideInProgress(t *testing.T) {
	p := NewParser()
	content := `---
id: E06/T001
status: done
phase: implementation
objective: "Style the board"
---

# Task`

	task, err := p.ParseTaskFile("test.md", content)
	if err != nil {
		t.Fatalf("ParseTaskFile() error = %v, want no error for stale implementation phase", err)
	}
	if task.Stage != "" {
		t.Fatalf("Task.Stage = %q, want empty for non-in-progress task", task.Stage)
	}
}

func TestParseTaskFile_rejectsMissingStageForInProgress(t *testing.T) {
	p := NewParser()
	content := `---
id: E06/T001
status: in_progress
objective: "Style the board"
---

# Task`

	_, err := p.ParseTaskFile("test.md", content)
	if err == nil {
		t.Fatal("ParseTaskFile() expected missing stage error")
	}
	if !strings.Contains(err.Error(), "stage is required") {
		t.Fatalf("ParseTaskFile() error = %v, want missing stage message", err)
	}
}

func TestParseTaskFile_prefersCanonicalStageOverLegacyPhase(t *testing.T) {
	p := NewParser()
	content := `---
id: E06/T001
status: in_progress
stage: build
phase: test
objective: "Style the board"
---

# Task`

	task, err := p.ParseTaskFile("test.md", content)
	if err != nil {
		t.Fatalf("ParseTaskFile() error = %v", err)
	}
	if task.Stage != StageBuild {
		t.Fatalf("Task.Stage = %q, want build from canonical stage", task.Stage)
	}
}

func TestParseTaskFile_readsLegacyPhaseForInProgress(t *testing.T) {
	p := NewParser()
	content := `---
id: E06/T001
status: in_progress
phase: test
objective: "Style the board"
---

# Task`

	task, err := p.ParseTaskFile("test.md", content)
	if err != nil {
		t.Fatalf("ParseTaskFile() error = %v", err)
	}
	if task.Stage != StageTest {
		t.Fatalf("Task.Stage = %q, want test from legacy phase", task.Stage)
	}
}

func TestParseTaskFile_rejectsInvalidLegacyPhaseForInProgress(t *testing.T) {
	p := NewParser()
	content := `---
id: E06/T001
status: in_progress
phase: done
objective: "Style the board"
---

# Task`

	_, err := p.ParseTaskFile("test.md", content)
	if err == nil {
		t.Fatal("ParseTaskFile() expected invalid stage error")
	}
	if !strings.Contains(err.Error(), `invalid stage "done"`) {
		t.Fatalf("ParseTaskFile() error = %v, want invalid stage done message", err)
	}
}

func TestParseTaskFile_rejectsLegacyImplementationPhaseForInProgress(t *testing.T) {
	p := NewParser()
	content := `---
id: E06/T001
status: in_progress
phase: implementation
objective: "Style the board"
---

# Task`

	_, err := p.ParseTaskFile("test.md", content)
	if err == nil {
		t.Fatal("ParseTaskFile() expected invalid legacy phase error")
	}
	if !strings.Contains(err.Error(), `invalid stage "implementation"`) {
		t.Fatalf("ParseTaskFile() error = %v, want invalid implementation phase message", err)
	}
}

func TestParseTaskFile_rejectsImplementationStageForInProgress(t *testing.T) {
	p := NewParser()
	content := `---
id: E06/T001
status: in_progress
stage: implementation
objective: "Style the board"
---

# Task`

	_, err := p.ParseTaskFile("test.md", content)
	if err == nil {
		t.Fatal("ParseTaskFile() expected invalid stage error")
	}
	if !strings.Contains(err.Error(), `invalid stage "implementation"`) {
		t.Fatalf("ParseTaskFile() error = %v, want invalid implementation stage message", err)
	}
}

func TestParseTaskFile_clearsStageOutsideInProgress(t *testing.T) {
	p := NewParser()
	content := `---
id: E06/T001
status: planned
stage: build
objective: "Style the board"
---

# Task`

	task, err := p.ParseTaskFile("test.md", content)
	if err != nil {
		t.Fatalf("ParseTaskFile() error = %v, want no error for stale stage field", err)
	}
	if task.Stage != "" {
		t.Fatalf("Task.Stage = %q, want empty for non-in-progress task", task.Stage)
	}
}

func TestParseTaskFile_clearsLegacyImplementationStageOutsideInProgress(t *testing.T) {
	p := NewParser()
	content := `---
id: E06/T001
status: done
stage: implementation
objective: "Style the board"
---

# Task`

	task, err := p.ParseTaskFile("test.md", content)
	if err != nil {
		t.Fatalf("ParseTaskFile() error = %v, want no error for stale implementation stage", err)
	}
	if task.Stage != "" {
		t.Fatalf("Task.Stage = %q, want empty for non-in-progress task", task.Stage)
	}
}

func TestParseTaskFile_extractsMarkdownAcceptanceAndChecklist(t *testing.T) {
	p := NewParser()
	content := `---
id: E06/T001
status: planned
objective: "Style the board"
---

# Task

## Acceptance Criteria

- First criterion.
- Second criterion.

## Implementation Plan

- [ ] First checklist item.
- [x] Second checklist item.

## Context Log

Notes here.`

	task, err := p.ParseTaskFile("test.md", content)
	if err != nil {
		t.Fatalf("ParseTaskFile() error = %v", err)
	}
	if len(task.Acceptance) != 2 || task.Acceptance[0] != "First criterion." || task.Acceptance[1] != "Second criterion." {
		t.Errorf("Task.Acceptance = %v, want markdown criteria", task.Acceptance)
	}
	if len(task.Checklist) != 2 {
		t.Fatalf("Task.Checklist len = %d, want 2", len(task.Checklist))
	}
	if task.Checklist[0].Text != "First checklist item." || task.Checklist[0].Done {
		t.Errorf("Task.Checklist[0] = %+v, want {Text:\"First checklist item.\", Done:false}", task.Checklist[0])
	}
	if task.Checklist[1].Text != "Second checklist item." || !task.Checklist[1].Done {
		t.Errorf("Task.Checklist[1] = %+v, want {Text:\"Second checklist item.\", Done:true}", task.Checklist[1])
	}
}

func TestParseTaskFile_complexityRoundTrip(t *testing.T) {
	p := NewParser()
	content := `---
id: E19/T001
status: planned
objective: "Complexity test"
complexity_tier: high
complexity_reason: "Requires coordinated changes across multiple packages."
depends_on: []
---

# Task

## Acceptance Criteria

- Works.`

	task, err := p.ParseTaskFile("test.md", content)
	if err != nil {
		t.Fatalf("ParseTaskFile() error = %v", err)
	}
	if task.ComplexityTier != ComplexityHigh {
		t.Errorf("ComplexityTier = %q, want high", task.ComplexityTier)
	}
	if task.ComplexityReason != "Requires coordinated changes across multiple packages." {
		t.Errorf("ComplexityReason = %q, want reason", task.ComplexityReason)
	}
}

func TestParseTaskFile_complexityAbsentIsCompatible(t *testing.T) {
	p := NewParser()
	content := `---
id: E19/T002
status: planned
objective: "No complexity"
depends_on: []
---

# Task

## Acceptance Criteria

- Works.`

	task, err := p.ParseTaskFile("test.md", content)
	if err != nil {
		t.Fatalf("ParseTaskFile() error = %v", err)
	}
	if task.ComplexityTier != "" {
		t.Errorf("ComplexityTier = %q, want empty", task.ComplexityTier)
	}
	if task.ComplexityReason != "" {
		t.Errorf("ComplexityReason = %q, want empty", task.ComplexityReason)
	}
}

func TestParseTaskFile_invalidComplexityTierRejected(t *testing.T) {
	p := NewParser()
	content := `---
id: E19/T003
status: planned
objective: "Bad tier"
complexity_tier: extreme
complexity_reason: "Some reason here."
depends_on: []
---

# Task

## Acceptance Criteria

- Works.`

	_, err := p.ParseTaskFile("test.md", content)
	if err == nil {
		t.Fatal("ParseTaskFile() expected error for invalid complexity_tier")
	}
	if !strings.Contains(err.Error(), "invalid complexity_tier") {
		t.Fatalf("ParseTaskFile() error = %v, want invalid complexity_tier message", err)
	}
}

func TestParseTaskFile_complexityTierWithoutReasonRejected(t *testing.T) {
	p := NewParser()
	content := `---
id: E19/T004
status: planned
objective: "Tier no reason"
complexity_tier: medium
depends_on: []
---

# Task

## Acceptance Criteria

- Works.`

	_, err := p.ParseTaskFile("test.md", content)
	if err == nil {
		t.Fatal("ParseTaskFile() expected error for complexity_tier without reason")
	}
	if !strings.Contains(err.Error(), "complexity_reason is required") {
		t.Fatalf("ParseTaskFile() error = %v, want reason required message", err)
	}
}

func TestParseTaskFile_complexityReasonWithoutTierRejected(t *testing.T) {
	p := NewParser()
	content := `---
id: E19/T005
status: planned
objective: "Reason no tier"
complexity_reason: "Some reason here."
depends_on: []
---

# Task

## Acceptance Criteria

- Works.`

	_, err := p.ParseTaskFile("test.md", content)
	if err == nil {
		t.Fatal("ParseTaskFile() expected error for complexity_reason without tier")
	}
	if !strings.Contains(err.Error(), "complexity_tier is required") {
		t.Fatalf("ParseTaskFile() error = %v, want tier required message", err)
	}
}

func TestParseTaskFile_complexityReasonTooLongRejected(t *testing.T) {
	p := NewParser()
	longReason := strings.Repeat("x", MaxComplexityReasonLen+1)
	content := "---\nid: E19/T006\nstatus: planned\nobjective: \"Long reason\"\ncomplexity_tier: low\ncomplexity_reason: \"" + longReason + "\"\ndepends_on: []\n---\n\n# Task\n\n## Acceptance Criteria\n\n- Works.\n"

	_, err := p.ParseTaskFile("test.md", content)
	if err == nil {
		t.Fatal("ParseTaskFile() expected error for oversized complexity_reason")
	}
	if !strings.Contains(err.Error(), "exceeds") {
		t.Fatalf("ParseTaskFile() error = %v, want exceeds message", err)
	}
}

func TestParseTaskFile_joinsHardWrappedChecklistItems(t *testing.T) {
	p := NewParser()
	content := `---
id: E06/T001
status: planned
objective: "Style the board"
---

# Task

## Implementation Plan

- [ ] First sentence spans across a hard markdown line break
  before it ends. Second sentence stays in the same checklist item.
- [x] Already checked sentence wraps
  without becoming another checklist item.
`

	task, err := p.ParseTaskFile("test.md", content)
	if err != nil {
		t.Fatalf("ParseTaskFile() error = %v", err)
	}
	if len(task.Checklist) != 2 {
		t.Fatalf("Task.Checklist len = %d, want 2", len(task.Checklist))
	}
	wantFirst := "First sentence spans across a hard markdown line break before it ends. Second sentence stays in the same checklist item."
	if task.Checklist[0].Text != wantFirst || task.Checklist[0].Done {
		t.Errorf("Task.Checklist[0] = %+v, want text %q and Done=false", task.Checklist[0], wantFirst)
	}
	wantSecond := "Already checked sentence wraps without becoming another checklist item."
	if task.Checklist[1].Text != wantSecond || !task.Checklist[1].Done {
		t.Errorf("Task.Checklist[1] = %+v, want text %q and Done=true", task.Checklist[1], wantSecond)
	}
}
