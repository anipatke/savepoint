package data

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

var ErrMtimeConflict = fmt.Errorf("file modified since last read")

// SplitFrontmatterBody splits content into frontmatter YAML and body.
func SplitFrontmatterBody(content string) (yamlStr string, body string, err error) {
	normalized := normalizeLineEndings(content)
	if !strings.HasPrefix(normalized, "---\n") {
		return "", "", ErrNoFrontmatter
	}
	end := strings.Index(normalized[4:], "\n---")
	if end == -1 {
		return "", "", ErrNoClosingFrontmatter
	}
	yamlStr = strings.TrimSpace(normalized[4 : 4+end])
	bodyStart := 4 + end + 4 // "---\n" + yaml + "\n---"
	body = ""
	if bodyStart < len(normalized) {
		body = normalized[bodyStart:]
	}
	return yamlStr, body, nil
}

func setMappingField(mapping *yaml.Node, key, value string) {
	for i := 0; i < len(mapping.Content)-1; i += 2 {
		if mapping.Content[i].Value == key {
			mapping.Content[i+1].Value = value
			mapping.Content[i+1].Tag = "!!str"
			return
		}
	}
	keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: key, Tag: "!!str"}
	valNode := &yaml.Node{Kind: yaml.ScalarNode, Value: value, Tag: "!!str"}
	mapping.Content = append(mapping.Content, keyNode, valNode)
}

func removeMappingField(mapping *yaml.Node, key string) bool {
	for i := 0; i < len(mapping.Content)-1; i += 2 {
		if mapping.Content[i].Value == key {
			mapping.Content = append(mapping.Content[:i], mapping.Content[i+2:]...)
			return true
		}
	}
	return false
}

func mappingFieldValue(mapping *yaml.Node, key string) (string, bool) {
	for i := 0; i < len(mapping.Content)-1; i += 2 {
		if mapping.Content[i].Value == key {
			return mapping.Content[i+1].Value, true
		}
	}
	return "", false
}

// v2FieldPatch is one owned frontmatter field a managed V2 write may change.
// Value patches a scalar field. Node patches a nested mapping field (an
// already-encoded sub-block, such as evidence's freshness or
// owner_validation) and takes precedence over Value when set. Remove deletes
// the key entirely — used for Task stage, which must not appear at all
// outside status: in_progress, and for absent evidence sub-blocks — rather
// than writing an empty value.
type v2FieldPatch struct {
	Key    string
	Value  string
	Node   *yaml.Node
	Remove bool
}

// patchV2Mapping applies patches to mapping and reports whether any patch
// actually changed an existing value or key presence. It only ever touches
// the keys named by patches via setMappingField/setMappingNode/
// removeMappingField, so every other mapping entry passes through untouched.
func patchV2Mapping(mapping *yaml.Node, patches []v2FieldPatch) bool {
	changed := false
	for _, patch := range patches {
		if patch.Remove {
			if removeMappingField(mapping, patch.Key) {
				changed = true
			}
			continue
		}
		if patch.Node != nil {
			if existing, ok := mappingFieldNode(mapping, patch.Key); ok && nodesEqualByEncoding(existing, patch.Node) {
				continue
			}
			setMappingNode(mapping, patch.Key, patch.Node)
			changed = true
			continue
		}
		if value, ok := mappingFieldValue(mapping, patch.Key); ok && value == patch.Value {
			continue
		}
		setMappingField(mapping, patch.Key, patch.Value)
		changed = true
	}
	return changed
}

// setMappingNode sets key's value to value, a fully-formed nested node, in
// place — used for structured evidence sub-blocks that setMappingField's
// scalar-only assignment cannot express. It replaces an existing value node
// wholesale rather than merging into it, matching the "patch only the named
// nested fields" contract: a sub-block is written or removed as a whole.
func setMappingNode(mapping *yaml.Node, key string, value *yaml.Node) {
	for i := 0; i < len(mapping.Content)-1; i += 2 {
		if mapping.Content[i].Value == key {
			mapping.Content[i+1] = value
			return
		}
	}
	keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: key, Tag: "!!str"}
	mapping.Content = append(mapping.Content, keyNode, value)
}

func mappingFieldNode(mapping *yaml.Node, key string) (*yaml.Node, bool) {
	for i := 0; i < len(mapping.Content)-1; i += 2 {
		if mapping.Content[i].Value == key {
			return mapping.Content[i+1], true
		}
	}
	return nil, false
}

