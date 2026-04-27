import { existsSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";
import {
  PROJECT_TEMPLATES,
  PROMPT_TEMPLATES,
  RELEASE_TEMPLATES,
  type TemplateName,
} from "./manifest.js";

export function defaultTemplateRoot(): string {
  let dir = dirname(fileURLToPath(import.meta.url));
  while (true) {
    const candidate = join(dir, "templates");
    if (existsSync(candidate) && existsSync(join(candidate, "project"))) {
      return candidate;
    }
    const parent = dirname(dir);
    if (parent === dir) break;
    dir = parent;
  }
  throw new Error("Template directory not found");
}

export function templatePath(name: TemplateName, root?: string): string {
  const base = root ?? defaultTemplateRoot();
  if ((PROJECT_TEMPLATES as readonly string[]).includes(name)) {
    return join(base, "project", name);
  }
  if ((RELEASE_TEMPLATES as readonly string[]).includes(name)) {
    return join(base, "release", name);
  }
  if ((PROMPT_TEMPLATES as readonly string[]).includes(name)) {
    return join(base, "prompts", name);
  }
  throw new Error(`Unknown template: ${name}`);
}
