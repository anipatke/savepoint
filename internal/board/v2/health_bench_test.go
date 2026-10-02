package v2

import (
	"fmt"
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/codehealth"
)

// Benchmarks for rendering the Code Health popover (T-094). They run only with
// -bench, so they never fail a gate on machine speed.
//
//	go test ./internal/board/v2 -run '^$' -bench 'HealthRender' -benchmem -count=5

// benchHistory is n earlier checks, newest first, as LoadDashboard would list
// them if it did not cap the list. Passing more than MaxDashboardHistory shows
// whether rendering stays bounded on its own.
func benchHistory(n int) []codehealth.DashboardHistoryEntry {
	out := make([]codehealth.DashboardHistoryEntry, n)
	for i := range out {
		origin, text := codehealth.OriginOfficial, "Official check"
		if i%3 == 2 {
			origin, text = codehealth.OriginManual, "Manual refresh"
		}
		out[i] = codehealth.DashboardHistoryEntry{
			CreatedAt: fmt.Sprintf("2026-10-01T%02d:%02d:00Z", i/60%24, i%60), WhenText: fmt.Sprintf("%d Oct 2026", i%28+1),
			Origin: origin, OriginText: text, Overall: codehealth.ClassificationGood, OverallText: "Good",
		}
	}
	return out
}

// benchRows is the five signals with instances copies of each, as a project
// with several configured instances per signal produces.
func benchRows(instances int) []codehealth.DashboardRow {
	var rows []codehealth.DashboardRow
	for _, row := range realCopyRows() {
		for i := 0; i < instances; i++ {
			r := row
			if instances > 1 {
				r.Name = fmt.Sprintf("inst%d", i)
			}
			rows = append(rows, r)
		}
	}
	return rows
}

func benchHealthModel(rows []codehealth.DashboardRow, history int, showHistory bool) Model {
	d := measuredDashboard(rows)
	d.History = benchHistory(history)
	return Model{Health: &HealthOverlay{
		Dashboard: &d, Loaded: true, History: showHistory,
		Freshness: &codehealth.CodeFreshness{Text: "Your code matches this check."},
		Rerun:     "Re-run: savepoint health check O-001",
	}}
}

// BenchmarkHealthRenderPopover is the fixed 72x17 popover for the normal
// view, the history view and a grouped multi-instance view.
func BenchmarkHealthRenderPopover(b *testing.B) {
	cases := []struct {
		name      string
		instances int
		history   int
		show      bool
	}{
		{"normal/instances=1", 1, codehealth.MaxDashboardHistory, false},
		{"normal/instances=4", 4, codehealth.MaxDashboardHistory, false},
		{"normal/instances=12", 12, codehealth.MaxDashboardHistory, false},
		{"history/entries=10", 1, 10, true},
		{"history/entries=1000", 1, 1000, true},
	}
	for _, c := range cases {
		m := benchHealthModel(benchRows(c.instances), c.history, c.show)
		b.Run(c.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				m.renderHealthPopover(120)
			}
		})
	}
}

// BenchmarkHealthRenderOverlay is the popover composed over a full-screen base
// view, which is what each frame draws while the popover is open.
func BenchmarkHealthRenderOverlay(b *testing.B) {
	for _, size := range [][2]int{{80, 24}, {200, 50}} {
		base := strings.TrimSuffix(strings.Repeat(strings.Repeat("board ", size[0]/6)+"\n", size[1]), "\n")
		m := benchHealthModel(benchRows(1), codehealth.MaxDashboardHistory, false)
		b.Run(fmt.Sprintf("%dx%d", size[0], size[1]), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				m.renderHealthOverlay(base, size[0], size[1])
			}
		})
	}
}

// TestHealthBenchFixturesFillTheFixedFrame keeps the benchmark inputs honest
// without timing anything: every case renders the same 17-line frame.
func TestHealthBenchFixturesFillTheFixedFrame(t *testing.T) {
	for _, instances := range []int{1, 4, 12} {
		for _, show := range []bool{false, true} {
			m := benchHealthModel(benchRows(instances), 1000, show)
			if got := strings.Count(m.renderHealthPopover(120), "\n") + 1; got != healthPopoverHeight {
				t.Errorf("instances=%d history=%v: %d lines, want %d", instances, show, got, healthPopoverHeight)
			}
		}
	}
}
