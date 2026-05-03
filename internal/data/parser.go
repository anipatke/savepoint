package data

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseFrontmatter(content string) (map[string]any, error) {
	frontmatter, err := extractFrontmatter(content)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := yaml.Unmarshal([]byte(frontmatter), &result); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return result, nil
}

func (p *Parser) ParseTaskFile(path string, content string) (*Task, error) {
	frontmatter, err := extractFrontmatter(content)
	if err != nil {
		return nil, fmt.Errorf("parse error for %s: %w", path, err)
	}

	var fields taskFrontmatter
	if err := yaml.Unmarshal([]byte(frontmatter), &fields); err != nil {
		return nil, fmt.Errorf("parse error for %s: failed to parse YAML: %w", path, err)
	}

	rawColumn := firstColumn(fields.Column, fields.Status)
	task := &Task{
		ID:          fields.ID,
		Title:       firstNonEmpty(fields.Title, fields.Objective),
		Description: fields.Description,
		Epic:        firstNonEmpty(fields.Epic, extractEpicFromID(fields.ID)),
		Release:     firstNonEmpty(fields.Release, "v1"),
		Column:      normalizeColumn(rawColumn),
		Stage:       firstStage(fields.Phase, fields.Stage),
		Priority:    fields.Priority,
		Points:      fields.Points,
		Tags:        fields.Tags,
		Acceptance:  firstList(fields.Acceptance, extractChecklistSection(content, "## Acceptance Criteria")),
		Checklist:   extractChecklistItems(content, "## Implementation Plan"),
		Notes:       fields.Notes,
		DependsOn:   fields.DependsOn,
		Progress:    fields.Progress,
	}

	if err := validateParsedTaskLifecycle(rawColumn, *task); err != nil {
		return nil, fmt.Errorf("parse error for %s: %w", path, err)
	}

	if task.Column == ColumnInProgress && task.Stage == "" {
		task.Stage = StageBuild
	}

	return task, nil
}

type taskFrontmatter struct {
	ID          string        `yaml:"id"`
	Title       string        `yaml:"title"`
	Objective   string        `yaml:"objective"`
	Description string        `yaml:"description"`
	Epic        string        `yaml:"epic"`
	Release     string        `yaml:"release"`
	Status      ColumnType    `yaml:"status"`
	Column      ColumnType    `yaml:"column"`
	Phase       ProgressStage `yaml:"phase"`
	Stage       ProgressStage `yaml:"stage"`
	Priority    string        `yaml:"priority"`
	Points      int           `yaml:"points"`
	Tags        []string      `yaml:"tags"`
	Acceptance  []string      `yaml:"acceptance"`
	Notes       string        `yaml:"notes"`
	DependsOn   []string      `yaml:"depends_on"`
	Progress    Progress      `yaml:"progress"`
}

// normalizeLineEndings replaces Windows line endings with Unix line endings.
func normalizeLineEndings(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}

func extractFrontmatter(content string) (string, error) {
	normalized := normalizeLineEndings(content)
	if !strings.HasPrefix(normalized, "---\n") {
		return "", ErrNoFrontmatter
	}

	end := strings.Index(normalized[len("---\n"):], "\n---")
	if end == -1 {
		return "", ErrNoClosingFrontmatter
	}

	return strings.TrimSpace(normalized[len("---\n") : len("---\n")+end]), nil
}

func extractEpicFromID(id string) string {
	parts := strings.Split(id, "/")
	if len(parts) >= 1 {
		return parts[0]
	}
	return ""
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func firstColumn(values ...ColumnType) ColumnType {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func normalizeColumn(value ColumnType) ColumnType {
	switch value {
	case "", legacyTodoColumn:
		return ColumnPlanned
	case ColumnPlanned, ColumnInProgress, ColumnDone:
		return value
	default:
		return value
	}
}

const legacyTodoColumn ColumnType = "todo"

func validateParsedTaskLifecycle(rawColumn ColumnType, task Task) error {
	if rawColumn != "" && rawColumn != legacyTodoColumn && !IsCanonicalColumn(rawColumn) {
		return fmt.Errorf("invalid task status %q: use planned, in_progress, or done. Add 'status: planned' or 'status: in_progress' to task frontmatter", rawColumn)
	}
	if task.Column == ColumnInProgress && !IsCanonicalStage(task.Stage) && task.Stage != "" {
		return fmt.Errorf("invalid phase %q: use build, test, or audit. Add 'phase: build' to task frontmatter", task.Stage)
	}
	if task.Column != ColumnInProgress && task.Stage != "" {
		return nil
	}
	return nil
}

func firstStage(values ...ProgressStage) ProgressStage {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func firstList(values ...[]string) []string {
	for _, value := range values {
		if len(value) > 0 {
			return value
		}
	}
	return nil
}

func extractChecklistItems(content, heading string) []CheckItem {
	normalized := normalizeLineEndings(content)
	start := strings.Index(normalized, heading)
	if start == -1 {
		return nil
	}

	section := normalized[start+len(heading):]
	if next := strings.Index(section, "\n## "); next != -1 {
		section = section[:next]
	}

	items := []CheckItem{}
	var current *CheckItem
	for _, line := range strings.Split(section, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "- [x] ") {
			items = append(items, CheckItem{Text: strings.TrimSpace(trimmed[6:]), Done: true})
			current = &items[len(items)-1]
			continue
		}
		if strings.HasPrefix(trimmed, "- [ ] ") {
			items = append(items, CheckItem{Text: strings.TrimSpace(trimmed[6:]), Done: false})
			current = &items[len(items)-1]
			continue
		}
		if strings.HasPrefix(trimmed, "- ") {
			items = append(items, CheckItem{Text: strings.TrimSpace(trimmed[2:]), Done: false})
			current = &items[len(items)-1]
			continue
		}
		if trimmed != "" && current != nil {
			current.Text = strings.TrimSpace(current.Text + " " + trimmed)
		}
	}
	return items
}

func extractChecklistSection(content, heading string) []string {
	normalized := normalizeLineEndings(content)
	start := strings.Index(normalized, heading)
	if start == -1 {
		return nil
	}

	section := normalized[start+len(heading):]
	if next := strings.Index(section, "\n## "); next != -1 {
		section = section[:next]
	}

	items := []string{}
	for _, line := range strings.Split(section, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "- [ ] ") || strings.HasPrefix(trimmed, "- [x] ") {
			items = append(items, strings.TrimSpace(trimmed[6:]))
			continue
		}
		if strings.HasPrefix(trimmed, "- ") {
			items = append(items, strings.TrimSpace(trimmed[2:]))
		}
	}
	return items
}
