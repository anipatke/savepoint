package v2

import (
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fsnotify/fsnotify"
)

const v2WatchDebounce = 100 * time.Millisecond

// These are the live V2 roots. They are data rather than a collection of
// path tests spread through the watcher so the file model stays visible in
// one place.
var v2WatchDirectories = []string{"objectives", "checks", "issues", "releases"}
var v2WatchFiles = []string{"router.md", "config.yml"}

// These trees contain historical or V1 material. They are deliberately not
// registered recursively, and events reaching the root watcher are rejected
// by isV2WatchedPath before they can request a reload.
var v2ExcludedDirectories = []string{"archive", "defects", "audit", "migrations"}

type v2FileChangeMsg struct{}

// newV2Watcher registers the V2 file model. The root itself catches direct
// config/router changes and the creation of a previously absent watched
// directory; live record trees are registered recursively below it.
func newV2Watcher(root string) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	if err := watcher.Add(root); err != nil {
		_ = watcher.Close()
		return nil, err
	}
	for _, directory := range v2WatchDirectories {
		path := filepath.Join(root, directory)
		if err := addV2DirsRecursive(watcher, path); err != nil {
			_ = watcher.Close()
			return nil, err
		}
	}
	return watcher, nil
}

// addV2DirsRecursive skips an absent optional tree. A newly created tree is
// picked up by watchV2CreatedDir when the root or its watched parent reports
// the create event.
func addV2DirsRecursive(watcher *fsnotify.Watcher, root string) error {
	return filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return filepath.SkipDir
			}
			return err
		}
		if !entry.IsDir() {
			return nil
		}
		return watcher.Add(path)
	})
}

// watchV2Files waits for a relevant filesystem event, drains the burst, and
// returns one message. The timer resets after each relevant event, so an
// editor's temporary-file/write/rename sequence produces one load after the
// filesystem has been quiet for the debounce interval.
func watchV2Files(watcher *fsnotify.Watcher, root string) tea.Cmd {
	return func() tea.Msg {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return nil
				}
				if !v2EventNeedsReload(root, event) {
					continue
				}
				watchV2CreatedDir(watcher, root, event)
				if debounceV2Events(watcher, root, v2WatchDebounce) {
					return v2FileChangeMsg{}
				}
				return nil
			case _, ok := <-watcher.Errors:
				if !ok {
					return nil
				}
			}
		}
	}
}

// debounceV2Events drains relevant events until the quiet interval elapses.
// It returns false when the watcher closes while the burst is being drained.
func debounceV2Events(watcher *fsnotify.Watcher, root string, delay time.Duration) bool {
	timer := time.NewTimer(delay)
	defer timer.Stop()

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return false
			}
			if !v2EventNeedsReload(root, event) {
				continue
			}
			watchV2CreatedDir(watcher, root, event)
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			timer.Reset(delay)
		case <-timer.C:
			return true
		case _, ok := <-watcher.Errors:
			if !ok {
				return false
			}
		}
	}
}

func v2EventNeedsReload(root string, event fsnotify.Event) bool {
	const relevant = fsnotify.Create | fsnotify.Write | fsnotify.Remove | fsnotify.Rename | fsnotify.Chmod
	if event.Op&relevant == 0 || !isV2WatchedPath(root, event.Name) {
		return false
	}
	// Windows reports a directory's metadata updates as writes to it, including
	// the ones caused by registering the watch itself. A real change inside the
	// directory also reports its own file, so a write or chmod on a directory
	// carries no information and would reload the board for nothing.
	if event.Op&(fsnotify.Create|fsnotify.Remove|fsnotify.Rename) == 0 {
		if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
			return false
		}
	}
	return true
}

// isV2WatchedPath maps an event to the V2 file model. It is intentionally
// lexical: no excluded directory is ever registered or followed merely
// because it exists under the same .savepoint root.
func isV2WatchedPath(root, path string) bool {
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return false
	}
	pathAbs, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	rel, err := filepath.Rel(rootAbs, pathAbs)
	if err != nil || rel == "." || rel == ".." || filepath.IsAbs(rel) {
		return false
	}
	parts := splitV2Path(rel)
	if len(parts) == 0 {
		return false
	}
	for _, excluded := range v2ExcludedDirectories {
		if parts[0] == excluded {
			return false
		}
	}
	if len(parts) == 1 {
		for _, file := range v2WatchFiles {
			if parts[0] == file {
				return true
			}
		}
	}
	for _, directory := range v2WatchDirectories {
		if parts[0] == directory {
			return true
		}
	}
	return false
}

func splitV2Path(path string) []string {
	clean := filepath.Clean(path)
	if clean == "." || clean == "" {
		return nil
	}
	parts := []string{}
	for clean != "." && clean != string(filepath.Separator) {
		parent, base := filepath.Split(clean)
		parts = append(parts, base)
		clean = filepath.Clean(parent)
	}
	for left, right := 0, len(parts)-1; left < right; left, right = left+1, right-1 {
		parts[left], parts[right] = parts[right], parts[left]
	}
	return parts
}

func watchV2CreatedDir(watcher *fsnotify.Watcher, root string, event fsnotify.Event) {
	if event.Op&fsnotify.Create == 0 || !isV2WatchedPath(root, event.Name) {
		return
	}
	info, err := os.Stat(event.Name)
	if err != nil || !info.IsDir() {
		return
	}
	_ = addV2DirsRecursive(watcher, event.Name)
}
