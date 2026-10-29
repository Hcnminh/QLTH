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
