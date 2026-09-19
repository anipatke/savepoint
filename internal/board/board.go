package board

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	xterm "github.com/charmbracelet/x/term"
	boardv2 "github.com/opencode/savepoint/internal/board/v2"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/migrate"
)

// Filters is the board's parsed filter surface, carried across both schemas.
// Release and Epic are the V1 filters; Objective is the V2 one. Which of them
// a project accepts is decided in runWithFilters, where the schema is known.
type Filters struct {
	Release   string
	Epic      string
	Objective string
}

func Run() error {
	return RunWithFilters(Filters{})
}

func RunWithFilters(filters Filters) error {
	return runWithFilters(".", filters, os.Stdout, xterm.IsTerminal(os.Stdout.Fd()))
}

// runWithFilters is the board's single dispatch point: it resolves the
// project root once, loads the project through data.LoadProject, and runs the
// board that the project's own schema_version names. The V2 board therefore
// never sees a V1 project and the V1 board never sees a V2 one, so neither
// needs a fallback into the other.
//
// start, stdout, and isTTY are parameters rather than process state so the
// dispatch is exercised against a temporary project directory without
// depending on the working directory (ARCH-03).
func runWithFilters(start string, filters Filters, stdout io.Writer, isTTY bool) error {
	deps := defaultModelDependencies()

	debugf("board dispatch: finding savepoint root from %q", start)
	root, err := deps.Discoverer.FindSavepointRoot(start)
	if err != nil {
		return err
	}
	debugf("board dispatch: root = %q", root)

	// A migration is the highest V2 Next rung even while the source project
	// still declares its old schema. Route it to the V2 board before schema
	// dispatch so board and resume report the same safe action instead of
	// trying to traverse the half-converted V1 tree.
	if pending, err := migrate.PendingOperation(filepath.Dir(root)); err != nil {
		return err
	} else if pending != nil {
		if filters.Release != "" || filters.Epic != "" {
			return fmt.Errorf("V1 board filters are unavailable while migration %s is pending", pending.OperationID)
		}
		return runV2Board(root, filters, stdout, isTTY)
	}

	project, loadErr := data.LoadProject(root)
	version := data.SchemaVersionV1
	if loadErr == nil {
		version = project.SchemaVersion
	} else {
		// A V2 project whose records will not load is a diagnostic screen, not
		// a crash and never a retry through V1 discovery, so route it to the V2
		// board by schema version alone and let its own load report the same
		// failure by name. Every other load failure — including a config.yml
		// whose schema_version itself is unreadable — is reported as it is.
		version = schemaVersionOrV1(root)
		if version != data.SchemaVersionV2 {
			return loadErr
		}
	}

	if err := rejectFiltersForSchema(filters, version); err != nil {
		return err
	}

	if version == data.SchemaVersionV2 {
		return runV2Board(root, filters, stdout, isTTY)
	}
	return runV1Board(root, filters, stdout, isTTY)
}

// schemaVersionOrV1 reports root's declared schema version, reading a version
// that cannot be determined as V1 so the caller reports its own, more specific
// diagnostic instead of this one.
func schemaVersionOrV1(root string) data.SchemaVersion {
	version, err := data.ReadSchemaVersion(filepath.Join(root, "config.yml"))
	if err != nil {
		return data.SchemaVersionV1
	}
	return version
}

// rejectFiltersForSchema refuses a filter flag that the resolved schema has no
// meaning for, naming both the flag and the schema version it was refused for
// rather than ignoring the flag and showing an unfiltered board (CFG-01).
func rejectFiltersForSchema(filters Filters, version data.SchemaVersion) error {
	if version == data.SchemaVersionV2 {
		for _, flag := range []struct{ name, value string }{{"--release", filters.Release}, {"--epic", filters.Epic}} {
			if flag.value != "" {
				return fmt.Errorf("%s is a schema_version 1 filter and this project is schema_version 2; select work with --objective instead", flag.name)
			}
		}
		return nil
	}
	if filters.Objective != "" {
		return fmt.Errorf("--objective is a schema_version 2 filter and this project is schema_version 1; select work with --release and --epic instead")
	}
	return nil
}

// runV2Board hands the resolved root to the V2 board, which owns its own load,
// its own diagnostic screen, and its own non-TTY output.
func runV2Board(root string, filters Filters, stdout io.Writer, isTTY bool) error {
	return boardv2.Run(boardv2.Options{
		Root:            root,
		ObjectiveFilter: filters.Objective,
		Stdout:          stdout,
		TTY:             isTTY,
	})
}

// runV1Board runs today's board over an already-resolved root. The TTY branch
// calls RunTUI with the call shape it already has: RunTUI resolves the root
// itself, which is the same root this dispatch just resolved, and rewiring it
// would edit a second V1 board file for no behavior change. E50 deletes it.
func runV1Board(root string, filters Filters, stdout io.Writer, isTTY bool) error {
	if !isTTY {
		return runPlainOutput(root, filters, stdout)
	}
	return RunTUI(filters.Release, filters.Epic)
}

func runPlainOutput(root string, filters Filters, stdout io.Writer) error {
	model, err := newProjectModelAtRoot(root, filters.Release, filters.Epic, defaultModelDependencies())
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(stdout, RenderPlainTable(model))
	return err
}

func newProjectModel(start, releaseFilter, epicFilter string) (Model, error) {
	return newProjectModelWithDependencies(start, releaseFilter, epicFilter, defaultModelDependencies())
}

