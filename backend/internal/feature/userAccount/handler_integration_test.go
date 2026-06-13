//go:build integration

package userAccount

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stringptr/SiGizi/backend/internal/feature/auditlog"
	"github.com/stringptr/SiGizi/backend/internal/feature/auth"
	"github.com/stringptr/SiGizi/backend/internal/jwtutils"
	"github.com/stringptr/SiGizi/backend/internal/middleware"
	"github.com/stringptr/SiGizi/backend/internal/testutils"
)

type userAccountTestFixture struct {
	handler   http.Handler
	jwtUtil   *jwtutils.JWT
	pool      *pgxpool.Pool
	seedIDs   *testutils.AuthSeedIDs
	posyanduID int32
}

func setupUserAccountIntegrationTest(t *testing.T) *userAccountTestFixture {
	t.Helper()

	handler, api, jwtUtil, blacklistRepo := testutils.SetupRouter(t)
	pool := testutils.NewTestDB(t)
	groups := testutils.CreateGroups(api, jwtUtil, blacklistRepo)

	superAdminGroup := huma.NewGroup(groups.AuthAccess, "")
	superAdminGroup.UseMiddleware(middleware.RequireRole(api, "SUPER_ADMIN"))

	authRepo := auth.NewRepo(pool)
	userAccountRepo := NewRepo(pool)
	auditLogRepo := auditlog.NewRepo(pool)

	userAccountService := NewService(userAccountRepo, authRepo, auditLogRepo)
	userAccountH := NewHandler(userAccountService)

	huma.Post(superAdminGroup, "/users", userAccountH.CreateUser)
	huma.Patch(superAdminGroup, "/users/{id}/role", userAccountH.UpdateUserRole)
	huma.Get(superAdminGroup, "/admin/audit-logs", userAccountH.GetAuditLogs)
	huma.Get(groups.UserGroup, "/users/{id}", userAccountH.GetUserByID)

	seedIDs := testutils.SeedAuthData(t, pool)
	posyanduID := testutils.SeedPosyandu(t, pool, 1, seedIDs.VerifiedUserID)

	return &userAccountTestFixture{
		handler:    handler,
		jwtUtil:    jwtUtil,
		pool:       pool,
		seedIDs:    seedIDs,
		posyanduID: posyanduID,
	}
}

func cleanupUserAccountTest(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	testutils.TruncateAuthTables(t, pool)
}

func superAdminCookie(f *userAccountTestFixture) *http.Cookie {
	return testutils.AccessCookie(f.jwtUtil, f.seedIDs.AdminUserID, []string{"USER", "DINKES", "SUPER_ADMIN", "ADMIN"})
}

func adminBidanCookie(f *userAccountTestFixture) *http.Cookie {
	return testutils.AccessCookie(f.jwtUtil, f.seedIDs.VerifiedUserID, []string{"USER", "BIDAN", "ADMIN"})
}

func userRolePath(idUser int32, suffix string) string {
	return fmt.Sprintf("/users/%d%s", idUser, suffix)
}

func TestUserAccountCreateUser_Success_Bidan(t *testing.T) {
	f := setupUserAccountIntegrationTest(t)
	defer cleanupUserAccountTest(t, f.pool)

	body := map[string]any{
		"email":          "newbidan@example.com",
		"password":       "password123",
		"nama":           "Bidan Baru",
		"nik":            "3201020304050101",
		"no_hp":          "081555555555",
		"jenis_kelamin":  "Perempuan",
		"tanggal_lahir":  "1990-01-01T00:00:00Z",
		"id_lokasi":      1,
		"role":           "Bidan",
		"no_str":         "67890/STR/2026",
		"wilayah_kerja":  1,
	}

	resp := testutils.DoRequest(f.handler, http.MethodPost, "/users", body, superAdminCookie(f))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusCreated

	testutils.TestResult{
		SRSRef:          "SRS-4.3(1), SRS-4.4",
		FSDRef:          "FSD-2.X",
		TSDRef:          "TSD-3.2 (baris 242)",
		NoTestScript:    "TC-UA-001",
		Functional:      "Manajemen Pengguna SuperAdmin",
		Endpoint:        "POST /users",
		ReqType:         "JSON Body + Cookie (SUPER_ADMIN)",
		Parameter:       `{"email":"newbidan@example.com","password":"password123","nama":"Bidan Baru","nik":"3201020304050101","no_hp":"081555555555","jenis_kelamin":"Perempuan","tanggal_lahir":"1990-01-01T00:00:00Z","id_lokasi":1,"role":"Bidan","no_str":"67890/STR/2026","wilayah_kerja":1}`,
		ShouldBeSuccess: "true",
		Expectation:     "Response 201, success:true",
	}.Log(t, pass, resp, respBody)
}

