package board

import (
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fsnotify/fsnotify"
	"github.com/opencode/savepoint/internal/data"
)

type fileChangeMsg struct{}
type reloadMsg struct {
	tasks        []data.Task
	releases     []string
	releaseEpics map[string][]string
	epicStatuses map[string]string
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
						watchCreatedDir(w, event)
					case <-timer.C:
						break drain
					}
				}
				return fileChangeMsg{}
			case _, ok := <-w.Errors:
				if !ok {
					return nil
				}
			}
		}
	}
}

func reloadTasks(root string) tea.Cmd {
	return func() tea.Msg {
		tasks, releases, releaseEpics, epicStatuses, err := loadBoardData(root)
		if err != nil {
			return nil
		}
		return reloadMsg{tasks: tasks, releases: releases, releaseEpics: releaseEpics, epicStatuses: epicStatuses}
	}
}

// newWatcher watches the releases directory by walking all subdirs (fsnotify v1.10 has no recursive opt).
func newWatcher(root string) (*fsnotify.Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	releasesPath := filepath.Join(root, "releases")
	if err := addDirsRecursive(w, releasesPath); err != nil {
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
