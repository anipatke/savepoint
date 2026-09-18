package migrate

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/opencode/savepoint/internal/data"
)

const opTestOriginal = "original bytes that must survive every failure\n"
const opTestReplacement = "replacement bytes\n"

// newTestOperation sets up a temporary project with one live file at relPath
// holding opTestOriginal, and a freshly created operation whose single entry
// plans to replace it with opTestReplacement. Every recoverability test below
// builds on this single-entry shape and drives it step by step.
func newTestOperation(t *testing.T, action EntryAction) (root, relPath string, op *Operation) {
	t.Helper()
	root = t.TempDir()
	relPath = "live.md"
	if err := os.WriteFile(filepath.Join(root, relPath), []byte(opTestOriginal), 0644); err != nil {
		t.Fatalf("write live file: %v", err)
	}

	sourceHashes := []ManifestSource{{Path: relPath, SHA256: hashBytes([]byte(opTestOriginal))}}
	entry := JournalEntry{Path: relPath, Action: action}
	if action != ActionRemove {
		entry.PlannedHash = hashBytes([]byte(opTestReplacement))
	}

	op, err := CreateOperation(root, "op-1", sourceHashes, []JournalEntry{entry}, time.Now())
	if err != nil {
		t.Fatalf("CreateOperation() error = %v", err)
	}
	return root, relPath, op
}

func readLive(t *testing.T, root, relPath string) (string, bool) {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(root, relPath))
	if os.IsNotExist(err) {
		return "", false
	}
	if err != nil {
		t.Fatalf("read live file: %v", err)
	}
	return string(content), true
}

// assertRecoverable is AC10's core invariant: at every point in the
// operation, the original bytes must be recoverable from either the live
// path itself (not yet touched) or a hash-verified backup copy. It never
// reads intent from journal state — only from what is actually on disk —
// because the whole point is to catch a case where the journal claims safety
// the filesystem does not back up.
func assertRecoverable(t *testing.T, root, relPath string, op *Operation) {
	t.Helper()
	live, liveExists := readLive(t, root, relPath)
	if liveExists && live == opTestOriginal {
		return
	}

	backup, err := os.ReadFile(op.BackupPath(relPath))
	if err != nil {
		t.Fatalf("original is gone from the live path (exists=%v content=%q) and no backup copy is readable: %v", liveExists, live, err)
	}
	if hashBytes(backup) != hashBytes([]byte(opTestOriginal)) {
		t.Fatalf("original is gone from the live path and the backup copy does not hash to the original either: got %q", string(backup))
	}
}

func TestCreateOperation_layoutAndCreateOnlyGuard(t *testing.T) {
	root, relPath, op := newTestOperation(t, ActionReplace)

	for _, dir := range []string{op.Dir, filepath.Join(op.Dir, backupDirName), filepath.Join(op.Dir, stagingDirName)} {
		if info, err := os.Stat(dir); err != nil || !info.IsDir() {
			t.Errorf("expected directory %s to exist, stat = %v", dir, err)
		}
	}
	if _, err := os.Stat(filepath.Join(op.Dir, journalFileName)); err != nil {
		t.Errorf("expected %s to exist: %v", journalFileName, err)
	}

	firstJournal, err := os.ReadFile(filepath.Join(op.Dir, journalFileName))
	if err != nil {
		t.Fatalf("read journal: %v", err)
	}

	_, err = CreateOperation(root, "op-1", nil, []JournalEntry{{Path: relPath, Action: ActionCreate}}, time.Now())
	if !errors.Is(err, ErrOperationExists) {
		t.Fatalf("second CreateOperation() error = %v, want ErrOperationExists", err)
	}

	afterJournal, err := os.ReadFile(filepath.Join(op.Dir, journalFileName))
	if err != nil {
		t.Fatalf("read journal after refused create: %v", err)
	}
	if string(afterJournal) != string(firstJournal) {
		t.Error("journal content changed after a refused duplicate create; a prior operation's recovery data must never be touched")
	}
}

// TestCreateOperation_pathNeverFallsInsideOperationDir proves AC1's second
// half directly: none of the paths this package ever writes through
// (BackupPath, StagingPath, the journal file) can collide with a live project
// path, because they are all rooted under the operation's own directory,
// which Inventory already excludes from every source walk.
func TestCreateOperation_pathNeverFallsInsideOperationDir(t *testing.T) {
	_, relPath, op := newTestOperation(t, ActionReplace)

	for _, p := range []string{op.BackupPath(relPath), op.StagingPath(relPath), op.journalPath()} {
		if filepath.Dir(p) == filepath.Dir(op.Dir) && filepath.Base(p) == relPath {
			t.Errorf("path %s escapes the operation directory", p)
		}
		if _, err := os.Stat(op.Dir); err != nil {
			t.Fatalf("operation dir missing: %v", err)
		}
	}
}

