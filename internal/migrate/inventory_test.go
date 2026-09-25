package migrate

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/opencode/savepoint/internal/testutil"
)

// TestInventory_matchesFixtureManifestHashesAndFiles proves, for each frozen
// fixture, that Inventory reports exactly the files manifest.yml records, each
// with the independently authored hash — the check that actually catches a
// hash taken from something other than the file's own bytes.
func TestInventory_matchesFixtureManifestHashesAndFiles(t *testing.T) {
	for _, fixture := range []string{"v1-basic", "v1-history"} {
		t.Run(fixture, func(t *testing.T) {
			manifest := loadFixtureManifest(t, fixture)
			files := mustInventory(t, fixture)

			got := map[string]SourceFile{}
			for _, f := range files {
				got[f.Path] = f
			}

			want := map[string]bool{}
			for _, mf := range manifest.Files {
				relPath := fixtureRelPath(mf.Path)
				want[relPath] = true

				sf, ok := got[relPath]
				if !ok {
					t.Errorf("Inventory(%s) missing %s, which manifest.yml records", fixture, relPath)
					continue
				}
				if sf.SHA256 != mf.SHA256 {
					t.Errorf("Inventory(%s) %s SHA256 = %s, want %s (from manifest.yml)", fixture, relPath, sf.SHA256, mf.SHA256)
				}
			}

			for path := range got {
				if !want[path] {
					t.Errorf("Inventory(%s) reports %s, which manifest.yml does not record", fixture, path)
				}
			}

			for _, a := range manifest.AbsentByDesign {
				if _, ok := got[fixtureRelPath(a.Path)]; ok {
					t.Errorf("Inventory(%s) reports %s, which manifest.yml marks absent by design (%s)", fixture, a.Path, a.Reason)
				}
			}
		})
	}
}

// TestInventory_hashesRawBytesNotAHealedReparse targets v1-basic's
// T002-follow-up.md, whose raw `phase: implementation` the V1 loader heals to
// `stage: build` on parse. Inventory never parses the file at all, so its
// hash must match the manifest's independently authored byte hash rather than
// whatever a re-marshalled, healed record would hash to.
func TestInventory_hashesRawBytesNotAHealedReparse(t *testing.T) {
	const fixture = "v1-basic"
	const path = ".savepoint/releases/v1/epics/E01-example/tasks/T002-follow-up.md"

	manifest := loadFixtureManifest(t, fixture)
	var want string
	for _, f := range manifest.Files {
		if fixtureRelPath(f.Path) == path {
			want = f.SHA256
		}
	}
	if want == "" {
		t.Fatalf("test setup: %s not found in %s manifest.yml", path, fixture)
	}

	sf := findSourceFile(t, mustInventory(t, fixture), path)
	if sf.SHA256 != want {
		t.Errorf("SHA256 = %s, want %s (the file's own authored bytes, not a healed reparse)", sf.SHA256, want)
	}
}

// TestInventory_excludesLegacyMigrationDirectory proves old .savepoint/.migration/
// journals are not reinterpreted as project source, while sibling files are.
func TestInventory_excludesLegacyMigrationDirectory(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".savepoint", "router.md"), "# Router\n")
	writeFile(t, filepath.Join(root, ".savepoint", legacyMigrationStateDir, "journal.yml"), "state: probing\n")

	files, err := Inventory(root)
	if err != nil {
		t.Fatalf("Inventory() error = %v", err)
	}

	for _, f := range files {
		if f.Path == ".savepoint/router.md" {
			continue
		}
		if filepath.Dir(f.Path) == filepath.ToSlash(filepath.Join(".savepoint", legacyMigrationStateDir)) {
			t.Errorf("Inventory() reported %s under .savepoint/%s, want it excluded", f.Path, legacyMigrationStateDir)
		}
	}
	if len(files) != 1 {
		t.Errorf("Inventory() = %v, want exactly [.savepoint/router.md]", files)
	}
}

