// Package migrate's apply.go performs the ordered publish that commits a
// previously computed ConversionPlan: it is the only place migration writes
// project content. Everything before this file computes; Apply commits.
//
// Apply revalidates the plan's raw source hashes against the live project
// before its first write (the reviewed preview may no longer describe the
// project), then walks a fixed, journalled publish order: back up every
// live path an entry will replace or remove, create additive V2 records,
// write the byte-preserved archive and the v1-to-v2.yml manifest, replace
// modified in-place documents, remove originals whose archive copies have
// verified, and finally activate schema_version: 2 — the one field, in one
// file, that is the actual commit point (see operation.go's package doc and
// activateSchema below for why that last write is deliberately outside the
// per-path journal every other write goes through).
//
// A second call against an incomplete operation resumes it from the durable
// operation record: every complete output is recorded before the first live
// install, so a new process does not need the original clock, operation-ID
// source, decisions file, or in-memory plan.
// Revalidation still refuses changed sources and installed outputs before any
// later live write.
package migrate

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/opencode/savepoint/internal/data"
)

var (
	// ErrPlanConflict means the plan itself reported a Conflict — most
	// commonly an existing v1-to-v2.yml on a project that is not yet at
	// schema_version 2. Apply refuses outright rather than guess which side
	// of the conflict is authoritative.
	ErrPlanConflict = errors.New("migrate: plan reports a conflict that must be resolved before applying")
	// ErrPlanNotAppliable means the plan has at least one unresolved
	// blocking ambiguity. Apply refuses rather than apply the resolved parts
	// of an otherwise-ambiguous plan.
	ErrPlanNotAppliable = errors.New("migrate: plan has unresolved blocking ambiguities")
	// ErrSourceConflict means a source's live bytes no longer hash to what
	// was recorded — either at the original preview (a fresh apply) or at
	// operation creation (a resume) — so the plan being applied no longer
	// describes the project. Apply aborts before its first write.
	ErrSourceConflict = errors.New("migrate: a source changed since it was recorded; the plan no longer describes the project")
	// ErrRecoveryPlanMissing means an old or manually-created operation has no
	// persisted plan and its caller supplied no compatible plan to reconstruct
	// output bytes. A real command-created operation always records the plan.
	ErrRecoveryPlanMissing = errors.New("migrate: pending operation has no recoverable plan")
)

// ApplyResult reports what Apply actually did.
type ApplyResult struct {
	// OperationID names the operation directory the apply ran (or resumed)
	// under. It is empty when AlreadyMigrated is true, since no operation
	// was created or needed.
	OperationID string
	// AlreadyMigrated means plan.SchemaAlreadyV2 was true: the project was
	// already at schema_version 2, and Apply wrote nothing.
	AlreadyMigrated bool
	// Resumed means Apply continued a previously created, incomplete
	// operation rather than creating a new one.
	Resumed bool
}

// Apply commits plan against the project at root. See the package doc above
// for the ordered publish and the resume contract.
func Apply(root string, plan *ConversionPlan) (*ApplyResult, error) {
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	pending, err := PendingOperation(rootAbs)
	if err != nil {
		return nil, err
	}

	resumed := pending != nil
	var op *Operation
	var batch *applyBatch
	if resumed {
		op, err = LoadOperation(pending.Dir)
		if err != nil {
			return nil, err
		}
		if err := revalidatePendingOperation(rootAbs, op); err != nil {
			return nil, err
		}
		batch, err = recoverApplyBatch(rootAbs, op, plan)
	} else {
		if plan == nil {
			return nil, fmt.Errorf("%w: no plan was supplied", ErrRecoveryPlanMissing)
		}
		if plan.SchemaAlreadyV2 {
			return &ApplyResult{AlreadyMigrated: true}, nil
		}
		if len(plan.Conflicts) > 0 {
			c := plan.Conflicts[0]
			return nil, fmt.Errorf("%w: %s: %s", ErrPlanConflict, c.Path, c.Detail)
		}
		if !plan.Appliable {
			return nil, fmt.Errorf("%w: %s", ErrPlanNotAppliable, strings.Join(plan.UnresolvedBlockingIDs, ", "))
		}

		sourcesToCheck := toManifestSources(plan.Sources)
		if err := revalidateSources(rootAbs, sourcesToCheck); err != nil {
			return nil, err
		}
		batch, err = buildApplyBatch(rootAbs, plan)
		if err != nil {
			return nil, err
		}
		op, err = CreateOperationWithPlan(rootAbs, plan.OperationID, sourcesToCheck, batch.journalEntries(), plan.GeneratedAt, plan)
	}
	if err != nil {
		return nil, err
	}

	if err := publish(op, rootAbs, batch); err != nil {
		return nil, err
	}
	if err := activateSchema(rootAbs); err != nil {
		return nil, err
	}
	if err := os.RemoveAll(op.Dir); err != nil {
		return nil, fmt.Errorf("clean up completed operation %s: %w", op.Journal.OperationID, err)
	}

	return &ApplyResult{OperationID: op.Journal.OperationID, Resumed: resumed}, nil
}

