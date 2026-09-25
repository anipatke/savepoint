package styles

import "github.com/charmbracelet/lipgloss"

func color(hex, ansi256, ansi string) lipgloss.CompleteColor {
	return lipgloss.CompleteColor{TrueColor: hex, ANSI256: ansi256, ANSI: ansi}
}

var (
	clrOrange         = color(AtariOrange, AtariOrange256, AtariOrange16)
	clrText           = color(PrimaryText, PrimaryText256, PrimaryText16)
	clrBorder         = color(BorderSubtle, BorderSubtle256, BorderSubtle16)
	clrBorderPlanned  = color(BorderPlannedFocused, BorderPlannedFocused256, BorderPlannedFocused16)
	clrSurface        = color(Surface2, Surface2256, Surface216) // intentionally black
	clrSurfaceDark    = color(Surface, Surface256, Surface16)    // intentionally black
	clrGreen          = color(NPPGreen, NPPGreen256, NPPGreen16)
	clrPurple         = color(VibePurple, VibePurple256, VibePurple16)
	clrRed            = color(IssueRed, IssueRed256, IssueRed16)
	clrWhite          = color(White, White256, White16)
	clrDim            = color(Dim, Dim256, Dim16)
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

	ColumnFocusedPlanned = lipgloss.NewStyle().
				BorderStyle(boxBorder).
				BorderForeground(clrBorderPlanned).
				Padding(0, 1)

	ColumnFocusedDone = lipgloss.NewStyle().
				BorderStyle(boxBorder).
				BorderForeground(clrGreen).
				Padding(0, 1)

	ColumnTitle = lipgloss.NewStyle().
			Foreground(clrText).
			Bold(true)

	ColumnTitleFocused = lipgloss.NewStyle().
				Foreground(clrOrange).
				Bold(true)

	ColumnTitleFocusedPlanned = lipgloss.NewStyle().
					Foreground(clrBorderPlanned).
					Bold(true)

	ColumnTitleFocusedDone = lipgloss.NewStyle().
					Foreground(clrGreen).
					Bold(true)

	TaskItem = lipgloss.NewStyle().
			Foreground(clrText)

	TaskItemFocused = lipgloss.NewStyle().
			Foreground(clrOrange)

	TaskItemFocusedDone = lipgloss.NewStyle().
				Foreground(clrGreen)

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

	// ObjectiveItemFocused, SidebarTitleFocused, and SidebarPanelFocused are
	// the V2 sidebar's own focus accent: purple, distinct from the orange the
	// Task columns wear, so an Objective under the cursor never reads as a
	// Task card.
	ObjectiveItemFocused = lipgloss.NewStyle().
				Foreground(clrPurple)

	SidebarTitleFocused = lipgloss.NewStyle().
				Foreground(clrPurple).
				Bold(true)

	SidebarPanelFocused = lipgloss.NewStyle().
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

	CardBoxFocusedPlanned = lipgloss.NewStyle().
				BorderStyle(boxBorder).
				BorderForeground(clrBorderPlanned).
				Padding(0, 1)

	CardBoxFocusedDone = lipgloss.NewStyle().
				BorderStyle(boxBorder).
				BorderForeground(clrGreen).
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

	// V2 router-phase pills: Idea/Design/Task/Check, colored to match the
	// public site (getsavepoint.dev): white, purple, orange, green.
	FooterPhaseIdea = lipgloss.NewStyle().
			Foreground(clrWhite).
			Bold(true)

	FooterPhaseDesign = lipgloss.NewStyle().
				Foreground(clrPurple).
				Bold(true)

	FooterPhaseTask = lipgloss.NewStyle().
			Foreground(clrOrange).
			Bold(true)

	FooterPhaseCheck = lipgloss.NewStyle().
				Foreground(clrGreen).
				Bold(true)

	HeaderRight = lipgloss.NewStyle().
			Foreground(clrDim)

	HeaderRelease = lipgloss.NewStyle().
			Foreground(clrText)

	// HeaderWhiteBold and HeaderWhite are a label/value pair for the
	// Release line in the V2 selection header: a bold white label, plain
	// white following text.
	HeaderWhiteBold = lipgloss.NewStyle().
			Foreground(clrWhite).
			Bold(true)

	HeaderWhite = lipgloss.NewStyle().
			Foreground(clrWhite)

	// NextLabel is the bold orange "NEXT:" prefix the Next area leads with —
	// the same accent weight the router's Task phase wears, so the panel a
	// reader glances at first carries the boldest color on the board.
	NextLabel = lipgloss.NewStyle().
			Foreground(clrOrange).
			Bold(true)

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

	// IssueAccent is the Issues surface's own accent: Issue IDs and the
	// Issues headings wear it so the panel reads as a different place from
	// the Task board it replaces on screen. Status columns keep their own
	// status colors.
	IssueAccent = lipgloss.NewStyle().
			Foreground(clrRed).
			Bold(true)

	// IssueItemFocused and IssueColumnFocused carry the same red to the Open
	// column's selected row title and border.
	IssueItemFocused = lipgloss.NewStyle().
				Foreground(clrRed)

	IssueColumnFocused = lipgloss.NewStyle().
				BorderStyle(boxBorder).
				BorderForeground(clrRed).
				Padding(0, 1)
)
