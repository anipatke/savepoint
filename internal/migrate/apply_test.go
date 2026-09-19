package migrate

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/data"
)

// copyFixtureProject copies fixture's project/ tree into a fresh temporary
// directory so apply tests can write into it without ever mutating the
// frozen fixtures under internal/data/testdata/migration/. It returns the
// copy's root (the directory that itself contains AGENTS.md and
// .savepoint/), matching fixtureProjectRoot's shape.
func copyFixtureProject(t *testing.T, fixture string) string {
	t.Helper()
	src := fixtureProjectRoot(fixture)
	dst := t.TempDir()

	err := filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, content, 0644)
	})
	if err != nil {
		t.Fatalf("copy fixture %s: %v", fixture, err)
	}
	return dst
}

func mustLoadProject(t *testing.T, root string) *data.Project {
	t.Helper()
	p, err := data.LoadProject(filepath.Join(root, ".savepoint"))
	if err != nil {
		t.Fatalf("LoadProject(%s) error = %v", root, err)
	}
	return p
}

func mustExist(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected %s to exist: %v", path, err)
	}
}

func mustNotExist(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Lstat(path); err == nil {
		t.Fatalf("expected %s to be gone, but it still exists", path)
	} else if !os.IsNotExist(err) {
		t.Fatalf("stat %s: %v", path, err)
	}
}

// assertNoPendingOperation checks via PendingOperation rather than a raw
// directory stat: CreateOperation's own os.MkdirAll leaves
// .savepoint/.migration/ itself in place (empty) once its one operation
// subdirectory is cleaned up, and PendingOperation already treats an empty
// operations root as "none pending" — the same check every real guard
// (upgrade-assets, board, doctor) would perform.
func assertNoPendingOperation(t *testing.T, root string) {
	t.Helper()
	report, err := PendingOperation(root)
	if err != nil {
		t.Fatalf("PendingOperation() error = %v", err)
	}
	if report != nil {
		t.Fatalf("PendingOperation() = %+v, want none", report)
	}
}

// --- fresh apply happy path --------------------------------------------

