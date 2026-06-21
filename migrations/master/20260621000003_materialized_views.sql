-- +goose Up
-- +goose StatementBegin

-- Drop old views
DROP VIEW IF EXISTS v_tumbuh_kembang;
DROP VIEW IF EXISTS v_riwayat_pemeriksaan;
DROP VIEW IF EXISTS v_public_stats;
DROP VIEW IF EXISTS v_dashboard_kehadiran_bulanan;
DROP VIEW IF EXISTS v_dashboard_stunting_per_wilayah;
DROP VIEW IF EXISTS v_dashboard_tren_stunting;
DROP VIEW IF EXISTS v_dashboard_distribusi_gizi;
DROP VIEW IF EXISTS v_dashboard_stats;
DROP VIEW IF EXISTS v_dashboard_jadwal_terdekat;

-- ============================================================
-- MATERIALIZED VIEWS (pre-computed, instant reads)
-- ============================================================

CREATE MATERIALIZED VIEW mv_dashboard_stats AS
WITH latest AS (
    SELECT DISTINCT ON (ji.id_pasien)
        hp.status_stunting,
        ji.id_pasien
    FROM hasil_pemeriksaan hp
    JOIN jadwal_imunisasi ji ON hp.id_jadwal_imunisasi = ji.id_imunisasi
    JOIN pasien p ON ji.id_pasien = p.id_pasien
    WHERE p.is_deleted = FALSE
    ORDER BY ji.id_pasien, hp.created_at DESC
)
SELECT
    (SELECT COUNT(*)::int FROM pasien WHERE is_deleted = FALSE) AS total_pasien,
    (SELECT COUNT(*)::int FROM hasil_pemeriksaan WHERE created_at >= CURRENT_DATE - INTERVAL '7 days') AS perlu_verifikasi,
    (SELECT COUNT(*)::int FROM tindak_lanjut WHERE status_pasien != 'Selesai Pemantauan') AS tindak_lanjut,
    (SELECT COUNT(*)::int FROM latest WHERE status_stunting IN ('Stunting', 'Stunting Berat')) AS kasus_stunting,
    (SELECT COUNT(*)::int FROM jadwal_imunisasi WHERE status_imunisasi = 'Belum' AND tanggal_jadwal >= CURRENT_DATE) AS jadwal_posyandu,
    (SELECT COUNT(*)::int FROM pasien p JOIN anak a ON p.id_pasien = a.id_pasien WHERE p.is_deleted = FALSE) AS total_balita,
    (SELECT ROUND(COUNT(*) FILTER (WHERE status_imunisasi = 'Sudah')::numeric * 100.0 / NULLIF(COUNT(*), 0), 1) FROM jadwal_imunisasi) AS cakupan_persentase;

CREATE MATERIALIZED VIEW mv_dashboard_distribusi_gizi AS
WITH latest AS (
    SELECT DISTINCT ON (ji.id_pasien) hp.status_gizi, ji.id_pasien
    FROM hasil_pemeriksaan hp
    JOIN jadwal_imunisasi ji ON hp.id_jadwal_imunisasi = ji.id_imunisasi
    JOIN pasien p ON ji.id_pasien = p.id_pasien
    WHERE p.is_deleted = FALSE
    ORDER BY ji.id_pasien, hp.created_at DESC
)
SELECT status_gizi, COUNT(*)::int AS jumlah FROM latest GROUP BY status_gizi ORDER BY jumlah DESC;

CREATE MATERIALIZED VIEW mv_dashboard_tren_stunting AS
SELECT
    TO_CHAR(hp.created_at, 'YYYY-MM') AS bulan,
    COUNT(DISTINCT ji.id_pasien)::int AS jumlah
FROM hasil_pemeriksaan hp
JOIN jadwal_imunisasi ji ON hp.id_jadwal_imunisasi = ji.id_imunisasi
WHERE hp.status_stunting IN ('Stunting', 'Stunting Berat')
GROUP BY TO_CHAR(hp.created_at, 'YYYY-MM')
ORDER BY bulan;

CREATE MATERIALIZED VIEW mv_dashboard_kehadiran_bulanan AS
SELECT
    TO_CHAR(hp.created_at, 'YYYY-MM') AS bulan,
    COUNT(DISTINCT ji.id_pasien)::int AS jumlah
FROM hasil_pemeriksaan hp
JOIN jadwal_imunisasi ji ON hp.id_jadwal_imunisasi = ji.id_imunisasi
GROUP BY TO_CHAR(hp.created_at, 'YYYY-MM')
ORDER BY bulan;

