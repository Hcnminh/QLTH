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

## Tài liệu dự án

| Tài liệu | Nội dung |
|---|---|
| [`DE_XUAT_CHUC_NANG_QLTH_KPI.md`](DE_XUAT_CHUC_NANG_QLTH_KPI.md) | Đề xuất chức năng, căn cứ pháp lý, KPI theo vị trí việc làm |
| [`KIEN_TRUC_CORE_GO_VA_ADMIN.md`](KIEN_TRUC_CORE_GO_VA_ADMIN.md) | Kiến trúc Go core, GoAdmin, tái sử dụng đa cấp học |
| [`KIEN_TRUC_THICH_UNG_THAY_DOI_QUY_DINH.md`](KIEN_TRUC_THICH_UNG_THAY_DOI_QUY_DINH.md) | Temporal Policy & Strategy Versioning khi luật thay đổi |
| [`docs/phap_ly/DAC_TA_THIET_KE_CHI_TIET_HE_THONG.md`](docs/phap_ly/DAC_TA_THIET_KE_CHI_TIET_HE_THONG.md) | Database schema, API endpoints, thuật toán |
| [`docs/phap_ly/`](docs/phap_ly/) | Toàn văn PDF các thông tư, nghị định làm căn cứ |

## Đóng góp

- Agent tự động: đọc [`AGENTS.md`](AGENTS.md) trước khi sửa mã nguồn.
- Báo lỗi: dùng template trong `.github/ISSUE_TEMPLATE/`.
- Gửi PR: điền theo `.github/pull_request_template.md`.

## Giấy phép

Chưa công bố giấy phép — liên hệ chủ sở hữu repository trước khi tái sử dụng.
