import { afterEach, beforeEach, describe, expect, it } from "vitest";
import { mkdtempSync, rmSync, mkdirSync, writeFileSync } from "node:fs";
import { join } from "node:path";
import { tmpdir } from "node:os";
import { readEpicTaskSet } from "../../src/readers/tasks.js";

let tmp: string;
let tasksDir: string;

beforeEach(() => {
  tmp = mkdtempSync(join(tmpdir(), "savepoint-tasks-test-"));
  tasksDir = join(
    tmp,
    ".savepoint",
    "releases",
    "v1",
    "epics",
    "E02-data-model",
    "tasks",
  );
  mkdirSync(tasksDir, { recursive: true });
});

afterEach(() => {
  rmSync(tmp, { recursive: true });
});

function writeTask(filename: string, content: string): void {
  writeFileSync(join(tasksDir, filename), content);
}

const T1 = `---
id: E02-data-model/T001-alpha
status: planned
objective: First task
depends_on: []
---
# T001`;

const T2 = `---
id: E02-data-model/T002-bravo
status: planned
objective: Second task
depends_on:
  - E02-data-model/T001-alpha
---
# T002`;

describe("readEpicTaskSet — success", () => {
  it("returns ok with all tasks when files parse and graph is valid", async () => {
    writeTask("T001-alpha.md", T1);
    writeTask("T002-bravo.md", T2);
    const r = await readEpicTaskSet(tmp, "v1", "E02-data-model");
    expect(r.ok).toBe(true);
    if (r.ok) {
      expect(r.tasks).toHaveLength(2);
      const ids = r.tasks.map((t) => t.frontmatter.id.raw);
      expect(ids).toContain("E02-data-model/T001-alpha");
      expect(ids).toContain("E02-data-model/T002-bravo");
    }
  });

  it("returns ok with empty tasks when directory has no .md files", async () => {
    const r = await readEpicTaskSet(tmp, "v1", "E02-data-model");
    expect(r.ok).toBe(true);
    if (r.ok) {
      expect(r.tasks).toHaveLength(0);
    }
  });

  it("ignores non-.md files in the tasks directory", async () => {
    writeTask("T001-alpha.md", T1);
    writeFileSync(join(tasksDir, "notes.txt"), "ignored");
    const r = await readEpicTaskSet(tmp, "v1", "E02-data-model");
    expect(r.ok).toBe(true);
    if (r.ok) {
      expect(r.tasks).toHaveLength(1);
    }
  });
});

describe("readEpicTaskSet — fs_errors", () => {
  it("returns fs_errors when the tasks directory does not exist", async () => {
    rmSync(tasksDir, { recursive: true });
    const r = await readEpicTaskSet(tmp, "v1", "E02-data-model");
    expect(r.ok).toBe(false);
    if (!r.ok) {
      expect(r.kind).toBe("fs_errors");
      if (r.kind === "fs_errors") {
        expect(r.errors[0].code).toBe("not_found");
      }
    }
  });

  it("returns fs_errors when a task file has no frontmatter", async () => {
    writeTask("T001-alpha.md", "Just markdown, no frontmatter.");
    const r = await readEpicTaskSet(tmp, "v1", "E02-data-model");
    expect(r.ok).toBe(false);
    if (!r.ok) {
      expect(r.kind).toBe("fs_errors");
      if (r.kind === "fs_errors") {
        expect(r.errors).toHaveLength(1);
        expect(r.errors[0].code).toBe("missing_frontmatter");
      }
    }
  });

  it("returns fs_errors when a task file fails schema validation", async () => {
    writeTask("T001-alpha.md", "---\nstatus: planned\n---\n# no id field");
    const r = await readEpicTaskSet(tmp, "v1", "E02-data-model");
    expect(r.ok).toBe(false);
    if (!r.ok) {
      expect(r.kind).toBe("fs_errors");
      if (r.kind === "fs_errors") {
        expect(r.errors[0].code).toBe("schema_rejection");
      }
    }
  });

  it("collects errors from all bad files before returning", async () => {
    writeTask("T001-alpha.md", "no frontmatter here");
    writeTask("T002-bravo.md", "no frontmatter here either");
    const r = await readEpicTaskSet(tmp, "v1", "E02-data-model");
    expect(r.ok).toBe(false);
    if (!r.ok) {
      expect(r.kind).toBe("fs_errors");
      if (r.kind === "fs_errors") {
        expect(r.errors).toHaveLength(2);
      }
    }
  });

  it("does not run graph validation when any file fails to parse", async () => {
    writeTask("T001-alpha.md", T1);
    writeTask("T002-bravo.md", "no frontmatter");
    const r = await readEpicTaskSet(tmp, "v1", "E02-data-model");
    expect(r.ok).toBe(false);
    if (!r.ok) {
      expect(r.kind).toBe("fs_errors");
    }
  });
});

describe("readEpicTaskSet — graph_errors", () => {
  it("reports duplicate IDs and includes them in errors", async () => {
    writeTask("T001-alpha.md", T1);
    writeTask("T001-alpha-dup.md", T1);
    const r = await readEpicTaskSet(tmp, "v1", "E02-data-model");
    expect(r.ok).toBe(false);
    if (!r.ok) {
      expect(r.kind).toBe("graph_errors");
      if (r.kind === "graph_errors") {
        expect(r.errors.duplicateIds).toHaveLength(1);
        expect(r.errors.duplicateIds[0].raw).toBe("E02-data-model/T001-alpha");
      }
    }
  });

  it("reports a missing dependency when a depends_on ID is not in the set", async () => {
    writeTask("T002-bravo.md", T2);
    const r = await readEpicTaskSet(tmp, "v1", "E02-data-model");
    expect(r.ok).toBe(false);
    if (!r.ok) {
      expect(r.kind).toBe("graph_errors");
      if (r.kind === "graph_errors") {
        expect(r.errors.missingDependencies).toHaveLength(1);
        expect(r.errors.missingDependencies[0].missing.raw).toBe(
          "E02-data-model/T001-alpha",
        );
      }
    }
  });

  it("detects a dependency cycle", async () => {
    const cycleT1 = `---
id: E02-data-model/T001-alpha
status: planned
objective: Cycle node A
depends_on:
  - E02-data-model/T002-bravo
---
# T001`;
    const cycleT2 = `---
id: E02-data-model/T002-bravo
status: planned
objective: Cycle node B
depends_on:
  - E02-data-model/T001-alpha
---
# T002`;
    writeTask("T001-alpha.md", cycleT1);
    writeTask("T002-bravo.md", cycleT2);
    const r = await readEpicTaskSet(tmp, "v1", "E02-data-model");
    expect(r.ok).toBe(false);
    if (!r.ok) {
      expect(r.kind).toBe("graph_errors");
      if (r.kind === "graph_errors") {
        expect(r.errors.cycle).not.toBeNull();
      }
    }
  });

  it("includes the parsed tasks in graph_errors so callers can inspect them", async () => {
    writeTask("T002-bravo.md", T2);
    const r = await readEpicTaskSet(tmp, "v1", "E02-data-model");
    expect(r.ok).toBe(false);
    if (!r.ok && r.kind === "graph_errors") {
      expect(r.tasks).toHaveLength(1);
    }
  });
});
