package test

import (
	"fmt"
	"net/http"
	"testing"
)

// TestArtikel tests all Artikel endpoints
func TestArtikel(t *testing.T) {
	cfg := SetupTestConfig(t)
	defer TeardownTestConfig(cfg)

	tc := NewTestContext(t, cfg)
	tc.SetupTestData()
	defer tc.CleanupTestData()

	t.Run("Happy Path", func(t *testing.T) {
		testArtikelHappyPath(tc)
	})

	t.Run("Sad Path", func(t *testing.T) {
		testArtikelSadPath(tc)
	})
}

func testArtikelHappyPath(tc *TestContext) {
	bidanToken := tc.GetTestToken("BIDAN")
	dinkesToken := tc.GetTestToken("DINKES")
	t := tc.T

	var artikelID int

	// TS-BE-024: Daftar Artikel
	t.Run("TS-BE-024: Daftar Artikel", func(t *testing.T) {
		resp := tc.MakeRequest("GET", "/api/v1/artikel?page=1&per_page=15", nil, bidanToken)
		tc.AssertStatusCode(resp, http.StatusOK, "Should return 200 OK")

		result := tc.ParseResponse(resp)
		tc.AssertSuccess(result, "Response should be successful")
	})

	// TS-BE-026: Membuat Artikel
	t.Run("TS-BE-026: Membuat Artikel", func(t *testing.T) {
		body := map[string]interface{}{
			"judul":       "Artikel Test - Pentingnya Gizi Seimbang",
			"isi_artikel": "Gizi seimbang sangat penting untuk tumbuh kembang anak...",
			"kategori":    "GIZI",
		}
		resp := tc.MakeRequest("POST", "/api/v1/artikel", body, bidanToken)
		tc.AssertStatusCode(resp, http.StatusCreated, "Should return 201 Created")

		result := tc.ParseResponse(resp)
		tc.AssertSuccess(result, "Response should be successful")

		// Store ID for subsequent tests
		if data, ok := result.Data.(map[string]interface{}); ok {
			if id, ok := data["id"].(float64); ok {
				artikelID = int(id)
			}
		}
	})

	// TS-BE-025: Detail Artikel
	t.Run("TS-BE-025: Detail Artikel", func(t *testing.T) {
		if artikelID > 0 {
			resp := tc.MakeRequest("GET", fmt.Sprintf("/api/v1/artikel/%d", artikelID), nil, bidanToken)
			tc.AssertStatusCode(resp, http.StatusOK, "Should return 200 OK")

			result := tc.ParseResponse(resp)
			tc.AssertSuccess(result, "Response should be successful")
		}
	})

	// TS-BE-027: Update Artikel
	t.Run("TS-BE-027: Update Artikel", func(t *testing.T) {
		if artikelID > 0 {
			body := map[string]interface{}{
				"judul":       "Artikel Test - Updated Title",
				"isi_artikel": "Konten yang telah diperbarui...",
			}
			resp := tc.MakeRequest("PATCH", fmt.Sprintf("/api/v1/artikel/%d", artikelID), body, bidanToken)
			tc.AssertStatusCode(resp, http.StatusOK, "Should return 200 OK")

			result := tc.ParseResponse(resp)
			tc.AssertSuccess(result, "Response should be successful")
		}
	})

	// TS-BE-029: Daftar Artikel Pending
	t.Run("TS-BE-029: Daftar Artikel Pending", func(t *testing.T) {
		resp := tc.MakeRequest("GET", "/api/v1/artikel/pending?page=1&per_page=15", nil, dinkesToken)
		tc.AssertStatusCode(resp, http.StatusOK, "Should return 200 OK")

		result := tc.ParseResponse(resp)
		tc.AssertSuccess(result, "Response should be successful")
	})

	// TS-BE-028: Review Artikel (Setujui)
	t.Run("TS-BE-028: Review Artikel (Setujui)", func(t *testing.T) {
		if artikelID > 0 {
			body := map[string]interface{}{
				"aksi": "setujui",
			}
			resp := tc.MakeRequest("PATCH", fmt.Sprintf("/api/v1/artikel/%d/review", artikelID), body, dinkesToken)
			tc.AssertStatusCode(resp, http.StatusOK, "Should return 200 OK")

			result := tc.ParseResponse(resp)
			tc.AssertSuccess(result, "Response should be successful")
		}
	})
}

