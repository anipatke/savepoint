# Savepoint guidance consolidation: before/after comparison

This is a maintenance comparison, not a Savepoint Check record or technical `CLEAR`. The baseline is the working tree immediately before consolidation, **including the completed Idea and Design interview additions**. Git HEAD was `328d513186c3a224f504ec53b7158acd650ca466`; using HEAD alone would incorrectly count those additions as part of this cleanup. The pre-existing README edits and generated temporary test files were excluded and left untouched.

The baseline copies are at `/tmp/savepoint-guidance-before`; a complete diff against those copies is at `/tmp/savepoint-guidance-before-after.patch`. Those temporary paths are session-local; this report records the durable comparison, decisions, and validation. Counts use whitespace-delimited words, not model tokens. The report itself is not a shipped template or a skill dependency and is excluded from savings.

## Size

| Document | Before words | After words | Reduction |
|---|---:|---:|---:|
| `AGENTS.md` | 3,151 | 2,793 | 11.4% |
| `.savepoint/router.md` | 412 | 121 | 70.6% |
| `agent-skills/savepoint-idea/SKILL.md` | 865 | 787 | 9.0% |
| `agent-skills/savepoint-design/SKILL.md` | 3,229 | 2,605 | 19.3% |
| `agent-skills/savepoint-task/SKILL.md` | 2,061 | 1,381 | 33.0% |
| `agent-skills/savepoint-check/SKILL.md` | 2,177 | 1,569 | 27.9% |
| `agent-skills/bubbletea-tui-design/SKILL.md` | 870 | 789 | 9.3% |
| `agent-skills/references/check-method.md` | 3,127 | 2,864 | 8.4% |
| `agent-skills/references/issue-capture.md` | 1,167 | 922 | 21.0% |
| `agent-skills/references/commands-and-procedures.md` | 506 | 292 | 42.3% |
| `templates/project-v2/AGENTS.md` | 2,727 | 2,190 | 19.7% |
| `templates/project-v2/.savepoint/router.md` | 407 | 119 | 70.8% |
| `templates/prompts/magic-prompt.prompt.md` | 86 | 58 | 32.6% |

The complete shipped `templates/` tree decreased from **17,457 to 13,892 words**, saving **3,565 words (20.4%)**. This includes all seven packaged skill/reference copies once each, the starter prompt, AGENTS, and hidden scaffold files. The live skills have the same reductions as their corresponding shipped copies. Bubble Tea is repository-only. Repo AGENTS retains its detailed Codebase Map and concrete build rules; the scaffold has neither, so its reduction differs.

The scaffold Idea (87 words), Design (265), Guardrails (644), config (52), initial Goal (57), and objectives `.gitkeep` are byte-identical to baseline. No rule table, severity definition, config default, Goal placeholder, or scaffold field was deleted.

## Ownership after consolidation

| Information | Canonical home | How other entrypoints reach it |
|---|---|---|
| Routing, selection, terminology, owner authority, verification policy, CLI permissions, worktree constraints | Project AGENTS | Router and skills point to its named sections |
| Intent intake and focused product questions | Idea skill | Router/AGENTS activate Idea |
| Current-Objective planning, readiness, decision interview, artifact templates, Goal creation and optional retrospective | Design skill | Router/AGENTS activate Design |
| Scoped execution, stage progress, extra reads, replan and waiver evidence | Task skill | Router/AGENTS activate Task |
| Independent verdict and exact Check/acceptance/exception record shapes | Check skill | Router/AGENTS activate Check |
| Quick/Full procedures, frozen scope, matrices, health collection, materiality, recheck convergence | Check method | Check loads it in full and applies the selected mode |
| Issue fields, deduplication, history, dispositions and direct-repair routing | Issue capture | Design, Task and Check enter it from their own workflow |
| Legacy command/procedure reconciliation | Commands reference | Design references it for reconciliation |
| TUI mechanics | Bubble Tea skill | Current Task scope, project identity and gates still govern |

No new skill, reference, router state, record type, runtime field, command, or approval checkpoint was introduced. This report is a requested review artifact, not another workflow input.

## Detailed requirement ledger

“Retained” means the obligation remains explicit at the named destination. “Moved” means repeated prose was removed while its authoritative requirement and route to it remain. The corrections section below separately records changes to contradictory or stale guidance.

### Routing and project policy

