import { afterEach, beforeEach, describe, expect, it } from "vitest";
import { mkdtempSync, rmSync, mkdirSync, writeFileSync } from "node:fs";
import { join } from "node:path";
import { tmpdir } from "node:os";
import { readReleasePrd } from "../../src/readers/release.js";

let tmp: string;

beforeEach(() => {
  tmp = mkdtempSync(join(tmpdir(), "savepoint-release-test-"));
  mkdirSync(join(tmp, ".savepoint", "releases", "v1"), { recursive: true });
});

afterEach(() => {
  rmSync(tmp, { recursive: true });
});

function prdPath(): string {
  return join(tmp, ".savepoint", "releases", "v1", "PRD.md");
}

describe("readReleasePrd — success", () => {
  it("parses a valid release PRD", async () => {
    writeFileSync(
      prdPath(),
      "---\nversion: 1\nname: MVP\nstatus: in_progress\n---\nBody",
    );
    const r = await readReleasePrd(tmp, "v1");
    expect(r.ok).toBe(true);
    if (r.ok) {
      expect(r.value.frontmatter.version).toBe(1);
      expect(r.value.frontmatter.name).toBe("MVP");
      expect(r.value.frontmatter.status).toBe("in_progress");
      expect(r.value.body).toContain("Body");
    }
  });
});

describe("readReleasePrd — not_found", () => {
  it("returns not_found for a missing PRD", async () => {
    const r = await readReleasePrd(tmp, "v99");
    expect(r.ok).toBe(false);
    if (!r.ok) {
      expect(r.error.code).toBe("not_found");
      expect(r.error.path).toContain("v99");
    }
  });
});

describe("readReleasePrd — missing_frontmatter", () => {
  it("returns missing_frontmatter when the file has no --- delimiters", async () => {
    writeFileSync(prdPath(), "No frontmatter here.");
    const r = await readReleasePrd(tmp, "v1");
    expect(r.ok).toBe(false);
    if (!r.ok) expect(r.error.code).toBe("missing_frontmatter");
  });
});

describe("readReleasePrd — schema_rejection", () => {
  it("returns schema_rejection when version is missing", async () => {
    writeFileSync(prdPath(), "---\nname: MVP\nstatus: in_progress\n---\nBody");
    const r = await readReleasePrd(tmp, "v1");
    expect(r.ok).toBe(false);
    if (!r.ok) {
      expect(r.error.code).toBe("schema_rejection");
      expect("reason" in r.error && r.error.reason).toMatch(/version/);
    }
  });
});
