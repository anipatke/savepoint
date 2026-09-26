// Package migrate's command.go is the behavior behind `savepoint migrate`:
// target validation, decisions loading, preview, and the guarded hand-off to
// Apply. It lives here rather than in main.go so the
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
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/opencode/savepoint/internal/data"
)

// Distinct migrate target diagnostics, so a caller (and a test) can tell a
// missing directory apart from an unwritable one or one that is not a
// Savepoint project at all, per FS-06. ErrTargetMissing and
// ErrTargetNotSavepoint are internal/data's own sentinels: the runtime
// callers that just need to locate a project (board, resume, doctor,
// upgrade-assets) read them from internal/data directly, without importing
// this package, and migrate keeps these names so its own existing callers
// and tests are unaffected.
var (
	ErrTargetMissing      = data.ErrTargetMissing
	ErrTargetNotSavepoint = data.ErrTargetNotSavepoint
	ErrTargetUnwritable   = errors.New("migrate: target directory is not writable")
	ErrGitUnavailable     = errors.New("migrate: git is not available on PATH")
	ErrNotGitWorkTree     = errors.New("migrate: project is not inside a git work tree")
	ErrDirtyGitPaths      = errors.New("migrate: migration paths have uncommitted or untracked changes")
)

// GitCommand runs git with dir as its working directory and returns combined
// output. RunCommand injects it so clean-tree refusals can be tested without
// requiring git in unit tests.
type GitCommand func(dir string, args ...string) (string, error)

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
	// Verbose prints the complete preview listing instead of the summary.
	Verbose bool

	Stdout         io.Writer
	Now            Clock
	NewOperationID OperationIDSource
	RunGit         GitCommand
}

// ResolveTarget validates dir as a migration target and returns its absolute
// path. It is a thin wrapper over data.ResolveTarget, the one implementation
// every runtime caller and migrate itself now share.
func ResolveTarget(dir string) (string, error) {
	return data.ResolveTarget(dir)
}

// FindProjectRoot locates the nearest Savepoint project while walking from
// start toward the filesystem root. It is a thin wrapper over
// data.FindProjectRoot, the one implementation every runtime caller and
// migrate itself now share.
func FindProjectRoot(start string) (string, error) {
	return data.FindProjectRoot(start)
}

// targetWriteabilityProbe is injected in tests so preview tests can prove the
// write-only preflight is never called. It is deliberately separate from
// ResolveTarget: resolving a target is a read-only operation shared by
// preview and apply.
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

// preview renders the plan as the summary, or the full listing when Verbose.
func (opts CommandOptions) preview(plan *ConversionPlan) string {
	if opts.Verbose {
		return FormatPreview(plan)
	}
	return FormatSummaryPreview(plan)
}

// RunCommand performs one `savepoint migrate` invocation and returns its exit
// code: 0 clean, 1 a named refusal, 2 an internal error.
func RunCommand(opts CommandOptions) (int, error) {
	opts = opts.withDefaults()

	root, err := ResolveTarget(opts.Dir)
	if err != nil {
		return 1, err
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
		if opts.Write {
			fmt.Fprintln(opts.Stdout, "project is already migrated (schema_version: 2); nothing to do")
			return 0, nil
		}
		fmt.Fprint(opts.Stdout, opts.preview(plan))
		return 0, nil
	}

	if len(plan.Conflicts) > 0 {
		fmt.Fprint(opts.Stdout, opts.preview(plan))
		return 1, fmt.Errorf("migrate: plan reports %d conflict(s); resolve before migrating", len(plan.Conflicts))
	}

	if !opts.Write {
		if !plan.Appliable {
			fmt.Fprint(opts.Stdout, opts.preview(plan))
			return 1, unresolvedAmbiguityError(plan)
		}
		if err := ensureCleanGitTree(root, plan, opts.RunGit); err != nil {
			fmt.Fprint(opts.Stdout, previewBlockedByGit(plan, opts, err))
			return 1, err
		}
		fmt.Fprint(opts.Stdout, opts.preview(plan))
		return 0, nil
	}

	if !plan.Appliable {
		fmt.Fprint(opts.Stdout, opts.preview(plan))
		return 1, unresolvedAmbiguityError(plan)
	}
	if err := ensureCleanGitTree(root, plan, opts.RunGit); err != nil {
		return 1, err
	}
	if err := ensureTargetWritable(root); err != nil {
		return 1, err
	}

	result, err := Apply(root, plan)
	if err != nil {
		return 1, err
	}
	fmt.Fprintln(opts.Stdout, applyOutcomeMessage(result, plan))
	return 0, nil
}

