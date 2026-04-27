import { describe, expect, it } from "vitest";
import {
  formatEpicId,
  formatTaskId,
  parseEpicId,
  parseReleaseId,
  parseTaskId,
} from "../../src/domain/ids.js";

describe("parseReleaseId", () => {
  it("accepts v1", () => {
    const r = parseReleaseId("v1");
    expect(r.ok).toBe(true);
    if (r.ok) {
      expect(r.value.raw).toBe("v1");
      expect(r.value.tag).toBe("release");
    }
  });

  it("accepts multi-digit versions", () => {
    const r = parseReleaseId("v12");
    expect(r.ok).toBe(true);
  });

  it("rejects bare number", () => {
    const r = parseReleaseId("1");
    expect(r.ok).toBe(false);
  });

  it("rejects empty string", () => {
    const r = parseReleaseId("");
    expect(r.ok).toBe(false);
  });

  it("rejects non-numeric suffix", () => {
    const r = parseReleaseId("vbeta");
    expect(r.ok).toBe(false);
  });
});

describe("parseEpicId", () => {
  it("accepts E02-data-model", () => {
    const r = parseEpicId("E02-data-model");
    expect(r.ok).toBe(true);
    if (r.ok) {
      expect(r.value.number).toBe(2);
      expect(r.value.slug).toBe("data-model");
      expect(r.value.tag).toBe("epic");
    }
  });

  it("accepts single-digit epic with leading zero", () => {
    const r = parseEpicId("E01-scaffolding");
    expect(r.ok).toBe(true);
    if (r.ok) expect(r.value.number).toBe(1);
  });

  it("rejects lowercase e", () => {
    const r = parseEpicId("e02-data-model");
    expect(r.ok).toBe(false);
  });

  it("rejects missing slug", () => {
    const r = parseEpicId("E02");
    expect(r.ok).toBe(false);
  });

  it("rejects slug starting with digit", () => {
    const r = parseEpicId("E02-1bad");
    expect(r.ok).toBe(false);
  });

  it("rejects empty string", () => {
    const r = parseEpicId("");
    expect(r.ok).toBe(false);
  });
});

describe("parseTaskId", () => {
  it("accepts E02-data-model/T001-domain-ids-status", () => {
    const r = parseTaskId("E02-data-model/T001-domain-ids-status");
    expect(r.ok).toBe(true);
    if (r.ok) {
      expect(r.value.tag).toBe("task");
      expect(r.value.number).toBe(1);
      expect(r.value.slug).toBe("domain-ids-status");
      expect(r.value.epic.number).toBe(2);
      expect(r.value.epic.slug).toBe("data-model");
    }
  });

  it("rejects missing slash separator", () => {
    const r = parseTaskId("E02-data-modelT001-slug");
    expect(r.ok).toBe(false);
  });

  it("rejects invalid epic prefix", () => {
    const r = parseTaskId("bad-epic/T001-slug");
    expect(r.ok).toBe(false);
  });

  it("rejects lowercase t in task segment", () => {
    const r = parseTaskId("E02-data-model/t001-slug");
    expect(r.ok).toBe(false);
  });

  it("rejects task number with fewer than 3 digits", () => {
    const r = parseTaskId("E02-data-model/T01-slug");
    expect(r.ok).toBe(false);
  });

  it("rejects task segment missing slug", () => {
    const r = parseTaskId("E02-data-model/T001");
    expect(r.ok).toBe(false);
  });

  it("rejects empty string", () => {
    const r = parseTaskId("");
    expect(r.ok).toBe(false);
  });
});

describe("formatEpicId", () => {
  it("round-trips the raw string", () => {
    const r = parseEpicId("E03-readers");
    expect(r.ok).toBe(true);
    if (r.ok) expect(formatEpicId(r.value)).toBe("E03-readers");
  });
});

describe("formatTaskId", () => {
  it("round-trips the raw string", () => {
    const r = parseTaskId("E02-data-model/T001-domain-ids-status");
    expect(r.ok).toBe(true);
    if (r.ok)
      expect(formatTaskId(r.value)).toBe(
        "E02-data-model/T001-domain-ids-status",
      );
  });
});
