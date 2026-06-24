# Test Script Backend - SIGizi

Dokumentasi lengkap untuk Test Script Backend sistem SIGizi yang mencakup 29 endpoint dengan skenario Happy Path dan Sad Path.

## 📁 File yang Tersedia

1. **TestScript-Backend-Complete.md** - Dokumentasi lengkap dalam format Markdown
2. **TestScript-Backend-HappyPath.csv** - Test cases positif dalam format CSV (bisa diimport ke Excel)
3. **TestScript-Backend-SadPath.csv** - Test cases negatif dalam format CSV (bisa diimport ke Excel)
4. **TestScript-README.md** - File ini

## 📊 Ringkasan Test Coverage

### Total Endpoints: 29
### Total Test Cases: 67
- **Happy Path (Positif)**: 29 test cases
- **Sad Path (Negatif)**: 38 test cases

## 🎯 Endpoint Coverage

### 1. Health Check (1 endpoint)
- `GET /health`

### 2. Monitoring - Pasien (4 endpoints)
- `GET /api/v1/monitoring/pasien` - Daftar Pasien
- `GET /api/v1/monitoring/pasien/{id}` - Detail Pasien
- `GET /api/v1/monitoring/pasien/search` - Cari Pasien
- `GET /api/v1/monitoring/pasien/export` - Export Excel

### 3. Monitoring - Statistik (2 endpoints)
- `GET /api/v1/monitoring/statistik` - Summary Dashboard
- `GET /api/v1/monitoring/statistik-bulanan` - Tren Pertumbuhan Bulanan

### 4. Monitoring - Pemeriksaan (6 endpoints)
- `POST /api/v1/monitoring/pemeriksaan` - Input Pemeriksaan
- `GET /api/v1/monitoring/pemeriksaan/{id}` - Detail Pemeriksaan
- `PUT /api/v1/monitoring/pemeriksaan/{id}` - Edit Pemeriksaan
- `DELETE /api/v1/monitoring/pemeriksaan/{id}` - Hapus Pemeriksaan
- `PATCH /api/v1/monitoring/pemeriksaan/{id}/verify` - Verifikasi Pemeriksaan
- `GET /api/v1/monitoring/pemeriksaan/pending` - Pemeriksaan Belum Diverifikasi

### 5. Imunisasi (3 endpoints)
- `GET /api/v1/imunisasi` - Daftar Jadwal Imunisasi
- `POST /api/v1/imunisasi` - Tambah Jadwal Imunisasi
- `PATCH /api/v1/imunisasi/{id}/realisasi` - Update Realisasi Vaksin

### 6. Tindak Lanjut (6 endpoints)
- `GET /api/v1/tindak-lanjut/pasien` - Daftar Pasien Tindak Lanjut
- `GET /api/v1/tindak-lanjut/pasien/{id}` - Detail Pasien Tindak Lanjut
- `POST /api/v1/tindak-lanjut` - Membuat Tindak Lanjut
- `PATCH /api/v1/rujukan/{id}/status` - Update Status Rujukan
- `GET /api/v1/tindak-lanjut/{id}` - Detail Tindak Lanjut
- `GET /api/v1/tindak-lanjut/status` - Status Tindak Lanjut (Kader)

### 7. Laporan (1 endpoint)
- `GET /api/v1/laporan/tindak-lanjut` - Laporan Tindak Lanjut

### 8. Artikel (6 endpoints)
- `GET /api/v1/artikel` - Daftar Artikel
- `GET /api/v1/artikel/{id}` - Detail Artikel
- `POST /api/v1/artikel` - Membuat Artikel
- `PATCH /api/v1/artikel/{id}` - Update Artikel
- `PATCH /api/v1/artikel/{id}/review` - Review Artikel
- `GET /api/v1/artikel/pending` - Daftar Artikel Pending
- `DELETE /api/v1/artikel/{id}` - Hapus Artikel

## 🔧 Cara Menggunakan File CSV

### Import ke Excel/Google Sheets:

1. Buka Excel atau Google Sheets
2. File > Import atau File > Open
3. Pilih file CSV yang ingin diimport:
   - `TestScript-Backend-HappyPath.csv` untuk skenario positif
   - `TestScript-Backend-SadPath.csv` untuk skenario negatif
4. Pilih delimiter: **Comma (,)**
5. Atur format kolom sesuai kebutuhan
6. Data akan muncul dalam bentuk tabel

### Struktur Kolom CSV:

| Kolom | Deskripsi |
|-------|-----------|
| No | Nomor urut test case |
| SRS Ref | Referensi ke dokumen SRS |
| FSD Ref | Referensi ke dokumen FSD |
| TSD Ref | Referensi ke dokumen TSD |
| No. Test Script | ID unik test script |
| Nama Fungsional | Nama fungsi yang ditest |
| URL / Endpoint | Endpoint API yang ditest |
| Parameter / Value | Parameter dan nilai yang digunakan |
| Ekspektasi | Hasil yang diharapkan |
| Hasil Aktual | Hasil yang didapat saat testing |
| Tgl Test | Tanggal pelaksanaan test |
| PIC | Person in Charge (tester) |
| Screenshot Hasil | Referensi ke screenshot |
| Pass / Failed | Status test (Pass/Failed) |

