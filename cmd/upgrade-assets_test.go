package cmd

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
)

func TestRunUpgradeAssetsHelp(t *testing.T) {
	var stdout bytes.Buffer
	called := false

	err := RunUpgradeAssets(context.Background(), []string{"--help"}, &stdout, func(context.Context, UpgradeAssetsOptions) error {
		called = true
		return nil
	})

	if err != nil {
		t.Fatalf("RunUpgradeAssets() error = %v", err)
	}
	if called {
		t.Fatal("RunUpgradeAssets() called runner for help")
	}
	if !strings.Contains(stdout.String(), "Usage: upgrade-assets [dir] [--dry-run] [--force]") {
		t.Fatalf("help output = %q", stdout.String())
	}
}

func TestRunUpgradeAssetsDefaultsToCurrentDirectory(t *testing.T) {
	got := runUpgradeAssetsOptions(t, nil)
	if got.Dir != "." {
		t.Fatalf("Dir = %q, want .", got.Dir)
	}
}

func TestRunUpgradeAssetsUsesSpecifiedDirectory(t *testing.T) {
	got := runUpgradeAssetsOptions(t, []string{"example"})
	if got.Dir != "example" {
		t.Fatalf("Dir = %q, want example", got.Dir)
	}
}

func TestRunUpgradeAssetsParsesDryRunAndForce(t *testing.T) {
	got := runUpgradeAssetsOptions(t, []string{"example", "--dry-run", "--force"})
	if !got.DryRun {
		t.Fatal("DryRun = false, want true")
	}
	if !got.Force {
		t.Fatal("Force = false, want true")
	}
}

func TestRunUpgradeAssetsRejectsUnknownFlags(t *testing.T) {
	var stdout bytes.Buffer
	called := false

	err := RunUpgradeAssets(context.Background(), []string{"--bogus"}, &stdout, func(context.Context, UpgradeAssetsOptions) error {
		called = true
		return nil
	})

	if err == nil {
		t.Fatal("RunUpgradeAssets() error = nil, want unknown flag error")
	}
	if called {
		t.Fatal("RunUpgradeAssets() called runner after invalid args")
	}
	if !strings.Contains(err.Error(), "unknown upgrade-assets flag") {
		t.Fatalf("error = %q, want unknown flag", err.Error())
	}
}

func TestRunUpgradeAssetsReturnsRunnerError(t *testing.T) {
	want := errors.New("runner failed")
	var stdout bytes.Buffer

	err := RunUpgradeAssets(context.Background(), nil, &stdout, func(context.Context, UpgradeAssetsOptions) error {
		return want
	})

	if !errors.Is(err, want) {
		t.Fatalf("RunUpgradeAssets() error = %v, want %v", err, want)
	}
}

func TestRunUpgradeAssetsRejectsMultipleDirectories(t *testing.T) {
	err := RunUpgradeAssets(context.Background(), []string{"dir1", "dir2"}, &bytes.Buffer{}, func(context.Context, UpgradeAssetsOptions) error {
		return nil
	})
	if err == nil {
		t.Fatal("expected error for multiple directories")
	}
	if !strings.Contains(err.Error(), "at most one directory") {
		t.Fatalf("error = %q, want 'at most one directory'", err.Error())
	}
}

func runUpgradeAssetsOptions(t *testing.T, args []string) UpgradeAssetsOptions {
	t.Helper()
	var stdout bytes.Buffer
	var got UpgradeAssetsOptions
	err := RunUpgradeAssets(context.Background(), args, &stdout, func(_ context.Context, options UpgradeAssetsOptions) error {
		got = options
		return nil
	})
	if err != nil {
		t.Fatalf("RunUpgradeAssets() error = %v", err)
	}
	return got
}
