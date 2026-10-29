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

# --- 10. Kế hoạch bài dạy (CV 2345/STEM TT 86) & Thời khóa biểu (TT 05/2025) ---
# Nộp giáo án chuẩn CV 2345 -> 201 Created
http_case "lesson-plan-submit-cv2345" POST "/api/v1/lesson-plans/submit" \
  '{"teacher_id":"b0000000-0000-0000-0000-000000000001","teacher_name":"Nguyễn Văn An","grade_level":5,"subject_id":"TOAN","title":"Bài 10: Phân số thập phân chuẩn CV 2345","week_number":10,"duration_slots":2,"curriculum_standard":"CV_2345","objectives":{"knowledge":"Hiểu khái niệm phân số thập phân","competencies":"Năng lực tư duy toán học","qualities":"Chăm chỉ"},"equipment":"Bộ thẻ số thập phân"}' \
  201 "GVBM"
grep -q '"current_status":"SUBMITTED"' "$LOG_DIR/last-response.json" \
  && log_case "lesson-plan-submitted-status" 0 0 "(trạng thái SUBMITTED)" \
  || log_case "lesson-plan-submitted-status" 0 1 "trạng thái không phải SUBMITTED"

LP_ID=$(python3 -c 'import json; d=json.load(open("'$LOG_DIR'/last-response.json")); print(d.get("data",{}).get("id",""))')

# Duyệt giáo án bước 1: Tổ trưởng chuyên môn duyệt -> DEPARTMENT_HEAD_APPROVED
http_case "lesson-plan-approve-dept-head" POST "/api/v1/lesson-plans/${LP_ID}/approve" \
  '{"action":"APPROVE","reviewer_role":"TO_TRUONG","reviewer_name":"Trần Thị Bình","comments":"Giáo án đạt chuẩn 4 bước CV 2345. Duyệt chuyển BGH."}' \
  200 "TO_TRUONG"
grep -q '"current_status":"DEPARTMENT_HEAD_APPROVED"' "$LOG_DIR/last-response.json" \
  && log_case "lesson-plan-dept-head-status" 0 0 "(duyệt cấp tổ thành công)" \
  || log_case "lesson-plan-dept-head-status" 0 1 "trạng thái không phải DEPARTMENT_HEAD_APPROVED"

# Duyệt giáo án bước 2: Ban Giám hiệu ký số duyệt hoàn tất -> APPROVED
http_case "lesson-plan-approve-principal" POST "/api/v1/lesson-plans/${LP_ID}/approve" \
  '{"action":"APPROVE","reviewer_role":"BGH","reviewer_name":"Lê Hoàng Cường","comments":"Ký duyệt chính thức theo quy chế chuyên môn.","digital_signature":"SIGN_SMARTCA_BGH_TSA_OK"}' \
  200 "BGH"
grep -q '"current_status":"APPROVED"' "$LOG_DIR/last-response.json" \
  && log_case "lesson-plan-principal-status" 0 0 "(BGH ký duyệt hoàn tất APPROVED)" \
  || log_case "lesson-plan-principal-status" 0 1 "trạng thái không phải APPROVED"

# Cấu hình tiết thời khóa biểu hợp lệ -> 201 Created
http_case "timetable-slot-configure-valid" POST "/api/v1/timetable/slots" \
  '{"academic_year_id":"2026-2027","term_id":"2026-2027-HK1","class_id":"f0000000-0000-0000-0000-000000000002","class_name":"5A1","teacher_id":"b0000000-0000-0000-0000-000000000001","teacher_name":"Nguyễn Văn An","subject_id":"TOAN","day_of_week":4,"period":1,"room_number":"Phòng 5A1","week_number":0}' \
  201 "BGH"

# Cấu hình trùng lịch giáo viên (Thứ 4, Tiết 1 cùng giáo viên ở lớp khác) -> 409 Conflict (TT 05/2025)
http_case "timetable-slot-teacher-conflict" POST "/api/v1/timetable/slots" \
  '{"academic_year_id":"2026-2027","term_id":"2026-2027-HK1","class_id":"f0000000-0000-0000-0000-000000000001","class_name":"1A1","teacher_id":"b0000000-0000-0000-0000-000000000001","teacher_name":"Nguyễn Văn An","subject_id":"TOAN","day_of_week":4,"period":1,"room_number":"Phòng 1A1","week_number":0}' \
  409 "BGH"
grep -q 'Thông tư 05/2025/TT-BGDĐT' "$LOG_DIR/last-response.json" \
  && log_case "timetable-conflict-citation" 0 0 "(báo lỗi xung đột lịch dạy TT 05/2025 chuẩn xác)" \
  || log_case "timetable-conflict-citation" 0 1 "không viện dẫn Thông tư 05/2025 trong lỗi xung đột"

