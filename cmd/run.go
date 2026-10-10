package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
)

const runUsage = "Usage: run --task <T-###> [dir]"

// RunOptions is run's whole parsed surface: the Task to run and an optional
// project directory, defaulting to the current one.
type RunOptions struct {
	Task string
	Dir  string
}

// RunRunner receives parsed options and returns an exit code: 0 the agent
// finished cleanly, 1 a project or run problem, 2 a refusal the owner can fix
// (the feature is off, or the arguments are wrong).
type RunRunner func(context.Context, RunOptions) (int, error)

func RunRun(ctx context.Context, args []string, stdout io.Writer, runner RunRunner) (int, error) {
	options, help, err := ParseRunArgs(args)
	if help {
		_, writeErr := fmt.Fprintln(stdout, runUsage)
		return 0, writeErr
	}
	if err != nil {
		fmt.Fprintln(stdout, runUsage)
		return 2, err
	}
	return runner(ctx, options)
}

func ParseRunArgs(args []string) (RunOptions, bool, error) {
	options := RunOptions{Dir: "."}
	var dirSet bool

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--help":
			return RunOptions{}, true, nil
		case "--task":
			if i+1 >= len(args) {
				return RunOptions{}, false, errors.New("run: --task needs a Task ID such as T-001")
			}
			i++
			options.Task = args[i]
		default:
			if len(arg) > 0 && arg[0] == '-' {
				return RunOptions{}, false, fmt.Errorf("unknown run flag %q", arg)
			}
			if dirSet {
				return RunOptions{}, false, errors.New("run accepts at most one directory")
			}
			options.Dir = arg
			dirSet = true
		}
	}

	if options.Task == "" {
		return RunOptions{}, false, errors.New("run: --task is required")
	}
	if len(options.Task) < 3 || options.Task[:2] != "T-" {
		return RunOptions{}, false, fmt.Errorf("run: %q is not a Task ID such as T-001", options.Task)
	}
	return options, false, nil
}
