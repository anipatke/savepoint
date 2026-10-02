package v2

import (
	"context"
	"errors"
	"fmt"

	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/opencode/savepoint/internal/codehealth"
)

// healthKey opens the Code Health popover from the board.
const healthKey = "H"

// HealthFuncs are the three things the Code Health screen asks of the outside
// world. Production wires them to codehealth in io.go; tests inject fakes so
// the reducer is driven by messages alone (ARCH-02).
type HealthFuncs struct {
	// Load describes the saved results. It reads files only.
	Load func(root string) (codehealth.Dashboard, error)
	// Freshness compares the measured code with the code now.
	Freshness func(ctx context.Context, root string, recorded codehealth.RepositoryIdentity) codehealth.CodeFreshness
	// Refresh runs one manual collection, reporting each signal as it starts.
	Refresh func(ctx context.Context, root string, progress func(codehealth.Progress)) (reportErr, err error)
}

// HealthOverlay is the Code Health popover. It is mutually exclusive
// with every other overlay and restores the board cursor it replaced.
type HealthOverlay struct {
	Origin detailOrigin

	// Dashboard is the last good result and stays on screen across a failed or
	// cancelled refresh and across a load error.
	Dashboard *codehealth.Dashboard
	Loaded    bool
	Freshness *codehealth.CodeFreshness

	Cursor int
	// History shows the last checks in place of the signals.
	History bool
	Notice  string
	// ReportWarning outlives the reload that follows a refresh.
	ReportWarning string
	// Rerun is the official-check command to copy, built when the popover
	// opens so rendering never looks anything up.
	Rerun   string
	Refresh *healthRefresh
}

// healthRefresh is a manual refresh in flight. Cancel is held here so Esc and
// quit can stop it; Events carries progress and the final result from the
// refresh goroutine to the program.
type healthRefresh struct {
	Cancel     context.CancelFunc
	Events     <-chan tea.Msg
	Progress   codehealth.Progress
	Started    bool
	Cancelling bool
}

type healthLoadedMsg struct {
	Dashboard codehealth.Dashboard
	Err       error
}

type healthFreshnessMsg struct {
	SnapshotID string
	Freshness  codehealth.CodeFreshness
}

type healthProgressMsg struct{ Progress codehealth.Progress }

type healthRefreshDoneMsg struct {
	Err error
	// ReportErr is a report that could not be rewritten after a good refresh.
	ReportErr error
}

func (m Model) healthFuncs() HealthFuncs {
	f := m.HealthFuncs
	if f.Load == nil {
		f.Load = loadHealthDashboard
	}
	if f.Freshness == nil {
		f.Freshness = healthFreshness
	}
	if f.Refresh == nil {
		f.Refresh = refreshHealth
	}
	return f
}

func (m *Model) openHealth() tea.Cmd {
	if m.Health != nil || m.Detail != nil || m.Issues != nil || m.ReleaseOverlay || m.Help || !m.Loaded {
		return nil
	}
	m.Health = &HealthOverlay{Rerun: m.healthRerunCommand(), Origin: detailOrigin{
		SidebarFocused:  m.SidebarFocused,
		ObjectiveCursor: m.ObjectiveCursor,
		FocusedColumn:   m.FocusedColumn,
		FocusedCard:     m.FocusedCard,
	}}
	return healthLoadCmd(m.healthFuncs(), m.Root)
}

func (m *Model) closeHealth() {
	if m.Health == nil {
		return
	}
	origin := m.Health.Origin
	m.Health = nil
	m.SidebarFocused = origin.SidebarFocused && m.sidebarVisible()
	m.ObjectiveCursor = origin.ObjectiveCursor
	m.FocusedColumn = origin.FocusedColumn
	m.FocusedCard = origin.FocusedCard
	m.clampObjectiveCursorToRows()
	m.clampFocus()
}

// cancelHealthRefresh stops a refresh in flight. It reports whether one was.
func (m *Model) cancelHealthRefresh() bool {
	if m.Health == nil || m.Health.Refresh == nil {
		return false
	}
	m.Health.Refresh.Cancelling = true
	m.Health.Refresh.Cancel()
	return true
}

// rows is the five signals the popover shows, one each. A signal with several
// configured instances is one row standing for its worst instance, so the frame
// never grows or scrolls.
func (h *HealthOverlay) rows() []codehealth.DashboardRow {
	if h.Dashboard == nil {
		return nil
	}
	var out []codehealth.DashboardRow
	index := map[codehealth.Capability]int{}
	count := map[codehealth.Capability]int{}
	for _, row := range h.Dashboard.Rows {
		count[row.Capability]++
		at, seen := index[row.Capability]
		switch {
		case !seen:
			index[row.Capability] = len(out)
			out = append(out, row)
		case codehealth.WorseInstance(row, out[at]):
			out[at] = row
		}
	}
	for i, row := range out {
		if n := count[row.Capability]; n > 1 {
			out[i] = groupedRow(row, n)
		}
	}
	return out
}

