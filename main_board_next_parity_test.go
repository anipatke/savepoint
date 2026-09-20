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

// This file is E49 T005's parity obligation, and it is the point of the whole
// Next area: section 11 of the release design says the board and resume share
// one interpretation, and until a test compares the two surfaces over the same
// project that is an intention rather than a property of the code.
//
// It reuses resumeMatrixCases — the on-disk rung matrix E48 T007 built — so the
// two surfaces are compared across every rung a board user reaches, not over a
// project shape invented here to make them agree.

// nextFacts are the three things both surfaces must report identically: which
// records the projection selected, the rung-specific evidence behind it, and
// the action itself. Identity and evidence are compared as whole lines, so a
// surface that named a different record, or softened a clearance statement,
// fails rather than passing on a shared substring.
type nextFacts struct {
	identity []string
	evidence []string
	action   string
}

func TestBoardNextAndResumeReportTheSameAnswer(t *testing.T) {
	for _, test := range resumeMatrixCases() {
		t.Run(test.name, func(t *testing.T) {
			dir := t.TempDir()
			test.build(t, dir)
			if test.wantKind == data.NextPendingMigration {
				// The live V2-only router refuses to render an interrupted
				// migration; it prints recovery guidance before either board or
				// resume can interpret indexed records. The package-level Next
				// matrix still covers this rung, while the ordinary surface parity
				// below applies only to an intact V2 runtime.
				return
			}

			// The projection both surfaces are claimed to share, resolved once
			// here through runResume's own read order. Every expectation below
			// is derived from it rather than written out by hand, so a rung
			// that changed meaning cannot leave this test asserting the old one.
			want := resolveNextFromDisk(t, dir)
			if want.Kind != test.wantKind {
				t.Fatalf("Next.Kind = %q, want %q", want.Kind, test.wantKind)
			}

			fromResume := resumeNextFacts(t, dir)
			fromBoard := boardNextFacts(t, dir)
			fromTUI := boardTUINextFacts(t, dir)

			assertSameLines(t, "selected records", fromResume.identity, fromBoard.identity)
			assertSameLines(t, "rung evidence", fromResume.evidence, fromBoard.evidence)
			if fromResume.action != fromBoard.action {
				t.Errorf("next action differs:\n resume: %q\n  board: %q", fromResume.action, fromBoard.action)
			}
			assertSameLines(t, "TUI selected records", fromResume.identity, fromTUI.identity)
			assertSameLines(t, "TUI rung evidence", fromResume.evidence, fromTUI.evidence)
			if fromResume.action != fromTUI.action {
				t.Errorf("TUI next action differs:\n resume: %q\n    TUI: %q", fromResume.action, fromTUI.action)
			}

			// Both agreeing on wording they each read from the same place is
			// necessary but not sufficient: pin them to the rung the projection
			// actually landed on, whose evidence lines no other rung produces.
			assertSameLines(t, "rung evidence against the projection", resume.EvidenceLines(want), fromBoard.evidence)
			if action := resume.ActionPhrase(want); action != fromBoard.action {
				t.Errorf("board action = %q, want the projection's own action %q", fromBoard.action, action)
			}
			assertSameLines(t, "selected records against the projection", projectionIdentity(want), fromBoard.identity)
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
	problems := doctor.CheckReleaseReadiness(savepointRoot)
	if len(problems) == 0 {
		t.Fatal("doctor reported no Release readiness diagnostic for the migrated incomplete fixture")
	}

	// Plain board, TUI board, and resume must name the same projection and
	// action after migration, not merely agree on a hand-built V2 fixture.
	want := resolveNextFromDisk(t, dir)
	fromResume := resumeNextFacts(t, dir)
	fromBoard := boardNextFacts(t, dir)
	fromTUI := boardTUINextFacts(t, dir)
	assertSameLines(t, "migrated selected records", fromResume.identity, fromBoard.identity)
	assertSameLines(t, "migrated rung evidence", fromResume.evidence, fromBoard.evidence)
	if fromResume.action != fromBoard.action {
		t.Fatalf("migrated next action differs: resume %q, board %q", fromResume.action, fromBoard.action)
	}
	assertSameLines(t, "migrated TUI selected records", fromResume.identity, fromTUI.identity)
	assertSameLines(t, "migrated TUI rung evidence", fromResume.evidence, fromTUI.evidence)
	if fromResume.action != fromTUI.action {
		t.Fatalf("migrated TUI next action differs: resume %q, TUI %q", fromResume.action, fromTUI.action)
	}
	assertSameLines(t, "migrated projection evidence", resume.EvidenceLines(want), fromBoard.evidence)
	if action := resume.ActionPhrase(want); action != fromBoard.action {
		t.Fatalf("migrated board action = %q, want projection action %q", fromBoard.action, action)
	}

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

// boardTUINextFacts drives the actual Bubble Tea model through its load
// command, then reads the rendered Next area. This makes the parity proof
// cover the drawn surface as well as the non-TTY formatter.
func boardTUINextFacts(t *testing.T, dir string) nextFacts {
	t.Helper()

	model := boardv2.NewModel(boardv2.Options{Root: filepath.Join(dir, ".savepoint")})
	// Use a wide enough virtual terminal that the comparison observes the
	// complete evidence line rather than the TUI's intentional narrow wrapping.
	sized, _ := model.Update(tea.WindowSizeMsg{Width: 200, Height: 48})
	sizedModel := sized.(boardv2.Model)
	loaded, _ := sizedModel.Update(sizedModel.Init()())
	final := loaded.(boardv2.Model)
	panel := nextPanelText(xansi.Strip(final.View()))
	return nextFacts{
		identity: linesWithTrimmedPrefix(panel, "Release: ", "Release outcome: ", "Release status: ", "Objective: ", "Task: "),
		evidence: linesWithTrimmedPrefix(panel,
			"Migration: ", "Replan: ", "Blocked: ", "Completion: ", "Ready: ",
			"Technical clearance: ", "Owner wait: ", "Release readiness: ", "Historical completion: "),
		action: strings.TrimPrefix(lineWithTrimmedPrefix(t, panel, "Action: "), "Action: "),
	}
}

func nextPanelText(rendered string) string {
	lines := strings.Split(rendered, "\n")
	start := -1
	end := len(lines)
	for i, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if start < 0 && strings.HasPrefix(line, "NEXT: ") {
			start = i
			continue
		}
		if start >= 0 && strings.HasPrefix(line, "Action: ") {
			end = i + 1
			break
		}
	}
	if start < 0 {
		return ""
	}
	return strings.Join(lines[start:end], "\n")
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

			boardResult := runMainInDirForTest(t, dir, []string{"board"})
			if boardResult.err != nil {
				t.Fatalf("built board failed: %v\nstderr: %s", boardResult.err, boardResult.stderr)
			}
			resumeResult := runMainInDirForTest(t, dir, []string{"resume"})
			if resumeResult.err != nil {
				t.Fatalf("built resume failed: %v\nstderr: %s", resumeResult.err, resumeResult.stderr)
			}

			boardFacts := factsFromCommandOutput(t, boardResult.stdout, "Action: ")
			resumeFacts := factsFromCommandOutput(t, resumeResult.stdout, "Next action: ")
			assertSameLines(t, "built-command selected records", resumeFacts.identity, boardFacts.identity)
			assertSameLines(t, "built-command rung evidence", resumeFacts.evidence, boardFacts.evidence)
			if boardFacts.action != resumeFacts.action {
				t.Errorf("built-command next action differs:\n board: %q\nresume: %q", boardFacts.action, resumeFacts.action)
			}
		})
	}
}

