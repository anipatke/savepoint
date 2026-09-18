// Package migrate's plan.go builds the deterministic, write-free conversion
// plan a user reviews before any V1-to-V2 migration writes anything. Plan
// reads project-owned source content (never writes it), allocates fresh
// global V2 identities for active work, and records everything else as a
// byte-preserved archive entry or a named ambiguity requiring an owner
// decision.
//
// Scope: Plan decides *identity and destination* — which V1 source becomes
// which V2 Objective/Task/Issue, or an archive entry, and the legacy
// reference map behind it. Rendering the converted file content (T004-T006)
// and applying the plan to disk (T008-T009) are later tasks. Plan therefore
// does not yet assign a fate to config.yml, router.md, PRD.md, Design.md,
// Health-Check.md, the managed guide, or skill files: those require content
// decisions this task does not make. Immutable, byte-preserve-only history
// (epic audits, the audit prompt/register/runs, resolved/verified/waived
// dispositions, and genuinely unclassified files) is archived here because
// archiving never requires a content decision.
package migrate

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/opencode/savepoint/internal/data"
	"gopkg.in/yaml.v3"
)

// TargetKind names which V2 record family a PlannedTarget becomes.
type TargetKind string

const (
	TargetObjective TargetKind = "objective"
	TargetTask      TargetKind = "task"
	TargetIssue     TargetKind = "issue"
)

// LegacyKey source-qualifies one V1 record's identity. Path is the actual
// uniqueness qualifier — a short V1 ID such as T001-shared can legitimately
// recur across releases (and, for defects and findings, is not even
// epic-scoped) — while Release, Epic, and OriginalID are carried alongside
// for display and for resolving a V1 short depends_on reference back to the
// record it named.
type LegacyKey struct {
	Release    string
	Epic       string // "" when the source is not epic-scoped
	Path       string
	OriginalID string
}

// PlannedTarget is one V1 source becoming one new V2 record with a freshly
// allocated global identity. TargetPath is relative to the .savepoint root.
type PlannedTarget struct {
	Kind       TargetKind
	GlobalID   string
	Legacy     LegacyKey
	TargetPath string
	// DependsOn lists the global Task IDs (only meaningful for TargetTask)
	// this target's converted depends_on will name. A dependency on
	// archived, completed work never appears here; see LegacyPrerequisite.
	DependsOn []string
	// DuplicateOfGlobalID is set only for a TargetIssue planned from a
	// `duplicate` finding whose canonical finding also converts: the
	// allocated I### of that canonical Issue. Empty for every other target,
	// including a duplicate finding whose canonical is archive-only (that
	// duplicate is archived too; see planFinding).
	DuplicateOfGlobalID string
}

// ArchiveEntry is one V1 source preserved byte-for-byte under archive/v1/
// instead of being converted, because it is settled history (a done Task, a
// closed epic, a resolved defect, a verified or waived finding, a duplicate
// finding whose canonical is itself archived) or because it matched no known
// role at all.
type ArchiveEntry struct {
	SourcePath  string
	ArchivePath string
	Role        Role
	Legacy      *LegacyKey // nil when the source carries no legacy identity (e.g. an unclassified file)
}

// LegacyPrerequisite records a dependency from a newly allocated active Task
// onto archived, completed V1 work that received no V2 identity. It is the
// mechanism T002's "no V2 schema change" boundary requires: the fact is
// recorded and resolvable, but it is never a fabricated depends_on entry and
// never a fabricated Check.
type LegacyPrerequisite struct {
	Task        string // the new T### that named the archived work as a dependency
	ArchivePath string
	Evidence    string // the original recorded completion (or waiver) evidence, best-effort from the source
}

// WaivedReference records that an active converted Task's V1 source (via its
// work_item or tasks reference) named a waived audit finding. A waiver is a
// closed owner decision, not follow-up, so it never becomes an Issue — but
// the fact that active work is aware of it must stay resolvable rather than
// silently vanish. This is the waiver counterpart to LegacyPrerequisite.
type WaivedReference struct {
	Task        string // the global Task ID whose V1 source referenced the waived finding
	ArchivePath string
	Reason      string // the original recorded waiver reason, best-effort from the source
}

// ConflictKind names a reason Plan refuses to proceed with normal planning.
type ConflictKind string

const (
	// ConflictExistingManifest means .savepoint/migrations/v1-to-v2.yml
	// already exists on a project whose schema is still V1 — distinct from
	// SchemaAlreadyV2, which is a clean already-migrated project.
	ConflictExistingManifest ConflictKind = "existing_manifest_conflict"
)

// Conflict is one named reason Plan (or a later apply) must not proceed
// without owner attention.
type Conflict struct {
	Kind   ConflictKind
	Path   string
	Detail string
}

