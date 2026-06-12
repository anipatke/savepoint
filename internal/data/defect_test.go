package data

import (
	"path/filepath"
	"testing"

	"github.com/opencode/savepoint/internal/testutil"
)

func TestParseDefectFile_Valid(t *testing.T) {
	p := NewParser()
	content := `---
id: v1.1/D001-auth-crash
release: v1.1
status: open
severity: high
title: Auth crash on empty token
---

## Problem

Token validation panics on empty input.
`
	defect, err := p.ParseDefectFile("D001.md", content)
	if err != nil {
		t.Fatalf("ParseDefectFile() error = %v", err)
	}
	if defect.ID != "v1.1/D001-auth-crash" {
		t.Errorf("ID = %q, want %q", defect.ID, "v1.1/D001-auth-crash")
	}
	if defect.Status != DefectOpen {
		t.Errorf("Status = %q, want open", defect.Status)
	}
	if defect.Severity != SeverityHigh {
		t.Errorf("Severity = %q, want high", defect.Severity)
	}
	if defect.Title != "Auth crash on empty token" {
		t.Errorf("Title = %q", defect.Title)
	}
}

func TestParseDefectFile_ObjectiveFallback(t *testing.T) {
	p := NewParser()
	content := `---
id: v1.1/D002-missing-title
status: resolved
severity: low
objective: Fix via objective field
---
`
	defect, err := p.ParseDefectFile("D002.md", content)
	if err != nil {
		t.Fatalf("ParseDefectFile() error = %v", err)
	}
	if defect.Title != "Fix via objective field" {
		t.Errorf("Title = %q, want objective fallback", defect.Title)
	}
}

func TestParseDefectFile_InProgressMissingStageHealsToBuild(t *testing.T) {
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
		t.Fatalf("ParseDefectFile() error = %v, want missing stage healed to build", err)
	}
	if defect.Stage != StageBuild {
		t.Errorf("Stage = %q, want build default for in_progress", defect.Stage)
	}
}

func TestParseDefectFile_InProgressWithStage(t *testing.T) {
	p := NewParser()
	content := `---
id: v1.1/D004-stage
status: in_progress
stage: test
severity: critical
title: Staged defect
---
`
	defect, err := p.ParseDefectFile("D004.md", content)
	if err != nil {
		t.Fatalf("ParseDefectFile() error = %v", err)
	}
	if defect.Stage != StageTest {
		t.Errorf("Stage = %q, want test", defect.Stage)
	}
}

func TestParseDefectFile_ResolvedStatus(t *testing.T) {
	p := NewParser()
	content := `---
id: v1.1/D005-resolved
status: resolved
severity: low
title: Fixed defect
---
`
	defect, err := p.ParseDefectFile("D005.md", content)
	if err != nil {
		t.Fatalf("ParseDefectFile() error = %v", err)
	}
	if defect.Status != DefectResolved {
		t.Errorf("Status = %q, want resolved", defect.Status)
	}
	if defect.Stage != "" {
		t.Errorf("Stage = %q, want empty for resolved defect", defect.Stage)
	}
}

func TestParseDefectFile_StaleStageOutsideInProgressIsDropped(t *testing.T) {
	p := NewParser()
	content := `---
id: v1.1/D005-resolved
status: resolved
stage: audit
severity: low
title: Fixed defect
---
`
	defect, err := p.ParseDefectFile("D005.md", content)
	if err != nil {
		t.Fatalf("ParseDefectFile() error = %v, want stale stage dropped", err)
	}
	if defect.Stage != "" {
		t.Errorf("Stage = %q, want empty outside in_progress", defect.Stage)
	}
}

func TestParseDefectFile_OptionalFields(t *testing.T) {
	p := NewParser()
	content := `---
id: v1.1/D006-optional
status: open
severity: medium
introduced: v1.0.5
reference: E12/T003
title: Optional fields defect
---
`
	defect, err := p.ParseDefectFile("D006.md", content)
	if err != nil {
		t.Fatalf("ParseDefectFile() error = %v", err)
	}
	if defect.Introduced != "v1.0.5" {
		t.Errorf("Introduced = %q, want v1.0.5", defect.Introduced)
	}
	if defect.Reference != "E12/T003" {
		t.Errorf("Reference = %q, want E12/T003", defect.Reference)
	}
}

func TestParseDefectFile_MissingFrontmatter(t *testing.T) {
	p := NewParser()
	_, err := p.ParseDefectFile("bad.md", "no frontmatter here")
	if err == nil {
		t.Fatal("ParseDefectFile() error = nil, want error for missing frontmatter")
	}
}

func TestParseDefectFile_MalformedYAML(t *testing.T) {
	p := NewParser()
	content := "---\n: invalid: yaml: [\n---\n"
	_, err := p.ParseDefectFile("bad.md", content)
	if err == nil {
		t.Fatal("ParseDefectFile() error = nil, want error for malformed YAML")
	}
}

