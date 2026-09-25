#!/usr/bin/env python3
"""Rename Savepoint V2 record identities to their hyphenated form.

A Savepoint V2 project numbers its records with a bare `[ROTCI]###` identity
(`O014`, `T019`, ...). This script converts a project's own active records to
the hyphenated form (`O-014`, `T-019`, ...): it renames the identity-named
directories and files under `.savepoint/{releases,objectives,checks,issues}`
and rewrites every word-bounded bare identity to its hyphenated form across
those renamed records plus the project's `.savepoint/router.md`, `AGENTS.md`,
`.savepoint/Design.md`, and `.savepoint/Guardrails.md`.

`.savepoint/archive/` is never touched, and only entries whose own name
starts with a bare identity are renamed or recursed into: a release
directory such as `.savepoint/releases/v1` that predates the V1-to-V2
conversion, and whose own numbering scheme is not this identity form, is
left exactly as it is.

Safe to run more than once: a name or reference already in hyphenated form
no longer matches the bare pattern, so a second run makes no changes.

Usage: scripts/hyphenate_record_ids.py [--root PATH] [--dry-run]
"""
import argparse
import re
import subprocess
import sys
from pathlib import Path

IDENTITY_KINDS = "ROTCI"
BARE_IDENTITY = re.compile(rf"^([{IDENTITY_KINDS}])([0-9]{{3,}})\b")
BARE_IDENTITY_ANYWHERE = re.compile(rf"\b([{IDENTITY_KINDS}])([0-9]{{3,}})\b")

RECORD_DIRS = ("releases", "objectives", "checks", "issues")
EXTRA_FILES = ("router.md", "Design.md", "Guardrails.md")


def hyphenated_name(name: str) -> str | None:
    m = BARE_IDENTITY.match(name)
    if not m:
        return None
    return f"{m.group(1)}-{m.group(2)}{name[m.end():]}"


def git_mv(src: Path, dst: Path, dry_run: bool) -> None:
    print(f"rename {src} -> {dst}")
    if dry_run:
        return
    result = subprocess.run(
        ["git", "mv", str(src), str(dst)], capture_output=True, text=True
    )
    if result.returncode != 0:
        src.rename(dst)


def collect_and_rename(dir_path: Path, restrict_to_matches: bool, dry_run: bool, files: list[Path]) -> None:
    for name in sorted(p.name for p in dir_path.iterdir()):
        child = dir_path / name
        new_name = hyphenated_name(name)
        if restrict_to_matches and new_name is None:
            continue  # not one of this project's identity-named records; leave untouched
        if new_name is not None and new_name != name:
            new_child = dir_path / new_name
            git_mv(child, new_child, dry_run)
            if not dry_run:
                child = new_child
        if child.is_dir():
            collect_and_rename(child, restrict_to_matches=False, dry_run=dry_run, files=files)
        else:
            files.append(child)


def rewrite_references(path: Path, dry_run: bool) -> bool:
    text = path.read_text(encoding="utf-8")
    new_text = BARE_IDENTITY_ANYWHERE.sub(lambda m: f"{m.group(1)}-{m.group(2)}", text)
    if new_text == text:
        return False
    print(f"rewrite {path}")
    if not dry_run:
        path.write_text(new_text, encoding="utf-8")
    return True


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--root", default=".", help="project root containing .savepoint (default: current directory)")
    parser.add_argument("--dry-run", action="store_true", help="print what would change without changing it")
    args = parser.parse_args()

    root = Path(args.root).resolve()
    savepoint_dir = root / ".savepoint"
    if not savepoint_dir.is_dir():
        print(f"{savepoint_dir}: not a Savepoint project (no .savepoint directory)", file=sys.stderr)
        return 1

    files: list[Path] = []
    for kind in RECORD_DIRS:
        base = savepoint_dir / kind
        if base.is_dir():
            collect_and_rename(base, restrict_to_matches=True, dry_run=args.dry_run, files=files)

    for extra in EXTRA_FILES:
        candidate = savepoint_dir / extra
        if candidate.is_file():
            files.append(candidate)
    agents_md = root / "AGENTS.md"
    if agents_md.is_file():
        files.append(agents_md)

    changed = 0
    for f in files:
        if rewrite_references(f, args.dry_run):
            changed += 1

    print(f"{'would rename/rewrite' if args.dry_run else 'renamed/rewrote'} identities across {len(files)} file(s), {changed} changed by content")
    return 0


if __name__ == "__main__":
    sys.exit(main())
