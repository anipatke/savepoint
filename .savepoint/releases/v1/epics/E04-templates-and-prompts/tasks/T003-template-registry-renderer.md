---
id: E04-templates-and-prompts/T003-template-registry-renderer
status: done
objective: "Add typed template lookup and rendering helpers for the markdown and YAML assets."
depends_on:
  [
    "E04-templates-and-prompts/T001-project-template-assets",
    "E04-templates-and-prompts/T002-release-and-prompt-assets",
  ]
---

# T003: template-registry-renderer

## Implementation Plan

- [x] Add `src/templates/paths.ts` or an equivalent module that centralizes template root and template path resolution.
- [x] Add a typed template manifest for project templates, release templates, and prompt templates without duplicating large file contents in TypeScript.
- [x] Add a read helper that loads a named template from disk and returns path-aware boundary errors for missing assets.
- [x] Add a small interpolation helper for supported variables such as project name and release number.
- [x] Export the public template helper surface from a narrow `src/templates` entrypoint.
- [x] Add focused unit tests for manifest lookup, missing-template failures, and interpolation behavior.
- [x] Run the focused template helper tests and typecheck.