// revalidateSources re-reads and re-hashes every recorded source against the
// live project, and aborts with ErrSourceConflict naming every path whose
// bytes differ (including one that vanished entirely) before Apply's first
// write. This is what "source changes invalidate preview" actually means:
// a real hash comparison, never a timestamp check.
func revalidateSources(root string, sources []ManifestSource) error {
	var changed []string
	for _, s := range sources {
		content, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(s.Path)))
		if err != nil || hashBytes(content) != s.SHA256 {
			changed = append(changed, s.Path)
		}
	}
	if len(changed) == 0 {
		return nil
	}
	sort.Strings(changed)
	return fmt.Errorf("%w: %s", ErrSourceConflict, strings.Join(changed, ", "))
}

func toManifestSources(sources []SourceFile) []ManifestSource {
	out := make([]ManifestSource, len(sources))
	for i, s := range sources {
		out[i] = ManifestSource{Path: s.Path, SHA256: s.SHA256}
	}
	return out
}

// revalidatePendingOperation checks only source paths the interrupted
// operation has not intentionally replaced or removed, then checks every
// installed/verified output. Removed sources are expected to be absent after
// their journalled removal; treating them as stale would make a valid
// post-removal recovery impossible.
func revalidatePendingOperation(root string, op *Operation) error {
	byPath := make(map[string]JournalEntry, len(op.Journal.Entries))
	for _, entry := range op.Journal.Entries {
		byPath[entry.Path] = entry
	}

	var sources []ManifestSource
	for _, source := range op.Journal.SourceHashes {
		entry, touched := byPath[source.Path]
		if touched && (entry.Action == ActionRemove || entry.Action == ActionReplace) &&
			(entry.State == StepInstalled || entry.State == StepVerified) {
			continue
		}
		sources = append(sources, source)
	}
	if err := revalidateSources(root, sources); err != nil {
		return err
	}

	for _, entry := range op.Journal.Entries {
		if entry.State != StepInstalled && entry.State != StepVerified {
			continue
		}
		if err := op.verifyRecorded(root, entry.Path); err != nil {
			return err
		}
	}
	return nil
}

// recoverApplyBatch reconstructs the publish order from the operation's
// journal and staged bytes. It does not re-read source files for a normal
// command-created operation, which is essential after archive removals have
// already happened. If an operation was interrupted while preparing its
// staging area, the persisted original plan fills only the missing entries;
// source revalidation has already proved that those reads are still safe.
func recoverApplyBatch(root string, op *Operation, fallback *ConversionPlan) (*applyBatch, error) {
	b, missing, err := batchFromJournal(op)
	if err != nil {
		return nil, err
	}
	if len(missing) == 0 {
		return b, nil
	}

	plan := op.Journal.Plan
	if plan == nil {
		plan = fallback
	}
	if plan == nil {
		return nil, fmt.Errorf("%w: missing staged outputs for %s", ErrRecoveryPlanMissing, strings.Join(missing, ", "))
	}
	rendered, err := buildApplyBatch(root, plan)
	if err != nil {
		return nil, fmt.Errorf("reconstruct pending operation: %w", err)
	}
	renderedByPath := make(map[string]stagedWrite)
	for _, write := range rendered.orderedWrites() {
		renderedByPath[write.Path] = write
	}

	for _, path := range missing {
		write, ok := renderedByPath[path]
		if !ok {
			return nil, fmt.Errorf("%w: journal path %s is absent from the recorded plan", ErrRecoveryPlanMissing, path)
		}
		entry, entryErr := op.entry(path)
		if entryErr != nil {
			return nil, entryErr
		}
		if write.Action != entry.Action || (write.Action != ActionRemove && hashBytes(write.Content) != entry.PlannedHash) {
			return nil, fmt.Errorf("%w: recorded output for %s no longer matches the operation journal", ErrRecoveryPlanMissing, path)
		}
		for i := range b.ordered {
			if b.ordered[i].Path == path {
				b.ordered[i].Content = write.Content
				break
			}
		}
	}
	return b, nil
}

