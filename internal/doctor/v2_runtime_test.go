package doctor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/testutil"
)

func TestRunV2ChecksDoesNotReadLegacyAuditRecords(t *testing.T) {
	root := t.TempDir()
	writeCompleteV2Project(t, root)
	auditPath := filepath.Join(root, "audit", "findings", "F001-legacy.md")
	testutil.WriteFile(t, auditPath, "not V2 audit data and intentionally malformed")

	report := RunV2Checks(root)
	if len(report.AuditRegister) != 0 {
		t.Fatalf("AuditRegister = %v, want no legacy audit check on the live V2 path", report.AuditRegister)
	}
	if report.HasProblems() {
		t.Fatalf("RunV2Checks() reported a problem from the legacy audit tree: %v", report.HealthFindings())
	}
}

func TestRunV2ChecksUsesStrictRouterReader(t *testing.T) {
	root := t.TempDir()
	writeCompleteV2Project(t, root)
	routerPath := filepath.Join(root, "router.md")
	if err := os.WriteFile(routerPath, []byte("## Current state\n\n```yaml\nstate: audit-pending\nrelease: v1\nepic: E01\nnext_action: legacy\n```\n"), 0644); err != nil {
		t.Fatal(err)
	}

	report := RunV2Checks(root)
	if report.RouterCheck == nil {
		t.Fatal("RouterCheck = nil, want strict V2 rejection of legacy router fields")
	}
	if !strings.Contains(report.RouterCheck.Error(), "V2") && !strings.Contains(report.RouterCheck.Error(), "unknown field") {
		t.Fatalf("RouterCheck = %v, want a named V2 router diagnostic", report.RouterCheck)
	}
}

func TestRunV2ChecksReportsRetiredNextActionWithoutItsValue(t *testing.T) {
	root := t.TempDir()
	writeCompleteV2Project(t, root)
	routerPath := filepath.Join(root, "router.md")
	raw, err := os.ReadFile(routerPath)
	if err != nil {
		t.Fatalf("ReadFile(router.md) error = %v", err)
	}
	const secret = "private handoff value"
	var cleanLines []string
	for _, line := range strings.Split(string(raw), "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "next_action:") {
			continue
		}
		cleanLines = append(cleanLines, line)
	}
	baseRouter := strings.Join(cleanLines, "\n")
	withRetiredField := strings.Replace(baseRouter, "```yaml\n", "```yaml\nnext_action: \""+secret+"\"\n", 1)
	if withRetiredField == baseRouter {
		t.Fatal("fixture router has no YAML state block to extend")
	}
	if err := os.WriteFile(routerPath, []byte(withRetiredField), 0644); err != nil {
		t.Fatalf("WriteFile(router.md) error = %v", err)
	}

	report := RunV2Checks(root)
	var findings []HealthFinding
	for _, finding := range report.HealthFindings() {
		if strings.Contains(finding.Message, "router-next-action-retired") {
			findings = append(findings, finding)
		}
	}
	if len(findings) != 1 {
		t.Fatalf("retired next_action findings = %+v, want exactly one; router error = %v; project findings = %+v", findings, report.RouterCheck, report.Project)
	}
	finding := findings[0]
	if finding.Category != HealthPendingReview || finding.File != routerPath ||
		!strings.Contains(strings.ToLower(finding.Repair), "delete") || !strings.Contains(finding.Repair, "next_action") {
		t.Errorf("retired next_action finding = %+v, want a repair hint to delete the line", finding)
	}
	formatted := report.Format()
	if strings.Contains(formatted, secret) {
		t.Errorf("doctor printed the retired next_action value:\n%s", formatted)
	}

	retiredLine := "next_action: \"" + secret + "\"\n"
	cleanRouter := strings.Replace(withRetiredField, retiredLine, "", 1)
	if cleanRouter == withRetiredField {
		t.Fatal("could not remove the retired next_action line from the fixture")
	}
	if err := os.WriteFile(routerPath, []byte(cleanRouter), 0644); err != nil {
		t.Fatalf("WriteFile(clean router.md) error = %v", err)
	}
	cleanReport := RunV2Checks(root)
	for _, finding := range cleanReport.HealthFindings() {
		if strings.Contains(finding.Message, "router-next-action-retired") {
			t.Errorf("doctor still reports next_action after cleanup: %+v", finding)
		}
	}
}

