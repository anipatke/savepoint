package v2

import (
	"go/build"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// v1BoardPackage is the package this one replaces. Importing it would put a V1
// board type back within reach of a V2 render path, which is the whole reason
// the V2 board is a sibling package rather than more files in that one.
const v1BoardPackage = "github.com/opencode/savepoint/internal/board"

// forbiddenReferences are the V1 record types that have no V2 meaning. Each is
// matched with a trailing word boundary, so data.TaskV2 and data.DefectStatus
// are unaffected while data.Task and data.Defect are caught.
var forbiddenReferences = []string{
	`data\.Task\b`,
	`data\.Defect\b`,
	`data\.ReleaseInfo\b`,
	`data\.EpicInfo\b`,
	`data\.AuditRegisterSet\b`,
	`data\.RouterState\b`,
}

func TestPackageDoesNotImportTheV1Board(t *testing.T) {
	pkg, err := build.ImportDir(".", build.IgnoreVendor)
	if err != nil {
		t.Fatalf("scan package imports: %v", err)
	}
	for _, imported := range append(pkg.Imports, pkg.TestImports...) {
		if imported == v1BoardPackage {
			t.Errorf("package imports %q; the V2 board must not reach a V1 board type", imported)
		}
	}
}

func TestPackageReferencesNoV1RecordType(t *testing.T) {
	patterns := make([]*regexp.Regexp, 0, len(forbiddenReferences))
	for _, pattern := range forbiddenReferences {
		patterns = append(patterns, regexp.MustCompile(pattern))
	}

	for _, path := range packageSourceFiles(t) {
		content := readCode(t, path)
		for _, pattern := range patterns {
			if match := pattern.FindString(content); match != "" {
				t.Errorf("%s references %q, a V1 record type with no V2 meaning", path, match)
			}
		}
	}
}

// retiredV1Surfaces are the navigation and overlay surfaces V2 replaces. The
// Objective sidebar is the whole of V2 navigation — one level, and not a
// directory — so no identifier here may name a release to pick, an epic to
// open, or the separate defect and audit-register collections Issues replace.
// The V1 styles carrying those names are included: reaching for one would put
// the word back in this package through the back door.
var retiredV1Surfaces = []string{"Release", "Epic", "Defect", "AuditRegister"}

func TestPackageCarriesNoReleaseOrEpicSurface(t *testing.T) {
	for _, path := range packageSourceFiles(t) {
		content := readCode(t, path)
		for _, name := range retiredV1Surfaces {
			if strings.Contains(content, name) {
				t.Errorf("%s names %q; V2 navigates Objectives and nothing else", path, name)
			}
		}
	}
}

// TestUpdateDoesNotPerformIO holds ARCH-02 structurally: the file declaring the
// reducer may not name the filesystem or the loaders, so load work cannot drift
// out of the load command and into a key branch.
func TestUpdateDoesNotPerformIO(t *testing.T) {
	reducer, content := reducerSource(t)
	for _, forbidden := range []string{"os.", "data.LoadProject", "migrate."} {
		if strings.Contains(content, forbidden) {
			t.Errorf("%s references %q; the reducer must do filesystem work through the load command only", reducer, forbidden)
		}
	}
}

// TestReloadUsesTheStartupLoadCommand proves the watcher cannot grow a
// second refresh path: its message branch must dispatch loadCmd, and the
// reducer must never call the filesystem loader directly.
func TestReloadUsesTheStartupLoadCommand(t *testing.T) {
	content := readCode(t, "update.go")
	if !strings.Contains(content, "case v2FileChangeMsg") {
		t.Fatal("update.go has no V2 watcher message branch")
	}
	if !strings.Contains(content, "loadCmd(") {
		t.Fatal("watcher reload branch does not dispatch loadCmd")
	}
	if strings.Contains(content, "loadProject(") {
		t.Fatal("update.go calls loadProject directly; reloads must use loadCmd")
	}
}

// reducerSource finds the file declaring Update, so the assertion follows the
// reducer if it moves between files rather than passing vacuously.
func reducerSource(t *testing.T) (string, string) {
	t.Helper()
	for _, path := range packageSourceFiles(t) {
		content := readCode(t, path)
		if strings.Contains(content, "func (m Model) Update(") {
			return path, content
		}
	}
	t.Fatal("no source file declares func (m Model) Update")
	return "", ""
}

// TestBadgeVocabularyLivesInOneFile holds the single-mapping rule: the
// clearance, blocker, and freshness vocabularies are translated in badges.go
// and nowhere else, and nothing else in the package reaches for a badge accent.
// A second surface that wanted its own wording would have to name one of these,
// which is exactly the drift this asserts against.
//
// Task status and stage are deliberately not listed. card.go reads them to
// choose which gate decision governs a Task — the same choice data.ResolveNext
// makes — which is dispatch, not a second copy of the badge vocabulary.
func TestBadgeVocabularyLivesInOneFile(t *testing.T) {
	const mapping = "badges.go"
	reserved := []string{
		"data.ClearanceCurrent", "data.ClearanceNeedsWork", "data.ClearanceStale",
		"data.ClearanceUnknown", "data.ClearanceMissing",
		"data.GateBlock", "data.FreshnessCurrent", "data.FreshnessStale",
		"data.FreshnessUnknown", "styles.Badge",
	}

	for _, path := range packageSourceFiles(t) {
		if path == mapping {
			continue
		}
		content := readCode(t, path)
		for _, name := range reserved {
			if strings.Contains(content, name) {
				t.Errorf("%s references %q; that vocabulary is translated in %s alone", path, name, mapping)
			}
		}
	}
}

// TestRenderingResolvesNothing holds the other half of the derive-nothing rule
// structurally: the files that draw columns and the board may not call a
// resolver or reach into an index. Resolution happens once, in card.go, before
// anything is rendered.
func TestRenderingResolvesNothing(t *testing.T) {
	forbidden := []string{"data.Resolve", "LatestCheck", "ScopeChecks", ".Freshness"}

	for _, path := range []string{"column.go", "view.go", "detail_view.go"} {
		content := readCode(t, path)
		for _, name := range forbidden {
			if strings.Contains(content, name) {
				t.Errorf("%s references %q; rendering reads resolved values and resolves nothing", path, name)
			}
		}
	}
}

// TestDetailRenderingResolvesNothing is the same rule for the detail overlay,
// which is the reason its resolution and its rendering are two files: detail.go
// reaches the index and settles a RecordDetail, and detail_view.go may not
// reach anything at all. A renderer that could look a record up is a renderer
// that can disagree with the badge on the card behind it.
func TestDetailRenderingResolvesNothing(t *testing.T) {
	const view = "detail_view.go"
	forbidden := []string{
		"data.Resolve", "V2Index", "index.", "LatestCheck", "ScopeChecks",
		"Tasks[", "Objectives[", "Checks[", "Issues[", "CheckIssues", "ObjectiveTasks",
	}

	content := readCode(t, view)
	for _, name := range forbidden {
		if strings.Contains(content, name) {
			t.Errorf("%s references %q; the overlay formats a resolved RecordDetail and consults nothing", view, name)
		}
	}
}

// TestProjectionIsResolvedOnlyInTheLoadCommand holds the shared-interpretation
// rule structurally: data.ResolveNext is called in the load command and nowhere
// else, so the Next area formats one resolved value per load rather than
// re-deriving an answer of its own at render time. A second call site in a
// rendering path is the exact divergence E48 built the projection to prevent.
func TestProjectionIsResolvedOnlyInTheLoadCommand(t *testing.T) {
	const loadCommand = "load.go"

	for _, path := range packageSourceFiles(t) {
		if path == loadCommand {
			continue
		}
		if strings.Contains(readCode(t, path), "data.ResolveNext") {
			t.Errorf("%s calls data.ResolveNext; the projection is resolved in %s alone", path, loadCommand)
		}
	}
}

// TestNextPanelDerivesNothing is the derive-nothing rule for the Next area: the
// file that formats it may not reach the index, a record, a resolver, or the
// sidebar's selection. Everything it says comes from the one data.Next it was
// handed, which is what makes the board and `savepoint resume` the same answer
// rather than two answers that currently agree.
func TestNextPanelDerivesNothing(t *testing.T) {
	const panel = "next_panel.go"
	forbidden := []string{
		"data.Resolve", "State.Index", "index.", "Router", "LatestCheck", "ScopeChecks",
		"m.Cards", "SelectedObjective", "Objectives[", "Tasks[", "Checks[",
	}

	content := readCode(t, panel)
	for _, name := range forbidden {
		if strings.Contains(content, name) {
			t.Errorf("%s references %q; the Next area formats the resolved projection and consults nothing else", panel, name)
		}
	}
}

// TestEvidenceWordingHasOneSource holds the other half of that: the board says
// nothing about clearance, freshness, exceptions, owner waits, or dependencies
// in words of its own. internal/resume owns that vocabulary for every surface
// reporting a data.Next, so a second copy cannot drift out of step with the
// first (STYLE-07, STYLE-09).
func TestEvidenceWordingHasOneSource(t *testing.T) {
	phrases := []string{
		"Technical clearance", "Owner acceptance", "Allowed by exception",
		"freshness", "has ever been recorded", "Waiting on Task",
	}

	for _, path := range packageSourceFiles(t) {
		content := readCode(t, path)
		for _, phrase := range phrases {
			if strings.Contains(content, phrase) {
				t.Errorf("%s spells out %q; that wording is internal/resume's, and the board calls it rather than restating it", path, phrase)
			}
		}
	}
}

// packageSourceFiles lists the package's non-test Go files. Test files are
// excluded because this file names the forbidden identifiers in order to look
// for them.
func packageSourceFiles(t *testing.T) []string {
	t.Helper()
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("read package directory: %v", err)
	}
	var files []string
	for _, entry := range entries {
		name := entry.Name()
		if filepath.Ext(name) != ".go" || strings.HasSuffix(name, "_test.go") {
			continue
		}
		files = append(files, name)
	}
	if len(files) == 0 {
		t.Fatal("no source files found; the boundary assertions would pass vacuously")
	}
	return files
}

// readCode returns a source file's code with its comments removed, so a
// structural assertion is about what the package does rather than about a
// sentence explaining it.
func readCode(t *testing.T, path string) string {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
	if err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	var out strings.Builder
	if err := printer.Fprint(&out, fset, file); err != nil {
		t.Fatalf("print %s: %v", path, err)
	}
	return out.String()
}
