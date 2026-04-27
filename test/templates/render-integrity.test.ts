import { describe, expect, it } from "vitest";
import { renderTemplate } from "../../src/templates/index.js";

describe("render integrity", () => {
  it("interpolates projectName", () => {
    const r = renderTemplate("Hello {{PROJECT_NAME}}", {
      projectName: "Savepoint",
    });
    expect(r).toBe("Hello Savepoint");
  });

  it("interpolates releaseNumber", () => {
    const r = renderTemplate("v{{RELEASE_NUMBER}}", { releaseNumber: "1" });
    expect(r).toBe("v1");
  });

  it("interpolates releaseName", () => {
    const r = renderTemplate("{{RELEASE_NAME}}", { releaseName: "Alpha" });
    expect(r).toBe("Alpha");
  });

  it("interpolates multiple supported variables", () => {
    const r = renderTemplate(
      "{{PROJECT_NAME}} v{{RELEASE_NUMBER}} — {{RELEASE_NAME}}",
      { projectName: "X", releaseNumber: "2", releaseName: "Beta" },
    );
    expect(r).toBe("X v2 — Beta");
  });

  it("leaves unresolved variable placeholders as a clear failure signal", () => {
    const r = renderTemplate("{{PROJECT_NAME}}", {});
    expect(r).toBe("{{PROJECT_NAME}}");
  });

  it("leaves multiple unresolved placeholders intact", () => {
    const r = renderTemplate("{{PROJECT_NAME}} v{{RELEASE_NUMBER}}", {});
    expect(r).toBe("{{PROJECT_NAME}} v{{RELEASE_NUMBER}}");
  });

  it("handles spaced braces", () => {
    const r = renderTemplate("{ { PROJECT_NAME } }", { projectName: "Y" });
    expect(r).toBe("Y");
  });
});