func TestApply_v1Basic_endToEnd(t *testing.T) {
	root := copyFixtureProject(t, "v1-basic")
	plan := mustPlan(t, root)
	if !plan.Appliable {
		t.Fatalf("plan.Appliable = false, want true; ambiguities = %+v", plan.Ambiguities)
	}

	result, err := Apply(root, plan)
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if result.AlreadyMigrated || result.Resumed {
		t.Errorf("result = %+v, want a fresh, non-resumed, non-already-migrated apply", result)
	}
	if result.OperationID == "" {
		t.Error("result.OperationID is empty")
	}

	// The operation is cleaned up once fully committed.
	assertNoPendingOperation(t, root)

	// Schema activated, and every other config.yml content preserved.
	version, err := data.ReadSchemaVersion(filepath.Join(root, ".savepoint", "config.yml"))
	if err != nil || version != data.SchemaVersionV2 {
		t.Fatalf("ReadSchemaVersion() = %v, %v, want SchemaVersionV2", version, err)
	}
	configContent := readFile(t, filepath.Join(root, ".savepoint", "config.yml"))
	if !strings.Contains(configContent, `bg: "#000000"`) {
		t.Errorf("config.yml lost its unrelated theme key:\n%s", configContent)
	}

	// The epic detail is always archived (and its original removed)
	// regardless of the epic's own active/done status; a done Task's
	// original is archived and removed the same way. An active Task's V1
	// source is converted but is not itself an ArchiveEntry, so it is left
	// in place — only settled work moves to archive/v1/.
	mustNotExist(t, filepath.Join(root, ".savepoint", "releases", "v1", "epics", "E01-example", "E01-Detail.md"))
	mustNotExist(t, filepath.Join(root, ".savepoint", "releases", "v1", "epics", "E01-example", "tasks", "T001-original.md"))
	mustExist(t, filepath.Join(root, ".savepoint", "releases", "v1", "epics", "E01-example", "tasks", "T002-follow-up.md"))
	mustNotExist(t, filepath.Join(root, ".savepoint", "PRD.md"))

	archived, ok := archiveByPath(plan, ".savepoint/releases/v1/epics/E01-example/tasks/T001-original.md")
	if !ok {
		t.Fatal("T001-original was not archived by the plan")
	}
	archiveContent := readFile(t, filepath.Join(root, filepath.FromSlash(archived.ArchivePath)))
	originalContent := readFile(t, filepath.Join(fixtureProjectRoot("v1-basic"), filepath.FromSlash(archived.SourcePath)))
	if archiveContent != originalContent {
		t.Errorf("archived T001-original does not match its original source bytes")
	}

	mustExist(t, filepath.Join(root, ".savepoint", "Idea.md"))
	mustExist(t, filepath.Join(root, ".savepoint", "migrations", "v1-to-v2.yml"))

	// Design.md is a preserved-in-place document: byte-identical, untouched.
	origDesign := readFile(t, filepath.Join(fixtureProjectRoot("v1-basic"), ".savepoint", "Design.md"))
	gotDesign := readFile(t, filepath.Join(root, ".savepoint", "Design.md"))
	if origDesign != gotDesign {
		t.Errorf("Design.md changed; preserved-in-place documents must stay byte-identical")
	}

	// The migrated project loads cleanly as V2, and no legacy directory
	// leaks into discovery.
	project := mustLoadProject(t, root)
	if project.SchemaVersion != data.SchemaVersionV2 {
		t.Fatalf("SchemaVersion = %v, want V2", project.SchemaVersion)
	}
	if len(project.V2.Objectives) != 1 {
		t.Errorf("Objectives = %d, want 1", len(project.V2.Objectives))
	}
	if len(project.V2.Tasks) != 1 {
		t.Errorf("Tasks = %d, want 1 (T001-original is archived, not converted)", len(project.V2.Tasks))
	}
	mustNotExist(t, filepath.Join(root, ".savepoint", "archive", "v1", "archive"))
}

func TestApply_v1History_distinctGlobalIDsAcrossReleases(t *testing.T) {
	root := copyFixtureProject(t, "v1-history")
	plan := mustPlan(t, root)
	if !plan.Appliable {
		t.Fatalf("plan.Appliable = false, want true; ambiguities = %+v", plan.Ambiguities)
	}

	if _, err := Apply(root, plan); err != nil {
		t.Fatalf("Apply() error = %v", err)
	}

	project := mustLoadProject(t, root)
	if project.SchemaVersion != data.SchemaVersionV2 {
		t.Fatalf("SchemaVersion = %v, want V2", project.SchemaVersion)
	}

	// v1's T001-shared (done, under a done epic) is archived; v1.1's
	// T001-shared (in_progress) and T002-follow-up (planned) convert to
	// exactly two active Tasks, each with its own distinct global ID.
	if len(project.V2.Tasks) != 2 {
		t.Errorf("Tasks = %d, want 2 (only v1.1's tasks are active)", len(project.V2.Tasks))
	}
	seen := map[string]bool{}
	for id := range project.V2.Tasks {
		if seen[id] {
			t.Errorf("global Task ID %s allocated more than once", id)
		}
		seen[id] = true
	}

	// v1's whole (done) epic is archived and its tasks' originals removed.
	// v1.1's active epic detail is archived and removed too (every epic
	// detail is, regardless of status), but its active tasks' V1 sources
	// are left in place — only settled work moves to .savepoint/archive/v1/.
	mustNotExist(t, filepath.Join(root, ".savepoint", "releases", "v1", "epics", "E01-example", "tasks", "T001-shared.md"))
	mustNotExist(t, filepath.Join(root, ".savepoint", "releases", "v1.1", "epics", "E01-example", "E01-Detail.md"))
	mustExist(t, filepath.Join(root, ".savepoint", "releases", "v1.1", "epics", "E01-example", "tasks", "T001-shared.md"))
	mustExist(t, filepath.Join(root, ".savepoint", "archive", "v1", ".savepoint", "releases", "v1", "epics", "E01-example", "tasks", "T001-shared.md"))
	mustExist(t, filepath.Join(root, ".savepoint", "archive", "v1", ".savepoint", "releases", "v1.1", "epics", "E01-example", "E01-Detail.md"))
}

