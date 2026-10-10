package run

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/resume"
)

// MinClaudeVersion is the oldest Claude Code with the --permission-prompts flag
// the default command relies on.
const MinClaudeVersion = "2.1.259"

const versionTimeout = 15 * time.Second

// Request is one invocation of the run command.
type Request struct {
	Dir  string // project directory; the nearest .savepoint must be in it
	Task string // T-###
	// Version reports the installed agent's version text for the command in
	// argv[0]. When nil the command is asked with --version.
	Version func(command string) (string, error)
	// Heartbeat, Grace and Now are passed to Start; zero values use its defaults.
	Heartbeat time.Duration
	Grace     time.Duration
	Now       func() time.Time
}

// Exit codes returned by Launch: a refusal the owner can fix by changing the
// command line or the setting is 2; a run or project problem is 1.
const (
	ExitOK      = 0
	ExitFailed  = 1
	ExitRefused = 2
)

// startableVerbs are the Next verbs a run may start from.
var startableVerbs = map[string]bool{"Start": true, "Build": true, "Test": true}

// Launch runs one agent for req.Task in the foreground and writes its record.
// It checks, in order and before any write, that agent_launcher.enabled is
// true, that the Task exists and its Next verb is Start, Build, or Test, and
// that the agent is new enough. It stops the agent when ctx ends, the way
// Ctrl-C in the terminal should. It never writes a Task, Objective, Check, or
// router record.
func Launch(ctx context.Context, req Request, stdout io.Writer) (int, error) {
	root, err := data.ResolveTarget(req.Dir)
	if err != nil {
		return ExitFailed, fmt.Errorf("run: %w", err)
	}
	if err := data.CheckRuntimeSchema(root); err != nil {
		return ExitFailed, fmt.Errorf("run: %w", err)
	}
	savepointRoot := filepath.Join(root, ".savepoint")
	cfg, err := data.NewConfigReader().Read(filepath.Join(savepointRoot, "config.yml"))
	if err != nil {
		return ExitFailed, fmt.Errorf("run: %w", err)
	}
	if !cfg.AgentLauncher.Enabled {
		return ExitRefused, errors.New("run: agent_launcher is off; set agent_launcher.enabled: true in .savepoint/config.yml to use this command")
	}

	index, err := data.LoadV2Index(savepointRoot)
	if err != nil {
		return ExitFailed, fmt.Errorf("run: loading project: %w", err)
	}
	nextLine, err := nextLineFor(index, req.Task)
	if err != nil {
		return ExitFailed, err
	}
	argv, err := BuildArgv(cfg.AgentLauncher.Builder, cfg.QualityGates, nextLine, root)
	if err != nil {
		return ExitFailed, fmt.Errorf("run: %w", err)
	}
	if err := checkAgentVersion(argv[0], req.Version); err != nil {
		return ExitFailed, err
	}

	r, err := Start(Spec{
		Dir: root, Task: req.Task, Provider: ProviderClaude, Argv: argv, NextLine: nextLine,
		Heartbeat: req.Heartbeat, Timeout: cfg.AgentLauncher.RunTimeout(), Grace: req.Grace, Now: req.Now,
	})
	if err != nil {
		return ExitFailed, err
	}
	fmt.Fprintf(stdout, "run: started %s\nrun: record in %s\n", req.Task, r.Dir)

	stopWatch := make(chan struct{})
	watched := make(chan struct{})
	go func() {
		defer close(watched)
		select {
		case <-ctx.Done():
			os.WriteFile(filepath.Join(r.Dir, cancelFile), nil, 0o644)
		case <-stopWatch:
		}
	}()
	final := r.Wait()
	close(stopWatch)
	<-watched

	fmt.Fprintf(stdout, "run: %s\n", final.State)
	if err := r.Err(); err != nil {
		return ExitFailed, fmt.Errorf("run: the record was not saved completely: %w", err)
	}
	if final.State != StateFinished || final.ExitCode == nil || *final.ExitCode != 0 {
		return ExitFailed, nil
	}
	return ExitOK, nil
}

// nextLineFor is the Next line for task as if the router selected it, or an
// error naming why the Task cannot be run.
func nextLineFor(index *data.V2Index, id string) (string, error) {
	task := index.Tasks[id]
	if task == nil {
		return "", fmt.Errorf("run: task %s does not exist", id)
	}
	if task.Status == data.ColumnDone {
		return "", fmt.Errorf("run: task %s is done; there is nothing to start", id)
	}
	objective := index.Objectives[task.Objective]
	if objective == nil {
		return "", fmt.Errorf("run: task %s names objective %s, which does not exist", id, task.Objective)
	}
	next := data.ResolveNext(data.NextInput{Index: index, Router: &data.RouterStateV2{
		State: data.RouterPhaseTask, Release: string(objective.Release), Objective: objective.ID, Task: id,
	}})
	if next.Task == nil || next.Task.ID != id {
		return "", fmt.Errorf("run: task %s cannot be selected: %s", id, resume.NextLine(next))
	}
	if verb := resume.NextVerb(next); !startableVerbs[verb] {
		return "", fmt.Errorf("run: task %s is not ready to start; its next step is %q (%s)", id, verb, resume.NextLine(next))
	}
	return resume.NextLine(next), nil
}

var versionNumber = regexp.MustCompile(`(\d+)\.(\d+)\.(\d+)`)

// checkAgentVersion refuses an agent older than MinClaudeVersion, or one whose
// version cannot be read, before anything is written (CFG-01).
func checkAgentVersion(command string, version func(string) (string, error)) error {
	if version == nil {
		version = askVersion
	}
	text, err := version(command)
	if err != nil {
		return fmt.Errorf("run: cannot read the version of %s: %w; Claude Code %s or newer is required", command, err, MinClaudeVersion)
	}
	got, ok := parseVersion(text)
	if !ok {
		return fmt.Errorf("run: cannot read a version from %q; Claude Code %s or newer is required", text, MinClaudeVersion)
	}
	min, _ := parseVersion(MinClaudeVersion)
	if got.less(min) {
		return fmt.Errorf("run: %s is version %d.%d.%d; Claude Code %s or newer is required", command, got[0], got[1], got[2], MinClaudeVersion)
	}
	return nil
}

type semver [3]int

func (a semver) less(b semver) bool {
	for i := range a {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return false
}

func parseVersion(text string) (semver, bool) {
	m := versionNumber.FindStringSubmatch(text)
	if m == nil {
		return semver{}, false
	}
	var v semver
	for i := range v {
		v[i], _ = strconv.Atoi(m[i+1])
	}
	return v, true
}

func askVersion(command string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), versionTimeout)
	defer cancel()
	out, err := exec.CommandContext(ctx, command, "--version").Output()
	return string(out), err
}
