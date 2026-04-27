import { afterEach, beforeEach, describe, expect, it } from "vitest";
import { mkdtempSync, rmSync, writeFileSync } from "node:fs";
import { join } from "node:path";
import { tmpdir } from "node:os";
import { readMarkdownFile } from "../../src/fs/markdown.js";

let tmp: string;

beforeEach(() => {
  tmp = mkdtempSync(join(tmpdir(), "savepoint-md-test-"));
});

afterEach(() => {
  rmSync(tmp, { recursive: true });
});

function acceptAll(raw: unknown): { ok: true; value: unknown } {
  return { ok: true, value: raw };
}

function requireTitle(
  raw: unknown,
): { ok: true; value: { title: string } } | { ok: false; reason: string } {
  if (
    typeof raw !== "object" ||
    raw === null ||
    !("title" in raw) ||
    typeof (raw as Record<string, unknown>).title !== "string"
  ) {
    return {
      ok: false,
      reason: "frontmatter must have a string 'title' field",
    };
  }
  return { ok: true, value: { title: (raw as { title: string }).title } };
}

describe("readMarkdownFile — success", () => {
  it("parses frontmatter and body", async () => {
    const file = join(tmp, "doc.md");
    writeFileSync(file, "---\ntitle: Hello\n---\nBody text here");
    const r = await readMarkdownFile(file, acceptAll);
    expect(r.ok).toBe(true);
    if (r.ok) {
      expect((r.value.frontmatter as { title: string }).title).toBe("Hello");
      expect(r.value.body).toBe("Body text here");
    }
  });

  it("preserves multiline body", async () => {
    const file = join(tmp, "multi.md");
    writeFileSync(file, "---\nkey: val\n---\nLine 1\nLine 2\n");
    const r = await readMarkdownFile(file, acceptAll);
    expect(r.ok).toBe(true);
    if (r.ok) expect(r.value.body).toContain("Line 2");
  });

  it("passes validated frontmatter through schema validator", async () => {
    const file = join(tmp, "valid.md");
    writeFileSync(file, "---\ntitle: My Task\n---\nContent");
    const r = await readMarkdownFile(file, requireTitle);
    expect(r.ok).toBe(true);
    if (r.ok) expect(r.value.frontmatter.title).toBe("My Task");
  });
});

describe("readMarkdownFile — not_found", () => {
  it("returns not_found for a missing file", async () => {
    const r = await readMarkdownFile(join(tmp, "ghost.md"), acceptAll);
    expect(r.ok).toBe(false);
    if (!r.ok) {
      expect(r.error.code).toBe("not_found");
      expect(r.error.path).toContain("ghost.md");
    }
  });
});

describe("readMarkdownFile — missing_frontmatter", () => {
  it("returns missing_frontmatter when file has no --- delimiters", async () => {
    const file = join(tmp, "no-fm.md");
    writeFileSync(file, "Just plain markdown, no frontmatter.");
    const r = await readMarkdownFile(file, acceptAll);
    expect(r.ok).toBe(false);
    if (!r.ok) {
      expect(r.error.code).toBe("missing_frontmatter");
      expect(r.error.path).toContain("no-fm.md");
    }
  });
});

describe("readMarkdownFile — malformed_frontmatter", () => {
  it("returns malformed_frontmatter for invalid YAML", async () => {
    const file = join(tmp, "bad-yaml.md");
    writeFileSync(file, "---\nkey: [\nbad yaml\n---\nBody");
    const r = await readMarkdownFile(file, acceptAll);
    expect(r.ok).toBe(false);
    if (!r.ok) {
      expect(r.error.code).toBe("malformed_frontmatter");
      expect(r.error.path).toContain("bad-yaml.md");
      expect("reason" in r.error).toBe(true);
    }
  });
});

describe("readMarkdownFile — schema_rejection", () => {
  it("returns schema_rejection when validator rejects frontmatter", async () => {
    const file = join(tmp, "schema-fail.md");
    writeFileSync(file, "---\nstatus: ok\n---\nBody");
    const r = await readMarkdownFile(file, requireTitle);
    expect(r.ok).toBe(false);
    if (!r.ok) {
      expect(r.error.code).toBe("schema_rejection");
      expect(r.error.path).toContain("schema-fail.md");
      expect("reason" in r.error).toBe(true);
    }
  });
});
