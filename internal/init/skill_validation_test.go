package init

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// skillRoots are the two trees that must stay in agreement: the skills this
// repository runs on, and the copies scaffolded into generated projects.
func skillRoots() map[string]string {
	return map[string]string{
		"live":     filepath.Join("..", "..", "agent-skills"),
		"template": filepath.Join("..", "..", "templates", "project", "agent-skills"),
	}
}

// splitAuditSkillNames are the two skills that replaced the generic audit skill.
var splitAuditSkillNames = []string{"savepoint-audit-task", "savepoint-audit-epic"}

func savepointSkillDirs(t *testing.T, root string) []string {
	t.Helper()

	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatalf("read skill root %s: %v", root, err)
	}

	var names []string
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "savepoint-") {
			names = append(names, entry.Name())
		}
	}
	if len(names) == 0 {
		t.Fatalf("no savepoint skills under %s", root)
	}
	return names
}

// frontmatterField reads one scalar key out of a leading YAML frontmatter block.
// It returns "" when the block or key is absent. CRLF checkouts are tolerated
// because skills ship to Windows projects.
func frontmatterField(content, key string) string {
	normalized := strings.ReplaceAll(content, "\r\n", "\n")
	if !strings.HasPrefix(normalized, "---\n") {
		return ""
	}
	end := strings.Index(normalized[4:], "\n---")
	if end < 0 {
		return ""
	}

	for _, line := range strings.Split(normalized[4:4+end], "\n") {
		name, value, found := strings.Cut(line, ":")
		if found && strings.TrimSpace(name) == key {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

// sectionBody returns the content between an H2 heading and the next H2, so a
// heading that exists but says nothing can be rejected.
func sectionBody(content, heading string) (string, bool) {
	normalized := strings.ReplaceAll(content, "\r\n", "\n")
	start := strings.Index(normalized, "\n"+heading+"\n")
	if start < 0 {
		return "", false
	}

	body := normalized[start+len(heading)+2:]
	if next := strings.Index(body, "\n## "); next >= 0 {
		body = body[:next]
	}
	return strings.TrimSpace(body), true
}

func TestSavepointSkillsHaveValidFrontmatter(t *testing.T) {
	for tree, root := range skillRoots() {
		for _, name := range savepointSkillDirs(t, root) {
			path := filepath.Join(root, name, "SKILL.md")
			data, err := os.ReadFile(path)
			if err != nil {
				t.Errorf("%s: read %s: %v", tree, path, err)
				continue
			}
			content := string(data)

			if got := frontmatterField(content, "name"); got != name {
				t.Errorf("%s: %s frontmatter name = %q, want folder name %q", tree, path, got, name)
			}
			if frontmatterField(content, "description") == "" {
				t.Errorf("%s: %s has no frontmatter description", tree, path)
			}
		}
	}
}

func TestSavepointSkillsHaveNonEmptyTriggerAndWorkflow(t *testing.T) {
	for tree, root := range skillRoots() {
		for _, name := range savepointSkillDirs(t, root) {
			path := filepath.Join(root, name, "SKILL.md")
			data, err := os.ReadFile(path)
			if err != nil {
				t.Errorf("%s: read %s: %v", tree, path, err)
				continue
			}

			for _, heading := range []string{"## Trigger", "## Workflow"} {
				body, found := sectionBody(string(data), heading)
				if !found {
					t.Errorf("%s: %s missing %s", tree, path, heading)
					continue
				}
				if body == "" {
					t.Errorf("%s: %s has an empty %s section", tree, path, heading)
				}
			}
		}
	}
}

func TestSplitAuditSkillsPassStructureValidation(t *testing.T) {
	required := []string{"## Purpose", "## Trigger", "## Read", "## Workflow", "## Rules"}

	for tree, root := range skillRoots() {
		for _, name := range splitAuditSkillNames {
			path := filepath.Join(root, name, "SKILL.md")
			data, err := os.ReadFile(path)
			if err != nil {
				t.Errorf("%s: split audit skill %s missing: %v", tree, path, err)
				continue
			}
			content := string(data)

			if got := frontmatterField(content, "name"); got != name {
				t.Errorf("%s: %s frontmatter name = %q, want %q", tree, path, got, name)
			}
			for _, heading := range required {
				body, found := sectionBody(content, heading)
				if !found {
					t.Errorf("%s: %s missing %s", tree, path, heading)
					continue
				}
				if body == "" {
					t.Errorf("%s: %s has an empty %s section", tree, path, heading)
				}
			}
		}
	}
}

// The shared audit method is loaded by both audit skills but must never trigger
// on its own, so it carries reference frontmatter rather than skill frontmatter.
func TestSharedAuditMethodIsNonTriggerableReference(t *testing.T) {
	for tree, root := range skillRoots() {
		path := filepath.Join(root, "references", "audit-method.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		if got := frontmatterField(content, "type"); got != "audit-method-reference" {
			t.Errorf("%s: %s type = %q, want audit-method-reference", tree, path, got)
		}
		if got := frontmatterField(content, "triggerable"); got != "false" {
			t.Errorf("%s: %s triggerable = %q, want false", tree, path, got)
		}
		if frontmatterField(content, "name") != "" {
			t.Errorf("%s: %s carries a skill name and would be discoverable as a skill", tree, path)
		}
	}
}
