import { describe, expect, it } from "vitest";
import { validateEpicFrontmatter } from "../../src/domain/epic.js";

const VALID = { type: "epic-design", status: "active" };

describe("validateEpicFrontmatter — valid", () => {
  it("accepts a complete valid frontmatter object", () => {
    const r = validateEpicFrontmatter(VALID);
    expect(r.ok).toBe(true);
    if (r.ok) {
      expect(r.value.type).toBe("epic-design");
      expect(r.value.status).toBe("active");
    }
  });
});

describe("validateEpicFrontmatter — non-object input", () => {
  it("rejects null", () => {
    expect(validateEpicFrontmatter(null).ok).toBe(false);
  });
  it("rejects a string", () => {
    expect(validateEpicFrontmatter("string").ok).toBe(false);
  });
  it("rejects an array", () => {
    expect(validateEpicFrontmatter([VALID]).ok).toBe(false);
  });
});

describe("validateEpicFrontmatter — type field", () => {
  it("rejects missing type", () => {
    const r = validateEpicFrontmatter({ status: "active" });
    expect(r.ok).toBe(false);
    if (!r.ok) expect(r.reason).toMatch(/type/);
  });
  it("rejects whitespace-only type", () => {
    const r = validateEpicFrontmatter({ ...VALID, type: "   " });
    expect(r.ok).toBe(false);
    if (!r.ok) expect(r.reason).toMatch(/type/);
  });
});

describe("validateEpicFrontmatter — status field", () => {
  it("rejects missing status", () => {
    const r = validateEpicFrontmatter({ type: "epic-design" });
    expect(r.ok).toBe(false);
    if (!r.ok) expect(r.reason).toMatch(/status/);
  });
  it("rejects whitespace-only status", () => {
    const r = validateEpicFrontmatter({ ...VALID, status: "  " });
    expect(r.ok).toBe(false);
    if (!r.ok) expect(r.reason).toMatch(/status/);
  });
});
