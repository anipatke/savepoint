package data

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/opencode/savepoint/internal/testutil"
)

func validCarryFrontmatter(check string) decisionCarryV2Frontmatter {
	applies := true
	return decisionCarryV2Frontmatter{
		Check:      check,
		Applies:    &applies,
		AssessedBy: evidenceActorFrontmatter{Role: "checker", Session: "check-1"},
		AssessedAt: "2026-10-10T00:00:00Z",
		Reason:     "only Design text changed",
	}
}

func ownerValidationWithCarry(carried ...decisionCarryV2Frontmatter) evidenceV2Frontmatter {
	return evidenceV2Frontmatter{OwnerValidation: &ownerValidationV2Frontmatter{
		Required:       true,
		AcceptedCheck:  "C-002",
		AcceptedBy:     &evidenceActorFrontmatter{Role: "owner", Session: "owner-1"},
		Scope:          []string{"AC-1", "AC-2"},
		CarriedForward: carried,
	}}
}

func TestDecodeEvidenceV2_ownerValidationScopeAndCarriedForward(t *testing.T) {
	notApplies := false
	changed := validCarryFrontmatter("C-004")
	changed.Applies = &notApplies
	changed.MaterialChange = "AC-2 behavior changed"
	renewed := validCarryFrontmatter("C-004")
	renewed.AssessedBy = evidenceActorFrontmatter{Role: "owner", Session: "owner-2"}

	evidence, err := decodeEvidenceV2("test.md", "task", "T-001", ownerValidationWithCarry(validCarryFrontmatter("C-003"), changed, renewed))
	if err != nil {
		t.Fatalf("decodeEvidenceV2() error = %v", err)
	}
	ov := evidence.OwnerValidation
	if len(ov.Scope) != 2 || ov.Scope[0] != "AC-1" {
		t.Errorf("Scope = %v, want [AC-1 AC-2]", ov.Scope)
	}
	if len(ov.CarriedForward) != 3 {
		t.Fatalf("CarriedForward = %+v, want 3 entries", ov.CarriedForward)
	}
	first := ov.CarriedForward[0]
	wantAt := time.Date(2026, 10, 10, 0, 0, 0, 0, time.UTC)
	if first.Check != "C-003" || !first.Applies || first.AssessedBy != (Actor{Role: ActorRoleChecker, Session: "check-1"}) ||
		!first.AssessedAt.Equal(wantAt) || first.Reason != "only Design text changed" {
		t.Errorf("CarriedForward[0] = %+v, want decoded checker entry", first)
	}
	if second := ov.CarriedForward[1]; second.Applies || second.MaterialChange != "AC-2 behavior changed" {
		t.Errorf("CarriedForward[1] = %+v, want not-applying entry naming the change", second)
	}
	if third := ov.CarriedForward[2]; third.AssessedBy.Role != ActorRoleOwner || !third.Applies {
		t.Errorf("CarriedForward[2] = %+v, want applying owner renewal", third)
	}
}

func TestDecodeEvidenceV2_exceptionCarriedForward(t *testing.T) {
	exception := validExceptionFrontmatter()
	exception.CarriedForward = []decisionCarryV2Frontmatter{validCarryFrontmatter("C-004")}
	evidence, err := decodeEvidenceV2("test.md", "objective", "O-001", evidenceV2Frontmatter{Exception: &exception})
	if err != nil {
		t.Fatalf("decodeEvidenceV2() error = %v", err)
	}
	if got := evidence.Exception.CarriedForward; len(got) != 1 || got[0].Check != "C-004" {
		t.Errorf("Exception.CarriedForward = %+v, want one C-004 entry", got)
	}
}

func TestDecodeEvidenceV2_decisionsWithoutNewFieldsDecodeUnchanged(t *testing.T) {
	ov := ownerValidationV2Frontmatter{Required: true, AcceptedCheck: "C-002", AcceptedBy: &evidenceActorFrontmatter{Role: "owner", Session: "owner-1"}}
	exception := validExceptionFrontmatter()
	evidence, err := decodeEvidenceV2("test.md", "task", "T-001", evidenceV2Frontmatter{OwnerValidation: &ov, Exception: &exception})
	if err != nil {
		t.Fatalf("decodeEvidenceV2() error = %v", err)
	}
	if evidence.OwnerValidation.Scope != nil || evidence.OwnerValidation.CarriedForward != nil || evidence.Exception.CarriedForward != nil {
		t.Errorf("evidence = %+v, want no scope or carried entries", evidence)
	}
}

