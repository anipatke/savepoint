package init

import (
	"runtime"
	"testing"
)

func TestCopyToClipboard_returnsWithoutPanic(t *testing.T) {
	result := CopyToClipboard("test content")
	if result.Status != ClipboardCopied &&
		result.Status != ClipboardSkipped &&
		result.Status != ClipboardFailed {
		t.Fatalf("CopyToClipboard() returned unexpected status %v", result.Status)
	}
}

func TestCopyToClipboard_skippedOnUnsupportedPlatform(t *testing.T) {
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		return
	}
	result := CopyToClipboard("test")
	if result.Status != ClipboardSkipped {
		t.Fatalf("on unsupported platform, expected ClipboardSkipped, got %v", result.Status)
	}
}

func TestCopyToClipboard_emptyString(t *testing.T) {
	result := CopyToClipboard("")
	if result.Status != ClipboardCopied &&
		result.Status != ClipboardSkipped &&
		result.Status != ClipboardFailed {
		t.Fatalf("CopyToClipboard('') returned unexpected status %v", result.Status)
	}
}

func TestCopyToClipboard_toolName(t *testing.T) {
	result := CopyToClipboard("test")
	if result.Status == ClipboardCopied && result.Tool == "" {
		t.Fatal("ClipboardCopied but Tool is empty")
	}
}

func TestCopyToClipboard_multipleCalls(t *testing.T) {
	for i := 0; i < 3; i++ {
		result := CopyToClipboard("test content")
		if result.Status != ClipboardCopied &&
			result.Status != ClipboardSkipped &&
			result.Status != ClipboardFailed {
			t.Fatalf("call %d: unexpected status %v", i, result.Status)
		}
	}
}
