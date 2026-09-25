package migrate

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

// TestLiveSourcesDoNotReachV1Readers is the package-boundary lock for the
// cutover. V1 readers remain in internal/data for explicit migration and
// frozen history, but the live command surfaces below must use strict V2
// loaders (or neutral project-root resolution) so a future refactor cannot
// quietly restore schema-dispatch fallback.
func TestLiveSourcesDoNotReachV1Readers(t *testing.T) {
	root := repositoryRoot(t)
	paths := []string{
		filepath.Join(root, "main.go"),
		filepath.Join(root, "internal", "board", "board.go"),
		filepath.Join(root, "internal", "doctor", "v2_runtime.go"),
	}
	paths = append(paths,
		productionGoFiles(t, filepath.Join(root, "cmd"))...,
	)
	paths = append(paths,
		productionGoFiles(t, filepath.Join(root, "internal", "board", "v2"))...,
	)
	paths = append(paths,
		productionGoFiles(t, filepath.Join(root, "internal", "resume"))...,
	)
	paths = append(paths,
		productionGoFiles(t, filepath.Join(root, "internal", "init"))...,
	)

	forbidden := []*regexp.Regexp{
		regexp.MustCompile(`data\.NewDiscover\s*\(`),
		regexp.MustCompile(`data\.NewParser\s*\(`),
		regexp.MustCompile(`data\.LoadProject\s*\(`),
		regexp.MustCompile(`\.ParseTaskFile\s*\(`),
		regexp.MustCompile(`\.ParseDefectFile\s*\(`),
		regexp.MustCompile(`\.ParseFindingFile\s*\(`),
		regexp.MustCompile(`data\.RouterState\b`),
		regexp.MustCompile(`data\.Task\b`),
		regexp.MustCompile(`data\.Defect\b`),
		regexp.MustCompile(`data\.AuditRegisterSet\b`),
		regexp.MustCompile(`data\.ReleaseInfo\b`),
		regexp.MustCompile(`data\.EpicInfo\b`),
	}

	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read live source %s: %v", path, err)
		}
		for _, pattern := range forbidden {
			if match := pattern.FindString(string(content)); match != "" {
				t.Errorf("%s reaches legacy V1 surface %q; live commands must use V2 readers", path, match)
			}
		}
	}
}

func TestFindProjectRootWalksWithoutV1Discovery(t *testing.T) {
	project := t.TempDir()
	if err := os.Mkdir(filepath.Join(project, ".savepoint"), 0o755); err != nil {
		t.Fatal(err)
	}
	nested := filepath.Join(project, "nested", "work")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}

	got, err := FindProjectRoot(nested)
	if err != nil {
		t.Fatalf("FindProjectRoot() error = %v", err)
	}
	if got != project {
		t.Fatalf("FindProjectRoot() = %q, want %q", got, project)
	}
}

func TestFindProjectRootRefusesMissingProject(t *testing.T) {
	_, err := FindProjectRoot(t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "no .savepoint directory") {
		t.Fatalf("FindProjectRoot() error = %v, want named missing-project diagnostic", err)
	}
}

func repositoryRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func productionGoFiles(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read live source directory %s: %v", dir, err)
	}
	var paths []string
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".go" || strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		paths = append(paths, filepath.Join(dir, entry.Name()))
	}
	return paths
}
