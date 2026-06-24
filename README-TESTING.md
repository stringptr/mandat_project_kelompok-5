# 🧪 Automated Testing - Backend SiGizi

> Framework automated testing lengkap untuk 29 endpoint backend dengan 67 test cases

---

## 🎯 Quick Start (3 Menit)

### 1. Start Test Database
```bash
docker run -d --name sigizi-test-db \
  -e POSTGRES_DB=sigizi_test \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 5433:5432 \
  postgres:16
```

### 2. Install Dependencies
```bash
cd backend
go mod download
```

### 3. Run Tests
**Windows:**
```powershell
.\run-tests.ps1
```

**Linux/Mac:**
```bash
chmod +x run-tests.sh && ./run-tests.sh
```

---

## 📊 Test Coverage

| Module | Happy Path | Sad Path | Total |
|--------|-----------|----------|-------|
| Health Check | 1 | 1 | 2 |
| Monitoring - Pasien | 4 | 5 | 9 |
| Monitoring - Statistik | 2 | 2 | 4 |
| Monitoring - Pemeriksaan | 6 | 8 | 14 |
| Imunisasi | 3 | 3 | 6 |
| Tindak Lanjut & Laporan | 7 | 7 | 13 |
| Artikel | 6 | 11 | 17 |
| **TOTAL** | **29** | **38** | **67** |

---

## 📁 File Structure

```
backend/
├── test/
│   ├── config_test.go           # Test helpers & utilities
│   ├── health_test.go           # Health check tests
│   ├── monitoring_test.go       # Monitoring tests
│   ├── imunisasi_test.go        # Imunisasi tests
│   ├── tindaklanjut_test.go     # Tindak lanjut tests
│   └── artikel_test.go          # Artikel tests
├── .env.test                     # Test environment config
├── run-tests.ps1                 # Windows test runner
├── run-tests.sh                  # Linux/Mac test runner
├── QUICK-START-TESTING.md        # Quick reference
└── test/README.md                # Technical docs
```

---

## 🚀 Cara Menggunakan

### Run All Tests
```bash
go test ./test/... -v
```

### Run Specific Module
```bash
go test ./test/ -v -run TestMonitoring
go test ./test/ -v -run TestArtikel
go test ./test/ -v -run TestTindakLanjut
```

### Run with Coverage
```bash
go test ./test/... -v -cover -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Using Helper Scripts

**Windows:**
```powershell
.\run-tests.ps1              # All tests
.\run-tests.ps1 monitoring   # Monitoring only
.\run-tests.ps1 coverage     # With coverage report
```

**Linux/Mac:**
```bash
./run-tests.sh              # All tests
./run-tests.sh monitoring   # Monitoring only
./run-tests.sh coverage     # With coverage report
```

---

## ✨ Features

- ✅ **67 Test Cases** covering all 29 endpoints
- ✅ **Happy & Sad Path** scenarios
- ✅ **Automated Setup** dengan helper scripts
- ✅ **Coverage Reports** dengan HTML visualization
- ✅ **Mock Authentication** untuk testing berbagai role
- ✅ **Database Isolation** setiap test independent
- ✅ **CI/CD Ready** dengan GitHub Actions support
- ✅ **Complete Documentation** dengan troubleshooting guide

---

## 📖 Dokumentasi Lengkap

| File | Deskripsi | Ukuran |
|------|-----------|--------|
| [`TESTING-SUMMARY.md`](TESTING-SUMMARY.md) | Overview & summary lengkap | ~200 baris |
| [`AUTOMATED-TEST-GUIDE.md`](AUTOMATED-TEST-GUIDE.md) | Panduan lengkap automated testing | ~500 baris |
| [`backend/QUICK-START-TESTING.md`](backend/QUICK-START-TESTING.md) | Quick reference guide | ~100 baris |
| [`backend/test/README.md`](backend/test/README.md) | Technical documentation | ~300 baris |

---

## 🔧 Test Helper Functions

```go
// Setup & Teardown
SetupTestConfig(t)              // Initialize test env
TeardownTestConfig(cfg)         // Cleanup

