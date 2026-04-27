---
id: E04-templates-and-prompts/T004-template-integrity-tests
status: done
objective: "Add integrity tests that lock required template files, markers, router states, and render behavior."
depends_on: ["E04-templates-and-prompts/T003-template-registry-renderer"]
---

# T004: template-integrity-tests

## Implementation Plan

- [x] Add `test/templates/project-templates.test.ts` covering required project template files, frontmatter where expected, and embedded agent instruction markers.
- [x] Add `test/templates/router-template.test.ts` covering required router states, read-order guidance, and the conditional visual-identity read instruction.
- [x] Add `test/templates/prompt-templates.test.ts` covering the required prompt template set and workflow-specific markers.
- [x] Add an audit prompt assertion that E04 onward audit guidance uses a single `.savepoint/audit/{epic}/proposals.md` bundle with delta-only edits where possible.
- [x] Add render assertions proving supported variables interpolate and unresolved required variables fail clearly.
- [x] Run `npm test -- test/templates` or the closest focused Vitest command available for the new tests.
