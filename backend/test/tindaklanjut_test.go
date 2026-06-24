package test

import (
	"fmt"
	"net/http"
	"testing"
)

// TestTindakLanjut tests all Tindak Lanjut endpoints
func TestTindakLanjut(t *testing.T) {
	cfg := SetupTestConfig(t)
	defer TeardownTestConfig(cfg)

	tc := NewTestContext(t, cfg)
	tc.SetupTestData()
	defer tc.CleanupTestData()

	t.Run("Happy Path", func(t *testing.T) {
		testTindakLanjutHappyPath(tc)
	})

	t.Run("Sad Path", func(t *testing.T) {
		testTindakLanjutSadPath(tc)
	})
}

func testTindakLanjutHappyPath(tc *TestContext) {
	bidanToken := tc.GetTestToken("BIDAN")
	kaderToken := tc.GetTestToken("KADER")
	dinkesToken := tc.GetTestToken("DINKES")
	t := tc.T

	var tindakLanjutID int
	var rujukanID int

	// TS-BE-017: Daftar Pasien Tindak Lanjut
	t.Run("TS-BE-017: Daftar Pasien Tindak Lanjut", func(t *testing.T) {
		resp := tc.MakeRequest("GET", "/api/v1/tindak-lanjut/pasien?page=1&per_page=15", nil, bidanToken)
		tc.AssertStatusCode(resp, http.StatusOK, "Should return 200 OK")

		result := tc.ParseResponse(resp)
		tc.AssertSuccess(result, "Response should be successful")
	})

	// TS-BE-018: Detail Pasien Tindak Lanjut
	t.Run("TS-BE-018: Detail Pasien Tindak Lanjut", func(t *testing.T) {
		resp := tc.MakeRequest("GET", "/api/v1/tindak-lanjut/pasien/1", nil, bidanToken)
		tc.AssertStatusCode(resp, http.StatusOK, "Should return 200 OK")

		result := tc.ParseResponse(resp)
		tc.AssertSuccess(result, "Response should be successful")
	})

	// TS-BE-019: Membuat Tindak Lanjut
	t.Run("TS-BE-019: Membuat Tindak Lanjut", func(t *testing.T) {
		body := map[string]interface{}{
			"id_hasil_pemeriksaan": 1,
			"jenis_tindak_lanjut":  "KONSELING",
			"catatan":              "Perlu konseling gizi untuk ibu",
			"tanggal_tindak_lanjut": "2026-06-15",
		}
		resp := tc.MakeRequest("POST", "/api/v1/tindak-lanjut", body, bidanToken)
		tc.AssertStatusCode(resp, http.StatusCreated, "Should return 201 Created")

		result := tc.ParseResponse(resp)
		tc.AssertSuccess(result, "Response should be successful")

		// Store ID for subsequent tests
		if data, ok := result.Data.(map[string]interface{}); ok {
			if id, ok := data["id"].(float64); ok {
				tindakLanjutID = int(id)
			}
		}
	})

	// TS-BE-021: Detail Tindak Lanjut
	t.Run("TS-BE-021: Detail Tindak Lanjut", func(t *testing.T) {
		if tindakLanjutID > 0 {
			resp := tc.MakeRequest("GET", fmt.Sprintf("/api/v1/tindak-lanjut/%d", tindakLanjutID), nil, bidanToken)
			tc.AssertStatusCode(resp, http.StatusOK, "Should return 200 OK")

			result := tc.ParseResponse(resp)
			tc.AssertSuccess(result, "Response should be successful")
		}
	})

	// TS-BE-020: Update Status Rujukan
	t.Run("TS-BE-020: Update Status Rujukan", func(t *testing.T) {
		// First create a rujukan
		body := map[string]interface{}{
			"id_hasil_pemeriksaan": 1,
			"jenis_tindak_lanjut":  "RUJUKAN",
			"tujuan_rujukan":       "Puskesmas Kecamatan",
			"alasan_rujukan":       "Gizi buruk memerlukan penanganan lanjut",
		}
		resp := tc.MakeRequest("POST", "/api/v1/tindak-lanjut", body, bidanToken)
		
		if resp.StatusCode == http.StatusCreated {
			result := tc.ParseResponse(resp)
			if data, ok := result.Data.(map[string]interface{}); ok {
				if id, ok := data["id"].(float64); ok {
					rujukanID = int(id)
				}
			}
		}

		// Now update the rujukan status
		if rujukanID > 0 {
			updateBody := map[string]interface{}{
				"status_rujukan": "SELESAI",
			}
			updateResp := tc.MakeRequest("PATCH", fmt.Sprintf("/api/v1/rujukan/%d/status", rujukanID), updateBody, bidanToken)
			tc.AssertStatusCode(updateResp, http.StatusOK, "Should return 200 OK")

			result := tc.ParseResponse(updateResp)
			tc.AssertSuccess(result, "Response should be successful")
		}
	})

	// TS-BE-022: Status Tindak Lanjut (Kader)
	t.Run("TS-BE-022: Status Tindak Lanjut (Kader)", func(t *testing.T) {
		resp := tc.MakeRequest("GET", "/api/v1/tindak-lanjut/status?page=1&per_page=15", nil, kaderToken)
		tc.AssertStatusCode(resp, http.StatusOK, "Should return 200 OK")

		result := tc.ParseResponse(resp)
		tc.AssertSuccess(result, "Response should be successful")
	})

	// TS-BE-023: Laporan Tindak Lanjut
	t.Run("TS-BE-023: Laporan Tindak Lanjut", func(t *testing.T) {
		resp := tc.MakeRequest("GET", "/api/v1/laporan/tindak-lanjut?periode_awal=2026-01-01&periode_akhir=2026-12-31", nil, dinkesToken)
		tc.AssertStatusCode(resp, http.StatusOK, "Should return 200 OK")

		result := tc.ParseResponse(resp)
		tc.AssertSuccess(result, "Response should be successful")
	})
}

