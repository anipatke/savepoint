package data

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

const stateBlockStart = "## Current state"
const stateBlockEnd = "```"

// RouterState is the legacy V1 router projection. It is read only by
// migration/history compatibility code; live commands decode RouterStateV2.
type RouterState struct {
	State      string `yaml:"state"`
	Release    string `yaml:"release"`
	Epic       string `yaml:"epic"`
	Task       string `yaml:"task"`
	Defect     string `yaml:"defect,omitempty"`
	NextAction string `yaml:"next_action"`
}

type RouterReader struct{}

func NewRouterReader() *RouterReader {
	return &RouterReader{}
}

func (r *RouterReader) ReadState(content string) (*RouterState, error) {
	yamlContent, err := extractStateBlock(content)
	if err != nil {
		return nil, err
	}

	var state RouterState
	if err := yaml.Unmarshal([]byte(yamlContent), &state); err != nil {
		return nil, fmt.Errorf("failed to parse router YAML: %w", err)
	}

	return &state, nil
}

// extractStateBlock locates the "## Current state" heading and its fenced
// ```yaml block, returning the block's trimmed content. Both the V1 reader
// above and the V2 reader in router_v2.go share this anchor-finding step;
// only the shape decoded out of the returned text differs between them.
func extractStateBlock(content string) (string, error) {
	normalized := normalizeLineEndings(content)
	lower := strings.ToLower(normalized)

	startIdx := strings.Index(lower, strings.ToLower(stateBlockStart))
	if startIdx == -1 {
		return "", fmt.Errorf("no Current state block found")
	}

	yamlStart := strings.Index(normalized[startIdx:], "```yaml")
	if yamlStart == -1 {
		return "", fmt.Errorf("no yaml code block found")
	}

	yamlStart += startIdx + len("```yaml")
	yamlEnd := strings.Index(normalized[yamlStart:], "```")
	if yamlEnd == -1 {
		return "", fmt.Errorf("no closing code block found")
	}

	return strings.TrimSpace(normalized[yamlStart : yamlStart+yamlEnd]), nil
}

// ReplaceStateBlock returns content with the "## Current state" fenced YAML
// block's contents replaced by newYAML (trimmed of surrounding whitespace),
// preserving every other byte — the prose above and below the anchor —
// untouched. It is the writer half of ReadState's anchor-finding: a migrated
// router still lives in the same document shape, only its vocabulary
// changes.
func ReplaceStateBlock(content, newYAML string) (string, error) {
	normalized := normalizeLineEndings(content)
	lower := strings.ToLower(normalized)

	startIdx := strings.Index(lower, strings.ToLower(stateBlockStart))
	if startIdx == -1 {
		return "", fmt.Errorf("no Current state block found")
	}

	yamlStart := strings.Index(normalized[startIdx:], "```yaml")
	if yamlStart == -1 {
		return "", fmt.Errorf("no yaml code block found")
	}
	yamlStart += startIdx + len("```yaml")

	yamlEnd := strings.Index(normalized[yamlStart:], "```")
	if yamlEnd == -1 {
		return "", fmt.Errorf("no closing code block found")
	}
	yamlEnd += yamlStart

	return normalized[:yamlStart] + "\n" + strings.TrimSpace(newYAML) + "\n" + normalized[yamlEnd:], nil
}
