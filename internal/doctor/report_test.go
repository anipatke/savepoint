package doctor

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/testutil"
)

func TestDiagnosticReport_HasProblems(t *testing.T) {
	root := t.TempDir()
	report := RunAllChecks(root, "")
	if !report.HasProblems() {
		t.Fatal("RunAllChecks() on empty dir should have problems")
	}
}

func TestDiagnosticReport_CleanProject(t *testing.T) {
	root := t.TempDir()
	writeReportProject(t, root)
	report := RunAllChecks(root, "")
	if report.HasProblems() {
		t.Fatalf("RunAllChecks() on valid project should have no problems, got: config=%v router=%v structure=%v deps=%v audit=%v orphans=%v gates=%v",
			report.ConfigCheck, report.RouterCheck, report.Structure, report.Dependencies, report.AuditState, report.Orphans, report.Gates.Results)
	}
}

func TestDiagnosticReport_FormatContainsSections(t *testing.T) {
	root := t.TempDir()
	report := RunAllChecks(root, "")
	output := report.Format()

	sections := []string{
		"Config Check",
		"Router Check",
		"Structure Check",
		"Dependency Check",
		"Audit State Check",
		"Orphan Check",
		"Defect Check",
		"Audit Register Check",
		"Quality Gates",
		"PROBLEMS FOUND",
	}
	for _, s := range sections {
		if !strings.Contains(output, s) {
			t.Errorf("report.Format() missing section %q", s)
		}
	}
}

func TestDiagnosticReport_FormatWithEpicFilter(t *testing.T) {
	root := t.TempDir()
	report := RunAllChecks(root, "E03")
	output := report.Format()
	if !strings.Contains(output, "filtering to epic: E03") {
		t.Errorf("report.Format() missing epic filter: %s", output)
	}
}

func TestDiagnosticReport_FormatAllClean(t *testing.T) {
	root := t.TempDir()
	writeReportProject(t, root)
	report := RunAllChecks(root, "")
	output := report.Format()
	if !strings.Contains(output, "ALL CLEAN") {
		t.Errorf("report.Format() on clean project should say ALL CLEAN, got: %s", output)
	}
	if strings.Contains(output, "PROBLEMS FOUND") {
		t.Errorf("report.Format() on clean project should not say PROBLEMS FOUND, got: %s", output)
	}
}

func TestDiagnosticReport_FormatShowsRepairs(t *testing.T) {
	root := t.TempDir()
	report := RunAllChecks(root, "")
	output := report.Format()
	if !strings.Contains(output, "repair:") {
		t.Errorf("report.Format() should include repair suggestions, got: %s", output)
	}
}

// TestDiagnosticReport_AuditRegisterAbsentStaysClean proves doctor output stays
// stable for a project that has not adopted the audit register: the section
// reports no problems and the overall result stays ALL CLEAN.
func TestDiagnosticReport_AuditRegisterAbsentStaysClean(t *testing.T) {
	root := t.TempDir()
	writeReportProject(t, root)
	report := RunAllChecks(root, "")
	output := report.Format()

	if len(report.AuditRegister) != 0 {
		t.Fatalf("AuditRegister = %v, want no problems without audit/ tree", report.AuditRegister)
	}
	if !strings.Contains(output, "Audit Register Check") {
		t.Fatalf("report.Format() missing Audit Register Check section: %s", output)
	}
	if !strings.Contains(output, "ALL CLEAN") {
		t.Fatalf("report.Format() should stay ALL CLEAN without audit register, got: %s", output)
	}
}

// TestDiagnosticReport_AuditRegisterProblemsInPlainOutput proves audit-register
// diagnostics reach the plain doctor output with the file, message, and typed
// repair suggestion.
func TestDiagnosticReport_AuditRegisterProblemsInPlainOutput(t *testing.T) {
	root := t.TempDir()
	writeReportProject(t, root)
	testutil.WriteFile(t, filepath.Join(root, "audit", "findings", "F001-verified.md"),
		`---
id: F001
title: "Verified without proof"
status: verified
severity: high
confidence: high
proof_needed: "regression test"
first_seen: "2026-07-01"
last_seen: "2026-07-01"
---

# Finding
`)

	report := RunAllChecks(root, "")
	output := report.Format()

	if !report.HasProblems() {
		t.Fatal("RunAllChecks() should report problems for verified finding without proof")
	}
	wants := []string{
		"✗ audit-register: " + filepath.Join(root, "audit", "findings", "F001-verified.md"),
		"verified finding has no named proof",
		"repair: A verified finding requires named proof",
	}
	for _, want := range wants {
		if !strings.Contains(output, want) {
			t.Errorf("report.Format() missing %q, got:\n%s", want, output)
		}
	}
}

func writeReportProject(t *testing.T, root string) {
	t.Helper()
	testutil.SetupMinimalProject(t, root, "v1", "E01-foo")
	testutil.WriteTask(t, root, "v1", "E01-foo", testutil.TaskFixture{
		Slug:      "T001-task",
		Status:    "planned",
		Objective: "Task",
	})
}
