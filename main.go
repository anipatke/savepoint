package main

import (
	"context"
	"crypto/rand"
	"embed"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/opencode/savepoint/cmd"
	"github.com/opencode/savepoint/internal/board"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/doctor"
	savepointinit "github.com/opencode/savepoint/internal/init"
	"github.com/opencode/savepoint/internal/migrate"
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
		case "migrate":
			code, err := cmd.RunMigrate(context.Background(), args[1:], os.Stdout, migrateRunner)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			os.Exit(code)
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

// Distinct migrate target diagnostics, so a caller (and a test) can tell a
// missing directory apart from an unwritable one or one that is not a
// Savepoint project at all, per FS-06.
var (
	ErrMigrateTargetMissing      = errors.New("migrate: target directory does not exist")
	ErrMigrateTargetNotSavepoint = errors.New("migrate: target directory is not a Savepoint project")
	ErrMigrateTargetUnwritable   = errors.New("migrate: target directory is not writable")
)

// resolveMigrateTarget validates dir as a migration target and returns its
// absolute path. Unlike data.Discover.FindSavepointRoot, it checks dir
// itself rather than walking up through parent directories: migrate must
// never silently operate on an unrelated ancestor project.
func resolveMigrateTarget(dir string) (string, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("%w: %s: %v", ErrMigrateTargetMissing, dir, err)
	}

	info, statErr := os.Stat(abs)
	if statErr != nil {
		if os.IsNotExist(statErr) {
			return "", fmt.Errorf("%w: %s", ErrMigrateTargetMissing, dir)
		}
		return "", fmt.Errorf("%w: %s: %v", ErrMigrateTargetMissing, dir, statErr)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("%w: %s is not a directory", ErrMigrateTargetMissing, dir)
	}

	savepointDir := filepath.Join(abs, ".savepoint")
	if sInfo, sErr := os.Stat(savepointDir); sErr != nil || !sInfo.IsDir() {
		return "", fmt.Errorf("%w: %s has no .savepoint directory", ErrMigrateTargetNotSavepoint, dir)
	}

	probe := filepath.Join(abs, ".savepoint-migrate-write-test")
	if writeErr := os.WriteFile(probe, []byte{}, 0644); writeErr != nil {
		return "", fmt.Errorf("%w: %s", ErrMigrateTargetUnwritable, dir)
	}
	os.Remove(probe)

	return abs, nil
}

// newMigrationOperationID generates a fresh, effectively-unique operation
// directory name for a real migration run. Tests inject their own
// deterministic migrate.OperationIDSource; this is the one production
// callers use.
func newMigrationOperationID() string {
	var suffix [4]byte
	if _, err := rand.Read(suffix[:]); err != nil {
		return fmt.Sprintf("op-%d", time.Now().UTC().UnixNano())
	}
	return fmt.Sprintf("op-%s-%s", time.Now().UTC().Format("20060102T150405Z"), hex.EncodeToString(suffix[:]))
}

// migrateRunner performs the actual migrate command work: cmd/migrate.go is
// argument parsing and dispatch only, so every filesystem interaction and
// every call into internal/migrate lives here, matching upgradeAssetsRunner
// and initRunner above.
func migrateRunner(ctx context.Context, opts cmd.MigrateOptions) (int, error) {
	root, err := resolveMigrateTarget(opts.Dir)
	if err != nil {
		return 1, err
	}

	willWrite := opts.WillWrite()

	if opts.Recover {
		report, pendingErr := migrate.PendingOperation(root)
		if pendingErr != nil {
			return 1, pendingErr
		}
		if report == nil {
			fmt.Println("savepoint migrate --recover: no incomplete migration operation found")
			if !willWrite {
				return 0, nil
			}
		} else {
			fmt.Println(report.RecoveryGuidance())
			if !willWrite {
				return 0, nil
			}
		}
	}

	now := time.Now

	var decisions migrate.Decisions
	if opts.DecisionsFile != "" {
		decisions, err = migrate.ReadDecisionsFile(opts.DecisionsFile, now())
		if err != nil {
			return 1, err
		}
	}

	plan, err := migrate.Plan(root, decisions, now, newMigrationOperationID)
	if err != nil {
		return 2, err
	}

	if plan.SchemaAlreadyV2 {
		fmt.Print(migrate.FormatPreview(plan))
		return 0, nil
	}

	if len(plan.Conflicts) > 0 {
		fmt.Print(migrate.FormatPreview(plan))
		return 1, fmt.Errorf("migrate: plan reports %d conflict(s); resolve before migrating", len(plan.Conflicts))
	}

	if !willWrite {
		fmt.Print(migrate.FormatPreview(plan))
		if !plan.Appliable {
			return 1, fmt.Errorf("migrate: unresolved blocking ambiguities: %s", strings.Join(plan.UnresolvedBlockingIDs, ", "))
		}
		return 0, nil
	}

	if !plan.Appliable {
		fmt.Print(migrate.FormatPreview(plan))
		return 1, fmt.Errorf("migrate: unresolved blocking ambiguities: %s", strings.Join(plan.UnresolvedBlockingIDs, ", "))
	}

	result, err := migrate.Apply(root, plan)
	if err != nil {
		return 1, err
	}
	switch {
	case result.AlreadyMigrated:
		fmt.Println("project is already at schema_version: 2; nothing to migrate")
	case result.Resumed:
		fmt.Printf("resumed and completed migration operation %s\n", result.OperationID)
	default:
		fmt.Printf("migration complete: operation %s\n", result.OperationID)
	}
	return 0, nil
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
