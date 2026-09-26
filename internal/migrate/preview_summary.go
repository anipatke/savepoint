package migrate

import (
	"fmt"
	"sort"
	"strings"
)

// FormatSummaryPreview renders plan as the short report `savepoint migrate`
// prints by default (I-068): what the migration creates and archives, what it
// still needs from the owner, and whether it can be applied. Every decision
// still needed is listed in full with the decisions-file entry that answers
// it; archives and advisory notes are counted rather than listed.
// FormatPreview remains the complete listing behind --verbose.
func FormatSummaryPreview(plan *ConversionPlan) string {
	var b strings.Builder

	b.WriteString("Migration preview (nothing has been written)\n============================================\n\n")

	if plan.SchemaAlreadyV2 {
		b.WriteString("project already declares schema_version: 2; nothing to migrate.\n")
		return b.String()
	}

	if len(plan.Conflicts) > 0 {
		b.WriteString("Conflicts (resolve these before migrating)\n------------------------------------------\n")
		conflicts := append([]Conflict{}, plan.Conflicts...)
		sort.Slice(conflicts, func(i, j int) bool { return conflicts[i].Path < conflicts[j].Path })
		for _, c := range conflicts {
			fmt.Fprintf(&b, "- %s: %s\n", c.Path, c.Detail)
		}
		b.WriteString("\n")
	}

	writeSummaryCounts(&b, plan)
	writeGoalSelection(&b, plan.GoalSelection)
	writeSummaryObjectives(&b, plan.Targets)
	writeDecisionsNeeded(&b, plan.Ambiguities)
	writeDecisionsApplied(&b, plan.Ambiguities)

	switch {
	case len(plan.Conflicts) > 0:
		b.WriteString("Status: blocked by the conflicts above.\n")
	case !plan.Appliable:
		fmt.Fprintf(&b, "Status: blocked until %d decision(s) above are answered. Add them to a decisions file and run\n  savepoint migrate --decisions <file>\n", len(plan.UnresolvedBlockingIDs))
	default:
		b.WriteString("Status: ready. To apply, run the same command with --apply.\n")
	}
	b.WriteString("Run with --verbose to list every planned record, archived file, and note.\n")
	return b.String()
}

func writeSummaryCounts(b *strings.Builder, plan *ConversionPlan) {
	counts := map[TargetKind]int{}
	var goals []string
	for _, t := range plan.Targets {
		counts[t.Kind]++
		if t.Kind == TargetRelease {
			goals = append(goals, fmt.Sprintf("%s %s", t.GlobalID, strings.ReplaceAll(t.ReleaseStatus, "_", " ")))
		}
	}
	sort.Strings(goals)

	advisories := 0
	for _, a := range plan.Ambiguities {
		if !a.Blocking {
			advisories++
		}
	}

	b.WriteString("Will create\n-----------\n")
	fmt.Fprintf(b, "  Goals       %4d", counts[TargetRelease])
	if len(goals) > 0 {
		fmt.Fprintf(b, "  (%s)", strings.Join(goals, ", "))
	}
	b.WriteString("\n")
	fmt.Fprintf(b, "  Objectives  %4d\n", counts[TargetObjective])
	fmt.Fprintf(b, "  Tasks       %4d\n", counts[TargetTask])
	fmt.Fprintf(b, "  Issues      %4d\n", counts[TargetIssue])
	if len(plan.Documents) > 0 {
		docs := make([]string, 0, len(plan.Documents))
		for _, d := range plan.Documents {
			docs = append(docs, d.TargetPath)
		}
		sort.Strings(docs)
		fmt.Fprintf(b, "  Documents   %4d  (%s)\n", len(plan.Documents), strings.Join(docs, ", "))
	}
	b.WriteString("\n")
	fmt.Fprintf(b, "Will archive  %d V1 file(s), unchanged, under .savepoint/archive/v1/\n", len(plan.Archives))
	if len(plan.Prereqs) > 0 {
		fmt.Fprintf(b, "Legacy prerequisites  %d (open Tasks depending on archived, completed work)\n", len(plan.Prereqs))
	}
	if len(plan.WaivedRefs) > 0 {
		fmt.Fprintf(b, "Waived references  %d\n", len(plan.WaivedRefs))
	}
	if advisories > 0 {
		fmt.Fprintf(b, "Notes  %d advisory (files migration does not recognise; archived as-is, nothing to do)\n", advisories)
	}
	b.WriteString("\n")
}

func writeSummaryObjectives(b *strings.Builder, targets []PlannedTarget) {
	var objectives []PlannedTarget
	tasks := map[string]int{}
	for _, t := range targets {
		switch t.Kind {
		case TargetObjective:
			objectives = append(objectives, t)
		case TargetTask:
			tasks[objectiveDirOf(t.TargetPath)]++
		}
	}
	if len(objectives) == 0 {
		return
	}
	sort.Slice(objectives, func(i, j int) bool { return objectives[i].GlobalID < objectives[j].GlobalID })

	b.WriteString("Objectives\n----------\n")
	for _, o := range objectives {
		source := o.Legacy.Epic
		if o.Legacy.Release != "" {
			source = o.Legacy.Release + "/" + o.Legacy.Epic
		}
		fmt.Fprintf(b, "  %s  from %s, %d open task(s)", o.GlobalID, source, tasks[objectiveDirOf(o.TargetPath)])
		if o.DecidedStatus != "" {
			fmt.Fprintf(b, ", status %s by your decision", strings.ReplaceAll(o.DecidedStatus, "_", " "))
		}
		b.WriteString("\n")
	}
	b.WriteString("\n")
}

// objectiveDirOf is the objectives/O-###-slug directory a target's path sits
// under, which is how a Task target is matched to its Objective.
func objectiveDirOf(path string) string {
	parts := strings.SplitN(path, "/", 3)
	if len(parts) < 2 {
		return path
	}
	return parts[0] + "/" + parts[1]
}

func writeDecisionsNeeded(b *strings.Builder, ambiguities []Ambiguity) {
	var needed []Ambiguity
	for _, a := range ambiguities {
		if a.Blocking && !a.Resolved {
			needed = append(needed, a)
		}
	}
	if len(needed) == 0 {
		return
	}
	fmt.Fprintf(b, "Decisions needed (%d)\n--------------------\n", len(needed))
	for _, a := range needed {
		fmt.Fprintf(b, "- %s\n    %s\n", a.Path, a.Detail)
		if len(a.Choices) > 0 {
			fmt.Fprintf(b, "    choices: %s\n", strings.Join(a.Choices, ", "))
		}
	}
	b.WriteString("\nDecisions file entries to fill in (pick one choice for each value):\n\n  decisions:\n")
	for _, a := range needed {
		value := "<choice>"
		if len(a.Choices) > 0 {
			value = "<" + strings.Join(a.Choices, " | ") + ">"
		}
		fmt.Fprintf(b, "    - id: %s\n      value: %s\n", a.ID, value)
	}
	b.WriteString("\n")
}

func writeDecisionsApplied(b *strings.Builder, ambiguities []Ambiguity) {
	var applied []Ambiguity
	for _, a := range ambiguities {
		if a.Resolved {
			applied = append(applied, a)
		}
	}
	if len(applied) == 0 {
		return
	}
	fmt.Fprintf(b, "Decisions applied (%d)\n---------------------\n", len(applied))
	for _, a := range applied {
		fmt.Fprintf(b, "- %s -> %s  (from %s)\n", a.Path, a.Decision, a.DecisionSourceFile)
	}
	b.WriteString("\n")
}
