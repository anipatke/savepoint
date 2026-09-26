// Package migrate's convert_issues.go renders the file content of V2 Issue
// records planned from V1 defects and audit findings by plan.go's
// planDefect/planFinding. It shares convert.go's three structural rules — no
// healing beyond what plan.go already applied, no fabricated evidence, no
// paraphrase of authored bodies — and adds a fourth specific to follow-up:
// no laundering. A `fixed` finding still awaits independent proof and can
// never render as `resolved`; a `duplicate` finding's resolution proves
// nothing and may only name a canonical Issue that itself converted; and no
// converted Issue ever names a Check, because migration produced it, not an
// evaluation.
//
// Like ConvertObjective/ConvertTask, ConvertIssue writes nothing: it renders
// content, validates it by decoding through the real strict V2 decoder, and
// returns it for a later apply step to write.
package migrate

import (
	"fmt"
	"strings"
	"time"

	"github.com/opencode/savepoint/internal/data"
	"gopkg.in/yaml.v3"
)

// defectKnownFrontmatterKeys and findingKnownFrontmatterKeys are the V1
// frontmatter keys convert_issues.go already interprets through
// data.Defect/data.AuditFinding. Anything else present on the source (for
// example a finding's reviewer_note) is preserved under legacy_fields
// instead of being silently dropped, exactly as convert.go does for
// Objectives and Tasks.
var defectKnownFrontmatterKeys = []string{
	"id", "release", "status", "severity", "introduced", "reference", "stage", "title",
}

var findingKnownFrontmatterKeys = []string{
	"id", "title", "status", "severity", "confidence", "source_auditor", "work_item",
	"releases", "epics", "tasks", "defects", "guardrail_ids", "locations", "duplicate_of",
	"deferral_reason", "waiver_reason", "verified_proof", "first_seen", "last_seen", "proof_needed",
}

// issueActorOutput, issueSourceOutput, issueResolutionOutput, and
// issueHistoryOutput mirror the yaml tags of internal/data's unexported
// issue frontmatter types exactly, so marshalled output decodes through the
// real strict decoder. None of them ever carries a check reference: a
// migrated Issue's origin, resolution, and history never borrow a Check's
// authority.
type issueActorOutput struct {
	Role    string `yaml:"role"`
	Session string `yaml:"session"`
}

type issueSourceOutput struct {
	Kind  string           `yaml:"kind"`
	Actor issueActorOutput `yaml:"actor"`
	At    string           `yaml:"at"`
}

type issueResolutionOutput struct {
	Disposition string           `yaml:"disposition"`
	Actor       issueActorOutput `yaml:"actor"`
	At          string           `yaml:"at"`
	Reason      string           `yaml:"reason,omitempty"`
}

type issueHistoryOutput struct {
	At    string           `yaml:"at"`
	Actor issueActorOutput `yaml:"actor"`
	Kind  string           `yaml:"kind"`
	Note  string           `yaml:"note,omitempty"`
}

type issueOutputFrontmatter struct {
	ID           string                 `yaml:"id"`
	Title        string                 `yaml:"title"`
	Type         string                 `yaml:"type"`
	Status       string                 `yaml:"status"`
	Source       issueSourceOutput      `yaml:"source"`
	Tasks        []string               `yaml:"tasks,omitempty"`
	GuardrailIDs []string               `yaml:"guardrail_ids,omitempty"`
	Severity     string                 `yaml:"severity,omitempty"`
	Resolution   *issueResolutionOutput `yaml:"resolution,omitempty"`
	DuplicateOf  string                 `yaml:"duplicate_of,omitempty"`
	History      []issueHistoryOutput   `yaml:"history,omitempty"`
	LegacyFields map[string]any         `yaml:"legacy_fields,omitempty"`
}

// findingIssueRule is one row of the disposition-to-Issue mapping table:
// which Issue status a finding status becomes, and — when the disposition
// requires a seeded history entry — which kind and which source field
// supplies its note. Expressed as data rather than as a chain of
// conditionals, per the epic's own implementation guidance.
type findingIssueRule struct {
	status      data.IssueStatus
	historyKind data.IssueHistoryKind // "" means no seeded history entry
	note        func(*data.AuditFinding) string
}

