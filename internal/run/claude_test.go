package run

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/opencode/savepoint/internal/data"
)

const fixtures = "testdata/claude-2.1.296"

// replay starts the fake provider so that it prints the file as Claude Code
// would, and returns the finished run's directory and record.
func replay(t *testing.T, file string) (string, Record) {
	t.Helper()
	spec := fakeSpec(t, project(t), "replay")
	spec.Provider = ProviderClaude
	abs, err := filepath.Abs(file) // the child runs in the project directory
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("FAKE_REPLAY", abs)
	r, err := Start(spec)
	if err != nil {
		t.Fatal(err)
	}
	return r.Dir, r.Wait()
}

func ptr[T any](v T) *T { return &v }

func TestReplay_realSessionWithEditAndDenialFillsTheRecord(t *testing.T) {
	file := filepath.Join(fixtures, "edit-and-denial.jsonl")
	dir, rec := replay(t, file)

	if rec.State != StateFinished {
		t.Fatalf("state = %q", rec.State)
	}
	if rec.SessionID != "99339342-d9e2-4a4f-a921-4a9f4469028b" || rec.Model != "claude-haiku-5-5" {
		t.Errorf("session/model = %q / %q", rec.SessionID, rec.Model)
	}
	if len(rec.FilesRead) != 1 || !strings.HasSuffix(rec.FilesRead[0], "/cap/notes.txt") {
		t.Errorf("files_read = %v", rec.FilesRead)
	}
	if len(rec.FilesEdited) != 1 || rec.FilesEdited[0] != rec.FilesRead[0] {
		t.Errorf("files_edited = %v", rec.FilesEdited)
	}
	wantCommands := []CommandRun{
		{Command: "echo done", ExitCode: ptr(0)},
		{Command: "ls /nonexistent-dir"}, // refused: it never ran, so no exit code
		{Command: "git status", ExitCode: ptr(0)},
	}
	if !reflect.DeepEqual(rec.Commands, wantCommands) {
		t.Errorf("commands = %+v, want %+v", rec.Commands, wantCommands)
	}
	if len(rec.Denials) != 1 || rec.Denials[0].Tool != "Bash" || rec.Denials[0].Target != "ls /nonexistent-dir" ||
		!strings.Contains(rec.Denials[0].Reason, "denied automatically") || rec.Denials[0].ToolUseID == "" {
		t.Errorf("denials = %+v", rec.Denials)
	}
	if rec.Usage["output_tokens"] != int64(2562) || rec.Usage["cache_read_input_tokens"] != int64(138195) {
		t.Errorf("usage = %v", rec.Usage)
	}
	if rec.TotalCostUSD == nil || *rec.TotalCostUSD <= 0 {
		t.Errorf("total_cost_usd = %v", rec.TotalCostUSD)
	}
	wantReport := "Four of the five steps ran, but the `ls /nonexistent-dir` check was blocked by a permission denial, so that step did not run."
	if rec.AgentReported != wantReport {
		t.Errorf("agent_reported = %q", rec.AgentReported)
	}
	if rec.Unparsed != 1 { // the rate_limit_event line
		t.Errorf("unparsed = %d, want 1", rec.Unparsed)
	}

	// The record on disk matches, and events.jsonl is the provider's bytes.
	wantYAML, _ := yaml.Marshal(rec)
	if gotYAML, _ := os.ReadFile(filepath.Join(dir, "run.yml")); string(gotYAML) != string(wantYAML) {
		t.Errorf("run.yml differs from the final record:\n got %s\nwant %s", gotYAML, wantYAML)
	}
	want, _ := os.ReadFile(file)
	if got, _ := os.ReadFile(filepath.Join(dir, "events.jsonl")); string(got) != string(want) {
		t.Error("events.jsonl is not the fixture verbatim")
	}
	summary, _ := os.ReadFile(filepath.Join(dir, "summary.md"))
	if !strings.HasPrefix(string(summary), "# Agent-reported summary\n") || !strings.Contains(string(summary), "\n"+wantReport+"\n") ||
		!strings.Contains(string(summary), "did not check it") {
		t.Errorf("summary.md = %q", summary)
	}
}

