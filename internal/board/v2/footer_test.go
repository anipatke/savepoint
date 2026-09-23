package v2

import (
	"strings"
	"testing"
)

// I-025 shrank the always-visible footer to the keys that actually do
// something on the focused surface right now. The full key map, including
// every key omitted here, stays available in Help (see actions_test.go's
// TestHelpListsOnlyFocusedOwnerCapabilities).

func TestFooterOmitsReleaseHintWithoutARelease(t *testing.T) {
	model := openSizedBoard(t, writeValidProject(t), 130, 40)
	if strings.Contains(model.hints(), "r:releases") {
		t.Errorf("footer offers the Release selector with no Release in the project:\n%q", model.hints())
	}

	withRelease := releaseBoard(t)
	if !strings.Contains(withRelease.hints(), "r:releases") {
		t.Errorf("footer omits the Release selector when a Release exists:\n%q", withRelease.hints())
	}
}

func TestFooterOmitsClearHintWithNothingSelected(t *testing.T) {
	inSidebar := press(t, sidebarBoard(t, writeNavigationProject(t)), "left")
	if !strings.Contains(inSidebar.hints(), "esc:clear") {
		t.Fatalf("footer omits esc:clear with an Objective selected:\n%q", inSidebar.hints())
	}

	cleared := press(t, inSidebar, "esc")
	if strings.Contains(cleared.hints(), "esc:clear") {
		t.Errorf("footer offers esc:clear with nothing selected:\n%q", cleared.hints())
	}
}

func TestFooterOmitsDetailHintOnAnEmptySurface(t *testing.T) {
	root := writeEmptyProjectFromTemplate(t)

	columns := sidebarBoard(t, root)
	if strings.Contains(columns.hints(), "enter:detail") {
		t.Errorf("footer offers enter:detail over an empty project:\n%q", columns.hints())
	}

	sidebar := press(t, columns, "left")
	if strings.Contains(sidebar.hints(), "v:detail") {
		t.Errorf("footer offers v:detail on an empty sidebar:\n%q", sidebar.hints())
	}
}

func TestFooterOmitsDeadTaskIssuesHint(t *testing.T) {
	model := press(t, issueBoard(t, writeIssuesProject(t)), "i")
	if strings.Contains(model.hints(), "I:task issues") {
		t.Errorf("Issues overlay footer still advertises the unhandled I key:\n%q", model.hints())
	}
}

func TestFooterCanonicalHintOnlyWhenTheIssueHasADuplicateTarget(t *testing.T) {
	root := writeIssuesProject(t)

	// I-001 is Open and carries no duplicate_of.
	notDuplicate := press(t, issueBoard(t, root), "i", "enter")
	if strings.Contains(notDuplicate.hints(), "enter:canonical") {
		t.Errorf("footer offers enter:canonical for an Issue with no canonical target:\n%q", notDuplicate.hints())
	}

	// I-004 is Resolved's second row and names I-001 as duplicate_of.
	duplicate := press(t, issueBoard(t, root), "i", "right", "right", "down", "enter")
	if !strings.Contains(duplicate.hints(), "enter:canonical") {
		t.Errorf("footer omits enter:canonical for an Issue that names a canonical target:\n%q", duplicate.hints())
	}
}

// Help closes on esc or q; neither reaches the global quit while Help is
// open, so the footer must not claim q as a quit key there.
func TestHelpCloseKeysDoNotQuit(t *testing.T) {
	model := press(t, openSizedBoard(t, writeValidProject(t), 100, 30), "?")
	if !model.Help {
		t.Fatal("? did not open help")
	}
	if strings.Contains(model.hints(), "quit") {
		t.Errorf("help footer claims a quit key while esc/q only close help:\n%q", model.hints())
	}

	closedModel, cmd := model.Update(keyMsg("q"))
	closed := closedModel.(Model)
	if closed.Help {
		t.Error("q did not close help")
	}
	if cmd != nil {
		t.Error("q while help is open ran a command instead of only closing help")
	}
}

// Tab no longer focuses a board surface; arrow keys alone do. No footer
// context may advertise it.
func TestNoFooterContextAdvertisesTab(t *testing.T) {
	contexts := []Model{
		openSizedBoard(t, writeValidProject(t), 130, 40),
		press(t, sidebarBoard(t, writeNavigationProject(t)), "left"),
		openSizedBoard(t, writeValidProject(t), sidebarBreakpoint-1, 40),
		press(t, openSizedBoard(t, writeValidProject(t), 130, 40), "enter"),
	}
	for _, model := range contexts {
		if strings.Contains(model.hints(), "tab") {
			t.Errorf("footer still advertises the retired Tab key:\n%q", model.hints())
		}
	}
}
