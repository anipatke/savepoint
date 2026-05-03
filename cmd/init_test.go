package cmd

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
)

func TestRunInitHelp(t *testing.T) {
	var stdout bytes.Buffer
	called := false

	err := RunInit(context.Background(), []string{"--help"}, &stdout, func(context.Context, InitOptions) error {
		called = true
		return nil
	})

	if err != nil {
		t.Fatalf("RunInit() error = %v", err)
	}
	if called {
		t.Fatal("RunInit() called runner for help")
	}
	if !strings.Contains(stdout.String(), "Usage: init [dir] [--force] [--install]") {
		t.Fatalf("help output = %q", stdout.String())
	}
}

func TestRunInitDefaultsToCurrentDirectory(t *testing.T) {
	got := runInitOptions(t, nil)

	if got.Dir != "." {
		t.Fatalf("Dir = %q, want .", got.Dir)
	}
}

func TestRunInitUsesSpecifiedDirectory(t *testing.T) {
	got := runInitOptions(t, []string{"example"})

	if got.Dir != "example" {
		t.Fatalf("Dir = %q, want example", got.Dir)
	}
}

func TestRunInitParsesForceAndInstall(t *testing.T) {
	got := runInitOptions(t, []string{"example", "--force", "--install"})

	if !got.Force {
		t.Fatal("Force = false, want true")
	}
	if !got.Install {
		t.Fatal("Install = false, want true")
	}
}

func TestRunInitRejectsUnknownFlags(t *testing.T) {
	var stdout bytes.Buffer
	called := false

	err := RunInit(context.Background(), []string{"--bogus"}, &stdout, func(context.Context, InitOptions) error {
		called = true
		return nil
	})

	if err == nil {
		t.Fatal("RunInit() error = nil, want unknown flag error")
	}
	if called {
		t.Fatal("RunInit() called runner after invalid args")
	}
	if !strings.Contains(err.Error(), "unknown init flag") {
		t.Fatalf("error = %q, want unknown flag", err.Error())
	}
}

func TestRunInitReturnsRunnerError(t *testing.T) {
	want := errors.New("runner failed")
	var stdout bytes.Buffer

	err := RunInit(context.Background(), nil, &stdout, func(context.Context, InitOptions) error {
		return want
	})

	if !errors.Is(err, want) {
		t.Fatalf("RunInit() error = %v, want %v", err, want)
	}
}

func runInitOptions(t *testing.T, args []string) InitOptions {
	t.Helper()

	var stdout bytes.Buffer
	var got InitOptions
	err := RunInit(context.Background(), args, &stdout, func(_ context.Context, options InitOptions) error {
		got = options
		return nil
	})
	if err != nil {
		t.Fatalf("RunInit() error = %v", err)
	}
	return got
}
