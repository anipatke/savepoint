package data

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

var ErrMtimeConflict = fmt.Errorf("file modified since last read")
var ErrProposalNotFound = fmt.Errorf("target text not found in file")

// ApplyProposal replaces the first occurrence of old with newText in the file at path.
func ApplyProposal(path, old, newText string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	normalized := normalizeLineEndings(string(content))
	if !strings.Contains(normalized, old) {
		return fmt.Errorf("%w: %s", ErrProposalNotFound, path)
	}
	updated := strings.Replace(normalized, old, newText, 1)
	return os.WriteFile(path, []byte(updated), 0644)
}

// UpdateEpicStatus sets the status field in the frontmatter of an E##-Detail.md file.
func UpdateEpicStatus(path, status string) error {
	return updateFrontmatterField(path, "status", status)
}

// UpdateLastAudited sets the last_audited field in the frontmatter of Design.md.
func UpdateLastAudited(path, value string) error {
	return updateFrontmatterField(path, "last_audited", value)
}

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

func updateFrontmatterField(path, key, value string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}

	normalized := normalizeLineEndings(string(content))

	raw, body, err := SplitFrontmatterBody(normalized)
	if err != nil {
		return fmt.Errorf("extract frontmatter: %w", err)
	}

	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(raw), &doc); err != nil {
		return fmt.Errorf("parse yaml: %w", err)
	}

	if doc.Kind != yaml.DocumentNode || len(doc.Content) == 0 {
		return fmt.Errorf("unexpected yaml structure")
	}

	mapping := doc.Content[0]
	if mapping.Kind != yaml.MappingNode {
		return fmt.Errorf("frontmatter is not a mapping")
	}

	setMappingField(mapping, key, value)

	out, err := yaml.Marshal(&doc)
	if err != nil {
		return fmt.Errorf("marshal yaml: %w", err)
	}

	newContent := "---\n" + strings.TrimSpace(string(out)) + "\n---" + body
	return os.WriteFile(path, []byte(newContent), 0644)
}

func WriteTaskStatus(path string, task *Task, expectedMtime time.Time) error {
	HealTaskMetadataForProgress(task)
	if err := ValidateTaskLifecycle(task); err != nil {
		return err
	}

	fi, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}

	if !fi.ModTime().Equal(expectedMtime) {
		return ErrMtimeConflict
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}

	normalized := normalizeLineEndings(string(content))

	raw, body, err := SplitFrontmatterBody(normalized)
	if err != nil {
		return fmt.Errorf("extract frontmatter: %w", err)
	}

	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(raw), &doc); err != nil {
		return fmt.Errorf("parse yaml: %w", err)
	}

	if doc.Kind != yaml.DocumentNode || len(doc.Content) == 0 {
		return fmt.Errorf("unexpected yaml structure")
	}

	mapping := doc.Content[0]
	if mapping.Kind != yaml.MappingNode {
		return fmt.Errorf("frontmatter is not a mapping")
	}

	setMappingField(mapping, "status", string(task.Column))

	if task.Stage == "" {
		removeMappingField(mapping, "phase")
		removeMappingField(mapping, "stage")
	} else {
		setMappingField(mapping, "stage", string(task.Stage))
		removeMappingField(mapping, "phase")
	}

	if task.ComplexityTier != "" {
		setMappingField(mapping, "complexity_tier", string(task.ComplexityTier))
	}
	if task.ComplexityReason != "" {
		setMappingField(mapping, "complexity_reason", task.ComplexityReason)
	}

	out, err := yaml.Marshal(&doc)
	if err != nil {
		return fmt.Errorf("marshal yaml: %w", err)
	}

	newContent := "---\n" + strings.TrimSpace(string(out)) + "\n---" + body

	return os.WriteFile(path, []byte(newContent), 0644)
}

