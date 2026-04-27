import { describe, expect, it } from "vitest";
import { PROMPT_TEMPLATES, loadTemplate } from "../../src/templates/index.js";

describe("prompt template integrity", () => {
  it("all required prompt templates exist", async () => {
    const required = [
      "prd.prompt.md",
      "design.prompt.md",
      "epic-design.prompt.md",
      "task-breakdown.prompt.md",
      "task-planning.prompt.md",
      "task-building.prompt.md",
      "audit-reconciliation.prompt.md",
    ];
    for (const name of required) {
      const r = await loadTemplate(name as (typeof PROMPT_TEMPLATES)[number]);
      expect(r.ok).toBe(true);
    }
  });

  it("every prompt contains an AGENT instruction marker", async () => {
    for (const name of PROMPT_TEMPLATES) {
      const r = await loadTemplate(name);
      expect(r.ok).toBe(true);
      if (!r.ok) continue;
      expect(r.value).toContain("<!-- AGENT:");
    }
  });

  it("prd.prompt.md contains workflow markers", async () => {
    const r = await loadTemplate("prd.prompt.md");
    expect(r.ok).toBe(true);
    if (!r.ok) return;
    expect(r.value).toContain("## Rules");
    expect(r.value).toContain("type: project-prd");
  });

  it("design.prompt.md contains workflow markers", async () => {
    const r = await loadTemplate("design.prompt.md");
    expect(r.ok).toBe(true);
    if (!r.ok) return;
    expect(r.value).toContain("## Rules");
    expect(r.value).toContain("type: project-design");
  });

  it("epic-design.prompt.md contains workflow markers", async () => {
    const r = await loadTemplate("epic-design.prompt.md");
    expect(r.ok).toBe(true);
    if (!r.ok) return;
    expect(r.value).toContain("## Rules");
    expect(r.value).toContain("type: epic-design");
  });

  it("task-breakdown.prompt.md contains depends_on guidance", async () => {
    const r = await loadTemplate("task-breakdown.prompt.md");
    expect(r.ok).toBe(true);
    if (!r.ok) return;
    expect(r.value).toContain("depends_on");
    expect(r.value).toContain("## Rules");
  });

  it("task-planning.prompt.md contains planned status guidance", async () => {
    const r = await loadTemplate("task-planning.prompt.md");
    expect(r.ok).toBe(true);
    if (!r.ok) return;
    expect(r.value).toContain("status: planned");
    expect(r.value).toContain("## Rules");
  });

  it("task-building.prompt.md contains workflow markers", async () => {
    const r = await loadTemplate("task-building.prompt.md");
    expect(r.ok).toBe(true);
    if (!r.ok) return;
    expect(r.value).toContain("## Rules");
    expect(r.value).toContain("status: review");
  });

  it("audit-reconciliation.prompt.md references single proposals bundle and delta-only edits", async () => {
    const r = await loadTemplate("audit-reconciliation.prompt.md");
    expect(r.ok).toBe(true);
    if (!r.ok) return;
    expect(r.value).toContain("proposals.md");
    expect(r.value).toContain("delta-only");
    expect(r.value).toContain("## Rules");
    expect(r.value).toContain("One proposals file only");
  });
});
