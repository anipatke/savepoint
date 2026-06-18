package board

import (
	"os"
	"path/filepath"
	"strings"
	"time"

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
		debugf("dispatch: fileChangeMsg")
		if m.Root != "" {
			return m, reloadTasks(m.Root, m.Dependencies)
		}
	case reloadMsg:
		debugf("dispatch: reloadMsg tasks=%d releases=%d", len(msg.tasks), len(msg.releases))
		m.AllTasks = msg.tasks
		m.AllDefects = msg.defects
		m.Releases = append([]string(nil), msg.releases...)
		m.ReleaseEpics = copyReleaseEpics(msg.releaseEpics)
		m.EpicStatus = msg.epicStatuses
		if msg.routerState != nil {
			m.RouterState = msg.routerState
			m.RouterTask = msg.routerState.Task
		}
		m.StatusMessage = msg.message
		m.SelectedRelease = firstKnown(m.SelectedRelease, m.Releases)
		m.refreshEpicsForRelease()
		m.refreshTasks()
		m.ensureFocusedTaskVisible()
		if m.Watcher != nil {
			return m, watchFiles(m.Watcher)
		}
	case tea.KeyMsg:
		if m.Overlay != OverlayNone {
			return m.handleOverlay(msg)
		}
		return m.handleBoardKey(msg)
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		if !m.epicPanelAvailable() {
			m.EpicPanelFocus = false
		}
		m.ensureFocusedTaskVisible()
	case routerWriteMsg:
		m.StatusMessage = msg.message
		m.RouterState = msg.state
		m.RouterTask = msg.taskID
	case taskWriteMsg:
		for i, t := range m.AllTasks {
			if sameTaskRecord(t, msg.next) {
				m.AllTasks[i] = msg.next
				break
			}
		}
		m.StatusMessage = taskTransitionMessage(msg.prefix, msg.next)
		m.refreshTasks()
		m.ensureFocusedTaskVisible()
	case taskRefreshMsg:
		for i, t := range m.AllTasks {
			if sameTaskRecord(t, msg.task) {
				m.AllTasks[i] = msg.task
				break
			}
		}
		m.StatusMessage = msg.message
		m.refreshTasks()
		m.ensureFocusedTaskVisible()
	case defectWriteMsg:
		for i, d := range m.AllDefects {
			if sameDefectRecord(d, msg.next) {
				m.AllDefects[i] = msg.next
				break
			}
		}
		m.StatusMessage = "Resolved " + shortID(msg.next.ID)
		m.clampDefectCursor()
	case epicDetailMsg:
		m.EpicDetailContent = msg.content
	case auditContentMsg:
		m.EpicAuditContent = msg.content
	case releaseDocsMsg:
		m.ReleaseDocs = msg.docs
		m.clampReleaseDocIndex()
	case epicStatusWrittenMsg:
		if m.EpicStatus == nil {
			m.EpicStatus = map[string]string{}
		}
		m.EpicStatus[msg.epicID] = msg.status
		m.StatusMessage = "Marked " + shortID(msg.epicID) + " " + msg.status
	case errorMsg:
		m.StatusMessage = msg.message
	}
	return m, nil
}

func (m Model) handleOverlay(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.Overlay {
	case OverlayHelp:
		return m.handleHelpOverlay(msg)
	case OverlayEpic:
		return m.handleEpicOverlay(msg)
	case OverlayRelease:
		return m.handleReleaseOverlay(msg)
	case OverlayDetail:
		return m.handleDetailOverlay(msg)
	case OverlayEpicDetail:
		return m.handleEpicDetailOverlay(msg)
	case OverlayDefect:
		return m.handleDefectOverlay(msg)
	case OverlayDefectDetail:
		return m.handleDefectDetailOverlay(msg)
	}
	return m, nil
}

