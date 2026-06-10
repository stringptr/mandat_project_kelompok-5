//go:build integration

package notification

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/stringptr/SiGizi/backend/internal/jwtutils"
	"github.com/stringptr/SiGizi/backend/internal/testutils"
)

type notifTestFixture struct {
	handler http.Handler
	jwtUtil *jwtutils.JWT
	pool    *pgxpool.Pool
}

func setupNotifIntegrationTest(t *testing.T) *notifTestFixture {
	t.Helper()

	pool := testutils.NewTestDB(t)
	handler, api, jwtUtil, br := testutils.SetupRouter(t)
	groups := testutils.CreateGroups(api, jwtUtil, br)

	notifRepo := NewRepo(pool)
	svc := NewService(notifRepo)
	h := NewHandler(svc)

	huma.Get(groups.BidanGroup, "/notifikasi/bidan", h.GetBidanDashboard)
	huma.Get(groups.AdminGroup, "/notifikasi/statistik", h.GetStatistics)
	huma.Get(groups.AdminGroup, "/notifikasi/aktivitas", h.GetActivity)
	huma.Get(groups.UserGroup, "/notifikasi", h.GetNotifikasi)
	huma.Get(groups.UserGroup, "/notifikasi/{id}", h.GetNotifikasiDetail)
	huma.Patch(groups.UserGroup, "/notifikasi/{id}/read", h.MarkRead)
	huma.Patch(groups.UserGroup, "/notifikasi/read-all", h.MarkAllRead)

	return &notifTestFixture{
		handler: handler,
		jwtUtil: jwtUtil,
		pool:    pool,
	}
}

func (f *notifTestFixture) cleanup(t *testing.T) {
	t.Helper()
	testutils.TruncateAuthTables(t, f.pool)
	testutils.TruncateNotifikasiTables(t, f.pool)
}

type notifSeedIDs struct {
	NotifikasiIDs []int32
}

func (f *notifTestFixture) seedNotif(t *testing.T, userID int32) *notifSeedIDs {
	t.Helper()
	ctx := context.Background()
	now := time.Now().Truncate(time.Microsecond)

	type notifData struct {
		judul string
		pesan string
		tipe  string
		read  bool
	}

	entries := []notifData{
		{"Jadwal Pemeriksaan", "Anda memiliki jadwal pemeriksaan rutin", "Pemeriksaan", false},
		{"Imunisasi Anak", "Imunisasi anak Anda sudah dekat", "Imunisasi", false},
		{"Rujukan ke Rumah Sakit", "Rujukan telah diajukan ke RSUD", "Rujukan", false},
		{"Edukasi Gizi", "Artikel edukasi gizi terbaru untuk ibu hamil", "Edukasi", true},
		{"Pengingat Kontrol", "Jangan lupa kontrol rutin ke posyandu", "Pengingat", false},
	}

	var ids []int32
	for _, e := range entries {
		var id int32
		err := f.pool.QueryRow(ctx, `
			INSERT INTO notifikasi (id_user, judul, pesan, tipe_notifikasi, status_baca, tanggal_kirim)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id_notifikasi`,
			userID, e.judul, e.pesan, e.tipe, e.read, now,
		).Scan(&id)
		if err != nil {
			t.Fatalf("failed to seed notifikasi '%s': %v", e.judul, err)
		}
		ids = append(ids, id)
	}

	return &notifSeedIDs{NotifikasiIDs: ids}
}

