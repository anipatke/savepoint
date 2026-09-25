// Package v2 is the board for a schema_version 2 project: Objectives, Task
// outcomes, and the one next action, read through internal/data's V2 index and
// resolvers. It is a sibling of the V1 board rather than a second personality
// inside it, so no V1 record type is reachable from a V2 render path.
package v2

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fsnotify/fsnotify"
	"github.com/opencode/savepoint/internal/data"
)

// Model is the V2 board's state. Root and ObjectiveFilter are the invocation
// it was given; State is the last completed load; Diagnostic is the named
// reason a load produced no project at all, and while it is set no board is
// drawn from project data.
type Model struct {
	// Root is the project's .savepoint directory.
	Root string
	// ObjectiveFilter is the --objective selection the command was invoked
	// with, resolved against the loaded index by exact global ID only.
	ObjectiveFilter string
	// Watcher is the V2 file-model watcher owned by the TUI run. It is nil in
	// plain mode and in unit models, where reloads are driven by messages.
	Watcher *fsnotify.Watcher

	State      ProjectState
	Loaded     bool
	Diagnostic string
	// ReloadDiagnostic keeps the last good State on screen when an external
	// edit is temporarily malformed. It is cleared by the next successful
	// load; Diagnostic is reserved for the initial no-board screen.
	ReloadDiagnostic string
	// lastLoadSeq is the Seq of the newest load result applied.
	lastLoadSeq uint64

	// FatalErr is set when the board cannot honor the invocation it was given
	// — today, an --objective naming no Objective in the loaded project. Run
	// returns it so the command exits nonzero instead of showing a board that
	// quietly ignored a flag.
	FatalErr error

	// SelectedObjective is the Objective in view: the --objective filter when
	// one was given, otherwise the router's selection when it resolves. It
	// filters the columns and nothing else — selecting a row here writes
	// nothing to router.md.
	SelectedObjective string

	// SelectedRelease is the optional delivery context in view. Its member
	// Objectives and their Tasks are derived from the loaded index; changing
	// this value never changes a lifecycle record.
	SelectedRelease string
	// Releases is the stable, ascending list of live R-### identities offered
	// by the selector. Titles remain on the indexed records and are looked up
	// only when the selector renders.
	Releases      []string
	ReleaseCursor int
	// ReleaseOverlay keeps the board visible behind the selector. Its origin
	// restores the surface that had focus when r opened it.
	ReleaseOverlay   bool
	releaseOrigin    detailOrigin
	releaseOriginSet bool
	// releaseRollback is set only while an optimistic Release selection is
	// waiting for the canonical router writer. A refusal can therefore restore
	// the prior in-memory view before the truthful reload begins.
	releaseRollback      *reloadSnapshot
	preserveReloadStatus bool

	// Objectives are the sidebar's rows in stable O-### order, each carrying the
	// values its badges read. Like Cards they are rebuilt by every load.
	Objectives []ObjectiveRow
	// ObjectiveCursor is the sidebar's cursor. Like the card cursor it selects
	// nothing by itself: moving it changes what is highlighted, not what the
	// columns show.
	ObjectiveCursor int
	// SidebarFocused reports which surface the keys drive, the sidebar or the
	// columns. Focus changes accents and glyphs only; no geometry depends on it.
	SidebarFocused bool

	// Cards are the loaded Tasks grouped by recorded status, each carrying the
	// resolved clearance and gate decision its badges read. They are rebuilt
	// by every load, so nothing on a card can outlive the records behind it.
	Cards map[data.ColumnType][]TaskCard
	// FocusedColumn and FocusedCard are the cursor: which column has focus and
	// which card within it. Focus is presentation state only — it selects
	// nothing in the project and writes nothing.
	FocusedColumn data.ColumnType
	FocusedCard   int

	// Detail is the open detail overlay, or nil when none is open. It holds
	// values resolved at the moment it opened and re-resolved by every later
	// load, so it can never show a record the project no longer has.
	Detail *RecordDetail
	// DetailOffset is how far the open overlay is scrolled.
	DetailOffset int
	// DetailOrigin is the surface the overlay was opened from, restored when it
	// closes so a reader returns to where they were rather than to the top of
	// the board.
	DetailOrigin detailOrigin

	// Issues is the one read-only follow-up overlay. It is mutually exclusive
	// with Detail and restores the board cursor it replaced when closed.
	Issues *IssueOverlay
	// Help is the keyboard reference overlay. It changes no cursor or project
	// state, and its action rows are derived from the currently focused record.
	Help bool

	Width         int
	Height        int
	StatusMessage string
}

