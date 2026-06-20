package board

import (
	"fmt"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/opencode/savepoint/internal/data"
)

func TestView_rendersWithoutPanic(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Width = 120
	m.Height = 40
	got := m.View()
	if got == "" {
		t.Error("View() returned empty string")
	}
}

func TestView_containsColumnTitles(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Width = 120
	got := m.View()

	for _, title := range []string{"PLANNED", "IN PROGRESS", "DONE"} {
		if !strings.Contains(got, title) {
			t.Errorf("View() missing column title %q", title)
		}
	}
}

func TestView_containsHeader(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Width = 120
	got := m.View()
	if !strings.Contains(got, "S A V E P O I N T") {
		t.Error("View() missing spaced header text")
	}
	if !strings.Contains(got, "▣") {
		t.Error("View() missing header icon")
	}
}

func TestView_containsDivider(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Width = 120
	got := m.View()
	if !strings.Contains(got, strings.Repeat("─", 120)) {
		t.Error("View() missing full-width horizontal divider")
	}
}

func TestView_containsFooterPhases(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Width = 120
	got := m.View()
	for _, phase := range []string{"PLAN", "BUILD", "AUDIT"} {
		if !strings.Contains(got, phase) {
			t.Errorf("View() missing footer phase %q", phase)
		}
	}
}

func TestView_containsFooterHints(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	footer := m.renderFooter(80)

	if !strings.Contains(footer, "←→:nav  p:priority  ctrl+r:refresh  R:release  d:defects  D:docs  ?:help  q:quit") {
		t.Fatal("renderFooter() missing navigation hints")
	}

	lines := strings.Split(footer, "\n")
	if len(lines) != 3 {
		t.Fatalf("renderFooter() returned %d lines, want 3", len(lines))
	}
	if strings.TrimSpace(plainTerminal(lines[1])) != "" {
		t.Fatalf("renderFooter() status line = %q, want blank", lines[1])
	}
	for i, line := range lines {
		if got := lipgloss.Width(line); got > 80 {
			t.Fatalf("renderFooter() line %d width = %d, want <= 80", i, got)
		}
	}
}

func TestView_footerRendersStatusMessage(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.StatusMessage = "Router set to v1.1 E05-tasking-permissions/T004"
	footer := plainTerminal(m.renderFooter(80))

	if !strings.Contains(footer, "Router set to v1.1 E05-tasking-permissions/T004") {
		t.Fatal("renderFooter() missing status message")
	}
}

func TestView_footerWrapsLongStatusMessage(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.StatusMessage = "Task not moved: complexity_reason has 24 words; maximum is 20 words. Shorten complexity_reason or adjust complexity_tier before retrying"
	footer := m.renderFooter(48)

	lines := strings.Split(footer, "\n")
	if len(lines) <= 3 {
		t.Fatalf("renderFooter() returned %d lines, want wrapped status lines", len(lines))
	}
	for i, line := range lines {
		if got := lipgloss.Width(line); got > 48 {
			t.Fatalf("renderFooter() line %d width = %d, want <= 48; line=%q", i, got, line)
		}
	}
	plain := plainTerminal(footer)
	if !strings.Contains(plain, "Task not moved: complexity_reason") {
		t.Fatalf("renderFooter() missing wrapped status start; got %q", plain)
	}
	if !strings.Contains(plain, "Shorten complexity_reason") {
		t.Fatalf("renderFooter() missing wrapped status continuation; got %q", plain)
	}
}

func TestView_containsBottomDivider(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Width = 120
	got := m.View()
	// There should be at least two divider lines (top and bottom)
	count := strings.Count(got, "─")
	if count < 2 {
		t.Errorf("View() expected at least 2 divider chars, got %d", count)
	}
}

func TestView_defaultWidthWhenZero(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	// Width=0: should use default and not panic
	got := m.View()
	if got == "" {
		t.Error("View() with zero Width returned empty string")
	}
}

func TestView_taskLabelFallback(t *testing.T) {
	tasks := []data.Task{{ID: "T1", Title: "", Column: data.ColumnPlanned}}
	m := NewModel(tasks, "v1", "E03")
	m.Width = 120
	got := m.View()
	if !strings.Contains(got, "T1") {
		t.Error("View() missing task ID when title is empty")
	}
}

func TestView_taskLabelWithTitle(t *testing.T) {
	tasks := []data.Task{{ID: "T2", Title: "My Task", Column: data.ColumnPlanned}}
	m := NewModel(tasks, "v1", "E03")
	m.Width = 120
	got := m.View()
	if !strings.Contains(got, "My Task") {
		t.Error("View() missing task title")
	}
}

func TestView_wideShowsEpicPanel(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Width = 120
	got := m.View()
	if !strings.Contains(got, "E03") {
		t.Error("View() at width>=120 missing epic panel content")
	}
}

