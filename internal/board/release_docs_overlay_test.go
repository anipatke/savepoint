package board

import (
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/testutil"
)

func sampleReleaseDocs() []data.ReleaseDoc {
	return []data.ReleaseDoc{
		{ID: data.ReleaseDocReleasePRD, Label: "Release PRD", Path: "releases/v1.2/v1.2-PRD.md", Body: "# Release PRD\n\nRelease scope text.", Available: true},
		{ID: data.ReleaseDocOverallPRD, Label: "Overall PRD", Path: "PRD.md", Body: "# Overall PRD\n\nVision text.", Available: true},
		{ID: data.ReleaseDocOverallDesign, Label: "Overall Design", Path: "Design.md", Body: "# Overall Design\n\nDesign text.", Available: true},
	}
}

// --- renderer ---

func TestRenderReleaseDocs_headerAndSelector(t *testing.T) {
	got := RenderReleaseDocs(sampleReleaseDocs(), 0, 70, 40, 0)
	if !strings.Contains(got, "RELEASE DOCS") {
		t.Error("RenderReleaseDocs missing RELEASE DOCS header")
	}
	for _, label := range []string{"Release PRD", "Overall PRD", "Overall Design"} {
		if !strings.Contains(got, label) {
			t.Errorf("RenderReleaseDocs missing selector label %q", label)
		}
	}
	if strings.Contains(got, "DOCS [3]") || strings.Contains(got, "DETAIL [1]") {
		t.Error("RenderReleaseDocs should not render the epic Detail/Audit tab strip")
	}
}

func TestRenderReleaseDocs_selectedBody(t *testing.T) {
	got := RenderReleaseDocs(sampleReleaseDocs(), 0, 70, 40, 0)
	if !strings.Contains(got, "Release scope text") {
		t.Error("RenderReleaseDocs should render the selected (Release PRD) body")
	}
	if strings.Contains(got, "Vision text") || strings.Contains(got, "Design text") {
		t.Error("RenderReleaseDocs should not render unselected document bodies")
	}
}

func TestRenderReleaseDocs_footer(t *testing.T) {
	got := RenderReleaseDocs(sampleReleaseDocs(), 0, 70, 40, 0)
	if !strings.Contains(got, "esc:close") {
		t.Error("RenderReleaseDocs missing esc:close footer")
	}
	if !strings.Contains(got, "[/]:doc") {
		t.Error("RenderReleaseDocs missing doc-switch hint")
	}
}

func TestRenderReleaseDocs_missingDocEmptyState(t *testing.T) {
	docs := []data.ReleaseDoc{
		{ID: data.ReleaseDocReleasePRD, Label: "Release PRD", Path: "releases/v1.2/v1.2-PRD.md", Available: false},
		{ID: data.ReleaseDocOverallPRD, Label: "Overall PRD", Path: "PRD.md", Body: "x", Available: true},
	}
	got := RenderReleaseDocs(docs, 0, 70, 40, 0)
	if !strings.Contains(got, "not found") || !strings.Contains(got, "v1.2-PRD.md") {
		t.Error("RenderReleaseDocs should show a missing-doc empty state naming the release PRD path")
	}
}

// --- board key + overlay wiring ---

func TestUpdate_DOpensReleaseDocsOverlayAndLoadsForSelectedRelease(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "releases", "v1.2", "v1.2-PRD.md"), "# Release PRD\nv1.2 scope")
	testutil.WriteFile(t, filepath.Join(root, "PRD.md"), "# Overall PRD\nvision")
	testutil.WriteFile(t, filepath.Join(root, "Design.md"), "# Overall Design\narch")

	m := NewModel(nil, "v1.2", "E01")
	m.Root = root

	got, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("D")})
	updated := requireModel(t, got)
	if updated.Overlay != OverlayReleaseDocs {
		t.Fatalf("Overlay = %q, want %q", updated.Overlay, OverlayReleaseDocs)
	}
	if cmd == nil {
		t.Fatal("expected a release-docs load command")
	}

	got2, _ := updated.Update(cmd())
	updated2 := requireModel(t, got2)
	if len(updated2.ReleaseDocs) != 3 {
		t.Fatalf("ReleaseDocs len = %d, want 3", len(updated2.ReleaseDocs))
	}
	relPRD, ok := docByID(updated2.ReleaseDocs, data.ReleaseDocReleasePRD)
	if !ok || !relPRD.Available || relPRD.Body != "# Release PRD\nv1.2 scope" {
		t.Errorf("Release PRD = %+v, want the v1.2 release PRD", relPRD)
	}
}

// TestUpdate_DLowercaseStillOpensDefects guards the case-sensitivity: lowercase
// d is defects, capital D is release docs.
func TestUpdate_DLowercaseStillOpensDefects(t *testing.T) {
	m := NewModel(nil, "v1.2", "E01")
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("d")})
	if requireModel(t, got).Overlay != OverlayDefect {
		t.Error("lowercase d should open the defects overlay, not release docs")
	}
}

func TestUpdate_releaseDocsOverlayClosesOnEsc(t *testing.T) {
	m := NewModel(nil, "v1.2", "E01")
	m.Overlay = OverlayReleaseDocs

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if requireModel(t, got).Overlay != OverlayNone {
		t.Error("esc should close the release docs overlay")
	}
}

func TestUpdate_releaseDocSelectionAndScroll(t *testing.T) {
	m := NewModel(nil, "v1.2", "E01")
	m.Overlay = OverlayReleaseDocs
	m.ReleaseDocs = sampleReleaseDocs()
	m.ReleaseDocOffsets = map[data.ReleaseDocID]int{}

	// Scroll the Release PRD twice, then switch to Overall PRD with ].
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	got, _ = requireModel(t, got).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	got, _ = requireModel(t, got).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("]")})
	updated := requireModel(t, got)

	if updated.ReleaseDocIndex != 1 {
		t.Fatalf("ReleaseDocIndex after ] = %d, want 1", updated.ReleaseDocIndex)
	}
	if updated.ReleaseDocOffsets[data.ReleaseDocReleasePRD] != 2 {
		t.Errorf("Release PRD offset = %d, want 2 (preserved across switch)", updated.ReleaseDocOffsets[data.ReleaseDocReleasePRD])
	}
	if updated.ReleaseDocOffsets[data.ReleaseDocOverallPRD] != 0 {
		t.Errorf("Overall PRD offset = %d, want 0", updated.ReleaseDocOffsets[data.ReleaseDocOverallPRD])
	}

	// Over-scroll up on Overall PRD clamps at the top.
	got, _ = updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	updated = requireModel(t, got)
	if updated.ReleaseDocOffsets[data.ReleaseDocOverallPRD] != 0 {
		t.Errorf("Overall PRD offset after over-scroll up = %d, want 0", updated.ReleaseDocOffsets[data.ReleaseDocOverallPRD])
	}
}

func TestView_releaseDocsOverlayRendered(t *testing.T) {
	m := NewModel(nil, "v1.2", "E01")
	m.Width = 120
	m.Height = 30
	m.Overlay = OverlayReleaseDocs
	m.ReleaseDocs = sampleReleaseDocs()

	got := m.View()
	if !strings.Contains(got, "RELEASE DOCS") {
		t.Error("View() with OverlayReleaseDocs should render the RELEASE DOCS overlay")
	}
}
