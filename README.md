# QLTH — Phần mềm Quản lý Trường Tiểu học tích hợp KPI (EduPrimary-KPI)

Backend cốt lõi (Go) cho hệ thống quản lý trường tiểu học Việt Nam, tuân thủ quy định pháp luật hiện hành của Bộ GD&ĐT và tích hợp bộ công cụ đánh giá KPI nhân viên 360 độ.

## Điểm nổi bật

- **Đánh giá học sinh theo Thông tư 27/2020/TT-BGDĐT** và sẵn sàng chuyển tiếp sang **Thông tư 89/2026/TT-BGDĐT** (hiệu lực 01/12/2026) nhờ cơ chế Strategy Versioning.
- **Định mức tiết dạy theo Thông tư 05/2025/TT-BGDĐT**: chuẩn 805 tiết/năm, GVCN tiểu học giảm 4 tiết/tuần, tự tính tiết thừa giờ.
- **KPI 360 độ** (hệ thống 40% + quản lý 30% + đồng nghiệp 15% + phụ huynh 15%) quy đổi tự động sang 4 mức xếp loại viên chức theo **Nghị định 90/2020/NĐ-CP** và **Nghị định 48/2023/NĐ-CP**, kèm trần 20% "hoàn thành xuất sắc".
- **Bảo vệ dữ liệu trẻ em** theo Điều 20 **Nghị định 13/2023/NĐ-CP**: bắt buộc Digital Consent của cha mẹ/người giám hộ trước khi kích hoạt hồ sơ học sinh.

## Yêu cầu

- Go **1.27.1** trở lên (`go version`)
- Docker + Docker Compose (tùy chọn, để chạy PostgreSQL/Redis cục bộ)

## Bắt đầu nhanh

```bash
git clone https://github.com/Hcnminh/QLTH.git
cd QLTH

# Cài dependency
make deps          # tương đương: go mod download

# Build + chạy toàn bộ test
make all           # tương đương: make build && make test

# Chạy API server cục bộ (cổng 8080)
make run

# Hoặc chạy full stack (API + PostgreSQL 16 + Redis 7)
make docker-up
```

Kiểm tra server đang chạy:

```bash
curl -s http://localhost:8080/health
# {"service":"eduprimary-kpi-core","status":"healthy","time":"..."}
```

## Các lệnh Makefile

| Lệnh | Mô tả |
|---|---|
| `make deps` | Tải Go modules |
| `make build` | Biên dịch ra `bin/api` |
| `make test` | Chạy toàn bộ unit + integration test |
| `make test-browser` | Chạy kiểm thử tự động trình duyệt E2E (Playwright) trên GoAdmin |
| `make qa` | Chạy kiểm thử tương tác end-to-end (QA harness 26 ca) |
| `make lint` | Chạy `go vet`, `gofmt` check và `golangci-lint` |
| `make run` | Chạy API server ở cổng 8098 (hoặc qua biến `PORT`) |
| `make docker-up` / `make docker-down` | Bật/tắt full stack qua Docker Compose |
| `make docs` | Validate OpenAPI spec + sinh `docs/api-reference.md` |
| `make clean` | Xóa artifacts |

## Cấu trúc thư mục

```
cmd/api/            # Entrypoint HTTP server
internal/api/       # Gin router, handlers, RBAC middleware
internal/core/      # Nghiệp vụ thuần túy (không phụ thuộc HTTP/DB)
  assessment/       # Strategy TT 27/2020 & TT 89/2026
  domain/           # Entities: School, Student, Staff, ...
  kpi/              # KPI 360 engine + xếp loại NĐ 90/48
  teaching/         # Định mức tiết dạy + thừa giờ TT 05/2025
docs/               # OpenAPI spec, API reference, runbooks, đặc tả thiết kế, văn bản pháp lý
scripts/            # Validate OpenAPI + AGENTS.md, sinh tài liệu API
.factory/skills/    # Skill cho AI agent: quy tắc nghiệp vụ pháp lý bất biến
```

## API

- OpenAPI 3.0 spec: [`docs/openapi.yaml`](docs/openapi.yaml)
- Tài liệu tham chiếu (sinh tự động từ spec): [`docs/api-reference.md`](docs/api-reference.md)

Sau khi sửa spec, chạy `make docs` để validate và sinh lại tài liệu, rồi commit cả hai file.

Tất cả endpoint `/api/v1/*` yêu cầu 2 header:

- `X-School-ID: <uuid>` — mã trường (multi-tenancy)
- `X-User-Role: <role>` — một trong `BGH`, `TO_TRUONG`, `GVCN`, `GVBM`, `KE_TOAN`, `VAN_THU`, `PHHS`, `SUPER_ADMIN`

Ví dụ:

```bash
curl -s -X POST http://localhost:8080/api/v1/teaching/calculate-quota \
  -H 'Content-Type: application/json' \
  -H 'X-School-ID: 00000000-0000-0000-0000-000000000001' \
  -H 'X-User-Role: BGH' \
  -d '{"teacher_id":"GV001","teacher_name":"Lê Thị Mai","school_type":"PRIMARY",
       "is_homeroom":true,"weeks_in_school_year":35,"actual_taught_hours":720}'
```

