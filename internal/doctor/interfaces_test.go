package doctor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/opencode/savepoint/internal/data"
)

type stubDoctorRouterReader struct {
	state *data.RouterState
	calls int
}

func (r *stubDoctorRouterReader) ReadState(content string) (*data.RouterState, error) {
	r.calls++
	return r.state, nil
}

type stubDoctorDiscoverer struct {
	releases []data.ReleaseInfo
	epics    map[string][]data.EpicInfo
	tasks    map[string][]data.TaskInfo
	calls    int
}

func (d *stubDoctorDiscoverer) ListRootDirs(root string) ([]string, error) {
	d.calls++
	return nil, nil
}

func (d *stubDoctorDiscoverer) ListReleases(root string) ([]data.ReleaseInfo, error) {
	d.calls++
	return d.releases, nil
}

func (d *stubDoctorDiscoverer) ListEpics(root, release string) ([]data.EpicInfo, error) {
	d.calls++
	return d.epics[release], nil
}

func (d *stubDoctorDiscoverer) ListTasks(root, release, epic string) ([]data.TaskInfo, error) {
	d.calls++
	return d.tasks[release+"/"+epic], nil
}

type countingDoctorParser struct {
	parser *data.Parser
	calls  int
}

func (p *countingDoctorParser) ParseFrontmatter(content string) (map[string]any, error) {
	p.calls++
	return p.parser.ParseFrontmatter(content)
}

func TestCheckRouterUsesInjectedRouterReader(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "router.md"), "# intentionally not a router state block")
	if err := os.MkdirAll(filepath.Join(root, "releases", "v9", "epics", "E01-mock"), 0755); err != nil {
		t.Fatal(err)
	}

	reader := &stubDoctorRouterReader{state: &data.RouterState{
		State:   "task-building",
		Release: "v9",
		Epic:    "E01-mock",
	}}

	if err := CheckRouter(root, "", DoctorDependencies{RouterReader: reader}); err != nil {
		t.Fatalf("CheckRouter() with injected reader = %v, want nil", err)
	}
	if reader.calls != 1 {
		t.Fatalf("ReadState calls = %d, want 1", reader.calls)
	}
}

func TestCheckDependenciesUsesInjectedDiscovererAndParser(t *testing.T) {
	root := t.TempDir()
	taskPath := filepath.Join(root, "virtual", "T001-task.md")
	if err := os.MkdirAll(filepath.Dir(taskPath), 0755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, taskPath, "---\nid: E01-mock/T001-task\nstatus: planned\nobjective: Mock\ndepends_on: []\n---\n")

	discoverer := &stubDoctorDiscoverer{
		releases: []data.ReleaseInfo{{ID: "v9", Path: filepath.Join(root, "virtual-release")}},
		epics: map[string][]data.EpicInfo{
			"v9": {{ID: "E01-mock", Path: filepath.Join(root, "virtual-epic")}},
		},
		tasks: map[string][]data.TaskInfo{
			"v9/E01-mock": {{ID: "T001-task", Path: taskPath}},
		},
	}
	parser := &countingDoctorParser{parser: data.NewParser()}

	problems := CheckDependencies(root, "", DoctorDependencies{Discoverer: discoverer, Parser: parser})
	if len(problems) > 0 {
		t.Fatalf("CheckDependencies() = %v, want no problems", problems)
	}
	if discoverer.calls == 0 {
		t.Fatal("injected discoverer was not used")
	}
	if parser.calls != 1 {
		t.Fatalf("ParseFrontmatter calls = %d, want 1", parser.calls)
	}
}