// findingDispositionRules covers every finding status that reaches
// ConvertIssue: open/triaged/mapped need no history entry, deferred and
// owner_decision each carry the shared deferral_reason field (matching the
// finding template's own "required for deferred / owner_decision" note),
// and in_progress/fixed both land as in_progress with a repair_attempted
// entry — fixed is never resolved, because it still awaits independent
// proof. verified and waived never reach this table: plan.go archives them
// outright. duplicate is handled separately in convertFindingIssue, since
// its shape (a resolution naming a canonical Issue) differs structurally
// from a seeded history entry.
var findingDispositionRules = map[data.FindingStatus]findingIssueRule{
	data.FindingOpen:    {status: data.IssueStatusOpen},
	data.FindingTriaged: {status: data.IssueStatusOpen},
	data.FindingMapped:  {status: data.IssueStatusOpen},
	data.FindingDeferred: {
		status: data.IssueStatusOpen, historyKind: data.IssueHistoryDeferred,
		note: func(f *data.AuditFinding) string { return f.DeferralReason },
	},
	data.FindingOwnerDecision: {
		status: data.IssueStatusOpen, historyKind: data.IssueHistoryOwnerDecision,
		note: func(f *data.AuditFinding) string { return f.DeferralReason },
	},
	data.FindingInProgress: {
		status: data.IssueStatusInProgress, historyKind: data.IssueHistoryRepairAttempted,
		note: func(f *data.AuditFinding) string { return f.ProofNeeded },
	},
	data.FindingFixed: {
		status: data.IssueStatusInProgress, historyKind: data.IssueHistoryRepairAttempted,
		note: func(f *data.AuditFinding) string { return f.ProofNeeded },
	},
}

// ConvertIssue renders the V2 Issue file content for target, reading the V1
// defect or finding source fresh from root and dispatching on which role the
// source classifies as. It writes nothing itself.
func ConvertIssue(root string, plan *ConversionPlan, target PlannedTarget) (string, error) {
	if target.Kind != TargetIssue {
		return "", fmt.Errorf("convert issue: target %s is not an issue target", target.GlobalID)
	}

	switch Classify(target.Legacy.Path) {
	case RoleDefect:
		return convertDefectIssue(root, plan, target)
	case RoleFinding:
		return convertFindingIssue(root, plan, target)
	default:
		return "", fmt.Errorf("convert issue %s: %s: source is neither a defect nor a finding", target.GlobalID, target.Legacy.Path)
	}
}

// convertDefectIssue renders an Issue from a V1 defect. plan.go only ever
// plans an Issue target for an open or in_progress defect (a resolved one is
// archived), but the status is re-derived and refused here rather than
// assumed, matching convert.go's own defensive re-checks.
func convertDefectIssue(root string, plan *ConversionPlan, target PlannedTarget) (string, error) {
	raw, err := readSourceFile(root, target.Legacy.Path)
	if err != nil {
		return "", err
	}
	crlf := strings.Contains(raw, "\r\n")
	fm, body, err := data.SplitFrontmatterBody(raw)
	if err != nil {
		return "", fmt.Errorf("convert issue %s: %s: %w", target.GlobalID, target.Legacy.Path, err)
	}

	defect, err := data.NewParser().ParseDefectFile(target.Legacy.Path, raw)
	if err != nil {
		return "", fmt.Errorf("convert issue %s: %s: %w", target.GlobalID, target.Legacy.Path, err)
	}

	var status data.IssueStatus
	switch defect.Status {
	case data.DefectOpen:
		status = data.IssueStatusOpen
	case data.DefectInProgress:
		status = data.IssueStatusInProgress
	default:
		return "", fmt.Errorf("%w: defect %s status %q does not convert to an Issue", ErrAmbiguousLifecycle, defect.ID, defect.Status)
	}

	legacyFields, err := unknownFields(fm, defectKnownFrontmatterKeys)
	if err != nil {
		return "", fmt.Errorf("convert issue %s: %s: %w", target.GlobalID, target.Legacy.Path, err)
	}

	var tasks []string
	if defect.Reference != "" {
		if id, ok := resolveTaskReference(plan.Targets, []string{defect.Release}, defect.Reference); ok {
			tasks = []string{id}
		}
	}

	actor := migrationActor(target.Legacy.Path)
	out := issueOutputFrontmatter{
		ID:           target.GlobalID,
		Title:        defect.Title,
		Type:         string(data.IssueTypeDefect),
		Status:       string(status),
		Source:       issueSourceOutput{Kind: string(data.IssueOriginMigration), Actor: actor, At: formatV2Time(plan.GeneratedAt)},
		Tasks:        tasks,
		Severity:     string(defect.Severity),
		LegacyFields: legacyFields,
	}

	renderedBody := relocatedBody("defect", target.Legacy.Path, defect.Release, body)
	return renderIssueContent(target, renderedBody, crlf, out)
}

