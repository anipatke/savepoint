package v2

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/opencode/savepoint/internal/migrate"
	"github.com/opencode/savepoint/internal/testutil"
)

// Every fixture below builds a whole project inside a temporary directory and
// returns its .savepoint root. No test in this package reads or writes the
// live Savepoint project (TEST-04).

// writeValidProject writes a V2 project whose router selects a planned Task
// with no dependencies: one Objective, two Tasks, one of them done, so counts
// and columns have something to prove.
func writeValidProject(t *testing.T) string {
	t.Helper()
	root := savepointRoot(t)
	writeConfig(t, root)
	writeRouter(t, root, "task", "O001", "T001")
	writeObjective(t, root, "O001", "First objective", "planned")
	writeTask(t, root, "O001", "T001", "Do the thing", "status: planned\n")
	writeTask(t, root, "O001", "T002", "Already finished", "status: done\n")
	return root
}

// writeEmptyProjectFromTemplate writes the project a user gets from
// `savepoint init`: the shipped templates/project-v2 .savepoint files, copied
// byte-for-byte, with no Objective, Task, Check, or Issue.
func writeEmptyProjectFromTemplate(t *testing.T) string {
	t.Helper()
	root := savepointRoot(t)
	template := filepath.Join("..", "..", "..", "templates", "project-v2", ".savepoint")
	for _, name := range []string{"config.yml", "router.md"} {
		content, err := os.ReadFile(filepath.Join(template, name))
		if err != nil {
			t.Fatalf("read shipped template %s: %v", name, err)
		}
		testutil.WriteFile(t, filepath.Join(root, name), string(content))
	}
	return root
}

func savepointRoot(t *testing.T) string {
	t.Helper()
	root := filepath.Join(t.TempDir(), ".savepoint")
	testutil.MkdirAll(t, root)
	return root
}

func writeConfig(t *testing.T, root string) {
	t.Helper()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
}

func writeRouter(t *testing.T, root, state, objective, task string) {
	t.Helper()
	testutil.WriteFile(t, filepath.Join(root, "router.md"), routerContent(state, objective, task))
}

func routerContent(state, objective, task string) string {
	return "# Router\n\n## Current state\n\n```yaml\nstate: " + state +
		"\nobjective: " + objective + "\ntask: " + task + "\nnext_action: \"do the next thing\"\n```\n"
}

func writeObjective(t *testing.T, root, id, title, status string) {
	t.Helper()
	writeObjectiveExtra(t, root, id, title, status, "")
}

// writeObjectiveExtra writes an Objective whose frontmatter carries extra
// lines — its own integration evidence, or a dependency on another Objective.
func writeObjectiveExtra(t *testing.T, root, id, title, status, extra string) {
	t.Helper()
	testutil.WriteFile(t, objectivePath(root, id),
		"---\nid: "+id+"\ntitle: \""+title+"\"\nstatus: "+status+"\n"+extra+"---\n\n# "+title+"\n")
}

// writeTask writes a Task under objective's directory. extra carries the
// frontmatter lines that vary per fixture, so an invalid shape is one line
// different from a valid one.
func writeTask(t *testing.T, root, objective, id, title, extra string) {
	t.Helper()
	writeTaskInDir(t, root, objective, objective, id, title, extra)
}

// writeTaskInDir writes a Task whose file lives under dir's Objective
// directory while its own objective field names owner. The two differ only in
// the fixtures that prove ownership is read from the record rather than from
// where the file happens to sit.
func writeTaskInDir(t *testing.T, root, dir, owner, id, title, extra string) {
	t.Helper()
	testutil.WriteFile(t, taskPath(root, dir, id),
		"---\nid: "+id+"\ntitle: \""+title+"\"\nobjective: "+owner+
			"\nplanned_by: {role: planner, session: board-fixture}\n"+extra+"---\n\n# "+title+"\n")
}

// writeCheck writes one recorded Check over a Task or Objective. result is
// CLEAR or NEEDS WORK; the checker and execution sessions differ, as
// DecodeCheckV2 requires.
func writeCheck(t *testing.T, root, id, scopeKind, target, result string) {
	t.Helper()
	writeCheckExtra(t, root, id, scopeKind, target, result, "")
}

