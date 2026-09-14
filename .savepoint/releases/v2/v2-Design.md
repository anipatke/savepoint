---
type: release-design
status: planned
release: v2
---

# V2 design proposal

## 1. Scope and evidence

This proposal applies the confirmed refactor prompt. It describes intended V2 behavior, not current runtime guarantees. Review covered source boundaries and representative historical records; no independent code audit or runtime verification is claimed here. The release PRD owns delivery order; this document owns technical contracts.

| Current evidence | Implication for V2 |
|---|---|
| `internal/data/discover.go`, `dependency.go`, `parser.go` traverse releases/epics and infer ownership from IDs; dependency matching includes release scope. | Replacing folder names alone is insufficient. Introduce one V2 project index with globally stable identities and explicit Objective ownership. |
| `internal/data/lifecycle.go` heals unknown status to planned and aliases to canonical values; `parser.go` uses `objective` as a title fallback. | Migration must consume raw metadata, not only healed model values. The V2 `objective` reference needs a versioned schema boundary. |
| `internal/data/write.go` preserves YAML nodes but normalizes line endings and uses direct writes. `internal/init/write.go` can fall back from rename to truncating copy. | Existing helpers are not proof of exact-byte or multi-file migration safety. Preserve raw source bytes and implement tested replacement/recovery semantics. |
| `internal/board/transitions.go` checks dependency status, not independent Check evidence; `interfaces.go` and `watch.go` expose release/epic, defect, and audit sets. | Centralize V2 gate decisions in data, then change consumers together. Do not enforce rules only in the board. |
| `internal/data/audit_finding.go` has ten finding states; `audit_run.go` stores immutable runs; register summaries are separately maintained. | Keep identity, proof, scope, dispositions, and history; derive lists/counts rather than maintaining another register. |
| `internal/init/upgrade.go`, `manifest.go`, `agents.go`, and failure/lifecycle tests implement ownership and recoverable upgrades. | Preserve provenance and managed-guide boundaries. Schema migration must be an explicit operation distinct from asset refresh. |
| `internal/doctor/checks.go` expects releases; config keys include `quality_gates`, and gates already have timeout handling. | Reuse diagnostics and command execution rather than inventing another health framework or command syntax. |
| `agent-skills/references/audit-method.md` contains scope, coverage, side-effect, materiality, and re-audit convergence mechanics. | Preserve rigor in shared references, while simplifying the public skill surface and changing checker write authority deliberately. |
| `templates/project/.savepoint` includes Concept, Health-Check, release and audit scaffolds. The live project lacks Health-Check and `.savepoint/audit/`. | Migration fixtures need both absent and populated optional artifacts. Live project absence is not a defect. |
| `main.go`, `bin/savepoint.js`, `internal/buildtool/main.go`, `package.json`, CI/publish workflows ship Go binaries through npm. | Preserve thin command dispatch and platform packaging. Do not rely on stale TypeScript sections in current Design.md as implementation evidence. |

V1 tests such as `TestLifecycle_upgradeHonoursOwnership`, `TestUpgrade_writeFailureLeavesRecoverableState`, and parser alias tests are useful starting evidence to extend, not assertions that V2 already passes. Historical E40 detail describes pre-fix problems; its Audit and source establish the later outcome.

## 2. File model and ownership

```text
project/
├── AGENTS.md                         # managed routing block; surrounding prose preserved
├── agent-skills/                     # four public skills plus shared method
└── .savepoint/
    ├── Idea.md                       # product intent, editable by planner/owner
    ├── Design.md                     # current architecture, planner/executor reconciliation
    ├── Guardrails.md                 # durable project policy
    ├── config.yml                    # schema_version: 2, theme, quality_gates
    ├── router.md                     # phase and selected Objective/Task, not completion truth
    ├── objectives/
    │   └── O001-project-recovery/
    │       ├── Objective.md
    │       └── tasks/T001-resume.md
    ├── checks/C001.md                # immutable completed evaluation, local provenance
    ├── issues/I001-resume-defect.md   # mutable follow-up plus append-only history
    ├── archive/v1/                   # byte-preserved historical source tree
    ├── migrations/v1-to-v2.yml       # source hashes, ID/reference map, operation recovery
    └── .upgrade-manifest.yml         # installed asset provenance, separate from schema version
```

