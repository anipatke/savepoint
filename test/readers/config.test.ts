import { afterEach, beforeEach, describe, expect, it } from "vitest";
import { mkdtempSync, rmSync, mkdirSync, writeFileSync } from "node:fs";
import { join } from "node:path";
import { tmpdir } from "node:os";
import { readConfig } from "../../src/readers/config.js";
import { CONFIG_DEFAULTS } from "../../src/domain/config.js";

let tmp: string;

beforeEach(() => {
  tmp = mkdtempSync(join(tmpdir(), "savepoint-config-test-"));
  mkdirSync(join(tmp, ".savepoint"), { recursive: true });
});

afterEach(() => {
  rmSync(tmp, { recursive: true });
});

function configPath(): string {
  return join(tmp, ".savepoint", "config.yml");
}

describe("readConfig — success", () => {
  it("returns defaults for an empty config file", async () => {
    writeFileSync(configPath(), "");
    const r = await readConfig(tmp);
    expect(r.ok).toBe(true);
    if (r.ok) expect(r.value).toEqual(CONFIG_DEFAULTS);
  });

  it("parses and overrides verify_strict", async () => {
    writeFileSync(configPath(), "verify_strict: true\n");
    const r = await readConfig(tmp);
    expect(r.ok).toBe(true);
    if (r.ok) expect(r.value.verify_strict).toBe(true);
  });

  it("applies defaults for unset quality_gates fields", async () => {
    writeFileSync(configPath(), "quality_gates:\n  lint: eslint .\n");
    const r = await readConfig(tmp);
    expect(r.ok).toBe(true);
    if (r.ok) {
      expect(r.value.quality_gates.lint).toBe("eslint .");
      expect(r.value.quality_gates.typecheck).toBe(
        CONFIG_DEFAULTS.quality_gates.typecheck,
      );
    }
  });

  it("reads a full config without applying defaults", async () => {
    writeFileSync(
      configPath(),
      [
        "verify_strict: true",
        "quality_gates:",
        "  lint: eslint .",
        "  typecheck: tsc",
        "  test: vitest",
        "  block_on_failure: false",
        "audit:",
        "  divergence_threshold: 0.3",
        "theme:",
        "  bg: '#000'",
        "  surface: '#111'",
        "  surface_2: '#222'",
        "  border: '#333'",
        "  text: '#fff'",
        "  borders: sharp",
        "  accents:",
        "    done: '#0f0'",
      ].join("\n"),
    );
    const r = await readConfig(tmp);
    expect(r.ok).toBe(true);
    if (r.ok) {
      expect(r.value.verify_strict).toBe(true);
      expect(r.value.quality_gates.lint).toBe("eslint .");
      expect(r.value.quality_gates.block_on_failure).toBe(false);
      expect(r.value.audit.divergence_threshold).toBe(0.3);
      expect(r.value.theme.bg).toBe("#000");
      expect(r.value.theme.borders).toBe("sharp");
      expect(r.value.theme.accents["done"]).toBe("#0f0");
    }
  });
});

describe("readConfig — not_found", () => {
  it("returns not_found when config.yml is missing", async () => {
    const r = await readConfig(tmp);
    expect(r.ok).toBe(false);
    if (!r.ok) {
      expect(r.error.code).toBe("not_found");
      expect(r.error.path).toContain("config.yml");
    }
  });
});

describe("readConfig — malformed_frontmatter", () => {
  it("returns malformed_frontmatter for invalid YAML", async () => {
    writeFileSync(configPath(), "key: [\nbad yaml");
    const r = await readConfig(tmp);
    expect(r.ok).toBe(false);
    if (!r.ok) {
      expect(r.error.code).toBe("malformed_frontmatter");
      expect("reason" in r.error).toBe(true);
    }
  });
});
