package init

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
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
	if string(data) != "release: v1" {
		t.Fatalf("AGENTS.md = %q, want %q", string(data), "release: v1")
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
	if err := os.WriteFile(existingPath, []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}

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
	if err := os.WriteFile(existingPath, []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}

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