| Before requirement or material | After location and disposition |
|---|---|
| Pasted Next is the selection; otherwise run read-only resume | AGENTS / Workflow, unchanged |
| Dispatch Start/Build/Test, Check, Plan/Replan, Pick a Task, Fix, owner-decision terminal words | AGENTS / Workflow, unchanged |
| Issue is context when Objective/Task is selected; do not substitute stale work | AGENTS / Workflow, unchanged |
| No named record means use router state, not a guessed record; missing tool fallback | AGENTS / Workflow and repo CLI Rules, unchanged |
| Four skill mappings, unavailable-skill file fallback, REPLAN is not a fifth state | AGENTS / Skill Activation, unchanged; router points there |
| Validate every router target and Task ownership; clear selection keys exactly; preserve release | AGENTS / Router Selection, unchanged |
| Direct Issue repair chooses only a uniquely linked Objective, otherwise clears issue alone; show resume | Issue capture / Out-Of-Scope Repair, retained; AGENTS and Task explicitly point there |
| Shared references never trigger independently | AGENTS reference inventory and every reference frontmatter, retained |
| Idea is original intent; Design is readiness/implemented reality | AGENTS / Router Selection and Context Budget, retained |
| Per-criterion Task evidence; configured gates; focused tests never satisfy handoff | AGENTS / Verification Policy and Task / Evidence And Handoff, retained |
| Ordinary, migration/platform-sensitive, Full Check and CI gate distinctions | Repo AGENTS / Verification Policy and Build unchanged in substance; scaffold generic gate ownership retained |
| Full-result reuse only for metadata; command/time/toolchain/result and unchanged code/tests/fixtures/dependencies/gate definitions | AGENTS / Verification Policy, literal retention; skill gate sections explicitly delegate |
| Optional Task Check needs explicit Task/reason/actor/time waiver; waiver is not CLEAR | AGENTS / Verification Policy and Task waiver block, retained |
| Waiver satisfies clear dependencies, never accepted; runtime determines blocks | AGENTS / Verification Policy, retained; Task and Check retain supporting explanation |
| Full Objective Check covers every Task including waived ones, integration and Design reconciliation | AGENTS policy, Check scope/closure and method depth, retained |
| Independent session, immutable runs, supersedes, signed CLEAR current without another assessment | Moved from AGENTS / Check into Verification Policy; Check record contract retained |
| Required live Goal, Objective release field, init G-001, migration live-selection/continuation and unresolved Apply block | AGENTS / Required Goal Context, byte-identical |
| Missing/unknown/archived Goal diagnostics, missing Objective release, no project-wide board fallback | Same Goal-context section, byte-identical |
| Stable R identities, new G identities, Release-compatible paths, derived completion, no Goal Tasks or publishing | Same Goal-context section, byte-identical |
| Task statuses and stages; audit is ready, not passed; no implementation stage | AGENTS / Terminology, retained |
| Executor may start Task and planned Objective, never complete/retreat status | AGENTS / Terminology, retained; owner stop instruction moved here |
| Checker-only verified Issue closure; owner acceptance/reopening; planner escalation | AGENTS / Terminology and Issue capture / Role Boundaries, retained |
| Capture durable out-of-scope defects/drift/guardrail gaps rather than expanding Design/Tasks | AGENTS / Issue Capture and three skill entry sections, retained |
| Lane detection, no router/identity writes, note follow-up, commit but no push/merge, Checks on main after merge | AGENTS / Worktree Lanes, byte-identical |
| Guardrails owns STYLE; advisory severity and absent policy do not block | AGENTS / Code Style, unchanged; role-specific reminders retained |
| Existing-code adoption uses targeted reads, no scan/automatic pass/model call, unknown unread areas | Template AGENTS / Existing Codebase Adoption, compressed but explicit |
| Adoption separates actual structure from owner intent; only managed-guide block exception preserves all outside bytes | Same adoption section, explicit; existing preservation tests still pass |
| Optional Concept/Health-Check/procedure files; required Goal and migration guidance | Same adoption section; identity/diagnostic details reference Required Goal Context |
| Only resume plus narrow create-task/Full-Check-health exceptions; setup/report human-only | AGENTS / CLI Rules, retained; Check explicitly names setup/report prohibition |
| Repo source-built CLI fallback and project-specific Codebase Map | Repo AGENTS / CLI Rules and Codebase Map, unchanged |
| Technical Context Log versus plain owner-facing reporting | Both AGENTS / Reporting to the Owner, unchanged |