func TestUserAccountCreateUser_Success_Kader(t *testing.T) {
	f := setupUserAccountIntegrationTest(t)
	defer cleanupUserAccountTest(t, f.pool)

	body := map[string]any{
		"email":          "newkader@example.com",
		"password":       "password123",
		"nama":           "Kader Baru",
		"nik":            "3201020304050102",
		"no_hp":          "081666666666",
		"jenis_kelamin":  "Laki-Laki",
		"tanggal_lahir":  "1992-05-15T00:00:00Z",
		"id_lokasi":      1,
		"role":           "Kader",
		"no_sk":          "SK/001/2026",
		"id_posyandu":    f.posyanduID,
	}

	resp := testutils.DoRequest(f.handler, http.MethodPost, "/users", body, superAdminCookie(f))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusCreated

	testutils.TestResult{
		SRSRef:          "SRS-4.3(1), SRS-4.4",
		FSDRef:          "FSD-2.X",
		TSDRef:          "TSD-3.2 (baris 242)",
		NoTestScript:    "TC-UA-002",
		Functional:      "Manajemen Pengguna SuperAdmin",
		Endpoint:        "POST /users",
		ReqType:         "JSON Body + Cookie (SUPER_ADMIN)",
		Parameter:       fmt.Sprintf(`{"email":"newkader@example.com","password":"password123","nama":"Kader Baru","nik":"3201020304050102","no_hp":"081666666666","jenis_kelamin":"Laki-Laki","tanggal_lahir":"1992-05-15T00:00:00Z","id_lokasi":1,"role":"Kader","no_sk":"SK/001/2026","id_posyandu":%d}`, f.posyanduID),
		ShouldBeSuccess: "true",
		Expectation:     "Response 201, success:true",
	}.Log(t, pass, resp, respBody)
}

func TestUserAccountCreateUser_Success_Dinkes(t *testing.T) {
	f := setupUserAccountIntegrationTest(t)
	defer cleanupUserAccountTest(t, f.pool)

	body := map[string]any{
		"email":          "newdinkes@example.com",
		"password":       "password123",
		"nama":           "Dinkes Baru",
		"nik":            "3201020304050103",
		"no_hp":          "081777777777",
		"jenis_kelamin":  "Perempuan",
		"tanggal_lahir":  "1988-10-20T00:00:00Z",
		"id_lokasi":      1,
		"role":           "Dinkes",
	}

	resp := testutils.DoRequest(f.handler, http.MethodPost, "/users", body, superAdminCookie(f))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusCreated

	testutils.TestResult{
		SRSRef:          "SRS-4.3(1), SRS-4.4",
		FSDRef:          "FSD-2.X",
		TSDRef:          "TSD-3.2 (baris 242)",
		NoTestScript:    "TC-UA-003",
		Functional:      "Manajemen Pengguna SuperAdmin",
		Endpoint:        "POST /users",
		ReqType:         "JSON Body + Cookie (SUPER_ADMIN)",
		Parameter:       `{"email":"newdinkes@example.com","password":"password123","nama":"Dinkes Baru","nik":"3201020304050103","no_hp":"081777777777","jenis_kelamin":"Perempuan","tanggal_lahir":"1988-10-20T00:00:00Z","id_lokasi":1,"role":"Dinkes"}`,
		ShouldBeSuccess: "true",
		Expectation:     "Response 201, success:true",
	}.Log(t, pass, resp, respBody)
}

