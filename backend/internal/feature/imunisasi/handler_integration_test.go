//go:build integration

package imunisasi

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

type imunisasiTestFixture struct {
	handler http.Handler
	jwtUtil *jwtutils.JWT
	pool    *pgxpool.Pool
}

type imunisasiSeedIDs struct {
	PosyanduID     int32
	PasienID       int32
	ImunisasiIDs   []int32
}

func setupImunisasiIntegrationTest(t *testing.T) *imunisasiTestFixture {
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

	huma.Get(groups.AdminGroup, "/imunisasi", h.GetAll)
	huma.Get(groups.AdminGroup, "/imunisasi/pasien/{id_pasien}", h.GetByPasienID)
	huma.Get(groups.AdminGroup, "/imunisasi/statistik", h.GetStatistik)
	huma.Get(groups.AdminGroup, "/imunisasi/{id}", h.GetByID)
	huma.Post(groups.BidanGroup, "/imunisasi", h.Create)
	huma.Put(groups.BidanGroup, "/imunisasi/{id}", h.Update)
	huma.Delete(groups.BidanGroup, "/imunisasi/{id}", h.Delete)
	huma.Patch(groups.BidanGroup, "/imunisasi/{id}/realisasi", h.Realisasi)

	return &imunisasiTestFixture{
		handler: handler,
		jwtUtil: jwtUtil,
		pool:    pool,
	}
}

func (f *imunisasiTestFixture) cleanup(t *testing.T) {
	t.Helper()
	testutils.TruncateImunisasiTables(t, f.pool)
	testutils.TruncateNotifikasiTables(t, f.pool)
	testutils.TruncatePasienTables(t, f.pool)
	testutils.TruncateAuthTables(t, f.pool)
}

func (f *imunisasiTestFixture) seed(t *testing.T) *testutils.AuthSeedIDs {
	t.Helper()
	return testutils.SeedAuthData(t, f.pool)
}

func (f *imunisasiTestFixture) seedImunisasi(t *testing.T, authIDs *testutils.AuthSeedIDs, pasienID int32) *imunisasiSeedIDs {
	t.Helper()
	ctx := context.Background()
	now := time.Now().Truncate(time.Microsecond)

	posyanduID := testutils.SeedPosyandu(t, f.pool, 1, authIDs.VerifiedUserID)

	_, err := f.pool.Exec(ctx,
		"INSERT INTO pasien (id_pasien, id_posyandu, created_at, updated_at) VALUES ($1, $2, $3, $3) ON CONFLICT (id_pasien) DO NOTHING",
		pasienID, posyanduID, now)
	if err != nil {
		t.Fatalf("failed to seed pasien: %v", err)
	}

	var ids []int32
	for _, vaksin := range []string{"BCG", "DPT"} {
		var id int32
		err := f.pool.QueryRow(ctx, `
			INSERT INTO jadwal_imunisasi (id_pasien, nama_vaksin, tanggal_jadwal, status_imunisasi, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $5)
			RETURNING id_imunisasi`,
			pasienID, vaksin, now.AddDate(0, 0, 30), "Belum", now,
		).Scan(&id)
		if err != nil {
			t.Fatalf("failed to seed jadwal_imunisasi '%s': %v", vaksin, err)
		}
		ids = append(ids, id)
	}

	return &imunisasiSeedIDs{
		PosyanduID:   posyanduID,
		PasienID:     pasienID,
		ImunisasiIDs: ids,
	}
}

// ---------------------------------------------------------------------------
// GET /imunisasi
// ---------------------------------------------------------------------------

func TestImunisasiGetAllSuccess(t *testing.T) {
	f := setupImunisasiIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	f.seedImunisasi(t, authIDs, authIDs.RegularUserID)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/imunisasi", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-3.1", FSDRef: "FSD-3.1",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-IMUNISASI-001",
		Functional: "Imunisasi List — Success", Endpoint: "GET /imunisasi",
		ReqType: "Cookie (access_token)", Parameter: "Role: ADMIN, 2 jadwal seeded",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data.jadwal has items",
	}.Log(t, pass, resp, respBody)
}

func TestImunisasiGetAllEmpty(t *testing.T) {
	f := setupImunisasiIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/imunisasi", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-3.1", FSDRef: "FSD-3.1",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-IMUNISASI-002",
		Functional: "Imunisasi List — Empty", Endpoint: "GET /imunisasi",
		ReqType: "Cookie (access_token)", Parameter: "Role: ADMIN, no jadwal seeded",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data.jadwal empty",
	}.Log(t, pass, resp, respBody)
}

