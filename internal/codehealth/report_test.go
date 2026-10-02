package codehealth

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// reportFor saves an official snapshot built from tweak and renders its report.
func reportFor(t *testing.T, tweak func(*CapabilityResult)) string {
	t.Helper()
	store, root := dashProject(t)
	cfg := dashConfig(t, store, allCapabilities()...)
	saveDash(t, store, cfg, OriginOfficial, 1, goodResults(cfg, tweak), dashRepo(1))
	text, err := RenderReport(mustLoad(t, root), "O-036")
	if err != nil {
		t.Fatalf("RenderReport: %v", err)
	}
	return text
}

// saveOfficialRuns stores enough identical official checks for a good value to
// be judged Good rather than Watch.
func saveOfficialRuns(t *testing.T, store Store, cfg Config, tweak func(*CapabilityResult)) {
	t.Helper()
	for n := 1; n <= 3; n++ {
		saveDash(t, store, cfg, OriginOfficial, n, goodResults(cfg, tweak), dashRepo(n))
	}
}

func headingOrder(text string) []string {
	var out []string
	for _, l := range strings.Split(text, "\n") {
		if h, ok := strings.CutPrefix(l, "## "); ok {
			out = append(out, h)
		}
	}
	return out
}

func TestReportBriefNamesRerunCommand(t *testing.T) {
	text := reportFor(t, nil)
	for _, want := range []string{"Needs attention", "blocking ones first", "propose a fix, then apply it", "Stop once a signal reaches Watch or better", "report any signals that remain", "`savepoint health check O-036`"} {
		if !strings.Contains(text, want) {
			t.Errorf("report lacks %q:\n%s", want, text)
		}
	}
	if !strings.HasPrefix(text, reportTitle+"\n\n"+strings.Split(reportBrief, "%s")[0]) {
		t.Errorf("report does not open with the brief:\n%s", text)
	}
}

func TestReportShowsEverySignalWithDashboardWording(t *testing.T) {
	store, root := dashProject(t)
	cfg := dashConfig(t, store, allCapabilities()...)
	saveDash(t, store, cfg, OriginOfficial, 1, goodResults(cfg, func(r *CapabilityResult) {
		if r.Capability == CapabilityTests {
			r.Value.Number = 3
		}
	}), dashRepo(1))
	d := mustLoad(t, root)
	text, err := RenderReport(d, "O-036")
	if err != nil {
		t.Fatal(err)
	}
	if len(headingOrder(text)) != len(Capabilities()) {
		t.Fatalf("headings = %v, want one per signal", headingOrder(text))
	}
	for _, row := range d.Rows {
		for _, want := range []string{row.CapabilityText, row.Question, row.LabelText, row.Value, row.Aim, row.Meaning, row.NextStep, row.SignOff} {
			if want == "" || !strings.Contains(text, want) {
				t.Errorf("%s: report lacks %q", row.Capability, want)
			}
		}
	}
	if !strings.Contains(text, textSignOffBlocks) {
		t.Errorf("a failing test should read %q:\n%s", textSignOffBlocks, text)
	}
}

func TestReportOrdersByLabelKeepingDashboardOrder(t *testing.T) {
	store, root := dashProject(t)
	// Tests, coverage, complexity, duplication are set up; vulnerabilities are not (Unknown).
	cfg := dashConfig(t, store, CapabilityTests, CapabilityCoverage, CapabilityComplexity, CapabilityDuplication)
	saveOfficialRuns(t, store, cfg, func(r *CapabilityResult) {
		switch r.Capability {
		case CapabilityComplexity:
			r.Value.Number = 500
		case CapabilityCoverage:
			r.Value.Number = 60
		}
	})
	d := mustLoad(t, root)
	text, err := RenderReport(d, "O-036")
	if err != nil {
		t.Fatal(err)
	}
	var wantLabels []string
	for _, l := range []Classification{ClassificationNeedsAttention, ClassificationWatch, ClassificationUnknown, ClassificationGood} {
		for _, row := range d.Rows {
			if row.Label == l {
				wantLabels = append(wantLabels, row.CapabilityText)
			}
		}
	}
	got := headingOrder(text)
	if strings.Join(got, "|") != strings.Join(wantLabels, "|") {
		t.Fatalf("order = %v, want %v", got, wantLabels)
	}
	if got[len(got)-1] == capabilityText[CapabilityDependencyVulnerability] {
		t.Errorf("Unknown signal ranked after Good: %v", got)
	}
	if !slicesIndexBefore(got, capabilityText[CapabilityComplexity], capabilityText[CapabilityTests]) {
		t.Errorf("red complexity should come before good tests: %v", got)
	}
}

func slicesIndexBefore(list []string, first, second string) bool {
	a, b := -1, -1
	for i, s := range list {
		if s == first {
			a = i
		}
		if s == second {
			b = i
		}
	}
	return a >= 0 && b >= 0 && a < b
}

