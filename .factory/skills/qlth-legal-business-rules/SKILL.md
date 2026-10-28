---
name: qlth-legal-business-rules
description: Quy tắc nghiệp vụ pháp lý bất biến của hệ thống QLTH (định mức tiết dạy TT 05/2025, KPI & xếp loại viên chức NĐ 90/48, đánh giá học sinh TT 27/TT 89, bảo vệ dữ liệu trẻ em NĐ 13/2023). Dùng khi sửa hoặc thêm logic trong internal/core/*, hoặc khi cần đối chiếu căn cứ pháp lý cho một thay đổi nghiệp vụ.
---

# Quy tắc nghiệp vụ pháp lý — EduPrimary-KPI (QLTH)

Mọi logic trong `internal/core/` suy ra từ văn bản pháp luật. **Không được sửa các hằng số và quy tắc dưới đây** trừ khi có văn bản pháp lý mới thay thế, và khi đó phải theo quy trình thêm strategy mới (mục 4), không sửa strategy cũ.

Văn bản gốc (PDF) nằm tại `docs/phap_ly/`. Bảng tổng hợp điều khoản: `docs/phap_ly/TONG_HOP_QUY_DINH_PHAP_LY_CHI_TIET.md`.

## 1. Định mức tiết dạy — Thông tư 05/2025/TT-BGDĐT (hiệu lực 22/04/2025)

Code: `internal/core/teaching/quota.go`

| Tham số | Giá trị | Ghi chú |
|---|---|---|
| Định mức chuẩn tiểu học | **805 tiết/năm** | = 23 tiết/tuần × 35 tuần |
| PTDT bán trú / lớp khuyết tật | 21 tiết/tuần | |
| Số tuần thực dạy mặc định | **35** | Điều 5 |
| Giảm cho GVCN tiểu học | **4 tiết/tuần** | KHÔNG phải 3 (mức cũ TT 28/2009 đã hết hiệu lực) |
| Giảm cho tổ trưởng / tổ phó | 3 / 1 tiết/tuần | |
| Nuôi con dưới 12 tháng | 3 tiết/tuần | |
| Tư vấn tâm lý kiêm nhiệm | 3–8 tiết/tuần | theo quy mô |
| Dạy hè tiếng Việt trẻ DTTS trước lớp 1 | hệ số 1.5 | Điều 13 |
| Dạy hòa nhập khuyết tật | hệ số quy đổi 0.2 | |

Công thức: `overtime_hours = max(0, total_effective_hours − net_annual_quota)`;
`net_annual_quota = max(0, standard − weekly_reduction × weeks)`.

## 2. KPI & xếp loại viên chức — NĐ 90/2020, NĐ 48/2023, Luật TĐKT 2022

Code: `internal/core/kpi/kpi.go`

- Trọng số cố định: **hệ thống 40% + quản lý 30% + đồng nghiệp 15% + phụ huynh 15%**.
- **Điểm liệt** (luôn xếp `KHÔNG HOÀN THÀNH NHIỆM VỤ`, xét TRƯỚC ngưỡng điểm):
  1. Bị kỷ luật Đảng/chính quyền trong năm.
  2. Vi phạm quy định cấm dạy thêm học sinh tiểu học (TT 29/2024 Điều 4, TT 19/2026).
  3. Để xảy ra tai nạn thương tích / mất ATVSTP nghiêm trọng.
- Ngưỡng điểm: ≥90 xuất sắc (điều kiện kèm: có sáng kiến/giải pháp đổi mới được công nhận — Luật TĐKT 2022 bỏ bắt buộc SKKN nhưng yêu cầu thành tích tiêu biểu); ≥75 tốt; ≥50 hoàn thành; <50 không hoàn thành.
- **≥90 mà thiếu sáng kiến → chỉ xếp "Hoàn thành tốt"** (không được xuất sắc).
- **Trần 20% xuất sắc** trên tổng số người đạt tốt trở lên (NĐ 90/2020). Vượt trần: giữ người điểm cao nhất, hạ số còn lại xuống "tốt" kèm lý do trong `reason`. Đơn vị nhỏ (<5 người đủ điều kiện) được tối thiểu 1 suất.

## 3. Bảo vệ dữ liệu trẻ em — NĐ 13/2023 Điều 20 + Luật Căn cước 2023 + NĐ 69/2024

Code: `internal/api/handlers/handlers.go` (StudentOnboarding), `internal/core/domain/models.go`

- Hồ sơ học sinh chỉ được kích hoạt khi `has_digital_consent = true`; thiếu → trả **422**, không tạo bản ghi.
- Consent thu qua: `VNEID_SSO` (phụ huynh định danh Mức 2 đại diện cho trẻ <14 tuổi), `SMS_OTP`, hoặc `PORTAL_SIGN`. Ghi lại `consent_signed_at` + `consent_method`.
- **Không log / không trả về API** các trường: `national_id`, `guardian_national_id`, hồ sơ y tế, ảnh chân dung — trừ khi caller có quyền và mục đích nghiệp vụ rõ ràng.
- Không dùng dữ liệu học sinh cho quảng cáo hay chia sẻ bên thứ ba.

## 4. Đánh giá học sinh — TT 27/2020, chuyển tiếp TT 89/2026 (hiệu lực 01/12/2026)

Code: `internal/core/assessment/strategy.go`

- `Registry.ResolveStrategy(evalDate)` là **nơi duy nhất** quyết định strategy theo ngày:
  trước 01/12/2026 → `TT27_2020`; từ 01/12/2026 → `TT89_2026`.
- Mức đạt môn học: `HTT` / `HT` / `CHT`. Phẩm chất (5) & năng lực (3 chung + 7 đặc thù): `T` / `Đ(D)` / `C`.
- Điểm kiểm tra định kỳ: 1.0–10.0, chỉ với môn có điểm số.
- Xếp loại cuối năm: Xuất sắc (tất cả HTT + phẩm chất/năng lực T + điểm ≥9) → Tốt (≥7) → Hoàn thành (≥5) → Chưa hoàn thành (có CHT hoặc <5).
- Học sinh khuyết tật hòa nhập (TT 89): đánh giá theo kế hoạch giáo dục cá nhân, học bạ ghi "Miễn" cho môn được miễn.

**Quy trình khi có thông tư mới**: tạo struct strategy MỚI implement `AssessmentStrategy`, đăng ký vào `Registry` kèm mốc hiệu lực, thêm unit test. KHÔNG sửa strategy cũ — học bạ lịch sử đã đóng băng snapshot ký số, không bao giờ tính lại.

## 5. Checklist khi sửa code nghiệp vụ

- [ ] Đối chiếu bảng hằng số ở trên — không đổi giá trị pháp lý.
- [ ] Test cũ vẫn pass: `go test ./internal/core/...`
- [ ] Thêm test mới cho hành vi thay đổi.
- [ ] Nếu đổi response API: cập nhật `docs/openapi.yaml`, chạy `make docs`.
- [ ] Nếu hành vi vận hành đổi: cập nhật `docs/runbooks/02-xu-ly-su-co.md`.
