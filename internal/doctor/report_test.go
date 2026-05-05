package doctor

import (
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

func writeReportProject(t *testing.T, root string) {
	t.Helper()
	testutil.SetupMinimalProject(t, root, "v1", "E01-foo")
	testutil.WriteTask(t, root, "v1", "E01-foo", testutil.TaskFixture{
		Slug:      "T001-task",
		Status:    "planned",
		Objective: "Task",
	})
}
