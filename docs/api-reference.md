<!-- File này được SINH TỰ ĐỘNG từ docs/openapi.yaml bằng scripts/generate_api_docs.py. -->
<!-- KHÔNG sửa tay. Sửa docs/openapi.yaml rồi chạy `make docs`. -->

# EduPrimary-KPI Core API — Tài liệu tham chiếu

Phiên bản spec: **1.0.0** · OpenAPI 3.0.3

API backend cho Phần mềm Quản lý Trường Tiểu học tích hợp đánh giá KPI (QLTH / EduPrimary-KPI).

## Căn cứ pháp lý của các nghiệp vụ
- **Thông tư 05/2025/TT-BGDĐT** — chế độ làm việc giáo viên phổ thông (805 tiết/năm, GVCN tiểu học giảm 4 tiết/tuần).
- **Thông tư 27/2020/TT-BGDĐT** — đánh giá học sinh tiểu học; **Thông tư 89/2026/TT-BGDĐT** sửa đổi, hiệu lực 01/12/2026.
- **Nghị định 90/2020/NĐ-CP** và **Nghị định 48/2023/NĐ-CP** — đánh giá, xếp loại viên chức.
- **Luật Thi đua, Khen thưởng 06/2022/QH15** — điều kiện sáng kiến/giải pháp đổi mới.
- **Nghị định 13/2023/NĐ-CP Điều 20** — bảo vệ dữ liệu cá nhân trẻ em (Digital Consent).
- **Thông tư 29/2024/TT-BGDĐT**, **Thông tư 19/2026/TT-BGDĐT** — cấm dạy thêm đối với học sinh tiểu học.

## Xác thực & phân quyền
Tất cả endpoint dưới `/api/v1` yêu cầu hai header:
- `X-School-ID` (bắt buộc): UUID của trường — phục vụ multi-tenancy. Thiếu header này trả `401`.
- `X-User-Role` (bắt buộc): vai trò người gọi. Role không nằm trong danh sách cho phép của endpoint trả `403`.

Vai trò hợp lệ: `BGH`, `TO_TRUONG`, `GVCN`, `GVBM`, `KE_TOAN`, `VAN_THU`, `PHHS`, `SUPER_ADMIN` (toàn quyền).

## Khuôn phản hồi chung
```json
{"success": true, "code": 200, "message": "...", "data": {}}
```

## Servers

| URL | Môi trường |
|---|---|
| `http://localhost:8080` | Local development (make run) |
| `https://api.qlth.example.vn` | Production (placeholder — thay bằng tên miền thật khi triển khai VPS1 qua Traefik) |

## Mục lục endpoint

