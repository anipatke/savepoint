// Claude Code SessionStart hook: tells Claude the current Savepoint Next line.
// Adds one line of context. It adds nothing and exits 0 when savepoint is
// missing, `resume` fails or times out, or this is not a Savepoint project.
"use strict";

const fs = require("node:fs");
const path = require("node:path");
const { spawnSync } = require("node:child_process");
const { findSavepoint } = require("./savepoint-find.js");

const TIMEOUT_MS = 5000;

function nextLine(projectDir, env = process.env) {
  if (!fs.existsSync(path.join(projectDir, ".savepoint"))) return null;
  const bin = findSavepoint(projectDir, env);
  if (!bin) return null;
  // A .cmd shim only runs through a shell, which joins the command line
  // unquoted; quote the path so a directory with a space still works.
  const shell = bin.toLowerCase().endsWith(".cmd");
  const result = spawnSync(shell ? `"${bin}"` : bin, ["resume"], {
    cwd: projectDir,
    encoding: "utf8",
    timeout: TIMEOUT_MS,
    shell,
    windowsHide: true,
  });
  if (result.error || result.status !== 0) return null;
  const line = result.stdout.split(/\r?\n/).find((l) => l.startsWith("Next"));
  return line ? line.trim() : null;
}

function main() {
  try {
    const projectDir = process.env.CLAUDE_PROJECT_DIR || process.cwd();
    const line = nextLine(projectDir);
    if (!line) return;
    process.stdout.write(
      JSON.stringify({
        hookSpecificOutput: {
          hookEventName: "SessionStart",
          additionalContext: `Savepoint ${line}`,
        },
      }),
    );
  } catch {
    // A hook must never break a session.
  }
}

if (require.main === module) main();

module.exports = { nextLine };
