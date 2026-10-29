package admin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupTestAdmin(t *testing.T) *gin.Engine {
	t.Helper()
	r := gin.New()
	eng, err := InitAdmin(r)
	if err != nil {
		t.Fatalf("Failed to initialize GoAdmin: %v", err)
	}
	if eng == nil {
		t.Fatal("Expected engine not to be nil")
	}
	return r
}

// VAL-ADM-002: Login screen and redirect tests
func TestAdminLoginInterface(t *testing.T) {
	r := setupTestAdmin(t)

	// GET /admin/login returns HTTP 200 with AdminLTE login interface
	req, _ := http.NewRequest(http.MethodGet, "/admin/login", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	body := w.Body.String()
	if !strings.Contains(body, "username") || !strings.Contains(body, "password") {
		t.Fatalf("login page missing username/password fields: %s", body)
	}

	// Unauthenticated GET /admin redirects to /admin/login
	reqAdmin, _ := http.NewRequest(http.MethodGet, "/admin", nil)
	wAdmin := httptest.NewRecorder()
	r.ServeHTTP(wAdmin, reqAdmin)

	if wAdmin.Code != http.StatusFound {
		t.Fatalf("expected status 302 redirect, got %d", wAdmin.Code)
	}
	loc := wAdmin.Header().Get("Location")
	if !strings.Contains(loc, "/admin/login") {
		t.Fatalf("expected redirect to /admin/login, got: %s", loc)
	}
}

// VAL-ADM-003: Administrative Authentication Flow
func TestAdminAuthenticationFlow(t *testing.T) {
	r := setupTestAdmin(t)

	// Valid login
	form := url.Values{}
	form.Add("username", "admin")
	form.Add("password", "admin")

	req, _ := http.NewRequest(http.MethodPost, "/admin/signin", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Code int                    `json:"code"`
		Msg  string                 `json:"msg"`
		Data map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse json response: %v", err)
	}
	if resp.Code != 200 {
		t.Fatalf("expected code 200, got %d (msg: %s)", resp.Code, resp.Msg)
	}
	if resp.Data["url"] != "/admin" {
		t.Fatalf("expected redirect url /admin, got %v", resp.Data["url"])
	}

	var sessionCookie *http.Cookie
	for _, c := range w.Result().Cookies() {
		if c.Name == "go_admin_session" {
			sessionCookie = c
			break
		}
	}
	if sessionCookie == nil {
		t.Fatalf("expected go_admin_session cookie to be set, got cookies: %v", w.Result().Cookies())
	}

	// Invalid login
	badForm := url.Values{}
	badForm.Add("username", "admin")
	badForm.Add("password", "wrongpassword")

	badReq, _ := http.NewRequest(http.MethodPost, "/admin/signin", strings.NewReader(badForm.Encode()))
	badReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	badW := httptest.NewRecorder()
	r.ServeHTTP(badW, badReq)

	if badW.Code == http.StatusOK {
		var badResp struct {
			Code int `json:"code"`
		}
		_ = json.Unmarshal(badW.Body.Bytes(), &badResp)
		if badResp.Code == 200 {
			t.Fatalf("expected failed authentication, got 200: %s", badW.Body.String())
		}
	}
}

// Authenticated dashboard access
func TestAdminDashboardAccess(t *testing.T) {
	r := setupTestAdmin(t)

	// Perform login first to obtain valid cookie
	form := url.Values{}
	form.Add("username", "admin")
	form.Add("password", "admin")

	loginReq, _ := http.NewRequest(http.MethodPost, "/admin/signin", strings.NewReader(form.Encode()))
	loginReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	loginW := httptest.NewRecorder()
	r.ServeHTTP(loginW, loginReq)

	var sessionCookie *http.Cookie
	for _, c := range loginW.Result().Cookies() {
		if c.Name == "go_admin_session" {
			sessionCookie = c
			break
		}
	}
	if sessionCookie == nil {
		t.Fatal("session cookie missing")
	}

	req, _ := http.NewRequest(http.MethodGet, "/admin", nil)
	req.AddCookie(sessionCookie)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}
	body := w.Body.String()
	if !strings.Contains(body, "EduPrimary") {
		t.Fatalf("dashboard missing EduPrimary content: %s", body)
	}
}

