package run

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

const stopBudget = 15 * time.Second // O-045 condition 3

// fakeStoppable is the fake agent for stop tests. idle waits with no child, as
// an agent does between turns; tree starts a grandchild, as an agent does
// mid-command, and records its pid in FAKE_PIDFILE; stubborn is tree that
// ignores SIGINT; sleeper is the grandchild.
func fakeStoppable(mode string) {
	if mode == "stubborn" {
		signal.Ignore(os.Interrupt)
	}
	if mode == "tree" || mode == "stubborn" {
		child := exec.Command(os.Args[0])
		child.Env = append(os.Environ(), "FAKE_AGENT=sleeper")
		if err := child.Start(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(5)
		}
		path := os.Getenv("FAKE_PIDFILE")
		os.WriteFile(path+".tmp", []byte(strconv.Itoa(child.Process.Pid)), 0o644)
		os.Rename(path+".tmp", path)
	}
	fmt.Println(`{"type":"system"}`)
	time.Sleep(time.Minute)
}

// startStoppable starts the fake in mode and, for modes with a grandchild,
// waits until its pid is known. It reports 0 for idle.
func startStoppable(t *testing.T, mode string, tweak func(*Spec)) (*Run, int) {
	t.Helper()
	pidfile := filepath.Join(t.TempDir(), "grandchild.pid")
	t.Setenv("FAKE_PIDFILE", pidfile)
	spec := fakeSpec(t, project(t), mode)
	if tweak != nil {
		tweak(&spec)
	}
	r, err := Start(spec)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { killLeftovers(r) })
	pid := 0
	if mode != "idle" {
		pid = waitForPid(t, pidfile)
		t.Cleanup(func() { killPid(pid) })
	}
	waitForEvent(t, r.Dir)
	return r, pid
}

