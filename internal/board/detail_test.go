package board

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/opencode/savepoint/internal/data"
)

func sampleTask() data.Task {
	return data.Task{
		ID:      "E04/T001",
		Title:   "My Task",
		Epic:    "E04-board-components",
		Release: "v1",
		Column:  data.ColumnInProgress,
		Stage:   data.StageBuild,
	}
}

func TestRenderDetail_containsID(t *testing.T) {
	got := RenderDetail(sampleTask(), 60)
	if !strings.Contains(got, "E04/T001") {
		t.Error("RenderDetail missing task ID")
	}
}

func TestRenderDetail_containsTitle(t *testing.T) {
	got := RenderDetail(sampleTask(), 60)
	if !strings.Contains(got, "My Task") {
		t.Error("RenderDetail missing task title")
	}
}

func TestRenderDetail_containsEpic(t *testing.T) {
	got := RenderDetail(sampleTask(), 60)
	if !strings.Contains(got, "E04-board-components") {
		t.Error("RenderDetail missing epic")
	}
}

func TestRenderDetail_containsRelease(t *testing.T) {
	got := RenderDetail(sampleTask(), 60)
	if !strings.Contains(got, "v1") {
		t.Error("RenderDetail missing release")
	}
}

func TestRenderDetail_containsStatus(t *testing.T) {
	got := RenderDetail(sampleTask(), 60)
	if !strings.Contains(got, "in_progress") {
		t.Error("RenderDetail missing status")
	}
}

func TestRenderDetail_containsPhase(t *testing.T) {
	got := RenderDetail(sampleTask(), 60)
	if !strings.Contains(got, "build") {
		t.Error("RenderDetail missing phase")
	}
}

func TestRenderDetail_containsEscHint(t *testing.T) {
	got := RenderDetail(sampleTask(), 60)
	if !strings.Contains(got, "esc") {
		t.Error("RenderDetail missing esc:close hint")
	}
}

func TestRenderDetail_containsDescription(t *testing.T) {
	tk := sampleTask()
	tk.Description = "some description text"
	got := RenderDetail(tk, 60)
	if !strings.Contains(got, "some description text") {
		t.Error("RenderDetail missing description text")
	}
}

func TestRenderDetail_noDescriptionSectionWhenEmpty(t *testing.T) {
	got := RenderDetail(sampleTask(), 60)
	if strings.Contains(got, "Description:") {
		t.Error("RenderDetail should not show Description section when empty")
	}
}

func TestRenderDetail_containsAcceptanceCriteria(t *testing.T) {
	tk := sampleTask()
	tk.Acceptance = []string{"criterion one", "criterion two"}
	got := RenderDetail(tk, 60)
	if !strings.Contains(got, "criterion one") {
		t.Error("RenderDetail missing first acceptance criterion")
	}
	if !strings.Contains(got, "criterion two") {
		t.Error("RenderDetail missing second acceptance criterion")
	}
}

func TestRenderDetail_containsChecklist(t *testing.T) {
	tk := sampleTask()
	tk.Checklist = []string{"first implementation item", "second implementation item"}
	got := RenderDetail(tk, 60)
	if !strings.Contains(got, "Implementation Plan:") {
		t.Error("RenderDetail missing implementation plan heading")
	}
	if !strings.Contains(got, "first implementation item") {
		t.Error("RenderDetail missing first checklist item")
	}
	if !strings.Contains(got, "second implementation item") {
		t.Error("RenderDetail missing second checklist item")
	}
}

func TestRenderDetail_wrapsLongDescription(t *testing.T) {
	tk := sampleTask()
	tk.Description = "alpha beta gamma delta epsilon zeta eta theta iota kappa lambda"
	got := RenderDetail(tk, 30)
	if strings.Contains(got, tk.Description) {
		t.Error("RenderDetail should wrap long description text")
	}
	if !strings.Contains(got, "alpha beta") || !strings.Contains(got, "lambda") {
		t.Error("RenderDetail should preserve wrapped description words")
	}
}

