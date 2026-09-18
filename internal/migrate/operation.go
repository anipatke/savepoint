// Package migrate's operation.go owns the on-disk recovery record for one
// migration apply. A migration can be interrupted at any point, and no
// filesystem this project targets can change many files atomically, so this
// file gives the operation somewhere durable to record what it has actually
// done: the operation directory, the per-path journal, backup capture and
// verification, staging writes, and the read-only detector other packages
// (upgrade-assets, the board's write commands, doctor) consult before they
// write anything of their own.
//
// The operation lives at .savepoint/.migration/<op-id>/, deliberately inside
// the project (so it survives with it and stays writable, including when the
// project directory is copied) and outside every path the operation itself
// writes to, so a partial apply can never damage its own recovery data. It
// reuses migrationStateDir (inventory.go) rather than a second constant for
// the identical path — Inventory already excludes it from every source walk.
//
// This file owns the generic, plan-agnostic mechanics only: how one path
// moves through planned -> backed_up -> staged -> installed -> verified.
// Deciding which paths to touch, in what order, and when to activate
// schema_version: 2 is the ordered publish in apply.go (a later task); this
// file never reads a ConversionPlan.
package migrate

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Distinct operation diagnostics, so a caller (and a test) can tell them
// apart with errors.Is rather than string-matching a message.
var (
	// ErrOperationExists means the create-only operation directory guard
	// refused to reuse or overwrite an existing operation: a prior backup or
	// user sidecar living there can never be destroyed this way.
	ErrOperationExists = errors.New("migrate: operation directory already exists")
	// ErrBackupVerificationFailed means a backup copy did not hash to the
	// recorded source hash once written, so the copy itself is untrustworthy
	// (a disk fault during the copy, not merely a stale source — freshness
	// revalidation against a changed source is a separate, earlier check).
	ErrBackupVerificationFailed = errors.New("migrate: backup copy failed hash verification")
	// ErrVerificationFailed means an installed path's live content (or, for a
	// removal, its continued existence) does not match what the journal
	// recorded as installed.
	ErrVerificationFailed = errors.New("migrate: installed content failed verification")
	// ErrMultipleOperations means more than one incomplete migration
	// operation directory exists. PendingOperation never heuristically picks
	// one — the same "never guess" rule the rest of this package applies to
	// ambiguous project content.
	ErrMultipleOperations = errors.New("migrate: more than one incomplete migration operation exists")
)

const (
	backupDirName   = "backup"
	stagingDirName  = "staging"
	journalFileName = "operation.yml"
)

// StepState is where one journalled path stands in its own lifecycle. The
// vocabulary is fixed and shared by every action kind; an action that has no
// use for a given state (a create has no backup, a remove has no staging)
// simply skips it rather than growing a parallel vocabulary.
type StepState string

const (
	StepPlanned   StepState = "planned"
	StepBackedUp  StepState = "backed_up"
	StepStaged    StepState = "staged"
	StepInstalled StepState = "installed"
	StepVerified  StepState = "verified"
)

// EntryAction names what an apply will do to one path.
type EntryAction string

const (
	// ActionCreate is an additive write to a path that must not already
	// exist: a new V2 record, an archive entry, the manifest. It has no
	// backup step — there is nothing at the destination to lose.
	ActionCreate EntryAction = "create"
	// ActionReplace overwrites an existing live file in place, through
	// ReplaceFile, and is backed up first.
	ActionReplace EntryAction = "replace"
	// ActionRemove deletes an existing live file (its content having been
	// preserved elsewhere, such as archive/v1/), and is backed up first.
	ActionRemove EntryAction = "remove"
)

// JournalEntry is one path's recorded progress through an operation: what
// will happen to it, the hash its final content is planned to have (empty for
// ActionRemove, which installs no content), the hash actually installed once
// written, and which state it has reached.
type JournalEntry struct {
	Path          string      `yaml:"path"`
	Action        EntryAction `yaml:"action"`
	PlannedHash   string      `yaml:"planned_hash,omitempty"`
	InstalledHash string      `yaml:"installed_hash,omitempty"`
	State         StepState   `yaml:"state"`
}

// Journal is operation.yml's content: the operation's identity, the recorded
// source hashes a backup is verified against (the same values a plan's
// inventory read — see ManifestSource), and one entry per path the operation
// will touch.
type Journal struct {
	OperationID  string           `yaml:"operation_id"`
	CreatedAt    time.Time        `yaml:"created_at"`
	SourceHashes []ManifestSource `yaml:"source_hashes"`
	Entries      []JournalEntry   `yaml:"entries"`
}

