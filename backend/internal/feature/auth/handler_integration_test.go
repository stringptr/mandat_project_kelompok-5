//go:build integration

package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go/jetstream"

	"github.com/stringptr/SiGizi/backend/internal/config"
	"github.com/stringptr/SiGizi/backend/internal/feature/auditlog"
	"github.com/stringptr/SiGizi/backend/internal/feature/bannedip"
	"github.com/stringptr/SiGizi/backend/internal/feature/jwtblacklist"
	"github.com/stringptr/SiGizi/backend/internal/feature/userAccount"
	"github.com/stringptr/SiGizi/backend/internal/feature/userSession"
	natsutil "github.com/stringptr/SiGizi/backend/internal/infrastructure/nats"
	"github.com/stringptr/SiGizi/backend/internal/jwtutils"
	"github.com/stringptr/SiGizi/backend/internal/middleware"
	"github.com/stringptr/SiGizi/backend/internal/testutils"
)

type integrationTestFixture struct {
	handler       http.Handler
	jwtUtil       *jwtutils.JWT
	pool          *pgxpool.Pool
	blacklistRepo *jwtblacklist.Repo
	banRepo       *bannedip.Repo
	bannedKV      jetstream.KeyValue
	blacklistKV   jetstream.KeyValue
}

func setupAuthIntegrationTest(t *testing.T) *integrationTestFixture {
	t.Helper()

	pool := testutils.NewTestDB(t)
	natsConn := testutils.NewTestNATS(t)

	ctx := context.Background()

	resetKV := func(bucket string, ttl time.Duration) jetstream.KeyValue {
		_ = natsConn.JetStream().DeleteKeyValue(ctx, bucket)
		kv, err := natsConn.CreateKeyValue(ctx, bucket, ttl)
		if err != nil {
			t.Fatalf("failed to create KV bucket %s: %v", bucket, err)
		}
		return kv
	}

	bannedKV := resetKV("banned_ips", 15*time.Minute)
	blacklistKV := resetKV("jwt_blacklist", 30*time.Minute)

	jwtUtil := jwtutils.New("test-secret-key")

	authCfg := &config.AuthConfig{
		JWTSecret:       "test-secret-key",
		AccessTokenTTL:  30 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
	}
	restrictCfg := &config.RestrictAuthConfig{
		MaxAttempt: 3,
		Duration:   1 * time.Hour,
	}

	authRepo := NewRepo(pool)
	userSessionRepo := userSession.NewRepo(pool)
	userAccountRepo := userAccount.NewRepo(pool)
	banRepo := bannedip.NewRepo(natsutil.NewKV(bannedKV))
	br := jwtblacklist.NewRepo(natsutil.NewKV(blacklistKV))
	auditLogRepo := auditlog.NewRepo(pool)

	svc := NewService(authRepo, userSessionRepo, userAccountRepo, jwtUtil, authCfg, restrictCfg, banRepo, br, auditLogRepo, nil)
	h := NewHandler(svc, &jwtUtil)

	testHandler, api := humatest.New(t, huma.DefaultConfig("Test", "1.0.0"))

	api.UseMiddleware(middleware.RealIPMiddleware())
	api.UseMiddleware(middleware.AccessTokenMiddleware(api, &jwtUtil, br))

	authAccess := huma.NewGroup(api, "")
	authRefresh := huma.NewGroup(api, "")
	userGroup := huma.NewGroup(authAccess, "")
	adminGroup := huma.NewGroup(authAccess, "")

	authAccess.UseMiddleware(middleware.AuthAccessMiddleware(api, &jwtUtil, br))
	authRefresh.UseMiddleware(middleware.AuthRefreshMiddleware(api, &jwtUtil))
	userGroup.UseMiddleware(middleware.RequireRole(authAccess, "USER"))
	adminGroup.UseMiddleware(middleware.RequireRole(authAccess, "ADMIN"))
	nonAuth := huma.NewGroup(api, "")
	nonAuth.UseMiddleware(middleware.NonAuthenticatedOnlyMiddleware(api, &jwtUtil, br))

	huma.Post(nonAuth, "/auth/register", h.Register)
	huma.Post(nonAuth, "/auth/login", h.Login)
	huma.Post(authRefresh, "/auth/refresh", h.Refresh)
	huma.Post(authRefresh, "/auth/logout", h.Logout)
	huma.Get(userGroup, "/auth/me", h.Me)
	huma.Patch(adminGroup, "/users/{id_user}/verification", h.VerifyUser)

	return &integrationTestFixture{
		handler:       testHandler,
		jwtUtil:       &jwtUtil,
		pool:          pool,
		blacklistRepo: br,
		banRepo:       banRepo,
		bannedKV:      bannedKV,
		blacklistKV:   blacklistKV,
	}
}

