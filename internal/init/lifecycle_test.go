package init

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Coverage for the E40 Template and skill lifecycle matrix, run against the
// real shipped templates rather than a hand-built fixture set: a new asset must
// land in the right ownership class the moment it ships, not the next time
// someone remembers to extend a test table.

// ownership is what the matrix promises about a shipped path on upgrade.
type ownership string

const (
	// ownPackage: Savepoint refreshes the file in place, and a customized one
	// conflicts rather than being overwritten.
	ownPackage ownership = "package-owned"
	// ownInstallIfMissing: written only when absent, never changed after.
	ownInstallIfMissing ownership = "install-if-missing"
	// ownProject: the user's file forever, `--force` included.
	ownProject ownership = "project-owned"
)

// wantOwnership is the matrix as a test expectation, kept independent of the
// production predicates so a change to those has to be a deliberate change
// here too.
func wantOwnership(path string) ownership {
	switch {
	case path == "AGENTS.md", strings.HasPrefix(path, "agent-skills/"):
		return ownPackage
	case path == ".savepoint/Guardrails.md", path == ".savepoint/Health-Check.md":
		return ownInstallIfMissing
	case strings.HasPrefix(path, ".savepoint/audit/"):
		return ownInstallIfMissing
	default:
		return ownProject
	}
}

// shippedTemplates is the real templates/project tree, the same one main.go
// embeds.
func shippedTemplates(t *testing.T) fs.FS {
	t.Helper()
	return os.DirFS(filepath.Join("..", "..", "templates", "project"))
}

func shippedPaths(t *testing.T, templates fs.FS) []string {
	t.Helper()
	var paths []string
	err := fs.WalkDir(templates, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk shipped templates: %v", err)
	}
	if len(paths) == 0 {
		t.Fatal("no shipped templates found")
	}
	return paths
}

// scaffoldProject installs the shipped templates into a fresh project.
func scaffoldProject(t *testing.T, templates fs.FS) string {
	t.Helper()
	dir := t.TempDir()
	if err := Scaffold(templates, dir, ProjectNameFromDir(dir), false); err != nil {
		t.Fatalf("Scaffold() error = %v", err)
	}
	return dir
}

func readProjectFile(t *testing.T, dir, path string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(path)))
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

// editEveryFile appends a local edit to every installed file and returns the
// resulting content per path, so an upgrade can be checked against what the
// user actually has on disk.
func editEveryFile(t *testing.T, dir string, paths []string) map[string]string {
	t.Helper()
	edited := make(map[string]string, len(paths))
	for _, path := range paths {
		content := readProjectFile(t, dir, path) + "\n<!-- local edit -->\n"
		if err := AtomicWrite(filepath.Join(dir, filepath.FromSlash(path)), []byte(content)); err != nil {
			t.Fatalf("edit %s: %v", path, err)
		}
		edited[path] = content
	}
	return edited
}

func TestLifecycle_installWritesEveryShippedAsset(t *testing.T) {
	templates := shippedTemplates(t)
	paths := shippedPaths(t, templates)
	dir := scaffoldProject(t, templates)

	for _, path := range paths {
		if _, err := os.Stat(filepath.Join(dir, filepath.FromSlash(path))); err != nil {
			t.Errorf("shipped asset %s missing after init: %v", path, err)
		}
	}

	// Fresh projects start with exact provenance for every skill, and for
	// nothing else: the manifest's scope is the skill entrypoints.
	manifest, err := LoadManifest(dir)
	if err != nil {
		t.Fatalf("LoadManifest() error = %v", err)
	}
	recorded := 0
	for _, path := range paths {
		hash, tracked := manifest.Hash(path)
		if !isManifestPath(path) {
			if tracked {
				t.Errorf("manifest tracks out-of-scope path %s", path)
			}
			continue
		}
		recorded++
		if !tracked {
			t.Errorf("manifest has no entry for %s", path)
			continue
		}
		if want := hashContent([]byte(readProjectFile(t, dir, path))); hash != want {
			t.Errorf("manifest hash for %s does not match the installed file", path)
		}
	}
	if recorded == 0 {
		t.Fatal("no skills recorded; the manifest scope check proved nothing")
	}
}

