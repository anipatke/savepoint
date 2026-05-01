package styles

import "github.com/charmbracelet/lipgloss"

func color(hex, ansi256, ansi string) lipgloss.CompleteColor {
	return lipgloss.CompleteColor{TrueColor: hex, ANSI256: ansi256, ANSI: ansi}
}

var (
	clrOrange      = color(AtariOrange, AtariOrange256, AtariOrange16)
	clrText        = color(PrimaryText, PrimaryText256, PrimaryText16)
	clrBorder      = color(BorderSubtle, BorderSubtle256, BorderSubtle16)
	clrSurface     = color(Surface2, Surface2256, Surface216) // #0F0F0F
	clrSurfaceDark = color(Surface, Surface256, Surface16)    // #0D0D0D
	clrGreen       = color(NPPGreen, NPPGreen256, NPPGreen16)
	clrPurple      = color(VibePurple, VibePurple256, VibePurple16)
	clrDim         = color(Dim, Dim256, Dim16)
)

var (
	HeaderIcon = lipgloss.NewStyle().
			Foreground(clrOrange).
			Bold(true)

	HeaderText = lipgloss.NewStyle().
			Foreground(clrText)

	Divider = lipgloss.NewStyle().
		Foreground(clrBorder)

	HeaderFrame = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(clrBorder).
			Padding(1, 1)

	BoardFrame = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(clrBorder).
			Padding(0, 1)

	Column = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(clrBorder).
		Background(clrSurfaceDark).
		Padding(0, 1)

	ColumnFocused = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(clrOrange).
			Background(clrSurfaceDark).
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
			Foreground(clrText).
			Background(clrSurface)

	EpicPanel = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(clrPurple).
			Background(clrSurface).
			Padding(0, 1)

	Card = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(clrBorder).
		Background(clrSurface).
		Padding(0, 1)

	CardFocused = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(clrOrange).
			Background(clrSurface).
			Padding(0, 1)

	CardMeta = lipgloss.NewStyle().Foreground(clrDim)

	GlyphBuild = lipgloss.NewStyle().Foreground(clrOrange)
	GlyphTest  = lipgloss.NewStyle().Foreground(clrGreen)
	GlyphAudit = lipgloss.NewStyle().Foreground(clrPurple)

	DetailOverlay = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(clrOrange).
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

	FooterDivider = lipgloss.NewStyle().
			Foreground(clrBorder)

	FooterHints = lipgloss.NewStyle().
			Foreground(clrDim)

	// Tag styles for semantic encoding
	TagDone = lipgloss.NewStyle().Foreground(clrGreen)
	TagAI   = lipgloss.NewStyle().Foreground(clrPurple)
)