func batchFromJournal(op *Operation) (*applyBatch, []string, error) {
	b := &applyBatch{}
	var missing []string
	for _, entry := range op.Journal.Entries {
		write := stagedWrite{Path: entry.Path, Action: entry.Action}
		if entry.Action != ActionRemove {
			content, err := os.ReadFile(op.StagingPath(entry.Path))
			if os.IsNotExist(err) {
				missing = append(missing, entry.Path)
			} else if err != nil {
				return nil, nil, fmt.Errorf("recover %s: read staged content: %w", entry.Path, err)
			} else if hashBytes(content) != entry.PlannedHash {
				return nil, nil, fmt.Errorf("recover %s: staged content hash does not match the operation journal", entry.Path)
			} else {
				write.Content = content
			}
		}
		b.ordered = append(b.ordered, write)
	}
	sort.Strings(missing)
	return b, missing, nil
}

// stagedWrite is one path's complete planned content, rendered once up
// front so both CreateOperation (which needs every entry's planned hash
// before the first write) and publish (which needs the bytes to stage) read
// it from the same place. Content is nil for an ActionRemove entry, which
// installs nothing.
type stagedWrite struct {
	Path    string
	Action  EntryAction
	Content []byte
}

func (w stagedWrite) journalEntry() JournalEntry {
	e := JournalEntry{Path: w.Path, Action: w.Action}
	if w.Action != ActionRemove {
		e.PlannedHash = hashBytes(w.Content)
	}
	return e
}

// applyBatch is every write Apply's ordered publish performs, grouped and
// ordered exactly as the publish order requires: additive creates, the
// archive and manifest creates, in-place document replaces, and the
// original-removal set that may only proceed once its paired archive entry
// verifies. Schema activation is deliberately not a field here — see
// activateSchema.
type applyBatch struct {
	creates        []stagedWrite // V2 Objective/Task/Issue records, plus Idea.md
	archiveCreates []stagedWrite // .savepoint/archive/v1/* byte-preserved copies
	manifest       stagedWrite   // migrations/v1-to-v2.yml
	replaces       []stagedWrite // router.md, rewritten in place
	removes        []stagedWrite // archived sources, removed after their archive copy verifies
	ordered        []stagedWrite // complete publish order, including the manifest
}

// afterPublishWriteHook is test-only fault injection for proving a real
// command invocation can leave a journal behind and a later invocation can
// recover it. Production callers leave it nil.
var afterPublishWriteHook func(string) error

func (b *applyBatch) journalEntries() []JournalEntry {
	var entries []JournalEntry
	for _, w := range b.orderedWrites() {
		entries = append(entries, w.journalEntry())
	}
	return entries
}

func (b *applyBatch) orderedWrites() []stagedWrite {
	if b.ordered != nil {
		return b.ordered
	}
	var ordered []stagedWrite
	ordered = append(ordered, b.creates...)
	ordered = append(ordered, b.archiveCreates...)
	if b.manifest.Path != "" {
		ordered = append(ordered, b.manifest)
	}
	ordered = append(ordered, b.replaces...)
	ordered = append(ordered, b.removes...)
	return ordered
}

// backupTargets is every entry that touches an existing live path and so
// must be backed up before its first write: every replace and every
// removal. Creates need no backup — nothing at their destination exists yet
// to lose.
func (b *applyBatch) backupTargets() []stagedWrite {
	all := make([]stagedWrite, 0, len(b.orderedWrites()))
	for _, w := range b.orderedWrites() {
		if w.Action == ActionReplace || w.Action == ActionRemove {
			all = append(all, w)
		}
	}
	return all
}

