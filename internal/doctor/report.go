package doctor

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/opencode/savepoint/internal/data"
)

// QualityGateReport wraps quality gate results.
type QualityGateReport struct {
	Results []GateResult
}

// DiagnosticReport is the complete output of all doctor checks.
type DiagnosticReport struct {
	ConfigCheck   error
	RouterCheck   error
	Migration     []Problem
	Project       []Problem
	Releases      []Problem
	ReleaseNotes  []string
	Structure     []Problem
	Dependencies  []Problem
	AuditState    []Problem
	Orphans       []Problem
	Defects       []Problem
	AuditRegister []Problem
	Issues        *IssuePosture
	Gates         QualityGateReport
	EpicFilter    string
}

// RunAllChecks is the compatibility report for migration/history fixtures and
// legacy doctor unit coverage. Live command execution uses RunV2Checks below;
// keeping this adapter named and documented prevents an accidental call from
// looking like the V2 runtime's health boundary. The
// releases/epics/tasks structural checks below are V1's own directory shape;
// a V2 project has no releases directory by design and CheckProject already
// reports its single structural diagnostic (see CheckProject), so those
// checks run only when the project is not V2 — never for a schema this
// project's own config declares it isn't using.
func RunAllChecks(root string, epicFilter string, overrides ...DoctorDependencies) *DiagnosticReport {
	deps := doctorDependencies(overrides)
	report := &DiagnosticReport{
		EpicFilter: epicFilter,
	}

	report.ConfigCheck = CheckConfig(root)
	report.RouterCheck = CheckRouter(root, epicFilter)
	report.Migration = CheckMigration(root)
	project, projectProblems := loadProjectChecks(root, deps)
	report.Project = projectProblems
	releaseDiagnostics := releaseDiagnosticsForProject(project)
	report.Releases = releaseDiagnostics.Problems
	report.ReleaseNotes = releaseDiagnostics.Notes

	version, _ := data.ReadSchemaVersion(filepath.Join(root, "config.yml"))
	if version != data.SchemaVersionV2 {
		report.Structure = CheckStructure(root, epicFilter)
		report.Dependencies = CheckDependencies(root, epicFilter)
		report.AuditState = CheckAuditState(root)
		report.Orphans = CheckOrphans(root)
		report.Defects = CheckDefects(root)
	}

	report.AuditRegister = CheckAuditRegister(root)
	report.Issues = IssuePostureReport(root)
	report.Gates.Results = RunQualityGates(root, deps)

	return report
}

// HealthCategory is one of the four ways doctor reports a finding: a record
// that cannot be interpreted correctly, evidence that has not been produced
// yet, a configured quality gate that failed, or work still waiting on a
// human or checker to look at it. Only the first three bear on structural
// soundness; HealthPendingReview never does — an open Issue backlog is a V2
// project's normal working state, not a fault (E48).
type HealthCategory string

const (
	HealthMalformedData   HealthCategory = "Malformed Data"
	HealthMissingEvidence HealthCategory = "Missing Evidence"
	HealthFailingGates    HealthCategory = "Failing Gates"
	HealthPendingReview   HealthCategory = "Pending Semantic Review"
)

// healthCategoryOrder is the fixed, reported order of the four categories.
var healthCategoryOrder = []HealthCategory{
	HealthMalformedData,
	HealthMissingEvidence,
	HealthFailingGates,
	HealthPendingReview,
}

// HealthFinding names one doctor finding assigned to exactly one health
// category. File is empty for a finding with no single record to point at.
type HealthFinding struct {
	Category HealthCategory
	File     string
	Message  string
	Repair   string
}

