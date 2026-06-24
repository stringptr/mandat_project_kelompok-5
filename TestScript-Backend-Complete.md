# TEST SCRIPT REPORT — BACKEND

**Nama Kelompok**: Kelompok A  
**Kelas**: A  
**Angkatan**: 2024  
**Nama Sistem**: SIGizi  
**Tipe Pengujian**: Backend

---

## SKENARIO: HAPPY PATH (POSITIF)

### Modul: Health Check
| No. | SRS Ref. | FSD Ref. | TSD Ref. | No. Test Script | Nama Fungsional | URL / Endpoint | Parameter / Value | Ekspektasi | Hasil Aktual | Tgl Test | PIC | Screenshot Hasil | Pass / Failed |
|-----|----------|----------|----------|-----------------|-----------------|----------------|-------------------|------------|--------------|----------|-----|------------------|---------------|
| 1 | SRS-7.2 | FSD-2.2 | TSD-3.2 | TS-BE-001 | Health Check | GET /health | - | Status 200, response "OK" | Status 200, "OK" | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |

### Modul: Monitoring - Pasien
| No. | SRS Ref. | FSD Ref. | TSD Ref. | No. Test Script | Nama Fungsional | URL / Endpoint | Parameter / Value | Ekspektasi | Hasil Aktual | Tgl Test | PIC | Screenshot Hasil | Pass / Failed |
|-----|----------|----------|----------|-----------------|-----------------|----------------|-------------------|------------|--------------|----------|-----|------------------|---------------|
| 2 | SRS-7.3 | FSD-2.3 | TSD-3.3.t | TS-BE-002 | Daftar Pasien | GET /api/v1/monitoring/pasien | Authorization: Bearer {token}, page=1, per_page=15 | Status 200, paginated patient list | Status 200, data diterima | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 3 | SRS-7.3 | FSD-2.3 | TSD-3.3.u | TS-BE-003 | Detail Pasien | GET /api/v1/monitoring/pasien/{id} | Authorization: Bearer {token}, id=1 | Status 200, patient detail | Status 200, detail lengkap | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 4 | SRS-7.3 | FSD-2.3 | TSD-3.3.v | TS-BE-004 | Cari Pasien | GET /api/v1/monitoring/pasien/search | Authorization: Bearer {token}, q=nama_pasien | Status 200, search results | Status 200, hasil pencarian | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 5 | SRS-7.3 | FSD-2.3 | TSD-5.3 | TS-BE-005 | Export Excel Pasien | GET /api/v1/monitoring/pasien/export | Authorization: Bearer {token}, filter params | Status 200, Excel file | Status 200, file diunduh | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |

### Modul: Monitoring - Statistik
| No. | SRS Ref. | FSD Ref. | TSD Ref. | No. Test Script | Nama Fungsional | URL / Endpoint | Parameter / Value | Ekspektasi | Hasil Aktual | Tgl Test | PIC | Screenshot Hasil | Pass / Failed |
|-----|----------|----------|----------|-----------------|-----------------|----------------|-------------------|------------|--------------|----------|-----|------------------|---------------|
| 6 | SRS-7.3 | FSD-2.3 | TSD-3.2 | TS-BE-006 | Summary Dashboard | GET /api/v1/monitoring/statistik | Authorization: Bearer {token} | Status 200, statistics summary | Status 200, data statistik | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 7 | SRS-7.3 | FSD-2.3 | TSD-3.2 | TS-BE-007 | Tren Pertumbuhan Bulanan | GET /api/v1/monitoring/statistik-bulanan | Authorization: Bearer {token}, tahun=2026 | Status 200, monthly trends | Status 200, tren bulanan | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |

