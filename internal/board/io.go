package board

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/opencode/savepoint/internal/data"
)

func writeRouterTaskCmd(root string, task data.Task, reader routerReader) tea.Cmd {
	return func() tea.Msg {
		routerPath := filepath.Join(root, "router.md")
		fi, err := os.Stat(routerPath)
		if err != nil {
			return errorMsg{message: err.Error()}
		}
		content, err := os.ReadFile(routerPath)
		if err != nil {
			return errorMsg{message: err.Error()}
		}
		state, err := reader.ReadState(string(content))
		if err != nil {
			return errorMsg{message: err.Error()}
		}
		state.Release = task.Release
		state.Epic = task.Epic
		state.State = "task-building"
		state.Task = task.ID
		state.NextAction = "Build " + task.ID + "."
		if err := data.WriteRouterState(root, state, fi.ModTime()); err != nil {
			return errorMsg{message: err.Error()}
		}
		message := "Router set to " + task.Release + " " + task.Epic + "/" + shortID(task.ID)
		return routerWriteMsg{message: message, state: state, taskID: task.ID}
	}
}

func writeRouterReleaseEpicCmd(root, selectedEpic, selectedRelease string, reader routerReader) tea.Cmd {
	return func() tea.Msg {
		routerPath := filepath.Join(root, "router.md")
		fi, err := os.Stat(routerPath)
		if err != nil {
			return errorMsg{message: err.Error()}
		}
		content, err := os.ReadFile(routerPath)
		if err != nil {
			return errorMsg{message: err.Error()}
		}
		state, err := reader.ReadState(string(content))
		if err != nil {
			return errorMsg{message: err.Error()}
		}
		state.Epic = shortID(selectedEpic)
		state.Release = selectedRelease
		if err := data.WriteRouterState(root, state, fi.ModTime()); err != nil {
			return errorMsg{message: err.Error()}
		}
		return routerWriteMsg{state: state}
	}
}

func writeTaskStatusCmd(orig, next data.Task, expectedMtime time.Time, prefix string) tea.Cmd {
	return func() tea.Msg {
		if err := data.WriteTaskStatus(next.Path, &next, expectedMtime); err != nil {
			if errors.Is(err, data.ErrMtimeConflict) {
				return retryTaskStatusAfterConflict(orig, next, prefix)
			}
			return errorMsg{message: taskWriteErrorMessage(err)}
		}
		fi, err := os.Stat(next.Path)
		if err != nil {
			return errorMsg{message: err.Error()}
		}
		next.Mtime = fi.ModTime()
		return taskWriteMsg{prefix: prefix, next: next}
	}
}

func writeDefectStatusCmd(next data.Defect, expectedMtime time.Time) tea.Cmd {
	return func() tea.Msg {
		if err := data.WriteDefectStatus(next.Path, &next, expectedMtime); err != nil {
			if errors.Is(err, data.ErrMtimeConflict) {
				return errorMsg{message: "defect changed on disk: refresh before retrying"}
			}
			return errorMsg{message: err.Error()}
		}
		fi, err := os.Stat(next.Path)
		if err != nil {
			return errorMsg{message: err.Error()}
		}
		next.Mtime = fi.ModTime()
		return defectWriteMsg{next: next}
	}
}

// writeEpicStatusCmd persists status to the epic's E##-Detail.md via
// data.UpdateEpicStatus, guarded by expectedMtime so a file changed since it was
// read does not get a partial overwrite. It mirrors writeDefectStatusCmd.
func writeEpicStatusCmd(epicID, path, status string, expectedMtime time.Time) tea.Cmd {
	return func() tea.Msg {
		fi, err := os.Stat(path)
		if err != nil {
			return errorMsg{message: err.Error()}
		}
		if !fi.ModTime().Equal(expectedMtime) {
			return errorMsg{message: "epic changed on disk: refresh before retrying"}
		}
		if err := data.UpdateEpicStatus(path, status); err != nil {
			return errorMsg{message: err.Error()}
		}
		return epicStatusWrittenMsg{epicID: epicID, status: status}
	}
}

func retryTaskStatusAfterConflict(orig, next data.Task, prefix string) tea.Msg {
	current, err := readTaskFromDisk(orig.Path, orig.Release, orig.Epic)
	if err != nil {
		return errorMsg{message: taskWriteErrorMessage(err)}
	}
	debugf("mtime conflict for %s: expected=%s current=%s", orig.ID, orig.Mtime.Format(time.RFC3339Nano), current.Mtime.Format(time.RFC3339Nano))
	if !sameTransitionBase(orig, current) {
		return taskRefreshMsg{
			task:    current,
			message: "task changed on disk: refreshed, retry if still intended",
		}
	}

	next.Mtime = current.Mtime
	if err := data.WriteTaskStatus(next.Path, &next, current.Mtime); err != nil {
		return errorMsg{message: taskWriteErrorMessage(err)}
	}
	fi, err := os.Stat(next.Path)
	if err != nil {
		return errorMsg{message: err.Error()}
	}
	next.Mtime = fi.ModTime()
	debugf("mtime conflict retry succeeded for %s", orig.ID)
	return taskWriteMsg{prefix: prefix, next: next}
}

func readTaskFromDisk(path, release, epic string) (data.Task, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return data.Task{}, err
	}
	task, err := data.NewParser().ParseTaskFile(path, string(raw))
	if err != nil {
		return data.Task{}, err
	}
	fi, err := os.Stat(path)
	if err != nil {
		return data.Task{}, err
	}
	task.Path = path
	task.Mtime = fi.ModTime()
	task.Release = release
	task.Epic = epic
	return *task, nil
}

func sameTransitionBase(orig, current data.Task) bool {
	return orig.ID == current.ID &&
		data.SameTaskLifecycleForTransition(orig, current)
}

func sameTaskRecord(a, b data.Task) bool {
	if a.Path != "" && b.Path != "" {
		return a.Path == b.Path
	}
	if a.ID != b.ID {
		return false
	}
	if a.Release != "" && b.Release != "" && a.Release != b.Release {
		return false
	}
	if a.Epic != "" && b.Epic != "" && a.Epic != b.Epic {
		return false
	}
	return true
}

func readEpicDetailCmd(epicDir, shortIDStr string) tea.Cmd {
	return func() tea.Msg {
		content := readEpicDetailFile(epicDir, shortIDStr)
		return epicDetailMsg{content: content}
	}
}

// loadReleaseDocsCmd reads the Release Docs overlay's supporting documents (the
// selected release's PRD plus the project-wide PRD/Design) from the .savepoint
// root through the data loader, keeping filesystem access out of Update(). A
// loader error surfaces as a status message rather than a panic.
func loadReleaseDocsCmd(root, release string) tea.Cmd {
	return func() tea.Msg {
		docs, err := data.LoadReleaseDocs(root, release)
		if err != nil {
			return errorMsg{message: err.Error()}
		}
		return releaseDocsMsg{docs: docs}
	}
}

func readEpicAuditCmd(epicDir, shortIDStr string) tea.Cmd {
	return func() tea.Msg {
		raw, err := os.ReadFile(filepath.Join(epicDir, shortIDStr+"-Audit.md"))
		if err != nil {
			return auditContentMsg{content: "(no audit available)"}
		}
		return auditContentMsg{content: string(raw)}
	}
}
