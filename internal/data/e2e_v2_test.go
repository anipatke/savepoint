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
	writeV2ObjectiveFixture(t, root, "O001-ship", "O001", "Ship it")

	// T001 is a technical Task: no owner_validation declared.
	t001Path := filepath.Join(root, v2ObjectivesDirName, "O001-ship", v2TasksDirName, "T001-alpha.md")
	testutil.WriteFile(t, t001Path, "---\nid: T001\ntitle: \"Alpha\"\nobjective: O001\nstatus: in_progress\nstage: audit\n---\n\n# Alpha\n\nAuthored plan notes for Alpha.\n")

	// T002 declares owner_validation.required and carries an unknown
	// frontmatter field and authored body content that every later evidence
	// patch must preserve untouched.
	t002Path := filepath.Join(root, v2ObjectivesDirName, "O001-ship", v2TasksDirName, "T002-beta.md")
	testutil.WriteFile(t, t002Path, "---\nid: T002\ntitle: \"Beta\"\nobjective: O001\nstatus: in_progress\nstage: audit\ncustom_note: keep-me\n---\n\n# Beta\n\nAuthored plan notes for Beta.\n")

	checkedAt := time.Date(2026, 9, 14, 0, 0, 0, 0, time.UTC)

	// --- records load: no Check yet, both Tasks resolve to missing clearance ---
	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	if got := ResolveClearance(index, "T001").State; got != ClearanceMissing {
		t.Fatalf("T001 clearance = %q, want missing before any Check is recorded", got)
	}
	if got := ResolveTaskCompletion(index, "T001"); got.Allowed {
		t.Fatalf("ResolveTaskCompletion(T001) = %+v, want blocked before any Check", got)
	}

	// --- a technical Task closes under checker authority once current ---
	c001, err := CreateCheckV2(root, index, NewCheckV2{
		Scope:     CheckScope{Kind: CheckScopeTask, ID: "T001"},
		Result:    CheckResultClear,
		CheckedBy: Actor{Role: ActorRoleChecker, Session: "sess-1"},
		CheckedAt: checkedAt,
		Body:      "\n\n# Check\n\nT001 reviewed clean.\n",
	})
	if err != nil {
		t.Fatalf("CreateCheckV2(C001) error = %v", err)
	}

	index, err = LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	t001 := index.Tasks["T001"]
	t001.Evidence = &Evidence{Freshness: &Freshness{
		State: FreshnessCurrent, Check: c001.ID,
		AssessedBy: Actor{Role: ActorRoleChecker, Session: "sess-1"},
		AssessedAt: checkedAt, Basis: "reviewed the diff",
	}}
	if err := WriteTaskEvidenceV2(t001); err != nil {
		t.Fatalf("WriteTaskEvidenceV2(T001) error = %v", err)
	}

	index, err = LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	if got := ResolveClearance(index, "T001").State; got != ClearanceCurrent {
		t.Fatalf("T001 clearance = %q, want current after freshness names the latest Check", got)
	}
	completion := ResolveTaskCompletion(index, "T001")
	if !completion.Allowed || completion.Actor != ActorRoleChecker {
		t.Fatalf("ResolveTaskCompletion(T001) = %+v, want allowed under checker authority", completion)
	}

	// --- an owner-validated Task waits until the owner accepts ---
	c002, err := CreateCheckV2(root, index, NewCheckV2{
		Scope:     CheckScope{Kind: CheckScopeTask, ID: "T002"},
		Result:    CheckResultClear,
		CheckedBy: Actor{Role: ActorRoleChecker, Session: "sess-1"},
		CheckedAt: checkedAt,
		Body:      "\n\n# Check\n\nT002 reviewed clean.\n",
	})
	if err != nil {
		t.Fatalf("CreateCheckV2(C002) error = %v", err)
	}

	index, err = LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	t002 := index.Tasks["T002"]
	t002.Evidence = &Evidence{
		Freshness: &Freshness{
			State: FreshnessCurrent, Check: c002.ID,
			AssessedBy: Actor{Role: ActorRoleChecker, Session: "sess-1"},
			AssessedAt: checkedAt, Basis: "reviewed the diff",
		},
		OwnerValidation: &OwnerValidation{Required: true},
	}
	if err := WriteTaskEvidenceV2(t002); err != nil {
		t.Fatalf("WriteTaskEvidenceV2(T002) error = %v", err)
	}

	index, err = LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	waiting := ResolveTaskCompletion(index, "T002")
	if waiting.Allowed {
		t.Fatalf("ResolveTaskCompletion(T002) = %+v, want blocked pending owner acceptance", waiting)
	}
	if len(waiting.Blockers) != 1 || waiting.Blockers[0].Kind != GateBlockOwnerAcceptance {
		t.Fatalf("ResolveTaskCompletion(T002) Blockers = %+v, want one GateBlockOwnerAcceptance", waiting.Blockers)
	}

	// The owner accepts C002.
	t002 = index.Tasks["T002"]
	t002.Evidence.OwnerValidation.AcceptedCheck = c002.ID
	t002.Evidence.OwnerValidation.AcceptedBy = Actor{Role: ActorRoleOwner, Session: "owner-1"}
	if err := WriteTaskEvidenceV2(t002); err != nil {
		t.Fatalf("WriteTaskEvidenceV2(T002) accept error = %v", err)
	}

	index, err = LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	accepted := ResolveTaskCompletion(index, "T002")
	if !accepted.Allowed || accepted.AllowedByException {
		t.Fatalf("ResolveTaskCompletion(T002) = %+v, want allowed by owner acceptance, not exception", accepted)
	}

	// --- a rerun supersedes the prior acceptance and freshness ---
	_, err = CreateCheckV2(root, index, NewCheckV2{
		Scope:      CheckScope{Kind: CheckScopeTask, ID: "T002"},
		Result:     CheckResultClear,
		CheckedBy:  Actor{Role: ActorRoleChecker, Session: "sess-2"},
		CheckedAt:  checkedAt.Add(24 * time.Hour),
		Supersedes: c002.ID,
		Body:       "\n\n# Check\n\nT002 re-reviewed after a rerun.\n",
	})
	if err != nil {
		t.Fatalf("CreateCheckV2(rerun) error = %v", err)
	}

	index, err = LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	rerunClearance := ResolveClearance(index, "T002")
	if rerunClearance.State != ClearanceStale {
		t.Fatalf("T002 clearance after rerun = %q, want stale (freshness still names the superseded Check)", rerunClearance.State)
	}
	blockedAgain := ResolveTaskCompletion(index, "T002")
	if blockedAgain.Allowed {
		t.Fatalf("ResolveTaskCompletion(T002) = %+v, want blocked again after the rerun supersedes the accepted Check", blockedAgain)
	}

	// --- evidence writes preserve authored content throughout ---
	rawT002, err := os.ReadFile(t002Path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(rawT002), "custom_note: keep-me") {
		t.Errorf("T002 file lost its unknown custom_note field:\n%s", rawT002)
	}
	if !strings.Contains(string(rawT002), "Authored plan notes for Beta.") {
		t.Errorf("T002 file lost its authored Markdown body:\n%s", rawT002)
	}
}
