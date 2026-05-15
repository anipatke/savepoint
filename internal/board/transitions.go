package board

import (
	"fmt"

	"github.com/opencode/savepoint/internal/data"
)

// Advance moves a task forward through the task lifecycle.
func Advance(t *data.Task) error {
	_, err := data.AdvanceTaskLifecycle(t)
	return err
}

// Retreat moves a task backward through the task lifecycle.
func Retreat(t *data.Task) error {
	_, err := data.RetreatTaskLifecycle(t)
	return err
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
func CanAdvance(t *data.Task, allTasks []data.Task, epicStatuses ...map[string]string) (bool, string) {
	scopedEpicStatuses := map[string]string(nil)
	if len(epicStatuses) > 0 {
		scopedEpicStatuses = epicStatuses[0]
	}
	switch t.Column {
	case data.ColumnPlanned:
		return dependenciesDone(t, allTasks, scopedEpicStatuses)
	case data.ColumnInProgress:
		next, err := data.AdvanceTaskLifecycleState(data.TaskLifecycleStateFromTask(*t))
		if err != nil {
			return false, err.Error()
		}
		switch next.Status {
		case data.ColumnInProgress:
			return true, ""
		case data.ColumnDone:
			return dependenciesDone(t, allTasks, scopedEpicStatuses)
		default:
			return false, fmt.Sprintf("unknown column %q", next.Status)
		}
	case data.ColumnDone:
		return false, "task is already done"
	default:
		return false, fmt.Sprintf("unknown column %q", t.Column)
	}
}

func dependenciesDone(t *data.Task, allTasks []data.Task, epicStatuses map[string]string) (bool, string) {
	for _, depID := range t.DependsOn {
		dep := data.ResolveDependency(depID, *t, allTasks, epicStatuses)
		switch dep.Kind {
		case data.DependencyTask:
			if dep.TaskStatus != data.ColumnDone {
				return false, fmt.Sprintf("dependency %q is not done", depID)
			}
		case data.DependencyEpic:
			if dep.EpicStatus != "audited" {
				return false, fmt.Sprintf("dependency %q is not audited", depID)
			}
		default:
			return false, fmt.Sprintf("dependency %q not found", depID)
		}
	}
	return true, ""
}
