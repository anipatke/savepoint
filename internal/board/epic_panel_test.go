package board

import (
	"os"
	"path/filepath"
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

func TestSliceIndex_found(t *testing.T) {
	epics := []string{"E01", "E02", "E03"}
	if got := sliceIndex(epics, "E02"); got != 1 {
		t.Errorf("sliceIndex = %d, want 1", got)
	}
}

func TestSliceIndex_notFound(t *testing.T) {
	if got := sliceIndex([]string{"E01"}, "E99"); got != 0 {
		t.Errorf("sliceIndex not-found = %d, want 0", got)
	}
}

func TestSliceIndex_empty(t *testing.T) {
	if got := sliceIndex(nil, "E01"); got != 0 {
		t.Errorf("sliceIndex empty = %d, want 0", got)
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
	got := RenderEpicDetail("E01-test", content, 60, 40, 0, 0)
	if !strings.Contains(got, "EPIC DETAIL") {
		t.Error("RenderEpicDetail missing EPIC DETAIL header")
	}
	if strings.Contains(got, "# Epic E01") {
		t.Error("RenderEpicDetail should strip raw markdown heading prefix")
	}
}

func TestRenderEpicDetail_noDetailFallback(t *testing.T) {
	got := RenderEpicDetail("E01-test", "(no detail available)", 60, 40, 0, 0)
	if !strings.Contains(got, "no detail available") {
		t.Error("RenderEpicDetail fallback message missing")
	}
}

func TestRenderEpicDetail_tabIndicatorDetailActive(t *testing.T) {
	got := RenderEpicDetail("E01-test", "content", 60, 40, 0, 0)
	if !strings.Contains(got, "DETAIL [1]") {
		t.Error("RenderEpicDetail tab=0: missing DETAIL [1] indicator")
	}
	if !strings.Contains(got, "AUDIT [2]") {
		t.Error("RenderEpicDetail tab=0: missing AUDIT [2] indicator")
	}
}

func TestRenderEpicDetail_tabIndicatorAuditActive(t *testing.T) {
	got := RenderEpicDetail("E01-test", "content", 60, 40, 0, 1)
	if !strings.Contains(got, "DETAIL [1]") {
		t.Error("RenderEpicDetail tab=1: missing DETAIL [1] indicator")
	}
	if !strings.Contains(got, "AUDIT [2]") {
		t.Error("RenderEpicDetail tab=1: missing AUDIT [2] indicator")
	}
}

func TestRenderEpicAuditTab_header(t *testing.T) {
	got := RenderEpicAuditTab("E06-test", "# Audit\n\n## Main Findings\nAll good.", 60, 40, 0, 1)
	if !strings.Contains(got, "EPIC AUDIT") {
		t.Error("RenderEpicAuditTab missing EPIC AUDIT header")
	}
}

func TestRenderEpicAuditTab_noContent(t *testing.T) {
	got := RenderEpicAuditTab("E06-test", "(no audit available)", 60, 40, 0, 1)
	if !strings.Contains(got, "no audit available") {
		t.Error("RenderEpicAuditTab fallback message missing")
	}
}

func TestRenderEpicAuditTab_emptyContent(t *testing.T) {
	got := RenderEpicAuditTab("E06-test", "", 60, 40, 0, 1)
	if !strings.Contains(got, "no audit available") {
		t.Error("RenderEpicAuditTab empty content should show fallback")
	}
}

func TestRenderEpicAuditTab_stripsFrontmatter(t *testing.T) {
	content := "---\ntype: audit\n---\n# E06 Audit\n\n## Main Findings\nLooks good."
	got := RenderEpicAuditTab("E06-test", content, 60, 40, 0, 1)
	if strings.Contains(got, "type: audit") {
		t.Error("RenderEpicAuditTab should strip frontmatter")
	}
	if !strings.Contains(got, "EPIC AUDIT") {
		t.Error("RenderEpicAuditTab missing header after frontmatter strip")
	}
}

func TestRenderEpicAuditTab_checkboxDonePresent(t *testing.T) {
	content := "## Code Style Review\n- [x] One job per file\n- [ ] One-sentence functions"
	got := RenderEpicAuditTab("E06-test", content, 60, 40, 0, 1)
	if !strings.Contains(got, "One job per file") {
		t.Error("RenderEpicAuditTab missing done checkbox text")
	}
	if !strings.Contains(got, "One-sentence functions") {
		t.Error("RenderEpicAuditTab missing undone checkbox text")
	}
}

func TestRenderEpicAuditTab_scrollFooter(t *testing.T) {
	got := RenderEpicAuditTab("E06-test", "# Audit", 60, 40, 0, 1)
	if !strings.Contains(got, "esc:close") {
		t.Error("RenderEpicAuditTab missing esc:close footer")
	}
}

func TestRenderEpicAuditTab_tabIndicator(t *testing.T) {
	got := RenderEpicAuditTab("E06-test", "# Audit", 60, 40, 0, 1)
	if !strings.Contains(got, "DETAIL [1]") {
		t.Error("RenderEpicAuditTab missing DETAIL [1] indicator")
	}
	if !strings.Contains(got, "AUDIT [2]") {
		t.Error("RenderEpicAuditTab missing AUDIT [2] indicator")
	}
}

func TestRenderEpicAuditTab_mainFindingsVisible(t *testing.T) {
	content := "## Main Findings\nAudit summary is visible.\n\n## Proposed Changes\n### Target File\nAGENTS.md\n"
	got := RenderEpicAuditTab("E06-test", content, 80, 50, 0, 1)
	if !strings.Contains(got, "Audit summary is visible") {
		t.Error("RenderEpicAuditTab should render Main Findings body")
	}
	if strings.Contains(got, "Target File") || strings.Contains(got, "AGENTS.md") {
		t.Error("RenderEpicAuditTab should not render Proposed Changes admin blocks")
	}
}

func TestRenderEpicAuditTab_qualityReviewHidden(t *testing.T) {
	content := "## Quality Review\nOld quality section.\n\n## Code Style Review\n- [ ] One job per file\n"
	got := RenderEpicAuditTab("E06-test", content, 80, 50, 0, 1)
	if strings.Contains(got, "Old quality section") {
		t.Error("RenderEpicAuditTab should not render superseded Quality Review section")
	}
	if !strings.Contains(got, "One job per file") {
		t.Error("RenderEpicAuditTab should render Code Style Review")
	}
}

func TestRenderEpicAuditTab_hiddenHeadingsRequireExactMatch(t *testing.T) {
	content := "## Proposed Changes Appendix\nNear-match section is visible.\n\n## Proposed Changes\nHidden admin section.\n"
	got := RenderEpicAuditTab("E06-test", content, 80, 50, 0, 1)
	if !strings.Contains(got, "Near-match section is visible") {
		t.Error("RenderEpicAuditTab should render headings that only partially match hidden headings")
	}
	if strings.Contains(got, "Hidden admin section") {
		t.Error("RenderEpicAuditTab should hide exact Proposed Changes section")
	}
}

func TestRenderEpicAuditTab_allCodeStyleRules(t *testing.T) {
	rules := []string{
		"One job per file",
		"One-sentence functions",
		"Test branches",
		"Types are documentation",
		"Build, don't speculate",
		"Errors at boundaries",
		"One source of truth",
		"Comments explain WHY",
		"Content in data files",
		"Small diffs",
	}
	content := "## Code Style Review\n"
	for _, r := range rules {
		content += "- [ ] " + r + "\n"
	}
	got := RenderEpicAuditTab("E06-test", content, 80, 50, 0, 1)
	for _, r := range rules {
		if !strings.Contains(got, r) {
			t.Errorf("RenderEpicAuditTab missing code style rule %q", r)
		}
	}
}

// TestView_epicAuditTabRendered verifies View() uses RenderEpicAuditTab when EpicDetailTab=1.
func TestView_epicAuditTabRendered(t *testing.T) {
	m := NewModel(nil, "v1.1", "E06-audit-command")
	m.Width = 120
	m.Height = 30
	m.Epics = []string{"E06-audit-command"}
	m.Overlay = OverlayEpicDetail
	m.EpicDetailTab = 1
	m.EpicAuditContent = "# Audit Findings: E06\n\n## Main Findings\nAll good.\n\n## Code Style Review\n- [x] One job per file\n"

	got := m.View()
	if !strings.Contains(got, "EPIC AUDIT") {
		t.Error("View() with EpicDetailTab=1 missing EPIC AUDIT header")
	}
	if strings.Contains(got, "EPIC DETAIL") {
		t.Error("View() with EpicDetailTab=1 should not render EPIC DETAIL header")
	}
}

// TestAuditWorkflow_fullEndToEnd exercises the full audit workflow:
// create E##-Audit.md on disk, open overlay, press 2, verify content loads and renders.
func TestAuditWorkflow_fullEndToEnd(t *testing.T) {
	root := t.TempDir()
	epicSlug := "E06-audit-command"
	epicDir := filepath.Join(root, "releases", "v1.1", "epics", epicSlug)
	if err := os.MkdirAll(epicDir, 0755); err != nil {
		t.Fatal(err)
	}

	auditContent := `---
type: audit-findings
audited: 2026-05-03
---
# Audit Findings: E06 Agent Audit + Audit Tab

## Main Findings
All acceptance criteria met.

## Code Style Review
- [x] One job per file
- [x] One-sentence functions
- [x] Test branches
- [x] Types are documentation
- [x] Build, don't speculate
- [x] Errors at boundaries
- [x] One source of truth
- [x] Comments explain WHY
- [x] Content in data files
- [x] Small diffs
`
	if err := os.WriteFile(filepath.Join(epicDir, "E06-Audit.md"), []byte(auditContent), 0644); err != nil {
		t.Fatal(err)
	}

	tasks := []data.Task{
		{ID: "E06-audit-command/T009-integration", Release: "v1.1", Epic: epicSlug, Column: data.ColumnPlanned},
	}
	m := NewModel(tasks, "v1.1", epicSlug)
	m.Root = root
	m.Epics = []string{epicSlug}
	m.EpicPanelCursor = 0
	m.Width = 120
	m.Height = 40

	// Open detail overlay (tab=0)
	m.openEpicDetailOverlay()
	if m.Overlay != OverlayEpicDetail {
		t.Fatal("overlay not opened")
	}
	if m.EpicDetailTab != 0 {
		t.Errorf("EpicDetailTab = %d, want 0 on open", m.EpicDetailTab)
	}

	// Press 2 → switch to audit tab, load content
	got, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("2")})
	updated := requireModel(t, got)

	if updated.EpicDetailTab != 1 {
		t.Errorf("EpicDetailTab = %d, want 1 after pressing 2", updated.EpicDetailTab)
	}

	msg := cmd()
	got2, _ := updated.Update(msg)
	updated2 := requireModel(t, got2)
	if updated2.EpicAuditContent == "" || updated2.EpicAuditContent == "(no audit available)" {
		t.Errorf("EpicAuditContent not loaded: %q", updated2.EpicAuditContent)
	}

	// Verify View() renders audit content
	view := updated2.View()
	if !strings.Contains(view, "EPIC AUDIT") {
		t.Error("View() after tab switch missing EPIC AUDIT")
	}
	if !strings.Contains(view, "One job per file") {
		t.Error("View() after tab switch missing code style rule")
	}

	// Press 1 → switch back to detail tab
	got, _ = updated2.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("1")})
	updated = requireModel(t, got)
	if updated.EpicDetailTab != 0 {
		t.Errorf("EpicDetailTab = %d, want 0 after pressing 1", updated.EpicDetailTab)
	}

	// Press esc → overlay closes
	got, _ = updated.Update(tea.KeyMsg{Type: tea.KeyEsc})
	updated = requireModel(t, got)
	if updated.Overlay != OverlayNone {
		t.Errorf("Overlay = %q, want none after esc", updated.Overlay)
	}
}

