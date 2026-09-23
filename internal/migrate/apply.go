package migrate

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/opencode/savepoint/internal/data"
)

var (
	ErrPlanConflict            = errors.New("migrate: plan reports a conflict that must be resolved before applying")
	ErrPlanNotAppliable        = errors.New("migrate: plan has unresolved blocking ambiguities")
	ErrCreateDestinationExists = errors.New("migrate: create destination already exists")
)

// EntryAction names the project-file change planned for one path.
type EntryAction string

const (
	ActionCreate  EntryAction = "create"
	ActionReplace EntryAction = "replace"
	ActionRemove  EntryAction = "remove"
)

// ApplyResult reports whether Apply found an already-migrated project.
type ApplyResult struct {
	AlreadyMigrated bool
	WrittenPaths    []string
}

type stagedWrite struct {
	Path    string
	Action  EntryAction
	Content []byte
}

type applyBatch struct {
	creates       []stagedWrite
	archiveCopies []stagedWrite
	replaces      []stagedWrite
	removes       []stagedWrite
	manifest      stagedWrite
}

// afterApplyMutationHook is test-only fault injection for the partial-write
// report. Production callers leave it nil.
var afterApplyMutationHook func(string) error

// Apply renders the plan before writing, then installs its outputs, removes
// archived V1 sources, writes the manifest, and activates schema version 2
// as the final project change. It relies on RunCommand's clean Git check for
// the user-facing safety boundary.
func Apply(root string, plan *ConversionPlan) (*ApplyResult, error) {
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	if plan == nil {
		return nil, fmt.Errorf("migrate: no conversion plan was supplied")
	}
	if plan.SchemaAlreadyV2 {
		return &ApplyResult{AlreadyMigrated: true}, nil
	}
	if len(plan.Conflicts) > 0 {
		conflict := plan.Conflicts[0]
		return nil, fmt.Errorf("%w: %s: %s", ErrPlanConflict, conflict.Path, conflict.Detail)
	}
	if !plan.Appliable {
		return nil, fmt.Errorf("%w: %s", ErrPlanNotAppliable, strings.Join(plan.UnresolvedBlockingIDs, ", "))
	}

	configPath := filepath.Join(rootAbs, ".savepoint", "config.yml")
	version, err := data.ReadSchemaVersion(configPath)
	if err != nil {
		return nil, err
	}
	if version == data.SchemaVersionV2 {
		return &ApplyResult{AlreadyMigrated: true}, nil
	}

	batch, err := buildApplyBatch(rootAbs, plan)
	if err != nil {
		return nil, applyFailure(plan, nil, err)
	}
	var written []string
	write := func(item stagedWrite) error {
		if err := writePlannedFile(rootAbs, item); err != nil {
			return fmt.Errorf("%s: %w", item.Path, err)
		}
		written = append(written, item.Path)
		if afterApplyMutationHook != nil {
			if err := afterApplyMutationHook(item.Path); err != nil {
				return fmt.Errorf("%s: %w", item.Path, err)
			}
		}
		return nil
	}
	for _, item := range batch.creates {
		if err := write(item); err != nil {
			return nil, applyFailure(plan, written, err)
		}
	}
	for _, item := range batch.archiveCopies {
		if err := write(item); err != nil {
			return nil, applyFailure(plan, written, err)
		}
	}
	for _, item := range batch.replaces {
		if err := write(item); err != nil {
			return nil, applyFailure(plan, written, err)
		}
	}
	for _, item := range batch.removes {
		path := filepath.Join(rootAbs, filepath.FromSlash(item.Path))
		if err := os.Remove(path); err != nil {
			return nil, applyFailure(plan, written, fmt.Errorf("remove %s: %w", item.Path, err))
		}
		written = append(written, item.Path)
		if afterApplyMutationHook != nil {
			if err := afterApplyMutationHook(item.Path); err != nil {
				return nil, applyFailure(plan, written, fmt.Errorf("%s: %w", item.Path, err))
			}
		}
	}
	if err := write(batch.manifest); err != nil {
		return nil, applyFailure(plan, written, err)
	}

	if err := activateSchema(rootAbs); err != nil {
		return nil, applyFailure(plan, written, fmt.Errorf(".savepoint/config.yml: %w", err))
	}
	written = append(written, ".savepoint/config.yml")
	return &ApplyResult{WrittenPaths: append([]string(nil), written...)}, nil
}

func applyFailure(plan *ConversionPlan, written []string, cause error) error {
	if len(written) == 0 {
		return fmt.Errorf("migrate: apply stopped: %w\nwritten paths: none\nundo from the project root: %s", cause, gitUndoCommand(plan))
	}
	return fmt.Errorf("migrate: apply stopped: %w\nwritten paths:\n - %s\nundo from the project root: %s",
		cause, strings.Join(uniqueSortedPaths(written), "\n - "), gitUndoCommand(plan))
}