func (f *integrationTestFixture) cleanup(t *testing.T) {
	t.Helper()
	testutils.TruncateAuthTables(t, f.pool)
	ctx := context.Background()
	_ = f.bannedKV.Delete(ctx, "127.0.0.1")
}

func (f *integrationTestFixture) seed(t *testing.T) *testutils.AuthSeedIDs {
	t.Helper()
	return testutils.SeedAuthData(t, f.pool)
}

func TestAuthRegisterSuccess(t *testing.T) {
	f := setupAuthIntegrationTest(t)
	defer f.cleanup(t)
	f.seed(t)

	body := map[string]any{
		"email": "newuser@example.com", "password": "password123",
		"nama": "Test User", "nik": "3201020304050010",
		"jenis_kelamin": "Perempuan", "tanggal_lahir": "2000-01-01T00:00:00Z",
		"id_lokasi": 1, "role": "Ibu Hamil", "no_hp": "08123456789",
	}
	resp := testutils.DoRequest(f.handler, http.MethodPost, "/auth/register", body)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated

	testutils.TestResult{
		SRSRef: "SRS-7.1(1,2,3)", FSDRef: "FSD-2.1",
		TSDRef: "TSD-3.3-1 (baris 250-266)", NoTestScript: "TC-AUTH-001",
		Functional: "Register User - Success", Endpoint: "POST /auth/register",
		ReqType: "JSON Body", Parameter: `{"email":"newuser@example.com","password":"***","nama":"Test User"}`,
		ShouldBeSuccess: "true",
		Expectation:     "Response 201, success:true, detail: 'Register berhasil. Akun sedang diverifikasi. Silahkan dicek secara berkala.'",
	}.Log(t, pass, resp, respBody)
}

func TestAuthRegisterDuplicate(t *testing.T) {
	f := setupAuthIntegrationTest(t)
	defer f.cleanup(t)
	f.seed(t)

	body := map[string]any{
		"email": "verified@example.com", "password": "password123",
		"nama": "Verified User", "nik": "3201020304050001",
		"jenis_kelamin": "Perempuan", "tanggal_lahir": "2000-01-01T00:00:00Z",
		"id_lokasi": 1, "role": "Ibu Hamil", "no_hp": "081111111111",
	}
	resp := testutils.DoRequest(f.handler, http.MethodPost, "/auth/register", body)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusConflict

	testutils.TestResult{
		SRSRef: "SRS-7.1(2)", FSDRef: "FSD-2.1",
		TSDRef: "TSD-3.3-1 (baris 250-266)", NoTestScript: "TC-AUTH-002",
		Functional: "Register User - Duplicate", Endpoint: "POST /auth/register",
		ReqType: "JSON Body", Parameter: `{"nik":"3201020304050001","email":"verified@example.com"}`,
		ShouldBeSuccess: "false",
		Expectation:     "Response 409, success:false, detail: 'Tidak dapat mendaftar dengan identitas yang diberikan.'",
	}.Log(t, pass, resp, respBody)
}

