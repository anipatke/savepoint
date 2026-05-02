package board

import (
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/opencode/savepoint/internal/data"
)

var columnOrder = []data.ColumnType{
	data.ColumnPlanned,
	data.ColumnInProgress,
	data.ColumnDone,
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case fileChangeMsg:
		if m.Root != "" {
			return m, reloadTasks(m.Root)
		}
	case reloadMsg:
		m.AllTasks = msg.tasks
		m.Releases = append([]string(nil), msg.releases...)
		m.ReleaseEpics = copyReleaseEpics(msg.releaseEpics)
		m.EpicStatus = msg.epicStatuses
		m.SelectedRelease = firstKnown(m.SelectedRelease, m.Releases)
		m.refreshEpicsForRelease()
		m.refreshTasks()
		m.ensureFocusedTaskVisible()
		if m.Watcher != nil {
			return m, watchFiles(m.Watcher)
		}
	case tea.KeyMsg:
		if m.Overlay != OverlayNone {
			return m.updateOverlay(msg)
		}
		switch msg.String() {
		case "q", "ctrl+c":
			if m.Watcher != nil {
				m.Watcher.Close()
			}
			return m, tea.Quit
		case "e":
			m.Overlay = OverlayEpic
			m.EpicCursor = epicIndex(m.Epics, m.SelectedEpic)
			return m, nil
		case "r":
			m.Overlay = OverlayRelease
			m.ReleaseCursor = releaseIndex(m.Releases, m.SelectedRelease)
			return m, nil
		case "?":
			m.Overlay = OverlayHelp
			return m, nil
		case "m":
			task, ok := m.focusedTask()
			if !ok {
				return m, nil
			}
			if taskDone(task) {
				m.StatusMessage = "Router not updated: focused task is done"
				return m, nil
			}
			if m.Root == "" {
				m.StatusMessage = "Router not updated: no savepoint root"
				return m, nil
			}
			message, err := m.writeRouterTask(task)
			if err != nil {
				m.StatusMessage = err.Error()
				return m, nil
			}
			m.StatusMessage = message
			return m, nil
		}
		if m.EpicPanelFocus {
			if !m.epicPanelAvailable() {
				m.EpicPanelFocus = false
			} else {
				return m.updateEpicPanel(msg)
			}
		}
		switch msg.String() {
		case "left", "h":
			if m.FocusedColumn == data.ColumnPlanned && m.epicPanelAvailable() {
				m.EpicPanelFocus = true
				m.EpicPanelCursor = epicIndex(m.Epics, m.SelectedEpic)
				m.StatusMessage = ""
				return m, nil
			}
			m.FocusedColumn = prevColumn(m.FocusedColumn)
			m.FocusedTask = 0
			m.ensureFocusedTaskVisible()
			m.StatusMessage = ""
		case "right", "l":
			m.FocusedColumn = nextColumn(m.FocusedColumn)
			m.FocusedTask = 0
			m.ensureFocusedTaskVisible()
			m.StatusMessage = ""
		case "up", "k":
			if m.FocusedTask > 0 {
				m.FocusedTask--
			}
			m.ensureFocusedTaskVisible()
			m.StatusMessage = ""
		case "down", "j":
			if m.FocusedTask < len(m.Tasks[m.FocusedColumn])-1 {
				m.FocusedTask++
			}
			m.ensureFocusedTaskVisible()
			m.StatusMessage = ""
		case "pgup":
			m.scrollFocusedColumn(-m.columnPageSize())
			m.StatusMessage = ""
		case "pgdown":
			m.scrollFocusedColumn(m.columnPageSize())
			m.StatusMessage = ""
		case "enter":
			tasks := m.Tasks[m.FocusedColumn]
			if len(tasks) > 0 && m.FocusedTask < len(tasks) {
				m.Overlay = OverlayDetail
				m.DetailOffset = 0
			}
			m.StatusMessage = ""
		case " ":
			tasks := m.Tasks[m.FocusedColumn]
			if len(tasks) > 0 && m.FocusedTask < len(tasks) {
				task := tasks[m.FocusedTask]
				if ok, reason := CanAdvance(&task, m.AllTasks); !ok {
					m.StatusMessage = reason
				} else {
					m.StatusMessage = ""
					for i, t := range m.AllTasks {
						if t.ID == task.ID {
							Advance(&m.AllTasks[i])
							if m.AllTasks[i].Path != "" {
								if err := data.WriteTaskStatus(m.AllTasks[i].Path, &m.AllTasks[i], task.Mtime); err != nil && err != data.ErrMtimeConflict {
									m.StatusMessage = err.Error()
								}
							}
							break
						}
					}
					m.refreshTasks()
					m.ensureFocusedTaskVisible()
				}
			}
		case "backspace":
			tasks := m.Tasks[m.FocusedColumn]
			if len(tasks) > 0 && m.FocusedTask < len(tasks) {
				task := tasks[m.FocusedTask]
				m.StatusMessage = ""
				for i, t := range m.AllTasks {
					if t.ID == task.ID {
						Retreat(&m.AllTasks[i])
						if m.AllTasks[i].Path != "" {
							if err := data.WriteTaskStatus(m.AllTasks[i].Path, &m.AllTasks[i], task.Mtime); err != nil && err != data.ErrMtimeConflict {
								m.StatusMessage = err.Error()
							}
						}
						break
					}
				}
				m.refreshTasks()
				m.ensureFocusedTaskVisible()
			}
		}
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		if !m.epicPanelAvailable() {
			m.EpicPanelFocus = false
		}
		m.ensureFocusedTaskVisible()
	}
	return m, nil
}

