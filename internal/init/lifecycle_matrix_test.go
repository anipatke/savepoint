package init

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/testutil"
)

// T007's cross-tree lifecycle matrix. Every mechanism exercised here —
// manifest provenance and the AGENTS.md marker pair — is already proven in
// isolation elsewhere, against synthetic template trees built to exercise one
// rule at a time. What none of those prove is the property the epic promises
// as a whole: that the real shipped V2 tree, run through the real loop a
// project actually takes, never loses a byte a user wrote. So these tests use
// realV2Tree, not fstest.MapFS, and run UpgradeProjectAssets, not the
// single-tree upgradeAssetsFromTree. V1 retirement against a real shipped V1
// tree no longer applies — that tree is gone (O-021) — and stays covered
// against synthetic trees in retire_v1_skills_test.go.

// realV2Tree returns the actual shipped V2 template tree, the same one
// main.go embeds.
func realV2Tree(t *testing.T) fs.FS {
	t.Helper()
	return os.DirFS(filepath.Join("..", "..", "templates", "project-v2"))
}

// skillEntrypoints lists every agent-skills/*/SKILL.md path a template tree
// ships, excluding the shared, non-triggerable references/ directory.
func skillEntrypoints(t *testing.T, templates fs.FS) []string {
	t.Helper()
	var paths []string
	err := fs.WalkDir(templates, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, "/SKILL.md") {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk templates: %v", err)
	}
	if len(paths) == 0 {
		t.Fatal("no skill entrypoints found in template tree")
	}
	return paths
}

// v2FreshProject scaffolds a project directly on the real V2 tree: `savepoint
// init`'s own default, carrying every V2 skill and no V1 skill ever.
func v2FreshProject(t *testing.T, v2Templates fs.FS) string {
	t.Helper()
	dir := t.TempDir()
	if err := Scaffold(v2Templates, dir, ProjectNameFromDir(dir), false); err != nil {
		t.Fatalf("Scaffold(v2) error = %v", err)
	}
	return dir
}

func TestLifecycleMatrix_fullLoopAcrossProjectKinds(t *testing.T) {
	v2Templates := realV2Tree(t)
	v2Skills := skillEntrypoints(t, v2Templates)
	dir := v2FreshProject(t, v2Templates)

	for _, skill := range retiredV1SkillDirs {
		if _, err := os.Stat(filepath.Join(dir, "agent-skills", skill, "SKILL.md")); err == nil {
			t.Errorf("before upgrade: V1 skill %s unexpectedly present on a fresh V2 project", skill)
		}
	}

	report, err := UpgradeProjectAssets(v2Templates, dir, false, false)
	if err != nil {
		t.Fatalf("first UpgradeProjectAssets() error = %v", err)
	}
	for _, e := range report.Actions {
		if e.Action == ActionRetired {
			t.Errorf("unexpected retirement action on a fresh V2 project: %+v", e)
		}
	}

	for _, path := range v2Skills {
		if _, err := os.Stat(filepath.Join(dir, filepath.FromSlash(path))); err != nil {
			t.Errorf("V2 skill %s not installed: %v", path, err)
		}
	}

	// The second run, against the now-settled project, must be a true
	// no-op: no content change, no mtime change.
	before := dirSnapshot(t, dir)
	beforeTimes := mtimeSnapshot(t, dir)

	second, err := UpgradeProjectAssets(v2Templates, dir, false, false)
	if err != nil {
		t.Fatalf("second UpgradeProjectAssets() error = %v", err)
	}
	for _, e := range second.Actions {
		if e.Action != ActionUnchanged && e.Action != ActionSkipped && e.Action != ActionInfo {
			t.Errorf("second run action %v for %s, want unchanged/skipped/info", e.Action, e.Path)
		}
	}
	assertNoChange(t, dir, before)
	after := mtimeSnapshot(t, dir)
	for path, mt := range beforeTimes {
		if after[path] != mt {
			t.Errorf("%s mtime changed on an idempotent rerun", path)
		}
	}
}