// NewOperationID generates the run identifier shown in migration previews.
// It is not persisted as recovery state. Tests inject their own deterministic
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
	if o.RunGit == nil {
		o.RunGit = defaultGitCommand
	}
	return o
}

func unresolvedAmbiguityError(plan *ConversionPlan) error {
	return fmt.Errorf("migrate: unresolved blocking ambiguities: %s", strings.Join(plan.UnresolvedBlockingIDs, ", "))
}

func previewBlockedByGit(plan *ConversionPlan, opts CommandOptions, checkErr error) string {
	preview := opts.preview(plan)
	if !opts.Verbose {
		preview = strings.Replace(preview, "Status: ready. To apply, run the same command with --apply.\n", "", 1)
	}

	var out strings.Builder
	out.WriteString(preview)
	if !strings.HasSuffix(preview, "\n") {
		out.WriteByte('\n')
	}
	if !strings.HasSuffix(preview, "\n\n") {
		out.WriteByte('\n')
	}

	var dirty *dirtyGitPathsError
	if errors.As(checkErr, &dirty) {
		out.WriteString("Status: blocked by working-tree changes.\n")
		fmt.Fprintf(&out, "%s\n", formatGitPathGroups(dirty.groups))
		fmt.Fprintf(&out, "Repair: %s\n", dirtyGitRepairAdvice)
		return out.String()
	}

	out.WriteString("Status: blocked by the Git working-tree check.\n")
	fmt.Fprintf(&out, "%s\n", checkErr)
	return out.String()
}

func applyOutcomeMessage(result *ApplyResult, plan *ConversionPlan) string {
	switch {
	case result.AlreadyMigrated:
		return "project is already migrated (schema_version: 2); nothing to do"
	default:
		return fmt.Sprintf("migration complete. Undo from the project root with: %s", gitUndoCommand(plan))
	}
}

func ensureCleanGitTree(root string, plan *ConversionPlan, runGit GitCommand) error {
	inside, err := runGit(root, "rev-parse", "--is-inside-work-tree")
	if errors.Is(err, ErrGitUnavailable) {
		return fmt.Errorf("%w: install git, then retry `savepoint migrate --apply`", ErrGitUnavailable)
	}
	if err != nil || strings.TrimSpace(inside) != "true" {
		return fmt.Errorf("%w: run `git init` in the project or move it inside a git work tree", ErrNotGitWorkTree)
	}

	paths := plannedGitPaths(plan)
	args := []string{"status", "--porcelain", "--untracked-files=all", "--ignored=matching", "--no-renames", "--"}
	args = append(args, paths...)
	status, err := runGit(root, args...)
	if err != nil {
		return fmt.Errorf("migrate: inspect git status for migration paths: %w", err)
	}
	if strings.TrimSpace(status) != "" {
		return &dirtyGitPathsError{groups: groupGitStatusPaths(status)}
	}
	return nil
}

const dirtyGitRepairAdvice = "stage and commit or stash changed paths, and move ignored files out of planned paths before retrying `savepoint migrate --apply`"

type gitStatusGroups struct {
	modified  []string
	untracked []string
	ignored   []string
}

type dirtyGitPathsError struct {
	groups gitStatusGroups
}

func (e *dirtyGitPathsError) Error() string {
	return fmt.Sprintf("%s:\n%s\nRepair: %s", ErrDirtyGitPaths, formatGitPathGroups(e.groups), dirtyGitRepairAdvice)
}

func (e *dirtyGitPathsError) Unwrap() error { return ErrDirtyGitPaths }

func groupGitStatusPaths(status string) gitStatusGroups {
	var groups gitStatusGroups
	for _, line := range strings.Split(status, "\n") {
		if line == "" {
			continue
		}
		if len(line) < 3 {
			groups.modified = append(groups.modified, line)
			continue
		}

		code := line[:2]
		path := line[3:]
		switch code {
		case "??":
			groups.untracked = append(groups.untracked, path)
		case "!!":
			groups.ignored = append(groups.ignored, path)
		default:
			groups.modified = append(groups.modified, path)
		}
	}
	sort.Strings(groups.modified)
	sort.Strings(groups.untracked)
	sort.Strings(groups.ignored)
	return groups
}

