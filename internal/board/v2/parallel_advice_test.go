package v2

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/resume"
	"github.com/opencode/savepoint/internal/testutil"
)

// writeAdviceProject writes O-001 with two independent lanes whose heads can
// start together (T-002 board, T-004 core), a second core Task behind T-004,
// a Task blocked on an unfinished dependency, and a Task whose title carries a
// terminal escape.
func writeAdviceProject(t *testing.T, enabled bool) string {
	t.Helper()
	root := writeValidProject(t)
	config := "schema_version: 2\n"
	if enabled {
		config += "features:\n  parallel_planning: true\n"
	}
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), config)
	writeFixtureRouter(t, root, "task", "O-001", "T-002")
	writeFixtureObjective(t, root, "O-001", "First objective", "in_progress", laneObjectiveExtra)
	writeTask(t, root, "O-001", "T-001", "Sequential work", "status: planned\n")
	writeTask(t, root, "O-001", "T-002", "Board work", "status: planned\nlane: board\nplanned_reads: [a.go]\nplanned_writes: [b.go]\ndepends_on: [{task: T-009, requires: clear}]\n")
	writeTask(t, root, "O-001", "T-009", "Finished prerequisite", "status: done\n")
	writeCheck(t, root, "C-001", "task", "T-009", "CLEAR")
	writeTask(t, root, "O-001", "T-004", "Core \\x1b[31mred\\x07 work", "status: planned\nlane: core\nplanned_reads: [c.go]\nplanned_writes: [d.go]\n")
	writeTask(t, root, "O-001", "T-006", "Core follow-up", "status: planned\nlane: core\nplanned_reads: [e.go]\nplanned_writes: [f.go]\n")
	writeTask(t, root, "O-001", "T-007", "Waiting work", "status: planned\nlane: board\nplanned_reads: [g.go]\nplanned_writes: [h.go]\ndepends_on: [{task: T-004, requires: clear}]\n")
	return root
}

func instructionBlock(t *testing.T, text, id string) string {
	t.Helper()
	start := strings.Index(text, "----- begin "+id+" -----")
	end := strings.Index(text, "----- end "+id+" -----")
	if start < 0 || end < start {
		t.Fatalf("no instruction block for %s in:\n%s", id, text)
	}
	return text[start : end+len("----- end "+id+" -----")]
}

func adviceSurfaces(t *testing.T, root string) (resumeText, objectiveText, taskText, plainText string) {
	t.Helper()
	state := loadProject(filepath.Clean(root)).State
	var out strings.Builder
	if err := resume.Render(&out, state.Next, state.Index); err != nil {
		t.Fatalf("resume.Render: %v", err)
	}
	objective, _ := newObjectiveDetail(state.Index, state.Health, "O-001")
	task, _ := newTaskDetail(state.Index, "T-002")
	enabled := state.Features.ParallelPlanning
	objective = withParallel(withEnabled(state, enabled), objective)
	task = withParallel(withEnabled(state, enabled), task)
	return out.String(), strings.Join(objective.Parallel, "\n"), strings.Join(task.Parallel, "\n"), renderPlain(state, "O-001")
}

func TestParallelAdviceAgreesAcrossResumeDetailsAndPlainBoard(t *testing.T) {
	resumeText, objectiveText, taskText, plainText := adviceSurfaces(t, writeAdviceProject(t, true))

	want := instructionBlock(t, objectiveText, "T-002")
	for name, text := range map[string]string{"resume": resumeText, "task detail": taskText, "plain board": plainText} {
		if got := instructionBlock(t, text, "T-002"); got != want {
			t.Errorf("%s instruction differs from the Objective detail:\n%s\n--- want ---\n%s", name, got, want)
		}
	}
	for _, line := range []string{
		"Start T-002 — Board work (O-001)",
		"Use savepoint-task for the Start line above.",
		"Objective: O-001; Task: T-002; lane: board (Lane / Proposed worktree — Board).",
		"Reads: a.go",
		"Writes: b.go",
		"Prerequisites: T-009 (requires clear).",
		"leave router.md and the Goal selection unchanged",
		"do not allocate Task, Check or Issue identities",
		"do not push or merge. Merging and Checks happen on main.",
		"Running the Task on main or in another worktree is fine",
	} {
		if !strings.Contains(want, line) {
			t.Errorf("instruction lacks %q:\n%s", line, want)
		}
	}
	for _, text := range []string{want, objectiveText} {
		if strings.Contains(text, "/home/") || strings.Contains(text, "git worktree") {
			t.Errorf("advice guesses a path or runs git:\n%s", text)
		}
	}
	for name, text := range map[string]string{"objective": objectiveText, "plain": plainText, "resume": resumeText} {
		if !strings.Contains(text, "May start together: T-004, T-002") && !strings.Contains(text, "May start together: T-002, T-004") {
			t.Errorf("%s does not say the two heads may start together:\n%s", name, text)
		}
	}
}

