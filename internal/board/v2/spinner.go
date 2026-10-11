package v2

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/opencode/savepoint/internal/data"
)

// The Task tile spinner is the owner-approved exception to the visual
// identity's "running animation skip" rule (.savepoint/visual-identity.md).
// It is driven by a frame counter on Model and a tick command that stops by
// itself, so an idle board runs no timer (ARCH-02).

// spinnerInterval is the time between spinner frames.
const spinnerInterval = 120 * time.Millisecond

// spinnerFrames are single-cell Braille frames, drawn in order and repeated.
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// liveStageWords are the words a live tile shows for each animating stage.
var liveStageWords = map[data.ProgressStage]string{
	data.StageBuild: "Building…",
	data.StageTest:  "Testing…",
}

// spinnerTickMsg advances the spinner by one frame.
type spinnerTickMsg struct{}

// spinnerFrame is the glyph for a frame counter, wrapping around the set.
func spinnerFrame(frame int) string {
	n := len(spinnerFrames)
	return spinnerFrames[((frame%n)+n)%n]
}

// spinnerTickCmd waits one interval and then reports a tick. It reads no file.
func spinnerTickCmd() tea.Cmd {
	return tea.Tick(spinnerInterval, func(time.Time) tea.Msg { return spinnerTickMsg{} })
}

// hasLiveTile reports whether any card in the board's current scope is at
// build or test, which is the only reason the spinner has to move.
func (m Model) hasLiveTile() bool {
	for _, cards := range m.Cards {
		for _, card := range cards {
			if _, ok := liveStageBadge(card.Task.Status, card.Task.Stage, 0); ok {
				return true
			}
		}
	}
	return false
}

// spinnerCmd starts the tick when a live tile is shown and the timer is not
// already running. It returns nil otherwise.
func (m *Model) spinnerCmd() tea.Cmd {
	if m.spinnerRunning || !m.hasLiveTile() {
		return nil
	}
	m.spinnerRunning = true
	return spinnerTickCmd()
}

// applySpinnerTick advances the frame and schedules the next tick only while
// a live tile is still shown; otherwise the timer ends.
func (m Model) applySpinnerTick() (tea.Model, tea.Cmd) {
	m.spinnerRunning = false
	if !m.hasLiveTile() {
		return m, nil
	}
	m.spinnerFrame++
	m.spinnerRunning = true
	return m, spinnerTickCmd()
}
