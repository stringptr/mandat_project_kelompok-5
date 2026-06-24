package test

import (
	"fmt"
	"net/http"
	"testing"
)

// TestImunisasi tests all Imunisasi endpoints
func TestImunisasi(t *testing.T) {
	cfg := SetupTestConfig(t)
	defer TeardownTestConfig(cfg)

	tc := NewTestContext(t, cfg)
	tc.SetupTestData()
	defer tc.CleanupTestData()

	t.Run("Happy Path", func(t *testing.T) {
		testImunisasiHappyPath(tc)
	})

	t.Run("Sad Path", func(t *testing.T) {
		testImunisasiSadPath(tc)
	})
}

func testImunisasiHappyPath(tc *TestContext) {
	bidanToken := tc.GetTestToken("BIDAN")
	t := tc.T

	var imunisasiID int

	// TS-BE-014: Daftar Jadwal Imunisasi
	t.Run("TS-BE-014: Daftar Jadwal Imunisasi", func(t *testing.T) {
		resp := tc.MakeRequest("GET", "/api/v1/imunisasi?page=1&per_page=15", nil, bidanToken)
		tc.AssertStatusCode(resp, http.StatusOK, "Should return 200 OK")

		result := tc.ParseResponse(resp)
		tc.AssertSuccess(result, "Response should be successful")
	})

	// TS-BE-015: Tambah Jadwal Imunisasi
	t.Run("TS-BE-015: Tambah Jadwal Imunisasi", func(t *testing.T) {
		body := map[string]interface{}{
			"id_anak":        1,
			"jenis_vaksin":   "BCG",
			"tanggal_jadwal": "2026-07-01",
			"lokasi":         "Posyandu RW 05",
		}
		resp := tc.MakeRequest("POST", "/api/v1/imunisasi", body, bidanToken)
		tc.AssertStatusCode(resp, http.StatusCreated, "Should return 201 Created")

		result := tc.ParseResponse(resp)
		tc.AssertSuccess(result, "Response should be successful")

		// Store ID for subsequent tests
		if data, ok := result.Data.(map[string]interface{}); ok {
			if id, ok := data["id"].(float64); ok {
				imunisasiID = int(id)
			}
		}
	})

	// TS-BE-016: Update Realisasi Vaksin
	t.Run("TS-BE-016: Update Realisasi Vaksin", func(t *testing.T) {
		if imunisasiID > 0 {
			body := map[string]interface{}{
				"tanggal_realisasi": "2026-07-01",
				"catatan":           "Vaksinasi berhasil dilakukan",
				"petugas":           "Bidan Test",
			}
			resp := tc.MakeRequest("PATCH", fmt.Sprintf("/api/v1/imunisasi/%d/realisasi", imunisasiID), body, bidanToken)
			tc.AssertStatusCode(resp, http.StatusOK, "Should return 200 OK")

			result := tc.ParseResponse(resp)
			tc.AssertSuccess(result, "Response should be successful")
		}
	})
}

func testImunisasiSadPath(tc *TestContext) {
	t := tc.T
	bidanToken := tc.GetTestToken("BIDAN")

	// TS-BE-NEG-017: Daftar Jadwal Imunisasi Tanpa Auth
	t.Run("TS-BE-NEG-017: Daftar Jadwal Imunisasi Tanpa Auth", func(t *testing.T) {
		resp := tc.MakeRequest("GET", "/api/v1/imunisasi", nil, "")
		tc.AssertStatusCode(resp, http.StatusUnauthorized, "Should return 401 Unauthorized")
	})

	// TS-BE-NEG-018: Tambah Jadwal Imunisasi Data Invalid
	t.Run("TS-BE-NEG-018: Tambah Jadwal Imunisasi Data Invalid", func(t *testing.T) {
		body := map[string]interface{}{
			"jenis_vaksin": "", // Empty vaccine type
		}
		resp := tc.MakeRequest("POST", "/api/v1/imunisasi", body, bidanToken)
		tc.AssertStatusCode(resp, http.StatusUnprocessableEntity, "Should return 422 Validation Error")
	})

	// TS-BE-NEG-019: Update Realisasi ID Tidak Ditemukan
	t.Run("TS-BE-NEG-019: Update Realisasi ID Tidak Ditemukan", func(t *testing.T) {
		body := map[string]interface{}{
			"tanggal_realisasi": "2026-07-01",
			"catatan":           "Test",
		}
		resp := tc.MakeRequest("PATCH", "/api/v1/imunisasi/99999/realisasi", body, bidanToken)
		tc.AssertStatusCode(resp, http.StatusNotFound, "Should return 404 Not Found")
	})
}
