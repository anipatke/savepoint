package init

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func AtomicWrite(target string, content []byte) error {
	dir := filepath.Dir(target)
	tmp, err := os.CreateTemp(dir, ".tmp-*.write")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmpName := tmp.Name()

	success := false
	defer func() {
		if !success {
			tmp.Close()
			os.Remove(tmpName)
		}
	}()

	if _, err := tmp.Write(content); err != nil {
		return fmt.Errorf("write temp file: %w", err)
	}

	if err := tmp.Sync(); err != nil {
		return fmt.Errorf("sync temp file: %w", err)
	}

	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temp file: %w", err)
	}

	if err := replaceFile(tmpName, target); err != nil {
		return fmt.Errorf("replace target with temp file: %w", err)
	}

	success = true
	return nil
}

func replaceFile(tmpName, target string) error {
	if err := os.Rename(tmpName, target); err == nil {
		return nil
	}

	src, err := os.Open(tmpName)
	if err != nil {
		return fmt.Errorf("open temp file: %w", err)
	}
	defer src.Close()

	dst, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("create target file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("copy content: %w", err)
	}
	if err := dst.Sync(); err != nil {
		return fmt.Errorf("sync target: %w", err)
	}
	if err := os.Remove(tmpName); err != nil {
		return fmt.Errorf("remove temp file: %w", err)
	}
	return nil
}