func TestLifecycleMatrix_editedV2SkillConflicts(t *testing.T) {
	v2Templates := realV2Tree(t)
	dir := v2FreshProject(t, v2Templates)

	target := skillEntrypoints(t, v2Templates)[0]
	targetPath := filepath.Join(dir, filepath.FromSlash(target))
	edited := readProjectFile(t, dir, target) + "\n<!-- local edit -->\n"
	if err := AtomicWrite(targetPath, []byte(edited)); err != nil {
		t.Fatal(err)
	}

	report, err := UpgradeProjectAssets(v2Templates, dir, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	entry, found := entryFor(report, target)
	if !found || entry.Action != ActionConflict || entry.Note != noteConflict {
		t.Fatalf("%s = %+v (found %v), want conflict with the standard note", target, entry, found)
	}
	if got := readProjectFile(t, dir, target); got != edited {
		t.Error("conflicted V2 skill's bytes changed, want unchanged")
	}

	incoming, err := os.ReadFile(targetPath + incomingSuffix)
	if err != nil {
		t.Fatalf("incoming sidecar not written: %v", err)
	}
	want, err := fs.ReadFile(v2Templates, target)
	if err != nil {
		t.Fatal(err)
	}
	if string(incoming) != string(want) {
		t.Errorf("incoming sidecar = %q, want the shipped V2 content", incoming)
	}
}

// TestLifecycleMatrix_untrackedV2SkillIsBackedUpAndReplaced covers the state
// T007's problem statement names: a V2 skill file on disk with no manifest
// entry for it. `--force` is not involved; this is the plain untracked path.
func TestLifecycleMatrix_untrackedV2SkillIsBackedUpAndReplaced(t *testing.T) {
	v2Templates := realV2Tree(t)
	dir := t.TempDir()
	testutil.MkdirAll(t, filepath.Join(dir, ".savepoint"))
	testutil.WriteFile(t, filepath.Join(dir, ".savepoint", "config.yml"), "schema_version: 2\n")

	target := skillEntrypoints(t, v2Templates)[0]
	targetPath := filepath.Join(dir, filepath.FromSlash(target))
	priorContent := "# Pre-manifest content the user already had\n"
	testutil.WriteFile(t, targetPath, priorContent)

	report, err := UpgradeProjectAssets(v2Templates, dir, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	entry, found := entryFor(report, target)
	if !found || entry.Action != ActionUpdated || entry.Note != noteBackup {
		t.Fatalf("%s = %+v (found %v), want updated with a backup note", target, entry, found)
	}

	backup, err := os.ReadFile(targetPath + backupSuffix)
	if err != nil {
		t.Fatalf("backup not written: %v", err)
	}
	if string(backup) != priorContent {
		t.Errorf("backup = %q, want byte-identical %q", backup, priorContent)
	}

	want, err := fs.ReadFile(v2Templates, target)
	if err != nil {
		t.Fatal(err)
	}
	if got := readProjectFile(t, dir, target); got != string(want) {
		t.Error("skill was not refreshed to the shipped V2 content")
	}
}

func TestLifecycleMatrix_forceAdoptsConflictedSkillAndUnmarkedGuideTogether(t *testing.T) {
	v2Templates := realV2Tree(t)
	dir := v2FreshProject(t, v2Templates)

	target := skillEntrypoints(t, v2Templates)[0]
	targetPath := filepath.Join(dir, filepath.FromSlash(target))
	editedSkill := readProjectFile(t, dir, target) + "\n<!-- local edit -->\n"
	if err := AtomicWrite(targetPath, []byte(editedSkill)); err != nil {
		t.Fatal(err)
	}

	guidePath := filepath.Join(dir, "AGENTS.md")
	unmarkedGuide := "# Team Guide\n\nOur own onboarding notes, no Savepoint markers.\n"
	if err := AtomicWrite(guidePath, []byte(unmarkedGuide)); err != nil {
		t.Fatal(err)
	}

	report, err := UpgradeProjectAssets(v2Templates, dir, false, true)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	skillEntry, found := entryFor(report, target)
	if !found || skillEntry.Action != ActionUpdated || skillEntry.Note != noteBackup {
		t.Fatalf("skill entry = %+v (found %v), want updated with a backup note", skillEntry, found)
	}
	skillBackup, err := os.ReadFile(targetPath + backupSuffix)
	if err != nil {
		t.Fatalf("skill backup not written: %v", err)
	}
	if string(skillBackup) != editedSkill {
		t.Errorf("skill backup = %q, want the user's edited content", skillBackup)
	}

	guideEntry, found := entryFor(report, "AGENTS.md")
	if !found || guideEntry.Action != ActionMerged || guideEntry.Note != noteBackup {
		t.Fatalf("guide entry = %+v (found %v), want merged with a backup note", guideEntry, found)
	}
	guideBackup, err := os.ReadFile(guidePath + backupSuffix)
	if err != nil {
		t.Fatalf("guide backup not written: %v", err)
	}
	if string(guideBackup) != unmarkedGuide {
		t.Errorf("guide backup = %q, want the user's unmarked guide", guideBackup)
	}
	if got := readProjectFile(t, dir, "AGENTS.md"); !strings.Contains(got, managedBegin) {
		t.Errorf("forced guide adoption missing managed block: %q", got)
	}
}

func TestLifecycleMatrix_markedGuideRefreshesOnlyBetweenMarkers(t *testing.T) {
	v2Templates := realV2Tree(t)
	dir := v2FreshProject(t, v2Templates)

	guidePath := filepath.Join(dir, "AGENTS.md")
	before := "# Our Team\n\nAbove-the-fold notes.\n\n"
	after := "\n\nBelow the marker: deploy on Fridays only with sign-off.\n"
	existing := before + managedBegin + "\n# Stale managed block\n" + managedEnd + after
	if err := AtomicWrite(guidePath, []byte(existing)); err != nil {
		t.Fatal(err)
	}

	report, err := UpgradeProjectAssets(v2Templates, dir, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	if action, found := actionFor(report, "AGENTS.md"); !found || action != ActionMerged {
		t.Fatalf("action = %v (found %v), want merged", action, found)
	}

	got := readProjectFile(t, dir, "AGENTS.md")
	if !strings.HasPrefix(got, before) {
		t.Errorf("content above the opening marker changed: %q", got)
	}
	if !strings.HasSuffix(got, after) {
		t.Errorf("content below the closing marker changed: %q", got)
	}
	if strings.Contains(got, "Stale managed block") {
		t.Error("stale managed block survived the refresh")
	}
	wantBlock, err := fs.ReadFile(v2Templates, "AGENTS.md")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, strings.TrimSpace(string(wantBlock))) {
		t.Error("managed block was not refreshed to the shipped V2 guide")
	}
}

func TestLifecycleMatrix_v2TreeGuideConflictCases(t *testing.T) {
	v2Templates := realV2Tree(t)
	wantBlock, err := fs.ReadFile(v2Templates, "AGENTS.md")
	if err != nil {
		t.Fatal(err)
	}

	cases := map[string]string{
		"unmarked":   "# My Guide\n\nExisting user content.\n",
		"begin only": "# My Guide\n\n" + managedBegin + "\n# Half a block\n",
		"end only":   "# My Guide\n\n# Half a block\n" + managedEnd + "\n",
		"reversed":   "# My Guide\n\n" + managedEnd + "\nInverted\n" + managedBegin + "\n",
	}

	for name, existingGuide := range cases {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			testutil.MkdirAll(t, filepath.Join(dir, ".savepoint"))
			testutil.WriteFile(t, filepath.Join(dir, ".savepoint", "config.yml"), "schema_version: 2\n")
			guidePath := filepath.Join(dir, "AGENTS.md")
			testutil.WriteFile(t, guidePath, existingGuide)

			report, err := UpgradeProjectAssets(v2Templates, dir, false, false)
			if err != nil {
				t.Fatalf("UpgradeProjectAssets() error = %v", err)
			}

			entry, found := entryFor(report, "AGENTS.md")
			if !found || entry.Action != ActionConflict {
				t.Fatalf("action = %+v (found %v), want conflict", entry, found)
			}
			if got := readProjectFile(t, dir, "AGENTS.md"); got != existingGuide {
				t.Error("guide changed, want byte-identical")
			}
			incoming, err := os.ReadFile(guidePath + incomingSuffix)
			if err != nil {
				t.Fatalf("incoming sidecar not written: %v", err)
			}
			if !strings.Contains(string(incoming), strings.TrimSpace(string(wantBlock))) {
				t.Error("sidecar missing the shipped V2 managed content")
			}
		})
	}
}

