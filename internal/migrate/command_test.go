package migrate

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestRunCommand_previewIsReadOnlyOnReadable0555Project(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("directory permissions are not enforced the same way on Windows")
	}
	if os.Geteuid() == 0 {
		t.Skip("running as root bypasses directory permissions")
	}

	root := copyFixtureProject(t, "v1-basic")
	if err := os.Chmod(root, 0555); err != nil {
		t.Fatalf("chmod project read-only: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(root, 0755) })

	before := snapshotTree(t, root)
	var preview strings.Builder
	probeCalls := 0
	oldProbe := targetWriteabilityProbe
	targetWriteabilityProbe = func(string) error {
		probeCalls++
		return nil
	}
	t.Cleanup(func() { targetWriteabilityProbe = oldProbe })

	code, err := RunCommand(CommandOptions{
		Dir:            root,
		Stdout:         &preview,
		Now:            fixedClock(time.Date(2026, 9, 19, 0, 0, 0, 0, time.UTC)),
		NewOperationID: fixedOperationID("preview"),
	})
	if err != nil || code != 0 {
		t.Fatalf("preview on a readable 0555 project: code = %d, err = %v\n%s", code, err, preview.String())
	}
	if !strings.Contains(preview.String(), "Migration preview") {
		t.Fatalf("preview output = %q, want the preview report", preview.String())
	}
	assertSnapshotsEqual(t, before, snapshotTree(t, root))
	if probeCalls != 0 {
		t.Fatalf("preview invoked the apply-only writeability probe %d time(s)", probeCalls)
	}
}

func TestRunCommand_applyWriteabilityProbePreservesPreExistingSentinelAndPrefix(t *testing.T) {
	root := copyFixtureProject(t, "v1-basic")
	sentinelPath := filepath.Join(root, ".savepoint-migrate-write-test")
	prefixPath := filepath.Join(root, ".savepoint-migrate-write-test-occupied")
	sentinel := []byte("pre-existing sentinel bytes\x00\n")
	prefix := []byte("pre-existing prefix collision bytes\xff\n")
	if err := os.WriteFile(sentinelPath, sentinel, 0644); err != nil {
		t.Fatalf("write sentinel: %v", err)
	}
	if err := os.WriteFile(prefixPath, prefix, 0644); err != nil {
		t.Fatalf("write prefix collision: %v", err)
	}

	var output strings.Builder
	code, err := RunCommand(CommandOptions{
		Dir:            root,
		Write:          true,
		Stdout:         &output,
		Now:            fixedClock(time.Date(2026, 9, 19, 12, 0, 0, 0, time.UTC)),
		NewOperationID: fixedOperationID("probe-collision"),
		RunGit:         cleanGitStub,
	})
	if err != nil || code != 0 {
		t.Fatalf("apply with pre-existing probe names: code = %d, err = %v\n%s", code, err, output.String())
	}
	if !strings.Contains(output.String(), "Undo from the project root with: git --literal-pathspecs restore") {
		t.Fatalf("apply output = %q, want the Git undo command", output.String())
	}

	for path, want := range map[string][]byte{sentinelPath: sentinel, prefixPath: prefix} {
		got, readErr := os.ReadFile(path)
		if readErr != nil {
			t.Fatalf("read preserved probe collision %s: %v", path, readErr)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("probe collision %s changed: got %q, want %q", path, got, want)
		}
	}
	matches, err := filepath.Glob(filepath.Join(root, ".savepoint-migrate-write-test-*"))
	if err != nil {
		t.Fatalf("glob probe files: %v", err)
	}
	if len(matches) != 1 || matches[0] != prefixPath {
		t.Fatalf("probe files after apply = %v, want only pre-existing %s", matches, prefixPath)
	}
}

func TestRunCommand_applyGitRefusalsHappenBeforeAnyWrite(t *testing.T) {
	cases := []struct {
		name     string
		runGit   GitCommand
		wantErr  error
		wantText string
	}{
		{
			name: "git missing",
			runGit: func(string, ...string) (string, error) {
				return "", ErrGitUnavailable
			},
			wantErr:  ErrGitUnavailable,
			wantText: "install git",
		},
		{
			name: "outside work tree",
			runGit: func(string, ...string) (string, error) {
				return "fatal: not a git repository", errors.New("exit status 128")
			},
			wantErr:  ErrNotGitWorkTree,
			wantText: "git init",
		},
		{
			name: "dirty migration path",
			runGit: func(_ string, args ...string) (string, error) {
				if args[0] == "rev-parse" {
					return "true\n", nil
				}
				return " M .savepoint/config.yml\n", nil
			},
			wantErr:  ErrDirtyGitPaths,
			wantText: "stage and commit or stash",
		},
		{
			name: "untracked migration source",
			runGit: func(_ string, args ...string) (string, error) {
				if args[0] == "rev-parse" {
					return "true\n", nil
				}
				return "?? .savepoint/releases/v1/epics/E01-example/E01-Detail.md\n", nil
			},
			wantErr:  ErrDirtyGitPaths,
			wantText: "stage and commit or stash",
		},
		{
			name: "ignored migration path",
			runGit: func(_ string, args ...string) (string, error) {
				if args[0] == "rev-parse" {
					return "true\n", nil
				}
				return "!! .savepoint/config.yml\n", nil
			},
			wantErr:  ErrDirtyGitPaths,
			wantText: "move ignored files",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := copyFixtureProject(t, "v1-basic")
			before := snapshotTree(t, root)
			probeCalls := 0
			oldProbe := targetWriteabilityProbe
			targetWriteabilityProbe = func(string) error {
				probeCalls++
				return nil
			}
			t.Cleanup(func() { targetWriteabilityProbe = oldProbe })

			code, err := RunCommand(CommandOptions{Dir: root, Write: true, RunGit: tc.runGit})
			if code != 1 || !errors.Is(err, tc.wantErr) {
				t.Fatalf("RunCommand() = %d, %v; want exit 1 with %v", code, err, tc.wantErr)
			}
			if !strings.Contains(err.Error(), tc.wantText) {
				t.Fatalf("RunCommand() error = %q, want next step %q", err, tc.wantText)
			}
			if probeCalls != 0 {
				t.Fatalf("writeability probe ran %d time(s) before Git refusal", probeCalls)
			}
			assertSnapshotsEqual(t, before, snapshotTree(t, root))
		})
	}
}

