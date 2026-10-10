package run

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// I-149: the runner, not the agent, decides whether a run is lost.
func TestRunnerLost_watchesTheRunnerNotTheAgent(t *testing.T) {
	now := time.Date(2026, 10, 11, 12, 0, 0, 0, time.UTC)
	at := func(ago time.Duration) string { return now.Add(-ago).Format(time.RFC3339) }
	dead, alive := deadPid(t), os.Getpid()
	cases := []struct {
		name string
		rec  Record
		want bool
	}{
		{"dead runner, live agent", Record{State: StateRunning, UpdatedAt: at(time.Minute), RunnerPID: dead, PID: alive}, true},
		{"live runner, dead agent", Record{State: StateRunning, UpdatedAt: at(time.Minute), RunnerPID: alive, PID: dead}, false},
		{"live both", Record{State: StateRunning, UpdatedAt: at(time.Minute), RunnerPID: alive, PID: alive}, false},
		{"dead both", Record{State: StateRunning, UpdatedAt: at(time.Minute), RunnerPID: dead, PID: dead}, true},
		{"dead runner, fresh", Record{State: StateRunning, UpdatedAt: at(10 * time.Second), RunnerPID: dead, PID: alive}, false},
		{"dead runner, exactly at the limit", Record{State: StateRunning, UpdatedAt: at(LostAfter), RunnerPID: dead, PID: alive}, false},
		{"dead runner, over the limit", Record{State: StateRunning, UpdatedAt: at(LostAfter + time.Second), RunnerPID: dead, PID: alive}, true},
	}
	for _, c := range cases {
		if got := RunnerLost(c.rec, now); got != c.want {
			t.Errorf("%s: RunnerLost = %v, want %v", c.name, got, c.want)
		}
	}
}

func TestStart_recordsTheRunnerPidSeparately(t *testing.T) {
	r, err := Start(fakeSpec(t, project(t), "ok"))
	if err != nil {
		t.Fatal(err)
	}
	final := r.Wait()
	if final.RunnerPID != os.Getpid() || final.PID == final.RunnerPID {
		t.Errorf("runner_pid = %d, pid = %d, want runner %d and a different agent pid", final.RunnerPID, final.PID, os.Getpid())
	}
}

// I-150: unknown or malformed shapes are counted, known housekeeping is not.
func TestParseClaudeLine_countsUnknownSubtypesAndMalformedMessages(t *testing.T) {
	unparsed := []string{
		`{"type":"system","subtype":"future_unknown"}`,
		`{"type":"system"}`,
		`{"type":"assistant","message":null}`,
		`{"type":"assistant"}`,
		`{"type":"assistant","message":{"content":7}}`,
		`{"type":"assistant","message":{"content":null}}`,
		`{"type":"assistant","message":{}}`,
		`{"type":"user","message":"text"}`,
	}
	for _, l := range unparsed {
		if got := ParseClaudeLine([]byte(l)); len(got) != 1 || got[0].Kind != KindUnparsed {
			t.Errorf("%s = %+v, want one unparsed event", l, got)
		}
	}
	ignored := []string{
		`{"type":"system","subtype":"hook_started"}`,
		`{"type":"system","subtype":"thinking_tokens"}`,
		`{"type":"assistant","message":{"content":"plain text"}}`,
		`{"type":"assistant","message":{"content":[]}}`,
		`{"type":"stream_event","event":{}}`,
	}
	for _, l := range ignored {
		if got := ParseClaudeLine([]byte(l)); len(got) != 0 {
			t.Errorf("%s = %+v, want none", l, got)
		}
	}
}

func TestStart_replayedBypassShapesIncrementUnparsedAndStayVerbatim(t *testing.T) {
	lines := `{"type":"system","subtype":"future_unknown"}
{"type":"assistant","message":null}
{"type":"assistant","message":{"content":7}}
{"type":"assistant","message":{"content":"hello"}}
{"type":"system","subtype":"hook_started"}
{"type":"stream_event","event":{}}
`
	replayFile := filepath.Join(t.TempDir(), "in.jsonl")
	if err := os.WriteFile(replayFile, []byte(lines), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("FAKE_REPLAY", replayFile)
	spec := fakeSpec(t, project(t), "replay")
	spec.Provider = ProviderClaude
	r, err := Start(spec)
	if err != nil {
		t.Fatal(err)
	}
	final := r.Wait()
	got, _ := os.ReadFile(filepath.Join(r.Dir, eventsFile))
	if string(got) != lines || final.Unparsed != 3 {
		t.Errorf("unparsed = %d (want 3), events verbatim = %v", final.Unparsed, string(got) == lines)
	}
}

// I-151: the saved command is cleaned; the child gets the original argv.
func TestStart_savedCommandHasCredentialsRemovedButChildGetsOriginal(t *testing.T) {
	spec := fakeSpec(t, project(t), "ok")
	spec.Argv = append(spec.Argv, "https://user:checksecret@example.com/x", "a\x1b[31mb")
	r, err := Start(spec)
	if err != nil {
		t.Fatal(err)
	}
	r.Wait()
	b, _ := os.ReadFile(filepath.Join(r.Dir, recordFile))
	if strings.Contains(string(b), "checksecret") || strings.Contains(string(b), "\x1b") {
		t.Errorf("run.yml keeps a credential or control sequence:\n%s", b)
	}
	if !strings.Contains(string(b), "https://example.com/x") {
		t.Errorf("run.yml lost the cleaned URL:\n%s", b)
	}
	if spec.Argv[len(spec.Argv)-2] != "https://user:checksecret@example.com/x" {
		t.Error("the caller's argv was changed")
	}
}

// I-152: a record that cannot be saved is reported, not called a success.
func TestStart_summaryWriteFailureIsReported(t *testing.T) {
	spec := fakeSpec(t, project(t), "slow")
	r, err := Start(spec)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(r.Dir, summaryFile), 0o755); err != nil {
		t.Fatal(err)
	}
	r.Wait()
	if r.Err() == nil {
		t.Error("Err() = nil after the summary could not be written")
	}
}

func TestStart_heartbeatWriteFailureStopsTheAgentAndIsReported(t *testing.T) {
	spec := fakeSpec(t, project(t), "idle")
	spec.Heartbeat = 20 * time.Millisecond
	spec.Timeout = time.Minute
	r, err := Start(spec)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(r.Dir, recordFile)
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(path, 0o755); err != nil {
		t.Fatal(err)
	}
	done := make(chan struct{})
	go func() { r.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(15 * time.Second):
		t.Fatal("the agent kept running after the record could not be written")
	}
	if r.Err() == nil {
		t.Error("Err() = nil after a heartbeat write failed")
	}
}
