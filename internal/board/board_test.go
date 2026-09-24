package board

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/testutil"
)

// Every test here builds its project in a temporary directory and drives the
// dispatch through runWithFilters, so neither the live Savepoint project nor
// the process working directory is involved (TEST-04, ARCH-03).

func writeV1ProjectForDispatch(t *testing.T) string {
	t.Helper()
	projectRoot := t.TempDir()
	savepointRoot := filepath.Join(projectRoot, ".savepoint")
	testutil.SetupMinimalProject(t, savepointRoot, "v1", "E01-alpha")
	testutil.WriteTask(t, savepointRoot, "v1", "E01-alpha", testutil.TaskFixture{
		Slug:      "T001-first",
		Release:   "v1",
		Status:    string(data.ColumnPlanned),
		Objective: "Test task",
	})
	return projectRoot
}

func writeV2ProjectForDispatch(t *testing.T) string {
	t.Helper()
	projectRoot := t.TempDir()
	savepointRoot := filepath.Join(projectRoot, ".savepoint")
	testutil.WriteFile(t, filepath.Join(savepointRoot, "config.yml"), "schema_version: 2\n")
	testutil.WriteFile(t, filepath.Join(savepointRoot, "router.md"),
		"# Router\n\n## Current state\n\n```yaml\nstate: task\nobjective: O-001\ntask: T-001\nnext_action: \"go\"\n```\n")
	testutil.WriteFile(t, filepath.Join(savepointRoot, "objectives", "O-001-alpha", "Objective.md"),
		"---\nid: O-001\ntitle: \"First objective\"\nstatus: planned\n---\n\n# First objective\n")
	testutil.WriteFile(t, filepath.Join(savepointRoot, "objectives", "O-001-alpha", "tasks", "T-001-first.md"),
		"---\nid: T-001\ntitle: \"Do the thing\"\nobjective: O-001\nplanned_by: {role: planner, session: dispatch-fixture}\nstatus: planned\n---\n\n# Do the thing\n")
	return projectRoot
}

func TestRunWithFiltersRefusesV1ProjectBeforeStartingBoard(t *testing.T) {
	projectRoot := writeV1ProjectForDispatch(t)
	var stdout bytes.Buffer

	err := runWithFilters(projectRoot, Filters{}, &stdout, false)
	if err == nil {
		t.Fatal("runWithFilters() error = nil, want the V1 project routed to migration")
	}
	if !strings.Contains(err.Error(), "schema_version 1") || !strings.Contains(err.Error(), "migrate --dry-run") {
		t.Errorf("error = %q, want the named migration preview route", err.Error())
	}
	if stdout.Len() != 0 {
		t.Errorf("stdout = %q, want no V1 board rendered", stdout.String())
	}
}

func TestRunWithFiltersDispatchesV2ProjectToTheV2Board(t *testing.T) {
	projectRoot := writeV2ProjectForDispatch(t)
	var stdout bytes.Buffer

	if err := runWithFilters(projectRoot, Filters{}, &stdout, false); err != nil {
		t.Fatalf("runWithFilters() error = %v", err)
	}

	got := stdout.String()
	if !strings.Contains(got, "Start T-001") || !strings.Contains(got, "Objectives: 1  Tasks: 1") {
		t.Errorf("V2 project did not reach the V2 board:\n%s", got)
	}
	if strings.Contains(got, "releases directory not found") {
		t.Errorf("V2 project reached V1 discovery:\n%s", got)
	}
}

