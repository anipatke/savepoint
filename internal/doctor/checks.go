package doctor

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/opencode/savepoint/internal/data"
	"gopkg.in/yaml.v3"
)

// CheckConfig validates config.yml: exists, valid YAML, required fields present.
func CheckConfig(root string) error {
	configPath := filepath.Join(root, "config.yml")
	raw, err := os.ReadFile(configPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("config.yml not found: %w", data.ErrConfigNotFound)
	}
	if err != nil {
		return fmt.Errorf("config.yml unreadable: %w", err)
	}

	var fields map[string]any
	if err := yaml.Unmarshal(raw, &fields); err != nil {
		return fmt.Errorf("config.yml invalid YAML: %w", err)
	}

	if _, ok := fields["quality_gates"]; !ok {
		return fmt.Errorf("config.yml missing required field: quality_gates")
	}
	if _, ok := fields["theme"]; !ok {
		return fmt.Errorf("config.yml missing required field: theme")
	}

	return nil
}

// CheckRouter validates router.md: valid state name, release/epic directories exist.
// epicFilter, if non-empty, skips directory checks when the router epic doesn't match.
func CheckRouter(root, epicFilter string) error {
	routerPath := filepath.Join(root, "router.md")
	raw, err := os.ReadFile(routerPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("router.md not found: %w", data.ErrConfigNotFound)
	}
	if err != nil {
		return fmt.Errorf("router.md unreadable: %w", err)
	}

	reader := data.NewRouterReader()
	state, err := reader.ReadState(string(raw))
	if err != nil {
		return fmt.Errorf("router.md invalid state block: %w", err)
	}

	if epicFilter != "" && state.Epic != epicFilter {
		return nil
	}

	if state.Release != "" && state.Release != "none" {
		releasePath := filepath.Join(root, "releases", state.Release)
		if _, err := os.Stat(releasePath); os.IsNotExist(err) {
			return fmt.Errorf("router.md release %q directory not found", state.Release)
		}
	}

	if state.Epic != "" && state.Epic != "none" {
		if state.Release == "" || state.Release == "none" {
			return fmt.Errorf("router.md has epic %q but no release", state.Epic)
		}
		epicPath := filepath.Join(root, "releases", state.Release, "epics", state.Epic)
		if _, err := os.Stat(epicPath); os.IsNotExist(err) {
			return fmt.Errorf("router.md epic %q directory not found", state.Epic)
		}
	}

	return nil
}

// Problem describes a single issue found during a structure check.
type Problem struct {
	File    string
	Line    int
	Message string
}

func (p Problem) Error() string {
	if p.Line > 0 {
		return fmt.Sprintf("%s:%d: %s", p.File, p.Line, p.Message)
	}
	if p.File != "" {
		return fmt.Sprintf("%s: %s", p.File, p.Message)
	}
	return p.Message
}

// CheckStructure validates release/epic/task structure and YAML across the project.
// epicFilter, if non-empty, restricts checks to matching epics.
func CheckStructure(root string, epicFilter string) []Problem {
	var problems []Problem
	d := data.NewDiscover()

	releasesPath := filepath.Join(root, "releases")
	if _, err := os.Stat(releasesPath); os.IsNotExist(err) {
		problems = append(problems, Problem{File: releasesPath, Message: "releases directory not found"})
		return problems
	}

	releases, err := d.ListReleases(root)
	if err != nil {
		problems = append(problems, Problem{File: releasesPath, Message: fmt.Sprintf("listing releases: %v", err)})
		return problems
	}

	if len(releases) == 0 {
		problems = append(problems, Problem{File: releasesPath, Message: "no release directories found"})
		return problems
	}

	for _, release := range releases {
		checkReleasePRD(release.Path, release.ID, &problems)

		epics, err := d.ListEpics(root, release.ID)
		if err != nil {
			problems = append(problems, Problem{
				File:    filepath.Join(release.Path, "epics"),
				Message: fmt.Sprintf("listing epics in release %q: %v", release.ID, err),
			})
			continue
		}

		for _, epic := range epics {
			if epicFilter != "" && epic.ID != epicFilter && !strings.HasPrefix(epic.ID, epicFilter) {
				continue
			}

			checkEpicDetail(epic.Path, epic.ID, &problems)

			tasks, err := d.ListTasks(root, release.ID, epic.ID)
			if err != nil {
				problems = append(problems, Problem{
					File:    filepath.Join(epic.Path, "tasks"),
					Message: fmt.Sprintf("listing tasks in epic %q: %v", epic.ID, err),
				})
				continue
			}

			for _, task := range tasks {
				checkTaskFile(task.Path, &problems)
			}
		}
	}

	return problems
}

