package init

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/data"
)

// TestV2ArtifactTemplatesDecodeThroughTypedContracts proves the provenance
// fields emitted by the public artifact templates are accepted and retained
// by the typed readers that own those schemas. A text-presence assertion alone
// would allow a template and decoder to drift while both tests still passed.
func TestV2ArtifactTemplatesDecodeThroughTypedContracts(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		designPath := filepath.Join(root, "savepoint-design", "SKILL.md")
		design := readProvenanceContractSource(t, designPath)
		taskContent := provenanceArtifactFence(t, design, "planned_by: {role: planner, session: planning-example}")

		task, err := data.DecodeTaskV2(designPath+"#task-template", taskContent)
		if err != nil {
			t.Errorf("%s: DecodeTaskV2(%s) error = %v", tree, designPath, err)
		} else if task.PlannedBy != (data.Actor{Role: data.ActorRolePlanner, Session: "planning-example"}) {
			t.Errorf("%s: Task PlannedBy = %+v, want planner/planning-example", tree, task.PlannedBy)
		}

		checkPath := filepath.Join(root, "savepoint-check", "SKILL.md")
		check := readProvenanceContractSource(t, checkPath)
		checkYAML := provenanceArtifactFence(t, check, "executed_session: build-001")
		checkYAML = strings.NewReplacer(
			"C-###", "C-001",
			"task|objective|release", "task",
			"T-###, O-###, or R-###", "T-001",
			"CLEAR|NEEDS WORK", "CLEAR",
		).Replace(checkYAML)
		checkContent := "---\n" + strings.TrimSpace(checkYAML) + "\n---\n\n# Check\n"

		decodedCheck, err := data.DecodeCheckV2(checkPath+"#check-template", checkContent)
		if err != nil {
			t.Errorf("%s: DecodeCheckV2(%s) error = %v", tree, checkPath, err)
		} else if decodedCheck.ExecutedSession != "build-001" {
			t.Errorf("%s: Check ExecutedSession = %q, want build-001", tree, decodedCheck.ExecutedSession)
		}
	}
}

func readProvenanceContractSource(t *testing.T, path string) string {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(content)
}

func provenanceArtifactFence(t *testing.T, content, marker string) string {
	t.Helper()

	markerAt := strings.Index(content, marker)
	if markerAt < 0 {
		t.Fatalf("artifact template does not contain marker %q", marker)
	}

	openingAt := strings.LastIndex(content[:markerAt], "```")
	if openingAt < 0 {
		t.Fatalf("artifact marker %q has no opening markdown fence", marker)
	}
	bodyOffset := strings.IndexByte(content[openingAt:], '\n')
	if bodyOffset < 0 {
		t.Fatalf("artifact opening fence for %q has no body", marker)
	}
	bodyStart := openingAt + bodyOffset + 1
	closingOffset := strings.Index(content[bodyStart:], "\n```")
	if closingOffset < 0 {
		t.Fatalf("artifact marker %q has no closing markdown fence", marker)
	}

	return content[bodyStart : bodyStart+closingOffset]
}
