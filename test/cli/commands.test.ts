import { describe, expect, it } from "vitest";
import { EXIT_NOT_IMPLEMENTED } from "../../src/cli/exit-codes.js";
import { runAudit } from "../../src/commands/audit.js";
import { runBoard } from "../../src/commands/board.js";
import { runDoctor } from "../../src/commands/doctor.js";
import { runInit } from "../../src/commands/init.js";

describe("runInit", () => {
  it("returns EXIT_NOT_IMPLEMENTED", () => {
    expect(runInit().exitCode).toBe(EXIT_NOT_IMPLEMENTED);
  });
  it("returns a non-empty message", () => {
    expect(runInit().message.length).toBeGreaterThan(0);
  });
});

describe("runBoard", () => {
  it("returns EXIT_NOT_IMPLEMENTED", () => {
    expect(runBoard().exitCode).toBe(EXIT_NOT_IMPLEMENTED);
  });
  it("returns a non-empty message", () => {
    expect(runBoard().message.length).toBeGreaterThan(0);
  });
});

describe("runAudit", () => {
  it("returns EXIT_NOT_IMPLEMENTED", () => {
    expect(runAudit().exitCode).toBe(EXIT_NOT_IMPLEMENTED);
  });
  it("returns a non-empty message", () => {
    expect(runAudit().message.length).toBeGreaterThan(0);
  });
});

describe("runDoctor", () => {
  it("returns EXIT_NOT_IMPLEMENTED", () => {
    expect(runDoctor().exitCode).toBe(EXIT_NOT_IMPLEMENTED);
  });
  it("returns a non-empty message", () => {
    expect(runDoctor().message.length).toBeGreaterThan(0);
  });
});
