package data

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// RouterPhaseV2 is the V2 router's closed state vocabulary. It replaces V1's
// free-form state string with a typed, exhaustive set: idea, design, task,
// and check are the only phases a V2 router can name.
type RouterPhaseV2 string

const (
	RouterPhaseIdea   RouterPhaseV2 = "idea"
	RouterPhaseDesign RouterPhaseV2 = "design"
	RouterPhaseTask   RouterPhaseV2 = "task"
	RouterPhaseCheck  RouterPhaseV2 = "check"
)

// RouterStateV2 is the decoded "## Current state" anchor for a V2 project:
// a phase, an optional Release context, the selected Objective, an optional
// selected Task, and prose for a human. It is a hint only — record identity
// and cross-selection ownership are resolved against the project index by
// ResolveSelection.
type RouterStateV2 struct {
	State      RouterPhaseV2
	Release    string // R### selection, or empty when no Release is selected
	Objective  string // O### selection, or empty when none is selected
	Task       string // T### selection, or empty when none is selected
	NextAction string
}

// routerV2Frontmatter is the raw shape decoded off the anchor before
// validation. "none" and "" are both the router's written spelling of "not
// selected" — templates/project-v2/.savepoint/router.md ships "none" for a
// fresh project's objective/task, the same sentinel V1's release/epic
// fields already use (see internal/doctor/checks.go).
type routerV2Frontmatter struct {
	State      string `yaml:"state"`
	Release    string `yaml:"release"`
	Objective  string `yaml:"objective"`
	Task       string `yaml:"task"`
	NextAction string `yaml:"next_action"`
}

// ReadStateV2 decodes the "## Current state" anchor into a V2 RouterStateV2.
// It reuses extractStateBlock's anchor-finding — the same heading and fenced
// ```yaml block the V1 reader locates — and then decodes strictly (DATA-03):
// an unrecognized or empty state, a malformed R###/O###/T### selection, a task
// selected without an objective, and an unknown key each return a named
// diagnostic instead of a healed default. Decoding performs no filesystem
// write and no repair of content.
func (r *RouterReader) ReadStateV2(content string) (*RouterStateV2, error) {
	yamlContent, err := extractStateBlock(content)
	if err != nil {
		return nil, err
	}

	dec := yaml.NewDecoder(strings.NewReader(yamlContent))
	dec.KnownFields(true)

	var fields routerV2Frontmatter
	if err := dec.Decode(&fields); err != nil {
		return nil, fmt.Errorf("%w: router state: %v", ErrV2Malformed, err)
	}

	phase := RouterPhaseV2(fields.State)
	switch phase {
	case RouterPhaseIdea, RouterPhaseDesign, RouterPhaseTask, RouterPhaseCheck:
	default:
		return nil, fmt.Errorf("%w: router state %q; use idea, design, task, or check", ErrV2InvalidLifecycle, fields.State)
	}

	release := normalizeRouterSelectionV2(fields.Release)
	if release != "" && !releaseIDPatternV2.MatchString(release) {
		return nil, fmt.Errorf("%w: router release %q must be a single R### selection", ErrV2InvalidID, fields.Release)
	}

	objective := normalizeRouterSelectionV2(fields.Objective)
	if objective != "" && !objectiveIDPattern.MatchString(objective) {
		return nil, fmt.Errorf("%w: router objective %q must be a single O### selection", ErrV2InvalidID, fields.Objective)
	}

	task := normalizeRouterSelectionV2(fields.Task)
	if task != "" && !taskIDPatternV2.MatchString(task) {
		return nil, fmt.Errorf("%w: router task %q must be a single T### selection", ErrV2InvalidID, fields.Task)
	}

	if task != "" && objective == "" {
		return nil, fmt.Errorf("%w: router task %q selected with no objective", ErrV2InvalidOwnership, fields.Task)
	}

	return &RouterStateV2{
		State:      phase,
		Release:    release,
		Objective:  objective,
		Task:       task,
		NextAction: fields.NextAction,
	}, nil
}

// normalizeRouterSelectionV2 treats "none" the same as an absent selection,
// matching the sentinel the shipped V2 template writes for a fresh project.
func normalizeRouterSelectionV2(raw string) string {
	if raw == "none" {
		return ""
	}
	return raw
}
