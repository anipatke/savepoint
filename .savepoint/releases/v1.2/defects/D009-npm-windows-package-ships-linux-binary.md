---
id: v1.2/D009-npm-windows-package-ships-linux-binary
release: v1.2
status: resolved
severity: high
title: "npm package ships Linux binary as savepoint.exe on Windows"
---

# D009: npm package ships Linux binary as savepoint.exe on Windows

## Symptom

Running `npx savepoint upgrade-assets` from a Windows project fails before Savepoint starts:

```text
This version of C:\Users\User\Branding\03-VIBE-LAB\payment-initiation\node_modules\savepoint\savepoint.exe is not compatible with the version of Windows you're running. Check your computer's system information and then contact the software publisher.
```

The installed package is `savepoint@1.2.2`, and its `package.json` maps the CLI bin to `./savepoint.exe`.

Investigation found the installed file at `payment-initiation\node_modules\savepoint\savepoint.exe` starts with an ELF header (`7F 45 4C 46`), which is a Linux executable format. A Windows executable should start with an `MZ` header.

## Expected Behavior

Windows users who install or run Savepoint through npm should receive a Windows-compatible executable, so `npx savepoint upgrade-assets` can launch normally.

## Reproduction

1. On Windows, install or invoke `savepoint@1.2.2` through npm/npx in an existing project.
2. Run:

   ```powershell
   npx savepoint upgrade-assets
   ```

3. Windows reports that `node_modules\savepoint\savepoint.exe` is not compatible.
4. Inspecting the first bytes of the installed `savepoint.exe` shows an ELF binary rather than a Windows PE binary.

## Impact

Windows users cannot run the npm-published Savepoint CLI, including `upgrade-assets`, because npm installs a Linux binary under a Windows `.exe` path. This blocks v1.2 asset upgrades for existing Windows projects.

## Fix Plan

- Update release packaging so npm does not publish a host-built binary under a universal `savepoint.exe` name.
- Ensure the Windows npm path receives a `GOOS=windows` binary, and non-Windows platforms receive the correct executable for their platform.
- Add a packaging or smoke check that verifies the npm artifact's executable header/platform before publish.

## Acceptance Criteria

- [x] On Windows, `npx savepoint upgrade-assets` launches without the Windows compatibility error.
- [x] The npm-installed Windows executable starts with a Windows PE `MZ` header.
- [x] The publish process cannot produce a package where `bin.savepoint` points to a Linux ELF binary named `savepoint.exe` for Windows users.

## Resolution Notes

Resolved by routing npm builds through `go run ./internal/buildtool build-npm`, which cross-compiles all supported npm runtime binaries under `dist/npm/<os>-<arch>/`.

Changed the npm CLI entrypoint to `bin/savepoint.js`, which selects the correct bundled binary for the current platform and architecture. Added a publish guard that reads each generated Windows executable header and fails the npm build unless it starts with the Windows `MZ` signature. This prevents a Linux ELF binary from being published under the Windows `.exe` path.

Verification:

- `go test ./internal/buildtool` passed.
- `npm.cmd run build` passed.
- `Format-Hex -Path dist\npm\windows-amd64\savepoint.exe -Count 16` showed `4D 5A` (`MZ`) at offset 0.
- `node .\bin\savepoint.js --version` launched the Windows npm binary and printed `v1.2.2`.
- `node .\bin\savepoint.js upgrade-assets C:\Users\User\Branding\03-VIBE-LAB\payment-initiation --dry-run` launched successfully and produced an upgrade report.
- `make build` passed.
- `make test` still fails on `TestBundledSavepointSkillsHaveDiscoveryFrontmatter` because `agent-skills\savepoint-audit\SKILL.md` is missing YAML frontmatter; that failure is outside this defect repair.
