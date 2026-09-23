package board

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	xterm "github.com/charmbracelet/x/term"
	boardv2 "github.com/opencode/savepoint/internal/board/v2"
	"github.com/opencode/savepoint/internal/migrate"
)

// Filters is the board's live filter surface. Objective is the only filter
// the V2 board accepts; Release and Epic are the V1 flags cmd still parses,
// carried here only so runWithFilters can refuse them by name.
type Filters struct {
	Release   string
	Epic      string
	Objective string
}

func Run() error {
	return RunWithFilters(Filters{})
}

func RunWithFilters(filters Filters) error {
	return runWithFilters(".", filters, os.Stdout, xterm.IsTerminal(os.Stdout.Fd()))
}

// runWithFilters is the board's live dispatch point. It resolves the project
// root once, runs the read-only cutover preflight before any V2 rendering or
// watcher starts, and hands only a valid V2 project to the V2 board. Legacy
// projects and pending/invalid projects receive a named refusal instead of a
// fallback into V1 discovery.
//
// start, stdout, and isTTY are parameters rather than process state so the
// dispatch is exercised against a temporary project directory without
// depending on the working directory (ARCH-03).
func runWithFilters(start string, filters Filters, stdout io.Writer, isTTY bool) error {
	debugf("board dispatch: finding savepoint root from %q", start)
	projectRoot, err := migrate.FindProjectRoot(start)
	if err != nil {
		return err
	}
	root := filepath.Join(projectRoot, ".savepoint")
	debugf("board dispatch: root = %q", root)

	if filters.Release != "" || filters.Epic != "" {
		return fmt.Errorf("--release and --epic are schema_version 1 filters and are unavailable in the V2-only runtime; use --objective")
	}

	preflight := migrate.PreflightCutover(projectRoot, migrate.CutoverPreflightOptions{})
	if diagnostic := preflight.RuntimeDiagnostic(); diagnostic != "" {
		return fmt.Errorf("board: %s", diagnostic)
	}
	return runV2Board(root, filters, stdout, isTTY)
}

// runV2Board hands the resolved root to the V2 board, which owns its own load,
// its own diagnostic screen, and its own non-TTY output.
func runV2Board(root string, filters Filters, stdout io.Writer, isTTY bool) error {
	return boardv2.Run(boardv2.Options{
		Root:            root,
		ObjectiveFilter: filters.Objective,
		Stdout:          stdout,
		TTY:             isTTY,
	})
}
