package data

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

var ErrMtimeConflict = fmt.Errorf("file modified since last read")
var ErrProposalNotFound = fmt.Errorf("target text not found in file")

// ApplyProposal replaces the first occurrence of old with newText in the file at path.
func ApplyProposal(path, old, newText string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	normalized := normalizeLineEndings(string(content))
	if !strings.Contains(normalized, old) {
		return fmt.Errorf("%w: %s", ErrProposalNotFound, path)
	}
	updated := strings.Replace(normalized, old, newText, 1)
	return os.WriteFile(path, []byte(updated), 0644)
}

// UpdateEpicStatus sets the status field in the frontmatter of an E##-Detail.md file.
func UpdateEpicStatus(path, status string) error {
	return updateFrontmatterField(path, "status", status)
}

// UpdateLastAudited sets the last_audited field in the frontmatter of Design.md.
func UpdateLastAudited(path, value string) error {
	return updateFrontmatterField(path, "last_audited", value)
}

// SplitFrontmatterBody splits content into frontmatter YAML and body.
func SplitFrontmatterBody(content string) (yamlStr string, body string, err error) {
	normalized := normalizeLineEndings(content)
	raw, err := extractFrontmatter(normalized)
	if err != nil {
		return "", "", err
	}
	delimLen := 4
	bodyStart := delimLen + len(raw) + delimLen
	body = ""
	if bodyStart < len(normalized) {
		body = normalized[bodyStart:]
	}
	return raw, body, nil
}

func updateFrontmatterField(path, key, value string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}

	normalized := normalizeLineEndings(string(content))

	raw, body, err := SplitFrontmatterBody(normalized)
	if err != nil {
		return fmt.Errorf("extract frontmatter: %w", err)
	}

	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(raw), &doc); err != nil {
		return fmt.Errorf("parse yaml: %w", err)
	}

	if doc.Kind != yaml.DocumentNode || len(doc.Content) == 0 {
		return fmt.Errorf("unexpected yaml structure")
	}

	mapping := doc.Content[0]
	if mapping.Kind != yaml.MappingNode {
		return fmt.Errorf("frontmatter is not a mapping")
	}

	setMappingField(mapping, key, value)

	out, err := yaml.Marshal(&doc)
	if err != nil {
		return fmt.Errorf("marshal yaml: %w", err)
	}

	newContent := "---\n" + strings.TrimSpace(string(out)) + "\n---" + body
	return os.WriteFile(path, []byte(newContent), 0644)
}

func WriteTaskStatus(path string, task *Task, expectedMtime time.Time) error {
	if err := ValidateTaskLifecycle(task); err != nil {
		return err
	}

	fi, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}

	if !fi.ModTime().Equal(expectedMtime) {
		return ErrMtimeConflict
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}

	normalized := normalizeLineEndings(string(content))

	raw, body, err := SplitFrontmatterBody(normalized)
	if err != nil {
		return fmt.Errorf("extract frontmatter: %w", err)
	}

	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(raw), &doc); err != nil {
		return fmt.Errorf("parse yaml: %w", err)
	}

	if doc.Kind != yaml.DocumentNode || len(doc.Content) == 0 {
		return fmt.Errorf("unexpected yaml structure")
	}

	mapping := doc.Content[0]
	if mapping.Kind != yaml.MappingNode {
		return fmt.Errorf("frontmatter is not a mapping")
	}

	setMappingField(mapping, "status", string(task.Column))

	if task.Stage == "" {
		removeMappingField(mapping, "phase")
		removeMappingField(mapping, "stage")
	} else {
		setMappingField(mapping, "phase", string(task.Stage))
		removeMappingField(mapping, "stage")
	}

	out, err := yaml.Marshal(&doc)
	if err != nil {
		return fmt.Errorf("marshal yaml: %w", err)
	}

	newContent := "---\n" + strings.TrimSpace(string(out)) + "\n---" + body

	return os.WriteFile(path, []byte(newContent), 0644)
}

func setMappingField(mapping *yaml.Node, key, value string) {
	for i := 0; i < len(mapping.Content)-1; i += 2 {
		if mapping.Content[i].Value == key {
			mapping.Content[i+1].Value = value
			mapping.Content[i+1].Tag = "!!str"
			return
		}
	}
	keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: key, Tag: "!!str"}
	valNode := &yaml.Node{Kind: yaml.ScalarNode, Value: value, Tag: "!!str"}
	mapping.Content = append(mapping.Content, keyNode, valNode)
}

func removeMappingField(mapping *yaml.Node, key string) {
	for i := 0; i < len(mapping.Content)-1; i += 2 {
		if mapping.Content[i].Value == key {
			mapping.Content = append(mapping.Content[:i], mapping.Content[i+2:]...)
			return
		}
	}
}

func WriteRouterState(root string, state *RouterState, expectedMtime time.Time) error {
	path := filepath.Join(root, "router.md")
	fi, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}

	if !fi.ModTime().Equal(expectedMtime) {
		return ErrMtimeConflict
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}

	normalized := normalizeLineEndings(string(content))

	startIdx := strings.Index(normalized, stateBlockStart)
	if startIdx == -1 {
		return fmt.Errorf("no Current state block found")
	}

	yamlStart := strings.Index(normalized[startIdx:], "```yaml")
	if yamlStart == -1 {
		return fmt.Errorf("no yaml code block found")
	}

	yamlStart += startIdx + len("```yaml")
	yamlEnd := strings.Index(normalized[yamlStart:], "```")
	if yamlEnd == -1 {
		return fmt.Errorf("no closing code block found")
	}

	out, err := yaml.Marshal(state)
	if err != nil {
		return fmt.Errorf("marshal yaml: %w", err)
	}

	newContent := normalized[:yamlStart] + "\n" + strings.TrimSpace(string(out)) + "\n" + normalized[yamlStart+yamlEnd:]

	return os.WriteFile(path, []byte(newContent), 0644)
}
