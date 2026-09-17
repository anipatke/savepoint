//go:build windows

package migrate

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"
)

// ReplaceFileW is bound here rather than taken from a library because neither
// the standard library's syscall package nor golang.org/x/sys/windows exposes
// it. Binding it through syscall.NewLazyDLL is the whole cost of the call, so
// the standard library sufficed and no new direct dependency was added.
var (
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	procReplaceFileW = kernel32.NewProc("ReplaceFileW")
)

// errorSharingViolation is the only transient replacement failure the Windows
// experiment for this contract observed: ReplaceFileW returns it while another
// process holds the destination open with a restrictive share mode, and the
// same call succeeds unchanged once that process lets go. ERROR_ACCESS_DENIED
// is deliberately absent — it was observed only from a destination carrying
// the read-only attribute, which is a standing condition that no amount of
// waiting clears.
const errorSharingViolation = syscall.Errno(32)

// replaceDestination replaces path with the completed temporary file at
// tempPath using ReplaceFileW, which preserves the destination's attributes
// and ACLs across the replacement. MoveFileEx — what os.Rename compiles to on
// Windows — was rejected for this path: it discards both, and it is refused
// with ERROR_ACCESS_DENIED whenever any other process holds the destination
// open at all, even under the permissive share mode Go's own os.Open uses.
// ReplaceFileW succeeds in that case, and reports genuine contention as a
// distinguishable ERROR_SHARING_VIOLATION rather than folding it into the same
// code a permission failure uses.
func replaceDestination(tempPath, path string) error {
	return replaceFileW(path, tempPath, "", 0)
}

// isTransientReplaceError reports whether a failed replacement could succeed
// on a retry, which on Windows means exactly one observed condition.
func isTransientReplaceError(err error) bool {
	var errno syscall.Errno
	return errors.As(err, &errno) && errno == errorSharingViolation
}

// replaceFileW calls ReplaceFileW(replaced, replacement, backup, flags). An
// empty backup passes NULL, which asks Windows not to keep a copy of the
// replaced file — migration keeps its own verified backup in the operation
// directory instead, where it survives the process rather than the call.
func replaceFileW(replaced, replacement, backup string, flags uint32) error {
	if err := procReplaceFileW.Find(); err != nil {
		return fmt.Errorf("ReplaceFileW unavailable: %w", err)
	}
	replacedPtr, err := syscall.UTF16PtrFromString(replaced)
	if err != nil {
		return err
	}
	replacementPtr, err := syscall.UTF16PtrFromString(replacement)
	if err != nil {
		return err
	}
	var backupPtr *uint16
	if backup != "" {
		if backupPtr, err = syscall.UTF16PtrFromString(backup); err != nil {
			return err
		}
	}

	r1, _, e1 := syscall.SyscallN(procReplaceFileW.Addr(),
		uintptr(unsafe.Pointer(replacedPtr)),
		uintptr(unsafe.Pointer(replacementPtr)),
		uintptr(unsafe.Pointer(backupPtr)),
		uintptr(flags), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			return e1
		}
		return syscall.EINVAL
	}
	return nil
}
