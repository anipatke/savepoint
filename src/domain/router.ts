import { type Result } from "./ids.js";

export const ROUTER_STATES = [
  "pre-implementation",
  "epic-design",
  "epic-task-breakdown",
  "task-planning",
  "task-building",
  "audit-pending",
] as const;

export type RouterStateValue = (typeof ROUTER_STATES)[number];

export interface RouterState {
  state: RouterStateValue;
  release: string;
  epic: string | undefined;
  next_action: string;
}

export function isRouterStateValue(value: string): value is RouterStateValue {
  return (ROUTER_STATES as readonly string[]).includes(value);
}

export function validateRouterState(raw: unknown): Result<RouterState> {
  if (typeof raw !== "object" || raw === null || Array.isArray(raw)) {
    return { ok: false, reason: "Router state must be a YAML mapping" };
  }
  const obj = raw as Record<string, unknown>;

  if (typeof obj["state"] !== "string") {
    return { ok: false, reason: "Missing or non-string 'state' field" };
  }
  if (!isRouterStateValue(obj["state"])) {
    return {
      ok: false,
      reason: `Unknown router state "${obj["state"]}": must be one of ${ROUTER_STATES.join(", ")}`,
    };
  }

  if (typeof obj["release"] !== "string" || obj["release"].trim() === "") {
    return { ok: false, reason: "Missing or empty 'release' field" };
  }

  const epic =
    typeof obj["epic"] === "string" && obj["epic"].trim() !== ""
      ? obj["epic"]
      : undefined;

  if (
    typeof obj["next_action"] !== "string" ||
    obj["next_action"].trim() === ""
  ) {
    return { ok: false, reason: "Missing or empty 'next_action' field" };
  }

  return {
    ok: true,
    value: {
      state: obj["state"],
      release: obj["release"],
      epic,
      next_action: obj["next_action"],
    },
  };
}