func TestRenderDetail_noAcceptanceSectionWhenEmpty(t *testing.T) {
	got := RenderDetail(sampleTask(), 60)
	if strings.Contains(got, "Acceptance Criteria:") {
		t.Error("RenderDetail should not show Acceptance section when empty")
	}
}

func TestPhaseLabel_build(t *testing.T) {
	if got := phaseLabel(data.StageBuild); got != "build" {
		t.Errorf("phaseLabel(StageBuild) = %q, want %q", got, "build")
	}
}

func TestPhaseLabel_test(t *testing.T) {
	if got := phaseLabel(data.StageTest); got != "test" {
		t.Errorf("phaseLabel(StageTest) = %q, want %q", got, "test")
	}
}

func TestPhaseLabel_audit(t *testing.T) {
	if got := phaseLabel(data.StageAudit); got != "audit" {
		t.Errorf("phaseLabel(StageAudit) = %q, want %q", got, "audit")
	}
}

func TestPhaseLabel_default(t *testing.T) {
	if got := phaseLabel(""); got != "build" {
		t.Errorf("phaseLabel(%q) = %q, want %q", "", got, "build")
	}
}

func TestUpdate_enterOpensDetailOverlay(t *testing.T) {
	tasks := []data.Task{sampleTask()}
	m := NewModel(tasks, "v1", "E04-board-components")
	m.FocusedColumn = data.ColumnInProgress
	m.FocusedTask = 0
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	updated := requireModel(t, got)
	if updated.Overlay != OverlayDetail {
		t.Errorf("Overlay = %q, want %q", updated.Overlay, OverlayDetail)
	}
}

func TestUpdate_enterNoOpWhenNoTasks(t *testing.T) {
	m := NewModel(nil, "v1", "E04-board-components")
	m.FocusedColumn = data.ColumnPlanned
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	updated := requireModel(t, got)
	if updated.Overlay != OverlayNone {
		t.Errorf("Overlay = %q, want none when column has no tasks", updated.Overlay)
	}
}

func TestUpdate_detailOverlayEscCloses(t *testing.T) {
	m := NewModel(nil, "v1", "E04-board-components")
	m.Overlay = OverlayDetail
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	updated := requireModel(t, got)
	if updated.Overlay != OverlayNone {
		t.Errorf("Overlay = %q after esc, want none", updated.Overlay)
	}
}

func TestUpdate_detailOverlayBlocksColumnNav(t *testing.T) {
	m := NewModel(nil, "v1", "E04-board-components")
	m.Overlay = OverlayDetail
	m.FocusedColumn = data.ColumnPlanned
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")})
	updated := requireModel(t, got)
	if updated.FocusedColumn != data.ColumnPlanned {
		t.Error("column nav should be blocked when detail overlay is open")
	}
}

func TestView_detailOverlayRendered(t *testing.T) {
	tasks := []data.Task{sampleTask()}
	m := NewModel(tasks, "v1", "E04-board-components")
	m.Width = 100
	m.Height = 30
	m.FocusedColumn = data.ColumnInProgress
	m.FocusedTask = 0
	m.Overlay = OverlayDetail
	got := m.View()
	if !strings.Contains(got, "TASK DETAIL") {
		t.Error("View() with OverlayDetail missing TASK DETAIL header")
	}
	if !strings.Contains(got, "E04/T001") {
		t.Error("View() with OverlayDetail missing task ID")
	}
}

func TestOverlayWidth_clampMax(t *testing.T) {
	if got := overlayWidth(120); got != 80 {
		t.Errorf("overlayWidth(120) = %d, want 80", got)
	}
}

func TestOverlayWidth_termMinus4(t *testing.T) {
	if got := overlayWidth(60); got != 56 {
		t.Errorf("overlayWidth(60) = %d, want 56", got)
	}
}

func TestOverlayWidth_clampMin(t *testing.T) {
	if got := overlayWidth(10); got != 20 {
		t.Errorf("overlayWidth(10) = %d, want 20", got)
	}
}
