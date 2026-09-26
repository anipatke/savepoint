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
// and applying the plan to disk (T008-T009) are later tasks. Plan also
// assigns a fate to the small, fixed set of top-level project documents
// (T006): PRD.md relocates to Idea.md (Documents) with the original
// archived, router.md rewrites its state anchor in place (Documents), and
// Health-Check.md is archived with its candidate commands recorded for a
// future preview. config.yml, Design.md, Guardrails.md, and
// visual-identity.md are deliberately never archived or targeted — staying
// out of both plan.Archives and plan.Documents *is* their "preserved in
// place, untouched" outcome, for apply to rely on. The managed guide and
// skill files stay out of scope for E45 entirely (E46/E47 own the V2
// routing block and the public skills). Immutable, byte-preserve-only
// history (epic audits, the audit prompt/register/runs,
// resolved/verified/waived dispositions, and genuinely unclassified files)
// is archived here because archiving never requires a content decision.
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
	TargetRelease   TargetKind = "release"
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

// PlannedTarget is one V1 source converted to a V2 record or one migration-
// generated V2 record, with a freshly allocated global identity. TargetPath
// is relative to the .savepoint root.
type PlannedTarget struct {
	Kind       TargetKind
	GlobalID   string
	Legacy     LegacyKey
	TargetPath string
	// Generated marks a V2 record created by migration without a V1 source
	// record. It is used for the continuation Goal when the V1 router does not
	// resolve to a live Goal; generated records have no legacy identity mapping.
	Generated bool
	// GeneratedTitle is populated only for generated records.
	GeneratedTitle string
	// ReleaseID is the allocated V2 Release identity an Objective belongs to.
	// It is empty for records that are not Objective-owned by a Release.
	ReleaseID string
	// ReleaseStatus is the resolved V2 lifecycle for a TargetRelease. It is
	// persisted with the plan so recovery never has to reinterpret legacy
	// release status after source files have been archived.
	ReleaseStatus string
	// DecidedStatus is the owner's decided lifecycle for an Objective or Task
	// whose V1 status migration does not recognise (I-067). It is persisted
	// with the plan so conversion and recovery use the decision rather than
	// reinterpreting the source. Empty when the source status is used.
	DecidedStatus string
	// DependsOn lists the global Task IDs (only meaningful for TargetTask)
	// this target's converted depends_on will name. A dependency on
	// archived, completed work never appears here; see LegacyPrerequisite.
	DependsOn []string
	// DuplicateOfGlobalID is set only for a TargetIssue planned from a
	// `duplicate` finding whose canonical finding also converts: the
	// allocated I-### of that canonical Issue. Empty for every other target,
	// including a duplicate finding whose canonical is archive-only (that
	// duplicate is archived too; see planFinding).
	DuplicateOfGlobalID string
}

// InstallPath returns t's actual filename relative to the .savepoint root,
// including the ".md" extension every V2 record family's discovery layout
// requires. TargetPath itself omits it for TargetTask — task identity
// allocation only needs the directory-qualified destination, not a filename
// — so apply and the manifest's recorded identity map resolve the real
// on-disk path through here rather than each re-deriving the same suffix
// rule. Objective and Issue TargetPath values are already complete
// filenames and pass through unchanged.
func (t PlannedTarget) InstallPath() string {
	if t.Kind == TargetTask && !strings.HasSuffix(t.TargetPath, ".md") {
		return t.TargetPath + ".md"
	}
	return t.TargetPath
}

// ArchiveEntry is one V1 source preserved byte-for-byte under .savepoint/archive/v1/
// instead of being converted, because it is settled history (a done Task, a
// closed epic, a resolved defect, a verified or waived finding, a duplicate
// finding whose canonical is itself archived) or because it matched no known
// role at all.
type ArchiveEntry struct {
	SourcePath  string
	ArchivePath string
	// SourceSHA256 is populated for release PRDs so the manifest can state
	// the exact archived source hash next to its accountable archive mapping.
	// Other archive entries continue to use ManifestSource for the shared
	// inventory hash table.
	SourceSHA256 string
	Role         Role
	Legacy       *LegacyKey // nil when the source carries no legacy identity (e.g. an unclassified file)
	// CandidateCommands lists command-looking lines convert_docs.go found in
	// an archived Health-Check.md body, for a future preview to report as an
	// owner decision. It is nil for every other role: nothing ever writes
	// this list into config.yml, so no quality_gates entry is ever inferred
	// from prose.
	CandidateCommands []string
}

// DocumentKind names one of the small, fixed set of top-level project
// documents migration relocates or rewrites in place, distinct from the
// Objective/Task/Issue record families TargetKind names.
type DocumentKind string