// buildApplyBatch renders every write plan describes, reading only —
// archive content is read fresh from each source's live bytes, and Convert*
// renders every V2 record and document — so the returned batch is a
// complete, self-contained description of everything Apply's ordered
// publish will do.
func buildApplyBatch(root string, plan *ConversionPlan) (*applyBatch, error) {
	b := &applyBatch{}

	for _, t := range plan.Targets {
		content, err := renderTargetContent(root, plan, t)
		if err != nil {
			return nil, fmt.Errorf("render %s: %w", t.GlobalID, err)
		}
		b.creates = append(b.creates, stagedWrite{
			Path:    savepointPath(t.InstallPath()),
			Action:  ActionCreate,
			Content: content,
		})
	}

	for _, d := range plan.Documents {
		content, err := renderDocumentContent(root, plan, d)
		if err != nil {
			return nil, fmt.Errorf("render document %s: %w", d.SourcePath, err)
		}
		w := stagedWrite{Path: savepointPath(d.TargetPath), Content: content}
		if d.Kind == DocumentRouter {
			w.Action = ActionReplace
			b.replaces = append(b.replaces, w)
		} else {
			w.Action = ActionCreate
			b.creates = append(b.creates, w)
		}
	}

	for _, a := range plan.Archives {
		raw, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(a.SourcePath)))
		if err != nil {
			return nil, fmt.Errorf("read archive source %s: %w", a.SourcePath, err)
		}
		b.archiveCreates = append(b.archiveCreates, stagedWrite{Path: a.ArchivePath, Action: ActionCreate, Content: raw})
		b.removes = append(b.removes, stagedWrite{Path: a.SourcePath, Action: ActionRemove})
	}

	manifestContent, err := BuildManifest(plan).Marshal()
	if err != nil {
		return nil, fmt.Errorf("build manifest: %w", err)
	}
	b.manifest = stagedWrite{Path: manifestRelPath(), Action: ActionCreate, Content: manifestContent}
	b.ordered = b.orderedWrites()

	return b, nil
}

func renderTargetContent(root string, plan *ConversionPlan, t PlannedTarget) ([]byte, error) {
	var content string
	var err error
	switch t.Kind {
	case TargetRelease:
		content, err = ConvertRelease(root, plan, t)
	case TargetObjective:
		content, err = ConvertObjective(root, plan, t)
	case TargetTask:
		content, err = ConvertTask(root, plan, t)
	case TargetIssue:
		content, err = ConvertIssue(root, plan, t)
	default:
		return nil, fmt.Errorf("migrate: target %s has unknown kind %q", t.GlobalID, t.Kind)
	}
	if err != nil {
		return nil, err
	}
	return []byte(content), nil
}

func renderDocumentContent(root string, plan *ConversionPlan, d PlannedDocument) ([]byte, error) {
	var content string
	var err error
	switch d.Kind {
	case DocumentIdea:
		content, err = ConvertIdea(root, d)
	case DocumentRouter:
		content, err = ConvertRouter(root, plan, d)
	default:
		return nil, fmt.Errorf("migrate: document %s has unknown kind %q", d.SourcePath, d.Kind)
	}
	if err != nil {
		return nil, err
	}
	return []byte(content), nil
}

func savepointPath(p string) string {
	return filepath.ToSlash(filepath.Join(".savepoint", p))
}

// publish runs the ordered walk: back up everything that will be replaced
// or removed, stage every complete output, install additive records and
// in-place documents, then remove originals now that their archive copies
// have verified. All staging happens before the first live install. Every
// per-entry call below is idempotent over the entry's already-recorded State,
// so calling publish again against a journal left partway through by an
// earlier call does only the remaining work.
func publish(op *Operation, root string, b *applyBatch) error {
	// An installed output may have been edited after the prior journal write
	// but before this recovery call. Check all such paths before staging or
	// installing anything else, so recovery never overwrites a user edit while
	// progressing another entry.
	for _, entry := range op.Journal.Entries {
		if entry.State != StepInstalled && entry.State != StepVerified {
			continue
		}
		if err := op.verifyRecorded(root, entry.Path); err != nil {
			return err
		}
	}
	for _, w := range b.backupTargets() {
		if err := ensureBackedUp(op, root, w.Path); err != nil {
			return err
		}
	}
	for _, w := range b.orderedWrites() {
		if w.Action == ActionRemove {
			continue
		}
		if err := stageWrite(op, w); err != nil {
			return err
		}
	}
	for _, w := range b.orderedWrites() {
		if w.Action == ActionRemove {
			continue
		}
		if err := installStagedWrite(op, root, w); err != nil {
			return err
		}
		if afterPublishWriteHook != nil {
			if err := afterPublishWriteHook(w.Path); err != nil {
				return err
			}
		}
	}
	for _, w := range b.orderedWrites() {
		if w.Action != ActionRemove {
			continue
		}
		if err := publishRemoval(op, root, w.Path); err != nil {
			return err
		}
	}
	return nil
}

