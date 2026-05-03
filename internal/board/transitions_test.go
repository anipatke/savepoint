package board

import (
	"testing"

	"github.com/opencode/savepoint/internal/data"
)

func TestAdvance(t *testing.T) {
	tests := []struct {
		name       string
		initialCol data.ColumnType
		initialSt  data.ProgressStage
		expectCol  data.ColumnType
		expectSt   data.ProgressStage
	}{
		{"planned to in_progress/build", data.ColumnPlanned, "", data.ColumnInProgress, data.StageBuild},
		{"in_progress/build to test", data.ColumnInProgress, data.StageBuild, data.ColumnInProgress, data.StageTest},
		{"in_progress/test to audit", data.ColumnInProgress, data.StageTest, data.ColumnInProgress, data.StageAudit},
		{"in_progress/audit to done", data.ColumnInProgress, data.StageAudit, data.ColumnDone, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := data.Task{Column: tt.initialCol, Stage: tt.initialSt}
			Advance(&task)
			if task.Column != tt.expectCol || task.Stage != tt.expectSt {
				t.Errorf("Advance() = %v/%v, want %v/%v", task.Column, task.Stage, tt.expectCol, tt.expectSt)
			}
			if task.Status != string(tt.expectCol) {
				t.Errorf("Advance() status = %q, want %q", task.Status, tt.expectCol)
			}
		})
	}
}

func TestRetreat(t *testing.T) {
	tests := []struct {
		name       string
		initialCol data.ColumnType
		initialSt  data.ProgressStage
		expectCol  data.ColumnType
		expectSt   data.ProgressStage
	}{
		{"done to in_progress/audit", data.ColumnDone, "", data.ColumnInProgress, data.StageAudit},
		{"in_progress/audit to test", data.ColumnInProgress, data.StageAudit, data.ColumnInProgress, data.StageTest},
		{"in_progress/test to build", data.ColumnInProgress, data.StageTest, data.ColumnInProgress, data.StageBuild},
		{"in_progress/build to planned", data.ColumnInProgress, data.StageBuild, data.ColumnPlanned, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := data.Task{Status: string(tt.initialCol), Column: tt.initialCol, Stage: tt.initialSt}
			Retreat(&task)
			if task.Column != tt.expectCol || task.Stage != tt.expectSt {
				t.Errorf("Retreat() = %v/%v, want %v/%v", task.Column, task.Stage, tt.expectCol, tt.expectSt)
			}
			if task.Status != string(tt.expectCol) {
				t.Errorf("Retreat() status = %q, want %q", task.Status, tt.expectCol)
			}
		})
	}
}

func TestCanAdvance_plannedAllowedWhenDependenciesDone(t *testing.T) {
	allTasks := []data.Task{
		{ID: "T1", Column: data.ColumnPlanned, DependsOn: []string{"T2"}},
		{ID: "T2", Column: data.ColumnDone},
	}
	ok, reason := CanAdvance(&allTasks[0], allTasks)
	if !ok {
		t.Errorf("CanAdvance(planned with done dep) = false %q, want true", reason)
	}
}

func TestCanAdvance_plannedBlockedByDependency(t *testing.T) {
	allTasks := []data.Task{
		{ID: "T1", Column: data.ColumnPlanned, DependsOn: []string{"T2"}},
		{ID: "T2", Column: data.ColumnInProgress},
	}
	ok, reason := CanAdvance(&allTasks[0], allTasks)
	if ok {
		t.Fatal("CanAdvance(planned with unfinished dep) = true, want false")
	}
	if reason != "dependency \"T2\" is not done" {
		t.Errorf("reason = %q, want dependency warning", reason)
	}
}

func TestCanAdvance_buildAlwaysAllowed(t *testing.T) {
	task := data.Task{ID: "T1", Column: data.ColumnInProgress, Stage: data.StageBuild}
	ok, reason := CanAdvance(&task, nil)
	if !ok {
		t.Errorf("CanAdvance(build) = false %q, want true", reason)
	}
}

func TestCanAdvance_testAlwaysAllowed(t *testing.T) {
	task := data.Task{ID: "T1", Column: data.ColumnInProgress, Stage: data.StageTest}
	ok, reason := CanAdvance(&task, nil)
	if !ok {
		t.Errorf("CanAdvance(test) = false %q, want true", reason)
	}
}

func TestCanAdvance_auditDoneBlockedByDependency(t *testing.T) {
	allTasks := []data.Task{
		{ID: "T1", Column: data.ColumnInProgress, Stage: data.StageAudit, DependsOn: []string{"T2"}},
		{ID: "T2", Column: data.ColumnInProgress, Stage: data.StageBuild},
	}
	ok, reason := CanAdvance(&allTasks[0], allTasks)
	if ok {
		t.Fatal("CanAdvance(audit with undep) = true, want false")
	}
	if reason == "" {
		t.Fatal("expected non-empty reason string")
	}
}

func TestCanAdvance_auditDoneAllowedWhenDepsDone(t *testing.T) {
	allTasks := []data.Task{
		{ID: "T1", Column: data.ColumnInProgress, Stage: data.StageAudit, DependsOn: []string{"T2"}},
		{ID: "T2", Column: data.ColumnDone},
	}
	ok, reason := CanAdvance(&allTasks[0], allTasks)
	if !ok {
		t.Errorf("CanAdvance(audit with dep done) = false %q, want true", reason)
	}
}

func TestCanAdvance_doneBlocked(t *testing.T) {
	task := data.Task{ID: "T1", Column: data.ColumnDone}
	ok, reason := CanAdvance(&task, nil)
	if ok {
		t.Fatal("CanAdvance(done) = true, want false")
	}
	if reason == "" {
		t.Fatal("expected non-empty reason string")
	}
}

func TestCanAdvance_unknownStageBlocked(t *testing.T) {
	task := data.Task{ID: "T1", Column: data.ColumnInProgress, Stage: "invalid"}
	ok, reason := CanAdvance(&task, nil)
	if ok {
		t.Fatal("CanAdvance(unknown stage) = true, want false")
	}
	if reason == "" {
		t.Fatal("expected non-empty reason string")
	}
}

func TestCanAdvance_auditDepsNotFoundBlocked(t *testing.T) {
	allTasks := []data.Task{
		{ID: "T1", Column: data.ColumnInProgress, Stage: data.StageAudit, DependsOn: []string{"T2"}},
	}
	ok, reason := CanAdvance(&allTasks[0], allTasks)
	if ok {
		t.Fatal("CanAdvance(audit missing dep) = true, want false")
	}
	if reason == "" {
		t.Fatal("expected non-empty reason string")
	}
}
