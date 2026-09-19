package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/opencode/savepoint/internal/data"
	savepointinit "github.com/opencode/savepoint/internal/init"
	"github.com/opencode/savepoint/internal/migrate"
	"github.com/opencode/savepoint/internal/resume"
	"github.com/opencode/savepoint/internal/testutil"
)

// This file proves the three claims T001-T006 each leave unproven on their
// own (E48 T007): every reachable project state lands on exactly one rung of
// the ladder, resume's no-write guarantee holds across the whole matrix and
// a second consecutive invocation, and the projection is a shared value with
// no resume-specific shape — a second, independent consumer built here reads
// the same selection and next action straight off it.

// matrixCase is one named project state the ladder must resolve to exactly
// one rung. build writes a real project tree under dir (dir/.savepoint/...,
// matching writeResumeV2Project's shape); wantKind is the single NextKind
// the whole pipeline — load, router read, projection, render — must reach.
type matrixCase struct {
	name       string
	build      func(t *testing.T, dir string)
	wantKind   data.NextKind
	wantAction string // substring resume's rendered "Next action:" line must contain
}

// resumeMatrixCases covers every rung of the precedence ladder documented on
// data.NextKind with a real, on-disk V2 project — not the in-memory V2Index
// literals internal/data/next_test.go builds for its own unit tests. Each
// case is independent and deliberately minimal: one project state, one
// expected rung.
func resumeMatrixCases() []matrixCase {
	return []matrixCase{
		{
			name:       "pending migration outranks everything",
			build:      matrixBuildPendingMigration,
			wantKind:   data.NextPendingMigration,
			wantAction: "Wait for the pending migration to finish",
		},
		{
			name:       "recorded replan blocks a planned task",
			build:      matrixBuildReplan,
			wantKind:   data.NextReplan,
			wantAction: "Resolve the recorded replan",
		},
		{
			name:       "unsatisfied task dependency blocks start",
			build:      matrixBuildTaskDependency,
			wantKind:   data.NextDependency,
			wantAction: "Wait on the named dependency",
		},
		{
			name:       "planned task with no blockers may start",
			build:      matrixBuildExecute,
			wantKind:   data.NextExecute,
			wantAction: "Start Task T001.",
		},
		{
			name:       "task at audit with no recorded check needs one",
			build:      matrixBuildCheckNeeded,
			wantKind:   data.NextCheckNeeded,
			wantAction: "Record a fresh Check",
		},
		{
			name:       "current clearance still needs owner acceptance",
			build:      matrixBuildOwnerValidation,
			wantKind:   data.NextOwnerValidationRequired,
			wantAction: "Ask the owner to accept",
		},
		{
			name:       "every owned task done, objective integration check missing",
			build:      matrixBuildObjectiveIntegration,
			wantKind:   data.NextObjectiveIntegration,
			wantAction: "Record the Objective O001 integration Check",
		},
		{
			name:       "no selection, a ready task exists elsewhere",
			build:      matrixBuildReady,
			wantKind:   data.NextReady,
			wantAction: "Start Task T001.",
		},
		{
			name:       "objective selected with no tasks yet falls through to planning",
			build:      matrixBuildObjectiveNoTasks,
			wantKind:   data.NextReady,
			wantAction: "Plan Tasks under Objective O001.",
		},
		{
			name:       "objective integration check outstanding but router selects nothing",
			build:      matrixBuildObjectiveIntegrationUnselected,
			wantKind:   data.NextObjectiveIntegration,
			wantAction: "Record the Objective O001 integration Check",
		},
		{
			name:       "empty project, nothing ready",
			build:      matrixBuildPlanObjective,
			wantKind:   data.NextPlanObjective,
			wantAction: "Plan the next Objective",
		},
		{
			name:       "fresh savepoint init scaffold",
			build:      matrixBuildFreshInitScaffold,
			wantKind:   data.NextPlanObjective,
			wantAction: "Plan the next Objective",
		},
	}
}

func matrixBuildPendingMigration(t *testing.T, dir string) {
	t.Helper()
	writeMigrateMinimalProject(t, dir)
	if _, err := migrate.CreateOperation(dir, "op-matrix-pending", nil, nil, time.Now()); err != nil {
		t.Fatalf("CreateOperation() error = %v", err)
	}
}

