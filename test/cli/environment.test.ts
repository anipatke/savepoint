import { describe, expect, it } from "vitest";
import { detectEnvironment } from "../../src/cli/environment.js";

const tty = { isTTY: true as const };
const notTTY = { isTTY: false as const };
const noStream = {};

describe("detectEnvironment — TTY detection", () => {
  it("stdout isTTY true when stream.isTTY is true", () => {
    const r = detectEnvironment(tty, noStream, {}, "linux");
    expect(r.stdoutIsTTY).toBe(true);
  });

  it("stdout isTTY false when stream.isTTY is false", () => {
    const r = detectEnvironment(notTTY, noStream, {}, "linux");
    expect(r.stdoutIsTTY).toBe(false);
  });

  it("stdout isTTY false when isTTY is absent", () => {
    const r = detectEnvironment(noStream, noStream, {}, "linux");
    expect(r.stdoutIsTTY).toBe(false);
  });

  it("stderr isTTY true when stream.isTTY is true", () => {
    const r = detectEnvironment(noStream, tty, {}, "linux");
    expect(r.stderrIsTTY).toBe(true);
  });

  it("stderr isTTY false when isTTY is absent", () => {
    const r = detectEnvironment(noStream, noStream, {}, "linux");
    expect(r.stderrIsTTY).toBe(false);
  });
});

describe("detectEnvironment — color detection", () => {
  it("color enabled when stdout is TTY and no overrides", () => {
    const r = detectEnvironment(tty, noStream, {}, "linux");
    expect(r.colorEnabled).toBe(true);
  });

  it("color disabled when stdout is not TTY and no overrides", () => {
    const r = detectEnvironment(notTTY, noStream, {}, "linux");
    expect(r.colorEnabled).toBe(false);
  });

  it("NO_COLOR disables color even when TTY", () => {
    const r = detectEnvironment(tty, noStream, { NO_COLOR: "" }, "linux");
    expect(r.colorEnabled).toBe(false);
  });

  it("FORCE_COLOR enables color even when not TTY", () => {
    const r = detectEnvironment(
      notTTY,
      noStream,
      { FORCE_COLOR: "1" },
      "linux",
    );
    expect(r.colorEnabled).toBe(true);
  });

  it("FORCE_COLOR=0 disables color even when TTY", () => {
    const r = detectEnvironment(tty, noStream, { FORCE_COLOR: "0" }, "linux");
    expect(r.colorEnabled).toBe(false);
  });

  it("CI disables color even when TTY", () => {
    const r = detectEnvironment(tty, noStream, { CI: "true" }, "linux");
    expect(r.colorEnabled).toBe(false);
  });

  it("TERM=dumb disables color even when TTY", () => {
    const r = detectEnvironment(tty, noStream, { TERM: "dumb" }, "linux");
    expect(r.colorEnabled).toBe(false);
  });

  it("NO_COLOR takes precedence over FORCE_COLOR", () => {
    const r = detectEnvironment(
      tty,
      noStream,
      { NO_COLOR: "", FORCE_COLOR: "1" },
      "linux",
    );
    expect(r.colorEnabled).toBe(false);
  });

  it("FORCE_COLOR takes precedence over CI", () => {
    const r = detectEnvironment(
      notTTY,
      noStream,
      { FORCE_COLOR: "1", CI: "true" },
      "linux",
    );
    expect(r.colorEnabled).toBe(true);
  });
});

describe("detectEnvironment — platform passthrough", () => {
  it("passes through platform unchanged", () => {
    expect(detectEnvironment(noStream, noStream, {}, "win32").platform).toBe(
      "win32",
    );
    expect(detectEnvironment(noStream, noStream, {}, "darwin").platform).toBe(
      "darwin",
    );
    expect(detectEnvironment(noStream, noStream, {}, "linux").platform).toBe(
      "linux",
    );
  });
});
