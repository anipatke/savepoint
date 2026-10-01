package cmd

import (
	"context"
	"fmt"
	"io"
)

const healthUsage = "Usage: health setup [dir] [--apply]"

// HealthSetupOptions is setup's whole parsed surface. Without Apply it only
// previews.
type HealthSetupOptions struct {
	Dir   string
	Apply bool
}

type HealthSetupRunner func(context.Context, HealthSetupOptions) error

func RunHealth(ctx context.Context, args []string, stdout io.Writer, runner HealthSetupRunner) error {
	options, help, err := ParseHealthArgs(args)
	if help {
		_, writeErr := fmt.Fprintln(stdout, healthUsage)
		return writeErr
	}
	if err != nil {
		return err
	}
	return runner(ctx, options)
}

func ParseHealthArgs(args []string) (HealthSetupOptions, bool, error) {
	options := HealthSetupOptions{Dir: "."}
	if len(args) > 0 && args[0] == "--help" {
		return options, true, nil
	}
	if len(args) == 0 || args[0] != "setup" {
		return options, false, fmt.Errorf("health needs a subcommand: setup\n%s", healthUsage)
	}
	var dirSet bool
	for _, arg := range args[1:] {
		switch arg {
		case "--help":
			return HealthSetupOptions{}, true, nil
		case "--apply":
			options.Apply = true
		default:
			if len(arg) > 0 && arg[0] == '-' {
				return HealthSetupOptions{}, false, fmt.Errorf("unknown health setup flag %q", arg)
			}
			if dirSet {
				return HealthSetupOptions{}, false, fmt.Errorf("health setup accepts at most one directory")
			}
			options.Dir = arg
			dirSet = true
		}
	}
	return options, false, nil
}