// AmbiguityKind names one of the four blocking-ambiguity shapes the epic
// defines. Plan detects the first two directly, since they fall out of
// identity allocation and dependency resolution; duplicate source identity
// and unresolved narrative findings are named here for the shared vocabulary
// but are not yet detected by this task — see T007.
type AmbiguityKind string

const (
	AmbiguityUnrecognizedLifecycle   AmbiguityKind = "unrecognized_lifecycle"
	AmbiguityMissingDependencyTarget AmbiguityKind = "missing_dependency_target"
	AmbiguityDuplicateSourceIdentity AmbiguityKind = "duplicate_source_identity"
	AmbiguityUnresolvedNarrativeFind AmbiguityKind = "unresolved_narrative_finding"
)

// Ambiguity is one named, stably identified condition that blocks apply
// until an owner decision resolves it. ID is derived from Kind and Path, so
// the same project previews the same ambiguity IDs every run.
type Ambiguity struct {
	ID       string
	Kind     AmbiguityKind
	Path     string
	Detail   string
	Resolved bool
	Decision string // the owner-supplied resolution, once Decisions names this ID
}

// Decisions maps an Ambiguity ID to the owner-supplied resolution recorded
// against it. The full decisions-file schema (T007) reads this from an
// explicit input file; Plan itself only consumes the resolved map.
type Decisions map[string]string

// Clock and OperationIDSource are injected so Plan and its callers can
// render byte-identical output across repeated runs in tests, and so a real
// apply can record when and under which operation a plan was produced.
type Clock func() time.Time
type OperationIDSource func() string

// Plan is the deterministic, write-free output of building a conversion
// plan. It is a value: reviewing it is reading a Go struct, never running a
// second pass over the filesystem.
type ConversionPlan struct {
	// SchemaAlreadyV2 means the project already declares schema_version: 2;
	// every other field is zero and there is nothing to migrate.
	SchemaAlreadyV2 bool

	GeneratedAt time.Time
	OperationID string

	Targets     []PlannedTarget
	Archives    []ArchiveEntry
	Prereqs     []LegacyPrerequisite
	WaivedRefs  []WaivedReference
	Conflicts   []Conflict
	Ambiguities []Ambiguity

	// Sources is the raw inventory Plan computed, retained so a manifest can
	// be built from the same read without walking the project a second time.
	Sources []SourceFile
}

const (
	migrationsDirRel = ".savepoint/migrations"
	manifestFileName = "v1-to-v2.yml"
	v2ObjectivesDir  = "objectives"
	v2ObjectiveFile  = "Objective.md"
	v2TasksDir       = "tasks"
	v2IssuesDir      = "issues"
)

func manifestRelPath() string {
	return migrationsDirRel + "/" + manifestFileName
}

// isPreservedMigrationsContent reports whether path is existing content
// under .savepoint/migrations/ that migration must leave exactly alone:
// internal/init's archived legacy skill copies and its README.md. It is
// never archived and never reported as unclassified, per the coexistence
// rule — it already lives at a V2-shaped location and rewriting or archiving
// it would duplicate content the project already preserved once.
func isPreservedMigrationsContent(path string) bool {
	return strings.HasPrefix(path, migrationsDirRel+"/") && path != manifestRelPath()
}

// Plan builds the deterministic conversion plan for the V1 project rooted at
// projectRoot (the directory containing AGENTS.md, agent-skills/, and
// .savepoint/). It opens no file for writing and creates no directory: every
// filesystem interaction is a read. now and newOperationID are injected so
// repeated calls over the same project produce byte-identical output.
func Plan(projectRoot string, decisions Decisions, now Clock, newOperationID OperationIDSource) (*ConversionPlan, error) {
	rootAbs, err := filepath.Abs(projectRoot)
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(rootAbs, ".savepoint", "config.yml")
	version, err := data.ReadSchemaVersion(configPath)
	if err != nil {
		return nil, err
	}
	if version == data.SchemaVersionV2 {
		return &ConversionPlan{SchemaAlreadyV2: true}, nil
	}

	sources, err := Inventory(rootAbs)
	if err != nil {
		return nil, err
	}

	for _, f := range sources {
		if f.Path == manifestRelPath() {
			return &ConversionPlan{
				Sources: sources,
				Conflicts: []Conflict{{
					Kind: ConflictExistingManifest,
					Path: f.Path,
					Detail: "a v1-to-v2.yml manifest already exists but schema_version is not 2; " +
						"resolve the conflict before migrating rather than overwrite recorded migration history",
				}},
			}, nil
		}
	}

	b := &planBuilder{
		root:      rootAbs,
		decisions: decisions,
		byPath:    indexSourcesByPath(sources),
		ids:       newIDAllocator(),
	}
	if err := b.build(); err != nil {
		return nil, err
	}

	return &ConversionPlan{
		GeneratedAt: now(),
		OperationID: newOperationID(),
		Targets:     b.targets,
		Archives:    b.archives,
		Prereqs:     b.prereqs,
		WaivedRefs:  b.waivedRefs,
		Ambiguities: b.ambiguities,
		Sources:     sources,
	}, nil
}

