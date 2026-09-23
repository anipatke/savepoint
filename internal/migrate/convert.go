// Package migrate's convert.go renders the file content of V1-to-V2
// Objective and Task records planned by plan.go. It is the boundary where a
// well-meaning converter would quietly lie, so three rules are structural
// here rather than merely followed:
//
//   - No healing. Status and stage recognition is re-derived from raw
//     frontmatter through the same non-healing helpers plan.go uses
//     (resolveEpicStatus, resolveTaskStatus) plus a local, equally
//     non-healing stage check; an unrecognized value returns
//     ErrAmbiguousLifecycle rather than a defaulted one. ParseTaskFile's
//     healing lifecycle path (ParseTaskLifecycle) is never called from here.
//   - No fabricated evidence. Rendered frontmatter never carries last_check,
//     freshness, owner_validation, exception, or replan: those keys are
//     simply never written, so an absent V1 fact stays absent.
//   - No paraphrase. The entire V1 body — every heading, comment, and
//     multiline block exactly as authored — is relocated byte-for-byte under
//     one clearly labeled heading naming its V1 source path.
//
// Convert*'s output is plain file content, validated by decoding it back
// through the real strict V2 decoder before it is returned. It writes
// nothing; T008/T009 own turning this content into files on disk.
package migrate

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/opencode/savepoint/internal/data"
	"gopkg.in/yaml.v3"
)

// ErrAmbiguousLifecycle means an active V1 record carries a status or stage
// value convert.go does not recognize. Rendering refuses rather than
// defaulting the value, matching plan.go's AmbiguityUnrecognizedLifecycle
// contract for the one lifecycle dimension plan.go itself does not check:
// Task stage. (Status is already screened out before a target ever reaches
// convert.go, since planTasksImpl only allocates a target for a recognized
// status; this error exists so that guarantee is enforced here too, rather
// than assumed.)
var ErrAmbiguousLifecycle = errors.New("migrate: active record has an unrecognized lifecycle value; owner decision required")

// ErrUnresolvedDependency means a declared V1 epic dependency does not name
// any epic converting to an Objective in this plan.
var ErrUnresolvedDependency = errors.New("migrate: dependency reference does not resolve to a converted record")

// objectiveSourceFrontmatter is the set of V1 epic-detail frontmatter fields
// convert.go interprets. Any other key present on the source is preserved
// under legacy_fields rather than silently dropped.
type objectiveSourceFrontmatter struct {
	Type      string   `yaml:"type"`
	Status    string   `yaml:"status"`
	DependsOn []string `yaml:"depends_on"`
}

var objectiveKnownFrontmatterKeys = []string{"type", "status", "depends_on"}

// taskSourceFrontmatter mirrors parser.go's unexported taskFrontmatter
// field-for-field: this is the set of V1 task frontmatter keys convert.go
// already understands. Anything else present on the source is preserved
// under legacy_fields.
type taskSourceFrontmatter struct {
	ID               string   `yaml:"id"`
	Title            string   `yaml:"title"`
	Objective        string   `yaml:"objective"`
	Description      string   `yaml:"description"`
	Epic             string   `yaml:"epic"`
	Release          string   `yaml:"release"`
	Status           string   `yaml:"status"`
	Column           string   `yaml:"column"`
	Phase            string   `yaml:"phase"`
	Stage            string   `yaml:"stage"`
	Priority         string   `yaml:"priority"`
	Points           int      `yaml:"points"`
	Tags             []string `yaml:"tags"`
	Acceptance       []string `yaml:"acceptance"`
	Notes            string   `yaml:"notes"`
	DependsOn        []string `yaml:"depends_on"`
	Progress         any      `yaml:"progress"`
	ComplexityTier   string   `yaml:"complexity_tier"`
	ComplexityReason string   `yaml:"complexity_reason"`
}

var taskKnownFrontmatterKeys = []string{
	"id", "title", "objective", "description", "epic", "release",
	"status", "column", "phase", "stage", "priority", "points", "tags",
	"acceptance", "notes", "depends_on", "progress", "complexity_tier",
	"complexity_reason",
}