// TestInventory_excludesArchiveDir is regression coverage for an audit
// finding: once the byte-preserved archive moved to .savepoint/archive/v1/
// (see archivePathFor in plan.go), a resumed migration's second Plan call
// would otherwise re-inventory a prior partial apply's already-written
// archive content as a brand-new, unclassified V1 source. Inventory must
// exclude .savepoint/archive/ exactly as it already excludes
// .savepoint/.migration/.
func TestInventory_excludesArchiveDir(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".savepoint", "router.md"), "# Router\n")
	writeFile(t, filepath.Join(root, ".savepoint", archiveDirName, "v1", ".savepoint", "PRD.md"), "archived content\n")

	files, err := Inventory(root)
	if err != nil {
		t.Fatalf("Inventory() error = %v", err)
	}

	archivePrefix := filepath.ToSlash(filepath.Join(".savepoint", archiveDirName)) + "/"
	for _, f := range files {
		if strings.HasPrefix(f.Path, archivePrefix) {
			t.Errorf("Inventory() reported %s under .savepoint/%s, want it excluded", f.Path, archiveDirName)
		}
	}
	if len(files) != 1 {
		t.Errorf("Inventory() = %v, want exactly [.savepoint/router.md]", files)
	}
}

// TestInventory_unclassifiedFileIsInventoriedNotDropped proves a file
// matching no known role still comes back from Inventory — a role-vocabulary
// gap must be visible in the report, never a silently shrunk file list.
func TestInventory_unclassifiedFileIsInventoriedNotDropped(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".savepoint", "scratch-notes.txt"), "not a known shape\n")

	files, err := Inventory(root)
	if err != nil {
		t.Fatalf("Inventory() error = %v", err)
	}

	sf := findSourceFile(t, files, ".savepoint/scratch-notes.txt")
	if role := Classify(sf.Path); role != RoleUnclassified {
		t.Errorf("Classify(%s) = %s, want %s", sf.Path, role, RoleUnclassified)
	}
}

// TestInventory_deterministicOrder proves repeated calls return files in the
// same, path-sorted order.
func TestInventory_deterministicOrder(t *testing.T) {
	for _, fixture := range []string{"v1-basic", "v1-history"} {
		files := mustInventory(t, fixture)
		if !sort.SliceIsSorted(files, func(i, j int) bool { return files[i].Path < files[j].Path }) {
			t.Errorf("Inventory(%s) is not sorted by Path: %v", fixture, files)
		}

		again := mustInventory(t, fixture)
		if len(again) != len(files) {
			t.Fatalf("Inventory(%s) returned %d files, then %d files", fixture, len(files), len(again))
		}
		for i := range files {
			if files[i].Path != again[i].Path {
				t.Errorf("Inventory(%s)[%d] = %s on first call, %s on second", fixture, i, files[i].Path, again[i].Path)
			}
		}
	}
}

// TestInventory_performsNoWrite snapshots every fixture file's bytes and
// modification time, runs Inventory (and Classify over its result), and
// asserts nothing changed — on the success path.
func TestInventory_performsNoWrite(t *testing.T) {
	for _, fixture := range []string{"v1-basic", "v1-history"} {
		t.Run(fixture, func(t *testing.T) {
			before := snapshotTree(t, fixtureProjectRoot(fixture))

			files := mustInventory(t, fixture)
			for _, f := range files {
				_ = Classify(f.Path)
			}

			after := snapshotTree(t, fixtureProjectRoot(fixture))
			assertSnapshotsEqual(t, before, after)
		})
	}
}

// TestInventory_performsNoWriteOnFailure repeats the same snapshot proof
// around each named failure path: a symlink, a case collision, and a path
// that resolves outside the project root.
func TestInventory_performsNoWriteOnFailure(t *testing.T) {
	t.Run("symlink present", func(t *testing.T) {
		root := newSymlinkFixture(t)
		before := snapshotTree(t, root)
		if _, err := Inventory(root); !errors.Is(err, ErrSymlinkNotAllowed) {
			t.Fatalf("Inventory() error = %v, want ErrSymlinkNotAllowed", err)
		}
		assertSnapshotsEqual(t, before, snapshotTree(t, root))
	})

	t.Run("case collision", func(t *testing.T) {
		root := newCaseCollisionFixture(t)
		before := snapshotTree(t, root)
		if _, err := Inventory(root); !errors.Is(err, ErrCaseCollision) {
			t.Fatalf("Inventory() error = %v, want ErrCaseCollision", err)
		}
		assertSnapshotsEqual(t, before, snapshotTree(t, root))
	})
}

