export const TASK_STATUSES = [
  "backlog",
  "planned",
  "in_progress",
  "review",
  "done",
] as const;
export type TaskStatus = (typeof TASK_STATUSES)[number];

export type Result<T> = { ok: true; value: T } | { ok: false; reason: string };

const ALLOWED_TRANSITIONS: Record<TaskStatus, readonly TaskStatus[]> = {
  backlog: ["planned"],
  planned: ["in_progress"],
  in_progress: ["review", "planned"],
  review: ["done", "in_progress"],
  done: ["review"],
};

export function isTaskStatus(value: string): value is TaskStatus {
  return (TASK_STATUSES as readonly string[]).includes(value);
}

export function parseTaskStatus(raw: string): Result<TaskStatus> {
  if (!isTaskStatus(raw)) {
    return {
      ok: false,
      reason: `Unknown task status "${raw}": must be one of ${TASK_STATUSES.join(", ")}`,
    };
  }
  return { ok: true, value: raw };
}

export function isTransitionAllowed(from: TaskStatus, to: TaskStatus): boolean {
  return (ALLOWED_TRANSITIONS[from] as readonly string[]).includes(to);
}

export function validateTransition(
  from: TaskStatus,
  to: TaskStatus,
): Result<{ from: TaskStatus; to: TaskStatus }> {
  if (!isTransitionAllowed(from, to)) {
    const allowed = ALLOWED_TRANSITIONS[from];
    const allowedStr =
      allowed.length > 0 ? allowed.join(", ") : "(none — terminal state)";
    return {
      ok: false,
      reason: `Transition "${from}" → "${to}" is not allowed. Allowed transitions from "${from}": ${allowedStr}`,
    };
  }
  return { ok: true, value: { from, to } };
}
