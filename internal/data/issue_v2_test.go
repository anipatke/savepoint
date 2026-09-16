package data

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// issueFixture builds a minimal valid Issue record so each test can override
// exactly the one field it is about, instead of restating the whole schema.
func issueFixture(overrides map[string]string) string {
	fields := map[string]string{
		"id":     "I001",
		"title":  "\"Follow-up that outlives one evaluation\"",
		"type":   "defect",
		"status": "open",
		"source": "{kind: check, check: C001, actor: {role: checker, session: sess-1}, at: '2026-09-15T00:00:00Z'}",
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
id: I001
title: "Router drifts from its documented state table"
type: drift
status: resolved
source:
  kind: check
  check: C001
  actor: {role: checker, session: sess-1}
  at: '2026-09-15T00:00:00Z'
tasks: [T001, T002]
checks: [C001, C002]
guardrail_ids: [DATA-02, STYLE-07]
severity: high
resolution:
  disposition: verified
  check: C002
  actor: {role: checker, session: sess-2}
  at: '2026-09-16T00:00:00Z'
  reason: repaired and rechecked
duplicate_of: I002
history:
  - {at: '2026-09-15T00:00:00Z', actor: {role: checker, session: sess-1}, kind: observed, note: first seen}
  - {at: '2026-09-16T00:00:00Z', actor: {role: checker, session: sess-2}, kind: rechecked, check: C002}
---

# Issue`

	issue, err := DecodeIssueV2("issues/I001-router-drift.md", content)
	if err != nil {
		t.Fatalf("DecodeIssueV2() error = %v", err)
	}

	if issue.ID != "I001" {
		t.Errorf("ID = %q, want I001", issue.ID)
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
		Check: "C001",
		Actor: Actor{Role: ActorRoleChecker, Session: "sess-1"},
		At:    mustParseTime(t, "2026-09-15T00:00:00Z"),
	}
	if issue.Origin != wantOrigin {
		t.Errorf("Origin = %+v, want %+v", issue.Origin, wantOrigin)
	}

	if len(issue.Tasks) != 2 || issue.Tasks[0] != "T001" || issue.Tasks[1] != "T002" {
		t.Errorf("Tasks = %v, want [T001 T002]", issue.Tasks)
	}
	if len(issue.Checks) != 2 || issue.Checks[0] != "C001" || issue.Checks[1] != "C002" {
		t.Errorf("Checks = %v, want [C001 C002]", issue.Checks)
	}
	if len(issue.GuardrailIDs) != 2 || issue.GuardrailIDs[0] != "DATA-02" {
		t.Errorf("GuardrailIDs = %v, want [DATA-02 STYLE-07]", issue.GuardrailIDs)
	}
	if issue.Severity != "high" {
		t.Errorf("Severity = %q, want high", issue.Severity)
	}
	if issue.DuplicateOf != "I002" {
		t.Errorf("DuplicateOf = %q, want I002", issue.DuplicateOf)
	}

	if issue.Resolution == nil {
		t.Fatal("Resolution = nil, want the decoded resolution block")
	}
	if issue.Resolution.Disposition != IssueDispositionVerified || issue.Resolution.Check != "C002" {
		t.Errorf("Resolution = %+v, want verified proved by C002", issue.Resolution)
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
	if issue.History[1].Kind != IssueHistoryRechecked || issue.History[1].Check != "C002" {
		t.Errorf("History[1] = %+v, want the rechecked entry naming C002", issue.History[1])
	}
	if !issue.History[0].At.Before(issue.History[1].At) {
		t.Errorf("History order = %v then %v, want recorded order preserved", issue.History[0].At, issue.History[1].At)
	}
}

// TestDecodeIssueV2_absentOptionalsStayAbsent proves an optional field that
// was never recorded decodes as absent rather than as a permissive default
// that would read as a deliberate declaration.
func TestDecodeIssueV2_absentOptionalsStayAbsent(t *testing.T) {
	issue, err := DecodeIssueV2("issues/I001-minimal.md", issueFixture(nil))
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
		{"wrong prefix letter", "C001"},
		{"lowercase prefix", "i001"},
		{"trailing garbage", "I001x"},
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
		{"missing kind", "{check: C001, actor: {role: checker, session: s}, at: '2026-09-15T00:00:00Z'}", ErrV2MissingField},
		{"unknown kind", "{kind: audit, actor: {role: checker, session: s}, at: '2026-09-15T00:00:00Z'}", ErrV2IssueMalformed},
		{"check kind without a check", "{kind: check, actor: {role: checker, session: s}, at: '2026-09-15T00:00:00Z'}", ErrV2MissingField},
		{"check kind with a malformed check", "{kind: check, check: T001, actor: {role: checker, session: s}, at: '2026-09-15T00:00:00Z'}", ErrV2InvalidID},
		{"report kind naming a check", "{kind: report, check: C001, actor: {role: owner, session: s}, at: '2026-09-15T00:00:00Z'}", ErrV2IssueMalformed},
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
		{"task reference with a check id", map[string]string{"tasks": "[C001]"}},
		{"task reference with too few digits", map[string]string{"tasks": "[T01]"}},
		{"check reference with a task id", map[string]string{"checks": "[T001]"}},
		{"check reference with an issue id", map[string]string{"checks": "[I001]"}},
		{"duplicate_of with a task id", map[string]string{"duplicate_of": "T001"}},
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
		{"malformed proof check", "{disposition: verified, check: I001, actor: {role: checker, session: s}, at: '2026-09-16T00:00:00Z'}", ErrV2InvalidID},
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
	for _, disposition := range []IssueDisposition{IssueDispositionVerified, IssueDispositionAccepted, IssueDispositionDuplicate} {
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
		{"malformed check reference", "{at: '2026-09-15T00:00:00Z', actor: {role: checker, session: s}, kind: rechecked, check: I001}", ErrV2InvalidID},
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

	_, err := DecodeIssueV2("issues/I001-x.md", issueFixture(map[string]string{"history": history}))
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
		IssueHistoryDeferred, IssueHistoryReopened, IssueHistoryOwnerDecision,
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
	writeV2ObjectiveFixture(t, root, "O001-first", "O001", "First objective")
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{fileName: "I001-first.md", id: "I001", issueType: "defect", status: "open"})
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{fileName: "I002-second.md", id: "I002", issueType: "drift", status: "in_progress"})
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{
		fileName: "I003-third.md", id: "I003", issueType: "defect", status: "resolved",
		resolution: "{disposition: accepted, actor: {role: owner, session: owner-1}, at: '2026-09-16T00:00:00Z', reason: accepted risk}",
	})
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{fileName: "I004-fourth.md", id: "I004", issueType: "guardrail", status: "open"})
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
		IssueStatusOpen:       {"I001", "I004"},
		IssueStatusInProgress: {"I002"},
		IssueStatusResolved:   {"I003"},
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

	index.Issues["I001"].Status = IssueStatusResolved
	want := []string{"I004"}
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