`archive/` and `migrations/` appear only after migration. Optional `visual-identity.md` and project verification procedures remain ordinary referenced project files. No default Concept, Health-Check, release PRD, audit register, or findings folder is needed in new V2 projects. Empty new projects may have no Objective yet; a Task can never exist without one.

Project records are authoritative markdown/YAML; directory indexes, issue counts, board columns, and Next are derived. Unknown fields and bodies survive managed edits; migration additionally archives original exact bytes. Filesystem paths are project-relative, normalized for storage and joined with `filepath` at I/O boundaries. Reject traversal, duplicate identity, and symlink escapes rather than following untrusted paths outside the project.

## 3. Identity, schema detection, and references

- `config.yml` owns `schema_version: 2`. Absence indicates legacy input only during transition/migration; unknown explicit versions fail with a named diagnostic. Never confuse npm/package versions, release metadata, or upgrade-manifest version with project schema.
- Global IDs use `O`, `T`, `I`, or `C` plus at least three digits. A title/path/Objective change does not change an ID. Discovery indexes explicit IDs; path mismatch is diagnostic, never silent reassignment.
- Allocate the next unused ID from active records and migration reservations; never reuse archived or deleted identities. Single-user planning records the reservation in the created record; concurrent duplicate creation is diagnosed, not silently merged.
- Task `objective` is an `O###` reference; `title` and Outcome supply display language. Every task belongs to exactly one Objective. Optional `release` is filtering/packaging metadata.
- Dependencies are records `{task: T###, requires: clear|accepted}`; `clear` is the default. Objective dependencies are `O###` IDs and require integration clearance. Detect missing targets, cycles, and self-dependencies.
- Legacy mapping keys include source path/release/epic and original ID. `T001` alone is not a migration key. Historical prerequisites use typed archive references, resolved through the migration manifest with preserved original completion/waiver evidence. Display them as legacy prerequisites, never new CLEAR results.

## 4. Task lifecycle and authority

Retain Task `status: planned|in_progress|done` to avoid turning board columns into five competing completion concepts. While `in_progress`, retain `stage: build|test|audit`; expose the last as Check in the UI. Omit stage otherwise. V2 lifecycle rules remain owned by `internal/data` and are distinct from current V1 ownership while the release is built.

| Event | Recorded effect | Gate / authority |
|---|---|---|
| Start | planned → in_progress, stage build | Executor; ready plan, satisfied Task and Objective dependencies. |
| Verify implementation | build → test → audit | Executor records AC evidence and required command results; audit means ready for Check, not passed. |
| Replan | Keep current status/stage; set `replan.reason` with handoff evidence | Executor pauses; planner resolves decisions and updates plan before clearing the flag. Preserve partial work. |
| Failed Check | New immutable NEEDS WORK Check; link from task | Checker records Issues. Executor resumes repair at build within the unfinished task. |
| Clear Check | New immutable CLEAR Check; link from task | Fresh checker session, required technical evidence, no unexcepted material blocker. |
| Complete technical Task | in_progress → done, remove stage | Checker may close only with current clearance and no required owner validation. |
| Complete owner-validated Task | in_progress → done, remove stage | Owner acceptance of the checked outcome, or checker finalization after previously recorded acceptance still applies. |
| Clearance stale/unknown | Derived completion invalid/needs attention; preserve previous records | No silent status rewrite. Fresh Check required; new repair task if completed work is defective. |

Task stores `last_check`, `owner_validation.required`, and optional `owner_validation.accepted_check`. Acceptance names the Check whose outcome the owner accepted; material changes require renewed acceptance. There is no separate mutable `checked: true` flag. Clearance freshness is a recorded assessment with scope/revision basis; lack of an assessment is unknown, not current.

