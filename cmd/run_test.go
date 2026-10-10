package cmd

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestParseRunArgs(t *testing.T) {
	cases := []struct {
		name    string
		args    []string
		want    RunOptions
		wantErr string
	}{
		{"task only", []string{"--task", "T-130"}, RunOptions{Task: "T-130", Dir: "."}, ""},
		{"task then dir", []string{"--task", "T-1", "proj"}, RunOptions{Task: "T-1", Dir: "proj"}, ""},
		{"dir then task", []string{"proj", "--task", "T-1"}, RunOptions{Task: "T-1", Dir: "proj"}, ""},
		{"missing task", nil, RunOptions{}, "--task is required"},
		{"task without value", []string{"--task"}, RunOptions{}, "--task needs"},
		{"not a task id", []string{"--task", "O-001"}, RunOptions{}, "not a Task ID"},
		{"unknown flag", []string{"--task", "T-1", "--force"}, RunOptions{}, `unknown run flag "--force"`},
		{"two dirs", []string{"--task", "T-1", "a", "b"}, RunOptions{}, "at most one directory"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, _, err := ParseRunArgs(c.args)
			if c.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), c.wantErr) {
					t.Fatalf("error = %v, want %q", err, c.wantErr)
				}
				return
			}
			if err != nil || got != c.want {
				t.Fatalf("got %+v, %v; want %+v", got, err, c.want)
			}
		})
	}
}

func TestRunRun_helpAndArgumentErrorsNeverCallTheRunner(t *testing.T) {
	called := false
	runner := func(context.Context, RunOptions) (int, error) { called = true; return 0, nil }

	var out bytes.Buffer
	if code, err := RunRun(context.Background(), []string{"--help"}, &out, runner); code != 0 || err != nil || !strings.Contains(out.String(), runUsage) {
		t.Errorf("help: code = %d, err = %v, out = %q", code, err, out.String())
	}
	out.Reset()
	if code, err := RunRun(context.Background(), []string{"--bogus"}, &out, runner); code != 2 || err == nil {
		t.Errorf("bad flag: code = %d, err = %v", code, err)
	}
	if called {
		t.Error("runner was called")
	}
}

func TestRunRun_passesParsedOptionsAndReturnsTheRunnersCode(t *testing.T) {
	var got RunOptions
	code, err := RunRun(context.Background(), []string{"--task", "T-7", "x"}, &bytes.Buffer{}, func(_ context.Context, o RunOptions) (int, error) {
		got = o
		return 1, nil
	})
	if code != 1 || err != nil || got != (RunOptions{Task: "T-7", Dir: "x"}) {
		t.Errorf("code = %d, err = %v, options = %+v", code, err, got)
	}
}