func WriteDefectStatus(path string, defect *Defect, expectedMtime time.Time) error {
	if err := validateDefectLifecycle(defect); err != nil {
		return err
	}

	fi, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}

	if !fi.ModTime().Equal(expectedMtime) {
		return ErrMtimeConflict
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}

	normalized := normalizeLineEndings(string(content))

	raw, body, err := SplitFrontmatterBody(normalized)
	if err != nil {
		return fmt.Errorf("extract frontmatter: %w", err)
	}

	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(raw), &doc); err != nil {
		return fmt.Errorf("parse yaml: %w", err)
	}

	if doc.Kind != yaml.DocumentNode || len(doc.Content) == 0 {
		return fmt.Errorf("unexpected yaml structure")
	}

	mapping := doc.Content[0]
	if mapping.Kind != yaml.MappingNode {
		return fmt.Errorf("frontmatter is not a mapping")
	}

	setMappingField(mapping, "status", string(defect.Status))

	if defect.Stage == "" {
		removeMappingField(mapping, "stage")
	} else {
		setMappingField(mapping, "stage", string(defect.Stage))
	}

	out, err := yaml.Marshal(&doc)
	if err != nil {
		return fmt.Errorf("marshal yaml: %w", err)
	}

	newContent := "---\n" + strings.TrimSpace(string(out)) + "\n---" + body

	return os.WriteFile(path, []byte(newContent), 0644)
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
	if filepath.IsAbs(source.Path) || source.ProjectRoot == "" {
		return filepath.Clean(source.Path), nil
	}

	root, err := filepath.Abs(source.ProjectRoot)
	if err != nil {
		return "", fmt.Errorf("resolve V2 project root %q: %w", source.ProjectRoot, err)
	}
	root = filepath.Clean(root)
	path := filepath.Clean(filepath.Join(root, source.Path))
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return "", fmt.Errorf("resolve V2 source %s: %w", source.Path, err)
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
// freshness, owner_validation, exception, and replan — of a V2 Task
// record's frontmatter to match task.Evidence, preserving every other YAML
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
	if evidence != nil {
		if evidence.LastCheck != "" {
			lastCheck = &evidence.LastCheck
		}
		freshness = evidence.Freshness
		ownerValidation = evidence.OwnerValidation
		exception = evidence.Exception
		replan = evidence.Replan
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

	return append(patches, freshnessPatch, ownerValidationPatch, exceptionPatch, replanPatch), nil
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
// change: which Release, Objective, and Task a V2 router points at. It is a type of
// its own rather than a *RouterStateV2 parameter so that the writer's
// signature cannot express a state or next_action write at all. A consumer
// records which record is selected; the phase names which skill owns the
// conversation, and next_action is prose those skills author, so neither is
// a selection writer's to touch.
type RouterSelectionV2 struct {
	Release   string // R### selection, or empty to clear the Release context
	Objective string // O### selection, or empty to clear the selection
	Task      string // T### selection, or empty to clear the selection
}

// routerSelectionNoneV2 is the router's written spelling of "not selected",
// the sentinel templates/project-v2/.savepoint/router.md ships and
// normalizeRouterSelectionV2 reads back as empty. Clearing a selection
// writes it rather than an empty value, so the field stays legible to a
// human reading the file.
const routerSelectionNoneV2 = "none"

// validate rejects a selection before any file is opened, applying the same
// three rules ReadStateV2 enforces on the way in: O###/T### shape, and a
// Task never selected without the Objective that owns it.
func (s RouterSelectionV2) validate() error {
	if s.Release != "" && !releaseIDPatternV2.MatchString(s.Release) {
		return fmt.Errorf("%w: router release %q must be a single R### selection", ErrV2InvalidID, s.Release)
	}
	if s.Objective != "" && !objectiveIDPattern.MatchString(s.Objective) {
		return fmt.Errorf("%w: router objective %q must be a single O### selection", ErrV2InvalidID, s.Objective)
	}
	if s.Task != "" && !taskIDPatternV2.MatchString(s.Task) {
		return fmt.Errorf("%w: router task %q must be a single T### selection", ErrV2InvalidID, s.Task)
	}
	if s.Task != "" && s.Objective == "" {
		return fmt.Errorf("%w: router task %q selected with no objective", ErrV2InvalidOwnership, s.Task)
	}
	return nil
}

// WriteRouterStateV2 records selection in root's router.md, changing the
// release, objective, and task keys of the "## Current state" anchor and
// nothing else.
// state, next_action, every other key the anchor carries, and the whole
// surrounding document survive byte-identical (DATA-01), because the anchor
// is edited as a YAML node tree and only those two keys are touched.
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

// patchRouterSelectionV2 sets release, objective, and task on mapping and reports
// whether either key's written value actually changed. It reaches no other
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

func WriteRouterState(root string, state *RouterState, expectedMtime time.Time) error {
	path := filepath.Join(root, "router.md")
	fi, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}

	if !fi.ModTime().Equal(expectedMtime) {
		return ErrMtimeConflict
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}

	normalized := normalizeLineEndings(string(content))

	startIdx := strings.Index(strings.ToLower(normalized), strings.ToLower(stateBlockStart))
	if startIdx == -1 {
		return fmt.Errorf("no Current state block found")
	}

	yamlStart := strings.Index(normalized[startIdx:], "```yaml")
	if yamlStart == -1 {
		return fmt.Errorf("no yaml code block found")
	}

	yamlStart += startIdx + len("```yaml")
	yamlEnd := strings.Index(normalized[yamlStart:], "```")
	if yamlEnd == -1 {
		return fmt.Errorf("no closing code block found")
	}

	out, err := yaml.Marshal(state)
	if err != nil {
		return fmt.Errorf("marshal yaml: %w", err)
	}

	newContent := normalized[:yamlStart] + "\n" + strings.TrimSpace(string(out)) + "\n" + normalized[yamlStart+yamlEnd:]

	return os.WriteFile(path, []byte(newContent), 0644)
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

// CreateCheckV2 allocates the next unused Check ID over index, marshals a
// new Check record from fields, and writes it create-only to
// root/checks/{id}.md. Checks are immutable once written: CreateCheckV2
// never rewrites or deletes an existing file, and refuses an ID collision
// with ErrV2CheckImmutable rather than overwriting it. The marshalled
// content is validated as a decodable CheckV2 before the file is created, so
// a rejected record leaves no file behind.
func CreateCheckV2(root string, index *V2Index, fields NewCheckV2) (*CheckV2, error) {
	if fields.Scope.Kind == CheckScopeRelease && !checkScopeTargetExists(index, fields.Scope) {
		return nil, fmt.Errorf("%w: check scope names missing release %s", ErrV2CheckMissingScopeTarget, fields.Scope.ID)
	}
	id := nextV2CheckID(index)

	raw := checkV2Frontmatter{
		ID:              id,
		Scope:           checkScopeFrontmatter{Kind: string(fields.Scope.Kind), ID: fields.Scope.ID},
		Result:          string(fields.Result),
		CheckedBy:       checkActorFrontmatter{Role: string(fields.CheckedBy.Role), Session: fields.CheckedBy.Session},
		ExecutedSession: fields.ExecutedSession,
		CheckedAt:       fields.CheckedAt.Format(time.RFC3339),
		Issues:          fields.Issues,
		Supersedes:      fields.Supersedes,
	}
	if fields.Reviewed != nil {
		raw.Reviewed = &reviewedFrontmatter{
			BaseCommit:   fields.Reviewed.BaseCommit,
			HeadCommit:   fields.Reviewed.HeadCommit,
			Files:        fields.Reviewed.Files,
			Dependencies: fields.Reviewed.Dependencies,
		}
	}

	out, err := yaml.Marshal(&raw)
	if err != nil {
		return nil, fmt.Errorf("create check %s: marshal yaml: %w", id, err)
	}

	relPath := filepath.Join(v2ChecksDirName, id+".md")
	content := "---\n" + strings.TrimSpace(string(out)) + "\n---" + fields.Body

	check, err := DecodeCheckV2(relPath, content)
	if err != nil {
		return nil, fmt.Errorf("create check %s: %w", id, err)
	}

	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("create check %s: resolve project root %q: %w", id, root, err)
	}
	path := filepath.Join(rootAbs, relPath)

	if err := createV2CheckFile(path, []byte(content)); err != nil {
		return nil, err
	}

	check.Source.ProjectRoot = rootAbs
	return check, nil
}

