package board

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestSetDebugToggle(t *testing.T) {
	t.Cleanup(func() { SetDebug(false) })

	if DebugEnabled() {
		t.Fatal("debug should be off by default")
	}

	SetDebug(true)
	if !DebugEnabled() {
		t.Fatal("debug should be on after SetDebug(true)")
	}

	SetDebug(false)
	if DebugEnabled() {
		t.Fatal("debug should be off after SetDebug(false)")
	}
}

func TestDebugfWritesToStderr(t *testing.T) {
	t.Cleanup(func() { SetDebug(false) })

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	orig := os.Stderr
	os.Stderr = w

	SetDebug(true)
	debugf("test message %d", 42)

	os.Stderr = orig
	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)
	r.Close()

	out := buf.String()
	if !strings.Contains(out, "test message 42") {
		t.Fatalf("expected debug output, got: %q", out)
	}
	if !strings.Contains(out, "[savepoint debug]") {
		t.Fatalf("expected debug prefix, got: %q", out)
	}
}

func TestDebugfSilentWhenDisabled(t *testing.T) {
	t.Cleanup(func() { SetDebug(false) })

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	orig := os.Stderr
	os.Stderr = w

	SetDebug(false)
	debugf("should not appear")

	os.Stderr = orig
	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)
	r.Close()

	if buf.Len() != 0 {
		t.Fatalf("expected no debug output when disabled, got: %q", buf.String())
	}
}

func TestDebugfFormat(t *testing.T) {
	t.Cleanup(func() { SetDebug(false) })

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	orig := os.Stderr
	os.Stderr = w

	SetDebug(true)
	debugf("key=%q value=%d", "hello", 7)

	os.Stderr = orig
	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)
	r.Close()

	want := fmt.Sprintf("[savepoint debug] key=%q value=%d\n", "hello", 7)
	if buf.String() != want {
		t.Fatalf("expected %q, got %q", want, buf.String())
	}
}
