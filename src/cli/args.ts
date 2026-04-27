export type AllowedCommand = "init" | "board" | "audit" | "doctor";

export type ParseResult =
  | { kind: "bare" }
  | { kind: "help" }
  | { kind: "version" }
  | { kind: "command"; command: AllowedCommand }
  | { kind: "unknown-command"; name: string }
  | { kind: "unknown-flag"; flag: string };

const ALLOWED_COMMANDS = new Set<string>(["init", "board", "audit", "doctor"]);
const HELP_FLAGS = new Set(["-h", "--help"]);
const VERSION_FLAGS = new Set(["-v", "--version"]);

export function parseArgs(argv: string[]): ParseResult {
  const args = argv.slice(2);

  if (args.length === 0) return { kind: "bare" };

  const first = args[0];

  if (HELP_FLAGS.has(first)) return { kind: "help" };
  if (VERSION_FLAGS.has(first)) return { kind: "version" };
  if (first.startsWith("-")) return { kind: "unknown-flag", flag: first };
  if (ALLOWED_COMMANDS.has(first))
    return { kind: "command", command: first as AllowedCommand };

  return { kind: "unknown-command", name: first };
}
