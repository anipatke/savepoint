package data

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"
)

// issueFixture builds a minimal valid Issue record so each test can override
// exactly the one field it is about, instead of restating the whole schema.
func issueFixture(overrides map[string]string) string {
	fields := map[string]string{
		"id":     "I-001",
		"title":  "\"Follow-up that outlives one evaluation\"",
		"type":   "defect",
		"status": "open",
		"source": "{kind: check, check: C-001, actor: {role: checker, session: sess-1}, at: '2026-09-15T00:00:00Z'}",
	}
	order := []string{"id", "title", "type", "status", "stage", "source", "tasks", "checks",
		"guardrail_ids", "severity", "resolution", "duplicate_of", "history"}

	for key, value := range overrides {
		fields[key] = value
	}

	var content strings.Builder
	content.WriteString("---\n")
	for _, key := range order {
		if value, ok := fields[key]; ok {
			content.WriteString(key + ": " + value + "\n")
		}
	}
	content.WriteString("---\n\n# Issue\n")
	return content.String()
}

func TestDecodeIssueV2_valid(t *testing.T) {
	content := `---
id: I-001
title: "Router drifts from its documented state table"
type: drift
status: resolved
source:
  kind: check
  check: C-001
  actor: {role: checker, session: sess-1}
  at: '2026-09-15T00:00:00Z'
tasks: [T-001, T-002]
checks: [C-001, C-002]
guardrail_ids: [DATA-02, STYLE-07]
severity: high
resolution:
  disposition: verified
  check: C-002
  actor: {role: checker, session: sess-2}
  at: '2026-09-16T00:00:00Z'
  reason: repaired and rechecked
duplicate_of: I-002
history:
  - {at: '2026-09-15T00:00:00Z', actor: {role: checker, session: sess-1}, kind: observed, note: first seen}
  - {at: '2026-09-16T00:00:00Z', actor: {role: checker, session: sess-2}, kind: rechecked, check: C-002}
---

# Issue`

	issue, err := DecodeIssueV2("issues/I-001-router-drift.md", content)
	if err != nil {
		t.Fatalf("DecodeIssueV2() error = %v", err)
	}

	if issue.ID != "I-001" {
		t.Errorf("ID = %q, want I-001", issue.ID)
	}
	if issue.Title != "Router drifts from its documented state table" {
		t.Errorf("Title = %q, want the recorded title", issue.Title)
	}
	if issue.Type != IssueTypeDrift {
		t.Errorf("Type = %q, want drift", issue.Type)
	}
	if issue.Status != IssueStatusResolved {
		t.Errorf("Status = %q, want resolved", issue.Status)
	}

	wantOrigin := IssueOrigin{
		Kind:  IssueOriginCheck,
		Check: "C-001",
		Actor: Actor{Role: ActorRoleChecker, Session: "sess-1"},
		At:    mustParseTime(t, "2026-09-15T00:00:00Z"),
	}
	if issue.Origin != wantOrigin {
		t.Errorf("Origin = %+v, want %+v", issue.Origin, wantOrigin)
	}

	if len(issue.Tasks) != 2 || issue.Tasks[0] != "T-001" || issue.Tasks[1] != "T-002" {
		t.Errorf("Tasks = %v, want [T-001 T-002]", issue.Tasks)
	}
	if len(issue.Checks) != 2 || issue.Checks[0] != "C-001" || issue.Checks[1] != "C-002" {
		t.Errorf("Checks = %v, want [C-001 C-002]", issue.Checks)
	}
	if len(issue.GuardrailIDs) != 2 || issue.GuardrailIDs[0] != "DATA-02" {
		t.Errorf("GuardrailIDs = %v, want [DATA-02 STYLE-07]", issue.GuardrailIDs)
	}
	if issue.Severity != "high" {
		t.Errorf("Severity = %q, want high", issue.Severity)
	}
	if issue.DuplicateOf != "I-002" {
		t.Errorf("DuplicateOf = %q, want I-002", issue.DuplicateOf)
	}

	if issue.Resolution == nil {
		t.Fatal("Resolution = nil, want the decoded resolution block")
	}
	if issue.Resolution.Disposition != IssueDispositionVerified || issue.Resolution.Check != "C-002" {
		t.Errorf("Resolution = %+v, want verified proved by C-002", issue.Resolution)
	}
	if issue.Resolution.Actor != (Actor{Role: ActorRoleChecker, Session: "sess-2"}) {
		t.Errorf("Resolution.Actor = %+v, want {checker sess-2}", issue.Resolution.Actor)
	}
	if !issue.Resolution.At.Equal(mustParseTime(t, "2026-09-16T00:00:00Z")) {
		t.Errorf("Resolution.At = %v, want 2026-09-16T00:00:00Z", issue.Resolution.At)
	}
	if issue.Resolution.Reason != "repaired and rechecked" {
		t.Errorf("Resolution.Reason = %q, want the recorded reason", issue.Resolution.Reason)
	}

	if len(issue.History) != 2 {
		t.Fatalf("History = %d entries, want 2", len(issue.History))
	}
	if issue.History[0].Kind != IssueHistoryObserved || issue.History[0].Note != "first seen" {
		t.Errorf("History[0] = %+v, want the observed entry first", issue.History[0])
	}
	if issue.History[1].Kind != IssueHistoryRechecked || issue.History[1].Check != "C-002" {
		t.Errorf("History[1] = %+v, want the rechecked entry naming C-002", issue.History[1])
	}
	if !issue.History[0].At.Before(issue.History[1].At) {
		t.Errorf("History order = %v then %v, want recorded order preserved", issue.History[0].At, issue.History[1].At)
	}
}

