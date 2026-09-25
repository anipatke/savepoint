package v2

import (
	"fmt"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/opencode/savepoint/internal/data"
)

// columnOrder is the left-to-right order of the three columns, and the order
// horizontal movement walks. It is the recorded status vocabulary and nothing
// else: there is no fourth column, and no lifecycle flag is promoted to one.
var columnOrder = []data.ColumnType{data.ColumnPlanned, data.ColumnInProgress, data.ColumnDone}

// Update is a reducer over typed messages only: it reads no file, starts no
// subprocess, and performs no IO (ARCH-02). Filesystem actions are returned as
// explicit Bubble Tea commands. Navigation clamps at each surface's edges.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		if !m.sidebarVisible() {
			// A terminal that shrank past the breakpoint stops drawing the
			// sidebar, so focus returns to the columns rather than staying on a
			// surface that is no longer on screen.
			m.SidebarFocused = false
		}
		// A shorter terminal gives an open overlay a shorter window, which can
		// leave it scrolled past its own end.
		m.clampDetailScroll()
		m.clampIssueScroll()
		return m, nil
	case projectLoadedMsg:
		return m.applyLoad(msg)
	case v2FileChangeMsg:
		// The watcher command has completed after its debounce window. Start it
		// again alongside the one load command so edits that arrive while the
		// index is rebuilding are not lost.
		if m.Watcher == nil {
			return m, loadCmd(m.Root)
		}
		return m, tea.Batch(loadCmd(m.Root), watchV2Files(m.Watcher, m.Root))
	case actionMsg:
		if msg.err != nil {
			if msg.releaseRollback {
				m.restoreReleaseWrite()
			}
			m.StatusMessage = msg.err.Error()
			m.preserveReloadStatus = msg.reload
			if msg.reload {
				return m, loadCmd(m.Root)
			}
			return m, nil
		}
		m.StatusMessage = msg.message
		m.releaseRollback = nil
		m.preserveReloadStatus = false
		if msg.issueSelectedID != "" && m.Issues != nil {
			m.Issues.SelectedID = msg.issueSelectedID
		}
		if msg.reload {
			return m, loadCmd(m.Root)
		}
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

// handleKey dispatches one key to the surface that has focus: an open detail
// overlay first, then the sidebar, then the columns. Surface focus is
// arrow-driven only — left from the Planned column enters the sidebar, and
// right from the sidebar returns to it — so arrow keys alone can walk every
// column end to end without a separate Tab key. While an overlay is open,
// keys do not cross at all, because they belong to what is on top.
func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()
	// A status message shares the footer with the key hints. The next key
	// dismisses it so the hints return; an action's own result replaces it.
	m.StatusMessage = ""
	if m.Help {
		if key == "esc" || key == "q" || key == "?" {
			m.Help = false
		}
		return m, nil
	}
	if m.ReleaseOverlay {
		return m.handleReleaseKey(key)
	}
	if key == "?" {
		m.Help = true
		return m, nil
	}
	if key == "q" || key == "ctrl+c" {
		if m.Watcher != nil {
			_ = m.Watcher.Close()
		}
		return m, tea.Quit
	}

	if m.Detail != nil {
		if action, ok := m.actionForKey(key); ok {
			return m, m.runAction(action)
		}
		m.handleDetailKey(key)
		return m, nil
	}
	if m.Issues != nil {
		return m, m.handleIssuesKey(key)
	}
	if action, ok := m.actionForKey(key); ok {
		return m, m.runAction(action)
	}
	if key == goalSelectorKey || key == goalSelectorAlias {
		m.openReleaseSelector()
		return m, nil
	}

	if !m.SidebarFocused {
		if key == " " || key == "backspace" {
			if taskID := m.focusedTaskID(); taskID != "" {
				if key == " " {
					return m, writeTaskAdvanceCmd(m.Root, taskID)
				}
				return m, writeTaskRetreatCmd(m.Root, taskID)
			}
			return m, nil
		}
	}
	if m.SidebarFocused {
		// A refusal explains why an open Objective cannot close yet; a done
		// Objective has nothing to close, so space leaves it alone.
		if key == objectiveCloseKey {
			if target, ok := m.focusedActionTarget(); ok && target.Kind == DetailObjective && !m.objectiveDone(target.ID) {
				return m, writeObjectiveCompletionCmd(m.Root, target.ID)
			}
		}
		return m, m.handleSidebarKey(key)
	}
	m.handleColumnKey(key)
	return m, nil
}