func TestImunisasiGetAllForbidden(t *testing.T) {
	f := setupImunisasiIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/imunisasi", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusForbidden

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-IMUNISASI-003",
		Functional: "Imunisasi List — Forbidden (USER)", Endpoint: "GET /imunisasi",
		ReqType: "Cookie (access_token)", Parameter: "Role: USER (requires ADMIN)",
		ShouldBeSuccess: "false",
		Expectation:     "Response 403, success:false",
	}.Log(t, pass, resp, respBody)
}

func TestImunisasiGetAllUnauthorized(t *testing.T) {
	f := setupImunisasiIntegrationTest(t)
	defer f.cleanup(t)
	f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/imunisasi", nil)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusUnauthorized

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-IMUNISASI-004",
		Functional: "Imunisasi List — Unauthorized", Endpoint: "GET /imunisasi",
		ReqType: "No Cookie", Parameter: "{}",
		ShouldBeSuccess: "false",
		Expectation:     "Response 401, success:false",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// GET /imunisasi/{id}
// ---------------------------------------------------------------------------

func TestImunisasiGetByIDSuccess(t *testing.T) {
	f := setupImunisasiIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedImunisasi(t, authIDs, authIDs.RegularUserID)

	path := fmt.Sprintf("/imunisasi/%d", seed.ImunisasiIDs[0])
	resp := testutils.DoRequest(f.handler, http.MethodGet, path, nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-3.2", FSDRef: "FSD-3.2",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-IMUNISASI-005",
		Functional: "Imunisasi Detail — Success", Endpoint: "GET /imunisasi/{id}",
		ReqType: "Cookie (access_token) + path param", Parameter: "Role: ADMIN, id=existing",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data contains detail",
	}.Log(t, pass, resp, respBody)
}

func TestImunisasiGetByIDNotFound(t *testing.T) {
	f := setupImunisasiIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/imunisasi/99999", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusNotFound

	testutils.TestResult{
		SRSRef: "SRS-3.2", FSDRef: "FSD-3.2",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-IMUNISASI-006",
		Functional: "Imunisasi Detail — Not Found", Endpoint: "GET /imunisasi/{id}",
		ReqType: "Cookie (access_token) + path param", Parameter: "Role: ADMIN, id=99999",
		ShouldBeSuccess: "false",
		Expectation:     "Response 404, success:false",
	}.Log(t, pass, resp, respBody)
}

func TestImunisasiGetByIDForbidden(t *testing.T) {
	f := setupImunisasiIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/imunisasi/1", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusForbidden

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-IMUNISASI-007",
		Functional: "Imunisasi Detail — Forbidden (USER)", Endpoint: "GET /imunisasi/{id}",
		ReqType: "Cookie (access_token)", Parameter: "Role: USER (requires ADMIN)",
		ShouldBeSuccess: "false",
		Expectation:     "Response 403, success:false",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// POST /imunisasi
// ---------------------------------------------------------------------------

func TestImunisasiCreateSuccess(t *testing.T) {
	f := setupImunisasiIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedImunisasi(t, authIDs, authIDs.RegularUserID)

	body := map[string]any{
		"id_pasien":     authIDs.RegularUserID,
		"nama_vaksin":   "Polio",
		"tanggal_jadwal": time.Now().AddDate(0, 0, 60).Format("2006-01-02"),
	}
	resp := testutils.DoRequest(f.handler, http.MethodPost, "/imunisasi", body,
		testutils.AccessCookie(f.jwtUtil, seed.PasienID, []string{"BIDAN", "USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusCreated

	testutils.TestResult{
		SRSRef: "SRS-3.3", FSDRef: "FSD-3.3",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-IMUNISASI-008",
		Functional: "Imunisasi Create — Success", Endpoint: "POST /imunisasi",
		ReqType: "JSON Body + Cookie (BIDAN)", Parameter: `{"id_pasien":"<ID>","nama_vaksin":"Polio","tanggal_jadwal":"...","status_imunisasi":"Belum"}`,
		ShouldBeSuccess: "true",
		Expectation:     "Response 201, success:true",
	}.Log(t, pass, resp, respBody)
}

func TestImunisasiCreatePasienNotFound(t *testing.T) {
	f := setupImunisasiIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	_ = authIDs

	body := map[string]any{
		"id_pasien":      99999,
		"nama_vaksin":    "Polio",
		"tanggal_jadwal": time.Now().AddDate(0, 0, 60).Format("2006-01-02"),
	}
	resp := testutils.DoRequest(f.handler, http.MethodPost, "/imunisasi", body,
		testutils.AccessCookie(f.jwtUtil, authIDs.VerifiedUserID, []string{"BIDAN", "USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusNotFound

	testutils.TestResult{
		SRSRef: "SRS-3.3", FSDRef: "FSD-3.3",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-IMUNISASI-009",
		Functional: "Imunisasi Create — Pasien Not Found", Endpoint: "POST /imunisasi",
		ReqType: "JSON Body + Cookie (BIDAN)", Parameter: `{"id_pasien":99999}`,
		ShouldBeSuccess: "false",
		Expectation:     "Response 404, success:false",
	}.Log(t, pass, resp, respBody)
}

func TestImunisasiCreateForbidden(t *testing.T) {
	f := setupImunisasiIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	body := map[string]any{
		"id_pasien":      authIDs.RegularUserID,
		"nama_vaksin":    "Polio",
		"tanggal_jadwal": time.Now().AddDate(0, 0, 60).Format("2006-01-02"),
	}
	resp := testutils.DoRequest(f.handler, http.MethodPost, "/imunisasi", body,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusForbidden

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-IMUNISASI-010",
		Functional: "Imunisasi Create — Forbidden (ADMIN)", Endpoint: "POST /imunisasi",
		ReqType: "JSON Body + Cookie (ADMIN)", Parameter: "Role: ADMIN (requires BIDAN)",
		ShouldBeSuccess: "false",
		Expectation:     "Response 403, success:false",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// PUT /imunisasi/{id}
// ---------------------------------------------------------------------------

func TestImunisasiUpdateSuccess(t *testing.T) {
	f := setupImunisasiIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedImunisasi(t, authIDs, authIDs.RegularUserID)

	body := map[string]any{"nama_vaksin": "BCG (Updated)"}
	path := fmt.Sprintf("/imunisasi/%d", seed.ImunisasiIDs[0])
	resp := testutils.DoRequest(f.handler, http.MethodPut, path, body,
		testutils.AccessCookie(f.jwtUtil, seed.PasienID, []string{"BIDAN", "USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-3.4", FSDRef: "FSD-3.4",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-IMUNISASI-011",
		Functional: "Imunisasi Update — Success", Endpoint: "PUT /imunisasi/{id}",
		ReqType: "JSON Body + Cookie (BIDAN)", Parameter: `{"nama_vaksin":"BCG (Updated)"}`,
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true",
	}.Log(t, pass, resp, respBody)
}

func TestImunisasiUpdateNotFound(t *testing.T) {
	f := setupImunisasiIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodPut, "/imunisasi/99999", map[string]any{"nama_vaksin": "Test"},
		testutils.AccessCookie(f.jwtUtil, authIDs.VerifiedUserID, []string{"BIDAN", "USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusNotFound

	testutils.TestResult{
		SRSRef: "SRS-3.4", FSDRef: "FSD-3.4",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-IMUNISASI-012",
		Functional: "Imunisasi Update — Not Found", Endpoint: "PUT /imunisasi/{id}",
		ReqType: "JSON Body + Cookie (BIDAN)", Parameter: "id=99999",
		ShouldBeSuccess: "false",
		Expectation:     "Response 404, success:false",
	}.Log(t, pass, resp, respBody)
}

func TestImunisasiUpdateForbidden(t *testing.T) {
	f := setupImunisasiIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodPut, "/imunisasi/1", map[string]any{"nama_vaksin": "Test"},
		testutils.AccessCookie(f.jwtUtil, authIDs.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusForbidden

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-IMUNISASI-013",
		Functional: "Imunisasi Update — Forbidden (USER)", Endpoint: "PUT /imunisasi/{id}",
		ReqType: "JSON Body + Cookie (USER)", Parameter: "Role: USER (requires BIDAN)",
		ShouldBeSuccess: "false",
		Expectation:     "Response 403, success:false",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// DELETE /imunisasi/{id}
// ---------------------------------------------------------------------------

func TestImunisasiDeleteSuccess(t *testing.T) {
	f := setupImunisasiIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedImunisasi(t, authIDs, authIDs.RegularUserID)

	path := fmt.Sprintf("/imunisasi/%d", seed.ImunisasiIDs[0])
	resp := testutils.DoRequest(f.handler, http.MethodDelete, path, nil,
		testutils.AccessCookie(f.jwtUtil, seed.PasienID, []string{"BIDAN", "USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-3.5", FSDRef: "FSD-3.5",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-IMUNISASI-014",
		Functional: "Imunisasi Delete — Success", Endpoint: "DELETE /imunisasi/{id}",
		ReqType: "Cookie (BIDAN)", Parameter: "id=existing",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true",
	}.Log(t, pass, resp, respBody)
}

func TestImunisasiDeleteNotFound(t *testing.T) {
	f := setupImunisasiIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodDelete, "/imunisasi/99999", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.VerifiedUserID, []string{"BIDAN", "USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusNotFound

	testutils.TestResult{
		SRSRef: "SRS-3.5", FSDRef: "FSD-3.5",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-IMUNISASI-015",
		Functional: "Imunisasi Delete — Not Found", Endpoint: "DELETE /imunisasi/{id}",
		ReqType: "Cookie (BIDAN)", Parameter: "id=99999",
		ShouldBeSuccess: "false",
		Expectation:     "Response 404, success:false",
	}.Log(t, pass, resp, respBody)
}

func TestImunisasiDeleteForbidden(t *testing.T) {
	f := setupImunisasiIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodDelete, "/imunisasi/1", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusForbidden

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-IMUNISASI-016",
		Functional: "Imunisasi Delete — Forbidden (USER)", Endpoint: "DELETE /imunisasi/{id}",
		ReqType: "Cookie (USER)", Parameter: "Role: USER (requires BIDAN)",
		ShouldBeSuccess: "false",
		Expectation:     "Response 403, success:false",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// PATCH /imunisasi/{id}/realisasi
// ---------------------------------------------------------------------------

func TestImunisasiRealisasiSuccess(t *testing.T) {
	f := setupImunisasiIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedImunisasi(t, authIDs, authIDs.RegularUserID)

	body := map[string]any{"tanggal_realisasi": time.Now().Format("2006-01-02")}
	path := fmt.Sprintf("/imunisasi/%d/realisasi", seed.ImunisasiIDs[0])
	resp := testutils.DoRequest(f.handler, http.MethodPatch, path, body,
		testutils.AccessCookie(f.jwtUtil, seed.PasienID, []string{"BIDAN", "USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-3.6", FSDRef: "FSD-3.6",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-IMUNISASI-017",
		Functional: "Imunisasi Realisasi — Success", Endpoint: "PATCH /imunisasi/{id}/realisasi",
		ReqType: "JSON Body + Cookie (BIDAN)", Parameter: `{"tanggal_realisasi":"..."}`,
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, status_imunisasi='Sudah'",
	}.Log(t, pass, resp, respBody)
}

func TestImunisasiRealisasiNotFound(t *testing.T) {
	f := setupImunisasiIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	body := map[string]any{"tanggal_realisasi": time.Now().Format("2006-01-02")}
	resp := testutils.DoRequest(f.handler, http.MethodPatch, "/imunisasi/99999/realisasi", body,
		testutils.AccessCookie(f.jwtUtil, authIDs.VerifiedUserID, []string{"BIDAN", "USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusNotFound

	testutils.TestResult{
		SRSRef: "SRS-3.6", FSDRef: "FSD-3.6",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-IMUNISASI-018",
		Functional: "Imunisasi Realisasi — Not Found", Endpoint: "PATCH /imunisasi/{id}/realisasi",
		ReqType: "JSON Body + Cookie (BIDAN)", Parameter: "id=99999",
		ShouldBeSuccess: "false",
		Expectation:     "Response 404, success:false",
	}.Log(t, pass, resp, respBody)
}

func TestImunisasiRealisasiForbidden(t *testing.T) {
	f := setupImunisasiIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	body := map[string]any{"tanggal_realisasi": time.Now().Format("2006-01-02")}
	resp := testutils.DoRequest(f.handler, http.MethodPatch, "/imunisasi/1/realisasi", body,
		testutils.AccessCookie(f.jwtUtil, authIDs.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusForbidden

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-IMUNISASI-019",
		Functional: "Imunisasi Realisasi — Forbidden (USER)", Endpoint: "PATCH /imunisasi/{id}/realisasi",
		ReqType: "JSON Body + Cookie (USER)", Parameter: "Role: USER (requires BIDAN)",
		ShouldBeSuccess: "false",
		Expectation:     "Response 403, success:false",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// GET /imunisasi/pasien/{id_pasien}
// ---------------------------------------------------------------------------

func TestImunisasiGetByPasienIDSuccess(t *testing.T) {
	f := setupImunisasiIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	f.seedImunisasi(t, authIDs, authIDs.RegularUserID)

	path := fmt.Sprintf("/imunisasi/pasien/%d", authIDs.RegularUserID)
	resp := testutils.DoRequest(f.handler, http.MethodGet, path, nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-3.7", FSDRef: "FSD-3.7",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-IMUNISASI-020",
		Functional: "Imunisasi Riwayat Pasien — Success", Endpoint: "GET /imunisasi/pasien/{id_pasien}",
		ReqType: "Cookie (access_token) + path param", Parameter: "Role: ADMIN, id=existing pasien",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data.riwayat_imunisasi has items",
	}.Log(t, pass, resp, respBody)
}

func TestImunisasiGetByPasienIDEmpty(t *testing.T) {
	f := setupImunisasiIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	ctx := context.Background()
	now := time.Now().Truncate(time.Microsecond)
	posyanduID := testutils.SeedPosyandu(t, f.pool, 1, authIDs.VerifiedUserID)
	_, err := f.pool.Exec(ctx,
		"INSERT INTO pasien (id_pasien, id_posyandu, created_at, updated_at) VALUES ($1, $2, $3, $3)",
		authIDs.RegularUserID, posyanduID, now)
	if err != nil {
		t.Fatalf("failed to seed pasien: %v", err)
	}

	path := fmt.Sprintf("/imunisasi/pasien/%d", authIDs.RegularUserID)
	resp := testutils.DoRequest(f.handler, http.MethodGet, path, nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-3.7", FSDRef: "FSD-3.7",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-IMUNISASI-021",
		Functional: "Imunisasi Riwayat Pasien — Empty", Endpoint: "GET /imunisasi/pasien/{id_pasien}",
		ReqType: "Cookie (access_token) + path param", Parameter: "Role: ADMIN, no jadwal for pasien",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data.riwayat_imunisasi empty",
	}.Log(t, pass, resp, respBody)
}

func TestImunisasiGetByPasienIDNotFound(t *testing.T) {
	f := setupImunisasiIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/imunisasi/pasien/99999", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusNotFound

	testutils.TestResult{
		SRSRef: "SRS-3.7", FSDRef: "FSD-3.7",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-IMUNISASI-022",
		Functional: "Imunisasi Riwayat Pasien — Not Found", Endpoint: "GET /imunisasi/pasien/{id_pasien}",
		ReqType: "Cookie (access_token) + path param", Parameter: "Role: ADMIN, id=99999",
		ShouldBeSuccess: "false",
		Expectation:     "Response 404, success:false",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// GET /imunisasi/statistik
// ---------------------------------------------------------------------------

func TestImunisasiGetStatistikSuccess(t *testing.T) {
	f := setupImunisasiIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	f.seedImunisasi(t, authIDs, authIDs.RegularUserID)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/imunisasi/statistik", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-3.8", FSDRef: "FSD-3.8",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-IMUNISASI-023",
		Functional: "Imunisasi Statistik — Success", Endpoint: "GET /imunisasi/statistik",
		ReqType: "Cookie (access_token)", Parameter: "Role: ADMIN, 2 jadwal seeded",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data contains statistik",
	}.Log(t, pass, resp, respBody)
}

func TestImunisasiGetStatistikEmpty(t *testing.T) {
	f := setupImunisasiIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/imunisasi/statistik", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-3.8", FSDRef: "FSD-3.8",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-IMUNISASI-024",
		Functional: "Imunisasi Statistik — Empty", Endpoint: "GET /imunisasi/statistik",
		ReqType: "Cookie (access_token)", Parameter: "Role: ADMIN, no jadwal",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data contains statistik with zero values",
	}.Log(t, pass, resp, respBody)
}

func TestImunisasiGetStatistikForbidden(t *testing.T) {
	f := setupImunisasiIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/imunisasi/statistik", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusForbidden

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-IMUNISASI-025",
		Functional: "Imunisasi Statistik — Forbidden (USER)", Endpoint: "GET /imunisasi/statistik",
		ReqType: "Cookie (access_token)", Parameter: "Role: USER (requires ADMIN)",
		ShouldBeSuccess: "false",
		Expectation:     "Response 403, success:false",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// Unauthorized — no token (representative)
// ---------------------------------------------------------------------------

func TestImunisasiCreateUnauthorized(t *testing.T) {
	f := setupImunisasiIntegrationTest(t)
	defer f.cleanup(t)
	f.seed(t)

	body := map[string]any{
		"id_pasien":       1,
		"nama_vaksin":     "Polio",
		"tanggal_jadwal":  time.Now().AddDate(0, 0, 60).Format("2006-01-02"),
	}
	resp := testutils.DoRequest(f.handler, http.MethodPost, "/imunisasi", body)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusUnauthorized

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-IMUNISASI-026",
		Functional: "Imunisasi — Unauthorized (No Token)", Endpoint: "POST /imunisasi",
		ReqType: "No Cookie", Parameter: "{}",
		ShouldBeSuccess: "false",
		Expectation:     "Response 401, success:false",
	}.Log(t, pass, resp, respBody)
}
