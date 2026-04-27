import { readFile } from "node:fs/promises";
import { templatePath } from "./paths.js";
import type { TemplateName } from "./manifest.js";

export type TemplateError = { code: "not_found"; path: string };

export type TemplateResult<T> =
  | { ok: true; value: T }
  | { ok: false; error: TemplateError };

export async function loadTemplate(
  name: TemplateName,
  options: { root?: string } = {},
): Promise<TemplateResult<string>> {
  const path = templatePath(name, options.root);
  try {
    const content = await readFile(path, "utf8");
    return { ok: true, value: content };
  } catch {
    return { ok: false, error: { code: "not_found", path } };
  }
}