func checkReleasePRD(releasePath string, releaseID string, problems *[]Problem) {
	prdPath := filepath.Join(releasePath, releaseID+"-PRD.md")
	raw, err := os.ReadFile(prdPath)
	if os.IsNotExist(err) {
		*problems = append(*problems, Problem{File: prdPath, Message: "release PRD file not found"})
		return
	}
	if err != nil {
		*problems = append(*problems, Problem{File: prdPath, Message: fmt.Sprintf("unreadable: %v", err)})
		return
	}
	validateFrontmatter(prdPath, string(raw), problems)
}

func checkEpicDetail(epicPath string, epicID string, problems *[]Problem) {
	prefix := extractPrefix(epicID)
	detailPath := filepath.Join(epicPath, prefix+"-Detail.md")
	raw, err := os.ReadFile(detailPath)
	if os.IsNotExist(err) {
		*problems = append(*problems, Problem{File: detailPath, Message: "epic detail file not found"})
		return
	}
	if err != nil {
		*problems = append(*problems, Problem{File: detailPath, Message: fmt.Sprintf("unreadable: %v", err)})
		return
	}
	validateFrontmatter(detailPath, string(raw), problems)
}

func extractPrefix(epicID string) string {
	if idx := strings.IndexByte(epicID, '-'); idx != -1 {
		return epicID[:idx]
	}
	return epicID
}

func checkTaskFile(path string, problems *[]Problem) {
	raw, err := os.ReadFile(path)
	if err != nil {
		*problems = append(*problems, Problem{File: path, Message: fmt.Sprintf("unreadable: %v", err)})
		return
	}

	content := string(raw)
	parser := data.NewParser()
	fm, err := parser.ParseFrontmatter(content)
	if err != nil {
		line := extractYAMLLine(err)
		*problems = append(*problems, Problem{File: path, Line: line, Message: fmt.Sprintf("invalid frontmatter: %v", err)})
		return
	}

	checkRequiredString(fm, path, "id", problems)
	checkRequiredString(fm, path, "status", problems)
	checkRequiredString(fm, path, "objective", problems)
	checkDependsOn(fm, path, problems)

	if !hasAcceptanceCriteria(content) {
		*problems = append(*problems, Problem{File: path, Message: "task missing ## Acceptance Criteria section"})
	}
}

func checkRequiredString(fm map[string]any, path, field string, problems *[]Problem) {
	val, ok := fm[field]
	if !ok {
		*problems = append(*problems, Problem{File: path, Message: fmt.Sprintf("task missing required frontmatter field: %s", field)})
		return
	}
	s, ok := val.(string)
	if !ok || s == "" {
		*problems = append(*problems, Problem{File: path, Message: fmt.Sprintf("task frontmatter field %q must be a non-empty string", field)})
	}
}

func checkDependsOn(fm map[string]any, path string, problems *[]Problem) {
	val, ok := fm["depends_on"]
	if !ok {
		return
	}
	switch val.(type) {
	case []any, []string:
	default:
		*problems = append(*problems, Problem{File: path, Message: "task frontmatter field depends_on must be a list"})
	}
}

func hasAcceptanceCriteria(content string) bool {
	normalized := strings.ReplaceAll(content, "\r\n", "\n")
	idx := strings.Index(normalized, "## Acceptance Criteria")
	if idx == -1 {
		return false
	}
	section := normalized[idx+len("## Acceptance Criteria"):]
	if next := strings.Index(section, "\n## "); next != -1 {
		section = section[:next]
	}
	section = strings.TrimSpace(section)
	return section != ""
}