// handleReleaseKey owns the selector while it is open. In particular q is a
// cancel key here, not the board's quit key, so a user can inspect a context
// and leave it unchanged with either Esc or q.
func (m Model) handleReleaseKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "esc", "q":
		m.closeReleaseSelector()
		return m, nil
	case "up", "k":
		if m.ReleaseCursor > 0 {
			m.ReleaseCursor--
		}
		return m, nil
	case "down", "j":
		if len(m.Releases) > 0 && m.ReleaseCursor < len(m.Releases)-1 {
			m.ReleaseCursor++
		}
		return m, nil
	case "enter":
		if len(m.Releases) == 0 {
			return m, nil
		}
		release := m.Releases[m.ReleaseCursor]
		before := m.snapshotReload()
		before.SelectedRelease = m.SelectedRelease
		m.applyReleaseSelection(release)
		m.ReleaseOverlay = false
		m.releaseOriginSet = false
		if m.Root == "" {
			return m, nil
		}
		m.releaseRollback = &before
		return m, writeReleaseSelectionCmd(m.Root, release, m.State.RouterMtime)
	case "v", "d":
		if len(m.Releases) == 0 || m.State.Index == nil {
			return m, nil
		}
		release, ok := newReleaseDetail(m.State.Index, m.Releases[m.ReleaseCursor])
		if !ok {
			return m, nil
		}
		// The detail belongs to the Release the selector was opened over, but
		// opening it is still a read: it does not select or persist that
		// Release. Closing therefore restores the pre-selector surface.
		m.DetailOrigin = m.releaseOrigin
		m.releaseOriginSet = false
		m.ReleaseOverlay = false
		m.Detail = &release
		m.DetailOffset = 0
		return m, nil
	}
	return m, nil
}

func (m *Model) closeReleaseSelector() {
	m.ReleaseOverlay = false
	if m.releaseOriginSet {
		m.SidebarFocused = m.releaseOrigin.SidebarFocused && m.sidebarVisible()
		m.ObjectiveCursor = m.releaseOrigin.ObjectiveCursor
		m.FocusedColumn = m.releaseOrigin.FocusedColumn
		m.FocusedCard = m.releaseOrigin.FocusedCard
		m.releaseOriginSet = false
		m.clampObjectiveCursorToRows()
		m.clampFocus()
	}
}

func (m *Model) openReleaseSelector() {
	m.releaseOrigin = detailOrigin{
		SidebarFocused:  m.SidebarFocused,
		ObjectiveCursor: m.ObjectiveCursor,
		FocusedColumn:   m.FocusedColumn,
		FocusedCard:     m.FocusedCard,
	}
	m.releaseOriginSet = true
	m.ReleaseOverlay = true
	m.ReleaseCursor = releaseIndex(m.Releases, m.SelectedRelease)
}

func (m Model) runAction(action BoardAction) tea.Cmd {
	target := actionTarget{Kind: action.TargetKind, ID: action.TargetID}
	switch action.Kind {
	case ActionAcceptCheck:
		return writeOwnerAcceptanceCmd(m.Root, target)
	case ActionCompleteByException:
		return writeExceptionCompletionCmd(m.Root, target)
	case ActionCompleteObjective:
		return writeObjectiveCompletionCmd(m.Root, target.ID)
	case ActionRecordSelection:
		if m.State.Router == nil {
			return func() tea.Msg { return actionMsg{err: fmt.Errorf("selection requires a loaded router state")} }
		}
		current := data.RouterSelectionV2{
			Release: m.State.Router.Release,
			Issue:   m.State.Router.Issue,
		}
		selection, err := selectionForTarget(m.State.Index, target, current)
		if err != nil {
			return func() tea.Msg { return actionMsg{err: err} }
		}
		return writeSelectionCmd(m.Root, selection, m.State.RouterMtime)
	default:
		return func() tea.Msg { return actionMsg{err: fmt.Errorf("unknown board action %q", action.Kind)} }
	}
}

// handleColumnKey moves the card cursor and opens the focused Task's detail.
// Enter has no other meaning on the columns, so it is the detail key there;
// v is the same action, and is what the sidebar uses because enter already
// selects an Objective. Left from the leftmost column crosses into the
// sidebar rather than clamping, when the sidebar is on screen to receive it.
func (m *Model) handleColumnKey(key string) {
	switch key {
	case "left", "h":
		if columnIndex(m.FocusedColumn) == 0 && m.sidebarVisible() {
			m.focusSidebar()
			return
		}
		m.focusColumn(-1)
	case "right", "l":
		m.focusColumn(1)
	case "up", "k":
		m.focusCard(-1)
	case "down", "j":
		m.focusCard(1)
	case "enter", "v":
		m.openDetail()
	case "i":
		m.openIssues("")
	case "I":
		m.openIssues(m.focusedTaskID())
	}
}