## 🧪 Cara Melakukan Testing

### Prerequisites:
1. Backend server sudah running di `http://localhost:8080`
2. Database sudah di-setup dan terisi data dummy
3. Miliki token JWT untuk setiap role:
   - Token BIDAN
   - Token DINKES (Dinas Kesehatan)
   - Token KADER
   - Token IBU/WALI

### Menggunakan Postman:

1. **Setup Environment**:
   ```
   BASE_URL = http://localhost:8080
   TOKEN_BIDAN = <your_bidan_token>
   TOKEN_DINKES = <your_dinkes_token>
   TOKEN_KADER = <your_kader_token>
   TOKEN_IBU = <your_ibu_token>
   ```

2. **Import Collection**:
   - Buat collection baru di Postman
   - Tambahkan request untuk setiap endpoint
   - Gunakan environment variables untuk token

3. **Running Tests**:
   - Jalankan Happy Path terlebih dahulu
   - Kemudian jalankan Sad Path
   - Dokumentasikan hasil setiap test

### Menggunakan cURL:

```bash
# Contoh: Test Health Check
curl -X GET http://localhost:8080/health

# Contoh: Test Daftar Pasien (dengan auth)
curl -X GET "http://localhost:8080/api/v1/monitoring/pasien?page=1&per_page=15" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"

# Contoh: Test Input Pemeriksaan
curl -X POST http://localhost:8080/api/v1/monitoring/pemeriksaan \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "id_pasien": 1,
    "berat_badan": 12.5,
    "tinggi_badan": 85.0,
    "lingkar_lengan": 15.0
  }'
```

## 📝 Skenario Testing

### Happy Path (Positif):
- Semua parameter valid
- Token authentication benar
- Role user sesuai dengan requirement
- Data yang direferensikan ada di database
- **Ekspektasi**: Response sukses (200, 201)

### Sad Path (Negatif):
- Tanpa token authentication → 401 Unauthorized
- Token invalid/expired → 401 Unauthorized
- Role tidak sesuai → 403 Forbidden
- Data tidak ditemukan → 404 Not Found
- Validasi data gagal → 422 Unprocessable Entity
- Format parameter salah → 400 Bad Request
- Method HTTP salah → 405 Method Not Allowed

## 🎯 Tips Testing

1. **Urutan Testing**:
   - Mulai dari Health Check
   - Test endpoint GET (read) dulu
   - Baru test endpoint POST/PUT/PATCH/DELETE (write)
   - Test authorization di setiap endpoint

2. **Data Preparation**:
   - Siapkan data dummy yang konsisten
   - Reset database setelah setiap test cycle
   - Gunakan transaction untuk rollback jika perlu

3. **Documentation**:
   - Screenshot setiap response
   - Catat response time
   - Dokumentasikan edge cases yang ditemukan

4. **Best Practices**:
   - Jalankan test di environment staging dulu
   - Jangan test di production
   - Gunakan automated testing untuk CI/CD
   - Monitor error logs saat testing

## 🔍 Validasi Response

Setiap response harus memiliki struktur:

```json
{
  "status": 200,
  "title": "OK",
  "success": true,
  "detail": "Deskripsi response",
  "data": { ... }
}
```

Untuk error response:

```json
{
  "status": 400,
  "title": "Bad Request",
  "success": false,
  "detail": "Deskripsi error",
  "errors": [
    {
      "location": "body.field_name",
      "message": "Error message",
      "value": "invalid_value"
    }
  ]
}
```

## 📊 Report Template

Gunakan template ini untuk melaporkan hasil testing:

```
TEST EXECUTION REPORT

Tanggal: [DD/MM/YYYY]
Tester: [Nama Tester]
Environment: [Development/Staging/Production]

SUMMARY:
- Total Test Cases: 67
- Passed: [X]
- Failed: [Y]
- Blocked: [Z]
- Pass Rate: [X/67 * 100]%

FAILED TEST CASES:
1. [Test Case ID] - [Alasan]
2. ...

NOTES:
- [Catatan tambahan]
- [Bug yang ditemukan]
- [Saran improvement]
```

## 🚀 Automation

Untuk automated testing, bisa menggunakan:

1. **Newman** (CLI runner untuk Postman):
   ```bash
   newman run postman_collection.json -e environment.json
   ```

2. **Jest + Supertest**:
   ```javascript
   describe('GET /api/v1/monitoring/pasien', () => {
     it('should return patient list', async () => {
       const response = await request(app)
         .get('/api/v1/monitoring/pasien')
         .set('Authorization', `Bearer ${token}`)
         .expect(200);
       
       expect(response.body.success).toBe(true);
     });
   });
   ```

3. **GitHub Actions CI/CD**:
   ```yaml
   name: Backend API Tests
   on: [push, pull_request]
   jobs:
     test:
       runs-on: ubuntu-latest
       steps:
         - uses: actions/checkout@v2
         - name: Run API Tests
           run: npm run test:api
   ```

## 📞 Contact

Jika ada pertanyaan atau menemukan bug, hubungi:
- Team: Kelompok A
- Kelas: A
- Angkatan: 2024

---

**Last Updated**: 10 Juni 2026
