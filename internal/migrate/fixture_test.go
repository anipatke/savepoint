package migrate

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

// --- Migration fixture helpers -----------------------------------------
//
// internal/data/testdata/migration/ holds frozen V1 projects kept as
// migration source evidence (see that directory's README.md). Each fixture
// carries a hand-authored manifest.yml recording every source file's exact
// role and SHA-256; these helpers read that manifest as the single source of
// truth for classification and inventory tests, never regenerating a hash to
// compare it against itself.

const migrationFixtureRoot = "../data/testdata/migration"

type fixtureManifest struct {
	Fixture        string                  `yaml:"fixture"`
	Files          []fixtureManifestFile   `yaml:"files"`
	AbsentByDesign []fixtureManifestAbsent `yaml:"absent_by_design"`
}

type fixtureManifestFile struct {
	Path   string `yaml:"path"`
	SHA256 string `yaml:"sha256"`
	Role   string `yaml:"role"`
}

type fixtureManifestAbsent struct {
	Path   string `yaml:"path"`
	Reason string `yaml:"reason"`
}

func fixtureDir(fixture string) string {
	return filepath.Join(migrationFixtureRoot, fixture)
}

func fixtureProjectRoot(fixture string) string {
	return filepath.Join(fixtureDir(fixture), "project")
}

func loadFixtureManifest(t *testing.T, fixture string) fixtureManifest {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(fixtureDir(fixture), "manifest.yml"))
	if err != nil {
		t.Fatalf("read manifest for fixture %s: %v", fixture, err)
	}
	var manifest fixtureManifest
	if err := yaml.Unmarshal(raw, &manifest); err != nil {
		t.Fatalf("parse manifest for fixture %s: %v", fixture, err)
	}
	return manifest
}

// fixtureRelPath turns a manifest path ("project/AGENTS.md") into the
// project-relative form ("AGENTS.md") Inventory and Classify use.
func fixtureRelPath(manifestPath string) string {
	rel, err := filepath.Rel("project", manifestPath)
	if err != nil {
		return manifestPath
	}
	return filepath.ToSlash(rel)
}

func mustInventory(t *testing.T, fixture string) []SourceFile {
	t.Helper()
	files, err := Inventory(fixtureProjectRoot(fixture))
	if err != nil {
		t.Fatalf("Inventory(%s) error = %v", fixture, err)
	}
	return files
}

func findSourceFile(t *testing.T, files []SourceFile, path string) SourceFile {
	t.Helper()
	for _, f := range files {
		if f.Path == path {
			return f
		}
	}
	t.Fatalf("no SourceFile with path %q among %d inventoried files", path, len(files))
	return SourceFile{}
}