func TestFormatNextActivity_nil(t *testing.T) {
	if got := FormatNextActivity(nil); got != "" {
		t.Errorf("FormatNextActivity(nil) = %q, want empty", got)
	}
}

func TestFormatNextActivity_states(t *testing.T) {
	cases := []struct {
		state *data.RouterState
		want  string
	}{
		{&data.RouterState{State: "task-building", Task: "T010", Epic: "E06", Release: "v1"}, "Build v1 E06/T010"},
		{&data.RouterState{State: "audit-pending", Epic: "E06"}, "Audit E06"},
		{&data.RouterState{State: "epic-design", Epic: "E06"}, "Design E06"},
		{&data.RouterState{State: "epic-task-breakdown", Epic: "E06"}, "Plan E06"},
		{&data.RouterState{State: "pre-implementation", Release: "v1"}, "Planning v1"},
	}
	for _, c := range cases {
		got := FormatNextActivity(c.state)
		if got != c.want {
			t.Errorf("FormatNextActivity(%q) = %q, want %q", c.state.State, got, c.want)
		}
	}
}

func TestFormatNextActivity_truncation(t *testing.T) {
	state := &data.RouterState{State: "task-building", Task: "T001", Epic: "E01-very-long-epic-name", Release: "v1.1"}
	got := FormatNextActivity(state)
	if lipgloss.Width(got) > 20 {
		t.Errorf("FormatNextActivity truncation: width %d > 20, got %q", lipgloss.Width(got), got)
	}
}

func TestFormatNextActivity_taskBuildingKeepsMinorReleaseVisible(t *testing.T) {
	state := &data.RouterState{State: "task-building", Task: "T001", Epic: "E03", Release: "v1.1"}
	got := FormatNextActivity(state)
	if got != "Build v1.1 E03/T001" {
		t.Errorf("FormatNextActivity() = %q, want Build v1.1 E03/T001", got)
	}
}

func TestView_headerShowsNextActivity(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Width = 120
	m.RouterState = &data.RouterState{State: "audit-pending", NextAction: "Audit E06"}
	got := m.View()
	if !strings.Contains(got, "AUDIT:") {
		t.Error("View() missing Next Activity phase tag")
	}
	if !strings.Contains(got, "Audit E06") {
		t.Error("View() missing activity text below header")
	}
	if strings.Contains(got, "Next Activity:") {
		t.Error("View() should not render Next Activity inside the header")
	}
}

func TestView_nextActivityLineImmediatelyBelowHeader(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Width = 120
	m.RouterState = &data.RouterState{State: "task-building", NextAction: "Build T010 (E06) v1"}
	got := m.View()

	lines := strings.Split(got, "\n")
	headerIndex := -1
	activityIndex := -1
	dividerIndex := -1
	for i, line := range lines {
		if strings.Contains(line, "S A V E P O I N T") {
			headerIndex = i
		}
		if strings.Contains(line, "BUILD:") && strings.Contains(line, "Build T010 (E06) v1") {
			activityIndex = i
		}
		if dividerIndex == -1 && strings.Contains(line, strings.Repeat("─", 120)) {
			dividerIndex = i
		}
	}
	if headerIndex == -1 || activityIndex == -1 || dividerIndex == -1 {
		t.Fatalf("View() missing expected header/activity/divider lines: header=%d activity=%d divider=%d", headerIndex, activityIndex, dividerIndex)
	}
	if !(headerIndex < activityIndex && activityIndex < dividerIndex) {
		t.Fatalf("Next Activity line order invalid: header=%d activity=%d divider=%d", headerIndex, activityIndex, dividerIndex)
	}
}

func TestView_headerNoActivityWhenNilState(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Width = 120
	m.RouterState = nil
	got := m.View()
	if strings.Contains(got, "Next Activity:") || strings.Contains(got, "BUILD:") || strings.Contains(got, "PLAN:") || strings.Contains(got, "AUDIT:") {
		t.Error("View() should not show Next Activity line when RouterState is nil")
	}
}

func TestView_headerNarrowWidth(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Width = 40
	m.RouterState = &data.RouterState{State: "audit-pending", NextAction: "Audit E06"}
	got := m.View()
	// Should not panic and header text should still be present
	if !strings.Contains(got, "S A V E P O I N T") {
		t.Error("View() at narrow width missing header text")
	}
}

