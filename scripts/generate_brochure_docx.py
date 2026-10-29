#!/usr/bin/env python3
"""Generate docs/brochure-qlth-kpi.docx - Executive Product Brochure for School Principals.

Produces a polished, professional Microsoft Word document (.docx) with:
  - Executive Cover Page with System Identity, Subtitle, and Metadata
  - Formatted Typography with Custom Palette (Navy, Royal Blue, Gold, Slate)
  - 5 Legal & Technological Breakthrough Highlights with Regulatory Citations
  - Comparative Analysis Table vs Legacy Systems (vnEdu, SMAS, MISA)
  - Detailed ROI Analysis & Unit Economics for a Standard Primary School
  - Embedded High-Resolution UI Mockups with Captions and Frames
  - SaaS Pricing Tiers (Lite, Plus, Enterprise) and Implementation Roadmap
  - Principal Executive Call to Action and Contact Information
"""
from __future__ import annotations

import os
from pathlib import Path
import docx
from docx import Document
from docx.enum.table import WD_ALIGN_VERTICAL, WD_TABLE_ALIGNMENT
from docx.enum.text import WD_ALIGN_PARAGRAPH
from docx.oxml import OxmlElement, parse_xml
from docx.oxml.ns import nsdecls, qn
from docx.shared import Inches, Pt, RGBColor

REPO_ROOT = Path(__file__).resolve().parent.parent
DOCS_DIR = REPO_ROOT / "docs"
SCREENSHOTS_DIR = DOCS_DIR / "screenshots"
OUTPUT_DOCX = DOCS_DIR / "brochure-qlth-kpi.docx"

# Palette
HEX_NAVY = "1A365D"       # Primary Dark Navy
HEX_ROYAL = "2B6CB0"      # Secondary Royal Blue
HEX_TEAL = "2C7A7B"       # Teal Accent
HEX_GOLD = "C05621"       # Dark Gold / Warm Amber
HEX_GRAY_BG = "F7FAFC"    # Soft Gray Background
HEX_GRAY_BORDER = "E2E8F0"# Border Gray
HEX_GRAY_TEXT = "4A5568"  # Slate Muted Text
HEX_DARK_TEXT = "1A202C"  # Body Dark Text
HEX_GREEN = "276749"      # Forest Green Accent

COLOR_NAVY = RGBColor(0x1A, 0x36, 0x5D)
COLOR_ROYAL = RGBColor(0x2B, 0x6C, 0xB0)
COLOR_GOLD = RGBColor(0xC0, 0x56, 0x21)
COLOR_GRAY = RGBColor(0x71, 0x80, 0x96)
COLOR_DARK = RGBColor(0x2D, 0x37, 0x48)


def set_cell_background(cell, hex_color: str) -> None:
    """Set background color of a table cell."""
    tc_pr = cell._tc.get_or_add_tcPr()
    shd = parse_xml(f'<w:shd {nsdecls("w")} w:fill="{hex_color}"/>')
    tc_pr.append(shd)


def set_cell_margins(cell, top=100, bottom=100, left=150, right=150) -> None:
    """Set inner padding of a cell in dxa (1 pt = 20 dxa)."""
    tc_pr = cell._tc.get_or_add_tcPr()
    tc_mar = parse_xml(
        f'<w:tcMar {nsdecls("w")}>'
        f'<w:top w:w="{top}" w:type="dxa"/>'
        f'<w:bottom w:w="{bottom}" w:type="dxa"/>'
        f'<w:left w:w="{left}" w:type="dxa"/>'
        f'<w:right w:w="{right}" w:type="dxa"/>'
        f'</w:tcMar>'
    )
    tc_pr.append(tc_mar)


def set_table_borders(table, hex_color="E2E8F0") -> None:
    """Apply clean light borders to a table."""
    tbl_pr = table._tbl.tblPr
    borders = parse_xml(
        f'<w:tblBorders {nsdecls("w")}>'
        f'<w:top w:val="single" w:sz="4" w:space="0" w:color="{hex_color}"/>'
        f'<w:bottom w:val="single" w:sz="4" w:space="0" w:color="{hex_color}"/>'
        f'<w:insideH w:val="single" w:sz="4" w:space="0" w:color="{hex_color}"/>'
        f'<w:insideV w:val="none"/>'
        f'<w:left w:val="none"/>'
        f'<w:right w:val="none"/>'
        f'</w:tblBorders>'
    )
    tbl_pr.append(borders)


def add_callout(doc: Document, title: str, text: str, border_hex=HEX_ROYAL, bg_hex=HEX_GRAY_BG) -> None:
    """Add a stylized callout box with a colored left bar and soft background."""
    tbl = doc.add_table(rows=1, cols=1)
    tbl.alignment = WD_TABLE_ALIGNMENT.CENTER
    tbl.autofit = False
    
    cell = tbl.cell(0, 0)
    cell.width = Inches(6.5)
    set_cell_background(cell, bg_hex)
    set_cell_margins(cell, top=140, bottom=140, left=200, right=180)
    
    # Left border only
    tc_pr = cell._tc.get_or_add_tcPr()
    borders = parse_xml(
        f'<w:tcBorders {nsdecls("w")}>'
        f'<w:top w:val="none"/>'
        f'<w:left w:val="single" w:sz="24" w:space="0" w:color="{border_hex}"/>'
        f'<w:bottom w:val="none"/>'
        f'<w:right w:val="none"/>'
        f'</w:tcBorders>'
    )
    tc_pr.append(borders)
    
    p = cell.paragraphs[0]
    p.paragraph_format.space_before = Pt(2)
    p.paragraph_format.space_after = Pt(4)
    run_t = p.add_run(f"📌 {title}\n")
    run_t.bold = True
    run_t.font.name = "Arial"
    run_t.font.size = Pt(10.5)
    run_t.font.color.rgb = COLOR_NAVY
    
    run_b = p.add_run(text)
    run_b.font.name = "Arial"
    run_b.font.size = Pt(9.5)
    run_b.font.color.rgb = COLOR_DARK

    p_after = doc.add_paragraph()
    p_after.paragraph_format.space_before = Pt(0)
    p_after.paragraph_format.space_after = Pt(4)


