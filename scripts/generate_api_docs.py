#!/usr/bin/env python3
"""Generate docs/api-reference.md from docs/openapi.yaml.

Single source of truth is the OpenAPI spec; the Markdown reference is a
generated artifact so it can never drift from the spec.

Usage:
    python3 scripts/generate_api_docs.py            # regenerate the file
    python3 scripts/generate_api_docs.py --check    # CI mode: fail if the
                                                    # committed file is stale
"""
from __future__ import annotations

import argparse
import sys
from pathlib import Path

try:
    import yaml
except ImportError:  # pragma: no cover
    print("ERROR: PyYAML is required. Install with: pip3 install pyyaml", file=sys.stderr)
    sys.exit(1)

REPO_ROOT = Path(__file__).resolve().parent.parent
SPEC_PATH = REPO_ROOT / "docs" / "openapi.yaml"
OUT_PATH = REPO_ROOT / "docs" / "api-reference.md"

HTTP_METHODS = ["get", "post", "put", "patch", "delete"]


def load_spec() -> dict:
    with SPEC_PATH.open(encoding="utf-8") as fh:
        return yaml.safe_load(fh)


def resolve_ref(doc: dict, ref: str):
    node = doc
    for part in ref.lstrip("#/").split("/"):
        node = node[part.replace("~1", "/").replace("~0", "~")]
    return node


def deref(doc: dict, node):
    if isinstance(node, dict) and "$ref" in node:
        return resolve_ref(doc, node["$ref"])
    return node


def schema_type_label(doc: dict, schema: dict) -> str:
    schema = deref(doc, schema) or {}
    if "allOf" in schema:
        return "object"
    typ = schema.get("type", "object")
    if typ == "array":
        items = deref(doc, schema.get("items", {})) or {}
        return f"array<{items.get('type', 'object')}>"
    fmt = schema.get("format")
    return f"{typ} ({fmt})" if fmt else str(typ)


def render_parameters(doc: dict, params: list) -> str:
    if not params:
        return "_Không có tham số._\n"
    lines = ["| Tên | Vị trí | Kiểu | Bắt buộc | Mô tả |", "|---|---|---|---|---|"]
    for p in params:
        p = deref(doc, p)
        sch = deref(doc, p.get("schema", {})) or {}
        enum = sch.get("enum")
        desc = (p.get("description") or "").replace("\n", " ")
        if enum:
            desc += f" Giá trị: `{'`, `'.join(map(str, enum))}`."
        lines.append(
            f"| `{p.get('name')}` | {p.get('in')} | {sch.get('type', 'string')} "
            f"| {'Có' if p.get('required') else 'Không'} | {desc} |"
        )
    return "\n".join(lines) + "\n"


def render_request_body(doc: dict, op: dict) -> str:
    body = op.get("requestBody")
    if not body:
        return ""
    body = deref(doc, body)
    content = body.get("content", {})
    json_ct = content.get("application/json", {})
    schema = deref(doc, json_ct.get("schema", {})) or {}
    out = ["### Request body\n"]
    props = schema.get("properties")
    if schema.get("type") == "array":
        item = deref(doc, schema.get("items", {})) or {}
        props = item.get("properties")
        out.append(f"Mảng các phần tử kiểu `{item.get('title') or 'object'}`.\n")
    if props:
        required = set(schema.get("required") or [])
        out.append("| Trường | Kiểu | Bắt buộc | Mô tả |")
        out.append("|---|---|---|---|")
        for name, spec in props.items():
            spec_d = deref(doc, spec) or {}
            desc = (spec_d.get("description") or "").replace("\n", " ")
            enum = spec_d.get("enum")
            if enum:
                desc += f" Giá trị: `{'`, `'.join(map(str, enum))}`."
            typ = schema_type_label(doc, spec)
            out.append(f"| `{name}` | {typ} | {'Có' if name in required else 'Không'} | {desc} |")
        out.append("")
    return "\n".join(out)


def render_responses(doc: dict, op: dict) -> str:
    out = ["### Responses\n", "| Mã | Mô tả |", "|---|---|"]
    for status, resp in (op.get("responses") or {}).items():
        resp = deref(doc, resp) or {}
        desc = (resp.get("description") or "").replace("\n", " ")
        out.append(f"| `{status}` | {desc} |")
    out.append("")
    return "\n".join(out)


