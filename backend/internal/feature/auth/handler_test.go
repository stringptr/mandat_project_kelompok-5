//go:build integration

package auth

import (
	"context"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofrs/uuid/v5"
	authDomain "github.com/stringptr/SiGizi/backend/internal/domain/auth"
	"github.com/stringptr/SiGizi/backend/internal/errorutils"
	"github.com/stringptr/SiGizi/backend/internal/jwtutils"
	"github.com/stringptr/SiGizi/backend/internal/testutils"
)

type mockAuthService struct {
	RegisterFunc func(context.Context, *authDomain.RegisterRequest, string) *errorutils.Error
	LoginFunc    func(context.Context, *authDomain.LoginRequest, string) (*authDomain.AuthResponse, *errorutils.Error)
	RefreshFunc  func(context.Context, uuid.UUID, string) (*authDomain.AuthResponse, *errorutils.Error)
	LogoutFunc   func(context.Context, uuid.UUID, string) *errorutils.Error
	VerifyFunc   func(context.Context, *authDomain.VerifyUserRequest) *errorutils.Error
}

func (m *mockAuthService) Register(ctx context.Context, dto *authDomain.RegisterRequest, ip string) *errorutils.Error {
	if m.RegisterFunc != nil {
		return m.RegisterFunc(ctx, dto, ip)
	}
	return nil
}

func (m *mockAuthService) Login(ctx context.Context, req *authDomain.LoginRequest, ip string) (*authDomain.AuthResponse, *errorutils.Error) {
	if m.LoginFunc != nil {
		return m.LoginFunc(ctx, req, ip)
	}
	return nil, nil
}

func (m *mockAuthService) Refresh(ctx context.Context, refreshToken uuid.UUID, ip string) (*authDomain.AuthResponse, *errorutils.Error) {
	if m.RefreshFunc != nil {
		return m.RefreshFunc(ctx, refreshToken, ip)
	}
	return nil, nil
}

func (m *mockAuthService) Logout(ctx context.Context, refreshToken uuid.UUID, accessTokenJTI string) *errorutils.Error {
	if m.LogoutFunc != nil {
		return m.LogoutFunc(ctx, refreshToken, accessTokenJTI)
	}
	return nil
}

func (m *mockAuthService) VerifyUser(ctx context.Context, req *authDomain.VerifyUserRequest) *errorutils.Error {
	if m.VerifyFunc != nil {
		return m.VerifyFunc(ctx, req)
	}
	return nil
}

func setupAuthTest(t *testing.T) (http.Handler, *mockAuthService, *jwtutils.JWT) {
	t.Helper()
	handler, api, jwtUtil, blacklistRepo := testutils.SetupRouter(t)
	groups := testutils.CreateGroups(api, jwtUtil, blacklistRepo)
	mockSvc := &mockAuthService{}
	authHandler := NewHandler(mockSvc, jwtUtil)

	huma.Post(groups.NonAuth, "/auth/register", authHandler.Register)
	huma.Post(groups.NonAuth, "/auth/login", authHandler.Login)
	huma.Post(groups.AuthRefresh, "/auth/refresh", authHandler.Refresh)
	huma.Post(groups.AuthRefresh, "/auth/logout", authHandler.Logout)
	huma.Get(groups.UserGroup, "/auth/me", authHandler.Me)
	huma.Patch(groups.AdminGroup, "/users/{id_user}/verification", authHandler.VerifyUser)

	return handler, mockSvc, jwtUtil
}