CREATE MATERIALIZED VIEW mv_dashboard_stunting_per_wilayah AS
WITH latest AS (
    SELECT DISTINCT ON (ji.id_pasien) hp.status_stunting, ji.id_pasien
    FROM hasil_pemeriksaan hp
    JOIN jadwal_imunisasi ji ON hp.id_jadwal_imunisasi = ji.id_imunisasi
    JOIN pasien p ON ji.id_pasien = p.id_pasien
    WHERE p.is_deleted = FALSE
    ORDER BY ji.id_pasien, hp.created_at DESC
),
pasien_lokasi AS (
    SELECT p.id_pasien,
        COALESCE(l_kab.nama_lokasi, 'Tidak Diketahui') AS nama_wilayah
    FROM pasien p
    JOIN posyandu pos ON p.id_posyandu = pos.id_posyandu
    JOIN lokasi l_kel ON pos.id_lokasi = l_kel.id_lokasi
    LEFT JOIN lokasi l_kec ON l_kel.bagian_dari = l_kec.id_lokasi AND l_kec.tipe_lokasi = 'Kecamatan'
    LEFT JOIN lokasi l_kab ON l_kec.bagian_dari = l_kab.id_lokasi AND l_kab.tipe_lokasi = 'Kabupaten'
    WHERE p.is_deleted = FALSE
)
SELECT
    pl.nama_wilayah,
    COUNT(DISTINCT pl.id_pasien)::int AS total_balita,
    COUNT(DISTINCT CASE WHEN l.status_stunting IN ('Stunting', 'Stunting Berat') THEN pl.id_pasien END)::int AS jumlah_kasus,
    ROUND(COUNT(DISTINCT CASE WHEN l.status_stunting IN ('Stunting', 'Stunting Berat') THEN pl.id_pasien END)::numeric * 100.0 / NULLIF(COUNT(DISTINCT pl.id_pasien), 0), 1) AS prevalensi,
    CASE
        WHEN COUNT(DISTINCT pl.id_pasien) > 0
            AND COUNT(DISTINCT CASE WHEN l.status_stunting IN ('Stunting', 'Stunting Berat') THEN pl.id_pasien END)::numeric * 100.0 / COUNT(DISTINCT pl.id_pasien) >= 20 THEN 'tinggi'
        WHEN COUNT(DISTINCT pl.id_pasien) > 0
            AND COUNT(DISTINCT CASE WHEN l.status_stunting IN ('Stunting', 'Stunting Berat') THEN pl.id_pasien END)::numeric * 100.0 / COUNT(DISTINCT pl.id_pasien) >= 10 THEN 'sedang'
        ELSE 'rendah'
    END AS level
FROM pasien_lokasi pl
LEFT JOIN latest l ON pl.id_pasien = l.id_pasien
GROUP BY pl.nama_wilayah
ORDER BY jumlah_kasus DESC;

CREATE MATERIALIZED VIEW mv_public_stats AS
WITH latest AS (
    SELECT DISTINCT ON (ji.id_pasien) hp.status_stunting, ji.id_pasien
    FROM hasil_pemeriksaan hp
    JOIN jadwal_imunisasi ji ON hp.id_jadwal_imunisasi = ji.id_imunisasi
    JOIN pasien p ON ji.id_pasien = p.id_pasien
    WHERE p.is_deleted = FALSE
    ORDER BY ji.id_pasien, hp.created_at DESC
)
SELECT
    (SELECT COUNT(*)::int FROM pasien WHERE is_deleted = FALSE) AS total_pasien,
    (SELECT COUNT(*)::int FROM pasien p JOIN anak a ON p.id_pasien = a.id_pasien WHERE p.is_deleted = FALSE) AS balita_dipantau,
    (SELECT COUNT(*)::int FROM latest WHERE status_stunting IN ('Stunting', 'Stunting Berat')) AS kasus_stunting,
    (SELECT COUNT(*)::int FROM artikel WHERE status_artikel = 'Dipublikasikan') AS total_artikel;

CREATE MATERIALIZED VIEW mv_riwayat_pemeriksaan AS
SELECT
    ji.id_pasien,
    hp.created_at::text AS tanggal,
    hp.berat_badan::numeric AS berat_badan,
    hp.tinggi_badan::numeric AS tinggi_badan,
    hp.status_gizi,
    hp.catatan,
    ua.nama AS petugas
FROM hasil_pemeriksaan hp
JOIN jadwal_imunisasi ji ON hp.id_jadwal_imunisasi = ji.id_imunisasi
JOIN user_account ua ON hp.id_petugas_input = ua.id_user
ORDER BY ji.id_pasien, hp.created_at DESC;