// convertFindingIssue renders an Issue from a V1 audit finding, covering
// every disposition plan.go can plan a target for: open/triaged/mapped,
// deferred, owner_decision, in_progress/fixed, and duplicate (whose
// canonical converted). verified and waived findings are archived by
// plan.go and never reach here.
func convertFindingIssue(root string, plan *ConversionPlan, target PlannedTarget) (string, error) {
	raw, err := readSourceFile(root, target.Legacy.Path)
	if err != nil {
		return "", err
	}
	crlf := strings.Contains(raw, "\r\n")
	fm, body, err := data.SplitFrontmatterBody(raw)
	if err != nil {
		return "", fmt.Errorf("convert issue %s: %s: %w", target.GlobalID, target.Legacy.Path, err)
	}

	finding, err := data.NewParser().ParseFindingFile(target.Legacy.Path, raw)
	if err != nil {
		return "", fmt.Errorf("convert issue %s: %s: %w", target.GlobalID, target.Legacy.Path, err)
	}

	legacyFields, err := unknownFields(fm, findingKnownFrontmatterKeys)
	if err != nil {
		return "", fmt.Errorf("convert issue %s: %s: %w", target.GlobalID, target.Legacy.Path, err)
	}

	issueType := data.IssueTypeOther
	if len(finding.GuardrailIDs) > 0 {
		issueType = data.IssueTypeGuardrail
	}

	actor := migrationActor(target.Legacy.Path)
	out := issueOutputFrontmatter{
		ID:           target.GlobalID,
		Title:        finding.Title,
		Type:         string(issueType),
		Source:       issueSourceOutput{Kind: string(data.IssueOriginMigration), Actor: actor, At: formatV2Time(plan.GeneratedAt)},
		Tasks:        resolveReferencedTasks(plan.Targets, finding.Releases, finding.Tasks),
		GuardrailIDs: finding.GuardrailIDs,
		Severity:     string(finding.Severity),
		LegacyFields: legacyFields,
	}

	if target.DuplicateOfGlobalID != "" {
		at, _ := seedHistoryTime(finding.LastSeen, plan.GeneratedAt)
		out.Status = string(data.IssueStatusResolved)
		out.DuplicateOf = target.DuplicateOfGlobalID
		out.Resolution = &issueResolutionOutput{
			Disposition: string(data.IssueDispositionDuplicate),
			Actor:       actor,
			At:          formatV2Time(at),
			Reason:      fmt.Sprintf("migrated V1 finding %s recorded as a duplicate of %s; no independent proof is tracked here", finding.ID, finding.DuplicateOf),
		}
	} else {
		findingStatus := finding.Status
		if target.DecidedStatus != "" {
			findingStatus = data.FindingStatus(target.DecidedStatus)
		}
		rule, ok := findingDispositionRules[findingStatus]
		if !ok {
			return "", fmt.Errorf("%w: finding %s status %q does not convert to an Issue", ErrAmbiguousLifecycle, finding.ID, finding.Status)
		}
		out.Status = string(rule.status)
		if rule.historyKind != "" {
			at, dateNote := seedHistoryTime(finding.LastSeen, plan.GeneratedAt)
			out.History = []issueHistoryOutput{{
				At:    formatV2Time(at),
				Actor: actor,
				Kind:  string(rule.historyKind),
				Note:  historyNote(rule.note(finding), dateNote),
			}}
		}
	}

	renderedBody := relocatedBody("finding", target.Legacy.Path, strings.Join(finding.Releases, ", "), body)
	renderedBody += findingMetadataBody(finding)
	return renderIssueContent(target, renderedBody, crlf, out)
}

