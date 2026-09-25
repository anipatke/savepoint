package cmd

import (
	"bytes"
	"context"
	"go/build"
	"strings"
	"testing"
)

func TestRunResumeHelp(t *testing.T) {
	var stdout bytes.Buffer
	called := false

	code, err := RunResume(context.Background(), []string{"--help"}, &stdout, func(context.Context, ResumeOptions) (int, error) {
		called = true
		return 0, nil
	})

	if err != nil {
		t.Fatalf("RunResume() error = %v", err)
	}
	if called {
		t.Fatal("RunResume() called runner for help")
	}
	if code != 0 {
		t.Fatalf("RunResume() code = %d, want 0", code)
	}
	if !strings.Contains(stdout.String(), "Usage: resume [dir]") {
		t.Fatalf("help output = %q", stdout.String())
	}
}

func TestRunResumeDefaultsDirectory(t *testing.T) {
	got := runResumeOptions(t, nil)
	if got.Dir != "." {
		t.Fatalf("Dir = %q, want .", got.Dir)
	}
}

func TestRunResumeUsesSpecifiedDirectory(t *testing.T) {
	got := runResumeOptions(t, []string{"example"})
	if got.Dir != "example" {
		t.Fatalf("Dir = %q, want example", got.Dir)
	}
}

func TestRunResumeRejectsUnknownFlag(t *testing.T) {
	var stdout bytes.Buffer
	called := false

	code, err := RunResume(context.Background(), []string{"--bogus"}, &stdout, func(context.Context, ResumeOptions) (int, error) {
		called = true
		return 0, nil
	})

	if err == nil {
		t.Fatal("RunResume() error = nil, want unknown flag error")
	}
	if called {
		t.Fatal("RunResume() called runner after invalid args")
	}
	if !strings.Contains(err.Error(), "unknown resume flag") {
		t.Fatalf("error = %q, want unknown flag", err.Error())
	}
	if code != 2 {
		t.Fatalf("code = %d, want 2", code)
	}
	if !strings.Contains(stdout.String(), "Usage: resume [dir]") {
		t.Fatalf("stdout = %q, want usage text printed alongside the error", stdout.String())
	}
}

func TestRunResumeRejectsMultipleDirectories(t *testing.T) {
	code, err := RunResume(context.Background(), []string{"dir1", "dir2"}, &bytes.Buffer{}, func(context.Context, ResumeOptions) (int, error) {
		return 0, nil
	})
	if err == nil {
		t.Fatal("expected error for multiple directories")
	}
	if !strings.Contains(err.Error(), "at most one directory") {
		t.Fatalf("error = %q, want 'at most one directory'", err.Error())
	}
	if code != 2 {
		t.Fatalf("code = %d, want 2", code)
	}
}

func TestRunResumeReturnsRunnerCodeAndError(t *testing.T) {
	code, err := RunResume(context.Background(), nil, &bytes.Buffer{}, func(context.Context, ResumeOptions) (int, error) {
		return 1, nil
	})
	if err != nil {
		t.Fatalf("RunResume() error = %v", err)
	}
	if code != 1 {
		t.Fatalf("code = %d, want 1", code)
	}
}

func runResumeOptions(t *testing.T, args []string) ResumeOptions {
	t.Helper()

	var stdout bytes.Buffer
	var got ResumeOptions
	code, err := RunResume(context.Background(), args, &stdout, func(_ context.Context, options ResumeOptions) (int, error) {
		got = options
		return 0, nil
	})
	if err != nil {
		t.Fatalf("RunResume() error = %v", err)
	}
	if code != 0 {
		t.Fatalf("RunResume() code = %d, want 0", code)
	}
	return got
}

// TestPackage_staysThinNoDomainImports proves the whole cmd package —
// resume.go included — stays argument parsing and dispatch only (ARCH-01):
// no record parsing, no gate reading, no rendering, no subprocess, and no
// network access anywhere in it.
func TestPackage_staysThinNoDomainImports(t *testing.T) {
	pkg, err := build.ImportDir(".", build.IgnoreVendor)
	if err != nil {
		t.Fatalf("scan package imports: %v", err)
	}
	allowed := map[string]bool{"context": true, "errors": true, "fmt": true, "io": true}
	for _, imported := range pkg.Imports {
		if !allowed[imported] {
			t.Errorf("cmd package imports %q, which ARCH-01 reserves for the internal packages behind the injected runner", imported)
		}
	}
}
