import { describe, expect, it } from "vitest";
import {
  applyConfigDefaults,
  CONFIG_DEFAULTS,
} from "../../src/domain/config.js";

describe("applyConfigDefaults — empty input", () => {
  it("returns all defaults when given an empty object", () => {
    const cfg = applyConfigDefaults({});
    expect(cfg).toEqual(CONFIG_DEFAULTS);
  });

  it("returns all defaults when given null", () => {
    const cfg = applyConfigDefaults(null);
    expect(cfg).toEqual(CONFIG_DEFAULTS);
  });

  it("returns all defaults when given a non-object", () => {
    const cfg = applyConfigDefaults("bad");
    expect(cfg).toEqual(CONFIG_DEFAULTS);
  });
});

describe("applyConfigDefaults — verify_strict", () => {
  it("uses provided true value", () => {
    const cfg = applyConfigDefaults({ verify_strict: true });
    expect(cfg.verify_strict).toBe(true);
  });

  it("uses provided false value", () => {
    const cfg = applyConfigDefaults({ verify_strict: false });
    expect(cfg.verify_strict).toBe(false);
  });

  it("falls back to default when non-boolean", () => {
    const cfg = applyConfigDefaults({ verify_strict: "yes" });
    expect(cfg.verify_strict).toBe(CONFIG_DEFAULTS.verify_strict);
  });
});

describe("applyConfigDefaults — quality_gates", () => {
  it("overrides lint with a string value", () => {
    const cfg = applyConfigDefaults({ quality_gates: { lint: "eslint ." } });
    expect(cfg.quality_gates.lint).toBe("eslint .");
  });

  it("accepts explicit null for lint", () => {
    const cfg = applyConfigDefaults({ quality_gates: { lint: null } });
    expect(cfg.quality_gates.lint).toBeNull();
  });

  it("falls back to default for unset lint", () => {
    const cfg = applyConfigDefaults({ quality_gates: {} });
    expect(cfg.quality_gates.lint).toBe(CONFIG_DEFAULTS.quality_gates.lint);
  });

  it("overrides block_on_failure", () => {
    const cfg = applyConfigDefaults({
      quality_gates: { block_on_failure: false },
    });
    expect(cfg.quality_gates.block_on_failure).toBe(false);
  });
});

describe("applyConfigDefaults — audit", () => {
  it("overrides divergence_threshold", () => {
    const cfg = applyConfigDefaults({ audit: { divergence_threshold: 0.8 } });
    expect(cfg.audit.divergence_threshold).toBe(0.8);
  });

  it("falls back to default for non-numeric threshold", () => {
    const cfg = applyConfigDefaults({
      audit: { divergence_threshold: "high" },
    });
    expect(cfg.audit.divergence_threshold).toBe(
      CONFIG_DEFAULTS.audit.divergence_threshold,
    );
  });
});

describe("applyConfigDefaults — theme", () => {
  it("overrides bg color", () => {
    const cfg = applyConfigDefaults({ theme: { bg: "#000000" } });
    expect(cfg.theme.bg).toBe("#000000");
  });

  it("falls back to default for non-string bg", () => {
    const cfg = applyConfigDefaults({ theme: { bg: 42 } });
    expect(cfg.theme.bg).toBe(CONFIG_DEFAULTS.theme.bg);
  });

  it("merges accents over defaults", () => {
    const cfg = applyConfigDefaults({
      theme: { accents: { done: "#00FF00", custom: "#AABBCC" } },
    });
    expect(cfg.theme.accents["done"]).toBe("#00FF00");
    expect(cfg.theme.accents["custom"]).toBe("#AABBCC");
    expect(cfg.theme.accents["backlog"]).toBe(
      CONFIG_DEFAULTS.theme.accents["backlog"],
    );
  });

  it("ignores non-string accent values", () => {
    const cfg = applyConfigDefaults({
      theme: { accents: { done: 42, custom: "#AABBCC" } },
    });
    expect(cfg.theme.accents["done"]).toBe(CONFIG_DEFAULTS.theme.accents.done);
    expect(cfg.theme.accents["custom"]).toBe("#AABBCC");
  });

  it("accepts a valid borders value", () => {
    const cfg = applyConfigDefaults({ theme: { borders: "sharp" } });
    expect(cfg.theme.borders).toBe("sharp");
  });

  it("falls back to default borders for an invalid value", () => {
    const cfg = applyConfigDefaults({ theme: { borders: "dotted" } });
    expect(cfg.theme.borders).toBe(CONFIG_DEFAULTS.theme.borders);
  });
});
