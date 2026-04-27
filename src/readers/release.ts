import {
  readMarkdownFile,
  type FsResult,
  type MarkdownDoc,
} from "../fs/markdown.js";
import { savepointPath } from "../fs/project.js";
import {
  validateReleaseFrontmatter,
  type ReleaseFrontmatter,
} from "../domain/release.js";

export async function readReleasePrd(
  root: string,
  releaseTag: string,
): Promise<FsResult<MarkdownDoc<ReleaseFrontmatter>>> {
  const path = savepointPath(root, "releases", releaseTag, "PRD.md");
  return readMarkdownFile(path, validateReleaseFrontmatter);
}