// Repeated InitAdmin calls (simulating multiple SetupRouter calls across test suites)
func TestAdminRepeatedInit(t *testing.T) {
	r2 := gin.New()
	eng2, err2 := InitAdmin(r2)
	if err2 != nil {
		t.Fatalf("Second InitAdmin failed: %v", err2)
	}
	if eng2 == nil {
		t.Fatal("Second engine is nil")
	}

	req, _ := http.NewRequest(http.MethodGet, "/admin/login", nil)
	w := httptest.NewRecorder()
	r2.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 on second router, got %d", w.Code)
	}
}

// VAL-CRUD-001: Staff & Teaching Quota Table Display and Filtering
func TestStaffTableDisplayAndFiltering(t *testing.T) {
	r := setupTestAdmin(t)

	// 1. Authenticate to get session cookie
	formVals := url.Values{}
	formVals.Add("username", "admin")
	formVals.Add("password", "admin")

	loginReq, _ := http.NewRequest(http.MethodPost, "/admin/signin", strings.NewReader(formVals.Encode()))
	loginReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	loginW := httptest.NewRecorder()
	r.ServeHTTP(loginW, loginReq)

	var sessionCookie *http.Cookie
	for _, c := range loginW.Result().Cookies() {
		if c.Name == "go_admin_session" {
			sessionCookie = c
			break
		}
	}
	if sessionCookie == nil {
		t.Fatal("session cookie missing")
	}

	// 2. GET /admin/info/staffs renders table
	req, _ := http.NewRequest(http.MethodGet, "/admin/info/staffs", nil)
	req.AddCookie(sessionCookie)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 for /admin/info/staffs, got %d. Body: %s", w.Code, w.Body.String())
	}

	body := w.Body.String()

	// Check table columns
	expectedColumns := []string{"Teacher Code", "Full Name", "Homeroom status", "Quota"}
	for _, col := range expectedColumns {
		if !strings.Contains(body, col) {
			t.Errorf("expected column '%s' to be present in table view HTML", col)
		}
	}

	// Check seeded teacher data
	expectedData := []string{"Nguyễn Văn An", "GV001", "Trần Thị Bình", "GV002"}
	for _, d := range expectedData {
		if !strings.Contains(body, d) {
			t.Errorf("expected seed data '%s' to be present in table view HTML", d)
		}
	}

	// 3. GET /admin/info/staffs/new renders new teacher form
	formReq, _ := http.NewRequest(http.MethodGet, "/admin/info/staffs/new", nil)
	formReq.AddCookie(sessionCookie)
	formW := httptest.NewRecorder()
	r.ServeHTTP(formW, formReq)

	if formW.Code != http.StatusOK {
		t.Fatalf("expected status 200 for /admin/info/staffs/new, got %d", formW.Code)
	}

	formBody := formW.Body.String()
	if !strings.Contains(formBody, "employee_code") || !strings.Contains(formBody, "full_name") {
		t.Errorf("expected employee_code and full_name fields in new teacher form")
	}

	// 4. Filter query: ?employee_code=GV001
	filterReq, _ := http.NewRequest(http.MethodGet, "/admin/info/staffs?employee_code=GV001", nil)
	filterReq.AddCookie(sessionCookie)
	filterW := httptest.NewRecorder()
	r.ServeHTTP(filterW, filterReq)

	if filterW.Code != http.StatusOK {
		t.Fatalf("expected status 200 for filtered /admin/info/staffs, got %d", filterW.Code)
	}

	filterBody := filterW.Body.String()
	if !strings.Contains(filterBody, "GV001") {
		t.Errorf("expected filtered body to contain GV001")
	}
}

