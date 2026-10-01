package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeGoManifest(t *testing.T, dir string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example\n\ngo 1.22\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}

func initGoProject(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	writeGoManifest(t, dir)
	if result := runMainForTest(t, []string{"init", dir}, ""); result.err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", result.err, result.stderr)
	}
	return dir
}

func TestMainHelpListsHealthSetup(t *testing.T) {
	result := runMainForTest(t, []string{"--help"}, "")
	if !strings.Contains(result.stdout, "health setup [dir] [--apply]") {
		t.Fatalf("help = %q", result.stdout)
	}
	result = runMainForTest(t, []string{"health", "setup", "--help"}, "")
	if result.err != nil || !strings.Contains(result.stdout, "Usage: health setup [dir] [--apply]") {
		t.Fatalf("health setup --help = %q, %v", result.stdout, result.err)
	}
}

func TestMainInitEndsWithHealthPreviewAndWritesNoHealthConfig(t *testing.T) {
	dir := t.TempDir()
	writeGoManifest(t, dir)

	result := runMainForTest(t, []string{"init", dir}, "")
	if result.err != nil {
		t.Fatalf("init failed: %v\nstderr: %s", result.err, result.stderr)
	}
	for _, want := range []string{"Health tools", "New:", "savepoint health setup --apply", "OSV-Scanner contacts OSV.dev"} {
		if !strings.Contains(result.stdout, want) {
			t.Errorf("init output missing %q", want)
		}
	}
	if _, err := os.Stat(filepath.Join(dir, ".savepoint", "health")); err == nil {
		t.Fatal("init wrote health files")
	}
}

func TestMainInitStillSucceedsWhenHealthDiscoveryFails(t *testing.T) {
	dir := initGoProject(t)
	health := filepath.Join(dir, ".savepoint", "health")
	if err := os.MkdirAll(health, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(health, "config.json"), []byte("{broken"), 0o644); err != nil {
		t.Fatal(err)
	}

	result := runMainForTest(t, []string{"init", dir, "--force"}, "")
	if result.err != nil {
		t.Fatalf("init --force failed: %v\nstderr: %s", result.err, result.stderr)
	}
	// The warning goes to stderr, which runMainForTest only captures on failure;
	// the preview must be absent and the broken file untouched.
	if strings.Contains(result.stdout, "Health tools") {
		t.Errorf("preview printed despite a broken config:\n%s", result.stdout)
	}
	if got, _ := os.ReadFile(filepath.Join(health, "config.json")); string(got) != "{broken" {
		t.Errorf("broken config was modified: %q", got)
	}
}

func TestMainHealthSetupPreviewWritesNothing(t *testing.T) {
	dir := initGoProject(t)
	before := snapshotDir(t, dir)

	result := runMainForTest(t, []string{"health", "setup", dir}, "")
	if result.err != nil {
		t.Fatalf("health setup failed: %v\nstderr: %s", result.err, result.stderr)
	}
	if !strings.Contains(result.stdout, "Nothing was written") || !strings.Contains(result.stdout, "--apply") {
		t.Fatalf("preview = %q", result.stdout)
	}
	assertSameSnapshot(t, before, snapshotDir(t, dir))
}

func TestMainHealthSetupApplyWritesOnlyConfigAndRepeatsAsUnchanged(t *testing.T) {
	dir := initGoProject(t)
	before := snapshotDir(t, dir)

	result := runMainForTest(t, []string{"health", "setup", dir, "--apply"}, "")
	if result.err != nil {
		t.Fatalf("apply failed: %v\nstderr: %s", result.err, result.stderr)
	}
	after := snapshotDir(t, dir)
	configRel := filepath.Join(".savepoint", "health", "config.json")
	for path := range after {
		if _, existed := before[path]; !existed && path != configRel {
			t.Errorf("apply wrote %s", path)
		}
	}
	if _, ok := after[configRel]; !ok {
		t.Fatal("apply did not write the health config")
	}
	if !strings.Contains(result.stdout, "Saved") {
		t.Errorf("apply output = %q", result.stdout)
	}

	result = runMainForTest(t, []string{"health", "setup", dir, "--apply"}, "")
	if result.err != nil {
		t.Fatalf("second apply failed: %v\nstderr: %s", result.err, result.stderr)
	}
	if !strings.Contains(result.stdout, "No changes") {
		t.Errorf("second apply output = %q", result.stdout)
	}
	assertSameSnapshot(t, after, snapshotDir(t, dir))
}

func TestMainHealthSetupRefusesMissingAndNonSavepointDirectories(t *testing.T) {
	plain := t.TempDir()
	writeGoManifest(t, plain)
	for name, dir := range map[string]string{
		"not a Savepoint project": plain,
		"does not exist":          filepath.Join(plain, "nope"),
	} {
		result := runMainForTest(t, []string{"health", "setup", dir, "--apply"}, "")
		if result.err == nil {
			t.Errorf("%s: setup succeeded, want a failure", name)
			continue
		}
		if !strings.Contains(result.stderr, name) {
			t.Errorf("%s: stderr = %q", name, result.stderr)
		}
	}
	if _, err := os.Stat(filepath.Join(plain, ".savepoint")); err == nil {
		t.Fatal("a refused setup created .savepoint")
	}
}

func TestMainHealthRequiresSetupSubcommand(t *testing.T) {
	if result := runMainForTest(t, []string{"health"}, ""); result.err == nil {
		t.Fatal("bare health succeeded")
	}
}

func TestMainHelpListsHealthCheck(t *testing.T) {
	result := runMainForTest(t, []string{"--help"}, "")
	if !strings.Contains(result.stdout, "health check O-### [dir]") {
		t.Fatalf("help = %q", result.stdout)
	}
	result = runMainForTest(t, []string{"health", "--help"}, "")
	if result.err != nil || !strings.Contains(result.stdout, "health setup") || !strings.Contains(result.stdout, "health check O-### [dir]") {
		t.Fatalf("health --help = %q, %v", result.stdout, result.err)
	}
}

func TestMainHealthCheckRefusesBadInputAndUnconfiguredProjects(t *testing.T) {
	dir := initGoProject(t)
	for name, args := range map[string][]string{
		"missing ID":        {"health", "check"},
		"malformed ID":      {"health", "check", "T-001"},
		"unknown flag":      {"health", "check", "O-001", "--manual"},
		"unknown Objective": {"health", "check", "O-999", dir},
	} {
		if result := runMainForTest(t, args, ""); result.err == nil {
			t.Errorf("%s: health check succeeded, want a failure", name)
		}
	}
	if _, err := os.Stat(filepath.Join(dir, ".savepoint", "health", "snapshots")); err == nil {
		t.Fatal("a refused health check saved a snapshot")
	}
}
