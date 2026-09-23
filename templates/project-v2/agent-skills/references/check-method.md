---
type: check-method-reference
triggerable: false
---

# Shared Savepoint Check Method

This reference is not a skill and never triggers on its own. It is loaded in
full by `savepoint-check`, which owns the Check trigger, scope, evidence mode,
and output contract. `savepoint-task` and `savepoint-design` reference it but
do not run it themselves.

Read and apply this method completely whenever a Check is run. A Task Check is
optional and runs only when requested or selected by the owner; the Full
Objective Check is mandatory at Objective closure, and the Full Goal Check is
mandatory whenever a Goal exists.

Goal Checks retain the existing serialized `scope.kind: release` and `R-###`
identity. The V2 board presents these optional contexts as Goals; `g` is the
canonical selector and `r` remains an undisplayed compatibility alias.

Where this method names `.savepoint/Guardrails.md` or an optional project
verification procedure, use them when the project has them and skip the
related step when they are absent; absence is not an Issue.

## Task Check And Objective Check Depth

A Task Check is focused and optional: one Task's outcome and evidence against
its own acceptance criteria, plan, and scoped files. If the owner skips this
local Check, the Task evidence must record an explicit waiver naming the Task,
reason, actor, and time. A waiver is not technical `CLEAR` and does not waive
an acceptance criterion or guardrail. It does satisfy a downstream Task
dependency that requires `clear` — the owner's own completion decision
stands in there — but never one that requires `accepted`, since there is no
Check for the owner to have accepted.

An Objective Check does everything a Task Check does, plus integration across
the Objective's Tasks and reconciliation against Design. It is mandatory and
must inspect every owned Task, including Tasks whose optional Task Check was
waived. A Task-only clearance never substitutes for the Objective's own Check.

A Goal Check is also mandatory whenever a Goal exists. It uses Full
evidence to evaluate every member Objective and cross-Objective integration;
the owner's acceptance of that exact current Check remains a separate step.

## Quick And Full Evidence Modes

Both modes apply the method below at different reach:

- **Quick** — run only when the optional Task Check is requested. Scope is the
  one Task: its acceptance criteria, its scoped files, and directly relevant
  `.savepoint/Guardrails.md` rules. Quick evidence is never an automatic gate
  at every Task handoff.
- **Full** — run for the mandatory Objective Check. Scope adds every member
  Task's outcome, cross-Task integration, and reconciliation against Design.
  Use the same Full depth for the mandatory Goal Check, extending scope to
  every member Objective and cross-Objective integration.

Both modes apply `.savepoint/Guardrails.md` when the project has it and skip
that step when it is absent. Both modes may run an optional project
verification procedure when the project defines one and skip that step when
it does not exist. Neither absence is an Issue; it is a skipped step.

Skipping an optional Task Check is not itself an Issue when the owner waiver is
present. The mandatory Objective Check still evaluates the Task's evidence and
any material guardrail or integration risk.

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
3. Check the normal case, boundary values, malformed input, failure behavior,
   and at least one bypass path.
4. Run at least one independent scenario that is not merely an existing unit
   test repeated unchanged.
5. Record the expected result, actual result, and concrete evidence.

A regression test proves its example. It does not prove the surrounding
invariant.

## Build The Mandatory Coverage Matrix

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

## Perform The Adversarial Pass

Ask every applicable question: can validation be bypassed through another
constructor, factory, direct public API, serialization form, or environment
path? Can state move backward, skip forward, overlap, revive after
completion, or continue after failure? Can switching modes or
representations bypass a monotonicity, ownership, authorization, idempotency,
or safety check? Can redirected, no-colour, non-interactive, failure, retry,
or timeout output lose required information or expose forbidden information?
What happens immediately below and above every documented limit? Does a
Unicode, parser, date/time, pagination, or numeric example cover the whole
input class, or only the tested sample? Can duplicate, empty, missing,
non-finite, mixed-type, or unexpected values create an impossible model or
unhandled exception? Are tests checking an independent outcome, or reusing
the implementation's own calculation and assumptions?

For state machines, structured events, parsers, renderers, auth boundaries,
billing, persistence, jobs, or multi-step writes, build a small behavior
matrix for relevant inputs and transitions. Do not rely on informal spot
checks.

## Re-check After Remediation

Use the immutable scope lock from the initial Check. Re-check every original
matrix cell, every original reproduction, the remediation's changed code
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
`git diff --check`, `make build`, and `make test` unless the invoking skill
names a narrower approved gate. Apply the evidence mode the invoking skill
requires: Quick only for a requested Task Check, Full for the mandatory
Objective or Goal Check. Treat passing
tests and gates as supporting evidence, never as a substitute for acceptance
review.

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
