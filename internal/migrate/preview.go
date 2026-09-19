// Package migrate's preview.go renders a ConversionPlan as the complete,
// deterministic, non-interactive report `savepoint migrate`'s default
// preview prints. It reads only the plan value passed to it: no filesystem
// access, so it never observes anything Plan itself did not already record.
package migrate

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// FormatPreview renders plan as the plan itself, not a summary that omits
// entries: every planned record and its identity mapping, every archive,
// document, legacy prerequisite, waived reference, conflict, and ambiguity
// Plan computed, in a stable order independent of Plan's internal build
// order.
func FormatPreview(plan *ConversionPlan) string {
	var b strings.Builder

	b.WriteString("Migration preview\n==================\n\n")

	if plan.SchemaAlreadyV2 {
		b.WriteString("project already declares schema_version: 2; nothing to migrate.\n")
		return b.String()
	}

	fmt.Fprintf(&b, "operation: %s\n", plan.OperationID)
	fmt.Fprintf(&b, "generated: %s\n\n", plan.GeneratedAt.Format(time.RFC3339))

	if len(plan.Conflicts) > 0 {
		b.WriteString("Conflicts\n---------\n")
		conflicts := append([]Conflict{}, plan.Conflicts...)
		sort.Slice(conflicts, func(i, j int) bool { return conflicts[i].Path < conflicts[j].Path })
		for _, c := range conflicts {
			fmt.Fprintf(&b, "- [%s] %s: %s\n", c.Kind, c.Path, c.Detail)
		}
		b.WriteString("\n")
	}

	writePlannedRecords(&b, plan.Targets)
	writeDocuments(&b, plan.Documents)
	writeArchives(&b, plan.Archives)
	writePrereqs(&b, plan.Prereqs)
	writeWaivedRefs(&b, plan.WaivedRefs)
	writeAmbiguities(&b, plan.Ambiguities)

	if !plan.Appliable {
		fmt.Fprintf(&b, "BLOCKED: unresolved blocking ambiguities: %s\n", strings.Join(plan.UnresolvedBlockingIDs, ", "))
	}

	return b.String()
}

func writePlannedRecords(b *strings.Builder, targets []PlannedTarget) {
	b.WriteString("Planned records\n---------------\n")
	if len(targets) == 0 {
		b.WriteString("(none)\n\n")
		return
	}

	sorted := append([]PlannedTarget{}, targets...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].GlobalID < sorted[j].GlobalID })

	for _, t := range sorted {
		legacy := t.Legacy
		scope := legacy.Path
		if legacy.Release != "" {
			scope = legacy.Release
			if legacy.Epic != "" {
				scope += "/" + legacy.Epic
			}
			scope += ": " + legacy.Path
		}
		fmt.Fprintf(b, "- %s %s <- %s -> %s\n", t.Kind, t.GlobalID, scope, t.InstallPath())
		if len(t.DependsOn) > 0 {
			fmt.Fprintf(b, "    depends_on: %s\n", strings.Join(t.DependsOn, ", "))
		}
		if t.Kind == TargetRelease {
			fmt.Fprintf(b, "    status: %s\n", t.ReleaseStatus)
			if t.ReleaseStatus == "done" {
				fmt.Fprintf(b, "    completion: historical archive reference at %s\n", archivePathFor(t.Legacy.Path))
			}
		}
		if t.DuplicateOfGlobalID != "" {
			fmt.Fprintf(b, "    duplicate_of: %s\n", t.DuplicateOfGlobalID)
		}
	}
	b.WriteString("\n")
}

func writeDocuments(b *strings.Builder, docs []PlannedDocument) {
	b.WriteString("Documents\n---------\n")
	if len(docs) == 0 {
		b.WriteString("(none)\n\n")
		return
	}

	sorted := append([]PlannedDocument{}, docs...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].SourcePath < sorted[j].SourcePath })

	for _, d := range sorted {
		fmt.Fprintf(b, "- %s: %s -> %s\n", d.Kind, d.SourcePath, d.TargetPath)
	}
	b.WriteString("\n")
}

func writeArchives(b *strings.Builder, archives []ArchiveEntry) {
	b.WriteString("Archives\n--------\n")
	if len(archives) == 0 {
		b.WriteString("(none)\n\n")
		return
	}

	sorted := append([]ArchiveEntry{}, archives...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].SourcePath < sorted[j].SourcePath })

	for _, a := range sorted {
		fmt.Fprintf(b, "- %s -> %s (%s)\n", a.SourcePath, a.ArchivePath, string(a.Role))
		if len(a.CandidateCommands) > 0 {
			fmt.Fprintf(b, "    candidate quality-gate commands found in the archived Health-Check.md, "+
				"for owner review only (none are added to config.yml automatically): %s\n",
				strings.Join(a.CandidateCommands, "; "))
		}
	}
	b.WriteString("\n")
}

func writePrereqs(b *strings.Builder, prereqs []LegacyPrerequisite) {
	if len(prereqs) == 0 {
		return
	}
	b.WriteString("Legacy prerequisites\n---------------------\n")
	sorted := append([]LegacyPrerequisite{}, prereqs...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Task < sorted[j].Task })
	for _, p := range sorted {
		fmt.Fprintf(b, "- %s depends on archived %s: %s\n", p.Task, p.ArchivePath, p.Evidence)
	}
	b.WriteString("\n")
}

func writeWaivedRefs(b *strings.Builder, refs []WaivedReference) {
	if len(refs) == 0 {
		return
	}
	b.WriteString("Waived references\n------------------\n")
	// plan.build() already sorts WaivedRefs by Task then ArchivePath.
	for _, w := range refs {
		fmt.Fprintf(b, "- %s references a waived finding archived at %s: %s\n", w.Task, w.ArchivePath, w.Reason)
	}
	b.WriteString("\n")
}

func writeAmbiguities(b *strings.Builder, ambiguities []Ambiguity) {
	b.WriteString("Ambiguities and decisions required\n-----------------------------------\n")
	if len(ambiguities) == 0 {
		b.WriteString("(none)\n\n")
		return
	}

	// plan.build() already sorts Ambiguities by ID.
	for _, a := range ambiguities {
		kind := "advisory"
		if a.Blocking {
			kind = "blocking"
		}
		status := "unresolved"
		if a.Resolved {
			status = fmt.Sprintf("resolved -> %s (from %s)", a.Decision, a.DecisionSourceFile)
		}
		fmt.Fprintf(b, "- [%s] %s %s: %s [%s]\n", a.ID, kind, a.Path, a.Detail, status)
		if len(a.Choices) > 0 {
			fmt.Fprintf(b, "    choices: %s\n", strings.Join(a.Choices, ", "))
		}
	}
	b.WriteString("\n")
}
