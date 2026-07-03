package board

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/opencode/savepoint/internal/data"
)

// findingsSet returns an audit set whose findings span two statuses and three
// severities so tests can exercise grouping, cursor movement, and detail content.
func findingsSet() data.AuditRegisterSet {
	return data.AuditRegisterSet{
		Findings: []data.AuditFinding{
			{ID: "F001", Title: "Nil deref in loader", Status: data.FindingOpen, Severity: data.SeverityCritical, Confidence: data.ConfidenceHigh, WorkItem: "E32/T001"},
			{ID: "F002", Title: "Unhandled error", Status: data.FindingOpen, Severity: data.SeverityMedium, Confidence: data.ConfidenceMedium},
			{ID: "F003", Title: "Stale cache", Status: data.FindingVerified, Severity: data.SeverityLow, Confidence: data.ConfidenceLow},
		},
	}
}

// detailFinding returns a fully populated finding for detail-content assertions.
func detailFinding() data.AuditFinding {
	return data.AuditFinding{
		ID:          "F007",
		Title:       "Race in watcher",
		Status:      data.FindingMapped,
		Severity:    data.SeverityHigh,
		Confidence:  data.ConfidenceHigh,
		WorkItem:    "E32/T003",
		Tasks:       []string{"T003", "T004"},
		Defects:     []string{"D012"},
		Locations:   []string{"internal/board/watch.go:42"},
		ProofNeeded: "Reproduce the data race under -race.",
		FirstSeen:   "2026-05-01",
		LastSeen:    "2026-06-15",
		Body:        "## Analysis\n\nThe watcher goroutine reads Model state.\n\n## Repro\n\nRun the board with two writers.",
	}
}

// --- findings tab rendering ---

func TestRenderAuditOverlay_findingsGroupedByStatus(t *testing.T) {
	got := RenderAuditOverlay(findingsSet(), auditTabFindings, 0, 70, 40, 0)

	for _, want := range []string{"OPEN", "VERIFIED", "F001", "Nil deref in loader", "conf:high", "E32/T001"} {
		if !strings.Contains(got, want) {
			t.Errorf("findings tab missing %q", want)
		}
	}
	// Findings arrive sorted active-before-resolved, so the OPEN group and its
	// critical F001 precede the VERIFIED group's F003.
	if strings.Index(got, "OPEN") > strings.Index(got, "VERIFIED") {
		t.Error("OPEN group should render before VERIFIED group")
	}
	if strings.Index(got, "F001") > strings.Index(got, "F003") {
		t.Error("F001 (open/critical) should render before F003 (verified/low)")
	}
}

func TestRenderAuditOverlay_findingsCursorHighlight(t *testing.T) {
	first := RenderAuditOverlay(findingsSet(), auditTabFindings, 0, 70, 40, 0)
	second := RenderAuditOverlay(findingsSet(), auditTabFindings, 1, 70, 40, 0)
	if first == second {
		t.Error("moving the finding cursor should change the rendered highlight")
	}
	if !strings.Contains(second, releaseActiveMarker) {
		t.Errorf("selected finding row should carry the cursor marker %q", releaseActiveMarker)
	}
}

// --- cursor navigation ---

func auditFindingsModel() Model {
	m := NewModel(nil, "v1.4", "E32")
	m.Overlay = OverlayAudit
	m.AuditTab = auditTabFindings
	m.AuditOffsets = map[auditTab]int{}
	m.Audit = findingsSet()
	return m
}

func TestUpdate_findingCursorMovesAndClamps(t *testing.T) {
	m := auditFindingsModel()

	// j / down advance; k / up retreat; both clamp at the ends.
	for _, key := range []string{"j", "down", "j"} { // 3 downs on a 3-finding list
		got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)})
		m = requireModel(t, got)
	}
	if m.FindingCursor != 2 {
		t.Fatalf("FindingCursor after three downs = %d, want 2 (clamped at last)", m.FindingCursor)
	}

	for i := 0; i < 5; i++ {
		got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
		m = requireModel(t, got)
	}
	if m.FindingCursor != 0 {
		t.Fatalf("FindingCursor after over-scroll up = %d, want 0 (clamped at first)", m.FindingCursor)
	}
}

func TestUpdate_findingCursorNoopWhenEmpty(t *testing.T) {
	m := NewModel(nil, "v1.4", "E32")
	m.Overlay = OverlayAudit
	m.AuditTab = auditTabFindings
	m.AuditOffsets = map[auditTab]int{}

	got, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	m = requireModel(t, got)
	if m.FindingCursor != 0 {
		t.Errorf("FindingCursor with no findings = %d, want 0", m.FindingCursor)
	}
	if cmd != nil {
		t.Error("finding cursor movement should not emit a command")
	}
}

// --- enter / detail / return ---

