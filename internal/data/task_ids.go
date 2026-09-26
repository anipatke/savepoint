package data

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	taskIDHighWaterFile = "task-ids.yml"
	taskIDLockFile      = "task-ids.lock"
	taskIDLockRetry     = 10 * time.Millisecond
)

// taskIDLockWait is how long an allocator waits for task-ids.lock before
// calling it leftover. A healthy hold takes a few milliseconds on Linux and
// tens on Windows, so even a long queue of live allocators finishes well
// inside it (I-061). It is a variable only so tests can shorten it.
var taskIDLockWait = 10 * time.Second

type taskIDHighWater struct {
	lastIssued string
	content    []byte
	info       os.FileInfo
	exists     bool
}

// AllocateTaskID reserves the next project-wide Task identity. root is the
// project's .savepoint directory. The reservation is durable before the ID is
// returned, so a caller that later fails to create its Task leaves a gap
// rather than making the identity available again.
func AllocateTaskID(root string) (taskID string, retErr error) {
	return withTaskIDReservation(root, nil, nil)
}

// withTaskIDReservation validates the project before touching allocation
// state, acquires the project lock, reserves one identity, and holds the lock
// until operation succeeds or its rollback has run. validate runs both before
// lock acquisition and again under the lock.
func withTaskIDReservation(
	root string,
	validate func(*V2Index) error,
	operation func(string, *V2Index) (func() error, error),
) (taskID string, retErr error) {
	preflight, err := LoadV2Index(root)
	if err != nil {
		return "", fmt.Errorf("allocate Task ID: strict-load project %s: %w", root, err)
	}
	if validate != nil {
		if err := validate(preflight); err != nil {
			return "", err
		}
	}

	lockPath := filepath.Join(root, taskIDLockFile)
	lock, err := acquireTaskIDLock(lockPath)
	if err != nil {
		return "", err
	}
	locked := true
	var rollback func() error
	defer func() {
		if !locked {
			return
		}
		if err := releaseTaskIDLock(lock, lockPath); err != nil {
			taskID = ""
			lockErr := fmt.Errorf("allocate Task ID: release lock %s: %w", lockPath, err)
			if rollback != nil {
				if cleanupErr := rollback(); cleanupErr != nil {
					lockErr = errors.Join(lockErr, fmt.Errorf("rollback created Task: %w", cleanupErr))
				}
			}
			if retErr == nil {
				retErr = lockErr
			} else {
				retErr = errors.Join(retErr, lockErr)
			}
		}
	}()

	// Re-load while holding the project lock so another allocator's work, or a
	// project edit made while this caller waited, cannot invalidate the
	// high-water calculation.
	index, err := LoadV2Index(root)
	if err != nil {
		return "", fmt.Errorf("allocate Task ID: strict-load project %s under lock: %w", root, err)
	}
	if validate != nil {
		if err := validate(index); err != nil {
			return "", err
		}
	}

	taskID, err = reserveTaskIDLocked(root, index)
	if err != nil {
		return "", err
	}
	if operation == nil {
		if err := releaseTaskIDLock(lock, lockPath); err != nil {
			locked = false
			return "", fmt.Errorf("allocate Task ID: release lock %s: %w", lockPath, err)
		}
		locked = false
		return taskID, nil
	}

	rollback, err = operation(taskID, index)
	if err != nil {
		taskID = ""
		if rollback != nil {
			if cleanupErr := rollback(); cleanupErr != nil {
				err = errors.Join(err, fmt.Errorf("rollback created Task: %w", cleanupErr))
			}
			rollback = nil
		}
		return "", err
	}

	if err := releaseTaskIDLock(lock, lockPath); err != nil {
		locked = false
		taskID = ""
		lockErr := fmt.Errorf("allocate Task ID: release lock %s: %w", lockPath, err)
		if rollback != nil {
			if cleanupErr := rollback(); cleanupErr != nil {
				lockErr = errors.Join(lockErr, fmt.Errorf("rollback created Task: %w", cleanupErr))
			}
			rollback = nil
		}
		return "", lockErr
	}
	locked = false
	rollback = nil
	return taskID, nil
}

