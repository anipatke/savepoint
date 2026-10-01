package cmd

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
)

func TestParseHealthArgs(t *testing.T) {
	cases := []struct {
		args []string
		want HealthSetupOptions
	}{
		{[]string{"setup"}, HealthSetupOptions{Dir: "."}},
		{[]string{"setup", "proj"}, HealthSetupOptions{Dir: "proj"}},
		{[]string{"setup", "--apply", "proj"}, HealthSetupOptions{Dir: "proj", Apply: true}},
	}
	for _, c := range cases {
		got, help, err := ParseHealthArgs(c.args)
		if err != nil || help || got != c.want {
			t.Errorf("ParseHealthArgs(%v) = %+v, %t, %v; want %+v", c.args, got, help, err, c.want)
		}
	}
}

func TestParseHealthArgsRejectsBadInput(t *testing.T) {
	for _, args := range [][]string{
		nil,
		{"refresh"},
		{"setup", "--bogus"},
		{"setup", "one", "two"},
	} {
		if _, help, err := ParseHealthArgs(args); err == nil || help {
			t.Errorf("ParseHealthArgs(%v) = help %t, err %v; want an error", args, help, err)
		}
	}
}

func TestRunHealthHelpDoesNotRun(t *testing.T) {
	for _, args := range [][]string{{"--help"}, {"setup", "--help"}} {
		var stdout bytes.Buffer
		called := false
		err := RunHealth(context.Background(), args, &stdout, func(context.Context, HealthSetupOptions) error {
			called = true
			return nil
		})
		if err != nil || called || !strings.Contains(stdout.String(), "Usage: health setup [dir] [--apply]") {
			t.Errorf("RunHealth(%v): err %v, called %t, stdout %q", args, err, called, stdout.String())
		}
	}
}

func TestRunHealthReturnsRunnerError(t *testing.T) {
	want := errors.New("runner failed")
	err := RunHealth(context.Background(), []string{"setup"}, &bytes.Buffer{}, func(context.Context, HealthSetupOptions) error {
		return want
	})
	if !errors.Is(err, want) {
		t.Fatalf("RunHealth() error = %v, want %v", err, want)
	}
}