func testTindakLanjutSadPath(tc *TestContext) {
	t := tc.T
	bidanToken := tc.GetTestToken("BIDAN")
	kaderToken := tc.GetTestToken("KADER")
	dinkesToken := tc.GetTestToken("DINKES")

	// TS-BE-NEG-020: Daftar Pasien Tindak Lanjut Bukan Bidan
	t.Run("TS-BE-NEG-020: Daftar Pasien Tindak Lanjut Bukan Bidan", func(t *testing.T) {
		resp := tc.MakeRequest("GET", "/api/v1/tindak-lanjut/pasien", nil, kaderToken)
		tc.AssertStatusCode(resp, http.StatusForbidden, "Should return 403 Forbidden")
	})

	// TS-BE-NEG-021: Detail Pasien Tindak Lanjut ID Invalid
	t.Run("TS-BE-NEG-021: Detail Pasien Tindak Lanjut ID Invalid", func(t *testing.T) {
		resp := tc.MakeRequest("GET", "/api/v1/tindak-lanjut/pasien/99999", nil, bidanToken)
		tc.AssertStatusCode(resp, http.StatusNotFound, "Should return 404 Not Found")
	})

	// TS-BE-NEG-022: Membuat Tindak Lanjut Data Tidak Lengkap
	t.Run("TS-BE-NEG-022: Membuat Tindak Lanjut Data Tidak Lengkap", func(t *testing.T) {
		body := map[string]interface{}{
			"id_hasil_pemeriksaan": 1,
			// Missing jenis_tindak_lanjut
		}
		resp := tc.MakeRequest("POST", "/api/v1/tindak-lanjut", body, bidanToken)
		tc.AssertStatusCode(resp, http.StatusUnprocessableEntity, "Should return 422 Validation Error")
	})

	// TS-BE-NEG-023: Membuat Tindak Lanjut Bukan Bidan
	t.Run("TS-BE-NEG-023: Membuat Tindak Lanjut Bukan Bidan", func(t *testing.T) {
		body := map[string]interface{}{
			"id_hasil_pemeriksaan": 1,
			"jenis_tindak_lanjut":  "KONSELING",
			"catatan":              "Test",
		}
		resp := tc.MakeRequest("POST", "/api/v1/tindak-lanjut", body, kaderToken)
		tc.AssertStatusCode(resp, http.StatusForbidden, "Should return 403 Forbidden")
	})

	// TS-BE-NEG-024: Update Status Rujukan Status Invalid
	t.Run("TS-BE-NEG-024: Update Status Rujukan Status Invalid", func(t *testing.T) {
		body := map[string]interface{}{
			"status_rujukan": "INVALID_STATUS",
		}
		resp := tc.MakeRequest("PATCH", "/api/v1/rujukan/1/status", body, bidanToken)
		tc.AssertStatusCode(resp, http.StatusUnprocessableEntity, "Should return 422 Validation Error")
	})

	// TS-BE-NEG-025: Detail Tindak Lanjut Akses Ditolak
	t.Run("TS-BE-NEG-025: Detail Tindak Lanjut Akses Ditolak", func(t *testing.T) {
		// Assuming user doesn't own this tindak lanjut
		resp := tc.MakeRequest("GET", "/api/v1/tindak-lanjut/1", nil, kaderToken)
		tc.AssertStatusCode(resp, http.StatusForbidden, "Should return 403 Forbidden")
	})

	// TS-BE-NEG-026: Status Tindak Lanjut Bukan Kader
	t.Run("TS-BE-NEG-026: Status Tindak Lanjut Bukan Kader", func(t *testing.T) {
		resp := tc.MakeRequest("GET", "/api/v1/tindak-lanjut/status", nil, bidanToken)
		tc.AssertStatusCode(resp, http.StatusForbidden, "Should return 403 Forbidden")
	})

	// TS-BE-NEG-027: Laporan Tindak Lanjut Bukan Dinkes
	t.Run("TS-BE-NEG-027: Laporan Tindak Lanjut Bukan Dinkes", func(t *testing.T) {
		resp := tc.MakeRequest("GET", "/api/v1/laporan/tindak-lanjut?periode_awal=2026-01-01&periode_akhir=2026-12-31", nil, bidanToken)
		tc.AssertStatusCode(resp, http.StatusForbidden, "Should return 403 Forbidden")
	})
}