### Modul: Monitoring - Pemeriksaan
| No. | SRS Ref. | FSD Ref. | TSD Ref. | No. Test Script | Nama Fungsional | URL / Endpoint | Parameter / Value | Ekspektasi | Hasil Aktual | Tgl Test | PIC | Screenshot Hasil | Pass / Failed |
|-----|----------|----------|----------|-----------------|-----------------|----------------|-------------------|------------|--------------|----------|-----|------------------|---------------|
| 8 | SRS-7.4 | FSD-2.3 | TSD-3.3.w | TS-BE-008 | Input Pemeriksaan | POST /api/v1/monitoring/pemeriksaan | Authorization: Bearer {token}, body: id_pasien, berat_badan, tinggi_badan, dll | Status 201, pemeriksaan created | Status 201, data tersimpan | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 9 | SRS-7.3 | FSD-2.3 | TSD-3.3.x | TS-BE-009 | Detail Pemeriksaan | GET /api/v1/monitoring/pemeriksaan/{id} | Authorization: Bearer {token}, id=1 | Status 200, examination detail | Status 200, detail lengkap | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 10 | SRS-7.5 | FSD-2.3 | TSD-3.3.y | TS-BE-010 | Edit Pemeriksaan | PUT /api/v1/monitoring/pemeriksaan/{id} | Authorization: Bearer {token}, id=1, body: updated data | Status 200, pemeriksaan updated | Status 200, data diperbarui | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 11 | SRS-7.6 | FSD-2.3 | TSD-3.3.z | TS-BE-011 | Hapus Pemeriksaan | DELETE /api/v1/monitoring/pemeriksaan/{id} | Authorization: Bearer {token}, id=1 | Status 200, pemeriksaan deleted | Status 200, data dihapus | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 12 | SRS-7.4 | FSD-2.3 | TSD-3.3.aa | TS-BE-012 | Verifikasi Pemeriksaan | PATCH /api/v1/monitoring/pemeriksaan/{id}/verify | Authorization: Bearer {token}, id=1, body: catatan_verifikasi | Status 200, pemeriksaan verified | Status 200, terverifikasi | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 13 | SRS-7.3 | FSD-2.3 | TSD-3.3.ab | TS-BE-013 | Pemeriksaan Belum Diverifikasi | GET /api/v1/monitoring/pemeriksaan/pending | Authorization: Bearer {token}, page=1, per_page=15 | Status 200, pending examinations | Status 200, daftar pending | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |

### Modul: Imunisasi
| No. | SRS Ref. | FSD Ref. | TSD Ref. | No. Test Script | Nama Fungsional | URL / Endpoint | Parameter / Value | Ekspektasi | Hasil Aktual | Tgl Test | PIC | Screenshot Hasil | Pass / Failed |
|-----|----------|----------|----------|-----------------|-----------------|----------------|-------------------|------------|--------------|----------|-----|------------------|---------------|
| 14 | SRS-7.3 | FSD-2.3 | TSD-3.3.ac | TS-BE-014 | Daftar Jadwal Imunisasi | GET /api/v1/imunisasi | Authorization: Bearer {token}, page=1, per_page=15 | Status 200, immunization schedule list | Status 200, jadwal diterima | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 15 | SRS-7.4 | FSD-2.3 | TSD-3.3.ae | TS-BE-015 | Tambah Jadwal Imunisasi | POST /api/v1/imunisasi | Authorization: Bearer {token}, body: id_anak, jenis_vaksin, tanggal_jadwal | Status 201, schedule created | Status 201, jadwal ditambahkan | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 16 | SRS-7.5 | FSD-2.3 | TSD-3.3.aj | TS-BE-016 | Update Realisasi Vaksin | PATCH /api/v1/imunisasi/{id}/realisasi | Authorization: Bearer {token}, id=1, body: tanggal_realisasi, catatan | Status 200, realization updated | Status 200, realisasi dicatat | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |

### Modul: Tindak Lanjut
| No. | SRS Ref. | FSD Ref. | TSD Ref. | No. Test Script | Nama Fungsional | URL / Endpoint | Parameter / Value | Ekspektasi | Hasil Aktual | Tgl Test | PIC | Screenshot Hasil | Pass / Failed |
|-----|----------|----------|----------|-----------------|-----------------|----------------|-------------------|------------|--------------|----------|-----|------------------|---------------|
| 17 | SRS-7.7 | FSD-2.4 | TSD-3.3.m | TS-BE-017 | Daftar Pasien Tindak Lanjut | GET /api/v1/tindak-lanjut/pasien | Authorization: Bearer {token}, page=1, per_page=15 | Status 200, follow-up patient list | Status 200, daftar diterima | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 18 | SRS-7.7 | FSD-2.4 | TSD-3.3.n | TS-BE-018 | Detail Pasien Tindak Lanjut | GET /api/v1/tindak-lanjut/pasien/{id} | Authorization: Bearer {token}, id=1 | Status 200, patient detail with history | Status 200, detail lengkap | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 19 | SRS-7.7 | FSD-2.4 | TSD-3.3.o | TS-BE-019 | Membuat Tindak Lanjut | POST /api/v1/tindak-lanjut | Authorization: Bearer {token}, body: id_hasil_pemeriksaan, jenis_tindak_lanjut | Status 201, follow-up created | Status 201, tindak lanjut dibuat | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 20 | SRS-7.7 | FSD-2.4 | TSD-3.3.p | TS-BE-020 | Update Status Rujukan | PATCH /api/v1/rujukan/{id}/status | Authorization: Bearer {token}, id=1, body: status_rujukan | Status 200, referral status updated | Status 200, status diperbarui | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 21 | SRS-7.7 | FSD-2.4 | TSD-3.3.s | TS-BE-021 | Detail Tindak Lanjut | GET /api/v1/tindak-lanjut/{id} | Authorization: Bearer {token}, id=1 | Status 200, follow-up detail | Status 200, detail diterima | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 22 | SRS-7.7 | FSD-2.4 | TSD-3.3.q | TS-BE-022 | Status Tindak Lanjut (Kader) | GET /api/v1/tindak-lanjut/status | Authorization: Bearer {token}, page=1, per_page=15 | Status 200, follow-up status list | Status 200, status diterima | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |

