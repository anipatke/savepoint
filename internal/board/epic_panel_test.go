package board

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/opencode/savepoint/internal/data"
)

func TestRenderEpicSidebar_containsEpicsHeader(t *testing.T) {
	got := RenderEpicSidebar([]string{"E01", "E02"}, "E01", 28)
	if !strings.Contains(got, "EPICS") {
		t.Error("RenderEpicSidebar missing EPICS header")
	}
}

func TestRenderEpicSidebar_activeEpicMarked(t *testing.T) {
	got := RenderEpicSidebar([]string{"E01", "E02"}, "E01", 28)
	if !strings.Contains(got, epicActiveMarker) {
		t.Errorf("RenderEpicSidebar missing active marker %q", epicActiveMarker)
	}
}

func TestRenderEpicSidebar_allEpicsPresent(t *testing.T) {
	epics := []string{"E01-foo", "E02-bar", "E03-baz"}
	got := RenderEpicSidebar(epics, "E01-foo", 32)
	for _, e := range epics {
		if !strings.Contains(got, e) {
			t.Errorf("RenderEpicSidebar missing epic %q", e)
		}
	}
}

func TestRenderEpicSidebar_emptyEpicsFallback(t *testing.T) {
	got := RenderEpicSidebar(nil, "E03", 28)
	if !strings.Contains(got, "E03") {
		t.Error("RenderEpicSidebar with empty list should show selected epic")
	}
}

func TestRenderEpicSidebar_emptyBothShowsNone(t *testing.T) {
	got := RenderEpicSidebar(nil, "", 28)
	if !strings.Contains(got, "(none)") {
		t.Error("RenderEpicSidebar with no epics and no selected should show (none)")
	}
}

func TestRenderEpicDropdown_containsHeader(t *testing.T) {
	got := RenderEpicDropdown([]string{"E01", "E02"}, 0, 32)
	if !strings.Contains(got, "SELECT EPIC") {
		t.Error("RenderEpicDropdown missing SELECT EPIC header")
	}
}

func TestRenderEpicDropdown_cursorMarked(t *testing.T) {
	got := RenderEpicDropdown([]string{"E01", "E02"}, 1, 32)
	if !strings.Contains(got, epicActiveMarker) {
		t.Errorf("RenderEpicDropdown missing cursor marker %q", epicActiveMarker)
	}
}

func TestRenderEpicDropdown_containsHint(t *testing.T) {
	got := RenderEpicDropdown([]string{"E01"}, 0, 32)
	if !strings.Contains(got, "esc") {
		t.Error("RenderEpicDropdown missing esc hint")
	}
}

func TestRenderEpicDropdown_emptyShowsNone(t *testing.T) {
	got := RenderEpicDropdown(nil, 0, 32)
	if !strings.Contains(got, "(none)") {
		t.Error("RenderEpicDropdown with no epics should show (none)")
	}
}

func TestEpicIndex_found(t *testing.T) {
	epics := []string{"E01", "E02", "E03"}
	if got := epicIndex(epics, "E02"); got != 1 {
		t.Errorf("epicIndex = %d, want 1", got)
	}
}

func TestEpicIndex_notFound(t *testing.T) {
	if got := epicIndex([]string{"E01"}, "E99"); got != 0 {
		t.Errorf("epicIndex not-found = %d, want 0", got)
	}
}

func TestEpicIndex_empty(t *testing.T) {
	if got := epicIndex(nil, "E01"); got != 0 {
		t.Errorf("epicIndex empty = %d, want 0", got)
	}
}

// Update integration tests for epic dropdown

func TestUpdate_eKeyOpensDropdownNarrow(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Width = 80 // narrow: < 120
	m.Epics = []string{"E01", "E03"}
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("e")})
	updated := requireModel(t, got)
	if updated.Overlay != OverlayEpic {
		t.Errorf("Overlay = %q, want %q", updated.Overlay, OverlayEpic)
	}
}

func TestUpdate_eKeyOpensDropdownWide(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Width = 120
	m.Epics = []string{"E01", "E03"}
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("e")})
	updated := requireModel(t, got)
	if updated.Overlay != OverlayEpic {
		t.Errorf("Overlay = %q, want %q", updated.Overlay, OverlayEpic)
	}
}

func TestUpdate_epicDropdownEscCloses(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Overlay = OverlayEpic
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	updated := requireModel(t, got)
	if updated.Overlay != OverlayNone {
		t.Errorf("Overlay = %q after esc, want none", updated.Overlay)
	}
}

