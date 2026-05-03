package board

import (
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/styles"
)

func TestRenderCard_containsID(t *testing.T) {
	task := data.Task{ID: "E04/T002", Title: "Build card", Stage: data.StageBuild}
	got := RenderCard(task, 30, false, nil)
	if !strings.Contains(got, "T002") {
		t.Error("RenderCard missing short task ID")
	}
}

func TestRenderCard_containsTitle(t *testing.T) {
	task := data.Task{ID: "T1", Title: "My title", Stage: data.StageBuild}
	got := RenderCard(task, 30, false, nil)
	if !strings.Contains(got, "My title") {
		t.Error("RenderCard missing task title")
	}
}

func TestRenderCard_containsBuildGlyph(t *testing.T) {
	task := data.Task{ID: "T1", Stage: data.StageBuild}
	got := RenderCard(task, 30, false, nil)
	if !strings.Contains(got, glyphBuild) {
		t.Errorf("RenderCard missing build glyph %q", glyphBuild)
	}
}

func TestRenderCard_containsTestGlyph(t *testing.T) {
	task := data.Task{ID: "T1", Stage: data.StageTest}
	got := RenderCard(task, 30, false, nil)
	if !strings.Contains(got, glyphTest) {
		t.Errorf("RenderCard missing test glyph %q", glyphTest)
	}
}

func TestRenderCard_containsAuditGlyph(t *testing.T) {
	task := data.Task{ID: "T1", Stage: data.StageAudit}
	got := RenderCard(task, 30, false, nil)
	if !strings.Contains(got, glyphAudit) {
		t.Errorf("RenderCard missing audit glyph %q", glyphAudit)
	}
}

func TestRenderCard_focusedDoesNotPanic(t *testing.T) {
	task := data.Task{ID: "T1", Title: "hello", Stage: data.StageBuild}
	got := RenderCard(task, 30, true, nil)
	if got == "" {
		t.Error("RenderCard focused returned empty string")
	}
}

func TestRenderCard_titleWraps(t *testing.T) {
	long := "This is a very long title that should be wrapped for sure"
	task := data.Task{ID: "T1", Title: long, Stage: data.StageBuild}
	got := RenderCard(task, 20, false, nil)
	// full title as one line does not fit; it must be broken up
	if strings.Contains(got, long) {
		t.Error("RenderCard should wrap long title, not render it as one line")
	}
	// words from the title must still appear somewhere in the output
	if !strings.Contains(got, "This") {
		t.Error("RenderCard wrapped title missing expected content")
	}
}

func TestRenderCard_idTruncated(t *testing.T) {
	long := "E04-board-components/T999-very-long-id"
	task := data.Task{ID: long, Stage: data.StageBuild}
	got := RenderCard(task, 20, false, nil)
	if strings.Contains(got, long) {
		t.Error("RenderCard should truncate long ID")
	}
}

func TestTruncate_shortString(t *testing.T) {
	if truncate("hi", 10) != "hi" {
		t.Error("truncate should not clip short string")
	}
}

func TestTruncate_exactLength(t *testing.T) {
	if truncate("hello", 5) != "hello" {
		t.Error("truncate should not clip string at exact max")
	}
}

func TestTruncate_clipsWithEllipsis(t *testing.T) {
	got := truncate("hello", 4)
	if got != "hel…" {
		t.Errorf("truncate got %q, want %q", got, "hel…")
	}
}

func TestTruncate_maxOne(t *testing.T) {
	got := truncate("hello", 1)
	if got != "…" {
		t.Errorf("truncate(max=1) got %q, want %q", got, "…")
	}
}

func TestRenderCard_defaultStageUsesBuildGlyph(t *testing.T) {
	task := data.Task{ID: "T1", Stage: ""}
	got := RenderCard(task, 30, false, nil)
	if !strings.Contains(got, glyphBuild) {
		t.Error("RenderCard with empty stage should use build glyph")
	}
}

func TestRenderCard_routerPriorityUsesGreenGlyph(t *testing.T) {
	task := data.Task{ID: "E06/T009", Release: "v1", Epic: "E06", Stage: data.StageTest}
	router := &data.RouterState{Release: "v1", Epic: "E06", Task: "E06/T009"}
	got := RenderCard(task, 30, false, router)
	if !isRouterPriority(task, router) {
		t.Error("router priority should match release, epic, and task")
	}
	if !strings.Contains(got, glyphBuild) {
		t.Error("router priority card should use build glyph")
	}
	nonPriority := RenderCard(task, 30, false, nil)
	if !strings.Contains(nonPriority, glyphTest) {
		t.Error("non-priority test card should use test glyph")
	}
}

func TestRenderCard_noBackgroundFillEscapes(t *testing.T) {
	task := data.Task{ID: "E06/T009", Title: "Router priority", Release: "v1", Epic: "E06", Stage: data.StageTest}
	router := &data.RouterState{Release: "v1", Epic: "E06", Task: "E06/T009"}
	got := RenderCard(task, 30, false, router)
	if strings.Contains(got, "\x1b[48;") || strings.Contains(got, "\x1b[40m") {
		t.Fatalf("RenderCard should not emit background fills; got %q", got)
	}
}

