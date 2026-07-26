package init

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/opencode/savepoint/internal/testutil"
)

func TestUpgradeProjectAssets_requiresSavepointProject(t *testing.T) {
	target := t.TempDir()
	templates := fstest.MapFS{}

	_, err := UpgradeProjectAssets(templates, target, false, false)
	if err == nil {
		t.Fatal("expected error for non-savepoint project")
	}
	if !strings.Contains(err.Error(), "not a Savepoint project") {
		t.Fatalf("error = %q, want 'not a Savepoint project'", err.Error())
	}
}

func TestUpgradeProjectAssets_requiresExistingDirectory(t *testing.T) {
	target := filepath.Join(t.TempDir(), "nonexistent")
	templates := fstest.MapFS{}

	_, err := UpgradeProjectAssets(templates, target, false, false)
	if err == nil {
		t.Fatal("expected error for missing directory")
	}
	if !strings.Contains(err.Error(), "does not exist") {
		t.Fatalf("error = %q, want 'does not exist'", err.Error())
	}
}

func TestUpgradeProjectAssets_skipsSavepointDir(t *testing.T) {
	target := t.TempDir()
	savepointDir := filepath.Join(target, ".savepoint")
	if err := os.MkdirAll(savepointDir, 0755); err != nil {
		t.Fatal(err)
	}

	templates := fstest.MapFS{
		".savepoint":           &fstest.MapFile{Mode: fs.ModeDir | 0755},
		".savepoint/PRD.md":    &fstest.MapFile{Data: []byte("# PRD")},
		".savepoint/Design.md": &fstest.MapFile{Data: []byte("# Design")},
		"agent-skills/savepoint-audit-epic/SKILL.md": &fstest.MapFile{Data: []byte("# Audit Skill")},
	}

	report, err := UpgradeProjectAssets(templates, target, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	for _, e := range report.Actions {
		if strings.HasPrefix(e.Path, ".savepoint/") {
			if e.Action != ActionSkipped {
				t.Errorf("path %s = %v, want skipped", e.Path, e.Action)
			}
		}
	}
}

func TestUpgradeProjectAssets_updatesAgentSkills(t *testing.T) {
	target := t.TempDir()
	savepointDir := filepath.Join(target, ".savepoint")
	if err := os.MkdirAll(savepointDir, 0755); err != nil {
		t.Fatal(err)
	}

	skillDir := filepath.Join(target, "agent-skills", "savepoint-audit-epic")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatal(err)
	}
	oldContent := "# Old Skill Content"
	testutil.WriteFile(t, filepath.Join(skillDir, "SKILL.md"), oldContent)

	newContent := "# New Audit Skill"
	templates := fstest.MapFS{
		"agent-skills/savepoint-audit-epic/SKILL.md": &fstest.MapFile{Data: []byte(newContent)},
	}

	report, err := UpgradeProjectAssets(templates, target, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	found := false
	for _, e := range report.Actions {
		if e.Path == "agent-skills/savepoint-audit-epic/SKILL.md" {
			found = true
			if e.Action != ActionUpdated {
				t.Errorf("action = %v, want updated", e.Action)
			}
		}
	}
	if !found {
		t.Fatal("skill path not in report")
	}

	data, err := os.ReadFile(filepath.Join(skillDir, "SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != newContent {
		t.Fatalf("content = %q, want %q", string(data), newContent)
	}
}

func TestUpgradeProjectAssets_skillIdempotent(t *testing.T) {
	target := t.TempDir()
	savepointDir := filepath.Join(target, ".savepoint")
	if err := os.MkdirAll(savepointDir, 0755); err != nil {
		t.Fatal(err)
	}

	skillDir := filepath.Join(target, "agent-skills", "savepoint-audit-epic")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatal(err)
	}
	content := "# Same Content"
	testutil.WriteFile(t, filepath.Join(skillDir, "SKILL.md"), content)

	templates := fstest.MapFS{
		"agent-skills/savepoint-audit-epic/SKILL.md": &fstest.MapFile{Data: []byte(content)},
	}

	report, err := UpgradeProjectAssets(templates, target, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	for _, e := range report.Actions {
		if e.Path == "agent-skills/savepoint-audit-epic/SKILL.md" {
			if e.Action != ActionUnchanged {
				t.Errorf("action = %v, want unchanged", e.Action)
			}
		}
	}
}

func TestUpgradeProjectAssets_mergesAgentsMd(t *testing.T) {
	target := t.TempDir()
	savepointDir := filepath.Join(target, ".savepoint")
	if err := os.MkdirAll(savepointDir, 0755); err != nil {
		t.Fatal(err)
	}

	existingGuide := "# My Guide\n\nExisting user content."
	testutil.WriteFile(t, filepath.Join(target, "AGENTS.md"), existingGuide)

	templates := fstest.MapFS{
		"AGENTS.md": &fstest.MapFile{Data: []byte("# Savepoint Managed Content")},
	}

	report, err := UpgradeProjectAssets(templates, target, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	found := false
	for _, e := range report.Actions {
		if e.Path == "AGENTS.md" {
			found = true
			if e.Action != ActionMerged {
				t.Errorf("action = %v, want merged", e.Action)
			}
		}
	}
	if !found {
		t.Fatal("AGENTS.md not in report")
	}

	data, err := os.ReadFile(filepath.Join(target, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(data)
	if !strings.Contains(got, "My Guide") {
		t.Errorf("missing user content: %q", got)
	}
	if !strings.Contains(got, managedBegin) {
		t.Errorf("missing managed block: %q", got)
	}
	if !strings.Contains(got, "Savepoint Managed Content") {
		t.Errorf("missing managed content: %q", got)
	}
}

func TestUpgradeProjectAssets_agentsMdIdempotent(t *testing.T) {
	target := t.TempDir()
	savepointDir := filepath.Join(target, ".savepoint")
	if err := os.MkdirAll(savepointDir, 0755); err != nil {
		t.Fatal(err)
	}

	templates := fstest.MapFS{
		"AGENTS.md": &fstest.MapFile{Data: []byte("# Savepoint Instructions")},
	}

	for range 2 {
		_, err := UpgradeProjectAssets(templates, target, false, false)
		if err != nil {
			t.Fatalf("UpgradeProjectAssets() run error = %v", err)
		}
	}

	data, err := os.ReadFile(filepath.Join(target, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(data)
	count := strings.Count(got, managedBegin)
	if count != 1 {
		t.Errorf("AGENTS.md has %d managed begin markers, want 1: %q", count, got)
	}
}

func TestUpgradeProjectAssets_dryRunDoesNotWrite(t *testing.T) {
	target := t.TempDir()
	savepointDir := filepath.Join(target, ".savepoint")
	if err := os.MkdirAll(savepointDir, 0755); err != nil {
		t.Fatal(err)
	}

	skillDir := filepath.Join(target, "agent-skills", "savepoint-audit-epic")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatal(err)
	}
	oldContent := "# Old Content"
	testutil.WriteFile(t, filepath.Join(skillDir, "SKILL.md"), oldContent)

	templates := fstest.MapFS{
		"agent-skills/savepoint-audit-epic/SKILL.md": &fstest.MapFile{Data: []byte("# New Content")},
	}

	report, err := UpgradeProjectAssets(templates, target, true, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() dry-run error = %v", err)
	}

	found := false
	for _, e := range report.Actions {
		if e.Path == "agent-skills/savepoint-audit-epic/SKILL.md" {
			found = true
			if e.Action != ActionUpdated {
				t.Errorf("dry-run action = %v, want updated", e.Action)
			}
		}
	}
	if !found {
		t.Fatal("skill path not in dry-run report")
	}

	data, err := os.ReadFile(filepath.Join(skillDir, "SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != oldContent {
		t.Fatalf("dry-run should not write: got %q, want %q", string(data), oldContent)
	}
}

func TestUpgradeProjectAssets_dryRunReportsUnchanged(t *testing.T) {
	target := t.TempDir()
	savepointDir := filepath.Join(target, ".savepoint")
	if err := os.MkdirAll(savepointDir, 0755); err != nil {
		t.Fatal(err)
	}

	content := "# Same Content"
	skillDir := filepath.Join(target, "agent-skills", "savepoint-audit-epic")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatal(err)
	}
	testutil.WriteFile(t, filepath.Join(skillDir, "SKILL.md"), content)

	templates := fstest.MapFS{
		"agent-skills/savepoint-audit-epic/SKILL.md": &fstest.MapFile{Data: []byte(content)},
	}

	report, err := UpgradeProjectAssets(templates, target, true, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() dry-run error = %v", err)
	}

	for _, e := range report.Actions {
		if e.Path == "agent-skills/savepoint-audit-epic/SKILL.md" {
			if e.Action != ActionUnchanged {
				t.Errorf("dry-run action = %v, want unchanged for same content", e.Action)
			}
		}
	}
}

func TestUpgradeProjectAssets_skipsNonAllowlistedFiles(t *testing.T) {
	target := t.TempDir()
	savepointDir := filepath.Join(target, ".savepoint")
	if err := os.MkdirAll(savepointDir, 0755); err != nil {
		t.Fatal(err)
	}

	templates := fstest.MapFS{
		"some-other-file.md": &fstest.MapFile{Data: []byte("content")},
		"README.md":          &fstest.MapFile{Data: []byte("# Readme")},
	}

	report, err := UpgradeProjectAssets(templates, target, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	for _, e := range report.Actions {
		if e.Action != ActionSkipped {
			t.Errorf("path %s = %v, want skipped", e.Path, e.Action)
		}
	}
}

func TestUpgradeProjectAssets_skipsPromptTemplates(t *testing.T) {
	target := t.TempDir()
	savepointDir := filepath.Join(target, ".savepoint")
	if err := os.MkdirAll(savepointDir, 0755); err != nil {
		t.Fatal(err)
	}

	templates := fstest.MapFS{
		"templates/prompts/task-building.prompt.md": &fstest.MapFile{Data: []byte("stale phase prompt")},
		"prompts/task-building.prompt.md":           &fstest.MapFile{Data: []byte("stale phase prompt")},
	}

	report, err := UpgradeProjectAssets(templates, target, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	if len(report.Actions) != len(templates) {
		t.Fatalf("actions = %d, want %d", len(report.Actions), len(templates))
	}
	for _, e := range report.Actions {
		if e.Action != ActionSkipped {
			t.Errorf("path %s = %v, want skipped", e.Path, e.Action)
		}
		if _, err := os.Stat(filepath.Join(target, e.Path)); !os.IsNotExist(err) {
			t.Errorf("prompt path %s should not be written, stat err = %v", e.Path, err)
		}
	}
}

func TestUpgradeProjectAssets_createsMissingSkillFile(t *testing.T) {
	target := t.TempDir()
	savepointDir := filepath.Join(target, ".savepoint")
	if err := os.MkdirAll(savepointDir, 0755); err != nil {
		t.Fatal(err)
	}

	templates := fstest.MapFS{
		"agent-skills/savepoint-audit-epic/SKILL.md": &fstest.MapFile{Data: []byte("# New Skill")},
	}

	report, err := UpgradeProjectAssets(templates, target, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	found := false
	for _, e := range report.Actions {
		if e.Path == "agent-skills/savepoint-audit-epic/SKILL.md" {
			found = true
			if e.Action != ActionUpdated {
				t.Errorf("action = %v, want updated", e.Action)
			}
		}
	}
	if !found {
		t.Fatal("skill path not in report")
	}

	if _, err := os.Stat(filepath.Join(target, "agent-skills", "savepoint-audit-epic", "SKILL.md")); err != nil {
		t.Errorf("skill file not created: %v", err)
	}
}

func TestUpgradeProjectAssets_casingVariantAgentsMd(t *testing.T) {
	target := t.TempDir()
	savepointDir := filepath.Join(target, ".savepoint")
	if err := os.MkdirAll(savepointDir, 0755); err != nil {
		t.Fatal(err)
	}

	variantPath := filepath.Join(target, "Agents.MD")
	testutil.WriteFile(t, variantPath, "# My Guide")

	templates := fstest.MapFS{
		"AGENTS.md": &fstest.MapFile{Data: []byte("# Savepoint Instructions")},
	}

	report, err := UpgradeProjectAssets(templates, target, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	found := false
	for _, e := range report.Actions {
		if e.Path == "AGENTS.md" {
			found = true
			if e.Action != ActionMerged {
				t.Errorf("action = %v, want merged", e.Action)
			}
		}
	}
	if !found {
		t.Fatal("AGENTS.md not in report")
	}

	data, err := os.ReadFile(variantPath)
	if err != nil {
		t.Fatal(err)
	}
	got := string(data)
	if !strings.Contains(got, "# My Guide") {
		t.Errorf("missing existing content: %q", got)
	}
	if !strings.Contains(got, managedBegin) {
		t.Errorf("missing managed block: %q", got)
	}

	entries, _ := os.ReadDir(target)
	guideCount := 0
	for _, e := range entries {
		if strings.ToLower(e.Name()) == "agents.md" {
			guideCount++
		}
	}
	if guideCount != 1 {
		t.Errorf("expected 1 agent guide file, found %d", guideCount)
	}
}

func TestUpgradeProjectAssets_dryRunUsesCasingVariantAgentGuide(t *testing.T) {
	target := t.TempDir()
	savepointDir := filepath.Join(target, ".savepoint")
	if err := os.MkdirAll(savepointDir, 0755); err != nil {
		t.Fatal(err)
	}

	variantPath := filepath.Join(target, "Agents.MD")
	testutil.WriteFile(t, variantPath, "# My Guide")

	templates := fstest.MapFS{
		"AGENTS.md": &fstest.MapFile{Data: []byte("# Savepoint Instructions")},
	}

	report, err := UpgradeProjectAssets(templates, target, true, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() dry-run error = %v", err)
	}

	for _, e := range report.Actions {
		if e.Path == "AGENTS.md" {
			if e.Action != ActionMerged {
				t.Errorf("dry-run action = %v, want merged", e.Action)
			}
			return
		}
	}
	t.Fatal("AGENTS.md not in dry-run report")
}

func TestUpgradeProjectAssets_multipleSkills(t *testing.T) {
	target := t.TempDir()
	savepointDir := filepath.Join(target, ".savepoint")
	if err := os.MkdirAll(savepointDir, 0755); err != nil {
		t.Fatal(err)
	}

	templates := fstest.MapFS{
		"agent-skills/skill-a/SKILL.md": &fstest.MapFile{Data: []byte("# Skill A")},
		"agent-skills/skill-b/SKILL.md": &fstest.MapFile{Data: []byte("# Skill B")},
	}

	report, err := UpgradeProjectAssets(templates, target, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	updated := 0
	for _, e := range report.Actions {
		if e.Action == ActionUpdated {
			updated++
		}
	}
	if updated != 2 {
		t.Errorf("expected 2 updated, got %d", updated)
	}
}

func auditTemplates() fstest.MapFS {
	return fstest.MapFS{
		".savepoint/audit/prompt.md":          &fstest.MapFile{Data: []byte("# Audit Prompt")},
		".savepoint/audit/register.md":        &fstest.MapFile{Data: []byte("# Audit Register")},
		".savepoint/audit/findings/README.md": &fstest.MapFile{Data: []byte("# Findings")},
		".savepoint/audit/runs/README.md":     &fstest.MapFile{Data: []byte("# Runs")},
	}
}

func actionFor(report *UpgradeReport, path string) (UpgradeAction, bool) {
	for _, e := range report.Actions {
		if e.Path == path {
			return e.Action, true
		}
	}
	return "", false
}

func TestUpgradeProjectAssets_addsMissingAuditAssets(t *testing.T) {
	target := t.TempDir()
	if err := os.MkdirAll(filepath.Join(target, ".savepoint"), 0755); err != nil {
		t.Fatal(err)
	}

	templates := auditTemplates()

	report, err := UpgradeProjectAssets(templates, target, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	for path, file := range templates {
		action, found := actionFor(report, path)
		if !found {
			t.Fatalf("audit asset %s not in report", path)
		}
		if action != ActionUpdated {
			t.Errorf("audit asset %s action = %v, want updated", path, action)
		}
		data, err := os.ReadFile(filepath.Join(target, filepath.FromSlash(path)))
		if err != nil {
			t.Errorf("audit asset %s not written: %v", path, err)
			continue
		}
		if string(data) != string(file.Data) {
			t.Errorf("audit asset %s = %q, want %q", path, string(data), string(file.Data))
		}
	}
}

func TestUpgradeProjectAssets_preservesEditedAuditAssets(t *testing.T) {
	target := t.TempDir()
	if err := os.MkdirAll(filepath.Join(target, ".savepoint"), 0755); err != nil {
		t.Fatal(err)
	}

	editedRegister := filepath.Join(target, ".savepoint", "audit", "register.md")
	userContent := "# User-maintained register\n\nF001 disposition work."
	if err := os.MkdirAll(filepath.Dir(editedRegister), 0755); err != nil {
		t.Fatal(err)
	}
	testutil.WriteFile(t, editedRegister, userContent)

	templates := auditTemplates()

	report, err := UpgradeProjectAssets(templates, target, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	action, found := actionFor(report, ".savepoint/audit/register.md")
	if !found {
		t.Fatal("register not in report")
	}
	if action != ActionUnchanged {
		t.Errorf("edited register action = %v, want unchanged", action)
	}

	data, err := os.ReadFile(editedRegister)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != userContent {
		t.Errorf("edited register overwritten: got %q, want %q", string(data), userContent)
	}

	// Missing siblings are still added alongside the preserved file.
	if _, err := os.Stat(filepath.Join(target, ".savepoint", "audit", "prompt.md")); err != nil {
		t.Errorf("missing prompt not added: %v", err)
	}
}

func TestUpgradeProjectAssets_pristineAuditAssetsUnchanged(t *testing.T) {
	target := t.TempDir()
	if err := os.MkdirAll(filepath.Join(target, ".savepoint"), 0755); err != nil {
		t.Fatal(err)
	}

	templates := auditTemplates()

	// First upgrade writes the assets; second sees them pristine.
	if _, err := UpgradeProjectAssets(templates, target, false, false); err != nil {
		t.Fatalf("first UpgradeProjectAssets() error = %v", err)
	}

	report, err := UpgradeProjectAssets(templates, target, false, false)
	if err != nil {
		t.Fatalf("second UpgradeProjectAssets() error = %v", err)
	}

	for path := range templates {
		action, found := actionFor(report, path)
		if !found {
			t.Fatalf("audit asset %s not in report", path)
		}
		if action != ActionUnchanged {
			t.Errorf("pristine audit asset %s action = %v, want unchanged", path, action)
		}
	}
}

func TestUpgradeProjectAssets_dryRunDoesNotAddAuditAssets(t *testing.T) {
	target := t.TempDir()
	if err := os.MkdirAll(filepath.Join(target, ".savepoint"), 0755); err != nil {
		t.Fatal(err)
	}

	templates := auditTemplates()

	report, err := UpgradeProjectAssets(templates, target, true, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() dry-run error = %v", err)
	}

	for path := range templates {
		action, found := actionFor(report, path)
		if !found {
			t.Fatalf("audit asset %s not in dry-run report", path)
		}
		if action != ActionUpdated {
			t.Errorf("dry-run audit asset %s action = %v, want updated", path, action)
		}
		if _, err := os.Stat(filepath.Join(target, filepath.FromSlash(path))); !os.IsNotExist(err) {
			t.Errorf("dry-run should not write %s, stat err = %v", path, err)
		}
	}
}

func policyTemplates() fstest.MapFS {
	return fstest.MapFS{
		".savepoint/Guardrails.md":   &fstest.MapFile{Data: []byte("# Guardrails for {{PROJECT_NAME}}\n\nSTYLE-01")},
		".savepoint/Health-Check.md": &fstest.MapFile{Data: []byte("# Health Check\n\n## Quick Check")},
		".savepoint/PRD.md":          &fstest.MapFile{Data: []byte("# PRD")},
		".savepoint/router.md":       &fstest.MapFile{Data: []byte("# Router")},
	}
}

func policyAssetPaths() []string {
	return []string{".savepoint/Guardrails.md", ".savepoint/Health-Check.md"}
}

// renderedPolicyTemplate is the template content as upgrade writes it: project
// name interpolated exactly as a fresh scaffold does.
func renderedPolicyTemplate(templates fstest.MapFS, path, target string) string {
	return interpolate(string(templates[path].Data), ProjectNameFromDir(target))
}

func TestUpgradeProjectAssets_installsMissingPolicyAssets(t *testing.T) {
	target := t.TempDir()
	if err := os.MkdirAll(filepath.Join(target, ".savepoint"), 0755); err != nil {
		t.Fatal(err)
	}

	templates := policyTemplates()

	report, err := UpgradeProjectAssets(templates, target, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	for _, path := range policyAssetPaths() {
		action, found := actionFor(report, path)
		if !found {
			t.Fatalf("policy asset %s not in report", path)
		}
		if action != ActionInstalled {
			t.Errorf("policy asset %s action = %v, want installed", path, action)
		}

		data, err := os.ReadFile(filepath.Join(target, filepath.FromSlash(path)))
		if err != nil {
			t.Errorf("policy asset %s not written: %v", path, err)
			continue
		}
		if want := renderedPolicyTemplate(templates, path, target); string(data) != want {
			t.Errorf("policy asset %s = %q, want %q", path, string(data), want)
		}
	}

	// Widening the gate must not widen it past the allowlist.
	for _, path := range []string{".savepoint/PRD.md", ".savepoint/router.md"} {
		action, found := actionFor(report, path)
		if !found {
			t.Fatalf("%s not in report", path)
		}
		if action != ActionSkipped {
			t.Errorf("non-policy .savepoint asset %s action = %v, want skipped", path, action)
		}
		if _, err := os.Stat(filepath.Join(target, filepath.FromSlash(path))); !os.IsNotExist(err) {
			t.Errorf("non-policy .savepoint asset %s written, stat err = %v", path, err)
		}
	}
}

func TestUpgradeProjectAssets_policyAssetBranches(t *testing.T) {
	templates := policyTemplates()

	cases := []struct {
		name       string
		existing   string // pre-existing project content; empty means the file is missing
		dryRun     bool
		wantAction UpgradeAction
	}{
		{name: "missing", wantAction: ActionInstalled},
		{name: "present pristine", existing: "template", wantAction: ActionUnchanged},
		{name: "user modified", existing: "# Locally tailored policy\n", wantAction: ActionUnchanged},
		{name: "missing dry run", dryRun: true, wantAction: ActionInstalled},
		{name: "present dry run", existing: "# Locally tailored policy\n", dryRun: true, wantAction: ActionUnchanged},
	}

	for _, c := range cases {
		for _, path := range policyAssetPaths() {
			t.Run(c.name+" "+path, func(t *testing.T) {
				target := t.TempDir()
				if err := os.MkdirAll(filepath.Join(target, ".savepoint"), 0755); err != nil {
					t.Fatal(err)
				}

				existing := c.existing
				if existing == "template" {
					existing = renderedPolicyTemplate(templates, path, target)
				}
				targetPath := filepath.Join(target, filepath.FromSlash(path))
				if existing != "" {
					testutil.WriteFile(t, targetPath, existing)
				}

				report, err := UpgradeProjectAssets(templates, target, c.dryRun, false)
				if err != nil {
					t.Fatalf("UpgradeProjectAssets() error = %v", err)
				}

				action, found := actionFor(report, path)
				if !found {
					t.Fatalf("policy asset %s not in report", path)
				}
				if action != c.wantAction {
					t.Errorf("policy asset %s action = %v, want %v", path, action, c.wantAction)
				}

				data, readErr := os.ReadFile(targetPath)
				switch {
				case existing != "":
					// Existing content, edited or not, must survive byte-identical.
					if readErr != nil {
						t.Fatalf("existing policy asset %s unreadable: %v", path, readErr)
					}
					if string(data) != existing {
						t.Errorf("policy asset %s = %q, want unchanged %q", path, string(data), existing)
					}
				case c.dryRun:
					if !os.IsNotExist(readErr) {
						t.Errorf("dry-run wrote %s, stat err = %v", path, readErr)
					}
				default:
					if readErr != nil {
						t.Fatalf("policy asset %s not installed: %v", path, readErr)
					}
				}
			})
		}
	}
}

func TestUpgradeProjectAssets_policyAssetsIdempotent(t *testing.T) {
	target := t.TempDir()
	if err := os.MkdirAll(filepath.Join(target, ".savepoint"), 0755); err != nil {
		t.Fatal(err)
	}

	templates := policyTemplates()

	if _, err := UpgradeProjectAssets(templates, target, false, false); err != nil {
		t.Fatalf("first UpgradeProjectAssets() error = %v", err)
	}

	before := map[string]string{}
	for _, path := range policyAssetPaths() {
		data, err := os.ReadFile(filepath.Join(target, filepath.FromSlash(path)))
		if err != nil {
			t.Fatalf("policy asset %s not installed by first run: %v", path, err)
		}
		before[path] = string(data)
	}

	report, err := UpgradeProjectAssets(templates, target, false, false)
	if err != nil {
		t.Fatalf("second UpgradeProjectAssets() error = %v", err)
	}

	for _, path := range policyAssetPaths() {
		action, found := actionFor(report, path)
		if !found {
			t.Fatalf("policy asset %s not in rerun report", path)
		}
		if action != ActionUnchanged {
			t.Errorf("rerun policy asset %s action = %v, want unchanged", path, action)
		}
		data, err := os.ReadFile(filepath.Join(target, filepath.FromSlash(path)))
		if err != nil {
			t.Fatalf("policy asset %s missing after rerun: %v", path, err)
		}
		if string(data) != before[path] {
			t.Errorf("rerun rewrote %s: got %q, want %q", path, string(data), before[path])
		}
	}
}

func TestUpgradeReport_format(t *testing.T) {
	r := &UpgradeReport{
		Actions: []UpgradeEntry{
			{Path: "agent-skills/a/SKILL.md", Action: ActionUpdated},
			{Path: "AGENTS.md", Action: ActionMerged},
			{Path: ".savepoint/PRD.md", Action: ActionSkipped},
			{Path: "agent-skills/b/SKILL.md", Action: ActionUnchanged},
			{Path: ".savepoint/Guardrails.md", Action: ActionInstalled},
			{Path: "agent-skills/c/SKILL.md", Action: ActionMigrated},
		},
	}

	output := r.Format()
	if !strings.Contains(output, "Installed: 1") {
		t.Errorf("missing installed count: %q", output)
	}
	if !strings.Contains(output, "Migrated: 1") {
		t.Errorf("missing migrated count: %q", output)
	}
	if !strings.Contains(output, "Updated: 1") {
		t.Errorf("missing updated count: %q", output)
	}
	if !strings.Contains(output, "Merged: 1") {
		t.Errorf("missing merged count: %q", output)
	}
	if !strings.Contains(output, "Skipped: 1") {
		t.Errorf("missing skipped count: %q", output)
	}
	if !strings.Contains(output, "Unchanged: 1") {
		t.Errorf("missing unchanged count: %q", output)
	}
}

func TestUpgradeReport_formatEmpty(t *testing.T) {
	r := &UpgradeReport{}
	output := r.Format()
	if !strings.Contains(output, "No assets to upgrade.") {
		t.Errorf("empty report = %q", output)
	}
}
