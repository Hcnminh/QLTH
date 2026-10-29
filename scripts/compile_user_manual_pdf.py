#!/usr/bin/env python3
"""Compile docs/huong-dan-su-dung.md to docs/huong-dan-su-dung.pdf with WeasyPrint.

Generates a professional, print-ready PDF with:
  - Elegant Cover Page with System Branding and Metadata
  - Automatic Page Numbering ("Trang X / Y") & Running Headers
  - Formatted Multi-Role Operational Guides
  - Embedded High-Resolution UI Screenshots with Borders & Captions
  - Color-Coded Legal Citations & Callout Warning Boxes
  - Clean Table Styling with Header Highlighting

Usage:
    python3 scripts/compile_user_manual_pdf.py
"""
from __future__ import annotations

import os
import subprocess
import sys
from pathlib import Path
import weasyprint

REPO_ROOT = Path(__file__).resolve().parent.parent
DOCS_DIR = REPO_ROOT / "docs"
MD_PATH = DOCS_DIR / "huong-dan-su-dung.md"
PDF_PATH = DOCS_DIR / "huong-dan-su-dung.pdf"

CSS_STYLES = """
@page {
    size: A4 portrait;
    margin: 22mm 18mm 22mm 18mm;
    @top-left {
        content: "EduPrimary-KPI: Hướng Dẫn Sử Dụng Hệ Thống";
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
        font-size: 8pt;
        color: #718096;
        border-bottom: 0.5pt solid #E2E8F0;
        padding-bottom: 4px;
    }
    @top-right {
        content: "Năm học 2026 - 2027";
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
        font-size: 8pt;
        color: #718096;
        border-bottom: 0.5pt solid #E2E8F0;
        padding-bottom: 4px;
    }
    @bottom-left {
        content: "Bảo mật nội bộ trường học • Tuân thủ NĐ 13/2023/NĐ-CP & TT 05/2025/TT-BGDĐT";
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
        font-size: 7.5pt;
        color: #A0AEC0;
        border-top: 0.5pt solid #E2E8F0;
        padding-top: 4px;
    }
    @bottom-right {
        content: "Trang " counter(page) " / " counter(pages);
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
        font-size: 8pt;
        font-weight: bold;
        color: #2D3748;
        border-top: 0.5pt solid #E2E8F0;
        padding-top: 4px;
    }
}

@page :first {
    margin: 0;
    @top-left { content: none; }
    @top-right { content: none; }
    @bottom-left { content: none; }
    @bottom-right { content: none; }
}

body {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, "Liberation Sans", sans-serif;
    font-size: 10pt;
    line-height: 1.55;
    color: #2D3748;
    background-color: #FFFFFF;
}

/* COVER PAGE */
.cover-page {
    page-break-after: always;
    height: 100vh;
    padding: 60mm 25mm 30mm 25mm;
    box-sizing: border-box;
    background: linear-gradient(135deg, #1A365D 0%, #2B6CB0 100%);
    color: #FFFFFF;
    position: relative;
}

.cover-badge {
    display: inline-block;
    padding: 6px 14px;
    background-color: rgba(255, 255, 255, 0.2);
    border: 1px solid rgba(255, 255, 255, 0.4);
    border-radius: 20px;
    font-size: 10pt;
    font-weight: 600;
    letter-spacing: 1px;
    text-transform: uppercase;
    margin-bottom: 25px;
}

.cover-title {
    font-size: 26pt;
    font-weight: 800;
    line-height: 1.25;
    margin: 0 0 15px 0;
    color: #FFFFFF;
}

.cover-subtitle {
    font-size: 13pt;
    font-weight: 400;
    line-height: 1.5;
    color: #E2E8F0;
    margin-bottom: 40px;
}

.cover-meta {
    margin-top: 80px;
    padding-top: 25px;
    border-top: 1px solid rgba(255, 255, 255, 0.3);
    font-size: 10pt;
    line-height: 1.8;
    color: #CBD5E0;
}

.cover-meta strong {
    color: #FFFFFF;
}

/* HEADINGS */
h1 {
    page-break-before: always;
    font-size: 17pt;
    font-weight: 700;
    color: #1A365D;
    border-bottom: 2pt solid #3182CE;
    padding-bottom: 6px;
    margin-top: 0;
    margin-bottom: 16px;
}

/* Don't break before first H1 after cover */
.document-content > h1:first-of-type {
    page-break-before: avoid;
}

h2 {
    font-size: 13pt;
    font-weight: 600;
    color: #2B6CB0;
    margin-top: 22px;
    margin-bottom: 10px;
    page-break-after: avoid;
}

h3 {
    font-size: 11pt;
    font-weight: 600;
    color: #2D3748;
    margin-top: 16px;
    margin-bottom: 8px;
    page-break-after: avoid;
}

h4 {
    font-size: 10pt;
    font-weight: 600;
    color: #4A5568;
    margin-top: 12px;
    margin-bottom: 6px;
    page-break-after: avoid;
}

p {
    margin-top: 0;
    margin-bottom: 10px;
    text-align: justify;
}

/* TABLES */
table {
    width: 100%;
    border-collapse: collapse;
    margin: 14px 0 18px 0;
    font-size: 8.5pt;
    line-height: 1.45;
    page-break-inside: avoid;
}

th, td {
    padding: 7px 9px;
    border: 0.5pt solid #CBD5E0;
    text-align: left;
    vertical-align: top;
}

th {
    background-color: #2B6CB0;
    color: #FFFFFF;
    font-weight: 600;
    font-size: 8.5pt;
}

tr:nth-child(even) td {
    background-color: #F7FAFC;
}

/* IMAGES & SCREENSHOTS */
img {
    max-width: 95%;
    height: auto;
    display: block;
    margin: 12px auto 6px auto;
    border: 1px solid #CBD5E0;
    border-radius: 4px;
    box-shadow: 0 2px 4px rgba(0,0,0,0.08);
    page-break-inside: avoid;
}

figure {
    margin: 14px 0;
    text-align: center;
    page-break-inside: avoid;
}

figcaption, p:has(> img) + p, img + em {
    font-size: 8.5pt;
    font-style: italic;
    color: #4A5568;
    text-align: center;
    margin-top: 4px;
    margin-bottom: 12px;
}

/* CODE BLOCKS & PRE */
pre {
    background-color: #F7FAFC;
    border: 0.5pt solid #E2E8F0;
    border-left: 3pt solid #4A5568;
    padding: 9px 12px;
    font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, Courier, monospace;
    font-size: 8pt;
    line-height: 1.4;
    overflow-x: hidden;
    white-space: pre-wrap;
    word-break: break-all;
    border-radius: 3px;
    margin: 10px 0;
    page-break-inside: avoid;
}

code {
    font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, Courier, monospace;
    font-size: 8.5pt;
    background-color: #EDF2F7;
    color: #C53030;
    padding: 1px 4px;
    border-radius: 3px;
}

pre code {
    background-color: transparent;
    color: #2D3748;
    padding: 0;
}

/* BLOCKQUOTES & CALLOUTS */
blockquote {
    margin: 12px 0;
    padding: 10px 14px;
    background-color: #EBF8FF;
    border-left: 3pt solid #3182CE;
    color: #2B6CB0;
    font-size: 9.5pt;
    border-radius: 0 4px 4px 0;
    page-break-inside: avoid;
}

/* LISTS */
ul, ol {
    margin-top: 0;
    margin-bottom: 10px;
    padding-left: 20px;
}

li {
    margin-bottom: 4px;
}

/* TOC STYLING */
#mục-lục-tổng-quan + ul, #mục-lục-tổng-quan ~ ul {
    list-style-type: none;
    padding-left: 0;
}

#mục-lục-tổng-quan + ul > li, #mục-lục-tổng-quan ~ ul > li {
    font-weight: 600;
    margin-top: 8px;
    color: #1A365D;
}

#mục-lục-tổng-quan + ul > li > ul, #mục-lục-tổng-quan ~ ul > li > ul {
    list-style-type: none;
    padding-left: 16px;
    font-weight: normal;
    color: #4A5568;
}

hr {
    border: none;
    border-top: 1pt solid #E2E8F0;
    margin: 20px 0;
}
"""