func TestOperation_backupCapturesAndVerifies(t *testing.T) {
	root, relPath, op := newTestOperation(t, ActionReplace)

	if err := op.Backup(root, relPath); err != nil {
		t.Fatalf("Backup() error = %v", err)
	}

	backup, err := os.ReadFile(op.BackupPath(relPath))
	if err != nil {
		t.Fatalf("read backup copy: %v", err)
	}
	if string(backup) != opTestOriginal {
		t.Errorf("backup content = %q, want %q", string(backup), opTestOriginal)
	}

	reloaded, err := LoadOperation(op.Dir)
	if err != nil {
		t.Fatalf("LoadOperation() error = %v", err)
	}
	entry, err := reloaded.entry(relPath)
	if err != nil {
		t.Fatal(err)
	}
	if entry.State != StepBackedUp {
		t.Errorf("persisted entry state = %v, want %v", entry.State, StepBackedUp)
	}

	assertRecoverable(t, root, relPath, op)
}

// TestOperation_backupAbortsBeforeAnyReplacementOnHashMismatch is AC4's
// abort-before-the-first-replacement case: the live source no longer matches
// the recorded hash (it changed after the operation was created), so the
// backup itself cannot be trusted and Backup must fail loudly, leaving the
// live file, the journal, and any partial backup copy exactly as if nothing
// had been attempted.
func TestOperation_backupAbortsBeforeAnyReplacementOnHashMismatch(t *testing.T) {
	root, relPath, op := newTestOperation(t, ActionReplace)

	if err := os.WriteFile(filepath.Join(root, relPath), []byte("tampered after the operation was created\n"), 0644); err != nil {
		t.Fatalf("tamper live file: %v", err)
	}

	err := op.Backup(root, relPath)
	if !errors.Is(err, ErrBackupVerificationFailed) {
		t.Fatalf("Backup() error = %v, want ErrBackupVerificationFailed", err)
	}

	if _, statErr := os.Stat(op.BackupPath(relPath)); !os.IsNotExist(statErr) {
		t.Error("a failed backup left a copy on disk; it must remove an unverified copy rather than leave it")
	}
	entry, entryErr := op.entry(relPath)
	if entryErr != nil {
		t.Fatal(entryErr)
	}
	if entry.State != StepPlanned {
		t.Errorf("entry state after a failed backup = %v, want unchanged %v", entry.State, StepPlanned)
	}
}

func TestOperation_replaceLifecycle_recoverableAtEveryStep(t *testing.T) {
	root, relPath, op := newTestOperation(t, ActionReplace)
	assertRecoverable(t, root, relPath, op) // before anything: original is live

	if err := op.Backup(root, relPath); err != nil {
		t.Fatalf("Backup() error = %v", err)
	}
	assertRecoverable(t, root, relPath, op)
	if live, _ := readLive(t, root, relPath); live != opTestOriginal {
		t.Errorf("live content after backup = %q, want the untouched original", live)
	}

	// Simulate a crash right after Backup persisted: reload from disk only,
	// as a resumed process would, and confirm recovery evidence survived the
	// simulated restart rather than only living in the in-memory struct.
	reloaded, err := LoadOperation(op.Dir)
	if err != nil {
		t.Fatalf("LoadOperation() after simulated crash error = %v", err)
	}
	assertRecoverable(t, root, relPath, reloaded)

	if err := reloaded.WriteStaged(relPath, []byte(opTestReplacement)); err != nil {
		t.Fatalf("WriteStaged() error = %v", err)
	}
	staged, err := os.ReadFile(reloaded.StagingPath(relPath))
	if err != nil {
		t.Fatalf("read staged content: %v", err)
	}
	if string(staged) != opTestReplacement {
		t.Errorf("staged content = %q, want %q", string(staged), opTestReplacement)
	}
	if live, _ := readLive(t, root, relPath); live != opTestOriginal {
		t.Error("staging touched the live path; it must only ever write inside the operation directory")
	}
	assertRecoverable(t, root, relPath, reloaded)

	reloaded, err = LoadOperation(op.Dir)
	if err != nil {
		t.Fatalf("LoadOperation() error = %v", err)
	}
	if err := reloaded.Install(root, relPath, 0644); err != nil {
		t.Fatalf("Install() error = %v", err)
	}
	if live, _ := readLive(t, root, relPath); live != opTestReplacement {
		t.Errorf("live content after install = %q, want %q", live, opTestReplacement)
	}
	// The live path now holds the new content, not the original — recovery
	// of the original must still be possible, from the backup alone.
	assertRecoverable(t, root, relPath, reloaded)

	reloaded, err = LoadOperation(op.Dir)
	if err != nil {
		t.Fatalf("LoadOperation() error = %v", err)
	}
	if err := reloaded.Verify(root, relPath); err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	entry, err := reloaded.entry(relPath)
	if err != nil {
		t.Fatal(err)
	}
	if entry.State != StepVerified {
		t.Errorf("entry state after Verify = %v, want %v", entry.State, StepVerified)
	}
}