// objectiveOutputFrontmatter and taskOutputFrontmatter mirror the yaml tags
// of internal/data's objectiveV2Frontmatter and taskV2Frontmatter exactly, so
// marshalled output decodes through the real strict decoders. legacy_fields
// is an extra key those decoders simply ignore.
type objectiveOutputFrontmatter struct {
	ID           string         `yaml:"id"`
	Title        string         `yaml:"title"`
	Status       string         `yaml:"status"`
	DependsOn    []string       `yaml:"depends_on,omitempty"`
	Release      string         `yaml:"release,omitempty"`
	LegacyFields map[string]any `yaml:"legacy_fields,omitempty"`
}

type taskDependsOnOutput struct {
	Task     string `yaml:"task"`
	Requires string `yaml:"requires,omitempty"`
}

type taskActorOutput struct {
	Role    string `yaml:"role"`
	Session string `yaml:"session"`
}

type taskOutputFrontmatter struct {
	ID           string                `yaml:"id"`
	Title        string                `yaml:"title"`
	Objective    string                `yaml:"objective"`
	PlannedBy    taskActorOutput       `yaml:"planned_by"`
	Status       string                `yaml:"status"`
	Stage        string                `yaml:"stage,omitempty"`
	DependsOn    []taskDependsOnOutput `yaml:"depends_on,omitempty"`
	Release      string                `yaml:"release,omitempty"`
	LegacyFields map[string]any        `yaml:"legacy_fields,omitempty"`
}

// ConvertObjective renders the V2 Objective file content for target, reading
// the V1 epic detail source fresh from root. It returns file content ready
// for a later write; it writes nothing itself. plan supplies the sibling
// Objective targets needed to resolve a declared epic dependency to its
// allocated O-###.
func ConvertObjective(root string, plan *ConversionPlan, target PlannedTarget) (string, error) {
	if target.Kind != TargetObjective {
		return "", fmt.Errorf("convert objective: target %s is not an objective target", target.GlobalID)
	}

	raw, err := readSourceFile(root, target.Legacy.Path)
	if err != nil {
		return "", err
	}
	crlf := strings.Contains(raw, "\r\n")
	fm, body, err := data.SplitFrontmatterBody(raw)
	if err != nil {
		return "", fmt.Errorf("convert objective %s: %s: %w", target.GlobalID, target.Legacy.Path, err)
	}

	var source objectiveSourceFrontmatter
	if err := yaml.Unmarshal([]byte(fm), &source); err != nil {
		return "", fmt.Errorf("convert objective %s: %s: %w", target.GlobalID, target.Legacy.Path, err)
	}

	status, _, recognized := resolveEpicStatus(source.Status)
	if !recognized {
		return "", fmt.Errorf("%w: epic %s status %q", ErrAmbiguousLifecycle, target.Legacy.Epic, source.Status)
	}

	dependsOn := make([]string, 0, len(source.DependsOn))
	for _, ref := range source.DependsOn {
		id, ok := resolveObjectiveDependency(plan, target.Legacy.Release, ref)
		if !ok {
			return "", fmt.Errorf("%w: epic %s depends_on %q", ErrUnresolvedDependency, target.Legacy.Epic, ref)
		}
		dependsOn = append(dependsOn, id)
	}

	legacyFields, err := unknownFields(fm, objectiveKnownFrontmatterKeys)
	if err != nil {
		return "", fmt.Errorf("convert objective %s: %s: %w", target.GlobalID, target.Legacy.Path, err)
	}

	title, ok := extractH1(body)
	if !ok {
		title = target.Legacy.Epic
	}

	out := objectiveOutputFrontmatter{
		ID:           target.GlobalID,
		Title:        title,
		Status:       status,
		DependsOn:    dependsOn,
		Release:      firstNonEmptyString(target.ReleaseID, target.Legacy.Release),
		LegacyFields: legacyFields,
	}

	marshalled, err := yaml.Marshal(&out)
	if err != nil {
		return "", fmt.Errorf("convert objective %s: marshal yaml: %w", target.GlobalID, err)
	}

	renderedBody := relocatedBody("epic", target.Legacy.Path, target.Legacy.Release, body)
	content := "---\n" + strings.TrimSpace(string(marshalled)) + "\n---" + renderedBody
	if crlf {
		content = strings.ReplaceAll(content, "\n", "\r\n")
	}

	if _, err := data.DecodeObjectiveV2(target.TargetPath, content); err != nil {
		return "", fmt.Errorf("convert objective %s: rendered content did not decode: %w", target.GlobalID, err)
	}

	return content, nil
}