// Operation is one migration operation directory bound to its journal.
// Methods on Operation are the only way this package moves a path through its
// lifecycle; every one of them persists the journal durably before returning
// on success, so the journal on disk is always exactly the state of the last
// step that actually finished. That is also what "written before and after
// each step" means in practice: the previous step's persisted state is the
// next step's on-disk "before", and each method's own persist is its "after".
// A step that fails leaves the journal exactly as the last successful step
// left it — never a state that claims work which did not finish.
type Operation struct {
	Dir     string
	Journal Journal
}

func operationsRootDir(projectRoot string) string {
	return filepath.Join(projectRoot, ".savepoint", migrationStateDir)
}

// OperationDir returns the operation directory for opID under projectRoot.
func OperationDir(projectRoot, opID string) string {
	return filepath.Join(operationsRootDir(projectRoot), opID)
}

// CreateOperation creates a brand-new operation directory at
// .savepoint/.migration/<opID>/ with backup/ and staging/ subdirectories and
// an initial operation.yml recording every entry as planned. It is
// create-only: an operation directory that already exists is refused with
// ErrOperationExists rather than reused or overwritten. Any failure after the
// directory itself is created removes it again, so a failed create never
// leaves a half-formed operation directory that PendingOperation would later
// mistake for one to resume.
func CreateOperation(projectRoot, opID string, sourceHashes []ManifestSource, entries []JournalEntry, createdAt time.Time) (op *Operation, err error) {
	dir := OperationDir(projectRoot, opID)

	if mkErr := os.MkdirAll(operationsRootDir(projectRoot), 0755); mkErr != nil {
		return nil, fmt.Errorf("create %s: %w", operationsRootDir(projectRoot), mkErr)
	}
	if mkErr := os.Mkdir(dir, 0755); mkErr != nil {
		if os.IsExist(mkErr) {
			return nil, fmt.Errorf("%w: %s", ErrOperationExists, opID)
		}
		return nil, fmt.Errorf("create operation %s: %w", opID, mkErr)
	}
	defer func() {
		if err != nil {
			os.RemoveAll(dir)
		}
	}()

	if mkErr := os.Mkdir(filepath.Join(dir, backupDirName), 0755); mkErr != nil {
		return nil, fmt.Errorf("create operation %s backup dir: %w", opID, mkErr)
	}
	if mkErr := os.Mkdir(filepath.Join(dir, stagingDirName), 0755); mkErr != nil {
		return nil, fmt.Errorf("create operation %s staging dir: %w", opID, mkErr)
	}

	plannedEntries := make([]JournalEntry, len(entries))
	copy(plannedEntries, entries)
	for i := range plannedEntries {
		plannedEntries[i].State = StepPlanned
		plannedEntries[i].InstalledHash = ""
	}

	op = &Operation{
		Dir: dir,
		Journal: Journal{
			OperationID:  opID,
			CreatedAt:    createdAt,
			SourceHashes: append([]ManifestSource{}, sourceHashes...),
			Entries:      plannedEntries,
		},
	}
	if createErr := op.createJournal(); createErr != nil {
		return nil, createErr
	}
	return op, nil
}

// LoadOperation reads an existing operation directory's journal, so a
// resumed process can pick up exactly where the journal last recorded.
func LoadOperation(dir string) (*Operation, error) {
	path := filepath.Join(dir, journalFileName)
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var j Journal
	if err := yaml.Unmarshal(content, &j); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return &Operation{Dir: dir, Journal: j}, nil
}

func (op *Operation) journalPath() string {
	return filepath.Join(op.Dir, journalFileName)
}

// createJournal writes the operation's first operation.yml, create-only: it
// runs exactly once, from CreateOperation, before any step has touched
// anything. Every later persist goes through persistJournal instead.
func (op *Operation) createJournal() error {
	content, err := yaml.Marshal(op.Journal)
	if err != nil {
		return fmt.Errorf("marshal operation.yml: %w", err)
	}
	f, err := os.OpenFile(op.journalPath(), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return fmt.Errorf("create operation.yml: %w", err)
	}
	defer f.Close()
	if _, err := f.Write(content); err != nil {
		return fmt.Errorf("write operation.yml: %w", err)
	}
	return nil
}