func TestDecodeEvidenceV2_carriedForwardRejections(t *testing.T) {
	notApplies := false
	tests := []struct {
		name    string
		mutate  func(*decisionCarryV2Frontmatter)
		wantErr error
		wantMsg string
	}{
		{"missing check", func(c *decisionCarryV2Frontmatter) { c.Check = "" }, ErrV2MissingField, "carried_forward[0].check"},
		{"malformed check", func(c *decisionCarryV2Frontmatter) { c.Check = "C3" }, ErrV2InvalidID, "carried_forward[0].check"},
		{"originating check", func(c *decisionCarryV2Frontmatter) { c.Check = "C-002" }, ErrV2EvidenceMalformed, "own check"},
		{"missing applies", func(c *decisionCarryV2Frontmatter) { c.Applies = nil }, ErrV2MissingField, "carried_forward[0].applies"},
		{"missing assessor role", func(c *decisionCarryV2Frontmatter) { c.AssessedBy.Role = "" }, ErrV2MissingField, "assessed_by.role"},
		{"executor assessor", func(c *decisionCarryV2Frontmatter) { c.AssessedBy.Role = "executor" }, ErrV2EvidenceMalformed, "use checker or owner"},
		{"missing assessor session", func(c *decisionCarryV2Frontmatter) { c.AssessedBy.Session = "" }, ErrV2MissingField, "assessed_by.session"},
		{"missing assessed_at", func(c *decisionCarryV2Frontmatter) { c.AssessedAt = "" }, ErrV2MissingField, "assessed_at"},
		{"bad assessed_at", func(c *decisionCarryV2Frontmatter) { c.AssessedAt = "yesterday" }, ErrV2EvidenceMalformed, "assessed_at"},
		{"blank reason", func(c *decisionCarryV2Frontmatter) { c.Reason = "  " }, ErrV2MissingField, "carried_forward[0].reason"},
		{"not applying without change", func(c *decisionCarryV2Frontmatter) { c.Applies = &notApplies }, ErrV2MissingField, "material_change"},
		{"owner entry that does not apply", func(c *decisionCarryV2Frontmatter) {
			c.AssessedBy.Role = "owner"
			c.Applies = &notApplies
			c.MaterialChange = "changed"
		}, ErrV2EvidenceMalformed, "must apply"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := validCarryFrontmatter("C-003")
			tt.mutate(&entry)
			_, err := decodeEvidenceV2("test.md", "task", "T-001", ownerValidationWithCarry(entry))
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("decodeEvidenceV2() error = %v, want %v", err, tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantMsg) || !strings.Contains(err.Error(), "test.md") {
				t.Errorf("error = %q, want it to name the file and contain %q", err, tt.wantMsg)
			}
		})
	}
}

func TestDecodeEvidenceV2_carriedForwardOrderAndDuplicates(t *testing.T) {
	ownerEntry := validCarryFrontmatter("C-003")
	ownerEntry.AssessedBy.Role = "owner"

	tests := []struct {
		name    string
		entries []decisionCarryV2Frontmatter
		wantErr bool
	}{
		{"in order", []decisionCarryV2Frontmatter{validCarryFrontmatter("C-003"), validCarryFrontmatter("C-010")}, false},
		{"checker then owner for one check", []decisionCarryV2Frontmatter{validCarryFrontmatter("C-003"), ownerEntry}, false},
		{"out of order", []decisionCarryV2Frontmatter{validCarryFrontmatter("C-010"), validCarryFrontmatter("C-003")}, true},
		{"two checker entries for one check", []decisionCarryV2Frontmatter{validCarryFrontmatter("C-003"), validCarryFrontmatter("C-003")}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := decodeEvidenceV2("test.md", "task", "T-001", ownerValidationWithCarry(tt.entries...))
			if tt.wantErr != (err != nil) {
				t.Fatalf("decodeEvidenceV2() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && !errors.Is(err, ErrV2EvidenceMalformed) {
				t.Errorf("error = %v, want ErrV2EvidenceMalformed", err)
			}
		})
	}
}