// HealthFindings assigns every finding r's checks already produced to
// exactly one of the four health categories. It invents no new detection: it
// only routes the same Problems, GateResults, and pending Issues the rest of
// this report already carries.
func (r *DiagnosticReport) HealthFindings() []HealthFinding {
	var findings []HealthFinding

	if r.ConfigCheck != nil {
		findings = append(findings, HealthFinding{
			Category: HealthMalformedData,
			Message:  fmt.Sprintf("config: %s", r.ConfigCheck),
			Repair:   SuggestRepair(r.ConfigCheck),
		})
	}
	if r.RouterCheck != nil {
		findings = append(findings, HealthFinding{
			Category: HealthMalformedData,
			Message:  fmt.Sprintf("router: %s", r.RouterCheck),
			Repair:   SuggestRepair(r.RouterCheck),
		})
	}
	findings = append(findings, problemFindings(r.Migration)...)
	findings = append(findings, problemFindings(r.Project)...)
	findings = append(findings, problemFindings(r.Releases)...)
	findings = append(findings, problemFindings(r.Structure)...)
	findings = append(findings, problemFindings(r.Dependencies)...)
	findings = append(findings, problemFindings(r.AuditState)...)
	findings = append(findings, problemFindings(r.Orphans)...)
	findings = append(findings, problemFindings(r.Defects)...)
	findings = append(findings, problemFindings(r.AuditRegister)...)

	for _, g := range r.Gates.Results {
		if g.Passed {
			continue
		}
		message := g.Name
		if g.Command != "" {
			message = fmt.Sprintf("%s (%s)", g.Name, g.Command)
		}
		findings = append(findings, HealthFinding{
			Category: HealthFailingGates,
			Message:  message,
			Repair:   GateSuggestion(g.Name),
		})
	}

	if r.Issues != nil {
		for _, entry := range r.Issues.Pending {
			findings = append(findings, HealthFinding{
				Category: HealthPendingReview,
				File:     entry.ID,
				Message:  fmt.Sprintf("%s issue, status %s", entry.Type, entry.Status),
			})
		}
	}

	return findings
}

// problemFindings converts one check's Problems into HealthFindings,
// defaulting an unset Category to malformed data — the category every
// structural check that predates the health split already belongs to.
func problemFindings(problems []Problem) []HealthFinding {
	var findings []HealthFinding
	for _, p := range problems {
		category := p.Category
		if category == "" {
			category = HealthMalformedData
		}
		findings = append(findings, HealthFinding{
			Category: category,
			File:     p.File,
			Message:  p.Message,
			Repair:   problemRepair(p),
		})
	}
	return findings
}

// HasProblems returns true if the report is structurally unsound: it carries
// a malformed-data, missing-evidence, or failing-gate finding. Pending
// semantic review — an advisory Issue backlog — never contributes; it is the
// normal state of a project being worked on, not a fault.
func (r *DiagnosticReport) HasProblems() bool {
	for _, f := range r.HealthFindings() {
		if f.Category != HealthPendingReview {
			return true
		}
	}
	return false
}

// Format produces the full human-readable diagnostic report.
func (r *DiagnosticReport) Format() string {
	var b strings.Builder

	b.WriteString("savepoint doctor report\n")
	if r.EpicFilter != "" {
		fmt.Fprintf(&b, "  filtering to epic: %s\n", r.EpicFilter)
	}
	b.WriteString("────────────────────────────────\n\n")

	sectionHeader(&b, "Config Check")
	printSingleCheck(&b, "config", r.ConfigCheck)

	sectionHeader(&b, "Router Check")
	printSingleCheck(&b, "router", r.RouterCheck)

	sectionHeader(&b, "Migration Check")
	printProblems(&b, "migration", r.Migration)

	sectionHeader(&b, "Project Check")
	printProblems(&b, "project", r.Project)

	if len(r.Releases) > 0 || len(r.ReleaseNotes) > 0 {
		sectionHeader(&b, "Release Check")
		printProblems(&b, "release", r.Releases)
		for _, note := range r.ReleaseNotes {
			fmt.Fprintf(&b, "  note: %s\n", note)
		}
		if len(r.ReleaseNotes) > 0 {
			b.WriteString("\n")
		}
	}

	sectionHeader(&b, "Structure Check")
	printProblems(&b, "structure", r.Structure)

	sectionHeader(&b, "Dependency Check")
	printProblems(&b, "dependency", r.Dependencies)

	sectionHeader(&b, "Audit State Check")
	printProblems(&b, "audit", r.AuditState)

	sectionHeader(&b, "Orphan Check")
	printProblems(&b, "orphan", r.Orphans)

	sectionHeader(&b, "Defect Check")
	printProblems(&b, "defect", r.Defects)

	sectionHeader(&b, "Audit Register Check")
	printProblems(&b, "audit-register", r.AuditRegister)

	sectionHeader(&b, "Issue Posture")
	printIssuePosture(&b, r.Issues)

	sectionHeader(&b, "Quality Gates")
	for _, g := range r.Gates.Results {
		status := "PASS"
		if !g.Passed {
			status = "FAIL"
		}
		fmt.Fprintf(&b, "  [%s] %s", status, g.Name)
		if g.Command != "" {
			fmt.Fprintf(&b, " (%s)", g.Command)
		}
		b.WriteString("\n")
		if !g.Passed {
			hint := GateSuggestion(g.Name)
			fmt.Fprintf(&b, "    repair: %s\n", hint)
			if g.Output != "" {
				for _, line := range strings.Split(g.Output, "\n") {
					b.WriteString("    ")
					b.WriteString(line)
					b.WriteString("\n")
				}
			}
		}
	}
	b.WriteString("\n")

	sectionHeader(&b, "Health Summary")
	printHealthSummary(&b, r.HealthFindings())

	b.WriteString("\n")
	if r.HasProblems() {
		b.WriteString("result: PROBLEMS FOUND (exit code 1)\n")
	} else {
		b.WriteString("result: ALL CLEAN (exit code 0)\n")
	}

	return b.String()
}

