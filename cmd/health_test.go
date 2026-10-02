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
		invocation, help, err := ParseHealthArgs(c.args)
		var got HealthSetupOptions
		if invocation.Setup != nil {
			got = *invocation.Setup
		}
		if err != nil || help || invocation.Setup == nil || got != c.want {
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
		{"check"},
		{"check", "T-001"},
		{"check", "O-1"},
		{"check", "O-001", "one", "two"},
		{"check", "O-001", "--manual"},
		{"check", "--bogus", "O-001"},
		{"report", "--bogus"},
		{"report", "one", "two"},
		{"report", "O-001", "proj"},
	} {
		if _, help, err := ParseHealthArgs(args); err == nil || help {
			t.Errorf("ParseHealthArgs(%v) = help %t, err %v; want an error", args, help, err)
		}
	}
}

func TestParseHealthCheckArgs(t *testing.T) {
	cases := []struct {
		args []string
		want HealthCheckOptions
	}{
		{[]string{"check", "O-030"}, HealthCheckOptions{Objective: "O-030", Dir: "."}},
		{[]string{"check", "O-030", "proj"}, HealthCheckOptions{Objective: "O-030", Dir: "proj"}},
	}
	for _, c := range cases {
		invocation, help, err := ParseHealthArgs(c.args)
		if err != nil || help || invocation.Check == nil || *invocation.Check != c.want {
			t.Errorf("ParseHealthArgs(%v) = %+v, %t, %v; want %+v", c.args, invocation, help, err, c.want)
		}
	}
}

func TestRunHealthHelpDoesNotRun(t *testing.T) {
	for _, args := range [][]string{{"--help"}, {"setup", "--help"}, {"check", "--help"}} {
		var stdout bytes.Buffer
		called := false
		runners := HealthRunners{
			Setup: func(context.Context, HealthSetupOptions) error { called = true; return nil },
			Check: func(context.Context, HealthCheckOptions) error { called = true; return nil },
		}
		err := RunHealth(context.Background(), args, &stdout, runners)
		out := stdout.String()
		if err != nil || called || !strings.Contains(out, "health setup [dir] [--apply]") || !strings.Contains(out, "health check O-### [dir]") {
			t.Errorf("RunHealth(%v): err %v, called %t, stdout %q", args, err, called, stdout.String())
		}
	}
}

func TestRunHealthReturnsRunnerError(t *testing.T) {
	want := errors.New("runner failed")
	runners := HealthRunners{
		Setup: func(context.Context, HealthSetupOptions) error { return want },
		Check: func(context.Context, HealthCheckOptions) error { return want },
	}
	for _, args := range [][]string{{"setup"}, {"check", "O-001"}} {
		if err := RunHealth(context.Background(), args, &bytes.Buffer{}, runners); !errors.Is(err, want) {
			t.Fatalf("RunHealth(%v) error = %v, want %v", args, err, want)
		}
	}
}

func TestParseHealthReportArgs(t *testing.T) {
	cases := []struct {
		args []string
		want HealthReportOptions
	}{
		{[]string{"report"}, HealthReportOptions{Dir: "."}},
		{[]string{"report", "proj"}, HealthReportOptions{Dir: "proj"}},
	}
	for _, c := range cases {
		invocation, help, err := ParseHealthArgs(c.args)
		if err != nil || help || invocation.Report == nil || *invocation.Report != c.want {
			t.Errorf("ParseHealthArgs(%v) = %+v, %t, %v; want %+v", c.args, invocation, help, err, c.want)
		}
	}
}

func TestParseHealthReportRejectionShowsUsage(t *testing.T) {
	for _, args := range [][]string{{"report", "--bogus"}, {"report", "one", "two"}} {
		_, _, err := ParseHealthArgs(args)
		if err == nil || !strings.Contains(err.Error(), "health report [dir]") {
			t.Errorf("ParseHealthArgs(%v) error = %v, want the usage line", args, err)
		}
	}
}

func TestRunHealthDispatchesReport(t *testing.T) {
	var got HealthReportOptions
	runners := HealthRunners{Report: func(_ context.Context, o HealthReportOptions) error { got = o; return nil }}
	if err := RunHealth(context.Background(), []string{"report", "proj"}, &bytes.Buffer{}, runners); err != nil || got.Dir != "proj" {
		t.Fatalf("RunHealth(report proj) = %v, options %+v", err, got)
	}
	var stdout bytes.Buffer
	if err := RunHealth(context.Background(), []string{"report", "--help"}, &stdout, HealthRunners{}); err != nil || !strings.Contains(stdout.String(), "health report [dir]") {
		t.Fatalf("report --help = %q, %v", stdout.String(), err)
	}
}