func (m Model) updateEpicPanel(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if len(m.Epics) == 0 {
		m.EpicPanelFocus = false
		return m, nil
	}

	m.StatusMessage = ""
	switch msg.String() {
	case "up", "k":
		if m.EpicPanelCursor > 0 {
			m.EpicPanelCursor--
			m.selectEpicPanelEpic()
		}
	case "down", "j":
		if m.EpicPanelCursor < len(m.Epics)-1 {
			m.EpicPanelCursor++
			m.selectEpicPanelEpic()
		}
	case "enter":
		m.openEpicDetailOverlay()
	case "right", "l":
		m.EpicPanelFocus = false
		m.FocusedColumn = data.ColumnPlanned
		m.FocusedTask = 0
		m.ensureFocusedTaskVisible()
	case "left", "h":
		// Stay in the panel; there is no focusable region further left.
	}
	return m, nil
}

func (m *Model) selectEpicPanelEpic() {
	if len(m.Epics) == 0 || m.EpicPanelCursor < 0 || m.EpicPanelCursor >= len(m.Epics) {
		return
	}
	m.SelectedEpic = m.Epics[m.EpicPanelCursor]
	m.FocusedTask = 0
	m.DetailOffset = 0
	m.refreshTasks()
	m.ensureFocusedTaskVisible()
	if m.Root != "" {
		if err := m.writeRouterReleaseEpic(); err != nil {
			m.StatusMessage = err.Error()
		}
	}
}

func (m *Model) openEpicDetailOverlay() {
	if len(m.Epics) == 0 || m.EpicPanelCursor < 0 || m.EpicPanelCursor >= len(m.Epics) {
		return
	}
	epicSlug := m.Epics[m.EpicPanelCursor]
	shortEpicID := epicSlug
	if idx := strings.Index(epicSlug, "-"); idx >= 0 {
		shortEpicID = epicSlug[:idx]
	}
	epicDir := filepath.Join(m.Root, "releases", m.SelectedRelease, "epics", epicSlug)
	m.EpicDetailContent = readEpicDetailFile(epicDir, shortEpicID)
	m.EpicDetailOffset = 0
	m.Overlay = OverlayEpicDetail
}

// readEpicDetailFile finds and reads the canonical detail file for an epic.
// Tries E##-Detail.md then E##-Design.md; falls back to any E##-*.md in the dir.
func readEpicDetailFile(epicDir, shortID string) string {
	for _, suffix := range []string{"-Detail.md", "-Design.md"} {
		if raw, err := os.ReadFile(filepath.Join(epicDir, shortID+suffix)); err == nil {
			return string(raw)
		}
	}
	entries, err := os.ReadDir(epicDir)
	if err != nil {
		return "(no detail available)"
	}
	prefix := shortID + "-"
	for _, e := range entries {
		if !e.IsDir() && strings.HasPrefix(e.Name(), prefix) && strings.HasSuffix(e.Name(), ".md") {
			if raw, err := os.ReadFile(filepath.Join(epicDir, e.Name())); err == nil {
				return string(raw)
			}
		}
	}
	return "(no detail available)"
}

func copyReleaseEpics(in map[string][]string) map[string][]string {
	out := make(map[string][]string, len(in))
	for release, epics := range in {
		out[release] = append([]string(nil), epics...)
	}
	return out
}

