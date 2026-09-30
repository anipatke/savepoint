package codehealth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"testing"
)

// gitEnv isolates tests from the developer's Git configuration.
func gitEnv(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed; repository-state tests need it")
	}
	t.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
	t.Setenv("GIT_CONFIG_SYSTEM", os.DevNull)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
}

// git runs the real git as an independent oracle and returns trimmed stdout.
func git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	full := append([]string{"-c", "user.name=Test", "-c", "user.email=test@example.com", "-c", "commit.gpgsign=false", "-c", "init.defaultBranch=main", "-c", "core.autocrlf=false"}, args...)
	cmd := exec.Command("git", full...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return strings.TrimSpace(string(out))
}

func write(t *testing.T, dir, rel, content string) {
	t.Helper()
	full := filepath.Join(dir, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// newRepo makes a repository with one commit holding a.go and sub/b.go.
func newRepo(t *testing.T) string {
	t.Helper()
	gitEnv(t)
	dir := t.TempDir()
	git(t, dir, "init")
	write(t, dir, "a.go", "package a\n")
	write(t, dir, "sub/b.go", "package b\n")
	git(t, dir, "add", ".")
	git(t, dir, "commit", "-m", "first")
	return dir
}

func commitFile(t *testing.T, dir, rel, content string) string {
	t.Helper()
	write(t, dir, rel, content)
	git(t, dir, "add", ".")
	git(t, dir, "commit", "-m", "change "+rel)
	return git(t, dir, "rev-parse", "HEAD")
}

func observe(t *testing.T, dir string, scope InputScope) Observation {
	t.Helper()
	obs, err := ObserveRepository(context.Background(), GitRunner{}, dir, scope)
	if err != nil {
		t.Fatalf("ObserveRepository(%s): %v", dir, err)
	}
	if err := obs.Identity.validate("repository"); err != nil {
		t.Fatalf("identity fails the snapshot validation: %v", err)
	}
	return obs
}

// oracleFingerprint recomputes the documented fingerprint format straight from
// file contents, independent of Git and of the implementation's hashing path.
func oracleFingerprint(files map[string]string) string {
	names := make([]string, 0, len(files))
	for n := range files {
		names = append(names, n)
	}
	slices.Sort(names)
	var b strings.Builder
	b.WriteString("codehealth-inputs-v1\x00")
	for _, n := range names {
		sum := sha256.Sum256([]byte(files[n]))
		b.WriteString("file\x00" + strconv.Itoa(len(n)) + "\x00" + n + "\x00" + hex.EncodeToString(sum[:]) + "\x00")
	}
	sum := sha256.Sum256([]byte(b.String()))
	return "sha256:" + hex.EncodeToString(sum[:])
}

func TestObserveCleanMatchesOracles(t *testing.T) {
	dir := newRepo(t)
	obs := observe(t, dir, InputScope{})
	if want := git(t, dir, "rev-parse", "HEAD"); obs.Identity.Commit != want {
		t.Errorf("commit = %s, want %s", obs.Identity.Commit, want)
	}
	if obs.Identity.Dirty || obs.Shallow || obs.Unborn() {
		t.Errorf("clean repository reported %+v", obs)
	}
	want := oracleFingerprint(map[string]string{"a.go": "package a\n", "sub/b.go": "package b\n"})
	if obs.Identity.InputFingerprint != want {
		t.Errorf("fingerprint = %s, want %s", obs.Identity.InputFingerprint, want)
	}
	if got := obs.Describe(); got != "Clean at commit "+obs.Identity.Commit[:7]+"." {
		t.Errorf("Describe() = %q", got)
	}
}

func TestFingerprintIgnoresEnumerationOrderAndLocation(t *testing.T) {
	files := []string{"sub/b.go", "a.go", "z.txt"}
	for _, n := range files {
		write(t, t.TempDir(), n, "x") // keep helper exercised without sharing state
	}
	mk := func() string {
		d := t.TempDir()
		for _, n := range files {
			write(t, d, n, "content of "+n)
		}
		return d
	}
	d1, d2 := mk(), mk()
	ctx := context.Background()
	a, err := fingerprintInputs(ctx, d1, files)
	if err != nil {
		t.Fatal(err)
	}
	b, err := fingerprintInputs(ctx, d2, []string{"z.txt", "a.go", "sub/b.go", "a.go"})
	if err != nil {
		t.Fatal(err)
	}
	if a != b {
		t.Errorf("order or location changed the fingerprint: %s vs %s", a, b)
	}
}

func TestObserveDirtyTracked(t *testing.T) {
	dir := newRepo(t)
	clean := observe(t, dir, InputScope{})

	write(t, dir, "a.go", "package a // changed\n")
	dirty := observe(t, dir, InputScope{})
	if !dirty.Identity.Dirty || dirty.Identity.Commit != clean.Identity.Commit {
		t.Errorf("unstaged change: %+v", dirty.Identity)
	}
	if dirty.Identity.InputFingerprint == clean.Identity.InputFingerprint {
		t.Error("changed content kept the fingerprint")
	}
	if !strings.HasPrefix(dirty.Describe(), "Uncommitted changes in scope") {
		t.Errorf("Describe() = %q", dirty.Describe())
	}

	git(t, dir, "add", "a.go")
	staged := observe(t, dir, InputScope{})
	if !staged.Identity.Dirty || staged.Identity.InputFingerprint != dirty.Identity.InputFingerprint {
		t.Errorf("staging must not change state: %+v", staged.Identity)
	}

	write(t, dir, "a.go", "package a\n")
	git(t, dir, "add", "a.go")
	if back := observe(t, dir, InputScope{}); back.Identity != clean.Identity {
		t.Errorf("reverting should restore the identity: %+v vs %+v", back.Identity, clean.Identity)
	}
}

func TestObserveUntrackedRelevanceAndScope(t *testing.T) {
	dir := newRepo(t)
	clean := observe(t, dir, InputScope{})

	write(t, dir, "new.go", "package n\n")
	relevant := observe(t, dir, InputScope{})
	if !relevant.Identity.Dirty || relevant.Identity.InputFingerprint == clean.Identity.InputFingerprint {
		t.Errorf("relevant untracked file must change identity: %+v", relevant.Identity)
	}

	if got := observe(t, dir, InputScope{Exclude: []string{"new.go"}}); got.Identity.Dirty || got.Identity.InputFingerprint != clean.Identity.InputFingerprint {
		t.Errorf("excluded untracked file leaked into identity: %+v", got.Identity)
	}
	if got := observe(t, dir, InputScope{Include: []string{"sub"}}); got.Identity.Dirty {
		t.Errorf("file outside the include scope made the tree dirty: %+v", got.Identity)
	}

	write(t, dir, ".gitignore", "ignored.log\n")
	git(t, dir, "add", ".gitignore")
	git(t, dir, "commit", "-m", "ignore")
	base := observe(t, dir, InputScope{})
	write(t, dir, "ignored.log", "noise")
	if got := observe(t, dir, InputScope{}); got.Identity != base.Identity {
		t.Errorf("git-ignored file changed identity: %+v", got.Identity)
	}
}

func TestObserveExcludesStoredHealthSnapshots(t *testing.T) {
	dir := newRepo(t)
	clean := observe(t, dir, InputScope{})
	write(t, dir, ".savepoint/health/snapshots/s.json", "{}")
	if got := observe(t, dir, InputScope{}); got.Identity != clean.Identity {
		t.Errorf("writing a snapshot changed the state it describes: %+v", got.Identity)
	}
}

func TestObserveDeletedFile(t *testing.T) {
	dir := newRepo(t)
	clean := observe(t, dir, InputScope{})
	if err := os.Remove(filepath.Join(dir, "a.go")); err != nil {
		t.Fatal(err)
	}
	unstaged := observe(t, dir, InputScope{})
	if !unstaged.Identity.Dirty || unstaged.Identity.InputFingerprint == clean.Identity.InputFingerprint {
		t.Errorf("deletion must change identity: %+v", unstaged.Identity)
	}
	git(t, dir, "rm", "--cached", "-q", "a.go")
	if staged := observe(t, dir, InputScope{}); staged.Identity != unstaged.Identity {
		t.Errorf("staged and unstaged deletion differ: %+v vs %+v", staged.Identity, unstaged.Identity)
	}
}

func TestObserveSymlinkHashesTargetWithoutFollowing(t *testing.T) {
	dir := newRepo(t)
	link := filepath.Join(dir, "link")
	if err := os.Symlink("a.go", link); err != nil {
		t.Skipf("symlinks unavailable on this host: %v", err)
	}
	first := observe(t, dir, InputScope{})
	if !first.Identity.Dirty {
		t.Error("new untracked symlink should be a relevant change")
	}
	if err := os.Remove(link); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("sub/b.go", link); err != nil {
		t.Fatal(err)
	}
	if second := observe(t, dir, InputScope{}); second.Identity.InputFingerprint == first.Identity.InputFingerprint {
		t.Error("retargeting a symlink kept the fingerprint")
	}
	// A dangling link must not fail the observation.
	if err := os.Remove(link); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("missing-target", link); err != nil {
		t.Fatal(err)
	}
	observe(t, dir, InputScope{})
}

func TestObserveUnusualFilenames(t *testing.T) {
	dir := newRepo(t)
	names := []string{"with space.txt", "ünï-cödé.txt", "-leading-dash.txt", "quote'\"s.txt", "tab\there.txt"}
	if runtime.GOOS != "windows" {
		names = append(names, "new\nline.txt")
	}
	files := map[string]string{"a.go": "package a\n", "sub/b.go": "package b\n"}
	for _, n := range names {
		write(t, dir, n, "data "+n)
		files[n] = "data " + n
	}
	if runtime.GOOS == "windows" {
		delete(files, "tab\there.txt") // tab is not a legal Windows filename character
	}
	got := observe(t, dir, InputScope{})
	if want := oracleFingerprint(files); got.Identity.InputFingerprint != want {
		t.Errorf("fingerprint with unusual names = %s, want %s", got.Identity.InputFingerprint, want)
	}
}

func TestObserveFromSubdirectoryStaysInsideIt(t *testing.T) {
	dir := newRepo(t)
	root := filepath.Join(dir, "sub")
	base := observe(t, root, InputScope{})
	if want := oracleFingerprint(map[string]string{"b.go": "package b\n"}); base.Identity.InputFingerprint != want {
		t.Errorf("subdirectory fingerprint = %s, want %s", base.Identity.InputFingerprint, want)
	}
	write(t, dir, "a.go", "package a // outside the project root\n")
	if got := observe(t, root, InputScope{}); got.Identity != base.Identity {
		t.Errorf("a change outside the project root altered identity: %+v", got.Identity)
	}
	write(t, root, "b.go", "package b // inside\n")
	if got := observe(t, root, InputScope{}); !got.Identity.Dirty {
		t.Error("a change inside the project root was missed")
	}
}

func TestObserveLinkedWorktree(t *testing.T) {
	dir := newRepo(t)
	wt := filepath.Join(t.TempDir(), "wt")
	git(t, dir, "worktree", "add", "-q", "-b", "lane", wt)
	got := observe(t, wt, InputScope{})
	if got.Identity.Dirty || got.Identity.Commit != git(t, dir, "rev-parse", "HEAD") {
		t.Errorf("linked worktree: %+v", got.Identity)
	}
	if want := observe(t, dir, InputScope{}); got.Identity != want.Identity {
		t.Errorf("same content must give the same identity in any worktree: %+v vs %+v", got.Identity, want.Identity)
	}
	write(t, wt, "lane.go", "package lane\n")
	if !observe(t, wt, InputScope{}).Identity.Dirty || observe(t, dir, InputScope{}).Identity.Dirty {
		t.Error("worktree changes must stay in their own worktree")
	}
}

func TestObserveUnborn(t *testing.T) {
	gitEnv(t)
	dir := t.TempDir()
	git(t, dir, "init")
	empty := observe(t, dir, InputScope{})
	if empty.Identity.Commit != "" || !empty.Identity.Dirty || !empty.Unborn() {
		t.Errorf("unborn repository: %+v", empty)
	}
	if !strings.HasPrefix(empty.Describe(), "No commits yet") {
		t.Errorf("Describe() = %q", empty.Describe())
	}
	write(t, dir, "a.go", "package a\n")
	git(t, dir, "add", "a.go")
	if got := observe(t, dir, InputScope{}); got.Identity.InputFingerprint != oracleFingerprint(map[string]string{"a.go": "package a\n"}) {
		t.Errorf("unborn fingerprint = %s", got.Identity.InputFingerprint)
	}
}

func TestObserveLeaksNothingSensitive(t *testing.T) {
	dir := newRepo(t)
	git(t, dir, "remote", "add", "origin", "https://user:s3cr3t-token@private-host.example/team/repo.git")
	write(t, dir, "a.go", "package a // TOP-SECRET-CONTENT\n")
	obs := observe(t, dir, InputScope{})
	dump := fmt.Sprintf("%+v %s", obs, obs.Describe())
	for _, leak := range []string{"s3cr3t", "private-host", "TOP-SECRET", "Test", "test@example.com", "first", dir} {
		if strings.Contains(dump, leak) {
			t.Errorf("observation leaks %q: %s", leak, dump)
		}
	}
}

func TestObserveBoundaryFailures(t *testing.T) {
	gitEnv(t)
	t.Run("not a repository", func(t *testing.T) {
		parent := t.TempDir()
		t.Setenv("GIT_CEILING_DIRECTORIES", parent)
		dir := filepath.Join(parent, "plain")
		write(t, dir, "a.go", "x")
		_, err := ObserveRepository(context.Background(), GitRunner{}, dir, InputScope{})
		if !errors.Is(err, ErrNotRepository) {
			t.Fatalf("err = %v, want ErrNotRepository", err)
		}
		if got := DescribeRepositoryError(err); !strings.Contains(got, "not a Git repository") {
			t.Errorf("description = %q", got)
		}
	})
	t.Run("git missing", func(t *testing.T) {
		_, err := ObserveRepository(context.Background(), GitRunner{Executable: "no-such-git-binary"}, t.TempDir(), InputScope{})
		if !errors.Is(err, ErrGitMissing) {
			t.Fatalf("err = %v, want ErrGitMissing", err)
		}
	})
	t.Run("cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := ObserveRepository(ctx, GitRunner{}, newRepo(t), InputScope{})
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("err = %v, want context.Canceled", err)
		}
	})
	t.Run("invalid scope", func(t *testing.T) {
		_, err := ObserveRepository(context.Background(), GitRunner{}, newRepo(t), InputScope{Include: []string{"../escape"}})
		if !errors.Is(err, ErrInvalidConfig) {
			t.Fatalf("err = %v, want ErrInvalidConfig", err)
		}
	})
	t.Run("unreadable input", func(t *testing.T) {
		if runtime.GOOS == "windows" || os.Geteuid() == 0 {
			t.Skip("file permission bits are not enforced here")
		}
		dir := newRepo(t)
		if err := os.Chmod(filepath.Join(dir, "a.go"), 0o000); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = os.Chmod(filepath.Join(dir, "a.go"), 0o644) })
		_, err := ObserveRepository(context.Background(), GitRunner{}, dir, InputScope{})
		if !errors.Is(err, ErrInputUnreadable) {
			t.Fatalf("err = %v, want ErrInputUnreadable", err)
		}
	})
}