// handleDetailKey drives the open overlay: it scrolls, clamping at both ends so
// a repeated press changes nothing, and closes back to the surface it came
// from. Nothing here writes, and nothing here reaches the filesystem — an
// overlay is a view of records the last load already resolved.
func (m *Model) handleDetailKey(key string) {
	switch key {
	case "up", "k":
		m.scrollDetail(-1)
	case "down", "j":
		m.scrollDetail(1)
	case "esc":
		m.closeDetail()
	}
}

// openDetail opens the record under the cursor on whichever surface has focus:
// the Objective the sidebar cursor is on, or the focused Task card. It records
// where it was opened from so closing returns the keys there, and it changes no
// selection — opening a record is not choosing one.
func (m *Model) openDetail() {
	detail, ok := m.detailUnderCursor()
	if !ok {
		return
	}
	m.DetailOrigin = detailOrigin{
		SidebarFocused:  m.SidebarFocused,
		ObjectiveCursor: m.ObjectiveCursor,
		FocusedColumn:   m.FocusedColumn,
		FocusedCard:     m.FocusedCard,
	}
	m.Detail = &detail
	m.DetailOffset = 0
}

// hasDetailTarget reports whether the surface holding focus has a record
// under its cursor at all — an empty column, an Objective-less project, or a
// project held back by a pending migration, none of which has one. The
// footer's detail hint asks this before ever offering the key.
func (m Model) hasDetailTarget() bool {
	if m.State.Index == nil {
		return false
	}
	if m.SidebarFocused {
		return m.ObjectiveCursor >= 0 && m.ObjectiveCursor < len(m.Objectives)
	}
	cards := m.Cards[m.FocusedColumn]
	return m.FocusedCard >= 0 && m.FocusedCard < len(cards)
}

// detailUnderCursor resolves the record the cursor is on, or reports that there
// is none, per hasDetailTarget.
func (m Model) detailUnderCursor() (RecordDetail, bool) {
	if !m.hasDetailTarget() {
		return RecordDetail{}, false
	}
	if m.SidebarFocused {
		return newObjectiveDetail(m.State.Index, m.Objectives[m.ObjectiveCursor].ID())
	}
	cards := m.Cards[m.FocusedColumn]
	return newTaskDetail(m.State.Index, cards[m.FocusedCard].Task.ID)
}

// closeDetail returns the keys to the surface the overlay was opened from, with
// the focus and cursor it had — clamped into the records that exist now, since
// a reload underneath the overlay may have removed some.
func (m *Model) closeDetail() {
	m.Detail = nil
	m.DetailOffset = 0
	m.SidebarFocused = m.DetailOrigin.SidebarFocused && m.sidebarVisible()
	m.ObjectiveCursor = m.DetailOrigin.ObjectiveCursor
	m.FocusedColumn = m.DetailOrigin.FocusedColumn
	m.FocusedCard = m.DetailOrigin.FocusedCard
	m.clampObjectiveCursorToRows()
	m.clampFocus()
}

// scrollDetail moves the overlay's window against the same bound the renderer
// clamps to, so both ends hold and a repeated press is a no-op.
func (m *Model) scrollDetail(delta int) {
	width, height := m.detailViewport()
	next := m.DetailOffset + delta
	if next < 0 || next > detailScrollLimit(*m.Detail, width, height) {
		return
	}
	m.DetailOffset = next
}

// clampDetailScroll pulls an open overlay back inside its content after
// anything that could have shortened the window or the record.
func (m *Model) clampDetailScroll() {
	if m.Detail == nil {
		return
	}
	width, height := m.detailViewport()
	if limit := detailScrollLimit(*m.Detail, width, height); m.DetailOffset > limit {
		m.DetailOffset = limit
	}
}

// refreshDetail re-resolves an open overlay against the load that just
// completed, so a record edited underneath it shows its new state. A record the
// load no longer holds closes the overlay: leaving the last copy on screen
// would show something older than the project.
func (m *Model) refreshDetail() {
	if m.Detail == nil {
		return
	}
	detail, ok := reopenDetail(m.State.Index, *m.Detail)
	if !ok {
		m.closeDetail()
		return
	}
	m.Detail = &detail
	m.clampDetailScroll()
}

