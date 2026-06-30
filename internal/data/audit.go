package data

import "sort"

// AuditRegisterSet is the repo-wide audit-register state for one .savepoint root:
// the canonical prompt, the register index, and the sorted finding and run
// records. It is the single data entry point the board and doctor read so they
// never duplicate audit/ traversal. A root without an audit/ tree yields a fully
// empty set (unavailable prompt/register, no findings, no runs) without error.
type AuditRegisterSet struct {
	Root     string
	Prompt   AuditPrompt
	Register AuditRegister
	Findings []AuditFinding
	Runs     []AuditRun
}

// LoadAuditRegisterSet loads the prompt, register, findings, and runs under the
// .savepoint root's audit/ tree, returning them with deterministic ordering.
// Each artifact is loaded by its existing tolerant loader, so a missing audit/
// directory (or any absent artifact) returns an empty/unavailable entry rather
// than an error; only a malformed finding or run aborts the load with actionable
// path context. Findings are sorted by SortFindings and runs by SortRuns.
func LoadAuditRegisterSet(root string) (AuditRegisterSet, error) {
	set := AuditRegisterSet{Root: root}

	prompt, err := LoadAuditPrompt(root)
	if err != nil {
		return AuditRegisterSet{}, err
	}
	set.Prompt = prompt

	register, err := LoadAuditRegister(root)
	if err != nil {
		return AuditRegisterSet{}, err
	}
	set.Register = register

	findings, err := LoadAuditFindings(root)
	if err != nil {
		return AuditRegisterSet{}, err
	}
	set.Findings = findings

	runs, err := LoadAuditRuns(root)
	if err != nil {
		return AuditRegisterSet{}, err
	}
	set.Runs = runs

	SortFindings(set.Findings)
	SortRuns(set.Runs)

	return set, nil
}

// SortFindings orders findings in place for the board: by lifecycle status
// (active before resolved), then severity (critical before low), then stable ID.
// The ordering is total and deterministic so repeated loads of the same tree
// render identically.
func SortFindings(findings []AuditFinding) {
	sort.SliceStable(findings, func(i, j int) bool {
		a, b := findings[i], findings[j]
		if pa, pb := findingStatusPriority(a.Status), findingStatusPriority(b.Status); pa != pb {
			return pa < pb
		}
		if sa, sb := findingSeverityPriority(a.Severity), findingSeverityPriority(b.Severity); sa != sb {
			return sa < sb
		}
		return a.ID < b.ID
	})
}

// SortRuns orders runs in place newest-first by date, with the label as a stable
// alphabetical tiebreak when two runs share a date. Dates are YYYY-MM-DD, so a
// descending string compare is chronological.
func SortRuns(runs []AuditRun) {
	sort.SliceStable(runs, func(i, j int) bool {
		a, b := runs[i], runs[j]
		if a.Date != b.Date {
			return a.Date > b.Date
		}
		return a.Label < b.Label
	})
}

// findingStatusPriority maps a finding status to its sort rank. The rank is the
// status's position in the canonical lifecycle order (findingStatuses), so open
// and other active states sort above resolved ones like waived or duplicate, and
// findingStatuses stays the single source of truth. An unrecognized status sorts
// last; load-time healing means callers normally never see one.
func findingStatusPriority(s FindingStatus) int {
	for i, st := range findingStatuses {
		if st == s {
			return i
		}
	}
	return len(findingStatuses)
}

// findingSeverityPriority maps a severity to its sort rank, critical highest. An
// unrecognized severity sorts after low; load-time healing means callers
// normally never see one.
func findingSeverityPriority(s DefectSeverity) int {
	switch s {
	case SeverityCritical:
		return 0
	case SeverityHigh:
		return 1
	case SeverityMedium:
		return 2
	case SeverityLow:
		return 3
	default:
		return 4
	}
}