// resumeNextFacts runs `savepoint resume` over dir and reads the three facts
// back out of its rendering.
func resumeNextFacts(t *testing.T, dir string) nextFacts {
	t.Helper()

	var out bytes.Buffer
	code, err := runResume(dir, &out)
	if err != nil {
		t.Fatalf("runResume() error = %v", err)
	}
	if code != 0 {
		t.Fatalf("runResume() exit = %d, want 0\n%s", code, out.String())
	}

	return nextFacts{
		identity: linesWithAnyPrefix(out.String(), "Release: ", "Release outcome: ", "Release status: ", "Objective: ", "Task: "),
		evidence: evidenceLinesIn(out.String()),
		action:   strings.TrimPrefix(lineWithPrefixIn(t, out.String(), "Next action: "), "Next action: "),
	}
}

// boardNextFacts runs the board's non-TTY rendering over the same project. It
// is the board's real output path, not a test-only view of it: the plain
// renderer and the drawn panel share one set of Next lines.
func boardNextFacts(t *testing.T, dir string) nextFacts {
	t.Helper()

	var out bytes.Buffer
	if err := boardv2.Run(boardv2.Options{Root: filepath.Join(dir, ".savepoint"), Stdout: &out, TTY: false}); err != nil {
		t.Fatalf("board Run() error = %v", err)
	}

	return nextFacts{
		identity: linesWithAnyPrefix(out.String(), "Release: ", "Release outcome: ", "Release status: ", "Objective: ", "Task: "),
		evidence: evidenceLinesIn(out.String()),
		action:   strings.TrimPrefix(lineWithPrefixIn(t, out.String(), "Action: "), "Action: "),
	}
}