// nodesEqualByEncoding reports whether a and b marshal to identical YAML
// once source formatting is normalized away, which is how patchV2Mapping
// detects a nested sub-block patch that changes nothing. A node read from
// disk keeps the quote/flow style its author wrote (e.g. single-quoted
// timestamps); a freshly encoded replacement node never carries that style.
// Comparing raw marshal output would report a style difference as a content
// change, so both sides are style-normalized first and marshalled again,
// which still quotes a value exactly when its tag requires it.
func nodesEqualByEncoding(a, b *yaml.Node) bool {
	if a == nil || b == nil {
		return a == b
	}
	encodedA, errA := yaml.Marshal(styleNormalizedV2Node(a))
	encodedB, errB := yaml.Marshal(styleNormalizedV2Node(b))
	if errA != nil || errB != nil {
		return false
	}
	return string(encodedA) == string(encodedB)
}

// styleNormalizedV2Node deep-copies n with every node's Style cleared so
// nodesEqualByEncoding compares semantic YAML content instead of incidental
// source formatting.
func styleNormalizedV2Node(n *yaml.Node) *yaml.Node {
	clone := cloneYAMLNode(n)
	clearV2NodeStyle(clone)
	return clone
}

func clearV2NodeStyle(n *yaml.Node) {
	if n == nil {
		return
	}
	n.Style = 0
	for _, child := range n.Content {
		clearV2NodeStyle(child)
	}
}

// encodeV2Node encodes v (a frontmatter-shaped struct) into a standalone
// yaml.Node, giving a managed V2 write a nested sub-block value it can hand
// to setMappingNode.
func encodeV2Node(v any) (*yaml.Node, error) {
	var node yaml.Node
	if err := node.Encode(v); err != nil {
		return nil, err
	}
	return &node, nil
}

// cloneYAMLNode deep-copies a yaml.Node so a managed V2 write can patch and
// marshal a scratch copy of a record's retained frontmatter without ever
// mutating the loaded record in memory, even when validation later refuses
// the patched content.
func cloneYAMLNode(n *yaml.Node) *yaml.Node {
	if n == nil {
		return nil
	}
	clone := *n
	if n.Content != nil {
		clone.Content = make([]*yaml.Node, len(n.Content))
		for i, child := range n.Content {
			clone.Content[i] = cloneYAMLNode(child)
		}
	}
	return &clone
}

func resolveV2SourcePath(source V2SourceDocument) (string, error) {
	if source.Path == "" {
		return "", fmt.Errorf("%w: empty V2 source path", ErrV2UnsafePath)
	}
	if source.ProjectRoot == "" {
		return filepath.Clean(source.Path), nil
	}

	root, err := filepath.Abs(source.ProjectRoot)
	if err != nil {
		return "", fmt.Errorf("resolve V2 project root %q: %w", source.ProjectRoot, err)
	}
	root = filepath.Clean(root)
	path := filepath.Clean(source.Path)
	if !filepath.IsAbs(path) {
		path = filepath.Join(root, path)
	}
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return "", fmt.Errorf("%w: %s escapes project root", ErrV2UnsafePath, source.Path)
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("%w: %s escapes project root", ErrV2UnsafePath, source.Path)
	}
	return path, nil
}

