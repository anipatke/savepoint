package data

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Parser is the legacy V1 frontmatter reader. It remains available to the
// explicit migration path and frozen historical fixtures; live V2 consumers
// use the strict record decoders and ParseV2Document instead.
type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
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

	lifecycle := ParseTaskLifecycle(TaskLifecycleMetadata{
		Status: fields.Status,
		Column: fields.Column,
		Stage:  fields.Stage,
		Phase:  fields.Phase,
	})

	task := &Task{
		ID:               fields.ID,
		Title:            firstNonEmpty(fields.Title, fields.Objective),
		Description:      fields.Description,
		Epic:             firstNonEmpty(fields.Epic, extractEpicFromID(fields.ID)),
		Release:          firstNonEmpty(fields.Release, "v1"),
		Column:           lifecycle.Status,
		Stage:            lifecycle.Stage,
		Priority:         fields.Priority,
		Points:           fields.Points,
		Tags:             fields.Tags,
		Acceptance:       firstList(fields.Acceptance, extractChecklistSection(content, "## Acceptance Criteria")),
		Checklist:        extractChecklistItems(content, "## Implementation Plan"),
		Notes:            fields.Notes,
		DependsOn:        fields.DependsOn,
		Progress:         fields.Progress,
		ComplexityTier:   fields.ComplexityTier,
		ComplexityReason: fields.ComplexityReason,
	}

	return task, nil
}

type taskFrontmatter struct {
	ID               string         `yaml:"id"`
	Title            string         `yaml:"title"`
	Objective        string         `yaml:"objective"`
	Description      string         `yaml:"description"`
	Epic             string         `yaml:"epic"`
	Release          string         `yaml:"release"`
	Status           ColumnType     `yaml:"status"`
	Column           ColumnType     `yaml:"column"`
	Phase            ProgressStage  `yaml:"phase"`
	Stage            ProgressStage  `yaml:"stage"`
	Priority         string         `yaml:"priority"`
	Points           int            `yaml:"points"`
	Tags             []string       `yaml:"tags"`
	Acceptance       []string       `yaml:"acceptance"`
	Notes            string         `yaml:"notes"`
	DependsOn        []string       `yaml:"depends_on"`
	Progress         Progress       `yaml:"progress"`
	ComplexityTier   ComplexityTier `yaml:"complexity_tier"`
	ComplexityReason string         `yaml:"complexity_reason"`
}

// normalizeLineEndings replaces Windows (CRLF) and legacy Mac (CR) line endings with LF.
func normalizeLineEndings(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return strings.ReplaceAll(s, "\r", "\n")
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

// V2SourceDocument couples a V2 record's project-relative path with its raw
// frontmatter YAML node and Markdown body. Strict V2 decoders project typed
// fields from Frontmatter via Node.Decode while retaining the document a
// later preserving rewrite needs, instead of going through the healing
// map[string]any path V1 parsing uses.
type V2SourceDocument struct {
	// Path is retained as the project-relative source path when the record
	// comes from discovery. ProjectRoot carries the resolution context for a
	// later managed write, so callers do not need to keep the same cwd.
	Path        string
	ProjectRoot string
	Frontmatter yaml.Node
	Body        string
	// frontmatterText is the exact frontmatter source Frontmatter was parsed
	// from, so a managed write can keep untouched fields byte-for-byte.
	frontmatterText string
	// CRLF records whether the original source used Windows line endings, so
	// a managed rewrite can reproduce the same line-ending form instead of
	// silently normalizing it to LF.
	CRLF bool

	// contentHash is the load-time freshness token. It is deliberately based
	// on the exact source bytes, rather than only mtime, because filesystems
	// may have coarse timestamp resolution and editors may preserve mtimes.
	contentHash    [sha256.Size]byte
	contentHashSet bool
}

// ParseV2Document splits content into its frontmatter YAML node and body,
// giving schema-specific V2 decoders a raw projection boundary distinct from
// the V1 struct-tag decode path in ParseTaskFile.
func ParseV2Document(path, content string) (V2SourceDocument, error) {
	fm, body, err := SplitFrontmatterBody(content)
	if err != nil {
		return V2SourceDocument{}, fmt.Errorf("parse error for %s: %w", path, err)
	}

	var node yaml.Node
	if err := yaml.Unmarshal([]byte(fm), &node); err != nil {
		return V2SourceDocument{}, fmt.Errorf("parse error for %s: failed to parse YAML: %w", path, err)
	}

	return V2SourceDocument{
		Path:            path,
		Frontmatter:     node,
		Body:            body,
		frontmatterText: fm,
		CRLF:            strings.Contains(content, "\r\n"),
		contentHash:     sha256.Sum256([]byte(content)),
		contentHashSet:  true,
	}, nil
}

type defectFrontmatter struct {
	ID         string         `yaml:"id"`
	Release    string         `yaml:"release"`
	Status     DefectStatus   `yaml:"status"`
	Severity   DefectSeverity `yaml:"severity"`
	Introduced string         `yaml:"introduced,omitempty"`
	Reference  string         `yaml:"reference,omitempty"`
	Stage      ProgressStage  `yaml:"stage,omitempty"`
	Title      string         `yaml:"title"`
	Objective  string         `yaml:"objective"`
}

func (p *Parser) ParseDefectFile(path string, content string) (*Defect, error) {
	fm, body, err := SplitFrontmatterBody(normalizeLineEndings(content))
	if err != nil {
		return nil, fmt.Errorf("parse error for %s: %w", path, err)
	}

	var fields defectFrontmatter
	if err := yaml.Unmarshal([]byte(fm), &fields); err != nil {
		return nil, fmt.Errorf("parse error for %s: failed to parse YAML: %w", path, err)
	}

	defect := &Defect{
		ID:         fields.ID,
		Release:    fields.Release,
		Status:     fields.Status,
		Severity:   fields.Severity,
		Introduced: fields.Introduced,
		Reference:  fields.Reference,
		Stage:      fields.Stage,
		Title:      firstNonEmpty(fields.Title, fields.Objective),
		Body:       body,
	}

	NormalizeDefectLifecycleForLoad(defect)

	return defect, nil
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