// nextV2CheckID allocates the next unused Check ID from index.Checks: one
// past the highest numeric ID already present, formatted with at least
// three digits. It never reuses an ID already present. Migration
// reservations extend this allocation rule in E45; E43 allocates over active
// records only.
func nextV2CheckID(index *V2Index) string {
	next := 1
	for id := range index.Checks {
		n, err := strconv.Atoi(strings.TrimPrefix(id, "C"))
		if err != nil {
			continue
		}
		if n >= next {
			next = n + 1
		}
	}
	return fmt.Sprintf("C%03d", next)
}

// v2CheckTempFile is the small file boundary needed by createV2CheckFile.
// Keeping the boundary narrower than *os.File makes every create-only
// operation failure deterministic in tests without changing the production
// write sequence.
type v2CheckTempFile interface {
	Name() string
	Chmod(os.FileMode) error
	Write([]byte) (int, error)
	Sync() error
	Close() error
}

type v2CheckFileOperations struct {
	mkdirAll   func(string, os.FileMode) error
	createTemp func(string, string) (v2CheckTempFile, error)
	remove     func(string) error
	link       func(string, string) error
}

var v2CheckFileOps = v2CheckFileOperations{
	mkdirAll: os.MkdirAll,
	createTemp: func(dir, pattern string) (v2CheckTempFile, error) {
		return os.CreateTemp(dir, pattern)
	},
	remove: os.Remove,
	link:   os.Link,
}

