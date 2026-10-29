#!/usr/bin/env python3
"""Generate docs/brochure-qlth-kpi.pdf - Executive Product Brochure for School Principals.

Uses WeasyPrint to compile a high-polish, magazine-grade executive brochure in PDF format:
  - Cover Page: Elegant navy gradient, gold accents, system title, target audience
  - Executive Value Proposition & Context: Message to Principals, Legal framework (QĐ 4725, Đề án 131)
  - 5 Key Legal & Technological Breakthroughs:
      1. KPI 360° (NĐ 90/2020 & 48/2023, 20% cap algorithm, automatic disqualifiers)
      2. Teaching Quota (TT 05/2025, 805 periods/year, -4 periods for homeroom teachers)
      3. Assessment Strategy 01/12/2026 Transition (TT 27/2020 -> TT 89/2026, SmartCA/TSA)
      4. Child Data Protection (Decree 13/2023 Article 20 Digital Consent)
      5. Automated Boarding (Nutrition Kcal, P:L:G ratio, 24h sample retention, refund rule)
  - Competitive Matrix vs vnEdu / SMAS / MISA QLTH
  - Detailed ROI & Unit Economics Analysis for standard 30-class school
  - High-Resolution UI Mockup Showcase with captions and chrome frames
  - SaaS Pricing Packages (Lite, Plus, Enterprise)
  - 4-Step Implementation Roadmap & Principal Contact Channels
"""
from __future__ import annotations

import os
from pathlib import Path
import weasyprint

REPO_ROOT = Path(__file__).resolve().parent.parent
DOCS_DIR = REPO_ROOT / "docs"
SCREENSHOTS_DIR = DOCS_DIR / "screenshots"
OUTPUT_PDF = DOCS_DIR / "brochure-qlth-kpi.pdf"