// fakeRunner scripts command results and records the argument vectors.
type fakeRunner struct {
	calls [][]string
	reply func(args []string) (CommandResult, error)
}

func (f *fakeRunner) Run(_ context.Context, _ string, args ...string) (CommandResult, error) {
	f.calls = append(f.calls, slices.Clone(args))
	return f.reply(args)
}

func TestRunnerBoundaryUsesArgumentVectors(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "my file; rm -rf.go", "x")
	f := &fakeRunner{reply: func(args []string) (CommandResult, error) {
		switch args[0] {
		case "rev-parse":
			switch args[1] {
			case "--is-inside-work-tree":
				return CommandResult{Stdout: []byte("true\n")}, nil
			case "--verify":
				return CommandResult{ExitCode: 1}, nil
			default:
				return CommandResult{Stdout: []byte("false\n")}, nil
			}
		case "ls-files":
			if slices.Contains(args, "--others") {
				return CommandResult{Stdout: []byte("my file; rm -rf.go\x00")}, nil
			}
			return CommandResult{}, nil
		}
		return CommandResult{}, errors.New("unexpected " + args[0])
	}}
	obs, err := ObserveRepository(context.Background(), f, dir, InputScope{})
	if err != nil {
		t.Fatal(err)
	}
	if !obs.Unborn() {
		t.Errorf("exit 1 from rev-parse --verify should mean unborn: %+v", obs)
	}
	for _, call := range f.calls {
		for _, a := range call {
			if strings.ContainsAny(a, " ;") {
				t.Errorf("argument %q looks shell-assembled in %v", a, call)
			}
		}
	}
}

