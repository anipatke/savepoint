package main

import (
	"os"
	"os/exec"
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
