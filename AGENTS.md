# AGENTS.md — Hướng dẫn cho AI Agent làm việc trên repo QLTH

Repo này là **backend Go** cho phần mềm Quản lý Trường Tiểu học tích hợp KPI (EduPrimary-KPI).
Đọc tài liệu này trước khi sửa mã nguồn. Mọi thay đổi phải giữ nguyên tính đúng đắn pháp lý nêu ở mục 4.

## 1. Lệnh thiết yếu

| Mục đích | Lệnh |
|---|---|
| Cài dependency | `make deps` (= `go mod download`) |
| Build | `make build` (= `go build -v -o bin/api ./cmd/api`) |
| Chạy toàn bộ test | `make test` (= `go test -v -count=1 ./...`) |
| Chạy 1 package test | `go test -v ./internal/core/kpi/...` |
| Lint + format check | `make lint` (= `go vet ./...` + `gofmt -l .`) |
| Chạy server local | `make run` (cổng 8080, biến `PORT` để đổi) |
| Full stack (Postgres + Redis + API) | `make docker-up` / `make docker-down` |
| Validate OpenAPI spec | `make docs` |

**Luôn chạy `make lint && make test` trước khi commit.** Không commit khi test đỏ.

## 2. Kiến trúc & quy ước

- Module Go tên `qlth`. Entrypoint: `cmd/api/main.go`.
- **`internal/core/` là nghiệp vụ thuần túy** — không import `gin`, `net/http`, hay package DB.
  Giữ nguyên ranh giới này: logic pháp lý phải test được bằng unit test không cần HTTP/DB.
- **`internal/api/`** là lớp transport: `router.go` (định tuyến), `handlers/` (bind JSON + gọi core), `middleware/` (RBAC).
  Handler chỉ được: validate input → gọi core engine → trả JSON. Không nhét logic nghiệp vụ vào handler.
- Mọi response theo khuôn: `{"success": bool, "code": int, "message": string, "data": ...}`.
- Endpoint `/api/v1/*` bắt buộc qua `middleware.RequireRole(...)` và yêu cầu header `X-School-ID` + `X-User-Role`.
  Khi thêm route mới, **phải** khai báo danh sách role được phép.
- Go chuẩn: `gofmt`, tên exported dùng PascalCase, lỗi trả về tường minh (không panic trong core).

## 3. Quy trình phát triển

1. Tạo issue từ template trong `.github/ISSUE_TEMPLATE/` (bug hoặc feature), gán nhãn `priority/*`, `type/*`, `area/*`.
2. Nhánh đặt tên `feat/<issue>-mo-ta-ngan` hoặc `fix/<issue>-mo-ta-ngan`, cắt từ `main`.
3. Viết test **trước hoặc cùng lúc** với code. Package core mới phải có file `*_test.go`.
4. Cập nhật `docs/openapi.yaml` khi thêm/đổi endpoint; cập nhật README/AGENTS.md khi đổi lệnh hoặc quy ước.
5. Mở PR theo `.github/pull_request_template.md`. Nhánh `main` được bảo vệ: cần PR + review, không push trực tiếp.
6. Commit message theo Conventional Commits: `feat:`, `fix:`, `docs:`, `test:`, `chore:`.

## 4. Ràng buộc nghiệp vụ — KHÔNG ĐƯỢC VI PHẠM

Đây là các quy tắc suy ra từ văn bản pháp luật. Sửa sai sẽ tạo ra kết quả đánh giá trái luật.

**Định mức tiết dạy (Thông tư 05/2025/TT-BGDĐT)** — `internal/core/teaching/quota.go`
- Chuẩn tiểu học: **805 tiết/năm**; trường PTDT bán trú / lớp khuyết tật: **21 tiết/tuần**.
- GVCN tiểu học giảm **4 tiết/tuần** (không phải 3 như TT 28/2009 cũ).
- Tổ trưởng −3, tổ phó −1, nuôi con dưới 12 tháng −3, tư vấn tâm lý −3..−8 tiết/tuần.
- Số tuần mặc định **35**. Tiết thừa giờ = tổng tiết quy đổi − định mức thực, không bao giờ âm.

