import { describe, expect, it } from "vitest";
import {
  EXIT_NOT_IMPLEMENTED,
  EXIT_SUCCESS,
  EXIT_USAGE,
} from "../../src/cli/exit-codes.js";
import { runCli } from "../../src/cli/run.js";
import type { RunnerInput } from "../../src/cli/run.js";
import { version } from "../../src/version.js";

function makeStream() {
  const chunks: string[] = [];
  return {
    write(chunk: string) {
      chunks.push(chunk);
    },
    get text() {
      return chunks.join("");
    },
  };
}

function run(args: string[]): {
  exitCode: number;
  stdout: string;
  stderr: string;
} {
  const out = makeStream();
  const err = makeStream();
  const input: RunnerInput = {
    argv: ["node", "savepoint", ...args],
    stdin: {},
    stdout: out,
    stderr: err,
    env: {},
    platform: "linux",
  };
  const result = runCli(input);
  return { exitCode: result.exitCode, stdout: out.text, stderr: err.text };
}

describe("runCli — bare invocation", () => {
  it("prints help to stdout", () => {
    const { stdout } = run([]);
    expect(stdout).toContain("savepoint");
    expect(stdout).toContain("Commands:");
  });

  it("exits 0", () => {
    expect(run([]).exitCode).toBe(EXIT_SUCCESS);
  });

  it("writes nothing to stderr", () => {
    expect(run([]).stderr).toBe("");
  });
});

describe("runCli — global help flags", () => {
  it("--help exits 0", () => {
    expect(run(["--help"]).exitCode).toBe(EXIT_SUCCESS);
  });

  it("-h exits 0", () => {
    expect(run(["-h"]).exitCode).toBe(EXIT_SUCCESS);
  });

  it("--help prints top-level help to stdout", () => {
    expect(run(["--help"]).stdout).toContain("Commands:");
  });

  it("--help writes nothing to stderr", () => {
    expect(run(["--help"]).stderr).toBe("");
  });
});

describe("runCli — version flags", () => {
  it("--version exits 0", () => {
    expect(run(["--version"]).exitCode).toBe(EXIT_SUCCESS);
  });

  it("-v exits 0", () => {
    expect(run(["-v"]).exitCode).toBe(EXIT_SUCCESS);
  });

  it("--version prints the version string to stdout", () => {
    expect(run(["--version"]).stdout).toContain(version);
  });

  it("--version writes nothing to stderr", () => {
    expect(run(["--version"]).stderr).toBe("");
  });
});

describe("runCli — command dispatch: init", () => {
  it("dispatches init and returns EXIT_NOT_IMPLEMENTED", () => {
    expect(run(["init"]).exitCode).toBe(EXIT_NOT_IMPLEMENTED);
  });

  it("writes stub output to stdout", () => {
    expect(run(["init"]).stdout.length).toBeGreaterThan(0);
  });
});

describe("runCli — command dispatch: board", () => {
  it("dispatches board and returns EXIT_NOT_IMPLEMENTED", () => {
    expect(run(["board"]).exitCode).toBe(EXIT_NOT_IMPLEMENTED);
  });

  it("writes stub output to stdout", () => {
    expect(run(["board"]).stdout.length).toBeGreaterThan(0);
  });
});

describe("runCli — command dispatch: audit", () => {
  it("dispatches audit and returns EXIT_NOT_IMPLEMENTED", () => {
    expect(run(["audit"]).exitCode).toBe(EXIT_NOT_IMPLEMENTED);
  });

  it("writes stub output to stdout", () => {
    expect(run(["audit"]).stdout.length).toBeGreaterThan(0);
  });
});

describe("runCli — command dispatch: doctor", () => {
  it("dispatches doctor and returns EXIT_NOT_IMPLEMENTED", () => {
    expect(run(["doctor"]).exitCode).toBe(EXIT_NOT_IMPLEMENTED);
  });

  it("writes stub output to stdout", () => {
    expect(run(["doctor"]).stdout.length).toBeGreaterThan(0);
  });
});

describe("runCli — command-level help", () => {
  it("init --help exits 0", () => {
    expect(run(["init", "--help"]).exitCode).toBe(EXIT_SUCCESS);
  });

  it("init -h exits 0", () => {
    expect(run(["init", "-h"]).exitCode).toBe(EXIT_SUCCESS);
  });

  it("init --help prints command-specific help to stdout", () => {
    expect(run(["init", "--help"]).stdout).toContain("savepoint init");
  });

  it("board --help prints command-specific help to stdout", () => {
    expect(run(["board", "--help"]).stdout).toContain("savepoint board");
  });

  it("audit --help prints command-specific help to stdout", () => {
    expect(run(["audit", "--help"]).stdout).toContain("savepoint audit");
  });

  it("doctor --help prints command-specific help to stdout", () => {
    expect(run(["doctor", "--help"]).stdout).toContain("savepoint doctor");
  });

  it("command --help writes nothing to stderr", () => {
    expect(run(["init", "--help"]).stderr).toBe("");
  });
});

describe("runCli — unknown command", () => {
  it("exits with EXIT_USAGE", () => {
    expect(run(["plan"]).exitCode).toBe(EXIT_USAGE);
  });

  it("writes error to stderr", () => {
    expect(run(["plan"]).stderr).toContain("unknown command");
    expect(run(["plan"]).stderr).toContain("plan");
  });

  it("writes nothing to stdout", () => {
    expect(run(["plan"]).stdout).toBe("");
  });
});

describe("runCli — unknown flag", () => {
  it("exits with EXIT_USAGE", () => {
    expect(run(["--verbose"]).exitCode).toBe(EXIT_USAGE);
  });

  it("writes error to stderr", () => {
    expect(run(["--verbose"]).stderr).toContain("unknown flag");
    expect(run(["--verbose"]).stderr).toContain("--verbose");
  });

  it("writes nothing to stdout", () => {
    expect(run(["--verbose"]).stdout).toBe("");
  });

  it("exits with EXIT_USAGE for unknown command-level flags", () => {
    expect(run(["init", "--bad"]).exitCode).toBe(EXIT_USAGE);
  });

  it("writes command-level unknown flags to stderr", () => {
    const result = run(["init", "--bad"]);
    expect(result.stderr).toContain("unknown flag");
    expect(result.stderr).toContain("--bad");
    expect(result.stderr).toContain("savepoint init");
  });

  it("writes nothing to stdout for unknown command-level flags", () => {
    expect(run(["init", "--bad"]).stdout).toBe("");
  });
});
