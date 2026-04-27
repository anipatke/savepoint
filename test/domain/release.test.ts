import { describe, expect, it } from "vitest";
import { validateReleaseFrontmatter } from "../../src/domain/release.js";

const VALID = { version: 1, name: "MVP", status: "in_progress" };

describe("validateReleaseFrontmatter — valid", () => {
  it("accepts a complete valid frontmatter object", () => {
    const r = validateReleaseFrontmatter(VALID);
    expect(r.ok).toBe(true);
    if (r.ok) {
      expect(r.value.version).toBe(1);
      expect(r.value.name).toBe("MVP");
      expect(r.value.status).toBe("in_progress");
    }
  });
});

describe("validateReleaseFrontmatter — non-object input", () => {
  it("rejects null", () => {
    expect(validateReleaseFrontmatter(null).ok).toBe(false);
  });
  it("rejects a string", () => {
    expect(validateReleaseFrontmatter("string").ok).toBe(false);
  });
  it("rejects an array", () => {
    expect(validateReleaseFrontmatter([VALID]).ok).toBe(false);
  });
});

describe("validateReleaseFrontmatter — version field", () => {
  it("rejects missing version", () => {
    const r = validateReleaseFrontmatter({ name: "MVP", status: "planned" });
    expect(r.ok).toBe(false);
    if (!r.ok) expect(r.reason).toMatch(/version/);
  });
  it("rejects string version", () => {
    const r = validateReleaseFrontmatter({ ...VALID, version: "1" });
    expect(r.ok).toBe(false);
    if (!r.ok) expect(r.reason).toMatch(/version/);
  });
});

describe("validateReleaseFrontmatter — name field", () => {
  it("rejects missing name", () => {
    const r = validateReleaseFrontmatter({ version: 1, status: "planned" });
    expect(r.ok).toBe(false);
    if (!r.ok) expect(r.reason).toMatch(/name/);
  });
  it("rejects whitespace-only name", () => {
    const r = validateReleaseFrontmatter({ ...VALID, name: "   " });
    expect(r.ok).toBe(false);
    if (!r.ok) expect(r.reason).toMatch(/name/);
  });
});

describe("validateReleaseFrontmatter — status field", () => {
  it("rejects missing status", () => {
    const r = validateReleaseFrontmatter({ version: 1, name: "MVP" });
    expect(r.ok).toBe(false);
    if (!r.ok) expect(r.reason).toMatch(/status/);
  });
  it("rejects whitespace-only status", () => {
    const r = validateReleaseFrontmatter({ ...VALID, status: "  " });
    expect(r.ok).toBe(false);
    if (!r.ok) expect(r.reason).toMatch(/status/);
  });
});
