package run

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Run states. They are owned by this package and never feed a lifecycle gate.
const (
	StateRunning  = "running"
	StateFinished = "finished"
	StateFailed   = "failed"
	// StateCancelled and StateTimedOut are written by the runner when it stops
	// the agent. StateRunnerLost is never written: a reader infers it with
	// RunnerLost from a run that stayed running after its runner died.
	StateCancelled  = "cancelled"
	StateTimedOut   = "timed_out"
	StateRunnerLost = "runner_lost"
)

// Record is the run.yml document: what the runner observed about one process.
type Record struct {
	Provider  string   `yaml:"provider"`
	Command   []string `yaml:"command"`
	Cwd       string   `yaml:"cwd"`
	Task      string   `yaml:"task"`
	NextLine  string   `yaml:"next_line"`
	StartedAt string   `yaml:"started_at"`
	UpdatedAt string   `yaml:"updated_at"`
	EndedAt   string   `yaml:"ended_at,omitempty"`
	PID       int      `yaml:"pid,omitempty"` // the agent process
	// RunnerPID is the process that keeps this record current. A reader tests
	// it, not PID, to tell whether the runner is lost.
	RunnerPID  int    `yaml:"runner_pid,omitempty"`
	State      string `yaml:"state"`
	ExitCode   *int   `yaml:"exit_code,omitempty"`
	StderrTail string `yaml:"stderr_tail,omitempty"`

	// What the runner read from the provider's event stream. Each value was
	// seen on a line in events.jsonl; none is a verdict on the work.
	SessionID   string       `yaml:"session_id,omitempty"`
	Model       string       `yaml:"model,omitempty"`
	FilesRead   []string     `yaml:"files_read,omitempty"`
	FilesEdited []string     `yaml:"files_edited,omitempty"`
	Commands    []CommandRun `yaml:"commands,omitempty"`
	Denials     []Denial     `yaml:"denials,omitempty"`
	// Usage is the provider's usage object verbatim; TotalCostUSD is the
	// provider's own client-side estimate, not a bill.
	Usage        map[string]any `yaml:"usage,omitempty"`
	TotalCostUSD *float64       `yaml:"total_cost_usd,omitempty"`
	// ListsTruncated is set when files_read, files_edited or commands hit
	// maxListed and later entries were left out of this record (they remain in
	// events.jsonl).
	ListsTruncated bool `yaml:"lists_truncated,omitempty"`
	// Unparsed counts stdout lines the adapter could not interpret. They stay
	// verbatim in events.jsonl.
	Unparsed int `yaml:"unparsed"`
	// AgentReported is the agent's own final message. It is the agent's claim,
	// not something the runner checked.
	AgentReported string `yaml:"agent_reported,omitempty"`
}

// CommandRun is one command the agent ran. ExitCode is nil until the provider
// reports how the command ended.
type CommandRun struct {
	Command  string `yaml:"command"`
	ExitCode *int   `yaml:"exit_code,omitempty"`
}

// Denial is one action the provider refused to run for lack of permission.
type Denial struct {
	Tool      string `yaml:"tool"`
	Target    string `yaml:"target,omitempty"`
	Reason    string `yaml:"reason,omitempty"`
	ToolUseID string `yaml:"tool_use_id,omitempty"`
}

// writeRecord replaces path with the encoded record, so a reader never sees a
// partial run.yml.
func writeRecord(path string, rec Record) error {
	content, err := yaml.Marshal(rec)
	if err != nil {
		return fmt.Errorf("encode %s: %w", path, err)
	}
	return writeFileAtomic(path, content)
}

// writeFileAtomic replaces path with content through a same-directory temporary
// file, so a reader never sees a partial file (FS-04).
func writeFileAtomic(path string, content []byte) (err error) {
	temp, err := os.CreateTemp(filepath.Dir(path), ".run-*")
	if err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	tempPath := temp.Name()
	defer func() {
		if err != nil {
			os.Remove(tempPath)
		}
	}()
	if _, err = temp.Write(content); err != nil {
		temp.Close()
		return fmt.Errorf("write %s: %w", path, err)
	}
	if err = temp.Sync(); err != nil {
		temp.Close()
		return fmt.Errorf("write %s: %w", path, err)
	}
	if err = temp.Close(); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	if err = os.Chmod(tempPath, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	if err = os.Rename(tempPath, path); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

// summaryNotice sits under the heading of summary.md so no reader mistakes the
// agent's message for something Savepoint checked.
const summaryNotice = "Agent-reported: the text below is the agent's own final message. Savepoint copied it and did not check it. What the runner observed is in run.yml."

// renderSummary is summary.md: the agent's final message, copied as it was
// reported (terminal control sequences and URL credentials removed).
func renderSummary(rec Record) string {
	body := rec.AgentReported
	if body == "" {
		body = "The agent reported no final message."
	}
	return "# Agent-reported summary\n\n" + summaryNotice + "\n\n" + body + "\n"
}