func matrixConfig(t *testing.T, dir string) {
	t.Helper()
	testutil.WriteFile(t, filepath.Join(dir, ".savepoint", "config.yml"), "schema_version: 2\n")
}

func matrixObjective(t *testing.T, dir, dirName, id, status string) {
	t.Helper()
	content := "---\nid: " + id + "\ntitle: \"Objective " + id + "\"\nstatus: " + status + "\n---\n\n# Objective " + id + "\n"
	testutil.WriteFile(t, filepath.Join(dir, ".savepoint", "objectives", dirName, "Objective.md"), content)
}

func matrixBuildReplan(t *testing.T, dir string) {
	t.Helper()
	matrixConfig(t, dir)
	matrixObjective(t, dir, "O001-first", "O001", "planned")
	content := "---\nid: T001\ntitle: \"Alpha\"\nobjective: O001\n" +
		"planned_by: {role: planner, session: planning-fixture}\nstatus: planned\n" +
		"replan: {reason: \"scope changed\", recorded_by: {role: owner, session: owner-1}, recorded_at: '2026-09-14T00:00:00Z'}\n" +
		"---\n\n# Alpha\n"
	testutil.WriteFile(t, filepath.Join(dir, ".savepoint", "objectives", "O001-first", "tasks", "T001-alpha.md"), content)
	testutil.WriteFile(t, filepath.Join(dir, ".savepoint", "router.md"), resumeRouterV2Content("task", "O001", "T001", "Resolve the replan."))
}

func matrixBuildTaskDependency(t *testing.T, dir string) {
	t.Helper()
	matrixConfig(t, dir)
	matrixObjective(t, dir, "O001-first", "O001", "planned")
	t001 := "---\nid: T001\ntitle: \"Alpha\"\nobjective: O001\n" +
		"planned_by: {role: planner, session: planning-fixture}\nstatus: planned\n" +
		"depends_on: [{task: T002}]\n---\n\n# Alpha\n"
	testutil.WriteFile(t, filepath.Join(dir, ".savepoint", "objectives", "O001-first", "tasks", "T001-alpha.md"), t001)
	t002 := "---\nid: T002\ntitle: \"Beta\"\nobjective: O001\n" +
		"planned_by: {role: planner, session: planning-fixture}\nstatus: planned\n---\n\n# Beta\n"
	testutil.WriteFile(t, filepath.Join(dir, ".savepoint", "objectives", "O001-first", "tasks", "T002-beta.md"), t002)
	testutil.WriteFile(t, filepath.Join(dir, ".savepoint", "router.md"), resumeRouterV2Content("task", "O001", "T001", "Wait on T002."))
}

func matrixBuildExecute(t *testing.T, dir string) {
	t.Helper()
	writeResumeV2Project(t, dir)
}

func matrixBuildCheckNeeded(t *testing.T, dir string) {
	t.Helper()
	matrixConfig(t, dir)
	matrixObjective(t, dir, "O001-first", "O001", "planned")
	content := "---\nid: T001\ntitle: \"Alpha\"\nobjective: O001\n" +
		"planned_by: {role: planner, session: planning-fixture}\nstatus: in_progress\nstage: audit\n---\n\n# Alpha\n"
	testutil.WriteFile(t, filepath.Join(dir, ".savepoint", "objectives", "O001-first", "tasks", "T001-alpha.md"), content)
	testutil.WriteFile(t, filepath.Join(dir, ".savepoint", "router.md"), resumeRouterV2Content("check", "O001", "T001", "Record a fresh Check."))
}

