package run

import (
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

const (
	// DefaultTimeout bounds a run when the Spec sets none.
	DefaultTimeout = 30 * time.Minute
	// DefaultGrace is how long a stopped agent gets to exit after SIGINT before
	// its whole process group is killed. Windows has no interrupt and skips it.
	DefaultGrace = 10 * time.Second
	// LostAfter is how stale updated_at must be before a dead pid means the
	// runner is lost. Several heartbeats, so a slow disk is not read as death.
	LostAfter = 30 * time.Second

	cancelFile = "cancel"
	cancelPoll = 200 * time.Millisecond
)

// stopper ends the child when a cancel file appears in the run directory or the
// timeout passes. It records which of the two came first for the final state.
type stopper struct {
	mu    sync.Mutex
	state string        // StateCancelled or StateTimedOut once a stop began
	done  chan struct{} // closed when the stopper has nothing left to do
}

// reason is the state the stop calls for, or "" when the child ended unaided.
func (s *stopper) reason() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.state
}

// startStopper watches runDir for a cancel file (checked every cancelPoll) and
// the clock for timeout. On either it interrupts the agent where the platform
// can, waits up to grace for it to go, then kills the whole process tree. It
// returns when the child has exited (procDone closed) and any kill is finished,
// so a caller that waits on done knows nothing is left running.
func startStopper(runDir string, cmd *exec.Cmd, timeout, grace time.Duration, procDone <-chan struct{}) *stopper {
	s := &stopper{done: make(chan struct{})}
	go func() {
		defer close(s.done)
		poll := time.NewTicker(cancelPoll)
		defer poll.Stop()
		limit := time.NewTimer(timeout)
		defer limit.Stop()
		path := filepath.Join(runDir, cancelFile)
		for {
			var state string
			select {
			case <-procDone:
				return
			case <-limit.C:
				state = StateTimedOut
			case <-poll.C:
				if _, err := os.Stat(path); err != nil {
					continue
				}
				state = StateCancelled
			}
			select {
			case <-procDone: // it ended on its own first: not a stop
				return
			default:
			}
			s.mu.Lock()
			s.state = state
			s.mu.Unlock()
			if interruptProcess(cmd) {
				select {
				case <-procDone:
				case <-time.After(grace):
				}
			}
			killTree(cmd)
			return
		}
	}()
	return s
}

// RunnerLost reports whether rec describes a run whose runner died without
// writing a final state: it still says running, its heartbeat is older than
// LostAfter, and the runner process is gone (whether or not the agent lives). It is an inference for readers to show;
// it never writes to run.yml. A record whose time cannot be read, or that has no
// pid, is not called lost.
func RunnerLost(rec Record, now time.Time) bool {
	pid := rec.RunnerPID
	if pid <= 0 {
		pid = rec.PID // records written before runner_pid existed
	}
	if rec.State != StateRunning || pid <= 0 {
		return false
	}
	updated, err := time.Parse(time.RFC3339, rec.UpdatedAt)
	if err != nil || now.Sub(updated) <= LostAfter {
		return false
	}
	return !processAlive(pid)
}
