#!/usr/bin/env bash
# QA end-to-end harness cho EduPrimary-KPI Core API.
#
# Mục đích: cho phép mission worker / QA agent (hoặc con người) khởi động app,
# điều khiển app như người dùng thật qua HTTP, quan sát log trên disk, rồi tắt sạch.
# Theo khuyến nghị Factory Missions: "one command to start the app", "route logs
# to the filesystem", "provide a way to send input".
#
# Cách dùng:
#   ./scripts/qa_e2e.sh            # build -> start -> exercise 22 ca -> log ra disk -> stop
#   PORT=9090 ./scripts/qa_e2e.sh  # đổi cổng nếu 8080 bận
#
# Kết quả:
#   - Log server:  logs/qa-server.log   (đọc được sau khi script chạy xong)
#   - Log request: logs/qa-requests.log (mỗi ca: tên, HTTP status, pass/fail)
#   - Exit code 0 nếu TẤT CẢ ca pass, 1 nếu có ca fail (kèm danh sách fail).
#
# Không cần database: core engines hiện chạy in-memory. Khi persistence được
# tích hợp (Giai đoạn 2), thêm bước `make docker-up` vào section START bên dưới
# và export DATABASE_URL/REDIS_URL trước khi start server.

set -uo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT"

find_free_port() {
  python3 -c 'import socket; s=socket.socket(); s.bind(("", 0)); print(s.getsockname()[1]); s.close()'
}

# Chọn cổng: nếu cổng được chỉ định bị chiếm thì tự động chọn cổng rảnh
DEFAULT_PORT="${PORT:-8098}"
if lsof -Pi :"$DEFAULT_PORT" -sTCP:LISTEN -t >/dev/null 2>&1; then
  PORT=$(find_free_port)
  echo "Cổng ${DEFAULT_PORT} đang bận, tự động chuyển sang cổng rảnh: ${PORT}"
else
  PORT="$DEFAULT_PORT"
fi

BASE="http://localhost:${PORT}"
LOG_DIR="$REPO_ROOT/logs"
SERVER_LOG="$LOG_DIR/qa-server.log"
REQ_LOG="$LOG_DIR/qa-requests.log"
SCHOOL_ID="00000000-0000-0000-0000-000000000001"

PASS=0
FAIL=0
FAILED_CASES=()

log_case() {
  # $1=tên ca, $2=http_code kỳ vọng, $3=http_code thực tế, $4=ghi chú
  local name="$1" want="$2" got="$3" note="${4:-}"
  if [ "$want" = "$got" ]; then
    PASS=$((PASS+1))
    echo "PASS  ${name}  (HTTP ${got}) ${note}" | tee -a "$REQ_LOG"
  else
    FAIL=$((FAIL+1))
    FAILED_CASES+=("${name}: expected ${want}, got ${got}")
    echo "FAIL  ${name}  (expected HTTP ${want}, got ${got}) ${note}" | tee -a "$REQ_LOG"
  fi
}

# http_case <tên> <method> <path> <json_body|-> <expected_code> <role> <extra_curl_args...>
http_case() {
  local name="$1" method="$2" path="$3" body="$4" want="$5" role="$6"
  shift 6
  local args=(-s -o "$LOG_DIR/last-response.json" -w '%{http_code}' -X "$method" "${BASE}${path}")
  args+=(-H 'Content-Type: application/json' -H "X-School-ID: ${SCHOOL_ID}")
  if [ -n "$role" ]; then
    args+=(-H "X-User-Role: ${role}")
  fi
  if [ "$body" != "-" ]; then
    args+=(-d "$body")
  fi
  args+=("$@")
  local code
  code=$(curl "${args[@]}" || echo "000")
  log_case "$name" "$want" "$code"
}

# ---------------------------------------------------------------- BUILD
echo "== BUILD =="
mkdir -p "$LOG_DIR" bin
: > "$SERVER_LOG"
: > "$REQ_LOG"
go build -o bin/api ./cmd/api || { echo "Build failed"; exit 1; }
echo "Binary: bin/api ($(go version))"

