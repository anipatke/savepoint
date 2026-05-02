package board

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/opencode/savepoint/internal/data"
)

func TestRenderEpicSidebar_containsEpicsHeader(t *testing.T) {
	got := RenderEpicSidebar([]string{"E01", "E02"}, "E01", 28, false, 0, nil)
	if !strings.Contains(got, "EPICS") {
		t.Error("RenderEpicSidebar missing EPICS header")
	}
}

func TestRenderEpicSidebar_activeEpicMarked(t *testing.T) {
	got := RenderEpicSidebar([]string{"E01", "E02"}, "E01", 28, false, 0, nil)
	if !strings.Contains(got, epicActiveMarker) {
		t.Errorf("RenderEpicSidebar missing active marker %q", epicActiveMarker)
	}
}

func TestRenderEpicSidebar_focusedCursorMarked(t *testing.T) {
	got := RenderEpicSidebar([]string{"E01", "E02"}, "E01", 28, true, 1, nil)
	if !strings.Contains(got, epicActiveMarker+"   E02") {
		t.Errorf("RenderEpicSidebar focused cursor missing marker, got %q", got)
	}
}

func TestRenderEpicSidebar_allEpicsPresent(t *testing.T) {
	epics := []string{"E01-foo", "E02-bar", "E03-baz"}
	got := RenderEpicSidebar(epics, "E01-foo", 32, false, 0, nil)
	for _, e := range epics {
		if !strings.Contains(got, e) {
			t.Errorf("RenderEpicSidebar missing epic %q", e)
		}
	}
}

func TestRenderEpicSidebar_emptyEpicsFallback(t *testing.T) {
	got := RenderEpicSidebar(nil, "E03", 28, false, 0, nil)
	if !strings.Contains(got, "E03") {
		t.Error("RenderEpicSidebar with empty list should show selected epic")
	}
}

func TestRenderEpicSidebar_emptyBothShowsNone(t *testing.T) {
	got := RenderEpicSidebar(nil, "", 28, false, 0, nil)
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

func TestUpdate_leftFromPlannedFocusesEpicPanelWide(t *testing.T) {
	m := NewModel(nil, "v1", "E02")
	m.Width = 120
	m.Epics = []string{"E01", "E02", "E03"}

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("h")})
	updated := requireModel(t, got)

	if !updated.EpicPanelFocus {
		t.Fatal("EpicPanelFocus = false, want true")
	}
	if updated.EpicPanelCursor != 1 {
		t.Errorf("EpicPanelCursor = %d, want 1", updated.EpicPanelCursor)
	}
	if updated.FocusedColumn != data.ColumnPlanned {
		t.Errorf("FocusedColumn = %q, want planned", updated.FocusedColumn)
	}
}

func TestUpdate_leftFromPlannedDoesNotFocusEpicPanelNarrow(t *testing.T) {
	m := NewModel(nil, "v1", "E02")
	m.Width = 100
	m.Epics = []string{"E01", "E02"}

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("h")})
	updated := requireModel(t, got)

	if updated.EpicPanelFocus {
		t.Fatal("EpicPanelFocus = true, want false")
	}
	if updated.FocusedColumn != data.ColumnDone {
		t.Errorf("FocusedColumn = %q, want done", updated.FocusedColumn)
	}
}

func TestUpdate_leftFromPlannedDoesNotFocusEmptyEpicPanel(t *testing.T) {
	m := NewModel(nil, "v1", "")
	m.Width = 120

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("h")})
	updated := requireModel(t, got)

	if updated.EpicPanelFocus {
		t.Fatal("EpicPanelFocus = true, want false with no epics")
	}
	if updated.EpicPanelCursor != 0 {
		t.Errorf("EpicPanelCursor = %d, want 0", updated.EpicPanelCursor)
	}
}

func TestUpdate_windowResizeClearsEpicPanelFocusWhenHidden(t *testing.T) {
	m := NewModel(nil, "v1", "E01")
	m.Width = 120
	m.Epics = []string{"E01"}
	m.EpicPanelFocus = true

	got, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 24})
	updated := requireModel(t, got)

	if updated.EpicPanelFocus {
		t.Fatal("EpicPanelFocus = true, want false when panel is hidden")
	}
}

func TestUpdate_epicPanelDownUpClamped(t *testing.T) {
	m := NewModel(nil, "v1", "E01")
	m.Width = 120
	m.EpicPanelFocus = true
	m.Epics = []string{"E01", "E02"}

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	updated := requireModel(t, got)
	if updated.EpicPanelCursor != 1 {
		t.Errorf("EpicPanelCursor after down = %d, want 1", updated.EpicPanelCursor)
	}

	got, _ = updated.Update(tea.KeyMsg{Type: tea.KeyDown})
	updated = requireModel(t, got)
	if updated.EpicPanelCursor != 1 {
		t.Errorf("EpicPanelCursor after clamped down = %d, want 1", updated.EpicPanelCursor)
	}

	got, _ = updated.Update(tea.KeyMsg{Type: tea.KeyUp})
	updated = requireModel(t, got)
	if updated.EpicPanelCursor != 0 {
		t.Errorf("EpicPanelCursor after up = %d, want 0", updated.EpicPanelCursor)
	}
}

