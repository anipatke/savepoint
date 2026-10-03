package legacydep

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func project(t *testing.T, pkg, installed string) string {
	t.Helper()
	dir := t.TempDir()
	if pkg != "" {
		if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte(pkg), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if installed != "" {
		nm := filepath.Join(dir, "node_modules", "savepoint")
		if err := os.MkdirAll(nm, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(nm, "package.json"), []byte(`{"version":"`+installed+`"}`), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func TestDetect(t *testing.T) {
	for _, tc := range []struct {
		name      string
		pkg       string
		installed string
		want      bool
		section   string
	}{
		{"no package.json", "", "", false, ""},
		{"unreadable package.json", "{", "", false, ""},
		{"unrelated project", `{"dependencies":{"react":"^18"}}`, "", false, ""},
		{"caret 1.x dependency", `{"dependencies":{"savepoint":"^1.3.0"}}`, "", true, "dependencies"},
		{"tilde 1.x devDependency", `{"devDependencies":{"savepoint":"~1.2.0"}}`, "", true, "devDependencies"},
		{"wildcard 1.x", `{"dependencies":{"savepoint":"1.x"}}`, "", true, "dependencies"},
		{"current range", `{"devDependencies":{"savepoint":"^2.1.3"}}`, "", false, ""},
		{"open-ended range", `{"devDependencies":{"savepoint":">=1.3.0"}}`, "", false, ""},
		{"range admitting 2", `{"devDependencies":{"savepoint":"^1.3.0 || ^2.0.0"}}`, "", false, ""},
		{"latest tag", `{"devDependencies":{"savepoint":"latest"}}`, "", false, ""},
		{"current range but old install", `{"devDependencies":{"savepoint":"^2.1.3"}}`, "1.3.0", true, ""},
		{"no declaration but old install", `{"name":"x"}`, "1.3.0", true, ""},
		{"current install", `{"devDependencies":{"savepoint":"^2.1.3"}}`, "2.1.3", false, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f, got := Detect(project(t, tc.pkg, tc.installed))
			if got != tc.want {
				t.Fatalf("Detect found = %v, want %v (%+v)", got, tc.want, f)
			}
			if got && f.Section != tc.section {
				t.Errorf("Section = %q, want %q", f.Section, tc.section)
			}
		})
	}
}

func TestFindingWordsNameVersionsAndRepair(t *testing.T) {
	f, ok := Detect(project(t, `{"dependencies":{"savepoint":"^1.3.0"}}`, "1.3.0"))
	if !ok {
		t.Fatal("expected a finding")
	}
	for _, want := range []string{"^1.3.0", "dependencies", "1.3.0", "npx savepoint"} {
		if !strings.Contains(f.Message(), want) {
			t.Errorf("Message() = %q, missing %q", f.Message(), want)
		}
	}
	if !strings.Contains(f.Repair(), "devDependencies") {
		t.Errorf("Repair() = %q", f.Repair())
	}
}

func TestDetectDoesNotEditPackageFiles(t *testing.T) {
	dir := project(t, `{"dependencies":{"savepoint":"^1.3.0"}}`, "")
	before, _ := os.ReadFile(filepath.Join(dir, "package.json"))
	Detect(dir)
	after, _ := os.ReadFile(filepath.Join(dir, "package.json"))
	if string(before) != string(after) {
		t.Error("Detect changed package.json")
	}
}