// TestInventory_rejectsSymlink proves a symlink anywhere under the walked
// trees fails the whole inventory closed, named distinctly from the other two
// confinement diagnostics.
func TestInventory_rejectsSymlink(t *testing.T) {
	root := newSymlinkFixture(t)
	if _, err := Inventory(root); !errors.Is(err, ErrSymlinkNotAllowed) {
		t.Fatalf("Inventory() error = %v, want ErrSymlinkNotAllowed", err)
	}
}

// TestInventory_rejectsCaseCollision proves two paths that would alias on a
// case-insensitive filesystem fail closed rather than silently picking one.
func TestInventory_rejectsCaseCollision(t *testing.T) {
	root := newCaseCollisionFixture(t)
	if _, err := Inventory(root); !errors.Is(err, ErrCaseCollision) {
		t.Fatalf("Inventory() error = %v, want ErrCaseCollision", err)
	}
}

// TestPathConfiner_rejectsEscapeOutsideRoot exercises the third confinement
// diagnostic directly: a resolved path outside the project root is a
// condition Inventory's own walk cannot produce (WalkDir never yields a path
// outside where it started), so it is proven at the confiner directly, the
// same way an escaping symlink would trip it in a live walk.
func TestPathConfiner_rejectsEscapeOutsideRoot(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	escaped := filepath.Join(outside, "escaped.md")
	writeFile(t, escaped, "# Escaped\n")
	confiner := newPathConfiner(root)

	_, err := confiner.confine(escaped, "escaped.md")
	if !errors.Is(err, ErrPathEscapesRoot) {
		t.Fatalf("confine() error = %v, want ErrPathEscapesRoot", err)
	}
}

// TestPathConfiner_rejectsCaseCollision exercises the case-collision
// diagnostic directly, independent of a filesystem walk.
func TestPathConfiner_rejectsCaseCollision(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "Foo.md"), "# Foo\n")
	writeFile(t, filepath.Join(root, "foo.md"), "# foo\n")
	confiner := newPathConfiner(root)

	if _, err := confiner.confine(filepath.Join(root, "Foo.md"), "Foo.md"); err != nil {
		t.Fatalf("confine() first call error = %v", err)
	}
	_, err := confiner.confine(filepath.Join(root, "foo.md"), "foo.md")
	if !errors.Is(err, ErrCaseCollision) {
		t.Fatalf("confine() error = %v, want ErrCaseCollision", err)
	}
}

// --- Fixture builders ----------------------------------------------------

func newSymlinkFixture(t *testing.T) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation requires elevated privileges on windows")
	}

	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".savepoint", "config.yml"), "schema: v1\n")

	outside := t.TempDir()
	target := filepath.Join(outside, "router.md")
	writeFile(t, target, "# Router\n")

	link := filepath.Join(root, ".savepoint", "router.md")
	if err := os.Symlink(target, link); err != nil {
		t.Fatalf("os.Symlink() error = %v", err)
	}
	return root
}

func newCaseCollisionFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	testutil.SkipIfCaseInsensitive(t, root)
	writeFile(t, filepath.Join(root, ".savepoint", "Notes.md"), "# Notes\n")
	writeFile(t, filepath.Join(root, ".savepoint", "notes.md"), "# notes\n")
	return root
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("MkdirAll(%s) error = %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile(%s) error = %v", path, err)
	}
}

// --- Snapshot helpers ------------------------------------------------------

type treeSnapshot map[string]fileSnapshot

type fileSnapshot struct {
	content []byte
	modTime time.Time
}

func snapshotTree(t *testing.T, root string) treeSnapshot {
	t.Helper()
	snap := treeSnapshot{}
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		snap[rel] = fileSnapshot{content: content, modTime: info.ModTime()}
		return nil
	})
	if err != nil {
		t.Fatalf("snapshotTree(%s) error = %v", root, err)
	}
	return snap
}

func assertSnapshotsEqual(t *testing.T, before, after treeSnapshot) {
	t.Helper()
	if len(before) != len(after) {
		t.Fatalf("file count changed: before %d, after %d", len(before), len(after))
	}
	for path, want := range before {
		got, ok := after[path]
		if !ok {
			t.Errorf("%s was present before and missing after", path)
			continue
		}
		if !want.modTime.Equal(got.modTime) {
			t.Errorf("%s modTime changed: before %v, after %v", path, want.modTime, got.modTime)
		}
		if string(want.content) != string(got.content) {
			t.Errorf("%s content changed", path)
		}
	}
}
