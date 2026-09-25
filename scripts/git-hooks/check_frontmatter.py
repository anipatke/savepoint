#!/usr/bin/env python3
"""Validate a single .savepoint markdown file's YAML frontmatter.

Reads the file content from stdin (so the pre-commit hook can check the
staged version, not the working-tree version) and takes the display path
as argv[1] for error messages.
"""
import sys

import yaml


def main() -> int:
    path = sys.argv[1] if len(sys.argv) > 1 else "<stdin>"
    text = sys.stdin.read()

    if not text.startswith("---\n"):
        print(f"{path}: frontmatter must start with '---' on the first line")
        return 1

    end = text.find("\n---", 4)
    if end == -1:
        print(f"{path}: no closing frontmatter delimiter found")
        return 1

    front = text[4:end]
    try:
        doc = yaml.safe_load(front)
    except yaml.YAMLError as exc:
        print(f"{path}: invalid YAML frontmatter: {exc}")
        return 1

    if not isinstance(doc, dict):
        print(f"{path}: frontmatter must decode to a mapping")
        return 1

    return 0


if __name__ == "__main__":
    sys.exit(main())