func TestUserAccountCreateUser_Forbidden_AsAdmin(t *testing.T) {
	f := setupUserAccountIntegrationTest(t)
	defer cleanupUserAccountTest(t, f.pool)

	body := map[string]any{
		"email":          "unused@example.com",
		"password":       "password123",
		"nama":           "Should Not Be Created",
		"nik":            "3201020304050199",
		"no_hp":          "081999999999",
		"jenis_kelamin":  "Laki-Laki",
		"tanggal_lahir":  "1990-01-01T00:00:00Z",
		"id_lokasi":      1,
		"role":           "Bidan",
	}

	resp := testutils.DoRequest(f.handler, http.MethodPost, "/users", body, adminBidanCookie(f))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusForbidden

	testutils.TestResult{
		SRSRef:          "SRS-4.3(1), SRS-4.4, SRS-SC-03",
		FSDRef:          "FSD-2.X",
		TSDRef:          "TSD-3.2 (baris 242)",
		NoTestScript:    "TC-UA-004",
		Functional:      "Manajemen Pengguna SuperAdmin",
		Endpoint:        "POST /users",
		ReqType:         "JSON Body + Cookie (ADMIN)",
		Parameter:       `{"email":"unused@example.com","password":"password123","nama":"Should Not Be Created","nik":"3201020304050199","no_hp":"081999999999","jenis_kelamin":"Laki-Laki","tanggal_lahir":"1990-01-01T00:00:00Z","id_lokasi":1,"role":"Bidan"}`,
		ShouldBeSuccess: "false",
		Expectation:     "Response 403, success:false, detail: 'Tidak mempunyai akses untuk halaman ini.'",
	}.Log(t, pass, resp, respBody)
}

func TestUserAccountCreateUser_Unauthorized(t *testing.T) {
	f := setupUserAccountIntegrationTest(t)
	defer cleanupUserAccountTest(t, f.pool)

	body := map[string]any{
		"email":          "noauth@example.com",
		"password":       "password123",
		"nama":           "No Auth User",
		"nik":            "3201020304050198",
		"no_hp":          "081888888888",
		"jenis_kelamin":  "Laki-Laki",
		"tanggal_lahir":  "1990-01-01T00:00:00Z",
		"id_lokasi":      1,
		"role":           "Bidan",
	}

	resp := testutils.DoRequest(f.handler, http.MethodPost, "/users", body)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusUnauthorized

	testutils.TestResult{
		SRSRef:          "SRS-SC-03",
		FSDRef:          "FSD-2.X",
		TSDRef:          "TSD-3.2 (baris 242)",
		NoTestScript:    "TC-UA-005",
		Functional:      "Manajemen Pengguna SuperAdmin",
		Endpoint:        "POST /users",
		ReqType:         "No Cookie",
		Parameter:       `{"email":"noauth@example.com","password":"password123","nama":"No Auth User","nik":"3201020304050198","no_hp":"081888888888","jenis_kelamin":"Laki-Laki","tanggal_lahir":"1990-01-01T00:00:00Z","id_lokasi":1,"role":"Bidan"}`,
		ShouldBeSuccess: "false",
		Expectation:     "Response 401, success:false",
	}.Log(t, pass, resp, respBody)
}

func TestUserAccountUpdateUserRole_Success_BidanToKader(t *testing.T) {
	f := setupUserAccountIntegrationTest(t)
	defer cleanupUserAccountTest(t, f.pool)

	body := map[string]any{
		"role":        "Kader",
		"no_sk":       "SK/002/2026",
		"id_posyandu": f.posyanduID,
	}
	path := userRolePath(f.seedIDs.VerifiedUserID, "/role")

	resp := testutils.DoRequest(f.handler, http.MethodPatch, path, body, superAdminCookie(f))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef:          "SRS-4.3(1), SRS-4.4",
		FSDRef:          "FSD-2.X",
		TSDRef:          "TSD-3.2 (baris 244)",
		NoTestScript:    "TC-UA-006",
		Functional:      "Manajemen Pengguna SuperAdmin",
		Endpoint:        "PATCH /users/{id}/role",
		ReqType:         "JSON Body + Cookie (SUPER_ADMIN)",
		Parameter:       fmt.Sprintf(`{"role":"Kader","no_sk":"SK/002/2026","id_posyandu":%d}`, f.posyanduID),
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true",
	}.Log(t, pass, resp, respBody)
}

