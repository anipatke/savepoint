import { EXIT_NOT_IMPLEMENTED } from "../cli/exit-codes.js";
import type { CommandResult } from "./result.js";

export function runAudit(): CommandResult {
  return {
    message: "savepoint audit: not implemented yet",
    exitCode: EXIT_NOT_IMPLEMENTED,
  };
}
