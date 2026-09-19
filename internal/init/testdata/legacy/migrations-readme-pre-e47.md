# Savepoint Migrations

Archived copies of package-owned assets that Savepoint retired during
`savepoint upgrade-assets`. Files here are records, not instructions:
nothing in this directory is loaded by an agent or triggerable as a skill.

- `savepoint-audit-SKILL.md` — the generic audit skill retired when audit
  split into `savepoint-audit-task` (read-only review of one in-progress
  task) and `savepoint-audit-epic` (audit-pending closeout). Kept so local
  edits are recoverable. A numbered suffix means a differing copy was archived by
  a later upgrade.

Delete anything here once you have salvaged what you need.