// createV2CheckFile writes a new Check record create-only. It is a thin call
// over createV2RecordFile so a collision is reported as ErrV2CheckImmutable,
// matching Checks' immutable-once-written contract.
func createV2CheckFile(path string, content []byte) error {
	return createV2RecordFile(path, content, ErrV2CheckImmutable)
}

// createV2RecordFile writes content to path create-only: it fully writes and
// syncs a same-directory temporary file, then hard-links it into place so an
// existing file at path is never rewritten, truncated, or replaced —
// os.Rename would silently overwrite, which a create-only write cannot
// allow. A collision at path is reported wrapping existsErr, the calling
// record family's own diagnostic, so a Check collision and an Issue
// collision are told apart. The operation boundary is injectable only for
// deterministic tests; production defaults remain the standard filesystem
// calls above.
func createV2RecordFile(path string, content []byte, existsErr error) (retErr error) {
	ops := v2CheckFileOps
	dir := filepath.Dir(path)
	if err := ops.mkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}

	temp, err := ops.createTemp(dir, ".savepoint-v2-check-*")
	if err != nil {
		return fmt.Errorf("create %s: create temporary file: %w", path, err)
	}
	tempPath := temp.Name()
	defer func() {
		if removeErr := ops.remove(tempPath); removeErr != nil && !os.IsNotExist(removeErr) {
			if retErr != nil {
				retErr = fmt.Errorf("%w (cleanup %s: %v)", retErr, tempPath, removeErr)
			} else {
				retErr = fmt.Errorf("cleanup %s: %w", tempPath, removeErr)
			}
		}
	}()

	if err := temp.Chmod(0644); err != nil {
		closeErr := temp.Close()
		if closeErr != nil {
			return fmt.Errorf("create %s: set temporary permissions: %w (close: %v)", path, err, closeErr)
		}
		return fmt.Errorf("create %s: set temporary permissions: %w", path, err)
	}
	written, writeErr := temp.Write(content)
	if writeErr != nil || written != len(content) {
		closeErr := temp.Close()
		if writeErr == nil {
			writeErr = fmt.Errorf("short write: wrote %d of %d bytes", written, len(content))
		}
		if closeErr != nil {
			return fmt.Errorf("create %s: %w (close: %v)", path, writeErr, closeErr)
		}
		return fmt.Errorf("create %s: %w", path, writeErr)
	}
	if err := temp.Sync(); err != nil {
		closeErr := temp.Close()
		if closeErr != nil {
			return fmt.Errorf("create %s: flush temporary file: %w (close: %v)", path, err, closeErr)
		}
		return fmt.Errorf("create %s: flush temporary file: %w", path, err)
	}
	if err := temp.Close(); err != nil {
		return fmt.Errorf("create %s: close temporary file: %w", path, err)
	}

	if err := ops.link(tempPath, path); err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("%w: %s", existsErr, path)
		}
		return fmt.Errorf("create %s: %w", path, err)
	}
	return nil
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

// CreateIssueV2 allocates the next unused Issue ID over index, marshals a new
// Issue record from fields with status: open, and writes it create-only to
// root/issues/{id}-{slug}.md, matching the {ID}-{slug}.md convention Tasks
// and Epics already use. CreateIssueV2 never rewrites or deletes an existing
// file, and refuses an ID collision with ErrV2IssueAlreadyExists rather than
// overwriting it. The marshalled content is validated as a decodable
// IssueV2 before the file is created, so a rejected record leaves no file
// behind.
func CreateIssueV2(root string, index *V2Index, fields NewIssueV2) (*IssueV2, error) {
	id := nextV2IssueID(index)

	raw := issueV2Frontmatter{
		ID:     id,
		Title:  fields.Title,
		Type:   string(fields.Type),
		Status: string(IssueStatusOpen),
		Source: &issueOriginFrontmatter{
			Kind:  string(fields.Origin.Kind),
			Check: fields.Origin.Check,
			Actor: evidenceActorFrontmatter{Role: string(fields.Origin.Actor.Role), Session: fields.Origin.Actor.Session},
			At:    fields.Origin.At.Format(time.RFC3339),
		},
		Tasks:        fields.Tasks,
		Checks:       fields.Checks,
		GuardrailIDs: fields.GuardrailIDs,
		Severity:     fields.Severity,
	}

	out, err := yaml.Marshal(&raw)
	if err != nil {
		return nil, fmt.Errorf("create issue %s: marshal yaml: %w", id, err)
	}

	filename := id
	if slug := issueV2Slug(fields.Title); slug != "" {
		filename = id + "-" + slug
	}
	relPath := filepath.Join(v2IssuesDirName, filename+".md")
	content := "---\n" + strings.TrimSpace(string(out)) + "\n---" + fields.Body

	issue, err := DecodeIssueV2(relPath, content)
	if err != nil {
		return nil, fmt.Errorf("create issue %s: %w", id, err)
	}

	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("create issue %s: resolve project root %q: %w", id, root, err)
	}
	path := filepath.Join(rootAbs, relPath)

	if err := createV2RecordFile(path, []byte(content), ErrV2IssueAlreadyExists); err != nil {
		return nil, err
	}

	issue.Source.ProjectRoot = rootAbs
	return issue, nil
}

