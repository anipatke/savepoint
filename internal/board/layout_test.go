package board

import "testing"

func TestCalculateLayout_wide(t *testing.T) {
	l := CalculateLayout(120, 40)
	if !l.EpicPanelVisible {
		t.Error("expected epic panel visible at width=120")
	}
	if l.EpicPanelWidth != epicPanelWidth {
		t.Errorf("EpicPanelWidth = %d, want %d", l.EpicPanelWidth, epicPanelWidth)
	}
	if l.ColCount != 3 {
		t.Errorf("ColCount = %d, want 3", l.ColCount)
	}
	if len(l.ColWidths) != 3 {
		t.Fatalf("len(ColWidths) = %d, want 3", len(l.ColWidths))
	}
	// inner = 120-4 = 116; available = 116-(28+4)-3*4 = 72; cw = 72/3 = 24
	for i, w := range l.ColWidths {
		if w != 24 {
			t.Errorf("ColWidths[%d] = %d, want 24", i, w)
		}
	}
}

func TestCalculateLayout_medium(t *testing.T) {
	l := CalculateLayout(100, 40)
	if l.EpicPanelVisible {
		t.Error("epic panel should be hidden at width=100")
	}
	if l.ColCount != 3 {
		t.Errorf("ColCount = %d, want 3", l.ColCount)
	}
	if len(l.ColWidths) != 3 {
		t.Fatalf("len(ColWidths) = %d, want 3", len(l.ColWidths))
	}
	for i, w := range l.ColWidths {
		if w != 28 {
			t.Errorf("ColWidths[%d] = %d, want 28", i, w)
		}
	}
}

func TestCalculateLayout_narrow(t *testing.T) {
	l := CalculateLayout(60, 40)
	if l.EpicPanelVisible {
		t.Error("epic panel should be hidden at width=60")
	}
	if l.ColCount != 1 {
		t.Errorf("ColCount = %d, want 1", l.ColCount)
	}
	// inner = 60 - 4 = 56; cw = 56 - 4 = 52
	if l.ColWidths[0] != 52 {
		t.Errorf("ColWidths[0] = %d, want 52", l.ColWidths[0])
	}
}

func TestCalculateLayout_tinyWidth_floorsAtMinColWidth(t *testing.T) {
	l := CalculateLayout(4, 40)
	if l.ColCount != 1 {
		t.Errorf("ColCount = %d, want 1", l.ColCount)
	}
	if l.ColWidths[0] != minColWidth {
		t.Errorf("ColWidths[0] = %d, want %d (minColWidth floor)", l.ColWidths[0], minColWidth)
	}
}

func TestCalculateLayout_breakpointBoundaries(t *testing.T) {
	cases := []struct {
		width        int
		wantColCount int
		wantEpic     bool
	}{
		{119, 3, false},
		{120, 3, true},
		{79, 1, false},
		{80, 3, false},
	}
	for _, tc := range cases {
		l := CalculateLayout(tc.width, 40)
		if l.ColCount != tc.wantColCount {
			t.Errorf("width=%d: ColCount = %d, want %d", tc.width, l.ColCount, tc.wantColCount)
		}
		if l.EpicPanelVisible != tc.wantEpic {
			t.Errorf("width=%d: EpicPanelVisible = %v, want %v", tc.width, l.EpicPanelVisible, tc.wantEpic)
		}
	}
}
