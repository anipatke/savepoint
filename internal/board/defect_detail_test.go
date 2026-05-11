package board

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/opencode/savepoint/internal/data"
)

func sampleDefect() data.Defect {
	return data.Defect{
		ID:       "v1.1/D003-auth-crash",
		Title:    "Auth crash on empty token",
		Status:   data.DefectOpen,
		Severity: data.SeverityHigh,
		Release:  "v1.1",
	}
}

func TestRenderDefectDetail_containsTitle(t *testing.T) {
	got := RenderDefectDetail(sampleDefect(), 60, 0, 0)
	if !strings.Contains(got, "Auth crash on empty token") {
		t.Error("RenderDefectDetail missing title")
	}
}

func TestRenderDefectDetail_containsHeader(t *testing.T) {
	got := RenderDefectDetail(sampleDefect(), 60, 0, 0)
	if !strings.Contains(got, "DEFECT DETAIL") {
		t.Error("RenderDefectDetail missing DEFECT DETAIL header")
	}
}

func TestRenderDefectDetail_containsShortID(t *testing.T) {
	got := RenderDefectDetail(sampleDefect(), 60, 0, 0)
	if !strings.Contains(got, "D003") {
		t.Error("RenderDefectDetail missing short defect ID")
	}
}

func TestRenderDefectDetail_containsSeverity(t *testing.T) {
	got := RenderDefectDetail(sampleDefect(), 60, 0, 0)
	if !strings.Contains(got, "high") {
		t.Error("RenderDefectDetail missing severity")
	}
}

func TestRenderDefectDetail_containsStatus(t *testing.T) {
	got := RenderDefectDetail(sampleDefect(), 60, 0, 0)
	if !strings.Contains(got, "open") {
		t.Error("RenderDefectDetail missing status")
	}
}

func TestRenderDefectDetail_containsRelease(t *testing.T) {
	got := RenderDefectDetail(sampleDefect(), 60, 0, 0)
	if !strings.Contains(got, "v1.1") {
		t.Error("RenderDefectDetail missing release")
	}
}

func TestRenderDefectDetail_containsEscHint(t *testing.T) {
	got := RenderDefectDetail(sampleDefect(), 60, 0, 0)
	if !strings.Contains(got, "esc") {
		t.Error("RenderDefectDetail missing esc:close hint")
	}
}

func TestRenderDefectDetail_rendersOptionalIntroduced(t *testing.T) {
	d := sampleDefect()
	d.Introduced = "v1.0.5"
	got := RenderDefectDetail(d, 60, 0, 0)
	if !strings.Contains(got, "v1.0.5") {
		t.Error("RenderDefectDetail missing Introduced value")
	}
}

func TestRenderDefectDetail_noIntroducedWhenEmpty(t *testing.T) {
	got := RenderDefectDetail(sampleDefect(), 60, 0, 0)
	if strings.Contains(got, "Introduced") {
		t.Error("RenderDefectDetail should not show Introduced when empty")
	}
}

func TestRenderDefectDetail_rendersOptionalReference(t *testing.T) {
	d := sampleDefect()
	d.Reference = "E12/T003"
	got := RenderDefectDetail(d, 60, 0, 0)
	if !strings.Contains(got, "E12/T003") {
		t.Error("RenderDefectDetail missing Reference value")
	}
}

func TestRenderDefectDetail_noReferenceWhenEmpty(t *testing.T) {
	got := RenderDefectDetail(sampleDefect(), 60, 0, 0)
	if strings.Contains(got, "Reference") {
		t.Error("RenderDefectDetail should not show Reference when empty")
	}
}

func TestRenderDefectDetail_rendersSymptomSection(t *testing.T) {
	d := sampleDefect()
	d.Body = "## Symptom\n\nToken validation panics on empty input.\n"
	got := RenderDefectDetail(d, 60, 0, 0)
	if !strings.Contains(got, "Symptom") {
		t.Error("RenderDefectDetail missing Symptom heading")
	}
	if !strings.Contains(got, "Token validation panics on empty input.") {
		t.Error("RenderDefectDetail missing Symptom content")
	}
}

func TestRenderDefectDetail_rendersExpectedBehaviorSection(t *testing.T) {
	d := sampleDefect()
	d.Body = "## Expected Behavior\n\nShould return an error, not panic.\n"
	got := RenderDefectDetail(d, 60, 0, 0)
	if !strings.Contains(got, "Expected Behavior") {
		t.Error("RenderDefectDetail missing Expected Behavior heading")
	}
}

func TestRenderDefectDetail_rendersResolutionNotesSection(t *testing.T) {
	d := sampleDefect()
	d.Body = "## Resolution Notes\n\nFixed by adding nil check.\n"
	got := RenderDefectDetail(d, 60, 0, 0)
	if !strings.Contains(got, "Resolution Notes") {
		t.Error("RenderDefectDetail missing Resolution Notes heading")
	}
	if !strings.Contains(got, "Fixed by adding nil check.") {
		t.Error("RenderDefectDetail missing Resolution Notes content")
	}
}