func TestUserAccountUpdateUserRole_Success_BidanToDinkes(t *testing.T) {
	f := setupUserAccountIntegrationTest(t)
	defer cleanupUserAccountTest(t, f.pool)

	body := map[string]any{
		"role": "Dinkes",
	}
	path := userRolePath(f.seedIDs.VerifiedUserID, "/role")

	resp := testutils.DoRequest(f.handler, http.MethodPatch, path, body, superAdminCookie(f))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef:          "SRS-4.3(1), SRS-4.4",
		FSDRef:          "FSD-2.X",
		TSDRef:          "TSD-3.2 (baris 244)",
		NoTestScript:    "TC-UA-007",
		Functional:      "Manajemen Pengguna SuperAdmin",
		Endpoint:        "PATCH /users/{id}/role",
		ReqType:         "JSON Body + Cookie (SUPER_ADMIN)",
		Parameter:       `{"role":"Dinkes"}`,
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true",
	}.Log(t, pass, resp, respBody)
}

func TestUserAccountUpdateUserRole_NotFound(t *testing.T) {
	f := setupUserAccountIntegrationTest(t)
	defer cleanupUserAccountTest(t, f.pool)

	body := map[string]any{
		"role": "Bidan",
	}

	resp := testutils.DoRequest(f.handler, http.MethodPatch, "/users/99999/role", body, superAdminCookie(f))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusNotFound

	testutils.TestResult{
		SRSRef:          "SRS-4.3(1), SRS-4.4",
		FSDRef:          "FSD-2.X",
		TSDRef:          "TSD-3.2 (baris 244)",
		NoTestScript:    "TC-UA-008",
		Functional:      "Manajemen Pengguna SuperAdmin",
		Endpoint:        "PATCH /users/{id}/role",
		ReqType:         "JSON Body + Cookie (SUPER_ADMIN)",
		Parameter:       `{"role":"Bidan"}`,
		ShouldBeSuccess: "false",
		Expectation:     "Response 404, success:false, detail: 'Pengguna tidak ditemukan.'",
	}.Log(t, pass, resp, respBody)
}

func TestUserAccountUpdateUserRole_Forbidden_AsAdmin(t *testing.T) {
	f := setupUserAccountIntegrationTest(t)
	defer cleanupUserAccountTest(t, f.pool)

	body := map[string]any{
		"role": "Dinkes",
	}
	path := userRolePath(f.seedIDs.VerifiedUserID, "/role")

	resp := testutils.DoRequest(f.handler, http.MethodPatch, path, body, adminBidanCookie(f))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusForbidden

	testutils.TestResult{
		SRSRef:          "SRS-4.3(1), SRS-4.4, SRS-SC-03",
		FSDRef:          "FSD-2.X",
		TSDRef:          "TSD-3.2 (baris 244)",
		NoTestScript:    "TC-UA-009",
		Functional:      "Manajemen Pengguna SuperAdmin",
		Endpoint:        "PATCH /users/{id}/role",
		ReqType:         "JSON Body + Cookie (ADMIN)",
		Parameter:       `{"role":"Dinkes"}`,
		ShouldBeSuccess: "false",
		Expectation:     "Response 403, success:false, detail: 'Tidak mempunyai akses untuk halaman ini.'",
	}.Log(t, pass, resp, respBody)
}

