import { EXIT_NOT_IMPLEMENTED } from "../cli/exit-codes.js";
import type { CommandResult } from "./result.js";

export function runBoard(): CommandResult {
  return {
    message: "savepoint board: not implemented yet",
    exitCode: EXIT_NOT_IMPLEMENTED,
  };
}
