package init

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// liveSkillRoot is the canonical skill source this repository runs on. The
// shipped V2 tree below must stay byte-identical to whichever of its skills
// and references it carries.
func liveSkillRoot() string {
	return filepath.Join("..", "..", "agent-skills")
}

// v2TemplateSkillRoot is the shipped tree for a V2 project: the four V2
// skills plus their three shared references.
func v2TemplateSkillRoot() string {
	return filepath.Join("..", "..", "templates", "project-v2", "agent-skills")
}

// v2SkillRoots pairs the live source with the V2 shipped tree. Use it for
// tests scoped to the four V2 skills and their three shared references, and
// for discovery-based checks that validate whatever skills a root happens to
// carry (frontmatter shape, non-empty sections).
func v2SkillRoots() map[string]string {
	return map[string]string{
		"live":     liveSkillRoot(),
		"template": v2TemplateSkillRoot(),
	}
}

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
	for tree, root := range v2SkillRoots() {
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
	for tree, root := range v2SkillRoots() {
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