// handleSidebarKey moves the Objective cursor, applies the new selection
// immediately, and starts explicit order writes as commands. Right crosses
// back into the columns, the mirror of the left key that crossed in from the
// Planned column.
func (m *Model) handleSidebarKey(key string) tea.Cmd {
	if change, ok := sidebarObjectiveOrderChange(m.Objectives, m.ObjectiveCursor, key); ok {
		if m.State.Index == nil || m.SelectedRelease == "" {
			return nil
		}
		return writeObjectiveGroupOrderCmd(m.State.Index, m.SelectedRelease, change)
	}

	switch key {
	case "right", "l":
		m.SidebarFocused = false
		m.FocusedColumn = data.ColumnPlanned
		m.clampFocus()
	case "up", "k":
		m.moveObjectiveCursor(-1)
	case "down", "j":
		m.moveObjectiveCursor(1)
	case "v":
		m.openDetail()
	case "i":
		m.openIssues("")
	case "esc":
		m.selectObjective("")
	}
	return nil
}

type objectiveGroupOrderChange struct {
	ObjectiveID  string
	Priority     data.ObjectivePriority
	ObjectiveIDs []string
	Message      string
}

// sidebarObjectiveOrderChange resolves a focused row's keyboard request using
// only the already ordered sidebar rows. It does not write or mutate them.
func sidebarObjectiveOrderChange(rows []ObjectiveRow, cursor int, key string) (objectiveGroupOrderChange, bool) {
	if cursor < 0 || cursor >= len(rows) || rows[cursor].Objective == nil {
		return objectiveGroupOrderChange{}, false
	}
	focused := rows[cursor]
	if focused.Objective.Status == data.ColumnDone {
		return objectiveGroupOrderChange{}, false
	}
	objectiveID := focused.ID()
	currentPriority := sidebarRowPriority(focused)

	if destination, ok := sidebarPriorityKey(key); ok {
		if destination == currentPriority {
			return objectiveGroupOrderChange{}, false
		}
		ids := sidebarObjectiveIDsInPriority(rows, destination, objectiveID)
		ids = append(ids, objectiveID)
		return objectiveGroupOrderChange{
			ObjectiveID:  objectiveID,
			Priority:     destination,
			ObjectiveIDs: ids,
			Message:      fmt.Sprintf("Objective %s moved to %s priority.", objectiveID, objectivePriorityLabel(destination)),
		}, true
	}

	var delta int
	var direction string
	switch key {
	case "K", "shift+up":
		delta, direction = -1, "up"
	case "J", "shift+down":
		delta, direction = 1, "down"
	default:
		return objectiveGroupOrderChange{}, false
	}

	ids := sidebarObjectiveIDsInPriority(rows, currentPriority, "")
	position := slices.Index(ids, objectiveID)
	neighbor := position + delta
	if position < 0 || neighbor < 0 || neighbor >= len(ids) {
		return objectiveGroupOrderChange{}, false
	}
	ids[position], ids[neighbor] = ids[neighbor], ids[position]
	return objectiveGroupOrderChange{
		ObjectiveID:  objectiveID,
		Priority:     currentPriority,
		ObjectiveIDs: ids,
		Message:      fmt.Sprintf("Objective %s moved %s within %s priority.", objectiveID, direction, objectivePriorityLabel(currentPriority)),
	}, true
}

func sidebarPriorityKey(key string) (data.ObjectivePriority, bool) {
	switch key {
	case "1":
		return data.ObjectivePriority("critical"), true
	case "2":
		return data.ObjectivePriority("high"), true
	case "3":
		return data.ObjectivePriority("medium"), true
	case "4":
		return data.ObjectivePriority("low"), true
	default:
		return "", false
	}
}

func sidebarRowPriority(row ObjectiveRow) data.ObjectivePriority {
	if row.Objective == nil || row.Objective.Priority == "" {
		return data.ObjectivePriority("medium")
	}
	return row.Objective.Priority
}

func sidebarObjectiveIDsInPriority(rows []ObjectiveRow, priority data.ObjectivePriority, exceptID string) []string {
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		if row.Objective != nil && row.Objective.Status != data.ColumnDone && row.ID() != exceptID && sidebarRowPriority(row) == priority {
			ids = append(ids, row.ID())
		}
	}
	return ids
}

func objectivePriorityLabel(priority data.ObjectivePriority) string {
	switch priority {
	case "critical":
		return "Critical"
	case "high":
		return "High"
	case "low":
		return "Low"
	default:
		return "Medium"
	}
}

// focusSidebar moves focus onto the sidebar, for the left-arrow crossing at
// the Planned column's edge, refusing to move it into a sidebar the current
// terminal is too narrow to draw.
func (m *Model) focusSidebar() {
	if !m.sidebarVisible() {
		return
	}
	m.SidebarFocused = true
}

