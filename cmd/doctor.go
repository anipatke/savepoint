package cmd

import (
	"context"
	"fmt"
	"io"
)

const doctorUsage = "Usage: doctor"

type DoctorOptions struct{}

// DoctorRunner receives parsed options and returns an exit code: 0=clean, 1=problems, 2=internal error.
type DoctorRunner func(DoctorOptions) (int, error)

func RunDoctor(ctx context.Context, args []string, stdout io.Writer, runner DoctorRunner) (int, error) {
	options, help, err := ParseDoctorArgs(args)
	if help {
		_, writeErr := fmt.Fprintln(stdout, doctorUsage)
		return 0, writeErr
	}
	if err != nil {
		return 2, err
	}
	return runner(options)
}

func ParseDoctorArgs(args []string) (DoctorOptions, bool, error) {
	var options DoctorOptions

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--help":
			return DoctorOptions{}, true, nil
		default:
			if len(arg) > 0 && arg[0] == '-' {
				return DoctorOptions{}, false, fmt.Errorf("unknown doctor flag %q", arg)
			}
			return DoctorOptions{}, false, fmt.Errorf("doctor takes no positional arguments, got %q", arg)
		}
	}

	return options, false, nil
}
