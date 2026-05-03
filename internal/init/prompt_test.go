package init

import (
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