func buildApplyBatch(root string, plan *ConversionPlan) (*applyBatch, error) {
	batch := &applyBatch{}
	for _, target := range plan.Targets {
		content, err := renderTargetContent(root, plan, target)
		if err != nil {
			return nil, fmt.Errorf("render %s: %w", target.GlobalID, err)
		}
		batch.creates = append(batch.creates, stagedWrite{
			Path: savepointPath(target.InstallPath()), Action: ActionCreate, Content: content,
		})
	}
	for _, document := range plan.Documents {
		content, err := renderDocumentContent(root, plan, document)
		if err != nil {
			return nil, fmt.Errorf("render document %s: %w", document.SourcePath, err)
		}
		item := stagedWrite{Path: savepointPath(document.TargetPath), Content: content}
		if document.Kind == DocumentRouter {
			item.Action = ActionReplace
			batch.replaces = append(batch.replaces, item)
		} else {
			item.Action = ActionCreate
			batch.creates = append(batch.creates, item)
		}
	}
	for _, archive := range plan.Archives {
		raw, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(archive.SourcePath)))
		if err != nil {
			return nil, fmt.Errorf("read archive source %s: %w", archive.SourcePath, err)
		}
		batch.archiveCopies = append(batch.archiveCopies, stagedWrite{
			Path: archive.ArchivePath, Action: ActionCreate, Content: raw,
		})
		batch.removes = append(batch.removes, stagedWrite{Path: archive.SourcePath, Action: ActionRemove})
	}
	manifestContent, err := BuildManifest(plan).Marshal()
	if err != nil {
		return nil, fmt.Errorf("build manifest: %w", err)
	}
	batch.manifest = stagedWrite{Path: manifestRelPath(), Action: ActionCreate, Content: manifestContent}
	return batch, nil
}

func renderTargetContent(root string, plan *ConversionPlan, target PlannedTarget) ([]byte, error) {
	var content string
	var err error
	switch target.Kind {
	case TargetRelease:
		content, err = ConvertRelease(root, plan, target)
	case TargetObjective:
		content, err = ConvertObjective(root, plan, target)
	case TargetTask:
		content, err = ConvertTask(root, plan, target)
	case TargetIssue:
		content, err = ConvertIssue(root, plan, target)
	default:
		return nil, fmt.Errorf("migrate: target %s has unknown kind %q", target.GlobalID, target.Kind)
	}
	if err != nil {
		return nil, err
	}
	return []byte(content), nil
}

func renderDocumentContent(root string, plan *ConversionPlan, document PlannedDocument) ([]byte, error) {
	var content string
	var err error
	switch document.Kind {
	case DocumentIdea:
		content, err = ConvertIdea(root, document)
	case DocumentRouter:
		content, err = ConvertRouter(root, plan, document)
	default:
		return nil, fmt.Errorf("migrate: document %s has unknown kind %q", document.SourcePath, document.Kind)
	}
	if err != nil {
		return nil, err
	}
	return []byte(content), nil
}

func savepointPath(path string) string {
	return filepath.ToSlash(filepath.Join(".savepoint", path))
}

func writePlannedFile(root string, item stagedWrite) error {
	path := filepath.Join(root, filepath.FromSlash(item.Path))
	mode := os.FileMode(0644)
	if item.Action == ActionReplace {
		info, err := os.Lstat(path)
		if err != nil {
			return fmt.Errorf("inspect replacement destination: %w", err)
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("replacement destination is not a regular file")
		}
		mode = info.Mode().Perm()
	} else if item.Action != ActionCreate {
		return fmt.Errorf("unsupported write action %q", item.Action)
	}
	return writeTempAndRename(path, item.Content, mode, item.Action == ActionCreate)
}

func writeTempAndRename(path string, content []byte, mode os.FileMode, createOnly bool) (retErr error) {
	if createOnly {
		if _, err := os.Lstat(path); err == nil {
			return fmt.Errorf("%w: %s", ErrCreateDestinationExists, path)
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("inspect create destination %s: %w", path, err)
		}
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create parent directory: %w", err)
	}
	temp, err := os.CreateTemp(filepath.Dir(path), ".savepoint-migrate-apply-*")
	if err != nil {
		return fmt.Errorf("create temporary file: %w", err)
	}
	tempPath := temp.Name()
	defer func() {
		if removeErr := os.Remove(tempPath); removeErr != nil && !os.IsNotExist(removeErr) && retErr == nil {
			retErr = fmt.Errorf("remove temporary file %s: %w", tempPath, removeErr)
		}
	}()
	if err := temp.Chmod(mode.Perm()); err != nil {
		_ = temp.Close()
		return fmt.Errorf("set temporary file permissions: %w", err)
	}
	if _, err := temp.Write(content); err != nil {
		_ = temp.Close()
		return fmt.Errorf("write temporary file: %w", err)
	}
	if err := temp.Sync(); err != nil {
		_ = temp.Close()
		return fmt.Errorf("sync temporary file: %w", err)
	}
	if err := temp.Close(); err != nil {
		return fmt.Errorf("close temporary file: %w", err)
	}
	if createOnly {
		if _, err := os.Lstat(path); err == nil {
			return fmt.Errorf("%w: %s", ErrCreateDestinationExists, path)
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("inspect create destination %s: %w", path, err)
		}
	}
	if err := os.Rename(tempPath, path); err != nil {
		return fmt.Errorf("rename temporary file into place: %w", err)
	}
	return nil
}

func activateSchema(root string) error {
	path := filepath.Join(root, ".savepoint", "config.yml")
	version, err := data.ReadSchemaVersion(path)
	if err != nil {
		return err
	}
	if version == data.SchemaVersionV2 {
		return nil
	}
	existing, err := os.ReadFile(path)
	createOnly := false
	mode := os.FileMode(0644)
	if os.IsNotExist(err) {
		createOnly = true
		existing = nil
	} else if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	activated, err := data.ActivateSchemaV2(existing)
	if err != nil {
		return err
	}
	if !createOnly {
		info, err := os.Lstat(path)
		if err != nil {
			return fmt.Errorf("inspect %s: %w", path, err)
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("%s is not a regular file", path)
		}
		mode = info.Mode().Perm()
	}
	return writeTempAndRename(path, activated, mode, createOnly)
}
