package data

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/opencode/savepoint/internal/testutil"
)

// TestE43_EpicScenario proves the E43 gate-check behavior end to end on a
// temporary project built from real files, not just in-memory index
// fixtures: records load, clearance resolves, a technical Task closes under
// checker authority, an owner-validated Task waits on the owner, a rerun
// supersedes the prior acceptance and freshness, and evidence writes
// preserve authored content and unknown fields.
func TestE43_EpicScenario(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O-001-ship", "O-001", "Ship it")

	// T-001 is a technical Task: no owner_validation declared.
	t001Path := filepath.Join(root, v2ObjectivesDirName, "O-001-ship", v2TasksDirName, "T-001-alpha.md")
	testutil.WriteFile(t, t001Path, "---\nid: T-001\ntitle: \"Alpha\"\nobjective: O-001\nplanned_by: {role: planner, session: planning-fixture}\nstatus: in_progress\nstage: audit\n---\n\n# Alpha\n\nAuthored plan notes for Alpha.\n")

	// T-002 declares owner_validation.required and carries an unknown
	// frontmatter field and authored body content that every later evidence
	// patch must preserve untouched.
	t002Path := filepath.Join(root, v2ObjectivesDirName, "O-001-ship", v2TasksDirName, "T-002-beta.md")
	testutil.WriteFile(t, t002Path, "---\nid: T-002\ntitle: \"Beta\"\nobjective: O-001\nplanned_by: {role: planner, session: planning-fixture}\nstatus: in_progress\nstage: audit\ncustom_note: keep-me\n---\n\n# Beta\n\nAuthored plan notes for Beta.\n")

	checkedAt := time.Date(2026, 9, 14, 0, 0, 0, 0, time.UTC)

	// --- records load: no Check yet, both Tasks resolve to missing clearance ---
	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	if got := ResolveClearance(index, "T-001").State; got != ClearanceMissing {
		t.Fatalf("T-001 clearance = %q, want missing before any Check is recorded", got)
	}
	if got := ResolveTaskCompletion(index, "T-001"); got.Allowed {
		t.Fatalf("ResolveTaskCompletion(T-001) = %+v, want blocked before any Check", got)
	}

	// --- a technical Task closes under checker authority once current ---
	c001 := writeScenarioCheck(t, root, "C-001", CheckScopeTask, "T-001", "sess-1", checkedAt, "")

	index, err = LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	t001 := index.Tasks["T-001"]
	t001.Evidence = &Evidence{Freshness: &Freshness{
		State: FreshnessCurrent, Check: c001.ID,
		AssessedBy: Actor{Role: ActorRoleChecker, Session: "sess-1"},
		AssessedAt: checkedAt, Basis: "reviewed the diff",
	}}
	if err := WriteTaskEvidenceV2(t001); err != nil {
		t.Fatalf("WriteTaskEvidenceV2(T-001) error = %v", err)
	}

	index, err = LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	if got := ResolveClearance(index, "T-001").State; got != ClearanceCurrent {
		t.Fatalf("T-001 clearance = %q, want current after freshness names the latest Check", got)
	}
	completion := ResolveTaskCompletion(index, "T-001")
	if !completion.Allowed || completion.Actor != ActorRoleChecker {
		t.Fatalf("ResolveTaskCompletion(T-001) = %+v, want allowed under checker authority", completion)
	}

	// --- an owner-validated Task waits until the owner accepts ---
	c002 := writeScenarioCheck(t, root, "C-002", CheckScopeTask, "T-002", "sess-1", checkedAt, "")

	index, err = LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	t002 := index.Tasks["T-002"]
	t002.Evidence = &Evidence{
		Freshness: &Freshness{
			State: FreshnessCurrent, Check: c002.ID,
			AssessedBy: Actor{Role: ActorRoleChecker, Session: "sess-1"},
			AssessedAt: checkedAt, Basis: "reviewed the diff",
		},
		OwnerValidation: &OwnerValidation{Required: true},
	}
	if err := WriteTaskEvidenceV2(t002); err != nil {
		t.Fatalf("WriteTaskEvidenceV2(T-002) error = %v", err)
	}

	index, err = LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	waiting := ResolveTaskCompletion(index, "T-002")
	if waiting.Allowed {
		t.Fatalf("ResolveTaskCompletion(T-002) = %+v, want blocked pending owner acceptance", waiting)
	}
	if len(waiting.Blockers) != 1 || waiting.Blockers[0].Kind != GateBlockOwnerAcceptance {
		t.Fatalf("ResolveTaskCompletion(T-002) Blockers = %+v, want one GateBlockOwnerAcceptance", waiting.Blockers)
	}

	// The owner accepts C-002.
	t002 = index.Tasks["T-002"]
	t002.Evidence.OwnerValidation.AcceptedCheck = c002.ID
	t002.Evidence.OwnerValidation.AcceptedBy = Actor{Role: ActorRoleOwner, Session: "owner-1"}
	if err := WriteTaskEvidenceV2(t002); err != nil {
		t.Fatalf("WriteTaskEvidenceV2(T-002) accept error = %v", err)
	}

	index, err = LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	accepted := ResolveTaskCompletion(index, "T-002")
	if !accepted.Allowed || accepted.AllowedByException {
		t.Fatalf("ResolveTaskCompletion(T-002) = %+v, want allowed by owner acceptance, not exception", accepted)
	}

	// --- a rerun supersedes the prior acceptance; its CLEAR is current on
	//     its own, but the owner has not accepted it yet ---
	writeScenarioCheck(t, root, "C-003", CheckScopeTask, "T-002", "sess-2", checkedAt.Add(24*time.Hour), c002.ID)

	index, err = LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	rerunClearance := ResolveClearance(index, "T-002")
	if rerunClearance.State != ClearanceCurrent || rerunClearance.Check != "C-003" {
		t.Fatalf("T-002 clearance after rerun = %+v, want current on C-003 (a CLEAR re-check needs no freshness record)", rerunClearance)
	}
	blockedAgain := ResolveTaskCompletion(index, "T-002")
	if blockedAgain.Allowed || len(blockedAgain.Blockers) != 1 || blockedAgain.Blockers[0].Kind != GateBlockOwnerAcceptance {
		t.Fatalf("ResolveTaskCompletion(T-002) = %+v, want blocked on owner acceptance after the rerun supersedes the accepted Check", blockedAgain)
	}

	// --- evidence writes preserve authored content throughout ---
	rawT002, err := os.ReadFile(t002Path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(rawT002), "custom_note: keep-me") {
		t.Errorf("T-002 file lost its unknown custom_note field:\n%s", rawT002)
	}
	if !strings.Contains(string(rawT002), "Authored plan notes for Beta.") {
		t.Errorf("T-002 file lost its authored Markdown body:\n%s", rawT002)
	}
}

