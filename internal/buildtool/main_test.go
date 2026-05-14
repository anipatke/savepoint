package main

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestVersion_override(t *testing.T) {
	versionOverride = "v1.2.3"
	defer func() { versionOverride = "" }()
	if got := version(); got != "v1.2.3" {
		t.Errorf("version() = %q, want %q", got, "v1.2.3")
	}
}

func TestVersion_env(t *testing.T) {
	versionOverride = ""
	os.Setenv("VERSION", "v2.0.0-env")
	defer os.Unsetenv("VERSION")
	if got := version(); got != "v2.0.0-env" {
		t.Errorf("version() = %q, want %q", got, "v2.0.0-env")
	}
}

func TestVersion_fallback(t *testing.T) {
	versionOverride = ""
	os.Unsetenv("VERSION")
	got := version()
	if got == "" {
		t.Error("version() returned empty string")
	}
}

func TestWriteChecksums(t *testing.T) {
	dir := t.TempDir()

	content := []byte("fake archive content")
	archive := filepath.Join(dir, "savepoint-v1.0.0-linux-amd64.tar.gz")
	if err := os.WriteFile(archive, content, 0o644); err != nil {
		t.Fatal(err)
	}

	dest := filepath.Join(dir, "checksums.txt")
	if err := writeChecksums(dest, []string{archive}); err != nil {
		t.Fatalf("writeChecksums: %v", err)
	}

	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}

	h := sha256.Sum256(content)
	wantHash := hex.EncodeToString(h[:])
	wantLine := wantHash + "  savepoint-v1.0.0-linux-amd64.tar.gz"

	lines := strings.Split(strings.TrimSpace(string(got)), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	if lines[0] != wantLine {
		t.Errorf("line = %q, want %q", lines[0], wantLine)
	}
}

func TestWriteChecksums_multiple(t *testing.T) {
	dir := t.TempDir()

	names := []string{"a.tar.gz", "b.tar.gz"}
	var paths []string
	for _, name := range names {
		p := filepath.Join(dir, name)
		if err := os.WriteFile(p, []byte(name), 0o644); err != nil {
			t.Fatal(err)
		}
		paths = append(paths, p)
	}

	dest := filepath.Join(dir, "checksums.txt")
	if err := writeChecksums(dest, paths); err != nil {
		t.Fatalf("writeChecksums: %v", err)
	}

	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(got)), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d: %s", len(lines), got)
	}
	for i, name := range names {
		h := sha256.Sum256([]byte(name))
		want := hex.EncodeToString(h[:]) + "  " + name
		if lines[i] != want {
			t.Errorf("line[%d] = %q, want %q", i, lines[i], want)
		}
	}
}

func TestWriteChecksums_missingFile(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "checksums.txt")
	err := writeChecksums(dest, []string{filepath.Join(dir, "nonexistent.tar.gz")})
	if err == nil {
		t.Error("expected error for missing archive, got nil")
	}
}

func TestTargets_includesWindows(t *testing.T) {
	var gotAMD64, gotARM64 bool
	for _, tgt := range targets {
		if tgt.os != "windows" {
			continue
		}
		switch tgt.arch {
		case "amd64":
			gotAMD64 = true
		case "arm64":
			gotARM64 = true
		}
	}
	if !gotAMD64 {
		t.Error("targets missing windows/amd64")
	}
	if !gotARM64 {
		t.Error("targets missing windows/arm64")
	}
}

func TestTargets_preservesLinuxDarwin(t *testing.T) {
	want := map[string]bool{
		"linux/amd64":  false,
		"linux/arm64":  false,
		"darwin/amd64": false,
		"darwin/arm64": false,
	}
	for _, tgt := range targets {
		key := tgt.os + "/" + tgt.arch
		if _, ok := want[key]; ok {
			want[key] = true
		}
	}
	for key, found := range want {
		if !found {
			t.Errorf("targets missing %s", key)
		}
	}
}

func TestExecutableName(t *testing.T) {
	if got := executableName("windows"); got != "savepoint.exe" {
		t.Errorf("executableName(windows) = %q, want savepoint.exe", got)
	}
	if got := executableName("linux"); got != "savepoint" {
		t.Errorf("executableName(linux) = %q, want savepoint", got)
	}
	if got := executableName("darwin"); got != "savepoint" {
		t.Errorf("executableName(darwin) = %q, want savepoint", got)
	}
}

func TestWriteTarGzPreservesWindowsExecutableName(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "savepoint.exe")
	if err := os.WriteFile(source, []byte("binary"), 0o755); err != nil {
		t.Fatal(err)
	}

	archive := filepath.Join(dir, "savepoint-windows-amd64.tar.gz")
	if err := writeTarGz(archive, source, executableName("windows")); err != nil {
		t.Fatalf("writeTarGz: %v", err)
	}

	f, err := os.Open(archive)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	header, err := tr.Next()
	if err != nil {
		t.Fatal(err)
	}
	if header.Name != "savepoint.exe" {
		t.Fatalf("archive member = %q, want savepoint.exe", header.Name)
	}

	content, err := io.ReadAll(tr)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "binary" {
		t.Fatalf("archive content = %q, want binary", content)
	}
}

func TestRequireWindowsExecutableAcceptsMZHeader(t *testing.T) {
	path := filepath.Join(t.TempDir(), "savepoint.exe")
	if err := os.WriteFile(path, []byte("MZfake-pe"), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := requireWindowsExecutable(path); err != nil {
		t.Fatalf("requireWindowsExecutable: %v", err)
	}
}

func TestRequireWindowsExecutableRejectsELFHeader(t *testing.T) {
	path := filepath.Join(t.TempDir(), "savepoint.exe")
	if err := os.WriteFile(path, []byte{0x7f, 'E', 'L', 'F'}, 0o755); err != nil {
		t.Fatal(err)
	}

	err := requireWindowsExecutable(path)
	if err == nil {
		t.Fatal("expected ELF header to be rejected")
	}
	if !strings.Contains(err.Error(), "not a Windows PE binary") {
		t.Fatalf("error = %q, want Windows PE binary message", err)
	}
}

func TestLocalExecutable(t *testing.T) {
	got := localExecutable()
	if runtime.GOOS == "windows" {
		if got != "savepoint.exe" {
			t.Errorf("localExecutable() = %q, want %q", got, "savepoint.exe")
		}
	} else {
		if got != "savepoint" {
			t.Errorf("localExecutable() = %q, want %q", got, "savepoint")
		}
	}
}
