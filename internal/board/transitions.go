package board

import (
	"fmt"

	"github.com/opencode/savepoint/internal/data"
)

// Advance moves a task forward through the phase lifecycle.
func Advance(t *data.Task) {
	switch t.Column {
	case data.ColumnPlanned:
		t.Column = data.ColumnInProgress
		t.Stage = data.StageBuild
	case data.ColumnInProgress:
		switch t.Stage {
		case data.StageBuild:
			t.Stage = data.StageTest
		case data.StageTest:
			t.Stage = data.StageAudit
		case data.StageAudit:
			t.Column = data.ColumnDone
			t.Stage = ""
		}
	}
}

// Retreat moves a task backward through the phase lifecycle.
func Retreat(t *data.Task) {
	switch t.Column {
	case data.ColumnDone:
		t.Column = data.ColumnInProgress
		t.Stage = data.StageAudit
	case data.ColumnInProgress:
		switch t.Stage {
		case data.StageAudit:
			t.Stage = data.StageTest
		case data.StageTest:
			t.Stage = data.StageBuild
		case data.StageBuild:
			t.Column = data.ColumnPlanned
			t.Stage = ""
		}
	}
}

// CanAdvance checks whether a task is allowed to advance to its next phase.
// It validates phase adjacency and dependency completion.
// Returns (true, "") if allowed, or (false, reason) if blocked.
func CanAdvance(t *data.Task, allTasks []data.Task) (bool, string) {
	switch t.Column {
	case data.ColumnPlanned:
		return true, ""
	case data.ColumnInProgress:
		switch t.Stage {
		case data.StageBuild:
			return true, ""
		case data.StageTest:
			return true, ""
		case data.StageAudit:
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
		default:
			return false, fmt.Sprintf("unknown stage %q", t.Stage)
		}
	case data.ColumnDone:
		return false, "task is already done"
	default:
		return false, fmt.Sprintf("unknown column %q", t.Column)
	}
}

func findTask(id string, tasks []data.Task) *data.Task {
	for i := range tasks {
		if tasks[i].ID == id {
			return &tasks[i]
		}
	}
	return nil
}