func (m Model) handleBoardKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		if m.Watcher != nil {
			m.Watcher.Close()
		}
		return m, tea.Quit
	case "e":
		m.Overlay = OverlayEpic
		m.EpicCursor = sliceIndex(m.Epics, m.SelectedEpic)
		return m, nil
	case "r":
		m.Overlay = OverlayRelease
		m.ReleaseCursor = sliceIndex(m.Releases, m.SelectedRelease)
		return m, nil
	case "d":
		m.Overlay = OverlayDefect
		m.DefectCursor = 0
		return m, nil
	case "?":
		m.Overlay = OverlayHelp
		return m, nil
	case "ctrl+r":
		if m.Root == "" {
			m.StatusMessage = "Refresh failed: no savepoint root"
			return m, nil
		}
		m.StatusMessage = ""
		return m, reloadTasksWithMessage(m.Root, m.Dependencies, "Refreshed")
	case "p":
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
		return m, writeRouterTaskCmd(m.Root, task, m.Dependencies.RouterReader)
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
			m.EpicPanelCursor = sliceIndex(m.Epics, m.SelectedEpic)
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
		return m.handleAdvanceTask()
	case "backspace":
		return m.handleRetreatTask()
	}
	return m, nil
}

func (m Model) handleAdvanceTask() (tea.Model, tea.Cmd) {
	tasks := m.Tasks[m.FocusedColumn]
	if len(tasks) > 0 && m.FocusedTask < len(tasks) {
		task := tasks[m.FocusedTask]
		if ok, reason := CanAdvance(&task, m.AllTasks, m.selectedReleaseEpicStatuses()); !ok {
			m.StatusMessage = reason
			return m, nil
		}
		m.StatusMessage = ""
		for i, t := range m.AllTasks {
			if sameTaskRecord(t, task) {
				next := m.AllTasks[i]
				if err := Advance(&next); err != nil {
					m.StatusMessage = err.Error()
					return m, nil
				}
				if next.Path != "" {
					return m, writeTaskStatusCmd(t, next, task.Mtime, "Moved")
				}
				m.AllTasks[i] = next
				m.StatusMessage = taskTransitionMessage("Moved", next)
				break
			}
		}
		m.refreshTasks()
		m.ensureFocusedTaskVisible()
	}
	return m, nil
}

func (m Model) selectedReleaseEpicStatuses() map[string]string {
	if len(m.ReleaseEpics) == 0 {
		return m.EpicStatus
	}
	epics := m.ReleaseEpics[m.SelectedRelease]
	statuses := make(map[string]string, len(epics))
	for _, epic := range epics {
		if status, ok := m.EpicStatus[epic]; ok {
			statuses[epic] = status
		}
	}
	return statuses
}

func (m Model) handleRetreatTask() (tea.Model, tea.Cmd) {
	tasks := m.Tasks[m.FocusedColumn]
	if len(tasks) > 0 && m.FocusedTask < len(tasks) {
		task := tasks[m.FocusedTask]
		m.StatusMessage = ""
		for i, t := range m.AllTasks {
			if sameTaskRecord(t, task) {
				next := m.AllTasks[i]
				if err := Retreat(&next); err != nil {
					m.StatusMessage = err.Error()
					return m, nil
				}
				if next.Path != "" {
					return m, writeTaskStatusCmd(t, next, task.Mtime, "Moved back")
				}
				m.AllTasks[i] = next
				m.StatusMessage = taskTransitionMessage("Moved back", next)
				break
			}
		}
		m.refreshTasks()
		m.ensureFocusedTaskVisible()
	}
	return m, nil
}

func (m Model) handleHelpOverlay(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.Overlay = OverlayNone
	}
	return m, nil
}

func (m Model) handleEpicOverlay(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.Overlay = OverlayNone
	case "up", "k":
		if m.EpicCursor > 0 {
			m.EpicCursor--
		}
	case "down", "j":
		if len(m.Epics) > 0 && m.EpicCursor < len(m.Epics)-1 {
			m.EpicCursor++
		}
	case "enter":
		if len(m.Epics) > 0 {
			m.SelectedEpic = m.Epics[m.EpicCursor]
			m.FocusedTask = 0
			m.DetailOffset = 0
			m.refreshTasks()
			m.ensureFocusedTaskVisible()
			m.Overlay = OverlayNone
			if m.Root != "" {
				return m, writeRouterReleaseEpicCmd(m.Root, m.SelectedEpic, m.SelectedRelease, m.Dependencies.RouterReader)
			}
		}
	}
	return m, nil
}