func TestRunCommand_secondApplyReportsAlreadyMigratedWithoutGit(t *testing.T) {
	root := copyFixtureProject(t, "v1-basic")
	if _, err := Apply(root, mustPlan(t, root)); err != nil {
		t.Fatalf("initial Apply() error = %v", err)
	}
	before := snapshotTree(t, root)
	gitCalls := 0
	var output strings.Builder
	code, err := RunCommand(CommandOptions{
		Dir:    root,
		Write:  true,
		Stdout: &output,
		RunGit: func(string, ...string) (string, error) {
			gitCalls++
			return "", ErrGitUnavailable
		},
	})
	if err != nil || code != 0 {
		t.Fatalf("second --apply = %d, %v\n%s", code, err, output.String())
	}
	if !strings.Contains(output.String(), "already migrated") {
		t.Fatalf("second --apply output = %q, want already-migrated message", output.String())
	}
	if gitCalls != 0 {
		t.Fatalf("second --apply called git %d times, want a no-op", gitCalls)
	}
	assertSnapshotsEqual(t, before, snapshotTree(t, root))
}

func TestRunCommand_realGitApplyIsReversibleAndRejectsDirtyPaths(t *testing.T) {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git is not available; real Git migration integration test skipped")
	}
	root := copyFixtureProject(t, "v1-history")
	initGitRepository(t, gitPath, root)
	before := snapshotTree(t, root)

	var preview strings.Builder
	code, err := RunCommand(CommandOptions{Dir: root, Stdout: &preview})
	if err != nil || code != 0 || !strings.Contains(preview.String(), "Migration preview") {
		t.Fatalf("preview in a clean Git project: code = %d, err = %v\n%s", code, err, preview.String())
	}
	assertSnapshotsEqual(t, before, snapshotTree(t, root))

	plan := mustPlan(t, root)
	var output strings.Builder
	code, err = RunCommand(CommandOptions{Dir: root, Write: true, Stdout: &output})
	if err != nil || code != 0 {
		t.Fatalf("apply in a clean Git project: code = %d, err = %v\n%s", code, err, output.String())
	}
	if !strings.Contains(output.String(), gitUndoCommand(plan)) {
		t.Fatalf("apply output omitted the undo command %q:\n%s", gitUndoCommand(plan), output.String())
	}
	wantChanged := make(map[string]bool)
	for _, path := range plannedGitPaths(plan) {
		wantChanged[path] = true
	}
	gotChanged := gitStatusPathsFor(t, gitPath, root, plannedGitPaths(plan))
	if len(gotChanged) != len(wantChanged) {
		t.Fatalf("Git changed paths = %v, want exactly planned paths %v", gotChanged, wantChanged)
	}
	for path := range wantChanged {
		if !gotChanged[path] {
			t.Errorf("Git status omitted planned change %s", path)
		}
	}

	tracked, untracked := gitUndoPathGroups(plan)
	runGitTest(t, gitPath, root, append([]string{"--literal-pathspecs", "restore", "--source=HEAD", "--staged", "--worktree", "--"}, tracked...)...)
	runGitTest(t, gitPath, root, append([]string{"--literal-pathspecs", "clean", "-fdx", "--"}, untracked...)...)
	if remaining := gitStatusPaths(t, gitPath, root); len(remaining) != 0 {
		t.Fatalf("Git status after printed undo sequence = %v, want clean", remaining)
	}
	assertProjectContentEqual(t, before, snapshotTree(t, root))

	configPath := filepath.Join(root, ".savepoint", "config.yml")
	config, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configPath, append(config, []byte("# owner edit\n")...), 0644); err != nil {
		t.Fatal(err)
	}
	dirtyBefore := snapshotTree(t, root)
	code, err = RunCommand(CommandOptions{Dir: root, Write: true, Stdout: &output})
	if code != 1 || !errors.Is(err, ErrDirtyGitPaths) {
		t.Fatalf("apply with an uncommitted .savepoint edit = %d, %v; want dirty-path refusal", code, err)
	}
	assertProjectContentEqual(t, dirtyBefore, snapshotTree(t, root))
}