// persistJournal durably overwrites operation.yml with the operation's
// current in-memory state, through ReplaceFile rather than a truncating
// write: a crash mid-write must never leave a half-written journal, because a
// half-written journal turns recovery into archaeology instead of a read.
func (op *Operation) persistJournal() error {
	content, err := yaml.Marshal(op.Journal)
	if err != nil {
		return fmt.Errorf("marshal operation.yml: %w", err)
	}
	info, err := os.Stat(op.journalPath())
	if err != nil {
		return fmt.Errorf("stat operation.yml: %w", err)
	}
	if err := ReplaceFile(op.journalPath(), content, info.Mode(), nil); err != nil {
		return fmt.Errorf("persist operation.yml: %w", err)
	}
	return nil
}

// BackupPath is where relPath's pre-operation bytes are copied, inside this
// operation's own directory.
func (op *Operation) BackupPath(relPath string) string {
	return filepath.Join(op.Dir, backupDirName, filepath.FromSlash(relPath))
}

// StagingPath is where relPath's complete new content is written before any
// live path is touched.
func (op *Operation) StagingPath(relPath string) string {
	return filepath.Join(op.Dir, stagingDirName, filepath.FromSlash(relPath))
}

func (op *Operation) entry(relPath string) (*JournalEntry, error) {
	for i := range op.Journal.Entries {
		if op.Journal.Entries[i].Path == relPath {
			return &op.Journal.Entries[i], nil
		}
	}
	return nil, fmt.Errorf("operation %s: no journal entry for %s", op.Journal.OperationID, relPath)
}

func (op *Operation) sourceHash(relPath string) (string, bool) {
	for _, s := range op.Journal.SourceHashes {
		if s.Path == relPath {
			return s.SHA256, true
		}
	}
	return "", false
}

func requireState(opID string, e *JournalEntry, want StepState) error {
	if e.State != want {
		return fmt.Errorf("operation %s: %s is in state %s, want %s before this step", opID, e.Path, e.State, want)
	}
	return nil
}

// Backup copies the live file at root/relPath into this operation's backup/
// directory and re-hashes the copy against the recorded source hash — the
// hash Plan's inventory read before anything changed, not the planned target
// content. A mismatch means the copy itself is untrustworthy (a disk fault
// during the copy, since freshness against a changed source is checked
// earlier, before an operation is even created) and aborts before advancing
// the entry's state, leaving both the live file and the journal exactly as
// they were. Only ActionReplace and ActionRemove entries have a backup step;
// nothing at an ActionCreate destination exists yet to lose.
func (op *Operation) Backup(root, relPath string) error {
	e, err := op.entry(relPath)
	if err != nil {
		return err
	}
	if e.Action != ActionReplace && e.Action != ActionRemove {
		return fmt.Errorf("operation %s: %s has action %s, which has no backup step", op.Journal.OperationID, relPath, e.Action)
	}
	if err := requireState(op.Journal.OperationID, e, StepPlanned); err != nil {
		return err
	}
	wantHash, ok := op.sourceHash(relPath)
	if !ok {
		return fmt.Errorf("operation %s: %s has no recorded source hash to verify a backup against", op.Journal.OperationID, relPath)
	}

	content, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(relPath)))
	if err != nil {
		return fmt.Errorf("backup %s: read source: %w", relPath, err)
	}

	dst := op.BackupPath(relPath)
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("backup %s: create backup dir: %w", relPath, err)
	}
	if err := writeCompleteFile(dst, content, 0644); err != nil {
		return fmt.Errorf("backup %s: %w", relPath, err)
	}

	if got := hashBytes(content); got != wantHash {
		// The copy on disk cannot be trusted; remove it rather than leave a
		// backup a later recovery might mistake for a good one.
		os.Remove(dst)
		return fmt.Errorf("%w: %s backup hash %s, want %s", ErrBackupVerificationFailed, relPath, got, wantHash)
	}

	e.State = StepBackedUp
	return op.persistJournal()
}

// WriteStaged writes content complete to this operation's staging area for
// relPath and advances the entry to StepStaged. Nothing under root is
// touched: staging happens entirely inside the operation directory, so every
// staged file exists before the first live path is touched (AC5). content
// must hash to the entry's recorded planned hash, or WriteStaged refuses
// before writing anything — a mismatch here is a caller error (staging the
// wrong bytes), not a recoverable operation failure.
func (op *Operation) WriteStaged(relPath string, content []byte) error {
	e, err := op.entry(relPath)
	if err != nil {
		return err
	}
	if e.Action == ActionRemove {
		return fmt.Errorf("operation %s: %s has action remove, which has nothing to stage", op.Journal.OperationID, relPath)
	}
	want := StepPlanned
	if e.Action == ActionReplace {
		want = StepBackedUp
	}
	if err := requireState(op.Journal.OperationID, e, want); err != nil {
		return err
	}
	if got := hashBytes(content); e.PlannedHash != "" && got != e.PlannedHash {
		return fmt.Errorf("operation %s: staged content for %s hashes to %s, want planned %s", op.Journal.OperationID, relPath, got, e.PlannedHash)
	}

	dst := op.StagingPath(relPath)
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("stage %s: create staging dir: %w", relPath, err)
	}
	if err := writeCompleteFile(dst, content, 0644); err != nil {
		return fmt.Errorf("stage %s: %w", relPath, err)
	}

	e.State = StepStaged
	return op.persistJournal()
}

