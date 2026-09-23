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
)

// This file proves two things across every rung the on-disk rung matrix
// (resumeMatrixCases, E48 T007) can produce:
//
//  1. The board's own two surfaces — the drawn TUI and the deterministic
//     non-TTY plain table — always render the same one-line Next summary as
//     each other, and that line matches what the resolved projection's own
//     Task or Objective says (expectedBoardLine).
//  2. `savepoint resume` still reports the full identity/evidence/action
//     narrative it always has, checked against the same projection.
//
// The board's one-line summary and resume's full narrative are deliberately
// NOT required to share wording anymore — the board is a glance, resume is
// the report (.savepoint/Design.md, "Layout"). This file used to also assert
// that cross-surface equality; it no longer does, on purpose.

func TestBoardNextAndResumeReportTheSameAnswer(t *testing.T) {
	for _, test := range resumeMatrixCases() {
		t.Run(test.name, func(t *testing.T) {
			dir := t.TempDir()
			test.build(t, dir)
			if test.wantKind == data.NextPendingMigration {
				// The live V2-only router refuses to render an interrupted
				// migration; it prints recovery guidance before either board or
				// resume can interpret indexed records.
				return
			}

			want := resolveNextFromDisk(t, dir)
			if want.Kind != test.wantKind {
				t.Fatalf("Next.Kind = %q, want %q", want.Kind, test.wantKind)
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
	cutover := data.ResolveReleaseCutover(index)
	if cutover.Allowed || len(cutover.Blockers) == 0 {
		t.Fatalf("migrated cutover decision = %+v, want the incomplete fixture Release blocked", cutover)
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
			if test.wantKind == data.NextPendingMigration {
				// A pending migration is a recovery state, not a V2 board
				// projection. The live command contract is tested by the focused
				// resume/dispatch tests; do not ask the built board to render it.
				return
			}

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
			if !strings.Contains(boardResult.stdout, expected) {
				t.Errorf("built board output missing %q:\n%s", expected, boardResult.stdout)
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

// expectedBoardLine mirrors internal/board/v2.nextLines' one-line format —
// not by calling it (that package's rendering is unexported), but by reading
// the same projection fields it reads, so this test proves the board's real
// output by comparison rather than by construction. Keep this in sync with
// nextLines/taskStageWord in internal/board/v2/next_panel.go.
func expectedBoardLine(next data.Next) string {
	if next.Task != nil {
		return taskStageWordForTest(next.Task) + " " + next.Task.ID + " — " + next.Task.Title
	}
	if next.Objective != nil {
		return next.Objective.ID + " — " + next.Objective.Title
	}
	return "Nothing selected yet"
}

func taskStageWordForTest(task *data.TaskV2) string {
	if task.Status == data.ColumnInProgress {
		switch task.Stage {
		case data.StageBuild:
			return "Build"
		case data.StageTest:
			return "Test"
		case data.StageAudit:
			return "Check"
		}
	}
	switch task.Status {
	case data.ColumnPlanned:
		return "Planned"
	case data.ColumnDone:
		return "Done"
	default:
		return string(task.Status)
	}
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
	if !strings.Contains(out.String(), expected) {
		t.Errorf("non-TTY board output missing %q:\n%s", expected, out.String())
	}
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
	if !strings.Contains(got, expected) {
		t.Errorf("TUI board view missing %q:\n%s", expected, got)
	}
}
