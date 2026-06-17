-- +goose Up
-- ============================================================
-- SEED DATA MASTER — bulk insert from CSV files
-- Urutan FK-safe agar constraint terpenuhi
-- CSVs di-mount dari host ./migrations/data/ ke /data-csv
-- ============================================================

-- ============================================================
-- 1. PENDIDIKAN
-- ============================================================
COPY pendidikan (id_pendidikan, nama_pendidikan, jenjang, lama_tahun)
FROM '/data-csv/pendidikan.csv' CSV HEADER;

-- ============================================================
-- 2. PEKERJAAN
-- ============================================================
COPY pekerjaan (id_pekerjaan, nama_pekerjaan, sektor)
FROM '/data-csv/pekerjaan.csv' CSV HEADER;

-- ============================================================
-- 3. KATEGORI PENDAPATAN
-- ============================================================
COPY kategori_pendapatan (id_pendapatan, kategori_pendapatan, pendapatan_min, pendapatan_max)
FROM '/data-csv/kategori_pendapatan.csv' CSV HEADER;

-- ============================================================
-- 4. LOKASI (self-referencing FK bagian_dari)
-- ============================================================
COPY lokasi (id_lokasi, nama_lokasi, tipe_lokasi, bagian_dari)
FROM '/data-csv/lokasi.csv' CSV HEADER;

-- ============================================================
-- 5. USER ACCOUNT (FK ke lokasi, pendidikan, pekerjaan, pendapatan)
-- ============================================================
COPY user_account (
    id_user, email, password, no_hp, status_verifikasi, nama, nik,
    jenis_kelamin, tanggal_lahir, id_lokasi, id_pendidikan,
    id_pekerjaan, id_pendapatan, jumlah_tanggungan, akun_ke,
    created_at, updated_at
)
FROM '/data-csv/user_account.csv' CSV HEADER;

-- ============================================================
-- 6. FASILITAS KESEHATAN (FK ke lokasi)
-- ============================================================
COPY fasilitas_kesehatan (id_faskes, nama_faskes, tipe_faskes, id_lokasi, created_at, updated_at)
FROM '/data-csv/fasilitas_kesehatan.csv' CSV HEADER;

-- ============================================================
-- 7. DINAS KESEHATAN (FK ke user_account)
-- ============================================================
COPY dinas_kesehatan (id_user, created_at, updated_at)
FROM '/data-csv/dinas_kesehatan.csv' CSV HEADER;

-- ============================================================
-- 8. BIDAN (FK ke user_account, lokasi)
-- ============================================================
COPY bidan (id_user, no_str, wilayah_kerja, created_at, updated_at)
FROM '/data-csv/bidan.csv' CSV HEADER;

-- ============================================================
-- 9. POSYANDU (FK ke lokasi, bidan)
-- ============================================================
COPY posyandu (id_posyandu, nama_posyandu, id_lokasi, id_bidan, created_at, updated_at)
FROM '/data-csv/posyandu.csv' CSV HEADER;

-- ============================================================
-- 10. KADER POSYANDU (FK ke user_account, posyandu)
-- ============================================================
COPY kader_posyandu (id_user, no_sk, id_posyandu, created_at, updated_at)
FROM '/data-csv/kader_posyandu.csv' CSV HEADER;

-- ============================================================
-- 11. PASIEN (FK ke user_account, posyandu)
-- ============================================================
COPY pasien (id_pasien, id_posyandu, created_at, updated_at)
FROM '/data-csv/pasien.csv' CSV HEADER;

-- ============================================================
-- 12. IBU HAMIL (FK ke pasien)
-- ============================================================
COPY ibu_hamil (id_ibu_hamil, id_pasien, hamil_ke, bulan_mulai_hamil, hpht, status_kehamilan, created_at, updated_at)
FROM '/data-csv/ibu_hamil.csv' CSV HEADER;

-- ============================================================
-- 13. ANAK (FK ke pasien, ibu_hamil, user_account)
-- ============================================================
COPY anak (id_pasien, id_ibu_hamil, id_wali, nama_anak, berat_lahir, panjang_lahir, hubungan_dengan_wali, created_at, updated_at)
FROM '/data-csv/anak.csv' CSV HEADER;

-- ============================================================
-- 14. JADWAL IMUNISASI (FK ke pasien)
-- ============================================================
COPY jadwal_imunisasi (id_imunisasi, id_pasien, nama_vaksin, tanggal_jadwal, tanggal_realisasi, status_imunisasi, created_at, updated_at)
FROM '/data-csv/jadwal_imunisasi.csv' CSV HEADER;

