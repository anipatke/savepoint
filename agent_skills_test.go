package main_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBundledSavepointSkillsHaveDiscoveryFrontmatter(t *testing.T) {
	assertSavepointSkillsHaveFrontmatter(t, filepath.Join("agent-skills"))
	assertSavepointSkillsHaveFrontmatter(t, filepath.Join("templates", "project", "agent-skills"))
	assertSavepointSkillsHaveFrontmatter(t, filepath.Join("templates", "project-v2", "agent-skills"))
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

	// V1 and V2 skills ship from separate trees (see .savepoint/Guardrails.md
	// TPL-01); each live skill is checked against whichever tree actually
	// carries it, rather than assuming a single shipped tree.
	shippedTrees := []string{
		filepath.Join("templates", "project", "agent-skills"),
		filepath.Join("templates", "project-v2", "agent-skills"),
	}

	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "savepoint-") {
			continue
		}

		sourcePath := filepath.Join(root, entry.Name(), "SKILL.md")
		source, err := os.ReadFile(sourcePath)
		if err != nil {
			t.Fatalf("ReadFile(%q) error = %v", sourcePath, err)
		}

		var scaffoldPath string
		var scaffold []byte
		for _, tree := range shippedTrees {
			candidate := filepath.Join(tree, entry.Name(), "SKILL.md")
			data, err := os.ReadFile(candidate)
			if err == nil {
				scaffoldPath, scaffold = candidate, data
				break
			}
			if !os.IsNotExist(err) {
				t.Fatalf("ReadFile(%q) error = %v", candidate, err)
			}
		}
		if scaffoldPath == "" {
			t.Fatalf("%s is not shipped in any known template tree", sourcePath)
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

		text := strings.ReplaceAll(string(content), "\r\n", "\n")
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