func TestRenderCard_routerPriorityMatchesShortID(t *testing.T) {
	// Router stores short IDs ("T009"); task ID is full slug — must still match.
	task := data.Task{ID: "E06-atari-noir-layout/T009-router-priority", Release: "v1", Epic: "E06-atari-noir-layout", Stage: data.StageTest}
	router := &data.RouterState{Release: "v1", Epic: "E06", Task: "T009"}
	got := RenderCard(task, 30, false, router)
	if !isRouterPriority(task, router) {
		t.Error("short router task ID should match full task ID slug")
	}
	if !strings.Contains(got, glyphBuild) {
		t.Error("router priority card should use build glyph")
	}
}

func TestRenderCard_staleRouterTaskNoMatch(t *testing.T) {
	// Task moved to a new epic; router still has old epic path — should NOT match a different task number.
	task := data.Task{ID: "E03-header-activity/T001-border-resize-fix", Release: "v1", Epic: "E03-header-activity", Stage: data.StageBuild}
	router := &data.RouterState{Release: "v1", Epic: "E03", Task: "T002"}
	got := RenderCard(task, 30, false, router)
	if isRouterPriority(task, router) {
		t.Error("stale router pointing to different task number should not show green glyph")
	}
	if !strings.Contains(got, styles.GlyphBuild.Render(glyphBuild)) {
		t.Error("non-priority build task should use orange build glyph")
	}
}

func TestRenderCard_routerSameTaskNumberDifferentEpicNoMatch(t *testing.T) {
	task := data.Task{ID: "E03-header-activity/T001-border-resize-fix", Release: "v1", Epic: "E03-header-activity", Stage: data.StageTest}
	router := &data.RouterState{Release: "v1", Epic: "E01", Task: "T001"}
	got := RenderCard(task, 30, false, router)
	if isRouterPriority(task, router) {
		t.Error("router priority should not match same task number in a different epic")
	}
	if !strings.Contains(got, styles.GlyphTest.Render(glyphTest)) {
		t.Error("non-priority test task should keep test glyph")
	}
}

func TestRenderCard_doneTaskUsesOrangeBuildGlyph(t *testing.T) {
	task := data.Task{ID: "E03/T001", Release: "v1", Epic: "E03", Column: data.ColumnDone, Stage: data.StageTest}
	router := &data.RouterState{Release: "v1", Epic: "E03", Task: "T001"}
	got := RenderCard(task, 30, false, router)
	if !isRouterPriority(task, router) {
		t.Error("router state should still identify the matching done task")
	}
	if !strings.Contains(got, styles.GlyphBuild.Render(glyphBuild)) {
		t.Error("done task should use orange build glyph")
	}
	if strings.Contains(got, glyphTest) {
		t.Error("done task should not use test glyph")
	}
}

func TestRenderCard_explicitStatusUsesUnifiedGlyph(t *testing.T) {
	tests := []struct {
		name   string
		status data.TaskStatus
		glyph  string
	}{
		{"planned", data.StatusPlanned, "○"},
		{"done", data.StatusDone, "◉"},
		{"audited", data.StatusAudited, "✓"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := data.Task{ID: "T1", Status: string(tt.status), Stage: data.StageAudit}
			got := RenderCard(task, 30, false, nil)
			if !strings.Contains(got, tt.glyph) {
				t.Errorf("RenderCard with status %q missing glyph %q", tt.status, tt.glyph)
			}
			if strings.Contains(got, glyphAudit) {
				t.Errorf("RenderCard with status %q should not fall back to audit glyph", tt.status)
			}
		})
	}
}

func TestRenderCard_inProgressShowsPhaseText(t *testing.T) {
	tests := []struct {
		name  string
		stage data.ProgressStage
		label string
		glyph string
	}{
		{"build", data.StageBuild, "BUILD", glyphBuild},
		{"test", data.StageTest, "TEST", glyphTest},
		{"audit", data.StageAudit, "AUDIT", glyphAudit},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := data.Task{ID: "T1", Column: data.ColumnInProgress, Status: string(data.StatusInProgress), Stage: tt.stage}
			got := RenderCard(task, 30, false, nil)
			if !strings.Contains(got, tt.label) {
				t.Errorf("RenderCard missing phase label %q", tt.label)
			}
			if !strings.Contains(got, tt.glyph) {
				t.Errorf("RenderCard missing phase glyph %q", tt.glyph)
			}
			if strings.Contains(got, "▶") {
				t.Error("RenderCard should not use generic in_progress glyph when phase is available")
			}
		})
	}
}

func TestRenderCard_doneShowsDoneText(t *testing.T) {
	task := data.Task{ID: "T1", Column: data.ColumnDone, Status: string(data.StatusDone)}
	got := RenderCard(task, 30, false, nil)
	if !strings.Contains(got, "DONE") {
		t.Error("RenderCard missing DONE phase label")
	}
}