### Idea and Design

| Before requirement or material | After location and disposition |
|---|---|
| Rough sentence is enough; no prepared requirements prerequisite | Idea / Workflow and starter prompt, retained |
| Idea reads only router, Idea, selected live Goal, stated intent, targeted existing evidence | Idea / Read, unchanged |
| Initial G-001 already exists; fill owner sections without another Goal, preserve existing Goal content | Idea / Goal paragraph and Workflow, retained |
| Idea captures intent, user, experience, scope, exclusions and success without architecture/components/interfaces | Idea / Purpose, Interview and Rules, retained |
| Ask owner about material product uncertainty rather than infer it | Moved from duplicate workflow/rule lines into Idea Interview, explicit |
| One question at a time, wait, concrete options/tool fallback, acknowledge, follow dependencies, avoid reopening settled choices | Idea and Design interview sections retained |
| Idea stops at sufficient owner detail and summarises before writing | Idea Interview retained; fresh-Goal clause clarified as applicable only when a fresh placeholder exists |
| Idea artifact fields/sections, write boundary, state design handoff without Objective detailing | Idea artifact and Rules unchanged; handoff retained |
| Design is implemented reality, Objective deltas planned change, Guardrails durable policy | Design / Purpose, Workflow, Design Template and Guardrails Rule, retained |
| Read only current planning context and targeted readiness evidence; no detailed distant backlog | Design / Read, Workflow and Rules, retained |
| One active Objective; inspect requirements, owner confirms design before readiness/Task detailing | Design / Workflow, retained; duplicate product-choice step removed |
| Confirmation is not execution approval or Task completion; owner reviews Task plan before routing | Design / Workflow steps 5 and 10, retained |
| Research Task for unknown approach, split unrelated outcomes/unresolved architecture, genuinely independent lanes | Design / Workflow, retained |
| Replan reassesses affected decisions; no speculative replacement or new approval checkpoint | Design / Workflow and interview, retained |
| Readiness: settled interfaces/ownership, scoped constraints, known dependency outcomes, verification approach | Design / Readiness Gate, byte-identical |
| Platform evidence for processes/paths/signals/replacement and named producer | Design / Workflow step 9, retained |
| ID-free draft, no predicted IDs/filenames, locked allocation/validation, retired reservation not reused | Design / Task Creation, unchanged |
| Strict resume after other identity creation/rename | Task Creation and AGENTS CLI Rules, retained |
| Goal diagnostics/history repeated in Design | Compressed to summary plus explicit AGENTS / Required Goal Context reference; planner creation rules remain local |
| New Goal first-unused global G identity, stable identity, duplicate refusal, four sections and release membership | Design / Required Goal Context creation steps, retained |
| No second membership list or Goal nesting/Task ownership; derived completion | Same creation steps and Objective-template membership text, retained |
| Optional, owner-requested Goal retrospective, explicit no-change outcome, packaged-skill suggestions as Issues, no new state/list, never holds a Goal open | Design / Optional Workflow Retrospective (made optional after this comparison) |
| Objective artifact metadata, allowed statuses, dependency shape, freshness actor/date/basis and body headings | Entire Objective artifact fence, byte-identical |
| Task artifact required title/objective/status/dependency/owner_validation/planned_by and every body heading | Task frontmatter fence byte-identical; complete worked-example body retained in shorter prose |
| Exact Context File paths, no globs/directory entries; design/guardrail references and scoped implementation plan | Task artifact, retained |
| Example covers selection/Next, owner-wait, evidence date, no writes/collection, malformed selection, failure and stale/unknown cases | Task example Outcome/User Check/Done When/Implementation Plan, retained in compact form |
| Example owner validation follows mandatory integration evidence | Task / Done When, restored explicitly during comparison |
| Evidence commands/outcomes/files/limitations and planner reconciliation of architecture deltas | Task example evidence/drift sections plus Task evidence contract, retained |
| Readable owner-facing title distinct from Outcome; exactly one O-reference; planning metadata completed | Design / post-template field contract, retained |
| E50 test-history commentary | Removed historical explanation; retained scenario-based readability requirement and changed its existing assertion accordingly |
| Design section inventory and durable 10–20 Guardrails with stable IDs/severity/exception authority | Design Template and Guardrails Rule, unchanged |
| Migration reconciliation and out-of-scope Issue entry/authority | Design reference-entry sections, unchanged |

