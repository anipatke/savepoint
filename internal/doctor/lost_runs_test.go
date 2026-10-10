package doctor

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/opencode/savepoint/internal/testutil"
)

func writeRunRecord(t *testing.T, root, name, state string, pid int, updated time.Time) {
	t.Helper()
	testutil.WriteFile(t, filepath.Join(root, "runs", name, "run.yml"),
		"provider: claude\ntask: T-001\nstate: "+state+"\npid: "+strconv.Itoa(pid)+"\nupdated_at: '"+updated.UTC().Format(time.RFC3339)+"'\n")
}

func TestRunV2Checks_namesALostRunAsAdviceOnly(t *testing.T) {
	root := filepath.Join(t.TempDir(), ".savepoint")
	writeCompleteV2Project(t, root)
	writeRunRecord(t, root, "20261011T100000Z-T-001", "running", 0x7ffffff0, time.Now().Add(-time.Hour))
	writeRunRecord(t, root, "20261011T110000Z-T-001", "finished", 0x7ffffff0, time.Now().Add(-time.Hour))

	report := RunV2Checks(root)
	if report.HasProblems() {
		t.Fatalf("a lost run must be advice, not a structural problem: %v", report.Project)
	}
	out := report.Format()
	if !strings.Contains(out, "[runner_lost] run 20261011T100000Z-T-001") || strings.Contains(out, "20261011T110000Z-T-001") {
		t.Errorf("report should name exactly the lost run:\n%s", out)
	}
}

func TestRunV2Checks_aProjectWithoutRunsPrintsNothingNew(t *testing.T) {
	root := filepath.Join(t.TempDir(), ".savepoint")
	writeCompleteV2Project(t, root)
	if out := RunV2Checks(root).Format(); strings.Contains(out, "runner_lost") {
		t.Errorf("unexpected runner_lost:\n%s", out)
	}
	if err := os.MkdirAll(filepath.Join(root, "runs"), 0o755); err != nil {
		t.Fatal(err)
	}
	if out := RunV2Checks(root).Format(); strings.Contains(out, "runner_lost") {
		t.Errorf("unexpected runner_lost for an empty runs directory:\n%s", out)
	}
}
