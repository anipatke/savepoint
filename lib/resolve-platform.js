"use strict";

const PLATFORM_MAP = Object.freeze({
  darwin: "darwin",
  linux: "linux",
  win32: "windows",
});

const ARCH_MAP = Object.freeze({
  arm64: "arm64",
  x64: "amd64",
});

const NPM_OS_MAP = Object.freeze({
  darwin: "darwin",
  linux: "linux",
  windows: "win32",
});

const NPM_CPU_MAP = Object.freeze({
  amd64: "x64",
  arm64: "arm64",
});

const SUPPORTED_TARGETS = Object.freeze([
  Object.freeze({ goos: "linux", goarch: "amd64" }),
  Object.freeze({ goos: "linux", goarch: "arm64" }),
  Object.freeze({ goos: "darwin", goarch: "amd64" }),
  Object.freeze({ goos: "darwin", goarch: "arm64" }),
  Object.freeze({ goos: "windows", goarch: "amd64" }),
  Object.freeze({ goos: "windows", goarch: "arm64" }),
]);

function platformPackageName(goos, goarch) {
  return `savepoint-${goos}-${goarch}`;
}

function executableName(goos) {
  return goos === "windows" ? "savepoint.exe" : "savepoint";
}

function resolvePlatform(input) {
  const platform = input && input.platform;
  const arch = input && input.arch;
  const goos = PLATFORM_MAP[platform];
  const goarch = ARCH_MAP[arch];
  if (!goos || !goarch) {
    return {
      supported: false,
      error: `savepoint does not support ${platform}/${arch}`,
    };
  }
  const executable = executableName(goos);
  return {
    supported: true,
    goos,
    goarch,
    slug: `${goos}-${goarch}`,
    packageName: platformPackageName(goos, goarch),
    binaryPath: `bin/${executable}`,
    executable,
  };
}

module.exports = {
  PLATFORM_MAP,
  ARCH_MAP,
  NPM_OS_MAP,
  NPM_CPU_MAP,
  SUPPORTED_TARGETS,
  platformPackageName,
  executableName,
  resolvePlatform,
};
