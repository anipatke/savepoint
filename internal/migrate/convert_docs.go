// Package migrate's convert_docs.go renders the small, fixed set of
// top-level project documents plan.go's planDocuments assigns a fate to:
// PRD.md's byte-preserving move to Idea.md, and router.md's in-place state
// anchor rewrite. Health-Check.md needs no renderer here — it is archived
// byte-for-byte like any other ArchiveEntry — but this file also owns
// extracting its candidate commands for a future preview to report, since
// that extraction is a document-content concern; no caller ever feeds the
// result into config.yml, so no quality_gates entry is ever inferred from
// prose.
package migrate

import (
	"fmt"
	"strings"

	"github.com/opencode/savepoint/internal/data"
	"gopkg.in/yaml.v3"
)

// ConvertIdea renders the V2 Idea.md content for doc: the exact V1 PRD.md
// bytes, unchanged. This is a relocation, not a re-render — no section is
// added, reordered, reworded, or dropped. plan.go's separate ArchiveEntry
// for the same source path preserves the identical bytes a second time,
// under archive/v1/, so V1 authorship survives both live and as history.
func ConvertIdea(root string, doc PlannedDocument) (string, error) {
	if doc.Kind != DocumentIdea {
		return "", fmt.Errorf("convert idea: doc %q is not an idea document", doc.SourcePath)
	}
	return readSourceFile(root, doc.SourcePath)
}

// routerStateMap is the fixed V1-to-V2 router state vocabulary: this is the
// one place it is written.
var routerStateMap = map[string]string{
	"pre-implementation":  "design",
	"epic-design":         "design",
	"epic-task-breakdown": "design",
	"task-building":       "task",
	"defect-building":     "task",
	"audit-pending":       "check",
}

// routerV2State is the "## Current state" YAML anchor's V2 shape: state,
// optional objective, optional task, and the preserved next_action prose.
// There is no V2 release/epic/defect field, so an archived selection is
// explained in an appended note rather than a fifth anchor field the design
// does not define.
type routerV2State struct {
	State      string `yaml:"state"`
	Objective  string `yaml:"objective,omitempty"`
	Task       string `yaml:"task,omitempty"`
	NextAction string `yaml:"next_action"`
}

// ConvertRouter renders the V2 router.md content for doc, reading the V1
// router fresh from root. Only the "## Current state" YAML anchor changes;
// every other authored byte — the read-order prose, headings, everything —
// is preserved verbatim via data.ReplaceStateBlock. plan resolves the
// previously selected epic/task to their allocated O###/T###, or names the
// archive path that explains why there is nothing left to select.
func ConvertRouter(root string, plan *ConversionPlan, doc PlannedDocument) (string, error) {
	if doc.Kind != DocumentRouter {
		return "", fmt.Errorf("convert router: doc %q is not a router document", doc.SourcePath)
	}

	raw, err := readSourceFile(root, doc.SourcePath)
	if err != nil {
		return "", err
	}
	crlf := strings.Contains(raw, "\r\n")

	v1State, err := data.NewRouterReader().ReadState(raw)
	if err != nil {
		return "", fmt.Errorf("convert router %s: %w", doc.SourcePath, err)
	}

	v2State, note, err := mapRouterState(plan, *v1State)
	if err != nil {
		return "", fmt.Errorf("convert router %s: %w", doc.SourcePath, err)
	}

	marshalled, err := yaml.Marshal(&v2State)
	if err != nil {
		return "", fmt.Errorf("convert router %s: marshal yaml: %w", doc.SourcePath, err)
	}

	content, err := data.ReplaceStateBlock(raw, string(marshalled))
	if err != nil {
		return "", fmt.Errorf("convert router %s: %w", doc.SourcePath, err)
	}

	if note != "" {
		content = strings.TrimRight(content, "\n") + "\n\n## Migration Note\n\n" + note + "\n"
	}

	if crlf {
		content = strings.ReplaceAll(content, "\n", "\r\n")
	}

	return content, nil
}