// --- refusals write nothing ---------------------------------------------

func TestApply_notAppliable_writesNothing(t *testing.T) {
	root := copyFixtureProject(t, "v1-basic")

	// Corrupt an active record's status into something unrecognized, which
	// Plan reports as a blocking ambiguity rather than healing it.
	taskPath := filepath.Join(root, ".savepoint", "releases", "v1", "epics", "E01-example", "tasks", "T002-follow-up.md")
	content := readFile(t, taskPath)
	corrupted := replaceOnce(t, content, "status: in_progress", "status: sideways")
	if err := os.WriteFile(taskPath, []byte(corrupted), 0644); err != nil {
		t.Fatalf("corrupt fixture copy: %v", err)
	}

	plan := mustPlan(t, root)
	if plan.Appliable {
		t.Fatal("plan.Appliable = true, want false: an unrecognized active status must block")
	}

	before := snapshotTree(t, root)
	_, err := Apply(root, plan)
	if err == nil {
		t.Fatal("Apply() error = nil, want ErrPlanNotAppliable")
	}
	assertSnapshotsEqual(t, before, snapshotTree(t, root))
}

func TestApply_planConflict_writesNothing(t *testing.T) {
	root := copyFixtureProject(t, "v1-basic")
	if err := os.MkdirAll(filepath.Join(root, ".savepoint", "migrations"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".savepoint", "migrations", "v1-to-v2.yml"), []byte("stale: true\n"), 0644); err != nil {
		t.Fatal(err)
	}

	plan := mustPlan(t, root)
	if len(plan.Conflicts) == 0 {
		t.Fatal("plan.Conflicts is empty, want the existing-manifest conflict")
	}

	before := snapshotTree(t, root)
	_, err := Apply(root, plan)
	if err == nil {
		t.Fatal("Apply() error = nil, want ErrPlanConflict")
	}
	assertSnapshotsEqual(t, before, snapshotTree(t, root))
}

// TestApply_preExistingCreateDestination_refusesCleanly is regression
// coverage for an audit finding: a file already sitting exactly where a plan
// intends to create one used to reach Apply undetected and crash with an
// internal invariant error ("... has action create, which has no backup
// step"), leaving an unresumable operation directory behind with no way out.
// Plan now catches this as a named ConflictDestinationExists, so Apply
// refuses before creating any operation at all.
func TestApply_preExistingCreateDestination_refusesCleanly(t *testing.T) {
	root := copyFixtureProject(t, "v1-basic")
	// v1-basic's PRD.md relocates to Idea.md; a user who started early
	// occupies that destination before migrating.
	if err := os.WriteFile(filepath.Join(root, ".savepoint", "Idea.md"), []byte("my early idea\n"), 0644); err != nil {
		t.Fatal(err)
	}

	plan := mustPlan(t, root)
	if len(plan.Conflicts) == 0 {
		t.Fatal("plan.Conflicts is empty, want a ConflictDestinationExists naming .savepoint/Idea.md")
	}

	before := snapshotTree(t, root)
	_, err := Apply(root, plan)
	if !errors.Is(err, ErrPlanConflict) {
		t.Fatalf("Apply() error = %v, want ErrPlanConflict", err)
	}
	assertSnapshotsEqual(t, before, snapshotTree(t, root))
	assertNoPendingOperation(t, root)
}

