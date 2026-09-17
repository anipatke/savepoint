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
		"Project Check",
		"Structure Check",
		"Dependency Check",
		"Audit State Check",
		"Orphan Check",
		"Defect Check",
		"Audit Register Check",
		"Issue Posture",
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

// TestDiagnosticReport_IssuePostureAdvisoryOnlyOnV1Project proves the Issue
// Posture section reports "(not a V2 project)" for a V1 project and never
// contributes to HasProblems.
func TestDiagnosticReport_IssuePostureAdvisoryOnlyOnV1Project(t *testing.T) {
	root := t.TempDir()
	writeReportProject(t, root)
	report := RunAllChecks(root, "")

	if report.Issues != nil {
		t.Fatalf("Issues = %+v, want nil for a V1 project", report.Issues)
	}
	output := report.Format()
	if !strings.Contains(output, "(not a V2 project)") {
		t.Errorf("report.Format() missing V1 Issue Posture note, got:\n%s", output)
	}
	if report.HasProblems() {
		t.Fatal("HasProblems() = true, want false: Issue posture must never contribute to health")
	}
}

// TestDiagnosticReport_IssuePostureCountsInPlainOutput proves an open Issue's
// counts reach the plain doctor output and never flip HasProblems, matching
// the rule that Issue posture is advisory only.
func TestDiagnosticReport_IssuePostureCountsInPlainOutput(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	testutil.WriteFile(t, filepath.Join(root, "issues", "I001-flaky.md"),
		"---\nid: I001\ntitle: \"Flaky\"\ntype: defect\nstatus: open\n"+
			"source: {kind: report, actor: {role: owner, session: owner-1}, at: '2026-09-15T00:00:00Z'}\n---\n\n# Issue\n")

	report := RunAllChecks(root, "")
	if report.Issues == nil || report.Issues.StatusCounts["open"] != 1 {
		t.Fatalf("Issues = %+v, want status open=1", report.Issues)
	}
	output := report.Format()
	if !strings.Contains(output, "open=1") {
		t.Errorf("report.Format() missing open Issue count, got:\n%s", output)
	}
	if len(report.Project) != 0 {
		t.Fatalf("Project = %v, want no problems: an open Issue alone does not fail the load", report.Project)
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
