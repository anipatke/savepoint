package board

import (
	"regexp"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/opencode/savepoint/internal/data"
)

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;?]*[ -/]*[@-~]`)

func plainTerminal(s string) string {
	return ansiPattern.ReplaceAllString(s, "")
}

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
	got := RenderDetail(sampleTask(), 60, nil, 0, 0)
	if !strings.Contains(got, "E04/T001") {
		t.Error("RenderDetail missing task ID")
	}
}

func TestRenderDetail_containsTitle(t *testing.T) {
	got := RenderDetail(sampleTask(), 60, nil, 0, 0)
	if !strings.Contains(got, "My Task") {
		t.Error("RenderDetail missing task title")
	}
}

func TestRenderDetail_containsEpic(t *testing.T) {
	got := RenderDetail(sampleTask(), 60, nil, 0, 0)
	if !strings.Contains(got, "E04-board-components") {
		t.Error("RenderDetail missing epic")
	}
}

func TestRenderDetail_containsRelease(t *testing.T) {
	got := RenderDetail(sampleTask(), 60, nil, 0, 0)
	if !strings.Contains(got, "v1") {
		t.Error("RenderDetail missing release")
	}
}

func TestRenderDetail_containsStatus(t *testing.T) {
	got := RenderDetail(sampleTask(), 60, nil, 0, 0)
	if !strings.Contains(got, "in_progress") {
		t.Error("RenderDetail missing status")
	}
}

func TestRenderDetail_containsStage(t *testing.T) {
	got := RenderDetail(sampleTask(), 60, nil, 0, 0)
	if !strings.Contains(got, "build") {
		t.Error("RenderDetail missing stage")
	}
}

func TestRenderDetail_containsEscHint(t *testing.T) {
	got := RenderDetail(sampleTask(), 60, nil, 0, 0)
	if !strings.Contains(got, "esc") {
		t.Error("RenderDetail missing esc:close hint")
	}
}

func TestRenderDetail_containsDescription(t *testing.T) {
	tk := sampleTask()
	tk.Description = "some description text"
	got := RenderDetail(tk, 60, nil, 0, 0)
	if !strings.Contains(got, "some description text") {
		t.Error("RenderDetail missing description text")
	}
}

func TestRenderDetail_noDescriptionSectionWhenEmpty(t *testing.T) {
	got := RenderDetail(sampleTask(), 60, nil, 0, 0)
	if strings.Contains(got, "Description:") {
		t.Error("RenderDetail should not show Description section when empty")
	}
}

func TestRenderDetail_containsChecklist(t *testing.T) {
	tk := sampleTask()
	tk.Checklist = []data.CheckItem{{Text: "first implementation item"}, {Text: "second implementation item", Done: true}}
	got := RenderDetail(tk, 60, nil, 0, 0)
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

func TestRenderDetail_checklistSingleSentenceGetsOneCheckbox(t *testing.T) {
	tk := sampleTask()
	tk.Checklist = []data.CheckItem{{Text: "single sentence task"}}

	got := plainTerminal(RenderDetail(tk, 60, nil, 0, 0))

	if count := strings.Count(got, "[ ]"); count != 1 {
		t.Fatalf("RenderDetail checkbox count = %d, want 1\n%s", count, got)
	}
	if strings.Contains(got, "[x]") {
		t.Fatal("RenderDetail should not render checked marker for unchecked item")
	}
}

func TestRenderDetail_checklistMultiSentenceGetsOneCheckboxPerSentence(t *testing.T) {
	tk := sampleTask()
	tk.Checklist = []data.CheckItem{{Text: "First sentence. Second sentence! Third sentence?"}}

	got := plainTerminal(RenderDetail(tk, 60, nil, 0, 0))

	if count := strings.Count(got, "[ ]"); count != 3 {
		t.Fatalf("RenderDetail checkbox count = %d, want 3\n%s", count, got)
	}
	for _, want := range []string{"[ ] First sentence.", "[ ] Second sentence!", "[ ] Third sentence?"} {
		if !strings.Contains(got, want) {
			t.Fatalf("RenderDetail missing sentence checkbox line %q\n%s", want, got)
		}
	}
}

func TestRenderDetail_checklistHardWrappedSentenceDoesNotDuplicateCheckbox(t *testing.T) {
	tk := sampleTask()
	tk.Checklist = []data.CheckItem{{
		Text: "This sentence is intentionally long enough to wrap inside a narrow detail overlay while remaining one semantic sentence.",
	}}

	got := plainTerminal(RenderDetail(tk, 34, nil, 0, 0))

	if count := strings.Count(got, "[ ]"); count != 1 {
		t.Fatalf("RenderDetail checkbox count = %d, want 1\n%s", count, got)
	}
	if !strings.Contains(got, "    intentionally long enough") {
		t.Fatalf("RenderDetail continuation line should align under checkbox text\n%s", got)
	}
}

func TestRenderDetail_checklistCheckedSentenceUsesCheckedMarker(t *testing.T) {
	tk := sampleTask()
	tk.Checklist = []data.CheckItem{{Text: "already done. still done.", Done: true}}

	got := plainTerminal(RenderDetail(tk, 60, nil, 0, 0))

	if count := strings.Count(got, "[x]"); count != 2 {
		t.Fatalf("RenderDetail checked checkbox count = %d, want 2\n%s", count, got)
	}
	if strings.Contains(got, "[ ]") {
		t.Fatal("RenderDetail should not render unchecked marker for checked item")
	}
}

func TestRenderDetail_wrapsLongDescription(t *testing.T) {
	tk := sampleTask()
	tk.Description = "alpha beta gamma delta epsilon zeta eta theta iota kappa lambda"
	got := RenderDetail(tk, 30, nil, 0, 0)
	if strings.Contains(got, tk.Description) {
		t.Error("RenderDetail should wrap long description text")
	}
	if !strings.Contains(got, "alpha beta") || !strings.Contains(got, "lambda") {
		t.Error("RenderDetail should preserve wrapped description words")
	}
}

func TestRenderDetail_noAcceptanceSectionWhenEmpty(t *testing.T) {
	got := RenderDetail(sampleTask(), 60, nil, 0, 0)
	if strings.Contains(got, "Acceptance Criteria:") {
		t.Error("RenderDetail should not show Acceptance section when empty")
	}
}

func TestStageLabel_build(t *testing.T) {
	if got := stageLabel(data.StageBuild); got != "build" {
		t.Errorf("stageLabel(StageBuild) = %q, want %q", got, "build")
	}
}

func TestStageLabel_test(t *testing.T) {
	if got := stageLabel(data.StageTest); got != "test" {
		t.Errorf("stageLabel(StageTest) = %q, want %q", got, "test")
	}
}

func TestStageLabel_audit(t *testing.T) {
	if got := stageLabel(data.StageAudit); got != "audit" {
		t.Errorf("stageLabel(StageAudit) = %q, want %q", got, "audit")
	}
}

func TestStageLabel_default(t *testing.T) {
	if got := stageLabel(""); got != "build" {
		t.Errorf("stageLabel(%q) = %q, want %q", "", got, "build")
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

func TestRenderDetail_routerPriorityLabel(t *testing.T) {
	task := sampleTask()
	router := &data.RouterState{Release: task.Release, Epic: task.Epic, Task: task.ID}
	got := RenderDetail(task, 60, router, 0, 0)
	if !strings.Contains(got, "(router priority)") {
		t.Error("RenderDetail missing router priority label for matching task")
	}
}

func TestRenderDetail_noRouterPriorityLabelWhenNoMatch(t *testing.T) {
	task := sampleTask()
	router := &data.RouterState{Release: task.Release, Epic: task.Epic, Task: "other-id"}
	got := RenderDetail(task, 60, router, 0, 0)
	if strings.Contains(got, "(router priority)") {
		t.Error("RenderDetail should not show router priority label for non-matching task")
	}
}

func TestRenderDetail_viewportShowsScrollIndicators(t *testing.T) {
	task := sampleTask()
	task.Description = "alpha beta gamma delta epsilon zeta eta theta iota kappa lambda mu nu xi omicron"

	got := RenderDetail(task, 32, nil, 8, 2)

	if !strings.Contains(got, "↑ 2 above") {
		t.Error("RenderDetail missing above indicator")
	}
	if !strings.Contains(got, "↓") || !strings.Contains(got, "more") {
		t.Error("RenderDetail missing more indicator")
	}
	if strings.Contains(got, "ID:") {
		t.Error("RenderDetail should not render body lines above viewport")
	}
}

func TestUpdate_detailOverlayScrollsWithJK(t *testing.T) {
	m := NewModel([]data.Task{sampleTask()}, "v1", "E04-board-components")
	m.Overlay = OverlayDetail

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	updated := requireModel(t, got)
	if updated.DetailOffset != 1 {
		t.Errorf("DetailOffset after j = %d, want 1", updated.DetailOffset)
	}

	got, _ = updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	updated = requireModel(t, got)
	if updated.DetailOffset != 0 {
		t.Errorf("DetailOffset after k = %d, want 0", updated.DetailOffset)
	}
}

func TestSplitChecklistSentences_abbreviationEg(t *testing.T) {
	got := splitChecklistSentences("Use e.g. a widget. Done.")
	if len(got) != 2 {
		t.Fatalf("splitChecklistSentences with e.g. = %d sentences, want 2: %v", len(got), got)
	}
	if got[0] != "Use e.g. a widget." {
		t.Errorf("sentence[0] = %q, want %q", got[0], "Use e.g. a widget.")
	}
}

func TestSplitChecklistSentences_abbreviationIe(t *testing.T) {
	got := splitChecklistSentences("Call i.e. the function. Done.")
	if len(got) != 2 {
		t.Fatalf("splitChecklistSentences with i.e. = %d sentences, want 2: %v", len(got), got)
	}
}

func TestSplitChecklistSentences_abbreviationDr(t *testing.T) {
	got := splitChecklistSentences("Dr. Smith approved it. Done.")
	if len(got) != 2 {
		t.Fatalf("splitChecklistSentences with Dr. = %d sentences, want 2: %v", len(got), got)
	}
	if got[0] != "Dr. Smith approved it." {
		t.Errorf("sentence[0] = %q, want %q", got[0], "Dr. Smith approved it.")
	}
}

func TestSplitChecklistSentences_abbreviationEtc(t *testing.T) {
	got := splitChecklistSentences("Add widgets, buttons, etc. to the panel. Done.")
	if len(got) != 2 {
		t.Fatalf("splitChecklistSentences with etc. = %d sentences, want 2: %v", len(got), got)
	}
}

func TestSplitChecklistSentences_abbreviationCaseInsensitive(t *testing.T) {
	got := splitChecklistSentences("See Fig. 3 for details. Done.")
	if len(got) != 2 {
		t.Fatalf("splitChecklistSentences with Fig. = %d sentences, want 2: %v", len(got), got)
	}
}

func TestSplitChecklistSentences_normalSplitUnaffected(t *testing.T) {
	got := splitChecklistSentences("First sentence. Second sentence.")
	if len(got) != 2 {
		t.Fatalf("splitChecklistSentences normal split = %d sentences, want 2: %v", len(got), got)
	}
}

func TestRenderDetail_showsComplexityTierAndReason(t *testing.T) {
	tk := sampleTask()
	tk.ComplexityTier = data.ComplexityHigh
	tk.ComplexityReason = "handles multiple parse and write paths"
	got := RenderDetail(tk, 60, nil, 0, 0)
	if !strings.Contains(got, "Complexity:") {
		t.Error("RenderDetail missing Complexity: label")
	}
	if !strings.Contains(got, "high") {
		t.Error("RenderDetail missing complexity tier")
	}
	if !strings.Contains(got, "handles multiple") {
		t.Error("RenderDetail missing complexity reason text")
	}
}

func TestRenderDetail_complexityNoReason(t *testing.T) {
	tk := sampleTask()
	tk.ComplexityTier = data.ComplexityMedium
	got := RenderDetail(tk, 60, nil, 0, 0)
	if !strings.Contains(got, "Complexity:") {
		t.Error("RenderDetail missing Complexity: label when reason is empty")
	}
	if !strings.Contains(got, "medium") {
		t.Error("RenderDetail missing complexity tier")
	}
	if strings.Contains(got, "—") {
		t.Error("RenderDetail should not include dash separator when reason is empty")
	}
}

func TestRenderDetail_noComplexityWhenEmpty(t *testing.T) {
	got := RenderDetail(sampleTask(), 60, nil, 0, 0)
	if strings.Contains(got, "Complexity:") {
		t.Error("RenderDetail should not show Complexity row when tier is empty")
	}
}

func TestRenderDetail_complexityReasonWrapsNarrow(t *testing.T) {
	tk := sampleTask()
	tk.ComplexityTier = data.ComplexitySpike
	tk.ComplexityReason = "requires deep investigation of unknown system boundaries and significant research"
	got := plainTerminal(RenderDetail(tk, 30, nil, 0, 0))
	if !strings.Contains(got, "Complexity:") {
		t.Error("RenderDetail missing Complexity: on narrow width")
	}
	if strings.Contains(got, tk.ComplexityReason) {
		t.Error("RenderDetail should wrap long complexity reason on narrow width")
	}
	if !strings.Contains(got, "spike") {
		t.Error("RenderDetail missing tier on narrow complexity row")
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
