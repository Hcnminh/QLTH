# Nhiệm vụ: Hoàn thiện Hệ thống Phần mềm Quản lý Trường Tiểu học tích hợp KPI (EduPrimary-KPI)

## 1. Tổng quan kế hoạch (Plan Overview)
Xây dựng và hoàn thiện nền tảng backend Go Core cho Phần mềm Quản lý Trường Tiểu học (EduPrimary-KPI), tuân thủ 100% các quy định pháp luật hiện hành của Bộ GD&ĐT, chuẩn hóa kiểm thử tương tác end-to-end, tích hợp cơ chế bảo vệ dữ liệu trẻ em theo Nghị định 13/2023/NĐ-CP và sẵn sàng triển khai trên hạ tầng VPS1 (Podman/Traefik).

- **Stack công nghệ**: Go 1.27+ (Gin Web Framework), PostgreSQL 16, Redis 7.
- **Ranh giới pháp lý**:
  - Thông tư 05/2025/TT-BGDĐT: Chuẩn định mức 805 tiết/năm, GVCN tiểu học giảm 4 tiết/tuần.
  - Thông tư 27/2020/TT-BGDĐT & Thông tư 89/2026/TT-BGDĐT (hiệu lực 01/12/2026): Đánh giá học sinh tiểu học, thi máy tính, đề chung Sở/Phòng lớp 5.
  - Nghị định 90/2020/NĐ-CP & Nghị định 48/2023/NĐ-CP: Đánh giá xếp loại viên chức, trần 20% xuất sắc.
  - Nghị định 13/2023/NĐ-CP Điều 20: Digital Consent bảo vệ dữ liệu trẻ em.
  - Thông tư 29/2024/TT-BGDĐT & Thông tư 19/2026/TT-BGDĐT: Kiểm soát nghiêm cấm dạy thêm tiểu học.

---

## 2. Các cột mốc và tính năng dự kiến (Milestones & Features)

### Cột mốc 1: Quản trị Dữ liệu & Nghiệp vụ Sư phạm Cốt lõi (Core Domain & Teaching Quota)
- **F-CORE-01: Định mức tiết dạy & Quyết toán thừa giờ (TT 05/2025)**
  - Tự động hóa tính toán: 805 tiết chuẩn, giảm 4 tiết/tuần cho GVCN, trừ tiết kiêm nhiệm (tổ trưởng, tư vấn tâm lý, con nhỏ <12 tháng).
  - Quy đổi tiết dạy tiếng Việt hè cho trẻ DTTS (1.5x) và dạy hòa nhập khuyết tật (0.2x).
  - Xuất bảng kê thừa giờ không bao giờ âm.
- **F-CORE-02: Đánh giá học sinh Strategy Versioning (TT 27/2020 -> TT 89/2026)**
  - Tự động chuyển đổi strategy theo mốc hiệu lực 01/12/2026.
  - Đánh giá thường xuyên HTT/HT/CHT, 5 phẩm chất, 10 năng lực.
  - Xếp loại cuối năm 4 mức, hỗ trợ học sinh khuyết tật hòa nhập (kế hoạch IEP).
  - Snapshot bất biến phục vụ ký số học bạ điện tử (SmartCA/Cloud CA).

### Cột mốc 2: Bộ công cụ Đánh giá Hiệu suất KPI 360 độ (KPI Engine)
- **F-KPI-01: Thuật toán tính điểm 4 chiều có trọng số**
  - Hệ thống tự động 40% + Quản lý 30% + Đồng nghiệp 15% + Phụ huynh CSAT 15%.
  - Tự động phát hiện điểm liệt (kỷ luật NĐ 48/2023, vi phạm dạy thêm TT 29/2024, tai nạn nghiêm trọng) -> xếp Không hoàn thành nhiệm vụ.
- **F-KPI-02: Xếp loại thi đua toàn trường & Áp trần 20% xuất sắc**
  - Ràng buộc sáng kiến/giải pháp đổi mới sáng tạo (Luật TĐKT 2022) cho mức xuất sắc.
  - Thuật toán khống chế trần 20% trên tổng số đạt tốt trở lên theo Nghị định 90/2020.

### Cột mốc 3: Bán trú, Dinh dưỡng & Tiếp nhận Học sinh an toàn
- **F-SVC-01: Bán trú thông minh & Báo hủy ăn tự động bù trừ tiền**
  - Báo hủy ăn trước 08:00 sáng -> tự động hoàn tiền/bù trừ hóa đơn tháng sau.
  - Báo hủy sau 08:00 sáng -> từ chối hoàn tiền theo quy chế bếp.
- **F-SVC-02: Tiếp nhận học sinh & Digital Consent (NĐ 13/2023 Điều 20)**
  - Kiểm tra bắt buộc sự đồng ý số của cha mẹ/người giám hộ trước khi kích hoạt hồ sơ.
  - Từ chối với mã HTTP 422 nếu thiếu Digital Consent.

### Cột mốc 4: Tích hợp Giao diện Quản trị GoAdmin & Tự động hóa Kiểm định
- **F-ADM-01: GoAdmin Back-Office nhúng trong Go binary**
  - Tích hợp `github.com/GoAdminGroup/go-admin` quản lý bảng biểu, danh mục dùng chung.
- **F-QA-01: Bộ kiểm thử tương tác end-to-end (QA Harness)**
  - Duy trì kịch bản `scripts/qa_e2e.sh` chạy 26+ ca kiểm thử tự động, route logs ra disk, tự giải phóng cổng.

---

## 3. Tiêu chí nghiệm thu & Ràng buộc thẩm định (Validation Contract)
1. **Kiểm thử tự động**: Toàn bộ unit tests và integration tests phải vượt qua (`make test`, coverage >= 70%).
2. **Kiểm thử tương tác**: `make qa` (scripts/qa_e2e.sh) chạy thành công 26/26 ca, log máy chủ ghi ra `logs/qa-server.log`.
3. **Tuân thủ quy định**: Mọi thay đổi logic phải giữ nguyên các hằng số pháp lý quy định trong `.factory/skills/qlth-legal-business-rules/SKILL.md`.
4. **Tài liệu & OpenAPI**: `make docs` xác thực OpenAPI 3.0.3 không có lỗi cấu trúc hoặc lệch route.
