// Claude Code PreToolUse hook: blocks three owner-authority breaches with a
// one-line reason. It adds nothing to context unless a rule fires, and it
// allows (exits 0, no output) on any error of its own so it never wedges a
// session. Managed by Savepoint and refreshed on upgrade.
//   1. An edit that sets a Task to `status: done` (only the owner may).
//   2. An edit to .savepoint/router.md inside a git worktree lane.
//   3. An agent run of a `savepoint` command other than resume, create-task,
//      or health check.
"use strict";

const fs = require("node:fs");
const path = require("node:path");
const { spawnSync } = require("node:child_process");

const TASK_FILE = /(^|\/)\.savepoint\/objectives\/[^/]+\/tasks\/T-[^/]*\.md$/;
const ROUTER_FILE = /(^|\/)\.savepoint\/router\.md$/;

const MSG_DONE = "Blocked: only the owner may set a Task to status: done; leave it for the owner on the board.";
const MSG_ROUTER = "Blocked: do not edit .savepoint/router.md in a worktree lane; the owner sets the router on main after merging.";
const MSG_CLI = "Blocked: only the owner runs this savepoint command (agents may run resume, create-task, and health check); the owner can type it with !.";

function normalize(file) {
  return String(file || "").replace(/\\/g, "/");
}

// Reports whether the YAML frontmatter of a Task file says `status: done`.
function isDone(text) {
  const match = /^---\r?\n([\s\S]*?)\r?\n---/.exec(String(text || ""));
  return !!match && /^status:\s*['"]?done['"]?\s*(#.*)?$/m.test(match[1]);
}

function applyEdit(text, edit) {
  const from = String(edit.old_string ?? "");
  const to = String(edit.new_string ?? "");
  if (from === "" || !text.includes(from)) return text;
  return edit.replace_all ? text.split(from).join(to) : text.replace(from, () => to);
}

// The file content that would exist after the tool call, or null if unknown.
function resultingText(tool, input, current) {
  if (tool === "Write") return String(input.content ?? "");
  if (tool === "Edit") return applyEdit(current, input);
  if (tool === "MultiEdit") return (input.edits || []).reduce(applyEdit, current);
  return null;
}

function readOrEmpty(file) {
  try {
    return fs.readFileSync(file, "utf8");
  } catch {
    return "";
  }
}

function checkTaskDone(tool, input, cwd) {
  const file = path.resolve(cwd, String(input.file_path || ""));
  const current = readOrEmpty(file);
  const next = resultingText(tool, input, current);
  return next !== null && isDone(next) && !isDone(current) ? MSG_DONE : null;
}

function gitDir(cwd, flag) {
  const result = spawnSync("git", ["rev-parse", flag], { cwd, encoding: "utf8", windowsHide: true });
  if (result.error || result.status !== 0) return null;
  return path.resolve(cwd, result.stdout.trim());
}

function inLane(cwd) {
  const dir = gitDir(cwd, "--git-dir");
  const common = gitDir(cwd, "--git-common-dir");
  if (!dir || !common) return false;
  return fs.realpathSync(dir) !== fs.realpathSync(common);
}

function checkRouter(input, cwd) {
  const dir = path.dirname(path.resolve(cwd, String(input.file_path || "")));
  return inLane(fs.existsSync(dir) ? dir : cwd) ? MSG_ROUTER : null;
}

// Splits a shell command into word lists, one per simple command.
function simpleCommands(command) {
  const words = [];
  const commands = [];
  let word = "";
  let quote = "";
  let has = false;
  const endWord = () => {
    if (has) words.push(word);
    word = "";
    has = false;
  };
  const endCommand = () => {
    endWord();
    if (words.length) commands.push(words.splice(0));
  };
  const text = String(command || "");
  for (let i = 0; i < text.length; i++) {
    const c = text[i];
    if (quote) {
      if (c === quote) quote = "";
      else word += c;
    } else if (c === '"' || c === "'") {
      quote = c;
      has = true;
    } else if (/\s/.test(c) && c !== "\n") {
      endWord();
    } else if (c === "\n" || c === ";" || c === "|" || c === "&" || c === "(" || c === ")" || c === "`") {
      endCommand();
    } else {
      word += c;
      has = true;
    }
  }
  endCommand();
  return commands;
}

function isSavepoint(word) {
  return /^savepoint(\.exe|\.cmd)?$/i.test(path.basename(normalize(word)));
}

// Returns the savepoint arguments of a simple command, or null if it is not one.
function savepointArgs(words) {
  let i = 0;
  while (i < words.length && /^[A-Za-z_][A-Za-z0-9_]*=/.test(words[i])) i++;
  if (/^(npx|bunx|pnpx)(\.cmd)?$/i.test(path.basename(normalize(words[i])))) {
    i++;
    while (i < words.length && words[i].startsWith("-")) i++;
  } else if (/^(pnpm|yarn)$/i.test(words[i]) && /^(dlx|exec)$/.test(words[i + 1] || "")) {
    i += 2;
  }
  return i < words.length && isSavepoint(words[i]) ? words.slice(i + 1) : null;
}

function allowedArgs(args) {
  const rest = args.filter((a) => !a.startsWith("-"));
  const [sub, second] = rest;
  return sub === "resume" || sub === "create-task" || (sub === "health" && second === "check");
}

function checkBash(input) {
  for (const words of simpleCommands(input.command)) {
    const args = savepointArgs(words);
    if (args && !allowedArgs(args)) return MSG_CLI;
  }
  return null;
}

// Returns the block reason for a PreToolUse payload, or null to allow.
function decide(payload) {
  const tool = payload.tool_name;
  const input = payload.tool_input || {};
  const cwd = payload.cwd || process.env.CLAUDE_PROJECT_DIR || process.cwd();
  if (tool === "Bash") return checkBash(input);
  if (tool !== "Write" && tool !== "Edit" && tool !== "MultiEdit") return null;
  const file = normalize(input.file_path);
  if (TASK_FILE.test(file)) return checkTaskDone(tool, input, cwd);
  if (ROUTER_FILE.test(file)) return checkRouter(input, cwd);
  return null;
}

function main() {
  try {
    const payload = JSON.parse(fs.readFileSync(0, "utf8"));
    const reason = decide(payload);
    if (!reason) return;
    process.stdout.write(
      JSON.stringify({
        hookSpecificOutput: {
          hookEventName: "PreToolUse",
          permissionDecision: "deny",
          permissionDecisionReason: reason,
        },
      }),
    );
  } catch {
    // Fail open: a hook must never break a session.
  }
}

if (require.main === module) main();

module.exports = { decide, isDone, simpleCommands, savepointArgs };
