package init

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProjectTemplatesUseCurrentWorkflow(t *testing.T) {
	root := filepath.Join("..", "..")
	agents := readTemplate(t, root, "templates", "project", "AGENTS.md")
	router := readTemplate(t, root, "templates", "project", ".savepoint", "router.md")
	auditSkill := readTemplate(t, root, "templates", "project", "agent-skills", "savepoint-audit", "SKILL.md")

	assertNotContains(t, agents, "`phase` (build/test/audit)")
	assertNotContains(t, agents, "npm run build && npm run test")
	assertContains(t, agents, "`stage` (build/test/audit): **required** when `status: in_progress`")
	assertContains(t, agents, "make build && make test")

	assertNotContains(t, router, ".savepoint/audit/{E##-epic}/snapshot.md")
	assertNotContains(t, router, ".savepoint/audit/{release}/{E##-epic}/proposals.md")
	assertNotContains(t, router, ".savepoint/audit/{E##-epic}/proposals.md")
	assertContains(t, router, ".savepoint/releases/{release}/epics/{E##-epic}/E##-Audit.md")
	assertContains(t, router, "`## Proposed Changes` — admin/apply metadata")
	assertContains(t, agents, "During audit apply/close, update the same `E##-Audit.md` visible sections")
	assertContains(t, auditSkill, "Update `E##-Audit.md` visible sections")
	assertContains(t, auditSkill, "Updated audit findings")
}

func readTemplate(t *testing.T, root string, parts ...string) string {
	t.Helper()

	path := filepath.Join(append([]string{root}, parts...)...)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read template %s: %v", path, err)
	}
	return string(data)
}

func assertContains(t *testing.T, content, want string) {
	t.Helper()

	if !strings.Contains(content, want) {
		t.Fatalf("template missing %q", want)
	}
}

func assertNotContains(t *testing.T, content, stale string) {
	t.Helper()

	if strings.Contains(content, stale) {
		t.Fatalf("template contains stale text %q", stale)
	}
}
