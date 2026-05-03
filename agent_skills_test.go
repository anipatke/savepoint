package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBundledSavepointSkillsHaveDiscoveryFrontmatter(t *testing.T) {
	assertSavepointSkillsHaveFrontmatter(t, filepath.Join("agent-skills"))
	assertSavepointSkillsHaveFrontmatter(t, filepath.Join("templates", "project", "agent-skills"))
}

func TestProjectAgentGuideIncludesLocalSkillFallback(t *testing.T) {
	path := filepath.Join("templates", "project", "AGENTS.md")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", path, err)
	}

	want := "If the agent says the skill is not found, read `agent-skills/{skill}/SKILL.md` directly"
	if !strings.Contains(string(content), want) {
		t.Fatalf("%s missing local skill fallback instruction", path)
	}
}

func TestScaffoldedSavepointSkillsMatchBundledSkills(t *testing.T) {
	root := filepath.Join("agent-skills")
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatalf("ReadDir(%q) error = %v", root, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "savepoint-") {
			continue
		}

		sourcePath := filepath.Join(root, entry.Name(), "SKILL.md")
		scaffoldPath := filepath.Join("templates", "project", "agent-skills", entry.Name(), "SKILL.md")
		source, err := os.ReadFile(sourcePath)
		if err != nil {
			t.Fatalf("ReadFile(%q) error = %v", sourcePath, err)
		}
		scaffold, err := os.ReadFile(scaffoldPath)
		if err != nil {
			t.Fatalf("ReadFile(%q) error = %v", scaffoldPath, err)
		}
		if string(scaffold) != string(source) {
			t.Fatalf("%s does not match %s", scaffoldPath, sourcePath)
		}
	}
}

func assertSavepointSkillsHaveFrontmatter(t *testing.T, root string) {
	t.Helper()

	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatalf("ReadDir(%q) error = %v", root, err)
	}

	var found int
	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "savepoint-") {
			continue
		}
		found++
		path := filepath.Join(root, entry.Name(), "SKILL.md")
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("ReadFile(%q) error = %v", path, err)
		}

		text := string(content)
		if !strings.HasPrefix(text, "---\n") {
			t.Fatalf("%s missing YAML frontmatter", path)
		}
		if !strings.Contains(text, "name: "+entry.Name()) {
			t.Fatalf("%s frontmatter name does not match directory", path)
		}
		if !strings.Contains(text, "description:") {
			t.Fatalf("%s missing frontmatter description", path)
		}
	}

	if found == 0 {
		t.Fatalf("%s contains %d savepoint skills, want at least 1", root, found)
	}
}
