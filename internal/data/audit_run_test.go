package data

import (
	"path/filepath"
	"strings"
	"testing"
)

// validRunContent is a complete, well-formed run record used as the baseline;
// individual tests mutate one field to exercise a parse or diagnose path.
const validRunContent = `---
date: 2026-06-29
auditor: agent
model: claude-opus-4-8
prompt_version: v1
commit: a1b2c3d
mode: full
coverage: examined internal/data; skipped internal/tui
source_audits:
  - 2026-06-01-seed
net_new: 2
reopened: 1
verified: 3
deferred: 4
coverage_gaps: 5
---

## Scope

Full sweep of the data package.

## Coverage

Examined internal/data. Unexamined internal/tui (time).
`

func TestParseRunFile_Valid(t *testing.T) {
	p := NewParser()
	run, err := p.ParseRunFile("2026-06-29-pre-release-sweep.md", validRunContent)
	if err != nil {
		t.Fatalf("ParseRunFile() error = %v", err)
	}
	if run.Date != "2026-06-29" {
		t.Errorf("Date = %q, want 2026-06-29", run.Date)
	}
	if run.Label != "pre-release-sweep" {
		t.Errorf("Label = %q, want pre-release-sweep", run.Label)
	}
	if run.Auditor != "agent" || run.Model != "claude-opus-4-8" {
		t.Errorf("auditor/model = %q/%q", run.Auditor, run.Model)
	}
	if run.PromptVersion != "v1" || run.Commit != "a1b2c3d" {
		t.Errorf("prompt_version/commit = %q/%q", run.PromptVersion, run.Commit)
	}
	if run.Mode != AuditModeFull {
		t.Errorf("Mode = %q, want full", run.Mode)
	}
	if run.Coverage == "" {
		t.Error("Coverage = empty, want one-line summary")
	}
	if len(run.SourceAudits) != 1 || run.SourceAudits[0] != "2026-06-01-seed" {
		t.Errorf("SourceAudits = %v", run.SourceAudits)
	}
	wantCounts := RunCounts{NetNew: 2, Reopened: 1, Verified: 3, Deferred: 4, CoverageGaps: 5}
	if run.Counts != wantCounts {
		t.Errorf("Counts = %+v, want %+v", run.Counts, wantCounts)
	}
	if !strings.Contains(run.Body, "## Scope") {
		t.Errorf("Body did not retain sections: %q", run.Body)
	}
	if diags := DiagnoseRun(run, "2026-06-29-pre-release-sweep.md"); len(diags) != 0 {
		t.Errorf("DiagnoseRun() = %v, want none for valid run", diags)
	}
}

func TestParseRunFile_AllModesValid(t *testing.T) {
	p := NewParser()
	for _, mode := range auditModes {
		t.Run(string(mode), func(t *testing.T) {
			content := strings.Replace(validRunContent, "mode: full", "mode: "+string(mode), 1)
			run, err := p.ParseRunFile("2026-06-29-sweep.md", content)
			if err != nil {
				t.Fatalf("ParseRunFile() error = %v", err)
			}
			if run.Mode != mode {
				t.Errorf("Mode = %q, want %q", run.Mode, mode)
			}
		})
	}
}

func TestParseRunFile_MissingFrontmatter(t *testing.T) {
	p := NewParser()
	if _, err := p.ParseRunFile("2026-06-29-x.md", "no frontmatter here"); err == nil {
		t.Fatal("ParseRunFile() error = nil, want error for missing frontmatter")
	}
}

func TestParseRunFile_MalformedYAML(t *testing.T) {
	p := NewParser()
	content := "---\n: invalid: yaml: [\n---\n"
	if _, err := p.ParseRunFile("2026-06-29-x.md", content); err == nil {
		t.Fatal("ParseRunFile() error = nil, want error for malformed YAML")
	}
}

func TestParseRunFilename(t *testing.T) {
	tests := []struct {
		path      string
		wantDate  string
		wantLabel string
		wantOK    bool
	}{
		{"2026-06-29-pre-release-sweep.md", "2026-06-29", "pre-release-sweep", true},
		{"runs/2026-01-02-seed.md", "2026-01-02", "seed", true},
		{"2026-06-29.md", "", "", false},
		{"scratch.md", "", "", false},
		{"June-29-sweep.md", "", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			date, label, ok := ParseRunFilename(tt.path)
			if date != tt.wantDate || label != tt.wantLabel || ok != tt.wantOK {
				t.Errorf("ParseRunFilename(%q) = %q, %q, %v; want %q, %q, %v",
					tt.path, date, label, ok, tt.wantDate, tt.wantLabel, tt.wantOK)
			}
		})
	}
}

func hasRunDiagnostic(diags []RunDiagnostic, code RunDiagnosticCode) bool {
	for _, d := range diags {
		if d.Code == code {
			return true
		}
	}
	return false
}

