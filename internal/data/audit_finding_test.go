package data

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// unmarshalRawFinding parses frontmatter into an un-normalized finding, the raw
// shape doctor passes to DiagnoseFinding before any load-time healing.
func unmarshalRawFinding(t *testing.T, fm string) *AuditFinding {
	t.Helper()
	var f AuditFinding
	if err := yaml.Unmarshal([]byte(fm), &f); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}
	return &f
}

// validFindingContent is a complete, well-formed finding record used as the
// baseline; individual tests mutate one field to exercise a heal/diagnose path.
const validFindingContent = `---
id: F001
title: Token leak in logs
status: open
severity: high
confidence: medium
source_auditor: agent
work_item: E12-auth/T003-redact-tokens
guardrail_ids:
  - no-secrets-in-logs
locations:
  - internal/log/writer.go:42-58
first_seen: 2026-06-01
last_seen: 2026-06-30
proof_needed: A regression test asserting tokens are redacted before write
---

## Summary

Tokens are written to logs in plaintext.
`

func TestParseFindingFile_Valid(t *testing.T) {
	p := NewParser()
	finding, err := p.ParseFindingFile("F001-token-leak-in-logs.md", validFindingContent)
	if err != nil {
		t.Fatalf("ParseFindingFile() error = %v", err)
	}
	if finding.ID != "F001" {
		t.Errorf("ID = %q, want F001", finding.ID)
	}
	if finding.Title != "Token leak in logs" {
		t.Errorf("Title = %q", finding.Title)
	}
	if finding.Status != FindingOpen {
		t.Errorf("Status = %q, want open", finding.Status)
	}
	if finding.Severity != SeverityHigh {
		t.Errorf("Severity = %q, want high", finding.Severity)
	}
	if finding.Confidence != ConfidenceMedium {
		t.Errorf("Confidence = %q, want medium", finding.Confidence)
	}
	if finding.ProofNeeded == "" {
		t.Error("ProofNeeded = empty, want non-empty")
	}
	if finding.FirstSeen != "2026-06-01" || finding.LastSeen != "2026-06-30" {
		t.Errorf("first/last seen = %q/%q", finding.FirstSeen, finding.LastSeen)
	}
	if !strings.Contains(finding.Body, "## Summary") {
		t.Errorf("Body did not retain sections: %q", finding.Body)
	}
	if diags := DiagnoseFinding(finding, "F001-token-leak-in-logs.md"); len(diags) != 0 {
		t.Errorf("DiagnoseFinding() = %v, want none for valid finding", diags)
	}
}

func TestParseFindingFile_AllStatusesValid(t *testing.T) {
	p := NewParser()
	for _, status := range findingStatuses {
		t.Run(string(status), func(t *testing.T) {
			content := strings.Replace(validFindingContent, "status: open", "status: "+string(status), 1)
			finding, err := p.ParseFindingFile("F001-token-leak-in-logs.md", content)
			if err != nil {
				t.Fatalf("ParseFindingFile() error = %v for status %q", err, status)
			}
			if finding.Status != status {
				t.Errorf("Status = %q, want %q", finding.Status, status)
			}
		})
	}
}

func TestParseFindingFile_OptionalFields(t *testing.T) {
	p := NewParser()
	content := `---
id: F042
title: Duplicate of token leak
status: duplicate
severity: high
confidence: high
first_seen: 2026-06-10
last_seen: 2026-06-30
proof_needed: n/a
releases:
  - v1.4
epics:
  - E12-auth
tasks:
  - E12-auth/T003-redact-tokens
defects:
  - v1.4/D009-token-leak
guardrail_ids:
  - no-secrets-in-logs
locations:
  - internal/log/writer.go:42
duplicate_of: F001
deferral_reason: superseded by canonical finding
waiver_reason: accepted risk for internal build
verified_proof: TestRedactTokens covers the redaction path
---
`
	finding, err := p.ParseFindingFile("F042-dup.md", content)
	if err != nil {
		t.Fatalf("ParseFindingFile() error = %v", err)
	}
	if got := []string{
		strings.Join(finding.Releases, ","),
		strings.Join(finding.Epics, ","),
		strings.Join(finding.Tasks, ","),
		strings.Join(finding.Defects, ","),
		strings.Join(finding.GuardrailIDs, ","),
		strings.Join(finding.Locations, ","),
	}; got[0] != "v1.4" || got[1] != "E12-auth" || got[2] != "E12-auth/T003-redact-tokens" ||
		got[3] != "v1.4/D009-token-leak" || got[4] != "no-secrets-in-logs" || got[5] != "internal/log/writer.go:42" {
		t.Errorf("optional link fields = %v", got)
	}
	if finding.DuplicateOf != "F001" {
		t.Errorf("DuplicateOf = %q, want F001", finding.DuplicateOf)
	}
	if finding.DeferralReason == "" || finding.WaiverReason == "" || finding.VerifiedProof == "" {
		t.Errorf("rationale/proof fields not parsed: %+v", finding)
	}
}

