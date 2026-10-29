package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"qlth/internal/db"

	"github.com/google/uuid"
)

func TestLessonPlan_SubmitAndApproveWorkflow(t *testing.T) {
	router := SetupRouter()
	schoolID := "a0000000-0000-0000-0000-000000000001"
	teacherID := "b0000000-0000-0000-0000-000000000001"
	deptHeadID := "b0000000-0000-0000-0000-000000000002"
	principalID := "b0000000-0000-0000-0000-000000000003"

	// 1. Submit Lesson Plan (CV 2345) -> 201 Created
	submitBody := map[string]interface{}{
		"teacher_id":          teacherID,
		"teacher_name":        "Nguyễn Văn An",
		"grade_level":         5,
		"subject_id":          "TOAN",
		"title":               "Bài 25: Hình thang và diện tích hình thang",
		"week_number":         18,
		"duration_slots":      2,
		"curriculum_standard": "CV_2345",
		"objectives": map[string]string{
			"knowledge":    "Nắm đặc điểm hình thang và công thức diện tích S = (a + b) x h / 2",
			"competencies": "Năng lực mô hình hóa và giải quyết vấn đề toán học",
			"qualities":    "Chăm chỉ, cẩn thận",
		},
		"equipment": "Mô hình hình thang cắt ghép, thước kẻ, ê-ke",
		"activities": []map[string]interface{}{
			{"phase": "Khởi động", "duration_mins": 5, "content": "Tìm hình thang trong thực tế"},
			{"phase": "Khám phá", "duration_mins": 15, "content": "Ghép hai hình thang thành hình bình hành"},
			{"phase": "Luyện tập", "duration_mins": 15, "content": "Tính diện tích thửa ruộng hình thang"},
			{"phase": "Vận dụng", "duration_mins": 5, "content": "Ứng dụng làm mái nhà đồ chơi"},
		},
	}
	bodyBytes, _ := json.Marshal(submitBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/lesson-plans/submit", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", schoolID)
	req.Header.Set("X-User-Role", "GVBM")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created on submit, got %d. Body: %s", w.Code, w.Body.String())
	}

	var submitResp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &submitResp)
	if submitResp["success"] != true {
		t.Fatalf("expected success: true, got %v", submitResp["success"])
	}
	planData := submitResp["data"].(map[string]interface{})
	planID := planData["id"].(string)
	if planData["current_status"] != "SUBMITTED" {
		t.Errorf("expected initial status SUBMITTED, got %s", planData["current_status"])
	}

	// 2. Multi-level Approval: Step 1 - Tổ trưởng chuyên môn duyệt
	deptHeadReview := map[string]interface{}{
		"action":        "APPROVE",
		"reviewer_id":   deptHeadID,
		"reviewer_role": "TO_TRUONG",
		"reviewer_name": "Trần Thị Bình",
		"comments":      "Giáo án đúng chuẩn CV 2345, chia thời gian phù hợp. Duyệt chuyển BGH.",
	}
	dhBytes, _ := json.Marshal(deptHeadReview)

	reqApprove1, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/lesson-plans/%s/approve", planID), bytes.NewBuffer(dhBytes))
	reqApprove1.Header.Set("Content-Type", "application/json")
	reqApprove1.Header.Set("X-School-ID", schoolID)
	reqApprove1.Header.Set("X-User-Role", "TO_TRUONG")

	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, reqApprove1)

	if w1.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on dept head approval, got %d. Body: %s", w1.Code, w1.Body.String())
	}

	var approve1Resp map[string]interface{}
	_ = json.Unmarshal(w1.Body.Bytes(), &approve1Resp)
	p1Data := approve1Resp["data"].(map[string]interface{})
	if p1Data["current_status"] != "DEPARTMENT_HEAD_APPROVED" {
		t.Errorf("expected DEPARTMENT_HEAD_APPROVED, got %s", p1Data["current_status"])
	}

	// 3. Multi-level Approval: Step 2 - Ban Giám hiệu ký duyệt số hoàn tất -> APPROVED
	bghReview := map[string]interface{}{
		"action":            "APPROVE",
		"reviewer_id":       principalID,
		"reviewer_role":     "BGH",
		"reviewer_name":     "Lê Hoàng Cường",
		"comments":          "Ký duyệt chính thức cho phép giảng dạy.",
		"digital_signature": "SMARTCA_BGH_SIGN_TSA_OK_12345",
	}
	bghBytes, _ := json.Marshal(bghReview)

	reqApprove2, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/lesson-plans/%s/approve", planID), bytes.NewBuffer(bghBytes))
	reqApprove2.Header.Set("Content-Type", "application/json")
	reqApprove2.Header.Set("X-School-ID", schoolID)
	reqApprove2.Header.Set("X-User-Role", "BGH")

	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, reqApprove2)

	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on principal approval, got %d. Body: %s", w2.Code, w2.Body.String())
	}

	var approve2Resp map[string]interface{}
	_ = json.Unmarshal(w2.Body.Bytes(), &approve2Resp)
	p2Data := approve2Resp["data"].(map[string]interface{})
	if p2Data["current_status"] != "APPROVED" {
		t.Errorf("expected status APPROVED, got %s", p2Data["current_status"])
	}

	// 4. GET /api/v1/lesson-plans/:id
	reqGet, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/lesson-plans/%s", planID), nil)
	reqGet.Header.Set("X-School-ID", schoolID)
	reqGet.Header.Set("X-User-Role", "GVBM")

	wGet := httptest.NewRecorder()
	router.ServeHTTP(wGet, reqGet)

	if wGet.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on GET, got %d", wGet.Code)
	}

	// 5. GET /api/v1/lesson-plans
	reqList, _ := http.NewRequest(http.MethodGet, "/api/v1/lesson-plans", nil)
	reqList.Header.Set("X-School-ID", schoolID)
	reqList.Header.Set("X-User-Role", "BGH")

	wList := httptest.NewRecorder()
	router.ServeHTTP(wList, reqList)

	if wList.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on list, got %d", wList.Code)
	}
}