### Task and Check

| Before requirement or material | After location and disposition |
|---|---|
| One scoped Task; no acceptance rewrite, silent scope expansion, self-clearance, Check or invented waiver | Task / Purpose, Workflow and consolidated Rules, retained |
| Blocked start is reported, not bypassed or substituted | Task / Workflow step 1, unchanged |
| Free reads versus necessary extra reads; log before reading/editing, files and reason | Task / Read and Workflow step 5, retained; duplicate Extra Reads section removed |
| Router selection preserves release; Task/Objective planned-to-in-progress start | Task / Workflow step 2, unchanged |
| Build/test/audit progress with acceptance evidence and gate outcomes | Task / Workflow and Evidence And Handoff, retained |
| Replan frontmatter reason/actor/time, body explanation, unchanged status/stage, preserved partial work and stop | Task / REPLAN REQUIRED including exact YAML block, retained |
| Body-only REPLAN does not route; material missing/contradictory/impossible examples | Same section, retained |
| Task NEEDS WORK resumes build; Objective NEEDS WORK never retreats done Tasks | Task / Lifecycle and Check / Workflow, retained |
| Direct repair proof and router handoff versus new planning | Task / Lifecycle explicitly enters Issue reference, whose exact routing/proof steps remain |
| No configured gates means record none, never invent; focused runs iteration-only; reuse limits | Task / Verification Gates plus shared AGENTS policy, retained |
| Health advisory, watch line not aim, risky production preferred, flat dispatch untouched, no branch-shuffling metric gaming | Task / Acting On A Code Health Report, byte-identical |
| Record remaining above-watch results; scope narrowing owner-only; cite Guardrail ID | Same health-report section, byte-identical |
| Every acceptance criterion, commands, read/changed files and limitations at handoff | Task / Evidence And Handoff, retained |
| Explicit owner waiver in exact frontmatter shape; prose alone insufficient | Same section including YAML block, byte-identical shape |
| Handoff to fresh Check only when requested; waived evidence consumed by mandatory Objective Check | Task / Workflow, Lifecycle and Evidence And Handoff, retained |
| Check reads scoped work/policy/prior evidence/changed code and complete method | Check / Read, retained |
| Method Quick for requested Task, Full for Objective, cross-Task integration and Design reconciliation | Check / Workflow and method mode/depth, retained |
| New immutable Check, strict index load, supersedes prior record; no remediation edits | Check / Workflow, artifact and consolidated Rules, retained |
| NEEDS WORK records all Issues, routes executor/planner; CLEAR alone never closes | Check / Workflow, retained |
| Check frontmatter and optional health_snapshot; reviewed metadata may be omitted, paths/hashes/empty lists not fake proof | Check / Artifact including exact YAML fence and explanations, retained |
| Check body outcome coverage, commands/probes, policy/style, needed owner validation and observations | Same artifact contract, retained |
| Technical Task clearance and no unexcepted blocker; owner-only done | Check / Closure Rules, byte-identical |
| Required owner validation names same current Check; superseded acceptance invalid; exact accepted_by shape | Same Closure Rules, byte-identical |
| Explicit Check-specific exception shape, not CLEAR/current clearance | Same Closure Rules, byte-identical |
| Objective all Tasks done, current integration Check, resolved material linked Issues and conditional owner acceptance | Same Closure Rules, byte-identical |
| No Objective exception excuses unfinished Task; no Objective Check closes a Task | Same Closure Rules, byte-identical |
| Goal derived completion, no Goal-level clearance/acceptance, not publishing | Same Closure Rules and shared Goal policy, retained |
| Missing scope/evidence cannot complete; stale/unknown blocks until new Check | Same Closure Rules, byte-identical |
| Only Full Check collects official health after full gate; manual snapshot not evidence | Check / Code Health Evidence and method / Collect Code Health Evidence, retained |
| Health test/coverage reuse gate reports, absent config nonfinding, blocking verdict prevents CLEAR | Method / Collect Code Health Evidence, retained; blocker examples moved here from skill |
| Optional failures/warnings never automatically create Issues; aim-to-watch nonblocking observation | Method health/Issue judgment and Check / Code Health Evidence, retained |
| Checker may verify Issue; owner acceptance/reopening is not CLEAR | Check / Issue Capture and Issue reference, retained |

