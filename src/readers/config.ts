import { readFile } from "node:fs/promises";
import { load } from "js-yaml";
import { savepointPath } from "../fs/project.js";
import { applyConfigDefaults, type SavepointConfig } from "../domain/config.js";
import { type FsResult } from "../fs/markdown.js";

export async function readConfig(
  root: string,
): Promise<FsResult<SavepointConfig>> {
  const path = savepointPath(root, "config.yml");

  let raw: string;
  try {
    raw = await readFile(path, "utf8");
  } catch {
    return { ok: false, error: { code: "not_found", path } };
  }

  let parsed: unknown;
  try {
    parsed = load(raw);
  } catch (err) {
    const reason = err instanceof Error ? err.message : String(err);
    return {
      ok: false,
      error: { code: "malformed_frontmatter", path, reason },
    };
  }

  return { ok: true, value: applyConfigDefaults(parsed) };
}