// issueV2Slug converts title into the lowercase, hyphen-separated slug
// CreateIssueV2 appends to the allocated ID, matching the {ID}-{slug}.md
// filename convention Tasks and Epics already use. Discovery only requires a
// filename to start with its declared ID, so an empty slug (an all-symbol
// title) is not an error — CreateIssueV2 falls back to the bare ID.
func issueV2Slug(title string) string {
	var b strings.Builder
	lastHyphen := true
	for _, r := range strings.ToLower(title) {
		switch {
		case r >= 'a' && r <= 'z' || r >= '0' && r <= '9':
			b.WriteRune(r)
			lastHyphen = false
		default:
			if !lastHyphen {
				b.WriteByte('-')
				lastHyphen = true
			}
		}
	}
	return strings.TrimSuffix(b.String(), "-")
}

// nextV2IssueID allocates the next unused Issue ID from index.Issues: one
// past the highest numeric ID already present, formatted with at least three
// digits. It never reuses an ID already present, and allocates over active
// records only, matching nextV2CheckID.
func nextV2IssueID(index *V2Index) string {
	next := 1
	for id := range index.Issues {
		n, err := strconv.Atoi(strings.TrimPrefix(id, "I"))
		if err != nil {
			continue
		}
		if n >= next {
			next = n + 1
		}
	}
	return fmt.Sprintf("I%03d", next)
}

// WriteIssueV2 patches only the status, resolution, duplicate_of, tasks,
// checks, and severity fields of a V2 Issue record's frontmatter to match
// issue's in-memory fields, preserving every other YAML key, unknown field,
// and the authored Markdown body unchanged. History is never touched here —
// it has its own append-only path in WriteIssueHistoryV2. The patched
// content is validated as a decodable IssueV2 before any file is replaced,
// so a malformed value is refused and the file is left untouched, and
// writing values the record already has is a no-op.
func WriteIssueV2(issue *IssueV2) error {
	patches, err := issueManagedPatches(issue)
	if err != nil {
		return err
	}

	return writeV2Record(&issue.Source, patches, func(content string) error {
		_, err := DecodeIssueV2(issue.Source.Path, content)
		return err
	})
}

func issueManagedPatches(issue *IssueV2) ([]v2FieldPatch, error) {
	tasksPatch, err := issueListV2Patch("tasks", issue.Tasks)
	if err != nil {
		return nil, err
	}
	checksPatch, err := issueListV2Patch("checks", issue.Checks)
	if err != nil {
		return nil, err
	}
	resolutionPatch, err := issueResolutionV2Patch(issue.Resolution)
	if err != nil {
		return nil, err
	}

	return []v2FieldPatch{
		{Key: "status", Value: string(issue.Status)},
		resolutionPatch,
		issueScalarV2Patch("duplicate_of", issue.DuplicateOf),
		tasksPatch,
		checksPatch,
		issueScalarV2Patch("severity", issue.Severity),
	}, nil
}

// issueListV2Patch patches an optional identity-reference list field,
// removing the key entirely when values is empty rather than writing an
// empty sequence — matching CreateIssueV2's omitempty output for the same
// field, so writing a record its own empty list back is a no-op.
func issueListV2Patch(key string, values []string) (v2FieldPatch, error) {
	if len(values) == 0 {
		return v2FieldPatch{Key: key, Remove: true}, nil
	}
	node, err := encodeV2Node(values)
	if err != nil {
		return v2FieldPatch{}, fmt.Errorf("encode issue %s: %w", key, err)
	}
	return v2FieldPatch{Key: key, Node: node}, nil
}

