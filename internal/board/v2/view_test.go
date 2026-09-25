package v2

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// openBoard runs the model through the load its Init would have run, so a view
// test asserts against the same state the program renders.
func openBoard(t *testing.T, root, objectiveFilter string) Model {
	t.Helper()
	model := NewModel(Options{Root: root, ObjectiveFilter: objectiveFilter})
	model.Width = 100
	model.Height = 30

	updated, _ := model.Update(loadCmd(root)().(projectLoadedMsg))
	next, ok := updated.(Model)
	if !ok {
		t.Fatalf("Update() returned %T, want Model", updated)
	}
	return next
}

func TestViewEmptyTemplateProjectOpensWithThreeEmptyColumns(t *testing.T) {
	root := writeEmptyProjectFromTemplate(t)

	got := openBoard(t, root, "").View()

	for _, label := range []string{"PLANNED (0)", "IN PROGRESS (0)", "DONE (0)"} {
		if !strings.Contains(got, label) {
			t.Errorf("view missing empty column %q:\n%s", label, got)
		}
	}
	if !strings.Contains(got, "savepoint doctor") || strings.Contains(got, "Choose a Goal") {
		t.Errorf("zero-Goal view does not point to savepoint doctor:\n%s", got)
	}
	for _, forbidden := range []string{diagnosticHeading, "error", "not found", "MIGRATION"} {
		if strings.Contains(got, forbidden) {
			t.Errorf("fresh project view contains %q, which reads as broken:\n%s", forbidden, got)
		}
	}
}

func TestViewReportsTheCountsItLoaded(t *testing.T) {
	root := writeValidProject(t)

	model := openBoard(t, root, "")
	got := model.View()

	for _, want := range []string{"◆ 0/1 Objectives", "▣ 1/2 Tasks", "✗ 0/0 Issues"} {
		if !strings.Contains(got, want) {
			t.Errorf("view does not report the loaded counts as %q:\n%s", want, got)
		}
	}
	if !strings.Contains(got, "PLANNED (1)") || !strings.Contains(got, "DONE (1)") {
		t.Errorf("view does not report Tasks in the columns their status names:\n%s", got)
	}
	if !strings.Contains(got, "T-001") || !strings.Contains(got, "Do the thing") {
		t.Errorf("view does not name the record the projection selected:\n%s", got)
	}
}

func TestViewDiagnosticScreenDrawsNoColumns(t *testing.T) {
	root := writeValidProject(t)
	// A Task with no title: DecodeTaskV2 refuses it, so no project loads.
	writeTitlelessTask(t, root)

	got := openBoard(t, root, "").View()

	if !strings.Contains(got, diagnosticHeading) {
		t.Errorf("view missing the load-diagnostic heading:\n%s", got)
	}
	if !strings.Contains(got, "T-001-fixture.md") || !strings.Contains(got, "missing required field title") {
		t.Errorf("diagnostic view does not name the file and the problem:\n%s", got)
	}
	for _, column := range []string{"PLANNED", "IN PROGRESS", "DONE"} {
		if strings.Contains(got, column) {
			t.Errorf("diagnostic view drew column %q; it must draw no board:\n%s", column, got)
		}
	}
}

func TestViewBeforeFirstLoadDrawsNoBoard(t *testing.T) {
	model := NewModel(Options{Root: writeValidProject(t)})

	got := model.View()

	if strings.Contains(got, "PLANNED") {
		t.Errorf("view drew columns before the first load completed:\n%s", got)
	}
}

func TestUpdateWindowSizeAndQuit(t *testing.T) {
	model := NewModel(Options{Root: writeValidProject(t)})

	sized, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	resized := sized.(Model)
	if resized.Width != 120 || resized.Height != 40 {
		t.Errorf("size = %dx%d, want 120x40", resized.Width, resized.Height)
	}

	if _, cmd := resized.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}); cmd == nil {
		t.Error("q returned no command, want quit")
	}
}