func checkV2SourceFresh(source V2SourceDocument, path string) (os.FileInfo, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, fmt.Errorf("stat %s: %w", path, err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("%w: %s is a symlink", ErrV2UnsafePath, path)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("write %s: target is a directory", path)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	if source.contentHashSet && sha256.Sum256(content) != source.contentHash {
		return nil, fmt.Errorf("%w: %w: %s", ErrV2SourceConflict, ErrMtimeConflict, path)
	}
	return info, nil
}

// replaceV2File writes content to a same-directory temporary file and
// replaces path only after the complete temporary file has been flushed and
// closed. The finalCheck runs after the temp file is durable and immediately
// before the rename, closing the load/replace freshness window as far as the
// platform's rename primitive permits.
func replaceV2File(path string, content []byte, mode os.FileMode, finalCheck func() error) (retErr error) {
	temp, err := os.CreateTemp(filepath.Dir(path), ".savepoint-v2-write-*")
	if err != nil {
		return fmt.Errorf("write %s: create temporary file: %w", path, err)
	}
	tempPath := temp.Name()
	defer func() {
		if tempPath == "" {
			return
		}
		if cleanupErr := os.Remove(tempPath); cleanupErr != nil && !os.IsNotExist(cleanupErr) {
			if retErr != nil {
				retErr = fmt.Errorf("%w (cleanup %s: %v)", retErr, tempPath, cleanupErr)
			} else {
				retErr = fmt.Errorf("cleanup %s: %w", tempPath, cleanupErr)
			}
		}
	}()

	if err := temp.Chmod(mode.Perm()); err != nil {
		closeErr := temp.Close()
		if closeErr != nil {
			return fmt.Errorf("write %s: set temporary permissions: %w (close: %v)", path, err, closeErr)
		}
		return fmt.Errorf("write %s: set temporary permissions: %w", path, err)
	}

	written, writeErr := temp.Write(content)
	if writeErr != nil || written != len(content) {
		closeErr := temp.Close()
		if writeErr == nil {
			writeErr = fmt.Errorf("short write: wrote %d of %d bytes", written, len(content))
		}
		if closeErr != nil {
			return fmt.Errorf("write %s: %w (close: %v)", path, writeErr, closeErr)
		}
		return fmt.Errorf("write %s: %w", path, writeErr)
	}

	if err := temp.Sync(); err != nil {
		closeErr := temp.Close()
		if closeErr != nil {
			return fmt.Errorf("write %s: flush temporary file: %w (close: %v)", path, err, closeErr)
		}
		return fmt.Errorf("write %s: flush temporary file: %w", path, err)
	}
	if err := temp.Close(); err != nil {
		return fmt.Errorf("write %s: close temporary file: %w", path, err)
	}

	if finalCheck != nil {
		if err := finalCheck(); err != nil {
			return err
		}
	}
	if err := os.Rename(tempPath, path); err != nil {
		return fmt.Errorf("write %s: replace file: %w", path, err)
	}
	tempPath = ""
	return nil
}

// writeV2Record is the shared no-op detection, preservation, and
// validate-before-replace boundary behind WriteObjectiveV2 and WriteTaskV2.
// Patches apply to a clone of source's retained document, so a validation
// failure never mutates the caller's loaded record, and the file is left
// completely untouched — no marshal, no write, no modified-time change —
// when no patch changes an existing value.
func writeV2Record(source *V2SourceDocument, patches []v2FieldPatch, validate func(content string) error) error {
	if source == nil {
		return fmt.Errorf("write V2 record: nil source")
	}
	path, err := resolveV2SourcePath(*source)
	if err != nil {
		return err
	}

	cloned := cloneYAMLNode(&source.Frontmatter)
	if cloned == nil || cloned.Kind != yaml.DocumentNode || len(cloned.Content) == 0 {
		return fmt.Errorf("%s: unexpected yaml structure", path)
	}

	mapping := cloned.Content[0]
	if mapping.Kind != yaml.MappingNode {
		return fmt.Errorf("%s: frontmatter is not a mapping", path)
	}

	if !patchV2Mapping(mapping, patches) {
		return nil
	}

	out, err := yaml.Marshal(cloned)
	if err != nil {
		return fmt.Errorf("%s: marshal yaml: %w", path, err)
	}

	newContent := "---\n" + strings.TrimSpace(string(out)) + "\n---" + source.Body
	if source.CRLF {
		newContent = strings.ReplaceAll(newContent, "\n", "\r\n")
	}

	if err := validate(newContent); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}

	info, err := checkV2SourceFresh(*source, path)
	if err != nil {
		return err
	}
	if info.Mode().Perm()&0222 == 0 {
		return fmt.Errorf("write %s: %w", path, os.ErrPermission)
	}
	if err := replaceV2File(path, []byte(newContent), info.Mode(), func() error {
		latest, err := checkV2SourceFresh(*source, path)
		if err != nil {
			return err
		}
		if latest.Mode().Perm()&0222 == 0 {
			return fmt.Errorf("write %s: %w", path, os.ErrPermission)
		}
		return nil
	}); err != nil {
		return err
	}

	// Keep a successfully managed record usable for another write in the
	// same process. Failed validation, freshness checks, and replacement leave
	// the caller's retained document untouched.
	source.Frontmatter = *cloned
	source.contentHash = sha256.Sum256([]byte(newContent))
	source.contentHashSet = true
	return nil
}

