"use strict";

const test = require("node:test");
const assert = require("node:assert/strict");

const {
  resolvePlatform,
  platformPackageName,
  executableName,
  SUPPORTED_TARGETS,
} = require("../lib/resolve-platform.js");
const { locateBinary } = require("../lib/locate-binary.js");

test("resolvePlatform maps linux/x64 to savepoint-linux-amd64", () => {
  const r = resolvePlatform({ platform: "linux", arch: "x64" });
  assert.equal(r.supported, true);
  assert.equal(r.packageName, "savepoint-linux-amd64");
  assert.equal(r.executable, "savepoint");
  assert.equal(r.binaryPath, "bin/savepoint");
  assert.equal(r.slug, "linux-amd64");
});

test("resolvePlatform maps darwin/arm64 to savepoint-darwin-arm64", () => {
  const r = resolvePlatform({ platform: "darwin", arch: "arm64" });
  assert.equal(r.packageName, "savepoint-darwin-arm64");
  assert.equal(r.binaryPath, "bin/savepoint");
});

test("resolvePlatform maps win32/x64 to savepoint-windows-amd64 with .exe", () => {
  const r = resolvePlatform({ platform: "win32", arch: "x64" });
  assert.equal(r.packageName, "savepoint-windows-amd64");
  assert.equal(r.executable, "savepoint.exe");
  assert.equal(r.binaryPath, "bin/savepoint.exe");
});

test("resolvePlatform maps win32/arm64 to savepoint-windows-arm64", () => {
  const r = resolvePlatform({ platform: "win32", arch: "arm64" });
  assert.equal(r.packageName, "savepoint-windows-arm64");
});

test("resolvePlatform returns unsupported for freebsd/x64", () => {
  const r = resolvePlatform({ platform: "freebsd", arch: "x64" });
  assert.equal(r.supported, false);
  assert.match(r.error, /freebsd\/x64/);
});

test("resolvePlatform returns unsupported for linux/mips", () => {
  const r = resolvePlatform({ platform: "linux", arch: "mips" });
  assert.equal(r.supported, false);
  assert.match(r.error, /linux\/mips/);
});

test("SUPPORTED_TARGETS covers six expected Go targets", () => {
  const keys = SUPPORTED_TARGETS.map((t) => `${t.goos}/${t.goarch}`).sort();
  assert.deepEqual(keys, [
    "darwin/amd64",
    "darwin/arm64",
    "linux/amd64",
    "linux/arm64",
    "windows/amd64",
    "windows/arm64",
  ]);
});

test("platformPackageName composes from goos and goarch", () => {
  assert.equal(platformPackageName("linux", "amd64"), "savepoint-linux-amd64");
});

test("executableName adds .exe only on windows", () => {
  assert.equal(executableName("windows"), "savepoint.exe");
  assert.equal(executableName("linux"), "savepoint");
  assert.equal(executableName("darwin"), "savepoint");
});

test("locateBinary returns installed platform package path when resolver succeeds", () => {
  const resolved = resolvePlatform({ platform: "linux", arch: "x64" });
  const got = locateBinary(resolved, {
    resolver: (id) => `/node_modules/${id}`,
    existsSync: () => false,
  });
  assert.equal(got, "/node_modules/savepoint-linux-amd64/bin/savepoint");
});

test("locateBinary falls back to dev dist when resolver throws and file exists", () => {
  const resolved = resolvePlatform({ platform: "linux", arch: "x64" });
  const got = locateBinary(resolved, {
    resolver: () => {
      throw new Error("not installed");
    },
    existsSync: () => true,
    devRoot: "/repo/dist/npm",
  });
  assert.equal(
    got.replace(/\\/g, "/"),
    "/repo/dist/npm/linux-amd64/bin/savepoint",
  );
});

test("locateBinary returns null when resolver throws and dev file missing", () => {
  const resolved = resolvePlatform({ platform: "linux", arch: "x64" });
  const got = locateBinary(resolved, {
    resolver: () => {
      throw new Error("not installed");
    },
    existsSync: () => false,
    devRoot: "/repo/dist/npm",
  });
  assert.equal(got, null);
});

test("locateBinary uses windows binary path with .exe for win32", () => {
  const resolved = resolvePlatform({ platform: "win32", arch: "x64" });
  const got = locateBinary(resolved, {
    resolver: (id) => `/node_modules/${id}`,
    existsSync: () => false,
  });
  assert.equal(
    got,
    "/node_modules/savepoint-windows-amd64/bin/savepoint.exe",
  );
});
