package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGetBoardingMenus_Success(t *testing.T) {
	router := SetupRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/boarding/menus", nil)
	req.Header.Set("X-School-ID", "a0000000-0000-0000-0000-000000000001")
	req.Header.Set("X-User-Role", "BGH")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Success bool `json:"success"`
		Data    []struct {
			ID                  string  `json:"id"`
			MealType            string  `json:"meal_type"`
			TotalKcal           float64 `json:"total_kcal"`
			NutritionalAnalysis struct {
				IsCompliant       bool   `json:"is_compliant"`
				GuidelineStandard string `json:"guideline_standard"`
				Breakdown         struct {
					TotalKcal           float64 `json:"total_kcal"`
					ProteinPercent      float64 `json:"protein_percent"`
					LipidPercent        float64 `json:"lipid_percent"`
					GlucidPercent       float64 `json:"glucid_percent"`
					EnergyBalanceStatus string  `json:"energy_balance_status"`
					PLGBalanceStatus    string  `json:"plg_balance_status"`
				} `json:"breakdown"`
			} `json:"nutritional_analysis"`
		} `json:"data"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	if !resp.Success {
		t.Errorf("Expected success to be true")
	}
	if len(resp.Data) == 0 {
		t.Fatalf("Expected at least 1 boarding menu")
	}

	firstMenu := resp.Data[0]
	if firstMenu.TotalKcal <= 0 {
		t.Errorf("Expected TotalKcal > 0, got %f", firstMenu.TotalKcal)
	}
	if firstMenu.NutritionalAnalysis.GuidelineStandard == "" {
		t.Errorf("Expected GuidelineStandard to be populated")
	}
	if firstMenu.NutritionalAnalysis.Breakdown.TotalKcal <= 0 {
		t.Errorf("Expected Breakdown.TotalKcal > 0")
	}
}

func TestCreateBoardingMenu_Success(t *testing.T) {
	router := SetupRouter()

	payload := map[string]interface{}{
		"menu_date": time.Now().Format("2006-01-02"),
		"meal_type": "LUNCH",
		"dishes": []map[string]interface{}{
			{"name": "Thịt kho tàu", "weight_g": 120, "protein_g": 20, "lipid_g": 18, "glucid_g": 4},
			{"name": "Canh cải thịt băm", "weight_g": 150, "protein_g": 5, "lipid_g": 2, "glucid_g": 8},
			{"name": "Cơm trắng", "weight_g": 180, "protein_g": 5, "lipid_g": 0.5, "glucid_g": 75},
		},
		"sample_retention_verified": true,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/boarding/menus", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", "a0000000-0000-0000-0000-000000000001")
	req.Header.Set("X-User-Role", "KE_TOAN")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestVerifyFoodSampleRetention_Success(t *testing.T) {
	router := SetupRouter()

	menuID := uuid.New()
	payload := map[string]interface{}{
		"sampled_at":          time.Now().Add(-26 * time.Hour).Format(time.RFC3339),
		"storage_temperature": 4.0,
		"keeper_name":         "Nguyễn Thị Y Tá",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/boarding/menus/"+menuID.String()+"/sample-retention", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", "a0000000-0000-0000-0000-000000000001")
	req.Header.Set("X-User-Role", "Y_TE")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetStudentHealthProfile_Success(t *testing.T) {
	router := SetupRouter()

	studentID := "c0000000-0000-0000-0000-000000000001"
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/students/"+studentID+"/health", nil)
	req.Header.Set("X-School-ID", "a0000000-0000-0000-0000-000000000001")
	req.Header.Set("X-User-Role", "GVCN")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Success bool `json:"success"`
		Data    struct {
			StudentID      string  `json:"student_id"`
			HeightCm       float64 `json:"height_cm"`
			WeightKg       float64 `json:"weight_kg"`
			BMI            float64 `json:"bmi"`
			BMICategory    string  `json:"bmi_category"`
			EyesightStatus string  `json:"eyesight_status"`
			MedicalEvents  []any   `json:"medical_events"`
		} `json:"data"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	if !resp.Success {
		t.Errorf("Expected success to be true")
	}
	if resp.Data.BMI <= 0 {
		t.Errorf("Expected BMI > 0, got %f", resp.Data.BMI)
	}
	if resp.Data.BMICategory == "" {
		t.Errorf("Expected BMICategory to be populated")
	}
	if resp.Data.EyesightStatus == "" {
		t.Errorf("Expected EyesightStatus to be populated")
	}
}

