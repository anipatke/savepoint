import { parseTaskId, type TaskId, type Result } from "./ids.js";
import { parseTaskStatus, type TaskStatus } from "./status.js";

export interface TaskFrontmatter {
  id: TaskId;
  status: TaskStatus;
  objective: string;
  depends_on: TaskId[];
}

export interface TaskDocument {
  frontmatter: TaskFrontmatter;
  body: string;
}

export function validateTaskFrontmatter(raw: unknown): Result<TaskFrontmatter> {
  if (typeof raw !== "object" || raw === null || Array.isArray(raw)) {
    return { ok: false, reason: "Frontmatter must be a YAML mapping" };
  }

  const obj = raw as Record<string, unknown>;

  if (typeof obj["id"] !== "string") {
    return { ok: false, reason: "Missing or non-string 'id' field" };
  }
  const idResult = parseTaskId(obj["id"]);
  if (!idResult.ok) return idResult;

  if (typeof obj["status"] !== "string") {
    return { ok: false, reason: "Missing or non-string 'status' field" };
  }
  const statusResult = parseTaskStatus(obj["status"]);
  if (!statusResult.ok) return statusResult;

  if (typeof obj["objective"] !== "string" || obj["objective"].trim() === "") {
    return { ok: false, reason: "Missing or empty 'objective' field" };
  }

  const rawDeps = obj["depends_on"];
  if (!Array.isArray(rawDeps)) {
    return { ok: false, reason: "'depends_on' must be a YAML sequence" };
  }

  const deps: TaskId[] = [];
  for (let i = 0; i < rawDeps.length; i++) {
    const dep = rawDeps[i];
    if (typeof dep !== "string") {
      return {
        ok: false,
        reason: `'depends_on[${i}]' must be a string, got ${typeof dep}`,
      };
    }
    const depResult = parseTaskId(dep);
    if (!depResult.ok) {
      return { ok: false, reason: `'depends_on[${i}]': ${depResult.reason}` };
    }
    deps.push(depResult.value);
  }

  return {
    ok: true,
    value: {
      id: idResult.value,
      status: statusResult.value,
      objective: obj["objective"],
      depends_on: deps,
    },
  };
}