// moveObjectiveCursor walks the sidebar, clamping at both ends rather than
// wrapping so a repeated press at either end changes nothing, and applies
// the row it lands on as the selection immediately — the V1 interaction,
// with no separate enter:select step.
func (m *Model) moveObjectiveCursor(delta int) {
	next := m.ObjectiveCursor + delta
	if next < 0 || next >= len(m.Objectives) {
		return
	}
	m.ObjectiveCursor = next
	m.selectObjective(m.Objectives[next].ID())
}

// selectObjective filters the columns to the Tasks objectiveID owns, or to the
// whole project when objectiveID is empty. Re-selecting what is already
// selected changes nothing, and the card cursor returns to the top of the
// filtered columns rather than to a position the new set may not have.
func (m *Model) selectObjective(objectiveID string) {
	if m.SelectedObjective == objectiveID {
		return
	}
	m.SelectedObjective = objectiveID
	m.Cards = groupTaskCardsForRelease(m.State.Index, m.SelectedRelease, objectiveID)
	m.FocusedCard = 0
	m.clampFocus()
}

// applyReleaseSelection changes only the board's delivery context. Objective
// and Task ownership stay in the index: the visible rows/cards are rebuilt
// from its reverse links, and cursors retain their identities where those
// identities still exist in the new context.
func (m *Model) applyReleaseSelection(releaseID string) {
	if m.State.Index == nil {
		m.SelectedRelease = releaseID
		return
	}

	view := m.snapshotReload()
	view.SelectedRelease = m.SelectedRelease
	m.SelectedRelease = releaseID
	if m.SelectedObjective != "" && !objectiveBelongsToRelease(m.State.Index, m.SelectedObjective, releaseID) {
		m.SelectedObjective = ""
	}
	m.Objectives = objectiveRowsForRelease(m.State.Index, releaseID)
	m.Cards = groupTaskCardsForRelease(m.State.Index, releaseID, m.SelectedObjective)
	m.restoreObjectiveCursor(view, true)
	m.restoreFocus(view, true)
}

func objectiveBelongsToRelease(index *data.V2Index, objectiveID, releaseID string) bool {
	if index == nil || objectiveID == "" || releaseID == "" {
		return false
	}
	objective, ok := index.Objectives[objectiveID]
	if !ok {
		return false
	}
	return string(objective.Release) == releaseID
}

type reloadSnapshot struct {
	SelectedRelease   string
	ReleaseCursorID   string
	SelectedObjective string
	ObjectiveCursorID string
	FocusedColumn     data.ColumnType
	FocusedCard       int
	FocusedTaskID     string
	RouterRelease     string
	RouterObjective   string
	DetailKind        DetailKind
	DetailID          string
	IssueSelectedID   string
	IssueDetailID     string
	IssueScopedTask   string
}

func (m Model) snapshotReload() reloadSnapshot {
	snapshot := reloadSnapshot{
		SelectedObjective: m.SelectedObjective,
		SelectedRelease:   m.SelectedRelease,
		FocusedColumn:     m.FocusedColumn,
		FocusedCard:       m.FocusedCard,
		FocusedTaskID:     m.focusedTaskID(),
	}
	if m.State.Router != nil {
		snapshot.RouterRelease = m.State.Router.Release
	}
	snapshot.RouterObjective = routerObjective(m.State)
	if m.ReleaseCursor >= 0 && m.ReleaseCursor < len(m.Releases) {
		snapshot.ReleaseCursorID = m.Releases[m.ReleaseCursor]
	}
	if m.ObjectiveCursor >= 0 && m.ObjectiveCursor < len(m.Objectives) {
		snapshot.ObjectiveCursorID = m.Objectives[m.ObjectiveCursor].ID()
	}
	if m.Detail != nil {
		snapshot.DetailKind = m.Detail.Kind
		snapshot.DetailID = m.Detail.ID
	}
	if m.Issues != nil {
		snapshot.IssueSelectedID = m.Issues.SelectedID
		snapshot.IssueScopedTask = m.Issues.ScopedTask
		if m.Issues.Detail != nil {
			snapshot.IssueDetailID = m.Issues.Detail.Issue.ID
		}
	}
	return snapshot
}