// WriteObjectiveV2 patches only the status field of a V2 Objective record's
// frontmatter to match objective.Status, preserving every other YAML
// key/value and the authored Markdown body unchanged. The patched content is
// validated as a decodable ObjectiveV2 before any file is replaced, so a
// malformed or unsupported status is refused and the file is left
// untouched. Writing the status the file already has is a no-op.
func WriteObjectiveV2(objective *ObjectiveV2) error {
	return writeV2Record(&objective.Source, []v2FieldPatch{
		{Key: "status", Value: string(objective.Status)},
	}, func(content string) error {
		_, err := DecodeObjectiveV2(objective.Source.Path, content)
		return err
	})
}

// WriteObjectiveGroupOrderV2 assigns priority and one-based rank values to
// the supplied ordered Objective IDs in a Goal. It validates membership and
// freshness for every source before writing any file. Replacements are
// atomic per Objective; if a later replacement fails, earlier successful
// writes remain loadable and duplicate ranks are reported by the index.
func WriteObjectiveGroupOrderV2(index *V2Index, goalID string, priority ObjectivePriority, objectiveIDs []string) error {
	if index == nil {
		return fmt.Errorf("write Objective order: nil V2 index")
	}
	if _, ok := index.Releases[goalID]; !ok {
		return fmt.Errorf("%w: Goal %s was not found", ErrV2MissingRelease, goalID)
	}
	if !isCanonicalObjectivePriority(priority) {
		return fmt.Errorf("%w: Objective priority %q; use critical, high, medium, or low", ErrV2Malformed, priority)
	}

	objectives := make([]*ObjectiveV2, len(objectiveIDs))
	seen := make(map[string]struct{}, len(objectiveIDs))
	for i, objectiveID := range objectiveIDs {
		if _, duplicate := seen[objectiveID]; duplicate {
			return fmt.Errorf("write Objective order: Objective %s appears more than once", objectiveID)
		}
		seen[objectiveID] = struct{}{}
		objective, ok := index.Objectives[objectiveID]
		if !ok || objective == nil {
			return fmt.Errorf("%w: Objective %s was not found", ErrV2InvalidID, objectiveID)
		}
		if string(objective.Release) != goalID {
			return fmt.Errorf("%w: Objective %s belongs to Goal %s, not Goal %s", ErrV2InvalidReleaseReference, objectiveID, objective.Release, goalID)
		}
		objectives[i] = objective
	}

	// Preflight every member before the first replacement. This prevents a
	// known stale source later in the order from leaving earlier ranks changed.
	for _, objective := range objectives {
		path, err := resolveV2SourcePath(objective.Source)
		if err != nil {
			return err
		}
		if _, err := checkV2SourceFresh(objective.Source, path); err != nil {
			return err
		}
	}

	for i, objective := range objectives {
		rank := i + 1
		if objective.Priority == priority && objective.Rank == rank {
			continue
		}
		patches := []v2FieldPatch{
			{Key: "priority", Node: &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: string(priority)}},
			{Key: "rank", Node: &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: fmt.Sprint(rank)}},
		}
		if err := writeV2Record(&objective.Source, patches, func(content string) error {
			updated, err := DecodeObjectiveV2(objective.Source.Path, content)
			if err != nil {
				return err
			}
			if string(updated.Release) != goalID {
				return fmt.Errorf("%w: Objective %s does not belong to Goal %s", ErrV2InvalidReleaseReference, objective.ID, goalID)
			}
			return nil
		}); err != nil {
			return fmt.Errorf("write Objective %s order: %w", objective.ID, err)
		}
		objective.Priority = priority
		objective.Rank = rank
	}
	return nil
}

// WriteTaskV2 patches only the status and stage fields of a V2 Task record's
// frontmatter to match task.Status and task.Stage, preserving every other
// YAML key/value (including depends_on and unknown fields) and the authored
// Markdown body unchanged. An empty task.Stage removes the stage key rather
// than writing it empty, matching the V2 rule that stage only applies to
// status: in_progress. The patched content is validated as a decodable
// TaskV2 before any file is replaced, so an unsafe or malformed lifecycle
// value is refused and the file is left untouched.
func WriteTaskV2(task *TaskV2) error {
	patches := []v2FieldPatch{
		{Key: "status", Value: string(task.Status)},
	}
	if task.Stage == "" {
		patches = append(patches, v2FieldPatch{Key: "stage", Remove: true})
	} else {
		patches = append(patches, v2FieldPatch{Key: "stage", Value: string(task.Stage)})
	}

	return writeV2Record(&task.Source, patches, func(content string) error {
		_, err := DecodeTaskV2(task.Source.Path, content)
		return err
	})
}

