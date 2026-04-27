export const PROJECT_TEMPLATES = [
  "AGENTS.md",
  ".savepoint/router.md",
  ".savepoint/PRD.md",
  ".savepoint/Design.md",
  ".savepoint/config.yml",
  ".savepoint/visual-identity.md",
] as const;

export const RELEASE_TEMPLATES = ["v1/PRD.md"] as const;

export const PROMPT_TEMPLATES = [
  "prd.prompt.md",
  "design.prompt.md",
  "epic-design.prompt.md",
  "task-breakdown.prompt.md",
  "task-planning.prompt.md",
  "task-building.prompt.md",
  "audit-reconciliation.prompt.md",
] as const;

export type ProjectTemplate = (typeof PROJECT_TEMPLATES)[number];
export type ReleaseTemplate = (typeof RELEASE_TEMPLATES)[number];
export type PromptTemplate = (typeof PROMPT_TEMPLATES)[number];
export type TemplateName = ProjectTemplate | ReleaseTemplate | PromptTemplate;
