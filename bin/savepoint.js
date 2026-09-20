#!/usr/bin/env node

const { spawnSync } = require("node:child_process");
const path = require("node:path");

const platformMap = {
  darwin: "darwin",
  linux: "linux",
  win32: "windows",
};

const archMap = {
  arm64: "arm64",
  x64: "amd64",
};

const supportedTargets = Object.freeze([
  "darwin/amd64",
  "darwin/arm64",
  "linux/amd64",
  "linux/arm64",
  "windows/amd64",
  "windows/arm64",
]);

function resolveTarget(platform = process.platform, arch = process.arch) {
  const goos = platformMap[platform];
  const goarch = archMap[arch];
  if (!goos || !goarch) {
    throw new Error(
      `savepoint does not support ${platform}/${arch}; supported targets: ${supportedTargets.join(", ")}`,
    );
  }
  return {
    goos,
    goarch,
    executable: goos === "windows" ? "savepoint.exe" : "savepoint",
  };
}

function main() {
  let target;
  try {
    target = resolveTarget();
  } catch (error) {
    console.error(error.message);
    return 1;
  }

  const binary = path.join(
    __dirname,
    "..",
    "dist",
    "npm",
    `${target.goos}-${target.goarch}`,
    target.executable,
  );
  const result = spawnSync(binary, process.argv.slice(2), { stdio: "inherit" });

  if (result.error) {
    console.error(`could not launch ${target.goos}/${target.goarch}: ${result.error.message}`);
    return 1;
  }
  if (result.status === null) {
    console.error(`savepoint terminated by ${result.signal ?? "an unknown signal"}`);
    return 1;
  }
  return result.status;
}

if (require.main === module) {
  process.exit(main());
}

module.exports = { resolveTarget, supportedTargets };