func matrixBuildOwnerValidation(t *testing.T, dir string) {
	t.Helper()
	matrixConfig(t, dir)
	matrixObjective(t, dir, "O001-first", "O001", "planned")
	check := "---\nid: C001\nscope: {kind: task, id: T001}\nresult: CLEAR\n" +
		"checked_by: {role: checker, session: sess-1}\nexecuted_session: build-fixture\nchecked_at: '2026-09-14T00:00:00Z'\n---\n\n# Check\n"
	testutil.WriteFile(t, filepath.Join(dir, ".savepoint", "checks", "C001-alpha.md"), check)
	task := "---\nid: T001\ntitle: \"Alpha\"\nobjective: O001\n" +
		"planned_by: {role: planner, session: planning-fixture}\nstatus: in_progress\nstage: audit\n" +
		"freshness: {state: current, check: C001, assessed_by: {role: checker, session: sess-1}, assessed_at: '2026-09-14T00:00:00Z', basis: reviewed}\n" +
		"owner_validation: {required: true}\n---\n\n# Alpha\n"
	testutil.WriteFile(t, filepath.Join(dir, ".savepoint", "objectives", "O001-first", "tasks", "T001-alpha.md"), task)
	testutil.WriteFile(t, filepath.Join(dir, ".savepoint", "router.md"), resumeRouterV2Content("check", "O001", "T001", "Ask the owner to accept."))
}

func matrixBuildObjectiveIntegration(t *testing.T, dir string) {
	t.Helper()
	matrixConfig(t, dir)
	matrixObjective(t, dir, "O001-first", "O001", "in_progress")
	content := "---\nid: T001\ntitle: \"Alpha\"\nobjective: O001\n" +
		"planned_by: {role: planner, session: planning-fixture}\nstatus: done\n---\n\n# Alpha\n"
	testutil.WriteFile(t, filepath.Join(dir, ".savepoint", "objectives", "O001-first", "tasks", "T001-alpha.md"), content)
	testutil.WriteFile(t, filepath.Join(dir, ".savepoint", "router.md"), resumeRouterV2Content("check", "O001", "T001", "Record the Objective integration Check."))
}

// matrixBuildObjectiveNoTasks selects an Objective that owns no Task yet —
// the normal state between defining an Objective and breaking it into Tasks.
// Rung seven must not ask for an integration Check over work that does not
// exist; the Objective falls through to rung eight's project-wide search,
// which finds this same Objective ready to be planned under.
func matrixBuildObjectiveNoTasks(t *testing.T, dir string) {
	t.Helper()
	matrixConfig(t, dir)
	matrixObjective(t, dir, "O001-first", "O001", "planned")
	testutil.WriteFile(t, filepath.Join(dir, ".savepoint", "router.md"), resumeRouterV2Content("design", "O001", "none", "Plan Tasks under Objective O001."))
}

// matrixBuildObjectiveIntegrationUnselected mirrors matrixBuildObjectiveIntegration
// — every owned Task done, integration Check missing — but the router selects
// nothing, proving the outstanding Check is still reported project-wide
// rather than disappearing behind an empty selection.
func matrixBuildObjectiveIntegrationUnselected(t *testing.T, dir string) {
	t.Helper()
	matrixConfig(t, dir)
	matrixObjective(t, dir, "O001-first", "O001", "in_progress")
	content := "---\nid: T001\ntitle: \"Alpha\"\nobjective: O001\n" +
		"planned_by: {role: planner, session: planning-fixture}\nstatus: done\n---\n\n# Alpha\n"
	testutil.WriteFile(t, filepath.Join(dir, ".savepoint", "objectives", "O001-first", "tasks", "T001-alpha.md"), content)
	testutil.WriteFile(t, filepath.Join(dir, ".savepoint", "router.md"), resumeRouterV2Content("idea", "none", "none", "Plan the next thing."))
}

func matrixBuildReady(t *testing.T, dir string) {
	t.Helper()
	matrixConfig(t, dir)
	matrixObjective(t, dir, "O001-first", "O001", "planned")
	content := "---\nid: T001\ntitle: \"Alpha\"\nobjective: O001\n" +
		"planned_by: {role: planner, session: planning-fixture}\nstatus: planned\n---\n\n# Alpha\n"
	testutil.WriteFile(t, filepath.Join(dir, ".savepoint", "objectives", "O001-first", "tasks", "T001-alpha.md"), content)
	testutil.WriteFile(t, filepath.Join(dir, ".savepoint", "router.md"), resumeRouterV2Content("idea", "none", "none", "Plan the next thing."))
}

func matrixBuildPlanObjective(t *testing.T, dir string) {
	t.Helper()
	matrixConfig(t, dir)
	testutil.WriteFile(t, filepath.Join(dir, ".savepoint", "router.md"), resumeRouterV2Content("idea", "none", "none", "Plan the first Objective."))
}

