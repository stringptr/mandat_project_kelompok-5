# ⚡ Quick Start: Automated Testing

## 🎯 3 Langkah Mudah Menjalankan Test

### 1️⃣ Start Test Database

**Menggunakan Docker (Recommended):**

```bash
docker run -d \
  --name sigizi-test-db \
  -e POSTGRES_DB=sigizi_test \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 5433:5432 \
  postgres:16
```

**Atau menggunakan PostgreSQL Lokal:**

```sql
CREATE DATABASE sigizi_test;
```

### 2️⃣ Install Dependencies

```bash
cd backend
go mod download
```

### 3️⃣ Run Tests!

**Windows (PowerShell):**

```powershell
# Run all tests
.\run-tests.ps1

# Run specific module
.\run-tests.ps1 monitoring
.\run-tests.ps1 artikel

# Run with coverage
.\run-tests.ps1 coverage
```

**Linux/Mac (Bash):**

```bash
# Make script executable
chmod +x run-tests.sh

# Run all tests
./run-tests.sh

# Run specific module
./run-tests.sh monitoring
./run-tests.sh artikel

# Run with coverage
./run-tests.sh coverage
```

**Manual (Semua Platform):**

```bash
# Run all tests
go test ./test/... -v

# Run with coverage
go test ./test/... -v -cover
```

---

## 📊 Test Coverage

Total: **67 Test Cases**
- ✅ 29 Happy Path
- ❌ 38 Sad Path

**Modules:**
- Health Check (2 tests)
- Monitoring (27 tests)
- Imunisasi (6 tests)
- Tindak Lanjut (13 tests)
- Artikel (17 tests)
- Laporan (2 tests)

---

## 🔍 Perintah Berguna

```bash
# Test specific module
go test ./test/ -v -run TestMonitoring
go test ./test/ -v -run TestArtikel

# Test specific case
go test ./test/ -v -run TestHealthCheck/Happy_Path/TS-BE-001

# Generate coverage HTML
go test ./test/... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# Clear test cache
go clean -testcache
```

---

## ❓ Troubleshooting

### Database Connection Error

```bash
# Check if Docker container is running
docker ps | grep sigizi-test-db

# Start if stopped
docker start sigizi-test-db

# Check logs
docker logs sigizi-test-db
```

### Module Not Found

```bash
go mod tidy
go mod download
```

### Tests Fail

1. Check `.env.test` configuration
2. Verify database is running
3. Run migrations if needed
4. Clear test cache: `go clean -testcache`

---

## 📚 Full Documentation

Lihat [AUTOMATED-TEST-GUIDE.md](../AUTOMATED-TEST-GUIDE.md) untuk dokumentasi lengkap.

---

## ✅ Expected Output

```
🧪 SiGizi Backend Automated Test Runner
========================================

📊 Checking test database...
✅ Test database is running

🏃 Running ALL tests...
=== RUN   TestHealthCheck
=== RUN   TestHealthCheck/Happy_Path
=== RUN   TestHealthCheck/Happy_Path/TS-BE-001:_Health_Check
--- PASS: TestHealthCheck (0.15s)
=== RUN   TestMonitoringPasien
=== RUN   TestMonitoringPasien/Happy_Path
=== RUN   TestMonitoringPasien/Happy_Path/TS-BE-002:_Daftar_Pasien
--- PASS: TestMonitoringPasien (1.23s)
...
PASS
coverage: 75.5% of statements
ok      myproject/test  5.234s

✅ All tests passed!
```

---

Happy Testing! 🎉
