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
		".savepoint":                            &fstest.MapFile{Mode: fs.ModeDir | 0755},
		".savepoint/PRD.md":                     &fstest.MapFile{Data: []byte("# PRD")},
		".savepoint/Design.md":                  &fstest.MapFile{Data: []byte("# Design")},
		"agent-skills/savepoint-audit/SKILL.md": &fstest.MapFile{Data: []byte("# Audit Skill")},
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

	skillDir := filepath.Join(target, "agent-skills", "savepoint-audit")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatal(err)
	}
	oldContent := "# Old Skill Content"
	testutil.WriteFile(t, filepath.Join(skillDir, "SKILL.md"), oldContent)

	newContent := "# New Audit Skill"
	templates := fstest.MapFS{
		"agent-skills/savepoint-audit/SKILL.md": &fstest.MapFile{Data: []byte(newContent)},
	}

	report, err := UpgradeProjectAssets(templates, target, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	found := false
	for _, e := range report.Actions {
		if e.Path == "agent-skills/savepoint-audit/SKILL.md" {
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

	skillDir := filepath.Join(target, "agent-skills", "savepoint-audit")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatal(err)
	}
	content := "# Same Content"
	testutil.WriteFile(t, filepath.Join(skillDir, "SKILL.md"), content)

	templates := fstest.MapFS{
		"agent-skills/savepoint-audit/SKILL.md": &fstest.MapFile{Data: []byte(content)},
	}

	report, err := UpgradeProjectAssets(templates, target, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	for _, e := range report.Actions {
		if e.Path == "agent-skills/savepoint-audit/SKILL.md" {
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

	skillDir := filepath.Join(target, "agent-skills", "savepoint-audit")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatal(err)
	}
	oldContent := "# Old Content"
	testutil.WriteFile(t, filepath.Join(skillDir, "SKILL.md"), oldContent)

	templates := fstest.MapFS{
		"agent-skills/savepoint-audit/SKILL.md": &fstest.MapFile{Data: []byte("# New Content")},
	}

	report, err := UpgradeProjectAssets(templates, target, true, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() dry-run error = %v", err)
	}

	found := false
	for _, e := range report.Actions {
		if e.Path == "agent-skills/savepoint-audit/SKILL.md" {
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
	skillDir := filepath.Join(target, "agent-skills", "savepoint-audit")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatal(err)
	}
	testutil.WriteFile(t, filepath.Join(skillDir, "SKILL.md"), content)

	templates := fstest.MapFS{
		"agent-skills/savepoint-audit/SKILL.md": &fstest.MapFile{Data: []byte(content)},
	}

	report, err := UpgradeProjectAssets(templates, target, true, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() dry-run error = %v", err)
	}

	for _, e := range report.Actions {
		if e.Path == "agent-skills/savepoint-audit/SKILL.md" {
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
		"agent-skills/savepoint-audit/SKILL.md": &fstest.MapFile{Data: []byte("# New Skill")},
	}

	report, err := UpgradeProjectAssets(templates, target, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	found := false
	for _, e := range report.Actions {
		if e.Path == "agent-skills/savepoint-audit/SKILL.md" {
			found = true
			if e.Action != ActionUpdated {
				t.Errorf("action = %v, want updated", e.Action)
			}
		}
	}
	if !found {
		t.Fatal("skill path not in report")
	}

	if _, err := os.Stat(filepath.Join(target, "agent-skills", "savepoint-audit", "SKILL.md")); err != nil {
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

func TestUpgradeReport_format(t *testing.T) {
	r := &UpgradeReport{
		Actions: []UpgradeEntry{
			{Path: "agent-skills/a/SKILL.md", Action: ActionUpdated},
			{Path: "AGENTS.md", Action: ActionMerged},
			{Path: ".savepoint/PRD.md", Action: ActionSkipped},
			{Path: "agent-skills/b/SKILL.md", Action: ActionUnchanged},
		},
	}

	output := r.Format()
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
