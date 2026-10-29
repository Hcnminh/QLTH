package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"qlth/internal/core/assessment"
	"qlth/internal/core/domain"

	"github.com/google/uuid"
)

func TestRecordSubjectScore_Success(t *testing.T) {
	router := SetupRouter()

	score9 := 9.0
	bodyData := map[string]interface{}{
		"student_id":        "c0000000-0000-0000-0000-000000000001",
		"term_id":           "2026-2027-HK1",
		"subject_id":        "TOAN",
		"achievement_level": "HTT",
		"exam_score":        score9,
		"teacher_comment":   "Tiếp thu bài nhanh, làm bài cẩn thận",
		"evaluation_date":   "2026-10-20",
	}
	body, _ := json.Marshal(bodyData)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/evaluations/subject-scores", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", uuid.New().String())
	req.Header.Set("X-User-Role", "GVBM")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d. Body: %s", w.Code, w.Body.String())
	}

	var res map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}
	if res["success"] != true {
		t.Fatalf("expected success: true, got %v", res["success"])
	}
	data := res["data"].(map[string]interface{})
	if data["regulation_applied"] != "TT27_2020" {
		t.Fatalf("expected TT27_2020 for 2026-10-20, got %v", data["regulation_applied"])
	}
}

func TestRecordSubjectScore_InvalidLevel(t *testing.T) {
	router := SetupRouter()

	bodyData := map[string]interface{}{
		"student_id":        uuid.New().String(),
		"term_id":           "2026-2027-HK1",
		"subject_id":        "TOAN",
		"achievement_level": "XUAT_SAC", // Không hợp lệ theo TT 27/TT 89 (phải là HTT, HT, CHT)
	}
	body, _ := json.Marshal(bodyData)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/evaluations/subject-scores", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", uuid.New().String())
	req.Header.Set("X-User-Role", "GVBM")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestRecordSubjectScore_ExamScoreOutOfRange(t *testing.T) {
	router := SetupRouter()

	bodyData := map[string]interface{}{
		"student_id":        uuid.New().String(),
		"term_id":           "2026-2027-HK1",
		"subject_id":        "TOAN",
		"achievement_level": "HT",
		"exam_score":        11.5, // Thang điểm kiểm tra là 1.0 - 10.0
	}
	body, _ := json.Marshal(bodyData)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/evaluations/subject-scores", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", uuid.New().String())
	req.Header.Set("X-User-Role", "GVBM")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", w.Code)
	}
}

func TestRecordSubjectScore_RBAC_Forbidden(t *testing.T) {
	router := SetupRouter()

	bodyData := map[string]interface{}{
		"student_id":        uuid.New().String(),
		"term_id":           "2026-2027-HK1",
		"subject_id":        "TOAN",
		"achievement_level": "HT",
	}
	body, _ := json.Marshal(bodyData)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/evaluations/subject-scores", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", uuid.New().String())
	req.Header.Set("X-User-Role", "KE_TOAN") // Kế toán không được nhập điểm môn học

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden, got %d", w.Code)
	}
}

func TestRecordBehaviorEvaluation_Success(t *testing.T) {
	router := SetupRouter()

	bodyData := map[string]interface{}{
		"student_id":                      "c0000000-0000-0000-0000-000000000001",
		"term_id":                         "2026-2027-HK1",
		"patriotism":                      "T",
		"compassion":                      "T",
		"diligence":                       "T",
		"honesty":                         "T",
		"responsibility":                  "T",
		"qualities_comment":               "Gương mẫu, có tinh thần tương thân tương ái",
		"autonomy_and_self_learning":      "T",
		"communication_and_collaboration": "T",
		"problem_solving_and_creativity":  "T",
		"language_competence":             "T",
		"mathematical_competence":         "T",
		"evaluation_date":                 "2026-10-20",
	}
	body, _ := json.Marshal(bodyData)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/evaluations/behaviors", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", uuid.New().String())
	req.Header.Set("X-User-Role", "GVCN")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d. Body: %s", w.Code, w.Body.String())
	}

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	if res["success"] != true {
		t.Fatalf("expected success: true, got %v", res["success"])
	}
}