func waitForPid(t *testing.T, path string) int {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		if b, err := os.ReadFile(path); err == nil {
			if pid, err := strconv.Atoi(strings.TrimSpace(string(b))); err == nil {
				return pid
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("the fake agent never recorded its grandchild")
	return 0
}

// waitForEvent returns once the agent has written a line, so a stop lands on a
// running agent rather than one still starting.
func waitForEvent(t *testing.T, dir string) {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		if b, _ := os.ReadFile(filepath.Join(dir, eventsFile)); len(b) > 0 {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("the fake agent wrote no event")
}

func killLeftovers(r *Run) {
	if rec := readRecordQuiet(r.Dir); rec.PID > 0 && processAlive(rec.PID) {
		if p, err := os.FindProcess(rec.PID); err == nil {
			p.Kill()
		}
	}
}

func killPid(pid int) {
	if p, err := os.FindProcess(pid); err == nil {
		p.Kill()
	}
}

func readRecordQuiet(dir string) Record {
	var rec Record
	if b, err := os.ReadFile(filepath.Join(dir, recordFile)); err == nil {
		yaml.Unmarshal(b, &rec)
	}
	return rec
}

func waitFinal(t *testing.T, r *Run) Record {
	t.Helper()
	done := make(chan Record, 1)
	go func() { done <- r.Wait() }()
	select {
	case rec := <-done:
		return rec
	case <-time.After(stopBudget):
		t.Fatalf("the run was still going %v after the stop", stopBudget)
		return Record{}
	}
}

// assertGone fails if pid is still a running process. A killed process can be
// visible for a moment while the system reaps it, so it polls briefly.
func assertGone(t *testing.T, what string, pid int) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for processAlive(pid) {
		if time.Now().After(deadline) {
			t.Fatalf("%s (pid %d) survived the stop", what, pid)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func cancel(t *testing.T, r *Run) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(r.Dir, cancelFile), nil, 0o644); err != nil {
		t.Fatal(err)
	}
}

func assertStopped(t *testing.T, r *Run, final Record, state string) {
	t.Helper()
	if final.State != state || final.EndedAt == "" || final.ExitCode == nil {
		t.Fatalf("final = %+v, want state %s with ended_at", final, state)
	}
	if got := readRecord(t, r.Dir); got.State != state || got.EndedAt != final.EndedAt {
		t.Errorf("run.yml = %+v", got)
	}
	events, err := os.ReadFile(filepath.Join(r.Dir, eventsFile))
	if err != nil || string(events) != "{\"type\":\"system\"}\n" {
		t.Errorf("events.jsonl = %q, %v", events, err)
	}
}

func TestStop_cancelFileEndsAnAgentBetweenTurns(t *testing.T) {
	r, _ := startStoppable(t, "idle", nil)
	cancel(t, r)
	final := waitFinal(t, r)
	assertStopped(t, r, final, StateCancelled)
	assertGone(t, "agent", final.PID)
}

func TestStop_cancelFileEndsAnAgentMidCommandAndItsGrandchild(t *testing.T) {
	r, grandchild := startStoppable(t, "tree", nil)
	cancel(t, r)
	final := waitFinal(t, r)
	assertStopped(t, r, final, StateCancelled)
	assertGone(t, "agent", final.PID)
	assertGone(t, "grandchild", grandchild)
}

func TestStop_cancelKillsAnAgentThatIgnoresTheInterrupt(t *testing.T) {
	r, grandchild := startStoppable(t, "stubborn", func(s *Spec) { s.Grace = 300 * time.Millisecond })
	cancel(t, r)
	final := waitFinal(t, r)
	assertStopped(t, r, final, StateCancelled)
	assertGone(t, "agent", final.PID)
	assertGone(t, "grandchild", grandchild)
}

func TestStop_timeoutEndsTheRunAndItsGrandchild(t *testing.T) {
	r, grandchild := startStoppable(t, "tree", func(s *Spec) { s.Timeout = 400 * time.Millisecond })
	final := waitFinal(t, r)
	assertStopped(t, r, final, StateTimedOut)
	assertGone(t, "agent", final.PID)
	assertGone(t, "grandchild", grandchild)
}

func TestStop_aRunThatEndsFirstIsNotStopped(t *testing.T) {
	spec := fakeSpec(t, project(t), "ok")
	spec.Timeout = time.Hour
	r, err := Start(spec)
	if err != nil {
		t.Fatal(err)
	}
	if final := r.Wait(); final.State != StateFinished || final.EndedAt == "" {
		t.Errorf("final = %+v", final)
	}
}

func TestRunnerLost_needsRunningStaleHeartbeatAndDeadPid(t *testing.T) {
	now := time.Date(2026, 10, 11, 12, 0, 0, 0, time.UTC)
	at := func(ago time.Duration) string { return now.Add(-ago).Format(time.RFC3339) }
	dead := deadPid(t)
	alive := os.Getpid()

	cases := []struct {
		name string
		rec  Record
		want bool
	}{
		{"running, stale, dead", Record{State: StateRunning, UpdatedAt: at(time.Minute), PID: dead}, true},
		{"fresh heartbeat", Record{State: StateRunning, UpdatedAt: at(10 * time.Second), PID: dead}, false},
		{"exactly at the limit", Record{State: StateRunning, UpdatedAt: at(LostAfter), PID: dead}, false},
		{"process still alive", Record{State: StateRunning, UpdatedAt: at(time.Minute), PID: alive}, false},
		{"already finished", Record{State: StateFinished, UpdatedAt: at(time.Minute), PID: dead}, false},
		{"already cancelled", Record{State: StateCancelled, UpdatedAt: at(time.Minute), PID: dead}, false},
		{"no pid recorded", Record{State: StateRunning, UpdatedAt: at(time.Minute)}, false},
		{"unreadable time", Record{State: StateRunning, UpdatedAt: "yesterday", PID: dead}, false},
	}
	for _, c := range cases {
		if got := RunnerLost(c.rec, now); got != c.want {
			t.Errorf("%s: RunnerLost = %v, want %v", c.name, got, c.want)
		}
	}
}

func TestRunnerLost_neverWritesTheRecord(t *testing.T) {
	dir := project(t)
	rec := Record{State: StateRunning, UpdatedAt: "2026-10-11T11:00:00Z", PID: deadPid(t)}
	path := filepath.Join(dir, recordFile)
	if err := writeRecord(path, rec); err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(path)
	if !RunnerLost(rec, time.Date(2026, 10, 11, 12, 0, 0, 0, time.UTC)) {
		t.Fatal("expected a lost runner")
	}
	if after, _ := os.ReadFile(path); string(after) != string(before) {
		t.Errorf("run.yml changed:\n%s\n%s", before, after)
	}
}

// deadPid is the pid of a process that has started and ended.
func deadPid(t *testing.T) int {
	t.Helper()
	t.Setenv("FAKE_AGENT", "quit")
	cmd := exec.Command(os.Args[0])
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
	return cmd.Process.Pid
}