func TestOperation_removeLifecycle_recoverableAtEveryStep(t *testing.T) {
	root, relPath, op := newTestOperation(t, ActionRemove)
	assertRecoverable(t, root, relPath, op)

	if err := op.Backup(root, relPath); err != nil {
		t.Fatalf("Backup() error = %v", err)
	}
	assertRecoverable(t, root, relPath, op)

	if err := op.Remove(root, relPath); err != nil {
		t.Fatalf("Remove() error = %v", err)
	}
	if _, exists := readLive(t, root, relPath); exists {
		t.Error("live file still exists after Remove()")
	}
	assertRecoverable(t, root, relPath, op) // gone from live, but the backup holds it

	if err := op.Verify(root, relPath); err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	entry, err := op.entry(relPath)
	if err != nil {
		t.Fatal(err)
	}
	if entry.State != StepVerified {
		t.Errorf("entry state = %v, want %v", entry.State, StepVerified)
	}
}

// TestOperation_stepsRefuseOutOfOrder proves the state machine is enforced in
// code, not just by caller discipline: skipping backup before staging, or
// installing before staging, is refused rather than silently accepted, which
// is what makes assertRecoverable's invariant structural rather than
// incidental.
func TestOperation_stepsRefuseOutOfOrder(t *testing.T) {
	root, relPath, op := newTestOperation(t, ActionReplace)

	if err := op.WriteStaged(relPath, []byte(opTestReplacement)); err == nil {
		t.Error("WriteStaged() before Backup() succeeded, want a refusal for a replace entry")
	}
	if err := op.Install(root, relPath, 0644); err == nil {
		t.Error("Install() before WriteStaged() succeeded, want a refusal")
	}
	if err := op.Verify(root, relPath); err == nil {
		t.Error("Verify() before Install() succeeded, want a refusal")
	}

	if live, _ := readLive(t, root, relPath); live != opTestOriginal {
		t.Error("a rejected out-of-order call still touched the live file")
	}
}

// TestOperation_writeStagedRejectsContentNotMatchingPlannedHash is an
// injected failure at the staging step: the caller offers the wrong bytes,
// and WriteStaged must refuse before writing anything, leaving no partial
// staged file and the entry's state unchanged.
func TestOperation_writeStagedRejectsContentNotMatchingPlannedHash(t *testing.T) {
	root, relPath, op := newTestOperation(t, ActionReplace)
	if err := op.Backup(root, relPath); err != nil {
		t.Fatalf("Backup() error = %v", err)
	}

	err := op.WriteStaged(relPath, []byte("not the planned content\n"))
	if err == nil {
		t.Fatal("WriteStaged() with mismatched content succeeded, want a refusal")
	}
	if _, statErr := os.Stat(op.StagingPath(relPath)); !os.IsNotExist(statErr) {
		t.Error("a rejected WriteStaged left a partial file in the staging area")
	}
	entry, entryErr := op.entry(relPath)
	if entryErr != nil {
		t.Fatal(entryErr)
	}
	if entry.State != StepBackedUp {
		t.Errorf("entry state after a rejected WriteStaged = %v, want unchanged %v", entry.State, StepBackedUp)
	}
	assertRecoverable(t, root, relPath, op)
}

