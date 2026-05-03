package init

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAtomicWrite_createsFile(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "output.txt")

	if err := AtomicWrite(target, []byte("hello")); err != nil {
		t.Fatalf("AtomicWrite() error = %v", err)
	}

	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello" {
		t.Fatalf("got %q, want %q", string(data), "hello")
	}
}

func TestAtomicWrite_replacesExistingFile(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "output.txt")
	if err := os.WriteFile(target, []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := AtomicWrite(target, []byte("new")); err != nil {
		t.Fatalf("AtomicWrite() error = %v", err)
	}

	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "new" {
		t.Fatalf("got %q, want %q", string(data), "new")
	}

	assertNoWriteArtifacts(t, dir)
}

func TestAtomicWrite_noTempFileLeftBehind(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "output.txt")

	if err := AtomicWrite(target, []byte("data")); err != nil {
		t.Fatalf("AtomicWrite() error = %v", err)
	}

	assertNoWriteArtifacts(t, dir)
}

func TestAtomicWrite_handlesNestedDirectories(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "deep", "nested", "output.txt")

	// Parent directories must exist before calling AtomicWrite
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		t.Fatal(err)
	}

	if err := AtomicWrite(target, []byte("nested")); err != nil {
		t.Fatalf("AtomicWrite() error = %v", err)
	}

	if _, err := os.Stat(target); err != nil {
		t.Errorf("target not created: %v", err)
	}
}

func assertNoWriteArtifacts(t *testing.T, dir string) {
	t.Helper()

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".tmp-") && strings.HasSuffix(e.Name(), ".write") {
			t.Errorf("temp file %q left behind after successful write", e.Name())
		}
		if strings.HasSuffix(e.Name(), ".savepoint-bak") {
			t.Errorf("backup file %q left behind after successful write", e.Name())
		}
	}
}
