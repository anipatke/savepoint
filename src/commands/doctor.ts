import { EXIT_NOT_IMPLEMENTED } from "../cli/exit-codes.js";
import type { CommandResult } from "./result.js";

export function runDoctor(): CommandResult {
  return {
    message: "savepoint doctor: not implemented yet",
    exitCode: EXIT_NOT_IMPLEMENTED,
  };
}
