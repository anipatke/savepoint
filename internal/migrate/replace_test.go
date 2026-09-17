package migrate

import (
	"errors"
	"go/build"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const originalContent = "original bytes that must survive every failure\n"

func writeDestination(t *testing.T, dir, content string) string {
	t.Helper()
	path := filepath.Join(dir, "destination.md")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write destination: %v", err)
	}
	return path
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(content)
}

// tempFiles lists the replacement temporary files left in dir, which is how
// every failure case checks that nothing was abandoned.
func tempFiles(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir %s: %v", dir, err)
	}
	var left []string
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".savepoint-migrate-") {
			left = append(left, entry.Name())
		}
	}
	return left
}

func TestReplaceFile_replacesExistingDestination(t *testing.T) {
	dir := t.TempDir()
	path := writeDestination(t, dir, originalContent)

	if err := ReplaceFile(path, []byte("replacement bytes\n"), 0644, nil); err != nil {
		t.Fatalf("ReplaceFile() error = %v", err)
	}

	if got := readFile(t, path); got != "replacement bytes\n" {
		t.Errorf("destination = %q, want the replacement content", got)
	}
	if left := tempFiles(t, dir); len(left) != 0 {
		t.Errorf("temporary files left behind: %v", left)
	}
}

// TestReplaceFile_finalCheckRunsAfterTempIsDurable proves the ordering the
// contract promises: the complete temporary file exists before finalCheck
// decides, and the destination is still the original when it is asked.
func TestReplaceFile_finalCheckRunsAfterTempIsDurable(t *testing.T) {
	dir := t.TempDir()
	path := writeDestination(t, dir, originalContent)

	var tempsAtCheck []string
	var destinationAtCheck string
	err := ReplaceFile(path, []byte("replacement bytes\n"), 0644, func() error {
		tempsAtCheck = tempFiles(t, dir)
		destinationAtCheck = readFile(t, path)
		return nil
	})
	if err != nil {
		t.Fatalf("ReplaceFile() error = %v", err)
	}

	if len(tempsAtCheck) != 1 {
		t.Errorf("temporary files during finalCheck = %v, want exactly one", tempsAtCheck)
	}
	if destinationAtCheck != originalContent {
		t.Errorf("destination during finalCheck = %q, want the original content", destinationAtCheck)
	}
}

// TestReplaceFile_rejectedFinalCheckLeavesDestinationIntact covers the failure
// injected in the one window where a complete replacement is already on disk.
func TestReplaceFile_rejectedFinalCheckLeavesDestinationIntact(t *testing.T) {
	dir := t.TempDir()
	path := writeDestination(t, dir, originalContent)
	refused := errors.New("source changed since the preview")

	err := ReplaceFile(path, []byte("replacement bytes\n"), 0644, func() error { return refused })
	if !errors.Is(err, refused) {
		t.Fatalf("ReplaceFile() error = %v, want the finalCheck error", err)
	}

	if got := readFile(t, path); got != originalContent {
		t.Errorf("destination = %q, want the original content", got)
	}
	if left := tempFiles(t, dir); len(left) != 0 {
		t.Errorf("temporary files left behind: %v", left)
	}
}

func TestReplaceFile_missingDestinationIsRefused(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "absent.md")

	err := ReplaceFile(path, []byte("replacement bytes\n"), 0644, nil)
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("ReplaceFile() error = %v, want a not-exist error", err)
	}
	if _, statErr := os.Stat(path); !errors.Is(statErr, os.ErrNotExist) {
		t.Error("ReplaceFile() created the destination it was asked to replace")
	}
	if left := tempFiles(t, dir); len(left) != 0 {
		t.Errorf("temporary files left behind: %v", left)
	}
}

func TestReplaceFile_directoryDestinationIsRefused(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "subdir")
	if err := os.Mkdir(path, 0755); err != nil {
		t.Fatal(err)
	}

	err := ReplaceFile(path, []byte("replacement bytes\n"), 0644, nil)
	if err == nil {
		t.Fatal("ReplaceFile() succeeded on a directory destination")
	}
	if !strings.Contains(err.Error(), "not a regular file") {
		t.Errorf("ReplaceFile() error = %v, want it to name the destination kind", err)
	}
}

