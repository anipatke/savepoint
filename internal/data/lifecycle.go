package data

import "fmt"

func ValidateTaskLifecycle(task Task) error {
	if !IsCanonicalColumn(task.Column) {
		return fmt.Errorf("invalid task status %q: use planned, in_progress, or done", task.Column)
	}

	if task.Column != ColumnInProgress && task.Stage != "" {
		return fmt.Errorf("phase %q is only valid when status is in_progress", task.Stage)
	}

	if task.Column == ColumnInProgress && !IsCanonicalStage(task.Stage) {
		return fmt.Errorf("invalid in_progress phase %q: use build, test, or audit", task.Stage)
	}

	return nil
}

func IsCanonicalColumn(value ColumnType) bool {
	switch value {
	case ColumnPlanned, ColumnInProgress, ColumnDone:
		return true
	default:
		return false
	}
}

func IsCanonicalStage(value ProgressStage) bool {
	switch value {
	case StageBuild, StageTest, StageAudit:
		return true
	default:
		return false
	}
}
