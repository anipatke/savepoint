package codehealth

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"
)

// Named repository-state errors. Callers match them with errors.Is. Messages
// never carry Git's own output, which can hold remote URLs or user names.
var (
	ErrGitMissing      = errors.New("git is not available")
	ErrNotRepository   = errors.New("not inside a Git work tree")
	ErrGitCommand      = errors.New("git command failed")
	ErrInputUnreadable = errors.New("relevant input cannot be read")
)

// DefaultGitTimeout bounds one whole repository observation or comparison.
const DefaultGitTimeout = 30 * time.Second

// maxGitOutput caps what one Git command may return, so a huge tree cannot
// exhaust memory.
const maxGitOutput = 32 << 20

// CommandResult is the outcome of a command that ran to completion. A non-zero
// ExitCode is a result, not an error, because some Git queries answer through
// it (for example merge-base --is-ancestor).
type CommandResult struct {
	Stdout   []byte
	ExitCode int
}

// CommandRunner is the only boundary to the outside world. It receives an
// argument vector, never a shell string. An error means the command could not
// run or was cancelled.
type CommandRunner interface {
	Run(ctx context.Context, dir string, args ...string) (CommandResult, error)
}

// GitRunner runs the git executable without a shell. Stderr is discarded and
// the environment is pinned so output does not depend on the user's locale or
// trigger prompts and index writes.
type GitRunner struct {
	// Executable defaults to "git".
	Executable string
}

// Run implements CommandRunner.
func (g GitRunner) Run(ctx context.Context, dir string, args ...string) (CommandResult, error) {
	exe := g.Executable
	if exe == "" {
		exe = "git"
	}
	cmd := exec.CommandContext(ctx, exe, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "LC_ALL=C", "GIT_TERMINAL_PROMPT=0", "GIT_OPTIONAL_LOCKS=0")
	out := &cappedBuffer{limit: maxGitOutput}
	cmd.Stdout = out
	err := cmd.Run()
	if ctx.Err() != nil {
		return CommandResult{}, ctx.Err()
	}
	var exit *exec.ExitError
	switch {
	case err == nil:
		return CommandResult{Stdout: out.Bytes()}, nil
	case out.overflow:
		return CommandResult{}, fmt.Errorf("%w: output exceeds %d bytes", ErrGitCommand, maxGitOutput)
	case errors.As(err, &exit):
		return CommandResult{Stdout: out.Bytes(), ExitCode: exit.ExitCode()}, nil
	case errors.Is(err, exec.ErrNotFound), errors.Is(err, fs.ErrNotExist):
		return CommandResult{}, ErrGitMissing
	default:
		return CommandResult{}, fmt.Errorf("%w: cannot start", ErrGitCommand)
	}
}

// cappedBuffer stops accepting output past limit and fails the command.
type cappedBuffer struct {
	bytes.Buffer
	limit    int
	overflow bool
}

func (b *cappedBuffer) Write(p []byte) (int, error) {
	if b.Len()+len(p) > b.limit {
		b.overflow = true
		return 0, errors.New("output limit reached")
	}
	return b.Buffer.Write(p)
}

// InputScope selects which repository files count as Code Health inputs.
// Paths are relative to the project root. Empty Include means every file.
// Patterns use forward slashes; "*" matches within one segment, "**" matches
// any number of segments, and a pattern naming a directory covers its contents.
type InputScope struct {
	Include []string
	Exclude []string
}

// healthExclusion keeps stored snapshots out of the fingerprint: writing a
// snapshot must not change the state it describes.
const healthExclusion = savepointDir + "/" + healthDir + "/**"

// Validate checks the patterns with the same rules as configured scope.
func (s InputScope) Validate() error {
	if err := validatePatternList("scope.include", s.Include); err != nil {
		return err
	}
	return validatePatternList("scope.exclude", s.Exclude)
}

func (s InputScope) relevant(p string) bool {
	if matchesAny(healthExclusion, p) || slices.ContainsFunc(s.Exclude, func(x string) bool { return matchesAny(x, p) }) {
		return false
	}
	return len(s.Include) == 0 || slices.ContainsFunc(s.Include, func(x string) bool { return matchesAny(x, p) })
}

func matchesAny(pattern, p string) bool {
	return matchSegments(splitPath(pattern), splitPath(p))
}

