package board

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/opencode/savepoint/internal/data"
)

func Run() error {
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

	releases, err := d.ListReleases(root)
	if err != nil {
		return Model{}, err
	}

	releaseIDs := make([]string, 0, len(releases))
	releaseEpics := make(map[string][]string, len(releases))
	tasks := []data.Task{}

	for _, release := range releases {
		releaseIDs = append(releaseIDs, release.ID)
		epics, err := d.ListEpics(root, release.ID)
		if err != nil {
			return Model{}, err
		}
		for _, epic := range epics {
			releaseEpics[release.ID] = append(releaseEpics[release.ID], epic.ID)
			epicTasks, err := loadEpicTasks(d, root, release.ID, epic.ID)
			if err != nil {
				return Model{}, err
			}
			tasks = append(tasks, epicTasks...)
		}
	}

	release := firstKnown(routerState.Release, releaseIDs)
	epic := firstKnown(routerState.Epic, releaseEpics[release])

	model := NewModel(tasks, release, epic)
	model.Root = root
	model.Releases = releaseIDs
	model.ReleaseEpics = releaseEpics
	model.refreshEpicsForRelease()
	model.refreshTasks()

	return model, nil
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
	if len(values) == 0 {
		return ""
	}
	return values[0]
}