func TestObserveRunnerFailuresAreNamed(t *testing.T) {
	cases := map[string]struct {
		reply func([]string) (CommandResult, error)
		want  error
	}{
		"runner error": {func([]string) (CommandResult, error) { return CommandResult{}, errors.New("boom /home/user/secret") }, ErrGitCommand},
		"timeout":      {func([]string) (CommandResult, error) { return CommandResult{}, context.DeadlineExceeded }, context.DeadlineExceeded},
		"head exit code": {func(a []string) (CommandResult, error) {
			if a[1] == "--is-inside-work-tree" {
				return CommandResult{Stdout: []byte("true")}, nil
			}
			return CommandResult{ExitCode: 128}, nil
		}, ErrGitCommand},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := ObserveRepository(context.Background(), &fakeRunner{reply: c.reply}, t.TempDir(), InputScope{})
			if !errors.Is(err, c.want) {
				t.Fatalf("err = %v, want %v", err, c.want)
			}
			if strings.Contains(err.Error(), "secret") {
				t.Errorf("error leaks runner detail: %v", err)
			}
		})
	}
}

func TestRelateAgainstGitOracles(t *testing.T) {
	dir := newRepo(t)
	c1 := git(t, dir, "rev-parse", "HEAD")
	commitFile(t, dir, "a.go", "package a // 2\n")
	c3 := commitFile(t, dir, "a.go", "package a // 3\n")
	git(t, dir, "branch", "side", c1)

	ctx := context.Background()
	current := observe(t, dir, InputScope{})
	ident := func(commit string) RepositoryIdentity {
		return RepositoryIdentity{Commit: commit, InputFingerprint: current.Identity.InputFingerprint}
	}

	t.Run("same", func(t *testing.T) {
		rel, err := Relate(ctx, GitRunner{}, dir, ident(c3), current)
		if err != nil || rel.Kind != RelationSameCommit || !rel.InputsMatch {
			t.Fatalf("got %+v, %v", rel, err)
		}
		if rel.Describe() != "Measured at the current commit." {
			t.Errorf("Describe() = %q", rel.Describe())
		}
	})
	t.Run("same commit, inputs changed", func(t *testing.T) {
		other := ident(c3)
		other.InputFingerprint = "sha256:" + strings.Repeat("0", 64)
		rel, _ := Relate(ctx, GitRunner{}, dir, other, current)
		if rel.InputsMatch || !strings.Contains(rel.Describe(), "have changed") {
			t.Errorf("got %+v %q", rel, rel.Describe())
		}
	})
	t.Run("behind counts commits", func(t *testing.T) {
		rel, err := Relate(ctx, GitRunner{}, dir, ident(c1), current)
		want := mustCount(t, dir, c1+".."+c3)
		if err != nil || rel.Kind != RelationBehind || !rel.Counted || rel.CurrentOnly != want || rel.RecordedOnly != 0 {
			t.Fatalf("got %+v, %v; want %d behind", rel, err, want)
		}
		if got := rel.Describe(); got != "Measured 2 commits behind the current commit. Files in scope are unchanged." {
			t.Errorf("Describe() = %q", got)
		}
	})
	t.Run("ahead", func(t *testing.T) {
		git(t, dir, "checkout", "-q", c1)
		t.Cleanup(func() { git(t, dir, "checkout", "-q", "main") })
		old := observe(t, dir, InputScope{})
		rel, err := Relate(ctx, GitRunner{}, dir, ident(c3), old)
		if err != nil || rel.Kind != RelationAhead || rel.RecordedOnly != mustCount(t, dir, c1+".."+c3) {
			t.Fatalf("got %+v, %v", rel, err)
		}
	})
	t.Run("diverged", func(t *testing.T) {
		git(t, dir, "checkout", "-q", "side")
		t.Cleanup(func() { git(t, dir, "checkout", "-q", "main") })
		side := commitFile(t, dir, "side.go", "package side\n")
		rel, err := Relate(ctx, GitRunner{}, dir, ident(c3), observe(t, dir, InputScope{}))
		if err != nil || rel.Kind != RelationDiverged || !rel.Counted || rel.RecordedOnly != 2 || rel.CurrentOnly != 1 {
			t.Fatalf("got %+v, %v (side %s)", rel, err, side[:7])
		}
		if !strings.Contains(rel.Describe(), "diverged") {
			t.Errorf("Describe() = %q", rel.Describe())
		}
	})
	t.Run("no useful commit", func(t *testing.T) {
		rel, _ := Relate(ctx, GitRunner{}, dir, RepositoryIdentity{Dirty: true, InputFingerprint: current.Identity.InputFingerprint}, current)
		if rel.Kind != RelationNoCommit {
			t.Errorf("got %+v", rel)
		}
	})
	t.Run("commit missing locally is unknown, never diverged", func(t *testing.T) {
		rel, err := Relate(ctx, GitRunner{}, dir, ident(strings.Repeat("a", 40)), current)
		if err != nil || rel.Kind != RelationUnknown || rel.Counted {
			t.Fatalf("got %+v, %v", rel, err)
		}
		if strings.Contains(rel.Describe(), "commits") {
			t.Errorf("unknown relation must not state a commit count: %q", rel.Describe())
		}
	})
}

