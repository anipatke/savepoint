package data

import (
	"os"
	"path/filepath"
	"testing"
)

// writeAuditFile writes content to <root>/audit/<rel>, creating parent dirs, and
// returns the root. It is the shared setup for the absent/empty/valid loader
// tests below.
func writeAuditFile(t *testing.T, root, rel, content string) {
	t.Helper()
	path := filepath.Join(root, auditDirName, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
}

const validPromptContent = `# Audit Prompt

Canonical reusable prompt.

## Changelog

- **v2 (coverage)** — Added examined/unexamined coverage accounting.
- **v1 (initial)** — Reconcile against the register, require stable IDs.
`

func TestLoadAuditPrompt_Absent(t *testing.T) {
	prompt, err := LoadAuditPrompt(t.TempDir())
	if err != nil {
		t.Fatalf("LoadAuditPrompt() error = %v", err)
	}
	if prompt.Available {
		t.Errorf("Available = true, want false for absent prompt")
	}
	if prompt.Body != "" || prompt.Version != "" {
		t.Errorf("absent prompt carried Body=%q Version=%q", prompt.Body, prompt.Version)
	}
	if prompt.Path != filepath.Join(auditDirName, auditPromptFile) {
		t.Errorf("Path = %q, want display path even when absent", prompt.Path)
	}
}

func TestLoadAuditPrompt_Valid(t *testing.T) {
	root := t.TempDir()
	writeAuditFile(t, root, auditPromptFile, validPromptContent)

	prompt, err := LoadAuditPrompt(root)
	if err != nil {
		t.Fatalf("LoadAuditPrompt() error = %v", err)
	}
	if !prompt.Available {
		t.Fatal("Available = false, want true")
	}
	if prompt.Version != "v2" {
		t.Errorf("Version = %q, want v2 (highest changelog version)", prompt.Version)
	}
	if prompt.Body == "" {
		t.Error("Body = empty, want prompt content")
	}
}

func TestLoadAuditPrompt_NoChangelogVersion(t *testing.T) {
	root := t.TempDir()
	writeAuditFile(t, root, auditPromptFile, "# Audit Prompt\n\nNo changelog here.\n")

	prompt, err := LoadAuditPrompt(root)
	if err != nil {
		t.Fatalf("LoadAuditPrompt() error = %v", err)
	}
	if !prompt.Available {
		t.Fatal("Available = false, want true")
	}
	if prompt.Version != "" {
		t.Errorf("Version = %q, want empty when no changelog version present", prompt.Version)
	}
}

const validRegisterContent = `# Audit Register

## Convergence summary

| Metric | Count |
|--------|-------|
| Net-new | 3 |
| Reopened | 1 |
| Verified | 2 |
| Deferred | 4 |
| Coverage gaps | 5 |

## Findings

| ID | Title | Status |
|----|-------|--------|
| F001 | Token leak | open |
`

func TestLoadAuditRegister_Absent(t *testing.T) {
	register, err := LoadAuditRegister(t.TempDir())
	if err != nil {
		t.Fatalf("LoadAuditRegister() error = %v", err)
	}
	if register.Available {
		t.Errorf("Available = true, want false for absent register")
	}
	if register.HasSummary {
		t.Errorf("HasSummary = true, want false for absent register")
	}
}

func TestLoadAuditRegister_Valid(t *testing.T) {
	root := t.TempDir()
	writeAuditFile(t, root, auditRegisterFile, validRegisterContent)

	register, err := LoadAuditRegister(root)
	if err != nil {
		t.Fatalf("LoadAuditRegister() error = %v", err)
	}
	if !register.Available {
		t.Fatal("Available = false, want true")
	}
	if !register.HasSummary {
		t.Fatal("HasSummary = false, want true")
	}
	want := RegisterSummary{NetNew: 3, Reopened: 1, Verified: 2, Deferred: 4, CoverageGaps: 5}
	if register.Summary != want {
		t.Errorf("Summary = %+v, want %+v", register.Summary, want)
	}
}

func TestLoadAuditRegister_NoSummary(t *testing.T) {
	root := t.TempDir()
	writeAuditFile(t, root, auditRegisterFile, "# Audit Register\n\nNo summary table yet.\n")

	register, err := LoadAuditRegister(root)
	if err != nil {
		t.Fatalf("LoadAuditRegister() error = %v", err)
	}
	if !register.Available {
		t.Fatal("Available = false, want true")
	}
	if register.HasSummary {
		t.Errorf("HasSummary = true, want false; absent summary must differ from an all-zero one")
	}
	if register.Summary != (RegisterSummary{}) {
		t.Errorf("Summary = %+v, want zero value", register.Summary)
	}
}

func TestParseRegisterSummary_AllZerosStillCounts(t *testing.T) {
	body := "| Net-new | 0 |\n| Reopened | 0 |\n| Verified | 0 |\n| Deferred | 0 |\n| Coverage gaps | 0 |\n"
	summary, found := parseRegisterSummary(body)
	if !found {
		t.Fatal("found = false, want true for an explicit all-zero summary")
	}
	if summary != (RegisterSummary{}) {
		t.Errorf("Summary = %+v, want zero value", summary)
	}
}

func TestParsePromptVersion(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{"single version", "## Changelog\n- **v1 (initial)** — seed.", "v1"},
		{"highest wins regardless of order", "## Changelog\n- **v1** — a\n- **v3** — c\n- **v2** — b", "v3"},
		{"no changelog", "# Audit Prompt\nbody with **v9** prose but no changelog heading", ""},
		{"empty", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parsePromptVersion(tt.body); got != tt.want {
				t.Errorf("parsePromptVersion() = %q, want %q", got, tt.want)
			}
		})
	}
}
