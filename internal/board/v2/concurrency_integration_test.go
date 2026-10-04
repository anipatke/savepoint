package v2

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/testutil"
)

// These tests change a temporary project between loads and require every
// surface — resume, the Objective and Task details and the plain board — to
// keep agreeing on the selected Objective's advice.

func TestIntegrationSurfacesAgreeAfterStatusChanges(t *testing.T) {
	root := writeAdviceProject(t, true)

	// T-004's lane head starts, so T-006 is still behind it and T-002 stays
	// suggested beside the running Task.
	writeTask(t, root, "O-001", "T-004", "Core work", "status: in_progress\nstage: build\nlane: core\nplanned_reads: [c.go]\nplanned_writes: [d.go]\n")
	resumeText, objectiveText, taskText, plainText := adviceSurfaces(t, root)
	// Resume is focused on the selected Task, so only the Objective-wide
	// surfaces list the running work; none may advertise T-006.
	for name, text := range map[string]string{"objective": objectiveText, "plain": plainText} {
		if !strings.Contains(text, "In progress: T-004") {
			t.Errorf("%s lacks the running Task:\n%s", name, text)
		}
	}
	for name, text := range map[string]string{"resume": resumeText, "objective": objectiveText, "plain": plainText} {
		if strings.Contains(text, "begin T-006") {
			t.Errorf("%s advertises T-006 while T-004 holds its lane:\n%s", name, text)
		}
	}
	want := instructionBlock(t, objectiveText, "T-002")
	for name, text := range map[string]string{"resume": resumeText, "task detail": taskText, "plain": plainText} {
		if got := instructionBlock(t, text, "T-002"); got != want {
			t.Errorf("%s instruction for T-002 differs after the status change:\n%s\n--- want ---\n%s", name, got, want)
		}
	}

	// The lane head finishes: T-006 becomes the lane's next Task, and every
	// surface says so.
	writeTask(t, root, "O-001", "T-004", "Core work", "status: done\nlane: core\nplanned_reads: [c.go]\nplanned_writes: [d.go]\n")
	_, objectiveText, _, plainText = adviceSurfaces(t, root)
	if !strings.Contains(objectiveText, "begin T-006") || !strings.Contains(plainText, "begin T-006") {
		t.Errorf("T-006 is not suggested once its lane head is done:\n%s", objectiveText)
	}
}

func TestIntegrationSurfacesWithholdTogetherAfterAReplan(t *testing.T) {
	root := writeAdviceProject(t, true)
	writeTask(t, root, "O-001", "T-006", "Core follow-up",
		"status: in_progress\nstage: build\nlane: core\nplanned_reads: [e.go]\nplanned_writes: [f.go]\n"+
			"replan:\n  reason: \"scope changed\"\n  recorded_by: {role: planner, session: planner-fixture}\n  recorded_at: 2026-10-04T00:00:00Z\n")
	resumeText, objectiveText, taskText, plainText := adviceSurfaces(t, root)
	for name, text := range map[string]string{"resume": resumeText, "objective": objectiveText, "task detail": taskText, "plain": plainText} {
		if !strings.Contains(text, "records a replan") {
			t.Errorf("%s does not explain the withheld advice:\n%s", name, text)
		}
		if strings.Contains(text, "----- begin") {
			t.Errorf("%s still offers a Start instruction during a replan:\n%s", name, text)
		}
	}
}

func TestIntegrationMalformedAdviceDegradesWithoutHidingTheBoard(t *testing.T) {
	root := writeAdviceProject(t, true)
	writeTask(t, root, "O-001", "T-002", "Board work", "status: planned\nlane: board\nplanned_reads: [a.go]\nplanned_writes: not-a-list\n")

	project := loadProject(filepath.Clean(root))
	if project.Diagnostic != "" {
		t.Fatalf("malformed advice refused the board: %s", project.Diagnostic)
	}
	resumeText, objectiveText, _, plainText := adviceSurfaces(t, root)
	for name, text := range map[string]string{"resume": resumeText, "plain": plainText} {
		if !strings.Contains(text, "T-002") {
			t.Errorf("%s dropped the Task with malformed advice:\n%s", name, text)
		}
	}
	if strings.Contains(objectiveText, "begin T-002") {
		t.Errorf("a Task with unusable advice is still offered a Start instruction:\n%s", objectiveText)
	}
	if !strings.Contains(resumeText, "Start T-002") {
		t.Errorf("ordinary Next changed because advice was malformed:\n%s", resumeText)
	}
}

func TestIntegrationGoalWideViewNeverImpliesCrossObjectiveConcurrency(t *testing.T) {
	root := writeAdviceProject(t, true)
	writeFixtureObjective(t, root, "O-002", "Second objective", "in_progress", laneObjectiveExtra)
	writeTask(t, root, "O-002", "T-020", "Other work", "status: planned\nlane: core\nplanned_reads: [x.go]\nplanned_writes: [y.go]\n")
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\nfeatures:\n  parallel_planning: true\n")

	state := loadProject(filepath.Clean(root)).State
	goalWide := renderPlain(state, "")
	for _, text := range []string{"Parallel planning", "May start together", "begin T-020", "begin T-002"} {
		if strings.Contains(goalWide, text) {
			t.Errorf("Goal-wide board contains %q:\n%s", text, goalWide)
		}
	}
	selected := renderPlain(state, "O-002")
	if !strings.Contains(selected, "T-020") || strings.Contains(selected, "begin T-002") {
		t.Errorf("selected-Objective board mixes Objectives:\n%s", selected)
	}
}
