---
type: audit-method-reference
triggerable: false
---

# Shared Savepoint Audit Method

This reference is not a skill and never triggers on its own. It is loaded in
full by `savepoint-audit-task` and `savepoint-audit-epic`, which own their
triggers, scope, health-check mode, and output contracts.

Read and apply this method completely for both task and epic audits.

Where this method names `.savepoint/Guardrails.md` or `.savepoint/Health-Check.md`,
use them when the project has them and skip the related step when they are
absent; their absence is not a finding.

## Establish Scope

1. Inspect current files and the diff. Treat context logs, checked boxes, prior
   audit notes, remediation claims, and green tests as claims to verify, not
   proof.
2. Verify file reality for every scoped file before referring to it.
3. Load only the acceptance criteria, guardrail rules, design context, and
   health-check requirements that govern the scoped work.

## Freeze The Audit Scope

Before the first adversarial probe, trace the real code and workflow, then write
a numbered scope lock containing:

1. the acceptance criteria, guardrails, and release gates being tested;
2. the changed files and supported public entry points;
3. the directly relied-on runtime orchestration, external effects, and dependency
   behavior needed for those entry points to keep their promise;
4. the selected matrix axes and cells, including explicit not-applicable cells;
5. the supported-path and materiality boundary used to admit findings.

The scope lock defines what the audit means by `complete`. Do not expand into
unrelated dependency internals, unsupported configurations, or increasingly
remote hypothetical states.

During an initial audit only, new source evidence may correct a factual mistake
in the scope lock. Record the amendment and restart the affected matrix pass.
Once the initial verdict is returned, the lock is immutable for every re-audit.
Do not silently introduce a new axis, dependency layer, acceptance
interpretation, or meaning of "adjacent" during remediation review.

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

Always include these axes when applicable:

1. **Public surfaces:** every constructor, factory, parser/validator, state
   reducer, pure renderer, output backend, and selection/helper entry point.
2. **Input shape:** normal, empty, missing, duplicate, wrong type, mixed type,
   non-finite, mutable input, and mutation after validation.
3. **State:** every state; every allowed transition; backward, skipped,
   overlapping, terminal-revival, post-failure, and representation-switch
   transitions.
4. **Environment/output:** interactive and redirected sinks crossed with normal,
   no-colour, dumb/no-cursor, failure, warning, timeout, and explicit public
   overrides. Exercise both the recommended factory and direct public backend
   construction.
5. **Boundaries:** the exact limit, immediately below and above it, and the full
   small finite range when practical. Do not test only named sample points.
6. **Sequences:** complete workflows, not isolated frames. Include initialization,
   intermediate milestones, completion, failure, retry/repeat, and final output.
7. **Representations:** serialization round trips, direct models, structured
   events, and rendered output where each exists.
8. **Text classes:** ASCII, control characters, combining marks, variation
   selectors, wide characters, emoji modifiers, regional flags, and joined emoji
   when text width or truncation is in scope. Use an independent oracle when one
   is available; otherwise use structural assertions that do not reuse the
   implementation's calculation.

For progress renderers and state-machine output, the minimum matrix is:

- every public model/event/view entry point;
- TTY and non-TTY crossed with `NO_COLOR`, `TERM=dumb`, default detection, and
  explicit overrides;
- all model states and valid/invalid transition classes;
- every integer width from the hard floor through the comfortable minimum,
  plus immediately below/at/above the maximum;
- counted progress across 0%, each milestone boundary on both sides, 100%, and
  completion with and without a result word;
- timed progress start, repeat, advance, equal value, rewind, timeout boundary,
  and counted/timed representation switches;
- malformed values and containers, including mutation after construction;
- the complete Unicode corpus listed above;
- interactive, redirected, ready, warning, failure, and no-colour final output.

Run the matrix through repeatable parameterized tests or one deterministic audit
harness where possible. Record its rows, cell classifications, and command/output
evidence. Ad-hoc probes may supplement the matrix but cannot replace it.

Once the initial audit starts its adversarial probes, this matrix is the frozen
scope lock. A missing required cell discovered during that initial audit is an
audit-process error: amend the lock explicitly and rerun the affected matrix
before returning a verdict. A re-audit may never add the missing cell
retroactively as a new blocking perimeter.

### Finite External-Boundary Matrix

When scoped code relies on a server, subprocess, browser runner, provider, or
other external boundary, classify this finite set during the initial audit:

- configured target and actual runtime target;
- startup, discovery, and reuse ordering;
- connection refusal or unavailable dependency;
- successful and non-success responses;
- redirect behavior;
- timeout and cancellation;
- malformed or unexpected response;
- retry, cleanup, and partial side effects where applicable;
- secret-safe failure output.

Mark a cell not-applicable only with a scope reason. Do not invent additional
network or toolchain edge classes during later re-audits unless they meet the
credible-blocker exception below.

## Workflow And Side-Effect Audit Lock

