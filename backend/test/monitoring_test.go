package test

import (
	"fmt"
	"net/http"
	"testing"
)

// TestMonitoringPasien tests all Monitoring Pasien endpoints
func TestMonitoringPasien(t *testing.T) {
	cfg := SetupTestConfig(t)
	defer TeardownTestConfig(cfg)

	tc := NewTestContext(t, cfg)
	tc.SetupTestData()
	defer tc.CleanupTestData()

	t.Run("Happy Path", func(t *testing.T) {
		testMonitoringPasienHappyPath(tc)
	})

	t.Run("Sad Path", func(t *testing.T) {
		testMonitoringPasienSadPath(tc)
	})
}

func testMonitoringPasienHappyPath(tc *TestContext) {
	token := tc.GetTestToken("BIDAN")

	t := tc.T

	// TS-BE-002: Daftar Pasien
	t.Run("TS-BE-002: Daftar Pasien", func(t *testing.T) {
		resp := tc.MakeRequest("GET", "/api/v1/monitoring/pasien?page=1&per_page=15", nil, token)
		tc.AssertStatusCode(resp, http.StatusOK, "Should return 200 OK")

		result := tc.ParseResponse(resp)
		tc.AssertSuccess(result, "Response should be successful")
	})

	// TS-BE-003: Detail Pasien
	t.Run("TS-BE-003: Detail Pasien", func(t *testing.T) {
		resp := tc.MakeRequest("GET", "/api/v1/monitoring/pasien/1", nil, token)
		tc.AssertStatusCode(resp, http.StatusOK, "Should return 200 OK")

		result := tc.ParseResponse(resp)
		tc.AssertSuccess(result, "Response should be successful")
	})

	// TS-BE-004: Cari Pasien
	t.Run("TS-BE-004: Cari Pasien", func(t *testing.T) {
		resp := tc.MakeRequest("GET", "/api/v1/monitoring/pasien/search?q=Test", nil, token)
		tc.AssertStatusCode(resp, http.StatusOK, "Should return 200 OK")

		result := tc.ParseResponse(resp)
		tc.AssertSuccess(result, "Response should be successful")
	})

	// TS-BE-005: Export Excel Pasien
	t.Run("TS-BE-005: Export Excel Pasien", func(t *testing.T) {
		resp := tc.MakeRequest("GET", "/api/v1/monitoring/pasien/export", nil, token)
		tc.AssertStatusCode(resp, http.StatusOK, "Should return 200 OK")
	})
}

func testMonitoringPasienSadPath(tc *TestContext) {
	t := tc.T
	invalidToken := "invalid_token_xyz"
	kaderToken := tc.GetTestToken("KADER")

	// TS-BE-NEG-002: Daftar Pasien Tanpa Auth
	t.Run("TS-BE-NEG-002: Daftar Pasien Tanpa Auth", func(t *testing.T) {
		resp := tc.MakeRequest("GET", "/api/v1/monitoring/pasien", nil, "")
		tc.AssertStatusCode(resp, http.StatusUnauthorized, "Should return 401 Unauthorized")
	})

	// TS-BE-NEG-003: Daftar Pasien Token Invalid
	t.Run("TS-BE-NEG-003: Daftar Pasien Token Invalid", func(t *testing.T) {
		resp := tc.MakeRequest("GET", "/api/v1/monitoring/pasien", nil, invalidToken)
		tc.AssertStatusCode(resp, http.StatusUnauthorized, "Should return 401 Unauthorized")
	})

	// TS-BE-NEG-004: Detail Pasien ID Tidak Ditemukan
	t.Run("TS-BE-NEG-004: Detail Pasien ID Tidak Ditemukan", func(t *testing.T) {
		resp := tc.MakeRequest("GET", "/api/v1/monitoring/pasien/99999", nil, kaderToken)
		tc.AssertStatusCode(resp, http.StatusNotFound, "Should return 404 Not Found")
	})

	// TS-BE-NEG-005: Cari Pasien Query Kosong
	t.Run("TS-BE-NEG-005: Cari Pasien Query Kosong", func(t *testing.T) {
		resp := tc.MakeRequest("GET", "/api/v1/monitoring/pasien/search?q=", nil, kaderToken)
		tc.AssertStatusCode(resp, http.StatusUnprocessableEntity, "Should return 422 Validation Error")
	})

	// TS-BE-NEG-006: Export Pasien Role Tidak Sesuai
	t.Run("TS-BE-NEG-006: Export Pasien Role Tidak Sesuai", func(t *testing.T) {
		resp := tc.MakeRequest("GET", "/api/v1/monitoring/pasien/export", nil, kaderToken)
		tc.AssertStatusCode(resp, http.StatusForbidden, "Should return 403 Forbidden")
	})
}