// ConvertTask renders the V2 Task file content for target, reading the V1
// task source fresh from root. plan supplies the legacy prerequisites (if
// any) an archived, completed dependency of this task requires an authored
// line for.
func ConvertTask(root string, plan *ConversionPlan, target PlannedTarget) (string, error) {
	if target.Kind != TargetTask {
		return "", fmt.Errorf("convert task: target %s is not a task target", target.GlobalID)
	}

	objectiveID, ok := ownerObjectiveID(target.TargetPath)
	if !ok {
		return "", fmt.Errorf("convert task %s: target path %q does not name an owning objective", target.GlobalID, target.TargetPath)
	}

	raw, err := readSourceFile(root, target.Legacy.Path)
	if err != nil {
		return "", err
	}
	crlf := strings.Contains(raw, "\r\n")
	fm, body, err := data.SplitFrontmatterBody(raw)
	if err != nil {
		return "", fmt.Errorf("convert task %s: %s: %w", target.GlobalID, target.Legacy.Path, err)
	}

	var source taskSourceFrontmatter
	if err := yaml.Unmarshal([]byte(fm), &source); err != nil {
		return "", fmt.Errorf("convert task %s: %s: %w", target.GlobalID, target.Legacy.Path, err)
	}

	rawStatus := source.Status
	if rawStatus == "" {
		rawStatus = source.Column
	}
	status, _, recognized := resolveTaskStatus(rawStatus)
	if !recognized {
		return "", fmt.Errorf("%w: task %s status %q", ErrAmbiguousLifecycle, target.Legacy.OriginalID, rawStatus)
	}

	var stage string
	if status == string(data.ColumnInProgress) {
		rawStage := source.Stage
		if rawStage == "" {
			rawStage = source.Phase
		}
		resolved, recognized := resolveTaskStage(data.ProgressStage(rawStage))
		if !recognized {
			return "", fmt.Errorf("%w: task %s stage %q", ErrAmbiguousLifecycle, target.Legacy.OriginalID, rawStage)
		}
		stage = string(resolved)
	}

	title := firstNonEmptyString(source.Title, source.Objective)

	dependsOn := make([]taskDependsOnOutput, 0, len(target.DependsOn))
	for _, id := range target.DependsOn {
		dependsOn = append(dependsOn, taskDependsOnOutput{Task: id})
	}

	legacyFields, err := unknownFields(fm, taskKnownFrontmatterKeys)
	if err != nil {
		return "", fmt.Errorf("convert task %s: %s: %w", target.GlobalID, target.Legacy.Path, err)
	}

	out := taskOutputFrontmatter{
		ID:           target.GlobalID,
		Title:        title,
		Objective:    objectiveID,
		PlannedBy:    taskActorOutput{Role: string(data.ActorRolePlanner), Session: "migration"},
		Status:       status,
		Stage:        stage,
		DependsOn:    dependsOn,
		Release:      target.Legacy.Release,
		LegacyFields: legacyFields,
	}

	marshalled, err := yaml.Marshal(&out)
	if err != nil {
		return "", fmt.Errorf("convert task %s: marshal yaml: %w", target.GlobalID, err)
	}

	renderedBody := relocatedBody("task", target.Legacy.Path, target.Legacy.Release, body)
	renderedBody += legacyPrerequisiteBody(plan, target.GlobalID)

	content := "---\n" + strings.TrimSpace(string(marshalled)) + "\n---" + renderedBody
	if crlf {
		content = strings.ReplaceAll(content, "\n", "\r\n")
	}

	if _, err := data.DecodeTaskV2(target.TargetPath, content); err != nil {
		return "", fmt.Errorf("convert task %s: rendered content did not decode: %w", target.GlobalID, err)
	}

	return content, nil
}