// printHealthSummary labels each of the four health categories that has at
// least one finding, in fixed order, and omits a category with none rather
// than printing an empty heading.
func printHealthSummary(b *strings.Builder, findings []HealthFinding) {
	if len(findings) == 0 {
		fmt.Fprintf(b, "  ✓ no findings\n\n")
		return
	}
	for _, category := range healthCategoryOrder {
		var matched []HealthFinding
		for _, f := range findings {
			if f.Category == category {
				matched = append(matched, f)
			}
		}
		if len(matched) == 0 {
			continue
		}
		fmt.Fprintf(b, "  %s:\n", category)
		for _, f := range matched {
			if f.File != "" {
				fmt.Fprintf(b, "    ✗ %s: %s\n", f.File, f.Message)
			} else {
				fmt.Fprintf(b, "    ✗ %s\n", f.Message)
			}
			if f.Repair != "" {
				fmt.Fprintf(b, "      repair: %s\n", f.Repair)
			}
		}
	}
	b.WriteString("\n")
}

func sectionHeader(b *strings.Builder, title string) {
	fmt.Fprintf(b, "◆ %s\n", title)
	b.WriteString(strings.Repeat("─", len(title)+2))
	b.WriteString("\n")
}

func printSingleCheck(b *strings.Builder, name string, err error) {
	if err == nil {
		fmt.Fprintf(b, "  ✓ %s\n\n", name)
	} else {
		fmt.Fprintf(b, "  ✗ %s: %s\n", name, err.Error())
		fmt.Fprintf(b, "    repair: %s\n\n", SuggestRepair(err))
	}
}

func printProblems(b *strings.Builder, category string, problems []Problem) {
	if len(problems) == 0 {
		fmt.Fprintf(b, "  ✓ no problems\n\n")
		return
	}
	for _, p := range problems {
		fmt.Fprintf(b, "  ✗ %s: %s\n", category, p.Error())
		fmt.Fprintf(b, "    repair: %s\n", problemRepair(p))
	}
	b.WriteString("\n")
}

// printIssuePosture prints the Issue backlog counts IssuePostureReport
// derived from the loaded index. An open or in_progress Issue is advisory: it
// is never printed as a problem and never affects HasProblems.
func printIssuePosture(b *strings.Builder, posture *IssuePosture) {
	if posture == nil {
		fmt.Fprintf(b, "  (not a V2 project)\n\n")
		return
	}
	fmt.Fprintf(b, "  status: open=%d in_progress=%d resolved=%d\n",
		posture.StatusCounts[data.IssueStatusOpen],
		posture.StatusCounts[data.IssueStatusInProgress],
		posture.StatusCounts[data.IssueStatusResolved])
	fmt.Fprintf(b, "  type: defect=%d drift=%d guardrail=%d verification=%d other=%d\n\n",
		posture.TypeCounts[data.IssueTypeDefect],
		posture.TypeCounts[data.IssueTypeDrift],
		posture.TypeCounts[data.IssueTypeGuardrail],
		posture.TypeCounts[data.IssueTypeVerification],
		posture.TypeCounts[data.IssueTypeOther])
}

// problemRepair prefers a problem's typed repair suggestion, falling back to
// message-matched suggestions for checks that predate typed repairs.
func problemRepair(p Problem) string {
	if p.Repair != "" {
		return p.Repair
	}
	return SuggestRepair(p)
}