func TestAuthRegisterSuccess(t *testing.T) {
	handler, mockSvc, _ := setupAuthTest(t)
	mockSvc.RegisterFunc = func(_ context.Context, _ *authDomain.RegisterRequest, _ string) *errorutils.Error {
		return nil
	}

	body := map[string]any{
		"email": "newuser@example.com", "password": "password123",
		"nama": "Test User", "nik": "3201020304050001",
		"jenis_kelamin": "Perempuan", "tanggal_lahir": "2000-01-01T00:00:00Z",
		"id_lokasi": 1, "role": "Ibu Hamil", "no_hp": "08123456789",
	}
	resp := testutils.DoRequest(handler, "POST", "/auth/register", body)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusCreated

	testutils.TestResult{
		SRSRef: "SRS-§7.1(1,2,3)", FSDRef: "FSD-2.1",
		TSDRef: "TSD-§3.3-1 (baris 250-266)", NoTestScript: "TC-AUTH-001",
		Functional: "Register User - Success", Endpoint: "POST /auth/register",
		ReqType: "JSON Body", Parameter: `{"email":"newuser@example.com","password":"***","nama":"Test User"}`,
		ShouldBeSuccess: "true",
		Expectation:     "Response 201, success:true, detail: 'Register berhasil. Akun sedang diverifikasi. Silahkan dicek secara berkala.'",
	}.Log(t, pass, resp, respBody)
}

func TestAuthRegisterDuplicate(t *testing.T) {
	handler, mockSvc, _ := setupAuthTest(t)
	mockSvc.RegisterFunc = func(_ context.Context, _ *authDomain.RegisterRequest, _ string) *errorutils.Error {
		return &errorutils.Error{Status: http.StatusConflict, Message: "Akun dengan Email, NIK, dan Password tersebut sudah terdaftar."}
	}
	body := map[string]any{
		"email": "dup@example.com", "password": "password123",
		"nama": "Dup User", "nik": "3201020304050002",
		"jenis_kelamin": "Laki-Laki", "tanggal_lahir": "2000-01-01T00:00:00Z",
		"id_lokasi": 1, "role": "Kader", "no_hp": "08123456780",
	}
	resp := testutils.DoRequest(handler, "POST", "/auth/register", body)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusConflict

	testutils.TestResult{
		SRSRef: "SRS-§7.1(2)", FSDRef: "FSD-2.1",
		TSDRef: "TSD-§3.3-1 (baris 250-266)", NoTestScript: "TC-AUTH-002",
		Functional: "Register User - Duplicate", Endpoint: "POST /auth/register",
		ReqType: "JSON Body", Parameter: `{"nik":"3201020304050002","email":"dup@example.com"}`,
		ShouldBeSuccess: "false",
		Expectation:     "Response 409, success:false, detail: 'Akun dengan Email, NIK, dan Password tersebut sudah terdaftar.'",
	}.Log(t, pass, resp, respBody)
}

func TestAuthRegisterValidationError(t *testing.T) {
	handler, mockSvc, _ := setupAuthTest(t)
	mockSvc.RegisterFunc = func(_ context.Context, _ *authDomain.RegisterRequest, _ string) *errorutils.Error {
		return &errorutils.Error{Status: http.StatusUnprocessableEntity, Message: "Terdapat kesalahan pada pengisian berkas. Mohon dicek ulang."}
	}
	body := map[string]any{"email": "invalid-email", "password": "short", "nama": "", "nik": "123", "jenis_kelamin": "Laki-aki"}
	resp := testutils.DoRequest(handler, "POST", "/auth/register", body)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusUnprocessableEntity

	testutils.TestResult{
		SRSRef: "SRS-§7.1(2)", FSDRef: "FSD-2.1",
		TSDRef: "TSD-§3.3-1 (baris 250-266)", NoTestScript: "TC-AUTH-003",
		Functional: "Register User - Validation Error", Endpoint: "POST /auth/register",
		ReqType: "JSON Body", Parameter: `{"email":"invalid-email","password":"short","nik":"123"}`,
		ShouldBeSuccess: "false",
		Expectation:     "Response 422, success:false, detail: 'Terdapat kesalahan pada pengisian berkas. Mohon dicek ulang.'",
	}.Log(t, pass, resp, respBody)
}

