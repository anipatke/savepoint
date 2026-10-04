package v2

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/styles"
)

const (
	// optionsKey opens Advanced Options from the board.
	optionsKey = "o"

	optionsTitle           = "ADVANCED OPTIONS (EXPERIMENTAL)"
	parallelPlanningLabel  = "Parallel planning"
	parallelPlanningDetail = "Saves your choice for optional suggestions about which Tasks could run side by side in separate worktrees. Only the choice is saved now; no suggestions appear yet. They arrive later. You can ignore them: they never block work, change Code Health, or decide what is done."
	optionsStorageNote     = "Saved in .savepoint/config.yml."
)

// FeatureState is the saved optional-feature choice one load read from
// config.yml. Source names the bytes it came from, so a toggle is refused if
// the file changed since. Diagnostic is set instead when config.yml could not
// be read; the board still loads, and the option explains why it cannot change.
type FeatureState struct {
	ParallelPlanning bool
	Source           data.FeatureSource
	Diagnostic       string
}

// loadFeatureState reads the saved preference. It is called by loadProject
// only, so rendering and Update never touch config.yml (ARCH-02).
func loadFeatureState(root string) FeatureState {
	config, err := data.NewConfigReader().Read(filepath.Join(root, "config.yml"))
	if err != nil {
		return FeatureState{Diagnostic: err.Error()}
	}
	return FeatureState{
		ParallelPlanning: config.ParallelPlanningEnabled(),
		Source:           config.FeatureSource(),
	}
}

// OptionsOverlay is the open Advanced Options screen. Its origin restores the
// surface that had focus when o opened it. Saving blocks a second toggle while
// a write is in flight, so two presses cannot race on one source identity.
type OptionsOverlay struct {
	Origin detailOrigin
	Saving bool
	// Notice is the last save result shown inside the screen; NoticeFailed
	// marks it as a refusal rather than a confirmation.
	Notice       string
	NoticeFailed bool
}

// optionsSavedMsg is the typed result of one toggle write. Features is the
// state read back after a successful write, so "saved" is only ever shown for
// what the file now says.
type optionsSavedMsg struct {
	enabled  bool
	features FeatureState
	err      error
}

// writeParallelPlanningCmd saves the preference through the data writer, which
// owns byte preservation and the stale-source refusal.
func writeParallelPlanningCmd(root string, enabled bool, source data.FeatureSource) tea.Cmd {
	return func() tea.Msg {
		path := filepath.Join(root, "config.yml")
		if err := data.WriteParallelPlanning(path, enabled, source); err != nil {
			return optionsSavedMsg{enabled: enabled, err: err}
		}
		return optionsSavedMsg{enabled: enabled, features: loadFeatureState(root)}
	}
}

func (m *Model) openOptions() {
	m.Options = &OptionsOverlay{Origin: detailOrigin{
		SidebarFocused:  m.SidebarFocused,
		ObjectiveCursor: m.ObjectiveCursor,
		FocusedColumn:   m.FocusedColumn,
		FocusedCard:     m.FocusedCard,
	}}
}

// closeOptions returns the keys to the surface the screen was opened from,
// clamped into whatever a reload underneath it left behind.
func (m *Model) closeOptions() {
	if m.Options == nil {
		return
	}
	origin := m.Options.Origin
	m.Options = nil
	m.SidebarFocused = origin.SidebarFocused && m.sidebarVisible()
	m.ObjectiveCursor = origin.ObjectiveCursor
	m.FocusedColumn = origin.FocusedColumn
	m.FocusedCard = origin.FocusedCard
	m.clampObjectiveCursorToRows()
	m.clampFocus()
}

// handleOptionsKey owns the keys while the screen is open. q closes it rather
// than quitting, as in the Goal selector.
func (m Model) handleOptionsKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "ctrl+c":
		return m.quit()
	case "esc", "q", optionsKey:
		m.closeOptions()
	case "enter", " ":
		return m.toggleParallelPlanning()
	}
	return m, nil
}

