---
id: v1.2/D010-skill-frontmatter-crlf-ci-failure
release: v1.2
status: resolved
severity: high
title: "make ci fails on savepoint-audit skill frontmatter with CRLF line endings"
reference: E18-template-skill-optimisation
---

# D010: make ci fails on savepoint-audit skill frontmatter with CRLF line endings

## Symptom

GitHub Actions fails during `make ci` at the repo root:

```text
--- FAIL: TestBundledSavepointSkillsHaveDiscoveryFrontmatter (0.00s)
    agent_skills_test.go:11: agent-skills\savepoint-audit\SKILL.md missing YAML frontmatter
FAIL
FAIL    github.com/opencode/savepoint    0.140s
make: *** [Makefile:9: test] Error 1
Error: Process completed with exit code 1.
```

The reported skill file visibly contains YAML frontmatter, but its first bytes are `2D 2D 2D 0D 0A` (`---\r\n`). The test currently checks `strings.HasPrefix(text, "---\n")`, so CRLF files are incorrectly reported as missing frontmatter.

## Expected Behavior

`make ci` should pass when a bundled Savepoint skill has valid YAML frontmatter regardless of whether the file uses LF or CRLF line endings.

## Reproduction

1. Run:

   ```powershell
   make ci
   ```

2. Observe `go test ./...` fail in `TestBundledSavepointSkillsHaveDiscoveryFrontmatter`.
3. Inspect the beginning of `agent-skills/savepoint-audit/SKILL.md`:

   ```powershell
   Format-Hex -Path agent-skills/savepoint-audit/SKILL.md -Count 16
   ```

4. The file starts with `---\r\n`, which is valid frontmatter but does not match the test's strict `---\n` prefix check.

## Impact

The publish workflow is blocked because `make ci` fails before package verification and publish can complete. The failure is misleading: the file has frontmatter, but the test treats Windows-style line endings as absent frontmatter.

## Fix Plan

- Update `agent_skills_test.go` to normalize CRLF to LF before checking frontmatter markers, or parse frontmatter using the shared markdown/frontmatter helper instead of a raw string prefix.
- Keep the bundled and scaffolded skill consistency check intact.
- Add or preserve coverage proving skill files with CRLF line endings are accepted when their frontmatter is otherwise valid.

## Acceptance Criteria

- [ ] `TestBundledSavepointSkillsHaveDiscoveryFrontmatter` passes for `---\r\n` skill files.
- [ ] `make ci` passes for the current repository checkout.
- [ ] The test still fails for a skill file that truly lacks YAML frontmatter.

## Resolution Notes

Pending.