// WriteTaskEvidenceV2 patches only the evidence fields — last_check,
// freshness, owner_validation, exception, replan, and check_waiver — of a V2
// Task record's frontmatter to match task.Evidence, preserving every other YAML
// key/value (including status, stage, depends_on, and unknown fields) and
// the authored Markdown body unchanged. A nil Evidence, or a nil sub-block
// within it, removes that key rather than writing an empty value — clearing
// a replan flag goes through this path. The patched content is validated as
// a decodable TaskV2 before any file is replaced, so a malformed evidence
// value is refused and the file is left untouched, and writing values the
// record already has is a no-op.
func WriteTaskEvidenceV2(task *TaskV2) error {
	patches, err := evidencePatches(task.Evidence)
	if err != nil {
		return err
	}

	return writeV2Record(&task.Source, patches, func(content string) error {
		_, err := DecodeTaskV2(task.Source.Path, content)
		return err
	})
}

// WriteObjectiveEvidenceV2 patches only the evidence fields — last_check,
// freshness, owner_validation, exception, and replan — of a V2 Objective
// record's frontmatter to match objective.Evidence, through the same
// evidencePatches builder WriteTaskEvidenceV2 uses, preserving every other
// YAML key/value (including status, depends_on, and unknown fields) and the
// authored Markdown body unchanged. The patched content is validated as a
// decodable ObjectiveV2 before any file is replaced.
func WriteObjectiveEvidenceV2(objective *ObjectiveV2) error {
	patches, err := evidencePatches(objective.Evidence)
	if err != nil {
		return err
	}

	return writeV2Record(&objective.Source, patches, func(content string) error {
		_, err := DecodeObjectiveV2(objective.Source.Path, content)
		return err
	})
}

// evidencePatches builds the shared last_check/freshness/owner_validation/
// exception/replan patches for the Evidence block common to Task and
// Objective records, so WriteTaskEvidenceV2 and WriteObjectiveEvidenceV2
// have exactly one encoding of it between them.
func evidencePatches(evidence *Evidence) ([]v2FieldPatch, error) {
	var lastCheck *string
	var freshness *Freshness
	var ownerValidation *OwnerValidation
	var exception *Exception
	var replan *Replan
	var checkWaiver *CheckWaiver
	if evidence != nil {
		if evidence.LastCheck != "" {
			lastCheck = &evidence.LastCheck
		}
		freshness = evidence.Freshness
		ownerValidation = evidence.OwnerValidation
		exception = evidence.Exception
		replan = evidence.Replan
		checkWaiver = evidence.CheckWaiver
	}

	patches := []v2FieldPatch{lastCheckV2Patch(lastCheck)}

	freshnessPatch, err := freshnessV2Patch(freshness)
	if err != nil {
		return nil, err
	}
	ownerValidationPatch, err := ownerValidationV2Patch(ownerValidation)
	if err != nil {
		return nil, err
	}
	exceptionPatch, err := exceptionV2Patch(exception)
	if err != nil {
		return nil, err
	}
	replanPatch, err := replanV2Patch(replan)
	if err != nil {
		return nil, err
	}
	checkWaiverPatch, err := checkWaiverV2Patch(checkWaiver)
	if err != nil {
		return nil, err
	}

	return append(patches, freshnessPatch, ownerValidationPatch, exceptionPatch, replanPatch, checkWaiverPatch), nil
}

func checkWaiverV2Patch(waiver *CheckWaiver) (v2FieldPatch, error) {
	if waiver == nil {
		return v2FieldPatch{Key: "check_waiver", Remove: true}, nil
	}
	node, err := encodeV2Node(checkWaiverV2Frontmatter{
		Task:       waiver.Task,
		Reason:     waiver.Reason,
		Actor:      evidenceActorFrontmatter{Role: string(waiver.Actor.Role), Session: waiver.Actor.Session},
		RecordedAt: waiver.RecordedAt.Format(time.RFC3339),
	})
	if err != nil {
		return v2FieldPatch{}, fmt.Errorf("encode check_waiver evidence: %w", err)
	}
	return v2FieldPatch{Key: "check_waiver", Node: node}, nil
}

