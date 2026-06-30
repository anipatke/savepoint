package data

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// AuditMode is how broadly a run examined the repo. Non-canonical values are
// left as-is and surfaced by DiagnoseRun rather than healed: a run is immutable
// history, so its recorded mode is preserved verbatim.
type AuditMode string

const (
	AuditModeFull        AuditMode = "full"
	AuditModeIncremental AuditMode = "incremental"
	AuditModeTargeted    AuditMode = "targeted"
)

var auditModes = []AuditMode{AuditModeFull, AuditModeIncremental, AuditModeTargeted}

// runDatePattern is the run date shape: a YYYY-MM-DD calendar day.
var runDatePattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

// runFilenamePattern is the immutable run filename shape: YYYY-MM-DD-label.md,
// capturing the date and the label slug.
var runFilenamePattern = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})-(.+)$`)

// RunCounts are the headline convergence counts a run reports for itself. They
// must match the dispositions in the run body and feed the register summary.
type RunCounts struct {
	NetNew       int
	Reopened     int
	Verified     int
	Deferred     int
	CoverageGaps int
}

// AuditRun is one immutable audit run record parsed from
// runs/YYYY-MM-DD-label.md. The register is derived from these append-only
// records plus disposition work, so a run is never edited after it is written.
// Date comes from the frontmatter and Label from the filename; DiagnoseRun
// reports when the two disagree.
type AuditRun struct {
	Date          string
	Label         string
	Auditor       string
	Model         string
	PromptVersion string
	Commit        string
	Mode          AuditMode
	Coverage      string
	SourceAudits  []string
	Counts        RunCounts

	Body string `yaml:"-"`
	Path string `yaml:"-"`
}

// runFrontmatter is the on-disk YAML shape of a run record. It is unmarshalled
// directly and copied into AuditRun so the public type stays free of yaml tags
// for the filename-derived and body fields.
type runFrontmatter struct {
	Date          string    `yaml:"date"`
	Auditor       string    `yaml:"auditor"`
	Model         string    `yaml:"model"`
	PromptVersion string    `yaml:"prompt_version"`
	Commit        string    `yaml:"commit"`
	Mode          AuditMode `yaml:"mode"`
	Coverage      string    `yaml:"coverage"`
	SourceAudits  []string  `yaml:"source_audits"`
	NetNew        int       `yaml:"net_new"`
	Reopened      int       `yaml:"reopened"`
	Verified      int       `yaml:"verified"`
	Deferred      int       `yaml:"deferred"`
	CoverageGaps  int       `yaml:"coverage_gaps"`
}

// ParseRunFile parses a run record from runs/YYYY-MM-DD-label.md content. Only a
// structural failure (missing frontmatter, malformed YAML) returns an error so a
// malformed run is reported with actionable path context. Field- and shape-level
// problems are not healed (a run is immutable history) and are surfaced by
// DiagnoseRun instead. The path supplies the label and, when the frontmatter
// date is absent, lets doctor recover the run's date.
func (p *Parser) ParseRunFile(path string, content string) (*AuditRun, error) {
	fm, body, err := SplitFrontmatterBody(normalizeLineEndings(content))
	if err != nil {
		return nil, fmt.Errorf("parse error for %s: %w", path, err)
	}

	var fields runFrontmatter
	if err := yaml.Unmarshal([]byte(fm), &fields); err != nil {
		return nil, fmt.Errorf("parse error for %s: failed to parse YAML: %w", path, err)
	}

	run := &AuditRun{
		Date:          fields.Date,
		Auditor:       fields.Auditor,
		Model:         fields.Model,
		PromptVersion: fields.PromptVersion,
		Commit:        fields.Commit,
		Mode:          fields.Mode,
		Coverage:      fields.Coverage,
		SourceAudits:  fields.SourceAudits,
		Counts: RunCounts{
			NetNew:       fields.NetNew,
			Reopened:     fields.Reopened,
			Verified:     fields.Verified,
			Deferred:     fields.Deferred,
			CoverageGaps: fields.CoverageGaps,
		},
		Body: body,
		Path: path,
	}
	if _, label, ok := ParseRunFilename(path); ok {
		run.Label = label
	}

	return run, nil
}

// ParseRunFilename extracts the date and label from a YYYY-MM-DD-label.md path.
// ok is false when the filename does not match the run naming convention, so
// callers can flag the file rather than guessing a date or label.
func ParseRunFilename(path string) (date, label string, ok bool) {
	base := strings.TrimSuffix(filepath.Base(path), ".md")
	m := runFilenamePattern.FindStringSubmatch(base)
	if m == nil {
		return "", "", false
	}
	return m[1], m[2], true
}

// RunDiagnosticCode classifies a recoverable run problem so callers can react to
// a class without string matching the message.
type RunDiagnosticCode string

const (
	RunMissingFieldCode    RunDiagnosticCode = "missing_field"
	RunInvalidDateCode     RunDiagnosticCode = "invalid_date"
	RunInvalidModeCode     RunDiagnosticCode = "invalid_mode"
	RunInvalidFilenameCode RunDiagnosticCode = "invalid_filename"
	RunDateMismatchCode    RunDiagnosticCode = "date_mismatch"
)

// RunDiagnostic is a warning about a run problem that ParseRunFile tolerated
// (missing or malformed field, or a filename that breaks the naming rule).
type RunDiagnostic struct {
	Code    RunDiagnosticCode
	Message string
}

// DiagnoseRun reports every actionable problem in a run record, mirroring
// DiagnoseFinding: ParseRunFile loads tolerantly, and doctor passes the parsed
// run here to recover missing required fields, an invalid mode or date, a
// filename that breaks the YYYY-MM-DD-label.md rule, and a frontmatter date that
// disagrees with the filename.
func DiagnoseRun(r *AuditRun, path string) []RunDiagnostic {
	var diagnostics []RunDiagnostic

	fileDate, _, filenameOK := ParseRunFilename(path)
	if !filenameOK {
		diagnostics = append(diagnostics, RunDiagnostic{
			Code:    RunInvalidFilenameCode,
			Message: fmt.Sprintf("audit run filename %q is not YYYY-MM-DD-label.md", filepath.Base(path)),
		})
	}

	switch {
	case r.Date == "":
		diagnostics = appendRunMissingField(diagnostics, "date")
	case !runDatePattern.MatchString(r.Date):
		diagnostics = append(diagnostics, RunDiagnostic{
			Code:    RunInvalidDateCode,
			Message: fmt.Sprintf("audit run date invalid %q; use YYYY-MM-DD", r.Date),
		})
	case filenameOK && fileDate != r.Date:
		diagnostics = append(diagnostics, RunDiagnostic{
			Code:    RunDateMismatchCode,
			Message: fmt.Sprintf("audit run date %q does not match filename date %q", r.Date, fileDate),
		})
	}

	if r.Auditor == "" {
		diagnostics = appendRunMissingField(diagnostics, "auditor")
	}
	if r.PromptVersion == "" {
		diagnostics = appendRunMissingField(diagnostics, "prompt_version")
	}
	if r.Commit == "" {
		diagnostics = appendRunMissingField(diagnostics, "commit")
	}

	if r.Mode == "" {
		diagnostics = appendRunMissingField(diagnostics, "mode")
	} else if !isValidAuditMode(r.Mode) {
		diagnostics = append(diagnostics, RunDiagnostic{
			Code:    RunInvalidModeCode,
			Message: fmt.Sprintf("audit run mode invalid %q; use %s", r.Mode, joinAuditModes()),
		})
	}

	return diagnostics
}

// LoadAuditRuns reads every run record under audit/runs/ from the .savepoint
// root. A missing runs directory yields no runs without error; a malformed run
// file aborts the load with actionable path context. README.md is skipped as
// guidance, not a record. Ordering is left to discovery (T003); runs are
// returned in directory order.
func LoadAuditRuns(root string) ([]AuditRun, error) {
	dir := filepath.Join(root, auditDirName, auditRunsDir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read audit runs dir: %w", err)
	}

	p := NewParser()
	var runs []AuditRun
	for _, entry := range entries {
		if !isAuditRecordFile(entry) {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read audit run %s: %w", entry.Name(), err)
		}
		run, err := p.ParseRunFile(path, string(content))
		if err != nil {
			return nil, err
		}
		runs = append(runs, *run)
	}
	return runs, nil
}

// LoadAuditFindings reads every finding record under audit/findings/ from the
// .savepoint root. A missing findings directory yields no findings without
// error; a malformed finding file aborts the load with actionable path context.
// README.md is skipped as guidance. Findings are returned in directory order;
// the deterministic sort order is owned by discovery (T003).
func LoadAuditFindings(root string) ([]AuditFinding, error) {
	dir := filepath.Join(root, auditDirName, auditFindingsDir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read audit findings dir: %w", err)
	}

	p := NewParser()
	var findings []AuditFinding
	for _, entry := range entries {
		if !isAuditRecordFile(entry) {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read audit finding %s: %w", entry.Name(), err)
		}
		finding, err := p.ParseFindingFile(path, string(content))
		if err != nil {
			return nil, err
		}
		findings = append(findings, *finding)
	}
	return findings, nil
}

// isAuditRecordFile reports whether a directory entry is a record markdown file
// (not a subdirectory, not the README guidance).
func isAuditRecordFile(entry os.DirEntry) bool {
	if entry.IsDir() {
		return false
	}
	name := entry.Name()
	return strings.HasSuffix(name, ".md") && !strings.EqualFold(name, "README.md")
}

func appendRunMissingField(diagnostics []RunDiagnostic, field string) []RunDiagnostic {
	return append(diagnostics, RunDiagnostic{
		Code:    RunMissingFieldCode,
		Message: fmt.Sprintf("audit run missing required field: %s", field),
	})
}

func isValidAuditMode(m AuditMode) bool {
	for _, valid := range auditModes {
		if m == valid {
			return true
		}
	}
	return false
}

func joinAuditModes() string {
	parts := make([]string, len(auditModes))
	for i, m := range auditModes {
		parts[i] = string(m)
	}
	return strings.Join(parts, ", ")
}
