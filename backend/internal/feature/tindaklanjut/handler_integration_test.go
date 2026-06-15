//go:build integration

package tindaklanjut

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

type tindaklanjutTestFixture struct {
	handler http.Handler
	jwtUtil *jwtutils.JWT
	pool    *pgxpool.Pool
}

type tindaklanjutSeedIDs struct {
	PosyanduID         int32
	PasienID           int32
	JadwalImunisasiID  int32
	HasilPemeriksaanID int32
	FaskesID           int32
	TindakLanjutID     int32
	RujukanID          int32
}

func setupTindakLanjutIntegrationTest(t *testing.T) *tindaklanjutTestFixture {
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

	huma.Get(groups.BidanGroup, "/tindak-lanjut/pasien", h.GetPasienTindakLanjut)
	huma.Get(groups.BidanGroup, "/tindak-lanjut/pasien/{id}", h.GetDetailPasienByID)
	huma.Post(groups.BidanGroup, "/tindak-lanjut", h.CreateTindakLanjut)
	huma.Patch(groups.BidanGroup, "/rujukan/{id}/status", h.UpdateStatusRujukan)
	huma.Get(groups.BidanGroup, "/tindak-lanjut/status", h.GetStatusTindakLanjut)
	huma.Get(groups.DinkesGroup, "/laporan/tindak-lanjut", h.GetLaporanTindakLanjut)
	huma.Get(groups.UserGroup, "/tindak-lanjut/{id}", h.GetDetailTindakLanjutByID)

	return &tindaklanjutTestFixture{
		handler: handler,
		jwtUtil: jwtUtil,
		pool:    pool,
	}
}

func (f *tindaklanjutTestFixture) cleanup(t *testing.T) {
	t.Helper()
	testutils.TruncateTindakLanjutTables(t, f.pool)
	testutils.TruncateImunisasiTables(t, f.pool)
	testutils.TruncateNotifikasiTables(t, f.pool)
	testutils.TruncatePasienTables(t, f.pool)
	testutils.TruncateAuthTables(t, f.pool)
}

func (f *tindaklanjutTestFixture) seed(t *testing.T) *testutils.AuthSeedIDs {
	t.Helper()
	return testutils.SeedAuthData(t, f.pool)
}

func (f *tindaklanjutTestFixture) seedTindakLanjut(t *testing.T, authIDs *testutils.AuthSeedIDs) *tindaklanjutSeedIDs {
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
		authIDs.AdminUserID, jadwalID, 8.0, 65.0, 42.0, "120/80", "Stunting", "Gizi Kurang", "Perlu pemantauan", now,
	).Scan(&hasilID)
	if err != nil {
		t.Fatalf("failed to seed hasil_pemeriksaan: %v", err)
	}

	faskesID := testutils.SeedFasilitasKesehatan(t, f.pool, 1)

	var tindakLanjutID int32
	err = f.pool.QueryRow(ctx, `
		INSERT INTO tindak_lanjut (id_hasil_pemeriksaan, id_bidan, catatan_medis, rekomendasi, jadwal_kontrol, status_pasien, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $7)
		RETURNING id_tindak_lanjut`,
		hasilID, authIDs.VerifiedUserID, "Catatan medis", "Rekomendasi kontrol", now.AddDate(0, 0, 14), "Dalam Pemantauan", now,
	).Scan(&tindakLanjutID)
	if err != nil {
		t.Fatalf("failed to seed tindak_lanjut: %v", err)
	}

	var rujukanID int32
	err = f.pool.QueryRow(ctx, `
		INSERT INTO rujukan (id_tindak_lanjut, alasan_rujukan, tanggal_rujukan, status_rujukan, id_faskes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $6)
		RETURNING id_rujukan`,
		tindakLanjutID, "Stunting berat", now, "Diajukan", faskesID, now,
	).Scan(&rujukanID)
	if err != nil {
		t.Fatalf("failed to seed rujukan: %v", err)
	}

	return &tindaklanjutSeedIDs{
		PosyanduID:         posyanduID,
		PasienID:           pasienID,
		JadwalImunisasiID:  jadwalID,
		HasilPemeriksaanID: hasilID,
		FaskesID:           faskesID,
		TindakLanjutID:     tindakLanjutID,
		RujukanID:          rujukanID,
	}
}

