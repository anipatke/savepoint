package cmd

import (
	"context"
	"fmt"
	"io"
)

const boardUsage = "Usage: board [--objective <objective>]"

// BoardOptions is the V2 board's whole parsed filter surface. Legacy Release
// and Epic filters are intentionally absent from the live command contract;
// migration remains the only V1 reader.
type BoardOptions struct {
	Objective string
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
		case "--objective":
			i++
			if i >= len(args) {
				return BoardOptions{}, false, fmt.Errorf("--objective requires a value")
			}
			options.Objective = args[i]
		default:
			if len(arg) > 0 && arg[0] == '-' {
				return BoardOptions{}, false, fmt.Errorf("unknown board flag %q", arg)
			}
			return BoardOptions{}, false, fmt.Errorf("board takes no positional arguments, got %q", arg)
		}
	}

	return options, false, nil
}
