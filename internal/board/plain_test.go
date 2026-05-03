package board

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/data"
)

func TestRenderPlainTable_warningBanner(t *testing.T) {
	m := NewModel(nil, "", "")
	got := RenderPlainTable(m)
	if !strings.Contains(got, plainNonTTYWarning) {
		t.Errorf("RenderPlainTable missing warning banner, got:\n%s", got)
	}
}

func TestRenderPlainTable_columnHeaders(t *testing.T) {
	m := NewModel(nil, "", "")
	got := RenderPlainTable(m)
	for _, header := range []string{"PLANNED", "IN PROGRESS", "DONE"} {
		if !strings.Contains(got, header) {
			t.Errorf("RenderPlainTable missing column header %q", header)
		}
	}
}

func TestRenderPlainTable_taskIDsAndTitles(t *testing.T) {
	tasks := []data.Task{
		{ID: "E08/T001", Title: "CLI entrypoint", Column: data.ColumnDone},
		{ID: "E08/T002", Title: "Non-TTY fallback", Column: data.ColumnInProgress},
		{ID: "E08/T003", Title: "TUI app shell", Column: data.ColumnPlanned},
	}
	m := NewModel(tasks, "", "")
	got := RenderPlainTable(m)

	for _, want := range []string{"E08/T001", "CLI entrypoint", "E08/T002", "Non-TTY fallback", "E08/T003", "TUI app shell"} {
		if !strings.Contains(got, want) {
			t.Errorf("RenderPlainTable missing %q", want)
		}
	}
}

func TestRenderPlainTable_noneWhenColumnEmpty(t *testing.T) {
	m := NewModel(nil, "", "")
	got := RenderPlainTable(m)
	if !strings.Contains(got, "(none)") {
		t.Errorf("RenderPlainTable should show (none) for empty columns")
	}
}

func TestRenderPlainTable_auditSignalWhenProposalsExist(t *testing.T) {
	root := t.TempDir()
	epicDir := filepath.Join(root, "releases", "v1", "epics", "E01-test")
	if err := os.MkdirAll(epicDir, 0755); err != nil {
		t.Fatal(err)
	}
	auditContent := "---\ntype: audit\n---\n## Main Findings\nOK\n\n## Proposed Changes\n\n### Target File\nfoo.go\n"
	if err := os.WriteFile(filepath.Join(epicDir, "E01-Audit.md"), []byte(auditContent), 0644); err != nil {
		t.Fatal(err)
	}

	m := NewModel(nil, "", "")
	m.Root = root
	got := RenderPlainTable(m)
	if !strings.Contains(got, plainAuditSignal) {
		t.Errorf("RenderPlainTable missing audit signal when proposals exist, got:\n%s", got)
	}
}

func TestRenderPlainTable_noAuditSignalWhenNone(t *testing.T) {
	root := t.TempDir()
	m := NewModel(nil, "", "")
	m.Root = root
	got := RenderPlainTable(m)
	if strings.Contains(got, plainAuditSignal) {
		t.Errorf("RenderPlainTable should not show audit signal when no proposals exist")
	}
}

func TestHasAuditProposals_detectsSection(t *testing.T) {
	root := t.TempDir()
	epicDir := filepath.Join(root, "releases", "v1", "epics", "E02-slug")
	if err := os.MkdirAll(epicDir, 0755); err != nil {
		t.Fatal(err)
	}
	content := "## Main Findings\nAll good.\n\n## Proposed Changes\n\n### Target File\nbar.go\n"
	if err := os.WriteFile(filepath.Join(epicDir, "E02-Audit.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	if !hasAuditProposals(root) {
		t.Error("hasAuditProposals should return true when Proposed Changes section exists")
	}
}

func TestHasAuditProposals_noProposals(t *testing.T) {
	root := t.TempDir()
	epicDir := filepath.Join(root, "releases", "v1", "epics", "E03-slug")
	if err := os.MkdirAll(epicDir, 0755); err != nil {
		t.Fatal(err)
	}
	content := "## Main Findings\nAll good.\n"
	if err := os.WriteFile(filepath.Join(epicDir, "E03-Audit.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	if hasAuditProposals(root) {
		t.Error("hasAuditProposals should return false when no Proposed Changes section")
	}
}

func TestHasAuditProposals_missingRoot(t *testing.T) {
	if hasAuditProposals("/nonexistent/path/xyz") {
		t.Error("hasAuditProposals should return false for missing root")
	}
}
