package v2

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/resume"
)

// issueFilterOrder is the complete filter vocabulary. The empty value is the
// unfiltered view; the remaining values are descriptive types, never a
// severity or a gate decision.
var issueFilterOrder = []data.IssueType{
	"",
	"defect",
	"drift",
	"guardrail",
	"verification",
	"other",
}

// issueColumnLabels is the three-column vocabulary the Issues overlay always
// renders in, keyed by the recorded Issue status it counts — the same
// grey/orange/green accents the Task board wears for planned/in-progress/done,
// carried over onto open/in_progress/resolved so the two surfaces read as one
// color language.
var issueColumnLabels = []struct {
	Label  string
	Status data.IssueStatus
}{
	{"OPEN", data.IssueStatusOpen},
	{"IN PROGRESS", data.IssueStatusInProgress},
	{"RESOLVED", data.IssueStatusResolved},
}

func issueColumnIndex(status data.IssueStatus) int {
	for i, column := range issueColumnLabels {
		if column.Status == status {
			return i
		}
	}
	return -1
}

// IssueLink is a resolved reference shown by an Issue detail. The loader
// settles its label once, so the renderer never looks through the index.
type IssueLink struct {
	ID     string
	Label  string
	Detail string
}

// IssueRow is the complete list projection for one Issue. Links are carried
// from the Issue's own reference lists and the index's reverse maps; no view
// walks records to reconstruct a relation.
type IssueRow struct {
	Issue  *data.IssueV2
	Tasks  []IssueLink
	Checks []IssueLink
}

// IssueCatalog is built by the load command and reused by every Issues view.
// Rows are stable I### order. TaskRows is the scoped entry-point projection,
// including direct Task links and Issues attached to that Task's Check chain.
type IssueCatalog struct {
	Rows     []IssueRow
	ByID     map[string]IssueRow
	TaskRows map[string][]IssueRow
}

// IssueDetail is a fully resolved, read-only Issue overlay value. Body and
// history are carried exactly as recorded; the view only formats them.
type IssueDetail struct {
	Issue            *data.IssueV2
	Tasks            []IssueLink
	Checks           []IssueLink
	GuardrailIDs     []string
	DuplicateTarget  *IssueLink
	EscalationTarget *IssueLink
}

// issueOrigin records the board cursor the Issues overlay replaces, so close
// restores the reader to the same Task or Objective focus.
type issueOrigin struct {
	SidebarFocused  bool
	ObjectiveCursor int
	FocusedColumn   data.ColumnType
	FocusedCard     int
}

// IssueOverlay is the read-only Issues surface: one type filter, one status
// column in focus, one cursor inside it, optional detail, and a stack for
// duplicate-target navigation.
type IssueOverlay struct {
	Filter        data.IssueType
	FocusedStatus data.IssueStatus
	Cursor        int
	SelectedID    string
	Detail        *IssueDetail
	DetailOffset  int
	DetailStack   []*IssueDetail
	ScopedTask    string
	Origin        issueOrigin
}

// issueCatalog resolves every Issue and the two scoped link directions needed
// by the overlay. The Issue's own Tasks and Checks remain authoritative for
// detail links; the index maps answer which Issues belong to a focused Task.
func issueCatalog(index *data.V2Index) IssueCatalog {
	catalog := IssueCatalog{
		ByID:     map[string]IssueRow{},
		TaskRows: map[string][]IssueRow{},
	}
	if index == nil {
		return catalog
	}

	for _, id := range slices.Sorted(maps.Keys(index.Issues)) {
		issue := index.Issues[id]
		row := IssueRow{
			Issue:  issue,
			Tasks:  issueTaskLinks(index, issue.Tasks),
			Checks: issueCheckLinks(index, issue.Checks),
		}
		catalog.Rows = append(catalog.Rows, row)
		catalog.ByID[id] = row
	}

	for _, taskID := range slices.Sorted(maps.Keys(index.Tasks)) {
		issueIDs := index.TaskIssues[taskID]
		ids := append([]string(nil), issueIDs...)
		for _, checkID := range index.ScopeChecks[taskID] {
			ids = append(ids, index.CheckIssues[checkID]...)
		}
		ids = uniqueSortedIDs(ids)
		for _, issueID := range ids {
			if row, ok := catalog.ByID[issueID]; ok {
				catalog.TaskRows[taskID] = append(catalog.TaskRows[taskID], row)
			}
		}
	}

	return catalog
}

