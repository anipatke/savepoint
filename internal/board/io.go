package board

import (
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

func readEpicDetailCmd(epicDir, shortIDStr string) tea.Cmd {
	return func() tea.Msg {
		content := readEpicDetailFile(epicDir, shortIDStr)
		return epicDetailMsg{content: content}
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
