---
type: check-method-reference
triggerable: false
---

# Shared Savepoint Check Method

This non-triggerable reference is loaded in full by `savepoint-check`, which owns the trigger, scope, mode, and record. Design and Task reference it but never run a Check themselves. Apply AGENTS.md's Verification Policy for waivers, dependencies, gates, and completion authority.

Use Guardrails and optional project verification procedures when present; skip the related steps when absent. Absence is not an Issue.

## Task Check And Objective Check Depth

A requested Task Check evaluates one Task's acceptance criteria, plan, evidence, and scoped files. The mandatory Objective Check additionally covers every owned Task, including waived Tasks, cross-Task integration, and Design reconciliation. Task clearance never substitutes for Objective clearance; an owner waiver never waives acceptance criteria or guardrails.

## Quick And Full Evidence Modes

- **Quick** — optional, only when the owner requests or selects a Task Check. Follow the Quick Check Procedure below; never add an automatic Check to every handoff.
- **Full** — mandatory for Objective closure; apply every section of this method.

### Quick Check Procedure

1. Establish Scope, below.
2. Write a short scope lock: the acceptance criteria and guardrails being
   tested, the changed files, and what is out of scope. It is frozen for any
   re-check the same way as a Full lock.
3. Turn Acceptance Into Invariants, below, for each acceptance criterion.
4. Probe the changed code paths directly: boundary values, malformed or
   missing input, failure behavior, and one bypass path where one exists.
5. Verify File Reality and Verify Evidence And Gates, below.
6. Complete The Issues Pass, Summarize Materiality, and Review Code Style,
   below.

Quick mode does not build the coverage matrix, the external-boundary matrix,
the workflow and side-effect lock, or the adversarial pass. If a probe shows
risk that reaches past the one Task, record it as an observation for the
mandatory Objective Check rather than widening the Quick Check.

## Establish Scope

1. Inspect current files and the diff. Treat context logs, checked boxes,
   prior Check records, remediation claims, and green tests as claims to
   verify, not proof.
2. Verify file reality for every scoped file before referring to it.
3. Load only the acceptance criteria, guardrail rules, Design context, and
   evidence-mode requirements that govern the scoped work.

## Freeze The Check Scope

Before the first adversarial probe, trace the real code and workflow, then
write a numbered scope lock containing:

1. the acceptance criteria, guardrails, and release gates being tested;
2. the changed files and supported public entry points;
3. the directly relied-on runtime orchestration, external effects, and
   dependency behavior needed for those entry points to keep their promise;
4. the selected matrix axes and cells, including explicit not-applicable
   cells;
5. the supported-path and materiality boundary used to admit an Issue.

The scope lock defines what the Check means by `complete`. Do not expand into
unrelated dependency internals, unsupported configurations, or increasingly
remote hypothetical states.

During the initial Check only, new source evidence may correct a factual
mistake in the scope lock. Record the amendment and restart the affected
matrix pass. Once the initial verdict is returned, the lock is immutable for
every re-check. Do not silently introduce a new axis, dependency layer,
acceptance interpretation, or meaning of "adjacent" during remediation review.

## Turn Acceptance Into Invariants

For every acceptance criterion:

1. Restate it internally as a general rule, not as one example.
2. List the inputs, state transitions, output paths, environment modes, and
   public entry points that can affect the rule.
3. Check normal, boundary, malformed-input, and failure cases, plus at least one bypass path. Full mode uses the applicable matrix axes below.
4. Run at least one independent scenario that is not merely an existing unit
   test repeated unchanged.
5. Record the expected result, actual result, and concrete evidence.

A regression test proves its example. It does not prove the surrounding
invariant.

## Build The Mandatory Coverage Matrix

Full mode only. A Quick Task Check uses the Quick Check Procedure instead.

Before running focused probes, create a concrete matrix for the scoped public
behavior. List every row and applicable axis; mark a cell not-applicable only
with a reason tied to scope or an acceptance rule. A prose checklist is not a
matrix.

