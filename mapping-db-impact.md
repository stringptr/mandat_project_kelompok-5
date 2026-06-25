# Mapping DB Impact — SiGizi

> Dampak perubahan objek database terhadap API dan fungsionalitas  
> Database: PostgreSQL 18 (Master DB) | ORM: Go-Jet v2 | Total: 21 tabel, 0 SP, 1 function, 11 MV, 39 index, 15 trigger

---

## Daftar Isi

1. [Tabel](#tabel-impact-mapping)
2. [Stored Procedure](#stored-procedure)
3. [Function](#function)
4. [View & Materialized View](#view--materialized-view)
5. [Index](#index)
6. [Trigger](#trigger)
7. [Impact Matrix (Ringkasan)](#impact-matrix-ringkasan)

---

## Tabel Impact Mapping

### 1. `user_account`
| Aspek | Detail |
|--------|--------|
| **Deskripsi** | Tabel utama seluruh pengguna (Bidan, Kader, Dinkes, Pasien, Ibu Hamil, Anak, SuperAdmin) |
| **API READ** | POST /auth/login, POST /auth/refresh, POST /auth/register, GET /auth/me, GET /users, GET /users/{id}, GET /monitoring/pasien, GET /monitoring/pasien/{id}, GET /monitoring/pasien/search, GET /monitoring/pemeriksaan, GET /imunisasi, GET /imunisasi/{id}, GET /monitoring/semua-pemeriksaan, GET /artikel, GET /artikel/{id}, GET /artikel/semua, GET /artikel/pending, GET /tindak-lanjut/status, GET /tindak-lanjut/pasien, GET /tindak-lanjut/pasien/{id}, GET /admin/audit-logs |
| **API WRITE** | POST /auth/register (INSERT), PATCH /users/{id} (UPDATE), POST /users (INSERT), PATCH /users/{id_user}/verification (UPDATE), PATCH /users/{id}/role (indirect via role tables), DELETE /users/{id} (soft-delete), POST /auth/logout (indirect via session) |
| **DML** | SELECT, INSERT, UPDATE (soft-delete via is_deleted) |
| **Fungsionalitas** | Registrasi, login, manajemen profil, manajemen user, pencarian pasien, penulisan artikel |
| **Dampak Perubahan** | **KRITIS** — seluruh sistem bergantung. Perubahan skema (tambah/hapus kolom) berdampak pada semua endpoint auth, user management, pasien, notifikasi. Soft-delete via `is_deleted` harus selalu difilter. |

### 2. `lokasi`
| Aspek | Detail |
|--------|--------|
| **Deskripsi** | Data hierarki wilayah (Provinsi → Kabupaten → Kota → Kecamatan → Kelurahan) |
| **API READ** | GET /lokasi, GET /users (GetLokasiNames via raw query), GET /laporan/tindak-lanjut (via posyandu JOIN) |
| **API WRITE** | — (read-only di aplikasi) |
| **DML** | SELECT |
| **Fungsionalitas** | Dropdown lokasi bertingkat, filter wilayah di dashboard & laporan |
| **Dampak Perubahan** | **SEDANG** — perubahan data (rename/restruktur) berdampak pada dropdown lokasi dan laporan per wilayah. Self-referencing FK (`bagian_dari`) harus dijaga integritasnya. |

### 3. `dinas_kesehatan`
| Aspek | Detail |
|--------|--------|
| **Deskripsi** | Data Dinas Kesehatan (extends user_account via FK `id_user`) |
| **API READ** | POST /auth/login (EXISTS check role), GET /users, GET /users/{id} (role detection) |
| **API WRITE** | POST /auth/register (INSERT jika role=Dinkes), PATCH /users/{id}/role (soft-delete lama → INSERT baru) |
| **DML** | SELECT (EXISTS), INSERT, UPDATE (soft-delete) |
| **Fungsionalitas** | Role assignment untuk Dinas Kesehatan |
| **Dampak Perubahan** | **RENDAH** — hanya untuk deteksi role. Cascade delete dari `user_account` (`ON DELETE CASCADE`). |

### 4. `bidan`
| Aspek | Detail |
|--------|--------|
| **Deskripsi** | Data Bidan (extends user_account, FK ke `lokasi` untuk wilayah_kerja) |
| **API READ** | POST /auth/login (EXISTS check role), GET /users, GET /users/{id} (role detection), GET /notifikasi/bidan (dashboard bidan), GET /monitoring/semua-pemeriksaan (filter by id_bidan), POST /tindak-lanjut (id_bidan dari JWT) |
| **API WRITE** | POST /auth/register (INSERT jika role=Bidan), PATCH /users/{id}/role (soft-delete + INSERT) |
| **DML** | SELECT (EXISTS), INSERT, UPDATE (soft-delete) |
| **Fungsionalitas** | Role Bidan, dashboard notifikasi bidan, filter pemeriksaan per bidan, tindak lanjut oleh bidan |
| **Dampak Perubahan** | **SEDANG** — perubahan berdampak pada dashboard bidan, filter laporan, dan assignment tindak lanjut. |

### 5. `posyandu`
| Aspek | Detail |
|--------|--------|
| **Deskripsi** | Data Posyandu (FK ke `lokasi` dan `bidan`) |
| **API READ** | GET /monitoring/pasien, GET /monitoring/pasien/{id}, GET /monitoring/pasien/search, GET /notifikasi/bidan, GET /laporan/tindak-lanjut, GET /monitoring/semua-pemeriksaan |
| **API WRITE** | — (read-only di aplikasi) |
| **DML** | SELECT |
| **Fungsionalitas** | Pendaftaran pasien ke posyandu, filter laporan per posyandu, dashboard bidan |
| **Dampak Perubahan** | **SEDANG** — struktur posyandu memengaruhi pengelompokan pasien dan laporan wilayah. |

### 6. `kader_posyandu`
| Aspek | Detail |
|--------|--------|
| **Deskripsi** | Data Kader Posyandu (extends user_account, FK ke `posyandu`) |
| **API READ** | POST /auth/login (EXISTS check role), GET /users, GET /users/{id} (role detection), GET /monitoring/semua-pemeriksaan (get id_posyandu dari kader) |
| **API WRITE** | POST /auth/register (INSERT jika role=Kader), PATCH /users/{id}/role (soft-delete + INSERT) |
| **DML** | SELECT (EXISTS), INSERT, UPDATE (soft-delete) |
| **Fungsionalitas** | Role Kader, penentuan posyandu kader untuk filter pemeriksaan |
| **Dampak Perubahan** | **RENDAH** — serupa dengan bidan, lebih terbatas cakupannya. |

### 7. `pasien`
| Aspek | Detail |
|--------|--------|
| **Deskripsi** | Data Pasien (extends user_account, FK ke `posyandu`) — parent entity untuk Ibu Hamil & Anak |
| **API READ** | POST /auth/login (EXISTS check role), GET /monitoring/pasien, GET /monitoring/pasien/{id}, GET /monitoring/pasien/search, GET /notifikasi/bidan, GET /tindak-lanjut/status, GET /tindak-lanjut/pasien, GET /tindak-lanjut/pasien/{id}, GET /laporan/tindak-lanjut, GET /monitoring/semua-pemeriksaan, GET /imunisasi, GET /imunisasi/{id}, GET /imunisasi/pasien/{id_pasien} |
| **API WRITE** | POST /pasien/ibu-hamil (INSERT), POST /pasien/anak (INSERT), PATCH /pasien/{id} (UPDATE), DELETE /pasien/{id} (soft-delete dalam TRANSACTION bersama anak & ibu_hamil) |
| **DML** | SELECT (JOIN), INSERT, UPDATE, UPDATE (soft-delete) |
| **Fungsionalitas** | Pendaftaran & manajemen pasien, pencarian, monitoring, imunisasi, tindak lanjut |
| **Dampak Perubahan** | **KRITIS** — tabel sentral data pasien. Semua fitur monitoring, imunisasi, dan tindak lanjut bergantung. Cascade delete ke `ibu_hamil` dan `anak`. |

### 8. `fasilitas_kesehatan`
| Aspek | Detail |
|--------|--------|
| **Deskripsi** | Data Fasilitas Kesehatan (FK ke `lokasi`) — Puskesmas, RS, Klinik |
| **API READ** | GET /faskes, GET /tindak-lanjut/status (JOIN via rujukan), GET /tindak-lanjut/{id} (LEFT JOIN rujukan) |
| **API WRITE** | — (read-only) |
| **DML** | SELECT |
| **Fungsionalitas** | Referensi faskes untuk dropdown, rujukan pasien ke faskes |
| **Dampak Perubahan** | **RENDAH** — tabel referensi. Perubahan data berdampak pada dropdown dan laporan rujukan. |

### 9. `pendidikan`
| Aspek | Detail |
|--------|--------|
| **Deskripsi** | Referensi tingkat pendidikan |
| **API READ** | — (hanya sebagai FK di user_account) |
| **API WRITE** | — |
| **DML** | — (tidak ada query langsung dari repository) |
| **Fungsionalitas** | Data demografi pengguna |
| **Dampak Perubahan** | **MINIMAL** — hanya referensi statis. Tidak ada endpoint khusus. |

### 10. `pekerjaan`
| Aspek | Detail |
|--------|--------|
| **Deskripsi** | Referensi jenis pekerjaan |
| **API READ** | — (FK di user_account) |
| **API WRITE** | — |
| **DML** | — |
| **Fungsionalitas** | Data demografi pengguna |
| **Dampak Perubahan** | **MINIMAL** — referensi statis. |

### 11. `kategori_pendapatan`
| Aspek | Detail |
|--------|--------|
| **Deskripsi** | Referensi kategori pendapatan |
| **API READ** | — (FK di user_account) |
| **API WRITE** | — |
| **DML** | — |
| **Fungsionalitas** | Data demografi pengguna |
| **Dampak Perubahan** | **MINIMAL** — referensi statis. |

### 12. `ibu_hamil`
| Aspek | Detail |
|--------|--------|
| **Deskripsi** | Data kehamilan (FK ke `pasien`, cascade delete) |
| **API READ** | POST /auth/login (EXISTS check role), GET /monitoring/pasien, GET /monitoring/pasien/{id} (LEFT JOIN LATERAL), GET /users, GET /users/{id} (role detection) |
| **API WRITE** | POST /pasien/ibu-hamil (INSERT), PATCH /pasien/{id} (UPDATE), DELETE /pasien/{id} (soft-delete dalam TRANSACTION) |
| **DML** | SELECT, INSERT, UPDATE, UPDATE (soft-delete) |
| **Fungsionalitas** | Pendaftaran ibu hamil, tracking status kehamilan, dashboard ibu hamil |
| **Dampak Perubahan** | **SEDANG** — perubahan kolom status_kehamilan berdampak pada dashboard ibu hamil dan filter pasien. |

### 13. `anak`
| Aspek | Detail |
|--------|--------|
| **Deskripsi** | Data anak (FK ke `pasien`, `ibu_hamil`, `wali`/user_account) |
| **API READ** | POST /auth/login (EXISTS check role), GET /monitoring/pasien, GET /monitoring/pasien/{id} (LEFT JOIN LATERAL), GET /users, GET /users/{id} (role detection), GET /imunisasi, GET /imunisasi/{id} (ownership check), GET /notifikasi/bidan, GET /tindak-lanjut/pasien, GET /tindak-lanjut/pasien/{id}, GET /monitoring/pemeriksaan (ownership check), GET /monitoring/semua-pemeriksaan |
| **API WRITE** | POST /pasien/anak (INSERT), PATCH /pasien/{id} (UPDATE), DELETE /pasien/{id} (soft-delete dalam TRANSACTION) |
| **DML** | SELECT, INSERT, UPDATE, UPDATE (soft-delete) |
| **Fungsionalitas** | Pendaftaran anak, relasi wali-anak, ownership check untuk akses data pasien |
| **Dampak Perubahan** | **KRITIS** — `id_wali` digunakan sebagai ownership check di banyak endpoint pasien, imunisasi, dan pemeriksaan. Perubahan struktur dapat merusak otorisasi akses. |

### 14. `jadwal_imunisasi`
| Aspek | Detail |
|--------|--------|
| **Deskripsi** | Jadwal imunisasi per pasien (FK ke `pasien`, cascade delete) |
| **API READ** | GET /imunisasi, GET /imunisasi/{id}, GET /imunisasi/pasien/{id_pasien}, GET /imunisasi/statistik, GET /monitoring/pemeriksaan, GET /monitoring/pemeriksaan/{id}, GET /notifikasi/bidan, GET /dashboard/jadwal-terdekat, GET /tindak-lanjut/status, GET /tindak-lanjut/pasien, GET /tindak-lanjut/pasien/{id}, GET /laporan/tindak-lanjut, GET /monitoring/semua-pemeriksaan |
| **API WRITE** | POST /imunisasi (INSERT), PUT /imunisasi/{id} (UPDATE), DELETE /imunisasi/{id} (HARD DELETE), PATCH /imunisasi/{id}/realisasi (UPDATE) |
| **DML** | SELECT, INSERT, UPDATE, DELETE (hard) |
| **Fungsionalitas** | Penjadwalan imunisasi, realisasi, statistik cakupan, dasar pembuatan pemeriksaan |
| **Dampak Perubahan** | **KRITIS** — tabel jembatan antara imunisasi dan pemeriksaan (FK di `hasil_pemeriksaan`). HARD DELETE — tidak bisa direcover. Cascade dari pasien. |

### 15. `artikel`
| Aspek | Detail |
|--------|--------|
| **Deskripsi** | Artikel edukasi kesehatan (FK ke `user_account` untuk penulis & verifikator) |
| **API READ** | GET /artikel, GET /artikel/{id}, GET /artikel/semua, GET /artikel/pending, GET /stats (via mv_public_stats) |
| **API WRITE** | POST /artikel (INSERT), PATCH /artikel/{id} (UPDATE), DELETE /artikel/{id} (HARD DELETE), PATCH /artikel/{id}/review (UPDATE) |
| **DML** | SELECT, INSERT, UPDATE, DELETE (hard) |
| **Fungsionalitas** | Manajemen konten edukasi, workflow verifikasi (Draft → Menunggu Verifikasi → Dipublikasikan/Ditolak) |
| **Dampak Perubahan** | **SEDANG** — HARD DELETE. Perubahan enum `status_artikel` berdampak pada workflow verifikasi. |

### 16. `hasil_pemeriksaan`
| Aspek | Detail |
|--------|--------|
| **Deskripsi** | Hasil pemeriksaan antropometri & status gizi (FK ke `jadwal_imunisasi` dan `user_account`) |
| **API READ** | GET /monitoring/pemeriksaan, GET /monitoring/pemeriksaan/{id}, GET /monitoring/pemeriksaan/pending, GET /notifikasi/bidan, GET /tindak-lanjut/pasien, GET /tindak-lanjut/pasien/{id} (LATERAL subquery), GET /tindak-lanjut/status, GET /laporan/tindak-lanjut, GET /monitoring/semua-pemeriksaan, GET /monitoring/pasien/{id}/riwayat-pemeriksaan (via MV), GET /monitoring/pasien/{id}/tumbuh-kembang (via MV), GET /dashboard/stats (via MV), GET /dashboard/distribusi-gizi (via MV), GET /dashboard/tren-stunting (via MV) |
| **API WRITE** | POST /monitoring/pemeriksaan (INSERT), PUT /monitoring/pemeriksaan/{id} (UPDATE), DELETE /monitoring/pemeriksaan/{id} (HARD DELETE), PATCH /monitoring/pemeriksaan/{id}/verify (UPDATE) |
| **DML** | SELECT, INSERT, UPDATE, DELETE (hard) |
| **Fungsionalitas** | Pencatatan hasil pemeriksaan, penentuan status gizi & stunting, verifikasi bidan, sumber data dashboard |
| **Dampak Perubahan** | **KRITIS** — HARD DELETE. Data sumber untuk 8 materialized view dashboard. Status gizi & stunting digunakan di tindak lanjut. |

### 17. `tindak_lanjut`
| Aspek | Detail |
|--------|--------|
| **Deskripsi** | Tindak lanjut hasil pemeriksaan (FK ke `hasil_pemeriksaan` cascade, `bidan`) |
| **API READ** | GET /tindak-lanjut/status, GET /tindak-lanjut/pasien, GET /tindak-lanjut/{id}, GET /laporan/tindak-lanjut, GET /notifikasi/bidan |
| **API WRITE** | POST /tindak-lanjut (INSERT), PATCH /rujukan/{id}/status (indirect via rujukan) |
| **DML** | SELECT, INSERT |
| **Fungsionalitas** | Pencatatan tindak lanjut (rujukan/kontrol ulang), tracking status pasien |
| **Dampak Perubahan** | **SEDANG** — cascade delete dari `hasil_pemeriksaan`. Terkait erat dengan `rujukan`. |

### 18. `rujukan`
| Aspek | Detail |
|--------|--------|
| **Deskripsi** | Data rujukan pasien (FK ke `tindak_lanjut` cascade, `fasilitas_kesehatan`) |
| **API READ** | GET /tindak-lanjut/status, GET /tindak-lanjut/{id}, GET /laporan/tindak-lanjut, GET /notifikasi/bidan |
| **API WRITE** | POST /tindak-lanjut (INSERT jika jenis=Rujukan), PATCH /rujukan/{id}/status (UPDATE) |
| **DML** | SELECT, INSERT, UPDATE |
| **Fungsionalitas** | Manajemen rujukan, tracking status (Diajukan → Diproses → Diterima/Ditolak → Selesai), laporan rujukan per wilayah |
| **Dampak Perubahan** | **SEDANG** — cascade dari `tindak_lanjut`. Enum `status_rujukan` memengaruhi workflow rujukan. |

### 19. `notifikasi`
| Aspek | Detail |
|--------|--------|
| **Deskripsi** | Notifikasi pengguna (FK ke `user_account` cascade) |
| **API READ** | GET /notifikasi, GET /notifikasi/{id}, GET /notifikasi/aktivitas, GET /notifikasi/bidan (via notifikasi + table JOIN) |
| **API WRITE** | PATCH /notifikasi/{id}/read (UPDATE), PATCH /notifikasi/read-all (UPDATE) |
| **DML** | SELECT, INSERT, UPDATE |
| **Fungsionalitas** | Sistem notifikasi real-time (via NATS), notifikasi bidan dashboard, aktivitas pengguna |
| **Dampak Perubahan** | **RENDAH** — INSERT dilakukan via NATS publisher (bukan langsung API). Notifikasi dibaca oleh end-user. |

### 20. `user_session`
| Aspek | Detail |
|--------|--------|
| **Deskripsi** | Session management untuk JWT refresh token (FK ke `user_account`) |
| **API READ** | POST /auth/refresh (validasi session) |
| **API WRITE** | POST /auth/login (INSERT), POST /auth/refresh (UPDATE rotate), POST /auth/logout (UPDATE cabut) |
| **DML** | SELECT, INSERT, UPDATE |
| **Fungsionalitas** | Session management, refresh token rotation, logout |
| **Dampak Perubahan** | **KRITIS** — perubahan struktur session berdampak pada seluruh flow autentikasi. UUID primary key. |

### 21. `audit_log`
| Aspek | Detail |
|--------|--------|
| **Deskripsi** | Log aktivitas untuk audit trail (FK ke `user_account` dan `user_session`) |
| **API READ** | GET /admin/audit-logs |
| **API WRITE** | POST /auth/login (INSERT), POST /auth/register (INSERT), berbagai aktivitas sistem (INSERT via auditlog repo) |
| **DML** | SELECT, INSERT |
| **Fungsionalitas** | Audit trail semua aktivitas (login, register, CRUD, error, banned IP) |
| **Dampak Perubahan** | **RENDAH** — append-only. Tidak memengaruhi fungsionalitas bisnis utama. |

---

## Stored Procedure

| Status |
|--------|
| **Tidak ada stored procedure** dalam database. Semua logika bisnis diimplementasikan di application layer (Go) menggunakan Go-Jet ORM dan raw SQL query. |

---

## Function

### `update_updated_at_column()`
| Aspek | Detail |
|--------|--------|
| **Tipe** | TRIGGER FUNCTION (plpgsql) |
| **Return** | TRIGGER |
| **Logic** | `NEW.updated_at = CURRENT_TIMESTAMP; RETURN NEW;` |
| **Digunakan Oleh** | 15 trigger BEFORE UPDATE |
| **API Terkait** | Semua endpoint yang melakukan UPDATE: PATCH /users/{id}, PUT /imunisasi/{id}, PATCH /imunisasi/{id}/realisasi, PUT /monitoring/pemeriksaan/{id}, PATCH /monitoring/pemeriksaan/{id}/verify, PATCH /artikel/{id}, PATCH /artikel/{id}/review, PATCH /pasien/{id}, PATCH /rujukan/{id}/status, POST /auth/refresh, POST /auth/logout, dll. |
| **Fungsionalitas** | Auto-set `updated_at` setiap kali row di-update |
| **Dampak Perubahan** | **KRITIS** — 15 trigger bergantung pada fungsi ini. Jika fungsi dihapus/diubah, semua trigger akan gagal dan `updated_at` tidak akan ter-update otomatis. |

---

## View & Materialized View

### Regular Views (9) — **Sudah digantikan oleh Materialized View**

Regular views yang dibuat di `20260621000001_dashboard_views.sql` telah di-drop dan diganti dengan materialized view di `20260621000003_materialized_views.sql`. Pada rollback, views dibuat ulang sebagai stub kosong.

| View Reguler | Pengganti MV | Status |
|-------------|-------------|--------|
| `v_dashboard_stats` | `mv_dashboard_stats` | Digantikan |
| `v_dashboard_distribusi_gizi` | `mv_dashboard_distribusi_gizi` | Digantikan |
| `v_dashboard_tren_stunting` | `mv_dashboard_tren_stunting` | Digantikan |
| `v_dashboard_stunting_per_wilayah` | `mv_dashboard_stunting_per_wilayah` | Digantikan |
| `v_dashboard_kehadiran_bulanan` | `mv_dashboard_kehadiran_bulanan` | Digantikan |
| `v_dashboard_jadwal_terdekat` | `mv_dashboard_jadwal_terdekat` | Digantikan |
| `v_public_stats` | `mv_public_stats` | Digantikan |
| `v_riwayat_pemeriksaan` | `mv_riwayat_pemeriksaan` | Digantikan |
| `v_tumbuh_kembang` | `mv_tumbuh_kembang` | Digantikan |

### Materialized Views (11) — Aktif digunakan

#### `mv_dashboard_stats`
| Aspek | Detail |
|--------|--------|
| **Deskripsi** | Ringkasan dashboard: total_pasien, perlu_verifikasi, tindak_lanjut, kasus_stunting, jadwal_posyandu, total_balita, cakupan_persentase |
| **Tabel Sumber** | `user_account`, `pasien`, `hasil_pemeriksaan`, `tindak_lanjut`, `jadwal_imunisasi`, `ibu_hamil`, `anak` |
| **Digunakan API** | GET /dashboard/stats |
| **Fungsionalitas** | Dashboard utama — statistik one-shot |
| **Dampak Perubahan** | **SEDANG** — perlu REFRESH MATERIALIZED VIEW setelah data berubah. Data tidak real-time. |

#### `mv_dashboard_distribusi_gizi`
| Aspek | Detail |
|--------|--------|
| **Deskripsi** | Distribusi status gizi (Gizi Baik, Gizi Kurang, dll) berdasarkan pemeriksaan terakhir per pasien |
| **Tabel Sumber** | `hasil_pemeriksaan`, `jadwal_imunisasi`, `pasien` |
| **Digunakan API** | GET /dashboard/distribusi-gizi |
| **Fungsionalitas** | Pie chart distribusi gizi di dashboard |

#### `mv_dashboard_tren_stunting`
| Aspek | Detail |
|--------|--------|
| **Deskripsi** | Tren stunting bulanan |
| **Tabel Sumber** | `hasil_pemeriksaan`, `jadwal_imunisasi` |
| **Digunakan API** | GET /dashboard/tren-stunting |
| **Fungsionalitas** | Line chart tren stunting |

#### `mv_dashboard_stunting_per_wilayah`
| Aspek | Detail |
|--------|--------|
| **Deskripsi** | Prevalensi stunting per kabupaten/kota |
| **Tabel Sumber** | `hasil_pemeriksaan`, `jadwal_imunisasi`, `pasien`, `posyandu`, `lokasi` |
| **Digunakan API** | GET /dashboard/stunting-per-wilayah |
| **Fungsionalitas** | Bar chart / map stunting per wilayah |

#### `mv_dashboard_kehadiran_bulanan`
| Aspek | Detail |
|--------|--------|
| **Deskripsi** | Tren kehadiran posyandu bulanan |
| **Tabel Sumber** | `hasil_pemeriksaan` |
| **Digunakan API** | GET /dashboard/kehadiran-bulanan |
| **Fungsionalitas** | Line chart kehadiran posyandu |

#### `mv_dashboard_jadwal_terdekat`
| Aspek | Detail |
|--------|--------|
| **Deskripsi** | 10 jadwal imunisasi terdekat |
| **Tabel Sumber** | `jadwal_imunisasi`, `pasien`, `user_account` |
| **Digunakan API** | GET /dashboard/jadwal-terdekat |
| **Fungsionalitas** | Widget jadwal mendatang di dashboard |

#### `mv_public_stats`
| Aspek | Detail |
|--------|--------|
| **Deskripsi** | Statistik publik: total_pasien, balita_dipantau, kasus_stunting, total_artikel |
| **Tabel Sumber** | `pasien`, `anak`, `hasil_pemeriksaan`, `artikel` |
| **Digunakan API** | GET /stats (public) |
| **Fungsionalitas** | Landing page — statistik publik tanpa login |

#### `mv_riwayat_pemeriksaan`
| Aspek | Detail |
|--------|--------|
| **Deskripsi** | Riwayat pemeriksaan per pasien |
| **Tabel Sumber** | `hasil_pemeriksaan`, `jadwal_imunisasi`, `user_account` |
| **Digunakan API** | GET /monitoring/pasien/{id}/riwayat-pemeriksaan |
| **Fungsionalitas** | Timeline riwayat pemeriksaan pasien |

#### `mv_tumbuh_kembang`
| Aspek | Detail |
|--------|--------|
| **Deskripsi** | Data tumbuh kembang per bulan (BB/TB) |
| **Tabel Sumber** | `hasil_pemeriksaan`, `jadwal_imunisasi` |
| **Digunakan API** | GET /monitoring/pasien/{id}/tumbuh-kembang |
| **Fungsionalitas** | Grafik pertumbuhan anak (WHO standards) |

#### `mv_ibu_hamil_stats`
| Aspek | Detail |
|--------|--------|
| **Deskripsi** | Statistik ibu hamil: total, per trimester, melahirkan, nifas, keguguran |
| **Tabel Sumber** | `ibu_hamil` |
| **Digunakan API** | GET /dashboard/ibu-hamil-stats |
| **Fungsionalitas** | Dashboard KIA — statistik kehamilan |

#### `mv_ibu_hamil_per_wilayah`
| Aspek | Detail |
|--------|--------|
| **Deskripsi** | Distribusi ibu hamil per kabupaten/kota |
| **Tabel Sumber** | `ibu_hamil`, `pasien`, `posyandu`, `lokasi` |
| **Digunakan API** | GET /dashboard/ibu-hamil-per-wilayah |
| **Fungsionalitas** | Dashboard KIA — distribusi wilayah |

---

## Index

### Performance Indexes — High Priority (Hot Path)

| # | Index | Tabel | Kolom | Menguntungkan API |
|---|-------|-------|-------|-------------------|
| 1 | `idx_user_account_lokasi` | `user_account` | `id_lokasi` | GET /users (filter by lokasi) |
| 2 | `idx_user_account_deleted` | `user_account` | `is_deleted` | Semua query user_account (filter soft-delete) |
| 3 | `idx_user_account_status_verifikasi` | `user_account` | `status_verifikasi` | GET /users (filter verifikasi), PATCH /users/{id}/verification |
| 4 | `idx_user_account_lokasi_deleted` | `user_account` | `id_lokasi, is_deleted` | GET /users (compound filter) |
| 5 | `idx_pasien_posyandu` | `pasien` | `id_posyandu` | GET /monitoring/pasien, GET /laporan/tindak-lanjut |
| 6 | `idx_jadwal_imunisasi_pasien` | `jadwal_imunisasi` | `id_pasien` | GET /imunisasi/pasien/{id}, GET /imunisasi |
| 7 | `idx_jadwal_imunisasi_pasien_tanggal` | `jadwal_imunisasi` | `id_pasien, tanggal_jadwal` | GET /imunisasi/pasien/{id} (ordered by date) |
| 8 | `idx_hasil_pemeriksaan_jadwal` | `hasil_pemeriksaan` | `id_jadwal_imunisasi` | GET /monitoring/pemeriksaan/{id}, GET /tindak-lanjut/status |
| 9 | `idx_hasil_pemeriksaan_petugas` | `hasil_pemeriksaan` | `id_petugas_input` | GET /monitoring/pemeriksaan (filter petugas), GET /monitoring/semua-pemeriksaan |
| 10 | `idx_hasil_pemeriksaan_petugas_tanggal` | `hasil_pemeriksaan` | `id_petugas_input, created_at` | GET /monitoring/pemeriksaan (filter + sort) |
| 11 | `idx_notifikasi_user` | `notifikasi` | `id_user` | GET /notifikasi, GET /notifikasi/{id}, PATCH /notifikasi/{id}/read |
| 12 | `idx_notifikasi_user_baca` | `notifikasi` | `id_user, status_baca` | PATCH /notifikasi/read-all, GET /notifikasi (filter unread) |
| 13 | `idx_ibu_hamil_pasien` | `ibu_hamil` | `id_pasien` | GET /monitoring/pasien/{id} (JOIN ibu_hamil) |

### Performance Indexes — Medium Priority (FK Joins)

| # | Index | Tabel | Kolom | Menguntungkan API |
|---|-------|-------|-------|-------------------|
| 14 | `idx_lokasi_bagian_dari` | `lokasi` | `bagian_dari` | GET /lokasi (parent-child lookup) |
| 15 | `idx_bidan_wilayah_kerja` | `bidan` | `wilayah_kerja` | GET /notifikasi/bidan (filter wilayah) |
| 16 | `idx_posyandu_lokasi` | `posyandu` | `id_lokasi` | GET /laporan/tindak-lanjut (JOIN lokasi via posyandu) |
| 17 | `idx_posyandu_bidan` | `posyandu` | `id_bidan` | GET /notifikasi/bidan, GET /monitoring/semua-pemeriksaan |
| 18 | `idx_kader_posyandu_posyandu` | `kader_posyandu` | `id_posyandu` | GET /monitoring/semua-pemeriksaan (filter kader) |
| 19 | `idx_fasilitas_kesehatan_lokasi` | `fasilitas_kesehatan` | `id_lokasi` | GET /faskes (filter by wilayah) |
| 20 | `idx_anak_ibu_hamil` | `anak` | `id_ibu_hamil` | GET /monitoring/pasien/{id} (JOIN anak-ibu_hamil) |
| 21 | `idx_anak_wali` | `anak` | `id_wali` | Ownership check: GET /monitoring/pasien, GET /imunisasi, GET /monitoring/pemeriksaan/{id} |
| 22 | `idx_tindak_lanjut_bidan` | `tindak_lanjut` | `id_bidan` | GET /notifikasi/bidan (jadwal kontrol) |
| 23 | `idx_user_session_user` | `user_session` | `id_user` | POST /auth/refresh, POST /auth/logout |
| 24 | `idx_audit_log_user` | `audit_log` | `id_user` | GET /admin/audit-logs (filter user) |
| 25 | `idx_audit_log_session` | `audit_log` | `id_user_session` | GET /admin/audit-logs (filter session) |

### Filter/Sort Indexes

| # | Index | Tabel | Kolom | Menguntungkan API |
|---|-------|-------|-------|-------------------|
| 26 | `idx_artikel_status` | `artikel` | `status_artikel` | GET /artikel (filter published), GET /artikel/pending |
| 27 | `idx_artikel_tanggal_publish` | `artikel` | `tanggal_publish DESC` | GET /artikel (sort by publish date) |
| 28 | `idx_rujukan_status` | `rujukan` | `status_rujukan` | GET /tindak-lanjut/status (filter rujukan), GET /notifikasi/bidan |
| 29 | `idx_tindak_lanjut_status` | `tindak_lanjut` | `status_pasien` | GET /tindak-lanjut/status, GET /tindak-lanjut/pasien |
| 30 | `idx_hasil_pemeriksaan_status_gizi` | `hasil_pemeriksaan` | `status_gizi` | GET /tindak-lanjut/pasien (filter gizi buruk) |

### Additional FK Indexes

| # | Index | Tabel | Kolom | Menguntungkan API |
|---|-------|-------|-------|-------------------|
| 31 | `idx_artikel_id_penulis` | `artikel` | `id_penulis` | GET /artikel/semua (JOIN penulis) |
| 32 | `idx_artikel_id_verifikator` | `artikel` | `id_verifikator` | GET /artikel/{id} (JOIN verifikator) |
| 33 | `idx_rujukan_id_faskes` | `rujukan` | `id_faskes` | GET /tindak-lanjut/status, GET /tindak-lanjut/{id} (JOIN faskes) |
| 34 | `idx_user_account_id_pendidikan` | `user_account` | `id_pendidikan` | GET /users (JOIN referensi) |
| 35 | `idx_user_account_id_pekerjaan` | `user_account` | `id_pekerjaan` | GET /users (JOIN referensi) |
| 36 | `idx_user_account_id_pendapatan` | `user_account` | `id_pendapatan` | GET /users (JOIN referensi) |

### Dashboard Indexes

| # | Index | Tabel | Kolom | Menguntungkan API |
|---|-------|-------|-------|-------------------|
| 37 | `idx_hasil_pemeriksaan_stunting_created` | `hasil_pemeriksaan` | `status_stunting, created_at DESC` | GET /dashboard/tren-stunting, GET /dashboard/stunting-per-wilayah |
| 38 | `idx_hasil_pemeriksaan_jadwal_created` | `hasil_pemeriksaan` | `id_jadwal_imunisasi, created_at DESC` | GET /tindak-lanjut/pasien (LATERAL subquery) |
| 39 | `idx_hasil_pemeriksaan_created` | `hasil_pemeriksaan` | `created_at DESC` | GET /monitoring/pemeriksaan (sorting), dashboard MVs |

---

## Trigger

Semua trigger menggunakan fungsi `update_updated_at_column()` dan fire pada event **BEFORE UPDATE**.

### Dampak perubahan trigger:
- **Menghapus trigger** → kolom `updated_at` tidak akan ter-update otomatis. Application layer harus set manual.
- **Mengubah trigger** → semua UPDATE query di 15 tabel terpengaruh.

| # | Trigger | Tabel | Event | API yang Memicu |
|---|---------|-------|-------|-----------------|
| 1 | `trg_user_account_updated_at` | `user_account` | BEFORE UPDATE | PATCH /users/{id}, PATCH /users/{id_user}/verification, DELETE /users/{id} (soft-delete), PATCH /users/{id}/role |
| 2 | `trg_dinas_kesehatan_updated_at` | `dinas_kesehatan` | BEFORE UPDATE | PATCH /users/{id}/role (soft-delete) |
| 3 | `trg_bidan_updated_at` | `bidan` | BEFORE UPDATE | PATCH /users/{id}/role (soft-delete) |
| 4 | `trg_posyandu_updated_at` | `posyandu` | BEFORE UPDATE | — (tidak ada UPDATE API ke posyandu) |
| 5 | `trg_kader_posyandu_updated_at` | `kader_posyandu` | BEFORE UPDATE | PATCH /users/{id}/role (soft-delete) |
| 6 | `trg_pasien_updated_at` | `pasien` | BEFORE UPDATE | PATCH /pasien/{id}, DELETE /pasien/{id} (soft-delete) |
| 7 | `trg_fasilitas_kesehatan_updated_at` | `fasilitas_kesehatan` | BEFORE UPDATE | — (tidak ada UPDATE API) |
| 8 | `trg_ibu_hamil_updated_at` | `ibu_hamil` | BEFORE UPDATE | PATCH /pasien/{id}, DELETE /pasien/{id} (soft-delete) |
| 9 | `trg_anak_updated_at` | `anak` | BEFORE UPDATE | PATCH /pasien/{id}, DELETE /pasien/{id} (soft-delete) |
| 10 | `trg_jadwal_imunisasi_updated_at` | `jadwal_imunisasi` | BEFORE UPDATE | PUT /imunisasi/{id}, PATCH /imunisasi/{id}/realisasi |
| 11 | `trg_artikel_updated_at` | `artikel` | BEFORE UPDATE | PATCH /artikel/{id}, PATCH /artikel/{id}/review |
| 12 | `trg_hasil_pemeriksaan_updated_at` | `hasil_pemeriksaan` | BEFORE UPDATE | PUT /monitoring/pemeriksaan/{id}, PATCH /monitoring/pemeriksaan/{id}/verify |
| 13 | `trg_tindak_lanjut_updated_at` | `tindak_lanjut` | BEFORE UPDATE | — (tidak ada UPDATE API langsung, hanya INSERT) |
| 14 | `trg_rujukan_updated_at` | `rujukan` | BEFORE UPDATE | PATCH /rujukan/{id}/status |
| 15 | `user_session_updated_at` | `user_session` | BEFORE UPDATE | POST /auth/refresh, POST /auth/logout |

### Tabel tanpa trigger `updated_at`:
- `audit_log` — append-only, tidak ada UPDATE
- `pendidikan` — referensi statis
- `pekerjaan` — referensi statis
- `kategori_pendapatan` — referensi statis
- `notifikasi` — updated_at di-set manual via application code (Jet ORM)
- `lokasi` — referensi, tidak ada kolom updated_at

---

## Impact Matrix (Ringkasan)

### Level Dampak

| Level | Definisi | Jumlah Objek |
|-------|----------|-------------|
| **KRITIS** | Perubahan merusak banyak endpoint dan fungsionalitas inti | 5 tabel, 1 function |
| **SEDANG** | Perubahan berdampak pada beberapa endpoint spesifik | 5 tabel, 11 MV |
| **RENDAH** | Perubahan berdampak terbatas | 8 tabel |
| **MINIMAL** | Tabel referensi, tidak ada query langsung | 3 tabel |

### Tabel Kritis
| Tabel | Alasan |
|-------|--------|
| `user_account` | Seluruh autentikasi, otorisasi, dan data pengguna |
| `pasien` | Sentral data pasien, cascade ke ibu_hamil & anak |
| `anak` | Ownership check di banyak endpoint (akses kontrol) |
| `jadwal_imunisasi` | Jembatan imunisasi-pemeriksaan, HARD DELETE, cascade dari pasien |
| `hasil_pemeriksaan` | Sumber data 8 dashboard MV, HARD DELETE, dasar tindak lanjut |

### Materialized View Refresh Strategy
Semua 11 MV menggunakan PostgreSQL `REFRESH MATERIALIZED VIEW`. Tidak ada `REFRESH MATERIALIZED VIEW CONCURRENTLY` karena tidak ada `UNIQUE INDEX` pada MV. Artinya:
- **Refresh memblokir pembacaan** selama proses refresh
- Perlu dijadwalkan (cron/pg_cron) atau di-trigger setelah batch update

---

## DML Operation Summary

| Operasi | Tabel yang Terkena |
|---------|-------------------|
| **INSERT** | `user_account`, `bidan`, `kader_posyandu`, `dinas_kesehatan`, `pasien`, `ibu_hamil`, `anak`, `jadwal_imunisasi`, `artikel`, `hasil_pemeriksaan`, `tindak_lanjut`, `rujukan`, `notifikasi`, `user_session`, `audit_log` |
| **UPDATE** | `user_account`, `bidan` (soft-delete), `kader_posyandu` (soft-delete), `dinas_kesehatan` (soft-delete), `pasien`, `ibu_hamil`, `anak`, `jadwal_imunisasi`, `artikel`, `hasil_pemeriksaan`, `rujukan`, `notifikasi`, `user_session` |
| **DELETE (HARD)** | `jadwal_imunisasi`, `hasil_pemeriksaan`, `artikel` |
| **DELETE (SOFT)** | `user_account`, `dinas_kesehatan`, `bidan`, `kader_posyandu`, `pasien`, `ibu_hamil`, `anak` |
| **TRANSACTION** | Delete pasien: `anak` + `ibu_hamil` + `pasien` dalam satu transaksi |

> **Peringatan:** 3 tabel menggunakan HARD DELETE (`jadwal_imunisasi`, `hasil_pemeriksaan`, `artikel`). Data yang terhapus tidak bisa direcover tanpa backup.
