package doctor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/testutil"
)

func TestRunV2ChecksDoesNotReadLegacyAuditRecords(t *testing.T) {
	root := t.TempDir()
	writeCompleteV2Project(t, root)
	auditPath := filepath.Join(root, "audit", "findings", "F001-legacy.md")
	testutil.WriteFile(t, auditPath, "not V2 audit data and intentionally malformed")

	report := RunV2Checks(root)
	if len(report.AuditRegister) != 0 {
		t.Fatalf("AuditRegister = %v, want no legacy audit check on the live V2 path", report.AuditRegister)
	}
	if report.HasProblems() {
		t.Fatalf("RunV2Checks() reported a problem from the legacy audit tree: %v", report.HealthFindings())
	}
}

func TestRunV2ChecksUsesStrictRouterReader(t *testing.T) {
	root := t.TempDir()
	writeCompleteV2Project(t, root)
	routerPath := filepath.Join(root, "router.md")
	if err := os.WriteFile(routerPath, []byte("## Current state\n\n```yaml\nstate: audit-pending\nrelease: v1\nepic: E01\nnext_action: legacy\n```\n"), 0644); err != nil {
		t.Fatal(err)
	}

	report := RunV2Checks(root)
	if report.RouterCheck == nil {
		t.Fatal("RouterCheck = nil, want strict V2 rejection of legacy router fields")
	}
	if !strings.Contains(report.RouterCheck.Error(), "V2") && !strings.Contains(report.RouterCheck.Error(), "unknown field") {
		t.Fatalf("RouterCheck = %v, want a named V2 router diagnostic", report.RouterCheck)
	}
}
