package board

import (
	"testing"

	"github.com/muesli/termenv"
)

func TestDetectProfileNOCOLOR(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	t.Setenv("COLORTERM", "")
	if got := detectProfile(); got != termenv.Ascii {
		t.Fatalf("want Ascii, got %v", got)
	}
}

func TestDetectProfileTruecolor(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	t.Setenv("COLORTERM", "truecolor")
	if got := detectProfile(); got != termenv.TrueColor {
		t.Fatalf("want TrueColor, got %v", got)
	}
}

func TestDetectProfile24bit(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	t.Setenv("COLORTERM", "24bit")
	if got := detectProfile(); got != termenv.TrueColor {
		t.Fatalf("want TrueColor, got %v", got)
	}
}
