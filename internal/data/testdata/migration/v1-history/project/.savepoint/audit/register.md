# Audit Register

Current, mutable index of audit findings for this frozen example project. Reconciled
once, on 2026-07-01, against `runs/2026-07-01-example.md`.

## Convergence summary

| Metric | Count |
|--------|-------|
| Net-new | 3 |
| Reopened | 0 |
| Verified | 0 |
| Deferred | 0 |
| Coverage gaps | 1 |

## Findings

| ID | Title | Status | Severity | Confidence | Work item | First seen | Last seen | Proof |
|----|-------|--------|----------|------------|-----------|------------|-----------|-------|
| F001 | Shared-epic dependency resolution lacked regression coverage | fixed | medium | high | E01-example/T001-shared | 2026-06-20 | 2026-07-01 | A regression test proving the short reference resolves within the active release |
| F002 | Health-Check custom procedure omits Deep mode | waived | low | medium | E01-example/T001-shared | 2026-06-22 | 2026-07-01 | n/a; owner accepted the gap |
| F003 | Duplicate of shared-epic dependency resolution finding | duplicate | medium | high | E01-example/T001-shared | 2026-06-25 | 2026-07-01 | n/a; canonical finding is F001 |
