const assert = require("node:assert/strict");
const test = require("node:test");

const { resolveTarget, supportedTargets } = require("./savepoint.js");

test("resolves every supported platform and architecture", () => {
  const inputs = [
    ["linux", "x64", "linux", "amd64", "savepoint"],
    ["linux", "arm64", "linux", "arm64", "savepoint"],
    ["darwin", "x64", "darwin", "amd64", "savepoint"],
    ["darwin", "arm64", "darwin", "arm64", "savepoint"],
    ["win32", "x64", "windows", "amd64", "savepoint.exe"],
    ["win32", "arm64", "windows", "arm64", "savepoint.exe"],
  ];

  for (const [platform, arch, goos, goarch, executable] of inputs) {
    assert.deepEqual(resolveTarget(platform, arch), { goos, goarch, executable });
  }
  assert.equal(supportedTargets.length, 6);
});

test("reports unsupported platform and architecture with the supported matrix", () => {
  assert.throws(
    () => resolveTarget("freebsd", "x64"),
    /savepoint does not support freebsd\/x64; supported targets: .*windows\/arm64/,
  );
});