const (
	// DocumentIdea is .savepoint/PRD.md becoming .savepoint/Idea.md: a
	// byte-preserving relocation, never a re-render. The original PRD.md is
	// also archived (see ArchiveEntry for the same source path), so the V1
	// authorship survives twice: once live, once as history.
	DocumentIdea DocumentKind = "idea"
	// DocumentRouter is .savepoint/router.md rewritten in place: only the
	// "## Current state" YAML anchor's vocabulary changes; every other byte
	// of authored prose is preserved.
	DocumentRouter DocumentKind = "router"
)

// PlannedDocument is one top-level project document migration relocates or
// rewrites, as opposed to converting into a new Objective/Task/Issue record.
// SourcePath is project-relative, matching LegacyKey.Path and
// ArchiveEntry.SourcePath. TargetPath is .savepoint-root-relative, matching
// PlannedTarget.TargetPath.
type PlannedDocument struct {
	Kind       DocumentKind
	SourcePath string
	TargetPath string
}

// LegacyPrerequisite records a dependency from a newly allocated active Task
// onto archived, completed V1 work that received no V2 identity. It is the
// mechanism T002's "no V2 schema change" boundary requires: the fact is
// recorded and resolvable, but it is never a fabricated depends_on entry and
// never a fabricated Check.
type LegacyPrerequisite struct {
	Task        string // the new T-### that named the archived work as a dependency
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
	// ConflictDestinationExists means a path this plan intends to create (a
	// converted record, Idea.md, or an archive entry) already exists on
	// disk. A single path cannot carry two fates — a live file plus a
	// planned create — so Apply refuses via this named Conflict at plan
	// time rather than fail later with an internal invariant error and no
	// way out (WriteStaged has no "overwrite an existing create" mode by
	// design: ActionCreate destinations are supposed to be freshly
	// allocated paths nothing else could legitimately occupy).
	ConflictDestinationExists ConflictKind = "destination_exists_conflict"
)

// Conflict is one named reason Plan (or a later apply) must not proceed
// without owner attention.
type Conflict struct {
	Kind   ConflictKind
	Path   string
	Detail string
}

// AmbiguityKind, Ambiguity, and Decisions live in decisions.go, alongside the
// stable ID derivation and the blocking-versus-advisory classification: this
// file only detects ambiguities as it walks the project.

// Clock and OperationIDSource are injected so issue dates, converted bytes,
// and preview output stay byte-identical across repeated test runs.
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

	Targets []PlannedTarget
	// GoalSelection records the single live Goal the converted router selects.
	GoalSelection RouterGoalSelection
	Documents     []PlannedDocument
	Archives      []ArchiveEntry
	Prereqs       []LegacyPrerequisite
	WaivedRefs    []WaivedReference
	Conflicts     []Conflict
	Ambiguities   []Ambiguity

	// Appliable is false exactly when at least one blocking Ambiguity is
	// unresolved; UnresolvedBlockingIDs names every one of them, sorted, so a
	// refusal to apply can name concrete IDs rather than say only that
	// ambiguities exist. It does not fold in Conflicts, which block apply for
	// an unrelated reason and are reported separately.
	Appliable             bool
	UnresolvedBlockingIDs []string

	// Sources is the raw inventory Plan computed, retained so a manifest can
	// be built from the same read without walking the project a second time.
	Sources []SourceFile
}

// RouterGoalSelection is the plan's durable decision about the V2 router's
// live Goal. SourceRelease is the V1 release named by the router, when any;
// GoalID is empty while a blocking release-lifecycle decision is unresolved.
// Generated is true when migration had to create a continuation Goal.
type RouterGoalSelection struct {
	SourceRelease string
	GoalID        string
	Generated     bool
	Reason        string
}

const (
	migrationsDirRel = ".savepoint/migrations"
	manifestFileName = "v1-to-v2.yml"
	v2ObjectivesDir  = "objectives"
	v2ObjectiveFile  = "Objective.md"
	v2TasksDir       = "tasks"
	v2IssuesDir      = "issues"
	v2ReleasesDir    = "releases"

	ideaTargetPath   = "Idea.md"
	routerTargetPath = "router.md"
)

// preservedInPlaceDocs names project documents V2 keeps at their V1 path and
// bytes, untouched by migration, even though the frozen Role vocabulary (see
// TestClassify_roleVocabularyIsComplete) has no dedicated role for them, so
// Classify reports RoleUnclassified for them. Adding a new Role would widen
// that frozen vocabulary just to say "don't touch this," so the exemption
// from archiveRemainingRoles's unclassified sweep is by exact path instead.
// Design.md needs no entry here: RoleArchitecture already isn't in that
// sweep's role switch.
var preservedInPlaceDocs = map[string]bool{
	".savepoint/Guardrails.md":      true,
	".savepoint/visual-identity.md": true,
}

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
	if err := validateDecisions(decisions, b.ambiguities); err != nil {
		return nil, err
	}

	var unresolvedBlocking []string
	for _, a := range b.ambiguities {
		if a.Blocking && !a.Resolved {
			unresolvedBlocking = append(unresolvedBlocking, a.ID)
		}
	}

	return &ConversionPlan{
		GeneratedAt:           now(),
		OperationID:           newOperationID(),
		Targets:               b.targets,
		GoalSelection:         b.goalSelection,
		Documents:             b.documents,
		Archives:              b.archives,
		Prereqs:               b.prereqs,
		WaivedRefs:            b.waivedRefs,
		Conflicts:             checkDestinationCollisions(rootAbs, b.targets, b.documents, b.archives),
		Ambiguities:           b.ambiguities,
		Appliable:             len(unresolvedBlocking) == 0,
		UnresolvedBlockingIDs: unresolvedBlocking,
		Sources:               sources,
	}, nil
}

