package board

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/opencode/savepoint/internal/data"
)

func TestRenderDefectsOverlay_emptyState(t *testing.T) {
	got := RenderDefectsOverlay(nil, 0, 50)
	if !strings.Contains(got, "DEFECTS") {
		t.Error("RenderDefectsOverlay missing title")
	}
	plain := plainTerminal(got)
	if !strings.Contains(plain, "(no defects)") {
		t.Errorf("RenderDefectsOverlay empty state missing message; got %q", plain)
	}
}

func TestRenderDefectsOverlay_showsDefectRows(t *testing.T) {
	defects := []data.Defect{
		{ID: "D001-crash", Title: "App crashes on start", Severity: data.SeverityCritical, Status: data.ColumnPlanned},
		{ID: "D002-auth", Title: "Auth broken", Severity: data.SeverityHigh, Status: data.ColumnInProgress},
		{ID: "D003-typo", Title: "Typo in message", Severity: data.SeverityLow, Status: data.ColumnDone},
	}
	got := plainTerminal(RenderDefectsOverlay(defects, 0, 60))
	if !strings.Contains(got, "D001") {
		t.Error("overlay missing defect ID D001")
	}
	if !strings.Contains(got, "[⚠]") {
		t.Error("overlay missing critical severity tag")
	}
	if !strings.Contains(got, "App crashes on start") {
		t.Error("overlay missing defect title")
	}
}

func TestRenderDefectsOverlay_showsStatusSections(t *testing.T) {
	defects := []data.Defect{
		{ID: "D001", Title: "open", Severity: data.SeverityMedium, Status: data.ColumnPlanned},
		{ID: "D002", Title: "in progress", Severity: data.SeverityMedium, Status: data.ColumnInProgress},
		{ID: "D003", Title: "done", Severity: data.SeverityMedium, Status: data.ColumnDone},
	}
	got := plainTerminal(RenderDefectsOverlay(defects, 0, 60))
	if !strings.Contains(got, "OPEN") {
		t.Error("overlay missing OPEN section")
	}
	if !strings.Contains(got, "IN PROGRESS") {
		t.Error("overlay missing IN PROGRESS section")
	}
	if !strings.Contains(got, "RESOLVED") {
		t.Error("overlay missing RESOLVED section")
	}
}

func TestRenderDefectsOverlay_showsReference(t *testing.T) {
	defects := []data.Defect{
		{ID: "D001", Title: "crash", Severity: data.SeverityHigh, Status: data.ColumnPlanned, Reference: "E03/T002"},
	}
	got := plainTerminal(RenderDefectsOverlay(defects, 0, 60))
	if !strings.Contains(got, "E03/T002") {
		t.Errorf("overlay missing reference; got %q", got)
	}
}

func TestRenderDefectsOverlay_containsNavHint(t *testing.T) {
	got := RenderDefectsOverlay(nil, 0, 50)
	if !strings.Contains(got, "esc:close") {
		t.Error("overlay missing esc:close hint")
	}
}

func TestDefectsForOverlay_filtersByRelease(t *testing.T) {
	defects := []data.Defect{
		{ID: "D001", Release: "v1", Status: data.ColumnPlanned},
		{ID: "D002", Release: "v2", Status: data.ColumnPlanned},
	}
	got := defectsForOverlay(defects, "v1")
	if len(got) != 1 || got[0].ID != "D001" {
		t.Errorf("defectsForOverlay filtered = %v, want [D001]", got)
	}
}

func TestDefectsForOverlay_ordersByStatus(t *testing.T) {
	defects := []data.Defect{
		{ID: "D003", Status: data.ColumnDone},
		{ID: "D001", Status: data.ColumnPlanned},
		{ID: "D002", Status: data.ColumnInProgress},
	}
	got := defectsForOverlay(defects, "")
	if len(got) != 3 {
		t.Fatalf("len = %d, want 3", len(got))
	}
	if got[0].ID != "D001" || got[1].ID != "D002" || got[2].ID != "D003" {
		t.Errorf("order = [%s %s %s], want [D001 D002 D003]", got[0].ID, got[1].ID, got[2].ID)
	}
}

func TestUpdate_dOpensDefectOverlay(t *testing.T) {
	m := NewModel(nil, "v1", "E01")
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("d")})
	updated := requireModel(t, got)
	if updated.Overlay != OverlayDefect {
		t.Errorf("Overlay = %q, want %q", updated.Overlay, OverlayDefect)
	}
	if updated.DefectCursor != 0 {
		t.Errorf("DefectCursor = %d, want 0", updated.DefectCursor)
	}
}

