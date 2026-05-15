package data

import "fmt"

func ValidateTaskLifecycle(task *Task) error {
	if err := ValidateComplexity(task.ComplexityTier, task.ComplexityReason); err != nil {
		return err
	}

	if !IsCanonicalColumn(task.Column) {
		return fmt.Errorf("invalid status %q: use planned, in_progress, or done. Add 'status: planned' or 'status: in_progress' to task frontmatter", task.Column)
	}

	if task.Column == ColumnInProgress {
		if task.Stage == "" {
			return fmt.Errorf("stage is required when task status is in_progress. Add 'stage: build' to task frontmatter")
		}
		if !IsCanonicalStage(task.Stage) {
			return fmt.Errorf("invalid stage %q: use build, test, or audit. Add 'stage: build' to task frontmatter", task.Stage)
		}
		return nil
	}

	if task.Stage != "" {
		return fmt.Errorf("stage field %q is only valid when status is in_progress. Remove 'stage' or change status to in_progress", task.Stage)
	}

	return nil
}

// ValidateComplexity checks that complexity_tier and complexity_reason are
// mutually consistent and within allowed bounds.
func ValidateComplexity(tier ComplexityTier, reason string) error {
	if tier == "" && reason == "" {
		return nil
	}
	if tier != "" && !IsCanonicalComplexityTier(tier) {
		return fmt.Errorf("invalid complexity_tier %q: use low, medium, high, or spike", tier)
	}
	if tier != "" && reason == "" {
		return fmt.Errorf("complexity_reason is required when complexity_tier is set")
	}
	if reason != "" && tier == "" {
		return fmt.Errorf("complexity_tier is required when complexity_reason is set")
	}
	if len(reason) > MaxComplexityReasonLen {
		return fmt.Errorf("complexity_reason exceeds %d characters", MaxComplexityReasonLen)
	}
	return nil
}

func validateDefectLifecycle(d *Defect) error {
	if d.Status == "" {
		d.Status = DefectOpen
	}
	if !IsCanonicalDefectStatus(d.Status) {
		return fmt.Errorf("invalid defect status %q: use open, in_progress, or resolved", d.Status)
	}
	if d.Status == DefectInProgress {
		if d.Stage == "" {
			return fmt.Errorf("stage is required when defect status is in_progress")
		}
		if !IsCanonicalStage(d.Stage) {
			return fmt.Errorf("invalid stage %q: use build, test, or audit", d.Stage)
		}
		return nil
	}
	if d.Stage != "" {
		return fmt.Errorf("stage field %q is only valid when defect status is in_progress", d.Stage)
	}
	return nil
}

func IsCanonicalDefectStatus(value DefectStatus) bool {
	switch value {
	case DefectOpen, DefectInProgress, DefectResolved:
		return true
	default:
		return false
	}
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

func IsCanonicalComplexityTier(tier ComplexityTier) bool {
	switch tier {
	case ComplexityLow, ComplexityMedium, ComplexityHigh, ComplexitySpike:
		return true
	default:
		return false
	}
}