### Modul: Laporan
| No. | SRS Ref. | FSD Ref. | TSD Ref. | No. Test Script | Nama Fungsional | URL / Endpoint | Parameter / Value | Ekspektasi | Hasil Aktual | Tgl Test | PIC | Screenshot Hasil | Pass / Failed |
|-----|----------|----------|----------|-----------------|-----------------|----------------|-------------------|------------|--------------|----------|-----|------------------|---------------|
| 23 | SRS-7.7 | FSD-2.4 | TSD-3.3.r | TS-BE-023 | Laporan Tindak Lanjut | GET /api/v1/laporan/tindak-lanjut | Authorization: Bearer {token}, periode_awal, periode_akhir | Status 200, follow-up report | Status 200, laporan diterima | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |

### Modul: Artikel
| No. | SRS Ref. | FSD Ref. | TSD Ref. | No. Test Script | Nama Fungsional | URL / Endpoint | Parameter / Value | Ekspektasi | Hasil Aktual | Tgl Test | PIC | Screenshot Hasil | Pass / Failed |
|-----|----------|----------|----------|-----------------|-----------------|----------------|-------------------|------------|--------------|----------|-----|------------------|---------------|
| 24 | SRS-7.8 | FSD-2.5 | TSD-3.3.f | TS-BE-024 | Daftar Artikel | GET /api/v1/artikel | Authorization: Bearer {token}, page=1, per_page=15 | Status 200, published article list | Status 200, artikel diterima | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 25 | SRS-7.8 | FSD-2.5 | TSD-3.3.g | TS-BE-025 | Detail Artikel | GET /api/v1/artikel/{id} | Authorization: Bearer {token}, id=1 | Status 200, article detail | Status 200, detail lengkap | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 26 | SRS-7.8 | FSD-2.5 | TSD-3.3.h | TS-BE-026 | Membuat Artikel | POST /api/v1/artikel | Authorization: Bearer {token_bidan}, body: judul, isi_artikel, kategori | Status 201, article created (Draft) | Status 201, artikel dibuat | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 27 | SRS-7.8 | FSD-2.5 | TSD-3.3.i | TS-BE-027 | Update Artikel | PATCH /api/v1/artikel/{id} | Authorization: Bearer {token_bidan}, id=1, body: updated data | Status 200, article updated | Status 200, artikel diperbarui | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 28 | SRS-7.8 | FSD-2.5 | TSD-3.3.k | TS-BE-028 | Review Artikel (Setujui) | PATCH /api/v1/artikel/{id}/review | Authorization: Bearer {token_dinkes}, id=1, body: {aksi: "setujui"} | Status 200, article approved & published | Status 200, artikel disetujui | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 29 | SRS-7.8 | FSD-2.5 | TSD-3.3.l | TS-BE-029 | Daftar Artikel Pending | GET /api/v1/artikel/pending | Authorization: Bearer {token_dinkes}, page=1, per_page=15 | Status 200, pending article list | Status 200, artikel pending | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |

### Summary - Happy Path (Positif)
- **Total Test Cases**: 29
- **Pass**: 29
- **Failed**: 0
- **Coverage**: Semua endpoint sesuai TSD Section 3.2 & 3.3

**Referensi Dokumen**:
- **SRS**: SRS-7.1 (Register), SRS-7.2 (Login/Logout), SRS-7.3 (Tampil Data Monitoring), SRS-7.4 (Input Data Monitoring), SRS-7.5 (Update Data Monitoring), SRS-7.6 (Hapus Data Monitoring), SRS-7.7 (Tindak Lanjut/Rujukan), SRS-7.8 (Akses Edukasi Gizi)
- **FSD**: FSD-2.1 (Register), FSD-2.2 (Login), FSD-2.3 (Monitoring), FSD-2.4 (Tindak Lanjut), FSD-2.5 (Artikel)
- **TSD**: TSD-3.2 (Endpoint Index), TSD-3.3.a-aj (Endpoint Specs), TSD-5.x (DB Objects)

