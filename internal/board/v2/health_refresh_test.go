package v2

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/opencode/savepoint/internal/codehealth"
)

func TestRefreshHealthReportsBlockedReportAsWarning(t *testing.T) {
	for _, blocked := range []bool{false, true} {
		dir := t.TempDir()
		if out, err := exec.Command("git", "-C", dir, "init", "-q").CombinedOutput(); err != nil {
			t.Fatalf("git init: %v %s", err, out)
		}
		if err := os.MkdirAll(filepath.Join(dir, ".savepoint"), 0o755); err != nil {
			t.Fatal(err)
		}
		if _, err := codehealth.NewStore(dir).SaveConfig(codehealth.Config{Version: codehealth.ConfigVersion}); err != nil {
			t.Fatal(err)
		}
		if blocked {
			if err := os.MkdirAll(filepath.Join(dir, ".savepoint", "health", "report.md"), 0o755); err != nil {
				t.Fatal(err)
			}
		}
		reportErr, err := refreshHealth(context.Background(), filepath.Join(dir, ".savepoint"), nil)
		if err != nil {
			t.Fatalf("blocked=%v: refresh failed: %v", blocked, err)
		}
		if (reportErr != nil) != blocked {
			t.Fatalf("blocked=%v: reportErr = %v", blocked, reportErr)
		}
		if snaps, _ := codehealth.NewStore(dir).LoadSnapshots(); len(snaps) != 1 {
			t.Fatalf("blocked=%v: snapshots = %d, want 1", blocked, len(snaps))
		}
	}
}