CREATE MATERIALIZED VIEW mv_tumbuh_kembang AS
SELECT DISTINCT ON (ji.id_pasien, TO_CHAR(hp.created_at, 'YYYY-MM'))
    ji.id_pasien,
    TO_CHAR(hp.created_at, 'YYYY-MM') AS bulan,
    hp.berat_badan::numeric AS berat_badan,
    hp.tinggi_badan::numeric AS tinggi_badan
FROM hasil_pemeriksaan hp
JOIN jadwal_imunisasi ji ON hp.id_jadwal_imunisasi = ji.id_imunisasi
ORDER BY ji.id_pasien, TO_CHAR(hp.created_at, 'YYYY-MM'), hp.created_at DESC;

CREATE MATERIALIZED VIEW mv_dashboard_jadwal_terdekat AS
SELECT
    ji.id_imunisasi AS id,
    ji.nama_vaksin,
    ji.tanggal_jadwal::text,
    COALESCE(a.nama_anak, ua.nama) AS nama_pasien
FROM jadwal_imunisasi ji
JOIN pasien p ON ji.id_pasien = p.id_pasien
JOIN user_account ua ON p.id_pasien = ua.id_user
LEFT JOIN anak a ON p.id_pasien = a.id_pasien
WHERE ji.status_imunisasi = 'Belum'
  AND ji.tanggal_jadwal >= CURRENT_DATE
  AND p.is_deleted = FALSE
ORDER BY ji.tanggal_jadwal ASC
LIMIT 10;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP MATERIALIZED VIEW IF EXISTS mv_dashboard_jadwal_terdekat;
DROP MATERIALIZED VIEW IF EXISTS mv_tumbuh_kembang;
DROP MATERIALIZED VIEW IF EXISTS mv_riwayat_pemeriksaan;
DROP MATERIALIZED VIEW IF EXISTS mv_public_stats;
DROP MATERIALIZED VIEW IF EXISTS mv_dashboard_stunting_per_wilayah;
DROP MATERIALIZED VIEW IF EXISTS mv_dashboard_kehadiran_bulanan;
DROP MATERIALIZED VIEW IF EXISTS mv_dashboard_tren_stunting;
DROP MATERIALIZED VIEW IF EXISTS mv_dashboard_distribusi_gizi;
DROP MATERIALIZED VIEW IF EXISTS mv_dashboard_stats;

-- Restore regular views
CREATE VIEW v_dashboard_stats AS SELECT 0::int AS total_pasien, 0::int AS perlu_verifikasi, 0::int AS tindak_lanjut, 0::int AS kasus_stunting, 0::int AS jadwal_posyandu, 0::int AS total_balita, 0.0 AS cakupan_persentase;
CREATE VIEW v_dashboard_distribusi_gizi AS SELECT ''::varchar AS status_gizi, 0::int AS jumlah WHERE FALSE;
CREATE VIEW v_dashboard_tren_stunting AS SELECT ''::varchar AS bulan, 0::int AS jumlah WHERE FALSE;
CREATE VIEW v_dashboard_kehadiran_bulanan AS SELECT ''::varchar AS bulan, 0::int AS jumlah WHERE FALSE;
CREATE VIEW v_dashboard_stunting_per_wilayah AS SELECT ''::varchar AS nama_wilayah, 0::int AS total_balita, 0::int AS jumlah_kasus, 0.0 AS prevalensi, ''::varchar AS level WHERE FALSE;
CREATE VIEW v_public_stats AS SELECT 0::int AS total_pasien, 0::int AS balita_dipantau, 0::int AS kasus_stunting, 0::int AS total_artikel;
CREATE VIEW v_riwayat_pemeriksaan AS SELECT 0::int AS id_pasien, ''::text AS tanggal, 0.0::numeric AS berat_badan, 0.0::numeric AS tinggi_badan, ''::varchar AS status_gizi, NULL::varchar AS catatan, ''::varchar AS petugas WHERE FALSE;
CREATE VIEW v_tumbuh_kembang AS SELECT 0::int AS id_pasien, ''::varchar AS bulan, 0.0::numeric AS berat_badan, 0.0::numeric AS tinggi_badan WHERE FALSE;
CREATE VIEW v_dashboard_jadwal_terdekat AS SELECT 0::int AS id, ''::varchar AS nama_vaksin, ''::text AS tanggal_jadwal, ''::varchar AS nama_pasien WHERE FALSE;

-- +goose StatementEnd
