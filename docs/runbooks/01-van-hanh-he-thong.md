# Runbook 01 — Vận hành hệ thống EduPrimary-KPI

Phạm vi: khởi động, giám sát, dừng dịch vụ API core (Go/Gin) ở môi trường local và VPS.
Người thực hiện: developer, trực vận hành, hoặc AI agent được cấp quyền.

## 1. Thông tin dịch vụ

| Thuộc tính | Giá trị |
|---|---|
| Service | `eduprimary-kpi-core` |
| Entrypoint | `cmd/api/main.go` |
| Cổng mặc định | `8080` (đổi qua biến môi trường `PORT`) |
| Health endpoint | `GET /health` (công khai, không cần header) |
| Stack phụ thuộc | PostgreSQL 16, Redis 7 (qua `docker-compose.yml`) |
| Binary sau build | `bin/api` |

## 2. Khởi động môi trường local

```bash
# 1. Cài dependency
make deps

# 2. Build + chạy toàn bộ test (bắt buộc trước khi vận hành)
make all

# 3a. Chạy API đơn lẻ (không cần DB — core engine hiện chạy in-memory)
make run

# 3b. Hoặc chạy full stack API + PostgreSQL + Redis
make docker-up
```

Xác nhận dịch vụ sống:

```bash
curl -sf http://localhost:8080/health
# Kỳ vọng: {"service":"eduprimary-kpi-core","status":"healthy","time":"<RFC3339>"}
```

## 3. Kiểm tra nhanh các endpoint nghiệp vụ (smoke test thủ công)

Mọi endpoint `/api/v1/*` cần 2 header: `X-School-ID` (UUID bất kỳ khi test) và `X-User-Role`.

```bash
SCHOOL='00000000-0000-0000-0000-000000000001'

# Định mức tiết dạy GVCN tiểu học (kỳ vọng: net 665, overtime 55 với 720 tiết thực dạy)
curl -sf -X POST http://localhost:8080/api/v1/teaching/calculate-quota \
  -H 'Content-Type: application/json' -H "X-School-ID: $SCHOOL" -H 'X-User-Role: BGH' \
  -d '{"teacher_id":"GV001","teacher_name":"Test","school_type":"PRIMARY","is_homeroom":true,"weeks_in_school_year":35,"actual_taught_hours":720}'

# KPI cá nhân (kỳ vọng: rank xuất sắc khi score >= 90 và has_innovative_solution=true)
curl -sf -X POST http://localhost:8080/api/v1/kpi/evaluate-individual \
  -H 'Content-Type: application/json' -H "X-School-ID: $SCHOOL" -H 'X-User-Role: BGH' \
  -d '{"staff_id":"'"$(uuidgen | tr A-Z a-z)"'","staff_name":"Test","automated_system_score":95,"manager_assessment_score":95,"peer_review_score":95,"parent_satisfaction_score":95,"has_innovative_solution":true}'

# Digital Consent: kỳ vọng 422 khi has_digital_consent=false (NĐ 13/2023 Điều 20)
curl -s -o /dev/null -w '%{http_code}\n' -X POST http://localhost:8080/api/v1/students/onboard \
  -H 'Content-Type: application/json' -H "X-School-ID: $SCHOOL" -H 'X-User-Role: VAN_THU' \
  -d '{"student_code":"HS999","full_name":"Test","birth_date":"2019-01-01T00:00:00Z","guardian_name":"Parent","guardian_phone":"0900000000","has_digital_consent":false}'
```

Bộ smoke test tự động tương đương nằm trong `internal/api/router_test.go` — chạy bằng `make test`.

## 4. Giám sát & log

- **Log**: service ghi log chuẩn `log` ra stdout; Gin access log in từng request theo dạng
  `[GIN] <time> | <status> | <latency> | <method> "<path>"`.
- **Local**: xem log trực tiếp ở terminal chạy `make run`.
- **Docker**: `docker compose logs -f api` (local) hoặc `podman logs -f <container>` trên VPS1.
- **Latency bất thường** (> 500ms cho các endpoint tính toán thuần): xem Runbook 02 mục "API chậm".

## 5. Dừng dịch vụ

```bash
# Docker stack
make docker-down

# Process chạy nền thủ công (macOS/Linux)
lsof -ti :8080 | xargs kill
```

Trên VPS1 (Podman, không có Docker daemon): dùng `podman-compose down` trong thư mục project,
chi tiết theo `HUONG-DAN-HA-TANG-VA-TRIEN-KHAI-VPS1.md` ở root repo (Traefik router, network `web`, SELinux).

## 6. Sao lưu & dữ liệu

- Core API hiện **chưa ghi DB** (engine thuần in-memory); dữ liệu bền nằm ở PostgreSQL của docker-compose.
- Volume cần backup: `qlth_postgres_data`, `qlth_redis_data` (`docker volume ls`).
- Khi tích hợp DB thật (lộ trình Giai đoạn 2): thêm lịch `pg_dump` vào runbook này và mục 4 Runbook 02.