// TestDecodeIssueV2_absentOptionalsStayAbsent proves an optional field that
// was never recorded decodes as absent rather than as a permissive default
// that would read as a deliberate declaration.
func TestDecodeIssueV2_absentOptionalsStayAbsent(t *testing.T) {
	issue, err := DecodeIssueV2("issues/I-001-minimal.md", issueFixture(nil))
	if err != nil {
		t.Fatalf("DecodeIssueV2() error = %v", err)
	}

	if issue.Resolution != nil {
		t.Errorf("Resolution = %+v, want nil for an absent block", issue.Resolution)
	}
	if issue.Severity != "" {
		t.Errorf("Severity = %q, want empty for an absent field", issue.Severity)
	}
	if issue.DuplicateOf != "" {
		t.Errorf("DuplicateOf = %q, want empty for an absent field", issue.DuplicateOf)
	}
	if len(issue.Tasks) != 0 || len(issue.Checks) != 0 || len(issue.GuardrailIDs) != 0 || len(issue.History) != 0 {
		t.Errorf("optional lists = %v/%v/%v/%v, want all empty", issue.Tasks, issue.Checks, issue.GuardrailIDs, issue.History)
	}
}

func TestDecodeIssueV2_malformedID(t *testing.T) {
	tests := []struct {
		name string
		id   string
	}{
		{"missing digits", "I"},
		{"too few digits", "I01"},
		{"unhyphenated identity", "I001"},
		{"hyphenated but too few digits", "I-01"},
		{"wrong prefix letter", "C-001"},
		{"lowercase prefix", "i001"},
		{"trailing garbage", "I-001x"},
		{"empty", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecodeIssueV2("test.md", issueFixture(map[string]string{"id": "\"" + tt.id + "\""}))
			if !errors.Is(err, ErrV2InvalidID) {
				t.Fatalf("DecodeIssueV2() error = %v, want ErrV2InvalidID", err)
			}
		})
	}
}

func TestDecodeIssueV2_missingTitle(t *testing.T) {
	for _, title := range []string{"\"\"", "\"   \""} {
		_, err := DecodeIssueV2("test.md", issueFixture(map[string]string{"title": title}))
		if !errors.Is(err, ErrV2MissingField) {
			t.Fatalf("DecodeIssueV2() title %s error = %v, want ErrV2MissingField", title, err)
		}
	}
}

