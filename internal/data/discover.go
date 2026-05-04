package data

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type ReleaseInfo struct {
	ID    string
	Path  string
	Epics []EpicInfo
}

type EpicInfo struct {
	ID    string
	Path  string
	Tasks []TaskInfo
}

type TaskInfo struct {
	ID   string
	Path string
}

type Discover struct{}

func NewDiscover() *Discover {
	return &Discover{}
}

func (d *Discover) FindSavepointRoot(start string) (string, error) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}

	for {
		savepointPath := filepath.Join(dir, ".savepoint")
		info, err := os.Stat(savepointPath)
		if err == nil && info.IsDir() {
			return savepointPath, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", ErrSavepointDirectoryMissing
		}
		dir = parent
	}
}

func (d *Discover) ListReleases(root string) ([]ReleaseInfo, error) {
	releasesPath := filepath.Join(root, "releases")
	info, err := os.Stat(releasesPath)
	if err != nil {
		return nil, fmt.Errorf("releases directory not found: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("releases is not a directory")
	}

	entries, err := os.ReadDir(releasesPath)
	if err != nil {
		return nil, err
	}

	var releases []ReleaseInfo
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		id := entry.Name()
		releases = append(releases, ReleaseInfo{
			ID:   id,
			Path: filepath.Join(releasesPath, id),
		})
	}

	sort.Slice(releases, func(i, j int) bool {
		return releases[i].ID < releases[j].ID
	})
	return releases, nil
}

// ListRootDirs returns sorted child directory names directly under root.
func (d *Discover) ListRootDirs(root string) ([]string, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", root)
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}

	sort.Strings(dirs)
	return dirs, nil
}

func (d *Discover) ListEpics(root, release string) ([]EpicInfo, error) {
	epicsPath := filepath.Join(root, "releases", release, "epics")
	info, err := os.Stat(epicsPath)
	if err != nil {
		return nil, fmt.Errorf("epics directory not found: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("epics is not a directory")
	}

	entries, err := os.ReadDir(epicsPath)
	if err != nil {
		return nil, err
	}

	var epics []EpicInfo
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), "_") {
			continue
		}
		id := entry.Name()
		epics = append(epics, EpicInfo{
			ID:   id,
			Path: filepath.Join(epicsPath, id),
		})
	}

	sort.Slice(epics, func(i, j int) bool {
		return epics[i].ID < epics[j].ID
	})
	return epics, nil
}

func (d *Discover) ListTasks(root, release, epic string) ([]TaskInfo, error) {
	tasksPath := filepath.Join(root, "releases", release, "epics", epic, "tasks")
	info, err := os.Stat(tasksPath)
	if err != nil {
		return nil, fmt.Errorf("tasks directory not found: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("tasks is not a directory")
	}

	entries, err := os.ReadDir(tasksPath)
	if err != nil {
		return nil, err
	}

	var tasks []TaskInfo
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}
		id := entry.Name()[:len(entry.Name())-3]
		tasks = append(tasks, TaskInfo{
			ID:   id,
			Path: filepath.Join(tasksPath, entry.Name()),
		})
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].ID < tasks[j].ID
	})
	return tasks, nil
}