func mustCount(t *testing.T, dir, rng string) int {
	t.Helper()
	n, err := strconv.Atoi(git(t, dir, "rev-list", "--count", rng))
	if err != nil {
		t.Fatal(err)
	}
	return n
}

func TestRelateInShallowClone(t *testing.T) {
	src := newRepo(t)
	c1 := git(t, src, "rev-parse", "HEAD")
	c2 := commitFile(t, src, "a.go", "package a // 2\n")
	c3 := commitFile(t, src, "a.go", "package a // 3\n")

	shallow := filepath.Join(t.TempDir(), "shallow")
	git(t, filepath.Dir(shallow), "clone", "-q", "--depth", "2", "file://"+filepath.ToSlash(src), shallow)
	current := observe(t, shallow, InputScope{})
	if !current.Shallow || current.Identity.Commit != c3 {
		t.Fatalf("clone setup: %+v", current)
	}
	ident := func(c string) RepositoryIdentity {
		return RepositoryIdentity{Commit: c, InputFingerprint: current.Identity.InputFingerprint}
	}

	rel, err := Relate(context.Background(), GitRunner{}, shallow, ident(c2), current)
	if err != nil || rel.Kind != RelationBehind || rel.Counted {
		t.Fatalf("proven ancestry without proven distance: %+v, %v", rel, err)
	}
	if got := rel.Describe(); strings.Contains(got, "1 commit") || !strings.Contains(got, "earlier commit") {
		t.Errorf("shallow wording must not invent a count: %q", got)
	}
	rel, err = Relate(context.Background(), GitRunner{}, shallow, ident(c1), current)
	if err != nil || rel.Kind != RelationUnknown {
		t.Fatalf("commit beyond the shallow boundary: %+v, %v", rel, err)
	}
}