func indexSourcesByPath(sources []SourceFile) map[string]SourceFile {
	byPath := make(map[string]SourceFile, len(sources))
	for _, f := range sources {
		byPath[f.Path] = f
	}
	return byPath
}

// idAllocator assigns the next unused O/T/I number in allocation order.
// Archived-only records never call allocate, so they reserve nothing —
// exactly the "never reused, never reserved for archive-only work" rule.
type idAllocator struct {
	next map[string]int
}

func newIDAllocator() *idAllocator {
	return &idAllocator{next: map[string]int{"O": 1, "T": 1, "I": 1}}
}

func (a *idAllocator) allocate(prefix string) string {
	n := a.next[prefix]
	a.next[prefix] = n + 1
	return fmt.Sprintf("%s%03d", prefix, n)
}

// planBuilder accumulates one Plan's worth of decisions while walking the
// project in a fixed, deterministic order: releases, then epics within a
// release, then tasks within an epic, then a release's defects, then the
// project's audit findings.
type planBuilder struct {
	root      string
	decisions Decisions
	byPath    map[string]SourceFile

	ids *idAllocator

	targets     []PlannedTarget
	archives    []ArchiveEntry
	prereqs     []LegacyPrerequisite
	waivedRefs  []WaivedReference
	ambiguities []Ambiguity

	// taskByLegacyPath tracks the outcome of every task this build has
	// planned so far (its global ID if active, or its archive path if not),
	// keyed by source path, so later tasks can resolve their depends_on
	// against project-wide state exactly as data.ResolveDependency does.
	taskByLegacyPath map[string]plannedTaskOutcome
	// issueByFindingID tracks whether a given finding ID converted to an
	// Issue, so a duplicate finding can tell whether its canonical did too.
	issueByFindingID map[string]string // finding ID -> allocated I### (absent = archived)
}

type plannedTaskOutcome struct {
	task        data.Task
	globalID    string // "" when archived
	archivePath string // "" when active
}

func (b *planBuilder) build() error {
	b.taskByLegacyPath = map[string]plannedTaskOutcome{}
	b.issueByFindingID = map[string]string{}

	savepointRoot := filepath.Join(b.root, ".savepoint")
	discover := data.NewDiscover()

	releases, err := discover.ListReleases(savepointRoot)
	if err != nil {
		if _, statErr := os.Stat(filepath.Join(savepointRoot, "releases")); os.IsNotExist(statErr) {
			releases = nil // an empty V1 project has no releases yet; nothing to plan
		} else {
			return fmt.Errorf("list releases: %w", err)
		}
	}

	// Pass 1: plan every epic and its tasks, in release/epic order, so task
	// dependency resolution (pass 1 itself, since a task may depend on an
	// earlier task in the same epic) and later defect/finding passes can see
	// the complete task outcome map.
	for _, release := range releases {
		epics, err := discover.ListEpics(savepointRoot, release.ID)
		if err != nil {
			if _, statErr := os.Stat(filepath.Join(savepointRoot, "releases", release.ID, "epics")); os.IsNotExist(statErr) {
				continue // a release with no epics yet has nothing to plan
			}
			return fmt.Errorf("list epics for release %s: %w", release.ID, err)
		}
		for _, epic := range epics {
			if err := b.planEpic(discover, savepointRoot, release.ID, epic.ID); err != nil {
				return err
			}
		}
	}

	// Pass 2: release-scoped defects.
	for _, release := range releases {
		defects, err := discover.ListDefects(savepointRoot, release.ID)
		if err != nil {
			return fmt.Errorf("list defects for release %s: %w", release.ID, err)
		}
		for _, d := range defects {
			if err := b.planDefect(release.ID, d); err != nil {
				return err
			}
		}
	}

	// Pass 3: project-level audit findings.
	if err := b.planFindings(); err != nil {
		return err
	}

	// Pass 4: immutable audit history (runs, prompt, register) and epic
	// audits, plus anything genuinely unclassified. These never affect
	// identity allocation, so they run last and independently of the passes
	// above.
	b.archiveRemainingRoles()

	sort.SliceStable(b.ambiguities, func(i, j int) bool { return b.ambiguities[i].ID < b.ambiguities[j].ID })
	sort.SliceStable(b.waivedRefs, func(i, j int) bool {
		if b.waivedRefs[i].Task != b.waivedRefs[j].Task {
			return b.waivedRefs[i].Task < b.waivedRefs[j].Task
		}
		return b.waivedRefs[i].ArchivePath < b.waivedRefs[j].ArchivePath
	})

	return nil
}