func TestAuthLoginSuccess(t *testing.T) {
	handler, mockSvc, j := setupAuthTest(t)
	mockSvc.LoginFunc = func(_ context.Context, _ *authDomain.LoginRequest, _ string) (*authDomain.AuthResponse, *errorutils.Error) {
		id, _ := uuid.NewV7()
		return &authDomain.AuthResponse{
			AccessToken:  testutils.GenAccessToken(j, 1, []string{"USER"}, "test@example.com", "3201020304050001"),
			RefreshToken: id, AccessTokenExpiresIn: 1800, RefreshTokenExpiresIn: 604800,
		}, nil
	}
	body := map[string]any{"email": "test@example.com", "nik": "3201020304050001", "password": "password123"}
	resp := testutils.DoRequest(handler, "POST", "/auth/login", body)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated

	testutils.TestResult{
		SRSRef: "SRS-§7.2(1,2,3)", FSDRef: "FSD-2.2",
		TSDRef: "TSD-§3.3-2 (baris 269-286)", NoTestScript: "TC-AUTH-004",
		Functional: "Login User - Success", Endpoint: "POST /auth/login",
		ReqType: "JSON Body", Parameter: `{"email":"test@example.com","password":"***"}`,
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data contains access_token & refresh_token, Set-Cookie present",
	}.Log(t, pass, resp, respBody)
}

func TestAuthLoginWrongCredentials(t *testing.T) {
	handler, mockSvc, _ := setupAuthTest(t)
	mockSvc.LoginFunc = func(_ context.Context, _ *authDomain.LoginRequest, _ string) (*authDomain.AuthResponse, *errorutils.Error) {
		return nil, &errorutils.Error{Status: http.StatusUnauthorized, Message: "Email, NIK, atau Password tidak valid"}
	}
	body := map[string]any{"email": "wrong@example.com", "nik": "3201020304050099", "password": "wrongpassword"}
	resp := testutils.DoRequest(handler, "POST", "/auth/login", body)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusUnauthorized

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-§3.3-2 (baris 269-286)", NoTestScript: "TC-AUTH-005",
		Functional: "Login User - Wrong Credentials", Endpoint: "POST /auth/login",
		ReqType: "JSON Body", Parameter: `{"email":"wrong@example.com","password":"***"}`,
		ShouldBeSuccess: "false",
		Expectation:     "Response 401, success:false, detail: 'Email, NIK, atau Password tidak valid'",
	}.Log(t, pass, resp, respBody)
}

func TestAuthLoginPendingVerification(t *testing.T) {
	handler, mockSvc, _ := setupAuthTest(t)
	mockSvc.LoginFunc = func(_ context.Context, _ *authDomain.LoginRequest, _ string) (*authDomain.AuthResponse, *errorutils.Error) {
		return nil, &errorutils.Error{Status: http.StatusUnauthorized, Message: "Akun sedang dalam proses verifikasi. Silahkan dicek secara berkala."}
	}
	body := map[string]any{"email": "pending@example.com", "nik": "3201020304050003", "password": "password123"}
	resp := testutils.DoRequest(handler, "POST", "/auth/login", body)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusUnauthorized

	testutils.TestResult{
		SRSRef: "SRS-§7.1(3)", FSDRef: "FSD-2.2",
		TSDRef: "TSD-§3.3-2 (baris 269-286)", NoTestScript: "TC-AUTH-006",
		Functional: "Login User - Pending Verification", Endpoint: "POST /auth/login",
		ReqType: "JSON Body", Parameter: `{"email":"pending@example.com","password":"***"}`,
		ShouldBeSuccess: "false",
		Expectation:     "Response 401, success:false, detail: 'Akun sedang dalam proses verifikasi. Silahkan dicek secara berkala.'",
	}.Log(t, pass, resp, respBody)
}

func TestAuthLoginIPLocked(t *testing.T) {
	handler, mockSvc, _ := setupAuthTest(t)
	mockSvc.LoginFunc = func(_ context.Context, _ *authDomain.LoginRequest, _ string) (*authDomain.AuthResponse, *errorutils.Error) {
		return nil, &errorutils.Error{Status: http.StatusForbidden, Message: "Akses ditolak. Terlalu banyak percobaan. Silahkan coba lagi dalam 15 menit 0 detik."}
	}
	body := map[string]any{"email": "locked@example.com", "nik": "3201020304050004", "password": "wrong"}
	resp := testutils.DoRequest(handler, "POST", "/auth/login", body)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusForbidden

	testutils.TestResult{
		SRSRef: "SRS-§7.2(2)", FSDRef: "FSD-2.2",
		TSDRef: "TSD-§3.3-2 (baris 269-286)", NoTestScript: "TC-AUTH-007",
		Functional: "Login User - IP Locked", Endpoint: "POST /auth/login",
		ReqType: "JSON Body", Parameter: `{"email":"locked@example.com","password":"***"}`,
		ShouldBeSuccess: "false",
		Expectation:     "Response 403, success:false, detail: 'Akses ditolak. Terlalu banyak percobaan...'",
	}.Log(t, pass, resp, respBody)
}