# ---------------------------------------------------------------- START
echo "== START server on :${PORT} =="
PORT="$PORT" GIN_MODE=release ./bin/api > "$SERVER_LOG" 2>&1 &
SERVER_PID=$!

cleanup() {
  if kill -0 "$SERVER_PID" 2>/dev/null; then
    kill "$SERVER_PID" 2>/dev/null
    wait "$SERVER_PID" 2>/dev/null
  fi
}
trap cleanup EXIT

# Chờ server sống (tối đa 15s)
READY=0
for i in $(seq 1 30); do
  if curl -sf "${BASE}/health" | grep -q '"service":"eduprimary-kpi-core"'; then READY=1; break; fi
  sleep 0.5
done
if [ "$READY" -ne 1 ]; then
  echo "Server không khởi động được trên cổng ${PORT}. Log:"
  cat "$SERVER_LOG"
  exit 1
fi
echo "Server sẵn sàng (PID ${SERVER_PID}) trên ${BASE}, log: ${SERVER_LOG}"

# ---------------------------------------------------------------- CASES
echo "== EXERCISE USER FLOWS =="

# --- 1. Health (public) ---
http_case "health-check" GET "/health" "-" 200 ""

# --- 2. RBAC / multi-tenancy (middleware) ---
# Thiếu X-School-ID -> 401
code=$(curl -s -o /dev/null -w '%{http_code}' -X POST "${BASE}/api/v1/teaching/calculate-quota" \
  -H 'Content-Type: application/json' -H 'X-User-Role: BGH' -d '{}' || echo 000)
log_case "rbac-missing-school-id" 401 "$code"

# Role không đủ quyền (PHHS gọi batch KPI) -> 403
code=$(curl -s -o /dev/null -w '%{http_code}' -X POST "${BASE}/api/v1/kpi/evaluate-school-batch" \
  -H 'Content-Type: application/json' -H "X-School-ID: ${SCHOOL_ID}" -H 'X-User-Role: PHHS' -d '[]' || echo 000)
log_case "rbac-forbidden-role" 403 "$code"

# --- 3. Định mức tiết dạy (TT 05/2025/TT-BGDĐT) ---
# GVCN tiểu học: 805 - 4*35 = 665; dạy 720 -> thừa 55
http_case "quota-gvcn-tieu-hoc" POST "/api/v1/teaching/calculate-quota" \
  '{"teacher_id":"GV001","teacher_name":"Lê Thị Mai","school_type":"PRIMARY","is_homeroom":true,"weeks_in_school_year":35,"actual_taught_hours":720}' \
  200 "BGH"
grep -q '"net_annual_quota":665' "$LOG_DIR/last-response.json" \
  && grep -q '"overtime_hours":55' "$LOG_DIR/last-response.json" \
  && log_case "quota-gvcn-values" 0 0 "(665 tiết thực dạy, thừa 55)" \
  || log_case "quota-gvcn-values" 0 1 "response không khớp kỳ vọng"

# PTDT bán trú: 21*35 = 735; dạy 700 -> không thừa
http_case "quota-ptdt-ban-tru" POST "/api/v1/teaching/calculate-quota" \
  '{"teacher_id":"GV010","teacher_name":"Giàng A Páo","school_type":"ETHNIC_BOARDING_PRIMARY","weeks_in_school_year":35,"actual_taught_hours":700}' \
  200 "BGH"
grep -q '"standard_annual_hours":735' "$LOG_DIR/last-response.json" \
  && grep -q '"overtime_hours":0' "$LOG_DIR/last-response.json" \
  && log_case "quota-ptdt-values" 0 0 "(chuẩn 735, thừa 0)" \
  || log_case "quota-ptdt-values" 0 1 "response không khớp kỳ vọng"

