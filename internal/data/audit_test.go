package data

import (
	"path/filepath"
	"testing"
)

func TestLoadAuditRegisterSet_Absent(t *testing.T) {
	set, err := LoadAuditRegisterSet(t.TempDir())
	if err != nil {
		t.Fatalf("LoadAuditRegisterSet() error = %v", err)
	}
	if set.Prompt.Available || set.Register.Available {
		t.Errorf("prompt/register available = %v/%v, want false for absent audit tree",
			set.Prompt.Available, set.Register.Available)
	}
	if len(set.Findings) != 0 || len(set.Runs) != 0 {
		t.Errorf("findings/runs = %d/%d, want empty for absent audit tree",
			len(set.Findings), len(set.Runs))
	}
}

func TestLoadAuditRegisterSet_Partial(t *testing.T) {
	// Findings present, but no prompt, register, or runs.
	root := t.TempDir()
	writeAuditFile(t, root, filepath.Join(auditFindingsDir, "F001-token-leak.md"), validFindingContent)

	set, err := LoadAuditRegisterSet(root)
	if err != nil {
		t.Fatalf("LoadAuditRegisterSet() error = %v", err)
	}
	if set.Prompt.Available || set.Register.Available {
		t.Errorf("prompt/register available = %v/%v, want false when only findings exist",
			set.Prompt.Available, set.Register.Available)
	}
	if len(set.Findings) != 1 || set.Findings[0].ID != "F001" {
		t.Fatalf("findings = %v, want one F001", set.Findings)
	}
	if len(set.Runs) != 0 {
		t.Errorf("runs = %d, want 0 when runs dir absent", len(set.Runs))
	}
}

func TestLoadAuditRegisterSet_Complete(t *testing.T) {
	root := t.TempDir()
	writeAuditFile(t, root, auditPromptFile, validPromptContent)
	writeAuditFile(t, root, auditRegisterFile, validRegisterContent)
	writeAuditFile(t, root, filepath.Join(auditFindingsDir, "README.md"), "# Findings\n")
	writeAuditFile(t, root, filepath.Join(auditFindingsDir, "F001-token-leak.md"), validFindingContent)
	writeAuditFile(t, root, filepath.Join(auditRunsDir, "README.md"), "# Runs\n")
	writeAuditFile(t, root, filepath.Join(auditRunsDir, "2026-06-29-pre-release-sweep.md"), validRunContent)

	set, err := LoadAuditRegisterSet(root)
	if err != nil {
		t.Fatalf("LoadAuditRegisterSet() error = %v", err)
	}
	if set.Root != root {
		t.Errorf("Root = %q, want %q", set.Root, root)
	}
	if !set.Prompt.Available || set.Prompt.Version != "v2" {
		t.Errorf("prompt available/version = %v/%q, want true/v2", set.Prompt.Available, set.Prompt.Version)
	}
	if !set.Register.Available || !set.Register.HasSummary {
		t.Errorf("register available/hasSummary = %v/%v, want true/true",
			set.Register.Available, set.Register.HasSummary)
	}
	if len(set.Findings) != 1 || set.Findings[0].ID != "F001" {
		t.Errorf("findings = %v, want one F001 (README skipped)", set.Findings)
	}
	if len(set.Runs) != 1 || set.Runs[0].Label != "pre-release-sweep" {
		t.Errorf("runs = %v, want one pre-release-sweep (README skipped)", set.Runs)
	}
}

func TestLoadAuditRegisterSet_MalformedAborts(t *testing.T) {
	root := t.TempDir()
	writeAuditFile(t, root, auditPromptFile, validPromptContent)
	writeAuditFile(t, root, filepath.Join(auditRunsDir, "2026-06-29-broken.md"), "no frontmatter")

	if _, err := LoadAuditRegisterSet(root); err == nil {
		t.Fatal("LoadAuditRegisterSet() error = nil, want actionable error for malformed run")
	}
}

func TestSortFindings(t *testing.T) {
	// Order is deliberately scrambled; each field exercises one tier of the
	// status > severity > ID precedence.
	findings := []AuditFinding{
		{ID: "F003", Status: FindingDuplicate, Severity: SeverityCritical},
		{ID: "F002", Status: FindingOpen, Severity: SeverityLow},
		{ID: "F005", Status: FindingOpen, Severity: SeverityCritical},
		{ID: "F004", Status: FindingOpen, Severity: SeverityCritical},
		{ID: "F001", Status: FindingVerified, Severity: SeverityHigh},
	}
	SortFindings(findings)

	got := []string{}
	for _, f := range findings {
		got = append(got, f.ID)
	}
	// open+critical (F004 before F005 by ID), open+low, verified+high, duplicate+critical.
	want := []string{"F004", "F005", "F002", "F001", "F003"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("SortFindings() order = %v, want %v", got, want)
		}
	}
}

func TestSortRuns(t *testing.T) {
	runs := []AuditRun{
		{Date: "2026-06-01", Label: "seed"},
		{Date: "2026-06-29", Label: "beta"},
		{Date: "2026-06-29", Label: "alpha"},
		{Date: "2026-06-15", Label: "mid"},
	}
	SortRuns(runs)

	got := []string{}
	for _, r := range runs {
		got = append(got, r.Date+"-"+r.Label)
	}
	// Newest date first; same-date runs break the tie by label ascending.
	want := []string{"2026-06-29-alpha", "2026-06-29-beta", "2026-06-15-mid", "2026-06-01-seed"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("SortRuns() order = %v, want %v", got, want)
		}
	}
}