func (m Model) handleReleaseOverlay(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.Overlay = OverlayNone
	case "up", "k":
		if m.ReleaseCursor > 0 {
			m.ReleaseCursor--
		}
	case "down", "j":
		if len(m.Releases) > 0 && m.ReleaseCursor < len(m.Releases)-1 {
			m.ReleaseCursor++
		}
	case "enter":
		if len(m.Releases) > 0 {
			m.SelectedRelease = m.Releases[m.ReleaseCursor]
			m.refreshEpicsForRelease()
			m.FocusedTask = 0
			m.DetailOffset = 0
			m.refreshTasks()
			m.ensureFocusedTaskVisible()
			m.Overlay = OverlayNone
			if m.Root != "" {
				return m, writeRouterReleaseEpicCmd(m.Root, m.SelectedEpic, m.SelectedRelease, m.Dependencies.RouterReader)
			}
		}
	}
	return m, nil
}

func (m Model) handleDetailOverlay(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.Overlay = OverlayNone
	case "up", "k":
		if m.DetailOffset > 0 {
			m.DetailOffset--
		}
	case "down", "j":
		m.DetailOffset++
	case "pgup":
		m.DetailOffset -= m.detailPageSize()
		if m.DetailOffset < 0 {
			m.DetailOffset = 0
		}
	case "pgdown":
		m.DetailOffset += m.detailPageSize()
	}
	return m, nil
}

func (m Model) handleEpicDetailOverlay(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.Overlay = OverlayNone
	case "1":
		m.EpicDetailTab = 0
		m.EpicDetailOffset = 0
	case "2":
		m.EpicDetailTab = 1
		m.EpicDetailOffset = 0
		if m.EpicAuditContent == "" {
			epicSlug := m.epicDetailEpic()
			epicDir := filepath.Join(m.Root, "releases", m.SelectedRelease, "epics", epicSlug)
			return m, readEpicAuditCmd(epicDir, shortID(epicSlug))
		}
	case "3":
		m.EpicDetailTab = 2
		if len(m.ReleaseDocs) == 0 && m.Root != "" {
			return m, loadReleaseDocsCmd(m.Root)
		}
	case "[", "left", "h":
		if m.EpicDetailTab == 2 {
			m.selectReleaseDoc(m.ReleaseDocIndex - 1)
		}
	case "]", "right", "l":
		if m.EpicDetailTab == 2 {
			m.selectReleaseDoc(m.ReleaseDocIndex + 1)
		}
	case "a":
		return m.markEpicAudited()
	case "up", "k":
		if m.EpicDetailTab == 2 {
			m.scrollReleaseDoc(-1)
		} else if m.EpicDetailOffset > 0 {
			m.EpicDetailOffset--
		}
	case "down", "j":
		if m.EpicDetailTab == 2 {
			m.scrollReleaseDoc(1)
		} else {
			m.EpicDetailOffset++
		}
	case "pgup":
		if m.EpicDetailTab == 2 {
			m.scrollReleaseDoc(-m.detailPageSize())
		} else {
			m.EpicDetailOffset -= m.detailPageSize()
			if m.EpicDetailOffset < 0 {
				m.EpicDetailOffset = 0
			}
		}
	case "pgdown":
		if m.EpicDetailTab == 2 {
			m.scrollReleaseDoc(m.detailPageSize())
		} else {
			m.EpicDetailOffset += m.detailPageSize()
		}
	}
	return m, nil
}

// selectReleaseDoc moves the Release Docs selection to idx, clamped to the
// loaded docs. Each doc keeps its own scroll offset, so switching does not
// disturb the other document's position.
func (m *Model) selectReleaseDoc(idx int) {
	if idx < 0 || idx >= len(m.ReleaseDocs) {
		return
	}
	m.ReleaseDocIndex = idx
}

// scrollReleaseDoc adjusts the selected document's scroll offset by delta,
// clamped at the top. The upper bound is left to the renderer, matching the
// Detail/Audit scroll behavior. Offsets are stored per doc ID so the
// unselected document retains its position.
func (m *Model) scrollReleaseDoc(delta int) {
	doc, ok := m.selectedReleaseDoc()
	if !ok {
		return
	}
	if m.ReleaseDocOffsets == nil {
		m.ReleaseDocOffsets = map[data.ReleaseDocID]int{}
	}
	offset := m.ReleaseDocOffsets[doc.ID] + delta
	if offset < 0 {
		offset = 0
	}
	m.ReleaseDocOffsets[doc.ID] = offset
}

