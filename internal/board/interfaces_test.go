package board

import (
	"path/filepath"
	"testing"

	"github.com/opencode/savepoint/internal/data"
)

type stubBoardDiscoverer struct {
	root      string
	releases  []data.ReleaseInfo
	epics     map[string][]data.EpicInfo
	tasks     map[string][]data.TaskInfo
	findCalls int
}

func (d *stubBoardDiscoverer) FindSavepointRoot(start string) (string, error) {
	d.findCalls++
	return d.root, nil
}

func (d *stubBoardDiscoverer) ListReleases(root string) ([]data.ReleaseInfo, error) {
	return d.releases, nil
}

func (d *stubBoardDiscoverer) ListEpics(root, release string) ([]data.EpicInfo, error) {
	return d.epics[release], nil
}

func (d *stubBoardDiscoverer) ListTasks(root, release, epic string) ([]data.TaskInfo, error) {
	return d.tasks[release+"/"+epic], nil
}

type countingBoardParser struct {
	parser           *data.Parser
	frontmatterCalls int
	taskFileCalls    int
}

func (p *countingBoardParser) ParseFrontmatter(content string) (map[string]any, error) {
	p.frontmatterCalls++
	return p.parser.ParseFrontmatter(content)
}

func (p *countingBoardParser) ParseTaskFile(path string, content string) (*data.Task, error) {
	p.taskFileCalls++
	return p.parser.ParseTaskFile(path, content)
}

type stubBoardRouterReader struct {
	state *data.RouterState
	calls int
}

func (r *stubBoardRouterReader) ReadState(content string) (*data.RouterState, error) {
	r.calls++
	return r.state, nil
}

func TestNewProjectModelUsesInjectedInterfaces(t *testing.T) {
	projectRoot := t.TempDir()
	savepointRoot := filepath.Join(projectRoot, ".savepoint")
	epicPath := filepath.Join(savepointRoot, "releases", "v9", "epics", "E01-mock")
	taskPath := filepath.Join(epicPath, "tasks", "T001-mock.md")

	writeFile(t, filepath.Join(savepointRoot, "router.md"), "# router")
	writeFile(t, filepath.Join(epicPath, "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# Epic\n")
	writeFile(t, taskPath, "---\nid: E01-mock/T001-mock\nstatus: planned\nobjective: Mock task\ndepends_on: []\n---\n\n# Task\n")

	discoverer := &stubBoardDiscoverer{
		root: savepointRoot,
		releases: []data.ReleaseInfo{{
			ID:   "v9",
			Path: filepath.Join(savepointRoot, "releases", "v9"),
		}},
		epics: map[string][]data.EpicInfo{
			"v9": {{ID: "E01-mock", Path: epicPath}},
		},
		tasks: map[string][]data.TaskInfo{
			"v9/E01-mock": {{ID: "T001-mock", Path: taskPath}},
		},
	}
	parser := &countingBoardParser{parser: data.NewParser()}
	router := &stubBoardRouterReader{state: &data.RouterState{
		State:   "task-building",
		Release: "v9",
		Epic:    "E01-mock",
		Task:    "E01-mock/T001-mock",
	}}

	model, err := newProjectModelWithDependencies(projectRoot, "", "", ModelDependencies{
		Discoverer:   discoverer,
		Parser:       parser,
		RouterReader: router,
	})
	if err != nil {
		t.Fatalf("newProjectModelWithDependencies() error = %v", err)
	}

	if discoverer.findCalls != 1 {
		t.Fatalf("FindSavepointRoot calls = %d, want 1", discoverer.findCalls)
	}
	if router.calls != 1 {
		t.Fatalf("ReadState calls = %d, want 1", router.calls)
	}
	if parser.frontmatterCalls != 1 || parser.taskFileCalls != 1 {
		t.Fatalf("parser calls = frontmatter:%d task:%d, want 1 each", parser.frontmatterCalls, parser.taskFileCalls)
	}
	if got := model.Tasks[data.ColumnPlanned][0].ID; got != "E01-mock/T001-mock" {
		t.Fatalf("loaded task = %q, want injected task", got)
	}
}