func TestAuthRegisterValidationError(t *testing.T) {
	f := setupAuthIntegrationTest(t)
	defer f.cleanup(t)

	body := map[string]any{"email": "invalid", "password": "short", "nama": "", "nik": "123"}
	resp := testutils.DoRequest(f.handler, http.MethodPost, "/auth/register", body)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusUnprocessableEntity

	testutils.TestResult{
		SRSRef: "SRS-7.1(2)", FSDRef: "FSD-2.1",
		TSDRef: "TSD-3.3-1 (baris 250-266)", NoTestScript: "TC-AUTH-003",
		Functional: "Register User - Validation Error", Endpoint: "POST /auth/register",
		ReqType:         "JSON Body",
		Parameter:       `{"email":"invalid","password":"short","nik":"123"}`,
		ShouldBeSuccess: "false",
		Expectation:     "Response 422, success:false",
	}.Log(t, pass, resp, respBody)
}

func TestAuthLoginSuccess(t *testing.T) {
	f := setupAuthIntegrationTest(t)
	defer f.cleanup(t)
	f.seed(t)

	body := map[string]any{"email": "verified@example.com", "nik": "3201020304050001", "password": "password123"}
	resp := testutils.DoRequest(f.handler, http.MethodPost, "/auth/login", body)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated

	testutils.TestResult{
		SRSRef: "SRS-7.2(1,2,3)", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3-2 (baris 269-286)", NoTestScript: "TC-AUTH-004",
		Functional: "Login User - Success", Endpoint: "POST /auth/login",
		ReqType: "JSON Body", Parameter: `{"email":"verified@example.com","password":"***"}`,
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data contains access_token & refresh_token, Set-Cookie present",
	}.Log(t, pass, resp, respBody)
}

func TestAuthLoginWrongCredentials(t *testing.T) {
	f := setupAuthIntegrationTest(t)
	defer f.cleanup(t)
	f.seed(t)

	body := map[string]any{"email": "verified@example.com", "nik": "3201020304050001", "password": "wrongpassword"}
	resp := testutils.DoRequest(f.handler, http.MethodPost, "/auth/login", body)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusUnauthorized

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3-2 (baris 269-286)", NoTestScript: "TC-AUTH-005",
		Functional: "Login User - Wrong Credentials", Endpoint: "POST /auth/login",
		ReqType: "JSON Body", Parameter: `{"email":"verified@example.com","password":"***"}`,
		ShouldBeSuccess: "false",
		Expectation:     "Response 401, success:false, detail: 'Email, NIK, atau Password tidak valid'",
	}.Log(t, pass, resp, respBody)
}

func TestAuthLoginPendingVerification(t *testing.T) {
	f := setupAuthIntegrationTest(t)
	defer f.cleanup(t)
	f.seed(t)

	body := map[string]any{"email": "unverified@example.com", "nik": "3201020304050002", "password": "password123"}
	resp := testutils.DoRequest(f.handler, http.MethodPost, "/auth/login", body)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusUnauthorized

	testutils.TestResult{
		SRSRef: "SRS-7.1(3)", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3-2 (baris 269-286)", NoTestScript: "TC-AUTH-006",
		Functional: "Login User - Pending Verification", Endpoint: "POST /auth/login",
		ReqType: "JSON Body", Parameter: `{"email":"unverified@example.com","password":"***"}`,
		ShouldBeSuccess: "false",
		Expectation:     "Response 401, success:false, detail: 'Akun sedang dalam proses verifikasi. Silahkan dicek secara berkala.'",
	}.Log(t, pass, resp, respBody)
}

func TestAuthRefreshSuccess(t *testing.T) {
	f := setupAuthIntegrationTest(t)
	defer f.cleanup(t)
	ids := f.seed(t)

	sessionUUID := testutils.SeedActiveSession(t, f.pool, ids.VerifiedUserID)
	sessionCookie := &http.Cookie{Name: "refresh_token", Value: sessionUUID.String(), Path: "/"}

	resp := testutils.DoRequest(f.handler, http.MethodPost, "/auth/refresh", nil, sessionCookie)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated

	testutils.TestResult{
		SRSRef: "SRS-7.2(3)", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3-3 (baris 289-307)", NoTestScript: "TC-AUTH-008",
		Functional: "Refresh Token - Success", Endpoint: "POST /auth/refresh",
		ReqType: "Cookie (refresh_token)", Parameter: "Cookie: refresh_token=uuid-v7",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data contains new access_token & refresh_token",
	}.Log(t, pass, resp, respBody)
}

