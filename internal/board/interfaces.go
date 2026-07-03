package board

import "github.com/opencode/savepoint/internal/data"

// taskDiscoverer provides project traversal for board loading.
type taskDiscoverer interface {
	FindSavepointRoot(start string) (string, error)
	ListReleases(root string) ([]data.ReleaseInfo, error)
	ListEpics(root, release string) ([]data.EpicInfo, error)
	ListTasks(root, release, epic string) ([]data.TaskInfo, error)
	ListDefects(root, release string) ([]data.DefectInfo, error)
}

// taskParser parses Savepoint frontmatter and task files for board loading.
type taskParser interface {
	ParseFrontmatter(content string) (map[string]any, error)
	ParseTaskFile(path string, content string) (*data.Task, error)
	ParseDefectFile(path string, content string) (*data.Defect, error)
}

// configReader reads board display configuration.
type configReader interface {
	Read(path string) (*data.Config, error)
}

// routerReader parses router state from router.md content.
type routerReader interface {
	ReadState(content string) (*data.RouterState, error)
}

// auditLoader loads the repo-wide audit-register data set for board display.
type auditLoader interface {
	LoadAuditRegisterSet(root string) (data.AuditRegisterSet, error)
}

// dataAuditLoader is the production auditLoader backed by the data package's
// tolerant audit-register loader.
type dataAuditLoader struct{}

func (dataAuditLoader) LoadAuditRegisterSet(root string) (data.AuditRegisterSet, error) {
	return data.LoadAuditRegisterSet(root)
}

// ModelDependencies contains board data-access dependencies.
type ModelDependencies struct {
	Discoverer   taskDiscoverer
	Parser       taskParser
	ConfigReader configReader
	RouterReader routerReader
	AuditLoader  auditLoader
}

func defaultModelDependencies() ModelDependencies {
	return ModelDependencies{
		Discoverer:   data.NewDiscover(),
		Parser:       data.NewParser(),
		ConfigReader: data.NewConfigReader(),
		RouterReader: data.NewRouterReader(),
		AuditLoader:  dataAuditLoader{},
	}
}

func modelDependencies(overrides []ModelDependencies) ModelDependencies {
	deps := defaultModelDependencies()
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
	if override.AuditLoader != nil {
		deps.AuditLoader = override.AuditLoader
	}
	return deps
}