func uniqueSortedIDs(ids []string) []string {
	if len(ids) == 0 {
		return nil
	}
	slices.Sort(ids)
	return slices.Compact(ids)
}

func issueTaskLinks(index *data.V2Index, ids []string) []IssueLink {
	links := make([]IssueLink, 0, len(ids))
	for _, id := range ids {
		if task, ok := index.Tasks[id]; ok {
			links = append(links, IssueLink{ID: id, Label: task.Title, Detail: fmt.Sprintf("%s", task.Status)})
			continue
		}
		links = append(links, IssueLink{ID: id})
	}
	return links
}

func issueCheckLinks(index *data.V2Index, ids []string) []IssueLink {
	links := make([]IssueLink, 0, len(ids))
	for _, id := range ids {
		if check, ok := index.Checks[id]; ok {
			links = append(links, IssueLink{ID: id, Label: string(check.Result), Detail: check.Scope.ID})
			continue
		}
		links = append(links, IssueLink{ID: id})
	}
	return links
}

func (m Model) issueRows() []IssueRow {
	if m.Issues == nil {
		return nil
	}
	if m.Issues.ScopedTask != "" {
		return m.State.Issues.TaskRows[m.Issues.ScopedTask]
	}
	return m.State.Issues.Rows
}

func (m Model) filteredIssueRows() []IssueRow {
	rows := m.issueRows()
	if m.Issues == nil || m.Issues.Filter == "" {
		return rows
	}
	filtered := make([]IssueRow, 0, len(rows))
	for _, row := range rows {
		if row.Issue.Type == m.Issues.Filter {
			filtered = append(filtered, row)
		}
	}
	return filtered
}

// groupedIssueRows buckets the scoped, type-filtered rows by status into the
// three columns the overlay renders, then sorts each bucket most-severe
// first (issueSeverityRank) so a reader scans blockers before cosmetic
// follow-ups within a column. The sort is stable, so Issues at the same
// rank — including every Issue when none records a severity at all — keep
// the underlying stable ID order, and a project renders the same way twice.
func (m Model) groupedIssueRows() map[data.IssueStatus][]IssueRow {
	grouped := map[data.IssueStatus][]IssueRow{
		data.IssueStatusOpen:       {},
		data.IssueStatusInProgress: {},
		data.IssueStatusResolved:   {},
	}
	for _, row := range m.filteredIssueRows() {
		grouped[row.Issue.Status] = append(grouped[row.Issue.Status], row)
	}
	for status, rows := range grouped {
		sorted := rows
		slices.SortStableFunc(sorted, func(a, b IssueRow) int {
			return issueSeverityRank(a.Issue.Severity) - issueSeverityRank(b.Issue.Severity)
		})
		grouped[status] = sorted
	}
	return grouped
}

// focusedIssueRows is the rows of whichever status column is focused — the
// only rows a list-mode key (cursor move, enter) ever acts on.
func (m Model) focusedIssueRows() []IssueRow {
	if m.Issues == nil {
		return nil
	}
	return m.groupedIssueRows()[m.Issues.FocusedStatus]
}

func issueFilterLabel(filter data.IssueType) string {
	if filter == "" {
		return "ALL"
	}
	return strings.ToUpper(string(filter))
}

func (m *Model) openIssues(taskID string) {
	if m.State.Index == nil {
		return
	}
	m.Issues = &IssueOverlay{
		ScopedTask:    taskID,
		FocusedStatus: data.IssueStatusOpen,
		Origin: issueOrigin{
			SidebarFocused:  m.SidebarFocused,
			ObjectiveCursor: m.ObjectiveCursor,
			FocusedColumn:   m.FocusedColumn,
			FocusedCard:     m.FocusedCard,
		},
	}
	m.clampIssueCursor()
}

func (m *Model) closeIssues() {
	if m.Issues == nil {
		return
	}
	origin := m.Issues.Origin
	m.Issues = nil
	m.SidebarFocused = origin.SidebarFocused && m.sidebarVisible()
	m.ObjectiveCursor = origin.ObjectiveCursor
	m.FocusedColumn = origin.FocusedColumn
	m.FocusedCard = origin.FocusedCard
	m.clampObjectiveCursorToRows()
	m.clampFocus()
}