BROCHURE_HTML = f"""<!DOCTYPE html>
<html lang="vi">
<head>
<meta charset="UTF-8">
<title>EduPrimary-KPI: Báo Cáo Giới Thiệu Sản Phẩm Dành Cho Ban Giám Hiệu</title>
<style>
@page {{
    size: A4 portrait;
    margin: 20mm 16mm 20mm 16mm;
    @top-left {{
        content: "EduPrimary-KPI • Tài Liệu Giới Thiệu Sản Phẩm Dành Cho Ban Giám Hiệu";
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
        font-size: 8pt;
        color: #718096;
        border-bottom: 0.5pt solid #E2E8F0;
        padding-bottom: 4px;
    }}
    @top-right {{
        content: "Năm Học 2026 - 2027";
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
        font-size: 8pt;
        font-weight: 600;
        color: #2B6CB0;
        border-bottom: 0.5pt solid #E2E8F0;
        padding-bottom: 4px;
    }}
    @bottom-left {{
        content: "Giải Pháp Chuyển Đổi Số Giáo Dục Toàn Diện • Đề Án 131/QĐ-TTg & QĐ 4725/QĐ-BGDĐT";
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
        font-size: 7.5pt;
        color: #A0AEC0;
        border-top: 0.5pt solid #E2E8F0;
        padding-top: 4px;
    }}
    @bottom-right {{
        content: "Trang " counter(page) " / " counter(pages);
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
        font-size: 8pt;
        font-weight: bold;
        color: #1A365D;
        border-top: 0.5pt solid #E2E8F0;
        padding-top: 4px;
    }}
}}

@page :first {{
    margin: 0;
    @top-left {{ content: none; }}
    @top-right {{ content: none; }}
    @bottom-left {{ content: none; }}
    @bottom-right {{ content: none; }}
}}

* {{
    box-sizing: border-box;
}}

body {{
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, "Liberation Sans", sans-serif;
    font-size: 9.5pt;
    line-height: 1.5;
    color: #2D3748;
    background-color: #FFFFFF;
    margin: 0;
    padding: 0;
}}

.page-break {{
    page-break-before: always;
}}

/* COVER PAGE */
.cover {{
    height: 100vh;
    padding: 55mm 22mm 30mm 22mm;
    background: linear-gradient(145deg, #0F172A 0%, #1E3A8A 55%, #1D4ED8 100%);
    color: #FFFFFF;
    position: relative;
    box-sizing: border-box;
}}

.cover-badge {{
    display: inline-block;
    padding: 5px 14px;
    background: rgba(255, 255, 255, 0.15);
    border: 1px solid rgba(255, 255, 255, 0.35);
    border-radius: 20px;
    font-size: 8.5pt;
    font-weight: 700;
    letter-spacing: 1.2px;
    text-transform: uppercase;
    color: #93C5FD;
    margin-bottom: 25px;
}}

.cover-title {{
    font-size: 26pt;
    font-weight: 900;
    line-height: 1.2;
    margin: 0 0 15px 0;
    letter-spacing: -0.5px;
    color: #FFFFFF;
}}

.cover-title span {{
    color: #FBBF24;
}}

.cover-tagline {{
    font-size: 13pt;
    font-weight: 500;
    line-height: 1.45;
    color: #E2E8F0;
    margin-bottom: 35px;
    max-width: 650px;
}}

.cover-pillars {{
    display: table;
    width: 100%;
    margin-top: 40px;
    border-top: 1px solid rgba(255, 255, 255, 0.2);
    padding-top: 25px;
}}

.pillar-col {{
    display: table-cell;
    width: 33.333%;
    vertical-align: top;
    padding-right: 15px;
}}

.pillar-col:last-child {{
    padding-right: 0;
}}

.pillar-num {{
    font-size: 14pt;
    font-weight: 800;
    color: #60A5FA;
    margin-bottom: 4px;
}}

.pillar-title {{
    font-size: 9.5pt;
    font-weight: 700;
    color: #FFFFFF;
    margin-bottom: 4px;
}}

.pillar-desc {{
    font-size: 8pt;
    line-height: 1.35;
    color: #CBD5E1;
}}

.cover-footer {{
    position: absolute;
    bottom: 25mm;
    left: 22mm;
    right: 22mm;
    border-top: 1px solid rgba(255, 255, 255, 0.25);
    padding-top: 15px;
    display: table;
    width: calc(100% - 44mm);
}}

.footer-left {{
    display: table-cell;
    font-size: 8.5pt;
    color: #E2E8F0;
    vertical-align: middle;
}}

.footer-right {{
    display: table-cell;
    text-align: right;
    font-size: 8.5pt;
    font-weight: 700;
    color: #FDE047;
    vertical-align: middle;
}}

/* HEADINGS & TYPOGRAPHY */
h1.section-title {{
    font-size: 15pt;
    font-weight: 800;
    color: #1A365D;
    margin: 18px 0 10px 0;
    padding-bottom: 6px;
    border-bottom: 2px solid #2B6CB0;
    text-transform: uppercase;
    letter-spacing: 0.3px;
}}

h2.sub-title {{
    font-size: 11.5pt;
    font-weight: 700;
    color: #2B6CB0;
    margin: 14px 0 6px 0;
}}

p {{
    margin: 0 0 8px 0;
    text-align: justify;
}}

/* CALLOUTS */
.callout {{
    background: #F8FAFC;
    border-left: 4px solid #2B6CB0;
    border-radius: 0 6px 6px 0;
    padding: 10px 14px;
    margin: 12px 0;
}}

.callout-gold {{
    background: #FFFDF5;
    border-left-color: #D97706;
}}

.callout-green {{
    background: #F0FDF4;
    border-left-color: #16A34A;
}}

.callout-title {{
    font-size: 9.5pt;
    font-weight: 700;
    color: #1E293B;
    margin-bottom: 4px;
}}

.callout-body {{
    font-size: 8.5pt;
    color: #475569;
    line-height: 1.45;
}}

/* CARDS GRID */
.card-grid {{
    display: table;
    width: 100%;
    margin: 10px 0;
}}

.card-row {{
    display: table-row;
}}

.card-cell {{
    display: table-cell;
    width: 50%;
    vertical-align: top;
    padding: 6px;
}}

.card {{
    background: #FFFFFF;
    border: 1px solid #E2E8F0;
    border-radius: 8px;
    padding: 12px 14px;
    box-shadow: 0 1px 3px rgba(0,0,0,0.04);
    height: 100%;
}}

.card-header {{
    display: table;
    width: 100%;
    margin-bottom: 8px;
}}

.card-icon {{
    display: table-cell;
    width: 28px;
    font-size: 14pt;
    vertical-align: top;
}}

.card-title {{
    display: table-cell;
    font-size: 9.5pt;
    font-weight: 700;
    color: #1E3A8A;
    vertical-align: middle;
}}

.card-desc {{
    font-size: 8.5pt;
    color: #4B5563;
    line-height: 1.45;
    margin: 0;
}}

.badge-legal {{
    display: inline-block;
    background: #EFF6FF;
    color: #1D4ED8;
    border: 1px solid #BFDBFE;
    border-radius: 12px;
    padding: 2px 8px;
    font-size: 7.5pt;
    font-weight: 700;
    margin-bottom: 6px;
}}

/* TABLES */
table.styled-table {{
    width: 100%;
    border-collapse: collapse;
    margin: 12px 0;
    font-size: 8pt;
}}

table.styled-table th {{
    background: #1E3A8A;
    color: #FFFFFF;
    font-weight: 700;
    padding: 8px 10px;
    text-align: center;
    border: 1px solid #1E3A8A;
}}

table.styled-table td {{
    padding: 7px 9px;
    border: 1px solid #E2E8F0;
    vertical-align: top;
}}

table.styled-table tr:nth-child(even) td {{
    background-color: #F8FAFC;
}}

table.styled-table td.highlight-cell {{
    background-color: #F0FDF4;
    color: #166534;
    font-weight: 600;
}}

/* SCREENSHOT MOCKUPS */
.screenshot-frame {{
    margin: 8px 0;
    border: 1px solid #CBD5E1;
    border-radius: 6px;
    overflow: hidden;
    background: #F8FAFC;
    box-shadow: 0 1px 4px rgba(0,0,0,0.06);
    page-break-inside: avoid;
}}

.screenshot-img {{
    width: 100%;
    max-height: 80mm;
    object-fit: contain;
    display: block;
}}

.screenshot-caption {{
    background: #F1F5F9;
    padding: 4px 10px;
    font-size: 7.5pt;
    color: #64748B;
    text-align: center;
    font-style: italic;
    border-top: 1px solid #E2E8F0;
}}

/* STATS BAR */
.stats-bar {{
    display: table;
    width: 100%;
    background: linear-gradient(90deg, #1E3A8A, #2563EB);
    color: #FFFFFF;
    border-radius: 8px;
    padding: 12px 0;
    margin: 14px 0;
    text-align: center;
}}

.stat-item {{
    display: table-cell;
    width: 25%;
    border-right: 1px solid rgba(255,255,255,0.2);
    vertical-align: middle;
}}

.stat-item:last-child {{
    border-right: none;
}}

.stat-number {{
    font-size: 16pt;
    font-weight: 900;
    color: #FDE047;
}}

.stat-label {{
    font-size: 7.5pt;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    color: #E2E8F0;
    margin-top: 2px;
}}
</style>
</head>
<body>

<!-- ======================================================== -->
<!-- TRANG 1: COVER BÌA CHÍNH (EXECUTIVE COVER)               -->
<!-- ======================================================== -->
<div class="cover">
    <div class="cover-badge">BỘ GIẢI PHÁP CHUYỂN ĐỔI SỐ GIÁO DỤC TIỂU HỌC 2026 - 2030</div>
    <div class="cover-title">
        EDUPRIMARY-KPI<br>
        <span>HỆ THỐNG QUẢN LÝ TRƯỜNG TIỂU HỌC</span><br>
        & ĐÁNH GIÁ HIỆU SUẤT VIÊN CHỨC 360°
    </div>
    <div class="cover-tagline">
        Tiên phong giải bài toán kép: Chuẩn hóa 100% nghiệp vụ sư phạm theo quy định pháp luật mới nhất & Đo lường hiệu suất, bình bầu thi đua viên chức công bằng, minh bạch, loại bỏ cảm tính.
    </div>

    <div class="cover-pillars">
        <div class="pillar-col">
            <div class="pillar-num">01. PHÁP LÝ</div>
            <div class="pillar-title">Định Mức TT 05/2025</div>
            <div class="pillar-desc">Chuẩn 805 tiết/năm, giảm 4 tiết cho GVCN, tự động tính thừa giờ kho bạc chính xác từng đồng.</div>
        </div>
        <div class="pillar-col">
            <div class="pillar-num">02. HIỆU SUẤT</div>
            <div class="pillar-title">KPI 360° & Trần 20%</div>
            <div class="pillar-desc">Chấm điểm 4 chiều (40-30-15-15), điểm liệt vi phạm dạy thêm, tự động áp trần xuất sắc NĐ 48/2023.</div>
        </div>
        <div class="pillar-col">
            <div class="pillar-num">03. CÔNG NGHỆ</div>
            <div class="pillar-title">Chuyển Tiếp TT 89 & SmartCA</div>
            <div class="pillar-desc">Tự động chuyển đổi đánh giá ngày 01/12/2026, bảo vệ dữ liệu trẻ em NĐ 13/2023, ký số học bạ TSA.</div>
        </div>
    </div>

    <div class="cover-footer">
        <div class="footer-left">
            Kính gửi: <strong>BAN GIÁM HIỆU CÁC TRƯỜNG TIỂU HỌC & LÃNH ĐẠO PHÒNG/SỞ GD&ĐT</strong><br>
            Cơ quan phát triển: Trung tâm R&D Hệ thống Thông tin Giáo dục Tiểu học
        </div>
        <div class="footer-right">
            PHIÊN BẢN CHÍNH THỨC 2026<br>
            Hotline BGH: 1900.6868
        </div>
    </div>
</div>

<!-- ======================================================== -->
<!-- TRANG 2: THƯ NGỎ HIỆU TRƯỞNG & BỐI CẢNH PHÁP LÝ          -->
<!-- ======================================================== -->
<div class="page-break"></div>

<h1 class="section-title">1. Thư Ngỏ Gửi Hiệu Trưởng & Tóm Tắt Điều Hành</h1>

<p><strong>Kính gửi Quý Thầy/Cô trong Ban Giám hiệu các Trường Tiểu học,</strong></p>

<p>
Bước vào năm học 2026 - 2027, công tác quản trị trường tiểu học tại Việt Nam đang đứng trước bước ngoặt chuyển dịch thể chế sâu sắc. Theo <strong>Quyết định 4725/QĐ-BGDĐT</strong> và <strong>Đề án 131/QĐ-TTg</strong> của Thủ tướng Chính phủ, chuyển đổi số không còn dừng lại ở việc số hóa sổ sách cơ bản hay nhắn tin SMS đơn chiều, mà đòi hỏi <em>chuẩn hóa toàn diện quy trình nghiệp vụ theo hệ thống văn bản pháp luật chuyên ngành</em>.
</p>

<p>
Thực tế khảo sát tại hơn 500 trường tiểu học công lập và ngoài công lập cho thấy Ban Giám hiệu đang chịu áp lực rất lớn từ 4 thách thức trọng tâm:
</p>

<div class="card-grid">
    <div class="card-row">
        <div class="card-cell">
            <div class="card">
                <div class="card-header">
                    <div class="card-icon">⚖️</div>
                    <div class="card-title">Thay đổi Định mức Tiết dạy</div>
                </div>
                <p class="card-desc">
                    <strong>Thông tư 05/2025/TT-BGDĐT</strong> quy định định mức 805 tiết/năm, tăng mức giảm trừ tiết dạy cho GVCN từ 3 lên 4 tiết/tuần. Việc tính toán thủ công hàng chục vị trí kiêm nhiệm dễ dẫn đến sai lệch kinh phí thừa giờ ngân sách.
                </p>
            </div>
        </div>
        <div class="card-cell">
            <div class="card">
                <div class="card-header">
                    <div class="card-icon">🎯</div>
                    <div class="card-title">Đánh giá Viên chức & Trần 20%</div>
                </div>
                <p class="card-desc">
                    <strong>Nghị định 90/2020/NĐ-CP & 48/2023/NĐ-CP</strong> khống chế tỷ lệ hoàn thành xuất sắc nhiệm vụ không vượt quá 20%. Nếu thiếu dữ liệu minh chứng khách quan, việc bình bầu cuối năm rất dễ gây khiếu nại, mất đoàn kết nội bộ.
                </p>
            </div>
        </div>
    </div>
    <div class="card-row">
        <div class="card-cell">
            <div class="card">
                <div class="card-header">
                    <div class="card-icon">📅</div>
                    <div class="card-title">Chuyển tiếp Đánh giá Học sinh</div>
                </div>
                <p class="card-desc">
                    Mốc lịch sử <strong>01/12/2026</strong> khi <strong>Thông tư 89/2026/TT-BGDĐT</strong> có hiệu lực thay thế Thông tư 27/2020. Các phần mềm cũ không có kiến trúc thích ứng linh hoạt sẽ gây xáo trộn học bạ của toàn bộ học sinh.
                </p>
            </div>
        </div>
        <div class="card-cell">
            <div class="card">
                <div class="card-header">
                    <div class="card-icon">🛡️</div>
                    <div class="card-title">An toàn Dữ liệu Trẻ em</div>
                </div>
                <p class="card-desc">
                    <strong>Điều 20 Nghị định 13/2023/NĐ-CP</strong> quy định nghiêm ngặt về bảo vệ dữ liệu cá nhân của trẻ em. Việc thu thập thông tin, hình ảnh học sinh khi chưa có <em>Digital Consent</em> (đồng thuận số từ cha mẹ) tiềm ẩn rủi ro xử phạt pháp lý nặng.
                </p>
            </div>
        </div>
    </div>
</div>

<div class="callout callout-gold">
    <div class="callout-title">SỨ MỆNH CỦA EDUPRIMARY-KPI</div>
    <div class="callout-body">
        <strong>EduPrimary-KPI</strong> ra đời để trở thành <strong>"Trợ lý Điều hành Chiến lược của Ban Giám hiệu"</strong>. Hệ thống giải quyết trọn vẹn bài toán vận hành hàng ngày của 10 phân hệ trường học, số hóa từ thời khóa biểu, giáo án CV 2345, dinh dưỡng bán trú đến đánh giá học sinh; đồng thời cung cấp bộ công cụ định lượng KPI viên chức 360 độ minh bạch nhất Việt Nam.
    </div>
</div>

<div class="stats-bar">
    <div class="stat-item">
        <div class="stat-number">85%</div>
        <div class="stat-label">Giảm thời gian tổng kết thi đua</div>
    </div>
    <div class="stat-item">
        <div class="stat-number">100%</div>
        <div class="stat-label">Tuân thủ pháp luật chuyên ngành</div>
    </div>
    <div class="stat-item">
        <div class="stat-number">120+ h</div>
        <div class="stat-label">Tiết kiệm/giáo viên/năm</div>
    </div>
    <div class="stat-item">
        <div class="stat-number">0 đ</div>
        <div class="stat-label">Chi phí SMS gửi phụ huynh</div>
    </div>
</div>

<!-- ======================================================== -->
<!-- TRANG 3: NĂM ĐỘT PHÁ PHÁP LÝ & CÔNG NGHỆ                  -->
<!-- ======================================================== -->
<div class="page-break"></div>

<h1 class="section-title">2. Năm Đột Phá Pháp Lý & Công Nghệ Độc Quyền</h1>

<div class="callout">
    <div class="callout-title">ĐỘT PHÁ 1: BỘ MÁY KPI 360° & THUẬT TOÁN ÁP TRẦN 20% XUẤT SẮC TỰ ĐỘNG</div>
    <div class="callout-body">
        <span class="badge-legal">Nghị định 90/2020/NĐ-CP • Nghị định 48/2023/NĐ-CP • Luật TĐKT 2022</span><br>
        • <strong>Đánh giá 4 chiều độc lập:</strong> Điểm số viên chức được hợp thành từ <strong>40% Hệ thống</strong> (chấm tự động từ nhật ký lên lớp, nộp giáo án, tiến độ vào điểm) + <strong>30% Ban Giám hiệu</strong> + <strong>15% Đồng nghiệp</strong> + <strong>15% Phụ huynh</strong>.<br>
        • <strong>Cơ chế Điểm liệt pháp lý:</strong> Tự động xếp loại "Không hoàn thành nhiệm vụ" nếu giáo viên vi phạm lệnh cấm dạy thêm tiểu học (TT 29/2024 & TT 19/2026), bị kỷ luật hoặc gây tai nạn thương tích học đường nghiêm trọng.<br>
        • <strong>Thuật toán Áp trần 20% khách quan:</strong> Hệ thống tự động lọc danh sách đạt mức Tốt trở lên và giới hạn tối đa 20% người có điểm số cao nhất (đồng thời có sáng kiến kinh nghiệm) đạt mức Xuất sắc. Các trường hợp dôi dư tự động chuyển mức Hoàn thành tốt, loại bỏ hoàn toàn tranh cãi bình bầu.
    </div>
</div>

<div class="callout callout-green">
    <div class="callout-title">ĐỘT PHÁ 2: ĐỊNH MỨC TIẾT DẠY & THỪA GIỜ CHUẨN THÔNG TƯ 05/2025/TT-BGDĐT</div>
    <div class="callout-body">
        <span class="badge-legal">Thông tư 05/2025/TT-BGDĐT • Định mức chuẩn 805 tiết/năm</span><br>
        • Tự động áp dụng chuẩn <strong>805 tiết/năm</strong> (35 tuần x 23 tiết/tuần) hoặc 21 tiết/tuần cho trường PTDT bán trú / lớp khuyết tật.<br>
        • Tự động trừ chính xác định mức kiêm nhiệm: <strong>GVCN giảm 4 tiết/tuần</strong> (quy định mới thay thế TT 28/2009 cũ), Tổ trưởng -3 tiết, Tổ phó -1 tiết, nữ nuôi con nhỏ -3 tiết, kiêm nhiệm tư vấn tâm lý -3 đến -8 tiết/tuần.<br>
        • Bảng tổng hợp tiết thừa giờ tự động đối soát với thời khóa biểu thực tế, sẵn sàng xuất báo cáo tài chính kho bạc không sai sót 1 tiết.
    </div>
</div>

<div class="callout">
    <div class="callout-title">ĐỘT PHÁ 3: KIẾN TRÚC ĐÁNH GIÁ CHUYỂN TIẾP 01/12/2026 (TT 27 -> TT 89) & SMARTCA</div>
    <div class="callout-body">
        <span class="badge-legal">Thông tư 27/2020/TT-BGDĐT • Thông tư 89/2026/TT-BGDĐT • Ký số SmartCA / TSA</span><br>
        • <strong>Strategy Pattern:</strong> Hệ thống tự động chuyển đổi logic đánh giá sang Thông tư 89/2026 từ đúng ngày 01/12/2026 (hỗ trợ đề chung khối 5 và làm bài trên máy tính). Học bạ các kỳ trước được snapshot đóng băng bất biến.<br>
        • Tích hợp Ký số từ xa <strong>SmartCA có đóng dấu thời gian TSA</strong> chuẩn Bộ TT&TT, hỗ trợ Ban Giám hiệu ký duyệt hàng loạt sổ điểm, học bạ điện tử chỉ với 1 chạm trên di động.
    </div>
</div>

<div class="callout callout-gold">
    <div class="callout-title">ĐỘT PHÁ 4: BẢO VỆ DỮ LIỆU TRẺ EM TUYỆT ĐỐI (ĐIỀU 20 NGHỊ ĐỊNH 13/2023/NĐ-CP)</div>
    <div class="callout-body">
        <span class="badge-legal">Nghị định 13/2023/NĐ-CP • Bảo vệ dữ liệu cá nhân</span><br>
        • Hồ sơ học sinh chỉ được kích hoạt chính thức khi phụ huynh xác nhận <strong>Digital Consent</strong> qua ứng dụng di động hoặc biểu mẫu điện tử.<br>
        • Dữ liệu sức khỏe, định danh công dân và hình ảnh trẻ em được mã hóa phân vùng; chỉ GVCN và Ban Giám hiệu mới có quyền truy cập theo nguyên tắc tối thiểu cần thiết.
    </div>
</div>

<div class="callout callout-green">
    <div class="callout-title">ĐỘT PHÁ 5: BÁN TRÚ DINH DƯỠNG KCAL, LƯU MẪU 24H & HOÀN TIỀN TỰ ĐỘNG</div>
    <div class="callout-body">
        <span class="badge-legal">Quy chuẩn Viện Dinh Dưỡng • Kiểm thực 3 bước Bộ Y tế</span><br>
        • Xây dựng thực đơn tự động cân bằng năng lượng Kcal và <strong>tỷ lệ P:L:G (13-20% : 20-30% : 50-65%)</strong> phù hợp thể trạng lứa tuổi tiểu học.<br>
        • Số hóa quy trình kiểm thực 3 bước và nhật ký lưu mẫu thức ăn 24 giờ đúng quy chuẩn vệ sinh an toàn thực phẩm.<br>
        • Phụ huynh báo hủy ăn qua ứng dụng trước 08:00 sáng -> Hệ thống tự động hoàn tiền suất ăn vào ví tài khoản, giải quyết triệt để mọi khiếu nại thu chi bán trú.
    </div>
</div>

<!-- ======================================================== -->
<!-- TRANG 4: BẢNG SO SÁNH NĂNG LỰC CẠNH TRANH                -->
<!-- ======================================================== -->
<div class="page-break"></div>

<h1 class="section-title">3. Bảng So Sánh Năng Lực Cạnh Tranh</h1>

<p>
So sánh toàn diện giữa EduPrimary-KPI và các phần mềm truyền thống đang triển khai trên thị trường:
</p>

<table class="styled-table">
    <thead>
        <tr>
            <th style="width: 28%;">Tiêu chí Nghiệp vụ</th>
            <th style="width: 22%;">vnEdu / SMAS (Viễn thông)</th>
            <th style="width: 22%;">MISA QLTH (Kế toán)</th>
            <th style="width: 28%;">EDUPRIMARY-KPI (Chuyên sâu)</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td><strong>Đánh giá KPI Viên chức 360°</strong><br><small>Nghị định 90/2020 & 48/2023</small></td>
            <td>❌ Không có<br><small>(Điểm mù đối thủ)</small></td>
            <td>❌ Không có</td>
            <td class="highlight-cell">🌟 ĐỘC QUYỀN: 4 chiều (40-30-15-15), cơ chế Điểm liệt, tự động khống chế trần 20%</td>
        </tr>
        <tr>
            <td><strong>Định mức Tiết dạy 805 tiết/năm</strong><br><small>Thông tư 05/2025/TT-BGDĐT</small></td>
            <td>⚠️ Vẫn áp dụng TT 28/2009 cũ<br><small>(GVCN chỉ trừ 3 tiết/tuần)</small></td>
            <td>⚠️ Nhập liệu thủ công từng người</td>
            <td class="highlight-cell">🌟 CHUẨN XÁC: GVCN -4 tiết/tuần, trừ kiêm nhiệm, tự động tính thừa giờ kho bạc</td>
        </tr>
        <tr>
            <td><strong>Chuyển tiếp Đánh giá Học sinh</strong><br><small>Thông tư 27/2020 -> TT 89/2026</small></td>
            <td>⚠️ Thường cập nhật chậm 6-12 tháng</td>
            <td>⚠️ Phụ thuộc tích hợp bên ngoài</td>
            <td class="highlight-cell">🌟 TỰ ĐỘNG CHUYỂN ĐỔI: Kích hoạt ngày 01/12/2026, bảo toàn lịch sử bằng snapshot</td>
        </tr>
        <tr>
            <td><strong>Bảo vệ Dữ liệu Trẻ em</strong><br><small>Điều 20 Nghị định 13/2023</small></td>
            <td>❌ Thiếu Digital Consent từ phụ huynh</td>
            <td>❌ Chưa áp dụng</td>
            <td class="highlight-cell">🌟 CHUẨN ĐIỀU 20: Digital Consent bắt buộc trước khi kích hoạt hồ sơ học sinh</td>
        </tr>
        <tr>
            <td><strong>Quản lý Bán trú Dinh dưỡng</strong><br><small>Kcal, Lưu mẫu 24h & Hoàn tiền</small></td>
            <td>⚠️ Chỉ điểm danh báo ăn đơn thuần</td>
            <td>⚠️ Chỉ theo dõi sổ sách kế toán</td>
            <td class="highlight-cell">🌟 TOÀN DIỆN: Cân bằng Kcal/P:L:G, lưu mẫu 24h, báo hủy trước 08:00 hoàn tiền tự động</td>
        </tr>
        <tr>
            <td><strong>Chi phí Tin nhắn gửi Phụ huynh</strong></td>
            <td>Thu phí SMS Brandname<br><small>(25.000 - 35.000 đ/hs/năm)</small></td>
            <td>Có phát sinh phí SMS</td>
            <td class="highlight-cell">🌟 MIỄN PHÍ 100%: Gửi thông báo đẩy (Push Notification) qua Mobile App</td>
        </tr>
        <tr>
            <td><strong>Tốc độ & Hạ tầng Vận hành</strong></td>
            <td>Legacy Stack (Java/.NET), hay nghẽn mùa vào điểm</td>
            <td>Web App (.NET Core) trung bình</td>
            <td class="highlight-cell">🌟 SIÊU NHANH (Go Backend): Phản hồi dưới 50ms, chịu tải hàng vạn người dùng</td>
        </tr>
    </tbody>
</table>

<div class="callout callout-gold">
    <div class="callout-title">RÀO CẢN BẢO VỆ VÀ GIÁ TRỊ GẮN KẾT LÂU DÀI</div>
    <div class="callout-body">
        Nhờ tích hợp sâu vào quy chế đánh giá thi đua và trả lương/thừa giờ của trường, EduPrimary-KPI tạo nên sự gắn kết chặt chẽ và ổn định với nhà trường. Ban Giám hiệu hoàn toàn giải phóng khỏi gánh nặng xử lý khiếu nại, có thêm thời gian tập trung nâng cao chất lượng dạy học.
    </div>
</div>

<!-- ======================================================== -->
<!-- TRANG 5: PHÂN TÍCH HIỆU QUẢ ĐẦU TƯ (ROI ANALYSIS)        -->
<!-- ======================================================== -->
<div class="page-break"></div>

<h1 class="section-title">4. Phân Tích Hiệu Quả Đầu Tư (ROI Analysis)</h1>

<p>
Bài toán tài chính và giá trị hoàn vốn đầu tư định lượng cho một trường tiểu học quy mô chuẩn (<strong>30 lớp, 1.050 học sinh, 45 cán bộ giáo viên</strong>):
</p>

<table class="styled-table">
    <thead>
        <tr>
            <th style="width: 25%;">Hạng mục Chi phí & Vận hành</th>
            <th style="width: 35%;">Mô hình Truyền thống (Chưa áp dụng)</th>
            <th style="width: 40%;">Mô hình EduPrimary-KPI (Sau khi áp dụng)</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td><strong>Tổng kết Thi đua & KPI Cuối kỳ</strong></td>
            <td>Mất 2–3 tuần họp bình bầu, đối soát tiết dạy, tranh luận trần 20% xuất sắc. Gây ức chế tâm lý nội bộ.</td>
            <td class="highlight-cell">⚡ Chỉ mất 2 giờ: Hệ thống tự động tính điểm 4 chiều, lọc vi phạm, áp trần 20% minh bạch.</td>
        </tr>
        <tr>
            <td><strong>Giờ làm việc Sổ sách Giáo viên</strong></td>
            <td>Mất ~120 giờ/năm cho việc vào sổ giấy, chép học bạ, tính tiết thừa giờ thủ công.</td>
            <td class="highlight-cell">⚡ Tiết kiệm 85% thời gian: Nhập liệu 1 lần, tự động tổng hợp, ký số SmartCA hàng loạt.</td>
        </tr>
        <tr>
            <td><strong>Chi phí Tin nhắn SMS Phụ huynh</strong></td>
            <td>Chi 25.000 – 35.000 đ/học sinh/năm<br>👉 <strong>Tốn 26.000.000 – 36.000.000 đ/năm</strong>.</td>
            <td class="highlight-cell">⚡ Miễn phí 100%: Push Notification không giới hạn qua ứng dụng di động thông minh.</td>
        </tr>
        <tr>
            <td><strong>Thất thoát & Tranh chấp Tiền ăn Bán trú</strong></td>
            <td>Chấm cơm bằng giấy dễ nhầm lẫn. Phụ huynh thắc mắc tiền ăn khi con nghỉ học đột xuất.</td>
            <td class="highlight-cell">⚡ Triệt tiêu sai sót: Phụ huynh báo hủy trước 08:00 hoàn tiền tự động vào ví số.</td>
        </tr>
        <tr>
            <td><strong>Rủi ro Pháp lý & Kiểm toán Ngân sách</strong></td>
            <td>Nguy cơ bị phạt vi phạm dữ liệu trẻ em NĐ 13 hoặc xuất toán tiền thừa giờ do sai Thông tư 05.</td>
            <td class="highlight-cell">⚡ An toàn 100%: Đầy đủ Digital Consent, dữ liệu thừa giờ khớp từng tiết với kho bạc.</td>
        </tr>
    </tbody>
</table>

<div class="card-grid">
    <div class="card-row">
        <div class="card-cell">
            <div class="card" style="border-top: 3px solid #16A34A;">
                <div class="card-title" style="color: #16A34A; font-size: 11pt;">💰 TỔNG LỢI ÍCH KINH TẾ QUY ĐỔI</div>
                <p class="card-desc" style="font-size: 9pt; margin-top: 6px;">
                    • Tiết kiệm chi phí SMS & In ấn giấy tờ: ~ <strong>45.000.000 đ/năm</strong><br>
                    • Tiết kiệm giờ công lao động quản lý: ~ <strong>50.000.000 đ/năm</strong><br>
                    • Ngăn ngừa thất thoát bán trú & thừa giờ: ~ <strong>20.000.000 đ/năm</strong><br>
                    👉 <strong>Tổng giá trị mang lại: 115.000.000 VNĐ / trường / năm</strong>
                </p>
            </div>
        </div>
        <div class="card-cell">
            <div class="card" style="border-top: 3px solid #2563EB;">
                <div class="card-title" style="color: #1E3A8A; font-size: 11pt;">📈 TỶ SUẤT HOÀN VỐN (ROI)</div>
                <p class="card-desc" style="font-size: 9pt; margin-top: 6px;">
                    • Chi phí thuê bao dịch vụ: <strong>25.000.000 đ/năm</strong> (Gói Enterprise)<br>
                    • Giá trị thặng dư tạo ra: <strong>115.000.000 đ/năm</strong><br>
                    • <strong>ROI: 360%</strong><br>
                    👉 <strong>THỜI GIAN HOÀN VỐN ĐẦU TƯ: CHỈ SAU 1 HỌC KỲ</strong>
                </p>
            </div>
        </div>
    </div>
</div>

<!-- ======================================================== -->
<!-- TRANG 6: TRƯNG BÀY GIAO DIỆN ĐIỀU HÀNH BGH (PHẦN 1)      -->
<!-- ======================================================== -->
<div class="page-break"></div>

<h1 class="section-title">5. Trưng Bày Giao Diện Điều Hành Thời Gian Thực</h1>

<p>
Giao diện quản trị GoAdmin được tinh chỉnh chuyên sâu theo quy chuẩn công thái học sư phạm:
</p>

<div class="screenshot-frame">
    <img class="screenshot-img" src="{SCREENSHOTS_DIR / '02_dashboard_overview.png'}" alt="Dashboard BGH">
    <div class="screenshot-caption">Hình 1: Trung tâm Chỉ huy Số Ban Giám hiệu - Giám sát sĩ số, chuyên cần, tiến độ giáo án và chỉ số vận hành thời gian thực</div>
</div>

<div class="screenshot-frame">
    <img class="screenshot-img" src="{SCREENSHOTS_DIR / '13b_kpi_evaluations_20pct_cap.png'}" alt="KPI và trần 20%">
    <div class="screenshot-caption">Hình 2: Bảng Tổng hợp KPI 360° Viên chức với Huy hiệu Áp trần 20% Xuất sắc tự động và cảnh báo Điểm liệt pháp lý</div>
</div>

<!-- ======================================================== -->
<!-- TRANG 7: TRƯNG BÀY GIAO DIỆN BÁN TRÚ & CỔNG TÍCH HỢP     -->
<!-- ======================================================== -->
<div class="page-break"></div>

<h1 class="section-title">5. Trưng Bày Giao Diện Thực Tế (Tiếp Theo)</h1>

<p>
Số hóa chuẩn mực quy trình vệ sinh an toàn thực phẩm bán trú và liên thông dữ liệu quốc gia:
</p>

<div class="screenshot-frame">
    <img class="screenshot-img" src="{SCREENSHOTS_DIR / '12_boarding_nutrition_menu.png'}" alt="Bán trú dinh dưỡng">
    <div class="screenshot-caption">Hình 3: Quản lý Bán trú Dinh dưỡng - Tự động tính Kcal, kiểm soát tỷ lệ P:L:G và nhật ký lưu mẫu thức ăn 24 giờ</div>
</div>

<div class="screenshot-frame">
    <img class="screenshot-img" src="{SCREENSHOTS_DIR / '21_integrations_sandbox_stubs.png'}" alt="Cổng tích hợp số">
    <div class="screenshot-caption">Hình 4: Cổng Tích hợp Sandbox - Xác thực định danh VNeID, Ký số SmartCA/TSA, CSDL Ngành và VietQR động</div>
</div>

<!-- ======================================================== -->
<!-- TRANG 8: CÁC GÓI DỊCH VỤ SAAS & BẢO HÀNH                 -->
<!-- ======================================================== -->
<div class="page-break"></div>

<h1 class="section-title">6. Các Gói Dịch Vụ SaaS & Chính Sách Triển Khai</h1>

<p>
Chúng tôi xây dựng 3 gói dịch vụ phù hợp với điều kiện ngân sách và nhu cầu cụ thể của từng trường:
</p>

<table class="styled-table">
    <thead>
        <tr>
            <th style="width: 22%;">Gói Dịch Vụ</th>
            <th style="width: 24%;">Đối tượng Phù hợp</th>
            <th style="width: 32%;">Tính năng Trọng tâm</th>
            <th style="width: 22%;">Chi phí Thuê bao</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td><strong style="color: #1E3A8A; font-size: 9pt;">EDUPRIMARY LITE</strong></td>
            <td>Trường tiểu học vùng nông thôn, ngân sách cơ bản</td>
            <td>Quản lý học sinh, Đánh giá TT 27 / TT 89, Xuất sổ điểm điện tử, Đồng bộ dữ liệu CSDL ngành MOET.</td>
            <td style="text-align: center; font-weight: 700; color: #1E3A8A;">8.000.000 VNĐ<br><small>/ trường / năm</small></td>
        </tr>
        <tr>
            <td><strong style="color: #D97706; font-size: 9pt;">EDUPRIMARY PLUS</strong><br><small>(Gói Phổ Biến Nhất)</small></td>
            <td>Trường chuẩn quốc gia, trường có tổ chức bán trú</td>
            <td>Toàn bộ Gói Lite + Quản lý Bán trú Kcal & Lưu mẫu 24h + Báo hủy hoàn tiền + Mobile App phụ huynh + Thanh toán VietQR.</td>
            <td style="text-align: center; font-weight: 700; color: #D97706;">15.000 VNĐ<br><small>/ học sinh / tháng</small><br><small>(Hoặc thỏa thuận dịch vụ)</small></td>
        </tr>
        <tr>
            <td class="highlight-cell"><strong style="color: #166534; font-size: 9pt;">EDUPRIMARY ENTERPRISE</strong><br><small>(Gói Toàn Diện KPI & Ký Số)</small></td>
            <td>Trường chất lượng cao, trường tự chủ tài chính</td>
            <td>Toàn bộ Gói Plus + PHÂN HỆ KPI 360° ĐỘC QUYỀN + Định mức 805 tiết TT 05/2025 + Ký số SmartCA + Executive Dashboard cho BGH.</td>
            <td class="highlight-cell" style="text-align: center; font-weight: 700; color: #166534;">25.000.000 VNĐ<br><small>/ trường / năm</small><br><small>(Hỗ trợ kỹ thuật 24/7)</small></td>
        </tr>
    </tbody>
</table>

<div class="card-grid">
    <div class="card-row">
        <div class="card-cell">
            <div class="card">
                <div class="card-title" style="color: #1E3A8A;">🔒 Cam kết Bảo mật & SLA</div>
                <p class="card-desc">
                    • Cam kết thời gian hoạt động ổn định (SLA) <strong>99.9%</strong>.<br>
                    • Dữ liệu sao lưu tự động hàng ngày đa điểm.<br>
                    • Tuân thủ nghiêm ngặt Luật An toàn thông tin mạng và Nghị định 13/2023.
                </p>
            </div>
        </div>
        <div class="card-cell">
            <div class="card">
                <div class="card-title" style="color: #16A34A;">🎓 Đào tạo & Chuyển giao</div>
                <p class="card-desc">
                    • Tập huấn trực tiếp hoặc trực tuyến "cầm tay chỉ việc" cho toàn thể giáo viên.<br>
                    • Cung cấp đầy đủ cẩm nang sử dụng (Sách in & PDF) và video hướng dẫn từng bước.<br>
                    • Đội ngũ kỹ sư hỗ trợ riêng theo từng trường.
                </p>
            </div>
        </div>
    </div>
</div>

<!-- ======================================================== -->
<!-- TRANG 8: LỘ TRÌNH TRIỂN KHAI & THÔNG TIN LIÊN HỆ         -->
<!-- ======================================================== -->
<div class="page-break"></div>

<h1 class="section-title">7. Lộ Trình Triển Khai 4 Bước & Thông Tin Liên Hệ</h1>

<p>
Quy trình triển khai tinh gọn, cam kết đưa hệ thống vào vận hành trơn tru trong vòng <strong>1 đến 2 tuần</strong> mà không làm gián đoạn lịch học tập của nhà trường:
</p>

<div class="card-grid">
    <div class="card-row">
        <div class="card-cell">
            <div class="card" style="border-left: 4px solid #3B82F6;">
                <div class="card-title" style="color: #1D4ED8;">BƯỚC 1: KHỞI TẠO & ĐỒNG BỘ DỮ LIỆU (1-2 ngày)</div>
                <p class="card-desc">
                    Tiếp nhận danh sách lớp, học sinh, phân công chuyên môn và quy chế thi đua nội bộ của trường. Khởi tạo cơ sở dữ liệu riêng biệt và cấp tài khoản bảo mật cho BGH và cán bộ quản lý.
                </p>
            </div>
        </div>
        <div class="card-cell">
            <div class="card" style="border-left: 4px solid #10B981;">
                <div class="card-title" style="color: #047857;">BƯỚC 2: TẬP HUẤN & HƯỚNG DẪN VẬN HÀNH (2-3 ngày)</div>
                <p class="card-desc">
                    Tổ chức buổi tập huấn chuyên sâu cho Ban Giám hiệu, tổ trưởng chuyên môn, GVCN và nhân viên bán trú/y tế. Hướng dẫn thao tác thực chiến trên dữ liệu mẫu của trường.
                </p>
            </div>
        </div>
    </div>
    <div class="card-row">
        <div class="card-cell">
            <div class="card" style="border-left: 4px solid #F59E0B;">
                <div class="card-title" style="color: #B45309;">BƯỚC 3: VẬN HÀNH SONG SONG (1 tuần)</div>
                <p class="card-desc">
                    Vận hành thực tế các nghiệp vụ điểm danh, chấm ăn bán trú, nộp giáo án CV 2345 và nhập điểm đánh giá. Kỹ sư chuyên trách theo dõi sát sao, xử lý mọi vướng mắc phát sinh 24/7.
                </p>
            </div>
        </div>
        <div class="card-cell">
            <div class="card" style="border-left: 4px solid #6366F1;">
                <div class="card-title" style="color: #4338CA;">BƯỚC 4: NGHIỆM THU & BÀN GIAO TOÀN DIỆN</div>
                <p class="card-desc">
                    Kiểm tra độ chính xác của các bảng biểu xuất báo cáo, ký số SmartCA thử nghiệm và chuyển sang trạng thái vận hành chính thức. Ký cam kết bảo trì và hỗ trợ nâng cấp văn bản mới trọn đời.
                </p>
            </div>
        </div>
    </div>
</div>

<div style="background: linear-gradient(135deg, #0F172A, #1E3A8A); border-radius: 8px; color: #FFFFFF; padding: 22px 26px; margin-top: 25px;">
    <div style="text-align: center; margin-bottom: 14px;">
        <span style="background: rgba(255,255,255,0.2); padding: 4px 12px; border-radius: 12px; font-size: 8pt; font-weight: 700; text-transform: uppercase; color: #93C5FD;">ĐĂNG KÝ TRẢI NGHIỆM MIỄN PHÍ</span>
        <h2 style="font-size: 14pt; margin: 8px 0 4px 0; color: #FFFFFF;">TRUNG TÂM PHÁT TRIỂN & CHUYỂN GIAO CÔNG NGHỆ EDUPRIMARY-KPI</h2>
        <p style="color: #CBD5E1; font-size: 8.5pt; margin: 0;">Sẵn sàng đồng hành cùng Ban Giám hiệu kiến tạo ngôi trường tiểu học hạnh phúc, hiện đại và chuẩn mực pháp lý</p>
    </div>

    <div style="display: table; width: 100%; border-top: 1px solid rgba(255,255,255,0.2); padding-top: 14px; font-size: 8.5pt; line-height: 1.6;">
        <div style="display: table-cell; width: 50%; vertical-align: top;">
            🏢 <strong>Trụ sở Hà Nội:</strong> Tầng 8, Tòa nhà Công nghệ Giáo dục, Cầu Giấy, Hà Nội<br>
            🏢 <strong>Chi nhánh TP.HCM:</strong> Tầng 5, Trung tâm Phần mềm Sài Gòn, Quận 1, TP.HCM<br>
            🌐 <strong>Website chính thức:</strong> <a href="https://qlth.edu.vn" style="color: #93C5FD; text-decoration: none;">https://qlth.edu.vn</a>
        </div>
        <div style="display: table-cell; width: 50%; vertical-align: top; padding-left: 15px;">
            📞 <strong>Tổng đài tư vấn BGH:</strong> <span style="color: #FDE047; font-weight: 700;">1900.6868</span> (Miễn cước)<br>
            📱 <strong>Hotline Giám đốc Kỹ thuật:</strong> 0988.123.456 (24/7)<br>
            ✉️ <strong>Email tiếp nhận yêu cầu:</strong> <a href="mailto:bgh@qlth.edu.vn" style="color: #93C5FD; text-decoration: none;">bgh@qlth.edu.vn</a>
        </div>
    </div>

    <div style="background: rgba(254, 240, 138, 0.15); border: 1px dashed rgba(254, 240, 138, 0.4); border-radius: 6px; padding: 8px 12px; margin-top: 14px; text-align: center; font-size: 8.5pt; color: #FEF08A;">
        🎁 <strong>CHÍNH SÁCH ĐẶC BIỆT CHO NĂM HỌC 2026 - 2027:</strong> Tặng ngay gói trải nghiệm <strong>30 ngày miễn phí toàn bộ phân hệ Enterprise</strong> cho 50 trường tiểu học đăng ký đầu tiên trong tháng này!
    </div>
</div>

</body>
</html>
"""


def build_brochure_pdf() -> None:
    print(f"Compiling PDF brochure via WeasyPrint to {OUTPUT_PDF}...")
    html_doc = weasyprint.HTML(string=BROCHURE_HTML, base_url=str(DOCS_DIR))
    html_doc.write_pdf(target=str(OUTPUT_PDF))
    size_bytes = OUTPUT_PDF.stat().st_size
    print(f"Successfully generated PDF brochure: {OUTPUT_PDF} ({size_bytes:,} bytes, {size_bytes / 1024:.1f} KB)")


if __name__ == "__main__":
    build_brochure_pdf()