func TestUpdate_defectOverlayEscCloses(t *testing.T) {
	m := NewModel(nil, "v1", "E01")
	m.Overlay = OverlayDefect
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	updated := requireModel(t, got)
	if updated.Overlay != OverlayNone {
		t.Errorf("Overlay = %q after esc, want none", updated.Overlay)
	}
}

func TestUpdate_defectOverlayQCloses(t *testing.T) {
	m := NewModel(nil, "v1", "E01")
	m.Overlay = OverlayDefect
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	updated := requireModel(t, got)
	if updated.Overlay != OverlayNone {
		t.Errorf("Overlay = %q after q, want none", updated.Overlay)
	}
}

func TestUpdate_defectOverlayNavigatesDown(t *testing.T) {
	m := NewModel(nil, "v1", "E01")
	m.SelectedRelease = "v1"
	m.AllDefects = []data.Defect{
		{Release: "v1", ID: "D001", Status: data.ColumnPlanned},
		{Release: "v1", ID: "D002", Status: data.ColumnPlanned},
	}
	m.Overlay = OverlayDefect
	m.DefectCursor = 0

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	updated := requireModel(t, got)
	if updated.DefectCursor != 1 {
		t.Errorf("DefectCursor = %d, want 1", updated.DefectCursor)
	}
}

func TestUpdate_defectOverlayNavigatesUp(t *testing.T) {
	m := NewModel(nil, "v1", "E01")
	m.SelectedRelease = "v1"
	m.AllDefects = []data.Defect{
		{Release: "v1", ID: "D001", Status: data.ColumnPlanned},
		{Release: "v1", ID: "D002", Status: data.ColumnPlanned},
	}
	m.Overlay = OverlayDefect
	m.DefectCursor = 1

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	updated := requireModel(t, got)
	if updated.DefectCursor != 0 {
		t.Errorf("DefectCursor = %d, want 0", updated.DefectCursor)
	}
}

func TestUpdate_defectOverlayCursorClampsAtTop(t *testing.T) {
	m := NewModel(nil, "v1", "E01")
	m.Overlay = OverlayDefect
	m.DefectCursor = 0

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	updated := requireModel(t, got)
	if updated.DefectCursor != 0 {
		t.Errorf("DefectCursor = %d, want 0 (clamped at top)", updated.DefectCursor)
	}
}

func TestUpdate_defectOverlayCursorClampsAtBottom(t *testing.T) {
	m := NewModel(nil, "v1", "E01")
	m.SelectedRelease = "v1"
	m.AllDefects = []data.Defect{
		{Release: "v1", ID: "D001", Status: data.ColumnPlanned},
	}
	m.Overlay = OverlayDefect
	m.DefectCursor = 0

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	updated := requireModel(t, got)
	if updated.DefectCursor != 0 {
		t.Errorf("DefectCursor = %d, want 0 (clamped at bottom, only 1 defect)", updated.DefectCursor)
	}
}

func TestView_defectOverlayRendered(t *testing.T) {
	m := NewModel(nil, "v1", "E01")
	m.Width = 100
	m.Height = 30
	m.SelectedRelease = "v1"
	m.AllDefects = []data.Defect{
		{Release: "v1", ID: "D001-crash", Title: "App crash", Severity: data.SeverityCritical, Status: data.ColumnPlanned},
	}
	m.Overlay = OverlayDefect
	got := plainTerminal(m.View())
	if !strings.Contains(got, "DEFECTS") {
		t.Error("View() with OverlayDefect missing DEFECTS header")
	}
	if !strings.Contains(got, "D001") {
		t.Error("View() with OverlayDefect missing defect ID")
	}
}

func TestView_defectOverlayEmptyState(t *testing.T) {
	m := NewModel(nil, "v1", "E01")
	m.Width = 100
	m.Height = 30
	m.SelectedRelease = "v1"
	m.Overlay = OverlayDefect
	got := plainTerminal(m.View())
	if !strings.Contains(got, "(no defects)") {
		t.Errorf("View() with empty defects missing empty state; got %q", got)
	}
}

func TestRenderHelp_containsDefectShortcut(t *testing.T) {
	got := RenderHelp(60)
	if !strings.Contains(got, "d") {
		t.Error("RenderHelp missing d shortcut")
	}
	if !strings.Contains(got, "defect") {
		t.Error("RenderHelp missing defects overlay description")
	}
}
