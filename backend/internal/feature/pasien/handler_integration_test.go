//go:build integration

package pasien

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/stringptr/SiGizi/backend/internal/feature/auditlog"
	"github.com/stringptr/SiGizi/backend/internal/jwtutils"
	"github.com/stringptr/SiGizi/backend/internal/testutils"
)

type pasienTestFixture struct {
	handler http.Handler
	jwtUtil *jwtutils.JWT
	pool    *pgxpool.Pool
}

type pasienSeedIDs struct {
	PosyanduID       int32
	PasienIbuHamilID int32
	PasienAnakID     int32
}

func setupPasienIntegrationTest(t *testing.T) *pasienTestFixture {
	t.Helper()

	pool := testutils.NewTestDB(t)
	handler, api, jwtUtil, br := testutils.SetupRouter(t)
	groups := testutils.CreateGroups(api, jwtUtil, br)

	pasienRepo := NewRepo(pool)
	auditRepo := auditlog.NewRepo(pool)
	svc := NewService(pasienRepo, auditRepo)
	h := NewHandler(svc)

	huma.Post(groups.AdminGroup, "/pasien/ibu-hamil", h.DaftarIbuHamil)
	huma.Post(groups.AdminGroup, "/pasien/anak", h.DaftarAnak)
	huma.Get(groups.AdminGroup, "/monitoring/pasien", h.GetAll)
	huma.Get(groups.AdminGroup, "/monitoring/pasien/search", h.Search)
	huma.Get(groups.AdminGroup, "/monitoring/pasien/{id}", h.GetByID)
	huma.Patch(groups.BidanGroup, "/pasien/{id}", h.Update)
	huma.Delete(groups.BidanGroup, "/pasien/{id}", h.Delete)

	return &pasienTestFixture{
		handler: handler,
		jwtUtil: jwtUtil,
		pool:    pool,
	}
}

func (f *pasienTestFixture) cleanup(t *testing.T) {
	t.Helper()
	testutils.TruncatePasienTables(t, f.pool)
	testutils.TruncateAuthTables(t, f.pool)
}

func (f *pasienTestFixture) seed(t *testing.T) *testutils.AuthSeedIDs {
	t.Helper()
	return testutils.SeedAuthData(t, f.pool)
}

func (f *pasienTestFixture) seedPasien(t *testing.T, authIDs *testutils.AuthSeedIDs) *pasienSeedIDs {
	t.Helper()
	posyanduID := testutils.SeedPosyandu(t, f.pool, 1, authIDs.VerifiedUserID)

	pasienIbuHamilID := testutils.SeedPasienIbuHamil(t, f.pool, authIDs.RegularUserID, posyanduID)
	pasienAnakID := testutils.SeedPasienAnak(t, f.pool, authIDs.UnverifiedUserID, posyanduID, authIDs.RegularUserID)

	return &pasienSeedIDs{
		PosyanduID:       posyanduID,
		PasienIbuHamilID: pasienIbuHamilID,
		PasienAnakID:     pasienAnakID,
	}
}

// ---------------------------------------------------------------------------
// POST /pasien/ibu-hamil
// ---------------------------------------------------------------------------