func TestLifecycleMatrix_guideCasingVariantPreservedAcrossSidecars(t *testing.T) {
	v2Templates := realV2Tree(t)

	t.Run("conflict sidecar", func(t *testing.T) {
		dir := t.TempDir()
		testutil.MkdirAll(t, filepath.Join(dir, ".savepoint"))
		testutil.WriteFile(t, filepath.Join(dir, ".savepoint", "config.yml"), "schema_version: 2\n")
		variantPath := filepath.Join(dir, "Agents.MD")
		testutil.WriteFile(t, variantPath, "# My Guide\n")

		report, err := UpgradeProjectAssets(v2Templates, dir, false, false)
		if err != nil {
			t.Fatalf("UpgradeProjectAssets() error = %v", err)
		}
		if action, found := actionFor(report, "AGENTS.md"); !found || action != ActionConflict {
			t.Fatalf("action = %v (found %v), want conflict", action, found)
		}
		if _, err := os.Stat(variantPath + incomingSuffix); err != nil {
			t.Errorf("incoming sidecar not written beside the on-disk casing: %v", err)
		}
		if _, err := os.Stat(filepath.Join(dir, "AGENTS.md")); !os.IsNotExist(err) {
			t.Error("a second, canonically cased guide was created")
		}
	})

	t.Run("force backup sidecar", func(t *testing.T) {
		dir := t.TempDir()
		testutil.MkdirAll(t, filepath.Join(dir, ".savepoint"))
		testutil.WriteFile(t, filepath.Join(dir, ".savepoint", "config.yml"), "schema_version: 2\n")
		variantPath := filepath.Join(dir, "Agents.MD")
		existing := "# My Guide\n"
		testutil.WriteFile(t, variantPath, existing)

		report, err := UpgradeProjectAssets(v2Templates, dir, false, true)
		if err != nil {
			t.Fatalf("UpgradeProjectAssets() error = %v", err)
		}
		if action, found := actionFor(report, "AGENTS.md"); !found || action != ActionMerged {
			t.Fatalf("action = %v (found %v), want merged", action, found)
		}
		backup, err := os.ReadFile(variantPath + backupSuffix)
		if err != nil {
			t.Fatalf("backup sidecar not written beside the on-disk casing: %v", err)
		}
		if string(backup) != existing {
			t.Errorf("backup = %q, want %q", backup, existing)
		}
		if _, err := os.Stat(filepath.Join(dir, "AGENTS.md")); !os.IsNotExist(err) {
			t.Error("a second, canonically cased guide was created")
		}
	})
}

