// Package data's runtime.go is the ordinary commands' front door: locating a
// project directory and deciding whether its schema is one the V2 runtime can
// read, without pulling in internal/migrate's conversion engine. Every command
// that only needs to know "is this a usable V2 project" calls CheckRuntimeSchema;
// internal/migrate itself calls ResolveTarget and FindProjectRoot through its
// own thin wrappers so there remains exactly one implementation of each.
package data

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Distinct project-target diagnostics, so a caller (and a test) can tell a
// missing directory apart from one that is not a Savepoint project at all.
var (
	ErrTargetMissing      = errors.New("target directory does not exist")
	ErrTargetNotSavepoint = errors.New("target directory is not a Savepoint project")
	// ErrSchemaMigrationRequired means a project's config.yml declares
	// schema_version 1, or no schema_version at all: the V2 runtime cannot
	// read it, and `savepoint migrate` is the only command allowed to change
	// what the project is.
	ErrSchemaMigrationRequired = errors.New("schema_version 1: run `savepoint migrate --dry-run`, then `savepoint migrate --apply`")
)

// ResolveTarget validates dir as a project target and returns its absolute
// path. It checks dir itself rather than walking up through parent
// directories: a command given an explicit directory must never silently
// operate on an unrelated ancestor project.
func ResolveTarget(dir string) (string, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("%w: %s: %v", ErrTargetMissing, dir, err)
	}

	info, statErr := os.Stat(abs)
	if statErr != nil {
		if os.IsNotExist(statErr) {
			return "", fmt.Errorf("%w: %s", ErrTargetMissing, dir)
		}
		return "", fmt.Errorf("%w: %s: %v", ErrTargetMissing, dir, statErr)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("%w: %s is not a directory", ErrTargetMissing, dir)
	}

	savepointDir := filepath.Join(abs, ".savepoint")
	if sInfo, sErr := os.Stat(savepointDir); sErr != nil || !sInfo.IsDir() {
		return "", fmt.Errorf("%w: %s has no .savepoint directory", ErrTargetNotSavepoint, dir)
	}

	return abs, nil
}

// FindProjectRoot locates the nearest Savepoint project while walking from
// start toward the filesystem root. It only identifies the project boundary
// and never discovers or parses any project record.
func FindProjectRoot(start string) (string, error) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", fmt.Errorf("resolve project root: %w", err)
	}

	for {
		info, statErr := os.Stat(filepath.Join(dir, ".savepoint"))
		if statErr == nil && info.IsDir() {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("%w: %s has no .savepoint directory", ErrTargetNotSavepoint, start)
		}
		dir = parent
	}
}

// CheckRuntimeSchema reads root's config.yml (root is a project root, the
// directory containing .savepoint) and reports whether the ordinary V2
// runtime may use it: nil for schema 2, ErrSchemaMigrationRequired for
// schema 1 (including a project with no config.yml at all), and the named
// ReadSchemaVersion sentinel for an unsupported or malformed schema_version.
// An unreadable config.yml is reported as its own wrapped error. It performs
// no other read and no write.
func CheckRuntimeSchema(root string) error {
	configPath := filepath.Join(root, ".savepoint", "config.yml")
	version, err := ReadSchemaVersion(configPath)
	if err != nil {
		if errors.Is(err, ErrUnsupportedSchemaVersion) || errors.Is(err, ErrMalformedSchemaVersion) {
			return err
		}
		return fmt.Errorf("cannot read schema at %s: %w", configPath, err)
	}
	if version != SchemaVersionV2 {
		return ErrSchemaMigrationRequired
	}
	return nil
}