// writeCheckExtra writes a Check whose frontmatter carries extra lines — a
// supersedes reference to the run it replaces, or the Issues it opened.
func writeCheckExtra(t *testing.T, root, id, scopeKind, target, result, extra string) {
	t.Helper()
	testutil.WriteFile(t, filepath.Join(root, "checks", id+".md"),
		"---\nid: "+id+"\nscope: {kind: "+scopeKind+", id: "+target+"}\nresult: "+result+
			"\nchecked_by: {role: checker, session: checker-fixture}\nexecuted_session: executor-fixture"+
			"\nchecked_at: 2026-01-02T00:00:00Z\n"+extra+"---\n\n# "+id+"\n")
}

// currentFreshness is the evidence block that makes a CLEAR Check count as
// current clearance: a checker-assessed freshness naming that same Check.
func currentFreshness(check string) string {
	return "freshness:\n  state: current\n  check: " + check +
		"\n  assessed_by: {role: checker, session: checker-fixture}" +
		"\n  assessed_at: 2026-01-02T01:00:00Z\n  basis: \"reran the suite\"\n"
}

// writeBadgeProject builds one project holding every card state the board can
// show: each stage, each clearance state, a dependency wait, a replan flag, an
// owner wait, an ordinary done, a done by recorded exception, and a done whose
// clearance has gone stale.
func writeBadgeProject(t *testing.T) string {
	t.Helper()
	root := savepointRoot(t)
	writeConfig(t, root)
	writeRouter(t, root, "task", "O001", "T001")
	writeObjective(t, root, "O001", "Ship the board", "in_progress")

	writeTask(t, root, "O001", "T001", "Planned with nothing recorded", "status: planned\n")
	writeTask(t, root, "O001", "T002", "Planned and waiting on T001", "status: planned\ndepends_on:\n  - {task: T001, requires: clear}\n")
	writeTask(t, root, "O001", "T003", "Being built right now", "status: in_progress\nstage: build\n")
	writeTask(t, root, "O001", "T004", "Being tested, and the plan slipped",
		"status: in_progress\nstage: test\nreplan:\n  reason: \"the plan no longer matches\"\n  recorded_by: {role: planner, session: planner-fixture}\n  recorded_at: 2026-01-01T00:00:00Z\n")

	writeCheck(t, root, "C001", "task", "T005", "NEEDS WORK")
	writeTask(t, root, "O001", "T005", "At audit with a check that found problems", "status: in_progress\nstage: audit\nlast_check: C001\n")

	writeCheck(t, root, "C002", "task", "T006", "CLEAR")
	writeTask(t, root, "O001", "T006", "At audit and waiting on the owner",
		"status: in_progress\nstage: audit\nlast_check: C002\n"+currentFreshness("C002")+
			"owner_validation:\n  required: true\n")

	writeCheck(t, root, "C003", "task", "T007", "CLEAR")
	writeTask(t, root, "O001", "T007", "Done and cleared", "status: done\nlast_check: C003\n"+currentFreshness("C003"))

	writeCheck(t, root, "C004", "task", "T008", "NEEDS WORK")
	writeTask(t, root, "O001", "T008", "Done because the owner accepted an exception",
		"status: done\nlast_check: C004\nexception:\n  requirements: [\"TEST-02\"]\n  reason: \"shipped with a known gap\"\n  owner: \"the owner\"\n  recorded_at: 2026-01-03T00:00:00Z\n  check: C004\n")

	writeCheck(t, root, "C005", "task", "T009", "CLEAR")
	writeTask(t, root, "O001", "T009", "Done, but the check was never assessed", "status: done\nlast_check: C005\n")

	return root
}

// staleFreshness is an assessment that names its Check but no longer calls it
// current, which is what ResolveClearance reads as stale.
func staleFreshness(check string) string {
	return "freshness:\n  state: stale\n  check: " + check +
		"\n  assessed_by: {role: checker, session: checker-fixture}" +
		"\n  assessed_at: 2026-01-02T01:00:00Z\n  basis: \"the code moved on\"\n"
}

