// Package migrate's command.go is the behavior behind `savepoint migrate`:
// target validation, the recover report, decisions loading, preview, and the
// guarded hand-off to Apply. It lives here rather than in main.go so the
// command's real body is exercised by tests (ARCH-01 keeps cmd/ and main.go
// to argument parsing and dispatch) and so the clock and operation-ID source
// stay injectable — the same injection Plan already requires to render
// byte-identical output across runs.
//
// The default-preview policy itself is not repeated here: cmd.MigrateOptions
// owns it and passes the already-resolved decision as CommandOptions.Write.
package migrate

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Distinct migrate target diagnostics, so a caller (and a test) can tell a
// missing directory apart from an unwritable one or one that is not a
// Savepoint project at all, per FS-06.
var (
	ErrTargetMissing      = errors.New("migrate: target directory does not exist")
	ErrTargetNotSavepoint = errors.New("migrate: target directory is not a Savepoint project")
	ErrTargetUnwritable   = errors.New("migrate: target directory is not writable")
)

// CommandOptions is everything RunCommand needs: the parsed user intent, the
// stream to report on, and the two injected sources that make a run
// reproducible.
type CommandOptions struct {
	Dir string
	// Write is the already-resolved decision to apply rather than preview.
	// The rule that produces it (apply without dry-run) belongs to the
	// command layer and is deliberately not duplicated here.
	Write         bool
	DecisionsFile string
	Recover       bool

	Stdout         io.Writer
	Now            Clock
	NewOperationID OperationIDSource
}

// ResolveTarget validates dir as a migration target and returns its absolute
// path. Unlike data.Discover.FindSavepointRoot, it checks dir itself rather
// than walking up through parent directories: migrate must never silently
// operate on an unrelated ancestor project.
func ResolveTarget(dir string) (string, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("%w: %s: %v", ErrTargetMissing, dir, err)
	}

	info, statErr := os.Stat(abs)
	if statErr != nil {
		if os.IsNotExist(statErr) {
			return "", fmt.Errorf("%w: %s", ErrTargetMissing, dir)
		}
		return "", fmt.Errorf("%w: %s: %v", ErrTargetMissing, dir, statErr)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("%w: %s is not a directory", ErrTargetMissing, dir)
	}

	savepointDir := filepath.Join(abs, ".savepoint")
	if sInfo, sErr := os.Stat(savepointDir); sErr != nil || !sInfo.IsDir() {
		return "", fmt.Errorf("%w: %s has no .savepoint directory", ErrTargetNotSavepoint, dir)
	}

	return abs, nil
}

// targetWriteabilityProbe is injected in tests so preview tests can prove the
// write-only preflight is never called. It is deliberately separate from
// ResolveTarget: resolving a target is a read-only operation shared by
// preview and recovery-report modes.
var targetWriteabilityProbe = probeTargetWriteability

func probeTargetWriteability(root string) (err error) {
	probe, err := os.CreateTemp(root, ".savepoint-migrate-write-test-*")
	if err != nil {
		return err
	}
	probePath := probe.Name()
	defer func() {
		// CreateTemp used O_CREATE|O_EXCL, so this cleanup can only target the
		// unique file this probe created, never a user's pre-existing sentinel.
		if removeErr := os.Remove(probePath); removeErr != nil && !os.IsNotExist(removeErr) {
			if err == nil {
				err = removeErr
				return
			}
			err = fmt.Errorf("%w; cleanup writeability probe: %v", err, removeErr)
		}
	}()
	err = probe.Close()
	return err
}

func ensureTargetWritable(root string) error {
	if err := targetWriteabilityProbe(root); err != nil {
		return fmt.Errorf("%w: %s", ErrTargetUnwritable, root)
	}
	return nil
}

