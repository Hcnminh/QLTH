package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"qlth/internal/core/domain"
	"qlth/internal/core/kpi"
	"qlth/internal/core/teaching"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestHealthCheck(t *testing.T) {
	router := SetupRouter()

	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}
}

func TestCalculateTeachingQuota_API(t *testing.T) {
	router := SetupRouter()

	profile := teaching.TeacherWorkloadProfile{
		TeacherID:         "GV001",
		TeacherName:       "Nguyễn Thị Hoa",
		SchoolType:        "PRIMARY",
		IsHomeroom:        true, // Reduces 4 hours/week (Circular 05/2025)
		WeeksInSchoolYear: 35,
		ActualTaughtHours: 720,
	}
	body, _ := json.Marshal(profile)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/teaching/calculate-quota", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", uuid.New().String())
	req.Header.Set("X-User-Role", "BGH")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	data := res["data"].(map[string]interface{})

	// Net quota = 805 - (4 * 35) = 665
	if int(data["net_annual_quota"].(float64)) != 665 {
		t.Fatalf("expected net quota 665, got %v", data["net_annual_quota"])
	}
	// Overtime = 720 - 665 = 55
	if int(data["overtime_hours"].(float64)) != 55 {
		t.Fatalf("expected overtime 55, got %v", data["overtime_hours"])
	}
}

func TestEvaluateKPI_API(t *testing.T) {
	router := SetupRouter()

	sub := kpi.StaffKPISubmission{
		StaffID:                 uuid.New(),
		StaffName:               "Lê Văn Cường",
		AutomatedSystemScore:    92.0,
		ManagerAssessmentScore:  90.0,
		PeerReviewScore:         90.0,
		ParentSatisfactionScore: 90.0,
		HasInnovativeSolution:   true,
	}
	body, _ := json.Marshal(sub)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/kpi/evaluate-individual", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", uuid.New().String())
	req.Header.Set("X-User-Role", "BGH")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	data := res["data"].(map[string]interface{})

	if data["official_rank"] != string(kpi.RankExcellent) {
		t.Fatalf("expected rank %s, got %v", kpi.RankExcellent, data["official_rank"])
	}
}

func TestAIPedagogicalCopilot_API(t *testing.T) {
	router := SetupRouter()

	payload := map[string]interface{}{
		"subject_name":      "Toán",
		"achievement_level": "HTT",
		"exam_score":        10.0,
		"diligence_level":   "T",
		"student_name":      "Minh Khôi",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/evaluations/ai-copilot", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", uuid.New().String())
	req.Header.Set("X-User-Role", "GVCN")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	suggestions := res["suggestions"].([]interface{})
	if len(suggestions) < 2 {
		t.Fatalf("expected at least 2 suggestions, got %d", len(suggestions))
	}
}

func TestCancelBoardingMeal_API(t *testing.T) {
	router := SetupRouter()

	// Case 1: Cancel before 8:00 AM (Hour 7) -> Refund eligible
	payloadRefund := map[string]interface{}{
		"student_id":        uuid.New().String(),
		"attendance_date":   "2026-10-28",
		"cancellation_hour": 7,
		"meal_fee_per_day":  35000,
	}
	body1, _ := json.Marshal(payloadRefund)

	req1, _ := http.NewRequest(http.MethodPost, "/api/v1/boarding/cancel-meal", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("X-School-ID", uuid.New().String())
	req1.Header.Set("X-User-Role", "PHHS")

	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w1.Code)
	}

	var res1 map[string]interface{}
	json.Unmarshal(w1.Body.Bytes(), &res1)
	data1 := res1["data"].(map[string]interface{})
	if data1["is_refund_eligible"] != true {
		t.Fatalf("expected refund eligible true for 7 AM cancellation")
	}

	// Case 2: Cancel after 8:00 AM (Hour 9) -> Refund NOT eligible
	payloadNoRefund := map[string]interface{}{
		"student_id":        uuid.New().String(),
		"attendance_date":   "2026-10-28",
		"cancellation_hour": 9,
		"meal_fee_per_day":  35000,
	}
	body2, _ := json.Marshal(payloadNoRefund)

	req2, _ := http.NewRequest(http.MethodPost, "/api/v1/boarding/cancel-meal", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("X-School-ID", uuid.New().String())
	req2.Header.Set("X-User-Role", "PHHS")

	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	var res2 map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &res2)
	data2 := res2["data"].(map[string]interface{})
	if data2["is_refund_eligible"] != false {
		t.Fatalf("expected refund eligible false for 9 AM cancellation")
	}
}

func TestStudentOnboarding_DigitalConsentValidation(t *testing.T) {
	router := SetupRouter()

	// Case: Without Digital Consent (Decree 13/2023) -> Should be rejected with 422
	studentWithoutConsent := domain.Student{
		FullName:          "Nguyễn Hoàng",
		NationalID:        "001206123456",
		StudentCode:       "HS1001",
		HasDigitalConsent: false,
	}
	body, _ := json.Marshal(studentWithoutConsent)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/students/onboard", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", uuid.New().String())
	req.Header.Set("X-User-Role", "VAN_THU")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422 Unprocessable Entity, got %d", w.Code)
	}
}

func TestRBAC_UnauthorizedRole(t *testing.T) {
	router := SetupRouter()

	// A Parent (PHHS) trying to call school batch KPI evaluation
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/kpi/evaluate-school-batch", bytes.NewBuffer([]byte("[]")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", uuid.New().String())
	req.Header.Set("X-User-Role", "PHHS") // Insufficient privilege

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden for Parent accessing KPI batch, got %d", w.Code)
	}
}

func TestAdminEndpoints_Integration(t *testing.T) {
	router := SetupRouter()

	// 1. VAL-ADM-002: GET /admin/login returns 200 with AdminLTE login interface
	t.Run("GET /admin/login returns 200 with login fields", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/admin/login", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 OK from /admin/login, got %d", w.Code)
		}
		body := w.Body.String()
		if !strings.Contains(body, "username") || !strings.Contains(body, "password") {
			t.Fatalf("expected login fields in html, got: %s", body)
		}
	})

	// 2. Unauthenticated GET /admin redirects to /admin/login
	t.Run("GET /admin redirects to /admin/login", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/admin", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusFound {
			t.Fatalf("expected 302 redirect from /admin, got %d", w.Code)
		}
		loc := w.Header().Get("Location")
		if !strings.Contains(loc, "/admin/login") {
			t.Fatalf("expected redirect to /admin/login, got: %s", loc)
		}
	})

	// 3. VAL-ADM-003: POST /admin/signin with valid admin credentials
	t.Run("POST /admin/signin authenticates and redirects to /admin", func(t *testing.T) {
		form := url.Values{}
		form.Add("username", "admin")
		form.Add("password", "admin")

		req, _ := http.NewRequest(http.MethodPost, "/admin/signin", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 OK from /admin/signin, got %d (body: %s)", w.Code, w.Body.String())
		}

		var resp struct {
			Code int                    `json:"code"`
			Msg  string                 `json:"msg"`
			Data map[string]interface{} `json:"data"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode signin response: %v", err)
		}
		if resp.Code != 200 {
			t.Fatalf("expected code 200, got %d", resp.Code)
		}
		if resp.Data["url"] != "/admin" {
			t.Fatalf("expected redirect to /admin, got %v", resp.Data["url"])
		}

		// Check session cookie
		hasCookie := false
		for _, c := range w.Result().Cookies() {
			if c.Name == "go_admin_session" {
				hasCookie = true
				break
			}
		}
		if !hasCookie {
			t.Fatal("expected go_admin_session cookie to be returned")
		}
	})
}
