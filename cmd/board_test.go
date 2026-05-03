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
	if !strings.Contains(stdout.String(), "board [--release <release>] [--epic <epic>]") {
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