| Method | Path | Tóm tắt |
|---|---|---|
| `GET` | [`/health`](#gethealth) | Health check |
| `POST` | [`/api/v1/boarding/cancel-meal`](#cancelboardingmeal) | Báo hủy suất ăn bán trú và xác định hoàn tiền |
| `POST` | [`/api/v1/evaluations/ai-copilot`](#aipedagogicalcopilot) | Trợ lý sư phạm gợi ý nhận xét học sinh |
| `POST` | [`/api/v1/kpi/evaluate-individual`](#evaluateindividualkpi) | Đánh giá KPI 360 độ cho một cán bộ/giáo viên |
| `POST` | [`/api/v1/kpi/evaluate-school-batch`](#evaluateschoolbatchkpi) | Xếp loại KPI toàn trường kèm trần 20% xuất sắc |
| `POST` | [`/api/v1/students/onboard`](#onboardstudent) | Tiếp nhận hồ sơ học sinh (tuyển sinh đầu cấp) |
| `POST` | [`/api/v1/teaching/calculate-quota`](#calculateteachingquota) | Tính định mức tiết dạy và số tiết thừa giờ |

## system

Health check và trạng thái dịch vụ

### `GET /health`

<a id="gethealth"></a>

**Health check**

Endpoint công khai, không yêu cầu header phân quyền. Dùng cho Docker/Compose healthcheck và load balancer probe.

### Responses

| Mã | Mô tả |
|---|---|
| `200` | Dịch vụ hoạt động bình thường |

---

## teaching

Định mức tiết dạy và tính thừa giờ (Thông tư 05/2025/TT-BGDĐT)

### `POST /api/v1/teaching/calculate-quota`

<a id="calculateteachingquota"></a>

**Tính định mức tiết dạy và số tiết thừa giờ**

Áp dụng **Thông tư 05/2025/TT-BGDĐT**:
- Định mức chuẩn tiểu học **805 tiết/năm** (PTDT bán trú / lớp khuyết tật: 21 tiết/tuần).
- Giảm trừ: GVCN tiểu học **4 tiết/tuần**, tổ trưởng 3, tổ phó 1, nuôi con dưới 12 tháng 3, tư vấn tâm lý 3–8.
- Quy đổi: dạy hè tiếng Việt cho trẻ DTTS trước lớp 1 hệ số 1.5; dạy hòa nhập khuyết tật hệ số 0.2.
- `overtime_hours` = tổng tiết quy đổi − định mức thực, không âm.

Role được phép: `BGH`, `KE_TOAN`, `GVCN`, `GVBM`, `SUPER_ADMIN`.

| Tên | Vị trí | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|---|
| `X-School-ID` | header | string | Có | UUID của trường (multi-tenancy). Thiếu header này trả 401. |
| `X-User-Role` | header | string | Có | Vai trò người gọi. `SUPER_ADMIN` bypass mọi kiểm tra role. Giá trị: `BGH`, `TO_TRUONG`, `GVCN`, `GVBM`, `KE_TOAN`, `VAN_THU`, `PHHS`, `SUPER_ADMIN`. |

### Request body

| Trường | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|
| `teacher_id` | string | Có |  |
| `teacher_name` | string | Có |  |
| `school_type` | string | Không | Loại trường quyết định định mức chuẩn Giá trị: `PRIMARY`, `ETHNIC_BOARDING_PRIMARY`, `SPECIAL_EDUCATION`. |
| `is_homeroom` | boolean | Không | Giáo viên chủ nhiệm — giảm 4 tiết/tuần (tiểu học) |
| `is_head_subject` | boolean | Không | Tổ trưởng chuyên môn — giảm 3 tiết/tuần |
| `is_deputy_head_subject` | boolean | Không | Tổ phó chuyên môn — giảm 1 tiết/tuần |
| `has_young_child` | boolean | Không | Nữ giáo viên nuôi con dưới 12 tháng — giảm 3 tiết/tuần |
| `counseling_hours_week` | integer | Không | Số tiết kiêm nhiệm tư vấn tâm lý học sinh mỗi tuần |
| `inclusive_teaching_wk` | integer | Không | Số tiết dạy hòa nhập trẻ khuyết tật mỗi tuần (quy đổi hệ số 0.2) |
| `pre_grade1_ethnic_hours` | integer | Không | Số tiết dạy tiếng Việt cho trẻ DTTS trước lớp 1 trong hè (hệ số 1.5) |
| `actual_taught_hours` | integer | Không | Tổng số tiết thực tế đã giảng dạy cả năm |
| `weeks_in_school_year` | integer | Không | Số tuần thực dạy; mặc định 35 nếu ≤ 0 |

### Responses

| Mã | Mô tả |
|---|---|
| `200` | Tính toán thành công |
| `400` | Dữ liệu đầu vào không hợp lệ |
| `401` | Thiếu header X-School-ID |
| `403` | Role không được phép truy cập endpoint này |

---

## kpi

Đánh giá KPI 360 độ và xếp loại viên chức (NĐ 90/2020, NĐ 48/2023)

### `POST /api/v1/kpi/evaluate-individual`

<a id="evaluateindividualkpi"></a>

**Đánh giá KPI 360 độ cho một cán bộ/giáo viên**

Tính điểm tổng hợp theo trọng số cố định **40% hệ thống + 30% quản lý + 15% đồng nghiệp + 15% phụ huynh**,
sau đó quy đổi sang 4 mức xếp loại theo **NĐ 90/2020/NĐ-CP** và **NĐ 48/2023/NĐ-CP**.

Điểm liệt (bắt buộc xếp `KHÔNG HOÀN THÀNH NHIỆM VỤ` bất kể điểm số):
bị kỷ luật trong năm; vi phạm quy định cấm dạy thêm tiểu học (TT 29/2024, TT 19/2026);
để xảy ra tai nạn thương tích hoặc mất an toàn thực phẩm nghiêm trọng.

Ngưỡng: ≥90 + có sáng kiến/giải pháp đổi mới → xuất sắc; ≥90 thiếu sáng kiến → tốt; ≥75 → tốt; ≥50 → hoàn thành; <50 → không hoàn thành.

Role được phép: `BGH`, `TO_TRUONG`, `SUPER_ADMIN`.

| Tên | Vị trí | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|---|
| `X-School-ID` | header | string | Có | UUID của trường (multi-tenancy). Thiếu header này trả 401. |
| `X-User-Role` | header | string | Có | Vai trò người gọi. `SUPER_ADMIN` bypass mọi kiểm tra role. Giá trị: `BGH`, `TO_TRUONG`, `GVCN`, `GVBM`, `KE_TOAN`, `VAN_THU`, `PHHS`, `SUPER_ADMIN`. |

### Request body

| Trường | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|
| `staff_id` | string (uuid) | Có |  |
| `staff_name` | string | Không |  |
| `job_position` | string | Không | Vị trí việc làm theo Thông tư 20/2023/TT-BGDĐT |
| `automated_system_score` | number (float) | Có | Điểm hệ thống đo tự động — trọng số 40% |
| `manager_assessment_score` | number (float) | Có | Điểm BGH/khối trưởng chấm — trọng số 30% |
| `peer_review_score` | number (float) | Có | Điểm đồng nghiệp đánh giá chéo — trọng số 15% |
| `parent_satisfaction_score` | number (float) | Có | Điểm khảo sát phụ huynh (CSAT) — trọng số 15% |
| `criteria_details` | array<object> | Không |  |
| `has_disciplinary_action` | boolean | Không | Bị kỷ luật Đảng/chính quyền trong năm (NĐ 48/2023) — điểm liệt |
| `violated_tutoring_rules` | boolean | Không | Vi phạm quy định cấm dạy thêm tiểu học (TT 29/2024, TT 19/2026) — điểm liệt |
| `has_severe_safety_incident` | boolean | Không | Để xảy ra tai nạn thương tích / mất ATTP nghiêm trọng — điểm liệt |
| `has_innovative_solution` | boolean | Không | Có sáng kiến hoặc giải pháp đổi mới sáng tạo được công nhận (Luật TĐKT 2022) |

### Responses

| Mã | Mô tả |
|---|---|
| `200` | Đánh giá thành công |
| `400` | Dữ liệu đầu vào không hợp lệ |
| `401` | Thiếu header X-School-ID |
| `403` | Role không được phép truy cập endpoint này |

---

### `POST /api/v1/kpi/evaluate-school-batch`

<a id="evaluateschoolbatchkpi"></a>

**Xếp loại KPI toàn trường kèm trần 20% xuất sắc**

Đánh giá hàng loạt toàn bộ cán bộ, sau đó áp dụng **trần 20%** trên tổng số người đạt "hoàn thành tốt" trở lên
theo **NĐ 90/2020/NĐ-CP**. Nếu số ứng viên xuất sắc vượt trần, hệ thống giữ lại những người điểm cao nhất
và hạ số còn lại xuống mức "hoàn thành tốt", có ghi rõ lý do trong `reason`.

Role được phép: `BGH`, `SUPER_ADMIN`.

| Tên | Vị trí | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|---|
| `X-School-ID` | header | string | Có | UUID của trường (multi-tenancy). Thiếu header này trả 401. |
| `X-User-Role` | header | string | Có | Vai trò người gọi. `SUPER_ADMIN` bypass mọi kiểm tra role. Giá trị: `BGH`, `TO_TRUONG`, `GVCN`, `GVBM`, `KE_TOAN`, `VAN_THU`, `PHHS`, `SUPER_ADMIN`. |

### Request body

Mảng các phần tử kiểu `object`.

| Trường | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|
| `staff_id` | string (uuid) | Không |  |
| `staff_name` | string | Không |  |
| `job_position` | string | Không | Vị trí việc làm theo Thông tư 20/2023/TT-BGDĐT |
| `automated_system_score` | number (float) | Không | Điểm hệ thống đo tự động — trọng số 40% |
| `manager_assessment_score` | number (float) | Không | Điểm BGH/khối trưởng chấm — trọng số 30% |
| `peer_review_score` | number (float) | Không | Điểm đồng nghiệp đánh giá chéo — trọng số 15% |
| `parent_satisfaction_score` | number (float) | Không | Điểm khảo sát phụ huynh (CSAT) — trọng số 15% |
| `criteria_details` | array<object> | Không |  |
| `has_disciplinary_action` | boolean | Không | Bị kỷ luật Đảng/chính quyền trong năm (NĐ 48/2023) — điểm liệt |
| `violated_tutoring_rules` | boolean | Không | Vi phạm quy định cấm dạy thêm tiểu học (TT 29/2024, TT 19/2026) — điểm liệt |
| `has_severe_safety_incident` | boolean | Không | Để xảy ra tai nạn thương tích / mất ATTP nghiêm trọng — điểm liệt |
| `has_innovative_solution` | boolean | Không | Có sáng kiến hoặc giải pháp đổi mới sáng tạo được công nhận (Luật TĐKT 2022) |

### Responses

| Mã | Mô tả |
|---|---|
| `200` | Xếp loại thành công |
| `400` | Dữ liệu đầu vào không hợp lệ |
| `401` | Thiếu header X-School-ID |
| `403` | Role không được phép truy cập endpoint này |
| `500` | Danh sách rỗng hoặc lỗi nội bộ |

---

## evaluations

Trợ lý sư phạm gợi ý nhận xét học sinh (TT 27/2020, TT 89/2026)

### `POST /api/v1/evaluations/ai-copilot`

<a id="aipedagogicalcopilot"></a>

**Trợ lý sư phạm gợi ý nhận xét học sinh**

Sinh 2–3 gợi ý nhận xét chuẩn mực sư phạm theo **Thông tư 27/2020/TT-BGDĐT**
(và **TT 89/2026/TT-BGDĐT** từ 01/12/2026), dựa trên mức đạt của môn học và điểm kiểm tra định kỳ.
Giúp giáo viên tránh nhận xét rập khuôn và giảm thời gian ghi sổ.

Role được phép: `BGH`, `GVCN`, `GVBM`, `SUPER_ADMIN`.

| Tên | Vị trí | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|---|
| `X-School-ID` | header | string | Có | UUID của trường (multi-tenancy). Thiếu header này trả 401. |
| `X-User-Role` | header | string | Có | Vai trò người gọi. `SUPER_ADMIN` bypass mọi kiểm tra role. Giá trị: `BGH`, `TO_TRUONG`, `GVCN`, `GVBM`, `KE_TOAN`, `VAN_THU`, `PHHS`, `SUPER_ADMIN`. |

### Request body

| Trường | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|
| `subject_name` | string | Có | Tên môn học hoặc hoạt động giáo dục |
| `achievement_level` | string | Có | Mức đạt theo TT 27/2020 Giá trị: `HTT`, `HT`, `CHT`. |
| `exam_score` | number (float) | Không | Điểm kiểm tra định kỳ (chỉ với môn có điểm số) |
| `diligence_level` | string | Không | Mức phẩm chất Chăm chỉ Giá trị: `T`, `D`, `C`. |
| `student_name` | string | Có |  |

### Responses

| Mã | Mô tả |
|---|---|
| `200` | Sinh gợi ý thành công |
| `400` | Dữ liệu đầu vào không hợp lệ |
| `401` | Thiếu header X-School-ID |
| `403` | Role không được phép truy cập endpoint này |

---

## boarding

Bán trú — báo hủy suất ăn và hoàn tiền

### `POST /api/v1/boarding/cancel-meal`

<a id="cancelboardingmeal"></a>

**Báo hủy suất ăn bán trú và xác định hoàn tiền**

Quy tắc nhà trường: báo hủy **trước 08:00 sáng** (`cancellation_hour < 8`) được hoàn/bù trừ tiền ăn
vào kỳ sau; báo hủy từ 08:00 trở đi không hoàn tiền do bếp đã chuẩn bị thực phẩm.
`meal_fee_per_day` mặc định 35000 VNĐ nếu bỏ trống hoặc ≤ 0.

Role được phép: `BGH`, `KE_TOAN`, `GVCN`, `PHHS`, `SUPER_ADMIN`.

| Tên | Vị trí | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|---|
| `X-School-ID` | header | string | Có | UUID của trường (multi-tenancy). Thiếu header này trả 401. |
| `X-User-Role` | header | string | Có | Vai trò người gọi. `SUPER_ADMIN` bypass mọi kiểm tra role. Giá trị: `BGH`, `TO_TRUONG`, `GVCN`, `GVBM`, `KE_TOAN`, `VAN_THU`, `PHHS`, `SUPER_ADMIN`. |

### Request body

| Trường | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|
| `student_id` | string (uuid) | Có |  |
| `attendance_date` | string (date) | Có | Ngày ăn theo định dạng YYYY-MM-DD |
| `cancellation_hour` | integer | Có | Giờ phụ huynh bấm báo hủy (0-23) |
| `meal_fee_per_day` | number (float) | Không | Tiền ăn mỗi ngày (VNĐ). Mặc định 35000. |

### Responses

| Mã | Mô tả |
|---|---|
| `200` | Đã ghi nhận hủy suất ăn |
| `400` | Dữ liệu đầu vào không hợp lệ |
| `401` | Thiếu header X-School-ID |
| `403` | Role không được phép truy cập endpoint này |

---

## students

Hồ sơ học sinh và tiếp nhận tuyển sinh

### `POST /api/v1/students/onboard`

<a id="onboardstudent"></a>

**Tiếp nhận hồ sơ học sinh (tuyển sinh đầu cấp)**

Tiếp nhận học sinh mới, đối soát số định danh cá nhân (VNeID / Luật Căn cước 2023).

**Ràng buộc pháp lý bắt buộc** theo Điều 20 Nghị định 13/2023/NĐ-CP:
`has_digital_consent` phải là `true` — tức đã thu được sự đồng ý của cha mẹ/người giám hộ.
Nếu thiếu, API trả **422 Unprocessable Entity** và không tạo hồ sơ.

Role được phép: `BGH`, `VAN_THU`, `SUPER_ADMIN`.

| Tên | Vị trí | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|---|
| `X-School-ID` | header | string | Có | UUID của trường (multi-tenancy). Thiếu header này trả 401. |
| `X-User-Role` | header | string | Có | Vai trò người gọi. `SUPER_ADMIN` bypass mọi kiểm tra role. Giá trị: `BGH`, `TO_TRUONG`, `GVCN`, `GVBM`, `KE_TOAN`, `VAN_THU`, `PHHS`, `SUPER_ADMIN`. |

### Request body

| Trường | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|
| `id` | string (uuid) | Không | Server tự sinh nếu bỏ trống |
| `school_id` | string (uuid) | Không |  |
| `national_id` | string | Không | Số định danh cá nhân 12 số (Luật Căn cước 2023). KHÔNG ghi vào log. |
| `student_code` | string | Có |  |
| `full_name` | string | Có |  |
| `birth_date` | string (date-time) | Có |  |
| `gender` | string | Không |  Giá trị: `MALE`, `FEMALE`. |
| `ethnicity` | string | Không |  |
| `current_class_id` | string (uuid) | Không |  |
| `status` | string | Không |  Giá trị: `STUDYING`, `TRANSFERRED`, `DROPPED`, `GRADUATED`. |
| `is_inclusive_education` | boolean | Không | Học sinh khuyết tật học hòa nhập (Thông tư 03/2018/TT-BGDĐT) |
| `guardian_name` | string | Có |  |
| `guardian_national_id` | string | Không |  |
| `guardian_phone` | string | Có |  |
| `guardian_relationship` | string | Không |  Giá trị: `FATHER`, `MOTHER`, `LEGAL_GUARDIAN`. |
| `has_digital_consent` | boolean | Có | BẮT BUỘC true — sự đồng ý của cha mẹ/người giám hộ theo Điều 20 NĐ 13/2023. Thiếu sẽ trả 422. |
| `consent_signed_at` | string (date-time) | Không |  |
| `consent_method` | string | Không | Phương thức thu thập đồng thuận Giá trị: `VNEID_SSO`, `SMS_OTP`, `PORTAL_SIGN`. |
| `created_at` | string (date-time) | Không |  |

### Responses

| Mã | Mô tả |
|---|---|
| `201` | Tiếp nhận thành công |
| `400` | Dữ liệu đầu vào không hợp lệ |
| `401` | Thiếu header X-School-ID |
| `403` | Role không được phép truy cập endpoint này |
| `422` | Thiếu sự đồng thuận số của cha mẹ/người giám hộ (NĐ 13/2023 Điều 20) |

---
