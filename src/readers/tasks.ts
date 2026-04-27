import { readdir } from "node:fs/promises";
import { join } from "node:path";
import { readMarkdownFile, type FsError } from "../fs/markdown.js";
import { savepointPath } from "../fs/project.js";
import { validateTaskFrontmatter, type TaskDocument } from "../domain/task.js";
import {
  findDuplicateIds,
  findMissingDependencies,
  findCycle,
  type MissingDependency,
} from "../validation/dependencies.js";
import { type TaskId } from "../domain/ids.js";

export interface GraphErrors {
  duplicateIds: TaskId[];
  missingDependencies: MissingDependency[];
  cycle: TaskId[] | null;
}

export type EpicTaskSetResult =
  | { ok: true; tasks: TaskDocument[] }
  | { ok: false; kind: "fs_errors"; errors: FsError[] }
  | {
      ok: false;
      kind: "graph_errors";
      tasks: TaskDocument[];
      errors: GraphErrors;
    };

export async function readEpicTaskSet(
  root: string,
  releaseTag: string,
  epicRaw: string,
): Promise<EpicTaskSetResult> {
  const tasksDir = savepointPath(
    root,
    "releases",
    releaseTag,
    "epics",
    epicRaw,
    "tasks",
  );

  let entries: string[];
  try {
    entries = await readdir(tasksDir);
  } catch {
    return {
      ok: false,
      kind: "fs_errors",
      errors: [{ code: "not_found", path: tasksDir }],
    };
  }

  const mdFiles = entries.filter((e) => e.endsWith(".md")).sort();

  const results = await Promise.all(
    mdFiles.map((file) =>
      readMarkdownFile(join(tasksDir, file), validateTaskFrontmatter),
    ),
  );

  const fsErrors: FsError[] = [];
  const tasks: TaskDocument[] = [];

  for (const result of results) {
    if (result.ok) {
      tasks.push(result.value);
    } else {
      fsErrors.push(result.error);
    }
  }

  if (fsErrors.length > 0) {
    return { ok: false, kind: "fs_errors", errors: fsErrors };
  }

  const frontmatters = tasks.map((t) => t.frontmatter);
  const duplicateIds = findDuplicateIds(frontmatters);
  const missingDependencies = findMissingDependencies(frontmatters);
  const cycle = findCycle(frontmatters);

  if (
    duplicateIds.length > 0 ||
    missingDependencies.length > 0 ||
    cycle !== null
  ) {
    return {
      ok: false,
      kind: "graph_errors",
      tasks,
      errors: { duplicateIds, missingDependencies, cycle },
    };
  }

  return { ok: true, tasks };
}