# Payload sai kiểu -> 400
http_case "quota-invalid-payload" POST "/api/v1/teaching/calculate-quota" '"khong-phai-json-object"' 400 "BGH"

# --- 4. KPI 360 + xếp loại viên chức (NĐ 90/2020, NĐ 48/2023) ---
# Điểm 92 + có sáng kiến -> Hoàn thành xuất sắc
http_case "kpi-excellent-with-innovation" POST "/api/v1/kpi/evaluate-individual" \
  '{"staff_id":"11111111-1111-1111-1111-111111111111","staff_name":"Lê Văn Cường","job_position":"HOMEROOM_TEACHER","automated_system_score":92,"manager_assessment_score":92,"peer_review_score":92,"parent_satisfaction_score":92,"has_innovative_solution":true}' \
  200 "BGH"
grep -q 'XUẤT SẮC' "$LOG_DIR/last-response.json" \
  && log_case "kpi-excellent-rank" 0 0 "(rank xuất sắc)" \
  || log_case "kpi-excellent-rank" 0 1 "rank không phải xuất sắc"

# Điểm 95 nhưng vi phạm dạy thêm tiểu học (TT 29/2024) -> điểm liệt, Không hoàn thành
http_case "kpi-fatal-tutoring-violation" POST "/api/v1/kpi/evaluate-individual" \
  '{"staff_id":"22222222-2222-2222-2222-222222222222","staff_name":"Trần Văn D","automated_system_score":95,"manager_assessment_score":95,"peer_review_score":95,"parent_satisfaction_score":95,"has_innovative_solution":true,"violated_tutoring_rules":true}' \
  200 "BGH"
grep -q 'KHÔNG HOÀN THÀNH' "$LOG_DIR/last-response.json" \
  && log_case "kpi-fatal-flaw-rank" 0 0 "(điểm liệt áp dụng đúng)" \
  || log_case "kpi-fatal-flaw-rank" 0 1 "điểm liệt không được áp dụng"

# Batch 6 người: 5 hợp lệ điểm cao + 1 vi phạm -> tối đa 1 suất xuất sắc (trần 20%)
http_case "kpi-batch-quota-cap" POST "/api/v1/kpi/evaluate-school-batch" \
  '[{"staff_id":"33333333-3333-3333-3333-333333333301","staff_name":"GV1","automated_system_score":99,"manager_assessment_score":99,"peer_review_score":99,"parent_satisfaction_score":99,"has_innovative_solution":true},
    {"staff_id":"33333333-3333-3333-3333-333333333302","staff_name":"GV2","automated_system_score":96,"manager_assessment_score":96,"peer_review_score":96,"parent_satisfaction_score":96,"has_innovative_solution":true},
    {"staff_id":"33333333-3333-3333-3333-333333333303","staff_name":"GV3","automated_system_score":93,"manager_assessment_score":93,"peer_review_score":93,"parent_satisfaction_score":93,"has_innovative_solution":true},
    {"staff_id":"33333333-3333-3333-3333-333333333304","staff_name":"GV4","automated_system_score":90,"manager_assessment_score":90,"peer_review_score":90,"parent_satisfaction_score":90,"has_innovative_solution":true},
    {"staff_id":"33333333-3333-3333-3333-333333333305","staff_name":"GV5","automated_system_score":80,"manager_assessment_score":80,"peer_review_score":80,"parent_satisfaction_score":80},
    {"staff_id":"33333333-3333-3333-3333-333333333306","staff_name":"GV6","automated_system_score":99,"manager_assessment_score":99,"peer_review_score":99,"parent_satisfaction_score":99,"has_innovative_solution":true,"has_disciplinary_action":true}]' \
  200 "BGH"
excellent_count=$(grep -o 'XUẤT SẮC' "$LOG_DIR/last-response.json" | wc -l | tr -d ' ')
if [ "$excellent_count" -eq 1 ]; then
  log_case "kpi-batch-20pct-cap" 0 0 "(đúng 1 suất xuất sắc / 5 người đạt tốt+)"