Before any completion control writes, resolve current evidence and prerequisites from files again. Rejected actions explain the missing requirement. Explicit owner exceptions specify requirement IDs, reason, owner provenance, time, and the affected Check/outcome. They never fabricate a CLEAR result or silently waive a new revision. Display completion by exception separately. Data/doctor detect inconsistent direct edits but do not claim authentication or prevention of external writes.

## 5. Check records and freshness

One immutable evaluation file per run; reruns get a new C ID and `supersedes`. Task and Objective Checks share a schema with `scope.kind: task|objective` and a target ID. Proposed minimum fields:

```yaml
id: C001
scope: {kind: task, id: T001}
result: CLEAR
checked_by: {role: checker, session: review-001}
executed_session: build-001
checked_at: '2026-09-14T00:00:00Z'
reviewed:
  base_commit: optional-base-sha
  head_commit: optional-head-sha
  files: [] # required populated path/content-hash or explicit absent entries in a real Check
  dependencies: [] # relevant source/config/lockfiles and requirement inputs
issues: []
supersedes: null
```

The body records outcome coverage, tests and exit results, negative/boundary probes, applicable Guardrails, owner validation still needed, and nonblocking observations. A record without sufficient scope/evidence cannot support completion. Git is useful provenance, not required as a database; exact file hashes or a named preserved snapshot can cover uncommitted or non-Git work. Agents collect this evidence explicitly in V2.0.

Task/Objective evidence records `freshness: {state: current|stale|unknown, check: C###, assessed_by: ..., assessed_at: ..., basis: ...}`. Agents compare the relevant implementation and requirement inputs to the recorded scope. New relevant files/dependency changes invalidate prior clearance, even if old reviewed files did not change. Administrative edits to evidence fields alone do not invalidate code review; changes to AC, Design constraints, or implementation do. No automatic Git hooks or whole-repository hash scan is required initially.

The board and resume say **recorded current as of ...**, not that the current tree has automatically been verified. Known stale/unknown evidence blocks normal completion. New sessions must reassess before relying on it. This is an explicit trust boundary of file-first, agent-agnostic operation.

Fresh session means independent from the executor conversation; model names are optional and same model is allowed. Checker writes Check/Issue and authorized evaluation/closure metadata only, never repairs implementation, rewrites AC to match a result, or updates Design as remediation. Corrections return to planner/executor, then a new Check verifies the repair. Preserve scope discipline, coverage reasoning, side-effect checks, materiality, and bounded re-audit convergence in shared methodology.

## 6. Objective and Issue contracts

Objective frontmatter: `id`, `title`, `status: planned|in_progress|done`, `depends_on: [O###]`, optional `release`, `last_check`, and the same freshness/exception evidence structure. Body: Outcome, Why, Success Conditions, Architectural Considerations, Boundaries. Task membership is derived from Task ownership, not a second manually maintained list. An Objective reaches done only after its Tasks meet completion rules and its independent integration Check is current (or an explicit scoped owner exception). Cross-task repair goes through Tasks. Dependent Objectives cannot start from Task-only clearance.

Issue frontmatter: `id`, `title`, `type: defect|drift|guardrail|verification|other`, `status: open|in_progress|resolved`, `source`, `tasks`, `checks`, `guardrail_ids`, optional severity, `resolution`, and `duplicate_of`. Body: Summary, Evidence, Proof Needed, History. Type is descriptive, never sufficient by itself to block a Task.

- Checker explains material blocking against AC, required evidence, or applicable policy. Advisory style is nonblocking. Pending owner acceptance lives on the Task, not a duplicate Issue unless durable follow-up is needed.
- Open → in_progress when repair starts; executor reports repair evidence without closing it. Checker verifies proof and closes. Owner acceptance is needed when declared; accepting unresolved risk requires an explicit owner decision.
- Resolved disposition distinguishes `verified`, `accepted`, and `duplicate`. Accepted is not fixed; duplicate points to canonical Issue and is not proof of repair. Reopen the same ID for the same recurring problem with dated evidence.
- Link out-of-scope repairs to new Tasks. Search existing Issues by symptom, location, violated requirement, and linked work before creating an ID. No probabilistic automatic deduplication service is required.
- History records observations, repair attempts, rechecks, deferrals, and owner decisions. Deferral stays open with reason, not another lifecycle state. Check records replace immutable audit run history; issue listings/counts are derived without a separate register.

