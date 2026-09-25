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
	if !strings.Contains(stdout.String(), "board [--objective <objective>]") {
		t.Fatalf("help output = %q", stdout.String())
	}
}

func TestRunBoardNoArgs(t *testing.T) {
	got := runBoardOptions(t, nil)

	if got.Objective != "" {
		t.Fatalf("Objective = %q, want empty", got.Objective)
	}
}

func TestRunBoardObjective(t *testing.T) {
	got := runBoardOptions(t, []string{"--objective", "O009"})

	if got.Objective != "O009" {
		t.Fatalf("Objective = %q, want O009", got.Objective)
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

func TestRunBoardRejectsLegacyReleaseAndEpicFlags(t *testing.T) {
	for _, flag := range []string{"--release", "--epic"} {
		t.Run(flag, func(t *testing.T) {
			var stdout bytes.Buffer
			err := RunBoard(context.Background(), []string{flag, "value"}, &stdout, func(BoardOptions) error {
				return nil
			})
			if err == nil {
				t.Fatalf("RunBoard(%q) error = nil, want legacy flag rejected", flag)
			}
			if !strings.Contains(err.Error(), "unknown board flag") {
				t.Fatalf("error = %q, want unknown board flag", err.Error())
			}
		})
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
