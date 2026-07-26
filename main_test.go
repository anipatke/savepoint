package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestMainVersionFlagPrintsVersion(t *testing.T) {
	result := runMainForTest(t, []string{"--version"}, "v9.8.7-test")

	if result.err != nil {
		t.Fatalf("savepoint --version failed: %v\nstderr: %s", result.err, result.stderr)
	}
	if strings.TrimSpace(result.stdout) != "v9.8.7-test" {
		t.Fatalf("stdout = %q, want version only", result.stdout)
	}
	if result.stderr != "" {
		t.Fatalf("stderr = %q, want empty", result.stderr)
	}
}

func TestMainInitHelpStillUsesNormalDispatch(t *testing.T) {
	result := runMainForTest(t, []string{"init", "--help"}, "")

	if result.err != nil {
		t.Fatalf("savepoint init --help failed: %v\nstderr: %s", result.err, result.stderr)
	}
	if !strings.Contains(result.stdout, "Usage: init [dir] [--force] [--install]") {
		t.Fatalf("stdout = %q, want init usage", result.stdout)
	}
}

func TestMainUpgradeAssetsPrintsPartialWorkOnFailure(t *testing.T) {
	// A write failure part-way through must not hide what was already applied:
	// the user needs the report to know which files changed and where any
	// backup went.
	if runtime.GOOS == "windows" {
		t.Skip("directory permissions are not enforced the same way on Windows")
	}
	if os.Geteuid() == 0 {
		t.Skip("running as root bypasses directory permissions")
	}

	dir := t.TempDir()
	mkdirAll(t, filepath.Join(dir, ".savepoint"))

	// A stale skill in a directory that cannot be written: the walk reaches it
	// after it has already installed earlier skills.
	blocked := filepath.Join(dir, "agent-skills", "savepoint-audit-epic")
	mkdirAll(t, blocked)
	if err := os.WriteFile(filepath.Join(blocked, "SKILL.md"), []byte("# Stale\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(blocked, 0500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(blocked, 0755) })

	result := runMainForTest(t, []string{"upgrade-assets", dir}, "")

	if result.err == nil {
		t.Fatal("upgrade-assets succeeded, want the blocked write to fail")
	}
	if !strings.Contains(result.stdout, "Upgrade Report:") {
		t.Errorf("stdout = %q, want the partial report", result.stdout)
	}
	if !strings.Contains(result.stdout, "failed  agent-skills/savepoint-audit-epic/SKILL.md") {
		t.Errorf("stdout = %q, want the failed path named", result.stdout)
	}
	if !strings.Contains(result.stdout, "agent-skills/references/audit-method.md") {
		t.Errorf("stdout = %q, want the already-applied work named", result.stdout)
	}
}

func mkdirAll(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
}

func TestMainHelperProcess(t *testing.T) {
	if os.Getenv("SAVEPOINT_TEST_MAIN") != "1" {
		return
	}
	if value := os.Getenv("SAVEPOINT_TEST_VERSION"); value != "" {
		version = value
	}
	os.Args = append([]string{"savepoint"}, helperArgs(os.Args)...)
	main()
}

type mainResult struct {
	stdout string
	stderr string
	err    error
}

func runMainForTest(t *testing.T, args []string, testVersion string) mainResult {
	t.Helper()

	cmdArgs := []string{"-test.run=TestMainHelperProcess", "--"}
	cmdArgs = append(cmdArgs, args...)
	cmd := exec.Command(os.Args[0], cmdArgs...)
	cmd.Env = append(os.Environ(), "SAVEPOINT_TEST_MAIN=1")
	if testVersion != "" {
		cmd.Env = append(cmd.Env, "SAVEPOINT_TEST_VERSION="+testVersion)
	}

	stdout, err := cmd.Output()
	stderr := ""
	if exitErr, ok := err.(*exec.ExitError); ok {
		stderr = string(exitErr.Stderr)
	}
	return mainResult{
		stdout: string(stdout),
		stderr: stderr,
		err:    err,
	}
}

func helperArgs(args []string) []string {
	for i, arg := range args {
		if arg == "--" {
			return args[i+1:]
		}
	}
	return nil
}