func TestUpdate_epicPanelEnterOpensDetailOverlay(t *testing.T) {
	tasks := []data.Task{
		{ID: "T1", Epic: "E01", Release: "v1", Column: data.ColumnPlanned},
		{ID: "T2", Epic: "E02", Release: "v1", Column: data.ColumnPlanned},
	}
	m := NewModel(tasks, "v1", "E01")
	m.Width = 120
	m.Epics = []string{"E01", "E02"}
	m.EpicPanelFocus = true
	m.EpicPanelCursor = 1
	m.FocusedTask = 3
	m.DetailOffset = 2

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	updated := requireModel(t, got)

	if updated.Overlay != OverlayEpicDetail {
		t.Errorf("Overlay = %q, want %q", updated.Overlay, OverlayEpicDetail)
	}
	if updated.EpicDetailOffset != 0 {
		t.Errorf("EpicDetailOffset = %d, want 0", updated.EpicDetailOffset)
	}
	if !updated.EpicPanelFocus {
		t.Error("EpicPanelFocus should remain true after enter")
	}
	// SelectedEpic unchanged; Enter now opens detail, not selects
	if updated.SelectedEpic != "E01" {
		t.Errorf("SelectedEpic = %q, want E01 (unchanged)", updated.SelectedEpic)
	}
}

func TestUpdate_epicPanelRightReturnsToPlanned(t *testing.T) {
	m := NewModel(nil, "v1", "E01")
	m.Width = 120
	m.EpicPanelFocus = true
	m.Epics = []string{"E01"}
	m.FocusedColumn = data.ColumnDone
	m.FocusedTask = 2

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")})
	updated := requireModel(t, got)

	if updated.EpicPanelFocus {
		t.Fatal("EpicPanelFocus = true, want false")
	}
	if updated.FocusedColumn != data.ColumnPlanned {
		t.Errorf("FocusedColumn = %q, want planned", updated.FocusedColumn)
	}
	if updated.FocusedTask != 0 {
		t.Errorf("FocusedTask = %d, want 0", updated.FocusedTask)
	}
}

func TestUpdate_overlayBlocksEpicPanelNav(t *testing.T) {
	m := NewModel(nil, "v1", "E01")
	m.Width = 120
	m.EpicPanelFocus = true
	m.Epics = []string{"E01", "E02"}
	m.Overlay = OverlayHelp

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	updated := requireModel(t, got)

	if updated.EpicPanelCursor != 0 {
		t.Errorf("EpicPanelCursor = %d, want 0 while overlay open", updated.EpicPanelCursor)
	}
}

func TestUpdate_epicPanelFocusAllowsGlobalQuit(t *testing.T) {
	m := NewModel(nil, "v1", "E01")
	m.Width = 120
	m.EpicPanelFocus = true
	m.Epics = []string{"E01"}

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if cmd == nil {
		t.Fatal("expected tea.Quit cmd from q while epic panel focused, got nil")
	}
}

func TestUpdate_epicPanelFocusAllowsEpicDropdown(t *testing.T) {
	m := NewModel(nil, "v1", "E02")
	m.Width = 120
	m.EpicPanelFocus = true
	m.Epics = []string{"E01", "E02"}

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("e")})
	updated := requireModel(t, got)

	if updated.Overlay != OverlayEpic {
		t.Errorf("Overlay = %q, want %q", updated.Overlay, OverlayEpic)
	}
	if updated.EpicCursor != 1 {
		t.Errorf("EpicCursor = %d, want 1", updated.EpicCursor)
	}
}

