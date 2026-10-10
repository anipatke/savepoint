// Package run starts one configured agent command as a child process, appends
// its stdout verbatim to events.jsonl, and keeps run.yml current with what the
// runner observed. A provider adapter reads the stream into the record; no
// other code interprets provider output, and nothing here scores the agent's
// own report. It never writes a Task, Objective, Check, or router record.
package run

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sync"
	"time"
)

const (
	// DefaultHeartbeat is how often updated_at is rewritten while the child runs.
	DefaultHeartbeat = 5 * time.Second
	killWait         = 2 * time.Second

	recordFile  = "run.yml"
	eventsFile  = "events.jsonl"
	summaryFile = "summary.md"
)

// Spec is one run to start: an argument vector, never a shell string.
type Spec struct {
	Dir       string   // project root; the child's working directory
	Task      string   // Task ID, used in the run directory name
	Provider  string   // label recorded in run.yml
	Argv      []string // executable followed by its arguments
	NextLine  string   // the Next line the agent was given, recorded verbatim
	Adapter   Adapter  // reads stdout lines into the record; the Provider's own when nil
	Heartbeat time.Duration
	Timeout   time.Duration    // longest the run may last; DefaultTimeout when zero
	Grace     time.Duration    // time to exit after SIGINT before the kill; DefaultGrace when zero
	Now       func() time.Time // clock; time.Now when nil
}

// Run is a started run. Its record lives in Dir.
type Run struct {
	Dir string // the run directory

	mu     sync.Mutex
	rec    Record
	done   chan struct{}
	stopHB chan struct{}
	final  Record
	err    error // first failure to save the record; guarded by mu
}