// checkDestinationCollisions reports a named Conflict for every path this
// plan intends to *create* that already exists on disk: a converted record,
// a created document (Idea.md), or an archive entry. A replaced document
// (router.md) is deliberately excluded — its destination is expected to
// already exist. This is a read-only check (os.Lstat only), so it costs Plan
// nothing of its write-free guarantee, and it is what turns "a stray file
// happens to sit where migration wants to create one" from an apply-time
// internal invariant crash with no way out into a named, reviewable plan-time
// refusal the owner can act on before anything is touched.
func checkDestinationCollisions(root string, targets []PlannedTarget, documents []PlannedDocument, archives []ArchiveEntry) []Conflict {
	var conflicts []Conflict
	check := func(relPath string) {
		abs := filepath.Join(root, filepath.FromSlash(relPath))
		if _, err := os.Lstat(abs); err == nil {
			conflicts = append(conflicts, Conflict{
				Kind: ConflictDestinationExists,
				Path: relPath,
				Detail: fmt.Sprintf("migration plans to create %s, but a file already exists there; "+
					"move or remove it before migrating rather than let it be silently overwritten or leave the operation with no way to finish", relPath),
			})
		}
	}

	for _, t := range targets {
		check(filepath.ToSlash(filepath.Join(".savepoint", t.InstallPath())))
	}
	for _, d := range documents {
		if d.Kind == DocumentRouter {
			continue
		}
		check(filepath.ToSlash(filepath.Join(".savepoint", d.TargetPath)))
	}
	for _, a := range archives {
		check(a.ArchivePath)
	}

	sort.Slice(conflicts, func(i, j int) bool { return conflicts[i].Path < conflicts[j].Path })
	return conflicts
}

func indexSourcesByPath(sources []SourceFile) map[string]SourceFile {
	byPath := make(map[string]SourceFile, len(sources))
	for _, f := range sources {
		byPath[f.Path] = f
	}
	return byPath
}

// idAllocator assigns the next unused O/T/I/R/G number in allocation order.
// Archived-only records never call allocate, so they reserve nothing —
// exactly the "never reused, never reserved for archive-only work" rule.
type idAllocator struct {
	next map[string]int
}

func newIDAllocator() *idAllocator {
	return &idAllocator{next: map[string]int{"O": 1, "T": 1, "I": 1, "R": 1, "G": 1}}
}

func (a *idAllocator) allocate(prefix string) string {
	n := a.next[prefix]
	for {
		id := fmt.Sprintf("%s-%03d", prefix, n)
		if !a.reserved(id) {
			a.next[prefix] = n + 1
			return id
		}
		n++
	}
}

func (a *idAllocator) reserved(id string) bool {
	return a.next["reserved:"+id] > 0
}

func (a *idAllocator) reserve(id string) {
	if id == "" {
		return
	}
	a.next["reserved:"+id] = 1
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
	documents   []PlannedDocument
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
	issueByFindingID map[string]string // finding ID -> allocated I-### (absent = archived)
	// taskIdentitySeen tracks the first source path to declare a given task
	// ID within one release+epic scope, so a second file declaring the same
	// ID raises AmbiguityDuplicateSourceIdentity instead of silently
	// shadowing or overwriting the first.
	taskIdentitySeen map[string]string // "release/epic/originalID" -> first source path
	// releaseIDs is the one source-qualified mapping from a V1 release
	// directory name to its allocated V2 R-### identity. Objectives and router
	// selections use this map; they never derive an R-### independently.
	releaseIDs    map[string]string
	goalSelection RouterGoalSelection
}

type plannedTaskOutcome struct {
	task        data.Task
	globalID    string // "" when archived
	archivePath string // "" when active
}

