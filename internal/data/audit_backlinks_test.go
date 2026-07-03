package data

import "testing"

func backlinkFindings() []AuditFinding {
	return []AuditFinding{
		{ID: "F001", Tasks: []string{"E32-audit-register-tui-review/T005-work-item-finding-backlinks"}, Epics: []string{"E32-audit-register-tui-review"}},
		{ID: "F002", Tasks: []string{"T005"}, Epics: []string{"E32"}},
		{ID: "F003", Tasks: []string{"E31-audit-register-data-model/T001-other"}, Epics: []string{"E31-audit-register-data-model"}},
		{ID: "F004"},
	}
}

func idsOf(findings []AuditFinding) []string {
	ids := make([]string, len(findings))
	for i, f := range findings {
		ids[i] = f.ID
	}
	return ids
}

func equalIDs(got []AuditFinding, want ...string) bool {
	ids := idsOf(got)
	if len(ids) != len(want) {
		return false
	}
	for i := range ids {
		if ids[i] != want[i] {
			return false
		}
	}
	return true
}

func TestFindingsForTask_matchesFullAndShortRefs(t *testing.T) {
	got := FindingsForTask(backlinkFindings(), "E32-audit-register-tui-review/T005-work-item-finding-backlinks")
	if !equalIDs(got, "F001", "F002") {
		t.Errorf("FindingsForTask = %v, want [F001 F002]", idsOf(got))
	}
}

func TestFindingsForTask_preservesInputOrder(t *testing.T) {
	findings := []AuditFinding{
		{ID: "F009", Tasks: []string{"T005"}},
		{ID: "F002", Tasks: []string{"T005"}},
	}
	got := FindingsForTask(findings, "E32-audit-register-tui-review/T005-slug")
	if !equalIDs(got, "F009", "F002") {
		t.Errorf("FindingsForTask order = %v, want [F009 F002]", idsOf(got))
	}
}

func TestFindingsForTask_noMatchAcrossDifferentTask(t *testing.T) {
	got := FindingsForTask(backlinkFindings(), "E31-audit-register-data-model/T099-ghost")
	if len(got) != 0 {
		t.Errorf("FindingsForTask for unlinked task = %v, want none", idsOf(got))
	}
}

func TestFindingsForTask_emptyIDOrFindings(t *testing.T) {
	if got := FindingsForTask(backlinkFindings(), ""); got != nil {
		t.Errorf("FindingsForTask with empty id = %v, want nil", idsOf(got))
	}
	if got := FindingsForTask(nil, "E32/T005"); got != nil {
		t.Errorf("FindingsForTask with no findings = %v, want nil", idsOf(got))
	}
}

func TestFindingsForEpic_matchesFullAndShortRefs(t *testing.T) {
	got := FindingsForEpic(backlinkFindings(), "E32-audit-register-tui-review")
	if !equalIDs(got, "F001", "F002") {
		t.Errorf("FindingsForEpic = %v, want [F001 F002]", idsOf(got))
	}
}

func TestFindingsForEpic_noMatchForDifferentEpic(t *testing.T) {
	got := FindingsForEpic(backlinkFindings(), "E31-audit-register-data-model")
	if !equalIDs(got, "F003") {
		t.Errorf("FindingsForEpic = %v, want [F003]", idsOf(got))
	}
}

func TestFindingsForEpic_emptyIDOrFindings(t *testing.T) {
	if got := FindingsForEpic(backlinkFindings(), ""); got != nil {
		t.Errorf("FindingsForEpic with empty id = %v, want nil", idsOf(got))
	}
	if got := FindingsForEpic(nil, "E32"); got != nil {
		t.Errorf("FindingsForEpic with no findings = %v, want nil", idsOf(got))
	}
}
