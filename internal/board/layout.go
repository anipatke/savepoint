package board

const (
	colOverhead = 4 // rounded border (1) + padding (1) each side

	minColWidth = 10

	epicPanelWidth    = 28
	epicPanelOverhead = 4

	boardFrameOverhead = 4 // rounded border (2) + padding (2×1)

	breakpointWide   = 120
	breakpointNarrow = 80
)

// Layout describes board geometry for a given terminal size.
type Layout struct {
	EpicPanelVisible bool
	EpicPanelWidth   int
	ColCount         int
	ColWidths        []int
}

// CalculateLayout returns the board layout for the given terminal dimensions.
//
//   - >=120 cols: epic panel (28w) + 3 columns
//   - 80–119 cols: 3 columns only
//   - <80 cols: 1 column
func CalculateLayout(width, _ int) Layout {
	inner := width - boardFrameOverhead
	switch {
	case width >= breakpointWide:
		available := inner - (epicPanelWidth + epicPanelOverhead) - 3*colOverhead
		cw := max(available/3, minColWidth)
		return Layout{
			EpicPanelVisible: true,
			EpicPanelWidth:   epicPanelWidth,
			ColCount:         3,
			ColWidths:        []int{cw, cw, cw},
		}
	case width >= breakpointNarrow:
		available := inner - 3*colOverhead
		cw := max(available/3, minColWidth)
		return Layout{
			EpicPanelVisible: false,
			ColCount:         3,
			ColWidths:        []int{cw, cw, cw},
		}
	default:
		cw := max(inner-colOverhead, minColWidth)
		return Layout{
			EpicPanelVisible: false,
			ColCount:         1,
			ColWidths:        []int{cw},
		}
	}
}