func TestRecordBehaviorEvaluation_InvalidQualityLevel(t *testing.T) {
	router := SetupRouter()

	bodyData := map[string]interface{}{
		"student_id":     uuid.New().String(),
		"term_id":        "2026-2027-HK1",
		"patriotism":     "X", // Không hợp lệ
		"compassion":     "T",
		"diligence":      "T",
		"honesty":        "T",
		"responsibility": "T",
	}
	body, _ := json.Marshal(bodyData)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/evaluations/behaviors", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", uuid.New().String())
	req.Header.Set("X-User-Role", "GVCN")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request for invalid quality level, got %d", w.Code)
	}
}

func TestRecordBehaviorEvaluation_RBAC_Forbidden(t *testing.T) {
	router := SetupRouter()

	bodyData := map[string]interface{}{
		"student_id":     uuid.New().String(),
		"term_id":        "2026-2027-HK1",
		"patriotism":     "T",
		"compassion":     "T",
		"diligence":      "T",
		"honesty":        "T",
		"responsibility": "T",
	}
	body, _ := json.Marshal(bodyData)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/evaluations/behaviors", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", uuid.New().String())
	req.Header.Set("X-User-Role", "GVBM") // GVBM không được đánh giá phẩm chất năng lực học sinh (chỉ BGH/GVCN)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden, got %d", w.Code)
	}
}

func TestGenerateTranscript_Circular27_Resolution(t *testing.T) {
	router := SetupRouter()

	score10 := 10.0
	score9 := 9.0
	studentID := uuid.New()
	reqBody := map[string]interface{}{
		"student_id":       studentID.String(),
		"academic_year_id": "2026-2027",
		"grade_level":      5,
		"evaluation_date":  "2026-10-25", // Trước 01/12/2026 -> phải phân giải TT27_2020
		"student": domain.Student{
			ID:          studentID,
			StudentCode: "HS_TT27",
			FullName:    "Lê Hoàng Long",
		},
		"subjects": []domain.StudentSubjectEvaluation{
			{SubjectID: "TOAN", AchievementLevel: "HTT", ExamScore: &score10},
			{SubjectID: "TIENG_VIET", AchievementLevel: "HTT", ExamScore: &score9},
		},
		"behavior": domain.StudentBehaviorEvaluation{
			Patriotism:               "T",
			Compassion:               "T",
			Diligence:                "T",
			Honesty:                  "T",
			Responsibility:           "T",
			AutonomySelfLearning:     "T",
			CommunicationCooperation: "T",
			ProblemSolvingCreativity: "T",
		},
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/transcripts/generate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", uuid.New().String())
	req.Header.Set("X-User-Role", "BGH")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d. Body: %s", w.Code, w.Body.String())
	}

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	data := res["data"].(map[string]interface{})

	if data["regulation_applied"] != "TT27_2020" {
		t.Fatalf("expected TT27_2020 before 01/12/2026, got %v", data["regulation_applied"])
	}
	if data["annual_rank"] != string(assessment.RankStudentExcellent) {
		t.Fatalf("expected %s, got %v", assessment.RankStudentExcellent, data["annual_rank"])
	}
	if data["annual_reward"] != "Học sinh Xuất sắc" {
		t.Fatalf("expected 'Học sinh Xuất sắc', got %v", data["annual_reward"])
	}
	if data["promotion_status"] != "PROMOTED" {
		t.Fatalf("expected 'PROMOTED', got %v", data["promotion_status"])
	}
}