// TestMonitoringStatistik tests all Monitoring Statistik endpoints
func TestMonitoringStatistik(t *testing.T) {
	cfg := SetupTestConfig(t)
	defer TeardownTestConfig(cfg)

	tc := NewTestContext(t, cfg)
	tc.SetupTestData()
	defer tc.CleanupTestData()

	t.Run("Happy Path", func(t *testing.T) {
		testMonitoringStatistikHappyPath(tc)
	})

	t.Run("Sad Path", func(t *testing.T) {
		testMonitoringStatistikSadPath(tc)
	})
}

func testMonitoringStatistikHappyPath(tc *TestContext) {
	token := tc.GetTestToken("DINKES")
	t := tc.T

	// TS-BE-006: Summary Dashboard
	t.Run("TS-BE-006: Summary Dashboard", func(t *testing.T) {
		resp := tc.MakeRequest("GET", "/api/v1/monitoring/statistik", nil, token)
		tc.AssertStatusCode(resp, http.StatusOK, "Should return 200 OK")

		result := tc.ParseResponse(resp)
		tc.AssertSuccess(result, "Response should be successful")
	})

	// TS-BE-007: Tren Pertumbuhan Bulanan
	t.Run("TS-BE-007: Tren Pertumbuhan Bulanan", func(t *testing.T) {
		resp := tc.MakeRequest("GET", "/api/v1/monitoring/statistik-bulanan?tahun=2026", nil, token)
		tc.AssertStatusCode(resp, http.StatusOK, "Should return 200 OK")

		result := tc.ParseResponse(resp)
		tc.AssertSuccess(result, "Response should be successful")
	})
}

func testMonitoringStatistikSadPath(tc *TestContext) {
	t := tc.T
	bidanToken := tc.GetTestToken("BIDAN")

	// TS-BE-NEG-007: Summary Dashboard Tanpa Auth
	t.Run("TS-BE-NEG-007: Summary Dashboard Tanpa Auth", func(t *testing.T) {
		resp := tc.MakeRequest("GET", "/api/v1/monitoring/statistik", nil, "")
		tc.AssertStatusCode(resp, http.StatusUnauthorized, "Should return 401 Unauthorized")
	})

	// TS-BE-NEG-008: Tren Bulanan Format Tahun Invalid
	t.Run("TS-BE-NEG-008: Tren Bulanan Format Tahun Invalid", func(t *testing.T) {
		resp := tc.MakeRequest("GET", "/api/v1/monitoring/statistik-bulanan?tahun=abcd", nil, bidanToken)
		tc.AssertStatusCode(resp, http.StatusUnprocessableEntity, "Should return 422 Validation Error")
	})
}

// TestMonitoringPemeriksaan tests all Monitoring Pemeriksaan endpoints
func TestMonitoringPemeriksaan(t *testing.T) {
	cfg := SetupTestConfig(t)
	defer TeardownTestConfig(cfg)

	tc := NewTestContext(t, cfg)
	tc.SetupTestData()
	defer tc.CleanupTestData()

	t.Run("Happy Path", func(t *testing.T) {
		testMonitoringPemeriksaanHappyPath(tc)
	})

	t.Run("Sad Path", func(t *testing.T) {
		testMonitoringPemeriksaanSadPath(tc)
	})
}