// applyLoad folds one load result into the model. A failed load is re-read
// once before it is reported. A failed reload keeps the last completed state
// visible and names the temporary data problem alongside it; only an initial
// failure has no board to preserve.
func (m Model) applyLoad(msg projectLoadedMsg) (tea.Model, tea.Cmd) {
	if msg.Seq != 0 {
		if msg.Seq < m.lastLoadSeq {
			return m, nil
		}
		m.lastLoadSeq = msg.Seq
	}
	if msg.Failed() && !msg.Retry {
		return m, retryLoadCmd(m.Root)
	}
	wasLoaded := m.Loaded
	snapshot := m.snapshotReload()
	statusBeforeReload := m.StatusMessage
	preserveStatus := m.preserveReloadStatus
	m.Loaded = true
	if msg.Failed() {
		if wasLoaded {
			m.ReloadDiagnostic = msg.Diagnostic
			return m, nil
		}
		m.Diagnostic = msg.Diagnostic
		m.State = ProjectState{}
		m.SelectedRelease = ""
		m.Releases = nil
		m.SelectedObjective = ""
		m.Objectives = nil
		m.ObjectiveCursor = 0
		m.Cards = groupTaskCards(nil)
		m.FocusedCard = 0
		m.Detail = nil
		m.DetailOffset = 0
		m.Issues = nil
		return m, nil
	}

	m.Diagnostic = ""
	m.ReloadDiagnostic = ""
	m.StatusMessage = ""
	if preserveStatus {
		m.StatusMessage = statusBeforeReload
	}
	m.preserveReloadStatus = false
	m.releaseRollback = nil
	m.State = msg.State
	if err := objectiveFilterError(msg.State.Index, m.ObjectiveFilter); err != nil {
		m.FatalErr = err
		return m, tea.Quit
	}
	m.Releases = orderedReleaseIDs(msg.State.Index)
	m.SelectedRelease = restoredRelease(msg.State)
	m.SelectedObjective = restoredObjectiveForRelease(m, snapshot, msg.State, wasLoaded, m.SelectedRelease)
	m.Objectives = objectiveRowsForRelease(msg.State.Index, m.SelectedRelease)
	m.Cards = groupTaskCardsForRelease(msg.State.Index, m.SelectedRelease, m.SelectedObjective)
	m.restoreReleaseCursor(snapshot, wasLoaded)
	m.restoreObjectiveCursor(snapshot, wasLoaded)
	m.restoreFocus(snapshot, wasLoaded)
	m.restoreOverlayOrigins()
	if wasLoaded && snapshot.SelectedObjective != "" && snapshot.RouterObjective == routerObjective(msg.State) &&
		!objectiveExists(msg.State.Index, snapshot.SelectedObjective) {
		m.noteStatus(fmt.Sprintf("Objective %s no longer exists; selection moved to %s.", snapshot.SelectedObjective, objectiveLabel(m.SelectedObjective)))
	}
	if wasLoaded && snapshot.SelectedRelease != "" && snapshot.RouterRelease == routerRelease(msg.State) &&
		!releaseExists(msg.State.Index, snapshot.SelectedRelease) {
		m.noteStatus(fmt.Sprintf("%s %s no longer exists; selection cleared.", goalLabel, snapshot.SelectedRelease))
	}
	m.refreshDetail()
	if snapshot.DetailID != "" && m.Detail == nil && recordExists(msg.State.Index, snapshot.DetailKind, snapshot.DetailID) == false {
		m.noteStatus(fmt.Sprintf("%s %s no longer exists; detail closed.", detailKindLabel(snapshot.DetailKind), snapshot.DetailID))
	}
	m.refreshIssues()
	if snapshot.IssueScopedTask != "" && !taskExists(msg.State.Index, snapshot.IssueScopedTask) {
		m.closeIssues()
		m.noteStatus(fmt.Sprintf("Task %s no longer exists; Issues overlay closed.", snapshot.IssueScopedTask))
	} else if snapshot.IssueDetailID != "" && !issueExists(msg.State.Index, snapshot.IssueDetailID) {
		m.noteStatus(fmt.Sprintf("Issue %s no longer exists; focus moved to a neighboring issue.", snapshot.IssueDetailID))
	} else if snapshot.IssueSelectedID != "" && !issueExists(msg.State.Index, snapshot.IssueSelectedID) {
		m.noteStatus(fmt.Sprintf("Issue %s no longer exists; focus moved to a neighboring issue.", snapshot.IssueSelectedID))
	}
	return m, nil
}

func routerRelease(state ProjectState) string {
	if state.Router == nil {
		return ""
	}
	return state.Router.Release
}

// restoredRelease treats the router as the persisted source of truth. A
// failed optimistic write therefore cannot leave its candidate context behind
// after the board reloads.
func restoredRelease(state ProjectState) string {
	return selectedRelease(state)
}

func objectiveExists(index *data.V2Index, id string) bool {
	if index == nil {
		return false
	}
	_, ok := index.Objectives[id]
	return ok
}

func taskExists(index *data.V2Index, id string) bool {
	if index == nil {
		return false
	}
	_, ok := index.Tasks[id]
	return ok
}

