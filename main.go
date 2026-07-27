package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"os"

	"github.com/opencode/savepoint/cmd"
	"github.com/opencode/savepoint/internal/board"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/doctor"
	savepointinit "github.com/opencode/savepoint/internal/init"
)

//go:embed templates/project
//go:embed all:templates/project/.savepoint
//go:embed templates/prompts
var projectTemplates embed.FS

var version = "dev"

func main() {
	args, debug := stripDebugFlag(os.Args[1:])
	if debug || os.Getenv("SAVEPOINT_DEBUG") != "" {
		board.SetDebug(true)
	}

	if len(args) > 0 {
		switch args[0] {
		case "--version":
			fmt.Println(version)
			os.Exit(0)
		case "init":
			if err := cmd.RunInit(context.Background(), args[1:], os.Stdout, initRunner); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			os.Exit(0)
		case "board":
			if err := cmd.RunBoard(context.Background(), args[1:], os.Stdout, func(opts cmd.BoardOptions) error {
				return board.RunWithFilters(opts.Release, opts.Epic)
			}); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			os.Exit(0)
		case "doctor":
			code, err := cmd.RunDoctor(context.Background(), args[1:], os.Stdout, func(opts cmd.DoctorOptions) (int, error) {
				return runDoctorChecks(opts)
			})
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			os.Exit(code)
		case "upgrade-assets":
			if err := cmd.RunUpgradeAssets(context.Background(), args[1:], os.Stdout, upgradeAssetsRunner); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			os.Exit(0)
		}
	}
	if err := board.Run(); err != nil {
		panic(err)
	}
}

// stripDebugFlag removes --debug from args and reports whether it was present.
func stripDebugFlag(args []string) ([]string, bool) {
	out := make([]string, 0, len(args))
	found := false
	for _, a := range args {
		if a == "--debug" {
			found = true
		} else {
			out = append(out, a)
		}
	}
	return out, found
}

func runDoctorChecks(opts cmd.DoctorOptions) (int, error) {
	discover := data.NewDiscover()
	root, err := discover.FindSavepointRoot(".")
	if err != nil {
		return 2, fmt.Errorf("savepoint root not found: %w", err)
	}

	report := doctor.RunAllChecks(root, opts.Epic)
	fmt.Fprint(os.Stdout, report.Format())

	if report.HasProblems() {
		return 1, nil
	}
	return 0, nil
}

func upgradeAssetsRunner(ctx context.Context, opts cmd.UpgradeAssetsOptions) error {
	sub, err := fs.Sub(projectTemplates, "templates/project")
	if err != nil {
		return fmt.Errorf("cannot load templates: %w", err)
	}

	report, err := savepointinit.UpgradeProjectAssets(sub, opts.Dir, opts.DryRun, opts.Force)
	if err != nil {
		// A failure part-way through still applied whatever came before it.
		// Print that work before the error so the user knows what changed.
		if report != nil && len(report.Actions) > 0 {
			fmt.Print(report.Format())
		}
		return err
	}

	fmt.Print(report.Format())
	return nil
}

func initRunner(ctx context.Context, opts cmd.InitOptions) error {
	if err := savepointinit.ValidateTarget(opts.Dir, opts.Force); err != nil {
		return err
	}

	sub, err := fs.Sub(projectTemplates, "templates/project")
	if err != nil {
		return fmt.Errorf("cannot load templates: %w", err)
	}

	projectName := savepointinit.ProjectNameFromDir(opts.Dir)
	if err := savepointinit.Scaffold(sub, opts.Dir, projectName, opts.Force); err != nil {
		return err
	}

	promptSub, err := fs.Sub(projectTemplates, "templates/prompts")
	if err != nil {
		return fmt.Errorf("cannot load prompt templates: %w", err)
	}

	prompt, err := savepointinit.RenderMagicPrompt(promptSub, projectName)
	if err != nil {
		return fmt.Errorf("render magic prompt: %w", err)
	}

	fmt.Println(prompt)

	result := savepointinit.CopyToClipboard(prompt)
	switch result.Status {
	case savepointinit.ClipboardCopied:
		fmt.Fprintf(os.Stderr, "prompt copied to clipboard via %s\n", result.Tool)
	case savepointinit.ClipboardFailed:
		fmt.Fprintf(os.Stderr, "warning: clipboard copy failed: %s\n", result.Message)
	}

	if opts.Install {
		if err := savepointinit.InstallDependencies(opts.Dir); err != nil {
			return err
		}
	}

	return nil
}
