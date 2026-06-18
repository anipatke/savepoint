package board

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/testutil"
)

func TestWriteEpicStatusCmd_WritesAuditedToFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "E31-Detail.md")
	testutil.WriteFile(t, path, "---\ntype: epic-design\nstatus: done\n---\n\n# Epic\n")
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	msg := writeEpicStatusCmd("E31-epic-audited-shortcut", path, string(data.EpicStatusAudited), fi.ModTime())()

	written, ok := msg.(epicStatusWrittenMsg)
	if !ok {
		t.Fatalf("msg type = %T, want epicStatusWrittenMsg", msg)
	}
	if written.epicID != "E31-epic-audited-shortcut" {
		t.Errorf("written.epicID = %q, want E31-epic-audited-shortcut", written.epicID)
	}
	if written.status != string(data.EpicStatusAudited) {
		t.Errorf("written.status = %q, want audited", written.status)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), "status: audited") {
		t.Fatalf("file not updated to audited:\n%s", raw)
	}
}

func TestWriteEpicStatusCmd_ReportsConflictAndLeavesFileUntouched(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "E31-Detail.md")
	testutil.WriteFile(t, path, "---\ntype: epic-design\nstatus: done\n---\n\n# Epic\n")
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	stale := fi.ModTime().Add(-time.Hour)
	msg := writeEpicStatusCmd("E31-epic-audited-shortcut", path, string(data.EpicStatusAudited), stale)()

	errMsg, ok := msg.(errorMsg)
	if !ok {
		t.Fatalf("msg type = %T, want errorMsg", msg)
	}
	if !strings.Contains(errMsg.message, "epic changed on disk") {
		t.Errorf("message = %q, want conflict message", errMsg.message)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), "status: done") {
		t.Fatalf("file should be untouched on conflict:\n%s", raw)
	}
}
