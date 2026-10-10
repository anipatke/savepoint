package init

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const hooksDir = ".claude/hooks"

func scaffoldedProject(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := Scaffold(shippedTemplates(t), dir, "p", false); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestScaffold_installsTrackedHooksAndSettings(t *testing.T) {
	dir := scaffoldedProject(t)
	manifest, err := LoadManifest(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"savepoint-find.js", "session-start.js", "guard.js"} {
		path := hooksDir + "/" + name
		content := readProjectFile(t, dir, path)
		if hash, ok := manifest.Hash(path); !ok || hash != hashContent([]byte(content)) {
			t.Errorf("manifest entry for %s missing or wrong", path)
		}
	}
	if _, ok := manifest.Hash(claudeSettingsPath); ok {
		t.Error("settings.json must not be manifest-tracked")
	}

	var settings struct {
		Hooks map[string][]struct {
			Hooks []struct {
				Type    string   `json:"type"`
				Command string   `json:"command"`
				Args    []string `json:"args"`
			} `json:"hooks"`
		} `json:"hooks"`
	}
	if err := json.Unmarshal([]byte(readProjectFile(t, dir, claudeSettingsPath)), &settings); err != nil {
		t.Fatalf("settings.json is not valid JSON: %v", err)
	}
	hook := settings.Hooks["SessionStart"][0].Hooks[0]
	if hook.Type != "command" || hook.Command != "node" || len(hook.Args) != 1 ||
		!strings.HasPrefix(hook.Args[0], "${CLAUDE_PROJECT_DIR}/") || !strings.HasSuffix(hook.Args[0], "session-start.js") {
		t.Errorf("unexpected hook command: %+v", hook)
	}
	if !strings.Contains(sessionStartEntry, hook.Args[0]) {
		t.Error("advice entry does not match the shipped settings template")
	}
	if strings.Contains(readProjectFile(t, dir, claudeSettingsPath), "statusLine") {
		t.Error("settings.json must not set a status line")
	}
}

func TestScaffold_keepsExistingSettings(t *testing.T) {
	for _, force := range []bool{false, true} {
		dir := t.TempDir()
		mine := []byte("{\n  \"statusLine\": {\"type\": \"command\", \"command\": \"mine\"}\n}\n")
		if err := os.MkdirAll(filepath.Join(dir, ".claude"), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, claudeSettingsPath), mine, 0644); err != nil {
			t.Fatal(err)
		}
		if err := Scaffold(shippedTemplates(t), dir, "p", force); err != nil {
			t.Fatal(err)
		}
		if got := mustReadFile(t, filepath.Join(dir, claudeSettingsPath)); string(got) != string(mine) {
			t.Errorf("force=%v: settings.json changed", force)
		}
		if advice := ClaudeSettingsAdvice(dir); !strings.Contains(advice, "session-start.js") {
			t.Errorf("force=%v: advice should name the entry to add, got %q", force, advice)
		}
	}
}

func TestClaudeSettingsAdvice_silentWhenAbsentOrWired(t *testing.T) {
	if got := ClaudeSettingsAdvice(t.TempDir()); got != "" {
		t.Errorf("absent file: advice = %q", got)
	}
	if got := ClaudeSettingsAdvice(scaffoldedProject(t)); got != "" {
		t.Errorf("wired file: advice = %q", got)
	}
}

func settingsEntry(t *testing.T, dir string, dryRun bool) UpgradeEntry {
	t.Helper()
	return upgradedPointerEntry(t, dir, dryRun, claudeSettingsPath)
}