func TestUserAccountUpdateUserRole_Unauthorized(t *testing.T) {
	f := setupUserAccountIntegrationTest(t)
	defer cleanupUserAccountTest(t, f.pool)

	body := map[string]any{
		"role": "Kader",
	}
	path := userRolePath(f.seedIDs.VerifiedUserID, "/role")

	resp := testutils.DoRequest(f.handler, http.MethodPatch, path, body)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusUnauthorized

	testutils.TestResult{
		SRSRef:          "SRS-SC-03",
		FSDRef:          "FSD-2.X",
		TSDRef:          "TSD-3.2 (baris 244)",
		NoTestScript:    "TC-UA-010",
		Functional:      "Manajemen Pengguna SuperAdmin",
		Endpoint:        "PATCH /users/{id}/role",
		ReqType:         "No Cookie",
		Parameter:       `{"role":"Kader"}`,
		ShouldBeSuccess: "false",
		Expectation:     "Response 401, success:false",
	}.Log(t, pass, resp, respBody)
}

func TestUserAccountGetAuditLogs_Success(t *testing.T) {
	f := setupUserAccountIntegrationTest(t)
	defer cleanupUserAccountTest(t, f.pool)

	createBody := map[string]any{
		"email":          "audittest@example.com",
		"password":       "password123",
		"nama":           "Audit Test User",
		"nik":            "3201020304050201",
		"no_hp":          "081111111112",
		"jenis_kelamin":  "Perempuan",
		"tanggal_lahir":  "1990-01-01T00:00:00Z",
		"id_lokasi":      1,
		"role":           "Bidan",
		"no_str":         "STR/001/2026",
		"wilayah_kerja":  1,
	}
	createResp := testutils.DoRequest(f.handler, http.MethodPost, "/users", createBody, superAdminCookie(f))
	testutils.ReadBody(createResp)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/admin/audit-logs?page=1&per_page=20", nil, superAdminCookie(f))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK &&
		strings.Contains(string(respBody), "audit_logs") &&
		strings.Contains(string(respBody), "total_data")

	testutils.TestResult{
		SRSRef:          "SRS-SC-04, SRS-4.3(1)",
		FSDRef:          "FSD-2.X",
		TSDRef:          "TSD-3.2 (baris 245)",
		NoTestScript:    "TC-UA-011",
		Functional:      "Manajemen Pengguna SuperAdmin",
		Endpoint:        "GET /admin/audit-logs",
		ReqType:         "Cookie (SUPER_ADMIN)",
		Parameter:       "?page=1&per_page=20",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data contains audit_logs and total_data",
	}.Log(t, pass, resp, respBody)
}

func TestUserAccountGetAuditLogs_Forbidden_AsAdmin(t *testing.T) {
	f := setupUserAccountIntegrationTest(t)
	defer cleanupUserAccountTest(t, f.pool)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/admin/audit-logs?page=1&per_page=20", nil, adminBidanCookie(f))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusForbidden

	testutils.TestResult{
		SRSRef:          "SRS-SC-04, SRS-4.3(1), SRS-SC-03",
		FSDRef:          "FSD-2.X",
		TSDRef:          "TSD-3.2 (baris 245)",
		NoTestScript:    "TC-UA-012",
		Functional:      "Manajemen Pengguna SuperAdmin",
		Endpoint:        "GET /admin/audit-logs",
		ReqType:         "Cookie (ADMIN)",
		Parameter:       "?page=1&per_page=20",
		ShouldBeSuccess: "false",
		Expectation:     "Response 403, success:false, detail: 'Tidak mempunyai akses untuk halaman ini.'",
	}.Log(t, pass, resp, respBody)
}

func TestUserAccountGetAuditLogs_Unauthorized(t *testing.T) {
	f := setupUserAccountIntegrationTest(t)
	defer cleanupUserAccountTest(t, f.pool)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/admin/audit-logs?page=1&per_page=20", nil)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusUnauthorized

	testutils.TestResult{
		SRSRef:          "SRS-SC-03",
		FSDRef:          "FSD-2.X",
		TSDRef:          "TSD-3.2 (baris 245)",
		NoTestScript:    "TC-UA-013",
		Functional:      "Manajemen Pengguna SuperAdmin",
		Endpoint:        "GET /admin/audit-logs",
		ReqType:         "No Cookie",
		Parameter:       "?page=1&per_page=20",
		ShouldBeSuccess: "false",
		Expectation:     "Response 401, success:false",
	}.Log(t, pass, resp, respBody)
}


