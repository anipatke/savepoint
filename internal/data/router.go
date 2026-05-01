package data

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

const stateBlockStart = "## Current state"
const stateBlockEnd = "```"

type RouterState struct {
	State      string `yaml:"state"`
	Release   string `yaml:"release"`
	Epic      string `yaml:"epic"`
	Task      string `yaml:"task"`
	NextAction string `yaml:"next_action"`
}

type RouterReader struct{}

func NewRouterReader() *RouterReader {
	return &RouterReader{}
}

func (r *RouterReader) ReadState(content string) (*RouterState, error) {
	startIdx := strings.Index(content, stateBlockStart)
	if startIdx == -1 {
		return nil, fmt.Errorf("no Current state block found")
	}

	yamlStart := strings.Index(content[startIdx:], "```yaml")
	if yamlStart == -1 {
		return nil, fmt.Errorf("no yaml code block found")
	}

	yamlStart += startIdx + len("```yaml")
	yamlEnd := strings.Index(content[yamlStart:], "```")
	if yamlEnd == -1 {
		return nil, fmt.Errorf("no closing code block found")
	}

	yamlContent := strings.TrimSpace(content[yamlStart : yamlStart+yamlEnd])

	var state RouterState
	if err := yaml.Unmarshal([]byte(yamlContent), &state); err != nil {
		return nil, fmt.Errorf("failed to parse router YAML: %w", err)
	}

	return &state, nil
}