// RunCommand performs one `savepoint migrate` invocation and returns its exit
// code: 0 clean, 1 a named refusal, 2 an internal error.
func RunCommand(opts CommandOptions) (int, error) {
	opts = opts.withDefaults()

	root, err := ResolveTarget(opts.Dir)
	if err != nil {
		return 1, err
	}

	if opts.Recover {
		code, done, err := reportRecovery(opts, root)
		if done {
			return code, err
		}
	}

	if opts.Write {
		if err := ensureTargetWritable(root); err != nil {
			return 1, err
		}
		// A pending operation is authoritative. Do not build a fresh plan:
		// doing so would generate a new operation ID/time and would also see
		// the pending operation's already-created destinations as collisions.
		pending, pendingErr := PendingOperation(root)
		if pendingErr != nil {
			return 1, pendingErr
		}
		if pending != nil {
			result, applyErr := Apply(root, nil)
			if applyErr != nil {
				return 1, applyErr
			}
			fmt.Fprintln(opts.Stdout, applyOutcomeMessage(result))
			return 0, nil
		}
	}

	var decisions Decisions
	if opts.DecisionsFile != "" {
		decisions, err = ReadDecisionsFile(opts.DecisionsFile, opts.Now())
		if err != nil {
			return 1, err
		}
	}

	plan, err := Plan(root, decisions, opts.Now, opts.NewOperationID)
	if err != nil {
		return 2, err
	}

	if plan.SchemaAlreadyV2 {
		fmt.Fprint(opts.Stdout, FormatPreview(plan))
		return 0, nil
	}

	if len(plan.Conflicts) > 0 {
		fmt.Fprint(opts.Stdout, FormatPreview(plan))
		return 1, fmt.Errorf("migrate: plan reports %d conflict(s); resolve before migrating", len(plan.Conflicts))
	}

	if !opts.Write {
		fmt.Fprint(opts.Stdout, FormatPreview(plan))
		if !plan.Appliable {
			return 1, unresolvedAmbiguityError(plan)
		}
		return 0, nil
	}

	if !plan.Appliable {
		fmt.Fprint(opts.Stdout, FormatPreview(plan))
		return 1, unresolvedAmbiguityError(plan)
	}

	result, err := Apply(root, plan)
	if err != nil {
		return 1, err
	}
	fmt.Fprintln(opts.Stdout, applyOutcomeMessage(result))
	return 0, nil
}

// NewOperationID generates a fresh, effectively-unique operation directory
// name for a real migration run. Tests inject their own deterministic
// OperationIDSource; this is the one production callers use.
func NewOperationID() string {
	var suffix [4]byte
	if _, err := rand.Read(suffix[:]); err != nil {
		return fmt.Sprintf("op-%d", time.Now().UTC().UnixNano())
	}
	return fmt.Sprintf("op-%s-%s", time.Now().UTC().Format("20060102T150405Z"), hex.EncodeToString(suffix[:]))
}

func (o CommandOptions) withDefaults() CommandOptions {
	if o.Stdout == nil {
		o.Stdout = os.Stdout
	}
	if o.Now == nil {
		o.Now = time.Now
	}
	if o.NewOperationID == nil {
		o.NewOperationID = NewOperationID
	}
	return o
}

// reportRecovery prints the pending-operation report --recover asks for. It
// returns done=true when the invocation is finished (a preview-mode
// --recover reports and stops); an --apply --recover run falls through to the
// recorded-operation recovery path before any fresh planning.
func reportRecovery(opts CommandOptions, root string) (code int, done bool, err error) {
	report, pendingErr := PendingOperation(root)
	if pendingErr != nil {
		return 1, true, pendingErr
	}
	if report == nil {
		fmt.Fprintln(opts.Stdout, "savepoint migrate --recover: no incomplete migration operation found")
	} else {
		fmt.Fprintln(opts.Stdout, report.RecoveryGuidance())
	}
	if !opts.Write {
		return 0, true, nil
	}
	return 0, false, nil
}

func unresolvedAmbiguityError(plan *ConversionPlan) error {
	return fmt.Errorf("migrate: unresolved blocking ambiguities: %s", strings.Join(plan.UnresolvedBlockingIDs, ", "))
}

func applyOutcomeMessage(result *ApplyResult) string {
	switch {
	case result.AlreadyMigrated:
		return "project is already at schema_version: 2; nothing to migrate"
	case result.Resumed:
		return fmt.Sprintf("resumed and completed migration operation %s", result.OperationID)
	default:
		return fmt.Sprintf("migration complete: operation %s", result.OperationID)
	}
}
