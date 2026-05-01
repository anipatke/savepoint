package board

import (
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/opencode/savepoint/internal/data"
)

type OverlayType string

const (
	OverlayNone    OverlayType = ""
	OverlayHelp    OverlayType = "help"
	OverlayEpic    OverlayType = "epic"
	OverlayRelease OverlayType = "release"
	OverlayDetail  OverlayType = "detail"
)

// Model holds all board state. Tasks are grouped by column for O(1) column access.
type Model struct {
	AllTasks        []data.Task
	Tasks           map[data.ColumnType][]data.Task
	FocusedColumn   data.ColumnType
	FocusedTask     int
	SelectedEpic    string
	SelectedRelease string
	Epics           []string
	EpicCursor      int
	Releases        []string
	ReleaseEpics    map[string][]string
	ReleaseCursor   int
	Overlay         OverlayType
	Width           int
	Height          int
	StatusMessage   string
	Root            string
}

// NewModel groups tasks by column and returns an initialized Model.
func NewModel(tasks []data.Task, release, epic string) Model {
	m := Model{
		AllTasks:        append([]data.Task(nil), tasks...),
		FocusedColumn:   data.ColumnPlanned,
		FocusedTask:     0,
		SelectedEpic:    epic,
		SelectedRelease: release,
		Overlay:         OverlayNone,
	}
	m.refreshTasks()
	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch()
}

func groupedTasks(tasks []data.Task) map[data.ColumnType][]data.Task {
	grouped := map[data.ColumnType][]data.Task{
		data.ColumnPlanned:    {},
		data.ColumnInProgress: {},
		data.ColumnDone:       {},
	}
	for _, t := range tasks {
		col := t.Column
		if col == "" {
			col = data.ColumnPlanned
		}
		grouped[col] = append(grouped[col], t)
	}
	return grouped
}

func (m *Model) refreshTasks() {
	visible := make([]data.Task, 0, len(m.AllTasks))
	for _, t := range m.AllTasks {
		if m.SelectedRelease != "" && t.Release != "" && t.Release != m.SelectedRelease {
			continue
		}
		if m.SelectedEpic != "" && t.Epic != "" && t.Epic != m.SelectedEpic {
			continue
		}
		visible = append(visible, t)
	}
	m.Tasks = groupedTasks(visible)
	m.clampFocusedTask()
}

func (m *Model) refreshEpicsForRelease() {
	if len(m.ReleaseEpics) == 0 {
		return
	}

	epics := m.ReleaseEpics[m.SelectedRelease]
	m.Epics = append([]string(nil), epics...)
	if len(m.Epics) == 0 {
		m.SelectedEpic = ""
		m.EpicCursor = 0
		return
	}

	for _, epic := range m.Epics {
		if epic == m.SelectedEpic {
			m.EpicCursor = epicIndex(m.Epics, m.SelectedEpic)
			return
		}
	}

	m.SelectedEpic = m.Epics[0]
	m.EpicCursor = 0
}

func (m *Model) clampFocusedTask() {
	tasks := m.Tasks[m.FocusedColumn]
	if len(tasks) == 0 {
		m.FocusedTask = 0
		return
	}
	if m.FocusedTask >= len(tasks) {
		m.FocusedTask = len(tasks) - 1
	}
	if m.FocusedTask < 0 {
		m.FocusedTask = 0
	}
}

func (m *Model) writeRouterReleaseEpic() error {
	routerPath := filepath.Join(m.Root, "router.md")

	fi, err := os.Stat(routerPath)
	if err != nil {
		return err
	}

	content, err := os.ReadFile(routerPath)
	if err != nil {
		return err
	}

	r := data.NewRouterReader()
	state, err := r.ReadState(string(content))
	if err != nil {
		return err
	}

	state.Epic = m.SelectedEpic
	state.Release = m.SelectedRelease

	return data.WriteRouterState(m.Root, state, fi.ModTime())
}
