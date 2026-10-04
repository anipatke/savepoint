package v2

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	xansi "github.com/charmbracelet/x/ansi"
	"github.com/opencode/savepoint/internal/testutil"
)

const optionsConfig = "# project config\nschema_version: 2\n"

// optionsBoard opens a sized board over a fixture whose config.yml carries a
// comment, so a save that restyles the file is visible.
func optionsBoard(t *testing.T) Model {
	t.Helper()
	root := writeValidProject(t)
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), optionsConfig)
	return openSizedBoard(t, root, 130, 40)
}

func readOptionsConfig(t *testing.T, model Model) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(model.Root, "config.yml"))
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}

// pressOptionsKey sends one key and returns the model and the command the
// reducer asked for, without running it.
func pressOptionsKey(t *testing.T, model Model, key string) (Model, tea.Cmd) {
	t.Helper()
	updated, cmd := model.Update(keyMsg(key))
	next, ok := updated.(Model)
	if !ok {
		t.Fatalf("Update() returned %T, want Model", updated)
	}
	return next, cmd
}

// toggleOptions presses enter on the open screen and settles the save and the
// reload that follows it.
func toggleOptions(t *testing.T, model Model) Model {
	t.Helper()
	next, cmd := pressOptionsKey(t, model, "enter")
	return applyBoardCommands(t, next, cmd)
}

func optionsScreen(model Model) string { return xansi.Strip(model.View()) }

func TestOptionsOpensWithParallelPlanningOffAndExplainsIt(t *testing.T) {
	model := press(t, optionsBoard(t), optionsKey)
	if model.Options == nil {
		t.Fatal("o did not open Advanced Options")
	}
	screen := optionsScreen(model)
	for _, want := range []string{"ADVANCED OPTIONS", "Parallel planning", "OFF", "ignore them", "never block work", "saved now", "arrive later", "yet"} {
		if !strings.Contains(screen, want) {
			t.Errorf("screen omits %q:\n%s", want, screen)
		}
	}
	if got := strings.Count(screen, "Parallel planning:"); got != 1 {
		t.Errorf("screen shows %d options, want exactly the one real option", got)
	}
	if model.hints() != "enter/space:toggle  esc/q:close" {
		t.Errorf("hints = %q", model.hints())
	}
}

func TestOptionsToggleSavesAndSurvivesRestart(t *testing.T) {
	model := press(t, optionsBoard(t), optionsKey)

	on := toggleOptions(t, model)
	screen := optionsScreen(on)
	if !strings.Contains(screen, "ON") || !strings.Contains(screen, "Saved. Parallel planning is on.") {
		t.Errorf("saved state not shown after a successful save:\n%s", screen)
	}
	if got := readOptionsConfig(t, on); got != optionsConfig+"features:\n    parallel_planning: true\n" {
		t.Errorf("config.yml after enabling = %q", got)
	}

	restarted := press(t, openSizedBoard(t, on.Root, 130, 40), optionsKey)
	if !strings.Contains(optionsScreen(restarted), "ON") {
		t.Errorf("restarted board lost the saved preference:\n%s", optionsScreen(restarted))
	}

	off := toggleOptions(t, restarted)
	if !strings.Contains(optionsScreen(off), "Saved. Parallel planning is off.") {
		t.Errorf("turning it off was not confirmed:\n%s", optionsScreen(off))
	}
	if got := readOptionsConfig(t, off); !strings.HasPrefix(got, optionsConfig) || !strings.Contains(got, "parallel_planning: false") {
		t.Errorf("config.yml after disabling = %q", got)
	}
}

func TestOptionsToggleLeavesSavedPlansAndCodeHealthAlone(t *testing.T) {
	model := optionsBoard(t)
	before := snapshotProject(t, model.Root)
	boardBefore := optionsScreen(model)
	chipBefore := model.State.HealthChip

	closed := press(t, toggleOptions(t, press(t, model, optionsKey)), "esc")

	after := snapshotProject(t, model.Root)
	for path, content := range before {
		if filepath.Base(path) == "config.yml" {
			continue
		}
		if after[path] != content {
			t.Errorf("%s changed when the option was toggled", path)
		}
	}
	if closed.State.HealthChip != chipBefore {
		t.Errorf("Code Health chip changed: %+v -> %+v", chipBefore, closed.State.HealthChip)
	}
	if got := optionsScreen(closed); got != boardBefore {
		t.Errorf("board changed after toggling and closing:\nbefore:\n%s\nafter:\n%s", boardBefore, got)
	}
}

func TestOptionsRefreshesAfterAnExternalConfigEdit(t *testing.T) {
	model := press(t, optionsBoard(t), optionsKey)
	testutil.WriteFile(t, filepath.Join(model.Root, "config.yml"), optionsConfig+"features:\n  parallel_planning: true\n")

	updated, _ := model.Update(v2FileChangeMsg{})
	reloaded, _ := updated.(Model).Update(loadCmd(model.Root)())
	if screen := optionsScreen(reloaded.(Model)); !strings.Contains(screen, "ON") {
		t.Errorf("open screen did not follow the external edit:\n%s", screen)
	}
}

func TestOptionsStaleSaveIsRefusedAndExplained(t *testing.T) {
	model := press(t, optionsBoard(t), optionsKey)
	external := optionsConfig + "# edited elsewhere\n"
	testutil.WriteFile(t, filepath.Join(model.Root, "config.yml"), external)

	refused := toggleOptions(t, model)
	screen := optionsScreen(refused)
	if !strings.Contains(screen, "Nothing was saved") || strings.Contains(screen, "Saved. Parallel") {
		t.Errorf("stale save not refused plainly:\n%s", screen)
	}
	if got := readOptionsConfig(t, refused); got != external {
		t.Errorf("stale save overwrote the external edit: %q", got)
	}

	retried := toggleOptions(t, refused)
	if !strings.Contains(optionsScreen(retried), "Saved. Parallel planning is on.") {
		t.Errorf("retry after the reload did not save:\n%s", optionsScreen(retried))
	}
	if got := readOptionsConfig(t, retried); !strings.Contains(got, "# edited elsewhere") {
		t.Errorf("retry lost the external edit: %q", got)
	}
}

