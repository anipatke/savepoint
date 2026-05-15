package data

import (
	"strings"
	"testing"
)

func TestParseTaskLifecycle_normalizesLoadCompatibleMetadata(t *testing.T) {
	tests := []struct {
		name       string
		metadata   TaskLifecycleMetadata
		wantStatus ColumnType
		wantStage  ProgressStage
	}{
		{
			name:       "missing status loads as planned",
			metadata:   TaskLifecycleMetadata{},
			wantStatus: ColumnPlanned,
		},
		{
			name:       "legacy todo status loads as planned",
			metadata:   TaskLifecycleMetadata{Status: LegacyTaskStatusTodo},
			wantStatus: ColumnPlanned,
		},
		{
			name:       "legacy phase supplies in progress stage",
			metadata:   TaskLifecycleMetadata{Status: ColumnInProgress, Phase: StageTest},
			wantStatus: ColumnInProgress,
			wantStage:  StageTest,
		},
		{
			name:       "stale stage outside in progress is cleared",
			metadata:   TaskLifecycleMetadata{Status: ColumnPlanned, Stage: StageBuild},
			wantStatus: ColumnPlanned,
		},
		{
			name:       "legacy stale implementation stage outside in progress is cleared",
			metadata:   TaskLifecycleMetadata{Status: ColumnDone, Stage: ProgressStage("implementation")},
			wantStatus: ColumnDone,
		},
		{
			name:       "legacy stale implementation phase outside in progress is cleared",
			metadata:   TaskLifecycleMetadata{Status: ColumnDone, Phase: ProgressStage("implementation")},
			wantStatus: ColumnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTaskLifecycle(tt.metadata)
			if err != nil {
				t.Fatalf("ParseTaskLifecycle() error = %v", err)
			}
			if got.Status != tt.wantStatus {
				t.Errorf("Status = %q, want %q", got.Status, tt.wantStatus)
			}
			if got.Stage != tt.wantStage {
				t.Errorf("Stage = %q, want %q", got.Stage, tt.wantStage)
			}
		})
	}
}

func TestParseTaskLifecycle_rejectsMalformedInProgressMetadata(t *testing.T) {
	tests := []struct {
		name     string
		metadata TaskLifecycleMetadata
		want     string
	}{
		{
			name:     "missing in progress stage",
			metadata: TaskLifecycleMetadata{Status: ColumnInProgress},
			want:     "stage is required",
		},
		{
			name:     "invalid in progress stage",
			metadata: TaskLifecycleMetadata{Status: ColumnInProgress, Stage: ProgressStage("implementation")},
			want:     `invalid stage "implementation"`,
		},
		{
			name:     "unknown status",
			metadata: TaskLifecycleMetadata{Status: ColumnType("review")},
			want:     `invalid status "review"`,
		},
		{
			name:     "invalid legacy phase outside in progress",
			metadata: TaskLifecycleMetadata{Status: ColumnPlanned, Phase: ProgressStage("done")},
			want:     `invalid legacy phase "done"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseTaskLifecycle(tt.metadata)
			if err == nil {
				t.Fatal("ParseTaskLifecycle() expected error")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("ParseTaskLifecycle() error = %v, want %q", err, tt.want)
			}
		})
	}
}

func TestValidateTaskLifecycleStateForWrite_rejectsLoadCompatibilityMetadata(t *testing.T) {
	tests := []struct {
		name  string
		state TaskLifecycleState
		want  string
	}{
		{
			name: "missing status",
			want: `invalid status ""`,
		},
		{
			name:  "legacy todo status",
			state: TaskLifecycleState{Status: LegacyTaskStatusTodo},
			want:  `invalid status "todo"`,
		},
		{
			name:  "stale non in progress stage",
			state: TaskLifecycleState{Status: ColumnDone, Stage: StageAudit},
			want:  `stage field "audit" is only valid`,
		},
		{
			name:  "invalid in progress stage",
			state: TaskLifecycleState{Status: ColumnInProgress, Stage: ProgressStage("implementation")},
			want:  `invalid stage "implementation"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTaskLifecycleStateForWrite(tt.state)
			if err == nil {
				t.Fatal("ValidateTaskLifecycleStateForWrite() expected error")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("ValidateTaskLifecycleStateForWrite() error = %v, want %q", err, tt.want)
			}
		})
	}
}