## 7. Idea, Design, Guardrails, and planning readiness

Idea template sections: Intent, User, Core Experience, Scope, Out of Scope, Success Criteria. Accept a rough sentence and refine through the planner; do not require a prepared PRD.

Design template sections: Architecture, Components/Codebase Map, Interfaces and Data Flow, Boundaries, Decisions, Current Technical State. Describe implemented reality concisely; Objective deltas distinguish planned changes until reconciliation. Avoid repeating Guardrails or maintaining a second backlog here.

Guardrails contain project-relevant durable constraints after Design is sufficiently mature. Prefer roughly 10–20 substantive rules, retaining more where justified. Stable category IDs reference constraints rather than repeated prose. Policy owns severity and exception authority; tasks/checks reference it. Current repository policy remains authoritative during the build, including FS/DATA/TPL/ARCH/TEST/REL and advisory STYLE rules.

Planned V2 policy changes: replace DATA-05's universal owner-only completion with conditional owner acceptance; extend DATA-02 to V2 semantics while keeping one owner package; retire hard context-byte limits from guidance; retain exact-byte user-content and provenance boundaries. These changes become active only at workflow cutover. Existing FS-06's preflight no-write expectation must be distinguished from explicit recoverable interruption after a valid migration starts; no promise of impossible multi-file atomicity.

Planner readiness for one Objective requires settled interfaces/data ownership, scoped constraints, known dependency outcomes, and a verification approach. Product choices go to the owner; technical readiness does not require owner code review. Unknown implementation approaches become bounded research Tasks with a decision deliverable. Detail only the next Objective. Split tasks with multiple unrelated outcomes or unresolved architectural decisions. Measure required/contextually added reads and replanning frequency rather than a universal <2KB cap.

## 8. Skill and verification model

| Public skill / capability | Reads | Writes | Forbidden / escalation |
|---|---|---|---|
| savepoint-idea / planner | Router, Idea, user intent, targeted existing-project evidence | Idea and routing handoff | No architecture/tasks/code; ask about material product uncertainty. |
| savepoint-design / planner | Idea, Design, Guardrails, current Objective, scoped implementation evidence | Design/Guardrails, Objective, next Objective's detailed Tasks, routing | No production code or wholesale backlog detail; research unresolved approach. |
| savepoint-task / executor | Router, Task, Objective boundaries, scoped files, relevant policy | Scoped implementation, evidence, lifecycle progress, replan handoff | No silent redesign or clearance/acceptance claims; targeted reads logged; material gaps return REPLAN REQUIRED. |
| savepoint-check / checker | Task or Objective scope, source/tests, policy, prior Checks/Issues, shared method | Check, Issues, evaluation metadata and authorized closure | Fresh session; no repairs or changed acceptance criteria; hand remediation back. |

Shared references own checking and Issue reconciliation methodology, not extra public phases. Issue capture is a supporting entry into design/task/check, with Defect still an explicit user term. Task Check is the focused depth; Objective Check includes integration and Design reconciliation. Retain proof and convergence, remove mandatory register bookkeeping.

Keep `quality_gates` as the config key to reuse existing execution and timeouts; add explicit build command support rather than introduce parallel `checks.technical`. Tasks name extra verification; optional project procedures cover reusable manual workflows. Health-Check migration maps commands into config and reusable custom prose into a preserved optional procedure. Universal Quick/Full mechanics go to shared method; release-specific readiness belongs to the release outcome. Absent legacy files need no generated substitute.

## 9. Migration and archival mapping

Use explicit `migrate [dir] --dry-run` and apply operation; asset upgrade never implicitly changes schema. No permanent --v1 runtime. Transitional schema dispatch exists only until consumers move to V2; E50 removes V1 live readers from board/doctor and retains migration input readers and frozen fixtures.