func TestDecodeEvidenceV2_scopeAndCarriedForwardRequireAcceptedCheck(t *testing.T) {
	tests := []struct {
		name string
		ov   ownerValidationV2Frontmatter
		want string
	}{
		{"scope", ownerValidationV2Frontmatter{Required: true, Scope: []string{"AC-1"}}, "owner_validation.scope"},
		{"carried_forward", ownerValidationV2Frontmatter{Required: true, CarriedForward: []decisionCarryV2Frontmatter{validCarryFrontmatter("C-003")}}, "owner_validation.carried_forward"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := decodeEvidenceV2("test.md", "task", "T-001", evidenceV2Frontmatter{OwnerValidation: &tt.ov})
			if !errors.Is(err, ErrV2EvidenceMalformed) || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("decodeEvidenceV2() error = %v, want ErrV2EvidenceMalformed naming %s", err, tt.want)
			}
		})
	}
}

func TestDecodeEvidenceV2_ownerValidationScopeBlankEntryRejected(t *testing.T) {
	evidence := ownerValidationWithCarry()
	evidence.OwnerValidation.Scope = []string{"AC-1", " "}
	_, err := decodeEvidenceV2("test.md", "task", "T-001", evidence)
	if !errors.Is(err, ErrV2EvidenceMalformed) || !strings.Contains(err.Error(), "owner_validation.scope") {
		t.Fatalf("decodeEvidenceV2() error = %v, want ErrV2EvidenceMalformed naming owner_validation.scope", err)
	}
}

func checkWithUnmet(result, unmet string) string {
	return "---\nid: C-007\nscope: {kind: objective, id: O-001}\nresult: " + result +
		"\nchecked_by: {role: checker, session: review-1}\nexecuted_session: build-1\nchecked_at: '2026-10-10T00:00:00Z'\n" +
		unmet + "---\n\n# Check\n"
}

func TestDecodeCheckV2_unmet(t *testing.T) {
	check, err := DecodeCheckV2("checks/C-007.md", checkWithUnmet("NEEDS WORK", "unmet: [TEST-08, AC-3]\n"))
	if err != nil {
		t.Fatalf("DecodeCheckV2() error = %v", err)
	}
	if len(check.Unmet) != 2 || check.Unmet[0] != "TEST-08" || check.Unmet[1] != "AC-3" {
		t.Errorf("Unmet = %v, want [TEST-08 AC-3]", check.Unmet)
	}

	check, err = DecodeCheckV2("checks/C-007.md", checkWithUnmet("NEEDS WORK", ""))
	if err != nil || check.Unmet != nil {
		t.Errorf("DecodeCheckV2() = %+v, %v; want nil Unmet and no error without the list", check, err)
	}
}

func TestDecodeCheckV2_unmetRejections(t *testing.T) {
	tests := []struct {
		name, result, unmet, want string
	}{
		{"on CLEAR", "CLEAR", "unmet: [TEST-08]\n", "only allowed on a NEEDS WORK"},
		{"blank entry", "NEEDS WORK", "unmet: [TEST-08, ' ']\n", "empty entry"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecodeCheckV2("checks/C-007.md", checkWithUnmet(tt.result, tt.unmet))
			if !errors.Is(err, ErrV2CheckMalformed) || !strings.Contains(err.Error(), tt.want) || !strings.Contains(err.Error(), "C-007.md") {
				t.Fatalf("DecodeCheckV2() error = %v, want ErrV2CheckMalformed naming the file and %q", err, tt.want)
			}
		})
	}
}

// writeCarryProject writes an Objective with Tasks T-001 and T-002 and Checks
// C-001 and C-002 on T-001, C-003 on T-002, then T-001 carrying the evidence
// under test.
func writeCarryProject(t *testing.T, evidence string) string {
	t.Helper()
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First objective")
	writeV2TaskFixture(t, root, "O-001-first", "T-002-beta.md", "T-002", "Beta", "O-001")
	writeV2CheckFixture(t, root, "C-001-a.md", "C-001", "task", "T-001")
	testutil.WriteFile(t, filepath.Join(root, v2ChecksDirName, "C-002-a.md"),
		"---\nid: C-002\nscope: {kind: task, id: T-001}\nresult: CLEAR\nchecked_by: {role: checker, session: sess-1}\nexecuted_session: build-fixture\nchecked_at: '2026-09-15T00:00:00Z'\nsupersedes: C-001\n---\n\n# Check\n")
	writeV2CheckFixture(t, root, "C-003-b.md", "C-003", "task", "T-002")
	testutil.WriteFile(t, filepath.Join(root, v2ObjectivesDirName, "O-001-first", v2TasksDirName, "T-001-alpha.md"),
		"---\nid: T-001\ntitle: \"Alpha\"\nobjective: O-001\nplanned_by: {role: planner, session: planning-fixture}\nstatus: planned\n"+evidence+"---\n\n# Alpha\n")
	return root
}

