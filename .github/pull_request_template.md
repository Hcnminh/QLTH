## Mô tả thay đổi

<!-- Thay đổi gì và vì sao. Link issue liên quan: Fixes #123 -->

## Loại thay đổi

- [ ] `feat` — tính năng mới
- [ ] `fix` — sửa lỗi
- [ ] `docs` — tài liệu, OpenAPI, runbook
- [ ] `refactor` — không đổi hành vi
- [ ] `test` — bổ sung/sửa test
- [ ] `chore` — build, CI, dependency

## Kiểm tra bắt buộc trước khi gửi

- [ ] Đã chạy `make lint` — pass, không file nào cần `gofmt`
- [ ] Đã chạy `make test` — toàn bộ test pass
- [ ] Đã chạy `make build` — biên dịch thành công
- [ ] Nếu đổi endpoint: đã cập nhật `docs/openapi.yaml` và `make docs` pass
- [ ] Nếu đổi lệnh/quy ước: đã cập nhật `README.md` và `AGENTS.md`

## Ràng buộc nghiệp vụ (đánh dấu nếu PR chạm vào các vùng này)

- [ ] **Định mức tiết dạy (TT 05/2025)** — giữ nguyên: 805 tiết/năm, GVCN tiểu học −4 tiết/tuần, 35 tuần mặc định
- [ ] **KPI & xếp loại (NĐ 90/2020, NĐ 48/2023)** — giữ nguyên trọng số 40/30/15/15, điểm liệt, trần 20% xuất sắc
- [ ] **Đánh giá học sinh (TT 27/2020, TT 89/2026)** — không sửa strategy cũ; thông tư mới = thêm strategy mới + đăng ký vào Registry
- [ ] **Dữ liệu trẻ em (NĐ 13/2023 Điều 20)** — giữ kiểm tra Digital Consent, không log dữ liệu nhạy cảm
- [ ] PR không liên quan các vùng trên

## Test đã bổ sung

<!-- Liệt kê file test mới/sửa và case đã cover. PR thêm logic nghiệp vụ mà không có test sẽ bị từ chối. -->

## Ảnh hưởng vận hành

- [ ] Không cần migration / thay đổi hạ tầng
- [ ] Cần thay đổi DB, Docker, env — đã ghi rõ trong phần mô tả và cập nhật runbook