func TestRunV2ChecksWarnsWhenRouterSelectsDoneTask(t *testing.T) {
	root := t.TempDir()
	writeCompleteV2Project(t, root)
	taskPath := filepath.Join(root, "objectives", "O-001-ship", "tasks", "T-001-write.md")
	raw, err := os.ReadFile(taskPath)
	if err != nil {
		t.Fatalf("ReadFile(Task) error = %v", err)
	}
	lines := strings.Split(string(raw), "\n")
	delimiters := 0
	statusFound := false
	for i, line := range lines {
		if line == "---" {
			delimiters++
			if delimiters == 2 {
				break
			}
			continue
		}
		if delimiters != 1 {
			continue
		}
		if strings.HasPrefix(line, "status:") {
			lines[i] = "status: done"
			statusFound = true
		}
		if strings.HasPrefix(line, "stage:") {
			lines[i] = ""
		}
	}
	if !statusFound {
		t.Fatal("Task frontmatter has no status to update")
	}
	if err := os.WriteFile(taskPath, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		t.Fatalf("WriteFile(Task) error = %v", err)
	}

	report := RunV2Checks(root)
	const phrase = "Warning: router still selects finished Task T-001."
	var staleSelection *HealthFinding
	for _, finding := range report.HealthFindings() {
		if finding.Message == phrase {
			staleSelection = &finding
			break
		}
	}
	if staleSelection == nil {
		t.Fatalf("HealthFindings() = %+v, want a stale-selection warning", report.HealthFindings())
	}
	if staleSelection.Category != HealthPendingReview || staleSelection.File != filepath.Join(root, "router.md") || staleSelection.Repair == "" {
		t.Errorf("stale selection finding = %+v, want a non-blocking warning on router.md with a repair hint", staleSelection)
	}
	formatted := report.Format()
	if !strings.Contains(formatted, "! project:") || !strings.Contains(formatted, phrase) || !strings.Contains(formatted, "repair: "+staleSelection.Repair) {
		t.Errorf("Format() omits the warning or repair hint:\n%s", formatted)
	}
}

// TestRunV2ChecksWarnsWhenPlannedObjectiveHasStartedTask covers I-044's
// backstop: work started without the board leaves its Objective planned, and
// doctor names it as a non-blocking warning with the status to set.
func TestRunV2ChecksWarnsWhenPlannedObjectiveHasStartedTask(t *testing.T) {
	root := t.TempDir()
	writeCompleteV2Project(t, root)
	taskPath := filepath.Join(root, "objectives", "O-001-ship", "tasks", "T-001-write.md")
	raw, err := os.ReadFile(taskPath)
	if err != nil {
		t.Fatalf("ReadFile(Task) error = %v", err)
	}
	started := strings.Replace(string(raw), "status: planned\n", "status: in_progress\nstage: build\n", 1)
	if started == string(raw) {
		t.Fatal("Task frontmatter has no planned status to update")
	}
	if err := os.WriteFile(taskPath, []byte(started), 0644); err != nil {
		t.Fatalf("WriteFile(Task) error = %v", err)
	}

	var warning *HealthFinding
	for _, finding := range RunV2Checks(root).HealthFindings() {
		if strings.Contains(finding.Message, "[v2-objective-planned-with-started-task] objective O-001") {
			warning = &finding
			break
		}
	}
	if warning == nil {
		t.Fatal("HealthFindings() has no planned-with-started-task warning for O-001")
	}
	if warning.Category != HealthPendingReview || !strings.Contains(warning.Message, "T-001") || !strings.Contains(warning.Repair, "in_progress") {
		t.Errorf("finding = %+v, want a pending-review warning naming T-001 with an in_progress repair", warning)
	}

	objectivePath := filepath.Join(root, "objectives", "O-001-ship", "Objective.md")
	objectiveRaw, err := os.ReadFile(objectivePath)
	if err != nil {
		t.Fatalf("ReadFile(Objective) error = %v", err)
	}
	if err := os.WriteFile(objectivePath, []byte(strings.Replace(string(objectiveRaw), "status: planned\n", "status: in_progress\n", 1)), 0644); err != nil {
		t.Fatalf("WriteFile(Objective) error = %v", err)
	}
	for _, finding := range RunV2Checks(root).HealthFindings() {
		if strings.Contains(finding.Message, "planned-with-started-task") {
			t.Errorf("warning %q remains after the Objective was set to in_progress", finding.Message)
		}
	}
}
