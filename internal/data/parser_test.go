package data

import (
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
phase: test
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

func TestParseTaskFile_rejectsPhaseOutsideInProgress(t *testing.T) {
	p := NewParser()
	content := `---
id: E06/T001
status: planned
phase: build
objective: "Style the board"
---

# Task`

	_, err := p.ParseTaskFile("test.md", content)
	if err == nil {
		t.Fatal("ParseTaskFile() expected invalid phase/status error")
	}
}

func TestParseTaskFile_rejectsInProgressWithoutPhase(t *testing.T) {
	p := NewParser()
	content := `---
id: E06/T001
status: in_progress
objective: "Style the board"
---

# Task`

	_, err := p.ParseTaskFile("test.md", content)
	if err == nil {
		t.Fatal("ParseTaskFile() expected missing phase error")
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
