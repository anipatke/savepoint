import { existsSync } from "node:fs";
import { join, dirname } from "node:path";

export type Result<T> = { ok: true; value: T } | { ok: false; reason: string };

const SAVEPOINT_DIR = ".savepoint";

export function findSavepointRoot(startDir: string): Result<string> {
  let dir = startDir;
  while (true) {
    if (existsSync(join(dir, SAVEPOINT_DIR))) {
      return { ok: true, value: dir };
    }
    const parent = dirname(dir);
    if (parent === dir) {
      return {
        ok: false,
        reason: `No ${SAVEPOINT_DIR}/ directory found searching up from "${startDir}"`,
      };
    }
    dir = parent;
  }
}

export function savepointPath(root: string, ...segments: string[]): string {
  return join(root, SAVEPOINT_DIR, ...segments);
}
