package main

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/opencode/savepoint/internal/migrate"
	"github.com/opencode/savepoint/internal/testutil"
)

// writeResumeV2Project writes a minimal, valid V2 project whose router
// selects a planned Task with no dependencies — the one rung every other
// resume fixture below varies from: a Task that ResolveTaskStart allows,
// reached through the ordinary schema/router/index path rather than the
// pending-migration shortcut.
func writeResumeV2Project(t *testing.T, root string) {
	t.Helper()
	savepointDir := filepath.Join(root, ".savepoint")
	testutil.WriteFile(t, filepath.Join(savepointDir, "config.yml"), "schema_version: 2\n")
	testutil.WriteFile(t, filepath.Join(savepointDir, "router.md"), resumeRouterV2Content("task", "O001", "T001", "Build T001."))
	testutil.WriteFile(t, filepath.Join(savepointDir, "objectives", "O001-first", "Objective.md"),
		"---\nid: O001\ntitle: \"First objective\"\nstatus: planned\n---\n\n# First objective\n")
	testutil.WriteFile(t, filepath.Join(savepointDir, "objectives", "O001-first", "tasks", "T001-alpha.md"),
		"---\nid: T001\ntitle: \"Do the thing\"\nobjective: O001\nplanned_by: {role: planner, session: planning-fixture}\nstatus: planned\n---\n\n# Do the thing\n")
}

func resumeRouterV2Content(state, objective, task, nextAction string) string {
	return "# Router\n\n## Current state\n\n```yaml\nstate: " + state + "\nobjective: " + objective + "\ntask: " + task + "\nnext_action: \"" + nextAction + "\"\n```\n"
}

func TestMainResumeV2ProjectRendersAndExitsZero(t *testing.T) {
	dir := t.TempDir()
	writeResumeV2Project(t, dir)
	before := snapshotDir(t, dir)

	result := runMainForTest(t, []string{"resume", dir}, "")

	if result.err != nil {
		t.Fatalf("savepoint resume failed: %v\nstderr: %s", result.err, result.stderr)
	}
	if !strings.Contains(result.stdout, "Task: T001 — Do the thing") {
		t.Errorf("stdout = %q, want the selected Task named", result.stdout)
	}
	if !strings.Contains(result.stdout, "Next action: Start Task T001.") {
		t.Errorf("stdout = %q, want the next action", result.stdout)
	}
	assertSameSnapshot(t, before, snapshotDir(t, dir))
}

func TestMainResumeV1ProjectPrintsMigrateRouteMessageAndExitsNonzero(t *testing.T) {
	dir := t.TempDir()
	writeMigrateMinimalProject(t, dir)
	before := snapshotDir(t, dir)

	result := runMainForTest(t, []string{"resume", dir}, "")

	if result.err == nil {
		t.Fatal("savepoint resume over a V1 project succeeded, want a nonzero exit")
	}
	if !strings.Contains(result.stdout, "schema_version 1") || !strings.Contains(result.stdout, "savepoint migrate") {
		t.Fatalf("stdout = %q, want the schema version and the migrate route named", result.stdout)
	}
	assertSameSnapshot(t, before, snapshotDir(t, dir))
}

func TestMainResumeMissingDirectory(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "does-not-exist")

	result := runMainForTest(t, []string{"resume", missing}, "")

	if result.err == nil {
		t.Fatal("savepoint resume over a missing directory succeeded, want a nonzero exit")
	}
	if !strings.Contains(result.stderr, "target directory does not exist") {
		t.Fatalf("stderr = %q, want a named missing-directory error", result.stderr)
	}
}

func TestMainResumeNotASavepointProject(t *testing.T) {
	dir := t.TempDir()
	before := snapshotDir(t, dir)

	result := runMainForTest(t, []string{"resume", dir}, "")

	if result.err == nil {
		t.Fatal("savepoint resume over a non-Savepoint directory succeeded, want a nonzero exit")
	}
	if !strings.Contains(result.stderr, "is not a Savepoint project") {
		t.Fatalf("stderr = %q, want a named not-a-project error", result.stderr)
	}
	assertSameSnapshot(t, before, snapshotDir(t, dir))
}

func TestMainResumeMalformedRouter(t *testing.T) {
	dir := t.TempDir()
	writeResumeV2Project(t, dir)
	testutil.WriteFile(t, filepath.Join(dir, ".savepoint", "router.md"), resumeRouterV2Content("bogus", "O001", "T001", "Build T001."))
	before := snapshotDir(t, dir)

	result := runMainForTest(t, []string{"resume", dir}, "")

	if result.err == nil {
		t.Fatal("savepoint resume over a malformed router succeeded, want a nonzero exit")
	}
	if !strings.Contains(result.stderr, "router state") {
		t.Fatalf("stderr = %q, want a named router diagnostic", result.stderr)
	}
	assertSameSnapshot(t, before, snapshotDir(t, dir))
}

