---
id: O-028
title: Configure and run project-owned health tools safely
status: planned
depends_on: [O-026, O-027]
release: R-007
priority: high
rank: 1
---

# O-028: Configure and run project-owned health tools safely

## Outcome

New and upgraded projects can discover, confirm, configure, and run supported project-owned providers through a safe deterministic service without guessed execution or hidden installation.

## Why

Automatic setup should reduce expertise requirements without turning discovery into permission to execute arbitrary commands or making Savepoint a package manager.

## Success Conditions

- Initial scaffolding produces confirmation-ready suggestions from bounded inspection of manifests, lockfiles, established scripts, provider configuration, and known reports.
- `Design.md` receives human-readable intentions; dedicated configuration owns exact provider, executable, arguments, working directory, timeout, report path, required/optional policy, exclusions, and thresholds.
- Existing projects use an explicit upgrade and reconciliation flow; opening a project writes and scans nothing.
- Direct no-shell execution is sequential, cancellable, output-bounded, and governed by provider-specific default timeouts with project overrides.
- Absence, execution failure, timeout, cancellation, malformed reports, and unhealthy measurements remain distinct.
- Fresh gate artifacts are reused without duplicating test/build work.
- Generated, vendored, and third-party code is excluded by default through provider mechanisms; confirmed overrides affect the scope fingerprint.
- Monorepos support several scoped instances of a capability without invalid aggregation.

## Architectural Considerations

Discovery and execution belong to Code Health. Init, upgrades, Checks, and TUI call narrow services. External tools stay installed and owned by the project.

## Boundaries

**In scope:** discovery proposals, configuration, confirmation, safe processes, timeouts, cancellation, artifact freshness, scope, multi-instance support, scaffolding, and upgrades.

**Out of scope:** installing tools, downloads, shell strings, background work, analysis logic, dashboard rendering, and implicit project writes.

## Confirmed Design Decisions

The owner confirmed discovery-and-suggest, Design/config separation, direct structured execution, sequential scheduling, provider timeouts, artifact reuse, default exclusions, polyglot support, and explicit upgrades on 2026-09-26.
