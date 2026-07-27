package init

import (
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

func manifestPath(dir string) string {
	return filepath.Join(dir, ".savepoint", ".upgrade-manifest.yml")
}

func writeManifestFile(t *testing.T, dir, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(dir, ".savepoint"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(manifestPath(dir), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestManifest_roundTrip(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".savepoint"), 0755); err != nil {
		t.Fatal(err)
	}

	saved := NewManifest()
	saved.Record("agent-skills/savepoint-build-task/SKILL.md", []byte("# Build Task"))
	if err := saved.Save(dir); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loaded, err := LoadManifest(dir)
	if err != nil {
		t.Fatalf("LoadManifest() error = %v", err)
	}
	if loaded.Version != ManifestVersion {
		t.Errorf("Version = %d, want %d", loaded.Version, ManifestVersion)
	}

	sum := sha256.Sum256([]byte("# Build Task"))
	want := hex.EncodeToString(sum[:])
	got, ok := loaded.Hash("agent-skills/savepoint-build-task/SKILL.md")
	if !ok {
		t.Fatal("skill hash not recorded")
	}
	if got != want {
		t.Errorf("hash = %q, want %q", got, want)
	}
}

func TestManifest_saveOfIdenticalContentDoesNotWrite(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".savepoint"), 0755); err != nil {
		t.Fatal(err)
	}

	manifest := NewManifest()
	manifest.Record("agent-skills/savepoint-build-task/SKILL.md", []byte("# Build Task"))
	if err := manifest.Save(dir); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	before, err := os.Stat(manifestPath(dir))
	if err != nil {
		t.Fatal(err)
	}

	// Saving unchanged provenance is not a change, and must not look like one
	// on disk.
	if err := manifest.Save(dir); err != nil {
		t.Fatalf("second Save() error = %v", err)
	}

	after, err := os.Stat(manifestPath(dir))
	if err != nil {
		t.Fatal(err)
	}
	if !os.SameFile(before, after) || !before.ModTime().Equal(after.ModTime()) {
		t.Error("Save() replaced the manifest despite identical content")
	}
}

func TestLoadManifest_missingFileIsEmpty(t *testing.T) {
	dir := t.TempDir()

	manifest, err := LoadManifest(dir)
	if err != nil {
		t.Fatalf("LoadManifest() error = %v", err)
	}
	if len(manifest.Skills) != 0 {
		t.Errorf("Skills = %v, want empty", manifest.Skills)
	}
	if manifest.Version != ManifestVersion {
		t.Errorf("Version = %d, want %d", manifest.Version, ManifestVersion)
	}
}

func TestLoadManifest_malformedFileNamesTheFile(t *testing.T) {
	dir := t.TempDir()
	writeManifestFile(t, dir, "version: 1\nskills: [not, a, map\n")

	_, err := LoadManifest(dir)
	if err == nil {
		t.Fatal("expected error for malformed manifest")
	}
	if !strings.Contains(err.Error(), ".upgrade-manifest.yml") {
		t.Errorf("error = %q, want it to name the manifest file", err.Error())
	}
}

func TestManifest_recordNormalizesToForwardSlashes(t *testing.T) {
	manifest := NewManifest()
	manifest.Record(filepath.Join("agent-skills", "savepoint-audit-task", "SKILL.md"), []byte("x"))

	if _, ok := manifest.Skills["agent-skills/savepoint-audit-task/SKILL.md"]; !ok {
		t.Errorf("Skills = %v, want a forward-slash key", manifest.Skills)
	}
}

func TestManifest_recordIgnoresPathsOutsideScope(t *testing.T) {
	manifest := NewManifest()
	manifest.Record("AGENTS.md", []byte("x"))
	manifest.Record("agent-skills/references/audit-method.md", []byte("x"))
	manifest.Record(".savepoint/router.md", []byte("x"))
	manifest.Record("agent-skills/savepoint-build-task/nested/SKILL.md", []byte("x"))

	if len(manifest.Skills) != 0 {
		t.Errorf("Skills = %v, want empty", manifest.Skills)
	}
}