func TestAuthRefreshSuccess(t *testing.T) {
	handler, mockSvc, j := setupAuthTest(t)
	mockSvc.RefreshFunc = func(_ context.Context, _ uuid.UUID, _ string) (*authDomain.AuthResponse, *errorutils.Error) {
		id, _ := uuid.NewV7()
		return &authDomain.AuthResponse{
			AccessToken:  testutils.GenAccessToken(j, 1, []string{"USER"}, "test@example.com", "3201020304050001"),
			RefreshToken: id, AccessTokenExpiresIn: 1800, RefreshTokenExpiresIn: 604800,
		}, nil
	}
	resp := testutils.DoRequest(handler, "POST", "/auth/refresh", nil, testutils.GenRefreshCookie())
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated

	testutils.TestResult{
		SRSRef: "SRS-§7.2(3)", FSDRef: "FSD-2.2",
		TSDRef: "TSD-§3.3-3 (baris 289-307)", NoTestScript: "TC-AUTH-008",
		Functional: "Refresh Token - Success", Endpoint: "POST /auth/refresh",
		ReqType: "Cookie (refresh_token)", Parameter: "Cookie: refresh_token=uuid-v7",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data contains new access_token & refresh_token",
	}.Log(t, pass, resp, respBody)
}

func TestAuthRefreshInvalidSession(t *testing.T) {
	handler, mockSvc, _ := setupAuthTest(t)
	mockSvc.RefreshFunc = func(_ context.Context, _ uuid.UUID, _ string) (*authDomain.AuthResponse, *errorutils.Error) {
		return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Sesi login tidak dapat ditemukan. Silahkan login ulang."}
	}
	resp := testutils.DoRequest(handler, "POST", "/auth/refresh", nil, testutils.GenRefreshCookie())
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusNotFound

	testutils.TestResult{
		SRSRef: "SRS-§7.2(3)", FSDRef: "FSD-2.2",
		TSDRef: "TSD-§3.3-3 (baris 289-307)", NoTestScript: "TC-AUTH-009",
		Functional: "Refresh Token - Invalid/Expired Session", Endpoint: "POST /auth/refresh",
		ReqType: "Cookie (refresh_token)", Parameter: "Cookie: refresh_token=invalid-uuid",
		ShouldBeSuccess: "false",
		Expectation:     "Response 404, success:false, detail: 'Sesi login tidak dapat ditemukan. Silahkan login ulang.'",
	}.Log(t, pass, resp, respBody)
}

func TestAuthLogoutSuccess(t *testing.T) {
	handler, mockSvc, _ := setupAuthTest(t)
	mockSvc.LogoutFunc = func(_ context.Context, _ uuid.UUID, _ string) *errorutils.Error {
		return nil
	}
	resp := testutils.DoRequest(handler, "POST", "/auth/logout", nil, testutils.GenRefreshCookie())
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-§7.2(4,5)", FSDRef: "FSD-2.2",
		TSDRef: "TSD-§3.3-4 (baris 310-324)", NoTestScript: "TC-AUTH-010",
		Functional: "Logout User - Success", Endpoint: "POST /auth/logout",
		ReqType: "Cookie (refresh_token)", Parameter: "Cookie: refresh_token=uuid-v7",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, detail: 'Logout berhasil.'",
	}.Log(t, pass, resp, respBody)
}