func TestDiagnoseRun_MissingRequiredFields(t *testing.T) {
	tests := []struct {
		name string
		find string
	}{
		{"date", "date: 2026-06-29\n"},
		{"auditor", "auditor: agent\n"},
		{"prompt_version", "prompt_version: v1\n"},
		{"commit", "commit: a1b2c3d\n"},
		{"mode", "mode: full\n"},
	}
	p := NewParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := strings.Replace(validRunContent, tt.find, "", 1)
			run, err := p.ParseRunFile("2026-06-29-sweep.md", content)
			if err != nil {
				t.Fatalf("ParseRunFile() error = %v", err)
			}
			diags := DiagnoseRun(run, "2026-06-29-sweep.md")
			if !hasRunDiagnostic(diags, RunMissingFieldCode) {
				t.Fatalf("diagnostics = %v, want a missing_field for %s", diags, tt.name)
			}
			var named bool
			for _, d := range diags {
				if strings.Contains(d.Message, tt.name) {
					named = true
				}
			}
			if !named {
				t.Errorf("diagnostics %v do not name missing field %q", diags, tt.name)
			}
		})
	}
}

func TestDiagnoseRun_InvalidValues(t *testing.T) {
	p := NewParser()
	tests := []struct {
		name     string
		find     string
		replace  string
		path     string
		wantCode RunDiagnosticCode
	}{
		{"invalid mode", "mode: full", "mode: deep", "2026-06-29-sweep.md", RunInvalidModeCode},
		{"invalid date", "date: 2026-06-29", "date: 2026/06/29", "2026-06-29-sweep.md", RunInvalidDateCode},
		{"date mismatch", "date: 2026-06-29", "date: 2026-06-30", "2026-06-29-sweep.md", RunDateMismatchCode},
		{"invalid filename", "mode: full", "mode: full", "scratch.md", RunInvalidFilenameCode},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := strings.Replace(validRunContent, tt.find, tt.replace, 1)
			run, err := p.ParseRunFile(tt.path, content)
			if err != nil {
				t.Fatalf("ParseRunFile() error = %v", err)
			}
			diags := DiagnoseRun(run, tt.path)
			if !hasRunDiagnostic(diags, tt.wantCode) {
				t.Fatalf("diagnostics = %v, want code %q", diags, tt.wantCode)
			}
		})
	}
}

func TestLoadAuditRuns_Absent(t *testing.T) {
	runs, err := LoadAuditRuns(t.TempDir())
	if err != nil {
		t.Fatalf("LoadAuditRuns() error = %v", err)
	}
	if len(runs) != 0 {
		t.Errorf("runs = %v, want empty for absent runs dir", runs)
	}
}

func TestLoadAuditRuns_Valid(t *testing.T) {
	root := t.TempDir()
	writeAuditFile(t, root, filepath.Join(auditRunsDir, "README.md"), "# Audit Runs\n\nGuidance, not a record.\n")
	writeAuditFile(t, root, filepath.Join(auditRunsDir, "2026-06-29-pre-release-sweep.md"), validRunContent)

	runs, err := LoadAuditRuns(root)
	if err != nil {
		t.Fatalf("LoadAuditRuns() error = %v", err)
	}
	if len(runs) != 1 {
		t.Fatalf("len(runs) = %d, want 1 (README must be skipped)", len(runs))
	}
	if runs[0].Label != "pre-release-sweep" {
		t.Errorf("Label = %q, want pre-release-sweep", runs[0].Label)
	}
}

func TestLoadAuditRuns_MalformedAborts(t *testing.T) {
	root := t.TempDir()
	writeAuditFile(t, root, filepath.Join(auditRunsDir, "2026-06-29-broken.md"), "no frontmatter")

	if _, err := LoadAuditRuns(root); err == nil {
		t.Fatal("LoadAuditRuns() error = nil, want actionable error for malformed run")
	}
}

func TestLoadAuditFindings_Absent(t *testing.T) {
	findings, err := LoadAuditFindings(t.TempDir())
	if err != nil {
		t.Fatalf("LoadAuditFindings() error = %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("findings = %v, want empty for absent findings dir", findings)
	}
}

func TestLoadAuditFindings_Valid(t *testing.T) {
	root := t.TempDir()
	writeAuditFile(t, root, filepath.Join(auditFindingsDir, "README.md"), "# Findings\n\nGuidance.\n")
	writeAuditFile(t, root, filepath.Join(auditFindingsDir, "F001-token-leak.md"), validFindingContent)

	findings, err := LoadAuditFindings(root)
	if err != nil {
		t.Fatalf("LoadAuditFindings() error = %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("len(findings) = %d, want 1 (README must be skipped)", len(findings))
	}
	if findings[0].ID != "F001" {
		t.Errorf("ID = %q, want F001", findings[0].ID)
	}
}

func TestLoadAuditFindings_MalformedAborts(t *testing.T) {
	root := t.TempDir()
	writeAuditFile(t, root, filepath.Join(auditFindingsDir, "F001-broken.md"), "no frontmatter")

	if _, err := LoadAuditFindings(root); err == nil {
		t.Fatal("LoadAuditFindings() error = nil, want actionable error for malformed finding")
	}
}
