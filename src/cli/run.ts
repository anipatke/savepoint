import type { StreamLike, EnvMap } from "./environment.js";
import { detectEnvironment } from "./environment.js";
import { parseArgs } from "./args.js";
import { topLevelHelp, commandHelp } from "./help.js";
import { EXIT_SUCCESS, EXIT_USAGE } from "./exit-codes.js";
import { version } from "../version.js";
import { runInit } from "../commands/init.js";
import { runBoard } from "../commands/board.js";
import { runAudit } from "../commands/audit.js";
import { runDoctor } from "../commands/doctor.js";

export interface WritableStream extends StreamLike {
  write(chunk: string): void;
}

export interface RunnerInput {
  argv: string[];
  stdin: StreamLike;
  stdout: WritableStream;
  stderr: WritableStream;
  env: EnvMap;
  platform: string;
}

export interface RunnerResult {
  exitCode: number;
}

const HELP_FLAGS = new Set(["-h", "--help"]);

const COMMAND_HANDLERS = {
  init: runInit,
  board: runBoard,
  audit: runAudit,
  doctor: runDoctor,
} as const;

export function runCli(input: RunnerInput): RunnerResult {
  const { argv, stdout, stderr, env, platform } = input;
  detectEnvironment(stdout, stderr, env, platform);

  const parsed = parseArgs(argv);

  if (parsed.kind === "bare" || parsed.kind === "help") {
    stdout.write(topLevelHelp());
    return { exitCode: EXIT_SUCCESS };
  }

  if (parsed.kind === "version") {
    stdout.write(`${version}\n`);
    return { exitCode: EXIT_SUCCESS };
  }

  if (parsed.kind === "command") {
    const nextArg = argv[3];
    if (nextArg !== undefined && HELP_FLAGS.has(nextArg)) {
      stdout.write(commandHelp(parsed.command));
      return { exitCode: EXIT_SUCCESS };
    }
    if (nextArg !== undefined && nextArg.startsWith("-")) {
      stderr.write(`savepoint ${parsed.command}: unknown flag '${nextArg}'\n`);
      stderr.write(commandHelp(parsed.command));
      return { exitCode: EXIT_USAGE };
    }
    const result = COMMAND_HANDLERS[parsed.command]();
    stdout.write(`${result.message}\n`);
    return { exitCode: result.exitCode };
  }

  if (parsed.kind === "unknown-command") {
    stderr.write(`savepoint: unknown command '${parsed.name}'\n`);
    stderr.write(topLevelHelp());
    return { exitCode: EXIT_USAGE };
  }

  // unknown-flag
  stderr.write(`savepoint: unknown flag '${parsed.flag}'\n`);
  stderr.write(topLevelHelp());
  return { exitCode: EXIT_USAGE };
}