Apply this lock to any command or workflow with multiple operations, external
calls, persistence, transactions, generated artifacts, cleanup, or structured
progress. Derive the inventory from the actual code path and external effects,
not from task notes, an implementation-owned lifecycle registry, or the test
table that is meant to verify it.

Before testing, record this table for the complete workflow:

| Order | Real operation | Side effect / state change | Failure timing | Failure owner and final state | Cleanup / secondary failure | Independent oracle |
|---|---|---|---|---|---|---|

The inventory must include, when applicable:

- input parsing, guards, setup, constructors, and connection/client creation;
- external calls and each operation that can partly succeed;
- transaction begin, write, commit, rollback, and retry behavior;
- artifact build, encode, open, partial write, flush/sync, replace, invalidation,
  and temporary-file cleanup;
- `finally` blocks, resource close, reporter/output close, and interruption;
- the point where success becomes externally visible, including terminal output,
  exit code, structured event, manifest, cache, or database state.

For every real operation:

1. Decide whether failure is command-fatal, a warning, or secondary cleanup
   information. The decision must follow acceptance criteria or guardrails and
   must never silently replace or hide the primary failure.
2. Exercise failure before, during, and after its side effect when those states
   differ. Include partial writes/updates and a primary failure combined with a
   rollback, close, or cleanup failure.
3. Trace the operations actually entered. Compare the success trace with the
   intended order and each failure trace with its expected prefix and cleanup.
4. Verify semantic state, not only file existence, valid syntax, or matching
   counts. Generated artifacts must describe the current external/database
   state; an old but well-formed artifact may still be false.
5. Use an independent oracle. Two registries, tables, reducers, or calculations
   copied from the implementation may agree while sharing the same omission.
   Implementation-owned declarations may support evidence but cannot define the
   audit scope or prove their own completeness.
6. For security filters, parsers, and redactors, derive a corpus from the full
   accepted input grammar: quoted/unquoted values, whitespace, escapes, encoded
   forms, mixed surrounding text, and configured secret values where applicable.
7. Mutate or bypass the general rule, not only the original reproduction. A
   mutation check is useful only if an adjacent omission would also fail.

Classify every operation row and every applicable failure timing as passed,
finding, unverified, or not-applicable with a scope reason. A workflow cannot
return `CLEAR` while a real operation is missing from the inventory, an output
can disagree with external state, a secondary failure is silently swallowed, or
the oracle is self-referential.

### Matrix Completion Lock

Do not begin the verdict or final findings write-up until:

1. every mandatory cell is classified as passed, finding, unverified, or
   not-applicable with a reason;
2. every prior remediation is reproduced inside the matrix or a named adjacent
   matrix cell;
3. every finding has been followed by completion of all remaining cells;
4. the acceptance-coverage classification has been reconciled against the
   matrix results; and
5. every applicable multi-step or side-effecting workflow satisfies the
   Workflow And Side-Effect Audit Lock above.

If time, tooling, or environment prevents completion, classify the affected
acceptance criterion as unverified and return `NEEDS WORK`; never silently shrink
the matrix.

## Perform The Adversarial Pass

Ask every applicable question:

- Can validation be bypassed through another constructor, factory, direct
  public API, serialization form, or environment path?
- Can state move backward, skip forward, overlap, revive after completion, or
  continue after failure?
- Can switching modes or representations bypass a monotonicity, ownership,
  authorization, idempotency, or safety check?
- Can redirected, no-colour, non-interactive, failure, retry, or timeout output
  lose required information or expose forbidden information?
- What happens immediately below and above every documented limit?
- Does a Unicode, parser, date/time, pagination, or numeric example cover the
  whole input class, or only the tested sample?
- Can duplicate, empty, missing, non-finite, mixed-type, or unexpected values
  create an impossible model or unhandled exception?
- Are tests checking an independent outcome, or reusing the implementation's
  own calculation and assumptions?

For state machines, structured events, parsers, renderers, auth boundaries,
billing, persistence, jobs, or multi-step writes, build a small behavior matrix
for relevant inputs and transitions. Do not rely on informal spot checks.

## Re-audit After Remediation

Use the immutable scope lock from the initial audit. Re-audit:

1. every original matrix cell;
2. every original reproduction;
3. the remediation's changed code paths;
4. only the adjacent cases already named in the original finding or scope lock;
5. the same focused and full gates.

Apply existing axes to remediation code, but do not add new axes, dependency
layers, or interpretations. A failure inside the frozen lock remains a finding.
A newly noticed issue outside it is an observation for follow-up and does not
change the verdict, unless it is a credible blocker: secret exposure,
cross-tenant or role-boundary access, sensitive-data or privacy harm,
destructive data loss, unsafe billing side effects, or an equivalent
`.savepoint/Guardrails.md` blocker.

Before running a re-audit probe, create an admission ledger with one row per
check:

| Re-audit check | Prior finding or remediation claim | Exact frozen matrix cell | Allowed result |
|---|---|---|---|