func (m Model) updateOverlay(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.Overlay = OverlayNone
	case "up", "k":
		if m.Overlay == OverlayDetail && m.DetailOffset > 0 {
			m.DetailOffset--
		}
		if m.Overlay == OverlayEpicDetail && m.EpicDetailOffset > 0 {
			m.EpicDetailOffset--
		}
		if m.Overlay == OverlayEpic && m.EpicCursor > 0 {
			m.EpicCursor--
		}
		if m.Overlay == OverlayRelease && m.ReleaseCursor > 0 {
			m.ReleaseCursor--
		}
	case "down", "j":
		if m.Overlay == OverlayDetail {
			m.DetailOffset++
		}
		if m.Overlay == OverlayEpicDetail {
			m.EpicDetailOffset++
		}
		if m.Overlay == OverlayEpic && len(m.Epics) > 0 && m.EpicCursor < len(m.Epics)-1 {
			m.EpicCursor++
		}
		if m.Overlay == OverlayRelease && len(m.Releases) > 0 && m.ReleaseCursor < len(m.Releases)-1 {
			m.ReleaseCursor++
		}
	case "enter":
		if m.Overlay == OverlayEpic && len(m.Epics) > 0 {
			m.SelectedEpic = m.Epics[m.EpicCursor]
			m.FocusedTask = 0
			m.DetailOffset = 0
			m.refreshTasks()
			m.ensureFocusedTaskVisible()
			m.Overlay = OverlayNone
			if m.Root != "" {
				if err := m.writeRouterReleaseEpic(); err != nil {
					m.StatusMessage = err.Error()
				}
			}
		}
		if m.Overlay == OverlayRelease && len(m.Releases) > 0 {
			m.SelectedRelease = m.Releases[m.ReleaseCursor]
			m.refreshEpicsForRelease()
			m.FocusedTask = 0
			m.DetailOffset = 0
			m.refreshTasks()
			m.ensureFocusedTaskVisible()
			m.Overlay = OverlayNone
			if m.Root != "" {
				if err := m.writeRouterReleaseEpic(); err != nil {
					m.StatusMessage = err.Error()
				}
			}
		}
	case "pgup":
		if m.Overlay == OverlayDetail {
			m.DetailOffset -= m.detailPageSize()
			if m.DetailOffset < 0 {
				m.DetailOffset = 0
			}
		}
		if m.Overlay == OverlayEpicDetail {
			m.EpicDetailOffset -= m.detailPageSize()
			if m.EpicDetailOffset < 0 {
				m.EpicDetailOffset = 0
			}
		}
	case "pgdown":
		if m.Overlay == OverlayDetail {
			m.DetailOffset += m.detailPageSize()
		}
		if m.Overlay == OverlayEpicDetail {
			m.EpicDetailOffset += m.detailPageSize()
		}
	}
	return m, nil
}

func (m *Model) ensureFocusedTaskVisible() {
	if m.ColumnOffsets == nil {
		m.ColumnOffsets = newColumnOffsets()
	}
	tasks := m.Tasks[m.FocusedColumn]
	if len(tasks) == 0 {
		m.ColumnOffsets[m.FocusedColumn] = 0
		return
	}
	pageSize := m.columnPageSize()
	offset := m.ColumnOffsets[m.FocusedColumn]
	if m.FocusedTask < offset {
		offset = m.FocusedTask
	}
	if m.FocusedTask >= offset+pageSize {
		offset = m.FocusedTask - pageSize + 1
	}
	maxOffset := max(len(tasks)-pageSize, 0)
	if offset > maxOffset {
		offset = maxOffset
	}
	if offset < 0 {
		offset = 0
	}
	m.ColumnOffsets[m.FocusedColumn] = offset
}

func (m *Model) scrollFocusedColumn(delta int) {
	if m.ColumnOffsets == nil {
		m.ColumnOffsets = newColumnOffsets()
	}
	tasks := m.Tasks[m.FocusedColumn]
	if len(tasks) == 0 {
		m.ColumnOffsets[m.FocusedColumn] = 0
		m.FocusedTask = 0
		return
	}
	pageSize := m.columnPageSize()
	maxOffset := max(len(tasks)-pageSize, 0)
	offset := m.ColumnOffsets[m.FocusedColumn] + delta
	if offset < 0 {
		offset = 0
	}
	if offset > maxOffset {
		offset = maxOffset
	}
	m.ColumnOffsets[m.FocusedColumn] = offset
	m.FocusedTask = min(offset, len(tasks)-1)
}

func (m Model) columnPageSize() int {
	h := m.Height
	if h == 0 {
		h = defaultTermH
	}
	return visibleColumnTaskLimit(CalculateLayout(m.Width, h).ContentHeight)
}

func (m Model) detailPageSize() int {
	return max(detailMaxHeight(m.Height)-3, 1)
}

func (m Model) epicPanelPageSize() int {
	h := m.Height
	if h == 0 {
		h = defaultTermH
	}
	return max(h/2, 1)
}

func (m Model) epicPanelAvailable() bool {
	return len(m.Epics) > 0 && CalculateLayout(m.Width, m.Height).EpicPanelVisible
}

func prevColumn(col data.ColumnType) data.ColumnType {
	for i, c := range columnOrder {
		if c == col {
			return columnOrder[(i+len(columnOrder)-1)%len(columnOrder)]
		}
	}
	return columnOrder[0]
}

func nextColumn(col data.ColumnType) data.ColumnType {
	for i, c := range columnOrder {
		if c == col {
			return columnOrder[(i+1)%len(columnOrder)]
		}
	}
	return columnOrder[0]
}