**KPI & xếp loại viên chức (NĐ 90/2020, NĐ 48/2023, Luật TĐKT 2022)** — `internal/core/kpi/kpi.go`
- Trọng số cố định: hệ thống **40%** + quản lý **30%** + đồng nghiệp **15%** + phụ huynh **15%**.
- Điểm liệt → "Không hoàn thành nhiệm vụ" bất kể điểm số: bị kỷ luật trong năm; vi phạm cấm dạy thêm tiểu học (TT 29/2024, TT 19/2026); để xảy ra tai nạn thương tích/mất ATTP nghiêm trọng.
- ≥90 điểm **và** có sáng kiến/giải pháp đổi mới → "Hoàn thành xuất sắc". ≥90 mà thiếu sáng kiến → chỉ "Hoàn thành tốt".
- Ngưỡng: ≥75 tốt, ≥50 hoàn thành, <50 không hoàn thành.
- Trần **20%** trên tổng số người đạt tốt trở lên; vượt trần thì giữ người điểm cao nhất, số còn lại hạ xuống "tốt".

**Đánh giá học sinh** — `internal/core/assessment/strategy.go`
- Mốc chuyển tiếp **01/12/2026**: trước đó dùng `TT27_2020`, từ đó trở đi dùng `TT89_2026`. Do `Registry.ResolveStrategy(evalDate)` quyết định — không hardcode ở nơi khác.
- Mức đạt: `HTT` / `HT` / `CHT`. Phẩm chất & năng lực: `T` / `D` / `C`. Điểm kiểm tra 1.0–10.0.
- Thêm thông tư mới = thêm struct strategy mới + đăng ký vào Registry. **Không sửa strategy cũ**, để học bạ lịch sử không bị tính lại.

**Dữ liệu trẻ em (Điều 20 NĐ 13/2023/NĐ-CP)**
- Hồ sơ học sinh chỉ được kích hoạt khi `HasDigitalConsent = true`; thiếu thì trả **422**.
- Không log, không trả về API các trường nhạy cảm (số định danh, bệnh án, ảnh) khi chưa cần thiết.

## 5. Test & QA tương tác

- Unit test core: `internal/core/{kpi,teaching,assessment}/*_test.go`.
- Integration test HTTP + RBAC: `internal/api/router_test.go` (dùng `httptest`, không cần server thật).
- QA thủ công: `make run` rồi gọi `GET /health` và các endpoint theo `docs/openapi.yaml`
  (đường dẫn QA đầy đủ trong `docs/runbooks/01-van-hanh-he-thong.md`).

## 6. File KHÔNG được sửa/xóa

- `docs/phap_ly/*.pdf` — toàn văn văn bản pháp luật, chỉ đọc.
- `HUONG-DAN-HA-TANG-VA-TRIEN-KHAI-VPS1.md` — tài liệu hạ tầng production của chủ repo.
- Các file `*.md` phân tích ở root là tài liệu đã thẩm định; sửa phải có lý do rõ ràng trong PR.

## 7. Tài liệu liên quan

- Runbook vận hành: `docs/runbooks/01-van-hanh-he-thong.md`
- Runbook xử lý sự cố: `docs/runbooks/02-xu-ly-su-co.md`
- Đặc tả thiết kế + DB schema: `docs/phap_ly/DAC_TA_THIET_KE_CHI_TIET_HE_THONG.md`
- OpenAPI spec: `docs/openapi.yaml` (bản Markdown sinh tự động: `docs/api-reference.md`)
- Biến môi trường: `.env.example`

## 8. Skill cho agent

Repo khai báo skill trong `.factory/skills/qlth-legal-business-rules/SKILL.md` —
chứa toàn bộ hằng số và quy tắc pháp lý bất biến (định mức tiết dạy, trọng số KPI, điểm liệt,
trần 20%, mốc chuyển tiếp TT 89, Digital Consent). Agent **phải** đọc skill này trước khi
sửa bất kỳ file nào trong `internal/core/`.