// TestReplaceFile_symlinkDestinationIsRefused keeps the contract from writing
// through a link to somewhere outside the project.
func TestReplaceFile_symlinkDestinationIsRefused(t *testing.T) {
	dir := t.TempDir()
	target := writeDestination(t, dir, originalContent)
	link := filepath.Join(dir, "link.md")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	err := ReplaceFile(link, []byte("replacement bytes\n"), 0644, nil)
	if err == nil {
		t.Fatal("ReplaceFile() succeeded on a symlink destination")
	}
	if got := readFile(t, target); got != originalContent {
		t.Errorf("symlink target = %q, want the original content", got)
	}
}

// TestReplaceFile_neverTruncatesTheDestination is the direct evidence for the
// rule this primitive exists to enforce: at no point in a failing replacement
// does the destination become shorter than it started.
func TestReplaceFile_neverTruncatesTheDestination(t *testing.T) {
	dir := t.TempDir()
	path := writeDestination(t, dir, originalContent)
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	originalSize := info.Size()

	err = ReplaceFile(path, []byte("x"), 0644, func() error {
		current, statErr := os.Stat(path)
		if statErr != nil {
			return statErr
		}
		if current.Size() != originalSize {
			t.Errorf("destination size during replacement = %d, want %d", current.Size(), originalSize)
		}
		return errors.New("refuse after the temporary file is durable")
	})
	if err == nil {
		t.Fatal("ReplaceFile() succeeded despite a rejected finalCheck")
	}

	final, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if final.Size() != originalSize {
		t.Errorf("destination size after the failure = %d, want %d", final.Size(), originalSize)
	}
	if got := readFile(t, path); got != originalContent {
		t.Errorf("destination = %q, want the original content", got)
	}
}

// TestReplaceFile_doesNotReachAtomicWrite proves the truncating-copy fallback
// in internal/init is unreachable from this path, structurally: the package
// does not import internal/init at all on the platform under test. It reads
// the package's own directory, so it only runs where the source sits beside
// the test binary — not when a cross-compiled binary is carried to another
// machine to record platform evidence.
func TestReplaceFile_doesNotReachAtomicWrite(t *testing.T) {
	if _, err := os.Stat("replace.go"); err != nil {
		t.Skipf("package source is not beside the test binary: %v", err)
	}
	pkg, err := build.ImportDir(".", build.IgnoreVendor)
	if err != nil {
		t.Fatalf("scan package imports: %v", err)
	}
	for _, imported := range append(pkg.Imports, pkg.TestImports...) {
		if strings.HasSuffix(imported, "/internal/init") {
			t.Errorf("internal/migrate imports %s, which exposes AtomicWrite's truncating-copy fallback", imported)
		}
	}
}

// TestReplaceFile_temporaryFileIsSameDirectory pins the same-volume guarantee
// the contract rests on: a cross-directory temporary file could land on
// another filesystem, where replacement degrades into a copy.
func TestReplaceFile_temporaryFileIsSameDirectory(t *testing.T) {
	dir := t.TempDir()
	path := writeDestination(t, dir, originalContent)

	var observed string
	err := ReplaceFile(path, []byte("replacement bytes\n"), 0644, func() error {
		left := tempFiles(t, dir)
		if len(left) == 1 {
			observed = filepath.Join(dir, left[0])
		}
		return nil
	})
	if err != nil {
		t.Fatalf("ReplaceFile() error = %v", err)
	}
	if observed == "" {
		t.Fatal("no temporary file was created in the destination's directory")
	}
	if filepath.Dir(observed) != filepath.Dir(path) {
		t.Errorf("temporary file directory = %s, want %s", filepath.Dir(observed), filepath.Dir(path))
	}
}

// TestIsTransientReplaceError_ignoresUnrelatedErrors guards the retry policy
// against widening past what the platform experiments observed.
func TestIsTransientReplaceError_ignoresUnrelatedErrors(t *testing.T) {
	for _, err := range []error{os.ErrNotExist, os.ErrPermission, errors.New("some other failure")} {
		if isTransientReplaceError(err) {
			t.Errorf("isTransientReplaceError(%v) = true, want false", err)
		}
	}
}