// epicRawFrontmatter reads only the fields plan.go needs from an epic detail
// file: its raw, unhealed status. Epic status is not part of the strict V2
// decoders (there is no EpicDetail record family), so this is a small local
// projection rather than a reuse of an existing typed loader.
type epicRawFrontmatter struct {
	Status string `yaml:"status"`
}

func (b *planBuilder) planEpic(discover *data.Discover, savepointRoot, release, epic string) error {
	epicDir := filepath.ToSlash(filepath.Join(".savepoint", "releases", release, "epics", epic))
	sf, ok := b.findRoleUnder(RoleEpicDetail, epicDir)
	if !ok {
		// No detail file at all is a structural gap outside this task's
		// scope; archive nothing and plan no tasks under an epic Plan cannot
		// name, rather than guess a status.
		return nil
	}

	content, err := b.readSource(sf.Path)
	if err != nil {
		return err
	}
	fm, _, err := data.SplitFrontmatterBody(content)
	if err != nil {
		return fmt.Errorf("parse epic detail %s: %w", sf.Path, err)
	}
	var raw epicRawFrontmatter
	if err := yaml.Unmarshal([]byte(fm), &raw); err != nil {
		return fmt.Errorf("parse epic detail %s: %w", sf.Path, err)
	}

	status, healed, recognized := resolveEpicStatus(raw.Status)
	legacy := LegacyKey{Release: release, Epic: epic, Path: sf.Path, OriginalID: epic}

	if !recognized {
		b.addAmbiguity(AmbiguityUnrecognizedLifecycle, sf.Path,
			fmt.Sprintf("epic %s/%s has unrecognized status %q; no Objective or Task under it can be planned until an owner decision resolves it", release, epic, raw.Status))
		return nil
	}
	_ = healed

	closed := status == "done" || status == "audited"

	if closed {
		b.archives = append(b.archives, ArchiveEntry{
			SourcePath:  sf.Path,
			ArchivePath: archivePathFor(sf.Path),
			Role:        RoleEpicDetail,
			Legacy:      &legacy,
		})
		return b.archiveAllTasks(discover, savepointRoot, release, epic)
	}

	objectiveID := b.ids.allocate("O")
	b.targets = append(b.targets, PlannedTarget{
		Kind:       TargetObjective,
		GlobalID:   objectiveID,
		Legacy:     legacy,
		TargetPath: filepath.ToSlash(filepath.Join(v2ObjectivesDir, objectiveID+"-"+slugOf(epic), v2ObjectiveFile)),
	})
	// The original epic detail is archived alongside conversion: migration
	// preserves the source bytes even though their content became an
	// Objective, per the epic's "original bytes archived" rule.
	b.archives = append(b.archives, ArchiveEntry{
		SourcePath:  sf.Path,
		ArchivePath: archivePathFor(sf.Path),
		Role:        RoleEpicDetail,
		Legacy:      &legacy,
	})

	return b.planTasks(discover, savepointRoot, release, epic, objectiveID, false)
}

// archiveAllTasks archives every task under a closed epic regardless of the
// task's own recorded status, since a closed epic's remaining work has no
// active V2 owner to attach a Task to.
func (b *planBuilder) archiveAllTasks(discover *data.Discover, savepointRoot, release, epic string) error {
	return b.planTasksArchiveOnly(discover, savepointRoot, release, epic)
}

func (b *planBuilder) planTasksArchiveOnly(discover *data.Discover, savepointRoot, release, epic string) error {
	return b.planTasksImpl(discover, savepointRoot, release, epic, "", true)
}

func (b *planBuilder) planTasks(discover *data.Discover, savepointRoot, release, epic, objectiveID string, archiveOnly bool) error {
	return b.planTasksImpl(discover, savepointRoot, release, epic, objectiveID, archiveOnly)
}