func TestGenerateTranscript_Circular89_Resolution(t *testing.T) {
	router := SetupRouter()

	score10 := 10.0
	studentID := uuid.New()
	reqBody := map[string]interface{}{
		"student_id":       studentID.String(),
		"academic_year_id": "2026-2027",
		"grade_level":      5,
		"evaluation_date":  "2026-12-05", // Từ 01/12/2026 trở đi -> phải phân giải TT89_2026
		"student": domain.Student{
			ID:                studentID,
			StudentCode:       "HS_TT89",
			NationalID:        "001206099999",
			FullName:          "Nguyễn Thị Lan",
			HasDigitalConsent: true,
		},
		"subjects": []domain.StudentSubjectEvaluation{
			{SubjectID: "TOAN", AchievementLevel: "HTT", ExamScore: &score10, IsComputerBasedTest: true},
			{SubjectID: "TIENG_VIET", AchievementLevel: "HTT", ExamScore: &score10, IsCentralizedExam: true},
		},
		"behavior": domain.StudentBehaviorEvaluation{
			Patriotism:               "T",
			Compassion:               "T",
			Diligence:                "T",
			Honesty:                  "T",
			Responsibility:           "T",
			AutonomySelfLearning:     "T",
			CommunicationCooperation: "T",
			ProblemSolvingCreativity: "T",
		},
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/transcripts/generate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", uuid.New().String())
	req.Header.Set("X-User-Role", "GVCN")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d. Body: %s", w.Code, w.Body.String())
	}

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	data := res["data"].(map[string]interface{})

	if data["regulation_applied"] != "TT89_2026" {
		t.Fatalf("expected TT89_2026 on or after 01/12/2026, got %v", data["regulation_applied"])
	}
	snapshotData := data["snapshot_data"].(map[string]interface{})
	if snapshotData["regulation"] != "TT89_2026" {
		t.Fatalf("snapshot_data regulation must be TT89_2026, got %v", snapshotData["regulation"])
	}
	if snapshotData["digital_transcript_status"] != "OFFICIAL_DIGITAL_TRANSCRIPT_CIRCULAR_89" {
		t.Fatalf("expected official digital transcript circular 89 status, got %v", snapshotData["digital_transcript_status"])
	}
}

func TestGenerateTranscript_InclusiveEducation(t *testing.T) {
	router := SetupRouter()

	studentID := uuid.New()
	reqBody := map[string]interface{}{
		"student_id":       studentID.String(),
		"academic_year_id": "2026-2027",
		"grade_level":      4,
		"evaluation_date":  "2026-12-10",
		"student": domain.Student{
			ID:                   studentID,
			StudentCode:          "HS_INCL",
			FullName:             "Phan Minh Tuấn",
			IsInclusiveEducation: true, // Học sinh khuyết tật hòa nhập
		},
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/transcripts/generate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", uuid.New().String())
	req.Header.Set("X-User-Role", "BGH")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d. Body: %s", w.Code, w.Body.String())
	}

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	data := res["data"].(map[string]interface{})

	if data["annual_rank"] != string(assessment.RankStudentSatisfactory) {
		t.Fatalf("expected Satisfactory for inclusive education student, got %v", data["annual_rank"])
	}
	if data["promotion_status"] != "PROMOTED" {
		t.Fatalf("expected PROMOTED, got %v", data["promotion_status"])
	}
}

func TestGenerateTranscript_RBAC_Forbidden(t *testing.T) {
	router := SetupRouter()

	reqBody := map[string]interface{}{
		"student_id":       uuid.New().String(),
		"academic_year_id": "2026-2027",
		"grade_level":      5,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/transcripts/generate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", uuid.New().String())
	req.Header.Set("X-User-Role", "PHHS") // Phụ huynh không có quyền xuất học bạ chính thức

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden, got %d", w.Code)
	}
}