func TestUpdate_epicPanelFocusAllowsReleaseDropdown(t *testing.T) {
	m := NewModel(nil, "v1", "E01")
	m.Width = 120
	m.EpicPanelFocus = true
	m.Epics = []string{"E01"}
	m.Releases = []string{"v1", "v1.1"}
	m.SelectedRelease = "v1.1"

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")})
	updated := requireModel(t, got)

	if updated.Overlay != OverlayRelease {
		t.Errorf("Overlay = %q, want %q", updated.Overlay, OverlayRelease)
	}
	if updated.ReleaseCursor != 1 {
		t.Errorf("ReleaseCursor = %d, want 1", updated.ReleaseCursor)
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

func TestUpdate_epicDetailOverlayEscCloses(t *testing.T) {
	m := NewModel(nil, "v1", "E01")
	m.Width = 120
	m.EpicPanelFocus = true
	m.Epics = []string{"E01"}
	m.Overlay = OverlayEpicDetail
	m.EpicDetailContent = "# E01 Detail"
	m.EpicDetailOffset = 3

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	updated := requireModel(t, got)

	if updated.Overlay != OverlayNone {
		t.Errorf("Overlay = %q, want none after esc", updated.Overlay)
	}
	if !updated.EpicPanelFocus {
		t.Error("EpicPanelFocus should remain true after closing detail overlay")
	}
}

func TestUpdate_epicDetailOverlayScrollUpDown(t *testing.T) {
	m := NewModel(nil, "v1", "E01")
	m.Width = 120
	m.EpicPanelFocus = true
	m.Epics = []string{"E01"}
	m.Overlay = OverlayEpicDetail
	m.EpicDetailOffset = 2

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	updated := requireModel(t, got)
	if updated.EpicDetailOffset != 3 {
		t.Errorf("EpicDetailOffset after down = %d, want 3", updated.EpicDetailOffset)
	}

	got, _ = updated.Update(tea.KeyMsg{Type: tea.KeyUp})
	updated = requireModel(t, got)
	if updated.EpicDetailOffset != 2 {
		t.Errorf("EpicDetailOffset after up = %d, want 2", updated.EpicDetailOffset)
	}

	got, _ = updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	updated = requireModel(t, got)
	if updated.EpicDetailOffset != 1 {
		t.Errorf("EpicDetailOffset after k = %d, want 1", updated.EpicDetailOffset)
	}

	// Clamp at 0
	updated.EpicDetailOffset = 0
	got, _ = updated.Update(tea.KeyMsg{Type: tea.KeyUp})
	updated = requireModel(t, got)
	if updated.EpicDetailOffset != 0 {
		t.Errorf("EpicDetailOffset should not go below 0, got %d", updated.EpicDetailOffset)
	}
}

func TestUpdate_epicDetailOverlayPgUpDown(t *testing.T) {
	m := NewModel(nil, "v1", "E01")
	m.Width = 120
	m.Height = 30
	m.Epics = []string{"E01"}
	m.Overlay = OverlayEpicDetail
	m.EpicDetailOffset = 20

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyPgUp})
	updated := requireModel(t, got)
	if updated.EpicDetailOffset >= 20 {
		t.Errorf("EpicDetailOffset after pgup = %d, should decrease from 20", updated.EpicDetailOffset)
	}
	if updated.EpicDetailOffset < 0 {
		t.Errorf("EpicDetailOffset = %d, should not go below 0", updated.EpicDetailOffset)
	}

	updated.EpicDetailOffset = 0
	got, _ = updated.Update(tea.KeyMsg{Type: tea.KeyPgDown})
	updated = requireModel(t, got)
	if updated.EpicDetailOffset <= 0 {
		t.Errorf("EpicDetailOffset after pgdown = %d, should increase from 0", updated.EpicDetailOffset)
	}
}

func TestView_epicDetailOverlayRendered(t *testing.T) {
	m := NewModel(nil, "v1", "E01")
	m.Width = 120
	m.Height = 30
	m.Epics = []string{"E01"}
	m.Overlay = OverlayEpicDetail
	m.EpicDetailContent = "# My Epic\n\n## Purpose\nDoes things."

	got := m.View()
	if !strings.Contains(got, "EPIC DETAIL") {
		t.Error("View() with OverlayEpicDetail missing EPIC DETAIL header")
	}
}

func TestView_epicDetailOverlayNoContent(t *testing.T) {
	m := NewModel(nil, "v1", "E01")
	m.Width = 120
	m.Height = 30
	m.Epics = []string{"E01"}
	m.Overlay = OverlayEpicDetail
	m.EpicDetailContent = "(no detail available)"

	got := m.View()
	if !strings.Contains(got, "no detail available") {
		t.Error("View() with missing epic detail should show 'no detail available'")
	}
}

func TestRenderEpicDetail_stripsMarkdownHeadings(t *testing.T) {
	content := "---\ntype: epic-design\n---\n# Epic E01\n\n## Purpose\nDoes things."
	got := RenderEpicDetail("E01-test", content, 60, 40, 0)
	if !strings.Contains(got, "EPIC DETAIL") {
		t.Error("RenderEpicDetail missing EPIC DETAIL header")
	}
	if strings.Contains(got, "# Epic E01") {
		t.Error("RenderEpicDetail should strip raw markdown heading prefix")
	}
}

func TestRenderEpicDetail_noDetailFallback(t *testing.T) {
	got := RenderEpicDetail("E01-test", "(no detail available)", 60, 40, 0)
	if !strings.Contains(got, "no detail available") {
		t.Error("RenderEpicDetail fallback message missing")
	}
}
