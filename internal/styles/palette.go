package styles

// Truecolor hex constants (Atari-Noir palette).
// Background, Surface, and Surface2 intentionally share one black value so the
// terminal stays visually flat; hierarchy comes from spacing, dividers, and accents.
const (
	Background   = "#000000"
	Surface      = "#000000"
	Surface2     = "#000000"
	Border       = "#1A1A1A"
	BorderSubtle = "#222222"
	PrimaryText  = "#F0E6DA"
	AtariOrange  = "#FC6323"
	NPPGreen     = "#A4C639"
	VibePurple   = "#B1A1DF"
)

// 256-color (ANSI256) fallbacks — nearest terminal approximations
const (
	Background256   = "16"
	Surface256      = "16"
	Surface2256     = "16"
	Border256       = "234"
	BorderSubtle256 = "235"
	PrimaryText256  = "230"
	AtariOrange256  = "208"
	NPPGreen256     = "148"
	VibePurple256   = "147"
)

// 16-color (basic ANSI) fallbacks
const (
	Background16   = "0"  // black
	Surface16      = "0"  // black
	Surface216     = "0"  // black
	Border16       = "0"  // black
	BorderSubtle16 = "8"  // dark gray
	PrimaryText16  = "15" // bright white
	AtariOrange16  = "9"  // bright red (closest to orange)
	NPPGreen16     = "2"  // green
	VibePurple16   = "5"  // magenta
	Dim16          = "8"  // dark gray
)

// Dim: muted foreground for metadata
const (
	Dim    = "#777777"
	Dim256 = "243"
)