func cleanGitStub(_ string, args ...string) (string, error) {
	if len(args) > 0 && args[0] == "rev-parse" {
		return "true\n", nil
	}
	return "", nil
}

func initGitRepository(t *testing.T, gitPath, root string) {
	t.Helper()
	runGitTest(t, gitPath, root, "init", "-q")
	runGitTest(t, gitPath, root, "config", "user.name", "Savepoint Test")
	runGitTest(t, gitPath, root, "config", "user.email", "savepoint-test@example.invalid")
	runGitTest(t, gitPath, root, "add", "-Af")
	runGitTest(t, gitPath, root, "commit", "-qm", "fixture baseline")
}

func runGitTest(t *testing.T, gitPath, root string, args ...string) string {
	t.Helper()
	command := exec.Command(gitPath, args...)
	command.Dir = root
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, output)
	}
	return string(output)
}

func gitStatusPaths(t *testing.T, gitPath, root string) map[string]bool {
	t.Helper()
	status := runGitTest(t, gitPath, root, "status", "--porcelain", "--untracked-files=all", "--no-renames")
	return parseGitStatusPaths(t, status)
}

func gitStatusPathsFor(t *testing.T, gitPath, root string, paths []string) map[string]bool {
	t.Helper()
	changed := make(map[string]bool)
	for _, path := range paths {
		status := runGitTest(t, gitPath, root, "--literal-pathspecs", "status", "--porcelain", "--untracked-files=all", "--ignored=matching", "--no-renames", "--", path)
		if strings.TrimSpace(status) != "" {
			changed[path] = true
		}
	}
	return changed
}

func parseGitStatusPaths(t *testing.T, status string) map[string]bool {
	t.Helper()
	paths := make(map[string]bool)
	for _, line := range strings.Split(strings.TrimSpace(status), "\n") {
		if line == "" {
			continue
		}
		if len(line) < 4 {
			t.Fatalf("unexpected git status line %q", line)
		}
		paths[line[3:]] = true
	}
	return paths
}

func assertProjectContentEqual(t *testing.T, before, after treeSnapshot) {
	t.Helper()
	projectFiles := func(snapshot treeSnapshot) map[string][]byte {
		files := make(map[string][]byte, len(snapshot))
		for path, item := range snapshot {
			if path == ".git" || strings.HasPrefix(path, ".git"+string(filepath.Separator)) {
				continue
			}
			files[path] = item.content
		}
		return files
	}
	want, got := projectFiles(before), projectFiles(after)
	if len(want) != len(got) {
		t.Fatalf("project file count changed: before %d, after %d", len(want), len(got))
	}
	for path, content := range want {
		if gotContent, ok := got[path]; !ok {
			t.Errorf("%s was present before and missing after", path)
		} else if !bytes.Equal(content, gotContent) {
			t.Errorf("%s content changed", path)
		}
	}
}