// Install applies relPath's staged content to the live project: a create
// writes it to a destination that must not already exist, a replace goes
// through ReplaceFile so the destination is never truncated or missing
// mid-write. A create destination is always a freshly allocated V2 path
// nothing else could have written before this operation, so overwriting a
// leftover from a previous incomplete attempt at the same operation is the
// intended, retry-safe behavior — never a collision with unrelated user
// content.
func (op *Operation) Install(root, relPath string, mode os.FileMode) error {
	e, err := op.entry(relPath)
	if err != nil {
		return err
	}
	if e.Action != ActionCreate && e.Action != ActionReplace {
		return fmt.Errorf("operation %s: %s has action %s, which has no install step", op.Journal.OperationID, relPath, e.Action)
	}
	if err := requireState(op.Journal.OperationID, e, StepStaged); err != nil {
		return err
	}

	staged, err := os.ReadFile(op.StagingPath(relPath))
	if err != nil {
		return fmt.Errorf("install %s: read staged content: %w", relPath, err)
	}

	dest := filepath.Join(root, filepath.FromSlash(relPath))
	switch e.Action {
	case ActionCreate:
		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return fmt.Errorf("install %s: create destination dir: %w", relPath, err)
		}
		if err := writeCompleteFile(dest, staged, mode); err != nil {
			return fmt.Errorf("install %s: %w", relPath, err)
		}
	case ActionReplace:
		if err := ReplaceFile(dest, staged, mode, nil); err != nil {
			return fmt.Errorf("install %s: %w", relPath, err)
		}
	}

	e.InstalledHash = hashBytes(staged)
	e.State = StepInstalled
	return op.persistJournal()
}

// Remove deletes the live file at root/relPath, which must already be backed
// up. It is the ActionRemove counterpart to Install: there is no content to
// install, so a successful removal advances the entry straight to
// StepInstalled, meaning "the action this entry names has been carried out."
// A destination already absent (a rerun after a removal that succeeded but
// whose state update did not persist) is not an error.
func (op *Operation) Remove(root, relPath string) error {
	e, err := op.entry(relPath)
	if err != nil {
		return err
	}
	if e.Action != ActionRemove {
		return fmt.Errorf("operation %s: %s has action %s, which has no remove step", op.Journal.OperationID, relPath, e.Action)
	}
	if err := requireState(op.Journal.OperationID, e, StepBackedUp); err != nil {
		return err
	}

	dest := filepath.Join(root, filepath.FromSlash(relPath))
	if err := os.Remove(dest); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove %s: %w", relPath, err)
	}

	e.State = StepInstalled
	return op.persistJournal()
}

// Verify confirms the outcome of an already-installed entry and advances it
// to StepVerified, the final state: for create/replace, the live file's bytes
// hash to the recorded installed hash; for remove, the live path is confirmed
// gone. PendingOperation treats any entry short of StepVerified as recovery
// work still to do.
func (op *Operation) Verify(root, relPath string) error {
	e, err := op.entry(relPath)
	if err != nil {
		return err
	}
	if err := requireState(op.Journal.OperationID, e, StepInstalled); err != nil {
		return err
	}

	dest := filepath.Join(root, filepath.FromSlash(relPath))
	if e.Action == ActionRemove {
		if _, statErr := os.Lstat(dest); statErr == nil {
			return fmt.Errorf("%w: %s still exists after removal", ErrVerificationFailed, relPath)
		} else if !os.IsNotExist(statErr) {
			return fmt.Errorf("verify %s: %w", relPath, statErr)
		}
	} else {
		content, readErr := os.ReadFile(dest)
		if readErr != nil {
			return fmt.Errorf("verify %s: %w", relPath, readErr)
		}
		if got := hashBytes(content); got != e.InstalledHash {
			return fmt.Errorf("%w: %s hash %s, want %s", ErrVerificationFailed, relPath, got, e.InstalledHash)
		}
	}

	e.State = StepVerified
	return op.persistJournal()
}

func hashBytes(content []byte) string {
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])
}

