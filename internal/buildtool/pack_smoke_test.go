package main

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestParseNpmPackFilename(t *testing.T) {
	output := []byte(`[{"id":"savepoint@1.2.3","name":"savepoint","version":"1.2.3","filename":"savepoint-1.2.3.tgz"}]`)
	got, err := parseNpmPackFilename(output)
	if err != nil {
		t.Fatalf("parseNpmPackFilename: %v", err)
	}
	if got != "savepoint-1.2.3.tgz" {
		t.Errorf("filename = %q, want savepoint-1.2.3.tgz", got)
	}
}

func TestParseNpmPackFilename_trimsWhitespace(t *testing.T) {
	output := []byte("\n  [{\"filename\":\"x.tgz\"}]  \n")
	got, err := parseNpmPackFilename(output)
	if err != nil {
		t.Fatalf("parseNpmPackFilename: %v", err)
	}
	if got != "x.tgz" {
		t.Errorf("filename = %q, want x.tgz", got)
	}
}

func TestParseNpmPackFilename_emptyOutput(t *testing.T) {
	if _, err := parseNpmPackFilename([]byte("   \n")); err == nil {
		t.Fatal("expected error on empty output")
	}
}

func TestParseNpmPackFilename_invalidJSON(t *testing.T) {
	if _, err := parseNpmPackFilename([]byte("not json")); err == nil {
		t.Fatal("expected error on invalid json")
	}
}

func TestParseNpmPackFilename_emptyArray(t *testing.T) {
	if _, err := parseNpmPackFilename([]byte("[]")); err == nil {
		t.Fatal("expected error on empty array")
	}
}

func TestParseNpmPackFilename_missingFilename(t *testing.T) {
	if _, err := parseNpmPackFilename([]byte(`[{"name":"x"}]`)); err == nil {
		t.Fatal("expected error when filename missing")
	}
}

func TestInstalledBinaryPath_windows(t *testing.T) {
	got := installedBinaryPath(filepath.Join("tmp", "install"), "windows")
	want := filepath.Join("tmp", "install", "node_modules", ".bin", "savepoint.cmd")
	if got != want {
		t.Errorf("path = %q, want %q", got, want)
	}
}

func TestInstalledBinaryPath_unix(t *testing.T) {
	got := installedBinaryPath("/tmp/install", "linux")
	want := filepath.Join("/tmp/install", "node_modules", ".bin", "savepoint")
	if got != want {
		t.Errorf("path = %q, want %q", got, want)
	}
	got = installedBinaryPath("/tmp/install", "darwin")
	if !strings.HasSuffix(got, filepath.Join(".bin", "savepoint")) {
		t.Errorf("darwin path = %q missing savepoint suffix", got)
	}
}

func TestNpmExecutable_matchesHost(t *testing.T) {
	got := npmExecutable()
	if runtime.GOOS == "windows" {
		if got != "npm.cmd" {
			t.Errorf("npmExecutable = %q, want npm.cmd", got)
		}
	} else {
		if got != "npm" {
			t.Errorf("npmExecutable = %q, want npm", got)
		}
	}
}

func TestHostTarget_matchesRuntime(t *testing.T) {
	tgt, err := hostTarget()
	if err != nil {
		t.Skipf("host %s/%s not in supported targets: %v", runtime.GOOS, runtime.GOARCH, err)
	}
	if tgt.os != runtime.GOOS || tgt.arch != runtime.GOARCH {
		t.Errorf("hostTarget = %+v, want {%s %s}", tgt, runtime.GOOS, runtime.GOARCH)
	}
}

func TestSmokeProjectManifestJSON_isValidPrivatePackage(t *testing.T) {
	if !strings.Contains(smokeProjectManifestJSON, `"private":true`) {
		t.Error("smoke manifest must be private to prevent accidental publish")
	}
	if !strings.HasSuffix(smokeProjectManifestJSON, "\n") {
		t.Error("smoke manifest must end with newline")
	}
}
