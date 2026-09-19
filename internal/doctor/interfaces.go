package doctor

import "github.com/opencode/savepoint/internal/data"

// taskDiscoverer provides project traversal for doctor checks.
type taskDiscoverer interface {
	ListRootDirs(root string) ([]string, error)
	ListReleases(root string) ([]data.ReleaseInfo, error)
	ListEpics(root, release string) ([]data.EpicInfo, error)
	ListTasks(root, release, epic string) ([]data.TaskInfo, error)
	ListDefects(root, release string) ([]data.DefectInfo, error)
}

// taskParser parses Savepoint frontmatter for doctor checks.
type taskParser interface {
	ParseFrontmatter(content string) (map[string]any, error)
	ParseDefectFile(path, content string) (*data.Defect, error)
	ParseRawFindingFile(path, content string) (*data.AuditFinding, error)
}

// configReader reads quality gate configuration.
type configReader interface {
	Read(path string) (*data.Config, error)
}

// routerReader parses router state from router.md content.
type routerReader interface {
	ReadState(content string) (*data.RouterState, error)
}

// projectLoader is the doctor's read-only consumer boundary for the
// schema-dispatched V2 index. It keeps doctor from re-parsing Release records
// or reimplementing the project's structural validation rules.
type projectLoader interface {
	Load(root string) (*data.Project, error)
}

type defaultProjectLoader struct{}

func (defaultProjectLoader) Load(root string) (*data.Project, error) {
	return data.LoadProject(root)
}

// DoctorDependencies contains doctor data-access dependencies.
type DoctorDependencies struct {
	Discoverer    taskDiscoverer
	Parser        taskParser
	ConfigReader  configReader
	RouterReader  routerReader
	ProjectLoader projectLoader
}

func defaultDoctorDependencies() DoctorDependencies {
	return DoctorDependencies{
		Discoverer:    data.NewDiscover(),
		Parser:        data.NewParser(),
		ConfigReader:  data.NewConfigReader(),
		RouterReader:  data.NewRouterReader(),
		ProjectLoader: defaultProjectLoader{},
	}
}

func doctorDependencies(overrides []DoctorDependencies) DoctorDependencies {
	deps := defaultDoctorDependencies()
	if len(overrides) == 0 {
		return deps
	}
	override := overrides[0]
	if override.Discoverer != nil {
		deps.Discoverer = override.Discoverer
	}
	if override.Parser != nil {
		deps.Parser = override.Parser
	}
	if override.ConfigReader != nil {
		deps.ConfigReader = override.ConfigReader
	}
	if override.RouterReader != nil {
		deps.RouterReader = override.RouterReader
	}
	if override.ProjectLoader != nil {
		deps.ProjectLoader = override.ProjectLoader
	}
	return deps
}
