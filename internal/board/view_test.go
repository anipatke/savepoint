package board

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/opencode/savepoint/internal/data"
)

func TestView_rendersWithoutPanic(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Width = 120
	m.Height = 40
	got := m.View()
	if got == "" {
		t.Error("View() returned empty string")
	}
}

func TestView_containsColumnTitles(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Width = 120
	got := m.View()

	for _, title := range []string{"PLANNED", "IN PROGRESS", "DONE"} {
		if !strings.Contains(got, title) {
			t.Errorf("View() missing column title %q", title)
		}
	}
}

func TestView_containsHeader(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Width = 120
	got := m.View()
	if !strings.Contains(got, "S A V E P O I N T") {
		t.Error("View() missing spaced header text")
	}
	if !strings.Contains(got, "▣") {
		t.Error("View() missing header icon")
	}
}

func TestView_containsDivider(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Width = 120
	got := m.View()
	if !strings.Contains(got, "─") {
		t.Error("View() missing horizontal divider")
	}
}

func TestView_containsFooterPhases(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Width = 120
	got := m.View()
	for _, phase := range []string{"PLAN", "BUILD", "AUDIT"} {
		if !strings.Contains(got, phase) {
			t.Errorf("View() missing footer phase %q", phase)
		}
	}
}

func TestView_containsFooterHints(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	footer := m.renderFooter(80)

	if !strings.Contains(footer, "←/→:nav  E:epic  R:release  ?:help  q:quit") {
		t.Fatal("renderFooter() missing navigation hints")
	}

	lines := strings.Split(footer, "\n")
	if len(lines) != 3 {
		t.Fatalf("renderFooter() returned %d lines, want 3", len(lines))
	}
	if strings.TrimSpace(lines[1]) != "" {
		t.Fatalf("renderFooter() spacer line = %q, want blank", lines[1])
	}
	for i, line := range lines {
		if got := lipgloss.Width(line); got > 80 {
			t.Fatalf("renderFooter() line %d width = %d, want <= 80", i, got)
		}
	}
}

func TestView_containsBottomDivider(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Width = 120
	got := m.View()
	// There should be at least two divider lines (top and bottom)
	count := strings.Count(got, "─")
	if count < 2 {
		t.Errorf("View() expected at least 2 divider chars, got %d", count)
	}
}

func TestView_defaultWidthWhenZero(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	// Width=0: should use default and not panic
	got := m.View()
	if got == "" {
		t.Error("View() with zero Width returned empty string")
	}
}

func TestView_taskLabelFallback(t *testing.T) {
	tasks := []data.Task{{ID: "T1", Title: "", Column: data.ColumnPlanned}}
	m := NewModel(tasks, "v1", "E03")
	m.Width = 120
	got := m.View()
	if !strings.Contains(got, "T1") {
		t.Error("View() missing task ID when title is empty")
	}
}

func TestView_taskLabelWithTitle(t *testing.T) {
	tasks := []data.Task{{ID: "T2", Title: "My Task", Column: data.ColumnPlanned}}
	m := NewModel(tasks, "v1", "E03")
	m.Width = 120
	got := m.View()
	if !strings.Contains(got, "My Task") {
		t.Error("View() missing task title")
	}
}

func TestView_wideShowsEpicPanel(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Width = 120
	got := m.View()
	if !strings.Contains(got, "E03") {
		t.Error("View() at width>=120 missing epic panel content")
	}
}

func TestView_narrowShowsSingleColumn(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Width = 60
	m.FocusedColumn = data.ColumnInProgress
	got := m.View()
	if !strings.Contains(got, "IN PROGRESS") {
		t.Error("View() at width<80 missing focused column title")
	}
	if strings.Contains(got, "PLANNED") {
		t.Error("View() at width<80 should not show non-focused columns")
	}
}