# --- 11. Bán trú, Dinh dưỡng Viện Dinh Dưỡng, Lưu mẫu 24h & Hồ sơ Y tế (BMI) ---
# Lấy danh sách thực đơn bán trú kèm phân tích dinh dưỡng chuẩn Viện Dinh Dưỡng
http_case "boarding-menus-nutrition-breakdown" GET "/api/v1/boarding/menus" \
  "-" 200 "BGH"
grep -q 'Viện Dinh dưỡng Quốc gia' "$LOG_DIR/last-response.json" \
  && log_case "boarding-nutrition-guidelines" 0 0 "(phân tích dinh dưỡng chuẩn Viện Dinh Dưỡng)" \
  || log_case "boarding-nutrition-guidelines" 0 1 "không có phân tích chuẩn Viện Dinh Dưỡng"

# Lập thực đơn bán trú mới tự động tính Kcal & cân bằng P:L:G -> 201 Created
http_case "boarding-menu-create-plan" POST "/api/v1/boarding/menus" \
  '{"menu_date":"2026-10-30","meal_type":"LUNCH","dishes":[{"name":"Thịt heo hầm đậu hũ","weight_g":120,"protein_g":22,"lipid_g":18,"glucid_g":5},{"name":"Canh rau ngót tôm tươi","weight_g":150,"protein_g":6,"lipid_g":2,"glucid_g":8},{"name":"Cơm tấm gạo tám","weight_g":180,"protein_g":5,"lipid_g":0.5,"glucid_g":75}],"sample_retention_verified":true}' \
  201 "KE_TOAN"
MENU_ID=$(python3 -c 'import json; d=json.load(open("'$LOG_DIR'/last-response.json")); print(d.get("data",{}).get("id",""))')

# Xác nhận kiểm tra lưu mẫu thức ăn 24h theo QĐ 1246/QĐ-BYT -> 200 OK
http_case "boarding-sample-retention-verify" POST "/api/v1/boarding/menus/${MENU_ID}/sample-retention" \
  '{"sampled_at":"2026-10-29T11:00:00Z","storage_temperature":4.0,"keeper_name":"Đặng Thị Mai (Y tế bán trú)"}' \
  200 "Y_TE"
grep -q '1246/QĐ-BYT' "$LOG_DIR/last-response.json" \
  && log_case "boarding-sample-retention-rule" 0 0 "(kiểm tra lưu mẫu 24h QĐ 1246/QĐ-BYT đạt)" \
  || log_case "boarding-sample-retention-rule" 0 1 "không xác nhận QĐ 1246/QĐ-BYT"

# Lấy hồ sơ sức khỏe điện tử học sinh (BMI, thị lực, nha, bệnh mãn tính) -> 200 OK
http_case "student-health-profile-get" GET "/api/v1/students/c0000000-0000-0000-0000-000000000001/health" \
  "-" 200 "GVCN"
grep -q '"bmi":' "$LOG_DIR/last-response.json" && grep -q '"bmi_category":' "$LOG_DIR/last-response.json" \
  && log_case "student-health-bmi-calculated" 0 0 "(tính toán BMI và phân loại thể trạng)" \
  || log_case "student-health-bmi-calculated" 0 1 "thiếu chỉ số BMI hoặc phân loại"

# Phụ huynh gửi đơn thuốc nhờ cô y tế cho uống tại trường -> 201 Created
http_case "student-health-prescription-add" POST "/api/v1/students/c0000000-0000-0000-0000-000000000001/health/prescriptions" \
  '{"medication_name":"Siro ho Astex","dosage":"5ml sau ăn trưa","administration_time":"12:30","instructions":"Uống cùng nước ấm","parent_notes":"Cháu ho dị ứng thời tiết"}' \
  201 "PHHS"

# Ghi nhật ký y tế và xử trí sơ cấp cứu tại phòng y tế -> 201 Created
http_case "student-health-medical-event-add" POST "/api/v1/students/c0000000-0000-0000-0000-000000000001/health/medical-events" \
  '{"event_type":"SƠ_CẤP_CỨU","symptoms":"Trầy da nhẹ ở đầu gối do chạy ngã","treatment":"Rửa nước muối sinh lý, sát khuẩn Betadine và dán gạc cá nhân","parent_notified":true}' \
  201 "Y_TE"

# --- 12. Cổng tích hợp ngoại vi Sandbox Stubs (VAL-DOGFOOD-005) ---
# Ký số tập trung SmartCA kèm dấu thời gian TSA -> 200 OK (CMS/PKCS#7 & TSA token)
http_case "integration-smartca-sign" POST "/api/v1/integrations/smartca/sign" \
  '{"document_id":"d1000000-0000-0000-0000-000000000001","document_type":"TRANSCRIPT_SNAPSHOT","signer_id":"b0000000-0000-0000-0000-000000000001","signer_name":"Trần Thị Bình","signer_role":"GVCN"}' \
  200 "GVCN"