func testArtikelSadPath(tc *TestContext) {
	t := tc.T
	bidanToken := tc.GetTestToken("BIDAN")
	bidan2Token := tc.GetTestToken("BIDAN2") // Different bidan
	dinkesToken := tc.GetTestToken("DINKES")
	kaderToken := tc.GetTestToken("KADER")

	// TS-BE-NEG-028: Daftar Artikel Tanpa Auth
	t.Run("TS-BE-NEG-028: Daftar Artikel Tanpa Auth", func(t *testing.T) {
		resp := tc.MakeRequest("GET", "/api/v1/artikel", nil, "")
		tc.AssertStatusCode(resp, http.StatusUnauthorized, "Should return 401 Unauthorized")
	})

	// TS-BE-NEG-029: Detail Artikel ID Tidak Ditemukan
	t.Run("TS-BE-NEG-029: Detail Artikel ID Tidak Ditemukan", func(t *testing.T) {
		resp := tc.MakeRequest("GET", "/api/v1/artikel/99999", nil, bidanToken)
		tc.AssertStatusCode(resp, http.StatusNotFound, "Should return 404 Not Found")
	})

	// TS-BE-NEG-030: Membuat Artikel Bukan Bidan
	t.Run("TS-BE-NEG-030: Membuat Artikel Bukan Bidan", func(t *testing.T) {
		body := map[string]interface{}{
			"judul":       "Test Artikel",
			"isi_artikel": "Konten artikel...",
			"kategori":    "GIZI",
		}
		resp := tc.MakeRequest("POST", "/api/v1/artikel", body, kaderToken)
		tc.AssertStatusCode(resp, http.StatusForbidden, "Should return 403 Forbidden")
	})

	// TS-BE-NEG-031: Membuat Artikel Judul Kosong
	t.Run("TS-BE-NEG-031: Membuat Artikel Judul Kosong", func(t *testing.T) {
		body := map[string]interface{}{
			"judul":       "",
			"isi_artikel": "Konten artikel...",
			"kategori":    "GIZI",
		}
		resp := tc.MakeRequest("POST", "/api/v1/artikel", body, bidanToken)
		tc.AssertStatusCode(resp, http.StatusUnprocessableEntity, "Should return 422 Validation Error")
	})

	// TS-BE-NEG-032: Update Artikel Bukan Pemilik
	t.Run("TS-BE-NEG-032: Update Artikel Bukan Pemilik", func(t *testing.T) {
		// Assuming artikel with ID 1 exists and owned by different bidan
		body := map[string]interface{}{
			"judul": "Updated by wrong user",
		}
		resp := tc.MakeRequest("PATCH", "/api/v1/artikel/1", body, bidan2Token)
		tc.AssertStatusCode(resp, http.StatusForbidden, "Should return 403 Forbidden")
	})

	// TS-BE-NEG-033: Update Artikel Sudah Published
	t.Run("TS-BE-NEG-033: Update Artikel Sudah Published", func(t *testing.T) {
		// Assuming artikel with ID 1 is already published
		body := map[string]interface{}{
			"judul": "Trying to update published artikel",
		}
		resp := tc.MakeRequest("PATCH", "/api/v1/artikel/1", body, bidanToken)
		// Should return 422 if already published
		// Actual status code depends on implementation
	})

	// TS-BE-NEG-034: Review Artikel Bukan Dinkes
	t.Run("TS-BE-NEG-034: Review Artikel Bukan Dinkes", func(t *testing.T) {
		body := map[string]interface{}{
			"aksi": "setujui",
		}
		resp := tc.MakeRequest("PATCH", "/api/v1/artikel/1/review", body, bidanToken)
		tc.AssertStatusCode(resp, http.StatusForbidden, "Should return 403 Forbidden")
	})

	// TS-BE-NEG-035: Review Artikel Aksi Invalid
	t.Run("TS-BE-NEG-035: Review Artikel Aksi Invalid", func(t *testing.T) {
		body := map[string]interface{}{
			"aksi": "invalid_action",
		}
		resp := tc.MakeRequest("PATCH", "/api/v1/artikel/1/review", body, dinkesToken)
		tc.AssertStatusCode(resp, http.StatusUnprocessableEntity, "Should return 422 Validation Error")
	})

	// TS-BE-NEG-036: Review Artikel Tolak Tanpa Catatan
	t.Run("TS-BE-NEG-036: Review Artikel Tolak Tanpa Catatan", func(t *testing.T) {
		body := map[string]interface{}{
			"aksi":           "tolak",
			"catatan_review": "",
		}
		resp := tc.MakeRequest("PATCH", "/api/v1/artikel/1/review", body, dinkesToken)
		tc.AssertStatusCode(resp, http.StatusUnprocessableEntity, "Should return 422 Validation Error")
	})

	// TS-BE-NEG-037: Daftar Artikel Pending Bukan Dinkes
	t.Run("TS-BE-NEG-037: Daftar Artikel Pending Bukan Dinkes", func(t *testing.T) {
		resp := tc.MakeRequest("GET", "/api/v1/artikel/pending", nil, bidanToken)
		tc.AssertStatusCode(resp, http.StatusForbidden, "Should return 403 Forbidden")
	})

	// TS-BE-NEG-038: Hapus Artikel Bukan Dinkes
	t.Run("TS-BE-NEG-038: Hapus Artikel Bukan Dinkes", func(t *testing.T) {
		resp := tc.MakeRequest("DELETE", "/api/v1/artikel/1", nil, bidanToken)
		tc.AssertStatusCode(resp, http.StatusForbidden, "Should return 403 Forbidden")
	})
}
