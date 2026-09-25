package cmd

import (
	"context"
	"fmt"
	"io"
)

const resumeUsage = "Usage: resume [dir]"

// ResumeOptions is resume's whole parsed surface: an optional project
// directory, defaulting to the current one. There is no other flag: resume
// has nothing to configure, only somewhere to look.
type ResumeOptions struct {
	Dir string
}

// ResumeRunner receives parsed options and returns an exit code: 0 clean
// (including an intact project whose router selection did not resolve — the
// project is fine even when the router line is not), 1 a named diagnostic,
// 2 an argument error.
type ResumeRunner func(context.Context, ResumeOptions) (int, error)

func RunResume(ctx context.Context, args []string, stdout io.Writer, runner ResumeRunner) (int, error) {
	options, help, err := ParseResumeArgs(args)
	if help {
		_, writeErr := fmt.Fprintln(stdout, resumeUsage)
		return 0, writeErr
	}
	if err != nil {
		fmt.Fprintln(stdout, resumeUsage)
		return 2, err
	}
	return runner(ctx, options)
}

func ParseResumeArgs(args []string) (ResumeOptions, bool, error) {
	options := ResumeOptions{Dir: "."}
	var dirSet bool

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--help":
			return ResumeOptions{}, true, nil
		default:
			if len(arg) > 0 && arg[0] == '-' {
				return ResumeOptions{}, false, fmt.Errorf("unknown resume flag %q", arg)
			}
			if dirSet {
				return ResumeOptions{}, false, fmt.Errorf("resume accepts at most one directory")
			}
			options.Dir = arg
			dirSet = true
		}
	}

	return options, false, nil
}