// VAL-CRUD-002: Student & Digital Consent Table Display and Child Data Privacy Compliance
func TestStudentTableDisplayAndConsentStatus(t *testing.T) {
	r := setupTestAdmin(t)

	// 1. Authenticate to get session cookie
	formVals := url.Values{}
	formVals.Add("username", "admin")
	formVals.Add("password", "admin")

	loginReq, _ := http.NewRequest(http.MethodPost, "/admin/signin", strings.NewReader(formVals.Encode()))
	loginReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	loginW := httptest.NewRecorder()
	r.ServeHTTP(loginW, loginReq)

	var sessionCookie *http.Cookie
	for _, c := range loginW.Result().Cookies() {
		if c.Name == "go_admin_session" {
			sessionCookie = c
			break
		}
	}
	if sessionCookie == nil {
		t.Fatal("session cookie missing")
	}

	// 2. GET /admin/info/students renders student table with consent fields
	req, _ := http.NewRequest(http.MethodGet, "/admin/info/students", nil)
	req.AddCookie(sessionCookie)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 for /admin/info/students, got %d. Body: %s", w.Code, w.Body.String())
	}

	body := w.Body.String()

	// Check table columns
	expectedColumns := []string{"Student Code", "Full Name", "Digital Consent", "Consent Verified"}
	for _, col := range expectedColumns {
		if !strings.Contains(body, col) {
			t.Errorf("expected column '%s' to be present in student table view HTML", col)
		}
	}

	// Check seeded student records
	expectedData := []string{"Nguyễn Hoàng Nam", "HS2026001", "Trần Mai Linh", "HS2026002", "Lê Quốc Bảo", "HS2026003"}
	for _, d := range expectedData {
		if !strings.Contains(body, d) {
			t.Errorf("expected seed student '%s' to be present in table view HTML", d)
		}
	}

	// Check visual indicator of Decree 13 compliance
	if !strings.Contains(body, "Tuân thủ NĐ 13") {
		t.Errorf("expected visual compliance indicator 'Tuân thủ NĐ 13' in table view HTML")
	}
	if !strings.Contains(body, "Chưa tuân thủ NĐ 13") {
		t.Errorf("expected visual non-compliance indicator 'Chưa tuân thủ NĐ 13' for unverified student in table view HTML")
	}

	// 3. GET /admin/info/students/new renders new student form
	formReq, _ := http.NewRequest(http.MethodGet, "/admin/info/students/new", nil)
	formReq.AddCookie(sessionCookie)
	formW := httptest.NewRecorder()
	r.ServeHTTP(formW, formReq)

	if formW.Code != http.StatusOK {
		t.Fatalf("expected status 200 for /admin/info/students/new, got %d", formW.Code)
	}

	formBody := formW.Body.String()
	for _, expectedFormField := range []string{"student_code", "full_name", "guardian_name", "has_digital_consent", "digital_consent_verified"} {
		if !strings.Contains(formBody, expectedFormField) {
			t.Errorf("expected field '%s' in new student form HTML", expectedFormField)
		}
	}

	// 4. Filter query: ?student_code=HS2026001
	filterReq, _ := http.NewRequest(http.MethodGet, "/admin/info/students?student_code=HS2026001", nil)
	filterReq.AddCookie(sessionCookie)
	filterW := httptest.NewRecorder()
	r.ServeHTTP(filterW, filterReq)

	if filterW.Code != http.StatusOK {
		t.Fatalf("expected status 200 for filtered /admin/info/students, got %d", filterW.Code)
	}

	filterBody := filterW.Body.String()
	if !strings.Contains(filterBody, "HS2026001") || !strings.Contains(filterBody, "Nguyễn Hoàng Nam") {
		t.Errorf("expected filtered student body to contain HS2026001 and Nguyễn Hoàng Nam")
	}
}