// TestApply_lateCreateDestinationsAreCreateOnly covers the installation
// boundary rather than only Plan's earlier collision scan. A user can create
// any additive destination after Plan returns; each platform primitive must
// refuse that late destination and preserve the user's bytes.
func TestApply_lateCreateDestinationsAreCreateOnly(t *testing.T) {
	cases := []struct {
		name string
		path func(*ConversionPlan) string
	}{
		{
			name: "target",
			path: func(plan *ConversionPlan) string {
				for _, target := range plan.Targets {
					return savepointPath(target.InstallPath())
				}
				return ""
			},
		},
		{
			name: "document",
			path: func(plan *ConversionPlan) string {
				for _, document := range plan.Documents {
					if document.Kind == DocumentIdea {
						return savepointPath(document.TargetPath)
					}
				}
				return ""
			},
		},
		{
			name: "archive",
			path: func(plan *ConversionPlan) string {
				if len(plan.Archives) == 0 {
					return ""
				}
				return plan.Archives[0].ArchivePath
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := copyFixtureProject(t, "v1-basic")
			plan := mustPlan(t, root)
			relPath := tc.path(plan)
			if relPath == "" {
				t.Fatal("fixture did not produce the create destination class")
			}
			path := filepath.Join(root, filepath.FromSlash(relPath))
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				t.Fatalf("create late destination parent: %v", err)
			}
			const userBytes = "user content created after planning\n"
			if err := os.WriteFile(path, []byte(userBytes), 0644); err != nil {
				t.Fatalf("create late destination: %v", err)
			}

			_, err := Apply(root, plan)
			if !errors.Is(err, ErrCreateDestinationExists) {
				t.Fatalf("Apply() error = %v, want ErrCreateDestinationExists", err)
			}
			if got := readFile(t, path); got != userBytes {
				t.Fatalf("late destination %s changed to %q; want the user's bytes %q", relPath, got, userBytes)
			}
		})
	}
}