func issueExists(index *data.V2Index, id string) bool {
	if index == nil {
		return false
	}
	_, ok := index.Issues[id]
	return ok
}

func recordExists(index *data.V2Index, kind DetailKind, id string) bool {
	if index == nil {
		return false
	}
	switch kind {
	case DetailObjective:
		return objectiveExists(index, id)
	case DetailRelease:
		return releaseExists(index, id)
	default:
		return taskExists(index, id)
	}
}

func objectiveLabel(id string) string {
	if id == "" {
		return "no Objective"
	}
	return "Objective " + id
}

func (m *Model) noteStatus(message string) {
	if strings.TrimSpace(message) == "" {
		return
	}
	if strings.TrimSpace(m.StatusMessage) == "" {
		m.StatusMessage = message
		return
	}
	m.StatusMessage += "; " + message
}

func restoredObjectiveForRelease(m Model, snapshot reloadSnapshot, state ProjectState, wasLoaded bool, releaseID string) string {
	if releaseID == "" {
		return ""
	}
	if m.ObjectiveFilter != "" || !wasLoaded {
		return selectedObjectiveForRelease(state, m.ObjectiveFilter, releaseID)
	}
	if routerObjective(state) != snapshot.RouterObjective {
		return selectedObjectiveForRelease(state, m.ObjectiveFilter, releaseID)
	}
	if snapshot.SelectedObjective == "" {
		return ""
	}
	if state.Index != nil {
		if objectiveBelongsToRelease(state.Index, snapshot.SelectedObjective, releaseID) {
			return snapshot.SelectedObjective
		}
	}
	return nearestObjectiveForRelease(state, m.ObjectiveCursor, releaseID)
}

func nearestObjectiveForRelease(state ProjectState, previousCursor int, releaseID string) string {
	rows := objectiveRowsForRelease(state.Index, releaseID)
	if len(rows) == 0 {
		return ""
	}
	if previousCursor < 0 {
		previousCursor = 0
	}
	if previousCursor >= len(rows) {
		previousCursor = len(rows) - 1
	}
	return rows[previousCursor].ID()
}

func (m *Model) restoreObjectiveCursor(snapshot reloadSnapshot, wasLoaded bool) {
	if !wasLoaded {
		m.clampObjectiveCursor()
		return
	}
	if snapshot.ObjectiveCursorID != "" {
		if cursor := m.objectiveIndex(snapshot.ObjectiveCursorID); cursor >= 0 {
			m.ObjectiveCursor = cursor
			return
		}
	}
	m.clampObjectiveCursorToRows()
}

func (m *Model) restoreReleaseCursor(snapshot reloadSnapshot, wasLoaded bool) {
	if wasLoaded && snapshot.ReleaseCursorID != "" {
		for i, id := range m.Releases {
			if id == snapshot.ReleaseCursorID {
				m.ReleaseCursor = i
				return
			}
		}
	}
	m.ReleaseCursor = releaseIndex(m.Releases, m.SelectedRelease)
}

// restoreReleaseWrite rolls back only the optimistic selector view. The
// following load still wins, so an external router edit or a removed Release
// is never hidden by this local recovery.
func (m *Model) restoreReleaseWrite() {
	if m.releaseRollback == nil {
		return
	}
	snapshot := *m.releaseRollback
	m.releaseRollback = nil
	m.SelectedRelease = snapshot.SelectedRelease
	m.SelectedObjective = snapshot.SelectedObjective
	m.Objectives = objectiveRowsForRelease(m.State.Index, m.SelectedRelease)
	m.Cards = groupTaskCardsForRelease(m.State.Index, m.SelectedRelease, m.SelectedObjective)
	m.restoreObjectiveCursor(snapshot, true)
	m.restoreFocus(snapshot, true)
	m.ReleaseCursor = releaseIndex(m.Releases, m.SelectedRelease)
}