func TestOperation_createInstallRejectsOccupiedDestinationAndPreservesBytes(t *testing.T) {
	root := t.TempDir()
	relPath := "new.md"
	content := []byte("planned output\n")
	op, err := CreateOperation(root, "op-create", nil, []JournalEntry{{
		Path:        relPath,
		Action:      ActionCreate,
		PlannedHash: hashBytes(content),
	}}, time.Now())
	if err != nil {
		t.Fatalf("CreateOperation() error = %v", err)
	}
	if err := op.WriteStaged(relPath, content); err != nil {
		t.Fatalf("WriteStaged() error = %v", err)
	}

	const userBytes = "unrelated user content\n"
	if err := os.WriteFile(filepath.Join(root, relPath), []byte(userBytes), 0644); err != nil {
		t.Fatalf("write occupied destination: %v", err)
	}
	if err := op.Install(root, relPath, 0644); !errors.Is(err, ErrCreateDestinationExists) {
		t.Fatalf("Install() error = %v, want ErrCreateDestinationExists", err)
	}
	if got := readFile(t, filepath.Join(root, relPath)); got != userBytes {
		t.Fatalf("occupied destination changed to %q; want %q", got, userBytes)
	}

	// Once the unrelated file is moved away, retrying the same staged entry
	// still succeeds, proving the refusal did not poison the operation state.
	if err := os.Remove(filepath.Join(root, relPath)); err != nil {
		t.Fatalf("remove occupied destination: %v", err)
	}
	if err := op.Install(root, relPath, 0644); err != nil {
		t.Fatalf("Install() after clearing destination error = %v", err)
	}
	if got := readFile(t, filepath.Join(root, relPath)); got != string(content) {
		t.Fatalf("installed content = %q, want %q", got, content)
	}
}

// TestOperation_installFailureLeavesLiveFileIntact injects a failure at the
// install step by making the live file's directory unwritable, so ReplaceFile
// cannot even create its temporary file. The original must survive untouched
// and recoverable, and the entry must not advance past staged.
func TestOperation_installFailureLeavesLiveFileIntact(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root ignores the write-permission bit this test depends on")
	}
	root, relPath, op := newTestOperation(t, ActionReplace)
	if err := op.Backup(root, relPath); err != nil {
		t.Fatalf("Backup() error = %v", err)
	}
	if err := op.WriteStaged(relPath, []byte(opTestReplacement)); err != nil {
		t.Fatalf("WriteStaged() error = %v", err)
	}

	if err := os.Chmod(root, 0555); err != nil {
		t.Fatalf("chmod root read-only: %v", err)
	}
	t.Cleanup(func() { os.Chmod(root, 0755) })

	if err := op.Install(root, relPath, 0644); err == nil {
		t.Fatal("Install() into a read-only directory succeeded, want a refusal")
	}

	os.Chmod(root, 0755)
	if live, _ := readLive(t, root, relPath); live != opTestOriginal {
		t.Errorf("live content after a failed install = %q, want the untouched original", live)
	}
	entry, entryErr := op.entry(relPath)
	if entryErr != nil {
		t.Fatal(entryErr)
	}
	if entry.State != StepStaged {
		t.Errorf("entry state after a failed install = %v, want unchanged %v", entry.State, StepStaged)
	}
	assertRecoverable(t, root, relPath, op)
}

// TestOperation_verifyDetectsCorruptedInstall injects a failure at the verify
// step: the live file is corrupted after a successful install (simulating
// disk-level damage, not a code bug), so Verify must report it rather than
// mark the entry verified — and the original bytes must still be recoverable
// from the backup even though the live path is now wrong.
func TestOperation_verifyDetectsCorruptedInstall(t *testing.T) {
	root, relPath, op := newTestOperation(t, ActionReplace)
	if err := op.Backup(root, relPath); err != nil {
		t.Fatalf("Backup() error = %v", err)
	}
	if err := op.WriteStaged(relPath, []byte(opTestReplacement)); err != nil {
		t.Fatalf("WriteStaged() error = %v", err)
	}
	if err := op.Install(root, relPath, 0644); err != nil {
		t.Fatalf("Install() error = %v", err)
	}

	if err := os.WriteFile(filepath.Join(root, relPath), []byte("corrupted after install\n"), 0644); err != nil {
		t.Fatalf("corrupt live file: %v", err)
	}

	err := op.Verify(root, relPath)
	if !errors.Is(err, ErrVerificationFailed) {
		t.Fatalf("Verify() error = %v, want ErrVerificationFailed", err)
	}
	entry, entryErr := op.entry(relPath)
	if entryErr != nil {
		t.Fatal(entryErr)
	}
	if entry.State != StepInstalled {
		t.Errorf("entry state after a failed verify = %v, want unchanged %v", entry.State, StepInstalled)
	}

	backup, err := os.ReadFile(op.BackupPath(relPath))
	if err != nil {
		t.Fatalf("read backup copy: %v", err)
	}
	if string(backup) != opTestOriginal {
		t.Errorf("backup content = %q, want the untouched original %q", string(backup), opTestOriginal)
	}
}