// projectionIdentity is the identity lines a surface must render for want,
// built from the projection itself.
func projectionIdentity(next data.Next) []string {
	var lines []string
	if next.Release != nil {
		lines = append(lines, "Release: "+next.Release.ID+" — "+next.Release.Title)
		lines = append(lines, "Release outcome: "+next.Release.Outcome)
		lines = append(lines, "Release status: "+string(next.Release.Status))
	}
	if next.Objective != nil {
		lines = append(lines, "Objective: "+next.Objective.ID+" — "+next.Objective.Title)
	}
	if next.Task != nil {
		lines = append(lines, "Task: "+next.Task.ID+" — "+next.Task.Title)
	}
	return lines
}

// evidenceLinesIn picks the rung-specific evidence out of a rendering by the
// prefixes internal/resume gives it. Both surfaces print these lines verbatim
// because both call resume.EvidenceLines for them.
func evidenceLinesIn(text string) []string {
	return linesWithAnyPrefix(text,
		"Migration: ", "Replan: ", "Blocked: ", "Completion: ", "Ready: ",
		"Technical clearance: ", "Owner wait: ", "Release readiness: ", "Historical completion: ")
}

func linesWithAnyPrefix(text string, prefixes ...string) []string {
	var found []string
	for _, rawLine := range strings.Split(strings.TrimRight(text, "\n"), "\n") {
		line := strings.TrimSpace(rawLine)
		for _, prefix := range prefixes {
			if strings.HasPrefix(line, prefix) {
				found = append(found, line)
				break
			}
		}
	}
	return found
}

func linesWithTrimmedPrefix(text string, prefixes ...string) []string {
	return linesWithAnyPrefix(text, prefixes...)
}

func lineWithTrimmedPrefix(t *testing.T, text, prefix string) string {
	t.Helper()
	lines := linesWithTrimmedPrefix(text, prefix)
	if len(lines) != 1 {
		t.Fatalf("found %d lines starting with %q, want exactly one:\n%s", len(lines), prefix, text)
	}
	return lines[0]
}

func factsFromCommandOutput(t *testing.T, text, actionPrefix string) nextFacts {
	t.Helper()
	return nextFacts{
		identity: linesWithTrimmedPrefix(text, "Release: ", "Release outcome: ", "Release status: ", "Objective: ", "Task: "),
		evidence: linesWithTrimmedPrefix(text,
			"Migration: ", "Replan: ", "Blocked: ", "Completion: ", "Ready: ",
			"Technical clearance: ", "Owner wait: ", "Release readiness: ", "Historical completion: "),
		action: strings.TrimPrefix(lineWithTrimmedPrefix(t, text, actionPrefix), actionPrefix),
	}
}

func lineWithPrefixIn(t *testing.T, text, prefix string) string {
	t.Helper()
	lines := linesWithAnyPrefix(text, prefix)
	if len(lines) != 1 {
		t.Fatalf("found %d lines starting with %q, want exactly one:\n%s", len(lines), prefix, text)
	}
	return lines[0]
}

func assertSameLines(t *testing.T, what string, want, got []string) {
	t.Helper()
	if len(want) != len(got) {
		t.Fatalf("%s: got %d lines, want %d:\n want: %q\n  got: %q", what, len(got), len(want), want, got)
	}
	for i := range want {
		if want[i] != got[i] {
			t.Errorf("%s line %d differs:\n want: %q\n  got: %q", what, i, want[i], got[i])
		}
	}
}
