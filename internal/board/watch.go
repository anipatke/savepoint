package board

import (
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fsnotify/fsnotify"
	"github.com/opencode/savepoint/internal/data"
)

// auditWatchDir is the audit-register subtree the board watches, mirroring the
// data package's audit/ directory layout so audit edits trigger a reload.
const auditWatchDir = "audit"

type fileChangeMsg struct{}
type reloadMsg struct {
	tasks        []data.Task
	defects      []data.Defect
	releases     []string
	releaseEpics map[string][]string
	epicStatuses map[string]string
	audit        data.AuditRegisterSet
	routerState  *data.RouterState
	message      string
}

// auditRegisterMsg carries a freshly loaded audit-register data set for the
// Audit Register overlay, produced by loadAuditRegisterCmd on overlay open.
type auditRegisterMsg struct {
	set data.AuditRegisterSet
}

type routerWriteMsg struct {
	message string
	state   *data.RouterState
	taskID  string
}

type taskWriteMsg struct {
	prefix string
	next   data.Task
	err    error
}

type defectWriteMsg struct {
	next data.Defect
}

type taskRefreshMsg struct {
	task    data.Task
	message string
}

type epicDetailMsg struct {
	content string
}

type auditContentMsg struct {
	content string
}

type epicStatusWrittenMsg struct {
	epicID string
	status string
}

// releaseDocsMsg carries the supporting documents loaded for the top-level
// Release Docs overlay (OverlayReleaseDocs).
type releaseDocsMsg struct {
	docs []data.ReleaseDoc
}

type errorMsg struct {
	message string
}

// watchFiles blocks until a file event arrives, debounces for 100ms, emits fileChangeMsg.
func watchFiles(w *fsnotify.Watcher) tea.Cmd {
	return func() tea.Msg {
		for {
			select {
			case event, ok := <-w.Events:
				if !ok {
					return nil
				}
				debugf("watcher: event %s", event)
				watchCreatedDir(w, event)
				timer := time.NewTimer(100 * time.Millisecond)
			drain:
				for {
					select {
					case event, ok := <-w.Events:
						if !ok {
							timer.Stop()
							return nil
						}
						debugf("watcher: event %s", event)
						watchCreatedDir(w, event)
					case <-timer.C:
						break drain
					}
				}
				debugf("watcher: emitting fileChangeMsg")
				return fileChangeMsg{}
			case _, ok := <-w.Errors:
				if !ok {
					return nil
				}
			}
		}
	}
}

func reloadTasks(root string, deps ModelDependencies) tea.Cmd {
	return reloadTasksWithMessage(root, deps, "")
}

func reloadTasksWithMessage(root string, deps ModelDependencies, message string) tea.Cmd {
	return func() tea.Msg {
		debugf("reload: starting task reload from %q", root)
		tasks, defects, releases, releaseEpics, epicStatuses, err := loadBoardData(root, deps.Discoverer, deps.Parser)
		if err != nil {
			debugf("reload: error: %v", err)
			return errorMsg{message: "reload failed: " + err.Error()}
		}
		debugf("reload: loaded %d tasks", len(tasks))
		routerState, _ := readRouterState(root, deps.RouterReader)
		audit := loadAuditBestEffort(root, deps.AuditLoader)
		return reloadMsg{tasks: tasks, defects: defects, releases: releases, releaseEpics: releaseEpics, epicStatuses: epicStatuses, audit: audit, routerState: routerState, message: message}
	}
}

// newWatcher watches the savepoint root (for router.md) and all releases subdirs.
func newWatcher(root string) (*fsnotify.Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	if err := w.Add(root); err != nil {
		w.Close()
		return nil, err
	}
	releasesPath := filepath.Join(root, "releases")
	if err := addDirsRecursive(w, releasesPath); err != nil {
		w.Close()
		return nil, err
	}
	// Watch the audit-register subtree so finding/run/prompt/register edits
	// drive the same reload path as task and release changes. A missing audit/
	// directory is skipped, so projects without an audit tree still start.
	auditPath := filepath.Join(root, auditWatchDir)
	if err := addDirsRecursive(w, auditPath); err != nil {
		w.Close()
		return nil, err
	}
	return w, nil
}

func watchCreatedDir(w *fsnotify.Watcher, event fsnotify.Event) {
	if !event.Has(fsnotify.Create) {
		return
	}
	info, err := os.Stat(event.Name)
	if err != nil || !info.IsDir() {
		return
	}
	_ = addDirsRecursive(w, event.Name)
}

func addDirsRecursive(w *fsnotify.Watcher, root string) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // skip unreadable dirs
		}
		if d.IsDir() {
			return w.Add(path)
		}
		return nil
	})
}
