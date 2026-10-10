package init

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const taskRel = ".savepoint/objectives/O-001-x/tasks/T-001-x.md"

func taskText(status string) string {
	return "---\nid: T-001\ntitle: X\nstatus: " + status + "\n---\n\n# X\n\nBody.\n"
}

// runGuard feeds a PreToolUse payload to the shipped guard script and returns
// the block reason, or "" when the call is allowed.
func runGuard(t *testing.T, stdin string, cwd string) string {
	t.Helper()
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node is not installed")
	}
	script, err := filepath.Abs(filepath.Join("..", "..", "templates", "project-v2", ".claude", "hooks", "guard.js"))
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(node, script)
	cmd.Dir = cwd
	cmd.Stdin = strings.NewReader(stdin)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("guard must exit 0, got %v", err)
	}
	if len(out) == 0 {
		return ""
	}
	var res struct {
		Hook struct {
			Event    string `json:"hookEventName"`
			Decision string `json:"permissionDecision"`
			Reason   string `json:"permissionDecisionReason"`
		} `json:"hookSpecificOutput"`
	}
	if err := json.Unmarshal(out, &res); err != nil {
		t.Fatalf("guard output is not JSON: %v: %s", err, out)
	}
	if res.Hook.Event != "PreToolUse" || res.Hook.Decision != "deny" {
		t.Fatalf("unexpected guard output: %s", out)
	}
	if strings.Contains(res.Hook.Reason, "\n") {
		t.Errorf("reason must be one line: %q", res.Hook.Reason)
	}
	return res.Hook.Reason
}

