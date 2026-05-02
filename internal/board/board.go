package board

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/opencode/savepoint/internal/data"
)

func Run() error {
	lipgloss.SetColorProfile(termenv.ANSI256)

	model, err := newProjectModel(".")
	if err != nil {
		return err
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}
	return nil
}

func newProgramModel() Model {
	return NewModel(nil, "v1", "E03-board-tui-core")
}

func newProjectModel(start string) (Model, error) {
	d := data.NewDiscover()
	root, err := d.FindSavepointRoot(start)
	if err != nil {
		return Model{}, err
	}

	routerState, err := readRouterState(root)
	if err != nil {
		return Model{}, err
	}

	tasks, releaseIDs, releaseEpics, epicStatuses, err := loadBoardData(root)
	if err != nil {
		return Model{}, err
	}

	release := firstKnown(routerState.Release, releaseIDs)
	epic := firstKnown(routerState.Epic, releaseEpics[release])

	model := NewModel(tasks, release, epic)
	model.Root = root
	model.RouterTask = routerState.Task
	model.RouterState = routerState
	model.Releases = releaseIDs
	model.ReleaseEpics = releaseEpics
	model.EpicStatus = epicStatuses
	model.refreshEpicsForRelease()
	model.refreshTasks()

	watcher, err := newWatcher(root)
	if err != nil {
		return Model{}, err
	}
	model.Watcher = watcher

	return model, nil
}

func loadBoardData(root string) ([]data.Task, []string, map[string][]string, map[string]string, error) {
	d := data.NewDiscover()
	releases, err := d.ListReleases(root)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	releaseIDs := make([]string, 0, len(releases))
	releaseEpics := make(map[string][]string, len(releases))
	var tasks []data.Task
	epicStatuses := make(map[string]string)

	for _, release := range releases {
		releaseIDs = append(releaseIDs, release.ID)
		epics, err := d.ListEpics(root, release.ID)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		for _, epic := range epics {
			releaseEpics[release.ID] = append(releaseEpics[release.ID], epic.ID)
			epicTasks, err := loadEpicTasks(d, root, release.ID, epic.ID)
			if err != nil {
				return nil, nil, nil, nil, err
			}
			tasks = append(tasks, epicTasks...)

			detailPath := filepath.Join(epic.Path, shortID(epic.ID)+"-Detail.md")
			if raw, err := os.ReadFile(detailPath); err == nil {
				parser := data.NewParser()
				if fm, err := parser.ParseFrontmatter(string(raw)); err == nil {
					if status, ok := fm["status"].(string); ok {
						epicStatuses[epic.ID] = status
					}
				}
			}
		}
	}

	return tasks, releaseIDs, releaseEpics, epicStatuses, nil
}

func loadAllTasks(root string) ([]data.Task, error) {
	tasks, _, _, _, err := loadBoardData(root)
	return tasks, err
}

func readRouterState(root string) (*data.RouterState, error) {
	content, err := os.ReadFile(filepath.Join(root, "router.md"))
	if err != nil {
		return nil, err
	}

	return data.NewRouterReader().ReadState(string(content))
}

func loadEpicTasks(d *data.Discover, root, release, epic string) ([]data.Task, error) {
	taskInfos, err := d.ListTasks(root, release, epic)
	if err != nil {
		return nil, err
	}

	parser := data.NewParser()
	tasks := make([]data.Task, 0, len(taskInfos))
	for _, taskInfo := range taskInfos {
		content, err := os.ReadFile(taskInfo.Path)
		if err != nil {
			return nil, err
		}
		task, err := parser.ParseTaskFile(taskInfo.Path, string(content))
		if err != nil {
			return nil, err
		}
		fi, err := os.Stat(taskInfo.Path)
		if err != nil {
			return nil, err
		}
		task.Path = taskInfo.Path
		task.Mtime = fi.ModTime()
		task.Release = release
		task.Epic = epic
		tasks = append(tasks, *task)
	}
	return tasks, nil
}

func firstKnown(preferred string, values []string) string {
	for _, value := range values {
		if value == preferred {
			return preferred
		}
	}
	for _, value := range values {
		if shortID(value) == shortID(preferred) {
			return value
		}
	}
	if len(values) == 0 {
		return ""
	}
	return values[0]
}