// TestE44_EpicScenario proves the E44 Issue and Objective-integration
// behavior end to end on a temporary project built from real files: a Task
// cannot start while its owning Objective waits on a dependency Objective
// with no current integration clearance, the dependency Objective closes
// under checker authority once its own Objective-scoped Check is current,
// readiness then clears for the downstream Task, an Issue observed against
// that Objective Check resolves as verified proof, and a later Check against
// the same scope leaves both the Objective's clearance and the Issue's proof
// stale — detected by InspectObjectiveConsistency and InspectIssueConsistency
// without rewriting any record.
func TestE44_EpicScenario(t *testing.T) {
	root := t.TempDir()
	checkedAt := time.Date(2026, 9, 15, 0, 0, 0, 0, time.UTC)

	writeV2ObjectiveFixture(t, root, "O-001-ship", "O-001", "Ship it")
	writeV2ObjectiveFixture(t, root, "O-002-follow-on", "O-002", "Follow on", "O-001")
	writeV2TaskFixture(t, root, "O-001-ship", "T-001-alpha.md", "T-001", "Alpha", "O-001")
	writeV2TaskFixture(t, root, "O-002-follow-on", "T-002-beta.md", "T-002", "Beta", "O-002")

	// --- T-002 cannot start while its owning Objective O-002 waits on O-001 ---
	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	start := ResolveTaskStart(index, "T-002")
	if start.Allowed {
		t.Fatalf("ResolveTaskStart(T-002) = %+v, want blocked on the owning objective's dependency", start)
	}
	if len(start.Blockers) != 1 || start.Blockers[0].Kind != GateBlockObjectiveDependency {
		t.Fatalf("ResolveTaskStart(T-002) Blockers = %+v, want one GateBlockObjectiveDependency", start.Blockers)
	}

	// --- close T-001 under checker authority ---
	c001 := writeScenarioCheck(t, root, "C-001", CheckScopeTask, "T-001", "sess-1", checkedAt, "")
	index, err = LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	t001 := index.Tasks["T-001"]
	t001.Evidence = &Evidence{Freshness: &Freshness{
		State: FreshnessCurrent, Check: c001.ID,
		AssessedBy: Actor{Role: ActorRoleChecker, Session: "sess-1"}, AssessedAt: checkedAt, Basis: "reviewed",
	}}
	if err := WriteTaskEvidenceV2(t001); err != nil {
		t.Fatalf("WriteTaskEvidenceV2(T-001) error = %v", err)
	}
	t001.Status = ColumnDone
	if err := WriteTaskV2(t001); err != nil {
		t.Fatalf("WriteTaskV2(T-001) error = %v", err)
	}

	// --- close O-001 itself under checker authority once its integration
	//     Check is current ---
	index, err = LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	co1 := writeScenarioCheck(t, root, "C-002", CheckScopeObjective, "O-001", "sess-1", checkedAt, "")
	index, err = LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	o001 := index.Objectives["O-001"]
	o001.Evidence = &Evidence{Freshness: &Freshness{
		State: FreshnessCurrent, Check: co1.ID,
		AssessedBy: Actor{Role: ActorRoleChecker, Session: "sess-1"}, AssessedAt: checkedAt, Basis: "integration reviewed",
	}}
	if err := WriteObjectiveEvidenceV2(o001); err != nil {
		t.Fatalf("WriteObjectiveEvidenceV2(O-001) error = %v", err)
	}

	index, err = LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	completion := ResolveObjectiveCompletion(index, "O-001")
	if !completion.Allowed || completion.Actor != ActorRoleChecker {
		t.Fatalf("ResolveObjectiveCompletion(O-001) = %+v, want allowed under checker authority", completion)
	}
	o001 = index.Objectives["O-001"]
	o001.Status = ColumnDone
	if err := WriteObjectiveV2(o001); err != nil {
		t.Fatalf("WriteObjectiveV2(O-001) error = %v", err)
	}

	// --- O-002's dependency on O-001 is now satisfied, so T-002 can start ---
	index, err = LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	dep := ResolveObjectiveDependency(index, "O-001")
	if !dep.Satisfied {
		t.Fatalf("ResolveObjectiveDependency(O-001) = %+v, want satisfied", dep)
	}
	start2 := ResolveTaskStart(index, "T-002")
	if !start2.Allowed {
		t.Fatalf("ResolveTaskStart(T-002) = %+v, want allowed once the owning objective's dependency clears", start2)
	}

	// --- an Issue observed against O-001's integration Check resolves as
	//     verified proof, using that same Check ---
	issueID := "I-001"
	testutil.WriteFile(t, filepath.Join(root, v2IssuesDirName, issueID+"-latency-regression.md"),
		"---\nid: "+issueID+"\ntitle: \"Latency regression\"\ntype: defect\nstatus: resolved\n"+
			"source: {kind: check, check: "+co1.ID+", actor: {role: checker, session: sess-1}, at: '"+checkedAt.Format(time.RFC3339)+"'}\n"+
			"checks: ["+co1.ID+"]\n"+
			"resolution: {disposition: verified, check: "+co1.ID+", actor: {role: checker, session: sess-1}, at: '"+checkedAt.Format(time.RFC3339)+"'}\n"+
			"---\n\n# Issue\n\nObserved during integration review.\n")

	index, err = LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	if got := index.IssueStatusCounts()[IssueStatusResolved]; got != 1 {
		t.Errorf("IssueStatusCounts()[resolved] = %d, want 1", got)
	}
	if got := index.IssueTypeCounts()[IssueTypeDefect]; got != 1 {
		t.Errorf("IssueTypeCounts()[defect] = %d, want 1", got)
	}
	if got := InspectObjectiveConsistency(index); len(got) != 0 {
		t.Fatalf("InspectObjectiveConsistency() = %+v, want none before a later Check supersedes O-001's proof", got)
	}
	if got := InspectIssueConsistency(index); len(got) != 0 {
		t.Fatalf("InspectIssueConsistency() = %+v, want none while the verified proof is still latest", got)
	}

	// --- a later CLEAR Objective-scoped Check keeps O-001 cleared on its
	//     own but leaves the Issue's verified proof superseded, without
	//     rewriting any record ---
	writeScenarioCheck(t, root, "C-003", CheckScopeObjective, "O-001", "sess-2", checkedAt.Add(24*time.Hour), co1.ID)

	index, err = LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	if objectiveProblems := InspectObjectiveConsistency(index); len(objectiveProblems) != 0 {
		t.Fatalf("InspectObjectiveConsistency() = %+v, want none (the later CLEAR Check is current on its own)", objectiveProblems)
	}
	issueProblems := InspectIssueConsistency(index)
	if len(issueProblems) != 1 || issueProblems[0].Kind != IssueConsistencyProofSuperseded || issueProblems[0].Issue != issueID {
		t.Fatalf("InspectIssueConsistency() = %+v, want one IssueConsistencyProofSuperseded naming %s", issueProblems, issueID)
	}
}

// writeScenarioCheck writes one CLEAR Check record as a checker session would
// author it, loads it back through the strict decoder, and returns it so the
// scenario can name its ID in later evidence.
func writeScenarioCheck(t *testing.T, root, id string, kind CheckScopeKind, scopeID, session string, at time.Time, supersedes string) *CheckV2 {
	t.Helper()
	content := "---\nid: " + id + "\nscope: {kind: " + string(kind) + ", id: " + scopeID + "}\nresult: CLEAR\n" +
		"checked_by: {role: checker, session: " + session + "}\nexecuted_session: build-001\nchecked_at: '" + at.Format(time.RFC3339) + "'\n"
	if supersedes != "" {
		content += "supersedes: " + supersedes + "\n"
	}
	content += "---\n\n# Check\n\n" + scopeID + " reviewed.\n"
	relPath := filepath.Join(v2ChecksDirName, id+".md")
	check, err := DecodeCheckV2(relPath, content)
	if err != nil {
		t.Fatalf("DecodeCheckV2(%s) error = %v", id, err)
	}
	testutil.WriteFile(t, filepath.Join(root, relPath), content)
	return check
}
