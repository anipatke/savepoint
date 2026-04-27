import { type TaskId } from "../domain/ids.js";
import { type TaskFrontmatter } from "../domain/task.js";

export interface MissingDependency {
  dependent: TaskId;
  missing: TaskId;
}

export function findDuplicateIds(tasks: TaskFrontmatter[]): TaskId[] {
  const seen = new Set<string>();
  const reported = new Set<string>();
  const dupes: TaskId[] = [];
  for (const task of tasks) {
    const raw = task.id.raw;
    if (seen.has(raw) && !reported.has(raw)) {
      dupes.push(task.id);
      reported.add(raw);
    } else {
      seen.add(raw);
    }
  }
  return dupes;
}

export function findMissingDependencies(
  tasks: TaskFrontmatter[],
): MissingDependency[] {
  const known = new Set(tasks.map((t) => t.id.raw));
  const missing: MissingDependency[] = [];
  for (const task of tasks) {
    for (const dep of task.depends_on) {
      if (!known.has(dep.raw)) {
        missing.push({ dependent: task.id, missing: dep });
      }
    }
  }
  return missing;
}

// Returns the cycle as an ordered list of task IDs (last entry equals first, showing the loop).
// Returns null when no cycle exists.
export function findCycle(tasks: TaskFrontmatter[]): TaskId[] | null {
  const taskMap = new Map(tasks.map((t) => [t.id.raw, t]));
  // 0 = unvisited, 1 = in stack, 2 = done
  const color = new Map<string, 0 | 1 | 2>();
  for (const task of tasks) color.set(task.id.raw, 0);

  const path: TaskId[] = [];

  function visit(raw: string): TaskId[] | null {
    const task = taskMap.get(raw);
    if (!task) return null;
    color.set(raw, 1);
    path.push(task.id);
    for (const dep of task.depends_on) {
      if (!taskMap.has(dep.raw)) continue;
      const c = color.get(dep.raw) ?? 0;
      if (c === 1) {
        const idx = path.findIndex((t) => t.raw === dep.raw);
        return [...path.slice(idx), dep];
      }
      if (c === 0) {
        const result = visit(dep.raw);
        if (result) return result;
      }
    }
    path.pop();
    color.set(raw, 2);
    return null;
  }

  for (const task of tasks) {
    if ((color.get(task.id.raw) ?? 0) === 0) {
      const result = visit(task.id.raw);
      if (result) return result;
    }
  }
  return null;
}
