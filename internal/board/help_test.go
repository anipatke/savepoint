package board

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestRenderHelp_containsTitle(t *testing.T) {
	got := RenderHelp(60)
	if !strings.Contains(got, "KEYBOARD SHORTCUTS") {
		t.Error("RenderHelp missing title")
	}
}

func TestRenderHelp_containsShortcuts(t *testing.T) {
	got := RenderHelp(60)
	for _, shortcut := range []string{
		"h / left",
		"l / right",
		"enter",
		"e",
		"r",
		"p",
		"up / k",
		"down / j",
		"?",
		"esc / q",
		"q / ctrl+c",
	} {
		if !strings.Contains(got, shortcut) {
			t.Errorf("RenderHelp missing shortcut %q", shortcut)
		}
	}
}

func TestRenderHelp_containsCloseHint(t *testing.T) {
	got := RenderHelp(60)
	if !strings.Contains(got, "esc/q:close") {
		t.Error("RenderHelp missing close hint")
	}
}

func TestUpdate_questionMarkOpensHelpOverlay(t *testing.T) {
	m := NewModel(nil, "v1", "E04")
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")})
	updated := requireModel(t, got)
	if updated.Overlay != OverlayHelp {
		t.Errorf("Overlay = %q, want %q", updated.Overlay, OverlayHelp)
	}
}

func TestUpdate_helpOverlayEscCloses(t *testing.T) {
	m := NewModel(nil, "v1", "E04")
	m.Overlay = OverlayHelp
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	updated := requireModel(t, got)
	if updated.Overlay != OverlayNone {
		t.Errorf("Overlay = %q after esc, want none", updated.Overlay)
	}
}

func TestUpdate_helpOverlayQCloses(t *testing.T) {
	m := NewModel(nil, "v1", "E04")
	m.Overlay = OverlayHelp
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	updated := requireModel(t, got)
	if updated.Overlay != OverlayNone {
		t.Errorf("Overlay = %q after q, want none", updated.Overlay)
	}
}

func TestView_helpOverlayRendered(t *testing.T) {
	m := NewModel(nil, "v1", "E04")
	m.Width = 100
	m.Height = 30
	m.Overlay = OverlayHelp
	got := m.View()
	if !strings.Contains(got, "KEYBOARD SHORTCUTS") {
		t.Error("View() with OverlayHelp missing help header")
	}
	if !strings.Contains(got, "q / ctrl+c") {
		t.Error("View() with OverlayHelp missing quit shortcut")
	}
}