### Check method and shared references

| Before requirement or material | After location and disposition |
|---|---|
| Load entire method, apply Quick procedure or every Full section, optional policy/procedure absence skipped | Method introduction/modes, retained; repeated introductory paragraphs removed |
| Quick scope lock/invariants/boundary/failure/bypass/file/gate/Issue/materiality/style steps; no matrix/adversarial widening | Quick Check Procedure, byte-identical |
| Current diff/evidence are claims, file reality, scoped acceptance/policy/Design evidence | Establish Scope, byte-identical |
| Numbered frozen lock, supported public/runtime/effect boundary, initial factual amendments only, immutable recheck | Freeze The Check Scope, byte-identical |
| Every criterion becomes invariant; affected inputs/state/output/environment/APIs, independent scenario and expected/actual proof | Turn Acceptance Into Invariants, equivalent wording; normal/boundary/malformed/failure/bypass obligations retained |
| Every coverage axis/cell, justified N/A, deterministic harness, no ad-hoc substitution, initial missing-cell correction | Build The Mandatory Coverage Matrix including finite external boundary subsection, byte-identical |
| Every real operation/effect/failure/cleanup/final state, independent oracle, semantic state and failure ordering | Workflow And Side-Effect Check Lock including Matrix Completion Lock, byte-identical |
| Every mandatory cell classified; complete remaining cells after finding; unverified/NEEDS WORK when incomplete | Same completion lock, byte-identical |
| Applicable validation bypass, backward/skip/revival, mode/ownership/safety, output secrecy, boundaries and input-class adversarial questions | Adversarial Pass converted into explicit list; matrix boundary axes retain exact below/above limits; no class deleted |
| Frozen-cell admission ledger, original cells/reproductions/changed paths/named adjacency; credible-blocker exception | Re-check After Remediation, byte-identical |
| Original Issue closure map and initial/full/targeted convergence limit, no third autonomous round | Same recheck section, byte-identical |
| Named files exist/deleted/discarded; unexplained phantom Issue | Verify File Reality, byte-identical |
| Focused and default/excluded-file lint/type gates, diff check, configured gates, evidence mode | Verify Evidence And Gates, byte-identical |
| Health official ID, config absent handling, command exit semantics, verdict reading, failed tool is not bad code | Collect Code Health Evidence, retained; explicit blocker examples added from removed duplicate skill text |
| Proven/Issue/Unverified classification, every criterion proven for CLEAR | Complete The Issues Pass, byte-identical |
| Supported-path, introduced/touched/promised, scope, reproducibility and credible consequence admission; all Issues together | Same Issues pass, byte-identical |
| Smallest repro, expected/actual, exact source evidence and inadequate tests | Same Issues pass, byte-identical |
| Likelihood/impact/materiality independently assessed, severity not priority, no implicit waiver, observations separate | Summarize Materiality, byte-identical |
| STYLE checklist follows project order/labels, unticked reason/evidence advisory, skip without rules | Review Code Style and example, byte-identical |
| Issue schema, optional fields, type alone never blocks | Issue artifact including YAML fence, byte-identical |
| Search before identity allocation, no automatic matching; strict resume after identities | Issue / Search Before Creating, retained |
| Exactly one of verified/accepted/duplicate/escalated; proof and authority distinctions; recurring Issue reuses ID | Issue / Dispositions and Role Boundaries, retained |
| Owner Space/Backspace, In Progress restriction, fixed reason/actor/time/history, explicit agent instruction only | Issue / Role Boundaries, retained after moving duplicate accepted text |
| Escalated retirement immediate, no prior Check needed, planner actor/Objective link/history; new Objective owns future proof | Issue / Escalation Retires The Issue, retained; duplicate prose removed |
| Append-only history, refusal of shortening/reordering/editing; deferral does not add a status | Issue / History Is Append-Only, retained; escalated kind made explicit |
| Before repair_attempted rerun Proof Needed/nearby cases, record each result, unavailable cases unverified | Issue / Out-Of-Scope Repair, unchanged |
| Exact unique-Objective router handoff, preserve release, show resume, leave repair open pending owner/checker | Same direct-repair procedure, unchanged |
| New Objective only for planning need; new bounded Task is not escalation and Issue remains open | Same repair procedure, retained |
| quality_gates exact keys, no checks.technical, build root/order/timeout/failure semantics | Commands / Config Contract, retained in shorter prose |
| Legacy Health-Check splits exactly three ways: configured commands, existing method, optional manual procedure | Commands / Three-Part Mapping, retained |
| Migration archives bytes and previews fenced command candidates, never writes config or procedure | Commands / What migrate actually does, retained |
| Procedure ordinary and Task-referenced, never default/required; absent file no substitute/finding | Commands / Preserved Procedure and Absence, retained |
| One-off verification/evidence belongs to Task, not shared policy/config/method | Commands / Extra Verification Belongs To The Task, retained |

