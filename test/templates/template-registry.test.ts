import { afterEach, beforeEach, describe, expect, it } from "vitest";
import { mkdtempSync, rmSync, mkdirSync, writeFileSync } from "node:fs";
import { join } from "node:path";
import { tmpdir } from "node:os";
import {
  PROJECT_TEMPLATES,
  RELEASE_TEMPLATES,
  PROMPT_TEMPLATES,
  templatePath,
  loadTemplate,
  renderTemplate,
} from "../../src/templates/index.js";

let tmp: string;

beforeEach(() => {
  tmp = mkdtempSync(join(tmpdir(), "savepoint-template-test-"));
});

afterEach(() => {
  rmSync(tmp, { recursive: true });
});

describe("templatePath", () => {
  it("resolves project templates under project/ subdirectory", () => {
    const p = templatePath("AGENTS.md", tmp);
    expect(p).toBe(join(tmp, "project", "AGENTS.md"));
  });

  it("resolves release templates under release/ subdirectory", () => {
    const p = templatePath("v1/PRD.md", tmp);
    expect(p).toBe(join(tmp, "release", "v1", "PRD.md"));
  });

  it("resolves prompt templates under prompts/ subdirectory", () => {
    const p = templatePath("prd.prompt.md", tmp);
    expect(p).toBe(join(tmp, "prompts", "prd.prompt.md"));
  });
});

describe("loadTemplate", () => {
  it("loads an existing template", async () => {
    mkdirSync(join(tmp, "project"), { recursive: true });
    writeFileSync(join(tmp, "project", "AGENTS.md"), "# Hello");
    const r = await loadTemplate("AGENTS.md", { root: tmp });
    expect(r.ok).toBe(true);
    if (r.ok) expect(r.value).toBe("# Hello");
  });

  it("returns not_found for a missing template", async () => {
    const r = await loadTemplate("AGENTS.md", { root: tmp });
    expect(r.ok).toBe(false);
    if (!r.ok) {
      expect(r.error.code).toBe("not_found");
      expect(r.error.path).toContain("AGENTS.md");
    }
  });
});

describe("renderTemplate", () => {
  it("interpolates projectName", () => {
    const r = renderTemplate("Hello {{PROJECT_NAME}}", {
      projectName: "World",
    });
    expect(r).toBe("Hello World");
  });

  it("interpolates releaseNumber", () => {
    const r = renderTemplate("v{{RELEASE_NUMBER}}", { releaseNumber: "2" });
    expect(r).toBe("v2");
  });

  it("interpolates releaseName", () => {
    const r = renderTemplate("{{RELEASE_NAME}}", { releaseName: "Alpha" });
    expect(r).toBe("Alpha");
  });

  it("interpolates multiple variables", () => {
    const r = renderTemplate("{{PROJECT_NAME}} v{{RELEASE_NUMBER}}", {
      projectName: "Savepoint",
      releaseNumber: "1",
    });
    expect(r).toBe("Savepoint v1");
  });

  it("ignores missing variables", () => {
    const r = renderTemplate("{{PROJECT_NAME}}", {});
    expect(r).toBe("{{PROJECT_NAME}}");
  });

  it("handles spaced braces like { { PROJECT_NAME } }", () => {
    const r = renderTemplate("{ { PROJECT_NAME } }", { projectName: "X" });
    expect(r).toBe("X");
  });
});

describe("template manifest integrity", () => {
  it("all project templates exist", async () => {
    for (const name of PROJECT_TEMPLATES) {
      const r = await loadTemplate(name);
      expect(r.ok).toBe(true);
    }
  });

  it("all release templates exist", async () => {
    for (const name of RELEASE_TEMPLATES) {
      const r = await loadTemplate(name);
      expect(r.ok).toBe(true);
    }
  });

  it("all prompt templates exist", async () => {
    for (const name of PROMPT_TEMPLATES) {
      const r = await loadTemplate(name);
      expect(r.ok).toBe(true);
    }
  });
});
