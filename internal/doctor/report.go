package doctor

import (
	"fmt"
	"strings"
)

// QualityGateReport wraps quality gate results.
type QualityGateReport struct {
	Results []GateResult
}

// DiagnosticReport is the complete output of all doctor checks.
type DiagnosticReport struct {
	ConfigCheck    error
	RouterCheck    error
	Structure      []Problem
	Dependencies   []Problem
	AuditState     []Problem
	Orphans        []Problem
	Gates          QualityGateReport
	EpicFilter     string
}

// RunAllChecks runs every doctor check and returns a full report.
func RunAllChecks(root string, epicFilter string) *DiagnosticReport {
	report := &DiagnosticReport{
		EpicFilter: epicFilter,
	}

	report.ConfigCheck = CheckConfig(root)
	report.RouterCheck = CheckRouter(root, epicFilter)
	report.Structure = CheckStructure(root, epicFilter)
	report.Dependencies = CheckDependencies(root, epicFilter)
	report.AuditState = CheckAuditState(root)
	report.Orphans = CheckOrphans(root)
	report.Gates.Results = RunQualityGates(root)

	return report
}

// HasProblems returns true if any check found issues.
func (r *DiagnosticReport) HasProblems() bool {
	if r.ConfigCheck != nil {
		return true
	}
	if r.RouterCheck != nil {
		return true
	}
	if len(r.Structure) > 0 {
		return true
	}
	if len(r.Dependencies) > 0 {
		return true
	}
	if len(r.AuditState) > 0 {
		return true
	}
	if len(r.Orphans) > 0 {
		return true
	}
	for _, g := range r.Gates.Results {
		if !g.Passed {
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

	sectionHeader(&b, "Structure Check")
	printProblems(&b, "structure", r.Structure)

	sectionHeader(&b, "Dependency Check")
	printProblems(&b, "dependency", r.Dependencies)

	sectionHeader(&b, "Audit State Check")
	printProblems(&b, "audit", r.AuditState)

	sectionHeader(&b, "Orphan Check")
	printProblems(&b, "orphan", r.Orphans)

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
	if r.HasProblems() {
		b.WriteString("result: PROBLEMS FOUND (exit code 1)\n")
	} else {
		b.WriteString("result: ALL CLEAN (exit code 0)\n")
	}

	return b.String()
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
		fmt.Fprintf(b, "    repair: %s\n", SuggestRepair(p))
	}
	b.WriteString("\n")
}
