package init

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

func runInitPipeline(t *testing.T, dir string, force bool) string {
	t.Helper()

	if err := ValidateTarget(dir, force); err != nil {
		t.Fatal(err)
	}

	templates := fstest.MapFS{
		".savepoint":                            &fstest.MapFile{Mode: fs.ModeDir | 0755},
		".savepoint/config.yml":                 &fstest.MapFile{Data: []byte("key: value")},
		".savepoint/Design.md":                  &fstest.MapFile{Data: []byte("# {{PROJECT_NAME}} Design")},
		".savepoint/PRD.md":                     &fstest.MapFile{Data: []byte("PRD: {{PROJECT_NAME}}")},
		".savepoint/router.md":                  &fstest.MapFile{Data: []byte("# Router")},
		".savepoint/visual-identity.md":         &fstest.MapFile{Data: []byte("# Visual Identity")},
		"AGENTS.md":                             &fstest.MapFile{Data: []byte("# Agents Guide\n\nBuild: npm run build")},
		"agent-skills/savepoint-audit/SKILL.md": &fstest.MapFile{Data: []byte("# Audit Skill")},
	}

	projectName := ProjectNameFromDir(dir)
	if err := Scaffold(templates, dir, projectName, force); err != nil {
		t.Fatal(err)
	}

	promptTemplates := fstest.MapFS{
		"magic-prompt.prompt.md": &fstest.MapFile{
			Data: []byte("<!-- AGENT: Read AGENTS.md -->\n\nProject: {{PROJECT_NAME}}\n\nStart by reading AGENTS.md"),
		},
	}

	prompt, err := RenderMagicPrompt(promptTemplates, projectName)
	if err != nil {
		t.Fatal(err)
	}

	return prompt
}

func TestIntegration_EmptyDirectory(t *testing.T) {
	dir := t.TempDir()
	prompt := runInitPipeline(t, dir, false)
	projectName := filepath.Base(dir)

	entries := []string{
		".savepoint/config.yml",
		".savepoint/Design.md",
		".savepoint/PRD.md",
		".savepoint/router.md",
		".savepoint/visual-identity.md",
		"AGENTS.md",
		"agent-skills/savepoint-audit/SKILL.md",
	}
	for _, e := range entries {
		if _, err := os.Stat(filepath.Join(dir, e)); err != nil {
			t.Errorf("missing %s: %v", e, err)
		}
	}

	data, err := os.ReadFile(filepath.Join(dir, ".savepoint", "Design.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), projectName) {
		t.Errorf("Design.md content = %q, want to contain %q", string(data), projectName)
	}

	if !strings.Contains(prompt, projectName) {
		t.Errorf("prompt = %q, want to contain %q", prompt, projectName)
	}
	if !strings.Contains(prompt, "AGENT") {
		t.Errorf("prompt = %q, want AGENT marker", prompt)
	}

	result := CopyToClipboard(prompt)
	if result.Status != ClipboardCopied &&
		result.Status != ClipboardSkipped &&
		result.Status != ClipboardFailed {
		t.Errorf("unexpected clipboard status: %v", result.Status)
	}
}

func TestIntegration_CompatibleProject(t *testing.T) {
	dir := t.TempDir()

	for _, name := range []string{"package.json", ".git", "README.md"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte{}, 0644); err != nil {
			t.Fatal(err)
		}
	}

	prompt := runInitPipeline(t, dir, false)

	if _, err := os.Stat(filepath.Join(dir, "package.json")); err != nil {
		t.Errorf("package.json missing: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, ".savepoint", "config.yml")); err != nil {
		t.Errorf(".savepoint/config.yml missing: %v", err)
	}

	if !strings.Contains(prompt, "AGENT") {
		t.Errorf("prompt should contain AGENT marker")
	}
}

func TestIntegration_NoForceOnExistingSavepoint(t *testing.T) {
	dir := t.TempDir()

	savepointDir := filepath.Join(dir, ".savepoint")
	if err := os.Mkdir(savepointDir, 0755); err != nil {
		t.Fatal(err)
	}

	err := ValidateTarget(dir, false)
	if err == nil {
		t.Fatal("expected error for existing .savepoint without --force")
	}
	if !strings.Contains(err.Error(), "already contains") {
		t.Errorf("error = %q, want 'already contains'", err.Error())
	}
}

func TestIntegration_ForceOverwritesExistingSavepoint(t *testing.T) {
	dir := t.TempDir()

	savepointDir := filepath.Join(dir, ".savepoint")
	if err := os.Mkdir(savepointDir, 0755); err != nil {
		t.Fatal(err)
	}

	prompt := runInitPipeline(t, dir, true)

	for _, path := range []string{
		".savepoint/config.yml", ".savepoint/Design.md",
		".savepoint/PRD.md", ".savepoint/router.md",
	} {
		if _, err := os.Stat(filepath.Join(dir, path)); err != nil {
			t.Errorf("expected %s to exist after --force: %v", path, err)
		}
	}

	if !strings.Contains(prompt, "AGENT") {
		t.Errorf("prompt should contain AGENT marker")
	}
}

func TestIntegration_InstallDependencies(t *testing.T) {
	if _, err := exec.LookPath("npm"); err != nil {
		t.Skip("npm not found in PATH, skipping install test")
	}

	dir := t.TempDir()

	packageJSON := filepath.Join(dir, "package.json")
	if err := os.WriteFile(packageJSON, []byte(`{"name":"test","version":"0.0.0"}`), 0644); err != nil {
		t.Fatal(err)
	}

	if err := ValidateTarget(dir, false); err != nil {
		t.Fatal(err)
	}

	if err := InstallDependencies(dir); err != nil {
		t.Fatalf("InstallDependencies() error = %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, "package-lock.json")); err != nil {
		t.Errorf("package-lock.json not created after npm install: %v", err)
	}
}

func TestIntegration_MagicPromptInOutput(t *testing.T) {
	dir := t.TempDir()
	prompt := runInitPipeline(t, dir, false)
	projectName := filepath.Base(dir)

	expectedParts := []string{
		projectName,
		"AGENT",
		"AGENTS.md",
		"Project",
	}
	for _, part := range expectedParts {
		if !strings.Contains(prompt, part) {
			t.Errorf("prompt missing %q", part)
		}
	}
}
