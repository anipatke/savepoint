// Package migrate owns the one-time conversion of a V1 Savepoint project into
// a V2 one. This file holds the replacement primitive every write in that
// operation goes through.
package migrate

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// replaceAttempts and replaceRetryDelay bound the retry for a replacement
// failure the platform reports as transient. Five attempts spaced 100ms apart
// give a scanner or editor that momentarily holds a record open under half a
// second to let go, which is long enough to clear the contention observed on
// Windows and short enough that an operator watching a migration does not read
// it as a hang. A condition the platform does not report as transient is never
// retried, so this bound only ever applies to real contention.
const (
	replaceAttempts    = 5
	replaceRetryDelay  = 100 * time.Millisecond
	replaceTempPattern = ".savepoint-migrate-*"
)

// ReplaceFile replaces the existing file at path with content, and is the only
// way migration overwrites a file a user already has.
//
// It writes a complete temporary file in path's own directory, flushes and
// closes it, runs finalCheck, and only then replaces the destination. Nothing
// truncates or removes the destination at any point: until the replacement
// succeeds the original file is untouched on disk, so every failure — a write
// error, a rejected finalCheck, or a refused replacement — leaves the
// destination byte-identical. The temporary file is removed on every failure
// path; when it cannot be removed, the returned error names it rather than
// leaving it unreported. A replacement the platform reports as transient is
// retried up to replaceAttempts times; any other failure is returned on the
// first attempt.
//
// path must already exist and be a regular file. A missing destination is an
// error on both platforms rather than a creation, because Windows'
// ReplaceFileW refuses it (ERROR_FILE_NOT_FOUND) while Unix' rename would
// silently create it; refusing uniformly keeps the contract from meaning two
// different things. Additive records are created through a create-only path,
// never through this one.
//
// mode sets the temporary file's permission bits. On Unix those bits become
// the replaced file's, because rename carries the temporary file's metadata
// across. On Windows they do not: ReplaceFileW keeps the destination's own
// attributes and ACLs, which is the documented platform difference here and
// the reason it was chosen. Callers that pass the destination's current mode —
// which is what migration does — get the same result on both platforms.
func ReplaceFile(path string, content []byte, mode os.FileMode, finalCheck func() error) (retErr error) {
	info, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("replace %s: %w", path, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("replace %s: not a regular file", path)
	}

	tempPath, err := writeReplacementTemp(path, content, mode)
	if err != nil {
		return err
	}
	defer func() {
		if tempPath == "" {
			return
		}
		if removeErr := os.Remove(tempPath); removeErr != nil && !os.IsNotExist(removeErr) {
			if retErr != nil {
				retErr = fmt.Errorf("%w (cleanup %s: %v)", retErr, tempPath, removeErr)
			} else {
				retErr = fmt.Errorf("cleanup %s: %w", tempPath, removeErr)
			}
		}
	}()

	if finalCheck != nil {
		if err := finalCheck(); err != nil {
			return err
		}
	}

	if err := replaceWithRetry(tempPath, path); err != nil {
		return err
	}
	tempPath = ""
	return nil
}

// writeReplacementTemp writes content to a complete, flushed, closed temporary
// file in path's own directory and returns its name. The directory is path's
// own so the replacement stays within one volume: a cross-directory temporary
// file can land on another filesystem, where the replacement stops being a
// metadata operation and becomes a copy — exactly the truncating fallback this
// package refuses. Every failure removes the temporary file before returning.
func writeReplacementTemp(path string, content []byte, mode os.FileMode) (string, error) {
	temp, err := os.CreateTemp(filepath.Dir(path), replaceTempPattern)
	if err != nil {
		return "", fmt.Errorf("replace %s: create temporary file: %w", path, err)
	}
	tempPath := temp.Name()

	fail := func(err error) (string, error) {
		if removeErr := os.Remove(tempPath); removeErr != nil && !os.IsNotExist(removeErr) {
			return "", fmt.Errorf("%w (cleanup %s: %v)", err, tempPath, removeErr)
		}
		return "", err
	}
	closeThenFail := func(err error) (string, error) {
		if closeErr := temp.Close(); closeErr != nil {
			return fail(fmt.Errorf("%w (close: %v)", err, closeErr))
		}
		return fail(err)
	}

	if err := temp.Chmod(mode.Perm()); err != nil {
		return closeThenFail(fmt.Errorf("replace %s: set temporary permissions: %w", path, err))
	}
	written, writeErr := temp.Write(content)
	if writeErr == nil && written != len(content) {
		writeErr = fmt.Errorf("short write: wrote %d of %d bytes", written, len(content))
	}
	if writeErr != nil {
		return closeThenFail(fmt.Errorf("replace %s: %w", path, writeErr))
	}
	if err := temp.Sync(); err != nil {
		return closeThenFail(fmt.Errorf("replace %s: flush temporary file: %w", path, err))
	}
	if err := temp.Close(); err != nil {
		return fail(fmt.Errorf("replace %s: close temporary file: %w", path, err))
	}
	return tempPath, nil
}

// replaceWithRetry calls the platform replacement, retrying only while the
// platform reports the failure as transient. The loop is shared so both
// platforms bound contention the same way; what counts as transient is the
// only part that varies, and each platform's isTransientReplaceError records
// what its experiment actually observed.
func replaceWithRetry(tempPath, path string) error {
	var err error
	for attempt := 1; attempt <= replaceAttempts; attempt++ {
		if err = replaceDestination(tempPath, path); err == nil {
			return nil
		}
		if !isTransientReplaceError(err) {
			return fmt.Errorf("replace %s: %w", path, err)
		}
		if attempt < replaceAttempts {
			time.Sleep(replaceRetryDelay)
		}
	}
	return fmt.Errorf("replace %s: still in use after %d attempts over %s: %w",
		path, replaceAttempts, replaceRetryDelay*(replaceAttempts-1), err)
}