// matchSegments matches pattern segments against path segments. Consuming the
// whole pattern matches any deeper path, so a directory pattern covers its
// contents.
func matchSegments(pat, segs []string) bool {
	if len(pat) == 0 {
		return true
	}
	if pat[0] == "**" {
		for i := 0; i <= len(segs); i++ {
			if matchSegments(pat[1:], segs[i:]) {
				return true
			}
		}
		return false
	}
	if len(segs) == 0 {
		return false
	}
	// Patterns are validated at the configuration boundary (validateGlob), so
	// a match error cannot occur here.
	ok, err := path.Match(pat[0], segs[0])
	return err == nil && ok && matchSegments(pat[1:], segs[1:])
}

// Observation is the described repository state at one moment. Shallow is kept
// beside the persisted identity because it only affects how a later comparison
// may be worded.
type Observation struct {
	Identity RepositoryIdentity
	Shallow  bool
}

// Unborn reports a repository with no useful commit.
func (o Observation) Unborn() bool { return o.Identity.Commit == "" }

// ObserveRepository describes the work tree rooted at root. It runs only
// read-only local Git queries and hashes file content in place; nothing is
// fetched, written, or retained.
func ObserveRepository(ctx context.Context, runner CommandRunner, root string, scope InputScope) (Observation, error) {
	if err := scope.Validate(); err != nil {
		return Observation{}, err
	}
	ctx, cancel := context.WithTimeout(ctx, DefaultGitTimeout)
	defer cancel()
	q := gitQuery{runner: runner, dir: root}

	if out, err := q.text(ctx, "rev-parse", "--is-inside-work-tree"); err != nil || out != "true" {
		if err != nil && !errors.Is(err, errExit) {
			return Observation{}, err
		}
		return Observation{}, ErrNotRepository
	}
	commit, err := q.headCommit(ctx)
	if err != nil {
		return Observation{}, err
	}
	shallow, err := q.shallow(ctx)
	if err != nil {
		return Observation{}, err
	}
	tracked, untracked, err := q.inputs(ctx, scope)
	if err != nil {
		return Observation{}, err
	}
	fingerprint, err := fingerprintInputs(ctx, root, slices.Concat(tracked, untracked))
	if err != nil {
		return Observation{}, err
	}
	dirty := commit == "" || len(untracked) > 0
	if !dirty {
		dirty, err = q.trackedChanged(ctx, scope)
		if err != nil {
			return Observation{}, err
		}
	}
	return Observation{
		Identity: RepositoryIdentity{Commit: commit, Dirty: dirty, InputFingerprint: fingerprint},
		Shallow:  shallow,
	}, nil
}

// errExit marks a Git query that ran but exited non-zero.
var errExit = errors.New("non-zero exit")

// gitQuery wraps a runner with the few Git questions this package asks.
type gitQuery struct {
	runner CommandRunner
	dir    string
}

func (q gitQuery) run(ctx context.Context, args ...string) (CommandResult, error) {
	res, err := q.runner.Run(ctx, q.dir, args...)
	switch {
	case err == nil:
		return res, nil
	case errors.Is(err, ErrGitMissing), errors.Is(err, ErrGitCommand),
		errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return res, err
	default:
		return res, fmt.Errorf("%w: %s", ErrGitCommand, args[0])
	}
}

// text runs a query expected to succeed and returns its trimmed output.
func (q gitQuery) text(ctx context.Context, args ...string) (string, error) {
	res, err := q.run(ctx, args...)
	if err != nil {
		return "", err
	}
	if res.ExitCode != 0 {
		return "", errExit
	}
	return strings.TrimSpace(string(res.Stdout)), nil
}

// headCommit returns the full HEAD commit id, or "" for an unborn branch.
func (q gitQuery) headCommit(ctx context.Context) (string, error) {
	res, err := q.run(ctx, "rev-parse", "--verify", "--quiet", "HEAD^{commit}")
	if err != nil {
		return "", err
	}
	switch res.ExitCode {
	case 0:
		return strings.TrimSpace(string(res.Stdout)), nil
	case 1:
		return "", nil
	}
	return "", fmt.Errorf("%w: rev-parse HEAD", ErrGitCommand)
}

func (q gitQuery) shallow(ctx context.Context) (bool, error) {
	out, err := q.text(ctx, "rev-parse", "--is-shallow-repository")
	if errors.Is(err, errExit) {
		return false, fmt.Errorf("%w: rev-parse shallow", ErrGitCommand)
	}
	return out == "true", err
}

