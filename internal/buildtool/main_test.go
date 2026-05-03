package main

import (
	"os"
	"runtime"
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
