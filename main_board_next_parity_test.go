package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	xansi "github.com/charmbracelet/x/ansi"
	boardv2 "github.com/opencode/savepoint/internal/board/v2"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/doctor"
	"github.com/opencode/savepoint/internal/migrate"
	"github.com/opencode/savepoint/internal/resume"
)

// This file proves three things across every rung the on-disk rung matrix
// (resumeMatrixCases, E48 T007) can produce:
//
//  1. The board's own two surfaces — the drawn TUI and the deterministic
//     non-TTY plain table — always render the same one-line Next summary as
//     each other, and that line matches what the resolved projection's own
//     Task or Objective says (expectedBoardLine).
//  2. `savepoint resume` starts with that exact same plain-text line, then
//     reports the full identity/evidence/action narrative below it.
//  3. A selected Task carries its owning Objective in the shared projection
//     when the Objective record exists.
//
// The rest of the board summary and resume's full narrative may still differ:
// the board is a glance, resume is the report (.savepoint/Design.md, "Layout").

func TestBoardNextAndResumeReportTheSameAnswer(t *testing.T) {
	for _, test := range resumeMatrixCases() {
		t.Run(test.name, func(t *testing.T) {
			dir := t.TempDir()
			test.build(t, dir)

			want := resolveNextFromDisk(t, dir)
			if want.Kind != test.wantKind {
				t.Fatalf("Next.Kind = %q, want %q", want.Kind, test.wantKind)
			}
			if want.Task != nil && (want.Objective == nil || want.Objective.ID != want.Task.Objective) {
				t.Errorf("Next Objective = %#v for Task owner %q, want the owning Objective record", want.Objective, want.Task.Objective)
			}

			assertResumeReportsTheProjection(t, dir, want)

			expected := expectedBoardLine(want)
			assertBoardLine(t, dir, expected)
			assertTUILine(t, dir, expected)
		})
	}
}

func TestMigratedReleaseFlowsThroughDoctorBoardSelectorPlainAndResume(t *testing.T) {
	dir := copyMigrateFixture(t)
	plan, err := migrate.Plan(dir, nil,
		func() time.Time { return time.Date(2026, 9, 18, 12, 0, 0, 0, time.UTC) },
		func() string { return "op-main-release-proof" })
	if err != nil {
		t.Fatalf("Plan() error = %v", err)
	}
	if _, err := migrate.Apply(dir, plan); err != nil {
		t.Fatalf("Apply() error = %v", err)
	}

	savepointRoot := filepath.Join(dir, ".savepoint")
	index, err := data.LoadV2Index(savepointRoot)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	if len(index.Releases) != 1 {
		t.Fatalf("migrated Releases = %d, want one first-class Release", len(index.Releases))
	}
	problems := doctor.RunV2Checks(savepointRoot).Releases
	if len(problems) == 0 {
		t.Fatal("doctor reported no Release readiness diagnostic for the migrated incomplete fixture")
	}

	// Plain board, TUI board, and resume must each still report the same
	// resolved projection after migration, not merely agree on a hand-built
	// V2 fixture — each checked against its own wording now, not each other.
	want := resolveNextFromDisk(t, dir)
	assertResumeReportsTheProjection(t, dir, want)
	expected := expectedBoardLine(want)
	assertBoardLine(t, dir, expected)
	assertTUILine(t, dir, expected)

	beforeSelector := snapshotDir(t, dir)
	model := boardv2.NewModel(boardv2.Options{Root: savepointRoot})
	sized, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 48})
	loaded, _ := sized.(boardv2.Model).Update(sized.(boardv2.Model).Init()())
	board := loaded.(boardv2.Model)
	opened, cmd := board.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	if cmd != nil {
		t.Fatal("opening migrated Release selector returned a command")
	}
	selector := opened.(boardv2.Model)
	if !selector.ReleaseOverlay || selector.SelectedRelease == "" || len(selector.Releases) != len(index.Releases) {
		t.Fatalf("migrated Release selector = overlay %t, selected %q, releases %v", selector.ReleaseOverlay, selector.SelectedRelease, selector.Releases)
	}
	view := xansi.Strip(selector.View())
	if !strings.Contains(view, selector.SelectedRelease) {
		t.Fatalf("migrated Release selector view omitted selected Release %s:\n%s", selector.SelectedRelease, view)
	}
	detailed, cmd := selector.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})
	if cmd != nil || detailed.(boardv2.Model).Detail == nil || detailed.(boardv2.Model).Detail.Kind != boardv2.DetailRelease {
		t.Fatalf("migrated Release detail = %T, command %t, want read-only Release detail", detailed, cmd != nil)
	}
	assertSameSnapshot(t, beforeSelector, snapshotDir(t, dir))
}