// treeSnapshot records every path in a project with the metadata a write would
// disturb — directories included, since creating and removing a probe file
// changes its parent's modification time without changing any file.
func treeSnapshot(t *testing.T, dir string) map[string]string {
	t.Helper()
	snapshot := map[string]string{}
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		snapshot[rel] = fmt.Sprintf("dir=%v size=%d mtime=%d", d.IsDir(), info.Size(), info.ModTime().UnixNano())
		return nil
	})
	if err != nil {
		t.Fatalf("snapshot %s: %v", dir, err)
	}
	return snapshot
}

func TestLifecycle_unchangedUpgradeTouchesNothing(t *testing.T) {
	templates := shippedTemplates(t)
	dir := scaffoldProject(t, templates)

	// The first upgrade settles any difference between scaffolding and upgrade
	// output; the second has genuinely nothing to do.
	if _, err := UpgradeProjectAssets(templates, dir, false, false); err != nil {
		t.Fatalf("first UpgradeProjectAssets() error = %v", err)
	}
	before := treeSnapshot(t, dir)

	report, err := UpgradeProjectAssets(templates, dir, false, false)
	if err != nil {
		t.Fatalf("second UpgradeProjectAssets() error = %v", err)
	}
	for _, e := range report.Actions {
		if e.Action != ActionUnchanged && e.Action != ActionSkipped {
			t.Errorf("second run %s = %v, want unchanged or skipped", e.Path, e.Action)
		}
	}

	after := treeSnapshot(t, dir)
	for path, want := range before {
		got, ok := after[path]
		if !ok {
			t.Errorf("%s disappeared during an unchanged upgrade", path)
			continue
		}
		if got != want {
			t.Errorf("%s changed during an unchanged upgrade: %s -> %s", path, want, got)
		}
	}
	for path := range after {
		if _, ok := before[path]; !ok {
			t.Errorf("%s was created by an unchanged upgrade", path)
		}
	}
}

func TestLifecycle_upgradeHonoursOwnership(t *testing.T) {
	templates := shippedTemplates(t)
	paths := shippedPaths(t, templates)

	for _, force := range []bool{false, true} {
		name := "without force"
		if force {
			name = "with force"
		}
		t.Run(name, func(t *testing.T) {
			dir := scaffoldProject(t, templates)
			edited := editEveryFile(t, dir, paths)

			report, err := UpgradeProjectAssets(templates, dir, false, force)
			if err != nil {
				t.Fatalf("UpgradeProjectAssets() error = %v", err)
			}

			checked := map[ownership]int{}
			for _, path := range paths {
				got := readProjectFile(t, dir, path)
				checked[wantOwnership(path)]++
				switch wantOwnership(path) {
				case ownProject, ownInstallIfMissing:
					// The permanent boundary: an edited project-owned file is
					// never rewritten, and `--force` does not widen the set of
					// files Savepoint is allowed to touch.
					if got != edited[path] {
						t.Errorf("%s (%s) was rewritten by upgrade", path, wantOwnership(path))
					}
					if action, found := actionFor(report, path); found && action != ActionSkipped && action != ActionUnchanged {
						t.Errorf("%s action = %v, want skipped or unchanged", path, action)
					}
				case ownPackage:
					// AGENTS.md is prose plus a managed block, so an edit
					// outside the block survives by design; its refresh is
					// asserted in TestLifecycle_upgradeRefreshesPackageOwnedAssets.
					if path == "AGENTS.md" {
						continue
					}
					if !force && isManifestPath(path) {
						continue // conflicts are asserted below
					}
					if strings.Contains(got, "<!-- local edit -->") {
						t.Errorf("%s kept the local edit, want the shipped content", path)
					}
				}
			}

			// Guard against a vacuous pass if a rename ever empties a class.
			for _, class := range []ownership{ownProject, ownInstallIfMissing, ownPackage} {
				if checked[class] == 0 {
					t.Errorf("no %s paths checked", class)
				}
			}
		})
	}
}

