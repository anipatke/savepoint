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
			if err := Advance(&task); err != nil {
				t.Fatalf("Advance() error = %v", err)
			}
			if task.Column != tt.expectCol || task.Stage != tt.expectSt {
				t.Errorf("Advance() = %v/%v, want %v/%v", task.Column, task.Stage, tt.expectCol, tt.expectSt)
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
			task := data.Task{Column: tt.initialCol, Stage: tt.initialSt}
			if err := Retreat(&task); err != nil {
				t.Fatalf("Retreat() error = %v", err)
			}
			if task.Column != tt.expectCol || task.Stage != tt.expectSt {
				t.Errorf("Retreat() = %v/%v, want %v/%v", task.Column, task.Stage, tt.expectCol, tt.expectSt)
			}
		})
	}
}

func TestAdvance_returnsLifecycleErrorForInvalidState(t *testing.T) {
	task := data.Task{ID: "T1", Column: data.ColumnInProgress, Stage: data.ProgressStage("review")}
	err := Advance(&task)
	if err == nil {
		t.Fatal("Advance() expected error")
	}
	if err.Error() != `unknown stage "review"` {
		t.Fatalf("Advance() error = %v, want unknown stage", err)
	}
}

func TestAdvance_clearsLegacyLoadedStaleStage(t *testing.T) {
	task := data.Task{ID: "T1", Column: data.ColumnPlanned, Stage: data.StageAudit}
	if err := Advance(&task); err != nil {
		t.Fatalf("Advance() error = %v", err)
	}
	if task.Column != data.ColumnInProgress || task.Stage != data.StageBuild {
		t.Fatalf("Advance() = %q/%q, want in_progress/build", task.Column, task.Stage)
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

func TestCanAdvance_shortTaskDependencyAllowedWithinSameEpic(t *testing.T) {
	allTasks := []data.Task{
		{ID: "E06-canvas-polish/T004-current", Release: "v1", Epic: "E06-canvas-polish", Column: data.ColumnPlanned, DependsOn: []string{"T003"}},
		{ID: "E06-canvas-polish/T003-prereq", Release: "v1", Epic: "E06-canvas-polish", Column: data.ColumnDone},
		{ID: "E07-other/T003-prereq", Release: "v1", Epic: "E07-other", Column: data.ColumnPlanned},
	}
	ok, reason := CanAdvance(&allTasks[0], allTasks)
	if !ok {
		t.Errorf("CanAdvance(short same-epic dep) = false %q, want true", reason)
	}
}

func TestCanAdvance_filenameStyleTaskDependencyAllowedWithinSameEpic(t *testing.T) {
	allTasks := []data.Task{
		{ID: "E17-defect-workflow-tui/T010-defect-resolve-hotkey", Release: "v1.2", Epic: "E17-defect-workflow-tui", Column: data.ColumnInProgress, Stage: data.StageAudit, DependsOn: []string{"T004-defects-overlay"}},
		{ID: "E17-defect-workflow-tui/T004-defects-overlay", Release: "v1.2", Epic: "E17-defect-workflow-tui", Column: data.ColumnDone},
		{ID: "E18-template-skill-optimisation/T004-defects-overlay", Release: "v1.2", Epic: "E18-template-skill-optimisation", Column: data.ColumnPlanned},
	}
	ok, reason := CanAdvance(&allTasks[0], allTasks)
	if !ok {
		t.Errorf("CanAdvance(filename-style same-epic dep) = false %q, want true", reason)
	}
}

func TestCanAdvance_shortTaskDependencyBlockedWhenNotDone(t *testing.T) {
	allTasks := []data.Task{
		{ID: "E06-canvas-polish/T004-current", Release: "v1", Epic: "E06-canvas-polish", Column: data.ColumnPlanned, DependsOn: []string{"T003"}},
		{ID: "E06-canvas-polish/T003-prereq", Release: "v1", Epic: "E06-canvas-polish", Column: data.ColumnInProgress},
	}
	ok, reason := CanAdvance(&allTasks[0], allTasks)
	if ok {
		t.Fatal("CanAdvance(short unfinished dep) = true, want false")
	}
	if reason != "dependency \"T003\" is not done" {
		t.Errorf("reason = %q, want dependency warning", reason)
	}
}

func TestCanAdvance_shortTaskDependencyMissingOutsideSameEpic(t *testing.T) {
	allTasks := []data.Task{
		{ID: "E06-canvas-polish/T004-current", Release: "v1", Epic: "E06-canvas-polish", Column: data.ColumnPlanned, DependsOn: []string{"T003"}},
		{ID: "E07-other/T003-prereq", Release: "v1", Epic: "E07-other", Column: data.ColumnDone},
	}
	ok, reason := CanAdvance(&allTasks[0], allTasks)
	if ok {
		t.Fatal("CanAdvance(short missing same-epic dep) = true, want false")
	}
	if reason != "dependency \"T003\" not found" {
		t.Errorf("reason = %q, want not found", reason)
	}
}

func TestCanAdvance_fullTaskDependencyStillWorks(t *testing.T) {
	allTasks := []data.Task{
		{ID: "E06-canvas-polish/T004-current", Release: "v1", Epic: "E06-canvas-polish", Column: data.ColumnPlanned, DependsOn: []string{"E06-canvas-polish/T003-prereq"}},
		{ID: "E06-canvas-polish/T003-prereq", Release: "v1", Epic: "E06-canvas-polish", Column: data.ColumnDone},
	}
	ok, reason := CanAdvance(&allTasks[0], allTasks)
	if !ok {
		t.Errorf("CanAdvance(full task dep) = false %q, want true", reason)
	}
}

func TestCanAdvance_shortEpicDependencyAllowedWhenAudited(t *testing.T) {
	task := data.Task{ID: "E06-canvas-polish/T004-current", Release: "v1", Epic: "E06-canvas-polish", Column: data.ColumnPlanned, DependsOn: []string{"E03"}}
	ok, reason := CanAdvance(&task, nil, map[string]string{"E03-canvas-baseline": "audited"})
	if !ok {
		t.Errorf("CanAdvance(short audited epic dep) = false %q, want true", reason)
	}
}

func TestCanAdvance_shortEpicDependencyBlockedWhenNotAudited(t *testing.T) {
	task := data.Task{ID: "E06-canvas-polish/T004-current", Release: "v1", Epic: "E06-canvas-polish", Column: data.ColumnPlanned, DependsOn: []string{"E03"}}
	ok, reason := CanAdvance(&task, nil, map[string]string{"E03-canvas-baseline": "planned"})
	if ok {
		t.Fatal("CanAdvance(short unaudited epic dep) = true, want false")
	}
	if reason != "dependency \"E03\" is not audited" {
		t.Errorf("reason = %q, want not audited", reason)
	}
}

func TestCanAdvance_shortEpicDependencyMissing(t *testing.T) {
	task := data.Task{ID: "E06-canvas-polish/T004-current", Release: "v1", Epic: "E06-canvas-polish", Column: data.ColumnPlanned, DependsOn: []string{"E03"}}
	ok, reason := CanAdvance(&task, nil, map[string]string{"E04-other": "audited"})
	if ok {
		t.Fatal("CanAdvance(short missing epic dep) = true, want false")
	}
	if reason != "dependency \"E03\" not found" {
		t.Errorf("reason = %q, want not found", reason)
	}
}

func TestCanAdvance_fullEpicDependencyStillWorks(t *testing.T) {
	task := data.Task{ID: "E06-canvas-polish/T004-current", Release: "v1", Epic: "E06-canvas-polish", Column: data.ColumnPlanned, DependsOn: []string{"E03-canvas-baseline"}}
	ok, reason := CanAdvance(&task, nil, map[string]string{"E03-canvas-baseline": "audited"})
	if !ok {
		t.Errorf("CanAdvance(full audited epic dep) = false %q, want true", reason)
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

func TestCanAdvance_auditDoneBlockedByUnauditedEpicDependency(t *testing.T) {
	task := data.Task{
		ID:        "E06-current/T004-current",
		Release:   "v1",
		Epic:      "E06-current",
		Column:    data.ColumnInProgress,
		Stage:     data.StageAudit,
		DependsOn: []string{"E03"},
	}
	ok, reason := CanAdvance(&task, nil, map[string]string{"E03-prereq": "planned"})
	if ok {
		t.Fatal("CanAdvance(audit with unaudited epic dep) = true, want false")
	}
	if reason != "dependency \"E03\" is not audited" {
		t.Errorf("reason = %q, want epic audit warning", reason)
	}
}

func TestCanAdvance_auditDoneAllowedWhenEpicDependencyAudited(t *testing.T) {
	task := data.Task{
		ID:        "E06-current/T004-current",
		Release:   "v1",
		Epic:      "E06-current",
		Column:    data.ColumnInProgress,
		Stage:     data.StageAudit,
		DependsOn: []string{"E03"},
	}
	ok, reason := CanAdvance(&task, nil, map[string]string{"E03-prereq": "audited"})
	if !ok {
		t.Errorf("CanAdvance(audit with audited epic dep) = false %q, want true", reason)
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
