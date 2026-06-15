//go:build integration

package pemeriksaan

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/stringptr/SiGizi/backend/internal/feature/auditlog"
	"github.com/stringptr/SiGizi/backend/internal/feature/notification"
	"github.com/stringptr/SiGizi/backend/internal/jwtutils"
	"github.com/stringptr/SiGizi/backend/internal/testutils"
)

type pemeriksaanTestFixture struct {
	handler http.Handler
	jwtUtil *jwtutils.JWT
	pool    *pgxpool.Pool
}

type pemeriksaanSeedIDs struct {
	PosyanduID          int32
	PasienID            int32
	JadwalImunisasiID   int32
	HasilPemeriksaanID  int32
}

func setupPemeriksaanIntegrationTest(t *testing.T) *pemeriksaanTestFixture {
	t.Helper()

	pool := testutils.NewTestDB(t)
	handler, api, jwtUtil, br := testutils.SetupRouter(t)
	groups := testutils.CreateGroups(api, jwtUtil, br)

	repo := NewRepo(pool)
	auditRepo := auditlog.NewRepo(pool)
	notifRepo := notification.NewRepo(pool)
	notifPub := &testutils.NoopNotifPublisher{}
	svc := NewService(repo, auditRepo, notifRepo, notifPub)
	h := NewHandler(svc)

	huma.Post(groups.AdminGroup, "/monitoring/pemeriksaan", h.Create)
	huma.Get(groups.BidanGroup, "/monitoring/pemeriksaan/pending", h.GetPending)
	huma.Get(groups.AdminGroup, "/monitoring/pemeriksaan/{id}", h.GetByID)
	huma.Put(groups.AdminGroup, "/monitoring/pemeriksaan/{id}", h.Update)
	huma.Delete(groups.BidanGroup, "/monitoring/pemeriksaan/{id}", h.Delete)
	huma.Patch(groups.BidanGroup, "/monitoring/pemeriksaan/{id}/verify", h.Verify)

	return &pemeriksaanTestFixture{
		handler: handler,
		jwtUtil: jwtUtil,
		pool:    pool,
	}
}

func (f *pemeriksaanTestFixture) cleanup(t *testing.T) {
	t.Helper()
	testutils.TruncateImunisasiTables(t, f.pool)
	testutils.TruncateNotifikasiTables(t, f.pool)
	testutils.TruncatePasienTables(t, f.pool)
	testutils.TruncateAuthTables(t, f.pool)
}

func (f *pemeriksaanTestFixture) seed(t *testing.T) *testutils.AuthSeedIDs {
	t.Helper()
	return testutils.SeedAuthData(t, f.pool)
}

func (f *pemeriksaanTestFixture) seedPemeriksaan(t *testing.T, authIDs *testutils.AuthSeedIDs) *pemeriksaanSeedIDs {
	t.Helper()
	ctx := context.Background()
	now := time.Now().Truncate(time.Microsecond)

	posyanduID := testutils.SeedPosyandu(t, f.pool, 1, authIDs.VerifiedUserID)
	pasienID := authIDs.RegularUserID

	_, err := f.pool.Exec(ctx,
		"INSERT INTO pasien (id_pasien, id_posyandu, created_at, updated_at) VALUES ($1, $2, $3, $3) ON CONFLICT (id_pasien) DO NOTHING",
		pasienID, posyanduID, now)
	if err != nil {
		t.Fatalf("failed to seed pasien: %v", err)
	}

	var jadwalID int32
	err = f.pool.QueryRow(ctx, `
		INSERT INTO jadwal_imunisasi (id_pasien, nama_vaksin, tanggal_jadwal, status_imunisasi, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $5)
		RETURNING id_imunisasi`,
		pasienID, "DPT", now.AddDate(0, 0, 30), "Belum", now,
	).Scan(&jadwalID)
	if err != nil {
		t.Fatalf("failed to seed jadwal_imunisasi: %v", err)
	}

	var hasilID int32
	err = f.pool.QueryRow(ctx, `
		INSERT INTO hasil_pemeriksaan (id_petugas_input, id_jadwal_imunisasi, berat_badan, tinggi_badan, lingkar_kepala, tekanan_darah, status_stunting, status_gizi, catatan, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $10)
		RETURNING id_hasil_pemeriksaan`,
		authIDs.AdminUserID, jadwalID, 12.5, 85.0, 45.0, "120/80", "Normal", "Gizi Baik", "Sehat", now,
	).Scan(&hasilID)
	if err != nil {
		t.Fatalf("failed to seed hasil_pemeriksaan: %v", err)
	}

	return &pemeriksaanSeedIDs{
		PosyanduID:         posyanduID,
		PasienID:           pasienID,
		JadwalImunisasiID:  jadwalID,
		HasilPemeriksaanID: hasilID,
	}
}