| V1 record | V2 destination / interpretation |
|---|---|
| Project PRD | Idea with authored content preserved; technical material remains available for planner reconciliation. |
| Active release PRDs / epic details | Objective scope and optional release metadata; original bytes archived. |
| Planned/in-progress Task | New global ID; preserve plan/AC/evidence, map Objective/dependencies, no fabricated clearance. |
| Done Task / completed audit | Archive intact; dependency lookup through source-qualified mapping. |
| Completed Tasks but unaudited epic | Keep an active Objective needing integration Check, referencing archived task evidence; do not declare it complete. |
| Unresolved defect | Issue type defect; active repair maps to an existing/new bounded Task in an Objective. |
| Finding open/triaged/mapped/deferred/owner_decision | Open Issue; preserve reasons, proof requirements, history and relations. |
| Finding in_progress/fixed | In-progress Issue; fixed still awaits independent verification. |
| Finding verified / resolved defect | Archive original proof/status; no new V2 Check implied. |
| Finding waived/duplicate | Archive disposition and canonical references; if active work needs it, preserve a typed historical decision mapping. |
| Audit register / runs / prompt | Archive exact history and custom instructions; unresolved entries map to Issues, custom active check procedures remain referenced. |
| Concept / Health-Check / visual identity | Preserve originals; carry useful active scope/procedures forward without mandatory new defaults. |

Do not guess whether two narrative findings are the same problem. Import explicit stable identities deterministically, preserve explicit duplicate links, and flag ambiguous unstructured audit findings for planner reconciliation. Preview enumerates ambiguous lifecycle values, missing targets, duplicate source IDs, and unresolved narrative findings before any write. An ambiguous active mapping blocks application until a concrete source-to-target decision is supplied; it is never silently discarded.

Migration protocol:

1. Read inventory of project-owned data and managed guide/skills, record raw bytes/hashes, permissions needed for restoration, and source-qualified references. Preflight version, writable destinations, path confinement, symlink/case collisions, and space/errors without modifying sources.
2. Preview deterministic target records, archives, ID mappings, conflicts, and required owner decisions. Dry-run creates no staging directories, probes, backups, or manifests.
3. On explicit apply, revalidate source hashes. Create a unique recoverable backup outside the live source tree, including impacted guide/skills; verify it before replacements. Never overwrite a prior backup or user sidecar.
4. Stage converted files and archive under a named migration operation; validate the full reference graph and preserve unknown content. Record per-path planned/installed hashes and recovery progress. Schema version activation happens last.
5. Publish through a tested platform-specific replace protocol. Do not reuse truncating-copy fallback for protected replacements. Multi-file changes are recoverable, not claimed atomic. Board/doctor/upgrade refuse normal writes while an operation is incomplete and provide recovery guidance.
6. Rerun recognizes the operation and either resumes validated identical work or reports a conflict if the user changed source/installed files. Recovery never overwrites those edits automatically. Successful second run changes no files, mtimes, or IDs. Validate board/doctor/resume against the result before recommending live use.

The exact Windows rename/recovery primitive and interruption schedule are an E45 design experiment before implementation. The required outcome is fixed: recoverable bytes and truthful failure reports, never silent loss. This uncertainty does not block E41's read-only fixture work.

## 10. Component boundaries and transition

| Component | Responsibility / proposed changes |
|---|---|
| internal/data | Canonical V2 records, schema/version detection, raw preservation, project index, dependency/evidence evaluation, named diagnostics, state writers. Temporary V1 input remains isolated and is retained only for migration after cutover. |
| internal/migrate (new) | One-time conversion planning, reference mapping, backups, operation recovery, schema activation. No board rendering or policy rules. Adds a justified Codebase Map row when implemented. |
| internal/init | V2 scaffolding and provenance-aware asset refresh; managed guide safety. Migration calls share low-level ownership helpers only where their contracts fit. |
| internal/doctor | Structural diagnosis and configured technical checks over the shared V2 data interpretation. Does not create Issues automatically or repair files. |
| internal/board | TUI and plain presentation, explicit asynchronous I/O, actions using canonical gate decisions. |
| internal/resume (new) | Read-only narrative presentation over the same project state/Next projection as board; command adapter remains thin. |
| cmd, main.go | migrate/resume dispatch and version-specific action wiring; no domain parsing in commands. |
| templates, agent-skills | Matching canonical/shipped skills, user-owned documents, tests for routing and obsolete terminology. |