func TestValidateTaskLifecycleTransition_validatesTargetLifecycle(t *testing.T) {
	err := ValidateTaskLifecycleTransition(TaskLifecycleTransition{
		From: TaskLifecycleState{Status: ColumnPlanned},
		To:   TaskLifecycleState{Status: ColumnInProgress, Stage: StageBuild},
	})
	if err != nil {
		t.Fatalf("ValidateTaskLifecycleTransition() error = %v", err)
	}

	err = ValidateTaskLifecycleTransition(TaskLifecycleTransition{
		From: TaskLifecycleState{Status: ColumnInProgress, Stage: StageBuild},
		To:   TaskLifecycleState{Status: ColumnDone, Stage: StageBuild},
	})
	if err == nil {
		t.Fatal("ValidateTaskLifecycleTransition() expected invalid target stage error")
	}
	if !strings.Contains(err.Error(), `stage field "build" is only valid`) {
		t.Fatalf("ValidateTaskLifecycleTransition() error = %v, want target stage message", err)
	}

	err = ValidateTaskLifecycleTransition(TaskLifecycleTransition{
		From: TaskLifecycleState{Status: ColumnInProgress},
		To:   TaskLifecycleState{Status: ColumnDone},
	})
	if err == nil {
		t.Fatal("ValidateTaskLifecycleTransition() expected invalid source lifecycle error")
	}
	if !strings.Contains(err.Error(), "invalid source lifecycle") {
		t.Fatalf("ValidateTaskLifecycleTransition() error = %v, want source lifecycle message", err)
	}
}

func TestAdvanceTaskLifecycle_movesThroughCanonicalStates(t *testing.T) {
	tests := []struct {
		name       string
		task       Task
		wantStatus ColumnType
		wantStage  ProgressStage
	}{
		{"planned to build", Task{Column: ColumnPlanned}, ColumnInProgress, StageBuild},
		{"build to test", Task{Column: ColumnInProgress, Stage: StageBuild}, ColumnInProgress, StageTest},
		{"test to audit", Task{Column: ColumnInProgress, Stage: StageTest}, ColumnInProgress, StageAudit},
		{"audit to done", Task{Column: ColumnInProgress, Stage: StageAudit}, ColumnDone, ""},
		{"legacy todo to build", Task{Column: LegacyTaskStatusTodo}, ColumnInProgress, StageBuild},
		{"stale planned stage to build", Task{Column: ColumnPlanned, Stage: StageAudit}, ColumnInProgress, StageBuild},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := tt.task
			transition, err := AdvanceTaskLifecycle(&task)
			if err != nil {
				t.Fatalf("AdvanceTaskLifecycle() error = %v", err)
			}
			if task.Column != tt.wantStatus || task.Stage != tt.wantStage {
				t.Fatalf("task lifecycle = %q/%q, want %q/%q", task.Column, task.Stage, tt.wantStatus, tt.wantStage)
			}
			if transition.To.Status != tt.wantStatus || transition.To.Stage != tt.wantStage {
				t.Fatalf("transition.To = %+v, want %q/%q", transition.To, tt.wantStatus, tt.wantStage)
			}
		})
	}
}