### TUI and remaining templates

| Before requirement or material | After location and disposition |
|---|---|
| Current Go/Bubble Tea/Lip Gloss stack, visual identity and smallest Task surface, no React/Ink/TS TUI replacement | Bubble Tea / Current Stack, Quick Workflow and package map, retained |
| Model explicit, no globals; Update deterministic typed reducer; Cmd owns I/O/timers/reload/write | Bubble Tea / Rules, combined without removal |
| Modal-first key dispatch, repeat safety, explicit router actions, quit precedence | Same TUI rules, unchanged |
| Cell geometry, focus-stable layout, border family, truncation, terminal clamp, shared styles and fallback colours | Lip Gloss Layout Rules, retained |
| Small helpers, selection/render separation, colour reinforced by text, stable-width glyphs, brief copy, deterministic non-TTY | Rendering Rules, unchanged |
| Changed branch/render/edge tests, messages/commands, temp fixtures, stable labels/truncation/order, narrow/wide/non-TTY | Testing Rules, retained |
| Common traps | Duplicates removed; unique monochrome/glyph and focused-versus-golden advice moved to Layout and Testing |
| Correct current package tests then required project gates | Testing Rules corrected to board/v2 and AGENTS gate ownership |
| Router state/Objective/Task/release/Issue selection values | Both YAML fences byte-identical; prose now delegates instead of forcing every phase to load Idea/Design |
| Router four mappings/lifecycle/verification descriptions | AGENTS authoritative mappings and policy retained and explicitly referenced |
| V1 archived sources not active routing inputs | Both router summaries explicit |
| Starter AGENTS/router startup and rough Idea intake | Prompt now routes via pasted Next or read-only resume and the active Read budget |
| Idea/Design/Guardrails/config/Goal scaffold content and defaults | All unchanged byte-for-byte; zero rules or fields removed |

## Intentional corrections, distinct from compression

1. **Quick versus Full:** AGENTS previously implied matrices and adversarial probing in both modes. It now delegates Quick to the existing Quick procedure and Full to every method section. The Quick procedure and all Full matrix obligations are unchanged. This resolves conflicting prose; it does not reduce the existing method's depth.
2. **Router read order:** removed unconditional Idea/Design/Objective reads that contradicted phase-specific Read sections. The router explicitly follows the active skill's budget. Current router selections are unchanged.
3. **Same-session review:** removed the ambiguous “continue anyway” route to an independent Check. A self-review may provide observations when requested, but writes no Check and cannot satisfy independence. This strengthens the existing independent-session requirement.
4. **Repository read budget:** replaced the conflicting blanket ban on all extra reads with the active Task skill's existing narrow allowance: necessary, targeted, logged before reading/editing. This is a reconciliation with the canonical skill, not permission for general exploration.
5. **TUI paths and gates:** corrected V1-era board paths to existing `internal/board/v2` paths; project gates remain authoritative rather than treating bare `go test ./...` as handoff.
6. **Issue vocabulary:** changed “not a fourth router state” to “not a fifth router state”; added `escalated` to the history list because it already appeared in escalation instructions and is an existing runtime kind (`internal/data/issue_v2.go`). No new vocabulary was invented.
7. **Fresh Goal interview:** made the stopping clause conditional on a fresh Goal placeholder being applicable; existing Goal content is still preserved.

## Behavioural scenario walkthroughs

These are manual comparisons of decisions required by the before/after text, **not executed model conversations or independent technical Checks**. Each expected decision is grounded in the retained contract above.

