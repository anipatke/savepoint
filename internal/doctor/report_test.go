package doctor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/opencode/savepoint/internal/migrate"
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
		"Migration Check",
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

// listFiles walks dir and returns every regular file's path relative to dir,
// so a before/after comparison proves a run touched nothing.
func listFiles(t *testing.T, dir string) []string {
	t.Helper()
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			rel, relErr := filepath.Rel(dir, path)
			if relErr != nil {
				return relErr
			}
			files = append(files, rel)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return files
}

func TestDiagnosticReport_MigrationOperationIncompleteIsAProblemAndWritesNothing(t *testing.T) {
	projectDir := t.TempDir()
	root := filepath.Join(projectDir, ".savepoint")
	if err := os.MkdirAll(root, 0755); err != nil {
		t.Fatal(err)
	}
	entries := []migrate.JournalEntry{{Path: "objectives/O001.md", Action: migrate.ActionCreate}}
	if _, err := migrate.CreateOperation(projectDir, "op-9", nil, entries, time.Now()); err != nil {
		t.Fatalf("CreateOperation() error = %v", err)
	}

	before := listFiles(t, projectDir)

	report := RunAllChecks(root, "")
	if !report.HasProblems() {
		t.Fatal("RunAllChecks() should report a problem for an incomplete migration operation")
	}
	output := report.Format()
	if !strings.Contains(output, "op-9") {
		t.Errorf("report.Format() = %q, want it to name the operation", output)
	}
	if !strings.Contains(output, "--recover") {
		t.Errorf("report.Format() = %q, want the recovery command", output)
	}

	after := listFiles(t, projectDir)
	if len(before) != len(after) {
		t.Fatalf("RunAllChecks() changed the file set: before=%v after=%v", before, after)
	}
	for i := range before {
		if before[i] != after[i] {
			t.Fatalf("RunAllChecks() changed the file set: before=%v after=%v", before, after)
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

func TestDiagnosticReport_NoReleaseOmitsReleaseSection(t *testing.T) {
	root := t.TempDir()
	writeCompleteV2Project(t, root)
	report := RunAllChecks(root, "")
	if strings.Contains(report.Format(), "Release Check") {
		t.Fatal("report.Format() contains Release Check for a project with no Releases")
	}
	if len(report.Releases) != 0 || len(report.ReleaseNotes) != 0 {
		t.Fatalf("Release diagnostics = %v / %v, want none", report.Releases, report.ReleaseNotes)
	}
}

func TestDiagnosticReport_HistoricalReleaseIsNotCurrentClear(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	archivePath := filepath.Join(root, "archive", "v1", "R001-PRD.md")
	testutil.WriteFile(t, archivePath, "historical release bytes\n")
	writeV2Release(t, root, "R001-history", "R001", "done",
		"legacy_completion:\n  source_path: .savepoint/releases/v1/v1-PRD.md\n  archive_path: archive/v1/R001-PRD.md\n  sha256: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\n")
	testutil.WriteFile(t, filepath.Join(root, "objectives", "O001-member", "Objective.md"),
		"---\nid: O001\ntitle: \"Member\"\nstatus: done\nrelease: R001\nfreshness: {state: current, check: C001, assessed_by: {role: checker, session: objective-checker}, assessed_at: '2026-09-14T00:00:00Z', basis: checked}\n---\n\n# Member\n")
	writeV2Check(t, root, "C001", "{kind: objective, id: O001}", "CLEAR", "")

	report := RunAllChecks(root, "")
	if len(report.Releases) != 0 {
		t.Fatalf("Release problems = %v, want none for a valid historical archive", report.Releases)
	}
	if len(report.ReleaseNotes) != 1 || !strings.Contains(report.ReleaseNotes[0], "historical evidence, not a current CLEAR Check") {
		t.Fatalf("ReleaseNotes = %v, want historical-not-current-CLEAR note", report.ReleaseNotes)
	}
	output := report.Format()
	for _, want := range []string{"Release Check", "historical completion", "not a current CLEAR Check"} {
		if !strings.Contains(output, want) {
			t.Errorf("report.Format() missing %q, got:\n%s", want, output)
		}
	}
}

func TestDiagnosticReport_DanglingHistoricalArchiveIsProblem(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Release(t, root, "R001-history", "R001", "done",
		"legacy_completion:\n  source_path: .savepoint/releases/v1/v1-PRD.md\n  archive_path: archive/v1/missing.md\n  sha256: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\n")
	testutil.WriteFile(t, filepath.Join(root, "objectives", "O001-member", "Objective.md"),
		"---\nid: O001\ntitle: \"Member\"\nstatus: done\nrelease: R001\n---\n\n# Member\n")

	problems := CheckReleaseReadiness(root)
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-release-legacy-dangling]") {
		t.Fatalf("CheckReleaseReadiness() = %v, want one dangling historical-archive problem", problems)
	}
	if problems[0].Category != HealthMalformedData || problems[0].Repair == "" {
		t.Fatalf("historical archive problem = %+v, want malformed category and repair", problems[0])
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

// writeCompleteV2Project writes a valid, complete V2 project: config.yml
// with the required fields, a V2-shaped router.md, one Objective, one Task,
// and no releases directory — the real shape of a V2 project, distinct from
// the V1 layout writeReportProject builds.
func writeCompleteV2Project(t *testing.T, root string) {
	t.Helper()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\nquality_gates:\n  lint: null\n  typecheck: null\n  build: null\n  test: null\ntheme: {}\n")
	testutil.WriteFile(t, filepath.Join(root, "router.md"),
		"## Current state\n\n```yaml\nstate: task\nobjective: O001\ntask: T001\nnext_action: \"build it\"\n```\n")
	writeV2Objective(t, root, "O001-ship", "O001", "Ship it")
	writeV2Task(t, root, "O001-ship", "T001-write.md", "T001", "Write it", "O001")
}

// TestDiagnosticReport_AdvisoryIssueBacklogIsStructurallySound proves a
// complete V2 project whose only finding is an open Issue is reported ALL
// CLEAN: the Issue is advisory, listed by name, type, and status under
// Pending Semantic Review, and never makes the project structurally unsound
// (E48 T006).
func TestDiagnosticReport_AdvisoryIssueBacklogIsStructurallySound(t *testing.T) {
	root := t.TempDir()
	writeCompleteV2Project(t, root)
	writeV2Issue(t, root, "I001-flaky.md", "I001", "open", "defect", "")

	report := RunAllChecks(root, "")
	if report.HasProblems() {
		t.Fatalf("RunAllChecks() on a complete V2 project with only an open Issue should have no problems, got: config=%v router=%v project=%v structure=%v deps=%v gates=%v",
			report.ConfigCheck, report.RouterCheck, report.Project, report.Structure, report.Dependencies, report.Gates.Results)
	}

	output := report.Format()
	if !strings.Contains(output, "ALL CLEAN") {
		t.Errorf("report.Format() should say ALL CLEAN for an advisory-only backlog, got:\n%s", output)
	}
	wants := []string{"Pending Semantic Review", "I001", "defect issue, status open"}
	for _, want := range wants {
		if !strings.Contains(output, want) {
			t.Errorf("report.Format() missing %q, got:\n%s", want, output)
		}
	}
}

// TestDiagnosticReport_CombinedMalformedMissingEvidenceAndAdvisoryIssue
// proves the four health categories stay independent even when a single
// project carries all three non-review kinds of finding at once: a
// malformed record (a Task's own owner-acceptance evidence disagrees with
// itself), a missing-evidence target (a done Task with no Check ever
// recorded), and an advisory open Issue. The first two must make the
// project structurally unsound; the Issue must not (E48 T007).
func TestDiagnosticReport_CombinedMalformedMissingEvidenceAndAdvisoryIssue(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O001-ship", "O001", "Ship it")

	// T001: done, no Check ever recorded — missing evidence.
	testutil.WriteFile(t, filepath.Join(root, "objectives", "O001-ship", "tasks", "T001-alpha.md"),
		"---\nid: T001\ntitle: \"Alpha\"\nobjective: O001\nplanned_by: {role: planner, session: planning-001}\nstatus: done\n---\n\n# Alpha\n")

	// T002: owner accepted C002, but C003 has since superseded it as the
	// latest Check — the record's own fields disagree, so this is malformed
	// data, not missing evidence.
	writeV2Check(t, root, "C002", "{kind: task, id: T002}", "CLEAR", "")
	writeV2Check(t, root, "C003", "{kind: task, id: T002}", "CLEAR", "C002")
	testutil.WriteFile(t, filepath.Join(root, "objectives", "O001-ship", "tasks", "T002-beta.md"),
		"---\nid: T002\ntitle: \"Beta\"\nobjective: O001\nplanned_by: {role: planner, session: planning-001}\nstatus: in_progress\nstage: audit\n"+
			"freshness: {state: current, check: C003, assessed_by: {role: checker, session: s}, assessed_at: '2026-09-14T00:00:00Z', basis: rechecked}\n"+
			"owner_validation: {required: true, accepted_check: C002, accepted_by: {role: owner, session: owner-1}}\n"+
			"---\n\n# Beta\n")

	writeV2Issue(t, root, "I001-flaky.md", "I001", "open", "defect", "")

	report := RunAllChecks(root, "")
	findings := report.HealthFindings()

	var sawMalformed, sawMissingEvidence, sawPendingReview bool
	for _, f := range findings {
		switch f.Category {
		case HealthMalformedData:
			if strings.Contains(f.Message, "v2-acceptance-superseded") {
				sawMalformed = true
			}
		case HealthMissingEvidence:
			if strings.Contains(f.Message, "clearance is missing") {
				sawMissingEvidence = true
			}
		case HealthPendingReview:
			if f.File == "I001" {
				sawPendingReview = true
			}
		}
	}
	if !sawMalformed {
		t.Errorf("HealthFindings() = %+v, want a HealthMalformedData finding naming v2-acceptance-superseded", findings)
	}
	if !sawMissingEvidence {
		t.Errorf("HealthFindings() = %+v, want a HealthMissingEvidence finding naming the missing clearance", findings)
	}
	if !sawPendingReview {
		t.Errorf("HealthFindings() = %+v, want a HealthPendingReview finding for Issue I001", findings)
	}

	if !report.HasProblems() {
		t.Fatal("HasProblems() = false, want true: malformed data and missing evidence both make the project structurally unsound")
	}

	var nonReview int
	for _, f := range findings {
		if f.Category != HealthPendingReview {
			nonReview++
		}
	}
	if nonReview == 0 {
		t.Fatal("expected findings outside Pending Semantic Review; HasProblems() must not depend on the advisory Issue")
	}
}

// TestDiagnosticReport_FullRunWritesNothing proves a full RunAllChecks and
// Format() over a complete V2 project with an open Issue leaves every file
// byte-identical, satisfying FS-03 for the health-category report path.
func TestDiagnosticReport_FullRunWritesNothing(t *testing.T) {
	projectDir := t.TempDir()
	root := filepath.Join(projectDir, ".savepoint")
	writeCompleteV2Project(t, root)
	writeV2Issue(t, root, "I001-flaky.md", "I001", "open", "defect", "")

	before := listFiles(t, projectDir)
	report := RunAllChecks(root, "")
	_ = report.Format()
	after := listFiles(t, projectDir)

	if len(before) != len(after) {
		t.Fatalf("RunAllChecks()+Format() changed the file set: before=%v after=%v", before, after)
	}
	for i := range before {
		if before[i] != after[i] {
			t.Fatalf("RunAllChecks()+Format() changed the file set: before=%v after=%v", before, after)
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