def add_heading_1(doc: Document, text: str) -> None:
    """Add styled Heading 1."""
    p = doc.add_paragraph()
    p.paragraph_format.space_before = Pt(16)
    p.paragraph_format.space_after = Pt(6)
    p.paragraph_format.keep_with_next = True
    run = p.add_run(text)
    run.font.name = "Arial"
    run.font.size = Pt(14)
    run.bold = True
    run.font.color.rgb = COLOR_NAVY


def add_heading_2(doc: Document, text: str) -> None:
    """Add styled Heading 2."""
    p = doc.add_paragraph()
    p.paragraph_format.space_before = Pt(12)
    p.paragraph_format.space_after = Pt(4)
    p.paragraph_format.keep_with_next = True
    run = p.add_run(text)
    run.font.name = "Arial"
    run.font.size = Pt(12)
    run.bold = True
    run.font.color.rgb = COLOR_ROYAL


def add_body_p(doc: Document, text: str, bold_prefix: str = "") -> None:
    """Add standard body paragraph."""
    p = doc.add_paragraph()
    p.paragraph_format.space_before = Pt(2)
    p.paragraph_format.space_after = Pt(4)
    p.paragraph_format.line_spacing = 1.15
    if bold_prefix:
        r_pre = p.add_run(bold_prefix)
        r_pre.font.name = "Arial"
        r_pre.font.size = Pt(10)
        r_pre.bold = True
        r_pre.font.color.rgb = COLOR_DARK
    r_body = p.add_run(text)
    r_body.font.name = "Arial"
    r_body.font.size = Pt(10)
    r_body.font.color.rgb = COLOR_DARK


def add_bullet_p(doc: Document, text: str, bold_prefix: str = "") -> None:
    """Add bullet paragraph."""
    p = doc.add_paragraph(style='List Bullet')
    p.paragraph_format.space_before = Pt(1)
    p.paragraph_format.space_after = Pt(3)
    p.paragraph_format.line_spacing = 1.15
    if bold_prefix:
        r_pre = p.add_run(bold_prefix)
        r_pre.font.name = "Arial"
        r_pre.font.size = Pt(9.5)
        r_pre.bold = True
        r_pre.font.color.rgb = COLOR_DARK
    r_body = p.add_run(text)
    r_body.font.name = "Arial"
    r_body.font.size = Pt(9.5)
    r_body.font.color.rgb = COLOR_DARK


def add_screenshot(doc: Document, img_filename: str, caption: str, width_inches=6.0) -> None:
    """Embed an annotated screenshot with caption and soft gray border."""
    img_path = SCREENSHOTS_DIR / img_filename
    if not img_path.exists():
        return

    tbl = doc.add_table(rows=2, cols=1)
    tbl.alignment = WD_TABLE_ALIGNMENT.CENTER
    tbl.autofit = False
    
    cell_img = tbl.cell(0, 0)
    cell_img.width = Inches(width_inches)
    set_cell_background(cell_img, "FAFAFA")
    set_cell_margins(cell_img, top=60, bottom=40, left=60, right=60)
    
    p_img = cell_img.paragraphs[0]
    p_img.alignment = WD_ALIGN_PARAGRAPH.CENTER
    p_img.paragraph_format.space_before = Pt(2)
    p_img.paragraph_format.space_after = Pt(2)
    run_img = p_img.add_run()
    run_img.add_picture(str(img_path), width=Inches(width_inches - 0.2))

    cell_cap = tbl.cell(1, 0)
    cell_cap.width = Inches(width_inches)
    set_cell_background(cell_cap, "EDF2F7")
    set_cell_margins(cell_cap, top=40, bottom=60, left=80, right=80)
    
    p_cap = cell_cap.paragraphs[0]
    p_cap.alignment = WD_ALIGN_PARAGRAPH.CENTER
    p_cap.paragraph_format.space_before = Pt(2)
    p_cap.paragraph_format.space_after = Pt(2)
    run_cap = p_cap.add_run(f"Hình: {caption}")
    run_cap.font.name = "Arial"
    run_cap.font.size = Pt(8.5)
    run_cap.italic = True
    run_cap.font.color.rgb = COLOR_GRAY

    # Border around table
    tbl_pr = tbl._tbl.tblPr
    borders = parse_xml(
        f'<w:tblBorders {nsdecls("w")}>'
        f'<w:top w:val="single" w:sz="6" w:space="0" w:color="CBD5E0"/>'
        f'<w:bottom w:val="single" w:sz="6" w:space="0" w:color="CBD5E0"/>'
        f'<w:left w:val="single" w:sz="6" w:space="0" w:color="CBD5E0"/>'
        f'<w:right w:val="single" w:sz="6" w:space="0" w:color="CBD5E0"/>'
        f'<w:insideH w:val="single" w:sz="4" w:space="0" w:color="E2E8F0"/>'
        f'<w:insideV w:val="none"/>'
        f'</w:tblBorders>'
    )
    tbl_pr.append(borders)

    p_after = doc.add_paragraph()
    p_after.paragraph_format.space_before = Pt(0)
    p_after.paragraph_format.space_after = Pt(6)