func TestUpgrade_claudeSettings(t *testing.T) {
	t.Run("absent file is installed", func(t *testing.T) {
		dir := scaffoldedProject(t)
		if err := os.Remove(filepath.Join(dir, claudeSettingsPath)); err != nil {
			t.Fatal(err)
		}
		if e := settingsEntry(t, dir, false); e.Action != ActionInstalled {
			t.Errorf("action = %s, want installed", e.Action)
		}
		if _, err := os.Stat(filepath.Join(dir, claudeSettingsPath)); err != nil {
			t.Errorf("settings not installed: %v", err)
		}
	})

	t.Run("dry run writes nothing", func(t *testing.T) {
		dir := scaffoldedProject(t)
		if err := os.Remove(filepath.Join(dir, claudeSettingsPath)); err != nil {
			t.Fatal(err)
		}
		if e := settingsEntry(t, dir, true); e.Action != ActionInstalled {
			t.Errorf("action = %s, want installed", e.Action)
		}
		if _, err := os.Stat(filepath.Join(dir, claudeSettingsPath)); !os.IsNotExist(err) {
			t.Error("dry run wrote settings.json")
		}
	})

	t.Run("wired file is unchanged", func(t *testing.T) {
		dir := scaffoldedProject(t)
		if e := settingsEntry(t, dir, false); e.Action != ActionUnchanged {
			t.Errorf("action = %s, want unchanged", e.Action)
		}
	})

	t.Run("user file is byte-identical and the entry is reported", func(t *testing.T) {
		dir := scaffoldedProject(t)
		mine := []byte("{ // mine\r\n  \"statusLine\": 1 }")
		if err := os.WriteFile(filepath.Join(dir, claudeSettingsPath), mine, 0644); err != nil {
			t.Fatal(err)
		}
		for _, dryRun := range []bool{true, false} {
			e := settingsEntry(t, dir, dryRun)
			if e.Action != ActionInfo || !strings.Contains(e.Note, sessionStartEntry) {
				t.Errorf("dryRun=%v: entry = %+v, want info naming the entry to add", dryRun, e)
			}
			if got := mustReadFile(t, filepath.Join(dir, claudeSettingsPath)); string(got) != string(mine) {
				t.Errorf("dryRun=%v: user settings changed", dryRun)
			}
		}
		if _, err := os.Stat(filepath.Join(dir, claudeSettingsPath+incomingSuffix)); err == nil {
			t.Error("no sidecar should be written for settings.json")
		}
	})
}

func TestUpgrade_claudeHooksFollowSkillPolicy(t *testing.T) {
	path := hooksDir + "/session-start.js"
	dir := scaffoldedProject(t)
	mine := []byte("// mine\n")
	if err := os.WriteFile(filepath.Join(dir, path), mine, 0644); err != nil {
		t.Fatal(err)
	}
	if e := upgradedPointerEntry(t, dir, false, path); e.Action != ActionConflict {
		t.Errorf("edited hook action = %s, want conflict", e.Action)
	}
	if got := mustReadFile(t, filepath.Join(dir, path)); string(got) != string(mine) {
		t.Error("edited hook was overwritten")
	}
	if err := os.Remove(filepath.Join(dir, path)); err != nil {
		t.Fatal(err)
	}
	if e := upgradedPointerEntry(t, dir, false, path); e.Action != ActionUpdated {
		t.Errorf("missing hook action = %s, want updated", e.Action)
	}
}

// runHook runs the shipped session-start script with node. It skips when node
// is unavailable, as the hook itself would be.
func runHook(t *testing.T, projectDir string, pathDirs ...string) (string, int) {
	t.Helper()
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node is not installed")
	}
	script, err := filepath.Abs(filepath.Join("..", "..", "templates", "project-v2", ".claude", "hooks", "session-start.js"))
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(node, script)
	cmd.Env = append(os.Environ(), "CLAUDE_PROJECT_DIR="+projectDir, "PATH="+strings.Join(pathDirs, string(os.PathListSeparator)))
	if runtime.GOOS == "windows" {
		cmd.Env = append(cmd.Env, "Path="+strings.Join(pathDirs, string(os.PathListSeparator)))
	}
	out, err := cmd.Output()
	code := 0
	if exitErr, ok := err.(*exec.ExitError); ok {
		code = exitErr.ExitCode()
	} else if err != nil {
		t.Fatal(err)
	}
	return string(out), code
}

// fakeName is the file name the hook looks for on this platform.
func fakeName() string {
	if runtime.GOOS == "windows" {
		return "savepoint.cmd"
	}
	return "savepoint"
}

// fakeSavepoint writes a savepoint executable that prints lines and exits with
// code, into a new directory named dirName (a temp directory when empty). On
// Windows it is a .cmd batch file, which is what npm installs there.
func fakeSavepoint(t *testing.T, dirName string, code int, lines ...string) string {
	t.Helper()
	dir := t.TempDir()
	if dirName != "" {
		dir = filepath.Join(dir, dirName)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
	}
	var body string
	if runtime.GOOS == "windows" {
		body = "@echo off\r\n"
		for _, l := range lines {
			body += "echo " + l + "\r\n"
		}
		body += fmt.Sprintf("exit /b %d\r\n", code)
	} else {
		body = "#!/bin/sh\n"
		for _, l := range lines {
			body += "echo '" + l + "'\n"
		}
		body += fmt.Sprintf("exit %d\n", code)
	}
	if err := os.WriteFile(filepath.Join(dir, fakeName()), []byte(body), 0755); err != nil {
		t.Fatal(err)
	}
	return dir
}

func hookProject(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".savepoint"), 0755); err != nil {
		t.Fatal(err)
	}
	return dir
}

