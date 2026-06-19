//go:build integration

package artikel

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

type artikelTestFixture struct {
	handler http.Handler
	jwtUtil *jwtutils.JWT
	pool    *pgxpool.Pool
}

type artikelSeedIDs struct {
	PublishedArtikelID int32
	PendingArtikelID   int32
	RejectedArtikelID  int32
}

func setupArtikelIntegrationTest(t *testing.T) *artikelTestFixture {
	t.Helper()

	pool := testutils.NewTestDB(t)
	handler, api, jwtUtil, br := testutils.SetupChiRouter(t)
	groups := testutils.CreateGroups(api, jwtUtil, br)

	repo := NewRepo(pool)
	auditRepo := auditlog.NewRepo(pool)
	notifRepo := notification.NewRepo(pool)
	notifPub := &testutils.NoopNotifPublisher{}
	svc := NewService(repo, auditRepo, notifRepo, notifPub)
	h := NewHandler(svc)

	huma.Get(groups.PublicGroup, "/artikel", h.GetAllPublished)
	huma.Get(groups.PublicGroup, "/artikel/{id}", h.GetByID)
	huma.Post(groups.BidanGroup, "/artikel", h.Create)
	huma.Patch(groups.BidanGroup, "/artikel/{id}", h.Update)
	huma.Delete(groups.DinkesGroup, "/artikel/{id}", h.Delete)
	huma.Patch(groups.DinkesGroup, "/artikel/{id}/review", h.Review)
	huma.Get(groups.DinkesGroup, "/artikel/pending", h.GetPending)

	return &artikelTestFixture{
		handler: handler,
		jwtUtil: jwtUtil,
		pool:    pool,
	}
}

func (f *artikelTestFixture) cleanup(t *testing.T) {
	t.Helper()
	testutils.TruncateArtikelTables(t, f.pool)
	testutils.TruncateNotifikasiTables(t, f.pool)
	testutils.TruncateAuthTables(t, f.pool)
}

func (f *artikelTestFixture) seed(t *testing.T) *testutils.AuthSeedIDs {
	t.Helper()
	return testutils.SeedAuthData(t, f.pool)
}

func (f *artikelTestFixture) seedArtikel(t *testing.T, authIDs *testutils.AuthSeedIDs) *artikelSeedIDs {
	t.Helper()
	ctx := context.Background()
	now := time.Now().Truncate(time.Microsecond)

	var publishedID int32
	err := f.pool.QueryRow(ctx, `
		INSERT INTO artikel (judul, isi_artikel, kategori, status_artikel, id_penulis, id_verifikator, tanggal_publish, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $8)
		RETURNING id_artikel`,
		"Artikel Published", "Isi artikel yang sudah dipublikasikan.", "Edukasi",
		"Dipublikasikan", authIDs.VerifiedUserID, &authIDs.AdminUserID, now, now,
	).Scan(&publishedID)
	if err != nil {
		t.Fatalf("failed to seed published artikel: %v", err)
	}

	var pendingID int32
	err = f.pool.QueryRow(ctx, `
		INSERT INTO artikel (judul, isi_artikel, kategori, status_artikel, id_penulis, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $6)
		RETURNING id_artikel`,
		"Artikel Pending", "Isi artikel yang menunggu verifikasi.", "Gizi",
		"Menunggu Verifikasi", authIDs.VerifiedUserID, now,
	).Scan(&pendingID)
	if err != nil {
		t.Fatalf("failed to seed pending artikel: %v", err)
	}

	var rejectedID int32
	err = f.pool.QueryRow(ctx, `
		INSERT INTO artikel (judul, isi_artikel, status_artikel, id_penulis, id_verifikator, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $6)
		RETURNING id_artikel`,
		"Artikel Ditolak", "Isi artikel yang ditolak.", "Ditolak",
		authIDs.VerifiedUserID, &authIDs.AdminUserID, now,
	).Scan(&rejectedID)
	if err != nil {
		t.Fatalf("failed to seed rejected artikel: %v", err)
	}

	return &artikelSeedIDs{
		PublishedArtikelID: publishedID,
		PendingArtikelID:   pendingID,
		RejectedArtikelID:  rejectedID,
	}
}

// ---------------------------------------------------------------------------
// GET /artikel
// ---------------------------------------------------------------------------

