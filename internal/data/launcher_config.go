package data

import "fmt"

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
// Builder is required when Enabled; Auditor is always optional.
type AgentLauncher struct {
	Enabled  bool           `yaml:"enabled"`
	Builder  *CommandSpec   `yaml:"builder"`
	Auditor  *CommandSpec   `yaml:"auditor"`
	Terminal TerminalConfig `yaml:"terminal"`
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

	if err := validateRequiredProfile("agent_launcher.builder", l.Builder); err != nil {
		return err
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