## Vận hành

- Runbook vận hành & xử lý sự cố: [`docs/runbooks/`](docs/runbooks/)
- Hướng dẫn triển khai VPS (Podman/Traefik): `HUONG-DAN-HA-TANG-VA-TRIEN-KHAI-VPS1.md`
- Biến môi trường: xem [`.env.example`](.env.example)

## Giao diện Quản trị & Kiểm thử Trình duyệt (Browser QA)

Hệ thống tích hợp bảng điều khiển quản trị trực quan **GoAdmin** tại `/admin`, cung cấp giao diện tương tác cho cán bộ quản lý trường học và giáo viên.

### Kiểm thử tự động E2E với Playwright

Chạy bộ kiểm thử tự động trình duyệt bằng lệnh:

```bash
make test-browser
```

Bộ kiểm thử [`scripts/qa_browser.js`](scripts/qa_browser.js) sử dụng Playwright headless Chromium để tự động xác thực toàn diện:
1. **Kiểm tra sức khỏe hệ thống**: Endpoint công khai `GET /health`.
2. **Quy trình xác thực**: Đăng nhập qua `/admin/login` với tài khoản quản trị mặc định (`admin` / `admin`).
3. **Bảng điều khiển trung tâm (Dashboard)**: Kiểm tra nhận diện thương hiệu EduPrimary-KPI và toàn bộ 11 menu nghiệp vụ.
4. **Toàn bộ 11 bảng giao diện GoAdmin (10 phân hệ nghiệp vụ)**:
   - Phân hệ 1: Hồ sơ Học sinh & Consent NĐ 13/2023 (`/admin/info/students`)
   - Phân hệ 2: Danh mục Môn học (`/admin/info/subjects`), Đánh giá Môn học (`/admin/info/student_subject_evaluations`), Đánh giá Năng lực - Phẩm chất TT 27/TT 89 (`/admin/info/student_behavior_evaluations`)
   - Phân hệ 3: Giáo viên & Định mức TT 05/2025 (`/admin/info/staffs`), Lớp học (`/admin/info/classes`), Kế hoạch bài dạy CV 2345/STEM (`/admin/info/lesson_plans`)
   - Phân hệ 4: Thực đơn Dinh dưỡng & Lưu mẫu 24h QĐ 1246/QĐ-BYT (`/admin/info/boarding_menus`)
   - Phân hệ 5: Hồ sơ Sức khỏe Học sinh & Chỉ số BMI (`/admin/info/student_health_records`)
   - Phân hệ 8: Đánh giá KPI 360 độ & Trần 20% NĐ 90/48 (`/admin/info/kpi_evaluations`)
   - Phân hệ 9: Học bạ số Bất biến & Chữ ký số SmartCA TSA (`/admin/info/academic_transcripts_snapshots`)
5. **Đa vai trò người dùng (Multi-Role Workflows)**:
   - **BGH**: Kiểm tra định mức 805h, xếp TKB chống trùng lịch, phân giải TT 27/89 mốc 01/12/2026, đồng bộ CSDL MOET Outbox.
   - **GVCN**: Nhập đánh giá năng lực - phẩm chất (T/D/C), gợi ý nhận xét AI Copilot, ký số học bạ SmartCA.
   - **GVBM**: Nộp giáo án CV 2345 4 bước, nhập điểm thi và đánh giá môn học (HTT/10.0).
   - **Kế toán**: Lập thực đơn dinh dưỡng chuẩn Kcal, sinh mã VietQR EMVCo thanh toán bán trú, xuất hóa đơn điện tử NĐ 123.
   - **Y tế**: Xác nhận lưu mẫu thức ăn 24h ở 4°C, ghi nhật ký sơ cấp cứu phòng y tế học đường.
   - **Văn thư**: Tuyển sinh học sinh có consent NĐ 13 Điều 20, xác thực định danh VNeID Mức 2 sinh trắc học.
   - **Phụ huynh**: Báo cắt ăn trưa trước 8h hoàn tiền tự động, gửi đơn thuốc nhờ cô y tế cho uống.

### Ảnh chụp màn hình giao diện (Visual Screenshots)

Bộ ảnh chụp trực quan độ phân giải cao được lưu trữ tại thư mục [`docs/screenshots/`](docs/screenshots/):