-- ============================================================
-- 15. ARTIKEL (FK ke user_account)
-- ============================================================
COPY artikel (id_artikel, judul, isi_artikel, kategori, status_artikel, id_penulis, id_verifikator, tanggal_publish, created_at, updated_at)
FROM '/data-csv/artikel.csv' CSV HEADER;

-- ============================================================
-- 16. HASIL PEMERIKSAAN (FK ke user_account, jadwal_imunisasi)
-- ============================================================
COPY hasil_pemeriksaan (id_hasil_pemeriksaan, id_petugas_input, id_jadwal_imunisasi, berat_badan, tinggi_badan, lingkar_kepala, tekanan_darah, status_stunting, status_gizi, catatan, created_at, updated_at)
FROM '/data-csv/hasil_pemeriksaan.csv' CSV HEADER;

-- ============================================================
-- 17. TINDAK LANJUT (FK ke hasil_pemeriksaan, bidan)
-- ============================================================
COPY tindak_lanjut (id_tindak_lanjut, id_hasil_pemeriksaan, id_bidan, catatan_medis, rekomendasi, jadwal_kontrol, status_pasien, created_at, updated_at)
FROM '/data-csv/tindak_lanjut.csv' CSV HEADER;

-- ============================================================
-- 18. RUJUKAN (FK ke tindak_lanjut, fasilitas_kesehatan)
-- ============================================================
COPY rujukan (id_rujukan, id_tindak_lanjut, alasan_rujukan, tanggal_rujukan, status_rujukan, id_faskes, created_at, updated_at)
FROM '/data-csv/rujukan.csv' CSV HEADER;

-- ============================================================
-- 19. NOTIFIKASI (FK ke user_account)
-- ============================================================
COPY notifikasi (id_notifikasi, id_user, judul, pesan, tipe_notifikasi, status_baca, tanggal_kirim)
FROM '/data-csv/notifikasi.csv' CSV HEADER;

-- ============================================================
-- 20. USER SESSION (FK ke user_account, UUID id_session)
-- ============================================================
COPY user_session (id_session, id_user, status_session, ip_address, expired_at, created_at, updated_at)
FROM '/data-csv/user_session.csv' CSV HEADER;

-- ============================================================
-- 21. AUDIT LOG (FK ke user_account, user_session)
-- ============================================================
COPY audit_log (id_log, tipe_aktor, id_user, id_user_session, tipe_aktivitas, berhasil, endpoint, table_name, record_id, old_value, new_value, detail, ip_address, user_agent, waktu_aktivitas)
FROM '/data-csv/audit_log.csv' CSV HEADER;

-- ============================================================
-- RESET ALL SEQUENCES after explicit ID inserts
-- ============================================================
SELECT setval('pendidikan_id_pendidikan_seq',               COALESCE((SELECT MAX(id_pendidikan)               FROM pendidikan),               1));
SELECT setval('pekerjaan_id_pekerjaan_seq',                 COALESCE((SELECT MAX(id_pekerjaan)                 FROM pekerjaan),                 1));
SELECT setval('kategori_pendapatan_id_pendapatan_seq',      COALESCE((SELECT MAX(id_pendapatan)                FROM kategori_pendapatan),        1));
SELECT setval('lokasi_id_lokasi_seq',                       COALESCE((SELECT MAX(id_lokasi)                    FROM lokasi),                    1));
SELECT setval('user_account_id_user_seq',                   COALESCE((SELECT MAX(id_user)                      FROM user_account),              1));
SELECT setval('fasilitas_kesehatan_id_faskes_seq',          COALESCE((SELECT MAX(id_faskes)                    FROM fasilitas_kesehatan),        1));
SELECT setval('posyandu_id_posyandu_seq',                   COALESCE((SELECT MAX(id_posyandu)                  FROM posyandu),                   1));
SELECT setval('ibu_hamil_id_ibu_hamil_seq',                 COALESCE((SELECT MAX(id_ibu_hamil)                 FROM ibu_hamil),                  1));
SELECT setval('jadwal_imunisasi_id_imunisasi_seq',          COALESCE((SELECT MAX(id_imunisasi)                 FROM jadwal_imunisasi),           1));
SELECT setval('artikel_id_artikel_seq',                     COALESCE((SELECT MAX(id_artikel)                   FROM artikel),                    1));
SELECT setval('hasil_pemeriksaan_id_hasil_pemeriksaan_seq', COALESCE((SELECT MAX(id_hasil_pemeriksaan)         FROM hasil_pemeriksaan),          1));
SELECT setval('tindak_lanjut_id_tindak_lanjut_seq',         COALESCE((SELECT MAX(id_tindak_lanjut)             FROM tindak_lanjut),              1));
SELECT setval('rujukan_id_rujukan_seq',                     COALESCE((SELECT MAX(id_rujukan)                   FROM rujukan),                    1));
SELECT setval('notifikasi_id_notifikasi_seq',               COALESCE((SELECT MAX(id_notifikasi)                FROM notifikasi),                 1));
SELECT setval('audit_log_id_log_seq',                       COALESCE((SELECT MAX(id_log)                       FROM audit_log),                  1));