func validateFrontmatter(path, content string, problems *[]Problem) {
	parser := data.NewParser()
	_, err := parser.ParseFrontmatter(content)
	if err != nil {
		line := extractYAMLLine(err)
		*problems = append(*problems, Problem{File: path, Line: line, Message: fmt.Sprintf("invalid frontmatter: %v", err)})
	}
}

// taskDep describes a parsed task's dependency information.
type taskDep struct {
	File      string
	ID        string
	DependsOn []string
}

// CheckDependencies validates task dependency integrity:
// missing deps, duplicate IDs, and dependency cycles.
// epicFilter restricts checks to matching epics if non-empty.
func CheckDependencies(root string, epicFilter string) []Problem {
	var problems []Problem
	d := data.NewDiscover()

	releases, err := d.ListReleases(root)
	if err != nil {
		problems = append(problems, Problem{Message: fmt.Sprintf("listing releases: %v", err)})
		return problems
	}

	var allTasks []taskDep
	idSet := make(map[string]string) // id -> first file seen

	for _, release := range releases {
		epics, err := d.ListEpics(root, release.ID)
		if err != nil {
			continue
		}
		for _, epic := range epics {
			if epicFilter != "" && epic.ID != epicFilter && !strings.HasPrefix(epic.ID, epicFilter) {
				continue
			}
			tasks, err := d.ListTasks(root, release.ID, epic.ID)
			if err != nil {
				continue
			}
			for _, t := range tasks {
				td := parseTaskDep(t.Path)
				if td == nil {
					continue
				}
				allTasks = append(allTasks, *td)
				if existing, ok := idSet[td.ID]; ok {
					problems = append(problems, Problem{
						File:    td.File,
						Message: fmt.Sprintf("duplicate task ID %q (first seen in %s)", td.ID, existing),
					})
				} else {
					idSet[td.ID] = td.File
				}
			}
		}
	}

	// Check for missing dependencies and cycles
	graph := make(map[string][]string) // id -> list of dependencies
	idToFile := make(map[string]string)

	for _, td := range allTasks {
		idToFile[td.ID] = td.File
		graph[td.ID] = td.DependsOn
	}

	for _, td := range allTasks {
		for _, dep := range td.DependsOn {
			if _, exists := idSet[dep]; !exists {
				problems = append(problems, Problem{
					File:    td.File,
					Message: fmt.Sprintf("depends_on references non-existent task %q", dep),
				})
			}
		}
	}

	// Cycle detection using DFS
	cycleProblems := detectCycles(graph, idToFile)
	problems = append(problems, cycleProblems...)

	return problems
}

func parseTaskDep(path string) *taskDep {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	parser := data.NewParser()
	fm, err := parser.ParseFrontmatter(string(raw))
	if err != nil {
		return nil
	}
	id, _ := fm["id"].(string)
	if id == "" {
		return nil
	}
	var deps []string
	switch v := fm["depends_on"].(type) {
	case []any:
		for _, d := range v {
			if s, ok := d.(string); ok {
				deps = append(deps, s)
			}
		}
	case []string:
		deps = v
	}
	return &taskDep{
		File:      path,
		ID:        id,
		DependsOn: deps,
	}
}

// detectCycles runs DFS on the dependency graph and returns cycle problems.
// Uses a path stack to accurately reconstruct cycle paths (avoids parent-map
// overwrite issues that produced inaccurate paths).
func detectCycles(graph map[string][]string, idToFile map[string]string) []Problem {
	const (
		white = 0
		gray  = 1
		black = 2
	)
	color := make(map[string]int)
	path := make([]string, 0)

	for id := range graph {
		color[id] = white
	}

	var problems []Problem

	var dfs func(id string)
	dfs = func(id string) {
		color[id] = gray
		path = append(path, id)
		for _, dep := range graph[id] {
			switch color[dep] {
			case white:
				dfs(dep)
			case gray:
				cycleStart := -1
				for i, n := range path {
					if n == dep {
						cycleStart = i
						break
					}
				}
				if cycleStart >= 0 {
					cycle := path[cycleStart:]
					cyclePath := make([]string, 0, len(cycle))
					for _, cid := range cycle {
						if f, ok := idToFile[cid]; ok {
							cyclePath = append(cyclePath, f)
						} else {
							cyclePath = append(cyclePath, cid)
						}
					}
					problems = append(problems, Problem{
						Message: fmt.Sprintf("dependency cycle detected: %s", strings.Join(cyclePath, " → ")),
					})
				}
			}
		}
		path = path[:len(path)-1]
		color[id] = black
	}

	for id := range graph {
		if color[id] == white {
			dfs(id)
		}
	}
	return problems
}

