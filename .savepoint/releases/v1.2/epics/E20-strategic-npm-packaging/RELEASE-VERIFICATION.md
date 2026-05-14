# Release Verification & Recovery — Strategic npm Packaging

Scope: the `Publish Package` workflow at `.github/workflows/publish.yml` and the npm package set produced by `npm run build` (root `savepoint` + six `savepoint-<os>-<arch>` platform packages).

## Publish Order

1. `verify` matrix runs on `ubuntu-latest`, `windows-latest`, and `macos-latest`. Each runner executes `make ci` and `make pack-smoke`, which exercises the packed root + host-platform tarball.
2. `publish` job runs only after every `verify` matrix leg succeeds (`needs: verify` + `fail-fast: true`).
3. `publish` rebuilds all platform artifacts, asserts each `dist/npm/<os-arch>/package.json` exists, then publishes platform packages in deterministic sorted order: `darwin-amd64, darwin-arm64, linux-amd64, linux-arm64, windows-amd64, windows-arm64`.
4. Root `savepoint` is published last so it never references a platform version that does not yet exist on the registry.

## Pre-Publish Checklist

- `package.json` `version` matches the tag and every `optionalDependencies.savepoint-*` entry uses the same version.
- Local `make build && make test` green.
- Local `make pack-smoke` green on at least the host platform.
- Tag pushed (`git push origin vX.Y.Z`) or `workflow_dispatch` triggered intentionally.

## Recovering From a Partial Publish

A failure between platform publishes and the root publish leaves some `savepoint-<os>-<arch>@X.Y.Z` versions live and the root unreleased. The root resolver is not affected because no end user has installed the new root yet.

Recovery options, in order of preference:

1. **Bump and republish.** Fastest, safest. Increment patch (`X.Y.Z+1`), update root `version` and every `optionalDependencies.*` entry to match, commit, tag, push. The previously published platform versions are simply unused. Do *not* `npm unpublish` — the 72-hour window plus dependency rules make this risky.
2. **Resume in place.** If the failure was transient (network, auth) and no code changed, re-run the failed `publish` job from the Actions UI. The platform publish step is idempotent only for versions that have *not* yet been published; already-published platform packages will return `EPUBLISHCONFLICT`. Add `|| true` is **not** acceptable — instead, edit the step locally to skip already-published versions:

   ```bash
   if npm view "$(node -p "require('./'+'$dir'+'package.json').name")@$(node -p "require('./'+'$dir'+'package.json').version")" version >/dev/null 2>&1; then
     echo "skipping already-published $dir"
     continue
   fi
   ```

   Use sparingly; prefer option 1.

## Recovering From a Failed Verify

`verify` failures block `publish` automatically. Read the failing matrix leg's logs:

- Go build errors → fix Go source, push a new tag.
- `make pack-smoke` failures on a single OS → likely platform-specific resolver/binary regression. Reproduce locally with `make pack-smoke` on that OS before re-tagging.
- Flaky network during `npm pack` → re-run the failed job; do not silently retry inside the script.

## Auth & Secret Hygiene

- `NODE_AUTH_TOKEN` is set only on the publish steps, never on `verify`.
- `id-token: write` is granted at workflow level for future npm provenance support; no other job uses it today.
- Never commit a `.npmrc` containing a token; CI writes auth via `setup-node`'s `registry-url` + `NODE_AUTH_TOKEN`.

## Manual Smoke After Publish

On each supported OS, in a clean directory:

```bash
npm install savepoint@X.Y.Z
npx savepoint --version
npx savepoint init ./tmp-fixture
```

A successful `--version` plus a populated `tmp-fixture/.savepoint/` confirms the platform package resolved correctly.
