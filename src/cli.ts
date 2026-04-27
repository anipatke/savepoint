import { runCli } from "./cli/run.js";

const result = runCli({
  argv: process.argv,
  stdin: process.stdin,
  stdout: process.stdout,
  stderr: process.stderr,
  env: process.env,
  platform: process.platform,
});

process.exit(result.exitCode);
