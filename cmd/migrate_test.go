package cmd

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunMigrateHelp(t *testing.T) {
	var stdout bytes.Buffer
	called := false

	code, err := RunMigrate(context.Background(), []string{"--help"}, &stdout, func(context.Context, MigrateOptions) (int, error) {
		called = true
		return 0, nil
	})

	if err != nil {
		t.Fatalf("RunMigrate() error = %v", err)
	}
	if called {
		t.Fatal("RunMigrate() called runner for help")
	}
	if code != 0 {
		t.Fatalf("RunMigrate() code = %d, want 0", code)
	}
	if !strings.Contains(stdout.String(), "Usage: migrate [dir]") {
		t.Fatalf("help output = %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "Preview is the default") {
		t.Fatalf("help output = %q, want it to state preview is the default", stdout.String())
	}
}

func TestRunMigrateDefaults(t *testing.T) {
	got := runMigrateOptions(t, nil)

	if got.Dir != "." {
		t.Fatalf("Dir = %q, want .", got.Dir)
	}
	if got.Apply || got.DryRun || got.Recover {
		t.Fatalf("options = %+v, want every flag false by default", got)
	}
	if got.WillWrite() {
		t.Fatal("WillWrite() = true, want false by default (preview)")
	}
}

func TestRunMigrateUsesSpecifiedDirectory(t *testing.T) {
	got := runMigrateOptions(t, []string{"example"})
	if got.Dir != "example" {
		t.Fatalf("Dir = %q, want example", got.Dir)
	}
}

func TestRunMigrateParsesApply(t *testing.T) {
	got := runMigrateOptions(t, []string{"--apply"})
	if !got.Apply {
		t.Fatal("Apply = false, want true")
	}
	if !got.WillWrite() {
		t.Fatal("WillWrite() = false, want true for --apply alone")
	}
}

func TestRunMigrateParsesDryRun(t *testing.T) {
	got := runMigrateOptions(t, []string{"--dry-run"})
	if !got.DryRun {
		t.Fatal("DryRun = false, want true")
	}
	if got.WillWrite() {
		t.Fatal("WillWrite() = true, want false for --dry-run")
	}
}

func TestRunMigrateApplyAndDryRunTogetherPreviews(t *testing.T) {
	got := runMigrateOptions(t, []string{"--apply", "--dry-run"})
	if !got.Apply || !got.DryRun {
		t.Fatalf("options = %+v, want both flags set", got)
	}
	if got.WillWrite() {
		t.Fatal("WillWrite() = true, want false when --dry-run accompanies --apply")
	}
}

func TestRunMigrateParsesRecover(t *testing.T) {
	got := runMigrateOptions(t, []string{"--recover"})
	if !got.Recover {
		t.Fatal("Recover = false, want true")
	}
}

func TestRunMigrateParsesDecisionsFile(t *testing.T) {
	got := runMigrateOptions(t, []string{"--decisions", "decisions.yml"})
	if got.DecisionsFile != "decisions.yml" {
		t.Fatalf("DecisionsFile = %q, want decisions.yml", got.DecisionsFile)
	}
}

func TestRunMigrateDecisionsMissingValue(t *testing.T) {
	var stdout bytes.Buffer

	code, err := RunMigrate(context.Background(), []string{"--decisions"}, &stdout, func(context.Context, MigrateOptions) (int, error) {
		return 0, nil
	})

	if err == nil {
		t.Fatal("RunMigrate() error = nil, want missing value error")
	}
	if !strings.Contains(err.Error(), "--decisions requires a value") {
		t.Fatalf("error = %q", err.Error())
	}
	if code != 2 {
		t.Fatalf("code = %d, want 2", code)
	}
}

func TestRunMigrateRejectsUnknownFlag(t *testing.T) {
	var stdout bytes.Buffer
	called := false

	code, err := RunMigrate(context.Background(), []string{"--bogus"}, &stdout, func(context.Context, MigrateOptions) (int, error) {
		called = true
		return 0, nil
	})

	if err == nil {
		t.Fatal("RunMigrate() error = nil, want unknown flag error")
	}
	if called {
		t.Fatal("RunMigrate() called runner after invalid args")
	}
	if !strings.Contains(err.Error(), "unknown migrate flag") {
		t.Fatalf("error = %q, want unknown flag", err.Error())
	}
	if code != 2 {
		t.Fatalf("code = %d, want 2", code)
	}
	if !strings.Contains(stdout.String(), "Usage: migrate [dir]") {
		t.Fatalf("stdout = %q, want usage text printed alongside the error", stdout.String())
	}
}

func TestRunMigrateRejectsMultipleDirectories(t *testing.T) {
	code, err := RunMigrate(context.Background(), []string{"dir1", "dir2"}, &bytes.Buffer{}, func(context.Context, MigrateOptions) (int, error) {
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

func TestRunMigrateReturnsRunnerCodeAndError(t *testing.T) {
	var stdout bytes.Buffer

	code, err := RunMigrate(context.Background(), nil, &stdout, func(context.Context, MigrateOptions) (int, error) {
		return 1, nil
	})
	if err != nil {
		t.Fatalf("RunMigrate() error = %v", err)
	}
	if code != 1 {
		t.Fatalf("code = %d, want 1", code)
	}
}

func runMigrateOptions(t *testing.T, args []string) MigrateOptions {
	t.Helper()

	var stdout bytes.Buffer
	var got MigrateOptions
	code, err := RunMigrate(context.Background(), args, &stdout, func(_ context.Context, options MigrateOptions) (int, error) {
		got = options
		return 0, nil
	})
	if err != nil {
		t.Fatalf("RunMigrate() error = %v", err)
	}
	if code != 0 {
		t.Fatalf("RunMigrate() code = %d, want 0", code)
	}
	return got
}
