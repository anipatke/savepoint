package styles

import "github.com/charmbracelet/lipgloss"

func color(hex, ansi256, ansi string) lipgloss.CompleteColor {
	return lipgloss.CompleteColor{TrueColor: hex, ANSI256: ansi256, ANSI: ansi}
}

var (
	clrOrange      = color(AtariOrange, AtariOrange256, AtariOrange16)
	clrText        = color(PrimaryText, PrimaryText256, PrimaryText16)
	clrBorder      = color(BorderSubtle, BorderSubtle256, BorderSubtle16)
	clrSurface     = color(Surface2, Surface2256, Surface216) // intentionally black
	clrSurfaceDark = color(Surface, Surface256, Surface16)    // intentionally black
	clrGreen       = color(NPPGreen, NPPGreen256, NPPGreen16)
	clrPurple      = color(VibePurple, VibePurple256, VibePurple16)
	clrDim         = color(Dim, Dim256, Dim16)
)

var boxBorder = lipgloss.NormalBorder()

var (
	HeaderIcon = lipgloss.NewStyle().
			Foreground(clrOrange).
			Bold(true)

	HeaderText = lipgloss.NewStyle().
			Foreground(clrText)

	Divider = lipgloss.NewStyle().
		Foreground(clrBorder)

	HeaderFrame = lipgloss.NewStyle().
			Padding(1, 1)

	BoardFrame = lipgloss.NewStyle()

	Column = lipgloss.NewStyle().
		Padding(0, 1)

	ColumnUnfocused = lipgloss.NewStyle().
			BorderStyle(boxBorder).
			BorderForeground(clrBorder).
			Padding(0, 1)

	ColumnFocused = lipgloss.NewStyle().
			BorderStyle(boxBorder).
			BorderForeground(clrOrange).
			Padding(0, 1)

	ColumnTitle = lipgloss.NewStyle().
			Foreground(clrText).
			Bold(true)

	ColumnTitleFocused = lipgloss.NewStyle().
				Foreground(clrOrange).
				Bold(true)

	TaskItem = lipgloss.NewStyle().
			Foreground(clrText)

	TaskItemFocused = lipgloss.NewStyle().
			Foreground(clrOrange)

	StatusBar = lipgloss.NewStyle().
			Foreground(clrText)

	EpicPanel = lipgloss.NewStyle().
			BorderStyle(boxBorder).
			BorderForeground(clrBorder).
			Padding(0, 1)

	EpicItemFocused = lipgloss.NewStyle().
			Foreground(clrPurple)

	EpicTitleFocused = lipgloss.NewStyle().
				Foreground(clrPurple).
				Bold(true)

	EpicPanelFocused = lipgloss.NewStyle().
				BorderStyle(boxBorder).
				BorderForeground(clrPurple).
				Padding(0, 1)

	Card = lipgloss.NewStyle().
		Padding(0, 1)

	CardFocused = lipgloss.NewStyle().
			BorderStyle(boxBorder).
			BorderForeground(clrOrange).
			Padding(0, 1)

	// CardBox and CardBoxFocused are one card frame in two accents. Both carry
	// the same border and padding, so focus changes color alone and a card
	// never occupies different cells focused than unfocused.
	CardBox = lipgloss.NewStyle().
		BorderStyle(boxBorder).
		BorderForeground(clrBorder).
		Padding(0, 1)

	CardBoxFocused = lipgloss.NewStyle().
			BorderStyle(boxBorder).
			BorderForeground(clrOrange).
			Padding(0, 1)

	CardMeta        = lipgloss.NewStyle().Foreground(clrDim)
	ScrollIndicator = lipgloss.NewStyle().
			Foreground(clrDim).
			Faint(true)

	GlyphBuild = lipgloss.NewStyle().Foreground(clrOrange)
	GlyphTest  = lipgloss.NewStyle().Foreground(clrGreen)
	GlyphAudit = lipgloss.NewStyle().Foreground(clrPurple)

	DetailOverlay = lipgloss.NewStyle().
			BorderStyle(boxBorder).
			BorderForeground(clrOrange).
			Padding(0, 1)

	EpicDetailOverlay = lipgloss.NewStyle().
				BorderStyle(boxBorder).
				BorderForeground(clrPurple).
				Padding(0, 1)

	// Footer phase styles
	FooterPhasePlan = lipgloss.NewStyle().
			Foreground(clrPurple).
			Bold(true)

	FooterPhaseBuild = lipgloss.NewStyle().
				Foreground(clrOrange).
				Bold(true)

	FooterPhaseAudit = lipgloss.NewStyle().
				Foreground(clrGreen).
				Bold(true)

	FooterPhaseDefect = lipgloss.NewStyle().
				Foreground(clrOrange).
				Bold(true)

	FooterDivider = lipgloss.NewStyle().
			Foreground(clrBorder)

	FooterHints = lipgloss.NewStyle().
			Foreground(clrDim)

	HeaderRight = lipgloss.NewStyle().
			Foreground(clrDim)

	HeaderRelease = lipgloss.NewStyle().
			Foreground(clrText)

	RootLine = lipgloss.NewStyle()

	// Tag styles for semantic encoding
	TagDone = lipgloss.NewStyle().Foreground(clrGreen)
	TagAI   = lipgloss.NewStyle().Foreground(clrPurple)

	// Badge styles are the four semantic accents a state badge can carry, in
	// the existing palette. They encode a category, never the state itself:
	// each badge also carries a glyph and a label, so every distinction
	// survives with color disabled.
	//
	// Clear: a requirement that is met. Attention: one that needs the reader.
	// Waiting: work that belongs to another record. Neutral: a state with
	// nothing recorded yet.
	// SidebarSelected marks the sidebar row whose record the board is filtered
	// to. It is a different accent from the focused-item style so a selection
	// the cursor has moved off stays visible, and the row carries its own glyph
	// besides, so the distinction survives with color stripped.
	SidebarSelected = lipgloss.NewStyle().
			Foreground(clrPurple).
			Bold(true)

	BadgeClear     = lipgloss.NewStyle().Foreground(clrGreen)
	BadgeAttention = lipgloss.NewStyle().Foreground(clrOrange)
	BadgeWaiting   = lipgloss.NewStyle().Foreground(clrPurple)
	BadgeNeutral   = lipgloss.NewStyle().Foreground(clrDim)
)