| Scenario | Required after behaviour | Comparison |
|---|---|---|
| Owner pastes a Next while router names different work | Use pasted selection; do not resume to reconfirm or substitute | Same |
| Named Task belongs to another Objective | Reject router edit; preserve selection/release | Same |
| Task dependency reported blocked | Report block; do not start or pick alternate work | Same |
| Owner supplies only a rough sentence | Begin scoped Idea questions; no prepared-document prerequisite | Same |
| Idea session sees unresolved component choice | Leave architecture to Design; ask only product intent | Same |
| Design has unconfirmed verification requirement | Ask owner; no Task detailing before confirmation/readiness | Same |
| Approach is technically unknown | Bounded research Task with decision deliverable, not guessed implementation | Same |
| Independent lane has a needed Issue | Note it in evidence; no identity or router write; commit locally | Same |
| Missing Context File discovered during Task | Record exact replan block, keep status/stage and partial work, stop | Same |
| Necessary extra read discovered | Log file/reason before read; remain targeted | Canonical skill preserved; conflicting repo wording resolved |
| Owner waives Task Check | Record exact explicit waiver; satisfies clear, never accepted; Full Check still required | Same |
| Executor asked to check own work | Self-review observations only; fresh checker needed for gate | Independence strengthened |
| Quick probe exposes risk beyond Task | Nonblocking observation for Full Check; do not widen Quick scope | Same existing method |
| Full Check unfinished matrix due unavailable environment | Classify unverified and NEEDS WORK; never shrink coverage | Same |
| Remediation reveals unrelated ordinary edge case | Outside frozen cell: observation; no new blocking perimeter | Same |
| Recheck finds credible Guardrail blocker outside lock | Apply exact credible-blocker exception and name rule | Same |
| Third autonomous remediation cycle proposed | Stop for owner decision at convergence limit | Same |
| Optional health tool fails | Distinguish failed tool from bad code; do not auto-open Issue | Same |
| Health command exits 0 with blocking verdict | Read verdict and withhold CLEAR; exit code does not establish clearance | Same |
| Objective fails after Tasks done | Direct Issue repair or planned work; do not retreat done Tasks | Same |
| Direct repair links zero or multiple Objectives | Clear issue only, preserve release, resume; do not guess | Same |
| Repair promoted into new Objective | Immediate escalated retirement; new Objective carries future proof | Same |
| Acceptance refers to superseded Check | Does not satisfy current owner validation | Same |
| Optional procedure or Guardrails absent | Skip related optional step; absence not a finding | Same |
| Focus colour changes border geometry | Correct layout; colour/glyph/style must preserve dimensions | Same |

## Validation and limits

- Comparative review inspected removed/changed prose against remaining source or explicitly named canonical destinations; restored health-blocker examples and worked-example owner-validation timing when the first draft compressed them too far.
- Automated comparison confirmed unchanged frontmatter for all five skills and three references; unchanged literal YAML contract blocks; unchanged Objective/Task artifact frontmatter; unchanged router selection fences; unchanged scaffold documents/config; and byte parity for all seven shipped skill/reference files.
- Ten critical Check-method sections remain literal copies: scope establishment/freezing, coverage/external matrix, workflow/completion lock, remediation recheck/convergence, file reality, evidence/gates, Issue verdict/admission, materiality, and code style. Quick procedure and Check closure rules also remain unchanged.
- Existing routing-template test now checks all four mappings at their authoritative AGENTS destination and rejects the retired broad read-order/policy duplication. Existing title test preserves scenario-based readability beyond field validation, replacing an E50 historical-reference assertion.
- `go test ./internal/init` passed after those changes, including scaffolding and preservation tests. All five live skills passed skill-creator validation. Changed-file whitespace checks passed.
- `make build && make test-fast` passed on 2026-10-02, completed before 22:43:15 UTC (2026-10-03 Australia/Sydney), with `go version go1.26.2 linux/amd64`. The first sandboxed build hit a read-only Go cache; the successful gate ran with approved cache access. No gate definitions or production runtime code were changed. This was a fresh ordinary handoff gate, not reused full evidence or a Full Objective Check.
- Final `savepoint resume` strict-loaded the project successfully and returned the same completed R-007 Goal/Next as before; no router selection or lifecycle status changed.

No unique normative obligation was identified as omitted or weakened by this comparison. Explicit conflict resolutions are listed above rather than being presented as byte-equivalent changes. This is a manual semantic audit supported by structural comparisons and existing tests, not a formal proof of agent behaviour. No independent Savepoint Check or owner completion decision is claimed.
