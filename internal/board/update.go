package board

import (
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
	case tea.KeyMsg:
		if m.Overlay != OverlayNone {
			return m.updateOverlay(msg)
		}
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "left", "h":
			m.FocusedColumn = prevColumn(m.FocusedColumn)
			m.FocusedTask = 0
			m.StatusMessage = ""
		case "right", "l":
			m.FocusedColumn = nextColumn(m.FocusedColumn)
			m.FocusedTask = 0
			m.StatusMessage = ""
		case "up", "k":
			if m.FocusedTask > 0 {
				m.FocusedTask--
			}
			m.StatusMessage = ""
		case "down", "j":
			if m.FocusedTask < len(m.Tasks[m.FocusedColumn])-1 {
				m.FocusedTask++
			}
			m.StatusMessage = ""
		case "enter":
			tasks := m.Tasks[m.FocusedColumn]
			if len(tasks) > 0 && m.FocusedTask < len(tasks) {
				m.Overlay = OverlayDetail
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
							break
						}
					}
					m.refreshTasks()
				}
			}
		case "backspace":
			tasks := m.Tasks[m.FocusedColumn]
			if len(tasks) > 0 && m.FocusedTask < len(tasks) {
				task := tasks[m.FocusedTask]
				for i, t := range m.AllTasks {
					if t.ID == task.ID {
						Retreat(&m.AllTasks[i])
						break
					}
				}
				m.refreshTasks()
			}
			m.StatusMessage = ""
		case "e":
			m.Overlay = OverlayEpic
			m.EpicCursor = epicIndex(m.Epics, m.SelectedEpic)
		case "r":
			m.Overlay = OverlayRelease
			m.ReleaseCursor = releaseIndex(m.Releases, m.SelectedRelease)
		case "?":
			m.Overlay = OverlayHelp
		}
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	}
	return m, nil
}

func (m Model) updateOverlay(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.Overlay = OverlayNone
	case "up", "k":
		if m.Overlay == OverlayEpic && m.EpicCursor > 0 {
			m.EpicCursor--
		}
		if m.Overlay == OverlayRelease && m.ReleaseCursor > 0 {
			m.ReleaseCursor--
		}
	case "down", "j":
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
			m.refreshTasks()
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
			m.refreshTasks()
			m.Overlay = OverlayNone
			if m.Root != "" {
				if err := m.writeRouterReleaseEpic(); err != nil {
					m.StatusMessage = err.Error()
				}
			}
		}
	}
	return m, nil
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
