# Frozen legacy fixtures

These files are pinned copies of shapes Savepoint has already written into real
projects. They are **byte-frozen**: never reformat them, never re-sync them with
the current templates, and never "fix up" their wording.

A test that fails against a fixture reports a compatibility break in the reader
or the upgrade path, not a stale fixture. Fixing such a failure means changing
the code — or making a deliberate, documented migration — not editing the file
here.

| File | Provenance |
| --- | --- |
| `router.md` | A pre-v1.5 project router: no `defect` key, a `phase`-era vocabulary in its prose, and an extra `updated` key the current reader does not model. |
| `AGENTS.unmarked.md` | A hand-written agent guide from before the `SAVEPOINT:BEGIN`/`END` marker pair existed. |
| `AGENTS.marked.md` | An agent guide carrying the marker pair around an older managed block, with user prose on both sides. |
| `SKILL.customized.md` | A bundled skill a project has locally tailored. |
| `task-legacy-phase.md` | A task file using the legacy `phase` frontmatter field instead of `stage`. |