func TestRecordStudentHealth_AndPrescriptionWorkflow(t *testing.T) {
	router := SetupRouter()
	studentID := uuid.New().String()

	// 1. Record physical health checkup
	checkupPayload := map[string]interface{}{
		"checkup_date":       "2026-09-20",
		"height_cm":          125.0,
		"weight_kg":          24.5,
		"left_eye_vision":    "10/10",
		"right_eye_vision":   "9/10",
		"dental_condition":   "Răng sữa thay bình thường",
		"chronic_conditions": "Không dị ứng",
	}
	cBody, _ := json.Marshal(checkupPayload)

	reqCheckup, _ := http.NewRequest(http.MethodPost, "/api/v1/students/"+studentID+"/health", bytes.NewReader(cBody))
	reqCheckup.Header.Set("Content-Type", "application/json")
	reqCheckup.Header.Set("X-School-ID", "a0000000-0000-0000-0000-000000000001")
	reqCheckup.Header.Set("X-User-Role", "Y_TE")

	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, reqCheckup)

	if w1.Code != http.StatusOK {
		t.Fatalf("Expected status 200 for health checkup, got %d: %s", w1.Code, w1.Body.String())
	}

	// 2. Submit prescription delivery request
	rxPayload := map[string]interface{}{
		"medication_name":     "Siro ho Prospan",
		"dosage":              "5ml sau ăn trưa",
		"administration_time": "12:30",
		"instructions":        "Uống trực tiếp",
		"parent_notes":        "Nhờ cô cho cháu uống đúng giờ",
	}
	rxBody, _ := json.Marshal(rxPayload)

	reqRx, _ := http.NewRequest(http.MethodPost, "/api/v1/students/"+studentID+"/health/prescriptions", bytes.NewReader(rxBody))
	reqRx.Header.Set("Content-Type", "application/json")
	reqRx.Header.Set("X-School-ID", "a0000000-0000-0000-0000-000000000001")
	reqRx.Header.Set("X-User-Role", "PHHS")

	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, reqRx)

	if w2.Code != http.StatusCreated {
		t.Fatalf("Expected status 201 for prescription request, got %d: %s", w2.Code, w2.Body.String())
	}

	var rxResp struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	_ = json.Unmarshal(w2.Body.Bytes(), &rxResp)
	rxID := rxResp.Data.ID

	// 3. Nurse administers prescription
	adminBody, _ := json.Marshal(map[string]string{"nurse_name": "Y tá Trưởng"})
	reqAdmin, _ := http.NewRequest(http.MethodPost, "/api/v1/students/"+studentID+"/health/prescriptions/"+rxID+"/administer", bytes.NewReader(adminBody))
	reqAdmin.Header.Set("Content-Type", "application/json")
	reqAdmin.Header.Set("X-School-ID", "a0000000-0000-0000-0000-000000000001")
	reqAdmin.Header.Set("X-User-Role", "Y_TE")

	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, reqAdmin)

	if w3.Code != http.StatusOK {
		t.Fatalf("Expected status 200 for administer prescription, got %d: %s", w3.Code, w3.Body.String())
	}

	// 4. Record medical event
	evPayload := map[string]interface{}{
		"event_type":      "SƠ_CẤP_CỨU",
		"symptoms":        "Chảy máu cam nhẹ",
		"treatment":       "Chườm lạnh và ấn cánh mũi 10 phút",
		"nurse_name":      "Y tá Trưởng",
		"parent_notified": true,
	}
	evBody, _ := json.Marshal(evPayload)

	reqEv, _ := http.NewRequest(http.MethodPost, "/api/v1/students/"+studentID+"/health/medical-events", bytes.NewReader(evBody))
	reqEv.Header.Set("Content-Type", "application/json")
	reqEv.Header.Set("X-School-ID", "a0000000-0000-0000-0000-000000000001")
	reqEv.Header.Set("X-User-Role", "Y_TE")

	w4 := httptest.NewRecorder()
	router.ServeHTTP(w4, reqEv)

	if w4.Code != http.StatusCreated {
		t.Fatalf("Expected status 201 for medical event, got %d: %s", w4.Code, w4.Body.String())
	}

	// 5. Query final profile
	reqFinal, _ := http.NewRequest(http.MethodGet, "/api/v1/students/"+studentID+"/health", nil)
	reqFinal.Header.Set("X-School-ID", "a0000000-0000-0000-0000-000000000001")
	reqFinal.Header.Set("X-User-Role", "GVCN")

	w5 := httptest.NewRecorder()
	router.ServeHTTP(w5, reqFinal)

	if w5.Code != http.StatusOK {
		t.Fatalf("Expected status 200 for health profile, got %d: %s", w5.Code, w5.Body.String())
	}
}

func TestHealthEndpoints_RBAC(t *testing.T) {
	router := SetupRouter()

	// Missing X-School-ID -> 401
	req1, _ := http.NewRequest(http.MethodGet, "/api/v1/boarding/menus", nil)
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	if w1.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 for missing school id, got %d", w1.Code)
	}

	// Forbidden role for POST /api/v1/students/:id/health -> 403
	req2, _ := http.NewRequest(http.MethodPost, "/api/v1/students/c0000000-0000-0000-0000-000000000001/health", bytes.NewReader([]byte("{}")))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("X-School-ID", "a0000000-0000-0000-0000-000000000001")
	req2.Header.Set("X-User-Role", "PHHS") // Parents cannot alter clinical health checkup
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	if w2.Code != http.StatusForbidden {
		t.Errorf("Expected 403 for unauthorized role, got %d", w2.Code)
	}
}