func lastCheckV2Patch(lastCheck *string) v2FieldPatch {
	if lastCheck == nil {
		return v2FieldPatch{Key: "last_check", Remove: true}
	}
	return v2FieldPatch{Key: "last_check", Value: *lastCheck}
}

func freshnessV2Patch(freshness *Freshness) (v2FieldPatch, error) {
	if freshness == nil {
		return v2FieldPatch{Key: "freshness", Remove: true}, nil
	}
	node, err := encodeV2Node(freshnessV2Frontmatter{
		State:      string(freshness.State),
		Check:      freshness.Check,
		AssessedBy: evidenceActorFrontmatter{Role: string(freshness.AssessedBy.Role), Session: freshness.AssessedBy.Session},
		AssessedAt: freshness.AssessedAt.Format(time.RFC3339),
		Basis:      freshness.Basis,
	})
	if err != nil {
		return v2FieldPatch{}, fmt.Errorf("encode freshness evidence: %w", err)
	}
	return v2FieldPatch{Key: "freshness", Node: node}, nil
}

func ownerValidationV2Patch(ownerValidation *OwnerValidation) (v2FieldPatch, error) {
	if ownerValidation == nil {
		return v2FieldPatch{Key: "owner_validation", Remove: true}, nil
	}
	raw := ownerValidationV2Frontmatter{
		Required:      ownerValidation.Required,
		AcceptedCheck: ownerValidation.AcceptedCheck,
	}
	if ownerValidation.AcceptedBy.Role != "" || ownerValidation.AcceptedBy.Session != "" {
		raw.AcceptedBy = &evidenceActorFrontmatter{
			Role:    string(ownerValidation.AcceptedBy.Role),
			Session: ownerValidation.AcceptedBy.Session,
		}
	}
	node, err := encodeV2Node(raw)
	if err != nil {
		return v2FieldPatch{}, fmt.Errorf("encode owner_validation evidence: %w", err)
	}
	return v2FieldPatch{Key: "owner_validation", Node: node}, nil
}

func exceptionV2Patch(exception *Exception) (v2FieldPatch, error) {
	if exception == nil {
		return v2FieldPatch{Key: "exception", Remove: true}, nil
	}
	node, err := encodeV2Node(exceptionV2Frontmatter{
		Requirements: exception.Requirements,
		Reason:       exception.Reason,
		Owner:        exception.Owner,
		RecordedAt:   exception.RecordedAt.Format(time.RFC3339),
		Check:        exception.Check,
	})
	if err != nil {
		return v2FieldPatch{}, fmt.Errorf("encode exception evidence: %w", err)
	}
	return v2FieldPatch{Key: "exception", Node: node}, nil
}

func replanV2Patch(replan *Replan) (v2FieldPatch, error) {
	if replan == nil {
		return v2FieldPatch{Key: "replan", Remove: true}, nil
	}
	node, err := encodeV2Node(replanV2Frontmatter{
		Reason:     replan.Reason,
		RecordedBy: evidenceActorFrontmatter{Role: string(replan.RecordedBy.Role), Session: replan.RecordedBy.Session},
		RecordedAt: replan.RecordedAt.Format(time.RFC3339),
	})
	if err != nil {
		return v2FieldPatch{}, fmt.Errorf("encode replan evidence: %w", err)
	}
	return v2FieldPatch{Key: "replan", Node: node}, nil
}

// RouterSelectionV2 is the whole of what WriteRouterStateV2 is able to
// change: which Release, Objective, Task, and Issue a V2 router points at. It
// is a type of its own rather than a *RouterStateV2 parameter so the writer's
// signature cannot express a state or next_action write at all. A consumer
// records which record is selected; the phase names which skill owns the
// conversation, and next_action is a retired compatibility field, so neither
// is a selection writer's to touch.
type RouterSelectionV2 struct {
	Release   string // R-### selection, or empty to clear the Release context
	Objective string // O-### selection, or empty to clear the selection
	Task      string // T-### selection, or empty to clear the selection
	Issue     string // I-### selection, or empty to clear the selection
}