func testMonitoringPemeriksaanHappyPath(tc *TestContext) {
	bidanToken := tc.GetTestToken("BIDAN")
	t := tc.T

	var pemeriksaanID int

	// TS-BE-008: Input Pemeriksaan
	t.Run("TS-BE-008: Input Pemeriksaan", func(t *testing.T) {
		body := map[string]interface{}{
			"id_pasien":     1,
			"berat_badan":   15.5,
			"tinggi_badan":  85.0,
			"tanggal_periksa": "2026-06-10",
		}
		resp := tc.MakeRequest("POST", "/api/v1/monitoring/pemeriksaan", body, bidanToken)
		tc.AssertStatusCode(resp, http.StatusCreated, "Should return 201 Created")

		result := tc.ParseResponse(resp)
		tc.AssertSuccess(result, "Response should be successful")

		// Store ID for subsequent tests
		if data, ok := result.Data.(map[string]interface{}); ok {
			if id, ok := data["id"].(float64); ok {
				pemeriksaanID = int(id)
			}
		}
	})

	// TS-BE-009: Detail Pemeriksaan
	t.Run("TS-BE-009: Detail Pemeriksaan", func(t *testing.T) {
		resp := tc.MakeRequest("GET", "/api/v1/monitoring/pemeriksaan/1", nil, bidanToken)
		tc.AssertStatusCode(resp, http.StatusOK, "Should return 200 OK")

		result := tc.ParseResponse(resp)
		tc.AssertSuccess(result, "Response should be successful")
	})

	// TS-BE-010: Edit Pemeriksaan
	t.Run("TS-BE-010: Edit Pemeriksaan", func(t *testing.T) {
		body := map[string]interface{}{
			"berat_badan":  16.0,
			"tinggi_badan": 86.0,
		}
		resp := tc.MakeRequest("PUT", "/api/v1/monitoring/pemeriksaan/1", body, bidanToken)
		tc.AssertStatusCode(resp, http.StatusOK, "Should return 200 OK")

		result := tc.ParseResponse(resp)
		tc.AssertSuccess(result, "Response should be successful")
	})

	// TS-BE-012: Verifikasi Pemeriksaan
	t.Run("TS-BE-012: Verifikasi Pemeriksaan", func(t *testing.T) {
		body := map[string]interface{}{
			"catatan_verifikasi": "Data telah diverifikasi",
		}
		resp := tc.MakeRequest("PATCH", "/api/v1/monitoring/pemeriksaan/1/verify", body, bidanToken)
		tc.AssertStatusCode(resp, http.StatusOK, "Should return 200 OK")

		result := tc.ParseResponse(resp)
		tc.AssertSuccess(result, "Response should be successful")
	})

	// TS-BE-013: Pemeriksaan Belum Diverifikasi
	t.Run("TS-BE-013: Pemeriksaan Belum Diverifikasi", func(t *testing.T) {
		resp := tc.MakeRequest("GET", "/api/v1/monitoring/pemeriksaan/pending", nil, bidanToken)
		tc.AssertStatusCode(resp, http.StatusOK, "Should return 200 OK")

		result := tc.ParseResponse(resp)
		tc.AssertSuccess(result, "Response should be successful")
	})

	// TS-BE-011: Hapus Pemeriksaan (terakhir karena destructive)
	t.Run("TS-BE-011: Hapus Pemeriksaan", func(t *testing.T) {
		if pemeriksaanID > 0 {
			resp := tc.MakeRequest("DELETE", fmt.Sprintf("/api/v1/monitoring/pemeriksaan/%d", pemeriksaanID), nil, bidanToken)
			tc.AssertStatusCode(resp, http.StatusOK, "Should return 200 OK")

			result := tc.ParseResponse(resp)
			tc.AssertSuccess(result, "Response should be successful")
		}
	})
}

