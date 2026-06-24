# 🚀 Panduan Automated Testing Backend SiGizi

## 📋 Daftar Isi

1. [Overview](#overview)
2. [Prerequisites](#prerequisites)
3. [Setup Environment](#setup-environment)
4. [Menjalankan Test](#menjalankan-test)
5. [Struktur Test](#struktur-test)
6. [Coverage & Reporting](#coverage--reporting)
7. [Troubleshooting](#troubleshooting)

---

## Overview

Automated test ini mengcover **29 endpoint** dengan total **67 test cases**:
- ✅ **29 Happy Path** (skenario sukses)
- ❌ **38 Sad Path** (skenario error)

### Test Coverage:
- Health Check (2 test cases)
- Monitoring Pasien (9 test cases)
- Monitoring Statistik (4 test cases)
- Monitoring Pemeriksaan (14 test cases)
- Imunisasi (6 test cases)
- Tindak Lanjut (13 test cases)
- Laporan (2 test cases)
- Artikel (17 test cases)

---

## Prerequisites

### 1. Install Go

Pastikan Go versi 1.21+ terinstall:

```bash
go version
```

Jika belum, download dari: https://go.dev/dl/

### 2. Install PostgreSQL

Anda bisa menggunakan:
- PostgreSQL lokal (v16+)
- Docker container

### 3. Install Dependencies

```bash
cd backend
go mod download
```

Install testing libraries:

```bash
go get -u github.com/stretchr/testify
go get -u github.com/joho/godotenv
```

---

## Setup Environment

### Option 1: Menggunakan Docker (Recommended)

**Step 1:** Jalankan PostgreSQL container untuk testing

```bash
docker run -d \
  --name sigizi-test-db \
  -e POSTGRES_DB=sigizi_test \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 5433:5432 \
  postgres:16
```

**Step 2:** Verifikasi container berjalan

```bash
docker ps | grep sigizi-test-db
```

**Step 3:** Setup database schema

```bash
# Jika ada migration files
docker exec -i sigizi-test-db psql -U postgres -d sigizi_test < migrations/schema.sql

# Atau jalankan migration tool
# goose -dir migrations postgres "postgres://postgres:postgres@localhost:5433/sigizi_test?sslmode=disable" up
```

### Option 2: Menggunakan PostgreSQL Lokal

**Step 1:** Buat database test

```sql
CREATE DATABASE sigizi_test;
```

**Step 2:** Jalankan migrations

```bash
# Sesuaikan dengan migration tool yang digunakan
psql -U postgres -d sigizi_test -f migrations/schema.sql
```

### Konfigurasi Environment Variables

File `.env.test` sudah tersedia di `backend/.env.test`:

```env
# Database Configuration
DB_MASTER_HOST=localhost
DB_MASTER_PORT=5433
DB_MASTER_USER=postgres
DB_MASTER_PASSWORD=postgres
DB_MASTER_NAME=sigizi_test
DB_MASTER_SSLMODE=disable

# Server Configuration
SERVER_HOST=localhost
SERVER_PORT=8081

# JWT Configuration
JWT_SECRET=test_secret_key_for_testing_only
JWT_EXPIRATION=3600

# Test Mode
TEST_MODE=true
```

**⚠️ Penting:** Jangan gunakan database production untuk testing!

---

## Menjalankan Test

### 1. Test Semua Module

```bash
cd backend
go test ./test/... -v
```

Output:
```
=== RUN   TestHealthCheck
=== RUN   TestHealthCheck/Happy_Path
=== RUN   TestHealthCheck/Happy_Path/TS-BE-001:_Health_Check
--- PASS: TestHealthCheck (0.15s)
...
PASS
ok      myproject/test  5.234s
```

### 2. Test Module Spesifik

#### Test Health Check
```bash
go test ./test/ -v -run TestHealthCheck
```

#### Test Monitoring
```bash
go test ./test/ -v -run TestMonitoring
```

#### Test Artikel
```bash
go test ./test/ -v -run TestArtikel
```

#### Test Tindak Lanjut
```bash
go test ./test/ -v -run TestTindakLanjut
```

#### Test Imunisasi
```bash
go test ./test/ -v -run TestImunisasi
```

### 3. Test Specific Test Case

```bash
# Test hanya TS-BE-001
go test ./test/ -v -run TestHealthCheck/Happy_Path/TS-BE-001

# Test semua Happy Path untuk Monitoring
go test ./test/ -v -run TestMonitoring/Happy_Path
```

### 4. Test dengan Timeout

```bash
# Set timeout 10 menit
go test ./test/... -v -timeout 10m
```

---

## Coverage & Reporting

### 1. Generate Coverage Report

```bash
# Run test dengan coverage
go test ./test/... -v -cover -coverprofile=coverage.out

# Output:
# ok      myproject/test  5.234s  coverage: 75.5% of statements
```

### 2. View Coverage HTML

```bash
# Generate HTML report
go tool cover -html=coverage.out -o coverage.html

# Buka di browser
# Windows:
start coverage.html

# Linux/Mac:
open coverage.html
```

### 3. Coverage by Function

```bash
go tool cover -func=coverage.out
```

Output:
```
myproject/internal/handler/monitoring_handler.go:15:    GetPasienList           85.7%
myproject/internal/handler/monitoring_handler.go:42:    GetPasienDetail         92.3%
myproject/internal/service/monitoring_service.go:20:    GetPasien              100.0%
...
total:                                                   (statements)            75.5%
```

### 4. Generate JSON Report

```bash
go test ./test/... -v -json > test-report.json
```

---

## Struktur Test

```
backend/
├── test/
│   ├── README.md                 # Dokumentasi test
│   ├── config_test.go           # Test configuration & helpers
│   ├── health_test.go           # Health check tests
│   ├── monitoring_test.go       # Monitoring tests (Pasien, Statistik, Pemeriksaan)
│   ├── imunisasi_test.go        # Imunisasi tests
│   ├── tindaklanjut_test.go     # Tindak Lanjut & Laporan tests
│   └── artikel_test.go          # Artikel tests
├── .env.test                     # Test environment variables
└── go.mod
```

### Test Naming Convention

Test cases mengikuti naming dari Test Script Manual:

- **Happy Path**: `TS-BE-001`, `TS-BE-002`, etc.
- **Sad Path**: `TS-BE-NEG-001`, `TS-BE-NEG-002`, etc.

### Test Structure

```go
func TestModuleName(t *testing.T) {
    // Setup
    cfg := SetupTestConfig(t)
    defer TeardownTestConfig(cfg)
    
    tc := NewTestContext(t, cfg)
    tc.SetupTestData()
    defer tc.CleanupTestData()
    
    // Run tests
    t.Run("Happy Path", func(t *testing.T) {
        // Happy path test cases
    })
    
    t.Run("Sad Path", func(t *testing.T) {
        // Sad path test cases
    })
}
```

---

## Troubleshooting

### ❌ Error: Connection Refused

```
Error: dial tcp [::1]:5433: connect: connection refused
```

**Solusi:**
1. Pastikan PostgreSQL running
2. Check port yang digunakan
3. Verifikasi credentials di `.env.test`

```bash
# Check Docker container
docker ps

# Check PostgreSQL local
pg_isready -h localhost -p 5433
```

### ❌ Error: Database Does Not Exist

```
Error: database "sigizi_test" does not exist
```

**Solusi:**

```bash
# Create database
docker exec -it sigizi-test-db psql -U postgres -c "CREATE DATABASE sigizi_test;"

# Atau via psql lokal
createdb sigizi_test
```

### ❌ Error: Module Not Found

```
Error: package myproject/internal/handler is not in GOROOT
```

**Solusi:**

```bash
# Re-download dependencies
go mod tidy
go mod download

# Clear module cache
go clean -modcache
```

### ❌ Test Timeout

```
panic: test timed out after 10m0s
```

**Solusi:**
1. Increase timeout: `go test -timeout 30m`
2. Check database performance
3. Run specific tests instead of all

### ❌ Import Errors

```
Error: cannot find package "github.com/stretchr/testify"
```

**Solusi:**

```bash
go get -u github.com/stretchr/testify
go get -u github.com/joho/godotenv
go mod tidy
```

---

## Best Practices

### 1. Isolasi Data
- Setiap test menggunakan transaction rollback
- Clean database setelah setiap test suite
- Gunakan database terpisah untuk testing

### 2. Test Independence
- Test tidak bergantung pada urutan eksekusi
- Setiap test bisa dijalankan sendiri
- Hindari shared state antar test

### 3. Mock Authentication
- Gunakan mock token untuk testing
- Simulasi berbagai role (BIDAN, KADER, DINKES)
- Test authorization dengan token berbeda

### 4. Assertion yang Jelas
```go
// ❌ Bad
assert.True(t, resp.StatusCode == 200)

// ✅ Good
tc.AssertStatusCode(resp, http.StatusOK, "Should return 200 OK")
```

### 5. Cleanup
```go
// Always cleanup after test
defer tc.CleanupTestData()
defer TeardownTestConfig(cfg)
```

---

## Continuous Integration

### GitHub Actions Example

Buat file `.github/workflows/test.yml`:

```yaml
name: Backend Tests

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:16
        env:
          POSTGRES_DB: sigizi_test
          POSTGRES_USER: postgres
          POSTGRES_PASSWORD: postgres
        ports:
          - 5433:5432
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Install dependencies
        run: |
          cd backend
          go mod download
      
      - name: Run tests
        run: |
          cd backend
          go test ./test/... -v -cover -coverprofile=coverage.out
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./backend/coverage.out
          flags: backend
```

---

## Quick Reference

### Perintah Umum

| Perintah | Deskripsi |
|----------|-----------|
| `go test ./test/... -v` | Jalankan semua test dengan verbose output |
| `go test ./test/... -cover` | Test dengan coverage report |
| `go test -run TestName` | Jalankan test spesifik |
| `go test -short` | Skip long-running tests |
| `go test -parallel 4` | Run tests in parallel |
| `go clean -testcache` | Clear test cache |

### Docker Commands

| Perintah | Deskripsi |
|----------|-----------|
| `docker start sigizi-test-db` | Start test database |
| `docker stop sigizi-test-db` | Stop test database |
| `docker rm sigizi-test-db` | Remove test database |
| `docker logs sigizi-test-db` | View database logs |

---

## 📞 Support

Jika mengalami masalah:

1. Check [Troubleshooting](#troubleshooting) section
2. Review test logs: `go test ./test/... -v > test.log 2>&1`
3. Hubungi tim development

---

## 📝 Summary

Automated testing memberikan:
- ✅ **Confidence**: Semua endpoint tertest
- ✅ **Speed**: Test berjalan dalam hitungan detik
- ✅ **Documentation**: Test sebagai dokumentasi hidup
- ✅ **Regression Prevention**: Deteksi bug lebih awal
- ✅ **CI/CD Ready**: Integrasi dengan pipeline

Happy Testing! 🎉