func (b *planBuilder) planTasksImpl(discover *data.Discover, savepointRoot, release, epic, objectiveID string, archiveOnly bool) error {
	taskInfos, err := discover.ListTasks(savepointRoot, release, epic)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil
		}
		return fmt.Errorf("list tasks for %s/%s: %w", release, epic, err)
	}

	parser := data.NewParser()
	for _, info := range taskInfos {
		rel := filepath.ToSlash(mustRel(b.root, info.Path))
		content, err := b.readSource(rel)
		if err != nil {
			return err
		}
		task, err := parser.ParseTaskFile(rel, content)
		if err != nil {
			return fmt.Errorf("parse task %s: %w", rel, err)
		}
		task.Release = release
		task.Epic = epic

		legacy := LegacyKey{Release: release, Epic: epic, Path: rel, OriginalID: task.ID}

		raw, err := readRawTaskStatus(content)
		if err != nil {
			return fmt.Errorf("parse task %s: %w", rel, err)
		}
		status, _, recognized := resolveTaskStatus(raw)
		if !recognized {
			b.addAmbiguity(AmbiguityUnrecognizedLifecycle, rel,
				fmt.Sprintf("task %s has unrecognized status %q; it cannot be planned as active or archived until an owner decision resolves it", task.ID, raw))
			continue
		}

		if archiveOnly || status == string(data.ColumnDone) {
			archivePath := archivePathFor(rel)
			b.archives = append(b.archives, ArchiveEntry{
				SourcePath:  rel,
				ArchivePath: archivePath,
				Role:        RoleTask,
				Legacy:      &legacy,
			})
			b.taskByLegacyPath[rel] = plannedTaskOutcome{task: *task, archivePath: archivePath}
			continue
		}

		taskID := b.ids.allocate("T")
		b.taskByLegacyPath[rel] = plannedTaskOutcome{task: *task, globalID: taskID}
		b.targets = append(b.targets, PlannedTarget{
			Kind:       TargetTask,
			GlobalID:   taskID,
			Legacy:     legacy,
			TargetPath: filepath.ToSlash(filepath.Join(v2ObjectivesDir, objectiveID+"-"+slugOf(epic), v2TasksDir, taskID+"-"+slugOf(task.ID))),
		})
	}

	// Dependencies are resolved in a second inner pass so a task depending on
	// a later-numbered sibling in the same epic still resolves correctly.
	for _, info := range taskInfos {
		rel := filepath.ToSlash(mustRel(b.root, info.Path))
		outcome, ok := b.taskByLegacyPath[rel]
		if !ok || outcome.globalID == "" {
			continue // archived, or blocked by an unrecognized-status ambiguity
		}
		if err := b.resolveTaskDependencies(rel, outcome); err != nil {
			return err
		}
	}

	return nil
}

// resolveTaskDependencies resolves one active task's depends_on references
// against every task outcome planned so far. A reference landing on an
// archived, completed task becomes a LegacyPrerequisite rather than a
// depends_on entry; a reference that resolves to nothing becomes a blocking
// ambiguity.
func (b *planBuilder) resolveTaskDependencies(rel string, outcome plannedTaskOutcome) error {
	if len(outcome.task.DependsOn) == 0 {
		return nil
	}

	allTasks := make([]data.Task, 0, len(b.taskByLegacyPath))
	for _, o := range b.taskByLegacyPath {
		allTasks = append(allTasks, o.task)
	}

	for i := range b.targets {
		if b.targets[i].Legacy.Path != rel || b.targets[i].Kind != TargetTask {
			continue
		}
		target := &b.targets[i]
		for _, ref := range outcome.task.DependsOn {
			resolution := data.ResolveDependency(ref, outcome.task, allTasks, map[string]string{})
			if resolution.Kind != data.DependencyTask {
				b.addAmbiguity(AmbiguityMissingDependencyTarget, rel,
					fmt.Sprintf("task %s depends_on %q, which does not resolve to any known task", outcome.task.ID, ref))
				continue
			}
			depOutcome := b.findTaskOutcomeByID(resolution.ID, resolution.TaskStatus)
			switch {
			case depOutcome.globalID != "":
				target.DependsOn = append(target.DependsOn, depOutcome.globalID)
			case depOutcome.archivePath != "":
				b.prereqs = append(b.prereqs, LegacyPrerequisite{
					Task:        target.GlobalID,
					ArchivePath: depOutcome.archivePath,
					Evidence:    completionEvidence(depOutcome.task),
				})
			default:
				b.addAmbiguity(AmbiguityMissingDependencyTarget, rel,
					fmt.Sprintf("task %s depends_on %q, which resolved but was never planned", outcome.task.ID, ref))
			}
		}
	}
	return nil
}

func (b *planBuilder) findTaskOutcomeByID(id string, status data.ColumnType) plannedTaskOutcome {
	for _, o := range b.taskByLegacyPath {
		if o.task.ID == id && o.task.Column == status {
			return o
		}
	}
	// Fall back to matching on ID alone: two releases never share a task
	// outcome for the same ID with different status in practice, but this
	// keeps the lookup total rather than panicking on an edge case.
	for _, o := range b.taskByLegacyPath {
		if o.task.ID == id {
			return o
		}
	}
	return plannedTaskOutcome{}
}

