package board

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/opencode/savepoint/internal/data"
)

func TestRenderReleaseDropdown_showsReleases(t *testing.T) {
	releases := []string{"v1", "v2", "v3"}
	out := RenderReleaseDropdown(releases, 0, 40)
	for _, r := range releases {
		if !strings.Contains(out, r) {
			t.Errorf("output missing release %q", r)
		}
	}
}

func TestRenderReleaseDropdown_marksCurrentCursor(t *testing.T) {
	releases := []string{"v1", "v2"}
	out := RenderReleaseDropdown(releases, 1, 40)
	if !strings.Contains(out, releaseActiveMarker) {
		t.Error("output missing active marker")
	}
}

func TestRenderReleaseDropdown_emptyList(t *testing.T) {
	out := RenderReleaseDropdown(nil, 0, 40)
	if !strings.Contains(out, "(none)") {
		t.Error("expected '(none)' for empty releases")
	}
}

func TestRenderReleaseDropdown_hintsPresent(t *testing.T) {
	out := RenderReleaseDropdown([]string{"v1"}, 0, 40)
	if !strings.Contains(out, "esc:cancel") {
		t.Error("expected hint text in dropdown")
	}
}

func TestReleaseIndex_found(t *testing.T) {
	releases := []string{"v1", "v2", "v3"}
	if got := releaseIndex(releases, "v2"); got != 1 {
		t.Errorf("releaseIndex = %d, want 1", got)
	}
}

func TestReleaseIndex_notFound(t *testing.T) {
	releases := []string{"v1", "v2"}
	if got := releaseIndex(releases, "v9"); got != 0 {
		t.Errorf("releaseIndex = %d, want 0", got)
	}
}

func TestUpdate_rOpensReleaseOverlay(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Releases = []string{"v1", "v2"}
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")})
	updated := requireModel(t, got)
	if updated.Overlay != OverlayRelease {
		t.Errorf("Overlay = %q, want %q", updated.Overlay, OverlayRelease)
	}
}

func TestUpdate_rSetsCursorToCurrentRelease(t *testing.T) {
	m := NewModel(nil, "v2", "E03")
	m.Releases = []string{"v1", "v2", "v3"}
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")})
	updated := requireModel(t, got)
	if updated.ReleaseCursor != 1 {
		t.Errorf("ReleaseCursor = %d, want 1", updated.ReleaseCursor)
	}
}

func TestUpdate_escClosesReleaseOverlay(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Overlay = OverlayRelease
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	updated := requireModel(t, got)
	if updated.Overlay != OverlayNone {
		t.Errorf("Overlay = %q, want empty", updated.Overlay)
	}
}

func TestUpdate_releaseNavDown(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Releases = []string{"v1", "v2"}
	m.Overlay = OverlayRelease
	m.ReleaseCursor = 0
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	updated := requireModel(t, got)
	if updated.ReleaseCursor != 1 {
		t.Errorf("ReleaseCursor = %d, want 1", updated.ReleaseCursor)
	}
}

func TestUpdate_releaseNavUp(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Releases = []string{"v1", "v2"}
	m.Overlay = OverlayRelease
	m.ReleaseCursor = 1
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	updated := requireModel(t, got)
	if updated.ReleaseCursor != 0 {
		t.Errorf("ReleaseCursor = %d, want 0", updated.ReleaseCursor)
	}
}

func TestUpdate_releaseEnterSelects(t *testing.T) {
	tasks := []data.Task{
		{ID: "T1", Epic: "E03", Release: "v1", Column: data.ColumnPlanned},
		{ID: "T2", Epic: "E03", Release: "v2", Column: data.ColumnPlanned},
	}
	m := NewModel(tasks, "v1", "E03")
	m.Releases = []string{"v1", "v2"}
	m.Overlay = OverlayRelease
	m.ReleaseCursor = 1
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	updated := requireModel(t, got)
	if updated.SelectedRelease != "v2" {
		t.Errorf("SelectedRelease = %q, want %q", updated.SelectedRelease, "v2")
	}
	if updated.Overlay != OverlayNone {
		t.Errorf("Overlay = %q, want empty after selection", updated.Overlay)
	}
	if got := len(updated.Tasks[data.ColumnPlanned]); got != 1 {
		t.Errorf("planned task count = %d, want 1 after release selection", got)
	}
	if updated.Tasks[data.ColumnPlanned][0].ID != "T2" {
		t.Errorf("visible task = %q, want T2", updated.Tasks[data.ColumnPlanned][0].ID)
	}
}

func TestUpdate_releaseEnterRefreshesEpics(t *testing.T) {
	tasks := []data.Task{
		{ID: "T1", Epic: "E01", Release: "v1", Column: data.ColumnPlanned},
		{ID: "T2", Epic: "E03", Release: "v2", Column: data.ColumnPlanned},
	}
	m := NewModel(tasks, "v1", "E01")
	m.Releases = []string{"v1", "v2"}
	m.ReleaseEpics = map[string][]string{
		"v1": []string{"E01"},
		"v2": []string{"E03"},
	}
	m.Overlay = OverlayRelease
	m.ReleaseCursor = 1

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	updated := requireModel(t, got)

	if updated.SelectedEpic != "E03" {
		t.Errorf("SelectedEpic = %q, want E03", updated.SelectedEpic)
	}
	if len(updated.Epics) != 1 || updated.Epics[0] != "E03" {
		t.Errorf("Epics = %v, want [E03]", updated.Epics)
	}
	if got := len(updated.Tasks[data.ColumnPlanned]); got != 1 {
		t.Fatalf("planned task count = %d, want 1", got)
	}
	if updated.Tasks[data.ColumnPlanned][0].ID != "T2" {
		t.Errorf("visible task = %q, want T2", updated.Tasks[data.ColumnPlanned][0].ID)
	}
}

func TestView_releaseDropdownKeepsBoardBehind(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Width = 100
	m.Height = 24
	m.Overlay = OverlayRelease
	m.Releases = []string{"v1", "v2"}
	got := m.View()
	if !strings.Contains(got, "S A V E P O I N T") {
		t.Error("View() with OverlayRelease should keep board visible behind overlay")
	}
}
