package board

import (
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/testutil"
)

func sampleAuditSet() data.AuditRegisterSet {
	return data.AuditRegisterSet{
		Prompt: data.AuditPrompt{
			Path:      "audit/prompt.md",
			Available: true,
			Body:      "# Audit Prompt\n\nReview the changed files for correctness.",
		},
		Findings: []data.AuditFinding{
			{ID: "F001", Title: "Nil deref in loader", Status: data.FindingOpen, Severity: data.SeverityCritical},
			{ID: "F002", Title: "Unhandled error", Status: data.FindingVerified, Severity: data.SeverityLow},
		},
		Runs: []data.AuditRun{
			{Date: "2026-06-15", Label: "full-sweep", Mode: data.AuditModeFull},
			{Date: "2026-05-01", Label: "targeted-io", Mode: data.AuditModeTargeted},
		},
	}
}

// --- renderer ---

func TestRenderAuditOverlay_headerAndTabStrip(t *testing.T) {
	got := RenderAuditOverlay(sampleAuditSet(), auditTabPrompt, 0, 70, 40, 0)
	if !strings.Contains(got, "AUDIT REGISTER") {
		t.Error("RenderAuditOverlay missing AUDIT REGISTER header")
	}
	for _, label := range []string{"Prompt", "Findings", "Runs"} {
		if !strings.Contains(got, label) {
			t.Errorf("RenderAuditOverlay missing tab label %q", label)
		}
	}
}

func TestRenderAuditOverlay_promptBody(t *testing.T) {
	got := RenderAuditOverlay(sampleAuditSet(), auditTabPrompt, 0, 70, 40, 0)
	if !strings.Contains(got, "Review the changed files") {
		t.Error("Prompt tab should render the prompt body")
	}
	if strings.Contains(got, "Nil deref") || strings.Contains(got, "full-sweep") {
		t.Error("Prompt tab should not render findings or runs bodies")
	}
}

func TestRenderAuditOverlay_findingsBody(t *testing.T) {
	got := RenderAuditOverlay(sampleAuditSet(), auditTabFindings, 0, 70, 40, 0)
	if !strings.Contains(got, "F001") || !strings.Contains(got, "Nil deref in loader") {
		t.Error("Findings tab should render finding ID and title")
	}
	if strings.Contains(got, "Review the changed files") {
		t.Error("Findings tab should not render the prompt body")
	}
}

func TestRenderAuditOverlay_runsBody(t *testing.T) {
	got := RenderAuditOverlay(sampleAuditSet(), auditTabRuns, 0, 70, 40, 0)
	if !strings.Contains(got, "2026-06-15") || !strings.Contains(got, "full-sweep") {
		t.Error("Runs tab should render run date and label")
	}
}

func TestRenderAuditOverlay_footer(t *testing.T) {
	got := RenderAuditOverlay(sampleAuditSet(), auditTabPrompt, 0, 70, 40, 0)
	if !strings.Contains(got, "esc:close") {
		t.Error("RenderAuditOverlay missing esc:close footer")
	}
	if !strings.Contains(got, "[/]:tab") {
		t.Error("RenderAuditOverlay missing tab-switch hint")
	}
}

func TestRenderAuditOverlay_emptyStates(t *testing.T) {
	empty := data.AuditRegisterSet{Prompt: data.AuditPrompt{Path: "audit/prompt.md"}}

	prompt := RenderAuditOverlay(empty, auditTabPrompt, 0, 70, 40, 0)
	if !strings.Contains(prompt, "no audit prompt") || !strings.Contains(prompt, "audit/prompt.md") {
		t.Error("missing prompt should render an empty state naming the prompt path")
	}

	findings := RenderAuditOverlay(empty, auditTabFindings, 0, 70, 40, 0)
	if !strings.Contains(findings, "no audit findings") {
		t.Error("missing findings should render an empty state")
	}

	runs := RenderAuditOverlay(empty, auditTabRuns, 0, 70, 40, 0)
	if !strings.Contains(runs, "no audit runs") {
		t.Error("missing runs should render an empty state")
	}
}

func TestRenderAuditOverlay_emptyPromptBody(t *testing.T) {
	set := data.AuditRegisterSet{Prompt: data.AuditPrompt{Path: "audit/prompt.md", Available: true, Body: "   "}}
	got := RenderAuditOverlay(set, auditTabPrompt, 0, 70, 40, 0)
	if !strings.Contains(got, "audit prompt is empty") {
		t.Error("available-but-blank prompt should render an empty-body state")
	}
}

// --- board key + overlay wiring ---