else
  log_case "kpi-batch-20pct-cap" 0 1 "(kỳ vọng 1, thực tế ${excellent_count})"
fi

# Batch rỗng -> 500 có message rõ
http_case "kpi-batch-empty" POST "/api/v1/kpi/evaluate-school-batch" '[]' 500 "BGH"

# --- 5. AI Copilot nhận xét (TT 27/2020 & TT 89/2026) ---
http_case "ai-copilot-htt" POST "/api/v1/evaluations/ai-copilot" \
  '{"subject_name":"Toán","achievement_level":"HTT","exam_score":10,"diligence_level":"T","student_name":"Minh Khôi"}' \
  200 "GVCN"
grep -q 'Minh Khôi' "$LOG_DIR/last-response.json" \
  && log_case "ai-copilot-personalized" 0 0 "(nhận xét có tên học sinh)" \
  || log_case "ai-copilot-personalized" 0 1 "nhận xét không cá nhân hóa"

http_case "ai-copilot-missing-fields" POST "/api/v1/evaluations/ai-copilot" \
  '{"subject_name":"Toán"}' 400 "GVCN"

# --- 6. Bán trú: hủy ăn trước/sau 8h ---
http_case "boarding-cancel-before-8am" POST "/api/v1/boarding/cancel-meal" \
  "{\"student_id\":\"44444444-4444-4444-4444-444444444444\",\"attendance_date\":\"$(date +%F)\",\"cancellation_hour\":7,\"meal_fee_per_day\":35000}" \
  200 "PHHS"
grep -q '"is_refund_eligible":true' "$LOG_DIR/last-response.json" \
  && grep -q '"refund_amount":35000' "$LOG_DIR/last-response.json" \
  && log_case "boarding-refund-before-8am" 0 0 "(hoàn 35000đ)" \
  || log_case "boarding-refund-before-8am" 0 1 "không hoàn tiền dù hủy trước 8h"

http_case "boarding-cancel-after-8am" POST "/api/v1/boarding/cancel-meal" \
  "{\"student_id\":\"44444444-4444-4444-4444-444444444444\",\"attendance_date\":\"$(date +%F)\",\"cancellation_hour\":9}" \
  200 "PHHS"
grep -q '"is_refund_eligible":false' "$LOG_DIR/last-response.json" \
  && log_case "boarding-no-refund-after-8am" 0 0 "(đúng quy chế bếp)" \
  || log_case "boarding-no-refund-after-8am" 0 1 "vẫn hoàn tiền sau 8h"

# --- 7. Tuyển sinh + Digital Consent (NĐ 13/2023 Điều 20) ---
http_case "student-onboard-no-consent" POST "/api/v1/students/onboard" \
  '{"student_code":"HS9001","full_name":"Nguyễn Văn A","birth_date":"2019-05-20T00:00:00Z","guardian_name":"Nguyễn Văn B","guardian_phone":"0912345678","has_digital_consent":false}' \
  422 "VAN_THU"

http_case "student-onboard-with-consent" POST "/api/v1/students/onboard" \
  '{"student_code":"HS9002","full_name":"Nguyễn Văn C","national_id":"001206012345","birth_date":"2019-06-01T00:00:00Z","guardian_name":"Nguyễn Văn D","guardian_national_id":"001075098765","guardian_phone":"0987654321","guardian_relationship":"FATHER","has_digital_consent":true,"consent_method":"VNEID_SSO"}' \
  201 "VAN_THU"
grep -q '"has_digital_consent":true' "$LOG_DIR/last-response.json" \
  && log_case "student-onboard-consent-persisted" 0 0 "(consent được ghi nhận)" \
  || log_case "student-onboard-consent-persisted" 0 1 "consent không xuất hiện trong response"

