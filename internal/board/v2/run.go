package v2

import (
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// Options is the V2 board's entry contract. Root is the project's .savepoint
// directory, already resolved by the dispatch; ObjectiveFilter is --objective;
// Stdout and TTY are passed in rather than read from the process so the board
// runs against a temporary project without depending on process state
// (ARCH-03).
type Options struct {
	Root            string
	ObjectiveFilter string
	Stdout          io.Writer
	TTY             bool
}

// Run opens the V2 board: the Bubble Tea program on a terminal, deterministic
// text otherwise. Both paths load through the same loadProject, so the two
// surfaces cannot disagree about what the project says.
func Run(opts Options) error {
	if opts.TTY {
		return runTUI(opts)
	}
	return runPlain(opts)
}

// runTUI starts the program, whose Init runs the load command. A load that
// produced a diagnostic is a screen the user reads and dismisses, so it is not
// an error here; an invocation the board cannot honor is, and it is returned
// from the final model so the command exits nonzero.
func runTUI(opts Options) error {
	watcher, err := newV2Watcher(opts.Root)
	if err != nil {
		return fmt.Errorf("watch V2 project: %w", err)
	}
	defer watcher.Close()

	model := NewModel(opts)
	model.Watcher = watcher
	program := tea.NewProgram(model, tea.WithAltScreen())
	final, err := program.Run()
	if err != nil {
		return err
	}
	if model, ok := final.(Model); ok && model.FatalErr != nil {
		return model.FatalErr
	}
	return nil
}

// runPlain renders the board for a writer that is not a terminal. A load
// diagnostic is returned as an error and nothing is written, so a piped board
// reports the failure on stderr, exits nonzero, and never emits a partial board
// that could be read as the project's real state.
func runPlain(opts Options) error {
	loaded := loadProject(opts.Root)
	if loaded.Failed() {
		return fmt.Errorf("invalid project data: %s", loaded.Diagnostic)
	}
	if err := objectiveFilterError(loaded.State.Index, opts.ObjectiveFilter); err != nil {
		return err
	}

	out := opts.Stdout
	if out == nil {
		out = os.Stdout
	}
	_, err := fmt.Fprint(out, renderPlain(loaded.State, selectedObjective(loaded.State, opts.ObjectiveFilter)))
	return err
}