// selectedReleaseDoc returns the currently selected release doc, if any.
func (m *Model) selectedReleaseDoc() (data.ReleaseDoc, bool) {
	if m.ReleaseDocIndex < 0 || m.ReleaseDocIndex >= len(m.ReleaseDocs) {
		return data.ReleaseDoc{}, false
	}
	return m.ReleaseDocs[m.ReleaseDocIndex], true
}

// clampReleaseDocIndex keeps the selection in range after docs (re)load.
func (m *Model) clampReleaseDocIndex() {
	if m.ReleaseDocIndex < 0 {
		m.ReleaseDocIndex = 0
	}
	if m.ReleaseDocIndex >= len(m.ReleaseDocs) {
		m.ReleaseDocIndex = max(len(m.ReleaseDocs)-1, 0)
	}
}

// markEpicAudited writes status audited for the open epic, guarded by the detail
// file's mtime captured when the overlay was opened. It no-ops with a message if
// the epic is already audited.
func (m Model) markEpicAudited() (tea.Model, tea.Cmd) {
	epicSlug := m.epicDetailEpic()
	if epicSlug == "" {
		return m, nil
	}
	if m.EpicStatus[epicSlug] == string(data.EpicStatusAudited) {
		m.StatusMessage = shortID(epicSlug) + " already audited"
		return m, nil
	}
	if m.Root == "" {
		m.StatusMessage = "Epic not updated: no savepoint root"
		return m, nil
	}
	path := epicDetailFilePath(m.Root, m.SelectedRelease, epicSlug)
	return m, writeEpicStatusCmd(epicSlug, path, string(data.EpicStatusAudited), m.EpicDetailMtime)
}

// epicDetailFilePath resolves the E##-Detail.md path for an epic slug, matching
// how loadBoardData and the audit-tab branch build it.
func epicDetailFilePath(root, release, epicSlug string) string {
	epicDir := filepath.Join(root, "releases", release, "epics", epicSlug)
	return filepath.Join(epicDir, shortID(epicSlug)+"-Detail.md")
}

func taskWriteErrorMessage(err error) string {
	if err == data.ErrMtimeConflict {
		return "mtime conflict: refresh before retrying"
	}
	if strings.Contains(err.Error(), "complexity_reason") {
		return "Task not moved: " + err.Error() + ". Shorten complexity_reason or adjust complexity_tier before retrying"
	}
	if strings.Contains(err.Error(), "invalid status") {
		return "Task not moved: " + err.Error()
	}
	return err.Error()
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
			return m, m.selectEpicPanelEpic()
		}
	case "down", "j":
		if m.EpicPanelCursor < len(m.Epics)-1 {
			m.EpicPanelCursor++
			return m, m.selectEpicPanelEpic()
		}
	case "enter":
		return m, m.openEpicDetailOverlay()
	case "right", "l":
		m.EpicPanelFocus = false
		m.FocusedColumn = data.ColumnPlanned
		m.FocusedTask = 0
		m.ensureFocusedTaskVisible()
	case "left", "h":
	}
	return m, nil
}

func (m *Model) selectEpicPanelEpic() tea.Cmd {
	if len(m.Epics) == 0 || m.EpicPanelCursor < 0 || m.EpicPanelCursor >= len(m.Epics) {
		return nil
	}
	m.SelectedEpic = m.Epics[m.EpicPanelCursor]
	m.FocusedTask = 0
	m.DetailOffset = 0
	m.refreshTasks()
	m.ensureFocusedTaskVisible()
	if m.Root != "" {
		return writeRouterReleaseEpicCmd(m.Root, m.SelectedEpic, m.SelectedRelease, m.Dependencies.RouterReader)
	}
	return nil
}

func (m *Model) openEpicDetailOverlay() tea.Cmd {
	if len(m.Epics) == 0 || m.EpicPanelCursor < 0 || m.EpicPanelCursor >= len(m.Epics) {
		return nil
	}
	epicSlug := m.Epics[m.EpicPanelCursor]
	shortEpicID := shortID(epicSlug)
	epicDir := filepath.Join(m.Root, "releases", m.SelectedRelease, "epics", epicSlug)
	m.EpicDetailEpic = epicSlug
	m.EpicDetailOffset = 0
	m.EpicDetailTab = 0
	m.EpicAuditContent = ""
	m.ReleaseDocs = nil
	m.ReleaseDocIndex = 0
	m.ReleaseDocOffsets = map[data.ReleaseDocID]int{}
	m.EpicDetailMtime = time.Time{}
	if fi, err := os.Stat(filepath.Join(epicDir, shortEpicID+"-Detail.md")); err == nil {
		m.EpicDetailMtime = fi.ModTime()
	}
	m.Overlay = OverlayEpicDetail
	cmd := readEpicDetailCmd(epicDir, shortEpicID)
	if m.Root != "" {
		return tea.Batch(cmd, loadReleaseDocsCmd(m.Root))
	}
	return cmd
}

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

