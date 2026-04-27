import { EXIT_NOT_IMPLEMENTED } from "../cli/exit-codes.js";
import type { CommandResult } from "./result.js";

export function runInit(): CommandResult {
  return {
    message: "savepoint init: not implemented yet",
    exitCode: EXIT_NOT_IMPLEMENTED,
  };
}