var taskID = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]*$`)

// Start creates .savepoint/runs/<yyyymmddThhmmssZ>-<task>/, writes run.yml with
// state running, and starts the child. Start returns once the child is running;
// Wait returns the final record. If the child cannot be started the record is
// closed as failed and the error is returned.
func Start(spec Spec) (*Run, error) {
	if len(spec.Argv) == 0 || spec.Argv[0] == "" {
		return nil, errors.New("run: no command to start")
	}
	if !taskID.MatchString(spec.Task) {
		return nil, fmt.Errorf("run: invalid task %q", spec.Task)
	}
	if st, err := os.Stat(filepath.Join(spec.Dir, ".savepoint")); err != nil || !st.IsDir() {
		return nil, fmt.Errorf("run: %s is not a Savepoint project", spec.Dir)
	}
	now := spec.Now
	if now == nil {
		now = time.Now
	}
	heartbeat := spec.Heartbeat
	if heartbeat <= 0 {
		heartbeat = DefaultHeartbeat
	}

	timeout, grace := spec.Timeout, spec.Grace
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	if grace <= 0 {
		grace = DefaultGrace
	}

	started := now().UTC()
	runsDir := filepath.Join(spec.Dir, ".savepoint", "runs")
	if err := os.MkdirAll(runsDir, 0o755); err != nil {
		return nil, fmt.Errorf("run: create runs directory: %w", err)
	}
	dir := filepath.Join(runsDir, started.Format("20060102T150405Z")+"-"+spec.Task)
	if err := os.Mkdir(dir, 0o755); err != nil {
		return nil, fmt.Errorf("run: create run directory: %w", err)
	}

	r := &Run{Dir: dir, done: make(chan struct{}), stopHB: make(chan struct{})}
	r.rec = Record{
		Provider:  spec.Provider,
		Command:   cleanArgv(spec.Argv),
		RunnerPID: os.Getpid(),
		Cwd:       spec.Dir,
		Task:      spec.Task,
		NextLine:  spec.NextLine,
		StartedAt: started.Format(time.RFC3339),
		UpdatedAt: started.Format(time.RFC3339),
		State:     StateRunning,
	}
	if err := writeRecord(filepath.Join(dir, recordFile), r.rec); err != nil {
		return nil, err
	}
	events, err := os.OpenFile(filepath.Join(dir, eventsFile), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("run: open events: %w", err)
	}

	cmd := exec.CommandContext(context.Background(), spec.Argv[0], spec.Argv[1:]...)
	cmd.Dir = spec.Dir
	cmd.Env = append(os.Environ(), "CI=1", "NO_COLOR=1")
	cmd.WaitDelay = killWait
	configureProcess(cmd)
	adapter := spec.Adapter
	if adapter == nil {
		adapter = adapters[spec.Provider]
	}
	obs := newObserver()
	lines := &lineWriter{out: events, observe: func(line []byte) {
		if adapter == nil {
			return
		}
		evs := adapter(line)
		r.mu.Lock()
		for _, ev := range evs {
			obs.apply(&r.rec, ev)
		}
		r.mu.Unlock()
	}}
	stderr := &tailBuffer{limit: maxStderrBytes}
	cmd.Stdout, cmd.Stderr = lines, stderr

	if err := cmd.Start(); err != nil {
		events.Close()
		r.mu.Lock()
		r.rec.State = StateFailed
		r.rec.StderrTail = sanitizeTail(err.Error())
		r.rec.UpdatedAt = now().UTC().Format(time.RFC3339)
		r.final = r.rec
		r.mu.Unlock()
		close(r.done)
		if werr := writeRecord(filepath.Join(dir, recordFile), r.final); werr != nil {
			return nil, werr
		}
		return nil, fmt.Errorf("run: start %s: %w", spec.Argv[0], err)
	}

	r.mu.Lock()
	r.rec.PID = cmd.Process.Pid
	r.rec.UpdatedAt = now().UTC().Format(time.RFC3339)
	err = writeRecord(filepath.Join(dir, recordFile), r.rec)
	r.mu.Unlock()
	if err != nil {
		// The record cannot be kept current, so the run must not go on unrecorded.
		cmd.Cancel()
		cmd.Wait()
		events.Close()
		close(r.done)
		return nil, err
	}

	procDone := make(chan struct{})
	stop := startStopper(dir, cmd, timeout, grace, procDone)
	hbDone := make(chan struct{})
	go func() {
		defer close(hbDone)
		t := time.NewTicker(heartbeat)
		defer t.Stop()
		for {
			select {
			case <-r.stopHB:
				return
			case <-t.C:
				r.mu.Lock()
				r.rec.UpdatedAt = now().UTC().Format(time.RFC3339)
				err := writeRecord(filepath.Join(dir, recordFile), r.rec)
				if err != nil && r.err == nil {
					r.err = err
				}
				r.mu.Unlock()
				if err != nil {
					// The record cannot be kept current, so the agent must not go on unrecorded.
					killTree(cmd)
					return
				}
			}
		}
	}()

	go func() {
		waitErr := cmd.Wait()
		close(procDone)
		<-stop.done // a stop in progress finishes killing the tree first
		lines.flush()
		events.Close()
		close(r.stopHB)
		<-hbDone

		r.mu.Lock()
		defer r.mu.Unlock()
		ended := now().UTC().Format(time.RFC3339)
		r.rec.UpdatedAt, r.rec.EndedAt = ended, ended
		code := cmd.ProcessState.ExitCode()
		r.rec.ExitCode = &code
		var exit *exec.ExitError
		switch {
		case stop.reason() != "":
			r.rec.State = stop.reason()
		case waitErr == nil, errors.Is(waitErr, exec.ErrWaitDelay):
			r.rec.State = StateFinished
		default:
			r.rec.State = StateFailed
			if !errors.As(waitErr, &exit) {
				r.rec.StderrTail = sanitizeTail(waitErr.Error())
			}
		}
		if r.rec.State == StateFailed && r.rec.StderrTail == "" {
			r.rec.StderrTail = sanitizeTail(stderr.Text())
		}
		r.final = r.rec
		for _, err := range []error{
			writeFileAtomic(filepath.Join(dir, summaryFile), []byte(renderSummary(r.rec))),
			writeRecord(filepath.Join(dir, recordFile), r.rec),
		} {
			if err != nil && r.err == nil {
				r.err = err
			}
		}
		close(r.done)
	}()
	return r, nil
}

// Wait blocks until the child has exited and run.yml holds its final state.
func (r *Run) Wait() Record {
	<-r.done
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.final
}

// Err is the first failure to save the run's record or summary, or nil. Check
// it after Wait: when it is set, run.yml may not hold the final state Wait
// returned.
func (r *Run) Err() error {
	<-r.done
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.err
}

// cleanArgv is argv as saved in run.yml: URL credentials and terminal controls
// removed. The child still receives the original argv.
func cleanArgv(argv []string) []string {
	out := make([]string, len(argv))
	for i, a := range argv {
		out[i] = cleanText(a)
	}
	return out
}

// lineWriter appends each complete stdout line verbatim to the events file, one
// Write per line so O_APPEND keeps lines whole.
type lineWriter struct {
	out     *os.File
	buf     []byte
	observe func(line []byte) // called with each line after it is on disk
}

func (w *lineWriter) Write(p []byte) (int, error) {
	w.buf = append(w.buf, p...)
	for {
		i := bytes.IndexByte(w.buf, '\n')
		if i < 0 {
			return len(p), nil
		}
		if _, err := w.out.Write(w.buf[:i+1]); err != nil {
			return 0, err
		}
		w.observe(w.buf[:i])
		w.buf = w.buf[i+1:]
	}
}

// flush writes a final line the child ended without a newline.
func (w *lineWriter) flush() {
	if len(w.buf) > 0 {
		w.out.Write(append(w.buf, '\n'))
		w.observe(w.buf)
		w.buf = nil
	}
}