func TestAuthLogoutNoSession(t *testing.T) {
	handler, mockSvc, _ := setupAuthTest(t)
	mockSvc.LogoutFunc = func(_ context.Context, _ uuid.UUID, _ string) *errorutils.Error {
		return &errorutils.Error{Status: http.StatusNotFound, Message: "Sesi login tidak dapat ditemukan."}
	}
	resp := testutils.DoRequest(handler, "POST", "/auth/logout", nil, testutils.GenRefreshCookie())
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusNotFound

	testutils.TestResult{
		SRSRef: "SRS-§7.2(5)", FSDRef: "FSD-2.2",
		TSDRef: "TSD-§3.3-4 (baris 310-324)", NoTestScript: "TC-AUTH-011",
		Functional: "Logout User - Session Not Found", Endpoint: "POST /auth/logout",
		ReqType: "Cookie (refresh_token)", Parameter: "Cookie: refresh_token=uuid-v7",
		ShouldBeSuccess: "false",
		Expectation:     "Response 404, success:false, detail: 'Sesi login tidak dapat ditemukan.'",
	}.Log(t, pass, resp, respBody)
}

func TestAuthMeSuccess(t *testing.T) {
	handler, _, j := setupAuthTest(t)
	resp := testutils.DoRequest(handler, "GET", "/auth/me", nil, testutils.AccessCookie(j, 1, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-§7.2(3), SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-§3.3-5 (baris 327-340)", NoTestScript: "TC-AUTH-012",
		Functional: "Get Current User - Success", Endpoint: "GET /auth/me",
		ReqType: "Cookie (access_token)", Parameter: "Cookie: access_token=jwt-valid",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data contains id_user, email, roles",
	}.Log(t, pass, resp, respBody)
}

func TestAuthMeNoToken(t *testing.T) {
	handler, _, _ := setupAuthTest(t)
	resp := testutils.DoRequest(handler, "GET", "/auth/me", nil)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusUnauthorized

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-§3.3-5 (baris 327-340)", NoTestScript: "TC-AUTH-013",
		Functional: "Get Current User - No Token", Endpoint: "GET /auth/me",
		ReqType: "No Cookie", Parameter: "{}",
		ShouldBeSuccess: "false",
		Expectation:     "Response 401, success:false, detail: 'Silahkan login terlebih dahulu.' or 'Mohon login terlebih dahulu.'",
	}.Log(t, pass, resp, respBody)
}

func TestAuthVerifyUserSuccess(t *testing.T) {
	handler, mockSvc, j := setupAuthTest(t)
	mockSvc.VerifyFunc = func(_ context.Context, _ *authDomain.VerifyUserRequest) *errorutils.Error {
		return nil
	}
	body := map[string]any{"status": "Aktif"}
	resp := testutils.DoRequest(handler, "PATCH", "/users/1/verification", body, testutils.AccessCookie(j, 2, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-§7.1(3), SRS-§4.3(d)", FSDRef: "FSD-2.1",
		TSDRef: "TSD-§3.3 (baris verifikasi admin)", NoTestScript: "TC-AUTH-014",
		Functional: "Verify User - Success (Admin)", Endpoint: "PATCH /users/{id_user}/verification",
		ReqType: "JSON Body + Cookie (ADMIN)", Parameter: `{"status":"Aktif"}, Role: ADMIN`,
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true",
	}.Log(t, pass, resp, respBody)
}

func TestAuthVerifyUserForbidden(t *testing.T) {
	handler, _, j := setupAuthTest(t)
	body := map[string]any{"status": "Aktif"}
	resp := testutils.DoRequest(handler, "PATCH", "/users/1/verification", body, testutils.AccessCookie(j, 3, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusForbidden

	testutils.TestResult{
		SRSRef: "SRS-SC-03, SRS-§4.3", FSDRef: "FSD-2.1",
		TSDRef: "TSD-§3.3 (baris verifikasi admin)", NoTestScript: "TC-AUTH-015",
		Functional: "Verify User - Forbidden (Non-Admin)", Endpoint: "PATCH /users/{id_user}/verification",
		ReqType: "JSON Body + Cookie (USER)", Parameter: `{"status":"Aktif"}, Role: USER`,
		ShouldBeSuccess: "false",
		Expectation:     "Response 403, success:false, detail: 'Tidak mempunyai akses untuk halaman ini.'",
	}.Log(t, pass, resp, respBody)
}