Require an exact frozen cell for every blocking result. A broad topic,
general-purpose axis, changed helper, plural configuration name, or newly
noticed supported value is not enough. Applying an existing axis to remediation
code means rechecking only the values, classes, and adjacent cases recorded in
the initial scope lock; it does not authorize filling an omitted initial matrix
cell later.

If a probe has no exact frozen cell:

- do not run it as a blocking probe;
- record a useful result as a non-blocking observation;
- do not include it in the finding count or remediation required for closure;
- promote it only when the credible-blocker exception above applies, naming the
  exact Blocker rule.

Start the re-audit result with a closure map of the prior findings: `closed`,
`still open`, or `unverified`. Do not renumber an observation as a new finding.
If remediation closes the recorded cells but owner-run evidence is still
impossible until commit, push, or deploy, report repository readiness separately
from the overall audit verdict.

Default convergence limit:

1. one initial audit;
2. one full re-audit after remediation;
3. if an in-scope failure remains, one targeted remediation and re-audit;
4. then stop and ask the owner to fix now, approve a permitted waiver, create
   follow-up work, or close with non-blocking observations outside the task.

Do not start a third autonomous remediation cycle or broaden scope to keep an
audit active. The owner may explicitly request further work, but that request
must name whether it extends the scope lock.

At the convergence limit, stop. Do not turn a newly noticed, non-blocking edge
case into another remediation round.

## Verify File Reality

Every file named in a task context log, drift note, or remediation claim must
either exist on disk now, be recorded as intentionally deleted, or be recorded
as discarded scratch work. Treat an unexplained phantom file as a finding.

## Verify Evidence And Gates

1. Run focused tests for changed behavior and relevant failure paths.
2. Run direct type or lint checks when the default gate excludes scoped files.
3. Run `git diff --check`, `make build`, and `make test` unless the active
   health-check instructions define a narrower approved gate.
4. Apply the health-check mode required by the invoking skill: Quick for task
   audit, Full for epic audit.
5. Treat passing tests and gates as supporting evidence, never as a substitute
   for acceptance review.

## Complete The Findings Pass

Classify every acceptance criterion before returning a verdict:

- **Proven:** independent evidence supports the general rule.
- **Finding:** a reproducible scenario violates the rule.
- **Unverified:** required evidence could not be obtained.

Any finding or material unverified criterion means `NEEDS WORK`. Return `CLEAR`
only when every criterion is proven, relevant guardrails are satisfied, and the
required gates pass.

Before admitting an item as a finding, require all of:

1. it violates a named acceptance criterion, guardrail, or release gate;
2. it is reproducible through a supported path;
3. it lies inside the frozen scope lock;
4. the task introduced, touched, widened, or explicitly promises the behavior;
5. it has a credible consequence rather than only a theoretical possibility.

If an item fails this test, record it as an observation or omit it. Observations
do not change `CLEAR`/`NEEDS WORK`, do not expand remediation scope, and should
name follow-up work only when useful. The credible-blocker exception in the
re-audit section overrides the frozen perimeter, but the auditor must state
which blocker rule makes the exception valid.

Each finding must include:

- the violated acceptance criterion or guardrail rule;
- the smallest reproducible scenario;
- expected and actual behavior;
- exact file and line evidence; and
- the missing or inadequate test evidence.

Report all findings from the completed pass together. Do not stop after the
first issue. Do not invent requirements: findings may rely only on acceptance
criteria, `.savepoint/Guardrails.md`, active Savepoint evidence gates, or
explicit release gates.

## Summarize Materiality

After the evidence-backed findings, summarize every finding in one compact
materiality table:

| Finding | Likelihood | Impact | Materiality | Recommendation |
|---|---|---|---|---|

Use `Low`, `Medium`, or `High` for likelihood, impact, and materiality, with a
short explanation where the rating is not self-evident.

- **Likelihood:** judge the realistic prerequisites, frequency, reachability,
  and whether normal users or only unusual operator states can trigger it.
- **Impact:** judge the credible outcome and existing containment, not only the
  theoretical worst case.
- **Materiality:** combine likelihood and impact with the task or epic's stated
  purpose and launch boundary.
- **Recommendation:** state the proportionate disposition, such as fix now,
  combine with another narrow fix, defer to named follow-up work, or accept with
  an explicit owner waiver.

Do not copy the finding order or guardrail severity into the materiality rating
without this separate check. Do not inflate a rare, contained developer-workflow
issue into a product-critical risk. Equally, do not use low likelihood to excuse
a credible blocker such as secret exposure, cross-tenant access, destructive
data loss, or sensitive-data harm.

Materiality guides priority and remediation scope; it does not silently waive an
acceptance criterion or guardrail. If the check shows that an item does not
actually violate an in-scope requirement, reclassify it as an observation rather
than leaving it as a finding. If any finding remains, the verdict remains
`NEEDS WORK` unless the owner explicitly approves the waiver allowed by policy.
When there are no findings, state that no materiality actions are required
instead of emitting an empty table.

List observations separately from findings. Do not use finding language,
`NEEDS WORK`, or an in-task fix recommendation for an out-of-scope observation
unless the owner explicitly expands the task.
