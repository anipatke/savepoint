package cmd

import (
	"context"
	"fmt"
	"io"
)

const boardUsage = "Usage: board [--release <release>] [--epic <epic>]"

type BoardOptions struct {
	Release string
	Epic    string
}

type BoardRunner func(BoardOptions) error

func RunBoard(ctx context.Context, args []string, stdout io.Writer, runner BoardRunner) error {
	options, help, err := ParseBoardArgs(args)
	if help {
		_, writeErr := fmt.Fprintln(stdout, boardUsage)
		return writeErr
	}
	if err != nil {
		return err
	}
	return runner(options)
}

func ParseBoardArgs(args []string) (BoardOptions, bool, error) {
	var options BoardOptions

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--help":
			return BoardOptions{}, true, nil
		case "--release":
			i++
			if i >= len(args) {
				return BoardOptions{}, false, fmt.Errorf("--release requires a value")
			}
			options.Release = args[i]
		case "--epic":
			i++
			if i >= len(args) {
				return BoardOptions{}, false, fmt.Errorf("--epic requires a value")
			}
			options.Epic = args[i]
		default:
			if len(arg) > 0 && arg[0] == '-' {
				return BoardOptions{}, false, fmt.Errorf("unknown board flag %q", arg)
			}
			return BoardOptions{}, false, fmt.Errorf("board takes no positional arguments, got %q", arg)
		}
	}

	return options, false, nil
}
