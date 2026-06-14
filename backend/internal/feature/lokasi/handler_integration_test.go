//go:build integration

package lokasi

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/stringptr/SiGizi/backend/internal/testutils"
)

type lokasiTestFixture struct {
	handler http.Handler
	pool    *pgxpool.Pool
}

func setupLokasiIntegrationTest(t *testing.T) *lokasiTestFixture {
	t.Helper()

	pool := testutils.NewTestDB(t)
	handler, api, jwtUtil, br := testutils.SetupRouter(t)
	groups := testutils.CreateGroups(api, jwtUtil, br)

	repo := NewRepo(pool)
	svc := NewService(repo)
	h := NewHandler(svc)

	huma.Get(groups.NonAuth, "/lokasi", h.GetLokasi)

	return &lokasiTestFixture{handler: handler, pool: pool}
}

func (f *lokasiTestFixture) cleanup(t *testing.T) {
	t.Helper()
	ctx := context.Background()
	_, err := f.pool.Exec(ctx, "TRUNCATE TABLE lokasi CASCADE; ALTER SEQUENCE lokasi_id_lokasi_seq RESTART WITH 1")
	if err != nil {
		t.Fatalf("failed to truncate lokasi: %v", err)
	}
}

func (f *lokasiTestFixture) seed(t *testing.T) {
	t.Helper()
	ctx := context.Background()

	_, err := f.pool.Exec(ctx, `
		INSERT INTO lokasi (id_lokasi, nama_lokasi, tipe_lokasi, bagian_dari) VALUES
		(1, 'Provinsi A', 'Provinsi', NULL),
		(2, 'Provinsi B', 'Provinsi', NULL),
		(3, 'Kabupaten A1', 'Kabupaten', 1),
		(4, 'Kabupaten B1', 'Kabupaten', 2),
		(5, 'Kecamatan A1.1', 'Kecamatan', 3),
		(6, 'Kelurahan A1.1.1', 'Kelurahan', 5)
	`)
	if err != nil {
		t.Fatalf("failed to seed lokasi: %v", err)
	}
}

// ---------------------------------------------------------------------------
// GET /lokasi?tipe=Provinsi
// ---------------------------------------------------------------------------

