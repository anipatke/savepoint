package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPlatformPackageName(t *testing.T) {
	got := platformPackageName(target{os: "linux", arch: "amd64"})
	if got != "savepoint-linux-amd64" {
		t.Errorf("platformPackageName = %q, want savepoint-linux-amd64", got)
	}
}

func TestBuildNPMManifest_linuxAmd64(t *testing.T) {
	m, err := buildNPMManifest(target{os: "linux", arch: "amd64"}, "1.2.3")
	if err != nil {
		t.Fatalf("buildNPMManifest: %v", err)
	}
	if m.Name != "savepoint-linux-amd64" {
		t.Errorf("Name = %q", m.Name)
	}
	if m.Version != "1.2.3" {
		t.Errorf("Version = %q", m.Version)
	}
	if len(m.OS) != 1 || m.OS[0] != "linux" {
		t.Errorf("OS = %v, want [linux]", m.OS)
	}
	if len(m.CPU) != 1 || m.CPU[0] != "x64" {
		t.Errorf("CPU = %v, want [x64]", m.CPU)
	}
	if len(m.Files) != 1 || m.Files[0] != "bin" {
		t.Errorf("Files = %v, want [bin]", m.Files)
	}
}

func TestBuildNPMManifest_windowsArm64(t *testing.T) {
	m, err := buildNPMManifest(target{os: "windows", arch: "arm64"}, "1.2.3")
	if err != nil {
		t.Fatalf("buildNPMManifest: %v", err)
	}
	if m.Name != "savepoint-windows-arm64" {
		t.Errorf("Name = %q", m.Name)
	}
	if m.OS[0] != "win32" {
		t.Errorf("OS = %v, want [win32]", m.OS)
	}
	if m.CPU[0] != "arm64" {
		t.Errorf("CPU = %v, want [arm64]", m.CPU)
	}
}

func TestBuildNPMManifest_darwinArm64(t *testing.T) {
	m, err := buildNPMManifest(target{os: "darwin", arch: "arm64"}, "1.2.3")
	if err != nil {
		t.Fatalf("buildNPMManifest: %v", err)
	}
	if m.OS[0] != "darwin" {
		t.Errorf("OS = %v, want [darwin]", m.OS)
	}
	if m.CPU[0] != "arm64" {
		t.Errorf("CPU = %v, want [arm64]", m.CPU)
	}
}

func TestBuildNPMManifest_allTargetsMappable(t *testing.T) {
	for _, tgt := range targets {
		m, err := buildNPMManifest(tgt, "9.9.9")
		if err != nil {
			t.Errorf("target %s/%s: %v", tgt.os, tgt.arch, err)
			continue
		}
		want := "savepoint-" + tgt.os + "-" + tgt.arch
		if m.Name != want {
			t.Errorf("Name = %q, want %q", m.Name, want)
		}
	}
}

func TestBuildNPMManifest_unknownOS(t *testing.T) {
	_, err := buildNPMManifest(target{os: "plan9", arch: "amd64"}, "1.0.0")
	if err == nil {
		t.Fatal("expected error for unknown os")
	}
}

func TestBuildNPMManifest_unknownArch(t *testing.T) {
	_, err := buildNPMManifest(target{os: "linux", arch: "mips"}, "1.0.0")
	if err == nil {
		t.Fatal("expected error for unknown arch")
	}
}

func TestWriteNPMManifest_writesJSON(t *testing.T) {
	dir := t.TempDir()
	m, err := buildNPMManifest(target{os: "linux", arch: "amd64"}, "1.0.0")
	if err != nil {
		t.Fatal(err)
	}
	if err := writeNPMManifest(dir, m); err != nil {
		t.Fatalf("writeNPMManifest: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "package.json"))
	if err != nil {
		t.Fatal(err)
	}
	var decoded npmManifest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if decoded.Name != "savepoint-linux-amd64" {
		t.Errorf("Name = %q", decoded.Name)
	}
	if !strings.HasSuffix(string(data), "\n") {
		t.Error("manifest should end with newline")
	}
}

func TestReadRootPackageVersion(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "package.json")
	if err := os.WriteFile(path, []byte(`{"name":"savepoint","version":"1.2.3"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := readRootPackageVersion(path)
	if err != nil {
		t.Fatalf("readRootPackageVersion: %v", err)
	}
	if got != "1.2.3" {
		t.Errorf("version = %q, want 1.2.3", got)
	}
}

func TestReadRootPackageVersion_missingVersion(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "package.json")
	if err := os.WriteFile(path, []byte(`{"name":"savepoint"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := readRootPackageVersion(path)
	if err == nil {
		t.Fatal("expected error for missing version")
	}
}

func TestReadRootPackageVersion_invalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "package.json")
	if err := os.WriteFile(path, []byte(`not json`), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := readRootPackageVersion(path)
	if err == nil {
		t.Fatal("expected parse error")
	}
}
