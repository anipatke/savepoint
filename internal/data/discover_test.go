package data

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindSavepointRoot(t *testing.T) {
	d := NewDiscover()
	savepointRoot := createDiscoveryFixture(t)
	start := filepath.Join(filepath.Dir(savepointRoot), "nested", "child")
	if err := os.MkdirAll(start, 0755); err != nil {
		t.Fatal(err)
	}

	root, err := d.FindSavepointRoot(start)
	if err != nil {
		t.Fatalf("FindSavepointRoot() error = %v", err)
	}
	if root != savepointRoot {
		t.Errorf("FindSavepointRoot() = %v, want %v", root, savepointRoot)
	}
}

func TestListReleases(t *testing.T) {
	d := NewDiscover()
	root := createDiscoveryFixture(t)

	releases, err := d.ListReleases(root)
	if err != nil {
		t.Fatalf("ListReleases() error = %v", err)
	}

	if len(releases) != 2 {
		t.Fatalf("ListReleases() returned %d releases, want 2", len(releases))
	}
	if releases[0].ID != "v1" || releases[1].ID != "v2" {
		t.Errorf("ListReleases() IDs = %v, want [v1 v2]", []string{releases[0].ID, releases[1].ID})
	}
}

func TestListRootDirs(t *testing.T) {
	d := NewDiscover()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "beta"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "alpha"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "notes.txt"), []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	dirs, err := d.ListRootDirs(root)
	if err != nil {
		t.Fatalf("ListRootDirs() error = %v", err)
	}

	if len(dirs) != 2 || dirs[0] != "alpha" || dirs[1] != "beta" {
		t.Fatalf("ListRootDirs() = %v, want [alpha beta]", dirs)
	}
}

func TestListRootDirsRejectsFile(t *testing.T) {
	d := NewDiscover()
	root := t.TempDir()
	path := filepath.Join(root, "not-dir")
	if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := d.ListRootDirs(path)
	if err == nil {
		t.Fatal("ListRootDirs() error = nil, want not directory error")
	}
}

func TestListEpics(t *testing.T) {
	d := NewDiscover()
	root := createDiscoveryFixture(t)

	epics, err := d.ListEpics(root, "v1")
	if err != nil {
		t.Fatalf("ListEpics() error = %v", err)
	}

	if len(epics) != 2 {
		t.Fatalf("ListEpics() returned %d epics, want 2", len(epics))
	}
	if epics[0].ID != "E01-go-setup" || epics[1].ID != "E02-data-readers" {
		t.Errorf("ListEpics() IDs = %v, want [E01-go-setup E02-data-readers]", []string{epics[0].ID, epics[1].ID})
	}
}

func TestListTasks(t *testing.T) {
	d := NewDiscover()
	root := createDiscoveryFixture(t)

	tasks, err := d.ListTasks(root, "v1", "E02-data-readers")
	if err != nil {
		t.Fatalf("ListTasks() error = %v", err)
	}

	if len(tasks) != 2 {
		t.Fatalf("ListTasks() returned %d tasks, want 2", len(tasks))
	}
	if tasks[0].ID != "T001-task-struct" || tasks[1].ID != "T002-frontmatter-parser" {
		t.Errorf("ListTasks() IDs = %v, want [T001-task-struct T002-frontmatter-parser]", []string{tasks[0].ID, tasks[1].ID})
	}
}

func createDiscoveryFixture(t *testing.T) string {
	t.Helper()

	root := t.TempDir()
	savepointRoot := filepath.Join(root, ".savepoint")
	paths := []string{
		filepath.Join(savepointRoot, "releases", "v1", "epics", "E02-data-readers", "tasks"),
		filepath.Join(savepointRoot, "releases", "v1", "epics", "E01-go-setup", "tasks"),
		filepath.Join(savepointRoot, "releases", "v1", "epics", "_archived"),
		filepath.Join(savepointRoot, "releases", "v2", "epics"),
	}
	for _, path := range paths {
		if err := os.MkdirAll(path, 0755); err != nil {
			t.Fatal(err)
		}
	}

	files := []string{
		filepath.Join(savepointRoot, "releases", "v1", "epics", "E02-data-readers", "tasks", "T002-frontmatter-parser.md"),
		filepath.Join(savepointRoot, "releases", "v1", "epics", "E02-data-readers", "tasks", "T001-task-struct.md"),
		filepath.Join(savepointRoot, "releases", "v1", "epics", "E02-data-readers", "tasks", "notes.txt"),
	}
	for _, file := range files {
		if err := os.WriteFile(file, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	return savepointRoot
}
