package board

import (
	"os"
	"path/filepath"
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
		{ID: "D001-crash", Title: "App crashes on start", Severity: data.SeverityCritical, Status: data.DefectOpen},
		{ID: "D002-auth", Title: "Auth broken", Severity: data.SeverityHigh, Status: data.DefectInProgress},
		{ID: "D003-typo", Title: "Typo in message", Severity: data.SeverityLow, Status: data.DefectResolved},
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
		{ID: "D001", Title: "open", Severity: data.SeverityMedium, Status: data.DefectOpen},
		{ID: "D002", Title: "in progress", Severity: data.SeverityMedium, Status: data.DefectInProgress},
		{ID: "D003", Title: "done", Severity: data.SeverityMedium, Status: data.DefectResolved},
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
		{ID: "D001", Title: "crash", Severity: data.SeverityHigh, Status: data.DefectOpen, Reference: "E03/T002"},
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
	if !strings.Contains(got, "space:resolve") {
		t.Error("overlay missing space resolve hint")
	}
}

func TestDefectsForOverlay_filtersByRelease(t *testing.T) {
	defects := []data.Defect{
		{ID: "D001", Release: "v1", Status: data.DefectOpen},
		{ID: "D002", Release: "v2", Status: data.DefectOpen},
	}
	got := defectsForOverlay(defects, "v1")
	if len(got) != 1 || got[0].ID != "D001" {
		t.Errorf("defectsForOverlay filtered = %v, want [D001]", got)
	}
}

func TestDefectsForOverlay_ordersByStatus(t *testing.T) {
	defects := []data.Defect{
		{ID: "D003", Status: data.DefectResolved},
		{ID: "D001", Status: data.DefectOpen},
		{ID: "D002", Status: data.DefectInProgress},
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
		{Release: "v1", ID: "D001", Status: data.DefectOpen},
		{Release: "v1", ID: "D002", Status: data.DefectOpen},
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
		{Release: "v1", ID: "D001", Status: data.DefectOpen},
		{Release: "v1", ID: "D002", Status: data.DefectOpen},
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
		{Release: "v1", ID: "D001", Status: data.DefectOpen},
	}
	m.Overlay = OverlayDefect
	m.DefectCursor = 0

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	updated := requireModel(t, got)
	if updated.DefectCursor != 0 {
		t.Errorf("DefectCursor = %d, want 0 (clamped at bottom, only 1 defect)", updated.DefectCursor)
	}
}

func TestUpdate_defectOverlaySpaceResolvesPlannedDefect(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "D001.md")
	content := `---
id: v1/D001
release: v1
status: open
severity: high
title: "Crash"
reference: E01/T001
---

# Body`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	m := NewModel(nil, "v1", "E01")
	m.SelectedRelease = "v1"
	m.AllDefects = []data.Defect{{
		ID: "v1/D001", Release: "v1", Status: data.DefectOpen,
		Severity: data.SeverityHigh, Title: "Crash", Path: path, Mtime: fi.ModTime(),
	}}
	m.Overlay = OverlayDefect

	got, cmd := m.Update(tea.KeyMsg{Type: tea.KeySpace})
	updated := requireModel(t, got)
	if cmd == nil {
		t.Fatal("space on planned defect returned nil command")
	}
	msg := cmd()
	got, _ = updated.Update(msg)
	updated = requireModel(t, got)

	if updated.Overlay != OverlayDefect {
		t.Errorf("Overlay = %q, want defect", updated.Overlay)
	}
	if updated.DefectCursor != 0 {
		t.Errorf("DefectCursor = %d, want 0", updated.DefectCursor)
	}
	if updated.AllDefects[0].Status != data.DefectResolved {
		t.Errorf("Status = %q, want resolved", updated.AllDefects[0].Status)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(result), "status: resolved") {
		t.Errorf("defect file missing status: resolved; got %s", result)
	}
	if !strings.Contains(string(result), "reference: E01/T001") {
		t.Error("unrelated frontmatter field not preserved")
	}
}

func TestUpdate_defectOverlaySpaceDoneIsNoop(t *testing.T) {
	m := NewModel(nil, "v1", "E01")
	m.SelectedRelease = "v1"
	m.AllDefects = []data.Defect{{Release: "v1", ID: "D001", Status: data.DefectResolved}}
	m.Overlay = OverlayDefect

	got, cmd := m.Update(tea.KeyMsg{Type: tea.KeySpace})
	updated := requireModel(t, got)

	if cmd != nil {
		t.Fatal("space on done defect returned command, want no-op")
	}
	if updated.StatusMessage != "Defect already resolved" {
		t.Errorf("StatusMessage = %q, want already resolved", updated.StatusMessage)
	}
	if updated.AllDefects[0].Status != data.DefectResolved {
		t.Errorf("Status = %q, want resolved", updated.AllDefects[0].Status)
	}
}

func TestUpdate_defectOverlaySpaceInProgressIsNoop(t *testing.T) {
	m := NewModel(nil, "v1", "E01")
	m.SelectedRelease = "v1"
	m.AllDefects = []data.Defect{{Release: "v1", ID: "D001", Status: data.DefectInProgress, Stage: data.StageBuild}}
	m.Overlay = OverlayDefect

	got, cmd := m.Update(tea.KeyMsg{Type: tea.KeySpace})
	updated := requireModel(t, got)

	if cmd != nil {
		t.Fatal("space on in-progress defect returned command, want clear no-op")
	}
	if !strings.Contains(updated.StatusMessage, "in progress") {
		t.Errorf("StatusMessage = %q, want in progress no-op", updated.StatusMessage)
	}
	if updated.AllDefects[0].Status != data.DefectInProgress {
		t.Errorf("Status = %q, want in_progress (noop)", updated.AllDefects[0].Status)
	}
}

func TestView_defectOverlayRendered(t *testing.T) {
	m := NewModel(nil, "v1", "E01")
	m.Width = 100
	m.Height = 30
	m.SelectedRelease = "v1"
	m.AllDefects = []data.Defect{
		{Release: "v1", ID: "D001-crash", Title: "App crash", Severity: data.SeverityCritical, Status: data.DefectOpen},
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