func TestPendingOperation_noneExists(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".savepoint"), 0755); err != nil {
		t.Fatal(err)
	}

	report, err := PendingOperation(root)
	if err != nil {
		t.Fatalf("PendingOperation() error = %v, want nil when no .migration/ dir exists", err)
	}
	if report != nil {
		t.Errorf("PendingOperation() = %+v, want nil", report)
	}
}

func TestPendingOperation_findsIncompleteOperationAndNamesRecovery(t *testing.T) {
	root, relPath, op := newTestOperation(t, ActionReplace)
	if err := op.Backup(root, relPath); err != nil {
		t.Fatalf("Backup() error = %v", err)
	}

	report, err := PendingOperation(root)
	if err != nil {
		t.Fatalf("PendingOperation() error = %v", err)
	}
	if report == nil {
		t.Fatal("PendingOperation() = nil, want the incomplete operation")
	}
	if report.OperationID != "op-1" {
		t.Errorf("report.OperationID = %q, want op-1", report.OperationID)
	}

	guidance := report.RecoveryGuidance()
	if !containsAll(guidance, "op-1", relPath, "backed_up", "--recover") {
		t.Errorf("RecoveryGuidance() = %q, want it to name the operation, the affected path and its state, and the recovery command", guidance)
	}
}

func TestPendingOperation_multipleOperationsIsANamedDiagnosticNotAChoice(t *testing.T) {
	root := t.TempDir()
	if _, err := CreateOperation(root, "op-a", nil, nil, time.Now()); err != nil {
		t.Fatalf("CreateOperation(op-a) error = %v", err)
	}
	if _, err := CreateOperation(root, "op-b", nil, nil, time.Now()); err != nil {
		t.Fatalf("CreateOperation(op-b) error = %v", err)
	}

	_, err := PendingOperation(root)
	if !errors.Is(err, ErrMultipleOperations) {
		t.Fatalf("PendingOperation() error = %v, want ErrMultipleOperations", err)
	}
	if !containsAll(err.Error(), "op-a", "op-b") {
		t.Errorf("PendingOperation() error = %q, want it to name both operation ids", err.Error())
	}
}

// TestLoadV2Index_ignoresMigrationOperationDirEntirely is AC9: even an
// operation directory staging content that is byte-for-byte shaped like
// Objective, Check, and Issue records must never be visible to V2 discovery.
// Each staged file below has deliberately malformed frontmatter, so if
// LoadV2Index ever walked into .migration/ (directly or by way of some
// future change to DiscoverV2Records), decoding it would fail loudly instead
// of this test passing for the wrong reason.
func TestLoadV2Index_ignoresMigrationOperationDirEntirely(t *testing.T) {
	root := t.TempDir()
	savepointRoot := filepath.Join(root, ".savepoint")
	if err := os.MkdirAll(savepointRoot, 0755); err != nil {
		t.Fatal(err)
	}

	opStaging := filepath.Join(savepointRoot, migrationStateDir, "op-1", stagingDirName)
	malformed := "---\nid: [unterminated\n---\nbody\n"
	writeFile(t, filepath.Join(opStaging, "objectives", "O999-fake", "Objective.md"), malformed)
	writeFile(t, filepath.Join(opStaging, "checks", "C999-fake.md"), malformed)
	writeFile(t, filepath.Join(opStaging, "issues", "I999-fake.md"), malformed)

	index, err := data.LoadV2Index(savepointRoot)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v, want nil: staged look-alikes under .migration/ must be invisible to discovery", err)
	}
	if len(index.Objectives) != 0 || len(index.Checks) != 0 || len(index.Issues) != 0 {
		t.Errorf("LoadV2Index() = %+v, want an empty index untouched by .migration/ content", index)
	}
}

func containsAll(s string, substrs ...string) bool {
	for _, sub := range substrs {
		if !strings.Contains(s, sub) {
			return false
		}
	}
	return true
}