func reserveTaskIDLocked(root string, index *V2Index) (string, error) {
	statePath := filepath.Join(root, taskIDHighWaterFile)
	state, err := readTaskIDHighWater(statePath)
	if err != nil {
		return "", err
	}

	activeMax, err := maxActiveTaskIDSuffix(index.Tasks)
	if err != nil {
		return "", err
	}
	lastIssued := state.lastIssued
	if compareTaskIDDecimal(activeMax, lastIssued) > 0 {
		lastIssued = activeMax
	}

	next := incrementTaskIDDecimal(lastIssued)
	if err := writeTaskIDHighWater(statePath, state, next); err != nil {
		return "", err
	}
	return formatTaskID(next), nil
}

// acquireTaskIDLock takes the lock, waiting up to taskIDLockWait for other
// allocators to finish. A lock still present after that is refused as
// leftover, since no live allocator holds it that long.
func acquireTaskIDLock(path string) (*os.File, error) {
	deadline := time.Now().Add(taskIDLockWait)
	for {
		lock, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
		if err == nil {
			return lock, nil
		}
		if !taskIDLockBusy(err, runtime.GOOS) {
			return nil, fmt.Errorf("allocate Task ID: create lock %s: %w", path, err)
		}
		remaining := time.Until(deadline)
		if remaining <= 0 {
			if !errors.Is(err, os.ErrExist) {
				return nil, fmt.Errorf("allocate Task ID: create lock %s: still failing after %s: %w", path, taskIDLockWait, err)
			}
			return nil, fmt.Errorf("allocate Task ID: lock file %s is present; refusing after %s (a leftover lock requires owner cleanup)", path, taskIDLockWait)
		}
		time.Sleep(min(taskIDLockRetry, remaining))
	}
}

// taskIDLockBusy reports whether an exclusive-create failure means another
// allocator holds or is releasing the lock. On Windows a lock another
// allocator has just removed stays "delete pending" until every handle to it
// closes, and creating it meanwhile fails with access denied rather than
// "exists" (I-061), so that error is retried within the same wait.
func taskIDLockBusy(err error, goos string) bool {
	if errors.Is(err, os.ErrExist) {
		return true
	}
	return goos == "windows" && errors.Is(err, os.ErrPermission)
}

func releaseTaskIDLock(lock *os.File, path string) error {
	closeErr := lock.Close()
	removeErr := os.Remove(path)
	if closeErr != nil && removeErr != nil {
		return errors.Join(fmt.Errorf("close: %w", closeErr), fmt.Errorf("remove: %w", removeErr))
	}
	if closeErr != nil {
		return fmt.Errorf("close: %w", closeErr)
	}
	if removeErr != nil {
		return fmt.Errorf("remove: %w", removeErr)
	}
	return nil
}

func readTaskIDHighWater(path string) (taskIDHighWater, error) {
	info, err := os.Lstat(path)
	if errors.Is(err, os.ErrNotExist) {
		return taskIDHighWater{lastIssued: "0"}, nil
	}
	if err != nil {
		return taskIDHighWater{}, fmt.Errorf("allocate Task ID: inspect high-water mark %s: %w", path, err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return taskIDHighWater{}, fmt.Errorf("allocate Task ID: high-water mark %s is a symlink", path)
	}
	if !info.Mode().IsRegular() {
		return taskIDHighWater{}, fmt.Errorf("allocate Task ID: high-water mark %s is not a regular file", path)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return taskIDHighWater{}, fmt.Errorf("allocate Task ID: read high-water mark %s: %w", path, err)
	}
	lastIssued, err := parseTaskIDHighWater(path, content)
	if err != nil {
		return taskIDHighWater{}, err
	}
	return taskIDHighWater{lastIssued: lastIssued, content: content, info: info, exists: true}, nil
}

func parseTaskIDHighWater(path string, content []byte) (string, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(content))
	var document yaml.Node
	if err := decoder.Decode(&document); err != nil {
		return "", fmt.Errorf("allocate Task ID: malformed high-water mark %s: %w", path, err)
	}
	if document.Kind != yaml.DocumentNode || len(document.Content) != 1 || document.Content[0].Kind != yaml.MappingNode {
		return "", fmt.Errorf("allocate Task ID: malformed high-water mark %s: expected a mapping with last_issued", path)
	}
	mapping := document.Content[0]
	if len(mapping.Content) != 2 || mapping.Content[0].Kind != yaml.ScalarNode || mapping.Content[0].Tag != "!!str" || mapping.Content[0].Value != "last_issued" {
		return "", fmt.Errorf("allocate Task ID: malformed high-water mark %s: expected only last_issued", path)
	}
	value := mapping.Content[1]
	if value.Kind != yaml.ScalarNode || value.Tag != "!!int" || !isTaskIDDecimal(value.Value) {
		return "", fmt.Errorf("allocate Task ID: malformed high-water mark %s: last_issued must be a non-negative decimal integer", path)
	}

	var trailing yaml.Node
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err != nil {
			return "", fmt.Errorf("allocate Task ID: malformed high-water mark %s: %w", path, err)
		}
		return "", fmt.Errorf("allocate Task ID: malformed high-water mark %s: multiple YAML documents", path)
	}
	return normalizeTaskIDDecimal(value.Value), nil
}