func TestReplay_failingCommandRecordsItsExitCode(t *testing.T) {
	_, rec := replay(t, filepath.Join(fixtures, "failing-command.jsonl"))
	wantCommands := []CommandRun{
		{Command: "echo done", ExitCode: ptr(0)},
		{Command: "ls missing-dir", ExitCode: ptr(2)},
		{Command: "git status", ExitCode: ptr(0)},
	}
	if !reflect.DeepEqual(rec.Commands, wantCommands) {
		t.Errorf("commands = %+v, want %+v", rec.Commands, wantCommands)
	}
	if len(rec.Denials) != 0 {
		t.Errorf("denials = %+v, want none", rec.Denials)
	}
	if !strings.Contains(rec.AgentReported, "failed as expected") {
		t.Errorf("agent_reported = %q", rec.AgentReported)
	}
}

func TestReplay_unknownAndMalformedLinesStayOnDiskAndAreCounted(t *testing.T) {
	lines := []string{
		`{"type":"system","subtype":"init","session_id":"s-1","model":"m"}`,
		`{"type":"brand_new_event","x":1}`,
		`{"type":"assistant","message":`, // cut off mid-write
		`not json at all`,
		`{"type":"stream_event","event":{"type":"content_block_delta"}}`, // ignored on purpose
		`{"type":"result","session_id":"s-1","result":"done"}`,
	}
	file := filepath.Join(t.TempDir(), "stream.jsonl")
	if err := os.WriteFile(file, []byte(strings.Join(lines, "\n")+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	dir, rec := replay(t, file)
	if rec.State != StateFinished || rec.Unparsed != 3 || rec.SessionID != "s-1" || rec.AgentReported != "done" {
		t.Errorf("record = %+v", rec)
	}
	got, _ := os.ReadFile(filepath.Join(dir, "events.jsonl"))
	if string(got) != strings.Join(lines, "\n")+"\n" {
		t.Errorf("events.jsonl = %q", got)
	}
}

func TestReplay_aRunWithoutAResultLineStillGetsASummaryAndNoReport(t *testing.T) {
	file := filepath.Join(t.TempDir(), "stream.jsonl")
	os.WriteFile(file, []byte(`{"type":"system","subtype":"init","session_id":"s-2"}`+"\n"), 0o644)
	dir, rec := replay(t, file)
	if rec.AgentReported != "" || rec.SessionID != "s-2" {
		t.Errorf("record = %+v", rec)
	}
	summary, _ := os.ReadFile(filepath.Join(dir, "summary.md"))
	if !strings.Contains(string(summary), "reported no final message") {
		t.Errorf("summary.md = %q", summary)
	}
}

func TestReplay_textCopiedIntoTheRecordHasControlSequencesAndCredentialsRemoved(t *testing.T) {
	lines := []string{
		`{"type":"assistant","message":{"content":[{"type":"tool_use","id":"t1","name":"Bash","input":{"command":"curl https://user:pw@example.com/x \u001b[31mred\u001b[0m\u0007"}}]}}`,
		`{"type":"result","session_id":"s\u001b[2J-3","result":"line one\u001b]0;title\u0007\nsee https://bob:hunter2@host.example/p\u001b[1mbold\u001b[0m\r\nend"}`,
	}
	file := filepath.Join(t.TempDir(), "stream.jsonl")
	os.WriteFile(file, []byte(strings.Join(lines, "\n")+"\n"), 0o644)
	dir, rec := replay(t, file)

	if len(rec.Commands) != 1 || rec.Commands[0].Command != "curl https://example.com/x red" {
		t.Errorf("commands = %+v", rec.Commands)
	}
	if rec.SessionID != "s-3" {
		t.Errorf("session_id = %q", rec.SessionID)
	}
	if want := "line one\nsee https://host.example/pbold\nend"; rec.AgentReported != want {
		t.Errorf("agent_reported = %q, want %q", rec.AgentReported, want)
	}
	for _, name := range []string{"run.yml", "summary.md"} {
		b, _ := os.ReadFile(filepath.Join(dir, name))
		if strings.ContainsAny(string(b), "\x1b\x07\r") || strings.Contains(string(b), "hunter2") || strings.Contains(string(b), "pw@") {
			t.Errorf("%s keeps a control sequence or credential:\n%s", name, b)
		}
	}
}

func TestParseClaudeLine_mapsEachDocumentedLineType(t *testing.T) {
	cases := []struct {
		name string
		line string
		want []Event
	}{
		{"init", `{"type":"system","subtype":"init","session_id":"s","model":"m"}`,
			[]Event{{Kind: KindInit, SessionID: "s", Model: "m"}}},
		{"read", `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"a","name":"Read","input":{"file_path":"/p/f.go"}}]}}`,
			[]Event{{Kind: KindRead, ToolUseID: "a", Tool: "Read", Text: "/p/f.go"}}},
		{"write is an edit", `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"a","name":"Write","input":{"file_path":"/p/f.go"}}]}}`,
			[]Event{{Kind: KindEdit, ToolUseID: "a", Tool: "Write", Text: "/p/f.go"}}},
		{"notebook edit", `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"a","name":"NotebookEdit","input":{"notebook_path":"/p/n.ipynb"}}]}}`,
			[]Event{{Kind: KindEdit, ToolUseID: "a", Tool: "NotebookEdit", Text: "/p/n.ipynb"}}},
		{"command", `{"type":"assistant","message":{"content":[{"type":"text","text":"ok"},{"type":"tool_use","id":"a","name":"Bash","input":{"command":"make  test\n"}}]}}`,
			[]Event{{Kind: KindCommand, ToolUseID: "a", Tool: "Bash", Text: "make test"}}},
		{"other tool", `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"a","name":"Grep","input":{"pattern":"x"}}]}}`,
			[]Event{{Kind: KindTool, ToolUseID: "a", Tool: "Grep", Text: "Grep"}}},
		{"read without a path", `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"a","name":"Read","input":{}}]}}`,
			[]Event{{Kind: KindTool, ToolUseID: "a", Tool: "Read", Text: "Read"}}},
		{"failed result with exit code", `{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"a","is_error":true,"content":"Exit code 127\nnope"}]}}`,
			[]Event{{Kind: KindToolResult, ToolUseID: "a", Failed: true, ExitCode: ptr(127)}}},
		{"result text in blocks", `{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"a","is_error":true,"content":[{"type":"text","text":"Exit code 1\nx"}]}]}}`,
			[]Event{{Kind: KindToolResult, ToolUseID: "a", Failed: true, ExitCode: ptr(1)}}},
		{"permission denied", `{"type":"system","subtype":"permission_denied","tool_name":"Bash","tool_use_id":"a","decision_reason":"no approval surface","message":"long text"}`,
			[]Event{{Kind: KindDenied, ToolUseID: "a", Tool: "Bash", Text: "no approval surface"}}},
		{"final", `{"type":"result","session_id":"s","total_cost_usd":0.5,"usage":{"output_tokens":3000000,"tier":"std"},"permission_denials":[{"tool_name":"Bash","tool_use_id":"a","tool_input":{"command":"git push"}}],"result":"all done"}`,
			[]Event{{Kind: KindFinal, SessionID: "s", Text: "all done", TotalCostUSD: ptr(0.5),
				Usage:   map[string]any{"output_tokens": int64(3000000), "tier": "std"},
				Denials: []Denial{{Tool: "Bash", Target: "git push", ToolUseID: "a"}}}}},
		{"plain-text user message has no blocks", `{"type":"user","message":{"content":"hello"}}`, nil},
		{"hooks are housekeeping", `{"type":"system","subtype":"hook_started"}`, nil},
		{"partial deltas are ignored", `{"type":"stream_event","event":{"type":"content_block_delta"}}`, nil},
		{"blank line", `   `, nil},
		{"unknown type", `{"type":"rate_limit_event"}`, []Event{{Kind: KindUnparsed}}},
		{"no type", `{"hello":1}`, []Event{{Kind: KindUnparsed}}},
		{"null", `null`, []Event{{Kind: KindUnparsed}}},
		{"malformed", `{"type":"assistant"`, []Event{{Kind: KindUnparsed}}},
		{"message is not an object", `{"type":"assistant","message":7}`, []Event{{Kind: KindUnparsed}}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := ParseClaudeLine([]byte(c.line)); !reflect.DeepEqual(got, c.want) {
				t.Errorf("got  %+v\nwant %+v", got, c.want)
			}
		})
	}
}

func TestCleanText_removesSequencesButKeepsLines(t *testing.T) {
	cases := map[string]string{
		"a\x1b[1;31mb\x1b[0mc":                  "abc",
		"a\x1b]8;;http://x\x07link\x1b]8;;\x07": "alink",
		"a\u009b31mb":                           "ab",
		"tab\there\nnext\r\nline\x00\x7f":       "tab\there\nnext\nline",
		"u https://a:b@h/p":                     "u https://h/p",
	}
	for in, want := range cases {
		if got := cleanText(in); got != want {
			t.Errorf("cleanText(%q) = %q, want %q", in, got, want)
		}
	}
	if got := oneLine(strings.Repeat("é", 400), 11); got != "ééééé…" {
		t.Errorf("oneLine cut = %q", got)
	}
}

func TestBuildArgv(t *testing.T) {
	gates := data.QualityGates{Build: ptr("make build"), Test: ptr("make test-fast"), Lint: ptr("golangci-lint run ./...")}
	next := "Build T-128 — Read Claude Code's event stream (O-045) {{project_root}} \"quoted\" $HOME"

	t.Run("default builder substitutes the Next line verbatim", func(t *testing.T) {
		argv, err := BuildArgv(nil, gates, next, ".")
		if err != nil {
			t.Fatal(err)
		}
		cwd, _ := os.Getwd()
		want := []string{"claude", "-p", next, "--output-format", "stream-json", "--verbose",
			"--permission-mode", "acceptEdits", "--allowedTools", "Bash(golangci-lint *)", "Bash(make *)",
			"--permission-prompts", "none"}
		if !reflect.DeepEqual(argv, want) {
			t.Errorf("argv = %q\nwant   %q", argv, want)
		}
		if argv[2] != next || cwd == "" {
			t.Error("the prompt was altered")
		}
	})
	t.Run("configured builder gets both placeholders", func(t *testing.T) {
		root := t.TempDir()
		b := &data.CommandSpec{Command: "my-agent", Args: []string{"--cwd={{project_root}}", "{{prompt}}", "plain"}}
		argv, err := BuildArgv(b, gates, next, root)
		if err != nil {
			t.Fatal(err)
		}
		if want := []string{"my-agent", "--cwd=" + root, next, "plain"}; !reflect.DeepEqual(argv, want) {
			t.Errorf("argv = %q, want %q", argv, want)
		}
	})
	t.Run("relative root becomes absolute", func(t *testing.T) {
		b := &data.CommandSpec{Command: "a", Args: []string{"{{project_root}}"}}
		argv, _ := BuildArgv(b, gates, "x", "sub")
		if !filepath.IsAbs(argv[1]) || filepath.Base(argv[1]) != "sub" {
			t.Errorf("project root = %q", argv[1])
		}
	})
	t.Run("builder without a command fails with a named error", func(t *testing.T) {
		for _, b := range []*data.CommandSpec{{}, {Command: "  ", Args: []string{"{{prompt}}"}}} {
			argv, err := BuildArgv(b, gates, next, t.TempDir())
			if err == nil || !strings.Contains(err.Error(), "agent_launcher.builder.command") || argv != nil {
				t.Errorf("argv = %q, err = %v", argv, err)
			}
		}
	})
}
