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
| `GET` | [`/api/v1/boarding/menus`](#listboardingmenus) | Lấy danh sách thực đơn và phân tích dinh dưỡng |
| `GET` | [`/api/v1/lesson-plans`](#listlessonplans) | Danh sách kế hoạch bài dạy của trường |
| `GET` | [`/api/v1/lesson-plans/{id}`](#getlessonplan) | Chi tiết kế hoạch bài dạy và lịch sử góp ý |
| `GET` | [`/api/v1/students/{id}/health`](#getstudenthealthprofile) | Lấy hồ sơ sức khỏe và nhật ký y tế học sinh |
| `GET` | [`/api/v1/timetable/slots`](#gettimetableslots) | Lấy danh sách các tiết thời khóa biểu |
| `GET` | [`/health`](#gethealth) | Health check |
| `POST` | [`/api/v1/boarding/cancel-meal`](#cancelboardingmeal) | Báo hủy suất ăn bán trú và xác định hoàn tiền |
| `POST` | [`/api/v1/boarding/menus`](#createboardingmenu) | Lập kế hoạch thực đơn bán trú mới |
| `POST` | [`/api/v1/boarding/menus/{id}/sample-retention`](#verifyfoodsampleretention) | Xác nhận kiểm tra lưu mẫu thức ăn 24h |
| `POST` | [`/api/v1/evaluations/ai-copilot`](#aipedagogicalcopilot) | Trợ lý sư phạm gợi ý nhận xét học sinh |
| `POST` | [`/api/v1/evaluations/behaviors`](#recordbehaviorevaluation) | Đánh giá phẩm chất và năng lực học sinh |
| `POST` | [`/api/v1/evaluations/subject-scores`](#recordsubjectscore) | Nhập đánh giá kết quả học tập môn học |
| `POST` | [`/api/v1/kpi/evaluate-individual`](#evaluateindividualkpi) | Đánh giá KPI 360 độ cho một cán bộ/giáo viên |
| `POST` | [`/api/v1/kpi/evaluate-school-batch`](#evaluateschoolbatchkpi) | Xếp loại KPI toàn trường kèm trần 20% xuất sắc |
| `POST` | [`/api/v1/lesson-plans/submit`](#submitlessonplan) | Nộp kế hoạch bài dạy (CV 2345 / STEM TT 86) |
| `POST` | [`/api/v1/lesson-plans/{id}/approve`](#approvelessonplan) | Phê duyệt hoặc góp ý kế hoạch bài dạy phân tầng |
| `POST` | [`/api/v1/students/onboard`](#onboardstudent) | Tiếp nhận hồ sơ học sinh (tuyển sinh đầu cấp) |
| `POST` | [`/api/v1/students/{id}/health`](#recordstudenthealth) | Ghi nhận kết quả khám sức khỏe học sinh |
| `POST` | [`/api/v1/students/{id}/health/medical-events`](#addmedicalevent) | Ghi nhận sự kiện y tế hoặc sơ cấp cứu tại trường |
| `POST` | [`/api/v1/students/{id}/health/prescriptions`](#addprescriptiondelivery) | Phụ huynh gửi yêu cầu cấp phát thuốc tại trường |
| `POST` | [`/api/v1/students/{id}/health/prescriptions/{rx_id}/administer`](#administerprescription) | Nhân viên y tế xác nhận cho uống thuốc |
| `POST` | [`/api/v1/teaching/calculate-quota`](#calculateteachingquota) | Tính định mức tiết dạy và số tiết thừa giờ |
| `POST` | [`/api/v1/timetable/slots`](#configuretimetableslot) | Cấu hình tiết thời khóa biểu và kiểm soát không trùng lịch |
| `POST` | [`/api/v1/transcripts/generate`](#generatetranscriptsnapshot) | Xuất snapshot học bạ số bất biến kèm giải quyết quy định TT 27 vs TT 89 |

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
| `X-User-Role` | header | string | Có | Vai trò người gọi. `SUPER_ADMIN` bypass mọi kiểm tra role. Giá trị: `BGH`, `TO_TRUONG`, `GVCN`, `GVBM`, `KE_TOAN`, `VAN_THU`, `PHHS`, `Y_TE`, `NHAN_VIEN_Y_TE`, `BEP_TRUONG`, `SUPER_ADMIN`. |

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
| `X-User-Role` | header | string | Có | Vai trò người gọi. `SUPER_ADMIN` bypass mọi kiểm tra role. Giá trị: `BGH`, `TO_TRUONG`, `GVCN`, `GVBM`, `KE_TOAN`, `VAN_THU`, `PHHS`, `Y_TE`, `NHAN_VIEN_Y_TE`, `BEP_TRUONG`, `SUPER_ADMIN`. |

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
| `X-User-Role` | header | string | Có | Vai trò người gọi. `SUPER_ADMIN` bypass mọi kiểm tra role. Giá trị: `BGH`, `TO_TRUONG`, `GVCN`, `GVBM`, `KE_TOAN`, `VAN_THU`, `PHHS`, `Y_TE`, `NHAN_VIEN_Y_TE`, `BEP_TRUONG`, `SUPER_ADMIN`. |

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

Đánh giá học sinh và trợ lý sư phạm (TT 27/2020, TT 89/2026)

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
| `X-User-Role` | header | string | Có | Vai trò người gọi. `SUPER_ADMIN` bypass mọi kiểm tra role. Giá trị: `BGH`, `TO_TRUONG`, `GVCN`, `GVBM`, `KE_TOAN`, `VAN_THU`, `PHHS`, `Y_TE`, `NHAN_VIEN_Y_TE`, `BEP_TRUONG`, `SUPER_ADMIN`. |

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

### `POST /api/v1/evaluations/subject-scores`

<a id="recordsubjectscore"></a>

**Nhập đánh giá kết quả học tập môn học**

Đánh giá thường xuyên và định kỳ môn học theo **Thông tư 27/2020** hoặc **Thông tư 89/2026**.
Mức đạt bắt buộc: `HTT`, `HT`, `CHT`. Điểm kiểm tra định kỳ từ 1.0 đến 10.0 (chỉ với môn có điểm số).

Role được phép: `BGH`, `GVCN`, `GVBM`, `SUPER_ADMIN`.

| Tên | Vị trí | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|---|
| `X-School-ID` | header | string | Có | UUID của trường (multi-tenancy). Thiếu header này trả 401. |
| `X-User-Role` | header | string | Có | Vai trò người gọi. `SUPER_ADMIN` bypass mọi kiểm tra role. Giá trị: `BGH`, `TO_TRUONG`, `GVCN`, `GVBM`, `KE_TOAN`, `VAN_THU`, `PHHS`, `Y_TE`, `NHAN_VIEN_Y_TE`, `BEP_TRUONG`, `SUPER_ADMIN`. |

### Request body

| Trường | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|
| `student_id` | string (uuid) | Có |  |
| `term_id` | string | Có |  |
| `subject_id` | string | Có |  |
| `teacher_id` | string (uuid) | Không |  |
| `achievement_level` | string | Có |  Giá trị: `HTT`, `HT`, `CHT`. |
| `exam_score` | number (float) | Không |  |
| `teacher_comment` | string | Không |  |
| `ai_suggested_comment` | string | Không |  |
| `is_ai_assisted` | boolean | Không |  |
| `is_computer_based_test` | boolean | Không |  |
| `is_centralized_exam` | boolean | Không |  |
| `evaluation_date` | string | Không |  |

### Responses

| Mã | Mô tả |
|---|---|
| `200` | Ghi nhận đánh giá thành công |
| `400` | Dữ liệu đầu vào không hợp lệ |
| `401` | Thiếu header X-School-ID |
| `403` | Role không được phép truy cập endpoint này |

---

### `POST /api/v1/evaluations/behaviors`

<a id="recordbehaviorevaluation"></a>

**Đánh giá phẩm chất và năng lực học sinh**

Đánh giá 5 phẩm chất chủ yếu (Yêu nước, Nhân ái, Chăm chỉ, Trung thực, Trách nhiệm)
và 10 năng lực (3 năng lực chung + 7 năng lực đặc thù) theo **Thông tư 27/2020** và **Thông tư 89/2026**.
Mức đánh giá: `T` (Tốt), `D` (Đạt), `C` (Cần cố gắng).

Role được phép: `BGH`, `GVCN`, `SUPER_ADMIN`.

| Tên | Vị trí | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|---|
| `X-School-ID` | header | string | Có | UUID của trường (multi-tenancy). Thiếu header này trả 401. |
| `X-User-Role` | header | string | Có | Vai trò người gọi. `SUPER_ADMIN` bypass mọi kiểm tra role. Giá trị: `BGH`, `TO_TRUONG`, `GVCN`, `GVBM`, `KE_TOAN`, `VAN_THU`, `PHHS`, `Y_TE`, `NHAN_VIEN_Y_TE`, `BEP_TRUONG`, `SUPER_ADMIN`. |

### Request body

| Trường | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|
| `student_id` | string (uuid) | Có |  |
| `term_id` | string | Có |  |
| `evaluator_id` | string (uuid) | Không |  |
| `patriotism` | string | Có |  Giá trị: `T`, `D`, `C`. |
| `compassion` | string | Có |  Giá trị: `T`, `D`, `C`. |
| `diligence` | string | Có |  Giá trị: `T`, `D`, `C`. |
| `honesty` | string | Có |  Giá trị: `T`, `D`, `C`. |
| `responsibility` | string | Có |  Giá trị: `T`, `D`, `C`. |
| `qualities_comment` | string | Không |  |
| `autonomy_and_self_learning` | string | Không |  Giá trị: `T`, `D`, `C`. |
| `communication_and_collaboration` | string | Không |  Giá trị: `T`, `D`, `C`. |
| `problem_solving_and_creativity` | string | Không |  Giá trị: `T`, `D`, `C`. |
| `language_competence` | string | Không |  Giá trị: `T`, `D`, `C`. |
| `mathematical_competence` | string | Không |  Giá trị: `T`, `D`, `C`. |
| `scientific_competence` | string | Không |  Giá trị: `T`, `D`, `C`. |
| `technological_competence` | string | Không |  Giá trị: `T`, `D`, `C`. |
| `computer_competence` | string | Không |  Giá trị: `T`, `D`, `C`. |
| `aesthetic_competence` | string | Không |  Giá trị: `T`, `D`, `C`. |
| `physical_competence` | string | Không |  Giá trị: `T`, `D`, `C`. |
| `competencies_comment` | string | Không |  |
| `evaluation_date` | string | Không |  |

### Responses

| Mã | Mô tả |
|---|---|
| `200` | Ghi nhận đánh giá thành công |
| `400` | Dữ liệu đầu vào không hợp lệ |
| `401` | Thiếu header X-School-ID |
| `403` | Role không được phép truy cập endpoint này |

---

## transcripts

Quản lý học bạ số và snapshot bất biến (TT 27/2020, TT 89/2026)

### `POST /api/v1/transcripts/generate`

<a id="generatetranscriptsnapshot"></a>

**Xuất snapshot học bạ số bất biến kèm giải quyết quy định TT 27 vs TT 89**

Tổng hợp đánh giá cả năm của học sinh, tự động giải quyết quy định áp dụng:
- Trước 01/12/2026: Áp dụng **Thông tư 27/2020/TT-BGDĐT**.
- Từ 01/12/2026 trở đi: Áp dụng **Thông tư 89/2026/TT-BGDĐT** (kiểm tra trên máy tính, đề chung lớp 5, giáo dục hòa nhập IEP).
Lưu snapshot bất biến (JSONB) vào `academic_transcripts_snapshots` với trạng thái khóa `is_locked = true`.

Role được phép: `BGH`, `GVCN`, `SUPER_ADMIN`.

| Tên | Vị trí | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|---|
| `X-School-ID` | header | string | Có | UUID của trường (multi-tenancy). Thiếu header này trả 401. |
| `X-User-Role` | header | string | Có | Vai trò người gọi. `SUPER_ADMIN` bypass mọi kiểm tra role. Giá trị: `BGH`, `TO_TRUONG`, `GVCN`, `GVBM`, `KE_TOAN`, `VAN_THU`, `PHHS`, `Y_TE`, `NHAN_VIEN_Y_TE`, `BEP_TRUONG`, `SUPER_ADMIN`. |

### Request body

| Trường | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|
| `student_id` | string (uuid) | Có |  |
| `academic_year_id` | string | Có |  |
| `grade_level` | integer | Không |  |
| `evaluation_date` | string | Không |  |
| `teacher_signature` | string | Không |  |
| `principal_signature` | string | Không |  |

### Responses

| Mã | Mô tả |
|---|---|
| `200` | Tạo snapshot học bạ thành công |
| `400` | Dữ liệu đầu vào không hợp lệ |
| `401` | Thiếu header X-School-ID |
| `403` | Role không được phép truy cập endpoint này |

---

## boarding

Bán trú — thực đơn, phân tích dinh dưỡng Viện Dinh Dưỡng, lưu mẫu 24h và hoàn tiền

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
| `X-User-Role` | header | string | Có | Vai trò người gọi. `SUPER_ADMIN` bypass mọi kiểm tra role. Giá trị: `BGH`, `TO_TRUONG`, `GVCN`, `GVBM`, `KE_TOAN`, `VAN_THU`, `PHHS`, `Y_TE`, `NHAN_VIEN_Y_TE`, `BEP_TRUONG`, `SUPER_ADMIN`. |

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

### `GET /api/v1/boarding/menus`

<a id="listboardingmenus"></a>

**Lấy danh sách thực đơn và phân tích dinh dưỡng**

Truy vấn danh sách thực đơn bán trú kèm phân tích dinh dưỡng theo chuẩn Viện Dinh Dưỡng
(Kcal: 600-750, P: 13-20%, L: 20-30%, G: 50-65%) và kiểm tra lưu mẫu 24h.

Role được phép: `BGH`, `KE_TOAN`, `GVCN`, `GVBM`, `PHHS`, `Y_TE`, `NHAN_VIEN_Y_TE`, `SUPER_ADMIN`.

| Tên | Vị trí | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|---|
| `X-School-ID` | header | string | Có | UUID của trường (multi-tenancy). Thiếu header này trả 401. |
| `X-User-Role` | header | string | Có | Vai trò người gọi. `SUPER_ADMIN` bypass mọi kiểm tra role. Giá trị: `BGH`, `TO_TRUONG`, `GVCN`, `GVBM`, `KE_TOAN`, `VAN_THU`, `PHHS`, `Y_TE`, `NHAN_VIEN_Y_TE`, `BEP_TRUONG`, `SUPER_ADMIN`. |
| `date` | query | string | Không | Lọc theo ngày (YYYY-MM-DD) |
| `meal_type` | query | string | Không | Lọc theo bữa (LUNCH, SNACK, BREAKFAST) |

### Responses

| Mã | Mô tả |
|---|---|
| `200` | Lấy danh sách thực đơn thành công |
| `401` | Thiếu header X-School-ID |
| `403` | Role không được phép truy cập endpoint này |

---

### `POST /api/v1/boarding/menus`

<a id="createboardingmenu"></a>

**Lập kế hoạch thực đơn bán trú mới**

Tạo thực đơn mới, tự động tính tổng Kcal và tỷ lệ P:L:G từ các món ăn thành phần.

Role được phép: `BGH`, `KE_TOAN`, `Y_TE`, `NHAN_VIEN_Y_TE`, `BEP_TRUONG`, `SUPER_ADMIN`.

| Tên | Vị trí | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|---|
| `X-School-ID` | header | string | Có | UUID của trường (multi-tenancy). Thiếu header này trả 401. |
| `X-User-Role` | header | string | Có | Vai trò người gọi. `SUPER_ADMIN` bypass mọi kiểm tra role. Giá trị: `BGH`, `TO_TRUONG`, `GVCN`, `GVBM`, `KE_TOAN`, `VAN_THU`, `PHHS`, `Y_TE`, `NHAN_VIEN_Y_TE`, `BEP_TRUONG`, `SUPER_ADMIN`. |

### Request body

| Trường | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|
| `menu_date` | string (date) | Có |  |
| `meal_type` | string | Không |  |
| `dishes` | array<object> | Không |  |
| `total_kcal` | number | Không |  |
| `protein_percent` | number | Không |  |
| `lipid_percent` | number | Không |  |
| `glucid_percent` | number | Không |  |
| `sample_retention_verified` | boolean | Không |  |

### Responses

| Mã | Mô tả |
|---|---|
| `201` | Thực đơn được tạo thành công |
| `400` | Dữ liệu đầu vào không hợp lệ |
| `401` | Thiếu header X-School-ID |
| `403` | Role không được phép truy cập endpoint này |

---

### `POST /api/v1/boarding/menus/{id}/sample-retention`

<a id="verifyfoodsampleretention"></a>

**Xác nhận kiểm tra lưu mẫu thức ăn 24h**

Kiểm tra thời gian lưu mẫu >= 24h và nhiệt độ 2-8°C theo Quyết định 1246/QĐ-BYT.

Role được phép: `BGH`, `Y_TE`, `NHAN_VIEN_Y_TE`, `BEP_TRUONG`, `SUPER_ADMIN`.

| Tên | Vị trí | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|---|
| `X-School-ID` | header | string | Có | UUID của trường (multi-tenancy). Thiếu header này trả 401. |
| `X-User-Role` | header | string | Có | Vai trò người gọi. `SUPER_ADMIN` bypass mọi kiểm tra role. Giá trị: `BGH`, `TO_TRUONG`, `GVCN`, `GVBM`, `KE_TOAN`, `VAN_THU`, `PHHS`, `Y_TE`, `NHAN_VIEN_Y_TE`, `BEP_TRUONG`, `SUPER_ADMIN`. |
| `id` | path | string | Có |  |

### Request body

| Trường | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|
| `sampled_at` | string (date-time) | Có |  |
| `storage_temperature` | number | Có |  |
| `keeper_name` | string | Có |  |

### Responses

| Mã | Mô tả |
|---|---|
| `200` | Xác nhận lưu mẫu thức ăn thành công |
| `400` | Dữ liệu đầu vào không hợp lệ |
| `401` | Thiếu header X-School-ID |
| `403` | Role không được phép truy cập endpoint này |

---

## health

Y tế học đường, khám sức khỏe, BMI, thị lực và cấp phát thuốc

### `GET /api/v1/students/{id}/health`

<a id="getstudenthealthprofile"></a>

**Lấy hồ sơ sức khỏe và nhật ký y tế học sinh**

Truy vấn thông tin thể lực, tính toán BMI, phân loại thể trạng, đánh giá thị lực,
nhật ký khám sức khỏe và các đơn thuốc được cấp phát.

Role được phép: `BGH`, `GVCN`, `Y_TE`, `NHAN_VIEN_Y_TE`, `PHHS`, `SUPER_ADMIN`.

| Tên | Vị trí | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|---|
| `X-School-ID` | header | string | Có | UUID của trường (multi-tenancy). Thiếu header này trả 401. |
| `X-User-Role` | header | string | Có | Vai trò người gọi. `SUPER_ADMIN` bypass mọi kiểm tra role. Giá trị: `BGH`, `TO_TRUONG`, `GVCN`, `GVBM`, `KE_TOAN`, `VAN_THU`, `PHHS`, `Y_TE`, `NHAN_VIEN_Y_TE`, `BEP_TRUONG`, `SUPER_ADMIN`. |
| `id` | path | string | Có |  |

### Responses

| Mã | Mô tả |
|---|---|
| `200` | Lấy hồ sơ sức khỏe thành công |
| `400` | Dữ liệu đầu vào không hợp lệ |
| `401` | Thiếu header X-School-ID |
| `403` | Role không được phép truy cập endpoint này |

---

### `POST /api/v1/students/{id}/health`

<a id="recordstudenthealth"></a>

**Ghi nhận kết quả khám sức khỏe học sinh**

Cập nhật chiều cao, cân nặng, tự động tính BMI và phân loại thể trạng (Suy dinh dưỡng / Bình thường / Thừa cân / Béo phì).

Role được phép: `BGH`, `Y_TE`, `NHAN_VIEN_Y_TE`, `SUPER_ADMIN`.

| Tên | Vị trí | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|---|
| `X-School-ID` | header | string | Có | UUID của trường (multi-tenancy). Thiếu header này trả 401. |
| `X-User-Role` | header | string | Có | Vai trò người gọi. `SUPER_ADMIN` bypass mọi kiểm tra role. Giá trị: `BGH`, `TO_TRUONG`, `GVCN`, `GVBM`, `KE_TOAN`, `VAN_THU`, `PHHS`, `Y_TE`, `NHAN_VIEN_Y_TE`, `BEP_TRUONG`, `SUPER_ADMIN`. |
| `id` | path | string | Có |  |

### Request body

| Trường | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|
| `checkup_date` | string (date) | Không |  |
| `height_cm` | number | Có |  |
| `weight_kg` | number | Có |  |
| `left_eye_vision` | string | Không |  |
| `right_eye_vision` | string | Không |  |
| `dental_condition` | string | Không |  |
| `chronic_conditions` | string | Không |  |

### Responses

| Mã | Mô tả |
|---|---|
| `200` | Cập nhật khám sức khỏe thành công |
| `400` | Dữ liệu đầu vào không hợp lệ |
| `401` | Thiếu header X-School-ID |
| `403` | Role không được phép truy cập endpoint này |

---

### `POST /api/v1/students/{id}/health/prescriptions`

<a id="addprescriptiondelivery"></a>

**Phụ huynh gửi yêu cầu cấp phát thuốc tại trường**

Gửi thông tin thuốc cần uống, liều lượng và thời gian cho học sinh.

Role được phép: `BGH`, `GVCN`, `PHHS`, `Y_TE`, `NHAN_VIEN_Y_TE`, `SUPER_ADMIN`.

| Tên | Vị trí | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|---|
| `X-School-ID` | header | string | Có | UUID của trường (multi-tenancy). Thiếu header này trả 401. |
| `X-User-Role` | header | string | Có | Vai trò người gọi. `SUPER_ADMIN` bypass mọi kiểm tra role. Giá trị: `BGH`, `TO_TRUONG`, `GVCN`, `GVBM`, `KE_TOAN`, `VAN_THU`, `PHHS`, `Y_TE`, `NHAN_VIEN_Y_TE`, `BEP_TRUONG`, `SUPER_ADMIN`. |
| `id` | path | string | Có |  |

### Request body

| Trường | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|
| `medication_name` | string | Có |  |
| `dosage` | string | Có |  |
| `administration_time` | string | Có |  |
| `instructions` | string | Không |  |
| `parent_notes` | string | Không |  |

### Responses

| Mã | Mô tả |
|---|---|
| `201` | Gửi yêu cầu đơn thuốc thành công |
| `400` | Dữ liệu đầu vào không hợp lệ |
| `401` | Thiếu header X-School-ID |
| `403` | Role không được phép truy cập endpoint này |

---

### `POST /api/v1/students/{id}/health/prescriptions/{rx_id}/administer`

<a id="administerprescription"></a>

**Nhân viên y tế xác nhận cho uống thuốc**

Xác nhận thời điểm và nhân viên thực hiện cấp phát thuốc cho học sinh.

Role được phép: `BGH`, `Y_TE`, `NHAN_VIEN_Y_TE`, `SUPER_ADMIN`.

| Tên | Vị trí | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|---|
| `X-School-ID` | header | string | Có | UUID của trường (multi-tenancy). Thiếu header này trả 401. |
| `X-User-Role` | header | string | Có | Vai trò người gọi. `SUPER_ADMIN` bypass mọi kiểm tra role. Giá trị: `BGH`, `TO_TRUONG`, `GVCN`, `GVBM`, `KE_TOAN`, `VAN_THU`, `PHHS`, `Y_TE`, `NHAN_VIEN_Y_TE`, `BEP_TRUONG`, `SUPER_ADMIN`. |
| `id` | path | string | Có |  |
| `rx_id` | path | string | Có |  |

### Request body

| Trường | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|
| `nurse_name` | string | Có |  |

### Responses

| Mã | Mô tả |
|---|---|
| `200` | Đã xác nhận cấp phát thuốc |
| `400` | Dữ liệu đầu vào không hợp lệ |
| `401` | Thiếu header X-School-ID |
| `403` | Role không được phép truy cập endpoint này |

---

### `POST /api/v1/students/{id}/health/medical-events`

<a id="addmedicalevent"></a>

**Ghi nhận sự kiện y tế hoặc sơ cấp cứu tại trường**

Lưu nhật ký sơ cấp cứu, xử lý sốt/tai nạn nhỏ tại phòng y tế học đường.

Role được phép: `BGH`, `Y_TE`, `NHAN_VIEN_Y_TE`, `SUPER_ADMIN`.

| Tên | Vị trí | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|---|
| `X-School-ID` | header | string | Có | UUID của trường (multi-tenancy). Thiếu header này trả 401. |
| `X-User-Role` | header | string | Có | Vai trò người gọi. `SUPER_ADMIN` bypass mọi kiểm tra role. Giá trị: `BGH`, `TO_TRUONG`, `GVCN`, `GVBM`, `KE_TOAN`, `VAN_THU`, `PHHS`, `Y_TE`, `NHAN_VIEN_Y_TE`, `BEP_TRUONG`, `SUPER_ADMIN`. |
| `id` | path | string | Có |  |

### Request body

| Trường | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|
| `event_type` | string | Có |  |
| `symptoms` | string | Không |  |
| `treatment` | string | Có |  |
| `medication_given` | string | Không |  |
| `nurse_name` | string | Không |  |
| `parent_notified` | boolean | Không |  |

### Responses

| Mã | Mô tả |
|---|---|
| `201` | Ghi nhận sự kiện y tế thành công |
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
| `X-User-Role` | header | string | Có | Vai trò người gọi. `SUPER_ADMIN` bypass mọi kiểm tra role. Giá trị: `BGH`, `TO_TRUONG`, `GVCN`, `GVBM`, `KE_TOAN`, `VAN_THU`, `PHHS`, `Y_TE`, `NHAN_VIEN_Y_TE`, `BEP_TRUONG`, `SUPER_ADMIN`. |

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

## curriculum

Quản lý kế hoạch bài dạy (CV 2345/STEM TT 86) và xếp thời khóa biểu (TT 05/2025)

### `POST /api/v1/lesson-plans/submit`

<a id="submitlessonplan"></a>

**Nộp kế hoạch bài dạy (CV 2345 / STEM TT 86)**

Nộp kế hoạch bài dạy điện tử (giáo án) tuân thủ:
- **Công văn 2345/BGDĐT-GDTH** (4 bước hoạt động: Khởi động, Khám phá, Luyện tập, Vận dụng; 3 mục tiêu: Kiến thức, Năng lực, Phẩm chất).
- **Thông tư 86/BGDĐT & CV 909/BGDĐT-GDTH** (giáo dục STEM tiểu học).
- Giới hạn từ tuần 1 đến tuần 35 theo Thông tư 05/2025/TT-BGDĐT.
- Khởi tạo trạng thái `SUBMITTED` (hoặc `DRAFT` nếu bật cờ lưu nháp).

Role được phép: `BGH`, `TO_TRUONG`, `GVCN`, `GVBM`, `GV`, `SUPER_ADMIN`.

| Tên | Vị trí | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|---|
| `X-School-ID` | header | string | Có | UUID của trường (multi-tenancy). Thiếu header này trả 401. |
| `X-User-Role` | header | string | Có | Vai trò người gọi. `SUPER_ADMIN` bypass mọi kiểm tra role. Giá trị: `BGH`, `TO_TRUONG`, `GVCN`, `GVBM`, `KE_TOAN`, `VAN_THU`, `PHHS`, `Y_TE`, `NHAN_VIEN_Y_TE`, `BEP_TRUONG`, `SUPER_ADMIN`. |

### Request body

| Trường | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|
| `school_id` | string (uuid) | Không |  |
| `teacher_id` | string (uuid) | Có |  |
| `teacher_name` | string | Không |  |
| `class_id` | string (uuid) | Không |  |
| `grade_level` | integer | Có |  |
| `subject_id` | string | Có |  |
| `title` | string | Có |  |
| `week_number` | integer | Có |  |
| `duration_slots` | integer | Không |  |
| `curriculum_standard` | string | Không |  Giá trị: `CV_2345`, `STEM_TT_86`. |
| `objectives` | object | Không |  |
| `equipment` | string | Không |  |
| `activities` | array<object> | Không |  |
| `attachment_url` | string | Không |  |
| `is_draft` | boolean | Không |  |

### Responses

| Mã | Mô tả |
|---|---|
| `201` | Nộp kế hoạch bài dạy thành công |
| `400` | Dữ liệu đầu vào không hợp lệ |
| `401` | Thiếu header X-School-ID |
| `403` | Role không được phép truy cập endpoint này |

---

### `POST /api/v1/lesson-plans/{id}/approve`

<a id="approvelessonplan"></a>

**Phê duyệt hoặc góp ý kế hoạch bài dạy phân tầng**

Quy trình thẩm định phân tầng:
- **Tổ trưởng chuyên môn (TO_TRUONG)**: Duyệt bước 1 chuyển trạng thái sang `DEPARTMENT_HEAD_APPROVED`.
- **Ban Giám hiệu (BGH)**: Duyệt hoàn tất chuyển trạng thái sang `APPROVED` kèm chữ ký số SmartCA/TSA.
- Khi từ chối (`REJECT`) hoặc yêu cầu sửa (`REQUEST_REVISION`): bắt buộc nhập `comments` góp ý sư phạm.

Role được phép: `BGH`, `TO_TRUONG`, `SUPER_ADMIN`.

| Tên | Vị trí | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|---|
| `id` | path | string | Có | UUID của kế hoạch bài dạy |
| `X-School-ID` | header | string | Có | UUID của trường (multi-tenancy). Thiếu header này trả 401. |
| `X-User-Role` | header | string | Có | Vai trò người gọi. `SUPER_ADMIN` bypass mọi kiểm tra role. Giá trị: `BGH`, `TO_TRUONG`, `GVCN`, `GVBM`, `KE_TOAN`, `VAN_THU`, `PHHS`, `Y_TE`, `NHAN_VIEN_Y_TE`, `BEP_TRUONG`, `SUPER_ADMIN`. |

### Request body

| Trường | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|
| `action` | string | Có |  Giá trị: `APPROVE`, `REJECT`, `REQUEST_REVISION`. |
| `reviewer_id` | string (uuid) | Không |  |
| `reviewer_role` | string | Không |  Giá trị: `TO_TRUONG`, `BGH`, `SUPER_ADMIN`. |
| `reviewer_name` | string | Không |  |
| `comments` | string | Không | Ý kiến nhận xét, góp ý sư phạm (bắt buộc khi REJECT hoặc REQUEST_REVISION) |
| `digital_signature` | string | Không |  |

### Responses

| Mã | Mô tả |
|---|---|
| `200` | Phê duyệt hoặc góp ý kế hoạch bài dạy thành công |
| `400` | Dữ liệu đầu vào không hợp lệ |
| `401` | Thiếu header X-School-ID |
| `403` | Role không được phép truy cập endpoint này |
| `404` | Không tìm thấy kế hoạch bài dạy |

---

### `GET /api/v1/lesson-plans/{id}`

<a id="getlessonplan"></a>

**Chi tiết kế hoạch bài dạy và lịch sử góp ý**

Lấy thông tin chi tiết của một kế hoạch bài dạy cùng toàn bộ lịch sử duyệt và ý kiến nhận xét của tổ chuyên môn và Ban Giám hiệu.

| Tên | Vị trí | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|---|
| `id` | path | string | Có | UUID của kế hoạch bài dạy |
| `X-School-ID` | header | string | Có | UUID của trường (multi-tenancy). Thiếu header này trả 401. |
| `X-User-Role` | header | string | Có | Vai trò người gọi. `SUPER_ADMIN` bypass mọi kiểm tra role. Giá trị: `BGH`, `TO_TRUONG`, `GVCN`, `GVBM`, `KE_TOAN`, `VAN_THU`, `PHHS`, `Y_TE`, `NHAN_VIEN_Y_TE`, `BEP_TRUONG`, `SUPER_ADMIN`. |

### Responses

| Mã | Mô tả |
|---|---|
| `200` | Lấy chi tiết kế hoạch bài dạy thành công |
| `401` | Thiếu header X-School-ID |
| `403` | Role không được phép truy cập endpoint này |
| `404` | Không tìm thấy kế hoạch bài dạy |

---

### `GET /api/v1/lesson-plans`

<a id="listlessonplans"></a>

**Danh sách kế hoạch bài dạy của trường**

Truy vấn danh sách kế hoạch bài dạy đã nộp, hỗ trợ bộ lọc theo môn học, giáo viên, trạng thái duyệt.

| Tên | Vị trí | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|---|
| `X-School-ID` | header | string | Có | UUID của trường (multi-tenancy). Thiếu header này trả 401. |
| `X-User-Role` | header | string | Có | Vai trò người gọi. `SUPER_ADMIN` bypass mọi kiểm tra role. Giá trị: `BGH`, `TO_TRUONG`, `GVCN`, `GVBM`, `KE_TOAN`, `VAN_THU`, `PHHS`, `Y_TE`, `NHAN_VIEN_Y_TE`, `BEP_TRUONG`, `SUPER_ADMIN`. |

### Responses

| Mã | Mô tả |
|---|---|
| `200` | Lấy danh sách thành công |
| `401` | Thiếu header X-School-ID |
| `403` | Role không được phép truy cập endpoint này |

---

### `GET /api/v1/timetable/slots`

<a id="gettimetableslots"></a>

**Lấy danh sách các tiết thời khóa biểu**

Truy vấn lịch thời khóa biểu theo năm học, học kỳ, lớp học hoặc giáo viên.

| Tên | Vị trí | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|---|
| `X-School-ID` | header | string | Có | UUID của trường (multi-tenancy). Thiếu header này trả 401. |
| `X-User-Role` | header | string | Có | Vai trò người gọi. `SUPER_ADMIN` bypass mọi kiểm tra role. Giá trị: `BGH`, `TO_TRUONG`, `GVCN`, `GVBM`, `KE_TOAN`, `VAN_THU`, `PHHS`, `Y_TE`, `NHAN_VIEN_Y_TE`, `BEP_TRUONG`, `SUPER_ADMIN`. |
| `academic_year_id` | query | string | Không |  |
| `term_id` | query | string | Không |  |
| `class_id` | query | string | Không |  |
| `teacher_id` | query | string | Không |  |

### Responses

| Mã | Mô tả |
|---|---|
| `200` | Lấy danh sách thành công |
| `401` | Thiếu header X-School-ID |
| `403` | Role không được phép truy cập endpoint này |

---

### `POST /api/v1/timetable/slots`

<a id="configuretimetableslot"></a>

**Cấu hình tiết thời khóa biểu và kiểm soát không trùng lịch**

Cấu hình tiết thời khóa biểu:
- Đảm bảo **không trùng lịch giảng dạy của giáo viên** giữa các lớp khác nhau tại cùng thứ/tiết theo **Thông tư 05/2025/TT-BGDĐT**.
- Đảm bảo một lớp không bị xếp 2 môn học đồng thời.
- Đảm bảo một phòng chức năng không bị trùng lịch sử dụng.
- Trả về `409 Conflict` nếu phát hiện xung đột lịch.

Role được phép: `BGH`, `TO_TRUONG`, `GVCN`, `GVBM`, `SUPER_ADMIN`.

| Tên | Vị trí | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|---|
| `X-School-ID` | header | string | Có | UUID của trường (multi-tenancy). Thiếu header này trả 401. |
| `X-User-Role` | header | string | Có | Vai trò người gọi. `SUPER_ADMIN` bypass mọi kiểm tra role. Giá trị: `BGH`, `TO_TRUONG`, `GVCN`, `GVBM`, `KE_TOAN`, `VAN_THU`, `PHHS`, `Y_TE`, `NHAN_VIEN_Y_TE`, `BEP_TRUONG`, `SUPER_ADMIN`. |

### Request body

| Trường | Kiểu | Bắt buộc | Mô tả |
|---|---|---|---|
| `school_id` | string (uuid) | Không |  |
| `academic_year_id` | string | Có |  |
| `term_id` | string | Có |  |
| `class_id` | string (uuid) | Có |  |
| `class_name` | string | Không |  |
| `teacher_id` | string (uuid) | Có |  |
| `teacher_name` | string | Không |  |
| `subject_id` | string | Có |  |
| `day_of_week` | integer | Có | Thứ trong tuần (2=Thứ Hai, 6=Thứ Sáu, 7=Thứ Bảy) |
| `period` | integer | Có | Tiết học trong ngày (1-5 sáng, 6-10 chiều) |
| `room_number` | string | Không |  |
| `week_number` | integer | Không | 0 cho cả học kỳ, hoặc số tuần cụ thể |

### Responses

| Mã | Mô tả |
|---|---|
| `201` | Cấu hình tiết thời khóa biểu thành công |
| `400` | Dữ liệu đầu vào không hợp lệ |
| `401` | Thiếu header X-School-ID |
| `403` | Role không được phép truy cập endpoint này |
| `409` | Phát hiện xung đột thời khóa biểu vi phạm Thông tư 05/2025/TT-BGDĐT |

---