// TestLifecycleMatrix_absentOptionalFilesProduceNoFinding is TPL-03 against
// the real V2 tree: a fresh project has no Health-Check.md, no Concept.md,
// and no procedures file (none shipped by the V2 tree), and the upgrade must
// complete with no failed entry and no note describing their absence as a
// problem.
func TestLifecycleMatrix_absentOptionalFilesProduceNoFinding(t *testing.T) {
	v2Templates := realV2Tree(t)
	dir := v2FreshProject(t, v2Templates)

	forbidden := []string{"missing", "absent", "warning", "diagnostic"}

	for _, optional := range []string{"Concept.md", "procedures.md"} {
		if _, err := os.Stat(filepath.Join(dir, ".savepoint", optional)); err == nil {
			t.Fatalf("test setup: %s unexpectedly present", optional)
		}
	}

	report, err := UpgradeProjectAssets(v2Templates, dir, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	for _, e := range report.Actions {
		if e.Action == ActionFailed {
			t.Errorf("absence of optional files produced a failure: %+v", e)
		}
		lower := strings.ToLower(e.Note)
		for _, word := range forbidden {
			if strings.Contains(lower, word) {
				t.Errorf("entry %+v treats an optional file's absence as a finding", e)
			}
		}
	}
}

func TestLifecycleMatrix_dryRunMatchesRealRunAcrossProjectKinds(t *testing.T) {
	v2Templates := realV2Tree(t)

	dryDir := v2FreshProject(t, v2Templates)
	before := dirSnapshot(t, dryDir)
	beforeTimes := mtimeSnapshot(t, dryDir)

	dryReport, err := UpgradeProjectAssets(v2Templates, dryDir, true, false)
	if err != nil {
		t.Fatalf("dry run error = %v", err)
	}
	assertNoChange(t, dryDir, before)
	after := mtimeSnapshot(t, dryDir)
	for path, mt := range beforeTimes {
		if after[path] != mt {
			t.Errorf("%s mtime changed by a dry run", path)
		}
	}

	realDir := v2FreshProject(t, v2Templates)
	realReport, err := UpgradeProjectAssets(v2Templates, realDir, false, false)
	if err != nil {
		t.Fatalf("real run error = %v", err)
	}

	if len(dryReport.Actions) != len(realReport.Actions) {
		t.Fatalf("dry-run actions = %d, real actions = %d", len(dryReport.Actions), len(realReport.Actions))
	}
	for i := range dryReport.Actions {
		if dryReport.Actions[i].Path != realReport.Actions[i].Path || dryReport.Actions[i].Action != realReport.Actions[i].Action {
			t.Errorf("action[%d] dry=%+v real=%+v", i, dryReport.Actions[i], realReport.Actions[i])
		}
	}
}

// TestLifecycleMatrix_retirementAndConflictReportedTogether covers the
// combined case: a V2 project that still carries the nine V1 skills (the
// state migrate leaves behind, since it only flips schema_version) whose
// AGENTS.md was rewritten by hand, stripping the marker pair, so the same
// upgrade that retires the nine V1 skills also finds an unmarked guide to
// conflict on. Both recoveries must be named.
func TestLifecycleMatrix_retirementAndConflictReportedTogether(t *testing.T) {
	v2Templates := realV2Tree(t)
	dir := v2ProjectWithLegacySkills(t, nil)

	guidePath := filepath.Join(dir, "AGENTS.md")
	unmarkedGuide := "# Our Rewritten Guide\n\nNo Savepoint markers left in this file.\n"
	if err := AtomicWrite(guidePath, []byte(unmarkedGuide)); err != nil {
		t.Fatal(err)
	}

	report, err := UpgradeProjectAssets(v2Templates, dir, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	formatted := report.Format()
	retired := 0
	for _, e := range report.Actions {
		if e.Action == ActionRetired {
			retired++
			wantNote := "archived copy saved to " + filepath.ToSlash(filepath.Join(migrationsDir, archiveStem(e.Path)+".md"))
			if e.Note != wantNote {
				t.Errorf("retired %s note = %q, want %q", e.Path, e.Note, wantNote)
			}
			if !strings.Contains(formatted, wantNote) {
				t.Errorf("formatted report omits recovery archive for %s: %q", e.Path, formatted)
			}
			archived, err := os.ReadFile(retiredArchivePath(dir, e.Path, 0))
			if err != nil {
				t.Errorf("retired %s has no readable archive: %v", e.Path, err)
			}
			if len(archived) == 0 {
				t.Errorf("retired %s archived empty content", e.Path)
			}
		}
	}
	if retired == 0 {
		t.Fatal("expected retirement actions alongside the guide conflict, found none")
	}

	guideEntry, found := entryFor(report, "AGENTS.md")
	if !found || guideEntry.Action != ActionConflict {
		t.Fatalf("AGENTS.md = %+v (found %v), want conflict", guideEntry, found)
	}
	if !strings.Contains(formatted, guideEntry.Note) || !strings.Contains(formatted, "AGENTS.md") {
		t.Errorf("formatted report omits the guide conflict recovery: %q", formatted)
	}
	if _, err := os.Stat(guidePath + incomingSuffix); err != nil {
		t.Errorf("conflicted guide has no incoming sidecar: %v", err)
	}
	if got := readProjectFile(t, dir, "AGENTS.md"); got != unmarkedGuide {
		t.Error("conflicted guide was not kept byte-identical")
	}
}
