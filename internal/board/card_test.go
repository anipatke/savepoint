package board

import (
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/data"
)

func TestRenderCard_containsID(t *testing.T) {
	task := data.Task{ID: "E04/T002", Title: "Build card", Stage: data.StageBuild}
	got := RenderCard(task, 30, false)
	if !strings.Contains(got, "T002") {
		t.Error("RenderCard missing short task ID")
	}
}

func TestRenderCard_containsTitle(t *testing.T) {
	task := data.Task{ID: "T1", Title: "My title", Stage: data.StageBuild}
	got := RenderCard(task, 30, false)
	if !strings.Contains(got, "My title") {
		t.Error("RenderCard missing task title")
	}
}

func TestRenderCard_containsBuildGlyph(t *testing.T) {
	task := data.Task{ID: "T1", Stage: data.StageBuild}
	got := RenderCard(task, 30, false)
	if !strings.Contains(got, glyphBuild) {
		t.Errorf("RenderCard missing build glyph %q", glyphBuild)
	}
}

func TestRenderCard_containsTestGlyph(t *testing.T) {
	task := data.Task{ID: "T1", Stage: data.StageTest}
	got := RenderCard(task, 30, false)
	if !strings.Contains(got, glyphTest) {
		t.Errorf("RenderCard missing test glyph %q", glyphTest)
	}
}

func TestRenderCard_containsAuditGlyph(t *testing.T) {
	task := data.Task{ID: "T1", Stage: data.StageAudit}
	got := RenderCard(task, 30, false)
	if !strings.Contains(got, glyphAudit) {
		t.Errorf("RenderCard missing audit glyph %q", glyphAudit)
	}
}

func TestRenderCard_focusedDoesNotPanic(t *testing.T) {
	task := data.Task{ID: "T1", Title: "hello", Stage: data.StageBuild}
	got := RenderCard(task, 30, true)
	if got == "" {
		t.Error("RenderCard focused returned empty string")
	}
}

func TestRenderCard_titleTruncated(t *testing.T) {
	long := "This is a very long title that should be truncated for sure"
	task := data.Task{ID: "T1", Title: long, Stage: data.StageBuild}
	got := RenderCard(task, 20, false)
	if strings.Contains(got, long) {
		t.Error("RenderCard should truncate long title")
	}
	if !strings.Contains(got, "…") {
		t.Error("RenderCard should include ellipsis when title truncated")
	}
}

func TestRenderCard_idTruncated(t *testing.T) {
	long := "E04-board-components/T999-very-long-id"
	task := data.Task{ID: long, Stage: data.StageBuild}
	got := RenderCard(task, 20, false)
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
	got := RenderCard(task, 30, false)
	if !strings.Contains(got, glyphBuild) {
		t.Error("RenderCard with empty stage should use build glyph")
	}
}
