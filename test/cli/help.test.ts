import { describe, expect, it } from "vitest";
import { commandHelp, topLevelHelp } from "../../src/cli/help.js";

describe("topLevelHelp", () => {
  it("returns a string", () => {
    expect(typeof topLevelHelp()).toBe("string");
  });

  it("lists all four commands", () => {
    const text = topLevelHelp();
    expect(text).toContain("init");
    expect(text).toContain("board");
    expect(text).toContain("audit");
    expect(text).toContain("doctor");
  });

  it("lists the global help flags", () => {
    expect(topLevelHelp()).toContain("--help");
    expect(topLevelHelp()).toContain("-h");
  });

  it("lists the global version flags", () => {
    expect(topLevelHelp()).toContain("--version");
    expect(topLevelHelp()).toContain("-v");
  });

  it("is deterministic", () => {
    expect(topLevelHelp()).toBe(topLevelHelp());
  });
});

describe("commandHelp — init", () => {
  it("returns a string containing init", () => {
    expect(commandHelp("init")).toContain("init");
  });

  it("is deterministic", () => {
    expect(commandHelp("init")).toBe(commandHelp("init"));
  });
});

describe("commandHelp — board", () => {
  it("returns a string containing board", () => {
    expect(commandHelp("board")).toContain("board");
  });

  it("is deterministic", () => {
    expect(commandHelp("board")).toBe(commandHelp("board"));
  });
});

describe("commandHelp — audit", () => {
  it("returns a string containing audit", () => {
    expect(commandHelp("audit")).toContain("audit");
  });

  it("is deterministic", () => {
    expect(commandHelp("audit")).toBe(commandHelp("audit"));
  });
});

describe("commandHelp — doctor", () => {
  it("returns a string containing doctor", () => {
    expect(commandHelp("doctor")).toContain("doctor");
  });

  it("is deterministic", () => {
    expect(commandHelp("doctor")).toBe(commandHelp("doctor"));
  });
});
