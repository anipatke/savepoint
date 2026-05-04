package board

import (
	"fmt"
	"os"
	"path/filepath"

	xterm "github.com/charmbracelet/x/term"
	"github.com/opencode/savepoint/internal/data"
)

func Run() error {
	return RunWithFilters("", "")
}

func RunWithFilters(release, epic string) error {
	if !xterm.IsTerminal(os.Stdout.Fd()) {
		return runPlainOutput(release, epic)
	}
	return RunTUI(release, epic)
}

func runPlainOutput(release, epic string) error {
	model, err := newProjectModel(".", release, epic)
	if err != nil {
		return err
	}
	fmt.Print(RenderPlainTable(model))
	return nil
}

func newProjectModel(start, releaseFilter, epicFilter string) (Model, error) {
	return newProjectModelWithDependencies(start, releaseFilter, epicFilter, defaultModelDependencies())
}

func newProjectModelWithDependencies(start, releaseFilter, epicFilter string, deps ModelDependencies) (Model, error) {
	deps = modelDependencies([]ModelDependencies{deps})

	root, err := deps.Discoverer.FindSavepointRoot(start)
	if err != nil {
		return Model{}, err
	}

	routerState, err := readRouterState(root, deps.RouterReader)
	if err != nil {
		return Model{}, err
	}

	tasks, releaseIDs, releaseEpics, epicStatuses, err := loadBoardData(root, deps.Discoverer, deps.Parser)
	if err != nil {
		return Model{}, err
	}

	preferredRelease := routerState.Release
	if releaseFilter != "" {
		preferredRelease = releaseFilter
	}
	preferredEpic := routerState.Epic
	if epicFilter != "" {
		preferredEpic = epicFilter
	}

	release := firstKnown(preferredRelease, releaseIDs)
	epic := firstKnown(preferredEpic, releaseEpics[release])

	model := NewModel(tasks, release, epic, deps)
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

func loadBoardData(root string, discoverer taskDiscoverer, parser taskParser) ([]data.Task, []string, map[string][]string, map[string]string, error) {
	releases, err := discoverer.ListReleases(root)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	releaseIDs := make([]string, 0, len(releases))
	releaseEpics := make(map[string][]string, len(releases))
	var tasks []data.Task
	epicStatuses := make(map[string]string)

	for _, release := range releases {
		releaseIDs = append(releaseIDs, release.ID)
		epics, err := discoverer.ListEpics(root, release.ID)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		for _, epic := range epics {
			releaseEpics[release.ID] = append(releaseEpics[release.ID], epic.ID)
			epicTasks, err := loadEpicTasks(discoverer, parser, root, release.ID, epic.ID)
			if err != nil {
				return nil, nil, nil, nil, err
			}
			tasks = append(tasks, epicTasks...)

			detailPath := filepath.Join(epic.Path, shortID(epic.ID)+"-Detail.md")
			if raw, err := os.ReadFile(detailPath); err == nil {
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

func readRouterState(root string, reader routerReader) (*data.RouterState, error) {
	content, err := os.ReadFile(filepath.Join(root, "router.md"))
	if err != nil {
		return nil, err
	}

	return reader.ReadState(string(content))
}

func loadEpicTasks(discoverer taskDiscoverer, parser taskParser, root, release, epic string) ([]data.Task, error) {
	taskInfos, err := discoverer.ListTasks(root, release, epic)
	if err != nil {
		return nil, err
	}

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