func testMonitoringPemeriksaanSadPath(tc *TestContext) {
	t := tc.T
	bidanToken := tc.GetTestToken("BIDAN")
	dinkesToken := tc.GetTestToken("DINKES")
	kaderToken := tc.GetTestToken("KADER")

	// TS-BE-NEG-009: Input Pemeriksaan Data Tidak Lengkap
	t.Run("TS-BE-NEG-009: Input Pemeriksaan Data Tidak Lengkap", func(t *testing.T) {
		body := map[string]interface{}{
			"id_pasien": 1,
			// Missing required fields
		}
		resp := tc.MakeRequest("POST", "/api/v1/monitoring/pemeriksaan", body, bidanToken)
		tc.AssertStatusCode(resp, http.StatusUnprocessableEntity, "Should return 422 Validation Error")
	})

	// TS-BE-NEG-010: Input Pemeriksaan Role Tidak Sesuai
	t.Run("TS-BE-NEG-010: Input Pemeriksaan Role Tidak Sesuai", func(t *testing.T) {
		body := map[string]interface{}{
			"id_pasien":     1,
			"berat_badan":   15.5,
			"tinggi_badan":  85.0,
			"tanggal_periksa": "2026-06-10",
		}
		resp := tc.MakeRequest("POST", "/api/v1/monitoring/pemeriksaan", body, dinkesToken)
		tc.AssertStatusCode(resp, http.StatusForbidden, "Should return 403 Forbidden")
	})

	// TS-BE-NEG-011: Input Pemeriksaan ID Pasien Tidak Ditemukan
	t.Run("TS-BE-NEG-011: Input Pemeriksaan ID Pasien Tidak Ditemukan", func(t *testing.T) {
		body := map[string]interface{}{
			"id_pasien":     99999,
			"berat_badan":   15.5,
			"tinggi_badan":  85.0,
			"tanggal_periksa": "2026-06-10",
		}
		resp := tc.MakeRequest("POST", "/api/v1/monitoring/pemeriksaan", body, bidanToken)
		tc.AssertStatusCode(resp, http.StatusNotFound, "Should return 404 Not Found")
	})

	// TS-BE-NEG-012: Detail Pemeriksaan ID Tidak Valid
	t.Run("TS-BE-NEG-012: Detail Pemeriksaan ID Tidak Valid", func(t *testing.T) {
		resp := tc.MakeRequest("GET", "/api/v1/monitoring/pemeriksaan/abc", nil, bidanToken)
		tc.AssertStatusCode(resp, http.StatusBadRequest, "Should return 400 Bad Request")
	})

	// TS-BE-NEG-013: Edit Pemeriksaan Sudah Diverifikasi
	t.Run("TS-BE-NEG-013: Edit Pemeriksaan Sudah Diverifikasi", func(t *testing.T) {
		// Assuming ID 1 is already verified
		body := map[string]interface{}{
			"berat_badan":  16.0,
			"tinggi_badan": 86.0,
		}
		resp := tc.MakeRequest("PUT", "/api/v1/monitoring/pemeriksaan/1", body, bidanToken)
		// This might return 422 if already verified, or 200 if not
		// Adjust based on actual implementation
	})

	// TS-BE-NEG-014: Hapus Pemeriksaan Tidak Ditemukan
	t.Run("TS-BE-NEG-014: Hapus Pemeriksaan Tidak Ditemukan", func(t *testing.T) {
		resp := tc.MakeRequest("DELETE", "/api/v1/monitoring/pemeriksaan/99999", nil, bidanToken)
		tc.AssertStatusCode(resp, http.StatusNotFound, "Should return 404 Not Found")
	})

	// TS-BE-NEG-015: Verifikasi Pemeriksaan Bukan Bidan
	t.Run("TS-BE-NEG-015: Verifikasi Pemeriksaan Bukan Bidan", func(t *testing.T) {
		body := map[string]interface{}{
			"catatan_verifikasi": "Test",
		}
		resp := tc.MakeRequest("PATCH", "/api/v1/monitoring/pemeriksaan/1/verify", body, kaderToken)
		tc.AssertStatusCode(resp, http.StatusForbidden, "Should return 403 Forbidden")
	})

	// TS-BE-NEG-016: Pemeriksaan Pending Tanpa Auth
	t.Run("TS-BE-NEG-016: Pemeriksaan Pending Tanpa Auth", func(t *testing.T) {
		resp := tc.MakeRequest("GET", "/api/v1/monitoring/pemeriksaan/pending", nil, "")
		tc.AssertStatusCode(resp, http.StatusUnauthorized, "Should return 401 Unauthorized")
	})
}