// clampIssueCursor keeps the overlay showing a real row. If the previously
// selected Issue still exists anywhere in view — even under a different
// status now, say a Check just resolved it — focus follows it to that
// column; otherwise the cursor is clamped inside whatever the focused status
// column holds. This is the restore path: it runs on open and after a
// reload, never on a reader's own left/right column choice.
func (m *Model) clampIssueCursor() {
	if m.Issues == nil {
		return
	}
	if issueColumnIndex(m.Issues.FocusedStatus) < 0 {
		m.Issues.FocusedStatus = data.IssueStatusOpen
	}
	if m.Issues.SelectedID != "" {
		grouped := m.groupedIssueRows()
		for _, column := range issueColumnLabels {
			for i, row := range grouped[column.Status] {
				if row.Issue.ID == m.Issues.SelectedID {
					m.Issues.FocusedStatus = column.Status
					m.Issues.Cursor = i
					return
				}
			}
		}
	}
	m.clampIssueCursorToStatus()
}

// clampIssueCursorToStatus clamps the cursor inside the focused status
// column alone, without searching other columns for SelectedID — the reader
// just chose this column or this filter, so a match elsewhere must never
// pull focus back to where it was.
func (m *Model) clampIssueCursorToStatus() {
	rows := m.focusedIssueRows()
	if len(rows) == 0 {
		m.Issues.Cursor = 0
		m.Issues.SelectedID = ""
		return
	}
	if m.Issues.Cursor < 0 {
		m.Issues.Cursor = 0
	}
	if m.Issues.Cursor >= len(rows) {
		m.Issues.Cursor = len(rows) - 1
	}
	m.Issues.SelectedID = rows[m.Issues.Cursor].Issue.ID
}

func (m *Model) moveIssueCursor(delta int) {
	if m.Issues == nil {
		return
	}
	rows := m.focusedIssueRows()
	next := m.Issues.Cursor + delta
	if next < 0 || next >= len(rows) {
		return
	}
	m.Issues.Cursor = next
	m.Issues.SelectedID = rows[next].Issue.ID
}

// focusIssueStatus moves the focused column left or right, clamping at both
// ends rather than wrapping, and re-clamps the cursor into the column it
// lands on without moving it to the top — the same left/right feel the Task
// board's own columns already have.
func (m *Model) focusIssueStatus(delta int) {
	if m.Issues == nil {
		return
	}
	current := issueColumnIndex(m.Issues.FocusedStatus)
	next := current + delta
	if next < 0 || next >= len(issueColumnLabels) {
		return
	}
	m.Issues.FocusedStatus = issueColumnLabels[next].Status
	m.clampIssueCursorToStatus()
}

func (m *Model) cycleIssueFilter() {
	if m.Issues == nil {
		return
	}
	current := 0
	for i, filter := range issueFilterOrder {
		if filter == m.Issues.Filter {
			current = i
			break
		}
	}
	m.Issues.Filter = issueFilterOrder[(current+1)%len(issueFilterOrder)]
	m.Issues.Cursor = 0
	m.clampIssueCursorToStatus()
}

func (m *Model) openSelectedIssue() {
	if m.Issues == nil {
		return
	}
	rows := m.focusedIssueRows()
	if m.Issues.Cursor < 0 || m.Issues.Cursor >= len(rows) {
		return
	}
	detail, ok := m.issueDetail(rows[m.Issues.Cursor].Issue.ID)
	if !ok {
		return
	}
	m.Issues.Detail = &detail
	m.Issues.DetailOffset = 0
	m.Issues.DetailStack = nil
}

func (m *Model) closeIssueDetail() {
	if m.Issues == nil || m.Issues.Detail == nil {
		return
	}
	if count := len(m.Issues.DetailStack); count > 0 {
		m.Issues.Detail = m.Issues.DetailStack[count-1]
		m.Issues.DetailStack = m.Issues.DetailStack[:count-1]
		m.Issues.DetailOffset = 0
		return
	}
	m.Issues.Detail = nil
	m.Issues.DetailOffset = 0
}