// writeNavigationProject builds the project the sidebar exists for: six
// Objectives covering every clearance state, an Objective whose Tasks are all
// done without its own integration Check, an Objective waiting on another, and
// one Task whose file sits under a different Objective's directory than the one
// its record names.
//
// The router selects O003, so a board opened over this project starts filtered
// to the Objective the router names.
func writeNavigationProject(t *testing.T) string {
	t.Helper()
	root := savepointRoot(t)
	writeConfig(t, root)
	writeRouter(t, root, "task", "O003", "T003")

	// O001: every owned Task done, and its own integration Check current.
	writeCheck(t, root, "C001", "objective", "O001", "CLEAR")
	writeObjectiveExtra(t, root, "O001", "Finished and integrated", "done",
		"last_check: C001\n"+currentFreshness("C001"))
	writeCheck(t, root, "C002", "task", "T001", "CLEAR")
	writeTask(t, root, "O001", "T001", "Finished work", "status: done\nlast_check: C002\n"+currentFreshness("C002"))

	// O002: every owned Task done, no integration Check of its own.
	writeObjective(t, root, "O002", "Every task done, nothing integrated", "in_progress")
	writeTask(t, root, "O002", "T002", "Also finished", "status: done\n")

	// O003: waiting on O002, which is not done.
	writeObjectiveExtra(t, root, "O003", "Waiting on the second", "planned", "depends_on: [\"O002\"]\n")
	writeTask(t, root, "O003", "T003", "Blocked by an objective wait", "status: planned\n")
	// Owned by O003, filed under O001's directory.
	writeTaskInDir(t, root, "O001", "O003", "T004", "Filed somewhere else entirely", "status: planned\n")

	// O004, O005, O006: the remaining clearance states.
	writeCheck(t, root, "C003", "objective", "O004", "NEEDS WORK")
	writeObjectiveExtra(t, root, "O004", "Integration found problems", "in_progress", "last_check: C003\n")
	writeCheck(t, root, "C004", "objective", "O005", "CLEAR")
	writeObjectiveExtra(t, root, "O005", "Integration never assessed", "in_progress", "last_check: C004\n")
	writeCheck(t, root, "C005", "objective", "O006", "CLEAR")
	writeObjectiveExtra(t, root, "O006", "Integration has gone stale", "in_progress",
		"last_check: C005\n"+staleFreshness("C005"))

	return root
}

// writeObjectiveDependencyProject builds a project whose second Objective waits
// on its first. The Task under that second Objective declares no dependency of
// its own, so anything the board says about a wait there came from the gate
// decision rather than from the Task's fields.
func writeObjectiveDependencyProject(t *testing.T) string {
	t.Helper()
	root := savepointRoot(t)
	writeConfig(t, root)
	writeRouter(t, root, "task", "O002", "T002")
	writeObjective(t, root, "O001", "Comes first", "in_progress")
	writeTask(t, root, "O001", "T001", "The work the other objective waits on", "status: planned\n")

	testutil.WriteFile(t, objectivePath(root, "O002"),
		"---\nid: O002\ntitle: \"Comes second\"\nstatus: planned\ndepends_on: [\"O001\"]\n---\n\n# Comes second\n")
	writeTask(t, root, "O002", "T002", "Blocked by its own objective's wait", "status: planned\n")

	return root
}

// writeIssue writes one open Issue recorded by a Check, naming the Task
// carrying its repair and the Check that observed it. The Check deliberately
// does not name the Issue back: that direction is the legal asymmetry, and it
// keeps the fixture one record shorter than the paired shape.
func writeIssue(t *testing.T, root, id, title, issueType, originCheck, task string) {
	t.Helper()
	testutil.WriteFile(t, filepath.Join(root, "issues", id+".md"),
		"---\nid: "+id+"\ntitle: \""+title+"\"\ntype: "+issueType+"\nstatus: open\n"+
			"source: {kind: check, check: "+originCheck+
			", actor: {role: checker, session: checker-fixture}, at: '2026-01-02T00:00:00Z'}\n"+
			"tasks: ["+task+"]\nchecks: ["+originCheck+"]\n---\n\n# "+title+"\n")
}