def compile_pdf() -> bool:
    print(f"Reading markdown source from {MD_PATH}...")
    if not MD_PATH.exists():
        print(f"ERROR: {MD_PATH} does not exist!", file=sys.stderr)
        return False

    # Run pandoc to convert Markdown to HTML body
    cmd = [
        "pandoc",
        str(MD_PATH),
        "-f", "gfm",
        "-t", "html5",
        "--mathml",
    ]
    print(f"Running pandoc: {' '.join(cmd)}")
    proc = subprocess.run(cmd, cwd=str(DOCS_DIR), capture_output=True, text=True)
    if proc.returncode != 0:
        print(f"Pandoc error: {proc.stderr}", file=sys.stderr)
        return False

    html_body = proc.stdout

    # Build cover page
    cover_html = """
    <div class="cover-page">
        <div class="cover-badge">Hệ Thống Quản Lý Giáo Dục Số Chuẩn Quốc Gia</div>
        <div class="cover-title">EDUPRIMARY-KPI</div>
        <div class="cover-subtitle">
            HỆ THỐNG QUẢN LÝ TRƯỜNG TIỂU HỌC TÍCH HỢP ĐÁNH GIÁ HIỆU SUẤT VIÊN CHỨC KPI<br>
            TÀI LIỆU HƯỚNG DẪN SỬ DỤNG TOÀN DIỆN ĐA VAI TRÒ
        </div>
        <div class="cover-meta">
            <strong>Phiên bản tài liệu:</strong> 2.4.0 (Năm học 2026 - 2027)<br>
            <strong>Căn cứ nghiệp vụ:</strong> Thông tư 05/2025/TT-BGDĐT • Nghị định 90/2020 & 48/2023 • Thông tư 27 & 89 • Nghị định 13/2023<br>
            <strong>Đối tượng áp dụng:</strong> BGH, GVCN, GVBM, Kế toán, Y tế học đường, Văn thư, Phụ huynh<br>
            <strong>Đơn vị ban hành:</strong> Phòng CNTT & Hội đồng Đổi mới Giáo dục Tiểu học
        </div>
    </div>
    """

    full_html = f"""<!DOCTYPE html>
<html lang="vi">
<head>
    <meta charset="utf-8">
    <title>Tài Liệu Hướng Dẫn Sử Dụng Hệ Thống EduPrimary-KPI</title>
    <style>
{CSS_STYLES}
    </style>
</head>
<body>
{cover_html}
<div class="document-content">
{html_body}
</div>
</body>
</html>
"""

    print("Rendering PDF via WeasyPrint...")
    # Base URL is DOCS_DIR so screenshots/01_login_page.png resolves cleanly
    html_doc = weasyprint.HTML(string=full_html, base_url=str(DOCS_DIR))
    html_doc.write_pdf(target=str(PDF_PATH))

    if not PDF_PATH.exists():
        print(f"ERROR: Failed to generate {PDF_PATH}!", file=sys.stderr)
        return False

    size = PDF_PATH.stat().st_size
    print(f"Successfully compiled {PDF_PATH}!")
    print(f"File size: {size} bytes ({size / (1024 * 1024):.2f} MB)")

    if size < 100 * 1024:
        print(f"ERROR: PDF size {size} is less than required 100KB!", file=sys.stderr)
        return False

    return True


if __name__ == "__main__":
    if not compile_pdf():
        sys.exit(1)
    print("User Manual PDF compilation completed successfully.")
    sys.exit(0)
