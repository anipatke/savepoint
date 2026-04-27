import { afterEach, beforeEach, describe, expect, it } from "vitest";
import { mkdtempSync, rmSync, mkdirSync, existsSync } from "node:fs";
import { join } from "node:path";
import { tmpdir } from "node:os";
import { findSavepointRoot, savepointPath } from "../../src/fs/project.js";

let tmp: string;

beforeEach(() => {
  tmp = mkdtempSync(join(tmpdir(), "savepoint-test-"));
});

afterEach(() => {
  rmSync(tmp, { recursive: true });
});

describe("findSavepointRoot", () => {
  it("finds root when .savepoint exists in the start directory", () => {
    mkdirSync(join(tmp, ".savepoint"));
    const r = findSavepointRoot(tmp);
    expect(r.ok).toBe(true);
    if (r.ok) expect(r.value).toBe(tmp);
  });

  it("finds root by walking up one level", () => {
    mkdirSync(join(tmp, ".savepoint"));
    const child = join(tmp, "subdir");
    mkdirSync(child);
    const r = findSavepointRoot(child);
    expect(r.ok).toBe(true);
    if (r.ok) expect(r.value).toBe(tmp);
  });

  it("returns error when no .savepoint exists anywhere", () => {
    const r = findSavepointRoot(tmp);
    expect(r.ok).toBe(false);
    if (!r.ok) expect(r.reason).toContain(".savepoint");
  });
});

describe("savepointPath", () => {
  it("joins root with .savepoint and segments", () => {
    const p = savepointPath(tmp, "releases", "v1");
    expect(p).toContain(".savepoint");
    expect(p).toContain("releases");
    expect(p.startsWith(tmp)).toBe(true);
  });

  it("works with no extra segments", () => {
    const p = savepointPath(tmp);
    expect(p).toContain(".savepoint");
    expect(existsSync(p)).toBe(false);
  });
});
