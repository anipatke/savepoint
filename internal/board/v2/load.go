package v2

import (
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/opencode/savepoint/internal/data"
)

// ProjectState is everything one load of a V2 project produced: the
// identity-keyed index, the decoded V2 router, and the one next action
// resolved from those two. Nothing here is derived by this package — Next is
// data.ResolveNext's own value, carried whole.
type ProjectState struct {
	Index       *data.V2Index
	Router      *data.RouterStateV2
	RouterMtime time.Time
	Issues      IssueCatalog
	Next        data.Next
}

// objectiveCount and taskCount report what the load put into the index. A nil
// index is zero of each rather than a missing value: an empty project and a
// project held back by a pending migration both count as nothing loaded, and
// both are states the board opens in normally.
func (s ProjectState) objectiveCount() int {
	_, total := s.objectiveCounts()
	return total
}

func (s ProjectState) taskCount() int {
	_, total := s.taskCounts()
	return total
}

// objectiveCounts, taskCounts, and issueCounts report how many of each record
// the load indexed and how many of those are finished: done Objectives and
// Tasks, resolved Issues.
func (s ProjectState) objectiveCounts() (finished, total int) {
	if s.Index == nil {
		return 0, 0
	}
	for _, objective := range s.Index.Objectives {
		if objective.Status == data.ColumnDone {
			finished++
		}
	}
	return finished, len(s.Index.Objectives)
}

func (s ProjectState) taskCounts() (finished, total int) {
	if s.Index == nil {
		return 0, 0
	}
	for _, task := range s.Index.Tasks {
		if task.Status == data.ColumnDone {
			finished++
		}
	}
	return finished, len(s.Index.Tasks)
}

func (s ProjectState) issueCounts() (finished, total int) {
	if s.Index == nil {
		return 0, 0
	}
	for _, issue := range s.Index.Issues {
		if issue.Status == data.IssueStatusResolved {
			finished++
		}
	}
	return finished, len(s.Index.Issues)
}

// projectLoadedMsg is the single message the load command returns, for both
// outcomes: a loaded ProjectState, or a named Diagnostic and nothing else. One
// message type keeps the reducer's load handling in one branch, and keeps a
// partially populated state from ever standing in for a failed load.
type projectLoadedMsg struct {
	State      ProjectState
	Diagnostic string
	// Seq orders results by when their load started, so a slow older load
	// cannot overwrite a newer one. Zero is unordered and always applied.
	Seq uint64
	// Retry marks the one delayed re-read that follows a failed load.
	Retry bool
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
		seq := loadSeq.Add(1)
		msg := loadProject(root)
		msg.Seq = seq
		return msg
	}
}

// loadSeq numbers loads in the order they start reading the project.
var loadSeq atomic.Uint64

// reloadRetryDelay gives an agent or editor that is part-way through a
// multi-step edit time to finish before a failed load is reported (I-031).
const reloadRetryDelay = 400 * time.Millisecond

// retryLoadCmd re-reads the project once after reloadRetryDelay. Its result
// is final: a second failure is shown rather than retried again.
func retryLoadCmd(root string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(reloadRetryDelay)
		msg := loadCmd(root)().(projectLoadedMsg)
		msg.Retry = true
		return msg
	}
}

// loadProject performs one load of the V2 project rooted at root — a
// .savepoint directory — and reports either its state or the first named
// diagnostic it hit. It is the body loadCmd wraps, called directly only where
// there is no Bubble Tea program to run it in (the non-TTY path).
//
// The board is only ever reached after the command dispatch above has
// already confirmed schema 2 (runWithFilters), so this is the ordinary path
// alone: load the index, decode the router, resolve the next action.
func loadProject(root string) projectLoadedMsg {
	index, err := data.LoadV2Index(root)
	if err != nil {
		return projectLoadedMsg{Diagnostic: err.Error()}
	}

	router, err := readRouter(root)
	if err != nil {
		return projectLoadedMsg{Diagnostic: err.Error()}
	}
	routerInfo, err := os.Stat(filepath.Join(root, "router.md"))
	if err != nil {
		return projectLoadedMsg{Diagnostic: err.Error()}
	}

	return projectLoadedMsg{State: ProjectState{
		Index:       index,
		Router:      router,
		RouterMtime: routerInfo.ModTime(),
		Issues:      issueCatalog(index),
		Next:        data.ResolveNext(data.NextInput{Index: index, Router: router}),
	}}
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