func maxActiveTaskIDSuffix(tasks map[string]*TaskV2) (string, error) {
	maximum := "0"
	for id, task := range tasks {
		if task == nil || !strings.HasPrefix(id, "T-") {
			return "", fmt.Errorf("allocate Task ID: strict index contains an invalid Task entry %q", id)
		}
		suffix := id[len("T-"):]
		if !isTaskIDDecimal(suffix) {
			return "", fmt.Errorf("allocate Task ID: strict index contains a non-decimal Task ID %q", id)
		}
		suffix = normalizeTaskIDDecimal(suffix)
		if compareTaskIDDecimal(suffix, maximum) > 0 {
			maximum = suffix
		}
	}
	return maximum, nil
}

func writeTaskIDHighWater(path string, before taskIDHighWater, next string) error {
	mode := os.FileMode(0o644)
	if before.exists {
		mode = before.info.Mode()
		if mode.Perm()&0222 == 0 {
			return fmt.Errorf("allocate Task ID: write high-water mark %s: %w", path, os.ErrPermission)
		}
	}
	content := []byte("last_issued: " + next + "\n")
	if err := replaceV2File(path, content, mode, func() error {
		return checkTaskIDHighWaterFresh(path, before)
	}); err != nil {
		return fmt.Errorf("allocate Task ID: persist high-water mark %s: %w", path, err)
	}
	return nil
}

func checkTaskIDHighWaterFresh(path string, before taskIDHighWater) error {
	latest, err := os.Lstat(path)
	if !before.exists {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("inspect high-water mark %s before replace: %w", path, err)
		}
		return fmt.Errorf("high-water mark %s appeared during allocation", path)
	}
	if err != nil {
		return fmt.Errorf("inspect high-water mark %s before replace: %w", path, err)
	}
	if latest.Mode()&os.ModeSymlink != 0 || !latest.Mode().IsRegular() || !os.SameFile(before.info, latest) || !latest.ModTime().Equal(before.info.ModTime()) || latest.Size() != before.info.Size() {
		return fmt.Errorf("high-water mark %s changed during allocation", path)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read high-water mark %s before replace: %w", path, err)
	}
	if !bytes.Equal(content, before.content) {
		return fmt.Errorf("high-water mark %s changed during allocation", path)
	}
	return nil
}

func isTaskIDDecimal(value string) bool {
	if value == "" {
		return false
	}
	for _, digit := range value {
		if digit < '0' || digit > '9' {
			return false
		}
	}
	return true
}

func normalizeTaskIDDecimal(value string) string {
	trimmed := strings.TrimLeft(value, "0")
	if trimmed == "" {
		return "0"
	}
	return trimmed
}

func compareTaskIDDecimal(left, right string) int {
	left = normalizeTaskIDDecimal(left)
	right = normalizeTaskIDDecimal(right)
	if len(left) < len(right) {
		return -1
	}
	if len(left) > len(right) {
		return 1
	}
	return strings.Compare(left, right)
}

func incrementTaskIDDecimal(value string) string {
	digits := []byte(normalizeTaskIDDecimal(value))
	for i := len(digits) - 1; i >= 0; i-- {
		if digits[i] < '9' {
			digits[i]++
			return string(digits)
		}
		digits[i] = '0'
	}
	return "1" + string(digits)
}

func formatTaskID(suffix string) string {
	if len(suffix) < 3 {
		suffix = strings.Repeat("0", 3-len(suffix)) + suffix
	}
	return "T-" + suffix
}
