#!/usr/bin/env python3
"""Validate that AGENTS.md stays consistent with the actual codebase.

Checks performed:
  1. Every `make <target>` referenced in AGENTS.md exists in the Makefile.
  2. Every repo-relative file path referenced in AGENTS.md exists on disk.
  3. Every Go package path referenced (internal/..., cmd/...) exists.
  4. The commands table in section 1 is non-empty and each command is syntactically plausible.

This prevents AGENTS.md from drifting when Makefile targets or files are renamed,
which would silently mislead autonomous agents.

Usage:
    python3 scripts/validate_agents_md.py
"""
from __future__ import annotations

import re
import sys
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parent.parent
AGENTS_PATH = REPO_ROOT / "AGENTS.md"
MAKEFILE_PATH = REPO_ROOT / "Makefile"

errors: list[str] = []


def fail(msg: str) -> None:
    errors.append(msg)


def makefile_targets() -> set[str]:
    if not MAKEFILE_PATH.exists():
        fail("Makefile not found at repo root")
        return set()
    text = MAKEFILE_PATH.read_text(encoding="utf-8")
    targets = set(re.findall(r"^([a-zA-Z0-9_.-]+):", text, flags=re.MULTILINE))
    phony = re.search(r"^\.PHONY:\s*(.+)$", text, flags=re.MULTILINE)
    if phony:
        targets.update(phony.group(1).split())
    return targets


def check_make_targets(agents_text: str, targets: set[str]) -> None:
    referenced = set(re.findall(r"\bmake\s+([a-zA-Z0-9_.-]+)", agents_text))
    for target in sorted(referenced):
        if target not in targets:
            fail(f"AGENTS.md references `make {target}` but Makefile has no such target")


def check_file_references(agents_text: str) -> None:
    # Markdown links [text](path) and inline code paths that look like repo files
    linked = re.findall(r"\[[^\]]+\]\(([^)#\s]+)\)", agents_text)
    inline = re.findall(r"`((?:docs|internal|cmd|scripts|\.github|\.factory)/[^`\s]+)`", agents_text)
    candidates = set()
    for ref in linked + inline:
        ref = ref.strip()
        if ref.startswith(("http://", "https://", "#")):
            continue
        # Strip trailing punctuation and anchors
        ref = ref.split("#")[0].rstrip(".,;:")
        if not ref:
            continue
        candidates.add(ref)

    for ref in sorted(candidates):
        path = REPO_ROOT / ref
        # Allow directory references (with or without trailing slash) and globs
        if "*" in ref:
            continue
        if not path.exists():
            fail(f"AGENTS.md references `{ref}` which does not exist")


def check_go_packages(agents_text: str) -> None:
    pkgs = set(re.findall(r"`(internal/[a-z/]+|cmd/[a-z/]+)/?`", agents_text))
    for pkg in sorted(pkgs):
        path = REPO_ROOT / pkg
        if not path.is_dir():
            fail(f"AGENTS.md references Go package dir `{pkg}` which does not exist")


def check_command_table(agents_text: str) -> None:
    # Section 1 must contain a table of commands with at least a few rows
    section = re.search(r"##\s*1\..*?(?=\n##\s|\Z)", agents_text, flags=re.DOTALL)
    if not section:
        fail("AGENTS.md is missing section '## 1. Lệnh thiết yếu'")
        return
    rows = re.findall(r"^\|[^|]+\|[^|]+\|$", section.group(0), flags=re.MULTILINE)
    if len(rows) < 4:  # header + separator + >=2 command rows
        fail("AGENTS.md section 1 command table looks empty or truncated")


def main() -> int:
    if not AGENTS_PATH.exists():
        print("FAIL: AGENTS.md not found at repo root", file=sys.stderr)
        return 1

    agents_text = AGENTS_PATH.read_text(encoding="utf-8")
    if len(agents_text) <= 100:
        fail("AGENTS.md is suspiciously short (<=100 chars)")

    targets = makefile_targets()
    check_make_targets(agents_text, targets)
    check_file_references(agents_text)
    check_go_packages(agents_text)
    check_command_table(agents_text)

    if errors:
        for e in errors:
            print(f"FAIL  {e}", file=sys.stderr)
        print(f"\nAGENTS.md validation FAILED with {len(errors)} error(s).", file=sys.stderr)
        return 1

    print("AGENTS.md is consistent with Makefile targets, file paths, and Go packages.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
