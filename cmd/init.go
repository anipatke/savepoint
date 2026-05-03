package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
)

const initUsage = "Usage: init [dir] [--force] [--install]"

type InitOptions struct {
	Dir     string
	Force   bool
	Install bool
}

type InitRunner func(context.Context, InitOptions) error

var ErrInitNotImplemented = errors.New("init scaffold is not implemented yet")

func RunInit(ctx context.Context, args []string, stdout io.Writer, runner InitRunner) error {
	options, help, err := ParseInitArgs(args)
	if help {
		_, writeErr := fmt.Fprintln(stdout, initUsage)
		return writeErr
	}
	if err != nil {
		return err
	}
	return runner(ctx, options)
}

func ParseInitArgs(args []string) (InitOptions, bool, error) {
	options := InitOptions{Dir: "."}
	var dirSet bool

	for _, arg := range args {
		switch arg {
		case "--help":
			return InitOptions{}, true, nil
		case "--force":
			options.Force = true
		case "--install":
			options.Install = true
		default:
			if len(arg) > 0 && arg[0] == '-' {
				return InitOptions{}, false, fmt.Errorf("unknown init flag %q", arg)
			}
			if dirSet {
				return InitOptions{}, false, fmt.Errorf("init accepts at most one directory")
			}
			options.Dir = arg
			dirSet = true
		}
	}

	return options, false, nil
}

func InitNotImplemented(context.Context, InitOptions) error {
	return ErrInitNotImplemented
}