// inputs lists relevant files, relative to the project root: tracked files and
// untracked files that Git does not ignore.
func (q gitQuery) inputs(ctx context.Context, scope InputScope) (tracked, untracked []string, err error) {
	if tracked, err = q.paths(ctx, scope, "ls-files", "-z", "--cached"); err != nil {
		return nil, nil, err
	}
	if untracked, err = q.paths(ctx, scope, "ls-files", "-z", "--others", "--exclude-standard"); err != nil {
		return nil, nil, err
	}
	trackedSet := make(map[string]bool, len(tracked))
	for _, p := range tracked {
		trackedSet[p] = true
	}
	untracked = slices.DeleteFunc(untracked, func(p string) bool { return trackedSet[p] })
	return tracked, untracked, nil
}

// trackedChanged reports whether any relevant tracked file differs from HEAD,
// staged or not. Deleted files count as changes.
func (q gitQuery) trackedChanged(ctx context.Context, scope InputScope) (bool, error) {
	changed, err := q.paths(ctx, scope, "diff", "--relative", "--name-only", "-z", "--no-renames", "HEAD", "--")
	return len(changed) > 0, err
}

// paths runs a NUL-separated path listing and keeps the relevant entries,
// sorted and de-duplicated.
func (q gitQuery) paths(ctx context.Context, scope InputScope, args ...string) ([]string, error) {
	res, err := q.run(ctx, args...)
	if err != nil {
		return nil, err
	}
	if res.ExitCode != 0 {
		return nil, fmt.Errorf("%w: %s", ErrGitCommand, args[0])
	}
	var out []string
	for _, raw := range strings.Split(string(res.Stdout), "\x00") {
		if raw == "" {
			continue
		}
		p := filepath.ToSlash(raw)
		if !scope.relevant(p) {
			continue
		}
		// A relevant input that cannot be fingerprinted safely must not vanish
		// from the identity; refuse to describe the repository instead.
		if !safeRelative(runtime.GOOS, p) {
			return nil, fieldError(ErrInputUnreadable, "input", "%q is not a safe project-relative path on %s", p, runtime.GOOS)
		}
		out = append(out, p)
	}
	slices.Sort(out)
	return slices.Compact(out), nil
}

// safeRelative refuses anything that could leave the project root. A colon is
// a drive or alternate-data-stream marker only on Windows; elsewhere it is an
// ordinary filename character.
func safeRelative(goos, p string) bool {
	if strings.HasPrefix(p, "/") || slices.Contains(splitPath(p), "..") {
		return false
	}
	return goos != "windows" || !strings.Contains(p, ":")
}

// fingerprintInputs hashes the working-tree bytes of every listed file. A file
// Git lists but that is gone from disk is skipped, so the fingerprint follows
// the work tree rather than the index. Entries are NUL-delimited and sorted, so
// neither enumeration order nor unusual filenames change the result. Only
// digests are kept; content and link targets are never retained.
func fingerprintInputs(ctx context.Context, root string, files []string) (string, error) {
	sorted := slices.Compact(slices.Sorted(slices.Values(files)))
	h := sha256.New()
	h.Write([]byte("codehealth-inputs-v1\x00"))
	for _, p := range sorted {
		if err := ctx.Err(); err != nil {
			return "", err
		}
		kind, digest, err := inputDigest(filepath.Join(root, filepath.FromSlash(p)))
		if err != nil {
			return "", fieldError(ErrInputUnreadable, "input", "%s", p)
		}
		if kind == "" {
			continue
		}
		h.Write([]byte(kind + "\x00" + strconv.Itoa(len(p)) + "\x00" + p + "\x00" + digest + "\x00"))
	}
	return digestPrefix + hex.EncodeToString(h.Sum(nil)), nil
}

// inputDigest returns the entry kind and content digest of one path, without
// following symlinks. Directories (submodules, nested repositories) and special
// files are not content inputs and return an empty kind.
func inputDigest(full string) (kind, digest string, err error) {
	info, err := os.Lstat(full)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return "", "", nil
	case err != nil:
		return "", "", err
	}
	switch {
	case info.Mode()&fs.ModeSymlink != 0:
		target, err := os.Readlink(full)
		if err != nil {
			return "", "", err
		}
		sum := sha256.Sum256([]byte(target))
		return "symlink", hex.EncodeToString(sum[:]), nil
	case info.Mode().IsRegular():
		f, err := os.Open(full)
		if err != nil {
			return "", "", err
		}
		defer f.Close()
		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			return "", "", err
		}
		return "file", hex.EncodeToString(h.Sum(nil)), nil
	}
	return "", "", nil
}