func TestAdvanceTaskLifecycle_rejectsInvalidStates(t *testing.T) {
	tests := []struct {
		name string
		task Task
		want string
	}{
		{"done", Task{Column: ColumnDone}, "already done"},
		{"unknown status", Task{Column: ColumnType("review")}, `unknown status "review"`},
		{"unknown stage", Task{Column: ColumnInProgress, Stage: ProgressStage("review")}, `unknown stage "review"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := AdvanceTaskLifecycle(&tt.task)
			if err == nil {
				t.Fatal("AdvanceTaskLifecycle() expected error")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("AdvanceTaskLifecycle() error = %v, want %q", err, tt.want)
			}
		})
	}
}

func TestRetreatTaskLifecycle_movesThroughCanonicalStates(t *testing.T) {
	tests := []struct {
		name       string
		task       Task
		wantStatus ColumnType
		wantStage  ProgressStage
	}{
		{"done to audit", Task{Column: ColumnDone}, ColumnInProgress, StageAudit},
		{"audit to test", Task{Column: ColumnInProgress, Stage: StageAudit}, ColumnInProgress, StageTest},
		{"test to build", Task{Column: ColumnInProgress, Stage: StageTest}, ColumnInProgress, StageBuild},
		{"build to planned", Task{Column: ColumnInProgress, Stage: StageBuild}, ColumnPlanned, ""},
		{"planned stays planned", Task{Column: ColumnPlanned}, ColumnPlanned, ""},
		{"stale done stage to audit", Task{Column: ColumnDone, Stage: StageBuild}, ColumnInProgress, StageAudit},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := tt.task
			transition, err := RetreatTaskLifecycle(&task)
			if err != nil {
				t.Fatalf("RetreatTaskLifecycle() error = %v", err)
			}
			if task.Column != tt.wantStatus || task.Stage != tt.wantStage {
				t.Fatalf("task lifecycle = %q/%q, want %q/%q", task.Column, task.Stage, tt.wantStatus, tt.wantStage)
			}
			if transition.To.Status != tt.wantStatus || transition.To.Stage != tt.wantStage {
				t.Fatalf("transition.To = %+v, want %q/%q", transition.To, tt.wantStatus, tt.wantStage)
			}
		})
	}
}

func TestSameTaskLifecycleForTransition_normalizesLoadCompatibleState(t *testing.T) {
	if !SameTaskLifecycleForTransition(
		Task{Column: ColumnPlanned, Stage: StageBuild},
		Task{Column: ColumnPlanned},
	) {
		t.Fatal("SameTaskLifecycleForTransition() should ignore stale stage outside in_progress")
	}

	if !SameTaskLifecycleForTransition(
		Task{Column: ColumnInProgress},
		Task{Column: ColumnInProgress, Stage: StageBuild},
	) {
		t.Fatal("SameTaskLifecycleForTransition() should compare missing in-progress stage as build")
	}

	if SameTaskLifecycleForTransition(
		Task{Column: ColumnInProgress, Stage: StageBuild},
		Task{Column: ColumnInProgress, Stage: StageTest},
	) {
		t.Fatal("SameTaskLifecycleForTransition() should detect real stage changes")
	}
}

func TestDiagnoseTaskLifecycle_reportsDoctorLifecycleProblems(t *testing.T) {
	tests := []struct {
		name  string
		input TaskLifecycleDiagnosticInput
		want  []TaskLifecycleDiagnosticCode
	}{
		{
			name: "missing status",
			input: TaskLifecycleDiagnosticInput{
				Metadata: TaskLifecycleMetadata{},
			},
			want: []TaskLifecycleDiagnosticCode{TaskLifecycleMissingStatus},
		},
		{
			name: "invalid status",
			input: TaskLifecycleDiagnosticInput{
				Metadata:  TaskLifecycleMetadata{Status: ColumnType("review")},
				HasStatus: true,
			},
			want: []TaskLifecycleDiagnosticCode{TaskLifecycleInvalidStatus},
		},
		{
			name: "missing in progress stage",
			input: TaskLifecycleDiagnosticInput{
				Metadata:  TaskLifecycleMetadata{Status: ColumnInProgress},
				HasStatus: true,
			},
			want: []TaskLifecycleDiagnosticCode{TaskLifecycleMissingStage},
		},
		{
			name: "invalid in progress stage",
			input: TaskLifecycleDiagnosticInput{
				Metadata:  TaskLifecycleMetadata{Status: ColumnInProgress, Stage: ProgressStage("done")},
				HasStatus: true,
				HasStage:  true,
			},
			want: []TaskLifecycleDiagnosticCode{TaskLifecycleInvalidStage},
		},
		{
			name: "legacy phase",
			input: TaskLifecycleDiagnosticInput{
				Metadata:  TaskLifecycleMetadata{Status: ColumnInProgress, Phase: StageBuild},
				HasStatus: true,
				HasPhase:  true,
			},
			want: []TaskLifecycleDiagnosticCode{TaskLifecycleLegacyPhase},
		},
		{
			name: "stale stage outside in progress",
			input: TaskLifecycleDiagnosticInput{
				Metadata:  TaskLifecycleMetadata{Status: ColumnPlanned, Stage: StageBuild},
				HasStatus: true,
				HasStage:  true,
			},
			want: []TaskLifecycleDiagnosticCode{TaskLifecycleStaleStage},
		},
		{
			name: "legacy phase and stale implementation stage outside in progress",
			input: TaskLifecycleDiagnosticInput{
				Metadata:  TaskLifecycleMetadata{Status: ColumnDone, Stage: LegacyTaskStageImplementation, Phase: LegacyTaskStageImplementation},
				HasStatus: true,
				HasStage:  true,
				HasPhase:  true,
			},
			want: []TaskLifecycleDiagnosticCode{TaskLifecycleLegacyPhase, TaskLifecycleStaleStage},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DiagnoseTaskLifecycle(tt.input)
			if len(got) != len(tt.want) {
				t.Fatalf("DiagnoseTaskLifecycle() = %+v, want codes %v", got, tt.want)
			}
			for i, want := range tt.want {
				if got[i].Code != want {
					t.Fatalf("DiagnoseTaskLifecycle()[%d].Code = %q, want %q", i, got[i].Code, want)
				}
				if got[i].Message == "" {
					t.Fatalf("DiagnoseTaskLifecycle()[%d].Message is empty", i)
				}
			}
		})
	}
}

func TestTaskLifecycleContract_exposesCanonicalValuesAndAliases(t *testing.T) {
	statuses := CanonicalTaskStatuses()
	if len(statuses) != 3 || statuses[0] != ColumnPlanned || statuses[1] != ColumnInProgress || statuses[2] != ColumnDone {
		t.Fatalf("CanonicalTaskStatuses() = %v, want planned/in_progress/done", statuses)
	}

	stages := CanonicalTaskStages()
	if len(stages) != 3 || stages[0] != StageBuild || stages[1] != StageTest || stages[2] != StageAudit {
		t.Fatalf("CanonicalTaskStages() = %v, want build/test/audit", stages)
	}

	status, ok := ResolveTaskStatusAlias(LegacyTaskStatusTodo)
	if !ok || status != ColumnPlanned {
		t.Fatalf("ResolveTaskStatusAlias(todo) = %q, %v; want planned, true", status, ok)
	}

	if !IsLegacyTaskStageAlias(LegacyTaskStageImplementation) {
		t.Fatal("IsLegacyTaskStageAlias(implementation) = false, want true")
	}
}

func TestValidateTaskLifecycle_allowsPlannedWithoutStage(t *testing.T) {
	task := Task{Column: ColumnPlanned}
	if err := ValidateTaskLifecycle(&task); err != nil {
		t.Fatalf("ValidateTaskLifecycle() error = %v", err)
	}
}

func TestValidateTaskLifecycle_rejectsInProgressWithoutStage(t *testing.T) {
	task := Task{Column: ColumnInProgress}
	if err := ValidateTaskLifecycle(&task); err == nil {
		t.Fatal("ValidateTaskLifecycle() expected missing stage error")
	}
}

func TestValidateTaskLifecycle_allowsInProgressWithStage(t *testing.T) {
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

func TestValidateTaskLifecycle_rejectsStageOutsideInProgress(t *testing.T) {
	task := Task{Column: ColumnPlanned, Stage: StageBuild}
	if err := ValidateTaskLifecycle(&task); err == nil {
		t.Fatal("ValidateTaskLifecycle() expected stage/status error")
	}
}