// CheckAuditState finds audit proposal files without matching audit-pending state in the router.
func CheckAuditState(root string) []Problem {
	var problems []Problem
	d := data.NewDiscover()

	routerPath := filepath.Join(root, "router.md")
	raw, err := os.ReadFile(routerPath)
	if err != nil {
		return problems
	}

	reader := data.NewRouterReader()
	state, err := reader.ReadState(string(raw))
	if err != nil {
		return problems
	}

	releases, err := d.ListReleases(root)
	if err != nil {
		return problems
	}

	for _, release := range releases {
		epics, err := d.ListEpics(root, release.ID)
		if err != nil {
			continue
		}
		for _, epic := range epics {
			prefix := extractPrefix(epic.ID)
			auditPath := filepath.Join(epic.Path, prefix+"-Audit.md")
			if _, err := os.Stat(auditPath); os.IsNotExist(err) {
				continue
			}
			if state.State != "audit-pending" || state.Epic != epic.ID {
				problems = append(problems, Problem{
					File:    auditPath,
					Message: fmt.Sprintf("audit proposal exists but router state is %q (epic: %q) — expected audit-pending for %q", state.State, state.Epic, epic.ID),
				})
			}
		}
	}

	return problems
}

// CheckOrphans finds tasks whose epic prefix in their ID does not match any existing epic directory.
func CheckOrphans(root string) []Problem {
	var problems []Problem

	existingEpics := make(map[string]bool)
	releasesPath := filepath.Join(root, "releases")
	rd, err := os.ReadDir(releasesPath)
	if err != nil {
		problems = append(problems, Problem{File: releasesPath, Message: fmt.Sprintf("listing releases: %v", err)})
		return problems
	}

	for _, releaseDir := range rd {
		if !releaseDir.IsDir() {
			continue
		}
		epicsPath := filepath.Join(releasesPath, releaseDir.Name(), "epics")
		epicEntries, err := os.ReadDir(epicsPath)
		if err != nil {
			continue
		}
		for _, e := range epicEntries {
			if e.IsDir() {
				existingEpics[e.Name()] = true
			}
		}
	}

	// Collect all tasks and check their epic references
	d := data.NewDiscover()
	allReleases, err := d.ListReleases(root)
	if err != nil {
		return problems
	}

	for _, release := range allReleases {
		epics, err := d.ListEpics(root, release.ID)
		if err != nil {
			continue
		}
		for _, epic := range epics {
			tasks, err := d.ListTasks(root, release.ID, epic.ID)
			if err != nil {
				continue
			}
			for _, t := range tasks {
				raw, err := os.ReadFile(t.Path)
				if err != nil {
					continue
				}
				parser := data.NewParser()
				fm, err := parser.ParseFrontmatter(string(raw))
				if err != nil {
					continue
				}
				id, _ := fm["id"].(string)
				if id == "" {
					continue
				}
				idx := strings.IndexByte(id, '/')
				if idx == -1 {
					continue
				}
				taskEpic := id[:idx]
				if !existingEpics[taskEpic] {
					problems = append(problems, Problem{
						File:    t.Path,
						Message: fmt.Sprintf("orphaned task: epic %q does not exist in any release — consider moving to .savepoint/orphans/", taskEpic),
					})
				}
			}
		}
	}

	return problems
}

func extractYAMLLine(err error) int {
	s := err.Error()
	const prefix = "yaml: line "
	if idx := strings.Index(s, prefix); idx != -1 {
		rest := s[idx+len(prefix):]
		if end := strings.IndexByte(rest, ':'); end != -1 {
			if line, err := strconv.Atoi(rest[:end]); err == nil {
				return line
			}
		}
	}
	return 0
}
