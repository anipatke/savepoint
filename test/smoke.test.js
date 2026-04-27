import { describe, expect, it } from "vitest";
import { version } from "../src/version.js";

describe("version", () => {
  it("exports the scaffold version", () => {
    expect(version).toBe("0.1.0");
  });
});
