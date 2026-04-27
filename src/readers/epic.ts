import {
  readMarkdownFile,
  type FsResult,
  type MarkdownDoc,
} from "../fs/markdown.js";
import { savepointPath } from "../fs/project.js";
import {
  validateEpicFrontmatter,
  type EpicFrontmatter,
} from "../domain/epic.js";

export async function readEpicDesign(
  root: string,
  releaseTag: string,
  epicRaw: string,
): Promise<FsResult<MarkdownDoc<EpicFrontmatter>>> {
  const path = savepointPath(
    root,
    "releases",
    releaseTag,
    "epics",
    epicRaw,
    "Design.md",
  );
  return readMarkdownFile(path, validateEpicFrontmatter);
}