func carryEvidence(field, check string) string {
	entry := "{check: " + check + ", applies: true, assessed_by: {role: checker, session: check-1}, assessed_at: '2026-10-10T00:00:00Z', reason: unchanged}"
	if field == "exception" {
		return "exception: {requirements: [TEST-08], reason: waived, owner: ani, recorded_at: '2026-10-09T00:00:00Z', check: C-001, carried_forward: [" + entry + "]}\n"
	}
	return "owner_validation: {required: true, accepted_check: C-001, accepted_by: {role: owner, session: owner-1}, carried_forward: [" + entry + "]}\n"
}

func TestLoadV2Index_carriedForwardReferences(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		check    string
		wantErr  error
		wantText string
	}{
		{"owner_validation same scope", "owner_validation", "C-002", nil, ""},
		{"exception same scope", "exception", "C-002", nil, ""},
		{"owner_validation dangling", "owner_validation", "C-999", ErrV2EvidenceMissingReference, "owner_validation.carried_forward[0].check"},
		{"exception dangling", "exception", "C-999", ErrV2EvidenceMissingReference, "exception.carried_forward[0].check"},
		{"owner_validation other scope", "owner_validation", "C-003", ErrV2EvidenceMalformed, "task T-002"},
		{"exception other scope", "exception", "C-003", ErrV2EvidenceMalformed, "task T-002"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := writeCarryProject(t, carryEvidence(tt.field, tt.check))
			index, err := LoadV2Index(root)
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("LoadV2Index() error = %v", err)
				}
				if index.Tasks["T-001"].Evidence == nil {
					t.Fatal("T-001 evidence = nil, want loaded")
				}
				return
			}
			if !errors.Is(err, tt.wantErr) || !strings.Contains(err.Error(), tt.wantText) || !strings.Contains(err.Error(), "T-001-alpha.md") {
				t.Fatalf("LoadV2Index() error = %v, want %v naming the file and %q", err, tt.wantErr, tt.wantText)
			}
		})
	}
}

func fullCarryEvidenceV2() *Evidence {
	evidence := fullTaskEvidenceV2()
	at := time.Date(2026, 10, 10, 0, 0, 0, 0, time.UTC)
	evidence.OwnerValidation.Scope = []string{"AC-1"}
	evidence.OwnerValidation.CarriedForward = []DecisionCarry{
		{Check: "C-002", Applies: true, AssessedBy: Actor{Role: ActorRoleChecker, Session: "check-1"}, AssessedAt: at, Reason: "unchanged"},
		{Check: "C-003", Applies: false, AssessedBy: Actor{Role: ActorRoleChecker, Session: "check-2"}, AssessedAt: at, Reason: "AC-1 changed", MaterialChange: "AC-1 now rejects empty input"},
	}
	evidence.Exception.CarriedForward = []DecisionCarry{
		{Check: "C-002", Applies: true, AssessedBy: Actor{Role: ActorRoleOwner, Session: "owner-2"}, AssessedAt: at, Reason: "renewed"},
	}
	return evidence
}

func assertCarryEvidenceRoundTrips(t *testing.T, got, want *Evidence) {
	t.Helper()
	if got == nil || got.OwnerValidation == nil || got.Exception == nil {
		t.Fatalf("Evidence = %+v, want owner_validation and exception", got)
	}
	if strings.Join(got.OwnerValidation.Scope, ",") != strings.Join(want.OwnerValidation.Scope, ",") {
		t.Errorf("Scope = %v, want %v", got.OwnerValidation.Scope, want.OwnerValidation.Scope)
	}
	assertCarriesEqual(t, "owner_validation", got.OwnerValidation.CarriedForward, want.OwnerValidation.CarriedForward)
	assertCarriesEqual(t, "exception", got.Exception.CarriedForward, want.Exception.CarriedForward)
}

