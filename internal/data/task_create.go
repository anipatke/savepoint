package data

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"
)

// TaskCreationResult identifies the new Task and its path relative to the
// project's .savepoint directory.
type TaskCreationResult struct {
	ID   string
	Path string
}

type idFreeTaskDraft struct {
	content      string
	title        string
	addObjective bool
}

// CreateTaskV2 creates one Task for objectiveID from an ID-free Markdown
// draft. The draft must contain V2 Task frontmatter without an id; objective
// may be omitted, or may match objectiveID. All existing draft bytes except
// the allocator supplied identity lines and the new filename are preserved.
// The Task ID lock remains held through exclusive creation and strict index
// validation. A post-reservation failure removes only the new Task file and
// leaves the ID retired.
func CreateTaskV2(root, objectiveID, draftContent string) (TaskCreationResult, error) {
	if !matchesV2Identity(objectiveID, 'O') {
		return TaskCreationResult{}, fmt.Errorf("create Task: objective %q must match O- plus at least three digits", objectiveID)
	}
	draft, err := parseIDFreeTaskDraft(draftContent, objectiveID)
	if err != nil {
		return TaskCreationResult{}, err
	}

	var result TaskCreationResult
	_, err = withTaskIDReservation(root, func(index *V2Index) error {
		if _, ok := index.Objectives[objectiveID]; !ok {
			return fmt.Errorf("create Task: objective %s does not exist in this project", objectiveID)
		}
		return nil
	}, func(taskID string, index *V2Index) (func() error, error) {
		objective := index.Objectives[objectiveID]
		objectivePath := filepath.Join(root, objective.Source.Path)
		objectiveDir := filepath.Dir(objectivePath)
		tasksDir := filepath.Join(objectiveDir, v2TasksDirName)
		createdTasksDir, err := ensureTaskDirectory(tasksDir)
		if err != nil {
			return nil, fmt.Errorf("create Task: prepare Objective task directory %s: %w", tasksDir, err)
		}

		filename := taskID + "-" + taskSlug(draft.title) + ".md"
		taskPath := filepath.Join(tasksDir, filename)
		relativePath := filepath.Join(filepath.Dir(objective.Source.Path), v2TasksDirName, filename)
		content := draft.withIdentity(taskID, objectiveID)
		if _, err := DecodeTaskV2(relativePath, content); err != nil {
			return rollbackTaskCreation(tasksDir, createdTasksDir, taskPath, nil), fmt.Errorf("create Task: validate draft for %s: %w", taskID, err)
		}

		info, err := createExclusiveTaskFile(taskPath, []byte(content))
		if err != nil {
			return rollbackTaskCreation(tasksDir, createdTasksDir, taskPath, nil), fmt.Errorf("create Task: write %s: %w", relativePath, err)
		}
		rollback := rollbackTaskCreation(tasksDir, createdTasksDir, taskPath, info)

		validated, err := LoadV2Index(root)
		if err != nil {
			return rollback, fmt.Errorf("create Task: strict-load project after creating %s: %w", relativePath, err)
		}
		if validated.Tasks[taskID] == nil {
			return rollback, fmt.Errorf("create Task: strict-load project did not index new Task %s at %s", taskID, relativePath)
		}

		result = TaskCreationResult{ID: taskID, Path: relativePath}
		return rollback, nil
	})
	if err != nil {
		return TaskCreationResult{}, err
	}
	return result, nil
}

func parseIDFreeTaskDraft(content, objectiveID string) (*idFreeTaskDraft, error) {
	yamlContent, _, err := SplitFrontmatterBody(content)
	if err != nil {
		return nil, fmt.Errorf("create Task: parse ID-free draft frontmatter: %w", err)
	}

	var document yaml.Node
	if err := yaml.Unmarshal([]byte(yamlContent), &document); err != nil {
		return nil, fmt.Errorf("create Task: malformed ID-free draft YAML: %w", err)
	}
	if document.Kind != yaml.DocumentNode || len(document.Content) != 1 || document.Content[0].Kind != yaml.MappingNode {
		return nil, fmt.Errorf("create Task: ID-free draft frontmatter must be a mapping")
	}

	mapping := document.Content[0]
	seen := make(map[string]bool, len(mapping.Content)/2)
	hasObjective := false
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		key := mapping.Content[i]
		value := mapping.Content[i+1]
		if key.Kind != yaml.ScalarNode || key.Tag != "!!str" {
			return nil, fmt.Errorf("create Task: ID-free draft frontmatter keys must be strings")
		}
		if seen[key.Value] {
			return nil, fmt.Errorf("create Task: ID-free draft has duplicate frontmatter field %q", key.Value)
		}
		seen[key.Value] = true
		switch key.Value {
		case "id":
			return nil, fmt.Errorf("create Task: draft contains reserved id field; the allocator assigns Task IDs")
		case "objective":
			if value.Kind != yaml.ScalarNode || value.Tag != "!!str" || value.Value != objectiveID {
				return nil, fmt.Errorf("create Task: draft objective must match supplied Objective %s", objectiveID)
			}
			hasObjective = true
		}
	}

	draft := &idFreeTaskDraft{content: content, addObjective: !hasObjective}
	validated, err := DecodeTaskV2("task-draft.md", draft.withIdentity("T-000", objectiveID))
	if err != nil {
		return nil, fmt.Errorf("create Task: invalid ID-free draft: %w", err)
	}
	draft.title = validated.Title
	return draft, nil
}