// Test Context
tc := NewTestContext(t, cfg)
tc.SetupTestData()              // Insert test data
tc.CleanupTestData()            // Remove test data

// HTTP Requests
resp := tc.MakeRequest(method, path, body, token)
result := tc.ParseResponse(resp)

// Assertions
tc.AssertStatusCode(resp, http.StatusOK)
tc.AssertSuccess(result)
tc.AssertError(result)

// Authentication
token := tc.GetTestToken("BIDAN")
```

---

## 🎯 Test Scenarios

### Happy Path (29 tests) ✅
- Semua request dengan data valid
- Authentication & authorization benar
- Expected response 200/201

### Sad Path (38 tests) ❌
- Request tanpa authentication (401)
- Request dengan token invalid (401)
- Request dengan role salah (403)
- Request dengan ID tidak ditemukan (404)
- Request dengan data invalid (422)
- Request dengan parameter salah (400)

---

## 🐛 Troubleshooting

### Database Connection Error
```bash
docker ps | grep sigizi-test-db      # Check status
docker start sigizi-test-db          # Start if stopped
docker logs sigizi-test-db           # Check logs
```

### Module Not Found
```bash
go mod tidy
go mod download
go clean -modcache
```

### Test Timeout
```bash
go test -timeout 30m ./test/...
```

### Clear Test Cache
```bash
go clean -testcache
```

---

## 🔄 CI/CD Integration

### GitHub Actions Example

```yaml
name: Backend Tests

on: [push, pull_request]

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
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - run: cd backend && go test ./test/... -v -cover
```

---

## 📈 Expected Output

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
=== RUN   TestMonitoringPasien/Happy_Path/TS-BE-003:_Detail_Pasien
--- PASS: TestMonitoringPasien (1.23s)

...

PASS
coverage: 75.5% of statements
ok      myproject/test  5.234s

✅ All tests passed!
```

---

## 💡 Tips & Best Practices

1. **Selalu gunakan database terpisah untuk testing**
2. **Clean test cache** jika test berperilaku aneh
3. **Run specific tests** saat development untuk speed
4. **Check coverage regularly** - target > 75%
5. **Review test failures** - jangan skip yang fail
6. **Update tests** saat menambah/ubah endpoint

---

## 🎓 Learning Resources

- [Go Testing Documentation](https://go.dev/doc/tutorial/add-a-test)
- [Testify Library](https://github.com/stretchr/testify)
- [HTTP Testing in Go](https://pkg.go.dev/net/http/httptest)
- [Test Coverage in Go](https://go.dev/blog/cover)

---

## 📞 Support

**Dokumentasi:**
- Baca [`TESTING-SUMMARY.md`](TESTING-SUMMARY.md) untuk overview
- Baca [`AUTOMATED-TEST-GUIDE.md`](AUTOMATED-TEST-GUIDE.md) untuk detail
- Baca [`QUICK-START-TESTING.md`](backend/QUICK-START-TESTING.md) untuk quick reference

**Issues:**
- Check troubleshooting section di dokumentasi
- Review test logs
- Contact development team

---

## ✅ Checklist

Sebelum commit/push, pastikan:

- [ ] Semua tests passing
- [ ] Coverage > 75%
- [ ] Tidak ada test yang di-skip
- [ ] Test data cleaned up properly
- [ ] Documentation updated (jika ada perubahan)

---

## 📊 Statistics

- **Total Endpoints**: 29
- **Total Test Cases**: 67
- **Happy Path Tests**: 29
- **Sad Path Tests**: 38
- **Test Files**: 6
- **Helper Functions**: 15+
- **Documentation Pages**: 4
- **Lines of Test Code**: ~2000+

---

**🎉 Happy Testing!**

Framework ini siap digunakan untuk automated testing backend SiGizi.  
Semua 67 test cases sudah ter-cover dan siap di-run! 🚀