func TestReportEveryLabelAppears(t *testing.T) {
	store, root := dashProject(t)
	cfg := dashConfig(t, store, CapabilityTests, CapabilityCoverage, CapabilityComplexity, CapabilityDuplication)
	saveOfficialRuns(t, store, cfg, func(r *CapabilityResult) {
		switch r.Capability {
		case CapabilityComplexity:
			r.Value.Number = 500
		case CapabilityCoverage:
			r.Value.Number = 60
		}
	})
	d := mustLoad(t, root)
	seen := map[Classification]bool{}
	for _, row := range d.Rows {
		seen[row.Label] = true
	}
	for _, l := range []Classification{ClassificationNeedsAttention, ClassificationUnknown, ClassificationGood} {
		if !seen[l] {
			t.Fatalf("fixture lacks label %s: %v", l, seen)
		}
	}
	text, _ := RenderReport(d, "O-036")
	for _, row := range d.Rows {
		if !strings.Contains(text, "Label: "+row.LabelText) {
			t.Errorf("%s: no %q line", row.Capability, "Label: "+row.LabelText)
		}
	}
}

func TestReportEvidenceOneManyAndNone(t *testing.T) {
	text := reportFor(t, func(r *CapabilityResult) {
		switch r.Capability {
		case CapabilityComplexity:
			r.Value.Number = 500
			r.Evidence = []EvidenceRef{
				{Path: "a/alpha.go", Line: 7, Note: "function Alpha"},
				{Path: "b/zeta.go", Line: 42, Note: "function Zeta"},
				{Path: "c/gamma.go", Note: "whole file"},
			}
		case CapabilityCoverage:
			r.Evidence = []EvidenceRef{{Path: "only.go", Line: 3, Note: "one note"}}
		case CapabilityDuplication:
			r.Evidence = nil
		}
	})
	alpha, zeta, gamma := strings.Index(text, "- a/alpha.go:7 function Alpha"), strings.Index(text, "- b/zeta.go:42 function Zeta"), strings.Index(text, "- c/gamma.go whole file")
	if alpha < 0 || zeta < alpha || gamma < zeta {
		t.Errorf("complexity evidence missing or out of stored order (%d %d %d):\n%s", alpha, zeta, gamma, text)
	}
	if strings.Contains(text, "gamma.go:") {
		t.Errorf("an entry with no line must show the path only:\n%s", text)
	}
	if !strings.Contains(text, "- only.go:3 one note") {
		t.Errorf("single evidence entry missing:\n%s", text)
	}
	if !strings.Contains(text, reportNoFiles) {
		t.Errorf("a signal with no evidence should say so:\n%s", text)
	}
}

func TestReportSeveralInstancesGetOneSectionEach(t *testing.T) {
	store, root := dashProject(t)
	cfg := Config{Version: ConfigVersion}
	for _, c := range Capabilities() {
		names := []string{""}
		if c == CapabilityComplexity {
			names = []string{"backend", "frontend"}
		}
		for _, n := range names {
			scope := "**/*.go"
			if n != "" {
				scope = n + "/**/*.go"
			}
			cfg.Capabilities = append(cfg.Capabilities, CapabilityConfig{Capability: c, Name: n, Provider: dashProviders[c], Scope: []string{scope}})
		}
	}
	if _, err := store.SaveConfig(cfg); err != nil {
		t.Fatal(err)
	}
	saveDash(t, store, cfg, OriginOfficial, 1, goodResults(cfg, nil), dashRepo(1))
	text, err := RenderReport(mustLoad(t, root), "O-036")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"## Complexity (backend)", "## Complexity (frontend)"} {
		if !strings.Contains(text, want) {
			t.Errorf("report lacks %q:\n%s", want, text)
		}
	}
	if got := len(headingOrder(text)); got != 6 {
		t.Errorf("sections = %d, want 6", got)
	}
}

func TestReportManualNewestSnapshot(t *testing.T) {
	store, root := dashProject(t)
	cfg := dashConfig(t, store, allCapabilities()...)
	saveDash(t, store, cfg, OriginOfficial, 1, goodResults(cfg, nil), dashRepo(1))
	saveDash(t, store, cfg, OriginManual, 2, goodResults(cfg, nil), dashRepo(2))
	text, err := RenderReport(mustLoad(t, root), "O-036")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(text, "manual refresh") || !strings.Contains(text, textSignOffManual) {
		t.Errorf("report should say the newest snapshot is a manual refresh:\n%s", text)
	}
	if strings.Contains(text, "Sign-off:") {
		t.Errorf("a manual snapshot has no per-signal sign-off:\n%s", text)
	}
}

func TestReportStatesOfficialAndDate(t *testing.T) {
	text := reportFor(t, nil)
	if !strings.Contains(text, "official check") {
		t.Errorf("report should say the snapshot is official:\n%s", text)
	}
	if !strings.Contains(text, "Oct") {
		t.Errorf("report should state the measured date:\n%s", text)
	}
}