Always include these axes when applicable: public surfaces (every
constructor, factory, parser/validator, state reducer, pure renderer, output
backend, and selection/helper entry point); input shape (normal, empty,
missing, duplicate, wrong type, mixed type, non-finite, mutable input, and
mutation after validation); state (every state, every allowed transition, and
backward, skipped, overlapping, terminal-revival, post-failure, and
representation-switch transitions); environment/output (interactive and
redirected sinks crossed with normal, no-colour, dumb/no-cursor, failure,
warning, timeout, and explicit public overrides); boundaries (the exact
limit, immediately below and above it, and the full small finite range when
practical); sequences (complete workflows including initialization,
intermediate milestones, completion, failure, retry/repeat, and final
output); representations (serialization round trips, direct models,
structured events, and rendered output where each exists); and text classes
(ASCII, control characters, combining marks, variation selectors, wide
characters, emoji modifiers, regional flags, and joined emoji when text width
or truncation is in scope, using an independent oracle when one is available).

Run the matrix through repeatable parameterized tests or one deterministic
Check harness where possible. Record its rows, cell classifications, and
command/output evidence. Ad-hoc probes may supplement the matrix but cannot
replace it.

Once the initial Check starts its adversarial probes, this matrix is the
frozen scope lock. A missing required cell discovered during that initial
Check is a Check-process error: amend the lock explicitly and rerun the
affected matrix before returning a verdict. A re-check may never add the
missing cell retroactively as a new blocking perimeter.

### Finite External-Boundary Matrix

When scoped code relies on a server, subprocess, browser runner, provider, or
other external boundary, classify this finite set during the initial Check:
configured target and actual runtime target; startup, discovery, and reuse
ordering; connection refusal or unavailable dependency; successful and
non-success responses; redirect behavior; timeout and cancellation;
malformed or unexpected response; retry, cleanup, and partial side effects
where applicable; and secret-safe failure output.

Mark a cell not-applicable only with a scope reason. Do not invent additional
network or toolchain edge classes during later re-checks unless they meet the
credible-blocker exception below.

## Workflow And Side-Effect Check Lock

Full mode only. A Quick Task Check uses the Quick Check Procedure instead.

Apply this lock to any command or workflow with multiple operations,
external calls, persistence, transactions, generated artifacts, cleanup, or
structured progress. Derive the inventory from the actual code path and
external effects, not from task notes or a test table that is meant to
verify it.

Before testing, record a table for the complete workflow: order, real
operation, side effect / state change, failure timing, failure owner and
final state, cleanup / secondary failure, and independent oracle.

The inventory must include, when applicable: input parsing, guards, setup,
constructors, and connection/client creation; external calls and each
operation that can partly succeed; transaction begin, write, commit,
rollback, and retry behavior; artifact build, encode, open, partial write,
flush/sync, replace, invalidation, and temporary-file cleanup; `finally`
blocks, resource close, reporter/output close, and interruption; and the
point where success becomes externally visible.

For every real operation: decide whether failure is command-fatal, a
warning, or secondary cleanup information, following acceptance criteria or
guardrails and never silently replacing or hiding the primary failure;
exercise failure before, during, and after its side effect when those states
differ; trace the operations actually entered against the intended order and
expected failure prefixes; verify semantic state, not only file existence or
matching counts; use an independent oracle rather than a second copy of the
implementation's own calculation; and mutate or bypass the general rule, not
only the original reproduction.

### Matrix Completion Lock

Do not begin the verdict or Issue write-up until every mandatory cell is
classified as passed, Issue, unverified, or not-applicable with a reason;
every prior remediation is reproduced inside the matrix or a named adjacent
cell; every Issue has been followed by completion of all remaining cells; the
acceptance-coverage classification has been reconciled against the matrix
results; and every applicable multi-step or side-effecting workflow satisfies
the Workflow And Side-Effect Check Lock above.

If time, tooling, or environment prevents completion, classify the affected
acceptance criterion as unverified and return `NEEDS WORK`; never silently
shrink the matrix.

## Parallel Planning Advice

Lane keys, `planned_reads`/`planned_writes` manifests and independence explanations are optional advice, and `features.parallel_planning` may be off. Ignoring a lane, running on main, using a different worktree, or changing files beyond a manifest is not a finding and never blocks `CLEAR`. Do not add a status, policy or gate for it, and never accept an executor's self-clearance. Evaluate it only when the scope's acceptance criteria make recommendations a feature: then verify their accuracy as behavior (dependencies respected, unexplained write/read overlap or unknown scope withholds advice, a stale explanation is not used). Malformed or stale advice is a nonblocking diagnostic. Actual-worktree filesystem and identity rules, dependencies and Guardrails stay in force.

## Perform The Adversarial Pass

Full mode only. Ask every applicable question below to challenge the completed coverage and workflow matrices:

