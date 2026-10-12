package init

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

var orchestrateAssets = []string{
	"agent-skills/savepoint-orchestrate/SKILL.md",
	".claude/skills/savepoint-orchestrate/SKILL.md",
}

func TestOrchestrateAssetsScaffoldAndUpgrade(t *testing.T) {
	templates := os.DirFS(filepath.Join("..", "..", "templates", "project-v2"))
	t.Run("scaffold", func(t *testing.T) {
		dir := t.TempDir()
		if err := Scaffold(templates, dir, "example", false); err != nil {
			t.Fatal(err)
		}
		for _, path := range orchestrateAssets {
			want := mustReadFile(t, filepath.Join("..", "..", "templates", "project-v2", filepath.FromSlash(path)))
			if got := mustReadFile(t, filepath.Join(dir, filepath.FromSlash(path))); !bytes.Equal(got, want) {
				t.Fatalf("%s differs from packaged asset", path)
			}
		}
	})
	t.Run("upgrade installs and repeats safely", func(t *testing.T) {
		dir := savepointProject(t)
		for attempt := 0; attempt < 2; attempt++ {
			report, err := upgradeAssetsFromTree(templates, dir, false, false)
			if err != nil {
				t.Fatal(err)
			}
			for _, path := range orchestrateAssets {
				want := mustReadFile(t, filepath.Join("..", "..", "templates", "project-v2", filepath.FromSlash(path)))
				if got := mustReadFile(t, filepath.Join(dir, filepath.FromSlash(path))); !bytes.Equal(got, want) {
					t.Fatalf("%s differs from packaged asset", path)
				}
				if attempt == 1 && upgradeActionFor(t, report, path) != ActionUnchanged {
					t.Fatalf("repeat upgrade changed %s", path)
				}
			}
		}
	})
	t.Run("upgrade preserves edited copies", func(t *testing.T) {
		dir := savepointProject(t)
		if _, err := upgradeAssetsFromTree(templates, dir, false, false); err != nil {
			t.Fatal(err)
		}
		const edited = "owner-authored orchestration guidance\n"
		for _, path := range orchestrateAssets {
			if err := os.WriteFile(filepath.Join(dir, filepath.FromSlash(path)), []byte(edited), 0644); err != nil {
				t.Fatal(err)
			}
		}
		report, err := upgradeAssetsFromTree(templates, dir, false, false)
		if err != nil {
			t.Fatal(err)
		}
		for _, path := range orchestrateAssets {
			if got := string(mustReadFile(t, filepath.Join(dir, filepath.FromSlash(path)))); got != edited {
				t.Fatalf("overwrote edited %s", path)
			}
			if upgradeActionFor(t, report, path) != ActionConflict {
				t.Fatalf("expected conflict for %s", path)
			}
		}
	})
}

func TestOrchestrateClaudePointerParity(t *testing.T) {
	assertFileMatches(t, filepath.Join("..", ".."),
		filepath.Join(".claude", "skills", "savepoint-orchestrate", "SKILL.md"),
		filepath.Join("templates", "project-v2", ".claude", "skills", "savepoint-orchestrate", "SKILL.md"))
}