func (m Model) epicDetailEpic() string {
	if m.EpicDetailEpic != "" {
		return m.EpicDetailEpic
	}
	if m.SelectedEpic != "" {
		return m.SelectedEpic
	}
	if len(m.Epics) > 0 && m.EpicPanelCursor >= 0 && m.EpicPanelCursor < len(m.Epics) {
		return m.Epics[m.EpicPanelCursor]
	}
	return ""
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

func (m *Model) clampDefectCursor() {
	defects := defectsForOverlay(m.AllDefects, m.SelectedRelease)
	if len(defects) == 0 {
		m.DefectCursor = 0
		return
	}
	if m.DefectCursor >= len(defects) {
		m.DefectCursor = len(defects) - 1
	}
	if m.DefectCursor < 0 {
		m.DefectCursor = 0
	}
}

func sameDefectRecord(a, b data.Defect) bool {
	if a.Path != "" && b.Path != "" {
		return a.Path == b.Path
	}
	if a.ID != b.ID {
		return false
	}
	if a.Release != "" && b.Release != "" && a.Release != b.Release {
		return false
	}
	return true
}

func (m Model) columnPageSize() int {
	h := m.Height
	if h == 0 {
		h = defaultTermH
	}
	return visibleColumnTaskLimit(CalculateLayout(m.Width, h).ContentHeight)
}

// conservativeColumnPageSize reserves 2 lines for scroll indicators so that
// ensureFocusedTaskVisible never sets an offset where the focused card is hidden
// by a top or bottom indicator that wasn't accounted for in the page budget.
func (m Model) conservativeColumnPageSize() int {
	h := m.Height
	if h == 0 {
		h = defaultTermH
	}
	contentHeight := CalculateLayout(m.Width, h).ContentHeight
	// contentHeight - 2 = card budget; subtract 2 more for both indicators.
	limit := (contentHeight - 4) / 3
	if limit < 1 {
		return 1
	}
	return limit
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

func (m Model) handleDefectOverlay(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	defects := defectsForOverlay(m.AllDefects, m.SelectedRelease)
	switch msg.String() {
	case "esc", "q":
		m.Overlay = OverlayNone
	case "up", "k":
		if m.DefectCursor > 0 {
			m.DefectCursor--
		}
	case "down", "j":
		if len(defects) > 0 && m.DefectCursor < len(defects)-1 {
			m.DefectCursor++
		}
	case "enter":
		if len(defects) > 0 && m.DefectCursor < len(defects) {
			m.DefectDetailOffset = 0
			m.Overlay = OverlayDefectDetail
		}
	case " ":
		if len(defects) == 0 || m.DefectCursor >= len(defects) {
			return m, nil
		}
		defect := defects[m.DefectCursor]
		switch defect.Status {
		case "", data.DefectOpen:
			if defect.Path == "" {
				m.StatusMessage = "Defect not updated: missing file path"
				return m, nil
			}
			next := defect
			next.Status = data.DefectResolved
			next.Stage = ""
			return m, writeDefectStatusCmd(next, defect.Mtime)
		case data.DefectResolved:
			m.StatusMessage = "Defect already resolved"
		case data.DefectInProgress:
			m.StatusMessage = "Defect in progress: resolve after lifecycle stage is closed"
		default:
			m.StatusMessage = "Defect not updated: invalid status " + string(defect.Status)
		}
	}
	return m, nil
}

func (m Model) handleDefectDetailOverlay(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.Overlay = OverlayDefect
	case "up", "k":
		if m.DefectDetailOffset > 0 {
			m.DefectDetailOffset--
		}
	case "down", "j":
		m.DefectDetailOffset++
	case "pgup":
		m.DefectDetailOffset -= m.detailPageSize()
		if m.DefectDetailOffset < 0 {
			m.DefectDetailOffset = 0
		}
	case "pgdown":
		m.DefectDetailOffset += m.detailPageSize()
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
