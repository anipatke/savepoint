package board

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/opencode/savepoint/internal/data"
)

// linkedFindings returns findings linking back to the task E32-audit/T005 (one by
// full ID, one by short ID) plus one unrelated finding, so tests can assert the
// reverse lookup, cursor bounds, and empty-state handling.
func linkedFindings() []data.AuditFinding {
	return []data.AuditFinding{
		{ID: "F001", Title: "Nil deref", Status: data.FindingOpen, Severity: data.SeverityCritical, Confidence: data.ConfidenceHigh, Tasks: []string{"E32-audit/T005-backlinks"}, Epics: []string{"E32-audit"}},
		{ID: "F002", Title: "Slow load", Status: data.FindingOpen, Severity: data.SeverityLow, Confidence: data.ConfidenceLow, Tasks: []string{"T005"}, Epics: []string{"E32"}},
		{ID: "F003", Title: "Unrelated", Status: data.FindingOpen, Severity: data.SeverityMedium, Confidence: data.ConfidenceMedium, Tasks: []string{"E31-other/T009-x"}, Epics: []string{"E31-other"}},
	}
}

func linkedDetailModel() Model {
	task := data.Task{ID: "E32-audit/T005-backlinks", Title: "Backlinks", Epic: "E32-audit", Release: "v1.4", Column: data.ColumnInProgress, Stage: data.StageBuild}
	m := NewModel([]data.Task{task}, "v1.4", "E32-audit")
	m.FocusedColumn = data.ColumnInProgress
	m.FocusedTask = 0
	m.Audit = data.AuditRegisterSet{Findings: linkedFindings()}
	m.Overlay = OverlayDetail
	return m
}

func linkedEpicModel() Model {
	m := NewModel(nil, "v1.4", "E32-audit")
	m.EpicDetailEpic = "E32-audit"
	m.EpicDetailTab = 0
	m.Audit = data.AuditRegisterSet{Findings: linkedFindings()}
	m.Overlay = OverlayEpicDetail
	return m
}

// --- rendering ---

func TestRenderDetail_linkedFindingsSection(t *testing.T) {
	findings := data.FindingsForTask(linkedFindings(), "E32-audit/T005-backlinks")
	got := plainTerminal(RenderDetail(sampleTask(), 70, nil, 40, 0, findings, 0))
	for _, want := range []string{"Linked Findings:", "F001", "Nil deref", "open", "F002"} {
		if !strings.Contains(got, want) {
			t.Errorf("task detail linked findings missing %q", want)
		}
	}
	if strings.Contains(got, "F003") {
		t.Error("task detail should not list findings that do not link to the task")
	}
	if !strings.Contains(got, "enter:open") {
		t.Error("task detail footer should advertise enter:open when findings are linked")
	}
}

func TestRenderDetail_linkedFindingsEmptyState(t *testing.T) {
	got := plainTerminal(RenderDetail(sampleTask(), 70, nil, 40, 0, nil, 0))
	if !strings.Contains(got, "Linked Findings:") {
		t.Error("task detail should always render the Linked Findings section")
	}
	if !strings.Contains(got, "(no linked findings)") {
		t.Error("task detail with no findings should render the empty state")
	}
	if strings.Contains(got, "enter:open") {
		t.Error("task detail without findings should not advertise enter:open")
	}
}

func TestRenderEpicDetail_linkedFindingsSection(t *testing.T) {
	findings := data.FindingsForEpic(linkedFindings(), "E32-audit")
	got := plainTerminal(RenderEpicDetail("E32-audit", "# Epic\n\nbody", 70, 40, 0, 0, findings, 0))
	for _, want := range []string{"Linked Findings:", "F001", "F002"} {
		if !strings.Contains(got, want) {
			t.Errorf("epic detail linked findings missing %q", want)
		}
	}
	if strings.Contains(got, "F003") {
		t.Error("epic detail should not list findings that do not link to the epic")
	}
}

func TestRenderEpicDetail_linkedFindingsEmptyState(t *testing.T) {
	got := plainTerminal(RenderEpicDetail("E32-audit", "# Epic", 70, 40, 0, 0, nil, 0))
	if !strings.Contains(got, "(no linked findings)") {
		t.Error("epic detail with no findings should render the empty state")
	}
}

// --- task detail navigation ---

func TestUpdate_detailLinkedFindingCursorMovesAndClamps(t *testing.T) {
	m := linkedDetailModel()
	for _, key := range []string{"j", "j"} { // two linked findings; second j clamps
		got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)})
		m = requireModel(t, got)
	}
	if m.LinkedFindingCursor != 1 {
		t.Fatalf("LinkedFindingCursor after two downs = %d, want 1 (clamped)", m.LinkedFindingCursor)
	}
	if m.DetailOffset != 0 {
		t.Errorf("DetailOffset = %d, want 0: findings should drive the cursor, not scroll", m.DetailOffset)
	}
	for i := 0; i < 3; i++ {
		got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
		m = requireModel(t, got)
	}
	if m.LinkedFindingCursor != 0 {
		t.Errorf("LinkedFindingCursor after over-scroll up = %d, want 0", m.LinkedFindingCursor)
	}
}