// TestBuiltBoardAndResumeReportTheSameAnswer uses the command dispatch in a
// subprocess for every on-disk rung. The subprocess is the test-built
// Savepoint binary, so this closes the gap between package-level parity and
// the user's `savepoint board` / `savepoint resume` paths.
func TestBuiltBoardAndResumeReportTheSameAnswer(t *testing.T) {
	for _, test := range resumeMatrixCases() {
		t.Run(test.name, func(t *testing.T) {
			dir := t.TempDir()
			test.build(t, dir)

			want := resolveNextFromDisk(t, dir)

			boardResult := runMainInDirForTest(t, dir, []string{"board"})
			if boardResult.err != nil {
				t.Fatalf("built board failed: %v\nstderr: %s", boardResult.err, boardResult.stderr)
			}
			resumeResult := runMainInDirForTest(t, dir, []string{"resume"})
			if resumeResult.err != nil {
				t.Fatalf("built resume failed: %v\nstderr: %s", resumeResult.err, resumeResult.stderr)
			}

			expected := expectedBoardLine(want)
			if !hasExactLine(boardResult.stdout, expected) {
				t.Errorf("built board output missing exact Next line %q:\n%s", expected, boardResult.stdout)
			}
			resumeFirstLine := strings.SplitN(resumeResult.stdout, "\n", 2)[0]
			if resumeFirstLine != expected {
				t.Errorf("built resume first line = %q, want the shared Next line %q", resumeFirstLine, expected)
			}
			if want.Task != nil && !strings.Contains(resumeResult.stdout, "Task: "+want.Task.ID+" — "+want.Task.Title) {
				t.Errorf("built resume output does not name the projection's own Task:\n%s", resumeResult.stdout)
			}
			if !strings.Contains(resumeResult.stdout, "Next action: ") {
				t.Errorf("built resume output carries no next action:\n%s", resumeResult.stdout)
			}
		})
	}
}

// expectedBoardLine is the shared wording used by the board and resume.
func expectedBoardLine(next data.Next) string {
	return resume.NextLine(next)
}

// assertResumeReportsTheProjection runs `savepoint resume` over dir and
// checks its output against want's own identity and action — resume's full
// narrative, unaffected by the board's simplified panel.
func assertResumeReportsTheProjection(t *testing.T, dir string, want data.Next) {
	t.Helper()
	var out bytes.Buffer
	code, err := runResume(dir, &out)
	if err != nil {
		t.Fatalf("runResume() error = %v", err)
	}
	if code != 0 {
		t.Fatalf("runResume() exit = %d, want 0\n%s", code, out.String())
	}
	text := out.String()
	firstLine := strings.SplitN(text, "\n", 2)[0]
	if firstLine != expectedBoardLine(want) {
		t.Errorf("resume first line = %q, want shared Next line %q", firstLine, expectedBoardLine(want))
	}
	if want.Task != nil && !strings.Contains(text, "Task: "+want.Task.ID+" — "+want.Task.Title) {
		t.Errorf("resume output does not name the projection's own Task:\n%s", text)
	}
	if want.Objective != nil && !strings.Contains(text, "Objective: "+want.Objective.ID+" — "+want.Objective.Title) {
		t.Errorf("resume output does not name the projection's own Objective:\n%s", text)
	}
	if !strings.Contains(text, "Next action: ") {
		t.Errorf("resume output carries no next action:\n%s", text)
	}
}

// assertBoardLine runs the board's non-TTY rendering over dir and checks it
// contains expected — the board's real output path, not a test-only view of
// it.
func assertBoardLine(t *testing.T, dir string, expected string) {
	t.Helper()
	var out bytes.Buffer
	if err := boardv2.Run(boardv2.Options{Root: filepath.Join(dir, ".savepoint"), Stdout: &out, TTY: false}); err != nil {
		t.Fatalf("board Run() error = %v", err)
	}
	if !hasExactLine(out.String(), expected) {
		t.Errorf("non-TTY board output missing exact Next line %q:\n%s", expected, out.String())
	}
}

func hasExactLine(output, expected string) bool {
	for _, line := range strings.Split(output, "\n") {
		if line == expected {
			return true
		}
	}
	return false
}

// assertTUILine drives the actual Bubble Tea model through its load command
// and checks the drawn Next area contains expected, proving the parity
// covers the drawn surface as well as the non-TTY formatter.
func assertTUILine(t *testing.T, dir string, expected string) {
	t.Helper()
	model := boardv2.NewModel(boardv2.Options{Root: filepath.Join(dir, ".savepoint")})
	sized, _ := model.Update(tea.WindowSizeMsg{Width: 200, Height: 48})
	sizedModel := sized.(boardv2.Model)
	loaded, _ := sizedModel.Update(sizedModel.Init()())
	final := loaded.(boardv2.Model)
	got := xansi.Strip(final.View())
	if !strings.Contains(got, "NEXT: "+expected) {
		t.Errorf("TUI board view missing %q:\n%s", "NEXT: "+expected, got)
	}
}
