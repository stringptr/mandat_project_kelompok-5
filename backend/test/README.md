# Automated Test Script untuk Backend SiGizi

## Overview

Automated test script ini mengcover semua 29 endpoint dengan skenario Happy Path dan Sad Path sesuai dengan Test Script Manual yang telah dibuat.

## Struktur Test

```
backend/test/
├── README.md                 # Dokumentasi ini
├── config_test.go           # Konfigurasi test & test helpers
├── setup_test.go            # Setup database & test data
├── monitoring_test.go       # Test untuk modul Monitoring
├── imunisasi_test.go        # Test untuk modul Imunisasi
├── tindaklanjut_test.go     # Test untuk modul Tindak Lanjut
├── artikel_test.go          # Test untuk modul Artikel
└── integration_test.go      # Integration test untuk semua modul
```

## Prerequisites

1. **Go 1.21+** terinstall
2. **PostgreSQL** database untuk testing
3. **Docker** (optional, untuk isolated test database)

## Setup Test Database

### Option 1: Menggunakan Database Lokal

Buat database khusus untuk testing:

```sql
CREATE DATABASE sigizi_test;
```

### Option 2: Menggunakan Docker

```bash
docker run -d \
  --name sigizi-test-db \
  -e POSTGRES_DB=sigizi_test \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 5433:5432 \
  postgres:16
```

## Environment Variables

Buat file `.env.test` di folder `backend/`:

```env
DB_MASTER_HOST=localhost
DB_MASTER_PORT=5433
DB_MASTER_USER=postgres
DB_MASTER_PASSWORD=postgres
DB_MASTER_NAME=sigizi_test
DB_MASTER_SSLMODE=disable
```

## Menjalankan Test

### 1. Jalankan Semua Test

```bash
cd backend
go test ./test/... -v
```

### 2. Jalankan Test Spesifik Module

```bash
# Test Monitoring saja
go test ./test/ -v -run TestMonitoring

# Test Artikel saja
go test ./test/ -v -run TestArtikel

# Test Tindak Lanjut saja
go test ./test/ -v -run TestTindakLanjut
```

### 3. Jalankan Test dengan Coverage

```bash
go test ./test/... -v -cover -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### 4. Jalankan Integration Test

```bash
go test ./test/ -v -run TestIntegration
```

## Test Categories

### Happy Path Tests
- ✅ Semua request dengan data valid
- ✅ Authentication & authorization yang benar
- ✅ Expected response 200/201

### Sad Path Tests
- ❌ Request tanpa authentication (401)
- ❌ Request dengan token invalid (401)
- ❌ Request dengan role yang salah (403)
- ❌ Request dengan ID yang tidak ditemukan (404)
- ❌ Request dengan data tidak lengkap/invalid (422)
- ❌ Request dengan parameter yang salah (400)

## Test Report

Setelah menjalankan test, laporan akan dibuat dalam format:

1. **Console Output** - Real-time test results
2. **coverage.html** - Coverage report (jika menggunakan -cover flag)
3. **test-report.json** - JSON format untuk CI/CD integration

## Continuous Integration

Test ini dapat diintegrasikan dengan GitHub Actions:

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
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run Tests
        run: |
          cd backend
          go test ./test/... -v -cover
```

## Troubleshooting

### Database Connection Error

```
Error: failed to connect to database
```

**Solution**: Pastikan database test sudah running dan credentials di `.env.test` benar.

### Port Already in Use

```
Error: bind: address already in use
```

**Solution**: Stop service yang menggunakan port 5433 atau ubah port di `.env.test`.

### Test Timeout

```
Error: test timed out
```

**Solution**: Increase timeout di test configuration atau check database performance.

## Best Practices

1. **Isolasi Data**: Setiap test menggunakan transaction yang di-rollback
2. **Test Order**: Test harus bisa dijalankan dalam urutan apapun
3. **Clean State**: Setiap test dimulai dengan clean database state
4. **Mock External Services**: Email, notifications, dll menggunakan mock
5. **Meaningful Assertions**: Setiap assertion harus jelas dan specific

## Contributing

Saat menambah endpoint baru:

1. Tambahkan test case di module yang sesuai
2. Include happy path dan sad path scenarios
3. Update test count di documentation
4. Pastikan semua test passing sebelum commit

## Contact

Untuk pertanyaan atau issues, hubungi tim development.
