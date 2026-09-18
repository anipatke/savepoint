package migrate

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Distinct confinement diagnostics. Each names exactly one way a discovered
// path could stop being safely project-relative, so a caller (and a test) can
// tell them apart with errors.Is rather than string-matching a message.
var (
	ErrPathEscapesRoot   = errors.New("migrate: path resolves outside the project root")
	ErrSymlinkNotAllowed = errors.New("migrate: symlink is not a valid source file")
	ErrCaseCollision     = errors.New("migrate: paths collide on a case-insensitive filesystem")
)

// migrationStateDir is the migration operation's own working state, not
// project source; Inventory never walks into it.
const migrationStateDir = ".migration"

// archiveDirName is the byte-preserved archive migration itself creates
// under .savepoint/ (see archivePathFor in plan.go, the single source of
// truth this name is shared with). Inventory never walks into it either: a
// resumed apply calls Plan again, and without this exclusion a prior partial
// apply's already-written archive content would be re-inventoried as a
// brand-new, unclassified V1 source on every resume.
const archiveDirName = "archive"

// SourceFile is one project-owned source file, recorded from its exact bytes
// rather than from any parsed and re-marshalled model: a hash taken from a
// loader's healed re-encoding would record what the loader wished the user
// had written, not what they actually wrote. Path is project-relative and
// forward-slash normalized, so a Windows target produces the same inventory
// keys as a Unix one.
type SourceFile struct {
	Path   string
	SHA256 string
	Size   int64
	Mode   fs.FileMode
}

// Inventory walks the project-owned source under projectRoot — .savepoint/
// content (excluding its own .savepoint/.migration/ operation state), the
// managed AGENTS.md, and agent-skills/ — and returns one SourceFile per file
// found, hashed from its exact bytes with no write of any kind. AGENTS.md and
// agent-skills/ are each optional; their absence is expected V1 shape and
// produces no error. Files are returned sorted by Path, so every consumer
// downstream of Inventory is deterministic for free.
func Inventory(projectRoot string) ([]SourceFile, error) {
	rootAbs, err := filepath.Abs(projectRoot)
	if err != nil {
		return nil, err
	}
	confiner := newPathConfiner(rootAbs)

	var files []SourceFile
	addFile := func(absPath, relPath string) error {
		info, err := confiner.confine(absPath, relPath)
		if err != nil {
			return err
		}
		content, err := os.ReadFile(absPath)
		if err != nil {
			return fmt.Errorf("read %s: %w", relPath, err)
		}
		sum := sha256.Sum256(content)
		files = append(files, SourceFile{
			Path:   filepath.ToSlash(relPath),
			SHA256: hex.EncodeToString(sum[:]),
			Size:   info.Size(),
			Mode:   info.Mode(),
		})
		return nil
	}

	agentsPath := filepath.Join(rootAbs, "AGENTS.md")
	switch _, statErr := os.Lstat(agentsPath); {
	case statErr == nil:
		if err := addFile(agentsPath, "AGENTS.md"); err != nil {
			return nil, err
		}
	case !os.IsNotExist(statErr):
		return nil, statErr
	}

	if err := walkTree(rootAbs, ".savepoint", addFile); err != nil {
		return nil, err
	}
	if err := walkTree(rootAbs, "agent-skills", addFile); err != nil {
		return nil, err
	}

	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	return files, nil
}

// walkTree walks rootAbs/relRoot, calling addFile with the absolute and
// project-relative path of every regular file found. A missing relRoot is not
// an error, since agent-skills/ in particular may not exist in a given V1
// project. .savepoint/.migration/ is skipped entirely — it is the migration
// operation's own state, never project source.
func walkTree(rootAbs, relRoot string, addFile func(absPath, relPath string) error) error {
	absRoot := filepath.Join(rootAbs, relRoot)
	switch _, err := os.Lstat(absRoot); {
	case os.IsNotExist(err):
		return nil
	case err != nil:
		return err
	}

	skipDirs := map[string]bool{
		filepath.Join(".savepoint", migrationStateDir): true,
		filepath.Join(".savepoint", archiveDirName):    true,
	}

	return filepath.WalkDir(absRoot, func(absPath string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(rootAbs, absPath)
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if skipDirs[relPath] {
				return filepath.SkipDir
			}
			return nil
		}
		return addFile(absPath, relPath)
	})
}

// pathConfiner rejects the three ways a discovered path could stop being a
// safe, project-relative source file: a symlink of any kind (its target is
// never trusted, escaping or not), a path that resolves outside the project
// root, and a path that aliases one already confined on a case-insensitive
// filesystem. Every file Inventory records passes through confine, so none of
// the three checks can be bypassed by an entry that is itself a symlink.
type pathConfiner struct {
	rootAbs string
	seen    map[string]string // lowercased relative path -> original relative path
}

func newPathConfiner(rootAbs string) *pathConfiner {
	return &pathConfiner{rootAbs: rootAbs, seen: map[string]string{}}
}

// confine validates absPath (identified for diagnostics by project-relative
// relPath) and returns its FileInfo once confined.
func (c *pathConfiner) confine(absPath, relPath string) (os.FileInfo, error) {
	info, err := os.Lstat(absPath)
	if err != nil {
		return nil, err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("%w: %s", ErrSymlinkNotAllowed, relPath)
	}

	absClean, err := filepath.Abs(absPath)
	if err != nil {
		return nil, err
	}
	if absClean != c.rootAbs && !strings.HasPrefix(absClean, c.rootAbs+string(filepath.Separator)) {
		return nil, fmt.Errorf("%w: %s", ErrPathEscapesRoot, relPath)
	}

	key := strings.ToLower(relPath)
	if existing, ok := c.seen[key]; ok && existing != relPath {
		return nil, fmt.Errorf("%w: %s and %s", ErrCaseCollision, existing, relPath)
	}
	c.seen[key] = relPath

	return info, nil
}