func TestLessonPlan_RejectionWorkflowAndValidation(t *testing.T) {
	router := SetupRouter()
	schoolID := "a0000000-0000-0000-0000-000000000001"
	teacherID := "b0000000-0000-0000-0000-000000000001"

	// 1. Submit invalid grade level -> 400 Bad Request
	invalidGrade := map[string]interface{}{
		"teacher_id":  teacherID,
		"grade_level": 7, // Invalid for primary school
		"subject_id":  "TOAN",
		"title":       "Bài tập lớp 7",
		"week_number": 5,
	}
	body, _ := json.Marshal(invalidGrade)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/lesson-plans/submit", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", schoolID)
	req.Header.Set("X-User-Role", "GVBM")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid grade level, got %d", w.Code)
	}

	// 2. Submit valid plan then reject without comment -> 400 Bad Request
	validSubmit := map[string]interface{}{
		"teacher_id":  teacherID,
		"grade_level": 3,
		"subject_id":  "TIENG_VIET",
		"title":       "Bài 8: Ôn tập chữ hoa",
		"week_number": 8,
	}
	body, _ = json.Marshal(validSubmit)
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/lesson-plans/submit", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", schoolID)
	req.Header.Set("X-User-Role", "GVBM")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", w.Code)
	}

	var res map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &res)
	planID := res["data"].(map[string]interface{})["id"].(string)

	// Reject without comments -> 400
	rejectNoComment := map[string]interface{}{
		"action":        "REJECT",
		"reviewer_role": "TO_TRUONG",
		"comments":      "   ",
	}
	body, _ = json.Marshal(rejectNoComment)
	req, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/lesson-plans/%s/approve", planID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", schoolID)
	req.Header.Set("X-User-Role", "TO_TRUONG")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 when rejecting without comments, got %d", w.Code)
	}

	// Reject with comments -> 200 OK and status REJECTED
	rejectWithComment := map[string]interface{}{
		"action":        "REJECT",
		"reviewer_role": "TO_TRUONG",
		"comments":      "Thiếu thiết bị mẫu chữ hoa theo CV 2345. Yêu cầu soạn lại mục Khám phá.",
	}
	body, _ = json.Marshal(rejectWithComment)
	req, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/lesson-plans/%s/approve", planID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", schoolID)
	req.Header.Set("X-User-Role", "TO_TRUONG")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on rejection, got %d", w.Code)
	}

	var rejectResp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &rejectResp)
	if rejectResp["data"].(map[string]interface{})["current_status"] != "REJECTED" {
		t.Errorf("expected status REJECTED, got %s", rejectResp["data"].(map[string]interface{})["current_status"])
	}
}