func TestRenderNextActivityLine_phaseMapping(t *testing.T) {
	cases := []struct {
		name  string
		state *data.RouterState
		tag   string
	}{
		{"build", &data.RouterState{State: "task-building", NextAction: "Build T010 (E06) v1"}, "BUILD:"},
		{"audit", &data.RouterState{State: "audit-pending", NextAction: "Audit E03"}, "AUDIT:"},
		{"pre implementation", &data.RouterState{State: "pre-implementation", NextAction: "Plan v1.1"}, "PLAN:"},
		{"epic design", &data.RouterState{State: "epic-design", NextAction: "Design E03"}, "PLAN:"},
		{"task breakdown", &data.RouterState{State: "epic-task-breakdown", NextAction: "Break down E03"}, "PLAN:"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := renderNextActivityLine(tc.state, 80)
			if !strings.Contains(got, tc.tag) {
				t.Fatalf("renderNextActivityLine() missing phase tag %q in %q", tc.tag, got)
			}
			if !strings.Contains(got, tc.state.NextAction) {
				t.Fatalf("renderNextActivityLine() missing next_action %q in %q", tc.state.NextAction, got)
			}
		})
	}
}

func TestFormatNextActivity_defectBuilding(t *testing.T) {
	state := &data.RouterState{State: "defect-building", Defect: "D001-nil-pointer"}
	got := FormatNextActivity(state)
	if got != "Defect D001" {
		t.Errorf("FormatNextActivity(defect-building) = %q, want %q", got, "Defect D001")
	}
}

func TestRenderNextActivityLine_defectBuilding(t *testing.T) {
	state := &data.RouterState{State: "defect-building", Defect: "D001-nil-pointer", NextAction: "Fix D001"}
	got := renderNextActivityLine(state, 80)
	if !strings.Contains(got, "DEFECT:") {
		t.Fatalf("renderNextActivityLine(defect-building) missing DEFECT: tag; got %q", got)
	}
	if !strings.Contains(got, "Fix D001") {
		t.Fatalf("renderNextActivityLine(defect-building) missing next_action; got %q", got)
	}
}

func TestIsRouterPriority_defectBuildingSkipsTaskCards(t *testing.T) {
	task := data.Task{ID: "E01/T001", Epic: "E01", Release: "v1.1", Column: data.ColumnPlanned}
	state := &data.RouterState{State: "defect-building", Task: "E01/T001", Epic: "E01", Release: "v1.1"}
	card := RenderCard(task, 30, false, state, "")
	// Priority glyph (green ▣) should NOT appear when state is defect-building
	plain := plainTerminal(card)
	_ = plain
	// isRouterPriority must return false — verify via FormatNextActivity not crashing and no green highlight
	// Direct function test
	if isRouterPriority(task, state) {
		t.Error("isRouterPriority() = true for defect-building state, want false")
	}
}

func TestRenderNextActivityLine_hiddenStates(t *testing.T) {
	cases := []*data.RouterState{
		nil,
		{State: "idle", NextAction: "Wait"},
		{State: "task-building", NextAction: ""},
	}
	for _, state := range cases {
		if got := renderNextActivityLine(state, 80); got != "" {
			t.Fatalf("renderNextActivityLine(%v) = %q, want empty", state, got)
		}
	}
}

func TestRenderNextActivityLine_truncatesAtNarrowWidth(t *testing.T) {
	got := renderNextActivityLine(&data.RouterState{
		State:      "pre-implementation",
		NextAction: "Build T010 (E06) v1 with a very long follow-up activity",
	}, 18)
	if lipgloss.Width(got) > 18 {
		t.Fatalf("renderNextActivityLine() width = %d, want <= 18; got %q", lipgloss.Width(got), got)
	}
	if !strings.Contains(got, "PLAN:") || !strings.Contains(got, "…") {
		t.Fatalf("renderNextActivityLine() = %q, want PLAN tag and ellipsis", got)
	}
}

func BenchmarkCalculateLayout_narrow(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		CalculateLayout(60, 24)
	}
}

func BenchmarkCalculateLayout_standard(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		CalculateLayout(80, 24)
	}
}

func BenchmarkCalculateLayout_wide(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		CalculateLayout(120, 24)
	}
}

func BenchmarkCalculateLayout_extraWide(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		CalculateLayout(220, 50)
	}
}

func BenchmarkView_empty(b *testing.B) {
	m := NewModel(nil, "v1", "E03")
	m.Width = 120
	m.Height = 40
	b.ReportAllocs()
	for b.Loop() {
		m.View()
	}
}

func BenchmarkView_withTasks(b *testing.B) {
	tasks := make([]data.Task, 15)
	stages := []data.ProgressStage{data.StageBuild, data.StageTest, data.StageAudit}
	cols := []data.ColumnType{data.ColumnPlanned, data.ColumnInProgress, data.ColumnDone}
	for i := range tasks {
		tasks[i] = data.Task{
			ID:      fmt.Sprintf("E06-layout/T%03d-task-slug", i+1),
			Title:   fmt.Sprintf("Task %d with a reasonable title length", i+1),
			Column:  cols[i%3],
			Stage:   stages[i%3],
			Release: "v1.1",
			Epic:    "E06-layout",
		}
	}
	m := NewModel(tasks, "v1.1", "E06")
	m.Width = 120
	m.Height = 40
	b.ReportAllocs()
	for b.Loop() {
		m.View()
	}
}

