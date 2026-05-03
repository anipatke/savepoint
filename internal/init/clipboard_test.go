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

func TestCopyToClipboard_resultHasNonEmptyStatusString(t *testing.T) {
	result := CopyToClipboard("test")
	if result.Status.String() == "" {
		t.Fatal("Status.String() returned empty")
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

func TestClipboardStatus_string(t *testing.T) {
	tests := []struct {
		status ClipboardStatus
		want   string
	}{
		{ClipboardCopied, "copied"},
		{ClipboardSkipped, "skipped"},
		{ClipboardFailed, "failed"},
	}
	for _, tt := range tests {
		if got := tt.status.String(); got != tt.want {
			t.Errorf("ClipboardStatus(%d).String() = %q, want %q", tt.status, got, tt.want)
		}
	}
}

func TestCopyToClipboard_multipleCalls(t *testing.T) {
	for i := 0; i < 3; i++ {
		result := CopyToClipboard("test content")
		if result.Status != ClipboardCopied &&
			result.Status != ClipboardSkipped {
			t.Fatalf("call %d: unexpected status %v", i, result.Status)
		}
	}
}