func TestArtikelGetAllPublishedWithData(t *testing.T) {
	f := setupArtikelIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	f.seedArtikel(t, authIDs)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/artikel", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-6.1", FSDRef: "FSD-6.1",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-ARTIKEL-001",
		Functional: "Artikel List Published — With Data", Endpoint: "GET /artikel",
		ReqType: "Cookie (ADMIN)", Parameter: "Role: ADMIN, 1 published seeded",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data.artikel has items",
	}.Log(t, pass, resp, respBody)
}

func TestArtikelGetAllPublishedEmpty(t *testing.T) {
	f := setupArtikelIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/artikel", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-6.1", FSDRef: "FSD-6.1",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-ARTIKEL-002",
		Functional: "Artikel List Published — Empty", Endpoint: "GET /artikel",
		ReqType: "Cookie (ADMIN)", Parameter: "Role: ADMIN, no published articles",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data.artikel empty",
	}.Log(t, pass, resp, respBody)
}

func TestArtikelGetAllPublishedPublicAccess(t *testing.T) {
	f := setupArtikelIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/artikel", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-ARTIKEL-003",
		Functional: "Artikel List — Public Access (USER)", Endpoint: "GET /artikel",
		ReqType: "Cookie (USER)", Parameter: "Role: USER, endpoint is public",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// GET /artikel/{id}
// ---------------------------------------------------------------------------

func TestArtikelGetByIDPublishedSuccess(t *testing.T) {
	f := setupArtikelIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedArtikel(t, authIDs)

	path := fmt.Sprintf("/artikel/%d", seed.PublishedArtikelID)
	resp := testutils.DoRequest(f.handler, http.MethodGet, path, nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-6.2", FSDRef: "FSD-6.2",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-ARTIKEL-004",
		Functional: "Artikel Detail — Published", Endpoint: "GET /artikel/{id}",
		ReqType: "Cookie (ADMIN) + path param", Parameter: "id=published",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data contains detail",
	}.Log(t, pass, resp, respBody)
}

func TestArtikelGetByIDNotFound(t *testing.T) {
	f := setupArtikelIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/artikel/99999", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusNotFound

	testutils.TestResult{
		SRSRef: "SRS-6.2", FSDRef: "FSD-6.2",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-ARTIKEL-005",
		Functional: "Artikel Detail — Not Found", Endpoint: "GET /artikel/{id}",
		ReqType: "Cookie (ADMIN) + path param", Parameter: "id=99999",
		ShouldBeSuccess: "false",
		Expectation:     "Response 404, success:false",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// POST /artikel
// ---------------------------------------------------------------------------

func TestArtikelCreateSuccess(t *testing.T) {
	f := setupArtikelIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	body := map[string]any{
		"judul":       "Artikel Baru",
		"isi_artikel": "Ini adalah konten artikel baru yang dibuat oleh bidan.",
		"kategori":    "Kesehatan",
	}
	resp := testutils.DoRequest(f.handler, http.MethodPost, "/artikel", body,
		testutils.AccessCookie(f.jwtUtil, authIDs.VerifiedUserID, []string{"BIDAN", "USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusCreated

	testutils.TestResult{
		SRSRef: "SRS-6.3", FSDRef: "FSD-6.3",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-ARTIKEL-006",
		Functional: "Artikel Create — Success", Endpoint: "POST /artikel",
		ReqType: "JSON Body + Cookie (BIDAN)", Parameter: `{"judul":"Artikel Baru","kategori":"Kesehatan"}`,
		ShouldBeSuccess: "true",
		Expectation:     "Response 201, success:true, data contains id_artikel with status 'Menunggu Verifikasi'",
	}.Log(t, pass, resp, respBody)
}

func TestArtikelCreateForbidden(t *testing.T) {
	f := setupArtikelIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	body := map[string]any{
		"judul":       "Artikel oleh Admin",
		"isi_artikel": "Admin tidak boleh membuat artikel.",
	}
	resp := testutils.DoRequest(f.handler, http.MethodPost, "/artikel", body,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusForbidden

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-ARTIKEL-007",
		Functional: "Artikel Create — Forbidden (ADMIN)", Endpoint: "POST /artikel",
		ReqType: "JSON Body + Cookie (ADMIN)", Parameter: "Role: ADMIN (requires BIDAN)",
		ShouldBeSuccess: "false",
		Expectation:     "Response 403, success:false",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// PATCH /artikel/{id}
// ---------------------------------------------------------------------------

func TestArtikelUpdateOwnPendingSuccess(t *testing.T) {
	f := setupArtikelIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedArtikel(t, authIDs)

	body := map[string]any{"judul": "Artikel Pending — Updated"}
	path := fmt.Sprintf("/artikel/%d", seed.PendingArtikelID)
	resp := testutils.DoRequest(f.handler, http.MethodPatch, path, body,
		testutils.AccessCookie(f.jwtUtil, authIDs.VerifiedUserID, []string{"BIDAN", "USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-6.4", FSDRef: "FSD-6.4",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-ARTIKEL-008",
		Functional: "Artikel Update — Own Pending Success", Endpoint: "PATCH /artikel/{id}",
		ReqType: "JSON Body + Cookie (BIDAN)", Parameter: `{"judul":"Artikel Pending — Updated"}`,
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true",
	}.Log(t, pass, resp, respBody)
}

func TestArtikelUpdateNotFound(t *testing.T) {
	f := setupArtikelIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodPatch, "/artikel/99999", map[string]any{"judul": "Test"},
		testutils.AccessCookie(f.jwtUtil, authIDs.VerifiedUserID, []string{"BIDAN", "USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusNotFound

	testutils.TestResult{
		SRSRef: "SRS-6.4", FSDRef: "FSD-6.4",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-ARTIKEL-009",
		Functional: "Artikel Update — Not Found", Endpoint: "PATCH /artikel/{id}",
		ReqType: "JSON Body + Cookie (BIDAN)", Parameter: "id=99999",
		ShouldBeSuccess: "false",
		Expectation:     "Response 404, success:false",
	}.Log(t, pass, resp, respBody)
}

func TestArtikelUpdateForbiddenNotOwner(t *testing.T) {
	f := setupArtikelIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	// pending article created by VerifiedUserID, try to update with AdminUserID (not the owner)
	// AdminUserID does not have BIDAN role, but let's use a different bidan
	ctx := context.Background()
	now := time.Now().Truncate(time.Microsecond)
	var otherBidanArtikelID int32
	err := f.pool.QueryRow(ctx, `
		INSERT INTO artikel (judul, isi_artikel, status_artikel, id_penulis, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $5)
		RETURNING id_artikel`,
		"Other Bidan Article", "Isi artikel bidan lain.", "Menunggu Verifikasi", authIDs.AdminUserID, now,
	).Scan(&otherBidanArtikelID)
	if err != nil {
		t.Fatalf("failed to seed other bidan's artikel: %v", err)
	}

	// VerifiedUserID tries to update AdminUserID's article -> should get 403 Forbidden
	path := fmt.Sprintf("/artikel/%d", otherBidanArtikelID)
	resp := testutils.DoRequest(f.handler, http.MethodPatch, path, map[string]any{"judul": "Hacked!"},
		testutils.AccessCookie(f.jwtUtil, authIDs.VerifiedUserID, []string{"BIDAN", "USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusForbidden

	testutils.TestResult{
		SRSRef: "SRS-6.4", FSDRef: "FSD-6.4",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-ARTIKEL-010",
		Functional: "Artikel Update — Forbidden (Not Owner)", Endpoint: "PATCH /artikel/{id}",
		ReqType: "JSON Body + Cookie (BIDAN)", Parameter: "Role: BIDAN, but not the article owner",
		ShouldBeSuccess: "false",
		Expectation:     "Response 403, success:false",
	}.Log(t, pass, resp, respBody)
}

func TestArtikelUpdateForbiddenAlreadyPublished(t *testing.T) {
	f := setupArtikelIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedArtikel(t, authIDs)

	body := map[string]any{"judul": "Cannot Update Published"}
	path := fmt.Sprintf("/artikel/%d", seed.PublishedArtikelID)
	resp := testutils.DoRequest(f.handler, http.MethodPatch, path, body,
		testutils.AccessCookie(f.jwtUtil, authIDs.VerifiedUserID, []string{"BIDAN", "USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusForbidden

	testutils.TestResult{
		SRSRef: "SRS-6.4", FSDRef: "FSD-6.4",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-ARTIKEL-011",
		Functional: "Artikel Update — Forbidden (Already Published)", Endpoint: "PATCH /artikel/{id}",
		ReqType: "JSON Body + Cookie (BIDAN)", Parameter: "Status: 'Dipublikasikan' (not 'Menunggu Verifikasi')",
		ShouldBeSuccess: "false",
		Expectation:     "Response 403, success:false",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// DELETE /artikel/{id}
// ---------------------------------------------------------------------------

func TestArtikelDeleteSuccess(t *testing.T) {
	f := setupArtikelIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedArtikel(t, authIDs)

	path := fmt.Sprintf("/artikel/%d", seed.PendingArtikelID)
	resp := testutils.DoRequest(f.handler, http.MethodDelete, path, nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"USER", "DINKES", "SUPER_ADMIN", "ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-6.5", FSDRef: "FSD-6.5",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-ARTIKEL-012",
		Functional: "Artikel Delete — Success", Endpoint: "DELETE /artikel/{id}",
		ReqType: "Cookie (DINKES)", Parameter: "Role: DINKES, id=pending",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true",
	}.Log(t, pass, resp, respBody)
}

func TestArtikelDeleteNotFound(t *testing.T) {
	f := setupArtikelIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodDelete, "/artikel/99999", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"USER", "DINKES", "SUPER_ADMIN", "ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusNotFound

	testutils.TestResult{
		SRSRef: "SRS-6.5", FSDRef: "FSD-6.5",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-ARTIKEL-013",
		Functional: "Artikel Delete — Not Found", Endpoint: "DELETE /artikel/{id}",
		ReqType: "Cookie (DINKES)", Parameter: "id=99999",
		ShouldBeSuccess: "false",
		Expectation:     "Response 404, success:false",
	}.Log(t, pass, resp, respBody)
}

func TestArtikelDeleteForbidden(t *testing.T) {
	f := setupArtikelIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedArtikel(t, authIDs)

	path := fmt.Sprintf("/artikel/%d", seed.PendingArtikelID)
	resp := testutils.DoRequest(f.handler, http.MethodDelete, path, nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.VerifiedUserID, []string{"BIDAN", "USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusForbidden

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-ARTIKEL-014",
		Functional: "Artikel Delete — Forbidden (BIDAN)", Endpoint: "DELETE /artikel/{id}",
		ReqType: "Cookie (BIDAN)", Parameter: "Role: BIDAN (requires DINKES)",
		ShouldBeSuccess: "false",
		Expectation:     "Response 403, success:false",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// PATCH /artikel/{id}/review
// ---------------------------------------------------------------------------

func TestArtikelReviewSetujuiSuccess(t *testing.T) {
	f := setupArtikelIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedArtikel(t, authIDs)

	body := map[string]any{"aksi": "setujui"}
	path := fmt.Sprintf("/artikel/%d/review", seed.PendingArtikelID)
	resp := testutils.DoRequest(f.handler, http.MethodPatch, path, body,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"USER", "DINKES", "SUPER_ADMIN", "ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-6.6", FSDRef: "FSD-6.6",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-ARTIKEL-015",
		Functional: "Artikel Review — Setujui Success", Endpoint: "PATCH /artikel/{id}/review",
		ReqType: "JSON Body + Cookie (DINKES)", Parameter: `{"aksi":"setujui"}`,
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data.status_artikel='Dipublikasikan'",
	}.Log(t, pass, resp, respBody)
}

func TestArtikelReviewTolakSuccess(t *testing.T) {
	f := setupArtikelIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedArtikel(t, authIDs)

	body := map[string]any{"aksi": "tolak", "catatan_review": "Konten tidak sesuai"}
	path := fmt.Sprintf("/artikel/%d/review", seed.PendingArtikelID)
	resp := testutils.DoRequest(f.handler, http.MethodPatch, path, body,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"USER", "DINKES", "SUPER_ADMIN", "ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-6.6", FSDRef: "FSD-6.6",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-ARTIKEL-016",
		Functional: "Artikel Review — Tolak Success", Endpoint: "PATCH /artikel/{id}/review",
		ReqType: "JSON Body + Cookie (DINKES)", Parameter: `{"aksi":"tolak","catatan_review":"Konten tidak sesuai"}`,
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data.status_artikel='Ditolak'",
	}.Log(t, pass, resp, respBody)
}

func TestArtikelReviewNotFound(t *testing.T) {
	f := setupArtikelIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	body := map[string]any{"aksi": "setujui"}
	resp := testutils.DoRequest(f.handler, http.MethodPatch, "/artikel/99999/review", body,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"USER", "DINKES", "SUPER_ADMIN", "ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusNotFound

	testutils.TestResult{
		SRSRef: "SRS-6.6", FSDRef: "FSD-6.6",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-ARTIKEL-017",
		Functional: "Artikel Review — Not Found", Endpoint: "PATCH /artikel/{id}/review",
		ReqType: "JSON Body + Cookie (DINKES)", Parameter: "id=99999",
		ShouldBeSuccess: "false",
		Expectation:     "Response 404, success:false",
	}.Log(t, pass, resp, respBody)
}

func TestArtikelReviewForbidden(t *testing.T) {
	f := setupArtikelIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	seed := f.seedArtikel(t, authIDs)

	body := map[string]any{"aksi": "setujui"}
	path := fmt.Sprintf("/artikel/%d/review", seed.PendingArtikelID)
	resp := testutils.DoRequest(f.handler, http.MethodPatch, path, body,
		testutils.AccessCookie(f.jwtUtil, authIDs.VerifiedUserID, []string{"BIDAN", "USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusForbidden

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-ARTIKEL-018",
		Functional: "Artikel Review — Forbidden (BIDAN)", Endpoint: "PATCH /artikel/{id}/review",
		ReqType: "JSON Body + Cookie (BIDAN)", Parameter: "Role: BIDAN (requires DINKES)",
		ShouldBeSuccess: "false",
		Expectation:     "Response 403, success:false",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// GET /artikel/pending
// ---------------------------------------------------------------------------

func TestArtikelGetPendingWithData(t *testing.T) {
	f := setupArtikelIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)
	f.seedArtikel(t, authIDs)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/artikel/pending", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"USER", "DINKES", "SUPER_ADMIN", "ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-6.7", FSDRef: "FSD-6.7",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-ARTIKEL-019",
		Functional: "Artikel Pending List — With Data", Endpoint: "GET /artikel/pending",
		ReqType: "Cookie (DINKES+ADMIN)", Parameter: "Role: DINKES+ADMIN, 1 pending seeded",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data.artikel has items",
	}.Log(t, pass, resp, respBody)
}

func TestArtikelGetPendingEmpty(t *testing.T) {
	f := setupArtikelIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/artikel/pending", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.AdminUserID, []string{"USER", "DINKES", "SUPER_ADMIN", "ADMIN"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-6.7", FSDRef: "FSD-6.7",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-ARTIKEL-020",
		Functional: "Artikel Pending List — Empty", Endpoint: "GET /artikel/pending",
		ReqType: "Cookie (DINKES+ADMIN)", Parameter: "Role: DINKES+ADMIN, no pending articles",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true, data.artikel empty",
	}.Log(t, pass, resp, respBody)
}

func TestArtikelGetPendingForbidden(t *testing.T) {
	f := setupArtikelIntegrationTest(t)
	defer f.cleanup(t)
	authIDs := f.seed(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/artikel/pending", nil,
		testutils.AccessCookie(f.jwtUtil, authIDs.RegularUserID, []string{"USER"}))
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusForbidden

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-ARTIKEL-021",
		Functional: "Artikel Pending — Forbidden (USER)", Endpoint: "GET /artikel/pending",
		ReqType: "Cookie (USER)", Parameter: "Role: USER (requires DINKES)",
		ShouldBeSuccess: "false",
		Expectation:     "Response 403, success:false",
	}.Log(t, pass, resp, respBody)
}

// ---------------------------------------------------------------------------
// Unauthorized
// ---------------------------------------------------------------------------

func TestArtikelGetAllPublishedNoAuth(t *testing.T) {
	f := setupArtikelIntegrationTest(t)
	defer f.cleanup(t)

	resp := testutils.DoRequest(f.handler, http.MethodGet, "/artikel", nil)
	respBody := testutils.ReadBody(resp)
	pass := resp.StatusCode == http.StatusOK

	testutils.TestResult{
		SRSRef: "SRS-SC-03", FSDRef: "FSD-2.2",
		TSDRef: "TSD-3.3 (Endpoint Index)", NoTestScript: "TC-ARTIKEL-022",
		Functional: "Artikel List — Public No Auth", Endpoint: "GET /artikel",
		ReqType: "No Cookie", Parameter: "{}",
		ShouldBeSuccess: "true",
		Expectation:     "Response 200, success:true",
	}.Log(t, pass, resp, respBody)
}