func (m *Model) openDuplicateTarget() {
	if m.Issues == nil || m.Issues.Detail == nil || m.Issues.Detail.DuplicateTarget == nil {
		return
	}
	target := m.Issues.Detail.DuplicateTarget.ID
	detail, ok := m.issueDetail(target)
	if !ok {
		return
	}
	m.Issues.DetailStack = append(m.Issues.DetailStack, m.Issues.Detail)
	m.Issues.Detail = &detail
	m.Issues.DetailOffset = 0
}

// scrollIssues and clampIssueScroll only ever move the Issue detail's own
// scroll offset: the list beneath it has no scroll offset of its own, since
// each status column's visible window follows its cursor exactly as a Task
// column's does.
func (m *Model) scrollIssues(delta int) {
	if m.Issues == nil || m.Issues.Detail == nil {
		return
	}
	_, height := m.detailViewport()
	limit := issueDetailScrollLimit(*m.Issues.Detail, m.terminalWidth(), height)
	next := m.Issues.DetailOffset + delta
	if next >= 0 && next <= limit {
		m.Issues.DetailOffset = next
	}
}

func (m *Model) clampIssueScroll() {
	if m.Issues == nil || m.Issues.Detail == nil {
		return
	}
	_, height := m.detailViewport()
	limit := issueDetailScrollLimit(*m.Issues.Detail, m.terminalWidth(), height)
	if m.Issues.DetailOffset > limit {
		m.Issues.DetailOffset = limit
	}
}

func (m *Model) handleIssuesKey(key string) {
	if m.Issues == nil {
		return
	}
	if m.Issues.Detail != nil {
		switch key {
		case "up", "k":
			m.scrollIssues(-1)
		case "down", "j":
			m.scrollIssues(1)
		case "esc":
			m.closeIssueDetail()
		case "enter", "d":
			m.openDuplicateTarget()
		}
		return
	}
	switch key {
	case "left", "h":
		m.focusIssueStatus(-1)
	case "right", "l":
		m.focusIssueStatus(1)
	case "up", "k":
		m.moveIssueCursor(-1)
	case "down", "j":
		m.moveIssueCursor(1)
	case "f", "/":
		m.cycleIssueFilter()
	case "enter", "v":
		m.openSelectedIssue()
	case "esc":
		m.closeIssues()
	}
}

func (m Model) issueDetail(id string) (IssueDetail, bool) {
	row, ok := m.State.Issues.ByID[id]
	if !ok || row.Issue == nil {
		return IssueDetail{}, false
	}
	detail := IssueDetail{
		Issue:        row.Issue,
		Tasks:        row.Tasks,
		Checks:       row.Checks,
		GuardrailIDs: append([]string(nil), row.Issue.GuardrailIDs...),
	}
	if row.Issue.DuplicateOf != "" {
		if target, exists := m.State.Issues.ByID[row.Issue.DuplicateOf]; exists {
			detail.DuplicateTarget = &IssueLink{ID: target.Issue.ID, Label: target.Issue.Title}
		}
	}
	if row.Issue.EscalatedTo != "" && m.State.Index != nil {
		if target, exists := m.State.Index.Objectives[row.Issue.EscalatedTo]; exists {
			detail.EscalationTarget = &IssueLink{ID: target.ID, Label: target.Title}
		}
	}
	return detail, true
}

func (m *Model) refreshIssues() {
	if m.Issues == nil {
		return
	}
	if m.Issues.Detail != nil {
		id := m.Issues.Detail.Issue.ID
		detail, ok := m.issueDetail(id)
		if !ok {
			m.Issues.Detail = nil
			m.Issues.DetailStack = nil
		} else {
			m.Issues.Detail = &detail
		}
	}
	m.clampIssueCursor()
	m.clampIssueScroll()
}

func (m Model) focusedTaskID() string {
	cards := m.Cards[m.FocusedColumn]
	if m.FocusedCard < 0 || m.FocusedCard >= len(cards) {
		return ""
	}
	return cards[m.FocusedCard].Task.ID
}

func issueHistoryLine(entry data.IssueHistoryEntry) string {
	line := fmt.Sprintf("%s  %s  %s", entry.At.Format(time.RFC3339), resume.ActorLabel(entry.Actor), entry.Kind)
	if entry.Check != "" {
		line += "  Check " + entry.Check
	}
	if strings.TrimSpace(entry.Note) != "" {
		line += " — " + entry.Note
	}
	return line
}