func TestView_headerHidesDefectSignalWhenZeroOpen(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Width = 120
	m.SelectedRelease = "v1"
	m.AllDefects = []data.Defect{
		{Release: "v1", Status: data.DefectResolved},
	}
	got := m.View()
	if strings.Contains(got, "⚠") {
		t.Error("header must not show defect signal when zero open defects")
	}
}

func TestView_headerShowsOpenDefectCount(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Width = 120
	m.SelectedRelease = "v1"
	m.AllDefects = []data.Defect{
		{Release: "v1", Status: data.DefectOpen},
		{Release: "v1", Status: data.DefectInProgress},
		{Release: "v1", Status: data.DefectResolved},
	}
	got := plainTerminal(m.View())
	if !strings.Contains(got, "⚠  2 open") {
		t.Errorf("header missing defect signal; got %q", got)
	}
}

func TestView_headerDefectSignalRespectsSelectedRelease(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Width = 120
	m.SelectedRelease = "v1"
	m.AllDefects = []data.Defect{
		{Release: "v2", Status: data.DefectOpen},
	}
	got := m.View()
	if strings.Contains(got, "⚑") {
		t.Error("defect signal must not show defects from other releases")
	}
}

func TestView_headerShowsReleaseWithDefects(t *testing.T) {
	m := NewModel(nil, "v1.2", "E03")
	m.Width = 120
	m.SelectedRelease = "v1.2"
	m.AllDefects = []data.Defect{
		{Release: "v1.2", Status: data.DefectOpen},
	}
	got := plainTerminal(m.View())
	if !strings.Contains(got, "v1.2 │ ⚠  1 open") {
		t.Errorf("header missing release+defect cluster; got %q", got)
	}
}

func TestView_headerShowsReleaseAtZeroDefects(t *testing.T) {
	m := NewModel(nil, "v1.2", "E03")
	m.Width = 120
	m.SelectedRelease = "v1.2"
	got := plainTerminal(m.renderHeader(m.Width))
	if !strings.Contains(got, "v1.2") {
		t.Errorf("header missing release at zero defects; got %q", got)
	}
	if strings.Contains(got, "│") {
		t.Errorf("header should not show a separator at zero defects; got %q", got)
	}
}

func TestView_headerEmptyReleaseNoSeparator(t *testing.T) {
	m := NewModel(nil, "", "E03")
	m.Width = 120
	m.SelectedRelease = ""
	m.AllDefects = []data.Defect{
		{Release: "", Status: data.DefectOpen},
	}
	got := plainTerminal(m.renderHeader(m.Width))
	if strings.Contains(got, "│") {
		t.Errorf("empty release must not emit a separator; got %q", got)
	}
	if !strings.Contains(got, "⚠  1 open") {
		t.Errorf("header should still show defect signal with empty release; got %q", got)
	}
}

func TestView_headerReleaseNarrowWidthNoOverflow(t *testing.T) {
	m := NewModel(nil, "v1.2", "E03")
	m.Width = 24
	m.SelectedRelease = "v1.2"
	m.AllDefects = []data.Defect{
		{Release: "v1.2", Status: data.DefectOpen},
	}
	header := plainTerminal(m.renderHeader(m.Width))
	for _, line := range strings.Split(header, "\n") {
		if lipgloss.Width(line) > m.Width {
			t.Errorf("header line overflows width %d: %q (width %d)", m.Width, line, lipgloss.Width(line))
		}
	}
}

func TestView_threeColumnLayoutUnchangedWithDefects(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Width = 100
	m.Height = 40
	m.AllDefects = []data.Defect{
		{Release: "v1", Status: data.DefectOpen},
	}
	l := CalculateLayout(m.Width, m.Height)
	if l.ColCount != 3 {
		t.Fatalf("ColCount = %d, want 3", l.ColCount)
	}
	if l.EpicPanelVisible {
		t.Error("epic panel should be hidden at width=100")
	}
	got := m.View()
	for _, title := range []string{"PLANNED", "IN PROGRESS", "DONE"} {
		if !strings.Contains(got, title) {
			t.Errorf("three-column layout missing %q with defects present", title)
		}
	}
}

func TestView_narrowShowsSingleColumn(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Width = 60
	m.FocusedColumn = data.ColumnInProgress
	got := m.View()
	if !strings.Contains(got, "IN PROGRESS") {
		t.Error("View() at width<80 missing focused column title")
	}
	if strings.Contains(got, "PLANNED") {
		t.Error("View() at width<80 should not show non-focused columns")
	}
}