| Phân hệ / Màn hình | Tệp hình ảnh | Mô tả |
|---|---|---|
| **Đăng nhập Quản trị** | [`01_login_page.png`](docs/screenshots/01_login_page.png) | Giao diện đăng nhập hệ thống EduPrimary-KPI GoAdmin |
| **Bảng điều khiển (Dashboard)** | [`02_dashboard.png`](docs/screenshots/02_dashboard.png) | Trung tâm điều hành chỉ số giáo viên, học sinh, KPI và RBAC |
| **Giáo viên & Định mức** | [`03_staff_table.png`](docs/screenshots/03_staff_table.png) | Quản lý giáo viên, phân công GVCN và định mức theo TT 05/2025 |
| **Hồ sơ Học sinh & Consent** | [`04_student_table.png`](docs/screenshots/04_student_table.png) | Hồ sơ học sinh và trạng thái bảo vệ dữ liệu trẻ em NĐ 13/2023 |
| **Đánh giá KPI 360** | [`05_kpi_table.png`](docs/screenshots/05_kpi_table.png) | Bảng tổng hợp KPI 360 độ, xếp loại thi đua và khống chế trần 20% |
| **Quản lý Lớp học** | [`06_classes_management.png`](docs/screenshots/06_classes_management.png) | Danh sách lớp học 1-5, phòng học và GVCN |
| **Kế hoạch Bài dạy CV 2345** | [`07_lesson_plans_cv2345.png`](docs/screenshots/07_lesson_plans_cv2345.png) | Kế hoạch bài dạy chuẩn CV 2345 & STEM TT 86, trạng thái phê duyệt |
| **Danh mục Môn học** | [`08_subjects_catalog.png`](docs/screenshots/08_subjects_catalog.png) | Danh mục môn học, hình thức kiểm tra điểm số / nhận xét CT GDPT 2018 |
| **Đánh giá Môn học** | [`09_subject_evaluations_tt27_89.png`](docs/screenshots/09_subject_evaluations_tt27_89.png) | Đánh giá môn học theo TT 27/2020 và TT 89/2026 (HTT/HT/CHT) |
| **Đánh giá Năng lực - Phẩm chất** | [`10_behavior_competency_evaluations.png`](docs/screenshots/10_behavior_competency_evaluations.png) | Đánh giá 5 phẩm chất, 10 năng lực cốt lõi theo thang T/D/C |
| **Học bạ số Snapshot SmartCA** | [`11_transcripts_smartca_snapshot.png`](docs/screenshots/11_transcripts_smartca_snapshot.png) | Snapshot học bạ số bất biến kèm chữ ký số và dấu thời gian TSA |
| **Thực đơn Bán trú Dinh dưỡng** | [`12_boarding_nutrition_menu.png`](docs/screenshots/12_boarding_nutrition_menu.png) | Cân bằng Kcal/P:L:G Viện Dinh Dưỡng và lưu mẫu 24h QĐ 1246/QĐ-BYT |
| **Hồ sơ Sức khỏe & BMI** | [`13_student_health_bmi.png`](docs/screenshots/13_student_health_bmi.png) | Theo dõi BMI, thị lực học đường, nha khoa và nhật ký y tế |

## Tài liệu dự án & Ấn phẩm Bàn giao

| Tài liệu / Ấn phẩm | Định dạng | Nội dung & Đối tượng sử dụng |
|---|---|---|
| **Brochure Giới thiệu Sản phẩm** | [`docs/brochure-qlth-kpi.pdf`](docs/brochure-qlth-kpi.pdf) / [`.docx`](docs/brochure-qlth-kpi.docx) | Bản in & giới thiệu tổng quan dành cho Ban Giám Hiệu, Phòng/Sở GD&ĐT: 5 đột phá pháp lý, so sánh đối thủ, phân tích ROI |
| **Hướng dẫn Sử dụng Toàn diện** | [`docs/huong-dan-su-dung.md`](docs/huong-dan-su-dung.md) / [`.pdf`](docs/huong-dan-su-dung.pdf) | Cẩm nang vận hành chi tiết step-by-step cho 7 vai trò (BGH, GVCN, GVBM, Kế toán, Y tế, Văn thư, Phụ huynh) kèm ảnh chụp giao diện |
| [`DE_XUAT_CHUC_NANG_QLTH_KPI.md`](DE_XUAT_CHUC_NANG_QLTH_KPI.md) | Markdown | Đề xuất chức năng, căn cứ pháp lý, KPI theo vị trí việc làm |
| [`KIEN_TRUC_CORE_GO_VA_ADMIN.md`](KIEN_TRUC_CORE_GO_VA_ADMIN.md) | Markdown | Kiến trúc Go core, GoAdmin, tái sử dụng đa cấp học |
| [`KIEN_TRUC_THICH_UNG_THAY_DOI_QUY_DINH.md`](KIEN_TRUC_THICH_UNG_THAY_DOI_QUY_DINH.md) | Markdown | Temporal Policy & Strategy Versioning khi luật thay đổi |
| [`docs/phap_ly/DAC_TA_THIET_KE_CHI_TIET_HE_THONG.md`](docs/phap_ly/DAC_TA_THIET_KE_CHI_TIET_HE_THONG.md) | Markdown | Database schema 16 bảng, API endpoints, thuật toán |
| [`docs/phap_ly/`](docs/phap_ly/) | PDF | Toàn văn PDF các thông tư, nghị định làm căn cứ |

## Đóng góp

- Agent tự động: đọc [`AGENTS.md`](AGENTS.md) trước khi sửa mã nguồn.
- Báo lỗi: dùng template trong `.github/ISSUE_TEMPLATE/`.
- Gửi PR: điền theo `.github/pull_request_template.md`.

## Giấy phép

Chưa công bố giấy phép — liên hệ chủ sở hữu repository trước khi tái sử dụng.