// issueScalarV2Patch patches an optional scalar Issue field, removing the
// key entirely when value is empty rather than writing it blank — matching
// the rule that an absent optional field stays absent.
func issueScalarV2Patch(key, value string) v2FieldPatch {
	if value == "" {
		return v2FieldPatch{Key: key, Remove: true}
	}
	return v2FieldPatch{Key: key, Value: value}
}

func issueResolutionV2Patch(resolution *IssueResolution) (v2FieldPatch, error) {
	if resolution == nil {
		return v2FieldPatch{Key: "resolution", Remove: true}, nil
	}
	node, err := encodeV2Node(issueResolutionFrontmatter{
		Disposition: string(resolution.Disposition),
		Check:       resolution.Check,
		Actor:       evidenceActorFrontmatter{Role: string(resolution.Actor.Role), Session: resolution.Actor.Session},
		At:          resolution.At.Format(time.RFC3339),
		Reason:      resolution.Reason,
	})
	if err != nil {
		return v2FieldPatch{}, fmt.Errorf("encode issue resolution: %w", err)
	}
	return v2FieldPatch{Key: "resolution", Node: node}, nil
}

// WriteIssueHistoryV2 replaces issue's recorded history with entries, which
// must extend it: entries shorter than the recorded history, or differing in
// any already-recorded entry, are refused with ErrV2IssueHistoryNotAppendOnly
// and leave the file untouched. Every earlier entry's fields are preserved
// byte-for-byte in the rewritten frontmatter because they are the same
// values re-marshalled, never edited. The patched content is validated as a
// decodable IssueV2 before any file is replaced, and appending the same
// entries the record already has is a no-op.
func WriteIssueHistoryV2(issue *IssueV2, entries []IssueHistoryEntry) error {
	if err := validateAppendOnlyIssueHistory(issue.Source.Path, issue.ID, issue.History, entries); err != nil {
		return err
	}

	// An empty history removes the key rather than recording an empty
	// sequence, so writing no entries to a record that has none stays a true
	// no-op instead of rewriting the file with `history: []`.
	patch := v2FieldPatch{Key: "history", Remove: true}
	if len(entries) > 0 {
		node, err := encodeV2Node(issueHistoryV2Frontmatter(entries))
		if err != nil {
			return fmt.Errorf("encode issue history: %w", err)
		}
		patch = v2FieldPatch{Key: "history", Node: node}
	}

	return writeV2Record(&issue.Source, []v2FieldPatch{patch}, func(content string) error {
		_, err := DecodeIssueV2(issue.Source.Path, content)
		return err
	})
}

// validateAppendOnlyIssueHistory refuses any write whose entries do not
// start with recorded's exact entries in the same order: a shorter list, a
// reordered entry, or an edited entry are each a named diagnostic rather
// than a silently accepted rewrite of the past.
func validateAppendOnlyIssueHistory(path, id string, recorded, entries []IssueHistoryEntry) error {
	if len(entries) < len(recorded) {
		return fmt.Errorf("%w: %s: issue %s history has %d recorded entries, write has only %d", ErrV2IssueHistoryNotAppendOnly, path, id, len(recorded), len(entries))
	}
	for i, existing := range recorded {
		if !issueHistoryEntryEqual(existing, entries[i]) {
			return fmt.Errorf("%w: %s: issue %s history entry %d differs from the recorded entry", ErrV2IssueHistoryNotAppendOnly, path, id, i)
		}
	}
	return nil
}

func issueHistoryEntryEqual(a, b IssueHistoryEntry) bool {
	return a.At.Equal(b.At) && a.Actor == b.Actor && a.Kind == b.Kind && a.Note == b.Note && a.Check == b.Check
}

// issueHistoryV2Frontmatter re-marshals entries in recorded order into the
// frontmatter shape decodeIssueHistory reads back.
func issueHistoryV2Frontmatter(entries []IssueHistoryEntry) []issueHistoryFrontmatter {
	out := make([]issueHistoryFrontmatter, 0, len(entries))
	for _, entry := range entries {
		out = append(out, issueHistoryFrontmatter{
			At:    entry.At.Format(time.RFC3339),
			Actor: evidenceActorFrontmatter{Role: string(entry.Actor.Role), Session: entry.Actor.Session},
			Kind:  string(entry.Kind),
			Note:  entry.Note,
			Check: entry.Check,
		})
	}
	return out
}
