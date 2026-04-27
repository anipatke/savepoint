---
version: { { RELEASE_NUMBER } }
name: "{{RELEASE_NAME}}"
status: in_progress
---

# Release v{{RELEASE_NUMBER}} — {{RELEASE_NAME}}

## What ships in v{{RELEASE_NUMBER}}

<!-- One paragraph describing the release theme and primary deliverables. -->

### In scope

<!-- Bullet list of features, commands, or capabilities included in this release. -->

### Out of v{{RELEASE_NUMBER}} scope (deferred to later releases)

<!-- Bullet list of what is explicitly not included, with target future release if known. -->

## Proposed epic breakdown

These are the proposed epics for v{{RELEASE_NUMBER}}. Order matters: each builds on the prior. The first epic is E01-scaffolding (per the savepoint convention). Epic folders use `E##-slug` prefixes so humans can navigate the filesystem without reopening this table.

| #   | Epic name | Purpose                         |
| --- | --------- | ------------------------------- |
| 01  | `E01-...` | <!-- Purpose of first epic. --> |

<!-- Add additional rows as needed. -->

## Success criteria for v{{RELEASE_NUMBER}}

<!-- How do you know the release is shippable? -->

## Risks tracked at this release level

<!-- What could derail the release and what is the mitigation? -->