func TestUpdate_epicDropdownDownMovesCursor(t *testing.T) {
	m := NewModel(nil, "v1", "E01")
	m.Overlay = OverlayEpic
	m.Epics = []string{"E01", "E02", "E03"}
	m.EpicCursor = 0
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	updated := requireModel(t, got)
	if updated.EpicCursor != 1 {
		t.Errorf("EpicCursor = %d, want 1", updated.EpicCursor)
	}
}

func TestUpdate_epicDropdownUpMovesCursor(t *testing.T) {
	m := NewModel(nil, "v1", "E02")
	m.Overlay = OverlayEpic
	m.Epics = []string{"E01", "E02", "E03"}
	m.EpicCursor = 2
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	updated := requireModel(t, got)
	if updated.EpicCursor != 1 {
		t.Errorf("EpicCursor = %d, want 1", updated.EpicCursor)
	}
}

func TestUpdate_epicDropdownEnterSelectsEpic(t *testing.T) {
	tasks := []data.Task{
		{ID: "T1", Epic: "E01", Release: "v1", Column: data.ColumnPlanned},
		{ID: "T3", Epic: "E03", Release: "v1", Column: data.ColumnPlanned},
	}
	m := NewModel(tasks, "v1", "E01")
	m.Overlay = OverlayEpic
	m.Epics = []string{"E01", "E02", "E03"}
	m.EpicCursor = 2
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	updated := requireModel(t, got)
	if updated.SelectedEpic != "E03" {
		t.Errorf("SelectedEpic = %q, want %q", updated.SelectedEpic, "E03")
	}
	if updated.Overlay != OverlayNone {
		t.Errorf("Overlay = %q after enter, want none", updated.Overlay)
	}
	if got := len(updated.Tasks[data.ColumnPlanned]); got != 1 {
		t.Errorf("planned task count = %d, want 1 after epic selection", got)
	}
	if updated.Tasks[data.ColumnPlanned][0].ID != "T3" {
		t.Errorf("visible task = %q, want T3", updated.Tasks[data.ColumnPlanned][0].ID)
	}
}

func TestUpdate_epicDropdownDownClampedAtEnd(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Overlay = OverlayEpic
	m.Epics = []string{"E01", "E02"}
	m.EpicCursor = 1
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	updated := requireModel(t, got)
	if updated.EpicCursor != 1 {
		t.Errorf("EpicCursor = %d, want 1 (clamped)", updated.EpicCursor)
	}
}

func TestUpdate_epicDropdownUpClampedAtStart(t *testing.T) {
	m := NewModel(nil, "v1", "E01")
	m.Overlay = OverlayEpic
	m.Epics = []string{"E01", "E02"}
	m.EpicCursor = 0
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	updated := requireModel(t, got)
	if updated.EpicCursor != 0 {
		t.Errorf("EpicCursor = %d, want 0 (clamped)", updated.EpicCursor)
	}
}

func TestUpdate_overlayBlocksColumnNav(t *testing.T) {
	m := NewModel(nil, "v1", "E01")
	m.Overlay = OverlayEpic
	m.FocusedColumn = "planned"
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")})
	updated := requireModel(t, got)
	if updated.FocusedColumn != "planned" {
		t.Error("column nav should be blocked when overlay is open")
	}
}

func TestView_epicDropdownOverlayRendered(t *testing.T) {
	m := NewModel(nil, "v1", "E01")
	m.Width = 80
	m.Height = 24
	m.Overlay = OverlayEpic
	m.Epics = []string{"E01", "E02"}
	got := m.View()
	if !strings.Contains(got, "SELECT EPIC") {
		t.Error("View() with OverlayEpic missing SELECT EPIC")
	}
}

func TestView_epicDropdownKeepsBoardBehind(t *testing.T) {
	m := NewModel(nil, "v1", "E01")
	m.Width = 100
	m.Height = 24
	m.Overlay = OverlayEpic
	m.Epics = []string{"E01", "E02"}
	got := m.View()
	if !strings.Contains(got, "S A V E P O I N T") {
		t.Error("View() with OverlayEpic should keep board visible behind overlay")
	}
}

func TestView_epicSidebarOnWide(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Width = 120
	m.Epics = []string{"E01", "E03"}
	got := m.View()
	if !strings.Contains(got, "EPICS") {
		t.Error("View() at width>=120 missing EPICS header in sidebar")
	}
}
