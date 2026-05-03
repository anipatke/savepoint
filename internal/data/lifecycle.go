package data

import "fmt"

func ValidateTaskLifecycle(task *Task) error {
	if !IsCanonicalColumn(task.Column) {
		return fmt.Errorf("invalid status %q: use planned, in_progress, or done. Add 'status: planned' or 'status: in_progress' to task frontmatter", task.Column)
	}

	if task.Column == ColumnInProgress {
		if task.Stage == "" {
			task.Stage = StageBuild
			return nil
		}
		if !IsCanonicalStage(task.Stage) {
			return fmt.Errorf("invalid phase %q: use build, test, or audit. Add 'phase: build' to task frontmatter", task.Stage)
		}
		return nil
	}

	if task.Stage != "" {
		return fmt.Errorf("phase field %q is only valid when status is in_progress. Remove 'phase' or change status to in_progress", task.Stage)
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
