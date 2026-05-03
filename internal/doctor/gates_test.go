package doctor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeConfig(t *testing.T, root, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, "config.yml"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestRunQualityGates_NoConfig(t *testing.T) {
	root := t.TempDir()
	results := RunQualityGates(root)
	if len(results) != 1 {
		t.Fatalf("RunQualityGates() = %d results, want 1", len(results))
	}
	if !results[0].Passed {
		t.Fatalf("RunQualityGates() = %v, want passed (default config has no gates)", results[0])
	}
	if !strings.Contains(results[0].Output, "no quality gates configured") {
		t.Fatalf("RunQualityGates() output = %q, want 'no quality gates configured'", results[0].Output)
	}
}

func TestRunQualityGates_AllNull(t *testing.T) {
	root := t.TempDir()
	writeConfig(t, root, "quality_gates:\n  lint: null\n  typecheck: null\n  test: null\ntheme:\n  bg: \"#000\"\n")

	results := RunQualityGates(root)
	if len(results) != 1 {
		t.Fatalf("RunQualityGates() = %d results, want 1 (no gates configured)", len(results))
	}
	if !results[0].Passed {
		t.Fatalf("RunQualityGates() = %v, want passed with 'no gates configured'", results[0])
	}
	if !strings.Contains(results[0].Output, "no quality gates configured") {
		t.Fatalf("RunQualityGates() output = %q, want 'no quality gates configured'", results[0].Output)
	}
}

func TestRunQualityGates_LintOnly(t *testing.T) {
	root := t.TempDir()
	writeConfig(t, root, "quality_gates:\n  lint: \"go version\"\n  typecheck: null\n  test: null\ntheme:\n  bg: \"#000\"\n")

	results := RunQualityGates(root)
	if len(results) != 1 {
		t.Fatalf("RunQualityGates() = %d results, want 1 (lint only)", len(results))
	}
	if !results[0].Passed {
		t.Fatalf("RunQualityGates() lint should pass: %v", results[0])
	}
	if results[0].Name != "lint" {
		t.Fatalf("RunQualityGates()[0].Name = %q, want \"lint\"", results[0].Name)
	}
}

func TestRunQualityGates_AllThree(t *testing.T) {
	root := t.TempDir()
	writeConfig(t, root, "quality_gates:\n  lint: \"go version\"\n  typecheck: \"go version\"\n  test: \"go version\"\ntheme:\n  bg: \"#000\"\n")

	results := RunQualityGates(root)
	if len(results) != 3 {
		t.Fatalf("RunQualityGates() = %d results, want 3", len(results))
	}
	for _, r := range results {
		if !r.Passed {
			t.Fatalf("RunQualityGates() %s should pass: %v", r.Name, r)
		}
	}
}

func TestRunQualityGates_FailingCommand(t *testing.T) {
	root := t.TempDir()
	writeConfig(t, root, "quality_gates:\n  lint: \"cmd-that-does-not-exist-12345\"\n  typecheck: null\n  test: null\ntheme:\n  bg: \"#000\"\n")

	results := RunQualityGates(root)
	if len(results) != 1 {
		t.Fatalf("RunQualityGates() = %d results, want 1", len(results))
	}
	if results[0].Passed {
		t.Fatal("RunQualityGates() should fail for non-existent command")
	}
}

func TestRunQualityGates_ExitCodeNonZero(t *testing.T) {
	root := t.TempDir()
	writeConfig(t, root, "quality_gates:\n  lint: \"go vet\"\n  typecheck: null\n  test: null\ntheme:\n  bg: \"#000\"\n")
	// Write a bad Go file so go vet fails
	badPath := filepath.Join(root, "bad.go")
	os.WriteFile(badPath, []byte("package x\n\nfunc f() { return 1 }\n"), 0644)

	results := RunQualityGates(root)
	if len(results) != 1 {
		t.Fatalf("RunQualityGates() = %d results, want 1", len(results))
	}
	if results[0].Passed {
		t.Fatal("RunQualityGates() should fail when go vet fails")
	}
}

func TestRunQualityGates_Timeout(t *testing.T) {
	root := t.TempDir()
	writeConfig(t, root, "quality_gates:\n  lint: \"go test -test.run TestNonexistent -count=1 ./...\"\n  gate_timeout: \"1ns\"\ntheme:\n  bg: \"#000\"\n")

	results := RunQualityGates(root)
	if len(results) != 1 {
		t.Fatalf("RunQualityGates() = %d results, want 1", len(results))
	}
	if results[0].Passed {
		t.Fatal("RunQualityGates() should fail due to timeout")
	}
	if !strings.Contains(results[0].Output, "timed out") {
		t.Fatalf("RunQualityGates() output = %q, want 'timed out'", results[0].Output)
	}
}

func TestRunQualityGates_DefaultTimeout(t *testing.T) {
	root := t.TempDir()
	writeConfig(t, root, "quality_gates:\n  lint: \"go version\"\n  # no gate_timeout set — uses 60s default\ntheme:\n  bg: \"#000\"\n")

	results := RunQualityGates(root)
	if len(results) != 1 {
		t.Fatalf("RunQualityGates() = %d results, want 1", len(results))
	}
	if !results[0].Passed {
		t.Fatalf("RunQualityGates() lint should pass with default timeout: %v", results[0])
	}
}

func TestSplitCommand(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"echo hello", []string{"echo", "hello"}},
		{"go test ./...", []string{"go", "test", "./..."}},
		{"\"c:\\program files\\go\\bin\\go\"", []string{"c:\\program files\\go\\bin\\go"}},
		{"", nil},
		{"   ", nil},
	}
	for _, tt := range tests {
		got := splitCommand(tt.input)
		if len(got) != len(tt.want) {
			t.Errorf("splitCommand(%q) = %v, want %v", tt.input, got, tt.want)
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("splitCommand(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
			}
		}
	}
}
