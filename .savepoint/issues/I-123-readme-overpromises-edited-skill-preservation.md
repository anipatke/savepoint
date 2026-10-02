---
id: I-123
title: README overpromises edited skill preservation
type: drift
status: resolved
source:
  kind: check
  check: C-961
  actor: {role: checker, session: check-o034-20261003}
  at: '2026-10-02T20:24:00Z'
tasks: [T-100]
checks: [C-961, C-962]
guardrail_ids: [TPL-02]
severity: medium
resolution:
  disposition: verified
  check: C-962
  actor: {role: checker, session: recheck-o034-20261003}
  at: '2026-10-02T20:29:00Z'
  reason: Both README claims now explain the pre-manifest replacement and backup exception; frozen upgrade matrix and fresh full gate pass.
history:
  - at: '2026-10-02T20:24:00Z'
    actor: {role: checker, session: check-o034-20261003}
    kind: observed
    check: C-961
    note: Independent real-template upgrade matrix reproduced replacement of an authored pre-manifest skill without force; a backup preserves the original.
  - at: '2026-10-02T20:29:00Z'
    actor: {role: checker, session: recheck-o034-20261003}
    kind: rechecked
    check: C-962
    note: Tracked clean, tracked edited and pre-manifest dry-run/apply/repeat scenarios agree with both revised README statements. Fresh full gate and official health pass; verified closure.
---

# I-123: README overpromises edited skill preservation

## Summary

T-100 added an unconditional promise that edited skills are never overwritten
without `--force`. Supported pre-manifest upgrades deliberately replace the
skill after backing it up. This violates TPL-02: guidance must describe current
behavior. The finding concerns new documentation, not a demand to change the
existing upgrade policy.

## Evidence

- `README.md:153` promises an edited skill is kept with the revision beside it
  unless force is passed; `README.md:436`–438 says it is never overwritten
  without force. Both statements were introduced in b27c17a.
- `internal/init/upgrade.go:487`–495 selects `updated` plus `.bak` when the
  skill has no manifest entry. The tracked edited branch at 502–504 conflicts.
- Independent scratch overlay `TestO034IndependentUpgradeMatrix/pre-manifest`:
  create a schema-2 temporary project, write an authored task skill, omit its
  manifest entry, invoke `upgradeAssetsFromTree(realTemplates, dir, false,
  false)`. Actual: live skill equals current package bytes; `.bak` equals
  original authored bytes. Expected from README: original live bytes remain,
  incoming text offered as `.new`. Repeat is unchanged; dry-run writes nothing.
- Existing `TestUpgradeSkill_untrackedFileIsBackedUpAndReplaced` and
  `TestUpgradeProjectAssets_legacyProjectWithUnmarkedGuide` already assert the
  actual policy. New `TestUpgradeKeepsAnEditedSkillAndOffersTheRevision` covers
  only a manifest-tracked edit, so it cannot substantiate the README's absolute
  promise.

## Proof Needed

Qualify both README claims to distinguish manifest-tracked edited skills from
pre-manifest skills, explaining replacement with a recoverable `.bak` for the
latter. Re-run the frozen R9 tracked-unedited/tracked-edited/pre-manifest ×
dry-run/apply/repeat cells against the wording; verify the named existing
legacy-policy tests and required gate. Preserve existing runtime behavior.
