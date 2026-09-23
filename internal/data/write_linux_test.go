//go:build linux

package data

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

// TestWriteObjectiveV2_midWriteFailureLeavesOriginalIntact uses the platform
// file-size limit to interrupt a temp-file write after it has begun. The
// destination must remain byte-identical because replacement happens only
// after the temp file is complete.
func TestWriteObjectiveV2_midWriteFailureLeavesOriginalIntact(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "Objective.md")
	content := "---\nid: O-005\ntitle: \"Large objective\"\nstatus: planned\n---\n\n" + strings.Repeat("authored body that must survive.\n", 300)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	objective, err := DecodeObjectiveV2(path, content)
	if err != nil {
		t.Fatalf("DecodeObjectiveV2() error = %v", err)
	}
	objective.Status = ColumnInProgress

	var oldLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_FSIZE, &oldLimit); err != nil {
		t.Skipf("file-size limit unavailable: %v", err)
	}
	limit := uint64(len(content) - 1)
	if oldLimit.Max < limit || oldLimit.Cur < limit {
		t.Skip("inherited file-size limit cannot run this interruption test")
	}
	defer func() {
		_ = syscall.Setrlimit(syscall.RLIMIT_FSIZE, &oldLimit)
	}()
	limited := oldLimit
	limited.Cur = limit
	if err := syscall.Setrlimit(syscall.RLIMIT_FSIZE, &limited); err != nil {
		t.Skipf("cannot set file-size limit: %v", err)
	}

	err = WriteObjectiveV2(objective)
	if err == nil {
		t.Fatal("WriteObjectiveV2() succeeded despite a mid-write file-size failure")
	}
	if !errors.Is(err, syscall.EFBIG) {
		t.Fatalf("WriteObjectiveV2() error = %v, want EFBIG", err)
	}
	if !strings.Contains(err.Error(), path) {
		t.Errorf("WriteObjectiveV2() error = %v, want destination path %q", err, path)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != content {
		t.Fatalf("destination was truncated or changed after interrupted write: got %d bytes, want %d", len(got), len(content))
	}
}