// matrixBuildFreshInitScaffold writes the real embedded V2 template tree via
// savepointinit.Scaffold — the same call initRunner makes — rather than
// reconstructing the shipped router.md/config.yml content by hand, so this
// case proves the actual `savepoint init` output resumes cleanly, not a
// hand-written approximation of it (E48-Detail's "onboarding rather than an
// error" claim, and this task's fresh-scaffold acceptance criterion).
func matrixBuildFreshInitScaffold(t *testing.T, dir string) {
	t.Helper()
	sub, err := fs.Sub(projectTemplatesV2, "templates/project-v2")
	if err != nil {
		t.Fatalf("fs.Sub() error = %v", err)
	}
	if err := savepointinit.Scaffold(sub, dir, "matrix-fixture", false); err != nil {
		t.Fatalf("Scaffold() error = %v", err)
	}
}

// resolveNextFromDisk loads dir exactly as runResume's own read order does —
// pending migration first, then schema version, project, and router — and
// returns the resolved projection directly, so the matrix can assert
// next.Kind, next.Objective, and next.Task without parsing them back out of
// rendered prose.
func resolveNextFromDisk(t *testing.T, dir string) data.Next {
	t.Helper()

	pending, err := migrate.PendingOperation(dir)
	if err != nil {
		t.Fatalf("PendingOperation() error = %v", err)
	}
	if pending != nil {
		return data.ResolveNext(data.NextInput{Migration: data.MigrationState{Pending: true, OperationID: pending.OperationID}})
	}

	savepointRoot := filepath.Join(dir, ".savepoint")
	version, err := data.ReadSchemaVersion(filepath.Join(savepointRoot, "config.yml"))
	if err != nil {
		t.Fatalf("ReadSchemaVersion() error = %v", err)
	}
	if version != data.SchemaVersionV2 {
		t.Fatalf("SchemaVersion = %v, want V2 for a matrix fixture", version)
	}

	project, err := data.LoadProject(savepointRoot)
	if err != nil {
		t.Fatalf("LoadProject() error = %v", err)
	}

	routerContent, err := os.ReadFile(filepath.Join(savepointRoot, "router.md"))
	if err != nil {
		t.Fatalf("reading router.md: %v", err)
	}
	router, err := data.NewRouterReader().ReadStateV2(string(routerContent))
	if err != nil {
		t.Fatalf("ReadStateV2() error = %v", err)
	}

	return data.ResolveNext(data.NextInput{Index: project.V2, Router: router})
}

// minimalSecondConsumer is deliberately not internal/resume: it reads
// data.Next directly and derives its own selection identity and next-action
// category using its own vocabulary, proving the projection alone — with no
// resume-specific shape — carries what a second surface (E49's board) would
// need. It renders nothing textual; it reports typed facts.
type minimalSecondConsumerResult struct {
	SelectedID string // the Task or Objective ID in view, or "" for neither
	ActionKind data.NextKind
}

func minimalSecondConsumerRender(next data.Next) minimalSecondConsumerResult {
	result := minimalSecondConsumerResult{ActionKind: next.Kind}
	switch {
	case next.Task != nil:
		result.SelectedID = next.Task.ID
	case next.Objective != nil:
		result.SelectedID = next.Objective.ID
	}
	return result
}