// completionEvidence is a best-effort human-readable record of why an
// archived task counted as complete. V1 has no dedicated completion-evidence
// field, so this names what is actually available rather than fabricating
// structured proof.
func completionEvidence(task data.Task) string {
	var done []string
	for _, item := range task.Acceptance {
		done = append(done, item)
	}
	if len(done) > 0 {
		return "recorded done; acceptance criteria: " + strings.Join(done, "; ")
	}
	return "recorded done in the V1 source; no further completion evidence was authored"
}

func (b *planBuilder) planDefect(release string, info data.DefectInfo) error {
	rel := filepath.ToSlash(mustRel(b.root, info.Path))
	content, err := b.readSource(rel)
	if err != nil {
		return err
	}
	defect, err := data.NewParser().ParseDefectFile(rel, content)
	if err != nil {
		return fmt.Errorf("parse defect %s: %w", rel, err)
	}
	legacy := LegacyKey{Release: release, Path: rel, OriginalID: defect.ID}

	if defect.Status == data.DefectResolved {
		b.archives = append(b.archives, ArchiveEntry{
			SourcePath:  rel,
			ArchivePath: archivePathFor(rel),
			Role:        RoleDefect,
			Legacy:      &legacy,
		})
		return nil
	}

	issueID := b.ids.allocate("I")
	b.targets = append(b.targets, PlannedTarget{
		Kind:       TargetIssue,
		GlobalID:   issueID,
		Legacy:     legacy,
		TargetPath: filepath.ToSlash(filepath.Join(v2IssuesDir, issueID+"-"+slugOf(defect.ID)+".md")),
	})
	return nil
}

func (b *planBuilder) planFindings() error {
	var findingPaths []string
	for path := range b.byPath {
		if Classify(path) == RoleFinding {
			findingPaths = append(findingPaths, path)
		}
	}
	sort.Strings(findingPaths)

	parser := data.NewParser()
	findings := make(map[string]*data.AuditFinding, len(findingPaths))
	for _, path := range findingPaths {
		content, err := b.readSource(path)
		if err != nil {
			return err
		}
		finding, err := parser.ParseFindingFile(path, content)
		if err != nil {
			return fmt.Errorf("parse finding %s: %w", path, err)
		}
		findings[finding.ID] = finding
	}

	for _, path := range findingPaths {
		finding := findingByPath(findings, path)
		if finding == nil {
			continue
		}
		if err := b.planFinding(path, finding, findings); err != nil {
			return err
		}
	}
	return nil
}

func findingByPath(findings map[string]*data.AuditFinding, path string) *data.AuditFinding {
	for _, f := range findings {
		if f.Path == path {
			return f
		}
	}
	return nil
}

func (b *planBuilder) planFinding(path string, finding *data.AuditFinding, all map[string]*data.AuditFinding) error {
	legacy := LegacyKey{Path: path, OriginalID: finding.ID}

	switch finding.Status {
	case data.FindingOpen, data.FindingTriaged, data.FindingMapped, data.FindingDeferred, data.FindingOwnerDecision,
		data.FindingInProgress, data.FindingFixed:
		issueID := b.ids.allocate("I")
		b.issueByFindingID[finding.ID] = issueID
		b.targets = append(b.targets, PlannedTarget{
			Kind:       TargetIssue,
			GlobalID:   issueID,
			Legacy:     legacy,
			TargetPath: filepath.ToSlash(filepath.Join(v2IssuesDir, issueID+"-"+slugOf(finding.ID)+".md")),
		})
		return nil

	case data.FindingVerified:
		b.archives = append(b.archives, ArchiveEntry{
			SourcePath:  path,
			ArchivePath: archivePathFor(path),
			Role:        RoleFinding,
			Legacy:      &legacy,
		})
		return nil

	case data.FindingWaived:
		archivePath := archivePathFor(path)
		b.archives = append(b.archives, ArchiveEntry{
			SourcePath:  path,
			ArchivePath: archivePath,
			Role:        RoleFinding,
			Legacy:      &legacy,
		})
		b.recordWaivedReferences(finding, archivePath)
		return nil

	case data.FindingDuplicate:
		canonical, exists := all[finding.DuplicateOf]
		if !exists {
			b.addAmbiguity(AmbiguityUnresolvedNarrativeFind, path,
				fmt.Sprintf("finding %s is marked duplicate of %q, which names no known finding in this project", finding.ID, finding.DuplicateOf))
			return nil
		}

		canonicalID, canonicalConverts := b.issueByFindingID[canonical.ID]
		if !canonicalConverts && !dispositionArchivesOutright(canonical.Status) {
			// The canonical finding sorts after this duplicate but will
			// still convert; resolve forward.
			id, err := b.planFindingForward(canonical)
			if err != nil {
				return err
			}
			canonicalID, canonicalConverts = id, id != ""
		}
		if canonicalConverts {
			issueID := b.ids.allocate("I")
			b.issueByFindingID[finding.ID] = issueID
			b.targets = append(b.targets, PlannedTarget{
				Kind:                TargetIssue,
				GlobalID:            issueID,
				Legacy:              legacy,
				TargetPath:          filepath.ToSlash(filepath.Join(v2IssuesDir, issueID+"-"+slugOf(finding.ID)+".md")),
				DuplicateOfGlobalID: canonicalID,
			})
			return nil
		}
		b.archives = append(b.archives, ArchiveEntry{
			SourcePath:  path,
			ArchivePath: archivePathFor(path),
			Role:        RoleFinding,
			Legacy:      &legacy,
		})
		return nil

	default:
		b.addAmbiguity(AmbiguityUnrecognizedLifecycle, path,
			fmt.Sprintf("finding %s has unrecognized status %q", finding.ID, finding.Status))
		return nil
	}
}

