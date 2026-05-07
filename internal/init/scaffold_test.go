package init

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/opencode/savepoint/internal/testutil"
)

func TestScaffold_createsDirectories(t *testing.T) {
	target := t.TempDir()
	templates := fstest.MapFS{
		".savepoint":            &fstest.MapFile{Mode: fs.ModeDir | 0755},
		".savepoint/config.yml": &fstest.MapFile{Data: []byte("key: value")},
		"AGENTS.md":             &fstest.MapFile{Data: []byte("# Agents Guide")},
	}

	if err := Scaffold(templates, target, "myapp", false); err != nil {
		t.Fatalf("Scaffold() error = %v", err)
	}

	if _, err := os.Stat(filepath.Join(target, ".savepoint", "config.yml")); err != nil {
		t.Errorf(".savepoint/config.yml not created: %v", err)
	}
	if _, err := os.Stat(filepath.Join(target, "AGENTS.md")); err != nil {
		t.Errorf("AGENTS.md not created: %v", err)
	}
}

func TestScaffold_interpolatesProjectName(t *testing.T) {
	target := t.TempDir()
	templates := fstest.MapFS{
		"Design.md": &fstest.MapFile{Data: []byte("# {{PROJECT_NAME}} Design")},
		"PRD.md":    &fstest.MapFile{Data: []byte("Project: {{PROJECT_NAME}}")},
	}

	if err := Scaffold(templates, target, "myapp", false); err != nil {
		t.Fatalf("Scaffold() error = %v", err)
	}

	data, err := os.ReadFile(filepath.Join(target, "Design.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "# myapp Design" {
		t.Fatalf("Design.md = %q, want %q", string(data), "# myapp Design")
	}

	data, err = os.ReadFile(filepath.Join(target, "PRD.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "Project: myapp" {
		t.Fatalf("PRD.md = %q, want %q", string(data), "Project: myapp")
	}
}

func TestScaffold_interpolatesReleaseNumber(t *testing.T) {
	target := t.TempDir()
	templates := fstest.MapFS{
		"AGENTS.md": &fstest.MapFile{Data: []byte("release: v{{RELEASE_NUMBER}}")},
	}

	if err := Scaffold(templates, target, "myapp", false); err != nil {
		t.Fatalf("Scaffold() error = %v", err)
	}

	data, err := os.ReadFile(filepath.Join(target, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "release: v1") {
		t.Fatalf("AGENTS.md = %q, want to contain %q", string(data), "release: v1")
	}
}

func TestScaffold_agentGuideInsertsBlock(t *testing.T) {
	target := t.TempDir()
	existingPath := filepath.Join(target, "AGENTS.md")
	testutil.WriteFile(t, existingPath, "# My Guide\n\nExisting content.")

	templates := fstest.MapFS{
		"AGENTS.md": &fstest.MapFile{Data: []byte("# Savepoint Instructions")},
	}

	if err := Scaffold(templates, target, "myapp", false); err != nil {
		t.Fatalf("Scaffold() error = %v", err)
	}

	data, err := os.ReadFile(existingPath)
	if err != nil {
		t.Fatal(err)
	}
	got := string(data)
	if !strings.Contains(got, "# My Guide") {
		t.Errorf("AGENTS.md missing existing content: %q", got)
	}
	if !strings.Contains(got, "# Savepoint Instructions") {
		t.Errorf("AGENTS.md missing managed block content: %q", got)
	}
	if !strings.Contains(got, managedBegin) {
		t.Errorf("AGENTS.md missing managed block begin marker: %q", got)
	}
}

func TestScaffold_agentGuideIdempotent(t *testing.T) {
	target := t.TempDir()
	templates := fstest.MapFS{
		"AGENTS.md": &fstest.MapFile{Data: []byte("# Savepoint Instructions")},
	}

	for i := range 2 {
		if err := Scaffold(templates, target, "myapp", false); err != nil {
			t.Fatalf("Scaffold() run %d error = %v", i+1, err)
		}
	}

	data, err := os.ReadFile(filepath.Join(target, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(data)
	count := strings.Count(got, managedBegin)
	if count != 1 {
		t.Errorf("AGENTS.md has %d managed begin markers, want 1: %q", count, got)
	}
}

func TestScaffold_agentGuideCasingVariant(t *testing.T) {
	target := t.TempDir()
	variantPath := filepath.Join(target, "Agents.MD")
	testutil.WriteFile(t, variantPath, "# My Guide")

	templates := fstest.MapFS{
		"AGENTS.md": &fstest.MapFile{Data: []byte("# Savepoint Instructions")},
	}

	if err := Scaffold(templates, target, "myapp", false); err != nil {
		t.Fatalf("Scaffold() error = %v", err)
	}

	// Exactly one agent guide file should exist (no duplicate created).
	entries, _ := os.ReadDir(target)
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
	got := string(data)
	if !strings.Contains(got, "# My Guide") {
		t.Errorf("Agents.MD missing existing content: %q", got)
	}
	if !strings.Contains(got, "# Savepoint Instructions") {
		t.Errorf("Agents.MD missing managed block content: %q", got)
	}
}

func TestScaffold_agentGuideForcePreservesUserContent(t *testing.T) {
	target := t.TempDir()
	existingPath := filepath.Join(target, "AGENTS.md")
	testutil.WriteFile(t, existingPath, "# My Guide\n\nUser content.\n\n"+managedBegin+"\nold block\n"+managedEnd+"\n")

	templates := fstest.MapFS{
		"AGENTS.md": &fstest.MapFile{Data: []byte("# New Savepoint Block")},
	}

	if err := Scaffold(templates, target, "myapp", true); err != nil {
		t.Fatalf("Scaffold() with force error = %v", err)
	}

	data, err := os.ReadFile(existingPath)
	if err != nil {
		t.Fatal(err)
	}
	got := string(data)
	if !strings.Contains(got, "User content.") {
		t.Errorf("AGENTS.md missing user content after force: %q", got)
	}
	if strings.Contains(got, "old block") {
		t.Errorf("AGENTS.md still has old managed block after force: %q", got)
	}
	if !strings.Contains(got, "# New Savepoint Block") {
		t.Errorf("AGENTS.md missing new managed content after force: %q", got)
	}
}

func TestScaffold_createsParentDirs(t *testing.T) {
	target := t.TempDir()
	templates := fstest.MapFS{
		"deep/nested/dir/file.txt": &fstest.MapFile{Data: []byte("content")},
	}

	if err := Scaffold(templates, target, "myapp", false); err != nil {
		t.Fatalf("Scaffold() error = %v", err)
	}

	path := filepath.Join(target, "deep", "nested", "dir", "file.txt")
	if _, err := os.Stat(path); err != nil {
		t.Errorf("deep/nested/dir/file.txt not created: %v", err)
	}
}

func TestScaffold_overwritesWithForce(t *testing.T) {
	target := t.TempDir()
	existingPath := filepath.Join(target, "file.txt")
	testutil.WriteFile(t, existingPath, "old")

	templates := fstest.MapFS{
		"file.txt": &fstest.MapFile{Data: []byte("new")},
	}

	if err := Scaffold(templates, target, "myapp", true); err != nil {
		t.Fatalf("Scaffold() with force error = %v", err)
	}

	data, err := os.ReadFile(existingPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "new" {
		t.Fatalf("file.txt = %q, want %q", string(data), "new")
	}
}

func TestScaffold_overwritesExistingAfterValidation(t *testing.T) {
	target := t.TempDir()
	existingPath := filepath.Join(target, "file.txt")
	testutil.WriteFile(t, existingPath, "old")

	templates := fstest.MapFS{
		"file.txt": &fstest.MapFile{Data: []byte("new")},
	}

	// Without force, scaffold still overwrites since validation
	// guarantees no conflicts. The force param is for explicit override.
	if err := Scaffold(templates, target, "myapp", false); err != nil {
		t.Fatalf("Scaffold() error = %v", err)
	}

	data, err := os.ReadFile(existingPath)
	if err != nil {
		t.Fatal(err)
	}
	// Without force we still write (validation has already cleared conflicts)
	if string(data) != "new" {
		t.Fatalf("file.txt = %q, want %q", string(data), "new")
	}
}

func TestScaffold_createsReleaseSkeleton(t *testing.T) {
	target := t.TempDir()
	templates := fstest.MapFS{
		".savepoint/releases/v1/epics":      &fstest.MapFile{Mode: fs.ModeDir | 0755},
		".savepoint/releases/v1/v1-PRD.md":  &fstest.MapFile{Data: []byte("# v{{RELEASE_NUMBER}} PRD for {{PROJECT_NAME}}")},
	}

	if err := Scaffold(templates, target, "myapp", false); err != nil {
		t.Fatalf("Scaffold() error = %v", err)
	}

	epicsPath := filepath.Join(target, ".savepoint", "releases", "v1", "epics")
	if info, err := os.Stat(epicsPath); err != nil || !info.IsDir() {
		t.Errorf(".savepoint/releases/v1/epics not created as directory: %v", err)
	}

	prdPath := filepath.Join(target, ".savepoint", "releases", "v1", "v1-PRD.md")
	data, err := os.ReadFile(prdPath)
	if err != nil {
		t.Errorf(".savepoint/releases/v1/v1-PRD.md not created: %v", err)
	}
	if got := string(data); !strings.Contains(got, "v1 PRD for myapp") {
		t.Errorf("v1-PRD.md = %q, want interpolated content", got)
	}
}

func TestProjectNameFromDir(t *testing.T) {
	dir := t.TempDir()
	name := filepath.Base(dir)
	got := ProjectNameFromDir(dir)
	if got != name {
		t.Fatalf("ProjectNameFromDir(%q) = %q, want %q", dir, got, name)
	}
}

func TestProjectNameFromDir_dot(t *testing.T) {
	got := ProjectNameFromDir(".")
	cwd, _ := os.Getwd()
	want := filepath.Base(cwd)
	if got != want {
		t.Fatalf("ProjectNameFromDir(\".\") = %q, want %q", got, want)
	}
}

func TestInterpolate(t *testing.T) {
	tests := []struct {
		input string
		name  string
		want  string
	}{
		{input: "# {{PROJECT_NAME}}", name: "myapp", want: "# myapp"},
		{input: "v{{RELEASE_NUMBER}}", name: "myapp", want: "v1"},
		{input: "{{PROJECT_NAME}} v{{RELEASE_NUMBER}}", name: "foo", want: "foo v1"},
		{input: "no variables", name: "myapp", want: "no variables"},
	}

	for _, tt := range tests {
		got := interpolate(tt.input, tt.name)
		if got != tt.want {
			t.Errorf("interpolate(%q, %q) = %q, want %q", tt.input, tt.name, got, tt.want)
		}
	}
}
