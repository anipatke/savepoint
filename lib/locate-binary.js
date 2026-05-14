"use strict";

const fs = require("node:fs");
const path = require("node:path");

function locateBinary(resolved, options) {
  const opts = options || {};
  const resolver = opts.resolver || ((id) => require.resolve(id));
  const existsSync = opts.existsSync || fs.existsSync;
  const devRoot =
    opts.devRoot || path.join(__dirname, "..", "dist", "npm");

  try {
    return resolver(`${resolved.packageName}/${resolved.binaryPath}`);
  } catch (_) {
    // Fall through to dev fallback.
  }

  const local = path.join(devRoot, resolved.slug, resolved.binaryPath);
  if (existsSync(local)) {
    return local;
  }
  return null;
}

module.exports = { locateBinary };
