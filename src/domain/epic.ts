import { type Result } from "./ids.js";

export interface EpicFrontmatter {
  type: string;
  status: string;
}

export function validateEpicFrontmatter(raw: unknown): Result<EpicFrontmatter> {
  if (typeof raw !== "object" || raw === null || Array.isArray(raw)) {
    return { ok: false, reason: "Frontmatter must be a YAML mapping" };
  }
  const obj = raw as Record<string, unknown>;

  if (typeof obj["type"] !== "string" || obj["type"].trim() === "") {
    return { ok: false, reason: "Missing or empty 'type' field" };
  }
  if (typeof obj["status"] !== "string" || obj["status"].trim() === "") {
    return { ok: false, reason: "Missing or empty 'status' field" };
  }

  return {
    ok: true,
    value: { type: obj["type"], status: obj["status"] },
  };
}