func TestUpdate_detailEnterOpensFindingAndReturns(t *testing.T) {
	m := linkedDetailModel()
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}) // cursor -> F002
	m = requireModel(t, got)

	got, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = requireModel(t, got)
	if cmd != nil {
		t.Error("opening finding detail should not emit a command")
	}
	if m.Overlay != OverlayFindingDetail {
		t.Fatalf("Overlay after enter = %q, want %q", m.Overlay, OverlayFindingDetail)
	}
	if m.FindingDetailOrigin != OverlayDetail {
		t.Errorf("FindingDetailOrigin = %q, want %q", m.FindingDetailOrigin, OverlayDetail)
	}
	if f, ok := m.activeFinding(); !ok || f.ID != "F002" {
		t.Errorf("activeFinding = %+v ok=%v, want F002", f, ok)
	}

	got, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = requireModel(t, got)
	if m.Overlay != OverlayDetail {
		t.Errorf("esc from finding detail = %q, want %q", m.Overlay, OverlayDetail)
	}
	if m.LinkedFindingCursor != 1 {
		t.Errorf("LinkedFindingCursor after round-trip = %d, want 1 (preserved)", m.LinkedFindingCursor)
	}
}

func TestUpdate_boardEnterResetsLinkedFindingCursor(t *testing.T) {
	m := linkedDetailModel()
	m.Overlay = OverlayNone
	m.LinkedFindingCursor = 5
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = requireModel(t, got)
	if m.Overlay != OverlayDetail || m.LinkedFindingCursor != 0 {
		t.Errorf("opening detail: overlay=%q cursor=%d, want detail/0", m.Overlay, m.LinkedFindingCursor)
	}
}

// --- epic detail navigation ---

func TestUpdate_epicDetailEnterOpensFindingAndReturns(t *testing.T) {
	m := linkedEpicModel()
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = requireModel(t, got)
	if m.Overlay != OverlayFindingDetail {
		t.Fatalf("Overlay after enter = %q, want %q", m.Overlay, OverlayFindingDetail)
	}
	if m.FindingDetailOrigin != OverlayEpicDetail {
		t.Errorf("FindingDetailOrigin = %q, want %q", m.FindingDetailOrigin, OverlayEpicDetail)
	}
	if f, ok := m.activeFinding(); !ok || f.ID != "F001" {
		t.Errorf("activeFinding = %+v ok=%v, want F001", f, ok)
	}

	got, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	m = requireModel(t, got)
	if m.Overlay != OverlayEpicDetail {
		t.Errorf("q from finding detail = %q, want %q", m.Overlay, OverlayEpicDetail)
	}
}

func TestUpdate_epicDetailAuditTabScrollsNotCursor(t *testing.T) {
	m := linkedEpicModel()
	m.EpicDetailTab = 1 // Audit tab: findings must not capture the cursor keys
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	m = requireModel(t, got)
	if m.EpicDetailOffset != 1 {
		t.Errorf("Audit tab j = EpicDetailOffset %d, want 1 (scroll)", m.EpicDetailOffset)
	}
	if m.LinkedFindingCursor != 0 {
		t.Errorf("Audit tab j moved LinkedFindingCursor to %d, want 0", m.LinkedFindingCursor)
	}
}

func TestUpdate_epicDetailTabSwitchResetsFindingCursor(t *testing.T) {
	m := linkedEpicModel()
	m.LinkedFindingCursor = 1
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("1")})
	m = requireModel(t, got)
	if m.LinkedFindingCursor != 0 {
		t.Errorf("LinkedFindingCursor after tab switch = %d, want 0", m.LinkedFindingCursor)
	}
}

// --- read-only guarantee ---

func TestUpdate_detailBacklinkKeysNeverMutateFindings(t *testing.T) {
	for _, overlay := range []OverlayType{OverlayDetail, OverlayFindingDetail} {
		for _, key := range []tea.KeyMsg{
			{Type: tea.KeyRunes, Runes: []rune(" ")},
			{Type: tea.KeyBackspace},
			{Type: tea.KeyRunes, Runes: []rune("p")},
			{Type: tea.KeyEnter},
		} {
			m := linkedDetailModel()
			m.Overlay = overlay
			m.FindingDetailOrigin = OverlayDetail
			before := append([]data.AuditFinding(nil), m.Audit.Findings...)

			got, cmd := m.Update(key)
			m = requireModel(t, got)
			if cmd != nil {
				t.Errorf("overlay %q key %v should not emit a write command", overlay, key)
			}
			for i, f := range m.Audit.Findings {
				if f.Status != before[i].Status {
					t.Errorf("overlay %q key %v mutated finding %s status", overlay, key, f.ID)
				}
			}
		}
	}
}