func formatGitPathGroups(groups gitStatusGroups) string {
	var out strings.Builder
	writeGroup := func(name string, paths []string) {
		if len(paths) == 0 {
			return
		}
		fmt.Fprintf(&out, "%s:\n", name)
		for _, path := range paths {
			fmt.Fprintf(&out, "  - %s\n", path)
		}
	}
	writeGroup("Modified", groups.modified)
	writeGroup("Untracked", groups.untracked)
	writeGroup("Ignored", groups.ignored)
	return strings.TrimSuffix(out.String(), "\n")
}

func defaultGitCommand(dir string, args ...string) (string, error) {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		return "", ErrGitUnavailable
	}
	commandArgs := append([]string{"--literal-pathspecs"}, args...)
	command := exec.Command(gitPath, commandArgs...)
	command.Dir = dir
	if len(args) > 0 && args[0] == "status" {
		command.Env = envWithGitOptionalLocksDisabled(os.Environ())
	}
	output, err := command.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(string(output)))
	}
	return string(output), nil
}

func envWithGitOptionalLocksDisabled(env []string) []string {
	const setting = "GIT_OPTIONAL_LOCKS=0"
	for i, entry := range env {
		if strings.HasPrefix(entry, "GIT_OPTIONAL_LOCKS=") {
			env[i] = setting
			return env
		}
	}
	return append(env, setting)
}

func plannedGitPaths(plan *ConversionPlan) []string {
	paths := []string{".savepoint/config.yml", manifestRelPath()}
	for _, target := range plan.Targets {
		paths = append(paths, savepointPath(target.InstallPath()))
	}
	for _, document := range plan.Documents {
		paths = append(paths, savepointPath(document.TargetPath))
	}
	for _, archive := range plan.Archives {
		paths = append(paths, archive.ArchivePath, archive.SourcePath)
	}
	return uniqueSortedPaths(paths)
}

func uniqueSortedPaths(paths []string) []string {
	seen := make(map[string]struct{}, len(paths))
	unique := make([]string, 0, len(paths))
	for _, path := range paths {
		path = filepath.ToSlash(filepath.Clean(filepath.FromSlash(path)))
		if _, exists := seen[path]; exists {
			continue
		}
		seen[path] = struct{}{}
		unique = append(unique, path)
	}
	sort.Strings(unique)
	return unique
}

func gitUndoCommand(plan *ConversionPlan) string {
	tracked, untracked := gitUndoPathGroups(plan)
	clean := "git --literal-pathspecs clean -fdx -- " + shellJoin(untracked)
	if len(tracked) == 0 {
		// git restore refuses an empty pathspec, and `&&` would then skip
		// the clean, so a plan with nothing tracked prints the clean alone.
		return clean
	}
	return "git --literal-pathspecs restore --source=HEAD --staged --worktree -- " + shellJoin(tracked) +
		" && " + clean
}

func gitUndoPathGroups(plan *ConversionPlan) (tracked, untracked []string) {
	if planHasSource(plan, ".savepoint/config.yml") {
		tracked = append(tracked, ".savepoint/config.yml")
	} else {
		untracked = append(untracked, ".savepoint/config.yml")
	}
	untracked = append(untracked, manifestRelPath())
	for _, target := range plan.Targets {
		untracked = append(untracked, savepointPath(target.InstallPath()))
	}
	for _, document := range plan.Documents {
		path := savepointPath(document.TargetPath)
		if document.Kind == DocumentRouter {
			tracked = append(tracked, path)
		} else {
			untracked = append(untracked, path)
		}
	}
	for _, archive := range plan.Archives {
		tracked = append(tracked, archive.SourcePath)
		untracked = append(untracked, archive.ArchivePath)
	}
	return uniqueSortedPaths(tracked), uniqueSortedPaths(untracked)
}

func planHasSource(plan *ConversionPlan, path string) bool {
	for _, source := range plan.Sources {
		if source.Path == path {
			return true
		}
	}
	return false
}

func shellJoin(paths []string) string {
	quoted := make([]string, len(paths))
	for i, path := range paths {
		quoted[i] = "'" + strings.ReplaceAll(path, "'", "'\\''") + "'"
	}
	return strings.Join(quoted, " ")
}
