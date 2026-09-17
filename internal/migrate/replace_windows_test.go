//go:build windows

package migrate

// This file is the Windows half of the E45 T001 platform experiment and the
// platform-tagged evidence CFG-02 requires for the replacement contract. The
// probe tests call the candidate primitives directly rather than through
// ReplaceFile, because the decision this task records is about what Windows
// does, and a test that only exercised the chosen wrapper could not show why
// the alternatives were rejected.

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
	"unsafe"
)

const (
	envHoldPath  = "SAVEPOINT_E45_HOLD_PATH"
	envHoldShare = "SAVEPOINT_E45_HOLD_SHARE"
	envCrashDir  = "SAVEPOINT_E45_CRASH_DIR"
)

// The probe calls the production ReplaceFileW binding in replace_windows.go
// rather than its own copy, so the evidence below describes the call the
// operation actually makes.
var procGetVolumeInformationW = kernel32.NewProc("GetVolumeInformationW")

const (
	errorFileNotFound = syscall.Errno(2)
	errorAccessDenied = syscall.Errno(5)
)

// errnoName names the Windows error behind err, so the recorded evidence
// distinguishes a sharing violation from an access denial rather than leaving
// two different conditions behind the same English sentence.
func errnoName(err error) string {
	if err == nil {
		return "none"
	}
	var errno syscall.Errno
	if !errors.As(err, &errno) {
		return "not-an-errno"
	}
	switch errno {
	case errorFileNotFound:
		return "ERROR_FILE_NOT_FOUND(2)"
	case errorAccessDenied:
		return "ERROR_ACCESS_DENIED(5)"
	case errorSharingViolation:
		return "ERROR_SHARING_VIOLATION(32)"
	default:
		return fmt.Sprintf("errno(%d)", uintptr(errno))
	}
}

// TestMain lets the test binary re-execute itself as the "other process" the
// acceptance criteria name: a real second process holding a handle with a
// chosen share mode, and a real process that dies between writing the
// temporary file and replacing the destination.
func TestMain(m *testing.M) {
	if path := os.Getenv(envHoldPath); path != "" {
		os.Exit(runHoldOpen(path, os.Getenv(envHoldShare)))
	}
	if dir := os.Getenv(envCrashDir); dir != "" {
		os.Exit(runCrashBeforeReplace(dir))
	}
	os.Exit(m.Run())
}

// runHoldOpen opens path with the requested share mode, announces readiness on
// stdout, and holds the handle until the parent closes its stdin pipe.
func runHoldOpen(path, share string) int {
	var shareMode uint32
	if _, err := fmt.Sscanf(share, "%d", &shareMode); err != nil {
		fmt.Printf("HOLD-ERROR bad share mode %q: %v\n", share, err)
		return 2
	}
	namePtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		fmt.Printf("HOLD-ERROR bad path: %v\n", err)
		return 2
	}
	handle, err := syscall.CreateFile(namePtr, syscall.GENERIC_READ, shareMode, nil, syscall.OPEN_EXISTING, syscall.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		fmt.Printf("HOLD-ERROR open: %v\n", err)
		return 2
	}
	defer syscall.CloseHandle(handle)
	fmt.Println("HOLD-READY")
	bufio.NewReader(os.Stdin).ReadString('\n')
	return 0
}

// runCrashBeforeReplace writes a complete temporary file next to the
// destination and then dies without replacing anything, which is the exact
// interruption window the operation journal has to survive.
func runCrashBeforeReplace(dir string) int {
	temp, err := os.CreateTemp(dir, ".savepoint-migrate-*")
	if err != nil {
		fmt.Printf("CRASH-ERROR create: %v\n", err)
		return 2
	}
	if _, err := temp.WriteString("replacement bytes"); err != nil {
		fmt.Printf("CRASH-ERROR write: %v\n", err)
		return 2
	}
	if err := temp.Sync(); err != nil {
		fmt.Printf("CRASH-ERROR sync: %v\n", err)
		return 2
	}
	if err := temp.Close(); err != nil {
		fmt.Printf("CRASH-ERROR close: %v\n", err)
		return 2
	}
	fmt.Printf("CRASH-TEMP %s\n", temp.Name())
	os.Exit(3)
	return 3
}