// RouterSelectionAfterClosureV2 advances only a selection that names the
// record just closed. Task closure picks the first unfinished Task in the
// owning Objective's already-sorted Task list, regardless of its gate; an
// Objective closure clears that Objective and Task. Release and Issue context
// travel through both transitions unchanged.
func RouterSelectionAfterClosureV2(index *V2Index, current, closed RouterSelectionV2) (RouterSelectionV2, bool) {
	if index == nil {
		return current, false
	}
	if closed.Task != "" {
		task, ok := index.Tasks[closed.Task]
		if !ok || task == nil || task.Status != ColumnDone || current.Task != closed.Task || current.Objective != task.Objective ||
			(closed.Objective != "" && closed.Objective != task.Objective) {
			return current, false
		}

		advanced := current
		advanced.Objective = task.Objective
		advanced.Task = ""
		for _, taskID := range index.ObjectiveTasks[task.Objective] {
			candidate := index.Tasks[taskID]
			if candidate != nil && candidate.Status != ColumnDone {
				advanced.Task = candidate.ID
				break
			}
		}
		return advanced, true
	}
	if closed.Objective != "" && current.Objective == closed.Objective {
		objective, ok := index.Objectives[closed.Objective]
		if ok && objective != nil && objective.Status == ColumnDone {
			cleared := current
			cleared.Objective = ""
			cleared.Task = ""
			return cleared, true
		}
	}
	return current, false
}

// routerSelectionNoneV2 is the router's written spelling of "not selected",
// the sentinel templates/project-v2/.savepoint/router.md ships and
// normalizeRouterSelectionV2 reads back as empty. Clearing a selection
// writes it rather than an empty value, so the field stays legible to a
// human reading the file.
const routerSelectionNoneV2 = "none"

// validate rejects a selection before any file is opened, applying the same
// three rules ReadStateV2 enforces on the way in: O-###/T-### shape, and a
// Task never selected without the Objective that owns it.
func (s RouterSelectionV2) validate() error {
	if s.Release != "" && !matchesV2Identity(s.Release, 'R') {
		return fmt.Errorf("%w: router release %q must be a single R-### selection", ErrV2InvalidID, s.Release)
	}
	if s.Objective != "" && !matchesV2Identity(s.Objective, 'O') {
		return fmt.Errorf("%w: router objective %q must be a single O-### selection", ErrV2InvalidID, s.Objective)
	}
	if s.Task != "" && !matchesV2Identity(s.Task, 'T') {
		return fmt.Errorf("%w: router task %q must be a single T-### selection", ErrV2InvalidID, s.Task)
	}
	if s.Issue != "" && !matchesV2Identity(s.Issue, 'I') {
		return fmt.Errorf("%w: router issue %q must be a single I-### selection", ErrV2InvalidID, s.Issue)
	}
	if s.Task != "" && s.Objective == "" {
		return fmt.Errorf("%w: router task %q selected with no objective", ErrV2InvalidOwnership, s.Task)
	}
	return nil
}

