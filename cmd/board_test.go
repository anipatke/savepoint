package cmd

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
)

func TestRunBoardHelp(t *testing.T) {
	var stdout bytes.Buffer
	called := false

	err := RunBoard(context.Background(), []string{"--help"}, &stdout, func(BoardOptions) error {
		called = true
		return nil
	})

	if err != nil {
		t.Fatalf("RunBoard() error = %v", err)
	}
	if called {
		t.Fatal("RunBoard() called runner for help")
	}
	if !strings.Contains(stdout.String(), "board [--release <release>] [--epic <epic>] [--objective <objective>]") {
		t.Fatalf("help output = %q", stdout.String())
	}
}

func TestRunBoardNoArgs(t *testing.T) {
	got := runBoardOptions(t, nil)

	if got.Release != "" {
		t.Fatalf("Release = %q, want empty", got.Release)
	}
	if got.Epic != "" {
		t.Fatalf("Epic = %q, want empty", got.Epic)
	}
}

func TestRunBoardRelease(t *testing.T) {
	got := runBoardOptions(t, []string{"--release", "v1"})

	if got.Release != "v1" {
		t.Fatalf("Release = %q, want v1", got.Release)
	}
}

func TestRunBoardEpic(t *testing.T) {
	got := runBoardOptions(t, []string{"--epic", "E03"})

	if got.Epic != "E03" {
		t.Fatalf("Epic = %q, want E03", got.Epic)
	}
}

func TestRunBoardReleaseAndEpic(t *testing.T) {
	got := runBoardOptions(t, []string{"--release", "v1", "--epic", "E03"})

	if got.Release != "v1" {
		t.Fatalf("Release = %q, want v1", got.Release)
	}
	if got.Epic != "E03" {
		t.Fatalf("Epic = %q, want E03", got.Epic)
	}
}

func TestRunBoardObjective(t *testing.T) {
	got := runBoardOptions(t, []string{"--objective", "O009"})

	if got.Objective != "O009" {
		t.Fatalf("Objective = %q, want O009", got.Objective)
	}
	if got.Release != "" || got.Epic != "" {
		t.Fatalf("Release/Epic = %q/%q, want both empty", got.Release, got.Epic)
	}
}

// TestRunBoardObjectiveAlongsideV1Filters proves parsing accepts filters from
// both schemas and judges neither: which of them applies is decided behind the
// runner, where the project's schema_version is known (ARCH-01).
func TestRunBoardObjectiveAlongsideV1Filters(t *testing.T) {
	got := runBoardOptions(t, []string{"--release", "v1", "--objective", "O009"})

	if got.Release != "v1" || got.Objective != "O009" {
		t.Fatalf("options = %+v, want both filters parsed", got)
	}
}

func TestRunBoardObjectiveMissingValue(t *testing.T) {
	var stdout bytes.Buffer

	err := RunBoard(context.Background(), []string{"--objective"}, &stdout, func(BoardOptions) error {
		return nil
	})

	if err == nil {
		t.Fatal("RunBoard() error = nil, want missing value error")
	}
	if !strings.Contains(err.Error(), "--objective requires a value") {
		t.Fatalf("error = %q", err.Error())
	}
}

func TestRunBoardRejectsUnknownFlag(t *testing.T) {
	var stdout bytes.Buffer

	err := RunBoard(context.Background(), []string{"--bogus"}, &stdout, func(BoardOptions) error {
		return nil
	})

	if err == nil {
		t.Fatal("RunBoard() error = nil, want unknown flag error")
	}
	if !strings.Contains(err.Error(), "unknown board flag") {
		t.Fatalf("error = %q, want unknown flag", err.Error())
	}
}

func TestRunBoardRejectsPositionalArgs(t *testing.T) {
	var stdout bytes.Buffer

	err := RunBoard(context.Background(), []string{"extra"}, &stdout, func(BoardOptions) error {
		return nil
	})

	if err == nil {
		t.Fatal("RunBoard() error = nil, want positional arg error")
	}
}

func TestRunBoardReleaseMissingValue(t *testing.T) {
	var stdout bytes.Buffer

	err := RunBoard(context.Background(), []string{"--release"}, &stdout, func(BoardOptions) error {
		return nil
	})

	if err == nil {
		t.Fatal("RunBoard() error = nil, want missing value error")
	}
	if !strings.Contains(err.Error(), "--release requires a value") {
		t.Fatalf("error = %q", err.Error())
	}
}

func TestRunBoardReturnsRunnerError(t *testing.T) {
	want := errors.New("runner failed")
	var stdout bytes.Buffer

	err := RunBoard(context.Background(), nil, &stdout, func(BoardOptions) error {
		return want
	})

	if !errors.Is(err, want) {
		t.Fatalf("RunBoard() error = %v, want %v", err, want)
	}
}

func runBoardOptions(t *testing.T, args []string) BoardOptions {
	t.Helper()

	var stdout bytes.Buffer
	var got BoardOptions
	err := RunBoard(context.Background(), args, &stdout, func(options BoardOptions) error {
		got = options
		return nil
	})
	if err != nil {
		t.Fatalf("RunBoard() error = %v", err)
	}
	return got
}
