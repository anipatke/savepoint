package board

import (
	"fmt"

	"github.com/opencode/savepoint/internal/data"
)

// Advance moves a task forward through the phase lifecycle.
func Advance(t *data.Task) {
	stage := t.Stage
	if stage == "" {
		stage = data.StageBuild
	}
	switch t.Column {
	case data.ColumnPlanned:
		t.Column = data.ColumnInProgress
		t.Stage = data.StageBuild
	case data.ColumnInProgress:
		switch stage {
		case data.StageBuild:
			t.Stage = data.StageTest
		case data.StageTest:
			t.Stage = data.StageAudit
		case data.StageAudit:
			t.Column = data.ColumnDone
			t.Stage = ""
		}
	}
	syncTaskStatus(t)
}

// Retreat moves a task backward through the phase lifecycle.
func Retreat(t *data.Task) {
	stage := t.Stage
	if stage == "" {
		stage = data.StageBuild
	}
	switch t.Column {
	case data.ColumnDone:
		t.Column = data.ColumnInProgress
		t.Stage = data.StageAudit
	case data.ColumnInProgress:
		switch stage {
		case data.StageAudit:
			t.Stage = data.StageTest
		case data.StageTest:
			t.Stage = data.StageBuild
		case data.StageBuild:
			t.Column = data.ColumnPlanned
			t.Stage = ""
		}
	}
	syncTaskStatus(t)
}

func syncTaskStatus(t *data.Task) {
	t.Status = string(t.Column)
}

func taskTransitionMessage(prefix string, task data.Task) string {
	if task.Column == data.ColumnInProgress {
		return fmt.Sprintf("%s %s to %s", prefix, shortID(task.ID), task.Stage)
	}
	return fmt.Sprintf("%s %s to %s", prefix, shortID(task.ID), task.Column)
}

// CanAdvance checks whether a task is allowed to advance to its next phase.
// It validates phase adjacency and dependency completion.
// Returns (true, "") if allowed, or (false, reason) if blocked.
func CanAdvance(t *data.Task, allTasks []data.Task) (bool, string) {
	switch t.Column {
	case data.ColumnPlanned:
		return dependenciesDone(t, allTasks)
	case data.ColumnInProgress:
		stage := t.Stage
		if stage == "" {
			stage = data.StageBuild
		}
		switch stage {
		case data.StageBuild:
			return true, ""
		case data.StageTest:
			return true, ""
		case data.StageAudit:
			return dependenciesDone(t, allTasks)
		default:
			return false, fmt.Sprintf("unknown stage %q", stage)
		}
	case data.ColumnDone:
		return false, "task is already done"
	default:
		return false, fmt.Sprintf("unknown column %q", t.Column)
	}
}

func dependenciesDone(t *data.Task, allTasks []data.Task) (bool, string) {
	for _, depID := range t.DependsOn {
		dep := findTask(depID, allTasks)
		if dep == nil {
			return false, fmt.Sprintf("dependency %q not found", depID)
		}
		if dep.Column != data.ColumnDone {
			return false, fmt.Sprintf("dependency %q is not done", depID)
		}
	}
	return true, ""
}

func findTask(id string, tasks []data.Task) *data.Task {
	for i := range tasks {
		if tasks[i].ID == id {
			return &tasks[i]
		}
	}
	return nil
}
