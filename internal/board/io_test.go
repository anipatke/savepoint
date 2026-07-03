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

// board-local audit fixtures mirror the data package's known-good records so
// these command tests exercise a genuine load without reaching into data's
// unexported test constants.
const boardAuditPromptContent = "# Audit Prompt\n\nCanonical reusable prompt.\n"

const boardAuditFindingContent = `---
id: F001
title: Token leak in logs
status: open
severity: high
confidence: medium
source_auditor: agent
work_item: E12-auth/T003-redact-tokens
first_seen: 2026-06-01
last_seen: 2026-06-30
proof_needed: A regression test asserting tokens are redacted before write
---

## Summary

Tokens are written to logs in plaintext.
`

const boardAuditRunContent = `---
date: 2026-06-29
auditor: agent
model: claude-opus-4-8
prompt_version: v1
commit: a1b2c3d
mode: full
coverage: examined internal/board
net_new: 1
reopened: 0
verified: 0
deferred: 0
coverage_gaps: 0
---

## Scope

Sweep of the board package.
`

// writeBoardAuditFile writes an audit-register artifact under root/audit/rel.
func writeBoardAuditFile(t *testing.T, root, rel, content string) {
	t.Helper()
	testutil.WriteFile(t, filepath.Join(root, "audit", rel), content)
}

// writeCompleteAuditTree lays down a prompt, register, one finding, and one run
// so a load returns a fully populated set.
func writeCompleteAuditTree(t *testing.T, root string) {
	t.Helper()
	writeBoardAuditFile(t, root, "prompt.md", boardAuditPromptContent)
	writeBoardAuditFile(t, root, "register.md", "# Audit Register\n")
	writeBoardAuditFile(t, root, filepath.Join("findings", "F001-token-leak.md"), boardAuditFindingContent)
	writeBoardAuditFile(t, root, filepath.Join("runs", "2026-06-29-board-sweep.md"), boardAuditRunContent)
}

func TestLoadAuditRegisterCmd_ReturnsPopulatedSetWhenAuditTreeExists(t *testing.T) {
	root := t.TempDir()
	writeCompleteAuditTree(t, root)

	msg := loadAuditRegisterCmd(root, dataAuditLoader{})()

	auditMsg, ok := msg.(auditRegisterMsg)
	if !ok {
		t.Fatalf("msg type = %T, want auditRegisterMsg", msg)
	}
	if !auditMsg.set.Prompt.Available || !auditMsg.set.Register.Available {
		t.Errorf("prompt/register available = %v/%v, want true/true",
			auditMsg.set.Prompt.Available, auditMsg.set.Register.Available)
	}
	if len(auditMsg.set.Findings) != 1 || auditMsg.set.Findings[0].ID != "F001" {
		t.Errorf("findings = %v, want one F001", auditMsg.set.Findings)
	}
	if len(auditMsg.set.Runs) != 1 {
		t.Errorf("runs = %d, want 1", len(auditMsg.set.Runs))
	}
}

// TestLoadAuditRegisterCmd_MissingTreeIsEmptyNotError proves an absent audit/
// tree yields an empty data state rather than an error, so a project without an
// audit register still opens the overlay cleanly.
func TestLoadAuditRegisterCmd_MissingTreeIsEmptyNotError(t *testing.T) {
	msg := loadAuditRegisterCmd(t.TempDir(), dataAuditLoader{})()

	auditMsg, ok := msg.(auditRegisterMsg)
	if !ok {
		t.Fatalf("msg type = %T, want auditRegisterMsg (missing tree must not error)", msg)
	}
	if auditMsg.set.Prompt.Available || auditMsg.set.Register.Available {
		t.Errorf("prompt/register available = %v/%v, want false for absent audit tree",
			auditMsg.set.Prompt.Available, auditMsg.set.Register.Available)
	}
	if len(auditMsg.set.Findings) != 0 || len(auditMsg.set.Runs) != 0 {
		t.Errorf("findings/runs = %d/%d, want empty for absent audit tree",
			len(auditMsg.set.Findings), len(auditMsg.set.Runs))
	}
}