func TestRelateFailurePaths(t *testing.T) {
	cur := Observation{Identity: RepositoryIdentity{Commit: strings.Repeat("b", 40), InputFingerprint: "sha256:" + strings.Repeat("1", 64)}}
	rec := RepositoryIdentity{Commit: strings.Repeat("a", 40), InputFingerprint: cur.Identity.InputFingerprint}

	failing := &fakeRunner{reply: func([]string) (CommandResult, error) { return CommandResult{}, ErrGitMissing }}
	rel, err := Relate(context.Background(), failing, t.TempDir(), rec, cur)
	if err != nil || rel.Kind != RelationUnavailable {
		t.Fatalf("git failure should be Unavailable, got %+v, %v", rel, err)
	}
	if !strings.Contains(rel.Describe(), "unavailable") {
		t.Errorf("Describe() = %q", rel.Describe())
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	cancelled := &fakeRunner{reply: func([]string) (CommandResult, error) { return CommandResult{}, ctx.Err() }}
	if _, err := Relate(ctx, cancelled, t.TempDir(), rec, cur); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation must surface as an error, got %v", err)
	}

	badCounts := &fakeRunner{reply: func(a []string) (CommandResult, error) {
		if a[0] == "rev-list" {
			return CommandResult{Stdout: []byte("not numbers")}, nil
		}
		return CommandResult{}, nil
	}}
	if rel, _ := Relate(context.Background(), badCounts, t.TempDir(), rec, cur); rel.Kind != RelationUnavailable || rel.Counted {
		t.Errorf("malformed rev-list output must not produce a count: %+v", rel)
	}
}

