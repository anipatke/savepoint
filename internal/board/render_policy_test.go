package board

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/opencode/savepoint/internal/data"
)

func TestRenderPolicy_noBackgroundEscapes(t *testing.T) {
	lipgloss.SetColorProfile(termenv.ANSI256)

	task := data.Task{
		ID:      "E05-tasking-permissions/T004-implement-m-hotkey",
		Title:   "Implement router priority",
		Release: "v1.1",
		Epic:    "E05-tasking-permissions",
		Column:  data.ColumnInProgress,
		Stage:   data.StageBuild,
	}
	m := NewModel([]data.Task{
		{
			ID:      task.ID,
			Title:   task.Title,
			Release: task.Release,
			Epic:    task.Epic,
			Column:  task.Column,
			Stage:   task.Stage,
		},
	}, "v1.1", "E05-tasking-permissions")
	m.Width = 120
	m.Height = 30
	m.FocusedColumn = data.ColumnInProgress
	m.RouterState = &data.RouterState{
		State:      "task-building",
		Release:    "v1.1",
		Epic:       "E05-tasking-permissions",
		Task:       "E05-tasking-permissions/T004-implement-m-hotkey",
		NextAction: "Build E05-tasking-permissions/T004-implement-m-hotkey.",
	}

	cases := map[string]string{
		"board":            m.View(),
		"card":             RenderCard(task, 30, true, m.RouterState),
		"detail":           RenderDetail(task, 60, m.RouterState, 0, 0),
		"epic dropdown":    RenderEpicDropdown([]string{"E05-tasking-permissions"}, 0, 40),
		"release dropdown": RenderReleaseDropdown([]string{"v1.1"}, 0, 40),
		"help":             RenderHelp(60),
	}
	for name, got := range cases {
		assertNoBackgroundEscapes(t, name, got)
	}
}

func assertNoBackgroundEscapes(t *testing.T, name, got string) {
	t.Helper()
	for _, escape := range []string{"\x1b[48;", "\x1b[40m"} {
		if strings.Contains(got, escape) {
			t.Fatalf("%s emitted background escape prefix %q in %q", name, escape, got)
		}
	}
}

func TestRenderPolicy_usesSingleLineBorders(t *testing.T) {
	m := NewModel([]data.Task{{ID: "T001", Title: "Task", Column: data.ColumnPlanned}}, "v1", "E01")
	m.Width = 120
	got := m.View()

	if strings.Contains(got, "╭") || strings.Contains(got, "╮") || strings.Contains(got, "╰") || strings.Contains(got, "╯") {
		t.Fatalf("View should use single-line borders, got rounded border glyphs in %q", got)
	}
	if !strings.Contains(got, "┌") || !strings.Contains(got, "┐") || !strings.Contains(got, "└") || !strings.Contains(got, "┘") {
		t.Fatalf("View missing expected single-line border glyphs in %q", got)
	}
}