-- ============================================================
-- VERIFICATION
-- ============================================================
SELECT 'seed_master' AS stage, COUNT(*) AS total_rows FROM (
    SELECT 'pendidikan' AS tbl, COUNT(*) FROM pendidikan UNION ALL
    SELECT 'pekerjaan', COUNT(*) FROM pekerjaan UNION ALL
    SELECT 'kategori_pendapatan', COUNT(*) FROM kategori_pendapatan UNION ALL
    SELECT 'lokasi', COUNT(*) FROM lokasi UNION ALL
    SELECT 'user_account', COUNT(*) FROM user_account UNION ALL
    SELECT 'fasilitas_kesehatan', COUNT(*) FROM fasilitas_kesehatan UNION ALL
    SELECT 'dinas_kesehatan', COUNT(*) FROM dinas_kesehatan UNION ALL
    SELECT 'bidan', COUNT(*) FROM bidan UNION ALL
    SELECT 'posyandu', COUNT(*) FROM posyandu UNION ALL
    SELECT 'kader_posyandu', COUNT(*) FROM kader_posyandu UNION ALL
    SELECT 'pasien', COUNT(*) FROM pasien UNION ALL
    SELECT 'ibu_hamil', COUNT(*) FROM ibu_hamil UNION ALL
    SELECT 'anak', COUNT(*) FROM anak UNION ALL
    SELECT 'jadwal_imunisasi', COUNT(*) FROM jadwal_imunisasi UNION ALL
    SELECT 'artikel', COUNT(*) FROM artikel UNION ALL
    SELECT 'hasil_pemeriksaan', COUNT(*) FROM hasil_pemeriksaan UNION ALL
    SELECT 'tindak_lanjut', COUNT(*) FROM tindak_lanjut UNION ALL
    SELECT 'rujukan', COUNT(*) FROM rujukan UNION ALL
    SELECT 'notifikasi', COUNT(*) FROM notifikasi UNION ALL
    SELECT 'user_session', COUNT(*) FROM user_session UNION ALL
    SELECT 'audit_log', COUNT(*) FROM audit_log
) AS counts;

-- +goose Down

-- ============================================================
-- ROLLBACK — truncate in reverse FK order
-- ============================================================
TRUNCATE TABLE audit_log CASCADE;
TRUNCATE TABLE user_session CASCADE;
TRUNCATE TABLE notifikasi CASCADE;
TRUNCATE TABLE rujukan CASCADE;
TRUNCATE TABLE tindak_lanjut CASCADE;
TRUNCATE TABLE hasil_pemeriksaan CASCADE;
TRUNCATE TABLE artikel CASCADE;
TRUNCATE TABLE jadwal_imunisasi CASCADE;
TRUNCATE TABLE anak CASCADE;
TRUNCATE TABLE ibu_hamil CASCADE;
TRUNCATE TABLE pasien CASCADE;
TRUNCATE TABLE kader_posyandu CASCADE;
TRUNCATE TABLE posyandu CASCADE;
TRUNCATE TABLE bidan CASCADE;
TRUNCATE TABLE dinas_kesehatan CASCADE;
TRUNCATE TABLE fasilitas_kesehatan CASCADE;
TRUNCATE TABLE user_account CASCADE;
TRUNCATE TABLE lokasi CASCADE;
TRUNCATE TABLE kategori_pendapatan CASCADE;
TRUNCATE TABLE pekerjaan CASCADE;
TRUNCATE TABLE pendidikan CASCADE;
