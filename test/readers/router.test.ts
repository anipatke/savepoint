import { afterEach, beforeEach, describe, expect, it } from "vitest";
import { mkdtempSync, rmSync, mkdirSync, writeFileSync } from "node:fs";
import { join } from "node:path";
import { tmpdir } from "node:os";
import { readRouterState } from "../../src/readers/router.js";

let tmp: string;

beforeEach(() => {
  tmp = mkdtempSync(join(tmpdir(), "savepoint-router-test-"));
  mkdirSync(join(tmp, ".savepoint"), { recursive: true });
});

afterEach(() => {
  rmSync(tmp, { recursive: true });
});

function routerPath(): string {
  return join(tmp, ".savepoint", "router.md");
}

function makeRouterMd(yaml: string): string {
  return (
    "# Router\n\n## Current state\n\n```yaml\n" +
    yaml +
    "\n```\n\n## Other section\n"
  );
}

describe("readRouterState — success", () => {
  it("parses a valid router state block", async () => {
    writeFileSync(
      routerPath(),
      makeRouterMd(
        'state: task-building\nrelease: v1\nepic: E02-data-model\nnext_action: "Build T004."',
      ),
    );
    const r = await readRouterState(tmp);
    expect(r.ok).toBe(true);
    if (r.ok) {
      expect(r.value.state).toBe("task-building");
      expect(r.value.release).toBe("v1");
      expect(r.value.epic).toBe("E02-data-model");
      expect(r.value.next_action).toBe("Build T004.");
    }
  });

  it("parses a state block without the optional epic field", async () => {
    writeFileSync(
      routerPath(),
      makeRouterMd(
        'state: pre-implementation\nrelease: v1\nnext_action: "Start here."',
      ),
    );
    const r = await readRouterState(tmp);
    expect(r.ok).toBe(true);
    if (r.ok) {
      expect(r.value.state).toBe("pre-implementation");
      expect(r.value.epic).toBeUndefined();
    }
  });
});

describe("readRouterState — not_found", () => {
  it("returns not_found when router.md is missing", async () => {
    const r = await readRouterState(tmp);
    expect(r.ok).toBe(false);
    if (!r.ok) {
      expect(r.error.code).toBe("not_found");
      expect(r.error.path).toContain("router.md");
    }
  });
});

describe("readRouterState — missing_frontmatter", () => {
  it("returns missing_frontmatter when there is no ## Current state block", async () => {
    writeFileSync(routerPath(), "# Router\n\nNo state block here.\n");
    const r = await readRouterState(tmp);
    expect(r.ok).toBe(false);
    if (!r.ok) expect(r.error.code).toBe("missing_frontmatter");
  });
});

describe("readRouterState — malformed_frontmatter", () => {
  it("returns malformed_frontmatter for invalid YAML in the state block", async () => {
    writeFileSync(routerPath(), makeRouterMd("state: [\nbad yaml"));
    const r = await readRouterState(tmp);
    expect(r.ok).toBe(false);
    if (!r.ok) {
      expect(r.error.code).toBe("malformed_frontmatter");
      expect("reason" in r.error).toBe(true);
    }
  });
});

describe("readRouterState — schema_rejection", () => {
  it("returns schema_rejection for an unknown state value", async () => {
    writeFileSync(
      routerPath(),
      makeRouterMd('state: not-a-state\nrelease: v1\nnext_action: "x"'),
    );
    const r = await readRouterState(tmp);
    expect(r.ok).toBe(false);
    if (!r.ok) {
      expect(r.error.code).toBe("schema_rejection");
      expect("reason" in r.error && r.error.reason).toMatch(/state/);
    }
  });

  it("returns schema_rejection when release is missing", async () => {
    writeFileSync(
      routerPath(),
      makeRouterMd('state: task-building\nnext_action: "x"'),
    );
    const r = await readRouterState(tmp);
    expect(r.ok).toBe(false);
    if (!r.ok) expect(r.error.code).toBe("schema_rejection");
  });
});