// WriteRouterStateV2 records selection in root's router.md, changing the
// release, objective, task, and issue keys of the "## Current state" anchor and
// nothing else. State, any existing next_action key, every other key the
// anchor carries, and the whole surrounding document survive byte-identical
// (DATA-01), because the anchor is edited as a YAML node tree and only
// selection keys are touched. The writer never adds next_action.
//
// expectedMtime guards the read the caller's selection was based on: a file
// modified since then is refused with ErrMtimeConflict, before the write and
// again immediately before the replacement lands. Writing the selection a
// file already records is a no-op that leaves bytes and modified time alone.
//
// The updated content must decode through ReadStateV2 before it replaces
// anything, so a router this package could not interpret afterwards is
// refused with that reader's own diagnostic and the file is left untouched.
// A router whose recorded state is already invalid is therefore refused too:
// a selection is not worth writing into a document whose meaning cannot be
// established (DATA-03).
func WriteRouterStateV2(root string, selection RouterSelectionV2, expectedMtime time.Time) error {
	if err := selection.validate(); err != nil {
		return err
	}

	path := filepath.Join(root, "router.md")
	info, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("%w: %s is a symlink", ErrV2UnsafePath, path)
	}
	if info.IsDir() {
		return fmt.Errorf("write %s: target is a directory", path)
	}
	if !info.ModTime().Equal(expectedMtime) {
		return ErrMtimeConflict
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	original := string(raw)

	block, err := extractStateBlock(original)
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}

	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(block), &doc); err != nil {
		return fmt.Errorf("%w: %s: router state: %v", ErrV2Malformed, path, err)
	}
	if doc.Kind != yaml.DocumentNode || len(doc.Content) == 0 || doc.Content[0].Kind != yaml.MappingNode {
		return fmt.Errorf("%w: %s: router state is not a mapping", ErrV2Malformed, path)
	}

	if !patchRouterSelectionV2(doc.Content[0], selection) {
		return nil
	}

	out, err := yaml.Marshal(&doc)
	if err != nil {
		return fmt.Errorf("%s: marshal yaml: %w", path, err)
	}

	updated, err := ReplaceStateBlock(original, string(out))
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	// ReplaceStateBlock normalizes line endings to find the anchor; a file
	// that arrived with CRLF keeps it, so "every other byte unchanged" holds
	// on Windows-authored routers too (CFG-02).
	if strings.Contains(original, "\r\n") {
		updated = strings.ReplaceAll(updated, "\n", "\r\n")
	}

	if _, err := NewRouterReader().ReadStateV2(updated); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}

	if info.Mode().Perm()&0222 == 0 {
		return fmt.Errorf("write %s: %w", path, os.ErrPermission)
	}

	return replaceV2File(path, []byte(updated), info.Mode(), func() error {
		latest, err := os.Lstat(path)
		if err != nil {
			return fmt.Errorf("stat %s: %w", path, err)
		}
		if !latest.ModTime().Equal(expectedMtime) {
			return ErrMtimeConflict
		}
		return nil
	})
}

// patchRouterSelectionV2 sets release, objective, task, and issue on mapping
// and reports whether any written value actually changed. It reaches no other
// key, which is how state, next_action, and any field a project added itself
// pass through a selection write untouched. A key the document does not have
// is added only to record a real selection: writing the "not selected"
// sentinel into a router that never carried the field would be a change with
// nothing behind it.
func patchRouterSelectionV2(mapping *yaml.Node, selection RouterSelectionV2) bool {
	changed := false
	for _, field := range []struct {
		key   string
		value string
	}{
		{key: "release", value: routerSelectionValueV2(selection.Release)},
		{key: "objective", value: routerSelectionValueV2(selection.Objective)},
		{key: "task", value: routerSelectionValueV2(selection.Task)},
		{key: "issue", value: routerSelectionValueV2(selection.Issue)},
	} {
		current, present := mappingFieldValue(mapping, field.key)
		if present && current == field.value {
			continue
		}
		if !present && field.value == routerSelectionNoneV2 {
			continue
		}
		setMappingField(mapping, field.key, field.value)
		changed = true
	}
	return changed
}

func routerSelectionValueV2(id string) string {
	if id == "" {
		return routerSelectionNoneV2
	}
	return id
}

// NewCheckV2 is the caller-supplied content for a Check creation: every
// field CreateCheckV2 needs except the ID and Source, which it derives
// itself — the ID from next-unused allocation over index, and Source from
// the file it writes. ExecutedSession identifies the build session being
// checked and is required even when CheckedBy names a different session.
// Body is the exact Markdown appended after the frontmatter delimiter, in
// the same form V2SourceDocument.Body holds it (including its leading
// newline).
type NewCheckV2 struct {
	Scope           CheckScope
	Result          CheckResult
	CheckedBy       Actor
	ExecutedSession string
	CheckedAt       time.Time
	Reviewed        *ReviewedBasis
	Issues          []string
	Supersedes      string
	Body            string
}

// NewIssueV2 is the caller-supplied content for an Issue creation: every
// field CreateIssueV2 needs except the ID and Source, which it derives
// itself — the ID from next-unused allocation over index, and Source from
// the file it writes. A newly created Issue always starts status: open,
// matching the design's "Open -> in_progress when repair starts" flow:
// resolution and history are recorded later through WriteIssueV2 and
// WriteIssueHistoryV2. Body is the exact Markdown appended after the
// frontmatter delimiter, in the same form V2SourceDocument.Body holds it
// (including its leading newline).
type NewIssueV2 struct {
	Title        string
	Type         IssueType
	Origin       IssueOrigin
	Tasks        []string
	Checks       []string
	GuardrailIDs []string
	Severity     string
	Body         string
}
