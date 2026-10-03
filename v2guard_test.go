package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func project(t *testing.T, config string) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".savepoint"), 0o755); err != nil {
		t.Fatal(err)
	}
	if config != "" {
		if err := os.WriteFile(filepath.Join(dir, ".savepoint", "config.yml"), []byte(config), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func TestV2ProjectNoticeNamesTheFix(t *testing.T) {
	dir := project(t, "verify_strict: false\nschema_version: 2\n")
	notice, ok := v2ProjectNotice(dir, []string{"health", "check", "O-002"})
	if !ok {
		t.Fatal("a schema_version: 2 project was not recognised")
	}
	for _, want := range []string{"Savepoint V2", "cannot read it", "npx savepoint@latest health check O-002", "npm install -D savepoint@latest"} {
		if !strings.Contains(notice, want) {
			t.Errorf("notice is missing %q:\n%s", want, notice)
		}
	}
}

func TestV2ProjectNoticeFindsTheProjectFromASubdirectory(t *testing.T) {
	dir := project(t, "schema_version: 2\n")
	sub := filepath.Join(dir, "src", "deep")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, ok := v2ProjectNotice(sub, nil); !ok {
		t.Error("V2 project above the working directory was not recognised")
	}
}

func TestV2ProjectNoticeLeavesV1AndUnreadableProjectsAlone(t *testing.T) {
	for name, config := range map[string]string{
		"no schema_version":       "verify_strict: false\n",
		"schema_version 1":        "schema_version: 1\n",
		"indented, not top-level": "nested:\n  schema_version: 2\n",
		"commented out":           "# schema_version: 2\n",
		"empty config":            "",
	} {
		if notice, ok := v2ProjectNotice(project(t, config), []string{"board"}); ok {
			t.Errorf("%s: unexpected notice:\n%s", name, notice)
		}
	}
	if _, ok := v2ProjectNotice(t.TempDir(), []string{"board"}); ok {
		t.Error("a directory with no .savepoint was treated as V2")
	}
}
