package cmd

import (
	"bytes"
	"context"
	"errors"
	"io"
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

func TestRunCreateTaskHelpAndSuccessOutput(t *testing.T) {
	var stdout bytes.Buffer
	called := false
	err := RunCreateTask(context.Background(), []string{"--help"}, &stdout, io.Discard, func(context.Context, CreateTaskOptions) (string, string, error) {
		called = true
		return "", "", nil
	})
	if err != nil {
		t.Fatalf("RunCreateTask(--help) error = %v", err)
	}
	if called {
		t.Fatal("RunCreateTask(--help) called runner")
	}
	if !strings.Contains(stdout.String(), "Usage: create-task --objective <O-###> --draft <path> [dir]") {
		t.Fatalf("help output = %q", stdout.String())
	}

	stdout.Reset()
	err = RunCreateTask(context.Background(), []string{"--objective", "O-001", "--draft", "task.md", "/tmp/project"}, &stdout, io.Discard, func(_ context.Context, options CreateTaskOptions) (string, string, error) {
		if options != (CreateTaskOptions{Dir: "/tmp/project", Objective: "O-001", Draft: "task.md"}) {
			t.Errorf("runner options = %+v", options)
		}
		return "T-014", "objectives/O-001-first/tasks/T-014-review.md", nil
	})
	if err != nil {
		t.Fatalf("RunCreateTask() error = %v", err)
	}
	if got, want := stdout.String(), "Created T-014 at .savepoint/objectives/O-001-first/tasks/T-014-review.md\n"; got != want {
		t.Fatalf("success output = %q, want %q", got, want)
	}
}

func TestParseCreateTaskArgs(t *testing.T) {
	got, help, err := ParseCreateTaskArgs([]string{"--draft", "task.md", "--objective", "O-001"})
	if err != nil || help {
		t.Fatalf("ParseCreateTaskArgs() = (%+v, %t, %v)", got, help, err)
	}
	if got != (CreateTaskOptions{Dir: ".", Objective: "O-001", Draft: "task.md"}) {
		t.Fatalf("options = %+v", got)
	}
}

func TestParseCreateTaskArgsRejectsMissingAndUnknownArguments(t *testing.T) {
	for _, args := range [][]string{
		{"--draft", "task.md"},
		{"--objective", "O-001"},
		{"--objective", "O-001", "--draft", "task.md", "--other"},
		{"--objective", "O-001", "--draft", "task.md", "one", "two"},
	} {
		if _, help, err := ParseCreateTaskArgs(args); err == nil || help {
			t.Errorf("ParseCreateTaskArgs(%v) = (help=%t, err=%v), want an argument error", args, help, err)
		}
	}
}

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) { return 0, errors.New("no space left on device") }

func TestRunCreateTaskOutputFailureSucceedsWithStderrWarning(t *testing.T) {
	var stderr bytes.Buffer
	calls := 0
	err := RunCreateTask(context.Background(), []string{"--objective", "O-001", "--draft", "task.md"}, failingWriter{}, &stderr, func(context.Context, CreateTaskOptions) (string, string, error) {
		calls++
		return "T-014", "objectives/O-001-first/tasks/T-014-review.md", nil
	})
	if err != nil {
		t.Fatalf("RunCreateTask() error = %v, want nil after the Task was persisted", err)
	}
	for _, want := range []string{"created T-014", "T-014-review.md", "no space left on device"} {
		if !strings.Contains(stderr.String(), want) {
			t.Errorf("stderr %q missing %q", stderr.String(), want)
		}
	}
	if calls != 1 {
		t.Fatalf("runner calls = %d, want 1", calls)
	}
}
