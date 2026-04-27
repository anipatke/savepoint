export interface StreamLike {
  isTTY?: boolean;
}

export interface EnvMap {
  NO_COLOR?: string;
  FORCE_COLOR?: string;
  CI?: string;
  TERM?: string;
}

export interface TerminalEnv {
  stdoutIsTTY: boolean;
  stderrIsTTY: boolean;
  colorEnabled: boolean;
  platform: string;
}

export function detectEnvironment(
  stdout: StreamLike,
  stderr: StreamLike,
  env: EnvMap,
  platform: string,
): TerminalEnv {
  const stdoutIsTTY = stdout.isTTY === true;
  const stderrIsTTY = stderr.isTTY === true;
  const colorEnabled = resolveColor(stdoutIsTTY, env);
  return { stdoutIsTTY, stderrIsTTY, colorEnabled, platform };
}

function resolveColor(stdoutIsTTY: boolean, env: EnvMap): boolean {
  if (env.NO_COLOR !== undefined) return false;
  if (env.FORCE_COLOR !== undefined) return env.FORCE_COLOR !== "0";
  if (env.CI !== undefined) return false;
  if (env.TERM === "dumb") return false;
  return stdoutIsTTY;
}
