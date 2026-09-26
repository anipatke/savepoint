---
id: O-026
title: Select trustworthy Code Health providers
status: planned
depends_on: []
release: R-007
priority: critical
rank: 1
---

# O-026: Select trustworthy Code Health providers

## Outcome

Savepoint has a confirmed modular boundary and an evidence-backed catalogue selecting one official project-owned tool or standard report path for each of the six Code Health capabilities.

## Why

Provider choices must not smuggle language engines, mandatory runtimes, unsafe downloads, or fragile parsing into Savepoint.

## Success Conditions

- Candidates for tests, coverage, complexity, duplication, change hotspots, and dependency vulnerabilities are compared for supported stacks, deterministic machine output, licence, runtime, platforms, offline behavior, performance, report size, and installation ownership.
- One official provider or report path is selected per capability with limitations and fixture strategy.
- Duplication support does not make Node or another ecosystem a mandatory Savepoint runtime.
- Complexity and change-hotspot analysis remain external rather than becoming Savepoint analysis engines.
- The narrow discovery, execution, normalization, snapshot, and dashboard interfaces owned by the separate Code Health module are settled.
- Compatibility rules identify which provider and configuration changes begin a new comparison series.

## Architectural Considerations

Code Health is compiled into Savepoint but owns its domain. Project-owned tools are invoked through structured executable and argument definitions without a shell. Scaffolding chooses from the supported catalogue and requires confirmation.

## Boundaries

**In scope:** bounded provider research, decisions, output contracts, fixtures, modular interfaces, packaging and licence implications, and unsupported cases.

**Out of scope:** installing scanners, production adapters, analysis algorithms, dashboard work, Check integration, custom-provider APIs, and a generic plugin framework.

## Confirmed Design Decisions

The owner confirmed the project-owned-tool model, one official path per capability, a supported catalogue, no automatic installation, no AI assessment, and no compromise to Savepoint's single-binary distribution on 2026-09-26.