func TestRenderEpicAuditTab_v11AuditFiles(t *testing.T) {
	files := []struct {
		path string
		want string
	}{
		{filepath.Join("..", "..", ".savepoint", "releases", "v1.1", "epics", "E02-cross-platform-compatibility", "E02-Audit.md"), "cross-platform build work"},
		{filepath.Join("..", "..", ".savepoint", "releases", "v1.1", "epics", "E03-ui-visual-refinement", "E03-Audit.md"), "visual refinement work"},
		{filepath.Join("..", "..", ".savepoint", "releases", "v1.1", "epics", "E04-epic-navigation", "E04-Audit.md"), "wide-screen epic navigation"},
		{filepath.Join("..", "..", ".savepoint", "releases", "v1.1", "epics", "E05-tasking-permissions", "E05-Audit.md"), "tasking-permissions shift"},
		{filepath.Join("..", "..", ".savepoint", "releases", "v1.1", "epics", "E06-audit-command", "E06-Audit.md"), "agent-led"},
	}

	for _, tt := range files {
		content, err := os.ReadFile(tt.path)
		if err != nil {
			t.Fatalf("read %s: %v", tt.path, err)
		}
		if !strings.Contains(string(content), tt.want) {
			t.Fatalf("fixture %s missing %q", tt.path, tt.want)
		}
		got := RenderEpicAuditTab(filepath.Base(filepath.Dir(tt.path)), string(content), 80, 40, 0, 1)
		if !strings.Contains(got, tt.want) {
			t.Errorf("RenderEpicAuditTab(%s) missing %q", tt.path, tt.want)
		}
		if strings.Contains(got, "Target File") {
			t.Errorf("RenderEpicAuditTab(%s) should not render Proposed Changes", tt.path)
		}
		if strings.Contains(got, "Boundaries") || strings.Contains(got, "Implemented as") || strings.Contains(got, "Implemented As") {
			t.Errorf("RenderEpicAuditTab(%s) should only render visible audit sections", tt.path)
		}
	}
}
