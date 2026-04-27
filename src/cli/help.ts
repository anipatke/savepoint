import type { AllowedCommand } from "./args.js";

const TOP_LEVEL = `\
savepoint — AI dev workflow checkpoint tool

Usage:
  savepoint <command> [options]

Commands:
  init     Initialize a savepoint project in the current directory
  board    Display the task board for the current epic
  audit    Run the audit loop for the completed epic
  doctor   Check project health and configuration

Global Flags:
  -h, --help       Show this help message
  -v, --version    Show the version number
`;

const COMMAND_HELP: Record<AllowedCommand, string> = {
  init: `\
savepoint init — Initialize a savepoint project

Usage:
  savepoint init

Initialize a .savepoint directory in the current working directory.
`,
  board: `\
savepoint board — Display the task board

Usage:
  savepoint board

Show the task board for the active epic in the current project.
`,
  audit: `\
savepoint audit — Run the audit loop for the completed epic

Usage:
  savepoint audit

Run the audit loop after all tasks in the current epic are done.
`,
  doctor: `\
savepoint doctor — Check project health

Usage:
  savepoint doctor

Validate the project configuration and surface any issues.
`,
};

export function topLevelHelp(): string {
  return TOP_LEVEL;
}

export function commandHelp(command: AllowedCommand): string {
  return COMMAND_HELP[command];
}
