#!/usr/bin/env node
"use strict";

const { spawnSync } = require("node:child_process");
const { resolvePlatform } = require("../lib/resolve-platform.js");
const { locateBinary } = require("../lib/locate-binary.js");

function main() {
  const resolved = resolvePlatform({
    platform: process.platform,
    arch: process.arch,
  });
  if (!resolved.supported) {
    console.error(resolved.error);
    process.exit(1);
  }

  const binary = locateBinary(resolved);
  if (!binary) {
    console.error(
      `savepoint platform package ${resolved.packageName} is not installed; reinstall savepoint to fetch the binary for ${resolved.goos}/${resolved.goarch}`,
    );
    process.exit(1);
  }

  const result = spawnSync(binary, process.argv.slice(2), { stdio: "inherit" });
  if (result.error) {
    console.error(result.error.message);
    process.exit(1);
  }
  process.exit(result.status ?? 0);
}

main();