func TestOptionsUnwritableConfigReportsAndSavesNothing(t *testing.T) {
	if runtime.GOOS == "windows" || os.Geteuid() == 0 {
		t.Skip("file modes do not stop this writer here")
	}
	model := press(t, optionsBoard(t), optionsKey)
	path := filepath.Join(model.Root, "config.yml")
	if err := os.Chmod(path, 0o444); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(path, 0o644) })

	refused := toggleOptions(t, model)
	screen := optionsScreen(refused)
	if !strings.Contains(screen, "Nothing was saved") || !strings.Contains(screen, "read-only") {
		t.Errorf("unwritable config not explained:\n%s", screen)
	}
	if !strings.Contains(screen, "OFF") {
		t.Errorf("screen shows a state that was not saved:\n%s", screen)
	}
	if got := readOptionsConfig(t, refused); got != optionsConfig {
		t.Errorf("config.yml changed despite the refusal: %q", got)
	}
}

func TestOptionsUnreadableConfigRefusesToggleWithoutWriting(t *testing.T) {
	model := optionsBoard(t)
	model.State.Features = FeatureState{Diagnostic: "failed to parse config YAML"}
	opened := press(t, model, optionsKey)

	refused, cmd := pressOptionsKey(t, opened, "enter")
	if cmd != nil {
		t.Fatal("toggle ran a write against a config.yml that could not be read")
	}
	if screen := optionsScreen(refused); !strings.Contains(screen, "Nothing was saved") || !strings.Contains(screen, "failed to parse") {
		t.Errorf("unreadable config not explained:\n%s", screen)
	}
}

func TestOptionsMalformedPreferenceIsReportedNotHealed(t *testing.T) {
	root := writeValidProject(t)
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\nfeatures:\n  parallel_planning: maybe\n")
	state := loadFeatureState(root)
	if !strings.Contains(state.Diagnostic, "parallel_planning") {
		t.Errorf("Diagnostic = %q, want it to name the key", state.Diagnostic)
	}
}

func TestOptionsCloseReturnsFocusToTheOriginSurface(t *testing.T) {
	sidebar := press(t, sidebarBoard(t, writeNavigationProject(t)), "left")
	if !sidebar.SidebarFocused {
		t.Fatal("fixture did not focus the sidebar")
	}
	for _, closeKey := range []string{"esc", "q", optionsKey} {
		opened := press(t, sidebar, optionsKey)
		closed, cmd := pressOptionsKey(t, opened, closeKey)
		if cmd != nil {
			t.Errorf("%s ran a command instead of only closing", closeKey)
		}
		if closed.Options != nil || !closed.SidebarFocused || closed.ObjectiveCursor != sidebar.ObjectiveCursor {
			t.Errorf("%s did not restore the sidebar focus", closeKey)
		}
	}

	columns := press(t, sidebarBoard(t, writeNavigationProject(t)), "right")
	closed := press(t, columns, optionsKey, "esc")
	if closed.SidebarFocused != columns.SidebarFocused || closed.FocusedColumn != columns.FocusedColumn || closed.FocusedCard != columns.FocusedCard {
		t.Errorf("column focus not restored: %+v vs %+v", closed.detailOriginSnapshot(), columns.detailOriginSnapshot())
	}
}

func (m Model) detailOriginSnapshot() detailOrigin {
	return detailOrigin{SidebarFocused: m.SidebarFocused, ObjectiveCursor: m.ObjectiveCursor, FocusedColumn: m.FocusedColumn, FocusedCard: m.FocusedCard}
}

func TestOptionsCtrlCStillQuits(t *testing.T) {
	_, cmd := pressOptionsKey(t, press(t, optionsBoard(t), optionsKey), "ctrl+c")
	if cmd == nil {
		t.Fatal("ctrl+c did not quit from Advanced Options")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Errorf("ctrl+c returned %T, want tea.QuitMsg", cmd())
	}
}

func TestOptionsSecondToggleWhileSavingIsIgnored(t *testing.T) {
	model := press(t, optionsBoard(t), optionsKey)
	saving, cmd := pressOptionsKey(t, model, "enter")
	if cmd == nil || !saving.Options.Saving {
		t.Fatal("first toggle did not start a save")
	}
	if _, again := pressOptionsKey(t, saving, " "); again != nil {
		t.Error("a second toggle started while the first save was in flight")
	}
}

func TestOptionsKeyIsDocumentedInFooterAndHelp(t *testing.T) {
	model := optionsBoard(t)
	if !strings.Contains(model.hints(), "o:options") {
		t.Errorf("footer omits the options key:\n%q", model.hints())
	}
	if help := optionsScreen(press(t, model, "?")); !strings.Contains(help, "Advanced Options") {
		t.Errorf("help omits Advanced Options:\n%s", help)
	}
}

func TestOptionsFitsTheTerminal(t *testing.T) {
	for _, size := range [][2]int{{100, 30}, {60, 20}, {40, 14}} {
		model := press(t, optionsBoard(t), optionsKey)
		model.Width, model.Height = size[0], size[1]
		for i, line := range strings.Split(model.View(), "\n") {
			if got := xansi.StringWidth(line); got > size[0] {
				t.Errorf("%dx%d: line %d is %d cells wide", size[0], size[1], i, got)
			}
		}
	}
}
