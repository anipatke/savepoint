package data

import (
	"os"
	"strings"
	"testing"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	tmpfile, err := os.CreateTemp("", "config-*.yml")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(tmpfile.Name()) })
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()
	return tmpfile.Name()
}

func TestLauncherDefaultsWhenAbsent(t *testing.T) {
	r := NewConfigReader()

	// Missing file and a config without an agent_launcher block both default
	// to a disabled launcher with auto terminal mode.
	for name, path := range map[string]string{
		"missing file":    "nonexistent.yml",
		"no launcher key": writeTempConfig(t, "theme:\n  bg: \"#000000\"\n"),
	} {
		t.Run(name, func(t *testing.T) {
			config, err := r.Read(path)
			if err != nil {
				t.Fatalf("Read() error = %v", err)
			}
			if config.AgentLauncher.Enabled {
				t.Errorf("AgentLauncher.Enabled = true, want false by default")
			}
			if config.AgentLauncher.Terminal.Mode != TerminalModeAuto {
				t.Errorf("Terminal.Mode = %q, want %q", config.AgentLauncher.Terminal.Mode, TerminalModeAuto)
			}
		})
	}
}

func TestLauncherDisabledSkipsValidation(t *testing.T) {
	// A disabled launcher with no builder is still valid.
	path := writeTempConfig(t, `agent_launcher:
  enabled: false
`)
	config, err := NewConfigReader().Read(path)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if config.AgentLauncher.Enabled {
		t.Errorf("Enabled = true, want false")
	}
}

func TestLauncherEnabledFull(t *testing.T) {
	path := writeTempConfig(t, `agent_launcher:
  enabled: true
  builder:
    command: claude
    args: ["{{prompt}}"]
  auditor:
    command: claude
    args: ["--audit", "{{prompt}}"]
  terminal:
    mode: override
    override:
      command: wezterm
      args: ["start", "--cwd", "{{project_root}}", "--"]
`)
	config, err := NewConfigReader().Read(path)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	l := config.AgentLauncher
	if !l.Enabled {
		t.Fatal("Enabled = false, want true")
	}
	if l.Builder == nil || l.Builder.Command != "claude" {
		t.Errorf("Builder = %+v, want command claude", l.Builder)
	}
	if l.Builder == nil || len(l.Builder.Args) != 1 || l.Builder.Args[0] != PlaceholderPrompt {
		t.Errorf("Builder.Args = %v, want [%s]", l.Builder.Args, PlaceholderPrompt)
	}
	if l.Auditor == nil || l.Auditor.Command != "claude" {
		t.Errorf("Auditor = %+v, want command claude", l.Auditor)
	}
	if l.Terminal.Mode != TerminalModeOverride {
		t.Errorf("Terminal.Mode = %q, want override", l.Terminal.Mode)
	}
	if l.Terminal.Override == nil || l.Terminal.Override.Command != "wezterm" {
		t.Errorf("Terminal.Override = %+v, want command wezterm", l.Terminal.Override)
	}
}

func TestLauncherEnabledOmitsOptionalAuditor(t *testing.T) {
	path := writeTempConfig(t, `agent_launcher:
  enabled: true
  builder:
    command: claude
`)
	config, err := NewConfigReader().Read(path)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if config.AgentLauncher.Auditor != nil {
		t.Errorf("Auditor = %+v, want nil when omitted", config.AgentLauncher.Auditor)
	}
}

func TestLauncherValidateErrors(t *testing.T) {
	tests := []struct {
		name      string
		launcher  AgentLauncher
		wantField string
	}{
		{
			name:      "enabled without builder",
			launcher:  AgentLauncher{Enabled: true, Terminal: TerminalConfig{Mode: TerminalModeAuto}},
			wantField: "agent_launcher.builder is required",
		},
		{
			name:      "builder without command",
			launcher:  AgentLauncher{Enabled: true, Builder: &CommandSpec{Args: []string{"x"}}, Terminal: TerminalConfig{Mode: TerminalModeAuto}},
			wantField: "agent_launcher.builder.command is required",
		},
		{
			name:      "auditor present without command",
			launcher:  AgentLauncher{Enabled: true, Builder: &CommandSpec{Command: "claude"}, Auditor: &CommandSpec{Args: []string{"x"}}, Terminal: TerminalConfig{Mode: TerminalModeAuto}},
			wantField: "agent_launcher.auditor.command is required",
		},
		{
			name:      "invalid terminal mode",
			launcher:  AgentLauncher{Enabled: true, Builder: &CommandSpec{Command: "claude"}, Terminal: TerminalConfig{Mode: "popup"}},
			wantField: "agent_launcher.terminal.mode",
		},
		{
			name:      "override mode without override command",
			launcher:  AgentLauncher{Enabled: true, Builder: &CommandSpec{Command: "claude"}, Terminal: TerminalConfig{Mode: TerminalModeOverride}},
			wantField: "agent_launcher.terminal.override is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.launcher.Validate()
			if err == nil {
				t.Fatalf("Validate() = nil, want error naming %q", tt.wantField)
			}
			if !strings.Contains(err.Error(), tt.wantField) {
				t.Errorf("Validate() error = %q, want to name %q", err.Error(), tt.wantField)
			}
		})
	}
}

func TestLauncherEnabledMissingBuilderFailsRead(t *testing.T) {
	path := writeTempConfig(t, `agent_launcher:
  enabled: true
`)
	if _, err := NewConfigReader().Read(path); err == nil {
		t.Fatal("Read() = nil error, want launcher validation error")
	}
}

func TestExistingConfigUnaffectedByLauncher(t *testing.T) {
	// Theme and quality-gate parsing must remain unchanged when no launcher
	// block is present.
	test := "go test ./..."
	path := writeTempConfig(t, `theme:
  bg: "#000000"
quality_gates:
  test: "go test ./..."
  block_on_failure: true
`)
	config, err := NewConfigReader().Read(path)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if config.Theme.BG != "#000000" {
		t.Errorf("Theme.BG = %v, want #000000", config.Theme.BG)
	}
	if config.QualityGates.Test == nil || *config.QualityGates.Test != test {
		t.Errorf("QualityGates.Test = %v, want %q", config.QualityGates.Test, test)
	}
	if !config.QualityGates.BlockOnFailure {
		t.Errorf("QualityGates.BlockOnFailure = false, want true")
	}
}

func TestSupportedPlaceholdersDefined(t *testing.T) {
	want := map[string]bool{PlaceholderPrompt: false, PlaceholderProjectRoot: false}
	for _, p := range SupportedPlaceholders {
		if _, ok := want[p]; !ok {
			t.Errorf("unexpected placeholder %q in SupportedPlaceholders", p)
		}
		want[p] = true
	}
	for p, seen := range want {
		if !seen {
			t.Errorf("placeholder %q missing from SupportedPlaceholders", p)
		}
	}
}
