package doctor

import "github.com/opencode/savepoint/internal/data"

// taskDiscoverer provides project traversal for doctor checks.
type taskDiscoverer interface {
	ListRootDirs(root string) ([]string, error)
	ListReleases(root string) ([]data.ReleaseInfo, error)
	ListEpics(root, release string) ([]data.EpicInfo, error)
	ListTasks(root, release, epic string) ([]data.TaskInfo, error)
}

// taskParser parses Savepoint frontmatter for doctor checks.
type taskParser interface {
	ParseFrontmatter(content string) (map[string]any, error)
}

// configReader reads quality gate configuration.
type configReader interface {
	Read(path string) (*data.Config, error)
}

// routerReader parses router state from router.md content.
type routerReader interface {
	ReadState(content string) (*data.RouterState, error)
}

// DoctorDependencies contains doctor data-access dependencies.
type DoctorDependencies struct {
	Discoverer   taskDiscoverer
	Parser       taskParser
	ConfigReader configReader
	RouterReader routerReader
}

func defaultDoctorDependencies() DoctorDependencies {
	return DoctorDependencies{
		Discoverer:   data.NewDiscover(),
		Parser:       data.NewParser(),
		ConfigReader: data.NewConfigReader(),
		RouterReader: data.NewRouterReader(),
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
	return deps
}
