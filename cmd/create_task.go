package cmd

import (
	"context"
	"fmt"
	"io"
)

const createTaskUsage = `Usage: create-task --objective <O-###> --draft <path> [dir]

Draft must be a V2 Task Markdown record without an id. Its objective may be omitted or must match --objective; the command assigns the ID and preserves the remaining draft content.`

type CreateTaskOptions struct {
	Dir       string
	Objective string
	Draft     string
}

type CreateTaskRunner func(context.Context, CreateTaskOptions) (id, path string, err error)

func RunCreateTask(ctx context.Context, args []string, stdout, stderr io.Writer, runner CreateTaskRunner) error {
	options, help, err := ParseCreateTaskArgs(args)
	if help {
		_, writeErr := fmt.Fprintln(stdout, createTaskUsage)
		return writeErr
	}
	if err != nil {
		return err
	}
	id, path, err := runner(ctx, options)
	if err != nil {
		return err
	}
	if _, err = fmt.Fprintf(stdout, "Created %s at .savepoint/%s\n", id, path); err != nil {
		// The Task is already persisted, so the command succeeded; a nonzero exit would invite a retry that duplicates it.
		fmt.Fprintf(stderr, "warning: created %s at .savepoint/%s, but writing the result to stdout failed: %v\n", id, path, err)
	}
	return nil
}

func ParseCreateTaskArgs(args []string) (CreateTaskOptions, bool, error) {
	options := CreateTaskOptions{Dir: "."}
	var dirSet bool
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--help":
			return CreateTaskOptions{}, true, nil
		case "--objective", "--draft":
			if i+1 >= len(args) || args[i+1] == "" || args[i+1][0] == '-' {
				return CreateTaskOptions{}, false, fmt.Errorf("create-task flag %q requires a value", arg)
			}
			i++
			if arg == "--objective" {
				if options.Objective != "" {
					return CreateTaskOptions{}, false, fmt.Errorf("create-task accepts --objective only once")
				}
				options.Objective = args[i]
			} else {
				if options.Draft != "" {
					return CreateTaskOptions{}, false, fmt.Errorf("create-task accepts --draft only once")
				}
				options.Draft = args[i]
			}
		default:
			if len(arg) > 0 && arg[0] == '-' {
				return CreateTaskOptions{}, false, fmt.Errorf("unknown create-task flag %q", arg)
			}
			if dirSet {
				return CreateTaskOptions{}, false, fmt.Errorf("create-task accepts at most one project directory")
			}
			options.Dir = arg
			dirSet = true
		}
	}
	if options.Objective == "" {
		return CreateTaskOptions{}, false, fmt.Errorf("create-task requires --objective")
	}
	if options.Draft == "" {
		return CreateTaskOptions{}, false, fmt.Errorf("create-task requires --draft")
	}
	return options, false, nil
}
