export type TemplateVars = {
  projectName?: string;
  releaseNumber?: string;
  releaseName?: string;
};

export function renderTemplate(content: string, vars: TemplateVars): string {
  let result = content;
  if (vars.projectName !== undefined) {
    result = result.replaceAll(
      /\{\s*\{\s*PROJECT_NAME\s*\}\s*\}/g,
      vars.projectName,
    );
  }
  if (vars.releaseNumber !== undefined) {
    result = result.replaceAll(
      /\{\s*\{\s*RELEASE_NUMBER\s*\}\s*\}/g,
      vars.releaseNumber,
    );
  }
  if (vars.releaseName !== undefined) {
    result = result.replaceAll(
      /\{\s*\{\s*RELEASE_NAME\s*\}\s*\}/g,
      vars.releaseName,
    );
  }
  return result;
}
