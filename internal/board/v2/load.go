package v2

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/migrate"
)

// ProjectState is everything one load of a V2 project produced: the
// identity-keyed index, the decoded V2 router, the injected migration state
// with the guidance that explains it, and the one next action resolved from
// those three. Nothing here is derived by this package — Next is
// data.ResolveNext's own value, carried whole.
type ProjectState struct {
	Index     *data.V2Index
	Router    *data.RouterStateV2
	Issues    IssueCatalog
	Migration data.MigrationState
	// MigrationGuidance is migrate.PendingOperation's own explanation of the
	// incomplete conversion, empty when none is pending.
	MigrationGuidance string
	Next              data.Next
}

// objectiveCount and taskCount report what the load put into the index. A nil
// index is zero of each rather than a missing value: an empty project and a
// project held back by a pending migration both count as nothing loaded, and
// both are states the board opens in normally.
func (s ProjectState) objectiveCount() int {
	if s.Index == nil {
		return 0
	}
	return len(s.Index.Objectives)
}

func (s ProjectState) taskCount() int {
	if s.Index == nil {
		return 0
	}
	return len(s.Index.Tasks)
}

// projectLoadedMsg is the single message the load command returns, for both
// outcomes: a loaded ProjectState, or a named Diagnostic and nothing else. One
// message type keeps the reducer's load handling in one branch, and keeps a
// partially populated state from ever standing in for a failed load.
type projectLoadedMsg struct {
	State      ProjectState
	Diagnostic string
}

// Failed reports whether the load produced a diagnostic instead of a project.
func (m projectLoadedMsg) Failed() bool {
	return m.Diagnostic != ""
}

// loadCmd is the board's only load path. Startup and every later reload run
// this same command, so there is no startup path and refresh path that can
// disagree, and no filesystem access happens in Update (ARCH-02).
func loadCmd(root string) tea.Cmd {
	return func() tea.Msg {
		return loadProject(root)
	}
}

// loadProject performs one load of the V2 project rooted at root — a
// .savepoint directory — and reports either its state or the first named
// diagnostic it hit. It is the body loadCmd wraps, called directly only where
// there is no Bubble Tea program to run it in (the non-TTY path).
//
// A pending migration is read first and answered on its own: a project midway
// through conversion holds records whose meaning is not yet settled, so they
// are not loaded and not interpreted, exactly as `savepoint resume` treats the
// same state. Everything else is the ordinary path: load the index, decode the
// router, resolve the next action.
func loadProject(root string) projectLoadedMsg {
	migration, guidance, err := pendingMigration(root)
	if err != nil {
		return projectLoadedMsg{Diagnostic: err.Error()}
	}
	if migration.Pending {
		return projectLoadedMsg{State: ProjectState{
			Migration:         migration,
			MigrationGuidance: guidance,
			Next:              data.ResolveNext(data.NextInput{Migration: migration}),
		}}
	}

	project, err := data.LoadProject(root)
	if err != nil {
		return projectLoadedMsg{Diagnostic: err.Error()}
	}

	router, err := readRouter(root)
	if err != nil {
		return projectLoadedMsg{Diagnostic: err.Error()}
	}

	return projectLoadedMsg{State: ProjectState{
		Index:  project.V2,
		Router: router,
		Issues: issueCatalog(project.V2),
		Next:   data.ResolveNext(data.NextInput{Index: project.V2, Router: router}),
	}}
}

// pendingMigration reads the incomplete-migration detector the V1 board,
// upgrade-assets, and doctor all consult, so the V2 board reports the same
// boundary rather than growing a second policy. root is the project's
// .savepoint directory; the project directory PendingOperation expects is one
// level up.
func pendingMigration(root string) (data.MigrationState, string, error) {
	report, err := migrate.PendingOperation(filepath.Dir(root))
	if err != nil {
		return data.MigrationState{}, "", err
	}
	if report == nil {
		return data.MigrationState{}, "", nil
	}
	return data.MigrationState{Pending: true, OperationID: report.OperationID}, report.RecoveryGuidance(), nil
}

// readRouter reads and strictly decodes the project's V2 router. Both an
// unreadable file and a router the decoder refuses are returned naming the
// file, so neither becomes a healed default (DATA-03).
func readRouter(root string) (*data.RouterStateV2, error) {
	path := filepath.Join(root, "router.md")
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	router, err := data.NewRouterReader().ReadStateV2(string(content))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return router, nil
}