func TestAuthRefreshInvalidSession(t *testing.T) {
	f := setupAuthIntegrationTest(t)
	defer f.cleanup(t)

	invalidUUID, _ := uuid.NewV7()
	sessionCookie := &http.Cookie{Name: "refresh_token", Value: invalidUUID.String(), Path: "/"}

	resp := testutils.DoRequest(f.handler, http.MethodPost, "/auth/refresh", nil, sessionCookie)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusNotFound

	testutils.TestResult{
		SRSRef: "SRS-7.2(3)", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3-3 (baris 289-307)", NoTestScript: "TC-AUTH-009",
		Functional: "Refresh Token - Invalid/Expired Session", Endpoint: "POST /auth/refresh",
		ReqType: "Cookie (refresh_token)", Parameter: "Cookie: refresh_token=invalid-uuid",
		ShouldBeSuccess: "false",
		Expectation:     "Response 404, success:false, detail: 'Sesi login tidak dapat ditemukan. Silahkan login ulang.'",
	}.Log(t, pass, resp, respBody)
}

func TestAuthLogoutSuccess(t *testing.T) {
	f := setupAuthIntegrationTest(t)
	defer f.cleanup(t)
	ids := f.seed(t)

	sessionUUID := testutils.SeedActiveSession(t, f.pool, ids.VerifiedUserID)
	sessionCookie := &http.Cookie{Name: "refresh_token", Value: sessionUUID.String(), Path: "/"}

	resp := testutils.DoRequest(f.handler, http.MethodPost, "/auth/logout", nil, sessionCookie)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-7.2(4,5)", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3-4 (baris 310-324)", NoTestScript: "TC-AUTH-010",
		Functional: "Logout User - Success", Endpoint: "POST /auth/logout",
		ReqType: "Cookie (refresh_token)", Parameter: "Cookie: refresh_token=uuid-v7",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, detail: 'Logout berhasil.'",
	}.Log(t, pass, resp, respBody)
}

func TestAuthLogoutNoSession(t *testing.T) {
	f := setupAuthIntegrationTest(t)
	defer f.cleanup(t)

	invalidUUID, _ := uuid.NewV7()
	sessionCookie := &http.Cookie{Name: "refresh_token", Value: invalidUUID.String(), Path: "/"}

	resp := testutils.DoRequest(f.handler, http.MethodPost, "/auth/logout", nil, sessionCookie)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusNotFound

	testutils.TestResult{
		SRSRef: "SRS-7.2(5)", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3-4 (baris 310-324)", NoTestScript: "TC-AUTH-011",
		Functional: "Logout User - Session Not Found", Endpoint: "POST /auth/logout",
		ReqType: "Cookie (refresh_token)", Parameter: "Cookie: refresh_token=uuid-v7",
		ShouldBeSuccess: "false",
		Expectation:     "Response 404, success:false, detail: 'Sesi login tidak dapat ditemukan.'",
	}.Log(t, pass, resp, respBody)
}

func TestAuthMeSuccess(t *testing.T) {
	f := setupAuthIntegrationTest(t)
	defer f.cleanup(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/auth/me", nil, testutils.AccessCookie(f.jwtUtil, 1, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-7.2(3), SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3-5 (baris 327-340)", NoTestScript: "TC-AUTH-012",
		Functional: "Get Current User - Success", Endpoint: "GET /auth/me",
		ReqType: "Cookie (access_token)", Parameter: "Cookie: access_token=jwt-valid",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data contains id_user, email, roles",
	}.Log(t, pass, resp, respBody)
}

