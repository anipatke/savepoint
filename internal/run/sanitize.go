package run

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// maxStderrBytes bounds the stderr kept in memory; maxTailLen bounds the single
// line stored in run.yml.
const (
	maxStderrBytes = 64 << 10
	maxTailLen     = 512
)

// tailBuffer keeps the last limit bytes and drops the earlier ones, so the error
// an agent prints last survives however much came first. It always accepts
// writes. Copied from the health runner, which this package must not import.
type tailBuffer struct {
	buf     []byte
	limit   int
	dropped bool
}

func (b *tailBuffer) Write(p []byte) (int, error) {
	b.buf = append(b.buf, p...)
	if len(b.buf) > 2*b.limit {
		b.dropped = true
		b.buf = append(b.buf[:0], b.buf[len(b.buf)-b.limit:]...)
	}
	return len(p), nil
}

// cutUserinfo matches the start of text that may be the end of a URL whose
// "scheme://" was cut away: some "user:password@".
var cutUserinfo = regexp.MustCompile(`^[^\s/?#@]*[:/]*[^\s/?#@]*@`)

// Text is the kept tail. When earlier bytes were dropped it starts mid-word, so
// a leading "user:password@" fragment is removed.
func (b *tailBuffer) Text() string {
	s := string(b.buf)
	if len(b.buf) > b.limit {
		s = s[len(s)-b.limit:]
	}
	if b.dropped || len(b.buf) > b.limit {
		s = cutUserinfo.ReplaceAllString(s, "")
	}
	return s
}

// urlUserinfo matches the credentials of a URL, "user:password@" after "://".
var urlUserinfo = regexp.MustCompile(`://[^\s/?#@]*@`)

// terminalSequence matches ANSI escape sequences in both the 7-bit (ESC) and
// 8-bit (C1) forms: CSI, OSC up to BEL or ST, and the two-byte escapes.
var terminalSequence = regexp.MustCompile(`\x1b\[[0-?]*[ -/]*[@-~]|\x{9b}[0-?]*[ -/]*[@-~]|\x1b\][^\x07\x1b]*(?:\x07|\x1b\\)?|\x1b[@-Z\\-_]`)

// cleanText removes terminal control sequences and control characters other
// than newline and tab, and strips URL credentials, so text copied from agent
// output into the run record is safe to show (ARCH-05). It keeps line structure.
func cleanText(s string) string {
	s = strings.ToValidUTF8(s, "")
	s = terminalSequence.ReplaceAllString(s, "")
	s = strings.Map(func(r rune) rune {
		if r != '\n' && r != '\t' && unicode.IsControl(r) {
			return -1
		}
		return r
	}, s)
	return urlUserinfo.ReplaceAllString(s, "://")
}

// oneLine is cleanText on a single line, cut to max bytes at a rune boundary
// with the start kept: a path or command is identified by how it begins.
func oneLine(s string, max int) string {
	s = strings.Join(strings.Fields(cleanText(s)), " ")
	if len(s) <= max {
		return s
	}
	cut := max
	for cut > 0 && !utf8.RuneStart(s[cut]) {
		cut--
	}
	return s[:cut] + "…"
}

// sanitizeTail reduces stderr to one printable line that keeps its end, with
// URL credentials removed, since the text is saved in the run record (ARCH-05).
func sanitizeTail(s string) string {
	s = strings.ToValidUTF8(s, "")
	s = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return ' '
		}
		return r
	}, s)
	s = strings.Join(strings.Fields(s), " ")
	s = urlUserinfo.ReplaceAllString(s, "://")
	if len(s) > maxTailLen {
		cut := len(s) - maxTailLen
		for cut < len(s) && !utf8.RuneStart(s[cut]) {
			cut++
		}
		s = "…" + s[cut:]
	}
	return s
}