func (m Model) toggleParallelPlanning() (tea.Model, tea.Cmd) {
	options := m.Options
	if options.Saving {
		return m, nil
	}
	if m.Root == "" {
		return m, nil
	}
	if diagnostic := m.State.Features.Diagnostic; diagnostic != "" {
		options.Notice = "Nothing was saved. config.yml could not be read: " + diagnostic
		options.NoticeFailed = true
		return m, nil
	}
	options.Saving = true
	options.Notice = ""
	options.NoticeFailed = false
	return m, writeParallelPlanningCmd(m.Root, !m.State.Features.ParallelPlanning, m.State.Features.Source)
}

// applyOptionsSaved folds one write result into the model. Every outcome
// starts a reload: a conflict then shows the file's current value and a fresh
// source to retry against.
func (m Model) applyOptionsSaved(msg optionsSavedMsg) (tea.Model, tea.Cmd) {
	notice, failed := optionsSaveNotice(msg)
	if !failed {
		m.State.Features = msg.features
	}
	if m.Options != nil {
		m.Options.Saving = false
		m.Options.Notice = notice
		m.Options.NoticeFailed = failed
	} else {
		m.StatusMessage = notice
		m.preserveReloadStatus = true
	}
	return m, loadCmd(m.Root)
}

// optionsSaveNotice words the result for the owner. A success is confirmed
// only when the state read back matches what was asked for.
func optionsSaveNotice(msg optionsSavedMsg) (notice string, failed bool) {
	switch {
	case msg.err == nil && msg.features.Diagnostic == "" && msg.features.ParallelPlanning == msg.enabled:
		return "Saved. " + parallelPlanningLabel + " is " + onOff(msg.enabled) + ".", false
	case msg.err == nil:
		return "Nothing was confirmed. config.yml does not show the change; reload and try again.", true
	case errors.Is(msg.err, data.ErrV2SourceConflict), errors.Is(msg.err, data.ErrMtimeConflict):
		return "Nothing was saved. config.yml changed on disk; the screen now shows its current value. Press enter to try again.", true
	case errors.Is(msg.err, os.ErrPermission):
		return "Nothing was saved. config.yml is read-only; make it writable and try again.", true
	case errors.Is(msg.err, data.ErrMalformedFeaturePreference):
		return "Nothing was saved. " + msg.err.Error() + ". Fix it by hand in config.yml.", true
	default:
		return "Nothing was saved. " + msg.err.Error(), true
	}
}

func onOff(enabled bool) string {
	if enabled {
		return "on"
	}
	return "off"
}

func (m Model) renderOptionsOverlay(base string, width, height int) string {
	screen := renderOptions(m.State.Features, *m.Options, releaseOverlayWidth(width))
	return overlayOnV2Base(dimV2Lines(base), screen, width, height)
}

// renderOptions draws the one real option with its state in words as well as
// color, the explanation, and the last save result.
func renderOptions(features FeatureState, options OptionsOverlay, width int) string {
	inner := width - 4
	if inner < 2 {
		inner = 2
	}
	wrap := lipgloss.NewStyle().Width(inner)

	state := styles.CardMeta.Render("OFF")
	if features.ParallelPlanning {
		state = styles.HealthGood.Render("ON")
	}
	lines := []string{
		styles.ColumnTitleFocused.Render(optionsTitle),
		strings.Repeat("─", inner),
		styles.TaskItemFocused.Render(releaseActiveMarker+" "+parallelPlanningLabel+": ") + state,
		"",
		wrap.Render(styles.CardMeta.Render(parallelPlanningDetail)),
		wrap.Render(styles.CardMeta.Render(optionsStorageNote)),
	}
	switch {
	case options.Saving:
		lines = append(lines, "", wrap.Render(styles.StatusBar.Render("Saving…")))
	case options.Notice != "" && options.NoticeFailed:
		lines = append(lines, "", wrap.Render(styles.HealthNeedsAttention.Render(options.Notice)))
	case options.Notice != "":
		lines = append(lines, "", wrap.Render(styles.HealthGood.Render(options.Notice)))
	}
	lines = append(lines, "", styles.CardMeta.Render("enter/space:toggle  esc/q:close"))
	return styles.DetailOverlay.Width(width).Render(strings.Join(lines, "\n"))
}