// helperCommand re-executes this test binary by its absolute path. os.Args[0]
// is relative when the binary is launched from its own directory, which
// os/exec refuses to resolve.
func helperCommand(t *testing.T) *exec.Cmd {
	t.Helper()
	exe, err := os.Executable()
	if err != nil {
		t.Fatalf("resolve test binary: %v", err)
	}
	return exec.Command(exe)
}

// holdOpen starts the helper process and returns once it reports the handle is
// open. The returned release closes the pipe, which unblocks and ends it.
func holdOpen(t *testing.T, path string, shareMode uint32) (release func()) {
	t.Helper()
	cmd := helperCommand(t)
	cmd.Env = append(os.Environ(), fmt.Sprintf("%s=%s", envHoldPath, path), fmt.Sprintf("%s=%d", envHoldShare, shareMode))
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("stdin pipe: %v", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("stdout pipe: %v", err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatalf("start holder: %v", err)
	}
	line, err := bufio.NewReader(stdout).ReadString('\n')
	if err != nil || !strings.HasPrefix(line, "HOLD-READY") {
		stdin.Close()
		cmd.Wait()
		t.Fatalf("holder did not open %s with share mode %d: %q (%v)", path, shareMode, strings.TrimSpace(line), err)
	}
	t.Logf("holder process %d holds %s open with share mode 0x%x", cmd.Process.Pid, filepath.Base(path), shareMode)
	released := false
	return func() {
		if released {
			return
		}
		released = true
		stdin.Close()
		cmd.Wait()
	}
}

// moveFileExReplace is the rejected candidate. os.Rename on Windows is
// MoveFileEx with MOVEFILE_REPLACE_EXISTING — internal/syscall/windows.Rename
// makes exactly that call — so exercising os.Rename exercises the candidate.
func moveFileExReplace(from, to string) error {
	return os.Rename(from, to)
}

// volumeFileSystem reports the filesystem backing dir, so the transcript can
// prove the run happened on NTFS rather than on a WSL ext4 mount.
func volumeFileSystem(t *testing.T, dir string) string {
	t.Helper()
	volume := filepath.VolumeName(dir) + `\`
	rootPtr, err := syscall.UTF16PtrFromString(volume)
	if err != nil {
		t.Fatalf("volume name %q: %v", volume, err)
	}
	nameBuf := make([]uint16, 261)
	fsBuf := make([]uint16, 261)
	r1, _, e1 := syscall.SyscallN(procGetVolumeInformationW.Addr(),
		uintptr(unsafe.Pointer(rootPtr)),
		uintptr(unsafe.Pointer(&nameBuf[0])), uintptr(len(nameBuf)),
		0, 0, 0,
		uintptr(unsafe.Pointer(&fsBuf[0])), uintptr(len(fsBuf)))
	if r1 == 0 {
		t.Fatalf("GetVolumeInformationW(%s): %v", volume, e1)
	}
	return syscall.UTF16ToString(fsBuf)
}

func writeTemp(t *testing.T, dir, content string) string {
	t.Helper()
	temp, err := os.CreateTemp(dir, ".savepoint-migrate-*")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	if _, err := temp.WriteString(content); err != nil {
		t.Fatalf("write temp: %v", err)
	}
	if err := temp.Sync(); err != nil {
		t.Fatalf("sync temp: %v", err)
	}
	if err := temp.Close(); err != nil {
		t.Fatalf("close temp: %v", err)
	}
	return temp.Name()
}

func fileAttributes(t *testing.T, path string) uint32 {
	t.Helper()
	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		t.Fatalf("path %q: %v", path, err)
	}
	attrs, err := syscall.GetFileAttributes(pathPtr)
	if err != nil {
		t.Fatalf("GetFileAttributes(%s): %v", path, err)
	}
	return attrs
}

func setFileAttributes(t *testing.T, path string, attrs uint32) {
	t.Helper()
	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		t.Fatalf("path %q: %v", path, err)
	}
	if err := syscall.SetFileAttributes(pathPtr, attrs); err != nil {
		t.Fatalf("SetFileAttributes(%s, 0x%x): %v", path, attrs, err)
	}
}

// icacls reports the destination's access control list as icacls prints it,
// which is the readable form of the ACL question the task has to answer.
func icacls(t *testing.T, args ...string) string {
	t.Helper()
	out, err := exec.Command("icacls.exe", args...).CombinedOutput()
	if err != nil {
		t.Fatalf("icacls %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

// aclLines keeps only the ACE lines, dropping icacls' trailing summary and the
// path prefix that differs between two temporary directories.
func aclLines(raw, path string) []string {
	var aces []string
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), path))
		if line == "" || strings.Contains(line, "Successfully processed") {
			continue
		}
		aces = append(aces, line)
	}
	return aces
}

func TestProbeVolumeIsNTFS(t *testing.T) {
	dir := t.TempDir()
	fs := volumeFileSystem(t, dir)
	t.Logf("EVIDENCE volume: temp dir %s is on filesystem %q", dir, fs)
	if fs != "NTFS" {
		t.Fatalf("experiment must run on NTFS, got %q for %s", fs, dir)
	}
}

// TestProbeDestinationHeldOpen records what each candidate does when another
// process holds the destination open, for both the permissive share mode Go's
// own os.Open uses and the deny-all share mode a hostile holder uses.
func TestProbeDestinationHeldOpen(t *testing.T) {
	shareModes := []struct {
		name  string
		share uint32
	}{
		{"share-read-write-delete", syscall.FILE_SHARE_READ | syscall.FILE_SHARE_WRITE | syscall.FILE_SHARE_DELETE},
		{"share-read-only", syscall.FILE_SHARE_READ},
		{"share-none", 0},
	}
	for _, mode := range shareModes {
		t.Run(mode.name, func(t *testing.T) {
			for _, primitive := range []string{"MoveFileEx", "ReplaceFileW"} {
				t.Run(primitive, func(t *testing.T) {
					dir := t.TempDir()
					dest := writeDestination(t, dir, "original bytes")
					temp := writeTemp(t, dir, "replacement bytes")
					release := holdOpen(t, dest, mode.share)
					defer release()

					var err error
					if primitive == "MoveFileEx" {
						err = moveFileExReplace(temp, dest)
					} else {
						err = replaceFileW(dest, temp, "", 0)
					}
					t.Logf("EVIDENCE open-dest %s share=%s: err=%v errno=%s", primitive, mode.name, err, errnoName(err))

					// The destination is only readable once the holder lets
					// go: a deny-all holder blocks this test's own read too.
					release()
					got, readErr := os.ReadFile(dest)
					if readErr != nil {
						t.Fatalf("read destination: %v", readErr)
					}
					if err == nil {
						if string(got) != "replacement bytes" {
							t.Fatalf("%s reported success but destination holds %q", primitive, got)
						}
					} else if string(got) != "original bytes" {
						t.Fatalf("%s failed but destination is no longer the original: %q", primitive, got)
					}
				})
			}
		})
	}
}

// TestProbeReadOnlyDestination records what each candidate does when the
// destination carries FILE_ATTRIBUTE_READONLY, and what the attribute looks
// like afterwards.
func TestProbeReadOnlyDestination(t *testing.T) {
	for _, primitive := range []string{"MoveFileEx", "ReplaceFileW"} {
		t.Run(primitive, func(t *testing.T) {
			dir := t.TempDir()
			dest := writeDestination(t, dir, "original bytes")
			temp := writeTemp(t, dir, "replacement bytes")
			setFileAttributes(t, dest, syscall.FILE_ATTRIBUTE_READONLY)
			defer func() {
				if _, err := os.Stat(dest); err == nil {
					setFileAttributes(t, dest, syscall.FILE_ATTRIBUTE_NORMAL)
				}
			}()

			var err error
			if primitive == "MoveFileEx" {
				err = moveFileExReplace(temp, dest)
			} else {
				err = replaceFileW(dest, temp, "", 0)
			}
			after := fileAttributes(t, dest)
			got, readErr := os.ReadFile(dest)
			if readErr != nil {
				t.Fatalf("read destination: %v", readErr)
			}
			t.Logf("EVIDENCE read-only %s: err=%v errno=%s attrs-after=0x%x readonly-after=%t content=%q",
				primitive, err, errnoName(err), after, after&syscall.FILE_ATTRIBUTE_READONLY != 0, got)
			if err != nil && string(got) != "original bytes" {
				t.Fatalf("%s failed but destination is no longer the original: %q", primitive, got)
			}
		})
	}
}

// TestProbeDestinationAttributesSurvive records whether a non-default
// attribute set on the destination survives the replacement.
func TestProbeDestinationAttributesSurvive(t *testing.T) {
	for _, primitive := range []string{"MoveFileEx", "ReplaceFileW"} {
		t.Run(primitive, func(t *testing.T) {
			dir := t.TempDir()
			dest := writeDestination(t, dir, "original bytes")
			temp := writeTemp(t, dir, "replacement bytes")
			setFileAttributes(t, dest, syscall.FILE_ATTRIBUTE_HIDDEN)
			before := fileAttributes(t, dest)

			var err error
			if primitive == "MoveFileEx" {
				err = moveFileExReplace(temp, dest)
			} else {
				err = replaceFileW(dest, temp, "", 0)
			}
			if err != nil {
				t.Fatalf("%s on a hidden destination: %v", primitive, err)
			}
			after := fileAttributes(t, dest)
			t.Logf("EVIDENCE attributes %s: before=0x%x after=0x%x hidden-preserved=%t",
				primitive, before, after, after&syscall.FILE_ATTRIBUTE_HIDDEN != 0)
		})
	}
}

// TestProbeDestinationACLSurvives records whether an explicit ACE added to the
// destination survives the replacement.
func TestProbeDestinationACLSurvives(t *testing.T) {
	for _, primitive := range []string{"MoveFileEx", "ReplaceFileW"} {
		t.Run(primitive, func(t *testing.T) {
			dir := t.TempDir()
			dest := writeDestination(t, dir, "original bytes")
			temp := writeTemp(t, dir, "replacement bytes")
			icacls(t, dest, "/grant", "*S-1-1-0:(RX)")
			before := aclLines(icacls(t, dest), dest)

			var err error
			if primitive == "MoveFileEx" {
				err = moveFileExReplace(temp, dest)
			} else {
				err = replaceFileW(dest, temp, "", 0)
			}
			if err != nil {
				t.Fatalf("%s on a destination with an explicit ACE: %v", primitive, err)
			}
			after := aclLines(icacls(t, dest), dest)
			t.Logf("EVIDENCE acl %s: before=%v after=%v preserved=%t",
				primitive, before, after, strings.Join(before, "|") == strings.Join(after, "|"))
		})
	}
}

// TestProbeMissingDestination records what each candidate does when there is
// nothing to replace, which decides whether the contract can cover creating a
// record as well as replacing one.
func TestProbeMissingDestination(t *testing.T) {
	for _, primitive := range []string{"MoveFileEx", "ReplaceFileW"} {
		t.Run(primitive, func(t *testing.T) {
			dir := t.TempDir()
			dest := filepath.Join(dir, "destination.md")
			temp := writeTemp(t, dir, "replacement bytes")

			var err error
			if primitive == "MoveFileEx" {
				err = moveFileExReplace(temp, dest)
			} else {
				err = replaceFileW(dest, temp, "", 0)
			}
			_, statErr := os.Stat(dest)
			t.Logf("EVIDENCE missing-dest %s: err=%v errno=%s destination-created=%t",
				primitive, err, errnoName(err), statErr == nil)
		})
	}
}

// TestProbeInterruptionBeforeReplace records the on-disk state when the
// process dies after the temporary file is durable and before the
// replacement, which is the window every later recovery task is built on.
func TestProbeInterruptionBeforeReplace(t *testing.T) {
	dir := t.TempDir()
	dest := writeDestination(t, dir, "original bytes")

	cmd := helperCommand(t)
	cmd.Env = append(os.Environ(), fmt.Sprintf("%s=%s", envCrashDir, dir))
	out, err := cmd.CombinedOutput()
	if cmd.ProcessState == nil || cmd.ProcessState.ExitCode() != 3 {
		t.Fatalf("crash helper exit = %v (err %v), output:\n%s", cmd.ProcessState, err, out)
	}
	tempName := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(string(out)), "CRASH-TEMP"))

	got, readErr := os.ReadFile(dest)
	if readErr != nil {
		t.Fatalf("read destination: %v", readErr)
	}
	_, tempErr := os.Stat(tempName)
	t.Logf("EVIDENCE interruption: destination=%q temp-present=%t temp=%s", got, tempErr == nil, tempName)
	if string(got) != "original bytes" {
		t.Fatalf("destination changed despite the interruption happening before replacement: %q", got)
	}
}

// TestProbeSharingViolationIsTransient records whether a failed replacement
// against a held-open destination succeeds once the holder releases it, which
// is what decides whether a bounded retry is honest or theatre.
func TestProbeSharingViolationIsTransient(t *testing.T) {
	for _, primitive := range []string{"MoveFileEx", "ReplaceFileW"} {
		t.Run(primitive, func(t *testing.T) {
			dir := t.TempDir()
			dest := writeDestination(t, dir, "original bytes")
			temp := writeTemp(t, dir, "replacement bytes")
			release := holdOpen(t, dest, 0)

			attempt := func() error {
				if primitive == "MoveFileEx" {
					return moveFileExReplace(temp, dest)
				}
				return replaceFileW(dest, temp, "", 0)
			}
			first := attempt()
			release()
			second := attempt()
			t.Logf("EVIDENCE transient %s: while-held=%v errno=%s after-release=%v",
				primitive, first, errnoName(first), second)
			if first == nil {
				t.Skipf("%s was not blocked by the holder, so there is nothing to retry", primitive)
			}
			if second != nil {
				t.Fatalf("%s still failed after the holder released: %v", primitive, second)
			}
		})
	}
}

// The tests below exercise ReplaceFile itself rather than the raw candidates:
// they are the platform-tagged coverage CFG-02 requires for the two behaviors
// the contract documents as differing from Unix.

// TestReplaceFile_preservesDestinationAttributesAndACL is the Windows half of
// the documented platform difference. The destination's own attributes and
// ACL survive the replacement and the mode argument does not decide them,
// which is the property MoveFileEx lacks — see TestProbeDestinationAttributes
// Survive and TestProbeDestinationACLSurvives for the rejected candidate's
// behavior on the same inputs. The Unix half is
// TestReplaceFile_modeBecomesTheReplacedFileMode.
func TestReplaceFile_preservesDestinationAttributesAndACL(t *testing.T) {
	dir := t.TempDir()
	path := writeDestination(t, dir, originalContent)
	icacls(t, path, "/grant", "*S-1-1-0:(RX)")
	setFileAttributes(t, path, syscall.FILE_ATTRIBUTE_HIDDEN)
	beforeACL := aclLines(icacls(t, path), path)

	if err := ReplaceFile(path, []byte("replacement bytes\n"), 0600, nil); err != nil {
		t.Fatalf("ReplaceFile() error = %v", err)
	}

	if got := readFile(t, path); got != "replacement bytes\n" {
		t.Errorf("destination = %q, want the replacement content", got)
	}
	if attrs := fileAttributes(t, path); attrs&syscall.FILE_ATTRIBUTE_HIDDEN == 0 {
		t.Errorf("destination attributes = 0x%x, want FILE_ATTRIBUTE_HIDDEN preserved", attrs)
	}
	afterACL := aclLines(icacls(t, path), path)
	if strings.Join(beforeACL, "|") != strings.Join(afterACL, "|") {
		t.Errorf("destination ACL = %v, want %v", afterACL, beforeACL)
	}
}

// TestReplaceFile_succeedsWhileDestinationIsOpen covers the case that decided
// the primitive: another process holding the destination open under the share
// mode Go's own os.Open uses does not stop the replacement.
func TestReplaceFile_succeedsWhileDestinationIsOpen(t *testing.T) {
	dir := t.TempDir()
	path := writeDestination(t, dir, originalContent)
	release := holdOpen(t, path, syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE|syscall.FILE_SHARE_DELETE)
	defer release()

	if err := ReplaceFile(path, []byte("replacement bytes\n"), 0644, nil); err != nil {
		t.Fatalf("ReplaceFile() error = %v", err)
	}

	release()
	if got := readFile(t, path); got != "replacement bytes\n" {
		t.Errorf("destination = %q, want the replacement content", got)
	}
}

// TestReplaceFile_exhaustedRetriesNameThePathAndCondition covers the other end
// of the retry policy: a holder that never lets go produces a truthful error
// naming the destination and the condition, leaves the destination unchanged,
// and abandons no temporary file.
func TestReplaceFile_exhaustedRetriesNameThePathAndCondition(t *testing.T) {
	dir := t.TempDir()
	path := writeDestination(t, dir, originalContent)
	release := holdOpen(t, path, 0)
	defer release()

	err := ReplaceFile(path, []byte("replacement bytes\n"), 0644, nil)
	if err == nil {
		t.Fatal("ReplaceFile() succeeded against a destination held open with no sharing")
	}
	if !strings.Contains(err.Error(), path) {
		t.Errorf("ReplaceFile() error = %v, want the destination path named", err)
	}
	if !strings.Contains(err.Error(), "still in use") {
		t.Errorf("ReplaceFile() error = %v, want the underlying condition named", err)
	}
	if !errors.Is(err, errorSharingViolation) {
		t.Errorf("ReplaceFile() error = %v, want it to wrap ERROR_SHARING_VIOLATION", err)
	}

	release()
	if got := readFile(t, path); got != originalContent {
		t.Errorf("destination = %q, want the original content", got)
	}
	if left := tempFiles(t, dir); len(left) != 0 {
		t.Errorf("temporary files left behind: %v", left)
	}
}

// TestReplaceFile_readOnlyDestinationIsNotRetried proves the retry policy does
// not widen past the one transient condition the experiment observed: a
// read-only destination fails with ERROR_ACCESS_DENIED on the first attempt
// and is returned immediately rather than waited on.
func TestReplaceFile_readOnlyDestinationIsNotRetried(t *testing.T) {
	dir := t.TempDir()
	path := writeDestination(t, dir, originalContent)
	setFileAttributes(t, path, syscall.FILE_ATTRIBUTE_READONLY)
	defer setFileAttributes(t, path, syscall.FILE_ATTRIBUTE_NORMAL)

	started := time.Now()
	err := ReplaceFile(path, []byte("replacement bytes\n"), 0644, nil)
	elapsed := time.Since(started)

	if err == nil {
		t.Fatal("ReplaceFile() succeeded on a read-only destination")
	}
	if !errors.Is(err, errorAccessDenied) {
		t.Errorf("ReplaceFile() error = %v, want it to wrap ERROR_ACCESS_DENIED", err)
	}
	if elapsed >= replaceRetryDelay {
		t.Errorf("ReplaceFile() took %s, want no retry delay for a non-transient failure", elapsed)
	}
	if got := readFile(t, path); got != originalContent {
		t.Errorf("destination = %q, want the original content", got)
	}
}