func (b *planBuilder) build() error {
	b.taskByLegacyPath = map[string]plannedTaskOutcome{}
	b.issueByFindingID = map[string]string{}
	b.taskIdentitySeen = map[string]string{}
	b.releaseIDs = map[string]string{}

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
	for _, release := range releases {
		// A V1 directory can itself already carry an R###-shaped name. Reserve
		// that source identity so allocation cannot silently reuse it as a new
		// V2 identity.
		b.ids.reserve(sourceReleaseIdentity(release.ID))
	}

	// Pass 0: every release and its PRD are resolved before any Objective is
	// planned. This makes the Objective release edge a lookup to the same
	// allocated R-### target rather than a second interpretation of the source.
	for _, release := range releases {
		if err := b.planRelease(release.ID); err != nil {
			return err
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
	if err := b.planRouterGoalSelection(); err != nil {
		return err
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

	// Pass 4: top-level project documents (PRD/Idea, router, Health-Check).
	if err := b.planDocuments(); err != nil {
		return err
	}

	// Pass 5: immutable audit history (runs, prompt, register) and epic
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

// planRouterGoalSelection keeps the live Goal that contains the selected work
// when possible. If active work belongs only to historical Goals, it plans a
// continuation and moves those active Objectives into it.
func (b *planBuilder) planRouterGoalSelection() error {
	router, ok := b.byPath[".savepoint/router.md"]
	if !ok {
		return nil
	}
	raw, err := b.readSource(router.Path)
	if err != nil {
		return err
	}
	v1, err := data.NewRouterReader().ReadState(raw)
	if err != nil {
		// Older planning-only callers can supply a skeletal router. The
		// existing conversion path still rejects it before Apply writes; leave
		// selection unset here rather than pretending the router was resolved.
		return nil
	}

	if v1.Release != "" {
		if id := b.releaseIDs[v1.Release]; id != "" {
			if b.isLiveGoal(id) {
				b.goalSelection = RouterGoalSelection{
					SourceRelease: v1.Release,
					GoalID:        id,
					Reason:        fmt.Sprintf("V1 router release %q resolves to live Goal %s", v1.Release, id),
				}
				return nil
			}
		}
		if b.releaseLifecycleDecisionPending(v1.Release) {
			b.goalSelection = RouterGoalSelection{
				SourceRelease: v1.Release,
				Reason:        fmt.Sprintf("V1 router release %q is awaiting an owner lifecycle decision", v1.Release),
			}
			return nil
		}
	}

	selectedGoalID, selectedEpicFound := b.liveGoalForRouterEpic(v1.Release, v1.Epic)
	if selectedEpicFound && b.isLiveGoal(selectedGoalID) {
		b.goalSelection = RouterGoalSelection{
			SourceRelease: v1.Release,
			GoalID:        selectedGoalID,
			Reason:        fmt.Sprintf("V1 router selection maps to existing live Goal %s containing its active Objective", selectedGoalID),
		}
		return nil
	}
	selectedEpicIsHistorical := selectedEpicFound && selectedGoalID != "" && !b.isLiveGoal(selectedGoalID)
	if !selectedEpicIsHistorical {
		if goalID, unique := b.uniqueLiveGoalWithObjectives(); unique {
			b.goalSelection = RouterGoalSelection{
				SourceRelease: v1.Release,
				GoalID:        goalID,
				Reason:        fmt.Sprintf("Selected existing live Goal %s because it contains the converted active Objectives", goalID),
			}
			return nil
		}
		if goalID, unique := b.uniqueLiveGoal(); unique {
			b.goalSelection = RouterGoalSelection{
				SourceRelease: v1.Release,
				GoalID:        goalID,
				Reason:        fmt.Sprintf("Selected the only existing live Goal %s after the V1 router selection did not resolve", goalID),
			}
			return nil
		}
	}

	reason := "The V1 router has no selected release"
	if v1.Release != "" {
		reason = fmt.Sprintf("V1 router release %q does not resolve to a live Goal", v1.Release)
		if id := b.releaseIDs[v1.Release]; id != "" && !b.isLiveGoal(id) {
			reason = fmt.Sprintf("V1 router release %q resolves only to historical Goal %s", v1.Release, id)
		}
	}
	if selectedEpicIsHistorical {
		reason = fmt.Sprintf("V1 router selection points to historical Goal %s; active Objectives move to the continuation Goal", selectedGoalID)
	}

	hasExistingLiveGoal := b.hasSourceLiveGoal()
	goalID := b.ids.allocate("G")
	b.targets = append(b.targets, PlannedTarget{
		Kind:           TargetRelease,
		GlobalID:       goalID,
		TargetPath:     filepath.ToSlash(filepath.Join(v2ReleasesDir, goalID+"-continued-after-migration", "Release.md")),
		ReleaseID:      goalID,
		ReleaseStatus:  string(data.ColumnInProgress),
		Generated:      true,
		GeneratedTitle: continuationGoalTitle,
	})
	b.goalSelection = RouterGoalSelection{
		SourceRelease: v1.Release,
		GoalID:        goalID,
		Generated:     true,
		Reason:        reason,
	}
	if selectedEpicIsHistorical || !hasExistingLiveGoal {
		for i := range b.targets {
			target := &b.targets[i]
			if target.Kind == TargetObjective && !b.isLiveGoal(target.ReleaseID) {
				target.ReleaseID = goalID
			}
		}
	}
	return nil
}

func (b *planBuilder) isLiveGoal(goalID string) bool {
	for _, target := range b.targets {
		if target.Kind == TargetRelease && target.GlobalID == goalID {
			return !target.Generated && target.ReleaseStatus == string(data.ColumnInProgress)
		}
	}
	return false
}

func (b *planBuilder) liveGoalForRouterEpic(release, epic string) (string, bool) {
	if epic == "" {
		return "", false
	}
	for _, target := range b.targets {
		if target.Kind == TargetObjective && target.Legacy.Release == release && target.Legacy.Epic == epic {
			return target.ReleaseID, true
		}
	}

	var goalID string
	matches := 0
	for _, target := range b.targets {
		if target.Kind != TargetObjective || target.Legacy.Epic != epic {
			continue
		}
		matches++
		goalID = target.ReleaseID
	}
	return goalID, matches == 1
}

func (b *planBuilder) uniqueLiveGoalWithObjectives() (string, bool) {
	var goalID string
	for _, target := range b.targets {
		if target.Kind != TargetObjective || !b.isLiveGoal(target.ReleaseID) {
			continue
		}
		if goalID != "" && goalID != target.ReleaseID {
			return "", false
		}
		goalID = target.ReleaseID
	}
	return goalID, goalID != ""
}

func (b *planBuilder) uniqueLiveGoal() (string, bool) {
	var goalID string
	for _, target := range b.targets {
		if target.Kind != TargetRelease || target.Generated || target.ReleaseStatus != string(data.ColumnInProgress) {
			continue
		}
		if goalID != "" && goalID != target.GlobalID {
			return "", false
		}
		goalID = target.GlobalID
	}
	return goalID, goalID != ""
}

func (b *planBuilder) hasSourceLiveGoal() bool {
	for _, target := range b.targets {
		if target.Kind == TargetRelease && !target.Generated && target.ReleaseStatus == string(data.ColumnInProgress) {
			return true
		}
	}
	return false
}

func (b *planBuilder) releaseLifecycleDecisionPending(release string) bool {
	releaseDir := filepath.ToSlash(filepath.Join(".savepoint", "releases", release))
	for _, ambiguity := range b.ambiguities {
		if !ambiguity.Blocking || ambiguity.Resolved {
			continue
		}
		if ambiguity.Kind != AmbiguityReleaseCompletion && ambiguity.Kind != AmbiguityReleaseLifecycle {
			continue
		}
		if strings.HasPrefix(ambiguity.Path, releaseDir+"/") {
			return true
		}
	}
	return false
}

func sourceReleaseIdentity(name string) string {
	if !strings.HasPrefix(name, "R") {
		return ""
	}
	digits := 1
	for digits < len(name) && name[digits] >= '0' && name[digits] <= '9' {
		digits++
	}
	if digits-1 < 3 {
		return ""
	}
	return name[:digits]
}

type releaseRawFrontmatter struct {
	Name   string `yaml:"name"`
	Title  string `yaml:"title"`
	Status string `yaml:"status"`
}

// planRelease allocates one stable R-### for one source release and gives its
// PRD both a live Release destination and an exact-byte archive destination.
// A missing PRD in the old container-only shape has no promise to convert.
// A duplicated PRD is a blocking source-identity ambiguity: the planner does
// not invent a promise or choose one authored document without an owner
// decision.
func (b *planBuilder) planRelease(release string) error {
	releaseDir := filepath.ToSlash(filepath.Join(".savepoint", "releases", release))
	prds := b.findRolesUnder(RoleReleasePRD, releaseDir)
	var source SourceFile
	switch len(prds) {
	case 0:
		// Some V1 projects used the release directory only as a container for
		// Epic records and never authored a release PRD. There is no promise
		// source to convert in that shape; leave the container untouched and
		// keep the existing E45 migration behavior. A directory with multiple
		// PRDs, by contrast, is an explicit source collision and blocks.
		return nil
	case 1:
		source = prds[0]
	default:
		paths := make([]string, 0, len(prds))
		for _, prd := range prds {
			paths = append(paths, prd.Path)
		}
		ambiguityID := string(AmbiguityDuplicateReleaseSource) + ":" + releaseDir
		b.addAmbiguity(AmbiguityDuplicateReleaseSource, releaseDir,
			fmt.Sprintf("release %s has multiple release PRD sources: %s; migration cannot choose one without an owner decision", release, strings.Join(paths, ", ")), paths)
		if decision, ok := b.decisions[ambiguityID]; ok {
			for _, prd := range prds {
				if prd.Path == decision.Value {
					source = prd
					break
				}
			}
		}
		if source.Path == "" {
			return nil
		}
	}

	content, err := b.readSource(source.Path)
	if err != nil {
		return err
	}
	fm, _, err := data.SplitFrontmatterBody(content)
	if err != nil {
		return fmt.Errorf("parse release PRD %s: %w", source.Path, err)
	}
	var raw releaseRawFrontmatter
	if err := yaml.Unmarshal([]byte(fm), &raw); err != nil {
		return fmt.Errorf("parse release PRD %s: %w", source.Path, err)
	}

	status, ok := b.resolveReleaseStatus(source.Path, raw.Status)
	if !ok {
		return nil
	}

	releaseID := b.ids.allocate("R")
	b.releaseIDs[release] = releaseID
	legacy := LegacyKey{Release: release, Path: source.Path, OriginalID: release}
	b.targets = append(b.targets, PlannedTarget{
		Kind:          TargetRelease,
		GlobalID:      releaseID,
		Legacy:        legacy,
		TargetPath:    filepath.ToSlash(filepath.Join(v2ReleasesDir, releaseID+"-"+slugOf(release), "Release.md")),
		ReleaseID:     releaseID,
		ReleaseStatus: status,
	})
	b.archives = append(b.archives, ArchiveEntry{
		SourcePath:   source.Path,
		ArchivePath:  archivePathFor(source.Path),
		SourceSHA256: source.SHA256,
		Role:         RoleReleasePRD,
		Legacy:       &legacy,
	})
	return nil
}

// resolveReleaseStatus maps legacy release lifecycle to the deliberately
// conservative V2 migration outcome. Any non-settled legacy status becomes
// in_progress so the migrated Release cannot be mistaken for completed work.
// A legacy `audited` disposition is not completion evidence by itself and
// therefore requires an explicit owner choice before a target is planned.
func (b *planBuilder) resolveReleaseStatus(path, raw string) (string, bool) {
	switch raw {
	case "planned", "in_progress", "todo":
		return string(data.ColumnInProgress), true
	case "done", "complete", "completed":
		return string(data.ColumnDone), true
	case "audited":
		id := string(AmbiguityReleaseCompletion) + ":" + path
		b.addAmbiguity(AmbiguityReleaseCompletion, path,
			fmt.Sprintf("release PRD %s records audited without a typed completion decision; choose whether the release is historical or active", path),
			[]string{"in_progress", "done"})
		if decision, ok := b.decisions[id]; ok {
			if decision.Value == "done" {
				return string(data.ColumnDone), true
			}
			if decision.Value == "in_progress" {
				return string(data.ColumnInProgress), true
			}
		}
		return "", false
	default:
		id := string(AmbiguityReleaseLifecycle) + ":" + path
		b.addAmbiguity(AmbiguityReleaseLifecycle, path,
			fmt.Sprintf("release PRD %s has unrecognized status %q; choose an active or historical V2 lifecycle", path, raw),
			[]string{"in_progress", "done"})
		if decision, ok := b.decisions[id]; ok {
			if decision.Value == "done" || decision.Value == "in_progress" {
				return decision.Value, true
			}
		}
		return "", false
	}
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

	decidedStatus := ""
	if !recognized {
		b.addAmbiguity(AmbiguityUnrecognizedLifecycle, sf.Path,
			fmt.Sprintf("epic %s/%s has unrecognized status %q; no Objective or Task under it can be planned until an owner decision resolves it", release, epic, raw.Status),
			epicStatusChoices())
		decided, ok := b.decidedLifecycle(sf.Path)
		if !ok {
			return nil
		}
		status, decidedStatus = decided, decided
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

	goalID := b.releaseIDs[release]
	if goalID == "" {
		if b.releaseLifecycleDecisionPending(release) {
			return nil
		}
		return fmt.Errorf("plan Objective for V1 epic %s/%s: V1 release %q has no resolvable V2 Goal", release, epic, release)
	}

	objectiveID := b.ids.allocate("O")
	b.targets = append(b.targets, PlannedTarget{
		Kind:          TargetObjective,
		GlobalID:      objectiveID,
		Legacy:        legacy,
		TargetPath:    filepath.ToSlash(filepath.Join(v2ObjectivesDir, objectiveID+"-"+slugOf(epic), v2ObjectiveFile)),
		ReleaseID:     goalID,
		DecidedStatus: decidedStatus,
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
		decidedStatus := ""
		if !recognized {
			b.addAmbiguity(AmbiguityUnrecognizedLifecycle, rel,
				fmt.Sprintf("task %s has unrecognized status %q; it cannot be planned as active or archived until an owner decision resolves it", task.ID, raw),
				taskStatusChoices())
			decided, ok := b.decidedLifecycle(rel)
			if !ok {
				continue
			}
			status, decidedStatus = decided, decided
			task.Column = data.ColumnType(decided)
		}

		identityScope := release + "/" + epic + "/" + task.ID
		reuseID := ""
		if firstPath, seen := b.taskIdentitySeen[identityScope]; seen {
			b.addAmbiguityDiscriminated(AmbiguityDuplicateSourceIdentity, rel, firstPath,
				fmt.Sprintf("task id %s is declared by both %s and %s; allocation cannot tell which one it actually names", task.ID, firstPath, rel),
				[]string{"keep_first", "keep_second"})
			decision, ok := b.decisions[string(AmbiguityDuplicateSourceIdentity)+":"+rel+"#"+firstPath]
			if !ok {
				continue
			}
			if decision.Value == "keep_first" {
				b.archives = append(b.archives, ArchiveEntry{
					SourcePath:  rel,
					ArchivePath: archivePathFor(rel),
					Role:        RoleTask,
					Legacy:      &legacy,
				})
				continue
			}
			// keep_second: archive the first declaration and plan this one
			// in its place, under the identity the first was given.
			reuseID = b.retireTaskDeclaration(firstPath)
		}
		b.taskIdentitySeen[identityScope] = rel

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

		taskID := reuseID
		if taskID == "" {
			taskID = b.ids.allocate("T")
		}
		b.taskByLegacyPath[rel] = plannedTaskOutcome{task: *task, globalID: taskID}
		b.targets = append(b.targets, PlannedTarget{
			Kind:          TargetTask,
			GlobalID:      taskID,
			Legacy:        legacy,
			TargetPath:    filepath.ToSlash(filepath.Join(v2ObjectivesDir, objectiveID+"-"+slugOf(epic), v2TasksDir, taskID+"-"+slugOf(task.ID))),
			DecidedStatus: decidedStatus,
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

// retireTaskDeclaration withdraws an earlier task declaration that an owner
// decision replaced (keep_second): its planned target, if any, is removed and
// the source is archived instead. It returns the Task identity the retired
// declaration held, or "" when it was already archive-only.
func (b *planBuilder) retireTaskDeclaration(path string) string {
	outcome := b.taskByLegacyPath[path]
	delete(b.taskByLegacyPath, path)
	if outcome.globalID == "" {
		return ""
	}
	kept := b.targets[:0]
	for _, target := range b.targets {
		if target.Kind == TargetTask && target.Legacy.Path == path {
			legacy := target.Legacy
			b.archives = append(b.archives, ArchiveEntry{
				SourcePath:  path,
				ArchivePath: archivePathFor(path),
				Role:        RoleTask,
				Legacy:      &legacy,
			})
			continue
		}
		kept = append(kept, target)
	}
	b.targets = kept
	return outcome.globalID
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
				b.addAmbiguityDiscriminated(AmbiguityMissingDependencyTarget, rel, ref,
					fmt.Sprintf("task %s depends_on %q, which does not resolve to any known task", outcome.task.ID, ref),
					missingDependencyChoices())
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
				b.addAmbiguityDiscriminated(AmbiguityMissingDependencyTarget, rel, ref,
					fmt.Sprintf("task %s depends_on %q, which resolved but was never planned", outcome.task.ID, ref),
					missingDependencyChoices())
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
				fmt.Sprintf("finding %s is marked duplicate of %q, which names no known finding in this project", finding.ID, finding.DuplicateOf),
				unresolvedNarrativeFindingChoices())
			decision, ok := b.decisions[string(AmbiguityUnresolvedNarrativeFind)+":"+path]
			if !ok {
				return nil
			}
			if decision.Value == "archive" {
				b.archives = append(b.archives, ArchiveEntry{
					SourcePath:  path,
					ArchivePath: archivePathFor(path),
					Role:        RoleFinding,
					Legacy:      &legacy,
				})
				return nil
			}
			// treat_as_original: the finding converts as its own open Issue.
			issueID := b.ids.allocate("I")
			b.issueByFindingID[finding.ID] = issueID
			b.targets = append(b.targets, PlannedTarget{
				Kind:          TargetIssue,
				GlobalID:      issueID,
				Legacy:        legacy,
				TargetPath:    filepath.ToSlash(filepath.Join(v2IssuesDir, issueID+"-"+slugOf(finding.ID)+".md")),
				DecidedStatus: string(data.FindingOpen),
			})
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
			fmt.Sprintf("finding %s has unrecognized status %q", finding.ID, finding.Status),
			findingStatusChoices())
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

// planDocuments plans the small, fixed set of top-level project documents:
// PRD.md relocates to Idea.md (byte-preserving; the original is also
// archived), router.md rewrites its state anchor in place, and
// Health-Check.md is archived with its candidate commands recorded for a
// future preview to report. config.yml, Design.md, Guardrails.md, and
// visual-identity.md are never archived or targeted by this method: their
// absence from plan.Archives and plan.Documents is itself the "preserved in
// place, untouched" outcome apply later relies on.
func (b *planBuilder) planDocuments() error {
	var paths []string
	for path := range b.byPath {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	for _, path := range paths {
		switch Classify(path) {
		case RoleProductPRD:
			b.archives = append(b.archives, ArchiveEntry{
				SourcePath:  path,
				ArchivePath: archivePathFor(path),
				Role:        RoleProductPRD,
			})
			b.documents = append(b.documents, PlannedDocument{
				Kind:       DocumentIdea,
				SourcePath: path,
				TargetPath: ideaTargetPath,
			})

		case RoleRouter:
			b.documents = append(b.documents, PlannedDocument{
				Kind:       DocumentRouter,
				SourcePath: path,
				TargetPath: routerTargetPath,
			})

		case RoleHealthCheck:
			content, err := b.readSource(path)
			if err != nil {
				return err
			}
			b.archives = append(b.archives, ArchiveEntry{
				SourcePath:        path,
				ArchivePath:       archivePathFor(path),
				Role:              RoleHealthCheck,
				CandidateCommands: CandidateHealthCheckCommands(content),
			})
		}
	}
	return nil
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
		if handled[path] || isPreservedMigrationsContent(path) || preservedInPlaceDocs[path] {
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
			if role == RoleUnclassified {
				// Advisory only: the file is archived intact regardless of
				// this ambiguity's resolution, so it never blocks apply. It
				// is named purely so the owner can see what migration could
				// not place, per AC4.
				b.addAmbiguity(AmbiguityUnclassifiedFile, path,
					fmt.Sprintf("%s matched no known V1 file role; it is archived byte-for-byte at %s rather than converted or dropped", path, archivePathFor(path)),
					nil)
			}
		}
	}
}

func (b *planBuilder) readSource(rel string) (string, error) {
	content, err := os.ReadFile(filepath.Join(b.root, filepath.FromSlash(rel)))
	if err != nil {
		return "", fmt.Errorf("read %s: %w", rel, err)
	}
	return string(content), nil
}

// archivePathFor is the deterministic .savepoint/archive/v1/ destination for
// a project-relative V1 source path: the same relative layout, moved under
// .savepoint/archive/v1/, so an archived path is always recoverable by
// inspection alone. It lives inside .savepoint/ — as v2-Design.md specifies —
// so the destination is inventoried, path-confined, and collision-checked
// like every other path the operation writes; a location outside .savepoint/
// would be invisible to the inventory and therefore to every conflict and
// backup guarantee this operation makes. archiveDirName (inventory.go) is
// the shared name Inventory also excludes from its own walk, so this
// directory's own content is never re-inventoried as a new source on resume.
func archivePathFor(sourcePath string) string {
	return filepath.ToSlash(filepath.Join(".savepoint", archiveDirName, "v1", sourcePath))
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
	matches := b.findRolesUnder(role, dir)
	if len(matches) == 1 {
		return matches[0], true
	}
	return SourceFile{}, false
}

func (b *planBuilder) findRolesUnder(role Role, dir string) []SourceFile {
	prefix := dir + "/"
	var matches []SourceFile
	for path, sf := range b.byPath {
		if !strings.HasPrefix(path, prefix) {
			continue
		}
		if strings.Contains(path[len(prefix):], "/") {
			continue // nested further (e.g. tasks/); not directly under dir
		}
		if Classify(path) == role {
			matches = append(matches, sf)
		}
	}
	sort.Slice(matches, func(i, j int) bool { return matches[i].Path < matches[j].Path })
	return matches
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

// The choices* helpers below are the concrete, enumerable resolutions T007's
// AC2 requires each ambiguity to state. Choices resolve the *ambiguity Plan
// raised* (which status the owner means, whether to drop a reference, which
// of two duplicate declarations is canonical, how to treat an unresolved
// duplicate finding) — none of them describe the deeper V1 source edit that
// naturally follows and belongs to the owner, not to migration.

func epicStatusChoices() []string {
	statuses := data.CanonicalEpicStatuses()
	choices := make([]string, len(statuses))
	for i, s := range statuses {
		choices[i] = string(s)
	}
	return choices
}

func taskStatusChoices() []string {
	statuses := data.CanonicalTaskStatuses()
	choices := make([]string, len(statuses))
	for i, s := range statuses {
		choices[i] = string(s)
	}
	return choices
}

func findingStatusChoices() []string {
	return []string{
		string(data.FindingOpen), string(data.FindingTriaged), string(data.FindingMapped),
		string(data.FindingInProgress), string(data.FindingFixed), string(data.FindingVerified),
		string(data.FindingDeferred), string(data.FindingOwnerDecision), string(data.FindingWaived),
		string(data.FindingDuplicate),
	}
}

// missingDependencyChoices is deliberately a single choice: migration never
// invents a replacement dependency target, so the only decision it can offer
// is to drop the reference. Naming a real target is a V1 source edit,
// followed by a rerun that previews the corrected reference cleanly.
func missingDependencyChoices() []string {
	return []string{"drop_dependency"}
}

// unresolvedNarrativeFindingChoices resolves a `duplicate` finding whose
// duplicate_of names no known finding: either it converts as its own
// Issue (treat_as_original) or it is archived like a settled duplicate
// (archive). Neither choice guesses which existing finding the owner meant.
func unresolvedNarrativeFindingChoices() []string {
	return []string{"treat_as_original", "archive"}
}