grep -q '"environment":"SANDBOX_STUB"' "$LOG_DIR/last-response.json" && grep -q '"signature_format":"CMS/PKCS#7"' "$LOG_DIR/last-response.json" && grep -q '"tsa_provider":' "$LOG_DIR/last-response.json" \
  && log_case "smartca-cms-tsa-stub" 0 0 "(chữ ký số SmartCA CMS/PKCS#7 và dấu thời gian TSA)" \
  || log_case "smartca-cms-tsa-stub" 0 1 "thiếu chữ ký CMS hoặc dấu thời gian TSA"

# Sinh mã VietQR động chuẩn EMVCo và URL ảnh -> 200 OK
http_case "integration-vietqr-generate" POST "/api/v1/integrations/vietqr/generate" \
  '{"bank_bin":"970422","account_number":"0123456789","account_name":"Truong Tieu Hoc EduPrimary","amount":550000,"description":"TIEN AN BAN TRU THANG 11"}' \
  200 "KE_TOAN"
grep -q '"environment":"SANDBOX_STUB"' "$LOG_DIR/last-response.json" && grep -q '"qr_string":"000201010212' "$LOG_DIR/last-response.json" && grep -q 'img.vietqr.io' "$LOG_DIR/last-response.json" \
  && log_case "vietqr-emvco-dynamic" 0 0 "(mã VietQR động chuẩn EMVCo và URL ảnh thanh toán)" \
  || log_case "vietqr-emvco-dynamic" 0 1 "chuỗi QR không đúng chuẩn EMVCo hoặc thiếu URL ảnh"

# Đồng bộ hàng đợi Outbox lên CSDL Ngành Giáo dục MOET -> 200 OK
http_case "integration-moet-sync" POST "/api/v1/integrations/moet/sync" \
  '{"school_id":"00000000-0000-0000-0000-000000000001","academic_year":"2026-2027","sync_type":"ALL","batch_size":20}' \
  200 "BGH"
grep -q '"environment":"SANDBOX_STUB"' "$LOG_DIR/last-response.json" && grep -q '"moet_transaction_id":' "$LOG_DIR/last-response.json" && grep -q '"status":"COMPLETED"' "$LOG_DIR/last-response.json" \
  && log_case "moet-outbox-sync-batch" 0 0 "(xử lý hàng đợi Outbox đồng bộ CSDL MOET thành công)" \
  || log_case "moet-outbox-sync-batch" 0 1 "đồng bộ Outbox MOET thất bại"

# Xác thực định danh điện tử VNeID Mức 2 (NĐ 69/2024 & NĐ 13/2023) -> 200 OK
http_case "integration-vneid-verify-lvl2" POST "/api/v1/integrations/vneid/verify" \
  '{"national_id":"001202012345","full_name":"Nguyễn Văn An","date_of_birth":"2018-05-15","gender":"MALE","consent_granted":true}' \
  200 "VAN_THU"
grep -q '"environment":"SANDBOX_STUB"' "$LOG_DIR/last-response.json" && grep -q '"verification_level":"LEVEL_2"' "$LOG_DIR/last-response.json" && grep -q '"biometric_matched":true' "$LOG_DIR/last-response.json" \
  && log_case "vneid-level2-verified" 0 0 "(xác thực VNeID Mức độ 2 sinh trắc học thành công)" \
  || log_case "vneid-level2-verified" 0 1 "xác thực VNeID Level 2 thất bại"

# Mô phỏng phát hành hóa đơn điện tử Nghị định 123/2020/NĐ-CP -> 200 OK
http_case "integration-einvoice-issue" POST "/api/v1/integrations/einvoice/issue" \
  '{"school_id":"00000000-0000-0000-0000-000000000001","buyer_name":"Trần Văn Hùng","student_name":"Trần Bảo Nam","student_class":"5A1","items":[{"name":"Tiền ăn bán trú tháng 11/2026","unit":"Tháng","quantity":1,"unit_price":660000,"vat_rate":0}]}' \
  200 "KE_TOAN"
grep -q '"environment":"SANDBOX_STUB"' "$LOG_DIR/last-response.json" && grep -q '"invoice_series":"1C26TAA"' "$LOG_DIR/last-response.json" && grep -q '"status":"ISSUED"' "$LOG_DIR/last-response.json" \
  && log_case "einvoice-decree123-issued" 0 0 "(phát hành hóa đơn điện tử NĐ 123 ký hiệu 1C26TAA)" \
  || log_case "einvoice-decree123-issued" 0 1 "phát hành hóa đơn điện tử thất bại"

# Kiểm tra trạng thái toàn bộ 5 cổng kết nối sandbox -> 200 OK
http_case "integration-sandbox-status-all" GET "/api/v1/integrations/sandbox/status" \
  "-" 200 "BGH"
grep -q '"environment":"SANDBOX_STUB"' "$LOG_DIR/last-response.json" && grep -q 'ACTIVE_SANDBOX' "$LOG_DIR/last-response.json" \
  && log_case "sandbox-status-verified" 0 0 "(toàn bộ 5 sandbox stubs hoạt động ACTIVE_SANDBOX)" \
  || log_case "sandbox-status-verified" 0 1 "trạng thái sandbox stubs không hợp lệ"

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