// writeCompleteFile writes content to path as a single complete file, via a
// temporary file in the same directory that is flushed, closed, and renamed
// into place. Unlike ReplaceFile it freely overwrites an existing
// destination, which is correct here because every caller uses it only for
// the operation's own internal files (staging and backup copies) or for an
// ActionCreate destination that is a freshly allocated path nothing else
// could legitimately occupy — never for replacing a live file a user already
// has. ReplaceFile owns that case.
func writeCompleteFile(path string, content []byte, mode os.FileMode) (err error) {
	temp, err := os.CreateTemp(filepath.Dir(path), replaceTempPattern)
	if err != nil {
		return fmt.Errorf("create temporary file: %w", err)
	}
	tempPath := temp.Name()
	defer func() {
		if err != nil {
			if removeErr := os.Remove(tempPath); removeErr != nil && !os.IsNotExist(removeErr) {
				err = fmt.Errorf("%w (cleanup %s: %v)", err, tempPath, removeErr)
			}
		}
	}()

	if err = temp.Chmod(mode.Perm()); err != nil {
		temp.Close()
		return fmt.Errorf("set temporary permissions: %w", err)
	}
	if _, err = temp.Write(content); err != nil {
		temp.Close()
		return fmt.Errorf("write: %w", err)
	}
	if err = temp.Sync(); err != nil {
		temp.Close()
		return fmt.Errorf("flush: %w", err)
	}
	if err = temp.Close(); err != nil {
		return fmt.Errorf("close: %w", err)
	}
	if err = os.Rename(tempPath, path); err != nil {
		return fmt.Errorf("rename into place: %w", err)
	}
	return nil
}

// PendingOperationReport is what PendingOperation returns for the single
// incomplete operation it found: enough identity and journal state to explain
// what happened and what to do next, without performing any recovery itself.
type PendingOperationReport struct {
	OperationID string
	Dir         string
	Journal     Journal
}

// RecoveryGuidance renders a human-readable explanation of the report's
// state, naming the operation, every path not yet verified, and the recovery
// command — never a bare claim that something went wrong. It distinguishes
// FS-06's preflight no-partial-write expectation (nothing here applies before
// an operation directory exists at all) from a recorded interruption after a
// valid apply began: this guidance only ever describes the latter, because
// PendingOperation only ever finds something once CreateOperation succeeded.
func (r *PendingOperationReport) RecoveryGuidance() string {
	var incomplete []string
	for _, e := range r.Journal.Entries {
		if e.State != StepVerified {
			incomplete = append(incomplete, fmt.Sprintf("%s (%s)", e.Path, e.State))
		}
	}
	sort.Strings(incomplete)

	if len(incomplete) == 0 {
		return fmt.Sprintf(
			"migration operation %s recorded every path verified but did not activate schema_version: 2; "+
				"this is a recorded interruption after a valid apply began, not a partial write — "+
				"run `savepoint migrate --recover` to resume and finish activation",
			r.OperationID,
		)
	}
	return fmt.Sprintf(
		"migration operation %s was interrupted after a valid apply began; the following paths are not yet "+
			"verified and are recoverable from their backup or staged copy: %s; run `savepoint migrate --recover` to resume",
		r.OperationID, strings.Join(incomplete, ", "),
	)
}

// PendingOperation is a read-only detector for an incomplete migration
// operation under projectRoot's .savepoint/.migration/ directory. It returns
// nil, nil when none exists. More than one incomplete operation is refused
// with ErrMultipleOperations naming every operation ID found, rather than
// heuristically picking one.
func PendingOperation(projectRoot string) (*PendingOperationReport, error) {
	rootAbs, err := filepath.Abs(projectRoot)
	if err != nil {
		return nil, err
	}
	opsRoot := operationsRootDir(rootAbs)

	dirEntries, err := os.ReadDir(opsRoot)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("list migration operations: %w", err)
	}

	var ids []string
	for _, e := range dirEntries {
		if e.IsDir() {
			ids = append(ids, e.Name())
		}
	}
	if len(ids) == 0 {
		return nil, nil
	}
	sort.Strings(ids)
	if len(ids) > 1 {
		return nil, fmt.Errorf("%w: %s", ErrMultipleOperations, strings.Join(ids, ", "))
	}

	op, err := LoadOperation(filepath.Join(opsRoot, ids[0]))
	if err != nil {
		return nil, fmt.Errorf("load operation %s: %w", ids[0], err)
	}
	return &PendingOperationReport{OperationID: op.Journal.OperationID, Dir: op.Dir, Journal: op.Journal}, nil
}
