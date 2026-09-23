package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/opencode/savepoint/internal/migrate"
	"github.com/opencode/savepoint/internal/testutil"
)

// writeBoardV2Project writes a valid V2 project at dir, the shape `savepoint
// init` has produced since E47 plus one Objective and one Task.
func writeBoardV2Project(t *testing.T, dir string) {
	t.Helper()
	savepointDir := filepath.Join(dir, ".savepoint")
	testutil.WriteFile(t, filepath.Join(savepointDir, "config.yml"), "schema_version: 2\n")
	testutil.WriteFile(t, filepath.Join(savepointDir, "router.md"),
		"# Router\n\n## Current state\n\n```yaml\nstate: task\nobjective: O-001\ntask: T-001\nnext_action: \"Build T-001.\"\n```\n")
	testutil.WriteFile(t, filepath.Join(savepointDir, "objectives", "O-001-first", "Objective.md"),
		"---\nid: O-001\ntitle: \"First objective\"\nstatus: planned\n---\n\n# First objective\n")
	testutil.WriteFile(t, filepath.Join(savepointDir, "objectives", "O-001-first", "tasks", "T-001-alpha.md"),
		"---\nid: T-001\ntitle: \"Do the thing\"\nobjective: O-001\nplanned_by: {role: planner, session: board-fixture}\nstatus: planned\n---\n\n# Do the thing\n")
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
	if !strings.Contains(result.stdout, "Planned T-001 — Do the thing") {
		t.Errorf("stdout = %q, want the resolved next action", result.stdout)
	}
	if !strings.Contains(result.stdout, "Objectives: 1  Tasks: 1") {
		t.Errorf("stdout = %q, want the loaded counts", result.stdout)
	}
	assertSameSnapshot(t, before, snapshotDir(t, dir))
}

func TestMainBareStartupUsesTheV2BoardPath(t *testing.T) {
	dir := t.TempDir()
	writeBoardV2Project(t, dir)

	result := runMainInDirForTest(t, dir, nil)

	if result.err != nil {
		t.Fatalf("bare startup failed: %v\nstderr: %s", result.err, result.stderr)
	}
	if !strings.Contains(result.stdout, "Planned T-001 — Do the thing") {
		t.Fatalf("stdout = %q, want the V2 Next projection", result.stdout)
	}
}

// TestMainBoardV2LoadFailureExitsNonzeroWithDiagnosticOnStderr closes the
// load-diagnostic contract end to end: without a TTY the failure is named on
// stderr, the exit is nonzero, and no board is written to stdout.
func TestMainBoardV2LoadFailureExitsNonzeroWithDiagnosticOnStderr(t *testing.T) {
	dir := t.TempDir()
	writeBoardV2Project(t, dir)
	testutil.WriteFile(t, filepath.Join(dir, ".savepoint", "objectives", "O-001-first", "tasks", "T-001-alpha.md"),
		"---\nid: T-001\nobjective: O-001\nplanned_by: {role: planner, session: board-fixture}\nstatus: planned\n---\n\n# Untitled\n")
	before := snapshotDir(t, dir)

	result := runMainInDirForTest(t, dir, []string{"board"})

	if result.err == nil {
		t.Fatal("savepoint board over an unloadable V2 project succeeded, want a nonzero exit")
	}
	if !strings.Contains(result.stderr, "T-001-alpha.md") || !strings.Contains(result.stderr, "missing required field title") {
		t.Fatalf("stderr = %q, want the file and the problem named", result.stderr)
	}
	if strings.Contains(result.stdout, "PLANNED") {
		t.Errorf("stdout = %q, want no partial board", result.stdout)
	}
	assertSameSnapshot(t, before, snapshotDir(t, dir))
}

func TestMainBoardRejectsLegacyFilterBeforeProjectRouting(t *testing.T) {
	dir := t.TempDir()
	writeBoardV2Project(t, dir)

	result := runMainInDirForTest(t, dir, []string{"board", "--release", "v1"})

	if result.err == nil {
		t.Fatal("savepoint board --release over a V2 project succeeded, want a nonzero exit")
	}
	if !strings.Contains(result.stderr, "--release") || !strings.Contains(result.stderr, "unknown board flag") {
		t.Fatalf("stderr = %q, want the legacy flag rejected by the V2-only parser", result.stderr)
	}
}

func TestMainBoardV1ProjectPrintsMigrationPreviewRouteAndDoesNotRender(t *testing.T) {
	dir := t.TempDir()
	writeMigrateMinimalProject(t, dir)
	before := snapshotDir(t, dir)

	result := runMainInDirForTest(t, dir, []string{"board"})

	if result.err == nil {
		t.Fatal("savepoint board over a V1 project exited zero, want migration refusal")
	}
	if !strings.Contains(result.stderr, "schema_version 1") || !strings.Contains(result.stderr, "migrate --dry-run") {
		t.Fatalf("stderr = %q, want the migration preview route", result.stderr)
	}
	if result.stdout != "" {
		t.Fatalf("stdout = %q, want no V1 board rendered", result.stdout)
	}
	assertSameSnapshot(t, before, snapshotDir(t, dir))
}

func TestMainDoctorV1ProjectRefusesBeforeChecks(t *testing.T) {
	dir := t.TempDir()
	writeMigrateMinimalProject(t, dir)
	before := snapshotDir(t, dir)

	result := runMainInDirForTest(t, dir, []string{"doctor"})

	if result.err == nil {
		t.Fatal("savepoint doctor over a V1 project exited zero, want migration refusal")
	}
	if !strings.Contains(result.stderr, "schema_version 1") || !strings.Contains(result.stderr, "migrate --dry-run") {
		t.Fatalf("stderr = %q, want the migration preview route", result.stderr)
	}
	if result.stdout != "" {
		t.Fatalf("stdout = %q, want no V1 quality-gate report", result.stdout)
	}
	assertSameSnapshot(t, before, snapshotDir(t, dir))
}

func TestMainDoctorPendingMigrationPrintsRecoveryGuidance(t *testing.T) {
	dir := t.TempDir()
	writeBoardV2Project(t, dir)
	if _, err := migrate.CreateOperation(dir, "op-doctor-pending", nil, nil, time.Now()); err != nil {
		t.Fatalf("CreateOperation() error = %v", err)
	}

	result := runMainInDirForTest(t, dir, []string{"doctor"})

	if result.err == nil {
		t.Fatal("savepoint doctor with a pending migration exited zero, want recovery refusal")
	}
	if !strings.Contains(result.stderr, "op-doctor-pending") || !strings.Contains(result.stderr, "migrate --recover") {
		t.Fatalf("stderr = %q, want recovery guidance", result.stderr)
	}
	if result.stdout != "" {
		t.Fatalf("stdout = %q, want no quality-gate report", result.stdout)
	}
}