func TestTimetable_ConfigureSlotsAndConflictDetection(t *testing.T) {
	cleanupDB := func() {
		database, _ := db.GetDB()
		if database != nil {
			_, _ = database.Exec("DELETE FROM timetable_slots WHERE id NOT IN ('e0000000-0000-0000-0000-000000000201', 'e0000000-0000-0000-0000-000000000202', 'e0000000-0000-0000-0000-000000000203')")
		}
	}
	cleanupDB()
	t.Cleanup(cleanupDB)

	router := SetupRouter()
	schoolID := "a0000000-0000-0000-0000-000000000001"
	class1ID := "f0000000-0000-0000-0000-000000000001"
	class2ID := "f0000000-0000-0000-0000-000000000002"
	teacherID := "b0000000-0000-0000-0000-000000000001"

	// 1. Configure slot 1: Teacher b01 in Class 5A1 on Thứ 3, Tiết 1 -> 201 Created
	slot1 := map[string]interface{}{
		"academic_year_id": "2026-2027",
		"term_id":          "2026-2027-HK1",
		"class_id":         class2ID,
		"class_name":       "5A1",
		"teacher_id":       teacherID,
		"teacher_name":     "Nguyễn Văn An",
		"subject_id":       "TOAN",
		"day_of_week":      3, // Thứ 3
		"period":           1, // Tiết 1
		"room_number":      "Phòng 5A1",
		"week_number":      0,
	}
	body, _ := json.Marshal(slot1)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/timetable/slots", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", schoolID)
	req.Header.Set("X-User-Role", "BGH")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for slot 1, got %d. Body: %s", w.Code, w.Body.String())
	}

	// 2. Configure conflicting slot: SAME teacher in Class 1A1 at exact same Thứ 3, Tiết 1 -> 409 Conflict
	conflictingTeacherSlot := map[string]interface{}{
		"academic_year_id": "2026-2027",
		"term_id":          "2026-2027-HK1",
		"class_id":         class1ID,
		"class_name":       "1A1",
		"teacher_id":       teacherID,
		"teacher_name":     "Nguyễn Văn An",
		"subject_id":       "TOAN",
		"day_of_week":      3,
		"period":           1,
		"room_number":      "Phòng 1A1",
		"week_number":      0,
	}
	body, _ = json.Marshal(conflictingTeacherSlot)
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/timetable/slots", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", schoolID)
	req.Header.Set("X-User-Role", "BGH")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409 Conflict for teacher overlap, got %d. Body: %s", w.Code, w.Body.String())
	}

	var conflictResp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &conflictResp)
	if conflictResp["success"] != false {
		t.Errorf("expected success: false on conflict, got %v", conflictResp["success"])
	}

	// 3. Configure valid slot 2: Same teacher, different period (Thứ 3, Tiết 2) -> 201 Created
	slot2 := map[string]interface{}{
		"academic_year_id": "2026-2027",
		"term_id":          "2026-2027-HK1",
		"class_id":         class2ID,
		"class_name":       "5A1",
		"teacher_id":       teacherID,
		"teacher_name":     "Nguyễn Văn An",
		"subject_id":       "TOAN",
		"day_of_week":      3,
		"period":           2, // Tiết 2
		"room_number":      "Phòng 5A1",
		"week_number":      0,
	}
	body, _ = json.Marshal(slot2)
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/timetable/slots", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", schoolID)
	req.Header.Set("X-User-Role", "GVCN")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for slot 2, got %d. Body: %s", w.Code, w.Body.String())
	}

	// 4. GET /api/v1/timetable/slots
	reqGet, _ := http.NewRequest(http.MethodGet, "/api/v1/timetable/slots?academic_year_id=2026-2027&term_id=2026-2027-HK1", nil)
	reqGet.Header.Set("X-School-ID", schoolID)
	reqGet.Header.Set("X-User-Role", "PHHS")

	wGet := httptest.NewRecorder()
	router.ServeHTTP(wGet, reqGet)
	if wGet.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for GET timetable slots, got %d. Body: %s", wGet.Code, wGet.Body.String())
	}
}

func TestLessonPlan_RBAC(t *testing.T) {
	router := SetupRouter()
	schoolID := uuid.New().String()

	// 1. Missing X-School-ID -> 401
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/lesson-plans/submit", bytes.NewBufferString("{}"))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got %d", w.Code)
	}

	// 2. PHHS (Parent) attempting to submit lesson plan -> 403 Forbidden
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/lesson-plans/submit", bytes.NewBufferString("{}"))
	req.Header.Set("X-School-ID", schoolID)
	req.Header.Set("X-User-Role", "PHHS")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 Forbidden for PHHS, got %d", w.Code)
	}

	// 3. GVBM attempting to approve lesson plan -> 403 Forbidden
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/lesson-plans/e0000000-0000-0000-0000-000000000001/approve", bytes.NewBufferString("{}"))
	req.Header.Set("X-School-ID", schoolID)
	req.Header.Set("X-User-Role", "GVBM")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 Forbidden for GVBM approval, got %d", w.Code)
	}
}