// restoreFocus follows a Task identity across a reload. If the identity no
// longer exists, the old column and ordinal are the first fallback; when that
// column is empty the nearest non-empty column wins, preferring the column on
// the left at equal distance.
func (m *Model) restoreFocus(snapshot reloadSnapshot, wasLoaded bool) {
	if !wasLoaded || snapshot.FocusedTaskID == "" {
		m.FocusedColumn = snapshot.FocusedColumn
		m.FocusedCard = snapshot.FocusedCard
		m.clampFocus()
		return
	}
	for _, column := range columnOrder {
		for card, candidate := range m.Cards[column] {
			if candidate.Task.ID == snapshot.FocusedTaskID {
				m.FocusedColumn = column
				m.FocusedCard = card
				return
			}
		}
	}

	m.FocusedColumn = snapshot.FocusedColumn
	m.FocusedCard = snapshot.FocusedCard
	if len(m.Cards[m.FocusedColumn]) == 0 {
		m.FocusedColumn = nearestNonEmptyColumn(m.Cards, snapshot.FocusedColumn)
		m.FocusedCard = 0
	} else {
		m.clampFocus()
	}
	if m.State.Index == nil {
		return
	}
	if _, stillExists := m.State.Index.Tasks[snapshot.FocusedTaskID]; !stillExists {
		m.StatusMessage = fmt.Sprintf("Task %s no longer exists; focus moved to %s.", snapshot.FocusedTaskID, m.focusedTaskLabel())
	}
}

func nearestNonEmptyColumn(cards map[data.ColumnType][]TaskCard, from data.ColumnType) data.ColumnType {
	start := columnIndex(from)
	if start < 0 {
		start = 0
	}
	for distance := 0; distance < len(columnOrder); distance++ {
		left := start - distance
		if left >= 0 && len(cards[columnOrder[left]]) > 0 {
			return columnOrder[left]
		}
		right := start + distance
		if distance > 0 && right < len(columnOrder) && len(cards[columnOrder[right]]) > 0 {
			return columnOrder[right]
		}
	}
	return columnOrder[start]
}

func (m Model) focusedTaskLabel() string {
	if id := m.focusedTaskID(); id != "" {
		return "Task " + id
	}
	return "an empty column"
}

func (m *Model) restoreOverlayOrigins() {
	if m.Detail != nil {
		m.DetailOrigin.SidebarFocused = m.SidebarFocused
		m.DetailOrigin.ObjectiveCursor = m.ObjectiveCursor
		m.DetailOrigin.FocusedColumn = m.FocusedColumn
		m.DetailOrigin.FocusedCard = m.FocusedCard
	}
	if m.Issues != nil {
		m.Issues.Origin.SidebarFocused = m.SidebarFocused
		m.Issues.Origin.ObjectiveCursor = m.ObjectiveCursor
		m.Issues.Origin.FocusedColumn = m.FocusedColumn
		m.Issues.Origin.FocusedCard = m.FocusedCard
	}
}

// clampObjectiveCursor puts the sidebar cursor on the selected Objective when
// there is one, and otherwise keeps it inside the rows that exist — so a load
// that removed Objectives, or a list shorter than the cursor, leaves a cursor
// something renders.
func (m *Model) clampObjectiveCursor() {
	if selected := m.objectiveIndex(m.SelectedObjective); selected >= 0 {
		m.ObjectiveCursor = selected
		return
	}
	m.clampObjectiveCursorToRows()
}

// clampObjectiveCursorToRows keeps the sidebar cursor inside the rows that
// exist without moving it onto the selection, for the restores that must put
// the cursor back exactly where the reader left it.
func (m *Model) clampObjectiveCursorToRows() {
	if m.ObjectiveCursor >= len(m.Objectives) {
		m.ObjectiveCursor = len(m.Objectives) - 1
	}
	if m.ObjectiveCursor < 0 {
		m.ObjectiveCursor = 0
	}
}

// focusColumn moves focus one column left or right, clamping at both ends
// rather than wrapping, and re-clamps the card cursor into the column it lands
// on.
func (m *Model) focusColumn(delta int) {
	current := columnIndex(m.FocusedColumn)
	next := current + delta
	if next < 0 || next >= len(columnOrder) {
		return
	}
	m.FocusedColumn = columnOrder[next]
	m.clampFocus()
}

// focusCard moves the card cursor within the focused column, clamping at both
// ends.
func (m *Model) focusCard(delta int) {
	next := m.FocusedCard + delta
	if next < 0 || next >= len(m.Cards[m.FocusedColumn]) {
		return
	}
	m.FocusedCard = next
}

// clampFocus keeps the cursor inside the cards that actually exist, so a reload
// that removed the focused Task, or a column shorter than the cursor, leaves a
// valid focus rather than an index nothing renders.
func (m *Model) clampFocus() {
	if columnIndex(m.FocusedColumn) < 0 {
		m.FocusedColumn = data.ColumnPlanned
	}
	count := len(m.Cards[m.FocusedColumn])
	if m.FocusedCard >= count {
		m.FocusedCard = count - 1
	}
	if m.FocusedCard < 0 {
		m.FocusedCard = 0
	}
}

func columnIndex(column data.ColumnType) int {
	for i, candidate := range columnOrder {
		if candidate == column {
			return i
		}
	}
	return -1
}
