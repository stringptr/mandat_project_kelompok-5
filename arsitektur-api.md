# Arsitektur API — SiGizi

> Sistem Informasi Gizi & Monitoring Imunisasi Anak  
> Backend: Go + Huma v2 + Chi Router | Database: PostgreSQL 18 | ORM: Go-Jet v2

---

## Daftar Isi

1. [Hirarki Role & Middleware](#hirarki-role--middleware)
2. [Auth & Session](#1-auth--session)
3. [Lokasi](#2-lokasi)
4. [Fasilitas Kesehatan](#3-fasilitas-kesehatan)
5. [User Account](#4-user-account)
6. [Notifikasi](#5-notifikasi)
7. [Pasien](#6-pasien)
8. [Pemeriksaan](#7-pemeriksaan)
9. [Imunisasi](#8-imunisasi)
10. [Artikel](#9-artikel)
11. [Tindak Lanjut](#10-tindak-lanjut)
12. [Dashboard](#11-dashboard)

---

## Hirarki Role & Middleware

```
publicGroup          → tanpa autentikasi
nonAuthOnlyGroup     → hanya non-autentikasi (register/login)
authAccess           → JWT access token valid
authRefresh          → JWT refresh token valid
  └─ userGroup       → role USER (semua role terautentikasi)
       ├─ adminGroup → role ADMIN (Bidan, Kader, Dinkes)
       │    ├─ bidanGroup    → role BIDAN
       │    ├─ kaderGroup    → role KADER
       │    └─ superAdminGroup → role SUPER_ADMIN
       └─ dinkesGroup → role DINKES
```

**File route:** `backend/internal/api/v1/routes.go`  
**Middleware:** `backend/internal/middleware/` (access.go, refresh.go, ratelimit.go, ip.go)

---

## 1. Auth & Session

### POST /auth/register
| Aspek | Detail |
|--------|--------|
| **Role** | Non-authenticated only |
| **Handler** | `feature/auth/handler.go` — `Register` |
| **Input** | Body: `RegisterRequest` |
| | `email` (string, format:email, required) |
| | `password` (string, min:8, max:255, required) |
| | `no_hp` (string, max:20, required) |
| | `nama` (string, min:1, max:255, required) |
| | `nik` (string, 16 digits, required) |
| | `jenis_kelamin` (enum: Laki-Laki, Perempuan, required) |
| | `tanggal_lahir` (date-time, required) |
| | `id_lokasi` (int32, min:1, required) |
| | `id_pendidikan` (int32, optional) |
| | `id_pekerjaan` (int32, optional) |
| | `id_pendapatan` (int32, optional) |
| | `jumlah_tanggungan` (int32, optional) |
| | `role` (enum: Bidan, Kader, Dinkes, Ibu Hamil, Anak) |
| | `no_str` (string, optional — Bidan) |
| | `wilayah_kerja` (int32, optional — Bidan) |
| | `no_sk` (string, optional — Kader) |
| | `id_posyandu` (int32, optional — Kader/Anak) |
| **Output** | `201` — `{ "message": "Register berhasil. Akun sedang diverifikasi." }` |
| **Table READ** | `user_account` (cek duplikat email/NIK) |
| **Table WRITE** | `user_account` (INSERT), `bidan` / `kader_posyandu` / `dinas_kesehatan` / `pasien` (INSERT role record) |
| **SP** | Tidak ada |
| **DML DB** | `INSERT INTO user_account (...) VALUES (...) RETURNING *` |
| | `INSERT INTO bidan / kader_posyandu / dinas_kesehatan / pasien (...) VALUES (...)` |

### POST /auth/login
| Aspek | Detail |
|--------|--------|
| **Role** | Non-authenticated only |
| **Handler** | `feature/auth/handler.go` — `Login` |
| **Input** | Body: `LoginRequest` |
| | `email` (string, email, optional — salah satu dengan NIK) |
| | `nik` (string, 16 digits, optional — salah satu dengan email) |
| | `password` (string, min:8, max:255, required) |
| **Output** | `200` — `{ "access_token": "...", "refresh_token": "...", "expires_in": 3600 }` + Set-Cookie header (HttpOnly, SameSite=Strict) |
| **Table READ** | `user_account` (by email/NIK, cek password hash), `bidan`, `kader_posyandu`, `dinas_kesehatan`, `pasien`, `ibu_hamil`, `anak` (EXISTS cek role) |
| **Table WRITE** | `user_session` (INSERT), `audit_log` (INSERT) |
| **SP** | Tidak ada |
| **DML DB** | `SELECT * FROM user_account WHERE (email = $1 OR nik = $2) AND is_deleted = false` |
| | `SELECT EXISTS(SELECT 1 FROM dinas_kesehatan WHERE id_user = $1 AND is_deleted = false)` — ×6 EXISTS untuk 6 role |
| | `INSERT INTO user_session (...) VALUES (...)` |
| | `INSERT INTO audit_log (...) VALUES (...)` |

### POST /auth/refresh
| Aspek | Detail |
|--------|--------|
| **Role** | Auth (refresh token) |
| **Handler** | `feature/auth/handler.go` — `Refresh` |
| **Input** | Cookie: `refresh_token` |
| **Output** | `200` — `{ "access_token": "...", "refresh_token": "...", "expires_in": 3600 }` + cookie baru |
| **Table READ** | `user_session` (validasi session by refresh token JTI) |
| **Table WRITE** | `user_session` (UPDATE rotate session) |
| **SP** | Tidak ada |
| **DML DB** | `SELECT * FROM user_session WHERE id_session = $1` |
| | `UPDATE user_session SET status_session = 'KADALUWARSA', updated_at = NOW() WHERE id_session = $1` |
| | `INSERT INTO user_session (...) VALUES (...)` — session baru |

### POST /auth/logout
| Aspek | Detail |
|--------|--------|
| **Role** | Public |
| **Handler** | `feature/auth/handler.go` — `Logout` |
| **Input** | Cookie: `refresh_token` |
| **Output** | `200` — `{ "message": "Logout berhasil." }` + clear cookie header |
| **Table WRITE** | `user_session` (UPDATE — cabut session) |
| **SP** | Tidak ada |
| **DML DB** | `UPDATE user_session SET status_session = 'DICABUT', updated_at = NOW() WHERE id_session = $1` |

### GET /auth/me
| Aspek | Detail |
|--------|--------|
| **Role** | USER |
| **Handler** | `feature/auth/handler.go` — `Me` |
| **Input** | Cookie: `access_token` (JWT claim dari context) |
| **Output** | `200` — `{ "id_user": ..., "email": "...", "roles": [...], "nama": "..." }` (JWT claim) |
| **Table READ** | — (dari JWT context, tidak query DB) |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | Tidak ada (data dari JWT token) |

### PATCH /users/{id_user}/verification
| Aspek | Detail |
|--------|--------|
| **Role** | ADMIN |
| **Handler** | `feature/auth/handler.go` — `VerifyUser` |
| **Input** | Path: `id_user` (int32) |
| | Body: `{ "status": "Aktif|Ditolak", "alasan_penolakan": "..." }` |
| **Output** | `200` — `{}` |
| **Table READ** | `user_account` (by id) |
| **Table WRITE** | `user_account` (UPDATE status_verifikasi) |
| **SP** | Tidak ada |
| **DML DB** | `UPDATE user_account SET status_verifikasi = $1, updated_at = NOW() WHERE id_user = $2 AND is_deleted = false` |

---

## 2. Lokasi

### GET /lokasi
| Aspek | Detail |
|--------|--------|
| **Role** | Public |
| **Handler** | `feature/lokasi/handler.go` — `GetLokasi` |
| **Input** | Query: `tipe` (enum: Provinsi, Kabupaten, Kota, Kecamatan, Kelurahan), `bagian_dari` (int32, parent ID) |
| **Output** | `200` — `[ { "id_lokasi": 1, "nama_lokasi": "...", "tipe_lokasi": "...", "bagian_dari": null } ]` |
| **Table READ** | `lokasi` |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | `SELECT id_lokasi, nama_lokasi, tipe_lokasi::text, bagian_dari FROM lokasi WHERE tipe_lokasi::text = $1 AND (bagian_dari = $2 OR bagian_dari IS NULL) ORDER BY nama_lokasi ASC` |

---

## 3. Fasilitas Kesehatan

### GET /faskes
| Aspek | Detail |
|--------|--------|
| **Role** | USER |
| **Handler** | `feature/faskes/handler.go` — `GetFaskes` |
| **Input** | Query: `search` (string, optional) |
| **Output** | `200` — `[ { "id_faskes": 1, "nama_faskes": "...", "tipe_faskes": "..." } ]` |
| **Table READ** | `fasilitas_kesehatan` |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | `SELECT id_faskes, nama_faskes, tipe_faskes FROM fasilitas_kesehatan WHERE is_deleted = false AND nama_faskes ILIKE $1 ORDER BY nama_faskes ASC` |

---

## 4. User Account

### GET /users
| Aspek | Detail |
|--------|--------|
| **Role** | ADMIN |
| **Handler** | `feature/userAccount/handler.go` — `GetAllUsers` |
| **Input** | Query: `page` (int, default:1), `per_page` (int, default:20, max:100), `q` (string), `role` (enum: Bidan,Kader,Dinkes,Pasien,Ibu Hamil,Anak), `status_verifikasi` (enum: Pending,Aktif,Ditolak) |
| **Output** | `200` — `{ "users": [...], "meta": { "page":1, "per_page":20, "total":100 } }` |
| | UserListItem: `id_user`, `nama`, `nik`, `email`, `no_hp`, `jenis_kelamin`, `status_verifikasi`, `roles`, `id_lokasi`, `nama_lokasi`, `created_at`, `updated_at` |
| **Table READ** | `user_account`, `dinas_kesehatan`, `bidan`, `kader_posyandu`, `pasien`, `ibu_hamil`, `anak` (LEFT JOIN for role detection), `lokasi` (GetLokasiNames) |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | `SELECT ua.* FROM user_account ua LEFT JOIN dinas_kesehatan dk ... LEFT JOIN bidan b ... WHERE ... AND role_check ORDER BY id_user DESC LIMIT $1 OFFSET $2` |
| | `SELECT COUNT(*) FROM user_account ua LEFT JOIN ...` |
| | `SELECT id_lokasi, nama_lokasi FROM lokasi WHERE id_lokasi = ANY($1)` |

### GET /users/{id}
| Aspek | Detail |
|--------|--------|
| **Role** | USER |
| **Handler** | `feature/userAccount/handler.go` — `GetUserByID` |
| **Input** | Path: `id` (int32, min:1) |
| **Output** | `200` — `UserDetailResponse` (id_user, email, no_hp, nama, nik, jenis_kelamin, tanggal_lahir, status_verifikasi, id_lokasi, id_pendidikan, id_pekerjaan, id_pendapatan, jumlah_tanggungan, roles, created_at, updated_at) |
| **Table READ** | `user_account`, `bidan`, `kader_posyandu`, `dinas_kesehatan`, `pasien`, `ibu_hamil`, `anak` (EXISTS for roles) |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | `SELECT * FROM user_account WHERE id_user = $1 AND is_deleted = false` |
| | `SELECT EXISTS(SELECT 1 FROM dinas_kesehatan WHERE id_user = $1 AND is_deleted = false)` — ×6 EXISTS |

### PATCH /users/{id}
| Aspek | Detail |
|--------|--------|
| **Role** | USER (self) atau ADMIN/SUPER_ADMIN |
| **Handler** | `feature/userAccount/handler.go` — `UpdateUser` |
| **Input** | Path: `id` (int32, min:1) |
| | Body: `UpdateUserRequest` — semua field optional: `email`, `no_hp`, `nama`, `nik`, `jenis_kelamin`, `tanggal_lahir`, `id_lokasi`, `id_pendidikan`, `id_pekerjaan`, `id_pendapatan`, `jumlah_tanggungan` |
| **Output** | `200` — `UserDetailResponse` (data terbaru setelah update) |
| **Table READ** | `user_account` (2× — before & after update) |
| **Table WRITE** | `user_account` (UPDATE) |
| **SP** | Tidak ada |
| **DML DB** | `UPDATE user_account SET email = $1, no_hp = $2, nama = $3, nik = $4, jenis_kelamin = $5, tanggal_lahir = $6, id_lokasi = $7, id_pendidikan = $8, id_pekerjaan = $9, id_pendapatan = $10, jumlah_tanggungan = $11, updated_at = NOW() WHERE id_user = $12 AND is_deleted = false` |

### POST /users
| Aspek | Detail |
|--------|--------|
| **Role** | SUPER_ADMIN |
| **Handler** | `feature/userAccount/handler.go` — `CreateUser` |
| **Input** | Body: `CreateUserRequest` — sama dengan RegisterRequest tanpa Ibu Hamil/Anak role |
| | `email`, `password`, `no_hp`, `nama`, `nik`, `jenis_kelamin`, `tanggal_lahir`, `id_lokasi`, `id_pendidikan`, `id_pekerjaan`, `id_pendapatan`, `jumlah_tanggungan`, `role` (Bidan,Kader,Dinkes), `no_str`, `wilayah_kerja`, `no_sk`, `id_posyandu` |
| **Output** | `201` — `{ "id_user": 123 }` |
| **Table READ** | `user_account` (cek duplikat) |
| **Table WRITE** | `user_account` (INSERT), `bidan` / `kader_posyandu` / `dinas_kesehatan` (INSERT) |
| **SP** | Tidak ada |
| **DML DB** | `INSERT INTO user_account (...) VALUES (...) RETURNING id_user` |
| | `INSERT INTO bidan / kader_posyandu / dinas_kesehatan (...) VALUES (...)` |

### PATCH /users/{id}/role
| Aspek | Detail |
|--------|--------|
| **Role** | SUPER_ADMIN |
| **Handler** | `feature/userAccount/handler.go` — `UpdateUserRole` |
| **Input** | Path: `id` (int32, min:1) |
| | Body: `role` (enum: Bidan,Kader,Dinkes), `no_str`, `wilayah_kerja`, `no_sk`, `id_posyandu` |
| **Output** | `200` — `{}` |
| **Table WRITE** | `bidan`, `kader_posyandu`, `dinas_kesehatan` (DELETE role lama → INSERT role baru) |
| **SP** | Tidak ada |
| **DML DB** | `UPDATE bidan SET is_deleted = true, deleted_at = NOW() WHERE id_user = $1` |
| | `UPDATE kader_posyandu SET is_deleted = true, deleted_at = NOW() WHERE id_user = $1` |
| | `UPDATE dinas_kesehatan SET is_deleted = true, deleted_at = NOW() WHERE id_user = $1` |
| | `INSERT INTO [role_baru] (...) VALUES (...)` |

### GET /admin/audit-logs
| Aspek | Detail |
|--------|--------|
| **Role** | SUPER_ADMIN |
| **Handler** | `feature/userAccount/handler.go` — `GetAuditLogs` |
| **Input** | Query: `page` (int, default:1), `per_page` (int, default:20, max:100) |
| **Output** | `200` — `{ "audit_logs": [...], "meta": {...} }` |
| | AuditLogItem: `id_log`, `tipe_aktor`, `id_user`, `tipe_aktivitas`, `berhasil`, `endpoint`, `table_name`, `record_id`, `detail`, `ip_address`, `user_agent`, `waktu_aktivitas` |
| **Table READ** | `audit_log` |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | `SELECT * FROM audit_log ORDER BY waktu_aktivitas DESC LIMIT $1 OFFSET $2` |
| | `SELECT COUNT(*) FROM audit_log` |

### DELETE /users/{id}
| Aspek | Detail |
|--------|--------|
| **Role** | SUPER_ADMIN |
| **Handler** | `feature/userAccount/handler.go` — `DeleteUser` |
| **Input** | Path: `id` (int32, min:1) |
| **Output** | `200` — `{}` |
| **Table WRITE** | `user_account` (soft-delete) |
| **SP** | Tidak ada |
| **DML DB** | `UPDATE user_account SET is_deleted = true, deleted_at = NOW() WHERE id_user = $1` |

---

## 5. Notifikasi

### GET /notifikasi
| Aspek | Detail |
|--------|--------|
| **Role** | USER |
| **Handler** | `feature/notification/handler.go` — `GetNotifikasi` |
| **Input** | Query: `page` (int, default:1), `per_page` (int, default:20), `q` (string, optional) |
| **Output** | `200` — `{ "notifikasi": [...], "meta": {...} }` |
| | NotifikasiItem: `id_notifikasi`, `judul`, `pesan`, `tipe_notifikasi`, `status_baca`, `tanggal_kirim` |
| **Table READ** | `notifikasi` |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | `SELECT * FROM notifikasi WHERE id_user = $1 ORDER BY tanggal_kirim DESC LIMIT $2 OFFSET $3` |
| | `SELECT COUNT(*) FROM notifikasi WHERE id_user = $1` |

### GET /notifikasi/{id}
| Aspek | Detail |
|--------|--------|
| **Role** | USER |
| **Handler** | `feature/notification/handler.go` — `GetNotifikasiDetail` |
| **Input** | Path: `id` (int32) |
| **Output** | `200` — NotifikasiDetail: `id_notifikasi`, `judul`, `pesan`, `tipe_notifikasi`, `status_baca`, `tanggal_kirim`, `aksi` (label + url) |
| **Table READ** | `notifikasi` |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | `SELECT * FROM notifikasi WHERE id_notifikasi = $1 AND id_user = $2` |

### PATCH /notifikasi/{id}/read
| Aspek | Detail |
|--------|--------|
| **Role** | USER |
| **Handler** | `feature/notification/handler.go` — `MarkRead` |
| **Input** | Path: `id` (int32) |
| **Output** | `200` — `{ "id_notifikasi": ..., "status_baca": true }` |
| **Table WRITE** | `notifikasi` (UPDATE) |
| **SP** | Tidak ada |
| **DML DB** | `UPDATE notifikasi SET status_baca = true WHERE id_notifikasi = $1 AND id_user = $2` |

### PATCH /notifikasi/read-all
| Aspek | Detail |
|--------|--------|
| **Role** | USER |
| **Handler** | `feature/notification/handler.go` — `MarkAllRead` |
| **Input** | — |
| **Output** | `200` — `{ "jumlah_diperbarui": 5, "status": "success" }` |
| **Table WRITE** | `notifikasi` (UPDATE massal) |
| **SP** | Tidak ada |
| **DML DB** | `UPDATE notifikasi SET status_baca = true WHERE id_user = $1 AND status_baca = false` |

### GET /notifikasi/bidan
| Aspek | Detail |
|--------|--------|
| **Role** | BIDAN |
| **Handler** | `feature/notification/handler.go` — `GetBidanDashboard` |
| **Input** | — (id_user dari JWT) |
| **Output** | `200` — BidanNotificationResponse: `statistik` (jadwal_kontrol, risiko_stunting, rujukan_mendesak), `notifikasi_risiko_stunting[]`, `jadwal_monitoring[]`, `rujukan_mendesak[]`, `laporan_bulanan` |
| **Table READ** | `bidan`, `notifikasi`, `tindak_lanjut`, `rujukan`, `hasil_pemeriksaan`, `jadwal_imunisasi`, `pasien`, `posyandu`, `anak` |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | `SELECT id_user FROM bidan WHERE id_user = $1 AND is_deleted = false` |
| | `SELECT COUNT(*) FROM tindak_lanjut WHERE id_bidan = $1 AND jadwal_kontrol >= CURRENT_DATE` |
| | `SELECT COUNT(DISTINCT ji.id_pasien) FROM hasil_pemeriksaan hp INNER JOIN jadwal_imunisasi ji ... INNER JOIN pasien p ... INNER JOIN posyandu pos ... WHERE status_stunting IN ('Stunting Berat','Stunting','Berisiko Stunting')` |
| | `SELECT COUNT(*) FROM rujukan r INNER JOIN tindak_lanjut tl ... WHERE status_rujukan IN ('Diajukan','Diproses')` |
| | + list queries untuk masing-masing kategori |

### GET /notifikasi/aktivitas
| Aspek | Detail |
|--------|--------|
| **Role** | ADMIN |
| **Handler** | `feature/notification/handler.go` — `GetActivity` |
| **Input** | — (id_user dari JWT) |
| **Output** | `200` — NotificationActivity: `hari_ini[]`, `kemarin[]` (AktivitasItem: id_notifikasi, judul, status, timestamp) |
| **Table READ** | `notifikasi` |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | `SELECT * FROM notifikasi WHERE id_user = $1 AND tanggal_kirim >= CURRENT_DATE` — hari ini |
| | `SELECT * FROM notifikasi WHERE id_user = $1 AND tanggal_kirim >= CURRENT_DATE - INTERVAL '1 day' AND tanggal_kirim < CURRENT_DATE` — kemarin |

---

## 6. Pasien

### POST /pasien/ibu-hamil
| Aspek | Detail |
|--------|--------|
| **Role** | ADMIN |
| **Handler** | `feature/pasien/handler.go` — `DaftarIbuHamil` |
| **Input** | Body: `DaftarIbuHamilRequest` |
| | `id_user` (int32, min:1) — user_account yang sudah ada |
| | `id_posyandu` (int32, min:1) |
| | `hamil_ke` (int32, min:1) |
| | `bulan_mulai_hamil` (string, date) |
| | `hpht` (string, date — hari pertama haid terakhir) |
| | `status_kehamilan` (enum: Trimester 1, Trimester 2, Trimester 3, Melahirkan, Nifas, Keguguran) |
| **Output** | `201` — `{}` |
| **Table READ** | `user_account` (validasi), `pasien` (cek existing) |
| **Table WRITE** | `pasien` (INSERT), `ibu_hamil` (INSERT) |
| **SP** | Tidak ada |
| **DML DB** | `INSERT INTO pasien (id_pasien, id_posyandu) VALUES ($1, $2)` |
| | `INSERT INTO ibu_hamil (id_pasien, hamil_ke, bulan_mulai_hamil, hpht, status_kehamilan) VALUES ($1, $2, $3, $4, $5)` |

### POST /pasien/anak
| Aspek | Detail |
|--------|--------|
| **Role** | ADMIN |
| **Handler** | `feature/pasien/handler.go` — `DaftarAnak` |
| **Input** | Body: `DaftarAnakRequest` |
| | `id_user` (int32, min:1) |
| | `id_posyandu` (int32, min:1) |
| | `id_ibu_hamil` (int32, optional) |
| | `id_wali` (int32, min:1) |
| | `nama_anak` (string, min:1, max:255) |
| | `berat_lahir` (float64, min:0) |
| | `panjang_lahir` (float64, min:0) |
| | `hubungan_dengan_wali` (enum: Kandung, Tiri, Angkat) |
| **Output** | `201` — `{}` |
| **Table READ** | `user_account` (validasi), `ibu_hamil` (validasi, optional) |
| **Table WRITE** | `pasien` (INSERT), `anak` (INSERT) |
| **SP** | Tidak ada |
| **DML DB** | `INSERT INTO pasien (id_pasien, id_posyandu) VALUES ($1, $2)` |
| | `INSERT INTO anak (id_pasien, id_ibu_hamil, id_wali, nama_anak, berat_lahir, panjang_lahir, hubungan_dengan_wali) VALUES (...)` |

### GET /monitoring/pasien
| Aspek | Detail |
|--------|--------|
| **Role** | USER |
| **Handler** | `feature/pasien/handler.go` — `GetAll` |
| **Input** | Query: `page` (int, default:1), `per_page` (int, default:20, max:100), `q` (string) |
| | Petugas → lihat semua pasien; Non-petugas → lihat pasien milik sendiri |
| **Output** | `200` — `{ "pasien": [...], "meta": {...} }` |
| | PasienListItem: `id_pasien`, `nama`, `nik`, `jenis_kelamin`, `umur`, `nama_posyandu`, `jenis_pasien`, `status_kehamilan` |
| **Table READ** | `pasien`, `user_account`, `posyandu`, `ibu_hamil`, `anak` |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | Raw SQL: `SELECT p.id_pasien, ua.nama, ua.nik, ua.jenis_kelamin, ..., pos.nama_posyandu, COALESCE(ih.status_kehamilan) FROM pasien p INNER JOIN user_account ua ON p.id_pasien = ua.id_user INNER JOIN posyandu pos ON p.id_posyandu = pos.id_posyandu LEFT JOIN ibu_hamil ih ... LEFT JOIN anak a ... WHERE ... ILIKE ... ORDER BY ... LIMIT $1 OFFSET $2` |
| | Atau dengan filter `a.id_wali = $1` untuk non-petugas |

### GET /monitoring/pasien/search
| Aspek | Detail |
|--------|--------|
| **Role** | ADMIN |
| **Handler** | `feature/pasien/handler.go` — `Search` |
| **Input** | Query: `q` (string, max:255), `page` (int, default:1), `per_page` (int, default:20, max:100) |
| **Output** | `200` — `{ "pasien": [...], "meta": {...} }` |
| **Table READ** | `pasien`, `user_account`, `posyandu`, `ibu_hamil`, `anak` |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | Raw SQL dengan ILIKE pada `ua.nama` atau `ua.nik` |

### GET /monitoring/pasien/{id}
| Aspek | Detail |
|--------|--------|
| **Role** | USER |
| **Handler** | `feature/pasien/handler.go` — `GetByID` |
| **Input** | Path: `id` (int32, min:1) |
| **Output** | `200` — PasienDetailResponse: `id_pasien`, `nama`, `nik`, `email`, `no_hp`, `jenis_kelamin`, `tanggal_lahir`, `id_lokasi`, `nama_posyandu`, `id_posyandu`, `jenis_pasien`, `data_ibu_hamil`, `data_anak`, `created_at`, `updated_at` |
| **Table READ** | `pasien`, `user_account`, `posyandu`, `ibu_hamil`, `anak`, `wali` (user_account as wali) |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | Raw SQL: `SELECT p.*, ua.*, pos.nama_posyandu, ih.*, a.*, w.nama as nama_wali FROM pasien p INNER JOIN user_account ua ... INNER JOIN posyandu pos ... LEFT JOIN LATERAL (SELECT * FROM ibu_hamil WHERE id_pasien = p.id_pasien LIMIT 1) ih ON true LEFT JOIN LATERAL (SELECT ... FROM anak LEFT JOIN user_account w ON a.id_wali = w.id_user ...) a ON true WHERE p.id_pasien = $1` |

### PATCH /pasien/{id}
| Aspek | Detail |
|--------|--------|
| **Role** | ADMIN |
| **Handler** | `feature/pasien/handler.go` — `Update` |
| **Input** | Path: `id` (int32, min:1) |
| | Body: `UpdatePasienRequest` — semua optional: `id_posyandu`, `hamil_ke`, `bulan_mulai_hamil`, `hpht`, `status_kehamilan`, `nama_anak`, `berat_lahir`, `panjang_lahir`, `hubungan_dengan_wali`, `id_wali` |
| **Output** | `200` — PasienDetailResponse (data terbaru) |
| **Table READ** | `pasien` (existing data) |
| **Table WRITE** | `pasien` (UPDATE), `ibu_hamil` (UPDATE jika ada), `anak` (UPDATE jika ada) |
| **SP** | Tidak ada |
| **DML DB** | `UPDATE pasien SET id_posyandu = $1, updated_at = NOW() WHERE id_pasien = $2` |
| | `UPDATE ibu_hamil SET hamil_ke = $1, bulan_mulai_hamil = $2, hpht = $3, status_kehamilan = $4, updated_at = NOW() WHERE id_pasien = $5` |
| | `UPDATE anak SET nama_anak = $1, berat_lahir = $2, panjang_lahir = $3, hubungan_dengan_wali = $4, id_wali = $5, updated_at = NOW() WHERE id_pasien = $6` |

### DELETE /pasien/{id}
| Aspek | Detail |
|--------|--------|
| **Role** | BIDAN |
| **Handler** | `feature/pasien/handler.go` — `Delete` |
| **Input** | Path: `id` (int32, min:1) |
| **Output** | `200` — `{}` |
| **Table WRITE** | `anak`, `ibu_hamil`, `pasien` (soft-delete dalam TRANSACTION) |
| **SP** | Tidak ada |
| **DML DB** | **Transaction:** |
| | `UPDATE anak SET is_deleted = true, deleted_at = NOW() WHERE id_pasien = $1` |
| | `UPDATE ibu_hamil SET is_deleted = true, deleted_at = NOW() WHERE id_pasien = $1` |
| | `UPDATE pasien SET is_deleted = true, deleted_at = NOW() WHERE id_pasien = $1` |

---

## 7. Pemeriksaan

### GET /monitoring/pemeriksaan
| Aspek | Detail |
|--------|--------|
| **Role** | ADMIN |
| **Handler** | `feature/pemeriksaan/handler.go` — `GetAll` |
| **Input** | Query: `page` (int, default:1), `per_page` (int, default:20, max:100), `q` (string) |
| **Output** | `200` — `{ "pemeriksaan": [...], "total_data": ..., "page": ..., "per_page": ... }` |
| | PemeriksaanListItem: `id_hasil_pemeriksaan`, `nama_pasien`, `diinput_oleh`, `status_stunting`, `status_gizi`, `status_verifikasi`, `tanggal_input` |
| **Table READ** | `hasil_pemeriksaan`, `jadwal_imunisasi`, `pasien`, `user_account` (×2 — pasien & petugas), `anak` |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | Raw SQL JOIN: `SELECT hp.*, ... FROM hasil_pemeriksaan hp INNER JOIN jadwal_imunisasi ji ON hp.id_jadwal_imunisasi = ji.id_imunisasi INNER JOIN pasien p ON ji.id_pasien = p.id_pasien INNER JOIN user_account ua ON p.id_pasien = ua.id_user INNER JOIN user_account petugas ON hp.id_petugas_input = petugas.id_user LEFT JOIN anak a ON p.id_pasien = a.id_pasien WHERE ILIKE ... ORDER BY ... LIMIT $1 OFFSET $2` |

### POST /monitoring/pemeriksaan
| Aspek | Detail |
|--------|--------|
| **Role** | ADMIN |
| **Handler** | `feature/pemeriksaan/handler.go` — `Create` |
| **Input** | Body: `CreatePemeriksaanRequest` |
| | `id_jadwal_imunisasi` (int32, min:1) |
| | `berat_badan` (float64, min:0) |
| | `tinggi_badan` (float64, min:0) |
| | `lingkar_kepala` (float64, min:0) |
| | `tekanan_darah` (string, min:1, max:20) |
| | `catatan` (string, max:1000, optional) |
| **Output** | `201` — `{ "id_hasil_pemeriksaan": ..., "status_stunting": "...", "status_gizi": "...", "created_at": "..." }` |
| **Table READ** | `jadwal_imunisasi` (validasi), `user_account` (petugas) |
| **Table WRITE** | `hasil_pemeriksaan` (INSERT) |
| **SP** | Tidak ada |
| **DML DB** | `INSERT INTO hasil_pemeriksaan (id_jadwal_imunisasi, id_petugas_input, berat_badan, tinggi_badan, lingkar_kepala, tekanan_darah, catatan, status_stunting, status_gizi, status_verifikasi, created_at, updated_at) VALUES (...) RETURNING id_hasil_pemeriksaan` |

### GET /monitoring/pemeriksaan/{id}
| Aspek | Detail |
|--------|--------|
| **Role** | USER |
| **Handler** | `feature/pemeriksaan/handler.go` — `GetByID` |
| **Input** | Path: `id` (int32, min:1) |
| | Non-petugas: dicek ownership pasien |
| **Output** | `200` — DetailPemeriksaanResponse: `id_hasil_pemeriksaan`, `pasien` (id, nama), `antropometri` (berat_badan, tinggi_badan, lingkar_kepala, tekanan_darah), `status_kesehatan` (status_stunting, status_gizi), `catatan` |
| **Table READ** | `hasil_pemeriksaan`, `jadwal_imunisasi`, `pasien`, `user_account`, `anak` |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | `SELECT EXISTS(...) FROM hasil_pemeriksaan hp INNER JOIN jadwal_imunisasi ji ... INNER JOIN pasien p ... LEFT JOIN anak a ... WHERE id_hasil_pemeriksaan = $1 AND (p.id_pasien = $2 OR a.id_wali = $2)` — ownership check |
| | Raw SQL JOIN untuk detail: `SELECT hp.*, COALESCE(a.nama_anak, ua.nama) as nama_pasien, ua.nama as nama_petugas FROM hasil_pemeriksaan hp INNER JOIN jadwal_imunisasi ji ... INNER JOIN pasien p ... LEFT JOIN anak a ... INNER JOIN user_account ua ON p.id_pasien = ua.id_user ... WHERE hp.id_hasil_pemeriksaan = $1` |

### PUT /monitoring/pemeriksaan/{id}
| Aspek | Detail |
|--------|--------|
| **Role** | ADMIN |
| **Handler** | `feature/pemeriksaan/handler.go` — `Update` |
| **Input** | Path: `id` (int32, min:1) |
| | Body: `UpdatePemeriksaanRequest` — semua optional: `berat_badan`, `tinggi_badan`, `lingkar_kepala`, `tekanan_darah`, `catatan` |
| **Output** | `200` — `{ "id_hasil_pemeriksaan": ..., "status_gizi_baru": "...", "updated_at": "..." }` |
| **Table WRITE** | `hasil_pemeriksaan` (UPDATE) |
| **SP** | Tidak ada |
| **DML DB** | `UPDATE hasil_pemeriksaan SET berat_badan = $1, tinggi_badan = $2, lingkar_kepala = $3, tekanan_darah = $4, catatan = $5, status_stunting = $6, status_gizi = $7, updated_at = NOW() WHERE id_hasil_pemeriksaan = $8` |

### DELETE /monitoring/pemeriksaan/{id}
| Aspek | Detail |
|--------|--------|
| **Role** | ADMIN |
| **Handler** | `feature/pemeriksaan/handler.go` — `Delete` |
| **Input** | Path: `id` (int32, min:1) |
| **Output** | `200` — `{}` |
| **Table WRITE** | `hasil_pemeriksaan` (**HARD DELETE**) |
| **SP** | Tidak ada |
| **DML DB** | `DELETE FROM hasil_pemeriksaan WHERE id_hasil_pemeriksaan = $1` |

### PATCH /monitoring/pemeriksaan/{id}/verify
| Aspek | Detail |
|--------|--------|
| **Role** | BIDAN |
| **Handler** | `feature/pemeriksaan/handler.go` — `Verify` |
| **Input** | Path: `id` (int32, min:1) |
| | id_petugas dari JWT |
| **Output** | `200` — `{ "id_hasil_pemeriksaan": ..., "diverifikasi_oleh": ..., "status_verifikasi": "Terverifikasi" }` |
| **Table WRITE** | `hasil_pemeriksaan` (UPDATE — tandai terverifikasi) |
| **SP** | Tidak ada |
| **DML DB** | `UPDATE hasil_pemeriksaan SET id_petugas_verifikasi = $1, status_verifikasi = 'Terverifikasi', updated_at = NOW() WHERE id_hasil_pemeriksaan = $2` |

### GET /monitoring/pemeriksaan/pending
| Aspek | Detail |
|--------|--------|
| **Role** | BIDAN |
| **Handler** | `feature/pemeriksaan/handler.go` — `GetPending` |
| **Input** | Query: `page` (int, default:1), `per_page` (int, default:20, max:100) |
| **Output** | `200` — `{ "pemeriksaan_pending": [...], "meta": {...} }` |
| | PendingPemeriksaanItem: `id_hasil_pemeriksaan`, `nama_pasien`, `diinput_oleh`, `tanggal_input` |
| **Table READ** | `hasil_pemeriksaan`, `jadwal_imunisasi`, `pasien`, `user_account`, `anak` (filter: status_verifikasi pending) |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | Raw SQL JOIN dengan filter `status_verifikasi = 'Pending' ORDER BY created_at DESC LIMIT $1 OFFSET $2` |

---

## 8. Imunisasi

### GET /imunisasi
| Aspek | Detail |
|--------|--------|
| **Role** | USER |
| **Handler** | `feature/imunisasi/handler.go` — `GetAll` |
| **Input** | Query: `page` (int, default:1), `per_page` (int, default:20, max:100), `q` (string) |
| | Petugas → semua jadwal; Non-petugas → jadwal milik sendiri |
| **Output** | `200` — `{ "jadwal": [...], "meta": {...} }` |
| | ImunisasiListItem: `id_imunisasi`, `nama_pasien`, `nama_vaksin`, `tanggal_jadwal`, `status_imunisasi` |
| **Table READ** | `jadwal_imunisasi`, `pasien`, `user_account`, `anak` |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | Raw SQL: `SELECT ji.*, ua.nama FROM jadwal_imunisasi ji INNER JOIN pasien p ON ji.id_pasien = p.id_pasien INNER JOIN user_account ua ON p.id_pasien = ua.id_user LEFT JOIN anak a ON p.id_pasien = a.id_pasien WHERE (p.id_pasien = $1 OR a.id_wali = $1) ... ORDER BY ... LIMIT $2 OFFSET $3` (non-petugas) |
| | Tanpa filter user untuk petugas |

### GET /imunisasi/{id}
| Aspek | Detail |
|--------|--------|
| **Role** | USER |
| **Handler** | `feature/imunisasi/handler.go` — `GetByID` |
| **Input** | Path: `id` (int32, min:1) |
| | Non-petugas: dicek ownership pasien |
| **Output** | `200` — ImunisasiDetail: `id_imunisasi`, `id_pasien`, `nama_pasien`, `nama_vaksin`, `tanggal_jadwal`, `tanggal_realisasi`, `status_imunisasi` |
| **Table READ** | `jadwal_imunisasi`, `pasien`, `user_account`, `anak` |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | `SELECT * FROM jadwal_imunisasi WHERE id_imunisasi = $1` |
| | Ownership check: `EXISTS(...) WHERE ji.id_imunisasi = $1 AND (p.id_pasien = $2 OR a.id_wali = $2)` |

### POST /imunisasi
| Aspek | Detail |
|--------|--------|
| **Role** | ADMIN |
| **Handler** | `feature/imunisasi/handler.go` — `Create` |
| **Input** | Body: `CreateImunisasiRequest` |
| | `id_pasien` (int32, min:1) |
| | `nama_vaksin` (string, min:1, max:100) |
| | `tanggal_jadwal` (string, date) |
| **Output** | `201` — `{ "id_imunisasi": ..., "status_imunisasi": "Belum" }` |
| **Table WRITE** | `jadwal_imunisasi` (INSERT) |
| **SP** | Tidak ada |
| **DML DB** | `INSERT INTO jadwal_imunisasi (id_pasien, nama_vaksin, tanggal_jadwal, status_imunisasi, created_at, updated_at) VALUES ($1, $2, $3, 'Belum', NOW(), NOW()) RETURNING id_imunisasi` |

### PUT /imunisasi/{id}
| Aspek | Detail |
|--------|--------|
| **Role** | ADMIN |
| **Handler** | `feature/imunisasi/handler.go` — `Update` |
| **Input** | Path: `id` (int32, min:1) |
| | Body: `UpdateImunisasiRequest` — semua optional: `id_pasien`, `nama_vaksin`, `tanggal_jadwal` |
| **Output** | `200` — `{ "id_imunisasi": ..., "updated_at": "..." }` |
| **Table WRITE** | `jadwal_imunisasi` (UPDATE) |
| **SP** | Tidak ada |
| **DML DB** | `UPDATE jadwal_imunisasi SET id_pasien = $1, nama_vaksin = $2, tanggal_jadwal = $3, updated_at = NOW() WHERE id_imunisasi = $4` |

### DELETE /imunisasi/{id}
| Aspek | Detail |
|--------|--------|
| **Role** | ADMIN |
| **Handler** | `feature/imunisasi/handler.go` — `Delete` |
| **Input** | Path: `id` (int32, min:1) |
| **Output** | `200` — `{}` |
| **Table WRITE** | `jadwal_imunisasi` (**HARD DELETE**) |
| **SP** | Tidak ada |
| **DML DB** | `DELETE FROM jadwal_imunisasi WHERE id_imunisasi = $1` |

### PATCH /imunisasi/{id}/realisasi
| Aspek | Detail |
|--------|--------|
| **Role** | ADMIN |
| **Handler** | `feature/imunisasi/handler.go` — `Realisasi` |
| **Input** | Path: `id` (int32, min:1) |
| | Body: `{ "tanggal_realisasi": "2025-01-15" }` |
| **Output** | `200` — `{ "id_imunisasi": ..., "status_imunisasi": "Sudah", "tanggal_realisasi": "2025-01-15" }` |
| **Table WRITE** | `jadwal_imunisasi` (UPDATE) |
| **SP** | Tidak ada |
| **DML DB** | `UPDATE jadwal_imunisasi SET tanggal_realisasi = $1, status_imunisasi = 'Sudah', updated_at = NOW() WHERE id_imunisasi = $2` |

### GET /imunisasi/pasien/{id_pasien}
| Aspek | Detail |
|--------|--------|
| **Role** | USER |
| **Handler** | `feature/imunisasi/handler.go` — `GetByPasienID` |
| **Input** | Path: `id_pasien` (int32, min:1) |
| **Output** | `200` — `{ "id_pasien": ..., "riwayat_imunisasi": [{ "id_imunisasi", "nama_vaksin", "tanggal_jadwal", "tanggal_realisasi", "status_imunisasi" }] }` |
| **Table READ** | `jadwal_imunisasi`, `anak` (ownership check) |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | `SELECT * FROM jadwal_imunisasi WHERE id_pasien = $1 ORDER BY tanggal_jadwal ASC` |

### GET /imunisasi/statistik
| Aspek | Detail |
|--------|--------|
| **Role** | ADMIN |
| **Handler** | `feature/imunisasi/handler.go` — `GetStatistik` |
| **Input** | — |
| **Output** | `200` — `{ "total_target_imunisasi": ..., "total_terealisasi": ..., "cakupan_persentase": ..., "vaksin_terbanyak": "..." }` |
| **Table READ** | `jadwal_imunisasi` |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | `SELECT COUNT(*), COUNT(*) FILTER (WHERE status_imunisasi = 'Sudah') FROM jadwal_imunisasi` |
| | `SELECT nama_vaksin FROM jadwal_imunisasi GROUP BY nama_vaksin ORDER BY COUNT(*) DESC LIMIT 1` |

---

## 9. Artikel

### GET /artikel
| Aspek | Detail |
|--------|--------|
| **Role** | Public |
| **Handler** | `feature/artikel/handler.go` — `GetAllPublished` |
| **Input** | Query: `page` (int, default:1), `per_page` (int, default:20, max:100) |
| **Output** | `200` — `{ "artikel": [...], "meta": {...} }` |
| | ArtikelListItem: `id_artikel`, `judul`, `kategori`, `ringkasan`, `nama_penulis`, `tanggal_publish`, `status_artikel` |
| **Table READ** | `artikel`, `user_account` (penulis) |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | `SELECT a.*, ua.nama FROM artikel a INNER JOIN user_account ua ON a.id_penulis = ua.id_user WHERE status_artikel = 'Dipublikasikan' ORDER BY tanggal_publish DESC LIMIT $1 OFFSET $2` |
| | `SELECT COUNT(*) FROM artikel WHERE status_artikel = 'Dipublikasikan'` |

### GET /artikel/{id}
| Aspek | Detail |
|--------|--------|
| **Role** | Public |
| **Handler** | `feature/artikel/handler.go` — `GetByID` |
| **Input** | Path: `id` (int32, min:1) |
| **Output** | `200` — ArtikelDetail: `id_artikel`, `judul`, `isi_artikel`, `kategori`, `nama_penulis`, `nama_verifikator`, `tanggal_publish`, `created_at`, `updated_at` |
| **Table READ** | `artikel`, `user_account` (penulis), `user_account` (verifikator, LEFT JOIN) |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | `SELECT a.*, penulis.nama, verifikator.nama FROM artikel a INNER JOIN user_account penulis ON a.id_penulis = penulis.id_user LEFT JOIN user_account verifikator ON a.id_verifikator = verifikator.id_user WHERE a.id_artikel = $1` |

### GET /artikel/semua
| Aspek | Detail |
|--------|--------|
| **Role** | ADMIN |
| **Handler** | `feature/artikel/handler.go` — `GetAll` |
| **Input** | Query: `page` (int, default:1), `per_page` (int, default:20, max:100) |
| **Output** | `200` — `{ "artikel": [...], "meta": {...} }` — semua status artikel |
| **Table READ** | `artikel`, `user_account` |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | Sama dengan GetAllPublished tapi tanpa filter `status_artikel = 'Dipublikasikan'` |

### POST /artikel
| Aspek | Detail |
|--------|--------|
| **Role** | ADMIN |
| **Handler** | `feature/artikel/handler.go` — `Create` |
| **Input** | Body: `CreateArtikelRequest` |
| | `judul` (string, min:1, max:255) |
| | `isi_artikel` (string, min:1) |
| | `kategori` (string, max:100, optional) |
| **Output** | `201` — `{ "id_artikel": ..., "status_artikel": "..." }` |
| | Dinkes → langsung "Dipublikasikan"; non-Dinkes → "Menunggu Verifikasi" |
| **Table WRITE** | `artikel` (INSERT) |
| **SP** | Tidak ada |
| **DML DB** | `INSERT INTO artikel (judul, isi_artikel, kategori, id_penulis, status_artikel, created_at, updated_at) VALUES (...) RETURNING id_artikel` |

### PATCH /artikel/{id}
| Aspek | Detail |
|--------|--------|
| **Role** | ADMIN |
| **Handler** | `feature/artikel/handler.go` — `Update` |
| **Input** | Path: `id` (int32, min:1) |
| | Body: `UpdateArtikelRequest` — semua optional: `judul`, `isi_artikel`, `kategori` |
| **Output** | `200` — ArtikelDetail |
| **Table WRITE** | `artikel` (UPDATE) |
| **SP** | Tidak ada |
| **DML DB** | `UPDATE artikel SET judul = $1, isi_artikel = $2, kategori = $3, updated_at = NOW() WHERE id_artikel = $4` |

### DELETE /artikel/{id}
| Aspek | Detail |
|--------|--------|
| **Role** | DINKES |
| **Handler** | `feature/artikel/handler.go` — `Delete` |
| **Input** | Path: `id` (int32, min:1) |
| **Output** | `200` — `{}` |
| **Table WRITE** | `artikel` (**HARD DELETE**) |
| **SP** | Tidak ada |
| **DML DB** | `DELETE FROM artikel WHERE id_artikel = $1` |

### PATCH /artikel/{id}/review
| Aspek | Detail |
|--------|--------|
| **Role** | DINKES |
| **Handler** | `feature/artikel/handler.go` — `Review` |
| **Input** | Path: `id` (int32, min:1) |
| | Body: `{ "aksi": "setujui|tolak", "catatan_review": "..." }` |
| **Output** | `200` — `{ "id_artikel": ..., "status_artikel": "...", "tanggal_publish": "..." }` |
| **Table WRITE** | `artikel` (UPDATE — review) |
| **SP** | Tidak ada |
| **DML DB** | Setujui: `UPDATE artikel SET status_artikel = 'Dipublikasikan', id_verifikator = $1, tanggal_publish = NOW(), updated_at = NOW() WHERE id_artikel = $2 RETURNING *` |
| | Tolak: `UPDATE artikel SET status_artikel = 'Ditolak', id_verifikator = $1, updated_at = NOW() WHERE id_artikel = $2 RETURNING *` |

### GET /artikel/pending
| Aspek | Detail |
|--------|--------|
| **Role** | DINKES |
| **Handler** | `feature/artikel/handler.go` — `GetPending` |
| **Input** | Query: `page` (int, default:1), `per_page` (int, default:20, max:100) |
| **Output** | `200` — `{ "artikel": [...], "meta": {...} }` |
| | ArtikelPendingItem: `id_artikel`, `judul`, `nama_penulis`, `created_at`, `status_artikel` |
| **Table READ** | `artikel`, `user_account` (filter: `status_artikel = 'Menunggu Verifikasi'`) |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | `SELECT a.*, ua.nama FROM artikel a INNER JOIN user_account ua ON a.id_penulis = ua.id_user WHERE status_artikel = 'Menunggu Verifikasi' ORDER BY created_at DESC LIMIT $1 OFFSET $2` |

---

## 10. Tindak Lanjut

### GET /tindak-lanjut/status
| Aspek | Detail |
|--------|--------|
| **Role** | Authenticated (authAccess) |
| **Handler** | `feature/tindaklanjut/handler.go` — `GetStatusTindakLanjut` |
| **Input** | Query: `page` (int, default:1), `per_page` (int, default:20, max:100) |
| **Output** | `200` — `{ "pasien": [...], "meta": {...} }` |
| | StatusTindakLanjutItem: `id_pasien`, `nama_pasien`, `status_pasien`, `status_rujukan`, `tanggal_rujukan`, `tanggal_deadline` (tanggal_rujukan + 7 hari), `status_deadline` (terlambat/mendekati/aman), `faskes`, `alasan_rujukan`, `jenis_tindakan` |
| **Table READ** | `tindak_lanjut`, `hasil_pemeriksaan`, `jadwal_imunisasi`, `pasien`, `user_account`, `rujukan`, `fasilitas_kesehatan` |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | Raw SQL dengan CASE: `SELECT ..., r.tanggal_rujukan, r.tanggal_rujukan + interval '7 days' as tanggal_deadline, CASE WHEN (r.tanggal_rujukan + interval '7 days') < NOW() THEN 'terlambat' WHEN ... END as status_deadline FROM tindak_lanjut tl LEFT JOIN hasil_pemeriksaan hp ... LEFT JOIN jadwal_imunisasi ji ... LEFT JOIN pasien p ... LEFT JOIN user_account ua ... LEFT JOIN rujukan r ON tl.id_tindak_lanjut = r.id_tindak_lanjut LEFT JOIN fasilitas_kesehatan fk ON r.id_faskes = fk.id_faskes ORDER BY ... LIMIT $1 OFFSET $2` |

### GET /tindak-lanjut/pasien
| Aspek | Detail |
|--------|--------|
| **Role** | BIDAN |
| **Handler** | `feature/tindaklanjut/handler.go` — `GetPasienTindakLanjut` |
| **Input** | Query: `page` (int, default:1), `per_page` (int, default:20, max:100) |
| **Output** | `200` — `{ "pasien": [...], "meta": {...} }` |
| | PasienTindakLanjutItem: `id_pasien`, `nama_pasien`, `status_gizi`, `status_pasien`, `tanggal_pemeriksaan` |
| **Table READ** | `pasien`, `user_account`, `anak`, `hasil_pemeriksaan` (LEFT JOIN LATERAL), `tindak_lanjut` (cek existing) |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | `SELECT p.*, ua.nama, hp.status_gizi FROM pasien p INNER JOIN user_account ua ... LEFT JOIN anak a ... LEFT JOIN LATERAL (SELECT * FROM hasil_pemeriksaan WHERE id_jadwal_imunisasi.pasien = p.id_pasien ORDER BY created_at DESC LIMIT 1) hp ON true WHERE hp.status_gizi IS NOT NULL AND hp.status_gizi != 'Gizi Baik' AND tl.id_tindak_lanjut IS NULL ... LIMIT $1 OFFSET $2` |

### GET /tindak-lanjut/pasien/{id}
| Aspek | Detail |
|--------|--------|
| **Role** | BIDAN |
| **Handler** | `feature/tindaklanjut/handler.go` — `GetDetailPasienByID` |
| **Input** | Path: `id` (int32, min:1) |
| **Output** | `200` — DetailPasienTindakLanjut: `id_pasien`, `nama_pasien`, `usia` (EXTRACT YEAR FROM AGE), `hasil_monitoring_terakhir` (status_gizi, status_stunting, catatan), `riwayat_pemeriksaan[]` (tanggal, berat_badan, tinggi_badan) |
| **Table READ** | `pasien`, `user_account`, `anak`, `hasil_pemeriksaan` (LATERAL), `jadwal_imunisasi` |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | `SELECT p.*, ua.*, EXTRACT(YEAR FROM AGE(ua.tanggal_lahir)) as usia FROM pasien p INNER JOIN user_account ua ...` |
| | Riwayat: `SELECT ji.tanggal_jadwal, hp.berat_badan, hp.tinggi_badan FROM hasil_pemeriksaan hp INNER JOIN jadwal_imunisasi ji ON hp.id_jadwal_imunisasi = ji.id_imunisasi WHERE ji.id_pasien = $1 ORDER BY ji.tanggal_jadwal` |

### POST /tindak-lanjut
| Aspek | Detail |
|--------|--------|
| **Role** | BIDAN |
| **Handler** | `feature/tindaklanjut/handler.go` — `CreateTindakLanjut` |
| **Input** | Body: `CreateTindakLanjutRequest` |
| | `id_hasil_pemeriksaan` (int32) |
| | `jenis_tindakan` (enum: Rujukan, Kontrol Ulang) |
| | `catatan_medis` (string, max:1000, optional) |
| | `rekomendasi` (string, max:1000, optional) |
| | `jadwal_kontrol` (string, optional) |
| | `alasan_rujukan` (string, max:1000, optional) |
| | `id_faskes` (int32, optional — untuk rujukan) |
| **Output** | `201` — `{ "id_tindak_lanjut": ..., "id_rujukan": ..., "status_pasien": "..." }` |
| **Table READ** | `hasil_pemeriksaan` (validasi), `jadwal_imunisasi` (get id_pasien) |
| **Table WRITE** | `tindak_lanjut` (INSERT), `rujukan` (INSERT jika jenis=Rujukan) |
| **SP** | Tidak ada |
| **DML DB** | `INSERT INTO tindak_lanjut (id_hasil_pemeriksaan, id_bidan, jenis_tindakan, catatan_medis, rekomendasi, jadwal_kontrol, status_pasien, created_at, updated_at) VALUES (...) RETURNING id_tindak_lanjut` |
| | Jika Rujukan: `INSERT INTO rujukan (id_tindak_lanjut, alasan_rujukan, id_faskes, status_rujukan, tanggal_rujukan, created_at, updated_at) VALUES (...) RETURNING id_rujukan` |

### PATCH /rujukan/{id}/status
| Aspek | Detail |
|--------|--------|
| **Role** | BIDAN |
| **Handler** | `feature/tindaklanjut/handler.go` — `UpdateStatusRujukan` |
| **Input** | Path: `id` (int32, min:1) |
| | Body: `{ "status_rujukan": "Diajukan|Diproses|Diterima|Ditolak|Selesai" }` |
| **Output** | `200` — `{ "id_rujukan": ..., "status_rujukan": "..." }` |
| **Table WRITE** | `rujukan` (UPDATE) |
| **SP** | Tidak ada |
| **DML DB** | `UPDATE rujukan SET status_rujukan = $1, updated_at = NOW() WHERE id_rujukan = $2 RETURNING *` |

### GET /laporan/tindak-lanjut
| Aspek | Detail |
|--------|--------|
| **Role** | DINKES |
| **Handler** | `feature/tindaklanjut/handler.go` — `GetLaporanTindakLanjut` |
| **Input** | — |
| **Output** | `200` — `{ "laporan": [...], "total_data": ... }` |
| | LaporanTindakLanjutItem: `wilayah`, `jumlah_pasien_dirujuk`, `jumlah_pasien_diterima`, `jumlah_pasien_diproses` |
| **Table READ** | `rujukan`, `tindak_lanjut`, `hasil_pemeriksaan`, `jadwal_imunisasi`, `pasien`, `posyandu`, `lokasi` |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | `SELECT l.nama_lokasi as wilayah, COUNT(DISTINCT r.id_rujukan) FILTER (WHERE r.status_rujukan = 'Diajukan') as jumlah_pasien_dirujuk, COUNT(DISTINCT r.id_rujukan) FILTER (WHERE r.status_rujukan = 'Diterima') as jumlah_pasien_diterima, COUNT(DISTINCT r.id_rujukan) FILTER (WHERE r.status_rujukan = 'Diproses') as jumlah_pasien_diproses FROM rujukan r INNER JOIN tindak_lanjut tl ... INNER JOIN hasil_pemeriksaan hp ... INNER JOIN jadwal_imunisasi ji ... INNER JOIN pasien p ... INNER JOIN posyandu pos ... INNER JOIN lokasi l ON pos.id_lokasi = l.id_lokasi GROUP BY l.nama_lokasi` |

### GET /tindak-lanjut/{id}
| Aspek | Detail |
|--------|--------|
| **Role** | USER |
| **Handler** | `feature/tindaklanjut/handler.go` — `GetDetailTindakLanjutByID` |
| **Input** | Path: `id` (int32, min:1) — id_tindak_lanjut |
| **Output** | `200` — DetailTindakLanjutPasien: `id_tindak_lanjut`, `status_pasien`, `catatan_medis`, `rekomendasi`, `jadwal_kontrol`, `status_rujukan`, `nama_faskes` |
| **Table READ** | `tindak_lanjut`, `rujukan` (LEFT JOIN), `fasilitas_kesehatan` (LEFT JOIN) |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | `SELECT tl.*, r.status_rujukan, fk.nama_faskes FROM tindak_lanjut tl LEFT JOIN rujukan r ON tl.id_tindak_lanjut = r.id_tindak_lanjut LEFT JOIN fasilitas_kesehatan fk ON r.id_faskes = fk.id_faskes WHERE tl.id_tindak_lanjut = $1` |

---

## 11. Dashboard

### GET /dashboard/stats
| Aspek | Detail |
|--------|--------|
| **Role** | USER |
| **Handler** | `feature/dashboard/handler.go` — `GetDashboardStats` |
| **Input** | — |
| **Output** | `200` — DashboardStatsResponse: `total_pasien`, `perlu_verifikasi`, `tindak_lanjut`, `kasus_stunting`, `jadwal_posyandu`, `total_balita`, `cakupan_persentase` |
| **Table READ** | `mv_dashboard_stats` (materialized view) |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | `SELECT total_pasien, perlu_verifikasi, tindak_lanjut, kasus_stunting, jadwal_posyandu, total_balita, cakupan_persentase FROM mv_dashboard_stats` |

### GET /dashboard/distribusi-gizi
| Aspek | Detail |
|--------|--------|
| **Role** | USER |
| **Handler** | `feature/dashboard/handler.go` — `GetDistribusiGizi` |
| **Input** | — |
| **Output** | `200` — `{ "distribusi": [{ "status_gizi": "...", "jumlah": ... }] }` |
| **Table READ** | `mv_dashboard_distribusi_gizi` |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | `SELECT status_gizi, jumlah FROM mv_dashboard_distribusi_gizi ORDER BY jumlah DESC` |

### GET /dashboard/tren-stunting
| Aspek | Detail |
|--------|--------|
| **Role** | USER |
| **Handler** | `feature/dashboard/handler.go` — `GetTrenStunting` |
| **Input** | — |
| **Output** | `200` — `{ "tren": [{ "bulan": "...", "jumlah": ... }] }` |
| **Table READ** | `mv_dashboard_tren_stunting` |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | `SELECT bulan, jumlah FROM mv_dashboard_tren_stunting ORDER BY bulan` |

### GET /dashboard/stunting-per-wilayah
| Aspek | Detail |
|--------|--------|
| **Role** | USER |
| **Handler** | `feature/dashboard/handler.go` — `GetStuntingPerWilayah` |
| **Input** | — |
| **Output** | `200` — `{ "wilayah": [{ "nama_wilayah", "prevalensi", "jumlah_kasus", "total_balita", "level" }] }` |
| **Table READ** | `mv_dashboard_stunting_per_wilayah` |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | `SELECT nama_wilayah, prevalensi, jumlah_kasus, total_balita, level FROM mv_dashboard_stunting_per_wilayah ORDER BY jumlah_kasus DESC` |

### GET /dashboard/kehadiran-bulanan
| Aspek | Detail |
|--------|--------|
| **Role** | USER |
| **Handler** | `feature/dashboard/handler.go` — `GetKehadiranBulanan` |
| **Input** | — |
| **Output** | `200` — `{ "tren": [{ "bulan": "...", "jumlah": ... }] }` |
| **Table READ** | `mv_dashboard_kehadiran_bulanan` |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | `SELECT bulan, jumlah FROM mv_dashboard_kehadiran_bulanan ORDER BY bulan` |

### GET /dashboard/jadwal-terdekat
| Aspek | Detail |
|--------|--------|
| **Role** | USER |
| **Handler** | `feature/dashboard/handler.go` — `GetJadwalTerdekat` |
| **Input** | — |
| **Output** | `200` — `{ "jadwal": [{ "id", "nama_vaksin", "tanggal_jadwal", "nama_pasien" }] }` (LIMIT 10) |
| **Table READ** | `mv_dashboard_jadwal_terdekat` |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | `SELECT id, nama_vaksin, tanggal_jadwal::text, nama_pasien FROM mv_dashboard_jadwal_terdekat` |

### GET /stats
| Aspek | Detail |
|--------|--------|
| **Role** | Public |
| **Handler** | `feature/dashboard/handler.go` — `GetPublicStats` |
| **Input** | — |
| **Output** | `200` — `{ "total_pasien", "balita_dipantau", "kasus_stunting", "total_artikel" }` |
| **Table READ** | `mv_public_stats` |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | `SELECT total_pasien, balita_dipantau, kasus_stunting, total_artikel FROM mv_public_stats` |

### GET /monitoring/pasien/{id}/riwayat-pemeriksaan
| Aspek | Detail |
|--------|--------|
| **Role** | USER |
| **Handler** | `feature/dashboard/handler.go` — `GetRiwayat` |
| **Input** | Path: `id` (int32, min:1) — id_pasien |
| **Output** | `200` — `{ "riwayat": [{ "tanggal", "berat_badan", "tinggi_badan", "status_gizi", "catatan", "petugas" }] }` |
| **Table READ** | `mv_riwayat_pemeriksaan` |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | `SELECT tanggal::text, berat_badan, tinggi_badan, status_gizi, catatan, petugas FROM mv_riwayat_pemeriksaan WHERE id_pasien = $1 ORDER BY tanggal DESC` |

### GET /monitoring/pasien/{id}/tumbuh-kembang
| Aspek | Detail |
|--------|--------|
| **Role** | USER |
| **Handler** | `feature/dashboard/handler.go` — `GetTumbuhKembang` |
| **Input** | Path: `id` (int32, min:1) — id_pasien |
| **Output** | `200` — `{ "data": [{ "bulan", "berat_badan", "tinggi_badan" }] }` |
| **Table READ** | `mv_tumbuh_kembang` |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | `SELECT bulan, berat_badan, tinggi_badan FROM mv_tumbuh_kembang WHERE id_pasien = $1 ORDER BY bulan` |

### GET /dashboard/ibu-hamil-stats
| Aspek | Detail |
|--------|--------|
| **Role** | USER |
| **Handler** | `feature/dashboard/handler.go` — `GetIbuHamilStats` |
| **Input** | — |
| **Output** | `200` — `{ "total_ibu_hamil", "trimester_1", "trimester_2", "trimester_3", "melahirkan", "nifas", "keguguran" }` |
| **Table READ** | `mv_ibu_hamil_stats` |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | `SELECT total_ibu_hamil, trimester_1, trimester_2, trimester_3, melahirkan, nifas, keguguran FROM mv_ibu_hamil_stats` |

### GET /dashboard/ibu-hamil-per-wilayah
| Aspek | Detail |
|--------|--------|
| **Role** | USER |
| **Handler** | `feature/dashboard/handler.go` — `GetIbuHamilPerWilayah` |
| **Input** | — |
| **Output** | `200` — `{ "wilayah": [{ "nama_wilayah", "total_ibu_hamil", "trimester_1", "trimester_2", "trimester_3", "melahirkan", "nifas", "keguguran" }] }` |
| **Table READ** | `mv_ibu_hamil_per_wilayah` |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | `SELECT * FROM mv_ibu_hamil_per_wilayah ORDER BY total_ibu_hamil DESC` |

### GET /monitoring/semua-pemeriksaan
| Aspek | Detail |
|--------|--------|
| **Role** | ADMIN |
| **Handler** | `feature/dashboard/handler.go` — `GetSemuaPemeriksaan` |
| **Input** | Query: `page` (int, default:1), `per_page` (int, default:20, max:50), `id_bidan` (int32, optional), `id_posyandu` (int32, optional), `id_kader` (int32, optional) |
| **Output** | `200` — `{ "pemeriksaan": [...], "meta": {...} }` |
| | PemeriksaanItem: `id_hasil_pemeriksaan`, `id_jadwal_imunisasi`, `nama_vaksin`, `nama_pasien`, `berat_badan`, `tinggi_badan`, `lingkar_kepala`, `tekanan_darah`, `status_stunting`, `status_gizi`, `catatan`, `tanggal`, `petugas` |
| **Table READ** | `hasil_pemeriksaan`, `jadwal_imunisasi`, `pasien`, `posyandu`, `user_account` (×2), `anak`, `kader_posyandu` (untuk filter id_posyandu dari id_kader) |
| **Table WRITE** | — |
| **SP** | Tidak ada |
| **DML DB** | Raw SQL JOIN dengan WHERE dinamis: kader→cari id_posyandu, filter by id_bidan dan/atau id_posyandu. ORDER BY created_at DESC LIMIT $1 OFFSET $2 |

---

## Ringkasan Statistik

| Domain | Jumlah Endpoint | Table READ | Table WRITE |
|--------|----------------|------------|-------------|
| Auth & Session | 6 | 9 | 4 |
| Lokasi | 1 | 1 | 0 |
| Faskes | 1 | 1 | 0 |
| User Account | 7 | 9 | 5 |
| Notifikasi | 7 | 9 | 1 |
| Pasien | 7 | 5 | 3 |
| Pemeriksaan | 7 | 5 | 1 |
| Imunisasi | 8 | 5 | 1 |
| Artikel | 7 | 2 | 1 |
| Tindak Lanjut | 7 | 10 | 2 |
| Dashboard | 12 | 11 MV + 7 tables | 0 |
| **Total** | **~70** | **21 tables + 11 MVs** | **11 tables** |

> **Catatan:** Tidak ada Stored Procedure (SP) dalam database ini.  
> Semua logic bisnis berada di application layer (Go service), bukan di database.