// seedPasienTanpaTindakLanjut creates a pasien with pemeriksaan but NO tindak_lanjut
// so they appear in the "perlu tindak lanjut" list.
func (f *tindaklanjutTestFixture) seedPasienTanpaTindakLanjut(t *testing.T, authIDs *testutils.AuthSeedIDs) int32 {
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
		pasienID, "Polio", now.AddDate(0, 0, 30), "Belum", now,
	).Scan(&jadwalID)
	if err != nil {
		t.Fatalf("failed to seed jadwal_imunisasi: %v", err)
	}

	var hasilID int32
	err = f.pool.QueryRow(ctx, `
		INSERT INTO hasil_pemeriksaan (id_petugas_input, id_jadwal_imunisasi, berat_badan, tinggi_badan, lingkar_kepala, tekanan_darah, status_stunting, status_gizi, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $9)
		RETURNING id_hasil_pemeriksaan`,
		authIDs.AdminUserID, jadwalID, 7.5, 62.0, 40.0, "110/70", "Stunting", "Gizi Buruk", now,
	).Scan(&hasilID)
	if err != nil {
		t.Fatalf("failed to seed hasil_pemeriksaan: %v", err)
	}

	_ = hasilID // not stored, this pasien has no follow-up yet
	return pasienID
}

// ---------------------------------------------------------------------------
// GET /tindak-lanjut/pasien
// ---------------------------------------------------------------------------

func TestTindakLanjutGetPasienWithData(t *testing.T) {
	f := setupTindakLanjutIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	f.seedPasienTanpaTindakLanjut(t, authIDs)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/tindak-lanjut/pasien", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.VerifiedUserID, []string{"BIDAN", "USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-5.1", FSDRef: "FSD-5.1",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-TL-001",
		Functional: "Tindak Lanjut Pasien List — With Data", Endpoint: "GET /tindak-lanjut/pasien",
		ReqType: "Cookie (BIDAN)", Parameter: "Role: BIDAN, pasien without follow-up seeded",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data.pasien has items",
	}.Log(t, pass, resp, respBody)
}

func TestTindakLanjutGetPasienEmpty(t *testing.T) {
	f := setupTindakLanjutIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/tindak-lanjut/pasien", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.VerifiedUserID, []string{"BIDAN", "USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-5.1", FSDRef: "FSD-5.1",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-TL-002",
		Functional: "Tindak Lanjut Pasien List — Empty", Endpoint: "GET /tindak-lanjut/pasien",
		ReqType: "Cookie (BIDAN)", Parameter: "Role: BIDAN, no eligible pasien",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data.pasien empty",
	}.Log(t, pass, resp, respBody)
}