// mapRouterState maps a V1 RouterState to its V2 anchor plus an optional
// migration note. Objective/Task selection is resolved against plan's own
// targets and archives — the same identity allocation the rest of migration
// uses — never re-derived from the raw strings a second way.
func mapRouterState(plan *ConversionPlan, v1 data.RouterState) (routerV2State, string, error) {
	v2, ok := routerStateMap[v1.State]
	if !ok {
		return routerV2State{}, "", fmt.Errorf("%w: router state %q", ErrAmbiguousLifecycle, v1.State)
	}

	out := routerV2State{State: v2, NextAction: v1.NextAction}
	var notes []string

	if v1.Epic != "" {
		objectiveID, archivePath := resolveRouterObjective(plan, v1.Release, v1.Epic)
		switch {
		case objectiveID != "":
			out.Objective = objectiveID
		case archivePath != "":
			notes = append(notes, fmt.Sprintf(
				"The previously selected epic %s/%s was archived rather than converted (see `%s`); the router now selects no Objective.",
				v1.Release, v1.Epic, archivePath))
		default:
			notes = append(notes, fmt.Sprintf(
				"The previously selected epic %s/%s does not resolve to any planned or archived record; the router now selects no Objective.",
				v1.Release, v1.Epic))
		}
	}

	if v1.Task != "" {
		taskID, archivePath := resolveRouterTask(plan, v1.Release, v1.Task)
		switch {
		case taskID != "":
			out.Task = taskID
		case archivePath != "":
			notes = append(notes, fmt.Sprintf(
				"The previously selected task %s was archived rather than converted (see `%s`); the router now selects no Task.",
				v1.Task, archivePath))
		default:
			notes = append(notes, fmt.Sprintf(
				"The previously selected task %s does not resolve to any planned or archived record; the router now selects no Task.",
				v1.Task))
		}
	}

	return out, strings.Join(notes, "\n\n"), nil
}

// resolveRouterObjective resolves a V1 router epic selection to the
// Objective's allocated O### when it converted, or the archive path that
// explains why it did not.
func resolveRouterObjective(plan *ConversionPlan, release, epic string) (globalID, archivePath string) {
	for _, t := range plan.Targets {
		if t.Kind == TargetObjective && t.Legacy.Release == release && t.Legacy.Epic == epic {
			return t.GlobalID, ""
		}
	}
	for _, a := range plan.Archives {
		if a.Role == RoleEpicDetail && a.Legacy != nil && a.Legacy.Release == release && a.Legacy.Epic == epic {
			return "", a.ArchivePath
		}
	}
	return "", ""
}

// resolveRouterTask resolves a V1 router task selection (an epic-qualified
// short id, e.g. "E01-example/T002-follow-up" — the same shape
// LegacyKey.OriginalID already carries for a task) to its allocated T###
// when it converted, or the archive path that explains why it did not.
func resolveRouterTask(plan *ConversionPlan, release, taskRef string) (globalID, archivePath string) {
	for _, t := range plan.Targets {
		if t.Kind == TargetTask && t.Legacy.Release == release && t.Legacy.OriginalID == taskRef {
			return t.GlobalID, ""
		}
	}
	for _, a := range plan.Archives {
		if a.Role == RoleTask && a.Legacy != nil && a.Legacy.Release == release && a.Legacy.OriginalID == taskRef {
			return "", a.ArchivePath
		}
	}
	return "", ""
}

// codeFence is the Markdown fenced-code-block delimiter
// CandidateHealthCheckCommands uses to find candidate command lines.
const codeFence = "```"

// CandidateHealthCheckCommands scans an archived Health-Check.md body for
// fenced-code-block lines that look like commands, so a future preview can
// report them to the owner as candidate quality_gates. It never feeds
// config.yml itself: parsing authored prose into executable configuration is
// exactly the guess migration exists to refuse, so this function's result is
// only ever attached to the ArchiveEntry for display.
func CandidateHealthCheckCommands(content string) []string {
	var commands []string
	inFence := false
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(strings.TrimRight(line, "\r"))
		if strings.HasPrefix(trimmed, codeFence) {
			inFence = !inFence
			continue
		}
		if !inFence || trimmed == "" {
			continue
		}
		commands = append(commands, trimmed)
	}
	return commands
}
