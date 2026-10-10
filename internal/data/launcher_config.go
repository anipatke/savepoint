package data

import (
	"fmt"
	"strings"
	"time"
)

// TerminalMode selects how the launcher opens an interactive terminal.
type TerminalMode string

const (
	// TerminalModeAuto lets the launcher pick a platform terminal adapter.
	TerminalModeAuto TerminalMode = "auto"
	// TerminalModeOverride uses an explicit, structured terminal command.
	TerminalModeOverride TerminalMode = "override"
)

// Argument placeholders the launcher substitutes inside command arguments.
// They are kept here as the single source of truth so prompts (T002) and
// terminal adapters (T003) agree on the supported set.
const (
	// PlaceholderPrompt is replaced with the generated workflow prompt text.
	PlaceholderPrompt = "{{prompt}}"
	// PlaceholderProjectRoot is replaced with the absolute project root path.
	PlaceholderProjectRoot = "{{project_root}}"
)

// SupportedPlaceholders is the complete set of placeholders the launcher will
// substitute in agent profile and terminal override arguments.
var SupportedPlaceholders = []string{PlaceholderPrompt, PlaceholderProjectRoot}

// CommandSpec is a structured executable plus an ordered argument list. Args
// carry placeholders rather than shell command strings so the launcher never
// interpolates into a shell.
type CommandSpec struct {
	Command string   `yaml:"command"`
	Args    []string `yaml:"args"`
}

// TerminalConfig controls how an interactive terminal is opened. In auto mode
// the launcher chooses a platform adapter; in override mode it uses Override.
type TerminalConfig struct {
	Mode     TerminalMode `yaml:"mode"`
	Override *CommandSpec `yaml:"override"`
}

// AgentLauncher is the disabled-by-default, vendor-neutral launcher contract.
// Builder and Auditor are optional; an absent Builder means DefaultAgentBuilder.
type AgentLauncher struct {
	Enabled  bool           `yaml:"enabled"`
	Builder  *CommandSpec   `yaml:"builder"`
	Auditor  *CommandSpec   `yaml:"auditor"`
	Terminal TerminalConfig `yaml:"terminal"`
	// Timeout bounds one headless run, written like quality_gates.gate_timeout
	// ("30m"). Use RunTimeout, which applies the default.
	Timeout string `yaml:"timeout"`
}

// DefaultRunTimeout is how long a headless run may last when
// agent_launcher.timeout is unset or not a positive duration.
const DefaultRunTimeout = 30 * time.Minute

// RunTimeout is the configured run timeout. It is parsed as the gate timeout
// is: a value that is not a positive duration falls back to the default.
func (l AgentLauncher) RunTimeout() time.Duration {
	if d, err := time.ParseDuration(l.Timeout); err == nil && d > 0 {
		return d
	}
	return DefaultRunTimeout
}

// fillLauncherDefaults applies defaults that are independent of YAML parsing.
// Terminal mode defaults to auto so callers always have a concrete mode.
func fillLauncherDefaults(launcher AgentLauncher) AgentLauncher {
	if launcher.Terminal.Mode == "" {
		launcher.Terminal.Mode = TerminalModeAuto
	}
	return launcher
}

// Validate reports the exact missing or invalid launcher field. A disabled
// launcher is always valid regardless of the other fields.
func (l AgentLauncher) Validate() error {
	if !l.Enabled {
		return nil
	}

	// A missing builder is allowed: headless runs then use DefaultAgentBuilder.
	if l.Builder != nil {
		if err := validateProfileCommand("agent_launcher.builder", l.Builder); err != nil {
			return err
		}
	}
	if l.Auditor != nil {
		if err := validateProfileCommand("agent_launcher.auditor", l.Auditor); err != nil {
			return err
		}
	}
	return validateTerminal(l.Terminal)
}

func validateRequiredProfile(field string, profile *CommandSpec) error {
	if profile == nil {
		return fmt.Errorf("%s is required when agent_launcher.enabled is true", field)
	}
	return validateProfileCommand(field, profile)
}

func validateProfileCommand(field string, profile *CommandSpec) error {
	if profile.Command == "" {
		return fmt.Errorf("%s.command is required", field)
	}
	return nil
}

func validateTerminal(terminal TerminalConfig) error {
	switch terminal.Mode {
	case TerminalModeAuto:
		return nil
	case TerminalModeOverride:
		return validateRequiredProfile("agent_launcher.terminal.override", terminal.Override)
	default:
		return fmt.Errorf("agent_launcher.terminal.mode %q is invalid (want auto or override)", terminal.Mode)
	}
}

// Default Claude Code command for headless runs, used when the owner has not
// set agent_launcher.builder. The flags are content, not logic: a later provider
// is another table (O-045 Confirmed Design). Argument order matters: the
// variadic --allowedTools follows the prompt and precedes a flag.
const (
	// DefaultBuilderCommand is the executable the default builder starts.
	DefaultBuilderCommand = "claude"
	// AllowedToolsFlag introduces the list DefaultAgentBuilder derives from
	// quality_gates; it is omitted when no gate names a command.
	AllowedToolsFlag = "--allowedTools"
)

var (
	defaultBuilderLeadArgs = []string{
		"-p", PlaceholderPrompt, "--output-format", "stream-json", "--verbose",
		"--permission-mode", "acceptEdits",
	}
	defaultBuilderTailArgs = []string{"--permission-prompts", "none"}
)

// DefaultAgentBuilder is the Claude command used when agent_launcher.builder is
// absent. It allows the first word of each configured quality gate as
// `Bash(<word> *)` so the agent can run its own gates and nothing else asks for
// approval; a gate that is unset allows nothing.
func DefaultAgentBuilder(gates QualityGates) CommandSpec {
	args := append([]string(nil), defaultBuilderLeadArgs...)
	var allowed []string
	seen := map[string]bool{}
	for _, gate := range []*string{gates.Lint, gates.Typecheck, gates.Build, gates.Test} {
		if gate == nil {
			continue
		}
		words := strings.Fields(*gate)
		if len(words) == 0 || seen[words[0]] {
			continue
		}
		seen[words[0]] = true
		allowed = append(allowed, "Bash("+words[0]+" *)")
	}
	if len(allowed) > 0 {
		args = append(args, AllowedToolsFlag)
		args = append(args, allowed...)
	}
	return CommandSpec{Command: DefaultBuilderCommand, Args: append(args, defaultBuilderTailArgs...)}
}
