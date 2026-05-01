package data

import (
	"testing"
)

func TestTaskString(t *testing.T) {
	task := Task{
		ID:      "E01/T001",
		Title:   "Test Task",
		Column:  ColumnPlanned,
	}
	want := "Task(E01/T001)"
	if got := task.String(); got != want {
		t.Errorf("Task.String() = %v, want %v", got, want)
	}
}

func TestColumnTypes(t *testing.T) {
	tests := []struct {
		input ColumnType
		want  string
	}{
		{ColumnPlanned, "planned"},
		{ColumnInProgress, "in_progress"},
		{ColumnDone, "done"},
	}

	for _, tt := range tests {
		if string(tt.input) != tt.want {
			t.Errorf("ColumnType = %v, want %v", tt.input, tt.want)
		}
	}
}

func TestProgressStage(t *testing.T) {
	tests := []struct {
		input ProgressStage
		want  string
	}{
		{StageBuild, "build"},
		{StageTest, "test"},
		{StageAudit, "audit"},
	}

	for _, tt := range tests {
		if string(tt.input) != tt.want {
			t.Errorf("ProgressStage = %v, want %v", tt.input, tt.want)
		}
	}
}