def generate() -> str:
    doc = load_spec()
    info = doc.get("info", {})
    lines: list[str] = []

    lines.append("<!-- File này được SINH TỰ ĐỘNG từ docs/openapi.yaml bằng scripts/generate_api_docs.py. -->")
    lines.append("<!-- KHÔNG sửa tay. Sửa docs/openapi.yaml rồi chạy `make docs`. -->")
    lines.append("")
    lines.append(f"# {info.get('title', 'API')} — Tài liệu tham chiếu")
    lines.append("")
    lines.append(f"Phiên bản spec: **{info.get('version', '?')}** · OpenAPI {doc.get('openapi', '?')}")
    lines.append("")
    if info.get("description"):
        lines.append(info["description"].strip())
        lines.append("")

    servers = doc.get("servers") or []
    if servers:
        lines.append("## Servers\n")
        lines.append("| URL | Môi trường |")
        lines.append("|---|---|")
        for s in servers:
            lines.append(f"| `{s.get('url')}` | {s.get('description', '')} |")
        lines.append("")

    tags = {t["name"]: t.get("description", "") for t in doc.get("tags", []) if t.get("name")}
    by_tag: dict[str, list] = {name: [] for name in tags}
    by_tag.setdefault("_other", [])

    for route, item in (doc.get("paths") or {}).items():
        for method in HTTP_METHODS:
            op = item.get(method)
            if not op:
                continue
            op_tags = op.get("tags") or ["_other"]
            for t in op_tags:
                by_tag.setdefault(t, []).append((method, route, op))

    lines.append("## Mục lục endpoint\n")
    lines.append("| Method | Path | Tóm tắt |")
    lines.append("|---|---|---|")
    for method, route, op in sorted(
        (m, r, o) for entries in by_tag.values() for (m, r, o) in entries
    ):
        anchor = op.get("operationId", route)
        lines.append(f"| `{method.upper()}` | [`{route}`](#{anchor.lower()}) | {op.get('summary', '')} |")
    lines.append("")

    for tag_name, entries in by_tag.items():
        if not entries:
            continue
        title = tag_name if tag_name != "_other" else "Khác"
        lines.append(f"## {title}")
        if tags.get(tag_name):
            lines.append("")
            lines.append(tags[tag_name])
        lines.append("")
        for method, route, op in entries:
            op_id = op.get("operationId", "")
            lines.append(f"### `{method.upper()} {route}`")
            lines.append("")
            lines.append(f"<a id=\"{op_id.lower()}\"></a>")
            lines.append("")
            if op.get("summary"):
                lines.append(f"**{op['summary']}**")
                lines.append("")
            if op.get("description"):
                lines.append(op["description"].strip())
                lines.append("")
            params = op.get("parameters") or []
            if params:
                lines.append(render_parameters(doc, params))
            body_md = render_request_body(doc, op)
            if body_md:
                lines.append(body_md)
            lines.append(render_responses(doc, op))
            lines.append("---")
            lines.append("")

    return "\n".join(lines).rstrip() + "\n"


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--check",
        action="store_true",
        help="fail if the committed docs/api-reference.md differs from generated output",
    )
    args = parser.parse_args()

    generated = generate()

    if args.check:
        if not OUT_PATH.exists():
            print(f"FAIL: {OUT_PATH} does not exist. Run `make docs`.", file=sys.stderr)
            return 1
        committed = OUT_PATH.read_text(encoding="utf-8")
        if committed != generated:
            print(
                f"FAIL: {OUT_PATH.relative_to(REPO_ROOT)} is out of date with docs/openapi.yaml.\n"
                "Run `make docs` and commit the regenerated file.",
                file=sys.stderr,
            )
            return 1
        print(f"OK: {OUT_PATH.relative_to(REPO_ROOT)} matches docs/openapi.yaml.")
        return 0

    OUT_PATH.write_text(generated, encoding="utf-8")
    print(f"Generated {OUT_PATH.relative_to(REPO_ROOT)} from {SPEC_PATH.relative_to(REPO_ROOT)}.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