func nodeDir(t *testing.T) string {
	t.Helper()
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node is not installed")
	}
	// Link node into a directory of its own so a savepoint installed beside it
	// on the developer's machine cannot take part in the lookup under test.
	dir := t.TempDir()
	if err := os.Symlink(node, filepath.Join(dir, filepath.Base(node))); err != nil {
		return filepath.Dir(node)
	}
	return dir
}

func TestSessionStartHook(t *testing.T) {
	node := nodeDir(t)

	t.Run("adds only the Next line", func(t *testing.T) {
		bin := fakeSavepoint(t, "", 0, "Goal: x", "Next action: Start Task T-9.", "more")
		out, code := runHook(t, hookProject(t), bin, node)
		if code != 0 {
			t.Fatalf("exit = %d", code)
		}
		var got struct {
			Out struct {
				Event   string `json:"hookEventName"`
				Context string `json:"additionalContext"`
			} `json:"hookSpecificOutput"`
		}
		if err := json.Unmarshal([]byte(out), &got); err != nil {
			t.Fatalf("output is not JSON: %q", out)
		}
		if got.Out.Event != "SessionStart" || got.Out.Context != "Savepoint Next action: Start Task T-9." {
			t.Errorf("unexpected output: %q", out)
		}
	})

	silent := map[string]func(t *testing.T) (string, []string){
		"savepoint missing": func(t *testing.T) (string, []string) { return "", nil },
		"resume fails": func(t *testing.T) (string, []string) {
			return fakeSavepoint(t, "", 3, "Next action: no"), nil
		},
		"no Next line": func(t *testing.T) (string, []string) {
			return fakeSavepoint(t, "", 0, "hello"), nil
		},
	}
	for name, setup := range silent {
		t.Run(name, func(t *testing.T) {
			bin, _ := setup(t)
			dirs := []string{node}
			if bin != "" {
				dirs = []string{bin, node}
			}
			out, code := runHook(t, hookProject(t), dirs...)
			if out != "" || code != 0 {
				t.Errorf("out=%q exit=%d, want silent success", out, code)
			}
		})
	}

	t.Run("not a Savepoint project", func(t *testing.T) {
		bin := fakeSavepoint(t, "", 0, "Next action: x")
		out, code := runHook(t, t.TempDir(), bin, node)
		if out != "" || code != 0 {
			t.Errorf("out=%q exit=%d, want silent success", out, code)
		}
	})

	t.Run("runs savepoint from a directory with a space", func(t *testing.T) {
		bin := fakeSavepoint(t, "First Last", 0, "Next action: spaced")
		out, code := runHook(t, hookProject(t), bin, node)
		if code != 0 || !strings.Contains(out, "Next action: spaced") {
			t.Errorf("out=%q exit=%d", out, code)
		}
	})

	t.Run("finds savepoint in node_modules/.bin", func(t *testing.T) {
		project := hookProject(t)
		binDir := filepath.Join(project, "node_modules", ".bin")
		if err := os.MkdirAll(binDir, 0755); err != nil {
			t.Fatal(err)
		}
		fake := fakeSavepoint(t, "", 0, "Next action: local")
		if err := os.Rename(filepath.Join(fake, fakeName()), filepath.Join(binDir, fakeName())); err != nil {
			t.Fatal(err)
		}
		out, _ := runHook(t, project, node)
		if !strings.Contains(out, "Next action: local") {
			t.Errorf("out = %q", out)
		}
	})
}

func TestFindSavepoint_windowsNames(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node is not installed")
	}
	helper, err := filepath.Abs(filepath.Join("..", "..", "templates", "project-v2", ".claude", "hooks", "savepoint-find.js"))
	if err != nil {
		t.Fatal(err)
	}
	project := t.TempDir()
	bin := filepath.Join(project, "node_modules", ".bin")
	if err := os.MkdirAll(bin, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(bin, "savepoint.cmd"), []byte("@echo off\r\n"), 0644); err != nil {
		t.Fatal(err)
	}
	script := `const {findSavepoint}=require(process.argv[1]);` +
		`console.log(findSavepoint(process.argv[2],{PATH:""},"win32")||"none")`
	out, err := exec.Command(node, "-e", script, helper, project).Output()
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSpace(string(out)); got != filepath.Join(bin, "savepoint.cmd") {
		t.Errorf("found %q, want the .cmd shim in node_modules/.bin", got)
	}
}

func TestRepoHooksAreByteIdenticalToTemplates(t *testing.T) {
	for _, name := range []string{"savepoint-find.js", "session-start.js", "guard.js"} {
		want := mustReadFile(t, filepath.Join("..", "..", "templates", "project-v2", ".claude", "hooks", name))
		got := mustReadFile(t, filepath.Join("..", "..", ".claude", "hooks", name))
		if string(got) != string(want) {
			t.Errorf(".claude/hooks/%s differs from its template", name)
		}
	}
}