// TestLoadAuditRegisterCmd_MalformedSurfacesErrorMsg proves a malformed record
// becomes a status message instead of aborting the board.
func TestLoadAuditRegisterCmd_MalformedSurfacesErrorMsg(t *testing.T) {
	root := t.TempDir()
	writeBoardAuditFile(t, root, "prompt.md", boardAuditPromptContent)
	writeBoardAuditFile(t, root, filepath.Join("runs", "2026-06-29-broken.md"), "no frontmatter")

	msg := loadAuditRegisterCmd(root, dataAuditLoader{})()

	errMsg, ok := msg.(errorMsg)
	if !ok {
		t.Fatalf("msg type = %T, want errorMsg on malformed record", msg)
	}
	if errMsg.message == "" {
		t.Error("errorMsg.message = empty, want actionable error text")
	}
}

// TestReloadTasksRefreshesAuditRegister proves the aggregate reload path
// (fileChangeMsg, ctrl+r, and audit-file watch events) carries fresh
// audit-register data into the model.
func TestReloadTasksRefreshesAuditRegister(t *testing.T) {
	root := t.TempDir()
	testutil.WriteRouter(t, root, "task-building", "v1", "E01-audit", "", "test")
	writeTask(t, root, "v1", "E01-audit", "T001-load", data.ColumnPlanned)
	writeCompleteAuditTree(t, root)

	msg := reloadTasksWithMessage(root, defaultModelDependencies(), "Refreshed")()

	reload, ok := msg.(reloadMsg)
	if !ok {
		t.Fatalf("msg type = %T, want reloadMsg", msg)
	}
	if !reload.audit.Prompt.Available {
		t.Error("reload.audit.Prompt.Available = false, want reload to refresh audit data")
	}
	if len(reload.audit.Findings) != 1 {
		t.Errorf("reload.audit.Findings = %d, want 1", len(reload.audit.Findings))
	}
}

// TestReloadTasksToleratesMalformedAudit proves a malformed audit record does
// not block a task/release refresh: audit degrades to an empty set while the
// reload still succeeds.
func TestReloadTasksToleratesMalformedAudit(t *testing.T) {
	root := t.TempDir()
	testutil.WriteRouter(t, root, "task-building", "v1", "E01-audit", "", "test")
	writeTask(t, root, "v1", "E01-audit", "T001-load", data.ColumnPlanned)
	writeBoardAuditFile(t, root, filepath.Join("runs", "2026-06-29-broken.md"), "no frontmatter")

	msg := reloadTasksWithMessage(root, defaultModelDependencies(), "Refreshed")()

	reload, ok := msg.(reloadMsg)
	if !ok {
		t.Fatalf("msg type = %T, want reloadMsg (malformed audit must not fail reload)", msg)
	}
	if len(reload.tasks) != 1 {
		t.Errorf("reload.tasks = %d, want 1 task still loaded", len(reload.tasks))
	}
	if reload.audit.Prompt.Available || len(reload.audit.Findings) != 0 {
		t.Error("malformed audit should degrade to an empty set on reload")
	}
}

// TestNewProjectModelLoadsAuditRegister proves board startup loads audit data
// into the model when an audit/ tree exists.
func TestNewProjectModelLoadsAuditRegister(t *testing.T) {
	projectRoot := t.TempDir()
	savepointRoot := filepath.Join(projectRoot, ".savepoint")
	testutil.WriteRouter(t, savepointRoot, "task-building", "v1", "E01-audit", "", "test")
	writeTask(t, savepointRoot, "v1", "E01-audit", "T001-load", data.ColumnPlanned)
	writeCompleteAuditTree(t, savepointRoot)

	model, err := newProjectModel(projectRoot, "", "")
	if err != nil {
		t.Fatalf("newProjectModel() error = %v", err)
	}
	if model.Watcher != nil {
		t.Cleanup(func() { model.Watcher.Close() })
	}

	if !model.Audit.Prompt.Available {
		t.Error("model.Audit.Prompt.Available = false, want startup to load audit data")
	}
	if len(model.Audit.Findings) != 1 || model.Audit.Findings[0].ID != "F001" {
		t.Errorf("model.Audit.Findings = %v, want one F001", model.Audit.Findings)
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