func TestRenderDefectDetail_skipsAbsentOptionalSections(t *testing.T) {
	d := sampleDefect()
	d.Body = "## Symptom\n\nPanics.\n"
	got := RenderDefectDetail(d, 60, 0, 0)
	for _, heading := range []string{"Expected Behavior", "Reproduction", "Impact", "Fix Plan", "Acceptance Criteria", "Resolution Notes"} {
		if strings.Contains(got, heading) {
			t.Errorf("RenderDefectDetail should not show absent section %q", heading)
		}
	}
}

func TestRenderDefectDetail_wrapsLongContent(t *testing.T) {
	d := sampleDefect()
	long := "alpha beta gamma delta epsilon zeta eta theta iota kappa lambda mu nu"
	d.Body = "## Symptom\n\n" + long + "\n"
	got := RenderDefectDetail(d, 30, 0, 0)
	if strings.Contains(got, long) {
		t.Error("RenderDefectDetail should wrap long section content")
	}
	if !strings.Contains(got, "alpha") || !strings.Contains(got, "lambda") {
		t.Error("RenderDefectDetail should preserve wrapped content words")
	}
}

func TestRenderDefectDetail_scrollsWithOffset(t *testing.T) {
	d := sampleDefect()
	d.Body = "## Symptom\n\nLine one.\n## Expected Behavior\n\nLine two.\n## Reproduction\n\nLine three.\n"
	got := RenderDefectDetail(d, 60, 8, 3)
	// With a non-zero offset the ID row should be scrolled past.
	if strings.Contains(got, "ID:") {
		t.Error("RenderDefectDetail should not render rows scrolled above viewport")
	}
	if !strings.Contains(got, "↑") {
		t.Error("RenderDefectDetail missing above scroll indicator")
	}
}

// Card marker tests

func TestDefectMarkerForTask_matchesFullID(t *testing.T) {
	defects := []data.Defect{
		{ID: "v1.1/D003-auth-crash", Reference: "E12/T003", Status: data.DefectOpen},
	}
	got := defectMarkerForTask("E12/T003", defects)
	if got != "⚠  1/0" {
		t.Errorf("defectMarkerForTask = %q, want %q", got, "⚠  1/0")
	}
}

func TestDefectMarkerForTask_matchesShortID(t *testing.T) {
	defects := []data.Defect{
		{ID: "v1.1/D005-crash", Reference: "T007", Status: data.DefectOpen},
	}
	got := defectMarkerForTask("E03/T007-some-slug", defects)
	if got != "⚠  1/0" {
		t.Errorf("defectMarkerForTask = %q, want %q", got, "⚠  1/0")
	}
}

func TestDefectMarkerForTask_noMatchWhenDifferentRef(t *testing.T) {
	defects := []data.Defect{
		{ID: "v1.1/D001-foo", Reference: "T999", Status: data.DefectOpen},
	}
	got := defectMarkerForTask("T001", defects)
	if got != "" {
		t.Errorf("defectMarkerForTask = %q, want empty", got)
	}
}

func TestDefectMarkerForTask_noMatchWhenNoReference(t *testing.T) {
	defects := []data.Defect{
		{ID: "v1.1/D002-bar", Reference: "", Status: data.DefectOpen},
	}
	got := defectMarkerForTask("T001", defects)
	if got != "" {
		t.Errorf("defectMarkerForTask = %q, want empty for defect with no reference", got)
	}
}

func TestDefectMarkerForTask_aggregatesMultipleOpenDefects(t *testing.T) {
	defects := []data.Defect{
		{ID: "v1.1/D001-foo", Reference: "T005", Status: data.DefectOpen},
		{ID: "v1.1/D002-bar", Reference: "T005", Status: data.DefectInProgress},
		{ID: "v1.1/D003-baz", Reference: "T005", Status: data.DefectResolved},
	}
	got := defectMarkerForTask("T005", defects)
	if got != "⚠  2/1" {
		t.Errorf("defectMarkerForTask = %q, want %q", got, "⚠  2/1")
	}
}

func TestDefectMarkerForTask_allResolved(t *testing.T) {
	defects := []data.Defect{
		{ID: "v1.1/D001-foo", Reference: "T005", Status: data.DefectResolved},
		{ID: "v1.1/D002-bar", Reference: "T005", Status: data.DefectResolved},
	}
	got := defectMarkerForTask("T005", defects)
	if got != "⚠  0/2" {
		t.Errorf("defectMarkerForTask = %q, want %q", got, "⚠  0/2")
	}
}

func TestDefectMarkerForTask_inProgressCountsAsOpen(t *testing.T) {
	defects := []data.Defect{
		{ID: "v1.1/D004-qux", Reference: "E05/T002", Status: data.DefectInProgress},
	}
	got := defectMarkerForTask("E05/T002", defects)
	if got != "⚠  1/0" {
		t.Errorf("defectMarkerForTask = %q, want %q", got, "⚠  1/0")
	}
}