func TestLifecycle_upgradeRefreshesPackageOwnedAssets(t *testing.T) {
	templates := shippedTemplates(t)
	paths := shippedPaths(t, templates)
	dir := scaffoldProject(t, templates)
	editEveryFile(t, dir, paths)
	staleGuide := "# Team Guide\n\nOur prose.\n\n" + managedBegin + "\n# Old managed content\n" + managedEnd + "\n\nMore of our prose.\n"
	if err := AtomicWrite(filepath.Join(dir, "AGENTS.md"), []byte(staleGuide)); err != nil {
		t.Fatal(err)
	}

	report, err := UpgradeProjectAssets(templates, dir, false, true)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	skills := 0
	for _, path := range paths {
		if wantOwnership(path) != ownPackage || path == "AGENTS.md" {
			continue
		}
		skills++

		want, err := fs.ReadFile(templates, path)
		if err != nil {
			t.Fatalf("read template %s: %v", path, err)
		}
		if got := readProjectFile(t, dir, path); got != string(want) {
			t.Errorf("%s was not refreshed to the shipped content", path)
		}
		if action, found := actionFor(report, path); !found || action != ActionUpdated {
			t.Errorf("%s action = %v (found %v), want updated", path, action, found)
		}
		// Nothing the user had is lost: replacing a tracked skill under
		// --force keeps the old copy beside it. Shared references under
		// agent-skills/references/ carry no provenance and refresh directly.
		_, err = os.Stat(filepath.Join(dir, filepath.FromSlash(path)) + backupSuffix)
		if isManifestPath(path) && err != nil {
			t.Errorf("%s replaced without a backup: %v", path, err)
		}
		if !isManifestPath(path) && err == nil {
			t.Errorf("%s is package-owned and needs no backup, but one was written", path)
		}
	}
	if skills == 0 {
		t.Fatal("no package-owned skill assets found")
	}

	// The marked guide keeps its user prose and takes the new managed block.
	guide := readProjectFile(t, dir, "AGENTS.md")
	for _, prose := range []string{"# Team Guide", "Our prose.", "More of our prose."} {
		if !strings.Contains(guide, prose) {
			t.Errorf("AGENTS.md lost user prose %q outside the managed block", prose)
		}
	}
	if strings.Contains(guide, "# Old managed content") {
		t.Error("AGENTS.md kept the stale managed block")
	}
	wantBlock, err := fs.ReadFile(templates, "AGENTS.md")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(guide, strings.TrimSpace(string(wantBlock))) {
		t.Error("AGENTS.md managed block was not refreshed to the shipped guide")
	}
	if strings.Count(guide, managedBegin) != 1 {
		t.Errorf("AGENTS.md has %d managed blocks, want 1", strings.Count(guide, managedBegin))
	}
	if action, found := actionFor(report, "AGENTS.md"); !found || action != ActionMerged {
		t.Errorf("AGENTS.md action = %v (found %v), want merged", action, found)
	}
}

func TestLifecycle_upgradeConflictsOnCustomizedSkills(t *testing.T) {
	templates := shippedTemplates(t)
	paths := shippedPaths(t, templates)
	dir := scaffoldProject(t, templates)
	edited := editEveryFile(t, dir, paths)

	report, err := UpgradeProjectAssets(templates, dir, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	conflicts := 0
	for _, path := range paths {
		if !isManifestPath(path) {
			continue
		}
		conflicts++

		if got := readProjectFile(t, dir, path); got != edited[path] {
			t.Errorf("customized %s was overwritten without --force", path)
		}
		entry, found := entryFor(report, path)
		if !found || entry.Action != ActionConflict {
			t.Errorf("%s entry = %+v, want conflict", path, entry)
		}
		if _, err := os.Stat(filepath.Join(dir, filepath.FromSlash(path)) + incomingSuffix); err != nil {
			t.Errorf("%s has no incoming sidecar: %v", path, err)
		}
	}
	if conflicts == 0 {
		t.Fatal("no manifest-tracked skills found")
	}
}

func TestLifecycle_upgradeReinstallsMissingInstallIfMissingAssets(t *testing.T) {
	templates := shippedTemplates(t)
	paths := shippedPaths(t, templates)
	dir := scaffoldProject(t, templates)

	var removed []string
	for _, path := range paths {
		if wantOwnership(path) != ownInstallIfMissing {
			continue
		}
		if err := os.Remove(filepath.Join(dir, filepath.FromSlash(path))); err != nil {
			t.Fatalf("remove %s: %v", path, err)
		}
		removed = append(removed, path)
	}
	if len(removed) == 0 {
		t.Fatal("no install-if-missing assets found")
	}

	if _, err := UpgradeProjectAssets(templates, dir, false, false); err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	for _, path := range removed {
		if _, err := os.Stat(filepath.Join(dir, filepath.FromSlash(path))); err != nil {
			t.Errorf("install-if-missing asset %s not restored: %v", path, err)
		}
	}
}
