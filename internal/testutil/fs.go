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