// writeTaskBody writes a Task with a body of its own, for the fixtures that
// prove the overlay displays author-owned markdown without reading it.
func writeTaskBody(t *testing.T, root, objective, id, title, extra, body string) {
	t.Helper()
	testutil.WriteFile(t, taskPath(root, objective, id),
		"---\nid: "+id+"\ntitle: \""+title+"\"\nobjective: "+objective+
			"\nplanned_by: {role: planner, session: board-fixture}\n"+extra+"---\n\n"+body)
}

// acceptedByOwner is the owner_validation block of a record whose owner has
// accepted a named Check.
func acceptedByOwner(check string) string {
	return "owner_validation:\n  required: true\n  accepted_check: " + check +
		"\n  accepted_by: {role: owner, session: owner-fixture}\n"
}

// writeEvidenceProject builds the project the detail surfaces exist for: an
// Objective with its own superseded integration chain and an unsatisfied
// Objective dependency, and Tasks covering a rerun chain, a satisfied
// accepted-level dependency, an unsatisfied clear-level one, a replan, a
// recorded exception, an outstanding owner wait, a record with no Check at
// all, and a linked Issue.
func writeEvidenceProject(t *testing.T) string {
	t.Helper()
	root := savepointRoot(t)
	writeConfig(t, root)
	writeRouter(t, root, "task", "O001", "T002")

	// O001's integration was checked twice; the rerun supersedes the first run.
	writeCheck(t, root, "C010", "objective", "O001", "CLEAR")
	writeCheckExtra(t, root, "C011", "objective", "O001", "CLEAR", "supersedes: C010\n")
	writeObjectiveExtra(t, root, "O001", "Ship the evidence surface", "in_progress",
		"depends_on: [\"O002\"]\nlast_check: C011\n"+currentFreshness("C011"))
	writeObjective(t, root, "O002", "Groundwork nobody has started", "planned")

	writeCheck(t, root, "C001", "task", "T001", "CLEAR")
	writeTask(t, root, "O001", "T001", "Cleared and accepted by the owner",
		"status: done\nlast_check: C001\n"+currentFreshness("C001")+acceptedByOwner("C001"))

	// T002's first run found problems; the rerun supersedes it.
	writeCheck(t, root, "C002", "task", "T002", "NEEDS WORK")
	writeCheckExtra(t, root, "C003", "task", "T002", "CLEAR", "supersedes: C002\n")
	writeTaskBody(t, root, "O001", "T002", "Rechecked after the first run found problems",
		"status: in_progress\nstage: audit\ndepends_on:\n  - {task: T001, requires: accepted}\n"+
			"last_check: C003\n"+currentFreshness("C003"),
		"# Outcome\n\nThe board shows the record behind the badge.\n\n- [x] status: done\n- [ ] still open\n")

	writeTask(t, root, "O001", "T003", "Waiting on work that is not done",
		"status: planned\ndepends_on:\n  - {task: T004, requires: clear}\n")
	writeTask(t, root, "O001", "T004", "Flagged for replan",
		"status: in_progress\nstage: build\nreplan:\n  reason: \"the approach no longer fits\"\n"+
			"  recorded_by: {role: planner, session: planner-fixture}\n  recorded_at: 2026-01-04T00:00:00Z\n")

	writeCheck(t, root, "C004", "task", "T005", "NEEDS WORK")
	writeTask(t, root, "O001", "T005", "Closed under a recorded exception",
		"status: done\nlast_check: C004\nexception:\n  requirements: [\"TEST-02\", \"TEST-05\"]\n"+
			"  reason: \"shipped with a known gap\"\n  owner: \"the owner\"\n"+
			"  recorded_at: 2026-01-05T00:00:00Z\n  check: C004\n")

	writeCheck(t, root, "C005", "task", "T006", "CLEAR")
	writeTask(t, root, "O001", "T006", "Waiting on the owner to accept",
		"status: in_progress\nstage: audit\nlast_check: C005\n"+currentFreshness("C005")+
			"owner_validation:\n  required: true\n")

	writeTask(t, root, "O001", "T007", "Nothing recorded against it yet", "status: planned\n")

	writeIssue(t, root, "I001", "The retry loop drops the last attempt", "defect", "C002", "T002")

	return root
}