func TestMainResumeIndexFailsToLoad(t *testing.T) {
	dir := t.TempDir()
	savepointDir := filepath.Join(dir, ".savepoint")
	testutil.WriteFile(t, filepath.Join(savepointDir, "config.yml"), "schema_version: 2\n")
	testutil.WriteFile(t, filepath.Join(savepointDir, "router.md"), resumeRouterV2Content("idea", "none", "none", ""))
	testutil.WriteFile(t, filepath.Join(savepointDir, "objectives", "O001-first", "Objective.md"),
		"---\nid: O001\ntitle: \"First objective\"\nstatus: planned\n---\n\n# First objective\n")
	// A Task missing its required title fails DecodeTaskV2, so the whole
	// index load fails closed rather than loading a project short one Task.
	testutil.WriteFile(t, filepath.Join(savepointDir, "objectives", "O001-first", "tasks", "T001-alpha.md"),
		"---\nid: T001\nobjective: O001\nplanned_by: {role: planner, session: planning-fixture}\nstatus: planned\n---\n\n# Untitled\n")
	before := snapshotDir(t, dir)

	result := runMainForTest(t, []string{"resume", dir}, "")

	if result.err == nil {
		t.Fatal("savepoint resume over an unloadable index succeeded, want a nonzero exit")
	}
	if !strings.Contains(result.stderr, "missing required field title") {
		t.Fatalf("stderr = %q, want the decode diagnostic named", result.stderr)
	}
	assertSameSnapshot(t, before, snapshotDir(t, dir))
}

func TestMainResumeUnresolvableSelectionStillExitsZero(t *testing.T) {
	dir := t.TempDir()
	writeResumeV2Project(t, dir)
	testutil.WriteFile(t, filepath.Join(dir, ".savepoint", "router.md"), resumeRouterV2Content("task", "O001", "T999", "Build T999."))
	before := snapshotDir(t, dir)

	result := runMainForTest(t, []string{"resume", dir}, "")

	if result.err != nil {
		t.Fatalf("savepoint resume over an unresolvable selection failed: %v\nstderr: %s", result.err, result.stderr)
	}
	if !strings.Contains(result.stdout, "Selection: The router names task T999") {
		t.Fatalf("stdout = %q, want the selection diagnostic named", result.stdout)
	}
	if !strings.Contains(result.stdout, "Next action:") {
		t.Fatalf("stdout = %q, want the available next action alongside the diagnostic", result.stdout)
	}
	assertSameSnapshot(t, before, snapshotDir(t, dir))
}

func TestMainResumePendingMigrationRendersRungAndLeavesOperationUnchanged(t *testing.T) {
	dir := t.TempDir()
	writeMigrateMinimalProject(t, dir)
	op, err := migrate.CreateOperation(dir, "op-test-resume-pending", nil, nil, time.Now())
	if err != nil {
		t.Fatalf("CreateOperation() error = %v", err)
	}
	before := snapshotDir(t, dir)

	result := runMainForTest(t, []string{"resume", dir}, "")

	if result.err != nil {
		t.Fatalf("savepoint resume over a pending migration failed: %v\nstderr: %s", result.err, result.stderr)
	}
	if !strings.Contains(result.stdout, "Migration:") || !strings.Contains(result.stdout, op.Journal.OperationID) {
		t.Fatalf("stdout = %q, want the migration rung naming the pending operation", result.stdout)
	}
	assertSameSnapshot(t, before, snapshotDir(t, dir))
}

// resumeFailingWriter always fails, so runResume's writer-error path is
// exercised directly rather than by trying to make a real stdout fail.
type resumeFailingWriter struct{}

func (resumeFailingWriter) Write([]byte) (int, error) { return 0, errors.New("disk full") }

// TestRunResumeWriterFailurePropagates proves a writer failure mid-render
// surfaces with context and a nonzero code, rather than being swallowed —
// checked against runResume directly since the production stdout used by the
// subprocess-based tests above cannot be made to fail on demand.
func TestRunResumeWriterFailurePropagates(t *testing.T) {
	dir := t.TempDir()
	writeResumeV2Project(t, dir)
	before := snapshotDir(t, dir)

	code, err := runResume(dir, resumeFailingWriter{})

	if err == nil {
		t.Fatal("runResume() error = nil, want the writer failure propagated")
	}
	if !strings.Contains(err.Error(), "disk full") {
		t.Fatalf("error = %q, want it to wrap the writer's own error", err.Error())
	}
	if code != 1 {
		t.Fatalf("code = %d, want 1", code)
	}
	assertSameSnapshot(t, before, snapshotDir(t, dir))
}
