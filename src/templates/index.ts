export {
  PROJECT_TEMPLATES,
  RELEASE_TEMPLATES,
  PROMPT_TEMPLATES,
  type ProjectTemplate,
  type ReleaseTemplate,
  type PromptTemplate,
  type TemplateName,
} from "./manifest.js";

export { defaultTemplateRoot, templatePath } from "./paths.js";
export { renderTemplate, type TemplateVars } from "./render.js";
export {
  loadTemplate,
  type TemplateResult,
  type TemplateError,
} from "./load.js";
