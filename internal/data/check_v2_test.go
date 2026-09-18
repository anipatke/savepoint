package data

import (
	"errors"
	"testing"
	"time"
)

func TestDecodeCheckV2_valid(t *testing.T) {
	content := `---
id: C001
scope: {kind: task, id: T001}
result: CLEAR
checked_by: {role: checker, session: review-001}
executed_session: build-001
checked_at: '2026-09-14T00:00:00Z'
reviewed:
  base_commit: abc123
  head_commit: def456
  files: [internal/data/check_v2.go]
  dependencies: [go.mod]
issues: [I001]
supersedes: C000
---

# Check`

	check, err := DecodeCheckV2("checks/C001-clear.md", content)
	if err != nil {
		t.Fatalf("DecodeCheckV2() error = %v", err)
	}

	if check.ID != "C001" {
		t.Errorf("ID = %q, want C001", check.ID)
	}
	if check.Scope != (CheckScope{Kind: CheckScopeTask, ID: "T001"}) {
		t.Errorf("Scope = %+v, want {task T001}", check.Scope)
	}
	if check.Result != CheckResultClear {
		t.Errorf("Result = %q, want CLEAR", check.Result)
	}
	if check.CheckedBy != (Actor{Role: ActorRoleChecker, Session: "review-001"}) {
		t.Errorf("CheckedBy = %+v, want {checker review-001}", check.CheckedBy)
	}
	if check.ExecutedSession != "build-001" {
		t.Errorf("ExecutedSession = %q, want build-001", check.ExecutedSession)
	}
	wantTime, _ := time.Parse(time.RFC3339, "2026-09-14T00:00:00Z")
	if !check.CheckedAt.Equal(wantTime) {
		t.Errorf("CheckedAt = %v, want %v", check.CheckedAt, wantTime)
	}
	if check.Reviewed == nil {
		t.Fatal("Reviewed = nil, want populated basis")
	}
	if check.Reviewed.BaseCommit != "abc123" || check.Reviewed.HeadCommit != "def456" {
		t.Errorf("Reviewed commits = %q/%q, want abc123/def456", check.Reviewed.BaseCommit, check.Reviewed.HeadCommit)
	}
	if len(check.Reviewed.Files) != 1 || check.Reviewed.Files[0] != "internal/data/check_v2.go" {
		t.Errorf("Reviewed.Files = %v, want [internal/data/check_v2.go]", check.Reviewed.Files)
	}
	if len(check.Reviewed.Dependencies) != 1 || check.Reviewed.Dependencies[0] != "go.mod" {
		t.Errorf("Reviewed.Dependencies = %v, want [go.mod]", check.Reviewed.Dependencies)
	}
	if len(check.Issues) != 1 || check.Issues[0] != "I001" {
		t.Errorf("Issues = %v, want [I001]", check.Issues)
	}
	if check.Supersedes != "C000" {
		t.Errorf("Supersedes = %q, want C000", check.Supersedes)
	}
}

func TestDecodeCheckV2_minimalValid(t *testing.T) {
	content := `---
id: C002
scope: {kind: objective, id: O001}
result: NEEDS WORK
checked_by: {role: owner, session: sess-1}
executed_session: build-001
checked_at: '2026-09-14T00:00:00Z'
---

# Check`

	check, err := DecodeCheckV2("test.md", content)
	if err != nil {
		t.Fatalf("DecodeCheckV2() error = %v", err)
	}
	if check.Scope != (CheckScope{Kind: CheckScopeObjective, ID: "O001"}) {
		t.Errorf("Scope = %+v, want {objective O001}", check.Scope)
	}
	if check.Result != CheckResultNeedsWork {
		t.Errorf("Result = %q, want NEEDS WORK", check.Result)
	}
	if check.Reviewed != nil {
		t.Errorf("Reviewed = %+v, want nil for an absent block, not a passing default", check.Reviewed)
	}
	if len(check.Issues) != 0 {
		t.Errorf("Issues = %v, want empty", check.Issues)
	}
	if check.Supersedes != "" {
		t.Errorf("Supersedes = %q, want empty", check.Supersedes)
	}
}

