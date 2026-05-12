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

const goos = platformMap[process.platform];
const goarch = archMap[process.arch];

if (!goos || !goarch) {
  console.error(`savepoint does not support ${process.platform}/${process.arch}`);
  process.exit(1);
}

const executable = goos === "windows" ? "savepoint.exe" : "savepoint";
const binary = path.join(__dirname, "..", "dist", "npm", `${goos}-${goarch}`, executable);
const result = spawnSync(binary, process.argv.slice(2), { stdio: "inherit" });

if (result.error) {
  console.error(result.error.message);
  process.exit(1);
}

process.exit(result.status ?? 0);
