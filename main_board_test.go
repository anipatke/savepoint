package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/testutil"
)

// writeBoardV2Project writes a valid V2 project at dir, the shape `savepoint
// init` has produced since E47 plus one Objective and one Task.
func writeBoardV2Project(t *testing.T, dir string) {
	t.Helper()
	savepointDir := filepath.Join(dir, ".savepoint")
	testutil.WriteFile(t, filepath.Join(savepointDir, "config.yml"), "schema_version: 2\n")
	testutil.WriteFile(t, filepath.Join(savepointDir, "router.md"),
		"# Router\n\n## Current state\n\n```yaml\nstate: task\nobjective: O001\ntask: T001\nnext_action: \"Build T001.\"\n```\n")
	testutil.WriteFile(t, filepath.Join(savepointDir, "objectives", "O001-first", "Objective.md"),
		"---\nid: O001\ntitle: \"First objective\"\nstatus: planned\n---\n\n# First objective\n")
	testutil.WriteFile(t, filepath.Join(savepointDir, "objectives", "O001-first", "tasks", "T001-alpha.md"),
		"---\nid: T001\ntitle: \"Do the thing\"\nobjective: O001\nplanned_by: {role: planner, session: board-fixture}\nstatus: planned\n---\n\n# Do the thing\n")
}

// runMainInDirForTest runs the built command in dir, so a command that reads the
// project from its working directory — board does — is exercised against a
// temporary project rather than this repository's own.
func runMainInDirForTest(t *testing.T, dir string, args []string) mainResult {
	t.Helper()

	self, err := filepath.Abs(os.Args[0])
	if err != nil {
		t.Fatalf("resolve test binary: %v", err)
	}

	cmdArgs := append([]string{"-test.run=TestMainHelperProcess", "--"}, args...)
	command := exec.Command(self, cmdArgs...)
	command.Dir = dir
	command.Env = append(os.Environ(), "SAVEPOINT_TEST_MAIN=1")

	stdout, err := command.Output()
	stderr := ""
	if exitErr, ok := err.(*exec.ExitError); ok {
		stderr = string(exitErr.Stderr)
	}
	return mainResult{stdout: string(stdout), stderr: stderr, err: err}
}

func TestMainBoardV2ProjectWithoutTTYReportsNextAndExitsZero(t *testing.T) {
	dir := t.TempDir()
	writeBoardV2Project(t, dir)
	before := snapshotDir(t, dir)

	result := runMainInDirForTest(t, dir, []string{"board"})

	if result.err != nil {
		t.Fatalf("savepoint board failed: %v\nstderr: %s", result.err, result.stderr)
	}
	if !strings.Contains(result.stdout, "NEXT:") || !strings.Contains(result.stdout, "Task: T001 — Do the thing") {
		t.Errorf("stdout = %q, want the resolved next action", result.stdout)
	}
	if !strings.Contains(result.stdout, "Objectives: 1  Tasks: 1") {
		t.Errorf("stdout = %q, want the loaded counts", result.stdout)
	}
	assertSameSnapshot(t, before, snapshotDir(t, dir))
}

// TestMainBoardV2LoadFailureExitsNonzeroWithDiagnosticOnStderr closes the
// load-diagnostic contract end to end: without a TTY the failure is named on
// stderr, the exit is nonzero, and no board is written to stdout.
func TestMainBoardV2LoadFailureExitsNonzeroWithDiagnosticOnStderr(t *testing.T) {
	dir := t.TempDir()
	writeBoardV2Project(t, dir)
	testutil.WriteFile(t, filepath.Join(dir, ".savepoint", "objectives", "O001-first", "tasks", "T001-alpha.md"),
		"---\nid: T001\nobjective: O001\nplanned_by: {role: planner, session: board-fixture}\nstatus: planned\n---\n\n# Untitled\n")
	before := snapshotDir(t, dir)

	result := runMainInDirForTest(t, dir, []string{"board"})

	if result.err == nil {
		t.Fatal("savepoint board over an unloadable V2 project succeeded, want a nonzero exit")
	}
	if !strings.Contains(result.stderr, "T001-alpha.md") || !strings.Contains(result.stderr, "missing required field title") {
		t.Fatalf("stderr = %q, want the file and the problem named", result.stderr)
	}
	if strings.Contains(result.stdout, "PLANNED") {
		t.Errorf("stdout = %q, want no partial board", result.stdout)
	}
	assertSameSnapshot(t, before, snapshotDir(t, dir))
}

func TestMainBoardRejectsV1FilterOnAV2Project(t *testing.T) {
	dir := t.TempDir()
	writeBoardV2Project(t, dir)

	result := runMainInDirForTest(t, dir, []string{"board", "--release", "v1"})

	if result.err == nil {
		t.Fatal("savepoint board --release over a V2 project succeeded, want a nonzero exit")
	}
	if !strings.Contains(result.stderr, "--release") || !strings.Contains(result.stderr, "schema_version 2") {
		t.Fatalf("stderr = %q, want the flag and the schema version named", result.stderr)
	}
}