// TestRunWithFiltersV2LoadFailureReportsTheDiagnosticWithoutV1Fallback proves a
// project the V2 index refuses is reported by name: the fixture has no
// releases/ directory at all, so any retry through V1 discovery would surface
// as "releases directory not found" instead.
func TestRunWithFiltersV2LoadFailureReportsTheDiagnosticWithoutV1Fallback(t *testing.T) {
	projectRoot := writeV2ProjectForDispatch(t)
	testutil.WriteFile(t, filepath.Join(projectRoot, ".savepoint", "objectives", "O-001-alpha", "tasks", "T-001-first.md"),
		"---\nid: T-001\nobjective: O-001\nplanned_by: {role: planner, session: dispatch-fixture}\nstatus: planned\n---\n\n# Untitled\n")
	var stdout bytes.Buffer

	err := runWithFilters(projectRoot, Filters{}, &stdout, false)

	if err == nil {
		t.Fatal("runWithFilters() error = nil, want the load diagnostic")
	}
	if !strings.Contains(err.Error(), "T-001-first.md") || !strings.Contains(err.Error(), "missing required field title") {
		t.Errorf("error = %q, want the file and the problem named", err.Error())
	}
	if strings.Contains(err.Error(), "releases directory not found") {
		t.Errorf("error = %q, want no V1 discovery fallback", err.Error())
	}
	if stdout.Len() != 0 {
		t.Errorf("stdout = %q, want no partial board", stdout.String())
	}
}

func TestRunWithFiltersRejectsFiltersTheSchemaHasNoMeaningFor(t *testing.T) {
	tests := []struct {
		name      string
		project   func(t *testing.T) string
		filters   Filters
		wantParts []string
	}{
		{
			name:      "release against a V2 project",
			project:   writeV2ProjectForDispatch,
			filters:   Filters{Release: "v1"},
			wantParts: []string{"--release", "V2-only runtime"},
		},
		{
			name:      "epic against a V2 project",
			project:   writeV2ProjectForDispatch,
			filters:   Filters{Epic: "E01-alpha"},
			wantParts: []string{"--epic", "V2-only runtime"},
		},
		{
			name:      "objective against a V1 project",
			project:   writeV1ProjectForDispatch,
			filters:   Filters{Objective: "O001"},
			wantParts: []string{"schema_version 1", "migrate --dry-run"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			projectRoot := test.project(t)
			var stdout bytes.Buffer

			err := runWithFilters(projectRoot, test.filters, &stdout, false)

			if err == nil {
				t.Fatal("runWithFilters() error = nil, want the filter refused")
			}
			for _, part := range test.wantParts {
				if !strings.Contains(err.Error(), part) {
					t.Errorf("error = %q, want it to name %q", err.Error(), part)
				}
			}
			if stdout.Len() != 0 {
				t.Errorf("stdout = %q, want no board started", stdout.String())
			}
		})
	}
}

func TestRunWithFiltersRejectsUnknownObjectiveOnAV2Project(t *testing.T) {
	projectRoot := writeV2ProjectForDispatch(t)
	var stdout bytes.Buffer

	err := runWithFilters(projectRoot, Filters{Objective: "O404"}, &stdout, false)

	if err == nil {
		t.Fatal("runWithFilters() error = nil, want the unknown objective refused")
	}
	if !strings.Contains(err.Error(), "--objective O404") {
		t.Errorf("error = %q, want the flag and value named", err.Error())
	}
	if stdout.Len() != 0 {
		t.Errorf("stdout = %q, want no board written", stdout.String())
	}
}

func TestRunWithFiltersReportsAMissingProject(t *testing.T) {
	err := runWithFilters(t.TempDir(), Filters{}, &bytes.Buffer{}, false)

	if err == nil {
		t.Fatal("runWithFilters() error = nil, want a missing .savepoint reported")
	}
}

func TestRunWithFiltersRejectsV1FiltersWithoutStartingBoard(t *testing.T) {
	projectRoot := writeV1ProjectForDispatch(t)
	var stdout bytes.Buffer

	err := runWithFilters(projectRoot, Filters{Release: "v1", Epic: "E01-alpha"}, &stdout, false)
	if err == nil {
		t.Fatal("runWithFilters() error = nil, want V1 filters refused")
	}
	if !strings.Contains(err.Error(), "V2-only runtime") {
		t.Errorf("error = %q, want the V2-only filter diagnostic", err.Error())
	}
	if stdout.Len() != 0 {
		t.Errorf("stdout = %q, want no V1 board rendered", stdout.String())
	}
}
