package data

import "testing"

func TestResolveConcurrencyV2_duplicateLaneSuppressesItsRecommendations(t *testing.T) {
	index := newConcurrencyIndex(planned("T-001", "core", []string{}, []string{"a.go"}), planned("T-002", "board", []string{}, []string{"b.go"}), planned("T-003", "extra", []string{}, []string{"c.go"}))
	objective := mustDecodeObjectiveV2(t, "O.md", "---\nid: O-001\ntitle: O\nstatus: in_progress\nlanes:\n  - {key: core, title: First}\n  - {key: core, title: Second}\n  - {key: board, title: Board}\n  - {key: extra, title: Extra}\n---\n")
	index.Objectives["O-001"].Plan = objective.Plan
	validateV2Planning(index)

	c := concurrencyFor(index)
	for _, group := range c.Groups {
		for _, id := range group {
			if id == "T-001" {
				t.Errorf("Groups = %v, want nothing recommended for the ambiguous lane", c.Groups)
			}
		}
	}
	if got := candidateIDs(c); len(got) != 2 || got[0] != "T-002" || got[1] != "T-003" {
		t.Errorf("candidates = %v, want only the unambiguous lanes", got)
	}
	if noteFor(c, "T-001", ConcurrencyUnusablePlan) == nil {
		t.Errorf("T-001 has no unusable-plan note: %+v", c.Notes)
	}
	if len(c.Diagnostics) == 0 {
		t.Error("duplicate-lane diagnostic was lost")
	}
}