// detailOrigin is the focus and cursor state one overlay was opened over.
// Opening a detail changes none of it; closing restores it, clamped into
// whatever records a reload underneath the overlay left behind.
type detailOrigin struct {
	SidebarFocused  bool
	ObjectiveCursor int
	FocusedColumn   data.ColumnType
	FocusedCard     int
}

// NewModel builds the V2 board model for one invocation. It performs no IO:
// every value the board displays arrives through the load command.
func NewModel(opts Options) Model {
	return Model{
		Root:            opts.Root,
		ObjectiveFilter: opts.ObjectiveFilter,
		Cards:           groupTaskCards(nil),
		FocusedColumn:   data.ColumnPlanned,
	}
}

func (m Model) Init() tea.Cmd {
	if m.Watcher == nil {
		return loadCmd(m.Root)
	}
	return tea.Batch(loadCmd(m.Root), watchV2Files(m.Watcher, m.Root))
}

// sidebarVisible reports whether the Objective sidebar is drawn at the current
// terminal size. The keys follow the screen: a terminal too narrow to carry the
// sidebar beside three columns does not move focus into a surface nothing
// draws.
func (m Model) sidebarVisible() bool {
	return m.terminalWidth() >= sidebarBreakpoint
}

// objectiveIndex is the sidebar row holding id, or -1 when no row does — which
// is the answer for an empty project and for a router naming a record the index
// no longer has.
func (m Model) objectiveIndex(id string) int {
	if id == "" {
		return -1
	}
	for i, row := range m.Objectives {
		if row.ID() == id {
			return i
		}
	}
	return -1
}

// objectiveFilterError names an --objective selection the loaded project has no
// record for. Resolution is an exact global-ID lookup: no prefix, numeric
// proximity, or title match, and nothing similarly named is offered in its
// place. A project with no index to resolve against — one held back by a
// pending migration — is not judged here; the migration is the report.
func objectiveFilterError(index *data.V2Index, filter string) error {
	if filter == "" || index == nil {
		return nil
	}
	if _, ok := index.Objectives[filter]; ok {
		return nil
	}
	return fmt.Errorf("--objective %s names no objective in this project", filter)
}

func selectedObjectiveForRelease(state ProjectState, filter, releaseID string) string {
	selected := selectedObjectiveUnscoped(state, filter)
	if selected == "" || releaseID == "" {
		return selected
	}
	if objectiveBelongsToRelease(state.Index, selected, releaseID) {
		return selected
	}
	return ""
}

func selectedObjectiveUnscoped(state ProjectState, filter string) string {
	if filter != "" {
		return filter
	}
	objective := routerObjective(state)
	if state.Index == nil || objective == "" {
		return ""
	}
	if _, ok := state.Index.Objectives[objective]; !ok {
		return ""
	}
	return objective
}

// routerObjective is the Objective the router's selection puts in view: its
// selected Objective, or, when it selects an Issue alone, the one Objective
// that Issue's linked Tasks and Objective-scoped Checks belong to. An Issue
// linked to no Objective, or to more than one, puts none in view rather than
// guessing between them.
func routerObjective(state ProjectState) string {
	router := state.Router
	if router == nil {
		return ""
	}
	if router.Objective != "" || router.Issue == "" || state.Index == nil {
		return router.Objective
	}
	issue, ok := state.Index.Issues[router.Issue]
	if !ok {
		return ""
	}
	owners := map[string]bool{}
	for _, taskID := range issue.Tasks {
		if task, ok := state.Index.Tasks[taskID]; ok {
			owners[task.Objective] = true
		}
	}
	for _, checkID := range issue.Checks {
		if check, ok := state.Index.Checks[checkID]; ok && check.Scope.Kind == data.CheckScopeObjective {
			owners[check.Scope.ID] = true
		}
	}
	if len(owners) != 1 {
		return ""
	}
	for objective := range owners {
		return objective
	}
	return ""
}