func TestAuthMeNoToken(t *testing.T) {
	f := setupAuthIntegrationTest(t)
	defer f.cleanup(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/auth/me", nil)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusUnauthorized

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3-5 (baris 327-340)", NoTestScript: "TC-AUTH-013",
		Functional: "Get Current User - No Token", Endpoint: "GET /auth/me",
		ReqType: "No Cookie", Parameter: "{}",
		ShouldBeSuccess: "false",
		Expectation:     "Response 401, success:false, detail: 'Silahkan login terlebih dahulu.' or 'Mohon login terlebih dahulu.'",
	}.Log(t, pass, resp, respBody)
}

func TestAuthVerifyUserSuccess(t *testing.T) {
	f := setupAuthIntegrationTest(t)
	defer f.cleanup(t)
	ids := f.seed(t)

	body := map[string]any{"status": "Aktif"}
	path := "/users/" + fmt.Sprint(ids.UnverifiedUserID) + "/verification"
	resp := testutils.DoRequest(f.handler, http.MethodPatch, path, body, testutils.AccessCookie(f.jwtUtil, ids.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-7.1(3), SRS-4.3(d)", FSDRef: "FSD-2.1",
		TSDRef: "TSD-3.3 (baris verifikasi admin)", NoTestScript: "TC-AUTH-014",
		Functional: "Verify User - Success (Admin)", Endpoint: "PATCH /users/{id_user}/verification",
		ReqType: "JSON Body + Cookie (ADMIN)", Parameter: `{"status":"Aktif"}, Role: ADMIN`,
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true",
	}.Log(t, pass, resp, respBody)
}

func TestAuthVerifyUserForbidden(t *testing.T) {
	f := setupAuthIntegrationTest(t)
	defer f.cleanup(t)
	ids := f.seed(t)

	body := map[string]any{"status": "Aktif"}
	path := "/users/" + fmt.Sprint(ids.UnverifiedUserID) + "/verification"
	resp := testutils.DoRequest(f.handler, http.MethodPatch, path, body, testutils.AccessCookie(f.jwtUtil, ids.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusForbidden

	testutils.TestResult{
		SRSRef: "SRS-SC-03, SRS-4.3", FSDRef: "FSD-2.1",
		TSDRef: "TSD-3.3 (baris verifikasi admin)", NoTestScript: "TC-AUTH-015",
		Functional: "Verify User - Forbidden (Non-Admin)", Endpoint: "PATCH /users/{id_user}/verification",
		ReqType: "JSON Body + Cookie (USER)", Parameter: `{"status":"Aktif"}, Role: USER`,
		ShouldBeSuccess: "false",
		Expectation:     "Response 403, success:false, detail: 'Tidak mempunyai akses untuk halaman ini.'",
	}.Log(t, pass, resp, respBody)
}

func doRequestWithIP(handler http.Handler, method, path string, body any, remoteAddr string, cookies ...*http.Cookie) *http.Response {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = remoteAddr
	for _, c := range cookies {
		req.AddCookie(c)
	}
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w.Result()
}

func TestAuthLoginIPLocked(t *testing.T) {
	f := setupAuthIntegrationTest(t)
	defer f.cleanup(t)
	f.seed(t)

	for i := 0; i < 3; i++ {
		body := map[string]any{"email": "verified@example.com", "nik": "3201020304050001", "password": "wrongpassword"}
		resp := doRequestWithIP(f.handler, http.MethodPost, "/auth/login", body, "192.0.2.1:12345")
		_ = resp.Body.Close()
	}

	body := map[string]any{"email": "verified@example.com", "nik": "3201020304050001", "password": "wrongpassword"}
	resp := doRequestWithIP(f.handler, http.MethodPost, "/auth/login", body, "192.0.2.1:12345")
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusForbidden

	testutils.TestResult{
		SRSRef: "SRS-7.2(2)", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3-2 (baris 269-286)", NoTestScript: "TC-AUTH-007",
		Functional: "Login User - IP Locked", Endpoint: "POST /auth/login",
		ReqType: "JSON Body", Parameter: `{"email":"verified@example.com","password":"***"}, 3x failed attempts`,
		ShouldBeSuccess: "false",
		Expectation:     "Response 403, success:false, detail: 'Akses ditolak. Terlalu banyak percobaan...'",
	}.Log(t, pass, resp, respBody)
}
