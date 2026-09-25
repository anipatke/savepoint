package data

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// FindingStatus is a point in the audit finding lifecycle. The supported values
// mirror the v1.4 PRD. Non-canonical values are healed to open at load time so a
// malformed finding never blocks load; doctor surfaces them via DiagnoseFinding.
type FindingStatus string

const (
	FindingOpen          FindingStatus = "open"
	FindingTriaged       FindingStatus = "triaged"
	FindingMapped        FindingStatus = "mapped"
	FindingInProgress    FindingStatus = "in_progress"
	FindingFixed         FindingStatus = "fixed"
	FindingVerified      FindingStatus = "verified"
	FindingDeferred      FindingStatus = "deferred"
	FindingOwnerDecision FindingStatus = "owner_decision"
	FindingWaived        FindingStatus = "waived"
	FindingDuplicate     FindingStatus = "duplicate"
)

var findingStatuses = []FindingStatus{
	FindingOpen, FindingTriaged, FindingMapped, FindingInProgress, FindingFixed,
	FindingVerified, FindingDeferred, FindingOwnerDecision, FindingWaived, FindingDuplicate,
}

// FindingConfidence is the auditor's confidence that a finding is real.
type FindingConfidence string

const (
	ConfidenceHigh   FindingConfidence = "high"
	ConfidenceMedium FindingConfidence = "medium"
	ConfidenceLow    FindingConfidence = "low"
)

var findingConfidences = []FindingConfidence{ConfidenceHigh, ConfidenceMedium, ConfidenceLow}

// findingIDPattern is the stable finding ID shape: F followed by at least three
// digits (e.g. F001). The number is assigned once and never reused.
var findingIDPattern = regexp.MustCompile(`^F\d{3,}$`)

// AuditFinding is a durable, one-file-per-finding audit record parsed from
// findings/F###-slug.md. Findings persist across audit runs so the register can
// converge instead of restarting from a cold scan. Severity reuses the shared
// DefectSeverity vocabulary so the critical/high/medium/low set has one source.
type AuditFinding struct {
	ID          string            `yaml:"id"`
	Title       string            `yaml:"title"`
	Status      FindingStatus     `yaml:"status"`
	Severity    DefectSeverity    `yaml:"severity"`
	Confidence  FindingConfidence `yaml:"confidence"`
	ProofNeeded string            `yaml:"proof_needed"`
	FirstSeen   string            `yaml:"first_seen"`
	LastSeen    string            `yaml:"last_seen"`

	SourceAuditor  string   `yaml:"source_auditor,omitempty"`
	WorkItem       string   `yaml:"work_item,omitempty"`
	Releases       []string `yaml:"releases,omitempty"`
	Epics          []string `yaml:"epics,omitempty"`
	Tasks          []string `yaml:"tasks,omitempty"`
	Defects        []string `yaml:"defects,omitempty"`
	GuardrailIDs   []string `yaml:"guardrail_ids,omitempty"`
	Locations      []string `yaml:"locations,omitempty"`
	DuplicateOf    string   `yaml:"duplicate_of,omitempty"`
	DeferralReason string   `yaml:"deferral_reason,omitempty"`
	WaiverReason   string   `yaml:"waiver_reason,omitempty"`
	VerifiedProof  string   `yaml:"verified_proof,omitempty"`

	Body  string    `yaml:"-"`
	Path  string    `yaml:"-"`
	Mtime time.Time `yaml:"-"`
}

// ParseFindingFile parses a finding record from F###-slug.md content. Only a
// structural failure (missing frontmatter, malformed YAML) returns an error;
// recoverable field problems are healed at load time so a finding never blocks
// the board. The path supplies the filename used to recover a stable ID. Doctor
// reports every healed condition via DiagnoseFinding read from raw frontmatter.
func (p *Parser) ParseFindingFile(path string, content string) (*AuditFinding, error) {
	finding, err := p.ParseRawFindingFile(path, content)
	if err != nil {
		return nil, err
	}
	NormalizeFindingForLoad(finding, path)
	return finding, nil
}

// ParseRawFindingFile parses a finding record without load-time healing so the
// original field values survive for DiagnoseFinding. Doctor uses it to report
// the problems that NormalizeFindingForLoad would silently heal.
func (p *Parser) ParseRawFindingFile(path string, content string) (*AuditFinding, error) {
	fm, body, err := SplitFrontmatterBody(normalizeLineEndings(content))
	if err != nil {
		return nil, fmt.Errorf("parse error for %s: %w", path, err)
	}

	var finding AuditFinding
	if err := yaml.Unmarshal([]byte(fm), &finding); err != nil {
		return nil, fmt.Errorf("parse error for %s: failed to parse YAML: %w", path, err)
	}
	finding.Body = body
	finding.Path = path

	return &finding, nil
}

// NormalizeFindingForLoad heals recoverable finding fields in place so a load
// never fails on bad data: non-canonical status, severity, and confidence fall
// back to their defaults, and the frontmatter ID is recovered from the filename
// when it is missing, malformed, or disagrees with the file it lives in. Text
// fields with no safe default (title, proof, dates) are left as-is and reported
// by DiagnoseFinding instead.
func NormalizeFindingForLoad(f *AuditFinding, path string) {
	if id := findingIDFromFilename(path); id != "" && f.ID != id {
		f.ID = id
	}
	if !isValidFindingStatus(f.Status) {
		f.Status = FindingOpen
	}
	if !isValidFindingSeverity(f.Severity) {
		f.Severity = SeverityMedium
	}
	if !isValidFindingConfidence(f.Confidence) {
		f.Confidence = ConfidenceMedium
	}
}

// findingIDFromFilename extracts the F### prefix from an F###-slug.md path.
// It returns "" when the filename has no recognizable finding ID, leaving the
// frontmatter ID as the sole authority in that case.
func findingIDFromFilename(path string) string {
	base := strings.TrimSuffix(filepath.Base(path), ".md")
	if base == "" {
		return ""
	}
	id := base
	if dash := strings.IndexByte(base, '-'); dash != -1 {
		id = base[:dash]
	}
	if !findingIDPattern.MatchString(id) {
		return ""
	}
	return id
}

func isValidFindingStatus(s FindingStatus) bool {
	for _, valid := range findingStatuses {
		if s == valid {
			return true
		}
	}
	return false
}

func isValidFindingSeverity(s DefectSeverity) bool {
	switch s {
	case SeverityCritical, SeverityHigh, SeverityMedium, SeverityLow:
		return true
	default:
		return false
	}
}

func isValidFindingConfidence(c FindingConfidence) bool {
	for _, valid := range findingConfidences {
		if c == valid {
			return true
		}
	}
	return false
}
