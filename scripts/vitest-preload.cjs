const childProcess = require("node:child_process");
const esbuild = require("esbuild");

const originalExec = childProcess.exec;
const originalExecFile = childProcess.execFile;

function shouldStub(command) {
  return (
    typeof command === "string" && command.trim().toLowerCase() === "net use"
  );
}

function fakeChildProcess() {
  return {
    pid: 0,
    kill() {
      return true;
    },
    stdin: null,
    stdout: null,
    stderr: null,
    on() {
      return this;
    },
    once() {
      return this;
    },
    addListener() {
      return this;
    },
    removeListener() {
      return this;
    },
  };
}

// Vitest probes `net use` on Windows to detect network drive mappings, which fails in this scaffold.
childProcess.exec = function exec(command, options, callback) {
  if (shouldStub(command)) {
    const done = typeof options === "function" ? options : callback;
    if (typeof done === "function") {
      setImmediate(() => done(null, "", ""));
    }
    return fakeChildProcess();
  }

  return originalExec.call(this, command, options, callback);
};

childProcess.execFile = function execFile(file, args, options, callback) {
  if (shouldStub(file)) {
    const done =
      typeof args === "function"
        ? args
        : typeof options === "function"
          ? options
          : callback;
    if (typeof done === "function") {
      setImmediate(() => done(null, "", ""));
    }
    return fakeChildProcess();
  }

  return originalExecFile.call(this, file, args, options, callback);
};

function identityTransform(input) {
  const code =
    typeof input === "string" ? input : (input?.toString("utf8") ?? "");
  return Promise.resolve({
    code,
    map: "",
    warnings: [],
    errors: [],
    mangleCache: {},
    legalComments: "none",
  });
}

function identityTransformSync(input) {
  const code =
    typeof input === "string" ? input : (input?.toString("utf8") ?? "");
  return {
    code,
    map: "",
    warnings: [],
    errors: [],
    mangleCache: {},
    legalComments: "none",
  };
}

// Vitest reaches into esbuild on Windows before the real transform path is ready.
esbuild.transform = identityTransform;
esbuild.transformSync = identityTransformSync;
