import { readFile } from "node:fs/promises";
import { load } from "js-yaml";

export type FsError =
  | { code: "not_found"; path: string }
  | { code: "missing_frontmatter"; path: string }
  | { code: "malformed_frontmatter"; path: string; reason: string }
  | { code: "schema_rejection"; path: string; reason: string };

export type FsResult<T> =
  | { ok: true; value: T }
  | { ok: false; error: FsError };

export interface MarkdownDoc<T> {
  frontmatter: T;
  body: string;
}

const FRONTMATTER_RE = /^---\r?\n([\s\S]*?)\r?\n---\r?\n?([\s\S]*)$/;

export async function readMarkdownFile<T>(
  path: string,
  validate: (
    raw: unknown,
  ) => { ok: true; value: T } | { ok: false; reason: string },
): Promise<FsResult<MarkdownDoc<T>>> {
  let raw: string;
  try {
    raw = await readFile(path, "utf8");
  } catch {
    return { ok: false, error: { code: "not_found", path } };
  }

  const match = FRONTMATTER_RE.exec(raw);
  if (!match) {
    return { ok: false, error: { code: "missing_frontmatter", path } };
  }

  const [, yamlBlock, body] = match;
  let parsed: unknown;
  try {
    parsed = load(yamlBlock);
  } catch (err) {
    const reason = err instanceof Error ? err.message : String(err);
    return {
      ok: false,
      error: { code: "malformed_frontmatter", path, reason },
    };
  }

  const result = validate(parsed);
  if (!result.ok) {
    return {
      ok: false,
      error: { code: "schema_rejection", path, reason: result.reason },
    };
  }

  return {
    ok: true,
    value: { frontmatter: result.value, body: body.trimEnd() },
  };
}