func TestParallelAdviceExplainsBlockedAndLaterTasksWithoutAdvertisingThem(t *testing.T) {
	_, objectiveText, _, _ := adviceSurfaces(t, writeAdviceProject(t, true))

	for _, id := range []string{"T-007", "T-006"} {
		if strings.Contains(objectiveText, "begin "+id) || strings.Contains(objectiveText, "Start "+id) {
			t.Errorf("%s is advertised with a Start instruction:\n%s", id, objectiveText)
		}
	}
	for _, want := range []string{"Not suggested now: T-007 cannot start yet", "Not suggested now: T-004 is suggested first in lane core"} {
		if !strings.Contains(objectiveText, want) {
			t.Errorf("advice lacks reason %q:\n%s", want, objectiveText)
		}
	}
}

func TestParallelAdviceIsAbsentWhileTheOptionIsOffAndKeepsRecords(t *testing.T) {
	root := writeAdviceProject(t, false)
	resumeText, objectiveText, taskText, plainText := adviceSurfaces(t, root)

	for name, text := range map[string]string{"resume": resumeText, "objective": objectiveText, "task": taskText, "plain": plainText} {
		if strings.Contains(text, "Parallel planning") || strings.Contains(text, "savepoint-task. Start") {
			t.Errorf("%s shows advice with the option off:\n%s", name, text)
		}
	}
	if !strings.Contains(resumeText, "Start T-002") || !strings.Contains(resumeText, "Next action: Start Task T-002.") {
		t.Errorf("Next meaning changed with the option off:\n%s", resumeText)
	}
	state := loadProject(root).State
	if task := state.Index.Tasks["T-002"]; task.Plan.Lane != "board" || !task.Plan.Writes.Contains("b.go") {
		t.Errorf("saved plan was dropped: %+v", task.Plan)
	}
}

func TestParallelAdviceKeepsTheNextLineAndStripsTerminalControls(t *testing.T) {
	enabled, _, _, plainOn := adviceSurfaces(t, writeAdviceProject(t, true))
	disabled, _, _, _ := adviceSurfaces(t, writeAdviceProject(t, false))

	if first := func(s string) string { return strings.SplitN(s, "\n", 2)[0] }; first(enabled) != first(disabled) {
		t.Errorf("Next line changed: %q vs %q", first(enabled), first(disabled))
	}
	state := loadProject(writeAdviceProject(t, true)).State
	objective, _ := newObjectiveDetail(state.Index, state.Health, "O-001")
	objective = withParallel(withEnabled(state, true), objective)
	for name, text := range map[string]string{"resume": enabled, "detail": strings.Join(objective.Parallel, "\n"), "plain": plainOn} {
		if strings.ContainsAny(text, "\x1b\x07") {
			t.Errorf("%s output carries a terminal control: %q", name, text)
		}
	}
	if !strings.Contains(strings.Join(objective.Parallel, "\n"), "Core red work") {
		t.Errorf("sanitised title missing:\n%s", strings.Join(objective.Parallel, "\n"))
	}
}

func TestParallelAdviceNamesNoCrossObjectiveOpportunityGoalWide(t *testing.T) {
	root := writeAdviceProject(t, true)
	state := loadProject(root).State
	if plain := renderPlain(state, ""); strings.Contains(plain, "Parallel planning") || strings.Contains(plain, "May start together") {
		t.Errorf("Goal-wide board compares Objectives:\n%s", plain)
	}
}

func TestParallelAdviceWithholdsInstructionsForABlockedFocusedTask(t *testing.T) {
	state := loadProject(writeAdviceProject(t, true)).State
	detail, _ := newTaskDetail(state.Index, "T-007")
	text := strings.Join(withParallel(withEnabled(state, true), detail).Parallel, "\n")
	if strings.Contains(text, "begin T-007") || !strings.Contains(text, "T-007 cannot start yet") {
		t.Errorf("blocked Task detail = %q, want the constraint and no Start instruction", text)
	}
}

// withEnabled returns state with the parallel-planning preference set.
func withEnabled(state ProjectState, enabled bool) ProjectState {
	state.Features.ParallelPlanning = enabled
	return state
}