func TestDescriptionsAreDeterministicPlainWords(t *testing.T) {
	c := strings.Repeat("c", 40)
	tests := []struct {
		got, want string
	}{
		{Observation{Identity: RepositoryIdentity{Commit: c}}.Describe(), "Clean at commit ccccccc."},
		{Observation{Identity: RepositoryIdentity{Commit: c, Dirty: true}}.Describe(), "Uncommitted changes in scope on top of commit ccccccc."},
		{Relation{Kind: RelationBehind, CurrentOnly: 1, Counted: true}.Describe(), "Measured 1 commit behind the current commit."},
		{Relation{Kind: RelationAhead, RecordedOnly: 3, Counted: true}.Describe(), "Measured 3 commits ahead of the current commit."},
		{Relation{Kind: RelationDiverged, RecordedOnly: 1, CurrentOnly: 4, Counted: true}.Describe(), "Measured on a diverged commit: 1 commit only there, 4 commits only here."},
		{Relation{Kind: RelationUnknown}.Describe(), "How the measured commit relates to the current one is unknown: history is incomplete."},
		{Relation{Kind: RelationNoCommit}.Describe(), "Measured without a commit to compare against."},
		{Relation{Kind: RelationUnavailable}.Describe(), "How the measured commit relates to the current one is unavailable."},
		{DescribeRepositoryError(ErrGitMissing), "Repository state is unavailable: Git is not installed."},
		{DescribeRepositoryError(context.DeadlineExceeded), "Repository state is unavailable: Git took too long."},
		{DescribeRepositoryError(errors.New("anything else")), "Repository state is unavailable: Git could not answer."},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("got %q, want %q", tt.got, tt.want)
		}
	}
}