func TestDecodeCheckV2_executedSessionValidation(t *testing.T) {
	tests := []struct {
		name            string
		executedSession string
		checkedSession  string
		wantErr         error
	}{
		{name: "missing", executedSession: "", checkedSession: "review-001", wantErr: ErrV2MissingField},
		{name: "blank", executedSession: "   ", checkedSession: "review-001", wantErr: ErrV2MissingField},
		{name: "malformed shape", executedSession: "{id: build-001}", checkedSession: "review-001", wantErr: ErrV2Malformed},
		{name: "same as checker", executedSession: "review-001", checkedSession: "review-001", wantErr: ErrV2CheckMalformed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executedSession := ""
			if tt.executedSession != "" {
				executedSession = "executed_session: " + tt.executedSession + "\n"
			}
			content := "---\nid: C001\nscope: {kind: task, id: T001}\nresult: CLEAR\nchecked_by: {role: checker, session: " + tt.checkedSession + "}\n" + executedSession + "checked_at: '2026-09-14T00:00:00Z'\n---\n\n# Check"
			_, err := DecodeCheckV2("test.md", content)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("DecodeCheckV2() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecodeCheckV2_reviewedEmptyEntriesStayAbsent(t *testing.T) {
	content := `---
id: C002
scope: {kind: task, id: T001}
result: CLEAR
checked_by: {role: checker, session: sess-1}
executed_session: build-001
checked_at: '2026-09-14T00:00:00Z'
reviewed:
  files: []
  dependencies: []
---

# Check`

	check, err := DecodeCheckV2("test.md", content)
	if err != nil {
		t.Fatalf("DecodeCheckV2() error = %v", err)
	}
	if check.Reviewed == nil {
		t.Fatal("Reviewed = nil, want a present-but-empty basis")
	}
	if check.Reviewed.BaseCommit != "" || check.Reviewed.HeadCommit != "" {
		t.Errorf("Reviewed commits = %q/%q, want both empty", check.Reviewed.BaseCommit, check.Reviewed.HeadCommit)
	}
	if len(check.Reviewed.Files) != 0 || len(check.Reviewed.Dependencies) != 0 {
		t.Errorf("Reviewed.Files/Dependencies = %v/%v, want both empty", check.Reviewed.Files, check.Reviewed.Dependencies)
	}
}

func TestDecodeCheckV2_malformedID(t *testing.T) {
	tests := []struct {
		name string
		id   string
	}{
		{"missing digits", "C"},
		{"too few digits", "C01"},
		{"wrong prefix letter", "T001"},
		{"lowercase prefix", "c001"},
		{"trailing garbage", "C001x"},
		{"empty", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := "---\nid: \"" + tt.id + "\"\nscope: {kind: task, id: T001}\nresult: CLEAR\nchecked_by: {role: checker, session: s}\nexecuted_session: build-001\nchecked_at: '2026-09-14T00:00:00Z'\n---\n\n# Check"
			_, err := DecodeCheckV2("test.md", content)
			if !errors.Is(err, ErrV2InvalidID) {
				t.Fatalf("DecodeCheckV2() error = %v, want ErrV2InvalidID", err)
			}
		})
	}
}

func TestDecodeCheckV2_scope(t *testing.T) {
	tests := []struct {
		name    string
		scope   string
		wantErr error
	}{
		{"missing kind", "{id: T001}", ErrV2MissingField},
		{"invalid kind", "{kind: epic, id: T001}", ErrV2CheckMalformed},
		{"missing id", "{kind: task}", ErrV2MissingField},
		{"task kind with objective id", "{kind: task, id: O001}", ErrV2InvalidID},
		{"objective kind with task id", "{kind: objective, id: T001}", ErrV2InvalidID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := "---\nid: C001\nscope: " + tt.scope + "\nresult: CLEAR\nchecked_by: {role: checker, session: s}\nexecuted_session: build-001\nchecked_at: '2026-09-14T00:00:00Z'\n---\n\n# Check"
			_, err := DecodeCheckV2("test.md", content)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("DecodeCheckV2() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecodeCheckV2_result(t *testing.T) {
	tests := []struct {
		name    string
		result  string
		wantErr error
	}{
		{"missing", "", ErrV2MissingField},
		{"garbage", "MAYBE", ErrV2CheckMalformed},
		{"legacy lowercase", "clear", ErrV2CheckMalformed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := "---\nid: C001\nscope: {kind: task, id: T001}\nresult: \"" + tt.result + "\"\nchecked_by: {role: checker, session: s}\nexecuted_session: build-001\nchecked_at: '2026-09-14T00:00:00Z'\n---\n\n# Check"
			_, err := DecodeCheckV2("test.md", content)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("DecodeCheckV2() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecodeCheckV2_checkedBy(t *testing.T) {
	tests := []struct {
		name      string
		checkedBy string
		wantErr   error
	}{
		{"missing role", "{session: s}", ErrV2MissingField},
		{"invalid role", "{role: reviewer, session: s}", ErrV2CheckMalformed},
		{"missing session", "{role: checker}", ErrV2MissingField},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := "---\nid: C001\nscope: {kind: task, id: T001}\nresult: CLEAR\nchecked_by: " + tt.checkedBy + "\nexecuted_session: build-001\nchecked_at: '2026-09-14T00:00:00Z'\n---\n\n# Check"
			_, err := DecodeCheckV2("test.md", content)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("DecodeCheckV2() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecodeCheckV2_clearRequiresCheckerProvenance(t *testing.T) {
	content := `---
id: C001
scope: {kind: task, id: T001}
result: CLEAR
checked_by: {role: executor, session: executor-1}
executed_session: build-001
checked_at: '2026-09-14T00:00:00Z'
---

# Check`

	_, err := DecodeCheckV2("test.md", content)
	if !errors.Is(err, ErrV2CheckMalformed) {
		t.Fatalf("DecodeCheckV2() error = %v, want ErrV2CheckMalformed", err)
	}
}

func TestDecodeCheckV2_checkedAt(t *testing.T) {
	tests := []struct {
		name      string
		checkedAt string
		wantErr   error
	}{
		{"missing", "", ErrV2MissingField},
		{"not a timestamp", "yesterday", ErrV2CheckMalformed},
		{"date only, not RFC3339", "2026-09-14", ErrV2CheckMalformed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := "---\nid: C001\nscope: {kind: task, id: T001}\nresult: CLEAR\nchecked_by: {role: checker, session: s}\nexecuted_session: build-001\nchecked_at: \"" + tt.checkedAt + "\"\n---\n\n# Check"
			_, err := DecodeCheckV2("test.md", content)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("DecodeCheckV2() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecodeCheckV2_invalidIssueReference(t *testing.T) {
	content := `---
id: C001
scope: {kind: task, id: T001}
result: CLEAR
checked_by: {role: checker, session: s}
executed_session: build-001
checked_at: '2026-09-14T00:00:00Z'
issues: [T001]
---

# Check`

	_, err := DecodeCheckV2("test.md", content)
	if !errors.Is(err, ErrV2InvalidID) {
		t.Fatalf("DecodeCheckV2() error = %v, want ErrV2InvalidID", err)
	}
}

func TestDecodeCheckV2_invalidSupersedes(t *testing.T) {
	content := `---
id: C001
scope: {kind: task, id: T001}
result: CLEAR
checked_by: {role: checker, session: s}
executed_session: build-001
checked_at: '2026-09-14T00:00:00Z'
supersedes: T001
---

# Check`

	_, err := DecodeCheckV2("test.md", content)
	if !errors.Is(err, ErrV2InvalidID) {
		t.Fatalf("DecodeCheckV2() error = %v, want ErrV2InvalidID", err)
	}
}

func TestDecodeCheckV2_malformedYAML(t *testing.T) {
	content := `---
id: [broken
---

# Check`

	_, err := DecodeCheckV2("test.md", content)
	if err == nil {
		t.Fatal("DecodeCheckV2() expected error for malformed YAML")
	}
}

func TestDecodeCheckV2_noFrontmatter(t *testing.T) {
	_, err := DecodeCheckV2("test.md", "# No frontmatter here")
	if !errors.Is(err, ErrNoFrontmatter) {
		t.Fatalf("DecodeCheckV2() error = %v, want ErrNoFrontmatter", err)
	}
}