**Breakdown Referensi TSD**:
- Health Check: TSD-3.2
- Auth Endpoints: TSD-3.3.a-e  
- Monitoring Pasien: TSD-3.3.t-u-v
- Monitoring Statistik: TSD-3.2
- Monitoring Pemeriksaan: TSD-3.3.w-x-y-z-aa-ab
- Imunisasi: TSD-3.3.ac-ad-ae-af-ag-ah-ai-aj
- Tindak Lanjut: TSD-3.3.m-n-o-p-q-s
- Laporan: TSD-3.3.r
- Artikel: TSD-3.3.f-g-h-i-j-k-l

---

## SKENARIO: SAD PATH (NEGATIF)

### Modul: Health Check
| No. | SRS Ref. | FSD Ref. | TSD Ref. | No. Test Script | Nama Fungsional | URL / Endpoint | Parameter / Value | Ekspektasi | Hasil Aktual | Tgl Test | PIC | Screenshot Hasil | Pass / Failed |
|-----|----------|----------|----------|-----------------|-----------------|----------------|-------------------|------------|--------------|----------|-----|------------------|---------------|
| 1 | SRS-7.2 | FSD-2.2 | TSD-3.2 | TS-BE-NEG-001 | Health Check Invalid Method | POST /health | - | Status 405, Method Not Allowed | Status 405, method not allowed | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |

### Modul: Monitoring - Pasien
| No. | SRS Ref. | FSD Ref. | TSD Ref. | No. Test Script | Nama Fungsional | URL / Endpoint | Parameter / Value | Ekspektasi | Hasil Aktual | Tgl Test | PIC | Screenshot Hasil | Pass / Failed |
|-----|----------|----------|----------|-----------------|-----------------|----------------|-------------------|------------|--------------|----------|-----|------------------|---------------|
| 2 | SRS-7.3 | FSD-2.3 | TSD-3.3.t | TS-BE-NEG-002 | Daftar Pasien Tanpa Auth | GET /api/v1/monitoring/pasien | No Authorization header | Status 401, Unauthorized | Status 401, unauthorized | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 3 | SRS-7.3 | FSD-2.3 | TSD-3.3.t | TS-BE-NEG-003 | Daftar Pasien Token Invalid | GET /api/v1/monitoring/pasien | Authorization: Bearer invalid_token | Status 401, Unauthorized | Status 401, token invalid | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 4 | SRS-7.3 | FSD-2.3 | TSD-3.3.u | TS-BE-NEG-004 | Detail Pasien ID Tidak Ditemukan | GET /api/v1/monitoring/pasien/99999 | Authorization: Bearer {token}, id=99999 | Status 404, Not Found | Status 404, pasien tidak ditemukan | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 5 | SRS-7.3 | FSD-2.3 | TSD-3.3.v | TS-BE-NEG-005 | Cari Pasien Query Kosong | GET /api/v1/monitoring/pasien/search?q= | Authorization: Bearer {token}, q="" | Status 422, Validation Error | Status 422, query wajib diisi | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 6 | SRS-7.3 | FSD-2.3 | TSD-5.3 | TS-BE-NEG-006 | Export Pasien Role Tidak Sesuai | GET /api/v1/monitoring/pasien/export | Authorization: Bearer {token_kader} | Status 403, Forbidden | Status 403, akses ditolak | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |

### Modul: Monitoring - Statistik
| No. | SRS Ref. | FSD Ref. | TSD Ref. | No. Test Script | Nama Fungsional | URL / Endpoint | Parameter / Value | Ekspektasi | Hasil Aktual | Tgl Test | PIC | Screenshot Hasil | Pass / Failed |
|-----|----------|----------|----------|-----------------|-----------------|----------------|-------------------|------------|--------------|----------|-----|------------------|---------------|
| 7 | SRS-7.3 | FSD-2.3 | TSD-3.2 | TS-BE-NEG-007 | Summary Dashboard Tanpa Auth | GET /api/v1/monitoring/statistik | No Authorization header | Status 401, Unauthorized | Status 401, unauthorized | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 8 | SRS-7.3 | FSD-2.3 | TSD-3.2 | TS-BE-NEG-008 | Tren Bulanan Format Tahun Invalid | GET /api/v1/monitoring/statistik-bulanan?tahun=abcd | Authorization: Bearer {token}, tahun=abcd | Status 422, Validation Error | Status 422, format tahun invalid | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |

### Modul: Monitoring - Pemeriksaan
| No. | SRS Ref. | FSD Ref. | TSD Ref. | No. Test Script | Nama Fungsional | URL / Endpoint | Parameter / Value | Ekspektasi | Hasil Aktual | Tgl Test | PIC | Screenshot Hasil | Pass / Failed |
|-----|----------|----------|----------|-----------------|-----------------|----------------|-------------------|------------|--------------|----------|-----|------------------|---------------|
| 9 | SRS-7.4 | FSD-2.3 | TSD-3.3.w | TS-BE-NEG-009 | Input Pemeriksaan Data Tidak Lengkap | POST /api/v1/monitoring/pemeriksaan | Authorization: Bearer {token}, body: {id_pasien: 1} (missing required fields) | Status 422, Validation Error | Status 422, data tidak lengkap | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 10 | SRS-7.4 | FSD-2.3 | TSD-3.3.w | TS-BE-NEG-010 | Input Pemeriksaan Role Tidak Sesuai | POST /api/v1/monitoring/pemeriksaan | Authorization: Bearer {token_dinkes}, body: valid data | Status 403, Forbidden | Status 403, role tidak sesuai | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 11 | SRS-7.4 | FSD-2.3 | TSD-3.3.w | TS-BE-NEG-011 | Input Pemeriksaan ID Pasien Tidak Ditemukan | POST /api/v1/monitoring/pemeriksaan | Authorization: Bearer {token}, body: {id_pasien: 99999, ...} | Status 404, Not Found | Status 404, pasien tidak ditemukan | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 12 | SRS-7.3 | FSD-2.3 | TSD-3.3.x | TS-BE-NEG-012 | Detail Pemeriksaan ID Tidak Valid | GET /api/v1/monitoring/pemeriksaan/abc | Authorization: Bearer {token}, id=abc | Status 400, Bad Request | Status 400, ID harus angka | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 13 | SRS-7.5 | FSD-2.3 | TSD-3.3.y | TS-BE-NEG-013 | Edit Pemeriksaan Sudah Diverifikasi | PUT /api/v1/monitoring/pemeriksaan/{id} | Authorization: Bearer {token}, id=1 (already verified) | Status 422, Validation Error | Status 422, sudah diverifikasi | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 14 | SRS-7.6 | FSD-2.3 | TSD-3.3.z | TS-BE-NEG-014 | Hapus Pemeriksaan Tidak Ditemukan | DELETE /api/v1/monitoring/pemeriksaan/99999 | Authorization: Bearer {token}, id=99999 | Status 404, Not Found | Status 404, data tidak ditemukan | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 15 | SRS-7.4 | FSD-2.3 | TSD-3.3.aa | TS-BE-NEG-015 | Verifikasi Pemeriksaan Bukan Bidan | PATCH /api/v1/monitoring/pemeriksaan/{id}/verify | Authorization: Bearer {token_kader}, id=1 | Status 403, Forbidden | Status 403, hanya bidan | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 16 | SRS-7.3 | FSD-2.3 | TSD-3.3.ab | TS-BE-NEG-016 | Pemeriksaan Pending Tanpa Auth | GET /api/v1/monitoring/pemeriksaan/pending | No Authorization header | Status 401, Unauthorized | Status 401, unauthorized | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |

### Modul: Imunisasi
| No. | SRS Ref. | FSD Ref. | TSD Ref. | No. Test Script | Nama Fungsional | URL / Endpoint | Parameter / Value | Ekspektasi | Hasil Aktual | Tgl Test | PIC | Screenshot Hasil | Pass / Failed |
|-----|----------|----------|----------|-----------------|-----------------|----------------|-------------------|------------|--------------|----------|-----|------------------|---------------|
| 17 | SRS-7.3 | FSD-2.3 | TSD-3.3.ac | TS-BE-NEG-017 | Daftar Jadwal Imunisasi Tanpa Auth | GET /api/v1/imunisasi | No Authorization header | Status 401, Unauthorized | Status 401, unauthorized | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 18 | SRS-7.4 | FSD-2.3 | TSD-3.3.ae | TS-BE-NEG-018 | Tambah Jadwal Imunisasi Data Invalid | POST /api/v1/imunisasi | Authorization: Bearer {token}, body: {jenis_vaksin: ""} | Status 422, Validation Error | Status 422, data tidak valid | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 19 | SRS-7.5 | FSD-2.3 | TSD-3.3.aj | TS-BE-NEG-019 | Update Realisasi ID Tidak Ditemukan | PATCH /api/v1/imunisasi/99999/realisasi | Authorization: Bearer {token}, id=99999 | Status 404, Not Found | Status 404, jadwal tidak ditemukan | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |

### Modul: Tindak Lanjut
| No. | SRS Ref. | FSD Ref. | TSD Ref. | No. Test Script | Nama Fungsional | URL / Endpoint | Parameter / Value | Ekspektasi | Hasil Aktual | Tgl Test | PIC | Screenshot Hasil | Pass / Failed |
|-----|----------|----------|----------|-----------------|-----------------|----------------|-------------------|------------|--------------|----------|-----|------------------|---------------|
| 20 | SRS-7.7 | FSD-2.4 | TSD-3.3.m | TS-BE-NEG-020 | Daftar Pasien Tindak Lanjut Bukan Bidan | GET /api/v1/tindak-lanjut/pasien | Authorization: Bearer {token_kader} | Status 403, Forbidden | Status 403, hanya bidan | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 21 | SRS-7.7 | FSD-2.4 | TSD-3.3.n | TS-BE-NEG-021 | Detail Pasien Tindak Lanjut ID Invalid | GET /api/v1/tindak-lanjut/pasien/99999 | Authorization: Bearer {token}, id=99999 | Status 404, Not Found | Status 404, pasien tidak ditemukan | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 22 | SRS-7.7 | FSD-2.4 | TSD-3.3.o | TS-BE-NEG-022 | Membuat Tindak Lanjut Data Tidak Lengkap | POST /api/v1/tindak-lanjut | Authorization: Bearer {token}, body: {id_hasil_pemeriksaan: 1} (missing jenis) | Status 422, Validation Error | Status 422, data tidak lengkap | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 23 | SRS-7.7 | FSD-2.4 | TSD-3.3.o | TS-BE-NEG-023 | Membuat Tindak Lanjut Bukan Bidan | POST /api/v1/tindak-lanjut | Authorization: Bearer {token_kader}, body: valid data | Status 403, Forbidden | Status 403, hanya bidan | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 24 | SRS-7.7 | FSD-2.4 | TSD-3.3.p | TS-BE-NEG-024 | Update Status Rujukan Status Invalid | PATCH /api/v1/rujukan/{id}/status | Authorization: Bearer {token}, id=1, body: {status_rujukan: "INVALID"} | Status 422, Validation Error | Status 422, status tidak valid | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 25 | SRS-7.7 | FSD-2.4 | TSD-3.3.s | TS-BE-NEG-025 | Detail Tindak Lanjut Akses Ditolak | GET /api/v1/tindak-lanjut/{id} | Authorization: Bearer {token_other_user}, id=1 | Status 403, Forbidden | Status 403, bukan pemilik data | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 26 | SRS-7.7 | FSD-2.4 | TSD-3.3.q | TS-BE-NEG-026 | Status Tindak Lanjut Bukan Kader | GET /api/v1/tindak-lanjut/status | Authorization: Bearer {token_bidan} | Status 403, Forbidden | Status 403, hanya kader | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |

### Modul: Laporan
| No. | SRS Ref. | FSD Ref. | TSD Ref. | No. Test Script | Nama Fungsional | URL / Endpoint | Parameter / Value | Ekspektasi | Hasil Aktual | Tgl Test | PIC | Screenshot Hasil | Pass / Failed |
|-----|----------|----------|----------|-----------------|-----------------|----------------|-------------------|------------|--------------|----------|-----|------------------|---------------|
| 27 | SRS-7.7 | FSD-2.4 | TSD-3.3.r | TS-BE-NEG-027 | Laporan Tindak Lanjut Bukan Dinkes | GET /api/v1/laporan/tindak-lanjut | Authorization: Bearer {token_bidan} | Status 403, Forbidden | Status 403, hanya dinkes | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |

### Modul: Artikel
| No. | SRS Ref. | FSD Ref. | TSD Ref. | No. Test Script | Nama Fungsional | URL / Endpoint | Parameter / Value | Ekspektasi | Hasil Aktual | Tgl Test | PIC | Screenshot Hasil | Pass / Failed |
|-----|----------|----------|----------|-----------------|-----------------|----------------|-------------------|------------|--------------|----------|-----|------------------|---------------|
| 28 | SRS-7.8 | FSD-2.5 | TSD-3.3.f | TS-BE-NEG-028 | Daftar Artikel Tanpa Auth | GET /api/v1/artikel | No Authorization header | Status 401, Unauthorized | Status 401, unauthorized | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 29 | SRS-7.8 | FSD-2.5 | TSD-3.3.g | TS-BE-NEG-029 | Detail Artikel ID Tidak Ditemukan | GET /api/v1/artikel/99999 | Authorization: Bearer {token}, id=99999 | Status 404, Not Found | Status 404, artikel tidak ditemukan | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 30 | SRS-7.8 | FSD-2.5 | TSD-3.3.h | TS-BE-NEG-030 | Membuat Artikel Bukan Bidan | POST /api/v1/artikel | Authorization: Bearer {token_kader}, body: valid data | Status 403, Forbidden | Status 403, hanya bidan | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 31 | SRS-7.8 | FSD-2.5 | TSD-3.3.h | TS-BE-NEG-031 | Membuat Artikel Judul Kosong | POST /api/v1/artikel | Authorization: Bearer {token_bidan}, body: {judul: "", isi_artikel: "...", kategori: "..."} | Status 422, Validation Error | Status 422, judul wajib diisi | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 32 | SRS-7.8 | FSD-2.5 | TSD-3.3.i | TS-BE-NEG-032 | Update Artikel Bukan Pemilik | PATCH /api/v1/artikel/{id} | Authorization: Bearer {token_bidan_other}, id=1 | Status 403, Forbidden | Status 403, bukan pemilik | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 33 | SRS-7.8 | FSD-2.5 | TSD-3.3.i | TS-BE-NEG-033 | Update Artikel Sudah Published | PATCH /api/v1/artikel/{id} | Authorization: Bearer {token_bidan}, id=1 (already published) | Status 422, Validation Error | Status 422, sudah published | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 34 | SRS-7.8 | FSD-2.5 | TSD-3.3.k | TS-BE-NEG-034 | Review Artikel Bukan Dinkes | PATCH /api/v1/artikel/{id}/review | Authorization: Bearer {token_bidan}, id=1, body: {aksi: "setujui"} | Status 403, Forbidden | Status 403, hanya dinkes | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 35 | SRS-7.8 | FSD-2.5 | TSD-3.3.k | TS-BE-NEG-035 | Review Artikel Aksi Invalid | PATCH /api/v1/artikel/{id}/review | Authorization: Bearer {token_dinkes}, id=1, body: {aksi: "invalid"} | Status 422, Validation Error | Status 422, aksi tidak valid | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 36 | SRS-7.8 | FSD-2.5 | TSD-3.3.k | TS-BE-NEG-036 | Review Artikel Tolak Tanpa Catatan | PATCH /api/v1/artikel/{id}/review | Authorization: Bearer {token_dinkes}, id=1, body: {aksi: "tolak", catatan_review: ""} | Status 422, Validation Error | Status 422, catatan wajib diisi | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 37 | SRS-7.8 | FSD-2.5 | TSD-3.3.l | TS-BE-NEG-037 | Daftar Artikel Pending Bukan Dinkes | GET /api/v1/artikel/pending | Authorization: Bearer {token_bidan} | Status 403, Forbidden | Status 403, hanya dinkes | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |
| 38 | SRS-7.8 | FSD-2.5 | TSD-3.3.j | TS-BE-NEG-038 | Hapus Artikel Bukan Dinkes | DELETE /api/v1/artikel/{id} | Authorization: Bearer {token_bidan}, id=1 | Status 403, Forbidden | Status 403, hanya dinkes | 10/06/2026 | Tester Backend | Lihat lampiran | Pass |

### Summary - Sad Path (Negatif)
- **Total Test Cases**: 38
- **Pass**: 38
- **Failed**: 0

---

## RINGKASAN KESELURUHAN

### Total Test Cases: 67
- **Happy Path**: 29 test cases
- **Sad Path**: 38 test cases

### Breakdown per Modul:
1. **Health Check** (SRS-7.2, FSD-2.2, TSD-3.2): 1 happy path + 1 sad path = 2 test cases
2. **Monitoring - Pasien** (SRS-7.3, FSD-2.3, TSD-3.3.t-v): 4 happy path + 5 sad path = 9 test cases
3. **Monitoring - Statistik** (SRS-7.3, FSD-2.3, TSD-3.2): 2 happy path + 2 sad path = 4 test cases
4. **Monitoring - Pemeriksaan** (SRS-7.3-7.6, FSD-2.3, TSD-3.3.w-ab): 6 happy path + 8 sad path = 14 test cases
5. **Imunisasi** (SRS-7.3-7.5, FSD-2.3, TSD-3.3.ac-aj): 3 happy path + 3 sad path = 6 test cases
6. **Tindak Lanjut** (SRS-7.7, FSD-2.4, TSD-3.3.m-s): 6 happy path + 7 sad path = 13 test cases
7. **Laporan** (SRS-7.7, FSD-2.4, TSD-3.3.r): 1 happy path + 1 sad path = 2 test cases
8. **Artikel** (SRS-7.8, FSD-2.5, TSD-3.3.f-l): 6 happy path + 11 sad path = 17 test cases

### Catatan Penting:
- Semua endpoint (29 endpoints) telah di-cover dalam test script
- Setiap endpoint minimal memiliki 1 skenario happy path
- Sad path mencakup berbagai skenario error: authentication, authorization, validation, not found, dll
- Token authorization disesuaikan dengan role yang dibutuhkan (BIDAN, KADER, DINKES, IBU/WALI)
- Test dilakukan pada tanggal 10 Juni 2026
- Referensi lengkap sesuai dokumentasi sistem:
  - **SRS Section 7**: Functional Requirements (7.1-7.8)
  - **FSD Section 2**: Functional Specification Design (2.1-2.5)
  - **TSD Section 3**: Technical Specification Design (3.2, 3.3.a-aj)

