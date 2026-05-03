package styles

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestPaletteConstants_present(t *testing.T) {
	// Truecolor tier
	if Background == "" {
		t.Error("Background constant is empty")
	}
	if Surface == "" {
		t.Error("Surface constant is empty")
	}
	if Surface2 == "" {
		t.Error("Surface2 constant is empty")
	}
	if Border == "" {
		t.Error("Border constant is empty")
	}
	if BorderSubtle == "" {
		t.Error("BorderSubtle constant is empty")
	}
	if PrimaryText == "" {
		t.Error("PrimaryText constant is empty")
	}
	if AtariOrange == "" {
		t.Error("AtariOrange constant is empty")
	}
	if NPPGreen == "" {
		t.Error("NPPGreen constant is empty")
	}
	if VibePurple == "" {
		t.Error("VibePurple constant is empty")
	}
	if Dim == "" {
		t.Error("Dim constant is empty")
	}
}

func TestPaletteConstants_256tier(t *testing.T) {
	if Background256 == "" {
		t.Error("Background256 constant is empty")
	}
	if Surface256 == "" {
		t.Error("Surface256 constant is empty")
	}
	if Surface2256 == "" {
		t.Error("Surface2256 constant is empty")
	}
	if Border256 == "" {
		t.Error("Border256 constant is empty")
	}
	if BorderSubtle256 == "" {
		t.Error("BorderSubtle256 constant is empty")
	}
	if PrimaryText256 == "" {
		t.Error("PrimaryText256 constant is empty")
	}
	if AtariOrange256 == "" {
		t.Error("AtariOrange256 constant is empty")
	}
	if NPPGreen256 == "" {
		t.Error("NPPGreen256 constant is empty")
	}
	if VibePurple256 == "" {
		t.Error("VibePurple256 constant is empty")
	}
	if Dim256 == "" {
		t.Error("Dim256 constant is empty")
	}
}

func TestPaletteConstants_16tier(t *testing.T) {
	if Background16 == "" {
		t.Error("Background16 constant is empty")
	}
	if Surface16 == "" {
		t.Error("Surface16 constant is empty")
	}
	if Surface216 == "" {
		t.Error("Surface216 constant is empty")
	}
	if Border16 == "" {
		t.Error("Border16 constant is empty")
	}
	if BorderSubtle16 == "" {
		t.Error("BorderSubtle16 constant is empty")
	}
	if PrimaryText16 == "" {
		t.Error("PrimaryText16 constant is empty")
	}
	if AtariOrange16 == "" {
		t.Error("AtariOrange16 constant is empty")
	}
	if NPPGreen16 == "" {
		t.Error("NPPGreen16 constant is empty")
	}
	if VibePurple16 == "" {
		t.Error("VibePurple16 constant is empty")
	}
}

func TestColor(t *testing.T) {
	c := color("#FF0000", "196", "9")
	expected := lipgloss.CompleteColor{TrueColor: "#FF0000", ANSI256: "196", ANSI: "9"}
	if c != expected {
		t.Errorf("color() = %+v, want %+v", c, expected)
	}
}

func TestColor_usesPaletteConstants(t *testing.T) {
	c := color(AtariOrange, AtariOrange256, AtariOrange16)
	if c.TrueColor != AtariOrange {
		t.Errorf("color TrueColor = %q, want %q", c.TrueColor, AtariOrange)
	}
	if c.ANSI256 != AtariOrange256 {
		t.Errorf("color ANSI256 = %q, want %q", c.ANSI256, AtariOrange256)
	}
	if c.ANSI != AtariOrange16 {
		t.Errorf("color ANSI = %q, want %q", c.ANSI, AtariOrange16)
	}
}

func TestColor_returnsCompleteColor(t *testing.T) {
	var cc lipgloss.CompleteColor
	cc = color(Background, Background256, Background16)
	if cc.TrueColor != Background {
		t.Errorf("expected TrueColor %q, got %q", Background, cc.TrueColor)
	}
}