// groupedRow labels the worst of n instances as such.
func groupedRow(row codehealth.DashboardRow, n int) codehealth.DashboardRow {
	name := row.CapabilityText
	if short, ok := healthShortName[row.Capability]; ok && lipgloss.Width(name) > healthRowNameCap {
		name = short
	}
	row.CapabilityText = fmt.Sprintf(textHealthInstances, name, n)
	if row.Name != "" {
		row.Question += fmt.Sprintf(textHealthWorstOf, n, row.Name)
		row.Name = ""
	} else {
		row.Question += fmt.Sprintf(textHealthWorstOfUnnamed, n)
	}
	return row
}

// canRefresh is true once the screen knows the project is configured.
func (h *HealthOverlay) canRefresh() bool {
	return h.Dashboard != nil && h.Dashboard.State != codehealth.DashboardNotConfigured && h.Refresh == nil
}

func (m *Model) handleHealthKey(key string) tea.Cmd {
	h := m.Health
	switch key {
	case "esc":
		switch {
		case m.cancelHealthRefresh():
		case h.History:
			h.History = false
		default:
			m.closeHealth()
		}
	case "h":
		if h.Loaded && h.Dashboard != nil {
			h.History = !h.History
		}
	case "up", "k":
		if !h.History {
			m.moveHealth(-1)
		}
	case "down", "j":
		if !h.History {
			m.moveHealth(1)
		}
	case "R":
		if h.canRefresh() {
			return m.startHealthRefresh()
		}
	}
	return nil
}

func (m *Model) moveHealth(delta int) {
	h := m.Health
	h.Cursor = min(max(h.Cursor+delta, 0), max(len(h.rows())-1, 0))
}

func (m *Model) startHealthRefresh() tea.Cmd {
	ctx, cancel := context.WithCancel(context.Background())
	events := make(chan tea.Msg, 16)
	m.Health.Refresh = &healthRefresh{Cancel: cancel, Events: events}
	m.Health.Notice, m.Health.ReportWarning = "", ""
	return tea.Batch(healthRefreshCmd(ctx, m.healthFuncs(), m.Root, events), waitHealthEvent(events))
}

// applyHealth folds one Code Health message into the overlay. A message that
// arrives after the overlay closed is dropped.
func (m Model) applyHealth(msg tea.Msg) (tea.Model, tea.Cmd) {
	// The header chip follows every good load, whether or not the screen is
	// still open to show it.
	if loaded, ok := msg.(healthLoadedMsg); ok && loaded.Err == nil {
		m.State.HealthChip = loaded.Dashboard.Chip()
	}
	h := m.Health
	if h == nil {
		return m, nil
	}
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case healthLoadedMsg:
		if msg.Err != nil {
			h.Notice = fmt.Sprintf(textHealthLoadFailed, msg.Err)
			break
		}
		h.Notice = ""
		h.Loaded = true
		dashboard := msg.Dashboard
		h.Dashboard = &dashboard
		h.Cursor = min(h.Cursor, max(len(dashboard.Rows)-1, 0))
		if dashboard.State == codehealth.DashboardMeasured {
			h.Freshness = nil
			cmd = healthFreshnessCmd(m.healthFuncs(), m.Root, dashboard)
		}
	case healthFreshnessMsg:
		if h.Dashboard != nil && h.Dashboard.SnapshotID == msg.SnapshotID {
			fresh := msg.Freshness
			h.Freshness = &fresh
		}
	case healthProgressMsg:
		if h.Refresh != nil {
			h.Refresh.Progress, h.Refresh.Started = msg.Progress, true
			cmd = waitHealthEvent(h.Refresh.Events)
		}
	case healthRefreshDoneMsg:
		h.Refresh = nil
		switch {
		case errors.Is(msg.Err, codehealth.ErrCollectionCancelled):
			h.Notice = textHealthCancelled
		case msg.Err != nil:
			h.Notice = fmt.Sprintf(textHealthRefreshFailed, msg.Err)
		default:
			if msg.ReportErr != nil {
				h.ReportWarning = fmt.Sprintf(textHealthReportFailed, msg.ReportErr)
			}
			cmd = healthLoadCmd(m.healthFuncs(), m.Root)
		}
	}
	return m, cmd
}

// healthRefreshHint offers R only when pressing it would start a refresh.
func (m Model) healthRefreshHint() string {
	if m.Health != nil && m.Health.canRefresh() {
		return "R:refresh"
	}
	return ""
}

// healthRerunCommand is the line that re-runs the official check: for the
// Objective the sidebar is on, else the router's, else a placeholder.
func (m Model) healthRerunCommand() string {
	id := "O-###"
	switch {
	case m.ObjectiveCursor >= 0 && m.ObjectiveCursor < len(m.Objectives):
		id = m.Objectives[m.ObjectiveCursor].ID()
	case m.State.Router != nil && m.State.Router.Objective != "":
		id = m.State.Router.Objective
	}
	return fmt.Sprintf(textHealthRerun, id)
}
