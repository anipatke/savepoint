import { afterEach, beforeEach, describe, expect, it } from "vitest";
import { mkdtempSync, rmSync, mkdirSync, writeFileSync } from "node:fs";
import { join } from "node:path";
import { tmpdir } from "node:os";
import { readEpicDesign } from "../../src/readers/epic.js";

let tmp: string;

beforeEach(() => {
  tmp = mkdtempSync(join(tmpdir(), "savepoint-epic-test-"));
  mkdirSync(
    join(tmp, ".savepoint", "releases", "v1", "epics", "E02-data-model"),
    { recursive: true },
  );
});

afterEach(() => {
  rmSync(tmp, { recursive: true });
});

function designPath(): string {
  return join(
    tmp,
    ".savepoint",
    "releases",
    "v1",
    "epics",
    "E02-data-model",
    "Design.md",
  );
}

describe("readEpicDesign — success", () => {
  it("parses a valid epic Design.md", async () => {
    writeFileSync(
      designPath(),
      "---\ntype: epic-design\nstatus: active\n---\n# Design",
    );
    const r = await readEpicDesign(tmp, "v1", "E02-data-model");
    expect(r.ok).toBe(true);
    if (r.ok) {
      expect(r.value.frontmatter.type).toBe("epic-design");
      expect(r.value.frontmatter.status).toBe("active");
    }
  });
});

describe("readEpicDesign — not_found", () => {
  it("returns not_found for a missing epic", async () => {
    const r = await readEpicDesign(tmp, "v1", "E99-missing");
    expect(r.ok).toBe(false);
    if (!r.ok) {
      expect(r.error.code).toBe("not_found");
      expect(r.error.path).toContain("E99-missing");
    }
  });
});

describe("readEpicDesign — missing_frontmatter", () => {
  it("returns missing_frontmatter when file has no --- delimiters", async () => {
    writeFileSync(designPath(), "Just markdown, no frontmatter.");
    const r = await readEpicDesign(tmp, "v1", "E02-data-model");
    expect(r.ok).toBe(false);
    if (!r.ok) expect(r.error.code).toBe("missing_frontmatter");
  });
});

describe("readEpicDesign — schema_rejection", () => {
  it("returns schema_rejection when type is missing", async () => {
    writeFileSync(designPath(), "---\nstatus: active\n---\n# Design");
    const r = await readEpicDesign(tmp, "v1", "E02-data-model");
    expect(r.ok).toBe(false);
    if (!r.ok) {
      expect(r.error.code).toBe("schema_rejection");
      expect("reason" in r.error && r.error.reason).toMatch(/type/);
    }
  });
});