func TestInputScopeMatching(t *testing.T) {
	tests := []struct {
		scope InputScope
		path  string
		want  bool
	}{
		{InputScope{}, "a/b.go", true},
		{InputScope{}, ".savepoint/health/snapshots/x.json", false},
		{InputScope{}, ".savepoint/router.md", true},
		{InputScope{Include: []string{"internal"}}, "internal/x/y.go", true},
		{InputScope{Include: []string{"internal"}}, "cmd/y.go", false},
		{InputScope{Include: []string{"**/*.go"}}, "a/b/c.go", true},
		{InputScope{Include: []string{"**/*.go"}}, "c.go", true},
		{InputScope{Include: []string{"**/*.go"}}, "c.txt", false},
		{InputScope{Include: []string{"internal"}, Exclude: []string{"internal/gen"}}, "internal/gen/z.go", false},
		{InputScope{Exclude: []string{"*.md"}}, "README.md", false},
		{InputScope{Exclude: []string{"*.md"}}, "docs/README.md", true},
		{InputScope{Include: []string{"src"}}, "srcs/a.go", false},
	}
	for _, tt := range tests {
		if got := tt.scope.relevant(tt.path); got != tt.want {
			t.Errorf("%+v relevant(%q) = %v, want %v", tt.scope, tt.path, got, tt.want)
		}
	}
}
