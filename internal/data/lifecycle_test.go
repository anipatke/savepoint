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
			name:       "agent complete status loads as done",
			metadata:   TaskLifecycleMetadata{Status: LegacyTaskStatusComplete},
			wantStatus: ColumnDone,
		},
		{
			name:       "agent completed status loads as done",
			metadata:   TaskLifecycleMetadata{Status: LegacyTaskStatusCompleted},
			wantStatus: ColumnDone,
		},
		{
			name:       "legacy phase supplies in progress stage",
			metadata:   TaskLifecycleMetadata{Status: ColumnInProgress, Phase: StageTest},
			wantStatus: ColumnInProgress,
			wantStage:  StageTest,
		},
		{
			name:       "legacy implementation stage loads as build",
			metadata:   TaskLifecycleMetadata{Status: ColumnInProgress, Stage: LegacyTaskStageImplementation},
			wantStatus: ColumnInProgress,
			wantStage:  StageBuild,
		},
		{
			name:       "legacy implementation phase loads as build",
			metadata:   TaskLifecycleMetadata{Status: ColumnInProgress, Phase: LegacyTaskStageImplementation},
			wantStatus: ColumnInProgress,
			wantStage:  StageBuild,
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
			got := ParseTaskLifecycle(tt.metadata)
			if got.Status != tt.wantStatus {
				t.Errorf("Status = %q, want %q", got.Status, tt.wantStatus)
			}
			if got.Stage != tt.wantStage {
				t.Errorf("Stage = %q, want %q", got.Stage, tt.wantStage)
			}
		})
	}
}

func TestParseTaskLifecycle_healsMalformedMetadata(t *testing.T) {
	tests := []struct {
		name       string
		metadata   TaskLifecycleMetadata
		wantStatus ColumnType
		wantStage  ProgressStage
	}{
		{
			name:       "missing in progress stage defaults to build",
			metadata:   TaskLifecycleMetadata{Status: ColumnInProgress},
			wantStatus: ColumnInProgress,
			wantStage:  StageBuild,
		},
		{
			name:       "invalid in progress stage heals to build",
			metadata:   TaskLifecycleMetadata{Status: ColumnInProgress, Stage: ProgressStage("review")},
			wantStatus: ColumnInProgress,
			wantStage:  StageBuild,
		},
		{
			name:       "unknown status heals to planned",
			metadata:   TaskLifecycleMetadata{Status: ColumnType("review")},
			wantStatus: ColumnPlanned,
		},
		{
			name:       "invalid legacy phase outside in progress is dropped",
			metadata:   TaskLifecycleMetadata{Status: ColumnPlanned, Phase: ProgressStage("done")},
			wantStatus: ColumnPlanned,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseTaskLifecycle(tt.metadata)
			if got.Status != tt.wantStatus {
				t.Errorf("Status = %q, want %q", got.Status, tt.wantStatus)
			}
			if got.Stage != tt.wantStage {
				t.Errorf("Stage = %q, want %q", got.Stage, tt.wantStage)
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
			name:  "agent complete status",
			state: TaskLifecycleState{Status: LegacyTaskStatusComplete},
			want:  `invalid status "complete"`,
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

func TestTaskLifecycleContract_exposesCanonicalValuesAndAliases(t *testing.T) {
	statuses := CanonicalTaskStatuses()
	if len(statuses) != 3 || statuses[0] != ColumnPlanned || statuses[1] != ColumnInProgress || statuses[2] != ColumnDone {
		t.Fatalf("CanonicalTaskStatuses() = %v, want planned/in_progress/done", statuses)
	}

	status, ok := ResolveTaskStatusAlias(LegacyTaskStatusTodo)
	if !ok || status != ColumnPlanned {
		t.Fatalf("ResolveTaskStatusAlias(todo) = %q, %v; want planned, true", status, ok)
	}

	status, ok = ResolveTaskStatusAlias(LegacyTaskStatusComplete)
	if !ok || status != ColumnDone {
		t.Fatalf("ResolveTaskStatusAlias(complete) = %q, %v; want done, true", status, ok)
	}

	status, ok = ResolveTaskStatusAlias(LegacyTaskStatusCompleted)
	if !ok || status != ColumnDone {
		t.Fatalf("ResolveTaskStatusAlias(completed) = %q, %v; want done, true", status, ok)
	}

	if !IsLegacyTaskStageAlias(LegacyTaskStageImplementation) {
		t.Fatal("IsLegacyTaskStageAlias(implementation) = false, want true")
	}

	if NormalizeTaskStageForLoad(LegacyTaskStageImplementation) != StageBuild {
		t.Fatal("NormalizeTaskStageForLoad(implementation) should return build")
	}
}

func TestResolveEpicStatusAlias_reportsOnlyKnownLeaks(t *testing.T) {
	if alias, ok := ResolveEpicStatusAlias("epic-design"); !ok || alias != EpicStatus(ColumnPlanned) {
		t.Fatalf("ResolveEpicStatusAlias(epic-design) = %q, %v; want planned, true", alias, ok)
	}
	if alias, ok := ResolveEpicStatusAlias("completed"); !ok || alias != EpicStatus(ColumnDone) {
		t.Fatalf("ResolveEpicStatusAlias(completed) = %q, %v; want done, true", alias, ok)
	}
	if alias, ok := ResolveEpicStatusAlias(EpicStatus(ColumnPlanned)); ok {
		t.Fatalf("ResolveEpicStatusAlias(planned) = %q, %v; want \"\", false", alias, ok)
	}
	if alias, ok := ResolveEpicStatusAlias("garbage"); ok {
		t.Fatalf("ResolveEpicStatusAlias(garbage) = %q, %v; want \"\", false", alias, ok)
	}
}