func TestScaffold_writesManifestForSkills(t *testing.T) {
	dir := t.TempDir()
	templates := fstest.MapFS{
		"AGENTS.md": &fstest.MapFile{Data: []byte("# Guide for {{PROJECT_NAME}}")},
		"agent-skills/savepoint-build-task/SKILL.md": &fstest.MapFile{Data: []byte("# Build {{PROJECT_NAME}}")},
		"agent-skills/references/audit-method.md":    &fstest.MapFile{Data: []byte("# Method")},
	}

	if err := Scaffold(templates, dir, "demo", false); err != nil {
		t.Fatalf("Scaffold() error = %v", err)
	}

	manifest, err := LoadManifest(dir)
	if err != nil {
		t.Fatalf("LoadManifest() error = %v", err)
	}
	if len(manifest.Skills) != 1 {
		t.Fatalf("Skills = %v, want exactly the skill entrypoint", manifest.Skills)
	}

	// The hash must match the interpolated bytes actually written to disk.
	written, err := os.ReadFile(filepath.Join(dir, "agent-skills", "savepoint-build-task", "SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(written)
	want := hex.EncodeToString(sum[:])
	got, _ := manifest.Hash("agent-skills/savepoint-build-task/SKILL.md")
	if got != want {
		t.Errorf("hash = %q, want %q for on-disk content", got, want)
	}
}

func TestUpgradeProjectAssets_recordsSkillHashes(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".savepoint"), 0755); err != nil {
		t.Fatal(err)
	}
	templates := fstest.MapFS{
		"agent-skills":                               &fstest.MapFile{Mode: fs.ModeDir | 0755},
		"agent-skills/savepoint-build-task/SKILL.md": &fstest.MapFile{Data: []byte("# Build Task")},
	}

	if _, err := UpgradeProjectAssets(templates, dir, false, false); err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	manifest, err := LoadManifest(dir)
	if err != nil {
		t.Fatalf("LoadManifest() error = %v", err)
	}
	sum := sha256.Sum256([]byte("# Build Task"))
	want := hex.EncodeToString(sum[:])
	if got, _ := manifest.Hash("agent-skills/savepoint-build-task/SKILL.md"); got != want {
		t.Errorf("hash = %q, want %q", got, want)
	}
}

func TestUpgradeProjectAssets_dryRunLeavesManifestUntouched(t *testing.T) {
	dir := t.TempDir()
	writeManifestFile(t, dir, "version: 1\nskills: {}\n")
	before, err := os.ReadFile(manifestPath(dir))
	if err != nil {
		t.Fatal(err)
	}

	templates := fstest.MapFS{
		"agent-skills":                               &fstest.MapFile{Mode: fs.ModeDir | 0755},
		"agent-skills/savepoint-build-task/SKILL.md": &fstest.MapFile{Data: []byte("# Build Task")},
	}

	if _, err := UpgradeProjectAssets(templates, dir, true, false); err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	after, err := os.ReadFile(manifestPath(dir))
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(before) {
		t.Errorf("manifest changed during dry run:\nbefore %q\nafter  %q", before, after)
	}
}

func TestUpgradeProjectAssets_dryRunDoesNotCreateManifest(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".savepoint"), 0755); err != nil {
		t.Fatal(err)
	}
	templates := fstest.MapFS{
		"agent-skills":                               &fstest.MapFile{Mode: fs.ModeDir | 0755},
		"agent-skills/savepoint-build-task/SKILL.md": &fstest.MapFile{Data: []byte("# Build Task")},
	}

	if _, err := UpgradeProjectAssets(templates, dir, true, false); err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	if _, err := os.Stat(manifestPath(dir)); !os.IsNotExist(err) {
		t.Errorf("manifest exists after dry run (err = %v), want it absent", err)
	}
}