// ---------------------------------------------------------------------------
// POST /monitoring/pemeriksaan
// ---------------------------------------------------------------------------

func TestPemeriksaanCreateSuccess(t *testing.T) {
	f := setupPemeriksaanIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedPemeriksaan(t, authIDs)

	body := map[string]any{
		"id_jadwal_imunisasi": seed.JadwalImunisasiID,
		"berat_badan":         14.0,
		"tinggi_badan":        88.0,
		"lingkar_kepala":      46.0,
		"tekanan_darah":       "110/70",
		"catatan":             "Perkembangan baik",
	}
	resp := testutils.DoRequest(f.handler, http.MethodPost, "/monitoring/pemeriksaan", body,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusCreated

	testutils.TestResult{
		SRSRef: "SRS-2.4", FSDRef: "FSD-2.4",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-PERIKSA-001",
		Functional: "Pemeriksaan Create — Success", Endpoint: "POST /monitoring/pemeriksaan",
		ReqType: "JSON Body + Cookie (ADMIN)", Parameter: `{"id_jadwal_imunisasi":"<ID>","berat_badan":14,"tinggi_badan":88,"lingkar_kepala":46,"tekanan_darah":"110/70"}`,
		ShouldBeSuccess: "true",
		Expectation:     "Response 201, success:true, data contains id_hasil_pemeriksaan",
	}.Log(t, pass, resp, respBody)
}

func TestPemeriksaanCreateJadwalNotFound(t *testing.T) {
	f := setupPemeriksaanIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	body := map[string]any{
		"id_jadwal_imunisasi": 99999,
		"berat_badan":         12.0,
		"tinggi_badan":        80.0,
		"lingkar_kepala":      44.0,
		"tekanan_darah":       "120/80",
	}
	resp := testutils.DoRequest(f.handler, http.MethodPost, "/monitoring/pemeriksaan", body,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusNotFound

	testutils.TestResult{
		SRSRef: "SRS-2.4", FSDRef: "FSD-2.4",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-PERIKSA-002",
		Functional: "Pemeriksaan Create — Jadwal Not Found", Endpoint: "POST /monitoring/pemeriksaan",
		ReqType: "JSON Body + Cookie (ADMIN)", Parameter: `{"id_jadwal_imunisasi":99999}`,
		ShouldBeSuccess: "false",
		Expectation:     "Response 404, success:false",
	}.Log(t, pass, resp, respBody)
}

func TestPemeriksaanCreateForbidden(t *testing.T) {
	f := setupPemeriksaanIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedPemeriksaan(t, authIDs)

	body := map[string]any{
		"id_jadwal_imunisasi": seed.JadwalImunisasiID,
		"berat_badan":         12.0,
		"tinggi_badan":        80.0,
		"lingkar_kepala":      44.0,
		"tekanan_darah":       "120/80",
	}
	resp := testutils.DoRequest(f.handler, http.MethodPost, "/monitoring/pemeriksaan", body,
		testutils.AccessCookie(f.jwtUtil, authIDs.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusForbidden

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-PERIKSA-003",
		Functional: "Pemeriksaan Create — Forbidden (USER)", Endpoint: "POST /monitoring/pemeriksaan",
		ReqType: "JSON Body + Cookie (USER)", Parameter: "Role: USER (requires ADMIN)",
		ShouldBeSuccess: "false",
		Expectation:     "Response 403, success:false",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// GET /monitoring/pemeriksaan/{id}
// ---------------------------------------------------------------------------

func TestPemeriksaanGetByIDSuccess(t *testing.T) {
	f := setupPemeriksaanIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedPemeriksaan(t, authIDs)

	path := fmt.Sprintf("/monitoring/pemeriksaan/%d", seed.HasilPemeriksaanID)
	resp := testutils.DoRequest(f.handler, http.MethodGet, path, nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-2.4", FSDRef: "FSD-2.4",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-PERIKSA-004",
		Functional: "Pemeriksaan Detail — Success", Endpoint: "GET /monitoring/pemeriksaan/{id}",
		ReqType: "Cookie (access_token) + path param", Parameter: "Role: ADMIN, id=existing",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data contains detail",
	}.Log(t, pass, resp, respBody)
}

func TestPemeriksaanGetByIDNotFound(t *testing.T) {
	f := setupPemeriksaanIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/monitoring/pemeriksaan/99999", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusNotFound

	testutils.TestResult{
		SRSRef: "SRS-2.4", FSDRef: "FSD-2.4",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-PERIKSA-005",
		Functional: "Pemeriksaan Detail — Not Found", Endpoint: "GET /monitoring/pemeriksaan/{id}",
		ReqType: "Cookie (access_token) + path param", Parameter: "Role: ADMIN, id=99999",
		ShouldBeSuccess: "false",
		Expectation:     "Response 404, success:false",
	}.Log(t, pass, resp, respBody)
}

func TestPemeriksaanGetByIDForbidden(t *testing.T) {
	f := setupPemeriksaanIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/monitoring/pemeriksaan/1", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusForbidden

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-PERIKSA-006",
		Functional: "Pemeriksaan Detail — Forbidden (USER)", Endpoint: "GET /monitoring/pemeriksaan/{id}",
		ReqType: "Cookie (access_token)", Parameter: "Role: USER (requires ADMIN)",
		ShouldBeSuccess: "false",
		Expectation:     "Response 403, success:false",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// PUT /monitoring/pemeriksaan/{id}
// ---------------------------------------------------------------------------

func TestPemeriksaanUpdateSuccess(t *testing.T) {
	f := setupPemeriksaanIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedPemeriksaan(t, authIDs)

	body := map[string]any{"berat_badan": 15.0, "tinggi_badan": 90.0}
	path := fmt.Sprintf("/monitoring/pemeriksaan/%d", seed.HasilPemeriksaanID)
	resp := testutils.DoRequest(f.handler, http.MethodPut, path, body,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-2.4", FSDRef: "FSD-2.4",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-PERIKSA-007",
		Functional: "Pemeriksaan Update — Success", Endpoint: "PUT /monitoring/pemeriksaan/{id}",
		ReqType: "JSON Body + Cookie (ADMIN)", Parameter: `{"berat_badan":15,"tinggi_badan":90}`,
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true",
	}.Log(t, pass, resp, respBody)
}

func TestPemeriksaanUpdateNotFound(t *testing.T) {
	f := setupPemeriksaanIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodPut, "/monitoring/pemeriksaan/99999", map[string]any{"berat_badan": 15.0},
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusNotFound

	testutils.TestResult{
		SRSRef: "SRS-2.4", FSDRef: "FSD-2.4",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-PERIKSA-008",
		Functional: "Pemeriksaan Update — Not Found", Endpoint: "PUT /monitoring/pemeriksaan/{id}",
		ReqType: "JSON Body + Cookie (ADMIN)", Parameter: "id=99999",
		ShouldBeSuccess: "false",
		Expectation:     "Response 404, success:false",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// DELETE /monitoring/pemeriksaan/{id}
// ---------------------------------------------------------------------------

func TestPemeriksaanDeleteSuccess(t *testing.T) {
	f := setupPemeriksaanIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedPemeriksaan(t, authIDs)

	path := fmt.Sprintf("/monitoring/pemeriksaan/%d", seed.HasilPemeriksaanID)
	resp := testutils.DoRequest(f.handler, http.MethodDelete, path, nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.VerifiedUserID, []string{"BIDAN", "USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-2.4", FSDRef: "FSD-2.4",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-PERIKSA-009",
		Functional: "Pemeriksaan Delete — Success", Endpoint: "DELETE /monitoring/pemeriksaan/{id}",
		ReqType: "Cookie (BIDAN)", Parameter: "id=existing",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true",
	}.Log(t, pass, resp, respBody)
}

func TestPemeriksaanDeleteNotFound(t *testing.T) {
	f := setupPemeriksaanIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodDelete, "/monitoring/pemeriksaan/99999", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.VerifiedUserID, []string{"BIDAN", "USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusNotFound

	testutils.TestResult{
		SRSRef: "SRS-2.4", FSDRef: "FSD-2.4",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-PERIKSA-010",
		Functional: "Pemeriksaan Delete — Not Found", Endpoint: "DELETE /monitoring/pemeriksaan/{id}",
		ReqType: "Cookie (BIDAN)", Parameter: "id=99999",
		ShouldBeSuccess: "false",
		Expectation:     "Response 404, success:false",
	}.Log(t, pass, resp, respBody)
}

func TestPemeriksaanDeleteForbidden(t *testing.T) {
	f := setupPemeriksaanIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodDelete, "/monitoring/pemeriksaan/1", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusForbidden

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-PERIKSA-011",
		Functional: "Pemeriksaan Delete — Forbidden (USER)", Endpoint: "DELETE /monitoring/pemeriksaan/{id}",
		ReqType: "Cookie (USER)", Parameter: "Role: USER (requires BIDAN)",
		ShouldBeSuccess: "false",
		Expectation:     "Response 403, success:false",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// PATCH /monitoring/pemeriksaan/{id}/verify
// ---------------------------------------------------------------------------

func TestPemeriksaanVerifySuccess(t *testing.T) {
	f := setupPemeriksaanIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedPemeriksaan(t, authIDs)

	path := fmt.Sprintf("/monitoring/pemeriksaan/%d/verify", seed.HasilPemeriksaanID)
	resp := testutils.DoRequest(f.handler, http.MethodPatch, path, nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.VerifiedUserID, []string{"BIDAN", "USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-2.4", FSDRef: "FSD-2.4",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-PERIKSA-012",
		Functional: "Pemeriksaan Verify — Success", Endpoint: "PATCH /monitoring/pemeriksaan/{id}/verify",
		ReqType: "Cookie (BIDAN) + path param", Parameter: "id=existing",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true",
	}.Log(t, pass, resp, respBody)
}

func TestPemeriksaanVerifyNotFound(t *testing.T) {
	f := setupPemeriksaanIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodPatch, "/monitoring/pemeriksaan/99999/verify", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.VerifiedUserID, []string{"BIDAN", "USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusNotFound

	testutils.TestResult{
		SRSRef: "SRS-2.4", FSDRef: "FSD-2.4",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-PERIKSA-013",
		Functional: "Pemeriksaan Verify — Not Found", Endpoint: "PATCH /monitoring/pemeriksaan/{id}/verify",
		ReqType: "Cookie (BIDAN) + path param", Parameter: "id=99999",
		ShouldBeSuccess: "false",
		Expectation:     "Response 404, success:false",
	}.Log(t, pass, resp, respBody)
}

func TestPemeriksaanVerifyForbidden(t *testing.T) {
	f := setupPemeriksaanIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodPatch, "/monitoring/pemeriksaan/1/verify", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusForbidden

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-PERIKSA-014",
		Functional: "Pemeriksaan Verify — Forbidden (USER)", Endpoint: "PATCH /monitoring/pemeriksaan/{id}/verify",
		ReqType: "Cookie (USER)", Parameter: "Role: USER (requires BIDAN)",
		ShouldBeSuccess: "false",
		Expectation:     "Response 403, success:false",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// GET /monitoring/pemeriksaan/pending
// ---------------------------------------------------------------------------

func TestPemeriksaanGetPendingWithData(t *testing.T) {
	f := setupPemeriksaanIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	f.seedPemeriksaan(t, authIDs)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/monitoring/pemeriksaan/pending", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.VerifiedUserID, []string{"BIDAN", "USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-2.4", FSDRef: "FSD-2.4",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-PERIKSA-015",
		Functional: "Pemeriksaan Pending — With Data", Endpoint: "GET /monitoring/pemeriksaan/pending",
		ReqType: "Cookie (BIDAN)", Parameter: "Role: BIDAN, 1 seeded",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data.pemeriksaan_pending has items",
	}.Log(t, pass, resp, respBody)
}

func TestPemeriksaanGetPendingEmpty(t *testing.T) {
	f := setupPemeriksaanIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/monitoring/pemeriksaan/pending", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.VerifiedUserID, []string{"BIDAN", "USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-2.4", FSDRef: "FSD-2.4",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-PERIKSA-016",
		Functional: "Pemeriksaan Pending — Empty", Endpoint: "GET /monitoring/pemeriksaan/pending",
		ReqType: "Cookie (BIDAN)", Parameter: "Role: BIDAN, no data",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data.pemeriksaan_pending empty",
	}.Log(t, pass, resp, respBody)
}

func TestPemeriksaanGetPendingForbidden(t *testing.T) {
	f := setupPemeriksaanIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/monitoring/pemeriksaan/pending", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusForbidden

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-PERIKSA-017",
		Functional: "Pemeriksaan Pending — Forbidden (USER)", Endpoint: "GET /monitoring/pemeriksaan/pending",
		ReqType: "Cookie (USER)", Parameter: "Role: USER (requires BIDAN)",
		ShouldBeSuccess: "false",
		Expectation:     "Response 403, success:false",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// Unauthorized
// ---------------------------------------------------------------------------

func TestPemeriksaanUnauthorized(t *testing.T) {
	f := setupPemeriksaanIntegrationTest(t)
	defer f.cleanup(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/monitoring/pemeriksaan", nil)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusUnauthorized

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-PERIKSA-018",
		Functional: "Pemeriksaan — Unauthorized (No Token)", Endpoint: "GET /monitoring/pemeriksaan",
		ReqType: "No Cookie", Parameter: "{}",
		ShouldBeSuccess: "false",
		Expectation:     "Response 401, success:false",
	}.Log(t, pass, resp, respBody)
}