func ensureBackedUp(op *Operation, root, relPath string) error {
	e, err := op.entry(relPath)
	if err != nil {
		return err
	}
	if e.State == StepPlanned {
		return op.Backup(root, relPath)
	}
	return nil
}

// publishWrite carries one create or replace entry from wherever its State
// already stands through staged, installed, and verified. Skipping a step
// already recorded as done is what makes a resumed call safe: WriteStaged
// additionally refuses if content does not hash to the entry's frozen
// planned hash, so a resume whose re-rendered content actually differs from
// the interrupted attempt's fails loudly here rather than silently
// installing the wrong bytes.
func publishWrite(op *Operation, root string, w stagedWrite) error {
	if err := stageWrite(op, w); err != nil {
		return err
	}
	return installStagedWrite(op, root, w)
}

func stageWrite(op *Operation, w stagedWrite) error {
	e, err := op.entry(w.Path)
	if err != nil {
		return err
	}
	if e.State == StepPlanned || e.State == StepBackedUp {
		if err := op.WriteStaged(w.Path, w.Content); err != nil {
			return err
		}
	}
	return nil
}

func installStagedWrite(op *Operation, root string, w stagedWrite) error {
	e, err := op.entry(w.Path)
	if err != nil {
		return err
	}
	if e.State == StepStaged {
		mode := os.FileMode(0644)
		if w.Action == ActionReplace {
			if info, statErr := os.Stat(filepath.Join(root, filepath.FromSlash(w.Path))); statErr == nil {
				mode = info.Mode()
			}
		}
		if err := op.Install(root, w.Path, mode); err != nil {
			return err
		}
	}
	if e.State == StepInstalled {
		if err := op.Verify(root, w.Path); err != nil {
			return err
		}
	}
	if e.State == StepVerified {
		if err := op.verifyRecorded(root, w.Path); err != nil {
			return err
		}
	}
	return nil
}

func publishRemoval(op *Operation, root, relPath string) error {
	e, err := op.entry(relPath)
	if err != nil {
		return err
	}
	if e.State == StepBackedUp {
		if err := op.Remove(root, relPath); err != nil {
			return err
		}
	}
	if e.State == StepInstalled {
		if err := op.Verify(root, relPath); err != nil {
			return err
		}
	}
	return nil
}

// activateSchema performs migration's actual commit: setting
// schema_version: 2 in config.yml. It is deliberately outside the
// operation's per-path journal — every other write is backed up, staged,
// and verified because a crash could otherwise lose or duplicate it, but
// this is a single field in a single file, replaced through the same
// temp-file-then-rename primitive (ReplaceFile) that already guarantees the
// destination is never left partially written: before this call succeeds, a
// reader sees a valid V1 project; after it, a valid V2 one. It is also
// idempotent — a config.yml this call already activated, or one a prior
// interrupted attempt activated just before crashing, is left untouched —
// which is what lets PendingOperationReport.RecoveryGuidance's "every path
// verified but not yet activated" case describe a real, recoverable, named
// state rather than a guess.
func activateSchema(root string) error {
	path := filepath.Join(root, ".savepoint", "config.yml")

	version, err := data.ReadSchemaVersion(path)
	if err != nil {
		return err
	}
	if version == data.SchemaVersionV2 {
		return nil
	}

	existing, statErr := os.ReadFile(path)
	if statErr != nil && !os.IsNotExist(statErr) {
		return fmt.Errorf("read %s: %w", path, statErr)
	}

	activated, err := data.ActivateSchemaV2(existing)
	if err != nil {
		return err
	}

	if os.IsNotExist(statErr) {
		return writeCompleteFile(path, activated, 0644)
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}
	return ReplaceFile(path, activated, info.Mode(), nil)
}
