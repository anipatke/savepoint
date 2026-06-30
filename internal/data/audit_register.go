package data

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Audit-register artifact locations relative to the .savepoint root. The audit
// domain lives under a single audit/ directory beside the release tree.
const (
	auditDirName      = "audit"
	auditPromptFile   = "prompt.md"
	auditRegisterFile = "register.md"
	auditRunsDir      = "runs"
	auditFindingsDir  = "findings"
)

// AuditPrompt is the canonical, reusable audit prompt loaded from
// audit/prompt.md. Available reports whether the file was found and read; an
// absent prompt carries an empty Body rather than aborting the load, mirroring
// ReleaseDoc. Version is the latest changelog version (e.g. "v1") when the
// prompt records one, so runs can be compared against the prompt they used.
type AuditPrompt struct {
	Path      string // path relative to the .savepoint root, for display
	Available bool
	Body      string
	Version   string
}

// LoadAuditPrompt reads the audit prompt from the .savepoint root. A missing
// prompt yields an unavailable entry; only an unexpected read error aborts the
// load and is returned with path context.
func LoadAuditPrompt(root string) (AuditPrompt, error) {
	rel := filepath.Join(auditDirName, auditPromptFile)
	prompt := AuditPrompt{Path: rel}

	content, err := os.ReadFile(filepath.Join(root, rel))
	switch {
	case err == nil:
		prompt.Available = true
		prompt.Body = string(content)
		prompt.Version = parsePromptVersion(prompt.Body)
	case os.IsNotExist(err):
		prompt.Available = false
	default:
		return AuditPrompt{}, fmt.Errorf("read audit prompt %s: %w", rel, err)
	}
	return prompt, nil
}

// RegisterSummary holds the convergence counts the register tracks for the most
// recent reconciled state, so progress is visible run over run.
type RegisterSummary struct {
	NetNew       int
	Reopened     int
	Verified     int
	Deferred     int
	CoverageGaps int
}

// AuditRegister is the current, mutable index of audit findings loaded from
// audit/register.md. Available reports whether the file was found; HasSummary
// reports whether a convergence summary was parsed from the body so callers can
// distinguish an all-zero summary from a missing one.
type AuditRegister struct {
	Path       string // path relative to the .savepoint root, for display
	Available  bool
	Body       string
	Summary    RegisterSummary
	HasSummary bool
}

// LoadAuditRegister reads the audit register from the .savepoint root. A missing
// register yields an unavailable entry; only an unexpected read error aborts the
// load. Convergence counts are parsed best-effort from the body so a malformed
// summary degrades to HasSummary=false rather than failing the load.
func LoadAuditRegister(root string) (AuditRegister, error) {
	rel := filepath.Join(auditDirName, auditRegisterFile)
	register := AuditRegister{Path: rel}

	content, err := os.ReadFile(filepath.Join(root, rel))
	switch {
	case err == nil:
		register.Available = true
		register.Body = string(content)
		register.Summary, register.HasSummary = parseRegisterSummary(register.Body)
	case os.IsNotExist(err):
		register.Available = false
	default:
		return AuditRegister{}, fmt.Errorf("read audit register %s: %w", rel, err)
	}
	return register, nil
}

// promptVersionPattern matches a bold changelog version token like **v1
// (initial)** in audit/prompt.md, capturing the vN identifier.
var promptVersionPattern = regexp.MustCompile(`\*\*\s*(v\d+)`)

// parsePromptVersion returns the highest-numbered changelog version recorded in
// the prompt, or "" when there is no Changelog section. Scoping strictly to the
// Changelog keeps prose elsewhere (e.g. an example version in the body) from
// masquerading as the prompt's version.
func parsePromptVersion(body string) string {
	idx := strings.Index(body, "## Changelog")
	if idx == -1 {
		return ""
	}
	section := body[idx:]

	best := ""
	bestN := -1
	for _, m := range promptVersionPattern.FindAllStringSubmatch(section, -1) {
		n, err := strconv.Atoi(m[1][1:])
		if err != nil {
			continue
		}
		if n > bestN {
			bestN, best = n, m[1]
		}
	}
	return best
}

// summaryRowPattern matches a two-column markdown table row whose second cell is
// an integer, e.g. "| Net-new | 3 |", capturing the label and count.
var summaryRowPattern = regexp.MustCompile(`^\|\s*([^|]+?)\s*\|\s*(-?\d+)\s*\|`)

// parseRegisterSummary extracts the convergence counts from the register body's
// summary table. It returns found=false when no recognized metric row is present
// so an absent summary is distinguishable from an all-zero one.
func parseRegisterSummary(body string) (RegisterSummary, bool) {
	var summary RegisterSummary
	found := false
	for _, line := range strings.Split(body, "\n") {
		m := summaryRowPattern.FindStringSubmatch(strings.TrimSpace(line))
		if m == nil {
			continue
		}
		count, err := strconv.Atoi(m[2])
		if err != nil {
			continue
		}
		switch strings.ToLower(strings.TrimSpace(m[1])) {
		case "net-new", "net new":
			summary.NetNew, found = count, true
		case "reopened":
			summary.Reopened, found = count, true
		case "verified":
			summary.Verified, found = count, true
		case "deferred":
			summary.Deferred, found = count, true
		case "coverage gaps", "coverage-gaps":
			summary.CoverageGaps, found = count, true
		}
	}
	return summary, found
}
