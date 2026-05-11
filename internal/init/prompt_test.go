package init

import (
	"os"
	"sort"
	"strings"
	"testing"
	"testing/fstest"
)

func TestRenderMagicPrompt_rendersTemplate(t *testing.T) {
	templates := fstest.MapFS{
		"magic-prompt.prompt.md": &fstest.MapFile{
			Data: []byte("Project: {{PROJECT_NAME}}"),
		},
	}

	got, err := RenderMagicPrompt(templates, "myapp")
	if err != nil {
		t.Fatalf("RenderMagicPrompt() error = %v", err)
	}

	want := "Project: myapp"
	if got != want {
		t.Fatalf("RenderMagicPrompt() = %q, want %q", got, want)
	}
}

func TestRenderMagicPrompt_interpolatesAllVariables(t *testing.T) {
	templates := fstest.MapFS{
		"magic-prompt.prompt.md": &fstest.MapFile{
			Data: []byte("{{PROJECT_NAME}} v{{RELEASE_NUMBER}}"),
		},
	}

	got, err := RenderMagicPrompt(templates, "myapp")
	if err != nil {
		t.Fatalf("RenderMagicPrompt() error = %v", err)
	}

	want := "myapp v1"
	if got != want {
		t.Fatalf("RenderMagicPrompt() = %q, want %q", got, want)
	}
}

func TestRenderMagicPrompt_handlesMissingTemplate(t *testing.T) {
	_, err := RenderMagicPrompt(fstest.MapFS{}, "myapp")
	if err == nil {
		t.Fatal("RenderMagicPrompt() expected error for missing template")
	}
}

func TestRenderMagicPrompt_usesEmbeddedTemplate(t *testing.T) {
	templates := fstest.MapFS{
		"magic-prompt.prompt.md": &fstest.MapFile{
			Data: []byte("<!-- AGENT: Read AGENTS.md -->\n\nProject: {{PROJECT_NAME}}"),
		},
	}

	got, err := RenderMagicPrompt(templates, "my-project")
	if err != nil {
		t.Fatalf("RenderMagicPrompt() error = %v", err)
	}

	if !strings.Contains(got, "my-project") {
		t.Fatalf("RenderMagicPrompt() = %q, does not contain project name", got)
	}
	if !strings.Contains(got, "AGENT") {
		t.Fatalf("RenderMagicPrompt() = %q, does not contain template content", got)
	}
}

func TestRenderMagicPrompt_handlesEmptyMapFS(t *testing.T) {
	_, err := RenderMagicPrompt(fstest.MapFS{}, "")
	if err == nil {
		t.Fatal("RenderMagicPrompt() expected error for empty MapFS")
	}
}

func TestPromptTemplates_onlyMagicPromptRemains(t *testing.T) {
	entries, err := os.ReadDir("../../templates/prompts")
	if err != nil {
		t.Fatalf("ReadDir() error = %v", err)
	}

	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			t.Fatalf("prompt templates include directory %s, want files only", entry.Name())
		}
		names = append(names, entry.Name())
	}
	sort.Strings(names)

	if len(names) != 1 || names[0] != "magic-prompt.prompt.md" {
		t.Fatalf("prompt templates = %v, want [magic-prompt.prompt.md]", names)
	}
}

func TestPromptTemplates_magicPromptIsBootstrapOnly(t *testing.T) {
	data, err := os.ReadFile("../../templates/prompts/magic-prompt.prompt.md")
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	content := string(data)
	for _, want := range []string{"{{PROJECT_NAME}}", "AGENTS.md", "AGENT"} {
		if !strings.Contains(content, want) {
			t.Fatalf("magic prompt missing %q", want)
		}
	}

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
		if strings.Contains(content, stale) {
			t.Fatalf("magic prompt contains stale phase instruction %q", stale)
		}
	}
}