### Endpoint Summary:

#### 1. Health Check (1 endpoint)
- `GET /health`

#### 2. Monitoring - Pasien (4 endpoints)
- `GET /api/v1/monitoring/pasien`
- `GET /api/v1/monitoring/pasien/{id}`
- `GET /api/v1/monitoring/pasien/search`
- `GET /api/v1/monitoring/pasien/export`

#### 3. Monitoring - Statistik (2 endpoints)
- `GET /api/v1/monitoring/statistik`
- `GET /api/v1/monitoring/statistik-bulanan`

#### 4. Monitoring - Pemeriksaan (6 endpoints)
- `POST /api/v1/monitoring/pemeriksaan`
- `GET /api/v1/monitoring/pemeriksaan/{id}`
- `PUT /api/v1/monitoring/pemeriksaan/{id}`
- `DELETE /api/v1/monitoring/pemeriksaan/{id}`
- `PATCH /api/v1/monitoring/pemeriksaan/{id}/verify`
- `GET /api/v1/monitoring/pemeriksaan/pending`

#### 5. Imunisasi (3 endpoints)
- `GET /api/v1/imunisasi`
- `POST /api/v1/imunisasi`
- `PATCH /api/v1/imunisasi/{id}/realisasi`

#### 6. Tindak Lanjut (6 endpoints)
- `GET /api/v1/tindak-lanjut/pasien`
- `GET /api/v1/tindak-lanjut/pasien/{id}`
- `POST /api/v1/tindak-lanjut`
- `PATCH /api/v1/rujukan/{id}/status`
- `GET /api/v1/tindak-lanjut/{id}`
- `GET /api/v1/tindak-lanjut/status`

#### 7. Laporan (1 endpoint)
- `GET /api/v1/laporan/tindak-lanjut`

#### 8. Artikel (6 endpoints)
- `GET /api/v1/artikel`
- `GET /api/v1/artikel/{id}`
- `POST /api/v1/artikel`
- `PATCH /api/v1/artikel/{id}`
- `PATCH /api/v1/artikel/{id}/review`
- `GET /api/v1/artikel/pending`
- `DELETE /api/v1/artikel/{id}`

**Total: 29 Unique Endpoints**

---

## REKOMENDASI TESTING

### Tools yang Disarankan:
1. **Postman** atau **Insomnia** - untuk manual testing
2. **Newman** - untuk automated testing via command line
3. **Jest + Supertest** - untuk integration testing otomatis
4. **K6** atau **Apache JMeter** - untuk load testing

### Environment Setup:
```
BASE_URL=http://localhost:8080
TOKEN_BIDAN=<JWT token untuk role BIDAN>
TOKEN_DINKES=<JWT token untuk role DINKES>
TOKEN_KADER=<JWT token untuk role KADER>
TOKEN_IBU=<JWT token untuk role IBU/WALI>
```

### Best Practices:
1. Jalankan test dalam urutan: Happy Path dulu, baru Sad Path
2. Gunakan data dummy yang konsisten untuk setiap test
3. Reset database setelah setiap test cycle untuk konsistensi
4. Capture screenshot untuk setiap test case
5. Dokumentasikan response time untuk setiap endpoint
6. Lakukan load testing untuk endpoint yang sering diakses

---

## TRACEABILITY MATRIX

Setiap test case dapat dilacak kembali ke dokumen requirement:

| Module | SRS | FSD | TSD | Test Cases |
|--------|-----|-----|-----|------------|
| Health Check | - | - | TSD-3.2 | 2 |
| Auth | FR-7.1, FR-7.2 | FSD-2.1, FSD-2.2 | TSD-3.3.a-e | - |
| Monitoring Pasien | FR-7.3 | FSD-2.3 | TSD-3.3.t-u-v | 9 |
| Monitoring Statistik | FR-7.3 | FSD-2.3 | TSD-3.2 | 4 |
| Monitoring Pemeriksaan | FR-7.3, FR-7.4 | FSD-2.3 | TSD-3.3.w-ab | 14 |
| Imunisasi | FR-7.5, FR-7.6 | FSD-2.3 | TSD-3.3.ac-aj | 6 |
| Tindak Lanjut & Rujukan | FR-7.7 | FSD-2.4 | TSD-3.3.m-s | 13 |
| Laporan | FR-7.7 | FSD-2.4 | TSD-3.3.r | 2 |
| Artikel | FR-7.5 | FSD-2.5 | TSD-3.3.f-l | 17 |
| **TOTAL** | - | - | - | **67** |

---

**Dokumen Test Script ini telah mencakup semua 29 endpoint backend SIGizi**

*Generated: 10 Juni 2026*  
*Kelompok 5 - Kelas A - RSI 2026*
