package run

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/opencode/savepoint/internal/data"
)

// BuildArgv turns the launcher's builder into the argument vector for a run:
// the command followed by its arguments, with {{prompt}} replaced by nextLine
// exactly and {{project_root}} by the absolute project root. No shell is
// involved, so neither value is quoted or interpreted. A nil builder means the
// owner has not set one and the default Claude command is used; a builder with
// no command is an error, returned before anything is written (CFG-01).
func BuildArgv(builder *data.CommandSpec, gates data.QualityGates, nextLine, projectRoot string) ([]string, error) {
	spec := data.DefaultAgentBuilder(gates)
	if builder != nil {
		spec = *builder
	}
	if strings.TrimSpace(spec.Command) == "" {
		return nil, errors.New("agent_launcher.builder.command is required: set the command that starts the agent")
	}
	root, err := filepath.Abs(projectRoot)
	if err != nil {
		return nil, fmt.Errorf("run: resolve project root %q: %w", projectRoot, err)
	}
	replacer := strings.NewReplacer(data.PlaceholderPrompt, nextLine, data.PlaceholderProjectRoot, root)
	argv := []string{spec.Command}
	for _, arg := range spec.Args {
		argv = append(argv, replacer.Replace(arg))
	}
	return argv, nil
}
