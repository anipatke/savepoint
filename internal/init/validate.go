package init

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	ErrTargetMissing    = errors.New("target directory does not exist")
	ErrNotADirectory    = errors.New("target is not a directory")
	ErrPermissionDenied = errors.New("permission denied")
	ErrAlreadyInit      = errors.New("target already contains .savepoint directory")
	ErrConflict         = errors.New("target has conflicting files")
)

var conflictingFiles = []string{
	"AGENTS.md",
	"agent-skills",
}

type ValidationError struct {
	Type    error
	Message string
}

func (e *ValidationError) Error() string { return e.Message }
func (e *ValidationError) Unwrap() error { return e.Type }

func ValidateTarget(path string, force bool) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return &ValidationError{Type: ErrPermissionDenied, Message: fmt.Sprintf("cannot resolve path %q: permission denied", path)}
	}

	info, err := os.Stat(abs)
	if err != nil {
		if os.IsNotExist(err) {
			return &ValidationError{Type: ErrTargetMissing, Message: fmt.Sprintf("target directory %q does not exist", path)}
		}
		if os.IsPermission(err) {
			return &ValidationError{Type: ErrPermissionDenied, Message: fmt.Sprintf("permission denied accessing %q", path)}
		}
		return &ValidationError{Type: ErrPermissionDenied, Message: fmt.Sprintf("cannot access %q: %v", path, err)}
	}

	if !info.IsDir() {
		return &ValidationError{Type: ErrNotADirectory, Message: fmt.Sprintf("target %q is not a directory", path)}
	}

	if err := checkWritable(abs); err != nil {
		return err
	}

	savepointDir := filepath.Join(abs, ".savepoint")
	if _, err := os.Stat(savepointDir); err == nil {
		if !force {
			return &ValidationError{Type: ErrAlreadyInit, Message: fmt.Sprintf("target %q already contains a .savepoint directory (use --force to overwrite)", path)}
		}
	}

	if !force {
		for _, name := range conflictingFiles {
			conflictPath := filepath.Join(abs, name)
			if _, err := os.Stat(conflictPath); err == nil {
				return &ValidationError{Type: ErrConflict, Message: fmt.Sprintf("target %q has conflicting file %q (use --force to overwrite)", path, name)}
			}
		}
	}

	return nil
}

func checkWritable(dir string) error {
	testFile := filepath.Join(dir, ".savepoint-write-test")
	if err := os.WriteFile(testFile, []byte{}, 0644); err != nil {
		if os.IsPermission(err) {
			return &ValidationError{Type: ErrPermissionDenied, Message: fmt.Sprintf("target directory %q is not writable", dir)}
		}
		return &ValidationError{Type: ErrPermissionDenied, Message: fmt.Sprintf("cannot write to %q: %v", dir, err)}
	}
	os.Remove(testFile)
	return nil
}