func TestReportRefusals(t *testing.T) {
	_, root := dashProject(t)
	if _, err := RenderReport(mustLoad(t, root), "O-036"); !errors.Is(err, ErrReportNotConfigured) {
		t.Errorf("not set up: err = %v, want ErrReportNotConfigured", err)
	}
	store, root := dashProject(t)
	dashConfig(t, store, allCapabilities()...)
	text, err := RenderReport(mustLoad(t, root), "O-036")
	if !errors.Is(err, ErrReportNoSnapshot) || text != "" {
		t.Errorf("no snapshot: (%q, %v), want refusal", text, err)
	}
}

func TestReportKeepsInternalsOut(t *testing.T) {
	store, root := dashProject(t)
	cfg := dashConfig(t, store, allCapabilities()...)
	snap := saveDash(t, store, cfg, OriginOfficial, 1, goodResults(cfg, nil), dashRepo(1))
	d := mustLoad(t, root)
	text, _ := RenderReport(d, "O-036")
	banned := []string{snap.ID, strings.TrimPrefix(snap.ID, digestPrefix), snap.CreatedAt, snap.Repository.Commit, "2026-10-01T"}
	for _, p := range dashProviders {
		banned = append(banned, string(p))
	}
	for _, b := range banned {
		if strings.Contains(text, b) {
			t.Errorf("report contains %q", b)
		}
	}
}

func readReport(t *testing.T, root string) string {
	t.Helper()
	b, err := os.ReadFile(healthPath(root, reportFile))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestWriteReportCreatesFilesAndIsRepeatable(t *testing.T) {
	st, root := newProject(t)
	changed, err := st.WriteReport("first\n")
	if err != nil || !changed {
		t.Fatalf("first write = (%v, %v)", changed, err)
	}
	if got := readReport(t, root); got != "first\n" {
		t.Errorf("report = %q", got)
	}
	ignore, err := os.ReadFile(healthPath(root, gitignoreFile))
	if err != nil || string(ignore) != "report.md\n" {
		t.Errorf(".gitignore = %q, %v", ignore, err)
	}
	info, _ := os.Stat(healthPath(root, reportFile))
	changed, err = st.WriteReport("first\n")
	if err != nil || changed {
		t.Errorf("identical write = (%v, %v), want unchanged", changed, err)
	}
	after, _ := os.Stat(healthPath(root, reportFile))
	if !after.ModTime().Equal(info.ModTime()) {
		t.Errorf("identical write touched the file")
	}
}

func TestWriteReportReplacesOlderReportAtomically(t *testing.T) {
	st, root := newProject(t)
	writeFile(t, healthPath(root, reportFile), "old report\n")
	changed, err := st.WriteReport("new report\n")
	if err != nil || !changed {
		t.Fatalf("write = (%v, %v)", changed, err)
	}
	if got := readReport(t, root); got != "new report\n" {
		t.Errorf("report = %q", got)
	}
	entries, _ := os.ReadDir(healthPath(root))
	for _, e := range entries {
		if isTempName(e.Name()) {
			t.Errorf("leftover temp file %s", e.Name())
		}
	}
}

func TestWriteReportNeverOverwritesOwnersIgnoreFile(t *testing.T) {
	st, root := newProject(t)
	writeFile(t, healthPath(root, gitignoreFile), "# mine\nsnapshots/\n")
	if _, err := st.WriteReport("text\n"); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(healthPath(root, gitignoreFile))
	if string(got) != "# mine\nsnapshots/\n" {
		t.Errorf(".gitignore = %q, want the owner's file untouched", got)
	}
}

func TestWriteReportCreatesIgnoreFileOnce(t *testing.T) {
	st, root := newProject(t)
	for _, text := range []string{"one\n", "two\n"} {
		if _, err := st.WriteReport(text); err != nil {
			t.Fatal(err)
		}
	}
	path := healthPath(root, gitignoreFile)
	if err := os.WriteFile(path, []byte("report.md\nextra\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := st.WriteReport("three\n"); err != nil {
		t.Fatal(err)
	}
	if got, _ := os.ReadFile(path); string(got) != "report.md\nextra\n" {
		t.Errorf(".gitignore = %q after a later write", got)
	}
}

func TestWriteReportUnwritableDirectory(t *testing.T) {
	st, root := newProject(t)
	if _, err := st.WriteReport("one\n"); err != nil {
		t.Fatal(err)
	}
	readOnlyOrSkip(t, healthPath(root))
	if _, err := st.WriteReport("two\n"); err == nil {
		t.Fatal("write into an unwritable directory succeeded")
	}
	if got := readReport(t, root); got != "one\n" {
		t.Errorf("failed write changed the report to %q", got)
	}
}

func TestWriteReportRefusesNonProject(t *testing.T) {
	root := t.TempDir()
	if _, err := NewStore(root).WriteReport("x\n"); !errors.Is(err, ErrNotProject) {
		t.Errorf("err = %v, want ErrNotProject", err)
	}
	if _, err := os.Stat(filepath.Join(root, savepointDir)); err == nil {
		t.Error("write scaffolded .savepoint")
	}
}
