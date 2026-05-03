package init

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
)

type ClipboardStatus int

const (
	ClipboardCopied  ClipboardStatus = iota
	ClipboardSkipped
	ClipboardFailed
)

func (s ClipboardStatus) String() string {
	switch s {
	case ClipboardCopied:
		return "copied"
	case ClipboardSkipped:
		return "skipped"
	case ClipboardFailed:
		return "failed"
	default:
		return "unknown"
	}
}

type ClipboardResult struct {
	Status  ClipboardStatus
	Tool    string
	Message string
}

func CopyToClipboard(text string) ClipboardResult {
	switch runtime.GOOS {
	case "windows":
		return copyWindows(text)
	case "darwin":
		return copyDarwin(text)
	case "linux":
		return copyLinux(text)
	default:
		return ClipboardResult{
			Status:  ClipboardSkipped,
			Tool:    "",
			Message: fmt.Sprintf("unsupported platform: %s", runtime.GOOS),
		}
	}
}

func copyWindows(text string) ClipboardResult {
	if _, err := exec.LookPath("clip"); err != nil {
		return ClipboardResult{
			Status:  ClipboardSkipped,
			Tool:    "clip.exe",
			Message: "clip.exe not found in PATH",
		}
	}

	cmd := exec.Command("clip")
	cmd.Stdin = bytes.NewReader([]byte(text))
	if err := cmd.Run(); err != nil {
		return ClipboardResult{
			Status:  ClipboardFailed,
			Tool:    "clip.exe",
			Message: fmt.Sprintf("clip.exe failed: %v", err),
		}
	}

	return ClipboardResult{
		Status:  ClipboardCopied,
		Tool:    "clip.exe",
		Message: "",
	}
}

func copyDarwin(text string) ClipboardResult {
	if _, err := exec.LookPath("pbcopy"); err != nil {
		return ClipboardResult{
			Status:  ClipboardSkipped,
			Tool:    "pbcopy",
			Message: "pbcopy not found in PATH",
		}
	}

	cmd := exec.Command("pbcopy")
	cmd.Stdin = bytes.NewReader([]byte(text))
	if err := cmd.Run(); err != nil {
		return ClipboardResult{
			Status:  ClipboardFailed,
			Tool:    "pbcopy",
			Message: fmt.Sprintf("pbcopy failed: %v", err),
		}
	}

	return ClipboardResult{
		Status:  ClipboardCopied,
		Tool:    "pbcopy",
		Message: "",
	}
}

func copyLinux(text string) ClipboardResult {
	if path, err := exec.LookPath("xclip"); err == nil {
		cmd := exec.Command(path, "-selection", "clipboard")
		cmd.Stdin = bytes.NewReader([]byte(text))
		if err := cmd.Run(); err != nil {
			return ClipboardResult{
				Status:  ClipboardFailed,
				Tool:    "xclip",
				Message: fmt.Sprintf("xclip failed: %v", err),
			}
		}
		return ClipboardResult{
			Status:  ClipboardCopied,
			Tool:    "xclip",
			Message: "",
		}
	}

	if path, err := exec.LookPath("xsel"); err == nil {
		cmd := exec.Command(path, "--input", "--clipboard")
		cmd.Stdin = bytes.NewReader([]byte(text))
		if err := cmd.Run(); err != nil {
			return ClipboardResult{
				Status:  ClipboardFailed,
				Tool:    "xsel",
				Message: fmt.Sprintf("xsel failed: %v", err),
			}
		}
		return ClipboardResult{
			Status:  ClipboardCopied,
			Tool:    "xsel",
			Message: "",
		}
	}

	return ClipboardResult{
		Status:  ClipboardSkipped,
		Tool:    "",
		Message: "no clipboard tool found (tried xclip, xsel)",
	}
}
