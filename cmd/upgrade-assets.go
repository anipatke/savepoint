package cmd

import (
	"context"
	"fmt"
	"io"
)

const upgradeUsage = "Usage: upgrade-assets [dir] [--dry-run] [--force]"

type UpgradeAssetsOptions struct {
	Dir    string
	DryRun bool
	Force  bool
}

type UpgradeAssetsRunner func(context.Context, UpgradeAssetsOptions) error

func RunUpgradeAssets(ctx context.Context, args []string, stdout io.Writer, runner UpgradeAssetsRunner) error {
	options, help, err := ParseUpgradeAssetsArgs(args)
	if help {
		_, writeErr := fmt.Fprintln(stdout, upgradeUsage)
		return writeErr
	}
	if err != nil {
		return err
	}
	return runner(ctx, options)
}

func ParseUpgradeAssetsArgs(args []string) (UpgradeAssetsOptions, bool, error) {
	options := UpgradeAssetsOptions{Dir: "."}
	var dirSet bool

	for _, arg := range args {
		switch arg {
		case "--help":
			return UpgradeAssetsOptions{}, true, nil
		case "--dry-run":
			options.DryRun = true
		case "--force":
			options.Force = true
		default:
			if len(arg) > 0 && arg[0] == '-' {
				return UpgradeAssetsOptions{}, false, fmt.Errorf("unknown upgrade-assets flag %q", arg)
			}
			if dirSet {
				return UpgradeAssetsOptions{}, false, fmt.Errorf("upgrade-assets accepts at most one directory")
			}
			options.Dir = arg
			dirSet = true
		}
	}

	return options, false, nil
}
