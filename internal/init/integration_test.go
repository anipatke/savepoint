package init

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/opencode/savepoint/internal/testutil"
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
		".savepoint/releases/v1/epics":          &fstest.MapFile{Mode: fs.ModeDir | 0755},
		".savepoint/releases/v1/v1-PRD.md":      &fstest.MapFile{Data: []byte("# v1 PRD for {{PROJECT_NAME}}")},
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
		".savepoint/releases/v1/v1-PRD.md",
		"AGENTS.md",
		"agent-skills/savepoint-audit/SKILL.md",
	}

	if info, err := os.Stat(filepath.Join(dir, ".savepoint", "releases", "v1", "epics")); err != nil || !info.IsDir() {
		t.Errorf(".savepoint/releases/v1/epics not created as directory: %v", err)
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
		testutil.WriteFile(t, filepath.Join(dir, name), "")
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

func TestIntegration_ExistingAgentGuide(t *testing.T) {
	dir := t.TempDir()
	existingPath := filepath.Join(dir, "AGENTS.md")
	testutil.WriteFile(t, existingPath, "# Team Guide\n\nOur custom rules.")

	prompt := runInitPipeline(t, dir, false)

	data, err := os.ReadFile(existingPath)
	if err != nil {
		t.Fatal(err)
	}
	got := string(data)
	if !strings.Contains(got, "# Team Guide") {
		t.Errorf("AGENTS.md missing existing content: %q", got)
	}
	if !strings.Contains(got, "Our custom rules.") {
		t.Errorf("AGENTS.md missing existing custom rules: %q", got)
	}
	if !strings.Contains(got, "<!-- SAVEPOINT:BEGIN -->") {
		t.Errorf("AGENTS.md missing managed block: %q", got)
	}
	if !strings.Contains(prompt, "AGENT") {
		t.Errorf("prompt should contain AGENT marker")
	}
}

func TestIntegration_ExistingAgentGuideCasingVariant(t *testing.T) {
	dir := t.TempDir()
	variantPath := filepath.Join(dir, "Agents.MD")
	testutil.WriteFile(t, variantPath, "# Team Guide")

	runInitPipeline(t, dir, false)

	// Exactly one agent guide file should exist (no duplicate created).
	entries, _ := os.ReadDir(dir)
	guideCount := 0
	for _, e := range entries {
		if strings.ToLower(e.Name()) == "agents.md" {
			guideCount++
		}
	}
	if guideCount != 1 {
		t.Errorf("expected 1 agent guide file, found %d", guideCount)
	}

	data, err := os.ReadFile(variantPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "# Team Guide") {
		t.Errorf("Agents.MD missing original content: %q", string(data))
	}
	if !strings.Contains(string(data), "<!-- SAVEPOINT:BEGIN -->") {
		t.Errorf("Agents.MD missing managed block: %q", string(data))
	}
}

func TestIntegration_InstallDependencies(t *testing.T) {
	if _, err := exec.LookPath("npm"); err != nil {
		t.Skip("npm not found in PATH, skipping install test")
	}

	dir := t.TempDir()

	packageJSON := filepath.Join(dir, "package.json")
	testutil.WriteFile(t, packageJSON, `{"name":"test","version":"0.0.0"}`)

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

func TestIntegration_MagicPromptDoesNotCarryPhaseWorkflow(t *testing.T) {
	dir := t.TempDir()
	prompt := runInitPipeline(t, dir, false)

	stalePhaseInstructions := []string{
		"pre-implementation",
		"epic-design",
		"epic-task-breakdown",
		"task-building",
		"audit-pending",
		"defect-building",
		"phase prompt",
		"prompt-based phase",
	}
	for _, stale := range stalePhaseInstructions {
		if strings.Contains(prompt, stale) {
			t.Fatalf("prompt contains stale phase workflow text %q", stale)
		}
	}
}