// RelationKind says how a recorded commit relates to the current one. Only
// ancestry Git can prove produces Behind, Ahead, or Diverged.
type RelationKind string

const (
	RelationSameCommit  RelationKind = "same_commit"
	RelationBehind      RelationKind = "behind"      // recorded commit is an ancestor of current
	RelationAhead       RelationKind = "ahead"       // current commit is an ancestor of recorded
	RelationDiverged    RelationKind = "diverged"    // neither is an ancestor of the other
	RelationUnknown     RelationKind = "unknown"     // history is incomplete (shallow or missing object)
	RelationNoCommit    RelationKind = "no_commit"   // either side has no useful commit
	RelationUnavailable RelationKind = "unavailable" // Git could not answer
)

// Relation compares a recorded repository identity with the current one.
// Counts are commits reachable from one side only and are meaningful only when
// Counted is true; a shallow clone can prove ancestry without proving distance.
type Relation struct {
	Kind         RelationKind
	RecordedOnly int
	CurrentOnly  int
	Counted      bool
	// InputsMatch is true when the configured-scope fingerprints are equal,
	// whatever the commits say.
	InputsMatch bool
}

// Relate compares recorded with current using actual commit ancestry from the
// local object database. It never fetches. Git failures yield Unavailable;
// only cancellation and timeout are returned as errors.
func Relate(ctx context.Context, runner CommandRunner, root string, recorded RepositoryIdentity, current Observation) (Relation, error) {
	rel := Relation{InputsMatch: recorded.InputFingerprint == current.Identity.InputFingerprint}
	switch {
	case recorded.Commit == "" || current.Identity.Commit == "":
		rel.Kind = RelationNoCommit
		return rel, nil
	case recorded.Commit == current.Identity.Commit:
		rel.Kind = RelationSameCommit
		return rel, nil
	}
	ctx, cancel := context.WithTimeout(ctx, DefaultGitTimeout)
	defer cancel()
	kind, err := relateCommits(ctx, gitQuery{runner: runner, dir: root}, recorded.Commit, current.Identity.Commit, current.Shallow)
	if err != nil {
		if ctx.Err() != nil {
			return Relation{}, ctx.Err()
		}
		rel.Kind = RelationUnavailable
		return rel, nil
	}
	rel.Kind = kind
	if kind == RelationUnknown || kind == RelationUnavailable {
		return rel, nil
	}
	if !current.Shallow {
		rel.RecordedOnly, rel.CurrentOnly, err = countSides(ctx, gitQuery{runner: runner, dir: root}, recorded.Commit, current.Identity.Commit)
		if err != nil {
			if ctx.Err() != nil {
				return Relation{}, ctx.Err()
			}
			rel.Kind = RelationUnavailable
			return rel, nil
		}
		rel.Counted = true
	}
	return rel, nil
}

// relateCommits decides the kind from ancestry. Without full history a
// negative ancestry answer proves nothing, so a shallow clone or an absent
// object yields Unknown rather than Diverged.
func relateCommits(ctx context.Context, q gitQuery, recorded, current string, shallow bool) (RelationKind, error) {
	for _, id := range []string{recorded, current} {
		res, err := q.run(ctx, "cat-file", "-e", id+"^{commit}")
		if err != nil {
			return "", err
		}
		if res.ExitCode != 0 {
			return RelationUnknown, nil
		}
	}
	recordedIsAncestor, err := q.isAncestor(ctx, recorded, current)
	if err != nil {
		return "", err
	}
	if recordedIsAncestor {
		return RelationBehind, nil
	}
	currentIsAncestor, err := q.isAncestor(ctx, current, recorded)
	if err != nil {
		return "", err
	}
	switch {
	case currentIsAncestor:
		return RelationAhead, nil
	case shallow:
		return RelationUnknown, nil
	}
	return RelationDiverged, nil
}

