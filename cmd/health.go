package cmd

import (
	"context"
	"fmt"
	"io"
)

const healthUsage = "Usage: health setup [dir] [--apply]\n       health check O-### [dir]\n       health report [dir]"

// HealthSetupOptions is setup's whole parsed surface. Without Apply it only
// previews.
type HealthSetupOptions struct {
	Dir   string
	Apply bool
}

// HealthCheckOptions is check's whole parsed surface. There is no flag that
// changes what kind of snapshot it collects.
type HealthCheckOptions struct {
	Objective string
	Dir       string
}

// HealthReportOptions is report's whole parsed surface. It takes no Objective.
type HealthReportOptions struct {
	Dir string
}

// HealthInvocation is one parsed health command; exactly one of Setup, Check
// and Report is set.
type HealthInvocation struct {
	Setup  *HealthSetupOptions
	Check  *HealthCheckOptions
	Report *HealthReportOptions
}

type HealthSetupRunner func(context.Context, HealthSetupOptions) error
type HealthCheckRunner func(context.Context, HealthCheckOptions) error
type HealthReportRunner func(context.Context, HealthReportOptions) error

// HealthRunners are the production behaviors RunHealth dispatches to.
type HealthRunners struct {
	Setup  HealthSetupRunner
	Check  HealthCheckRunner
	Report HealthReportRunner
}

// isObjectiveID reports whether s is O- followed by at least three digits.
func isObjectiveID(s string) bool {
	if len(s) < 5 || s[:2] != "O-" {
		return false
	}
	for _, c := range s[2:] {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func RunHealth(ctx context.Context, args []string, stdout io.Writer, runners HealthRunners) error {
	invocation, help, err := ParseHealthArgs(args)
	if help {
		_, writeErr := fmt.Fprintln(stdout, healthUsage)
		return writeErr
	}
	if err != nil {
		return err
	}
	if invocation.Check != nil {
		return runners.Check(ctx, *invocation.Check)
	}
	if invocation.Report != nil {
		return runners.Report(ctx, *invocation.Report)
	}
	return runners.Setup(ctx, *invocation.Setup)
}

func ParseHealthArgs(args []string) (HealthInvocation, bool, error) {
	if len(args) > 0 && args[0] == "--help" {
		return HealthInvocation{}, true, nil
	}
	if len(args) == 0 {
		return HealthInvocation{}, false, fmt.Errorf("health needs a subcommand: setup, check or report\n%s", healthUsage)
	}
	switch args[0] {
	case "setup":
		return parseHealthSetup(args[1:])
	case "check":
		return parseHealthCheck(args[1:])
	case "report":
		return parseHealthReport(args[1:])
	}
	return HealthInvocation{}, false, fmt.Errorf("health needs a subcommand: setup, check or report\n%s", healthUsage)
}

func parseHealthSetup(args []string) (HealthInvocation, bool, error) {
	options := HealthSetupOptions{Dir: "."}
	var dirSet bool
	for _, arg := range args {
		switch arg {
		case "--help":
			return HealthInvocation{}, true, nil
		case "--apply":
			options.Apply = true
		default:
			if len(arg) > 0 && arg[0] == '-' {
				return HealthInvocation{}, false, fmt.Errorf("unknown health setup flag %q", arg)
			}
			if dirSet {
				return HealthInvocation{}, false, fmt.Errorf("health setup accepts at most one directory")
			}
			options.Dir = arg
			dirSet = true
		}
	}
	return HealthInvocation{Setup: &options}, false, nil
}

func parseHealthCheck(args []string) (HealthInvocation, bool, error) {
	options := HealthCheckOptions{Dir: "."}
	var dirSet, objectiveSet bool
	for _, arg := range args {
		switch {
		case arg == "--help":
			return HealthInvocation{}, true, nil
		case len(arg) > 0 && arg[0] == '-':
			return HealthInvocation{}, false, fmt.Errorf("unknown health check flag %q", arg)
		case !objectiveSet:
			if !isObjectiveID(arg) {
				return HealthInvocation{}, false, fmt.Errorf("health check needs an Objective ID like O-001, got %q", arg)
			}
			options.Objective = arg
			objectiveSet = true
		case !dirSet:
			options.Dir = arg
			dirSet = true
		default:
			return HealthInvocation{}, false, fmt.Errorf("health check accepts one Objective and at most one directory")
		}
	}
	if !objectiveSet {
		return HealthInvocation{}, false, fmt.Errorf("health check needs an Objective ID like O-001\n%s", healthUsage)
	}
	return HealthInvocation{Check: &options}, false, nil
}

func parseHealthReport(args []string) (HealthInvocation, bool, error) {
	options := HealthReportOptions{Dir: "."}
	var dirSet bool
	for _, arg := range args {
		switch {
		case arg == "--help":
			return HealthInvocation{}, true, nil
		case len(arg) > 0 && arg[0] == '-':
			return HealthInvocation{}, false, fmt.Errorf("unknown health report flag %q\n%s", arg, healthUsage)
		case dirSet:
			return HealthInvocation{}, false, fmt.Errorf("health report accepts at most one directory\n%s", healthUsage)
		default:
			options.Dir = arg
			dirSet = true
		}
	}
	return HealthInvocation{Report: &options}, false, nil
}