// findingMetadataBody renders the V1 finding frontmatter fields that convert
// into no Issue field of their own — proof_needed, first_seen/last_seen,
// confidence, source_auditor, and the work_item/releases/epics/tasks/
// defects/locations relations — as an authored section, so AC3's "reasons,
// proof requirements, and relations" survive conversion instead of being
// dropped simply because they are recognized (and therefore not caught by
// legacy_fields) but have no dedicated Issue field.
func findingMetadataBody(finding *data.AuditFinding) string {
	var b strings.Builder
	b.WriteString("\n\n## V1 Finding Metadata\n")
	writeMetaLine(&b, "Confidence", string(finding.Confidence))
	writeMetaLine(&b, "Source auditor", finding.SourceAuditor)
	writeMetaLine(&b, "Proof needed", finding.ProofNeeded)
	writeMetaLine(&b, "First seen", finding.FirstSeen)
	writeMetaLine(&b, "Last seen", finding.LastSeen)
	writeMetaLine(&b, "Work item", finding.WorkItem)
	writeMetaListLine(&b, "Releases", finding.Releases)
	writeMetaListLine(&b, "Epics", finding.Epics)
	writeMetaListLine(&b, "Tasks", finding.Tasks)
	writeMetaListLine(&b, "Defects", finding.Defects)
	writeMetaListLine(&b, "Locations", finding.Locations)
	return b.String()
}

func writeMetaLine(b *strings.Builder, label, value string) {
	if value == "" {
		return
	}
	fmt.Fprintf(b, "\n- %s: %s", label, value)
}

func writeMetaListLine(b *strings.Builder, label string, values []string) {
	if len(values) == 0 {
		return
	}
	fmt.Fprintf(b, "\n- %s: %s", label, strings.Join(values, ", "))
}

// renderIssueContent marshals out, appends renderedBody (already relocated
// under the standard "Migrated from V1" heading convert.go's Objective/Task
// rendering also uses), assembles the final file content, and validates it
// by decoding through the real strict decoder before returning it.
func renderIssueContent(target PlannedTarget, renderedBody string, crlf bool, out issueOutputFrontmatter) (string, error) {
	marshalled, err := yaml.Marshal(&out)
	if err != nil {
		return "", fmt.Errorf("convert issue %s: marshal yaml: %w", target.GlobalID, err)
	}

	content := "---\n" + strings.TrimSpace(string(marshalled)) + "\n---" + renderedBody
	if crlf {
		content = strings.ReplaceAll(content, "\n", "\r\n")
	}

	if _, err := data.DecodeIssueV2(target.TargetPath, content); err != nil {
		return "", fmt.Errorf("convert issue %s: rendered content did not decode: %w", target.GlobalID, err)
	}

	return content, nil
}

// migrationActor is the {role, session} provenance every migrated Issue's
// source, resolution, and history entries attribute their facts to. Session
// names the exact V1 source path — LegacyKey's own "actual uniqueness
// qualifier" — so the source-qualified legacy key behind kind: migration is
// recorded directly in the frontmatter, not left implicit.
func migrationActor(legacyPath string) issueActorOutput {
	return issueActorOutput{Role: string(data.ActorRoleOwner), Session: legacyPath}
}

// formatV2Time renders t as the RFC 3339 string decodeV2Timestamp requires.
func formatV2Time(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

// seedHistoryTime returns the date a seeded history entry should carry: the
// V1 source's own recorded date when raw parses as one, or fallback (the
// migration time) with an explanatory note when raw is absent or
// unparseable — never a fabricated date presented as if the source recorded
// it.
func seedHistoryTime(raw string, fallback time.Time) (time.Time, string) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return fallback, "the V1 source recorded no date for this entry; recorded here as the migration time instead"
	}
	t, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return fallback, fmt.Sprintf("the V1 source's recorded date %q did not parse; recorded here as the migration time instead", raw)
	}
	return t, ""
}

// historyNote joins a finding's own authored reason with an explanatory note
// about the date substitution, when one applies; either half may be empty.
func historyNote(reason, dateNote string) string {
	switch {
	case reason == "":
		return dateNote
	case dateNote == "":
		return reason
	default:
		return reason + " (" + dateNote + ")"
	}
}

// resolveReferencedTasks resolves every legacy task reference in refs to its
// converted global Task ID, scoped to releases, skipping (never guessing at)
// a reference that does not resolve to a converted active Task and
// deduplicating repeats. Order follows refs, since that is the order the V1
// source itself named them in.
func resolveReferencedTasks(targets []PlannedTarget, releases, refs []string) []string {
	var out []string
	seen := map[string]bool{}
	for _, ref := range refs {
		id, ok := resolveTaskReference(targets, releases, ref)
		if !ok || seen[id] {
			continue
		}
		seen[id] = true
		out = append(out, id)
	}
	return out
}
