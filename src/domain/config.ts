export interface QualityGates {
  lint: string | null;
  typecheck: string | null;
  test: string | null;
  block_on_failure: boolean;
}

export interface AuditConfig {
  divergence_threshold: number;
}

export type BorderStyle = "subtle" | "sharp" | "none";

export interface ThemeConfig {
  bg: string;
  surface: string;
  surface_2: string;
  border: string;
  text: string;
  accents: Record<string, string>;
  borders: BorderStyle;
}

export interface SavepointConfig {
  verify_strict: boolean;
  quality_gates: QualityGates;
  audit: AuditConfig;
  theme: ThemeConfig;
}

export const CONFIG_DEFAULTS: SavepointConfig = {
  verify_strict: false,
  quality_gates: {
    lint: null,
    typecheck: null,
    test: null,
    block_on_failure: true,
  },
  audit: {
    divergence_threshold: 0.5,
  },
  theme: {
    bg: "#121212",
    surface: "#0D0D0D",
    surface_2: "#0F0F0F",
    border: "#1A1A1A",
    text: "#F0E6DA",
    accents: {
      backlog: "#B1A1DF",
      planned: "#B1A1DF",
      in_progress: "#FC6323",
      review: "#FC6323",
      done: "#A4C639",
      blocked: "#FF4444",
    },
    borders: "subtle",
  },
};

function isRecord(val: unknown): val is Record<string, unknown> {
  return typeof val === "object" && val !== null && !Array.isArray(val);
}

function stringEntriesOnly(
  raw: Record<string, unknown>,
): Record<string, string> {
  return Object.fromEntries(
    Object.entries(raw).filter((entry): entry is [string, string] => {
      return typeof entry[1] === "string";
    }),
  );
}

export function applyConfigDefaults(raw: unknown): SavepointConfig {
  const obj = isRecord(raw) ? raw : {};
  const qg = isRecord(obj["quality_gates"]) ? obj["quality_gates"] : {};
  const audit_cfg = isRecord(obj["audit"]) ? obj["audit"] : {};
  const theme = isRecord(obj["theme"]) ? obj["theme"] : {};
  const accentRaw = isRecord(theme["accents"])
    ? stringEntriesOnly(theme["accents"])
    : {};

  const borderRaw = theme["borders"];
  const borders: BorderStyle =
    borderRaw === "subtle" || borderRaw === "sharp" || borderRaw === "none"
      ? borderRaw
      : CONFIG_DEFAULTS.theme.borders;

  return {
    verify_strict:
      typeof obj["verify_strict"] === "boolean"
        ? obj["verify_strict"]
        : CONFIG_DEFAULTS.verify_strict,
    quality_gates: {
      lint:
        qg["lint"] === null || typeof qg["lint"] === "string"
          ? (qg["lint"] as string | null)
          : CONFIG_DEFAULTS.quality_gates.lint,
      typecheck:
        qg["typecheck"] === null || typeof qg["typecheck"] === "string"
          ? (qg["typecheck"] as string | null)
          : CONFIG_DEFAULTS.quality_gates.typecheck,
      test:
        qg["test"] === null || typeof qg["test"] === "string"
          ? (qg["test"] as string | null)
          : CONFIG_DEFAULTS.quality_gates.test,
      block_on_failure:
        typeof qg["block_on_failure"] === "boolean"
          ? qg["block_on_failure"]
          : CONFIG_DEFAULTS.quality_gates.block_on_failure,
    },
    audit: {
      divergence_threshold:
        typeof audit_cfg["divergence_threshold"] === "number"
          ? audit_cfg["divergence_threshold"]
          : CONFIG_DEFAULTS.audit.divergence_threshold,
    },
    theme: {
      bg:
        typeof theme["bg"] === "string"
          ? theme["bg"]
          : CONFIG_DEFAULTS.theme.bg,
      surface:
        typeof theme["surface"] === "string"
          ? theme["surface"]
          : CONFIG_DEFAULTS.theme.surface,
      surface_2:
        typeof theme["surface_2"] === "string"
          ? theme["surface_2"]
          : CONFIG_DEFAULTS.theme.surface_2,
      border:
        typeof theme["border"] === "string"
          ? theme["border"]
          : CONFIG_DEFAULTS.theme.border,
      text:
        typeof theme["text"] === "string"
          ? theme["text"]
          : CONFIG_DEFAULTS.theme.text,
      accents: { ...CONFIG_DEFAULTS.theme.accents, ...accentRaw },
      borders,
    },
  };
}