No new database, service, or mandatory dependency is proposed. Any filesystem helper extraction must have actual consumers, not anticipate a generic framework. Keep current V1 guidance usable during implementation. Add V2 public skills under new names; retain old canonical/shipped pairs until the explicit cutover. Before making V2 the scaffold default, data, migration, and consumers must support it. Transitional readers/adapters have a removal owner (E50), not an indefinite compatibility promise.

## 11. Board, Next, resume, and health

Keep the three-column board: Planned, In Progress, Done. Show implementation/check/owner-wait/replan distinctions as badges and the prominent Next area, not new columns. Done-by-exception is visibly distinct; stale completion needs attention. Use Objective sidebar/selector, Task Outcome and User Check details, Check history, and one Issues overlay filtered by type. Preserve Defect labeling, related links, narrow widths, scrolling, stable focus geometry, monochrome and non-TTY output.

Shared Next projection gives deterministic precedence: incomplete migration/invalid target → replan → unsatisfied dependency/owner prerequisite → execute/verify → fresh Check → required owner validation → Objective integration Check → next ready Task/Objective → plan next Objective. Explicit router selection is retained; never silently select a similarly numbered Task. Archived/missing selection yields a named diagnostic and available next action. Fresh init with no Tasks yields Idea/Design planning rather than an error.

Router keeps its `## Current state` YAML anchor with `state: idea|design|task|check`, `objective`, optional `task`, and a short human next_action. These are routing hints; Task/Check/Issue records decide readiness. `REPLAN REQUIRED` derives a planner action without a fifth public phase. Resume and board use the same interpretation; resume performs no writes or verification subprocesses. It shows recorded evidence dates, freshness basis, owner waits, relevant Issues, and what changed since the last recorded handoff only when that information exists.

Health reports actionable structural/gate evidence rather than declaring every project with open Issues unhealthy. Advisory backlog does not block progress. Doctor failure output distinguishes malformed data, missing evidence, and failing configured checks from semantic review still needed. Keep archive content out of live diagnostics except reference-integrity checks.

## 12. Verification and delivery readiness

Software coverage: frozen V1 source fixtures; strict/raw parsing and byte preservation; global IDs and scoped legacy mappings; lifecycle authority and exception gates; dependency cycles/approval requirements; stale/unknown evidence; fresh vs repeated Checks; Issue proof/duplicate history; Objective integration; crash/retry/conflict/backup migration paths; init/upgrade ownership; shared Next; board/narrow/non-TTY; resume no-write behavior; six target packages/checksums.

Run required `make build && make test` at implementation handoffs, focused package tests for changed behavior, and the existing CI/distribution checks at release scope. Tests use temporary projects. Agent evaluations are separate evidence:

| Scenario | Pass evidence |
|---|---|
| Execute a fully planned task | Outcome met, explicit reads/extra reads recorded, no hidden architecture decisions, technical evidence and Check handoff. |
| Encounter a materially invalid plan | REPLAN REQUIRED, preserved partial work, clear planner decision, revised plan before resuming. |
| Fresh checker reviews seeded defect | Detect material seeded failure with reproduction; avoid advisory false blockers; record Issue and verify repair through a subsequent Check. |

Record session/model where supplied, task scope, context usage, replanning count, findings, and limitations. Do not claim all-model reliability from three scenarios. Use one tiny and one existing-codebase trial, including migration of copies. No live maintainer project mutation until recovery and the full core loop pass. Package release and publishing remain separate actions from planning and build validation.

## 13. Complete Task example

