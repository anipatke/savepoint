package data

import "testing"

func TestValidateTaskLifecycle_allowsPlannedWithoutPhase(t *testing.T) {
	task := Task{Column: ColumnPlanned}
	if err := ValidateTaskLifecycle(&task); err != nil {
		t.Fatalf("ValidateTaskLifecycle() error = %v", err)
	}
}

func TestValidateTaskLifecycle_defaultsInProgressWithoutPhase(t *testing.T) {
	task := Task{Column: ColumnInProgress}
	if err := ValidateTaskLifecycle(&task); err != nil {
		t.Fatalf("ValidateTaskLifecycle() error = %v", err)
	}
	if task.Stage != StageBuild {
		t.Fatalf("Task.Stage = %q, want %q", task.Stage, StageBuild)
	}
}

func TestValidateTaskLifecycle_allowsInProgressWithPhase(t *testing.T) {
	task := Task{Column: ColumnInProgress, Stage: StageAudit}
	if err := ValidateTaskLifecycle(&task); err != nil {
		t.Fatalf("ValidateTaskLifecycle() error = %v", err)
	}
}

func TestValidateTaskLifecycle_rejectsUnknownStatus(t *testing.T) {
	task := Task{Column: "review"}
	if err := ValidateTaskLifecycle(&task); err == nil {
		t.Fatal("ValidateTaskLifecycle() expected unknown status error")
	}
}

func TestValidateTaskLifecycle_rejectsPhaseOutsideInProgress(t *testing.T) {
	task := Task{Column: ColumnPlanned, Stage: StageBuild}
	if err := ValidateTaskLifecycle(&task); err == nil {
		t.Fatal("ValidateTaskLifecycle() expected phase/status error")
	}
}
