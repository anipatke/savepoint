package doctor

import (
	"errors"
	"fmt"
	"strings"

	"github.com/opencode/savepoint/internal/data"
)

func SuggestRepair(err error) string {
	switch {
	case errors.Is(err, data.ErrConfigNotFound):
		return "Run `savepoint init` to scaffold a new project"
	case errors.Is(err, data.ErrInvalidStatus):
		return "Set router state to a recognized workflow state (see router.md State → action section)"
	case errors.Is(err, data.ErrMissingFrontmatter):
		return "Fix the YAML frontmatter between the --- delimiters"
	case errors.Is(err, data.ErrStructureProblem):
		return "Review the file and fix the reported issue"
	}

	msg := err.Error()
	switch {
	case strings.Contains(msg, "config.yml not found"):
		return "Run `savepoint init` to scaffold a new project"
	case strings.Contains(msg, "config.yml missing required field"):
		return "Add the missing field to config.yml — see the project template for reference"
	case strings.Contains(msg, "invalid YAML"):
		return "Fix the YAML syntax error at the indicated line"
	case strings.Contains(msg, "router.md not found"):
		return "Run `savepoint init` to scaffold a new project"
	case strings.Contains(msg, "unknown state"):
		return "Set router state to a recognized workflow state (see router.md State → action section)"
	case strings.Contains(msg, "release PRD file not found"):
		return "Create a {release}-PRD.md file with frontmatter for the release"
	case strings.Contains(msg, "release"):
		return "Create the release directory at releases/<release-id>/"
	case strings.Contains(msg, "epic") && strings.Contains(msg, "directory not found"):
		return "Create the epic directory at releases/<release>/epics/<epic-id>/"
	case strings.Contains(msg, "epic detail file not found"):
		return "Create an E##-Detail.md with frontmatter for the epic"
	case strings.Contains(msg, "invalid frontmatter"):
		return "Fix the YAML frontmatter between the --- delimiters"
	case strings.Contains(msg, "task missing required frontmatter field"):
		return "Add the missing field to the task frontmatter"
	case strings.Contains(msg, "missing ## Acceptance Criteria"):
		return "Add an ## Acceptance Criteria section with checkable items"
	case strings.Contains(msg, "depends_on must be a list"):
		return "Change depends_on to a YAML list format"
	case strings.Contains(msg, "references non-existent"):
		return "Create the referenced task or remove the dependency"
	case strings.Contains(msg, "duplicate task ID"):
		return "Rename one of the tasks to have a unique ID"
	case strings.Contains(msg, "dependency cycle"):
		return "Break the circular dependency chain between tasks"
	case strings.Contains(msg, "audit proposal exists"):
		return "Set router state to audit-pending for the matching epic, or remove stale audit files"
	case strings.Contains(msg, "orphaned"):
		return "Move the task directory to the correct epic or create the referenced epic"
	case strings.Contains(msg, "quality gate"):
		return "Fix the issue reported by the quality gate tool"
	default:
		return "Review the file and fix the reported issue"
	}
}

// GateSuggestion returns a command-specific repair hint.
func GateSuggestion(name string) string {
	switch name {
	case "lint":
		return "Run `make lint` locally and fix reported issues"
	case "typecheck":
		return "Run `make typecheck` locally and fix type errors"
	case "test":
		return "Run `make test` locally and fix failing tests"
	default:
		return fmt.Sprintf("Run %q locally and fix reported issues", name)
	}
}
