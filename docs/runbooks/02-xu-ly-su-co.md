# Runbook 02 — Xử lý sự cố (Incident Response)

Quy trình phản ứng khi hệ thống EduPrimary-KPI gặp sự cố trong giờ vận hành.
Mức độ ưu tiên: **P0** (dừng dịch vụ / sai kết quả pháp lý / lộ dữ liệu trẻ em) → **P3** (lỗi nhỏ).

## 1. Phân loại mức độ sự cố

| Mức | Định nghĩa | Ví dụ | Thời gian phản hồi |
|---|---|---|---|
| **P0** | Service chết hoàn toàn, hoặc trả kết quả pháp lý sai, hoặc rò rỉ dữ liệu trẻ em | API 5xx toàn bộ; KPI xếp loại sai mức; endpoint trả số định danh học sinh ra log | Ngay lập tức |
| **P1** | Một phân hệ nghiệp vụ lỗi, ảnh hưởng nhiều trường | Tính định mức tiết dạy sai; không ký được học bạ | Trong 2 giờ |
| **P2** | Lỗi có workaround | Một endpoint trả 400 với payload hợp lệ | Trong 1 ngày |
| **P3** | cosmetic, log noise | Sai message tiếng Việt | Theo lịch |

## 2. Quy trình xử lý chung

1. **Xác nhận & phân loại** — kiểm tra `GET /health`, đọc log (`docker compose logs -f api`), đối chiếu bảng trên.
2. **Chặn ảnh hưởng lan rộng**:
   - Nếu kết quả **sai về pháp lý** (định mức, KPI, xếp loại học sinh): dừng endpoint liên quan hoặc đưa hệ thống về chế độ chỉ-đọc; thông báo BGH các trường đang dùng dừng nhập liệu.
   - Nếu **rò rỉ dữ liệu trẻ em**: thu hồi token, tắt endpoint, giữ nguyên hiện trường log để điều tra theo Nghị định 13/2023/NĐ-CP (Điều 23 — thông báo vi phạm trong 72 giờ).
3. **Tái hiện bằng test**: viết một test fail trong package tương ứng (`internal/core/...` hoặc `internal/api/router_test.go`) trước khi sửa.
4. **Sửa & verify**: `make lint && make test && make build`. Với lỗi logic pháp lý, đối chiếu văn bản gốc trong `docs/phap_ly/`.
5. **Rollback nếu cần**: xem mục 5.
6. **Ghi nhận**: tạo issue theo template `bug_report.yml`, gắn nhãn `priority/Px` và `type/bug`.

## 3. Sự cố thường gặp và cách xử lý

### 3.1 API không khởi động / cổng 8080 bị chiếm

```bash
lsof -ti :8080            # tìm PID đang giữ cổng
make run                  # khởi động lại; nếu vẫn lỗi, đổi PORT=8081 make run
```

### 3.2 Kết quả định mức tiết dạy sai (TT 05/2025)

Đối chiếu `internal/core/teaching/quota.go` với bảng chuẩn:

| Trường hợp | Kỳ vọng |
|---|---|
| Tiểu học thường | chuẩn 805 tiết/năm |
| PTDT bán trú / lớp khuyết tật | 21 tiết/tuần |
| GVCN tiểu học | **giảm 4 tiết/tuần** (KHÔNG phải 3) |
| Tổ trưởng / tổ phó | giảm 3 / 1 tiết/tuần |
| Nuôi con < 12 tháng | giảm 3 tiết/tuần |

Test chuẩn để chạy lại: `go test -run TestCalculateQuota ./internal/core/teaching/...`

### 3.3 Xếp loại KPI / viên chức sai (NĐ 90/2020, NĐ 48/2023)

Kiểm tra theo thứ tự trong `internal/core/kpi/kpi.go`:
1. Điểm liệt có được xét **trước** ngưỡng điểm không (kỷ luật, vi phạm dạy thêm tiểu học, tai nạn nghiêm trọng → luôn "Không hoàn thành").
2. Trọng số có đúng 40/30/15/15 không.
3. Mức ≥90 có yêu cầu `has_innovative_solution` không (thiếu → chỉ "Hoàn thành tốt").
4. Batch có áp trần 20% trên tổng số đạt tốt trở lên không.

Test chuẩn: `go test -run 'TestEvaluate|TestCalculateStaffScore' ./internal/core/kpi/...`

### 3.4 Đánh giá học sinh dùng sai thông tư (TT 27 vs TT 89)

Mốc chuyển tiếp là **01/12/2026**. `Registry.ResolveStrategy(evalDate)` trong
`internal/core/assessment/strategy.go` quyết định strategy. Nếu học bạ năm cũ bị tính lại theo luật mới
→ kiểm tra xem có ai gọi nhầm strategy thay vì để Registry phân giải theo ngày.

Test chuẩn: `go test -run TestRegistry ./internal/core/assessment/...`

### 3.5 Tiếp nhận học sinh bị từ chối 422

Đây là **hành vi đúng** theo Điều 20 Nghị định 13/2023/NĐ-CP: thiếu Digital Consent của cha mẹ/người giám hộ.
Hướng dẫn đơn vị thu thập đồng thuận (`VNEID_SSO`, `SMS_OTP`, hoặc `PORTAL_SIGN`) rồi gửi lại. **Không** bypass kiểm tra này.

### 3.6 API chậm bất thường

- Xem access log Gin: latency từng request.
- Endpoint tính toán thuần (quota, KPI) phải trả về < 50ms; nếu chậm hơn → nghi ngờ allocation vòng lặp hoặc log quá mức.
- Khi đã tích hợp DB: kiểm tra slow query, chỉ mục trên `student_id`, `staff_id`, `term_id`.

## 4. Kênh thông báo & leo thang

- **Nội bộ dev**: issue GitHub gắn nhãn `priority/P0`..`P3`, mention CODEOWNERS của vùng mã nguồn bị lỗi.
- **Tới nhà trường (BGH)**: thông báo qua kênh triển khai dự án, nêu rõ phân hệ ảnh hưởng và thời gian dự kiến khắc phục.
- **P0 liên quan dữ liệu cá nhân**: báo cáo theo Điều 23 Nghị định 13/2023/NĐ-CP — thông báo vi phạm tới cơ quan chuyên trách bảo vệ dữ liệu cá nhân trong **72 giờ**.

## 5. Rollback

- **Mã nguồn**: nhánh `main` được bảo vệ, mọi thay đổi qua PR → revert PR bằng `git revert <merge-commit>` rồi mở PR revert.
- **Container**: giữ image tag theo commit (`bin/api` build từ `git rev-parse --short HEAD`); rollback bằng cách chạy lại image của tag trước.
- **Database** (khi đã tích hợp): chỉ rollback sau khi có bản dump; migration phải thuận nghịch và được test trên bản sao.

## 6. Sau sự cố (postmortem)

Với P0/P1, ghi lại trong issue: nguyên nhân gốc, tác động (trường nào, khoảng thời gian), cách phát hiện,
hành động khắc phục, và **hành động phòng ngừa** (test bổ sung, cảnh báo mới, sửa tài liệu).
