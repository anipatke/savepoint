# Audit Register

Current, mutable index of audit findings for this repo. The register is derived from the
immutable run records under `runs/` plus user and agent disposition work. Reconcile it on
every audit using `prompt.md`; never duplicate a finding that already has a stable `F###`
ID. Full finding records live in `findings/` (see `findings/README.md`).

## Convergence summary

Counts for the most recent reconciled state, so progress is visible run over run.

| Metric | Count |
|--------|-------|
| Net-new | 0 |
| Reopened | 0 |
| Verified | 0 |
| Deferred | 0 |
| Coverage gaps | 0 |

## Findings

One row per finding. `Status` uses the lifecycle values: `open`, `triaged`, `mapped`,
`in_progress`, `fixed`, `verified`, `deferred`, `owner_decision`, `waived`, `duplicate`.
`Proof` names the evidence for closure (required before `verified`).

| ID | Title | Status | Severity | Confidence | Work item | First seen | Last seen | Proof |
|----|-------|--------|----------|------------|-----------|------------|-----------|-------|
<!-- Add the first real finding as `F001`; do not keep placeholder rows in the register. -->