- Can another constructor, factory, public API, serialization form, or environment path bypass validation?
- Can backward, skipped, overlapping, revived, or post-failure transitions violate the invariant?
- Can mode or representation changes bypass monotonicity, ownership, authorization, idempotency, or safety?
- Can redirected, no-colour, non-interactive, failure, retry, or timeout output lose required information or expose forbidden information?
- Do boundary, Unicode, parser, date/time, pagination, and numeric probes establish the relevant input class rather than one sample? Can duplicate, empty, missing, non-finite, mixed-type, or unexpected values create an impossible model or unhandled exception?
- Do tests use an independent outcome, rather than the implementation's own calculation and assumptions?

For state machines, structured events, parsers, renderers, auth, billing, persistence, jobs, and multi-step writes, verify relevant inputs and transitions through the required matrices; informal spot checks do not suffice.

## Re-check After Remediation

Use the immutable scope lock from the initial Check. Re-check every original
matrix cell (for a Quick Check, every original scope-lock item), every original reproduction, the remediation's changed code
paths, and only the adjacent cases already named in the original Issue or
scope lock, against the same focused and full gates.

Apply existing axes to remediation code, but do not add new axes, dependency
layers, or interpretations. A failure inside the frozen lock remains an
Issue. A newly noticed problem outside it is an observation for follow-up and
does not change the verdict, unless it is a credible blocker: secret
exposure, cross-tenant or role-boundary access, sensitive-data or privacy
harm, destructive data loss, unsafe billing side effects, or an equivalent
`.savepoint/Guardrails.md` blocker.

Before running a re-check probe, create an admission ledger with one row per
check: the re-check item, the prior Issue or remediation claim, the exact
frozen matrix cell, and the allowed result. Require an exact frozen cell for
every blocking result. A broad topic, general-purpose axis, changed helper,
plural configuration name, or newly noticed supported value is not enough.

If a probe has no exact frozen cell: do not run it as a blocking probe;
record a useful result as a non-blocking observation; do not count it toward
remediation required for closure; promote it only when the credible-blocker
exception above applies, naming the exact Blocker rule.

Start the re-check result with a closure map of the prior Issues: closed,
still open, or unverified.

Also assess every prior owner decision (acceptance or exception) on the scope.
Compare each decision's scope with what changed since its last assessment and
record the outcome as a `carried_forward` entry under `savepoint-check`
Closure Rules. A decision whose scope is unchanged still applies; state why
and do not ask the owner to renew it. List `unmet` requirement IDs on a
`NEEDS WORK` result so exception coverage can be checked.

Default convergence limit: one initial Check; one full re-check after
remediation; if an in-scope failure remains, one targeted remediation and
re-check; then stop and ask the owner to fix now, approve a permitted waiver,
create follow-up work, or close with non-blocking observations outside the
Task or Objective. Do not start a third autonomous remediation cycle or
broaden scope to keep a Check active. At the convergence limit, stop; do not
turn a newly noticed, non-blocking edge case into another remediation round.

## Verify File Reality

Every file named in a Task's evidence, drift note, or remediation claim must
either exist on disk now, be recorded as intentionally deleted, or be
recorded as discarded scratch work. Treat an unexplained phantom file as an
Issue.

## Verify Evidence And Gates

Run focused tests for changed behavior and relevant failure paths. Run
direct type or lint checks when the default gate excludes scoped files. Run
`git diff --check` and the project's configured build and test gates
(`quality_gates` in `.savepoint/config.yml`) unless the invoking skill names
a narrower approved gate. Apply the evidence mode the invoking skill
requires: Quick only for a requested Task Check, Full for the mandatory
Objective Check. Treat passing
tests and gates as supporting evidence, never as a substitute for acceptance
review.

## Collect Code Health Evidence

Full mode only. A Quick Task Check, and every other Savepoint activity, never
collects health.

1. Run the project's full gate first (see Verify Evidence And Gates). Test and
   coverage health instances read the reports that gate wrote, so tests do not
   run twice.
2. Run `savepoint health check O-### [dir]` for the Objective under review. It
   saves one `official` snapshot and prints the snapshot ID, whether a new
   snapshot was created, and a plain verdict. It exits non-zero only when no
   snapshot was saved; a blocking verdict still exits 0, so read the verdict.
   Never use a manual snapshot (from the board refresh) as Check evidence.