// planFindingForward allocates a not-yet-visited canonical finding out of
// sorted order, so a duplicate that sorts before its canonical still resolves
// deterministically (allocation order is by ID text, so a forward reference
// is resolved once, memoized in issueByFindingID, and never re-allocated).
func (b *planBuilder) planFindingForward(canonical *data.AuditFinding) (string, error) {
	if id, ok := b.issueByFindingID[canonical.ID]; ok {
		return id, nil
	}
	if dispositionArchivesOutright(canonical.Status) {
		return "", nil
	}
	issueID := b.ids.allocate("I")
	b.issueByFindingID[canonical.ID] = issueID
	b.targets = append(b.targets, PlannedTarget{
		Kind:       TargetIssue,
		GlobalID:   issueID,
		Legacy:     LegacyKey{Path: canonical.Path, OriginalID: canonical.ID},
		TargetPath: filepath.ToSlash(filepath.Join(v2IssuesDir, issueID+"-"+slugOf(canonical.ID)+".md")),
	})
	return issueID, nil
}

// recordWaivedReferences checks whether a waived finding's own work_item or
// tasks reference names an active, already-planned Task, and records a
// WaivedReference when it does — the mechanism behind AC7: a waiver never
// becomes an Issue, but active work's awareness of it stays resolvable
// rather than silently dropped.
func (b *planBuilder) recordWaivedReferences(finding *data.AuditFinding, archivePath string) {
	refs := append([]string{}, finding.Tasks...)
	if finding.WorkItem != "" {
		refs = append(refs, finding.WorkItem)
	}

	seen := map[string]bool{}
	for _, ref := range refs {
		id, ok := resolveTaskReference(b.targets, finding.Releases, ref)
		if !ok || seen[id] {
			continue
		}
		seen[id] = true
		b.waivedRefs = append(b.waivedRefs, WaivedReference{
			Task:        id,
			ArchivePath: archivePath,
			Reason:      finding.WaiverReason,
		})
	}
}

// resolveTaskReference resolves a V1 defect/finding task reference (an
// epic-qualified task id, e.g. "E01-example/T001-shared" — exactly the shape
// task.ID and therefore LegacyKey.OriginalID already carry) to the global
// Task ID of the matching *active* Task target, scoped to one of releases
// when releases is non-empty. Only active targets are ever searched, so an
// archived task of the same short id in another release can never be
// mistaken for the one referenced — no release filter is even required for
// that case, since an archived task never becomes a TargetTask at all; the
// filter guards the rarer case of the same id being active in two releases
// at once. A reference that resolves to nothing is not an error: the epic
// requires no link be invented, only that a named one is not dropped.
func resolveTaskReference(targets []PlannedTarget, releases []string, ref string) (string, bool) {
	for _, t := range targets {
		if t.Kind != TargetTask || t.Legacy.OriginalID != ref {
			continue
		}
		if len(releases) == 0 || containsString(releases, t.Legacy.Release) {
			return t.GlobalID, true
		}
	}
	return "", false
}