func TestUpdate_findingEnterOpensDetail(t *testing.T) {
	m := auditFindingsModel()
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	m = requireModel(t, got)

	got, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = requireModel(t, got)
	if m.Overlay != OverlayFindingDetail {
		t.Fatalf("Overlay after enter = %q, want %q", m.Overlay, OverlayFindingDetail)
	}
	if m.FindingDetailOffset != 0 {
		t.Errorf("FindingDetailOffset on open = %d, want 0", m.FindingDetailOffset)
	}
	if f, ok := m.selectedFinding(); !ok || f.ID != "F002" {
		t.Errorf("selectedFinding = %+v ok=%v, want F002", f, ok)
	}
}

func TestUpdate_findingEnterNoopWhenEmpty(t *testing.T) {
	m := NewModel(nil, "v1.4", "E32")
	m.Overlay = OverlayAudit
	m.AuditTab = auditTabFindings

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if requireModel(t, got).Overlay != OverlayAudit {
		t.Error("enter with no findings should keep the audit overlay open")
	}
}

func TestUpdate_findingDetailReturnsToFindingsTab(t *testing.T) {
	for _, key := range []tea.KeyMsg{{Type: tea.KeyEsc}, {Type: tea.KeyRunes, Runes: []rune("q")}} {
		m := auditFindingsModel()
		m.Overlay = OverlayFindingDetail
		got, _ := m.Update(key)
		m = requireModel(t, got)
		if m.Overlay != OverlayAudit {
			t.Errorf("key %v should return to the audit overlay, got %q", key, m.Overlay)
		}
		if m.AuditTab != auditTabFindings {
			t.Errorf("returning from detail should stay on the Findings tab, got %d", m.AuditTab)
		}
	}
}

func TestUpdate_findingDetailScrolls(t *testing.T) {
	m := auditFindingsModel()
	m.Overlay = OverlayFindingDetail

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	m = requireModel(t, got)
	if m.FindingDetailOffset != 1 {
		t.Fatalf("FindingDetailOffset after j = %d, want 1", m.FindingDetailOffset)
	}
	// Over-scroll up clamps at the top.
	for i := 0; i < 3; i++ {
		got, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
		m = requireModel(t, got)
	}
	if m.FindingDetailOffset != 0 {
		t.Errorf("FindingDetailOffset after over-scroll up = %d, want 0", m.FindingDetailOffset)
	}
}

// --- detail content ---

func TestRenderFindingDetail_content(t *testing.T) {
	got := RenderFindingDetail(detailFinding(), 70, 40, 0)

	wants := []string{
		"FINDING DETAIL",
		"F007", "Race in watcher", "mapped", "high",
		"E32/T003",                   // work item link
		"T003, T004",                 // task links
		"D012",                       // defect link
		"internal/board/watch.go:42", // location
		"Reproduce the data race",    // proof needed
		"2026-05-01", "2026-06-15",   // first / last seen
		"Analysis", "watcher goroutine", "Repro", // body sections
		"esc:close",
	}
	for _, want := range wants {
		if !strings.Contains(got, want) {
			t.Errorf("finding detail missing %q", want)
		}
	}
}

func TestRenderFindingDetail_omitsAbsentLinks(t *testing.T) {
	f := data.AuditFinding{ID: "F009", Title: "Bare finding", Status: data.FindingOpen, Severity: data.SeverityLow, Confidence: data.ConfidenceLow, FirstSeen: "2026-01-01", LastSeen: "2026-01-02"}
	got := RenderFindingDetail(f, 70, 40, 0)
	for _, absent := range []string{"Work Item", "Tasks", "Locations", "Proof Needed"} {
		if strings.Contains(got, absent) {
			t.Errorf("bare finding detail should omit %q row", absent)
		}
	}
}

// --- read-only guarantee ---

func TestUpdate_findingsTabNeverMutatesOrWrites(t *testing.T) {
	// Keys that mutate or persist on other overlays must be inert on findings.
	for _, key := range []tea.KeyMsg{
		{Type: tea.KeyRunes, Runes: []rune(" ")},
		{Type: tea.KeyBackspace},
		{Type: tea.KeyRunes, Runes: []rune("p")},
		{Type: tea.KeyRunes, Runes: []rune("a")},
	} {
		m := auditFindingsModel()
		before := append([]data.AuditFinding(nil), m.Audit.Findings...)

		got, cmd := m.Update(key)
		m = requireModel(t, got)
		if cmd != nil {
			t.Errorf("key %v on findings tab should not emit a write command", key)
		}
		for i, f := range m.Audit.Findings {
			if f.Status != before[i].Status {
				t.Errorf("key %v mutated finding %s status: %q -> %q", key, f.ID, before[i].Status, f.Status)
			}
		}
	}
}
