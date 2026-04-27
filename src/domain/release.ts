import { type Result } from "./ids.js";

export interface ReleaseFrontmatter {
  version: number;
  name: string;
  status: string;
}

export function validateReleaseFrontmatter(
  raw: unknown,
): Result<ReleaseFrontmatter> {
  if (typeof raw !== "object" || raw === null || Array.isArray(raw)) {
    return { ok: false, reason: "Frontmatter must be a YAML mapping" };
  }
  const obj = raw as Record<string, unknown>;

  if (typeof obj["version"] !== "number") {
    return { ok: false, reason: "Missing or non-numeric 'version' field" };
  }
  if (typeof obj["name"] !== "string" || obj["name"].trim() === "") {
    return { ok: false, reason: "Missing or empty 'name' field" };
  }
  if (typeof obj["status"] !== "string" || obj["status"].trim() === "") {
    return { ok: false, reason: "Missing or empty 'status' field" };
  }

  return {
    ok: true,
    value: {
      version: obj["version"],
      name: obj["name"],
      status: obj["status"],
    },
  };
}