func TestParseFindingFile_MissingFrontmatter(t *testing.T) {
	p := NewParser()
	if _, err := p.ParseFindingFile("F001-x.md", "no frontmatter here"); err == nil {
		t.Fatal("ParseFindingFile() error = nil, want error for missing frontmatter")
	}
}

func TestParseFindingFile_MalformedYAML(t *testing.T) {
	p := NewParser()
	content := "---\n: invalid: yaml: [\n---\n"
	if _, err := p.ParseFindingFile("F001-x.md", content); err == nil {
		t.Fatal("ParseFindingFile() error = nil, want error for malformed YAML")
	}
}

func TestParseFindingFile_HealsInvalidEnums(t *testing.T) {
	p := NewParser()
	content := validFindingContent
	content = strings.Replace(content, "status: open", "status: closed", 1)
	content = strings.Replace(content, "severity: high", "severity: blocker", 1)
	content = strings.Replace(content, "confidence: medium", "confidence: certain", 1)

	finding, err := p.ParseFindingFile("F001-token-leak-in-logs.md", content)
	if err != nil {
		t.Fatalf("ParseFindingFile() error = %v, want healing not error", err)
	}
	if finding.Status != FindingOpen {
		t.Errorf("Status = %q, want healed to open", finding.Status)
	}
	if finding.Severity != SeverityMedium {
		t.Errorf("Severity = %q, want healed to medium", finding.Severity)
	}
	if finding.Confidence != ConfidenceMedium {
		t.Errorf("Confidence = %q, want healed to medium", finding.Confidence)
	}
}

func TestParseFindingFile_RecoversIDFromFilename(t *testing.T) {
	p := NewParser()
	tests := []struct {
		name     string
		frontID  string
		filename string
		wantID   string
	}{
		{"mismatch heals to filename", "id: F001", "F002-renamed.md", "F002"},
		{"invalid heals to filename", "id: not-an-id", "F003-bad-front-id.md", "F003"},
		{"slug-only filename keeps frontmatter id", "id: F001", "scratch.md", "F001"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := strings.Replace(validFindingContent, "id: F001", tt.frontID, 1)
			finding, err := p.ParseFindingFile(tt.filename, content)
			if err != nil {
				t.Fatalf("ParseFindingFile() error = %v", err)
			}
			if finding.ID != tt.wantID {
				t.Errorf("ID = %q, want %q", finding.ID, tt.wantID)
			}
		})
	}
}

func hasDiagnostic(diags []FindingDiagnostic, code FindingDiagnosticCode) bool {
	for _, d := range diags {
		if d.Code == code {
			return true
		}
	}
	return false
}

func TestDiagnoseFinding_MissingRequiredFields(t *testing.T) {
	tests := []struct {
		name string
		find string
	}{
		{"id", "id: F001\n"},
		{"title", "title: Token leak in logs\n"},
		{"status", "status: open\n"},
		{"severity", "severity: high\n"},
		{"confidence", "confidence: medium\n"},
		{"proof_needed", "proof_needed: A regression test asserting tokens are redacted before write\n"},
		{"first_seen", "first_seen: 2026-06-01\n"},
		{"last_seen", "last_seen: 2026-06-30\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, _, err := SplitFrontmatterBody(strings.Replace(validFindingContent, tt.find, "", 1))
			if err != nil {
				t.Fatalf("SplitFrontmatterBody() error = %v", err)
			}
			raw := unmarshalRawFinding(t, fm)
			diags := DiagnoseFinding(raw, "scratch.md")
			if !hasDiagnostic(diags, FindingMissingFieldCode) {
				t.Fatalf("diagnostics = %v, want a missing_field for %s", diags, tt.name)
			}
			var found bool
			for _, d := range diags {
				if strings.Contains(d.Message, tt.name) {
					found = true
				}
			}
			if !found {
				t.Errorf("diagnostics %v do not name missing field %q", diags, tt.name)
			}
		})
	}
}

func TestDiagnoseFinding_InvalidValues(t *testing.T) {
	tests := []struct {
		name     string
		find     string
		replace  string
		path     string
		wantCode FindingDiagnosticCode
	}{
		{"invalid id", "id: F001", "id: finding-1", "scratch.md", FindingInvalidIDCode},
		{"short id", "id: F001", "id: F1", "scratch.md", FindingInvalidIDCode},
		{"id mismatch", "id: F001", "id: F001", "F002-other.md", FindingIDMismatchCode},
		{"invalid status", "status: open", "status: closed", "F001-x.md", FindingInvalidStatusCode},
		{"invalid severity", "severity: high", "severity: blocker", "F001-x.md", FindingInvalidSeverityCode},
		{"invalid confidence", "confidence: medium", "confidence: certain", "F001-x.md", FindingInvalidConfidenceCode},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, _, err := SplitFrontmatterBody(strings.Replace(validFindingContent, tt.find, tt.replace, 1))
			if err != nil {
				t.Fatalf("SplitFrontmatterBody() error = %v", err)
			}
			raw := unmarshalRawFinding(t, fm)
			diags := DiagnoseFinding(raw, tt.path)
			if !hasDiagnostic(diags, tt.wantCode) {
				t.Fatalf("diagnostics = %v, want code %q", diags, tt.wantCode)
			}
		})
	}
}