// VAL-CRUD-003: KPI Evaluations Table Display, 4-Dimension Scores, and Decree 90/2020 20% Cap
func TestKPIEvaluationTableDisplayAndRankingCap(t *testing.T) {
	r := setupTestAdmin(t)

	// 1. Authenticate to get session cookie
	formVals := url.Values{}
	formVals.Add("username", "admin")
	formVals.Add("password", "admin")

	loginReq, _ := http.NewRequest(http.MethodPost, "/admin/signin", strings.NewReader(formVals.Encode()))
	loginReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	loginW := httptest.NewRecorder()
	r.ServeHTTP(loginW, loginReq)

	var sessionCookie *http.Cookie
	for _, c := range loginW.Result().Cookies() {
		if c.Name == "go_admin_session" {
			sessionCookie = c
			break
		}
	}
	if sessionCookie == nil {
		t.Fatal("session cookie missing")
	}

	// 2. GET /admin/info/kpi_evaluations renders KPI evaluations table
	req, _ := http.NewRequest(http.MethodGet, "/admin/info/kpi_evaluations", nil)
	req.AddCookie(sessionCookie)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 for /admin/info/kpi_evaluations, got %d. Body: %s", w.Code, w.Body.String())
	}

	body := w.Body.String()

	// Check table columns for weighted dimensions (40/30/15/15), fatal flaws, rank, and 20% cap
	expectedColumns := []string{"Staff Name", "Final Score", "Final Rank", "Honors Status", "Fatal Flaws", "20% Cap Status"}
	for _, col := range expectedColumns {
		if !strings.Contains(body, col) {
			t.Errorf("expected column '%s' to be present in KPI evaluations table view HTML", col)
		}
	}

	// Check seeded KPI records
	expectedData := []string{"Nguyễn Văn An", "Trần Thị Bình", "Lê Hoàng Cường"}
	for _, d := range expectedData {
		if !strings.Contains(body, d) {
			t.Errorf("expected seed staff '%s' to be present in KPI evaluations table view HTML", d)
		}
	}

	// Check rank badges exist
	if !strings.Contains(body, "rank-badge") {
		t.Errorf("expected 'rank-badge' class to be present in HTML for official ranks")
	}
	if !strings.Contains(body, "Xuất sắc (NĐ 90/2020)") {
		t.Errorf("expected excellent rank badge 'Xuất sắc (NĐ 90/2020)' in table view HTML")
	}

	// Check Decree 90/2020 20% cap indication
	if !strings.Contains(body, "Hạ do trần 20% (NĐ 90/2020)") {
		t.Errorf("expected 20-percent cap warning indicator 'Hạ do trần 20%% (NĐ 90/2020)' in table view HTML")
	}

	// 3. GET /admin/info/kpi_evaluations/new renders new evaluation form
	formReq, _ := http.NewRequest(http.MethodGet, "/admin/info/kpi_evaluations/new", nil)
	formReq.AddCookie(sessionCookie)
	formW := httptest.NewRecorder()
	r.ServeHTTP(formW, formReq)

	if formW.Code != http.StatusOK {
		t.Fatalf("expected status 200 for /admin/info/kpi_evaluations/new, got %d", formW.Code)
	}

	formBody := formW.Body.String()
	for _, expectedFormField := range []string{"staff_name", "final_score", "official_rank", "is_capped_by_20pct", "has_innovative_solution", "has_disciplinary_action"} {
		if !strings.Contains(formBody, expectedFormField) {
			t.Errorf("expected field '%s' in new KPI evaluation form HTML", expectedFormField)
		}
	}

	// 4. Filter query: ?staff_name=Nguyễn+Văn+An
	filterReq, _ := http.NewRequest(http.MethodGet, "/admin/info/kpi_evaluations?staff_name=Nguy%E1%BB%85n+V%C4%83n+An", nil)
	filterReq.AddCookie(sessionCookie)
	filterW := httptest.NewRecorder()
	r.ServeHTTP(filterW, filterReq)

	if filterW.Code != http.StatusOK {
		t.Fatalf("expected status 200 for filtered /admin/info/kpi_evaluations, got %d", filterW.Code)
	}

	filterBody := filterW.Body.String()
	if !strings.Contains(filterBody, "Nguyễn Văn An") {
		t.Errorf("expected filtered KPI body to contain Nguyễn Văn An")
	}
}
