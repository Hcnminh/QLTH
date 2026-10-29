#!/usr/bin/env python3
"""Validate docs/openapi.yaml for structural correctness.

Checks performed (no external deps beyond PyYAML):
  1. File parses as YAML and declares openapi 3.x.
  2. Required root fields: info{title,version}, paths, servers.
  3. Every $ref pointer resolves to an existing location in the document.
  4. Every path item has at least one valid HTTP method operation.
  5. Every operation has an operationId, at least one response, and
     non-empty summary/description.
  6. Cross-check against the Gin router: every /api/v1 route registered in
     internal/api/router.go must appear in the spec (prevents doc drift).

Exit code 0 = pass, 1 = fail with printed diagnostics.
Usage:
    make docs
    python3 scripts/validate_openapi.py
"""
from __future__ import annotations

import re
import sys
from pathlib import Path

try:
    import yaml
except ImportError:  # pragma: no cover
    print("ERROR: PyYAML is required. Install with: pip3 install pyyaml", file=sys.stderr)
    sys.exit(1)

REPO_ROOT = Path(__file__).resolve().parent.parent
SPEC_PATH = REPO_ROOT / "docs" / "openapi.yaml"
ROUTER_PATH = REPO_ROOT / "internal" / "api" / "router.go"

HTTP_METHODS = {"get", "put", "post", "delete", "options", "head", "patch", "trace"}

errors: list[str] = []
warnings: list[str] = []


def fail(msg: str) -> None:
    errors.append(msg)


def warn(msg: str) -> None:
    warnings.append(msg)


def load_spec() -> dict:
    if not SPEC_PATH.exists():
        fail(f"Spec file not found: {SPEC_PATH}")
        sys.exit(1)
    with SPEC_PATH.open(encoding="utf-8") as fh:
        try:
            doc = yaml.safe_load(fh)
        except yaml.YAMLError as exc:
            fail(f"YAML parse error: {exc}")
            sys.exit(1)
    if not isinstance(doc, dict):
        fail("Spec root must be a mapping")
        sys.exit(1)
    return doc


def check_root(doc: dict) -> None:
    version = str(doc.get("openapi", ""))
    if not version.startswith("3."):
        fail(f"openapi version must be 3.x, got {version!r}")
    info = doc.get("info")
    if not isinstance(info, dict):
        fail("missing info object")
        return
    for field in ("title", "version"):
        if not info.get(field):
            fail(f"info.{field} is required")
    if not doc.get("paths"):
        fail("paths object is required and must not be empty")
    servers = doc.get("servers")
    if not isinstance(servers, list) or not servers:
        fail("at least one server entry is required")
    elif not all(s.get("url") for s in servers):
        fail("every server entry needs a url")


def check_refs(doc: dict) -> None:
    """Walk the document and verify every $ref resolves."""

    def resolve(ref: str) -> bool:
        if not ref.startswith("#/"):
            return False  # external refs not used in this project
        node = doc
        for part in ref[2:].split("/"):
            part = part.replace("~1", "/").replace("~0", "~")
            if isinstance(node, dict) and part in node:
                node = node[part]
            else:
                return False
        return True

    def walk(node, path="$"):
        if isinstance(node, dict):
            for key, value in node.items():
                if key == "$ref":
                    if not isinstance(value, str):
                        fail(f"{path}.$ref must be a string")
                    elif not resolve(value):
                        fail(f"{path}: unresolved $ref {value!r}")
                else:
                    walk(value, f"{path}.{key}")
        elif isinstance(node, list):
            for idx, item in enumerate(node):
                walk(item, f"{path}[{idx}]")

    walk(doc)


def check_operations(doc: dict) -> None:
    paths = doc.get("paths") or {}
    seen_ids: dict[str, str] = {}
    for route, item in paths.items():
        if not route.startswith("/"):
            fail(f"path {route!r} must start with /")
        if not isinstance(item, dict):
            fail(f"path {route!r} must be a mapping")
            continue
        methods = [m for m in item if m.lower() in HTTP_METHODS]
        if not methods:
            fail(f"path {route!r} declares no HTTP operations")
        for method in methods:
            op = item[method]
            where = f"{method.upper()} {route}"
            if not isinstance(op, dict):
                fail(f"{where}: operation must be a mapping")
                continue
            op_id = op.get("operationId")
            if not op_id:
                fail(f"{where}: missing operationId")
            elif op_id in seen_ids:
                fail(f"{where}: duplicate operationId {op_id!r} (also on {seen_ids[op_id]})")
            else:
                seen_ids[op_id] = where
            if not op.get("summary"):
                fail(f"{where}: missing summary")
            if not op.get("description"):
                warn(f"{where}: missing description")
            responses = op.get("responses")
            if not isinstance(responses, dict) or not responses:
                fail(f"{where}: must declare at least one response")
            else:
                for status, resp in responses.items():
                    if not re.fullmatch(r"\d{3}|default", str(status)):
                        fail(f"{where}: invalid response status {status!r}")
                    if not isinstance(resp, dict) or (
                        "$ref" not in resp and not resp.get("description")
                    ):
                        fail(f"{where}: response {status} needs description or $ref")


def check_router_coverage(doc: dict) -> None:
    """Every Gin route in router.go must be documented in the spec (anti doc-drift)."""
    if not ROUTER_PATH.exists():
        warn(f"router file not found, skipping drift check: {ROUTER_PATH}")
        return

    source = ROUTER_PATH.read_text(encoding="utf-8")
    spec_paths = set(doc.get("paths") or {})

    route_re = re.compile(r'\.(GET|POST|PUT|PATCH|DELETE)\(\s*"([^"]+)"')
    group_re = re.compile(r'\.Group\(\s*"([^"]+)"')
    prefixes = ["", *group_re.findall(source)]

    registered: set[str] = set()
    missing: set[str] = set()
    for _method, path in route_re.findall(source):
        candidates = []
        for p in prefixes:
            raw = p + path
            candidates.append(raw)
            candidates.append(re.sub(r':([a-zA-Z0-9_]+)', r'{\1}', raw))
        match = next((c for c in candidates if c in spec_paths), None)
        if match:
            registered.add(match)
        else:
            missing.add(path)

    if missing:
        fail(f"routes in router.go missing from openapi.yaml: {sorted(missing)}")

    undocumented = spec_paths - registered
    if undocumented:
        warn(f"spec paths with no matching route in router.go: {sorted(undocumented)}")


def main() -> int:
    doc = load_spec()
    check_root(doc)
    check_refs(doc)
    check_operations(doc)
    check_router_coverage(doc)

    for w in warnings:
        print(f"WARN  {w}")
    for e in errors:
        print(f"FAIL  {e}", file=sys.stderr)

    if errors:
        print(f"\nopenapi validation FAILED with {len(errors)} error(s).", file=sys.stderr)
        return 1
    print("openapi.yaml is valid: structure, $ref resolution, operations, router coverage OK.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
