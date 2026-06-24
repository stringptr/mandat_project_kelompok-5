package test

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHealthCheck tests health check endpoint
func TestHealthCheck(t *testing.T) {
	cfg := SetupTestConfig(t)
	defer TeardownTestConfig(cfg)

	tc := NewTestContext(t, cfg)

	t.Run("Happy Path", func(t *testing.T) {
		testHealthCheckHappyPath(tc)
	})

	t.Run("Sad Path", func(t *testing.T) {
		testHealthCheckSadPath(tc)
	})
}

func testHealthCheckHappyPath(tc *TestContext) {
	t := tc.T

	// TS-BE-001: Health Check
	t.Run("TS-BE-001: Health Check", func(t *testing.T) {
		resp := tc.MakeRequest("GET", "/health", nil, "")
		tc.AssertStatusCode(resp, http.StatusOK, "Should return 200 OK")

		// Check response body
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err, "Should read response body")
		assert.Equal(t, "OK", string(body), "Should return OK")
	})
}

func testHealthCheckSadPath(tc *TestContext) {
	t := tc.T

	// TS-BE-NEG-001: Health Check Invalid Method
	t.Run("TS-BE-NEG-001: Health Check Invalid Method", func(t *testing.T) {
		resp := tc.MakeRequest("POST", "/health", nil, "")
		tc.AssertStatusCode(resp, http.StatusMethodNotAllowed, "Should return 405 Method Not Allowed")
	})
}
