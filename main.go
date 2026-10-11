package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/opencode/savepoint/cmd"
	"github.com/opencode/savepoint/internal/board"
	"github.com/opencode/savepoint/internal/codehealth"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/doctor"
	"github.com/opencode/savepoint/internal/healthcheck"
	savepointinit "github.com/opencode/savepoint/internal/init"
	"github.com/opencode/savepoint/internal/legacydep"
	"github.com/opencode/savepoint/internal/migrate"
	"github.com/opencode/savepoint/internal/resume"
)

//go:embed templates/prompts
var promptTemplates embed.FS

//go:embed all:templates/project-v2
var projectTemplatesV2 embed.FS

var version = "dev"

func main() {
	args, debug := stripDebugFlag(os.Args[1:])
	if debug || os.Getenv("SAVEPOINT_DEBUG") != "" {
		board.SetDebug(true)
	}

	if len(args) > 0 {
		os.Exit(runCommand(args))
	}
	// A board refusal, such as an unmigrated project, is a routine message for
	// the user, not a crash.
	if err := board.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// exitFor prints err, when there is one, and returns the exit status for a
// command that reports only success or failure.
func exitFor(err error) int {
	if err == nil {
		return 0
	}
	fmt.Fprintln(os.Stderr, err)
	return 1
}

// exitWith prints err, when there is one, and returns the command's own exit
// code, which does not depend on whether an error was printed.
func exitWith(code int, err error) int {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return code
}

// runCommand runs the named subcommand and returns its exit status.
func runCommand(args []string) int {
	switch args[0] {
	case "--help", "-h", "help":
		fmt.Print(mainUsage)
		return 0
	case "--version":
		fmt.Println(version)
		return 0
	case "init":
		return exitFor(cmd.RunInit(context.Background(), args[1:], os.Stdout, initRunner))
	case "create-task":
		return exitFor(cmd.RunCreateTask(context.Background(), args[1:], os.Stdout, os.Stderr, createTaskRunner))
	case "board":
		return exitFor(cmd.RunBoard(context.Background(), args[1:], os.Stdout, func(opts cmd.BoardOptions) error {
			return board.RunWithFilters(board.Filters{
				Objective: opts.Objective,
			})
		}))
	case "doctor":
		code, err := cmd.RunDoctor(context.Background(), args[1:], os.Stdout, func(opts cmd.DoctorOptions) (int, error) {
			return runDoctorChecks(opts)
		})
		return exitWith(code, err)
	case "upgrade-assets":
		return exitFor(cmd.RunUpgradeAssets(context.Background(), args[1:], os.Stdout, upgradeAssetsRunner))
	case "migrate":
		code, err := cmd.RunMigrate(context.Background(), args[1:], os.Stdout, migrateRunner)
		return exitWith(code, err)
	case "health":
		// Ctrl-C cancels collection through the context so a running tool is stopped.
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()
		return exitFor(cmd.RunHealth(ctx, args[1:], os.Stdout, cmd.HealthRunners{Setup: healthSetupRunner, Check: healthCheckRunner, Report: healthReportRunner}))
	case "resume":
		code, err := cmd.RunResume(context.Background(), args[1:], os.Stdout, resumeRunner)
		return exitWith(code, err)
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n%s", args[0], mainUsage)
		return 2
	}
}

const mainUsage = `Usage: savepoint <command> [options]

Commands:
  init [dir] [--force] [--install]       Create a V2 project
  create-task --objective <O-###> --draft <path> [dir]  Create a Task with the next ID
  board [--objective <objective>]        Open the V2 board
  doctor                                Check project health
  resume [dir]                           Print the next V2 action
  health setup [dir] [--apply]           Preview or save suggested health tools
  health check O-### [dir]               Collect an official health snapshot for an Objective
  health report [dir]                    Rewrite .savepoint/health/report.md from the newest snapshot
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
	projectDir := opts.Dir
	if root, err := data.FindProjectRoot(opts.Dir); err == nil {
		projectDir = root
	}
	if found, ok := legacydep.Detect(projectDir); ok {
		fmt.Printf("\nWarning: %s\n%s\n", found.Message(), found.Repair())
	}
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

	// An unreadable config.yml leaves the optional advice off; resume's own
	// answer never depends on it.
	parallel := false
	if config, err := data.NewConfigReader().Read(filepath.Join(savepointRoot, "config.yml")); err == nil {
		parallel = config.ParallelPlanningEnabled()
	}
	next := data.ResolveNext(data.NextInput{Index: index, Router: router, ParallelPlanning: parallel})

	if err := resume.Render(stdout, next, index); err != nil {
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

	if advice := savepointinit.ClaudeSettingsAdvice(opts.Dir); advice != "" {
		fmt.Fprintln(os.Stderr, advice)
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

	// Discovery trouble never fails an init that already scaffolded.
	if err := previewHealthSetup(ctx, opts.Dir, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "warning: health tool suggestions unavailable: %s\n", err)
	}

	return nil
}

// healthSetupRunner is the production wiring for `health setup`. Discovery,
// reconciliation, and rendering live in internal/codehealth.
func healthSetupRunner(ctx context.Context, opts cmd.HealthSetupOptions) error {
	root, plan, err := planHealthSetup(ctx, opts.Dir)
	if err != nil {
		return err
	}
	if !opts.Apply {
		_, err = fmt.Fprint(os.Stdout, plan.Preview())
		return err
	}
	changed, err := plan.Apply(codehealth.NewStore(root))
	if err != nil {
		return fmt.Errorf("health setup: %w", err)
	}
	_, err = fmt.Fprint(os.Stdout, plan.Applied(changed))
	return err
}

// healthCheckRunner is the production wiring for `health check`; the
// collection itself lives in internal/healthcheck.
func healthCheckRunner(ctx context.Context, opts cmd.HealthCheckOptions) error {
	return healthcheck.Run(ctx, healthcheck.Request{Dir: opts.Dir, Objective: opts.Objective, Stderr: os.Stderr}, os.Stdout)
}

// healthReportRunner is the production wiring for `health report`.
func healthReportRunner(_ context.Context, opts cmd.HealthReportOptions) error {
	return healthcheck.RunReport(healthcheck.ReportRequest{Dir: opts.Dir}, os.Stdout)
}

// previewHealthSetup prints what setup would suggest without writing anything.
func previewHealthSetup(ctx context.Context, dir string, stdout io.Writer) error {
	_, plan, err := planHealthSetup(ctx, dir)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(stdout, "\n"+plan.Preview())
	return err
}

func planHealthSetup(ctx context.Context, dir string) (string, codehealth.SetupPlan, error) {
	root, err := data.ResolveTarget(dir)
	switch {
	case errors.Is(err, data.ErrTargetMissing):
		return "", codehealth.SetupPlan{}, fmt.Errorf("health setup: target directory does not exist: %s", dir)
	case errors.Is(err, data.ErrTargetNotSavepoint):
		return "", codehealth.SetupPlan{}, fmt.Errorf("health setup: target directory is not a Savepoint project: %s has no .savepoint directory", dir)
	case err != nil:
		return "", codehealth.SetupPlan{}, fmt.Errorf("health setup: %w", err)
	}
	if err := data.CheckRuntimeSchema(root); err != nil {
		return "", codehealth.SetupPlan{}, fmt.Errorf("health setup: %w", err)
	}
	plan, err := codehealth.PlanProject(ctx, root, exec.LookPath)
	if err != nil {
		return "", codehealth.SetupPlan{}, fmt.Errorf("health setup: %w", err)
	}
	return root, plan, nil
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
