package init

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/opencode/savepoint/internal/testutil"
)

const claudeTemplate = "# Claude Code\n\n@AGENTS.md\n\nSavepoint manages this section."

func claudeTemplates() fstest.MapFS {
	return fstest.MapFS{
		"CLAUDE.md": &fstest.MapFile{Data: []byte(claudeTemplate)},
	}
}

func claudeBlock() string {
	return managedBegin + "\n" + claudeTemplate + "\n" + managedEnd
}

func TestScaffold_writesClaudeGuideWithImport(t *testing.T) {
	dir := t.TempDir()
	if err := Scaffold(claudeTemplates(), dir, "p", false); err != nil {
		t.Fatal(err)
	}
	got := string(mustReadFile(t, filepath.Join(dir, "CLAUDE.md")))
	if got != claudeBlock()+"\n" {
		t.Errorf("CLAUDE.md = %q", got)
	}
	if !strings.Contains(got, "@AGENTS.md") {
		t.Error("import missing")
	}
}

func TestScaffold_keepsUserClaudeText(t *testing.T) {
	dir := t.TempDir()
	user := "# Mine\n\nUse tabs.\n"
	testutil.WriteFile(t, filepath.Join(dir, "CLAUDE.md"), user)
	if err := Scaffold(claudeTemplates(), dir, "p", false); err != nil {
		t.Fatal(err)
	}
	got := string(mustReadFile(t, filepath.Join(dir, "CLAUDE.md")))
	if !strings.HasPrefix(got, user) || !strings.HasSuffix(got, claudeBlock()+"\n") {
		t.Errorf("CLAUDE.md = %q", got)
	}
}

func TestScaffold_noSecondImport(t *testing.T) {
	dir := t.TempDir()
	user := "# Mine\n@AGENTS.md\n"
	testutil.WriteFile(t, filepath.Join(dir, "CLAUDE.md"), user)
	if err := Scaffold(claudeTemplates(), dir, "p", false); err != nil {
		t.Fatal(err)
	}
	if got := string(mustReadFile(t, filepath.Join(dir, "CLAUDE.md"))); got != user {
		t.Errorf("CLAUDE.md = %q, want unchanged %q", got, user)
	}
}

func TestUpgradeProjectAssets_claudeGuideLifecycle(t *testing.T) {
	target := savepointProject(t)
	path := filepath.Join(target, "CLAUDE.md")
	user := "# Mine\n\nUse tabs.\n"
	testutil.WriteFile(t, path, user)

	// Dry run reports the merge and writes nothing.
	report, err := upgradeAssetsFromTree(claudeTemplates(), target, true, false)
	if err != nil {
		t.Fatal(err)
	}
	if a, _ := actionFor(report, "CLAUDE.md"); a != ActionMerged {
		t.Errorf("dry run action = %v, want merged", a)
	}
	if got := string(mustReadFile(t, path)); got != user {
		t.Errorf("dry run changed CLAUDE.md: %q", got)
	}

	report = mustUpgrade(t, claudeTemplates(), target)
	if a, _ := actionFor(report, "CLAUDE.md"); a != ActionMerged {
		t.Errorf("action = %v, want merged", a)
	}
	got := string(mustReadFile(t, path))
	if want := user + "\n" + claudeBlock() + "\n"; got != want {
		t.Errorf("CLAUDE.md = %q, want %q", got, want)
	}

	// Second run on unchanged input changes nothing.
	report = mustUpgrade(t, claudeTemplates(), target)
	if a, _ := actionFor(report, "CLAUDE.md"); a != ActionUnchanged {
		t.Errorf("second action = %v, want unchanged", a)
	}
	if again := string(mustReadFile(t, path)); again != got {
		t.Errorf("second run changed CLAUDE.md: %q", again)
	}
}

func TestUpgradeProjectAssets_claudeGuideRefreshesOnlyBlock(t *testing.T) {
	target := savepointProject(t)
	path := filepath.Join(target, "CLAUDE.md")
	before, after := "top\n", "\nbottom\n"
	testutil.WriteFile(t, path, before+managedBegin+"\nold\n"+managedEnd+after)

	mustUpgrade(t, claudeTemplates(), target)
	if got, want := string(mustReadFile(t, path)), before+claudeBlock()+after; got != want {
		t.Errorf("CLAUDE.md = %q, want %q", got, want)
	}
}

func TestUpgradeProjectAssets_claudeGuideAlreadyImportedUnchanged(t *testing.T) {
	target := savepointProject(t)
	path := filepath.Join(target, "CLAUDE.md")
	user := "@AGENTS.md\n\nmine\n"
	testutil.WriteFile(t, path, user)

	report := mustUpgrade(t, claudeTemplates(), target)
	if a, _ := actionFor(report, "CLAUDE.md"); a != ActionUnchanged {
		t.Errorf("action = %v, want unchanged", a)
	}
	if got := string(mustReadFile(t, path)); got != user {
		t.Errorf("CLAUDE.md = %q", got)
	}
}

func TestUpgradeProjectAssets_claudeGuideMissingIsCreated(t *testing.T) {
	target := savepointProject(t)
	report := mustUpgrade(t, claudeTemplates(), target)
	if a, _ := actionFor(report, "CLAUDE.md"); a != ActionUpdated {
		t.Errorf("action = %v, want updated", a)
	}
	if got := string(mustReadFile(t, filepath.Join(target, "CLAUDE.md"))); got != claudeBlock()+"\n" {
		t.Errorf("CLAUDE.md = %q", got)
	}
}

func TestUpgradeProjectAssets_claudeGuideWriteFailureReported(t *testing.T) {
	target := savepointProject(t)
	// A directory where the file belongs makes the read fail.
	if err := os.Mkdir(filepath.Join(target, "CLAUDE.md"), 0755); err != nil {
		t.Fatal(err)
	}
	report, err := upgradeAssetsFromTree(claudeTemplates(), target, false, false)
	if err == nil {
		t.Fatal("expected an error")
	}
	if a, _ := actionFor(report, "CLAUDE.md"); a != ActionFailed {
		t.Errorf("action = %v, want failed", a)
	}
}

func TestUpgradeProjectAssets_claudeGuideHalfMarkerPairKeepsUserText(t *testing.T) {
	for name, user := range map[string]string{
		"lone begin": "# Mine\n" + managedBegin + "\nKEEP ME\n",
		"lone end":   "# Mine\nKEEP ME\n" + managedEnd + "\n",
		"end first":  managedEnd + "\nKEEP ME\n" + managedBegin + "\n",
	} {
		t.Run(name, func(t *testing.T) {
			target := savepointProject(t)
			path := filepath.Join(target, "CLAUDE.md")
			testutil.WriteFile(t, path, user)

			for i := 0; i < 2; i++ {
				report := mustUpgrade(t, claudeTemplates(), target)
				if a, _ := actionFor(report, "CLAUDE.md"); a != ActionUnchanged {
					t.Errorf("run %d: action = %v, want unchanged", i+1, a)
				}
				if got := string(mustReadFile(t, path)); got != user {
					t.Fatalf("run %d: CLAUDE.md = %q, want it untouched", i+1, got)
				}
			}
		})
	}
}
