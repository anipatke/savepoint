package board

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fsnotify/fsnotify"
	"github.com/opencode/savepoint/internal/data"
)

type OverlayType string

const (
	OverlayNone       OverlayType = ""
	OverlayHelp       OverlayType = "help"
	OverlayEpic       OverlayType = "epic"
	OverlayRelease    OverlayType = "release"
	OverlayDetail     OverlayType = "detail"
	OverlayEpicDetail OverlayType = "detail-epic"
)

// Model holds all board state. Tasks are grouped by column for O(1) column access.
type Model struct {
	Theme             data.Theme
	AllTasks          []data.Task
	Tasks             map[data.ColumnType][]data.Task
	FocusedColumn     data.ColumnType
	FocusedTask       int
	ColumnOffsets     map[data.ColumnType]int
	DetailOffset      int
	SelectedEpic      string
	SelectedRelease   string
	Epics             []string
	EpicCursor        int
	EpicPanelFocus    bool
	EpicPanelCursor   int
	EpicDetailOffset  int
	EpicDetailEpic    string
	EpicDetailContent string
	EpicDetailTab     int    // 0=Detail, 1=Audit
	EpicAuditContent  string // cached E##-Audit.md content
	Releases          []string
	ReleaseEpics      map[string][]string
	ReleaseCursor     int
	Overlay           OverlayType
	Width             int
	Height            int
	StatusMessage     string
	Root              string
	EpicStatus        map[string]string
	RouterTask        string
	RouterState       *data.RouterState
	Watcher           *fsnotify.Watcher
}

// NewModel groups tasks by column and returns an initialized Model.
func NewModel(tasks []data.Task, release, epic string) Model {
	m := Model{
		AllTasks:        append([]data.Task(nil), tasks...),
		FocusedColumn:   data.ColumnPlanned,
		FocusedTask:     0,
		ColumnOffsets:   newColumnOffsets(),
		SelectedEpic:    epic,
		SelectedRelease: release,
		Overlay:         OverlayNone,
	}
	m.refreshTasks()
	return m
}

func (m Model) Init() tea.Cmd {
	if m.Watcher == nil {
		return nil
	}
	return watchFiles(m.Watcher)
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
	m.clampColumnOffsets()
}

func newColumnOffsets() map[data.ColumnType]int {
	return map[data.ColumnType]int{
		data.ColumnPlanned:    0,
		data.ColumnInProgress: 0,
		data.ColumnDone:       0,
	}
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
		m.EpicPanelCursor = 0
		m.EpicPanelFocus = false
		return
	}

	for _, epic := range m.Epics {
		if epic == m.SelectedEpic {
			m.EpicCursor = sliceIndex(m.Epics, m.SelectedEpic)
			m.clampEpicPanelCursor()
			return
		}
	}

	m.SelectedEpic = m.Epics[0]
	m.EpicCursor = 0
	m.clampEpicPanelCursor()
}

func (m *Model) clampEpicPanelCursor() {
	if len(m.Epics) == 0 {
		m.EpicPanelCursor = 0
		m.EpicPanelFocus = false
		return
	}
	if m.EpicPanelCursor >= len(m.Epics) {
		m.EpicPanelCursor = len(m.Epics) - 1
	}
	if m.EpicPanelCursor < 0 {
		m.EpicPanelCursor = 0
	}
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

func (m *Model) clampColumnOffsets() {
	if m.ColumnOffsets == nil {
		m.ColumnOffsets = newColumnOffsets()
	}
	for _, col := range columnOrder {
		tasks := m.Tasks[col]
		offset := m.ColumnOffsets[col]
		if offset < 0 || len(tasks) == 0 {
			m.ColumnOffsets[col] = 0
			continue
		}
		if offset >= len(tasks) {
			m.ColumnOffsets[col] = len(tasks) - 1
		}
	}
}

func taskDone(task data.Task) bool {
	return task.Column == data.ColumnDone || task.Status == string(data.StatusDone)
}
