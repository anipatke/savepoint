// Shared by the Savepoint Claude Code hooks: finds the savepoint executable.
// Order: PATH, then node_modules/.bin under the project. Returns null when it
// is not found. No dependencies; managed by Savepoint and refreshed on upgrade.
"use strict";

const fs = require("node:fs");
const path = require("node:path");

function candidates(dir, windows) {
  return windows
    ? [path.join(dir, "savepoint.exe"), path.join(dir, "savepoint.cmd")]
    : [path.join(dir, "savepoint")];
}

function isFile(file) {
  try {
    return fs.statSync(file).isFile();
  } catch {
    return false;
  }
}

function findSavepoint(projectDir, env = process.env, platform = process.platform) {
  const windows = platform === "win32";
  const pathVar = env.PATH || env.Path || "";
  const dirs = pathVar.split(path.delimiter).filter(Boolean);
  dirs.push(path.join(projectDir, "node_modules", ".bin"));
  for (const dir of dirs) {
    for (const file of candidates(dir, windows)) {
      if (isFile(file)) return file;
    }
  }
  return null;
}

module.exports = { findSavepoint };