func TestLokasiGetProvinsiSuccess(t *testing.T) {
	f := setupLokasiIntegrationTest(t)
	defer f.cleanup(t)
	f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/lokasi?tipe=Provinsi", nil)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-1.1", FSDRef: "FSD-1.1",
		TSDRef: "TSD-3.3 (Endpoint Index — baris lokasi)", NoTestScript: "TC-LOKASI-001",
		Functional: "Get Provinsi — Success", Endpoint: "GET /lokasi?tipe=Provinsi",
		ReqType: "Query param", Parameter: "tipe=Provinsi",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data berisi 2 provinsi",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// GET /lokasi?tipe=Kabupaten&bagian_dari={id}
// ---------------------------------------------------------------------------

func TestLokasiGetKabupatenWithParentSuccess(t *testing.T) {
	f := setupLokasiIntegrationTest(t)
	defer f.cleanup(t)
	f.seed(t)

	path := "/lokasi?tipe=Kabupaten&bagian_dari=1"
	resp := testutils.DoRequest(f.handler, http.MethodGet, path, nil)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-1.1", FSDRef: "FSD-1.1",
		TSDRef: "TSD-3.3 (Endpoint Index — baris lokasi)", NoTestScript: "TC-LOKASI-002",
		Functional: "Get Kabupaten — With Parent Filter", Endpoint: "GET /lokasi?tipe=Kabupaten&bagian_dari=1",
		ReqType: "Query param", Parameter: "tipe=Kabupaten, bagian_dari=1",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data berisi 1 kabupaten",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// GET /lokasi?tipe=Kabupaten (tanpa bagian_dari)
// ---------------------------------------------------------------------------

func TestLokasiGetKabupatenWithoutParentEmpty(t *testing.T) {
	f := setupLokasiIntegrationTest(t)
	defer f.cleanup(t)
	f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/lokasi?tipe=Kabupaten", nil)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-1.1", FSDRef: "FSD-1.1",
		TSDRef: "TSD-3.3 (Endpoint Index — baris lokasi)", NoTestScript: "TC-LOKASI-003",
		Functional: "Get Kabupaten — Without Parent (IS NULL)", Endpoint: "GET /lokasi?tipe=Kabupaten",
		ReqType: "Query param", Parameter: "tipe=Kabupaten (tanpa bagian_dari)",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data empty array (semua kabupaten punya parent)",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// GET /lokasi?tipe=Kecamatan&bagian_dari={id}
// ---------------------------------------------------------------------------

func TestLokasiGetKecamatanWithParentSuccess(t *testing.T) {
	f := setupLokasiIntegrationTest(t)
	defer f.cleanup(t)
	f.seed(t)

	path := "/lokasi?tipe=Kecamatan&bagian_dari=3"
	resp := testutils.DoRequest(f.handler, http.MethodGet, path, nil)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-1.1", FSDRef: "FSD-1.1",
		TSDRef: "TSD-3.3 (Endpoint Index — baris lokasi)", NoTestScript: "TC-LOKASI-004",
		Functional: "Get Kecamatan — With Parent Filter", Endpoint: "GET /lokasi?tipe=Kecamatan&bagian_dari=3",
		ReqType: "Query param", Parameter: "tipe=Kecamatan, bagian_dari=3",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data berisi 1 kecamatan",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// GET /lokasi?tipe=Kelurahan&bagian_dari={id}
// ---------------------------------------------------------------------------

func TestLokasiGetKelurahanWithParentSuccess(t *testing.T) {
	f := setupLokasiIntegrationTest(t)
	defer f.cleanup(t)
	f.seed(t)

	path := "/lokasi?tipe=Kelurahan&bagian_dari=5"
	resp := testutils.DoRequest(f.handler, http.MethodGet, path, nil)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-1.1", FSDRef: "FSD-1.1",
		TSDRef: "TSD-3.3 (Endpoint Index — baris lokasi)", NoTestScript: "TC-LOKASI-005",
		Functional: "Get Kelurahan — With Parent Filter", Endpoint: "GET /lokasi?tipe=Kelurahan&bagian_dari=5",
		ReqType: "Query param", Parameter: "tipe=Kelurahan, bagian_dari=5",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data berisi 1 kelurahan",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// GET /lokasi?tipe=Kabupaten&bagian_dari=99999 (not found)
// ---------------------------------------------------------------------------

func TestLokasiGetEmptyResult(t *testing.T) {
	f := setupLokasiIntegrationTest(t)
	defer f.cleanup(t)
	f.seed(t)

	path := "/lokasi?tipe=Kabupaten&bagian_dari=99999"
	resp := testutils.DoRequest(f.handler, http.MethodGet, path, nil)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-1.1", FSDRef: "FSD-1.1",
		TSDRef: "TSD-3.3 (Endpoint Index — baris lokasi)", NoTestScript: "TC-LOKASI-006",
		Functional: "Get Lokasi — Empty Result (parent not found)", Endpoint: "GET /lokasi?tipe=Kabupaten&bagian_dari=99999",
		ReqType: "Query param", Parameter: "tipe=Kabupaten, bagian_dari=99999 (nonexistent)",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data empty array",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// GET /lokasi?tipe=Invalid (validation error)
// ---------------------------------------------------------------------------

func TestLokasiGetValidationErrorInvalidTipe(t *testing.T) {
	f := setupLokasiIntegrationTest(t)
	defer f.cleanup(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/lokasi?tipe=Invalid", nil)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusUnprocessableEntity

	testutils.TestResult{
		SRSRef: "SRS-1.1", FSDRef: "FSD-1.1",
		TSDRef: "TSD-3.3 (Endpoint Index — baris lokasi)", NoTestScript: "TC-LOKASI-007",
		Functional: "Get Lokasi — Validation Error (Invalid tipe)", Endpoint: "GET /lokasi?tipe=Invalid",
		ReqType: "Query param", Parameter: "tipe=Invalid (not in enum)",
		ShouldBeSuccess: "false",
		Expectation:     "Response 422, success:false, errors about invalid enum value",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// GET /lokasi (missing required tipe)
// ---------------------------------------------------------------------------

func TestLokasiGetValidationErrorMissingTipe(t *testing.T) {
	f := setupLokasiIntegrationTest(t)
	defer f.cleanup(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/lokasi", nil)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusUnprocessableEntity

	testutils.TestResult{
		SRSRef: "SRS-1.1", FSDRef: "FSD-1.1",
		TSDRef: "TSD-3.3 (Endpoint Index — baris lokasi)", NoTestScript: "TC-LOKASI-008",
		Functional: "Get Lokasi — Validation Error (Missing tipe)", Endpoint: "GET /lokasi",
		ReqType: "No query param", Parameter: "(tanpa query parameter)",
		ShouldBeSuccess: "false",
		Expectation:     fmt.Sprintf("Response %d, success:false", http.StatusUnprocessableEntity),
	}.Log(t, pass, resp, respBody)
}