func TestUpdate_AOpensAuditOverlayAndLoads(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "audit", "prompt.md"), "# Audit Prompt\nreview")

	m := NewModel(nil, "v1.4", "E32")
	m.Root = root

	got, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("A")})
	updated := requireModel(t, got)
	if updated.Overlay != OverlayAudit {
		t.Fatalf("Overlay = %q, want %q", updated.Overlay, OverlayAudit)
	}
	if updated.AuditTab != auditTabPrompt {
		t.Errorf("AuditTab = %d, want Prompt tab", updated.AuditTab)
	}
	if cmd == nil {
		t.Fatal("expected an audit-register load command")
	}

	got2, _ := updated.Update(cmd())
	updated2 := requireModel(t, got2)
	if !updated2.Audit.Prompt.Available || !strings.Contains(updated2.Audit.Prompt.Body, "review") {
		t.Errorf("audit prompt = %+v, want the loaded prompt", updated2.Audit.Prompt)
	}
}

func TestUpdate_auditOverlayTabSwitching(t *testing.T) {
	m := NewModel(nil, "v1.4", "E32")
	m.Overlay = OverlayAudit
	m.AuditOffsets = map[auditTab]int{}

	// ] advances Prompt -> Findings -> Runs and clamps at the last tab.
	for _, key := range []string{"]", "]", "]"} {
		got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)})
		m = requireModel(t, got)
	}
	if m.AuditTab != auditTabRuns {
		t.Fatalf("AuditTab after three ] = %d, want Runs (clamped)", m.AuditTab)
	}

	// [ retreats and h retreats consistently with document overlays.
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("[")})
	m = requireModel(t, got)
	if m.AuditTab != auditTabFindings {
		t.Fatalf("AuditTab after [ = %d, want Findings", m.AuditTab)
	}
	got, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("h")})
	m = requireModel(t, got)
	if m.AuditTab != auditTabPrompt {
		t.Fatalf("AuditTab after h = %d, want Prompt", m.AuditTab)
	}
	// [ at the first tab clamps (no wrap).
	got, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("[")})
	if requireModel(t, got).AuditTab != auditTabPrompt {
		t.Error("[ at first tab should clamp at Prompt")
	}
	// l advances like right arrow.
	got, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")})
	if requireModel(t, got).AuditTab != auditTabFindings {
		t.Error("l should advance to Findings")
	}
}

func TestUpdate_auditOverlayScrollPerTab(t *testing.T) {
	m := NewModel(nil, "v1.4", "E32")
	m.Overlay = OverlayAudit
	m.AuditOffsets = map[auditTab]int{}

	// Scroll the Prompt tab down twice, then switch to Findings with ].
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	got, _ = requireModel(t, got).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	got, _ = requireModel(t, got).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("]")})
	m = requireModel(t, got)

	if m.AuditOffsets[auditTabPrompt] != 2 {
		t.Errorf("Prompt offset = %d, want 2 (preserved across switch)", m.AuditOffsets[auditTabPrompt])
	}
	if m.AuditOffsets[auditTabFindings] != 0 {
		t.Errorf("Findings offset = %d, want 0", m.AuditOffsets[auditTabFindings])
	}

	// Over-scroll up on Findings clamps at the top.
	got, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	m = requireModel(t, got)
	if m.AuditOffsets[auditTabFindings] != 0 {
		t.Errorf("Findings offset after over-scroll up = %d, want 0", m.AuditOffsets[auditTabFindings])
	}

	// pgdown advances by a page; pgup returns to the top.
	got, _ = m.Update(tea.KeyMsg{Type: tea.KeyPgDown})
	m = requireModel(t, got)
	if m.AuditOffsets[auditTabFindings] <= 0 {
		t.Errorf("Findings offset after pgdown = %d, want > 0", m.AuditOffsets[auditTabFindings])
	}
	got, _ = m.Update(tea.KeyMsg{Type: tea.KeyPgUp})
	m = requireModel(t, got)
	if m.AuditOffsets[auditTabFindings] != 0 {
		t.Errorf("Findings offset after pgup = %d, want 0", m.AuditOffsets[auditTabFindings])
	}
}

func TestUpdate_auditOverlayClosesOnEscAndQ(t *testing.T) {
	for _, key := range []tea.KeyMsg{{Type: tea.KeyEsc}, {Type: tea.KeyRunes, Runes: []rune("q")}} {
		m := NewModel(nil, "v1.4", "E32")
		m.Overlay = OverlayAudit
		got, _ := m.Update(key)
		if requireModel(t, got).Overlay != OverlayNone {
			t.Errorf("key %v should close the audit overlay", key)
		}
	}
}

func TestView_auditOverlayRendered(t *testing.T) {
	m := NewModel(nil, "v1.4", "E32")
	m.Width = 120
	m.Height = 30
	m.Overlay = OverlayAudit
	m.Audit = sampleAuditSet()

	got := m.View()
	if !strings.Contains(got, "AUDIT REGISTER") {
		t.Error("View() with OverlayAudit should render the AUDIT REGISTER overlay")
	}
}