func containsString(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func dispositionArchivesOutright(status data.FindingStatus) bool {
	switch status {
	case data.FindingVerified, data.FindingWaived:
		return true
	default:
		return false
	}
}

// archiveRemainingRoles archives every inventoried file whose role is
// immutable byte-preserved history (epic audits, the audit prompt/register,
// audit runs) or genuinely unclassified, skipping anything already archived
// or converted above and anything the migrations-directory coexistence rule
// preserves in place.
func (b *planBuilder) archiveRemainingRoles() {
	handled := map[string]bool{}
	for _, a := range b.archives {
		handled[a.SourcePath] = true
	}
	for _, t := range b.targets {
		handled[t.Legacy.Path] = true
	}

	var paths []string
	for path := range b.byPath {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	for _, path := range paths {
		if handled[path] || isPreservedMigrationsContent(path) {
			continue
		}
		role := Classify(path)
		switch role {
		case RoleEpicAudit, RoleAuditPrompt, RoleAuditRegister, RoleAuditRun, RoleUnclassified:
			b.archives = append(b.archives, ArchiveEntry{
				SourcePath:  path,
				ArchivePath: archivePathFor(path),
				Role:        role,
			})
		}
	}
}

func (b *planBuilder) addAmbiguity(kind AmbiguityKind, path, detail string) {
	id := string(kind) + ":" + path
	decision, resolved := b.decisions[id]
	b.ambiguities = append(b.ambiguities, Ambiguity{
		ID:       id,
		Kind:     kind,
		Path:     path,
		Detail:   detail,
		Resolved: resolved,
		Decision: decision,
	})
}

func (b *planBuilder) readSource(rel string) (string, error) {
	content, err := os.ReadFile(filepath.Join(b.root, filepath.FromSlash(rel)))
	if err != nil {
		return "", fmt.Errorf("read %s: %w", rel, err)
	}
	return string(content), nil
}

// archivePathFor is the deterministic archive/v1/ destination for a
// project-relative V1 source path: the same relative layout, moved under
// archive/v1/, so an archived path is always recoverable by inspection alone.
func archivePathFor(sourcePath string) string {
	return filepath.ToSlash(filepath.Join("archive", "v1", sourcePath))
}

// slugOf derives a deterministic destination slug from an already-authored
// V1 identifier (an epic directory name, a task/defect/finding id) by
// dropping its numeric ID prefix. Reusing the author's own slug — rather
// than deriving one from title text — keeps renaming/paraphrasing entirely
// out of identity planning.
func slugOf(id string) string {
	if slash := strings.LastIndexByte(id, '/'); slash >= 0 {
		id = id[slash+1:]
	}
	if dash := strings.IndexByte(id, '-'); dash >= 0 && dash+1 < len(id) {
		return id[dash+1:]
	}
	return id
}

// findRoleUnder returns the sole inventoried file with the given role
// directly inside dir, if any. Epic detail and audit filenames are named
// after the epic's short prefix (e.g. E01-Detail.md under epics/E01-example/),
// not its full directory name, so lookup goes by role and directory rather
// than by constructing an assumed filename.
func (b *planBuilder) findRoleUnder(role Role, dir string) (SourceFile, bool) {
	prefix := dir + "/"
	for path, sf := range b.byPath {
		if !strings.HasPrefix(path, prefix) {
			continue
		}
		if strings.Contains(path[len(prefix):], "/") {
			continue // nested further (e.g. tasks/); not directly under dir
		}
		if Classify(path) == role {
			return sf, true
		}
	}
	return SourceFile{}, false
}

func mustRel(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return path
	}
	return rel
}

// resolveEpicStatus canonicalizes a raw epic status string. planned and
// in_progress are active; done and audited are closed; a legacy task-style
// alias (todo/complete/completed) is healed exactly as
// data.ResolveEpicStatusAlias would; anything else is unrecognized.
func resolveEpicStatus(raw string) (status string, healed bool, recognized bool) {
	switch raw {
	case "planned", "in_progress", "done", "audited":
		return raw, false, true
	}
	if alias, ok := data.ResolveEpicStatusAlias(data.EpicStatus(raw)); ok {
		return string(alias), true, true
	}
	return "", false, false
}

// resolveTaskStatus canonicalizes a raw task status string using the same
// alias table data.ResolveTaskStatusAlias applies at V1 load time. Unlike
// ParseTaskFile, this never heals a genuinely unrecognized value to planned:
// migration must refuse it as a blocking ambiguity instead.
func resolveTaskStatus(raw string) (status string, healed bool, recognized bool) {
	col := data.ColumnType(raw)
	if data.IsCanonicalTaskStatus(col) {
		return raw, false, true
	}
	if alias, ok := data.ResolveTaskStatusAlias(col); ok {
		return string(alias), true, true
	}
	return "", false, false
}

// rawTaskStatusFrontmatter reads only status/column so the caller can decide
// recognition itself instead of going through ParseTaskFile's healing.
type rawTaskStatusFrontmatter struct {
	Status string `yaml:"status"`
	Column string `yaml:"column"`
}

func readRawTaskStatus(content string) (string, error) {
	fm, _, err := data.SplitFrontmatterBody(content)
	if err != nil {
		return "", err
	}
	var raw rawTaskStatusFrontmatter
	if err := yaml.Unmarshal([]byte(fm), &raw); err != nil {
		return "", err
	}
	if raw.Status != "" {
		return raw.Status, nil
	}
	return raw.Column, nil
}
