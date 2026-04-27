import { describe, expect, it } from "vitest";
import { parseArgs } from "../../src/cli/args.js";

const argv = (args: string[]) => ["node", "savepoint", ...args];

describe("parseArgs — bare invocation", () => {
  it("returns bare when no args", () => {
    expect(parseArgs(argv([]))).toEqual({ kind: "bare" });
  });
});

describe("parseArgs — help flags", () => {
  it("recognises --help", () => {
    expect(parseArgs(argv(["--help"]))).toEqual({ kind: "help" });
  });

  it("recognises -h", () => {
    expect(parseArgs(argv(["-h"]))).toEqual({ kind: "help" });
  });
});

describe("parseArgs — version flags", () => {
  it("recognises --version", () => {
    expect(parseArgs(argv(["--version"]))).toEqual({ kind: "version" });
  });

  it("recognises -v", () => {
    expect(parseArgs(argv(["-v"]))).toEqual({ kind: "version" });
  });
});

describe("parseArgs — known commands", () => {
  it("recognises init", () => {
    expect(parseArgs(argv(["init"]))).toEqual({
      kind: "command",
      command: "init",
    });
  });

  it("recognises board", () => {
    expect(parseArgs(argv(["board"]))).toEqual({
      kind: "command",
      command: "board",
    });
  });

  it("recognises audit", () => {
    expect(parseArgs(argv(["audit"]))).toEqual({
      kind: "command",
      command: "audit",
    });
  });

  it("recognises doctor", () => {
    expect(parseArgs(argv(["doctor"]))).toEqual({
      kind: "command",
      command: "doctor",
    });
  });
});

describe("parseArgs — unknown command", () => {
  it("returns unknown-command for unrecognised word", () => {
    expect(parseArgs(argv(["plan"]))).toEqual({
      kind: "unknown-command",
      name: "plan",
    });
  });
});

describe("parseArgs — unknown flag", () => {
  it("returns unknown-flag for unrecognised dash-prefixed arg", () => {
    expect(parseArgs(argv(["--verbose"]))).toEqual({
      kind: "unknown-flag",
      flag: "--verbose",
    });
  });

  it("returns unknown-flag for single-dash unrecognised arg", () => {
    expect(parseArgs(argv(["-x"]))).toEqual({
      kind: "unknown-flag",
      flag: "-x",
    });
  });
});
