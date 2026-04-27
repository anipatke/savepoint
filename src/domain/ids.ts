type Ok<T> = { ok: true; value: T };
type Err = { ok: false; reason: string };
export type Result<T> = Ok<T> | Err;

export interface ReleaseId {
  tag: "release";
  raw: string;
}

export interface EpicId {
  tag: "epic";
  number: number;
  slug: string;
  raw: string;
}

export interface TaskId {
  tag: "task";
  epic: EpicId;
  number: number;
  slug: string;
  raw: string;
}

const RELEASE_RE = /^v\d+$/;
const EPIC_RE = /^(E(\d{2,})-([a-z][a-z0-9-]*))$/;
const TASK_SEGMENT_RE = /^T(\d{3,})-([a-z][a-z0-9-]*)$/;

export function parseReleaseId(raw: string): Result<ReleaseId> {
  if (!RELEASE_RE.test(raw)) {
    return {
      ok: false,
      reason: `Invalid release ID "${raw}": expected "v<number>" (e.g. "v1")`,
    };
  }
  return { ok: true, value: { tag: "release", raw } };
}

export function parseEpicId(raw: string): Result<EpicId> {
  const m = EPIC_RE.exec(raw);
  if (!m) {
    return {
      ok: false,
      reason: `Invalid epic ID "${raw}": expected "E<nn>-<slug>" (e.g. "E02-data-model")`,
    };
  }
  return {
    ok: true,
    value: { tag: "epic", number: parseInt(m[2], 10), slug: m[3], raw: m[1] },
  };
}

export function parseTaskId(raw: string): Result<TaskId> {
  const slash = raw.indexOf("/");
  if (slash === -1) {
    return {
      ok: false,
      reason: `Invalid task ID "${raw}": expected "<epic>/<task>" (e.g. "E02-data-model/T001-slug")`,
    };
  }
  const epicRaw = raw.slice(0, slash);
  const taskRaw = raw.slice(slash + 1);

  const epicResult = parseEpicId(epicRaw);
  if (!epicResult.ok) return epicResult;

  const m = TASK_SEGMENT_RE.exec(taskRaw);
  if (!m) {
    return {
      ok: false,
      reason: `Invalid task segment "${taskRaw}": expected "T<nnn>-<slug>" (e.g. "T001-domain-ids")`,
    };
  }

  return {
    ok: true,
    value: {
      tag: "task",
      epic: epicResult.value,
      number: parseInt(m[1], 10),
      slug: m[2],
      raw,
    },
  };
}

export function formatEpicId(id: EpicId): string {
  return id.raw;
}

export function formatTaskId(id: TaskId): string {
  return id.raw;
}