func newProjectModelWithDependencies(start, releaseFilter, epicFilter string, deps ModelDependencies) (Model, error) {
	deps = modelDependencies([]ModelDependencies{deps})

	debugf("board init: finding savepoint root from %q", start)
	root, err := deps.Discoverer.FindSavepointRoot(start)
	if err != nil {
		return Model{}, err
	}
	debugf("board init: root = %q", root)

	return newProjectModelAtRoot(root, releaseFilter, epicFilter, deps)
}

// newProjectModelAtRoot builds the V1 board model for an already-resolved
// project root, so the dispatch above resolves it exactly once.
func newProjectModelAtRoot(root, releaseFilter, epicFilter string, deps ModelDependencies) (Model, error) {
	deps = modelDependencies([]ModelDependencies{deps})

	routerState, err := readRouterState(root, deps.RouterReader)
	if err != nil {
		return Model{}, err
	}

	tasks, defects, releaseIDs, releaseEpics, epicStatuses, err := loadBoardData(root, deps.Discoverer, deps.Parser)
	if err != nil {
		return Model{}, err
	}
	debugf("board init: loaded %d tasks across %d releases", len(tasks), len(releaseIDs))

	preferredRelease := routerState.Release
	if releaseFilter != "" {
		preferredRelease = releaseFilter
	}
	preferredEpic := routerState.Epic
	if epicFilter != "" {
		preferredEpic = epicFilter
	}

	release := firstKnown(preferredRelease, releaseIDs)
	if releaseFilter != "" {
		var ok bool
		release, ok = knownValue(releaseFilter, releaseIDs)
		if !ok {
			return Model{}, fmt.Errorf("release %q not found", releaseFilter)
		}
	}

	epic := firstKnown(preferredEpic, releaseEpics[release])
	if epicFilter != "" {
		var ok bool
		epic, ok = knownValue(epicFilter, releaseEpics[release])
		if !ok {
			return Model{}, fmt.Errorf("epic %q not found in release %q", epicFilter, release)
		}
	}

	model := NewModel(tasks, release, epic, deps)
	model.Root = root
	model.RouterTask = routerState.Task
	model.RouterState = routerState
	model.AllDefects = defects
	model.Audit = loadAuditBestEffort(root, deps.AuditLoader)
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
	debugf("board init: file watcher started at %q", root)

	return model, nil
}

func loadBoardData(root string, discoverer taskDiscoverer, parser taskParser) ([]data.Task, []data.Defect, []string, map[string][]string, map[string]string, error) {
	releases, err := discoverer.ListReleases(root)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	releaseIDs := make([]string, 0, len(releases))
	releaseEpics := make(map[string][]string, len(releases))
	var tasks []data.Task
	var defects []data.Defect
	epicStatuses := make(map[string]string)

	for _, release := range releases {
		releaseIDs = append(releaseIDs, release.ID)
		epics, err := discoverer.ListEpics(root, release.ID)
		if err != nil {
			return nil, nil, nil, nil, nil, err
		}
		for _, epic := range epics {
			releaseEpics[release.ID] = append(releaseEpics[release.ID], epic.ID)
			epicTasks, err := loadEpicTasks(discoverer, parser, root, release.ID, epic.ID)
			if err != nil {
				return nil, nil, nil, nil, nil, err
			}
			tasks = append(tasks, epicTasks...)

			detailPath := filepath.Join(epic.Path, shortID(epic.ID)+"-Detail.md")
			if raw, err := os.ReadFile(detailPath); err == nil {
				if fm, err := parser.ParseFrontmatter(string(raw)); err == nil {
					if status, ok := fm["status"].(string); ok {
						epicStatuses[epic.ID] = string(data.NormalizeEpicStatusForLoad(data.EpicStatus(status)))
					}
				}
			}
		}

		releaseDefects, err := loadReleaseDefects(discoverer, parser, root, release.ID)
		if err != nil {
			return nil, nil, nil, nil, nil, err
		}
		defects = append(defects, releaseDefects...)
	}

	return tasks, defects, releaseIDs, releaseEpics, epicStatuses, nil
}

func loadReleaseDefects(discoverer taskDiscoverer, parser taskParser, root, release string) ([]data.Defect, error) {
	infos, err := discoverer.ListDefects(root, release)
	if err != nil {
		return nil, err
	}
	defects := make([]data.Defect, 0, len(infos))
	for _, info := range infos {
		raw, err := os.ReadFile(info.Path)
		if err != nil {
			return nil, err
		}
		d, err := parser.ParseDefectFile(info.Path, string(raw))
		if err != nil {
			return nil, err
		}
		fi, err := os.Stat(info.Path)
		if err != nil {
			return nil, err
		}
		d.Path = info.Path
		d.Mtime = fi.ModTime()
		if d.Release == "" {
			d.Release = release
		}
		defects = append(defects, *d)
	}
	return defects, nil
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
		if os.IsNotExist(err) || strings.Contains(err.Error(), "tasks directory not found") {
			return nil, nil
		}
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
	if value, ok := knownValue(preferred, values); ok {
		return value
	}
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func knownValue(preferred string, values []string) (string, bool) {
	for _, value := range values {
		if value == preferred {
			return value, true
		}
	}
	for _, value := range values {
		if shortID(value) == shortID(preferred) {
			return value, true
		}
	}
	return "", false
}