func TestParseDefectFile_InvalidStatusHealsToOpen(t *testing.T) {
	p := NewParser()
	content := `---
id: v1.1/D007-bad-status
status: blocked
severity: low
title: Bad status defect
---
`
	defect, err := p.ParseDefectFile("D007.md", content)
	if err != nil {
		t.Fatalf("ParseDefectFile() error = %v, want invalid status healed to open", err)
	}
	if defect.Status != DefectOpen {
		t.Errorf("Status = %q, want open for invalid status", defect.Status)
	}
}

func TestParseDefectFile_TaskStyleStatusAliasesResolve(t *testing.T) {
	p := NewParser()
	tests := []struct {
		status string
		want   DefectStatus
	}{
		{status: "planned", want: DefectOpen},
		{status: "todo", want: DefectOpen},
		{status: "done", want: DefectResolved},
		{status: "complete", want: DefectResolved},
		{status: "completed", want: DefectResolved},
	}
	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			content := `---
id: v1.1/D009-alias
status: ` + tt.status + `
severity: low
title: Alias defect
---
`
			defect, err := p.ParseDefectFile("D009.md", content)
			if err != nil {
				t.Fatalf("ParseDefectFile() error = %v", err)
			}
			if defect.Status != tt.want {
				t.Errorf("Status = %q, want %q for alias %q", defect.Status, tt.want, tt.status)
			}
		})
	}
}

func TestParseDefectFile_InvalidStageWhenInProgressHealsToBuild(t *testing.T) {
	p := NewParser()
	content := `---
id: v1.1/D008-bad-stage
status: in_progress
stage: review
severity: high
title: Bad stage defect
---
`
	defect, err := p.ParseDefectFile("D008.md", content)
	if err != nil {
		t.Fatalf("ParseDefectFile() error = %v, want invalid stage healed to build", err)
	}
	if defect.Stage != StageBuild {
		t.Errorf("Stage = %q, want build for invalid in_progress stage", defect.Stage)
	}
}

func TestParseDefectFile_EmptyStatusDefaultsToOpen(t *testing.T) {
	p := NewParser()
	content := `---
id: v1.1/D009-no-status
severity: low
title: No status defect
---
`
	defect, err := p.ParseDefectFile("D009.md", content)
	if err != nil {
		t.Fatalf("ParseDefectFile() error = %v", err)
	}
	if defect.Status != DefectOpen {
		t.Errorf("Status = %q, want open (default)", defect.Status)
	}
}

func TestListDefects_MissingDirectoryIsZeroDefects(t *testing.T) {
	d := NewDiscover()
	root := t.TempDir()
	savepointRoot := filepath.Join(root, ".savepoint")
	testutil.MkdirAll(t, filepath.Join(savepointRoot, "releases", "v1.1"))

	defects, err := d.ListDefects(savepointRoot, "v1.1")
	if err != nil {
		t.Fatalf("ListDefects() error = %v, want nil for missing defects dir", err)
	}
	if len(defects) != 0 {
		t.Errorf("ListDefects() returned %d defects, want 0", len(defects))
	}
}

func TestListDefects_ReturnsSortedMdFiles(t *testing.T) {
	d := NewDiscover()
	root := t.TempDir()
	savepointRoot := filepath.Join(root, ".savepoint")
	defectsDir := filepath.Join(savepointRoot, "releases", "v1.1", "defects")
	testutil.MkdirAll(t, defectsDir)
	testutil.WriteFile(t, filepath.Join(defectsDir, "D002-slow-query.md"), "test")
	testutil.WriteFile(t, filepath.Join(defectsDir, "D001-crash.md"), "test")
	testutil.WriteFile(t, filepath.Join(defectsDir, "notes.txt"), "not a defect")

	defects, err := d.ListDefects(savepointRoot, "v1.1")
	if err != nil {
		t.Fatalf("ListDefects() error = %v", err)
	}
	if len(defects) != 2 {
		t.Fatalf("ListDefects() returned %d defects, want 2", len(defects))
	}
	if defects[0].ID != "D001-crash" || defects[1].ID != "D002-slow-query" {
		t.Errorf("ListDefects() IDs = [%s, %s], want [D001-crash, D002-slow-query]", defects[0].ID, defects[1].ID)
	}
}

func TestListDefects_EmptyDirectory(t *testing.T) {
	d := NewDiscover()
	root := t.TempDir()
	savepointRoot := filepath.Join(root, ".savepoint")
	testutil.MkdirAll(t, filepath.Join(savepointRoot, "releases", "v1.1", "defects"))

	defects, err := d.ListDefects(savepointRoot, "v1.1")
	if err != nil {
		t.Fatalf("ListDefects() error = %v", err)
	}
	if len(defects) != 0 {
		t.Errorf("ListDefects() returned %d defects, want 0 for empty dir", len(defects))
	}
}

func TestDefectSeverityConstants(t *testing.T) {
	tests := []struct {
		input DefectSeverity
		want  string
	}{
		{SeverityCritical, "critical"},
		{SeverityHigh, "high"},
		{SeverityMedium, "medium"},
		{SeverityLow, "low"},
	}
	for _, tt := range tests {
		if string(tt.input) != tt.want {
			t.Errorf("DefectSeverity = %v, want %v", tt.input, tt.want)
		}
	}
}
