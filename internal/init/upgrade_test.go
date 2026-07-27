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

	before := "# My Guide\n\nExisting user content.\n\n"
	after := "\n\nTrailing user note.\n"
	existingGuide := before + managedBegin + "\n# Old Managed Content\n" + managedEnd + after
	testutil.WriteFile(t, filepath.Join(target, "AGENTS.md"), existingGuide)

	templates := fstest.MapFS{
		"AGENTS.md": &fstest.MapFile{Data: []byte("# Savepoint Managed Content")},
	}

	report, err := UpgradeProjectAssets(templates, target, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	action, found := actionFor(report, "AGENTS.md")
	if !found {
		t.Fatal("AGENTS.md not in report")
	}
	if action != ActionMerged {
		t.Errorf("action = %v, want merged", action)
	}

	data, err := os.ReadFile(filepath.Join(target, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	// Everything outside the markers must survive byte for byte; only the
	// block between them is Savepoint's to rewrite.
	want := before + managedBegin + "\n# Savepoint Managed Content\n" + managedEnd + after
	if string(data) != want {
		t.Errorf("merged guide = %q, want %q", string(data), want)
	}
}

func TestUpgradeProjectAssets_absentAgentsMdWritesWholeFile(t *testing.T) {
	target := t.TempDir()
	savepointDir := filepath.Join(target, ".savepoint")
	if err := os.MkdirAll(savepointDir, 0755); err != nil {
		t.Fatal(err)
	}

	templates := fstest.MapFS{
		"AGENTS.md": &fstest.MapFile{Data: []byte("# Savepoint Instructions")},
	}

	report, err := UpgradeProjectAssets(templates, target, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	action, found := actionFor(report, "AGENTS.md")
	if !found {
		t.Fatal("AGENTS.md not in report")
	}
	if action != ActionUpdated {
		t.Errorf("action = %v, want updated", action)
	}

	data, err := os.ReadFile(filepath.Join(target, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	want := managedBegin + "\n# Savepoint Instructions\n" + managedEnd + "\n"
	if string(data) != want {
		t.Errorf("guide = %q, want %q", string(data), want)
	}
}

func TestUpgradeProjectAssets_unmarkedAgentsMdConflicts(t *testing.T) {
	target := t.TempDir()
	savepointDir := filepath.Join(target, ".savepoint")
	if err := os.MkdirAll(savepointDir, 0755); err != nil {
		t.Fatal(err)
	}

	guidePath := filepath.Join(target, "AGENTS.md")
	existingGuide := "# My Guide\n\nExisting user content.\n"
	testutil.WriteFile(t, guidePath, existingGuide)

	templates := fstest.MapFS{
		"AGENTS.md": &fstest.MapFile{Data: []byte("# Savepoint Instructions")},
	}

	report, err := UpgradeProjectAssets(templates, target, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	entry, found := entryFor(report, "AGENTS.md")
	if !found {
		t.Fatal("AGENTS.md not in report")
	}
	if entry.Action != ActionConflict {
		t.Errorf("action = %v, want conflict", entry.Action)
	}
	if entry.Note != noteConflict {
		t.Errorf("note = %q, want %q", entry.Note, noteConflict)
	}

	data, err := os.ReadFile(guidePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != existingGuide {
		t.Errorf("guide changed: got %q, want %q", string(data), existingGuide)
	}

	incoming, err := os.ReadFile(guidePath + incomingSuffix)
	if err != nil {
		t.Fatalf("incoming sidecar not written: %v", err)
	}
	if !strings.Contains(string(incoming), "# My Guide") {
		t.Errorf("sidecar missing user content: %q", string(incoming))
	}
	if !strings.Contains(string(incoming), "# Savepoint Instructions") {
		t.Errorf("sidecar missing managed content: %q", string(incoming))
	}
	if strings.Count(string(incoming), managedBegin) != 1 {
		t.Errorf("sidecar should carry exactly one managed block: %q", string(incoming))
	}
}

func TestUpgradeProjectAssets_halfMarkedAgentsMdConflicts(t *testing.T) {
	// One marker of the pair has no block to replace, so treating it as marked
	// would splice the file at a single offset and corrupt it.
	cases := map[string]string{
		"begin only": "# My Guide\n\n" + managedBegin + "\n# Half a block\n",
		"end only":   "# My Guide\n\n# Half a block\n" + managedEnd + "\n",
		"reversed":   "# My Guide\n\n" + managedEnd + "\nInverted\n" + managedBegin + "\n",
	}

	for name, existingGuide := range cases {
		t.Run(name, func(t *testing.T) {
			target := t.TempDir()
			savepointDir := filepath.Join(target, ".savepoint")
			if err := os.MkdirAll(savepointDir, 0755); err != nil {
				t.Fatal(err)
			}

			guidePath := filepath.Join(target, "AGENTS.md")
			testutil.WriteFile(t, guidePath, existingGuide)

			templates := fstest.MapFS{
				"AGENTS.md": &fstest.MapFile{Data: []byte("# Savepoint Instructions")},
			}

			report, err := UpgradeProjectAssets(templates, target, false, false)
			if err != nil {
				t.Fatalf("UpgradeProjectAssets() error = %v", err)
			}

			action, found := actionFor(report, "AGENTS.md")
			if !found {
				t.Fatal("AGENTS.md not in report")
			}
			if action != ActionConflict {
				t.Errorf("action = %v, want conflict", action)
			}

			data, err := os.ReadFile(guidePath)
			if err != nil {
				t.Fatal(err)
			}
			if string(data) != existingGuide {
				t.Errorf("guide changed: got %q, want %q", string(data), existingGuide)
			}
			if _, err := os.Stat(guidePath + incomingSuffix); err != nil {
				t.Errorf("incoming sidecar not written: %v", err)
			}
		})
	}
}

func TestUpgradeProjectAssets_forceAdoptsUnmarkedAgentsMd(t *testing.T) {
	target := t.TempDir()
	savepointDir := filepath.Join(target, ".savepoint")
	if err := os.MkdirAll(savepointDir, 0755); err != nil {
		t.Fatal(err)
	}

	guidePath := filepath.Join(target, "AGENTS.md")
	existingGuide := "# My Guide\n\nExisting user content.\n"
	testutil.WriteFile(t, guidePath, existingGuide)

	templates := fstest.MapFS{
		"AGENTS.md": &fstest.MapFile{Data: []byte("# Savepoint Instructions")},
	}

	report, err := UpgradeProjectAssets(templates, target, false, true)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	entry, found := entryFor(report, "AGENTS.md")
	if !found {
		t.Fatal("AGENTS.md not in report")
	}
	if entry.Action != ActionMerged {
		t.Errorf("action = %v, want merged", entry.Action)
	}
	if entry.Note != noteBackup {
		t.Errorf("note = %q, want %q", entry.Note, noteBackup)
	}

	data, err := os.ReadFile(guidePath)
	if err != nil {
		t.Fatal(err)
	}
	got := string(data)
	if !strings.Contains(got, "# My Guide") {
		t.Errorf("force dropped user content: %q", got)
	}
	if !strings.Contains(got, "# Savepoint Instructions") {
		t.Errorf("force missing managed content: %q", got)
	}

	backup, err := os.ReadFile(guidePath + backupSuffix)
	if err != nil {
		t.Fatalf("backup not written: %v", err)
	}
	if string(backup) != existingGuide {
		t.Errorf("backup = %q, want %q", string(backup), existingGuide)
	}
	if _, err := os.Stat(guidePath + incomingSuffix); !os.IsNotExist(err) {
		t.Errorf("force should not write an incoming sidecar")
	}
}

func TestUpgradeProjectAssets_dryRunUnmarkedAgentsMdConflicts(t *testing.T) {
	for _, tc := range []struct {
		name  string
		force bool
		want  UpgradeAction
	}{
		{name: "conflict", force: false, want: ActionConflict},
		{name: "force", force: true, want: ActionMerged},
	} {
		t.Run(tc.name, func(t *testing.T) {
			target := t.TempDir()
			savepointDir := filepath.Join(target, ".savepoint")
			if err := os.MkdirAll(savepointDir, 0755); err != nil {
				t.Fatal(err)
			}

			guidePath := filepath.Join(target, "AGENTS.md")
			existingGuide := "# My Guide\n\nExisting user content.\n"
			testutil.WriteFile(t, guidePath, existingGuide)

			templates := fstest.MapFS{
				"AGENTS.md": &fstest.MapFile{Data: []byte("# Savepoint Instructions")},
			}

			report, err := UpgradeProjectAssets(templates, target, true, tc.force)
			if err != nil {
				t.Fatalf("UpgradeProjectAssets() dry-run error = %v", err)
			}

			action, found := actionFor(report, "AGENTS.md")
			if !found {
				t.Fatal("AGENTS.md not in dry-run report")
			}
			if action != tc.want {
				t.Errorf("dry-run action = %v, want %v", action, tc.want)
			}

			data, err := os.ReadFile(guidePath)
			if err != nil {
				t.Fatal(err)
			}
			if string(data) != existingGuide {
				t.Errorf("dry run changed the guide: got %q", string(data))
			}
			for _, suffix := range []string{incomingSuffix, backupSuffix} {
				if _, err := os.Stat(guidePath + suffix); !os.IsNotExist(err) {
					t.Errorf("dry run wrote %s sidecar", suffix)
				}
			}
		})
	}
}

func TestUpgradeProjectAssets_conflictSidecarKeepsGuideCasing(t *testing.T) {
	target := t.TempDir()
	savepointDir := filepath.Join(target, ".savepoint")
	if err := os.MkdirAll(savepointDir, 0755); err != nil {
		t.Fatal(err)
	}

	variantPath := filepath.Join(target, "Agents.MD")
	testutil.WriteFile(t, variantPath, "# My Guide\n")

	templates := fstest.MapFS{
		"AGENTS.md": &fstest.MapFile{Data: []byte("# Savepoint Instructions")},
	}

	report, err := UpgradeProjectAssets(templates, target, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	action, found := actionFor(report, "AGENTS.md")
	if !found {
		t.Fatal("AGENTS.md not in report")
	}
	if action != ActionConflict {
		t.Errorf("action = %v, want conflict", action)
	}

	if _, err := os.Stat(variantPath + incomingSuffix); err != nil {
		t.Errorf("sidecar not written beside the on-disk casing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(target, "AGENTS.md")); !os.IsNotExist(err) {
		t.Errorf("upgrade created a second agent guide under canonical casing")
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
	testutil.WriteFile(t, variantPath, "# My Guide\n\n"+managedBegin+"\n# Old\n"+managedEnd+"\n")

	templates := fstest.MapFS{
		"AGENTS.md": &fstest.MapFile{Data: []byte("# Savepoint Instructions")},
	}

	report, err := UpgradeProjectAssets(templates, target, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	action, found := actionFor(report, "AGENTS.md")
	if !found {
		t.Fatal("AGENTS.md not in report")
	}
	if action != ActionMerged {
		t.Errorf("action = %v, want merged", action)
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
	existingGuide := "# My Guide\n\n" + managedBegin + "\n# Old\n" + managedEnd + "\n"
	testutil.WriteFile(t, variantPath, existingGuide)

	templates := fstest.MapFS{
		"AGENTS.md": &fstest.MapFile{Data: []byte("# Savepoint Instructions")},
	}

	report, err := UpgradeProjectAssets(templates, target, true, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() dry-run error = %v", err)
	}

	action, found := actionFor(report, "AGENTS.md")
	if !found {
		t.Fatal("AGENTS.md not in dry-run report")
	}
	if action != ActionMerged {
		t.Errorf("dry-run action = %v, want merged", action)
	}

	data, err := os.ReadFile(variantPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != existingGuide {
		t.Errorf("dry run changed the guide: got %q", string(data))
	}
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

func entryFor(report *UpgradeReport, path string) (UpgradeEntry, bool) {
	for _, e := range report.Actions {
		if e.Path == path {
			return e, true
		}
	}
	return UpgradeEntry{}, false
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

// --- Skill conflict policy -------------------------------------------------

const policySkillPath = "agent-skills/savepoint-build-task/SKILL.md"

func skillTemplates(content string) fstest.MapFS {
	return fstest.MapFS{policySkillPath: &fstest.MapFile{Data: []byte(content)}}
}

// skillProject creates a Savepoint project whose skill file holds onDisk, with
// the manifest recording the hash of recorded. An empty string means the file
// or the manifest entry is absent.
func skillProject(t *testing.T, onDisk, recorded string) string {
	t.Helper()
	dir := t.TempDir()
	testutil.MkdirAll(t, filepath.Join(dir, ".savepoint"))
	if onDisk != "" {
		testutil.WriteFile(t, filepath.Join(dir, filepath.FromSlash(policySkillPath)), onDisk)
	}
	if recorded != "" {
		manifest := NewManifest()
		manifest.Record(policySkillPath, []byte(recorded))
		if err := manifest.Save(dir); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func skillFile(t *testing.T, dir, suffix string) (string, bool) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(policySkillPath)) + suffix)
	if os.IsNotExist(err) {
		return "", false
	}
	if err != nil {
		t.Fatal(err)
	}
	return string(data), true
}

func assertSkill(t *testing.T, dir, suffix, want string) {
	t.Helper()
	got, ok := skillFile(t, dir, suffix)
	if !ok {
		t.Fatalf("%s%s missing, want %q", policySkillPath, suffix, want)
	}
	if got != want {
		t.Errorf("%s%s = %q, want %q", policySkillPath, suffix, got, want)
	}
}

func assertNoSidecars(t *testing.T, dir string) {
	t.Helper()
	for _, suffix := range []string{incomingSuffix, backupSuffix} {
		if got, ok := skillFile(t, dir, suffix); ok {
			t.Errorf("%s%s exists (%q), want none", policySkillPath, suffix, got)
		}
	}
}

func assertSkillHash(t *testing.T, dir, want string) {
	t.Helper()
	manifest, err := LoadManifest(dir)
	if err != nil {
		t.Fatal(err)
	}
	hash, ok := manifest.Hash(policySkillPath)
	if !ok {
		t.Fatalf("no manifest entry for %s", policySkillPath)
	}
	if hash != hashContent([]byte(want)) {
		t.Errorf("manifest hash does not match %q", want)
	}
}

// upgradeSkill runs an upgrade of a single skill and returns its report entry.
func upgradeSkill(t *testing.T, dir, template string, dryRun, force bool) UpgradeEntry {
	t.Helper()
	report, err := UpgradeProjectAssets(skillTemplates(template), dir, dryRun, force)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}
	for _, e := range report.Actions {
		if e.Path == policySkillPath {
			return e
		}
	}
	t.Fatalf("%s not in report: %+v", policySkillPath, report.Actions)
	return UpgradeEntry{}
}

func TestUpgradeSkill_missingFileIsWritten(t *testing.T) {
	dir := skillProject(t, "", "")

	if entry := upgradeSkill(t, dir, "# New", false, false); entry.Action != ActionUpdated {
		t.Errorf("action = %v, want updated", entry.Action)
	}

	assertSkill(t, dir, "", "# New")
	assertSkillHash(t, dir, "# New")
	assertNoSidecars(t, dir)
}

func TestUpgradeSkill_identicalFileIsUnchanged(t *testing.T) {
	dir := skillProject(t, "# Same", "")

	if entry := upgradeSkill(t, dir, "# Same", false, false); entry.Action != ActionUnchanged {
		t.Errorf("action = %v, want unchanged", entry.Action)
	}

	assertSkill(t, dir, "", "# Same")
	assertSkillHash(t, dir, "# Same")
	assertNoSidecars(t, dir)
}

func TestUpgradeSkill_trackedAndOutdatedIsReplaced(t *testing.T) {
	dir := skillProject(t, "# Ours v1", "# Ours v1")

	if entry := upgradeSkill(t, dir, "# Ours v2", false, false); entry.Action != ActionUpdated {
		t.Errorf("action = %v, want updated", entry.Action)
	}

	assertSkill(t, dir, "", "# Ours v2")
	assertSkillHash(t, dir, "# Ours v2")
	assertNoSidecars(t, dir)
}

func TestUpgradeSkill_customizedFileConflicts(t *testing.T) {
	dir := skillProject(t, "# Ours v1 plus my edits", "# Ours v1")

	entry := upgradeSkill(t, dir, "# Ours v2", false, false)
	if entry.Action != ActionConflict {
		t.Errorf("action = %v, want conflict", entry.Action)
	}
	if entry.Note == "" {
		t.Error("conflict entry has no note naming the incoming file")
	}

	assertSkill(t, dir, "", "# Ours v1 plus my edits")
	assertSkill(t, dir, incomingSuffix, "# Ours v2")
	if _, ok := skillFile(t, dir, backupSuffix); ok {
		t.Error("conflict wrote a backup, want none")
	}
	assertSkillHash(t, dir, "# Ours v1")
}

func TestUpgradeSkill_untrackedFileIsBackedUpAndReplaced(t *testing.T) {
	dir := skillProject(t, "# Pre-manifest content", "")

	entry := upgradeSkill(t, dir, "# Ours v2", false, false)
	if entry.Action != ActionUpdated {
		t.Errorf("action = %v, want updated", entry.Action)
	}
	if entry.Note == "" {
		t.Error("migration entry has no note naming the backup")
	}

	assertSkill(t, dir, "", "# Ours v2")
	assertSkill(t, dir, backupSuffix, "# Pre-manifest content")
	assertSkillHash(t, dir, "# Ours v2")
}

func TestUpgradeSkill_forceOverridesConflict(t *testing.T) {
	dir := skillProject(t, "# Ours v1 plus my edits", "# Ours v1")

	entry := upgradeSkill(t, dir, "# Ours v2", false, true)
	if entry.Action != ActionUpdated {
		t.Errorf("action = %v, want updated", entry.Action)
	}

	assertSkill(t, dir, "", "# Ours v2")
	assertSkill(t, dir, backupSuffix, "# Ours v1 plus my edits")
	assertSkillHash(t, dir, "# Ours v2")
	if _, ok := skillFile(t, dir, incomingSuffix); ok {
		t.Error("force wrote a .new file, want none")
	}
}

func TestUpgradeSkill_forceLeavesUnrelatedCasesAlone(t *testing.T) {
	dir := skillProject(t, "# Same", "# Same")

	if entry := upgradeSkill(t, dir, "# Same", false, true); entry.Action != ActionUnchanged {
		t.Errorf("action = %v, want unchanged", entry.Action)
	}
	assertNoSidecars(t, dir)
}

func TestUpgradeSkill_dryRunMatchesRealRunAndWritesNothing(t *testing.T) {
	cases := []struct {
		name     string
		onDisk   string
		recorded string
		template string
		force    bool
		want     UpgradeAction
	}{
		{"missing", "", "", "# New", false, ActionUpdated},
		{"identical", "# Same", "# Same", "# Same", false, ActionUnchanged},
		{"tracked outdated", "# Ours v1", "# Ours v1", "# Ours v2", false, ActionUpdated},
		{"customized", "# My edits", "# Ours v1", "# Ours v2", false, ActionConflict},
		{"untracked", "# Pre-manifest", "", "# Ours v2", false, ActionUpdated},
		{"forced", "# My edits", "# Ours v1", "# Ours v2", true, ActionUpdated},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			live := skillProject(t, tc.onDisk, tc.recorded)
			liveEntry := upgradeSkill(t, live, tc.template, false, tc.force)
			if liveEntry.Action != tc.want {
				t.Fatalf("real-run action = %v, want %v", liveEntry.Action, tc.want)
			}

			dir := skillProject(t, tc.onDisk, tc.recorded)
			manifestBefore, _ := os.ReadFile(manifestPath(dir))

			dryEntry := upgradeSkill(t, dir, tc.template, true, tc.force)
			if dryEntry.Action != liveEntry.Action {
				t.Errorf("dry-run action = %v, want %v", dryEntry.Action, liveEntry.Action)
			}

			if got, ok := skillFile(t, dir, ""); ok != (tc.onDisk != "") || got != tc.onDisk {
				t.Errorf("dry run changed the skill file: got %q (exists %v), want %q", got, ok, tc.onDisk)
			}
			assertNoSidecars(t, dir)
			manifestAfter, _ := os.ReadFile(manifestPath(dir))
			if string(manifestAfter) != string(manifestBefore) {
				t.Errorf("dry run changed the manifest: %q -> %q", manifestBefore, manifestAfter)
			}
		})
	}
}

func TestUpgradeSkill_secondRunIsIdempotent(t *testing.T) {
	dir := skillProject(t, "# Ours v1", "# Ours v1")
	upgradeSkill(t, dir, "# Ours v2", false, false)

	// An upgrade that changes nothing must change nothing on disk, manifest
	// included: a new file identity is still an observable write.
	manifestBefore, err := os.Stat(manifestPath(dir))
	if err != nil {
		t.Fatal(err)
	}

	report, err := UpgradeProjectAssets(skillTemplates("# Ours v2"), dir, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}
	for _, e := range report.Actions {
		if e.Action != ActionUnchanged {
			t.Errorf("second run %s = %v, want unchanged", e.Path, e.Action)
		}
	}
	assertNoSidecars(t, dir)

	manifestAfter, err := os.Stat(manifestPath(dir))
	if err != nil {
		t.Fatal(err)
	}
	if !os.SameFile(manifestBefore, manifestAfter) {
		t.Error("second run replaced the unchanged manifest")
	}
	if !manifestBefore.ModTime().Equal(manifestAfter.ModTime()) {
		t.Error("second run rewrote the unchanged manifest in place")
	}
}

func TestUpgradeSkill_unresolvedConflictRefreshesTheSameSidecar(t *testing.T) {
	dir := skillProject(t, "# My edits", "# Ours v1")
	upgradeSkill(t, dir, "# Ours v2", false, false)

	if entry := upgradeSkill(t, dir, "# Ours v3", false, false); entry.Action != ActionConflict {
		t.Errorf("second action = %v, want conflict", entry.Action)
	}

	assertSkill(t, dir, "", "# My edits")
	assertSkill(t, dir, incomingSuffix, "# Ours v3")
	if _, ok := skillFile(t, dir, incomingSuffix+incomingSuffix); ok {
		t.Error("upgrade accumulated a second .new variant")
	}
}

// --- Legacy projects -------------------------------------------------------

// legacyFixture loads a frozen fixture from testdata/legacy. The fixtures are
// files already in the wild; see that directory's README for the freeze rule.
func legacyFixture(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "legacy", name))
	if err != nil {
		t.Fatalf("read legacy fixture %s: %v", name, err)
	}
	return string(data)
}

// legacyProject assembles a pre-manifest project: a customized skill, the given
// agent guide fixture, and no recorded provenance for either.
func legacyProject(t *testing.T, guideFixture string) (dir, guide, skill string) {
	t.Helper()
	dir = t.TempDir()
	testutil.MkdirAll(t, filepath.Join(dir, ".savepoint"))

	guide = legacyFixture(t, guideFixture)
	testutil.WriteFile(t, filepath.Join(dir, "AGENTS.md"), guide)

	skill = legacyFixture(t, "SKILL.customized.md")
	testutil.WriteFile(t, filepath.Join(dir, filepath.FromSlash(policySkillPath)), skill)

	return dir, guide, skill
}

func legacyTemplates() fstest.MapFS {
	return fstest.MapFS{
		"AGENTS.md":     &fstest.MapFile{Data: []byte("# Agents Guide\n\nCurrent managed guidance.")},
		policySkillPath: &fstest.MapFile{Data: []byte("# Savepoint Skill: Build Task\n\nCurrent skill body.")},
	}
}

func TestUpgradeProjectAssets_legacyProjectWithUnmarkedGuide(t *testing.T) {
	dir, guide, skill := legacyProject(t, "AGENTS.unmarked.md")
	templates := legacyTemplates()

	report, err := UpgradeProjectAssets(templates, dir, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	guideEntry, found := entryFor(report, "AGENTS.md")
	if !found {
		t.Fatal("AGENTS.md not in report")
	}
	if guideEntry.Action != ActionConflict {
		t.Errorf("unmarked guide action = %v, want conflict", guideEntry.Action)
	}

	guidePath := filepath.Join(dir, "AGENTS.md")
	data, err := os.ReadFile(guidePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != guide {
		t.Errorf("unmarked guide changed:\n got %q\nwant %q", string(data), guide)
	}
	if _, err := os.Stat(guidePath + incomingSuffix); err != nil {
		t.Errorf("incoming guide sidecar not written: %v", err)
	}

	// A pre-manifest skill has no provenance, so the migration rule applies:
	// keep a recoverable copy, replace, and record the hash from here on.
	skillEntry, found := entryFor(report, policySkillPath)
	if !found {
		t.Fatal("skill not in report")
	}
	if skillEntry.Action != ActionUpdated || skillEntry.Note != noteBackup {
		t.Errorf("skill entry = %+v, want updated with a backup note", skillEntry)
	}
	assertSkill(t, dir, "", string(templates[policySkillPath].Data))
	assertSkill(t, dir, backupSuffix, skill)
	assertSkillHash(t, dir, string(templates[policySkillPath].Data))
}

func TestUpgradeProjectAssets_legacyProjectWithMarkedGuide(t *testing.T) {
	dir, guide, _ := legacyProject(t, "AGENTS.marked.md")
	templates := legacyTemplates()

	report, err := UpgradeProjectAssets(templates, dir, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	action, found := actionFor(report, "AGENTS.md")
	if !found {
		t.Fatal("AGENTS.md not in report")
	}
	if action != ActionMerged {
		t.Errorf("marked guide action = %v, want merged", action)
	}

	data, err := os.ReadFile(filepath.Join(dir, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	merged := string(data)

	// Prose on both sides of the markers survives; only the block is rewritten.
	for _, want := range []string{
		"Prose the team wrote, above the managed block.",
		"Prose the team wrote, below the managed block.",
		"Current managed guidance.",
	} {
		if !strings.Contains(merged, want) {
			t.Errorf("merged guide missing %q:\n%s", want, merged)
		}
	}
	if strings.Contains(merged, "Follow the phase prompt") {
		t.Errorf("merged guide kept the old managed block:\n%s", merged)
	}
	if got := strings.Count(merged, managedBegin); got != 1 {
		t.Errorf("merged guide has %d managed blocks, want 1", got)
	}
	if merged == guide {
		t.Error("merged guide unchanged, want the managed block refreshed")
	}
	if _, err := os.Stat(filepath.Join(dir, "AGENTS.md"+incomingSuffix)); !os.IsNotExist(err) {
		t.Error("marked guide wrote a conflict sidecar")
	}
}

func TestUpgradeProjectAssets_legacyCustomizedSkillConflictsOnceTracked(t *testing.T) {
	// After the one-time migration the project has exact provenance, so the
	// same local tailoring must conflict rather than be overwritten again.
	dir, _, skill := legacyProject(t, "AGENTS.marked.md")
	templates := legacyTemplates()

	if _, err := UpgradeProjectAssets(templates, dir, false, false); err != nil {
		t.Fatalf("first UpgradeProjectAssets() error = %v", err)
	}
	testutil.WriteFile(t, filepath.Join(dir, filepath.FromSlash(policySkillPath)), skill)

	entry := upgradeSkill(t, dir, "# Savepoint Skill: Build Task\n\nNewer skill body.", false, false)
	if entry.Action != ActionConflict {
		t.Errorf("action = %v, want conflict", entry.Action)
	}
	assertSkill(t, dir, "", skill)
	assertSkill(t, dir, incomingSuffix, "# Savepoint Skill: Build Task\n\nNewer skill body.")
}

func TestUpgradeReport_formatConflict(t *testing.T) {
	r := &UpgradeReport{
		Actions: []UpgradeEntry{
			{Path: policySkillPath, Action: ActionConflict, Note: noteConflict},
			{Path: "agent-skills/other/SKILL.md", Action: ActionUpdated, Note: noteBackup},
		},
	}

	output := r.Format()
	if !strings.Contains(output, "Conflicts: 1") {
		t.Errorf("missing conflict count: %q", output)
	}
	if !strings.Contains(output, "conflict  "+policySkillPath) {
		t.Errorf("missing conflict path: %q", output)
	}
	if !strings.Contains(output, noteConflict) || !strings.Contains(output, noteBackup) {
		t.Errorf("missing sidecar notes: %q", output)
	}
}