# --- 8. Đánh giá học sinh & Học bạ số (TT 27/2020 & TT 89/2026 Strategy Resolution) ---
# Đánh giá môn học thành công
http_case "eval-subject-score-success" POST "/api/v1/evaluations/subject-scores" \
  '{"student_id":"c0000000-0000-0000-0000-000000000001","term_id":"2026-2027-HK1","subject_id":"TOAN","achievement_level":"HTT","exam_score":10.0,"teacher_comment":"Hoàn thành xuất sắc môn Toán"}' \
  200 "GVBM"

# Đánh giá môn học sai mức đạt -> 400
http_case "eval-subject-score-invalid" POST "/api/v1/evaluations/subject-scores" \
  '{"student_id":"c0000000-0000-0000-0000-000000000001","term_id":"2026-2027-HK1","subject_id":"TOAN","achievement_level":"GIOI"}' \
  400 "GVBM"

# Đánh giá phẩm chất và năng lực (5 phẩm chất, 10 năng lực)
http_case "eval-behavior-success" POST "/api/v1/evaluations/behaviors" \
  '{"student_id":"c0000000-0000-0000-0000-000000000001","term_id":"2026-2027-HK1","patriotism":"T","compassion":"T","diligence":"T","honesty":"T","responsibility":"T","qualities_comment":"Chăm ngoan","autonomy_and_self_learning":"T","communication_and_collaboration":"T","problem_solving_and_creativity":"T","language_competence":"T","mathematical_competence":"T","scientific_competence":"T","technological_competence":"T","computer_competence":"T","aesthetic_competence":"T","physical_competence":"T"}' \
  200 "GVCN"

# Xuất snapshot học bạ trước 01/12/2026 -> phân giải TT27_2020
http_case "transcript-generate-tt27" POST "/api/v1/transcripts/generate" \
  '{"student_id":"c0000000-0000-0000-0000-000000000001","academic_year_id":"2026-2027","grade_level":5,"evaluation_date":"2026-10-25"}' \
  200 "BGH"
grep -q '"regulation_applied":"TT27_2020"' "$LOG_DIR/last-response.json" \
  && log_case "transcript-tt27-resolved" 0 0 "(phân giải đúng TT 27/2020 trước 01/12/2026)" \
  || log_case "transcript-tt27-resolved" 0 1 "không phân giải TT 27/2020"

# Xuất snapshot học bạ từ 01/12/2026 trở đi -> phân giải TT89_2026
http_case "transcript-generate-tt89" POST "/api/v1/transcripts/generate" \
  '{"student_id":"c0000000-0000-0000-0000-000000000001","academic_year_id":"2026-2027","grade_level":5,"evaluation_date":"2026-12-05"}' \
  200 "BGH"
grep -q '"regulation_applied":"TT89_2026"' "$LOG_DIR/last-response.json" \
  && log_case "transcript-tt89-resolved" 0 0 "(phân giải đúng TT 89/2026 từ 01/12/2026)" \
  || log_case "transcript-tt89-resolved" 0 1 "không phân giải TT 89/2026"

# --- 9. SUPER_ADMIN bypass RBAC ---
http_case "super-admin-bypass" POST "/api/v1/kpi/evaluate-individual" \
  '{"staff_id":"55555555-5555-5555-5555-555555555555","staff_name":"Admin test","automated_system_score":80,"manager_assessment_score":80,"peer_review_score":80,"parent_satisfaction_score":80}' \
  200 "SUPER_ADMIN"

# ---------------------------------------------------------------- REPORT
echo
echo "== KẾT QUẢ =="
echo "PASS: ${PASS}   FAIL: ${FAIL}"
if [ "$FAIL" -gt 0 ]; then
  echo "Các ca fail:"
  for c in "${FAILED_CASES[@]}"; do echo "  - $c"; done
fi
echo "Server log:     ${SERVER_LOG}"
echo "Request log:    ${REQ_LOG}"
echo "Last response:  ${LOG_DIR}/last-response.json"

[ "$FAIL" -eq 0 ]
