package data

import (
	"fmt"
	"regexp"
	"strings"
)

// taskIDPatternV2 anchors V2 global Task identity: T plus at least three
// digits.
var taskIDPatternV2 = regexp.MustCompile(`^T[0-9]{3,}$`)

// TaskDependencyRequirement is the clearance a V2 Task dependency needs
// before it can be treated as satisfied. It is decoded strictly: an omitted
// requirement defaults to clear, but every other unrecognized value is
// rejected rather than healed.
type TaskDependencyRequirement string

const (
	TaskDependencyClear    TaskDependencyRequirement = "clear"
	TaskDependencyAccepted TaskDependencyRequirement = "accepted"
)

// TaskDependencyV2 is one decoded {task: T###, requires: clear|accepted}
// dependency record.
type TaskDependencyV2 struct {
	Task     string
	Requires TaskDependencyRequirement
}

// TaskV2 is a strict V2 Task record decoded by DecodeTaskV2. It is a
// distinct type from Task (task.go): V2 identity, ownership, and lifecycle
// rules never broaden or leak into the active V1 Task contract.
type TaskV2 struct {
	ID        string
	Title     string
	Objective string // the single O### owner
	PlannedBy Actor  // planner provenance recorded when the Task is created
	Status    ColumnType
	Stage     ProgressStage
	DependsOn []TaskDependencyV2
	// Release is retained only for transitional V2 migration compatibility.
	// Release-aware indexing derives context from the owning Objective and
	// never reads this legacy packaging field as an ownership edge.
	Release  string
	Evidence *Evidence
	Source   V2SourceDocument
}

type taskDependencyV2Frontmatter struct {
	Task     string `yaml:"task"`
	Requires string `yaml:"requires"`
}

type taskV2Frontmatter struct {
	ID                    string                        `yaml:"id"`
	Title                 string                        `yaml:"title"`
	Objective             string                        `yaml:"objective"`
	PlannedBy             evidenceActorFrontmatter      `yaml:"planned_by"`
	Status                ColumnType                    `yaml:"status"`
	Stage                 ProgressStage                 `yaml:"stage"`
	DependsOn             []taskDependencyV2Frontmatter `yaml:"depends_on"`
	Release               string                        `yaml:"release"`
	evidenceV2Frontmatter `yaml:",inline"`
}

// DecodeTaskV2 strictly decodes a V2 Task record from content. It requires a
// valid global T### ID, a non-empty title distinct from any Objective
// reference, exactly one O### objective owner, a canonical lifecycle
// status/stage combination, planner provenance, and well-formed dependency
// records. Nothing here is healed: missing titles, malformed IDs, unknown
// lifecycle values, invalid provenance, and invalid dependency requirements
// all return named diagnostics.
func DecodeTaskV2(path, content string) (*TaskV2, error) {
	doc, err := ParseV2Document(path, content)
	if err != nil {
		return nil, err
	}

	var fields taskV2Frontmatter
	if err := doc.Frontmatter.Decode(&fields); err != nil {
		return nil, fmt.Errorf("%w: %s: %v", ErrV2Malformed, path, err)
	}

	if !taskIDPatternV2.MatchString(fields.ID) {
		return nil, fmt.Errorf("%w: %s: task id %q must match T plus at least three digits", ErrV2InvalidID, path, fields.ID)
	}

	// Title is required on its own terms. V1's objective-as-title fallback
	// (parser.go firstNonEmpty(Title, Objective)) never applies here: V2
	// objective is an O### identity reference, not display language.
	if strings.TrimSpace(fields.Title) == "" {
		return nil, fmt.Errorf("%w: %s: task %s missing required field title", ErrV2MissingField, path, fields.ID)
	}

	if !objectiveIDPattern.MatchString(fields.Objective) {
		return nil, fmt.Errorf("%w: %s: task %s objective %q must be a single O### owner", ErrV2InvalidOwnership, path, fields.ID, fields.Objective)
	}

	plannedBy, err := decodeTaskPlanner(path, fields.ID, fields.PlannedBy)
	if err != nil {
		return nil, err
	}

	// Reuse the strict, non-healing write-path validator so V2 lifecycle
	// rules stay in internal/data's single canonical vocabulary (DATA-02)
	// instead of re-deriving status/stage rules or routing through V1's
	// load-time defaulting/alias path in ParseTaskLifecycle.
	if err := ValidateTaskLifecycleStateForWrite(TaskLifecycleState{Status: fields.Status, Stage: fields.Stage}); err != nil {
		return nil, fmt.Errorf("%w: %s: task %s: %v", ErrV2InvalidLifecycle, path, fields.ID, err)
	}

	dependsOn := make([]TaskDependencyV2, 0, len(fields.DependsOn))
	for _, ref := range fields.DependsOn {
		if !taskIDPatternV2.MatchString(ref.Task) {
			return nil, fmt.Errorf("%w: %s: task %s depends_on %q must match T plus at least three digits", ErrV2InvalidDependency, path, fields.ID, ref.Task)
		}

		requires := TaskDependencyRequirement(ref.Requires)
		if requires == "" {
			requires = TaskDependencyClear
		} else if requires != TaskDependencyClear && requires != TaskDependencyAccepted {
			return nil, fmt.Errorf("%w: %s: task %s depends_on %s requires %q; use clear or accepted", ErrV2InvalidDependency, path, fields.ID, ref.Task, ref.Requires)
		}

		dependsOn = append(dependsOn, TaskDependencyV2{Task: ref.Task, Requires: requires})
	}

	evidence, err := decodeEvidenceV2(path, "task", fields.ID, fields.evidenceV2Frontmatter)
	if err != nil {
		return nil, err
	}

	return &TaskV2{
		ID:        fields.ID,
		Title:     fields.Title,
		Objective: fields.Objective,
		PlannedBy: plannedBy,
		Status:    fields.Status,
		Stage:     fields.Stage,
		DependsOn: dependsOn,
		Release:   fields.Release,
		Evidence:  evidence,
		Source:    doc,
	}, nil
}

// decodeTaskPlanner decodes the required planner provenance attached to a
// Task. The shared actor decoder enforces the role vocabulary and scalar
// session shape; this boundary additionally narrows the role to planner so a
// Task cannot borrow execution, checking, or owner authority as its plan.
func decodeTaskPlanner(path, taskID string, raw evidenceActorFrontmatter) (Actor, error) {
	plannedBy, err := decodeV2Actor(ErrV2TaskMalformed, path, "task", taskID, "planned_by", raw)
	if err != nil {
		return Actor{}, err
	}
	if plannedBy.Role != ActorRolePlanner {
		return Actor{}, fmt.Errorf("%w: %s: task %s planned_by.role %q; use planner", ErrV2TaskMalformed, path, taskID, plannedBy.Role)
	}
	if strings.TrimSpace(raw.Session) == "" {
		return Actor{}, fmt.Errorf("%w: %s: task %s missing required field planned_by.session", ErrV2MissingField, path, taskID)
	}

	return plannedBy, nil
}
