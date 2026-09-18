//go:build !windows

package migrate

import "os"

// replaceDestination replaces path with the completed temporary file at
// tempPath. Within a single directory rename is atomic on Unix and is not
// refused because another process holds the destination open — an open handle
// keeps referring to the replaced inode — so the primitive needs nothing
// around it.
func replaceDestination(tempPath, path string) error {
	return os.Rename(tempPath, path)
}

// createDestinationAtomically links the complete temporary inode into the
// destination name without replacement. POSIX link(2) is atomic with respect
// to a competing creator: if the destination appeared after planning, the
// link fails with EEXIST and the existing inode is left untouched.
func createDestinationAtomically(tempPath, path string) error {
	return os.Link(tempPath, path)
}

func isCreateDestinationExistsError(err error) bool {
	return os.IsExist(err)
}

// isTransientReplaceError reports whether a failed replacement could succeed
// on a retry. On Unix it never does: the experiment behind this contract found
// no condition where rename failed for contention and then succeeded
// unchanged, because an open destination does not block rename in the first
// place. Every failure here is a standing condition — a missing directory, a
// permission denial, a cross-device path — so retrying would only delay a
// truthful error.
func isTransientReplaceError(error) bool {
	return false
}
