import { readFile } from "node:fs/promises";
import { load } from "js-yaml";
import { savepointPath } from "../fs/project.js";
import { validateRouterState, type RouterState } from "../domain/router.js";
import { type FsResult } from "../fs/markdown.js";

// Matches the ```yaml ... ``` code block under the "## Current state" heading.
const STATE_BLOCK_RE = /## Current state\s+```yaml\r?\n([\s\S]*?)\r?\n```/;

export async function readRouterState(
  root: string,
): Promise<FsResult<RouterState>> {
  const path = savepointPath(root, "router.md");

  let raw: string;
  try {
    raw = await readFile(path, "utf8");
  } catch {
    return { ok: false, error: { code: "not_found", path } };
  }

  const match = STATE_BLOCK_RE.exec(raw);
  if (!match) {
    return { ok: false, error: { code: "missing_frontmatter", path } };
  }

  let parsed: unknown;
  try {
    parsed = load(match[1]);
  } catch (err) {
    const reason = err instanceof Error ? err.message : String(err);
    return {
      ok: false,
      error: { code: "malformed_frontmatter", path, reason },
    };
  }

  const result = validateRouterState(parsed);
  if (!result.ok) {
    return {
      ok: false,
      error: { code: "schema_rejection", path, reason: result.reason },
    };
  }

  return { ok: true, value: result.value };
}
