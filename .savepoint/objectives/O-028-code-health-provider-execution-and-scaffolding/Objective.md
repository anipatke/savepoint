---
id: O-028
title: Configure and run project-owned health tools safely
status: done
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

On 2026-10-01, the owner confirmed the O-028 detail design:

- **Setup flow.** A new human-only command, `savepoint health setup [dir]`, previews what it found and what it would configure without writing. `--apply` writes `.savepoint/health/config.json`. `savepoint init` prints the same preview at the end. Re-running setup on a configured project is the upgrade and reconciliation path: it lists new suggestions and configured tools whose inputs or executables went missing, and changes nothing without `--apply`. Opening a project, the board, resume, and doctor never discover, scan, or write.
- **Discovery** reads a bounded set of project files (manifests, lockfiles, test and coverage settings, provider configuration, known report paths) and checks PATH for executables. It never runs a tool, installs anything, or uses the network.
- **Execution.** Savepoint runs only the configured analysis tools (Lizard, jscpd, OSV-Scanner): sequentially, as an executable plus argument vector with no shell, cancellable, with bounded output. It never runs tests. Tests and coverage read the report the project's own gate already produced; a report older than the relevant inputs it describes is stale, a missing one is absent.
- **Required policy.** Each configured tool is optional by default; the owner may set `required: true`. Optional failures are reported, never hidden.
- **Default exclusions.** Setup proposes vendored, generated, and third-party exclusions (for example `vendor/**`, `node_modules/**`, `**/*.pb.go`) for confirmation and passes them through each provider's own exclusion mechanism. Exclusions are part of the scope fingerprint, so a change starts a new series.
- **Monorepos.** Several instances of one capability and provider are allowed when each has a distinct `name` and scope; instances are never averaged or summed together.
- **Vulnerabilities.** Setup suggests OSV-Scanner's normal online mode and states plainly that the scanner, not Savepoint, sends package names and versions to OSV.dev.
- **Design.md.** Setup writes only the health configuration. The design skill describes confirmed tools in Design.md in plain words; setup never edits Design.md.

Technical decisions recorded with that confirmation: provider default timeouts are 60 seconds for Lizard and jscpd and 120 seconds for OSV-Scanner, with project overrides up to 600 seconds; the config and snapshot schemas stay at version 1 (nothing has shipped) while gaining the instance `name` and `required` fields; report normalization is O-029's work, so O-028 defines the per-provider reader seam and tests it with fakes, and a provider without a reader yields `unsupported`; O-028 adds no collection command, which arrives with O-030 (Check) and O-031 (TUI refresh).