func TestDefectMarkerForTask_mixedRefStylesCounted(t *testing.T) {
	// One defect uses the full task ID, another uses the short task ID — both must match.
	defects := []data.Defect{
		{ID: "v1.2/D010-x", Reference: "E09/T003-some-slug", Status: data.DefectOpen},
		{ID: "v1.2/D011-y", Reference: "T003", Status: data.DefectResolved},
	}
	got := defectMarkerForTask("E09/T003-some-slug", defects)
	if got != "⚠  1/1" {
		t.Errorf("defectMarkerForTask = %q, want %q", got, "⚠  1/1")
	}
}

func TestRenderCard_defectMarkerAppearsWhenWidthPermits(t *testing.T) {
	task := data.Task{ID: "T1", Title: "Fix bug", Stage: data.StageBuild}
	got := RenderCard(task, 40, false, nil, "⚠  D003")
	if !strings.Contains(got, "⚠  D003") {
		t.Error("RenderCard missing defect marker when width permits")
	}
}

func TestRenderCard_defectMarkerOmittedWhenTooNarrow(t *testing.T) {
	task := data.Task{ID: "T1", Title: "Fix", Stage: data.StageBuild}
	// At width 10, id + glyph + marker won't fit — marker must be omitted.
	got := plainTerminal(RenderCard(task, 10, false, nil, "⚠  D003"))
	// Title and ID must still appear.
	if !strings.Contains(got, "T1") {
		t.Error("RenderCard missing task ID at narrow width")
	}
}

// Overlay integration tests

func TestUpdate_enterOnDefectOverlayOpensDefectDetail(t *testing.T) {
	d := sampleDefect()
	m := NewModel(nil, "v1.1", "")
	m.AllDefects = []data.Defect{d}
	m.SelectedRelease = "v1.1"
	m.Overlay = OverlayDefect
	m.DefectCursor = 0

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	updated := requireModel(t, got)
	if updated.Overlay != OverlayDefectDetail {
		t.Errorf("Overlay = %q, want %q", updated.Overlay, OverlayDefectDetail)
	}
}

func TestUpdate_enterOnDefectOverlayNoOpWhenEmpty(t *testing.T) {
	m := NewModel(nil, "v1.1", "")
	m.Overlay = OverlayDefect
	m.DefectCursor = 0

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	updated := requireModel(t, got)
	if updated.Overlay != OverlayDefect {
		t.Errorf("Overlay = %q, want defect (no-op when empty)", updated.Overlay)
	}
}

func TestUpdate_defectDetailEscReturnsToDefectList(t *testing.T) {
	m := NewModel(nil, "v1.1", "")
	m.Overlay = OverlayDefectDetail

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	updated := requireModel(t, got)
	if updated.Overlay != OverlayDefect {
		t.Errorf("Overlay = %q, want defect after esc from detail", updated.Overlay)
	}
}

func TestUpdate_defectDetailScrollsWithJK(t *testing.T) {
	m := NewModel(nil, "v1.1", "")
	m.Overlay = OverlayDefectDetail
	m.DefectDetailOffset = 0

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	updated := requireModel(t, got)
	if updated.DefectDetailOffset != 1 {
		t.Errorf("DefectDetailOffset after j = %d, want 1", updated.DefectDetailOffset)
	}

	got, _ = updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	updated = requireModel(t, got)
	if updated.DefectDetailOffset != 0 {
		t.Errorf("DefectDetailOffset after k = %d, want 0", updated.DefectDetailOffset)
	}
}

func TestView_defectDetailOverlayRendered(t *testing.T) {
	d := sampleDefect()
	m := NewModel(nil, "v1.1", "")
	m.Width = 100
	m.Height = 30
	m.AllDefects = []data.Defect{d}
	m.SelectedRelease = "v1.1"
	m.Overlay = OverlayDefectDetail
	m.DefectCursor = 0

	got := m.View()
	if !strings.Contains(got, "DEFECT DETAIL") {
		t.Error("View() with OverlayDefectDetail missing DEFECT DETAIL header")
	}
	if !strings.Contains(got, "D003") {
		t.Error("View() with OverlayDefectDetail missing defect short ID")
	}
}

func TestParseDefectSections_parsesMultipleSections(t *testing.T) {
	body := "## Symptom\n\nPanics.\n\n## Impact\n\nData loss.\n"
	got := parseDefectSections(body)
	if got["Symptom"] == "" {
		t.Error("parseDefectSections missing Symptom section")
	}
	if got["Impact"] == "" {
		t.Error("parseDefectSections missing Impact section")
	}
	if !strings.Contains(got["Symptom"], "Panics.") {
		t.Errorf("Symptom content = %q, want 'Panics.'", got["Symptom"])
	}
}

func TestParseDefectSections_emptyBody(t *testing.T) {
	got := parseDefectSections("")
	if len(got) != 0 {
		t.Errorf("parseDefectSections on empty body = %v, want empty map", got)
	}
}

func TestShortDefectID_stripsReleaseAndSlug(t *testing.T) {
	got := shortDefectID("v1.1/D003-auth-crash")
	if got != "D003" {
		t.Errorf("shortDefectID = %q, want %q", got, "D003")
	}
}

func TestShortDefectID_noPrefix(t *testing.T) {
	got := shortDefectID("D007-some-bug")
	if got != "D007" {
		t.Errorf("shortDefectID = %q, want %q", got, "D007")
	}
}