def build_brochure_docx() -> None:
    doc = Document()

    # Set page layout to A4 with 0.8 inch margins
    sections = doc.sections
    for section in sections:
        section.page_width = Inches(8.27)
        section.page_height = Inches(11.69)
        section.top_margin = Inches(0.8)
        section.bottom_margin = Inches(0.8)
        section.left_margin = Inches(0.8)
        section.right_margin = Inches(0.8)
        
        # Configure Header and Footer
        header = section.header
        p_hdr = header.paragraphs[0]
        p_hdr.text = "EDUPRIMARY-KPI: GIẢI PHÁP QUẢN LÝ TRƯỜNG TIỂU HỌC TOÀN DIỆN & ĐÁNH GIÁ HIỆU SUẤT"
        p_hdr.alignment = WD_ALIGN_PARAGRAPH.RIGHT
        p_hdr.runs[0].font.name = "Arial"
        p_hdr.runs[0].font.size = Pt(8)
        p_hdr.runs[0].font.color.rgb = COLOR_GRAY

        footer = section.footer
        p_ftr = footer.paragraphs[0]
        p_ftr.text = "Tài liệu giới thiệu sản phẩm dành cho Ban Giám hiệu • Hotline: 0988.123.456 • Website: qlth.edu.vn"
        p_ftr.alignment = WD_ALIGN_PARAGRAPH.CENTER
        p_ftr.runs[0].font.name = "Arial"
        p_ftr.runs[0].font.size = Pt(8)
        p_ftr.runs[0].font.color.rgb = COLOR_GRAY

    # -------------------------------------------------------------
    # 1. COVER / HERO BANNER (Executive Presentation Box)
    # -------------------------------------------------------------
    hero_tbl = doc.add_table(rows=1, cols=1)
    hero_tbl.alignment = WD_TABLE_ALIGNMENT.CENTER
    hero_tbl.autofit = False
    hero_cell = hero_tbl.cell(0, 0)
    hero_cell.width = Inches(6.67)
    set_cell_background(hero_cell, HEX_NAVY)
    set_cell_margins(hero_cell, top=300, bottom=300, left=300, right=300)
    
    p_badge = hero_cell.paragraphs[0]
    p_badge.paragraph_format.space_after = Pt(8)
    r_badge = p_badge.add_run("BỘ GIẢI PHÁP CHUYỂN ĐỔI SỐ GIÁO DỤC TIỂU HỌC 2026 - 2030")
    r_badge.font.name = "Arial"
    r_badge.font.size = Pt(9.5)
    r_badge.font.bold = True
    r_badge.font.color.rgb = RGBColor(0x90, 0xCD, 0xF4) # Light Blue
    
    p_title = hero_cell.add_paragraph()
    p_title.paragraph_format.space_after = Pt(8)
    r_title = p_title.add_run("EDUPRIMARY-KPI\nHỆ THỐNG QUẢN LÝ TRƯỜNG TIỂU HỌC TOÀN DIỆN\n& ĐÁNH GIÁ HIỆU SUẤT VIÊN CHỨC 360°")
    r_title.font.name = "Arial"
    r_title.font.size = Pt(17)
    r_title.font.bold = True
    r_title.font.color.rgb = RGBColor(0xFF, 0xFF, 0xFF)

    p_sub = hero_cell.add_paragraph()
    p_sub.paragraph_format.space_after = Pt(12)
    r_sub = p_sub.add_run("Tiên phong giải bài toán kép: Chuẩn hóa 100% nghiệp vụ sư phạm theo quy định pháp luật mới nhất & Đo lường hiệu suất, thi đua viên chức công bằng, minh bạch.")
    r_sub.font.name = "Arial"
    r_sub.font.size = Pt(10.5)
    r_sub.font.color.rgb = RGBColor(0xE2, 0xE8, 0xF0)

    p_meta = hero_cell.add_paragraph()
    p_meta.paragraph_format.space_after = Pt(0)
    r_meta = p_meta.add_run("Kính gửi: BAN GIÁM HIỆU CÁC TRƯỜNG TIỂU HỌC & LÃNH ĐẠO PHÒNG/SỞ GD&ĐT")
    r_meta.font.name = "Arial"
    r_meta.font.size = Pt(9.5)
    r_meta.font.bold = True
    r_meta.font.color.rgb = RGBColor(0xFE, 0xD7, 0xAA) # Light amber

    doc.add_paragraph().paragraph_format.space_after = Pt(8)

    # -------------------------------------------------------------
    # 2. THƯ NGỎ GỬI HIỆU TRƯỞNG & TÓM TẮT ĐIỀU HÀNH (EXECUTIVE SUMMARY)
    # -------------------------------------------------------------
    add_heading_1(doc, "1. THƯ NGỎ GỬI QUÝ HIỆU TRƯỞNG & TÓM TẮT ĐIỀU HÀNH")
    
    add_body_p(doc, 
        "Kính gửi Quý Thầy/Cô trong Ban Giám hiệu,\n\n"
        "Bước vào giai đoạn thực thi Quyết định 4725/QĐ-BGDĐT và Đề án 131/QĐ-TTg của Thủ tướng Chính phủ về tăng cường ứng dụng CNTT và chuyển đổi số trong giáo dục, các trường tiểu học trên toàn quốc đang đối diện với những áp lực cải cách thể chế chưa từng có. Việc quản lý một nhà trường tiểu học hiện đại không còn đơn thuần là vào điểm hay điểm danh, mà đòi hỏi sự kết nối chặt chẽ giữa: chuẩn hóa thời lượng giảng dạy theo Thông tư 05/2025/TT-BGDĐT, chuẩn bị chuyển đổi cách đánh giá học sinh theo Thông tư 89/2026/TT-BGDĐT, đảm bảo an toàn tuyệt đối dữ liệu trẻ em theo Nghị định 13/2023/NĐ-CP, và đặc biệt là đánh giá xếp loại viên chức công bằng, tránh cào bằng theo Nghị định 90/2020 và 48/2023/NĐ-CP.",
        "Kính thưa Ban Giám hiệu: "
    )
    
    add_body_p(doc,
        "EduPrimary-KPI được thiết kế chuyên biệt từ thực tiễn sư phạm Việt Nam, cung cấp cho Ban Giám hiệu một Trung tâm Chỉ huy Số (Executive Command Dashboard) thời gian thực. Hệ thống giúp giảm 85% thời gian tổng kết thi đua, loại bỏ hoàn toàn các khiếu nại về tiết thừa giờ, và nâng tầm chất lượng dạy học toàn trường."
    )

    add_callout(doc, 
        "Thông điệp cốt lõi dành cho Nhà quản lý Giáo dục",
        "\"EduPrimary-KPI không chỉ là một phần mềm quản lý học sinh thông thường. Đây là Trợ lý Chiến lược của Hiệu trưởng, biến toàn bộ các văn bản pháp luật phức tạp của Bộ GD&ĐT thành các luồng vận hành tự động, chính xác, giúp Ban Giám hiệu an tâm lãnh đạo và nâng cao vị thế nhà trường.\"",
        border_hex=HEX_ROYAL,
        bg_hex=HEX_GRAY_BG
    )

    # -------------------------------------------------------------
    # 3. NĂM ĐỘT PHÁ PHÁP LÝ & CÔNG NGHỆ (KEY LEGAL DIFFERENTIATORS)
    # -------------------------------------------------------------
    add_heading_1(doc, "2. NĂM ĐỘT PHÁ PHÁP LÝ & CÔNG NGHỆ ĐỘC QUYỀN")
    add_body_p(doc, "EduPrimary-KPI tạo lập sự khác biệt vượt trội so với các phần mềm hiện có nhờ việc nhúng sâu các quy tắc nghiệp vụ pháp lý bất biến ngay trong lõi hệ thống:")

    add_heading_2(doc, "Đột phá 1: Phân hệ KPI 360° Viên chức & Thuật toán Áp trần 20% Xuất sắc")
    add_bullet_p(doc, "Hệ thống tự động tổng hợp đánh giá 4 chiều theo tỷ trọng chuẩn: 40% Hệ thống (chấm tự động từ dữ liệu bài giảng, đúng giờ, hồ sơ) + 30% Ban Giám hiệu + 15% Đồng nghiệp + 15% Phụ huynh học sinh.", "• Đánh giá 4 chiều khách quan (40-30-15-15): ")
    add_bullet_p(doc, "Tự động kích hoạt điểm liệt (xếp loại Không hoàn thành nhiệm vụ) nếu giáo viên bị kỷ luật, vi phạm quy định dạy thêm tiểu học (Thông tư 29/2024/TT-BGDĐT & Thông tư 19/2026/TT-BGDĐT), hoặc để xảy ra sự cố an toàn học đường nghiêm trọng.", "• Cơ chế Điểm liệt pháp lý: ")
    add_bullet_p(doc, "Tuân thủ nghiêm ngặt Điều 2 Nghị định 48/2023/NĐ-CP: Tỷ lệ hoàn thành xuất sắc nhiệm vụ không vượt quá 20% tổng số viên chức đạt mức Tốt. Hệ thống tự động xếp hạng điểm số và hạ ngưỡng tự động đối với các trường hợp vượt trần, loại bỏ hoàn toàn hiện tượng cảm tính hoặc xung đột nội bộ.", "• Thuật toán khống chế trần 20% tự động: ")

    add_screenshot(doc, "13b_kpi_evaluations_20pct_cap.png", "Giao diện Quản lý Đánh giá KPI 360° với Huy hiệu khống chế trần 20% xuất sắc và Điểm liệt pháp lý", width_inches=6.2)

    add_heading_2(doc, "Đột phá 2: Định mức Tiết dạy & Tiết thừa giờ Chuẩn Thông tư 05/2025/TT-BGDĐT")
    add_bullet_p(doc, "Tự động áp dụng chuẩn 805 tiết/năm (35 tuần x 23 tiết/tuần) cho trường phổ thông và 21 tiết/tuần cho trường PTDT bán trú / lớp khuyết tật.", "• Chuẩn định mức 805 tiết/năm: ")
    add_bullet_p(doc, "Cập nhật chính xác quy định mới: Giáo viên chủ nhiệm (GVCN) tiểu học được giảm trừ 4 tiết/tuần (thay vì 3 tiết như Thông tư 28/2009 cũ); Tổ trưởng chuyên môn giảm 3 tiết; Tổ phó giảm 1 tiết; Nữ nuôi con dưới 12 tháng giảm 3 tiết; Kiêm nhiệm tư vấn tâm lý giảm 3–8 tiết.", "• Giảm trừ theo vị trí kiêm nhiệm: ")
    add_bullet_p(doc, "Đối soát số tiết quy đổi thực tế với định mức nghĩa vụ. Số tiết thừa giờ được tính toán tự động, không bao giờ âm, xuất biểu mẫu thanh toán thừa giờ chuẩn kho bạc.", "• Tự động tính thừa giờ: ")

    add_screenshot(doc, "03_staffs_teaching_quota.png", "Bảng Quản lý Định mức Tiết dạy và Phân công Chuyên môn chuẩn Thông tư 05/2025/TT-BGDĐT", width_inches=6.2)

    add_heading_2(doc, "Đột phá 3: Kiến trúc Đánh giá Học sinh Chuyển tiếp 01/12/2026 (TT 27/2020 & TT 89/2026)")
    add_bullet_p(doc, "Thiết kế Strategy Pattern đón đầu mốc chuyển tiếp pháp lý ngày 01/12/2026. Trước mốc này, hệ thống áp dụng Thông tư 27/2020/TT-BGDĐT; từ ngày 01/12/2026, tự động kích hoạt Thông tư 89/2026/TT-BGDĐT với yêu cầu thi trên máy tính và ngân hàng đề chung lớp 5.", "• Tự động chuyển đổi thông tư: ")
    add_bullet_p(doc, "Lịch sử đánh giá của các năm học trước được lưu trữ dưới dạng snapshot bất biến, không bao giờ bị tính lại sai lệch khi thông tư mới có hiệu lực.", "• Bảo toàn học bạ số: ")
    add_bullet_p(doc, "Tích hợp ký số từ xa SmartCA có dấu thời gian TSA chuẩn Bộ TT&TT, sẵn sàng cấp bằng và học bạ điện tử liên thông CSDL quốc gia.", "• Ký số SmartCA & TSA: ")

    add_screenshot(doc, "11_transcripts_smartca_snapshot.png", "Học bạ Điện tử và Nhật ký Ký số SmartCA với mã băm SHA-256 bảo toàn tính pháp lý", width_inches=6.2)

    add_heading_2(doc, "Đột phá 4: Bảo vệ Dữ liệu Trẻ em Tuyệt đối (Điều 20 Nghị định 13/2023/NĐ-CP)")
    add_bullet_p(doc, "Hồ sơ học sinh chỉ được kích hoạt trên hệ thống khi có xác nhận đồng thuận số (Digital Consent) từ cha/mẹ hoặc người giám hộ hợp pháp.", "• Digital Consent bắt buộc: ")
    add_bullet_p(doc, "Mã hóa dữ liệu định danh cá nhân, phân quyền truy cập nghiêm ngặt. Phụ huynh chỉ được xem dữ liệu con em mình; tuyệt đối không để lộ dữ liệu y tế, hình ảnh nhạy cảm.", "• An toàn thông tin học đường: ")

    add_heading_2(doc, "Đột phá 5: Quản lý Bán trú Dinh dưỡng, Lưu mẫu 24h & Hoàn tiền Suất ăn Tự động")
    add_bullet_p(doc, "Thiết lập thực đơn khoa học đảm bảo năng lượng Kcal và tỷ lệ cân đối P:L:G (13-20% Đạm, 20-30% Béo, 50-65% Đường bột) theo khuyến nghị Viện Dinh dưỡng Quốc gia.", "• Cân đối Kcal & P:L:G tự động: ")
    add_bullet_p(doc, "Theo dõi nguồn gốc thực phẩm, nhật ký kiểm thực 3 bước và quy trình niêm phong lưu mẫu thực phẩm 24 giờ đúng quy chuẩn Bộ Y tế.", "• Kiểm thực 3 bước & Lưu mẫu 24h: ")
    add_bullet_p(doc, "Phụ huynh báo hủy ăn qua ứng dụng trước 08:00 sáng; hệ thống tự động hoàn tiền suất ăn vào ví học sinh, minh bạch và dứt điểm mọi tranh chấp thu chi.", "• Báo hủy hoàn tiền tự động: ")

    add_screenshot(doc, "12_boarding_nutrition_menu.png", "Phân hệ Bán trú Dinh dưỡng với Bảng tính Kcal, Tỷ lệ P:L:G và Quản lý Lưu mẫu thức ăn 24h", width_inches=6.2)

    # -------------------------------------------------------------
    # 4. BẢNG SO SÁNH VỚI CÁC PHẦN MỀM TRUYỀN THỐNG
    # -------------------------------------------------------------
    add_heading_1(doc, "3. BẢNG SO SÁNH NĂNG LỰC CẠNH TRANH")
    add_body_p(doc, "Bảng đối chiếu tính năng giữa EduPrimary-KPI và các giải pháp đang có trên thị trường:")

    comp_table = doc.add_table(rows=6, cols=4)
    comp_table.alignment = WD_TABLE_ALIGNMENT.CENTER
    comp_table.autofit = False
    set_table_borders(comp_table, "CBD5E0")

    headers = [
        "Tiêu chí Nghiệp vụ",
        "vnEdu / SMAS\n(Viễn thông)",
        "MISA QLTH\n(Kế toán)",
        "EDUPRIMARY-KPI\n(Giải pháp chuyên sâu)"
    ]
    col_widths = [Inches(1.8), Inches(1.5), Inches(1.4), Inches(2.0)]
    
    hdr_row = comp_table.rows[0]
    for idx, name in enumerate(headers):
        cell = hdr_row.cells[idx]
        cell.width = col_widths[idx]
        set_cell_background(cell, HEX_NAVY)
        set_cell_margins(cell, top=100, bottom=100, left=80, right=80)
        p = cell.paragraphs[0]
        p.alignment = WD_ALIGN_PARAGRAPH.CENTER
        r = p.add_run(name)
        r.font.name = "Arial"
        r.font.size = Pt(8.5)
        r.bold = True
        r.font.color.rgb = RGBColor(0xFF, 0xFF, 0xFF)

    rows_data = [
        ("Bộ công cụ KPI Giáo viên 360°\n(Nghị định 90/2020 & 48/2023)", "❌ Không có\n(Điểm mù lớn)", "❌ Không có", "🌟 ĐỘC QUYỀN: Đánh giá 4 chiều, điểm liệt, tự động áp trần 20%"),
        ("Định mức Tiết dạy 805 tiết/năm\n(Thông tư 05/2025/TT-BGDĐT)", "⚠️ Chưa cập nhật\n(Vẫn dùng TT 28/2009 cũ)", "⚠️ Nhập liệu thủ công", "🌟 TỰ ĐỘNG: Giảm 4 tiết GVCN, trừ kiêm nhiệm, tính thừa giờ tự động"),
        ("Chuyển tiếp Đánh giá Học sinh\n(Thông tư 27/2020 & 89/2026)", "⚠️ Cập nhật chậm 6-12 tháng", "⚠️ Phụ thuộc module ngoài", "🌟 ĐÓN ĐẦU: Kiến trúc Strategy tự kích hoạt ngày 01/12/2026"),
        ("Bảo vệ Dữ liệu Trẻ em\n(Nghị định 13/2023/NĐ-CP)", "❌ Thiếu Digital Consent", "❌ Không có", "🌟 CHUẨN ĐIỀU 20: Digital Consent xác nhận mới kích hoạt hồ sơ"),
        ("Quản lý Bán trú & Hoàn tiền\n(Quy chuẩn Dinh dưỡng & Bộ Y tế)", "⚠️ Chỉ điểm danh báo ăn", "⚠️ Chỉ ghi sổ kế toán", "🌟 CÂN ĐỐI KCAL, tỷ lệ P:L:G, lưu mẫu 24h & tự động hoàn tiền"),
    ]

    for row_idx, rdata in enumerate(rows_data, start=1):
        row = comp_table.rows[row_idx]
        bg_col = HEX_GRAY_BG if row_idx % 2 == 1 else "FFFFFF"
        for col_idx, text in enumerate(rdata):
            cell = row.cells[col_idx]
            cell.width = col_widths[col_idx]
            set_cell_background(cell, bg_col if col_idx < 3 else "F0FFF4") # Light green highlight for EduPrimary
            set_cell_margins(cell, top=80, bottom=80, left=80, right=80)
            p = cell.paragraphs[0]
            p.alignment = WD_ALIGN_PARAGRAPH.LEFT if col_idx == 0 else WD_ALIGN_PARAGRAPH.CENTER
            r = p.add_run(text)
            r.font.name = "Arial"
            r.font.size = Pt(8)
            if col_idx == 3:
                r.bold = True
                r.font.color.rgb = RGBColor(0x27, 0x67, 0x49) # Green
            else:
                r.font.color.rgb = COLOR_DARK

    doc.add_paragraph().paragraph_format.space_after = Pt(8)

    # -------------------------------------------------------------
    # 5. PHÂN TÍCH HIỆU QUẢ ĐẦU TƯ (ROI ANALYSIS) & LỢI ÍCH KINH TẾ
    # -------------------------------------------------------------
    add_heading_1(doc, "4. PHÂN TÍCH HIỆU QUẢ ĐẦU TƯ (ROI ANALYSIS) & LỢI ÍCH KINH TẾ")
    add_body_p(doc, 
        "Áp dụng EduPrimary-KPI mang lại giá trị hoàn vốn đầu tư (ROI) vượt trội, có thể định lượng cụ thể cho một trường tiểu học quy mô chuẩn (30 lớp, 1.050 học sinh, 45 cán bộ giáo viên):"
    )

    roi_table = doc.add_table(rows=6, cols=3)
    roi_table.alignment = WD_TABLE_ALIGNMENT.CENTER
    roi_table.autofit = False
    set_table_borders(roi_table, "CBD5E0")

    roi_widths = [Inches(2.5), Inches(2.2), Inches(2.0)]
    roi_headers = ["Hạng mục Tiết kiệm & Tối ưu", "Trước khi dùng EduPrimary-KPI", "Khi triển khai EduPrimary-KPI"]
    
    for idx, name in enumerate(roi_headers):
        cell = roi_table.rows[0].cells[idx]
        cell.width = roi_widths[idx]
        set_cell_background(cell, HEX_NAVY)
        set_cell_margins(cell, top=100, bottom=100, left=80, right=80)
        p = cell.paragraphs[0]
        p.alignment = WD_ALIGN_PARAGRAPH.CENTER
        r = p.add_run(name)
        r.font.name = "Arial"
        r.font.size = Pt(8.5)
        r.bold = True
        r.font.color.rgb = RGBColor(0xFF, 0xFF, 0xFF)

    roi_rows = [
        ("Thời gian tổng kết thi đua BGH\n& Xếp loại viên chức", "Mất 2–3 tuần họp bình bầu, khiếu nại định mức & trần 20%", "⚡ Chỉ mất 2 giờ: Hệ thống tự động tính điểm 4 chiều & áp trần"),
        ("Thời gian lập báo cáo & sổ sách\ncủa giáo viên", "Mất ~120 giờ/năm cho việc vào sổ điểm giấy, tính tiết thừa giờ", "⚡ Giảm 85% thời gian: Dữ liệu đồng bộ tự động, xuất báo cáo 1 click"),
        ("Chi phí tin nhắn SMS thông báo\ngửi phụ huynh", "Chi 25.000 – 35.000 VNĐ/học sinh/năm (tốn 26–36 triệu VNĐ/năm)", "⚡ Miễn phí 100%: Gửi thông báo Push Notification qua ứng dụng di động"),
        ("Thất thoát & Tranh chấp tiền ăn\nbán trú học sinh", "Thường xuyên nhầm lẫn chấm ăn giấy, phụ huynh thắc mắc hoàn tiền", "⚡ Không sai lệch: Phụ huynh báo hủy trước 08:00 hoàn tiền tự động"),
        ("Rủi ro pháp lý & Phạt hành chính\n(Nghị định 13/2023 & TT 05/2025)", "Nguy cơ bị phạt vi phạm dữ liệu trẻ em hoặc tính sai tiền thừa giờ", "⚡ An tâm 100%: 100% hồ sơ có Digital Consent, kiểm soát chặt chẽ"),
    ]

    for row_idx, rdata in enumerate(roi_rows, start=1):
        row = roi_table.rows[row_idx]
        bg_col = HEX_GRAY_BG if row_idx % 2 == 1 else "FFFFFF"
        for col_idx, text in enumerate(rdata):
            cell = row.cells[col_idx]
            cell.width = roi_widths[col_idx]
            set_cell_background(cell, bg_col if col_idx < 2 else "F0FFF4")
            set_cell_margins(cell, top=80, bottom=80, left=80, right=80)
            p = cell.paragraphs[0]
            p.alignment = WD_ALIGN_PARAGRAPH.LEFT
            r = p.add_run(text)
            r.font.name = "Arial"
            r.font.size = Pt(8)
            if col_idx == 2:
                r.bold = True
                r.font.color.rgb = RGBColor(0x27, 0x67, 0x49)
            else:
                r.font.color.rgb = COLOR_DARK

    doc.add_paragraph().paragraph_format.space_after = Pt(4)

    add_callout(doc,
        "Đo lường Chỉ số Hoàn vốn Đầu tư (ROI Metric)",
        "Ước tính giá trị kinh tế trực tiếp nhà trường thu được mỗi năm: từ 65.000.000 đến 120.000.000 VNĐ (thông qua tiết kiệm chi phí in ấn, SMS, giảm giờ công quản lý và triệt tiêu sai sót tài chính). Với chi phí thuê bao SaaS chỉ từ 15.000.000 – 25.000.000 VNĐ/năm, TỶ SUẤT HOÀN VỐN ĐẠT TRÊN 350%, HOÀN VỐN TOÀN BỘ CHỈ SAU 1 HỌC KỲ.",
        border_hex=HEX_GREEN,
        bg_hex="F0FFF4"
    )

    # -------------------------------------------------------------
    # 6. MOCKUP GIAO DIỆN & TRẢI NGHIỆM ĐIỀU HÀNH THỜI GIAN THỰC
    # -------------------------------------------------------------
    add_heading_1(doc, "5. MOCKUP GIAO DIỆN & TRẢI NGHIỆM ĐIỀU HÀNH THỜI GIAN THỰC")
    add_body_p(doc, "Giao diện quản trị được thiết kế theo tiêu chuẩn công thái học sư phạm, hiện đại, trực quan, hỗ trợ đa nền tảng (Web, Máy tính bảng, Điện thoại di động):")

    add_screenshot(doc, "02_dashboard_overview.png", "Trung tâm Điều hành BGH: Thống kê số lượng học sinh, tỷ lệ chuyên cần, tiến độ giáo án và cảnh báo tức thì", width_inches=6.2)
    add_screenshot(doc, "21_integrations_sandbox_stubs.png", "Cổng Tích hợp Sandbox: Xác thực VNeID Mức 2, Ký số SmartCA/TSA, Thanh toán VietQR và Đồng bộ CSDL ngành", width_inches=6.2)

    # -------------------------------------------------------------
    # 7. CÁC GÓI DỊCH VỤ SAAS & CHÍNH SÁCH BẢO HÀNH
    # -------------------------------------------------------------
    add_heading_1(doc, "6. CÁC GÓI DỊCH VỤ SAAS & CHÍNH SÁCH BẢO HÀNH")
    add_body_p(doc, "EduPrimary-KPI áp dụng mô hình thuê bao linh hoạt theo khả năng ngân sách của từng trường:")

    price_table = doc.add_table(rows=4, cols=4)
    price_table.alignment = WD_TABLE_ALIGNMENT.CENTER
    price_table.autofit = False
    set_table_borders(price_table, "CBD5E0")

    p_widths = [Inches(1.5), Inches(1.5), Inches(2.2), Inches(1.5)]
    p_headers = ["Gói Dịch vụ", "Đối tượng", "Tính năng Trọng tâm", "Chi phí Thuê bao"]

    for idx, name in enumerate(p_headers):
        cell = price_table.rows[0].cells[idx]
        cell.width = p_widths[idx]
        set_cell_background(cell, HEX_NAVY)
        set_cell_margins(cell, top=100, bottom=100, left=80, right=80)
        p = cell.paragraphs[0]
        p.alignment = WD_ALIGN_PARAGRAPH.CENTER
        r = p.add_run(name)
        r.font.name = "Arial"
        r.font.size = Pt(8.5)
        r.bold = True
        r.font.color.rgb = RGBColor(0xFF, 0xFF, 0xFF)

    packages = [
        ("EDUPRIMARY LITE", "Trường công lập nông thôn, ngân sách cơ bản", "Quản lý học sinh, Đánh giá TT 27/TT 89, Xuất sổ điểm điện tử, Đồng bộ CSDL ngành", "8.000.000 VNĐ\n/ trường / năm"),
        ("EDUPRIMARY PLUS\n(Gói Phổ biến)", "Trường tiểu học chuẩn quốc gia, trường bán trú", "Toàn bộ gói Lite + Bán trú Kcal & Lưu mẫu 24h + Báo hủy hoàn tiền + Mobile App phụ huynh + VietQR gạch nợ tức thì", "15.000 VNĐ\n/ học sinh / tháng\n(hoặc trọn gói ngân sách)"),
        ("EDUPRIMARY ENTERPRISE\n(Gói Chữ ký Số & KPI)", "Trường chất lượng cao, trường tự chủ tài chính", "Toàn bộ gói Plus + PHÂN HỆ KPI 360° ĐỘC QUYỀN + Quản lý 805 tiết dạy TT 05/2025 + Ký số SmartCA + Executive Dashboard", "25.000.000 VNĐ\n/ trường / năm\n(Hỗ trợ kỹ thuật 24/7 VIP)")
    ]

    for row_idx, rdata in enumerate(packages, start=1):
        row = price_table.rows[row_idx]
        bg_col = "FFFFFF" if row_idx % 2 == 1 else HEX_GRAY_BG
        for col_idx, text in enumerate(rdata):
            cell = row.cells[col_idx]
            cell.width = p_widths[col_idx]
            set_cell_background(cell, "FEFCBF" if row_idx == 2 and col_idx == 0 else bg_col)
            set_cell_margins(cell, top=80, bottom=80, left=80, right=80)
            p = cell.paragraphs[0]
            p.alignment = WD_ALIGN_PARAGRAPH.LEFT if col_idx in (1, 2) else WD_ALIGN_PARAGRAPH.CENTER
            r = p.add_run(text)
            r.font.name = "Arial"
            r.font.size = Pt(8)
            if col_idx == 0:
                r.bold = True
                r.font.color.rgb = COLOR_NAVY
            elif col_idx == 3:
                r.bold = True
                r.font.color.rgb = COLOR_GOLD

    doc.add_paragraph().paragraph_format.space_after = Pt(8)

    # -------------------------------------------------------------
    # 8. LỘ TRÌNH TRIỂN KHAI 4 BƯỚC & CAM KẾT ĐỒNG HÀNH
    # -------------------------------------------------------------
    add_heading_1(doc, "7. LỘ TRÌNH TRIỂN KHAI 4 BƯỚC & THÔNG TIN LIÊN HỆ")
    
    add_bullet_p(doc, "Tiếp nhận danh sách lớp, cán bộ và quy chế thi đua riêng của nhà trường; khởi tạo hệ thống và tài khoản bảo mật trong 24 giờ.", "• Bước 1 - Thiết lập & Đồng bộ ban đầu (1–2 ngày): ")
    add_bullet_p(doc, "Tổ chức buổi hướng dẫn trực tiếp/trực tuyến cho BGH, GVCN và các tổ chuyên môn; có tài liệu và video cầm tay chỉ việc.", "• Bước 2 - Tập huấn & Hướng dẫn vận hành (2–3 ngày): ")
    add_bullet_p(doc, "Đưa hệ thống vào vận hành thực tế; đội ngũ kỹ sư đồng hành giám sát và giải đáp mọi vướng mắc 24/7.", "• Bước 3 - Vận hành song song (1 tuần): ")
    add_bullet_p(doc, "Chính thức bàn giao hệ thống; ký cam kết bảo mật thông tin và bảo trì trọn đời theo các thông tư mới của Bộ GD&ĐT.", "• Bước 4 - Nghiệm thu & Chuyển giao chính thức: ")

    contact_tbl = doc.add_table(rows=1, cols=1)
    contact_tbl.alignment = WD_TABLE_ALIGNMENT.CENTER
    contact_tbl.autofit = False
    c_cell = contact_tbl.cell(0, 0)
    c_cell.width = Inches(6.67)
    set_cell_background(c_cell, "EDF2F7")
    set_cell_margins(c_cell, top=200, bottom=200, left=250, right=250)
    
    # Border
    c_tc_pr = c_cell._tc.get_or_add_tcPr()
    borders = parse_xml(
        f'<w:tcBorders {nsdecls("w")}>'
        f'<w:top w:val="single" w:sz="16" w:space="0" w:color="{HEX_NAVY}"/>'
        f'<w:left w:val="single" w:sz="16" w:space="0" w:color="{HEX_NAVY}"/>'
        f'<w:bottom w:val="single" w:sz="16" w:space="0" w:color="{HEX_NAVY}"/>'
        f'<w:right w:val="single" w:sz="16" w:space="0" w:color="{HEX_NAVY}"/>'
        f'</w:tcBorders>'
    )
    c_tc_pr.append(borders)

    p_contact = c_cell.paragraphs[0]
    p_contact.alignment = WD_ALIGN_PARAGRAPH.CENTER
    p_contact.paragraph_format.space_after = Pt(4)
    r_ct1 = p_contact.add_run("TRUNG TÂM PHÁT TRIỂN & CHUYỂN GIAO CÔNG NGHỆ EDUPRIMARY-KPI\n")
    r_ct1.font.name = "Arial"
    r_ct1.font.size = Pt(11)
    r_ct1.bold = True
    r_ct1.font.color.rgb = COLOR_NAVY

    r_ct2 = p_contact.add_run(
        "🏢 Trụ sở chính: Tầng 8, Tòa nhà Công nghệ Giáo dục, Hà Nội / Chi nhánh TP. Hồ Chí Minh\n"
        "📞 Tổng đài tư vấn BGH: 1900.6868  |  Hotline Kỹ thuật: 0988.123.456 (24/7)\n"
        "✉️ Email tiếp nhận đăng ký: bgh@qlth.edu.vn  |  🌐 Cổng thông tin: https://qlth.edu.vn\n"
        "🎁 ƯU ĐÃI ĐẶC BIỆT: MIỄN PHÍ TRẢI NGHIỆM 30 NGÀY TOÀN BỘ PHÂN HỆ CHO CÁC TRƯỜNG ĐĂNG KÝ SỚM"
    )
    r_ct2.font.name = "Arial"
    r_ct2.font.size = Pt(9.5)
    r_ct2.font.color.rgb = COLOR_DARK

    OUTPUT_DOCX.parent.mkdir(parents=True, exist_ok=True)
    doc.save(str(OUTPUT_DOCX))
    size_bytes = OUTPUT_DOCX.stat().st_size
    print(f"Successfully generated DOCX brochure: {OUTPUT_DOCX} ({size_bytes:,} bytes, {size_bytes / 1024:.1f} KB)")


if __name__ == "__main__":
    build_brochure_docx()
