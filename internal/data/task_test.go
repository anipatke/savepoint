package data

import (
	"testing"
)

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