func payload(t *testing.T, tool string, input map[string]any, cwd string) string {
	t.Helper()
	b, err := json.Marshal(map[string]any{"tool_name": tool, "tool_input": input, "cwd": cwd})
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func writeTask(t *testing.T, dir, status string) string {
	t.Helper()
	path := filepath.Join(dir, filepath.FromSlash(taskRel))
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(taskText(status)), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestGuard_taskDone(t *testing.T) {
	dir := t.TempDir()
	path := writeTask(t, dir, "in_progress")
	cases := []struct {
		name  string
		tool  string
		input map[string]any
		block bool
	}{
		{"edit to done blocks", "Edit", map[string]any{"file_path": path, "old_string": "status: in_progress", "new_string": "status: done"}, true},
		{"write done blocks", "Write", map[string]any{"file_path": path, "content": taskText("done")}, true},
		{"multiedit to done blocks", "MultiEdit", map[string]any{"file_path": path, "edits": []map[string]any{{"old_string": "title: X", "new_string": "title: Y"}, {"old_string": "status: in_progress", "new_string": "status: done"}}}, true},
		{"quoted done blocks", "Edit", map[string]any{"file_path": path, "old_string": "status: in_progress", "new_string": "status: 'done'"}, true},
		{"edit to another status allows", "Edit", map[string]any{"file_path": path, "old_string": "status: in_progress", "new_string": "status: planned"}, false},
		{"edit to a body line allows", "Edit", map[string]any{"file_path": path, "old_string": "Body.", "new_string": "status: done is a phrase"}, false},
		{"other file allows", "Write", map[string]any{"file_path": filepath.Join(dir, "notes.md"), "content": taskText("done")}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			reason := runGuard(t, payload(t, c.tool, c.input, dir), dir)
			if c.block && !strings.Contains(reason, "status: done") {
				t.Errorf("want a status: done block, got %q", reason)
			}
			if !c.block && reason != "" {
				t.Errorf("want allow, got %q", reason)
			}
		})
	}

	t.Run("edits to a done task allow", func(t *testing.T) {
		done := writeTask(t, t.TempDir(), "done")
		in := map[string]any{"file_path": done, "old_string": "Body.", "new_string": "More."}
		if reason := runGuard(t, payload(t, "Edit", in, dir), dir); reason != "" {
			t.Errorf("want allow, got %q", reason)
		}
	})
}

func gitIn(t *testing.T, dir string, args ...string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}
	cmd := exec.Command("git", append([]string{"-c", "user.email=t@t", "-c", "user.name=t"}, args...)...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

func TestGuard_routerInLane(t *testing.T) {
	main := t.TempDir()
	gitIn(t, main, "init", "-q")
	if err := os.MkdirAll(filepath.Join(main, ".savepoint"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(main, ".savepoint", "router.md"), []byte("x\n"), 0644); err != nil {
		t.Fatal(err)
	}
	gitIn(t, main, "add", ".")
	gitIn(t, main, "commit", "-q", "-m", "init")
	lane := filepath.Join(t.TempDir(), "lane")
	gitIn(t, main, "worktree", "add", "-q", lane, "-b", "lane")

	edit := func(root string) string {
		return payload(t, "Edit", map[string]any{"file_path": filepath.Join(root, ".savepoint", "router.md"), "old_string": "x", "new_string": "y"}, root)
	}
	if reason := runGuard(t, edit(lane), lane); !strings.Contains(reason, "router.md") {
		t.Errorf("lane: want a router block, got %q", reason)
	}
	if reason := runGuard(t, edit(main), main); reason != "" {
		t.Errorf("main checkout: want allow, got %q", reason)
	}
	if reason := runGuard(t, edit(t.TempDir()), t.TempDir()); reason != "" {
		t.Errorf("not a git repo: want allow, got %q", reason)
	}
}

func TestGuard_savepointCommands(t *testing.T) {
	dir := t.TempDir()
	cases := []struct {
		command string
		block   bool
	}{
		{"savepoint resume", false},
		{"./savepoint resume", false},
		{"npx savepoint resume", false},
		{"npx --yes savepoint create-task --objective O-001 --draft d.md", false},
		{"savepoint health check O-001", false},
		{"cd x && savepoint resume | head -3", false},
		{"go test ./... && make build", false},
		{"echo savepoint board", false},
		{"cat savepoint.md", false},
		{"savepoint board", true},
		{"./savepoint init .", true},
		{"npx savepoint upgrade-assets", true},
		{"savepoint health setup --apply", true},
		{"savepoint health report", true},
		{"savepoint migrate", true},
		{"make build; savepoint doctor", true},
		{"FOO=1 /usr/bin/savepoint board", true},
		{"savepoint.exe board", true},
		{"savepoint", true},
		{"echo hi\nsavepoint board", true},
		{"true && (savepoint board)", true},
	}
	for _, c := range cases {
		t.Run(c.command, func(t *testing.T) {
			reason := runGuard(t, payload(t, "Bash", map[string]any{"command": c.command}, dir), dir)
			if c.block && !strings.Contains(reason, "owner") {
				t.Errorf("want a block, got %q", reason)
			}
			if !c.block && reason != "" {
				t.Errorf("want allow, got %q", reason)
			}
		})
	}
}

func TestGuard_failsOpen(t *testing.T) {
	dir := t.TempDir()
	for name, stdin := range map[string]string{
		"empty input":     "",
		"not json":        "not json",
		"no tool input":   `{"tool_name":"Edit"}`,
		"unknown tool":    `{"tool_name":"Read","tool_input":{"file_path":"` + taskRel + `"}}`,
		"missing file":    payload(t, "Edit", map[string]any{"file_path": filepath.Join(dir, taskRel), "old_string": "a", "new_string": "status: done"}, dir),
		"bash no command": `{"tool_name":"Bash","tool_input":{}}`,
	} {
		t.Run(name, func(t *testing.T) {
			if reason := runGuard(t, stdin, dir); reason != "" {
				t.Errorf("want allow, got %q", reason)
			}
		})
	}
}

func TestSettingsNote_namesOnlyMissingEntries(t *testing.T) {
	both := []byte("session-start.js guard.js")
	if note := settingsNote(both); note != "" {
		t.Errorf("both wired: note = %q", note)
	}
	onlyStart := settingsNote([]byte("session-start.js"))
	if !strings.Contains(onlyStart, guardEntry) || strings.Contains(onlyStart, sessionStartEntry) {
		t.Errorf("only the guard entry should be named: %q", onlyStart)
	}
	if note := settingsNote([]byte("{}")); !strings.Contains(note, guardEntry) || !strings.Contains(note, sessionStartEntry) {
		t.Errorf("both entries should be named: %q", note)
	}
}
