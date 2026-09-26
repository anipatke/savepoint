package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/opencode/savepoint/cmd"
	"github.com/opencode/savepoint/internal/board"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/doctor"
	savepointinit "github.com/opencode/savepoint/internal/init"
	"github.com/opencode/savepoint/internal/migrate"
	"github.com/opencode/savepoint/internal/resume"
)

//go:embed templates/prompts
var promptTemplates embed.FS

//go:embed templates/project-v2
//go:embed all:templates/project-v2/.savepoint
var projectTemplatesV2 embed.FS

var version = "dev"

func main() {
	args, debug := stripDebugFlag(os.Args[1:])
	if debug || os.Getenv("SAVEPOINT_DEBUG") != "" {
		board.SetDebug(true)
	}

	if len(args) > 0 {
		switch args[0] {
		case "--help", "-h", "help":
			fmt.Print(mainUsage)
			os.Exit(0)
		case "--version":
			fmt.Println(version)
			os.Exit(0)
		case "init":
			if err := cmd.RunInit(context.Background(), args[1:], os.Stdout, initRunner); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			os.Exit(0)
		case "create-task":
			if err := cmd.RunCreateTask(context.Background(), args[1:], os.Stdout, os.Stderr, createTaskRunner); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			os.Exit(0)
		case "board":
			if err := cmd.RunBoard(context.Background(), args[1:], os.Stdout, func(opts cmd.BoardOptions) error {
				return board.RunWithFilters(board.Filters{
					Objective: opts.Objective,
				})
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
		case "migrate":
			code, err := cmd.RunMigrate(context.Background(), args[1:], os.Stdout, migrateRunner)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			os.Exit(code)
		case "resume":
			code, err := cmd.RunResume(context.Background(), args[1:], os.Stdout, resumeRunner)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			os.Exit(code)
		default:
			fmt.Fprintf(os.Stderr, "unknown command %q\n\n%s", args[0], mainUsage)
			os.Exit(2)
		}
	}
	// A board refusal, such as an unmigrated project, is a routine message for
	// the user, not a crash.
	if err := board.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

const mainUsage = `Usage: savepoint <command> [options]

Commands:
  init [dir] [--force] [--install]       Create a V2 project
  create-task --objective <O-###> --draft <path> [dir]  Create a Task with the next ID
  board [--objective <objective>]        Open the V2 board
  doctor                                Check project health
  resume [dir]                           Print the next V2 action
  migrate [dir] [--apply]                Convert a legacy project to V2
  upgrade-assets [dir]                   Refresh assets in an existing V2 project

Run ` + "`savepoint <command> --help`" + ` for command-specific options.

Migration and asset upgrades are separate: use migrate for a legacy project,
then use upgrade-assets to refresh its V2-managed assets.
`

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

func runDoctorChecks(_ cmd.DoctorOptions) (int, error) {
	projectRoot, err := data.FindProjectRoot(".")
	if err != nil {
		return 2, fmt.Errorf("savepoint root not found: %w", err)
	}
	root := filepath.Join(projectRoot, ".savepoint")

	if err := data.CheckRuntimeSchema(projectRoot); err != nil {
		return 1, fmt.Errorf("doctor: %s", err)
	}

	report := doctor.RunV2Checks(root)
	fmt.Fprint(os.Stdout, report.Format())

	if report.HasProblems() {
		return 1, nil
	}
	return 0, nil
}

func upgradeAssetsRunner(ctx context.Context, opts cmd.UpgradeAssetsOptions) error {
	subV2, err := fs.Sub(projectTemplatesV2, "templates/project-v2")
	if err != nil {
		return fmt.Errorf("cannot load templates: %w", err)
	}

	// V1 projects are deliberately refused inside UpgradeProjectAssets; the
	// explicit migrate command is the only path allowed to change their
	// workflow.
	report, err := savepointinit.UpgradeProjectAssets(subV2, opts.Dir, opts.DryRun, opts.Force)
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

// migrateRunner is the production wiring for the migrate command: it maps the
// parsed options onto migrate.CommandOptions and supplies the real clock and
// operation-ID source. All behavior lives in internal/migrate, matching
// upgradeAssetsRunner and initRunner above.
func migrateRunner(ctx context.Context, opts cmd.MigrateOptions) (int, error) {
	return migrate.RunCommand(migrate.CommandOptions{
		Dir:            opts.Dir,
		Write:          opts.WillWrite(),
		DecisionsFile:  opts.DecisionsFile,
		Verbose:        opts.Verbose,
		Stdout:         os.Stdout,
		Now:            time.Now,
		NewOperationID: migrate.NewOperationID,
	})
}

// resumeRunner is the production wiring for the resume command: real stdout,
// matching upgradeAssetsRunner and migrateRunner above. All behavior lives in
// runResume, which is unit-testable with an injected writer.
func resumeRunner(ctx context.Context, opts cmd.ResumeOptions) (int, error) {
	return runResume(opts.Dir, os.Stdout)
}

// runResume performs one `savepoint resume` invocation: it resolves opts.Dir
// exactly as migrate does (ARCH-03), runs the small runtime schema check
// before interpreting any records, and gives a legacy project named migration
// guidance instead of rendering it as V2. For an intact V2 project it loads
// the project, reads its V2 router, resolves the Next projection, and renders
// it to stdout. It performs no write on any path, including every failure
// path: every read here is os.ReadFile, os.Stat, or a pure computation,
// never a write. The schema-1 case reports through stdout as an explanatory
// message (exit 1, no error); every other nonzero path returns an error so
// the dispatch above prints it to stderr.
func runResume(dir string, stdout io.Writer) (int, error) {
	root, err := data.ResolveTarget(dir)
	if err != nil {
		return 1, resumeTargetError(dir, err)
	}
	savepointRoot := filepath.Join(root, ".savepoint")

	if err := data.CheckRuntimeSchema(root); err != nil {
		if errors.Is(err, data.ErrSchemaMigrationRequired) {
			fmt.Fprintf(stdout, "resume: %s\n", err)
			return 1, nil
		}
		return 1, fmt.Errorf("resume: %s", err)
	}

	index, err := data.LoadV2Index(savepointRoot)
	if err != nil {
		return 1, fmt.Errorf("resume: loading project: %w", err)
	}

	routerContent, err := os.ReadFile(filepath.Join(savepointRoot, "router.md"))
	if err != nil {
		return 1, fmt.Errorf("resume: reading router: %w", err)
	}
	router, err := data.NewRouterReader().ReadStateV2(string(routerContent))
	if err != nil {
		return 1, fmt.Errorf("resume: reading router: %w", err)
	}

	next := data.ResolveNext(data.NextInput{Index: index, Router: router})

	if err := resume.Render(stdout, next); err != nil {
		return 1, fmt.Errorf("resume: writing output: %w", err)
	}
	return 0, nil
}

// resumeTargetError names data.ResolveTarget's two target diagnostics for
// the command the user actually ran, so a failed `savepoint resume /typo`
// reports a resume-specific message rather than a generic wrapped error.
func resumeTargetError(dir string, err error) error {
	switch {
	case errors.Is(err, data.ErrTargetMissing):
		return fmt.Errorf("resume: target directory does not exist: %s", dir)
	case errors.Is(err, data.ErrTargetNotSavepoint):
		return fmt.Errorf("resume: target directory is not a Savepoint project: %s has no .savepoint directory", dir)
	default:
		return fmt.Errorf("resume: %w", err)
	}
}

func initRunner(ctx context.Context, opts cmd.InitOptions) error {
	if err := savepointinit.ValidateTarget(opts.Dir, opts.Force); err != nil {
		return err
	}

	sub, err := fs.Sub(projectTemplatesV2, "templates/project-v2")
	if err != nil {
		return fmt.Errorf("cannot load templates: %w", err)
	}

	projectName := savepointinit.ProjectNameFromDir(opts.Dir)
	if err := savepointinit.Scaffold(sub, opts.Dir, projectName, opts.Force); err != nil {
		return err
	}

	promptSub, err := fs.Sub(promptTemplates, "templates/prompts")
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

func createTaskRunner(_ context.Context, opts cmd.CreateTaskOptions) (string, string, error) {
	projectRoot, err := data.ResolveTarget(opts.Dir)
	if err != nil {
		return "", "", fmt.Errorf("create-task: %w", err)
	}
	if err := data.CheckRuntimeSchema(projectRoot); err != nil {
		return "", "", fmt.Errorf("create-task: %w", err)
	}
	draft, err := os.ReadFile(opts.Draft)
	if err != nil {
		return "", "", fmt.Errorf("create-task: read draft %s: %w", opts.Draft, err)
	}
	created, err := data.CreateTaskV2(filepath.Join(projectRoot, ".savepoint"), opts.Objective, string(draft))
	if err != nil {
		return "", "", err
	}
	return created.ID, filepath.ToSlash(created.Path), nil
}
