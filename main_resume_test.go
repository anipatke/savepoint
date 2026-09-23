package main

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/testutil"
)

// writeResumeV2Project writes a minimal, valid V2 project whose router
// selects a planned Task with no dependencies — the one rung every other
// resume fixture below varies from: a Task that ResolveTaskStart allows,
// reached through the ordinary schema/router/index path.
func writeResumeV2Project(t *testing.T, root string) {
	t.Helper()
	savepointDir := filepath.Join(root, ".savepoint")
	testutil.WriteFile(t, filepath.Join(savepointDir, "config.yml"), "schema_version: 2\n")
	testutil.WriteFile(t, filepath.Join(savepointDir, "router.md"), resumeRouterV2Content("task", "O-001", "T-001", "Build T-001."))
	testutil.WriteFile(t, filepath.Join(savepointDir, "objectives", "O-001-first", "Objective.md"),
		"---\nid: O-001\ntitle: \"First objective\"\nstatus: planned\n---\n\n# First objective\n")
	testutil.WriteFile(t, filepath.Join(savepointDir, "objectives", "O-001-first", "tasks", "T-001-alpha.md"),
		"---\nid: T-001\ntitle: \"Do the thing\"\nobjective: O-001\nplanned_by: {role: planner, session: planning-fixture}\nstatus: planned\n---\n\n# Do the thing\n")
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
	if !strings.Contains(result.stdout, "Task: T-001 — Do the thing") {
		t.Errorf("stdout = %q, want the selected Task named", result.stdout)
	}
	if !strings.Contains(result.stdout, "Next action: Start Task T-001.") {
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
	testutil.WriteFile(t, filepath.Join(dir, ".savepoint", "router.md"), resumeRouterV2Content("bogus", "O-001", "T-001", "Build T-001."))
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
	testutil.WriteFile(t, filepath.Join(savepointDir, "objectives", "O-001-first", "Objective.md"),
		"---\nid: O-001\ntitle: \"First objective\"\nstatus: planned\n---\n\n# First objective\n")
	// A Task missing its required title fails DecodeTaskV2, so the whole
	// index load fails closed rather than loading a project short one Task.
	testutil.WriteFile(t, filepath.Join(savepointDir, "objectives", "O-001-first", "tasks", "T-001-alpha.md"),
		"---\nid: T-001\nobjective: O-001\nplanned_by: {role: planner, session: planning-fixture}\nstatus: planned\n---\n\n# Untitled\n")
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
	testutil.WriteFile(t, filepath.Join(dir, ".savepoint", "router.md"), resumeRouterV2Content("task", "O-001", "T-999", "Build T-999."))
	before := snapshotDir(t, dir)

	result := runMainForTest(t, []string{"resume", dir}, "")

	if result.err != nil {
		t.Fatalf("savepoint resume over an unresolvable selection failed: %v\nstderr: %s", result.err, result.stderr)
	}
	if !strings.Contains(result.stdout, "Selection: The router names task T-999") {
		t.Fatalf("stdout = %q, want the selection diagnostic named", result.stdout)
	}
	if !strings.Contains(result.stdout, "Next action:") {
		t.Fatalf("stdout = %q, want the available next action alongside the diagnostic", result.stdout)
	}
	assertSameSnapshot(t, before, snapshotDir(t, dir))
}

// TestMainResumeV1ProjectWithLegacyMigrationJournalGetsTheSameMigrateRouteMessage
// proves a leftover journal file does not change the generic schema-1 route.
func TestMainResumeV1ProjectWithLegacyMigrationJournalGetsTheSameMigrateRouteMessage(t *testing.T) {
	dir := t.TempDir()
	writeMigrateMinimalProject(t, dir)
	legacyDir := filepath.Join(dir, ".savepoint", ".migration", "op-test-resume-pending")
	testutil.WriteFile(t, filepath.Join(legacyDir, "operation.yml"), "legacy journal\n")
	before := snapshotDir(t, dir)

	result := runMainForTest(t, []string{"resume", dir}, "")

	if result.err == nil {
		t.Fatal("savepoint resume over a V1 project with a pending operation succeeded, want a nonzero exit")
	}
	if !strings.Contains(result.stdout, "schema_version 1") || !strings.Contains(result.stdout, "savepoint migrate") {
		t.Fatalf("stdout = %q, want the ordinary schema-1 migrate-route message", result.stdout)
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
