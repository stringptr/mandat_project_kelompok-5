package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"myproject/internal/config"
	"myproject/internal/handler"
	"myproject/internal/repository"
	"myproject/internal/service"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConfig holds test configuration
type TestConfig struct {
	DBPool *pgxpool.Pool
	Router *chi.Mux
	API    huma.API
	Server *httptest.Server
}

// TestContext holds test execution context
type TestContext struct {
	T      *testing.T
	Config *TestConfig
	Tokens map[string]string // role -> token
}

// Response represents API response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SetupTestConfig initializes test environment
func SetupTestConfig(t *testing.T) *TestConfig {
	// Load test environment variables
	if err := godotenv.Load("../.env.test"); err != nil {
		t.Logf("Warning: .env.test not found, using default config")
	}

	// Load configuration
	cfg := config.Load()

	// Setup test database
	dbPool, err := setupTestDatabase(cfg.DBMasterConfig)
	require.NoError(t, err, "Failed to setup test database")

	// Initialize repositories
	monitoringRepo := repository.NewMonitoringRepository(dbPool)
	tindakLanjutRepo := repository.NewTindakLanjutRepository(dbPool)
	artikelRepo := repository.NewArtikelRepository(dbPool)

	// Initialize services
	monitoringService := service.NewMonitoringService(monitoringRepo)
	tindakLanjutService := service.NewTindakLanjutService(tindakLanjutRepo)
	artikelService := service.NewArtikelService(artikelRepo)

	// Initialize handlers
	monitoringHandler := handler.NewMonitoringHandler(monitoringService)
	tindakLanjutHandler := handler.NewTindakLanjutHandler(tindakLanjutService)
	artikelHandler := handler.NewArtikelHandler(artikelService)

	// Setup router
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Health check
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Setup Huma API
	humaConfig := huma.DefaultConfig("SiGizi Test API", "1.0.0")
	api := humachi.New(router, humaConfig)

	// Register routes
	monitoringHandler.RegisterRoutes(api)
	tindakLanjutHandler.RegisterRoutes(api)
	artikelHandler.RegisterRoutes(api)

	// Create test server
	server := httptest.NewServer(router)

	return &TestConfig{
		DBPool: dbPool,
		Router: router,
		API:    api,
		Server: server,
	}
}

// setupTestDatabase creates test database connection
func setupTestDatabase(cfg config.PostgresConfig) (*pgxpool.Pool, error) {
	dsn := cfg.DSN()

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database config: %w", err)
	}

	// Test pool settings
	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return pool, nil
}

// TeardownTestConfig cleans up test environment
func TeardownTestConfig(cfg *TestConfig) {
	if cfg.Server != nil {
		cfg.Server.Close()
	}
	if cfg.DBPool != nil {
		cfg.DBPool.Close()
	}
}

// NewTestContext creates a new test context
func NewTestContext(t *testing.T, cfg *TestConfig) *TestContext {
	return &TestContext{
		T:      t,
		Config: cfg,
		Tokens: make(map[string]string),
	}
}

// MakeRequest performs HTTP request and returns response
func (tc *TestContext) MakeRequest(method, path string, body interface{}, token string) *http.Response {
	var reqBody io.Reader

	if body != nil {
		jsonBytes, err := json.Marshal(body)
		require.NoError(tc.T, err, "Failed to marshal request body")
		reqBody = bytes.NewBuffer(jsonBytes)
	}

	req, err := http.NewRequest(method, tc.Config.Server.URL+path, reqBody)
	require.NoError(tc.T, err, "Failed to create request")

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	require.NoError(tc.T, err, "Failed to perform request")

	return resp
}

// ParseResponse parses HTTP response to Response struct
func (tc *TestContext) ParseResponse(resp *http.Response) *Response {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(tc.T, err, "Failed to read response body")

	var response Response
	if len(body) > 0 {
		err = json.Unmarshal(body, &response)
		if err != nil {
			// If not JSON, treat as plain text
			response.Message = string(body)
		}
	}

	return &response
}

// AssertStatusCode checks HTTP status code
func (tc *TestContext) AssertStatusCode(resp *http.Response, expected int, msgAndArgs ...interface{}) {
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("TEST_JSON_DUMP|%s|%s\n", tc.T.Name(), string(body))
	resp.Body = io.NopCloser(bytes.NewBuffer(body))
	assert.Equal(tc.T, expected, resp.StatusCode, msgAndArgs...)
}

// AssertSuccess checks if response is successful
func (tc *TestContext) AssertSuccess(resp *Response, msgAndArgs ...interface{}) {
	assert.True(tc.T, resp.Success, msgAndArgs...)
}

// AssertError checks if response contains error
func (tc *TestContext) AssertError(resp *Response, msgAndArgs ...interface{}) {
	assert.False(tc.T, resp.Success, msgAndArgs...)
	assert.NotEmpty(tc.T, resp.Error, msgAndArgs...)
}

// SetupTestData inserts test data into database
func (tc *TestContext) SetupTestData() {
	ctx := context.Background()

	// Insert test users
	_, err := tc.Config.DBPool.Exec(ctx, `
		INSERT INTO user_account (id, email, role, password_hash, is_verified, created_at)
		VALUES 
			(1, 'bidan@test.com', 'BIDAN', '$argon2id$v=19$m=65536,t=3,p=4$test', true, NOW()),
			(2, 'kader@test.com', 'KADER', '$argon2id$v=19$m=65536,t=3,p=4$test', true, NOW()),
			(3, 'dinkes@test.com', 'DINKES', '$argon2id$v=19$m=65536,t=3,p=4$test', true, NOW()),
			(4, 'ibu@test.com', 'IBU', '$argon2id$v=19$m=65536,t=3,p=4$test', true, NOW())
		ON CONFLICT DO NOTHING
	`)
	require.NoError(tc.T, err, "Failed to insert test users")

	// Insert test pasien data
	_, err = tc.Config.DBPool.Exec(ctx, `
		INSERT INTO pasien (id, nik, nama_lengkap, tanggal_lahir, jenis_kelamin, created_at)
		VALUES 
			(1, '1234567890123456', 'Test Pasien 1', '2020-01-01', 'L', NOW()),
			(2, '1234567890123457', 'Test Pasien 2', '2020-02-01', 'P', NOW())
		ON CONFLICT DO NOTHING
	`)
	require.NoError(tc.T, err, "Failed to insert test pasien")
}

// CleanupTestData removes test data from database
func (tc *TestContext) CleanupTestData() {
	ctx := context.Background()

	// Clean up in reverse order of foreign keys
	tables := []string{
		"tindaklanjut",
		"hasil_pemeriksaan",
		"artikel",
		"jadwal_imunisasi",
		"pasien",
		"user_session",
		"user_account",
	}

	for _, table := range tables {
		_, err := tc.Config.DBPool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			tc.T.Logf("Warning: failed to truncate %s: %v", table, err)
		}
	}
}

// GenerateTestToken generates a mock JWT token for testing
func (tc *TestContext) GenerateTestToken(role string) string {
	// In real implementation, this should generate actual JWT
	// For testing, we use a simple mock token
	token := fmt.Sprintf("test_token_%s_%d", role, time.Now().Unix())
	tc.Tokens[role] = token
	return token
}

// GetTestToken retrieves or generates test token for role
func (tc *TestContext) GetTestToken(role string) string {
	if token, exists := tc.Tokens[role]; exists {
		return token
	}
	return tc.GenerateTestToken(role)
}