func TestPasienDaftarIbuHamilSuccess(t *testing.T) {
	f := setupPasienIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedPasien(t, authIDs)

	body := map[string]any{
		"id_user":           authIDs.VerifiedUserID,
		"id_posyandu":       seed.PosyanduID,
		"hamil_ke":          1,
		"bulan_mulai_hamil": "2026-01-01",
		"hpht":              "2026-01-15",
		"status_kehamilan":  "Trimester 1",
	}
	resp := testutils.DoRequest(f.handler, http.MethodPost, "/pasien/ibu-hamil", body,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusCreated

	testutils.TestResult{
		SRSRef: "SRS-2.3", FSDRef: "FSD-2.3",
		TSDRef: "TSD-3.3 (Endpoint Index — baris pasien)", NoTestScript: "TC-PASIEN-001",
		Functional: "Daftar Ibu Hamil — Success", Endpoint: "POST /pasien/ibu-hamil",
		ReqType: "JSON Body + Cookie (ADMIN)", Parameter: "VerifiedUserID, valid posyandu",
		ShouldBeSuccess: "true",
		Expectation:     "Response 201, success:true",
	}.Log(t, pass, resp, respBody)
}

func TestPasienDaftarIbuHamilConflict(t *testing.T) {
	f := setupPasienIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedPasien(t, authIDs)

	body := map[string]any{
		"id_user":           authIDs.RegularUserID,
		"id_posyandu":       seed.PosyanduID,
		"hamil_ke":          1,
		"bulan_mulai_hamil": "2026-01-01",
		"hpht":              "2026-01-15",
		"status_kehamilan":  "Trimester 1",
	}
	resp := testutils.DoRequest(f.handler, http.MethodPost, "/pasien/ibu-hamil", body,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusConflict

	testutils.TestResult{
		SRSRef: "SRS-2.3", FSDRef: "FSD-2.3",
		TSDRef: "TSD-3.3 (Endpoint Index — baris pasien)", NoTestScript: "TC-PASIEN-002",
		Functional: "Daftar Ibu Hamil — Conflict (Already Pasien)", Endpoint: "POST /pasien/ibu-hamil",
		ReqType: "JSON Body + Cookie (ADMIN)", Parameter: "RegularUserID yang sudah terdaftar sebagai pasien",
		ShouldBeSuccess: "false",
		Expectation:     "Response 409, success:false",
	}.Log(t, pass, resp, respBody)
}

func TestPasienDaftarIbuHamilUserNotFound(t *testing.T) {
	f := setupPasienIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedPasien(t, authIDs)

	body := map[string]any{
		"id_user":           99999,
		"id_posyandu":       seed.PosyanduID,
		"hamil_ke":          1,
		"bulan_mulai_hamil": "2026-01-01",
		"hpht":              "2026-01-15",
		"status_kehamilan":  "Trimester 1",
	}
	resp := testutils.DoRequest(f.handler, http.MethodPost, "/pasien/ibu-hamil", body,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusNotFound

	testutils.TestResult{
		SRSRef: "SRS-2.3", FSDRef: "FSD-2.3",
		TSDRef: "TSD-3.3 (Endpoint Index — baris pasien)", NoTestScript: "TC-PASIEN-003",
		Functional: "Daftar Ibu Hamil — User Not Found", Endpoint: "POST /pasien/ibu-hamil",
		ReqType: "JSON Body + Cookie (ADMIN)", Parameter: "id_user=99999 (nonexistent)",
		ShouldBeSuccess: "false",
		Expectation:     "Response 404, success:false, detail: 'User tidak ditemukan.'",
	}.Log(t, pass, resp, respBody)
}

func TestPasienDaftarIbuHamilPosyanduNotFound(t *testing.T) {
	f := setupPasienIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	body := map[string]any{
		"id_user":           authIDs.VerifiedUserID,
		"id_posyandu":       99999,
		"hamil_ke":          1,
		"bulan_mulai_hamil": "2026-01-01",
		"hpht":              "2026-01-15",
		"status_kehamilan":  "Trimester 1",
	}
	resp := testutils.DoRequest(f.handler, http.MethodPost, "/pasien/ibu-hamil", body,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusNotFound

	testutils.TestResult{
		SRSRef: "SRS-2.3", FSDRef: "FSD-2.3",
		TSDRef: "TSD-3.3 (Endpoint Index — baris pasien)", NoTestScript: "TC-PASIEN-004",
		Functional: "Daftar Ibu Hamil — Posyandu Not Found", Endpoint: "POST /pasien/ibu-hamil",
		ReqType: "JSON Body + Cookie (ADMIN)", Parameter: "id_posyandu=99999 (nonexistent)",
		ShouldBeSuccess: "false",
		Expectation:     "Response 404, success:false, detail: 'Posyandu tidak ditemukan.'",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// POST /pasien/anak
// ---------------------------------------------------------------------------

func TestPasienDaftarAnakSuccess(t *testing.T) {
	f := setupPasienIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedPasien(t, authIDs)

	body := map[string]any{
		"id_user":             authIDs.VerifiedUserID,
		"id_posyandu":         seed.PosyanduID,
		"id_wali":             authIDs.RegularUserID,
		"nama_anak":           "Bayi Test",
		"berat_lahir":         3.2,
		"panjang_lahir":       50.0,
		"hubungan_dengan_wali": "Kandung",
	}
	resp := testutils.DoRequest(f.handler, http.MethodPost, "/pasien/anak", body,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusCreated

	testutils.TestResult{
		SRSRef: "SRS-2.3", FSDRef: "FSD-2.3",
		TSDRef: "TSD-3.3 (Endpoint Index — baris pasien)", NoTestScript: "TC-PASIEN-005",
		Functional: "Daftar Anak — Success", Endpoint: "POST /pasien/anak",
		ReqType: "JSON Body + Cookie (ADMIN)", Parameter: "VerifiedUserID, valid posyandu & wali",
		ShouldBeSuccess: "true",
		Expectation:     "Response 201, success:true",
	}.Log(t, pass, resp, respBody)
}

func TestPasienDaftarAnakConflict(t *testing.T) {
	f := setupPasienIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedPasien(t, authIDs)

	body := map[string]any{
		"id_user":             authIDs.RegularUserID,
		"id_posyandu":         seed.PosyanduID,
		"id_wali":             authIDs.RegularUserID,
		"nama_anak":           "Bayi Test",
		"berat_lahir":         3.2,
		"panjang_lahir":       50.0,
		"hubungan_dengan_wali": "Kandung",
	}
	resp := testutils.DoRequest(f.handler, http.MethodPost, "/pasien/anak", body,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusConflict

	testutils.TestResult{
		SRSRef: "SRS-2.3", FSDRef: "FSD-2.3",
		TSDRef: "TSD-3.3 (Endpoint Index — baris pasien)", NoTestScript: "TC-PASIEN-006",
		Functional: "Daftar Anak — Conflict (Already Pasien)", Endpoint: "POST /pasien/anak",
		ReqType: "JSON Body + Cookie (ADMIN)", Parameter: "RegularUserID yang sudah terdaftar sebagai pasien",
		ShouldBeSuccess: "false",
		Expectation:     "Response 409, success:false",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// GET /monitoring/pasien
// ---------------------------------------------------------------------------

func TestPasienGetAllSuccessWithData(t *testing.T) {
	f := setupPasienIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	f.seedPasien(t, authIDs)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/monitoring/pasien", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-2.3", FSDRef: "FSD-2.3",
		TSDRef: "TSD-3.3-20 (baris 589-603)", NoTestScript: "TC-PASIEN-007",
		Functional: "Monitoring Pasien List — With Data", Endpoint: "GET /monitoring/pasien",
		ReqType: "Cookie (access_token)", Parameter: "Role: ADMIN, 2 pasien seeded",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data.pasien has items, total_data=2",
	}.Log(t, pass, resp, respBody)
}

func TestPasienGetAllEmpty(t *testing.T) {
	f := setupPasienIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/monitoring/pasien", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-2.3", FSDRef: "FSD-2.3",
		TSDRef: "TSD-3.3-20 (baris 589-603)", NoTestScript: "TC-PASIEN-008",
		Functional: "Monitoring Pasien List — Empty", Endpoint: "GET /monitoring/pasien",
		ReqType: "Cookie (access_token)", Parameter: "Role: ADMIN, no pasien seeded",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data.pasien empty, total_data=0",
	}.Log(t, pass, resp, respBody)
}

func TestPasienGetAllPagination(t *testing.T) {
	f := setupPasienIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	f.seedPasien(t, authIDs)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/monitoring/pasien?page=1&per_page=1", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-2.3", FSDRef: "FSD-2.3",
		TSDRef: "TSD-3.3-20 (baris 589-603)", NoTestScript: "TC-PASIEN-009",
		Functional: "Monitoring Pasien List — Pagination", Endpoint: "GET /monitoring/pasien?page=1&per_page=1",
		ReqType: "Cookie (access_token) + query", Parameter: "Role: ADMIN, page=1, per_page=1",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data.pasien has ≤1 item, total_data=2",
	}.Log(t, pass, resp, respBody)
}

func TestPasienGetAllForbidden(t *testing.T) {
	f := setupPasienIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/monitoring/pasien", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusForbidden

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3-20 (baris 589-603)", NoTestScript: "TC-PASIEN-010",
		Functional: "Monitoring Pasien List — Forbidden (Wrong Role)", Endpoint: "GET /monitoring/pasien",
		ReqType: "Cookie (access_token)", Parameter: "Role: USER (requires ADMIN)",
		ShouldBeSuccess: "false",
		Expectation:     "Response 403, success:false",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// GET /monitoring/pasien/search
// ---------------------------------------------------------------------------

func TestPasienSearchSuccessByName(t *testing.T) {
	f := setupPasienIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	f.seedPasien(t, authIDs)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/monitoring/pasien/search?q=Regular", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-2.3", FSDRef: "FSD-2.3",
		TSDRef: "TSD-3.3-22 (baris 623-636)", NoTestScript: "TC-PASIEN-011",
		Functional: "Cari Pasien — By Name", Endpoint: "GET /monitoring/pasien/search?q=Regular",
		ReqType: "Cookie (access_token) + query", Parameter: "Role: ADMIN, q=Regular",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data contains matching pasien",
	}.Log(t, pass, resp, respBody)
}

func TestPasienSearchSuccessByNIK(t *testing.T) {
	f := setupPasienIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	f.seedPasien(t, authIDs)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/monitoring/pasien/search?q=3201020304050004", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-2.3", FSDRef: "FSD-2.3",
		TSDRef: "TSD-3.3-22 (baris 623-636)", NoTestScript: "TC-PASIEN-012",
		Functional: "Cari Pasien — By NIK", Endpoint: "GET /monitoring/pasien/search?q=3201020304050004",
		ReqType: "Cookie (access_token) + query", Parameter: "Role: ADMIN, q=NIK RegularUser",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data contains matching pasien",
	}.Log(t, pass, resp, respBody)
}

func TestPasienSearchEmptyResult(t *testing.T) {
	f := setupPasienIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	f.seedPasien(t, authIDs)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/monitoring/pasien/search?q=ZZZZ", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-2.3", FSDRef: "FSD-2.3",
		TSDRef: "TSD-3.3-22 (baris 623-636)", NoTestScript: "TC-PASIEN-013",
		Functional: "Cari Pasien — Empty Result", Endpoint: "GET /monitoring/pasien/search?q=ZZZZ",
		ReqType: "Cookie (access_token) + query", Parameter: "Role: ADMIN, q=ZZZZ (no match)",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data is empty array",
	}.Log(t, pass, resp, respBody)
}

func TestPasienSearchBadRequest(t *testing.T) {
	f := setupPasienIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/monitoring/pasien/search", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusBadRequest

	testutils.TestResult{
		SRSRef: "SRS-2.3", FSDRef: "FSD-2.3",
		TSDRef: "TSD-3.3-22 (baris 623-636)", NoTestScript: "TC-PASIEN-014",
		Functional: "Cari Pasien — Bad Request (No Query)", Endpoint: "GET /monitoring/pasien/search",
		ReqType: "Cookie (access_token)", Parameter: "Role: ADMIN, no q param",
		ShouldBeSuccess: "false",
		Expectation:     "Response 400, success:false, detail: 'Parameter pencarian (q) wajib diisi.'",
	}.Log(t, pass, resp, respBody)
}

func TestPasienSearchForbidden(t *testing.T) {
	f := setupPasienIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/monitoring/pasien/search?q=test", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusForbidden

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3-22 (baris 623-636)", NoTestScript: "TC-PASIEN-015",
		Functional: "Cari Pasien — Forbidden (Wrong Role)", Endpoint: "GET /monitoring/pasien/search",
		ReqType: "Cookie (access_token) + query", Parameter: "Role: USER (requires ADMIN)",
		ShouldBeSuccess: "false",
		Expectation:     "Response 403, success:false",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// GET /monitoring/pasien/{id}
// ---------------------------------------------------------------------------

func TestPasienGetByIDIbuHamil(t *testing.T) {
	f := setupPasienIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedPasien(t, authIDs)

	path := fmt.Sprintf("/monitoring/pasien/%d", seed.PasienIbuHamilID)
	resp := testutils.DoRequest(f.handler, http.MethodGet, path, nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-2.3", FSDRef: "FSD-2.3",
		TSDRef: "TSD-3.3-21 (baris 606-619)", NoTestScript: "TC-PASIEN-016",
		Functional: "Detail Pasien — Ibu Hamil", Endpoint: "GET /monitoring/pasien/{id}",
		ReqType: "Cookie (access_token) + path param", Parameter: "Role: ADMIN, id=ibu_hamil pasien",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data contains jenis_pasien='Ibu Hamil' and data_ibu_hamil",
	}.Log(t, pass, resp, respBody)
}

func TestPasienGetByIDAnak(t *testing.T) {
	f := setupPasienIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedPasien(t, authIDs)

	path := fmt.Sprintf("/monitoring/pasien/%d", seed.PasienAnakID)
	resp := testutils.DoRequest(f.handler, http.MethodGet, path, nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-2.3", FSDRef: "FSD-2.3",
		TSDRef: "TSD-3.3-21 (baris 606-619)", NoTestScript: "TC-PASIEN-017",
		Functional: "Detail Pasien — Anak", Endpoint: "GET /monitoring/pasien/{id}",
		ReqType: "Cookie (access_token) + path param", Parameter: "Role: ADMIN, id=anak pasien",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data contains jenis_pasien='Anak' and data_anak",
	}.Log(t, pass, resp, respBody)
}

func TestPasienGetByIDNotFound(t *testing.T) {
	f := setupPasienIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/monitoring/pasien/99999", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusNotFound

	testutils.TestResult{
		SRSRef: "SRS-2.3", FSDRef: "FSD-2.3",
		TSDRef: "TSD-3.3-21 (baris 606-619)", NoTestScript: "TC-PASIEN-018",
		Functional: "Detail Pasien — Not Found", Endpoint: "GET /monitoring/pasien/99999",
		ReqType: "Cookie (access_token) + path param", Parameter: "Role: ADMIN, id=99999 (nonexistent)",
		ShouldBeSuccess: "false",
		Expectation:     "Response 404, success:false, detail: 'Pasien tidak ditemukan.'",
	}.Log(t, pass, resp, respBody)
}

func TestPasienGetByIDForbidden(t *testing.T) {
	f := setupPasienIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/monitoring/pasien/1", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusForbidden

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3-21 (baris 606-619)", NoTestScript: "TC-PASIEN-019",
		Functional: "Detail Pasien — Forbidden (Wrong Role)", Endpoint: "GET /monitoring/pasien/{id}",
		ReqType: "Cookie (access_token)", Parameter: "Role: USER (requires ADMIN)",
		ShouldBeSuccess: "false",
		Expectation:     "Response 403, success:false",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// PATCH /pasien/{id}
// ---------------------------------------------------------------------------

func TestPasienUpdateSuccess(t *testing.T) {
	f := setupPasienIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedPasien(t, authIDs)

	body := map[string]any{
		"status_kehamilan": "Trimester 2",
	}
	path := fmt.Sprintf("/pasien/%d", seed.PasienIbuHamilID)
	resp := testutils.DoRequest(f.handler, http.MethodPatch, path, body,
		testutils.AccessCookie(f.jwtUtil, authIDs.VerifiedUserID, []string{"BIDAN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-2.3", FSDRef: "FSD-2.3",
		TSDRef: "TSD-3.3 (Endpoint Index — baris pasien)", NoTestScript: "TC-PASIEN-020",
		Functional: "Update Pasien — Success", Endpoint: "PATCH /pasien/{id}",
		ReqType: "JSON Body + Cookie (BIDAN)", Parameter: "Update status_kehamilan on existing ibu hamil pasien",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data.status_kehamilan updated",
	}.Log(t, pass, resp, respBody)
}

func TestPasienUpdateNotFound(t *testing.T) {
	f := setupPasienIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	body := map[string]any{"status_kehamilan": "Trimester 2"}
	resp := testutils.DoRequest(f.handler, http.MethodPatch, "/pasien/99999", body,
		testutils.AccessCookie(f.jwtUtil, authIDs.VerifiedUserID, []string{"BIDAN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusNotFound

	testutils.TestResult{
		SRSRef: "SRS-2.3", FSDRef: "FSD-2.3",
		TSDRef: "TSD-3.3 (Endpoint Index — baris pasien)", NoTestScript: "TC-PASIEN-021",
		Functional: "Update Pasien — Not Found", Endpoint: "PATCH /pasien/99999",
		ReqType: "JSON Body + Cookie (BIDAN)", Parameter: "id=99999 (nonexistent)",
		ShouldBeSuccess: "false",
		Expectation:     "Response 404, success:false, detail: 'Pasien tidak ditemukan.'",
	}.Log(t, pass, resp, respBody)
}

func TestPasienUpdateForbidden(t *testing.T) {
	f := setupPasienIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedPasien(t, authIDs)

	body := map[string]any{"status_kehamilan": "Trimester 2"}
	path := fmt.Sprintf("/pasien/%d", seed.PasienIbuHamilID)
	resp := testutils.DoRequest(f.handler, http.MethodPatch, path, body,
		testutils.AccessCookie(f.jwtUtil, authIDs.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusForbidden

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3 (Endpoint Index — baris pasien)", NoTestScript: "TC-PASIEN-022",
		Functional: "Update Pasien — Forbidden (Wrong Role)", Endpoint: "PATCH /pasien/{id}",
		ReqType: "JSON Body + Cookie (USER)", Parameter: "Role: USER (requires BIDAN)",
		ShouldBeSuccess: "false",
		Expectation:     "Response 403, success:false",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// DELETE /pasien/{id}
// ---------------------------------------------------------------------------

func TestPasienDeleteSuccess(t *testing.T) {
	f := setupPasienIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedPasien(t, authIDs)

	path := fmt.Sprintf("/pasien/%d", seed.PasienAnakID)
	resp := testutils.DoRequest(f.handler, http.MethodDelete, path, nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.VerifiedUserID, []string{"BIDAN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-2.3", FSDRef: "FSD-2.3",
		TSDRef: "TSD-3.3 (Endpoint Index — baris pasien)", NoTestScript: "TC-PASIEN-023",
		Functional: "Delete Pasien — Success", Endpoint: "DELETE /pasien/{id}",
		ReqType: "Cookie (BIDAN)", Parameter: "id=existing anak pasien",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true",
	}.Log(t, pass, resp, respBody)
}

func TestPasienDeleteNotFound(t *testing.T) {
	f := setupPasienIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodDelete, "/pasien/99999", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.VerifiedUserID, []string{"BIDAN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusNotFound

	testutils.TestResult{
		SRSRef: "SRS-2.3", FSDRef: "FSD-2.3",
		TSDRef: "TSD-3.3 (Endpoint Index — baris pasien)", NoTestScript: "TC-PASIEN-024",
		Functional: "Delete Pasien — Not Found", Endpoint: "DELETE /pasien/99999",
		ReqType: "Cookie (BIDAN)", Parameter: "id=99999 (nonexistent)",
		ShouldBeSuccess: "false",
		Expectation:     "Response 404, success:false, detail: 'Pasien tidak ditemukan.'",
	}.Log(t, pass, resp, respBody)
}

func TestPasienDeleteForbidden(t *testing.T) {
	f := setupPasienIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodDelete, "/pasien/1", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusForbidden

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3 (Endpoint Index — baris pasien)", NoTestScript: "TC-PASIEN-025",
		Functional: "Delete Pasien — Forbidden (Wrong Role)", Endpoint: "DELETE /pasien/{id}",
		ReqType: "Cookie (USER)", Parameter: "Role: USER (requires BIDAN)",
		ShouldBeSuccess: "false",
		Expectation:     "Response 403, success:false",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// Unauthorized — no token
// ---------------------------------------------------------------------------

func TestPasienUnauthorized(t *testing.T) {
	f := setupPasienIntegrationTest(t)
	defer f.cleanup(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/monitoring/pasien", nil)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusUnauthorized

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3-20 (baris 589-603)", NoTestScript: "TC-PASIEN-026",
		Functional: "Monitoring Pasien — Unauthorized (No Token)", Endpoint: "GET /monitoring/pasien",
		ReqType: "No Cookie", Parameter: "{}",
		ShouldBeSuccess: "false",
		Expectation:     "Response 401, success:false",
	}.Log(t, pass, resp, respBody)
}