// relocatedBody wraps the exact V1 body (already newline-normalized by
// SplitFrontmatterBody) under one heading naming the V1 source it was
// relocated from verbatim, byte for byte, comments and all.
func relocatedBody(kind, legacyPath, release, body string) string {
	var b strings.Builder
	b.WriteString("\n## Migrated from V1\n\n")
	fmt.Fprintf(&b, "Relocated verbatim from the V1 %s body at `%s`", kind, legacyPath)
	if release != "" {
		fmt.Fprintf(&b, " (release `%s`)", release)
	}
	b.WriteString(".\n\n## V1 Body (verbatim)")
	b.WriteString(body)
	return b.String()
}

// legacyPrerequisiteBody renders one authored paragraph per LegacyPrerequisite
// recorded against taskGlobalID, naming the archive path and the original
// recorded completion evidence — never a fabricated depends_on entry.
func legacyPrerequisiteBody(plan *ConversionPlan, taskGlobalID string) string {
	var matches []LegacyPrerequisite
	for _, p := range plan.Prereqs {
		if p.Task == taskGlobalID {
			matches = append(matches, p)
		}
	}
	if len(matches) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString("\n\n## Legacy Prerequisite\n")
	for _, p := range matches {
		b.WriteString("\nThis Task's V1 source depended on completed V1 work that received no V2 identity, ")
		fmt.Fprintf(&b, "archived at `%s`. Recorded completion evidence: %s\n", p.ArchivePath, p.Evidence)
	}
	return b.String()
}

// resolveObjectiveDependency maps a declared V1 epic dependency reference
// (an epic directory id, e.g. "E01-example") to the allocated O-### of the
// Objective converted from the same release's epic of that id.
func resolveObjectiveDependency(plan *ConversionPlan, release, ref string) (string, bool) {
	for _, t := range plan.Targets {
		if t.Kind != TargetObjective {
			continue
		}
		if t.Legacy.Release == release && t.Legacy.Epic == ref {
			return t.GlobalID, true
		}
	}
	return "", false
}

// resolveTaskStage canonicalizes a raw stage value using the same alias
// table data.NormalizeTaskStageForLoad applies at V1 load time, but — like
// resolveTaskStatus in plan.go — never heals a genuinely unrecognized value.
func resolveTaskStage(raw data.ProgressStage) (data.ProgressStage, bool) {
	if data.IsCanonicalStage(raw) {
		return raw, true
	}
	if data.IsLegacyTaskStageAlias(raw) {
		return data.NormalizeTaskStageForLoad(raw), true
	}
	return "", false
}

// ownerObjectiveID extracts the O-### owner from a Task target path shaped
// "objectives/{O-###}-{slug}/tasks/{T-###}-{slug}", the layout plan.go always
// builds — rather than re-deriving ownership by any other means.
func ownerObjectiveID(targetPath string) (string, bool) {
	parts := strings.Split(targetPath, "/")
	if len(parts) < 2 || parts[0] != v2ObjectivesDir {
		return "", false
	}
	seg := parts[1]
	if !strings.HasPrefix(seg, "O-") {
		return "", false
	}
	slugOffset := strings.IndexByte(seg[2:], '-')
	if slugOffset < 3 {
		return "", false
	}
	identityEnd := 2 + slugOffset
	for _, digit := range seg[2:identityEnd] {
		if digit < '0' || digit > '9' {
			return "", false
		}
	}
	return seg[:identityEnd], true
}

// unknownFields parses raw frontmatter YAML and returns every key not in
// known, so a V1 source field convert.go does not otherwise interpret is
// preserved on the converted record rather than silently dropped. A source
// carrying no unknown field returns a nil map.
func unknownFields(rawFrontmatter string, known []string) (map[string]any, error) {
	var all map[string]any
	if err := yaml.Unmarshal([]byte(rawFrontmatter), &all); err != nil {
		return nil, err
	}
	for _, key := range known {
		delete(all, key)
	}
	if len(all) == 0 {
		return nil, nil
	}
	return all, nil
}

// extractH1 returns the text of the first level-1 Markdown heading in body,
// trimmed, so an Objective's title can come from the epic's own authored
// heading rather than a derived or invented string.
func extractH1(body string) (string, bool) {
	for _, line := range strings.Split(body, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(trimmed, "# ")), true
		}
	}
	return "", false
}

func firstNonEmptyString(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func readSourceFile(root, relPath string) (string, error) {
	content, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(relPath)))
	if err != nil {
		return "", fmt.Errorf("read %s: %w", relPath, err)
	}
	return string(content), nil
}
