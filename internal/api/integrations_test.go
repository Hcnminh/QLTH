package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"qlth/internal/integrations/sandbox"

	"github.com/google/uuid"
)

func TestIntegrations_SmartCASign_Success(t *testing.T) {
	router := SetupRouter()

	reqPayload := sandbox.SmartCASignRequest{
		DocumentID:   "d1000000-0000-0000-0000-000000000001",
		DocumentType: "TRANSCRIPT_SNAPSHOT",
		SignerID:     "b0000000-0000-0000-0000-000000000001",
		SignerName:   "Trần Thị Bình",
		SignerRole:   "GVCN",
	}
	body, _ := json.Marshal(reqPayload)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/integrations/smartca/sign", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", uuid.New().String())
	req.Header.Set("X-User-Role", "GVCN")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d, body: %s", w.Code, w.Body.String())
	}

	var res map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to parse json response: %v", err)
	}

	if res["success"] != true {
		t.Errorf("expected success true")
	}
	data, ok := res["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected data object in response")
	}

	if data["environment"] != "SANDBOX_STUB" {
		t.Errorf("expected SANDBOX_STUB environment, got: %v", data["environment"])
	}
	if data["signature_format"] != "CMS/PKCS#7" {
		t.Errorf("expected CMS/PKCS#7 signature format, got: %v", data["signature_format"])
	}
	cmsSig, ok := data["cms_pkcs7_signature"].(string)
	if !ok || len(cmsSig) < 20 {
		t.Errorf("expected valid simulated cms_pkcs7_signature")
	}
	tsaToken, ok := data["tsa_token"].(map[string]interface{})
	if !ok || tsaToken["tsa_provider"] == nil {
		t.Errorf("expected valid tsa_token object")
	}
}

func TestIntegrations_VietQRGenerate_Success(t *testing.T) {
	router := SetupRouter()

	reqPayload := sandbox.VietQRGenerateRequest{
		BankBIN:       "970422",
		AccountNumber: "0123456789",
		AccountName:   "Truong Tieu Hoc EduPrimary",
		Amount:        550000,
		Description:   "TIEN AN BAN TRU THANG 11",
	}
	body, _ := json.Marshal(reqPayload)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/integrations/vietqr/generate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", uuid.New().String())
	req.Header.Set("X-User-Role", "KE_TOAN")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d, body: %s", w.Code, w.Body.String())
	}

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	data := res["data"].(map[string]interface{})

	if data["environment"] != "SANDBOX_STUB" {
		t.Errorf("expected SANDBOX_STUB, got: %v", data["environment"])
	}
	qrStr, ok := data["qr_string"].(string)
	if !ok || len(qrStr) < 20 {
		t.Errorf("expected non-empty EMVCo qr_string")
	}
	imgURL, ok := data["image_url"].(string)
	if !ok || len(imgURL) < 10 {
		t.Errorf("expected non-empty image_url")
	}
}

func TestIntegrations_MOETSync_Success(t *testing.T) {
	router := SetupRouter()

	reqPayload := sandbox.MOETSyncRequest{
		SchoolID:     "00000000-0000-0000-0000-000000000001",
		AcademicYear: "2026-2027",
		SyncType:     "ALL",
		BatchSize:    20,
	}
	body, _ := json.Marshal(reqPayload)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/integrations/moet/sync", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", "00000000-0000-0000-0000-000000000001")
	req.Header.Set("X-User-Role", "BGH")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d, body: %s", w.Code, w.Body.String())
	}

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	data := res["data"].(map[string]interface{})

	if data["environment"] != "SANDBOX_STUB" {
		t.Errorf("expected SANDBOX_STUB, got: %v", data["environment"])
	}
	if data["status"] != "COMPLETED" {
		t.Errorf("expected COMPLETED status, got: %v", data["status"])
	}
	if data["moet_transaction_id"] == nil || data["moet_transaction_id"] == "" {
		t.Errorf("expected valid moet_transaction_id")
	}
}

func TestIntegrations_MOETQueue_Success(t *testing.T) {
	router := SetupRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/integrations/moet/queue", nil)
	req.Header.Set("X-School-ID", "00000000-0000-0000-0000-000000000001")
	req.Header.Set("X-User-Role", "VAN_THU")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d, body: %s", w.Code, w.Body.String())
	}

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	data := res["data"].(map[string]interface{})

	if data["environment"] != "SANDBOX_STUB" {
		t.Errorf("expected SANDBOX_STUB, got: %v", data["environment"])
	}
	if data["status"] != "HEALTHY" {
		t.Errorf("expected HEALTHY status, got: %v", data["status"])
	}
}