func assertCarriesEqual(t *testing.T, name string, got, want []DecisionCarry) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s carried_forward = %+v, want %+v", name, got, want)
	}
	for i := range want {
		if got[i].Check != want[i].Check || got[i].Applies != want[i].Applies || got[i].AssessedBy != want[i].AssessedBy ||
			!got[i].AssessedAt.Equal(want[i].AssessedAt) || got[i].Reason != want[i].Reason || got[i].MaterialChange != want[i].MaterialChange {
			t.Errorf("%s carried_forward[%d] = %+v, want %+v", name, i, got[i], want[i])
		}
	}
}

func TestWriteEvidenceV2_roundTripsScopeAndCarriedForwardForTaskAndObjective(t *testing.T) {
	want := fullCarryEvidenceV2()
	dir := t.TempDir()

	taskPath := filepath.Join(dir, "T-020.md")
	taskContent := "---\nid: T-020\ntitle: \"Plain\"\nobjective: O-002\nplanned_by: {role: planner, session: planning-fixture}\nstatus: in_progress\nstage: build\nunknown_key: keep me\n---\n\n# Task\n\nAuthored notes.\n"
	writeTestFile(t, taskPath, taskContent)
	task := mustDecodeTaskV2(t, taskPath, taskContent)
	task.Evidence = want
	if err := WriteTaskEvidenceV2(task); err != nil {
		t.Fatalf("WriteTaskEvidenceV2() error = %v", err)
	}
	written := string(readTestBytes(t, taskPath))
	reparsed, err := DecodeTaskV2(taskPath, written)
	if err != nil {
		t.Fatalf("DecodeTaskV2() after write error = %v", err)
	}
	assertCarryEvidenceRoundTrips(t, reparsed.Evidence, want)
	if !strings.Contains(written, "unknown_key: keep me") || !strings.Contains(written, "Authored notes.") {
		t.Errorf("written task lost unknown frontmatter or body:\n%s", written)
	}

	objectivePath := filepath.Join(dir, "Objective.md")
	objectiveContent := "---\nid: O-002\ntitle: \"Plain\"\nstatus: in_progress\nunknown_key: keep me\n---\n\n# Objective\n\nAuthored notes.\n"
	writeTestFile(t, objectivePath, objectiveContent)
	objective, err := DecodeObjectiveV2(objectivePath, objectiveContent)
	if err != nil {
		t.Fatalf("DecodeObjectiveV2() error = %v", err)
	}
	objective.Evidence = want
	objective.Evidence.CheckWaiver = nil
	if err := WriteObjectiveEvidenceV2(objective); err != nil {
		t.Fatalf("WriteObjectiveEvidenceV2() error = %v", err)
	}
	written = string(readTestBytes(t, objectivePath))
	reloaded, err := DecodeObjectiveV2(objectivePath, written)
	if err != nil {
		t.Fatalf("DecodeObjectiveV2() after write error = %v", err)
	}
	assertCarryEvidenceRoundTrips(t, reloaded.Evidence, want)
	if !strings.Contains(written, "unknown_key: keep me") || !strings.Contains(written, "Authored notes.") {
		t.Errorf("written objective lost unknown frontmatter or body:\n%s", written)
	}
}

func TestWriteTaskEvidenceV2_omitsScopeAndCarriedForwardWhenEmpty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "T-020.md")
	content := "---\nid: T-020\ntitle: \"Plain\"\nobjective: O-002\nplanned_by: {role: planner, session: planning-fixture}\nstatus: in_progress\nstage: build\n---\n\n# Task\n"
	writeTestFile(t, path, content)
	task := mustDecodeTaskV2(t, path, content)
	task.Evidence = fullTaskEvidenceV2()
	if err := WriteTaskEvidenceV2(task); err != nil {
		t.Fatalf("WriteTaskEvidenceV2() error = %v", err)
	}
	written := string(readTestBytes(t, path))
	for _, key := range []string{"carried_forward", "scope:", "material_change"} {
		if strings.Contains(written, key) {
			t.Errorf("written task contains %q; a record without the new fields must write as before:\n%s", key, written)
		}
	}
}
