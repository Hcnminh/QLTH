package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"qlth/internal/core/kpi"
)

// BGH gọi batch evaluation hợp lệ → 200 và áp trần 20% xuất sắc
func TestBatchEvaluateKPI_Success(t *testing.T) {
	router := SetupRouter()

	subs := make([]kpi.StaffKPISubmission, 5)
	for i := range subs {
		score := 90.0 + float64(i)
		subs[i] = kpi.StaffKPISubmission{
			StaffID:                 uuid.New(),
			StaffName:               "Giáo viên",
			AutomatedSystemScore:    score,
			ManagerAssessmentScore:  score,
			PeerReviewScore:         score,
			ParentSatisfactionScore: score,
			HasInnovativeSolution:   true,
		}
	}
	body, _ := json.Marshal(subs)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/kpi/evaluate-school-batch", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", uuid.New().String())
	req.Header.Set("X-User-Role", "BGH")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var res struct {
		Success bool `json:"success"`
		Data    []struct {
			OfficialRank string `json:"official_rank"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if !res.Success {
		t.Fatal("expected success=true")
	}
	if len(res.Data) != 5 {
		t.Fatalf("expected 5 evaluations, got %d", len(res.Data))
	}

	// Trần 20% của 5 người đạt tốt trở lên = tối đa 1 suất xuất sắc
	excellent := 0
	for _, d := range res.Data {
		if d.OfficialRank == string(kpi.RankExcellent) {
			excellent++
		}
	}
	if excellent != 1 {
		t.Fatalf("expected exactly 1 excellent (20%% cap of 5), got %d", excellent)
	}
}

// Danh sách rỗng → 500 với message rõ ràng (không panic)
func TestBatchEvaluateKPI_EmptyList(t *testing.T) {
	router := SetupRouter()

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/kpi/evaluate-school-batch", bytes.NewBuffer([]byte("[]")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", uuid.New().String())
	req.Header.Set("X-User-Role", "BGH")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 for empty batch, got %d", w.Code)
	}
}

// JSON body sai cấu trúc → 400, không crash
func TestBatchEvaluateKPI_InvalidBody(t *testing.T) {
	router := SetupRouter()

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/kpi/evaluate-school-batch", bytes.NewBuffer([]byte(`{"not":"an array"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", uuid.New().String())
	req.Header.Set("X-User-Role", "BGH")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for malformed body, got %d", w.Code)
	}
}

// Hồ sơ học sinh hợp lệ kèm Digital Consent → 201 và server tự sinh UUID
func TestStudentOnboarding_SuccessWithConsent(t *testing.T) {
	router := SetupRouter()

	payload := map[string]interface{}{
		"student_code":          "HS2001",
		"full_name":             "Trần Bảo Ngọc",
		"national_id":           "001206999888",
		"birth_date":            "2019-05-20T00:00:00Z",
		"guardian_name":         "Trần Văn Phúc",
		"guardian_national_id":  "001075012345",
		"guardian_phone":        "0912345678",
		"guardian_relationship": "FATHER",
		"has_digital_consent":   true,
		"consent_method":        "VNEID_SSO",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/students/onboard", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", uuid.New().String())
	req.Header.Set("X-User-Role", "VAN_THU")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d. Body: %s", w.Code, w.Body.String())
	}

	var res struct {
		Success bool `json:"success"`
		Data    struct {
			ID          string `json:"id"`
			StudentCode string `json:"student_code"`
			Consent     bool   `json:"has_digital_consent"`
			CreatedAt   string `json:"created_at"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if res.Data.ID == "" {
		t.Fatal("server must auto-generate student UUID")
	}
	if res.Data.StudentCode != "HS2001" {
		t.Fatalf("expected student_code HS2001, got %s", res.Data.StudentCode)
	}
	if !res.Data.Consent {
		t.Fatal("consent flag must be preserved")
	}
}

// Thiếu header X-School-ID → 401 (multi-tenancy bắt buộc)
func TestRBAC_MissingSchoolID(t *testing.T) {
	router := SetupRouter()

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/kpi/evaluate-individual", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Role", "BGH")
	// Không set X-School-ID

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without X-School-ID, got %d", w.Code)
	}
}

// SUPER_ADMIN bypass mọi kiểm tra role
func TestRBAC_SuperAdminBypass(t *testing.T) {
	router := SetupRouter()

	sub := kpi.StaffKPISubmission{
		StaffID:                 uuid.New(),
		AutomatedSystemScore:    80,
		ManagerAssessmentScore:  80,
		PeerReviewScore:         80,
		ParentSatisfactionScore: 80,
	}
	body, _ := json.Marshal(sub)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/kpi/evaluate-individual", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", uuid.New().String())
	req.Header.Set("X-User-Role", "SUPER_ADMIN")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected SUPER_ADMIN to pass RBAC, got %d", w.Code)
	}
}

// Thiếu role header → mặc định GUEST và bị từ chối 403
func TestRBAC_MissingRoleDefaultsToGuest(t *testing.T) {
	router := SetupRouter()

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/evaluations/ai-copilot", bytes.NewBuffer([]byte(`{"subject_name":"Toán","achievement_level":"HTT","student_name":"Test"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", uuid.New().String())
	// Không set X-User-Role → GUEST

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for GUEST, got %d", w.Code)
	}
}

// Hủy suất ăn thiếu trường bắt buộc → 400
func TestCancelBoardingMeal_InvalidPayload(t *testing.T) {
	router := SetupRouter()

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/boarding/cancel-meal", bytes.NewBuffer([]byte(`{"student_id":"not-a-uuid"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", uuid.New().String())
	req.Header.Set("X-User-Role", "PHHS")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid UUID, got %d", w.Code)
	}
}

// Mặc định tiền ăn 35000 VNĐ khi không truyền meal_fee_per_day
func TestCancelBoardingMeal_DefaultMealFee(t *testing.T) {
	router := SetupRouter()

	payload := map[string]interface{}{
		"student_id":        uuid.New().String(),
		"attendance_date":   "2026-10-28",
		"cancellation_hour": 6, // trước 8h → được hoàn tiền
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/boarding/cancel-meal", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", uuid.New().String())
	req.Header.Set("X-User-Role", "PHHS")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var res struct {
		Data struct {
			IsRefundEligible bool    `json:"is_refund_eligible"`
			RefundAmount     float64 `json:"refund_amount"`
		} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &res)

	if !res.Data.IsRefundEligible {
		t.Fatal("cancellation at 6AM must be refund eligible")
	}
	if res.Data.RefundAmount != 35000 {
		t.Fatalf("expected default meal fee 35000, got %f", res.Data.RefundAmount)
	}
}