func TestDecodeIssueV2_type(t *testing.T) {
	for _, valid := range []IssueType{IssueTypeDefect, IssueTypeDrift, IssueTypeGuardrail, IssueTypeVerification, IssueTypeOther} {
		issue, err := DecodeIssueV2("test.md", issueFixture(map[string]string{"type": string(valid)}))
		if err != nil {
			t.Fatalf("DecodeIssueV2() type %q error = %v", valid, err)
		}
		if issue.Type != valid {
			t.Errorf("Type = %q, want %q", issue.Type, valid)
		}
	}

	tests := []struct {
		name      string
		issueType string
		wantErr   error
	}{
		{"missing", "\"\"", ErrV2MissingField},
		{"unknown", "regression", ErrV2IssueMalformed},
		{"wrong case", "Defect", ErrV2IssueMalformed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecodeIssueV2("test.md", issueFixture(map[string]string{"type": tt.issueType}))
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("DecodeIssueV2() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecodeIssueV2_statusVocabulary(t *testing.T) {
	for _, valid := range []IssueStatus{IssueStatusOpen, IssueStatusInProgress, IssueStatusResolved} {
		overrides := map[string]string{"status": string(valid)}
		if valid == IssueStatusResolved {
			overrides["resolution"] = "{disposition: accepted, actor: {role: owner, session: owner-1}, at: '2026-09-16T00:00:00Z', reason: accepted risk}"
		}
		issue, err := DecodeIssueV2("test.md", issueFixture(overrides))
		if err != nil {
			t.Fatalf("DecodeIssueV2() status %q error = %v", valid, err)
		}
		if issue.Status != valid {
			t.Errorf("Status = %q, want %q", issue.Status, valid)
		}
	}

	_, err := DecodeIssueV2("test.md", issueFixture(map[string]string{"status": "\"\""}))
	if !errors.Is(err, ErrV2MissingField) {
		t.Fatalf("DecodeIssueV2() missing status error = %v, want ErrV2MissingField", err)
	}
}

// TestDecodeIssueV2_rejectsTaskLifecycleVocabulary proves Issue status is its
// own vocabulary: a Task lifecycle value is refused outright rather than
// healed into a valid-looking Issue status, so the two vocabularies cannot
// merge by accident.
func TestDecodeIssueV2_rejectsTaskLifecycleVocabulary(t *testing.T) {
	tests := []struct {
		name      string
		overrides map[string]string
	}{
		{"planned", map[string]string{"status": string(ColumnPlanned)}},
		{"done", map[string]string{"status": string(ColumnDone)}},
		{"legacy todo", map[string]string{"status": string(LegacyTaskStatusTodo)}},
		{"legacy complete", map[string]string{"status": string(LegacyTaskStatusComplete)}},
		{"in_progress with a build stage", map[string]string{"status": "in_progress", "stage": string(StageBuild)}},
		{"in_progress with an audit stage", map[string]string{"status": "in_progress", "stage": string(StageAudit)}},
		{"open with a stage", map[string]string{"status": "open", "stage": string(StageBuild)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issue, err := DecodeIssueV2("test.md", issueFixture(tt.overrides))
			if !errors.Is(err, ErrV2InvalidLifecycle) {
				t.Fatalf("DecodeIssueV2() error = %v, want ErrV2InvalidLifecycle", err)
			}
			if issue != nil {
				t.Fatalf("DecodeIssueV2() = %+v, want no healed record alongside the diagnostic", issue)
			}
		})
	}
}

func TestDecodeIssueV2_source(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		wantErr error
	}{
		{"missing block", "", ErrV2MissingField},
		{"missing kind", "{check: C-001, actor: {role: checker, session: s}, at: '2026-09-15T00:00:00Z'}", ErrV2MissingField},
		{"unknown kind", "{kind: audit, actor: {role: checker, session: s}, at: '2026-09-15T00:00:00Z'}", ErrV2IssueMalformed},
		{"check kind without a check", "{kind: check, actor: {role: checker, session: s}, at: '2026-09-15T00:00:00Z'}", ErrV2MissingField},
		{"check kind with a malformed check", "{kind: check, check: T-001, actor: {role: checker, session: s}, at: '2026-09-15T00:00:00Z'}", ErrV2InvalidID},
		{"report kind naming a check", "{kind: report, check: C-001, actor: {role: owner, session: s}, at: '2026-09-15T00:00:00Z'}", ErrV2IssueMalformed},
		{"missing actor role", "{kind: report, actor: {session: s}, at: '2026-09-15T00:00:00Z'}", ErrV2MissingField},
		{"unknown actor role", "{kind: report, actor: {role: reviewer, session: s}, at: '2026-09-15T00:00:00Z'}", ErrV2IssueMalformed},
		{"missing actor session", "{kind: report, actor: {role: owner}, at: '2026-09-15T00:00:00Z'}", ErrV2MissingField},
		{"missing at", "{kind: report, actor: {role: owner, session: s}}", ErrV2MissingField},
		{"unparseable at", "{kind: report, actor: {role: owner, session: s}, at: yesterday}", ErrV2IssueMalformed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			overrides := map[string]string{"source": tt.source}
			if tt.source == "" {
				overrides = map[string]string{"source": "null"}
			}
			_, err := DecodeIssueV2("test.md", issueFixture(overrides))
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("DecodeIssueV2() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

// TestDecodeIssueV2_sourceReportAndMigration proves an Issue can name an
// origin other than a Check, so a directly reported or migrated Issue does
// not have to borrow an evaluation's authority to exist.
func TestDecodeIssueV2_sourceReportAndMigration(t *testing.T) {
	for _, kind := range []IssueOriginKind{IssueOriginReport, IssueOriginMigration} {
		source := "{kind: " + string(kind) + ", actor: {role: owner, session: owner-1}, at: '2026-09-15T00:00:00Z'}"
		issue, err := DecodeIssueV2("test.md", issueFixture(map[string]string{"source": source}))
		if err != nil {
			t.Fatalf("DecodeIssueV2() source kind %q error = %v", kind, err)
		}
		if issue.Origin.Kind != kind || issue.Origin.Check != "" {
			t.Errorf("Origin = %+v, want kind %q with no check", issue.Origin, kind)
		}
	}
}

func TestDecodeIssueV2_identityReferences(t *testing.T) {
	tests := []struct {
		name      string
		overrides map[string]string
	}{
		{"task reference with a check id", map[string]string{"tasks": "[C-001]"}},
		{"task reference with too few digits", map[string]string{"tasks": "[T01]"}},
		{"check reference with a task id", map[string]string{"checks": "[T-001]"}},
		{"check reference with an issue id", map[string]string{"checks": "[I-001]"}},
		{"duplicate_of with a task id", map[string]string{"duplicate_of": "T-001"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecodeIssueV2("test.md", issueFixture(tt.overrides))
			if !errors.Is(err, ErrV2InvalidID) {
				t.Fatalf("DecodeIssueV2() error = %v, want ErrV2InvalidID", err)
			}
		})
	}
}

// TestDecodeIssueV2_guardrailIDsAreOpaque proves policy IDs are recorded as
// given and never resolved against a Guardrails file: any shape of policy
// identifier decodes, because Savepoint does not own policy vocabulary.
func TestDecodeIssueV2_guardrailIDsAreOpaque(t *testing.T) {
	issue, err := DecodeIssueV2("test.md", issueFixture(map[string]string{
		"guardrail_ids": "[DATA-02, \"some.project/rule-7\", NOT-A-REAL-RULE]",
	}))
	if err != nil {
		t.Fatalf("DecodeIssueV2() error = %v", err)
	}
	want := []string{"DATA-02", "some.project/rule-7", "NOT-A-REAL-RULE"}
	if len(issue.GuardrailIDs) != len(want) {
		t.Fatalf("GuardrailIDs = %v, want %v", issue.GuardrailIDs, want)
	}
	for i, id := range want {
		if issue.GuardrailIDs[i] != id {
			t.Errorf("GuardrailIDs[%d] = %q, want %q", i, issue.GuardrailIDs[i], id)
		}
	}

	_, err = DecodeIssueV2("test.md", issueFixture(map[string]string{"guardrail_ids": "[DATA-02, \"  \"]"}))
	if !errors.Is(err, ErrV2IssueMalformed) {
		t.Fatalf("DecodeIssueV2() empty guardrail id error = %v, want ErrV2IssueMalformed", err)
	}
}

func TestDecodeIssueV2_blankSeverityRejected(t *testing.T) {
	_, err := DecodeIssueV2("test.md", issueFixture(map[string]string{"severity": "\"   \""}))
	if !errors.Is(err, ErrV2IssueMalformed) {
		t.Fatalf("DecodeIssueV2() error = %v, want ErrV2IssueMalformed", err)
	}
}

func TestDecodeIssueV2_resolutionShape(t *testing.T) {
	tests := []struct {
		name       string
		resolution string
		wantErr    error
	}{
		{"missing disposition", "{actor: {role: checker, session: s}, at: '2026-09-16T00:00:00Z'}", ErrV2MissingField},
		{"unknown disposition", "{disposition: closed, actor: {role: checker, session: s}, at: '2026-09-16T00:00:00Z'}", ErrV2IssueMalformed},
		{"malformed proof check", "{disposition: verified, check: I-001, actor: {role: checker, session: s}, at: '2026-09-16T00:00:00Z'}", ErrV2InvalidID},
		{"missing actor", "{disposition: verified, at: '2026-09-16T00:00:00Z'}", ErrV2MissingField},
		{"unknown actor role", "{disposition: verified, actor: {role: reviewer, session: s}, at: '2026-09-16T00:00:00Z'}", ErrV2IssueMalformed},
		{"missing at", "{disposition: verified, actor: {role: checker, session: s}}", ErrV2MissingField},
		{"unparseable at", "{disposition: verified, actor: {role: checker, session: s}, at: '2026-09-16'}", ErrV2IssueMalformed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecodeIssueV2("test.md", issueFixture(map[string]string{
				"status":     "resolved",
				"resolution": tt.resolution,
			}))
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("DecodeIssueV2() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

// TestDecodeIssueV2_resolutionDispositionsDecodeAsShape proves every
// disposition decodes here. The proof obligations each one carries are
// cross-record rules resolved against the full index, not decoding rules.
func TestDecodeIssueV2_resolutionDispositionsDecodeAsShape(t *testing.T) {
	for _, disposition := range []IssueDisposition{IssueDispositionVerified, IssueDispositionAccepted, IssueDispositionDuplicate, IssueDispositionEscalated} {
		resolution := "{disposition: " + string(disposition) + ", actor: {role: checker, session: s}, at: '2026-09-16T00:00:00Z', reason: closed}"
		issue, err := DecodeIssueV2("test.md", issueFixture(map[string]string{
			"status":     "resolved",
			"resolution": resolution,
		}))
		if err != nil {
			t.Fatalf("DecodeIssueV2() disposition %q error = %v", disposition, err)
		}
		if issue.Resolution == nil || issue.Resolution.Disposition != disposition {
			t.Errorf("Resolution = %+v, want disposition %q", issue.Resolution, disposition)
		}
	}
}

func TestDecodeIssueV2_historyEntryShapes(t *testing.T) {
	tests := []struct {
		name    string
		entry   string
		wantErr error
	}{
		{"missing kind", "{at: '2026-09-15T00:00:00Z', actor: {role: checker, session: s}}", ErrV2MissingField},
		{"unrecognized kind", "{at: '2026-09-15T00:00:00Z', actor: {role: checker, session: s}, kind: noted}", ErrV2IssueMalformed},
		{"missing at", "{actor: {role: checker, session: s}, kind: observed}", ErrV2MissingField},
		{"unparseable at", "{at: someday, actor: {role: checker, session: s}, kind: observed}", ErrV2IssueMalformed},
		{"missing actor role", "{at: '2026-09-15T00:00:00Z', actor: {session: s}, kind: observed}", ErrV2MissingField},
		{"unknown actor role", "{at: '2026-09-15T00:00:00Z', actor: {role: reviewer, session: s}, kind: observed}", ErrV2IssueMalformed},
		{"missing actor session", "{at: '2026-09-15T00:00:00Z', actor: {role: checker}, kind: observed}", ErrV2MissingField},
		{"malformed check reference", "{at: '2026-09-15T00:00:00Z', actor: {role: checker, session: s}, kind: rechecked, check: I-001}", ErrV2InvalidID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecodeIssueV2("test.md", issueFixture(map[string]string{"history": "[" + tt.entry + "]"}))
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("DecodeIssueV2() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

// TestDecodeIssueV2_historyDiagnosticNamesTheEntry proves a malformed entry
// in a long history names its own position, so a record with many entries
// does not have to be searched by hand.
func TestDecodeIssueV2_historyDiagnosticNamesTheEntry(t *testing.T) {
	history := "[" +
		"{at: '2026-09-15T00:00:00Z', actor: {role: checker, session: s}, kind: observed}, " +
		"{at: '2026-09-16T00:00:00Z', actor: {role: executor, session: s}, kind: sorted-out}]"

	_, err := DecodeIssueV2("issues/I-001-x.md", issueFixture(map[string]string{"history": history}))
	if !errors.Is(err, ErrV2IssueMalformed) {
		t.Fatalf("DecodeIssueV2() error = %v, want ErrV2IssueMalformed", err)
	}
	if !strings.Contains(err.Error(), "history[1]") {
		t.Errorf("DecodeIssueV2() error = %v, want the offending entry index named", err)
	}
}

func TestDecodeIssueV2_historyKinds(t *testing.T) {
	kinds := []IssueHistoryKind{
		IssueHistoryObserved, IssueHistoryRepairAttempted, IssueHistoryRechecked,
		IssueHistoryDeferred, IssueHistoryReopened, IssueHistoryOwnerDecision, IssueHistoryEscalated,
	}

	for _, kind := range kinds {
		entry := "{at: '2026-09-15T00:00:00Z', actor: {role: checker, session: s}, kind: " + string(kind) + "}"
		issue, err := DecodeIssueV2("test.md", issueFixture(map[string]string{"history": "[" + entry + "]"}))
		if err != nil {
			t.Fatalf("DecodeIssueV2() history kind %q error = %v", kind, err)
		}
		if len(issue.History) != 1 || issue.History[0].Kind != kind {
			t.Errorf("History = %+v, want one %q entry", issue.History, kind)
		}
	}
}

func TestDecodeIssueV2_malformedYAML(t *testing.T) {
	_, err := DecodeIssueV2("test.md", "---\nid: [broken\n---\n\n# Issue")
	if err == nil {
		t.Fatal("DecodeIssueV2() expected error for malformed YAML")
	}
}

func TestDecodeIssueV2_noFrontmatter(t *testing.T) {
	_, err := DecodeIssueV2("test.md", "# No frontmatter here")
	if !errors.Is(err, ErrNoFrontmatter) {
		t.Fatalf("DecodeIssueV2() error = %v, want ErrNoFrontmatter", err)
	}
}

// writeV2MixedIssueProject writes a project whose Issues span three statuses
// and three types, so one fixture exercises every listing and count.
func writeV2MixedIssueProject(t *testing.T, root string) {
	t.Helper()
	writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First objective")
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{fileName: "I-001-first.md", id: "I-001", issueType: "defect", status: "open"})
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{fileName: "I-002-second.md", id: "I-002", issueType: "drift", status: "in_progress"})
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{
		fileName: "I-003-third.md", id: "I-003", issueType: "defect", status: "resolved",
		resolution: "{disposition: accepted, actor: {role: owner, session: owner-1}, at: '2026-09-16T00:00:00Z', reason: accepted risk}",
	})
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{fileName: "I-004-fourth.md", id: "I-004", issueType: "guardrail", status: "open"})
}

// TestV2Index_issueListingsAndCounts proves the open, in_progress, and
// resolved sets, and the status and type counts, are answered from the index
// alone, with no stored summary anywhere in the project.
func TestV2Index_issueListingsAndCounts(t *testing.T) {
	root := t.TempDir()
	writeV2MixedIssueProject(t, root)

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}

	listings := map[IssueStatus][]string{
		IssueStatusOpen:       {"I-001", "I-004"},
		IssueStatusInProgress: {"I-002"},
		IssueStatusResolved:   {"I-003"},
	}
	for status, want := range listings {
		if got := index.IssueIDsWithStatus(status); !reflect.DeepEqual(got, want) {
			t.Errorf("IssueIDsWithStatus(%s) = %v, want %v", status, got, want)
		}
	}

	wantStatusCounts := map[IssueStatus]int{IssueStatusOpen: 2, IssueStatusInProgress: 1, IssueStatusResolved: 1}
	if got := index.IssueStatusCounts(); !reflect.DeepEqual(got, wantStatusCounts) {
		t.Errorf("IssueStatusCounts() = %v, want %v", got, wantStatusCounts)
	}

	wantTypeCounts := map[IssueType]int{
		IssueTypeDefect: 2, IssueTypeDrift: 1, IssueTypeGuardrail: 1,
		IssueTypeVerification: 0, IssueTypeOther: 0,
	}
	if got := index.IssueTypeCounts(); !reflect.DeepEqual(got, wantTypeCounts) {
		t.Errorf("IssueTypeCounts() = %v, want %v", got, wantTypeCounts)
	}
}

// TestV2Index_issueCountsCoverEveryVocabularyValue proves a project with no
// Issues still reports every status and type at zero, so a caller never has
// to read an absent key as none.
func TestV2Index_issueCountsCoverEveryVocabularyValue(t *testing.T) {
	index := &V2Index{Issues: map[string]*IssueV2{}}

	statusCounts := index.IssueStatusCounts()
	for _, status := range issueStatuses {
		if count, ok := statusCounts[status]; !ok || count != 0 {
			t.Errorf("IssueStatusCounts()[%s] = %d, ok = %v, want 0 and present", status, count, ok)
		}
	}

	typeCounts := index.IssueTypeCounts()
	for _, issueType := range issueTypes {
		if count, ok := typeCounts[issueType]; !ok || count != 0 {
			t.Errorf("IssueTypeCounts()[%s] = %d, ok = %v, want 0 and present", issueType, count, ok)
		}
	}

	if got := index.IssueIDsWithStatus(IssueStatusOpen); len(got) != 0 {
		t.Errorf("IssueIDsWithStatus(open) = %v, want empty", got)
	}
}

// TestV2Index_issueListingsAreDerivedNotPersisted proves a listing is computed
// on each call: reading one writes nothing to the project, and a change to the
// loaded records changes the next answer without any file being refreshed.
func TestV2Index_issueListingsAreDerivedNotPersisted(t *testing.T) {
	root := t.TempDir()
	writeV2MixedIssueProject(t, root)

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}

	before := projectFiles(t, root)
	index.IssueIDsWithStatus(IssueStatusOpen)
	index.IssueStatusCounts()
	index.IssueTypeCounts()
	if after := projectFiles(t, root); !reflect.DeepEqual(before, after) {
		t.Errorf("project files after listing = %v, want unchanged %v", after, before)
	}

	index.Issues["I-001"].Status = IssueStatusResolved
	want := []string{"I-004"}
	if got := index.IssueIDsWithStatus(IssueStatusOpen); !reflect.DeepEqual(got, want) {
		t.Errorf("IssueIDsWithStatus(open) = %v, want %v recomputed from the index", got, want)
	}
}

// projectFiles lists every file under root, so a test can prove a read-only
// operation left the project exactly as it found it.
func projectFiles(t *testing.T, root string) []string {
	t.Helper()
	var files []string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		info, err := os.Stat(path)
		if err != nil {
			return err
		}
		files = append(files, path+":"+info.ModTime().String())
		return nil
	})
	if err != nil {
		t.Fatalf("WalkDir(%s) error = %v", root, err)
	}
	return files
}

func mustVerifiedIssue(id, proofCheck string) *IssueV2 {
	return &IssueV2{
		ID:     id,
		Status: IssueStatusResolved,
		Checks: []string{proofCheck},
		Resolution: &IssueResolution{
			Disposition: IssueDispositionVerified,
			Check:       proofCheck,
			Actor:       Actor{Role: ActorRoleChecker, Session: "sess-1"},
		},
		Source: V2SourceDocument{Path: "issues/" + id + "-x.md"},
	}
}

func TestInspectIssueConsistency_ignoresUnresolvedAndNonVerifiedIssues(t *testing.T) {
	index := &V2Index{
		Issues: map[string]*IssueV2{
			"I-001": {ID: "I-001", Status: IssueStatusOpen},
			"I-002": {ID: "I-002", Status: IssueStatusResolved, Resolution: &IssueResolution{Disposition: IssueDispositionAccepted}},
		},
		Checks:      map[string]*CheckV2{},
		LatestCheck: map[string]string{},
	}

	if got := InspectIssueConsistency(index); len(got) != 0 {
		t.Fatalf("InspectIssueConsistency() = %+v, want none for open/accepted issues", got)
	}
}

func TestInspectIssueConsistency_verifiedProofStillLatestReportsNothing(t *testing.T) {
	index := &V2Index{
		Issues: map[string]*IssueV2{"I-001": mustVerifiedIssue("I-001", "C-001")},
		Checks: map[string]*CheckV2{
			"C-001": {ID: "C-001", Scope: CheckScope{Kind: CheckScopeTask, ID: "T-001"}, Result: CheckResultClear},
		},
		LatestCheck: map[string]string{"T-001": "C-001"},
	}

	if got := InspectIssueConsistency(index); len(got) != 0 {
		t.Fatalf("InspectIssueConsistency() = %+v, want none when the proof check is still latest", got)
	}
}

func TestInspectIssueConsistency_verifiedProofSuperseded(t *testing.T) {
	index := &V2Index{
		Issues: map[string]*IssueV2{"I-001": mustVerifiedIssue("I-001", "C-001")},
		Checks: map[string]*CheckV2{
			"C-001": {ID: "C-001", Scope: CheckScope{Kind: CheckScopeTask, ID: "T-001"}, Result: CheckResultClear},
			"C-002": {ID: "C-002", Scope: CheckScope{Kind: CheckScopeTask, ID: "T-001"}, Result: CheckResultClear, Supersedes: "C-001"},
		},
		LatestCheck: map[string]string{"T-001": "C-002"},
	}

	got := InspectIssueConsistency(index)
	if len(got) != 1 || got[0].Kind != IssueConsistencyProofSuperseded || got[0].Issue != "I-001" {
		t.Fatalf("InspectIssueConsistency() = %+v, want one IssueConsistencyProofSuperseded naming I-001", got)
	}
	if !strings.Contains(got[0].Detail, "C-001") || !strings.Contains(got[0].Detail, "C-002") {
		t.Errorf("Detail = %q, want both the proof and its superseder named", got[0].Detail)
	}
}

func TestInspectIssueConsistency_sortedOrderReturnsEveryProblem(t *testing.T) {
	index := &V2Index{
		Issues: map[string]*IssueV2{
			"I-002": mustVerifiedIssue("I-002", "C-001"),
			"I-001": mustVerifiedIssue("I-001", "C-001"),
		},
		Checks: map[string]*CheckV2{
			"C-001": {ID: "C-001", Scope: CheckScope{Kind: CheckScopeTask, ID: "T-001"}, Result: CheckResultClear},
			"C-002": {ID: "C-002", Scope: CheckScope{Kind: CheckScopeTask, ID: "T-001"}, Result: CheckResultClear, Supersedes: "C-001"},
		},
		LatestCheck: map[string]string{"T-001": "C-002"},
	}

	got := InspectIssueConsistency(index)
	if len(got) != 2 {
		t.Fatalf("InspectIssueConsistency() = %+v, want 2 problems", got)
	}
	if got[0].Issue != "I-001" || got[1].Issue != "I-002" {
		t.Fatalf("InspectIssueConsistency() order = [%s, %s], want sorted [I-001, I-002]", got[0].Issue, got[1].Issue)
	}
}

func TestAdvanceAndRetreatIssueV2_transitions(t *testing.T) {
	actor := Actor{Role: ActorRoleOwner, Session: "board-owner"}
	at := time.Date(2026, 9, 25, 2, 3, 4, 0, time.UTC)
	tests := []struct {
		name                   string
		status                 string
		fields                 string
		advance                bool
		wantStatus             IssueStatus
		wantKind               IssueHistoryKind
		wantDispositionInNote  string
		wantResolved           bool
		wantClearClosureFields bool
	}{
		{
			name:       "open to in progress",
			status:     "open",
			advance:    true,
			wantStatus: IssueStatusInProgress,
			wantKind:   IssueHistoryOwnerDecision,
		},
		{
			name:       "in progress to open",
			status:     "in_progress",
			wantStatus: IssueStatusOpen,
			wantKind:   IssueHistoryOwnerDecision,
		},
		{
			name:         "in progress to resolved",
			status:       "in_progress",
			advance:      true,
			wantStatus:   IssueStatusResolved,
			wantKind:     IssueHistoryOwnerDecision,
			wantResolved: true,
		},
		{
			name:   "verified issue reopens",
			status: "resolved",
			fields: `resolution:
  disposition: verified
  check: C-003
  actor: {role: checker, session: checker-session}
  at: '2026-09-24T10:00:00Z'
  reason: repaired`,
			wantStatus:             IssueStatusInProgress,
			wantKind:               IssueHistoryReopened,
			wantDispositionInNote:  "verified",
			wantClearClosureFields: true,
		},
		{
			name:   "duplicate issue reopens and clears duplicate target",
			status: "resolved",
			fields: `resolution:
  disposition: duplicate
  actor: {role: planner, session: planner-session}
  at: '2026-09-24T10:00:00Z'
  reason: duplicate report
duplicate_of: I-002`,
			wantStatus:             IssueStatusInProgress,
			wantKind:               IssueHistoryReopened,
			wantDispositionInNote:  "duplicate",
			wantClearClosureFields: true,
		},
		{
			name:   "escalated issue reopens and clears escalation target",
			status: "resolved",
			fields: `resolution:
  disposition: escalated
  actor: {role: planner, session: planner-session}
  at: '2026-09-24T10:00:00Z'
  reason: promoted
escalated_to: O-002`,
			wantStatus:             IssueStatusInProgress,
			wantKind:               IssueHistoryReopened,
			wantDispositionInNote:  "escalated",
			wantClearClosureFields: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issue, index, path := newIssueTransitionFixture(t, tt.status, tt.fields)
			beforeBody := issue.Source.Body
			beforeHistory := append([]IssueHistoryEntry(nil), issue.History...)
			if tt.advance {
				if err := AdvanceIssueV2(index, issue.ID, actor, at); err != nil {
					t.Fatalf("AdvanceIssueV2() error = %v", err)
				}
			} else if err := RetreatIssueV2(index, issue.ID, actor, at); err != nil {
				t.Fatalf("RetreatIssueV2() error = %v", err)
			}

			if issue.Status != tt.wantStatus {
				t.Errorf("status = %q, want %q", issue.Status, tt.wantStatus)
			}
			if len(issue.History) != len(beforeHistory)+1 {
				t.Fatalf("history length = %d, want %d", len(issue.History), len(beforeHistory)+1)
			}
			if !reflect.DeepEqual(issue.History[:len(beforeHistory)], beforeHistory) {
				t.Errorf("existing history changed: got %+v, want original entries %+v", issue.History[:len(beforeHistory)], beforeHistory)
			}
			entry := issue.History[len(issue.History)-1]
			if entry.Kind != tt.wantKind || entry.Actor != actor || !entry.At.Equal(at) {
				t.Errorf("appended history entry = %+v, want kind %q, actor %+v, time %s", entry, tt.wantKind, actor, at.Format(time.RFC3339))
			}
			if tt.wantDispositionInNote != "" && !strings.Contains(entry.Note, tt.wantDispositionInNote) {
				t.Errorf("reopened note = %q, want removed disposition %q named", entry.Note, tt.wantDispositionInNote)
			}
			if tt.wantResolved {
				if issue.Resolution == nil {
					t.Fatal("resolution = nil, want accepted owner resolution")
				}
				if issue.Resolution.Disposition != IssueDispositionAccepted || issue.Resolution.Actor != actor || !issue.Resolution.At.Equal(at) || issue.Resolution.Reason != issueBoardResolutionReason || issue.Resolution.Check != "" {
					t.Errorf("resolution = %+v, want accepted owner resolution without a proof Check", issue.Resolution)
				}
			} else if tt.wantClearClosureFields {
				if issue.Resolution != nil || issue.DuplicateOf != "" || issue.EscalatedTo != "" {
					t.Errorf("reopened closure fields = resolution:%+v duplicate_of:%q escalated_to:%q, want all cleared", issue.Resolution, issue.DuplicateOf, issue.EscalatedTo)
				}
			}

			written, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			updated, err := DecodeIssueV2(path, string(written))
			if err != nil {
				t.Fatalf("DecodeIssueV2(written) error = %v", err)
			}
			if updated.Source.Body != beforeBody {
				t.Errorf("Markdown body changed: got %q, want %q", updated.Source.Body, beforeBody)
			}
			if !strings.Contains(string(written), "project_extension") || !strings.Contains(string(written), "preserve-this-value") {
				t.Errorf("unknown frontmatter field was not preserved:\n%s", written)
			}
			if !strings.Contains(string(written), "history_extension") {
				t.Errorf("unknown existing history field was not preserved:\n%s", written)
			}
			if !reflect.DeepEqual(updated.History[:len(beforeHistory)], beforeHistory) {
				t.Errorf("written history changed existing entries: got %+v, want %+v", updated.History[:len(beforeHistory)], beforeHistory)
			}
			if len(updated.History) != len(beforeHistory)+1 || updated.History[len(updated.History)-1].Kind != tt.wantKind {
				t.Errorf("written history = %+v, want one appended %q entry", updated.History, tt.wantKind)
			}
			if tt.wantClearClosureFields && (strings.Contains(string(written), "duplicate_of:") || strings.Contains(string(written), "escalated_to:")) {
				t.Errorf("reopened file retained disposition targets:\n%s", written)
			}
		})
	}
}

func TestAdvanceAndRetreatIssueV2_refusalsAndUnknownIDDoNotWrite(t *testing.T) {
	actor := Actor{Role: ActorRoleOwner, Session: "board-owner"}
	at := time.Date(2026, 9, 25, 2, 3, 4, 0, time.UTC)
	tests := []struct {
		name   string
		status string
		fields string
		call   func(*V2Index, string) error
		want   error
	}{
		{
			name:   "advance resolved",
			status: "resolved",
			fields: `resolution:
  disposition: accepted
  actor: {role: owner, session: prior-owner}
  at: '2026-09-24T10:00:00Z'
  reason: previously accepted`,
			call: func(index *V2Index, id string) error { return AdvanceIssueV2(index, id, actor, at) },
			want: ErrV2IssueCannotAdvance,
		},
		{
			name:   "retreat open",
			status: "open",
			call:   func(index *V2Index, id string) error { return RetreatIssueV2(index, id, actor, at) },
			want:   ErrV2IssueCannotRetreat,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issue, index, path := newIssueTransitionFixture(t, tt.status, tt.fields)
			before, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if err := tt.call(index, issue.ID); !errors.Is(err, tt.want) {
				t.Fatalf("transition error = %v, want %v", err, tt.want)
			}
			after, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(after, before) {
				t.Errorf("refused transition changed file:\n%s", after)
			}
		})
	}

	if err := AdvanceIssueV2(&V2Index{Issues: map[string]*IssueV2{}}, "I-999", actor, at); !errors.Is(err, ErrV2IssueNotFound) {
		t.Errorf("unknown Issue error = %v, want %v", err, ErrV2IssueNotFound)
	}
}

func TestAdvanceIssueV2_refusesStaleIssueFile(t *testing.T) {
	actor := Actor{Role: ActorRoleOwner, Session: "board-owner"}
	at := time.Date(2026, 9, 25, 2, 3, 4, 0, time.UTC)
	issue, index, path := newIssueTransitionFixture(t, "open", "")
	external := []byte("external edit that must survive")
	if err := os.WriteFile(path, external, 0644); err != nil {
		t.Fatal(err)
	}

	err := AdvanceIssueV2(index, issue.ID, actor, at)
	if !errors.Is(err, ErrV2SourceConflict) {
		t.Fatalf("AdvanceIssueV2() error = %v, want ErrV2SourceConflict", err)
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(after, external) {
		t.Errorf("stale write replaced external file content: %q", after)
	}
}

func TestAdvanceIssueV2_preservesExistingHistoryBytes(t *testing.T) {
	actor := Actor{Role: ActorRoleOwner, Session: "board-owner"}
	at := time.Date(2026, 9, 25, 2, 3, 4, 0, time.UTC)
	issue, index, path := newIssueTransitionFixture(t, "open", "")
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	const existingHistory = `history:
  - at: '2026-09-21T00:00:00Z'
    actor: {role: checker, session: existing-checker}
    kind: observed
    note: original history entry
    history_extension: preserve-this-value`
	if !strings.Contains(string(before), existingHistory) {
		t.Fatalf("fixture history does not contain the expected 2-space source block:\n%s", before)
	}

	if err := AdvanceIssueV2(index, issue.ID, actor, at); err != nil {
		t.Fatalf("AdvanceIssueV2() error = %v", err)
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(after), existingHistory) {
		t.Errorf("existing history bytes changed during Issue transition:\n%s", after)
	}
	if !strings.Contains(string(after), "note: Moved from open to in_progress by the owner from the board.") {
		t.Errorf("new history entry is missing:\n%s", after)
	}
}

// TestIssueTransitionsKeepEarlierFrontmatterBytes runs a full move sequence
// on a 4-space Issue that already holds a board-written block-style entry and
// a wrapped plain note. Every move must keep all earlier frontmatter lines
// byte-for-byte, apart from the status line and a removed resolution block.
func TestIssueTransitionsKeepEarlierFrontmatterBytes(t *testing.T) {
	actor := Actor{Role: ActorRoleOwner, Session: "board-owner"}
	at := time.Date(2026, 9, 25, 2, 3, 4, 0, time.UTC)
	path := filepath.Join(t.TempDir(), "I-001-bytes.md")
	const body = "\n\n# Body\n\nKeep this body.\n"
	content := `---
id: I-001
title: Issue byte preservation
type: defect
status: open
source:
    kind: check
    check: C-001
    actor: {role: checker, session: checker-1}
    at: '2026-09-23T01:42:44Z'
checks: [C-001]
history:
    - at: '2026-09-23T01:42:44Z'
      actor: {role: checker, session: checker-1}
      kind: observed
      check: C-001
      note: Owner reported a RELOAD error decoding router.md because of a YAML
        mapping-values parse error on line 5. Exact diagnostic is recorded in the
        Evidence section.
    - at: "2026-09-25T09:26:40Z"
      actor:
        role: owner
        session: board-owner
      kind: owner_decision
      note: Moved from open to in_progress by the owner from the board.
---` + body
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	issue, err := DecodeIssueV2(path, content)
	if err != nil {
		t.Fatalf("DecodeIssueV2() error = %v", err)
	}
	index := &V2Index{Issues: map[string]*IssueV2{issue.ID: issue}}

	moves := []struct {
		name    string
		advance bool
		status  string
	}{
		{"open to in_progress", true, "in_progress"},
		{"in_progress to resolved", true, "resolved"},
		{"reopen", false, "in_progress"},
		{"in_progress to open", false, "open"},
	}
	for i, move := range moves {
		before, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if i == 2 {
			// Reload from disk so the move also runs on a freshly parsed record.
			issue, err = DecodeIssueV2(path, string(before))
			if err != nil {
				t.Fatalf("reload: %v", err)
			}
			index.Issues[issue.ID] = issue
		}
		moveAt := at.Add(time.Duration(i) * time.Minute)
		if move.advance {
			err = AdvanceIssueV2(index, issue.ID, actor, moveAt)
		} else {
			err = RetreatIssueV2(index, issue.ID, actor, moveAt)
		}
		if err != nil {
			t.Fatalf("%s: %v", move.name, err)
		}
		after, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}

		kept := string(before)
		kept = kept[:strings.Index(kept, "\n---\n")+1]
		if cut := strings.Index(kept, "\nresolution:\n"); cut >= 0 {
			kept = kept[:cut+1]
		}
		kept = regexp.MustCompile(`(?m)^status: .*$`).ReplaceAllString(kept, "status: "+move.status)
		if !strings.HasPrefix(string(after), kept) {
			t.Fatalf("%s changed earlier frontmatter bytes.\nwant prefix:\n%s\ngot:\n%s", move.name, kept, after)
		}
		if !strings.HasSuffix(string(after), "\n---"+body) {
			t.Fatalf("%s changed the body:\n%s", move.name, after)
		}
		if _, err := DecodeIssueV2(path, string(after)); err != nil {
			t.Fatalf("%s wrote an undecodable record: %v", move.name, err)
		}
	}
}

func newIssueTransitionFixture(t *testing.T, status, extraFields string) (*IssueV2, *V2Index, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "I-001-transition.md")
	content := `---
id: I-001
title: "Issue transition test"
type: defect
status: ` + status + `
source: {kind: report, actor: {role: owner, session: reporter}, at: '2026-09-20T00:00:00Z'}
project_extension: {source: test, value: preserve-this-value}
`
	if extraFields != "" {
		content += extraFields + "\n"
	}
	content += `history:
  - at: '2026-09-21T00:00:00Z'
    actor: {role: checker, session: existing-checker}
    kind: observed
    note: original history entry
    history_extension: preserve-this-value
---

# Owner-authored Issue body

Keep this exact body when writing.
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	issue, err := DecodeIssueV2(path, content)
	if err != nil {
		t.Fatalf("DecodeIssueV2() error = %v", err)
	}
	return issue, &V2Index{Issues: map[string]*IssueV2{issue.ID: issue}}, path
}
