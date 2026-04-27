import { describe, expect, it } from "vitest";
import { validateTaskFrontmatter } from "../../src/domain/task.js";

const VALID_FM = {
  id: "E02-data-model/T003-task-documents",
  status: "planned",
  objective: "A valid objective",
  depends_on: ["E02-data-model/T001-domain-ids-status"],
};

describe("validateTaskFrontmatter", () => {
  it("accepts valid frontmatter with one dependency", () => {
    const r = validateTaskFrontmatter(VALID_FM);
    expect(r.ok).toBe(true);
    if (r.ok) {
      expect(r.value.id.raw).toBe("E02-data-model/T003-task-documents");
      expect(r.value.status).toBe("planned");
      expect(r.value.objective).toBe("A valid objective");
      expect(r.value.depends_on).toHaveLength(1);
      expect(r.value.depends_on[0].raw).toBe(
        "E02-data-model/T001-domain-ids-status",
      );
    }
  });

  it("accepts valid frontmatter with empty depends_on", () => {
    const r = validateTaskFrontmatter({ ...VALID_FM, depends_on: [] });
    expect(r.ok).toBe(true);
    if (r.ok) expect(r.value.depends_on).toHaveLength(0);
  });

  it("accepts valid frontmatter with multiple dependencies", () => {
    const r = validateTaskFrontmatter({
      ...VALID_FM,
      depends_on: [
        "E02-data-model/T001-domain-ids-status",
        "E02-data-model/T002-markdown-frontmatter-boundary",
      ],
    });
    expect(r.ok).toBe(true);
    if (r.ok) expect(r.value.depends_on).toHaveLength(2);
  });

  it("rejects null", () => {
    const r = validateTaskFrontmatter(null);
    expect(r.ok).toBe(false);
  });

  it("rejects a non-object string", () => {
    const r = validateTaskFrontmatter("not an object");
    expect(r.ok).toBe(false);
  });

  it("rejects an array", () => {
    const r = validateTaskFrontmatter([VALID_FM]);
    expect(r.ok).toBe(false);
  });

  it("rejects missing id", () => {
    const r = validateTaskFrontmatter({
      status: VALID_FM.status,
      objective: VALID_FM.objective,
      depends_on: VALID_FM.depends_on,
    });
    expect(r.ok).toBe(false);
    if (!r.ok) expect(r.reason).toMatch(/id/);
  });

  it("rejects a numeric id", () => {
    const r = validateTaskFrontmatter({ ...VALID_FM, id: 42 });
    expect(r.ok).toBe(false);
    if (!r.ok) expect(r.reason).toMatch(/id/);
  });

  it("rejects an invalid task id format", () => {
    const r = validateTaskFrontmatter({ ...VALID_FM, id: "not-a-task-id" });
    expect(r.ok).toBe(false);
  });

  it("rejects a task id missing the epic prefix", () => {
    const r = validateTaskFrontmatter({
      ...VALID_FM,
      id: "T003-task-documents",
    });
    expect(r.ok).toBe(false);
  });

  it("rejects missing status", () => {
    const r = validateTaskFrontmatter({
      id: VALID_FM.id,
      objective: VALID_FM.objective,
      depends_on: VALID_FM.depends_on,
    });
    expect(r.ok).toBe(false);
    if (!r.ok) expect(r.reason).toMatch(/status/);
  });

  it("rejects an unknown status value", () => {
    const r = validateTaskFrontmatter({ ...VALID_FM, status: "waiting" });
    expect(r.ok).toBe(false);
  });

  it("rejects missing objective", () => {
    const r = validateTaskFrontmatter({
      id: VALID_FM.id,
      status: VALID_FM.status,
      depends_on: VALID_FM.depends_on,
    });
    expect(r.ok).toBe(false);
    if (!r.ok) expect(r.reason).toMatch(/objective/);
  });

  it("rejects a whitespace-only objective", () => {
    const r = validateTaskFrontmatter({ ...VALID_FM, objective: "   " });
    expect(r.ok).toBe(false);
    if (!r.ok) expect(r.reason).toMatch(/objective/);
  });

  it("rejects a non-array depends_on", () => {
    const r = validateTaskFrontmatter({
      ...VALID_FM,
      depends_on: "E02-data-model/T001-domain-ids-status",
    });
    expect(r.ok).toBe(false);
    if (!r.ok) expect(r.reason).toMatch(/depends_on/);
  });

  it("rejects absent depends_on", () => {
    const r = validateTaskFrontmatter({
      id: VALID_FM.id,
      status: VALID_FM.status,
      objective: VALID_FM.objective,
    });
    expect(r.ok).toBe(false);
    if (!r.ok) expect(r.reason).toMatch(/depends_on/);
  });

  it("rejects a non-string entry inside depends_on", () => {
    const r = validateTaskFrontmatter({ ...VALID_FM, depends_on: [42] });
    expect(r.ok).toBe(false);
    if (!r.ok) expect(r.reason).toMatch(/depends_on\[0\]/);
  });

  it("rejects a malformed dependency id inside depends_on", () => {
    const r = validateTaskFrontmatter({
      ...VALID_FM,
      depends_on: ["not-a-valid-id"],
    });
    expect(r.ok).toBe(false);
    if (!r.ok) expect(r.reason).toMatch(/depends_on\[0\]/);
  });
});