// TestResumeMatrix_everyRungReachedExactlyOnce is E48 T007's central proof:
// each named project state above resolves through the real load-to-render
// pipeline to exactly the one rung it was built for, resume's rendered
// output names the matching next action, a second independent consumer
// derives the same selection and action kind from the same projection value,
// no matrix project's files are touched by any of it, and a second
// consecutive resume over the same project produces byte-identical output.
func TestResumeMatrix_everyRungReachedExactlyOnce(t *testing.T) {
	seenKinds := map[data.NextKind]bool{}

	for _, tc := range resumeMatrixCases() {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			tc.build(t, dir)
			before := snapshotDir(t, dir)

			next := resolveNextFromDisk(t, dir)
			if next.Kind != tc.wantKind {
				t.Fatalf("resolveNextFromDisk() Kind = %q, want exactly %q", next.Kind, tc.wantKind)
			}
			if next.SelectionDiagnostic != nil {
				t.Errorf("SelectionDiagnostic = %+v, want nil: no matrix case names a router selection that fails to resolve", next.SelectionDiagnostic)
			}
			seenKinds[next.Kind] = true

			var buf strings.Builder
			if err := resume.Render(&buf, next); err != nil {
				t.Fatalf("resume.Render() error = %v", err)
			}
			rendered := buf.String()
			if !strings.Contains(rendered, tc.wantAction) {
				t.Fatalf("rendered output = %q, want it to contain %q", rendered, tc.wantAction)
			}
			assertNarrowWidthReadable(t, rendered)

			consumed := minimalSecondConsumerRender(next)
			if consumed.ActionKind != next.Kind {
				t.Errorf("minimalSecondConsumerRender() ActionKind = %q, want %q", consumed.ActionKind, next.Kind)
			}
			if consumed.SelectedID != "" && !strings.Contains(rendered, consumed.SelectedID) {
				t.Errorf("second consumer selected %q, which resume's own rendering %q never names", consumed.SelectedID, rendered)
			}

			// A second consecutive resolve+render over the same untouched
			// project must reach the same rung and produce byte-identical
			// output: determinism holds across a repeat invocation, not just
			// within one.
			again := resolveNextFromDisk(t, dir)
			var again_buf strings.Builder
			if err := resume.Render(&again_buf, again); err != nil {
				t.Fatalf("second resume.Render() error = %v", err)
			}
			if again_buf.String() != rendered {
				t.Fatalf("second invocation rendered %q, want byte-identical to the first %q", again_buf.String(), rendered)
			}

			assertSameSnapshot(t, before, snapshotDir(t, dir))
		})
	}

	for _, kind := range []data.NextKind{
		data.NextPendingMigration, data.NextReplan, data.NextDependency, data.NextExecute,
		data.NextCheckNeeded, data.NextOwnerValidationRequired, data.NextObjectiveIntegration,
		data.NextReady, data.NextPlanObjective,
	} {
		if !seenKinds[kind] {
			t.Errorf("no matrix case reached rung %q; the ladder has a gap this matrix does not cover", kind)
		}
	}
}

// assertNarrowWidthReadable proves rendered text stays a sequence of short,
// unpadded sentences rather than a fixed-width table, matching the property
// internal/resume's own tests already prove on synthetic data.Next values —
// here re-checked against the text the real load-to-render pipeline
// produced for a matrix project.
func assertNarrowWidthReadable(t *testing.T, rendered string) {
	t.Helper()
	for _, glyph := range []string{"─", "│", "┌", "\t"} {
		if strings.Contains(rendered, glyph) {
			t.Errorf("rendered output contains fixed-width layout glyph %q, which does not stay readable at 40 columns", glyph)
		}
	}
}

// TestResumeMatrix_runResumeThroughTheRealCommandAlsoWritesNothingTwice
// exercises the same rungs through runResume itself (the exact function
// `savepoint resume` dispatches to), rather than the pipeline reassembled by
// resolveNextFromDisk, and proves two consecutive invocations through that
// entrypoint are also byte-identical and leave the project untouched.
func TestResumeMatrix_runResumeThroughTheRealCommandAlsoWritesNothingTwice(t *testing.T) {
	for _, tc := range resumeMatrixCases() {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			tc.build(t, dir)
			before := snapshotDir(t, dir)

			var first strings.Builder
			code, err := runResume(dir, &first)
			if err != nil {
				t.Fatalf("runResume() error = %v", err)
			}
			if code != 0 {
				t.Fatalf("runResume() code = %d, want 0 for an intact matrix project", code)
			}

			var second strings.Builder
			code, err = runResume(dir, &second)
			if err != nil {
				t.Fatalf("second runResume() error = %v", err)
			}
			if code != 0 {
				t.Fatalf("second runResume() code = %d, want 0", code)
			}
			if second.String() != first.String() {
				t.Fatalf("second runResume() output = %q, want byte-identical to the first %q", second.String(), first.String())
			}

			assertSameSnapshot(t, before, snapshotDir(t, dir))
		})
	}
}
