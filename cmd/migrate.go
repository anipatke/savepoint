package cmd

import (
	"context"
	"fmt"
	"io"
)

const migrateUsage = "Usage: migrate [dir] [--apply] [--dry-run] [--decisions FILE] [--verbose]\n" +
	"  Preview is the default: nothing is written unless --apply is given.\n" +
	"  --dry-run is an explicit synonym for the default preview behavior;\n" +
	"  passing both --apply and --dry-run previews.\n" +
	"  --verbose lists every planned record, archived file, and note instead\n" +
	"  of the summary."

type MigrateOptions struct {
	Dir           string
	Apply         bool
	DryRun        bool
	DecisionsFile string
	Verbose       bool
}

// WillWrite reports whether these options ask the runner to write:
// --apply without --dry-run. Passing both is accepted and previews.
func (o MigrateOptions) WillWrite() bool {
	return o.Apply && !o.DryRun
}

// MigrateRunner receives parsed options and returns an exit code: 0=clean,
// 1=named refusal, 2=argument or internal error.
type MigrateRunner func(context.Context, MigrateOptions) (int, error)

func RunMigrate(ctx context.Context, args []string, stdout io.Writer, runner MigrateRunner) (int, error) {
	options, help, err := ParseMigrateArgs(args)
	if help {
		_, writeErr := fmt.Fprintln(stdout, migrateUsage)
		return 0, writeErr
	}
	if err != nil {
		fmt.Fprintln(stdout, migrateUsage)
		return 2, err
	}
	return runner(ctx, options)
}

func ParseMigrateArgs(args []string) (MigrateOptions, bool, error) {
	options := MigrateOptions{Dir: "."}
	var dirSet bool

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--help":
			return MigrateOptions{}, true, nil
		case "--apply":
			options.Apply = true
		case "--dry-run":
			options.DryRun = true
		case "--verbose":
			options.Verbose = true
		case "--decisions":
			i++
			if i >= len(args) {
				return MigrateOptions{}, false, fmt.Errorf("--decisions requires a value")
			}
			options.DecisionsFile = args[i]
		default:
			if len(arg) > 0 && arg[0] == '-' {
				return MigrateOptions{}, false, fmt.Errorf("unknown migrate flag %q", arg)
			}
			if dirSet {
				return MigrateOptions{}, false, fmt.Errorf("migrate accepts at most one directory")
			}
			options.Dir = arg
			dirSet = true
		}
	}

	return options, false, nil
}