This is a proposed V2 execution contract example, not an active V1 task. Paths marked new are introduced by the preceding schema/Next work. Final task planning must confirm those APIs against the then-current code.

```yaml
id: T014
title: Resume unfinished work without changing project files
objective: O008
status: planned
depends_on: [{task: T013, requires: clear}]
owner_validation: {required: true}
planned_by: {role: planner, session: planning-example}
```

**Outcome:** reopening an existing V2 project shows selected work, recorded Check freshness, and one understandable next action without changing project files.

**User Check:** open a project with a Task awaiting Check; invoke resume as the human; confirm Task/outcome, owner-wait distinction, evidence date, and next action match the board. No automatic command execution or newly written evidence should appear.

**Done When:** correct Task/Objective and next action; stale/unknown evidence is explicit; malformed/missing selection is named; file bytes and mtimes unchanged; required owner validation recorded after technical clearance.

**Context Files:** `cmd/board.go`, `cmd/board_test.go` (command pattern); `internal/data/project.go` and `internal/data/next.go` (new in preceding work, shared index/projection); `main.go`, `main_test.go`; `internal/resume/resume.go`, `internal/resume/resume_test.go`, `cmd/resume.go`, `cmd/resume_test.go` (new in this example).

**Design References:** sections 4, 5, 10, and 11. **Guardrails:** FS-01, FS-03, DATA-02, DATA-03, ARCH-01, ARCH-03, TEST-01..04, TEST-08; apply V2 policy wording only after cutover.

**Implementation Plan:**

1. Confirm the preceding project/Next APIs exist; return REPLAN REQUIRED if they do not. Reuse their gate decisions rather than parsing Task frontmatter in the command.
2. Add thin resume argument handling using existing command runner injection, supporting help and optional project directory with named errors.
3. Add a renderer accepting the resolved immutable project projection and output writer. Show selected Objective/Task, outcome, implementation vs technical vs owner state, relevant Issues, evidence basis/date, and Next.
4. Keep presentation deterministic and usable without TTY. Reuse width-safe formatting for interactive output where needed; do not start a board process or subprocess to generate evidence.
5. Wire main dispatch through the injected runner pattern. Propagate read/output errors with path context and nonzero status.
6. Verify awaiting Check, awaiting owner, stale/unknown evidence, missing target, malformed router, no Task yet, and writer failure against the shared projection. Snapshot temporary project bytes/mtimes before and after invocation.

**Boundaries:** no new Task states, evidence collection, Objective creation, automatic model routing, or duplicate gate logic.

**Technical Verification:** focused cmd/resume/data tests, no filesystem changes on success/failure, parity with board Next fixtures, `make build && make test`.

**Technical Evidence:** pending execution; record named cases, results, reviewed source basis, files read/changed, and limitations. **Drift Notes:** new resume package/map update if not already introduced; reconcile through executor/planner before Check. This example intentionally retains implementation detail and separate owner validation.

## 14. Delete, merge, retain, and review boundary

Rename PRD → Idea and Epic → Objective in V2 public assets. Remove mandatory releases and release PRDs from new projects; retain release metadata. Merge finding/defect infrastructure into Issues and audit/run evaluations into Checks. Replace register bookkeeping with derived listings; preserve historical register/runs/waivers in archive. Retain Guardrails, doctor, router anchor, detailed Task bodies, upgrade provenance, shared check rigor, visual identity where relevant, and byte-owned user documents.

Each retained abstraction has one purpose: Idea defines intent; Design explains technical reality; Guardrails protect constraints; Objective bounds integration; Task bounds execution; Check records independent evidence; Issue holds durable follow-up; router selects current work; config runs project commands; archive preserves source history; migration manifest makes conversion recoverable; upgrade manifest protects customized assets. No separate owner-validation Issue, duplicate completion flag, second health policy, or parallel authoritative backlog is needed by default.

Readiness: E41 can proceed after this design/backlog review. No unresolved architectural choice is embedded in its fixture Tasks. Later epics need scoped design refinement before task creation, particularly platform recovery and exact writer APIs. This proposal is not an independent Check of its own implementation.
