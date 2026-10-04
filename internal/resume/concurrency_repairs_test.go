package resume

import (
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/data"
)

func TestParallelLinesShowPlanningDiagnosticsEvenWithNoOtherAdvice(t *testing.T) {
	c := &data.ConcurrencyV2{Enabled: true, Diagnostics: []data.PlanDiagnosticV2{
		{Record: "T-004", Path: "tasks/T-004.md", Field: "planned_writes", Message: "glob \"bad/*.go\" is not allowed\x1b[2J"},
		{Record: "O-001", Path: "O-001.md", Field: "lanes", Message: "lane key \"core\" is declared more than once"},
	}}
	objective := strings.Join(ParallelLines(&data.V2Index{}, c, ""), "\n")
	for _, want := range []string{"T-004 tasks/T-004.md, field planned_writes", "bad/*.go", "O-001 O-001.md, field lanes"} {
		if !strings.Contains(objective, want) {
			t.Errorf("objective advice lacks %q:\n%s", want, objective)
		}
	}
	task := strings.Join(ParallelLines(&data.V2Index{}, c, "T-004"), "\n")
	if !strings.Contains(task, "planned_writes") || strings.Contains(task, "field lanes") {
		t.Errorf("task advice = %q, want only the Task's own diagnostic", task)
	}
	if strings.ContainsAny(objective, "\x1b") {
		t.Errorf("diagnostic keeps a terminal control: %q", objective)
	}
}

func TestInstructionStartsWithAStandaloneNextLine(t *testing.T) {
	task := &data.TaskV2{ID: "T-004", Title: "Core work", Status: data.ColumnPlanned}
	objective := &data.ObjectiveV2{ID: "O-001", Title: "O"}
	lines := instructionLines(&data.V2Index{
		Tasks:      map[string]*data.TaskV2{"T-004": task},
		Objectives: map[string]*data.ObjectiveV2{"O-001": objective},
	}, "O-001", data.ConcurrencyCandidateV2{Task: "T-004", Lane: "core"})
	want := NextLine(data.Next{Task: task, Objective: objective})
	if !strings.HasPrefix(want, "Start T-004") || !contains(lines, want) {
		t.Errorf("instruction lacks standalone %q:\n%s", want, strings.Join(lines, "\n"))
	}
}

func contains(lines []string, want string) bool {
	for _, line := range lines {
		if line == want {
			return true
		}
	}
	return false
}
