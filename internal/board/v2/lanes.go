package v2

import (
	"slices"

	"github.com/opencode/savepoint/internal/data"
)

// LaneHeading is the saved lane a card sits under, resolved while cards are
// grouped so rendering stays pure. Key identifies the lane across columns,
// namespaced by Objective, and an empty Key means the card has no heading.
// Headings are labels over cards: they are never selectable, never counted,
// and never part of a status write.
type LaneHeading struct {
	Key   string
	Title string
	// Readiness is the canonical wording for the lane's state, taken from
	// data.ResolveConcurrencyV2. It is empty when that projection withholds
	// advice, so a heading never states more than the projection did.
	Readiness string
}

// Present reports whether the card has a heading.
func (h LaneHeading) Present() bool { return h.Key != "" }

// ungroupedTitle heads the sequential work of an Objective that also has
// laned Tasks. It carries no readiness, since no lane is suggested for it.
const ungroupedTitle = "Ungrouped — sequential work"

// laneLayout orders taskIDs so each Objective's Tasks sit together in its
// declared lane order, and names the heading every Task falls under. The
// order is one sort applied before Tasks are split into columns, so lane order
// and Task order agree across all three. When no Objective in view has a laned
// Task it returns taskIDs untouched and no headings, which keeps the legacy
// board exactly as it was.
//
// namespaced prefixes each title with its Objective, for a view that spans
// Objectives. Readiness is always read per Objective: no suggestion ever
// compares lanes of different Objectives.
func laneLayout(index *data.V2Index, taskIDs []string, namespaced bool) ([]string, map[string]LaneHeading) {
	if index == nil {
		return taskIDs, nil
	}

	type placed struct {
		id        string
		objective string
		rank      int
		position  int
	}
	headings := map[string]LaneHeading{}
	projections := map[string]*data.ConcurrencyV2{}
	laned := map[string]bool{}
	for _, id := range taskIDs {
		task := index.Tasks[id]
		if _, ok := index.Objectives[task.Objective].Plan.Lane(task.Plan.Lane); ok && task.Plan.Lane != "" {
			laned[task.Objective] = true
		}
	}
	if len(laned) == 0 {
		return taskIDs, nil
	}

	order := make([]placed, 0, len(taskIDs))
	for position, id := range taskIDs {
		task := index.Tasks[id]
		entry := placed{id: id, objective: task.Objective, position: position}
		if !laned[task.Objective] {
			order = append(order, entry)
			continue
		}
		objective := index.Objectives[task.Objective]
		entry.rank = len(objective.Plan.Lanes)
		heading := LaneHeading{Key: task.Objective + "/", Title: ungroupedTitle}
		if lane, ok := objective.Plan.Lane(task.Plan.Lane); ok && task.Plan.Lane != "" {
			entry.rank = slices.IndexFunc(objective.Plan.Lanes, func(l data.ObjectiveLaneV2) bool { return l.Key == lane.Key })
			projection, seen := projections[task.Objective]
			if !seen {
				projection = data.ResolveConcurrencyV2(index, task.Objective, data.ConcurrencyOptionsV2{Enabled: true})
				projections[task.Objective] = projection
			}
			heading = LaneHeading{Key: task.Objective + "/" + lane.Key, Title: lane.Title, Readiness: laneReadiness(projection, lane.Key)}
		}
		if namespaced {
			heading.Title = task.Objective + " · " + heading.Title
		}
		headings[id] = heading
		order = append(order, entry)
	}

	slices.SortStableFunc(order, func(a, b placed) int {
		switch {
		case a.objective != b.objective:
			return compareStrings(a.objective, b.objective)
		case a.rank != b.rank:
			return a.rank - b.rank
		}
		return a.position - b.position
	})
	ids := make([]string, len(order))
	for i, entry := range order {
		ids[i] = entry.id
	}
	return ids, headings
}

func compareStrings(a, b string) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

// laneReadiness words the state of one lane from the canonical projection. It
// reads what the projection recorded and decides nothing itself.
func laneReadiness(c *data.ConcurrencyV2, key string) string {
	if c == nil || c.Withheld != nil {
		return ""
	}
	var lane data.ConcurrencyLaneV2
	for _, candidate := range c.Lanes {
		if candidate.Key == key {
			lane = candidate
		}
	}
	switch {
	case len(lane.InProgress) > 0:
		return "in progress"
	case slices.ContainsFunc(c.Candidates, func(candidate data.ConcurrencyCandidateV2) bool { return candidate.Lane == key }):
		return "ready to start"
	case len(lane.Planned) == 0:
		if len(lane.Done) > 0 {
			return "done"
		}
		return ""
	}
	for _, note := range c.Notes {
		if len(note.Tasks) == 0 || !slices.Contains(lane.Planned, note.Tasks[0]) {
			continue
		}
		switch note.Kind {
		case data.ConcurrencyStartBlocked:
			return "waiting on earlier work"
		case data.ConcurrencyUnknownScope, data.ConcurrencyPathAlias, data.ConcurrencyUnusablePlan:
			return "scope unclear"
		}
	}
	return "not suggested now"
}