// invalidProjectCase is one structural refusal LoadV2Index makes over an
// otherwise valid project, with the file path and wording its diagnostic must
// name. The cases are shared so every surface that reports a refused load — the
// screen, and the non-TTY path — is proven against the same four shapes.
type invalidProjectCase struct {
	Name      string
	Build     func(t *testing.T, root string)
	WantPath  string
	WantParts []string
}

func invalidProjectCases() []invalidProjectCase {
	return []invalidProjectCase{
		{
			Name:      "task with no title",
			Build:     writeTitlelessTask,
			WantPath:  filepath.Join("objectives", "O001-fixture", "tasks", "T001-fixture.md"),
			WantParts: []string{"missing required field title"},
		},
		{
			Name: "duplicate objective id",
			Build: func(t *testing.T, root string) {
				testutil.WriteFile(t, filepath.Join(root, "objectives", "O001-again", "Objective.md"),
					"---\nid: O001\ntitle: \"Second claim on O001\"\nstatus: planned\n---\n\n# Second\n")
			},
			WantPath:  filepath.Join("objectives", "O001-again", "Objective.md"),
			WantParts: []string{"declared more than once", "O001"},
		},
		{
			Name: "task owned by a missing objective",
			Build: func(t *testing.T, root string) {
				testutil.WriteFile(t, taskPath(root, "O001", "T003"),
					"---\nid: T003\ntitle: \"Orphan work\"\nobjective: O999\nplanned_by: {role: planner, session: board-fixture}\nstatus: planned\n---\n\n# Orphan work\n")
			},
			WantPath:  filepath.Join("objectives", "O001-fixture", "tasks", "T003-fixture.md"),
			WantParts: []string{"missing objective", "O999"},
		},
		{
			Name: "task depending on a record that does not exist",
			Build: func(t *testing.T, root string) {
				writeTask(t, root, "O001", "T004", "Blocked work", "status: planned\ndepends_on:\n  - {task: T900, requires: clear}\n")
			},
			WantPath:  filepath.Join("objectives", "O001-fixture", "tasks", "T004-fixture.md"),
			WantParts: []string{"T900"},
		},
	}
}

// writeTitlelessTask replaces the fixture's T001 with a Task carrying no
// title. DecodeTaskV2 refuses it, so a project containing it does not load at
// all — which is what makes the epic's title requirement structural.
func writeTitlelessTask(t *testing.T, root string) {
	t.Helper()
	testutil.WriteFile(t, taskPath(root, "O001", "T001"),
		"---\nid: T001\nobjective: O001\nplanned_by: {role: planner, session: board-fixture}\nstatus: planned\n---\n\n# Untitled\n")
}

// createPendingOperation records an incomplete migration operation over the
// project root holding root, and returns its operation ID.
func createPendingOperation(t *testing.T, root string) string {
	t.Helper()
	operation, err := migrate.CreateOperation(filepath.Dir(root), "op-test-board-fixture", nil, nil, time.Now())
	if err != nil {
		t.Fatalf("CreateOperation() error = %v", err)
	}
	return operation.Journal.OperationID
}

// removeTask deletes a Task file, for tests that reload a project whose records
// changed underneath the board.
func removeTask(t *testing.T, root, objective, id string) {
	t.Helper()
	if err := os.Remove(taskPath(root, objective, id)); err != nil {
		t.Fatalf("remove task %s: %v", id, err)
	}
}

// removeObjective and removeCheck delete a record, for tests that reload a
// project whose records changed underneath the board.
func removeObjective(t *testing.T, root, id string) {
	t.Helper()
	if err := os.RemoveAll(filepath.Dir(objectivePath(root, id))); err != nil {
		t.Fatalf("remove objective %s: %v", id, err)
	}
}

func removeCheck(t *testing.T, root, id string) {
	t.Helper()
	if err := os.Remove(filepath.Join(root, "checks", id+".md")); err != nil {
		t.Fatalf("remove check %s: %v", id, err)
	}
}

func objectivePath(root, id string) string {
	return filepath.Join(root, "objectives", id+"-fixture", "Objective.md")
}

func taskPath(root, objective, id string) string {
	return filepath.Join(root, "objectives", objective+"-fixture", "tasks", id+"-fixture.md")
}