// ---------------------------------------------------------------------------
// TC-NOTIF-001: List notifications — success (has data)
// ---------------------------------------------------------------------------
func TestNotifikasiListSuccess(t *testing.T) {
	f := setupNotifIntegrationTest(t)
	defer f.cleanup(t)
	ids := testutils.SeedAuthData(t, f.pool)
	f.seedNotif(t, ids.RegularUserID)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/notifikasi", nil,
		testutils.AccessCookie(f.jwtUtil, ids.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated

	testutils.TestResult{
		SRSRef: "SRS-4.1", FSDRef: "FSD-4.1",
		TSDRef: "TSD-4.1-1", NoTestScript: "TC-NOTIF-001",
		Functional: "Notifikasi List — Success", Endpoint: "GET /notifikasi",
		ReqType: "Cookie (access_token)", Parameter: "Role: USER, 5 notifications seeded",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data contains notifikasi array and meta with total=5",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// TC-NOTIF-002: List notifications — empty (no data)
// ---------------------------------------------------------------------------
func TestNotifikasiListEmpty(t *testing.T) {
	f := setupNotifIntegrationTest(t)
	defer f.cleanup(t)
	ids := testutils.SeedAuthData(t, f.pool)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/notifikasi", nil,
		testutils.AccessCookie(f.jwtUtil, ids.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated

	testutils.TestResult{
		SRSRef: "SRS-4.1", FSDRef: "FSD-4.1",
		TSDRef: "TSD-4.1-1", NoTestScript: "TC-NOTIF-002",
		Functional: "Notifikasi List — Empty", Endpoint: "GET /notifikasi",
		ReqType: "Cookie (access_token)", Parameter: "Role: USER, no notifications seeded",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data.notifikasi empty, meta.total=0",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// TC-NOTIF-003: List notifications — search
// ---------------------------------------------------------------------------
func TestNotifikasiListSearch(t *testing.T) {
	f := setupNotifIntegrationTest(t)
	defer f.cleanup(t)
	ids := testutils.SeedAuthData(t, f.pool)
	f.seedNotif(t, ids.RegularUserID)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/notifikasi?q=Imunisasi", nil,
		testutils.AccessCookie(f.jwtUtil, ids.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated

	testutils.TestResult{
		SRSRef: "SRS-4.1", FSDRef: "FSD-4.1",
		TSDRef: "TSD-4.1-1", NoTestScript: "TC-NOTIF-003",
		Functional: "Notifikasi List — Search", Endpoint: "GET /notifikasi?q=Imunisasi",
		ReqType: "Cookie (access_token) + query", Parameter: "Role: USER, q=Imunisasi",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, data.notifikasi filtered to matching records",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// TC-NOTIF-004: List notifications — pagination
// ---------------------------------------------------------------------------
func TestNotifikasiListPagination(t *testing.T) {
	f := setupNotifIntegrationTest(t)
	defer f.cleanup(t)
	ids := testutils.SeedAuthData(t, f.pool)
	f.seedNotif(t, ids.RegularUserID)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/notifikasi?page=1&per_page=2", nil,
		testutils.AccessCookie(f.jwtUtil, ids.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated

	testutils.TestResult{
		SRSRef: "SRS-4.1", FSDRef: "FSD-4.1",
		TSDRef: "TSD-4.1-1", NoTestScript: "TC-NOTIF-004",
		Functional: "Notifikasi List — Pagination", Endpoint: "GET /notifikasi?page=1&per_page=2",
		ReqType: "Cookie (access_token) + query", Parameter: "Role: USER, page=1, per_page=2",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, data.notifikasi has ≤2 items, meta.current_page=1, meta.per_page=2",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// TC-NOTIF-005: Detail — success
// ---------------------------------------------------------------------------
func TestNotifikasiDetailSuccess(t *testing.T) {
	f := setupNotifIntegrationTest(t)
	defer f.cleanup(t)
	ids := testutils.SeedAuthData(t, f.pool)
	seed := f.seedNotif(t, ids.RegularUserID)

	path := fmt.Sprintf("/notifikasi/%d", seed.NotifikasiIDs[0])
	resp := testutils.DoRequest(f.handler, http.MethodGet, path, nil,
		testutils.AccessCookie(f.jwtUtil, ids.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated

	testutils.TestResult{
		SRSRef: "SRS-4.2", FSDRef: "FSD-4.2",
		TSDRef: "TSD-4.2-1", NoTestScript: "TC-NOTIF-005",
		Functional: "Notifikasi Detail — Success", Endpoint: "GET /notifikasi/{id}",
		ReqType: "Cookie (access_token) + path param", Parameter: "Role: USER, id=first notifikasi (Pemeriksaan)",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data contains detail with aksi",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// TC-NOTIF-006: Detail — not found
// ---------------------------------------------------------------------------
func TestNotifikasiDetailNotFound(t *testing.T) {
	f := setupNotifIntegrationTest(t)
	defer f.cleanup(t)
	ids := testutils.SeedAuthData(t, f.pool)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/notifikasi/99999", nil,
		testutils.AccessCookie(f.jwtUtil, ids.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusNotFound

	testutils.TestResult{
		SRSRef: "SRS-4.2", FSDRef: "FSD-4.2",
		TSDRef: "TSD-4.2-1", NoTestScript: "TC-NOTIF-006",
		Functional: "Notifikasi Detail — Not Found", Endpoint: "GET /notifikasi/{id}",
		ReqType: "Cookie (access_token) + path param", Parameter: "Role: USER, id=99999 (nonexistent)",
		ShouldBeSuccess: "false",
		Expectation:     "Response 404, success:false, detail: 'Notifikasi tidak ditemukan.'",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// TC-NOTIF-007: Mark read — success
// ---------------------------------------------------------------------------
func TestNotifikasiMarkReadSuccess(t *testing.T) {
	f := setupNotifIntegrationTest(t)
	defer f.cleanup(t)
	ids := testutils.SeedAuthData(t, f.pool)
	seed := f.seedNotif(t, ids.RegularUserID)

	path := fmt.Sprintf("/notifikasi/%d/read", seed.NotifikasiIDs[0])
	resp := testutils.DoRequest(f.handler, http.MethodPatch, path, nil,
		testutils.AccessCookie(f.jwtUtil, ids.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated

	testutils.TestResult{
		SRSRef: "SRS-4.3", FSDRef: "FSD-4.3",
		TSDRef: "TSD-4.3-1", NoTestScript: "TC-NOTIF-007",
		Functional: "Notifikasi Mark Read — Success", Endpoint: "PATCH /notifikasi/{id}/read",
		ReqType: "Cookie (access_token) + path param", Parameter: "Role: USER, id=first unread notifikasi",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data.status_baca:true",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// TC-NOTIF-008: Mark read — not found
// ---------------------------------------------------------------------------
func TestNotifikasiMarkReadNotFound(t *testing.T) {
	f := setupNotifIntegrationTest(t)
	defer f.cleanup(t)
	ids := testutils.SeedAuthData(t, f.pool)

	resp := testutils.DoRequest(f.handler, http.MethodPatch, "/notifikasi/99999/read", nil,
		testutils.AccessCookie(f.jwtUtil, ids.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusNotFound

	testutils.TestResult{
		SRSRef: "SRS-4.3", FSDRef: "FSD-4.3",
		TSDRef: "TSD-4.3-1", NoTestScript: "TC-NOTIF-008",
		Functional: "Notifikasi Mark Read — Not Found", Endpoint: "PATCH /notifikasi/{id}/read",
		ReqType: "Cookie (access_token) + path param", Parameter: "Role: USER, id=99999 (nonexistent)",
		ShouldBeSuccess: "false",
		Expectation:     "Response 404, success:false, detail: 'Notifikasi tidak ditemukan.'",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// TC-NOTIF-009: Mark all read — success (has unread)
// ---------------------------------------------------------------------------
func TestNotifikasiMarkAllReadSuccess(t *testing.T) {
	f := setupNotifIntegrationTest(t)
	defer f.cleanup(t)
	ids := testutils.SeedAuthData(t, f.pool)
	f.seedNotif(t, ids.RegularUserID)

	resp := testutils.DoRequest(f.handler, http.MethodPatch, "/notifikasi/read-all", nil,
		testutils.AccessCookie(f.jwtUtil, ids.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated

	testutils.TestResult{
		SRSRef: "SRS-4.3", FSDRef: "FSD-4.3",
		TSDRef: "TSD-4.3-2", NoTestScript: "TC-NOTIF-009",
		Functional: "Notifikasi Mark All Read — Success", Endpoint: "PATCH /notifikasi/read-all",
		ReqType: "Cookie (access_token)", Parameter: "Role: USER, 4 unread + 1 read seeded",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data.jumlah_diperbarui=4, data.status='SEMUA_DIBACA'",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// TC-NOTIF-010: Mark all read — no unread
// ---------------------------------------------------------------------------
func TestNotifikasiMarkAllReadNoUnread(t *testing.T) {
	f := setupNotifIntegrationTest(t)
	defer f.cleanup(t)
	ids := testutils.SeedAuthData(t, f.pool)

	resp := testutils.DoRequest(f.handler, http.MethodPatch, "/notifikasi/read-all", nil,
		testutils.AccessCookie(f.jwtUtil, ids.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated

	testutils.TestResult{
		SRSRef: "SRS-4.3", FSDRef: "FSD-4.3",
		TSDRef: "TSD-4.3-2", NoTestScript: "TC-NOTIF-010",
		Functional: "Notifikasi Mark All Read — No Unread", Endpoint: "PATCH /notifikasi/read-all",
		ReqType: "Cookie (access_token)", Parameter: "Role: USER, no notifications seeded",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data.jumlah_diperbarui=0",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// TC-NOTIF-011: Unauthorized — no token
// ---------------------------------------------------------------------------
func TestNotifikasiUnauthorized(t *testing.T) {
	f := setupNotifIntegrationTest(t)
	defer f.cleanup(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/notifikasi", nil)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusUnauthorized

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3-5", NoTestScript: "TC-NOTIF-011",
		Functional: "Notifikasi — Unauthorized (No Token)", Endpoint: "GET /notifikasi",
		ReqType: "No Cookie", Parameter: "{}",
		ShouldBeSuccess: "false",
		Expectation:     "Response 401, success:false",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// TC-NOTIF-012: Forbidden — wrong role (USER trying BIDAN endpoint)
// ---------------------------------------------------------------------------
func TestNotifikasiForbiddenRole(t *testing.T) {
	f := setupNotifIntegrationTest(t)
	defer f.cleanup(t)
	ids := testutils.SeedAuthData(t, f.pool)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/notifikasi/bidan", nil,
		testutils.AccessCookie(f.jwtUtil, ids.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusForbidden

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3-5", NoTestScript: "TC-NOTIF-012",
		Functional: "Notifikasi — Forbidden (Wrong Role)", Endpoint: "GET /notifikasi/bidan",
		ReqType: "Cookie (access_token)", Parameter: "Role: USER (requires BIDAN)",
		ShouldBeSuccess: "false",
		Expectation:     "Response 403, success:false",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// TC-NOTIF-013: Bidan dashboard — success
// ---------------------------------------------------------------------------
func TestNotifikasiBidanDashboard(t *testing.T) {
	f := setupNotifIntegrationTest(t)
	defer f.cleanup(t)
	ids := testutils.SeedAuthData(t, f.pool)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/notifikasi/bidan", nil,
		testutils.AccessCookie(f.jwtUtil, ids.VerifiedUserID, []string{"BIDAN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated

	testutils.TestResult{
		SRSRef: "SRS-4.5", FSDRef: "FSD-4.5",
		TSDRef: "TSD-4.5-1", NoTestScript: "TC-NOTIF-013",
		Functional: "Notifikasi Bidan Dashboard — Success", Endpoint: "GET /notifikasi/bidan",
		ReqType: "Cookie (access_token)", Parameter: "Role: BIDAN (verifiedUser with bidan row)",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data contains statistik and lists",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// TC-NOTIF-014: Statistics — admin success
// ---------------------------------------------------------------------------
func TestNotifikasiStatistics(t *testing.T) {
	f := setupNotifIntegrationTest(t)
	defer f.cleanup(t)
	ids := testutils.SeedAuthData(t, f.pool)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/notifikasi/statistik", nil,
		testutils.AccessCookie(f.jwtUtil, ids.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated

	testutils.TestResult{
		SRSRef: "SRS-4.4", FSDRef: "FSD-4.4",
		TSDRef: "TSD-4.4-1", NoTestScript: "TC-NOTIF-014",
		Functional: "Notifikasi Statistics — Success", Endpoint: "GET /notifikasi/statistik",
		ReqType: "Cookie (access_token)", Parameter: "Role: ADMIN (adminUser with dinkes row)",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data contains notification stats",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// TC-NOTIF-015: Activity — with data
// ---------------------------------------------------------------------------
func TestNotifikasiActivityWithData(t *testing.T) {
	f := setupNotifIntegrationTest(t)
	defer f.cleanup(t)
	ids := testutils.SeedAuthData(t, f.pool)
	f.seedNotif(t, ids.AdminUserID)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/notifikasi/aktivitas", nil,
		testutils.AccessCookie(f.jwtUtil, ids.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated

	testutils.TestResult{
		SRSRef: "SRS-4.4", FSDRef: "FSD-4.4",
		TSDRef: "TSD-4.4-2", NoTestScript: "TC-NOTIF-015",
		Functional: "Notifikasi Activity — With Data", Endpoint: "GET /notifikasi/aktivitas",
		ReqType: "Cookie (access_token)", Parameter: "Role: ADMIN, 5 notifications seeded for adminUser",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data.hari_ini has items, data.kemarin may be empty",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// TC-NOTIF-016: Activity — empty (no data)
// ---------------------------------------------------------------------------
func TestNotifikasiActivityEmpty(t *testing.T) {
	f := setupNotifIntegrationTest(t)
	defer f.cleanup(t)
	ids := testutils.SeedAuthData(t, f.pool)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/notifikasi/aktivitas", nil,
		testutils.AccessCookie(f.jwtUtil, ids.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated

	testutils.TestResult{
		SRSRef: "SRS-4.4", FSDRef: "FSD-4.4",
		TSDRef: "TSD-4.4-2", NoTestScript: "TC-NOTIF-016",
		Functional: "Notifikasi Activity — Empty", Endpoint: "GET /notifikasi/aktivitas",
		ReqType: "Cookie (access_token)", Parameter: "Role: ADMIN, no notifications seeded",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data.hari_ini empty, data.kemarin empty",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// TC-NOTIF-017: Statistics — forbidden (USER trying ADMIN endpoint)
// ---------------------------------------------------------------------------
func TestNotifikasiStatisticsForbidden(t *testing.T) {
	f := setupNotifIntegrationTest(t)
	defer f.cleanup(t)
	ids := testutils.SeedAuthData(t, f.pool)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/notifikasi/statistik", nil,
		testutils.AccessCookie(f.jwtUtil, ids.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusForbidden

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3-5", NoTestScript: "TC-NOTIF-017",
		Functional: "Notifikasi Statistics — Forbidden (Wrong Role)", Endpoint: "GET /notifikasi/statistik",
		ReqType: "Cookie (access_token)", Parameter: "Role: USER (requires ADMIN)",
		ShouldBeSuccess: "false",
		Expectation:     "Response 403, success:false",
	}.Log(t, pass, resp, respBody)
}
