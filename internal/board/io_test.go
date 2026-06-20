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

func docByID(docs []data.ReleaseDoc, id data.ReleaseDocID) (data.ReleaseDoc, bool) {
	for _, doc := range docs {
		if doc.ID == id {
			return doc, true
		}
	}
	return data.ReleaseDoc{}, false
}

func TestLoadReleaseDocsCmd_ReturnsReleaseAndOverallDocs(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "releases", "v1.2", "v1.2-PRD.md"), "# Release PRD\nv1.2")
	testutil.WriteFile(t, filepath.Join(root, "PRD.md"), "# Overall PRD\nvision")
	testutil.WriteFile(t, filepath.Join(root, "Design.md"), "# Overall Design\narch")

	msg := loadReleaseDocsCmd(root, "v1.2")()

	docsMsg, ok := msg.(releaseDocsMsg)
	if !ok {
		t.Fatalf("msg type = %T, want releaseDocsMsg", msg)
	}
	relPRD, ok := docByID(docsMsg.docs, data.ReleaseDocReleasePRD)
	if !ok || !relPRD.Available || relPRD.Body != "# Release PRD\nv1.2" {
		t.Errorf("Release PRD = %+v, want available v1.2 release PRD", relPRD)
	}
	overallPRD, ok := docByID(docsMsg.docs, data.ReleaseDocOverallPRD)
	if !ok || !overallPRD.Available || overallPRD.Body != "# Overall PRD\nvision" {
		t.Errorf("Overall PRD = %+v, want available root PRD", overallPRD)
	}
	overallDesign, ok := docByID(docsMsg.docs, data.ReleaseDocOverallDesign)
	if !ok || !overallDesign.Available || overallDesign.Body != "# Overall Design\narch" {
		t.Errorf("Overall Design = %+v, want available root Design", overallDesign)
	}
}

// TestLoadReleaseDocsCmd_MissingDocIsNotFatal proves a missing supporting
// document yields an unavailable entry rather than aborting the overlay load.
func TestLoadReleaseDocsCmd_MissingDocIsNotFatal(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "PRD.md"), "# Overall PRD")

	msg := loadReleaseDocsCmd(root, "v1.2")()

	docsMsg, ok := msg.(releaseDocsMsg)
	if !ok {
		t.Fatalf("msg type = %T, want releaseDocsMsg (missing doc must not be fatal)", msg)
	}
	relPRD, ok := docByID(docsMsg.docs, data.ReleaseDocReleasePRD)
	if !ok {
		t.Fatal("releaseDocsMsg missing Release PRD entry")
	}
	if relPRD.Available {
		t.Errorf("Release PRD Available = true, want false when file absent")
	}
}

// TestLoadReleaseDocsCmd_ReadErrorSurfacesAsErrorMsg proves an unexpected read
// error (a directory where a doc file is expected) becomes a status message
// instead of crashing the overlay.
func TestLoadReleaseDocsCmd_ReadErrorSurfacesAsErrorMsg(t *testing.T) {
	root := t.TempDir()
	testutil.MkdirAll(t, filepath.Join(root, "PRD.md"))

	msg := loadReleaseDocsCmd(root, "v1.2")()

	errMsg, ok := msg.(errorMsg)
	if !ok {
		t.Fatalf("msg type = %T, want errorMsg on read failure", msg)
	}
	if !strings.Contains(errMsg.message, "PRD.md") {
		t.Errorf("message = %q, want path context containing PRD.md", errMsg.message)
	}
}