func TestApply_sourceChangedSincePreview_conflictNamesPath(t *testing.T) {
	root := copyFixtureProject(t, "v1-basic")
	plan := mustPlan(t, root)

	taskPath := filepath.Join(root, ".savepoint", "releases", "v1", "epics", "E01-example", "tasks", "T001-original.md")
	original := readFile(t, taskPath)
	if err := os.WriteFile(taskPath, []byte(original+"\nedited after preview\n"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Apply(root, plan)
	if err == nil {
		t.Fatal("Apply() error = nil, want a source conflict")
	}
	assertNoPendingOperation(t, root)
}

// --- interruption and resume ---------------------------------------------

func TestApply_interruptedBeforeActivation_stillLoadsAsV1(t *testing.T) {
	root := copyFixtureProject(t, "v1-basic")
	plan := mustPlan(t, root)

	batch, err := buildApplyBatch(root, plan)
	if err != nil {
		t.Fatalf("buildApplyBatch() error = %v", err)
	}
	op, err := CreateOperation(root, plan.OperationID, toManifestSources(plan.Sources), batch.journalEntries(), plan.GeneratedAt)
	if err != nil {
		t.Fatalf("CreateOperation() error = %v", err)
	}
	if err := publish(op, root, batch); err != nil {
		t.Fatalf("publish() error = %v", err)
	}
	// Deliberately stop here: every entry is verified, but activateSchema
	// has not run yet.

	project := mustLoadProject(t, root)
	if project.SchemaVersion != data.SchemaVersionV1 {
		t.Fatalf("SchemaVersion = %v, want V1 before activation", project.SchemaVersion)
	}

	report, err := PendingOperation(root)
	if err != nil {
		t.Fatalf("PendingOperation() error = %v", err)
	}
	if report == nil {
		t.Fatal("PendingOperation() = nil, want the still-open operation")
	}
}

func TestApply_interruptedAfterActivation_loadsAsV2(t *testing.T) {
	root := copyFixtureProject(t, "v1-basic")
	plan := mustPlan(t, root)

	batch, err := buildApplyBatch(root, plan)
	if err != nil {
		t.Fatalf("buildApplyBatch() error = %v", err)
	}
	op, err := CreateOperation(root, plan.OperationID, toManifestSources(plan.Sources), batch.journalEntries(), plan.GeneratedAt)
	if err != nil {
		t.Fatalf("CreateOperation() error = %v", err)
	}
	if err := publish(op, root, batch); err != nil {
		t.Fatalf("publish() error = %v", err)
	}
	if err := activateSchema(root); err != nil {
		t.Fatalf("activateSchema() error = %v", err)
	}
	// Deliberately skip cleanup, simulating a crash right after activation.

	project := mustLoadProject(t, root)
	if project.SchemaVersion != data.SchemaVersionV2 {
		t.Fatalf("SchemaVersion = %v, want V2 after activation", project.SchemaVersion)
	}
	if len(project.V2.Objectives) != 1 {
		t.Errorf("Objectives = %d, want 1 even though the operation directory was never cleaned up", len(project.V2.Objectives))
	}
}

func TestApply_resumesAfterInterruption_convergesToSameFinalState(t *testing.T) {
	root := copyFixtureProject(t, "v1-basic")
	plan := mustPlan(t, root)

	batch, err := buildApplyBatch(root, plan)
	if err != nil {
		t.Fatalf("buildApplyBatch() error = %v", err)
	}
	op, err := CreateOperation(root, plan.OperationID, toManifestSources(plan.Sources), batch.journalEntries(), plan.GeneratedAt)
	if err != nil {
		t.Fatalf("CreateOperation() error = %v", err)
	}

	// Simulate a crash partway through: back everything up and create the
	// additive records, but go no further.
	for _, w := range batch.backupTargets() {
		if err := ensureBackedUp(op, root, w.Path); err != nil {
			t.Fatalf("ensureBackedUp(%s) error = %v", w.Path, err)
		}
	}
	for _, w := range batch.creates {
		if err := publishWrite(op, root, w); err != nil {
			t.Fatalf("publishWrite(%s) error = %v", w.Path, err)
		}
	}

	pending, err := PendingOperation(root)
	if err != nil || pending == nil {
		t.Fatalf("PendingOperation() = %+v, %v, want the operation still open", pending, err)
	}

	// The same in-memory plan stands in for a caller reconstructing
	// byte-identical content across the resume, exactly as the package doc
	// on Apply requires.
	result, err := Apply(root, plan)
	if err != nil {
		t.Fatalf("resuming Apply() error = %v", err)
	}
	if !result.Resumed {
		t.Error("result.Resumed = false, want true")
	}

	assertNoPendingOperation(t, root)
	project := mustLoadProject(t, root)
	if project.SchemaVersion != data.SchemaVersionV2 {
		t.Fatalf("SchemaVersion = %v, want V2 after resume completes", project.SchemaVersion)
	}
	if len(project.V2.Objectives) != 1 || len(project.V2.Tasks) != 1 {
		t.Errorf("Objectives = %d, Tasks = %d, want 1 and 1", len(project.V2.Objectives), len(project.V2.Tasks))
	}
}

func TestApply_recoversAtEveryPublishBoundaryWithoutOverwritingUserEdits(t *testing.T) {
	baselineRoot := copyFixtureProject(t, "v1-history")
	baselinePlan := mustPlan(t, baselineRoot)
	if _, err := Apply(baselineRoot, baselinePlan); err != nil {
		t.Fatalf("baseline Apply() error = %v", err)
	}
	baseline := snapshotTree(t, baselineRoot)

	probeRoot := copyFixtureProject(t, "v1-history")
	probePlan := mustPlan(t, probeRoot)
	probeBatch, err := buildApplyBatch(probeRoot, probePlan)
	if err != nil {
		t.Fatalf("buildApplyBatch() error = %v", err)
	}
	boundaries := probeBatch.orderedWrites()
	if len(boundaries) == 0 {
		t.Fatal("migration produced no publish boundaries")
	}

	for _, boundary := range boundaries {
		boundary := boundary
		t.Run(string(boundary.Action)+" "+boundary.Path, func(t *testing.T) {
			root := copyFixtureProject(t, "v1-history")
			plan := mustPlan(t, root)
			failed := false
			oldWriteHook := afterPublishWriteHook
			oldRemovalHook := afterPublishRemovalHook
			t.Cleanup(func() {
				afterPublishWriteHook = oldWriteHook
				afterPublishRemovalHook = oldRemovalHook
			})
			if boundary.Action == ActionRemove {
				afterPublishRemovalHook = func(path string) error {
					if path == boundary.Path && !failed {
						failed = true
						return errors.New("synthetic publish interruption")
					}
					return nil
				}
			} else {
				afterPublishWriteHook = func(path string) error {
					if path == boundary.Path && !failed {
						failed = true
						return errors.New("synthetic publish interruption")
					}
					return nil
				}
			}

			if _, err := Apply(root, plan); err == nil {
				t.Fatalf("Apply() at %s completed; want synthetic interruption", boundary.Path)
			}
			if !failed {
				t.Fatalf("boundary hook for %s was never reached", boundary.Path)
			}

			// Clearing the fault injection proves the durable operation resumes
			// at every boundary and reaches the same bytes as a clean apply.
			afterPublishWriteHook = nil
			afterPublishRemovalHook = nil
			result, err := Apply(root, plan)
			if err != nil {
				t.Fatalf("resuming Apply() at %s: %v", boundary.Path, err)
			}
			if !result.Resumed {
				t.Fatalf("resume result at %s = %+v, want Resumed", boundary.Path, result)
			}
			assertSnapshotContentsEqual(t, baseline, snapshotTree(t, root))
		})
	}

	// A user edit after an interrupted install must be refused rather than
	// overwritten. This separate probe keeps the every-boundary recovery loop
	// above able to prove successful convergence as well.
	var editedBoundary stagedWrite
	for _, boundary := range boundaries {
		if boundary.Action != ActionRemove {
			editedBoundary = boundary
			break
		}
	}
	if editedBoundary.Path == "" {
		t.Fatal("migration produced no install boundary for the user-edit probe")
	}
	t.Run("user edit after interruption", func(t *testing.T) {
		root := copyFixtureProject(t, "v1-history")
		plan := mustPlan(t, root)
		oldWriteHook := afterPublishWriteHook
		oldRemovalHook := afterPublishRemovalHook
		t.Cleanup(func() {
			afterPublishWriteHook = oldWriteHook
			afterPublishRemovalHook = oldRemovalHook
		})
		afterPublishWriteHook = func(path string) error {
			if path == editedBoundary.Path {
				return errors.New("synthetic publish interruption")
			}
			return nil
		}
		if _, err := Apply(root, plan); err == nil {
			t.Fatalf("Apply() at %s completed; want synthetic interruption", editedBoundary.Path)
		}
		path := filepath.Join(root, filepath.FromSlash(editedBoundary.Path))
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read interrupted output %s: %v", editedBoundary.Path, err)
		}
		if err := os.WriteFile(path, append(content, []byte("\nuser edit after interruption\n")...), 0644); err != nil {
			t.Fatalf("edit interrupted output %s: %v", editedBoundary.Path, err)
		}
		afterPublishWriteHook = nil
		if _, err := Apply(root, plan); err == nil {
			t.Fatalf("resume after user edit at %s succeeded; want conflict", editedBoundary.Path)
		} else if !errors.Is(err, ErrVerificationFailed) {
			t.Fatalf("resume after user edit at %s error = %v, want ErrVerificationFailed", editedBoundary.Path, err)
		}
	})
}

func TestApply_resumeAfterUserEditedInstalledFile_reportsConflict(t *testing.T) {
	root := copyFixtureProject(t, "v1-basic")
	plan := mustPlan(t, root)

	batch, err := buildApplyBatch(root, plan)
	if err != nil {
		t.Fatalf("buildApplyBatch() error = %v", err)
	}
	op, err := CreateOperation(root, plan.OperationID, toManifestSources(plan.Sources), batch.journalEntries(), plan.GeneratedAt)
	if err != nil {
		t.Fatalf("CreateOperation() error = %v", err)
	}
	for _, w := range batch.backupTargets() {
		if err := ensureBackedUp(op, root, w.Path); err != nil {
			t.Fatalf("ensureBackedUp(%s) error = %v", w.Path, err)
		}
	}
	if len(batch.creates) == 0 {
		t.Fatal("test fixture produced no additive creates to tamper with")
	}
	// Drive the first create through staged and installed by hand, stopping
	// short of Verify — simulating a crash between Install and Verify,
	// before publishWrite's own call would have caught the tamper itself.
	first := batch.creates[0]
	if err := op.WriteStaged(first.Path, first.Content); err != nil {
		t.Fatalf("WriteStaged(%s) error = %v", first.Path, err)
	}
	if err := op.Install(root, first.Path, 0644); err != nil {
		t.Fatalf("Install(%s) error = %v", first.Path, err)
	}

	// Tamper with the just-installed file before it was ever verified.
	installed := filepath.Join(root, filepath.FromSlash(first.Path))
	if err := os.WriteFile(installed, []byte("tampered after install\n"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err = Apply(root, plan)
	if err == nil {
		t.Fatal("resuming Apply() error = nil, want a verification conflict over the tampered installed file")
	}
}

// --- idempotence -----------------------------------------------------

func TestApply_secondFullRun_isNoOp(t *testing.T) {
	root := copyFixtureProject(t, "v1-basic")
	plan := mustPlan(t, root)
	if _, err := Apply(root, plan); err != nil {
		t.Fatalf("first Apply() error = %v", err)
	}

	before := snapshotTree(t, root)

	plan2 := mustPlan(t, root)
	if !plan2.SchemaAlreadyV2 {
		t.Fatal("second Plan() over an already-migrated project did not report SchemaAlreadyV2")
	}
	result, err := Apply(root, plan2)
	if err != nil {
		t.Fatalf("second Apply() error = %v", err)
	}
	if !result.AlreadyMigrated {
		t.Error("result.AlreadyMigrated = false, want true")
	}

	assertSnapshotsEqual(t, before, snapshotTree(t, root))
}

// --- migrations/ coexistence -------------------------------------------

func TestApply_preservesExistingMigrationsDirectoryContent(t *testing.T) {
	root := copyFixtureProject(t, "v1-basic")
	if err := os.MkdirAll(filepath.Join(root, ".savepoint", "migrations"), 0755); err != nil {
		t.Fatal(err)
	}
	readmePath := filepath.Join(root, ".savepoint", "migrations", "README.md")
	skillPath := filepath.Join(root, ".savepoint", "migrations", "savepoint-audit-SKILL.md")
	if err := os.WriteFile(readmePath, []byte("# preexisting readme\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(skillPath, []byte("# archived legacy skill\n"), 0644); err != nil {
		t.Fatal(err)
	}

	plan := mustPlan(t, root)
	if _, err := Apply(root, plan); err != nil {
		t.Fatalf("Apply() error = %v", err)
	}

	if got := readFile(t, readmePath); got != "# preexisting readme\n" {
		t.Errorf("migrations/README.md changed: %q", got)
	}
	if got := readFile(t, skillPath); got != "# archived legacy skill\n" {
		t.Errorf("migrations/savepoint-audit-SKILL.md changed: %q", got)
	}
	mustExist(t, filepath.Join(root, ".savepoint", "migrations", "v1-to-v2.yml"))
}

// --- helpers ------------------------------------------------------------

func replaceOnce(t *testing.T, content, old, new string) string {
	t.Helper()
	if !strings.Contains(content, old) {
		t.Fatalf("content does not contain %q", old)
	}
	return strings.Replace(content, old, new, 1)
}

func assertSnapshotContentsEqual(t *testing.T, before, after treeSnapshot) {
	t.Helper()
	if len(before) != len(after) {
		t.Fatalf("file count changed: before %d, after %d", len(before), len(after))
	}
	for path, want := range before {
		got, ok := after[path]
		if !ok {
			t.Errorf("%s was present before and missing after", path)
			continue
		}
		if string(want.content) != string(got.content) {
			t.Errorf("%s content changed", path)
		}
	}
}