3. Record the printed snapshot ID in the Check frontmatter as
   `health_snapshot`. A project with no health configuration prints "Code
   Health is not configured"; record "Code Health not configured" in the Check
   body, omit `health_snapshot`, and do not treat it as a finding.
4. Treat the verdict as supporting evidence, never as `CLEAR` or as proof the
   Objective is healthy. Blockers include failing tests, high/critical vulnerabilities, required instances that failed or are stale, and opt-in `blocking` rules. A verdict that says "blocks clearance" prevents
   `CLEAR`: record the Issues behind it and result `NEEDS WORK`. A verdict that
   "does not block clearance" neither grants nor withholds `CLEAR`.
5. Keep collection failure, incomplete coverage, stale reports, unhealthy
   measurement, and Check findings as separate statements. A failed tool is not
   bad code.
6. Optional-instance failures, warnings, Watch results, and unknown-severity
   vulnerabilities do not block. Open an Issue for one only through the Issue
   Capture judgment below, when durable follow-up is warranted; collection
   never creates Issues.

## Complete The Issues Pass

Classify every acceptance criterion before returning a verdict: **Proven**
(independent evidence supports the general rule), **Issue** (a reproducible
scenario violates the rule), or **Unverified** (required evidence could not
be obtained). Any Issue or material unverified criterion means `NEEDS WORK`.
Return `CLEAR` only when every criterion is proven, relevant guardrails are
satisfied, and the required gates pass.

Before admitting an item as an Issue, require all of: it violates a named
acceptance criterion, guardrail, or release gate; it is reproducible through
a supported path; it lies inside the frozen scope lock; the Task or
Objective introduced, touched, widened, or explicitly promises the behavior;
and it has a credible consequence rather than only a theoretical
possibility. If an item fails this test, record it as an observation or omit
it.

Each Issue must include the violated acceptance criterion or guardrail rule,
the smallest reproducible scenario, expected and actual behavior, exact file
and line evidence, and the missing or inadequate test evidence.

Report every Issue from the completed pass together. Do not stop after the
first one. Do not invent requirements: Issues may rely only on acceptance
criteria, `.savepoint/Guardrails.md`, the active evidence mode, or explicit
release gates.

## Summarize Materiality

After the evidence-backed Issues, summarize every one in a compact
materiality table: Issue, Likelihood, Impact, Materiality, Recommendation.
Use Low, Medium, or High for likelihood, impact, and materiality, with a
short explanation where the rating is not self-evident.

Judge likelihood by realistic prerequisites, frequency, reachability, and
whether normal users or only unusual operator states can trigger it. Judge
impact by the credible outcome and existing containment, not only the
theoretical worst case. Combine likelihood and impact with the Task or
Objective's stated purpose and launch boundary to get materiality.
Recommendation states the proportionate disposition: fix now, combine with
another narrow fix, defer to named follow-up work, or accept with an
explicit owner waiver.

Do not copy the Issue order or guardrail severity into the materiality
rating without this separate check. Do not inflate a rare, contained
developer-workflow issue into a product-critical risk. Equally, do not use
low likelihood to excuse a credible blocker such as secret exposure,
cross-tenant access, destructive data loss, or sensitive-data harm.

Materiality guides priority and remediation scope; it does not silently
waive an acceptance criterion or guardrail. If the check shows that an item
does not actually violate an in-scope requirement, reclassify it as an
observation rather than leaving it as an Issue. If any Issue remains, the
verdict remains `NEEDS WORK` unless the owner explicitly approves the waiver
allowed by policy. When there are no Issues, state that no materiality
actions are required instead of emitting an empty table.

List observations separately from Issues. Do not use Issue language, `NEEDS
WORK`, or an in-task fix recommendation for an out-of-scope observation
unless the owner explicitly expands the Task or Objective.

## Review Code Style

When `.savepoint/Guardrails.md` defines `STYLE-*` rules, every Check record
carries a `## Code Style Review` section with one checkbox per rule, in the
order Guardrails lists them. Take rule IDs and labels from Guardrails; do not
hardcode them here. Tick a rule the scoped code follows. Leave a rule
unticked with a one-line reason and file evidence when the scoped code
departs from it. Style is advisory: an unticked rule never creates an Issue
or changes the result on its own. Skip the section when Guardrails defines
no `STYLE-*` rules.

```markdown
## Code Style Review

- [x] STYLE-01 **One job per file**
- [ ] STYLE-10 **Small diffs** — one commit adds ~2,000 lines across 30 files.
```