func TestTindakLanjutGetPasienForbidden(t *testing.T) {
	f := setupTindakLanjutIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/tindak-lanjut/pasien", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusForbidden

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-TL-003",
		Functional: "Tindak Lanjut Pasien List — Forbidden (USER)", Endpoint: "GET /tindak-lanjut/pasien",
		ReqType: "Cookie (USER)", Parameter: "Role: USER (requires BIDAN)",
		ShouldBeSuccess: "false",
		Expectation:     "Response 403, success:false",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// GET /tindak-lanjut/pasien/{id}
// ---------------------------------------------------------------------------

func TestTindakLanjutGetDetailPasienSuccess(t *testing.T) {
	f := setupTindakLanjutIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	f.seedTindakLanjut(t, authIDs)

	path := fmt.Sprintf("/tindak-lanjut/pasien/%d", authIDs.RegularUserID)
	resp := testutils.DoRequest(f.handler, http.MethodGet, path, nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.VerifiedUserID, []string{"BIDAN", "USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-5.2", FSDRef: "FSD-5.2",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-TL-004",
		Functional: "Tindak Lanjut Detail Pasien — Success", Endpoint: "GET /tindak-lanjut/pasien/{id}",
		ReqType: "Cookie (BIDAN) + path param", Parameter: "Role: BIDAN, id=existing",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data contains detail",
	}.Log(t, pass, resp, respBody)
}

func TestTindakLanjutGetDetailPasienNotFound(t *testing.T) {
	f := setupTindakLanjutIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/tindak-lanjut/pasien/99999", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.VerifiedUserID, []string{"BIDAN", "USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusNotFound

	testutils.TestResult{
		SRSRef: "SRS-5.2", FSDRef: "FSD-5.2",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-TL-005",
		Functional: "Tindak Lanjut Detail Pasien — Not Found", Endpoint: "GET /tindak-lanjut/pasien/{id}",
		ReqType: "Cookie (BIDAN) + path param", Parameter: "id=99999",
		ShouldBeSuccess: "false",
		Expectation:     "Response 404, success:false",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// POST /tindak-lanjut
// ---------------------------------------------------------------------------

func TestTindakLanjutCreateSuccess(t *testing.T) {
	f := setupTindakLanjutIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedTindakLanjut(t, authIDs)

	// Need a fresh hasil_pemeriksaan without existing tindak_lanjut
	ctx := context.Background()
	now := time.Now().Truncate(time.Microsecond)
	var newHasilID int32
	err := f.pool.QueryRow(ctx, `
		INSERT INTO hasil_pemeriksaan (id_petugas_input, id_jadwal_imunisasi, berat_badan, tinggi_badan, lingkar_kepala, tekanan_darah, status_stunting, status_gizi, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $9)
		RETURNING id_hasil_pemeriksaan`,
		authIDs.AdminUserID, seed.JadwalImunisasiID, 8.0, 65.0, 42.0, "120/80", "Stunting", "Gizi Kurang", now,
	).Scan(&newHasilID)
	if err != nil {
		t.Fatalf("failed to seed additional hasil_pemeriksaan: %v", err)
	}

	body := map[string]any{
		"id_hasil_pemeriksaan": newHasilID,
		"jenis_tindakan":       "Kontrol Ulang",
		"catatan_medis":        "Pantau perkembangan",
		"rekomendasi":          "Kontrol 2 minggu lagi",
		"jadwal_kontrol":       now.AddDate(0, 0, 14).Format("2006-01-02"),
	}
	resp := testutils.DoRequest(f.handler, http.MethodPost, "/tindak-lanjut", body,
		testutils.AccessCookie(f.jwtUtil, authIDs.VerifiedUserID, []string{"BIDAN", "USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusCreated

	testutils.TestResult{
		SRSRef: "SRS-5.3", FSDRef: "FSD-5.3",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-TL-006",
		Functional: "Tindak Lanjut Create — Success", Endpoint: "POST /tindak-lanjut",
		ReqType: "JSON Body + Cookie (BIDAN)", Parameter: `{"jenis_tindakan":"Kontrol Ulang"}`,
		ShouldBeSuccess: "true",
		Expectation:     "Response 201, success:true, data contains id_tindak_lanjut",
	}.Log(t, pass, resp, respBody)
}

func TestTindakLanjutCreateWithRujukanSuccess(t *testing.T) {
	f := setupTindakLanjutIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedTindakLanjut(t, authIDs)

	ctx := context.Background()
	now := time.Now().Truncate(time.Microsecond)
	var newHasilID int32
	err := f.pool.QueryRow(ctx, `
		INSERT INTO hasil_pemeriksaan (id_petugas_input, id_jadwal_imunisasi, berat_badan, tinggi_badan, lingkar_kepala, tekanan_darah, status_stunting, status_gizi, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $9)
		RETURNING id_hasil_pemeriksaan`,
		authIDs.AdminUserID, seed.JadwalImunisasiID, 7.0, 60.0, 40.0, "110/70", "Stunting Berat", "Gizi Buruk", now,
	).Scan(&newHasilID)
	if err != nil {
		t.Fatalf("failed to seed additional hasil_pemeriksaan: %v", err)
	}

	body := map[string]any{
		"id_hasil_pemeriksaan": newHasilID,
		"jenis_tindakan":       "Rujukan",
		"catatan_medis":        "Rujukan ke puskesmas",
		"rekomendasi":          "Penanganan lebih lanjut",
		"alasan_rujukan":       "Stunting berat",
		"id_faskes":            seed.FaskesID,
	}
	resp := testutils.DoRequest(f.handler, http.MethodPost, "/tindak-lanjut", body,
		testutils.AccessCookie(f.jwtUtil, authIDs.VerifiedUserID, []string{"BIDAN", "USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusCreated

	testutils.TestResult{
		SRSRef: "SRS-5.3", FSDRef: "FSD-5.3",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-TL-007",
		Functional: "Tindak Lanjut Create — With Rujukan", Endpoint: "POST /tindak-lanjut",
		ReqType: "JSON Body + Cookie (BIDAN)", Parameter: `{"jenis_tindakan":"Rujukan","alasan_rujukan":"Stunting berat","id_faskes":"<ID>"}`,
		ShouldBeSuccess: "true",
		Expectation:     "Response 201, success:true, data contains id_rujukan",
	}.Log(t, pass, resp, respBody)
}

func TestTindakLanjutCreateConflict(t *testing.T) {
	f := setupTindakLanjutIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedTindakLanjut(t, authIDs)

	body := map[string]any{
		"id_hasil_pemeriksaan": seed.HasilPemeriksaanID,
		"jenis_tindakan":       "Kontrol Ulang",
	}
	resp := testutils.DoRequest(f.handler, http.MethodPost, "/tindak-lanjut", body,
		testutils.AccessCookie(f.jwtUtil, authIDs.VerifiedUserID, []string{"BIDAN", "USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusConflict

	testutils.TestResult{
		SRSRef: "SRS-5.3", FSDRef: "FSD-5.3",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-TL-008",
		Functional: "Tindak Lanjut Create — Conflict (Duplicate)", Endpoint: "POST /tindak-lanjut",
		ReqType: "JSON Body + Cookie (BIDAN)", Parameter: "id_hasil_pemeriksaan already has follow-up",
		ShouldBeSuccess: "false",
		Expectation:     "Response 409, success:false",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// PATCH /rujukan/{id}/status
// ---------------------------------------------------------------------------

func TestTindakLanjutUpdateStatusRujukanSuccess(t *testing.T) {
	f := setupTindakLanjutIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedTindakLanjut(t, authIDs)

	body := map[string]any{"status_rujukan": "Diproses"}
	path := fmt.Sprintf("/rujukan/%d/status", seed.RujukanID)
	resp := testutils.DoRequest(f.handler, http.MethodPatch, path, body,
		testutils.AccessCookie(f.jwtUtil, authIDs.VerifiedUserID, []string{"BIDAN", "USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-5.4", FSDRef: "FSD-5.4",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-TL-009",
		Functional: "Rujukan Update Status — Success", Endpoint: "PATCH /rujukan/{id}/status",
		ReqType: "JSON Body + Cookie (BIDAN)", Parameter: `{"status_rujukan":"Diproses"}`,
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data.status_rujukan='Diproses'",
	}.Log(t, pass, resp, respBody)
}

func TestTindakLanjutUpdateStatusRujukanInvalid(t *testing.T) {
	f := setupTindakLanjutIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedTindakLanjut(t, authIDs)

	body := map[string]any{"status_rujukan": "InvalidStatus"}
	path := fmt.Sprintf("/rujukan/%d/status", seed.RujukanID)
	resp := testutils.DoRequest(f.handler, http.MethodPatch, path, body,
		testutils.AccessCookie(f.jwtUtil, authIDs.VerifiedUserID, []string{"BIDAN", "USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusUnprocessableEntity

	testutils.TestResult{
		SRSRef: "SRS-5.4", FSDRef: "FSD-5.4",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-TL-010",
		Functional: "Rujukan Update Status — Invalid", Endpoint: "PATCH /rujukan/{id}/status",
		ReqType: "JSON Body + Cookie (BIDAN)", Parameter: `{"status_rujukan":"InvalidStatus"}`,
		ShouldBeSuccess: "false",
		Expectation:     "Response 422, success:false (validated by Huma enum)",
	}.Log(t, pass, resp, respBody)
}

func TestTindakLanjutUpdateStatusRujukanNotFound(t *testing.T) {
	f := setupTindakLanjutIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	body := map[string]any{"status_rujukan": "Diproses"}
	resp := testutils.DoRequest(f.handler, http.MethodPatch, "/rujukan/99999/status", body,
		testutils.AccessCookie(f.jwtUtil, authIDs.VerifiedUserID, []string{"BIDAN", "USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusNotFound

	testutils.TestResult{
		SRSRef: "SRS-5.4", FSDRef: "FSD-5.4",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-TL-011",
		Functional: "Rujukan Update Status — Not Found", Endpoint: "PATCH /rujukan/{id}/status",
		ReqType: "JSON Body + Cookie (BIDAN)", Parameter: "id=99999",
		ShouldBeSuccess: "false",
		Expectation:     "Response 404, success:false",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// GET /tindak-lanjut/status
// ---------------------------------------------------------------------------

func TestTindakLanjutGetStatusWithData(t *testing.T) {
	f := setupTindakLanjutIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	f.seedTindakLanjut(t, authIDs)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/tindak-lanjut/status", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.VerifiedUserID, []string{"BIDAN", "USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-5.5", FSDRef: "FSD-5.5",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-TL-012",
		Functional: "Tindak Lanjut Status List — With Data", Endpoint: "GET /tindak-lanjut/status",
		ReqType: "Cookie (BIDAN)", Parameter: "Role: BIDAN, data seeded",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data.pasien has items",
	}.Log(t, pass, resp, respBody)
}

func TestTindakLanjutGetStatusEmpty(t *testing.T) {
	f := setupTindakLanjutIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/tindak-lanjut/status", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.VerifiedUserID, []string{"BIDAN", "USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-5.5", FSDRef: "FSD-5.5",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-TL-013",
		Functional: "Tindak Lanjut Status List — Empty", Endpoint: "GET /tindak-lanjut/status",
		ReqType: "Cookie (BIDAN)", Parameter: "Role: BIDAN, no data",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data.pasien empty",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// GET /laporan/tindak-lanjut
// ---------------------------------------------------------------------------

func TestTindakLanjutGetLaporanWithData(t *testing.T) {
	f := setupTindakLanjutIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	f.seedTindakLanjut(t, authIDs)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/laporan/tindak-lanjut", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"USER", "DINKES", "SUPER_ADMIN", "ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-5.6", FSDRef: "FSD-5.6",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-TL-014",
		Functional: "Laporan Tindak Lanjut — With Data", Endpoint: "GET /laporan/tindak-lanjut",
		ReqType: "Cookie (DINKES+ADMIN)", Parameter: "Role: DINKES+ADMIN, data seeded",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data.laporan has items",
	}.Log(t, pass, resp, respBody)
}

func TestTindakLanjutGetLaporanEmpty(t *testing.T) {
	f := setupTindakLanjutIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/laporan/tindak-lanjut", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"USER", "DINKES", "SUPER_ADMIN", "ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-5.6", FSDRef: "FSD-5.6",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-TL-015",
		Functional: "Laporan Tindak Lanjut — Empty", Endpoint: "GET /laporan/tindak-lanjut",
		ReqType: "Cookie (DINKES+ADMIN)", Parameter: "Role: DINKES+ADMIN, no data",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data.laporan empty",
	}.Log(t, pass, resp, respBody)
}

func TestTindakLanjutGetLaporanForbidden(t *testing.T) {
	f := setupTindakLanjutIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/laporan/tindak-lanjut", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusForbidden

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-TL-016",
		Functional: "Laporan Tindak Lanjut — Forbidden (USER)", Endpoint: "GET /laporan/tindak-lanjut",
		ReqType: "Cookie (USER)", Parameter: "Role: USER (requires DINKES)",
		ShouldBeSuccess: "false",
		Expectation:     "Response 403, success:false",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// GET /tindak-lanjut/{id}
// ---------------------------------------------------------------------------

func TestTindakLanjutGetDetailByIDSuccess(t *testing.T) {
	f := setupTindakLanjutIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedTindakLanjut(t, authIDs)

	path := fmt.Sprintf("/tindak-lanjut/%d", seed.TindakLanjutID)
	resp := testutils.DoRequest(f.handler, http.MethodGet, path, nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-5.7", FSDRef: "FSD-5.7",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-TL-017",
		Functional: "Tindak Lanjut Detail By ID — Success", Endpoint: "GET /tindak-lanjut/{id}",
		ReqType: "Cookie (USER) + path param", Parameter: "Role: USER, id=existing",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data contains detail",
	}.Log(t, pass, resp, respBody)
}

func TestTindakLanjutGetDetailByIDNotFound(t *testing.T) {
	f := setupTindakLanjutIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/tindak-lanjut/99999", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusNotFound

	testutils.TestResult{
		SRSRef: "SRS-5.7", FSDRef: "FSD-5.7",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-TL-018",
		Functional: "Tindak Lanjut Detail By ID — Not Found", Endpoint: "GET /tindak-lanjut/{id}",
		ReqType: "Cookie (USER) + path param", Parameter: "id=99999",
		ShouldBeSuccess: "false",
		Expectation:     "Response 404, success:false",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// Unauthorized
// ---------------------------------------------------------------------------

func TestTindakLanjutUnauthorized(t *testing.T) {
	f := setupTindakLanjutIntegrationTest(t)
	defer f.cleanup(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/tindak-lanjut/pasien", nil)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusUnauthorized

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-TL-019",
		Functional: "Tindak Lanjut — Unauthorized (No Token)", Endpoint: "GET /tindak-lanjut/pasien",
		ReqType: "No Cookie", Parameter: "{}",
		ShouldBeSuccess: "false",
		Expectation:     "Response 401, success:false",
	}.Log(t, pass, resp, respBody)
}