func TestIntegrations_VNeIDVerify_Success(t *testing.T) {
	router := SetupRouter()

	reqPayload := sandbox.VNeIDVerificationRequest{
		NationalID:     "001202012345",
		FullName:       "Nguyễn Văn An",
		DateOfBirth:    "2018-05-15",
		Gender:         "MALE",
		ConsentGranted: true,
	}
	body, _ := json.Marshal(reqPayload)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/integrations/vneid/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", uuid.New().String())
	req.Header.Set("X-User-Role", "VAN_THU")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d, body: %s", w.Code, w.Body.String())
	}

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	data := res["data"].(map[string]interface{})

	if data["environment"] != "SANDBOX_STUB" {
		t.Errorf("expected SANDBOX_STUB, got: %v", data["environment"])
	}
	if data["verification_level"] != "LEVEL_2" {
		t.Errorf("expected LEVEL_2, got: %v", data["verification_level"])
	}
	if data["status"] != "VERIFIED" {
		t.Errorf("expected status VERIFIED, got: %v", data["status"])
	}
}

func TestIntegrations_VNeIDVerify_MissingConsent(t *testing.T) {
	router := SetupRouter()

	reqPayload := sandbox.VNeIDVerificationRequest{
		NationalID:     "001202012345",
		FullName:       "Nguyễn Văn An",
		ConsentGranted: false, // Decree 13/2023 violation
	}
	body, _ := json.Marshal(reqPayload)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/integrations/vneid/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", uuid.New().String())
	req.Header.Set("X-User-Role", "VAN_THU")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422 Unprocessable Entity for missing consent, got %d", w.Code)
	}
}

func TestIntegrations_EInvoiceIssue_Success(t *testing.T) {
	router := SetupRouter()

	reqPayload := sandbox.EInvoiceIssueRequest{
		SchoolID:     "00000000-0000-0000-0000-000000000001",
		BuyerName:    "Trần Văn Hùng",
		StudentName:  "Trần Bảo Nam",
		StudentClass: "5A1",
		Items: []sandbox.InvoiceItem{
			{
				Name:      "Tiền ăn bán trú tháng 11/2026",
				Unit:      "Tháng",
				Quantity:  1,
				UnitPrice: 660000,
				VATRate:   0,
			},
		},
	}
	body, _ := json.Marshal(reqPayload)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/integrations/einvoice/issue", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", "00000000-0000-0000-0000-000000000001")
	req.Header.Set("X-User-Role", "KE_TOAN")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d, body: %s", w.Code, w.Body.String())
	}

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	data := res["data"].(map[string]interface{})

	if data["environment"] != "SANDBOX_STUB" {
		t.Errorf("expected SANDBOX_STUB, got: %v", data["environment"])
	}
	if data["invoice_series"] != "1C26TAA" {
		t.Errorf("expected Decree 123 series 1C26TAA, got: %v", data["invoice_series"])
	}
	if data["status"] != "ISSUED" {
		t.Errorf("expected ISSUED status, got: %v", data["status"])
	}
}

func TestIntegrations_SandboxStatus(t *testing.T) {
	router := SetupRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/integrations/sandbox/status", nil)
	req.Header.Set("X-School-ID", uuid.New().String())
	req.Header.Set("X-User-Role", "BGH")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	data := res["data"].(map[string]interface{})

	if data["environment"] != "SANDBOX_STUB" {
		t.Errorf("expected SANDBOX_STUB environment, got: %v", data["environment"])
	}
	integrations, ok := data["integrations"].([]interface{})
	if !ok || len(integrations) < 5 {
		t.Errorf("expected at least 5 sandbox integrations, got: %d", len(integrations))
	}
}

func TestIntegrations_RBAC_Forbidden(t *testing.T) {
	router := SetupRouter()

	// BEP_TRUONG is not allowed to sync MOET
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/integrations/moet/sync", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-School-ID", uuid.New().String())
	req.Header.Set("X-User-Role", "BEP_TRUONG")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden, got %d", w.Code)
	}
}
