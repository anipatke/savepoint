//go:build !windows

package migrate

import (
	"os"
	"strings"
	"testing"
)

// TestReplaceFile_modeBecomesTheReplacedFileMode is the Unix half of the
// documented platform difference: rename carries the temporary file's
// permission bits onto the destination, so the mode argument decides what the
// replaced file ends up with. The Windows half is
// TestReplaceFile_preservesDestinationAttributesAndACL, where the
// destination's own attributes win instead.
func TestReplaceFile_modeBecomesTheReplacedFileMode(t *testing.T) {
	dir := t.TempDir()
	path := writeDestination(t, dir, originalContent)
	if err := os.Chmod(path, 0644); err != nil {
		t.Fatal(err)
	}

	if err := ReplaceFile(path, []byte("replacement bytes\n"), 0600, nil); err != nil {
		t.Fatalf("ReplaceFile() error = %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0600 {
		t.Errorf("destination mode = %04o, want 0600", got)
	}
}

// TestReplaceFile_openDestinationDoesNotBlockReplacement records the Unix
// behavior the retry policy's absence rests on: an open handle keeps referring
// to the replaced file, so rename is never refused for contention and there is
// nothing for a retry to clear.
func TestReplaceFile_openDestinationDoesNotBlockReplacement(t *testing.T) {
	dir := t.TempDir()
	path := writeDestination(t, dir, originalContent)

	held, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer held.Close()

	if err := ReplaceFile(path, []byte("replacement bytes\n"), 0644, nil); err != nil {
		t.Fatalf("ReplaceFile() error = %v", err)
	}
	if got := readFile(t, path); got != "replacement bytes\n" {
		t.Errorf("destination = %q, want the replacement content", got)
	}
}

// TestReplaceFile_refusedReplacementLeavesDestinationIntact fails the
// replacement itself by making the destination's directory unwritable, which
// is the case that would reach AtomicWrite's truncating copy if this path had
// one. It is Unix-tagged because a directory permission bit does not restrict
// this on Windows; the Windows cases that refuse a replacement are
// TestReplaceFile_exhaustedRetriesNameThePathAndCondition and
// TestReplaceFile_readOnlyDestinationIsNotRetried.
func TestReplaceFile_refusedReplacementLeavesDestinationIntact(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root ignores the directory permission this case depends on")
	}
	dir := t.TempDir()
	path := writeDestination(t, dir, originalContent)

	// The temporary file has to exist before the directory is sealed, so the
	// failure lands on the replacement rather than on creating the temp file.
	var sealed bool
	err := ReplaceFile(path, []byte("replacement bytes\n"), 0644, func() error {
		if err := os.Chmod(dir, 0500); err != nil {
			return err
		}
		sealed = true
		return nil
	})
	if sealed {
		defer os.Chmod(dir, 0700)
	}
	if err == nil {
		t.Fatal("ReplaceFile() succeeded despite an unwritable directory")
	}
	if !strings.Contains(err.Error(), path) {
		t.Errorf("ReplaceFile() error = %v, want the destination path named", err)
	}

	if got := readFile(t, path); got != originalContent {
		t.Errorf("destination = %q, want the original content", got)
	}
}