func (draft *idFreeTaskDraft) withIdentity(taskID, objectiveID string) string {
	lineEnding := "\n"
	openingLength := len("---\n")
	if strings.HasPrefix(draft.content, "---\r\n") {
		lineEnding = "\r\n"
		openingLength = len("---\r\n")
	}
	identity := "id: " + taskID + lineEnding
	if draft.addObjective {
		identity += "objective: " + objectiveID + lineEnding
	}
	return draft.content[:openingLength] + identity + draft.content[openingLength:]
}

func ensureTaskDirectory(path string) (bool, error) {
	info, err := os.Lstat(path)
	if err == nil {
		if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
			return false, fmt.Errorf("path is not a real directory")
		}
		return false, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return false, err
	}
	if err := os.Mkdir(path, 0o755); err == nil {
		return true, nil
	} else if !errors.Is(err, os.ErrExist) {
		return false, err
	}

	info, err = os.Lstat(path)
	if err != nil {
		return false, err
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
		return false, fmt.Errorf("path is not a real directory")
	}
	return false, nil
}

func createExclusiveTaskFile(path string, content []byte) (os.FileInfo, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	info, statErr := file.Stat()
	written, writeErr := file.Write(content)
	if writeErr == nil && written != len(content) {
		writeErr = io.ErrShortWrite
	}
	var syncErr error
	if writeErr == nil {
		syncErr = file.Sync()
	}
	closeErr := file.Close()
	if err := errors.Join(statErr, writeErr, syncErr, closeErr); err != nil {
		if info != nil {
			if cleanupErr := removeCreatedTaskFile(path, info); cleanupErr != nil {
				err = errors.Join(err, fmt.Errorf("cleanup new Task file: %w", cleanupErr))
			}
		} else if cleanupErr := os.Remove(path); cleanupErr != nil && !errors.Is(cleanupErr, os.ErrNotExist) {
			err = errors.Join(err, fmt.Errorf("cleanup new Task file: %w", cleanupErr))
		}
		return nil, err
	}
	return info, nil
}

func rollbackTaskCreation(tasksDir string, createdTasksDir bool, taskPath string, info os.FileInfo) func() error {
	return func() error {
		var rollbackErr error
		if info != nil {
			rollbackErr = removeCreatedTaskFile(taskPath, info)
		}
		if createdTasksDir {
			if err := os.Remove(tasksDir); err != nil && !errors.Is(err, os.ErrNotExist) {
				rollbackErr = errors.Join(rollbackErr, fmt.Errorf("remove newly created task directory %s: %w", tasksDir, err))
			}
		}
		return rollbackErr
	}
}

func removeCreatedTaskFile(path string, created os.FileInfo) error {
	current, err := os.Lstat(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("inspect %s before rollback: %w", path, err)
	}
	if current.Mode()&os.ModeSymlink != 0 || !os.SameFile(created, current) {
		return fmt.Errorf("refusing to remove changed path %s during rollback", path)
	}
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("remove %s during rollback: %w", path, err)
	}
	return nil
}

func taskSlug(title string) string {
	var slug strings.Builder
	separatorPending := false
	for _, char := range strings.ToLower(strings.TrimSpace(title)) {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			if separatorPending && slug.Len() > 0 {
				slug.WriteByte('-')
			}
			separatorPending = false
			slug.WriteRune(char)
		} else if slug.Len() > 0 {
			separatorPending = true
		}
	}
	if slug.Len() == 0 {
		return "task"
	}
	return slug.String()
}
