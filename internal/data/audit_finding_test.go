package data

import (
	"strings"
	"testing"
)

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
