package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

// WriteFile writes content to path, creating parent directories if needed.
func WriteFile(t testing.TB, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

// MkdirAll creates directories with mode 0755, fatal on error.
func MkdirAll(t testing.TB, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatal(err)
	}
}

// SkipIfCaseInsensitive skips a test whose fixture needs two names that
// differ only by case. On Windows and default macOS filesystems the second
// name opens the first file, so the fixture cannot exist there.
func SkipIfCaseInsensitive(t testing.TB, dir string) {
	t.Helper()
	probe := filepath.Join(dir, "Case-Probe")
	if err := os.WriteFile(probe, nil, 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(probe)
	if _, err := os.Stat(filepath.Join(dir, "case-probe")); err == nil {
		t.Skip("filesystem is case-insensitive; names differing only by case are one file")
	}
}