func TestGenerateTranscript_FullDatabaseFlow(t *testing.T) {
	router := SetupRouter()
	studentID := "c0000000-0000-0000-0000-000000000002"
	schoolID := uuid.New().String()

	// 1. Post TOAN evaluation
	score10 := 10.0
	toanData := map[string]interface{}{
		"student_id":        studentID,
		"term_id":           "2026-2027-HK1",
		"subject_id":        "TOAN",
		"achievement_level": "HTT",
		"exam_score":        score10,
		"teacher_comment":   "Xuất sắc môn Toán",
	}
	toanBody, _ := json.Marshal(toanData)
	req1, _ := http.NewRequest(http.MethodPost, "/api/v1/evaluations/subject-scores", bytes.NewBuffer(toanBody))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("X-School-ID", schoolID)
	req1.Header.Set("X-User-Role", "GVBM")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	if w1.Code != http.StatusOK {
		t.Fatalf("expected 200 for subject evaluation, got %d. Body: %s", w1.Code, w1.Body.String())
	}

	// 2. Post TIENG_VIET evaluation
	score95 := 9.5
	tvData := map[string]interface{}{
		"student_id":        studentID,
		"term_id":           "2026-2027-HK1",
		"subject_id":        "TIENG_VIET",
		"achievement_level": "HTT",
		"exam_score":        score95,
		"teacher_comment":   "Chữ viết đẹp, tiếp thu nhanh",
	}
	tvBody, _ := json.Marshal(tvData)
	req2, _ := http.NewRequest(http.MethodPost, "/api/v1/evaluations/subject-scores", bytes.NewBuffer(tvBody))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("X-School-ID", schoolID)
	req2.Header.Set("X-User-Role", "GVBM")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200 for TIENG_VIET evaluation, got %d", w2.Code)
	}

	// 3. Post Behavior evaluation
	behData := map[string]interface{}{
		"student_id":                      studentID,
		"term_id":                         "2026-2027-HK1",
		"patriotism":                      "T",
		"compassion":                      "T",
		"diligence":                       "T",
		"honesty":                         "T",
		"responsibility":                  "T",
		"qualities_comment":               "Gương mẫu trong lớp",
		"autonomy_and_self_learning":      "T",
		"communication_and_collaboration": "T",
		"problem_solving_and_creativity":  "T",
		"language_competence":             "T",
		"mathematical_competence":         "T",
		"scientific_competence":           "T",
		"technological_competence":        "T",
		"computer_competence":             "T",
		"aesthetic_competence":            "T",
		"physical_competence":             "T",
	}
	behBody, _ := json.Marshal(behData)
	req3, _ := http.NewRequest(http.MethodPost, "/api/v1/evaluations/behaviors", bytes.NewBuffer(behBody))
	req3.Header.Set("Content-Type", "application/json")
	req3.Header.Set("X-School-ID", schoolID)
	req3.Header.Set("X-User-Role", "GVCN")
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)
	if w3.Code != http.StatusOK {
		t.Fatalf("expected 200 for behavior evaluation, got %d", w3.Code)
	}

	// 4. Generate Transcript reading from DB
	genData := map[string]interface{}{
		"student_id":       studentID,
		"academic_year_id": "2026-2027",
		"grade_level":      5,
		"evaluation_date":  "2026-10-25",
	}
	genBody, _ := json.Marshal(genData)
	req4, _ := http.NewRequest(http.MethodPost, "/api/v1/transcripts/generate", bytes.NewBuffer(genBody))
	req4.Header.Set("Content-Type", "application/json")
	req4.Header.Set("X-School-ID", schoolID)
	req4.Header.Set("X-User-Role", "BGH")
	w4 := httptest.NewRecorder()
	router.ServeHTTP(w4, req4)
	if w4.Code != http.StatusOK {
		t.Fatalf("expected 200 for transcript generation, got %d. Body: %s", w4.Code, w4.Body.String())
	}

	var res map[string]interface{}
	json.Unmarshal(w4.Body.Bytes(), &res)
	data := res["data"].(map[string]interface{})
	if data["regulation_applied"] != "TT27_2020" {
		t.Fatalf("expected TT27_2020, got %v", data["regulation_applied"])
	}
	if data["annual_rank"] != string(assessment.RankStudentExcellent) {
		t.Fatalf("expected %s, got %v", assessment.RankStudentExcellent, data["annual_rank"])
	}
	if data["is_locked"] != true {
		t.Fatalf("transcript snapshot must be locked (immutable)")
	}
}