func (q gitQuery) isAncestor(ctx context.Context, ancestor, descendant string) (bool, error) {
	res, err := q.run(ctx, "merge-base", "--is-ancestor", ancestor, descendant)
	if err != nil {
		return false, err
	}
	switch res.ExitCode {
	case 0:
		return true, nil
	case 1:
		return false, nil
	}
	return false, fmt.Errorf("%w: merge-base", ErrGitCommand)
}

// countSides returns commits only on the recorded side and only on the current
// side of the symmetric difference.
func countSides(ctx context.Context, q gitQuery, recorded, current string) (recordedOnly, currentOnly int, err error) {
	out, err := q.text(ctx, "rev-list", "--left-right", "--count", recorded+"..."+current)
	if err != nil {
		return 0, 0, fmt.Errorf("%w: rev-list", ErrGitCommand)
	}
	left, right, ok := strings.Cut(out, "\t")
	if !ok {
		return 0, 0, fmt.Errorf("%w: rev-list output", ErrGitCommand)
	}
	if recordedOnly, err = strconv.Atoi(left); err != nil {
		return 0, 0, fmt.Errorf("%w: rev-list output", ErrGitCommand)
	}
	if currentOnly, err = strconv.Atoi(right); err != nil {
		return 0, 0, fmt.Errorf("%w: rev-list output", ErrGitCommand)
	}
	return recordedOnly, currentOnly, nil
}

// Describe renders the observation in plain words. It names only the abbreviated
// commit, never paths, remotes, authors, or messages.
func (o Observation) Describe() string {
	switch {
	case o.Unborn():
		return "No commits yet. Health is tied to the current files in scope only."
	case o.Identity.Dirty:
		return fmt.Sprintf("Uncommitted changes in scope on top of commit %s.", shortCommit(o.Identity.Commit))
	}
	return fmt.Sprintf("Clean at commit %s.", shortCommit(o.Identity.Commit))
}

// DescribeRepositoryError renders a failed observation in plain words.
func DescribeRepositoryError(err error) string {
	switch {
	case errors.Is(err, ErrGitMissing):
		return "Repository state is unavailable: Git is not installed."
	case errors.Is(err, ErrNotRepository):
		return "Repository state is unavailable: this is not a Git repository."
	case errors.Is(err, context.DeadlineExceeded):
		return "Repository state is unavailable: Git took too long."
	case errors.Is(err, context.Canceled):
		return "Repository state is unavailable: the check was cancelled."
	case errors.Is(err, ErrInputUnreadable):
		return "Repository state is unavailable: a file in scope could not be read."
	}
	return "Repository state is unavailable: Git could not answer."
}

// Describe renders the relation in plain words. A number of commits appears
// only when Git proved it.
func (r Relation) Describe() string {
	var s string
	switch r.Kind {
	case RelationSameCommit:
		s = "Measured at the current commit."
	case RelationBehind:
		s = fmt.Sprintf("Measured %s.", distance(r.Counted, r.CurrentOnly, "behind the current commit", "on an earlier commit than the current one"))
	case RelationAhead:
		s = fmt.Sprintf("Measured %s.", distance(r.Counted, r.RecordedOnly, "ahead of the current commit", "on a later commit than the current one"))
	case RelationDiverged:
		s = fmt.Sprintf("Measured on a diverged commit: %s only there, %s only here.", plural(r.RecordedOnly, "commit"), plural(r.CurrentOnly, "commit"))
	case RelationUnknown:
		s = "How the measured commit relates to the current one is unknown: history is incomplete."
	case RelationNoCommit:
		s = "Measured without a commit to compare against."
	default:
		s = "How the measured commit relates to the current one is unavailable."
	}
	if r.Kind != RelationSameCommit && r.Kind != RelationUnknown && r.Kind != RelationUnavailable && r.InputsMatch {
		s += " Files in scope are unchanged."
	} else if r.Kind == RelationSameCommit && !r.InputsMatch {
		s += " Files in scope have changed since."
	}
	return s
}

func distance(counted bool, n int, withCount, without string) string {
	if counted {
		return plural(n, "commit") + " " + withCount
	}
	return without
}

func plural(n int, noun string) string {
	if n == 1 {
		return "1 " + noun
	}
	return strconv.Itoa(n) + " " + noun + "s"
}

func shortCommit(c string) string {
	if len(c) > 7 {
		return c[:7]
	}
	return c
}
